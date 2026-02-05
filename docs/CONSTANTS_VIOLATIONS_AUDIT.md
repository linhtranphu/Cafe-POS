# 🚨 Constants Pattern Violations - Full Audit Report

**Ngày kiểm tra:** 2026-02-04  
**Phạm vi:** Toàn bộ frontend và backend  
**Kết quả:** 15+ files vi phạm (chủ yếu frontend)

---

## 📊 Tổng quan

### Mức độ nghiêm trọng
- 🔴 **CRITICAL:** Stores (state management) - 5 files
- 🟠 **HIGH:** Vue components - 4 files  
- 🟡 **MEDIUM:** Vue views - 1 file
- 🟢 **BACKEND:** Không có vi phạm

### Tổng số vi phạm
- **Frontend:** 60+ locations
- **Backend:** 0 violations ✅

---

## 🔴 CRITICAL - Stores (State Management)

### 1. `frontend/src/stores/cashier.js`
**Vi phạm:** Hardcoded payment methods

```javascript
// ❌ Lines 30, 33, 36, 40
cashPayments: (state) => 
  state.payments.filter(p => p.payment_method === 'CASH'),

transferPayments: (state) => 
  state.payments.filter(p => p.payment_method === 'TRANSFER'),

qrPayments: (state) => 
  state.payments.filter(p => p.payment_method === 'QR'),

totalCashAmount: (state) => 
  state.payments
    .filter(p => p.payment_method === 'CASH')
    .reduce((sum, p) => sum + p.amount, 0),
```

**Sửa:**
```javascript
import { PAYMENT_METHOD } from '../constants/order'

cashPayments: (state) => 
  state.payments.filter(p => p.payment_method === PAYMENT_METHOD.CASH),

transferPayments: (state) => 
  state.payments.filter(p => p.payment_method === PAYMENT_METHOD.TRANSFER),

qrPayments: (state) => 
  state.payments.filter(p => p.payment_method === PAYMENT_METHOD.QR),

totalCashAmount: (state) => 
  state.payments
    .filter(p => p.payment_method === PAYMENT_METHOD.CASH)
    .reduce((sum, p) => sum + p.amount, 0),
```

---

### 2. `frontend/src/stores/shift.js`
**Vi phạm:** Hardcoded shift status

```javascript
// ❌ Lines 33, 41, 49
hasOpenShift: (state) => {
  return state.currentShift !== null && state.currentShift.status === 'OPEN'
},

openShifts: (state) => {
  return state.shifts.filter(s => s.status === 'OPEN')
},

closedShifts: (state) => {
  return state.shifts.filter(s => s.status === 'CLOSED')
},
```

**Sửa:**
```javascript
import { SHIFT_STATUS } from '../constants/shift'

hasOpenShift: (state) => {
  return state.currentShift !== null && state.currentShift.status === SHIFT_STATUS.OPEN
},

openShifts: (state) => {
  return state.shifts.filter(s => s.status === SHIFT_STATUS.OPEN)
},

closedShifts: (state) => {
  return state.shifts.filter(s => s.status === SHIFT_STATUS.CLOSED)
},
```

---

### 3. `frontend/src/stores/cashierShift.js`
**Vi phạm:** Hardcoded cashier shift status

```javascript
// ❌ Lines 36, 44, 52, 60
hasOpenCashierShift: (state) => {
  return state.currentCashierShift && state.currentCashierShift.status === 'OPEN'
},

canStartCashierShift: (state) => {
  return !state.currentCashierShift || state.currentCashierShift.status === 'CLOSED'
},

isClosureInitiated: (state) => {
  return state.currentCashierShift && state.currentCashierShift.status === 'CLOSURE_INITIATED'
},

isClosed: (state) => {
  return state.currentCashierShift && state.currentCashierShift.status === 'CLOSED'
},
```

**Sửa:**
```javascript
import { CASHIER_SHIFT_STATUS } from '../constants/shift'

hasOpenCashierShift: (state) => {
  return state.currentCashierShift && state.currentCashierShift.status === CASHIER_SHIFT_STATUS.OPEN
},

canStartCashierShift: (state) => {
  return !state.currentCashierShift || state.currentCashierShift.status === CASHIER_SHIFT_STATUS.CLOSED
},

isClosureInitiated: (state) => {
  return state.currentCashierShift && state.currentCashierShift.status === CASHIER_SHIFT_STATUS.CLOSURE_INITIATED
},

isClosed: (state) => {
  return state.currentCashierShift && state.currentCashierShift.status === CASHIER_SHIFT_STATUS.CLOSED
},
```

---

### 4. `frontend/src/stores/order.js`
**Vi phạm:** Hardcoded order status

