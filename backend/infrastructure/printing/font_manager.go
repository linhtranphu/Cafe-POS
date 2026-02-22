package printing

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// FontConfig represents configuration for font management
type FontConfig struct {
	FontPaths  []string // Paths to TrueType font files
	NormalSize float64  // Font size for normal text (in points)
	BoldSize   float64  // Font size for bold text (in points)
}

// RendererConfig represents configuration for the text renderer
type RendererConfig struct {
	PixelWidth  int     // Width of the printable area in pixels (384 for 58mm, 576 for 80mm)
	FontPath    string  // Path to the TrueType font file to use for rendering
	FontSize    float64 // Font size in points for normal text
	LineSpacing int     // Vertical spacing between lines in pixels
	Margin      int     // Left and right margin in pixels
}

// FontManager manages font loading and provides font faces for rendering
type FontManager struct {
	fontPaths  []string
	normalSize float64
	boldSize   float64
}

// NewFontManager creates a new FontManager instance
func NewFontManager(config *FontConfig) *FontManager {
	return &FontManager{
		fontPaths:  config.FontPaths,
		normalSize: config.NormalSize,
		boldSize:   config.BoldSize,
	}
}

// LoadFonts loads normal and bold font faces
// Returns normal font face, bold font face, and error if any
func (m *FontManager) LoadFonts() (normalFont, boldFont font.Face, err error) {
	// Try to load from configured font paths first
	for _, fontPath := range m.fontPaths {
		if fontPath == "" {
			continue // Skip empty paths
		}
		
		normalFont, err = m.loadFontFromPath(fontPath, m.normalSize)
		if err == nil {
			// Successfully loaded normal font, try to load bold variant
			// For now, use the same font with larger size for bold
			boldFont, err = m.loadFontFromPath(fontPath, m.boldSize)
			if err == nil {
				return normalFont, boldFont, nil
			}
			// If bold loading failed, wrap error with context
			return nil, nil, fmt.Errorf("font loading error: normal font loaded successfully from %s, but bold variant failed: %w", fontPath, err)
		}
	}

	// If configured paths failed, try to find system fonts
	systemFontPath, err := m.findSystemFont()
	if err != nil {
		return nil, nil, fmt.Errorf("font loading error: no configured fonts found (tried %d paths) and system font discovery failed: %w", len(m.fontPaths), err)
	}

	// Load from system font
	normalFont, err = m.loadFontFromPath(systemFontPath, m.normalSize)
	if err != nil {
		return nil, nil, fmt.Errorf("font loading error: failed to load normal font from system font at %s (size: %.1fpt): %w", systemFontPath, m.normalSize, err)
	}

	boldFont, err = m.loadFontFromPath(systemFontPath, m.boldSize)
	if err != nil {
		return nil, nil, fmt.Errorf("font loading error: failed to load bold variant from system font at %s (size: %.1fpt): %w", systemFontPath, m.boldSize, err)
	}

	return normalFont, boldFont, nil
}

// loadFontFromPath loads a TrueType font from the specified file path
func (m *FontManager) loadFontFromPath(path string, size float64) (font.Face, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("font file not found: %s does not exist", path)
	} else if err != nil {
		return nil, fmt.Errorf("font file access error: cannot access %s: %w", path, err)
	}

	// Read font file
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("font file read error: failed to read %s: %w", path, err)
	}

	// Parse TrueType/OpenType font
	parsedFont, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("font parsing error: %s is not a valid TrueType/OpenType font: %w", path, err)
	}

	// Create font face with specified size and anti-aliasing
	// DPI is set to 72 (standard for font rendering)
	// HintingFull enables anti-aliasing for smooth character edges (Requirement 9.1)
	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull, // Enables anti-aliasing for improved readability
	})
	if err != nil {
		return nil, fmt.Errorf("font face creation error: failed to create font face from %s at size %.1fpt: %w", path, size, err)
	}

	return face, nil
}

// findSystemFont searches for Vietnamese-compatible system fonts
// Returns the path to the first found font
func (m *FontManager) findSystemFont() (string, error) {
	var searchPaths []string

	// Define search paths based on operating system
	switch runtime.GOOS {
	case "windows":
		// Windows font directories
		windowsDir := os.Getenv("WINDIR")
		if windowsDir == "" {
			windowsDir = "C:\\Windows"
		}
		searchPaths = []string{
			filepath.Join(windowsDir, "Fonts", "arialuni.ttf"),  // Arial Unicode MS
			filepath.Join(windowsDir, "Fonts", "arial.ttf"),     // Arial (fallback)
			filepath.Join(windowsDir, "Fonts", "DejaVuSans.ttf"), // DejaVu Sans
		}

	case "darwin":
		// macOS font directories
		searchPaths = []string{
			"/Library/Fonts/Arial Unicode.ttf",                    // Arial Unicode MS
			"/System/Library/Fonts/Supplemental/Arial Unicode.ttf", // Arial Unicode MS (alternative)
			"/Library/Fonts/Arial.ttf",                            // Arial (fallback)
			"/System/Library/Fonts/Supplemental/Arial.ttf",        // Arial (alternative)
		}

	case "linux":
		// Linux font directories
		searchPaths = []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",           // DejaVu Sans
			"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf", // Liberation Sans
			"/usr/share/fonts/truetype/noto/NotoSans-Regular.ttf",       // Noto Sans
			"/usr/share/fonts/truetype/roboto/Roboto-Regular.ttf",       // Roboto
			"/usr/share/fonts/TTF/DejaVuSans.ttf",                       // DejaVu Sans (Arch Linux)
		}

	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Search for the first available font
	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("no Vietnamese-compatible system fonts found. Searched paths: %v. Please install Arial Unicode MS, Roboto, or DejaVu Sans", searchPaths)
}
