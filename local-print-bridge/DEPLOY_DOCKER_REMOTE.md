# Hướng Dẫn Deploy Local Print Bridge Trên Máy Docker Khác

## Tổng Quan

Hướng dẫn này giúp bạn deploy Local Print Bridge trên một máy Docker riêng biệt (không phải máy chạy backend), thường là:
- Máy tính tại quầy thu ngân
- Máy tính kết nối trực tiếp với máy in
- Server riêng trong mạng LAN

## Yêu Cầu

### Máy Docker (Print Bridge Server)
- Docker và Docker Compose đã cài đặt
- Kết nối mạng LAN với máy in
- Có thể truy cập từ máy backend
- RAM: tối thiểu 512MB
- CPU: 1 core

### Mạng
- Máy in và máy Docker trong cùng mạng LAN
- Máy backend có thể truy cập máy Docker qua HTTP
- Port 3001 không bị firewall chặn

## Bước 1: Chuẩn Bị Máy Docker

### 1.1. Cài Docker (nếu chưa có)

**Ubuntu/Debian:**
```bash
# Update package list
sudo apt update

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo apt install docker-compose-plugin

# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker
```

**macOS:**
```bash
# Install Docker Desktop
brew install --cask docker
```

**Windows:**
- Download Docker Desktop từ https://www.docker.com/products/docker-desktop

### 1.2. Kiểm Tra Docker

```bash
# Check Docker version
docker --version

# Check Docker Compose
docker compose version

# Test Docker
docker run hello-world
```

## Bước 2: Chuẩn Bị Files

### Cách 1: Pull Image Từ DockerHub (Khuyến nghị - Đơn giản nhất)

**Không cần clone code, chỉ cần tạo 2 files!**

```bash
# Tạo thư mục
mkdir -p ~/print-bridge
cd ~/print-bridge
```

Xem chi tiết tại [DOCKERHUB_GUIDE.md](./DOCKERHUB_GUIDE.md)

### Cách 2: Dùng Git

```bash
# Clone repository
git clone <your-repo-url>
cd <repo-name>/local-print-bridge
```

### Cách 3: Copy Thủ Công

Trên máy có code, tạo archive:
```bash
cd local-print-bridge
tar -czf print-bridge.tar.gz .
```

Copy sang máy Docker:
```bash
# Dùng scp
scp print-bridge.tar.gz user@docker-machine:/home/user/

# Trên máy Docker, extract
cd /home/user
tar -xzf print-bridge.tar.gz
cd local-print-bridge
```

### Cách 4: Dùng USB

```bash
# Copy folder local-print-bridge vào USB
# Cắm USB vào máy Docker
# Copy vào máy Docker
cp -r /media/usb/local-print-bridge ~/
cd ~/local-print-bridge
```

## Bước 3: Cấu Hình Print Bridge

### 3.1. Tạo File .env

```bash
cd local-print-bridge

# Copy từ example
cp .env.example .env

# Edit cấu hình
nano .env
```

### 3.2. Cấu Hình Printer IPs

Sửa file `.env`:

```bash
# Server Port
PORT=3001

# Printer IPs trong mạng LAN của bạn
DEFAULT_BILL_PRINTER_IP=192.168.1.115
DEFAULT_BILL_PRINTER_PORT=9100

DEFAULT_LABEL_PRINTER_IP=192.168.1.116
DEFAULT_LABEL_PRINTER_PORT=9100

# Logging
LOG_LEVEL=info

# Timeout
PRINTER_TIMEOUT=5
```

**Lưu ý:** Thay đổi IP addresses theo máy in thực tế của bạn.

### 3.3. Kiểm Tra Kết Nối Máy In

```bash
# Ping máy in
ping 192.168.1.115

# Check port máy in
nc -zv 192.168.1.115 9100

# Hoặc dùng telnet
telnet 192.168.1.115 9100
```

## Bước 4: Build và Deploy

### 4.1. Build Docker Image

```bash
cd local-print-bridge

# Build image
docker compose build

# Hoặc build không dùng cache
docker compose build --no-cache
```

### 4.2. Start Service

```bash
# Start service
docker compose up -d

# Check logs
docker compose logs -f

# Check status
docker compose ps
```

### 4.3. Verify Service

```bash
# Health check
curl http://localhost:3001/health

# Expected response:
# {"status":"ok","timestamp":"..."}
```

