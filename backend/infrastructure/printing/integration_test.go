package printing

import (
	"testing"

	"cafe-pos/backend/domain/printing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegration_FullWorkflow tests the complete FormatParser → TextRenderer → ImageConverter flow
func TestIntegration_FullWorkflow(t *testing.T) {
	// Sample Vietnamese receipt content
	content := `Cafe ABC
123 Nguyen Hue, Q1, TP.HCM
Tel: 028-1234-5678
================================
HOA DON BAN HANG
Order: ORD-001
Date: 2024-01-15 14:30
================================
Cafe Latte
  2 x 45,000 = 90,000
Banh Mi Thit
  1 x 35,000 = 35,000
================================
TONG CONG: 125,000 VND
Tien mat: 200,000 VND
Tien thua: 75,000 VND
================================
Cam on quy khach!
Hen gap lai!`

	// Create printer configuration
	config := &printing.PrinterConfig{
		Type:           printing.PrinterTypeBill,
		ConnectionType: printing.ConnectionTypeNetwork,
		IPAddress:      "192.168.1.100",
		Port:           9100,
		PaperWidth:     80,
	}

	// Create printer instance
	printer, err := NewESCPOSPrinter(config)
	require.NoError(t, err, "Failed to create printer")

	escposPrinter := printer.(*ESCPOSPrinter)

	// Verify all components are initialized
	assert.NotNil(t, escposPrinter.formatParser, "FormatParser should be initialized")
	assert.NotNil(t, escposPrinter.textRenderer, "TextRenderer should be initialized")
	assert.NotNil(t, escposPrinter.imageConverter, "ImageConverter should be initialized")

	// Test the full conversion workflow
	commands, err := escposPrinter.convertToESCPOS(content)
	require.NoError(t, err, "Failed to convert content to ESC/POS")
	require.NotNil(t, commands, "Commands should not be nil")

	// Verify command structure
	assert.True(t, len(commands) > 100, "Commands should contain substantial data (got %d bytes)", len(commands))

	// Verify ESC_INIT at start
	assert.Equal(t, ESC_INIT, commands[:len(ESC_INIT)], "Should start with ESC_INIT")

	// Verify GS v 0 command is present (image data)
	hasGSV0 := false
	for i := 0; i < len(commands)-len(GS_V_0); i++ {
		if commands[i] == GS_V_0[0] && commands[i+1] == GS_V_0[1] && commands[i+2] == GS_V_0[2] {
			hasGSV0 = true
			break
		}
	}
	assert.True(t, hasGSV0, "Should contain GS v 0 command for image data")

	// Verify GS_CUT at end
	endPos := len(commands) - len(GS_CUT)
	assert.Equal(t, GS_CUT, commands[endPos:], "Should end with GS_CUT")
}

// TestIntegration_ComponentFlow tests each component in the workflow
func TestIntegration_ComponentFlow(t *testing.T) {
	content := `Test Header
Item 1
================================
Total: 100,000 VND`

	// Step 1: FormatParser
	formatParser := NewFormatParser(80)
	lines := formatParser.Parse(content)
	require.NotEmpty(t, lines, "FormatParser should return lines")
	assert.Equal(t, 4, len(lines), "Should parse 4 lines")

	// Step 2: TextRenderer
	pixelWidth := CalculatePixelWidth(80)
	rendererConfig := &RendererConfig{
		PixelWidth:  pixelWidth,
		FontPath:    "",
		FontSize:    14.0,
		LineSpacing: 4,
		Margin:      8,
	}

	textRenderer, err := NewTextRenderer(rendererConfig)
	require.NoError(t, err, "Failed to create TextRenderer")

	img, err := textRenderer.Render(lines)
	require.NoError(t, err, "Failed to render text to image")
	require.NotNil(t, img, "Image should not be nil")

	// Verify image dimensions
	bounds := img.Bounds()
	assert.Equal(t, pixelWidth, bounds.Dx(), "Image width should match pixel width")
	assert.Greater(t, bounds.Dy(), 0, "Image height should be positive")

	// Step 3: ImageConverter
	imageConverter := NewImageConverter(pixelWidth)
	escposData, err := imageConverter.ConvertToESCPOS(img)
	require.NoError(t, err, "Failed to convert image to ESC/POS")
	require.NotEmpty(t, escposData, "ESC/POS data should not be empty")

	// Verify ESC/POS format
	assert.Equal(t, byte(0x1D), escposData[0], "Should start with GS")
	assert.Equal(t, byte(0x76), escposData[1], "Should have v")
	assert.Equal(t, byte(0x30), escposData[2], "Should have 0")
}

// TestIntegration_VietnameseCharacters tests Vietnamese character handling through the full workflow
func TestIntegration_VietnameseCharacters(t *testing.T) {
	// Vietnamese text with various diacritics
	content := `Cà phê sữa đá
Bánh mì thịt
Phở bò tái
Cơm tấm sườn
Bún chả Hà Nội`

	config := &printing.PrinterConfig{
		Type:           printing.PrinterTypeBill,
		ConnectionType: printing.ConnectionTypeNetwork,
		IPAddress:      "192.168.1.100",
		Port:           9100,
		PaperWidth:     58,
	}

	printer, err := NewESCPOSPrinter(config)
	require.NoError(t, err, "Failed to create printer")

	escposPrinter := printer.(*ESCPOSPrinter)

	// Convert to ESC/POS
	commands, err := escposPrinter.convertToESCPOS(content)
	require.NoError(t, err, "Failed to convert Vietnamese content")
	require.NotEmpty(t, commands, "Commands should not be empty")

	// Verify the workflow completed successfully
	assert.Greater(t, len(commands), 50, "Should generate substantial command data")
}

// TestIntegration_PaperWidthCalculation tests pixel width calculation for different paper sizes
func TestIntegration_PaperWidthCalculation(t *testing.T) {
	tests := []struct {
		name         string
		paperWidth   int
		expectedPixelWidth int
	}{
		{
			name:         "58mm paper",
			paperWidth:   58,
			expectedPixelWidth: 463,
		},
		{
			name:         "80mm paper",
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

			escposPrinter := printer.(*ESCPOSPrinter)

			// Verify pixel width is calculated correctly
			calculatedWidth := CalculatePixelWidth(tt.paperWidth)
			assert.Equal(t, tt.expectedPixelWidth, calculatedWidth, "Pixel width calculation mismatch")

			// Verify components use the correct pixel width
			assert.Equal(t, calculatedWidth, escposPrinter.textRenderer.pixelWidth, "TextRenderer pixel width mismatch")
			assert.Equal(t, calculatedWidth, escposPrinter.imageConverter.pixelWidth, "ImageConverter pixel width mismatch")
		})
	}
}
