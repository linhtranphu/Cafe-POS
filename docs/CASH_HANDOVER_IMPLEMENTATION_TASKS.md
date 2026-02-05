# 🚀 Cash Handover Implementation Tasks

## 📋 Tổng Quan

Triển khai tính năng bàn giao tiền từ Waiter → Cashier với đối soát chi tiết và xử lý chênh lệch.

**Quyết định quan trọng:** XÓA handover cũ (cashier-to-cashier) vì hoàn toàn khác với yêu cầu mới.

---

## 🗑️ PHASE 0: Cleanup Old Handover (CRITICAL - DO FIRST)

### Task 0.1: Remove Old Backend Handover
**Files to modify:**
- `backend/application/services/cashier_report_service.go`
- `backend/interfaces/http/cashier_handler.go`
- `backend/main.go`

**Actions:**
```go
// 1. Remove from cashier_report_service.go
// - Delete HandoverData struct (lines ~104-108)
// - Delete HandoverShift() method (lines ~111-160)

// 2. Remove from cashier_handler.go
// - Delete HandoverShift() handler (lines ~156-172)

// 3. Remove from main.go
// - Delete route: cashier.POST("/handover", cashierHandler.HandoverShift)
```

**Verification:**
- [ ] No compilation errors
- [ ] No references to old HandoverShift in codebase
- [ ] Run: `cd backend && go build`

---

### Task 0.2: Remove Old Frontend Handover
**Files to modify:**
- `frontend/src/views/CashierReports.vue`
- `frontend/src/stores/cashier.js`
- `frontend/src/services/cashier.js`

**Actions:**
```javascript
// 1. Remove from CashierReports.vue
// - Delete "Shift Handover" section (lines ~71-94)
// - Delete handoverForm ref (lines ~265-268)
// - Delete performHandover() method (lines ~301-316)

// 2. Remove from cashier.js store
// - Delete handoverShift() action (lines ~193-204)

// 3. Remove from cashier.js service
// - Delete handoverShift line from exports
```

**Verification:**
- [ ] No console errors
- [ ] CashierReports page loads correctly
- [ ] Run: `cd frontend && npm run build`

---

## 🏗️ PHASE 1: Backend Domain Layer

### Task 1.1: Create Handover Domain Package
**Create:** `backend/domain/handover/cash_handover.go`

**Content:**
```go
package handover

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Enums
type HandoverStatus string
type HandoverType string
type DiscrepancyType string
type ResponsibilityType string

const (
    StatusPending     HandoverStatus = "PENDING"
    StatusConfirmed   HandoverStatus = "CONFIRMED"
    StatusRejected    HandoverStatus = "REJECTED"
    StatusDiscrepancy HandoverStatus = "DISCREPANCY"
    
    TypePartial   HandoverType = "PARTIAL"
    TypeFull      HandoverType = "FULL"
    TypeEndShift  HandoverType = "END_SHIFT"
    
    DiscrepancyShortage DiscrepancyType = "SHORTAGE"
    DiscrepancyOverage  DiscrepancyType = "OVERAGE"
    
    ResponsibilityWaiter   ResponsibilityType = "WAITER"
    ResponsibilityCashier  ResponsibilityType = "CASHIER"
    ResponsibilitySystem   ResponsibilityType = "SYSTEM"
    ResponsibilityCustomer ResponsibilityType = "CUSTOMER"
    ResponsibilityUnknown  ResponsibilityType = "UNKNOWN"
)

// Main struct - copy from analysis doc
type CashHandover struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    WaiterShiftID   primitive.ObjectID `bson:"waiter_shift_id" json:"waiter_shift_id"`
    CashierShiftID  primitive.ObjectID `bson:"cashier_shift_id" json:"cashier_shift_id"`
    // ... (copy full struct from doc)
}

// Helper methods
func (h *CashHandover) HasDiscrepancy() bool {
    return h.Discrepancy != 0
}

// ... (copy all methods from doc)
```

**Checklist:**
- [ ] Create package directory
- [ ] Copy all structs from analysis doc
- [ ] Copy all const definitions
- [ ] Copy all helper methods
- [ ] Add proper imports
- [ ] Test compilation: `go build`

---

### Task 1.2: Create Discrepancy Domain
**Create:** `backend/domain/handover/cash_discrepancy.go`

**Content:**
```go
package handover

// Copy CashDiscrepancy struct from analysis doc
type CashDiscrepancy struct {
    ID                      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    HandoverID              primitive.ObjectID `bson:"handover_id" json:"handover_id"`
    // ... (copy full struct)
}
```

