# Hướng dẫn sử dụng Template mặc định

## Template đã tạo

Đã tạo thành công template hóa đơn mặc định với các thông tin:

### Thông tin Template

- **Tên**: Hóa đơn mặc định
- **Loại**: BILL (Hóa đơn bán hàng)
- **Trạng thái**: Mặc định (is_default = true)

### Nội dung Template

Template bao gồm các thông tin:

1. **Header (Đầu hóa đơn)**
   - Tên cửa hàng
   - Địa chỉ (nếu có)
   - Số điện thoại (nếu có)

2. **Thông tin đơn hàng**
   - Số order (Order Number)
   - Ngày giờ tạo đơn (định dạng: DD/MM/YYYY HH:mm)
   - Số bàn (nếu có)
   - Tên khách hàng (nếu có)

3. **Chi tiết món**
   - Tên món
   - Variant (nếu có)
   - Số lượng x Đơn giá = Thành tiền

4. **Tổng tiền**
   - Tổng tiền hàng (Subtotal)
   - Giảm giá (nếu có)
   - Thuế (nếu có)
   - **TỔNG CỘNG** (Total)

5. **Footer (Cuối hóa đơn)**
   - Lời cảm ơn
   - Tin nhắn tùy chỉnh (nếu có)

## Cách sử dụng

### 1. Xem template từ Frontend

Truy cập: http://localhost:5173/#/print-management

Trong phần "Print Templates", bạn sẽ thấy template "Hóa đơn mặc định"

### 2. Chỉnh sửa template

Bạn có thể chỉnh sửa template bằng cách:
- Click vào template trong danh sách
- Sửa nội dung theo ý muốn
- Lưu lại

### 3. Test in template

Khi tạo order mới, hệ thống sẽ tự động sử dụng template mặc định này để in hóa đơn.

## Template Variables

Template sử dụng Go template syntax với các biến sau:

### Shop Information
- `{{.ShopName}}` - Tên cửa hàng
- `{{.ShopAddress}}` - Địa chỉ
- `{{.ShopPhone}}` - Số điện thoại
- `{{.ShowAddress}}` - Hiển thị địa chỉ (true/false)
- `{{.ShowPhone}}` - Hiển thị SĐT (true/false)
- `{{.ShowCustomMessage}}` - Hiển thị tin nhắn tùy chỉnh (true/false)
- `{{.CustomMessage}}` - Tin nhắn tùy chỉnh

### Order Information
- `{{.Order.OrderNumber}}` - Số order
- `{{.Order.CreatedAt}}` - Thời gian tạo
- `{{.Order.TableNumber}}` - Số bàn
- `{{.Order.CustomerName}}` - Tên khách
- `{{.Order.Subtotal}}` - Tổng tiền hàng
- `{{.Order.DiscountAmount}}` - Số tiền giảm giá
- `{{.Order.TaxAmount}}` - Tiền thuế
- `{{.Order.Total}}` - Tổng cộng

### Order Items (Loop)
```
{{range .Order.Items}}
{{.Name}} - {{.VariantName}}
  {{.Quantity}} x {{formatPrice .UnitPrice}} = {{formatPrice .TotalPrice}}
{{end}}
```

### Helper Functions
- `{{formatPrice .Order.Total}}` - Format số tiền (thêm dấu phẩy)
- `{{.Order.CreatedAt.Format "02/01/2006 15:04"}}` - Format ngày giờ

### Conditional Display
```
{{if .ShowAddress}}{{.ShopAddress}}{{end}}
{{if gt .Order.DiscountAmount 0.0}}Giảm giá: ...{{end}}
```

## Ví dụ Output

```
Cafe ABC
123 Nguyen Hue, Q1, TP.HCM
Tel: 028-1234-5678
================================
HOA DON BAN HANG
Order: ORD-001
Ngay: 22/02/2026 09:15
Ban: 5
Khach: Nguyen Van A
================================
Cafe Latte
  2 x 45,000 = 90,000
Banh Mi Thit
  1 x 35,000 = 35,000
================================
Tong tien: 125,000 VND
Giam gia: -10,000 VND
TONG CONG: 115,000 VND
================================
Cam on quy khach!
Hen gap lai!
```

## Lưu ý

1. **Font Size**: Hiện tại font size đã được tăng lên 18pt để dễ đọc hơn

2. **Tiếng Việt**: Template hỗ trợ đầy đủ tiếng Việt với tất cả dấu thanh

3. **Separator Lines**: Sử dụng `===` hoặc `---` để tạo đường kẻ ngang

4. **Alignment**: 
   - Text ngắn (< 30 ký tự) không có dấu `:` sẽ tự động căn giữa
   - Text có chứa `TOTAL`, `TỔNG` sẽ tự động căn giữa và in đậm
   - Text thông thường sẽ căn trái

5. **Bold Text**: Text có chứa `TOTAL`, `TỔNG`, `GIẢM GIÁ`, `DISCOUNT` sẽ tự động in đậm

## Tạo lại Template

Nếu cần tạo lại template mặc định, chạy:

```bash
cd backend
MONGODB_URI="mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin" \
MONGODB_DATABASE="cafe_pos" \
go run cmd/create-default-template/main.go
```

## Troubleshooting

### Template không hiển thị
- Kiểm tra MongoDB đang chạy
- Kiểm tra backend đang chạy
- Refresh trang frontend

### In không ra tiếng Việt
- Đảm bảo backend đã restart sau khi update code
- Kiểm tra font đã được cài đặt (Arial Unicode MS, Roboto, hoặc DejaVu Sans)

### Font quá nhỏ/lớn
- Chỉnh sửa `FontSize` trong `backend/infrastructure/printing/escpos_printer.go`
- Hiện tại đang dùng 18pt (có thể tăng lên 20pt hoặc 22pt nếu cần)
