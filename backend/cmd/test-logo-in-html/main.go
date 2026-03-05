package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	log.Println("=== Test Logo in HTML with Chromedp ===")
	
	// Đọc logo và convert sang base64
	logoData, err := os.ReadFile("uploads/logo.png")
	if err != nil {
		log.Fatalf("Failed to read logo: %v", err)
	}
	
	logoBase64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(logoData)
	log.Printf("Logo base64 length: %d", len(logoBase64))
	log.Printf("Logo base64 preview: %s...", logoBase64[:100])
	
	// Tạo HTML đơn giản với logo
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body {
            width: 576px;
            background: white;
            padding: 20px;
            font-family: Arial, sans-serif;
        }
        .logo-container {
            border: 2px solid red;
            padding: 10px;
            margin-bottom: 20px;
        }
        .logo {
            width: 200px;
            border: 2px solid blue;
        }
        .info {
            font-size: 16px;
            margin: 10px 0;
        }
    </style>
</head>
<body>
    <div class="logo-container">
        <h2>Logo Test</h2>
        <img class="logo" src="%s" alt="Logo">
    </div>
    <div class="info">Nếu thấy logo "CAFE" ở trên = SUCCESS ✅</div>
    <div class="info">Nếu không thấy logo = FAILED ❌</div>
</body>
</html>`, logoBase64)
	
	log.Printf("HTML length: %d bytes", len(html))
	
	// Setup chromedp
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoSandbox,
		chromedp.DisableGPU,
		chromedp.WindowSize(800, 600),
	)
	
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	var buf []byte
	
	// Test 1: Data URL approach
	log.Println("\n=== Test 1: Data URL Approach ===")
	dataURL := "data:text/html;charset=utf-8," + url.PathEscape(html)
	log.Printf("Data URL length: %d", len(dataURL))
	
	err = chromedp.Run(ctx,
		chromedp.Navigate(dataURL),
		chromedp.WaitReady("body"),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.FullScreenshot(&buf, 100),
	)
	
	if err != nil {
		log.Printf("ERROR: %v", err)
	} else {
		log.Printf("Screenshot captured: %d bytes", len(buf))
		
		// Save
		if err := os.WriteFile("test_logo_dataurl.png", buf, 0644); err != nil {
			log.Printf("ERROR saving: %v", err)
		} else {
			log.Printf("✅ Saved: test_logo_dataurl.png")
		}
		
		// Decode to check
		img, err := png.Decode(bytes.NewReader(buf))
		if err != nil {
			log.Printf("ERROR decoding: %v", err)
		} else {
			log.Printf("Image size: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
		}
	}
	
	// Test 2: Temp file approach
	log.Println("\n=== Test 2: Temp File Approach ===")
	tmpFile, err := os.CreateTemp("", "test_logo_*.html")
	if err != nil {
		log.Fatalf("Failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)
	
	if _, err := tmpFile.WriteString(html); err != nil {
		log.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()
	
	fileURL := "file://" + tmpPath
	log.Printf("File URL: %s", fileURL)
	
	buf = nil
	err = chromedp.Run(ctx,
		chromedp.Navigate(fileURL),
		chromedp.WaitReady("body"),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.FullScreenshot(&buf, 100),
	)
	
	if err != nil {
		log.Printf("ERROR: %v", err)
	} else {
		log.Printf("Screenshot captured: %d bytes", len(buf))
		
		// Save
		if err := os.WriteFile("test_logo_file.png", buf, 0644); err != nil {
			log.Printf("ERROR saving: %v", err)
		} else {
			log.Printf("✅ Saved: test_logo_file.png")
		}
		
		// Decode to check
		img, err := png.Decode(bytes.NewReader(buf))
		if err != nil {
			log.Printf("ERROR decoding: %v", err)
		} else {
			log.Printf("Image size: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
		}
	}
	
	log.Println("\n=== Test completed ===")
	log.Println("Kiểm tra các file:")
	log.Println("- test_logo_dataurl.png (data URL approach)")
	log.Println("- test_logo_file.png (file URL approach)")
	log.Println("\nMở file để xem logo có hiển thị không")
}