**Checklist:**
- [ ] Create file
- [ ] Copy CashDiscrepancy struct
- [ ] Add proper imports
- [ ] Test compilation

---

### Task 1.3: Update Shift Domain Models
**Modify:** `backend/domain/order/shift.go`

**Add to Shift struct:**
```go
type Shift struct {
    // ... existing fields
    
    // Cash handover fields (add at end)
    CurrentCash         float64 `bson:"current_cash" json:"current_cash"`
    HandedOverCash      float64 `bson:"handed_over_cash" json:"handed_over_cash"`
    RemainingCash       float64 `bson:"remaining_cash" json:"remaining_cash"`
    TotalDiscrepancy    float64 `bson:"total_discrepancy" json:"total_discrepancy"`
    HandoverCount       int     `bson:"handover_count" json:"handover_count"`
}
```

**Modify:** `backend/domain/cashier/cashier_shift.go`

**Add to CashierShift struct:**
```go
type CashierShift struct {
    // ... existing fields
    
    // Cash handover fields (add at end)
    ReceivedCash        float64 `bson:"received_cash" json:"received_cash"`
    TotalDiscrepancy    float64 `bson:"total_discrepancy" json:"total_discrepancy"`
    HandoverCount       int     `bson:"handover_count" json:"handover_count"`
    DiscrepancyCount    int     `bson:"discrepancy_count" json:"discrepancy_count"`
}
```

**Checklist:**
- [ ] Add fields to Shift struct
- [ ] Add fields to CashierShift struct
- [ ] Test compilation
- [ ] Verify no breaking changes

---

## 🗄️ PHASE 2: Backend Repository Layer

### Task 2.1: Create Handover Repository
**Create:** `backend/infrastructure/mongodb/cash_handover_repository.go`

**Implementation:**
```go
package mongodb

import (
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "your-project/backend/domain/handover"
)

type CashHandoverRepository struct {
    collection *mongo.Collection
}

func NewCashHandoverRepository(db *mongo.Database) *CashHandoverRepository {
    return &CashHandoverRepository{
        collection: db.Collection("cash_handovers"),
    }
}

// Implement all methods from analysis doc:
// - Create()
// - FindByID()
// - Update()
// - FindByWaiterShift()
// - FindByCashierShift()
// - FindPendingByCashier()
// - FindByDateRange()
// - FindWithDiscrepancies()
// - FindRequiringApproval()
```

**Checklist:**
- [ ] Create file
- [ ] Implement all CRUD methods
- [ ] Add proper error handling
- [ ] Add indexes for performance
- [ ] Test compilation

---

### Task 2.2: Create Discrepancy Repository
**Create:** `backend/infrastructure/mongodb/cash_discrepancy_repository.go`

**Implementation:**
```go
package mongodb

type CashDiscrepancyRepository struct {
    collection *mongo.Collection
}

// Implement all methods from analysis doc:
// - Create()
// - FindByID()
// - Update()
// - FindByHandoverID()
// - FindPendingResolution()
// - FindRequiringApproval()
// - FindByDateRange()
```

**Checklist:**
- [ ] Create file
- [ ] Implement all methods
- [ ] Add proper error handling
- [ ] Test compilation

---

## 🔧 PHASE 3: Backend Service Layer

### Task 3.1: Create Handover Service
**Create:** `backend/application/services/cash_handover_service.go`

**Key methods to implement:**
1. `CreateHandover()` - Waiter creates handover request
2. `ConfirmHandoverWithReconciliation()` - Cashier confirms with actual amount
3. `createDiscrepancyRecord()` - Internal helper
4. `updateCashAmounts()` - Update shift cash amounts
5. `ApproveDiscrepancy()` - Manager approval for large discrepancies
6. `GetDiscrepancyStats()` - Statistics

**Implementation structure:**
```go
package services

type CashHandoverService struct {
    handoverRepo        *mongodb.CashHandoverRepository
    discrepancyRepo     *mongodb.CashDiscrepancyRepository
    shiftRepo           *mongodb.ShiftRepository
    cashierShiftRepo    *mongodb.CashierShiftRepository
    orderRepo           *mongodb.OrderRepository
    discrepancyThreshold float64
}

func NewCashHandoverService(
    handoverRepo *mongodb.CashHandoverRepository,
    // ... other repos
) *CashHandoverService {
    return &CashHandoverService{
        handoverRepo: handoverRepo,
        // ...
        discrepancyThreshold: 100000, // 100k VND
    }
}

// Copy all methods from analysis doc
```

