package printing

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"time"

	"cafe-pos/backend/domain/printing"
	
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// ESC/POS command constants
var (
	// ESC @ - Initialize printer
	ESC_INIT = []byte{0x1B, 0x40}

	// ESC a - Set alignment
	ESC_ALIGN_LEFT   = []byte{0x1B, 0x61, 0x00}
	ESC_ALIGN_CENTER = []byte{0x1B, 0x61, 0x01}
	ESC_ALIGN_RIGHT  = []byte{0x1B, 0x61, 0x02}

	// ESC E - Set bold
	ESC_BOLD_ON  = []byte{0x1B, 0x45, 0x01}
	ESC_BOLD_OFF = []byte{0x1B, 0x45, 0x00}

	// FS . - Cancel Kanji mode (enable Unicode)
	FS_CANCEL_KANJI = []byte{0x1C, 0x2E}

	// FS & - Select Kanji mode (for some printers, this enables Unicode)
	FS_SELECT_KANJI = []byte{0x1C, 0x26}

	// ESC t - Select character code table
	ESC_SELECT_CODE_TABLE = []byte{0x1B, 0x74}

	// LF - Line feed
	LF = []byte{0x0A}

	// GS V - Paper cut
	GS_CUT = []byte{0x1D, 0x56, 0x00}

	// ESC d - Print and feed n lines
	ESC_FEED_LINES = []byte{0x1B, 0x64}

	// GS v 0 - Print raster bit image
	GS_V_0 = []byte{0x1D, 0x76, 0x30}
)

// Image mode constants for GS v 0 command
const (
	IMAGE_MODE_NORMAL = 0x00
)

// Pixel width constants based on paper width and DPI
const (
	PIXEL_WIDTH_58MM = 384 // 58mm paper at 203 DPI
	PIXEL_WIDTH_80MM = 576 // 80mm paper at 203 DPI
	DPI              = 203 // Standard thermal printer resolution
)

// CalculatePixelWidth calculates the pixel width from paper width in millimeters.
// Formula: (paper_width_mm / 25.4) * DPI
// This converts millimeters to inches (divide by 25.4) and then to pixels (multiply by DPI).
// For 58mm paper: (58 / 25.4) * 203 ≈ 384 pixels
// For 80mm paper: (80 / 25.4) * 203 ≈ 576 pixels
func CalculatePixelWidth(paperWidthMM int) int {
	return int((float64(paperWidthMM) / 25.4) * float64(DPI))
}

// PrinterStatus represents the status of a printer
type PrinterStatus struct {
	IsOnline    bool
	PaperStatus string
	ErrorMsg    string
}

// Printer defines the interface for printer operations
type Printer interface {
	Connect() error
	Disconnect() error
	Print(content string) error
	GetStatus() (PrinterStatus, error)
}

// ESCPOSPrinter implements the Printer interface for ESC/POS thermal printers
type ESCPOSPrinter struct {
	config         *printing.PrinterConfig
	conn           net.Conn
	formatParser   *FormatParser
	textRenderer   *TextRenderer
	imageConverter *ImageConverter
}

// NewESCPOSPrinter creates a new ESC/POS printer instance
func NewESCPOSPrinter(config *printing.PrinterConfig) (Printer, error) {
	if config == nil {
		return nil, fmt.Errorf("printer initialization error: config cannot be nil")
	}

	if config.PaperWidth <= 0 {
		return nil, fmt.Errorf("printer initialization error: invalid paper width %d mm (must be positive)", config.PaperWidth)
	}

	// Calculate pixel width from paper width
	pixelWidth := CalculatePixelWidth(config.PaperWidth)

	// Initialize FormatParser
	formatParser := NewFormatParser(config.PaperWidth)

	// Initialize TextRenderer with font configuration
	rendererConfig := &RendererConfig{
		PixelWidth:  pixelWidth,
		FontPath:    "", // Empty path will trigger system font discovery
		FontSize:    26.0, // Increased from 14.0 for better readability
		LineSpacing: 10,    // Increased proportionally
		Margin:      14,
	}

	textRenderer, err := NewTextRenderer(rendererConfig)
	if err != nil {
		return nil, fmt.Errorf("printer initialization error: failed to initialize text renderer: %w", err)
	}

	// Initialize ImageConverter
	imageConverter := NewImageConverter(pixelWidth)

	return &ESCPOSPrinter{
		config:         config,
		formatParser:   formatParser,
		textRenderer:   textRenderer,
		imageConverter: imageConverter,
	}, nil
}

// Connect establishes a TCP/IP connection to the printer
func (p *ESCPOSPrinter) Connect() error {
	// Validate configuration
	if p.config.ConnectionType != printing.ConnectionTypeNetwork {
		return fmt.Errorf("printer connection error: ESC/POS printer only supports network connection (got: %s)", p.config.ConnectionType)
	}

	if p.config.IPAddress == "" {
		return fmt.Errorf("printer connection error: IP address is required for network printer")
	}

	if p.config.Port == 0 {
		return fmt.Errorf("printer connection error: port is required for network printer")
	}

	// Build connection address
	address := fmt.Sprintf("%s:%d", p.config.IPAddress, p.config.Port)

	// Establish TCP connection with timeout
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return fmt.Errorf("printer connection error: failed to connect to printer at %s: %w", address, err)
	}

	// Set read/write timeouts
	if err := conn.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		conn.Close()
		return fmt.Errorf("printer connection error: failed to set connection deadline for %s: %w", address, err)
	}

	p.conn = conn
	return nil
}

// Disconnect closes the connection to the printer
func (p *ESCPOSPrinter) Disconnect() error {
	if p.conn == nil {
		return nil // Already disconnected
	}

	err := p.conn.Close()
	p.conn = nil
	return err
}

