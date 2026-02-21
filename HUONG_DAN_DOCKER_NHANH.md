# Hướng Dẫn Deploy Local Print Bridge Bằng Docker

## 🚀 Triển Khai Nhanh (3 Phút)

### Bước 1: Cài Docker

**Windows/Mac:**
- Tải Docker Desktop: https://www.docker.com/products/docker-desktop/
- Cài đặt và khởi động

**Linux:**
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
newgrp docker
```

### Bước 2: Cấu Hình

```bash
cd local-print-bridge
cp .env.docker .env
nano .env
```

**Sửa 3 dòng này:**
```env
BACKEND_URL=https://your-ec2-domain.com     # ← Sửa thành domain EC2 thật
DEFAULT_BILL_PRINTER_IP=192.168.1.100       # ← Sửa thành IP máy in bill
DEFAULT_LABEL_PRINTER_IP=192.168.1.101      # ← Sửa thành IP máy in tem
```

### Bước 3: Khởi Động

```bash
./docker-start.sh
```

**Hoặc:**
```bash
docker-compose up -d
```

### Bước 4: Kiểm Tra

```bash
curl http://localhost:3001/health
```

**Kết quả mong đợi:**
```json
{"status":"ok","service":"Local Print Bridge",...}
```

## ✅ Xong! Service Đã Chạy

Mở POS trong browser, vào Print Management, phải thấy:
- **"Local Bridge Online"** (màu xanh) ở góc trên

## 📋 Các Lệnh Thường Dùng

```bash
# Xem logs
docker-compose logs -f

# Dừng service
docker-compose stop

# Khởi động lại
docker-compose restart

# Xem trạng thái
docker-compose ps

# Xóa container
docker-compose down
```

## 🔧 Troubleshooting

### "Local Bridge Offline" trong POS

```bash
# Kiểm tra service có chạy không
docker-compose ps

# Nếu không chạy, start lại
docker-compose up -d

# Xem logs để biết lỗi
docker-compose logs
```

### Máy In Không In

```bash
# Test kết nối máy in
curl -X POST http://localhost:3001/test-connection \
  -H "Content-Type: application/json" \
  -d '{"printerIP":"192.168.1.100","printerPort":9100}'
```

**Nếu lỗi:**
- Kiểm tra IP máy in đúng chưa
- Ping máy in: `ping 192.168.1.100`
- Kiểm tra máy in có bật không

### Backend Không Nhận Status

```bash
# Kiểm tra BACKEND_URL
cat .env | grep BACKEND_URL

# Test kết nối backend
curl https://your-ec2-domain.com/health
```

**Nếu lỗi:**
- Sửa BACKEND_URL trong .env
- Restart: `docker-compose restart`

## 🎯 Test End-to-End

1. Mở POS trong browser
2. Tạo đơn hàng mới
3. Thanh toán (mark as PAID)
4. Máy in sẽ tự động in bill
5. Vào Print Jobs tab → thấy job status "COMPLETED"

## 📊 Xem Thống Kê

```bash
curl http://localhost:3001/status
```

**Kết quả:**
```json
{
  "stats": {
    "totalPrints": 150,
    "successfulPrints": 148,
    "failedPrints": 2,
    "successRate": "98.67%"
  }
}
```

## 🔄 Cập Nhật Code Mới

```bash
# Pull code mới
git pull

# Rebuild và restart
docker-compose up -d --build
```

## 🛡️ Tự Động Khởi Động

Service sẽ **tự động start** khi máy khởi động nhờ:
```yaml
restart: unless-stopped
```

Không cần làm gì thêm!

## 💡 Tips

### Xem Logs Real-Time
```bash
docker-compose logs -f
```

### Xem 50 Dòng Logs Cuối
```bash
docker-compose logs --tail=50
```

### Xem Resource Usage
```bash
docker stats local-print-bridge
```

### Vào Trong Container
```bash
docker-compose exec print-bridge sh
```

## 📚 Tài Liệu Chi Tiết

- [Hướng dẫn Docker đầy đủ](LOCAL_PRINT_BRIDGE_DOCKER_GUIDE.md)
- [Tổng quan Docker](LOCAL_PRINT_BRIDGE_DOCKER_SUMMARY.md)
- [Integration Guide](LOCAL_PRINT_BRIDGE_INTEGRATION.md)

## ❓ Cần Giúp Đỡ?

1. Xem logs: `docker-compose logs`
2. Check status: `docker-compose ps`
3. Test health: `curl localhost:3001/health`
4. Đọc troubleshooting trong [Docker Guide](LOCAL_PRINT_BRIDGE_DOCKER_GUIDE.md)

## 🎉 Hoàn Thành!

Bây giờ bạn đã có Local Print Bridge chạy trong Docker:
- ✅ Tự động khởi động khi máy bật
- ✅ Tự động restart nếu crash
- ✅ Logs tự động rotate
- ✅ Dễ dàng quản lý với docker-compose

**Chúc in ấn thành công! 🖨️**
