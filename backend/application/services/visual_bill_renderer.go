package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"runtime"
	"strings"

	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/domain/settings"

	"github.com/fogleman/gg"
	"github.com/kenshaw/escpos"
	"github.com/kenshaw/escpos/raster"
	"github.com/nfnt/resize"
	_ "image/jpeg"
)

const (
	ImageWidthPixels = 576
	Margin           = 20
)

// VisualBillRenderer renders bills as images using the exact layout from preview.go
type VisualBillRenderer struct {
	fontPath string
}

// NewVisualBillRenderer creates a new visual bill renderer
func NewVisualBillRenderer() (*VisualBillRenderer, error) {
	fontPath, err := findSystemFont()
	if err != nil {
		return nil, fmt.Errorf("failed to find system font: %w", err)
	}
	
	return &VisualBillRenderer{
		fontPath: fontPath,
	}, nil
}

// RenderBillToESCPOS renders a bill as ESC/POS commands
func (r *VisualBillRenderer) RenderBillToESCPOS(ord *order.Order, shopSettings *settings.ShopSettings) ([]byte, error) {
	// Create bill image
	img, err := r.createBillImage(ord, shopSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to create bill image: %w", err)
	}
	
	// Convert to ESC/POS commands
	escposData, err := r.imageToESCPOS(img)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to ESC/POS: %w", err)
	}
	
	return escposData, nil
}

// RenderBillToBase64 renders a bill as base64-encoded ESC/POS commands
func (r *VisualBillRenderer) RenderBillToBase64(ord *order.Order, shopSettings *settings.ShopSettings) (string, error) {
	escposData, err := r.RenderBillToESCPOS(ord, shopSettings)
	if err != nil {
		return "", err
	}
	
	// Encode to base64
	return base64.StdEncoding.EncodeToString(escposData), nil
}