**Checklist:**
- [ ] Create service struct
- [ ] Implement CreateHandover()
- [ ] Implement ConfirmHandoverWithReconciliation()
- [ ] Implement createDiscrepancyRecord()
- [ ] Implement updateCashAmounts()
- [ ] Implement ApproveDiscrepancy()
- [ ] Implement GetDiscrepancyStats()
- [ ] Add proper validation
- [ ] Add error handling
- [ ] Test compilation

---

### Task 3.2: Add Handover Methods to Shift Service
**Modify:** `backend/application/services/shift_service.go`

**Add methods:**
```go
// GetPendingHandover returns pending handover for a shift
func (s *ShiftService) GetPendingHandover(ctx context.Context, shiftID primitive.ObjectID) (*handover.CashHandover, error)

// GetHandoverHistory returns all handovers for a shift
func (s *ShiftService) GetHandoverHistory(ctx context.Context, shiftID primitive.ObjectID) ([]*handover.CashHandover, error)

// CancelHandover cancels a pending handover
func (s *ShiftService) CancelHandover(ctx context.Context, handoverID primitive.ObjectID) error
```

**Checklist:**
- [ ] Add handoverRepo to ShiftService
- [ ] Implement GetPendingHandover()
- [ ] Implement GetHandoverHistory()
- [ ] Implement CancelHandover()
- [ ] Test compilation

---

## 🌐 PHASE 4: Backend HTTP Layer

### Task 4.1: Create Handover Handler
**Create:** `backend/interfaces/http/cash_handover_handler.go`

**Endpoints to implement:**
```go
type CashHandoverHandler struct {
    handoverService *services.CashHandoverService
}

// Waiter endpoints
func (h *CashHandoverHandler) CreateHandover(c *gin.Context)
func (h *CashHandoverHandler) CreateHandoverAndEndShift(c *gin.Context)
func (h *CashHandoverHandler) GetPendingHandover(c *gin.Context)
func (h *CashHandoverHandler) GetHandoverHistory(c *gin.Context)
func (h *CashHandoverHandler) CancelHandover(c *gin.Context)

// Cashier endpoints
func (h *CashHandoverHandler) GetPendingHandovers(c *gin.Context)
func (h *CashHandoverHandler) GetTodayHandovers(c *gin.Context)
func (h *CashHandoverHandler) ConfirmHandover(c *gin.Context)
func (h *CashHandoverHandler) QuickConfirm(c *gin.Context)

// Manager endpoints
func (h *CashHandoverHandler) GetPendingApprovals(c *gin.Context)
func (h *CashHandoverHandler) ApproveDiscrepancy(c *gin.Context)
func (h *CashHandoverHandler) GetDiscrepancyStats(c *gin.Context)
```

**Checklist:**
- [ ] Create handler struct
- [ ] Implement all waiter endpoints
- [ ] Implement all cashier endpoints
- [ ] Implement all manager endpoints
- [ ] Add proper validation
- [ ] Add error handling
- [ ] Test compilation

---

### Task 4.2: Register Routes
**Modify:** `backend/main.go`

**Add routes:**
```go
// Initialize handler
handoverHandler := http.NewCashHandoverHandler(handoverService)

// Waiter routes
waiter := api.Group("/api/shifts")
waiter.Use(authMiddleware, roleMiddleware("waiter"))
{
    waiter.POST("/:id/handover", handoverHandler.CreateHandover)
    waiter.POST("/:id/handover-and-end", handoverHandler.CreateHandoverAndEndShift)
    waiter.GET("/:id/pending-handover", handoverHandler.GetPendingHandover)
    waiter.GET("/:id/handovers", handoverHandler.GetHandoverHistory)
}

// Cashier routes
cashier := api.Group("/api/cash-handovers")
cashier.Use(authMiddleware, roleMiddleware("cashier", "manager"))
{
    cashier.GET("/pending", handoverHandler.GetPendingHandovers)
    cashier.GET("/today", handoverHandler.GetTodayHandovers)
    cashier.POST("/:id/confirm", handoverHandler.ConfirmHandover)
    cashier.POST("/:id/quick-confirm", handoverHandler.QuickConfirm)
    cashier.DELETE("/:id", handoverHandler.CancelHandover)
}

// Manager routes
manager := api.Group("/api/cash-handovers")
manager.Use(authMiddleware, roleMiddleware("manager"))
{
    manager.GET("/pending-approval", handoverHandler.GetPendingApprovals)
    manager.POST("/:id/approve", handoverHandler.ApproveDiscrepancy)
    manager.GET("/discrepancy-stats", handoverHandler.GetDiscrepancyStats)
}
```

