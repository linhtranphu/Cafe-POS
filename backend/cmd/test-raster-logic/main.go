package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/kenshaw/escpos"
	"github.com/kenshaw/escpos/raster"
)

func main() {
	fmt.Println("Testing raster logic...")
	
	// Create a simple test image: 8x8 pixels
	// Top half: WHITE (255, 255, 255)
	// Bottom half: BLACK (0, 0, 0)
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	
	// Top 4 rows: WHITE
	for y := 0; y < 4; y++ {
		for x := 0; x < 8; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}
	
	// Bottom 4 rows: BLACK
	for y := 4; y < 8; y++ {
		for x := 0; x < 8; x++ {
			img.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
		}
	}
	
	// Save test image
	f, _ := os.Create("test_image.png")
	png.Encode(f, img)
	f.Close()
	fmt.Println("Created test_image.png (top=white, bottom=black)")
	
	// Test with different thresholds
	thresholds := []float64{0.1, 0.3, 0.5, 0.7, 0.9}
	
	for _, threshold := range thresholds {
		fmt.Printf("\n=== Testing threshold %.1f ===\n", threshold)
		
		converter := &raster.Converter{
			MaxWidth:  8,
			Threshold: threshold,
		}
		
		data, width, bytesWidth := converter.ToRaster(img)
		
		fmt.Printf("Width: %d, BytesWidth: %d, DataLen: %d\n", width, bytesWidth, len(data))
		
		// Print bitmap
		for y := 0; y < 8; y++ {
			fmt.Printf("Row %d: ", y)
			for x := 0; x < 8; x++ {
				byteIndex := y*bytesWidth + x/8
				bitIndex := 7 - (x % 8)
				bit := (data[byteIndex] >> bitIndex) & 1
				if bit == 1 {
					fmt.Print("█") // Black
				} else {
					fmt.Print("░") // White
				}
			}
			fmt.Println()
		}
	}
	
	// Test lightness calculation
	fmt.Println("\n=== Lightness values ===")
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	gray := color.RGBA{R: 128, G: 128, B: 128, A: 255}
	
	fmt.Printf("WHITE (255,255,255): lightness = %.3f\n", calculateLightness(white))
	fmt.Printf("BLACK (0,0,0): lightness = %.3f\n", calculateLightness(black))
	fmt.Printf("GRAY (128,128,128): lightness = %.3f\n", calculateLightness(gray))
	
	// Test ESC/POS output
	fmt.Println("\n=== ESC/POS output test ===")
	var buf bytes.Buffer
	e := escpos.New(&buf)
	e.Init()
	
	converter := &raster.Converter{
		MaxWidth:  8,
		Threshold: 0.5,
	}
	converter.Print(img, e)
	
	e.FormfeedN(2)
	e.Cut()
	e.End()
	
	fmt.Printf("ESC/POS data length: %d bytes\n", buf.Len())
	fmt.Printf("First 50 bytes: %v\n", buf.Bytes()[:min(50, buf.Len())])
}

func calculateLightness(c color.Color) float64 {
	const (
		lumR, lumG, lumB = 55, 182, 18
	)
	r, g, b, _ := c.RGBA()
	return float64(lumR*r+lumG*g+lumB*b) / float64(0xffff*(lumR+lumG+lumB))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
