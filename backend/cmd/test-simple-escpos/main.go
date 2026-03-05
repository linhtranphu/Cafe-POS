package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
)

func main() {
	log.Println("=== Test Simple ESC/POS Commands ===")
	
	// Create a simple test image: 576x100 with text "TEST"
	width := 576
	height := 100
	img := image.NewGray(image.Rect(0, 0, width, height))
	
	// Fill with white
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetGray(x, y, color.Gray{255})
		}
	}
	
	// Draw black rectangle (simple test pattern)
	for y := 20; y < 80; y++ {
		for x := 100; x < 476; x++ {
			img.SetGray(x, y, color.Gray{0})
		}
	}
	
	log.Println("Created test image: 576x100 with black rectangle")
	
	// Convert to ESC/POS using GS v 0 (raster bit image)
	escposData := imageToESCPOSRaster(img)
	
	log.Printf("ESC/POS data size: %d bytes", len(escposData))
	
	// Show first 50 bytes
	log.Println("\nFirst 50 bytes (hex):")
	for i := 0; i < 50 && i < len(escposData); i++ {
		if i%16 == 0 {
			fmt.Printf("\n%04x: ", i)
		}
		fmt.Printf("%02x ", escposData[i])
	}
	fmt.Println()
	
	// Save to file
	filename := "test_simple_escpos.bin"
	if err := os.WriteFile(filename, escposData, 0644); err != nil {
		log.Fatalf("Failed to save: %v", err)
	}
	
	log.Printf("\n✅ Saved to %s", filename)
	log.Println("\nThis file uses GS v 0 (raster bit image) command")
	log.Println("which is more widely supported than GS 8 L")
}

func imageToESCPOSRaster(img image.Image) []byte {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	var buf bytes.Buffer
	
	// ESC @ - Initialize printer
	buf.Write([]byte{0x1B, 0x40})
	
	// Process image line by line
	bytesPerLine := (width + 7) / 8
	
	for y := 0; y < height; y++ {
		// GS v 0 - Print raster bit image
		// Format: GS v 0 m xL xH yL yH d1...dk
		buf.Write([]byte{0x1D, 0x76, 0x30, 0x00})
		
		// xL xH - width in bytes (little endian)
		buf.WriteByte(byte(bytesPerLine & 0xFF))
		buf.WriteByte(byte((bytesPerLine >> 8) & 0xFF))
		
		// yL yH - height (1 line at a time)
		buf.WriteByte(0x01)
		buf.WriteByte(0x00)
		
		// Convert line to bitmap
		lineData := make([]byte, bytesPerLine)
		for x := 0; x < width; x++ {
			gray := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			// If pixel is dark (< 128), set bit to 1 (print black)
			if gray.Y < 128 {
				byteIndex := x / 8
				bitIndex := 7 - (x % 8)
				lineData[byteIndex] |= (1 << bitIndex)
			}
		}
		buf.Write(lineData)
	}
	
	// Feed paper
	buf.Write([]byte{0x1B, 0x64, 0x03}) // ESC d 3 (feed 3 lines)
	
	// Cut paper
	buf.Write([]byte{0x1D, 0x56, 0x00}) // GS V 0 (full cut)
	
	return buf.Bytes()
}
