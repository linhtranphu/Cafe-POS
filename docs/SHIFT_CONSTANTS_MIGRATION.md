# 🔄 Shift Constants Migration Guide

## 📋 Tổng quan

**Ngày tạo:** 2026-02-04  
**Mục đích:** Migration từ hardcoded shift status strings sang constants pattern

## ⚠️ Vấn đề phát hiện

Nhiều Vue components đang **hardcode shift status strings** thay vì sử dụng constants, vi phạm pattern đã được thiết lập trong dự án.

### Files cần migration:
1. ✅ `frontend/src/views/CashierDashboard.vue` - **ĐÃ HOÀN THÀNH**
2. ❌ `frontend/src/views/ShiftView.vue` - Cần migration
3. ❌ `frontend/src/views/CashierShiftClosure.vue` - Cần migration
4. ❌ `frontend/src/views/ManagerShiftView.vue` - Cần migration
5. ❌ `frontend/src/views/BaristaView.vue` - Cần kiểm tra
6. ❌ `frontend/src/components/CashierShiftManager.vue` - Cần kiểm tra

## ✅ Giải pháp đã triển khai

### 1. Tạo Shift Constants File

**File:** `frontend/src/constants/shift.js`

```javascript
// Shift Status Constants
export const SHIFT_STATUS = {
  OPEN: 'OPEN',
  CLOSED: 'CLOSED'
}

// Cashier Shift Status Constants
export const CASHIER_SHIFT_STATUS = {
  OPEN: 'OPEN',
  CLOSURE_INITIATED: 'CLOSURE_INITIATED',
  CLOSED: 'CLOSED'
}

// Shift Type Constants
export const SHIFT_TYPE = {
  MORNING: 'MORNING',
  AFTERNOON: 'AFTERNOON',
  EVENING: 'EVENING'
}

// Role Type Constants
export const ROLE_TYPE = {
  WAITER: 'waiter',
  BARISTA: 'barista'
}
```

### 2. Backend Constants Reference

**Waiter/Barista Shifts:**
- File: `backend/domain/order/shift.go`
- Constants: `ShiftOpen`, `ShiftClosed`

**Cashier Shifts:**
- File: `backend/domain/cashier/cashier_shift.go`
- Constants: `CashierShiftOpen`, `CashierShiftClosureInitiated`, `CashierShiftClosed`

## 📝 Migration Pattern

### Before (❌ Hardcoded):
```vue
<script setup>
// No imports

const isOpen = (shift) => shift.status === 'OPEN'
</script>

<template>
  <div v-if="shift.status === 'OPEN'">
    <span :class="shift.status === 'OPEN' ? 'bg-green-100' : 'bg-gray-100'">
      {{ shift.status === 'OPEN' ? 'Đang mở' : 'Đã đóng' }}
    </span>
  </div>
</template>
```

### After (✅ Using Constants):
```vue
<script setup>
import { SHIFT_STATUS, CASHIER_SHIFT_STATUS } from '../constants/shift'

const isOpen = (shift) => shift.status === SHIFT_STATUS.OPEN
</script>

<template>
  <div v-if="shift.status === SHIFT_STATUS.OPEN">
    <span :class="shift.status === SHIFT_STATUS.OPEN ? 'bg-green-100' : 'bg-gray-100'">
      {{ shift.status === SHIFT_STATUS.OPEN ? 'Đang mở' : 'Đã đóng' }}
    </span>
  </div>
</template>
```

## 🔍 Common Patterns to Replace

### 1. Status Comparison
```javascript
// ❌ Before
if (shift.status === 'OPEN') { }
if (shift.status === 'CLOSED') { }
if (shift.status === 'CLOSURE_INITIATED') { }

// ✅ After
if (shift.status === SHIFT_STATUS.OPEN) { }
if (shift.status === SHIFT_STATUS.CLOSED) { }
if (shift.status === CASHIER_SHIFT_STATUS.CLOSURE_INITIATED) { }
```

### 2. Template Conditionals
```vue
<!-- ❌ Before -->
<div v-if="shift.status === 'OPEN'">

<!-- ✅ After -->
<div v-if="shift.status === SHIFT_STATUS.OPEN">
```

### 3. Computed Properties
```javascript
// ❌ Before
const openShifts = computed(() => 
  shifts.value.filter(s => s.status === 'OPEN')
)

// ✅ After
const openShifts = computed(() => 
  shifts.value.filter(s => s.status === SHIFT_STATUS.OPEN)
)
```

