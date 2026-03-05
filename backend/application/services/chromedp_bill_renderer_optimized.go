package services

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"image/color"
	"image/png"
	"log"
	"net/url"
	"os"
	"time"

	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/domain/settings"

	"github.com/chromedp/chromedp"
)

//go:embed templates/bill_template_optimized.html
var billTemplateOptimizedHTML string

const (
	BillWidthPixelsOptimized = 576 // K80 printer width
	ThresholdValue           = 128 // Binarization threshold (0-255)
)

// ChromedpBillRendererOptimized renders bills using Chromedp with optimizations:
// 1. go:embed for single executable
// 2. Reusable Chrome context for speed
// 3. Binarization for sharp black & white output
type ChromedpBillRendererOptimized struct {
	ctx    context.Context
	cancel context.CancelFunc
	tmpl   *template.Template
}

// NewChromedpBillRendererOptimized creates a new optimized renderer
func NewChromedpBillRendererOptimized() (*ChromedpBillRendererOptimized, error) {
	// Parse embedded template
	tmpl, err := template.New("bill").Parse(billTemplateOptimizedHTML)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Create reusable Chrome context with optimized settings for bill printing
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		// Set font rendering flags for better Vietnamese support
		chromedp.Flag("font-render-hinting", "none"),
		chromedp.Flag("disable-font-subpixel-positioning", false),
		// Window size for 80mm bill (576px width)
		// Height set large enough for long bills, FullScreenshot will capture actual content
		chromedp.WindowSize(576, 4000),
	)
	
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, _ := chromedp.NewContext(allocCtx)

	// Initialize Chrome (warm up)
	if err := chromedp.Run(ctx); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize Chrome: %w", err)
	}

	return &ChromedpBillRendererOptimized{
		ctx:    ctx,
		cancel: cancel,
		tmpl:   tmpl,
	}, nil
}

// Close cleans up Chrome instance
func (r *ChromedpBillRendererOptimized) Close() {
	if r.cancel != nil {
		r.cancel()
	}
}

// RenderBillToESCPOS renders bill as ESC/POS commands
func (r *ChromedpBillRendererOptimized) RenderBillToESCPOS(ord *order.Order, shopSettings *settings.ShopSettings) ([]byte, error) {
	// 1. Prepare data
	billData, err := r.prepareBillData(ord, shopSettings)
	if err != nil {
		return nil, fmt.Errorf("prepare data failed: %w", err)
	}

	// DEBUG: Log logo status
	log.Printf("DEBUG: ShowLogo=%v, LogoBase64 length=%d", billData.ShowLogo, len(billData.LogoBase64))

	// 2. Render HTML
	var htmlBuf bytes.Buffer
	if err := r.tmpl.Execute(&htmlBuf, billData); err != nil {
		return nil, fmt.Errorf("template execution failed: %w", err)
	}

	htmlStr := htmlBuf.String()

	// DEBUG: Save HTML to file for inspection
	if err := os.WriteFile("debug_rendered.html", []byte(htmlStr), 0644); err == nil {
		log.Printf("DEBUG: Saved rendered HTML to debug_rendered.html")
	}

	// DEBUG: Check if logo is in HTML
	if billData.ShowLogo {
		if len(billData.LogoBase64) > 0 {
			log.Printf("DEBUG: Logo base64 is in template data")
			// Check if it's in rendered HTML
			if bytes.Contains([]byte(htmlStr), []byte("data:image/png;base64,")) {
				log.Printf("DEBUG: ✅ Logo base64 found in rendered HTML")
			} else {
				log.Printf("DEBUG: ❌ Logo base64 NOT found in rendered HTML")
			}
		} else {
			log.Printf("DEBUG: ⚠️ LogoBase64 is empty, will use fallback")
		}
	}

	// 3. Capture with Chromedp
	img, err := r.captureHTML(htmlStr)
	if err != nil {
		return nil, fmt.Errorf("capture failed: %w", err)
	}

	// 4. Convert to grayscale (better for logos than binarization)
	grayImg := convertToGrayscale(img)

	// 5. Convert to ESC/POS
	escposData := imageToESCPOSOptimized(grayImg)

	return escposData, nil
}

