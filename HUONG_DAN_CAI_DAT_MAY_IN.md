# Hướng Dẫn Cài Đặt và Sử Dụng Máy In

## Mục Lục
1. [Yêu Cầu Hệ Thống](#yêu-cầu-hệ-thống)
2. [Chuẩn Bị Máy In](#chuẩn-bị-máy-in)
3. [Cài Đặt Local Print Bridge](#cài-đặt-local-print-bridge)
4. [Cấu Hình Máy In Trong Hệ Thống](#cấu-hình-máy-in-trong-hệ-thống)
5. [Kiểm Tra và Test](#kiểm-tra-và-test)
6. [Sử Dụng Hàng Ngày](#sử-dụng-hàng-ngày)
7. [Xử Lý Sự Cố](#xử-lý-sự-cố)

---

## Yêu Cầu Hệ Thống

### Phần Cứng
- ✅ Máy in nhiệt ESC/POS (khổ 80mm hoặc 58mm)
- ✅ Router WiFi hoặc Switch mạng
- ✅ Máy tính/Laptop chạy Windows/macOS/Linux
- ✅ Cáp mạng LAN (nếu dùng kết nối có dây)

### Phần Mềm
- ✅ Node.js phiên bản 16 trở lên
- ✅ Trình duyệt web hiện đại (Chrome, Firefox, Safari, Edge)
- ✅ Kết nối internet ổn định

### Mạng
- ✅ Máy in và máy tính phải cùng mạng LAN
- ✅ Không cần mở port ra internet
- ✅ Không cần cấu hình firewall đặc biệt

---

## Chuẩn Bị Máy In

### Bước 1: Kết Nối Máy In Vào Mạng

#### Với Máy In WiFi:
1. Bật máy in
2. Nhấn giữ nút WiFi trên máy in (thường là 3-5 giây)
3. Đèn WiFi sẽ nhấp nháy
4. Kết nối điện thoại/laptop vào WiFi của máy in (tên thường là "Printer-XXXX")
5. Mở trình duyệt, truy cập `http://192.168.192.168` hoặc `http://192.168.1.1`
6. Chọn mạng WiFi quán cafe và nhập mật khẩu
7. Máy in sẽ tự động kết nối và in ra giấy thông báo địa chỉ IP

#### Với Máy In LAN (Cáp Mạng):
1. Cắm cáp mạng từ máy in vào router/switch
2. Bật máy in
3. Nhấn nút "Feed" hoặc "Test" để in thông tin cấu hình
4. Ghi lại địa chỉ IP (ví dụ: `192.168.1.100`)

### Bước 2: Kiểm Tra Kết Nối

Mở Terminal/Command Prompt và chạy:

```bash
ping 192.168.1.100
```

Thay `192.168.1.100` bằng IP của máy in. Nếu thấy phản hồi thì máy in đã kết nối thành công.

### Bước 3: Ghi Chú Thông Tin

Ghi lại các thông tin sau:
- **IP máy in Bill:** `192.168.1.100` (ví dụ)
- **IP máy in Label:** `192.168.1.101` (ví dụ)
- **Port:** `9100` (mặc định cho ESC/POS)

---

## Cài Đặt Local Print Bridge

Local Print Bridge là phần mềm chạy trên máy tính tại quán, giúp kết nối từ hệ thống POS trên cloud xuống máy in local.

### Phương Án 1: Cài Đặt Bằng Docker (Khuyến Nghị) ⭐

**Ưu điểm:**
- Cài đặt nhanh nhất (1 lệnh)
- Tự động khởi động khi máy tính bật
- Dễ quản lý và cập nhật

**Các bước:**

1. **Cài đặt Docker:**
   - Windows/Mac: Tải Docker Desktop từ https://www.docker.com/products/docker-desktop
   - Linux: `curl -fsSL https://get.docker.com | sh`

2. **Tạo file cấu hình:**
   ```bash
   cd local-print-bridge
   cp .env.docker .env
   ```

3. **Chỉnh sửa file `.env`:**
   ```bash
   nano .env
   ```
   
   Nội dung:
   ```env
   PORT=3001
   BACKEND_URL=https://api.yourdomain.com
   DEFAULT_BILL_PRINTER_IP=192.168.1.100
   DEFAULT_BILL_PRINTER_PORT=9100
   DEFAULT_LABEL_PRINTER_IP=192.168.1.101
   DEFAULT_LABEL_PRINTER_PORT=9100
   ```

4. **Khởi động service:**
   ```bash
   ./docker-start.sh
   ```

5. **Kiểm tra:**
   ```bash
   docker ps
   ```
   
   Bạn sẽ thấy container `local-print-bridge` đang chạy.

### Phương Án 2: Cài Đặt Trực Tiếp (Node.js)

**Các bước:**

1. **Cài đặt Node.js:**
   - Tải từ https://nodejs.org (chọn phiên bản LTS)
   - Cài đặt theo hướng dẫn

2. **Cài đặt dependencies:**
   ```bash
   cd local-print-bridge
   npm install
   ```

3. **Tạo file cấu hình:**
   ```bash
   cp .env.example .env
   nano .env
   ```
   
   Điền thông tin như Phương Án 1.

4. **Khởi động service:**
   ```bash
   npm start
   ```

5. **Chạy tự động khi khởi động máy:**
   
   **Windows:**
   ```bash
   npm install -g pm2-windows-startup
   pm2-startup install
   pm2 start src/index.js --name print-bridge
   pm2 save
   ```
   
   **macOS/Linux:**
   ```bash
   npm install -g pm2
   pm2 start src/index.js --name print-bridge
   pm2 startup
   pm2 save
   ```

### Kiểm Tra Local Print Bridge

Mở trình duyệt và truy cập:
```
http://localhost:3001/health
```

Bạn sẽ thấy:
```json
{
  "status": "ok",
  "service": "Local Print Bridge",
  "version": "1.0.0"
}
```

---

## Cấu Hình Máy In Trong Hệ Thống

### Bước 1: Đăng Nhập Hệ Thống POS

1. Mở trình duyệt
2. Truy cập `https://pos.yourdomain.com` (hoặc `http://localhost:5173` nếu chạy local)
3. Đăng nhập với tài khoản **Manager**

### Bước 2: Vào Trang Quản Lý In Ấn

1. Sau khi đăng nhập, bạn sẽ thấy trang Dashboard
2. Cuộn xuống phần **"👥 Nhân sự & Hệ thống"**
3. Click vào ô **🖨️ In ấn** (màu xanh lá - teal)
4. Trang "Quản Lý In Ấn" sẽ mở ra với 4 tab:
   - **📄 Print Jobs** - Lịch sử các lệnh in
   - **🖨️ Máy In** - Cấu hình máy in
   - **📝 Templates** - Mẫu in
   - **⚙️ Cài Đặt** - Cài đặt chung

### Bước 3: Thêm Cấu Hình Máy In

1. Click vào tab **🖨️ Máy In**
2. Click nút **Thêm Máy In** (hoặc **+ Thêm Máy In**)

#### Cấu Hình Máy In Bill:

```
Tên máy in: Máy In Bill Quầy
Loại: Bill (Hóa Đơn)
Địa chỉ IP: 192.168.1.100
Port: 9100
Mặc định: ✓ (Bật)
```

#### Cấu Hình Máy In Label:

```
Tên máy in: Máy In Label Bar
Loại: Label (Nhãn)
Địa chỉ IP: 192.168.1.101
Port: 9100
Mặc định: ✓ (Bật)
```

3. Click **Lưu**

### Bước 4: Kiểm Tra Kết Nối

1. Trong danh sách máy in, click nút **Test Kết Nối** (hoặc icon 🔌)
2. Nếu thành công, máy in sẽ in ra một trang test
3. Nếu thất bại, kiểm tra lại IP và kết nối mạng
4. Ở góc trên bên phải, bạn sẽ thấy trạng thái **Local Bridge**:
   - 🟢 **Local Bridge Online** - Đã kết nối
   - ⚪ **Local Bridge Offline** - Chưa kết nối

### Bước 5: Cấu Hình Template In (Tùy Chọn)

1. Click vào tab **📝 Templates**
2. Chọn mẫu mặc định hoặc tạo mẫu mới
3. Chỉnh sửa:
   - Tên quán
   - Địa chỉ
   - Số điện thoại
   - Logo (nếu có)
   - Độ rộng giấy (80mm hoặc 58mm)

### Bước 6: Bật Tự Động In

1. Click vào tab **⚙️ Cài Đặt**
2. Tìm mục **In Tự Động**
3. Bật các tùy chọn:
   - ✓ Tự động in bill khi thanh toán
   - ✓ Tự động in label khi gửi order
4. Click **Lưu**

---

## Kiểm Tra và Test

### Test 1: In Thử Từ Hệ Thống

1. Vào Dashboard, cuộn xuống phần **"👥 Nhân sự & Hệ thống"**
2. Click ô **🖨️ In ấn**
3. Chọn tab **🖨️ Máy In**
4. Tìm máy in Bill trong danh sách
5. Click nút **Test Kết Nối** hoặc **In Thử**
6. Kiểm tra:
   - ✅ Máy in có in ra không?
   - ✅ Nội dung có rõ ràng không?
   - ✅ Căn lề có đúng không?
   - ✅ Trạng thái Local Bridge có hiển thị "Online" không?

### Test 2: Tạo Order Thử

1. Vào **Bán Hàng**
2. Tạo một order test:
   - Chọn 1-2 món
   - Nhập tên khách: "Test"
   - Click **Gửi Đơn**
3. Kiểm tra:
   - ✅ Label có in ra tại bar không?
   - ✅ Thông tin món có đúng không?

### Test 3: Thanh Toán Thử

1. Với order test ở trên
2. Click **Thanh Toán**
3. Chọn phương thức thanh toán
4. Click **Hoàn Tất**
5. Kiểm tra:
   - ✅ Bill có in ra không?
   - ✅ Tổng tiền có đúng không?
   - ✅ Thông tin khách có đầy đủ không?

### Test 4: In Lại

1. Vào **Lịch Sử Order**
2. Chọn order vừa tạo
3. Click **In Lại Bill**
4. Kiểm tra máy in có in lại không

---

## Sử Dụng Hàng Ngày

### Quy Trình Bán Hàng Chuẩn

```
1. Nhân viên tạo order trên POS
   ↓
2. Click "Gửi Đơn"
   ↓
3. Label tự động in ra tại bar
   ↓
4. Barista nhận label và pha chế
   ↓
5. Khách thanh toán
   ↓
6. Bill tự động in ra tại quầy
   ↓
7. Giao bill cho khách
```

### Các Tính Năng In

#### 1. In Tự Động
- Bill in tự động khi thanh toán
- Label in tự động khi gửi order
- Không cần thao tác thêm

#### 2. In Lại
- Vào **Lịch Sử Order**
- Chọn order cần in lại
- Click **In Lại Bill** hoặc **In Lại Label**

#### 3. In Nhiều Bản
- Khi in lại, có thể chọn số lượng bản in
- Hữu ích khi cần nhiều label cho cùng món

### Theo Dõi Trạng Thái In

1. Vào Dashboard, click ô **🖨️ In ấn**
2. Tab **📄 Print Jobs** sẽ hiển thị danh sách các lệnh in:
   - ✅ **Hoàn Thành (COMPLETED):** In thành công
   - ⏳ **Đang Chờ (PENDING):** Đang xử lý
   - ❌ **Thất Bại (FAILED):** Có lỗi xảy ra

3. Với lệnh thất bại:
   - Click vào dòng đó để **Xem Chi Tiết** lỗi
   - Click nút **In Lại** (🔄) để thử lại

4. Kiểm tra trạng thái Local Bridge:
   - Ở góc trên bên phải trang
   - 🟢 **Local Bridge Online** = Sẵn sàng in
   - ⚪ **Local Bridge Offline** = Cần khởi động lại service

---

## Xử Lý Sự Cố

### Sự Cố 1: Máy In Không In

**Triệu chứng:**
- Order đã gửi nhưng không có gì in ra
- Trạng thái in hiển thị "Thất Bại"

**Cách khắc phục:**

1. **Kiểm tra máy in:**
   ```bash
   ping 192.168.1.100
   ```
   - Nếu không ping được → Kiểm tra kết nối mạng
   - Nếu ping được → Tiếp tục bước 2

2. **Kiểm tra Local Print Bridge:**
   - Mở `http://localhost:3001/health`
   - Nếu không mở được → Khởi động lại service:
     ```bash
     # Docker
     docker restart local-print-bridge
     
     # PM2
     pm2 restart print-bridge
     ```

3. **Kiểm tra giấy in:**
   - Mở nắp máy in
   - Kiểm tra còn giấy không
   - Lắp giấy đúng hướng (mặt nhiệt hướng lên)

4. **Test kết nối:**
   - Vào Dashboard, click **🖨️ In ấn**
   - Chọn tab **🖨️ Máy In**
   - Click **Test Kết Nối** trên máy in
   - Nếu thất bại → Kiểm tra lại IP

### Sự Cố 2: In Ra Giấy Trắng

**Nguyên nhân:**
- Giấy lắp ngược chiều
- Giấy hết mực nhiệt

**Cách khắc phục:**
1. Lấy giấy ra
2. Lật ngược lại (mặt nhiệt hướng lên)
3. Lắp lại và thử in

**Cách nhận biết mặt nhiệt:**
- Dùng móng tay cào nhẹ lên giấy
- Mặt nào có vết đen → Đó là mặt nhiệt

### Sự Cố 3: In Bị Lỗi Font/Ký Tự

**Triệu chứng:**
- Tiếng Việt hiển thị sai
- Ký tự đặc biệt bị lỗi

**Cách khắc phục:**
1. Vào Dashboard, click **🖨️ In ấn**
2. Chọn tab **📝 Templates**
3. Chọn template đang dùng
4. Tìm mục **Encoding** và chọn: **UTF-8**
5. Click **Lưu** và thử in lại

### Sự Cố 4: In Chậm

**Nguyên nhân:**
- Mạng WiFi yếu
- Máy in quá xa router

**Cách khắc phục:**
1. Đổi sang kết nối LAN (cáp mạng)
2. Đặt máy in gần router hơn
3. Nâng cấp router WiFi

### Sự Cố 5: Local Print Bridge Không Chạy

**Triệu chứng:**
- Không mở được `http://localhost:3001/health`
- Lỗi "Connection refused"

**Cách khắc phục:**

**Với Docker:**
```bash
# Kiểm tra container
docker ps -a

# Xem log
docker logs local-print-bridge

# Khởi động lại
docker restart local-print-bridge

# Nếu không có container
cd local-print-bridge
./docker-start.sh
```

**Với PM2:**
```bash
# Kiểm tra status
pm2 status

# Xem log
pm2 logs print-bridge

# Khởi động lại
pm2 restart print-bridge

# Nếu không có process
pm2 start src/index.js --name print-bridge
pm2 save
```

### Sự Cố 6: Máy In Bị Kẹt Giấy

**Cách khắc phục:**
1. Tắt máy in
2. Mở nắp máy in
3. Nhẹ nhàng kéo giấy ra theo chiều in
4. Đóng nắp
5. Bật máy in lại
6. Nhấn nút Feed để kiểm tra

---

## Bảo Trì Định Kỳ

### Hàng Ngày
- ✅ Kiểm tra giấy in còn đủ không
- ✅ Kiểm tra máy in có hoạt động không

### Hàng Tuần
- ✅ Lau sạch đầu in bằng cồn
- ✅ Kiểm tra kết nối mạng
- ✅ Xem log lỗi in (nếu có)

### Hàng Tháng
- ✅ Vệ sinh máy in (bên trong và bên ngoài)
- ✅ Kiểm tra cáp mạng có hỏng không
- ✅ Cập nhật firmware máy in (nếu có)
- ✅ Backup cấu hình hệ thống

### Hàng Quý
- ✅ Thay đầu in (nếu in mờ)
- ✅ Kiểm tra tổng thể hệ thống
- ✅ Đào tạo lại nhân viên (nếu cần)

---

## Lưu Ý Quan Trọng

### Về Giấy In
- ✅ Dùng giấy nhiệt chất lượng tốt (tránh phai màu)
- ✅ Bảo quản giấy nơi khô ráo, tránh ánh nắng
- ✅ Không để giấy gần nguồn nhiệt
- ✅ Kiểm tra hạn sử dụng của giấy

### Về Máy In
- ✅ Không để máy in gần nước
- ✅ Không để máy in nơi ẩm ướt
- ✅ Tắt máy in khi không sử dụng lâu
- ✅ Không tự ý tháo máy in

### Về Mạng
- ✅ Giữ IP máy in cố định (không đổi)
- ✅ Không thay đổi mật khẩu WiFi mà không cập nhật máy in
- ✅ Đảm bảo router luôn bật

### Về Bảo Mật
- ✅ Không chia sẻ thông tin cấu hình ra bên ngoài
- ✅ Chỉ cho phép máy tính tin cậy kết nối
- ✅ Thay đổi mật khẩu admin định kỳ

---

## Liên Hệ Hỗ Trợ

Nếu gặp vấn đề không giải quyết được, liên hệ:

**Hỗ Trợ Kỹ Thuật:**
- Email: support@yourdomain.com
- Hotline: 1900-xxxx-xxx
- Giờ làm việc: 8:00 - 22:00 (Thứ 2 - Chủ Nhật)

**Tài Liệu Tham Khảo:**
- [Hướng Dẫn Nhanh](LOCAL_PRINT_BRIDGE_QUICK_START.md)
- [Hướng Dẫn Docker](LOCAL_PRINT_BRIDGE_DOCKER_GUIDE.md)
- [Tài Liệu Kỹ Thuật](LOCAL_PRINT_BRIDGE_INTEGRATION.md)

---

## Phụ Lục

### A. Danh Sách Máy In Tương Thích

**Máy In Khổ 80mm:**
- Xprinter XP-80C
- HPRT TP808
- Rongta RP80
- Antech AP80

**Máy In Khổ 58mm:**
- Xprinter XP-58IIH
- HPRT TP58
- Rongta RP58

### B. Các Lệnh Hữu Ích

**Kiểm tra IP máy in:**
```bash
arp -a | grep "192.168.1"
```

**Test kết nối:**
```bash
telnet 192.168.1.100 9100
```

**Xem log Local Print Bridge:**
```bash
# Docker
docker logs -f local-print-bridge

# PM2
pm2 logs print-bridge
```

**Khởi động lại service:**
```bash
# Docker
docker restart local-print-bridge

# PM2
pm2 restart print-bridge
```

### C. Checklist Cài Đặt Ban Đầu

```
□ Máy in đã kết nối mạng
□ Đã ghi lại IP máy in
□ Đã cài đặt Local Print Bridge
□ Đã test kết nối máy in
□ Đã cấu hình máy in trong hệ thống
□ Đã test in thử
□ Đã bật tự động in
□ Đã đào tạo nhân viên
□ Đã backup cấu hình
```

---

**Phiên bản:** 1.0  
**Cập nhật:** February 16, 2026  
**Tác giả:** Cafe POS System Team

