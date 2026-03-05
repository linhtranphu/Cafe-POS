package printing

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// Helper function to create a test PNG image
func createTestPNG(path string, width, height int, fillColor color.Color) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Fill with specified color
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, fillColor)
		}
	}
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	// Save as PNG
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	
	return png.Encode(file, img)
}

func TestNewLogoRenderer(t *testing.T) {
	tests := []struct {
		name            string
		maxWidthPercent float64
		margin          int
		expectedPercent float64
		expectedMargin  int
	}{
		{
			name:            "Valid parameters",
			maxWidthPercent: 25.0,
			margin:          10,
			expectedPercent: 25.0,
			expectedMargin:  10,
		},
		{
			name:            "Zero percent defaults to 25%",
			maxWidthPercent: 0,
			margin:          10,
			expectedPercent: 25.0,
			expectedMargin:  10,
		},
		{
			name:            "Negative percent defaults to 25%",
			maxWidthPercent: -10,
			margin:          10,
			expectedPercent: 25.0,
			expectedMargin:  10,
		},
		{
			name:            "Over 100% defaults to 25%",
			maxWidthPercent: 150,
			margin:          10,
			expectedPercent: 25.0,
			expectedMargin:  10,
		},
		{
			name:            "Negative margin defaults to 0",
			maxWidthPercent: 25.0,
			margin:          -5,
			expectedPercent: 25.0,
			expectedMargin:  0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewLogoRenderer(tt.maxWidthPercent, tt.margin)
			
			if renderer == nil {
				t.Fatal("NewLogoRenderer returned nil")
			}
			
			if renderer.maxWidthPercent != tt.expectedPercent {
				t.Errorf("Expected maxWidthPercent %.1f, got %.1f", tt.expectedPercent, renderer.maxWidthPercent)
			}
			
			if renderer.margin != tt.expectedMargin {
				t.Errorf("Expected margin %d, got %d", tt.expectedMargin, renderer.margin)
			}
		})
	}
}

func TestRenderLogo_MissingFile(t *testing.T) {
	renderer := NewLogoRenderer(25.0, 10)
	
	_, err := renderer.RenderLogo("/nonexistent/path/logo.png", 384)
	if err == nil {
		t.Error("Expected error for missing file, got nil")
	}
	
	if !os.IsNotExist(err) && err.Error() == "" {
		t.Errorf("Expected 'file not found' error, got: %v", err)
	}
}

func TestRenderLogo_InvalidPaperWidth(t *testing.T) {
	renderer := NewLogoRenderer(25.0, 10)
	
	// Create a temporary test image
	tmpDir := t.TempDir()
	logoPath := filepath.Join(tmpDir, "test_logo.png")
	if err := createTestPNG(logoPath, 100, 100, color.Black); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	
	_, err := renderer.RenderLogo(logoPath, 0)
	if err == nil {
		t.Error("Expected error for invalid paper width, got nil")
	}
	
	_, err = renderer.RenderLogo(logoPath, -100)
	if err == nil {
		t.Error("Expected error for negative paper width, got nil")
	}
}

func TestRenderLogo_EmptyPath(t *testing.T) {
	renderer := NewLogoRenderer(25.0, 10)
	
	_, err := renderer.RenderLogo("", 384)
	if err == nil {
		t.Error("Expected error for empty path, got nil")
	}
}

func TestRenderLogo_ValidPNG(t *testing.T) {
	renderer := NewLogoRenderer(25.0, 10)
	
	// Create a temporary test image
	tmpDir := t.TempDir()
	logoPath := filepath.Join(tmpDir, "test_logo.png")
	if err := createTestPNG(logoPath, 200, 100, color.RGBA{R: 100, G: 150, B: 200, A: 255}); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	
	paperWidth := 384
	result, err := renderer.RenderLogo(logoPath, paperWidth)
	if err != nil {
		t.Fatalf("RenderLogo failed: %v", err)
	}
	
	if result == nil {
		t.Fatal("RenderLogo returned nil image")
	}
	
	// Check that result is grayscale
	bounds := result.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		t.Errorf("Invalid result dimensions: %dx%d", bounds.Dx(), bounds.Dy())
	}
	
	// Verify logo width constraint (should be <= 25% of paper width)
	maxWidth := int(float64(paperWidth) * 25.0 / 100.0)
	if bounds.Dx() > maxWidth {
		t.Errorf("Logo width %d exceeds max width %d", bounds.Dx(), maxWidth)
	}
}

