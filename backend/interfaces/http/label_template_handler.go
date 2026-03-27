package http

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"cafe-pos/backend/domain/settings"
	"cafe-pos/backend/infrastructure/printing"
	"cafe-pos/backend/infrastructure/printbridge"

	"github.com/gin-gonic/gin"
)

// LabelTemplateHandler handles label template management
type LabelTemplateHandler struct {
	orderRepo        OrderRepository
	shopSettingsRepo settings.ShopSettingsRepository
	templatePath     string
}

// NewLabelTemplateHandler creates a new handler
func NewLabelTemplateHandler(
	orderRepo OrderRepository,
	shopSettingsRepo settings.ShopSettingsRepository,
	templatePath string,
) *LabelTemplateHandler {
	return &LabelTemplateHandler{
		orderRepo:        orderRepo,
		shopSettingsRepo: shopSettingsRepo,
		templatePath:     templatePath,
	}
}

// getPrintBridgeClient creates a print bridge client from settings
func (h *LabelTemplateHandler) getPrintBridgeClient(ctx context.Context) (*printbridge.Client, error) {
	shopSettings, err := h.shopSettingsRepo.GetSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop settings: %w", err)
	}

	if shopSettings.PrintBridgeURL == "" {
		return nil, fmt.Errorf("print bridge URL not configured in settings")
	}

	return printbridge.NewClient(shopSettings.PrintBridgeURL, 30*time.Second), nil
}


// GetLabelTemplate handles GET /api/label-templates/order-item
func (h *LabelTemplateHandler) GetLabelTemplate(c *gin.Context) {
	content, err := os.ReadFile(h.templatePath)
	if err != nil {
		log.Printf("Failed to read label template: %v", err)
		c.JSON(500, gin.H{"error": "Failed to read template"})
		return
	}

	c.JSON(200, gin.H{
		"success":  true,
		"content":  string(content),
		"path":     h.templatePath,
		"filename": filepath.Base(h.templatePath),
	})
}

// UpdateLabelTemplate handles PUT /api/label-templates/order-item
func (h *LabelTemplateHandler) UpdateLabelTemplate(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Backup current template
	backupPath := h.templatePath + ".backup"
	if err := copyFile(h.templatePath, backupPath); err != nil {
		log.Printf("Warning: Failed to create backup: %v", err)
	}

	// Write new template
	if err := os.WriteFile(h.templatePath, []byte(req.Content), 0644); err != nil {
		log.Printf("Failed to write label template: %v", err)
		c.JSON(500, gin.H{"error": "Failed to save template"})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Label template saved successfully",
		"backup":  backupPath,
	})
}

// TestPrintLabel handles POST /api/label-templates/test-print
func (h *LabelTemplateHandler) TestPrintLabel(c *gin.Context) {
	var req struct {
		ItemName     string `json:"item_name" binding:"required"`
		Note         string `json:"note"`
		CustomerName string `json:"customer_name" binding:"required"`
		PrinterIP    string `json:"printer_ip" binding:"required"`
		Port         int    `json:"port"`
		PaperWidth   int    `json:"paper_width"` // from printer config (mm)
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Default port
	if req.Port == 0 {
		req.Port = 9100
	}

	// Get print bridge client
	printBridgeClient, err := h.getPrintBridgeClient(c.Request.Context())
	if err != nil {
		c.JSON(503, gin.H{
			"error": "Print bridge not configured: " + err.Error(),
			"hint":  "Please configure Print Bridge URL in Settings",
		})
		return
	}

	// Check if print bridge is available
	if !printBridgeClient.IsAvailable() {
		c.JSON(503, gin.H{
			"error": "Print bridge is not available. Please ensure local print bridge is running.",
		})
		return
	}

	// Determine label dimensions from paper_width in request
	// paper_width is the width of the label roll (e.g. 40mm)
	// height defaults to 30mm for standard labels
	labelWidth := req.PaperWidth
	labelHeight := 30 // default height
	if labelWidth == 0 {
		labelWidth = 40 // default 40mm
	}

	// Initialize TSPL generator
	tsplGenerator := printing.NewTSPLGenerator(
		labelWidth,
		labelHeight,
		203,
	)

	// Load template
	templateContent, err := os.ReadFile(h.templatePath)
	if err != nil {
		log.Printf("Failed to read label template: %v", err)
		c.JSON(500, gin.H{"error": "Failed to read template"})
		return
	}

	// Prepare test data
	data := printing.LabelData{
		OrderNumber:  "TEST-001",
		ItemName:     req.ItemName,
		Note:         req.Note,
		Time:         time.Now().Format("15:04"),
		CustomerName: req.CustomerName,
	}

	// Generate TSPL
	tsplCommands, err := tsplGenerator.GenerateLabelCommands(data, string(templateContent))
	if err != nil {
		log.Printf("Failed to generate TSPL: %v", err)
		c.JSON(500, gin.H{"error": "Failed to generate TSPL: " + err.Error()})
		return
	}

	// Send to printer
	err = printBridgeClient.SendTSPLCommands(c.Request.Context(), tsplCommands, req.PrinterIP, req.Port)
	if err != nil {
		log.Printf("Failed to print label: %v", err)
		c.JSON(500, gin.H{"error": "Failed to print: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Test print successful",
	})
}

