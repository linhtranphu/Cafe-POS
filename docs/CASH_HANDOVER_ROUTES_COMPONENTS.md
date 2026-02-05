# 🗺️ Cash Handover - Routes & Components Reference

## 📍 Routes

### Frontend Routes

| Route | Component | Role | Description |
|-------|-----------|------|-------------|
| `/shift` | `ShiftView.vue` | Waiter | Waiter shift management & handover creation |
| `/cashier/handovers` | `CashierHandoverView.vue` | Cashier, Manager | Handover management & approval |
| `/cashier` | `CashierDashboard.vue` | Cashier | Dashboard with handover notifications |

### Backend API Routes

#### Waiter Routes (Shift-based)
```
POST   /api/shifts/:id/handover
POST   /api/shifts/:id/handover-and-end
GET    /api/shifts/:id/pending-handover
GET    /api/shifts/:id/handovers
```

#### Cashier Routes
```
GET    /api/cash-handovers/pending
GET    /api/cash-handovers/all-pending
GET    /api/cash-handovers/today
POST   /api/cash-handovers/:id/confirm
POST   /api/cash-handovers/:id/quick-confirm
```

#### Manager Routes
```
GET    /api/cash-handovers/pending-approvals
POST   /api/cash-handovers/:id/approve-discrepancy
GET    /api/cash-handovers/discrepancy-stats
```

---

## 🧩 Components

### Vue Components

#### 1. ShiftView.vue
**Location:** `frontend/src/views/ShiftView.vue`

**Purpose:** Waiter interface for shift management and handover creation

**Key Sections:**
- Shift info display
- Partial handover modal
- Handover & end shift modal
- Pending handover status
- Handover history

**Key Methods:**
```javascript
createPartialHandover()
createHandoverAndEndShift()
cancelHandover()
fetchPendingHandover()
fetchHandoverHistory()
```

**Store Dependencies:**
- `useShiftStore()` - Shift management
- `useAuthStore()` - User authentication

---

#### 2. CashierHandoverView.vue
**Location:** `frontend/src/views/CashierHandoverView.vue`

**Purpose:** Cashier/Manager interface for handover management

**Key Sections:**
- Tab navigation (Pending, Today, Approvals)
- Pending handovers list
- Today's handovers list
- Reconciliation modal
- Approval modal (manager only)
- Discrepancy stats (manager only)

**Key Methods:**
```javascript
fetchPendingHandovers()
fetchTodayHandovers()
fetchPendingApprovals()
confirmHandover()
quickConfirm()
approveDiscrepancy()
getDiscrepancyStats()
```

**Store Dependencies:**
- `useCashierStore()` - Cashier operations
- `useAuthStore()` - User authentication & role check

---

#### 3. CashierDashboard.vue
**Location:** `frontend/src/views/CashierDashboard.vue`

**Purpose:** Cashier dashboard with handover notifications

**Key Sections:**
- Handover notification banner
- Quick handover actions
- Pending handovers preview (top 3)

**Key Methods:**
```javascript
fetchPendingHandovers()
quickConfirm()
navigateToHandovers()
```

**Store Dependencies:**
- `useCashierStore()` - Handover data
- `useCashierShiftStore()` - Shift data

---

## 🏪 Stores

### 1. Shift Store
**Location:** `frontend/src/stores/shift.js`

**Handover-related State:**
```javascript
{
  pendingHandover: null,
  handoverHistory: []
}
```

**Handover-related Actions:**
```javascript
createPartialHandover(shiftId, data)
createHandoverAndEndShift(shiftId, data)
cancelHandover(handoverId)
fetchPendingHandover(shiftId)
fetchHandoverHistory(shiftId)
```

---

### 2. Cashier Store
**Location:** `frontend/src/stores/cashier.js`

**Handover-related State:**
```javascript
{
  pendingHandovers: [],
  todayHandovers: [],
  pendingApprovals: [],
  discrepancyStats: null
}
```

**Handover-related Actions:**
```javascript
fetchPendingHandovers()
fetchAllPendingHandovers()
fetchTodayHandovers()
confirmHandover(handoverId, data)
quickConfirm(handoverId, status)
fetchPendingApprovals()
approveDiscrepancy(handoverId, data)
getDiscrepancyStats(startDate, endDate)
```