func TestRenderLogo_ResizeLargeLogo(t *testing.T) {
	renderer := NewLogoRenderer(25.0, 10)
	
	// Create a large test image (larger than 25% of paper width)
	tmpDir := t.TempDir()
	logoPath := filepath.Join(tmpDir, "large_logo.png")
	if err := createTestPNG(logoPath, 500, 250, color.White); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	
	paperWidth := 384
	maxWidth := int(float64(paperWidth) * 25.0 / 100.0) // 96 pixels
	
	result, err := renderer.RenderLogo(logoPath, paperWidth)
	if err != nil {
		t.Fatalf("RenderLogo failed: %v", err)
	}
	
	// Verify logo was resized
	bounds := result.Bounds()
	if bounds.Dx() > maxWidth {
		t.Errorf("Logo width %d exceeds max width %d", bounds.Dx(), maxWidth)
	}
	
	// Verify aspect ratio is maintained (approximately)
	originalAspect := 500.0 / 250.0 // 2.0
	resultAspect := float64(bounds.Dx()) / float64(bounds.Dy())
	
	// Allow 5% tolerance for rounding
	if resultAspect < originalAspect*0.95 || resultAspect > originalAspect*1.05 {
		t.Errorf("Aspect ratio not maintained: original %.2f, result %.2f", originalAspect, resultAspect)
	}
}

func TestConvertToGrayscale(t *testing.T) {
	renderer := NewLogoRenderer(25.0, 10)
	
	// Create a color image
	colorImg := image.NewRGBA(image.Rect(0, 0, 10, 10))
	
	// Fill with different colors
	colorImg.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})     // Red
	colorImg.Set(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255})     // Green
	colorImg.Set(2, 0, color.RGBA{R: 0, G: 0, B: 255, A: 255})     // Blue
	colorImg.Set(3, 0, color.RGBA{R: 255, G: 255, B: 255, A: 255}) // White
	colorImg.Set(4, 0, color.RGBA{R: 0, G: 0, B: 0, A: 255})       // Black
	
	// Convert to grayscale
	grayImg := renderer.convertToGrayscale(colorImg)
	
	if grayImg == nil {
		t.Fatal("convertToGrayscale returned nil")
	}
	
	// Verify all pixels are grayscale (R=G=B)
	bounds := grayImg.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := grayImg.At(x, y)
			r, g, b, _ := c.RGBA()
			
			// In grayscale, R=G=B
			if r != g || g != b {
				t.Errorf("Pixel at (%d,%d) is not grayscale: R=%d, G=%d, B=%d", x, y, r>>8, g>>8, b>>8)
			}
		}
	}
}

func TestResizeLogo_SmallImage(t *testing.T) {
	renderer := NewLogoRenderer(25.0, 10)
	
	// Create a small image that doesn't need resizing
	smallImg := image.NewRGBA(image.Rect(0, 0, 50, 50))
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			smallImg.Set(x, y, color.Gray{Y: 128})
		}
	}
	
	maxWidth := 100
	result := renderer.resizeLogo(smallImg, maxWidth)
	
	if result == nil {
		t.Fatal("resizeLogo returned nil")
	}
	
	// Image should not be enlarged, just converted to grayscale
	bounds := result.Bounds()
	if bounds.Dx() != 50 || bounds.Dy() != 50 {
		t.Errorf("Expected dimensions 50x50, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestResizeLogo_LargeImage(t *testing.T) {
	renderer := NewLogoRenderer(25.0, 10)
	
	// Create a large image that needs resizing
	largeImg := image.NewRGBA(image.Rect(0, 0, 200, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 200; x++ {
			largeImg.Set(x, y, color.Gray{Y: 200})
		}
	}
	
	maxWidth := 100
	result := renderer.resizeLogo(largeImg, maxWidth)
	
	if result == nil {
		t.Fatal("resizeLogo returned nil")
	}
	
	// Image should be resized to maxWidth
	bounds := result.Bounds()
	if bounds.Dx() != maxWidth {
		t.Errorf("Expected width %d, got %d", maxWidth, bounds.Dx())
	}
	
	// Height should be scaled proportionally (200x100 -> 100x50)
	expectedHeight := 50
	if bounds.Dy() != expectedHeight {
		t.Errorf("Expected height %d, got %d", expectedHeight, bounds.Dy())
	}
}

func TestLoadImage_UnsupportedFormat(t *testing.T) {
	renderer := NewLogoRenderer(25.0, 10)
	
	// Create a file with unsupported extension
	tmpDir := t.TempDir()
	unsupportedPath := filepath.Join(tmpDir, "test.bmp")
	
	// Create an empty file
	file, err := os.Create(unsupportedPath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()
	
	_, err = renderer.loadImage(unsupportedPath)
	if err == nil {
		t.Error("Expected error for unsupported format, got nil")
	}
}