**Checklist:**
- [ ] Initialize handler with dependencies
- [ ] Register waiter routes
- [ ] Register cashier routes
- [ ] Register manager routes
- [ ] Test routes with Postman/curl
- [ ] Verify authentication works

---

## 🎨 PHASE 5: Frontend Services & Stores

### Task 5.1: Create Handover Service
**Create:** `frontend/src/services/handover.js`

**Implementation:**
```javascript
import api from './api'

export const handoverService = {
  // Waiter endpoints
  createHandover: (shiftId, data) => 
    api.post(`/api/shifts/${shiftId}/handover`, data),
  
  createHandoverAndEndShift: (shiftId, data) => 
    api.post(`/api/shifts/${shiftId}/handover-and-end`, data),
  
  getPendingHandover: (shiftId) => 
    api.get(`/api/shifts/${shiftId}/pending-handover`),
  
  getHandoverHistory: (shiftId) => 
    api.get(`/api/shifts/${shiftId}/handovers`),
  
  cancelHandover: (handoverId) => 
    api.delete(`/api/cash-handovers/${handoverId}`),
  
  // Cashier endpoints
  getPendingHandovers: () => 
    api.get('/api/cash-handovers/pending'),
  
  getTodayHandovers: () => 
    api.get('/api/cash-handovers/today'),
  
  confirmHandover: (handoverId, data) => 
    api.post(`/api/cash-handovers/${handoverId}/confirm`, data),
  
  quickConfirm: (handoverId, status) => 
    api.post(`/api/cash-handovers/${handoverId}/quick-confirm`, { status }),
  
  // Manager endpoints
  getPendingApprovals: () => 
    api.get('/api/cash-handovers/pending-approval'),
  
  approveDiscrepancy: (handoverId, data) => 
    api.post(`/api/cash-handovers/${handoverId}/approve`, data),
  
  getDiscrepancyStats: (startDate, endDate) => 
    api.get(`/api/cash-handovers/discrepancy-stats?start=${startDate}&end=${endDate}`)
}
```

**Checklist:**
- [ ] Create service file
- [ ] Implement all API calls
- [ ] Test API calls work
- [ ] Handle errors properly

---

### Task 5.2: Update Shift Store
**Modify:** `frontend/src/stores/shift.js`

**Add methods:**
```javascript
// Add to shift store
const createCashHandover = async (shiftId, handoverData) => {
  try {
    const response = await handoverService.createHandover(shiftId, handoverData)
    return response.data
  } catch (error) {
    console.error('Error creating cash handover:', error)
    throw error
  }
}

const createHandoverAndEndShift = async (shiftId, handoverData) => {
  try {
    const response = await handoverService.createHandoverAndEndShift(shiftId, handoverData)
    return response.data
  } catch (error) {
    console.error('Error creating handover and end shift:', error)
    throw error
  }
}

const getPendingHandover = async (shiftId) => {
  try {
    const response = await handoverService.getPendingHandover(shiftId)
    return response.data
  } catch (error) {
    console.error('Error fetching pending handover:', error)
    return null
  }
}

const getHandoverHistory = async (shiftId) => {
  try {
    const response = await handoverService.getHandoverHistory(shiftId)
    return response.data
  } catch (error) {
    console.error('Error fetching handover history:', error)
    return []
  }
}

const cancelHandover = async (handoverId) => {
  try {
    const response = await handoverService.cancelHandover(handoverId)
    return response.data
  } catch (error) {
    console.error('Error canceling handover:', error)
    throw error
  }
}

// Add to return statement
return {
  // ... existing
  createCashHandover,
  createHandoverAndEndShift,
  getPendingHandover,
  getHandoverHistory,
  cancelHandover
}
```

**Checklist:**
- [ ] Add handover methods to shift store
- [ ] Import handoverService
- [ ] Export new methods
- [ ] Test store methods work

---

### Task 5.3: Update Cashier Store
**Modify:** `frontend/src/stores/cashier.js`

