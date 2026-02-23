# Cashier Cancel Closure với Rollback

## Vấn đề

Khi đóng ca cashier:
1. Bước 1: Bắt đầu đóng ca → Status = CLOSURE_INITIATED ✅
2. Bước 2: Nhập tiền thực tế → `actual_cash` và `variance` được lưu ✅
3. Bấm "Quay lại" hoặc muốn hủy → KHÔNG có rollback ❌

**Hậu quả**: 
- Domain model có `actual_cash` và `variance`
- Nhưng status vẫn là CLOSURE_INITIATED
- Không thể tiếp tục đóng ca, không thể hủy
- Data stuck ở trạng thái inconsistent

## Root Cause

### Domain Model Restriction (Trước khi fix)
```go
func (cs *CashierShift) CancelClosure(...) error {
    // Cannot cancel if actual cash has been recorded
    if cs.ActualCash != nil {
        return errors.New("cannot cancel closure: actual cash has been recorded")
    }
    // ...
}
```

### State Machine Restriction (Trước khi fix)
```go
func (sm *ShiftStateMachine) CanCancelClosure(shift *CashierShift) bool {
    if shift.Status != CashierShiftClosureInitiated {
        return false
    }
    
    // Cannot cancel if actual cash has been recorded
    if shift.ActualCash != nil {
        return false
    }
    
    return true
}
```

### Frontend Limitation (Trước khi fix)
```vue
<!-- Only show cancel button when NO actual_cash -->
<div v-if="shift.status === CLOSURE_INITIATED && !shift.actual_cash">
  <button @click="cancelClosure">Hủy đóng ca</button>
</div>
```

## Giải pháp

### 1. Update Domain Model - Allow Cancel with Rollback

**File**: `backend/domain/cashier/cashier_shift.go`

```go
func (cs *CashierShift) CancelClosure(userID, deviceID string, timestamp time.Time) error {
    // Validate current status is ClosureInitiated
    if cs.Status != CashierShiftClosureInitiated {
        return errors.New("cannot cancel closure: shift status must be ClosureInitiated")
    }

    // Create audit log with rollback info
    auditData := map[string]interface{}{
        "had_actual_cash": cs.ActualCash != nil,
        "had_variance":    cs.Variance != nil,
    }
    auditEntry, err := NewAuditLogEntry("closure_cancelled", userID, deviceID, timestamp, auditData)
    if err != nil {
        return err
    }

    // ✅ NEW: Rollback closure data
    cs.ActualCash = nil
    cs.Variance = nil

    // Transition status back to Open
    cs.Status = CashierShiftOpen

    // Add audit log entry
    cs.AuditLog = append(cs.AuditLog, *auditEntry)
    cs.UpdatedAt = timestamp

    return nil
}
```

**Key Changes**:
- ❌ Removed: Check `if cs.ActualCash != nil`
- ✅ Added: Rollback `cs.ActualCash = nil` and `cs.Variance = nil`
- ✅ Added: Audit log with rollback info

### 2. Update State Machine - Allow Cancel Anytime

**File**: `backend/domain/cashier/cashier_shift_state_machine.go`

```go
func (sm *ShiftStateMachine) CanCancelClosure(shift *CashierShift) bool {
    // ✅ Can cancel anytime during CLOSURE_INITIATED state
    // Will rollback actual_cash and variance if they exist
    return shift.Status == CashierShiftClosureInitiated
}
```

**Key Changes**:
- ❌ Removed: Check `if shift.ActualCash != nil`
- ✅ Simplified: Only check status

### 3. Update Frontend - Show Cancel Button Always

**File**: `frontend/src/views/CashierShiftClosure.vue`

```vue
<!-- Show cancel button ANYTIME during CLOSURE_INITIATED -->
<div v-if="shift.status === CASHIER_SHIFT_STATUS.CLOSURE_INITIATED && 
           shift.status !== CASHIER_SHIFT_STATUS.CLOSED" 
     class="bg-orange-50 border-2 border-orange-200 rounded-2xl p-4 shadow-sm">
  <div class="flex items-start gap-3 mb-3">
    <span class="text-2xl">⚠️</span>
    <div class="flex-1">
      <p class="font-bold text-orange-800 mb-1">Đang trong quy trình đóng ca</p>
      <p class="text-sm text-orange-700">
        Nếu bạn muốn hủy quy trình đóng ca và quay về trạng thái mở ca, bấm nút bên dưới.
        <span v-if="shift.actual_cash" class="font-semibold">
          Lưu ý: Tiền thực tế và chênh lệch đã nhập sẽ bị xóa.
        </span>
      </p>
    </div>
  </div>
  <button @click="cancelClosure" :disabled="processing">
    {{ processing ? 'Đang hủy...' : '↩️ Hủy đóng ca' }}
  </button>
</div>
```

