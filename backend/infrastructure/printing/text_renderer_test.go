package printing

import (
	"image"
	"strings"
	"testing"
)

func TestNewTextRenderer(t *testing.T) {
	tests := []struct {
		name        string
		config      *RendererConfig
		expectError bool
	}{
		{
			name: "valid config with system font",
			config: &RendererConfig{
				PixelWidth:  384,
				FontPath:    "", // Will use system font
				FontSize:    14.0,
				LineSpacing: 4,
				Margin:      8,
			},
			expectError: false,
		},
		{
			name:        "nil config",
			config:      nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer, err := NewTextRenderer(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if renderer == nil {
				t.Fatal("expected renderer but got nil")
			}

			if renderer.pixelWidth != tt.config.PixelWidth {
				t.Errorf("expected pixelWidth %d, got %d", tt.config.PixelWidth, renderer.pixelWidth)
			}
		})
	}
}

func TestRender_EmptyContent(t *testing.T) {
	renderer := createTestRenderer(t)

	_, err := renderer.Render([]LineFormat{})
	if err == nil {
		t.Error("expected error for empty content")
	}
}

func TestRender_SingleLine(t *testing.T) {
	renderer := createTestRenderer(t)

	lines := []LineFormat{
		{Text: "Hello World", Bold: false, Alignment: AlignLeft, IsSeparator: false},
	}

	img, err := renderer.Render(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if img == nil {
		t.Fatal("expected image but got nil")
	}

	// Check image dimensions
	bounds := img.Bounds()
	if bounds.Dx() != 384 {
		t.Errorf("expected width 384, got %d", bounds.Dx())
	}

	if bounds.Dy() <= 0 {
		t.Error("expected positive height")
	}
}

func TestRender_VietnameseText(t *testing.T) {
	renderer := createTestRenderer(t)

	lines := []LineFormat{
		{Text: "Xin chào", Bold: false, Alignment: AlignLeft, IsSeparator: false},
		{Text: "Cảm ơn quý khách", Bold: true, Alignment: AlignCenter, IsSeparator: false},
		{Text: "Hẹn gặp lại", Bold: false, Alignment: AlignRight, IsSeparator: false},
	}

	img, err := renderer.Render(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if img == nil {
		t.Fatal("expected image but got nil")
	}

	// Verify image has content (not all white)
	hasBlackPixels := false
	bounds := img.Bounds()
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

func TestRender_MultipleLines(t *testing.T) {
	renderer := createTestRenderer(t)

	lines := []LineFormat{
		{Text: "Line 1", Bold: false, Alignment: AlignLeft, IsSeparator: false},
		{Text: "Line 2", Bold: true, Alignment: AlignCenter, IsSeparator: false},
		{Text: "Line 3", Bold: false, Alignment: AlignRight, IsSeparator: false},
	}

	img, err := renderer.Render(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if img == nil {
		t.Fatal("expected image but got nil")
	}
}

func TestRender_EmptyLines(t *testing.T) {
	renderer := createTestRenderer(t)

	lines := []LineFormat{
		{Text: "Line 1", Bold: false, Alignment: AlignLeft, IsSeparator: false},
		{Text: "", Bold: false, Alignment: AlignLeft, IsSeparator: false},
		{Text: "Line 2", Bold: false, Alignment: AlignLeft, IsSeparator: false},
	}

	img, err := renderer.Render(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if img == nil {
		t.Fatal("expected image but got nil")
	}
}

func TestRender_SeparatorLines(t *testing.T) {
	renderer := createTestRenderer(t)

	lines := []LineFormat{
		{Text: "Header", Bold: true, Alignment: AlignCenter, IsSeparator: false},
		{Text: "===", Bold: false, Alignment: AlignCenter, IsSeparator: true},
		{Text: "Content", Bold: false, Alignment: AlignLeft, IsSeparator: false},
		{Text: "---", Bold: false, Alignment: AlignCenter, IsSeparator: true},
		{Text: "Footer", Bold: false, Alignment: AlignCenter, IsSeparator: false},
	}

	img, err := renderer.Render(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if img == nil {
		t.Fatal("expected image but got nil")
	}

	// Verify separator lines are drawn (should have horizontal black pixels)
	bounds := img.Bounds()
	foundSeparator := false
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		blackCount := 0
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if img.GrayAt(x, y).Y == 0 {
				blackCount++
			}
		}
		// A separator line should have many consecutive black pixels
		if blackCount > bounds.Dx()/2 {
			foundSeparator = true
			break
		}
	}

	if !foundSeparator {
		t.Error("expected to find separator line in image")
	}
}

func TestRender_BoldText(t *testing.T) {
	renderer := createTestRenderer(t)

	normalLine := []LineFormat{
		{Text: "Normal Text", Bold: false, Alignment: AlignLeft, IsSeparator: false},
	}

	boldLine := []LineFormat{
		{Text: "Bold Text", Bold: true, Alignment: AlignLeft, IsSeparator: false},
	}

	normalImg, err := renderer.Render(normalLine)
	if err != nil {
		t.Fatalf("unexpected error rendering normal text: %v", err)
	}

	boldImg, err := renderer.Render(boldLine)
	if err != nil {
		t.Fatalf("unexpected error rendering bold text: %v", err)
	}

	// Count black pixels in both images
	normalBlackPixels := countBlackPixels(normalImg)
	boldBlackPixels := countBlackPixels(boldImg)

	// Bold text should have more black pixels (thicker strokes)
	if boldBlackPixels <= normalBlackPixels {
		t.Errorf("expected bold text to have more black pixels than normal text, got normal=%d, bold=%d", normalBlackPixels, boldBlackPixels)
	}
}

func TestRender_Alignment(t *testing.T) {
	renderer := createTestRenderer(t)

	tests := []struct {
		name      string
		alignment Alignment
	}{
		{"left aligned", AlignLeft},
		{"center aligned", AlignCenter},
		{"right aligned", AlignRight},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := []LineFormat{
				{Text: "Test", Bold: false, Alignment: tt.alignment, IsSeparator: false},
			}

			img, err := renderer.Render(lines)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if img == nil {
				t.Fatal("expected image but got nil")
			}

			// Verify text was rendered
			if countBlackPixels(img) == 0 {
				t.Error("expected text to be rendered")
			}
		})
	}
}

func TestRender_TextWrapping(t *testing.T) {
	renderer := createTestRenderer(t)

	// Create a very long line that should wrap
	longText := "This is a very long line of text that should definitely wrap to multiple lines when rendered because it exceeds the pixel width of the image"

	lines := []LineFormat{
		{Text: longText, Bold: false, Alignment: AlignLeft, IsSeparator: false},
	}

	img, err := renderer.Render(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if img == nil {
		t.Fatal("expected image but got nil")
	}

	// Image height should be greater than a single line
	// (indicating text was wrapped to multiple lines)
	bounds := img.Bounds()
	singleLineHeight := int(renderer.fontSize) + renderer.lineSpacing + 2*renderer.margin

	if bounds.Dy() <= singleLineHeight {
		t.Errorf("expected wrapped text to have height > %d, got %d", singleLineHeight, bounds.Dy())
	}
}

func TestMeasureText(t *testing.T) {
	renderer := createTestRenderer(t)

	tests := []struct {
		name     string
		text     string
		expected int // We just check it's positive
	}{
		{"empty string", "", 0},
		{"single character", "A", 1},
		{"word", "Hello", 1},
		{"vietnamese text", "Xin chào", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := renderer.measureText(tt.text, renderer.normalFont)

			if tt.text == "" {
				if width != 0 {
					t.Errorf("expected width 0 for empty string, got %d", width)
				}
			} else {
				if width <= 0 {
					t.Errorf("expected positive width for '%s', got %d", tt.text, width)
				}
			}
		})
	}
}

func TestCalculateImageHeight(t *testing.T) {
	renderer := createTestRenderer(t)

	tests := []struct {
		name  string
		lines []LineFormat
	}{
		{
			name: "single line",
			lines: []LineFormat{
				{Text: "Test", Bold: false, Alignment: AlignLeft, IsSeparator: false},
			},
		},
		{
			name: "multiple lines",
			lines: []LineFormat{
				{Text: "Line 1", Bold: false, Alignment: AlignLeft, IsSeparator: false},
				{Text: "Line 2", Bold: false, Alignment: AlignLeft, IsSeparator: false},
				{Text: "Line 3", Bold: false, Alignment: AlignLeft, IsSeparator: false},
			},
		},
		{
			name: "with empty lines",
			lines: []LineFormat{
				{Text: "Line 1", Bold: false, Alignment: AlignLeft, IsSeparator: false},
				{Text: "", Bold: false, Alignment: AlignLeft, IsSeparator: false},
				{Text: "Line 2", Bold: false, Alignment: AlignLeft, IsSeparator: false},
			},
		},
		{
			name: "with separator",
			lines: []LineFormat{
				{Text: "Header", Bold: false, Alignment: AlignLeft, IsSeparator: false},
				{Text: "===", Bold: false, Alignment: AlignCenter, IsSeparator: true},
				{Text: "Content", Bold: false, Alignment: AlignLeft, IsSeparator: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			height := renderer.calculateImageHeight(tt.lines)

			if height <= 0 {
				t.Errorf("expected positive height, got %d", height)
			}

			// Height should at least include margins
			minHeight := 2 * renderer.margin
			if height < minHeight {
				t.Errorf("expected height >= %d (margins), got %d", minHeight, height)
			}
		})
	}
}

func TestWrapText(t *testing.T) {
	renderer := createTestRenderer(t)

	tests := []struct {
		name          string
		text          string
		expectWrapped bool
	}{
		{
			name:          "short text",
			text:          "Hello",
			expectWrapped: false,
		},
		{
			name:          "long text",
			text:          "This is a very long line of text that should definitely wrap to multiple lines when rendered",
			expectWrapped: true,
		},
		{
			name:          "empty text",
			text:          "",
			expectWrapped: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := renderer.wrapText(tt.text, renderer.normalFont)

			if len(lines) == 0 {
				t.Fatal("expected at least one line")
			}

			if tt.expectWrapped {
				if len(lines) == 1 {
					t.Error("expected text to be wrapped to multiple lines")
				}
			} else {
				if len(lines) != 1 {
					t.Errorf("expected single line, got %d lines", len(lines))
				}
			}

			// Verify all words are preserved
			originalWords := len(strings.Fields(tt.text))
			wrappedWords := 0
			for _, line := range lines {
				wrappedWords += len(strings.Fields(line))
			}

			if originalWords != wrappedWords {
				t.Errorf("expected %d words, got %d words after wrapping", originalWords, wrappedWords)
			}
		})
	}
}

func TestWrapText_WordBoundaries(t *testing.T) {
	renderer := createTestRenderer(t)

	text := "Hello World Test"
	lines := renderer.wrapText(text, renderer.normalFont)

	// Verify no words are split
	for _, line := range lines {
		words := strings.Fields(line)
		for _, word := range words {
			if !strings.Contains(text, word) {
				t.Errorf("wrapped line contains partial word: %s", word)
			}
		}
	}
}

// Helper functions

func createTestRenderer(t *testing.T) *TextRenderer {
	t.Helper()

	config := &RendererConfig{
		PixelWidth:  384,
		FontPath:    "", // Use system font
		FontSize:    14.0,
		LineSpacing: 4,
		Margin:      8,
	}

	renderer, err := NewTextRenderer(config)
	if err != nil {
		t.Fatalf("failed to create test renderer: %v", err)
	}

	return renderer
}

func countBlackPixels(img *image.Gray) int {
	count := 0
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if img.GrayAt(x, y).Y < 128 {
				count++
			}
		}
	}

	return count
}