**Add state and methods:**
```javascript
// Add to state
const pendingHandovers = ref([])
const todayHandovers = ref([])

// Add methods
const fetchPendingHandovers = async () => {
  loading.value = true
  try {
    const response = await handoverService.getPendingHandovers()
    pendingHandovers.value = response.data
  } catch (error) {
    console.error('Error fetching pending handovers:', error)
    throw error
  } finally {
    loading.value = false
  }
}

const fetchTodayHandovers = async () => {
  try {
    const response = await handoverService.getTodayHandovers()
    todayHandovers.value = response.data
  } catch (error) {
    console.error('Error fetching today handovers:', error)
    throw error
  }
}

const confirmHandover = async (handoverId, confirmData) => {
  try {
    const response = await handoverService.confirmHandover(handoverId, confirmData)
    return response.data
  } catch (error) {
    console.error('Error confirming handover:', error)
    throw error
  }
}

const quickConfirm = async (handoverId, status) => {
  try {
    const response = await handoverService.quickConfirm(handoverId, status)
    await fetchPendingHandovers()
    return response.data
  } catch (error) {
    console.error('Error quick confirming handover:', error)
    throw error
  }
}

// Add to return
return {
  // ... existing
  pendingHandovers,
  todayHandovers,
  fetchPendingHandovers,
  fetchTodayHandovers,
  confirmHandover,
  quickConfirm
}
```

**Checklist:**
- [ ] Add state variables
- [ ] Add handover methods
- [ ] Import handoverService
- [ ] Export new state and methods
- [ ] Test store methods work

---

## 🎨 PHASE 6: Frontend UI Components

### Task 6.1: Update ShiftView.vue (Waiter Interface)
**Modify:** `frontend/src/views/ShiftView.vue`

**Changes needed:**
1. Add cash status display (3 cards: Tiền hiện có, Đã bàn giao, Tổng thu)
2. Add pending handover banner
3. Add handover action buttons
4. Add handover history section
5. Add partial handover modal
6. Add handover-and-end-shift modal

**Script additions:**
```javascript
// Add refs
const showPartialHandoverForm = ref(false)
const showHandoverEndShiftForm = ref(false)
const pendingHandover = ref(null)
const handoverHistory = ref([])

const partialHandoverForm = ref({
  declared_amount: 0,
  waiter_note: ''
})

const handoverEndShiftForm = ref({
  end_cash: 0,
  waiter_note: ''
})

const isWaiter = computed(() => authStore.user?.role === 'waiter')

// Add methods
const fetchHandoverData = async () => { /* ... */ }
const createPartialHandover = async () => { /* ... */ }
const createHandoverAndEndShift = async () => { /* ... */ }
const cancelHandover = async (handoverId) => { /* ... */ }
const getHandoverTypeText = (type) => { /* ... */ }
const getHandoverStatusText = (status) => { /* ... */ }
const getHandoverStatusClass = (status) => { /* ... */ }
```

**Checklist:**
- [ ] Add cash status cards to template
- [ ] Add pending handover banner
- [ ] Add handover buttons
- [ ] Add handover history section
- [ ] Create partial handover modal
- [ ] Create handover-and-end-shift modal
- [ ] Add script logic
- [ ] Test UI interactions
- [ ] Test with real data

**Reference:** Copy full template and script from analysis doc sections.

---

### Task 6.2: Create CashierHandoverView.vue
**Create:** `frontend/src/views/CashierHandoverView.vue`

**Sections needed:**
1. Mobile header
2. Pending handovers section
3. Today's handovers section
4. Confirm/Reject modal with reconciliation

**Key features:**
- Display pending handovers with waiter info
- Show declared amount
- Buttons for confirm/reject
- Modal for entering actual amount
- Auto-calculate discrepancy
- Discrepancy reason selection
- Large discrepancy warning

**Checklist:**
- [ ] Create file
- [ ] Add mobile header
- [ ] Add pending handovers section
- [ ] Add today's handovers section
- [ ] Create confirm modal
- [ ] Add reconciliation logic
- [ ] Add discrepancy handling
- [ ] Test UI interactions
- [ ] Test with real data

**Reference:** Copy full template from analysis doc.

---

### Task 6.3: Update CashierDashboard.vue
**Modify:** `frontend/src/views/CashierDashboard.vue`

**Add sections:**
1. Handover notification banner (if pending > 0)
2. Quick handover actions section

