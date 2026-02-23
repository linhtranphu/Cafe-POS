# Frontend ShiftView - Tách bạch Cash và Transfer

## VẤN ĐỀ

Frontend có logic fallback phức tạp:
```vue
{{ currentShift.remaining_cash || currentShift.current_cash || 0 }}
{{ currentShift.remaining_transfer || currentShift.transfer_revenue || 0 }}
```

**Hậu quả:**
- Khi `remaining_transfer = 0` → fallback sang `transfer_revenue`
- `transfer_revenue` là tổng doanh thu, KHÔNG bị trừ khi bàn giao
- User thấy số tiền CK không giảm sau khi bàn giao

## GIẢI PHÁP

Tách bạch hoàn toàn, KHÔNG fallback:

### 1. Hiển thị 3 sections riêng biệt

**TRƯỚC (1 section, 3 cột):**
```vue
<div class="grid grid-cols-3">
  <div>💵 Tiền mặt: {{ remaining_cash || current_cash }}</div>
  <div>💳 Tiền CK: {{ remaining_transfer || transfer_revenue }}</div>
  <div>Đã bàn giao: {{ handed_over_cash + handed_over_transfer }}</div>
</div>
```

**SAU (3 sections riêng):**
```vue
<!-- 1. Doanh thu -->
<div>
  <p>📊 Doanh thu</p>
  <div>💵 Tiền mặt: {{ current_cash || 0 }}</div>
  <div>💳 Tiền CK: {{ transfer_revenue || 0 }}</div>
</div>

<!-- 2. Còn lại (chưa bàn giao) -->
<div>
  <p>💰 Còn lại (chưa bàn giao)</p>
  <div>💵 Tiền mặt: {{ remaining_cash || 0 }}</div>
  <div>💳 Tiền CK: {{ remaining_transfer || 0 }}</div>
</div>

<!-- 3. Đã bàn giao -->
<div>
  <p>✅ Đã bàn giao</p>
  <div>💵 Tiền mặt: {{ handed_over_cash || 0 }}</div>
  <div>💳 Tiền CK: {{ handed_over_transfer || 0 }}</div>
</div>
```

### 2. Buttons logic đơn giản

**TRƯỚC:**
```vue
<button v-if="(remaining_cash || current_cash || 0) > 0">
  Bàn giao
</button>
```

**SAU:**
```vue
<button v-if="remaining_cash > 0 || remaining_transfer > 0">
  Bàn giao
</button>
```

### 3. Form inputs với max rõ ràng

**TRƯỚC:**
```vue
<input :max="remaining_cash || current_cash || 0" />
<input :max="remaining_transfer || transfer_revenue || 0" />
```

**SAU:**
```vue
<input :max="remaining_cash || 0" />
<p>Tối đa: {{ formatPrice(remaining_cash || 0) }}</p>

<input :max="remaining_transfer || 0" />
<p>Tối đa: {{ formatPrice(remaining_transfer || 0) }}</p>
```

### 4. Function createHandoverAndEndShift

**TRƯỚC:**
```js
const handoverData = {
  declared_amount: currentShift.value?.remaining_cash || currentShift.value?.current_cash || 0,
  waiter_note: ...,
  end_cash: ...
}
await shiftStore.createHandoverAndEndShift(...)
```

**SAU:**
```js
const handoverData = {
  cash_amount: currentShift.value?.remaining_cash || 0,
  transfer_amount: currentShift.value?.remaining_transfer || 0,
  handover_type: 'END_SHIFT',
  waiter_note: ...
}
await shiftStore.createCashHandover(...)
```

## THAY ĐỔI CHI TIẾT

### Các chỗ đã sửa:

1. **Line 48-62**: Cash Status section - tách thành 3 sections riêng
2. **Line 80-100**: Action buttons - logic đơn giản hơn
3. **Line 330-340**: Partial handover form - balance info rõ ràng
4. **Line 368-390**: Input fields - thêm max hint
5. **Line 450-460**: Handover end shift form - cash summary
6. **Line 707-730**: createHandoverAndEndShift function - dùng cash_amount/transfer_amount

### Các field được dùng:

**Doanh thu (Revenue):**
- `current_cash` - tổng tiền mặt đã thu
- `transfer_revenue` - tổng tiền CK đã thu

**Còn lại (Remaining):**
- `remaining_cash` - tiền mặt chưa bàn giao
- `remaining_transfer` - tiền CK chưa bàn giao

**Đã bàn giao (Handed Over):**
- `handed_over_cash` - tiền mặt đã bàn giao
- `handed_over_transfer` - tiền CK đã bàn giao

## LƯU Ý

**Vấn đề còn lại cần fix ở backend:**

Khi thu tiền từ order, backend cần cập nhật CẢ 2 fields:
```go
// Trong order_service.go - CollectPayment()
if paymentMethod == TRANSFER || paymentMethod == QR {
    shift.TransferRevenue += order.Total      // ✅ Đã có
    shift.RemainingTransfer += order.Total    // ❌ THIẾU - CẦN THÊM
}

if paymentMethod == CASH {
    shift.CurrentCash += order.Total          // ✅ Đã có
    shift.RemainingCash += order.Total        // ✅ Đã có
}
```

**Nếu không fix backend:**
- `remaining_transfer` sẽ luôn = 0
- Frontend sẽ hiển thị "Tiền CK còn lại: 0 ₫"
- User không thể bàn giao tiền CK

## KẾT QUẢ

✅ Frontend hiển thị rõ ràng 3 loại số liệu
✅ Không còn fallback logic phức tạp
✅ User thấy chính xác số tiền còn lại để bàn giao
✅ Dễ debug khi có vấn đề

❌ Backend vẫn cần fix để cập nhật `remaining_transfer` khi thu tiền
