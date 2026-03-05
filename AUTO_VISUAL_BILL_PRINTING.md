# Tự động In Bill với Visual Template

## Tổng quan

Hệ thống đã được cập nhật để **tự động in bill với visual template** khi tạo order mới. Bill sẽ được render chính xác như trong `preview.go` với layout pixel-perfect.

## Cách hoạt động

### 1. Khi tạo Order

Khi một order mới được tạo (status COMPLETED), hệ thống tự động:

1. **Tạo Print Job**: Gọi `CreatePrintJobsForOrder()` để tạo:
   - 1 bill print job (sử dụng visual template)
   - N label print jobs (cho mỗi item)

2. **Render Visual Bill**: 
   - Sử dụng `VisualBillRenderer` để tạo hình ảnh bill
   - Layout chính xác như `preview.go`:
     - 576px width (72mm @ 203 DPI)
     - Logo 200px bên trái
     - Shop info bên phải
     - Order details, items table, total
   - Convert sang ESC/POS commands (GS v 0 raster bit image)
   - Encode thành base64 và lưu vào `PrintJob.Content`

3. **Print Worker xử lý**:
   - Background worker poll pending jobs mỗi 10 giây
   - Detect base64 content và decode
   - Gửi raw ESC/POS commands đến máy in
   - Retry tự động nếu thất bại (max 3 lần)

### 2. Flow chi tiết

```
Order Created (COMPLETED)
    ↓
CreatePrintJobsForOrder()
    ↓
createBillJob()
    ↓
VisualBillRenderer.RenderBillToBase64()
    ├─ createBillImage() - Tạo hình ảnh với gg library
    ├─ imageToESCPOS() - Convert sang ESC/POS
    └─ base64.Encode() - Encode thành string
    ↓
Save PrintJob (Content = base64, ContentType = "binary")
    ↓
Print Worker picks up job
    ↓
ESCPOSPrinter.Print()
    ├─ Detect base64 content
    ├─ Decode base64 → raw ESC/POS bytes
    └─ Send to printer via TCP
    ↓
Bill printed! ✅
```

## Cấu hình

### 1. Cài đặt Printer

Vào **Print Management** → **Máy In**:

1. Thêm máy in mới:
   - Type: **Bill**
   - Name: "Máy in Bill"
   - Connection Type: **Network**
   - IP Address: `192.168.1.115`
   - Port: `9100`
   - Enabled: ✅

2. Đặt làm **Default Printer** cho Bill

### 2. Cài đặt Shop Settings

Vào **Print Management** → **Cài Đặt**:

1. Upload logo (sẽ được resize thành 200px width)
2. Nhập thông tin:
   - Shop Name: "Tiệm cà phê Ông Tạ"
   - Address: "Đ/c: 10/8 Trần Nhật Duật, P. Tân Định, Quận 1, HCM"
   - Phone: "0906990602"
   - Custom Message: "Cảm ơn quý khách!"

3. Bật các options:
   - ✅ Show Logo
   - ✅ Show Address
   - ✅ Show Phone
   - ✅ Show Custom Message

## Thay đổi Code

### Backend

#### 1. Domain Model
- **`print_job.go`**: Thêm field `ContentType` ("text" hoặc "binary")

#### 2. Services
- **`visual_bill_renderer.go`**: Service mới để render bill thành hình ảnh và ESC/POS
  - `RenderBillToBase64()`: Render và encode base64
  - `createBillImage()`: Tạo hình ảnh (giống preview.go)
  - `imageToESCPOS()`: Convert sang ESC/POS GS v 0 commands

- **`print_service.go`**: Cập nhật để sử dụng visual renderer
  - Thêm `visualRenderer` field
  - `createBillJob()` ưu tiên dùng visual renderer, fallback về text template nếu lỗi

- **`print_worker.go`**: Xử lý binary content
  - Detect và log content type