```javascript
// ❌ Lines 145, 149, 153, 157
createdOrders: (state) => {
  return state.orders.filter(o => o.status === 'CREATED')
},

paidOrders: (state) => {
  return state.orders.filter(o => o.status === 'PAID')
},

inProgressOrders: (state) => {
  return state.orders.filter(o => o.status === 'IN_PROGRESS')
},

servedOrders: (state) => {
  return state.orders.filter(o => o.status === 'SERVED')
}
```

**Sửa:**
```javascript
import { ORDER_STATUS } from '../constants/order'

createdOrders: (state) => {
  return state.orders.filter(o => o.status === ORDER_STATUS.CREATED)
},

paidOrders: (state) => {
  return state.orders.filter(o => o.status === ORDER_STATUS.PAID)
},

inProgressOrders: (state) => {
  return state.orders.filter(o => o.status === ORDER_STATUS.IN_PROGRESS)
},

servedOrders: (state) => {
  return state.orders.filter(o => o.status === ORDER_STATUS.SERVED)
}
```

---

### 5. `frontend/src/stores/barista.js`
**Vi phạm:** Hardcoded order status

```javascript
// ❌ Lines 89, 93, 97, 105, 109, 113
inProgressOrders: (state) => {
  return state.myOrders.filter(o => o.status === 'IN_PROGRESS')
},

readyOrders: (state) => {
  return state.myOrders.filter(o => o.status === 'READY')
},

servedOrders: (state) => {
  return state.myOrders.filter(o => o.status === 'SERVED')
},

inProgressCount: (state) => {
  return state.myOrders.filter(o => o.status === 'IN_PROGRESS').length
},

readyCount: (state) => {
  return state.myOrders.filter(o => o.status === 'READY').length
},

servedCount: (state) => {
  return state.myOrders.filter(o => o.status === 'SERVED').length
}
```

**Sửa:**
```javascript
import { ORDER_STATUS } from '../constants/order'

inProgressOrders: (state) => {
  return state.myOrders.filter(o => o.status === ORDER_STATUS.IN_PROGRESS)
},

readyOrders: (state) => {
  return state.myOrders.filter(o => o.status === ORDER_STATUS.READY)
},

servedOrders: (state) => {
  return state.myOrders.filter(o => o.status === ORDER_STATUS.SERVED)
},

inProgressCount: (state) => {
  return state.myOrders.filter(o => o.status === ORDER_STATUS.IN_PROGRESS).length
},

readyCount: (state) => {
  return state.myOrders.filter(o => o.status === ORDER_STATUS.READY).length
},

servedCount: (state) => {
  return state.myOrders.filter(o => o.status === ORDER_STATUS.SERVED).length
}
```

---

## 🟠 HIGH - Vue Components

### 6. `frontend/src/views/CashierShiftClosure.vue`
**Vi phạm:** Hardcoded cashier shift status (multiple locations)

```javascript
// ❌ Lines 66, 81, 228, 268, 282, 422-424
<div v-if="shift.status === 'OPEN'">
<div v-if="shift.status === 'CLOSURE_INITIATED' && !shift.actual_cash">
<div v-if="shift.status === 'CLOSED'">

const canConfirm = computed(() => {
  if (shift.value.status !== 'CLOSURE_INITIATED') return false
})

const canCloseShift = computed(() => {
  if (shift.value.status !== 'CLOSURE_INITIATED') return false
})

const getStatusText = (status) => {
  const statusMap = {
    'OPEN': '🟢 Đang mở',
    'CLOSURE_INITIATED': '🟡 Đang đóng',
    'CLOSED': '🔴 Đã đóng'
  }
  return statusMap[status] || status
}
```

**Sửa:**
```javascript
import { CASHIER_SHIFT_STATUS } from '../constants/shift'

<div v-if="shift.status === CASHIER_SHIFT_STATUS.OPEN">
<div v-if="shift.status === CASHIER_SHIFT_STATUS.CLOSURE_INITIATED && !shift.actual_cash">
<div v-if="shift.status === CASHIER_SHIFT_STATUS.CLOSED">

const canConfirm = computed(() => {
  if (shift.value.status !== CASHIER_SHIFT_STATUS.CLOSURE_INITIATED) return false
})

const canCloseShift = computed(() => {
  if (shift.value.status !== CASHIER_SHIFT_STATUS.CLOSURE_INITIATED) return false
})

const getStatusText = (status) => {
  const statusMap = {
    [CASHIER_SHIFT_STATUS.OPEN]: '🟢 Đang mở',
    [CASHIER_SHIFT_STATUS.CLOSURE_INITIATED]: '🟡 Đang đóng',
    [CASHIER_SHIFT_STATUS.CLOSED]: '🔴 Đã đóng'
  }
  return statusMap[status] || status
}
```

---

### 7. `frontend/src/views/ManagerShiftView.vue`
**Vi phạm:** Hardcoded shift status and role types (multiple locations)

