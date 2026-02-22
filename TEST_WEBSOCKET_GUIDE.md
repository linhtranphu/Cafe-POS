# Hướng Dẫn Test WebSocket Máy In

## Tình Trạng Hiện Tại

✅ Backend đã cài đặt Socket.IO server (go-socket.io v1.7.0)
✅ Print Bridge đã cài đặt WebSocket client (socket.io-client v2.5.0)
✅ Import statement đã được sửa từ `const { io }` thành `const io`
⏳ Cần restart Print Bridge để test

## Bước 1: Restart Print Bridge

### Nếu đang chạy local (npm start):
```bash
# Dừng process hiện tại (Ctrl+C)
# Sau đó start lại:
cd local-print-bridge
npm start
```

### Nếu đang chạy Docker:
```bash
# Restart container
docker restart local-print-bridge

# Hoặc stop và start lại
docker stop local-print-bridge
docker start local-print-bridge
```

## Bước 2: Kiểm Tra Logs

### Print Bridge Logs
Sau khi restart, bạn sẽ thấy:
```
[WebSocket] Connecting to: http://localhost:3000
[WebSocket] ✅ Connected to backend
```

Nếu thấy lỗi:
```
[WebSocket] Connection error (attempt 1/10): ...
```
→ Có vấn đề, xem phần Troubleshooting bên dưới

### Backend Logs
```bash
tail -f backend.log | grep Socket
```

Bạn sẽ thấy:
```
[Socket.IO] Client connected: <socket-id>
```

## Bước 3: Test Kết Nối Đơn Giản

Chạy test script:
```bash
node test-websocket-connection.js
```

Kết quả mong đợi:
```
Testing WebSocket connection to backend...
✅ Connected to backend!
Socket ID: <some-id>
Waiting for events...
Test complete. Disconnecting...
```

## Bước 4: Test Print Job Thực Tế

### Cách 1: Qua Frontend
1. Mở frontend: http://localhost:5173
2. Login với admin/admin123
3. Tạo một order mới
4. Xem logs của Print Bridge:
   ```
   [WebSocket] 📨 New print job received: { job: { id: '...', content: '...', printer_ip: '...', ... } }
   [PrintJobHandler] Processing job <id> - Type: ORDER, Printer: <ip>:<port>
   [PrintJobHandler] ✅ Job <id> printed successfully
   [PrintJobHandler] ✅ Backend updated - Job <id> -> COMPLETED
   ```

### Cách 2: Qua API (Manual Test)
```bash
# Tạo print job qua API
curl -X POST http://localhost:3000/api/manager/print-jobs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "type": "ORDER",
    "content": "Test Print\nWebSocket Test\n",
    "printer_ip": "192.168.1.100",
    "printer_port": 9100
  }'
```

## Kết Quả Mong Đợi

### ✅ Thành Công
1. Print Bridge kết nối WebSocket thành công
2. Nhận được event `print-job-created` khi có job mới
3. Gửi lệnh in đến máy in
4. Cập nhật status về backend (COMPLETED hoặc FAILED)
5. Máy in thực sự in ra giấy

### ❌ Lỗi Có Thể Gặp

#### Lỗi 1: Connection Timeout
```
[WebSocket] Connection error (attempt 1/10): timeout
```

**Nguyên nhân:**
- Backend chưa chạy
- Port 3000 bị block
- Firewall chặn kết nối

**Giải pháp:**
```bash
# Kiểm tra backend đang chạy
curl http://localhost:3000/health

# Kiểm tra Socket.IO endpoint
curl "http://localhost:3000/socket.io/?EIO=3&transport=polling"
```

#### Lỗi 2: Module Not Found
```
Error: Cannot find module 'socket.io-client'
```

**Giải pháp:**
```bash
cd local-print-bridge
npm install
```

#### Lỗi 3: Printer Connection Failed
```
[PrintJobHandler] ❌ Job <id> failed: connect ETIMEDOUT
```

**Nguyên nhân:**
- Máy in không online
- IP address sai
- Port sai (thường là 9100)

**Giải pháp:**
```bash
# Test kết nối máy in
nc -zv 192.168.1.100 9100

# Hoặc
telnet 192.168.1.100 9100
```

## Bước 5: Kiểm Tra Toàn Bộ Flow

### Checklist
- [ ] Backend chạy trên port 3000
- [ ] MongoDB chạy với replica set
- [ ] Print Bridge kết nối WebSocket thành công
- [ ] Tạo order từ frontend
- [ ] Print Bridge nhận được event
- [ ] Print job được gửi đến máy in
- [ ] Máy in thực sự in ra
- [ ] Backend status được cập nhật

## Debug Commands

### Xem tất cả logs
```bash
# Backend
tail -f backend.log

# Print Bridge (nếu chạy Docker)
docker logs -f local-print-bridge

# MongoDB
docker logs cafe-pos-mongodb
```

### Kiểm tra kết nối
```bash
# Backend health
curl http://localhost:3000/health

# Print Bridge health
curl http://localhost:3001/health

# Socket.IO endpoint
curl "http://localhost:3000/socket.io/?EIO=3&transport=polling"
```

### Kiểm tra processes
```bash
# Backend
lsof -i :3000

# Print Bridge
lsof -i :3001

# MongoDB
docker ps | grep mongodb
```

## Lưu Ý Quan Trọng

1. **Socket.IO Version Compatibility**
   - Backend: go-socket.io v1.7.0 (Engine.IO v3)
   - Client: socket.io-client v2.5.0 (Engine.IO v3)
   - ✅ Compatible!

2. **Import Syntax**
   - ❌ `const { io } = require('socket.io-client')` (ES6 destructuring - không work với v2.x)
   - ✅ `const io = require('socket.io-client')` (CommonJS default export)

3. **Network Configuration**
   - Print Bridge kết nối RA ngoài (outbound) đến backend
   - Không cần mở port 3001 trên router
   - Chỉ cần backend port 3000 accessible

4. **Printer Configuration**
   - Máy in phải cùng mạng LAN với Print Bridge
   - Port mặc định: 9100 (ESC/POS)
   - Test kết nối trước: `nc -zv <printer-ip> 9100`

## Next Steps Sau Khi Test Thành Công

1. Build Docker image mới với WebSocket:
   ```bash
   cd local-print-bridge
   ./build-print-bridge-docker.sh
   ```

2. Deploy lên production (Windows machine):
   - Xem `local-print-bridge/DEPLOY_WINDOWS_WEBSOCKET.md`

3. Cập nhật documentation:
   - Thêm WebSocket setup vào README
   - Cập nhật deployment guide

## Liên Hệ

Nếu gặp vấn đề, cung cấp:
1. Print Bridge logs
2. Backend logs
3. Lỗi cụ thể
4. Môi trường (local/Docker/production)
