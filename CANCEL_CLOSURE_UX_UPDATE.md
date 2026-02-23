# UX Update: Nút "Hủy đóng ca" thay vì Confirm Dialog

## Thay đổi

Thay vì hiển thị confirm dialog khi user bấm "Quay lại", giờ có một nút "Hủy đóng ca" rõ ràng trong UI.

## UI Changes

### Trước (Old UX)
```
[← Quay lại] → Confirm dialog: "Bạn có muốn hủy đóng ca?"
```

### Sau (New UX)
```
┌─────────────────────────────────────────┐
│ ⚠️ Đã bắt đầu đóng ca                   │
│                                         │
│ Nếu bạn muốn hủy quy trình đóng ca     │
│ và quay về trạng thái mở ca, bấm nút   │
│ bên dưới.                               │
│                                         │
│ [↩️ Hủy đóng ca]                        │
└─────────────────────────────────────────┘

Bước 2: Nhập tiền thực tế
[Nhập số tiền...]
[✓ Xác nhận tiền mặt]
```

## Ưu điểm

1. **Rõ ràng hơn**: User thấy ngay option để hủy đóng ca
2. **Ít nhầm lẫn**: Nút "Quay lại" chỉ đơn giản là quay về, không có side effect
3. **Intentional action**: User phải chủ động bấm nút "Hủy đóng ca", không phải quyết định trong confirm dialog
4. **Visual feedback**: Card màu cam nổi bật, dễ nhận biết trạng thái

## Khi nào hiển thị nút "Hủy đóng ca"?

- ✅ Status = `CLOSURE_INITIATED`
- ✅ Chưa nhập tiền thực tế (`actual_cash == null`)

## Khi nào ẩn nút "Hủy đóng ca"?

- ❌ Status = `OPEN` (chưa bắt đầu đóng ca)
- ❌ Đã nhập tiền thực tế (`actual_cash != null`)
- ❌ Status = `CLOSED` (đã đóng ca)

## Flow

```
1. User: Bấm "Bắt đầu đóng ca"
   → Status: OPEN → CLOSURE_INITIATED
   → Hiển thị card "Hủy đóng ca"

2. User: Bấm "↩️ Hủy đóng ca"
   → Confirm: "Bạn có chắc muốn hủy quy trình đóng ca?"
   → User: OK
   → API: POST /cancel-closure
   → Status: CLOSURE_INITIATED → OPEN
   → Alert: "✅ Đã hủy quy trình đóng ca thành công!"
   → Reload shift data
   → Ẩn card "Hủy đóng ca"
   → Hiển thị lại "Bước 1: Bắt đầu đóng ca"

3. User: Bấm "← Quay lại"
   → Navigate về /cashier
   → Không có confirm dialog
   → Ca vẫn giữ nguyên status
```

## Code Changes

### Template
```vue
<!-- Cancel Closure Option -->
<div v-if="shift.status === CASHIER_SHIFT_STATUS.CLOSURE_INITIATED && !shift.actual_cash" 
     class="bg-orange-50 border-2 border-orange-200 rounded-2xl p-4 shadow-sm">
  <div class="flex items-start gap-3 mb-3">
    <span class="text-2xl">⚠️</span>
    <div class="flex-1">
      <p class="font-bold text-orange-800 mb-1">Đã bắt đầu đóng ca</p>
      <p class="text-sm text-orange-700">
        Nếu bạn muốn hủy quy trình đóng ca và quay về trạng thái mở ca, bấm nút bên dưới.
      </p>
    </div>
  </div>
  <button
    @click="cancelClosure"
    :disabled="processing"
    class="w-full py-3 bg-orange-500 text-white rounded-xl font-bold active:scale-95 transition-transform disabled:opacity-50"
  >
    {{ processing ? 'Đang hủy...' : '↩️ Hủy đóng ca' }}
  </button>
</div>
```

### Script
```javascript
const cancelClosure = async () => {
  if (!confirm('Bạn có chắc muốn hủy quy trình đóng ca?\n\nCa sẽ quay về trạng thái mở.')) {
    return
  }
  
  processing.value = true
  error.value = null
  
  try {
    await cashierShiftService.cancelClosure(shift.value.id)
    await loadShift()
    alert('✅ Đã hủy quy trình đóng ca thành công!')
  } catch (err) {
    error.value = err.response?.data?.error || 'Không thể hủy đóng ca'
  } finally {
    processing.value = false
  }
}

const goBack = () => {
  router.push('/cashier')  // Simple navigation, no logic
}
```

## Testing

### Test Case 1: Cancel thành công
1. Bấm "Bắt đầu đóng ca"
2. Verify: Card "Hủy đóng ca" xuất hiện
3. Bấm "↩️ Hủy đóng ca"
4. Confirm: OK
5. Verify: Success message, card biến mất, status = OPEN

### Test Case 2: Không cho phép cancel
1. Bấm "Bắt đầu đóng ca"
2. Nhập tiền thực tế: 500000
3. Bấm "✓ Xác nhận tiền mặt"
4. Verify: Card "Hủy đóng ca" biến mất

### Test Case 3: User không confirm
1. Bấm "Bắt đầu đóng ca"
2. Bấm "↩️ Hủy đóng ca"
3. Confirm: Cancel
4. Verify: Vẫn ở trang đóng ca, status không đổi

### Test Case 4: Nút "Quay lại"
1. Bấm "Bắt đầu đóng ca"
2. Bấm "← Quay lại" (ở header)
3. Verify: Navigate về /cashier, không có confirm dialog
4. Quay lại trang đóng ca
5. Verify: Status vẫn là CLOSURE_INITIATED

## Screenshots

### Before (chưa nhập tiền)
```
┌────────────────────────────────────────┐
│ 📊 Thông tin ca làm                    │
│ [Ca details...]                        │
└────────────────────────────────────────┘

┌────────────────────────────────────────┐
│ ⚠️ Đã bắt đầu đóng ca                  │
│                                        │
│ Nếu bạn muốn hủy quy trình...         │
│                                        │
│ [↩️ Hủy đóng ca]                       │
└────────────────────────────────────────┘

┌────────────────────────────────────────┐
│ Bước 2: Nhập tiền thực tế              │
│                                        │
│ [Input: ________]                      │
│                                        │
│ [✓ Xác nhận tiền mặt]                  │
└────────────────────────────────────────┘
```

### After (đã nhập tiền)
```
┌────────────────────────────────────────┐
│ 📊 Thông tin ca làm                    │
│ [Ca details...]                        │
└────────────────────────────────────────┘

┌────────────────────────────────────────┐
│ Bước 3: Giải trình chênh lệch          │
│ (hoặc Bước 4: Xác nhận trách nhiệm)    │
│                                        │
│ [Continue workflow...]                 │
└────────────────────────────────────────┘
```

## Notes

- Card "Hủy đóng ca" chỉ xuất hiện khi có thể cancel (chưa nhập tiền thực tế)
- Màu cam (orange) để cảnh báo nhưng không quá nghiêm trọng như đỏ
- Icon ↩️ thể hiện action "quay lại" trạng thái trước
- Vẫn có confirm dialog để tránh bấm nhầm
- Success message để user biết action đã thành công