---

## 🔌 Services

### Handover Service
**Location:** `frontend/src/services/handover.js`

**API Methods:**
```javascript
// Waiter APIs
createHandover(shiftId, data)
createHandoverAndEndShift(shiftId, data)
getPendingHandover(shiftId)
getHandoverHistory(shiftId)
cancelHandover(handoverId)

// Cashier APIs
getPendingHandovers()
getAllPendingHandovers()
getTodayHandovers()
confirmHandover(handoverId, data)
quickConfirm(handoverId, status)

// Manager APIs
getPendingApprovals()
approveDiscrepancy(handoverId, data)
getDiscrepancyStats(startDate, endDate)
```

---

## 🎯 Navigation Flow

### Waiter Flow
```
Login (waiter)
  ↓
/shift (ShiftView.vue)
  ↓
Click "💰 Bàn giao tiền"
  ↓
Modal: Partial Handover
  ↓
Submit → API: POST /api/shifts/:id/handover
  ↓
Success → Show pending status
```

### Cashier Flow
```
Login (cashier)
  ↓
/cashier (CashierDashboard.vue)
  ↓
See notification: "⚠️ 3 yêu cầu bàn giao"
  ↓
Click "Xem ngay"
  ↓
/cashier/handovers (CashierHandoverView.vue)
  ↓
Tab: Đang chờ
  ↓
Click "Xác nhận chi tiết"
  ↓
Modal: Reconciliation
  ↓
Submit → API: POST /api/cash-handovers/:id/confirm
  ↓
Success → Update list
```

### Manager Flow
```
Login (manager)
  ↓
/cashier/handovers (CashierHandoverView.vue)
  ↓
Tab: Cần duyệt (2)
  ↓
See discrepancy list
  ↓
Click "Chi tiết"
  ↓
Modal: Approval
  ↓
Submit → API: POST /api/cash-handovers/:id/approve-discrepancy
  ↓
Success → Update list
```

---

## 🔐 Role-Based Access

### Waiter
- ✅ Create partial handover
- ✅ Create handover & end shift
- ✅ View pending handover
- ✅ View handover history
- ✅ Cancel pending handover
- ❌ Confirm handover
- ❌ Approve discrepancy

### Cashier
- ❌ Create handover
- ✅ View pending handovers (assigned to them)
- ✅ View all pending handovers
- ✅ View today's handovers
- ✅ Confirm handover with reconciliation
- ✅ Quick confirm/reject
- ❌ Approve discrepancy (if > threshold)

### Manager
- ❌ Create handover
- ✅ View all pending handovers
- ✅ View today's handovers
- ✅ Confirm handover
- ✅ Approve/reject discrepancy
- ✅ View discrepancy statistics
- ✅ Full access to all handover data

---

## 📦 Data Models

### Handover Object (Frontend)
```javascript
{
  id: "507f1f77bcf86cd799439011",
  shift_id: "507f1f77bcf86cd799439012",
  waiter_id: "507f1f77bcf86cd799439013",
  waiter_name: "Nguyễn Văn A",
  cashier_id: "507f1f77bcf86cd799439014",
  cashier_name: "Trần Thị B",
  declared_amount: 500000,
  actual_amount: 500000,
  discrepancy_amount: 0,
  status: "PENDING", // PENDING, CONFIRMED, REJECTED, REQUIRES_APPROVAL
  waiter_note: "Bàn giao tiền ca chiều",
  cashier_note: "",
  manager_note: "",
  created_at: "2026-02-04T14:30:00Z",
  confirmed_at: null,
  approved_at: null
}
```

### Discrepancy Stats Object
```javascript
{
  total_handovers: 150,
  handovers_with_discrepancy: 12,
  total_discrepancy_amount: -500000,
  average_discrepancy: -41667,
  max_discrepancy: -200000,
  approved_count: 10,
  rejected_count: 2
}
```

---

## 🎨 UI States

### Handover Status Colors

| Status | Color | Badge | Description |
|--------|-------|-------|-------------|
| PENDING | Yellow | ⏳ | Waiting for cashier confirmation |
| CONFIRMED | Green | ✅ | Confirmed by cashier, no discrepancy |
| REJECTED | Red | ❌ | Rejected by cashier |
| REQUIRES_APPROVAL | Orange | ⚠️ | Discrepancy > threshold, needs manager approval |

