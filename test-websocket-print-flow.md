# Test WebSocket Print Flow - Kết Quả

## Test Date: 2026-02-21

## ✅ Test 1: HTTP Direct Print
**Status**: SUCCESS

**Command**:
```bash
curl -X POST http://localhost:3001/print \
  -H "Content-Type: application/json" \
  -d '{
    "jobId": "test-001",
    "content": "=== TEST PRINT ===\nWebSocket Test\nTest from API\n",
    "printerIP": "192.168.1.115",
    "printerPort": 9100,
    "type": "TEST"
  }'
```

**Response**:
```json
{
  "success": true,
  "jobId": "test-001",
  "message": "Print completed successfully",
  "timestamp": "2026-02-21T15:52:06.204Z"
}
```

**Print Bridge Logs**:
```
[2026-02-21T15:52:06.170Z] [INFO ] Print request received - Job ID: test-001, Printer: 192.168.1.115:9100, Type: TEST
[2026-02-21T15:52:06.203Z] [INFO ] Print successful - Job ID: test-001
[2026-02-21T15:52:06.204Z] [INFO ] Backend updated - Job ID: test-001 -> COMPLETED
```

**Kết luận**: ✅ HTTP endpoint hoạt động tốt, máy in nhận được lệnh in

---

## ✅ Test 2: WebSocket Connection
**Status**: SUCCESS

**Print Bridge Logs**:
```
[2026-02-21T15:48:26.739Z] [INFO ] [WebSocket] Connecting to: http://localhost:3000
[2026-02-21T15:48:26.783Z] [INFO ] [WebSocket] ✅ Connected to backend
```

**Backend Logs**:
```
2026/02/21 22:48:26 [Socket.IO] Client connected: <socket-id>
```

**Kết luận**: ✅ WebSocket connection thành công giữa Print Bridge và Backend

---

## ⏳ Test 3: WebSocket Print Job Broadcast
**Status**: PENDING - Cần test với real order

**Cách test**:
1. Mở frontend: http://localhost:5173
2. Login với admin/admin123
3. Tạo một order mới
4. Xem Print Bridge logs để verify nhận được event `print-job-created`

**Hoặc test qua backend API**:
```bash
# Cần implement endpoint để trigger broadcast
# Backend sẽ broadcast event khi có order mới
```

**Expected Print Bridge Logs**:
```
[WebSocket] 📨 New print job received: { job: { id: '...', content: '...', ... } }
[PrintJobHandler] Processing job <id> - Type: ORDER, Printer: <ip>:<port>
[PrintJobHandler] ✅ Job <id> printed successfully
[PrintJobHandler] ✅ Backend updated - Job <id> -> COMPLETED
```

---

## Tổng Kết

### ✅ Đã Hoạt Động
1. Print Bridge HTTP endpoint (`/print`) - ✅ Working
2. WebSocket connection giữa Print Bridge và Backend - ✅ Connected
3. Printer communication (192.168.1.115:9100) - ✅ Working

### ⏳ Cần Test Thêm
1. WebSocket broadcast từ backend khi có order mới
2. Print Bridge nhận event qua WebSocket
3. End-to-end flow: Order → Backend → WebSocket → Print Bridge → Printer

### 🎯 Next Steps
1. Tạo order từ frontend để test WebSocket broadcast
2. Hoặc implement test endpoint ở backend để trigger broadcast manually
3. Verify máy in thực sự in ra giấy

### 📊 Success Criteria
- [x] Print Bridge starts successfully
- [x] WebSocket connects to backend
- [x] HTTP print endpoint works
- [x] Printer receives print commands
- [ ] WebSocket broadcast works
- [ ] Print Bridge receives WebSocket events
- [ ] End-to-end order → print flow works

---

## Commands Reference

### Start Services
```bash
./restart_local.sh
```

### Test HTTP Print
```bash
curl -X POST http://localhost:3001/print \
  -H "Content-Type: application/json" \
  -d @test-print-payload.json
```

### Watch Logs
```bash
# Print Bridge (if running locally)
# Check terminal output

# Backend
tail -f backend.log | grep -i socket

# Print Bridge (if Docker)
docker logs -f local-print-bridge
```

### Test Printer Connection
```bash
nc -zv 192.168.1.115 9100
```

---

## Notes

1. **WebSocket Version Compatibility**: 
   - Backend: go-socket.io v1.7.0 (Engine.IO v3) ✅
   - Client: socket.io-client v2.5.0 (Engine.IO v3) ✅
   - Compatible!

2. **Backend URL Configuration**:
   - Print Bridge `.env` has `BACKEND_URL=http://localhost:3000` ✅
   - WebSocket connects successfully ✅

3. **Printer Configuration**:
   - Bill Printer: 192.168.1.115:9100 ✅
   - Label Printer: 192.168.1.101:9100 ✅

4. **Known Issues**:
   - Backend logs show some clients connecting with EIO=4 then disconnecting
   - This might be from browser/frontend, not Print Bridge
   - Print Bridge itself connects successfully with EIO=3

---

**Overall Status**: 🟢 WebSocket infrastructure working, ready for end-to-end test
