package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net"
	"os"
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
	fmt.Println("🖨️  Test in hóa đơn với logo")
	
	printerIP := "192.168.1.115:9100"
	conn, err := net.Dial("tcp", printerIP)
	if err != nil {
		log.Fatalf("Lỗi kết nối: %v", err)
	}
	defer conn.Close()
	fmt.Println("✓ Đã kết nối máy in")

	billData := map[string]string{
		"address":     "Đ/c: 10/8 Trần Nhật Duật, P. Tân Định, Quận 1, HCM",
		"thanks":      "Cảm ơn quý khách!",
		"title":       "Tiệm cà phê Ông Tạ",
		"sdt":         "Hotline: 0906990602",
		"orderNo":     "Order: 20260222-095703-168",
		"Waiter":      "Waiter: Waiter1",
		"CreatedDate": "Ngày tạo: 26/02/2026 12:45 AM",
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

	fmt.Println("Đang tạo ảnh...")
	img, err := createBillImage(billData, items, totalAmount)
	if err != nil {
		log.Fatalf("Lỗi tạo ảnh: %v", err)
	}
	fmt.Printf("✓ Ảnh: %dx%d px\n", img.Bounds().Dx(), img.Bounds().Dy())

	saveImageDebug(img, "debug_bill.png")
	fmt.Println("✓ Đã lưu debug_bill.png")

	conn.Write([]byte{0x1B, 0x40})
	fmt.Println("✓ Đã init máy in")

	fmt.Println("Đang chuyển đổi ảnh sang ESC/POS...")
	bitmapData := convertImageToESCPOS(img)
	
	fmt.Printf("Đang gửi %d bytes...\n", len(bitmapData))
	_, err = conn.Write(bitmapData)
	if err != nil {
		log.Fatalf("Lỗi gửi dữ liệu: %v", err)
	}
	fmt.Println("✓ Đã gửi ảnh")

	conn.Write([]byte{0x1B, 0x64, 0x03})
	conn.Write([]byte{0x1D, 0x56, 0x42, 0x00})
	
	fmt.Println("✓ Hoàn tất!")
}

func convertImageToESCPOS(img image.Image) []byte {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	widthBytes := (width + 7) / 8
	bitmap := make([]byte, widthBytes*height)
	
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := (r + g + b) / 3
			if gray < 32768 {
				byteIndex := y*widthBytes + x/8
				bitIndex := 7 - (x % 8)
				bitmap[byteIndex] |= (1 << bitIndex)
			}
		}
	}
	
	var buf bytes.Buffer
	buf.WriteByte(0x1D)
	buf.WriteByte(0x76)
	buf.WriteByte(0x30)
	buf.WriteByte(0x00)
	buf.WriteByte(byte(widthBytes & 0xFF))
	buf.WriteByte(byte((widthBytes >> 8) & 0xFF))
	buf.WriteByte(byte(height & 0xFF))
	buf.WriteByte(byte((height >> 8) & 0xFF))
	buf.Write(bitmap)
	
	return buf.Bytes()
}

func createBillImage(data map[string]string, items []OrderItem, totalAmount int) (image.Image, error) {
	dc := gg.NewContext(ImageWidthPixels, 900)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	fontPath, err := findSystemFont()
	if err != nil {
		return nil, err
	}

	dc.SetColor(color.Black)
	y := 20.0

	// Logo bên trái, địa chỉ bên phải
	logoPath := "../../../backend/uploads/logos/logo_24094.jpeg"
	logo, err := loadImage(logoPath)
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
	titleLines := wrapText(dc, data["title"], maxWidth)
	textY := y + 20.0
	for i, line := range titleLines {
		dc.DrawString(line, textX, textY+float64(i*18))
		dc.DrawString(line, textX+2, textY+float64(i*18)) //fake bold
	}

	// Vẽ địa chỉ bên dứoi title
	textX = float64(Margin + 285)
	dc.LoadFontFace(fontPath, 16)
	maxWidth = float64(ImageWidthPixels - Margin - 360)
	addressLines := wrapText(dc, data["address"], maxWidth)
	textY = y + 50.0
	for i, line := range addressLines {
		dc.DrawString(line, textX, textY+float64(i*18))
	}

	// Vẽ sdt bên dứoi address
	textX = float64(Margin + 285)
	dc.LoadFontFace(fontPath, 16)
	maxWidth = float64(ImageWidthPixels - Margin - 360)
	sdtLines := wrapText(dc, data["sdt"], maxWidth)
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
		dc.DrawString(formatMoney(item.UnitPrice), colX[3], y)
		totalStr := formatMoney(item.Total)
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
	totalStr := formatMoney(totalAmount)
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

func formatMoney(amount int) string {
	if amount >= 1000000 {
		return fmt.Sprintf("%d,%03d,%03d", amount/1000000, (amount%1000000)/1000, amount%1000)
	} else if amount >= 1000 {
		return fmt.Sprintf("%d,%03d", amount/1000, amount%1000)
	}
	return fmt.Sprintf("%d", amount)
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
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

func saveImageDebug(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
