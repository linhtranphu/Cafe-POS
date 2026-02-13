# Hướng Dẫn Chạy Đầy Đủ - Test Chi Phí & Lợi Nhuận

## Bước 1: Start Backend

### Option 1: Chạy trực tiếp
```bash
cd backend
go run main.go
```

### Option 2: Sử dụng Docker
```bash
docker-compose up -d
```

Backend sẽ chạy tại: `http://localhost:8080`

## Bước 2: Kiểm Tra Backend

```bash
curl http://localhost:8080/api/menu
```

Nếu thấy response JSON → Backend đã sẵn sàng ✅

## Bước 3: Seed Ingredients (Nếu Chưa Có)

```bash
go run backend/cmd/seed/main.go
```

Hoặc:

```bash
cd backend
go run cmd/seed/main.go
```

Điều này tạo dữ liệu nguyên liệu với giá và wastage.

## Bước 4: Seed Menu Items với Variants

```bash
./seed-menu-variants-auto.sh
```

Hoặc:

```bash
cd backend
echo "y" | go run cmd/seed-menu-variants/main.go
```

Kết quả: 10 món (6 multi-size, 4 single-size) ✅

## Bước 5: Tính Chi Phí

```bash
./calculate-costs-simple.sh
```

Script này sẽ:
- Lấy tất cả menu items
- Gọi API calculate-cost cho từng món
- Hiển thị kết quả

## Bước 6: Xem Kết Quả

### Option 1: Qua API

```bash
# Xem tất cả món với chi phí
curl http://localhost:8080/api/menu

# Xem chi tiết chi phí một món
curl http://localhost:8080/api/menu/{MENU_ID}/cost-breakdown

# Xem phân tích lợi nhuận
curl http://localhost:8080/api/menu/{MENU_ID}/profit-analysis
```

### Option 2: Qua Frontend

1. Start frontend:
```bash
cd frontend
npm run dev
```

2. Truy cập:
```
http://localhost:5173/cost-analysis
```

## Troubleshooting

### Lỗi: "Failed to connect to backend"

**Giải pháp:**
```bash
# Kiểm tra backend có chạy không
ps aux | grep "go run main.go"

# Hoặc kiểm tra port
lsof -i :8080

# Start backend nếu chưa chạy
cd backend && go run main.go
```

### Lỗi: "No menu items found"

**Giải pháp:**
```bash
# Seed lại menu
./seed-menu-variants-auto.sh
```

### Lỗi: "Cost status INCOMPLETE"

**Nguyên nhân:** Thiếu dữ liệu nguyên liệu

**Giải pháp:**
```bash
# Seed ingredients
cd backend
go run cmd/seed/main.go
```

### Lỗi: "jq: command not found"

**Giải pháp:** Sử dụng script mới không cần jq
```bash
./calculate-costs-simple.sh
```

## Kiểm Tra Nhanh

### 1. Backend đang chạy?
```bash
curl http://localhost:8080/api/menu
```

### 2. Có menu items?
```bash
curl http://localhost:8080/api/menu | grep -c "id"
```

### 3. Có ingredients?
```bash
curl http://localhost:8080/api/ingredients | grep -c "id"
```

## Workflow Hoàn Chỉnh

```
1. Start Backend
   ↓
2. Seed Ingredients (nếu chưa có)
   ↓
3. Seed Menu Items
   ↓
4. Calculate Costs
   ↓
5. View Analysis (API hoặc Frontend)
```

## Scripts Có Sẵn

- ✅ `seed-menu-variants-auto.sh` - Seed menu tự động
- ✅ `calculate-costs-simple.sh` - Tính chi phí (không cần jq)
- ⚠️ `calculate-all-menu-costs.sh` - Cần cài jq
- ⚠️ `view-menu-cost-analysis.sh` - Cần cài jq

## Cài Đặt jq (Optional)

Nếu muốn dùng scripts có format đẹp hơn:

### macOS:
```bash
brew install jq
```

### Ubuntu/Debian:
```bash
sudo apt-get install jq
```

Sau khi cài jq, có thể dùng:
```bash
./calculate-all-menu-costs.sh
./view-menu-cost-analysis.sh
```

## Kết Quả Mong Đợi

Sau khi hoàn thành, bạn sẽ thấy:

```
📊 Cà phê sữa đá
  Size M:  Cost: 13,800đ, Margin: 44.8%
  Size L:  Cost: 20,700đ, Margin: 31.0%
  Size XL: Cost: 27,600đ, Margin: 21.1%
```

Chúc bạn test thành công! 🎉
