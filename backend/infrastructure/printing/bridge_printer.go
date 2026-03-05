package printing

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"cafe-pos/backend/domain/printing"
	"cafe-pos/backend/infrastructure/printbridge"
)

// BridgePrinter wraps any printer to use print bridge instead of direct connection
type BridgePrinter struct {
	config       *printing.PrinterConfig
	bridgeClient *printbridge.Client
	innerPrinter Printer // The actual printer (ESCPOSPrinter or LabelPrinter)
}

// NewBridgePrinter creates a printer that uses print bridge
func NewBridgePrinter(config *printing.PrinterConfig, bridgeClient *printbridge.Client, innerPrinter Printer) Printer {
	return &BridgePrinter{
		config:       config,
		bridgeClient: bridgeClient,
		innerPrinter: innerPrinter,
	}
}

// Connect checks if print bridge is available (no actual connection needed)
func (p *BridgePrinter) Connect() error {
	if p.bridgeClient == nil {
		return fmt.Errorf("print bridge client not configured")
	}

	// Check if print bridge is available
	if !p.bridgeClient.IsAvailable() {
		return fmt.Errorf("print bridge is not available")
	}

	log.Printf("[BRIDGE PRINTER] Print bridge is available for printer: %s", p.config.Name)
	return nil
}

// Disconnect is a no-op for bridge printer
func (p *BridgePrinter) Disconnect() error {
	// No connection to close
	return nil
}

// Print sends content to print bridge for printing
func (p *BridgePrinter) Print(content string) error {
	if content == "" {
		return fmt.Errorf("print error: content cannot be empty")
	}

	if p.bridgeClient == nil {
		return fmt.Errorf("print error: print bridge client not configured")
	}

	ctx := context.Background()

	// Check content type and handle accordingly
	if isHTMLContent(content) {
		// HTML content - send to print bridge for rendering and printing
		log.Printf("[BRIDGE PRINTER] Sending HTML to print bridge - printer: %s, IP: %s:%d, size: %d bytes",
			p.config.Name, p.config.IPAddress, p.config.Port, len(content))

		if err := p.bridgeClient.RenderAndPrint(ctx, content, p.config.IPAddress, p.config.Port, 576); err != nil {
			return fmt.Errorf("print error: failed to render and print HTML via bridge: %w", err)
		}

		log.Printf("[BRIDGE PRINTER] HTML print successful via bridge - printer: %s", p.config.Name)
		return nil
	}

	// Binary or text content - convert to ESC/POS and send
	var escposData []byte
	var err error

	if isBase64Content(content) {
		// Decode base64 to get raw ESC/POS commands
		escposData, err = base64.StdEncoding.DecodeString(content)
		if err != nil {
			return fmt.Errorf("print error: failed to decode base64 content: %w", err)
		}
		log.Printf("[BRIDGE PRINTER] Using pre-rendered binary content (%d bytes)", len(escposData))
	} else {
		// Plain text content
		log.Printf("[BRIDGE PRINTER] Converting text to ESC/POS (%d chars)", len(content))
		escposData = []byte(content)
	}

	// Send ESC/POS data to print bridge
	log.Printf("[BRIDGE PRINTER] Sending ESC/POS to print bridge - printer: %s, IP: %s:%d, size: %d bytes",
		p.config.Name, p.config.IPAddress, p.config.Port, len(escposData))

	if err := p.bridgeClient.PrintESCPOS(ctx, escposData, p.config.IPAddress, p.config.Port); err != nil {
		return fmt.Errorf("print error: failed to print via bridge: %w", err)
	}

	log.Printf("[BRIDGE PRINTER] Print successful via bridge - printer: %s", p.config.Name)
	return nil
}

// GetStatus returns printer status (always online for bridge printer)
func (p *BridgePrinter) GetStatus() (PrinterStatus, error) {
	// Check if print bridge is available
	if p.bridgeClient == nil || !p.bridgeClient.IsAvailable() {
		return PrinterStatus{
			IsOnline:    false,
			PaperStatus: "UNKNOWN",
			ErrorMsg:    "Print bridge not available",
		}, nil
	}

	// Assume printer is online if bridge is available
	return PrinterStatus{
		IsOnline:    true,
		PaperStatus: "OK",
		ErrorMsg:    "",
	}, nil
}

// isHTMLContent checks if content is HTML
func isHTMLContent(content string) bool {
	// Simple check: HTML content typically starts with <!DOCTYPE or <html
	return len(content) > 10 && (
		content[:9] == "<!DOCTYPE" ||
		content[:5] == "<html" ||
		content[:5] == "<HTML")
}

