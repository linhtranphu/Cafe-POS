# 🚀 START HERE - Quick Setup

## Bước 1: Start Backend (Terminal 1)

Mở terminal mới và chạy:

```bash
cd backend
go run main.go
```

**Giữ terminal này chạy!** Backend sẽ chạy tại `http://localhost:3000`

## Bước 2: Run Setup Script (Terminal 2)

Mở terminal mới (khác với terminal backend) và chạy:

```bash
./setup-and-calculate-costs.sh
```

Script này sẽ tự động:
- ✅ Kiểm tra backend đang chạy
- ✅ Seed ingredients (nếu chưa có)
- ✅ Seed menu items (10 món)
- ✅ Tính chi phí cho tất cả món

## Kết Quả

Sau khi hoàn thành, bạn sẽ thấy:

```
🎉 SETUP COMPLETE!

✅ Backend is running
✅ Ingredients seeded
✅ Menu items seeded (10 items)
✅ Costs calculated
```

## Xem Kết Quả

### Option 1: Qua API

```bash
# Lấy danh sách món
curl http://localhost:3000/api/menu

# Lấy ID món đầu tiên
MENU_ID=$(curl -s http://localhost:3000/api/menu | grep -o '"id":"[^"]*"' | head -1 | sed 's/"id":"//g' | sed 's/"//g')

# Xem chi tiết chi phí
curl http://localhost:3000/api/menu/$MENU_ID/cost-breakdown

# Xem phân tích lợi nhuận
curl http://localhost:3000/api/menu/$MENU_ID/profit-analysis
```

### Option 2: Qua Frontend

1. Start frontend (Terminal 3):
```bash
cd frontend
npm run dev
```

2. Truy cập:
```
http://localhost:5173/cost-analysis
```

## Troubleshooting

### Backend không start được?

Kiểm tra:
```bash
# MongoDB có chạy không?
ps aux | grep mongod

# Port 3000 có bị chiếm không?
lsof -i :3000

# Kiểm tra .env file
cat backend/.env
```

### Script báo lỗi?

Chạy từng bước thủ công:

```bash
# 1. Seed ingredients
cd backend
go run cmd/seed/main.go

# 2. Seed menu
cd ..
./seed-menu-variants-auto.sh

# 3. Calculate costs
./calculate-costs-simple.sh
```

## Tóm Tắt

```
Terminal 1: cd backend && go run main.go
Terminal 2: ./setup-and-calculate-costs.sh
Terminal 3: cd frontend && npm run dev (optional)
```

Chúc bạn test thành công! 🎉