// prepareBillData converts order to template data
func (r *ChromedpBillRendererOptimized) prepareBillData(ord *order.Order, shopSettings *settings.ShopSettings) (*BillTemplateDataOptimized, error) {
	data := &BillTemplateDataOptimized{
		ShopName:      shopSettings.ShopName,
		ShopAddress:   shopSettings.ShopAddress,
		ShopPhone:     shopSettings.ShopPhone,
		ShowLogo:      shopSettings.ShowLogo,
		ShowAddress:   shopSettings.ShowAddress,
		ShowPhone:     shopSettings.ShowPhone,
		OrderNumber:   ord.OrderNumber,
		CustomerName:  ord.CustomerName,
		WaiterName:    ord.WaiterName,
		PaymentMethod: string(ord.PaymentMethod),
		CreatedDate:   ord.CreatedAt.Format("02/01/2006 03:04 PM"),
		Total:         formatMoneyVN(ord.Total),
		CustomMessage: shopSettings.CustomMessage,
		ShowCustomMsg: shopSettings.ShowCustomMessage,
	}

	// Load logo as base64 (embedded in HTML)
	if shopSettings.ShowLogo && shopSettings.LogoURL != "" {
		logoBase64, err := loadImageAsBase64(shopSettings.LogoURL)
		if err != nil {
			log.Printf("ERROR: Failed to load logo from %s: %v", shopSettings.LogoURL, err)
		} else {
			previewLen := 50
			if len(logoBase64) < previewLen {
				previewLen = len(logoBase64)
			}
			log.Printf("SUCCESS: Logo loaded: %d bytes (preview: %s...)", len(logoBase64), logoBase64[:previewLen])
			data.LogoBase64 = template.URL(logoBase64) // Convert to template.URL for src attribute
		}
	} else {
		log.Printf("INFO: Logo disabled or no URL (ShowLogo=%v, LogoURL=%s)", shopSettings.ShowLogo, shopSettings.LogoURL)
	}

	// Convert items
	for i, item := range ord.Items {
		itemName := item.Name
		if item.VariantName != "" {
			itemName = fmt.Sprintf("%s (%s)", item.Name, item.VariantName)
		}

		data.Items = append(data.Items, BillItemDataOptimized{
			STT:       i + 1,
			Name:      itemName,
			Quantity:  item.Quantity,
			UnitPrice: formatMoneyVN(item.Price),
			Total:     formatMoneyVN(item.Subtotal),
		})
	}

	return data, nil
}

// captureHTML uses Chromedp to render HTML and capture screenshot
func (r *ChromedpBillRendererOptimized) captureHTML(html string) (image.Image, error) {
	var buf []byte

	// Sử dụng data URL với các chốt chặn để đảm bảo render hoàn chỉnh
	dataURL := "data:text/html;charset=utf-8," + url.PathEscape(html)

	log.Printf("Capturing HTML with chromedp (data URL approach)")

	// Navigate và capture với các chốt chặn để đảm bảo render hoàn toàn
	err := chromedp.Run(r.ctx,
		// Điều hướng tới data URL chứa HTML
		chromedp.Navigate(dataURL),

		// Set viewport to exact bill width (80mm = 576px at 72 DPI)
		chromedp.EmulateViewport(576, 1),

		// CHỐT CHẶN 1: Đợi body element xuất hiện
		chromedp.WaitReady("body"),

		// CHỐT CHẶN 2: Đợi fonts load xong
		chromedp.Evaluate(`
			document.fonts.ready.then(() => {
				console.log('Fonts loaded');
			});
		`, nil),

		// CHỐT CHẶN 3: Đợi tất cả images load xong (nếu có)
		chromedp.Evaluate(`
			new Promise((resolve) => {
				const images = document.querySelectorAll('img');
				if (images.length === 0) {
					resolve();
					return;
				}
				let loaded = 0;
				const checkComplete = () => {
					loaded++;
					if (loaded === images.length) resolve();
				};
				images.forEach(img => {
					if (img.complete) {
						checkComplete();
					} else {
						img.onload = checkComplete;
						img.onerror = checkComplete;
					}
				});
				// Timeout sau 2 giây
				setTimeout(resolve, 2000);
			});
		`, nil),

		// CHỐT CHẶN 4: Đợi thêm để font/CSS render hoàn toàn
		chromedp.Sleep(300 * time.Millisecond),

		// CHỐT CHẶN 5: Chụp toàn bộ trang với quality 100
		chromedp.FullScreenshot(&buf, 100),
	)
	if err != nil {
		return nil, fmt.Errorf("chromedp run failed: %w", err)
	}

	log.Printf("Screenshot captured: %d bytes", len(buf))

	// Decode PNG
	img, err := png.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, fmt.Errorf("png decode failed: %w", err)
	}

	log.Printf("Image decoded: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())

	// Crop to bill width (should already be 576px but ensure it)
	return cropToWidthOptimized(img, BillWidthPixelsOptimized), nil
}

// convertToGrayscale converts image to grayscale (preserves more detail than binarization)
func convertToGrayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray.Set(x, y, img.At(x, y))
		}
	}

	return gray
}

// binarizeImageOptimized converts to pure black & white (threshold)
// Pixel > threshold → White (255)
// Pixel ≤ threshold → Black (0)
func binarizeImageOptimized(img image.Image, threshold uint8) *image.Gray {
	bounds := img.Bounds()
	bw := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			
			// Convert to grayscale (0-255)
			gray := uint8((r + g + b) / 3 / 257)
			
			// Apply threshold
			if gray > threshold {
				bw.SetGray(x, y, color.Gray{Y: 255}) // White
			} else {
				bw.SetGray(x, y, color.Gray{Y: 0}) // Black
			}
		}
	}

	return bw
}

// cropToWidthOptimized crops image to specified width (centered)
func cropToWidthOptimized(img image.Image, width int) image.Image {
	bounds := img.Bounds()
	if bounds.Dx() <= width {
		return img
	}

	// Center crop
	offsetX := (bounds.Dx() - width) / 2
	cropped := image.NewRGBA(image.Rect(0, 0, width, bounds.Dy()))

	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < width; x++ {
			cropped.Set(x, y, img.At(x+offsetX, y))
		}
	}

	return cropped
}

