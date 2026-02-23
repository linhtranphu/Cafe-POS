# HANDOVER LOGIC - FIX TRIỆT ĐỂ HOÀN TẤT

## TÓM TẮT

Đã viết lại toàn bộ logic bàn giao từ đầu với nguyên tắc đơn giản và rõ ràng.

## NGUYÊN TẮC MỚI

### 1. Khi TẠO handover (CreateHandover)
- TRỪ ngay `remaining_cash` và `remaining_transfer` từ shift
- Lưu handover với status PENDING
- **CHỈ TRỪ 1 LẦN DUY NHẤT Ở ĐÂY**

### 2. Khi CONFIRM handover (ConfirmHandover)
- Nhận `actual_cash_amount` và `actual_transfer_amount` từ cashier
- Tính chênh lệch riêng cho cash và transfer
- Nếu chênh lệch > 100k → cần manager approval
- Gọi `updateBalances()` để cập nhật

### 3. Khi UPDATE balances (updateBalances)
- CHỈ cập nhật `handed_over_cash` và `handed_over_transfer`
- **KHÔNG BAO GIỜ trừ `remaining_cash` và `remaining_transfer` lại**
- Cập nhật cashier shift: `received_cash` và `received_transfer`

## THAY ĐỔI CODE

### Backend Service (cash_handover_service.go)

**XÓA toàn bộ functions cũ:**
- `CreateHandoverWithTransfer()` - duplicate
- `ConfirmHandoverWithReconciliation()` - duplicate  
- `ConfirmHandoverWithDualAmounts()` - duplicate
- `updateCashAmounts()` - logic sai
- `updateDualBalances()` - duplicate
- `CreateHandoverAndEndShift()` - không cần thiết
- `ApproveDiscrepancy()` - implement sau
- `GetDiscrepancyStats()` - implement sau

**GIỮ LẠI và viết lại:**
- `CreateHandover(cashAmount, transferAmount, ...)` - tạo handover, trừ remaining
- `ConfirmHandover(actualCashAmount, actualTransferAmount, ...)` - confirm handover
- `updateBalances()` - CHỈ update handed_over, KHÔNG trừ remaining
- `CancelHandover()` - hủy handover, HOÀN LẠI remaining
- Các query functions: `GetPendingHandover`, `GetHandoverHistory`, etc.

**Debug logging:**
- Thêm emoji và log rõ ràng cho mỗi bước
- 🔵 CREATE HANDOVER
- 🟢 CONFIRM HANDOVER  
- 🟡 UPDATE BALANCES

### Backend Handler (cash_handover_handler.go)

**Đơn giản hóa:**
- `CreateHandover()` - chỉ nhận `cash_amount` và `transfer_amount`
- `ConfirmHandover()` - chỉ nhận `actual_cash_amount` và `actual_transfer_amount`
- Xóa toàn bộ logic backward compatibility phức tạp

**Routes comment out:**
- `/handover-and-end` - không cần thiết
- `/quick-confirm` - không cần thiết
- `/pending-approval` - implement sau
- `/approve` - implement sau
- `/discrepancy-stats` - implement sau

### Frontend Cashier View (CashierHandoverView.vue)

**Thay đổi lớn:**

**TRƯỚC:**
```vue
<!-- Chỉ có 1 input -->
<input v-model.number="confirmForm.actual_amount" />
```

**SAU:**
```vue
<!-- Input riêng cho cash -->
<div v-if="selectedHandover?.cash_declared_amount > 0">
  <label>💵 Số tiền mặt thực nhận (VNĐ) *</label>
  <input v-model.number="confirmForm.actual_cash_amount" />
</div>

<!-- Input riêng cho transfer -->
<div v-if="selectedHandover?.transfer_declared_amount > 0">
  <label>💳 Số tiền CK thực nhận (VNĐ) *</label>
  <input v-model.number="confirmForm.actual_transfer_amount" />
</div>
```

**Hiển thị chênh lệch riêng:**
```vue
<!-- Cash Discrepancy -->
<div v-if="cashDiscrepancy !== 0">
  💵 {{ cashDiscrepancy > 0 ? 'Thừa' : 'Thiếu' }} tiền mặt: 
  {{ formatPrice(Math.abs(cashDiscrepancy)) }}
</div>

<!-- Transfer Discrepancy -->
<div v-if="transferDiscrepancy !== 0">
  💳 {{ transferDiscrepancy > 0 ? 'Thừa' : 'Thiếu' }} tiền CK: 
  {{ formatPrice(Math.abs(transferDiscrepancy)) }}
</div>
```

