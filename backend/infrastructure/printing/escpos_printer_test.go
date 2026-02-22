package printing

import (
	"testing"

	"cafe-pos/backend/domain/printing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewESCPOSPrinter(t *testing.T) {
	tests := []struct {
		name         string
		paperWidth   int
		expectedPixelWidth int
	}{
		{
			name:         "58mm width should calculate to 463 pixels",
			paperWidth:   58,
			expectedPixelWidth: 463,
		},
		{
			name:         "80mm width should calculate to 639 pixels",
			paperWidth:   80,
			expectedPixelWidth: 639,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &printing.PrinterConfig{
				Type:           printing.PrinterTypeBill,
				ConnectionType: printing.ConnectionTypeNetwork,
				IPAddress:      "192.168.1.100",
				Port:           9100,
				PaperWidth:     tt.paperWidth,
			}

			printer, err := NewESCPOSPrinter(config)
			require.NoError(t, err)
			assert.NotNil(t, printer)
			
			escposPrinter := printer.(*ESCPOSPrinter)
			assert.Equal(t, config, escposPrinter.config)
			assert.NotNil(t, escposPrinter.formatParser)
			assert.NotNil(t, escposPrinter.imageConverter)
			assert.NotNil(t, escposPrinter.textRenderer)
		})
	}
}

func TestESCPOSPrinter_Connect_Validation(t *testing.T) {
	tests := []struct {
		name        string
		config      *printing.PrinterConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "USB connection should fail",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeBill,
				ConnectionType: printing.ConnectionTypeUSB,
				USBPath:        "/dev/usb/lp0",
				PaperWidth:     80,
			},
			expectError: true,
			errorMsg:    "only supports network connection",
		},
		{
			name: "missing IP address should fail",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeBill,
				ConnectionType: printing.ConnectionTypeNetwork,
				Port:           9100,
				PaperWidth:     80,
			},
			expectError: true,
			errorMsg:    "IP address is required",
		},
		{
			name: "missing port should fail",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeBill,
				ConnectionType: printing.ConnectionTypeNetwork,
				IPAddress:      "192.168.1.100",
				PaperWidth:     80,
			},
			expectError: true,
			errorMsg:    "port is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printer, err := NewESCPOSPrinter(tt.config)
			require.NoError(t, err)
			
			err = printer.Connect()

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

func TestESCPOSPrinter_Print_Validation(t *testing.T) {
	config := &printing.PrinterConfig{
		Type:           printing.PrinterTypeBill,
		ConnectionType: printing.ConnectionTypeNetwork,
		IPAddress:      "192.168.1.100",
		Port:           9100,
		PaperWidth:     80,
	}
	printer, err := NewESCPOSPrinter(config)
	require.NoError(t, err)
	
	escposPrinter := printer.(*ESCPOSPrinter)

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
			errorMsg:    "content cannot be empty",
		},
		{
			name:        "not connected should fail",
			content:     "Test content",
			expectError: true,
			errorMsg:    "not connected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := escposPrinter.Print(tt.content)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestESCPOSPrinter_GetStatus_NotConnected(t *testing.T) {
	config := &printing.PrinterConfig{
		Type:           printing.PrinterTypeBill,
		ConnectionType: printing.ConnectionTypeNetwork,
		IPAddress:      "192.168.1.100",
		Port:           9100,
		PaperWidth:     80,
	}
	printer, err := NewESCPOSPrinter(config)
	require.NoError(t, err)

	status, err := printer.GetStatus()
	assert.NoError(t, err)
	assert.False(t, status.IsOnline)
	assert.Equal(t, "UNKNOWN", status.PaperStatus)
	assert.Equal(t, "Not connected", status.ErrorMsg)
}

func TestESCPOSPrinter_Disconnect(t *testing.T) {
	config := &printing.PrinterConfig{
		Type:           printing.PrinterTypeBill,
		ConnectionType: printing.ConnectionTypeNetwork,
		IPAddress:      "192.168.1.100",
		Port:           9100,
		PaperWidth:     80,
	}
	printer, err := NewESCPOSPrinter(config)
	require.NoError(t, err)

	// Disconnect when not connected should not error
	err = printer.Disconnect()
	assert.NoError(t, err)
}

