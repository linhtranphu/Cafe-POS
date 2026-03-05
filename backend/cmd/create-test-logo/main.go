package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func main() {
	// Tạo logo test 200x100 với text "CAFE KIRO"
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
	
	// Vẽ text "CAFE KIRO" đơn giản (pixel art)
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
	
	// Save
	file, err := os.Create("uploads/logo.png")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()
	
	if err := png.Encode(file, img); err != nil {
		log.Fatalf("Failed to encode PNG: %v", err)
	}
	
	log.Println("✅ Test logo created: uploads/logo.png")
}

func drawRect(img *image.RGBA, x, y, w, h int, c color.RGBA) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.Set(x+dx, y+dy, c)
		}
	}
}