```javascript
// ❌ Lines 21, 26, 396, 401, 406, 466, 472
<button @click="filterStatus = 'OPEN'">
<button @click="filterStatus = 'CLOSED'">

const openWaiterShifts = computed(() => {
  return shifts.filter(s => s.status === 'OPEN' && s.role_type === 'waiter')
})

const openBaristaShifts = computed(() => {
  return shifts.filter(s => s.status === 'OPEN' && s.role_type === 'barista')
})

const openCashierShifts = computed(() => {
  return shifts.filter(s => s.status === 'OPEN')
})

const getStatusClass = (status) => {
  return status === 'OPEN' ? 'bg-green-100' : 'bg-gray-100'
}

const getStatusText = (status) => {
  return status === 'OPEN' ? 'Đang mở' : 'Đã đóng'
}
```

**Sửa:**
```javascript
import { SHIFT_STATUS, CASHIER_SHIFT_STATUS, ROLE_TYPE } from '../constants/shift'

<button @click="filterStatus = SHIFT_STATUS.OPEN">
<button @click="filterStatus = SHIFT_STATUS.CLOSED">

const openWaiterShifts = computed(() => {
  return shifts.filter(s => s.status === SHIFT_STATUS.OPEN && s.role_type === ROLE_TYPE.WAITER)
})

const openBaristaShifts = computed(() => {
  return shifts.filter(s => s.status === SHIFT_STATUS.OPEN && s.role_type === ROLE_TYPE.BARISTA)
})

const openCashierShifts = computed(() => {
  return shifts.filter(s => s.status === CASHIER_SHIFT_STATUS.OPEN)
})

const getStatusClass = (status) => {
  return status === SHIFT_STATUS.OPEN ? 'bg-green-100' : 'bg-gray-100'
}

const getStatusText = (status) => {
  return status === SHIFT_STATUS.OPEN ? 'Đang mở' : 'Đã đóng'
}
```

---

### 8. `frontend/src/views/ShiftView.vue`
**Vi phạm:** Hardcoded shift status (multiple locations)

```vue
<!-- ❌ Lines 186, 188, 197, 201, 205, 211 -->
<span :class="shift.status === 'OPEN' ? 'bg-green-100' : 'bg-gray-100'">
  {{ shift.status === 'OPEN' ? 'Đang mở' : 'Đã đóng' }}
</span>

<div v-if="shift.status === 'CLOSED'">
<div v-if="shift.status === 'CLOSED'">
<div v-if="shift.status === 'CLOSED'">
<button v-if="isCashier && shift.status === 'OPEN'">
```

**Sửa:**
```vue
<script setup>
import { SHIFT_STATUS } from '../constants/shift'
</script>

<span :class="shift.status === SHIFT_STATUS.OPEN ? 'bg-green-100' : 'bg-gray-100'">
  {{ shift.status === SHIFT_STATUS.OPEN ? 'Đang mở' : 'Đã đóng' }}
</span>

<div v-if="shift.status === SHIFT_STATUS.CLOSED">
<div v-if="shift.status === SHIFT_STATUS.CLOSED">
<div v-if="shift.status === SHIFT_STATUS.CLOSED">
<button v-if="isCashier && shift.status === SHIFT_STATUS.OPEN">
```

---

### 9. `frontend/src/components/CashierShiftManager.vue`
**Vi phạm:** Hardcoded shift status

```javascript
// ❌ Line 136
const canCloseShift = computed(() => {
  return currentShift.value && currentShift.value.status === 'OPEN'
})
```

**Sửa:**
```javascript
import { CASHIER_SHIFT_STATUS } from '../constants/shift'

const canCloseShift = computed(() => {
  return currentShift.value && currentShift.value.status === CASHIER_SHIFT_STATUS.OPEN
})
```

---

## 🟡 MEDIUM - Vue Views

### 10. `frontend/src/views/DashboardView.vue`
**Vi phạm:** Hardcoded order status

```javascript
// ❌ Lines 523, 530, 538, 542
const todayRevenue = computed(() => {
  return orders.value
    .filter(o => new Date(o.created_at).toDateString() === today && o.status !== 'CANCELLED')
    .reduce((sum, o) => sum + o.total, 0)
})

const completedOrders = computed(() => {
  return orders.value.filter(o => 
    new Date(o.created_at).toDateString() === today && o.status === 'SERVED'
  ).length
})

const pendingOrders = computed(() => {
  if (user.value?.role === 'manager') {
    return orders.value.filter(o => 
      o.status !== 'SERVED' && o.status !== 'CANCELLED'
    ).length
  }
  return orders.value.filter(o => o.status === 'CREATED').length
})
```

