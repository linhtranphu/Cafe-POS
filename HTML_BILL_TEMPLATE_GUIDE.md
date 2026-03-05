# Hướng dẫn In Bill với HTML Template

## Tổng quan

Hệ thống in bill đã được cập nhật để sử dụng **HTML template** với CSS chính xác theo thông số của `preview.go`. Template HTML được render thành hình ảnh bằng headless Chrome, sau đó convert sang ESC/POS và in.

## Ưu điểm HTML Template

✅ **Dễ customize**: Chỉnh sửa HTML/CSS thay vì code Go  
✅ **Chính xác**: CSS matching pixel-perfect với preview.go  
✅ **Maintain dễ**: Không cần recompile khi thay đổi layout  
✅ **Preview nhanh**: Mở HTML trong browser để xem  
✅ **Flexible**: Dễ thêm/bớt elements, thay đổi style  

## Cấu trúc

### 1. HTML Template

**File**: `backend/templates/bill_template.html`

**Thông số chính xác từ preview.go**:
- Width: 576px (72mm @ 203 DPI)
- Margin: 20px
- Container width: 536px (576 - 40)

**Layout**:
```
┌─────────────────────────────────────┐
│  [Logo 200px]  [Shop Info]          │ ← Header
│                                     │
│     HÓA ĐƠN THANH TOÁN             │ ← Title (34px, centered)
│                                     │
│  Order: xxx                         │
│  Waiter: xxx                        │ ← Order Info (16px)
│  Thanh Toán: xxx                    │
│  Ngày tạo: xxx                      │
│  ─────────────────────────────────  │
│  STT | Tên món | SL | Giá | Total  │ ← Table Header (17px)
│  ─────────────────────────────────  │
│  1   | Item 1  | 2  | 25k | 50k    │
│  2   | Item 2  | 1  | 35k | 35k    │ ← Items (17px)
│  ─────────────────────────────────  │
│           TỔNG TIỀN:        105k    │ ← Total (24px)
│  ─────────────────────────────────  │
│       Cảm ơn quý khách!             │ ← Thanks (22px, centered)
└─────────────────────────────────────┘
```

**Font sizes** (matching preview.go):
- Shop title: 25px (fake bold với text-shadow)
- Shop address/phone: 16px
- Bill title: 34px
- Order info: 16px
- Table header/items: 17px
- Total: 24px
- Thanks: 22px

**Column positions** (matching preview.go):
- STT: 10px from left
- Tên món: 50px from left
- SL: 290px from left
- Đơn giá: 340px from left
- Thành tiền: right-aligned, 10px from right

### 2. HTML Renderer Service

**File**: `backend/application/services/html_bill_renderer.go`

**Flow**:
```
Order Data → Template Data → HTML → Chromedp → PNG Image → ESC/POS → Base64
```

**Key functions**:
- `NewHTMLBillRenderer()`: Khởi tạo với template path
- `RenderBillToBase64()`: Render bill → base64 ESC/POS
- `RenderBillToESCPOS()`: Render bill → raw ESC/POS bytes
- `renderHTMLToImage()`: HTML → PNG image (chromedp)
- `imageToESCPOSCommands()`: PNG → ESC/POS GS v 0 commands

### 3. Template Data

**Struct**: `BillTemplateData`

```go
type BillTemplateData struct {
    ShopName          string
    ShopAddress       string
    ShopPhone         string
    ShowLogo          bool
    ShowAddress       bool
    ShowPhone         bool
    ShowCustomMessage bool
    CustomMessage     string
    LogoPath          string  // file:// URL
    OrderNumber       string
    WaiterName        string
    PaymentMethod     string
    CreatedDate       string  // formatted
    Items             []BillItem
    Total             float64
}
```

**Template functions**:
- `{{formatMoney .Total}}`: Format số tiền với dấu phẩy
- `{{add $index 1}}`: Tính STT (index + 1)

## Cách hoạt động

### 1. Khi tạo Order

```
Order Created → CreatePrintJobsForOrder()
    ↓
createBillJob()
    ↓
HTMLBillRenderer.RenderBillToBase64()
    ├─ Prepare template data
    ├─ Execute HTML template
    ├─ Render HTML → PNG (chromedp)
    ├─ Convert PNG → ESC/POS
    └─ Encode base64
    ↓
Save PrintJob (ContentType="binary")
    ↓
Print Worker → Decode base64 → Send to printer
    ↓
Bill printed! ✅
```

### 2. Chromedp Rendering

Chromedp là headless Chrome browser:
- Tạo context với timeout 30s
- Navigate to "about:blank"
- Inject HTML content
- Wait 500ms for rendering
- Capture full screenshot (quality 100)
- Return PNG image

## Customize Template

### Thay đổi Layout

Edit `backend/templates/bill_template.html`:

