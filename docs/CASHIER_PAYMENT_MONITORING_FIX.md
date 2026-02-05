# 🔧 Fix: Cashier Payment Monitoring - Distribution Model

**Date:** 2026-02-04  
**Issue:** Cashier không thấy payments trong dashboard  
**Root Cause:** Logic sai - tìm payments theo cashier shift thay vì waiter shift

---

## 🎯 Vấn đề

### Mô hình Distribution

Trong hệ thống POS này, sử dụng **mô hình distribution**:

```
Waiter (Phục vụ)
  ↓
Thu tiền từ khách
  ↓
Giữ tiền trong ca của waiter
  ↓
Bàn giao tiền cho Cashier
  ↓
Cashier (Thu ngân)
  ↓
Nhận tiền qua handover
```

### Logic Cũ (SAI)

```javascript
// Frontend: Chọn CASHIER shift
selectedShift = cashierShift.id

// Backend: Tìm orders theo cashier shift
orders = findByShiftID(cashierShift.id)

// ❌ KẾT QUẢ: Không có orders
// Vì orders thuộc về WAITER shift, không phải cashier shift
```

### Tại sao sai?

1. **Orders thuộc về waiter shift**
   - Khi waiter tạo order → `order.shift_id = waiter_shift_id`
   - Waiter thu tiền → Tiền vào waiter shift

2. **Cashier shift là shift riêng**
   - Cashier có shift riêng để quản lý ca làm việc
   - Cashier KHÔNG tạo orders trực tiếp
   - Cashier chỉ nhận tiền qua handover

3. **Payments không thuộc cashier shift**
   - Payments gắn với orders
   - Orders gắn với waiter shifts
   - → Payments KHÔNG có trong cashier shift

---

## ✅ Giải pháp

### Logic Mới (ĐÚNG)

Cashier cần xem payments từ **WAITER SHIFTS**, không phải cashier shift.

```javascript
// Frontend: Chọn WAITER shift (không phải cashier shift)
selectedShift = waiterShift.id

// Backend: Tìm orders theo waiter shift
orders = findByShiftID(waiterShift.id)

// ✅ KẾT QUẢ: Có payments
// Vì orders thuộc về waiter shift
```

---

## 🔧 Code Changes

### 1. Backend Service

**File:** `backend/application/services/payment_oversight_service.go`

**Before:**
```go
// Tìm theo shift_id (có thể là cashier shift - SAI)
func (s *PaymentOversightService) GetPaymentsByShift(shiftID string) ([]*PaymentSummary, error) {
    shiftObjID, _ := primitive.ObjectIDFromHex(shiftID)
    orders, _ := s.orderRepo.FindByShiftID(ctx, shiftObjID)
    // ...
}
```

**After:**
```go
// Tìm theo shift_id (WAITER shift - ĐÚNG)
// In distribution model: Cashier monitors payments from waiter shifts
// Orders belong to waiter shifts, cashier receives cash via handovers
func (s *PaymentOversightService) GetPaymentsByShift(shiftID string) ([]*PaymentSummary, error) {
    shiftObjID, _ := primitive.ObjectIDFromHex(shiftID)
    
    // Find orders by waiter shift ID
    orders, _ := s.orderRepo.FindByShiftID(ctx, shiftObjID)
    
    var payments []*PaymentSummary
    for _, ord := range orders {
        // Include all paid orders
        if ord.Status == PAID || ord.Status == IN_PROGRESS || 
           ord.Status == SERVED || ord.Status == QUEUED {
            payments = append(payments, &PaymentSummary{...})
        }
    }
    return payments, nil
}
```

**Changes:**
- ✅ Added comment explaining distribution model
- ✅ Clarified that shift_id is WAITER shift
- ✅ Added QUEUED status to include orders being prepared

### 2. Frontend View

**File:** `frontend/src/views/CashierDashboard.vue`

**Before:**
```vue
<!-- Shift Selector - Only Cashier Shifts -->
<select v-model="selectedShift">
  <option value="">-- Chọn ca thu ngân --</option>
  <option v-for="shift in cashierShifts" :key="shift.id">
    {{ shift.cashier_name }} - {{ shift.status }}
  </option>
</select>
```

