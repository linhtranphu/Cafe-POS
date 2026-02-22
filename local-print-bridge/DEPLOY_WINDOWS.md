# Hướng Dẫn Deploy Local Print Bridge trên Windows PC tại Quán

## Yêu Cầu

- Windows 10/11
- Docker Desktop for Windows (hoặc Node.js 18+)
- Kết nối internet để truy cập EC2
- Kết nối LAN đến máy in

## Bước 1: Cài Đặt Docker Desktop (Khuyến nghị)

1. Download Docker Desktop: https://www.docker.com/products/docker-desktop/
2. Cài đặt và khởi động Docker Desktop
3. Đảm bảo Docker đang chạy (icon Docker ở system tray)

## Bước 2: Tạo Thư Mục và File Cấu Hình

1. Tạo thư mục `C:\cafe-print-bridge`
2. Tạo file `.env` trong thư mục đó với nội dung:

```bash
# Server Configuration
PORT=3001

# Backend URL - QUAN TRỌNG: Thay YOUR_EC2_IP bằng IP/domain thực tế
BACKEND_URL=http://YOUR_EC2_IP:3000

# Printer Configuration
PRINTER_TIMEOUT=5000
DEFAULT_BILL_PRINTER_IP=192.168.1.115
DEFAULT_LABEL_PRINTER_IP=192.168.1.101

# Node Environment
NODE_ENV=production
```

3. **Cập nhật các giá trị:**
   - `YOUR_EC2_IP` → IP hoặc domain của EC2 server
   - `DEFAULT_BILL_PRINTER_IP` → IP máy in hóa đơn (80mm)
   - `DEFAULT_LABEL_PRINTER_IP` → IP máy in nhãn (58mm)

## Bước 3: Chạy Print Bridge

### Option A: Dùng Docker (Khuyến nghị)

Mở PowerShell hoặc Command Prompt và chạy:

```powershell
cd C:\cafe-print-bridge

# Pull image từ Docker Hub
docker pull linhtranphu/local-print-bridge:latest

# Chạy container
docker run -d `
  --name print-bridge `
  --restart unless-stopped `
  -p 3001:3001 `
  --env-file .env `
  linhtranphu/local-print-bridge:latest
```

### Option B: Dùng Node.js (Nếu không dùng Docker)

```powershell
# Clone hoặc copy source code
cd C:\cafe-print-bridge

# Cài đặt dependencies
npm install

# Chạy service
npm start
```

## Bước 4: Kiểm Tra Hoạt Động

1. Mở trình duyệt: http://localhost:3001/health
2. Nên thấy: `{"status":"ok","service":"Local Print Bridge",...}`

## Bước 5: Test In

```powershell
# Test kết nối đến backend
curl http://YOUR_EC2_IP:3000/api/manager/print-jobs/pending

# Xem logs
docker logs -f print-bridge
```

## Bước 6: Cấu Hình Tự Động Khởi Động

### Docker Desktop:
- Docker Desktop sẽ tự động start container khi Windows khởi động
- Container có flag `--restart unless-stopped`

### Node.js:
Dùng `pm2` hoặc Windows Service:

```powershell
# Cài pm2
npm install -g pm2

# Start với pm2
pm2 start src/index.js --name print-bridge

# Cấu hình auto-start
pm2 startup
pm2 save
```

## Troubleshooting

### Lỗi: Cannot connect to backend

**Nguyên nhân:** EC2 IP/domain sai hoặc firewall chặn

**Giải pháp:**
1. Kiểm tra EC2 Security Group cho phép inbound port 3000
2. Test: `curl http://YOUR_EC2_IP:3000/health`
3. Kiểm tra `.env` có đúng BACKEND_URL không

### Lỗi: Cannot connect to printer

**Nguyên nhân:** Máy in không online hoặc IP sai

**Giải pháp:**
1. Ping máy in: `ping 192.168.1.115`
2. Kiểm tra máy in có bật không
3. Kiểm tra IP máy in (in test page từ máy in)
4. Đảm bảo Windows PC và máy in cùng mạng LAN

### Lỗi: Port 3001 already in use

**Giải pháp:**
```powershell
# Tìm process đang dùng port 3001
netstat -ano | findstr :3001

# Kill process (thay PID bằng số thực tế)
taskkill /PID <PID> /F
```

## Cập Nhật Version Mới

```powershell
# Stop container hiện tại
docker stop print-bridge
docker rm print-bridge

# Pull version mới
docker pull linhtranphu/local-print-bridge:latest

# Chạy lại (dùng lệnh ở Bước 3)
docker run -d --name print-bridge --restart unless-stopped -p 3001:3001 --env-file .env linhtranphu/local-print-bridge:latest
```

## Monitoring

### Xem logs:
```powershell
docker logs -f print-bridge
```

### Xem status:
```powershell
docker ps | findstr print-bridge
```

### Restart service:
```powershell
docker restart print-bridge
```

## Lưu Ý Quan Trọng

1. ✅ Windows PC phải luôn bật khi quán hoạt động
2. ✅ Đảm bảo internet ổn định (để poll backend)
3. ✅ Máy in phải cùng mạng LAN với Windows PC
4. ✅ Không tắt Docker Desktop khi đang hoạt động
5. ✅ Backup file `.env` để dễ restore khi cần

## Liên Hệ Hỗ Trợ

Nếu gặp vấn đề, kiểm tra logs và liên hệ IT support với thông tin:
- Logs: `docker logs print-bridge > logs.txt`
- Config: Nội dung file `.env`
- Error message cụ thể
