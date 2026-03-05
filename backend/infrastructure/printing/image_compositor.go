package printing

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

// ImageCompositor combines logo and text content into a single image for printing
type ImageCompositor struct {
	paperWidth int // Width of the paper in pixels
	margin     int // Margin around content in pixels
}

// NewImageCompositor creates a new ImageCompositor instance
func NewImageCompositor(paperWidth, margin int) *ImageCompositor {
	// Validate inputs
	if paperWidth <= 0 {
		paperWidth = 576 // Default to 80mm paper width
	}
	if margin < 0 {
		margin = 0
	}

	return &ImageCompositor{
		paperWidth: paperWidth,
		margin:     margin,
	}
}

// Compose combines logo (optional) and text content into a single grayscale image
// If logo is nil, returns the text content as-is
// Layout: [margin] -> [logo] -> [spacing: 20px] -> [text content] -> [margin]
func (c *ImageCompositor) Compose(logo *image.Gray, textContent *image.Gray) (*image.Gray, error) {
	// Validate text content (required)
	if textContent == nil {
		return nil, fmt.Errorf("image composition error: text content cannot be nil")
	}

	// If no logo, return text content as-is
	if logo == nil {
		return textContent, nil
	}

	// Calculate total height needed
	totalHeight := c.calculateTotalHeight(logo, textContent)
	if totalHeight <= 0 {
		return nil, fmt.Errorf("image composition error: calculated total height is invalid (%d pixels)", totalHeight)
	}

	// Create combined image
	combinedImg := image.NewGray(image.Rect(0, 0, c.paperWidth, totalHeight))
	if combinedImg == nil {
		return nil, fmt.Errorf("image composition error: failed to allocate image buffer (%dx%d pixels)", c.paperWidth, totalHeight)
	}

	// Fill with white background
	draw.Draw(combinedImg, combinedImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Draw logo at top with margin
	y := c.margin
	c.drawLogo(combinedImg, logo, c.margin, y)
	y += logo.Bounds().Dy()

	// Add spacing between logo and text content
	y += 20

	// Draw text content
	c.drawTextContent(combinedImg, textContent, 0, y)

	return combinedImg, nil
}

// calculateTotalHeight calculates the total height needed for the combined image
// Layout: [margin] + [logo height] + [spacing: 20px] + [text content height]
func (c *ImageCompositor) calculateTotalHeight(logo *image.Gray, textContent *image.Gray) int {
	height := c.margin // Top margin

	// Add logo height if present
	if logo != nil {
		height += logo.Bounds().Dy()
		height += 20 // Spacing between logo and text
	}

	// Add text content height
	if textContent != nil {
		height += textContent.Bounds().Dy()
	}

	return height
}

// drawLogo draws the logo onto the destination image at the specified position
// Logo is aligned to the left with the specified x offset
func (c *ImageCompositor) drawLogo(dst *image.Gray, logo *image.Gray, x, y int) {
	if logo == nil || dst == nil {
		return
	}

	// Get logo bounds
	logoBounds := logo.Bounds()
	logoWidth := logoBounds.Dx()
	logoHeight := logoBounds.Dy()

	// Draw logo pixel by pixel
	for ly := 0; ly < logoHeight; ly++ {
		for lx := 0; lx < logoWidth; lx++ {
			// Calculate destination coordinates
			dstX := x + lx
			dstY := y + ly

			// Check bounds
			if dstX >= 0 && dstX < dst.Bounds().Dx() && dstY >= 0 && dstY < dst.Bounds().Dy() {
				// Get pixel from logo
				pixel := logo.GrayAt(logoBounds.Min.X+lx, logoBounds.Min.Y+ly)
				// Set pixel in destination
				dst.SetGray(dstX, dstY, pixel)
			}
		}
	}
}

// drawTextContent draws the text content onto the destination image at the specified position
func (c *ImageCompositor) drawTextContent(dst *image.Gray, textContent *image.Gray, x, y int) {
	if textContent == nil || dst == nil {
		return
	}

	// Get text content bounds
	textBounds := textContent.Bounds()
	textWidth := textBounds.Dx()
	textHeight := textBounds.Dy()

	// Draw text content pixel by pixel
	for ty := 0; ty < textHeight; ty++ {
		for tx := 0; tx < textWidth; tx++ {
			// Calculate destination coordinates
			dstX := x + tx
			dstY := y + ty

			// Check bounds
			if dstX >= 0 && dstX < dst.Bounds().Dx() && dstY >= 0 && dstY < dst.Bounds().Dy() {
				// Get pixel from text content
				pixel := textContent.GrayAt(textBounds.Min.X+tx, textBounds.Min.Y+ty)
				// Set pixel in destination
				dst.SetGray(dstX, dstY, pixel)
			}
		}
	}
}
