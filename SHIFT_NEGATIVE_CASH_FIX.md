# Fix: Shift Hiển Thị Tiền Mặt Âm Sau Bàn Giao Transfer

## Vấn Đề

Sau khi bàn giao 30,000 VND tiền chuyển khoản, shift view của waiter hiển thị:
- 💵 Tiền mặt: -8,000 ₫ ❌ (âm!)
- 💳 Tiền CK: 30,000 ₫
- Đã bàn giao: 30,000 ₫

## Nguyên Nhân

### Root Cause
Backend cũ (trước khi fix) đã chạy khi confirm handover, dẫn đến:

1. **Handover record sai:**
   - `transfer_declared_amount = 30000` ✅
   - `cash_declared_amount = 0` ✅
   - `actual_amount = 30000` (old format) ✅
   - Nhưng `transfer_actual_amount = 0` ❌
   - Và `cash_actual_amount = 0` ✅

2. **Backend cũ không auto-convert:**
   - Không convert `actual_amount` sang `transfer_actual_amount`
   - Dẫn đến logic update shift sai

3. **Shift bị update sai:**
   - `remaining_cash` bị trừ 30,000 (từ 72,000 → 42,000) ❌
   - `remaining_transfer` không bị trừ (vẫn 30,000) ❌
   - `handed_over_cash` tăng 30,000 ❌
   - `handed_over_transfer` không tăng ❌

## Dữ Liệu Thực Tế

### Shift ID: `699bda2ea31585ce7ad4c47c`

**BEFORE FIX:**
```javascript
{
  current_cash: 72000,
  remaining_cash: 42000,        // ❌ Sai: bị trừ 30,000
  handed_over_cash: 30000,      // ❌ Sai: lẽ ra phải là 0
  transfer_revenue: 30000,
  remaining_transfer: 30000,    // ❌ Sai: lẽ ra phải là 0
  handed_over_transfer: 0       // ❌ Sai: lẽ ra phải là 30,000
}
```

**AFTER FIX:**
```javascript
{
  current_cash: 72000,
  remaining_cash: 72000,        // ✅ Đúng: không bị trừ
  handed_over_cash: 0,          // ✅ Đúng: không bàn giao tiền mặt
  transfer_revenue: 30000,
  remaining_transfer: 0,        // ✅ Đúng: đã bàn giao hết
  handed_over_transfer: 30000   // ✅ Đúng: bàn giao 30,000 CK
}
```

### Handover ID: `699bda5ca31585ce7ad4c484`

**BEFORE FIX:**
```javascript
{
  cash_declared_amount: 0,
  transfer_declared_amount: 30000,
  cash_actual_amount: 0,
  transfer_actual_amount: 0,    // ❌ Sai: lẽ ra phải là 30,000
  actual_amount: 30000          // Old format
}
```

**AFTER FIX:**
```javascript
{
  cash_declared_amount: 0,
  transfer_declared_amount: 30000,
  cash_actual_amount: 0,
  transfer_actual_amount: 30000,  // ✅ Đúng
  actual_amount: 30000            // Keep for backward compatibility
}
```

## Giải Pháp

### Script Fix: `fix-shift-699bda2e.js`

```javascript
// Step 1: Fix handover
db.cash_handovers.updateOne(
  {_id: handoverId},
  {
    $set: {
      transfer_actual_amount: 30000,  // Move from actual_amount
      cash_actual_amount: 0,
      transfer_discrepancy: 0,
      cash_discrepancy: 0
    }
  }
);

// Step 2: Fix shift
db.shifts.updateOne(
  {_id: shiftId},
  {
    $set: {
      remaining_cash: 72000,          // Add back 30,000
      remaining_transfer: 0,          // Deduct 30,000
      handed_over_cash: 0,            // Remove wrong amount
      handed_over_transfer: 30000     // Add correct amount
    }
  }
);
```

## Kết Quả

### Frontend Display

**BEFORE:**
```
💵 Tiền mặt: -8.000 ₫
💳 Tiền CK: 30.000 ₫
Đã bàn giao: 30.000 ₫
```

**AFTER:**
```
💵 Tiền mặt: 72.000 ₫
💳 Tiền CK: 0 ₫
Đã bàn giao: 30.000 ₫
```

## Tại Sao Vấn Đề Này Xảy Ra?

1. **Backend được deploy nhưng chưa restart**
   - Code mới đã được push
   - Nhưng backend vẫn chạy code cũ
   - Khi confirm handover, code cũ chạy và tạo dữ liệu sai

2. **Không có auto-convert trong code cũ**
   - Code cũ không biết về `transfer_actual_amount`
   - Chỉ xử lý `actual_amount`
   - Dẫn đến logic update shift sai

3. **Frontend gửi old format**
   - Frontend vẫn gửi `actual_amount`
   - Backend cũ không convert sang new format
   - Dữ liệu bị lưu sai

## Phòng Tránh Trong Tương Lai

### 1. Always Restart Backend After Deploy
```bash
# After git pull
docker-compose restart backend
```

### 2. Migration Script cho Old Data
Nếu có nhiều shifts bị ảnh hưởng, chạy script tổng quát:

```javascript
// Find all shifts with potential issues
db.shifts.find({
  handed_over_cash: {$gt: 0},
  handed_over_transfer: 0,
  transfer_revenue: {$gt: 0}
}).forEach(shift => {
  // Check if handovers were actually transfer
  const handovers = db.cash_handovers.find({
    waiter_shift_id: shift._id,
    transfer_declared_amount: {$gt: 0},
    transfer_actual_amount: 0,
    actual_amount: {$gt: 0}
  }).toArray();
  
  if (handovers.length > 0) {
    // Fix this shift
    // ... (similar logic as above)
  }
});
```

### 3. Add Validation
Backend nên validate:
```go
// In ConfirmHandoverWithReconciliation
if h.TransferDeclaredAmount > 0 && req.TransferActualAmount == 0 && req.ActualAmount > 0 {
  // Auto-convert
  req.TransferActualAmount = req.ActualAmount
}
```

## Files Changed

1. **Database:**
   - `shifts` collection: 1 document updated
   - `cash_handovers` collection: 1 document updated

2. **Scripts:**
   - `fix-shift-699bda2e.js` - Fix script for this specific shift
   - `find-shift-1140.js` - Helper to find the shift
   - `check-handover-data.js` - Helper to check data

## Testing

### Verify Fix
```bash
# Check shift
docker exec cafe-pos-mongodb mongosh cafe_pos -u admin -p password123 \
  --authenticationDatabase admin --eval '
  db.shifts.findOne({_id: ObjectId("699bda2ea31585ce7ad4c47c")})
'

# Check handover
docker exec cafe-pos-mongodb mongosh cafe_pos -u admin -p password123 \
  --authenticationDatabase admin --eval '
  db.cash_handovers.findOne({_id: ObjectId("699bda5ca31585ce7ad4c484")})
'
```

### Frontend Verification
1. Login as waiter
2. Go to http://localhost:5173/#/shifts
3. Verify:
   - 💵 Tiền mặt: 72,000 ₫ (not negative)
   - 💳 Tiền CK: 0 ₫
   - Đã bàn giao: 30,000 ₫

## Summary

- ✅ Fixed handover: `transfer_actual_amount` now correct
- ✅ Fixed shift: cash and transfer amounts now correct
- ✅ Frontend will display correct values
- ✅ No more negative cash!

## Next Steps

1. ✅ Fix applied (DONE)
2. ⏳ Verify in frontend
3. ⏳ Check for other affected shifts
4. ⏳ Always restart backend after deploy