// Print sends content to the printer with ESC/POS commands
func (p *ESCPOSPrinter) Print(content string) error {
	if content == "" {
		return fmt.Errorf("print error: content cannot be empty")
	}

	if p.conn == nil {
		return fmt.Errorf("print error: printer not connected (call Connect() first)")
	}

	var commands []byte
	var err error

	// Check if content is base64-encoded binary data
	// Base64 strings are typically much longer and contain only alphanumeric + / =
	if isBase64Content(content) {
		// Decode base64 to get raw ESC/POS commands
		commands, err = base64.StdEncoding.DecodeString(content)
		if err != nil {
			return fmt.Errorf("print error: failed to decode base64 content: %w", err)
		}
		log.Printf("[PRINTER] Using pre-rendered binary content (%d bytes)", len(commands))
	} else {
		// Plain text content - convert to ESC/POS
		if p.textRenderer == nil {
			return fmt.Errorf("print error: text renderer not initialized (font loading may have failed during initialization)")
		}

		commands, err = p.convertToESCPOS(content)
		if err != nil {
			return fmt.Errorf("print error: failed to convert content to ESC/POS format: %w", err)
		}
		log.Printf("[PRINTER] Converted text to ESC/POS (%d bytes)", len(commands))
	}

	// Send commands to printer
	bytesWritten, err := p.conn.Write(commands)
	if err != nil {
		return fmt.Errorf("print error: failed to send data to printer at %s:%d: %w", p.config.IPAddress, p.config.Port, err)
	}

	if bytesWritten != len(commands) {
		return fmt.Errorf("print error: incomplete data transmission (sent %d of %d bytes) to printer at %s:%d", bytesWritten, len(commands), p.config.IPAddress, p.config.Port)
	}

	return nil
}

// GetStatus retrieves the current status of the printer
func (p *ESCPOSPrinter) GetStatus() (PrinterStatus, error) {
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
	// Full status query would require DLE EOT commands which vary by printer model
	return PrinterStatus{
		IsOnline:    true,
		PaperStatus: "OK",
		ErrorMsg:    "",
	}, nil
}

// convertToESCPOS converts plain text template output to ESC/POS commands
// IMAGE RENDERING DISABLED - Using simple text mode with Windows-1258 encoding
// To re-enable image rendering, uncomment the code block below and comment out the simple text mode
func (p *ESCPOSPrinter) convertToESCPOS(content string) ([]byte, error) {
	// ============================================================================
	// SIMPLE TEXT MODE WITH WINDOWS-1258 ENCODING (CURRENT)
	// ============================================================================
	var commands []byte
	
	// Initialize printer
	commands = append(commands, ESC_INIT...)
	
	// Set code page to Vietnamese
	// ESC t n - Select character code table
	// For Zywell ZY303 and similar Vietnamese thermal printers:
	//   n = 255: Internal Vietnamese code page (Zywell/Xprinter specific) - RECOMMENDED
	//   n = 31: Windows-1258 (Vietnamese) - Standard
	//   n = 30: Alternative internal Vietnamese code page
	// Using n = 255 (0xFF) for Zywell internal Vietnamese code page
	commands = append(commands, 0x1B, 0x74, 0xFF) // ESC t 255
	
	// Optional: Select Font A (12x24) for better clarity
	// ESC M n - Select character font
	// n = 0: Font A (12x24), n = 1: Font B (9x17)
	commands = append(commands, 0x1B, 0x4D, 0x00) // ESC M 0 (Font A)
	
	// Encode content from UTF-8 to Windows-1258
	encodedContent := encodeVietnamese(content)
	commands = append(commands, encodedContent...)
	
	// Feed a few lines before cutting
	commands = append(commands, ESC_FEED_LINES...)
	commands = append(commands, 0x03) // Feed 3 lines
	
	// Cut paper
	commands = append(commands, GS_CUT...)
	
	return commands, nil
	
	// ============================================================================
	// IMAGE RENDERING MODE (DISABLED)
	// Uncomment this section to re-enable Vietnamese font rendering via images
	// ============================================================================
	/*
	// Check if text renderer is available
	if p.textRenderer == nil {
		return nil, fmt.Errorf("conversion error: text renderer not initialized (font loading may have failed)")
	}

	// Parse content to identify formatting
	lines := p.formatParser.Parse(content)
	if len(lines) == 0 {
		return nil, fmt.Errorf("conversion error: format parser returned no lines from content")
	}

	// Render formatted lines to bitmap image
	img, err := p.textRenderer.Render(lines)
	if err != nil {
		return nil, fmt.Errorf("conversion error: text rendering failed: %w", err)
	}

	// Convert image to ESC/POS format
	imageData, err := p.imageConverter.ConvertToESCPOS(img)
	if err != nil {
		return nil, fmt.Errorf("conversion error: image to ESC/POS conversion failed: %w", err)
	}

	// Build complete command sequence
	var commands []byte

	// Initialize printer
	commands = append(commands, ESC_INIT...)

	// Append image data
	commands = append(commands, imageData...)

	// Feed a few lines before cutting
	commands = append(commands, ESC_FEED_LINES...)
	commands = append(commands, 0x03) // Feed 3 lines

	// Cut paper
	commands = append(commands, GS_CUT...)

	return commands, nil
	*/
}

// encodeVietnamese converts UTF-8 string to Windows-1258 encoding for Vietnamese text
func encodeVietnamese(input string) []byte {
	// Chuyển đổi từ UTF-8 sang Windows-1258
	encoder := charmap.Windows1258.NewEncoder()
	output, _, _ := transform.String(encoder, input)
	return []byte(output)
}