**Template additions:**
```vue
<!-- Handover Notifications Section -->
<div v-if="pendingHandovers.length > 0" class="bg-yellow-50 border-l-4 border-yellow-400 p-4 mb-4">
  <div class="flex items-center justify-between">
    <div class="flex items-center">
      <div class="flex-shrink-0">
        <svg class="h-5 w-5 text-yellow-400" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
        </svg>
      </div>
      <div class="ml-3">
        <p class="text-sm text-yellow-700">
          <strong>{{ pendingHandovers.length }} yêu cầu bàn giao</strong> đang chờ xác nhận
        </p>
      </div>
    </div>
    <button @click="$router.push('/cashier/handovers')" 
      class="bg-yellow-500 hover:bg-yellow-600 text-white px-4 py-2 rounded-lg text-sm font-medium">
      Xem ngay
    </button>
  </div>
</div>

<!-- Quick Handover Actions -->
<div v-if="pendingHandovers.length > 0" class="bg-white rounded-2xl p-6 shadow-sm mb-4">
  <h3 class="text-lg font-bold mb-4">⚡ Bàn giao nhanh</h3>
  <div class="space-y-3">
    <div v-for="handover in pendingHandovers.slice(0, 3)" :key="handover.id" 
      class="flex items-center justify-between p-3 bg-gray-50 rounded-xl">
      <div>
        <p class="font-medium">{{ handover.waiter_name }}</p>
        <p class="text-sm text-gray-500">{{ formatPrice(handover.declared_amount) }}</p>
      </div>
      <div class="flex gap-2">
        <button @click="quickConfirm(handover.id, 'CONFIRMED')"
          class="bg-green-500 hover:bg-green-600 text-white px-3 py-1 rounded-lg text-sm">
          ✅
        </button>
        <button @click="quickConfirm(handover.id, 'REJECTED')"
          class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded-lg text-sm">
          ❌
        </button>
      </div>
    </div>
  </div>
</div>
```

**Script additions:**
```javascript
const pendingHandovers = computed(() => cashierStore.pendingHandovers)

onMounted(async () => {
  await cashierStore.fetchPendingHandovers()
})

const quickConfirm = async (handoverId, status) => {
  try {
    await cashierStore.quickConfirm(handoverId, status)
    alert(status === 'CONFIRMED' ? 'Đã xác nhận!' : 'Đã từ chối!')
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}
```

**Checklist:**
- [ ] Add notification banner
- [ ] Add quick actions section
- [ ] Add script logic
- [ ] Test notifications appear
- [ ] Test quick actions work

---

### Task 6.4: Add Router Configuration
**Modify:** `frontend/src/router/index.js`

**Add route:**
```javascript
{
  path: '/cashier/handovers',
  name: 'CashierHandovers',
  component: () => import('../views/CashierHandoverView.vue'),
  meta: { 
    requiresAuth: true, 
    roles: ['cashier', 'manager'] 
  }
}
```

**Checklist:**
- [ ] Add route
- [ ] Test navigation works
- [ ] Test role-based access

---

### Task 6.5: Update Navigation
**Modify:** `frontend/src/components/Navigation.vue` (or BottomNav.vue)

**Add for cashier:**
```vue
<router-link v-if="isCashier" to="/cashier/handovers" 
  class="flex items-center px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-lg">
  <span class="mr-3">💰</span>
  <span>Bàn giao tiền</span>
  <span v-if="pendingHandoversCount > 0" 
    class="ml-auto bg-red-500 text-white text-xs px-2 py-1 rounded-full">
    {{ pendingHandoversCount }}
  </span>
</router-link>
```

**Checklist:**
- [ ] Add navigation link
- [ ] Add badge for pending count
- [ ] Test navigation works
- [ ] Test badge updates

---

## 🧪 PHASE 7: Testing & Validation

### Task 7.1: Backend Unit Tests
**Create test files:**
- `backend/application/services/cash_handover_service_test.go`
- `backend/domain/handover/cash_handover_test.go`

**Test cases:**
1. CreateHandover - success
2. CreateHandover - exceeds remaining cash
3. CreateHandover - no cashier shift open
4. ConfirmHandover - no discrepancy
5. ConfirmHandover - with discrepancy
6. ConfirmHandover - large discrepancy requires approval
7. ApproveDiscrepancy - manager approval
8. UpdateCashAmounts - correct calculations

**Checklist:**
- [ ] Write unit tests
- [ ] Run tests: `go test ./...`
- [ ] All tests pass
- [ ] Coverage > 80%

---

### Task 7.2: API Integration Tests
**Test with Postman/curl:**