// createBillImage creates the bill image matching preview.go layout
func (r *VisualBillRenderer) createBillImage(ord *order.Order, shopSettings *settings.ShopSettings) (image.Image, error) {
	dc := gg.NewContext(ImageWidthPixels, 900)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetColor(color.Black)
	
	y := 20.0
	
	// Logo and shop info
	var logoH int
	if shopSettings.ShowLogo && shopSettings.LogoURL != "" {
		logo, err := loadImage(shopSettings.LogoURL)
		if err == nil {
			resizedLogo := resize.Resize(200, 0, logo, resize.Lanczos3)
			logoH = resizedLogo.Bounds().Dy()
			dc.DrawImage(resizedLogo, Margin+20, int(y))
		} else {
			log.Printf("Failed to load logo: %v", err)
		}
	}
	
	// Shop title (right side of logo)
	textX := float64(Margin + 280)
	if err := dc.LoadFontFace(r.fontPath, 25); err != nil {
		return nil, err
	}
	maxWidth := float64(ImageWidthPixels - Margin - 210)
	titleLines := wrapText(dc, shopSettings.ShopName, maxWidth)
	textY := y + 20.0
	for i, line := range titleLines {
		dc.DrawString(line, textX, textY+float64(i*18))
		dc.DrawString(line, textX+2, textY+float64(i*18)) // fake bold
	}
	
	// Address (below title)
	if shopSettings.ShowAddress {
		textX = float64(Margin + 285)
		if err := dc.LoadFontFace(r.fontPath, 16); err != nil {
			return nil, err
		}
		maxWidth = float64(ImageWidthPixels - Margin - 360)
		addressLines := wrapText(dc, shopSettings.ShopAddress, maxWidth)
		textY = y + 50.0
		for i, line := range addressLines {
			dc.DrawString(line, textX, textY+float64(i*18))
		}
	}
	
	// Phone (below address)
	if shopSettings.ShowPhone {
		textX = float64(Margin + 285)
		if err := dc.LoadFontFace(r.fontPath, 16); err != nil {
			return nil, err
		}
		maxWidth = float64(ImageWidthPixels - Margin - 360)
		phoneText := fmt.Sprintf("Hotline: %s", shopSettings.ShopPhone)
		sdtLines := wrapText(dc, phoneText, maxWidth)
		textY = y + 87.0
		for i, line := range sdtLines {
			dc.DrawString(line, textX, textY+float64(i*18))
		}
	}
	
	// Update y after logo
	if logoH > 60 {
		y += float64(logoH) + 45
	} else {
		y += 75
	}
	
	// Title "HÓA ĐƠN THANH TOÁN"
	if err := dc.LoadFontFace(r.fontPath, 34); err != nil {
		return nil, err
	}
	title := "HÓA ĐƠN THANH TOÁN"
	tw, _ := dc.MeasureString(title)
	dc.DrawString(title, (float64(ImageWidthPixels)-tw)/2, y)
	y += 40
	
	// Order info
	if err := dc.LoadFontFace(r.fontPath, 16); err != nil {
		return nil, err
	}
	
	orderNoLine := fmt.Sprintf("Order: %s", ord.OrderNumber)
	dc.DrawString(orderNoLine, Margin+10, y)
	y += 20
	
	waiterLine := fmt.Sprintf("Waiter: %s", ord.WaiterName)
	dc.DrawString(waiterLine, Margin+10, y)
	y += 20
	
	paymentMethodLine := fmt.Sprintf("Thanh Toán: %s", ord.PaymentMethod)
	dc.DrawString(paymentMethodLine, Margin+10, y)
	y += 20
	
	createdDateLine := fmt.Sprintf("Ngày tạo: %s", ord.CreatedAt.Format("02/01/2006 03:04 PM"))
	dc.DrawString(createdDateLine, Margin+10, y)
	y += 20
	
	// Line
	dc.DrawLine(float64(Margin), y, float64(ImageWidthPixels-Margin), y)
	dc.Stroke()
	y += 25
	
	// Table header
	if err := dc.LoadFontFace(r.fontPath, 17); err != nil {
		return nil, err
	}
	colX := []float64{
		float64(Margin + 10),
		float64(Margin + 50),
		float64(Margin + 290),
		float64(Margin + 340),
		float64(Margin + 450),
	}
	
	dc.DrawString("STT", colX[0], y)
	dc.DrawString("Tên món", colX[1], y)
	dc.DrawString("SL", colX[2], y)
	dc.DrawString("Đơn giá", colX[3], y)
	dc.DrawString("Thành tiền", colX[4], y)
	y += 8
	
	dc.DrawLine(float64(Margin), y, float64(ImageWidthPixels-Margin), y)
	dc.Stroke()
	y += 20
	
	// Items
	for i, item := range ord.Items {
		dc.DrawString(fmt.Sprintf("%d", i+1), colX[0], y)
		
		itemName := item.Name
		if item.VariantName != "" {
			itemName = fmt.Sprintf("%s (%s)", item.Name, item.VariantName)
		}
		dc.DrawString(itemName, colX[1], y)
		dc.DrawString(fmt.Sprintf("%d", item.Quantity), colX[2], y)
		dc.DrawString(formatMoney(item.Price), colX[3], y)
		
		totalStr := formatMoney(item.Subtotal)
		tw, _ = dc.MeasureString(totalStr)
		dc.DrawString(totalStr, float64(ImageWidthPixels-Margin)-tw-10, y)
		
		y += 28
	}
	
	y += 5
	dc.DrawLine(float64(Margin), y, float64(ImageWidthPixels-Margin), y)
	dc.Stroke()
	y += 30
	
	// Total
	if err := dc.LoadFontFace(r.fontPath, 24); err != nil {
		return nil, err
	}
	dc.DrawString("TỔNG TIỀN:", colX[2]-50, y)
	totalStr := formatMoney(ord.Total)
	tw, _ = dc.MeasureString(totalStr)
	dc.DrawString(totalStr, float64(ImageWidthPixels-Margin)-tw-10, y)
	y += 30
	
	dc.DrawLine(float64(Margin), y, float64(ImageWidthPixels-Margin), y)
	dc.Stroke()
	y += 30
	
	// Thanks message
	if shopSettings.ShowCustomMessage && shopSettings.CustomMessage != "" {
		if err := dc.LoadFontFace(r.fontPath, 22); err != nil {
			return nil, err
		}
		tw, _ = dc.MeasureString(shopSettings.CustomMessage)
		dc.DrawString(shopSettings.CustomMessage, (float64(ImageWidthPixels)-tw)/2, y)
		y += 30
	}
	
	// Crop to actual height
	finalH := int(y) + 10
	finalImg := dc.Image()
	cropped := image.NewRGBA(image.Rect(0, 0, ImageWidthPixels, finalH))
	for py := 0; py < finalH; py++ {
		for px := 0; px < ImageWidthPixels; px++ {
			cropped.Set(px, py, finalImg.At(px, py))
		}
	}
	
	return cropped, nil
}

