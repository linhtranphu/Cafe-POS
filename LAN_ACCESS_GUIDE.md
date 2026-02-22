# LAN Access Guide - Truy cập từ máy khác trong mạng

## Tổng quan

Frontend có thể được truy cập từ các máy khác trong cùng mạng LAN (ví dụ: điện thoại, tablet, máy tính khác).

## Cách sử dụng

### 1. Start services

```bash
./restart_local.sh
```

Script sẽ hiển thị IP address của máy host:

```
🌐 Access Information:
  Frontend (Local):  http://localhost:5173
  Frontend (LAN):    http://192.168.1.100:5173  ← Dùng IP này
  Backend:           http://localhost:3000
  Print Bridge:      http://localhost:3001/health

📱 Access from other devices in LAN:
  Open browser and go to: http://192.168.1.100:5173
```

### 2. Truy cập từ máy khác

Trên máy khác trong cùng mạng:
1. Mở browser (Chrome, Safari, Firefox, etc.)
2. Nhập địa chỉ: `http://192.168.1.100:5173` (thay IP bằng IP của máy host)
3. Frontend sẽ load và hoạt động bình thường

## Kiến trúc

```
┌─────────────────────────────────────────────┐
│         Máy Host (Dev Machine)              │
│                                             │
│  ┌──────────────┐      ┌──────────────┐   │
│  │  Frontend    │      │   Backend    │   │
│  │  Port 5173   │─────▶│   Port 3000  │   │
│  │  (LAN OK)    │      │  (localhost) │   │
│  └──────────────┘      └──────────────┘   │
│         │                                   │
└─────────┼───────────────────────────────────┘
          │
          │ LAN (192.168.1.x)
          │
    ┌─────▼──────┐
    │  Máy khác  │
    │  (Phone,   │
    │   Tablet,  │
    │   Laptop)  │
    └────────────┘
```

## Cách hoạt động

1. **Frontend**: Chạy với `--host` flag, cho phép bind vào tất cả network interfaces (0.0.0.0)
2. **Backend**: Vẫn chỉ chạy trên localhost (127.0.0.1)
3. **API Calls**: Frontend gọi API qua `http://localhost:3000` từ browser

## Lưu ý quan trọng

### Backend không expose ra LAN

Backend chỉ chạy trên localhost, không thể truy cập từ máy khác. Điều này là:
- ✅ **An toàn**: Backend không bị expose ra mạng
- ✅ **Đúng**: Frontend gọi API từ browser, không phải từ server

### API calls từ browser

Khi truy cập frontend từ máy khác:
- Browser load frontend từ `http://192.168.1.100:5173`
- Frontend (JavaScript) chạy trong browser của máy khác
- API calls được gửi từ browser đó đến `http://localhost:3000`
- **Vấn đề**: `localhost` trên máy khác không phải là máy host!

### Giải pháp cho production

Để frontend hoạt động đúng khi truy cập từ LAN, cần:

1. **Option 1: Sử dụng IP của máy host**
   ```bash
   # frontend/.env.local
   VITE_API_URL=http://192.168.1.100:3000
   ```
   
   Nhưng backend cũng cần expose port 3000:
   ```bash
   # Backend cần bind vào 0.0.0.0 thay vì localhost
   ```

2. **Option 2: Sử dụng proxy** (Khuyến nghị cho dev)
   Vite dev server có thể proxy API requests:
   
   ```javascript
   // frontend/vite.config.js
   export default {
     server: {
       host: true,
       proxy: {
         '/api': {
           target: 'http://localhost:3000',
           changeOrigin: true
         }
       }
     }
   }
   ```

## Setup cho LAN access hoàn chỉnh

### Bước 1: Cấu hình Vite proxy

Tạo/cập nhật `frontend/vite.config.js`:

```javascript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    host: true, // Listen on all network interfaces
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
        secure: false
      },
      '/socket.io': {
        target: 'http://localhost:3000',
        changeOrigin: true,
        ws: true // Enable WebSocket proxy
      }
    }
  }
})
```

### Bước 2: Cập nhật frontend API URL

```bash
# frontend/.env
VITE_API_URL=  # Leave empty to use relative URLs
```

Hoặc trong code:

```javascript
// Use relative URL instead of absolute
const API_URL = '/api'  // Will be proxied to http://localhost:3000/api
```

### Bước 3: Restart frontend

```bash
./restart_local.sh
```

## Testing

### Test từ máy host

```bash
curl http://localhost:5173
curl http://192.168.1.100:5173
```

### Test từ máy khác

1. Mở browser
2. Truy cập `http://192.168.1.100:5173`
3. Kiểm tra Network tab trong DevTools
4. API calls nên đi đến `/api/...` (relative URL)
5. Vite proxy sẽ forward đến `http://localhost:3000/api/...`

## Troubleshooting

### Không truy cập được từ máy khác

**Kiểm tra firewall:**
```bash
# macOS - Allow port 5173
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --add /usr/local/bin/node
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --unblockapp /usr/local/bin/node
```

**Kiểm tra network:**
```bash
# Ping máy host từ máy khác
ping 192.168.1.100

# Check port 5173 open
nc -zv 192.168.1.100 5173
```

### API calls fail (CORS error)

**Giải pháp**: Sử dụng Vite proxy (xem Setup ở trên)

### WebSocket không kết nối

**Giải pháp**: Đảm bảo proxy WebSocket được cấu hình:
```javascript
proxy: {
  '/socket.io': {
    target: 'http://localhost:3000',
    ws: true
  }
}
```

## Use Cases

### Development

- Test trên nhiều devices (phone, tablet)
- Demo cho team members
- Test responsive design
- Test trên real devices

### Production

Không nên dùng setup này cho production. Thay vào đó:
- Deploy frontend và backend lên server
- Sử dụng domain name
- Cấu hình HTTPS
- Sử dụng reverse proxy (nginx)

## Security Notes

⚠️ **Chỉ dùng trong mạng LAN tin cậy**
- Không expose ra internet
- Không dùng cho production
- Backend vẫn chỉ chạy localhost (an toàn)
- Frontend có thể truy cập từ LAN (cần cẩn thận)

## Summary

✅ Frontend có thể truy cập từ LAN với `--host` flag
✅ Backend vẫn an toàn trên localhost
⚠️ Cần Vite proxy để API calls hoạt động đúng
⚠️ Chỉ dùng cho development, không dùng production
