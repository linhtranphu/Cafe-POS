# Fix WebSocket Connection trên EC2

## Vấn đề

Khi deploy lên EC2, WebSocket không thể kết nối với các lỗi:
```
WebSocket connection to 'wss://tacafe.store/socket.io/?EIO=3&transport=websocket' failed: 
WebSocket is closed before the connection is established.
[WebSocket] Connection error: timeout
```

## Nguyên nhân

Nginx config thiếu proxy cho Socket.IO endpoint `/socket.io/`. Hiện tại chỉ có `/api/` được proxy sang backend.

## Giải pháp

### 1. Đã fix nginx.conf

File `frontend/nginx.conf` đã được thêm config:

```nginx
# Socket.IO WebSocket proxy
location /socket.io/ {
    proxy_pass http://backend:3000;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_set_header Host $host;
    proxy_cache_bypass $http_upgrade;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_read_timeout 86400;
}
```

### 2. Deploy lên EC2

Chạy script tự động:
```bash
./fix-websocket-ec2.sh
```

Hoặc thực hiện thủ công:

#### Bước 1: Rebuild frontend (local)
```bash
cd frontend
docker build --no-cache -t linhtranphu/cafe-pos-frontend:latest .
cd ..
```

#### Bước 2: Push lên Docker Hub (local)
```bash
docker push linhtranphu/cafe-pos-frontend:latest
```

#### Bước 3: Update trên EC2
```bash
# SSH vào EC2
ssh -i your-key.pem ubuntu@tacafe.store

# Pull image mới
docker pull linhtranphu/cafe-pos-frontend:latest

# Restart frontend
docker-compose -f docker-compose.prod.yml up -d --force-recreate frontend

# Kiểm tra logs
docker logs -f cafe-pos-frontend
```

### 3. Verify

Mở https://tacafe.store và kiểm tra browser console:
- ✅ Thấy: `[WebSocket] Connected`
- ❌ Không còn: `WebSocket is closed before the connection is established`

## Vấn đề phụ: Print Bridge CORS

Lỗi thứ 2 trong log:
```
Access to fetch at 'http://192.168.1.19:3001/health' from origin 'http://52.77.228.154' 
has been blocked by CORS policy: The request client is not a secure context and the 
resource is in more-private address space `local`.
```

### Nguyên nhân
Browser hiện đại (Chrome 94+) chặn request từ public IP (HTTP) sang local network vì lý do bảo mật.

### Giải pháp

**Option 1: Sử dụng HTTPS cho EC2 (Khuyến nghị)**
```bash
# Cài đặt SSL certificate cho tacafe.store
# Browser sẽ cho phép HTTPS → Local HTTP
```

**Option 2: Truy cập qua local network**
```
# Thay vì truy cập qua 52.77.228.154
# Truy cập qua IP local: http://192.168.1.X
```

**Option 3: Disable security trong Chrome (Chỉ để test)**
```bash
# macOS
open -na "Google Chrome" --args --disable-web-security --user-data-dir="/tmp/chrome_dev"

# Windows
chrome.exe --disable-web-security --user-data-dir="C:\tmp\chrome_dev"
```

**Option 4: Cloudflare Tunnel cho Print Bridge**
Expose print bridge qua HTTPS để tránh CORS:
```bash
# Trong local-print-bridge
cloudflare tunnel --url http://localhost:3001
```

## Tóm tắt

1. ✅ WebSocket fix: Thêm `/socket.io/` proxy trong nginx → rebuild → redeploy
2. ⚠️ Print Bridge CORS: Cần HTTPS hoặc truy cập qua local network
3. 🔧 Script: `./fix-websocket-ec2.sh` để tự động hóa

## Testing

Sau khi deploy, test các chức năng realtime:
- Tạo order mới → Kiểm tra Kitchen Display cập nhật realtime
- Thay đổi trạng thái order → Kiểm tra notification
- Print job → Kiểm tra status update realtime