**After:**
```vue
<!-- Shift Selector - Waiter Shifts for Payment Monitoring -->
<label>📅 Chọn ca phục vụ để xem thanh toán</label>
<p class="text-xs text-gray-500 mb-2">
  💡 Chọn ca của waiter để xem các thanh toán trong ca đó
</p>
<select v-model="selectedShift">
  <option value="">-- Chọn ca phục vụ --</option>
  <option v-for="shift in waiterShifts" :key="shift.id">
    {{ shift.waiter_name }} ({{ shift.role }}) - {{ shift.status }}
  </option>
</select>
```

**Changes:**
- ✅ Changed label from "ca thu ngân" to "ca phục vụ"
- ✅ Added helper text explaining to select waiter shift
- ✅ Changed data source from `cashierShifts` to `waiterShifts`
- ✅ Display waiter name and role

### 3. Frontend Script

**File:** `frontend/src/views/CashierDashboard.vue`

**Before:**
```javascript
import { useCashierShiftStore } from '../stores/cashierShift'

const cashierShiftStore = useCashierShiftStore()
const cashierShifts = computed(() => cashierShiftStore.cashierShifts)

onMounted(async () => {
  await cashierShiftStore.fetchMyCashierShifts()
})
```

**After:**
```javascript
import { useCashierShiftStore } from '../stores/cashierShift'
import { useShiftStore } from '../stores/shift'

const cashierShiftStore = useCashierShiftStore()
const shiftStore = useShiftStore()

const cashierShifts = computed(() => cashierShiftStore.cashierShifts)
const waiterShifts = computed(() => shiftStore.shifts) // All shifts including waiter

onMounted(async () => {
  // Fetch cashier shifts for cashier shift manager
  await cashierShiftStore.fetchMyCashierShifts()
  // Fetch all shifts (including waiter shifts) for payment monitoring
  await shiftStore.fetchAllShifts()
})
```

**Changes:**
- ✅ Import `useShiftStore`
- ✅ Add `waiterShifts` computed property
- ✅ Fetch all shifts on mount

---

## 🎨 UI Changes

### Before (SAI)

```
┌─────────────────────────────────────────┐
│  📅 Chọn ca thu ngân để xem            │
│  [-- Chọn ca thu ngân --]             │
│  [Ca 1: Nguyễn Văn A (Cashier)]       │
│  [Ca 2: Trần Thị B (Cashier)]         │
└─────────────────────────────────────────┘
         ↓ Chọn cashier shift
         ↓
┌─────────────────────────────────────────┐
│  💳 Danh sách thanh toán               │
│  0 giao dịch                           │
│         📭                              │
│  Chưa có thanh toán nào                │
└─────────────────────────────────────────┘
```

### After (ĐÚNG)

```
┌─────────────────────────────────────────┐
│  📅 Chọn ca phục vụ để xem thanh toán  │
│  💡 Chọn ca của waiter để xem các      │
│     thanh toán trong ca đó             │
│  [-- Chọn ca phục vụ --]              │
│  [Ca 1: Lê Văn C (waiter) - 🟢 Đang mở]│
│  [Ca 2: Phạm Thị D (waiter) - 🔴 Đã đóng]│
└─────────────────────────────────────────┘
         ↓ Chọn waiter shift
         ↓
┌─────────────────────────────────────────┐
│  💳 Danh sách thanh toán               │
│  5 giao dịch                           │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │ Trần Thị B          50,000₫      │ │
│  │ 04/02 - 10:05       💵 Tiền mặt  │ │
│  │ [✏️ Điều chỉnh] [⚠️ Báo lỗi] [🔒]│ │
│  └───────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

---

## 🔄 Flow Mới

### 1. Cashier Login

```
Login với role: cashier
  ↓
Navigate to: /cashier
  ↓
Dashboard loads
```

### 2. Load Shifts

```
onMounted()
  ↓
Fetch cashier shifts (for shift manager)
  ↓
Fetch ALL shifts (for payment monitoring)
  ↓