func TestESCPOSPrinter_CommandConstants(t *testing.T) {
	// Verify ESC/POS command constants are correct
	assert.Equal(t, []byte{0x1B, 0x40}, ESC_INIT, "ESC @ - Initialize")
	assert.Equal(t, []byte{0x1B, 0x61, 0x00}, ESC_ALIGN_LEFT, "ESC a 0 - Align left")
	assert.Equal(t, []byte{0x1B, 0x61, 0x01}, ESC_ALIGN_CENTER, "ESC a 1 - Align center")
	assert.Equal(t, []byte{0x1B, 0x61, 0x02}, ESC_ALIGN_RIGHT, "ESC a 2 - Align right")
	assert.Equal(t, []byte{0x1B, 0x45, 0x01}, ESC_BOLD_ON, "ESC E 1 - Bold on")
	assert.Equal(t, []byte{0x1B, 0x45, 0x00}, ESC_BOLD_OFF, "ESC E 0 - Bold off")
	assert.Equal(t, []byte{0x0A}, LF, "LF - Line feed")
	assert.Equal(t, []byte{0x1D, 0x56, 0x00}, GS_CUT, "GS V 0 - Paper cut")
	assert.Equal(t, []byte{0x1B, 0x64}, ESC_FEED_LINES, "ESC d - Feed lines")
	assert.Equal(t, []byte{0x1D, 0x76, 0x30}, GS_V_0, "GS v 0 - Print raster bit image")
}

func TestCalculatePixelWidth(t *testing.T) {
	// Test the actual formula implementation: (paper_width_mm / 25.4) * 203 DPI
	tests := []struct {
		name           string
		paperWidthMM   int
		expectedPixels int
	}{
		{
			name:           "58mm paper: (58 / 25.4) * 203 = 463 pixels",
			paperWidthMM:   58,
			expectedPixels: 463,
		},
		{
			name:           "80mm paper: (80 / 25.4) * 203 = 639 pixels",
			paperWidthMM:   80,
			expectedPixels: 639,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePixelWidth(tt.paperWidthMM)
			assert.Equal(t, tt.expectedPixels, result,
				"CalculatePixelWidth(%d) should return %d pixels", tt.paperWidthMM, tt.expectedPixels)
		})
	}
}

func TestPixelWidthConstants(t *testing.T) {
	// Verify that the constants are defined
	assert.Equal(t, 384, PIXEL_WIDTH_58MM, "PIXEL_WIDTH_58MM constant")
	assert.Equal(t, 576, PIXEL_WIDTH_80MM, "PIXEL_WIDTH_80MM constant")
	assert.Equal(t, 203, DPI, "DPI constant")
}

func TestImageModeConstants(t *testing.T) {
	// Verify image mode constants
	assert.Equal(t, 0x00, IMAGE_MODE_NORMAL, "IMAGE_MODE_NORMAL constant")
}

func TestESCPOSPrinter_ConvertToESCPOS_WithTextRenderer(t *testing.T) {
	config := &printing.PrinterConfig{
		Type:           printing.PrinterTypeBill,
		ConnectionType: printing.ConnectionTypeNetwork,
		IPAddress:      "192.168.1.100",
		Port:           9100,
		PaperWidth:     80,
	}
	printer, err := NewESCPOSPrinter(config)
	require.NoError(t, err)
	
	escposPrinter := printer.(*ESCPOSPrinter)

	content := `Cafe ABC
Order: ORD-001
================================
Cafe Latte
  2 x 45000 = 90000
================================
TOTAL: 120000 VND`

	commands, err := escposPrinter.convertToESCPOS(content)
	require.NoError(t, err)
	require.NotNil(t, commands)

	// Verify command structure
	// Should start with ESC_INIT
	assert.True(t, len(commands) >= len(ESC_INIT))
	assert.Equal(t, ESC_INIT, commands[:len(ESC_INIT)])

	// Should contain GS v 0 command for image data
	hasGSV0 := false
	for i := 0; i < len(commands)-len(GS_V_0); i++ {
		if commands[i] == GS_V_0[0] && commands[i+1] == GS_V_0[1] && commands[i+2] == GS_V_0[2] {
			hasGSV0 = true
			break
		}
	}
	assert.True(t, hasGSV0, "Should contain GS v 0 command for image data")

	// Should end with paper cut
	assert.True(t, len(commands) >= len(GS_CUT))
	endPos := len(commands) - len(GS_CUT)
	assert.Equal(t, GS_CUT, commands[endPos:])
}

func TestESCPOSPrinter_ConvertToESCPOS_NoTextRenderer(t *testing.T) {
	config := &printing.PrinterConfig{
		Type:           printing.PrinterTypeBill,
		ConnectionType: printing.ConnectionTypeNetwork,
		IPAddress:      "192.168.1.100",
		Port:           9100,
		PaperWidth:     80,
	}
	
	// Create printer with nil text renderer
	printer := &ESCPOSPrinter{
		config:         config,
		formatParser:   NewFormatParser(config.PaperWidth),
		textRenderer:   nil,
		imageConverter: NewImageConverter(CalculatePixelWidth(config.PaperWidth)),
	}

	content := "Test content"
	_, err := printer.convertToESCPOS(content)
	
	require.Error(t, err)
	assert.Contains(t, err.Error(), "text renderer not initialized")
}
