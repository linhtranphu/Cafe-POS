# Task Implementation Plan: Tách riêng Cashier Shifts và Waiter Shifts

## Phân tích vấn đề

### Tình trạng hiện tại:
- **1 collection `shifts`** dùng chung cho waiter, cashier, barista
- Phân biệt bằng field `role_type` (waiter/cashier/barista)
- Logic đơn giản: start shift → end shift
- Không có logic đối soát tiền mặt phức tạp cho cashier
- Cashier shift và waiter shift bị nhầm lẫn

### Mục tiêu:
- **2 hệ thống shift riêng biệt:**
  1. **Waiter/Barista Shifts** - Ca làm việc thông thường
  2. **Cashier Shifts** - Ca thu ngân với logic đối soát phức tạp
- Cashier chỉ đóng ca khi tất cả waiter shifts đã đóng
- Mỗi loại shift có collection và domain model riêng

---

## Implementation Tasks

### Phase 1: Backend - Domain Layer

#### Task 1.1: Tạo CashierShift Domain Model
**File:** `backend/domain/cashier/cashier_shift.go`

**Nội dung:**
```go
package cashier

import (
    "errors"
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type CashierShiftStatus string

const (
    CashierShiftOpen              CashierShiftStatus = "OPEN"
    CashierShiftClosureInitiated  CashierShiftStatus = "CLOSURE_INITIATED"
    CashierShiftClosed            CashierShiftStatus = "CLOSED"
)

type CashierShift struct {
    ID              primitive.ObjectID           `bson:"_id,omitempty" json:"id"`
    CashierID       primitive.ObjectID           `bson:"cashier_id" json:"cashier_id"`
    CashierName     string                       `bson:"cashier_name" json:"cashier_name"`
    StartTime       time.Time                    `bson:"start_time" json:"start_time"`
    EndTime         *time.Time                   `bson:"end_time,omitempty" json:"end_time,omitempty"`
    Status          CashierShiftStatus           `bson:"status" json:"status"`
    StartingFloat   float64                      `bson:"starting_float" json:"starting_float"`
    SystemCash      float64                      `bson:"system_cash" json:"system_cash"`
    ActualCash      *float64                     `bson:"actual_cash,omitempty" json:"actual_cash,omitempty"`
    Variance        *Variance                    `bson:"variance,omitempty" json:"variance,omitempty"`
    Confirmation    *ResponsibilityConfirmation  `bson:"confirmation,omitempty" json:"confirmation,omitempty"`
    AuditLog        []AuditLogEntry              `bson:"audit_log" json:"audit_log"`
    CreatedAt       time.Time                    `bson:"created_at" json:"created_at"`
    UpdatedAt       time.Time                    `bson:"updated_at" json:"updated_at"`
}

// Domain methods
func NewCashierShift(cashierID primitive.ObjectID, cashierName string, startingFloat float64) *CashierShift
func (cs *CashierShift) InitiateClosure(userID, deviceID string, timestamp time.Time) error
func (cs *CashierShift) RecordActualCash(actualCash float64, userID, deviceID string, timestamp time.Time) (*Variance, error)
func (cs *CashierShift) DocumentVariance(reason VarianceReason, notes string, userID, deviceID string, timestamp time.Time) error
func (cs *CashierShift) ConfirmResponsibility(userID, deviceID string, timestamp time.Time) error
func (cs *CashierShift) CanClose() error
func (cs *CashierShift) Close(userID, deviceID string, timestamp time.Time) error
```

**Lý do:** Tạo domain model riêng cho cashier shift với logic phức tạp

---

#### Task 1.2: Giữ nguyên Waiter/Barista Shift
**File:** `backend/domain/order/shift.go`

**Thay đổi:**
- Xóa `RoleCashier` khỏi `RoleType` enum
- Chỉ giữ `RoleWaiter` và `RoleBarista`
- Xóa các field `CashierID`, `CashierName` (legacy)
- Đơn giản hóa logic: chỉ start/end shift

```go
type RoleType string

const (
    RoleWaiter  RoleType = "waiter"
    RoleBarista RoleType = "barista"
)
```

**Lý do:** Tách rời logic cashier khỏi waiter/barista shifts

---

### Phase 2: Backend - Repository Layer

#### Task 2.1: Tạo CashierShiftRepository
**File:** `backend/infrastructure/mongodb/cashier_shift_repository.go`

