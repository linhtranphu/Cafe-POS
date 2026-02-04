# Cash Handover Implementation Plan

## 🎯 Tổng Quan Implementation

Kế hoạch triển khai tính năng bàn giao tiền giữa Waiter và Cashier với đối soát chi tiết, được chia thành 4 phases chính.

---

## 📋 Phase 1: Backend Foundation (2-3 ngày)

### 1.1 Database Schema & Models

#### 1.1.1 Tạo Domain Models
**File**: `backend/domain/handover/cash_handover.go`
```go
- HandoverStatus enum (PENDING, CONFIRMED, REJECTED, DISCREPANCY)
- HandoverType enum (PARTIAL, FULL, END_SHIFT)  
- ResponsibilityType enum (WAITER, CASHIER, SYSTEM, UNKNOWN)
- CashHandover struct với tất cả fields
- CashDiscrepancy struct cho tracking chênh lệch
- Request/Response structs
- Business logic methods (HasDiscrepancy, RequiresManagerApproval, etc.)
```

#### 1.1.2 Cập Nhật Shift Models
**File**: `backend/domain/order/shift.go`
```go
- Thêm fields: CurrentCash, HandedOverCash, RemainingCash, TotalDiscrepancy, HandoverCount
```

**File**: `backend/domain/cashier/cashier_shift.go`  
```go
- Thêm fields: ReceivedCash, TotalDiscrepancy, HandoverCount, DiscrepancyCount
```

#### 1.1.3 MongoDB Collections
```javascript
// Tạo collections mới
- cash_handovers
- cash_discrepancies

// Cập nhật existing collections
- shifts (thêm cash tracking fields)
- cashier_shifts (thêm received cash fields)
```

### 1.2 Repository Layer

#### 1.2.1 Cash Handover Repository
**File**: `backend/infrastructure/mongodb/cash_handover_repository.go`
```go
- Create(handover) error
- FindByID(id) (*CashHandover, error)
- Update(id, handover) error
- FindByWaiterShift(shiftID) ([]*CashHandover, error)
- FindByCashierShift(shiftID) ([]*CashHandover, error)
- FindPendingByCashier(cashierID) ([]*CashHandover, error)
- FindByDateRange(start, end) ([]*CashHandover, error)
- FindWithDiscrepancies() ([]*CashHandover, error)
- FindRequiringApproval() ([]*CashHandover, error)
```

#### 1.2.2 Cash Discrepancy Repository
**File**: `backend/infrastructure/mongodb/cash_discrepancy_repository.go`
```go
- Create(discrepancy) error
- FindByID(id) (*CashDiscrepancy, error)
- Update(id, discrepancy) error
- FindByHandoverID(handoverID) (*CashDiscrepancy, error)
- FindPendingResolution() ([]*CashDiscrepancy, error)
- FindRequiringApproval() ([]*CashDiscrepancy, error)
```

### 1.3 Service Layer

#### 1.3.1 Cash Handover Service
**File**: `backend/application/services/cash_handover_service.go`
```go
- CreateHandover(waiterShiftID, req, waiterID, waiterName) (*CashHandover, error)
- CreateHandoverAndEndShift(waiterShiftID, req, waiterID, waiterName) (*CashHandover, error)
- ConfirmHandoverWithReconciliation(handoverID, req, cashierID) error
- ApproveDiscrepancy(handoverID, managerID, approved, note) error
- GetDiscrepancyStats(startDate, endDate) (*DiscrepancyStats, error)
- createDiscrepancyRecord(handover) error
- updateCashAmounts(handover) error
```

**Validation Logic:**
- Validate waiter shift ownership
- Check remaining cash limits
- Validate cashier authorization
- Calculate discrepancies
- Handle manager approval thresholds

---

## 📋 Phase 2: API Layer (1-2 ngày)

### 2.1 HTTP Handlers

#### 2.1.1 Cash Handover Handler
**File**: `backend/interfaces/http/cash_handover_handler.go`
```go
- CreateHandover(c *gin.Context)
- CreateHandoverAndEndShift(c *gin.Context)  
- ConfirmHandover(c *gin.Context)
- GetPendingHandovers(c *gin.Context)
- GetTodayHandovers(c *gin.Context)
- GetHandoverHistory(c *gin.Context)
- CancelHandover(c *gin.Context)
- ReconcileHandover(c *gin.Context)
- QuickConfirm(c *gin.Context)
```

