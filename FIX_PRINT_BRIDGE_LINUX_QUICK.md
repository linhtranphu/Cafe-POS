# Fix Print Bridge Trên Linux - Quick Guide

## Vấn Đề

```
✗ chromium-browser not found
✗ Cannot POST /print-html (404)
```

## Nguyên Nhân

1. Container đang chạy là bản cũ KHÔNG có Chromium
2. Endpoint đúng là `/render-and-print` chứ không phải `/print-html`

## Giải Pháp

### Trên Máy Linux

```bash
# 1. Copy script sang máy Linux
scp rebuild-print-bridge-linux.sh user@linux-machine:/path/to/project/

# 2. SSH vào máy Linux
ssh user@linux-machine

# 3. Chạy script rebuild
cd /path/to/project
chmod +x rebuild-print-bridge-linux.sh
./rebuild-print-bridge-linux.sh

# Script sẽ:
# - Backup files hiện tại
# - Stop container cũ
# - Update Dockerfile với Chromium
# - Update docker-compose.yml với cấu hình Linux
# - Build image mới (mất 3-5 phút)
# - Start container mới
# - Verify Chromium
# - Test endpoints
```

## Sau Khi Rebuild

### 1. Kiểm Tra Container

```bash
# Check container status
docker ps | grep print-bridge

# Should show:
# local-print-bridge   Up X minutes (healthy)
```

### 2. Kiểm Tra Chromium

```bash
# Verify Chromium exists
docker exec local-print-bridge which chromium-browser

# Should show:
# /usr/bin/chromium-browser

# Check version
docker exec local-print-bridge chromium-browser --version

# Should show:
# Chromium 119.x.x.x
```

### 3. Test Endpoints

```bash
# Health check
curl http://localhost:3001/health

# Test connection
curl -X POST http://localhost:3001/test-connection \
  -H "Content-Type: application/json" \
  -d '{"printerIP": "192.168.1.100", "printerPort": 9100}'

# Test render and print (endpoint đúng)
curl -X POST http://localhost:3001/render-and-print \
  -H "Content-Type: application/json" \
  -d '{
    "html": "<html><body><h1>Test Print</h1></body></html>",
    "printerIP": "192.168.1.100",
    "printerPort": 9100,
    "paperWidth": 80
  }'
```

### 4. Test Từ Web Interface

1. Truy cập https://tacafe.store/#/print-management
2. Kiểm tra kết nối (should be OK)
3. Bấm "In thử" (should work now!)

## Endpoints Đúng

Print Bridge Go có các endpoints sau:

| Endpoint | Method | Mô tả |
|----------|--------|-------|
| `/health` | GET | Health check |
| `/status` | GET | Service status |
| `/test-connection` | POST | Test printer connection |
| `/render-and-print` | POST | Render HTML và in (ĐÚNG) |
| `/print` | POST | Direct print ESC/POS |
| `/test-render` | POST | Test render (debug) |

**Lưu ý:** Không có endpoint `/print-html`!

## Nếu Vẫn Lỗi

### Check Logs

```bash
# Real-time logs
docker logs -f local-print-bridge

# Last 50 lines
docker logs --tail=50 local-print-bridge

# Filter errors
docker logs local-print-bridge 2>&1 | grep -i error
```

### Common Issues

**1. Build failed**
```bash
# Check Docker space
docker system df

# Clean up
docker system prune -a

# Rebuild
cd local-print-bridge
docker compose build --no-cache
```

**2. Container keeps restarting**
```bash
# Check logs
docker logs local-print-bridge

# Check if port 3001 is in use
sudo lsof -i :3001

# Try different port
# Edit .env: PORT=3002
# Edit docker-compose.yml: ports: "3002:3002"
```

**3. Chromium errors**
```bash
# Increase shared memory
# Edit docker-compose.yml:
shm_size: '1gb'

# Restart
docker compose down
docker compose up -d
```

**4. Cannot connect to printer**
```bash
# From container, test network
docker exec local-print-bridge ping 192.168.1.100

# Test port
docker exec local-print-bridge nc -zv 192.168.1.100 9100
```

## Diagnostic Commands

```bash
# Full diagnostic
./diagnose-print-bridge-linux.sh

# Quick checks
docker ps | grep print-bridge
docker exec local-print-bridge which chromium-browser
curl http://localhost:3001/health
docker logs --tail=20 local-print-bridge
```

## Rollback

Nếu cần rollback về version cũ:

```bash
cd local-print-bridge

# Stop current
docker compose down

# Restore backup
cp docker-compose.yml.backup.YYYYMMDD_HHMMSS docker-compose.yml
cp Dockerfile.backup.YYYYMMDD_HHMMSS Dockerfile

# Rebuild
docker compose build
docker compose up -d
```

## Checklist

- [ ] Script rebuild đã chạy thành công
- [ ] Container đang chạy và healthy
- [ ] Chromium đã được cài đặt trong container
- [ ] Health endpoint trả về OK
- [ ] Test connection thành công
- [ ] Logs không có errors
- [ ] Test in thử từ web interface thành công

## Thời Gian Ước Tính

- Download script: 1 phút
- Chạy rebuild script: 5-10 phút (tùy tốc độ mạng)
- Test và verify: 2-3 phút
- **Tổng: ~15 phút**

## Support

Nếu vẫn gặp vấn đề sau khi rebuild:

1. Chạy diagnostic và lưu output:
   ```bash
   ./diagnose-print-bridge-linux.sh > diagnostic.txt
   ```

2. Collect logs:
   ```bash
   docker logs local-print-bridge > print-bridge.log
   ```

3. Check container config:
   ```bash
   docker inspect local-print-bridge > container-info.json
   ```

4. Gửi kèm:
   - diagnostic.txt
   - print-bridge.log
   - container-info.json

---

**TL;DR:** Chạy `./rebuild-print-bridge-linux.sh` trên máy Linux, đợi 5-10 phút, test lại!