## Bước 5: Test Print Bridge

### 5.1. Test Connection

```bash
# Test kết nối với máy in
curl -X POST http://localhost:3001/test-connection \
  -H "Content-Type: application/json" \
  -d '{
    "printerIP": "192.168.1.115",
    "printerPort": 9100
  }'
```

### 5.2. Test Print HTML

```bash
# Test in HTML
curl -X POST http://localhost:3001/print-html \
  -H "Content-Type: application/json" \
  -d '{
    "html": "<html><body><h1>Test Print</h1></body></html>",
    "printerIP": "192.168.1.115",
    "printerPort": 9100,
    "paperWidth": 80
  }'
```

## Bước 6: Cấu Hình Backend

### 6.1. Lấy IP Của Máy Docker

```bash
# Trên máy Docker, lấy IP
ip addr show | grep "inet " | grep -v 127.0.0.1

# Hoặc
hostname -I
```

Ví dụ: `192.168.1.200`

### 6.2. Cấu Hình Backend

Trên máy backend, sửa file `.env`:

```bash
# Print Bridge URL (IP của máy Docker)
PRINT_BRIDGE_URL=http://192.168.1.200:3001
```

### 6.3. Test Từ Backend

```bash
# Từ máy backend, test kết nối
curl http://192.168.1.200:3001/health

# Test print
curl -X POST http://192.168.1.200:3001/test-connection \
  -H "Content-Type: application/json" \
  -d '{
    "printerIP": "192.168.1.115",
    "printerPort": 9100
  }'
```

## Bước 7: Cấu Hình Auto-Start

### 7.1. Docker Compose Auto-Restart

File `docker-compose.yml` đã có:
```yaml
restart: unless-stopped
```

Service sẽ tự động restart khi:
- Container crash
- Máy reboot
- Docker daemon restart

### 7.2. Systemd Service (Optional)

Tạo file `/etc/systemd/system/print-bridge.service`:

```bash
sudo nano /etc/systemd/system/print-bridge.service
```

Nội dung:
```ini
[Unit]
Description=Local Print Bridge
Requires=docker.service
After=docker.service network-online.target
Wants=network-online.target

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/home/user/local-print-bridge
ExecStart=/usr/bin/docker compose up -d
ExecStop=/usr/bin/docker compose down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
```

Enable service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable print-bridge
sudo systemctl start print-bridge

# Check status
sudo systemctl status print-bridge
```

## Bước 8: Monitoring và Maintenance

### 8.1. Check Logs

```bash
# View logs
docker compose logs -f

# Last 100 lines
docker compose logs --tail=100

# Logs since 1 hour ago
docker compose logs --since 1h
```

### 8.2. Check Status

```bash
# Container status
docker compose ps

# Health check
curl http://localhost:3001/health

# Docker stats
docker stats local-print-bridge
```

### 8.3. Restart Service

```bash
# Restart
docker compose restart

# Stop
docker compose down

# Start
docker compose up -d
```

### 8.4. Update Service

```bash
# Pull latest code
git pull

# Rebuild
docker compose build --no-cache

# Restart
docker compose up -d

# Check logs
docker compose logs -f
```

## Troubleshooting

### Lỗi: Cannot connect to printer

```bash
# 1. Check network connectivity
ping 192.168.1.115

# 2. Check port
nc -zv 192.168.1.115 9100

# 3. Check from container
docker compose exec print-bridge ping 192.168.1.115
docker compose exec print-bridge nc -zv 192.168.1.115 9100

# 4. Check firewall
sudo ufw status
sudo ufw allow 9100/tcp
```

### Lỗi: Backend cannot connect to Print Bridge

```bash
# 1. Check Print Bridge is running
curl http://localhost:3001/health

# 2. Check from backend machine
curl http://192.168.1.200:3001/health

# 3. Check firewall on Docker machine
sudo ufw allow 3001/tcp

# 4. Check Docker network
docker compose exec print-bridge ip addr
```

### Lỗi: Container keeps restarting

```bash
# Check logs
docker compose logs

# Common issues:
# - Port 3001 already in use
# - Invalid .env configuration
# - Chromium dependencies missing

# Check port
sudo lsof -i :3001
sudo netstat -tulpn | grep 3001
```

### Lỗi: High memory usage

```bash
# Check memory
docker stats local-print-bridge

