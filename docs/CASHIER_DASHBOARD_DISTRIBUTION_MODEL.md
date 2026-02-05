# 📋 Cashier Dashboard - Distribution Model Design

**Date:** 2026-02-04  
**Context:** Cashier dashboard trong mô hình distribution

---

## 🎯 Mô hình Distribution

### Vai trò và trách nhiệm

```
┌─────────────────────────────────────────┐
│           WAITER (Phục vụ)              │
├─────────────────────────────────────────┤
│ • Tạo orders                            │
│ • Thu tiền từ khách                     │
│ • Giữ tiền trong ca của mình            │
│ • Bàn giao tiền cho cashier             │
└─────────────────────────────────────────┘
                  ↓
            [HANDOVER]
                  ↓
┌─────────────────────────────────────────┐
│          CASHIER (Thu ngân)             │
├─────────────────────────────────────────┤
│ • KHÔNG tạo orders trực tiếp            │
│ • KHÔNG thu tiền từ khách               │
│ • Nhận tiền qua handover                │
│ • Giám sát payments của waiters         │
│ • Xử lý discrepancies                   │
│ • Quản lý ca thu ngân riêng             │
└─────────────────────────────────────────┘
```

---

## 🎨 Cashier Dashboard Sections

### 1. ✅ Handover Notifications (Hiển thị)

**Mục đích:** Thông báo các yêu cầu bàn giao chờ xử lý

**Hiển thị khi:**
- Có handovers với status PENDING
- Có handovers với discrepancy chưa approve

**Lý do:** Đây là chức năng chính của cashier trong mô hình distribution

```vue
<!-- Handover Notifications -->
<div v-if="pendingHandovers.length > 0">
  {{ pendingHandovers.length }} yêu cầu bàn giao đang chờ
</div>

<!-- Handover Discrepancy Alert -->
<div v-if="handoversWithDiscrepancy.length > 0">
  {{ handoversWithDiscrepancy.length }} bàn giao có chênh lệch
</div>
```

### 2. ✅ Shift Selector (Hiển thị - Đã sửa)

**Mục đích:** Chọn waiter shift để xem payments

**Hiển thị:** Dropdown với **waiter shifts** (không phải cashier shifts)

**Lý do:** 
- Orders thuộc về waiter shifts
- Cashier cần xem payments từ waiter shifts để giám sát

```vue
<!-- Shift Selector - Waiter Shifts -->
<label>📅 Chọn ca phục vụ để xem thanh toán</label>
<p>💡 Chọn ca của waiter để xem các thanh toán trong ca đó</p>
<select v-model="selectedShift">
  <option v-for="shift in waiterShifts">
    {{ shift.waiter_name }} ({{ shift.role }})
  </option>
</select>
```

### 3. ✅ Shift Status Card (Hiển thị - Đã cập nhật label)

**Mục đích:** Hiển thị thống kê của waiter shift đang xem

**Hiển thị khi:** `shiftStatus !== null` (đã chọn shift)

**Lý do:** Thông tin hữu ích để cashier biết tổng quan shift đang giám sát

**Changes:**
- ✅ Title: "📊 Tổng quan ca phục vụ" (thay vì "ca làm")
- ✅ Subtitle: "Thống kê của ca đang xem"

```vue
<!-- Shift Status Card -->
<div v-if="shiftStatus">
  <h2>📊 Tổng quan ca phục vụ</h2>
  <p>Thống kê của ca đang xem</p>
  <div>Tổng đơn: {{ shiftStatus.total_orders }}</div>
  <div>Doanh thu: {{ shiftStatus.total_revenue }}</div>
</div>
```

### 4. ❌ Cash Reconciliation (ẨN)

**Mục đích:** Đối soát tiền mặt

**Hiển thị:** `v-if="false"` - DISABLED

**Lý do:**
1. **Selected shift là WAITER shift**, không phải cashier shift
2. **Cashier không giữ tiền** từ waiter shift
3. **Đối soát xảy ra trong handover flow**, không phải ở đây
4. Cashier không thể đối soát tiền mà họ không cầm

**Trong mô hình distribution:**
- Waiter giữ tiền → Waiter tự đối soát khi bàn giao
- Cashier nhận tiền qua handover → Đối soát trong handover confirmation
- Cashier KHÔNG đối soát tiền của waiter shift

```vue
<!-- Cash Reconciliation - DISABLED -->
<!-- 
  In distribution model:
  - Waiter holds cash in their shift
  - Cashier receives cash via handover
  - Reconciliation happens during handover confirmation
-->
<div v-if="false">
  <!-- Hidden section -->
</div>
```

### 5. ✅ Payment List (Hiển thị)

**Mục đích:** Hiển thị danh sách payments từ waiter shift

**Hiển thị khi:** `payments.length > 0`

**Lý do:** Đây là chức năng giám sát chính của cashier

**Actions available:**
- ✏️ Điều chỉnh (Override payment)
- ⚠️ Báo lỗi (Report discrepancy)
- 🔒 Khóa order (Lock order)