**Waiter flow:**
```bash
# 1. Create partial handover
POST /api/shifts/{shift_id}/handover
{
  "declared_amount": 500000,
  "handover_type": "PARTIAL",
  "waiter_note": "Test handover"
}

# 2. Get pending handover
GET /api/shifts/{shift_id}/pending-handover

# 3. Get handover history
GET /api/shifts/{shift_id}/handovers

# 4. Cancel handover
DELETE /api/cash-handovers/{handover_id}
```

**Cashier flow:**
```bash
# 1. Get pending handovers
GET /api/cash-handovers/pending

# 2. Confirm handover
POST /api/cash-handovers/{handover_id}/confirm
{
  "actual_amount": 500000,
  "status": "CONFIRMED",
  "cashier_note": "Confirmed"
}

# 3. Confirm with discrepancy
POST /api/cash-handovers/{handover_id}/confirm
{
  "actual_amount": 480000,
  "status": "CONFIRMED",
  "cashier_note": "Short 20k",
  "discrepancy_reason": "COUNTING_ERROR",
  "discrepancy_responsibility": "WAITER"
}
```

**Checklist:**
- [ ] Test all waiter endpoints
- [ ] Test all cashier endpoints
- [ ] Test all manager endpoints
- [ ] Test error cases
- [ ] Test validation
- [ ] Document API responses

---

### Task 7.3: Frontend E2E Tests
**Test scenarios:**

**Waiter:**
1. Open shift → See cash status
2. Create partial handover → See pending banner
3. Cancel handover → Banner disappears
4. Create handover and end shift → Shift closes after cashier confirms
5. View handover history

**Cashier:**
1. See notification when handover pending
2. View pending handovers list
3. Confirm handover without discrepancy
4. Confirm handover with small discrepancy
5. Confirm handover with large discrepancy → Requires manager approval
6. Reject handover with reason
7. Quick confirm from dashboard

**Checklist:**
- [ ] Test waiter flow
- [ ] Test cashier flow
- [ ] Test manager approval flow
- [ ] Test error handling
- [ ] Test UI responsiveness
- [ ] Test on mobile devices

---

## 📊 PHASE 8: Database Migration & Indexes

### Task 8.1: Add Indexes
**Run in MongoDB:**
```javascript
// Cash handovers collection
db.cash_handovers.createIndex({ "waiter_shift_id": 1 })
db.cash_handovers.createIndex({ "cashier_shift_id": 1 })
db.cash_handovers.createIndex({ "status": 1 })
db.cash_handovers.createIndex({ "handover_at": -1 })
db.cash_handovers.createIndex({ "cashier_id": 1, "status": 1 })
db.cash_handovers.createIndex({ "requires_approval": 1, "status": 1 })

// Cash discrepancies collection
db.cash_discrepancies.createIndex({ "handover_id": 1 })
db.cash_discrepancies.createIndex({ "resolution_status": 1 })
db.cash_discrepancies.createIndex({ "requires_manager_approval": 1 })
db.cash_discrepancies.createIndex({ "created_at": -1 })
```

**Checklist:**
- [ ] Create indexes
- [ ] Verify index creation
- [ ] Test query performance

---

### Task 8.2: Update Existing Shifts
**Migration script:**
```javascript
// Add new fields to existing shifts
db.shifts.updateMany(
  { type: "WAITER" },
  {
    $set: {
      current_cash: 0,
      handed_over_cash: 0,
      remaining_cash: 0,
      total_discrepancy: 0,
      handover_count: 0
    }
  }
)

// Add new fields to existing cashier shifts
db.cashier_shifts.updateMany(
  {},
  {
    $set: {
      received_cash: 0,
      total_discrepancy: 0,
      handover_count: 0,
      discrepancy_count: 0
    }
  }
)
```

**Checklist:**
- [ ] Run migration script
- [ ] Verify all shifts updated
- [ ] Test with existing data

---

## 📝 PHASE 9: Documentation

### Task 9.1: API Documentation
**Create:** `docs/API_CASH_HANDOVER.md`

**Content:**
- All endpoints with request/response examples
- Authentication requirements
- Error codes and messages
- Rate limiting (if any)

**Checklist:**
- [ ] Document all endpoints
- [ ] Add request examples
- [ ] Add response examples
- [ ] Add error cases

---

### Task 9.2: User Guide
**Create:** `docs/USER_GUIDE_CASH_HANDOVER.md`

