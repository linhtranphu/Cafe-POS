# Hướng Dẫn Sử Dụng Docker Image Từ DockerHub

## Phần 1: Push Image Lên DockerHub (Dành Cho Developer)

### 1.1. Tạo Tài Khoản DockerHub

1. Truy cập https://hub.docker.com
2. Đăng ký tài khoản miễn phí
3. Ghi nhớ username của bạn (ví dụ: `yourusername`)

### 1.2. Login Docker CLI

```bash
# Login vào DockerHub
docker login

# Nhập username và password khi được hỏi
```

### 1.3. Build và Tag Image

```bash
cd local-print-bridge

# Build image với tag
docker build -t yourusername/local-print-bridge:latest .

# Hoặc build với version cụ thể
docker build -t yourusername/local-print-bridge:v1.0.0 .

# Build cả 2 tags
docker build -t yourusername/local-print-bridge:latest \
             -t yourusername/local-print-bridge:v1.0.0 .
```

### 1.4. Push Image Lên DockerHub

```bash
# Push latest tag
docker push yourusername/local-print-bridge:latest

# Push version tag
docker push yourusername/local-print-bridge:v1.0.0

# Push tất cả tags
docker push yourusername/local-print-bridge --all-tags
```

### 1.5. Verify Image Trên DockerHub

1. Truy cập https://hub.docker.com/r/yourusername/local-print-bridge
2. Kiểm tra image đã được push thành công
3. Xem tags và size của image

---

## Phần 2: Pull và Sử Dụng Image (Dành Cho User)

### 2.1. Chuẩn Bị Máy Docker

**Chỉ cần cài Docker, KHÔNG cần code!**

```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
newgrp docker

# macOS
brew install --cask docker

# Windows
# Download Docker Desktop từ https://www.docker.com/products/docker-desktop
```

### 2.2. Tạo Thư Mục Làm Việc

```bash
# Tạo thư mục
mkdir -p ~/print-bridge
cd ~/print-bridge
```

### 2.3. Tạo File docker-compose.yml

```bash
# Tạo file docker-compose.yml
nano docker-compose.yml
```

Nội dung:

```yaml
version: '3.8'

services:
  print-bridge:
    image: yourusername/local-print-bridge:latest
    container_name: local-print-bridge
    ports:
      - "3001:3001"
    environment:
      - PORT=3001
      - HOST=0.0.0.0
      - DEFAULT_BILL_PRINTER_IP=${DEFAULT_BILL_PRINTER_IP:-192.168.1.115}
      - DEFAULT_BILL_PRINTER_PORT=${DEFAULT_BILL_PRINTER_PORT:-9100}
      - DEFAULT_LABEL_PRINTER_IP=${DEFAULT_LABEL_PRINTER_IP:-192.168.1.116}
      - DEFAULT_LABEL_PRINTER_PORT=${DEFAULT_LABEL_PRINTER_PORT:-9100}
      - LOG_LEVEL=${LOG_LEVEL:-info}
      - PRINTER_TIMEOUT=${PRINTER_TIMEOUT:-5}
    restart: unless-stopped
    security_opt:
      - seccomp:unconfined
    shm_size: '256mb'
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:3001/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
```

### 2.4. Tạo File .env

```bash
# Tạo file .env
nano .env
```

Nội dung:

```bash
# Printer Configuration
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

### 2.5. Pull và Start Service

```bash
# Pull image từ DockerHub
docker compose pull

# Start service
docker compose up -d

# Check logs
docker compose logs -f

# Check status
docker compose ps
```

### 2.6. Verify Service

```bash
# Health check
curl http://localhost:3001/health

# Expected response:
# {"status":"ok","timestamp":"..."}

# Test printer connection
curl -X POST http://localhost:3001/test-connection \
  -H "Content-Type: application/json" \
  -d '{
    "printerIP": "192.168.1.115",
    "printerPort": 9100
  }'
```

---

## Phần 3: Quick Start Script

### 3.1. Tạo Script Tự Động

```bash
# Tạo file setup.sh
nano setup.sh
```

Nội dung:

```bash
#!/bin/bash

echo "=== Local Print Bridge Setup ==="
echo ""

# Tạo thư mục
mkdir -p ~/print-bridge
cd ~/print-bridge

# Tạo docker-compose.yml
cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  print-bridge:
    image: yourusername/local-print-bridge:latest
    container_name: local-print-bridge
    ports:
      - "3001:3001"
    environment:
      - PORT=3001
      - HOST=0.0.0.0
      - DEFAULT_BILL_PRINTER_IP=${DEFAULT_BILL_PRINTER_IP:-192.168.1.115}
      - DEFAULT_BILL_PRINTER_PORT=${DEFAULT_BILL_PRINTER_PORT:-9100}
      - DEFAULT_LABEL_PRINTER_IP=${DEFAULT_LABEL_PRINTER_IP:-192.168.1.116}
      - DEFAULT_LABEL_PRINTER_PORT=${DEFAULT_LABEL_PRINTER_PORT:-9100}
      - LOG_LEVEL=${LOG_LEVEL:-info}
      - PRINTER_TIMEOUT=${PRINTER_TIMEOUT:-5}
    restart: unless-stopped
    security_opt:
      - seccomp:unconfined
    shm_size: '256mb'
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:3001/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
EOF

