# Local Print Bridge - Docker Deployment Guide

## Tổng Quan

Hướng dẫn deploy Local Print Bridge bằng Docker để dễ dàng cài đặt và quản lý.

## Ưu Điểm Docker

✅ **Dễ cài đặt** - Không cần cài Node.js thủ công  
✅ **Tự động khởi động** - Service tự động start khi máy khởi động  
✅ **Dễ quản lý** - Dùng docker-compose để start/stop/restart  
✅ **Logs tập trung** - Xem logs dễ dàng với docker-compose logs  
✅ **Cô lập** - Không ảnh hưởng đến hệ thống khác  

## Yêu Cầu

- Docker Desktop (Windows/Mac) hoặc Docker Engine (Linux)
- Docker Compose
- Máy in kết nối mạng LAN
- Kết nối internet (để kết nối backend EC2)

## Cài Đặt Docker

### Windows

1. Tải Docker Desktop: https://www.docker.com/products/docker-desktop/
2. Cài đặt và khởi động Docker Desktop
3. Verify: Mở Command Prompt và chạy `docker --version`

### macOS

```bash
# Dùng Homebrew
brew install --cask docker

# Hoặc tải từ: https://www.docker.com/products/docker-desktop/
```

### Linux (Ubuntu/Debian)

```bash
# Cài Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Cài Docker Compose
sudo apt-get install docker-compose-plugin

# Add user vào docker group
sudo usermod -aG docker $USER
newgrp docker
```

## Triển Khai Nhanh (3 Phút)

### Bước 1: Chuẩn Bị

```bash
cd local-print-bridge
```

### Bước 2: Cấu Hình

```bash
# Copy file cấu hình mẫu
cp .env.docker .env

# Chỉnh sửa file .env
nano .env  # hoặc dùng text editor bất kỳ
```

**Cập nhật các giá trị sau:**

```env
BACKEND_URL=https://your-actual-ec2-domain.com
DEFAULT_BILL_PRINTER_IP=192.168.1.100
DEFAULT_LABEL_PRINTER_IP=192.168.1.101
```

### Bước 3: Khởi Động

**Cách 1: Dùng script tự động**

```bash
./docker-start.sh
```

**Cách 2: Dùng docker-compose thủ công**

```bash
# Build image
docker-compose build

# Start service
docker-compose up -d

# Xem logs
docker-compose logs -f
```

### Bước 4: Kiểm Tra

```bash
# Test health endpoint
curl http://localhost:3001/health

# Kết quả mong đợi:
# {"status":"ok","service":"Local Print Bridge",...}
```

## Quản Lý Service

### Xem Trạng Thái

```bash
docker-compose ps
```

### Xem Logs

```bash
# Xem logs real-time
docker-compose logs -f

# Xem 50 dòng cuối
docker-compose logs --tail=50

# Xem logs của 1 giờ qua
docker-compose logs --since 1h
```

### Dừng Service

```bash
docker-compose stop
```

### Khởi Động Lại

```bash
docker-compose restart
```

### Dừng và Xóa Container

```bash
docker-compose down
```

### Cập Nhật Code

```bash
# Pull code mới
git pull

# Rebuild và restart
docker-compose up -d --build
```

## Cấu Hình Network

### Chế Độ Host Network (Mặc Định)

File `docker-compose.yml` sử dụng `network_mode: host`:

**Ưu điểm:**
- Container truy cập trực tiếp vào mạng LAN
- Không cần port mapping
- Dễ kết nối với máy in

**Nhược điểm:**
- Chỉ hoạt động trên Linux
- Trên Windows/Mac cần dùng bridge network

### Chế Độ Bridge Network (Windows/Mac)

Nếu bạn dùng Windows hoặc Mac, sửa `docker-compose.yml`:

```yaml
services:
  print-bridge:
    # Xóa dòng này:
    # network_mode: host
    
    # Thêm port mapping:
    ports:
      - "3001:3001"
    
    # Thêm network:
    networks:
      - print-network

networks:
  print-network:
    driver: bridge
```

**Lưu ý:** Với bridge network, máy in phải accessible từ Docker container.

## Test Kết Nối Máy In

### Từ Host Machine

```bash
# Test ping
ping 192.168.1.100

# Test port 9100
nc -zv 192.168.1.100 9100
```

### Từ Docker Container

```bash
# Vào container
docker-compose exec print-bridge sh

# Test ping (nếu có ping)
ping 192.168.1.100

# Test với Node.js
node src/test-printer.js 192.168.1.100 9100

# Thoát container
exit
```

## Tự Động Khởi Động

### Docker Desktop (Windows/Mac)

1. Mở Docker Desktop Settings
2. Chọn "General"
3. Bật "Start Docker Desktop when you log in"
4. Service sẽ tự động start với `restart: unless-stopped`

### Linux (systemd)

Docker service tự động start:

```bash
# Enable Docker service
sudo systemctl enable docker

# Service sẽ tự động start với restart policy
```

## Monitoring

### Health Check

Docker tự động check health mỗi 30 giây:

```bash
# Xem health status
docker-compose ps

# Kết quả:
# NAME                  STATUS
# local-print-bridge    Up (healthy)
```