```vue
<!-- Payment List -->
<div v-if="payments.length > 0">
  <div v-for="payment in payments">
    <button @click="showOverrideModal(payment)">✏️ Điều chỉnh</button>
    <button @click="showDiscrepancyModal(payment)">⚠️ Báo lỗi</button>
    <button @click="lockOrder(payment.order_id)">🔒</button>
  </div>
</div>
```

---

## 🔄 User Flow

### Cashier Workflow

```
1. Login as cashier
   ↓
2. Navigate to /cashier
   ↓
3. See dashboard with:
   - Handover notifications (if any)
   - Shift selector (waiter shifts)
   - Empty payment list
   ↓
4. Select a waiter shift from dropdown
   ↓
5. View payments from that shift
   ↓
6. Monitor and take actions:
   - Override incorrect payments
   - Report discrepancies
   - Lock problematic orders
   ↓
7. Process handovers from waiters
   ↓
8. Receive cash via handover confirmation
```

---

## 🆚 Comparison: Waiter vs Cashier Dashboard

| Feature | Waiter Dashboard | Cashier Dashboard |
|---------|-----------------|-------------------|
| **Create Orders** | ✅ Yes | ❌ No |
| **Collect Payments** | ✅ Yes | ❌ No |
| **Hold Cash** | ✅ Yes | ❌ No (receives via handover) |
| **View Own Orders** | ✅ Yes | ❌ No |
| **View Others' Orders** | ❌ No | ✅ Yes (monitoring) |
| **Handover Cash** | ✅ Yes (to cashier) | ❌ No |
| **Receive Handovers** | ❌ No | ✅ Yes (from waiters) |
| **Reconcile Own Cash** | ✅ Yes (during handover) | ❌ No (not applicable) |
| **Override Payments** | ❌ No | ✅ Yes |
| **Report Discrepancies** | ❌ No | ✅ Yes |
| **Lock Orders** | ❌ No | ✅ Yes |

---

## 📊 Data Ownership

### Waiter Shift Data

```javascript
{
  id: "waiter_shift_id",
  role: "waiter",
  waiter_id: "...",
  waiter_name: "Lê Văn C",
  
  // Waiter owns this data
  current_cash: 150000,      // Cash waiter is holding
  remaining_cash: 150000,    // Cash not yet handed over
  total_revenue: 150000,     // Total collected
  
  // Orders belong to waiter shift
  orders: [
    { id: "order1", total: 50000, payment_method: "CASH" },
    { id: "order2", total: 100000, payment_method: "CASH" }
  ]
}
```

### Cashier Shift Data

```javascript
{
  id: "cashier_shift_id",
  role: "cashier",
  cashier_id: "...",
  cashier_name: "Nguyễn Văn A",
  
  // Cashier owns this data
  received_cash: 0,          // Cash received via handovers
  total_handovers: 0,        // Number of handovers processed
  
  // NO orders directly
  orders: []                 // Empty - cashier doesn't create orders
}
```

---

## 🔍 Why This Design?

### 1. Clear Accountability

- **Waiter** responsible for their orders and cash
- **Cashier** responsible for receiving and securing cash
- Clear audit trail of who collected what

### 2. Separation of Concerns

- **Waiter** focuses on customer service and order fulfillment
- **Cashier** focuses on cash management and oversight
- No confusion about responsibilities

### 3. Security

- Cash changes hands through formal handover process
- Both parties confirm amounts
- Discrepancies are tracked and resolved
- Manager approval for significant discrepancies

### 4. Flexibility

- Cashier can monitor any waiter shift
- Multiple waiters can work simultaneously
- Cashier provides oversight without interfering

---

## 🐛 Common Misunderstandings

### ❌ Misconception 1: Cashier Reconciles Waiter Cash

**Wrong:** Cashier đối soát tiền mặt của waiter shift

**Right:** Cashier chỉ xác nhận số tiền nhận được qua handover

**Why:** Cashier không cầm tiền của waiter cho đến khi handover xảy ra

### ❌ Misconception 2: Cashier Shift Has Orders

**Wrong:** Chọn cashier shift để xem orders

**Right:** Chọn waiter shift để xem orders

**Why:** Orders thuộc về waiter shift, không phải cashier shift

### ❌ Misconception 3: Cashier Creates Orders

**Wrong:** Cashier có thể tạo orders

**Right:** Chỉ waiter tạo orders

**Why:** Trong mô hình distribution, cashier không tương tác trực tiếp với khách

---

## ✅ Summary

### Cashier Dashboard Purpose

**Primary Functions:**
1. ✅ Monitor payments from waiter shifts
2. ✅ Process handovers from waiters
3. ✅ Handle payment discrepancies
4. ✅ Provide oversight and control

**NOT for:**
1. ❌ Creating orders
2. ❌ Collecting payments directly
3. ❌ Reconciling waiter's cash
4. ❌ Managing waiter shifts

### Key Changes Made

1. ✅ Shift selector shows **waiter shifts** (not cashier shifts)
2. ✅ Added helper text explaining to select waiter shift
3. ✅ Updated labels: "ca phục vụ" instead of "ca thu ngân"
4. ✅ Hidden cash reconciliation section (not applicable)
5. ✅ Updated shift status card title to clarify it's waiter shift stats

---

**Last Updated:** 2026-02-04  
**Model:** Distribution  
**Status:** Implemented and documented
