package http

import (
	"net/http"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/batch"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BatchRecordHandler handles batch record HTTP endpoints
type BatchRecordHandler struct {
	batchRecordService *services.BatchRecordService
}

// NewBatchRecordHandler creates a new batch record handler
func NewBatchRecordHandler(batchRecordService *services.BatchRecordService) *BatchRecordHandler {
	return &BatchRecordHandler{
		batchRecordService: batchRecordService,
	}
}

// CreateBatchRecordRequest represents the request body for creating a batch record
type CreateBatchRecordRequest struct {
	BatchDefinitionID string  `json:"batch_definition_id" binding:"required"`
	QuantityProduced  float64 `json:"quantity_produced" binding:"required,gt=0"`
	PreparedBy        string  `json:"prepared_by" binding:"required"`
}

// CreateBatchRecord handles POST /api/batch-records
func (h *BatchRecordHandler) CreateBatchRecord(c *gin.Context) {
	var req CreateBatchRecordRequest
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

	// Create batch record
	createReq := services.CreateBatchRequest{
		BatchDefinitionID: batchDefID,
		QuantityProduced:  req.QuantityProduced,
		PreparedBy:        req.PreparedBy,
	}

	result, err := h.batchRecordService.CreateBatch(c.Request.Context(), createReq)
	if err != nil {
		// Check for specific error types
		if err.Error() == "insufficient ingredients" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetBatchRecords handles GET /api/batch-records
func (h *BatchRecordHandler) GetBatchRecords(c *gin.Context) {
	// Parse query parameters
	batchDefID := c.Query("batch_definition_id")
	status := c.Query("status")
	preparedBy := c.Query("prepared_by")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	// Create filter
	filter := batch.BatchRecordFilter{
		Status:     status,
		PreparedBy: preparedBy,
	}

	// Parse batch definition ID if provided
	if batchDefID != "" {
		id, err := primitive.ObjectIDFromHex(batchDefID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch_definition_id"})
			return
		}
		filter.BatchDefinitionID = &id
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

	// Get batch records
	records, err := h.batchRecordService.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Simple pagination (in production, use proper pagination)
	c.JSON(http.StatusOK, gin.H{
		"data":  records,
		"total": len(records),
		"page":  page,
		"limit": limit,
	})
}

// GetBatchRecord handles GET /api/batch-records/:id
func (h *BatchRecordHandler) GetBatchRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	record, err := h.batchRecordService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "batch record not found"})
		return
	}

	c.JSON(http.StatusOK, record)
}

// UpdateQuantityRequest represents the request body for updating batch quantity
type UpdateQuantityRequest struct {
	QuantityRemaining float64 `json:"quantity_remaining" binding:"required,min=0"`
}

// UpdateBatchQuantity handles PATCH /api/batch-records/:id/quantity
func (h *BatchRecordHandler) UpdateBatchQuantity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateQuantityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.batchRecordService.UpdateQuantity(c.Request.Context(), id, req.QuantityRemaining)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get updated record
	record, err := h.batchRecordService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

// MarkBatchExpired handles PATCH /api/batch-records/:id/expire
func (h *BatchRecordHandler) MarkBatchExpired(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.batchRecordService.MarkAsExpired(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get updated record
	record, err := h.batchRecordService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

// DeleteBatchRecord handles DELETE /api/batch-records/:id
func (h *BatchRecordHandler) DeleteBatchRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.batchRecordService.Delete(c.Request.Context(), id)
	if err != nil {
		// Check for specific error types
		if err.Error() == "cannot delete batch that has been partially used" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
