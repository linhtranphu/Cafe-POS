# Hướng dẫn In Bill bằng Chromedp (HTML → Ảnh → ESC/POS)

## Tổng quan

Phương pháp này sử dụng Chromedp để render HTML thành ảnh, sau đó chuyển sang ESC/POS để in. Ưu điểm:

- ✅ Thiết kế bill bằng HTML/CSS (dễ dàng, trực quan)
- ✅ Không lo lỗi font tiếng Việt (tất cả đã render thành ảnh)
- ✅ Logo và text được gộp thành 1 ảnh duy nhất
- ✅ Chất lượng in sắc nét nhờ binarization
- ✅ Single executable nhờ go:embed

## Kiến trúc

```
HTML Template (go:embed)
    ↓
Chromedp (render HTML)
    ↓
Screenshot (PNG)
    ↓
Binarization (Black & White)
    ↓
ESC/POS Commands
    ↓
Network Printer (TCP 9100)
```

## Các bước triển khai

### Bước 1: Template HTML

File: `backend/application/services/templates/bill_template.html`

- Chiều rộng cố định: 576px (K80 printer)
- Chỉ dùng màu đen/trắng
- Logo nhúng dưới dạng base64
- Font: Arial (hỗ trợ tiếng Việt tốt)

### Bước 2: Chromedp Renderer

File: `backend/application/services/chromedp_bill_renderer_optimized.go`

Tính năng:
- **go:embed**: Nhúng template HTML vào binary
- **Reusable Context**: Giữ 1 Chrome instance xuyên suốt (tăng tốc)
- **Binarization**: Chuyển ảnh sang đen/trắng thuần túy (threshold = 128)
- **ESC/POS**: Convert ảnh thành lệnh GS v 0 (raster bit image)

### Bước 3: Tối ưu hóa

#### 3.1. go:embed - Single Executable

```go
//go:embed templates/bill_template.html
var billTemplateHTML string
```

Khách hàng chỉ cần 1 file `.exe` duy nhất, không cần copy template riêng.

#### 3.2. Binarization - Ảnh sắc nét

```go
func binarizeImageOptimized(img image.Image, threshold uint8) *image.Gray {
    // Pixel > 128 → White (255)
    // Pixel ≤ 128 → Black (0)
}
```

Loại bỏ hiện tượng "muỗi" (dithering), chữ in ra cực kỳ sắc nét.

#### 3.3. Reusable Chrome Context

```go
// Khởi tạo 1 lần
renderer, _ := NewChromedpBillRendererOptimized()
defer renderer.Close()

// Dùng nhiều lần (nhanh)
for _, order := range orders {
    escpos, _ := renderer.RenderBillToESCPOS(order, settings)
    SendToPrinter(printerIP, escpos)
}
```

Tránh khởi động Chrome mỗi lần in (tiết kiệm ~1s/bill).

## Cách sử dụng

### Test cơ bản

```bash
cd backend/cmd/test-chromedp-print
go run main.go
```

### Tích hợp vào API

```go
// Khởi tạo renderer (1 lần khi start server)
renderer, err := services.NewChromedpBillRendererOptimized()
if err != nil {
    log.Fatal(err)
}
defer renderer.Close()

// Handler
func PrintBill(c *gin.Context) {
    // Fetch order & settings
    order := ...
    settings := ...
    
    // Render
    escpos, err := renderer.RenderBillToESCPOS(order, settings)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    // Send to printer
    err = services.SendToPrinter("192.168.1.100:9100", escpos)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"success": true})
}
```

### Preview (Debug)

```go
// Lưu ảnh preview để kiểm tra
err := renderer.SavePreviewImage(order, settings, "preview.png")
```

## Cấu trúc ESC/POS

```
ESC @ (0x1B 0x40)           - Initialize printer
GS v 0 (0x1D 0x76 0x30 0x00) - Print raster image
  xL xH                      - Width in bytes (little endian)
  yL yH                      - Height (1 line at a time)
  [bitmap data]              - 1 bit per pixel (1=black, 0=white)
ESC d 3 (0x1B 0x64 0x03)    - Feed 3 lines
GS V A 0 (0x1D 0x56 0x41 0x00) - Cut paper
```

## Lưu ý quan trọng

### Logo

- Logo được load từ file và convert sang base64
- Nhúng trực tiếp vào HTML: `<img src="data:image/jpeg;base64,...">`
- Chromedp render logo cùng với text thành 1 ảnh duy nhất

### Font

- Chromedp sử dụng font hệ thống (Arial)
- Không cần cài font riêng vì đã render thành ảnh
- Tiếng Việt hiển thị hoàn hảo

### Threshold (Ngưỡng)

- Mặc định: 128 (giữa 0-255)
- Tăng threshold → nhiều trắng hơn (chữ mỏng)
- Giảm threshold → nhiều đen hơn (chữ đậm)
- Điều chỉnh tùy theo máy in

### Performance

- Lần đầu: ~1-2s (khởi động Chrome)
- Lần sau: ~300-500ms (reuse context)
- Nếu in liên tục, giữ renderer alive

## Dependencies

```bash
go get github.com/chromedp/chromedp
```

Chromedp sẽ tự động download Chrome binary khi chạy lần đầu.

## So sánh với phương pháp cũ

| Tiêu chí | gg (fogleman) | Chromedp |
|----------|---------------|----------|
| Thiết kế | Code Go | HTML/CSS |
| Font Việt | Cần font file | Tự động |
| Logo | Load riêng | Nhúng HTML |
| Độ phức tạp | Cao | Thấp |
| Tốc độ | Nhanh (~100ms) | Trung bình (~500ms) |
| Linh hoạt | Thấp | Cao |

## Kết luận

Chromedp phù hợp khi:
- Cần thiết kế bill phức tạp
- Muốn dễ dàng chỉnh sửa layout
- Không muốn lo font tiếng Việt
- Chấp nhận trade-off về tốc độ (~500ms)

Nếu cần tốc độ tối đa và layout đơn giản, dùng `gg` (fogleman).
