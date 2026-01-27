# 💵 Cashier System - Implementation Plan

## 📋 Phase Overview

### Phase 1: Domain Layer (Backend)
- Cash reconciliation entities
- Payment discrepancy tracking
- Shift closure with audit

### Phase 2: Repository Layer (Backend)
- Cash reconciliation CRUD
- Payment audit queries
- Shift closure operations

### Phase 3: Service Layer (Backend)
- Reconciliation business logic
- Payment oversight services
- Audit trail management

### Phase 4: Handler Layer (Backend)
- Cashier HTTP endpoints
- Payment control APIs
- Reconciliation endpoints

### Phase 5: Routes Integration (Backend)
- Cashier-specific routes
- Authorization middleware
- API documentation

### Phase 6: Frontend Services & Stores
- Cashier API services
- State management
- Real-time updates

### Phase 7: Frontend Views
- Shift management UI
- Payment oversight dashboard
- Reconciliation interface

---

## ✅ Phase 1: Domain Layer

### Files to Create:

#### 1. `backend/domain/cashier/cash_reconciliation.go`
```go
type CashReconciliation struct {
    ID                string    `json:"id" bson:"_id,omitempty"`
    ShiftID           string    `json:"shift_id" bson:"shift_id"`
    CashierID         string    `json:"cashier_id" bson:"cashier_id"`
    ExpectedCash      float64   `json:"expected_cash" bson:"expected_cash"`
    ActualCash        float64   `json:"actual_cash" bson:"actual_cash"`
    Difference        float64   `json:"difference" bson:"difference"`
    Status            string    `json:"status" bson:"status"` // MATCH, OVER, SHORT
    Notes             string    `json:"notes" bson:"notes"`
    ReconciliationAt  time.Time `json:"reconciliation_at" bson:"reconciliation_at"`
}

type PaymentDiscrepancy struct {
    ID          string    `json:"id" bson:"_id,omitempty"`
    OrderID     string    `json:"order_id" bson:"order_id"`
    CashierID   string    `json:"cashier_id" bson:"cashier_id"`
    Reason      string    `json:"reason" bson:"reason"`
    Amount      float64   `json:"amount" bson:"amount"`
    Status      string    `json:"status" bson:"status"` // PENDING, RESOLVED
    ReportedAt  time.Time `json:"reported_at" bson:"reported_at"`
    ResolvedAt  *time.Time `json:"resolved_at,omitempty" bson:"resolved_at,omitempty"`
}
```

#### 2. `backend/domain/cashier/shift_closure.go`
```go
type ShiftClosure struct {
    ID              string    `json:"id" bson:"_id,omitempty"`
    ShiftID         string    `json:"shift_id" bson:"shift_id"`
    CashierID       string    `json:"cashier_id" bson:"cashier_id"`
    TotalOrders     int       `json:"total_orders" bson:"total_orders"`
    TotalRevenue    float64   `json:"total_revenue" bson:"total_revenue"`
    CashRevenue     float64   `json:"cash_revenue" bson:"cash_revenue"`
    TransferRevenue float64   `json:"transfer_revenue" bson:"transfer_revenue"`
    QRRevenue       float64   `json:"qr_revenue" bson:"qr_revenue"`
    Discrepancies   []string  `json:"discrepancies" bson:"discrepancies"`
    ClosedAt        time.Time `json:"closed_at" bson:"closed_at"`
}
```

#### 3. `backend/domain/cashier/payment_audit.go`
```go
type PaymentAudit struct {
    ID           string    `json:"id" bson:"_id,omitempty"`
    OrderID      string    `json:"order_id" bson:"order_id"`
    Action       string    `json:"action" bson:"action"` // CANCEL, REFUND, OVERRIDE
    CashierID    string    `json:"cashier_id" bson:"cashier_id"`
    Reason       string    `json:"reason" bson:"reason"`
    OldStatus    string    `json:"old_status" bson:"old_status"`
    NewStatus    string    `json:"new_status" bson:"new_status"`
    Amount       float64   `json:"amount" bson:"amount"`
    AuditedAt    time.Time `json:"audited_at" bson:"audited_at"`
}
```

---

## ✅ Phase 2: Repository Layer

### Files to Create:

