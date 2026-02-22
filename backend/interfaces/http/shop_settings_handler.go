package http

import (
	"net/http"

	"cafe-pos/backend/domain/settings"
	"cafe-pos/backend/infrastructure/mongodb"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ShopSettingsHandler handles HTTP requests for shop settings
type ShopSettingsHandler struct {
	repo *mongodb.ShopSettingsRepository
}

// NewShopSettingsHandler creates a new shop settings handler
func NewShopSettingsHandler(repo *mongodb.ShopSettingsRepository) *ShopSettingsHandler {
	return &ShopSettingsHandler{
		repo: repo,
	}
}

// GetSettings retrieves the shop settings
// GET /api/shop-settings
func (h *ShopSettingsHandler) GetSettings(c *gin.Context) {
	shopSettings, err := h.repo.GetSettings(c.Request.Context())
	if err != nil {
		if err == settings.ErrSettingsNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shop settings not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shopSettings)
}

// UpdateSettingsRequest represents the request body for updating shop settings
type UpdateSettingsRequest struct {
	ShopName          string `json:"shop_name" binding:"required"`
	ShopAddress       string `json:"shop_address"`
	ShopPhone         string `json:"shop_phone"`
	LogoURL           string `json:"logo_url"`
	CustomMessage     string `json:"custom_message"`
	ShowLogo          bool   `json:"show_logo"`
	ShowAddress       bool   `json:"show_address"`
	ShowPhone         bool   `json:"show_phone"`
	ShowCustomMessage bool   `json:"show_custom_message"`
	AutoPrintEnabled  bool   `json:"auto_print_enabled"`
}

// UpdateSettings updates the shop settings
// PUT /api/shop-settings/:id
func (h *ShopSettingsHandler) UpdateSettings(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing settings
	shopSettings, err := h.repo.GetSettingsByID(c.Request.Context(), id)
	if err != nil {
		if err == settings.ErrSettingsNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shop settings not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	shopSettings.ShopName = req.ShopName
	shopSettings.UpdatePrintSettings(
		req.ShopAddress,
		req.ShopPhone,
		req.LogoURL,
		req.CustomMessage,
	)

	shopSettings.SetFieldVisibility(
		req.ShowLogo,
		req.ShowAddress,
		req.ShowPhone,
		req.ShowCustomMessage,
	)
	shopSettings.SetAutoPrintEnabled(req.AutoPrintEnabled)

	// Save to database
	if err := h.repo.UpdateSettings(c.Request.Context(), shopSettings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shopSettings)
}
