# Fix Hoàn Chỉnh: Bàn Giao Tiền Chuyển Khoản

## Vấn Đề Ban Đầu

Khi bàn giao tiền chuyển khoản (30,000 VND):
- ❌ Tiền mặt bị trừ: `remaining_cash`: 22,000 → -8,000
- ❌ Tiền CK không bị trừ: `remaining_transfer`: 30,000 → 30,000 (không đổi)
- ❌ `handed_over_cash`: 30,000 (sai, nên là 0)
- ❌ `handed_over_transfer`: 0 (sai, nên là 30,000)

## Nguyên Nhân

### 1. Frontend gửi old format
```javascript
{
  actual_amount: 30000,  // Old format
  status: "CONFIRMED"
}
```

Thay vì:
```javascript
{
  cash_actual_amount: 0,
  transfer_actual_amount: 30000,  // New format
  status: "CONFIRMED"
}
```

### 2. Backend không tự động convert
Backend validation yêu cầu new format khi có `cash_declared_amount` hoặc `transfer_declared_amount`, nhưng không tự động convert từ `actual_amount`.

### 3. Logic update shift sai
Trong `updateCashAmounts`, code đã được sửa để xử lý separate amounts, nhưng do handover không có `cash_actual_amount`/`transfer_actual_amount` đúng, nên vẫn update sai.

## Giải Pháp

### Backend: Auto-convert old format to new format

#### File: `backend/application/services/cash_handover_service.go`

**Bước 1: Auto-convert actual_amount**

```go
// 4. Validate and normalize actual amounts for CONFIRMED status
if req.Status == handover.StatusConfirmed {
    // Auto-convert old format to new format based on declared amounts
    if req.CashActualAmount == 0 && req.TransferActualAmount == 0 && req.ActualAmount > 0 {
        // Old format: distribute actual_amount based on what was declared
        if h.CashDeclaredAmount > 0 && h.TransferDeclaredAmount == 0 {
            // Cash only
            req.CashActualAmount = req.ActualAmount
        } else if h.TransferDeclaredAmount > 0 && h.CashDeclaredAmount == 0 {
            // Transfer only
            req.TransferActualAmount = req.ActualAmount
        } else if h.CashDeclaredAmount > 0 && h.TransferDeclaredAmount > 0 {
            // Both: need separate amounts
            return errors.New("cash_actual_amount and transfer_actual_amount are required")
        } else {
            // Fallback: assume cash for backward compatibility
            req.CashActualAmount = req.ActualAmount
        }
    }
    
    // Validate required amounts
    if h.CashDeclaredAmount > 0 && req.CashActualAmount == 0 {
        return errors.New("cash_actual_amount is required when cash was declared")
    }
    if h.TransferDeclaredAmount > 0 && req.TransferActualAmount == 0 {
        return errors.New("transfer_actual_amount is required when transfer was declared")
    }
}
```

**Bước 2: Calculate discrepancies correctly**

```go
// 5. Calculate discrepancies
var cashDiscrepancy, transferDiscrepancy, totalDiscrepancy float64
if req.Status == handover.StatusConfirmed {
    // Calculate separate discrepancies
    cashDiscrepancy = req.CashActualAmount - h.CashDeclaredAmount
    transferDiscrepancy = req.TransferActualAmount - h.TransferDeclaredAmount
    totalDiscrepancy = cashDiscrepancy + transferDiscrepancy
}
```

**Bước 3: Update handover with correct amounts**

```go
// 6. Update handover
if req.Status == handover.StatusConfirmed {
    // Update separate amounts
    h.CashActualAmount = req.CashActualAmount
    h.TransferActualAmount = req.TransferActualAmount
    h.CashDiscrepancy = cashDiscrepancy
    h.TransferDiscrepancy = transferDiscrepancy
    
    // Update deprecated fields for backward compatibility
    h.ActualAmount = req.CashActualAmount + req.TransferActualAmount
    h.Discrepancy = totalDiscrepancy
    
    h.ReconciledAt = &now
}
```

**Bước 4: Update shift correctly (đã sửa trước đó)**

```go
// updateCashAmounts
func (s *CashHandoverService) updateCashAmounts(ctx context.Context, h *handover.CashHandover) error {
    // Update cash amounts
    if h.CashDeclaredAmount > 0 {
        waiterShift.HandedOverCash += h.CashActualAmount
        waiterShift.RemainingCash -= h.CashDeclaredAmount
        waiterShift.TotalDiscrepancy += h.CashDiscrepancy
    }
    
    // Update transfer amounts
    if h.TransferDeclaredAmount > 0 {
        waiterShift.HandedOverTransfer += h.TransferActualAmount
        waiterShift.RemainingTransfer -= h.TransferDeclaredAmount
    }
    
    // ... rest of the code
}
```

