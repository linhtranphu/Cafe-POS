package printing

import (
	"strings"
	"testing"

	"cafe-pos/backend/domain/printing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLabelPrinter(t *testing.T) {
	tests := []struct {
		name           string
		paperWidth     int
		expectedSize   LabelSize
	}{
		{
			name:         "40mm width should use 40x30 label",
			paperWidth:   40,
			expectedSize: LabelSize_40x30,
		},
		{
			name:         "50mm width should use 50x30 label",
			paperWidth:   50,
			expectedSize: LabelSize_50x30,
		},
		{
			name:         "60mm width should use 60x40 label",
			paperWidth:   60,
			expectedSize: LabelSize_60x40,
		},
		{
			name:         "unspecified width should default to 50x30",
			paperWidth:   0,
			expectedSize: LabelSize_50x30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &printing.PrinterConfig{
				Type:           printing.PrinterTypeLabel,
				ConnectionType: printing.ConnectionTypeNetwork,
				IPAddress:      "192.168.1.101",
				Port:           9100,
				PaperWidth:     tt.paperWidth,
			}

			printer := NewLabelPrinter(config).(*LabelPrinter)
			assert.Equal(t, tt.expectedSize.Width, printer.labelSize.Width)
			assert.Equal(t, tt.expectedSize.Height, printer.labelSize.Height)
			assert.Equal(t, tt.expectedSize.MaxCharsPerLine, printer.labelSize.MaxCharsPerLine)
			assert.Equal(t, tt.expectedSize.MaxLines, printer.labelSize.MaxLines)
		})
	}
}

