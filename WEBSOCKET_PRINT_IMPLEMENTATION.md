# WebSocket Print Implementation - Summary

## Tổng quan

Đã implement WebSocket cho Print Bridge để nhận print job real-time từ Backend, thay thế polling cũ.

## Thay đổi

### 1. Print Bridge (local-print-bridge/)

#### Files mới:
- `src/services/websocketClient.js` - WebSocket client kết nối đến Backend
- `src/services/printJobHandler.js` - Xử lý print job nhận từ WebSocket
- `test-websocket.sh` - Script test WebSocket connection
- `README.md` - Documentation
- `DEPLOY_WINDOWS_WEBSOCKET.md` - Hướng dẫn deploy với WebSocket

#### Files cập nhật:
- `package.json` - Thêm `socket.io-client@^4.7.2`
- `src/index.js` - Kết nối WebSocket khi start, graceful shutdown
- `.env` - Thêm comment về WebSocket usage
- `.env.production` - Cập nhật notes về WebSocket

#### Tính năng:
- ✅ Auto-connect WebSocket khi start
- ✅ Auto-reconnect khi mất kết nối (max 10 attempts)
- ✅ Nhận event `print-job-created` real-time
- ✅ Tự động in ngay khi nhận job
- ✅ Cập nhật status về Backend qua HTTP
- ✅ Graceful shutdown

### 2. Backend (backend/)

#### Files cập nhật:
- `infrastructure/websocket/broadcaster.go`:
  - Thêm `PrinterConfigRepository` interface
  - Thêm method `SetPrinterRepository()`
  - Cập nhật `BroadcastPrintJobCreated()` để gửi `printer_ip` và `printer_port`
  
- `main.go`:
  - Wire up `printerConfigRepo` vào `wsBroadcaster`

#### Event format:
```json
{
  "job": {
    "id": "...",
    "type": "BILL",
    "order_id": "...",
    "order_number": "...",
    "printer_id": "...",
    "printer_ip": "192.168.1.115",
    "printer_port": 9100,
    "content": "...",
    "status": "PENDING",
    "created_at": "..."
  }
}
```

### 3. Documentation

#### Files mới:
- `WEBSOCKET_ARCHITECTURE.md` - Kiến trúc chi tiết
- `WEBSOCKET_PRINT_IMPLEMENTATION.md` - Summary này
- `local-print-bridge/README.md` - Print Bridge docs
- `local-print-bridge/DEPLOY_WINDOWS_WEBSOCKET.md` - Deploy guide

## Luồng hoạt động

### Auto-print (WebSocket)

```
1. User tạo order trong Frontend
   ↓
2. Backend tạo print job
   ↓
3. Backend broadcast event "print-job-created" qua WebSocket
   ↓
4. Print Bridge nhận event real-time
   ↓
5. Print Bridge tự động in ngay
   ↓
6. Print Bridge gọi HTTP PUT về Backend để update status
   ↓
7. Backend broadcast "print-job-status-changed" đến Frontend
   ↓
8. Frontend hiển thị notification
```

### Manual print (HTTP)

```
1. User click "Reprint" trong Frontend
   ↓
2. Frontend gọi Backend API
   ↓
3. Backend tạo print job
   ↓
4. Backend broadcast qua WebSocket (giống auto-print)
```

## Ưu điểm so với Polling

| Tính năng | Polling (Cũ) | WebSocket (Mới) |
|-----------|--------------|-----------------|
| Tốc độ | 10-30s delay | < 1s (real-time) |
| Bandwidth | Cao (check liên tục) | Thấp (chỉ khi có job) |
| CPU | Cao | Thấp |
| Độ tin cậy | Có thể bỏ sót | Không bỏ sót |
| Scalability | Kém | Tốt |

## Cấu hình Production

### EC2 Backend

```bash
# .env
PORT=3000
MONGODB_URI=mongodb://localhost:27017/cafe_pos?replicaSet=rs0
JWT_SECRET=your-production-secret

# Security Group
Inbound Rules:
- Port 3000 (HTTP + WebSocket)
```

### Print Bridge (Windows PC)

