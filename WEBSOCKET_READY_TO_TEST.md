# WebSocket Sẵn Sàng Test! 🚀

## Tình Trạng Hiện Tại

### ✅ Đã Hoàn Thành
1. **Backend Socket.IO Server**
   - Library: `github.com/googollee/go-socket.io v1.7.0`
   - Protocol: Engine.IO v3
   - Endpoint: `http://localhost:3000/socket.io/`
   - Status: ✅ Running và responding

2. **Print Bridge WebSocket Client**
   - Library: `socket.io-client v2.5.0`
   - Protocol: Engine.IO v3 (compatible!)
   - Import: Fixed từ `const { io }` → `const io`
   - Status: ⏳ Chưa start

3. **Version Compatibility**
   - Backend: go-socket.io v1.7.0 (Engine.IO v3) ✅
   - Client: socket.io-client v2.5.0 (Engine.IO v3) ✅
   - **100% Compatible!**

### 🔧 Đã Sửa
- Import statement trong `websocketClient.js`
- Downgrade socket.io-client từ v4.x → v2.5.0
- Chuyển từ custom Socket.IO implementation sang standard library

## Cách Test

### Option 1: Start Print Bridge Manually (Recommended)

```bash
# Terminal 1: Watch backend logs
tail -f backend.log | grep -i socket

# Terminal 2: Start Print Bridge
cd local-print-bridge
npm start
```

**Kết quả mong đợi trong Terminal 2:**
```
🖨️  Local Print Bridge Server
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ Server running on port 3001
📡 Backend URL: http://localhost:3000
🔄 Polling interval: 30000ms
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

[WebSocket] Connecting to: http://localhost:3000
[WebSocket] ✅ Connected to backend
```

**Kết quả mong đợi trong Terminal 1 (backend logs):**
```
[Socket.IO] Client connected: <socket-id>
```

### Option 2: Use Test Script

```bash
./test-print-bridge-websocket.sh
```

## Test Print Job Flow

### Bước 1: Verify Connection
Sau khi Print Bridge start, kiểm tra logs có dòng:
```
[WebSocket] ✅ Connected to backend
```

### Bước 2: Create Test Print Job

#### Via Frontend:
1. Mở http://localhost:5173
2. Login: admin/admin123
3. Tạo order mới
4. Xem Print Bridge logs

#### Via API:
```bash
# Get auth token first
TOKEN=$(curl -s -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.token')

# Create print job
curl -X POST http://localhost:3000/api/manager/print-jobs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "type": "ORDER",
    "content": "=== TEST PRINT ===\nWebSocket Test\nTimestamp: '$(date)'\n",
    "printer_ip": "192.168.1.100",
    "printer_port": 9100
  }'
```

### Bước 3: Verify Print Bridge Receives Event

**Print Bridge logs sẽ hiển thị:**
```
[WebSocket] 📨 New print job received: {
  job: {
    id: '...',
    content: '=== TEST PRINT ===...',
    printer_ip: '192.168.1.100',
    printer_port: 9100,
    type: 'ORDER'
  }
}
[PrintJobHandler] Processing job <id> - Type: ORDER, Printer: 192.168.1.100:9100
```

### Bước 4: Verify Printer Prints

Nếu máy in online và configured đúng:
```
[PrintJobHandler] ✅ Job <id> printed successfully
[PrintJobHandler] ✅ Backend updated - Job <id> -> COMPLETED
```

Nếu máy in offline hoặc có lỗi:
```
[PrintJobHandler] ❌ Job <id> failed: connect ETIMEDOUT
[PrintJobHandler] Backend updated - Job <id> -> FAILED
```

## Troubleshooting

### Lỗi: Connection Timeout
```
[WebSocket] Connection error (attempt 1/10): timeout
```

**Kiểm tra:**
```bash
# Backend có chạy không?
curl http://localhost:3000/health

# Socket.IO endpoint có hoạt động không?
curl "http://localhost:3000/socket.io/?EIO=3&transport=polling"

# Backend logs có lỗi không?
tail -f backend.log | grep -i error
```

### Lỗi: Module Not Found
```
Error: Cannot find module 'socket.io-client'
```

