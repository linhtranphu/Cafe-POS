package http

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"

	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/domain/settings"
	"cafe-pos/backend/infrastructure/printbridge"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HTMLTemplateHandlerBridge handles HTML template management using print bridge
type HTMLTemplateHandlerBridge struct {
	orderRepo        OrderRepository
	shopSettingsRepo settings.ShopSettingsRepository
	templatePath     string
}

// OrderRepository interface for fetching orders
type OrderRepository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*order.Order, error)
}

// NewHTMLTemplateHandlerBridge creates a new handler using print bridge
func NewHTMLTemplateHandlerBridge(
	orderRepo OrderRepository,
	shopSettingsRepo settings.ShopSettingsRepository,
	templatePath string,
) *HTMLTemplateHandlerBridge {
	return &HTMLTemplateHandlerBridge{
		orderRepo:        orderRepo,
		shopSettingsRepo: shopSettingsRepo,
		templatePath:     templatePath,
	}
}

// getPrintBridgeClient creates a print bridge client from settings
func (h *HTMLTemplateHandlerBridge) getPrintBridgeClient(ctx context.Context) (*printbridge.Client, error) {
	shopSettings, err := h.shopSettingsRepo.GetSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop settings: %w", err)
	}

	if shopSettings.PrintBridgeURL == "" {
		return nil, fmt.Errorf("print bridge URL not configured in settings")
	}

	return printbridge.NewClient(shopSettings.PrintBridgeURL, 30*time.Second), nil
}

