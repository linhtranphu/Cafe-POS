# Hướng Dẫn Khắc Phục Lỗi Print Bridge

## Lỗi: "Print bridge is not available"

### Cách Kiểm Tra Nhanh

#### 1. Chạy Script Chẩn Đoán Tự Động

```bash
./diagnose-print-bridge.sh
```

Script này sẽ kiểm tra:
- ✓ Backend configuration (.env)
- ✓ Print Bridge service status
- ✓ Network connectivity
- ✓ Health endpoint
- ✓ Backend service
- ✓ Test endpoints
- ✓ Đưa ra khuyến nghị

#### 2. Kiểm Tra Thủ Công

**Bước 1: Kiểm tra Print Bridge có chạy không**

```bash
# Check Docker container
docker ps | grep print-bridge

# Hoặc
cd local-print-bridge
docker compose ps
```

Kết quả mong đợi:
```
local-print-bridge   running   0.0.0.0:3001->3001/tcp
```

**Bước 2: Test health endpoint**

```bash
# Nếu Print Bridge chạy trên cùng máy
curl http://localhost:3001/health

# Nếu Print Bridge chạy trên máy khác
curl http://192.168.1.X:3001/health
```

Kết quả mong đợi:
```json
{"status":"ok","timestamp":"2024-..."}
```

**Bước 3: Kiểm tra backend config**

```bash
# Check PRINT_BRIDGE_URL trong .env
grep PRINT_BRIDGE_URL .env
```

Phải có:
```bash
PRINT_BRIDGE_URL=http://localhost:3001
# Hoặc
PRINT_BRIDGE_URL=http://192.168.1.X:3001
```

---

## Các Nguyên Nhân Thường Gặp

### 1. Print Bridge Không Chạy

**Triệu chứng:**
```bash
$ docker ps | grep print-bridge
# Không có kết quả
```

**Giải pháp:**

```bash
# Start Print Bridge
cd local-print-bridge
docker compose up -d

# Check logs
docker compose logs -f

# Verify
curl http://localhost:3001/health
```

---

### 2. PRINT_BRIDGE_URL Không Đúng

**Triệu chứng:**
- Backend log: "connection refused"
- Health check fail

**Kiểm tra:**

```bash
# Check backend .env
grep PRINT_BRIDGE_URL .env

# Test URL
PRINT_BRIDGE_URL=$(grep PRINT_BRIDGE_URL .env | cut -d'=' -f2)
curl $PRINT_BRIDGE_URL/health
```

**Giải pháp:**

Sửa file `.env` của backend:

```bash
# Nếu Print Bridge chạy trên CÙNG máy với backend
PRINT_BRIDGE_URL=http://localhost:3001

# Nếu Print Bridge chạy trên máy KHÁC
# (Thay 192.168.1.200 bằng IP thực tế của máy Print Bridge)
PRINT_BRIDGE_URL=http://192.168.1.200:3001
```

Restart backend:
```bash
docker compose restart backend
```

---

### 3. Firewall Chặn Port 3001

**Triệu chứng:**
- Ping được host nhưng không connect được port
- `nc -zv host 3001` fail

**Kiểm tra:**

```bash
# Test port
nc -zv localhost 3001

# Hoặc từ máy backend
nc -zv 192.168.1.200 3001
```

**Giải pháp:**

Trên máy Print Bridge:

```bash
# Ubuntu/Debian
sudo ufw allow 3001/tcp
sudo ufw reload

# CentOS/RHEL
sudo firewall-cmd --add-port=3001/tcp --permanent
sudo firewall-cmd --reload

# macOS
# Thường không cần, macOS không có firewall mặc định
```

---

### 4. Print Bridge Crashed

**Triệu chứng:**
- Container tồn tại nhưng status là "Exited"
- Logs có errors

**Kiểm tra:**

```bash
cd local-print-bridge

# Check container status
docker compose ps

# Check logs
docker compose logs --tail=50
```

**Giải pháp:**

```bash
# Restart container
docker compose restart

# Nếu vẫn crash, rebuild
docker compose down
docker compose build --no-cache
docker compose up -d

# Check logs
docker compose logs -f
```

---

### 5. Network Không Kết Nối (Máy Khác Nhau)

**Triệu chứng:**
- Backend và Print Bridge ở 2 máy khác nhau
- Không ping được

**Kiểm tra:**

```bash
# Từ máy backend, ping máy Print Bridge
ping 192.168.1.200

# Check route
traceroute 192.168.1.200

# Check port
nc -zv 192.168.1.200 3001
```

**Giải pháp:**

1. Đảm bảo 2 máy trong cùng mạng LAN
2. Check IP address của máy Print Bridge:
   ```bash
   # Trên máy Print Bridge
   ip addr show
   # Hoặc
   hostname -I
   ```
3. Update PRINT_BRIDGE_URL trong backend .env
4. Check firewall trên cả 2 máy

---

### 6. Port 3001 Đã Được Sử Dụng

**Triệu chứng:**
- Container không start
- Logs: "port already in use"

**Kiểm tra:**

```bash
# Check port usage
sudo lsof -i :3001
# Hoặc
sudo netstat -tulpn | grep 3001
```

**Giải pháp:**

**Option 1: Kill process đang dùng port**
```bash
# Find PID
sudo lsof -i :3001

# Kill process
sudo kill -9 <PID>

# Restart Print Bridge
cd local-print-bridge
docker compose up -d
```

**Option 2: Đổi port**

