package printing

import (
	"fmt"
	"net"
	"strings"
	"time"

	"cafe-pos/backend/domain/printing"
)

// Label size configurations (width x height in mm)
type LabelSize struct {
	Width  int // mm
	Height int // mm
	MaxCharsPerLine int
	MaxLines int
}

var (
	// Supported label sizes
	LabelSize_40x30 = LabelSize{Width: 40, Height: 30, MaxCharsPerLine: 20, MaxLines: 5}
	LabelSize_50x30 = LabelSize{Width: 50, Height: 30, MaxCharsPerLine: 25, MaxLines: 5}
	LabelSize_60x40 = LabelSize{Width: 60, Height: 40, MaxCharsPerLine: 30, MaxLines: 7}
)

// LabelPrinter implements the Printer interface for label printers
type LabelPrinter struct {
	config    *printing.PrinterConfig
	conn      net.Conn
	labelSize LabelSize
}

// NewLabelPrinter creates a new label printer instance
func NewLabelPrinter(config *printing.PrinterConfig) Printer {
	// Determine label size based on paper width
	// Default to 50x30mm if not specified
	labelSize := LabelSize_50x30
	
	switch config.PaperWidth {
	case 40:
		labelSize = LabelSize_40x30
	case 50:
		labelSize = LabelSize_50x30
	case 60:
		labelSize = LabelSize_60x40
	}

	return &LabelPrinter{
		config:    config,
		labelSize: labelSize,
	}
}

// Connect establishes a TCP/IP connection to the label printer
func (p *LabelPrinter) Connect() error {
	// Validate configuration
	if p.config.ConnectionType != printing.ConnectionTypeNetwork {
		return fmt.Errorf("label printer only supports network connection")
	}

	if p.config.IPAddress == "" {
		return fmt.Errorf("IP address is required for network printer")
	}

	if p.config.Port == 0 {
		return fmt.Errorf("port is required for network printer")
	}

	// Build connection address
	address := fmt.Sprintf("%s:%d", p.config.IPAddress, p.config.Port)

	// Establish TCP connection with timeout
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to label printer at %s: %w", address, err)
	}

	// Set read/write timeouts
	if err := conn.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		conn.Close()
		return fmt.Errorf("failed to set connection deadline: %w", err)
	}

	p.conn = conn
	return nil
}

// Disconnect closes the connection to the label printer
func (p *LabelPrinter) Disconnect() error {
	if p.conn == nil {
		return nil // Already disconnected
	}

	err := p.conn.Close()
	p.conn = nil
	return err
}

// Print sends content to the label printer
func (p *LabelPrinter) Print(content string) error {
	if content == "" {
		return fmt.Errorf("print content cannot be empty")
	}

	if p.conn == nil {
		return fmt.Errorf("printer not connected")
	}

	// Validate content fits within label size constraints
	if err := p.validateContent(content); err != nil {
		return fmt.Errorf("content validation failed: %w", err)
	}

	// Convert plain text to label printer commands
	commands := p.convertToLabelCommands(content)

	// Send commands to printer
	_, err := p.conn.Write(commands)
	if err != nil {
		return fmt.Errorf("failed to send data to label printer: %w", err)
	}

	return nil
}

// GetStatus retrieves the current status of the label printer
func (p *LabelPrinter) GetStatus() (PrinterStatus, error) {
	// Check if connection is alive
	if p.conn == nil {
		return PrinterStatus{
			IsOnline:    false,
			PaperStatus: "UNKNOWN",
			ErrorMsg:    "Not connected",
		}, nil
	}

	// Try to set a deadline to test connection
	if err := p.conn.SetDeadline(time.Now().Add(2 * time.Second)); err != nil {
		return PrinterStatus{
			IsOnline:    false,
			PaperStatus: "UNKNOWN",
			ErrorMsg:    "Connection error",
		}, nil
	}

	// For basic implementation, assume printer is online if connected
	// Full status query would require printer-specific commands
	return PrinterStatus{
		IsOnline:    true,
		PaperStatus: "OK",
		ErrorMsg:    "",
	}, nil
}

// validateContent checks if content fits within label size constraints
func (p *LabelPrinter) validateContent(content string) error {
	lines := strings.Split(content, "\n")
	
	// Check number of lines
	nonEmptyLines := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines++
		}
	}
	
	if nonEmptyLines > p.labelSize.MaxLines {
		return fmt.Errorf("content has %d lines, but label size %dx%d supports max %d lines", 
			nonEmptyLines, p.labelSize.Width, p.labelSize.Height, p.labelSize.MaxLines)
	}

	// Check line width
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > p.labelSize.MaxCharsPerLine {
			return fmt.Errorf("line %d has %d characters, but label size %dx%d supports max %d characters per line",
				i+1, len(trimmed), p.labelSize.Width, p.labelSize.Height, p.labelSize.MaxCharsPerLine)
		}
	}

	return nil
}

// convertToLabelCommands converts plain text to label printer commands
// This is a basic implementation that sends raw text
// Can be enhanced with specific command sets (ZPL, EPL, TSPL) later
func (p *LabelPrinter) convertToLabelCommands(content string) []byte {
	var commands []byte

	// For basic implementation, we'll send formatted text
	// In production, this would be replaced with printer-specific commands
	// like ZPL (Zebra), EPL (Eltron), or TSPL (TSC)

	// Add label start marker (generic)
	commands = append(commands, []byte("^XA")...) // ZPL-style start (can be adapted)
	commands = append(commands, []byte("\n")...)

	// Process content line by line
	lines := strings.Split(content, "\n")
	yPosition := 20 // Starting Y position

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Center text on label
		xPosition := p.calculateCenterPosition(trimmed)

		// Add text command (simplified format)
		// In production, use proper ZPL/EPL/TSPL commands
		textCmd := fmt.Sprintf("^FO%d,%d^A0N,25,25^FD%s^FS\n", xPosition, yPosition, trimmed)
		commands = append(commands, []byte(textCmd)...)

		yPosition += 30 // Move to next line
	}

	// Add label end marker
	commands = append(commands, []byte("^XZ")...) // ZPL-style end
	commands = append(commands, []byte("\n")...)

	return commands
}

// calculateCenterPosition calculates X position to center text on label
func (p *LabelPrinter) calculateCenterPosition(text string) int {
	// Approximate character width in dots (assuming 8 dots per mm, 25pt font)
	charWidth := 15 // dots per character
	textWidth := len(text) * charWidth
	
	// Label width in dots (assuming 8 dots per mm)
	labelWidth := p.labelSize.Width * 8
	
	// Calculate center position
	xPosition := (labelWidth - textWidth) / 2
	if xPosition < 10 {
		xPosition = 10 // Minimum margin
	}
	
	return xPosition
}
