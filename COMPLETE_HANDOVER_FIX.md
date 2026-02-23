# COMPLETE HANDOVER FIX - Sửa triệt để logic bàn giao

## VẤN ĐỀ

1. Backend có 2 luồng xử lý không nhất quán (old vs new)
2. Frontend cashier view chỉ có 1 input actual_amount - không phân biệt cash/transfer
3. Logic trừ tiền không rõ ràng - có chỗ trừ 2 lần, có chỗ không trừ

## GIẢI PHÁP

### 1. Backend - Đơn giản hóa thành 1 luồng duy nhất

**Nguyên tắc:**
- Khi CREATE handover: TRỪ ngay remaining_cash và remaining_transfer
- Khi CONFIRM handover: CHỈ cập nhật handed_over_cash và handed_over_transfer
- KHÔNG BAO GIỜ trừ 2 lần

**Functions giữ lại:**
- `CreateHandover()` - tạo handover (hỗ trợ cả cash và transfer riêng biệt)
- `ConfirmHandover()` - confirm handover (hỗ trợ cả cash và transfer riêng biệt)
- `updateBalances()` - cập nhật balances (CHỈ update handed_over, KHÔNG trừ remaining)

**Functions XÓA:**
- `CreateHandoverWithTransfer()` - duplicate
- `ConfirmHandoverWithReconciliation()` - duplicate
- `ConfirmHandoverWithDualAmounts()` - duplicate
- `updateCashAmounts()` - logic sai
- `updateDualBalances()` - duplicate

### 2. Frontend Cashier View - Thêm input riêng cho cash và transfer

**Hiện tại:**
```vue
<input v-model.number="confirmForm.actual_amount" />
```

**Sửa thành:**
```vue
<!-- Nếu handover có cash -->
<input v-model.number="confirmForm.actual_cash_amount" />

<!-- Nếu handover có transfer -->
<input v-model.number="confirmForm.actual_transfer_amount" />
```

### 3. Frontend Waiter View - Đã OK

Form bàn giao đã có 3 options rõ ràng:
- 💵 Tiền mặt only
- 💳 Tiền CK only
- 💰 Cả hai

## IMPLEMENTATION

Sẽ fix theo thứ tự:
1. Backend service - đơn giản hóa
2. Backend handler - đơn giản hóa
3. Frontend cashier view - thêm input riêng
4. Test toàn bộ flow
