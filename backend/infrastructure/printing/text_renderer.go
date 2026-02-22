package printing

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// TextRenderer renders formatted text lines to a bitmap image
type TextRenderer struct {
	pixelWidth  int
	normalFont  font.Face
	boldFont    font.Face
	fontSize    float64
	lineSpacing int
	margin      int
}

// NewTextRenderer creates a new TextRenderer instance with font loading
func NewTextRenderer(config *RendererConfig) (*TextRenderer, error) {
	if config == nil {
		return nil, fmt.Errorf("text renderer initialization error: renderer config cannot be nil")
	}

	if config.PixelWidth <= 0 {
		return nil, fmt.Errorf("text renderer initialization error: invalid pixel width %d (must be positive)", config.PixelWidth)
	}

	if config.FontSize <= 0 {
		return nil, fmt.Errorf("text renderer initialization error: invalid font size %.1f (must be positive)", config.FontSize)
	}

	// Create font manager to load fonts
	fontConfig := &FontConfig{
		FontPaths:  []string{config.FontPath},
		NormalSize: config.FontSize,
		BoldSize:   config.FontSize * 1.2, // Bold font is 20% larger
	}

	fontManager := NewFontManager(fontConfig)
	normalFont, boldFont, err := fontManager.LoadFonts()
	if err != nil {
		return nil, fmt.Errorf("text renderer initialization error: %w", err)
	}

	return &TextRenderer{
		pixelWidth:  config.PixelWidth,
		normalFont:  normalFont,
		boldFont:    boldFont,
		fontSize:    config.FontSize,
		lineSpacing: config.LineSpacing,
		margin:      config.Margin,
	}, nil
}

// Render renders formatted lines to a grayscale image
func (r *TextRenderer) Render(lines []LineFormat) (*image.Gray, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("text rendering error: cannot render empty content (no lines provided)")
	}

	// Calculate total image height
	height := r.calculateImageHeight(lines)
	if height <= 0 {
		return nil, fmt.Errorf("text rendering error: calculated image height is invalid (%d pixels)", height)
	}

	// Create grayscale image
	img := image.NewGray(image.Rect(0, 0, r.pixelWidth, height))
	if img == nil {
		return nil, fmt.Errorf("text rendering error: failed to allocate image buffer (%dx%d pixels)", r.pixelWidth, height)
	}

	// Fill with white background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Draw each line
	y := r.margin
	for i, line := range lines {
		newY := r.drawLine(img, line, y)
		if newY < y {
			return nil, fmt.Errorf("text rendering error: line %d rendering failed (invalid Y position)", i+1)
		}
		y = newY
	}

	return img, nil
}

// measureText measures the width of text in pixels
func (r *TextRenderer) measureText(text string, fontFace font.Face) int {
	if text == "" {
		return 0
	}

	// Measure text width using font metrics
	drawer := &font.Drawer{
		Face: fontFace,
	}

	advance := drawer.MeasureString(text)
	return advance.Ceil()
}

// calculateImageHeight calculates the total height needed for all lines
// Optimized to minimize paper usage with tight bounds around content (Requirement 9.3)
func (r *TextRenderer) calculateImageHeight(lines []LineFormat) int {
	height := r.margin // Top margin

	for i, line := range lines {
		if line.Text == "" {
			// Empty line: add minimal spacing
			height += r.lineSpacing / 2
		} else if line.IsSeparator {
			// Separator line: add minimal padding above and below
			height += r.lineSpacing / 2
			height += 2 // Separator line height
			height += r.lineSpacing / 2
		} else {
			// Regular text line: calculate actual text height
			fontFace := r.normalFont
			if line.Bold {
				fontFace = r.boldFont
			}

			metrics := fontFace.Metrics()
			// Use Height only (ascent + descent) without extra descent
			lineHeight := metrics.Height.Ceil()
			
			// Add line spacing only between lines, not after the last line
			if i < len(lines)-1 {
				height += lineHeight + r.lineSpacing
			} else {
				height += lineHeight
			}
		}
	}

	// Add minimal bottom padding (10 pixels) instead of full margin
	height += 10
	return height
}

// drawLine draws a single line and returns the next Y position
func (r *TextRenderer) drawLine(img *image.Gray, line LineFormat, y int) int {
	// Handle empty lines
	if line.Text == "" {
		return y + r.lineSpacing/2
	}

	// Handle separator lines
	if line.IsSeparator {
		y += r.lineSpacing / 2
		separatorY := y + 1
		// Draw horizontal line (centered, with margins)
		for x := r.margin; x < r.pixelWidth-r.margin; x++ {
			img.SetGray(x, separatorY, color.Gray{Y: 0})
		}
		return y + 2 + r.lineSpacing/2
	}

	// Select font based on bold flag
	fontFace := r.normalFont
	if line.Bold {
		fontFace = r.boldFont
	}

	// Handle text wrapping
	wrappedLines := r.wrapText(line.Text, fontFace)

	// Draw each wrapped line
	for _, wrappedText := range wrappedLines {
		// Measure text width
		textWidth := r.measureText(wrappedText, fontFace)

		// Calculate X position based on alignment
		var x int
		switch line.Alignment {
		case AlignLeft:
			x = r.margin
		case AlignCenter:
			x = (r.pixelWidth - textWidth) / 2
		case AlignRight:
			x = r.pixelWidth - textWidth - r.margin
		}

		// Ensure x is not negative
		if x < r.margin {
			x = r.margin
		}

		// Get font metrics for baseline calculation
		metrics := fontFace.Metrics()
		ascent := metrics.Ascent.Ceil()

		// Draw text
		drawer := &font.Drawer{
			Dst:  img,
			Src:  image.Black,
			Face: fontFace,
			Dot:  fixed.P(x, y+ascent),
		}
		drawer.DrawString(wrappedText)

		// Move to next line - use optimized height calculation
		lineHeight := metrics.Height.Ceil()
		y += lineHeight + r.lineSpacing
	}

	return y
}

// wrapText wraps text at word boundaries to fit within pixel width
func (r *TextRenderer) wrapText(text string, fontFace font.Face) []string {
	// Calculate available width (excluding margins)
	availableWidth := r.pixelWidth - 2*r.margin

	// If text fits, return as-is
	textWidth := r.measureText(text, fontFace)
	if textWidth <= availableWidth {
		return []string{text}
	}

	// Split text into words
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
	}

	var lines []string
	var currentLine strings.Builder
	currentWidth := 0

	for i, word := range words {
		wordWidth := r.measureText(word, fontFace)
		spaceWidth := r.measureText(" ", fontFace)

		// Check if adding this word would exceed available width
		widthWithWord := currentWidth
		if currentLine.Len() > 0 {
			widthWithWord += spaceWidth
		}
		widthWithWord += wordWidth

		if widthWithWord > availableWidth && currentLine.Len() > 0 {
			// Current line is full, start new line
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
			currentWidth = wordWidth
		} else {
			// Add word to current line
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
				currentWidth += spaceWidth
			}
			currentLine.WriteString(word)
			currentWidth += wordWidth
		}

		// If this is the last word, add the current line
		if i == len(words)-1 && currentLine.Len() > 0 {
			lines = append(lines, currentLine.String())
		}
	}

	// If no lines were created (shouldn't happen), return original text
	if len(lines) == 0 {
		return []string{text}
	}

	return lines
}
