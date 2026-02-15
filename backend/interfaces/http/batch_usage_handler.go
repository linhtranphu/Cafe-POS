package http

import (
	"net/http"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/batch"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BatchUsageHandler handles batch usage HTTP endpoints
type BatchUsageHandler struct {
	batchUsageService *services.BatchUsageService
}

// NewBatchUsageHandler creates a new batch usage handler
func NewBatchUsageHandler(batchUsageService *services.BatchUsageService) *BatchUsageHandler {
	return &BatchUsageHandler{
		batchUsageService: batchUsageService,
	}
}

// UseBatchRequest represents the request body for using batch
type UseBatchRequest struct {
	BatchDefinitionID string `json:"batch_definition_id" binding:"required"`
	QuantityNeeded    float64 `json:"quantity_needed" binding:"required,gt=0"`
	OrderID           string `json:"order_id" binding:"required"`
	MenuItemID        string `json:"menu_item_id" binding:"required"`
	MenuItemName      string `json:"menu_item_name" binding:"required"`
}

// UseBatchResponse represents the response for batch usage
type UseBatchResponse struct {
	Success     bool                      `json:"success"`
	BatchesUsed []BatchUsageDetailResponse `json:"batches_used"`
	TotalCost   float64                   `json:"total_cost"`
	Message     string                    `json:"message"`
}

// BatchUsageDetailResponse represents details of a single batch usage in response
type BatchUsageDetailResponse struct {
	BatchRecordID string  `json:"batch_record_id"`
	QuantityUsed  float64 `json:"quantity_used"`
	CostPerUnit   float64 `json:"cost_per_unit"`
}

// UseBatch handles POST /api/batch-usage
func (h *BatchUsageHandler) UseBatch(c *gin.Context) {
	var req UseBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse batch definition ID
	batchDefID, err := primitive.ObjectIDFromHex(req.BatchDefinitionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch_definition_id"})
		return
	}

	// Parse order ID
	orderID, err := primitive.ObjectIDFromHex(req.OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order_id"})
		return
	}

	// Parse menu item ID
	menuItemID, err := primitive.ObjectIDFromHex(req.MenuItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu_item_id"})
		return
	}

	// Create service request
	serviceReq := services.UseBatchRequest{
		BatchDefinitionID: batchDefID,
		QuantityNeeded:    req.QuantityNeeded,
		OrderID:           orderID,
		MenuItemID:        menuItemID,
		MenuItemName:      req.MenuItemName,
	}

	// Use batch
	result, err := h.batchUsageService.UseBatch(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert result to response format
	batchesUsed := make([]BatchUsageDetailResponse, len(result.BatchesUsed))
	for i, bu := range result.BatchesUsed {
		batchesUsed[i] = BatchUsageDetailResponse{
			BatchRecordID: bu.BatchRecordID.Hex(),
			QuantityUsed:  bu.QuantityUsed,
			CostPerUnit:   bu.CostPerUnit,
		}
	}

	response := UseBatchResponse{
		Success:     result.Success,
		BatchesUsed: batchesUsed,
		TotalCost:   result.TotalCost,
		Message:     result.Message,
	}

	// Return 200 even if not successful (business logic failure, not server error)
	c.JSON(http.StatusOK, response)
}

// GetUsageHistory handles GET /api/batch-usage/history
func (h *BatchUsageHandler) GetUsageHistory(c *gin.Context) {
	// Parse query parameters
	batchRecordID := c.Query("batch_record_id")
	orderID := c.Query("order_id")
	menuItemID := c.Query("menu_item_id")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	// Create filter
	filter := batch.BatchUsageLogFilter{}

	// Parse batch record ID if provided
	if batchRecordID != "" {
		id, err := primitive.ObjectIDFromHex(batchRecordID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch_record_id"})
			return
		}
		filter.BatchRecordID = &id
	}

	// Parse order ID if provided
	if orderID != "" {
		id, err := primitive.ObjectIDFromHex(orderID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order_id"})
			return
		}
		filter.OrderID = &id
	}

	// Parse menu item ID if provided
	if menuItemID != "" {
		id, err := primitive.ObjectIDFromHex(menuItemID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu_item_id"})
			return
		}
		filter.MenuItemID = &id
	}

	// Parse dates if provided
	if fromDate != "" {
		t, err := time.Parse(time.RFC3339, fromDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from_date format"})
			return
		}
		filter.FromDate = &t
	}

	if toDate != "" {
		t, err := time.Parse(time.RFC3339, toDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to_date format"})
			return
		}
		filter.ToDate = &t
	}

	// Get usage history
	usageLogs, err := h.batchUsageService.GetUsageHistory(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Simple pagination (in production, use proper pagination)
	c.JSON(http.StatusOK, gin.H{
		"data":  usageLogs,
		"total": len(usageLogs),
		"page":  page,
		"limit": limit,
	})
}