### Button States

| Action | Color | Icon | Condition |
|--------|-------|------|-----------|
| Create Handover | Blue | 💰 | Shift is ACTIVE |
| End Shift | Red | 🏁 | Shift is ACTIVE |
| Confirm | Green | ✅ | Handover is PENDING |
| Reject | Red | ❌ | Handover is PENDING |
| Approve | Green | ✅ | Handover REQUIRES_APPROVAL |
| Cancel | Gray | ❌ | Handover is PENDING (waiter only) |

---

## 🔄 State Transitions

### Handover State Machine

```
PENDING
  ↓
  ├─→ CONFIRMED (if discrepancy <= threshold)
  ├─→ REJECTED (if cashier rejects)
  └─→ REQUIRES_APPROVAL (if discrepancy > threshold)
       ↓
       ├─→ CONFIRMED (if manager approves)
       └─→ REJECTED (if manager rejects)
```

### Shift State with Handover

```
ACTIVE
  ↓
  ├─→ (Partial Handover) → ACTIVE (continues)
  └─→ (Handover & End) → ENDING
       ↓
       └─→ (Cashier confirms) → ENDED
```

---

## 📊 API Response Examples

### Create Handover Response
```json
{
  "handover": {
    "id": "507f1f77bcf86cd799439011",
    "shift_id": "507f1f77bcf86cd799439012",
    "declared_amount": 500000,
    "status": "PENDING",
    "created_at": "2026-02-04T14:30:00Z"
  },
  "message": "Handover created successfully"
}
```

### Confirm Handover Response
```json
{
  "handover": {
    "id": "507f1f77bcf86cd799439011",
    "status": "CONFIRMED",
    "actual_amount": 500000,
    "discrepancy_amount": 0,
    "confirmed_at": "2026-02-04T14:35:00Z"
  },
  "message": "Handover confirmed successfully"
}
```

### Discrepancy Requires Approval Response
```json
{
  "handover": {
    "id": "507f1f77bcf86cd799439011",
    "status": "REQUIRES_APPROVAL",
    "actual_amount": 350000,
    "discrepancy_amount": -150000,
    "confirmed_at": "2026-02-04T14:35:00Z"
  },
  "message": "Discrepancy requires manager approval"
}
```

---

## 🧪 Testing Routes

### Manual Testing Checklist

**Waiter:**
- [ ] Can create partial handover
- [ ] Can create handover & end shift
- [ ] Can view pending handover
- [ ] Can view handover history
- [ ] Can cancel pending handover
- [ ] Cannot access cashier routes

**Cashier:**
- [ ] Can view pending handovers
- [ ] Can confirm handover (no discrepancy)
- [ ] Can confirm handover (small discrepancy)
- [ ] Can send to manager (large discrepancy)
- [ ] Can quick confirm/reject
- [ ] Cannot approve discrepancy

**Manager:**
- [ ] Can view all pending handovers
- [ ] Can view pending approvals
- [ ] Can approve discrepancy
- [ ] Can reject discrepancy
- [ ] Can view discrepancy stats
- [ ] Has full access to all features

---

## 🔗 Related Files

### Backend
- `backend/domain/handover/cash_handover.go`
- `backend/domain/handover/cash_discrepancy.go`
- `backend/application/services/cash_handover_service.go`
- `backend/interfaces/http/cash_handover_handler.go`
- `backend/infrastructure/mongodb/cash_handover_repository.go`

### Frontend
- `frontend/src/views/ShiftView.vue`
- `frontend/src/views/CashierHandoverView.vue`
- `frontend/src/views/CashierDashboard.vue`
- `frontend/src/stores/shift.js`
- `frontend/src/stores/cashier.js`
- `frontend/src/services/handover.js`
- `frontend/src/router/index.js`

### Documentation
- `docs/CASH_HANDOVER_UI_GUIDE.md`
- `docs/CASH_HANDOVER_API_DOCUMENTATION.md`
- `docs/CASH_HANDOVER_USER_GUIDE.md`

---

**Last Updated:** 2026-02-04  
**Version:** 1.0  
**Status:** ✅ Complete
