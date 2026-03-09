# 🔧 Hướng dẫn sửa lỗi "Failed to fetch dynamically imported module"

## Nguyên nhân
Lỗi này xảy ra khi:
- Deploy phiên bản mới nhưng browser vẫn cache phiên bản cũ
- File JS/CSS đã bị thay đổi hash nhưng index.html cũ vẫn reference hash cũ
- Browser đã load index.html cũ và cố gắng load các module với hash không tồn tại

## Giải pháp cho User (Người dùng)

### 1. Hard Refresh Browser (Khuyến nghị)
**Chrome/Edge/Firefox:**
- Mac: `Cmd + Shift + R`
- Windows/Linux: `Ctrl + Shift + R`

**Safari:**
- Mac: `Cmd + Option + R`

### 2. Clear Cache và Reload
1. Mở DevTools (F12)
2. Right-click vào nút Reload
3. Chọn "Empty Cache and Hard Reload"

### 3. Clear Browser Data (Nếu vẫn lỗi)
**Chrome/Edge:**
1. Settings → Privacy and Security → Clear browsing data
2. Chọn "Cached images and files"
3. Time range: "Last hour"
4. Click "Clear data"

**Safari:**
1. Safari → Settings → Advanced
2. Check "Show Develop menu"
3. Develop → Empty Caches

## Giải pháp cho Developer (Triển khai)

### 1. Redeploy Frontend
```bash
# Rebuild và push frontend mới
./redeploy-frontend.sh
```

### 2. Trên Production Server
```bash
# Pull image mới
docker-compose pull frontend

# Restart frontend container
docker-compose up -d frontend

# Cleanup
docker system prune -f
```

### 3. Verify Deployment
```bash
# Check container logs
docker logs cafe-pos-frontend

# Check nginx is serving new files
docker exec cafe-pos-frontend ls -la /usr/share/nginx/html/assets/
```

## Phòng ngừa

### 1. Nginx Configuration
File `frontend/nginx.conf` đã được cấu hình:
- HTML files: KHÔNG cache (always fresh)
- JS/CSS with hash: Cache 1 year (immutable)
- Other assets: Cache 30 days

### 2. Service Worker
Nếu app sử dụng Service Worker, cần:
```javascript
// Unregister old service worker
navigator.serviceWorker.getRegistrations().then(registrations => {
  registrations.forEach(registration => registration.unregister())
})
```

### 3. Meta Tags
Thêm vào `index.html`:
```html
<meta http-equiv="Cache-Control" content="no-cache, no-store, must-revalidate">
<meta http-equiv="Pragma" content="no-cache">
<meta http-equiv="Expires" content="0">
```

## Kiểm tra Cache Headers

### Sử dụng curl
```bash
# Check index.html (should NOT cache)
curl -I https://tacafe.store/

# Check JS file (should cache)
curl -I https://tacafe.store/assets/index-3131f1d4.js
```

### Sử dụng Browser DevTools
1. Mở DevTools (F12)
2. Network tab
3. Reload page
4. Click vào file
5. Xem Headers → Response Headers → Cache-Control

## Troubleshooting

### Lỗi vẫn còn sau khi clear cache?
1. Check xem có đang dùng proxy/CDN không (Cloudflare, etc.)
2. Clear CDN cache nếu có
3. Check service worker: DevTools → Application → Service Workers
4. Try incognito/private mode

### Multiple users báo lỗi?
→ Cần redeploy frontend mới (xem phần Developer ở trên)

### Chỉ 1 user báo lỗi?
→ User cần clear cache (xem phần User ở trên)

## Scripts hữu ích

### Rebuild toàn bộ
```bash
./build_docker_hub.sh
```

### Chỉ rebuild frontend
```bash
./redeploy-frontend.sh
```

### Check Docker disk usage
```bash
docker system df
```

### Cleanup Docker
```bash
docker system prune -a --volumes
```
