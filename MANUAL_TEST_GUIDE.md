# Manual End-to-End Test Guide

## Mục Tiêu Test

Kiểm tra flow hoàn chỉnh:
1. ✅ Order với tiền mặt → Tiền mặt tăng
2. ✅ Order với chuyển khoản → Tiền CK tăng
3. ✅ Handover tiền CK → Tiền CK giảm, tiền mặt KHÔNG đổi
4. ✅ Cashier confirm → Verify không có tiền âm

## Prerequisites

- ✅ Backend running: http://localhost:3000
- ✅ Frontend running: http://localhost:5173
- ✅ MongoDB running with replica set
- ✅ Users: waiter/waiter123, cashier/cashier123

## Test Steps

### PHASE 1: Setup & Initial State

#### Step 1.1: Login as Waiter
```
URL: http://localhost:5173
Username: waiter
Password: waiter123
```

#### Step 1.2: Check Current Shift
```
Navigate to: Shifts (/#/shifts)

Expected:
- Shift status: ĐANG MỞ
- Start cash: 12,000 ₫ (or whatever was set)
- Current cash: 72,000 ₫
- Remaining cash: 72,000 ₫
- Transfer revenue: 0 ₫
- Remaining transfer: 0 ₫
```

**📸 Screenshot 1:** Initial shift state

**Record Initial Values:**
```
Initial State:
- Current Cash: _________ VND
- Remaining Cash: _________ VND
- Transfer Revenue: _________ VND
- Remaining Transfer: _________ VND
```

---

### PHASE 2: Create Orders

#### Step 2.1: Create Order #1 - CASH Payment
```
1. Navigate to: Dashboard (/#/)
2. Click "Tạo đơn mới"
3. Select menu item (any item)
4. Enter customer name: "Test Cash"
5. Click "Tạo đơn"
6. Note order total: _________ VND
7. Click "Thanh toán"
8. Select payment method: TIỀN MẶT
9. Enter amount: [order total]
10. Click "Xác nhận thanh toán"
```

**Expected Result:**
```
✅ Order created successfully
✅ Payment successful
✅ Order status: PAID
```

#### Step 2.2: Verify Shift After Cash Payment
```
Navigate to: Shifts (/#/shifts)

Expected Changes:
- Current Cash: [Initial + Order Total] ✅
- Remaining Cash: [Initial + Order Total] ✅
- Transfer Revenue: [No change] ✅
- Remaining Transfer: [No change] ✅
```

**📸 Screenshot 2:** Shift after cash payment

**Record After Cash Payment:**
```
After Cash Payment (Order: _____ VND):
- Current Cash: _________ VND (should increase)
- Remaining Cash: _________ VND (should increase)
- Transfer Revenue: _________ VND (should NOT change)
- Remaining Transfer: _________ VND (should NOT change)
```

#### Step 2.3: Create Order #2 - TRANSFER Payment
```
1. Navigate to: Dashboard (/#/)
2. Click "Tạo đơn mới"
3. Select menu item (any item)
4. Enter customer name: "Test Transfer"
5. Click "Tạo đơn"
6. Note order total: _________ VND
7. Click "Thanh toán"
8. Select payment method: CHUYỂN KHOẢN
9. Enter amount: [order total]
10. Click "Xác nhận thanh toán"
```

**Expected Result:**
```
✅ Order created successfully
✅ Payment successful
✅ Order status: PAID
```

#### Step 2.4: Verify Shift After Transfer Payment
```
Navigate to: Shifts (/#/shifts)

Expected Changes:
- Current Cash: [No change from Step 2.2] ✅
- Remaining Cash: [No change from Step 2.2] ✅
- Transfer Revenue: [Order Total] ✅
- Remaining Transfer: [Order Total] ✅
```

**📸 Screenshot 3:** Shift after transfer payment

**Record After Transfer Payment:**
```
After Transfer Payment (Order: _____ VND):
- Current Cash: _________ VND (should NOT change)
- Remaining Cash: _________ VND (should NOT change)
- Transfer Revenue: _________ VND (should increase)
- Remaining Transfer: _________ VND (should increase)
```

**🔍 CRITICAL CHECK #1:**
```
✅ Cash did NOT increase after transfer payment
✅ Transfer revenue DID increase
```