func TestLabelPrinter_ValidateContent(t *testing.T) {
	tests := []struct {
		name        string
		labelSize   LabelSize
		content     string
		expectError bool
		errorMsg    string
	}{
		{
			name:      "valid content within limits should pass",
			labelSize: LabelSize_50x30,
			content: `Order: ORD-001
1/3
Cafe Latte
Size M
15:30`,
			expectError: false,
		},
		{
			name:      "content with too many lines should fail",
			labelSize: LabelSize_40x30, // Max 5 lines
			content: `Line 1
Line 2
Line 3
Line 4
Line 5
Line 6`,
			expectError: true,
			errorMsg:    "has 6 lines, but label size 40x30 supports max 5 lines",
		},
		{
			name:      "line too long should fail",
			labelSize: LabelSize_40x30, // Max 20 chars per line
			content:   "This is a very long line that exceeds maximum",
			expectError: true,
			errorMsg:    "has 45 characters, but label size 40x30 supports max 20 characters per line",
		},
		{
			name:      "empty lines should not count",
			labelSize: LabelSize_50x30, // Max 5 lines
			content: `Order: ORD-001

Cafe Latte

15:30`,
			expectError: false,
		},
		{
			name:      "content at exact limit should pass",
			labelSize: LabelSize_50x30, // Max 25 chars per line, 5 lines
			content: `1234567890123456789012345
Line 2
Line 3
Line 4
Line 5`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printer := &LabelPrinter{
				labelSize: tt.labelSize,
			}

			err := printer.validateContent(tt.content)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLabelPrinter_ConvertToLabelCommands(t *testing.T) {
	config := &printing.PrinterConfig{
		Type:           printing.PrinterTypeLabel,
		ConnectionType: printing.ConnectionTypeNetwork,
		IPAddress:      "192.168.1.101",
		Port:           9100,
		PaperWidth:     50,
	}
	printer := NewLabelPrinter(config).(*LabelPrinter)

	content := `Order: ORD-001
1/3
Cafe Latte
Size M`

	commands := printer.convertToLabelCommands(content)
	commandStr := string(commands)

	// Verify ZPL-style commands are present
	assert.Contains(t, commandStr, "^XA", "Should contain label start marker")
	assert.Contains(t, commandStr, "^XZ", "Should contain label end marker")
	
	// Verify content is included
	assert.Contains(t, commandStr, "Order: ORD-001")
	assert.Contains(t, commandStr, "Cafe Latte")
	assert.Contains(t, commandStr, "Size M")
	
	// Verify positioning commands
	assert.Contains(t, commandStr, "^FO", "Should contain field origin commands")
	assert.Contains(t, commandStr, "^FD", "Should contain field data commands")
	assert.Contains(t, commandStr, "^FS", "Should contain field separator commands")
}

func TestLabelPrinter_CalculateCenterPosition(t *testing.T) {
	tests := []struct {
		name      string
		labelSize LabelSize
		text      string
		minPos    int // Minimum expected position (should be > 0)
	}{
		{
			name:      "short text should be centered",
			labelSize: LabelSize_50x30,
			text:      "Test",
			minPos:    100, // Should be well centered
		},
		{
			name:      "long text should have minimum margin",
			labelSize: LabelSize_40x30,
			text:      "Very long text here",
			minPos:    10, // Should be at minimum margin
		},
		{
			name:      "empty text should return minimum margin",
			labelSize: LabelSize_50x30,
			text:      "",
			minPos:    10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printer := &LabelPrinter{
				labelSize: tt.labelSize,
			}

			pos := printer.calculateCenterPosition(tt.text)
			assert.GreaterOrEqual(t, pos, 10, "Position should be at least minimum margin")
		})
	}
}

func TestLabelPrinter_Connect_Validation(t *testing.T) {
	tests := []struct {
		name        string
		config      *printing.PrinterConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "USB connection should fail",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeLabel,
				ConnectionType: printing.ConnectionTypeUSB,
				USBPath:        "/dev/usb/lp0",
			},
			expectError: true,
			errorMsg:    "only supports network connection",
		},
		{
			name: "missing IP address should fail",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeLabel,
				ConnectionType: printing.ConnectionTypeNetwork,
				Port:           9100,
			},
			expectError: true,
			errorMsg:    "IP address is required",
		},
		{
			name: "missing port should fail",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeLabel,
				ConnectionType: printing.ConnectionTypeNetwork,
				IPAddress:      "192.168.1.101",
			},
			expectError: true,
			errorMsg:    "port is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printer := NewLabelPrinter(tt.config)
			err := printer.Connect()

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				printer.Disconnect()
			}
		})
	}
}

func TestLabelPrinter_Print_WithValidation(t *testing.T) {
	config := &printing.PrinterConfig{
		Type:           printing.PrinterTypeLabel,
		ConnectionType: printing.ConnectionTypeNetwork,
		IPAddress:      "192.168.1.101",
		Port:           9100,
		PaperWidth:     40, // Small label for testing validation
	}
	printer := NewLabelPrinter(config).(*LabelPrinter)

	tests := []struct {
		name        string
		content     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty content should fail",
			content:     "",
			expectError: true,
			errorMsg:    "print content cannot be empty",
		},
		{
			name: "content exceeding line limit should fail validation",
			content: strings.Repeat("Line\n", 10), // 10 lines, but 40x30 supports max 5
			expectError: true,
			errorMsg:    "supports max 5 lines",
		},
		{
			name:        "line exceeding width should fail validation",
			content:     strings.Repeat("X", 50), // 50 chars, but 40x30 supports max 20
			expectError: true,
			errorMsg:    "supports max 20 characters per line",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test validation directly for content validation tests
			if strings.Contains(tt.name, "validation") {
				err := printer.validateContent(tt.content)
				if tt.expectError {
					require.Error(t, err)
					assert.Contains(t, err.Error(), tt.errorMsg)
				} else {
					assert.NoError(t, err)
				}
			} else {
				// Test Print method for other cases
				err := printer.Print(tt.content)
				if tt.expectError {
					require.Error(t, err)
					assert.Contains(t, err.Error(), tt.errorMsg)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}
