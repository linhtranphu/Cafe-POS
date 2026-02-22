# Print Bridge Docker Network Fix

## Vấn Đề

Frontend hiển thị "Local Bridge Offline" mặc dù Docker container đang chạy.

## Nguyên Nhân

Docker container được start với `--network host` mode, nhưng trên **macOS**, host network mode không hoạt động như Linux. Docker Desktop trên macOS chạy trong VM, nên `host` mode không expose port ra host machine.

**Kết quả**:
- Container chạy và healthy ✅
- Port 3001 không accessible từ localhost ❌
- Frontend không thể connect đến Print Bridge ❌

## ✅ Giải Pháp

Thay `--network host` bằng port mapping `-p 3001:3001`

### Before (Không hoạt động trên macOS)
```bash
docker run -d \
    --name local-print-bridge \
    --restart unless-stopped \
    --network host \
    --env-file .env \
    linhtranphu/local-print-bridge:latest
```

### After (Hoạt động trên macOS)
```bash
docker run -d \
    --name local-print-bridge \
    --restart unless-stopped \
    -p 3001:3001 \
    --env-file local-print-bridge/.env \
    linhtranphu/local-print-bridge:latest
```

## Đã Thực Hiện

1. **Stop và remove container cũ**:
   ```bash
   docker stop local-print-bridge
   docker rm local-print-bridge
   ```

2. **Start container mới với port mapping**:
   ```bash
   docker run -d \
       --name local-print-bridge \
       --restart unless-stopped \
       -p 3001:3001 \
       --env-file local-print-bridge/.env \
       linhtranphu/local-print-bridge:latest
   ```

3. **Verify**:
   ```bash
   curl http://localhost:3001/health
   # {"status":"ok","service":"Local Print Bridge","version":"1.0.0"}
   ```

4. **Update restart_local.sh**:
   - Sửa từ `--network host` → `-p 3001:3001`
   - Fix path `.env` → `local-print-bridge/.env`

## Verify

### 1. Check Container
```bash
docker ps | grep print-bridge
# Should show: Up X minutes (healthy)
```

### 2. Test Health Endpoint
```bash
curl http://localhost:3001/health
# Should return: {"status":"ok",...}
```

### 3. Check Frontend
- Mở http://localhost:5173/#/print-management
- Phải thấy: "🟢 Local Bridge Online"

## Lưu Ý Quan Trọng

### macOS vs Linux

| Platform | Network Mode | Port Access |
|----------|--------------|-------------|
| Linux | `--network host` | ✅ Works |
| macOS | `--network host` | ❌ Doesn't work |
| macOS | `-p 3001:3001` | ✅ Works |

**Lý do**: Docker Desktop trên macOS chạy trong VM (HyperKit/QEMU), không phải native như Linux.

### Production (Windows)

Trên Windows cũng cần dùng port mapping:
```bash
docker run -d ^
    --name local-print-bridge ^
    --restart unless-stopped ^
    -p 3001:3001 ^
    --env-file .env ^
    linhtranphu/local-print-bridge:latest
```

## Next Steps

### 1. Build Image Mới với WebSocket

Container hiện tại chưa có WebSocket code. Cần build image mới:

```bash
cd local-print-bridge
./build-print-bridge-docker.sh
```

### 2. Restart với Image Mới

```bash
docker stop local-print-bridge
docker rm local-print-bridge
docker run -d \
    --name local-print-bridge \
    --restart unless-stopped \
    -p 3001:3001 \
    --env-file .env \
    linhtranphu/local-print-bridge:latest
```

### 3. Verify WebSocket

```bash
docker logs local-print-bridge | grep WebSocket
# Should see: [WebSocket] ✅ Connected to backend
```

## Files Changed

- `restart_local.sh` - Updated to use port mapping instead of host network

## Troubleshooting

### Still showing "Offline"?

1. **Hard refresh browser**: Ctrl+Shift+R (Windows) or Cmd+Shift+R (Mac)
2. **Check container**: `docker ps | grep print-bridge`
3. **Check health**: `curl http://localhost:3001/health`
4. **Check logs**: `docker logs local-print-bridge`

### Port 3001 already in use?

```bash
# Find process using port 3001
lsof -i :3001

# Kill it
kill -9 <PID>

# Or stop all Print Bridge processes
docker stop local-print-bridge
pkill -f "node.*print-bridge"
```

### Container not starting?

```bash
# Check logs
docker logs local-print-bridge

# Check if .env file exists
ls -la local-print-bridge/.env

# Try running without -d to see errors
docker run --rm -p 3001:3001 --env-file local-print-bridge/.env linhtranphu/local-print-bridge:latest
```

---

**Status**: ✅ Print Bridge accessible on localhost:3001
**Next**: Build new image with WebSocket support
