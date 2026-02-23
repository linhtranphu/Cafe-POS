# Cashier Shift Closure Simplification

## Tổng quan

Đơn giản hóa quy trình đóng ca thu ngân từ 5 bước xuống còn 2-3 bước, bỏ các bước xác nhận thủ công.

## Thay đổi

### Quy trình cũ (5 bước):
1. Bắt đầu đóng ca
2. Nhập tiền thực tế
3. Giải trình chênh lệch (nếu có)
4. Xác nhận trách nhiệm ✂️ **BỎ**
5. Hoàn tất đóng ca ✂️ **BỎ**

### Quy trình mới (2-3 bước):
1. **Bắt đầu đóng ca** (kiểm tra ca waiter)
2. **Nhập tiền thực tế** → Tự động đóng ca (nếu không có chênh lệch)
3. **Giải trình chênh lệch** (nếu có) → Tự động đóng ca

## Chi tiết thay đổi

### 1. Frontend Changes

**File**: `frontend/src/views/CashierShiftClosure.vue`

#### Bước 1: Thêm kiểm tra waiter shifts
```vue
<!-- Step 1: Initiate Closure -->
<div v-if="shift.status === CASHIER_SHIFT_STATUS.OPEN">
  <!-- Waiter Shifts Warning -->
  <div v-if="waiterShiftsStatus && !waiterShiftsStatus.can_close">
    ⚠️ Không thể đóng ca! Còn X ca waiter đang mở
  </div>
  
  <button 
    @click="initiateClosure"
    :disabled="waiterShiftsStatus && !waiterShiftsStatus.can_close">
    ▶️ Bắt đầu đóng ca
  </button>
</div>
```

#### Bước 2: Tự động đóng ca sau khi nhập tiền
```javascript
const recordActualCash = async () => {
  await cashierShiftService.recordActualCash(shift.value.id, actualCash.value)
  await loadShift()
  
  // Auto close if no variance
  if (shift.value.variance && shift.value.variance.amount === 0) {
    await autoCloseShift()
  }
}
```

#### Bước 3: Giải trình và đóng ca
```javascript
const documentVarianceAndClose = async () => {
  // Document variance first
  await cashierShiftService.documentVariance(shift.value.id, {
    reason: varianceReason.value,
    notes: varianceNotes.value
  })
  
  // Then auto close
  await autoCloseShift()
}
```

#### Computed properties mới
```javascript
const needsVarianceDocumentation = computed(() => {
  return shift.value?.status === CLOSURE_INITIATED &&
         shift.value?.actual_cash &&
         shift.value?.variance?.amount !== 0 &&
         !shift.value?.variance?.reason
})

const readyToAutoClose = computed(() => {
  return shift.value?.status === CLOSURE_INITIATED &&
         shift.value?.actual_cash &&
         (!shift.value?.variance || 
          shift.value?.variance?.amount === 0 ||
          (shift.value?.variance?.reason && shift.value?.variance?.notes))
})
```

#### Watcher để tự động đóng ca
```javascript
watch(readyToAutoClose, async (isReady) => {
  if (isReady && !autoClosing.value) {
    autoClosing.value = true
    await autoCloseShift()
    autoClosing.value = false
  }
})
```

### 2. Backend Changes

**File**: `backend/domain/cashier/cashier_shift.go`

#### Bỏ yêu cầu confirmation trong CanClose()

**Trước**:
```go
func (cs *CashierShift) CanClose() error {
    if cs.Status != CashierShiftClosureInitiated {
        return errors.New("status must be ClosureInitiated")
    }
    
    // ❌ Yêu cầu confirmation
    if cs.Confirmation == nil {
        return errors.New("responsibility confirmation is required")
    }
    
    if cs.Variance != nil && cs.Variance.RequiresDocumentation() {
        if cs.Variance.Reason == nil || cs.Variance.Notes == "" {
            return errors.New("variance must be documented")
        }
    }
    
    return nil
}
```

