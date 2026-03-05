# Deployment Checklist - Print System với WebSocket

## Tổng quan

Checklist này giúp deploy hệ thống in với kiến trúc WebSocket mới:
- Backend (EC2) ← WebSocket → Print Bridge (Local) → Printers

## Phase 1: Backend (EC2)

### 1.1 Fix Nginx WebSocket Proxy

- [ ] Nginx config đã có `/socket.io/` location block
  ```bash
  # Kiểm tra
  grep -A 10 "location /socket.io/" frontend/nginx.conf
  ```

- [ ] Rebuild frontend image
  ```bash
  ./fix-websocket-ec2.sh
  ```

- [ ] Push image lên Docker Hub
  ```bash
  docker push linhtranphu/cafe-pos-frontend:latest
  ```

### 1.2 Deploy lên EC2

- [ ] SSH vào EC2
  ```bash
  ssh -i your-key.pem ubuntu@tacafe.store
  ```

- [ ] Pull image mới
  ```bash
  docker pull linhtranphu/cafe-pos-frontend:latest
  ```

- [ ] Restart frontend container
  ```bash
  docker-compose -f docker-compose.prod.yml up -d --force-recreate frontend
  ```

- [ ] Kiểm tra logs
  ```bash
  docker logs -f cafe-pos-frontend
  docker logs -f cafe-pos-backend
  ```

### 1.3 Verify Backend WebSocket

- [ ] Test từ browser console tại https://tacafe.store
  ```javascript
  // Mở DevTools Console
  // Kiểm tra không còn lỗi WebSocket timeout
  ```

- [ ] Test từ Print Bridge machine
  ```bash
  curl https://tacafe.store/api/state-machines
  # Should return 200 OK
  ```

### 1.4 Security Group (EC2)

- [ ] Port 80 (HTTP) - Open to 0.0.0.0/0
- [ ] Port 443 (HTTPS) - Open to 0.0.0.0/0 (nếu có SSL)
- [ ] Port 3000 (Backend) - Open to Print Bridge IP hoặc 0.0.0.0/0

## Phase 2: Print Bridge (Local Machine)

### 2.1 Cài đặt môi trường

- [ ] Node.js đã cài (v16+)
  ```bash
  node --version
  ```

- [ ] Git đã cài (để clone repo)
  ```bash
  git --version
  ```

- [ ] Clone repository
  ```bash
  git clone <repo-url>
  cd local-print-bridge
  ```

### 2.2 Cấu hình

- [ ] Copy .env.example
  ```bash
  cp .env.example .env
  ```

- [ ] Cập nhật BACKEND_URL trong .env
  ```bash
  # .env
  BACKEND_URL=https://tacafe.store
  # Hoặc: BACKEND_URL=http://52.77.228.154:3000
  ```

- [ ] Cập nhật Printer IPs
  ```bash
  # .env
  DEFAULT_BILL_PRINTER_IP=192.168.1.100
  DEFAULT_BILL_PRINTER_PORT=9100
  DEFAULT_LABEL_PRINTER_IP=192.168.1.101
  DEFAULT_LABEL_PRINTER_PORT=9100
  ```

### 2.3 Test kết nối

- [ ] Test backend reachable
  ```bash
  curl https://tacafe.store/api/state-machines
  ```

- [ ] Test WebSocket connection
  ```bash
  node test-backend-websocket.js https://tacafe.store
  # Should see: ✅ Connected to backend
  ```

- [ ] Test printer connection
  ```bash
  npm start
  # Trong terminal khác:
  curl -X POST http://localhost:3001/test-connection \
    -H "Content-Type: application/json" \
    -d '{"printerIP": "192.168.1.100"}'
  ```

### 2.4 Khởi động service

**Option A: Manual start (Development)**
- [ ] Start Print Bridge
  ```bash
  npm start
  ```
- [ ] Kiểm tra logs
  ```
  [WebSocket] ✅ Connected to backend
  Ready to accept print requests!
  ```

**Option B: PM2 (Production)**
- [ ] Install PM2
  ```bash
  npm install -g pm2
  ```
- [ ] Start with PM2
  ```bash
  pm2 start src/index.js --name print-bridge
  pm2 save
  pm2 startup
  ```
- [ ] Verify running
  ```bash
  pm2 status
  pm2 logs print-bridge
  ```

**Option C: Quick Start Script**
- [ ] Run quick start
  ```bash
  ./quick-start.sh
  ```

## Phase 3: Testing End-to-End

### 3.1 Test Print Flow

