# 📋 Constants Pattern - Frontend/Backend Synchronization

## 🎯 Mục đích

Document này mô tả pattern sử dụng constants để đảm bảo đồng bộ giữa frontend và backend, tránh lỗi do hardcode string values.

## ⚠️ Vấn đề đã gặp

### Bug: Payment Method Case Mismatch
**Ngày phát hiện:** 2026-02-04

**Mô tả:**
- Frontend gửi `payment_method: 'CASH'` (uppercase)
- Backend check `if req.PaymentMethod == "cash"` (lowercase)
- Kết quả: Shift cash không được update sau khi thu tiền

**Root cause:** Hardcode string values ở cả frontend và backend, không có single source of truth.

## ✅ Giải pháp: Constants Pattern

### 1. Backend Constants (Go)

**File:** `backend/domain/order/order.go`

```go
// Payment Method Constants
type PaymentMethod string

const (
	PaymentCash     PaymentMethod = "CASH"
	PaymentTransfer PaymentMethod = "TRANSFER"
	PaymentQR       PaymentMethod = "QR"
)

// Order Status Constants
type OrderStatus string

const (
	StatusCreated    OrderStatus = "CREATED"
	StatusPaid       OrderStatus = "PAID"
	StatusQueued     OrderStatus = "QUEUED"
	StatusInProgress OrderStatus = "IN_PROGRESS"
	StatusReady      OrderStatus = "READY"
	StatusServed     OrderStatus = "SERVED"
	StatusCancelled  OrderStatus = "CANCELLED"
	StatusLocked     OrderStatus = "LOCKED"
)
```

**Sử dụng trong code:**
```go
// ✅ ĐÚNG - Sử dụng constant
if req.PaymentMethod == order.PaymentCash {
    // Update shift cash
}

// ❌ SAI - Hardcode string
if req.PaymentMethod == "cash" {
    // Có thể sai case
}
```

### 2. Frontend Constants (JavaScript)

**File:** `frontend/src/constants/order.js`

```javascript
// Payment Method Constants
// Must match backend: backend/domain/order/order.go (PaymentMethod type)
export const PAYMENT_METHOD = {
  CASH: 'CASH',
  TRANSFER: 'TRANSFER',
  QR: 'QR'
}

// Order Status Constants
// Must match backend: backend/domain/order/order.go (OrderStatus type)
export const ORDER_STATUS = {
  CREATED: 'CREATED',
  PAID: 'PAID',
  QUEUED: 'QUEUED',
  IN_PROGRESS: 'IN_PROGRESS',
  READY: 'READY',
  SERVED: 'SERVED',
  CANCELLED: 'CANCELLED',
  LOCKED: 'LOCKED'
}

// Display configurations
export const PAYMENT_METHOD_DISPLAY = [
  { value: PAYMENT_METHOD.CASH, label: 'Tiền mặt', icon: '💵' },
  { value: PAYMENT_METHOD.QR, label: 'QR', icon: '📱' },
  { value: PAYMENT_METHOD.TRANSFER, label: 'CK', icon: '🏦' }
]

export const ORDER_STATUS_DISPLAY = {
  [ORDER_STATUS.CREATED]: { 
    label: 'Mới tạo', 
    icon: '🆕', 
    badge: 'bg-gray-100 text-gray-800' 
  },
  [ORDER_STATUS.PAID]: { 
    label: 'Đã thanh toán', 
    icon: '💰', 
    badge: 'bg-green-100 text-green-800' 
  },
  // ... other statuses
}
```

**Sử dụng trong Vue component:**
```vue
<script setup>
import { PAYMENT_METHOD, ORDER_STATUS } from '../constants/order'

// ✅ ĐÚNG - Sử dụng constant
const paymentMethod = ref(PAYMENT_METHOD.CASH)

const isCreated = (order) => order.status === ORDER_STATUS.CREATED

// ❌ SAI - Hardcode string
const paymentMethod = ref('CASH')
const isCreated = (order) => order.status === 'CREATED'
</script>

<template>
  <!-- ✅ ĐÚNG -->
  <button v-if="order.status === ORDER_STATUS.CREATED">
    Thu tiền
  </button>
  
  <!-- ❌ SAI -->
  <button v-if="order.status === 'CREATED'">
    Thu tiền
  </button>
</template>
```

## 📁 Cấu trúc Constants Files

### Frontend Constants Structure
```
frontend/src/constants/
├── order.js          # Order status, payment methods
├── expense.js        # Expense categories, types
├── facility.js       # Facility types, areas
└── ingredient.js     # Ingredient categories, units
```

### Backend Constants Structure
```
backend/domain/
├── order/
│   └── order.go      # OrderStatus, PaymentMethod
├── expense/
│   └── expense.go    # ExpenseCategory, ExpenseType
├── facility/
│   └── facility.go   # FacilityType, FacilityArea
└── ingredient/
    └── ingredient.go # IngredientCategory, Unit
```

## 🔄 Quy trình thêm constant mới

### 1. Thêm vào Backend (Go)
```go
// backend/domain/order/order.go
const (
    PaymentCash     PaymentMethod = "CASH"
    PaymentTransfer PaymentMethod = "TRANSFER"
    PaymentQR       PaymentMethod = "QR"
    PaymentCard     PaymentMethod = "CARD"  // ← Thêm mới
)
```

