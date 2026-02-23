# Test: Transfer Handover Flow

## Mục Đích

Test end-to-end flow bàn giao tiền chuyển khoản để đảm bảo:
1. ✅ Payment chuyển khoản cập nhật `transfer_revenue`, không ảnh hưởng `current_cash`
2. ✅ Handover tiền CK trừ vào `remaining_transfer`, không trừ vào `remaining_cash`
3. ✅ Tiền mặt của waiter không bị âm sau khi bàn giao tiền CK
4. ✅ Frontend hiển thị đúng số tiền

## Prerequisites

1. Backend đang chạy tại `http://localhost:8080`
2. Frontend đang chạy tại `http://localhost:5173`
3. MongoDB đang chạy trong Docker
4. Có users: `waiter` và `cashier` với password `password123`
5. Có ít nhất 1 menu item

## Test Scripts

### 1. API Test Script

**File:** `test-transfer-handover-flow.sh`

**Chức năng:**
- Tạo order với payment TRANSFER
- Bàn giao tiền CK
- Cashier xác nhận
- Verify shift state

**Chạy:**
```bash
./test-transfer-handover-flow.sh
```

**Expected Output:**
```
===================================
🧪 TEST: Transfer Handover Flow
===================================

Step 1: Login as waiter
✅ Logged in as waiter

Step 2: Get waiter's current shift
✅ Found shift: 699bda2ea31585ce7ad4c47c

📊 Shift BEFORE:
   Current Cash: 72000 VND
   Remaining Cash: 72000 VND
   Remaining Transfer: 0 VND

Step 3: Create order with TRANSFER payment
✅ Created order: 699bdb1ea31585ce7ad4c485
   Total: 30000 VND

Step 4: Pay order with TRANSFER
✅ Payment successful (TRANSFER)

📊 Shift AFTER PAYMENT:
   Current Cash: 72000 VND
   Remaining Cash: 72000 VND
   Remaining Transfer: 30000 VND

✅ Payment correctly updated transfer, not cash

Step 5: Create handover for TRANSFER
✅ Created handover: 699bdb2fa31585ce7ad4c486
   Transfer Amount: 30000 VND

Step 6: Login as cashier
✅ Logged in as cashier

Step 7: Confirm handover
✅ Handover confirmed

Step 8: Verify final shift state

📊 Shift FINAL STATE:
   Current Cash: 72000 VND
   Remaining Cash: 72000 VND
   Remaining Transfer: 0 VND
   Handed Over Cash: 0 VND
   Handed Over Transfer: 30000 VND

Step 9: Verify results

✅ PASS: Cash unchanged (72000 VND)
✅ PASS: Transfer fully handed over (0 VND)
✅ PASS: No cash handed over (0 VND)
✅ PASS: Transfer handed over correctly (30000 VND)

===================================
✅ ALL TESTS PASSED!
===================================
```

### 2. Database Verification Script

**File:** `verify-transfer-handover-db.js`

**Chức năng:**
- Kiểm tra handover record trong DB
- Kiểm tra shift record trong DB
- Verify tất cả amounts đúng

**Chạy:**
```bash
docker cp verify-transfer-handover-db.js cafe-pos-mongodb:/tmp/
docker exec cafe-pos-mongodb mongosh cafe_pos -u admin -p password123 \
  --authenticationDatabase admin /tmp/verify-transfer-handover-db.js
```

**Expected Output:**
```
=================================
🔍 VERIFY: Transfer Handover Data
=================================

📋 Most Recent Transfer Handover:
   ID: 699bdb2fa31585ce7ad4c486
   Status: CONFIRMED
   Type: PARTIAL

💰 Declared Amounts:
   Cash: 0 VND
   Transfer: 30000 VND
   Total (old): 30000 VND

✅ Actual Amounts:
   Cash: 0 VND
   Transfer: 30000 VND
   Total (old): 30000 VND

✅ PASS: Cash declared is 0
✅ PASS: Transfer declared is 30000
✅ PASS: Cash actual is 0
✅ PASS: Transfer actual matches declared

📊 Checking Shift Data:
   Shift ID: 699bda2ea31585ce7ad4c47c
   User: waiter
   Status: OPEN

💵 Cash Amounts:
   Current Cash: 72000 VND
   Remaining Cash: 72000 VND
   Handed Over Cash: 0 VND

💳 Transfer Amounts:
   Transfer Revenue: 30000 VND
   Remaining Transfer: 0 VND
   Handed Over Transfer: 30000 VND

✅ PASS: No cash handed over (0 VND)
✅ PASS: Handed over transfer matches handovers (30000 VND)
✅ PASS: Remaining transfer is correct (0 VND)
✅ PASS: Current cash is not negative (72000 VND)
✅ PASS: Remaining cash is not negative (72000 VND)

=================================
✅ ALL CHECKS PASSED!
=================================
```

## Manual Frontend Testing

### 1. Waiter View

**URL:** `http://localhost:5173/#/shifts`

**Login:** `waiter` / `password123`

