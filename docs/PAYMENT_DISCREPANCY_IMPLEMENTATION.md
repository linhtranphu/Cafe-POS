# 📋 Payment Discrepancy Implementation Status

**Date:** 2026-02-04  
**Status:** ✅ **FULLY IMPLEMENTED**

---

## 🎯 Overview

Payment Discrepancy là tính năng cho phép cashier báo cáo và theo dõi các sai lệch/vấn đề liên quan đến thanh toán của orders.

---

## ✅ Implementation Checklist

### Backend (Go)

- [x] **Domain Model** - `backend/domain/cashier/cash_reconciliation.go`
  - [x] `PaymentDiscrepancy` struct
  - [x] Status constants (PENDING, RESOLVED)
  - [x] `Resolve()` method

- [x] **Repository** - `backend/infrastructure/mongodb/payment_discrepancy_repository.go`
  - [x] `Create()` - Tạo discrepancy mới
  - [x] `FindByOrderID()` - Tìm theo order
  - [x] `FindPendingDiscrepancies()` - Lấy discrepancies chưa giải quyết
  - [x] `UpdateStatus()` - Cập nhật status
  - [x] `FindByShiftID()` - Tìm theo shift

- [x] **Service** - `backend/application/services/payment_oversight_service.go`
  - [x] `ReportDiscrepancy()` - Báo cáo sai lệch
  - [x] `GetPendingDiscrepancies()` - Lấy danh sách pending
  - [x] `ResolveDiscrepancy()` - Giải quyết sai lệch

- [x] **Handler** - `backend/interfaces/http/cashier_handler.go`
  - [x] `ReportDiscrepancy()` - POST endpoint
  - [x] `GetPendingDiscrepancies()` - GET endpoint
  - [x] `ResolveDiscrepancy()` - POST endpoint

- [x] **Routes** - `backend/main.go`
  ```go
  cashier.POST("/discrepancies", cashierHandler.ReportDiscrepancy)
  cashier.GET("/discrepancies/pending", cashierHandler.GetPendingDiscrepancies)
  cashier.POST("/discrepancies/:id/resolve", cashierHandler.ResolveDiscrepancy)
  ```

### Frontend (Vue.js)

- [x] **Service** - `frontend/src/services/cashier.js`
  - [x] `reportDiscrepancy(data)` - API call
  - [x] `getPendingDiscrepancies()` - API call
  - [x] `resolveDiscrepancy(discrepancyId)` - API call

- [x] **Store** - `frontend/src/stores/cashier.js`
  - [x] State: `discrepancies[]`
  - [x] Getter: `pendingDiscrepancies`
  - [x] Action: `reportDiscrepancy()`
  - [x] Action: `getPendingDiscrepancies()`
  - [x] Action: `resolveDiscrepancy()`

- [x] **Component** - `frontend/src/components/DiscrepancyModal.vue`
  - [x] Form để nhập lý do và số tiền
  - [x] Validation
  - [x] Emit confirm event

- [x] **View** - `frontend/src/views/CashierDashboard.vue`
  - [x] Button "⚠️ Báo lỗi" trên mỗi payment
  - [x] Alert hiển thị pending discrepancies
  - [x] Expandable list để xem chi tiết
  - [x] Button "✓ Giải quyết" cho mỗi discrepancy
  - [x] Integration với modal

---

## 📊 Data Flow

### 1. Report Discrepancy

```
User clicks "⚠️ Báo lỗi" on payment
  ↓
DiscrepancyModal opens
  ↓
User enters reason & amount
  ↓
Click "Báo cáo"
  ↓
handleReportDiscrepancy() called
  ↓
cashierStore.reportDiscrepancy({
  order_id: payment.order_id,
  reason: "...",
  amount: 5000
})
  ↓
POST /api/cashier/discrepancies
  ↓
PaymentOversightService.ReportDiscrepancy()
  ↓
Create PaymentDiscrepancy with status=PENDING
  ↓
Save to MongoDB collection "payment_discrepancies"
  ↓
Refresh pending discrepancies
  ↓
Dashboard shows alert
```

### 2. View Pending Discrepancies

```
CashierDashboard mounted
  ↓
cashierStore.getPendingDiscrepancies()
  ↓
GET /api/cashier/discrepancies/pending
  ↓
PaymentOversightService.GetPendingDiscrepancies()
  ↓
Query MongoDB: { status: "PENDING" }
  ↓
Return array of discrepancies
  ↓
Display in dashboard alert (yellow box)
  ↓
User clicks "Xem chi tiết"
  ↓
Expandable list shows all discrepancies
```

