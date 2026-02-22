package printing

import (
	"fmt"
	"image"
	"image/color"
)

// ImageConverter converts bitmap images to ESC/POS raster format
type ImageConverter struct {
	pixelWidth int
}

// NewImageConverter creates a new ImageConverter instance
func NewImageConverter(pixelWidth int) *ImageConverter {
	return &ImageConverter{
		pixelWidth: pixelWidth,
	}
}

// ConvertToESCPOS converts a grayscale image to ESC/POS GS v 0 format
// GS v 0 format: GS v 0 m xL xH yL yH [d]k
// - GS v 0: Command prefix (0x1D 0x76 0x30)
// - m: Mode (0x00 for normal)
// - xL, xH: Width in bytes (little-endian)
// - yL, yH: Height in pixels (little-endian)
// - [d]k: Raster data (width_bytes * height bytes)
func (c *ImageConverter) ConvertToESCPOS(img *image.Gray) ([]byte, error) {
	if img == nil {
		return nil, fmt.Errorf("image conversion error: input image is nil")
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("image conversion error: invalid image dimensions (%dx%d pixels)", width, height)
	}

	if width > c.pixelWidth {
		return nil, fmt.Errorf("image conversion error: image width %d exceeds printer pixel width %d", width, c.pixelWidth)
	}

	// Apply threshold to create pure black/white image
	bwImage := c.applyThreshold(img)
	
	// Convert image to raster byte array
	rasterData := c.imageToRaster(bwImage)
	
	// Calculate dimensions
	widthBytes := (width + 7) / 8 // Round up to nearest byte
	
	// Validate raster data size
	expectedSize := widthBytes * height
	if len(rasterData) != expectedSize {
		return nil, fmt.Errorf("image conversion error: raster data size mismatch (expected %d bytes, got %d bytes)", expectedSize, len(rasterData))
	}

	// Build ESC/POS command
	// GS v 0 command prefix
	result := []byte{0x1D, 0x76, 0x30}
	
	// Mode: 0x00 for normal
	result = append(result, 0x00)
	
	// Width in bytes (little-endian)
	result = append(result, byte(widthBytes&0xFF))
	result = append(result, byte((widthBytes>>8)&0xFF))
	
	// Height in pixels (little-endian)
	result = append(result, byte(height&0xFF))
	result = append(result, byte((height>>8)&0xFF))
	
	// Append raster data
	result = append(result, rasterData...)
	
	return result, nil
}

// imageToRaster converts image pixels to raster byte array
// Each byte represents 8 horizontal pixels (1 = black, 0 = white)
func (c *ImageConverter) imageToRaster(img *image.Gray) []byte {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	widthBytes := (width + 7) / 8 // Round up to nearest byte
	
	raster := make([]byte, widthBytes*height)
	
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Get pixel value (0 = black, 255 = white in grayscale)
			pixel := img.GrayAt(x+bounds.Min.X, y+bounds.Min.Y).Y
			
			// If pixel is black (0), set corresponding bit to 1
			if pixel == 0 {
				byteIndex := y*widthBytes + x/8
				bitIndex := 7 - (x % 8) // MSB first
				raster[byteIndex] |= 1 << bitIndex
			}
		}
	}
	
	return raster
}

// applyThreshold applies threshold to create pure black/white image
// Pixels >= 128 → white (255)
// Pixels < 128 → black (0)
func (c *ImageConverter) applyThreshold(img *image.Gray) *image.Gray {
	bounds := img.Bounds()
	result := image.NewGray(bounds)
	
	const threshold = 128
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.GrayAt(x, y).Y
			
			if pixel >= threshold {
				// White
				result.SetGray(x, y, color.Gray{Y: 255})
			} else {
				// Black
				result.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}
	
	return result
}