# Tạo .env
cat > .env << 'EOF'
# Printer Configuration
DEFAULT_BILL_PRINTER_IP=192.168.1.115
DEFAULT_BILL_PRINTER_PORT=9100

DEFAULT_LABEL_PRINTER_IP=192.168.1.116
DEFAULT_LABEL_PRINTER_PORT=9100

# Logging
LOG_LEVEL=info

# Timeout
PRINTER_TIMEOUT=5
EOF

echo "✓ Files created"
echo ""
echo "Please edit .env file to configure your printer IPs:"
echo "  nano .env"
echo ""
echo "Then run:"
echo "  docker compose pull"
echo "  docker compose up -d"
echo ""
```

### 3.2. Chạy Setup Script

```bash
# Make executable
chmod +x setup.sh

# Run
./setup.sh

# Edit .env với printer IPs của bạn
nano ~/print-bridge/.env

# Start service
cd ~/print-bridge
docker compose pull
docker compose up -d
```

---

## Phần 4: Update Image

### 4.1. Update Lên Version Mới

```bash
cd ~/print-bridge

# Pull latest image
docker compose pull

# Restart với image mới
docker compose up -d

# Check logs
docker compose logs -f
```

### 4.2. Rollback Về Version Cũ

```bash
# Edit docker-compose.yml, thay đổi tag
# image: yourusername/local-print-bridge:v1.0.0

# Pull version cũ
docker compose pull

# Restart
docker compose up -d
```

### 4.3. Auto-Update Script

```bash
# Tạo update script
nano update.sh
```

Nội dung:

```bash
#!/bin/bash

echo "=== Updating Print Bridge ==="

cd ~/print-bridge

# Backup current version
docker tag local-print-bridge local-print-bridge:backup

# Pull latest
docker compose pull

# Restart
docker compose up -d

# Check health
sleep 5
curl -s http://localhost:3001/health

echo ""
echo "✓ Update complete"
```

---

## Phần 5: Commands Thường Dùng

### Start/Stop Service

```bash
cd ~/print-bridge

# Start
docker compose up -d

# Stop
docker compose down

# Restart
docker compose restart

# View logs
docker compose logs -f

# View status
docker compose ps
```

### Check Health

```bash
# Health check
curl http://localhost:3001/health

# Test printer
curl -X POST http://localhost:3001/test-connection \
  -H "Content-Type: application/json" \
  -d '{"printerIP": "192.168.1.115", "printerPort": 9100}'

# Docker stats
docker stats local-print-bridge
```

### Troubleshooting

```bash
# View logs
docker compose logs --tail=100

# Check container
docker compose ps

# Restart service
docker compose restart

# Rebuild (pull latest)
docker compose pull
docker compose up -d

# Check network
docker compose exec print-bridge ping 192.168.1.115
```

---

## Phần 6: Cấu Hình Backend

### 6.1. Lấy IP Máy Docker

```bash
# Lấy IP
ip addr show | grep "inet " | grep -v 127.0.0.1

# Hoặc
hostname -I
```

Ví dụ: `192.168.1.200`

### 6.2. Config Backend

Sửa file `.env` của backend:

```bash
# Print Bridge URL
PRINT_BRIDGE_URL=http://192.168.1.200:3001
```

### 6.3. Test Từ Backend

```bash
# Test connection
curl http://192.168.1.200:3001/health

# Test print
curl -X POST http://192.168.1.200:3001/test-connection \
  -H "Content-Type: application/json" \
  -d '{"printerIP": "192.168.1.115", "printerPort": 9100}'
```

---

## Phần 7: Multi-Architecture Support

### 7.1. Build Multi-Arch Image (Developer)

```bash
# Setup buildx
docker buildx create --name multiarch --use
docker buildx inspect --bootstrap

# Build cho nhiều platforms
docker buildx build \
  --platform linux/amd64,linux/arm64,linux/arm/v7 \
  -t yourusername/local-print-bridge:latest \
  --push \
  .

# Verify
docker buildx imagetools inspect yourusername/local-print-bridge:latest
```

### 7.2. Pull Đúng Architecture (User)

Docker tự động pull đúng architecture:

```bash
# Trên x86_64 (Intel/AMD)
docker compose pull  # Tự động pull amd64

# Trên ARM (Raspberry Pi, M1 Mac)
docker compose pull  # Tự động pull arm64

