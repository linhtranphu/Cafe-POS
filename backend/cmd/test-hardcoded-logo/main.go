package main

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

// HTML with hardcoded small logo
const htmlWithLogo = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body {
            width: 576px;
            background: white;
            padding: 20px;
            font-family: Arial;
        }
        .header {
            display: flex;
            margin-bottom: 20px;
            border: 2px solid blue;
            padding: 10px;
        }
        .logo {
            width: 200px;
            margin-right: 20px;
            background: #f0f0f0;
        }
        .logo img {
            width: 200px;
            height: auto;
            display: block;
        }
        .info {
            flex: 1;
        }
    </style>
</head>
<body>
    <h1>Test Hardcoded Logo</h1>
    
    <div class="header">
        <div class="logo">
            <!-- Small 10x10 red PNG -->
            <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAAFUlEQVR42mP8z8BQz0AEYBxVSF+FABJADveWkH6oAAAAAElFTkSuQmCC" 
                 alt="Logo">
        </div>
        <div class="info">
            <h2>Shop Name</h2>
            <p>Address Line 1</p>
            <p>Phone: 123456789</p>
        </div>
    </div>
    
    <p>If you see a red square in the logo area, chromedp is working!</p>
    <p>If not, there's an issue with image rendering in chromedp.</p>
</body>
</html>`

func main() {
	fmt.Println("Testing hardcoded logo in chromedp...")
	
	// Save HTML for manual check
	if err := os.WriteFile("test_hardcoded.html", []byte(htmlWithLogo), 0644); err != nil {
		log.Printf("Warning: Failed to save HTML: %v", err)
	} else {
		fmt.Println("✓ Saved test_hardcoded.html")
	}
	
	// Capture with chromedp
	if err := captureHTML(htmlWithLogo, "test_hardcoded_capture.png"); err != nil {
		log.Fatal("Capture failed:", err)
	}
	
	fmt.Println("\n✓ Test complete!")
	fmt.Println("  1. Open test_hardcoded.html in browser")
	fmt.Println("  2. Open test_hardcoded_capture.png")
	fmt.Println("  3. Both should show a red square logo")
}

func captureHTML(html string, outputFile string) error {
	// Save to temp file
	tmpFile, err := os.CreateTemp("", "test_*.html")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)
	
	if _, err := tmpFile.WriteString(html); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()
	
	// Create Chrome context
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.WindowSize(576, 1200),
	)
	
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	
	// Navigate and capture
	fileURL := "file://" + tmpPath
	fmt.Printf("  Capturing from: %s\n", fileURL)
	
	var buf []byte
	err = chromedp.Run(ctx,
		chromedp.Navigate(fileURL),
		chromedp.WaitReady("body"),
		chromedp.Sleep(2*time.Second), // Wait for images
		chromedp.FullScreenshot(&buf, 100),
	)
	if err != nil {
		return fmt.Errorf("chromedp run failed: %w", err)
	}
	
	fmt.Printf("  Screenshot captured: %d bytes\n", len(buf))
	
	// Decode and save
	img, err := png.Decode(bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("png decode failed: %w", err)
	}
	
	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output: %w", err)
	}
	defer f.Close()
	
	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("png encode failed: %w", err)
	}
	
	fmt.Printf("✓ Saved %s (%dx%d)\n", outputFile, img.Bounds().Dx(), img.Bounds().Dy())
	return nil
}
