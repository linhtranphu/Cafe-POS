# Fix Modal "Bàn giao và đóng ca" - Hiển thị cả Cash và Transfer

## VẤN ĐỀ

Modal "Bàn giao và đóng ca" chỉ hiển thị tiền mặt, không hiển thị tiền CK.

**Trước:**
```
💰 Tiền sẽ bàn giao: 50,000 ₫  (chỉ tiền mặt)
```

**Thực tế shift có:**
- Tiền mặt còn lại: 50,000 ₫
- Tiền CK còn lại: 30,000 ₫

**Kết quả**: User không biết sẽ bàn giao cả tiền CK.

## GIẢI PHÁP

Hiển thị chi tiết cả cash và transfer trong modal.

### Code cũ (CHỈ HIỂN THỊ TIỀN MẶT)
```vue
<div class="bg-orange-50 p-4 rounded-xl mb-4">
  <div class="space-y-2">
    <div class="flex justify-between items-center">
      <span class="text-sm text-gray-600">Tiền sẽ bàn giao</span>
      <span class="font-bold text-2xl text-orange-600">
        {{ formatPrice(currentShift?.remaining_cash || 0) }}
      </span>
    </div>
    <div class="flex justify-between items-center text-sm">
      <span class="text-gray-500">Tiền cuối ca</span>
      <span class="font-medium">{{ formatPrice(handoverEndShiftForm.end_cash) }}</span>
    </div>
  </div>
</div>
```

### Code mới (HIỂN THỊ ĐẦY ĐỦ)
```vue
<div class="bg-orange-50 p-4 rounded-xl mb-4">
  <p class="text-sm font-medium text-gray-700 mb-3">💰 Số tiền sẽ bàn giao</p>
  <div class="space-y-2">
    <!-- Cash -->
    <div v-if="(currentShift?.remaining_cash || 0) > 0" 
      class="flex justify-between items-center">
      <span class="text-sm text-gray-600">💵 Tiền mặt</span>
      <span class="font-bold text-lg text-green-600">
        {{ formatPrice(currentShift?.remaining_cash || 0) }}
      </span>
    </div>
    
    <!-- Transfer -->
    <div v-if="(currentShift?.remaining_transfer || 0) > 0" 
      class="flex justify-between items-center">
      <span class="text-sm text-gray-600">💳 Tiền CK</span>
      <span class="font-bold text-lg text-blue-600">
        {{ formatPrice(currentShift?.remaining_transfer || 0) }}
      </span>
    </div>
    
    <!-- Total (if both exist) -->
    <div v-if="(currentShift?.remaining_cash || 0) > 0 && (currentShift?.remaining_transfer || 0) > 0" 
      class="flex justify-between items-center pt-2 border-t border-orange-200">
      <span class="text-sm font-medium text-gray-700">Tổng cộng</span>
      <span class="font-bold text-xl text-orange-600">
        {{ formatPrice((currentShift?.remaining_cash || 0) + (currentShift?.remaining_transfer || 0)) }}
      </span>
    </div>
    
    <!-- End cash -->
    <div class="flex justify-between items-center text-sm pt-2 border-t border-orange-200">
      <span class="text-gray-500">Tiền cuối ca</span>
      <span class="font-medium">{{ formatPrice(handoverEndShiftForm.end_cash) }}</span>
    </div>
  </div>
</div>
```

## HIỂN THỊ THEO TRƯỜNG HỢP

### Trường hợp 1: Chỉ có tiền mặt
```
💰 Số tiền sẽ bàn giao
💵 Tiền mặt: 50,000 ₫
─────────────────────
Tiền cuối ca: 0 ₫
```

### Trường hợp 2: Chỉ có tiền CK
```
💰 Số tiền sẽ bàn giao
💳 Tiền CK: 30,000 ₫
─────────────────────
Tiền cuối ca: 0 ₫
```

### Trường hợp 3: Có cả hai
```
💰 Số tiền sẽ bàn giao
💵 Tiền mặt: 50,000 ₫
💳 Tiền CK: 30,000 ₫
─────────────────────
Tổng cộng: 80,000 ₫
─────────────────────
Tiền cuối ca: 0 ₫
```

## LOGIC BACKEND (ĐÃ ĐÚNG)

Function `createHandoverAndEndShift` đã gửi đúng cả 2 amounts:
```js
const handoverData = {
  cash_amount: currentShift.value?.remaining_cash || 0,
  transfer_amount: currentShift.value?.remaining_transfer || 0,
  handover_type: 'END_SHIFT',
  waiter_note: handoverEndShiftForm.value.waiter_note
}
```

Backend `CreateHandover` sẽ:
1. Trừ `remaining_cash` theo `cash_amount`
2. Trừ `remaining_transfer` theo `transfer_amount`
3. Tạo handover với status PENDING

Khi cashier confirm:
1. Cập nhật `handed_over_cash` và `handed_over_transfer`
2. Đóng shift nếu `handover_type = END_SHIFT`

## WARNING TEXT CẬP NHẬT

**Trước:**
> Thao tác này sẽ bàn giao toàn bộ tiền còn lại và tự động đóng ca sau khi cashier xác nhận.

**Sau:**
> Thao tác này sẽ bàn giao toàn bộ tiền mặt và tiền CK còn lại, sau đó tự động đóng ca khi cashier xác nhận.

## KẾT QUẢ

✅ User thấy rõ sẽ bàn giao bao nhiêu tiền mặt
✅ User thấy rõ sẽ bàn giao bao nhiêu tiền CK
✅ User thấy tổng cộng nếu có cả 2 loại
✅ Không còn nhầm lẫn về số tiền bàn giao
✅ UI responsive - chỉ hiển thị loại tiền nào có giá trị > 0

## FILES THAY ĐỔI

- `frontend/src/views/ShiftView.vue` - Modal "Bàn giao và đóng ca"
  - Line ~470-500: Cập nhật hiển thị money summary
  - Thêm conditional rendering cho cash và transfer
  - Thêm tổng cộng khi có cả 2 loại
