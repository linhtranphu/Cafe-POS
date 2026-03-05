# Print Bridge WebSocket Setup - Kết nối trực tiếp với Backend

## Kiến trúc

```
┌─────────────┐         ┌─────────────┐         ┌──────────────┐
│   Backend   │◄────────│Print Bridge │◄────────│   Printer    │
│   (EC2)     │WebSocket│  (Local)    │  ESC/POS│  (Thermal)   │
└─────────────┘         └─────────────┘         └──────────────┘
      │
      │ HTTP API
      ▼
┌─────────────┐
│  Frontend   │
│  (Browser)  │
└─────────────┘
```

## Flow hoạt động

1. User tạo order trên Frontend → Backend
2. Backend tạo print job → Emit `print-job-created` qua WebSocket
3. Print Bridge nhận event → In ra máy in
4. Print Bridge update status về Backend qua HTTP API

## Ưu điểm so với Frontend → Print Bridge

✅ Không bị CORS (backend → print bridge đều là server-side)
✅ Không cần HTTPS cho EC2
✅ Print Bridge chỉ cần kết nối 1 lần với backend
✅ Realtime hơn, không phụ thuộc browser

## Cấu hình

### 1. Backend (Đã sẵn sàng)

Backend đã có:
- Socket.IO server listening trên port 3000
- Emit event `print-job-created` khi tạo print job
- Endpoint `/api/print-jobs/:id/status` để update status

File: `backend/infrastructure/websocket/broadcaster.go`
```go
func (b *Broadcaster) BroadcastPrintJobCreated(job *printing.PrintJob) {
    jobData := map[string]interface{}{
        "job": map[string]interface{}{
            "id":           job.ID.Hex(),
            "content":      job.Content,
            "printer_ip":   printerIP,
            "printer_port": printerPort,
            "type":         string(job.Type),
        },
    }
    b.socketIOServer.BroadcastToNamespace("/", "print-job-created", jobData)
}
```

### 2. Print Bridge (Đã sẵn sàng)

Print Bridge đã có:
- WebSocket client kết nối với backend
- Listen event `print-job-created`
- Auto print và update status về backend

File: `local-print-bridge/src/services/websocketClient.js`

### 3. Cấu hình môi trường

#### Option A: Print Bridge trong mạng LAN (Khuyến nghị)

Máy tính chạy Print Bridge trong cùng mạng với máy in:

```bash
# local-print-bridge/.env
PORT=3001

# Backend URL - Dùng domain hoặc IP public của EC2
BACKEND_URL=https://tacafe.store
# Hoặc: BACKEND_URL=http://52.77.228.154:3000

# Printer IPs (trong mạng LAN)
DEFAULT_BILL_PRINTER_IP=192.168.1.100
DEFAULT_BILL_PRINTER_PORT=9100

DEFAULT_LABEL_PRINTER_IP=192.168.1.101
DEFAULT_LABEL_PRINTER_PORT=9100

LOG_LEVEL=info
PRINTER_TIMEOUT=5000
```

#### Option B: Print Bridge trên cùng máy EC2 (Testing)

Nếu muốn test trên EC2 (không có máy in thật):

```bash
# local-print-bridge/.env
PORT=3001
BACKEND_URL=http://localhost:3000
DEFAULT_BILL_PRINTER_IP=192.168.1.100
DEFAULT_BILL_PRINTER_PORT=9100
```

### 4. Khởi động Print Bridge

```bash
cd local-print-bridge

# Install dependencies
npm install

# Start service
npm start
```

Log sẽ hiển thị:
```
🖨️  Local Print Bridge Started
Server running on: http://localhost:3001
Backend URL: https://tacafe.store
[WebSocket] Connecting to backend: https://tacafe.store
[WebSocket] ✅ Connected to backend
Ready to accept print requests!
```

## Testing

### Test 1: Kiểm tra kết nối