# Check architecture
docker inspect yourusername/local-print-bridge:latest | grep Architecture
```

---

## Phần 8: Private Registry (Optional)

### 8.1. Sử Dụng Private DockerHub Repository

```bash
# User cần login
docker login

# Pull private image
docker compose pull
```

### 8.2. Sử Dụng Self-Hosted Registry

Edit `docker-compose.yml`:

```yaml
services:
  print-bridge:
    image: registry.yourcompany.com/local-print-bridge:latest
    # ... rest of config
```

Pull:

```bash
# Login vào registry
docker login registry.yourcompany.com

# Pull
docker compose pull
```

---

## Phần 9: Offline Deployment

### 9.1. Save Image (Máy Có Internet)

```bash
# Pull image
docker pull yourusername/local-print-bridge:latest

# Save to file
docker save yourusername/local-print-bridge:latest | gzip > print-bridge.tar.gz

# Copy file này sang máy offline (USB, network share, etc.)
```

### 9.2. Load Image (Máy Offline)

```bash
# Load image
docker load < print-bridge.tar.gz

# Verify
docker images | grep print-bridge

# Start service (dùng docker-compose.yml như bình thường)
docker compose up -d
```

---

## Phần 10: Best Practices

### 10.1. Version Tagging

```bash
# Developer nên push cả latest và version tag
docker push yourusername/local-print-bridge:latest
docker push yourusername/local-print-bridge:v1.0.0

# User nên dùng version tag cho production
image: yourusername/local-print-bridge:v1.0.0  # Stable
# image: yourusername/local-print-bridge:latest  # Latest features
```

### 10.2. Image Size Optimization

Image đã được optimize:
- Multi-stage build
- Alpine base image
- Minimal dependencies
- Size: ~150MB (compressed)

### 10.3. Security

```bash
# Scan image for vulnerabilities
docker scan yourusername/local-print-bridge:latest

# Check for updates regularly
docker compose pull
```

---

## Phần 11: Troubleshooting

### Lỗi: Cannot pull image

```bash
# Check Docker login
docker login

# Check network
ping hub.docker.com

# Try with full image name
docker pull docker.io/yourusername/local-print-bridge:latest

# Check Docker Hub status
curl https://status.docker.com
```

### Lỗi: Wrong architecture

```bash
# Check your architecture
uname -m

# Pull specific platform
docker pull --platform linux/amd64 yourusername/local-print-bridge:latest

# Or edit docker-compose.yml
services:
  print-bridge:
    image: yourusername/local-print-bridge:latest
    platform: linux/amd64
```

### Lỗi: Image too large

```bash
# Check image size
docker images yourusername/local-print-bridge

# Clean up old images
docker image prune -a

# Use specific version instead of latest
docker pull yourusername/local-print-bridge:v1.0.0
```

---

## Phần 12: CI/CD Integration (Developer)

### 12.1. GitHub Actions

```yaml
# .github/workflows/docker-publish.yml
name: Docker Build and Push

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Login to DockerHub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      
      - name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: ./local-print-bridge
          push: true
          tags: |
            yourusername/local-print-bridge:latest
            yourusername/local-print-bridge:${{ github.ref_name }}
```

---

## So Sánh: Build Local vs Pull DockerHub

### Build Local (Cách Cũ)
```bash
# Cần:
- Clone code repository
- Cài đặt Git
- Build mất 5-10 phút
- Tốn bandwidth download dependencies
- Cần hiểu về Docker build

# Commands:
git clone <repo>
cd local-print-bridge
docker compose build
docker compose up -d
```

### Pull DockerHub (Cách Mới - Khuyến Nghị)
```bash
# Cần:
- Chỉ cần Docker
- Không cần code
- Pull mất 1-2 phút
- Image đã được optimize
- Đơn giản hơn nhiều

# Commands:
mkdir ~/print-bridge
cd ~/print-bridge
# Tạo docker-compose.yml và .env
docker compose pull
docker compose up -d
```

---

## Checklist Deploy Với DockerHub

- [ ] Docker đã cài đặt
- [ ] Tạo thư mục ~/print-bridge
- [ ] Tạo file docker-compose.yml
- [ ] Tạo file .env với printer IPs
- [ ] Pull image: `docker compose pull`
- [ ] Start service: `docker compose up -d`
- [ ] Check health: `curl http://localhost:3001/health`
- [ ] Test printer connection
- [ ] Config backend PRINT_BRIDGE_URL
- [ ] Test print từ backend

---

## Tài Liệu Tham Khảo

- [DEPLOY_DOCKER_REMOTE.md](./DEPLOY_DOCKER_REMOTE.md) - Deploy trên máy khác
- [DOCKER_GUIDE.md](./DOCKER_GUIDE.md) - Chi tiết Docker
- DockerHub: https://hub.docker.com/r/yourusername/local-print-bridge

---

**Version:** 1.0.0  
**Last Updated:** March 2026
