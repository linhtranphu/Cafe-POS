# Debug 500 Internal Server Error

## Lỗi

```
POST http://localhost:5173/api/cashier-shifts/699c82b…/close-with-fund-handover 500 (Internal Server Error)
```

## Cách Debug

### 1. Kiểm Tra Backend Logs

Xem terminal đang chạy backend (`go run main.go`), sẽ thấy error message chi tiết.

Các lỗi có thể gặp:

#### A. MongoDB Connection Error
```
Error: no reachable servers
```
**Fix**: Khởi động MongoDB

#### B. Validation Error
```
Error: variance requires documentation: reason and notes are required
```
**Fix**: Đảm bảo gửi variance_reason và variance_notes khi có chênh lệch

#### C. Transaction Error
```
Error: Transaction numbers are only allowed on a replica set member or mongos
```
**Fix**: MongoDB phải chạy ở chế độ replica set

#### D. Field Mapping Error
```
Error: cannot unmarshal...
```
**Fix**: Kiểm tra JSON field names

### 2. Test API Trực Tiếp

```bash
# Lấy token
TOKEN="your_jwt_token"

# Lấy shift ID
SHIFT_ID="699c82b..."

# Test với curl
curl -X POST "http://localhost:3000/api/cashier-shifts/$SHIFT_ID/close-with-fund-handover" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "actual_cash": 2000000,
    "variance_reason": null,
    "variance_notes": null
  }' \
  -v
```

### 3. Kiểm Tra Request Payload

Mở DevTools > Network > Click vào request failed > Headers tab > Request Payload

Đảm bảo:
```json
{
  "actual_cash": 2000000,  // number, not string
  "variance_reason": null,  // hoặc "COUNTING_ERROR"
  "variance_notes": null    // hoặc "Ghi chú..."
}
```

### 4. Các Lỗi Thường Gặp

#### Lỗi 1: MongoDB Không Chạy Replica Set

**Triệu chứng**:
```
Transaction numbers are only allowed on a replica set member
```

**Fix**:
```bash
# Dừng MongoDB
mongod --shutdown

# Khởi động với replica set
mongod --replSet rs0

# Trong mongosh
rs.initiate()
```

#### Lỗi 2: Thiếu Variance Documentation

**Triệu chứng**:
```
variance requires documentation: reason and notes are required
```

**Fix**: Khi actual_cash ≠ expected_cash, phải gửi:
```json
{
  "actual_cash": 1995000,
  "variance_reason": "COUNTING_ERROR",
  "variance_notes": "Đếm nhầm tờ 50k thành 100k"
}
```

#### Lỗi 3: Waiter Shifts Chưa Đóng

**Triệu chứng**:
```
cannot close cashier shift: X waiter shift(s) still open
```

**Fix**: Đóng tất cả waiter shifts trước

#### Lỗi 4: Shift Không Ở Trạng Thái OPEN

**Triệu chứng**:
```
shift must be in OPEN status, current: CLOSING
```

**Fix**: Shift phải ở trạng thái OPEN, không thể đóng shift đang CLOSING

### 5. Enable Debug Logging

Thêm vào backend để xem chi tiết:

```go
// backend/interfaces/http/cashier_shift_closure_handler.go
func (h *CashierShiftClosureHandler) CloseShiftWithFundHandover(c *gin.Context) {
    // ... existing code ...
    
    // Add logging
    log.Printf("Request: %+v", req)
    log.Printf("Shift ID: %s", shiftID)
    log.Printf("User ID: %s", userID)
    
    // ... rest of code ...
}
```

### 6. Kiểm Tra Database

```javascript
// MongoDB shell
use cafe_pos

// Kiểm tra shift
db.cashier_shifts.findOne({ _id: ObjectId("699c82b...") })

// Kiểm tra status
db.cashier_shifts.findOne(
  { _id: ObjectId("699c82b...") },
  { status: 1, received_cash: 1, received_transfer: 1, starting_float: 1 }
)

// Kiểm tra waiter shifts
db.shifts.find({ status: "OPEN", role_type: "WAITER" }).count()
```

## Quick Fix

### Nếu Lỗi MongoDB Transaction

```bash
# 1. Dừng backend (Ctrl+C)

# 2. Khởi động MongoDB với replica set
mongod --replSet rs0

# 3. Trong mongosh
mongosh
rs.initiate()

# 4. Khởi động lại backend
cd backend
go run main.go
```

### Nếu Lỗi Validation

Kiểm tra frontend gửi đúng format:

```javascript
// frontend/src/views/CashierShiftClosureV2.vue
const payload = {
  actual_cash: actualCash,  // Phải là number
  variance_reason: hasVariance ? varianceReason : null,
  variance_notes: hasVariance ? varianceNotes : null
}
```

## Restart Backend

Sau khi fix, restart backend:

```bash
cd backend
# Ctrl+C để dừng
go run main.go
```

## Test Lại

1. Refresh trình duyệt
2. Thử đóng ca lại
3. Kiểm tra backend logs
4. Kiểm tra response trong DevTools

---

**Lưu ý**: Hãy xem backend terminal để biết lỗi cụ thể!
