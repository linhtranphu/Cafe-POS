# ✅ WebSocket Removed - HTTP Only Architecture

## Thay đổi

Đã chuyển từ kiến trúc WebSocket sang HTTP-only để đơn giản hóa hệ thống.

## Kiến trúc cũ (WebSocket - Phức tạp)

```
Frontend → Backend (EC2) ←WebSocket→ Print Bridge (Local) → Printer
                ↓
         Socket.IO Server
```

**Vấn đề:**
- Cần maintain WebSocket connection
- Phức tạp với reconnection logic
- Timeout errors khi network không ổn định
- Cần expose WebSocket port qua firewall

## Kiến trúc mới (HTTP Only - Đơn giản)

```
Frontend → Backend (EC2) --HTTP POST--> Print Bridge (Local) → Printer
                                        ↓
                                   Puppeteer/Chromium
                                        ↓
                                   HTML → ESC/POS
```

**Ưu điểm:**
- Đơn giản, dễ debug
- Không cần maintain persistent connection
- HTTP request/response rõ ràng
- Dễ monitor và log

## Files đã thay đổi

### 1. Print Bridge (local-print-bridge/)

**src/index.js:**
- ❌ Removed: `websocketClient`, `backendSync`, `printJobHandler`
- ✅ Kept: `htmlRenderer` (core functionality)
- ✅ Simplified `/print` endpoint (no job tracking)
- ✅ Kept `/render-html` endpoint (main feature)

**Endpoints còn lại:**
```
POST /render-html      - Render HTML to ESC/POS (MAIN)
POST /print            - Direct print ESC/POS data
POST /test-connection  - Test printer connection
GET  /health           - Health check
GET  /status           - Service status
```

**.env.example:**
- ❌ Removed: `BACKEND_URL` (không cần nữa)
- ✅ Kept: Printer IPs, ports, timeouts

### 2. Frontend (frontend/)

**src/main.js:**
- ❌ Disabled: `websocketService.connect()`
- ❌ Disabled: WebSocket auth subscription
- ✅ Kept: `localPrintService` (HTTP client)

**Files không cần nữa (có thể xóa sau):**
- `src/services/websocket.js`
- `src/composables/usePrintJobWebSocket.js`

### 3. Backend (backend/)

**Không thay đổi** - Backend đã có HTTP client sẵn:
- `infrastructure/printbridge/client.go` - HTTP POST client
- `interfaces/http/html_template_handler_bridge.go` - Handler

## Flow hoạt động

### HTML Template Printing:

1. **User clicks "Test Print" trong Settings**
   ```
   Frontend → POST /api/manager/html-templates/test-print
   ```

2. **Backend nhận request**
   ```go
   // html_template_handler_bridge.go
   func (h *HTMLTemplateHandlerBridge) TestPrintHTMLTemplate(c *gin.Context)
   ```

3. **Backend gọi Print Bridge**
   ```go
   // printbridge/client.go
   escposData := client.RenderHTMLToESCPOS(ctx, html, 576)
   // HTTP POST http://localhost:3001/render-html
   ```

4. **Print Bridge render HTML**
   ```javascript
   // htmlRenderer.js
   Puppeteer → HTML → Screenshot → Grayscale → ESC/POS
   ```

5. **Backend nhận ESC/POS data và in**
   ```go
   // Send to printer via TCP
   conn.Write(escposData)
   ```

## Deployment Steps

### 1. Deploy Print Bridge (Local Machine)

```bash
cd local-print-bridge

# Update .env (remove BACKEND_URL)
cp .env.example .env
nano .env
# Set printer IPs only

# Install dependencies
npm install

# Start service
npm start
# Or with PM2:
pm2 start src/index.js --name print-bridge
```

**Expected logs:**
```
🖨️  Local Print Bridge Started (HTTP Only)
Server running on: http://localhost:3001
Endpoints:
  POST /render-html - Render HTML to ESC/POS
  POST /print - Direct print ESC/POS data
  ...
Ready to accept requests!
```

### 2. Deploy Backend (EC2)

```bash
# SSH to EC2
ssh ubuntu@tacafe.store

# Update .env
cd ~/cafe-pos
nano .env

# Add this line:
PRINT_BRIDGE_URL=http://192.168.1.X:3001
# Replace X with your local machine IP

# Rebuild and restart backend
docker-compose build backend
docker-compose up -d backend

# Check logs
docker logs -f cafe-pos-backend
```

**Expected logs:**
```
🔗 Initializing Print Bridge client: http://192.168.1.X:3001
✅ Print Bridge connected successfully
✅ HTML template handler (bridge) initialized
✅ HTML template routes registered (using print bridge)
```

### 3. Deploy Frontend (EC2)

```bash
# Frontend không cần thay đổi gì
# WebSocket code đã bị disable, không ảnh hưởng

# Nếu muốn rebuild:
docker-compose build frontend
docker-compose up -d frontend
```

## Testing

### 1. Test Print Bridge

```bash
# On local machine
curl http://localhost:3001/health
# Should return: {"status":"ok",...}

# Test HTML rendering
curl -X POST http://localhost:3001/test-render
# Should return: {"success":true,...}
```

### 2. Test Backend → Print Bridge

```bash
# From EC2 or local
curl -X POST https://tacafe.store/api/manager/html-templates/test-print \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### 3. Test from Frontend

1. Login to https://tacafe.store
2. Go to Settings → Print Settings
3. Click "Test Print"
4. Should see success message
5. Check printer output

## Troubleshooting

### 404 on /api/manager/html-templates/bill

**Cause:** Backend không có `PRINT_BRIDGE_URL` trong .env

**Fix:**
```bash
# On EC2
echo "PRINT_BRIDGE_URL=http://192.168.1.X:3001" >> ~/cafe-pos/.env
docker-compose restart backend
```

### Print Bridge not reachable

**Cause:** Network issue hoặc firewall

**Fix:**
```bash
# Test from EC2
curl http://192.168.1.X:3001/health

# If fails, check:
# 1. Print bridge is running
# 2. Firewall allows port 3001
# 3. IP address is correct
```

### WebSocket timeout errors (in browser console)

**Cause:** Frontend vẫn còn WebSocket code cũ

**Fix:** Đã disable trong main.js, errors sẽ biến mất sau khi rebuild frontend

## Rollback Plan

Nếu cần rollback về WebSocket:

```bash
# Frontend
git revert <commit-hash>

# Print Bridge
git revert <commit-hash>

# Backend - không cần thay đổi (vẫn support cả 2)
```

## Summary

✅ **Removed:**
- WebSocket connection logic
- Socket.IO dependencies (có thể xóa sau)
- Backend sync service
- Print job handler

✅ **Kept:**
- HTML rendering (core feature)
- HTTP endpoints
- Printer service
- Direct print capability

✅ **Benefits:**
- Simpler architecture
- Easier to debug
- No connection management
- Clear request/response flow
- Less network complexity

---

**Status:** Ready for deployment  
**Architecture:** HTTP Only  
**Next:** Deploy and test end-to-end