**Nội dung:**
```go
type CashierShiftRepository struct {
    collection *mongo.Collection
}

func NewCashierShiftRepository(db *mongo.Database) *CashierShiftRepository {
    return &CashierShiftRepository{
        collection: db.Collection("cashier_shifts"), // Collection riêng
    }
}

// Methods
func (r *CashierShiftRepository) Create(ctx context.Context, shift *cashier.CashierShift) error
func (r *CashierShiftRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*cashier.CashierShift, error)
func (r *CashierShiftRepository) Save(ctx context.Context, shift *cashier.CashierShift) error
func (r *CashierShiftRepository) FindOpenByCashier(ctx context.Context, cashierID primitive.ObjectID) (*cashier.CashierShift, error)
func (r *CashierShiftRepository) FindAll(ctx context.Context) ([]*cashier.CashierShift, error)
func (r *CashierShiftRepository) FindByCashierID(ctx context.Context, cashierID primitive.ObjectID) ([]*cashier.CashierShift, error)
```

**Indexes:**
```javascript
db.cashier_shifts.createIndex({ "cashier_id": 1, "start_time": -1 })
db.cashier_shifts.createIndex({ "status": 1 })
db.cashier_shifts.createIndex({ "end_time": -1 })
```

**Lý do:** Collection riêng cho cashier shifts, tách biệt hoàn toàn với waiter shifts

---

#### Task 2.2: Cập nhật ShiftRepository
**File:** `backend/infrastructure/mongodb/shift_repository.go`

**Thay đổi:**
- Xóa các method liên quan đến cashier
- Chỉ giữ method cho waiter và barista
- Xóa `FindOpenShiftByUser` với roleType cashier

**Lý do:** Tách rời logic repository

---

### Phase 3: Backend - Service Layer

#### Task 3.1: Tạo CashierShiftService
**File:** `backend/application/services/cashier_shift_service.go`

**Nội dung:**
```go
type CashierShiftService struct {
    cashierShiftRepo *mongodb.CashierShiftRepository
    waiterShiftRepo  *mongodb.ShiftRepository // Để kiểm tra waiter shifts
}

func NewCashierShiftService(
    cashierShiftRepo *mongodb.CashierShiftRepository,
    waiterShiftRepo *mongodb.ShiftRepository,
) *CashierShiftService

// Methods
func (s *CashierShiftService) StartCashierShift(ctx context.Context, cashierID primitive.ObjectID, cashierName string, startingFloat float64) (*cashier.CashierShift, error)
func (s *CashierShiftService) GetCurrentCashierShift(ctx context.Context, cashierID primitive.ObjectID) (*cashier.CashierShift, error)
func (s *CashierShiftService) GetCashierShiftsByUser(ctx context.Context, cashierID primitive.ObjectID) ([]*cashier.CashierShift, error)
func (s *CashierShiftService) GetAllCashierShifts(ctx context.Context) ([]*cashier.CashierShift, error)

// Shift closure methods (sử dụng ShiftClosureService đã có)
```

**Logic quan trọng:**
- `StartCashierShift`: Kiểm tra không có cashier shift nào đang open
- Không có method `EndCashierShift` đơn giản - phải dùng shift closure workflow

**Lý do:** Service riêng cho cashier shifts với logic nghiệp vụ riêng

---

#### Task 3.2: Cập nhật ShiftClosureService
**File:** `backend/application/services/shift_closure_service.go`

**Thay đổi:**
```go
type ShiftClosureService struct {
    cashierShiftRepo *mongodb.CashierShiftRepository
    waiterShiftRepo  *mongodb.ShiftRepository // Để kiểm tra waiter shifts
    shiftReportRepo  *mongodb.ShiftReportRepository
}

func (s *ShiftClosureService) InitiateShiftClosure(ctx context.Context, cashierShiftID primitive.ObjectID) (*ShiftSummary, error) {
    // 1. Load cashier shift
    cashierShift, err := s.cashierShiftRepo.FindByID(ctx, cashierShiftID)
    
    // 2. Kiểm tra tất cả waiter shifts đã đóng
    openWaiterShifts, err := s.waiterShiftRepo.FindOpenShifts(ctx)
    if len(openWaiterShifts) > 0 {
        return nil, errors.New("cannot close cashier shift: waiter shifts are still open")
    }
    
    // 3. Tiếp tục logic closure...
}
```

**Lý do:** Kết nối 2 hệ thống shift - cashier chỉ đóng khi waiter shifts đã đóng

---

#### Task 3.3: Cập nhật ShiftService
**File:** `backend/application/services/shift_service.go`

**Thay đổi:**
- Xóa logic xử lý cashier từ `StartShift`
- Chỉ xử lý waiter và barista
- Thêm validation: reject nếu roleType là cashier

