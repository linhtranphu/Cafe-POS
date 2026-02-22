# Local Print Bridge

Bridge service để kết nối Cafe POS Backend với máy in nhiệt qua mạng LAN.

## Tính năng

- ✅ **WebSocket Client**: Nhận print job real-time từ Backend
- ✅ **HTTP Client**: Cập nhật trạng thái job về Backend
- ✅ **HTTP Server**: API để manual print/reprint
- ✅ **Auto-reconnect**: Tự động kết nối lại khi mất kết nối
- ✅ **Error handling**: Xử lý lỗi và retry
- ✅ **Logging**: Chi tiết logs để debug

## Kiến trúc

```
Backend (EC2)
    │
    │ WebSocket (real-time push)
    ▼
Print Bridge (Windows PC)
    │
    │ Raw TCP (Port 9100)
    ▼
Thermal Printers (LAN)
```

## Cài đặt

### Development

```bash
npm install
cp .env.example .env
# Cập nhật BACKEND_URL trong .env
npm start
```

### Production (Docker)

```bash
# Build image
docker build -t local-print-bridge .

# Hoặc pull từ Docker Hub
docker pull linhtranphu/local-print-bridge:latest

# Run
docker run -d \
  --name local-print-bridge \
  --restart unless-stopped \
  --network host \
  --env-file .env \
  linhtranphu/local-print-bridge:latest
```

## Cấu hình

File `.env`:

```bash
# Server
PORT=3001

# Backend (WebSocket + HTTP)
BACKEND_URL=http://YOUR_EC2_IP:3000

# Printers
DEFAULT_BILL_PRINTER_IP=192.168.1.115
DEFAULT_LABEL_PRINTER_IP=192.168.1.101
```

## API Endpoints

### Health Check
```
GET /health
```

### Manual Print
```
POST /print
Body: {
  "jobId": "string",
  "content": "string",
  "printerIP": "string",
  "printerPort": 9100,
  "type": "bill|label"
}
```

### Test Connection
```
POST /test-connection
Body: {
  "printerIP": "string",
  "printerPort": 9100
}
```

### Status
```
GET /status
```

## WebSocket Events

Print Bridge lắng nghe các events từ Backend:

- `print-job-created`: Job mới được tạo → Tự động in
- `print-job-status-changed`: Trạng thái job thay đổi
- `print-job-failed`: Job thất bại

## Logs

```bash
# Docker
docker logs -f local-print-bridge

# Development
npm start
```

Logs sẽ hiển thị:
- WebSocket connection status
- Print job processing
- Printer communication
- Errors and warnings

## Troubleshooting

### WebSocket không kết nối

```
[WebSocket] Connection error: Error: timeout
```

**Giải pháp:**
- Kiểm tra BACKEND_URL đúng chưa
- Kiểm tra EC2 expose port 3000
- Test: `curl http://YOUR_EC2_IP:3000/api/login`

### Máy in không in

```
[PrintJobHandler] ❌ Job failed: connect ETIMEDOUT
```

**Giải pháp:**
- Kiểm tra printer IP đúng chưa
- Kiểm tra printer bật và kết nối LAN
- Test: `POST /test-connection`

### Backend không nhận status update

```
[PrintJobHandler] Failed to update backend
```

**Giải pháp:**
- Kiểm tra BACKEND_URL có http:// prefix
- Kiểm tra EC2 backend đang chạy
- Kiểm tra firewall/security group

## Development

### Structure

```
src/
├── index.js                    # Main server
├── services/
│   ├── websocketClient.js      # WebSocket client
│   ├── printJobHandler.js      # Job processing
│   ├── printerService.js       # Printer communication
│   └── backendSync.js          # Backend API calls
└── utils/
    └── logger.js               # Logging utility
```

### Testing

```bash
# Test WebSocket
./test-websocket.sh

# Test manual print
./test-print.sh
```

## Deployment

Xem [DEPLOY_WINDOWS_WEBSOCKET.md](./DEPLOY_WINDOWS_WEBSOCKET.md) để hướng dẫn chi tiết deploy trên Windows PC tại quán cafe.

### Build Docker Image

```bash
./build-print-bridge-docker.sh 1.1.0
```

## License

MIT