### 3. Resolve Discrepancy

```
User clicks "✓ Giải quyết" on discrepancy
  ↓
Confirm dialog
  ↓
cashierStore.resolveDiscrepancy(discrepancyId)
  ↓
POST /api/cashier/discrepancies/:id/resolve
  ↓
PaymentOversightService.ResolveDiscrepancy()
  ↓
Update status to "RESOLVED"
  ↓
Set resolved_at timestamp
  ↓
Refresh pending discrepancies
  ↓
Discrepancy removed from alert
```

---

## 🎨 UI Components

### CashierDashboard Alert (Lines 149-177)

```vue
<!-- Pending Discrepancies Alert -->
<div v-if="pendingDiscrepancies.length > 0" 
  class="bg-yellow-50 border-2 border-yellow-300 rounded-2xl p-4 mb-4">
  <div class="flex items-center gap-3 mb-3">
    <span class="text-2xl">⚠️</span>
    <div>
      <h3 class="font-bold text-yellow-800">Sai lệch cần xử lý</h3>
      <p class="text-sm text-yellow-700">
        {{ pendingDiscrepancies.length }} sai lệch đang chờ
      </p>
    </div>
  </div>
  <button @click="showDiscrepancyList = !showDiscrepancyList">
    {{ showDiscrepancyList ? 'Ẩn' : 'Xem chi tiết' }} →
  </button>
</div>
```

### Discrepancy List (Lines 179-206)

```vue
<!-- Expandable list showing each discrepancy -->
<div v-if="showDiscrepancyList && pendingDiscrepancies.length > 0">
  <div v-for="discrepancy in pendingDiscrepancies" :key="discrepancy.id">
    <h4>Order #{{ discrepancy.order_id?.slice(-6) }}</h4>
    <p>{{ discrepancy.reason }}</p>
    <span>{{ formatPrice(discrepancy.amount) }}</span>
    <button @click="resolveDiscrepancy(discrepancy.id)">
      ✓ Giải quyết
    </button>
  </div>
</div>
```

### Payment Actions (Lines 330-350)

```vue
<!-- Each payment has "Báo lỗi" button -->
<button @click="showDiscrepancyModal(payment)">
  ⚠️ Báo lỗi
</button>
```

---

## 🔧 API Endpoints

### 1. Report Discrepancy

**Endpoint:** `POST /api/cashier/discrepancies`

**Auth:** Required (Cashier role)

**Request Body:**
```json
{
  "order_id": "507f1f77bcf86cd799439011",
  "reason": "Khách hàng khiếu nại thiếu tiền thừa",
  "amount": 5000
}
```

**Response:**
```json
{
  "message": "Discrepancy reported successfully"
}
```

### 2. Get Pending Discrepancies

**Endpoint:** `GET /api/cashier/discrepancies/pending`

**Auth:** Required (Cashier role)

**Response:**
```json
{
  "discrepancies": [
    {
      "id": "507f1f77bcf86cd799439012",
      "order_id": "507f1f77bcf86cd799439011",
      "cashier_id": "507f1f77bcf86cd799439010",
      "reason": "Khách hàng khiếu nại thiếu tiền thừa",
      "amount": 5000,
      "status": "PENDING",
      "reported_at": "2026-02-04T10:30:00Z",
      "created_at": "2026-02-04T10:30:00Z",
      "updated_at": "2026-02-04T10:30:00Z"
    }
  ]
}
```

### 3. Resolve Discrepancy

**Endpoint:** `POST /api/cashier/discrepancies/:id/resolve`

**Auth:** Required (Cashier role)

**Response:**
```json
{
  "message": "Discrepancy resolved successfully"
}
```

---

## 💾 Database Schema

**Collection:** `payment_discrepancies`

```javascript
{
  _id: ObjectId,
  order_id: String,           // Reference to order
  cashier_id: String,         // Cashier who reported
  reason: String,             // Description of issue
  amount: Number,             // Amount involved
  status: String,             // "PENDING" or "RESOLVED"
  reported_at: Date,          // When reported
  resolved_at: Date,          // When resolved (nullable)
  created_at: Date,
  updated_at: Date
}
```

**Indexes:**
- `order_id` - For finding discrepancies by order
- `status` - For filtering pending discrepancies
- `cashier_id` - For filtering by cashier

---

## 🧪 Testing

### Manual Test Flow

