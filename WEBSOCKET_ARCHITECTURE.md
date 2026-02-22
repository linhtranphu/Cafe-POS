# WebSocket Architecture - Print System

## Câu trả lời: Có cần expose port 3000 không?

**CÓ - Port 3000 VẪN CẦN EXPOSE** vì:

1. **WebSocket dùng cho cả Frontend và Print Bridge** (real-time notifications)
2. **Print Bridge vẫn cần HTTP REST API** để cập nhật trạng thái về Backend

## Kiến trúc mới (với WebSocket)

```
┌─────────────────────────────────────────────────────────────┐
│                         EC2 Server                          │
│  ┌──────────────┐         ┌──────────────┐                 │
│  │   Frontend   │◄────────┤   Backend    │                 │
│  │  (Vue.js)    │ WebSocket│   (Go)      │                 │
│  │              │         │   Port 3000  │                 │
│  └──────────────┘         └──────┬───────┘                 │
│                                   │                          │
│                                   │ WebSocket + HTTP         │
│                                   │ (exposed)                │
└───────────────────────────────────┼──────────────────────────┘
                                    │
                                    │ Internet
                                    │
                    ┌───────────────▼────────────────┐
                    │   Windows PC (Quán Cafe)       │
                    │  ┌──────────────────────────┐  │
                    │  │  Local Print Bridge      │  │
                    │  │  Port 3001               │  │
                    │  │                          │  │
                    │  │  ✅ WebSocket Client     │  │
                    │  │     - Nhận job real-time │  │
                    │  │     - Tự động in ngay    │  │
                    │  │                          │  │
                    │  │  ✅ HTTP Client          │  │
                    │  │     - Cập nhật status    │  │
                    │  │                          │  │
                    │  │  ✅ HTTP Server          │  │
                    │  │     - Manual print API   │  │
                    │  └──────────┬───────────────┘  │
                    │             │                   │
                    │             │ Raw TCP          │
                    │             │ Port 9100        │
                    │  ┌──────────▼───────────────┐  │
                    │  │  Thermal Printers        │  │
                    │  │  - Bill: 192.168.1.115   │  │
                    │  │  - Label: 192.168.1.101  │  │
                    │  └──────────────────────────┘  │
                    └──────────────────────────────┘
```

## Chi tiết luồng hoạt động

### 1. WebSocket (Frontend ↔ Backend)

**Mục đích**: Thông báo real-time cho người dùng về trạng thái in

**Backend Implementation**:
- File: `backend/infrastructure/websocket/hub.go`
- File: `backend/infrastructure/websocket/socketio_handler.go`
- File: `backend/infrastructure/websocket/broadcaster.go`
- Endpoint: `GET /socket.io/` (Socket.IO compatible)
- Đã được khởi động trong `main.go`:
  ```go
  wsHub := websocket.NewHub()
  go wsHub.Run()
  wsBroadcaster := websocket.NewBroadcaster(wsHub)
  wsBroadcaster.SetPrinterRepository(printerConfigRepo)
  ```

**Frontend Implementation**:
- File: `frontend/src/services/websocket.js`
- Sử dụng Socket.IO client
- Kết nối đến: `http://localhost:3000` (hoặc EC2 domain)
- Lắng nghe events: `print-job-created`, `print-job-status-changed`, etc.

**Luồng hoạt động**:
```
1. User tạo order → Backend tạo print job
2. Backend broadcast event qua WebSocket → Frontend nhận thông báo
3. Frontend hiển thị notification real-time cho user
```

### 2. WebSocket (Print Bridge ↔ Backend) - ✨ MỚI

**Mục đích**: Nhận lệnh in real-time, không cần polling

**Print Bridge Implementation**:
- File: `local-print-bridge/src/services/websocketClient.js`
- File: `local-print-bridge/src/services/printJobHandler.js`
- Sử dụng Socket.IO client
- Kết nối đến: `http://[EC2_DOMAIN]:3000`
- Lắng nghe event: `print-job-created`

