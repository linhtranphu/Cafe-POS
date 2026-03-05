package printing

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

// LogoRenderer handles loading and rendering logos for bill printing
type LogoRenderer struct {
	maxWidthPercent float64 // Percentage of paper width for logo (default: 25%)
	margin          int     // Margin around logo in pixels
}

// NewLogoRenderer creates a new LogoRenderer instance
func NewLogoRenderer(maxWidthPercent float64, margin int) *LogoRenderer {
	// Set default values if not provided
	if maxWidthPercent <= 0 || maxWidthPercent > 100 {
		maxWidthPercent = 25.0 // Default to 25% of paper width
	}
	if margin < 0 {
		margin = 0
	}

	return &LogoRenderer{
		maxWidthPercent: maxWidthPercent,
		margin:          margin,
	}
}

// RenderLogo loads a logo from the specified path and renders it as a grayscale image
// The logo is resized to fit within maxWidthPercent of the paper width
// Returns the rendered logo image or an error if loading/processing fails
func (r *LogoRenderer) RenderLogo(logoPath string, paperWidth int) (*image.Gray, error) {
	// Validate inputs
	if logoPath == "" {
		return nil, fmt.Errorf("logo rendering error: logo path is empty")
	}
	if paperWidth <= 0 {
		return nil, fmt.Errorf("logo rendering error: invalid paper width %d (must be positive)", paperWidth)
	}

	// Check if file exists
	if _, err := os.Stat(logoPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("logo file not found: %s does not exist", logoPath)
	} else if err != nil {
		return nil, fmt.Errorf("logo file access error: cannot access %s: %w", logoPath, err)
	}

	// Open and decode the image file
	img, err := r.loadImage(logoPath)
	if err != nil {
		return nil, fmt.Errorf("logo loading error: %w", err)
	}

	// Calculate maximum logo width
	maxLogoWidth := int(float64(paperWidth) * r.maxWidthPercent / 100.0)
	if maxLogoWidth <= 0 {
		return nil, fmt.Errorf("logo rendering error: calculated max width is invalid (%d pixels)", maxLogoWidth)
	}

	// Resize logo if necessary
	resizedLogo := r.resizeLogo(img, maxLogoWidth)
	if resizedLogo == nil {
		return nil, fmt.Errorf("logo resize error: failed to resize logo")
	}

	return resizedLogo, nil
}

// loadImage loads an image from the specified file path
// Supports PNG, JPG, and JPEG formats
func (r *LogoRenderer) loadImage(path string) (image.Image, error) {
	// Open file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	// Determine file format from extension
	ext := strings.ToLower(filepath.Ext(path))
	
	var img image.Image
	switch ext {
	case ".png":
		img, err = png.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("failed to decode PNG file %s (file may be corrupt): %w", path, err)
		}
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("failed to decode JPEG file %s (file may be corrupt): %w", path, err)
		}
	default:
		return nil, fmt.Errorf("unsupported image format: %s (only PNG, JPG, and JPEG are supported)", ext)
	}

	return img, nil
}

// resizeLogo resizes the logo to fit within maxWidth while maintaining aspect ratio
// Uses bilinear interpolation for smooth scaling
// Returns a grayscale image
func (r *LogoRenderer) resizeLogo(img image.Image, maxWidth int) *image.Gray {
	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	// If image already fits, just convert to grayscale
	if originalWidth <= maxWidth {
		return r.convertToGrayscale(img)
	}

	// Calculate new dimensions maintaining aspect ratio
	scale := float64(maxWidth) / float64(originalWidth)
	newWidth := maxWidth
	newHeight := int(float64(originalHeight) * scale)

	// Ensure minimum height
	if newHeight <= 0 {
		newHeight = 1
	}

	// Create destination image (RGBA for intermediate processing)
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Resize using bilinear interpolation
	draw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	// Convert to grayscale
	return r.convertToGrayscale(dst)
}

// convertToGrayscale converts any image to grayscale
// All pixels will have R=G=B values
func (r *LogoRenderer) convertToGrayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get original color
			originalColor := img.At(x, y)
			
			// Convert to grayscale using the standard luminance formula
			// Gray = 0.299*R + 0.587*G + 0.114*B
			grayColor := color.GrayModel.Convert(originalColor)
			
			// Set grayscale pixel
			grayImg.Set(x, y, grayColor)
		}
	}

	return grayImg
}
