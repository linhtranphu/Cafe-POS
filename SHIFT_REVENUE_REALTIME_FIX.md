# Fix: Shift Revenue Should Use Realtime Values, Not Recalculate

## Vấn đề

Các field như `remaining_cash`, `remaining_transfer`, `current_cash`, `transfer_revenue` đang bị tính lại từ orders mỗi khi GET shift, thay vì sử dụng giá trị realtime đã được update khi collect payment.

## Root Cause

### 1. Realtime Update (ĐÚNG) ✅

Khi collect payment, các field được update đúng:

**File**: `backend/application/services/order_service.go` - Method `CollectPayment`

```go
// Update shift revenue based on payment method
if req.PaymentMethod == order.PaymentCash {
    shift.RemainingCash += req.Amount      // ✅ Update realtime
    shift.CurrentCash += req.Amount         // ✅ Update realtime
} else if req.PaymentMethod == order.PaymentTransfer || 
          req.PaymentMethod == order.PaymentQR {
    shift.TransferRevenue += req.Amount     // ✅ Update realtime
    shift.RemainingTransfer += req.Amount   // ✅ Update realtime
}
shift.TotalRevenue += req.Amount            // ✅ Update realtime
```

### 2. Recalculation (SAI) ❌

Nhưng mỗi khi GET shift, hệ thống gọi `CalculateTransferRevenue` để tính lại:

**File**: `backend/interfaces/http/shift_handler.go`

```go
// GetShift - Line 177-181
if shift.RoleType == order.RoleWaiter {
    _ = h.shiftService.CalculateTransferRevenue(c.Request.Context(), shift.ID)
    // Refetch shift to get updated values
    shift, _ = h.shiftService.GetShift(c.Request.Context(), shift.ID)
}

// GetAllShifts - Line 212-215
for _, shift := range shifts {
    if shift.RoleType == order.RoleWaiter && shift.Status == order.ShiftOpen {
        _ = h.shiftService.CalculateTransferRevenue(c.Request.Context(), shift.ID)
    }
}

// GetCurrentShift - Line 242-246
if shift.RoleType == order.RoleWaiter {
    _ = h.shiftService.CalculateTransferRevenue(c.Request.Context(), shift.ID)
    // Refetch shift to get updated values
    shift, _ = h.shiftService.GetShift(c.Request.Context(), shift.ID)
}
```

**File**: `backend/application/services/shift_service.go` - Method `CalculateTransferRevenue`

```go
// Calculate cash and transfer revenue separately
cashRevenue := 0.0
transferRevenue := 0.0
for _, o := range orders {
    // ❌ Chỉ tính 3 status
    if o.Status == order.StatusPaid || 
       o.Status == order.StatusInProgress || 
       o.Status == order.StatusServed {
        if o.PaymentMethod == order.PaymentCash {
            cashRevenue += o.Total
        } else if o.PaymentMethod == order.PaymentTransfer || 
                  o.PaymentMethod == order.PaymentQR {
            transferRevenue += o.Total
        }
    }
}

// ❌ GHI ĐÈ lên giá trị realtime
shift.TransferRevenue = transferRevenue
shift.RemainingTransfer = transferRevenue - shift.HandedOverTransfer
shift.CurrentCash = shift.StartCash + cashRevenue
shift.RemainingCash = shift.CurrentCash - shift.HandedOverCash
shift.TotalRevenue = cashRevenue + transferRevenue
```

## Flow vấn đề

```
1. Waiter collect payment 50k CASH
   → shift.CurrentCash = 100k + 50k = 150k ✅
   → shift.RemainingCash = 150k ✅
   → Lưu vào DB ✅

2. Waiter xem shift view
   → GET /api/shifts/current
   → Handler gọi CalculateTransferRevenue()
   → Tính lại từ orders (chỉ tính PAID, IN_PROGRESS, SERVED)
   → Nếu order đang ở QUEUED hoặc READY → KHÔNG tính
   → shift.CurrentCash = 100k + 0k = 100k ❌ (GHI ĐÈ)
   → shift.RemainingCash = 100k ❌ (GHI ĐÈ)
   → Lưu vào DB ❌
   → Return về frontend: 100k (SAI!)
```

## Giải pháp

**BỎ** tất cả lời gọi `CalculateTransferRevenue` trong `shift_handler.go` vì:

1. ✅ Các field đã được update realtime khi collect payment
2. ✅ Không cần tính lại từ orders
3. ✅ Tránh ghi đè giá trị đúng
4. ✅ Performance tốt hơn (không cần query orders)

### Implementation

**File**: `backend/interfaces/http/shift_handler.go`

#### Change 1: GetShift method

**Trước**:
```go
// Get shift
shift, err := h.shiftService.GetShift(c.Request.Context(), shiftObjID)
if err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "shift not found"})
    return
}

// Calculate transfer revenue for waiter shifts
if shift.RoleType == order.RoleWaiter {
    _ = h.shiftService.CalculateTransferRevenue(c.Request.Context(), shift.ID)
    // Refetch shift to get updated values
    shift, _ = h.shiftService.GetShift(c.Request.Context(), shift.ID)
}

c.JSON(http.StatusOK, shift)
```