**Verify:**
```
Ca đang mở
☀️ Ca sáng
👨‍💼 Phục vụ
ĐANG MỞ

Bắt đầu: 11:40
Tiền đầu ca: 50.000 ₫

💵 Tiền mặt: 72.000 ₫        ← Should NOT be negative!
💳 Tiền CK: 0 ₫              ← Should be 0 after handover
Đã bàn giao: 30.000 ₫        ← Should show transfer amount
```

**Expected:**
- ✅ Tiền mặt không âm
- ✅ Tiền CK = 0 (đã bàn giao hết)
- ✅ Đã bàn giao = số tiền CK đã bàn giao

### 2. Cashier View

**URL:** `http://localhost:5173/#/cashier/handovers`

**Login:** `cashier` / `password123`

**Verify:**
```
📋 Bàn giao hôm nay

┌─────────────────────────────────────┐
│ waiter                              │
│ 11:45                               │
│ [Một phần] [💳 CK]                  │  ← Should show transfer badge
│                                     │
│                     💳 30.000 ₫     │  ← Should show transfer amount
│                     [Đã xác nhận]   │
└─────────────────────────────────────┘
```

**Expected:**
- ✅ Badge hiển thị "💳 CK" (không phải "💵 Mặt")
- ✅ Số tiền hiển thị với icon 💳
- ✅ Status "Đã xác nhận"

## Test Cases

### Test Case 1: Transfer Payment Updates Transfer Revenue

**Steps:**
1. Create order
2. Pay with TRANSFER method
3. Check shift

**Expected:**
- `current_cash` không đổi
- `remaining_cash` không đổi
- `transfer_revenue` tăng
- `remaining_transfer` tăng

### Test Case 2: Transfer Handover Deducts Transfer

**Steps:**
1. Create handover with `transfer_amount` > 0, `cash_amount` = 0
2. Cashier confirms
3. Check shift

**Expected:**
- `remaining_cash` không đổi
- `remaining_transfer` giảm
- `handed_over_cash` không đổi
- `handed_over_transfer` tăng

### Test Case 3: Mixed Handover (Cash + Transfer)

**Steps:**
1. Create handover with both `cash_amount` and `transfer_amount`
2. Cashier confirms
3. Check shift

**Expected:**
- `remaining_cash` giảm theo `cash_amount`
- `remaining_transfer` giảm theo `transfer_amount`
- `handed_over_cash` tăng theo `cash_actual_amount`
- `handed_over_transfer` tăng theo `transfer_actual_amount`

### Test Case 4: Old Format Backward Compatibility

**Steps:**
1. Create handover with only `declared_amount` (old format)
2. Cashier confirms with only `actual_amount`
3. Check shift

**Expected:**
- Backend auto-converts based on declared amounts
- Shift updates correctly

## Troubleshooting

### Issue: Cash becomes negative

**Symptoms:**
```
💵 Tiền mặt: -8.000 ₫
```

**Cause:**
- Backend cũ chạy khi confirm handover
- `actual_amount` không được convert sang `transfer_actual_amount`
- Shift bị trừ vào `remaining_cash` thay vì `remaining_transfer`

**Fix:**
1. Restart backend để apply code mới
2. Chạy script fix DB: `fix-shift-699bda2e.js`

### Issue: Transfer not deducted

**Symptoms:**
```
💳 Tiền CK: 30.000 ₫  (should be 0)
```

**Cause:**
- Handover không có `transfer_actual_amount`
- Backend không update `remaining_transfer`

**Fix:**
1. Check handover record trong DB
2. Verify `transfer_actual_amount` có giá trị
3. Nếu không, chạy script fix

### Issue: Wrong badge in cashier view

**Symptoms:**
- Hiển thị "💵 Mặt" thay vì "💳 CK"

**Cause:**
- Frontend logic xác định payment type sai
- Handover có `declared_amount` nhưng không có `transfer_declared_amount`

**Fix:**
- Check frontend logic trong `CashierHandoverView.vue`
- Verify handover có `transfer_declared_amount` > 0

## Success Criteria

✅ All API tests pass
✅ All DB verification checks pass
✅ Frontend displays correct amounts
✅ No negative cash values
✅ Correct badges and colors in UI
✅ Backward compatibility maintained

## Related Files

- `backend/application/services/cash_handover_service.go` - Handover logic
- `backend/application/services/order_service.go` - Payment logic
- `frontend/src/views/CashierHandoverView.vue` - Cashier UI
- `frontend/src/views/ShiftView.vue` - Waiter shift UI
- `TRANSFER_HANDOVER_COMPLETE_FIX.md` - Fix documentation
- `CASHIER_WARNING_LOGIC_FIX.md` - Warning logic fix
- `SHIFT_NEGATIVE_CASH_FIX.md` - Negative cash fix

## Next Steps

1. ✅ Run API test
2. ✅ Run DB verification
3. ✅ Manual frontend testing
4. ⏳ User acceptance testing
5. ⏳ Deploy to production