#### 2.1.2 Manager Handler Extensions
**File**: `backend/interfaces/http/manager_handler.go`
```go
- GetPendingApprovals(c *gin.Context)
- ApproveDiscrepancy(c *gin.Context)
- GetDiscrepancyStats(c *gin.Context)
```

### 2.2 API Routes

#### 2.2.1 Route Registration
**File**: `backend/main.go` hoặc routes file
```go
// Waiter routes
POST   /api/shifts/:id/handover
POST   /api/shifts/:id/handover-and-end
GET    /api/shifts/:id/pending-handover
GET    /api/shifts/:id/handovers
DELETE /api/cash-handovers/:id

// Cashier routes  
GET    /api/cash-handovers/pending
GET    /api/cash-handovers/today
POST   /api/cash-handovers/:id/reconcile
POST   /api/cash-handovers/:id/quick-confirm
GET    /api/cash-handovers/discrepancy-stats

// Manager routes
GET    /api/cash-handovers/pending-approval
POST   /api/cash-handovers/:id/approve
GET    /api/discrepancies/stats
GET    /api/discrepancies/history
```

### 2.3 Middleware & Validation

#### 2.3.1 Authorization Middleware
```go
- Waiter: Chỉ được access handover của shift mình
- Cashier: Chỉ được confirm handover assigned cho mình
- Manager: Full access cho approval và stats
```

#### 2.3.2 Request Validation
```go
- Validate amount > 0 và <= remaining_cash
- Validate handover_type enum values
- Validate discrepancy_reason khi có chênh lệch
- Validate manager_note khi approve/reject
```

---

## 📋 Phase 3: Frontend Core (3-4 ngày)

### 3.1 Store Layer (Pinia)

#### 3.1.1 Enhanced Shift Store
**File**: `frontend/src/stores/shift.js`
```javascript
// State additions
- pendingHandover: ref(null)
- handoverHistory: ref([])

// Actions
- createCashHandover(shiftId, handoverData)
- createHandoverAndEndShift(shiftId, handoverData)
- getPendingHandover(shiftId)
- getHandoverHistory(shiftId)
- cancelHandover(handoverId)
```

#### 3.1.2 Enhanced Cashier Store
**File**: `frontend/src/stores/cashier.js`
```javascript
// State additions
- pendingHandovers: ref([])
- todayHandovers: ref([])
- discrepancyThreshold: ref(50000)

// Actions
- fetchPendingHandovers()
- fetchTodayHandovers()
- reconcileHandover(handoverId, reconcileData)
- quickConfirm(handoverId, status)
- getDiscrepancyStats(startDate, endDate)
```

#### 3.1.3 New Manager Store
**File**: `frontend/src/stores/manager.js`
```javascript
// State
- pendingApprovals: ref([])
- discrepancyStats: ref({})

// Actions
- fetchPendingApprovals()
- approveDiscrepancy(handoverId, approved, note)
- getDiscrepancyStats(startDate, endDate)
```

### 3.2 Waiter Interface (ShiftView.vue)

#### 3.2.1 Template Updates
```vue
<!-- Cash Status Display -->
- Tiền hiện có, Đã bàn giao, Tổng thu

<!-- Pending Handover Status -->
- Banner hiển thị handover đang chờ
- Nút hủy yêu cầu

<!-- Action Buttons -->
- "💰 Bàn giao một phần" button
- "🏁 Bàn giao và đóng ca" button  
- "Kết thúc ca" button (khi remaining_cash = 0)

<!-- Handover History Section -->
- Danh sách lịch sử bàn giao
- Status và phản hồi từ cashier
```

#### 3.2.2 Modal Components
```vue
<!-- Partial Handover Modal -->
- Form nhập số tiền và ghi chú
- Validation amount <= remaining_cash

<!-- Handover and End Shift Modal -->
- Cảnh báo về thao tác không thể hoàn tác
- Form nhập tiền cuối ca và ghi chú
- Hiển thị số tiền sẽ bàn giao
```

