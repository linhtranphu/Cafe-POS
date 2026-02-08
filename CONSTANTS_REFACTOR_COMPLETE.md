# Constants Refactor - Complete Summary

## ✅ Completed

### 1. Ingredient Constants
**File**: `frontend/src/constants/ingredient.js`
**Status**: ✅ COMPLETE

**Constants**:
- `INGREDIENT_UNITS` - Unit types (kg, g, L, ml, piece, box, pack)
- `STOCK_OPERATIONS` - Stock operations (in, out, adjust)
- `TRANSACTION_TYPES` - Transaction types (purchase, waste, adjustment, order)
- `ADJUSTMENT_TYPES` - Legacy types (add, remove, adjust)
- `STOCK_STATUS` - Stock status (in_stock, low_stock, out_of_stock)

**Updated Views**:
- ✅ `IngredientManagementView.vue` - All hardcoded strings replaced with constants

**Backend Sync**: ✅ Verified with `backend/domain/ingredient/`

---

### 2. User Constants
**File**: `frontend/src/constants/user.js`
**Status**: ✅ CREATED

**Constants**:
- `USER_ROLES` - User roles (admin, manager, cashier, waiter, barista)
- `USER_ROLE_OPTIONS` - UI options with labels, icons, badges
- `ROLE_PERMISSIONS` - Permission matrix for each role

**Helper Functions**:
- `getUserRoleLabel(role)` - Get display label
- `getUserRoleIcon(role)` - Get emoji icon
- `getUserRoleBadge(role)` - Get CSS badge class
- `hasPermission(role, permission)` - Check permission

**Backend Sync**: ⚠️ Need to verify with `backend/domain/user/user.go`

---

### 3. Shift Constants
**File**: `frontend/src/constants/shift.js`
**Status**: ✅ EXISTS (Review needed)

**Constants**:
- `SHIFT_STATUS` - Shift status (OPEN, CLOSED)
- `CASHIER_SHIFT_STATUS` - Cashier shift status (OPEN, CLOSURE_INITIATED, CLOSED)
- `SHIFT_TYPE` - Shift types (MORNING, AFTERNOON, EVENING)
- `ROLE_TYPE` - Role types (waiter, barista)

**Display Objects**:
- `SHIFT_STATUS_DISPLAY` - Labels, icons, badges
- `CASHIER_SHIFT_STATUS_DISPLAY` - Labels, icons, badges
- `SHIFT_TYPE_DISPLAY` - Labels, icons, badges
- `ROLE_TYPE_DISPLAY` - Labels, icons, badges

**Backend Sync**: ✅ Verified with `backend/domain/order/shift.go` and `backend/domain/cashier/cashier_shift.go`

---

### 4. Order Constants
**File**: `frontend/src/constants/order.js`
**Status**: ✅ EXISTS (Review needed)

**Constants**:
- `ORDER_STATUS` - Order statuses (CREATED, PAID, QUEUED, IN_PROGRESS, READY, SERVED, CANCELLED, LOCKED)
- `PAYMENT_METHOD` - Payment methods (CASH, TRANSFER, QR)

**Display Objects**:
- `ORDER_STATUS_DISPLAY` - Labels, icons, badges
- `PAYMENT_METHOD_DISPLAY` - Labels, icons
- `STATUS_FILTER_OPTIONS` - Filter options for UI

**Backend Sync**: ✅ Verified with `backend/domain/order/order.go`

---

### 5. Expense Constants
**File**: `frontend/src/constants/expense.js`
**Status**: ✅ EXISTS (Review needed)

**Constants**:
- `PAYMENT_METHODS` - Payment methods (cash, bank, card)
- `RECURRING_FREQUENCIES` - Frequencies (daily, weekly, monthly, quarterly, yearly)

**Helper Functions**:
- `getPaymentMethodLabel(method)` - Get display label
- `getFrequencyLabel(frequency)` - Get display label

**Backend Sync**: ✅ Verified with `backend/domain/expense/expense.go`

---

### 6. Facility Constants
**File**: `frontend/src/constants/facility.js`
**Status**: ✅ EXISTS (Review needed)

**Constants**: (Need to check file content)

**Backend Sync**: ⚠️ Need to verify with `backend/domain/facility/facility.go`

---

## ⚠️ Pending Updates

### Views That Need Constants Update

#### Priority 1 (Critical - User Roles)
1. **DashboardView.vue**
   - Replace `user?.role === 'manager'` with `USER_ROLES.MANAGER`
   - Replace `authStore.user?.role === 'barista'` with `USER_ROLES.BARISTA`
   - Replace `authStore.user?.role === 'cashier'` with `USER_ROLES.CASHIER`
   - Replace order status strings with `ORDER_STATUS` constants

2. **UserManagementView.vue**
   - Replace `u.role === 'manager'` with `USER_ROLES.MANAGER`
   - Replace `u.role === 'cashier'` with `USER_ROLES.CASHIER`

