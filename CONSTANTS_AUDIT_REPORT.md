# Constants Audit Report

## Overview
This report identifies all hardcoded strings in Vue files that should use constants instead.

## Files to Review

### 1. ✅ IngredientManagementView.vue
**Status**: FIXED
- Replaced `'add'`, `'remove'`, `'adjust'` with `ADJUSTMENT_TYPES.ADD`, `ADJUSTMENT_TYPES.REMOVE`, `ADJUSTMENT_TYPES.ADJUST`

### 2. ⚠️ CashierShiftClosure.vue
**Hardcoded Strings Found**:
- `shift.status === 'OPEN'`
- `shift.status === 'CLOSURE_INITIATED'`
- `shift.status === 'CLOSED'`

**Should Use**: `SHIFT_STATUS` constants from `frontend/src/constants/shift.js`

**Backend Constants**: Check `backend/domain/cashier/cashier_shift.go`

### 3. ⚠️ UserManagementView.vue
**Hardcoded Strings Found**:
- `u.role === 'manager'`
- `u.role === 'cashier'`

**Should Use**: `USER_ROLES` constants

**Backend Constants**: Check `backend/domain/user/user.go`

### 4. ⚠️ ShiftView.vue
**Hardcoded Strings Found**:
- `shift.status === 'OPEN'`
- `shift.status === 'CLOSED'`
- `authStore.user?.role === 'cashier'`
- `authStore.user?.role === 'manager'`
- `authStore.user?.role === 'waiter'`

**Should Use**: 
- `SHIFT_STATUS` constants
- `USER_ROLES` constants

### 5. ⚠️ DashboardView.vue
**Hardcoded Strings Found**:
- `user?.role === 'manager'`
- `authStore.user?.role === 'barista'`
- `authStore.user?.role === 'cashier'`
- `o.status === 'SERVED'`
- `o.status === 'CANCELLED'`
- `o.status === 'CREATED'`

**Should Use**:
- `USER_ROLES` constants
- `ORDER_STATUS` constants from `frontend/src/constants/order.js`

### 6. ⚠️ ManagerShiftView.vue
**Hardcoded Strings Found**:
- `filterStatus === 'all'`
- `filterStatus === 'OPEN'`
- `filterStatus === 'CLOSED'`
- `shift.status === 'CLOSED'`
- `selectedShiftType === 'waiter'`
- `selectedShiftType === 'cashier'`

**Should Use**:
- `SHIFT_STATUS` constants
- `SHIFT_TYPE` constants

### 7. ⚠️ OrderView.vue
**Likely Hardcoded Strings**:
- Order status strings
- Payment method strings

### 8. ⚠️ BaristaView.vue
**Likely Hardcoded Strings**:
- Order status strings
- Shift status strings

### 9. ⚠️ CashierDashboard.vue
**Likely Hardcoded Strings**:
- Payment method strings
- Shift status strings

### 10. ⚠️ ExpenseManagementView.vue
**Likely Hardcoded Strings**:
- Expense category strings
- Payment method strings

### 11. ⚠️ FacilityManagementView.vue
**Likely Hardcoded Strings**:
- Facility type strings
- Facility status strings

## Required Constants Files

### 1. User Constants (`frontend/src/constants/user.js`)
```javascript
// Must match backend: backend/domain/user/user.go
export const USER_ROLES = {
  ADMIN: 'admin',
  MANAGER: 'manager',
  CASHIER: 'cashier',
  WAITER: 'waiter',
  BARISTA: 'barista'
}

export const USER_ROLE_OPTIONS = [
  { value: USER_ROLES.ADMIN, label: 'Admin' },
  { value: USER_ROLES.MANAGER, label: 'Quản lý' },
  { value: USER_ROLES.CASHIER, label: 'Thu ngân' },
  { value: USER_ROLES.WAITER, label: 'Phục vụ' },
  { value: USER_ROLES.BARISTA, label: 'Pha chế' }
]
```