---

### PHASE 3: Handover Transfer

#### Step 3.1: Create Handover for Transfer
```
1. Stay in Shifts view (/#/shifts)
2. Scroll to "Bàn giao tiền" section
3. Select handover type: "Một phần"
4. Enter cash amount: 0
5. Enter transfer amount: [Transfer Revenue from Step 2.4]
6. Enter note: "Test transfer handover"
7. Click "Tạo yêu cầu bàn giao"
```

**Expected Result:**
```
✅ Handover created successfully
✅ Status: Chờ xác nhận
```

**📸 Screenshot 4:** Handover created

#### Step 3.2: Verify Shift After Creating Handover
```
Navigate to: Shifts (/#/shifts)

Expected Changes:
- Current Cash: [No change] ✅
- Remaining Cash: [No change] ✅
- Transfer Revenue: [No change] ✅
- Remaining Transfer: [Should be 0] ✅
- Handed Over: [Transfer amount] ✅
```

**Record After Creating Handover:**
```
After Creating Handover (Amount: _____ VND):
- Current Cash: _________ VND (should NOT change)
- Remaining Cash: _________ VND (should NOT change)
- Transfer Revenue: _________ VND (should NOT change)
- Remaining Transfer: _________ VND (should be 0)
- Handed Over: _________ VND (should show transfer amount)
```

**🔍 CRITICAL CHECK #2:**
```
✅ Cash did NOT decrease after handover
✅ Transfer remaining is now 0
✅ No negative values
```

---

### PHASE 4: Cashier Confirmation

#### Step 4.1: Logout and Login as Cashier
```
1. Click profile icon → Logout
2. Login with:
   Username: cashier
   Password: cashier123
```

#### Step 4.2: View Pending Handovers
```
Navigate to: Cashier Handovers (/#/cashier/handovers)

Expected:
- See pending handover from waiter
- Badge shows: 💳 Chuyển khoản (BLUE badge)
- Amount shows: 💳 [amount] ₫ (with blue icon)
- NOT showing: 💵 Tiền mặt (green badge)
```

**📸 Screenshot 5:** Cashier view - pending handover

**Verify Badge:**
```
✅ Badge color: Blue (not green)
✅ Badge text: "💳 Chuyển khoản" or "💳 CK"
✅ Amount icon: 💳 (not 💵)
```

#### Step 4.3: Confirm Handover
```
1. Click "Xác nhận" on the pending handover
2. Enter actual amount: [same as declared amount]
3. Enter note: "Test confirmation"
4. Click "Xác nhận"
```

**Expected Result:**
```
✅ Handover confirmed successfully
✅ Status changed to: Đã xác nhận
✅ Moved to "Bàn giao hôm nay" section
```

**📸 Screenshot 6:** Handover confirmed

---

### PHASE 5: Final Verification

#### Step 5.1: Logout and Login as Waiter Again
```
1. Logout from cashier
2. Login as waiter/waiter123
```

#### Step 5.2: Check Final Shift State
```
Navigate to: Shifts (/#/shifts)

Expected Final State:
- Current Cash: [Same as Step 2.2] ✅
- Remaining Cash: [Same as Step 2.2] ✅
- Transfer Revenue: [Same as Step 2.4] ✅
- Remaining Transfer: 0 ✅
- Handed Over Transfer: [Transfer amount] ✅
- Handed Over Cash: 0 ✅
```

**📸 Screenshot 7:** Final shift state

**Record Final State:**
```
Final State:
- Current Cash: _________ VND
- Remaining Cash: _________ VND
- Transfer Revenue: _________ VND
- Remaining Transfer: _________ VND (should be 0)
- Handed Over Transfer: _________ VND
- Handed Over Cash: _________ VND (should be 0)
```

**🔍 CRITICAL CHECK #3:**
```
✅ Current Cash = [Initial + Cash Order Total]
✅ Remaining Cash = [Initial + Cash Order Total]
✅ NO negative values anywhere
✅ Transfer fully handed over (remaining = 0)
✅ Cash NOT affected by transfer handover
```

---

## Verification Checklist