**Sau**:
```go
// Get shift
shift, err := h.shiftService.GetShift(c.Request.Context(), shiftObjID)
if err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "shift not found"})
    return
}

// No need to recalculate - values are updated realtime on payment
c.JSON(http.StatusOK, shift)
```

#### Change 2: GetAllShifts method

**Trước**:
```go
shifts, err := h.shiftService.GetAllShifts(c.Request.Context())
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
}

// Calculate transfer revenue for open waiter shifts
for _, shift := range shifts {
    if shift.RoleType == order.RoleWaiter && shift.Status == order.ShiftOpen {
        _ = h.shiftService.CalculateTransferRevenue(c.Request.Context(), shift.ID)
    }
}

// Refetch to get updated values
shifts, _ = h.shiftService.GetAllShifts(c.Request.Context())

c.JSON(http.StatusOK, shifts)
```

**Sau**:
```go
shifts, err := h.shiftService.GetAllShifts(c.Request.Context())
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
}

// No need to recalculate - values are updated realtime on payment
c.JSON(http.StatusOK, shifts)
```

#### Change 3: GetCurrentShift method

**Trước**:
```go
shift, err := h.shiftService.GetCurrentShift(c.Request.Context(), userObjID)
if err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "no open shift found"})
    return
}

// Calculate transfer revenue for waiter shifts
if shift.RoleType == order.RoleWaiter {
    _ = h.shiftService.CalculateTransferRevenue(c.Request.Context(), shift.ID)
    // Refetch shift to get updated values
    shift, _ = h.shiftService.GetShift(c.Request.Context(), shift.ID)
}

c.JSON(http.StatusOK, shift)
```

**Sau**:
```go
shift, err := h.shiftService.GetCurrentShift(c.Request.Context(), userObjID)
if err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "no open shift found"})
    return
}

// No need to recalculate - values are updated realtime on payment
c.JSON(http.StatusOK, shift)
```

## Về method CalculateTransferRevenue

Method này có thể:
1. **Xóa hoàn toàn** - Nếu không dùng ở đâu khác
2. **Giữ lại** - Để dùng cho migration/fix data cũ nếu cần

Kiểm tra xem còn dùng ở đâu không:
```bash
grep -r "CalculateTransferRevenue" backend/
```

Nếu chỉ dùng trong `shift_handler.go` → Có thể xóa method này luôn.

## Testing

### Test Case 1: Collect payment và xem shift
```
1. Start shift
   → CurrentCash: 100k
   → RemainingCash: 100k

2. Collect payment 50k CASH
   → CurrentCash: 150k ✅
   → RemainingCash: 150k ✅

3. GET /api/shifts/current
   → CurrentCash: 150k ✅ (không bị ghi đè)
   → RemainingCash: 150k ✅ (không bị ghi đè)
```

### Test Case 2: Order ở các status khác nhau
```
1. Start shift: 100k
2. Collect payment 30k (order A - PAID)
3. Send to barista (order A → QUEUED)
4. GET shift
   → CurrentCash: 130k ✅ (không phụ thuộc status)
```

### Test Case 3: Handover
```
1. CurrentCash: 150k
2. Handover 50k
   → HandedOverCash: 50k
   → RemainingCash: 100k
3. GET shift
   → CurrentCash: 150k ✅
   → RemainingCash: 100k ✅
```

## Benefits

1. ✅ **Chính xác**: Hiển thị đúng số tiền realtime
2. ✅ **Đơn giản**: Không cần logic phức tạp tính từ orders
3. ✅ **Performance**: Không cần query orders mỗi lần GET shift
4. ✅ **Maintainable**: Single source of truth (update khi payment)
5. ✅ **Không phụ thuộc status**: Chỉ quan tâm đến payment events

## Migration (Optional)

Nếu data cũ bị sai, có thể chạy script để fix:

```go
// Fix old shifts by recalculating from orders
func FixOldShifts() {
    shifts := getAllShifts()
    for _, shift := range shifts {
        if shift.Status == "CLOSED" {
            // Recalculate for closed shifts only
            CalculateTransferRevenue(shift.ID)
        }
    }
}
```

## Deployment

1. ✅ Phân tích vấn đề
2. ⏳ Remove CalculateTransferRevenue calls from shift_handler.go
3. ⏳ (Optional) Remove CalculateTransferRevenue method if not used elsewhere
4. ⏳ Backend restart
5. ⏳ Test với data thực tế
6. ⏳ Verify shift view hiển thị đúng

## Files to Change

- `backend/interfaces/http/shift_handler.go` - Remove 3 calls to CalculateTransferRevenue
- `backend/application/services/shift_service.go` - (Optional) Remove method if not used
