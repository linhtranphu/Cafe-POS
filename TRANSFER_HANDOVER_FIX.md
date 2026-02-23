# Fix: Bàn giao tiền chuyển khoản bị trừ sai vào tiền mặt

## Vấn đề

Khi bàn giao tiền chuyển khoản (transfer), hệ thống đang:
1. **Hiển thị cảnh báo sai**: So sánh với `remaining_cash` thay vì `remaining_transfer`
2. **Trừ tiền sai**: Trừ từ `remaining_cash` thay vì `remaining_transfer`
3. **Đóng ca sai**: Do tiền mặt bị âm (vì trừ nhầm số tiền transfer lớn hơn tiền mặt hiện có)

## Nguyên nhân

### Backend
- `CreateHandoverRequest` chỉ có `declared_amount` chung, không phân biệt cash và transfer
- Logic validation so sánh `declared_amount` với `remaining_cash` cho TẤT CẢ handover
- Logic `updateCashAmounts` trừ từ `remaining_cash` cho TẤT CẢ handover

### Frontend
- `CashierHandoverView` hiển thị warning bằng cách so sánh `declared_amount` với `remaining_cash`
- Không kiểm tra loại handover (cash/transfer/both)

## Giải pháp

### 1. Backend Domain Updates

#### `backend/domain/handover/cash_handover.go`

Thêm các trường mới vào `CreateHandoverRequest`:

```go
type CreateHandoverRequest struct {
    // Separate amounts for cash and transfer
    CashAmount     float64      `json:"cash_amount" binding:"gte=0"`
    TransferAmount float64      `json:"transfer_amount" binding:"gte=0"`
    
    // DEPRECATED: Keep for backward compatibility
    DeclaredAmount float64      `json:"declared_amount" binding:"gte=0"`
    
    HandoverType   HandoverType `json:"handover_type" binding:"required"`
    WaiterNote     string       `json:"waiter_note"`
}
```

Thêm helper methods:
- `GetCashAmount()` - Lấy cash amount với backward compatibility
- `GetTransferAmount()` - Lấy transfer amount
- `GetTotalAmount()` - Tổng cash + transfer

Tương tự cho `ConfirmHandoverRequest`:

```go
type ConfirmHandoverRequest struct {
    // Separate actual amounts
    CashActualAmount     float64 `json:"cash_actual_amount" binding:"gte=0"`
    TransferActualAmount float64 `json:"transfer_actual_amount" binding:"gte=0"`
    
    // DEPRECATED
    ActualAmount         float64 `json:"actual_amount" binding:"gte=0"`
    
    Status                    HandoverStatus     `json:"status" binding:"required"`
    CashierNote               string             `json:"cashier_note"`
    DiscrepancyReason         string             `json:"discrepancy_reason"`
    DiscrepancyResponsibility ResponsibilityType `json:"discrepancy_responsibility"`
}
```

### 2. Backend Service Updates

#### `backend/application/services/cash_handover_service.go`

**Sửa validation trong `CreateHandover`:**

```go
// Validate cash amount doesn't exceed remaining cash
if cashAmount > waiterShift.RemainingCash {
    return nil, errors.New("cash amount exceeds remaining cash in shift")
}

// Validate transfer amount doesn't exceed remaining transfer
if transferAmount > waiterShift.RemainingTransfer {
    return nil, errors.New("transfer amount exceeds remaining transfer in shift")
}
```

**Sửa logic tạo handover:**

```go
h := &handover.CashHandover{
    // New separate amounts
    CashDeclaredAmount:     cashAmount,
    TransferDeclaredAmount: transferAmount,
    CashActualAmount:       0,
    TransferActualAmount:   0,
    CashDiscrepancy:        0,
    TransferDiscrepancy:    0,
    
    // DEPRECATED: Keep for backward compatibility
    DeclaredAmount: req.GetTotalAmount(),
    // ...
}
```

**Sửa logic confirm:**

```go
// Calculate separate discrepancies
cashDiscrepancy = req.CashActualAmount - h.CashDeclaredAmount
transferDiscrepancy = req.TransferActualAmount - h.TransferDeclaredAmount
totalDiscrepancy = cashDiscrepancy + transferDiscrepancy
```

**Sửa `updateCashAmounts`:**

```go
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
```

### 3. Frontend Updates

#### `frontend/src/views/CashierHandoverView.vue`

**Sửa `hasShiftCashMismatch`:**

```javascript
const hasShiftCashMismatch = (handover) => {
  const shift = shiftsMap.value[handover.waiter_shift_id]
  if (!shift) return false
  
  // Check cash mismatch
  const cashDeclared = handover.cash_declared_amount || handover.declared_amount || 0
  const cashRemaining = shift.remaining_cash || 0
  const hasCashMismatch = cashDeclared > 0 && cashDeclared !== cashRemaining
  
  // Check transfer mismatch
  const transferDeclared = handover.transfer_declared_amount || 0
  const transferRemaining = shift.remaining_transfer || 0
  const hasTransferMismatch = transferDeclared > 0 && transferDeclared !== transferRemaining
  
  return hasCashMismatch || hasTransferMismatch
}
```

