package http

import (
	"net/http"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/batch"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BatchDefinitionHandler handles batch definition HTTP endpoints
type BatchDefinitionHandler struct {
	batchDefinitionService *services.BatchDefinitionService
}

// NewBatchDefinitionHandler creates a new batch definition handler
func NewBatchDefinitionHandler(batchDefinitionService *services.BatchDefinitionService) *BatchDefinitionHandler {
	return &BatchDefinitionHandler{
		batchDefinitionService: batchDefinitionService,
	}
}

// CreateBatchDefinitionRequest represents the request body for creating a batch definition
type CreateBatchDefinitionRequest struct {
	Name                string                      `json:"name" binding:"required"`
	Unit                string                      `json:"unit" binding:"required"`
	ShelfLifeHours      int                         `json:"shelf_life_hours" binding:"required,min=1"`
	ConversionRates     []ConversionRateRequest     `json:"conversion_rates" binding:"required,min=1"`
	LowStockThreshold   float64                     `json:"low_stock_threshold" binding:"min=0"`
	ExpiryWarningHours  int                         `json:"expiry_warning_hours" binding:"min=0"`
}

// ConversionRateRequest represents a conversion rate in the request
type ConversionRateRequest struct {
	SourceIngredientID   string  `json:"source_ingredient_id" binding:"required"`
	SourceIngredientName string  `json:"source_ingredient_name" binding:"required"`
	SourceQuantity       float64 `json:"source_quantity" binding:"required,gt=0"`
	SourceUnit           string  `json:"source_unit" binding:"required"`
	BatchQuantity        float64 `json:"batch_quantity" binding:"required,gt=0"`
	WastageRate          float64 `json:"wastage_rate" binding:"min=0,max=1"`
}

// CreateBatchDefinition handles POST /api/batch-definitions
func (h *BatchDefinitionHandler) CreateBatchDefinition(c *gin.Context) {
	var req CreateBatchDefinitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert request to domain model
	conversionRates := make([]batch.ConversionRate, len(req.ConversionRates))
	for i, cr := range req.ConversionRates {
		sourceID, err := primitive.ObjectIDFromHex(cr.SourceIngredientID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source_ingredient_id"})
			return
		}

		conversionRates[i] = batch.ConversionRate{
			SourceIngredientID:   sourceID,
			SourceIngredientName: cr.SourceIngredientName,
			SourceQuantity:       cr.SourceQuantity,
			SourceUnit:           cr.SourceUnit,
			BatchQuantity:        cr.BatchQuantity,
			WastageRate:          cr.WastageRate,
		}
	}

	createReq := &batch.CreateBatchDefinitionRequest{
		Name:               req.Name,
		Unit:               req.Unit,
		ShelfLifeHours:     req.ShelfLifeHours,
		ConversionRates:    conversionRates,
		LowStockThreshold:  req.LowStockThreshold,
		ExpiryWarningHours: req.ExpiryWarningHours,
	}

	// Create batch definition
	result, err := h.batchDefinitionService.Create(c.Request.Context(), createReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetBatchDefinitions handles GET /api/batch-definitions
func (h *BatchDefinitionHandler) GetBatchDefinitions(c *gin.Context) {
	// Parse query parameters
	search := c.Query("search")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	// Create filter
	filter := batch.BatchDefinitionFilter{
		Search: search,
	}

	// Get batch definitions
	definitions, total, err := h.batchDefinitionService.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Simple pagination (in production, use proper pagination)
	c.JSON(http.StatusOK, gin.H{
		"data":  definitions,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetBatchDefinition handles GET /api/batch-definitions/:id
func (h *BatchDefinitionHandler) GetBatchDefinition(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	definition, err := h.batchDefinitionService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "batch definition not found"})
		return
	}

	c.JSON(http.StatusOK, definition)
}

// UpdateBatchDefinition handles PUT /api/batch-definitions/:id
func (h *BatchDefinitionHandler) UpdateBatchDefinition(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req CreateBatchDefinitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert request to domain model
	conversionRates := make([]batch.ConversionRate, len(req.ConversionRates))
	for i, cr := range req.ConversionRates {
		sourceID, err := primitive.ObjectIDFromHex(cr.SourceIngredientID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source_ingredient_id"})
			return
		}

		conversionRates[i] = batch.ConversionRate{
			SourceIngredientID:   sourceID,
			SourceIngredientName: cr.SourceIngredientName,
			SourceQuantity:       cr.SourceQuantity,
			SourceUnit:           cr.SourceUnit,
			BatchQuantity:        cr.BatchQuantity,
			WastageRate:          cr.WastageRate,
		}
	}

	updateReq := &batch.UpdateBatchDefinitionRequest{
		Name:               req.Name,
		Unit:               req.Unit,
		ShelfLifeHours:     &req.ShelfLifeHours,
		ConversionRates:    conversionRates,
		LowStockThreshold:  &req.LowStockThreshold,
		ExpiryWarningHours: &req.ExpiryWarningHours,
	}

	// Update batch definition
	result, err := h.batchDefinitionService.Update(c.Request.Context(), id, updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteBatchDefinition handles DELETE /api/batch-definitions/:id
func (h *BatchDefinitionHandler) DeleteBatchDefinition(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.batchDefinitionService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
