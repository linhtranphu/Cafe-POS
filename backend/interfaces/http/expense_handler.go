package http

import (
	"errors"
	"net/http"
	"time"
	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/expense"
	"cafe-pos/backend/domain/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExpenseHandler struct {
	service                       *services.ExpenseService
	fundExpenseIntegrationService *services.FundExpenseIntegrationService
}

func NewExpenseHandler(service *services.ExpenseService, fundExpenseIntegrationService *services.FundExpenseIntegrationService) *ExpenseHandler {
	return &ExpenseHandler{
		service:                       service,
		fundExpenseIntegrationService: fundExpenseIntegrationService,
	}
}

func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var req struct {
		Date          string  `json:"date"`
		CategoryID    string  `json:"category_id"`
		Amount        float64 `json:"amount"`
		Description   string  `json:"description"`
		PaymentMethod string  `json:"payment_method"`
		Vendor        string  `json:"vendor,omitempty"`
		Notes         string  `json:"notes,omitempty"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Parse date - support both ISO format and YYYY-MM-DD
	var date time.Time
	var err error
	
	// Try ISO format first (YYYY-MM-DDTHH:MM:SSZ)
	date, err = time.Parse(time.RFC3339, req.Date)
	if err != nil {
		// Fallback to YYYY-MM-DD format
		date, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use ISO format (YYYY-MM-DDTHH:MM:SSZ) or YYYY-MM-DD"})
			return
		}
	}
	
	// Parse category ID
	categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id"})
		return
	}
	
	// Get username from context (set by auth middleware)
	username, _ := c.Get("username")
	createdBy := ""
	if u, ok := username.(string); ok {
		createdBy = u
	}
	
	e := expense.Expense{
		Date:          date,
		CategoryID:    categoryID,
		Amount:        req.Amount,
		Description:   req.Description,
		PaymentMethod: req.PaymentMethod,
		Vendor:        req.Vendor,
		Notes:         req.Notes,
		CreatedBy:     createdBy,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	
	if err := h.service.CreateExpense(c.Request.Context(), &e); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, e)
}

// CreateExpenseFromFund creates an expense paid from fund
// Requirements: 1.1, 1.2, 1.3, 6.4
func (h *ExpenseHandler) CreateExpenseFromFund(c *gin.Context) {
	var req struct {
		Date        string  `json:"date" binding:"required"`
		CategoryID  string  `json:"category_id" binding:"required"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description" binding:"required"`
		Vendor      string  `json:"vendor,omitempty"`
		Notes       string  `json:"notes,omitempty"`
		MoneyType   string  `json:"money_type,omitempty"` // "cash" or "transfer"
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Parse date - support both ISO format and YYYY-MM-DD
	var date time.Time
	var err error
	
	// Try ISO format first (YYYY-MM-DDTHH:MM:SSZ)
	date, err = time.Parse(time.RFC3339, req.Date)
	if err != nil {
		// Fallback to YYYY-MM-DD format
		date, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use ISO format (YYYY-MM-DDTHH:MM:SSZ) or YYYY-MM-DD"})
			return
		}
	}
	
	// Parse category ID
	categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id"})
		return
	}
	
	// Extract user info from JWT token (set by auth middleware)
	username, _ := c.Get("username")
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("role")
	
	userName := ""
	if u, ok := username.(string); ok {
		userName = u
	}
	
	userIDObj := primitive.NilObjectID
	if uid, ok := userID.(string); ok {
		if oid, err := primitive.ObjectIDFromHex(uid); err == nil {
			userIDObj = oid
		}
	}
	
	// Role is stored as user.Role type in context, convert to string
	role := ""
	if r, ok := userRole.(user.Role); ok {
		role = string(r)
	}
	
	// Create request for service
	serviceReq := services.CreateExpenseFromFundRequest{
		Date:        date,
		CategoryID:  categoryID,
		Amount:      req.Amount,
		Description: req.Description,
		Vendor:      req.Vendor,
		Notes:       req.Notes,
		MoneyType:   req.MoneyType,
		UserID:      userIDObj,
		UserName:    userName,
		UserRole:    role,
	}
	
	// Call service to create expense from fund
	// Requirements 1.2, 1.3: Validate balance and create withdrawal
	result, err := h.fundExpenseIntegrationService.CreateExpenseFromFund(c.Request.Context(), serviceReq)
	if err != nil {
		if errors.Is(err, services.ErrInsufficientFundBalance) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Return expense and fund transaction details
	c.JSON(http.StatusCreated, gin.H{
		"expense":       result.Expense,
		"journal_entry": result.JournalEntry,
	})
}