```bash
# Health check
curl http://localhost:3001/health

# Expected:
{
  "status": "ok",
  "service": "Local Print Bridge",
  "version": "1.0.0"
}
```

### Test 2: Kiểm tra WebSocket connection

Xem log của Print Bridge:
```
[WebSocket] ✅ Connected to backend
```

Xem log của Backend:
```
[Socket.IO] Client connected: <socket-id>
```

### Test 3: Test in thử

Tạo order trên Frontend → Kiểm tra log Print Bridge:
```
[WebSocket] 📨 New print job received: { job: { id: '...', content: '...' } }
[PrintJobHandler] Processing job <id> - Type: BILL, Printer: 192.168.1.100:9100
[PrintJobHandler] ✅ Job <id> printed successfully
[PrintJobHandler] ✅ Backend updated - Job <id> -> COMPLETED
```

## Troubleshooting

### Lỗi: WebSocket connection failed

**Nguyên nhân:** Backend chưa có Socket.IO proxy trong nginx

**Giải pháp:** Đã fix trong `WEBSOCKET_EC2_FIX.md`
```bash
./fix-websocket-ec2.sh
```

### Lỗi: ECONNREFUSED connecting to backend

**Nguyên nhân:** BACKEND_URL sai hoặc backend chưa chạy

**Giải pháp:**
```bash
# Kiểm tra backend có chạy không
curl https://tacafe.store/api/state-machines

# Kiểm tra BACKEND_URL trong .env
cat local-print-bridge/.env | grep BACKEND_URL
```

### Lỗi: Print job received but printer not responding

**Nguyên nhân:** Printer IP sai hoặc printer offline

**Giải pháp:**
```bash
# Test printer connection
curl -X POST http://localhost:3001/test-connection \
  -H "Content-Type: application/json" \
  -d '{"printerIP": "192.168.1.100", "printerPort": 9100}'

# Ping printer
ping 192.168.1.100
```

### Lỗi: Backend update failed

**Nguyên nhân:** Backend API không accessible từ Print Bridge

**Giải pháp:**
```bash
# Test backend API
curl https://tacafe.store/api/print-jobs/<job-id>/status \
  -X PUT \
  -H "Content-Type: application/json" \
  -d '{"status": "COMPLETED"}'
```

## Security Notes

### Production Setup

1. **Firewall rules:**
   - EC2 Security Group: Mở port 3000 cho IP của Print Bridge
   - Print Bridge machine: Chỉ cho phép kết nối đến printer IPs

2. **Authentication (Optional):**
   Thêm token authentication cho WebSocket:
   ```javascript
   // Print Bridge
   this.socket = io(backendUrl, {
     auth: {
       token: process.env.BRIDGE_AUTH_TOKEN
     }
   })
   ```

3. **HTTPS:**
   - Backend nên dùng HTTPS (wss://)
   - Update BACKEND_URL: `wss://tacafe.store`

## Monitoring

### Check Print Bridge status

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
  "uptime": 86400,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Log files

Print Bridge logs to console. Để lưu log:
```bash
npm start > print-bridge.log 2>&1
```

Hoặc dùng PM2:
```bash
npm install -g pm2
pm2 start src/index.js --name print-bridge
pm2 logs print-bridge
```

## Next Steps

1. ✅ Deploy backend với WebSocket fix
2. ✅ Cấu hình Print Bridge với BACKEND_URL
3. ✅ Test kết nối WebSocket
4. ✅ Test in thử với order thật
5. 🔄 Setup PM2 hoặc systemd để auto-start Print Bridge
6. 🔄 Monitor logs và performance

## Related Files

- Backend WebSocket: `backend/infrastructure/websocket/`
- Print Bridge WebSocket Client: `local-print-bridge/src/services/websocketClient.js`
- Print Job Handler: `local-print-bridge/src/services/printJobHandler.js`
- Backend Sync: `local-print-bridge/src/services/backendSync.js`
