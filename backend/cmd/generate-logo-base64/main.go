package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"bytes"
)

func main() {
	// Tạo logo đơn giản 200x100
	width := 200
	height := 100
	
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Background trắng
	white := color.RGBA{255, 255, 255, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, white)
		}
	}
	
	// Vẽ border đen
	black := color.RGBA{0, 0, 0, 255}
	for x := 0; x < width; x++ {
		img.Set(x, 0, black)
		img.Set(x, height-1, black)
	}
	for y := 0; y < height; y++ {
		img.Set(0, y, black)
		img.Set(width-1, y, black)
	}
	
	// Vẽ text "CAFE" đơn giản
	// C
	drawRect(img, 20, 30, 5, 40, black)
	drawRect(img, 20, 30, 20, 5, black)
	drawRect(img, 20, 65, 20, 5, black)
	
	// A
	drawRect(img, 50, 30, 5, 40, black)
	drawRect(img, 70, 30, 5, 40, black)
	drawRect(img, 50, 30, 25, 5, black)
	drawRect(img, 50, 50, 25, 5, black)
	
	// F
	drawRect(img, 85, 30, 5, 40, black)
	drawRect(img, 85, 30, 20, 5, black)
	drawRect(img, 85, 50, 15, 5, black)
	
	// E
	drawRect(img, 115, 30, 5, 40, black)
	drawRect(img, 115, 30, 20, 5, black)
	drawRect(img, 115, 50, 15, 5, black)
	drawRect(img, 115, 65, 20, 5, black)
	
	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}
	
	// Convert to base64
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	dataURL := fmt.Sprintf("data:image/png;base64,%s", encoded)
	
	fmt.Println("Logo Base64 Data URL:")
	fmt.Println(dataURL)
	fmt.Printf("\nLength: %d bytes\n", len(dataURL))
}

func drawRect(img *image.RGBA, x, y, w, h int, c color.RGBA) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.Set(x+dx, y+dy, c)
		}
	}
}
