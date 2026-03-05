package http

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/settings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HTMLTemplateHandler handles HTML template management
type HTMLTemplateHandler struct {
	chromedpRenderer *services.ChromedpBillRendererOptimized
	orderRepo        services.OrderRepository
	shopSettingsRepo settings.ShopSettingsRepository
	templatePath     string
}

// NewHTMLTemplateHandler creates a new HTML template handler
func NewHTMLTemplateHandler(
	chromedpRenderer *services.ChromedpBillRendererOptimized,
	orderRepo services.OrderRepository,
	shopSettingsRepo settings.ShopSettingsRepository,
	templatePath string,
) *HTMLTemplateHandler {
	return &HTMLTemplateHandler{
		chromedpRenderer: chromedpRenderer,
		orderRepo:        orderRepo,
		shopSettingsRepo: shopSettingsRepo,
		templatePath:     templatePath,
	}
}

// GetHTMLTemplate handles GET /api/html-templates/bill
func (h *HTMLTemplateHandler) GetHTMLTemplate(c *gin.Context) {
	// Read template file
	content, err := os.ReadFile(h.templatePath)
	if err != nil {
		log.Printf("Failed to read template: %v", err)
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

// UpdateHTMLTemplate handles PUT /api/html-templates/bill
type UpdateHTMLTemplateRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *HTMLTemplateHandler) UpdateHTMLTemplate(c *gin.Context) {
	var req UpdateHTMLTemplateRequest
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
		log.Printf("Failed to write template: %v", err)
		c.JSON(500, gin.H{"error": "Failed to save template"})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Template saved successfully",
		"backup":  backupPath,
	})
}

// TestPrintHTMLTemplate handles POST /api/html-templates/test-print
type TestPrintHTMLTemplateRequest struct {
	OrderID   string `json:"order_id" binding:"required"`
	PrinterIP string `json:"printer_ip" binding:"required"`
}

func (h *HTMLTemplateHandler) TestPrintHTMLTemplate(c *gin.Context) {
	var req TestPrintHTMLTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Parse order ID
	orderID, err := primitive.ObjectIDFromHex(req.OrderID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid order ID"})
		return
	}

	// Fetch order
	ord, err := h.orderRepo.FindByID(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Order not found"})
		return
	}

	// Fetch shop settings
	shopSettings, err := h.shopSettingsRepo.GetSettings(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch shop settings"})
		return
	}

	// Render bill using current template
	escposData, err := h.chromedpRenderer.RenderBillToESCPOS(ord, shopSettings)
	if err != nil {
		log.Printf("Failed to render bill: %v", err)
		c.JSON(500, gin.H{"error": "Failed to render bill: " + err.Error()})
		return
	}

	// Send to printer
	if err := sendToPrinterHTML(req.PrinterIP, escposData); err != nil {
		log.Printf("Failed to send to printer: %v", err)
		c.JSON(500, gin.H{"error": "Failed to send to printer: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"success":      true,
		"message":      "Test print successful",
		"order_number": ord.OrderNumber,
	})
}

// PreviewHTMLTemplate handles POST /api/html-templates/preview
type PreviewHTMLTemplateRequest struct {
	OrderID string `json:"order_id" binding:"required"`
}

func (h *HTMLTemplateHandler) PreviewHTMLTemplate(c *gin.Context) {
	var req PreviewHTMLTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Parse order ID
	orderID, err := primitive.ObjectIDFromHex(req.OrderID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid order ID"})
		return
	}

	// Fetch order
	ord, err := h.orderRepo.FindByID(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Order not found"})
		return
	}

	// Fetch shop settings
	shopSettings, err := h.shopSettingsRepo.GetSettings(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch shop settings"})
		return
	}

	// Save preview image
	filename := fmt.Sprintf("preview_html_template_%s.png", ord.OrderNumber)
	if err := h.chromedpRenderer.SavePreviewImage(ord, shopSettings, filename); err != nil {
		log.Printf("Failed to save preview: %v", err)
		c.JSON(500, gin.H{"error": "Failed to create preview: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"success":      true,
		"message":      "Preview created successfully",
		"filename":     filename,
		"order_number": ord.OrderNumber,
	})
}

// Helper functions

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func sendToPrinterHTML(printerIP string, data []byte) error {
	// Reuse the function from chromedp_print_handler
	return sendToPrinterChromedp(printerIP, data)
}
