package printing

import (
	"testing"
)

// TestNewTextRenderer_WithFontSizes tests creating a renderer with multiple font sizes
func TestNewTextRenderer_WithFontSizes(t *testing.T) {
	config := &RendererConfig{
		PixelWidth:  384,
		FontPath:    "", // Use system font
		FontSize:    18.0,
		FontSizes: map[string]float64{
			"normal": 18.0,
			"header": 22.0,
			"total":  20.0,
		},
		LineSpacing: 4,
		Margin:      8,
	}

	renderer, err := NewTextRenderer(config)
	if err != nil {
		t.Fatalf("failed to create renderer with font sizes: %v", err)
	}

	if renderer == nil {
		t.Fatal("expected renderer but got nil")
	}

	// Verify fonts map contains all three sizes
	if len(renderer.fonts) != 3 {
		t.Errorf("expected 3 font sizes, got %d", len(renderer.fonts))
	}

	// Verify each font size is loaded
	expectedSizes := []float64{18.0, 20.0, 22.0}
	for _, size := range expectedSizes {
		fontPair, exists := renderer.fonts[size]
		if !exists {
			t.Errorf("expected font size %.1fpt to be loaded", size)
			continue
		}

		if fontPair.Normal == nil {
			t.Errorf("expected normal font for size %.1fpt to be loaded", size)
		}

		if fontPair.Bold == nil {
			t.Errorf("expected bold font for size %.1fpt to be loaded", size)
		}
	}
}

// TestNewTextRenderer_WithoutFontSizes tests backward compatibility when FontSizes is not provided
func TestNewTextRenderer_WithoutFontSizes(t *testing.T) {
	config := &RendererConfig{
		PixelWidth:  384,
		FontPath:    "", // Use system font
		FontSize:    14.0,
		LineSpacing: 4,
		Margin:      8,
		// FontSizes is nil - should use default behavior
	}

	renderer, err := NewTextRenderer(config)
	if err != nil {
		t.Fatalf("failed to create renderer without font sizes: %v", err)
	}

	if renderer == nil {
		t.Fatal("expected renderer but got nil")
	}

	// Verify fonts map contains default size
	if len(renderer.fonts) != 1 {
		t.Errorf("expected 1 font size (default), got %d", len(renderer.fonts))
	}

	// Verify default font size is loaded
	fontPair, exists := renderer.fonts[14.0]
	if !exists {
		t.Error("expected default font size 14.0pt to be loaded")
	}

	if fontPair.Normal == nil {
		t.Error("expected normal font for default size to be loaded")
	}

	if fontPair.Bold == nil {
		t.Error("expected bold font for default size to be loaded")
	}
}

// TestRender_WithDifferentFontSizes tests rendering lines with different font sizes
func TestRender_WithDifferentFontSizes(t *testing.T) {
	config := &RendererConfig{
		PixelWidth:  384,
		FontPath:    "", // Use system font
		FontSize:    18.0,
		FontSizes: map[string]float64{
			"normal": 18.0,
			"header": 22.0,
			"total":  20.0,
		},
		LineSpacing: 4,
		Margin:      8,
	}

	renderer, err := NewTextRenderer(config)
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	lines := []LineFormat{
		{Text: "HÓA ĐƠN BÁN HÀNG", Bold: true, Alignment: AlignCenter, IsSeparator: false, FontSize: 22.0},
		{Text: "===", Bold: false, Alignment: AlignCenter, IsSeparator: true},
		{Text: "Cafe Latte", Bold: false, Alignment: AlignLeft, IsSeparator: false, FontSize: 18.0},
		{Text: "Banh Mi Thit", Bold: false, Alignment: AlignLeft, IsSeparator: false, FontSize: 18.0},
		{Text: "===", Bold: false, Alignment: AlignCenter, IsSeparator: true},
		{Text: "TỔNG CỘNG: 125,000 VND", Bold: true, Alignment: AlignRight, IsSeparator: false, FontSize: 20.0},
	}

	img, err := renderer.Render(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if img == nil {
		t.Fatal("expected image but got nil")
	}

	// Verify image has content
	bounds := img.Bounds()
	if bounds.Dx() != 384 {
		t.Errorf("expected width 384, got %d", bounds.Dx())
	}

	if bounds.Dy() <= 0 {
		t.Error("expected positive height")
	}

	// Verify image contains text (has black pixels)
	hasBlackPixels := false
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if img.GrayAt(x, y).Y < 255 {
				hasBlackPixels = true
				break
			}
		}
		if hasBlackPixels {
			break
		}
	}

	if !hasBlackPixels {
		t.Error("expected image to contain black pixels (text)")
	}
}

// TestRender_WithFontSizeFallback tests that renderer falls back to default font when size not found
func TestRender_WithFontSizeFallback(t *testing.T) {
	config := &RendererConfig{
		PixelWidth:  384,
		FontPath:    "", // Use system font
		FontSize:    18.0,
		FontSizes: map[string]float64{
			"normal": 18.0,
		},
		LineSpacing: 4,
		Margin:      8,
	}

	renderer, err := NewTextRenderer(config)
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	// Try to render with a font size that wasn't loaded (should fallback to default)
	lines := []LineFormat{
		{Text: "Test with unknown size", Bold: false, Alignment: AlignLeft, IsSeparator: false, FontSize: 24.0},
	}

	img, err := renderer.Render(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if img == nil {
		t.Fatal("expected image but got nil")
	}
}

// TestCalculateImageHeight_WithDifferentFontSizes tests height calculation with mixed font sizes
func TestCalculateImageHeight_WithDifferentFontSizes(t *testing.T) {
	config := &RendererConfig{
		PixelWidth:  384,
		FontPath:    "", // Use system font
		FontSize:    18.0,
		FontSizes: map[string]float64{
			"normal": 18.0,
			"header": 22.0,
			"total":  20.0,
		},
		LineSpacing: 4,
		Margin:      8,
	}

	renderer, err := NewTextRenderer(config)
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	lines := []LineFormat{
		{Text: "Header", Bold: true, Alignment: AlignCenter, IsSeparator: false, FontSize: 22.0},
		{Text: "Normal text", Bold: false, Alignment: AlignLeft, IsSeparator: false, FontSize: 18.0},
		{Text: "Total", Bold: true, Alignment: AlignRight, IsSeparator: false, FontSize: 20.0},
	}

	height := renderer.calculateImageHeight(lines)

	if height <= 0 {
		t.Errorf("expected positive height, got %d", height)
	}

	// Height should be greater than with all same-size fonts
	// because larger fonts take more space
	minHeight := 2 * renderer.margin
	if height < minHeight {
		t.Errorf("expected height >= %d (margins), got %d", minHeight, height)
	}
}