**Key Changes**:
- ❌ Removed: Condition `&& !shift.actual_cash`
- ✅ Added: Warning message when `actual_cash` exists
- ✅ Shows button at ANY step during closure

### 4. Enhanced Confirmation Message

```javascript
const cancelClosure = async () => {
  // Build confirmation message based on what will be rolled back
  let confirmMessage = 'Bạn có chắc muốn hủy quy trình đóng ca?\n\nCa sẽ quay về trạng thái mở.'
  
  if (shift.value.actual_cash) {
    confirmMessage += '\n\n⚠️ Các dữ liệu sau sẽ bị xóa:'
    confirmMessage += '\n• Tiền thực tế đã nhập'
    if (shift.value.variance) {
      confirmMessage += '\n• Chênh lệch đã tính'
      if (shift.value.variance.reason) {
        confirmMessage += '\n• Giải trình chênh lệch'
      }
    }
  }
  
  if (!confirm(confirmMessage)) {
    return
  }
  
  // Call API...
}
```

## Flow Comparison

### Before (BUG):
```
1. Bắt đầu đóng ca
   → Status: OPEN → CLOSURE_INITIATED ✅

2. Nhập tiền thực tế: 500k
   → actual_cash: 500k ✅
   → variance: -50k ✅
   → Status: CLOSURE_INITIATED ✅

3. Bấm "Quay lại" hoặc muốn hủy
   → ❌ Nút "Hủy đóng ca" không hiển thị
   → ❌ Không thể cancel
   → ❌ Data stuck: có actual_cash nhưng status = CLOSURE_INITIATED
   → ❌ Phải tiếp tục đóng ca hoặc manual fix DB
```

### After (FIXED):
```
1. Bắt đầu đóng ca
   → Status: OPEN → CLOSURE_INITIATED ✅
   → Nút "Hủy đóng ca" hiển thị ✅

2. Nhập tiền thực tế: 500k
   → actual_cash: 500k ✅
   → variance: -50k ✅
   → Status: CLOSURE_INITIATED ✅
   → Nút "Hủy đóng ca" vẫn hiển thị ✅
   → Warning: "Tiền thực tế và chênh lệch sẽ bị xóa" ✅

3. Bấm "Hủy đóng ca"
   → Confirm dialog với chi tiết rollback ✅
   → START TRANSACTION
   → actual_cash: 500k → null ✅
   → variance: -50k → null ✅
   → Status: CLOSURE_INITIATED → OPEN ✅
   → Audit log: "closure_cancelled" with rollback info ✅
   → COMMIT TRANSACTION
   → Success message ✅
```

## Transaction Protection

Vì đã có transaction từ implementation trước:

```go
func (s *CashierShiftService) CancelClosureWithTransaction(
    ctx context.Context,
    shiftID primitive.ObjectID,
    userID, deviceID string,
) (*cashier.CashierShift, error) {
    session, err := s.mongoClient.StartSession()
    defer session.EndSession(ctx)
    
    _, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
        // 1. Get shift
        // 2. Validate
        // 3. Cancel closure (with rollback)
        // 4. Save shift
        return nil, nil
    })
    
    return updatedShift, err
}
```

**Benefits**:
- ✅ Atomicity: Rollback thành công hoặc không thay đổi gì
- ✅ Consistency: Domain và DB luôn sync
- ✅ Error Recovery: Auto rollback nếu có lỗi

## Audit Trail

Audit log ghi nhận rollback:

```json
{
  "action": "closure_cancelled",
  "user_id": "cashier123",
  "device_id": "web",
  "timestamp": "2026-02-23T21:30:00Z",
  "data": {
    "had_actual_cash": true,
    "had_variance": true
  }
}
```

Giúp:
- Track được khi nào cancel
- Biết có rollback data không
- Audit compliance

## User Experience

