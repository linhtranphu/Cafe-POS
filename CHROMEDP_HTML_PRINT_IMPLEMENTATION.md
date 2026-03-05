# Chromedp HTML Print Implementation

## Tổng quan

Đã implement tính năng in hóa đơn bằng HTML template với Chromedp vào Print Management. Tính năng này cho phép:

1. Render HTML template thành ảnh bằng Chromedp (headless Chrome)
2. Convert ảnh thành ESC/POS commands
3. Gửi trực tiếp đến máy in qua network

## Kiến trúc

### Frontend

**Component mới:** `frontend/src/components/printing/ChromedpBillPrinter.vue`
- UI để chọn order và in bill
- Tìm kiếm orders
- Preview bill trước khi in
- Gửi lệnh in đến backend

**Cập nhật:** `frontend/src/views/PrintManagementView.vue`
- Thêm tab "HTML Print" (🖼️)
- Import và hiển thị ChromedpBillPrinter component

### Backend

**Handler mới:** `backend/interfaces/http/chromedp_print_handler.go`
- `POST /api/manager/chromedp-print/bill` - In bill bằng chromedp
- `GET /api/manager/chromedp-print/preview/:order_id` - Tạo preview PNG

**Service:** `backend/application/services/chromedp_bill_renderer_optimized.go`
- Sử dụng go:embed để nhúng HTML template
- Reusable Chrome context để tăng tốc độ
- Binarization (ngưỡng đen/trắng) cho output sắc nét
- Convert image sang ESC/POS raster commands

**Template:** `backend/application/services/templates/bill_template_optimized.html`
- HTML template với CSS inline
- Chiều rộng cố định 576px (K80 printer)
- Layout giống preview.go:
  - Logo + thông tin shop
  - Tiêu đề "HÓA ĐƠN THANH TOÁN"
  - Thông tin order (order number, waiter, payment method, date)
  - Bảng items
  - Tổng tiền
  - Lời cảm ơn

**Routes:** `backend/main.go`
- Khởi tạo chromedpPrintHandler
- Register routes `/chromedp-print/bill` và `/chromedp-print/preview/:order_id`
- Cleanup handler khi shutdown

## Workflow

### 1. In Bill

```
User chọn order → Click "In Bill (Chromedp)" 
→ Frontend gửi POST /api/manager/chromedp-print/bill
→ Backend:
  1. Fetch order từ database
  2. Fetch shop settings
  3. Render HTML template với data
  4. Chromedp capture HTML thành PNG
  5. Binarization (đen/trắng)
  6. Convert PNG → ESC/POS commands
  7. Gửi đến máy in qua TCP (port 9100)
→ Response success/error
```

### 2. Preview Bill

```
User chọn order → Click "Preview"
→ Frontend gửi GET /api/manager/chromedp-print/preview/:order_id
→ Backend:
  1. Fetch order từ database
  2. Fetch shop settings
  3. Render HTML template với data
  4. Chromedp capture HTML thành PNG
  5. Binarization
  6. Lưu file PNG (preview_chromedp_[order_number].png)
→ Response với filename
```

## So sánh với Visual Print

| Feature | Visual Print (gg) | Chromedp Print |
|---------|------------------|----------------|
| Rendering | fogleman/gg (Go graphics) | Chromedp (HTML/CSS) |
| Template | Go code | HTML template |
| Flexibility | Thấp (hard-coded) | Cao (HTML/CSS) |
| Performance | Nhanh hơn | Chậm hơn (Chrome overhead) |
| Font support | System fonts | Web fonts |
| Maintenance | Khó (phải code layout) | Dễ (edit HTML/CSS) |

## Ưu điểm Chromedp

1. **Dễ customize**: Chỉnh layout bằng HTML/CSS thay vì code Go
2. **WYSIWYG**: Preview trong browser = output cuối cùng
3. **Rich formatting**: Hỗ trợ đầy đủ CSS (gradients, shadows, etc.)
4. **Web fonts**: Có thể dùng Google Fonts, custom fonts
5. **Reusable**: Template có thể dùng cho nhiều mục đích khác

## Nhược điểm Chromedp

1. **Performance**: Chậm hơn gg renderer (Chrome startup overhead)
2. **Memory**: Tốn RAM hơn (Chrome process)
3. **Dependencies**: Cần Chrome/Chromium installed
4. **Complexity**: Thêm layer phức tạp (HTML → Image → ESC/POS)

## Cấu hình

### Printer IP
Mặc định: `192.168.1.115:9100`

### Template Location
`backend/application/services/templates/bill_template_optimized.html`

### Bill Width
576 pixels (K80 thermal printer - 72mm paper)

## Testing

### Test Backend
```bash
cd backend
go run cmd/test-chromedp-print/main.go
```

### Test Frontend
1. Mở http://localhost:5173/#/print-management
2. Click tab "HTML Print"
3. Chọn order
4. Click "Preview" để xem ảnh
5. Click "In Bill (Chromedp)" để in

## Files Changed/Created

### Frontend
- ✅ `frontend/src/components/printing/ChromedpBillPrinter.vue` (NEW)
- ✅ `frontend/src/views/PrintManagementView.vue` (UPDATED)

### Backend
- ✅ `backend/interfaces/http/chromedp_print_handler.go` (NEW)
- ✅ `backend/application/services/chromedp_bill_renderer_optimized.go` (EXISTING)
- ✅ `backend/application/services/templates/bill_template_optimized.html` (EXISTING)
- ✅ `backend/main.go` (UPDATED - routes & handler init)
- ❌ `backend/application/services/chromedp_bill_renderer.go` (DELETED - duplicate)

## Next Steps

1. **Template Editor**: Tạo UI để edit HTML template trực tiếp
2. **Multiple Templates**: Hỗ trợ nhiều templates cho các loại bill khác nhau
3. **Custom CSS**: Cho phép user customize CSS
4. **Font Management**: Upload và quản lý custom fonts
5. **Performance**: Cache Chrome context, optimize rendering
6. **Error Handling**: Better error messages và retry logic

## Troubleshooting

### Chrome not found
```bash
# macOS
brew install chromium

# Linux
apt-get install chromium-browser

# Windows
Download from https://www.chromium.org/getting-involved/download-chromium
```

### Slow rendering
- Giảm `chromedp.Sleep()` duration
- Sử dụng cached Chrome context (đã implement)
- Preload fonts

### Print quality issues
- Điều chỉnh `ThresholdValue` trong chromedp_bill_renderer_optimized.go
- Mặc định: 128 (0-255)
- Tăng = nhiều trắng hơn
- Giảm = nhiều đen hơn

### Network printer connection failed
- Kiểm tra IP và port (9100)
- Ping printer: `ping 192.168.1.115`
- Telnet test: `telnet 192.168.1.115 9100`
- Kiểm tra firewall

## Kết luận

Đã implement thành công tính năng in hóa đơn bằng HTML template với Chromedp. Tính năng này cung cấp flexibility cao cho việc customize layout và styling, phù hợp cho các use case cần rich formatting và dễ dàng maintenance.