### ✅ Payment Verification
- [ ] Cash payment increases current_cash
- [ ] Cash payment increases remaining_cash
- [ ] Transfer payment increases transfer_revenue
- [ ] Transfer payment increases remaining_transfer
- [ ] Transfer payment does NOT affect cash

### ✅ Handover Verification
- [ ] Transfer handover decreases remaining_transfer
- [ ] Transfer handover does NOT affect remaining_cash
- [ ] Transfer handover does NOT affect current_cash
- [ ] Handover shows correct badge (💳 CK)
- [ ] Handover shows correct amount with icon

### ✅ Cashier Verification
- [ ] Cashier sees correct payment type badge
- [ ] Cashier sees correct amount
- [ ] Confirmation updates shift correctly
- [ ] No warnings about wrong payment type

### ✅ Final State Verification
- [ ] No negative cash values
- [ ] No negative transfer values
- [ ] Cash unchanged by transfer operations
- [ ] Transfer correctly handed over
- [ ] All amounts add up correctly

---

## Expected Calculations

### Example Scenario:
```
Initial State:
- Start Cash: 12,000
- Current Cash: 72,000
- Remaining Cash: 72,000
- Transfer Revenue: 0
- Remaining Transfer: 0

Order #1 (Cash): 25,000
After:
- Current Cash: 97,000 (72,000 + 25,000)
- Remaining Cash: 97,000
- Transfer Revenue: 0
- Remaining Transfer: 0

Order #2 (Transfer): 30,000
After:
- Current Cash: 97,000 (NO CHANGE)
- Remaining Cash: 97,000 (NO CHANGE)
- Transfer Revenue: 30,000
- Remaining Transfer: 30,000

Handover Transfer: 30,000
After:
- Current Cash: 97,000 (NO CHANGE)
- Remaining Cash: 97,000 (NO CHANGE)
- Transfer Revenue: 30,000 (NO CHANGE)
- Remaining Transfer: 0 (30,000 - 30,000)
- Handed Over Transfer: 30,000

Final State:
- Current Cash: 97,000 ✅
- Remaining Cash: 97,000 ✅
- Transfer Revenue: 30,000 ✅
- Remaining Transfer: 0 ✅
- Handed Over Transfer: 30,000 ✅
- Handed Over Cash: 0 ✅
```

---

## Troubleshooting

### Issue: Cash becomes negative
**Symptom:** `remaining_cash` shows negative value

**Cause:** Transfer handover incorrectly deducted from cash

**Fix:** Backend not restarted with new code

**Solution:**
```bash
./restart_local.sh
```

### Issue: Wrong badge in cashier view
**Symptom:** Shows "💵 Mặt" instead of "💳 CK"

**Cause:** Frontend logic not updated

**Fix:** Clear browser cache and refresh

### Issue: Warning shows wrong payment type
**Symptom:** Warning says "Tiền mặt" for transfer handover

**Cause:** Frontend warning logic not fixed

**Fix:** Already fixed in code, refresh browser

---

## Success Criteria

All must be TRUE:
- [x] Cash payment updates cash correctly
- [x] Transfer payment updates transfer correctly
- [x] Transfer payment does NOT affect cash
- [x] Transfer handover deducts from transfer
- [x] Transfer handover does NOT affect cash
- [x] Cashier sees correct badge
- [x] No negative values anywhere
- [x] All calculations correct

---

## Test Report Template

```
Test Date: __________
Tester: __________
Environment: Local Development

Initial State:
- Current Cash: _________ VND
- Transfer Revenue: _________ VND

Test Results:
1. Cash Payment: ☐ PASS ☐ FAIL
2. Transfer Payment: ☐ PASS ☐ FAIL
3. Transfer Handover: ☐ PASS ☐ FAIL
4. Cashier Confirmation: ☐ PASS ☐ FAIL
5. Final Verification: ☐ PASS ☐ FAIL

Issues Found:
_________________________________
_________________________________

Overall Result: ☐ PASS ☐ FAIL

Notes:
_________________________________
_________________________________
```

---

**Frontend URLs:**
- Waiter: http://localhost:5173
- Shifts: http://localhost:5173/#/shifts
- Cashier: http://localhost:5173/#/cashier/handovers

**Credentials:**
- Waiter: waiter / waiter123
- Cashier: cashier / cashier123