#### 1. `backend/infrastructure/mongodb/cash_reconciliation_repository.go`
```go
type CashReconciliationRepository interface {
    Create(reconciliation *domain.CashReconciliation) error
    FindByShiftID(shiftID string) (*domain.CashReconciliation, error)
    FindByCashierID(cashierID string) ([]*domain.CashReconciliation, error)
    Update(reconciliation *domain.CashReconciliation) error
}
```

#### 2. `backend/infrastructure/mongodb/payment_discrepancy_repository.go`
```go
type PaymentDiscrepancyRepository interface {
    Create(discrepancy *domain.PaymentDiscrepancy) error
    FindByOrderID(orderID string) ([]*domain.PaymentDiscrepancy, error)
    FindPendingDiscrepancies() ([]*domain.PaymentDiscrepancy, error)
    UpdateStatus(id string, status string) error
    FindByShiftID(shiftID string) ([]*domain.PaymentDiscrepancy, error)
}
```

#### 3. `backend/infrastructure/mongodb/payment_audit_repository.go`
```go
type PaymentAuditRepository interface {
    Create(audit *domain.PaymentAudit) error
    FindByOrderID(orderID string) ([]*domain.PaymentAudit, error)
    FindByCashierID(cashierID string) ([]*domain.PaymentAudit, error)
    FindByDateRange(start, end time.Time) ([]*domain.PaymentAudit, error)
}
```

---

## ✅ Phase 3: Service Layer

### Files to Create:

#### 1. `backend/application/services/cash_reconciliation_service.go`
```go
type CashReconciliationService struct {
    reconciliationRepo domain.CashReconciliationRepository
    shiftRepo         domain.ShiftRepository
    orderRepo         domain.OrderRepository
}

// FR-CASH-06: Đối soát tiền mặt
func (s *CashReconciliationService) ReconcileCash(shiftID string, actualCash float64, notes string) error

// FR-CASH-03: Chốt ca
func (s *CashReconciliationService) CloseShift(shiftID string, cashierID string) (*domain.ShiftClosure, error)

// FR-CASH-02: Theo dõi trạng thái ca
func (s *CashReconciliationService) GetShiftStatus(shiftID string) (*ShiftStatusResponse, error)
```

#### 2. `backend/application/services/payment_oversight_service.go`
```go
type PaymentOversightService struct {
    orderRepo       domain.OrderRepository
    discrepancyRepo domain.PaymentDiscrepancyRepository
    auditRepo       domain.PaymentAuditRepository
}

// FR-CASH-04: Giám sát thanh toán
func (s *PaymentOversightService) GetPaymentsByShift(shiftID string) ([]*PaymentSummary, error)

// FR-CASH-05: Xử lý sai lệch thanh toán
func (s *PaymentOversightService) ReportDiscrepancy(orderID, reason string, amount float64) error

// FR-CASH-08: Hủy/điều chỉnh thanh toán
func (s *PaymentOversightService) OverridePayment(orderID, reason string, cashierID string) error

// FR-CASH-09: Khóa order
func (s *PaymentOversightService) LockOrder(orderID string, cashierID string) error
```

#### 3. `backend/application/services/cashier_report_service.go`
```go
type CashierReportService struct {
    orderRepo         domain.OrderRepository
    reconciliationRepo domain.CashReconciliationRepository
    shiftRepo         domain.ShiftRepository
}

// FR-CASH-10: Báo cáo ca
func (s *CashierReportService) GenerateShiftReport(shiftID string) (*ShiftReport, error)

// FR-CASH-11: Bàn giao ca
func (s *CashierReportService) HandoverShift(fromCashierID, toCashierID string) error
```

---

## ✅ Phase 4: Handler Layer

### Files to Create:

