# Print Bridge Configuration

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Internet                              │
│                                                              │
│  ┌────────────────────────────────────────────────────┐     │
│  │              EC2 Server (tacafe.store)             │     │
│  │                                                     │     │
│  │  ┌──────────────┐         ┌──────────────┐        │     │
│  │  │   Frontend   │         │   Backend    │        │     │
│  │  │  (Static)    │         │   (API)      │        │     │
│  │  └──────────────┘         └──────────────┘        │     │
│  └────────────────────────────────────────────────────┘     │
│                                                              │
└──────────────────────────────────────────────────────────────┘
                              │
                              │ HTTPS
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Cafe Local Network                        │
│                                                              │
│  ┌──────────────────┐                                       │
│  │   Staff Browser  │                                       │
│  │  (Chrome/Edge)   │                                       │
│  └────────┬─────────┘                                       │
│           │                                                  │
│           │ HTTP (Local)                                    │
│           ▼                                                  │
│  ┌──────────────────┐         ┌──────────────┐             │
│  │  Print Bridge    │────────▶│  ESC/POS     │             │
│  │  localhost:3001  │         │  Printer     │             │
│  │                  │         │  :9100       │             │
│  └──────────────────┘         └──────────────┘             │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## Problem

Browser truy cập `https://tacafe.store` (HTTPS) không thể kết nối tới `http://localhost:3001` (HTTP local) vì:
1. **Mixed Content**: HTTPS không cho phép gọi HTTP
2. **CORS**: Localhost không accessible từ remote domain
3. **Security**: Browser block cross-origin requests

## Solution: Expose Print Bridge

Print bridge cần được expose qua một URL mà browser có thể truy cập.

### Option 1: Local Network IP (Simplest)

**Khi nào dùng:** Browser và Print Bridge cùng mạng local

**Setup:**

1. **Tìm IP của máy chạy print bridge:**
```bash
# Windows
ipconfig

# macOS/Linux
ifconfig
# hoặc
ip addr show

# Ví dụ: 192.168.1.100
```

2. **Cấu hình print bridge listen trên tất cả interfaces:**

Trong print bridge code, thay vì:
```javascript
app.listen(3001, 'localhost')
```

Dùng:
```javascript
app.listen(3001, '0.0.0.0')  // Listen on all interfaces
```

3. **Test từ browser:**
```
http://192.168.1.100:3001/health
```

4. **Cấu hình frontend:**

Tạo file `frontend/.env.production.local`:
```env
VITE_PRINT_BRIDGE_URL=http://192.168.1.100:3001
```

5. **Rebuild frontend:**
```bash
cd frontend
npm run build

# Build Docker image
docker build -t linhtranphu/cafe-pos-frontend:2.0.1 .
docker push linhtranphu/cafe-pos-frontend:2.0.1
```

6. **Deploy:**
```bash
# On EC2
sudo docker-compose pull frontend
sudo docker-compose up -d frontend
```

**Ưu điểm:**
- ✅ Đơn giản nhất
- ✅ Không cần service bên ngoài
- ✅ Nhanh, không qua internet

**Nhược điểm:**
- ❌ Chỉ work khi browser truy cập từ cùng mạng
- ❌ IP có thể thay đổi (cần DHCP reservation)
- ❌ Mixed content warning (HTTPS → HTTP)

---

### Option 2: Cloudflare Tunnel (Recommended)

**Khi nào dùng:** Muốn truy cập từ bất kỳ đâu, secure

**Setup:**

1. **Install Cloudflare Tunnel:**
```bash
# macOS
brew install cloudflare/cloudflare/cloudflared

# Windows: Download from
# https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/install-and-setup/installation/
```

2. **Login:**
```bash
cloudflared tunnel login
```

3. **Create tunnel:**
```bash
cloudflared tunnel create cafe-print-bridge
```

4. **Configure** (`~/.cloudflared/config.yml`):
```yaml
tunnel: <TUNNEL_ID_FROM_STEP_3>
credentials-file: /path/to/<TUNNEL_ID>.json

ingress:
  - hostname: print.tacafe.store
    service: http://localhost:3001
  - service: http_status:404
```

5. **Route DNS:**
```bash
cloudflared tunnel route dns cafe-print-bridge print.tacafe.store
```

6. **Run tunnel:**
```bash
# Test first
cloudflared tunnel run cafe-print-bridge

# Run as service (macOS)
brew services start cloudflared

# Run as service (Windows)
cloudflared service install
```

7. **Configure frontend** (`frontend/.env.production.local`):
```env
VITE_PRINT_BRIDGE_URL=https://print.tacafe.store
```