```go
func (s *ShiftService) StartShift(ctx context.Context, req *order.StartShiftRequest, userID, userName string, roleType order.RoleType) (*order.Shift, error) {
    // Reject cashier role
    if roleType == order.RoleCashier {
        return nil, errors.New("use cashier shift service for cashier shifts")
    }
    
    // Chỉ xử lý waiter và barista
    // ...
}
```

**Lý do:** Tách rời logic xử lý cashier shifts

---

### Phase 4: Backend - API Layer

#### Task 4.1: Tạo CashierShiftHandler
**File:** `backend/interfaces/http/cashier_shift_handler.go`

**Endpoints:**
```
POST   /api/v1/cashier-shifts              - Start cashier shift
GET    /api/v1/cashier-shifts/current      - Get current cashier shift
GET    /api/v1/cashier-shifts              - Get all cashier shifts
GET    /api/v1/cashier-shifts/:id          - Get cashier shift by ID
GET    /api/v1/cashier-shifts/my-shifts    - Get my cashier shifts
```

**Lý do:** API riêng cho cashier shifts

---

#### Task 4.2: Giữ nguyên ShiftClosureHandler
**File:** `backend/interfaces/http/shift_closure_handler.go`

**Không thay đổi** - đã đúng logic:
```
POST   /api/v1/cashier-shifts/:id/initiate-closure
POST   /api/v1/cashier-shifts/:id/record-actual-cash
POST   /api/v1/cashier-shifts/:id/document-variance
POST   /api/v1/cashier-shifts/:id/confirm-responsibility
POST   /api/v1/cashier-shifts/:id/close
GET    /api/v1/shift-reports/:id
```

---

#### Task 4.3: Cập nhật ShiftHandler
**File:** `backend/interfaces/http/shift_handler.go`

**Thay đổi:**
- Thêm validation: reject nếu role là cashier
- Chỉ xử lý waiter và barista

```go
func (h *ShiftHandler) StartShift(c *gin.Context) {
    roleType := order.ParseRoleType(string(role.(user.Role)))
    
    // Reject cashier
    if roleType == order.RoleCashier {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Use /api/v1/cashier-shifts endpoint for cashier shifts"
        })
        return
    }
    
    // Continue with waiter/barista logic
}
```

---

#### Task 4.4: Cập nhật Routes
**File:** `backend/main.go`

**Thêm routes:**
```go
// Cashier shift routes (separate from waiter shifts)
cashierShiftHandler := http.NewCashierShiftHandler(cashierShiftService)
api.POST("/cashier-shifts", authMiddleware, cashierShiftHandler.StartCashierShift)
api.GET("/cashier-shifts/current", authMiddleware, cashierShiftHandler.GetCurrentCashierShift)
api.GET("/cashier-shifts", authMiddleware, cashierShiftHandler.GetAllCashierShifts)
api.GET("/cashier-shifts/:id", authMiddleware, cashierShiftHandler.GetCashierShift)
api.GET("/cashier-shifts/my-shifts", authMiddleware, cashierShiftHandler.GetMyCashierShifts)

// Shift closure routes (already exist)
shiftClosureHandler := http.NewShiftClosureHandler(shiftClosureService)
api.POST("/cashier-shifts/:id/initiate-closure", authMiddleware, shiftClosureHandler.InitiateShiftClosure)
// ... other closure endpoints
```

---

### Phase 5: Frontend - Services

#### Task 5.1: Tạo cashierShift.js service
**File:** `frontend/src/services/cashierShift.js`

**Nội dung:**
```javascript
import api from './api'

export default {
  // Start cashier shift
  async startCashierShift(startingFloat) {
    const response = await api.post('/v1/cashier-shifts', {
      starting_float: startingFloat
    })
    return response.data
  },

  // Get current cashier shift
  async getCurrentCashierShift() {
    const response = await api.get('/v1/cashier-shifts/current')
    return response.data
  },

  // Get all cashier shifts
  async getAllCashierShifts() {
    const response = await api.get('/v1/cashier-shifts')
    return response.data
  },

  // Get my cashier shifts
  async getMyCashierShifts() {
    const response = await api.get('/v1/cashier-shifts/my-shifts')
    return response.data
  },

  // Get cashier shift by ID
  async getCashierShift(id) {
    const response = await api.get(`/v1/cashier-shifts/${id}`)
    return response.data
  }
}
```

**Lý do:** Service riêng cho cashier shifts API

---

#### Task 5.2: Cập nhật shift.js service
**File:** `frontend/src/services/shift.js`

**Thay đổi:**
- Chỉ xử lý waiter và barista shifts
- Xóa logic cashier

---

### Phase 6: Frontend - Stores

#### Task 6.1: Tạo cashierShift.js store
**File:** `frontend/src/stores/cashierShift.js`

