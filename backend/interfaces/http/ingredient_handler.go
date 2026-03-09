package http

import (
	"fmt"
	"net/http"
	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IngredientHandler struct {
	ingredientService             *services.IngredientService
	fundExpenseIntegrationService *services.FundExpenseIntegrationService
}

func NewIngredientHandler(ingredientService *services.IngredientService) *IngredientHandler {
	return &IngredientHandler{
		ingredientService:             ingredientService,
		fundExpenseIntegrationService: nil, // Will be set via SetFundExpenseIntegrationService
	}
}

// SetFundExpenseIntegrationService sets the fund expense integration service
func (h *IngredientHandler) SetFundExpenseIntegrationService(service *services.FundExpenseIntegrationService) {
	h.fundExpenseIntegrationService = service
}

func (h *IngredientHandler) CreateIngredient(c *gin.Context) {
	var req ingredient.CreateIngredientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user info from context
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	
	userIDStr := ""
	if uid, ok := userID.(string); ok {
		userIDStr = uid
	}
	
	createdBy := ""
	if u, ok := username.(string); ok {
		createdBy = u
	}

	item, err := h.ingredientService.CreateIngredient(c.Request.Context(), &req, userIDStr, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *IngredientHandler) GetAllIngredients(c *gin.Context) {
	items, err := h.ingredientService.GetAllIngredients(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *IngredientHandler) GetIngredient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	item, err := h.ingredientService.GetIngredient(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ingredient not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *IngredientHandler) UpdateIngredient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req ingredient.UpdateIngredientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.ingredientService.UpdateIngredient(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *IngredientHandler) DeleteIngredient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.ingredientService.DeleteIngredient(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ingredient deleted"})
}

func (h *IngredientHandler) AdjustStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req ingredient.StockAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user info from context
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	req.UserID = userID.(string)
	req.Username = username.(string)

	item, err := h.ingredientService.AdjustStock(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// StockIn - Add stock (purchase/receive)
func (h *IngredientHandler) StockIn(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req ingredient.StockInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user info from context
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	req.UserID = userID.(string)
	req.Username = username.(string)

	item, err := h.ingredientService.StockIn(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// StockOut - Remove stock (usage/waste)
func (h *IngredientHandler) StockOut(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req ingredient.StockOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user info from context
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	req.UserID = userID.(string)
	req.Username = username.(string)

	item, err := h.ingredientService.StockOut(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// StockAdjust - Set stock to specific quantity (inventory correction)
func (h *IngredientHandler) StockAdjust(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req ingredient.StockAdjustRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user info from context
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	req.UserID = userID.(string)
	req.Username = username.(string)

	item, err := h.ingredientService.StockAdjust(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *IngredientHandler) GetLowStock(c *gin.Context) {
	items, err := h.ingredientService.GetLowStockIngredients(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *IngredientHandler) GetStockHistory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	histories, err := h.ingredientService.GetStockHistory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, histories)
}

// Category handlers
func (h *IngredientHandler) CreateCategory(c *gin.Context) {
	var req ingredient.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat, err := h.ingredientService.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

func (h *IngredientHandler) GetCategories(c *gin.Context) {
	categories, err := h.ingredientService.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *IngredientHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.ingredientService.DeleteCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete category that is in use"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})
}

// RestockFromFund handles POST /api/ingredients/:id/restock/from-fund
// Requirements: 2.1, 2.2, 2.3, 2.4, 6.5
func (h *IngredientHandler) RestockFromFund(c *gin.Context) {
	// Check if fund expense integration service is available
	if h.fundExpenseIntegrationService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "fund expense integration service not available"})
		return
	}

	// Parse ingredient ID
	idStr := c.Param("id")
	ingredientID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ingredient id"})
		return
	}

	// Parse request body
	var req struct {
		Quantity    float64 `json:"quantity" binding:"required,gt=0"`
		CostPerUnit float64 `json:"cost_per_unit" binding:"required,gte=0"`
		Reason      string  `json:"reason" binding:"required"`
		MoneyType   string  `json:"money_type"` // "cash" or "transfer"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("ERROR: Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Default to cash if not specified
	if req.MoneyType == "" {
		req.MoneyType = "cash"
	}
	
	// Validate money_type
	if req.MoneyType != "cash" && req.MoneyType != "transfer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "money_type must be 'cash' or 'transfer'"})
		return
	}
	
	fmt.Printf("DEBUG: Restock request - Ingredient: %s, Quantity: %.2f, Cost: %.2f, MoneyType: %s\n", 
		idStr, req.Quantity, req.CostPerUnit, req.MoneyType)

	// Extract user info from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userName, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found"})
		return
	}

	userRole, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user role not found"})
		return
	}

	// Convert user ID to ObjectID
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id format"})
		return
	}

	userObjID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}
	
	// Convert userRole to string
	var userRoleStr string
	if roleType, ok := userRole.(user.Role); ok {
		userRoleStr = string(roleType)
	} else if roleStr, ok := userRole.(string); ok {
		userRoleStr = roleStr
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user role format"})
		return
	}

	// Call FundExpenseIntegrationService.RestockIngredientFromFund
	restockReq := services.RestockFromFundRequest{
		IngredientID: ingredientID,
		Quantity:     req.Quantity,
		CostPerUnit:  req.CostPerUnit,
		Reason:       req.Reason,
		MoneyType:    req.MoneyType,
		UserID:       userObjID,
		UserName:     userName.(string),
		UserRole:     userRoleStr,
	}

	result, err := h.fundExpenseIntegrationService.RestockIngredientFromFund(c.Request.Context(), restockReq)
	if err != nil {
		fmt.Printf("ERROR: RestockIngredientFromFund failed: %v\n", err)
		// Handle specific errors
		if err == services.ErrInsufficientFundBalance {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return restock record, expense, and fund transaction details
	c.JSON(http.StatusOK, gin.H{
		"restock_record":   result.RestockRecord,
		"expense":          result.Expense,
		"fund_transaction": result.FundTransaction,
		"updated_stock":    result.UpdatedStock,
	})
}

// GetRestockHistory handles GET /api/ingredients/:id/restock-history
// Requirements: 2.5, 2.6, 2.7
func (h *IngredientHandler) GetRestockHistory(c *gin.Context) {
	// Check if fund expense integration service is available
	if h.fundExpenseIntegrationService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "fund expense integration service not available"})
		return
	}

	// Parse ingredient ID
	idStr := c.Param("id")
	ingredientID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ingredient id"})
		return
	}

	// Parse pagination parameters
	limit := 50 // default limit
	offset := 0 // default offset

	if limitStr := c.Query("limit"); limitStr != "" {
		var l int
		if _, err := fmt.Sscanf(limitStr, "%d", &l); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		var o int
		if _, err := fmt.Sscanf(offsetStr, "%d", &o); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get restock history from service
	records, err := h.fundExpenseIntegrationService.GetRestockHistory(c.Request.Context(), ingredientID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return restock history with fund transaction info
	c.JSON(http.StatusOK, gin.H{
		"records": records,
		"limit":   limit,
		"offset":  offset,
	})
}