// GetHTMLTemplate handles GET /api/html-templates/bill
func (h *HTMLTemplateHandlerBridge) GetHTMLTemplate(c *gin.Context) {
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
func (h *HTMLTemplateHandlerBridge) UpdateHTMLTemplate(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Backup current template
	backupPath := h.templatePath + ".backup"
	if err := copyFileBridge(h.templatePath, backupPath); err != nil {
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
func (h *HTMLTemplateHandlerBridge) TestPrintHTMLTemplate(c *gin.Context) {
	var req struct {
		OrderID   string `json:"order_id"`
		PrinterIP string `json:"printer_ip" binding:"required"`
		UseTestData bool `json:"use_test_data"` // If true, use mock data instead of real order
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get print bridge client from settings
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

	// Fetch shop settings
	shopSettings, err := h.shopSettingsRepo.GetSettings(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch shop settings"})
		return
	}

	var ord *order.Order
	
	// Use test data or real order
	if req.UseTestData || req.OrderID == "" {
		// Create mock order for testing
		ord = createMockOrder(shopSettings)
	} else {
		// Parse order ID
		orderID, err := primitive.ObjectIDFromHex(req.OrderID)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid order ID"})
			return
		}

		// Fetch order
		ord, err = h.orderRepo.FindByID(c.Request.Context(), orderID)
		if err != nil {
			c.JSON(404, gin.H{"error": "Order not found"})
			return
		}
	}

	// Render HTML with template
	htmlContent, err := h.renderBillHTML(ord, shopSettings)
	if err != nil {
		log.Printf("Failed to render HTML: %v", err)
		c.JSON(500, gin.H{"error": "Failed to render HTML: " + err.Error()})
		return
	}

	// Send HTML to print bridge for rendering and printing
	if err := printBridgeClient.RenderAndPrint(c.Request.Context(), htmlContent, req.PrinterIP, 9100, 576); err != nil {
		log.Printf("Failed to render and print via print bridge: %v", err)
		c.JSON(500, gin.H{"error": "Failed to print: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"success":      true,
		"message":      "Test print successful (via print bridge)",
		"order_number": ord.OrderNumber,
		"is_test_data": req.UseTestData || req.OrderID == "",
	})
}

// PreviewHTMLTemplate handles POST /api/html-templates/preview
func (h *HTMLTemplateHandlerBridge) PreviewHTMLTemplate(c *gin.Context) {
	var req struct {
		OrderID string `json:"order_id" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get print bridge client from settings
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

	// For preview, we just return success
	// The actual preview rendering happens in print bridge
	c.JSON(200, gin.H{
		"success":      true,
		"message":      "Preview will be rendered by print bridge",
		"order_number": ord.OrderNumber,
		"note":         "Preview rendering moved to local print bridge",
	})
}

// renderBillHTML renders the bill HTML from template
func (h *HTMLTemplateHandlerBridge) renderBillHTML(ord *order.Order, shopSettings *settings.ShopSettings) (string, error) {
	// Read template
	templateContent, err := os.ReadFile(h.templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}

	// Parse template
	tmpl, err := template.New("bill").Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Prepare data (similar to chromedp renderer)
	data := prepareBillData(ord, shopSettings)

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// prepareBillData prepares data for template rendering
func prepareBillData(ord *order.Order, settings *settings.ShopSettings) map[string]interface{} {
	// Format items
	items := make([]map[string]interface{}, 0, len(ord.Items))
	for i, item := range ord.Items {
		itemName := item.Name
		if item.VariantName != "" {
			itemName = fmt.Sprintf("%s (%s)", item.Name, item.VariantName)
		}

		items = append(items, map[string]interface{}{
			"STT":       i + 1,
			"Name":      itemName,
			"Quantity":  item.Quantity,
			"UnitPrice": formatMoneyVNBridge(item.Price),
			"Total":     formatMoneyVNBridge(item.Subtotal),
			"Note":      item.Note,
		})
	}

	data := map[string]interface{}{
		"ShopName":      settings.ShopName,
		"ShopAddress":   settings.ShopAddress,
		"ShopPhone":     settings.ShopPhone,
		"ShowLogo":      settings.ShowLogo,
		"ShowAddress":   settings.ShowAddress,
		"ShowPhone":     settings.ShowPhone,
		"OrderNumber":   ord.OrderNumber,
		"CustomerName":  ord.CustomerName,
		"WaiterName":    ord.WaiterName,
		"PaymentMethod": string(ord.PaymentMethod),
		"CreatedDate":   ord.CreatedAt.Format("02/01/2006 03:04 PM"),
		"Items":         items,
		"Total":         formatMoneyVNBridge(ord.Total),
		"CustomMessage": settings.CustomMessage,
		"ShowCustomMsg": settings.ShowCustomMessage,
		"LogoURL":       settings.LogoURL,
	}
	
	return data
}

// Helper functions

func createMockOrder(settings *settings.ShopSettings) *order.Order {
	now := time.Now()
	return &order.Order{
		ID:           primitive.NewObjectID(),
		OrderNumber:  fmt.Sprintf("TEST-%d", now.Unix()),
		CustomerName: "Nguyễn Văn A",
		WaiterName:   "Test Waiter",
		Items: []order.OrderItem{
			{
				Name:         "Cà phê sữa đá",
				VariantName:  "Size M",
				Quantity:     2,
				Price:        25000,
				Subtotal:     50000,
			},
			{
				Name:         "Trà sữa trân châu",
				VariantName:  "Size L",
				Quantity:     1,
				Price:        35000,
				Subtotal:     35000,
			},
			{
				Name:         "Bánh mì",
				VariantName:  "",
				Quantity:     1,
				Price:        20000,
				Subtotal:     20000,
			},
		},
		Total:         105000,
		PaymentMethod: order.PaymentCash,
		Status:        order.StatusServed,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func formatMoneyVNBridge(amount float64) string {
	amountInt := int(amount)
	if amountInt >= 1000000 {
		return fmt.Sprintf("%d,%03d,%03d", amountInt/1000000, (amountInt%1000000)/1000, amountInt%1000)
	} else if amountInt >= 1000 {
		return fmt.Sprintf("%d,%03d", amountInt/1000, amountInt%1000)
	}
	return fmt.Sprintf("%d", amountInt)
}

func copyFileBridge(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func sendToPrinterBridge(printerIP string, data []byte) error {
	// Reuse existing printer sending logic
	return sendToPrinterChromedp(printerIP, data)
}

// GetTempBillTemplate handles GET /api/html-templates/temp-bill
func (h *HTMLTemplateHandlerBridge) GetTempBillTemplate(c *gin.Context) {
	// Read temp bill template file
	tempBillPath := "./application/services/templates/temp_bill_template.html"
	content, err := os.ReadFile(tempBillPath)
	if err != nil {
		log.Printf("Failed to read temp bill template: %v", err)
		c.JSON(500, gin.H{"error": "Failed to read temp bill template"})
		return
	}

	c.JSON(200, gin.H{
		"success":  true,
		"content":  string(content),
		"path":     tempBillPath,
		"filename": filepath.Base(tempBillPath),
	})
}

// UpdateTempBillTemplate handles PUT /api/html-templates/temp-bill
func (h *HTMLTemplateHandlerBridge) UpdateTempBillTemplate(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	tempBillPath := "./application/services/templates/temp_bill_template.html"

	// Backup current template
	backupPath := tempBillPath + ".backup"
	if err := copyFileBridge(tempBillPath, backupPath); err != nil {
		log.Printf("Warning: Failed to create backup: %v", err)
	}

	// Write new template
	if err := os.WriteFile(tempBillPath, []byte(req.Content), 0644); err != nil {
		log.Printf("Failed to write temp bill template: %v", err)
		c.JSON(500, gin.H{"error": "Failed to save temp bill template"})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Temp bill template saved successfully",
		"backup":  backupPath,
	})
}