**Content:**
- Waiter guide: How to handover cash
- Cashier guide: How to confirm handovers
- Manager guide: How to approve discrepancies
- Screenshots/diagrams

**Checklist:**
- [ ] Write waiter guide
- [ ] Write cashier guide
- [ ] Write manager guide
- [ ] Add screenshots

---

### Task 9.3: Update INDEX.md
**Modify:** `docs/INDEX.md`

**Add entry:**
```markdown
### 💰 Quản Lý Bàn Giao Tiền
- **[CASH_HANDOVER_FEATURE_ANALYSIS.md](./CASH_HANDOVER_FEATURE_ANALYSIS.md)** - Phân tích tính năng bàn giao tiền
- **[CASH_HANDOVER_IMPLEMENTATION_TASKS.md](./CASH_HANDOVER_IMPLEMENTATION_TASKS.md)** - Tasks triển khai
- **[API_CASH_HANDOVER.md](./API_CASH_HANDOVER.md)** - API documentation
- **[USER_GUIDE_CASH_HANDOVER.md](./USER_GUIDE_CASH_HANDOVER.md)** - Hướng dẫn sử dụng
```

**Checklist:**
- [ ] Update INDEX.md
- [ ] Verify links work

---

## ✅ PHASE 10: Final Verification

### Task 10.1: Complete Flow Test
**Test complete workflow:**
1. Waiter opens shift
2. Waiter receives cash from orders
3. Waiter creates partial handover
4. Cashier receives notification
5. Cashier confirms with exact amount
6. Both shifts updated correctly
7. Waiter creates another handover
8. Cashier confirms with discrepancy
9. Manager receives approval request
10. Manager approves
11. Waiter creates handover and end shift
12. Cashier confirms
13. Waiter shift closes automatically

**Checklist:**
- [ ] Test complete flow
- [ ] Verify all data correct
- [ ] Check audit trail
- [ ] Verify cash amounts
- [ ] Test edge cases

---

### Task 10.2: Performance Testing
**Test scenarios:**
- [ ] 100 concurrent handover requests
- [ ] Large handover history (1000+ records)
- [ ] Multiple pending handovers
- [ ] Database query performance
- [ ] Frontend rendering performance

---

### Task 10.3: Security Audit
**Verify:**
- [ ] Authentication on all endpoints
- [ ] Role-based authorization
- [ ] Input validation
- [ ] SQL injection prevention
- [ ] XSS prevention
- [ ] CSRF protection

---

## 🚀 Deployment Checklist

### Pre-deployment:
- [ ] All tests pass
- [ ] Code review completed
- [ ] Documentation updated
- [ ] Database migration script ready
- [ ] Rollback plan prepared

### Deployment:
- [ ] Backup database
- [ ] Run migration script
- [ ] Deploy backend
- [ ] Deploy frontend
- [ ] Verify deployment
- [ ] Monitor for errors

### Post-deployment:
- [ ] Test in production
- [ ] Monitor performance
- [ ] Check error logs
- [ ] User acceptance testing
- [ ] Gather feedback

---

## 📊 Success Metrics

**Functional:**
- [ ] Waiter can create handover requests
- [ ] Cashier can confirm/reject handovers
- [ ] Discrepancies are tracked correctly
- [ ] Manager approval works for large discrepancies
- [ ] Cash amounts update correctly
- [ ] Audit trail is complete

**Performance:**
- [ ] API response time < 500ms
- [ ] Frontend load time < 2s
- [ ] No memory leaks
- [ ] Database queries optimized

**User Experience:**
- [ ] UI is intuitive
- [ ] Mobile responsive
- [ ] Error messages clear
- [ ] Loading states visible
- [ ] Success feedback immediate

---

## 🐛 Known Issues & Future Enhancements

### Known Issues:
- None yet (to be updated during implementation)

### Future Enhancements:
1. Real-time notifications using WebSocket
2. Bulk handover operations
3. Advanced analytics dashboard
4. Export handover reports to PDF/Excel
5. SMS notifications for large discrepancies
6. Biometric authentication for confirmations
7. Integration with accounting systems

---

## 📞 Support & Contact

**For implementation questions:**
- Check analysis doc: `CASH_HANDOVER_FEATURE_ANALYSIS.md`
- Review existing code patterns
- Ask team lead

**For bugs:**
- Create issue with reproduction steps
- Include error logs
- Tag as `handover` and `bug`

---

**Last Updated:** 2026-02-04
**Status:** Ready for Implementation
**Estimated Time:** 5-7 days (1 developer)
