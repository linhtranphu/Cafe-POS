# Hướng Dẫn In Logo Trên Hóa Đơn

## Tổng Quan
Hệ thống đã được tích hợp đầy đủ chức năng in logo trên hóa đơn. Logo sẽ được chuyển đổi sang định dạng ESC/POS và in ra máy in nhiệt.

## Các Bước Cài Đặt

### 1. Upload Logo
1. Truy cập: http://localhost:5173/#/print-management
2. Tìm phần "Logo Cửa Hàng"
3. Click "Chọn file" và chọn logo (PNG, JPG, JPEG)
4. Logo tối đa 2MB
5. Click "Upload" để lưu

### 2. Tạo Template Có Logo
1. Vẫn ở trang Print Management
2. Click "Tạo Template Mới"
3. Copy nội dung từ file `BILL_TEMPLATE_WITH_LOGO.txt`
4. Paste vào ô "Nội dung template"
5. Đặt tên template (ví dụ: "Hóa đơn có logo")
6. Chọn loại: "BILL"
7. Tick vào "Đặt làm mặc định"
8. Click "Lưu"

### 3. Kiểm Tra Template
Template phải có marker `[LOGO]` ở đầu:
```
{{if .ShowLogo}}[LOGO]{{end}}
{{.ShopName}}
...
```

### 4. Test In Logo
1. Tạo order mới từ giao diện
2. Hoàn thành order
3. Hệ thống sẽ tự động in hóa đơn với logo

## Cách Hoạt Động

### Quy Trình Xử Lý Logo
1. **Template Rendering**: Khi render template, marker `[LOGO]` được giữ lại
2. **Logo Loading**: LogoRenderer load file logo từ `./uploads/logos/`
3. **Resize**: Logo được resize về max 25% chiều rộng giấy (144px cho giấy 80mm)
4. **Convert to Grayscale**: Chuyển sang ảnh xám
5. **Convert to ESC/POS**: ImageConverter chuyển sang lệnh ESC/POS (GS v 0)
6. **Replace Marker**: Marker `[LOGO]` được thay bằng binary ESC/POS commands
7. **Print**: Gửi toàn bộ nội dung (logo + text) đến máy in

### Các Module Đã Implement
- ✅ `LogoRenderer`: Load và resize logo
- ✅ `ImageConverter`: Convert ảnh sang ESC/POS
- ✅ `TableFormatter`: Format bảng món ăn
- ✅ `FontSizeManager`: Quản lý font size
- ✅ `ImageCompositor`: Kết hợp logo và text
- ✅ Template integration trong `template_renderer.go`

## Troubleshooting

### Logo không hiển thị
1. Kiểm tra file logo tồn tại:
   ```bash
   ls -lh ./uploads/logos/
   ```

2. Kiểm tra shop settings:
   ```bash
   curl http://localhost:3000/api/settings | jq '.show_logo, .logo_url'
   ```

3. Kiểm tra template có marker `[LOGO]`:
   ```bash
   curl http://localhost:3000/api/manager/print-templates?type=BILL | jq '.templates[].content' | grep LOGO
   ```

4. Xem backend logs:
   ```bash
   tail -f backend.log | grep -i logo
   ```

### Preview không hoạt động
- Đã fix: Frontend giờ tự động fetch template content và type trước khi gọi preview API
- Preview API yêu cầu `content` và `type` trong request body

### Logo quá lớn/nhỏ
- Logo tự động resize về max 25% chiều rộng giấy
- Với giấy 80mm (576px): logo max 144px
- Điều chỉnh trong `LogoRenderer` nếu cần thay đổi tỷ lệ

## Test Script
Chạy script test để kiểm tra toàn bộ:
```bash
./test-logo-rendering.sh
```

Script sẽ kiểm tra:
- Backend đang chạy
- Logo file tồn tại
- Shop settings đúng
- Template có marker [LOGO]
- Preview API hoạt động

## Kết Quả Mong Đợi
Khi in hóa đơn, bạn sẽ thấy:
1. Logo ở đầu hóa đơn (nếu có)
2. Thông tin cửa hàng
3. Bảng món ăn được format đẹp
4. Tổng tiền và footer

## Lưu Ý Kỹ Thuật
- Logo được convert sang binary ESC/POS commands
- Commands này được embed trực tiếp vào content string
- Máy in nhiệt sẽ nhận diện và in logo
- Format: GS v 0 (ESC/POS raster image command)
