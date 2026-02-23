# Test Transfer Revenue Feature

## Tóm tắt
Tài liệu này hướng dẫn test tính năng hiển thị tiền chuyển khoản (Transfer Revenue) trong ShiftView.

## Thay đổi đã thực hiện

### Backend Changes
**File: `backend/interfaces/http/shift_handler.go`**

Đã thêm tự động tính toán transfer revenue khi fetch shift data:

```go
// GetCurrentShift - Tự động tính transfer revenue cho waiter shifts
if shift.RoleType == order.RoleWaiter {
    _ = h.shiftService.CalculateTransferRevenue(c.Request.Context(), shift.ID)
    // Refetch shift to get updated values
    shift, _ = h.shiftService.GetShift(c.Request.Context(), shift.ID)
}

// GetShift - Tự động tính transfer revenue cho waiter shifts
if shift.RoleType == order.RoleWaiter {
    _ = h.shiftService.CalculateTransferRevenue(c.Request.Context(), shift.ID)
    // Refetch shift to get updated values
    shift, _ = h.shiftService.GetShift(c.Request.Context(), shift.ID)
}
```

### Frontend (Không thay đổi)
Frontend đã sẵn sàng hiển thị:
- `transfer_revenue` - Tổng tiền CK đã thu
- `remaining_transfer` - Tiền CK còn lại chưa bàn giao
- `handed_over_transfer` - Tiền CK đã bàn giao

## Hướng dẫn Test

### Bước 1: Đảm bảo services đang chạy

```bash
# Kiểm tra MongoDB
docker ps | grep mongo

# Nếu chưa chạy, start MongoDB
docker-compose up -d mongodb

# Kiểm tra Backend
ps aux | grep "go run main.go"

# Nếu chưa chạy, start Backend
cd backend
go run main.go
```

### Bước 2: Tạo user waiter (nếu chưa có)

1. Mở browser: http://localhost:5173
2. Login với admin:
   - Username: `admin`
   - Password: `admin123`
3. Vào menu Users → Create User
4. Tạo user mới:
   - Username: `waiter1`
   - Password: `password123`
   - Full Name: `Waiter 1`
   - Role: `waiter`

### Bước 3: Test bằng Go script

```bash
cd backend
go run cmd/test-transfer-revenue/main.go
```

### Bước 4: Test thủ công trên UI

1. **Logout admin và login waiter1**
   - Username: `waiter1`
   - Password: `password123`

2. **Mở ca làm việc**
   - Vào menu "Ca làm việc"
   - Click "Mở ca"
   - Chọn ca (Sáng/Chiều/Tối)
   - Nhập tiền đầu ca: 100000
   - Click "Mở ca"

3. **Tạo order với payment method Transfer**
   - Vào menu "Order"
   - Tạo order mới
   - Chọn món
   - Chọn payment method: "Chuyển khoản" hoặc "QR"
   - Submit order

4. **Kiểm tra ShiftView**
   - Vào menu "Ca làm việc"
   - Xem phần "Ca đang mở"
   - Kiểm tra 3 ô hiển thị:
     - 💵 Tiền mặt: Hiển thị tiền mặt còn lại
     - 💳 Tiền CK: **Phải hiển thị số tiền từ orders có payment method Transfer/QR**
     - Đã bàn giao: Tổng tiền đã bàn giao

5. **Refresh trang để test lại**
   - Pull to refresh hoặc F5
   - Tiền CK phải vẫn hiển thị đúng

## Kết quả mong đợi

### API Response
Khi gọi `GET /api/waiter/shifts/current`, response phải có:

```json
{
  "id": "...",
  "transfer_revenue": 50000,      // Tổng tiền CK từ orders
  "remaining_transfer": 50000,     // Tiền CK chưa bàn giao
  "handed_over_transfer": 0,       // Tiền CK đã bàn giao
  "current_cash": 150000,
  "remaining_cash": 150000,
  ...
}
```

### UI Display
Trong ShiftView, phần "Ca đang mở" phải hiển thị:

```
💵 Tiền mặt        💳 Tiền CK         Đã bàn giao
   150,000₫           50,000₫            0₫
```

## Troubleshooting

### Vấn đề: Tiền CK vẫn hiển thị 0

**Nguyên nhân có thể:**
1. Backend chưa được restart sau khi sửa code
2. Orders không có payment method Transfer/QR
3. Browser cache

**Giải pháp:**
```bash
# 1. Restart backend
pkill -f "go run main.go"
cd backend
go run main.go

# 2. Clear browser cache
# Trong browser: Ctrl+Shift+R (hard refresh)

# 3. Kiểm tra orders có payment method đúng không
# Vào MongoDB và check:
docker exec -it cafe-pos-mongodb mongosh
use cafe_pos
db.orders.find({payment_method: {$in: ["TRANSFER", "QR"]}}).pretty()
```

### Vấn đề: Backend không compile

```bash
cd backend
go build -o backend main.go
# Nếu có lỗi, check file shift_handler.go
```

## Test Script Output

Khi chạy `go run cmd/test-transfer-revenue/main.go`, output mong đợi:

```
=== Testing Transfer Revenue Calculation ===

1. Logging in as waiter...
✅ Logged in successfully
Token: eyJhbGciOiJIUzI1NiIs...

2. Fetching current shift...

=== Shift Data ===
{
  "id": "65f1234567890abcdef12345",
  "transfer_revenue": 50000,
  "remaining_transfer": 50000,
  "handed_over_transfer": 0,
  "current_cash": 150000,
  "remaining_cash": 150000
}

=== Results ===
✅ transfer_revenue: 50000 VND
✅ remaining_transfer: 50000 VND
✅ handed_over_transfer: 0 VND
✅ current_cash: 150000 VND
✅ remaining_cash: 150000 VND

🎉 Transfer revenue is being calculated correctly!
```

## Kết luận

Sau khi restart backend với code mới, tính năng transfer revenue sẽ hoạt động tự động:
- Backend tự động tính toán transfer revenue từ orders khi fetch shift
- Frontend hiển thị đúng giá trị trong ShiftView
- Không cần thay đổi gì ở frontend
