package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	_ "image/jpeg"
)

const (
	ImageWidthPixels = 576
	Margin           = 20
)

type OrderItem struct {
	STT       int
	Name      string
	Quantity  int
	UnitPrice int
	Total     int
}

func main() {
	fmt.Println("🖼️  Preview template hóa đơn")

	billData := map[string]string{
		"address": "Đ/c: 10/8 Trần Nhật Duật, P. Tân Định, Quận 1, HCM",
		"thanks":  "Cảm ơn quý khách!",
		"title" : "Tiệm cà phê Ông Tạ",
		"sdt" : "Hotline: 0906990602",
		"orderNo" : "Order: 20260222-095703-168",
		"Waiter" : "Waiter: Waiter1",
		"CreatedDate" :"Ngày tạo: 26/02/2026 12:45 AM",
		"PaymentMethod" :"Thanh Toán: Tiền mặt",

	}

	items := []OrderItem{
		{STT: 1, Name: "Cà phê sữa đá", Quantity: 2, UnitPrice: 25000, Total: 50000},
		{STT: 2, Name: "Trà đào cam sả", Quantity: 1, UnitPrice: 35000, Total: 35000},
		{STT: 3, Name: "Bánh mì", Quantity: 1, UnitPrice: 20000, Total: 20000},
	}

	totalAmount := 0
	for _, item := range items {
		totalAmount += item.Total
	}

	fmt.Println("Đang tạo preview...")
	img, err := createBillImagePreview(billData, items, totalAmount)
	if err != nil {
		log.Fatalf("Lỗi: %v", err)
	}

	filename := "preview_bill.png"
	err = saveImagePreview(img, filename)
	if err != nil {
		log.Fatalf("Lỗi lưu file: %v", err)
	}

	fmt.Printf("✓ Đã tạo preview: %s\n", filename)
	fmt.Printf("  Kích thước: %dx%d px\n", img.Bounds().Dx(), img.Bounds().Dy())
	
	openFilePreview(filename)
}