// imageToESCPOS converts an image to ESC/POS raster bit image commands using escpos library
func (r *VisualBillRenderer) imageToESCPOS(img image.Image) ([]byte, error) {
	// CRITICAL: The raster.Converter has INVERTED logic!
	// - WHITE pixels (lightness=1.0) → bit=1 → prints BLACK
	// - BLACK pixels (lightness=0.0) → bit=0 → leaves WHITE
	//
	// Solution: INVERT the image before conversion
	invertedImg := invertImageVisual(img)
	
	var buf bytes.Buffer
	e := escpos.New(&buf)
	e.Init()
	
	converter := &raster.Converter{
		MaxWidth:  576,
		Threshold: 0.5,
	}
	
	converter.Print(invertedImg, e)
	
	e.FormfeedN(3)
	e.Cut()
	e.End()
	
	return buf.Bytes(), nil
}

// invertImageVisual inverts an image (black <-> white)
func invertImageVisual(img image.Image) image.Image {
	bounds := img.Bounds()
	inverted := image.NewRGBA(bounds)
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			
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

// Helper functions

func wrapText(dc *gg.Context, text string, maxWidth float64) []string {
	words := strings.Fields(text)
	var lines []string
	var currentLine string
	
	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word
		
		w, _ := dc.MeasureString(testLine)
		if w > maxWidth && currentLine != "" {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			currentLine = testLine
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return lines
}

func formatMoney(amount float64) string {
	amountInt := int(amount)
	if amountInt >= 1000000 {
		return fmt.Sprintf("%d,%03d,%03d", amountInt/1000000, (amountInt%1000000)/1000, amountInt%1000)
	} else if amountInt >= 1000 {
		return fmt.Sprintf("%d,%03d", amountInt/1000, amountInt%1000)
	}
	return fmt.Sprintf("%d", amountInt)
}

func loadImage(path string) (image.Image, error) {
	// Try original path first
	file, err := os.Open(path)
	
	// If failed and path starts with /, try prepending "."
	if err != nil && len(path) > 0 && path[0] == '/' {
		file, err = os.Open("." + path)
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to open image from %s: %w", path, err)
	}
	defer file.Close()
	
	img, _, err := image.Decode(file)
	return img, err
}

func findSystemFont() (string, error) {
	var paths []string
	switch runtime.GOOS {
	case "darwin":
		paths = []string{
			"/System/Library/Fonts/Supplemental/Arial.ttf",
			"/System/Library/Fonts/Helvetica.ttc",
		}
	case "windows":
		paths = []string{
			"C:\\Windows\\Fonts\\arial.ttf",
		}
	case "linux":
		paths = []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		}
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("font not found")
}

// SaveImagePreview saves the bill image as PNG for preview
func (r *VisualBillRenderer) SaveImagePreview(ord *order.Order, shopSettings *settings.ShopSettings, filename string) error {
	img, err := r.createBillImage(ord, shopSettings)
	if err != nil {
		return err
	}
	
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	
	return png.Encode(f, img)
}
