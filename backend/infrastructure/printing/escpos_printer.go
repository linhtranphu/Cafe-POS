package printing

import (
	"fmt"
	"net"
	"strings"
	"time"

	"cafe-pos/backend/domain/printing"
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

	// LF - Line feed
	LF = []byte{0x0A}

	// GS V - Paper cut
	GS_CUT = []byte{0x1D, 0x56, 0x00}

	// ESC d - Print and feed n lines
	ESC_FEED_LINES = []byte{0x1B, 0x64}
)

// PrinterStatus represents the status of a printer
type PrinterStatus struct {
	IsOnline    bool   `json:"is_online"`
	PaperStatus string `json:"paper_status"` // OK, LOW, OUT
	ErrorMsg    string `json:"error_msg,omitempty"`
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
	config     *printing.PrinterConfig
	conn       net.Conn
	paperWidth int // characters per line
}

// NewESCPOSPrinter creates a new ESC/POS printer instance
func NewESCPOSPrinter(config *printing.PrinterConfig) Printer {
	// Calculate characters per line based on paper width
	// Typical thermal printers: 58mm ≈ 32 chars, 80mm ≈ 48 chars
	charsPerLine := 48 // default for 80mm
	if config.PaperWidth == 58 {
		charsPerLine = 32
	}

	return &ESCPOSPrinter{
		config:     config,
		paperWidth: charsPerLine,
	}
}

// Connect establishes a TCP/IP connection to the printer
func (p *ESCPOSPrinter) Connect() error {
	// Validate configuration
	if p.config.ConnectionType != printing.ConnectionTypeNetwork {
		return fmt.Errorf("ESC/POS printer only supports network connection")
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
		return fmt.Errorf("failed to connect to printer at %s: %w", address, err)
	}

	// Set read/write timeouts
	if err := conn.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		conn.Close()
		return fmt.Errorf("failed to set connection deadline: %w", err)
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
		return fmt.Errorf("print content cannot be empty")
	}

	if p.conn == nil {
		return fmt.Errorf("printer not connected")
	}

	// Convert plain text to ESC/POS commands
	commands := p.convertToESCPOS(content)

	// Send commands to printer
	_, err := p.conn.Write(commands)
	if err != nil {
		return fmt.Errorf("failed to send data to printer: %w", err)
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
func (p *ESCPOSPrinter) convertToESCPOS(content string) []byte {
	var commands []byte

	// Initialize printer
	commands = append(commands, ESC_INIT...)

	// Process content line by line
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		// Detect formatting based on line content
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			// Empty line - just line feed
			commands = append(commands, LF...)
			continue
		}

		// Check for separator lines (===, ---)
		if isSeparatorLine(trimmed) {
			commands = append(commands, ESC_ALIGN_CENTER...)
			commands = append(commands, []byte(trimmed)...)
			commands = append(commands, LF...)
			commands = append(commands, ESC_ALIGN_LEFT...)
			continue
		}

		// Check if line should be centered (typically headers, totals)
		if shouldCenter(trimmed) {
			commands = append(commands, ESC_ALIGN_CENTER...)
			commands = append(commands, ESC_BOLD_ON...)
			commands = append(commands, []byte(trimmed)...)
			commands = append(commands, ESC_BOLD_OFF...)
			commands = append(commands, LF...)
			commands = append(commands, ESC_ALIGN_LEFT...)
			continue
		}

		// Check if line should be bold (typically labels like "TOTAL:")
		if shouldBold(trimmed) {
			commands = append(commands, ESC_BOLD_ON...)
			commands = append(commands, []byte(trimmed)...)
			commands = append(commands, ESC_BOLD_OFF...)
			commands = append(commands, LF...)
			continue
		}

		// Regular line - ensure it fits within paper width
		wrapped := p.wrapLine(line)
		for _, wrappedLine := range wrapped {
			commands = append(commands, []byte(wrappedLine)...)
			commands = append(commands, LF...)
		}
	}

	// Feed a few lines before cutting
	commands = append(commands, ESC_FEED_LINES...)
	commands = append(commands, 0x03) // Feed 3 lines

	// Cut paper
	commands = append(commands, GS_CUT...)

	return commands
}

// isSeparatorLine checks if a line is a separator (===, ---)
func isSeparatorLine(line string) bool {
	if len(line) == 0 {
		return false
	}
	// Check if line consists only of = or - characters
	for _, ch := range line {
		if ch != '=' && ch != '-' {
			return false
		}
	}
	return true
}

// shouldCenter checks if a line should be centered
func shouldCenter(line string) bool {
	// Center lines that look like headers or totals
	upper := strings.ToUpper(line)
	
	// Shop name/info (typically at top)
	if !strings.Contains(line, ":") && len(line) < 30 {
		return true
	}
	
	// Total line
	if strings.HasPrefix(upper, "TOTAL") {
		return true
	}
	
	// Thank you message
	if strings.Contains(upper, "THANK") || strings.Contains(upper, "CẢM ƠN") {
		return true
	}
	
	return false
}

// shouldBold checks if a line should be bold
func shouldBold(line string) bool {
	upper := strings.ToUpper(line)
	
	// Lines with TOTAL, SUBTOTAL, DISCOUNT
	if strings.Contains(upper, "TOTAL") || 
	   strings.Contains(upper, "DISCOUNT") || 
	   strings.Contains(upper, "GIẢM GIÁ") {
		return true
	}
	
	return false
}

// wrapLine wraps a line to fit within paper width
func (p *ESCPOSPrinter) wrapLine(line string) []string {
	// If line fits, return as-is
	if len(line) <= p.paperWidth {
		return []string{line}
	}

	// Split long lines
	var wrapped []string
	remaining := line
	
	for len(remaining) > p.paperWidth {
		// Try to break at a space
		breakPoint := p.paperWidth
		for i := p.paperWidth - 1; i > p.paperWidth/2; i-- {
			if remaining[i] == ' ' {
				breakPoint = i
				break
			}
		}
		
		wrapped = append(wrapped, remaining[:breakPoint])
		remaining = strings.TrimLeft(remaining[breakPoint:], " ")
	}
	
	if len(remaining) > 0 {
		wrapped = append(wrapped, remaining)
	}
	
	return wrapped
}