1. **Login as Cashier**
   - Navigate to `/cashier`

2. **Select a Shift**
   - Choose a cashier shift from dropdown
   - View payments list

3. **Report Discrepancy**
   - Click "⚠️ Báo lỗi" on any payment
   - Enter reason: "Test discrepancy"
   - Enter amount: 5000
   - Click "Báo cáo"
   - ✅ Should see success message
   - ✅ Should see yellow alert appear

4. **View Discrepancies**
   - ✅ Alert shows "1 sai lệch đang chờ"
   - Click "Xem chi tiết"
   - ✅ Should see discrepancy details

5. **Resolve Discrepancy**
   - Click "✓ Giải quyết"
   - Confirm dialog
   - ✅ Discrepancy removed from list
   - ✅ Alert disappears

### API Test

```bash
# Get pending discrepancies (requires auth token)
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/cashier/discrepancies/pending

# Expected: {"discrepancies": [...]}
```

---

## 🆚 Comparison: Payment vs Handover Discrepancy

| Feature | Payment Discrepancy | Handover Discrepancy |
|---------|-------------------|---------------------|
| **Purpose** | Track payment issues | Track cash handover differences |
| **Reported by** | Cashier | Cashier (during reconciliation) |
| **Affects** | Individual orders | Shift cash handover |
| **Approval needed** | No | Yes (Manager) |
| **Status flow** | PENDING → RESOLVED | PENDING → CONFIRMED → APPROVED |
| **Collection** | `payment_discrepancies` | `cash_discrepancies` |
| **UI Color** | Yellow (⚠️) | Red (🚨) |
| **Resolution** | Cashier resolves | Manager approves |

---

## 📝 Use Cases

### Use Case 1: Customer Complaint
**Scenario:** Customer says they didn't receive correct change

**Flow:**
1. Cashier finds the order in payment list
2. Clicks "⚠️ Báo lỗi"
3. Enters: "Khách khiếu nại thiếu tiền thừa 5,000đ"
4. Amount: 5000
5. Reports discrepancy
6. Investigates and resolves issue
7. Marks as resolved

### Use Case 2: Wrong Payment Method
**Scenario:** Order recorded as CASH but customer paid by transfer

**Flow:**
1. Cashier notices error
2. Reports discrepancy: "Sai phương thức - thực tế là chuyển khoản"
3. Uses "✏️ Điều chỉnh" to override payment
4. Re-collects payment with correct method
5. Resolves discrepancy

### Use Case 3: Duplicate Payment
**Scenario:** Customer accidentally paid twice

**Flow:**
1. Cashier reports discrepancy: "Thanh toán trùng lặp"
2. Amount: [duplicate amount]
3. Processes refund
4. Marks discrepancy as resolved

---

## 🔮 Future Enhancements

### Potential Improvements

1. **Discrepancy Categories**
   - Add predefined categories (Wrong amount, Wrong method, Duplicate, etc.)
   - Easier reporting and analytics

2. **Auto-Resolution**
   - Some discrepancies could auto-resolve when order is overridden
   - Link discrepancy to override action

3. **Manager Review**
   - High-value discrepancies require manager approval
   - Threshold-based escalation

4. **Analytics Dashboard**
   - Track discrepancy trends
   - Identify problematic patterns
   - Cashier performance metrics

5. **Notifications**
   - Alert manager when discrepancy reported
   - Remind cashier of pending discrepancies

6. **Audit Trail**
   - Link to payment audit records
   - Full history of actions taken

---

## 📚 Related Documentation

- [CASHIER_IMPLEMENTATION.md](./CASHIER_IMPLEMENTATION.md) - Overall cashier system
- [CASH_HANDOVER_COMPLETE_SUMMARY.md](./CASH_HANDOVER_COMPLETE_SUMMARY.md) - Handover discrepancies
- [ORDER_IMPLEMENTATION.md](./ORDER_IMPLEMENTATION.md) - Order and payment flow

---

## ✅ Conclusion

Payment Discrepancy feature is **FULLY IMPLEMENTED** and **PRODUCTION READY**.

**Components:**
- ✅ Backend API (3 endpoints)
- ✅ Database schema
- ✅ Frontend service
- ✅ Store management
- ✅ UI components
- ✅ Modal for reporting
- ✅ Dashboard alerts
- ✅ Resolution workflow

**Status:** Ready for use in production environment.

---

**Last Updated:** 2026-02-04  
**Verified By:** Development Team  
**Implementation Time:** Already complete