### 2. Shift Constants (`frontend/src/constants/shift.js`)
```javascript
// Must match backend: backend/domain/cashier/cashier_shift.go
export const SHIFT_STATUS = {
  OPEN: 'OPEN',
  CLOSURE_INITIATED: 'CLOSURE_INITIATED',
  CLOSED: 'CLOSED'
}

export const SHIFT_TYPES = {
  WAITER: 'waiter',
  CASHIER: 'cashier',
  BARISTA: 'barista'
}
```

### 3. Order Constants (`frontend/src/constants/order.js`)
Already exists - needs review for completeness

### 4. Expense Constants (`frontend/src/constants/expense.js`)
Already exists - needs review for completeness

### 5. Facility Constants (`frontend/src/constants/facility.js`)
Already exists - needs review for completeness

## Action Plan

### Phase 1: Create Missing Constants Files
- [ ] Create `frontend/src/constants/user.js`
- [ ] Review and update `frontend/src/constants/shift.js`
- [ ] Review `frontend/src/constants/order.js`
- [ ] Review `frontend/src/constants/expense.js`
- [ ] Review `frontend/src/constants/facility.js`

### Phase 2: Update Views (Priority Order)
1. [ ] UserManagementView.vue - User roles
2. [ ] DashboardView.vue - User roles + Order status
3. [ ] ShiftView.vue - Shift status + User roles
4. [ ] ManagerShiftView.vue - Shift status + types
5. [ ] CashierShiftClosure.vue - Shift status
6. [ ] OrderView.vue - Order status
7. [ ] BaristaView.vue - Order status + Shift status
8. [ ] CashierDashboard.vue - Payment methods
9. [ ] ExpenseManagementView.vue - Expense categories
10. [ ] FacilityManagementView.vue - Facility types

### Phase 3: Verify Backend Sync
- [ ] Compare all frontend constants with backend
- [ ] Document any mismatches
- [ ] Update CONSTANTS_SYNC.md

### Phase 4: Testing
- [ ] Test each view after constants update
- [ ] Verify no broken functionality
- [ ] Check console for errors

## Benefits

### 1. Type Safety
- Autocomplete in IDE
- Catch typos at development time
- Easier refactoring

### 2. Maintainability
- Single source of truth
- Easy to update values
- Clear documentation

### 3. Consistency
- Same values across all files
- No typos or variations
- Backend-frontend sync

### 4. Debugging
- Easier to search for usage
- Clear intent in code
- Better error messages

## Example Refactor

### Before (Hardcoded)
```vue
<template>
  <div v-if="user.role === 'manager'">
    Manager Dashboard
  </div>
  <div v-if="shift.status === 'OPEN'">
    Shift is open
  </div>
</template>

<script>
const isManager = computed(() => user.value?.role === 'manager')
const isOpen = computed(() => shift.value?.status === 'OPEN')
</script>
```

### After (Using Constants)
```vue
<template>
  <div v-if="user.role === USER_ROLES.MANAGER">
    Manager Dashboard
  </div>
  <div v-if="shift.status === SHIFT_STATUS.OPEN">
    Shift is open
  </div>
</template>

<script>
import { USER_ROLES } from '@/constants/user'
import { SHIFT_STATUS } from '@/constants/shift'

const isManager = computed(() => user.value?.role === USER_ROLES.MANAGER)
const isOpen = computed(() => shift.value?.status === SHIFT_STATUS.OPEN)

// Export for template
return {
  USER_ROLES,
  SHIFT_STATUS,
  isManager,
  isOpen
}
</script>
```

## Priority: HIGH

This refactor should be done ASAP to:
1. Prevent bugs from typos
2. Ensure backend-frontend sync
3. Improve code maintainability
4. Make future changes easier

## Estimated Effort

- Phase 1 (Create constants): 2 hours
- Phase 2 (Update views): 4-6 hours
- Phase 3 (Verify sync): 1 hour
- Phase 4 (Testing): 2 hours

**Total**: 9-11 hours

## Next Steps

1. Start with creating user.js constants
2. Update DashboardView.vue (most critical)
3. Continue with other views in priority order
4. Test thoroughly after each update
5. Document any issues found