#### 1. `backend/interfaces/http/cashier_handler.go`
```go
type CashierHandler struct {
    reconciliationService *services.CashReconciliationService
    oversightService     *services.PaymentOversightService
    reportService        *services.CashierReportService
}

// FR-CASH-01: Mở ca
func (h *CashierHandler) OpenShift(c *gin.Context)

// FR-CASH-02: Theo dõi trạng thái ca
func (h *CashierHandler) GetShiftStatus(c *gin.Context)

// FR-CASH-03: Chốt ca
func (h *CashierHandler) CloseShift(c *gin.Context)

// FR-CASH-04: Giám sát thanh toán
func (h *CashierHandler) GetPaymentsByShift(c *gin.Context)

// FR-CASH-05: Xử lý sai lệch thanh toán
func (h *CashierHandler) ReportDiscrepancy(c *gin.Context)

// FR-CASH-06: Đối soát tiền mặt
func (h *CashierHandler) ReconcileCash(c *gin.Context)

// FR-CASH-07: Đối soát chuyển khoản
func (h *CashierHandler) ConfirmTransfers(c *gin.Context)

// FR-CASH-08: Hủy/điều chỉnh thanh toán
func (h *CashierHandler) OverridePayment(c *gin.Context)

// FR-CASH-09: Khóa order
func (h *CashierHandler) LockOrder(c *gin.Context)

// FR-CASH-10: Báo cáo ca
func (h *CashierHandler) GenerateShiftReport(c *gin.Context)

// FR-CASH-11: Bàn giao ca
func (h *CashierHandler) HandoverShift(c *gin.Context)
```

---

## ✅ Phase 5: Routes Integration

### File to Update:

#### 1. `backend/main.go` - Add Cashier Routes
```go
// Cashier Routes (/api/cashier/*)
cashierGroup := api.Group("/cashier")
cashierGroup.Use(middleware.RequireRole("cashier", "manager"))
{
    // Shift Management
    cashierGroup.POST("/shifts/open", cashierHandler.OpenShift)
    cashierGroup.GET("/shifts/:id/status", cashierHandler.GetShiftStatus)
    cashierGroup.POST("/shifts/:id/close", cashierHandler.CloseShift)
    
    // Payment Oversight
    cashierGroup.GET("/shifts/:id/payments", cashierHandler.GetPaymentsByShift)
    cashierGroup.POST("/discrepancies", cashierHandler.ReportDiscrepancy)
    cashierGroup.POST("/orders/:id/override", cashierHandler.OverridePayment)
    cashierGroup.POST("/orders/:id/lock", cashierHandler.LockOrder)
    
    // Reconciliation
    cashierGroup.POST("/reconcile/cash", cashierHandler.ReconcileCash)
    cashierGroup.POST("/reconcile/transfers", cashierHandler.ConfirmTransfers)
    
    // Reports
    cashierGroup.GET("/reports/shift/:id", cashierHandler.GenerateShiftReport)
    cashierGroup.POST("/handover", cashierHandler.HandoverShift)
}
```

---

## ✅ Phase 6: Frontend Services & Stores

### Files to Create:

#### 1. `frontend/src/services/cashier.js`
```javascript
export const cashierService = {
  // Shift Management
  openShift: (data) => api.post('/cashier/shifts/open', data),
  getShiftStatus: (shiftId) => api.get(`/cashier/shifts/${shiftId}/status`),
  closeShift: (shiftId, data) => api.post(`/cashier/shifts/${shiftId}/close`, data),
  
  // Payment Oversight
  getPaymentsByShift: (shiftId) => api.get(`/cashier/shifts/${shiftId}/payments`),
  reportDiscrepancy: (data) => api.post('/cashier/discrepancies', data),
  overridePayment: (orderId, data) => api.post(`/cashier/orders/${orderId}/override`, data),
  lockOrder: (orderId) => api.post(`/cashier/orders/${orderId}/lock`),
  
  // Reconciliation
  reconcileCash: (data) => api.post('/cashier/reconcile/cash', data),
  confirmTransfers: (data) => api.post('/cashier/reconcile/transfers', data),
  
  // Reports
  generateShiftReport: (shiftId) => api.get(`/cashier/reports/shift/${shiftId}`),
  handoverShift: (data) => api.post('/cashier/handover', data)
}
```

#### 2. `frontend/src/stores/cashier.js`
```javascript
export const useCashierStore = defineStore('cashier', {
  state: () => ({
    currentShift: null,
    shiftStatus: null,
    payments: [],
    discrepancies: [],
    reconciliation: null,
    reports: []
  }),
  
  actions: {
    async openShift(data) { /* FR-CASH-01 */ },
    async getShiftStatus(shiftId) { /* FR-CASH-02 */ },
    async closeShift(shiftId, data) { /* FR-CASH-03 */ },
    async getPaymentsByShift(shiftId) { /* FR-CASH-04 */ },
    async reportDiscrepancy(data) { /* FR-CASH-05 */ },
    async reconcileCash(data) { /* FR-CASH-06 */ },
    async confirmTransfers(data) { /* FR-CASH-07 */ },
    async overridePayment(orderId, data) { /* FR-CASH-08 */ },
    async lockOrder(orderId) { /* FR-CASH-09 */ },
    async generateShiftReport(shiftId) { /* FR-CASH-10 */ },
    async handoverShift(data) { /* FR-CASH-11 */ }
  }
})
```