Dropdown shows waiter shifts
```

### 3. Select Waiter Shift

```
User chọn waiter shift từ dropdown
  ↓
@change="loadPayments"
  ↓
GET /api/cashier/shifts/:waiter_shift_id/payments
  ↓
Backend finds orders by waiter_shift_id
  ↓
Return payments
  ↓
Display in dashboard
```

### 4. Monitor Payments

```
Cashier xem payments từ waiter shift
  ↓
Có thể:
- ✏️ Điều chỉnh payment
- ⚠️ Báo lỗi (discrepancy)
- 🔒 Khóa order
```

---

## 📊 Data Model

### Order Structure

```javascript
{
  id: "507f1f77bcf86cd799439011",
  order_number: "20260204-100530-123",
  shift_id: "507f1f77bcf86cd799439010", // ← WAITER shift ID
  waiter_id: "507f1f77bcf86cd799439009",
  waiter_name: "Lê Văn C",
  payment_method: "CASH",
  status: "PAID",
  total: 50000,
  paid_at: "2026-02-04T10:05:30Z"
}
```

### Shift Types

```javascript
// Waiter Shift (có orders)
{
  id: "507f1f77bcf86cd799439010",
  role: "waiter",
  waiter_id: "507f1f77bcf86cd799439009",
  waiter_name: "Lê Văn C",
  status: "OPEN",
  current_cash: 150000, // Tiền waiter đang giữ
  orders: [...] // Orders trong shift này
}

// Cashier Shift (KHÔNG có orders)
{
  id: "507f1f77bcf86cd799439020",
  role: "cashier",
  cashier_id: "507f1f77bcf86cd799439019",
  cashier_name: "Nguyễn Văn A",
  status: "OPEN",
  received_cash: 0, // Tiền nhận từ handover
  handovers: [...] // Handovers đã nhận
}
```

---

## ✅ Testing

### Test Case 1: View Waiter Payments

**Steps:**
1. Login as cashier
2. Go to `/cashier`
3. Select a waiter shift from dropdown
4. ✅ Should see payments from that waiter shift

**Expected:**
- Dropdown shows waiter shifts (not cashier shifts)
- Payments display correctly
- Can perform actions (override, report discrepancy, lock)

### Test Case 2: Multiple Waiter Shifts

**Steps:**
1. Create 2 waiter shifts with different waiters
2. Each waiter creates orders and collects payments
3. Cashier selects shift 1
4. ✅ Should see only shift 1 payments
5. Cashier selects shift 2
6. ✅ Should see only shift 2 payments

---

## 🎯 Key Takeaways

### Distribution Model

1. **Waiter owns orders**
   - Orders created by waiter
   - Payments collected by waiter
   - Cash held in waiter shift

2. **Cashier monitors payments**
   - Views payments from waiter shifts
   - Can override/adjust payments
   - Reports discrepancies

3. **Cash flow via handover**
   - Waiter hands over cash to cashier
   - Cashier receives via handover process
   - Cashier shift tracks received cash

### Why This Design?

- ✅ **Accountability:** Each waiter responsible for their orders
- ✅ **Traceability:** Clear audit trail of who collected what
- ✅ **Flexibility:** Cashier can monitor any waiter shift
- ✅ **Security:** Handover process ensures proper cash transfer

---

## 📚 Related Documentation

- [CASHIER_PAYMENT_LIST_LOGIC.md](./CASHIER_PAYMENT_LIST_LOGIC.md) - Original logic (now updated)
- [CASH_HANDOVER_COMPLETE_SUMMARY.md](./CASH_HANDOVER_COMPLETE_SUMMARY.md) - Handover process
- [PAYMENT_DISCREPANCY_IMPLEMENTATION.md](./PAYMENT_DISCREPANCY_IMPLEMENTATION.md) - Discrepancy handling

---

**Status:** ✅ **FIXED**  
**Backend Rebuilt:** Yes  
**Backend Restarted:** Yes  
**Ready for Testing:** Yes

---

**Last Updated:** 2026-02-04  
**Fixed By:** Development Team
