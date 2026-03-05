# Hướng Dẫn Test Logo Upload

## Tổng quan
Feature logo upload cho phép upload logo quán và hiển thị trên bill. Logo sẽ xuất hiện ở góc trên bên trái của hóa đơn.

## Cài đặt Backend

### 1. Tạo thư mục uploads
```bash
mkdir -p uploads/logos
```

### 2. Khởi động backend
```bash
cd backend
go run main.go
```

Backend sẽ chạy ở `http://localhost:8080`

## Cài đặt Frontend

### 1. Khởi động frontend
```bash
cd frontend
npm run dev
```

Frontend sẽ chạy ở `http://localhost:5173`

## Test Cases

### Test 1: Upload Logo PNG
1. Mở trình duyệt: `http://localhost:5173/#/print-management`
2. Click tab "Cài Đặt" (⚙️)
3. Tìm section "Logo Quán"
4. Click button "📤 Upload logo"
5. Chọn file PNG (< 2MB)
6. Verify:
   - Preview hiển thị ngay sau khi chọn
   - Thông báo "Upload logo thành công"
   - Logo được lưu vào `uploads/logos/`
   - Logo URL được cập nhật trong database

### Test 2: Upload Logo JPG/JPEG
1. Lặp lại Test 1 với file JPG hoặc JPEG
2. Verify kết quả tương tự

### Test 3: Upload File Quá Lớn (> 2MB)
1. Chọn file > 2MB
2. Verify:
   - Hiển thị lỗi: "Kích thước file tối đa 2MB"
   - File không được upload
   - Logo cũ không bị thay đổi

### Test 4: Upload File Sai Format
1. Chọn file không phải PNG/JPG/JPEG (ví dụ: PDF, GIF)
2. Verify:
   - Hiển thị lỗi: "Chỉ hỗ trợ định dạng PNG, JPG, JPEG"
   - File không được upload

### Test 5: Thay Đổi Logo
1. Upload logo lần đầu
2. Upload logo mới
3. Verify:
   - Logo cũ bị xóa khỏi `uploads/logos/`
   - Logo mới được hiển thị
   - Chỉ có 1 logo trong thư mục

### Test 6: Xóa Logo
1. Upload logo
2. Click button "X" (đỏ) ở góc logo preview
3. Confirm xóa
4. Verify:
   - Logo bị xóa khỏi `uploads/logos/`
   - Preview biến mất
   - Logo URL trong database = ""
   - Button đổi thành "📤 Upload logo"

### Test 7: Toggle "Hiển thị logo trên bill"
1. Upload logo
2. Uncheck "Hiển thị logo trên bill"
3. Click "💾 Lưu cài đặt"
4. Tạo order và in bill
5. Verify: Bill không có logo
6. Check lại "Hiển thị logo trên bill"
7. Click "💾 Lưu cài đặt"
8. Tạo order và in bill
9. Verify: Bill có logo ở góc trên bên trái

### Test 8: Logo Sizing
1. Upload logo rất lớn (ví dụ: 2000x2000px)
2. Tạo order và in bill
3. Verify:
   - Logo được resize xuống <= 25% paper width
   - Logo không bị méo
   - Logo vẫn rõ nét

### Test 9: Logo với 58mm Paper
1. Vào Settings
2. Chọn "Khổ Giấy Bill" = "58mm"
3. Upload logo
4. Tạo order và in bill
5. Verify:
   - Logo fit trong 58mm paper
   - Logo không bị cắt

### Test 10: Logo với 80mm Paper
1. Vào Settings
2. Chọn "Khổ Giấy Bill" = "80mm"
3. Upload logo
4. Tạo order và in bill
5. Verify:
   - Logo fit trong 80mm paper
   - Logo có kích thước phù hợp

## API Endpoints

### Upload Logo
```bash
curl -X POST http://localhost:8080/api/settings/logo \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "logo=@/path/to/logo.png"
```

Response:
```json
{
  "logo_url": "/uploads/logos/logo_12345.png",
  "message": "Upload logo thành công"
}
```

### Delete Logo
```bash
curl -X DELETE http://localhost:8080/api/settings/logo \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Response:
```json
{
  "message": "Đã xóa logo thành công"
}
```

### Get Logo (Static File)
```
http://localhost:8080/uploads/logos/logo_12345.png
```

## Troubleshooting

### Lỗi: "No file uploaded"
- Kiểm tra form data có field "logo" không
- Kiểm tra Content-Type header = "multipart/form-data"

### Lỗi: "Không thể lưu file"
- Kiểm tra thư mục `uploads/logos/` có tồn tại không
- Kiểm tra quyền write cho thư mục

### Lỗi: "Không thể cập nhật cài đặt"
- Kiểm tra MongoDB connection
- Kiểm tra shop_settings collection có data không

### Logo không hiển thị trên bill
- Kiểm tra `show_logo` = true trong shop_settings
- Kiểm tra logo file tồn tại trong `uploads/logos/`
- Kiểm tra template có marker `[LOGO]` không
- Kiểm tra LogoRenderer được khởi tạo trong TemplateRenderer

### Logo bị méo
- Kiểm tra LogoRenderer resize logic
- Kiểm tra aspect ratio được giữ nguyên

## Files Liên Quan

### Backend
- `backend/interfaces/http/logo_upload_handler.go` - Logo upload handler
- `backend/main.go` - Route registration
- `backend/infrastructure/printing/logo_renderer.go` - Logo rendering
- `backend/application/services/template_renderer.go` - Template integration

### Frontend
- `frontend/src/components/printing/ShopSettingsForm.vue` - Logo upload UI
- `frontend/src/views/PrintManagementView.vue` - Print management page

### Storage
- `uploads/logos/` - Logo storage directory

## Next Steps

Sau khi test xong logo upload, tiếp tục với:
1. Test template management (switch giữa templates)
2. Test end-to-end bill rendering với logo + table
3. Test với máy in thật
