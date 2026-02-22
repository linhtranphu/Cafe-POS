# WebSocket Deployment Checklist

## Pre-deployment

- [ ] Backend đã build với broadcaster mới
- [ ] Print Bridge đã install `socket.io-client`
- [ ] Docker image mới đã build và push
- [ ] Documentation đã đọc và hiểu

## Backend (EC2)

- [ ] Upload backend binary mới lên EC2
- [ ] Restart backend service
- [ ] Check logs: WebSocket hub started
  ```bash
  # Expected logs:
  ✅ WebSocket hub started
  ✅ WebSocket endpoint registered at /socket.io/
  ```
- [ ] Test WebSocket endpoint
  ```bash
  curl http://YOUR_EC2_IP:3000/socket.io/
  # Should return Socket.IO handshake response
  ```
- [ ] Verify port 3000 exposed in Security Group
  - [ ] Inbound rule: Port 3000, Source: 0.0.0.0/0

## Print Bridge (Windows PC tại quán)

- [ ] Stop container cũ
  ```bash
  docker stop local-print-bridge
  docker rm local-print-bridge
  ```
- [ ] Pull image mới
  ```bash
  docker pull linhtranphu/local-print-bridge:latest
  ```
- [ ] Cập nhật `.env.production`
  ```bash
  BACKEND_URL=http://YOUR_EC2_IP:3000
  DEFAULT_BILL_PRINTER_IP=192.168.1.115
  DEFAULT_LABEL_PRINTER_IP=192.168.1.101
  ```
- [ ] Run container mới
  ```bash
  docker run -d \
    --name local-print-bridge \
    --restart unless-stopped \
    --network host \
    --env-file .env.production \
    linhtranphu/local-print-bridge:latest
  ```
- [ ] Check logs
  ```bash
  docker logs -f local-print-bridge
  ```
- [ ] Verify WebSocket connected
  ```
  # Expected logs:
  [WebSocket] Connecting to backend: http://YOUR_EC2_IP:3000
  [WebSocket] ✅ Connected to backend
  Ready to accept print requests!
  ```

## Testing

### Test 1: WebSocket Connection
- [ ] Print Bridge logs show "✅ Connected to backend"
- [ ] No connection errors in logs

### Test 2: Auto-print
- [ ] Tạo order trong Frontend
- [ ] Print Bridge logs show:
  ```
  [WebSocket] 📨 New print job received: {job_id}
  [PrintJobHandler] Processing job {job_id}
  [PrintJobHandler] ✅ Job printed successfully
  ```
- [ ] Máy in đã in hóa đơn
- [ ] Frontend hiển thị notification "In thành công"

### Test 3: Manual print
- [ ] Click "Reprint" trong Frontend
- [ ] Print Bridge nhận job qua WebSocket
- [ ] Máy in đã in lại

### Test 4: Reconnection
- [ ] Restart backend
- [ ] Print Bridge logs show reconnection attempts
- [ ] Print Bridge reconnect thành công
- [ ] Tạo order mới vẫn in được

### Test 5: Error handling
- [ ] Tắt máy in
- [ ] Tạo order
- [ ] Print Bridge logs show error
- [ ] Backend nhận status "FAILED"
- [ ] Frontend hiển thị error notification

## Monitoring

### Backend logs
- [ ] Check WebSocket connections
  ```bash
  docker logs backend | grep WebSocket
  ```
- [ ] Check broadcast events
  ```bash
  docker logs backend | grep "Broadcasted event"
  ```

### Print Bridge logs
- [ ] Check WebSocket status
  ```bash
  docker logs local-print-bridge | grep WebSocket
  ```
- [ ] Check print jobs
  ```bash
  docker logs local-print-bridge | grep PrintJobHandler
  ```

### Frontend
- [ ] Check WebSocket connection in browser console
- [ ] Check print notifications

## Rollback Plan

Nếu có vấn đề:

### Backend
```bash
# Revert to old binary
cp backend.old backend
systemctl restart backend
```

### Print Bridge
```bash
# Use old version
docker stop local-print-bridge
docker rm local-print-bridge
docker run -d \
  --name local-print-bridge \
  --restart unless-stopped \
  --network host \
  --env-file .env.production \
  linhtranphu/local-print-bridge:1.0.0
```

## Post-deployment

- [ ] Monitor logs for 1 hour
- [ ] Test multiple orders
- [ ] Verify no errors
- [ ] Document any issues
- [ ] Update team about new feature

## Success Criteria

✅ WebSocket connected
✅ Print jobs nhận real-time (< 1s)
✅ Không có errors trong logs
✅ Máy in hoạt động bình thường
✅ Frontend notifications hoạt động
✅ Auto-reconnect hoạt động

## Notes

- WebSocket dùng cùng port 3000 với HTTP
- Print Bridge không cần mở port 3001 trên router (chỉ kết nối outbound)
- WebSocket tự động reconnect khi mất kết nối
- HTTP POST /print vẫn hoạt động cho manual print
