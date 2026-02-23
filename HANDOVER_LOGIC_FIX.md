# Fix Bàn Giao Logic - Handover Logic Fix

## Vấn đề (Problem)

Sau khi cashier chấp nhận bàn giao:
1. **Bàn giao tiền mặt**: Tiền mặt của waiter vẫn còn nguyên (không bị trừ)
2. **Bàn giao tiền chuyển khoản**: Tiền CK của waiter bị tăng gấp đôi (thay vì giảm)

After cashier accepts handover:
1. **Cash handover**: Waiter's cash remains unchanged (should decrease)
2. **Transfer handover**: Waiter's transfer DOUBLES (should decrease)

## Nguyên nhân (Root Cause)

Có 2 luồng xử lý bàn giao trong code:

### Luồng MỚI (NEW Flow) - ĐÚNG ✅
1. `CreateHandoverWithTransfer()`: Tạo handover và **TRỪ ngay** `RemainingCash` và `RemainingTransfer`
2. `ConfirmHandoverWithDualAmounts()`: Xác nhận handover
3. `updateDualBalances()`: Chỉ cập nhật `HandedOverCash` và `HandedOverTransfer`, **KHÔNG trừ lại**

### Luồng CŨ (OLD Flow) - SAI ❌
1. `CreateHandover()`: Tạo handover nhưng **KHÔNG trừ** `RemainingCash` và `RemainingTransfer`
2. `ConfirmHandoverWithReconciliation()`: Xác nhận handover
3. `updateCashAmounts()`: Cập nhật `HandedOverCash/Transfer` và **TRỪ** `RemainingCash/Transfer`

**Vấn đề**: Luồng CŨ được cập nhật để hỗ trợ cash/transfer riêng biệt, nhưng logic không nhất quán:
- `CreateHandover()` KHÔNG trừ remaining (giống luồng cũ)
- `updateCashAmounts()` VẪN trừ remaining (giống luồng cũ)
- Nhưng frontend MỚI gọi `CreateHandoverWithTransfer()` (đã trừ remaining rồi)
- Khi confirm, nếu gọi `updateCashAmounts()` sẽ bị **TRỪ 2 LẦN**!

## Giải pháp (Solution)

Thống nhất logic cho cả 2 luồng:
1. **Khi tạo handover**: LUÔN trừ `RemainingCash` và `RemainingTransfer`
2. **Khi confirm handover**: CHỈ cập nhật `HandedOverCash` và `HandedOverTransfer`, KHÔNG trừ lại

### Thay đổi code (Code Changes)

#### 1. Fix `CreateHandover()` - Thêm logic trừ remaining
```go
// backend/application/services/cash_handover_service.go
func (s *CashHandoverService) CreateHandover(...) {
    // ... existing code ...
    
    if err := s.handoverRepo.Create(ctx, h); err != nil {
        return nil, err
    }

    // 7. Update shift balances - reduce remaining amounts
    // This ensures consistency with CreateHandoverWithTransfer
    waiterShift.RemainingCash -= cashAmount
    waiterShift.RemainingTransfer -= transferAmount
    waiterShift.UpdatedAt = time.Now()
    if err := s.shiftRepo.Update(ctx, waiterShiftID, waiterShift); err != nil {
        return nil, err
    }

    return h, nil
}
```

#### 2. Fix `updateCashAmounts()` - Bỏ logic trừ remaining
```go
// backend/application/services/cash_handover_service.go
func (s *CashHandoverService) updateCashAmounts(ctx context.Context, h *handover.CashHandover) error {
    // ... existing code ...
    
    // Update cash amounts
    // NOTE: RemainingCash and RemainingTransfer were already reduced in CreateHandover
    // Here we only update HandedOverCash/Transfer to track what was actually handed over
    if h.CashDeclaredAmount > 0 {
        waiterShift.HandedOverCash += h.CashActualAmount
        // DO NOT reduce RemainingCash again - already done in CreateHandover
        waiterShift.TotalDiscrepancy += h.CashDiscrepancy
    }
    
    // Update transfer amounts
    if h.TransferDeclaredAmount > 0 {
        waiterShift.HandedOverTransfer += h.TransferActualAmount
        // DO NOT reduce RemainingTransfer again - already done in CreateHandover
    }
    
    // ... rest of code ...
}
```

## Luồng xử lý sau khi fix (Fixed Flow)

### Khi tạo bàn giao (Create Handover)
```
Waiter có: RemainingCash = 50,000, RemainingTransfer = 30,000
Waiter bàn giao: 20,000 tiền CK

Sau khi tạo:
- RemainingCash = 50,000 (không đổi)
- RemainingTransfer = 10,000 (30,000 - 20,000) ✅
- HandedOverTransfer = 0 (chưa confirm)
```

### Khi cashier xác nhận (Confirm Handover)
```
Cashier nhận: 20,000 tiền CK (đúng với khai báo)

Sau khi confirm:
- RemainingCash = 50,000 (không đổi)
- RemainingTransfer = 10,000 (không đổi) ✅
- HandedOverTransfer = 20,000 (cập nhật) ✅
```

## Test lại (Re-test)

1. Khởi động lại backend (đã làm)
2. Login waiter, tạo order và thanh toán bằng tiền CK
3. Bàn giao một phần tiền CK
4. Login cashier, chấp nhận bàn giao
5. Kiểm tra:
   - Tiền CK của waiter giảm đúng số tiền bàn giao ✅
   - Tiền mặt của waiter không bị ảnh hưởng ✅
   - Không có số âm ✅

## Files đã sửa (Modified Files)

- `backend/application/services/cash_handover_service.go`
  - `CreateHandover()`: Thêm logic trừ remaining amounts
  - `updateCashAmounts()`: Bỏ logic trừ remaining amounts (đã trừ ở create rồi)
