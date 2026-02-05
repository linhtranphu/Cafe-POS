# ✅ Constants Migration Checklist

**Ngày bắt đầu:** 2026-02-04  
**Mục tiêu:** Migrate tất cả hardcoded strings sang constants pattern  
**Tổng số files:** 10 files cần migration

---

## 📊 Tiến độ tổng thể

- **Hoàn thành:** 1/10 (10%)
- **Đang làm:** 0/10 (0%)
- **Chưa bắt đầu:** 9/10 (90%)

```
Progress: [██░░░░░░░░░░░░░░░░░░] 10%
```

---

## 🔴 Phase 1: CRITICAL - Stores (Priority 1)

### ✅ 1. frontend/src/stores/cashier.js
- [ ] Import `PAYMENT_METHOD` from constants
- [ ] Replace `'CASH'` → `PAYMENT_METHOD.CASH` (4 locations)
- [ ] Replace `'TRANSFER'` → `PAYMENT_METHOD.TRANSFER` (2 locations)
- [ ] Replace `'QR'` → `PAYMENT_METHOD.QR` (2 locations)
- [ ] Test: Verify payment filtering works
- [ ] Test: Verify cash amount calculation works

**Locations:**
- Line 30: `cashPayments` getter
- Line 33: `transferPayments` getter
- Line 36: `qrPayments` getter
- Line 40: `totalCashAmount` getter

---

### ✅ 2. frontend/src/stores/shift.js
- [ ] Import `SHIFT_STATUS` from constants
- [ ] Replace `'OPEN'` → `SHIFT_STATUS.OPEN` (2 locations)
- [ ] Replace `'CLOSED'` → `SHIFT_STATUS.CLOSED` (1 location)
- [ ] Test: Verify `hasOpenShift` computed works
- [ ] Test: Verify `openShifts` filter works
- [ ] Test: Verify `closedShifts` filter works

**Locations:**
- Line 33: `hasOpenShift` getter
- Line 41: `openShifts` getter
- Line 49: `closedShifts` getter

---

### ✅ 3. frontend/src/stores/cashierShift.js
- [ ] Import `CASHIER_SHIFT_STATUS` from constants
- [ ] Replace `'OPEN'` → `CASHIER_SHIFT_STATUS.OPEN` (1 location)
- [ ] Replace `'CLOSED'` → `CASHIER_SHIFT_STATUS.CLOSED` (1 location)
- [ ] Replace `'CLOSURE_INITIATED'` → `CASHIER_SHIFT_STATUS.CLOSURE_INITIATED` (2 locations)
- [ ] Test: Verify `hasOpenCashierShift` works
- [ ] Test: Verify `canStartCashierShift` works
- [ ] Test: Verify `isClosureInitiated` works
- [ ] Test: Verify `isClosed` works

**Locations:**
- Line 36: `hasOpenCashierShift` getter
- Line 44: `canStartCashierShift` getter
- Line 52: `isClosureInitiated` getter
- Line 60: `isClosed` getter

---

### ✅ 4. frontend/src/stores/order.js
- [ ] Import `ORDER_STATUS` from constants
- [ ] Replace `'CREATED'` → `ORDER_STATUS.CREATED` (1 location)
- [ ] Replace `'PAID'` → `ORDER_STATUS.PAID` (1 location)
- [ ] Replace `'IN_PROGRESS'` → `ORDER_STATUS.IN_PROGRESS` (1 location)
- [ ] Replace `'SERVED'` → `ORDER_STATUS.SERVED` (1 location)
- [ ] Test: Verify order filtering by status works
- [ ] Test: Verify order counts are correct

**Locations:**
- Line 145: `createdOrders` getter
- Line 149: `paidOrders` getter
- Line 153: `inProgressOrders` getter
- Line 157: `servedOrders` getter

---

### ✅ 5. frontend/src/stores/barista.js
- [ ] Import `ORDER_STATUS` from constants
- [ ] Replace `'IN_PROGRESS'` → `ORDER_STATUS.IN_PROGRESS` (2 locations)
- [ ] Replace `'READY'` → `ORDER_STATUS.READY` (2 locations)
- [ ] Replace `'SERVED'` → `ORDER_STATUS.SERVED` (2 locations)
- [ ] Test: Verify barista order filtering works
- [ ] Test: Verify order counts are correct
- [ ] Test: Verify queue management works

**Locations:**
- Line 89: `inProgressOrders` getter
- Line 93: `readyOrders` getter
- Line 97: `servedOrders` getter
- Line 105: `inProgressCount` getter
- Line 109: `readyCount` getter
- Line 113: `servedCount` getter

---

## 🟠 Phase 2: HIGH - Vue Components (Priority 2)

### ✅ 6. frontend/src/views/CashierShiftClosure.vue
- [ ] Import `CASHIER_SHIFT_STATUS` from constants
- [ ] Replace template conditionals (3 locations)
- [ ] Replace computed properties (2 locations)
- [ ] Replace `getStatusText()` function (1 location)
- [ ] Test: Verify shift closure flow works
- [ ] Test: Verify status display is correct
- [ ] Test: Verify step transitions work

**Locations:**
- Line 66: Template `v-if="shift.status === 'OPEN'"`
- Line 81: Template `v-if="shift.status === 'CLOSURE_INITIATED'"`
- Line 228: Template `v-if="shift.status === 'CLOSED'"`
- Line 268: `canConfirm` computed
- Line 282: `canCloseShift` computed
- Lines 422-424: `getStatusText()` function