**Luồng hoạt động**:
```
1. Backend tạo print job
2. Backend broadcast event "print-job-created" với data:
   {
     job: {
       id: "...",
       content: "...",
       printer_ip: "192.168.1.115",
       printer_port: 9100,
       type: "BILL"
     }
   }
3. Print Bridge nhận event qua WebSocket
4. Print Bridge tự động in ngay lập tức
5. Print Bridge gọi HTTP PUT về Backend để cập nhật status
```

**Ưu điểm so với polling**:
- ⚡ In nhanh hơn (real-time, không delay)
- 💰 Tiết kiệm bandwidth (không cần poll liên tục)
- 🔋 Tiết kiệm CPU (không cần check liên tục)
- 🎯 Chính xác hơn (không bỏ sót job)

### 3. HTTP REST API (Print Bridge ↔ Backend)

**Mục đích**: Cập nhật trạng thái và manual print

**Print Bridge → Backend** (Cập nhật trạng thái):
```
PUT http://[EC2_DOMAIN]:3000/api/manager/print-jobs/:id/status
Body: {
  "status": "COMPLETED",  // hoặc "FAILED"
  "error_msg": ""
}
```

**Backend → Print Bridge** (Manual print - optional):
```
POST http://[PRINT_BRIDGE_IP]:3001/print
Body: {
  "jobId": "...",
  "content": "...",
  "printerIP": "192.168.1.115",
  "printerPort": 9100,
  "type": "bill"
}
```

**⚠️ QUAN TRỌNG**: 
- WebSocket dùng cho auto-print (real-time)
- HTTP POST dùng cho manual print/reprint
- HTTP PUT dùng cho status update

## Cấu hình Production

### EC2 Server (.env)
```bash
PORT=3000
MONGODB_URI=mongodb://localhost:27017/cafe_pos?replicaSet=rs0
JWT_SECRET=your-production-secret-key

# WebSocket sẽ tự động chạy trên cùng port 3000
```

### Print Bridge (.env.production)
```bash
# Server Configuration
PORT=3001

# Backend URL - PHẢI là EC2 public domain/IP
BACKEND_URL=http://[EC2_PUBLIC_IP]:3000

# Printer Configuration
DEFAULT_BILL_PRINTER_IP=192.168.1.115
DEFAULT_BILL_PRINTER_PORT=9100
DEFAULT_LABEL_PRINTER_IP=192.168.1.101
DEFAULT_LABEL_PRINTER_PORT=9100

# Logging
LOG_LEVEL=info
```

### EC2 Security Group
```
Inbound Rules:
- Port 80 (HTTP) - Frontend
- Port 443 (HTTPS) - Frontend (nếu có SSL)
- Port 3000 (HTTP) - Backend API + WebSocket
  * Frontend sẽ kết nối WebSocket qua port này
  * Print Bridge sẽ gọi REST API qua port này
```

## Tóm tắt

| Kết nối | Protocol | Port | Mục đích | Chiều |
|---------|----------|------|----------|-------|
| Frontend → Backend | WebSocket | 3000 | Real-time notifications | Inbound to EC2 |
| Frontend → Backend | HTTP | 3000 | REST API calls | Inbound to EC2 |
| Print Bridge → Backend | WebSocket | 3000 | Nhận job real-time | Outbound from cafe |
| Print Bridge → Backend | HTTP | 3000 | Cập nhật trạng thái | Outbound from cafe |
| Backend → Print Bridge | HTTP | 3001 | Manual print (optional) | Inbound to cafe |
| Print Bridge → Printers | Raw TCP | 9100 | In nhiệt | LAN only |

**Kết luận**: 
- Port 3000 PHẢI EXPOSE cho cả WebSocket và HTTP
- Print Bridge kết nối RA NGOÀI (outbound) nên không cần mở port 3001 trên router cafe
- WebSocket giúp in nhanh hơn và tiết kiệm tài nguyên so với polling
