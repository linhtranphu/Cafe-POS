# Fix Frontend Socket.IO Version Mismatch

## Vấn Đề

Frontend đang dùng `socket.io-client v4.8.3` (Engine.IO v4) nhưng backend server dùng `go-socket.io v1.7.0` (Engine.IO v3), gây ra lỗi:

```
Error: It seems you are trying to reach a Socket.IO server in v2.x with a v3.x client, 
but they are not compatible
```

## ✅ Đã Sửa

1. Downgrade socket.io-client trong `frontend/package.json`:
   ```json
   "socket.io-client": "^2.5.0"
   ```

2. Đã chạy `npm install` trong frontend directory

## 🔄 Cần Restart Frontend

Frontend cần được restart để áp dụng version mới:

```bash
# Stop frontend hiện tại (nếu đang chạy)
# Ctrl+C trong terminal đang chạy frontend

# Hoặc kill process:
kill $(lsof -t -i:5173)

# Start lại frontend
cd frontend
npm run dev -- --host
```

## Verify Fix

Sau khi restart frontend:

1. Mở browser console (F12)
2. Không còn thấy lỗi Socket.IO version mismatch
3. WebSocket connection thành công

## Version Compatibility

| Component | Library | Version | Engine.IO |
|-----------|---------|---------|-----------|
| Backend | go-socket.io | v1.7.0 | v3 |
| Print Bridge | socket.io-client | v2.5.0 | v3 |
| Frontend | socket.io-client | v2.5.0 | v3 |

✅ Tất cả đều dùng Engine.IO v3 - Compatible!

## Print Bridge Status

Print Bridge vẫn hoạt động tốt và không bị ảnh hưởng:

```
[WebSocket] ✅ Connected to backend
```

Health check:
```bash
curl http://localhost:3001/health
# {"status":"ok","service":"Local Print Bridge","version":"1.0.0"}
```

## Next Steps

1. Restart frontend với command trên
2. Test WebSocket connection từ frontend
3. Tạo order để test end-to-end flow
4. Verify Print Bridge nhận được event qua WebSocket

## Files Changed

- `frontend/package.json` - Downgraded socket.io-client to v2.5.0
- `frontend/node_modules/` - Updated dependencies

## Notes

- Lỗi ban đầu là từ frontend (browser console), không phải Print Bridge
- Print Bridge đã kết nối WebSocket thành công từ đầu
- Chỉ cần fix frontend để hoàn thành toàn bộ WebSocket infrastructure