Edit `local-print-bridge/.env`:
```bash
PORT=3002
```

Edit `local-print-bridge/docker-compose.yml`:
```yaml
ports:
  - "3002:3002"
```

Update backend `.env`:
```bash
PRINT_BRIDGE_URL=http://localhost:3002
```

Restart:
```bash
cd local-print-bridge
docker compose down
docker compose up -d
```

---

### 7. Docker Không Chạy

**Triệu chứng:**
```bash
$ docker ps
Cannot connect to the Docker daemon
```

**Giải pháp:**

```bash
# Start Docker
sudo systemctl start docker

# Enable auto-start
sudo systemctl enable docker

# Check status
sudo systemctl status docker
```

---

### 8. Chromium Errors

**Triệu chứng:**
- Container chạy nhưng không render được HTML
- Logs có Chromium errors

**Kiểm tra:**

```bash
cd local-print-bridge
docker compose logs | grep -i chromium
```

**Giải pháp:**

```bash
# Rebuild với no cache
docker compose down
docker compose build --no-cache
docker compose up -d

# Increase shared memory nếu cần
# Edit docker-compose.yml:
shm_size: '512mb'  # Tăng từ 256mb
```

---

## Checklist Kiểm Tra Đầy Đủ

### Trên Máy Print Bridge

- [ ] Docker đã cài đặt và chạy
- [ ] Print Bridge container đang chạy (`docker ps`)
- [ ] Health endpoint trả về OK (`curl http://localhost:3001/health`)
- [ ] Port 3001 không bị firewall chặn
- [ ] Logs không có errors (`docker compose logs`)

### Trên Máy Backend

- [ ] File .env có PRINT_BRIDGE_URL
- [ ] PRINT_BRIDGE_URL đúng (localhost hoặc IP máy Print Bridge)
- [ ] Có thể ping được máy Print Bridge (nếu khác máy)
- [ ] Có thể curl được health endpoint từ máy backend
- [ ] Backend đã restart sau khi sửa .env

### Network (Nếu 2 Máy Khác Nhau)

- [ ] 2 máy trong cùng mạng LAN
- [ ] Firewall không chặn port 3001
- [ ] IP address của Print Bridge đúng
- [ ] Có thể ping và nc được từ backend đến Print Bridge

---

## Commands Hữu Ích

### Kiểm Tra Status

```bash
# Print Bridge status
cd local-print-bridge && docker compose ps

# Backend status
docker compose ps backend

# Health check
curl http://localhost:3001/health

# Test connection
curl -X POST http://localhost:3001/test-connection \
  -H "Content-Type: application/json" \
  -d '{"printerIP": "192.168.1.100", "printerPort": 9100}'
```

### View Logs

```bash
# Print Bridge logs
cd local-print-bridge && docker compose logs -f

# Backend logs
docker compose logs -f backend

# Last 50 lines
docker compose logs --tail=50 backend
```

### Restart Services

```bash
# Restart Print Bridge
cd local-print-bridge && docker compose restart

# Restart Backend
docker compose restart backend

# Restart all
docker compose restart
```

### Network Debugging

```bash
# Ping
ping 192.168.1.200

# Check port
nc -zv 192.168.1.200 3001

# Traceroute
traceroute 192.168.1.200

# Check listening ports
sudo netstat -tulpn | grep 3001
```

---

## Test Flow Hoàn Chỉnh

```bash
# 1. Check Print Bridge
cd local-print-bridge
docker compose ps
docker compose logs --tail=20

# 2. Test health
curl http://localhost:3001/health

# 3. Test from backend machine (if different)
curl http://192.168.1.200:3001/health

# 4. Check backend config
cd ..
grep PRINT_BRIDGE_URL .env

# 5. Test print
curl -X POST http://localhost:3001/print-html \
  -H "Content-Type: application/json" \
  -d '{
    "html": "<html><body><h1>Test</h1></body></html>",
    "printerIP": "192.168.1.115",
    "printerPort": 9100,
    "paperWidth": 80
  }'

# 6. Check backend can reach Print Bridge
# From backend machine:
curl http://192.168.1.200:3001/health
```

---

## Khi Tất Cả Đều Fail

### 1. Rebuild Everything

```bash
# Stop all
cd local-print-bridge
docker compose down
cd ..
docker compose down

# Clean Docker
docker system prune -a

# Rebuild Print Bridge
cd local-print-bridge
docker compose build --no-cache
docker compose up -d
docker compose logs -f

# Restart Backend
cd ..
docker compose up -d backend
```

### 2. Check từ đầu

```bash
# Run diagnostic
./diagnose-print-bridge.sh

# Follow recommendations
```

### 3. Manual Test

```bash
# Test Print Bridge manually
cd local-print-bridge
docker compose up

# Trong terminal khác, test
curl http://localhost:3001/health

# Check logs real-time
```

---

## Liên Hệ Support

Nếu vẫn không giải quyết được:

1. Chạy diagnostic script và lưu output:
   ```bash
   ./diagnose-print-bridge.sh > diagnostic-output.txt
   ```

2. Collect logs:
   ```bash
   cd local-print-bridge
   docker compose logs > print-bridge-logs.txt
   cd ..
   docker compose logs backend > backend-logs.txt
   ```

3. Gửi kèm:
   - diagnostic-output.txt
   - print-bridge-logs.txt
   - backend-logs.txt
   - File .env (remove sensitive data)

---

**Version:** 1.0.0  
**Last Updated:** March 2026