---

## ✅ Phase 7: Frontend Views

### Files to Create:

#### 1. `frontend/src/views/CashierDashboard.vue`
```vue
<template>
  <div class="cashier-dashboard">
    <!-- FR-CASH-02: Theo dõi trạng thái ca -->
    <ShiftStatusCard :shift="currentShift" />
    
    <!-- FR-CASH-04: Giám sát thanh toán -->
    <PaymentOversightPanel :payments="payments" />
    
    <!-- FR-CASH-05: Xử lý sai lệch thanh toán -->
    <DiscrepancyPanel :discrepancies="discrepancies" />
  </div>
</template>
```

#### 2. `frontend/src/views/ShiftReconciliation.vue`
```vue
<template>
  <div class="shift-reconciliation">
    <!-- FR-CASH-06: Đối soát tiền mặt -->
    <CashReconciliationForm @submit="reconcileCash" />
    
    <!-- FR-CASH-07: Đối soát chuyển khoản -->
    <TransferConfirmationPanel @confirm="confirmTransfers" />
    
    <!-- FR-CASH-03: Chốt ca -->
    <ShiftClosurePanel @close="closeShift" />
  </div>
</template>
```

#### 3. `frontend/src/views/CashierReports.vue`
```vue
<template>
  <div class="cashier-reports">
    <!-- FR-CASH-10: Báo cáo ca -->
    <ShiftReportGenerator @generate="generateReport" />
    
    <!-- FR-CASH-11: Bàn giao ca -->
    <ShiftHandoverForm @handover="handoverShift" />
  </div>
</template>
```

#### 4. `frontend/src/router/index.js` - Add Routes
```javascript
{
  path: '/cashier',
  component: () => import('@/layouts/CashierLayout.vue'),
  meta: { requiresAuth: true, roles: ['cashier', 'manager'] },
  children: [
    { path: 'dashboard', component: () => import('@/views/CashierDashboard.vue') },
    { path: 'reconciliation', component: () => import('@/views/ShiftReconciliation.vue') },
    { path: 'reports', component: () => import('@/views/CashierReports.vue') }
  ]
}
```

---

## 📊 Implementation Summary

### Backend Files (13 files):
- **Domain Layer**: 3 files (cash_reconciliation.go, shift_closure.go, payment_audit.go)
- **Repository Layer**: 3 files (reconciliation, discrepancy, audit repositories)
- **Service Layer**: 3 files (reconciliation, oversight, report services)
- **Handler Layer**: 1 file (cashier_handler.go)
- **Routes**: 1 file updated (main.go)

### Frontend Files (8 files):
- **Services**: 1 file (cashier.js)
- **Stores**: 1 file (cashier.js)
- **Views**: 3 files (Dashboard, Reconciliation, Reports)
- **Router**: 1 file updated (index.js)
- **Components**: 2 files (ShiftStatusCard, PaymentOversightPanel)

### API Endpoints (11 endpoints):
- Shift Management: 3 endpoints
- Payment Oversight: 4 endpoints
- Reconciliation: 2 endpoints
- Reports: 2 endpoints

### Key Features:
- ✅ Real-time shift status monitoring
- ✅ Payment discrepancy tracking
- ✅ Cash reconciliation with variance detection
- ✅ Transfer confirmation workflow
- ✅ Payment override with audit trail
- ✅ Order locking mechanism
- ✅ Comprehensive shift reports
- ✅ Shift handover process

### Business Rules Enforced:
- ✅ Cashier can only close shifts after reconciliation
- ✅ All payment overrides require reason and audit log
- ✅ Discrepancies must be resolved before shift closure
- ✅ Locked orders cannot be modified
- ✅ Transfer confirmations required for shift completion
- ✅ Handover process tracks cashier transitions

**Total Implementation: 21 files (13 backend + 8 frontend)**