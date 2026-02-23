# Khởi Động Lại Frontend

## ✅ Đã Cập Nhật

Đã đổi lại cấu hình Vite về cổng 3000 (đúng với backend của bạn):

```javascript
// frontend/vite.config.js
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:3000',  // ✅ Đúng với backend
      changeOrigin: true
    }
  }
}
```

## ⚠️ BẠN CẦN LÀM NGAY

**KHỞI ĐỘNG LẠI FRONTEND DEV SERVER**

```bash
# 1. Vào thư mục frontend
cd frontend

# 2. Dừng server hiện tại (nếu đang chạy)
# Nhấn Ctrl+C trong terminal đang chạy npm run dev

# 3. Khởi động lại
npm run dev
```

## 🧪 Kiểm Tra

Sau khi khởi động lại:

1. Mở trình duyệt: `http://localhost:5173`
2. Đăng nhập với tài khoản cashier
3. Vào Cashier Dashboard
4. Kiểm tra console (F12):
   - Không còn lỗi 404
   - API calls thành công (200 OK)
5. Xem phần "💰 Tiền đang quản lý" hiển thị đúng

## 📊 Xác Nhận Backend

Đảm bảo backend đang chạy trên cổng 3000:

```bash
# Kiểm tra backend
curl http://localhost:3000/api/health

# Hoặc xem terminal backend, phải thấy:
# Server running on :3000
```

## ✅ Checklist

- [ ] Backend chạy trên port 3000
- [ ] Frontend dev server đã khởi động lại
- [ ] Mở http://localhost:5173
- [ ] Đăng nhập thành công
- [ ] Vào Cashier Dashboard
- [ ] Thấy "Tiền đang quản lý"
- [ ] Không có lỗi 404

## 🎯 Nếu Vẫn Lỗi

### Kiểm tra cổng backend

```bash
# Xem backend đang chạy trên cổng nào
lsof -i :3000
# Hoặc
netstat -an | grep 3000
```

### Xóa cache trình duyệt

1. Mở DevTools (F12)
2. Chuột phải vào nút refresh
3. Chọn "Empty Cache and Hard Reload"

### Test API trực tiếp

```bash
# Lấy token từ localStorage
TOKEN="your_token_here"

# Test endpoint
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/cashier-shifts/current
```

---

**QUAN TRỌNG**: Phải khởi động lại frontend dev server để thay đổi có hiệu lực!

```bash
cd frontend
# Ctrl+C để dừng
npm run dev  # Khởi động lại
```
