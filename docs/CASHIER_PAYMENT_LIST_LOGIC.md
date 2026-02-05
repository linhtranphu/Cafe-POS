# 💳 Logic Hiển Thị Danh Sách Thanh Toán - Cashier Dashboard

**File:** `frontend/src/views/CashierDashboard.vue`  
**URL:** `http://localhost:5173/#/cashier`

---

## 🎯 Điều kiện hiển thị

### ✅ Khi nào hiển thị danh sách thanh toán?

Danh sách thanh toán hiển thị khi:

1. **User đã chọn shift** từ dropdown
   - `selectedShift.value` có giá trị (không rỗng)
   - Dropdown: "📅 Chọn ca thu ngân để xem"

2. **Payments đã được load**
   - `payments.length > 0` → Hiển thị danh sách
   - `payments.length === 0` → Hiển thị empty state

### ❌ Khi nào KHÔNG hiển thị?

- Chưa chọn shift (`selectedShift.value === ''`)
- Đang loading (`loading === true`)
- Có lỗi (`error !== null`)

---

## 🔄 Flow Logic

### 1. Page Load (onMounted)

```javascript
onMounted(async () => {
  // Fetch cashier shifts instead of all shifts
  await cashierShiftStore.fetchMyCashierShifts()
  await cashierStore.getPendingDiscrepancies()
  await cashierStore.fetchPendingHandovers()
})
```

**Kết quả:**
- ✅ Dropdown có danh sách cashier shifts
- ❌ Payments chưa load (chưa chọn shift)
- 📭 Empty state hiển thị: "Chưa có thanh toán nào - Chọn ca làm việc để xem"

### 2. User Chọn Shift

**Trigger:** `@change="loadPayments"` trên dropdown

```javascript
const loadPayments = async () => {
  if (selectedShift.value) {
    await Promise.all([
      cashierStore.getShiftStatus(selectedShift.value),
      cashierStore.getPaymentsByShift(selectedShift.value)
    ])
  }
}
```

**API Calls:**
1. `GET /api/cashier/shifts/:id/status` → Lấy shift status
2. `GET /api/cashier/shifts/:id/payments` → Lấy payments

**Kết quả:**
- ✅ `shiftStatus` được populate
- ✅ `payments[]` được populate
- ✅ Danh sách thanh toán hiển thị

### 3. Refresh Data

**Trigger:** Click button 🔄

```javascript
const refreshData = async () => {
  if (selectedShift.value) {
    await Promise.all([
      cashierStore.getShiftStatus(selectedShift.value),
      cashierStore.getPaymentsByShift(selectedShift.value)
    ])
  }
  await cashierStore.getPendingDiscrepancies()
  await cashierStore.fetchPendingHandovers()
}
```

---

## 📊 Data Structure

### Shift Dropdown Data

**Source:** `cashierShiftStore.cashierShifts`

```javascript
[
  {
    id: "507f1f77bcf86cd799439011",
    start_time: "2026-02-04T08:00:00Z",
    cashier_name: "Nguyễn Văn A",
    status: "OPEN"
  },
  // ...
]
```

### Payment Data

**Source:** `cashierStore.payments`

```javascript
[
  {
    order_id: "507f1f77bcf86cd799439012",
    order_number: "20260204-100530-123",
    customer_name: "Trần Thị B",
    amount: 50000,
    payment_method: "CASH",
    status: "PAID",
    paid_at: "2026-02-04T10:05:30Z"
  },
  // ...
]
```

---

## 🎨 UI States

### State 1: Initial Load (No Shift Selected)

```
┌─────────────────────────────────────────┐
│  📅 Chọn ca thu ngân để xem            │
│  [-- Chọn ca thu ngân --]             │
├─────────────────────────────────────────┤
│  💳 Danh sách thanh toán               │
│  0 giao dịch                           │
│                                         │
│         📭                              │
│  Chưa có thanh toán nào                │
│  Chọn ca làm việc để xem               │
└─────────────────────────────────────────┘
```

**Condition:** `payments.length === 0`

### State 2: Shift Selected, Has Payments

