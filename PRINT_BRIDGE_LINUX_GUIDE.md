# Hướng Dẫn Fix Print Bridge Trên Linux

## Vấn Đề

Print Bridge chạy OK trên macOS nhưng không in được trên Linux:
- ✓ Kiểm tra kết nối: Thành công
- ✗ Bấm "In thử": Không in được

## Nguyên Nhân

Chromium trên Linux cần thêm:
1. Flags đặc biệt (`--no-sandbox`, `--disable-dev-shm-usage`)
2. Shared memory lớn hơn (512MB thay vì 256MB)
3. Capabilities (`SYS_ADMIN`)
4. Dependencies bổ sung (nss, freetype, harfbuzz)

## Giải Pháp Nhanh

### Trên Máy Linux

```bash
# 1. Copy script fix sang máy Linux
scp fix-print-bridge-linux.sh user@linux-machine:/path/to/project/

# 2. SSH vào máy Linux
ssh user@linux-machine

# 3. Chạy script fix
cd /path/to/project
./fix-print-bridge-linux.sh

# 4. Kiểm tra logs
cd local-print-bridge
docker logs -f local-print-bridge

# 5. Test in thử từ web interface
```

## Giải Pháp Thủ Công

### Bước 1: Update docker-compose.yml

```yaml
version: '3.8'

services:
  print-bridge:
    build:
      context: .
      dockerfile: Dockerfile
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
      # Chromium environment
      - CHROME_BIN=/usr/bin/chromium-browser
      - CHROME_PATH=/usr/bin/chromium-browser
    restart: unless-stopped
    # QUAN TRỌNG: Security options cho Chromium
    security_opt:
      - seccomp:unconfined
    # QUAN TRỌNG: Tăng shared memory
    shm_size: '512mb'
    # QUAN TRỌNG: Capabilities cho Chromium
    cap_add:
      - SYS_ADMIN
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:3001/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
```

### Bước 2: Update Dockerfile

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o print-bridge main.go

# Runtime stage
FROM alpine:latest

# Install Chromium và dependencies đầy đủ cho Linux
RUN apk --no-cache add \
    chromium \
    chromium-chromedriver \
    font-noto \
    font-noto-cjk \
    ttf-dejavu \
    fontconfig \
    ca-certificates \
    # Dependencies bổ sung cho Linux
    nss \
    freetype \
    harfbuzz \
    ttf-freefont \
    # Tools để debug
    wget \
    curl

RUN fc-cache -f

# Environment variables
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/bin/chromium-browser
ENV CHROMIUM_FLAGS="--no-sandbox --disable-dev-shm-usage --disable-gpu --headless"

WORKDIR /root/

COPY --from=builder /app/print-bridge .

EXPOSE 3001

CMD ["./print-bridge"]
```

### Bước 3: Rebuild

```bash
cd local-print-bridge

# Stop container
docker compose down

# Rebuild without cache
docker compose build --no-cache

# Start
docker compose up -d

# Check logs
docker logs -f local-print-bridge
```

## Kiểm Tra

### 1. Check Container Status

```bash
docker ps | grep print-bridge
```

Kết quả mong đợi:
```
local-print-bridge   Up X minutes (healthy)   0.0.0.0:3001->3001/tcp
```

### 2. Check Health

```bash
curl http://localhost:3001/health
```

Kết quả mong đợi:
```json
{"status":"ok","timestamp":"..."}
```

### 3. Test Chromium Trong Container

```bash
docker exec local-print-bridge chromium-browser \
  --headless \
  --disable-gpu \
  --no-sandbox \
  --disable-dev-shm-usage \
  --dump-dom \
  about:blank
```

Kết quả mong đợi: HTML output

### 4. Test Print Endpoint

```bash
curl -X POST http://localhost:3001/print-html \
  -H "Content-Type: application/json" \
  -d '{
    "html": "<html><body><h1>Test Print</h1></body></html>",
    "printerIP": "192.168.1.100",
    "printerPort": 9100,
    "paperWidth": 80
  }'
```

### 5. Check Logs

```bash
docker logs --tail=50 local-print-bridge
```

Tìm errors liên quan đến Chromium hoặc rendering.

## Troubleshooting

### Lỗi: "Failed to move to new namespace"

**Nguyên nhân:** Thiếu `--no-sandbox` flag

**Giải pháp:**
```yaml
# Trong docker-compose.yml
security_opt:
  - seccomp:unconfined
cap_add:
  - SYS_ADMIN
```

### Lỗi: "Shared memory too small"

**Nguyên nhân:** Chromium cần nhiều shared memory

**Giải pháp:**
```yaml
# Trong docker-compose.yml
shm_size: '512mb'  # Hoặc '1gb' nếu cần
```

### Lỗi: "Cannot find Chrome binary"

**Nguyên nhân:** Chromium không được cài đặt đúng

**Giải pháp:**
```bash
# Rebuild container
docker compose build --no-cache
docker compose up -d

# Verify Chromium
docker exec local-print-bridge which chromium-browser
docker exec local-print-bridge chromium-browser --version
```

### Lỗi: "Font rendering issues"

**Nguyên nhân:** Thiếu fonts

**Giải pháp:**
```dockerfile
# Trong Dockerfile, thêm fonts
RUN apk --no-cache add \
    font-noto \
    font-noto-cjk \
    ttf-dejavu \
    ttf-freefont \
    fontconfig

