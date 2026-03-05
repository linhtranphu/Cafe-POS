package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/fogleman/gg"
)

const (
	ImageWidth = 576
	Threshold  = 128
)

func main() {
	fmt.Println("🧪 Testing ESC/POS Image Conversion")
	fmt.Println("====================================")

	// Create simple test image
	img := createTestImage()

	// Save original
	saveImage(img, "test_original.png")
	fmt.Println("✓ Saved test_original.png")

	// Binarize
	bwImg := binarizeImage(img, Threshold)
	saveImage(bwImg, "test_binarized.png")
	fmt.Println("✓ Saved test_binarized.png")

	// Convert to ESC/POS - Method 1 (current)
	escpos1 := imageToESCPOSMethod1(bwImg)
	os.WriteFile("test_escpos_method1.bin", escpos1, 0644)
	fmt.Printf("✓ Method 1: %d bytes\n", len(escpos1))

	// Convert to ESC/POS - Method 2 (alternative)
	escpos2 := imageToESCPOSMethod2(bwImg)
	os.WriteFile("test_escpos_method2.bin", escpos2, 0644)
	fmt.Printf("✓ Method 2: %d bytes\n", len(escpos2))

	// Convert to ESC/POS - Method 3 (bit image)
	escpos3 := imageToESCPOSMethod3(bwImg)
	os.WriteFile("test_escpos_method3.bin", escpos3, 0644)
	fmt.Printf("✓ Method 3: %d bytes\n", len(escpos3))

	fmt.Println("\n✅ Test files created!")
	fmt.Println("Send to printer to test:")
	fmt.Println("  cat test_escpos_method1.bin | nc 192.168.1.115 9100")
	fmt.Println("  cat test_escpos_method2.bin | nc 192.168.1.115 9100")
	fmt.Println("  cat test_escpos_method3.bin | nc 192.168.1.115 9100")
}

func createTestImage() image.Image {
	dc := gg.NewContext(ImageWidth, 400)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	// Draw test patterns
	dc.DrawRectangle(50, 50, 100, 100)
	dc.Fill()

	dc.DrawString("Test Text 123", 200, 100)

	// Draw lines
	for i := 0; i < 5; i++ {
		y := float64(200 + i*20)
		dc.DrawLine(50, y, 500, y)
		dc.Stroke()
	}

	return dc.Image()
}

func binarizeImage(img image.Image, threshold uint8) *image.Gray {
	bounds := img.Bounds()
	bw := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := uint8((r + g + b) / 3 / 257)

			if gray > threshold {
				bw.SetGray(x, y, color.Gray{Y: 255})
			} else {
				bw.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	return bw
}

// Method 1: GS v 0 (current implementation)
func imageToESCPOSMethod1(img image.Image) []byte {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	var buf bytes.Buffer

	// ESC @ - Initialize
	buf.Write([]byte{0x1B, 0x40})

	bytesPerLine := (width + 7) / 8

	for y := 0; y < height; y++ {
		// GS v 0 - Print raster bit image
		buf.Write([]byte{0x1D, 0x76, 0x30, 0x00})

		// Width in bytes (little endian)
		buf.WriteByte(byte(bytesPerLine & 0xFF))
		buf.WriteByte(byte((bytesPerLine >> 8) & 0xFF))

		// Height (1 line)
		buf.WriteByte(0x01)
		buf.WriteByte(0x00)

		// Image data
		lineData := make([]byte, bytesPerLine)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := (r + g + b) / 3

			if gray < 32768 {
				byteIndex := x / 8
				bitIndex := 7 - (x % 8)
				lineData[byteIndex] |= (1 << bitIndex)
			}
		}
		buf.Write(lineData)
	}

	// Feed and cut
	buf.Write([]byte{0x1B, 0x64, 0x03})
	buf.Write([]byte{0x1D, 0x56, 0x41, 0x00})

	return buf.Bytes()
}

// Method 2: ESC * (bit image mode)
func imageToESCPOSMethod2(img image.Image) []byte {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	var buf bytes.Buffer

	// ESC @ - Initialize
	buf.Write([]byte{0x1B, 0x40})

	// Process in 24-dot lines (3 bytes per column)
	for y := 0; y < height; y += 24 {
		// ESC * m nL nH - Select bit image mode
		// m = 33 (24-dot double-density)
		buf.Write([]byte{0x1B, 0x2A, 0x21})

		// Width in dots (little endian)
		buf.WriteByte(byte(width & 0xFF))
		buf.WriteByte(byte((width >> 8) & 0xFF))

		// Image data (3 bytes per column for 24 dots)
		for x := 0; x < width; x++ {
			for k := 0; k < 3; k++ {
				var b byte
				for bit := 0; bit < 8; bit++ {
					py := y + k*8 + bit
					if py < height {
						r, g, b2, _ := img.At(x, py).RGBA()
						gray := (r + g + b2) / 3
						if gray < 32768 {
							b |= (1 << uint(7-bit))
						}
					}
				}
				buf.WriteByte(b)
			}
		}

		// Line feed
		buf.Write([]byte{0x0A})
	}

	// Feed and cut
	buf.Write([]byte{0x1B, 0x64, 0x03})
	buf.Write([]byte{0x1D, 0x56, 0x41, 0x00})

	return buf.Bytes()
}

// Method 3: GS v 0 with different mode
func imageToESCPOSMethod3(img image.Image) []byte {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	var buf bytes.Buffer

	// ESC @ - Initialize
	buf.Write([]byte{0x1B, 0x40})

	// Process in chunks of 256 lines
	chunkHeight := 256
	for startY := 0; startY < height; startY += chunkHeight {
		endY := startY + chunkHeight
		if endY > height {
			endY = height
		}
		currentHeight := endY - startY

		bytesPerLine := (width + 7) / 8

		// GS v 0 - Print raster bit image
		buf.Write([]byte{0x1D, 0x76, 0x30, 0x00})

		// Width in bytes (little endian)
		buf.WriteByte(byte(bytesPerLine & 0xFF))
		buf.WriteByte(byte((bytesPerLine >> 8) & 0xFF))

		// Height in dots (little endian)
		buf.WriteByte(byte(currentHeight & 0xFF))
		buf.WriteByte(byte((currentHeight >> 8) & 0xFF))

		// Image data
		for y := startY; y < endY; y++ {
			lineData := make([]byte, bytesPerLine)
			for x := 0; x < width; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				gray := (r + g + b) / 3

				if gray < 32768 {
					byteIndex := x / 8
					bitIndex := 7 - (x % 8)
					lineData[byteIndex] |= (1 << bitIndex)
				}
			}
			buf.Write(lineData)
		}
	}

	// Feed and cut
	buf.Write([]byte{0x1B, 0x64, 0x03})
	buf.Write([]byte{0x1D, 0x56, 0x41, 0x00})

	return buf.Bytes()
}

func saveImage(img image.Image, filename string) {
	f, _ := os.Create(filename)
	defer f.Close()
	png.Encode(f, img)
}
