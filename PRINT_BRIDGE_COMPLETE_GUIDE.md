# Print Bridge - Complete Setup Guide

## Tổng quan

Hướng dẫn đầy đủ để setup Print Bridge kết nối với Backend qua WebSocket.

## Kiến trúc

```
Backend (EC2) ←→ WebSocket ←→ Print Bridge (Local) ←→ ESC/POS ←→ Printers
```

## Phần 1: Backend Setup (EC2)

### 1.1 Fix Nginx WebSocket Proxy

```bash
# Trên máy dev
./fix-websocket-ec2.sh
```

Script này sẽ:
- Rebuild frontend với nginx config mới (có `/socket.io/` proxy)
- Push image lên Docker Hub

### 1.2 Deploy lên EC2

```bash
# SSH vào EC2
ssh -i your-key.pem ubuntu@tacafe.store

# Pull và restart
docker pull linhtranphu/cafe-pos-frontend:latest
docker-compose -f docker-compose.prod.yml up -d --force-recreate frontend

# Verify
docker logs -f cafe-pos-frontend
docker logs -f cafe-pos-backend | grep Socket.IO
```

**Expected logs:**
```
[Socket.IO] Server initialized
[Socket.IO] Client connected: <socket-id>
```

## Phần 2: Print Bridge Setup (Local Machine)

### Option A: Docker Deployment (Khuyến nghị)

**Ưu điểm:**
- ✅ Dễ deploy và maintain
- ✅ Auto-restart khi máy khởi động lại
- ✅ Isolated environment
- ✅ Dễ backup và restore

**Các bước:**

#### 1. Cài Docker

**Windows:**
- Download Docker Desktop: https://www.docker.com/products/docker-desktop
- Install và restart máy
- Verify: `docker --version`

**Linux:**
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

**macOS:**
- Download Docker Desktop: https://www.docker.com/products/docker-desktop
- Install và start Docker Desktop

#### 2. Clone Repository

```bash
git clone <repo-url>
cd local-print-bridge
```

#### 3. Cấu hình

```bash
# Copy template
cp .env.docker .env

# Edit với thông tin của bạn
nano .env  # hoặc notepad .env trên Windows
```

File `.env`:
```bash
PORT=3001
BACKEND_URL=https://tacafe.store
DEFAULT_BILL_PRINTER_IP=192.168.1.100
DEFAULT_BILL_PRINTER_PORT=9100
DEFAULT_LABEL_PRINTER_IP=192.168.1.101
DEFAULT_LABEL_PRINTER_PORT=9100
LOG_LEVEL=info
PRINTER_TIMEOUT=5000
```

#### 4. Deploy

```bash
# Sử dụng script tự động
./docker-deploy.sh

# Chọn option 1: Build and start
```

Hoặc manual:
```bash
docker-compose build
docker-compose up -d
```

#### 5. Verify

```bash
# Check container
docker-compose ps

# Check logs
docker-compose logs -f

# Test health
curl http://localhost:3001/health

# Test WebSocket
docker-compose logs | grep "WebSocket.*Connected"
```

**Expected output:**
```
[WebSocket] Connecting to: https://tacafe.store
[WebSocket] ✅ Connected to backend
Ready to accept print requests!
```

### Option B: Node.js Deployment (Development)

**Ưu điểm:**
- ✅ Dễ debug
- ✅ Dễ modify code
- ✅ Không cần Docker

**Các bước:**

#### 1. Cài Node.js

Download từ: https://nodejs.org/ (LTS version)

Verify:
```bash
node --version  # Should be v16+
npm --version
```

#### 2. Clone và Install

```bash
git clone <repo-url>
cd local-print-bridge
npm install
```

#### 3. Cấu hình

```bash
cp .env.example .env
# Edit .env với thông tin của bạn
```

#### 4. Start

**Quick start:**
```bash
./quick-start.sh
```

**Manual start:**
```bash
npm start
```

**With PM2 (auto-restart):**
```bash
npm install -g pm2
pm2 start src/index.js --name print-bridge
pm2 save
pm2 startup  # Auto-start on boot
```

#### 5. Verify

```bash
# Test health
curl http://localhost:3001/health

# Test WebSocket
node test-backend-websocket.js https://tacafe.store
```

## Phần 3: Testing End-to-End

### 3.1 Test WebSocket Connection

**Từ Print Bridge:**
```bash
# Docker
docker-compose logs | grep WebSocket

# Node.js
# Check console logs
```

**Expected:**
```
[WebSocket] ✅ Connected to backend
```

### 3.2 Test Printer Connection

```bash
curl -X POST http://localhost:3001/test-connection \
  -H "Content-Type: application/json" \
  -d '{"printerIP": "192.168.1.100", "printerPort": 9100}'
```

**Expected:**
```json
{
  "success": true,
  "message": "Printer connection successful",
  "printer": "192.168.1.100:9100"
}
```

### 3.3 Test Print Flow

1. Mở https://tacafe.store
2. Login vào hệ thống
3. Tạo order mới
4. Kiểm tra Print Bridge logs:

```bash
# Docker
docker-compose logs -f

# Node.js
# Check console
```