#### 3. Infrastructure
- **`escpos_printer.go`**: Xử lý base64 content
  - `Print()`: Detect base64, decode và gửi raw bytes
  - `isBase64()`: Helper function để detect base64 string

#### 4. Main
- **`main.go`**: Khởi tạo visual renderer và inject vào print service

### Frontend

Không cần thay đổi! Frontend vẫn tạo order như bình thường, backend tự động in.

## Testing

### 1. Test tạo Order

```bash
# Tạo order mới qua API hoặc UI
# Bill sẽ tự động được in
```

### 2. Kiểm tra Print Jobs

Vào **Print Management** → **Print Jobs** để xem:
- Status: PENDING → PRINTING → COMPLETED
- Content Type: "binary"
- Retry count nếu có lỗi

### 3. Kiểm tra Logs

```bash
# Backend logs
tail -f backend.log | grep PRINT

# Logs quan trọng:
# [PRINT] Using visual template for bill
# [PRINTER] Using pre-rendered binary content
# [PRINT SUCCESS] Print completed
```

## Fallback

Nếu visual renderer không khả dụng (lỗi font, library, etc.):
- Hệ thống tự động fallback về **text template**
- Log warning: "Using text template fallback"
- Bill vẫn được in nhưng dùng text-based rendering

## So sánh Text vs Visual Template

| Tính năng | Text Template | Visual Template |
|-----------|---------------|-----------------|
| Render | Go template text | Image-based (gg library) |
| Layout | Text-based, khó control | Pixel-perfect, giống preview.go |
| Font size | Không chính xác | Chính xác (25, 16, 34, 17, 24, 22) |
| Logo | Không hỗ trợ tốt | Hỗ trợ đầy đủ, resize 200px |
| Kích thước | 80mm (640px) | 72mm (576px) - đúng Zywell ZY303 |
| ESC/POS | Text commands | Raster bit image (GS v 0) |
| Content Type | "text" | "binary" (base64) |

## Troubleshooting

### Bill không in ra

1. **Kiểm tra Print Jobs**:
   - Vào Print Management → Print Jobs
   - Xem status và error message

2. **Kiểm tra Printer Config**:
   - Printer có enabled không?
   - IP và port đúng không?
   - Có phải default printer không?

3. **Kiểm tra Logs**:
   ```bash
   grep "PRINT ERROR" backend.log
   ```

4. **Test kết nối**:
   - Vào Print Management → Máy In
   - Click "Test Connection"

### Visual rendering lỗi

1. **Kiểm tra font**:
   - macOS: `/System/Library/Fonts/Supplemental/Arial.ttf`
   - Windows: `C:\Windows\Fonts\arial.ttf`
   - Linux: `/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf`

2. **Kiểm tra logo**:
   - Logo có tồn tại trong `backend/uploads/logos/`?
   - Format: JPEG hoặc PNG

3. **Xem logs**:
   ```bash
   grep "Visual rendering failed" backend.log
   ```

### Bill in ra nhưng không đúng layout

1. **Kiểm tra máy in**:
   - Có phải Zywell ZY303 hoặc tương thích?
   - Giấy in 72mm?

2. **Kiểm tra shop settings**:
   - Logo có quá lớn không?
   - Address quá dài?

## Lợi ích

✅ **Tự động**: Không cần thao tác thủ công  
✅ **Chính xác**: Layout pixel-perfect như preview.go  
✅ **Đáng tin cậy**: Retry tự động, fallback text template  
✅ **Dễ maintain**: Code tập trung, dễ debug  
✅ **Scalable**: Background worker xử lý queue  

## Kế hoạch phát triển

- [ ] Cho phép chọn text hoặc visual template per printer
- [ ] Customize font sizes, margins qua UI
- [ ] Hỗ trợ nhiều loại máy in
- [ ] Preview bill trước khi in
- [ ] Export PDF
- [ ] Email bill cho khách hàng
