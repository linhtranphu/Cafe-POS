package printing

import (
	"image"
	"image/color"
	"testing"
)

// TestNewImageCompositor tests the constructor
func TestNewImageCompositor(t *testing.T) {
	tests := []struct {
		name       string
		paperWidth int
		margin     int
		wantWidth  int
		wantMargin int
	}{
		{
			name:       "valid parameters",
			paperWidth: 576,
			margin:     10,
			wantWidth:  576,
			wantMargin: 10,
		},
		{
			name:       "zero paper width uses default",
			paperWidth: 0,
			margin:     10,
			wantWidth:  576,
			wantMargin: 10,
		},
		{
			name:       "negative paper width uses default",
			paperWidth: -100,
			margin:     10,
			wantWidth:  576,
			wantMargin: 10,
		},
		{
			name:       "negative margin becomes zero",
			paperWidth: 576,
			margin:     -5,
			wantWidth:  576,
			wantMargin: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compositor := NewImageCompositor(tt.paperWidth, tt.margin)
			if compositor.paperWidth != tt.wantWidth {
				t.Errorf("paperWidth = %d, want %d", compositor.paperWidth, tt.wantWidth)
			}
			if compositor.margin != tt.wantMargin {
				t.Errorf("margin = %d, want %d", compositor.margin, tt.wantMargin)
			}
		})
	}
}

// TestComposeWithoutLogo tests composing with nil logo
func TestComposeWithoutLogo(t *testing.T) {
	compositor := NewImageCompositor(576, 10)

	// Create a simple text content image
	textContent := image.NewGray(image.Rect(0, 0, 576, 100))
	// Fill with white
	for y := 0; y < 100; y++ {
		for x := 0; x < 576; x++ {
			textContent.SetGray(x, y, color.Gray{Y: 255})
		}
	}

	// Compose without logo
	result, err := compositor.Compose(nil, textContent)
	if err != nil {
		t.Fatalf("Compose() error = %v, want nil", err)
	}

	// Should return text content as-is
	if result != textContent {
		t.Error("Compose() with nil logo should return text content as-is")
	}
}

// TestComposeWithLogo tests composing with logo
func TestComposeWithLogo(t *testing.T) {
	compositor := NewImageCompositor(576, 10)

	// Create a simple logo (50x50 black square)
	logo := image.NewGray(image.Rect(0, 0, 50, 50))
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			logo.SetGray(x, y, color.Gray{Y: 0}) // Black
		}
	}

	// Create a simple text content image (576x100 white)
	textContent := image.NewGray(image.Rect(0, 0, 576, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 576; x++ {
			textContent.SetGray(x, y, color.Gray{Y: 255}) // White
		}
	}

	// Compose with logo
	result, err := compositor.Compose(logo, textContent)
	if err != nil {
		t.Fatalf("Compose() error = %v, want nil", err)
	}

	// Check result dimensions
	bounds := result.Bounds()
	expectedWidth := 576
	expectedHeight := 10 + 50 + 20 + 100 // margin + logo + spacing + text

	if bounds.Dx() != expectedWidth {
		t.Errorf("result width = %d, want %d", bounds.Dx(), expectedWidth)
	}
	if bounds.Dy() != expectedHeight {
		t.Errorf("result height = %d, want %d", bounds.Dy(), expectedHeight)
	}

	// Verify logo is at the correct position (margin, margin)
	// Logo should be black (Y=0) at position (10, 10)
	logoPixel := result.GrayAt(10, 10)
	if logoPixel.Y != 0 {
		t.Errorf("logo pixel at (10, 10) = %d, want 0 (black)", logoPixel.Y)
	}

	// Verify text content is at the correct position
	// Text should be white (Y=255) at position (0, 10+50+20) = (0, 80)
	textPixel := result.GrayAt(0, 80)
	if textPixel.Y != 255 {
		t.Errorf("text pixel at (0, 80) = %d, want 255 (white)", textPixel.Y)
	}
}

// TestComposeWithNilTextContent tests error handling for nil text content
func TestComposeWithNilTextContent(t *testing.T) {
	compositor := NewImageCompositor(576, 10)

	// Create a simple logo
	logo := image.NewGray(image.Rect(0, 0, 50, 50))

	// Compose with nil text content
	_, err := compositor.Compose(logo, nil)
	if err == nil {
		t.Error("Compose() with nil text content should return error")
	}
}