### 4. Status Display Maps
```javascript
// ❌ Before
const statusMap = {
  'OPEN': '🟢 Đang mở',
  'CLOSED': '🔴 Đã đóng'
}

// ✅ After
const statusMap = {
  [SHIFT_STATUS.OPEN]: '🟢 Đang mở',
  [SHIFT_STATUS.CLOSED]: '🔴 Đã đóng'
}
```

## 📋 Migration Checklist

### CashierDashboard.vue ✅
- [x] Import constants
- [x] Update `getStatusText()` function
- [x] Update `getShiftTypeText()` function
- [x] Update payment method functions
- [x] Update status badge functions
- [x] Test functionality

### ShiftView.vue ❌
- [ ] Import constants
- [ ] Replace hardcoded 'OPEN' comparisons
- [ ] Replace hardcoded 'CLOSED' comparisons
- [ ] Update template conditionals
- [ ] Test functionality

### CashierShiftClosure.vue ❌
- [ ] Import constants
- [ ] Replace status comparisons
- [ ] Update template conditionals
- [ ] Test functionality

### ManagerShiftView.vue ❌
- [ ] Import constants
- [ ] Update filterStatus logic
- [ ] Replace status comparisons
- [ ] Update computed properties
- [ ] Test functionality

### Components to Check ❌
- [ ] CashierShiftManager.vue
- [ ] BaristaView.vue
- [ ] Any other components using shift status

## 🧪 Testing

### Manual Testing Checklist
- [ ] Shift status displays correctly
- [ ] Status filters work properly
- [ ] Conditional rendering works
- [ ] No console errors
- [ ] Backend API calls still work

### Test Cases
1. Open shift → Status shows "🟢 Đang mở"
2. Closed shift → Status shows "🔴 Đã đóng"
3. Closure initiated → Status shows "🟡 Đang đóng"
4. Filter by status → Shows correct shifts
5. Status transitions → Work correctly

## 📚 Related Documentation

- **Pattern Guide:** `docs/CONSTANTS_PATTERN.md`
- **Quick Reference:** `docs/CONSTANTS_QUICK_REFERENCE.md`
- **Implementation Summary:** `docs/CONSTANTS_IMPLEMENTATION_SUMMARY.md`
- **Bug Case Study:** `docs/BUG_FIX_SHIFT_CASH_NOT_UPDATED.md`

## 🎯 Benefits

1. **Type Safety:** Constants prevent typos
2. **Consistency:** Single source of truth
3. **Maintainability:** Easy to update values
4. **Debugging:** Easier to track issues
5. **Backend Sync:** Guaranteed match with backend

## ⚠️ Important Notes

### Case Sensitivity
```javascript
// ✅ CORRECT - Match backend exactly
SHIFT_STATUS.OPEN = 'OPEN'  // Backend: ShiftOpen = "OPEN"

// ❌ WRONG - Case mismatch
SHIFT_STATUS.OPEN = 'open'  // Backend: ShiftOpen = "OPEN"
```

### Waiter vs Cashier Shifts
```javascript
// Waiter/Barista shifts (2 states)
SHIFT_STATUS.OPEN
SHIFT_STATUS.CLOSED

// Cashier shifts (3 states)
CASHIER_SHIFT_STATUS.OPEN
CASHIER_SHIFT_STATUS.CLOSURE_INITIATED
CASHIER_SHIFT_STATUS.CLOSED
```

### Role Types
```javascript
// Note: Role types are lowercase in backend
ROLE_TYPE.WAITER = 'waiter'   // Not 'WAITER'
ROLE_TYPE.BARISTA = 'barista' // Not 'BARISTA'
```

## 🚀 Next Steps

1. **Priority 1:** Complete migration for remaining Vue files
2. **Priority 2:** Add unit tests for constant usage
3. **Priority 3:** Document in team wiki
4. **Priority 4:** Add ESLint rule to prevent hardcoded strings

## 📞 Support

If you encounter issues during migration:
1. Check `docs/CONSTANTS_PATTERN.md` for detailed guide
2. Review `docs/CONSTANTS_QUICK_REFERENCE.md` for quick help
3. Look at `CashierDashboard.vue` as reference implementation
4. Check backend constants in `backend/domain/order/shift.go`

---

**Status:** 🟡 In Progress (1/6 files completed)  
**Last Updated:** 2026-02-04  
**Next File:** ShiftView.vue