func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
	filter := bson.M{}
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filter["date"] = bson.M{"$gte": t}
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			if filter["date"] != nil {
				filter["date"].(bson.M)["$lte"] = t
			} else {
				filter["date"] = bson.M{"$lte": t}
			}
		}
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, err := primitive.ObjectIDFromHex(categoryID); err == nil {
			filter["category_id"] = id
		}
	}
	
	// Requirement 4.3: Add paid_from_fund filter
	if paidFromFund := c.Query("paid_from_fund"); paidFromFund != "" {
		if paidFromFund == "true" {
			filter["paid_from_fund"] = true
		} else if paidFromFund == "false" {
			filter["paid_from_fund"] = false
		}
	}
	
	expenses, err := h.service.GetExpenses(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}

// GetExpenseByID retrieves a single expense by ID
// Requirement 4.2: Include fund_transaction_id in response
func (h *ExpenseHandler) GetExpenseByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}
	
	// Get expense by ID
	expenses, err := h.service.GetExpenses(c.Request.Context(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	if len(expenses) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}
	
	expense := expenses[0]
	
	// Requirement 4.2: Include fund transaction details if expense was paid from fund
	response := gin.H{
		"expense": expense,
	}
	
	// If expense was paid from fund, include fund transaction details
	if expense.PaidFromFund && !expense.FundTransactionID.IsZero() {
		response["has_fund_transaction"] = true
		response["fund_transaction_id"] = expense.FundTransactionID.Hex()
	}
	
	c.JSON(http.StatusOK, response)
}

func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}
	
	var req struct {
		Date          string  `json:"date"`
		CategoryID    string  `json:"category_id"`
		Amount        float64 `json:"amount"`
		Description   string  `json:"description"`
		PaymentMethod string  `json:"payment_method"`
		Vendor        string  `json:"vendor,omitempty"`
		Notes         string  `json:"notes,omitempty"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get existing expense to preserve CreatedBy and CreatedAt
	existingExpenses, err := h.service.GetExpenses(c.Request.Context(), bson.M{"_id": id})
	if err != nil || len(existingExpenses) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}
	existingExpense := existingExpenses[0]
	
	// Parse date - support both ISO format and YYYY-MM-DD
	var date time.Time
	
	// Try ISO format first (YYYY-MM-DDTHH:MM:SSZ)
	date, err = time.Parse(time.RFC3339, req.Date)
	if err != nil {
		// Fallback to YYYY-MM-DD format
		date, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use ISO format (YYYY-MM-DDTHH:MM:SSZ) or YYYY-MM-DD"})
			return
		}
	}
	
	// Parse category ID
	categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id"})
		return
	}
	
	// Preserve CreatedBy, CreatedAt, SourceType, and SourceID from existing expense
	e := expense.Expense{
		ID:            id,
		Date:          date,
		CategoryID:    categoryID,
		Amount:        req.Amount,
		Description:   req.Description,
		PaymentMethod: req.PaymentMethod,
		Vendor:        req.Vendor,
		Notes:         req.Notes,
		CreatedBy:     existingExpense.CreatedBy,     // Preserve original creator
		CreatedAt:     existingExpense.CreatedAt,     // Preserve creation time
		SourceType:    existingExpense.SourceType,    // Preserve source type
		SourceID:      existingExpense.SourceID,      // Preserve source ID
		UpdatedAt:     time.Now(),
	}
	
	if err := h.service.UpdateExpense(c.Request.Context(), id, &e); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, e)
}

func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	if err := h.service.DeleteExpense(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ExpenseHandler) CreateCategory(c *gin.Context) {
	var cat expense.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateCategory(c.Request.Context(), &cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func (h *ExpenseHandler) GetCategories(c *gin.Context) {
	categoryType := c.Query("type") // "" = general, "operating" = operating
	categories, err := h.service.GetCategories(c.Request.Context(), categoryType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (h *ExpenseHandler) DeleteCategory(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	if err := h.service.DeleteCategory(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ExpenseHandler) CreateRecurring(c *gin.Context) {
	var re expense.RecurringExpense
	if err := c.ShouldBindJSON(&re); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateRecurring(c.Request.Context(), &re); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, re)
}

func (h *ExpenseHandler) GetRecurring(c *gin.Context) {
	recurring, err := h.service.GetRecurring(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recurring)
}

func (h *ExpenseHandler) DeleteRecurring(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	if err := h.service.DeleteRecurring(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ExpenseHandler) CreatePrepaid(c *gin.Context) {
	var pe expense.PrepaidExpense
	if err := c.ShouldBindJSON(&pe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreatePrepaid(c.Request.Context(), &pe); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pe)
}

func (h *ExpenseHandler) GetPrepaid(c *gin.Context) {
	prepaid, err := h.service.GetPrepaid(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prepaid)
}

func (h *ExpenseHandler) DeletePrepaid(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	if err := h.service.DeletePrepaid(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
