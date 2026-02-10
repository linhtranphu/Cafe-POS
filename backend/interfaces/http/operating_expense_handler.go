package http

import (
	"net/http"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/expense"

	"github.com/gin-gonic/gin"
)

// OperatingExpenseHandler handles HTTP requests for operating expense management
type OperatingExpenseHandler struct {
	operatingExpenseService *services.OperatingExpenseService
}

// NewOperatingExpenseHandler creates a new operating expense handler
func NewOperatingExpenseHandler(operatingExpenseService *services.OperatingExpenseService) *OperatingExpenseHandler {
	return &OperatingExpenseHandler{
		operatingExpenseService: operatingExpenseService,
	}
}

// CreateOperatingExpense handles POST /api/operating-expenses
// Creates or updates an operating expense for a period
// Requirements: 6.5.2, 6.5.7
func (h *OperatingExpenseHandler) CreateOperatingExpense(c *gin.Context) {
	var req expense.OperatingExpenseRequest
	
	// Bind and validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	// Call service to upsert operating expense
	operatingExpense, err := h.operatingExpenseService.UpsertOperatingExpense(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create/update operating expense",
			"details": err.Error(),
		})
		return
	}
	
	// Return created/updated expense
	c.JSON(http.StatusOK, operatingExpense)
}

// GetOperatingExpenses handles GET /api/operating-expenses
// Retrieves operating expenses with optional date range filter
// Requirements: 6.5.7
func (h *OperatingExpenseHandler) GetOperatingExpenses(c *gin.Context) {
	// Parse optional query parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	
	var startDate, endDate *time.Time
	
	// Parse start_date if provided
	if startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid start_date format. Expected YYYY-MM-DD",
				"details": err.Error(),
			})
			return
		}
		startDate = &parsed
	}
	
	// Parse end_date if provided
	if endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid end_date format. Expected YYYY-MM-DD",
				"details": err.Error(),
			})
			return
		}
		endDate = &parsed
	}
	
	// Call service to get operating expenses
	expenses, err := h.operatingExpenseService.GetOperatingExpenses(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve operating expenses",
			"details": err.Error(),
		})
		return
	}
	
	// Return expenses array
	c.JSON(http.StatusOK, gin.H{
		"expenses": expenses,
	})
}
