# Hướng dẫn In Bill với Visual Template

## Tổng quan

Hệ thống in bill mới sử dụng template visual chính xác như trong `preview.go`, tạo hình ảnh và in qua ESC/POS commands.

## Đặc điểm

- **Kích thước chính xác**: 576px width (72mm @ 203 DPI) cho máy in Zywell ZY303
- **Layout giống hệt preview.go**:
  - Logo 200px bên trái
  - Tên quán, địa chỉ, SĐT bên phải
  - Tiêu đề "HÓA ĐƠN THANH TOÁN" căn giữa
  - Thông tin order: Order number, Waiter, Payment method, Date
  - Bảng items với các cột: STT, Tên món, SL, Đơn giá, Thành tiền
  - Tổng tiền và lời cảm ơn
- **Font sizes chính xác**: 25, 16, 34, 17, 24, 22 (giống preview.go)
- **Margin**: 20px
- **In qua ESC/POS**: Sử dụng GS v 0 command cho raster bit image

## Cách sử dụng

### 1. Truy cập trang In Visual Bill

Vào `http://localhost:5173/#/print-management` và chọn tab **"In Visual Bill"**

### 2. Cài đặt máy in

- Nhập IP máy in (VD: `192.168.1.115`)
- Port mặc định: 9100

### 3. Chọn Order

- Tìm kiếm order theo số order hoặc tên khách
- Click vào order để chọn
- Xem preview thông tin order

### 4. In Bill

- Click nút **"🖨️ In Bill"** để in trực tiếp
- Hoặc click **"👁️ Preview"** để tạo file PNG preview

## API Endpoints

### POST /api/manager/visual-print/bill

In bill với visual template

**Request:**
```json
{
  "order_id": "507f1f77bcf86cd799439011",
  "printer_ip": "192.168.1.115"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Bill printed successfully",
  "order_number": "20260222-095703-168"
}
```

### GET /api/manager/visual-print/preview/:order_id

Tạo preview PNG của bill

**Response:**
```json
{
  "success": true,
  "message": "Preview created successfully",
  "filename": "preview_bill_20260222-095703-168.png",
  "order_number": "20260222-095703-168"
}
```

## Cấu trúc Code

### Backend

- **`visual_bill_renderer.go`**: Service render bill thành hình ảnh và ESC/POS
  - `RenderBillToESCPOS()`: Tạo ESC/POS commands
  - `createBillImage()`: Tạo hình ảnh bill (giống preview.go)
  - `imageToESCPOS()`: Convert image sang ESC/POS raster bit image
  - `SaveImagePreview()`: Lưu preview PNG

- **`visual_print_handler.go`**: HTTP handler
  - `PrintVisualBill()`: In bill
  - `PreviewVisualBill()`: Tạo preview

### Frontend

- **`VisualBillPrinter.vue`**: Component chính để chọn order và in
  - Tìm kiếm orders
  - Chọn order
  - In bill hoặc tạo preview

- **`BillTemplateVisual.vue`**: Component preview template (canvas-based)
  - Hiển thị preview 576px width
  - Chỉnh sửa template settings
  - Test với dữ liệu mẫu

## Yêu cầu

### Backend
- Go 1.21+
- Libraries:
  - `github.com/fogleman/gg` (image rendering)
  - `github.com/nfnt/resize` (image resizing)

### Frontend
- Vue 3
- Vite

### Máy in
- Zywell ZY303 hoặc tương thích
- Hỗ trợ ESC/POS
- Kết nối mạng (port 9100)
- Giấy in 72mm

## Troubleshooting

### Bill không in ra

1. Kiểm tra IP máy in đúng chưa
2. Kiểm tra máy in đã bật và kết nối mạng
3. Kiểm tra port 9100 có mở không
4. Xem logs backend để biết lỗi chi tiết

### Hình ảnh bị cắt hoặc lỗi

1. Kiểm tra logo có tồn tại không (trong `backend/uploads/logos/`)
2. Kiểm tra font system có cài đặt không
3. Xem logs để biết lỗi render

### Preview không tạo được

1. Kiểm tra quyền ghi file trong thư mục backend
2. Kiểm tra order ID có đúng không
3. Xem logs backend

## So sánh với hệ thống cũ

| Tính năng | Hệ thống cũ (Text Template) | Hệ thống mới (Visual Template) |
|-----------|----------------------------|-------------------------------|
| Render | Go template text | Image-based (gg library) |
| Layout | Text-based, khó control | Pixel-perfect, giống preview.go |
| Font size | Không chính xác | Chính xác (25, 16, 34, 17, 24, 22) |
| Logo | Không hỗ trợ tốt | Hỗ trợ đầy đủ, resize 200px |
| Kích thước | 80mm (640px) | 72mm (576px) - đúng với Zywell ZY303 |
| ESC/POS | Text commands | Raster bit image (GS v 0) |

## Kế hoạch phát triển

- [ ] Lưu template settings vào database
- [ ] Cho phép customize font sizes, margins
- [ ] Hỗ trợ nhiều loại máy in
- [ ] In batch nhiều bills cùng lúc
- [ ] Export PDF
- [ ] Email bill cho khách hàng

## Liên hệ

Nếu có vấn đề, vui lòng kiểm tra logs hoặc liên hệ team phát triển.
