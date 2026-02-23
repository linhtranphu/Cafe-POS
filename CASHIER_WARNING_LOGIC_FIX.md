# Fix: Cảnh Báo Sai Loại Tiền Trong Cashier View

## Vấn Đề

Khi bàn giao tiền chuyển khoản (30,000 VND), modal xác nhận của cashier hiển thị cảnh báo sai:

```
⚠️ Tiền mặt khai báo nhiều hơn (8.000 ₫)
💵 Tiền mặt còn lại: 22.000 ₫
💳 Tiền CK còn lại: 30.000 ₫
```

Lẽ ra phải hiển thị:
```
✅ Không có cảnh báo
(vì 30,000 CK = 30,000 CK còn lại)
```

## Nguyên Nhân

Trong `shiftCashWarning` computed property, logic xác định `cashDeclared` sai:

```javascript
// ❌ SAI
const cashDeclared = handover.cash_declared_amount || handover.declared_amount || 0
```

Khi handover có:
- `cash_declared_amount = 0`
- `transfer_declared_amount = 30000`
- `declared_amount = 30000` (deprecated field)

Logic trên sẽ fallback về `declared_amount` và coi như là tiền mặt, dẫn đến:
- So sánh 30,000 (CK) với 22,000 (tiền mặt còn lại)
- Hiển thị cảnh báo sai "Tiền mặt khai báo nhiều hơn"

## Giải Pháp

### Bước 1: Xác định format đang dùng

```javascript
// Determine if using new format (separate amounts) or old format
const usingNewFormat = (handover.cash_declared_amount > 0 || handover.transfer_declared_amount > 0)
```

### Bước 2: Lấy cashDeclared đúng

```javascript
// ✅ ĐÚNG
const cashDeclared = usingNewFormat 
  ? (handover.cash_declared_amount || 0)  // New format: chỉ lấy cash_declared_amount
  : (handover.declared_amount || 0)        // Old format: fallback về declared_amount
```

### Bước 3: Chỉ check transfer cho new format

```javascript
// Check transfer mismatch (only for new format)
if (usingNewFormat) {
  const transferDeclared = handover.transfer_declared_amount || 0
  const transferRemaining = shift.remaining_transfer || 0
  if (transferDeclared > 0 && transferDeclared !== transferRemaining) {
    // ... show transfer warning
  }
}
```

### Bước 4: Hiển thị đúng số tiền còn lại

```javascript
return {
  message: warnings.join(' | '),
  cashRemaining: cashDeclared > 0 ? cashRemaining : undefined,
  transferRemaining: (usingNewFormat && handover.transfer_declared_amount > 0) 
    ? transferRemaining 
    : undefined
}
```

## Kết Quả

### Scenario 1: Bàn giao chỉ tiền CK (30,000)

**BEFORE:**
```
⚠️ Tiền mặt khai báo nhiều hơn (8.000 ₫)
💵 Tiền mặt còn lại: 22.000 ₫
💳 Tiền CK còn lại: 30.000 ₫
```

**AFTER:**
```
✅ Không có cảnh báo
(vì 30,000 CK = 30,000 CK còn lại)
```

### Scenario 2: Bàn giao chỉ tiền mặt (50,000)

**Shift có:**
- `remaining_cash = 22,000`
- `remaining_transfer = 30,000`

**Hiển thị:**
```
⚠️ Tiền mặt khai báo nhiều hơn (28.000 ₫)
💵 Tiền mặt còn lại: 22.000 ₫
```

### Scenario 3: Bàn giao cả hai

**Handover:**
- `cash_declared_amount = 20,000`
- `transfer_declared_amount = 30,000`

**Shift có:**
- `remaining_cash = 22,000`
- `remaining_transfer = 30,000`

**Hiển thị:**
```
⚠️ Tiền mặt khai báo ít hơn (2.000 ₫)
💵 Tiền mặt còn lại: 22.000 ₫
💳 Tiền CK còn lại: 30.000 ₫
```

### Scenario 4: Old format handover

**Handover:**
- `declared_amount = 50,000`
- `cash_declared_amount = 0`
- `transfer_declared_amount = 0`

**Shift có:**
- `remaining_cash = 22,000`

**Hiển thị:**
```
⚠️ Tiền mặt khai báo nhiều hơn (28.000 ₫)
💵 Tiền mặt còn lại: 22.000 ₫
```

## Logic Flow

```
┌─────────────────────────────────────┐
│ Check handover format               │
└──────────────┬──────────────────────┘
               │
               ├─ cash_declared_amount > 0 OR transfer_declared_amount > 0?
               │
       ┌───────┴───────┐
       │               │
      YES             NO
       │               │
  New Format      Old Format
       │               │
       ├─ cashDeclared = cash_declared_amount
       │  transferDeclared = transfer_declared_amount
       │  
       │  Check cash if cashDeclared > 0
       │  Check transfer if transferDeclared > 0
       │
       └─ cashDeclared = declared_amount
          
          Check cash only
```

## Files Changed

- `frontend/src/views/CashierHandoverView.vue`
  - Fixed `shiftCashWarning` computed property
  - Added `usingNewFormat` detection
  - Separated cash and transfer logic
  - Only show relevant remaining amounts

## Testing

### Test 1: Transfer-only handover
```bash
# Create handover with transfer only
curl -X POST http://localhost:8080/api/handovers \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "cash_amount": 0,
    "transfer_amount": 30000,
    "handover_type": "PARTIAL"
  }'

# Expected: No warning (if shift has 30,000 transfer remaining)
```

### Test 2: Cash-only handover
```bash
# Create handover with cash only
curl -X POST http://localhost:8080/api/handovers \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "cash_amount": 50000,
    "transfer_amount": 0,
    "handover_type": "PARTIAL"
  }'

# Expected: Warning if shift has less than 50,000 cash remaining
```

### Test 3: Both cash and transfer
```bash
# Create handover with both
curl -X POST http://localhost:8080/api/handovers \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "cash_amount": 20000,
    "transfer_amount": 30000,
    "handover_type": "PARTIAL"
  }'

# Expected: Separate warnings for each type if mismatch
```

## Backward Compatibility

✅ Old handovers (chỉ có `declared_amount`) vẫn hoạt động
✅ Hiển thị cảnh báo tiền mặt cho old format
✅ New handovers hiển thị cảnh báo đúng loại tiền

## Benefits

1. ✅ **Chính xác**: Cảnh báo đúng loại tiền đang bàn giao
2. ✅ **Rõ ràng**: Không còn nhầm lẫn giữa tiền mặt và CK
3. ✅ **Tránh lỗi**: Cashier không bị nhầm khi xác nhận
4. ✅ **Backward compatible**: Vẫn hoạt động với handovers cũ

## Next Steps

1. ✅ Fix warning logic (DONE)
2. ⏳ Test với real data
3. ⏳ User feedback
4. ⏳ Monitor for any edge cases