```bash
# .env.production
PORT=3001
BACKEND_URL=http://YOUR_EC2_IP:3000
DEFAULT_BILL_PRINTER_IP=192.168.1.115
DEFAULT_LABEL_PRINTER_IP=192.168.1.101
```

### Docker Run

```bash
docker run -d \
  --name local-print-bridge \
  --restart unless-stopped \
  --network host \
  --env-file .env.production \
  linhtranphu/local-print-bridge:latest
```

## Testing

### 1. Test WebSocket connection

```bash
# Print Bridge logs
docker logs -f local-print-bridge

# Expect:
[WebSocket] Connecting to backend: http://YOUR_EC2_IP:3000
[WebSocket] ✅ Connected to backend
```

### 2. Test auto-print

```bash
# 1. Tạo order trong Frontend
# 2. Xem Print Bridge logs:
[WebSocket] 📨 New print job received: {job_id}
[PrintJobHandler] Processing job {job_id}
[PrintJobHandler] ✅ Job printed successfully
[PrintJobHandler] ✅ Backend updated - Job {job_id} -> COMPLETED
```

### 3. Test manual print

```bash
curl -X POST http://localhost:3001/print \
  -H "Content-Type: application/json" \
  -d '{
    "jobId": "test",
    "content": "Test print",
    "printerIP": "192.168.1.115",
    "printerPort": 9100,
    "type": "bill"
  }'
```

## Troubleshooting

### WebSocket không kết nối

**Logs:**
```
[WebSocket] Connection error: Error: timeout
```

**Fix:**
1. Kiểm tra BACKEND_URL đúng chưa
2. Kiểm tra EC2 expose port 3000
3. Test: `curl http://YOUR_EC2_IP:3000/api/login`

### WebSocket disconnect liên tục

**Logs:**
```
[WebSocket] Disconnected: transport close
[WebSocket] Reconnecting...
```

**Fix:**
1. Kiểm tra internet ổn định
2. Kiểm tra Backend đang chạy
3. Restart: `docker restart local-print-bridge`

### Không nhận print job

**Logs:**
- WebSocket connected nhưng không có message

**Fix:**
1. Kiểm tra Backend có broadcast event không
2. Xem Backend logs: `docker logs backend`
3. Test tạo order trong Frontend

## Next Steps

### Để deploy production:

1. **Build và push Docker image mới:**
   ```bash
   cd local-print-bridge
   ./build-print-bridge-docker.sh
   ```

2. **Deploy Backend lên EC2:**
   - Upload binary mới
   - Restart backend service
   - Kiểm tra WebSocket endpoint: `http://YOUR_EC2_IP:3000/socket.io/`

3. **Deploy Print Bridge tại quán:**
   - Follow `DEPLOY_WINDOWS_WEBSOCKET.md`
   - Cập nhật `.env.production` với EC2 IP
   - Run Docker container
   - Test WebSocket connection

4. **Monitor:**
   - Backend logs: WebSocket connections
   - Print Bridge logs: Job processing
   - Frontend: Print notifications

## Files Changed

### Backend
- `backend/infrastructure/websocket/broadcaster.go`
- `backend/main.go`

### Print Bridge
- `local-print-bridge/package.json`
- `local-print-bridge/src/index.js`
- `local-print-bridge/src/services/websocketClient.js` (NEW)
- `local-print-bridge/src/services/printJobHandler.js` (NEW)
- `local-print-bridge/.env`
- `local-print-bridge/.env.production`
- `local-print-bridge/test-websocket.sh` (NEW)
- `local-print-bridge/README.md` (NEW)
- `local-print-bridge/DEPLOY_WINDOWS_WEBSOCKET.md` (NEW)

### Documentation
- `WEBSOCKET_ARCHITECTURE.md` (NEW)
- `WEBSOCKET_PRINT_IMPLEMENTATION.md` (NEW)

## Kết luận

✅ WebSocket implementation hoàn tất
✅ Print Bridge nhận job real-time
✅ Backend broadcast với printer info
✅ Documentation đầy đủ
✅ Ready for production deployment

**Port 3000 VẪN PHẢI EXPOSE** cho cả WebSocket và HTTP API.