### Scenario 1: Cancel Before Recording Cash
```
User: Bắt đầu đóng ca
UI: Hiển thị nút "Hủy đóng ca" (màu cam)
User: Bấm "Hủy đóng ca"
UI: Confirm: "Bạn có chắc muốn hủy quy trình đóng ca?"
User: OK
Result: ✅ Ca quay về OPEN, không mất data
```

### Scenario 2: Cancel After Recording Cash
```
User: Bắt đầu đóng ca
User: Nhập tiền thực tế: 500k
UI: Hiển thị nút "Hủy đóng ca" với warning
    "Lưu ý: Tiền thực tế và chênh lệch đã nhập sẽ bị xóa"
User: Bấm "Hủy đóng ca"
UI: Confirm: "Bạn có chắc muốn hủy quy trình đóng ca?
    
    ⚠️ Các dữ liệu sau sẽ bị xóa:
    • Tiền thực tế đã nhập
    • Chênh lệch đã tính"
User: OK
Result: ✅ Ca quay về OPEN, actual_cash và variance bị xóa
```

### Scenario 3: Cancel After Documenting Variance
```
User: Bắt đầu đóng ca
User: Nhập tiền thực tế: 500k (chênh lệch -50k)
User: Giải trình chênh lệch: "Lỗi đếm tiền"
UI: Hiển thị nút "Hủy đóng ca" với warning
User: Bấm "Hủy đóng ca"
UI: Confirm: "Bạn có chắc muốn hủy quy trình đóng ca?
    
    ⚠️ Các dữ liệu sau sẽ bị xóa:
    • Tiền thực tế đã nhập
    • Chênh lệch đã tính
    • Giải trình chênh lệch"
User: OK
Result: ✅ Ca quay về OPEN, tất cả data bị xóa
```

## Benefits

### 1. Flexibility ✅
- Có thể cancel bất kỳ lúc nào trong quá trình đóng ca
- Không bị stuck ở trạng thái inconsistent

### 2. Data Integrity ✅
- Rollback đảm bảo data clean
- Không có orphan data (actual_cash without closure)

### 3. User Control ✅
- User có quyền kiểm soát quy trình
- Clear warning về consequences

### 4. Audit Compliance ✅
- Audit log ghi nhận mọi rollback
- Traceability đầy đủ

### 5. Error Recovery ✅
- Nếu nhập sai, có thể cancel và nhập lại
- Không cần manual DB fix

## Testing

### Test Case 1: Cancel Before Recording Cash
```
1. Start closure
2. Click "Hủy đóng ca"
3. Confirm
Expected: Status = OPEN, no data lost
```

### Test Case 2: Cancel After Recording Cash
```
1. Start closure
2. Record actual cash: 500k
3. Click "Hủy đóng ca"
4. Confirm (with warning)
Expected: Status = OPEN, actual_cash = null, variance = null
```

### Test Case 3: Cancel After Documenting Variance
```
1. Start closure
2. Record actual cash: 500k (variance -50k)
3. Document variance: "Counting error"
4. Click "Hủy đóng ca"
5. Confirm (with detailed warning)
Expected: Status = OPEN, all closure data cleared
```

### Test Case 4: Transaction Rollback on Error
```
1. Start closure
2. Record actual cash: 500k
3. Simulate DB error during cancel
Expected: Transaction rollback, shift unchanged
```

## Files Changed

1. `backend/domain/cashier/cashier_shift.go`
   - Updated `CancelClosure` to rollback actual_cash and variance
   - Added audit data for rollback tracking

2. `backend/domain/cashier/cashier_shift_state_machine.go`
   - Simplified `CanCancelClosure` to allow cancel anytime

3. `frontend/src/views/CashierShiftClosure.vue`
   - Show cancel button during entire CLOSURE_INITIATED state
   - Enhanced confirmation message with rollback details
   - Added warning when actual_cash exists

## Deployment

- [x] Domain model updated
- [x] State machine updated
- [x] Frontend updated
- [x] Backend restarted
- [ ] Manual testing required
- [ ] Production deployment

## Related

- `CASHIER_CLOSURE_TRANSACTION.md` - Transaction implementation
- `CANCEL_CLOSURE_IMPLEMENTATION.md` - Original cancel implementation
- `TRANSACTION_IMPLEMENTATION_SUMMARY.md` - Overall transaction summary
