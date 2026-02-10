package http

import (
	"net/http"

	"cafe-pos/backend/application/services"
	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
	settingsService *services.ShopSettingsService
}

func NewSettingsHandler(settingsService *services.ShopSettingsService) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
	}
}

// GetSettings retrieves the shop settings
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	settings, err := h.settingsService.GetSettings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateSettings updates the shop settings
func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	var req services.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate low_margin_threshold if provided
	if req.LowMarginThreshold != nil && *req.LowMarginThreshold < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "low_margin_threshold must be >= 0"})
		return
	}

	settings, err := h.settingsService.UpdateSettings(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}