- [ ] Mở https://tacafe.store trên browser
- [ ] Login vào hệ thống
- [ ] Tạo order mới
- [ ] Kiểm tra Print Bridge logs
  ```
  [WebSocket] 📨 Received print job via WebSocket: <job-id>
  [PrintJobHandler] Processing job <job-id>
  [PrintJobHandler] ✅ Job <job-id> printed successfully
  ```
- [ ] Kiểm tra máy in có in ra không
- [ ] Kiểm tra status trong UI (order history)

### 3.2 Test Reconnection

- [ ] Restart backend
  ```bash
  # Trên EC2
  docker-compose restart backend
  ```
- [ ] Kiểm tra Print Bridge tự động reconnect
  ```
  [WebSocket] Disconnected: transport close
  [WebSocket] ✅ Reconnected after 1 attempts
  ```

### 3.3 Test Error Handling

- [ ] Tắt máy in
- [ ] Tạo order mới
- [ ] Kiểm tra Print Bridge logs
  ```
  [PrintJobHandler] ❌ Job <id> failed: connect ETIMEDOUT
  [PrintJobHandler] Backend updated - Job <id> -> FAILED
  ```
- [ ] Kiểm tra UI hiển thị lỗi

## Phase 4: Monitoring

### 4.1 Backend Monitoring

- [ ] Check backend logs
  ```bash
  docker logs -f cafe-pos-backend | grep Socket.IO
  ```
- [ ] Monitor connections
  ```
  [Socket.IO] Client connected: <socket-id>
  [Socket.IO] Broadcasted print-job-created event
  ```

### 4.2 Print Bridge Monitoring

- [ ] Check Print Bridge logs
  ```bash
  # PM2
  pm2 logs print-bridge
  
  # Manual
  # Logs in console
  ```

- [ ] Check status endpoint
  ```bash
  curl http://localhost:3001/status
  ```

- [ ] Check health endpoint
  ```bash
  curl http://localhost:3001/health
  ```

### 4.3 Network Monitoring

- [ ] Test backend từ Print Bridge
  ```bash
  ping tacafe.store
  curl https://tacafe.store/api/state-machines
  ```

- [ ] Test printers từ Print Bridge
  ```bash
  ping 192.168.1.100
  telnet 192.168.1.100 9100
  ```

## Phase 5: Documentation

- [ ] Document printer IPs và locations
- [ ] Document backend URL
- [ ] Document troubleshooting steps
- [ ] Train staff on monitoring

## Troubleshooting Quick Reference

### WebSocket không kết nối

```bash
# 1. Check backend reachable
curl https://tacafe.store/api/state-machines

# 2. Check nginx config
ssh ec2
grep -A 10 "location /socket.io/" /path/to/nginx.conf

# 3. Check backend logs
docker logs cafe-pos-backend | grep Socket.IO

# 4. Test WebSocket
node test-backend-websocket.js https://tacafe.store
```

### Máy in không in

```bash
# 1. Test printer connection
curl -X POST http://localhost:3001/test-connection \
  -H "Content-Type: application/json" \
  -d '{"printerIP": "192.168.1.100"}'

# 2. Ping printer
ping 192.168.1.100

# 3. Check printer power and network cable

# 4. Check printer IP (print test page from printer)
```

### Backend không nhận status update

```bash
# 1. Check BACKEND_URL in .env
cat .env | grep BACKEND_URL

# 2. Test backend API
curl -X PUT https://tacafe.store/api/print-jobs/test-id/status \
  -H "Content-Type: application/json" \
  -d '{"status": "COMPLETED"}'

# 3. Check Print Bridge logs
pm2 logs print-bridge | grep "Backend updated"
```

## Success Criteria

✅ Backend WebSocket server running và accessible
✅ Print Bridge connected to backend via WebSocket
✅ Print Bridge can reach printers
✅ End-to-end print flow works
✅ Auto-reconnect works after network issues
✅ Error handling works (printer offline, etc.)
✅ Monitoring in place

## Related Documentation

- `PRINT_ARCHITECTURE.md` - Architecture overview
- `PRINT_BRIDGE_WEBSOCKET_SETUP.md` - Detailed setup guide
- `WEBSOCKET_EC2_FIX.md` - WebSocket troubleshooting
- `local-print-bridge/README.md` - Print Bridge documentation
- `local-print-bridge/quick-start.sh` - Quick start script

## Support

Nếu gặp vấn đề:
1. Check logs (backend + print bridge)
2. Test network connectivity
3. Review troubleshooting section
4. Check related documentation
