package http

import (
	"fmt"
	"log"
	"net"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/settings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChromedpPrintHandler handles chromedp-based bill printing requests
type ChromedpPrintHandler struct {
	chromedpRenderer *services.ChromedpBillRendererOptimized
	orderRepo        services.OrderRepository
	shopSettingsRepo settings.ShopSettingsRepository
}

// NewChromedpPrintHandler creates a new chromedp print handler
func NewChromedpPrintHandler(
	orderRepo services.OrderRepository,
	shopSettingsRepo settings.ShopSettingsRepository,
) (*ChromedpPrintHandler, error) {
	// Initialize Chromedp renderer
	chromedpRenderer, err := services.NewChromedpBillRendererOptimized()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Chromedp renderer: %w", err)
	}

	return &ChromedpPrintHandler{
		chromedpRenderer: chromedpRenderer,
		orderRepo:        orderRepo,
		shopSettingsRepo: shopSettingsRepo,
	}, nil
}

// Close cleans up resources
func (h *ChromedpPrintHandler) Close() {
	if h.chromedpRenderer != nil {
		h.chromedpRenderer.Close()
	}
}

// GetRenderer returns the chromedp renderer instance
func (h *ChromedpPrintHandler) GetRenderer() *services.ChromedpBillRendererOptimized {
	return h.chromedpRenderer
}

// PrintChromedpBillRequest represents the request to print a bill using chromedp
type PrintChromedpBillRequest struct {
	OrderID   string `json:"order_id" binding:"required"`
	PrinterIP string `json:"printer_ip" binding:"required"`
}

// PrintChromedpBill handles POST /api/chromedp-print/bill
func (h *ChromedpPrintHandler) PrintChromedpBill(c *gin.Context) {
	var req PrintChromedpBillRequest
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

	// Render bill to ESC/POS using chromedp
	escposData, err := h.chromedpRenderer.RenderBillToESCPOS(ord, shopSettings)
	if err != nil {
		log.Printf("Failed to render chromedp bill: %v", err)
		c.JSON(500, gin.H{"error": "Failed to render bill: " + err.Error()})
		return
	}

	// Send to printer
	if err := sendToPrinterChromedp(req.PrinterIP, escposData); err != nil {
		log.Printf("Failed to send to printer: %v", err)
		c.JSON(500, gin.H{"error": "Failed to send to printer: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"success":      true,
		"message":      "Bill printed successfully using Chromedp",
		"order_number": ord.OrderNumber,
	})
}

// PreviewChromedpBill handles GET /api/chromedp-print/preview/:order_id
func (h *ChromedpPrintHandler) PreviewChromedpBill(c *gin.Context) {
	orderIDStr := c.Param("order_id")

	// Parse order ID
	orderID, err := primitive.ObjectIDFromHex(orderIDStr)
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
	filename := fmt.Sprintf("preview_chromedp_%s.png", ord.OrderNumber)
	if err := h.chromedpRenderer.SavePreviewImage(ord, shopSettings, filename); err != nil {
		log.Printf("Failed to save preview: %v", err)
		c.JSON(500, gin.H{"error": "Failed to create preview: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"success":      true,
		"message":      "Preview created successfully using Chromedp",
		"filename":     filename,
		"order_number": ord.OrderNumber,
	})
}

// sendToPrinterChromedp sends ESC/POS data to a network printer
func sendToPrinterChromedp(printerIP string, data []byte) error {
	// Add port if not specified
	if _, _, err := net.SplitHostPort(printerIP); err != nil {
		printerIP = printerIP + ":9100"
	}

	// Connect to printer
	conn, err := net.DialTimeout("tcp", printerIP, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to printer: %w", err)
	}
	defer conn.Close()

	// Set write deadline
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// Send data
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to printer: %w", err)
	}

	return nil
}
