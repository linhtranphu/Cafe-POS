package printing

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadFonts_ValidTTFFile tests loading a valid TTF file
func TestLoadFonts_ValidTTFFile(t *testing.T) {
	// Skip if no system fonts available
	manager := NewFontManager(&FontConfig{
		FontPaths:  []string{},
		NormalSize: 12.0,
		BoldSize:   14.0,
	})

	systemFontPath, err := manager.findSystemFont()
	if err != nil {
		t.Skip("No system fonts available for testing")
	}

	// Create manager with valid font path
	manager = NewFontManager(&FontConfig{
		FontPaths:  []string{systemFontPath},
		NormalSize: 12.0,
		BoldSize:   14.0,
	})

	normalFont, boldFont, err := manager.LoadFonts()
	require.NoError(t, err, "Should successfully load fonts from valid TTF file")
	assert.NotNil(t, normalFont, "Normal font should not be nil")
	assert.NotNil(t, boldFont, "Bold font should not be nil")

	// Clean up
	if normalFont != nil {
		normalFont.Close()
	}
	if boldFont != nil {
		boldFont.Close()
	}
}

// TestLoadFonts_MissingFont tests error handling for missing font files
func TestLoadFonts_MissingFont(t *testing.T) {
	manager := NewFontManager(&FontConfig{
		FontPaths:  []string{"/nonexistent/path/to/font.ttf"},
		NormalSize: 12.0,
		BoldSize:   14.0,
	})

	normalFont, boldFont, err := manager.LoadFonts()
	
	// Should either find a system font or return an error
	if err != nil {
		assert.Error(t, err, "Should return error when font not found")
		assert.Nil(t, normalFont, "Normal font should be nil on error")
		assert.Nil(t, boldFont, "Bold font should be nil on error")
		assert.Contains(t, err.Error(), "failed to load fonts", "Error message should indicate font loading failure")
	} else {
		// If no error, it means system font was found as fallback
		assert.NotNil(t, normalFont, "Normal font should not be nil when system font is used")
		assert.NotNil(t, boldFont, "Bold font should not be nil when system font is used")
		
		// Clean up
		if normalFont != nil {
			normalFont.Close()
		}
		if boldFont != nil {
			boldFont.Close()
		}
	}
}

// TestLoadFonts_InvalidFontFile tests error handling for corrupted font files
func TestLoadFonts_InvalidFontFile(t *testing.T) {
	// Create a temporary invalid font file
	tempDir := t.TempDir()
	invalidFontPath := filepath.Join(tempDir, "invalid.ttf")
	err := os.WriteFile(invalidFontPath, []byte("This is not a valid TTF file"), 0644)
	require.NoError(t, err, "Should create temporary invalid font file")

	manager := NewFontManager(&FontConfig{
		FontPaths:  []string{invalidFontPath},
		NormalSize: 12.0,
		BoldSize:   14.0,
	})

	normalFont, boldFont, err := manager.LoadFonts()
	
	// Should either find a system font or return an error
	if err != nil {
		assert.Error(t, err, "Should return error for invalid font file")
		assert.Nil(t, normalFont, "Normal font should be nil on error")
		assert.Nil(t, boldFont, "Bold font should be nil on error")
	} else {
		// If no error, it means system font was found as fallback
		assert.NotNil(t, normalFont, "Normal font should not be nil when system font is used")
		assert.NotNil(t, boldFont, "Bold font should not be nil when system font is used")
		
		// Clean up
		if normalFont != nil {
			normalFont.Close()
		}
		if boldFont != nil {
			boldFont.Close()
		}
	}
}

// TestFindSystemFont tests system font discovery
func TestFindSystemFont(t *testing.T) {
	manager := NewFontManager(&FontConfig{
		FontPaths:  []string{},
		NormalSize: 12.0,
		BoldSize:   14.0,
	})

	fontPath, err := manager.findSystemFont()
	
	// System font discovery may fail on systems without fonts installed
	if err != nil {
		t.Logf("System font discovery failed (expected on minimal systems): %v", err)
		assert.Contains(t, err.Error(), "no Vietnamese-compatible system fonts found", 
			"Error message should indicate no system fonts found")
		return
	}

	// If a font was found, verify it exists
	assert.NoError(t, err, "Should successfully find system font")
	assert.NotEmpty(t, fontPath, "Font path should not be empty")
	
	// Verify the file exists
	_, err = os.Stat(fontPath)
	assert.NoError(t, err, "System font file should exist at returned path")
	
	t.Logf("Found system font at: %s", fontPath)
}

// TestFindSystemFont_ByOS tests system font discovery for specific operating systems
func TestFindSystemFont_ByOS(t *testing.T) {
	manager := NewFontManager(&FontConfig{
		FontPaths:  []string{},
		NormalSize: 12.0,
		BoldSize:   14.0,
	})

	fontPath, err := manager.findSystemFont()
	
	switch runtime.GOOS {
	case "windows":
		if err == nil {
			assert.Contains(t, fontPath, "Windows\\Fonts", 
				"Windows font should be in Windows\\Fonts directory")
		}
	case "darwin":
		if err == nil {
			assert.True(t, 
				filepath.Dir(fontPath) == "/Library/Fonts" || 
				filepath.Dir(fontPath) == "/System/Library/Fonts/Supplemental",
				"macOS font should be in standard font directories")
		}
	case "linux":
		if err == nil {
			assert.Contains(t, fontPath, "/usr/share/fonts", 
				"Linux font should be in /usr/share/fonts directory")
		}
	}
	
	if err != nil {
		t.Logf("System font not found on %s (this is acceptable): %v", runtime.GOOS, err)
	} else {
		t.Logf("Found system font on %s: %s", runtime.GOOS, fontPath)
	}
}