// TestCalculateTotalHeight tests height calculation
func TestCalculateTotalHeight(t *testing.T) {
	compositor := NewImageCompositor(576, 10)

	tests := []struct {
		name         string
		logoHeight   int
		textHeight   int
		wantHeight   int
	}{
		{
			name:       "with logo and text",
			logoHeight: 50,
			textHeight: 100,
			wantHeight: 10 + 50 + 20 + 100, // margin + logo + spacing + text
		},
		{
			name:       "without logo",
			logoHeight: 0,
			textHeight: 100,
			wantHeight: 10 + 100, // margin + text
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logo *image.Gray
			if tt.logoHeight > 0 {
				logo = image.NewGray(image.Rect(0, 0, 50, tt.logoHeight))
			}

			textContent := image.NewGray(image.Rect(0, 0, 576, tt.textHeight))

			height := compositor.calculateTotalHeight(logo, textContent)
			if height != tt.wantHeight {
				t.Errorf("calculateTotalHeight() = %d, want %d", height, tt.wantHeight)
			}
		})
	}
}

// TestDrawLogo tests logo drawing
func TestDrawLogo(t *testing.T) {
	compositor := NewImageCompositor(576, 10)

	// Create destination image (white background)
	dst := image.NewGray(image.Rect(0, 0, 576, 200))
	for y := 0; y < 200; y++ {
		for x := 0; x < 576; x++ {
			dst.SetGray(x, y, color.Gray{Y: 255})
		}
	}

	// Create logo (black square)
	logo := image.NewGray(image.Rect(0, 0, 30, 30))
	for y := 0; y < 30; y++ {
		for x := 0; x < 30; x++ {
			logo.SetGray(x, y, color.Gray{Y: 0})
		}
	}

	// Draw logo at position (10, 10)
	compositor.drawLogo(dst, logo, 10, 10)

	// Verify logo was drawn
	// Check a pixel inside the logo area
	if dst.GrayAt(15, 15).Y != 0 {
		t.Error("logo pixel should be black (0)")
	}

	// Check a pixel outside the logo area
	if dst.GrayAt(50, 50).Y != 255 {
		t.Error("background pixel should be white (255)")
	}
}

// TestDrawTextContent tests text content drawing
func TestDrawTextContent(t *testing.T) {
	compositor := NewImageCompositor(576, 10)

	// Create destination image (white background)
	dst := image.NewGray(image.Rect(0, 0, 576, 200))
	for y := 0; y < 200; y++ {
		for x := 0; x < 576; x++ {
			dst.SetGray(x, y, color.Gray{Y: 255})
		}
	}

	// Create text content (gray rectangle)
	textContent := image.NewGray(image.Rect(0, 0, 576, 50))
	for y := 0; y < 50; y++ {
		for x := 0; x < 576; x++ {
			textContent.SetGray(x, y, color.Gray{Y: 128}) // Gray
		}
	}

	// Draw text content at position (0, 100)
	compositor.drawTextContent(dst, textContent, 0, 100)

	// Verify text content was drawn
	// Check a pixel inside the text area
	if dst.GrayAt(100, 120).Y != 128 {
		t.Errorf("text pixel should be gray (128), got %d", dst.GrayAt(100, 120).Y)
	}

	// Check a pixel outside the text area
	if dst.GrayAt(100, 50).Y != 255 {
		t.Error("background pixel should be white (255)")
	}
}

// TestDrawLogoWithNilInputs tests error handling for nil inputs
func TestDrawLogoWithNilInputs(t *testing.T) {
	compositor := NewImageCompositor(576, 10)

	// Create valid images
	dst := image.NewGray(image.Rect(0, 0, 576, 200))
	logo := image.NewGray(image.Rect(0, 0, 30, 30))

	// Test with nil logo (should not panic)
	compositor.drawLogo(dst, nil, 10, 10)

	// Test with nil destination (should not panic)
	compositor.drawLogo(nil, logo, 10, 10)

	// Test with both nil (should not panic)
	compositor.drawLogo(nil, nil, 10, 10)
}

// TestDrawTextContentWithNilInputs tests error handling for nil inputs
func TestDrawTextContentWithNilInputs(t *testing.T) {
	compositor := NewImageCompositor(576, 10)

	// Create valid images
	dst := image.NewGray(image.Rect(0, 0, 576, 200))
	textContent := image.NewGray(image.Rect(0, 0, 576, 50))

	// Test with nil text content (should not panic)
	compositor.drawTextContent(dst, nil, 0, 100)

	// Test with nil destination (should not panic)
	compositor.drawTextContent(nil, textContent, 0, 100)

	// Test with both nil (should not panic)
	compositor.drawTextContent(nil, nil, 0, 100)
}