func createBillImagePreview(data map[string]string, items []OrderItem, totalAmount int) (image.Image, error) {
	dc := gg.NewContext(ImageWidthPixels, 900)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	fontPath, err := findSystemFontPreview()
	if err != nil {
		return nil, err
	}

	dc.SetColor(color.Black)
	y := 20.0

	// Logo bên trái, địa chỉ bên phải
	logoPath := "../../../backend/uploads/logos/logo_24094.jpeg"
	logo, err := loadImagePreview(logoPath)
	var logoH int
	if err == nil {
		resizedLogo := resize.Resize(200, 0, logo, resize.Lanczos3)
		_ = resizedLogo.Bounds().Dx()
		logoH = resizedLogo.Bounds().Dy()
		
		// Vẽ logo bên trái
		dc.DrawImage(resizedLogo, Margin+20, int(y))
	}

	// Vẽ title bên phải logo
	textX := float64(Margin + 280)
	dc.LoadFontFace(fontPath, 25)
	maxWidth := float64(ImageWidthPixels - Margin - 210)
	titleLines := wrapTextPreview(dc, data["title"], maxWidth)
	textY := y + 20.0
	for i, line := range titleLines {
		dc.DrawString(line, textX, textY+float64(i*18))
		dc.DrawString(line, textX+2, textY+float64(i*18)) //fake bold
	}

	// Vẽ địa chỉ bên dứoi title
	textX = float64(Margin + 285)
	dc.LoadFontFace(fontPath, 16)
	maxWidth = float64(ImageWidthPixels - Margin - 360)
	addressLines := wrapTextPreview(dc, data["address"], maxWidth)
	textY = y + 50.0
	for i, line := range addressLines {
		dc.DrawString(line, textX, textY+float64(i*18))
	}

	// Vẽ sdt bên dứoi address
	textX = float64(Margin + 285)
	dc.LoadFontFace(fontPath, 16)
	maxWidth = float64(ImageWidthPixels - Margin - 360)
	sdtLines := wrapTextPreview(dc, data["sdt"], maxWidth)
	textY = y + 87.0
	for i, line := range sdtLines {
		dc.DrawString(line, textX, textY+float64(i*18))
	}


	// Cập nhật y sau logo
	if logoH > 60 {
		y += float64(logoH) + 45
	} else {
		y += 75
	}

	// Tiêu đề "HÓA ĐƠN THANH TOÁN"
	dc.LoadFontFace(fontPath, 34)
	title := "HÓA ĐƠN THANH TOÁN"
	tw, _ := dc.MeasureString(title)
	dc.DrawString(title, (float64(ImageWidthPixels)-tw)/2, y)
	y += 40

	// orderno
	dc.LoadFontFace(fontPath, 16)
	orderNoLine := data["orderNo"]
	tw, _ = dc.MeasureString(orderNoLine)
	dc.DrawString(orderNoLine, Margin+10, y)
	y += 20

	// Waiter
	dc.LoadFontFace(fontPath, 16)
	waiterLine := data["Waiter"]
	tw, _ = dc.MeasureString(waiterLine)
	dc.DrawString(waiterLine, Margin+10, y)
	y += 20

	// PaymentMethod
	dc.LoadFontFace(fontPath, 16)
	paymentMethodLine := data["PaymentMethod"]
	tw, _ = dc.MeasureString(paymentMethodLine)
	dc.DrawString(paymentMethodLine, Margin+10, y)
	y += 20

	// CreatedDate
	dc.LoadFontFace(fontPath, 16)
	createdDateLine := data["CreatedDate"]
	tw, _ = dc.MeasureString(createdDateLine)
	dc.DrawString(createdDateLine, Margin+10, y)
	y += 20



	
	dc.DrawLine(float64(Margin), y, float64(ImageWidthPixels-Margin), y)
	dc.Stroke()
	y += 25

	dc.LoadFontFace(fontPath, 17)
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

	for _, item := range items {
		dc.DrawString(fmt.Sprintf("%d", item.STT), colX[0], y)
		dc.DrawString(item.Name, colX[1], y)
		dc.DrawString(fmt.Sprintf("%d", item.Quantity), colX[2], y)
		dc.DrawString(formatMoneyPreview(item.UnitPrice), colX[3], y)
		totalStr := formatMoneyPreview(item.Total)
		tw, _ = dc.MeasureString(totalStr)
		dc.DrawString(totalStr, float64(ImageWidthPixels-Margin)-tw-10, y)
		y += 28
	}

	y += 5
	dc.DrawLine(float64(Margin), y, float64(ImageWidthPixels-Margin), y)
	dc.Stroke()
	y += 30

	dc.LoadFontFace(fontPath, 24)
	dc.DrawString("TỔNG TIỀN:", colX[2]-50, y)
	totalStr := formatMoneyPreview(totalAmount)
	tw, _ = dc.MeasureString(totalStr)
	dc.DrawString(totalStr, float64(ImageWidthPixels-Margin)-tw-10, y)
	y += 30

	dc.DrawLine(float64(Margin), y, float64(ImageWidthPixels-Margin), y)
	dc.Stroke()
	y += 30

	dc.LoadFontFace(fontPath, 22)
	tw, _ = dc.MeasureString(data["thanks"])
	dc.DrawString(data["thanks"], (float64(ImageWidthPixels)-tw)/2, y)
	y += 30

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

func wrapTextPreview(dc *gg.Context, text string, maxWidth float64) []string {
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

func formatMoneyPreview(amount int) string {
	if amount >= 1000000 {
		return fmt.Sprintf("%d,%03d,%03d", amount/1000000, (amount%1000000)/1000, amount%1000)
	} else if amount >= 1000 {
		return fmt.Sprintf("%d,%03d", amount/1000, amount%1000)
	}
	return fmt.Sprintf("%d", amount)
}

func loadImagePreview(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	return img, err
}

func findSystemFontPreview() (string, error) {
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

func saveImagePreview(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func openFilePreview(filename string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", filename)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", filename)
	case "linux":
		cmd = exec.Command("xdg-open", filename)
	}
	if cmd != nil {
		cmd.Run()
	}
}