RUN fc-cache -f
```

### Container Crash Liên Tục

**Kiểm tra:**
```bash
# Check logs
docker logs local-print-bridge

# Check container status
docker inspect local-print-bridge

# Try running manually
docker run -it --rm \
  -p 3001:3001 \
  --security-opt seccomp:unconfined \
  --cap-add SYS_ADMIN \
  --shm-size 512m \
  local-print-bridge-print-bridge
```

## So Sánh: macOS vs Linux

### macOS
- Chromium chạy OK với cấu hình mặc định
- Không cần `--no-sandbox`
- Shared memory 256MB đủ
- Không cần capabilities đặc biệt

### Linux
- Cần `--no-sandbox` flag
- Cần `SYS_ADMIN` capability
- Cần shared memory 512MB+
- Cần `seccomp:unconfined`
- Cần thêm dependencies (nss, freetype, harfbuzz)

## Scripts Hữu Ích

### Chẩn Đoán

```bash
# Chạy script chẩn đoán
./diagnose-print-bridge-linux.sh
```

Script sẽ kiểm tra:
- OS type
- Docker container status
- Chromium installation
- Chromium execution test
- Shared memory
- Environment variables
- Recent logs
- Print endpoint test

### Fix Tự Động

```bash
# Chạy script fix
./fix-print-bridge-linux.sh
```

Script sẽ:
- Backup files hiện tại
- Update docker-compose.yml
- Update Dockerfile
- Rebuild container
- Test Chromium
- Show logs

## Checklist Deploy Trên Linux

- [ ] Docker và Docker Compose đã cài đặt
- [ ] Code đã copy sang máy Linux
- [ ] File .env đã cấu hình printer IPs
- [ ] docker-compose.yml có `shm_size: '512mb'`
- [ ] docker-compose.yml có `security_opt: seccomp:unconfined`
- [ ] docker-compose.yml có `cap_add: SYS_ADMIN`
- [ ] Dockerfile có đầy đủ dependencies (nss, freetype, harfbuzz)
- [ ] Container build thành công
- [ ] Container chạy và healthy
- [ ] Health check trả về OK
- [ ] Chromium test thành công
- [ ] Print endpoint test thành công
- [ ] Test in thử từ web interface thành công

## Best Practices

### 1. Luôn Test Chromium Sau Khi Deploy

```bash
docker exec local-print-bridge chromium-browser \
  --headless --no-sandbox --disable-gpu \
  --dump-dom about:blank | head -5
```

### 2. Monitor Logs

```bash
# Real-time logs
docker logs -f local-print-bridge

# Filter for errors
docker logs local-print-bridge 2>&1 | grep -i error

# Filter for Chromium issues
docker logs local-print-bridge 2>&1 | grep -i chromium
```

### 3. Resource Monitoring

```bash
# Check memory usage
docker stats local-print-bridge

# Check if hitting shared memory limit
docker exec local-print-bridge df -h /dev/shm
```

### 4. Regular Updates

```bash
# Pull latest code
git pull

# Rebuild
cd local-print-bridge
docker compose build --no-cache
docker compose up -d

# Verify
curl http://localhost:3001/health
```

## Production Deployment

### Systemd Service

```ini
# /etc/systemd/system/print-bridge.service
[Unit]
Description=Local Print Bridge
Requires=docker.service
After=docker.service network-online.target
Wants=network-online.target

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/local-print-bridge
ExecStart=/usr/bin/docker compose up -d
ExecStop=/usr/bin/docker compose down
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable:
```bash
sudo systemctl daemon-reload
sudo systemctl enable print-bridge
sudo systemctl start print-bridge
sudo systemctl status print-bridge
```

### Auto-Restart on Failure

Already configured in docker-compose.yml:
```yaml
restart: unless-stopped
```

### Monitoring

```bash
# Check if service is running
systemctl status print-bridge

# Check container health
docker ps | grep print-bridge

# Check logs
journalctl -u print-bridge -f
docker logs -f local-print-bridge
```

## Tài Liệu Tham Khảo

- [DEPLOY_DOCKER_REMOTE.md](local-print-bridge/DEPLOY_DOCKER_REMOTE.md) - Deploy trên máy khác
- [DOCKERHUB_GUIDE.md](local-print-bridge/DOCKERHUB_GUIDE.md) - Sử dụng DockerHub
- [PRINT_BRIDGE_TROUBLESHOOTING.md](PRINT_BRIDGE_TROUBLESHOOTING.md) - Troubleshooting chung

## Support

Nếu vẫn gặp vấn đề:

1. Chạy diagnostic:
   ```bash
   ./diagnose-print-bridge-linux.sh > diagnostic.txt
   ```

2. Collect logs:
   ```bash
   docker logs local-print-bridge > print-bridge.log
   ```

3. Gửi kèm:
   - diagnostic.txt
   - print-bridge.log
   - docker-compose.yml
   - Dockerfile

---

**Version:** 1.0.0  
**Last Updated:** March 2026
