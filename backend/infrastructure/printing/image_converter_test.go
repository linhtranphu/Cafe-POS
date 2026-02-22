package printing

import (
	"image"
	"image/color"
	"testing"
)

func TestNewImageConverter(t *testing.T) {
	converter := NewImageConverter(384)
	
	if converter == nil {
		t.Fatal("NewImageConverter returned nil")
	}
	
	if converter.pixelWidth != 384 {
		t.Errorf("Expected pixelWidth 384, got %d", converter.pixelWidth)
	}
}

func TestApplyThreshold(t *testing.T) {
	converter := NewImageConverter(384)
	
	// Create a test image with various gray values
	img := image.NewGray(image.Rect(0, 0, 4, 1))
	img.SetGray(0, 0, color.Gray{Y: 0})    // Black
	img.SetGray(1, 0, color.Gray{Y: 64})   // Dark gray (< 128)
	img.SetGray(2, 0, color.Gray{Y: 128})  // Mid gray (>= 128)
	img.SetGray(3, 0, color.Gray{Y: 255})  // White
	
	result := converter.applyThreshold(img)
	
	// Check that threshold was applied correctly
	if result.GrayAt(0, 0).Y != 0 {
		t.Errorf("Expected pixel 0 to be black (0), got %d", result.GrayAt(0, 0).Y)
	}
	if result.GrayAt(1, 0).Y != 0 {
		t.Errorf("Expected pixel 1 to be black (0), got %d", result.GrayAt(1, 0).Y)
	}
	if result.GrayAt(2, 0).Y != 255 {
		t.Errorf("Expected pixel 2 to be white (255), got %d", result.GrayAt(2, 0).Y)
	}
	if result.GrayAt(3, 0).Y != 255 {
		t.Errorf("Expected pixel 3 to be white (255), got %d", result.GrayAt(3, 0).Y)
	}
}

func TestImageToRaster(t *testing.T) {
	converter := NewImageConverter(384)
	
	// Create a simple 8x2 black and white image
	// Row 0: alternating black/white (10101010)
	// Row 1: all black (11111111)
	img := image.NewGray(image.Rect(0, 0, 8, 2))
	
	// Row 0: alternating pattern
	for x := 0; x < 8; x++ {
		if x%2 == 0 {
			img.SetGray(x, 0, color.Gray{Y: 0}) // Black
		} else {
			img.SetGray(x, 0, color.Gray{Y: 255}) // White
		}
	}
	
	// Row 1: all black
	for x := 0; x < 8; x++ {
		img.SetGray(x, 1, color.Gray{Y: 0})
	}
	
	raster := converter.imageToRaster(img)
	
	// Expected: 2 bytes (1 byte per row for 8 pixels)
	if len(raster) != 2 {
		t.Fatalf("Expected 2 bytes, got %d", len(raster))
	}
	
	// Row 0: 10101010 in binary = 0xAA
	if raster[0] != 0xAA {
		t.Errorf("Expected first byte to be 0xAA (alternating pattern), got 0x%02X", raster[0])
	}
	
	// Row 1: 11111111 in binary = 0xFF
	if raster[1] != 0xFF {
		t.Errorf("Expected second byte to be 0xFF (all black), got 0x%02X", raster[1])
	}
}

func TestConvertToESCPOS(t *testing.T) {
	converter := NewImageConverter(384)
	
	// Create a simple 8x2 image
	img := image.NewGray(image.Rect(0, 0, 8, 2))
	
	// Fill with a pattern
	for y := 0; y < 2; y++ {
		for x := 0; x < 8; x++ {
			img.SetGray(x, y, color.Gray{Y: 0}) // All black
		}
	}
	
	result, err := converter.ConvertToESCPOS(img)
	if err != nil {
		t.Fatalf("ConvertToESCPOS failed: %v", err)
	}
	
	// Check command prefix: GS v 0 (0x1D 0x76 0x30)
	if len(result) < 8 {
		t.Fatalf("Result too short: %d bytes", len(result))
	}
	
	if result[0] != 0x1D || result[1] != 0x76 || result[2] != 0x30 {
		t.Errorf("Expected GS v 0 command prefix (0x1D 0x76 0x30), got 0x%02X 0x%02X 0x%02X",
			result[0], result[1], result[2])
	}
	
	// Check mode: 0x00
	if result[3] != 0x00 {
		t.Errorf("Expected mode 0x00, got 0x%02X", result[3])
	}
	
	// Check width in bytes (8 pixels = 1 byte)
	widthBytes := int(result[4]) | (int(result[5]) << 8)
	if widthBytes != 1 {
		t.Errorf("Expected width 1 byte, got %d", widthBytes)
	}
	
	// Check height (2 pixels)
	height := int(result[6]) | (int(result[7]) << 8)
	if height != 2 {
		t.Errorf("Expected height 2 pixels, got %d", height)
	}
	
	// Check raster data length (1 byte * 2 rows = 2 bytes)
	rasterDataLen := len(result) - 8
	if rasterDataLen != 2 {
		t.Errorf("Expected 2 bytes of raster data, got %d", rasterDataLen)
	}
}

func TestConvertToESCPOS_WithThreshold(t *testing.T) {
	converter := NewImageConverter(384)
	
	// Create an image with gray values that need thresholding
	img := image.NewGray(image.Rect(0, 0, 8, 1))
	
	// Set pixels with various gray values
	for x := 0; x < 4; x++ {
		img.SetGray(x, 0, color.Gray{Y: 50}) // Dark gray -> black
	}
	for x := 4; x < 8; x++ {
		img.SetGray(x, 0, color.Gray{Y: 200}) // Light gray -> white
	}
	
	result, err := converter.ConvertToESCPOS(img)
	if err != nil {
		t.Fatalf("ConvertToESCPOS failed: %v", err)
	}
	
	// Extract raster data (skip 8-byte header)
	rasterData := result[8:]
	
	// Expected: first 4 pixels black (1111), last 4 pixels white (0000)
	// Binary: 11110000 = 0xF0
	if len(rasterData) != 1 {
		t.Fatalf("Expected 1 byte of raster data, got %d", len(rasterData))
	}
	
	if rasterData[0] != 0xF0 {
		t.Errorf("Expected raster byte 0xF0 (half black, half white), got 0x%02X", rasterData[0])
	}
}

func TestImageToRaster_WidthNotMultipleOf8(t *testing.T) {
	converter := NewImageConverter(384)
	
	// Create an image with width not multiple of 8 (10 pixels)
	img := image.NewGray(image.Rect(0, 0, 10, 1))
	
	// Set all pixels to black
	for x := 0; x < 10; x++ {
		img.SetGray(x, 0, color.Gray{Y: 0})
	}
	
	raster := converter.imageToRaster(img)
	
	// Expected: 2 bytes (10 pixels needs 2 bytes, with 6 bits padding)
	if len(raster) != 2 {
		t.Fatalf("Expected 2 bytes for 10 pixels, got %d", len(raster))
	}
	
	// First byte: 11111111 = 0xFF (8 black pixels)
	if raster[0] != 0xFF {
		t.Errorf("Expected first byte 0xFF, got 0x%02X", raster[0])
	}
	
	// Second byte: 11000000 = 0xC0 (2 black pixels + 6 padding bits)
	if raster[1] != 0xC0 {
		t.Errorf("Expected second byte 0xC0 (2 pixels + padding), got 0x%02X", raster[1])
	}
}