### 2. Thêm vào Frontend (JS)
```javascript
// frontend/src/constants/order.js
export const PAYMENT_METHOD = {
  CASH: 'CASH',
  TRANSFER: 'TRANSFER',
  QR: 'QR',
  CARD: 'CARD'  // ← Thêm mới (phải match backend)
}

// Thêm display config
export const PAYMENT_METHOD_DISPLAY = [
  // ... existing
  { value: PAYMENT_METHOD.CARD, label: 'Thẻ', icon: '💳' }  // ← Thêm mới
]
```

### 3. Sử dụng trong code
```javascript
// ✅ ĐÚNG
paymentMethod.value = PAYMENT_METHOD.CARD

// ❌ SAI
paymentMethod.value = 'CARD'
```

## 📝 Checklist khi thêm constant mới

- [ ] Thêm constant vào backend (Go)
- [ ] Thêm constant vào frontend (JS)
- [ ] Verify values match EXACTLY (case-sensitive)
- [ ] Thêm display config nếu cần (label, icon, badge)
- [ ] Update tất cả hardcoded strings thành constants
- [ ] Test API call với constant mới
- [ ] Document trong file này nếu là pattern mới

## 🚨 Lưu ý quan trọng

### 1. Case Sensitivity
```javascript
// ✅ ĐÚNG - Match backend exactly
PAYMENT_METHOD.CASH = 'CASH'  // Backend: PaymentCash = "CASH"

// ❌ SAI - Case mismatch
PAYMENT_METHOD.CASH = 'cash'  // Backend: PaymentCash = "CASH"
```

### 2. String Values phải match 100%
```javascript
// Backend (Go)
const StatusInProgress OrderStatus = "IN_PROGRESS"

// Frontend (JS) - Must match exactly
ORDER_STATUS.IN_PROGRESS = 'IN_PROGRESS'  // ✅ ĐÚNG

ORDER_STATUS.IN_PROGRESS = 'in_progress'  // ❌ SAI
ORDER_STATUS.IN_PROGRESS = 'InProgress'   // ❌ SAI
ORDER_STATUS.IN_PROGRESS = 'IN-PROGRESS'  // ❌ SAI
```

### 3. Comment reference trong frontend
```javascript
// Payment Method Constants
// Must match backend: backend/domain/order/order.go (PaymentMethod type)
export const PAYMENT_METHOD = {
  CASH: 'CASH',
  // ...
}
```

## 🔍 Debugging Constants Issues

### Kiểm tra backend log
```bash
# Thêm debug log trong backend
fmt.Printf("DEBUG: PaymentMethod received: %s\n", req.PaymentMethod)
fmt.Printf("DEBUG: Comparing with constant: %s\n", order.PaymentCash)
```

### Kiểm tra frontend network tab
```javascript
// Check request payload
{
  "payment_method": "CASH",  // ← Verify value
  "amount": 50000
}
```

### Verify constant values
```bash
# Backend
grep -r "PaymentCash.*=" backend/domain/order/

# Frontend
grep -r "CASH:" frontend/src/constants/
```

## 📚 Related Files

### Backend
- `backend/domain/order/order.go` - Order constants
- `backend/domain/expense/expense.go` - Expense constants
- `backend/domain/facility/facility.go` - Facility constants
- `backend/domain/ingredient/ingredient.go` - Ingredient constants

### Frontend
- `frontend/src/constants/order.js` - Order constants
- `frontend/src/constants/expense.js` - Expense constants
- `frontend/src/constants/facility.js` - Facility constants
- `frontend/src/constants/ingredient.js` - Ingredient constants

### Documentation
- `docs/BUG_FIX_SHIFT_CASH_NOT_UPDATED.md` - Bug case study
- `docs/API_CONTRACTS.md` - API request/response formats

## 🎓 Best Practices

### 1. Always import constants
```javascript
// ✅ ĐÚNG
import { PAYMENT_METHOD, ORDER_STATUS } from '../constants/order'

// ❌ SAI - Không import
const status = 'CREATED'
```

### 2. Use constants in comparisons
```javascript
// ✅ ĐÚNG
if (order.status === ORDER_STATUS.CREATED) { }

// ❌ SAI
if (order.status === 'CREATED') { }
```

### 3. Use constants in assignments
```javascript
// ✅ ĐÚNG
paymentMethod.value = PAYMENT_METHOD.CASH

// ❌ SAI
paymentMethod.value = 'CASH'
```

### 4. Use constants in API calls
```javascript
// ✅ ĐÚNG
await orderService.collectPayment(id, {
  payment_method: PAYMENT_METHOD.CASH,
  amount: 50000
})

// ❌ SAI
await orderService.collectPayment(id, {
  payment_method: 'CASH',
  amount: 50000
})
```

## 🔄 Migration từ hardcoded strings

### Step 1: Tạo constants file
```javascript
// frontend/src/constants/order.js
export const ORDER_STATUS = {
  CREATED: 'CREATED',
  PAID: 'PAID',
  // ...
}
```

### Step 2: Find & Replace
```bash
# Find all hardcoded strings
grep -r "status === 'CREATED'" frontend/src/

# Replace with constant
# Before: order.status === 'CREATED'
# After:  order.status === ORDER_STATUS.CREATED
```

### Step 3: Import constants
```javascript
import { ORDER_STATUS } from '../constants/order'
```

### Step 4: Test thoroughly
- Test all status transitions
- Test all payment methods
- Verify API calls work correctly

---

**Last Updated:** 2026-02-04  
**Author:** Development Team  
**Related Bug:** Payment method case mismatch causing shift cash not updated
