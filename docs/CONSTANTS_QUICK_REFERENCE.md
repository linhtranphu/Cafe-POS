# 📋 Constants Quick Reference Card

> **TL;DR:** Always use constants, never hardcode strings. Values must match backend exactly.

---

## 🚀 Quick Start

### Backend (Go)
```go
// Use constant from domain
if req.PaymentMethod == order.PaymentCash {
    // ✅ Correct
}

// Don't hardcode
if req.PaymentMethod == "cash" {
    // ❌ Wrong - case mismatch possible
}
```

### Frontend (JavaScript)
```javascript
// Import constants
import { PAYMENT_METHOD, ORDER_STATUS } from '../constants/order'

// Use in code
paymentMethod.value = PAYMENT_METHOD.CASH  // ✅ Correct
paymentMethod.value = 'CASH'               // ❌ Wrong

// Use in comparisons
if (order.status === ORDER_STATUS.CREATED) { }  // ✅ Correct
if (order.status === 'CREATED') { }             // ❌ Wrong
```

---

## 📁 Available Constants

### Order Constants
**File:** `frontend/src/constants/order.js`

```javascript
import { 
  ORDER_STATUS,           // Status constants
  PAYMENT_METHOD,         // Payment method constants
  PAYMENT_METHOD_DISPLAY, // Display configs
  ORDER_STATUS_DISPLAY,   // Display configs
  STATUS_FILTER_OPTIONS   // Filter options
} from '../constants/order'
```

**Values:**
- `ORDER_STATUS.CREATED`, `PAID`, `QUEUED`, `IN_PROGRESS`, `READY`, `SERVED`, `CANCELLED`, `LOCKED`
- `PAYMENT_METHOD.CASH`, `TRANSFER`, `QR`

### Shift Constants
**File:** `frontend/src/constants/shift.js`

```javascript
import { 
  SHIFT_STATUS,                  // Waiter/Barista shift status
  CASHIER_SHIFT_STATUS,          // Cashier shift status
  SHIFT_TYPE,                    // Shift types
  ROLE_TYPE,                     // Role types
  SHIFT_STATUS_DISPLAY,          // Display configs
  CASHIER_SHIFT_STATUS_DISPLAY,  // Display configs
  SHIFT_TYPE_DISPLAY,            // Display configs
  ROLE_TYPE_DISPLAY              // Display configs
} from '../constants/shift'
```

**Values:**
- `SHIFT_STATUS.OPEN`, `CLOSED`
- `CASHIER_SHIFT_STATUS.OPEN`, `CLOSURE_INITIATED`, `CLOSED`
- `SHIFT_TYPE.MORNING`, `AFTERNOON`, `EVENING`
- `ROLE_TYPE.WAITER`, `BARISTA`

### Expense Constants
**File:** `frontend/src/constants/expense.js`

```javascript
import { 
  EXPENSE_CATEGORIES,
  EXPENSE_TYPES 
} from '../constants/expense'
```

### Facility Constants
**File:** `frontend/src/constants/facility.js`

```javascript
import { 
  FACILITY_TYPES,
  FACILITY_AREAS 
} from '../constants/facility'
```

### Ingredient Constants
**File:** `frontend/src/constants/ingredient.js`

```javascript
import { 
  INGREDIENT_CATEGORIES,
  INGREDIENT_UNITS 
} from '../constants/ingredient'
```

---

## ✅ Do's

```javascript
// ✅ Import constants
import { ORDER_STATUS } from '../constants/order'

// ✅ Use in refs
const status = ref(ORDER_STATUS.CREATED)

// ✅ Use in comparisons
if (order.status === ORDER_STATUS.PAID) { }

// ✅ Use in templates
<button v-if="order.status === ORDER_STATUS.CREATED">

// ✅ Use in API calls
await orderService.collectPayment(id, {
  payment_method: PAYMENT_METHOD.CASH
})
```

## ❌ Don'ts

```javascript
// ❌ Don't hardcode strings
const status = ref('CREATED')

// ❌ Don't use magic strings
if (order.status === 'PAID') { }

// ❌ Don't hardcode in templates
<button v-if="order.status === 'CREATED'">

// ❌ Don't hardcode in API calls
await orderService.collectPayment(id, {
  payment_method: 'CASH'
})
```

---

## 🔍 Common Patterns

### Status Check
```javascript
// Single status
const isCreated = (order) => order.status === ORDER_STATUS.CREATED

// Multiple statuses
const isProcessing = (order) => 
  [ORDER_STATUS.QUEUED, ORDER_STATUS.IN_PROGRESS].includes(order.status)
```

### Display Text
```javascript
// Get display label
const statusText = ORDER_STATUS_DISPLAY[order.status]?.label

// Get badge class
const badgeClass = ORDER_STATUS_DISPLAY[order.status]?.badge
```

### Payment Method Selection
```javascript
// Use display array for UI
const paymentMethods = PAYMENT_METHOD_DISPLAY
// [
//   { value: 'CASH', label: 'Tiền mặt', icon: '💵' },
//   { value: 'QR', label: 'QR', icon: '📱' },
//   { value: 'TRANSFER', label: 'CK', icon: '🏦' }
// ]

// Selected value
const selected = ref(PAYMENT_METHOD.CASH)
```

---

## 🐛 Debugging

### Check Constant Value
```javascript
console.log('Expected:', ORDER_STATUS.CREATED)  // "CREATED"
console.log('Actual:', order.status)            // "CREATED"
console.log('Match:', order.status === ORDER_STATUS.CREATED)  // true
```

### Check Backend Response
```javascript
// Network tab → Response
{
  "status": "CREATED",        // ← Must match constant exactly
  "payment_method": "CASH"    // ← Must match constant exactly
}
```

### Verify Import
```javascript
// If constant is undefined, check import
import { ORDER_STATUS } from '../constants/order'
console.log(ORDER_STATUS)  // Should show object with all statuses
```

---

## 📚 Full Documentation

**See:** `docs/CONSTANTS_PATTERN.md` for complete guide

**Topics:**
- Why use constants
- Backend/frontend sync
- Adding new constants
- Migration guide
- Best practices
- Troubleshooting

---

## 🆘 Quick Help

### "Constant is undefined"
→ Check import path: `'../constants/order'`

### "Value doesn't match"
→ Check case: Backend uses `"CASH"` not `"cash"`

### "Where to add new constant?"
→ Backend: `backend/domain/*/` → Frontend: `frontend/src/constants/`

### "How to find backend constant?"
→ Check comment in frontend constant file

---

**Last Updated:** 2026-02-04  
**Full Guide:** [CONSTANTS_PATTERN.md](./CONSTANTS_PATTERN.md)