# Increase shared memory in docker-compose.yml:
shm_size: '512mb'

# Restart
docker compose down
docker compose up -d
```

### Lỗi: Chromium errors

```bash
# Check Chromium
docker compose exec print-bridge which chromium-browser
docker compose exec print-bridge chromium-browser --version

# Rebuild with no cache
docker compose build --no-cache
docker compose up -d
```

## Network Configurations

### Cấu Hình 1: Same LAN (Khuyến nghị)

```
Backend (192.168.1.100) ──┐
                          ├─── LAN ─── Print Bridge (192.168.1.200) ─── Printer (192.168.1.115)
Client Browser ───────────┘
```

- Backend config: `PRINT_BRIDGE_URL=http://192.168.1.200:3001`
- Đơn giản, tốc độ nhanh
- Không cần cấu hình phức tạp

### Cấu Hình 2: Different Networks

```
Backend (Public IP) ─── Internet ─── Router ─── Print Bridge (192.168.1.200) ─── Printer
```

Cần:
- Port forwarding trên router: `public_ip:3001 → 192.168.1.200:3001`
- Hoặc dùng VPN/Tunnel (xem phần Cloudflare Tunnel)

### Cấu Hình 3: Cloudflare Tunnel (Remote Access)

Nếu backend và print bridge ở khác mạng, dùng Cloudflare Tunnel:

```bash
# Xem file CLOUDFLARE_TUNNEL_SETUP.md
cd local-print-bridge
./cloudflare-tunnel-setup.sh
```

## Security Best Practices

### 1. Firewall Rules

```bash
# Chỉ cho phép backend access
sudo ufw allow from 192.168.1.100 to any port 3001

# Hoặc cho phép cả subnet
sudo ufw allow from 192.168.1.0/24 to any port 3001
```

### 2. HTTPS (Optional)

Dùng reverse proxy như Nginx:

```nginx
server {
    listen 443 ssl;
    server_name print-bridge.local;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://localhost:3001;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 3. Authentication (Optional)

Thêm API key vào Print Bridge nếu cần.

## Performance Tuning

### 1. Resource Limits

Edit `docker-compose.yml`:

```yaml
services:
  print-bridge:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

### 2. Logging

Giới hạn log size:

```yaml
services:
  print-bridge:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## Backup và Recovery

### Backup Configuration

```bash
# Backup .env và docker-compose.yml
tar -czf print-bridge-config-backup.tar.gz .env docker-compose.yml

# Backup Docker image
docker save local-print-bridge:latest | gzip > print-bridge-image.tar.gz
```

### Restore

```bash
# Restore config
tar -xzf print-bridge-config-backup.tar.gz

# Restore image
docker load < print-bridge-image.tar.gz

# Start service
docker compose up -d
```

## Checklist Deploy

- [ ] Docker và Docker Compose đã cài đặt
- [ ] Code đã copy sang máy Docker
- [ ] File .env đã cấu hình đúng IP máy in
- [ ] Máy Docker ping được máy in
- [ ] Port 9100 của máy in accessible
- [ ] Docker image build thành công
- [ ] Container chạy và healthy
- [ ] Health check endpoint trả về OK
- [ ] Test connection với máy in thành công
- [ ] Backend có thể access Print Bridge
- [ ] Backend .env đã config PRINT_BRIDGE_URL
- [ ] Test print từ backend thành công
- [ ] Auto-restart đã cấu hình
- [ ] Firewall rules đã setup
- [ ] Monitoring và logs đã check

## Tài Liệu Tham Khảo

- [DOCKER_GUIDE.md](./DOCKER_GUIDE.md) - Chi tiết về Docker deployment
- [CLOUDFLARE_TUNNEL_SETUP.md](../CLOUDFLARE_TUNNEL_SETUP.md) - Setup remote access
- [PRINT_BRIDGE_COMPLETE_GUIDE.md](../PRINT_BRIDGE_COMPLETE_GUIDE.md) - Hướng dẫn tổng quan

## Support

Nếu gặp vấn đề:
1. Check logs: `docker compose logs -f`
2. Check health: `curl http://localhost:3001/health`
3. Check network: `ping` và `nc -zv`
4. Check firewall: `sudo ufw status`
5. Rebuild: `docker compose build --no-cache && docker compose up -d`

---

**Version:** 1.0.0  
**Last Updated:** March 2026
