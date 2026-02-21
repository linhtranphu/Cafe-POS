# Hướng Dẫn Truy Cập Trang Quản Lý In

## Cách Vào Trang Quản Lý In

### Bước 1: Đăng Nhập
1. Mở trình duyệt
2. Truy cập hệ thống POS
3. Đăng nhập với tài khoản **Manager**

### Bước 2: Vào Trang Quản Lý In
1. Sau khi đăng nhập, bạn sẽ thấy trang chủ với các ô menu dạng lưới
2. Tìm ô **🖨️ In ấn** (màu xanh lá - teal)
3. Click vào ô đó

### Bước 3: Các Tab Có Sẵn

Trang "Quản Lý In Ấn" có 4 tab chính:

#### 1. 📄 Print Jobs
- Xem lịch sử các lệnh in
- Trạng thái: PENDING, COMPLETED, FAILED
- Có thể in lại các job đã thất bại
- Xem chi tiết lỗi

#### 2. 🖨️ Máy In
- Danh sách tất cả máy in đã cấu hình
- Thêm máy in mới
- Chỉnh sửa cấu hình máy in
- Test kết nối máy in
- Xóa máy in

#### 3. 📝 Templates
- Quản lý mẫu in (Bill và Label)
- Tạo mẫu mới
- Chỉnh sửa mẫu hiện có
- Preview mẫu in
- Cấu hình:
  - Tên quán
  - Địa chỉ
  - Số điện thoại
  - Logo
  - Độ rộng giấy (80mm/58mm)

#### 4. ⚙️ Cài Đặt
- Bật/tắt tự động in
- Cấu hình in tự động khi:
  - Thanh toán (Bill)
  - Gửi order (Label)
- Các cài đặt chung khác

## Trạng Thái Local Bridge

Ở góc trên bên phải trang, bạn sẽ thấy trạng thái kết nối:

### 🟢 Local Bridge Online
- Màu xanh, có chấm nhấp nháy
- Nghĩa là: Đã kết nối với Local Print Bridge
- Có thể in được

### ⚪ Local Bridge Offline
- Màu xám
- Nghĩa là: Chưa kết nối với Local Print Bridge
- Không thể in được
- Cần khởi động Local Print Bridge

## Các Thao Tác Thường Dùng

### Thêm Máy In Mới
1. Vào tab **🖨️ Máy In**
2. Click **+ Thêm Máy In**
3. Điền thông tin:
   - Tên máy in
   - Loại (Bill/Label)
   - IP address
   - Port (thường là 9100)
   - Đặt làm mặc định (nếu cần)
4. Click **Lưu**

### Test Kết Nối Máy In
1. Vào tab **🖨️ Máy In**
2. Tìm máy in cần test
3. Click nút **Test Kết Nối** (hoặc icon 🔌)
4. Máy in sẽ in ra trang test nếu thành công

### Xem Lịch Sử In
1. Vào tab **📄 Print Jobs**
2. Xem danh sách các lệnh in
3. Lọc theo:
   - Trạng thái (All/Pending/Completed/Failed)
   - Loại (Bill/Label)
   - Thời gian

### In Lại Job Thất Bại
1. Vào tab **📄 Print Jobs**
2. Tìm job có trạng thái FAILED
3. Click nút **🔄 In Lại**
4. Hoặc click vào dòng đó để xem chi tiết lỗi

### Chỉnh Sửa Template
1. Vào tab **📝 Templates**
2. Chọn template cần chỉnh sửa
3. Click **Chỉnh Sửa**
4. Thay đổi nội dung
5. Click **Preview** để xem trước
6. Click **Lưu**

### Bật Tự Động In
1. Vào tab **⚙️ Cài Đặt**
2. Tìm mục "In Tự Động"
3. Bật các tùy chọn:
   - ✓ Tự động in bill khi thanh toán
   - ✓ Tự động in label khi gửi order
4. Click **Lưu**

## Lưu Ý Quan Trọng

### Quyền Truy Cập
- Chỉ tài khoản **Manager** mới có quyền truy cập trang này
- Các role khác (Waiter, Barista, Cashier) không thấy menu này

### Local Bridge
- Phải khởi động Local Print Bridge trước khi in
- Kiểm tra trạng thái ở góc trên bên phải
- Nếu Offline, khởi động lại service:
  ```bash
  # Docker
  docker restart local-print-bridge
  
  # PM2
  pm2 restart print-bridge
  ```

### Kết Nối Mạng
- Máy in và máy tính phải cùng mạng LAN
- IP máy in phải cố định (không đổi)
- Port mặc định là 9100 cho ESC/POS

## Troubleshooting

### Không Thấy Menu "In ấn"
- Kiểm tra bạn đã đăng nhập với tài khoản Manager chưa
- Refresh lại trang (F5)
- Xóa cache trình duyệt

### Local Bridge Luôn Offline
- Kiểm tra Local Print Bridge có đang chạy không:
  ```bash
  # Docker
  docker ps | grep print-bridge
  
  # PM2
  pm2 status
  ```
- Kiểm tra port 3001 có bị chiếm không:
  ```bash
  lsof -i :3001
  ```
- Xem log để biết lỗi:
  ```bash
  # Docker
  docker logs local-print-bridge
  
  # PM2
  pm2 logs print-bridge
  ```

### Không Test Được Máy In
- Kiểm tra IP máy in có đúng không
- Ping thử máy in:
  ```bash
  ping 192.168.1.100
  ```
- Kiểm tra máy in có bật không
- Kiểm tra giấy in còn không

## Liên Hệ Hỗ Trợ

Nếu gặp vấn đề, liên hệ:
- Email: support@yourdomain.com
- Hotline: 1900-xxxx-xxx

---

**Cập nhật:** February 16, 2026
