# Quick Start - Manager Features

## 🎯 Tính Năng Mới Cho Manager

### 1. Quản Lý Thiết Bị (Facility Management)
**Đường dẫn**: `/facilities`  
**Icon**: 🏢 Cơ sở vật chất

**Chức năng**:
- 📊 Dashboard thống kê thiết bị
- ➕ Thêm thiết bị mới
- ✏️ Sửa thông tin thiết bị
- 🗑️ Xóa thiết bị
- 📅 Xem lịch bảo trì
- ⚠️ Xem báo cáo sự cố
- 🔍 Tìm kiếm thiết bị

### 2. Quản Lý Nguyên Liệu (Ingredient Management)
**Đường dẫn**: `/ingredients`  
**Icon**: 🥬 Nguyên liệu

**Chức năng**:
- 📊 Dashboard thống kê tồn kho
- ➕ Thêm nguyên liệu mới
- ✏️ Sửa thông tin nguyên liệu
- 🗑️ Xóa nguyên liệu
- 📦 Điều chỉnh tồn kho (Nhập/Xuất)
- 📊 Xem lịch sử tồn kho
- ⚠️ Cảnh báo sắp hết hàng
- 🔍 Tìm kiếm nguyên liệu

## 🚀 Cách Sử Dụng

### Bước 1: Đăng Nhập
```
1. Mở trình duyệt
2. Truy cập: http://localhost:5173
3. Đăng nhập với tài khoản Manager
```

### Bước 2: Truy Cập Tính Năng
```
1. Sau khi đăng nhập, xem Dashboard
2. Click vào "🏢 Cơ sở vật chất" hoặc "🥬 Nguyên liệu"
3. Bắt đầu quản lý!
```

## 📋 Ví Dụ Sử Dụng

### Thêm Thiết Bị Mới
```
1. Vào /facilities
2. Click "➕ Thêm Thiết Bị"
3. Điền thông tin:
   - Tên: Máy pha cà phê
   - Loại: Equipment
   - Vị trí: Quầy bar
   - Trạng thái: Hoạt động
4. Click "Thêm Mới"
```

### Điều Chỉnh Tồn Kho
```
1. Vào /ingredients
2. Tìm nguyên liệu cần điều chỉnh
3. Click "📦 Điều Chỉnh"
4. Chọn loại:
   - Nhập Hàng: Khi nhập thêm
   - Xuất Hàng: Khi sử dụng
   - Điều Chỉnh: Khi kiểm kê
5. Nhập số lượng và lý do
6. Click "Xác Nhận"
```

## 🎨 Màu Sắc Trạng Thái

### Thiết Bị
- 🟢 Hoạt Động (Operational)
- 🟡 Bảo Trì (Maintenance)
- 🔴 Hỏng Hóc (Broken)
- ⚫ Ngừng Sử Dụng (Retired)

### Nguyên Liệu
- 🟢 Đủ Hàng (In Stock)
- 🟡 Sắp Hết (Low Stock)
- 🔴 Hết Hàng (Out of Stock)

## ✅ Checklist Kiểm Tra

- [ ] Backend đang chạy (port 8080)
- [ ] Frontend đang chạy (port 5173)
- [ ] Đã đăng nhập với tài khoản Manager
- [ ] Thấy 2 menu mới trong navigation
- [ ] Click vào menu và views load thành công
- [ ] Có thể thêm/sửa/xóa dữ liệu

## 🐛 Troubleshooting

### Views không load?
```bash
# 1. Check backend
curl http://localhost:8080/api/manager/facilities

# 2. Check console (F12)
# Xem có error không

# 3. Check authentication
# Đảm bảo đã login với role Manager
```

### Không thấy menu?
```
- Kiểm tra role: Phải là Manager
- Refresh trang (F5)
- Clear cache và login lại
```

### API errors?
```
- Kiểm tra backend đang chạy
- Kiểm tra token còn hạn không
- Xem Network tab trong DevTools
```

## 📞 Support

Nếu gặp vấn đề:
1. Check console errors (F12)
2. Check network requests
3. Verify backend is running
4. Verify logged in as Manager

## 🎉 Hoàn Thành!

Manager giờ có thể:
- ✅ Quản lý thiết bị đầy đủ
- ✅ Quản lý nguyên liệu đầy đủ
- ✅ Theo dõi tồn kho
- ✅ Lên lịch bảo trì
- ✅ Xử lý sự cố

Chúc sử dụng vui vẻ! 🚀