#### 3.2.3 Script Logic
```javascript
// Reactive data
- showPartialHandoverForm: ref(false)
- showHandoverEndShiftForm: ref(false)
- partialHandoverForm: ref({})
- handoverEndShiftForm: ref({})

// Methods
- createPartialHandover()
- createHandoverAndEndShift()
- cancelHandover(handoverId)
- fetchHandoverData()

// Helper functions
- getHandoverTypeText(type)
- getHandoverStatusText(status)
- getHandoverStatusClass(status)
```

### 3.3 Cashier Interface

#### 3.3.1 Enhanced CashierDashboard.vue
```vue
<!-- Handover Notifications -->
- Alert banner cho pending handovers
- Quick action buttons (✅/❌)

<!-- Quick Handover Section -->
- Top 3 pending handovers
- Quick confirm buttons
- Link to full handover management
```

#### 3.3.2 New CashierHandoverView.vue
```vue
<!-- Pending Handovers Section -->
- Danh sách yêu cầu chờ xác nhận
- Thông tin chi tiết handover
- Action buttons (Xác nhận/Từ chối)

<!-- Today's Handovers Section -->
- Lịch sử bàn giao hôm nay
- Status và ghi chú

<!-- Reconciliation Modal -->
- Form nhập actual amount
- Discrepancy calculation và display
- Discrepancy reason selection
- Responsibility assignment
- Large discrepancy warning
```

---

## 📋 Phase 4: Advanced Features (2-3 ngày)

### 4.1 Manager Interface

#### 4.1.1 New DiscrepancyApprovalView.vue
```vue
<!-- Pending Approvals Section -->
- Danh sách chênh lệch cần phê duyệt
- Chi tiết handover và discrepancy
- Approval/rejection form

<!-- Discrepancy Statistics -->
- Thống kê chênh lệch (shortage/overage)
- Charts và metrics
- Performance indicators

<!-- Approval Modal -->
- Form phê duyệt/từ chối
- Manager note input
- Confirmation workflow
```

### 4.2 Navigation & Routing

#### 4.2.1 Router Updates
**File**: `frontend/src/router/index.js`
```javascript
// New routes
{
  path: '/cashier/handovers',
  name: 'CashierHandovers',
  component: CashierHandoverView,
  meta: { requiresAuth: true, roles: ['cashier', 'manager'] }
},
{
  path: '/manager/discrepancies',
  name: 'DiscrepancyApproval', 
  component: DiscrepancyApprovalView,
  meta: { requiresAuth: true, roles: ['manager'] }
}
```

#### 4.2.2 Navigation Menu Updates
**File**: `frontend/src/components/Navigation.vue`
```vue
<!-- Cashier menu item -->
<router-link to="/cashier/handovers">
  💰 Bàn giao tiền
  <badge v-if="pendingCount > 0">{{ pendingCount }}</badge>
</router-link>

<!-- Manager menu item -->
<router-link to="/manager/discrepancies">
  ⚖️ Phê duyệt chênh lệch
  <badge v-if="approvalCount > 0">{{ approvalCount }}</badge>
</router-link>
```

### 4.3 Real-time Features

#### 4.3.1 WebSocket Integration (Optional)
```javascript
// Real-time notifications
- Waiter → Cashier: New handover request
- Cashier → Waiter: Handover confirmed/rejected
- System → Manager: Large discrepancy needs approval

// Auto-refresh triggers
- Handover status changes
- New pending requests
- Approval notifications
```

### 4.4 Reporting & Analytics

#### 4.4.1 Discrepancy Reports
```vue
<!-- Report Components -->
- Daily discrepancy summary
- User performance metrics
- Trend analysis charts
- Export functionality
```

---

## 📋 Phase 5: Testing & Polish (2-3 ngày)

### 5.1 Backend Testing

#### 5.1.1 Unit Tests
**Files**: `backend/application/services/*_test.go`
```go
// CashHandoverService tests
- TestCreateHandover()
- TestConfirmHandoverWithReconciliation()
- TestApproveDiscrepancy()
- TestDiscrepancyCalculation()
- TestManagerApprovalThreshold()

// Repository tests
- TestCashHandoverRepository()
- TestCashDiscrepancyRepository()
```