### Resource Usage

```bash
# Xem CPU, Memory usage
docker stats local-print-bridge
```

### Logs Rotation

Logs tự động rotate với cấu hình:
- Max size: 10MB per file
- Max files: 3 files
- Total: ~30MB logs

## Troubleshooting

### Container Không Start

```bash
# Xem logs chi tiết
docker-compose logs

# Xem logs của Docker daemon
docker logs local-print-bridge
```

**Nguyên nhân thường gặp:**
- Port 3001 đã được sử dụng
- File .env không đúng format
- BACKEND_URL không hợp lệ

### Không Kết Nối Được Máy In

```bash
# Kiểm tra network mode
docker-compose config | grep network_mode

# Test từ container
docker-compose exec print-bridge sh
ping 192.168.1.100
```

**Giải pháp:**
- Đảm bảo máy in và máy tính cùng mạng LAN
- Kiểm tra firewall không block port 9100
- Thử dùng host network mode (Linux)

### Backend Không Nhận Status Update

```bash
# Kiểm tra BACKEND_URL
docker-compose exec print-bridge sh
echo $BACKEND_URL

# Test kết nối backend
curl $BACKEND_URL/health
```

**Giải pháp:**
- Cập nhật BACKEND_URL đúng trong .env
- Restart container: `docker-compose restart`

### Container Bị Crash

```bash
# Xem logs trước khi crash
docker-compose logs --tail=100

# Xem exit code
docker-compose ps -a
```

## Backup và Restore

### Backup Configuration

```bash
# Backup .env file
cp .env .env.backup

# Backup toàn bộ folder
tar -czf print-bridge-backup.tar.gz local-print-bridge/
```

### Restore

```bash
# Restore .env
cp .env.backup .env

# Restart service
docker-compose restart
```

## Performance Tuning

### Giới Hạn Resources

Thêm vào `docker-compose.yml`:

```yaml
services:
  print-bridge:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 256M
        reservations:
          cpus: '0.25'
          memory: 128M
```

### Tối Ưu Logs

```yaml
logging:
  driver: "json-file"
  options:
    max-size: "5m"    # Giảm từ 10m
    max-file: "2"     # Giảm từ 3
```

## Security Best Practices

### 1. Không Expose Port Ra Internet

```yaml
# ĐÚNG: Chỉ bind localhost
ports:
  - "127.0.0.1:3001:3001"

# SAI: Expose ra tất cả interfaces
ports:
  - "3001:3001"
```

### 2. Chạy Với Non-Root User

Dockerfile đã cấu hình sẵn:

```dockerfile
USER nodejs  # Non-root user
```

### 3. Giới Hạn Capabilities

```yaml
security_opt:
  - no-new-privileges:true
cap_drop:
  - ALL
```

## So Sánh: Docker vs PM2

| Feature | Docker | PM2 |
|---------|--------|-----|
| Cài đặt | Cần Docker | Cần Node.js |
| Khởi động | docker-compose up | pm2 start |
| Logs | docker-compose logs | pm2 logs |
| Auto-restart | ✅ Built-in | ✅ Built-in |
| Resource limit | ✅ Dễ dàng | ❌ Khó |
| Isolation | ✅ Hoàn toàn | ❌ Không |
| Network | ⚠️ Phức tạp hơn | ✅ Đơn giản |

**Khuyến nghị:**
- **Docker**: Nếu đã có Docker, hoặc muốn isolation tốt
- **PM2**: Nếu đã có Node.js, hoặc cần network đơn giản

## Cheat Sheet

```bash
# Start
docker-compose up -d

# Stop
docker-compose stop

# Restart
docker-compose restart

# Logs
docker-compose logs -f

# Status
docker-compose ps

# Update
docker-compose up -d --build

# Remove
docker-compose down

# Shell access
docker-compose exec print-bridge sh

# Health check
curl http://localhost:3001/health

# Statistics
curl http://localhost:3001/status
```

## Tích Hợp Với POS

Sau khi deploy Docker:

1. **Mở POS trong browser**
2. **Vào Print Management**
3. **Kiểm tra header** - Phải hiện "Local Bridge Online" (màu xanh)
4. **Tạo đơn hàng test**
5. **Verify in tự động**

## Support

**Tài liệu liên quan:**
- [Integration Guide](LOCAL_PRINT_BRIDGE_INTEGRATION.md)
- [Quick Start Guide](LOCAL_PRINT_BRIDGE_QUICK_START.md)
- [Deployment Checklist](LOCAL_PRINT_BRIDGE_DEPLOYMENT_CHECKLIST.md)

**Docker Documentation:**
- https://docs.docker.com/
- https://docs.docker.com/compose/

## Kết Luận

Docker deployment giúp việc cài đặt và quản lý Local Print Bridge trở nên đơn giản hơn nhiều. Chỉ cần 3 phút để setup và service sẽ tự động chạy mỗi khi khởi động máy.

**Các bước chính:**
1. Cài Docker
2. Copy và sửa .env
3. Chạy `./docker-start.sh`
4. Done! ✅
