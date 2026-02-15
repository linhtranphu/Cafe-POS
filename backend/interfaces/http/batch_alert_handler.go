package http

import (
	"net/http"

	"cafe-pos/backend/application/services"
	"github.com/gin-gonic/gin"
)

// BatchAlertHandler handles batch alert HTTP endpoints
type BatchAlertHandler struct {
	batchAlertService *services.BatchAlertService
}

// NewBatchAlertHandler creates a new batch alert handler
func NewBatchAlertHandler(batchAlertService *services.BatchAlertService) *BatchAlertHandler {
	return &BatchAlertHandler{
		batchAlertService: batchAlertService,
	}
}

// GetAlerts handles GET /api/batch-alerts
func (h *BatchAlertHandler) GetAlerts(c *gin.Context) {
	alerts, err := h.batchAlertService.GetAlerts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alerts)
}