**Nội dung:**
```javascript
import { defineStore } from 'pinia'
import cashierShiftService from '../services/cashierShift'

export const useCashierShiftStore = defineStore('cashierShift', {
  state: () => ({
    currentCashierShift: null,
    cashierShifts: [],
    loading: false,
    error: null
  }),

  getters: {
    hasOpenCashierShift: (state) => {
      return state.currentCashierShift && state.currentCashierShift.status === 'OPEN'
    },
    canStartCashierShift: (state) => {
      return !state.currentCashierShift || state.currentCashierShift.status === 'CLOSED'
    }
  },

  actions: {
    async startCashierShift(startingFloat) {
      this.loading = true
      this.error = null
      try {
        const shift = await cashierShiftService.startCashierShift(startingFloat)
        this.currentCashierShift = shift
        return shift
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to start cashier shift'
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchCurrentCashierShift() {
      this.loading = true
      this.error = null
      try {
        const shift = await cashierShiftService.getCurrentCashierShift()
        this.currentCashierShift = shift
        return shift
      } catch (error) {
        this.currentCashierShift = null
        // Don't set error if no shift found (404 is expected)
        if (error.response?.status !== 404) {
          this.error = error.response?.data?.message || 'Failed to fetch current cashier shift'
        }
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchAllCashierShifts() {
      this.loading = true
      this.error = null
      try {
        const shifts = await cashierShiftService.getAllCashierShifts()
        this.cashierShifts = shifts
        return shifts
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to fetch cashier shifts'
        throw error
      } finally {
        this.loading = false
      }
    }
  }
})
```

**Lý do:** Store riêng cho cashier shifts state management

---

#### Task 6.2: Cập nhật shift.js store
**File:** `frontend/src/stores/shift.js`

**Thay đổi:**
- Xóa logic cashier shifts
- Chỉ xử lý waiter và barista shifts

---

### Phase 7: Frontend - UI Components

#### Task 7.1: Tạo CashierShiftManager Component
**File:** `frontend/src/components/CashierShiftManager.vue`

**Chức năng:**
- Hiển thị cashier shift hiện tại
- Button "Bắt đầu ca thu ngân" (nếu chưa có ca)
- Button "Đóng ca thu ngân" (nếu có ca đang mở)
- Hiển thị thông tin ca: thời gian bắt đầu, tiền đầu ca, trạng thái

**Lý do:** Component riêng để quản lý cashier shifts

---

#### Task 7.2: Cập nhật CashierDashboard
**File:** `frontend/src/views/CashierDashboard.vue`

**Thay đổi:**
```vue
<template>
  <div>
    <!-- Cashier Shift Manager -->
    <CashierShiftManager />
    
    <!-- Shift Selector - CHỈ hiển thị CASHIER shifts -->
    <select v-model="selectedCashierShift">
      <option value="">-- Chọn ca thu ngân --</option>
      <option v-for="shift in cashierShifts" :key="shift.id" :value="shift.id">
        {{ formatDate(shift.start_time) }} - {{ shift.cashier_name }}
      </option>
    </select>
    
    <!-- Close Shift Button -->
    <button v-if="canCloseCashierShift" @click="goToShiftClosure">
      🔒 Đóng ca thu ngân
    </button>
    
    <!-- Rest of dashboard -->
  </div>
</template>

<script setup>
import { useCashierShiftStore } from '../stores/cashierShift'

const cashierShiftStore = useCashierShiftStore()

const cashierShifts = computed(() => cashierShiftStore.cashierShifts)
const currentCashierShift = computed(() => cashierShiftStore.currentCashierShift)

const canCloseCashierShift = computed(() => {
  return currentCashierShift.value && currentCashierShift.value.status === 'OPEN'
})

onMounted(async () => {
  await cashierShiftStore.fetchCurrentCashierShift()
  await cashierShiftStore.fetchAllCashierShifts()
})
</script>
```

**Lý do:** Chỉ hiển thị cashier shifts, không còn nhầm lẫn với waiter shifts

---

#### Task 7.3: Cập nhật ShiftView
**File:** `frontend/src/views/ShiftView.vue`

**Thay đổi:**
- Chỉ hiển thị waiter và barista shifts
- Xóa logic cashier shifts
- Filter: `shifts.filter(s => s.role_type !== 'cashier')`

**Lý do:** Tách rời UI cho waiter/barista shifts

---

### Phase 8: Database Migration

#### Task 8.1: Tạo Migration Script
**File:** `backend/cmd/migrate/separate_cashier_shifts.go`

