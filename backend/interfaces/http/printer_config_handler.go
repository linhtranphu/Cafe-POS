package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/printing"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PrinterConfigHandler handles HTTP requests for printer configuration management
type PrinterConfigHandler struct {
	printerConfigRepo printing.PrinterConfigRepository
	printerManager    services.PrinterManager
}

// NewPrinterConfigHandler creates a new PrinterConfigHandler
func NewPrinterConfigHandler(printerConfigRepo printing.PrinterConfigRepository, printerManager services.PrinterManager) *PrinterConfigHandler {
	return &PrinterConfigHandler{
		printerConfigRepo: printerConfigRepo,
		printerManager:    printerManager,
	}
}

// CreatePrinterRequest represents the request body for creating a printer
type CreatePrinterRequest struct {
	Name           string                   `json:"name" binding:"required"`
	Type           printing.PrinterType     `json:"type" binding:"required"`
	ConnectionType printing.ConnectionType  `json:"connection_type" binding:"required"`
	IPAddress      string                   `json:"ip_address,omitempty"`
	Port           int                      `json:"port,omitempty"`
	USBPath        string                   `json:"usb_path,omitempty"`
	PaperWidth     int                      `json:"paper_width" binding:"required"`
	IsDefault      bool                     `json:"is_default"`
	IsEnabled      bool                     `json:"is_enabled"`
}

// UpdatePrinterRequest represents the request body for updating a printer
type UpdatePrinterRequest struct {
	Name           string                   `json:"name"`
	Type           printing.PrinterType     `json:"type"`
	ConnectionType printing.ConnectionType  `json:"connection_type"`
	IPAddress      string                   `json:"ip_address,omitempty"`
	Port           int                      `json:"port,omitempty"`
	USBPath        string                   `json:"usb_path,omitempty"`
	PaperWidth     int                      `json:"paper_width"`
	IsDefault      bool                     `json:"is_default"`
	IsEnabled      bool                     `json:"is_enabled"`
}

// ListPrinters handles GET /api/printers
func (h *PrinterConfigHandler) ListPrinters(c *gin.Context) {
	ctx := c.Request.Context()

	// Get query parameter for filtering by type
	printerType := c.Query("type")

	var printers []*printing.PrinterConfig
	var err error

	if printerType != "" {
		printers, err = h.printerConfigRepo.FindByType(ctx, printing.PrinterType(printerType))
	} else {
		printers, err = h.printerConfigRepo.FindAll(ctx)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch printers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"printers": printers})
}

// GetPrinter handles GET /api/printers/:id
func (h *PrinterConfigHandler) GetPrinter(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid printer ID format"})
		return
	}

	printer, err := h.printerConfigRepo.FindByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Printer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch printer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"printer": printer})
}

// CreatePrinter handles POST /api/printers
func (h *PrinterConfigHandler) CreatePrinter(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreatePrinterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate connection details based on connection type
	if req.ConnectionType == printing.ConnectionTypeNetwork {
		if req.IPAddress == "" || req.Port == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "IP address and port are required for network connection"})
			return
		}
	} else if req.ConnectionType == printing.ConnectionTypeUSB {
		if req.USBPath == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "USB path is required for USB connection"})
			return
		}
	}

	// If setting as default, unset other defaults of the same type
	if req.IsDefault {
		if err := h.unsetDefaultPrinters(ctx, req.Type); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update default printers"})
			return
		}
	}

	// Create printer config
	config := &printing.PrinterConfig{
		Name:           req.Name,
		Type:           req.Type,
		ConnectionType: req.ConnectionType,
		IPAddress:      req.IPAddress,
		Port:           req.Port,
		USBPath:        req.USBPath,
		PaperWidth:     req.PaperWidth,
		IsDefault:      req.IsDefault,
		IsEnabled:      req.IsEnabled,
	}

	if err := h.printerConfigRepo.Create(ctx, config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create printer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"printer": config})
}

// UpdatePrinter handles PUT /api/printers/:id
func (h *PrinterConfigHandler) UpdatePrinter(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid printer ID format"})
		return
	}

	// Fetch existing printer
	existing, err := h.printerConfigRepo.FindByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Printer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch printer"})
		return
	}

	var req UpdatePrinterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Type != "" {
		existing.Type = req.Type
	}
	if req.ConnectionType != "" {
		existing.ConnectionType = req.ConnectionType
	}
	if req.IPAddress != "" {
		existing.IPAddress = req.IPAddress
	}
	if req.Port != 0 {
		existing.Port = req.Port
	}
	if req.USBPath != "" {
		existing.USBPath = req.USBPath
	}
	if req.PaperWidth != 0 {
		existing.PaperWidth = req.PaperWidth
	}
	existing.IsDefault = req.IsDefault
	existing.IsEnabled = req.IsEnabled

	// If setting as default, unset other defaults of the same type
	if req.IsDefault {
		if err := h.unsetDefaultPrinters(ctx, existing.Type); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update default printers"})
			return
		}
	}

	if err := h.printerConfigRepo.Update(ctx, existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update printer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"printer": existing})
}

// DeletePrinter handles DELETE /api/printers/:id
func (h *PrinterConfigHandler) DeletePrinter(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid printer ID format"})
		return
	}

	if err := h.printerConfigRepo.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete printer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Printer deleted successfully"})
}

// TestConnection handles POST /api/printers/:id/test
func (h *PrinterConfigHandler) TestConnection(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid printer ID format"})
		return
	}

	// Fetch printer config
	config, err := h.printerConfigRepo.FindByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Printer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch printer"})
		return
	}

	// Get printer instance
	printer, err := h.printerManager.GetPrinter(config)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Failed to get printer: " + err.Error(),
		})
		return
	}

	// Connect to printer
	if err := printer.Connect(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Connection failed: " + err.Error(),
		})
		return
	}
	defer printer.Disconnect()

	// Print test receipt
	testContent := `================================
       TEST PRINT
================================
Printer: ` + config.Name + `
Type: ` + string(config.Type) + `
Connection: ` + string(config.ConnectionType) + `
IP: ` + config.IPAddress + `:` + fmt.Sprintf("%d", config.Port) + `
================================
This is a test print
Đây là bản in thử nghiệm
================================
Date: ` + time.Now().Format("2006-01-02 15:04:05") + `
================================


`

	if err := printer.Print(testContent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Print test failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Test print successful - Check your printer",
	})
}

// unsetDefaultPrinters unsets the is_default flag for all printers of a specific type
func (h *PrinterConfigHandler) unsetDefaultPrinters(ctx context.Context, printerType printing.PrinterType) error {
	printers, err := h.printerConfigRepo.FindByType(ctx, printerType)
	if err != nil {
		return err
	}

	for _, printer := range printers {
		if printer.IsDefault {
			printer.IsDefault = false
			if err := h.printerConfigRepo.Update(ctx, printer); err != nil {
				return err
			}
		}
	}

	return nil
}