// TestLoadFontFromPath_ValidPath tests loading a font from a valid path
func TestLoadFontFromPath_ValidPath(t *testing.T) {
	manager := NewFontManager(&FontConfig{
		FontPaths:  []string{},
		NormalSize: 12.0,
		BoldSize:   14.0,
	})

	// Try to find a system font for testing
	systemFontPath, err := manager.findSystemFont()
	if err != nil {
		t.Skip("No system fonts available for testing")
	}

	face, err := manager.loadFontFromPath(systemFontPath, 12.0)
	require.NoError(t, err, "Should successfully load font from valid path")
	assert.NotNil(t, face, "Font face should not be nil")
	
	// Clean up
	if face != nil {
		face.Close()
	}
}

// TestLoadFontFromPath_InvalidPath tests error handling for invalid paths
func TestLoadFontFromPath_InvalidPath(t *testing.T) {
	manager := NewFontManager(&FontConfig{
		FontPaths:  []string{},
		NormalSize: 12.0,
		BoldSize:   14.0,
	})

	face, err := manager.loadFontFromPath("/nonexistent/font.ttf", 12.0)
	assert.Error(t, err, "Should return error for nonexistent font file")
	assert.Nil(t, face, "Font face should be nil on error")
	assert.Contains(t, err.Error(), "font file not found", "Error should indicate file not found")
}

// TestLoadFontFromPath_DifferentSizes tests loading fonts with different sizes
func TestLoadFontFromPath_DifferentSizes(t *testing.T) {
	manager := NewFontManager(&FontConfig{
		FontPaths:  []string{},
		NormalSize: 12.0,
		BoldSize:   14.0,
	})

	// Try to find a system font for testing
	systemFontPath, err := manager.findSystemFont()
	if err != nil {
		t.Skip("No system fonts available for testing")
	}

	// Load with different sizes
	sizes := []float64{10.0, 12.0, 14.0, 16.0, 18.0}
	for _, size := range sizes {
		face, err := manager.loadFontFromPath(systemFontPath, size)
		require.NoError(t, err, "Should successfully load font with size %.1f", size)
		assert.NotNil(t, face, "Font face should not be nil for size %.1f", size)
		
		// Clean up
		if face != nil {
			face.Close()
		}
	}
}

// TestNewFontManager tests FontManager creation
func TestNewFontManager(t *testing.T) {
	config := &FontConfig{
		FontPaths:  []string{"/path/to/font1.ttf", "/path/to/font2.ttf"},
		NormalSize: 12.0,
		BoldSize:   14.0,
	}

	manager := NewFontManager(config)
	assert.NotNil(t, manager, "FontManager should not be nil")
	assert.Equal(t, config.FontPaths, manager.fontPaths, "Font paths should match")
	assert.Equal(t, config.NormalSize, manager.normalSize, "Normal size should match")
	assert.Equal(t, config.BoldSize, manager.boldSize, "Bold size should match")
}

// TestLoadFonts_FallbackToSystemFont tests fallback to system fonts when configured paths fail
func TestLoadFonts_FallbackToSystemFont(t *testing.T) {
	// Create manager with invalid configured paths
	manager := NewFontManager(&FontConfig{
		FontPaths:  []string{"/invalid/path1.ttf", "/invalid/path2.ttf"},
		NormalSize: 12.0,
		BoldSize:   14.0,
	})

	normalFont, boldFont, err := manager.LoadFonts()
	
	// Should either find a system font or return an error
	if err != nil {
		// No system fonts available
		assert.Error(t, err, "Should return error when no fonts available")
		assert.Contains(t, err.Error(), "failed to load fonts", "Error should indicate font loading failure")
	} else {
		// System font was found as fallback
		assert.NotNil(t, normalFont, "Normal font should not be nil when system font is used")
		assert.NotNil(t, boldFont, "Bold font should not be nil when system font is used")
		
		// Clean up
		if normalFont != nil {
			normalFont.Close()
		}
		if boldFont != nil {
			boldFont.Close()
		}
	}
}

// TestLoadFonts_EmptyFontPaths tests loading fonts with empty configured paths
func TestLoadFonts_EmptyFontPaths(t *testing.T) {
	manager := NewFontManager(&FontConfig{
		FontPaths:  []string{},
		NormalSize: 12.0,
		BoldSize:   14.0,
	})

	normalFont, boldFont, err := manager.LoadFonts()
	
	// Should try to find system fonts
	if err != nil {
		// No system fonts available
		assert.Error(t, err, "Should return error when no fonts available")
		assert.Contains(t, err.Error(), "failed to load fonts", "Error should indicate font loading failure")
	} else {
		// System font was found
		assert.NotNil(t, normalFont, "Normal font should not be nil when system font is used")
		assert.NotNil(t, boldFont, "Bold font should not be nil when system font is used")
		
		// Clean up
		if normalFont != nil {
			normalFont.Close()
		}
		if boldFont != nil {
			boldFont.Close()
		}
	}
}
