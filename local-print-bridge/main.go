package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"service":   "Local Print Bridge (Go + chromedp)",
			"version":   "1.0.0",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Status endpoint
	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success":   true,
			"uptime":    time.Since(startTime).Seconds(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Test printer connection
	r.POST("/test-connection", handleTestConnection)

	// Render HTML and print
	r.POST("/render-and-print", handleRenderAndPrint)

	// Test render (for debugging)
	r.POST("/test-render", handleTestRender)

	// Direct print (ESC/POS data)
	r.POST("/print", handleDirectPrint)

	log.Println(strings.Repeat("=", 50))
	log.Println("🖨️  Local Print Bridge Started (Go + chromedp)")
	log.Println(strings.Repeat("=", 50))
	log.Printf("Server running on: http://%s:%s\n", host, port)
	log.Printf("Default Bill Printer: %s\n", os.Getenv("DEFAULT_BILL_PRINTER_IP"))
	log.Printf("Default Label Printer: %s\n", os.Getenv("DEFAULT_LABEL_PRINTER_IP"))
	log.Println(strings.Repeat("=", 50))
	log.Println("Endpoints:")
	log.Println("  POST /render-and-print - Render HTML and print")
	log.Println("  POST /print - Direct print ESC/POS data")
	log.Println("  POST /test-connection - Test printer connection")
	log.Println("  GET /health - Health check")
	log.Println("  GET /status - Service status")
	log.Println(strings.Repeat("=", 50))
	log.Println("Ready to accept requests!")
	log.Println(strings.Repeat("=", 50))

	addr := host + ":" + port
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

var startTime = time.Now()

// handleTestConnection tests printer connection
func handleTestConnection(c *gin.Context) {
	var req struct {
		PrinterIP   string `json:"printerIP" binding:"required"`
		PrinterPort int    `json:"printerPort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"success": false, "error": err.Error()})
		return
	}

	if req.PrinterPort == 0 {
		req.PrinterPort = 9100
	}

	log.Printf("[Test Connection] Testing: %s:%d", req.PrinterIP, req.PrinterPort)

	// Try to connect
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", req.PrinterIP, req.PrinterPort), 5*time.Second)
	if err != nil {
		log.Printf("[Test Connection] Failed: %v", err)
		c.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
			"printer": fmt.Sprintf("%s:%d", req.PrinterIP, req.PrinterPort),
		})
		return
	}
	defer conn.Close()

	log.Printf("[Test Connection] Success")
	c.JSON(200, gin.H{
		"success": true,
		"message": "Printer connection successful",
		"printer": fmt.Sprintf("%s:%d", req.PrinterIP, req.PrinterPort),
	})
}

// handleRenderAndPrint renders HTML and prints
func handleRenderAndPrint(c *gin.Context) {
	var req struct {
		HTML        string `json:"html" binding:"required"`
		Width       int    `json:"width"`
		PrinterIP   string `json:"printerIP" binding:"required"`
		PrinterPort int    `json:"printerPort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"success": false, "error": err.Error()})
		return
	}

	if req.Width == 0 {
		req.Width = 576
	}
	if req.PrinterPort == 0 {
		req.PrinterPort = 9100
	}

	log.Printf("[Render & Print] Request - Printer: %s:%d", req.PrinterIP, req.PrinterPort)

	// Render HTML to ESC/POS
	escposData, err := renderHTMLToESCPOS(req.HTML, req.Width)
	if err != nil {
		log.Printf("[Render & Print] Render failed: %v", err)
		c.JSON(500, gin.H{"success": false, "error": fmt.Sprintf("Render failed: %v", err)})
		return
	}

	log.Printf("[Render & Print] Rendered - ESC/POS: %d bytes", len(escposData))

	// Send to printer
	if err := sendToPrinter(req.PrinterIP, req.PrinterPort, escposData); err != nil {
		log.Printf("[Render & Print] Print failed: %v", err)
		c.JSON(500, gin.H{"success": false, "error": fmt.Sprintf("Print failed: %v", err)})
		return
	}

	log.Printf("[Render & Print] Print successful")
	c.JSON(200, gin.H{
		"success":   true,
		"message":   "Render and print completed successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// handleTestRender tests HTML rendering
func handleTestRender(c *gin.Context) {
	testHTML := `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body {
      width: 576px;
      margin: 0;
      padding: 20px;
      font-family: Arial, sans-serif;
    }
    h1 { text-align: center; font-size: 24px; }
    p { font-size: 14px; }
  </style>
</head>
<body>
  <h1>Test Bill</h1>
  <p>Order #12345</p>
  <p>Date: ` + time.Now().Format("02/01/2006 03:04 PM") + `</p>
  <p>Total: 100,000 VND</p>
  <p>Thank you!</p>
</body>
</html>
`

	escposData, err := renderHTMLToESCPOS(testHTML, 576)
	if err != nil {
		c.JSON(500, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Test render successful",
		"size":    len(escposData),
	})
}

// handleDirectPrint prints ESC/POS data directly
func handleDirectPrint(c *gin.Context) {
	var req struct {
		Content     string `json:"content" binding:"required"`
		PrinterIP   string `json:"printerIP" binding:"required"`
		PrinterPort int    `json:"printerPort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"success": false, "error": err.Error()})
		return
	}

	if req.PrinterPort == 0 {
		req.PrinterPort = 9100
	}

	log.Printf("[Direct Print] Request - Printer: %s:%d", req.PrinterIP, req.PrinterPort)

	// Decode base64 if needed
	data := []byte(req.Content)

	if err := sendToPrinter(req.PrinterIP, req.PrinterPort, data); err != nil {
		log.Printf("[Direct Print] Failed: %v", err)
		c.JSON(500, gin.H{"success": false, "error": err.Error()})
		return
	}

	log.Printf("[Direct Print] Success")
	c.JSON(200, gin.H{
		"success":   true,
		"message":   "Print completed successfully",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// renderHTMLToESCPOS renders HTML to ESC/POS using chromedp
func renderHTMLToESCPOS(html string, width int) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var buf []byte
	var contentHeight int64
	
	if err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.EmulateViewport(int64(width), 2000),
		chromedp.WaitReady("body"),
		chromedp.Evaluate(`document.body.scrollHeight`, &contentHeight),
	); err != nil {
		return nil, fmt.Errorf("chromedp failed: %w", err)
	}

	// Add padding to avoid cutting off bottom content (margins, padding, etc.)
	contentHeight += 20

	log.Printf("[Render] Content height (with padding): %dpx", contentHeight)

	// Capture screenshot with exact dimensions using CDP
	if err := chromedp.Run(ctx,
		chromedp.EmulateViewport(int64(width), contentHeight),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Use page.CaptureScreenshot with clip to get exact dimensions
			screenshot, err := page.CaptureScreenshot().
				WithClip(&page.Viewport{
					X:      0,
					Y:      0,
					Width:  float64(width),
					Height: float64(contentHeight),
					Scale:  1,
				}).
				Do(ctx)
			if err != nil {
				return err
			}
			buf = screenshot
			return nil
		}),
	); err != nil {
		return nil, fmt.Errorf("chromedp screenshot failed: %w", err)
	}

	// Convert screenshot to ESC/POS
	escposData, err := imageToESCPOS(buf, width)
	if err != nil {
		return nil, fmt.Errorf("image conversion failed: %w", err)
	}

	return escposData, nil
}

// imageToESCPOS converts image to ESC/POS raster commands
func imageToESCPOS(imageData []byte, width int) ([]byte, error) {
	// Decode PNG image
	img, err := png.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG: %w", err)
	}

	// Convert to grayscale
	grayImg := imaging.Grayscale(img)

	// Get image dimensions
	bounds := grayImg.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	log.Printf("[Image Processing] Original: %dx%d", imgWidth, imgHeight)

	// Resize if needed to match printer width
	if imgWidth != width {
		grayImg = imaging.Resize(grayImg, width, 0, imaging.Lanczos)
		bounds = grayImg.Bounds()
		imgWidth = bounds.Dx()
		imgHeight = bounds.Dy()
		log.Printf("[Image Processing] Resized: %dx%d", imgWidth, imgHeight)
	}

	// Apply Floyd-Steinberg dithering
	ditheredImg := floydSteinbergDithering(grayImg)

	// Convert to ESC/POS raster commands
	commands := []byte{}

	// ESC @ - Initialize printer
	commands = append(commands, 0x1B, 0x40)

	// Process image line by line
	bytesPerLine := (imgWidth + 7) / 8 // Round up to nearest byte

	for y := 0; y < imgHeight; y++ {
		// GS v 0 - Print raster bit image
		commands = append(commands, 0x1D, 0x76, 0x30, 0x00)

		// xL xH - width in bytes (little endian)
		commands = append(commands, byte(bytesPerLine&0xFF), byte((bytesPerLine>>8)&0xFF))

		// yL yH - height (1 line at a time)
		commands = append(commands, 0x01, 0x00)

		// Convert line to bitmap
		lineData := make([]byte, bytesPerLine)

		for x := 0; x < imgWidth; x++ {
			// Get pixel value (0 = black, 255 = white after dithering)
			r, _, _, _ := ditheredImg.At(x, y).RGBA()
			isBlack := r < 32768 // Threshold at 50%

			if isBlack {
				byteIndex := x / 8
				bitIndex := 7 - (x % 8)
				lineData[byteIndex] |= (1 << bitIndex)
			}
		}

		commands = append(commands, lineData...)
	}

	// Feed paper (3 lines)
	commands = append(commands, 0x1B, 0x64, 0x03)

	// Cut paper (full cut)
	commands = append(commands, 0x1D, 0x56, 0x00)

	log.Printf("[Image Processing] ESC/POS generated: %d bytes", len(commands))

	return commands, nil
}

// floydSteinbergDithering applies Floyd-Steinberg dithering algorithm
func floydSteinbergDithering(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create a new grayscale image
	dithered := image.NewGray(bounds)

	// Copy original image to dithered
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dithered.Set(x, y, img.At(x, y))
		}
	}

	// Apply Floyd-Steinberg dithering
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldPixel := dithered.GrayAt(x, y).Y
			var newPixel uint8
			if oldPixel > 128 {
				newPixel = 255
			} else {
				newPixel = 0
			}

			dithered.SetGray(x, y, color.Gray{Y: newPixel})

			quantError := int(oldPixel) - int(newPixel)

			// Distribute error to neighboring pixels
			if x+1 < width {
				distributeError(dithered, x+1, y, quantError, 7.0/16.0)
			}
			if x-1 >= 0 && y+1 < height {
				distributeError(dithered, x-1, y+1, quantError, 3.0/16.0)
			}
			if y+1 < height {
				distributeError(dithered, x, y+1, quantError, 5.0/16.0)
			}
			if x+1 < width && y+1 < height {
				distributeError(dithered, x+1, y+1, quantError, 1.0/16.0)
			}
		}
	}

	return dithered
}

// distributeError distributes quantization error to neighboring pixel
func distributeError(img *image.Gray, x, y int, quantError int, factor float64) {
	oldValue := img.GrayAt(x, y).Y
	newValue := int(oldValue) + int(float64(quantError)*factor)

	// Clamp to valid range
	if newValue < 0 {
		newValue = 0
	}
	if newValue > 255 {
		newValue = 255
	}

	img.SetGray(x, y, color.Gray{Y: uint8(newValue)})
}

// sendToPrinter sends data to printer via TCP
func sendToPrinter(ip string, port int, data []byte) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to printer: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("failed to send data to printer: %w", err)
	}

	return nil
}