#### 5.1.2 Integration Tests
```go
// API endpoint tests
- TestHandoverWorkflow()
- TestDiscrepancyHandling()
- TestManagerApprovalFlow()
- TestAuthorizationRules()
```

### 5.2 Frontend Testing

#### 5.2.1 Component Tests
```javascript
// Vue component tests
- ShiftView handover functionality
- CashierHandoverView reconciliation
- DiscrepancyApprovalView approval flow
- Store action tests
```

#### 5.2.2 E2E Testing
```javascript
// Cypress tests
- Complete handover workflow
- Discrepancy handling flow
- Manager approval process
- Error scenarios
```

### 5.3 Data Migration

#### 5.3.1 Database Migration Scripts
```javascript
// Migration scripts
- Add new collections (cash_handovers, cash_discrepancies)
- Update existing shift collections with new fields
- Create indexes for performance
- Data validation scripts
```

### 5.4 Configuration & Deployment

#### 5.4.1 Environment Configuration
```go
// Config additions
- DISCREPANCY_THRESHOLD (default: 50000)
- MANAGER_APPROVAL_REQUIRED (default: true)
- HANDOVER_TIMEOUT (default: 24h)
```

#### 5.4.2 Documentation Updates
```markdown
- API documentation
- User guides for each role
- Admin configuration guide
- Troubleshooting guide
```

---

## 🚀 Deployment Checklist

### Pre-deployment
- [ ] All unit tests passing
- [ ] Integration tests passing
- [ ] E2E tests passing
- [ ] Code review completed
- [ ] Database migration scripts ready
- [ ] Configuration updated

### Deployment Steps
1. [ ] Run database migrations
2. [ ] Deploy backend changes
3. [ ] Deploy frontend changes
4. [ ] Update environment variables
5. [ ] Verify API endpoints
6. [ ] Test critical workflows
7. [ ] Monitor error logs

### Post-deployment
- [ ] User acceptance testing
- [ ] Performance monitoring
- [ ] Error rate monitoring
- [ ] User feedback collection
- [ ] Documentation updates

---

## 📊 Estimated Timeline

| Phase | Duration | Dependencies |
|-------|----------|--------------|
| Phase 1: Backend Foundation | 2-3 ngày | Database design approval |
| Phase 2: API Layer | 1-2 ngày | Phase 1 complete |
| Phase 3: Frontend Core | 3-4 ngày | Phase 2 complete |
| Phase 4: Advanced Features | 2-3 ngày | Phase 3 complete |
| Phase 5: Testing & Polish | 2-3 ngày | All phases complete |

**Total Estimated Time: 10-15 ngày**

---

## 🎯 Success Criteria

### Functional Requirements
- [ ] Waiter có thể tạo handover request (partial/full)
- [ ] Cashier có thể reconcile và confirm/reject
- [ ] Discrepancy được detect và handle đúng
- [ ] Manager approval workflow hoạt động
- [ ] Audit trail đầy đủ và chính xác

### Performance Requirements  
- [ ] API response time < 500ms
- [ ] Real-time updates < 2s delay
- [ ] Database queries optimized
- [ ] Frontend loading < 3s

### Security Requirements
- [ ] Role-based access control
- [ ] Data validation đầy đủ
- [ ] Audit logging secure
- [ ] Sensitive data encrypted

### User Experience
- [ ] Intuitive interface cho tất cả roles
- [ ] Clear error messages
- [ ] Responsive design
- [ ] Accessibility compliance

---

## 🔧 Technical Considerations

### Database Performance
- Index trên các query fields thường dùng
- Pagination cho large datasets
- Archive old handover records

### Security
- Input validation và sanitization
- Rate limiting cho API endpoints
- Audit log immutability
- Data encryption at rest

### Scalability
- Horizontal scaling support
- Caching strategy cho frequent queries
- Background job processing cho heavy operations
- Database connection pooling

### Monitoring
- Application metrics
- Error tracking
- Performance monitoring
- Business metrics dashboard