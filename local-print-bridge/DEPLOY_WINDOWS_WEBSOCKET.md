# Hướng Dẫn Deploy Local Print Bridge với WebSocket

## Tổng quan

Local Print Bridge với WebSocket giúp:
- ⚡ **In nhanh hơn**: Nhận job real-time qua WebSocket (không cần polling)
- 💰 **Tiết kiệm bandwidth**: Không cần check liên tục
- 🎯 **Chính xác hơn**: Không bỏ sót job

## Kiến trúc

```
Internet
   │
   │ WebSocket (real-time) + HTTP (status update)
   ▼
Windows PC (Quán Cafe)
   │
   ├─ Local Print Bridge (Docker Container)
   │  ├─ ✅ WebSocket Client → Nhận job real-time
   │  ├─ ✅ HTTP Client → Cập nhật status
   │  └─ ✅ HTTP Server → Manual print API
   │
   │ Raw TCP (Port 9100)
   ▼
Máy in nhiệt (LAN)
   ├─ Bill Printer: 192.168.1.115
   └─ Label Printer: 192.168.1.101
```

## Yêu Cầu

- Windows 10/11
- Docker Desktop for Windows
- Kết nối internet ổn định (cho WebSocket)
- Kết nối LAN đến máy in

## Bước 1: Cài Đặt Docker Desktop

1. Download: https://www.docker.com/products/docker-desktop/
2. Cài đặt và khởi động Docker Desktop
3. Đảm bảo Docker đang chạy

## Bước 2: Tạo File Cấu Hình

Tạo thư mục `C:\cafe-print-bridge` và file `.env`:

```bash
# Server Configuration
PORT=3001

# Backend URL (REQUIRED)
# Used for:
# - WebSocket connection (real-time job notifications)
# - HTTP API calls (status updates)
BACKEND_URL=http://YOUR_EC2_IP:3000

# Printer Configuration
PRINTER_TIMEOUT=5000
DEFAULT_BILL_PRINTER_IP=192.168.1.115
DEFAULT_LABEL_PRINTER_IP=192.168.1.101

# Node Environment
NODE_ENV=production
```

**Cập nhật:**
- `YOUR_EC2_IP` → IP hoặc domain của EC2 server
- `DEFAULT_BILL_PRINTER_IP` → IP máy in hóa đơn
- `DEFAULT_LABEL_PRINTER_IP` → IP máy in nhãn

## Bước 3: Chạy Print Bridge

Mở PowerShell và chạy:

```powershell
cd C:\cafe-print-bridge

docker run -d `
  --name local-print-bridge `
  --restart unless-stopped `
  --network host `
  --env-file .env `
  linhtranphu/local-print-bridge:latest
```

## Bước 4: Kiểm Tra WebSocket Connection

### 4.1. Xem logs

```powershell
docker logs -f local-print-bridge
```

Bạn sẽ thấy:
```
🖨️  Local Print Bridge Started
Server running on: http://localhost:3001
Backend URL: http://YOUR_EC2_IP:3000
[WebSocket] Connecting to backend: http://YOUR_EC2_IP:3000
[WebSocket] ✅ Connected to backend
Ready to accept print requests!
```

### 4.2. Test health check

```powershell
curl http://localhost:3001/health
```

Kết quả:
```json
{
  "status": "ok",
  "service": "Local Print Bridge",
  "version": "1.0.0"
}
```

### 4.3. Test WebSocket

1. Tạo order trong frontend
2. Xem logs Print Bridge:
   ```
   [WebSocket] 📨 New print job received: {job_id}
   [PrintJobHandler] Processing job {job_id}
   [PrintJobHandler] ✅ Job printed successfully
   [PrintJobHandler] ✅ Backend updated - Job {job_id} -> COMPLETED
   ```

## Bước 5: Troubleshooting

### WebSocket không kết nối được

**Triệu chứng:**
```
[WebSocket] Connection error: Error: timeout
[WebSocket] Reconnection error
```

**Giải pháp:**
1. Kiểm tra BACKEND_URL trong `.env` đúng chưa
2. Kiểm tra EC2 có expose port 3000 không
3. Kiểm tra firewall Windows có block không
4. Test kết nối: `curl http://YOUR_EC2_IP:3000/api/login`

### WebSocket bị disconnect liên tục

**Triệu chứng:**
```
[WebSocket] Disconnected: transport close
[WebSocket] Reconnecting...
```

**Giải pháp:**
1. Kiểm tra kết nối internet ổn định không
2. Kiểm tra EC2 backend có đang chạy không
3. Restart Print Bridge: `docker restart local-print-bridge`

### Không nhận được print job

**Triệu chứng:**
- WebSocket connected nhưng không có message

**Giải pháp:**
1. Kiểm tra backend có broadcast event không (xem backend logs)
2. Kiểm tra printer IP đúng chưa
3. Test manual print: 
   ```powershell
   curl -X POST http://localhost:3001/print `
     -H "Content-Type: application/json" `
     -d '{\"jobId\":\"test\",\"content\":\"Test\",\"printerIP\":\"192.168.1.115\"}'
   ```

## Bước 6: Cập Nhật Image Mới

Khi có version mới:

```powershell
# Stop container
docker stop local-print-bridge
docker rm local-print-bridge

# Pull image mới
docker pull linhtranphu/local-print-bridge:latest

# Chạy lại (dùng lệnh ở Bước 3)
docker run -d ...
```

## Bước 7: Auto-start khi Windows khởi động

Docker Desktop sẽ tự động start container khi Windows khởi động nhờ flag `--restart unless-stopped`.

Để đảm bảo:
1. Mở Docker Desktop Settings
2. General → "Start Docker Desktop when you log in" ✅
3. Resources → "Start Docker when Windows starts" ✅

## So sánh: Polling vs WebSocket

| Tính năng | Polling (Cũ) | WebSocket (Mới) |
|-----------|--------------|-----------------|
| Tốc độ in | 10-30 giây delay | < 1 giây (real-time) |
| Bandwidth | Cao (check liên tục) | Thấp (chỉ khi có job) |
| CPU usage | Cao | Thấp |
| Độ tin cậy | Có thể bỏ sót | Không bỏ sót |
| Kết nối | HTTP polling | WebSocket persistent |

## Lưu ý quan trọng

1. **Port 3000 phải expose trên EC2** cho cả WebSocket và HTTP
2. **Kết nối internet phải ổn định** để WebSocket hoạt động tốt
3. **Không cần mở port 3001 trên router** vì Print Bridge kết nối RA NGOÀI (outbound)
4. **WebSocket tự động reconnect** khi mất kết nối
5. **HTTP POST /print vẫn hoạt động** cho manual/reprint

## Hỗ trợ

Nếu gặp vấn đề:
1. Xem logs: `docker logs -f local-print-bridge`
2. Kiểm tra health: `curl http://localhost:3001/health`
3. Test backend: `curl http://YOUR_EC2_IP:3000/api/login`
4. Restart: `docker restart local-print-bridge`