**Sau**:
```go
func (cs *CashierShift) CanClose() error {
    if cs.Status != CashierShiftClosureInitiated {
        return errors.New("status must be ClosureInitiated")
    }
    
    // ✅ Chỉ yêu cầu actual cash
    if cs.ActualCash == nil {
        return errors.New("actual cash must be recorded")
    }
    
    if cs.Variance != nil && cs.Variance.RequiresDocumentation() {
        if cs.Variance.Reason == nil || cs.Variance.Notes == "" {
            return errors.New("variance must be documented")
        }
    }
    
    return nil
}
```

## Flow mới

### Trường hợp 1: Không có chênh lệch

```
1. User: Bấm "Bắt đầu đóng ca"
   → Kiểm tra waiter shifts
   → Nếu có ca mở: Hiển thị warning, disable button
   → Nếu không: Cho phép bắt đầu
   → Status: OPEN → CLOSURE_INITIATED

2. User: Nhập tiền thực tế = system_cash
   → Gọi API record-actual-cash
   → Variance = 0
   → Frontend tự động gọi API close
   → Status: CLOSURE_INITIATED → CLOSED
   → Hiển thị "✅ Ca làm đã đóng"
```

### Trường hợp 2: Có chênh lệch

```
1. User: Bấm "Bắt đầu đóng ca"
   → Status: OPEN → CLOSURE_INITIATED

2. User: Nhập tiền thực tế ≠ system_cash
   → Gọi API record-actual-cash
   → Variance ≠ 0
   → Hiển thị "Bước 3: Giải trình chênh lệch"

3. User: Chọn lý do + nhập ghi chú
   → Bấm "🔒 Ghi nhận và đóng ca"
   → Gọi API document-variance
   → Gọi API close
   → Status: CLOSURE_INITIATED → CLOSED
   → Hiển thị "✅ Ca làm đã đóng"
```

## Ưu điểm

1. **Nhanh hơn**: Giảm từ 5 bước xuống 2-3 bước
2. **Ít click hơn**: Không cần bấm "Xác nhận trách nhiệm" và "Hoàn tất đóng ca"
3. **Tự động hóa**: Tự động đóng ca sau khi hoàn thành các bước bắt buộc
4. **UX tốt hơn**: Rõ ràng, không gây nhầm lẫn
5. **Kiểm tra sớm**: Kiểm tra waiter shifts ngay ở bước 1

## Validation

### Frontend
- ✅ Kiểm tra waiter shifts trước khi bắt đầu đóng ca
- ✅ Tự động đóng ca khi đủ điều kiện
- ✅ Hiển thị loading state khi đang đóng ca

### Backend
- ✅ Status = CLOSURE_INITIATED
- ✅ Actual cash đã được nhập
- ✅ Variance đã được giải trình (nếu có)
- ❌ Không yêu cầu confirmation nữa

## Testing

### Test Case 1: Không có chênh lệch
1. Bắt đầu đóng ca
2. Nhập tiền = system_cash
3. Verify: Ca tự động đóng

### Test Case 2: Có chênh lệch
1. Bắt đầu đóng ca
2. Nhập tiền ≠ system_cash
3. Verify: Hiển thị form giải trình
4. Giải trình chênh lệch
5. Verify: Ca tự động đóng

### Test Case 3: Có ca waiter mở
1. Verify: Hiển thị warning
2. Verify: Button "Bắt đầu đóng ca" bị disable
3. Đóng tất cả ca waiter
4. Verify: Button được enable

## Files Changed

### Frontend
- `frontend/src/views/CashierShiftClosure.vue`
  - Bỏ Step 4 (Confirm Responsibility)
  - Bỏ Step 5 (Close Shift)
  - Thêm waiter shifts check vào Step 1
  - Thêm auto-close logic
  - Thêm watcher cho readyToAutoClose

### Backend
- `backend/domain/cashier/cashier_shift.go`
  - Cập nhật `CanClose()` method
  - Bỏ yêu cầu `Confirmation != nil`
  - Chỉ yêu cầu `ActualCash != nil`

## Deployment

1. ✅ Backend changes deployed
2. ✅ Backend restarted
3. ⏳ Frontend rebuild required
4. ⏳ Manual testing

## Notes

- Method `ConfirmResponsibility()` vẫn tồn tại trong code nhưng không được sử dụng
- Có thể xóa method này và related code trong tương lai nếu không cần
- Audit log vẫn ghi nhận đầy đủ các actions
- State machine vẫn hoạt động bình thường