```
┌─────────────────────────────────────────┐
│  📅 Chọn ca thu ngân để xem            │
│  [04/02/2026 - Nguyễn Văn A - 🟢...]  │
├─────────────────────────────────────────┤
│  💳 Danh sách thanh toán               │
│  3 giao dịch                           │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │ Trần Thị B          50,000₫      │ │
│  │ 04/02 - 10:05       💵 Tiền mặt  │ │
│  │ ✓ Đã thu                          │ │
│  │ [✏️ Điều chỉnh] [⚠️ Báo lỗi] [🔒]│ │
│  └───────────────────────────────────┘ │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │ Nguyễn Văn C        35,000₫      │ │
│  │ 04/02 - 10:15       💳 CK        │ │
│  │ ✓ Đã thu                          │ │
│  │ [✏️ Điều chỉnh] [⚠️ Báo lỗi] [🔒]│ │
│  └───────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

**Condition:** `payments.length > 0`

### State 3: Loading

```
┌─────────────────────────────────────────┐
│  💳 Danh sách thanh toán               │
│                                         │
│         🔄 (spinning)                   │
│         Đang tải...                     │
└─────────────────────────────────────────┘
```

**Condition:** `loading === true`

---

## 🔍 Backend API

### Get Payments by Shift

**Endpoint:** `GET /api/cashier/shifts/:id/payments`

**Handler:** `cashierHandler.GetPaymentsByShift()`

**Service:** `paymentOversightService.GetPaymentsByShift()`

**Logic:**
```go
func (s *PaymentOversightService) GetPaymentsByShift(shiftID string) ([]*PaymentSummary, error) {
    // 1. Convert shiftID to ObjectID
    shiftObjID, _ := primitive.ObjectIDFromHex(shiftID)
    
    // 2. Find all orders in shift
    orders, _ := s.orderRepo.FindByShiftID(ctx, shiftObjID)
    
    // 3. Filter orders with payment status
    var payments []*PaymentSummary
    for _, ord := range orders {
        if ord.Status == PAID || ord.Status == IN_PROGRESS || ord.Status == SERVED {
            payments = append(payments, &PaymentSummary{
                OrderID:       ord.ID.Hex(),
                OrderNumber:   ord.OrderNumber,
                Amount:        ord.Total,
                PaymentMethod: ord.PaymentMethod,
                Status:        ord.Status,
                PaidAt:        *ord.PaidAt,
            })
        }
    }
    
    return payments, nil
}
```

**Filtered Statuses:**
- ✅ `PAID` - Đã thanh toán
- ✅ `IN_PROGRESS` - Đang pha (đã thanh toán)
- ✅ `SERVED` - Đã phục vụ (đã thanh toán)
- ❌ `CREATED` - Chưa thanh toán (không hiển thị)
- ❌ `QUEUED` - Chờ pha (không hiển thị nếu chưa thanh toán)

---

## 🐛 Troubleshooting

### Vấn đề 1: Không thấy dropdown shifts

**Nguyên nhân:**
- Chưa có cashier shifts nào
- User không phải cashier
- API call failed

**Giải pháp:**
1. Kiểm tra role: `localStorage.getItem('role')` === 'cashier'
2. Kiểm tra API: `GET /api/cashier/shifts/my-cashier-shifts`
3. Tạo cashier shift mới nếu chưa có

### Vấn đề 2: Chọn shift nhưng không có payments

**Nguyên nhân:**
- Shift chưa có orders nào
- Orders chưa được thanh toán (status = CREATED)
- API call failed

**Giải pháp:**
1. Kiểm tra shift có orders: `GET /api/cashier/shifts/:id/status`
2. Tạo orders và thanh toán
3. Kiểm tra console log errors

### Vấn đề 3: Empty state hiển thị mãi

**Nguyên nhân:**
- `payments.length === 0`
- API không trả về data
- Filter logic loại bỏ tất cả orders

**Debug:**
```javascript
// Check in browser console
console.log('Selected Shift:', selectedShift.value)
console.log('Payments:', cashierStore.payments)
console.log('Loading:', cashierStore.loading)
console.log('Error:', cashierStore.error)
```

---

## 📝 Code References

### Template (Lines 290-360)

```vue
<!-- Payment List -->
<div class="mb-4">
  <div class="flex items-center justify-between mb-3">
    <h2 class="text-lg font-bold text-gray-800">💳 Danh sách thanh toán</h2>
    <span class="text-sm text-gray-600">{{ payments.length }} giao dịch</span>
  </div>

  <!-- Empty State -->
  <div v-if="payments.length === 0" class="text-center py-12 bg-white rounded-2xl">
    <div class="text-5xl mb-3">📭</div>
    <p class="text-gray-500">Chưa có thanh toán nào</p>
    <p class="text-sm text-gray-400 mt-1">Chọn ca làm việc để xem</p>
  </div>

  <!-- Payment Cards -->
  <div v-else class="space-y-3">
    <div v-for="payment in payments" :key="payment.order_id">
      <!-- Payment card content -->
    </div>
  </div>
</div>
```

### Script (Lines 389-438)

```javascript
const selectedShift = ref('')
const payments = computed(() => cashierStore.payments)

const loadPayments = async () => {
  if (selectedShift.value) {
    await Promise.all([
      cashierStore.getShiftStatus(selectedShift.value),
      cashierStore.getPaymentsByShift(selectedShift.value)
    ])
  }
}

onMounted(async () => {
  await cashierShiftStore.fetchMyCashierShifts()
  await cashierStore.getPendingDiscrepancies()
  await cashierStore.fetchPendingHandovers()
})
```

---

## ✅ Summary

**Danh sách thanh toán hiển thị khi:**

1. ✅ User đã login với role **cashier**
2. ✅ Đã chọn **shift** từ dropdown
3. ✅ Shift có **orders đã thanh toán** (PAID/IN_PROGRESS/SERVED)
4. ✅ API call thành công

**Không hiển thị khi:**

1. ❌ Chưa chọn shift
2. ❌ Shift không có orders
3. ❌ Orders chưa thanh toán (CREATED)
4. ❌ API call failed

**Empty state hiển thị khi:**
- `payments.length === 0` (chưa chọn shift HOẶC shift không có payments)

---

**Last Updated:** 2026-02-04  
**Related Files:**
- `frontend/src/views/CashierDashboard.vue`
- `frontend/src/stores/cashier.js`
- `backend/application/services/payment_oversight_service.go`