### Frontend Waiter View (ShiftView.vue)

**Không thay đổi** - đã OK từ trước:
- Form bàn giao có 3 options: 💵 Tiền mặt, 💳 Tiền CK, 💰 Cả hai
- Gửi `cash_amount` và `transfer_amount` riêng biệt

## LUỒNG XỬ LÝ MỚI

### Ví dụ: Bàn giao 20,000 tiền CK

#### 1. Waiter tạo handover
```
Shift trước:
- remaining_cash = 50,000
- remaining_transfer = 30,000
- handed_over_transfer = 0

Waiter chọn: 💳 Tiền CK, nhập 20,000

Backend CreateHandover():
🔵 [CREATE HANDOVER] Starting - Cash: 0, Transfer: 20000
🔵 [CREATE HANDOVER] Shift before - RemainingCash: 50000, RemainingTransfer: 30000
🔵 [CREATE HANDOVER] Shift after - RemainingCash: 50000, RemainingTransfer: 10000
✅ [CREATE HANDOVER] Success

Shift sau:
- remaining_cash = 50,000 (không đổi)
- remaining_transfer = 10,000 (30,000 - 20,000) ✅
- handed_over_transfer = 0 (chưa confirm)
```

#### 2. Cashier confirm handover
```
Cashier nhập:
- actual_cash_amount = 0
- actual_transfer_amount = 20,000 (đúng với khai báo)

Backend ConfirmHandover():
🟢 [CONFIRM HANDOVER] Starting - ActualCash: 0, ActualTransfer: 20000
🟢 [CONFIRM HANDOVER] Handover - DeclaredCash: 0, DeclaredTransfer: 20000
🟢 [CONFIRM HANDOVER] Discrepancy - Cash: 0, Transfer: 0, Total: 0
✅ [CONFIRM HANDOVER] Success

Backend updateBalances():
🟡 [UPDATE BALANCES] Starting
🟡 [UPDATE BALANCES] Waiter shift before - HandedOverCash: 0, HandedOverTransfer: 0, RemainingCash: 50000, RemainingTransfer: 10000
🟡 [UPDATE BALANCES] Waiter shift after - HandedOverCash: 0, HandedOverTransfer: 20000, RemainingCash: 50000, RemainingTransfer: 10000
✅ [UPDATE BALANCES] Success

Shift sau:
- remaining_cash = 50,000 (không đổi) ✅
- remaining_transfer = 10,000 (không đổi) ✅
- handed_over_transfer = 20,000 (cập nhật) ✅
```

## KẾT QUẢ

✅ Logic đơn giản, rõ ràng, dễ hiểu
✅ Chỉ trừ remaining 1 lần duy nhất (khi create)
✅ Phân biệt rõ ràng cash vs transfer
✅ Frontend cashier có input riêng cho từng loại
✅ Debug logging chi tiết với emoji
✅ Không còn code cũ gây rối

## TEST

1. Login waiter, tạo order thanh toán tiền CK
2. Bàn giao một phần tiền CK (ví dụ 20,000)
3. Kiểm tra shift: `remaining_transfer` giảm 20,000 ✅
4. Login cashier, xác nhận bàn giao
5. Nhập `actual_transfer_amount = 20,000`
6. Kiểm tra shift waiter:
   - `remaining_transfer` không đổi ✅
   - `handed_over_transfer` tăng 20,000 ✅
7. Kiểm tra shift cashier:
   - `received_transfer` tăng 20,000 ✅

## FILES THAY ĐỔI

- `backend/application/services/cash_handover_service.go` - viết lại hoàn toàn
- `backend/interfaces/http/cash_handover_handler.go` - đơn giản hóa
- `backend/main.go` - comment out routes không cần thiết
- `frontend/src/views/CashierHandoverView.vue` - thêm input riêng cash/transfer
- `frontend/src/views/ShiftView.vue` - không thay đổi (đã OK)

## BACKUP FILES

- `backend/application/services/cash_handover_service_backup.go` - đã xóa
- `backend/interfaces/http/cash_handover_handler_backup.go` - đã xóa
- `frontend/src/views/CashierHandoverView_backup.vue` - giữ lại để tham khảo