**Sửa:**
```javascript
import { ORDER_STATUS } from '../constants/order'

const todayRevenue = computed(() => {
  return orders.value
    .filter(o => new Date(o.created_at).toDateString() === today && o.status !== ORDER_STATUS.CANCELLED)
    .reduce((sum, o) => sum + o.total, 0)
})

const completedOrders = computed(() => {
  return orders.value.filter(o => 
    new Date(o.created_at).toDateString() === today && o.status === ORDER_STATUS.SERVED
  ).length
})

const pendingOrders = computed(() => {
  if (user.value?.role === 'manager') {
    return orders.value.filter(o => 
      o.status !== ORDER_STATUS.SERVED && o.status !== ORDER_STATUS.CANCELLED
    ).length
  }
  return orders.value.filter(o => o.status === ORDER_STATUS.CREATED).length
})
```

---

## 🟢 Backend Status

### ✅ Backend Code - NO VIOLATIONS

Backend code correctly uses constants throughout:

```go
// ✅ backend/domain/order/shift.go
const (
	ShiftOpen   ShiftStatus = "OPEN"
	ShiftClosed ShiftStatus = "CLOSED"
)

// ✅ backend/domain/cashier/cashier_shift.go
const (
	CashierShiftOpen             CashierShiftStatus = "OPEN"
	CashierShiftClosureInitiated CashierShiftStatus = "CLOSURE_INITIATED"
	CashierShiftClosed           CashierShiftStatus = "CLOSED"
)

// ✅ backend/domain/order/order.go
const (
	PaymentCash     PaymentMethod = "CASH"
	PaymentTransfer PaymentMethod = "TRANSFER"
	PaymentQR       PaymentMethod = "QR"
)

// ✅ Usage in backend code
if req.PaymentMethod == order.PaymentCash {
    // Correct usage
}

if shift.Status == order.ShiftOpen {
    // Correct usage
}
```

---

## 📋 Migration Priority

### Phase 1: CRITICAL (Stores) - **DO FIRST**
1. ✅ `frontend/src/stores/cashier.js`
2. ✅ `frontend/src/stores/shift.js`
3. ✅ `frontend/src/stores/cashierShift.js`
4. ✅ `frontend/src/stores/order.js`
5. ✅ `frontend/src/stores/barista.js`

### Phase 2: HIGH (Components)
6. ✅ `frontend/src/views/CashierShiftClosure.vue`
7. ✅ `frontend/src/views/ManagerShiftView.vue`
8. ✅ `frontend/src/views/ShiftView.vue`
9. ✅ `frontend/src/components/CashierShiftManager.vue`

### Phase 3: MEDIUM (Views)
10. ✅ `frontend/src/views/DashboardView.vue`

### Phase 4: Verification
- [ ] Run all tests
- [ ] Manual testing of all features
- [ ] Verify no console errors
- [ ] Check API calls still work

---

## 🎯 Impact Analysis

### Why This Matters

1. **Type Safety:** Hardcoded strings are prone to typos
2. **Maintainability:** Changes require updating multiple files
3. **Bug Risk:** Case mismatches cause silent failures (like the CASH vs cash bug)
4. **Consistency:** No guarantee frontend matches backend
5. **Debugging:** Harder to track where values are used

### Real Bug Example

**Bug:** Payment method case mismatch
- Frontend sent: `payment_method: 'CASH'`
- Backend checked: `if req.PaymentMethod == "cash"`
- Result: Shift cash not updated ❌

**Fix:** Use constants
- Frontend: `payment_method: PAYMENT_METHOD.CASH`
- Backend: `if req.PaymentMethod == order.PaymentCash`
- Result: Works correctly ✅

---

## 🚀 Next Steps

### Immediate Actions
1. **Create migration tasks** for each file
2. **Assign priorities** based on criticality
3. **Start with stores** (highest impact)
4. **Test thoroughly** after each migration

### Long-term Solutions
1. **Add ESLint rules** to prevent hardcoded strings
2. **Add unit tests** for constant values
3. **Document pattern** in team wiki
4. **Code review checklist** for new code

### ESLint Rule Example
```javascript
// .eslintrc.js
rules: {
  'no-restricted-syntax': [
    'error',
    {
      selector: "BinaryExpression[operator='==='][right.value=/^(OPEN|CLOSED|CASH|TRANSFER)$/]",
      message: 'Use constants instead of hardcoded strings'
    }
  ]
}
```

---

## 📚 Related Documentation

- **Pattern Guide:** `docs/CONSTANTS_PATTERN.md`
- **Quick Reference:** `docs/CONSTANTS_QUICK_REFERENCE.md`
- **Migration Guide:** `docs/SHIFT_CONSTANTS_MIGRATION.md`
- **Bug Case Study:** `docs/BUG_FIX_SHIFT_CASH_NOT_UPDATED.md`

---

**Report Generated:** 2026-02-04  
**Status:** 🔴 Critical violations found  
**Action Required:** Immediate migration needed  
**Estimated Effort:** 4-6 hours for complete migration
