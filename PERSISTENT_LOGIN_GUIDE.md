# Hướng Dẫn: Lưu Thông Tin Đăng Nhập (Persistent Login)

## Tính Năng
Hệ thống giờ đây sẽ lưu thông tin đăng nhập vào localStorage, cho phép người dùng không bị bắt đăng nhập lại khi refresh trang.

## Cách Hoạt Động

### 1. Khi Đăng Nhập
- Token JWT được lưu vào `localStorage.token`
- Thông tin user được lưu vào `localStorage.user`
- Token được set vào header `Authorization: Bearer <token>` cho tất cả API requests

### 2. Khi Refresh Trang
- App tự động gọi `authStore.initAuth()` trong `main.js`
- Khôi phục token và user từ localStorage
- Set token vào API headers
- User vẫn ở trạng thái đăng nhập

### 3. Khi Đăng Xuất
- Xóa token và user từ localStorage
- Xóa Authorization header
- Redirect về trang login

## Các File Được Cập Nhật

### `frontend/src/main.js`
```javascript
// Khôi phục auth từ localStorage khi app load
const authStore = useAuthStore()
authStore.initAuth()
```

### `frontend/src/stores/auth.js`
- Thêm method `initAuth()` - khôi phục auth từ localStorage
- Thêm method `validateToken()` - validate token với backend (optional)
- Cập nhật `login()` - set Authorization header
- Cập nhật `logout()` - xóa Authorization header

### `frontend/src/services/api.js`
- Đã có interceptor để tự động thêm token vào requests

## Cách Test

### Test 1: Đăng Nhập và Refresh
1. Truy cập http://localhost
2. Đăng nhập với `admin` / `admin123`
3. Mở DevTools (F12) → Application → LocalStorage
4. Xác nhận có `token` và `user` được lưu
5. Refresh trang (Ctrl+R)
6. Xác nhận vẫn ở trạng thái đăng nhập, không bị redirect về login

### Test 2: Đăng Xuất
1. Đăng nhập
2. Click Logout
3. Xác nhận localStorage bị xóa
4. Xác nhận redirect về login

### Test 3: Token Hết Hạn
1. Đăng nhập
2. Chờ token hết hạn (hoặc xóa token từ localStorage)
3. Refresh trang
4. Xác nhận redirect về login

## Bảo Mật

⚠️ **Lưu Ý:**
- Token được lưu trong localStorage (không phải httpOnly cookie)
- Điều này có thể bị XSS attack nếu có lỗ hổng
- Trong production, nên sử dụng httpOnly cookies
- Luôn validate token với backend khi cần

## Tương Lai

Có thể thêm:
- Refresh token mechanism
- Token expiration check
- Auto-logout khi token hết hạn
- Session timeout
