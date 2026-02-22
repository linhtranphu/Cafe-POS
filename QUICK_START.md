# Quick Start - Local Development

## Khởi động tất cả services

```bash
./restart_local.sh
```

Script này sẽ tự động:
- ✅ Start MongoDB với replica set
- ✅ Start Print Bridge (Docker)
- ✅ Start Backend (Go)
- ✅ Start Frontend (Vue.js)

## Dừng tất cả services

```bash
./stop_local.sh
```

## Truy cập

- **Frontend (Local)**: http://localhost:5173
- **Frontend (LAN)**: http://YOUR_LOCAL_IP:5173 (hiển thị khi chạy restart_local.sh)
- **Backend**: http://localhost:3000
- **Print Bridge**: http://localhost:3001/health
- **MongoDB**: mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin

### Truy cập từ máy khác trong LAN

Frontend đã được cấu hình để cho phép truy cập từ các máy khác trong mạng LAN:

1. Chạy `./restart_local.sh`
2. Xem IP address trong output (ví dụ: `http://192.168.1.100:5173`)
3. Mở browser trên máy khác và truy cập IP đó

**Lưu ý**: 
- Đảm bảo firewall không block port 5173
- Backend vẫn chỉ chạy trên localhost (port 3000)
- Frontend sẽ gọi API qua localhost từ browser

## Xem logs

```bash
# Backend
tail -f backend.log

# Frontend
tail -f frontend.log

# MongoDB
docker logs -f cafe-pos-mongodb

# Print Bridge
docker logs -f local-print-bridge
```

## Test Print Bridge WebSocket

1. Xem logs Print Bridge:
   ```bash
   docker logs -f local-print-bridge
   ```

2. Tạo order trong Frontend

3. Kiểm tra logs:
   ```
   [WebSocket] ✅ Connected to backend
   [WebSocket] 📨 New print job received: {job_id}
   [PrintJobHandler] ✅ Job printed successfully
   ```

## Troubleshooting

### MongoDB không start

```bash
docker-compose -f docker-compose.replica-set.yml down
docker-compose -f docker-compose.replica-set.yml up -d mongodb
```

### Print Bridge không kết nối WebSocket

```bash
# Check backend đang chạy
curl http://localhost:3000/api/login

# Restart Print Bridge
docker restart local-print-bridge
docker logs -f local-print-bridge
```

### Backend không start

```bash
# Check logs
cat backend.log

# Check MongoDB
docker ps | grep mongodb

# Restart
./restart_local.sh
```

## Build Print Bridge Docker Image

```bash
cd local-print-bridge
./build-print-bridge-docker.sh 1.1.0
```

## Cấu trúc Project

```
.
├── backend/              # Go backend
├── frontend/             # Vue.js frontend
├── local-print-bridge/   # Print Bridge service
├── restart_local.sh      # Start all services
├── stop_local.sh         # Stop all services
└── docker-compose.replica-set.yml  # MongoDB config
```

## Default Login

- Username: `admin`
- Password: `admin123`

## Next Steps

- Đọc [WEBSOCKET_ARCHITECTURE.md](./WEBSOCKET_ARCHITECTURE.md) để hiểu kiến trúc
- Đọc [WEBSOCKET_SUMMARY.md](./WEBSOCKET_SUMMARY.md) để hiểu WebSocket implementation
- Đọc [local-print-bridge/README.md](./local-print-bridge/README.md) để hiểu Print Bridge
