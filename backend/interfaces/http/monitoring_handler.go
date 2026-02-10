package http

import (
	"net/http"
	"strconv"

	"cafe-pos/backend/application/services"
	"github.com/gin-gonic/gin"
)

// MonitoringHandler handles monitoring and metrics endpoints
type MonitoringHandler struct {
	monitoringService *services.MonitoringService
}

// NewMonitoringHandler creates a new monitoring handler
func NewMonitoringHandler(monitoringService *services.MonitoringService) *MonitoringHandler {
	return &MonitoringHandler{
		monitoringService: monitoringService,
	}
}

// GetMetrics returns metrics filtered by type
// GET /api/monitoring/metrics?type=cost_calculation&limit=100
func (h *MonitoringHandler) GetMetrics(c *gin.Context) {
	metricType := services.MetricType(c.Query("type"))
	limitStr := c.DefaultQuery("limit", "100")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000 // Cap at 1000
	}
	
	metrics, err := h.monitoringService.GetMetrics(c.Request.Context(), metricType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch metrics",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
		"count":   len(metrics),
	})
}

// GetAlerts returns recent alerts filtered by level
// GET /api/monitoring/alerts?level=critical&limit=50
func (h *MonitoringHandler) GetAlerts(c *gin.Context) {
	level := services.AlertLevel(c.Query("level"))
	limitStr := c.DefaultQuery("limit", "50")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500 // Cap at 500
	}
	
	alerts, err := h.monitoringService.GetAlerts(c.Request.Context(), level, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch alerts",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"alerts": alerts,
		"count":  len(alerts),
	})
}

// GetAggregatedMetrics returns aggregated metrics summary
// GET /api/monitoring/metrics/aggregated
func (h *MonitoringHandler) GetAggregatedMetrics(c *gin.Context) {
	metrics, err := h.monitoringService.GetAggregatedMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch aggregated metrics",
		})
		return
	}
	
	c.JSON(http.StatusOK, metrics)
}

// GetHealthStatus returns the overall health status
// GET /api/monitoring/health
func (h *MonitoringHandler) GetHealthStatus(c *gin.Context) {
	health, err := h.monitoringService.GetHealthStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch health status",
		})
		return
	}
	
	// Set HTTP status based on health
	status := http.StatusOK
	if healthStatus, ok := health["status"].(string); ok {
		if healthStatus == "critical" {
			status = http.StatusServiceUnavailable
		} else if healthStatus == "degraded" {
			status = http.StatusOK // Still return 200 but with degraded status
		}
	}
	
	c.JSON(status, health)
}