**Nội dung:**
```go
// Migration: Tách cashier shifts từ collection 'shifts' sang 'cashier_shifts'

func MigrateCashierShifts(db *mongo.Database) error {
    shiftsCollection := db.Collection("shifts")
    cashierShiftsCollection := db.Collection("cashier_shifts")
    
    // 1. Find all cashier shifts
    cursor, err := shiftsCollection.Find(context.Background(), bson.M{
        "role_type": "cashier",
    })
    
    // 2. Transform and insert into cashier_shifts
    for cursor.Next(context.Background()) {
        var oldShift order.Shift
        cursor.Decode(&oldShift)
        
        newCashierShift := &cashier.CashierShift{
            ID:            oldShift.ID,
            CashierID:     oldShift.UserID,
            CashierName:   oldShift.UserName,
            StartTime:     oldShift.StartedAt,
            EndTime:       oldShift.EndedAt,
            Status:        convertStatus(oldShift.Status),
            StartingFloat: oldShift.StartCash,
            SystemCash:    oldShift.EndCash, // Approximate
            AuditLog:      []cashier.AuditLogEntry{},
            CreatedAt:     oldShift.CreatedAt,
            UpdatedAt:     oldShift.UpdatedAt,
        }
        
        cashierShiftsCollection.InsertOne(context.Background(), newCashierShift)
    }
    
    // 3. Delete cashier shifts from old collection
    shiftsCollection.DeleteMany(context.Background(), bson.M{
        "role_type": "cashier",
    })
    
    return nil
}
```

**Lý do:** Migrate dữ liệu cũ sang collection mới

---

## Summary

### Collections sau khi tách:

1. **`shifts`** - Waiter và Barista shifts
   - Dùng cho waiter và barista
   - Logic đơn giản: start/end
   - Không có đối soát tiền mặt

2. **`cashier_shifts`** - Cashier shifts
   - Chỉ dùng cho cashier
   - Logic phức tạp: closure workflow
   - Có đối soát tiền mặt, variance, confirmation

### API Endpoints:

**Waiter/Barista Shifts:**
```
POST   /api/v1/shifts              - Start waiter/barista shift
POST   /api/v1/shifts/:id/end      - End waiter/barista shift
GET    /api/v1/shifts/current      - Get current shift
GET    /api/v1/shifts              - Get all shifts
```

**Cashier Shifts:**
```
POST   /api/v1/cashier-shifts                        - Start cashier shift
GET    /api/v1/cashier-shifts/current                - Get current cashier shift
GET    /api/v1/cashier-shifts                        - Get all cashier shifts
POST   /api/v1/cashier-shifts/:id/initiate-closure   - Initiate closure
POST   /api/v1/cashier-shifts/:id/record-actual-cash - Record actual cash
POST   /api/v1/cashier-shifts/:id/document-variance  - Document variance
POST   /api/v1/cashier-shifts/:id/confirm-responsibility - Confirm responsibility
POST   /api/v1/cashier-shifts/:id/close              - Close shift
GET    /api/v1/shift-reports/:id                     - Get shift report
```

### Ưu điểm:

1. ✅ **Tách biệt rõ ràng** - Không còn nhầm lẫn giữa cashier và waiter shifts
2. ✅ **Logic riêng** - Mỗi loại shift có domain model và business logic riêng
3. ✅ **Dễ maintain** - Code rõ ràng, dễ hiểu, dễ test
4. ✅ **Scalable** - Dễ mở rộng thêm tính năng cho từng loại shift
5. ✅ **Đúng nghiệp vụ** - Cashier chỉ đóng ca khi tất cả waiter shifts đã đóng

---

## Execution Order

1. **Phase 1-2:** Backend Domain & Repository (2-3 giờ)
2. **Phase 3:** Backend Service Layer (2-3 giờ)
3. **Phase 4:** Backend API Layer (1-2 giờ)
4. **Phase 5-6:** Frontend Services & Stores (1-2 giờ)
5. **Phase 7:** Frontend UI Components (2-3 giờ)
6. **Phase 8:** Database Migration (1 giờ)

**Total estimate:** 9-14 giờ

---

## Testing Checklist

- [ ] Cashier có thể start cashier shift
- [ ] Cashier không thể start 2 shifts cùng lúc
- [ ] Waiter có thể start waiter shift độc lập
- [ ] Cashier không thể đóng ca khi còn waiter shifts mở
- [ ] Cashier có thể đóng ca khi tất cả waiter shifts đã đóng
- [ ] UI hiển thị đúng cashier shifts trong CashierDashboard
- [ ] UI hiển thị đúng waiter shifts trong ShiftView
- [ ] Migration script chạy thành công
- [ ] Tất cả API endpoints hoạt động đúng