```html
<!-- Thay đổi font size -->
<style>
.bill-title {
    font-size: 40px;  /* Thay vì 34px */
}
</style>

<!-- Thêm element mới -->
<div class="promo-message">
    Giảm giá 10% cho lần mua tiếp theo!
</div>
```

### Thay đổi Colors

```css
body {
    background: #f5f5f5;  /* Màu nền */
    color: #333;          /* Màu chữ */
}

.bill-title {
    color: #0066cc;       /* Màu xanh cho title */
}
```

### Thêm Logo Border

```css
.logo img {
    border: 2px solid #000;
    border-radius: 8px;
}
```

### Thay đổi Table Style

```css
.table-header {
    background: #f0f0f0;
    padding: 5px 0;
}

.item-row {
    border-bottom: 1px dashed #ccc;
}
```

## Testing

### 1. Preview HTML trong Browser

```bash
# Copy template và tạo file test với data mẫu
cp backend/templates/bill_template.html /tmp/test_bill.html

# Edit /tmp/test_bill.html, replace {{.ShopName}} với "Tiệm cà phê Ông Tạ", etc.
# Mở trong browser
open /tmp/test_bill.html
```

### 2. Test với Order thực

```bash
# Tạo order mới qua UI hoặc API
# Bill sẽ tự động in

# Hoặc test qua visual print endpoint
curl -X POST http://localhost:8080/api/manager/visual-print/bill \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "507f1f77bcf86cd799439011",
    "printer_ip": "192.168.1.115"
  }'
```

### 3. Preview PNG

```bash
# Tạo preview PNG
curl http://localhost:8080/api/manager/visual-print/preview/507f1f77bcf86cd799439011 \
  -H "Authorization: Bearer $TOKEN"

# Mở file preview_bill_*.png
```

## Yêu cầu

### Backend

1. **Chromedp**: Headless Chrome
   ```bash
   go get github.com/chromedp/chromedp
   ```

2. **Chrome/Chromium**: Phải cài đặt trên server
   - macOS: Chrome tự động detect
   - Linux: `apt-get install chromium-browser`
   - Docker: Cần base image có Chrome

### Template

- File `backend/templates/bill_template.html` phải tồn tại
- Logo path phải đúng (absolute path hoặc file:// URL)

## Troubleshooting

### HTML không render

1. **Kiểm tra Chrome**:
   ```bash
   which google-chrome
   which chromium-browser
   ```

2. **Kiểm tra logs**:
   ```bash
   grep "chromedp error" backend.log
   ```

3. **Test chromedp**:
   ```go
   ctx, cancel := chromedp.NewContext(context.Background())
   defer cancel()
   
   var buf []byte
   err := chromedp.Run(ctx,
       chromedp.Navigate("https://google.com"),
       chromedp.FullScreenshot(&buf, 100),
   )
   ```

### Template không parse

1. **Kiểm tra syntax**:
   - Go template syntax: `{{.Variable}}`
   - Closing tags: `{{end}}`
   - Functions: `{{formatMoney .Total}}`

2. **Kiểm tra path**:
   ```bash
   ls -la backend/templates/bill_template.html
   ```

### Logo không hiển thị

1. **Kiểm tra path**:
   - Phải là absolute path
   - Hoặc file:// URL
   - VD: `file:///path/to/logo.jpg`

2. **Kiểm tra file**:
   ```bash
   ls -la backend/uploads/logos/
   ```

### Bill in ra bị lỗi layout

1. **Kiểm tra width**: Phải đúng 576px
2. **Kiểm tra margin**: 20px
3. **Kiểm tra font sizes**: Match với preview.go
4. **Test trong browser** trước khi in

## So sánh với Visual Renderer

| Tính năng | Visual Renderer (gg) | HTML Renderer (chromedp) |
|-----------|---------------------|-------------------------|
| Render | Go code với gg library | HTML/CSS template |
| Customize | Phải edit Go code | Edit HTML/CSS |
| Preview | Cần compile & run | Mở HTML trong browser |
| Maintain | Khó, cần hiểu Go | Dễ, chỉ cần HTML/CSS |
| Performance | Nhanh hơn | Chậm hơn (headless Chrome) |
| Dependencies | gg, resize | chromedp, Chrome |
| Flexibility | Hạn chế | Rất linh hoạt |

## Kế hoạch phát triển

- [ ] Cache rendered images để tăng performance
- [ ] Hỗ trợ nhiều templates (template A, B, C...)
- [ ] Template editor trong UI
- [ ] Preview real-time khi edit template
- [ ] Export template thành PDF
- [ ] QR code trên bill
- [ ] Barcode cho order number

## Kết luận

HTML template approach cho phép customize dễ dàng và maintain tốt hơn. Layout chính xác 100% với preview.go nhờ CSS pixel-perfect. Hệ thống tự động in bill khi tạo order với template này.