8. **Rebuild and deploy** (same as Option 1)

**Ưu điểm:**
- ✅ HTTPS (no mixed content warning)
- ✅ Truy cập từ bất kỳ đâu
- ✅ Free
- ✅ Secure

**Nhược điểm:**
- ❌ Phụ thuộc internet
- ❌ Setup phức tạp hơn

---

### Option 3: Ngrok (Quick Test)

**Khi nào dùng:** Test nhanh, demo

**Setup:**

1. **Install ngrok:**
```bash
# Download from https://ngrok.com/download
```

2. **Start tunnel:**
```bash
ngrok http 3001
```

3. **Copy URL** (e.g., `https://abc123.ngrok-free.app`)

4. **Configure frontend:**
```env
VITE_PRINT_BRIDGE_URL=https://abc123.ngrok-free.app
```

**Ưu điểm:**
- ✅ Rất nhanh để test
- ✅ HTTPS

**Nhược điểm:**
- ❌ URL thay đổi mỗi lần restart (free tier)
- ❌ Không stable cho production

---

## Recommended Setup

**For your case (Cafe POS):**

### Development (Local testing)
```env
# frontend/.env.development
VITE_PRINT_BRIDGE_URL=http://localhost:3001
```

### Production (Deployed on EC2)

**Option A: Same network access (Simplest)**
```env
# frontend/.env.production.local
VITE_PRINT_BRIDGE_URL=http://192.168.1.100:3001
```

**Option B: Remote access (Most flexible)**
```env
# frontend/.env.production.local
VITE_PRINT_BRIDGE_URL=https://print.tacafe.store
```

---

## Implementation Steps

### 1. Update Print Bridge to Listen on All Interfaces

```javascript
// print-bridge/server.js
const PORT = process.env.PORT || 3001
const HOST = process.env.HOST || '0.0.0.0'  // Listen on all interfaces

app.listen(PORT, HOST, () => {
  console.log(`Print Bridge running on ${HOST}:${PORT}`)
})
```

### 2. Add CORS Support

```javascript
// print-bridge/server.js
const cors = require('cors')

app.use(cors({
  origin: [
    'http://localhost:5173',      // Development
    'https://tacafe.store',       // Production
    /^http:\/\/192\.168\.\d+\.\d+:\d+$/  // Local network
  ],
  credentials: true
}))
```

### 3. Configure Frontend

Create `frontend/.env.production.local`:
```env
# Use local IP (Option A)
VITE_PRINT_BRIDGE_URL=http://192.168.1.100:3001

# OR use Cloudflare Tunnel (Option B)
# VITE_PRINT_BRIDGE_URL=https://print.tacafe.store
```

### 4. Rebuild Frontend

```bash
cd frontend
npm run build
docker build -t linhtranphu/cafe-pos-frontend:2.0.1 .
docker push linhtranphu/cafe-pos-frontend:2.0.1
```

### 5. Deploy

```bash
# On EC2
sudo docker-compose pull frontend
sudo docker-compose up -d frontend
```

### 6. Test

1. Open `https://tacafe.store` from cafe
2. Go to Print Management
3. Check browser console:
   - Should see: `[LocalPrint] Bridge available: true`
4. Test printer connection

---

## Troubleshooting

### "Bridge not available"

**Check 1:** Print bridge running?
```bash
curl http://localhost:3001/health
```

**Check 2:** Firewall blocking?
```bash
# Windows: Allow port 3001 in Windows Firewall
# macOS: System Preferences → Security → Firewall → Options
```

**Check 3:** Can access from browser?
```
http://192.168.1.100:3001/health
```

### Mixed Content Warning

Browser console shows: "Mixed Content: The page at 'https://tacafe.store' was loaded over HTTPS, but requested an insecure resource 'http://192.168.1.100:3001'"

**Solution:** Use HTTPS (Cloudflare Tunnel) or allow insecure content in browser settings (not recommended)

### CORS Error

**Solution:** Add CORS headers in print bridge (see step 2 above)

---

## Quick Start (Recommended)

**For immediate use with local IP:**

1. **Find print bridge machine IP:**
   ```bash
   ipconfig  # Windows
   ifconfig  # macOS/Linux
   ```

2. **Update print bridge to listen on 0.0.0.0**

3. **Create frontend/.env.production.local:**
   ```env
   VITE_PRINT_BRIDGE_URL=http://YOUR_LOCAL_IP:3001
   ```

4. **Rebuild and deploy:**
   ```bash
   ./build_docker_hub.sh
   # Then deploy on EC2
   ```

Done! Browser tại quán sẽ kết nối được với print bridge local.