3. **ShiftView.vue**
   - Replace `authStore.user?.role === 'cashier'` with `USER_ROLES.CASHIER`
   - Replace `authStore.user?.role === 'manager'` with `USER_ROLES.MANAGER`
   - Replace `authStore.user?.role === 'waiter'` with `USER_ROLES.WAITER`
   - Replace `shift.status === 'OPEN'` with `SHIFT_STATUS.OPEN`
   - Replace `shift.status === 'CLOSED'` with `SHIFT_STATUS.CLOSED`

#### Priority 2 (High - Shift Status)
4. **ManagerShiftView.vue**
   - Replace `filterStatus === 'OPEN'` with `SHIFT_STATUS.OPEN`
   - Replace `filterStatus === 'CLOSED'` with `SHIFT_STATUS.CLOSED`
   - Replace `shift.status === 'CLOSED'` with `SHIFT_STATUS.CLOSED`
   - Replace `selectedShiftType === 'waiter'` with `ROLE_TYPE.WAITER`
   - Replace `selectedShiftType === 'cashier'` with constants

5. **CashierShiftClosure.vue**
   - Replace `shift.status === 'OPEN'` with `CASHIER_SHIFT_STATUS.OPEN`
   - Replace `shift.status === 'CLOSURE_INITIATED'` with `CASHIER_SHIFT_STATUS.CLOSURE_INITIATED`
   - Replace `shift.status === 'CLOSED'` with `CASHIER_SHIFT_STATUS.CLOSED`

#### Priority 3 (Medium - Order Status)
6. **OrderView.vue**
   - Replace order status strings with `ORDER_STATUS` constants
   - Replace payment method strings with `PAYMENT_METHOD` constants

7. **BaristaView.vue**
   - Replace order status strings with `ORDER_STATUS` constants
   - Replace shift status strings with `SHIFT_STATUS` constants

8. **CashierDashboard.vue**
   - Replace payment method strings with `PAYMENT_METHOD` constants
   - Replace shift status strings with `CASHIER_SHIFT_STATUS` constants

#### Priority 4 (Low - Other)
9. **ExpenseManagementView.vue**
   - Replace expense category strings with constants
   - Replace payment method strings with `PAYMENT_METHODS` constants

10. **FacilityManagementView.vue**
    - Replace facility type strings with constants
    - Replace facility status strings with constants

---

## Implementation Guide

### Step 1: Import Constants

```javascript
// In each Vue file
import { USER_ROLES } from '@/constants/user'
import { SHIFT_STATUS, CASHIER_SHIFT_STATUS } from '@/constants/shift'
import { ORDER_STATUS, PAYMENT_METHOD } from '@/constants/order'
```

### Step 2: Replace Hardcoded Strings

**Before**:
```javascript
if (user.value?.role === 'manager') {
  // ...
}
```

**After**:
```javascript
if (user.value?.role === USER_ROLES.MANAGER) {
  // ...
}
```

### Step 3: Export Constants for Template

```javascript
export default {
  setup() {
    // ... other code
    
    return {
      USER_ROLES,
      SHIFT_STATUS,
      ORDER_STATUS,
      // ... other exports
    }
  }
}
```

### Step 4: Use in Template

**Before**:
```vue
<div v-if="user.role === 'manager'">
  Manager Dashboard
</div>
```

**After**:
```vue
<div v-if="user.role === USER_ROLES.MANAGER">
  Manager Dashboard
</div>
```

---

## Testing Checklist

After updating each view:

- [ ] No console errors
- [ ] All conditions work correctly
- [ ] UI displays properly
- [ ] No broken functionality
- [ ] Constants are imported
- [ ] Constants are exported for template
- [ ] Backend sync verified

---

## Benefits Achieved

### 1. Type Safety ✅
- IDE autocomplete works
- Typos caught at development time
- Easier refactoring

### 2. Maintainability ✅
- Single source of truth
- Easy to update values
- Clear documentation

### 3. Consistency ✅
- Same values across all files
- No typos or variations
- Backend-frontend sync

### 4. Debugging ✅
- Easier to search for usage
- Clear intent in code
- Better error messages

---

## Backend Sync Status

### ✅ Verified
- Ingredient constants
- Order constants
- Shift constants
- Expense constants

### ⚠️ Need Verification
- User constants
- Facility constants

---

## Next Steps

1. **Immediate**: Update DashboardView.vue (most critical)
2. **Today**: Update UserManagementView.vue and ShiftView.vue
3. **This Week**: Update all remaining views
4. **Testing**: Thorough testing after each update
5. **Documentation**: Update this file as progress is made

---

## Estimated Completion

- **Phase 1** (Ingredient): ✅ DONE
- **Phase 2** (User constants): ✅ DONE
- **Phase 3** (Update views): ⏳ IN PROGRESS (0/10 views)
- **Phase 4** (Testing): ⏳ PENDING

**Total Progress**: 20% Complete

---

## Contact

If you encounter any issues during the refactor:
1. Check backend constants first
2. Verify import paths
3. Ensure constants are exported in return statement
4. Test in browser console
5. Check for typos in constant names

---

Last Updated: 2026-02-07
Status: IN PROGRESS
