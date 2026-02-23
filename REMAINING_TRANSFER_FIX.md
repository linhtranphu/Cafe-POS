# Fix RemainingTransfer - Khởi tạo khi Start Shift

## VẤN ĐỀ

**Hiện tượng**: Sau khi bàn giao tiền CK, frontend vẫn hiển thị số tiền CK không giảm.

**Nguyên nhân**: 
1. Backend `CollectPayment()` đã cập nhật `remaining_transfer` khi thu tiền ✅
2. Backend `CreateHandover()` đã trừ `remaining_transfer` khi bàn giao ✅
3. **NHƯNG** `StartShift()` KHÔNG khởi tạo `remaining_transfer = 0` ❌

**Kết quả**:
- Shift mới có `remaining_transfer = null` hoặc `undefined`
- Khi thu tiền CK: `null + 30000 = NaN` hoặc không cập nhật
- Frontend fallback sang `transfer_revenue` (không bị trừ khi bàn giao)

## GIẢI PHÁP

Khởi tạo TẤT CẢ fields liên quan khi start shift:

### Code cũ (THIẾU)
```go
shift := &order.Shift{
    Type:          req.Type,
    Status:        order.ShiftOpen,
    RoleType:      roleType,
    UserID:        userOID,
    UserName:      userName,
    StartCash:     req.StartCash,
    CurrentCash:   req.StartCash,
    RemainingCash: req.StartCash,
    StartedAt:     time.Now(),
}
```

### Code mới (ĐẦY ĐỦ)
```go
shift := &order.Shift{
    Type:               req.Type,
    Status:             order.ShiftOpen,
    RoleType:           roleType,
    UserID:             userOID,
    UserName:           userName,
    StartCash:          req.StartCash,
    CurrentCash:        req.StartCash,       // ✅ Tiền mặt hiện tại
    RemainingCash:      req.StartCash,       // ✅ Tiền mặt còn lại
    TransferRevenue:    0,                   // ✅ Doanh thu CK
    RemainingTransfer:  0,                   // ✅ Tiền CK còn lại
    HandedOverCash:     0,                   // ✅ Tiền mặt đã bàn giao
    HandedOverTransfer: 0,                   // ✅ Tiền CK đã bàn giao
    StartedAt:          time.Now(),
}
```

## LUỒNG XỬ LÝ SAU KHI FIX

### 1. Start Shift
```
Waiter mở ca với start_cash = 10,000

Shift được tạo:
- start_cash = 10,000
- current_cash = 10,000
- remaining_cash = 10,000
- transfer_revenue = 0          ✅ Khởi tạo
- remaining_transfer = 0        ✅ Khởi tạo
- handed_over_cash = 0          ✅ Khởi tạo
- handed_over_transfer = 0      ✅ Khởi tạo
```

### 2. Thu tiền CK
```
Order thanh toán 30,000 bằng CK

CollectPayment() cập nhật:
- transfer_revenue = 0 + 30,000 = 30,000      ✅
- remaining_transfer = 0 + 30,000 = 30,000    ✅
```

### 3. Bàn giao tiền CK
```
Waiter bàn giao 20,000 tiền CK

CreateHandover() cập nhật:
- remaining_transfer = 30,000 - 20,000 = 10,000  ✅

ConfirmHandover() cập nhật:
- handed_over_transfer = 0 + 20,000 = 20,000     ✅
```

### 4. Frontend hiển thị
```vue
📊 Doanh thu
- 💵 Tiền mặt: 10,000 (current_cash)
- 💳 Tiền CK: 30,000 (transfer_revenue)

💰 Còn lại (chưa bàn giao)
- 💵 Tiền mặt: 10,000 (remaining_cash)
- 💳 Tiền CK: 10,000 (remaining_transfer) ✅ ĐÚNG!

✅ Đã bàn giao
- 💵 Tiền mặt: 0 (handed_over_cash)
- 💳 Tiền CK: 20,000 (handed_over_transfer) ✅ ĐÚNG!
```

## FIELDS TRONG SHIFT

### Doanh thu (Revenue)
- `current_cash` - Tổng tiền mặt đã thu (bao gồm start_cash)
- `transfer_revenue` - Tổng tiền CK đã thu

### Còn lại (Remaining)
- `remaining_cash` - Tiền mặt chưa bàn giao
- `remaining_transfer` - Tiền CK chưa bàn giao

### Đã bàn giao (Handed Over)
- `handed_over_cash` - Tiền mặt đã bàn giao
- `handed_over_transfer` - Tiền CK đã bàn giao

### Công thức
```
current_cash = start_cash + (tiền mặt từ orders)
remaining_cash = current_cash - handed_over_cash

transfer_revenue = (tiền CK từ orders)
remaining_transfer = transfer_revenue - handed_over_transfer
```

## LƯU Ý CHO SHIFTS CŨ

Các shift đã tồn tại trước khi fix này sẽ có:
- `remaining_transfer = null` hoặc `undefined`
- Frontend sẽ hiển thị "Tiền CK còn lại: 0 ₫"

**Giải pháp**: Đóng tất cả shifts cũ và mở shift mới để có đầy đủ fields.

## FILES THAY ĐỔI

- `backend/application/services/shift_service.go` - Thêm khởi tạo fields trong `StartShift()`
- `backend/application/services/order_service.go` - Đã có logic cập nhật `remaining_transfer` ✅
- `backend/application/services/cash_handover_service.go` - Đã có logic trừ `remaining_transfer` ✅
- `frontend/src/views/ShiftView.vue` - Đã tách bạch hiển thị ✅

## KẾT QUẢ

✅ Shift mới có đầy đủ fields được khởi tạo
✅ Thu tiền CK cập nhật `remaining_transfer` đúng
✅ Bàn giao tiền CK trừ `remaining_transfer` đúng
✅ Frontend hiển thị chính xác số tiền còn lại
✅ Không còn fallback logic gây nhầm lẫn