---

### ✅ 7. frontend/src/views/ManagerShiftView.vue
- [ ] Import `SHIFT_STATUS`, `CASHIER_SHIFT_STATUS`, `ROLE_TYPE` from constants
- [ ] Replace filter buttons (2 locations)
- [ ] Replace computed properties (3 locations)
- [ ] Replace utility functions (2 locations)
- [ ] Test: Verify shift filtering works
- [ ] Test: Verify role type filtering works
- [ ] Test: Verify status display is correct

**Locations:**
- Line 21: Filter button `'OPEN'`
- Line 26: Filter button `'CLOSED'`
- Line 396: `openWaiterShifts` computed
- Line 401: `openBaristaShifts` computed
- Line 406: `openCashierShifts` computed
- Line 466: `getStatusClass()` function
- Line 472: `getStatusText()` function

---

### ✅ 8. frontend/src/views/ShiftView.vue
- [ ] Import `SHIFT_STATUS` from constants
- [ ] Replace template conditionals (6 locations)
- [ ] Test: Verify shift display works
- [ ] Test: Verify status badges are correct
- [ ] Test: Verify conditional rendering works

**Locations:**
- Line 186: Status badge class binding
- Line 188: Status text display
- Line 197: `v-if="shift.status === 'CLOSED'"`
- Line 201: `v-if="shift.status === 'CLOSED'"`
- Line 205: `v-if="shift.status === 'CLOSED'"`
- Line 211: `v-if="shift.status === 'OPEN'"`

---

### ✅ 9. frontend/src/components/CashierShiftManager.vue
- [ ] Import `CASHIER_SHIFT_STATUS` from constants
- [ ] Replace `canCloseShift` computed (1 location)
- [ ] Test: Verify shift manager component works
- [ ] Test: Verify close shift button state is correct

**Locations:**
- Line 136: `canCloseShift` computed

---

## 🟡 Phase 3: MEDIUM - Vue Views (Priority 3)

### ✅ 10. frontend/src/views/DashboardView.vue
- [ ] Import `ORDER_STATUS` from constants
- [ ] Replace `todayRevenue` computed (1 location)
- [ ] Replace `completedOrders` computed (1 location)
- [ ] Replace `pendingOrders` computed (3 locations)
- [ ] Test: Verify dashboard statistics are correct
- [ ] Test: Verify revenue calculation works
- [ ] Test: Verify order counts are correct

**Locations:**
- Line 523: `todayRevenue` - filter by `!== 'CANCELLED'`
- Line 530: `completedOrders` - filter by `=== 'SERVED'`
- Line 538: `pendingOrders` - filter by `!== 'SERVED'`
- Line 542: `pendingOrders` - filter by `!== 'CANCELLED'` and `=== 'CREATED'`

---

## 🧪 Phase 4: Testing & Verification

### Unit Tests
- [ ] Add tests for constant values
- [ ] Verify constants match backend
- [ ] Test constant imports work

### Integration Tests
- [ ] Test order status transitions
- [ ] Test shift status transitions
- [ ] Test payment method filtering
- [ ] Test role type filtering

### Manual Testing
- [ ] Test all order flows
- [ ] Test all shift flows
- [ ] Test cashier dashboard
- [ ] Test manager dashboard
- [ ] Test barista view
- [ ] Verify no console errors
- [ ] Verify API calls work correctly

---

## 📝 Migration Template

Use this template for each file:

```markdown
### File: [filename]

**Status:** 🔴 Not Started / 🟡 In Progress / 🟢 Complete

**Changes:**
1. Import constants
2. Replace hardcoded strings
3. Test functionality

**Testing:**
- [ ] Unit tests pass
- [ ] Manual testing complete
- [ ] No console errors
- [ ] API calls work

**Notes:**
[Any issues or observations]
```

---

## 🚀 Quick Start Guide

### For Each File:

1. **Open the file**
2. **Add imports at top:**
   ```javascript
   import { ORDER_STATUS, PAYMENT_METHOD } from '../constants/order'
   import { SHIFT_STATUS, CASHIER_SHIFT_STATUS, ROLE_TYPE } from '../constants/shift'
   ```

3. **Find & Replace:**
   - Use IDE search: `'OPEN'` → `SHIFT_STATUS.OPEN`
   - Use IDE search: `'CASH'` → `PAYMENT_METHOD.CASH`
   - etc.

4. **Update object keys:**
   ```javascript
   // Before
   const map = { 'OPEN': 'text' }
   
   // After
   const map = { [SHIFT_STATUS.OPEN]: 'text' }
   ```

5. **Test thoroughly**

6. **Mark as complete** in this checklist

---

## 🎯 Success Criteria

Migration is complete when:
- [ ] All 10 files updated
- [ ] All tests passing
- [ ] No hardcoded strings remain
- [ ] No console errors
- [ ] All features work correctly
- [ ] Code review approved
- [ ] Documentation updated

---

## 📞 Support

**Questions?** Check:
1. `docs/CONSTANTS_PATTERN.md` - Full pattern guide
2. `docs/CONSTANTS_QUICK_REFERENCE.md` - Quick reference
3. `docs/CONSTANTS_VIOLATIONS_AUDIT.md` - Detailed violations
4. `frontend/src/views/CashierDashboard.vue` - Reference implementation ✅

---

**Last Updated:** 2026-02-04  
**Next Review:** After Phase 1 completion  
**Estimated Time:** 4-6 hours total
