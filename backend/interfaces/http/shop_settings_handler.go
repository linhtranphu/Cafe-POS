package http

import (
	"context"
	"net/http"
	"time"

	"cafe-pos/backend/domain/settings"
	"cafe-pos/backend/infrastructure/mongodb"
	"cafe-pos/backend/infrastructure/printbridge"

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
	PrintBridgeURL    string `json:"print_bridge_url"`
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
	shopSettings.PrintBridgeURL = req.PrintBridgeURL
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

// CreateSettings creates new shop settings
// POST /api/shop-settings
func (h *ShopSettingsHandler) CreateSettings(c *gin.Context) {
	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if settings already exist
	existing, err := h.repo.FindFirst(c.Request.Context())
	if err == nil && existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Shop settings already exist. Use PUT to update."})
		return
	}

	// Create new settings
	shopSettings := settings.NewShopSettings(req.ShopName)
	shopSettings.PrintBridgeURL = req.PrintBridgeURL
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
	if err := h.repo.CreateSettings(c.Request.Context(), shopSettings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shopSettings)
}

// TestPrintBridgeRequest represents the request body for testing print bridge
type TestPrintBridgeRequest struct {
	BridgeURL string `json:"bridge_url" binding:"required"`
}

// TestPrintBridge tests connection to print bridge
// POST /api/manager/print-bridge/test
func (h *ShopSettingsHandler) TestPrintBridge(c *gin.Context) {
	var req TestPrintBridgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Import print bridge client
	printBridgeClient := printbridge.NewClient(req.BridgeURL, 5*time.Second)
	
	// Test connection
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	
	if err := printBridgeClient.TestConnection(ctx); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Print bridge is not available",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Print bridge is available",
	})
}