## Kết Quả Sau Khi Fix

### Scenario: Bàn giao 30,000 VND tiền CK

**Handover record:**
```json
{
  "cash_declared_amount": 0,
  "transfer_declared_amount": 30000,
  "cash_actual_amount": 0,        // ✅ Correct
  "transfer_actual_amount": 30000, // ✅ Correct
  "cash_discrepancy": 0,
  "transfer_discrepancy": 0
}
```

**Shift BEFORE:**
```json
{
  "current_cash": 22000,
  "remaining_cash": 22000,
  "handed_over_cash": 0,
  "transfer_revenue": 30000,
  "remaining_transfer": 30000,
  "handed_over_transfer": 0
}
```

**Shift AFTER:**
```json
{
  "current_cash": 22000,           // ✅ Không đổi
  "remaining_cash": 22000,         // ✅ Không đổi
  "handed_over_cash": 0,           // ✅ Không đổi
  "transfer_revenue": 30000,       // ✅ Không đổi
  "remaining_transfer": 0,         // ✅ Trừ 30,000
  "handed_over_transfer": 30000    // ✅ Cộng 30,000
}
```

## Test Cases

### Test 1: Bàn giao chỉ tiền mặt
```
Declared: cash=50000, transfer=0
Actual: actual_amount=50000
Expected: 
  - cash_actual_amount=50000
  - transfer_actual_amount=0
  - remaining_cash -= 50000
  - remaining_transfer không đổi
```

### Test 2: Bàn giao chỉ tiền CK
```
Declared: cash=0, transfer=30000
Actual: actual_amount=30000
Expected:
  - cash_actual_amount=0
  - transfer_actual_amount=30000
  - remaining_cash không đổi
  - remaining_transfer -= 30000
```

### Test 3: Bàn giao cả hai (future)
```
Declared: cash=50000, transfer=30000
Actual: cash_actual_amount=50000, transfer_actual_amount=30000
Expected:
  - remaining_cash -= 50000
  - remaining_transfer -= 30000
```

## Files Changed

1. `backend/application/services/cash_handover_service.go`
   - Auto-convert `actual_amount` to correct type
   - Calculate separate discrepancies
   - Update handover with correct amounts
   - Update shift with correct amounts (đã sửa trước đó)

2. `frontend/src/views/CashierHandoverView.vue`
   - Fix warning display logic (đã sửa trước đó)
   - Support both old and new format

## Backward Compatibility

✅ Old handovers (chỉ có `declared_amount`) vẫn hoạt động
✅ New handovers (có `cash_declared_amount`/`transfer_declared_amount`) hoạt động đúng
✅ Frontend không cần thay đổi (vẫn gửi `actual_amount`)
✅ Backend tự động convert sang đúng format

## Migration cho handovers cũ

Nếu có handovers cũ bị sai, chạy script fix:

```javascript
// Fix handovers that updated wrong amounts
db.cash_handovers.find({
  status: "CONFIRMED",
  transfer_declared_amount: {$gt: 0},
  cash_declared_amount: 0,
  transfer_actual_amount: 0
}).forEach(h => {
  // This handover was transfer-only but didn't set transfer_actual_amount
  db.cash_handovers.updateOne(
    {_id: h._id},
    {$set: {
      transfer_actual_amount: h.actual_amount,
      cash_actual_amount: 0
    }}
  );
  
  // Fix the shift
  db.shifts.updateOne(
    {_id: h.waiter_shift_id},
    {
      $inc: {
        handed_over_cash: -h.actual_amount,
        remaining_cash: h.actual_amount,
        handed_over_transfer: h.actual_amount,
        remaining_transfer: -h.actual_amount
      }
    }
  );
});
```

## Lưu Ý

1. **Restart backend** sau khi deploy code mới
2. **Test kỹ** với các scenarios khác nhau
3. **Backup database** trước khi chạy migration
4. **Monitor logs** để đảm bảo không có lỗi

## Next Steps

1. ✅ Backend auto-convert (DONE)
2. ⏳ Frontend gửi separate amounts (optional, for future)
3. ⏳ UI hiển thị rõ loại tiền đang bàn giao
4. ⏳ Reports phân biệt cash và transfer handovers
