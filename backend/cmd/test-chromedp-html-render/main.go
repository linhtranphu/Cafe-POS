package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

const testHTML = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            width: 576px;
            background: white;
            color: black;
            font-family: Arial, sans-serif;
            padding: 20px;
        }
        
        .header {
            display: flex;
            align-items: flex-start;
            margin-bottom: 20px;
        }
        
        .logo {
            width: 200px;
            margin-right: 20px;
        }
        
        .logo img {
            width: 200px;
            height: auto;
        }
        
        .shop-info {
            flex: 1;
        }
        
        .shop-name {
            font-size: 24px;
            font-weight: bold;
            margin-bottom: 10px;
        }
        
        .content {
            margin-top: 20px;
        }
        
        .test-text {
            font-size: 18px;
            margin: 10px 0;
        }
    </style>
</head>
<body>
    <div class="header">
        {{if .LogoBase64}}
        <div class="logo">
            <img src="{{.LogoBase64}}" alt="Logo">
        </div>
        {{end}}
        <div class="shop-info">
            <div class="shop-name">{{.ShopName}}</div>
            <div>{{.ShopAddress}}</div>
        </div>
    </div>
    
    <div class="content">
        <div class="test-text">This is a test render</div>
        <div class="test-text">Testing chromedp HTML capture</div>
        <div class="test-text">With base64 logo image</div>
    </div>
</body>
</html>`

type TestData struct {
	ShopName    string
	ShopAddress string
	LogoBase64  string
}

func main() {
	fmt.Println("Testing chromedp HTML rendering with base64 images...")
	
	// Load logo as base64
	logoPath := "./uploads/logos/logo_24094.jpeg"
	logoBase64, err := loadImageAsBase64(logoPath)
	if err != nil {
		log.Printf("Warning: Failed to load logo: %v", err)
		logoBase64 = ""
	} else {
		fmt.Printf("✓ Logo loaded: %d bytes\n", len(logoBase64))
	}
	
	// Prepare data
	data := TestData{
		ShopName:    "Test Coffee Shop",
		ShopAddress: "123 Test Street",
		LogoBase64:  logoBase64,
	}
	
	// Render HTML
	tmpl, err := template.New("test").Parse(testHTML)
	if err != nil {
		log.Fatal("Failed to parse template:", err)
	}
	
	var htmlBuf bytes.Buffer
	if err := tmpl.Execute(&htmlBuf, data); err != nil {
		log.Fatal("Failed to execute template:", err)
	}
	
	htmlContent := htmlBuf.String()
	fmt.Printf("✓ HTML rendered: %d bytes\n", len(htmlContent))
	
	// Save HTML for inspection
	if err := os.WriteFile("test_render.html", []byte(htmlContent), 0644); err != nil {
		log.Printf("Warning: Failed to save HTML: %v", err)
	} else {
		fmt.Println("✓ Saved test_render.html")
	}
	
	// Test Method 1: Using temp file with file:// URL
	fmt.Println("\n=== Method 1: file:// URL ===")
	if err := testFileURL(htmlContent); err != nil {
		log.Printf("Method 1 failed: %v", err)
	}
	
	// Test Method 2: Using data URL
	fmt.Println("\n=== Method 2: data URL ===")
	if err := testDataURL(htmlContent); err != nil {
		log.Printf("Method 2 failed: %v", err)
	}
	
	// Test Method 3: Using SetDocumentContent
	fmt.Println("\n=== Method 3: SetDocumentContent ===")
	if err := testSetDocumentContent(htmlContent); err != nil {
		log.Printf("Method 3 failed: %v", err)
	}
	
	fmt.Println("\n✓ All tests completed. Check output PNG files.")
}

func testFileURL(html string) error {
	// Create temp file
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
	var buf []byte
	
	err = chromedp.Run(ctx,
		chromedp.Navigate(fileURL),
		chromedp.Sleep(2*time.Second),
		chromedp.FullScreenshot(&buf, 100),
	)
	if err != nil {
		return fmt.Errorf("chromedp run failed: %w", err)
	}
	
	// Decode and save
	img, err := png.Decode(bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("png decode failed: %w", err)
	}
	
	f, err := os.Create("test_method1_file_url.png")
	if err != nil {
		return fmt.Errorf("failed to create output: %w", err)
	}
	defer f.Close()
	
	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("png encode failed: %w", err)
	}
	
	fmt.Printf("✓ Saved test_method1_file_url.png (%dx%d)\n", img.Bounds().Dx(), img.Bounds().Dy())
	return nil
}

func testDataURL(html string) error {
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
	
	// Use data URL
	dataURL := "data:text/html;charset=utf-8," + html
	var buf []byte
	
	err := chromedp.Run(ctx,
		chromedp.Navigate(dataURL),
		chromedp.Sleep(2*time.Second),
		chromedp.FullScreenshot(&buf, 100),
	)
	if err != nil {
		return fmt.Errorf("chromedp run failed: %w", err)
	}
	
	// Decode and save
	img, err := png.Decode(bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("png decode failed: %w", err)
	}
	
	f, err := os.Create("test_method2_data_url.png")
	if err != nil {
		return fmt.Errorf("failed to create output: %w", err)
	}
	defer f.Close()
	
	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("png encode failed: %w", err)
	}
	
	fmt.Printf("✓ Saved test_method2_data_url.png (%dx%d)\n", img.Bounds().Dx(), img.Bounds().Dy())
	return nil
}

func testSetDocumentContent(html string) error {
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
	
	// Use ActionFunc with page.SetDocumentContent
	var buf []byte
	
	err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Use chromedp's built-in method
			return chromedp.Run(ctx, chromedp.Evaluate(`
				document.open();
				document.write(`+"`"+html+"`"+`);
				document.close();
			`, nil))
		}),
		chromedp.Sleep(2*time.Second),
		chromedp.FullScreenshot(&buf, 100),
	)
	if err != nil {
		return fmt.Errorf("chromedp run failed: %w", err)
	}
	
	// Decode and save
	img, err := png.Decode(bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("png decode failed: %w", err)
	}
	
	f, err := os.Create("test_method3_set_content.png")
	if err != nil {
		return fmt.Errorf("failed to create output: %w", err)
	}
	defer f.Close()
	
	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("png encode failed: %w", err)
	}
	
	fmt.Printf("✓ Saved test_method3_set_content.png (%dx%d)\n", img.Bounds().Dx(), img.Bounds().Dy())
	return nil
}

func loadImageAsBase64(path string) (string, error) {
	// Try original path first
	data, err := os.ReadFile(path)
	
	// If failed and path starts with /, try prepending "."
	if err != nil && len(path) > 0 && path[0] == '/' {
		data, err = os.ReadFile("." + path)
	}
	
	if err != nil {
		return "", fmt.Errorf("failed to read image from %s: %w", path, err)
	}

	// Detect MIME type
	mimeType := "image/jpeg"
	if len(data) > 4 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		mimeType = "image/png"
	}

	// Use standard library base64 encoding
	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}
