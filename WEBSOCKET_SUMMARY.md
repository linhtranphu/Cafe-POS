# WebSocket Implementation - Quick Summary

## Đã làm gì?

Thêm WebSocket vào Print Bridge để nhận print job real-time từ Backend.

## Thay đổi chính

### Print Bridge
- ✅ Thêm WebSocket client (`socket.io-client`)
- ✅ Auto-connect khi start
- ✅ Auto-reconnect khi mất kết nối
- ✅ Nhận event `print-job-created` và tự động in

### Backend
- ✅ Broadcaster gửi thêm `printer_ip` và `printer_port` trong event
- ✅ Wire up printer repository

## Ưu điểm

| Trước (Polling) | Sau (WebSocket) |
|-----------------|-----------------|
| 10-30s delay | < 1s real-time |
| Bandwidth cao | Bandwidth thấp |
| Có thể bỏ sót | Không bỏ sót |

## Câu trả lời: Port 3000 có cần expose không?

**CÓ** - Port 3000 vẫn cần expose vì:
1. WebSocket dùng cùng port với HTTP (port 3000)
2. Print Bridge cần cập nhật status về Backend qua HTTP
3. Frontend cũng dùng WebSocket trên port 3000

## Luồng hoạt động

```
User tạo order
    ↓
Backend tạo print job
    ↓
Backend broadcast qua WebSocket → Print Bridge nhận real-time
    ↓
Print Bridge in ngay
    ↓
Print Bridge update status về Backend (HTTP)
    ↓
Backend broadcast status → Frontend hiển thị notification
```

## Deploy

### 1. Build Docker image mới
```bash
cd local-print-bridge
./build-print-bridge-docker.sh 1.1.0
```

### 2. Update Print Bridge tại quán
```bash
docker stop local-print-bridge
docker rm local-print-bridge
docker pull linhtranphu/local-print-bridge:latest
docker run -d --name local-print-bridge --restart unless-stopped --network host --env-file .env.production linhtranphu/local-print-bridge:latest
```

### 3. Check logs
```bash
docker logs -f local-print-bridge
```

Expect:
```
[WebSocket] Connecting to backend: http://YOUR_EC2_IP:3000
[WebSocket] ✅ Connected to backend
```

## Files quan trọng

- `WEBSOCKET_ARCHITECTURE.md` - Kiến trúc chi tiết
- `WEBSOCKET_PRINT_IMPLEMENTATION.md` - Implementation details
- `local-print-bridge/README.md` - Print Bridge docs
- `local-print-bridge/DEPLOY_WINDOWS_WEBSOCKET.md` - Deploy guide

## Test

Tạo order trong Frontend và xem logs Print Bridge:
```
[WebSocket] 📨 New print job received: {job_id}
[PrintJobHandler] ✅ Job printed successfully
```

Done! 🎉
