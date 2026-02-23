# Fix: Bước 5 hiển thị sai thời điểm

## Vấn đề

"Bước 5: Hoàn tất đóng ca" hiển thị ngay cả khi các bước 2, 3, 4 chưa hoàn thành.

## Nguyên nhân

Computed `canCloseShift` chỉ kiểm tra:
```javascript
const canCloseShift = computed(() => {
  if (!shift.value) return false
  if (shift.value.status !== CASHIER_SHIFT_STATUS.CLOSURE_INITIATED) return false
  return shift.value.confirmation !== null  // ❌ Chỉ kiểm tra confirmation
})
```

Điều này có nghĩa là nếu `confirmation` có giá trị (từ lần đóng ca trước hoặc lỗi data), bước 5 sẽ hiển thị ngay cả khi chưa nhập tiền thực tế hoặc chưa giải trình chênh lệch.

## Giải pháp

Cập nhật logic `canCloseShift` để kiểm tra đầy đủ tất cả các bước:

```javascript
const canCloseShift = computed(() => {
  if (!shift.value) return false
  if (shift.value.status !== CASHIER_SHIFT_STATUS.CLOSURE_INITIATED) return false
  
  // Must have recorded actual cash (Bước 2)
  if (!shift.value.actual_cash) return false
  
  // Must have confirmed responsibility (Bước 4)
  if (!shift.value.confirmation) return false
  
  // If there's variance, it must be documented (Bước 3)
  if (shift.value.variance && shift.value.variance.amount !== 0) {
    if (!shift.value.variance.reason || !shift.value.variance.notes) {
      return false
    }
  }
  
  return true
})
```

## Quy trình đúng

### Trường hợp không có chênh lệch:
1. Bước 1: Bắt đầu đóng ca ✅
2. Bước 2: Nhập tiền thực tế ✅
3. (Bước 3: Bỏ qua - không có chênh lệch)
4. Bước 4: Xác nhận trách nhiệm ✅
5. Bước 5: Hoàn tất đóng ca ✅ (hiển thị sau bước 4)

### Trường hợp có chênh lệch:
1. Bước 1: Bắt đầu đóng ca ✅
2. Bước 2: Nhập tiền thực tế ✅
3. Bước 3: Giải trình chênh lệch ✅ (bắt buộc)
4. Bước 4: Xác nhận trách nhiệm ✅
5. Bước 5: Hoàn tất đóng ca ✅ (hiển thị sau bước 4)

## Validation Logic

```
canCloseShift = true khi:
  ✅ status === CLOSURE_INITIATED
  ✅ actual_cash !== null (đã nhập tiền)
  ✅ confirmation !== null (đã xác nhận)
  ✅ Nếu variance.amount !== 0:
     ✅ variance.reason !== null
     ✅ variance.notes !== ""
```

## Testing

### Test Case 1: Không có chênh lệch
1. Bắt đầu đóng ca
2. Verify: Bước 5 KHÔNG hiển thị
3. Nhập tiền thực tế = system_cash
4. Verify: Bước 5 KHÔNG hiển thị (chưa confirm)
5. Xác nhận trách nhiệm
6. Verify: Bước 5 HIỂN THỊ

### Test Case 2: Có chênh lệch
1. Bắt đầu đóng ca
2. Nhập tiền thực tế khác system_cash
3. Verify: Bước 3 hiển thị, Bước 5 KHÔNG hiển thị
4. Giải trình chênh lệch
5. Verify: Bước 4 hiển thị, Bước 5 KHÔNG hiển thị
6. Xác nhận trách nhiệm
7. Verify: Bước 5 HIỂN THỊ

### Test Case 3: Chưa nhập tiền
1. Bắt đầu đóng ca
2. Verify: Chỉ hiển thị Bước 2, KHÔNG hiển thị Bước 5

## Files Changed

- `frontend/src/views/CashierShiftClosure.vue` - Updated `canCloseShift` computed

## Notes

- Logic này đảm bảo user phải đi qua tất cả các bước theo thứ tự
- Không thể skip bước nào
- UI rõ ràng, không gây nhầm lẫn
- Backend vẫn có validation riêng để đảm bảo data integrity