**Sửa hiển thị warning:**

```vue
<!-- Cash warning -->
<div v-if="handover.cash_declared_amount > 0">
  <p>💵 {{ handover.cash_declared_amount > shift.remaining_cash 
    ? 'Tiền mặt: Khai báo nhiều hơn' 
    : 'Tiền mặt: Khai báo ít hơn' }}</p>
  <p>Còn lại: {{ formatPrice(shift.remaining_cash) }}</p>
</div>

<!-- Transfer warning -->
<div v-if="handover.transfer_declared_amount > 0">
  <p>💳 {{ handover.transfer_declared_amount > shift.remaining_transfer 
    ? 'Tiền CK: Khai báo nhiều hơn' 
    : 'Tiền CK: Khai báo ít hơn' }}</p>
  <p>Còn lại: {{ formatPrice(shift.remaining_transfer) }}</p>
</div>
```

**Sửa `shiftCashWarning` computed:**

```javascript
const shiftCashWarning = computed(() => {
  if (!selectedHandover.value) return null
  
  const handover = selectedHandover.value
  const shift = shiftsMap.value[handover.waiter_shift_id]
  if (!shift) return null
  
  const warnings = []
  
  // Check cash mismatch
  const cashDeclared = handover.cash_declared_amount || handover.declared_amount || 0
  if (cashDeclared > 0 && cashDeclared !== shift.remaining_cash) {
    warnings.push(`Tiền mặt khai báo ${cashDeclared > shift.remaining_cash ? 'nhiều' : 'ít'} hơn`)
  }
  
  // Check transfer mismatch
  const transferDeclared = handover.transfer_declared_amount || 0
  if (transferDeclared > 0 && transferDeclared !== shift.remaining_transfer) {
    warnings.push(`Tiền CK khai báo ${transferDeclared > shift.remaining_transfer ? 'nhiều' : 'ít'} hơn`)
  }
  
  if (warnings.length > 0) {
    return {
      message: warnings.join(' | '),
      cashRemaining: cashDeclared > 0 ? shift.remaining_cash : undefined,
      transferRemaining: transferDeclared > 0 ? shift.remaining_transfer : undefined
    }
  }
  
  return null
})
```

## Backward Compatibility

Tất cả thay đổi đều duy trì backward compatibility:

1. **Deprecated fields vẫn được giữ**: `declared_amount`, `actual_amount`, `discrepancy`
2. **Helper methods fallback**: Nếu không có `cash_amount`/`transfer_amount`, sử dụng `declared_amount`
3. **Frontend fallback**: Hiển thị `declared_amount` nếu không có separate amounts
4. **Old handovers vẫn hoạt động**: Logic kiểm tra cả old và new format

## Testing

### Test Case 1: Bàn giao chỉ tiền mặt
```json
{
  "cash_amount": 50000,
  "transfer_amount": 0,
  "handover_type": "PARTIAL"
}
```
- ✅ Validate với `remaining_cash`
- ✅ Trừ từ `remaining_cash`
- ✅ Hiển thị warning đúng

### Test Case 2: Bàn giao chỉ tiền chuyển khoản
```json
{
  "cash_amount": 0,
  "transfer_amount": 75000,
  "handover_type": "PARTIAL"
}
```
- ✅ Validate với `remaining_transfer`
- ✅ Trừ từ `remaining_transfer`
- ✅ Hiển thị warning đúng
- ✅ KHÔNG trừ từ `remaining_cash`

### Test Case 3: Bàn giao cả hai
```json
{
  "cash_amount": 50000,
  "transfer_amount": 75000,
  "handover_type": "PARTIAL"
}
```
- ✅ Validate cả hai
- ✅ Trừ đúng từ mỗi loại
- ✅ Hiển thị warning cho cả hai

### Test Case 4: Backward compatibility
```json
{
  "declared_amount": 50000,
  "handover_type": "PARTIAL"
}
```
- ✅ Vẫn hoạt động như cũ
- ✅ Được xử lý như cash-only handover

## Files Changed

### Backend
- `backend/domain/handover/cash_handover.go` - Thêm separate amounts
- `backend/application/services/cash_handover_service.go` - Sửa validation và update logic

### Frontend
- `frontend/src/views/CashierHandoverView.vue` - Sửa warning display logic

## Impact

- ✅ Bàn giao tiền chuyển khoản không còn trừ nhầm vào tiền mặt
- ✅ Cảnh báo hiển thị đúng loại tiền
- ✅ Ca shift không bị đóng sai do tiền âm
- ✅ Backward compatible với handovers cũ
- ✅ Hỗ trợ đầy đủ cho bank transfer handover feature

## Next Steps

1. Test với các scenarios thực tế
2. Cập nhật frontend confirm form để nhập separate actual amounts
3. Thêm UI để hiển thị separate amounts trong history
4. Cập nhật reports để phân biệt cash và transfer revenue