**Expected logs:**
```
[WebSocket] 📨 New print job received: { job: { id: '...' } }
[PrintJobHandler] Processing job <id>
[PrintJobHandler] ✅ Job <id> printed successfully
[PrintJobHandler] ✅ Backend updated - Job <id> -> COMPLETED
```

5. Kiểm tra máy in có in ra không
6. Kiểm tra status trong UI

## Phần 4: Troubleshooting

### 4.1 WebSocket không kết nối

**Triệu chứng:**
```
[WebSocket] Connection error: timeout
```

**Giải pháp:**

1. Check backend reachable:
```bash
curl https://tacafe.store/api/state-machines
```

2. Check nginx có `/socket.io/` proxy:
```bash
# Trên EC2
grep -A 10 "location /socket.io/" /path/to/nginx.conf
```

3. Check backend logs:
```bash
docker logs cafe-pos-backend | grep Socket.IO
```

4. Check firewall/security group:
- EC2 Security Group: Port 3000 open
- Local firewall: Allow outbound to port 3000

### 4.2 Máy in không in

**Triệu chứng:**
```
[PrintJobHandler] ❌ Job failed: connect ETIMEDOUT
```

**Giải pháp:**

1. Check printer IP:
```bash
ping 192.168.1.100
```

2. Check printer port:
```bash
telnet 192.168.1.100 9100
# hoặc
nc -zv 192.168.1.100 9100
```

3. Check printer power và network cable

4. Print test page từ printer để xác nhận IP

5. Check network mode (Docker):
```yaml
# docker-compose.yml MUST have:
network_mode: host
```

### 4.3 Backend không nhận status update

**Triệu chứng:**
```
[PrintJobHandler] Failed to update backend
```

**Giải pháp:**

1. Check BACKEND_URL:
```bash
cat .env | grep BACKEND_URL
```

2. Test backend API:
```bash
curl -X PUT https://tacafe.store/api/print-jobs/test-id/status \
  -H "Content-Type: application/json" \
  -d '{"status": "COMPLETED"}'
```

3. Check backend logs:
```bash
docker logs cafe-pos-backend | grep "print-jobs"
```

## Phần 5: Monitoring

### 5.1 Health Check

```bash
# Print Bridge
curl http://localhost:3001/health

# Backend
curl https://tacafe.store/api/state-machines
```

### 5.2 Status Check

```bash
curl http://localhost:3001/status
```

Response:
```json
{
  "success": true,
  "stats": {
    "totalJobs": 150,
    "successfulJobs": 148,
    "failedJobs": 2
  },
  "uptime": 86400
}
```

### 5.3 Logs

**Docker:**
```bash
docker-compose logs -f
docker-compose logs --tail=100
```

**Node.js:**
```bash
# PM2
pm2 logs print-bridge

# Manual
# Logs in console
```

**Backend:**
```bash
docker logs -f cafe-pos-backend
```

## Phần 6: Maintenance

### 6.1 Update Print Bridge

**Docker:**
```bash
git pull
docker-compose build --no-cache
docker-compose up -d
```

**Node.js:**
```bash
git pull
npm install
pm2 restart print-bridge  # if using PM2
```

### 6.2 Backup Configuration

```bash
# Backup .env
cp .env .env.backup

# Backup logs
tar -czf logs-backup-$(date +%Y%m%d).tar.gz logs/
```

### 6.3 Update Printer IPs

```bash
# Edit .env
nano .env

# Restart
docker-compose restart  # Docker
pm2 restart print-bridge  # PM2
```

## Checklist

### Backend (EC2)
- [ ] Nginx có `/socket.io/` proxy config
- [ ] Frontend image đã rebuild và deploy
- [ ] Backend logs show Socket.IO initialized
- [ ] Security group allows port 3000

### Print Bridge (Local)
- [ ] Docker hoặc Node.js đã cài
- [ ] `.env` đã cấu hình đúng
- [ ] Service đang chạy
- [ ] WebSocket connected to backend
- [ ] Có thể ping được printers

### Testing
- [ ] Health check OK
- [ ] WebSocket connection OK
- [ ] Printer connection OK
- [ ] End-to-end print flow OK
- [ ] Status update về backend OK

## Related Documentation

- `PRINT_ARCHITECTURE.md` - Kiến trúc tổng quan
- `PRINT_BRIDGE_WEBSOCKET_SETUP.md` - Setup chi tiết
- `WEBSOCKET_EC2_FIX.md` - Fix WebSocket trên EC2
- `DEPLOYMENT_CHECKLIST.md` - Checklist deploy
- `local-print-bridge/DOCKER_DEPLOYMENT.md` - Docker guide
- `local-print-bridge/DOCKER_QUICK_START.md` - Docker quick start
- `local-print-bridge/README.md` - Print Bridge overview

## Support

Nếu gặp vấn đề:
1. Check logs (backend + print bridge)
2. Test network connectivity
3. Review troubleshooting section
4. Check related documentation

## Summary

**Backend:** WebSocket server on port 3000 with `/socket.io/` nginx proxy
**Print Bridge:** Docker hoặc Node.js, kết nối backend qua WebSocket
**Flow:** Backend → WebSocket → Print Bridge → ESC/POS → Printer

Hệ thống hoạt động realtime, auto-reconnect, và không có CORS issues! 🎉
