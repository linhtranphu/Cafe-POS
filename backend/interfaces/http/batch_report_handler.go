package http

import (
	"net/http"
	"time"

	"cafe-pos/backend/application/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BatchReportHandler handles batch report HTTP endpoints
type BatchReportHandler struct {
	batchReportService *services.BatchReportService
}

// NewBatchReportHandler creates a new batch report handler
func NewBatchReportHandler(batchReportService *services.BatchReportService) *BatchReportHandler {
	return &BatchReportHandler{
		batchReportService: batchReportService,
	}
}

// GetProductionReport handles GET /api/batch-reports/production
func (h *BatchReportHandler) GetProductionReport(c *gin.Context) {
	// Parse required query parameters
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	if fromDateStr == "" || toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from_date and to_date are required"})
		return
	}

	// Parse dates
	fromDate, err := time.Parse(time.RFC3339, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from_date format, use RFC3339"})
		return
	}

	toDate, err := time.Parse(time.RFC3339, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to_date format, use RFC3339"})
		return
	}

	// Parse optional parameters
	params := services.ProductionReportParams{
		FromDate:   fromDate,
		ToDate:     toDate,
		PreparedBy: c.Query("prepared_by"),
	}

	// Parse batch definition ID if provided
	batchDefIDStr := c.Query("batch_definition_id")
	if batchDefIDStr != "" {
		batchDefID, err := primitive.ObjectIDFromHex(batchDefIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch_definition_id"})
			return
		}
		params.BatchDefinitionID = &batchDefID
	}

	// Get production report
	report, err := h.batchReportService.GetProductionReport(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetWastageReport handles GET /api/batch-reports/wastage
func (h *BatchReportHandler) GetWastageReport(c *gin.Context) {
	// Parse required query parameters
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	if fromDateStr == "" || toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from_date and to_date are required"})
		return
	}

	// Parse dates
	fromDate, err := time.Parse(time.RFC3339, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from_date format, use RFC3339"})
		return
	}

	toDate, err := time.Parse(time.RFC3339, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to_date format, use RFC3339"})
		return
	}

	// Parse optional parameters
	params := services.WastageReportParams{
		FromDate: fromDate,
		ToDate:   toDate,
	}

	// Parse batch definition ID if provided
	batchDefIDStr := c.Query("batch_definition_id")
	if batchDefIDStr != "" {
		batchDefID, err := primitive.ObjectIDFromHex(batchDefIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch_definition_id"})
			return
		}
		params.BatchDefinitionID = &batchDefID
	}

	// Get wastage report
	report, err := h.batchReportService.GetWastageReport(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetUsageReport handles GET /api/batch-reports/usage
func (h *BatchReportHandler) GetUsageReport(c *gin.Context) {
	// Parse required query parameters
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	if fromDateStr == "" || toDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from_date and to_date are required"})
		return
	}

	// Parse dates
	fromDate, err := time.Parse(time.RFC3339, fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from_date format, use RFC3339"})
		return
	}

	toDate, err := time.Parse(time.RFC3339, toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to_date format, use RFC3339"})
		return
	}

	// Parse optional parameters
	params := services.UsageReportParams{
		FromDate: fromDate,
		ToDate:   toDate,
	}

	// Parse batch definition ID if provided
	batchDefIDStr := c.Query("batch_definition_id")
	if batchDefIDStr != "" {
		batchDefID, err := primitive.ObjectIDFromHex(batchDefIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch_definition_id"})
			return
		}
		params.BatchDefinitionID = &batchDefID
	}

	// Parse menu item ID if provided
	menuItemIDStr := c.Query("menu_item_id")
	if menuItemIDStr != "" {
		menuItemID, err := primitive.ObjectIDFromHex(menuItemIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu_item_id"})
			return
		}
		params.MenuItemID = &menuItemID
	}

	// Get usage report
	report, err := h.batchReportService.GetUsageReport(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}
