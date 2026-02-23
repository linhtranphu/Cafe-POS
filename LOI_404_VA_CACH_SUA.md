# Lỗi 404 API và Cách Sửa

## ❌ Vấn Đề

Frontend đang gặp lỗi 404 khi gọi API:
```
GET http://localhost:5173/api/cashier-shifts/.../managed-funds 404 (Not Found)
```

## 🔍 Nguyên Nhân

File cấu hình Vite (`frontend/vite.config.js`) đang trỏ sai cổng backend:
- **Đã cấu hình**: `target: 'http://localhost:3000'`
- **Backend thực tế**: Chạy trên `http://localhost:8080`

## ✅ Đã Sửa

Đã cập nhật `frontend/vite.config.js`:
```javascript
server: {
  port: 5173,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',  // Đổi từ 3000 sang 8080
      changeOrigin: true,
      rewrite: (path) => path
    }
  }
}
```

## ⚠️ BẠN CẦN LÀM GÌ

**KHỞI ĐỘNG LẠI FRONTEND DEV SERVER**

```bash
# 1. Dừng server hiện tại
# Nhấn Ctrl+C trong terminal đang chạy npm run dev

# 2. Khởi động lại
cd frontend
npm run dev
```

**Tại sao?** Vite chỉ đọc file cấu hình khi khởi động. Thay đổi trong `vite.config.js` cần khởi động lại mới có hiệu lực.

## 🧪 Kiểm Tra Sau Khi Sửa

1. **Mở Console trình duyệt** (F12)
2. **Vào Cashier Dashboard**
3. **Kiểm tra tab Network**:
   - Trước: `GET http://localhost:5173/api/... → 404`
   - Sau: `GET http://localhost:5173/api/... → 200 OK`

4. **Xác nhận hiển thị**:
   - Thấy phần "💰 Tiền đang quản lý"
   - Hiển thị số tiền mặt và chuyển khoản
   - Không có thông báo lỗi

## 🚀 Hướng Dẫn Nhanh

### Cách 1: Khởi động lại thủ công

```bash
# Terminal 1: Backend
cd backend
go run main.go

# Terminal 2: Frontend (KHỞI ĐỘNG LẠI)
cd frontend
npm run dev
```

### Cách 2: Dùng script tự động

```bash
./RESTART_SERVERS.sh
```

## ✅ Checklist

Sau khi khởi động lại:

- [ ] Backend chạy trên port 8080
- [ ] Frontend chạy trên port 5173
- [ ] Mở http://localhost:5173
- [ ] Đăng nhập với tài khoản cashier
- [ ] Vào Cashier Dashboard
- [ ] Thấy phần "Tiền đang quản lý"
- [ ] Không có lỗi 404 trong console
- [ ] Dữ liệu hiển thị đúng

## 🐛 Nếu Vẫn Lỗi

### Kiểm tra Backend

```bash
# Xem backend có chạy không
curl http://localhost:8080/api/health

# Nếu lỗi "connection refused", khởi động backend:
cd backend
go run main.go
```

### Kiểm tra Cổng

```bash
# Xem process nào đang dùng cổng
lsof -i :8080  # Backend
lsof -i :5173  # Frontend
```

### Xóa Cache Trình Duyệt

1. Mở DevTools (F12)
2. Chuột phải vào nút refresh
3. Chọn "Empty Cache and Hard Reload"

### Kiểm tra Token

```javascript
// Trong console trình duyệt
localStorage.getItem('token')

// Nếu null hoặc hết hạn, đăng nhập lại
```

## 📚 Tài Liệu Chi Tiết

- `API_PROXY_FIX.md` - Chi tiết về fix
- `TROUBLESHOOTING_404_ERRORS.md` - Hướng dẫn xử lý lỗi
- `API_404_FIX_SUMMARY.md` - Tóm tắt (tiếng Anh)

## 📝 Tóm Tắt

1. ✅ Đã sửa file `vite.config.js`
2. ⏳ **BẠN CẦN**: Khởi động lại frontend dev server
3. ⏳ **SAU ĐÓ**: Test lại trên trình duyệt

---

**QUAN TRỌNG**: Fix đã hoàn tất, nhưng bạn PHẢI khởi động lại frontend dev server để có hiệu lực!

```bash
# Làm ngay bây giờ:
cd frontend
# Nhấn Ctrl+C để dừng
npm run dev  # Khởi động lại
```
