# Build và Push Print Bridge Lên DockerHub

## Quick Start

### 1. Build và Push (Trên Máy Dev)

```bash
cd local-print-bridge

# Login DockerHub
docker login

# Build và push (thay 'yourusername' bằng username DockerHub của bạn)
./publish-to-dockerhub.sh yourusername v1.0.0

# Script sẽ:
# - Build image với Chromium
# - Verify Chromium trong image
# - Push lên DockerHub với tag v1.0.0 và latest
```

### 2. Pull và Sử Dụng (Trên Máy Linux)

```bash
# Tạo thư mục
mkdir ~/print-bridge && cd ~/print-bridge

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
      - DEFAULT_BILL_PRINTER_IP=192.168.1.115
      - DEFAULT_BILL_PRINTER_PORT=9100
    restart: unless-stopped
    security_opt:
      - seccomp:unconfined
    shm_size: '512mb'
    cap_add:
      - SYS_ADMIN
EOF

# Pull và start
docker compose pull
docker compose up -d

# Check
docker logs -f local-print-bridge
```

Done! Không cần build trên máy Linux.

## Chi Tiết

### Build và Push

```bash
cd local-print-bridge

# Cách 1: Dùng script (khuyến nghị)
./publish-to-dockerhub.sh yourusername v1.0.0

# Cách 2: Manual
docker build -t yourusername/local-print-bridge:v1.0.0 .
docker tag yourusername/local-print-bridge:v1.0.0 yourusername/local-print-bridge:latest
docker push yourusername/local-print-bridge:v1.0.0
docker push yourusername/local-print-bridge:latest
```

### Verify Image

```bash
# Check image size
docker images yourusername/local-print-bridge

# Test image locally
docker run -d -p 3001:3001 yourusername/local-print-bridge:latest
curl http://localhost:3001/health
```