**Giải pháp:**
```bash
cd local-print-bridge
npm install
```

### Lỗi: Version Mismatch
```
It seems you are trying to reach a Socket.IO server in v2.x with a v3.x client
```

**Kiểm tra version:**
```bash
# Client version (phải là 2.5.0)
cat local-print-bridge/package.json | grep socket.io-client

# Backend version (phải là 1.7.0)
grep go-socket.io backend/go.mod
```

### Lỗi: Printer Connection Failed
```
[PrintJobHandler] ❌ Job <id> failed: connect ETIMEDOUT
```

**Kiểm tra máy in:**
```bash
# Test kết nối
nc -zv 192.168.1.100 9100

# Hoặc
telnet 192.168.1.100 9100

# Ping máy in
ping 192.168.1.100
```

## Success Criteria

### ✅ WebSocket Connection Success
- [ ] Print Bridge logs: `[WebSocket] ✅ Connected to backend`
- [ ] Backend logs: `[Socket.IO] Client connected: <id>`
- [ ] Không có connection errors
- [ ] Connection stable (không disconnect/reconnect liên tục)

### ✅ Print Job Flow Success
- [ ] Backend broadcasts `print-job-created` event
- [ ] Print Bridge receives event
- [ ] Print job handler processes job
- [ ] Printer receives print command
- [ ] Paper prints out (if printer online)
- [ ] Backend status updated to COMPLETED

## Next Steps After Success

### 1. Build Docker Image
```bash
cd local-print-bridge
./build-print-bridge-docker.sh
```

### 2. Test Docker Container
```bash
# Stop local process
# Ctrl+C in Print Bridge terminal

# Start Docker container
docker run -d \
  --name local-print-bridge \
  --restart unless-stopped \
  --network host \
  --env-file local-print-bridge/.env \
  linhtranphu/local-print-bridge:latest

# Watch logs
docker logs -f local-print-bridge
```

### 3. Deploy to Production
- Xem `local-print-bridge/DEPLOY_WINDOWS_WEBSOCKET.md`
- Update Windows machine với Docker image mới
- Configure `.env` với production backend URL

### 4. Update Documentation
- [ ] Update README với WebSocket setup
- [ ] Update deployment guide
- [ ] Add troubleshooting section

## Files Changed

### Backend
- `backend/infrastructure/websocket/socketio_server.go` (NEW)
- `backend/infrastructure/websocket/broadcaster.go` (UPDATED)
- `backend/main.go` (UPDATED)
- `backend/go.mod` (UPDATED - added go-socket.io)

### Print Bridge
- `local-print-bridge/src/services/websocketClient.js` (NEW)
- `local-print-bridge/src/services/printJobHandler.js` (NEW)
- `local-print-bridge/src/index.js` (UPDATED)
- `local-print-bridge/package.json` (UPDATED - added socket.io-client v2.5.0)

### Documentation
- `TEST_WEBSOCKET_GUIDE.md` (NEW)
- `WEBSOCKET_READY_TO_TEST.md` (NEW)
- `test-print-bridge-websocket.sh` (NEW)
- `test-websocket-simple.js` (NEW)

## Quick Commands Reference

```bash
# Start everything
./restart_local.sh

# Start Print Bridge only
cd local-print-bridge && npm start

# Watch backend logs
tail -f backend.log | grep -i socket

# Watch Print Bridge logs (Docker)
docker logs -f local-print-bridge

# Test WebSocket connection
./test-print-bridge-websocket.sh

# Test printer connection
nc -zv <printer-ip> 9100

# Create test print job (need auth token)
curl -X POST http://localhost:3000/api/manager/print-jobs \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"type":"ORDER","content":"Test","printer_ip":"192.168.1.100"}'
```

## Support

Nếu gặp vấn đề, cung cấp:
1. Print Bridge logs (toàn bộ output)
2. Backend logs (grep Socket)
3. Lỗi cụ thể
4. Environment (local/Docker/production)
5. Printer IP và port

---

**Status**: ✅ Ready to test!
**Next Action**: Start Print Bridge và verify connection