// imageToESCPOSOptimized converts image to ESC/POS raster commands using escpos library
func imageToESCPOSOptimized(img image.Image) []byte {
	// Use GS v 0 (raster bit image) instead of GS 8 L
	// GS v 0 is more widely supported by ESC/POS printers
	
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
		
		// yL yH - height (1 line at a time for better compatibility)
		buf.WriteByte(0x01)
		buf.WriteByte(0x00)
		
		// Convert line to bitmap
		lineData := make([]byte, bytesPerLine)
		for x := 0; x < width; x++ {
			// Get pixel color
			r, g, b, _ := img.At(x, y).RGBA()
			// Convert to grayscale (weighted average)
			gray := (r*299 + g*587 + b*114) / 1000
			
			// If pixel is dark (< 50%), set bit to 1 (print black)
			// Note: RGBA values are 16-bit (0-65535), so 50% = 32768
			if gray < 32768 {
				byteIndex := x / 8
				bitIndex := 7 - (x % 8)
				lineData[byteIndex] |= (1 << bitIndex)
			}
		}
		buf.Write(lineData)
	}
	
	// Feed paper (3 lines)
	buf.Write([]byte{0x1B, 0x64, 0x03})
	
	// Cut paper (full cut)
	buf.Write([]byte{0x1D, 0x56, 0x00})
	
	return buf.Bytes()
}

// invertImage inverts an image (black <-> white)
func invertImage(img image.Image) image.Image {
	bounds := img.Bounds()
	inverted := image.NewRGBA(bounds)
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			
			// Invert RGB values (keep alpha)
			inverted.Set(x, y, color.RGBA{
				R: uint8(255 - (r >> 8)),
				G: uint8(255 - (g >> 8)),
				B: uint8(255 - (b >> 8)),
				A: uint8(a >> 8),
			})
		}
	}
	
	return inverted
}

// BillTemplateDataOptimized represents data for HTML template
type BillTemplateDataOptimized struct {
	ShopName      string
	ShopAddress   string
	ShopPhone     string
	LogoBase64    template.URL // Changed to template.URL for src attribute
	ShowLogo      bool
	ShowAddress   bool
	ShowPhone     bool
	OrderNumber   string
	CustomerName  string
	WaiterName    string
	PaymentMethod string
	CreatedDate   string
	Items         []BillItemDataOptimized
	Total         string
	CustomMessage string
	ShowCustomMsg bool
}

type BillItemDataOptimized struct {
	STT       int
	Name      string
	Quantity  int
	UnitPrice string
	Total     string
}

// Helper functions

func formatMoneyVN(amount float64) string {
	amountInt := int(amount)
	if amountInt >= 1000000 {
		return fmt.Sprintf("%d,%03d,%03d", amountInt/1000000, (amountInt%1000000)/1000, amountInt%1000)
	} else if amountInt >= 1000 {
		return fmt.Sprintf("%d,%03d", amountInt/1000, amountInt%1000)
	}
	return fmt.Sprintf("%d", amountInt)
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

	// Decode image to resize it
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}
	
	// Resize to max width 200px to reduce base64 size
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	if width > 200 {
		// Calculate new height maintaining aspect ratio
		newWidth := 200
		newHeight := (height * newWidth) / width
		
		// Simple resize by sampling
		resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		for y := 0; y < newHeight; y++ {
			for x := 0; x < newWidth; x++ {
				srcX := (x * width) / newWidth
				srcY := (y * height) / newHeight
				resized.Set(x, y, img.At(srcX, srcY))
			}
		}
		img = resized
		log.Printf("Logo resized from %dx%d to %dx%d", width, height, newWidth, newHeight)
	}
	
	// Encode resized image to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", fmt.Errorf("failed to encode resized image: %w", err)
	}
	
	// Convert to base64
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	result := fmt.Sprintf("data:image/png;base64,%s", encoded)
	
	log.Printf("Logo base64 size: %d bytes (original: %d bytes)", len(result), len(data))
	return result, nil
}

// SavePreviewImage saves bill as PNG for debugging
func (r *ChromedpBillRendererOptimized) SavePreviewImage(ord *order.Order, shopSettings *settings.ShopSettings, filename string) error {
	billData, err := r.prepareBillData(ord, shopSettings)
	if err != nil {
		return err
	}

	var htmlBuf bytes.Buffer
	if err := r.tmpl.Execute(&htmlBuf, billData); err != nil {
		return err
	}

	img, err := r.captureHTML(htmlBuf.String())
	if err != nil {
		return err
	}

	// Save RAW capture first (for debugging)
	rawFilename := "raw_" + filename
	fRaw, err := os.Create(rawFilename)
	if err == nil {
		png.Encode(fRaw, img)
		fRaw.Close()
		log.Printf("Saved raw capture: %s", rawFilename)
	}

	// Convert to grayscale (same as print)
	grayImg := convertToGrayscale(img)

	// Save as PNG
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, grayImg)
}
