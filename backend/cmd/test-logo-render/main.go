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

const simpleHTML = `<!DOCTYPE html>
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
        .test-section {
            margin: 20px 0;
            border: 2px solid black;
            padding: 10px;
        }
        .logo-container {
            background: #f0f0f0;
            padding: 10px;
            margin: 10px 0;
        }
        .logo-container img {
            display: block;
            width: 200px;
            height: auto;
        }
    </style>
</head>
<body>
    <div class="test-section">
        <h2>Test 1: Text Only</h2>
        <p>This is plain text - should always work</p>
    </div>
    
    <div class="test-section">
        <h2>Test 2: Logo Image</h2>
        <div class="logo-container">
            {{if .LogoBase64}}
            <img src="{{.LogoBase64}}" alt="Logo">
            <p>Logo should appear above this text</p>
            {{else}}
            <p>NO LOGO DATA</p>
            {{end}}
        </div>
    </div>
    
    <div class="test-section">
        <h2>Test 3: Small Base64 Image</h2>
        <div class="logo-container">
            <!-- 1x1 red pixel -->
            <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8DwHwAFBQIAX8jx0gAAAABJRU5ErkJggg==" 
                 style="width:50px;height:50px;background:red;" alt="Red pixel">
            <p>Red square should appear above</p>
        </div>
    </div>
</body>
</html>`

type TestData struct {
	LogoBase64 string
}

func main() {
	fmt.Println("Testing logo rendering in chromedp...")
	
	// Load actual logo
	logoPath := "./uploads/logos/logo_24094.jpeg"
	logoBase64, err := loadImageAsBase64(logoPath)
	if err != nil {
		log.Printf("Warning: Failed to load logo: %v", err)
		logoBase64 = ""
	} else {
		fmt.Printf("✓ Logo loaded: %d bytes\n", len(logoBase64))
	}
	
	// Render HTML
	tmpl, err := template.New("test").Parse(simpleHTML)
	if err != nil {
		log.Fatal("Failed to parse template:", err)
	}
	
	data := TestData{LogoBase64: logoBase64}
	var htmlBuf bytes.Buffer
	if err := tmpl.Execute(&htmlBuf, data); err != nil {
		log.Fatal("Failed to execute template:", err)
	}
	
	htmlContent := htmlBuf.String()
	fmt.Printf("✓ HTML rendered: %d bytes\n", len(htmlContent))
	
	// Save HTML for manual inspection
	if err := os.WriteFile("test_logo_render.html", []byte(htmlContent), 0644); err != nil {
		log.Printf("Warning: Failed to save HTML: %v", err)
	} else {
		fmt.Println("✓ Saved test_logo_render.html (open in browser to verify)")
	}
	
	// Capture with chromedp
	fmt.Println("\nCapturing with chromedp...")
	if err := captureHTML(htmlContent, "test_logo_capture.png"); err != nil {
		log.Fatal("Capture failed:", err)
	}
	
	fmt.Println("\n✓ Test complete!")
	fmt.Println("  1. Open test_logo_render.html in browser - logo should show")
	fmt.Println("  2. Open test_logo_capture.png - logo should show")
	fmt.Println("  3. If logo shows in HTML but not PNG, chromedp has issue")
	fmt.Println("  4. If logo shows in both, check ESC/POS conversion")
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
	fmt.Printf("  Navigating to: %s\n", fileURL)
	
	var buf []byte
	err = chromedp.Run(ctx,
		chromedp.Navigate(fileURL),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("  Waiting for images to load...")
			return chromedp.Evaluate(`
				new Promise((resolve) => {
					const images = document.querySelectorAll('img');
					console.log('Found', images.length, 'images');
					if (images.length === 0) {
						resolve();
						return;
					}
					let loaded = 0;
					const checkComplete = () => {
						loaded++;
						console.log('Image loaded:', loaded, '/', images.length);
						if (loaded === images.length) resolve();
					};
					images.forEach((img, i) => {
						console.log('Image', i, ':', img.src.substring(0, 50));
						if (img.complete) {
							console.log('  Already complete');
							checkComplete();
						} else {
							img.onload = () => {
								console.log('  Loaded successfully');
								checkComplete();
							};
							img.onerror = (e) => {
								console.log('  Failed to load:', e);
								checkComplete();
							};
						}
					});
					setTimeout(() => {
						console.log('Timeout reached, proceeding anyway');
						resolve();
					}, 3000);
				});
			`, nil).Do(ctx)
		}),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("  Taking screenshot...")
			return nil
		}),
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

func loadImageAsBase64(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil && len(path) > 0 && path[0] == '/' {
		data, err = os.ReadFile("." + path)
	}
	if err != nil {
		return "", fmt.Errorf("failed to read image: %w", err)
	}

	mimeType := "image/jpeg"
	if len(data) > 4 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		mimeType = "image/png"
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}
