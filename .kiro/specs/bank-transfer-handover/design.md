# Design Document - Bank Transfer Handover

## Overview

Tính năng này mở rộng hệ thống bàn giao tiền hiện tại để hỗ trợ cả tiền chuyển khoản (bank transfer) bên cạnh tiền mặt. Thiết kế tập trung vào việc:

1. **Phân tách rõ ràng**: Theo dõi riêng biệt tiền mặt và tiền chuyển khoản trong toàn bộ quy trình
2. **Quy trình đối soát**: Cashier đối soát với tài khoản ngân hàng bên ngoài trước khi xác nhận
3. **Xác nhận đồng thời**: Xác nhận cả hai loại tiền trong một giao dịch để đảm bảo tính nguyên tử
4. **Tương thích ngược**: Duy trì khả năng hoạt động với quy trình bàn giao tiền mặt hiện tại

### Key Design Decisions

- **Separate Fields**: Sử dụng các trường riêng biệt cho cash và transfer thay vì một trường chung để tránh nhầm lẫn
- **Atomic Updates**: Cập nhật cả cash và transfer trong một transaction để đảm bảo tính nhất quán
- **Backward Compatible**: Các trường transfer là optional, handover chỉ có cash vẫn hoạt động bình thường
- **External Reconciliation**: Hệ thống không tự động đối soát với ngân hàng, Cashier phải đối soát thủ công

## Architecture

### System Components

```
┌─────────────────────────────────────────────────────────────┐
│                     Frontend Layer                          │
├─────────────────────────────────────────────────────────────┤
│  ShiftView (Waiter)                                         │
│  - Create handover with cash + transfer amounts             │
│  - View handover history                                    │
│                                                             │
│  CashierShiftView (Cashier)                                 │
│  - View pending handovers (cash + transfer)                 │
│  - Confirm with actual amounts for both                     │
│  - External bank reconciliation                             │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                     API Layer                               │
├─────────────────────────────────────────────────────────────┤
│  POST /api/waiter/shifts/:id/handover                       │
│  - Create handover with cash + transfer                     │
│                                                             │
│  POST /api/waiter/shifts/:id/handover-and-end               │
│  - Create end-shift handover with both amounts              │
│                                                             │
│  PATCH /api/cashier/handovers/:id/confirm                   │
│  - Confirm with actual_cash and actual_transfer             │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                  Service Layer                              │
├─────────────────────────────────────────────────────────────┤
│  CashHandoverService                                        │
│  - CreateHandoverWithTransfer()                             │
│  - ConfirmHandoverWithDualAmounts()                         │
│  - CalculateSeparateDiscrepancies()                         │
│  - UpdateDualBalances()                                     │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   Domain Layer                              │
├─────────────────────────────────────────────────────────────┤
│  CashHandover (Extended)                                    │
│  - cash_declared_amount                                     │
│  - transfer_declared_amount                                 │
│  - cash_actual_amount                                       │
│  - transfer_actual_amount                                   │
│  - cash_discrepancy                                         │
│  - transfer_discrepancy                                     │
│                                                             │
│  Shift (Extended)                                           │
│  - transfer_revenue                                         │
│  - remaining_transfer                                       │
│  - handed_over_transfer                                     │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                  Repository Layer                           │
├─────────────────────────────────────────────────────────────┤
│  MongoDB Collections:                                       │
│  - shifts (with transfer fields)                            │
│  - cash_handovers (with transfer fields)                    │
│  - cash_discrepancies (with transfer fields)                │
└─────────────────────────────────────────────────────────────┘
```

### Data Flow

#### Waiter Creates Handover with Transfer

```
1. Waiter opens handover form
2. System displays:
   - Remaining cash: X VND
   - Total transfer collected: Y VND
3. Waiter selects handover type:
   - Cash only
   - Transfer only
   - Both cash and transfer
4. Waiter enters amounts:
   - Cash amount (if selected)
   - Transfer amount (if selected)
5. System validates:
   - Cash amount <= remaining_cash
   - Transfer amount <= transfer_revenue
6. System creates handover record with:
   - cash_declared_amount
   - transfer_declared_amount
   - status = PENDING
7. System updates shift:
   - remaining_cash -= cash_declared_amount
   - remaining_transfer -= transfer_declared_amount
```

#### Cashier Confirms Handover

```
1. Cashier views pending handover
2. System displays:
   - Declared cash: X VND
   - Declared transfer: Y VND
3. Cashier performs external reconciliation:
   - Counts physical cash
   - Checks bank account for transfers
4. Cashier enters actual amounts:
   - Actual cash: X' VND
   - Actual transfer: Y' VND
5. System calculates discrepancies:
   - cash_discrepancy = X' - X
   - transfer_discrepancy = Y' - Y
   - total_discrepancy = cash_discrepancy + transfer_discrepancy
6. If total_discrepancy > threshold:
   - Require manager approval
   - status = DISCREPANCY
7. Else:
   - Update shift balances atomically
   - Update cashier shift balances
   - status = CONFIRMED
```

## Components and Interfaces

### Backend Components

#### 1. Extended Domain Models

**CashHandover (Extended)**
```go
type CashHandover struct {
    // Existing fields...
    ID             primitive.ObjectID
    WaiterShiftID  primitive.ObjectID
    CashierShiftID primitive.ObjectID
    
    // NEW: Separate cash amounts
    CashDeclaredAmount    float64 `bson:"cash_declared_amount" json:"cash_declared_amount"`
    CashActualAmount      float64 `bson:"cash_actual_amount" json:"cash_actual_amount"`
    CashDiscrepancy       float64 `bson:"cash_discrepancy" json:"cash_discrepancy"`
    
    // NEW: Separate transfer amounts
    TransferDeclaredAmount float64 `bson:"transfer_declared_amount" json:"transfer_declared_amount"`
    TransferActualAmount   float64 `bson:"transfer_actual_amount" json:"transfer_actual_amount"`
    TransferDiscrepancy    float64 `bson:"transfer_discrepancy" json:"transfer_discrepancy"`
    
    // DEPRECATED: Keep for backward compatibility
    DeclaredAmount float64 `bson:"declared_amount" json:"declared_amount"`
    ActualAmount   float64 `bson:"actual_amount" json:"actual_amount"`
    Discrepancy    float64 `bson:"discrepancy" json:"discrepancy"`
    
    // Existing fields...
    Status       HandoverStatus
    HandoverType HandoverType
    // ... other fields
}

// NEW: Calculate total amounts
func (h *CashHandover) TotalDeclaredAmount() float64 {
    return h.CashDeclaredAmount + h.TransferDeclaredAmount
}

func (h *CashHandover) TotalActualAmount() float64 {
    return h.CashActualAmount + h.TransferActualAmount
}

func (h *CashHandover) TotalDiscrepancy() float64 {
    return h.CashDiscrepancy + h.TransferDiscrepancy
}

// NEW: Check if handover includes transfer
func (h *CashHandover) HasTransfer() bool {
    return h.TransferDeclaredAmount > 0
}
```

**Shift (Extended)**
```go
type Shift struct {
    // Existing fields...
    ID            primitive.ObjectID
    UserID        primitive.ObjectID
    
    // Existing cash fields
    StartCash        float64
    CurrentCash      float64
    HandedOverCash   float64
    RemainingCash    float64
    
    // NEW: Transfer tracking fields
    TransferRevenue     float64 `bson:"transfer_revenue" json:"transfer_revenue"`
    HandedOverTransfer  float64 `bson:"handed_over_transfer" json:"handed_over_transfer"`
    RemainingTransfer   float64 `bson:"remaining_transfer" json:"remaining_transfer"`
    
    // Existing fields...
    TotalRevenue     float64
    TotalDiscrepancy float64
    // ... other fields
}

// NEW: Calculate total revenue from all payment methods
func (s *Shift) CalculateTotalRevenue(orders []*Order) {
    cashRevenue := 0.0
    transferRevenue := 0.0
    
    for _, order := range orders {
        if order.Status == StatusPaid {
            if order.PaymentMethod == PaymentMethodCash {
                cashRevenue += order.Total
            } else if order.PaymentMethod == PaymentMethodTransfer || 
                      order.PaymentMethod == PaymentMethodQR {
                transferRevenue += order.Total
            }
        }
    }
    
    s.CurrentCash = s.StartCash + cashRevenue
    s.RemainingCash = s.CurrentCash - s.HandedOverCash
    s.TransferRevenue = transferRevenue
    s.RemainingTransfer = s.TransferRevenue - s.HandedOverTransfer
    s.TotalRevenue = cashRevenue + transferRevenue
}
```

#### 2. Request/Response Structures

**CreateHandoverWithTransferRequest**
```go
type CreateHandoverWithTransferRequest struct {
    CashAmount     float64      `json:"cash_amount" binding:"gte=0"`
    TransferAmount float64      `json:"transfer_amount" binding:"gte=0"`
    HandoverType   HandoverType `json:"handover_type" binding:"required"`
    WaiterNote     string       `json:"waiter_note"`
}

// Validation
func (r *CreateHandoverWithTransferRequest) Validate() error {
    if r.CashAmount == 0 && r.TransferAmount == 0 {
        return errors.New("at least one amount must be greater than 0")
    }
    return nil
}
```

**ConfirmHandoverWithDualAmountsRequest**
```go
type ConfirmHandoverWithDualAmountsRequest struct {
    ActualCashAmount     float64            `json:"actual_cash_amount" binding:"gte=0"`
    ActualTransferAmount float64            `json:"actual_transfer_amount" binding:"gte=0"`
    Status               HandoverStatus     `json:"status" binding:"required"`
    CashierNote          string             `json:"cashier_note"`
    DiscrepancyReason    string             `json:"discrepancy_reason"`
    DiscrepancyResponsibility ResponsibilityType `json:"discrepancy_responsibility"`
}
```

#### 3. Service Methods

**CashHandoverService (Extended)**
```go
// CreateHandoverWithTransfer creates a handover with both cash and transfer
func (s *CashHandoverService) CreateHandoverWithTransfer(
    ctx context.Context,
    waiterShiftID primitive.ObjectID,
    req *CreateHandoverWithTransferRequest,
    waiterID, waiterName string,
) (*CashHandover, error) {
    // 1. Validate shift
    // 2. Validate amounts against shift balances
    // 3. Find active cashier shift
    // 4. Create handover record with separate amounts
    // 5. Update shift balances (remaining_cash, remaining_transfer)
}

// ConfirmHandoverWithDualAmounts confirms handover with separate actual amounts
func (s *CashHandoverService) ConfirmHandoverWithDualAmounts(
    ctx context.Context,
    handoverID primitive.ObjectID,
    req *ConfirmHandoverWithDualAmountsRequest,
    cashierID string,
) error {
    // 1. Get handover record
    // 2. Validate cashier authorization
    // 3. Calculate separate discrepancies
    // 4. Check if requires manager approval
    // 5. Update handover record
    // 6. If approved, update shift balances atomically
}

// UpdateDualBalances updates both cash and transfer balances
func (s *CashHandoverService) UpdateDualBalances(
    ctx context.Context,
    h *CashHandover,
) error {
    // 1. Update waiter shift:
    //    - handed_over_cash += actual_cash_amount
    //    - handed_over_transfer += actual_transfer_amount
    // 2. Update cashier shift:
    //    - received_cash += actual_cash_amount
    //    - received_transfer += actual_transfer_amount
    // 3. Handle END_SHIFT type (close shift)
}
```

### Frontend Components

#### 1. ShiftView (Waiter) - Handover Form

**Template Structure**
```vue
<template>
  <div class="handover-form">
    <!-- Shift Summary -->
    <div class="shift-summary">
      <div class="amount-card cash">
        <label>💵 Tiền mặt còn lại</label>
        <div class="amount">{{ formatPrice(shift.remaining_cash) }}</div>
      </div>
      <div class="amount-card transfer">
        <label>💳 Tiền CK đã thu</label>
        <div class="amount">{{ formatPrice(shift.remaining_transfer) }}</div>
      </div>
    </div>

    <!-- Handover Type Selection -->
    <div class="handover-type">
      <label>Chọn loại bàn giao</label>
      <div class="type-buttons">
        <button @click="selectType('cash')" 
          :class="{ active: handoverType === 'cash' }">
          💵 Chỉ tiền mặt
        </button>
        <button @click="selectType('transfer')" 
          :class="{ active: handoverType === 'transfer' }">
          💳 Chỉ tiền CK
        </button>
        <button @click="selectType('both')" 
          :class="{ active: handoverType === 'both' }">
          💰 Cả hai
        </button>
      </div>
    </div>

    <!-- Amount Inputs -->
    <div v-if="handoverType === 'cash' || handoverType === 'both'" 
      class="amount-input">
      <label>Số tiền mặt bàn giao (VNĐ)</label>
      <input v-model.number="cashAmount" type="number" 
        :max="shift.remaining_cash" />
    </div>

    <div v-if="handoverType === 'transfer' || handoverType === 'both'" 
      class="amount-input">
      <label>Số tiền CK bàn giao (VNĐ)</label>
      <input v-model.number="transferAmount" type="number" 
        :max="shift.remaining_transfer" />
    </div>

    <!-- Note -->
    <div class="note-input">
      <label>Ghi chú (tùy chọn)</label>
      <textarea v-model="waiterNote"></textarea>
    </div>

    <!-- Submit Button -->
    <button @click="submitHandover" class="submit-btn">
      Bàn giao
    </button>
  </div>
</template>
```

**Component Logic**
```javascript
const handoverType = ref('both') // 'cash' | 'transfer' | 'both'
const cashAmount = ref(0)
const transferAmount = ref(0)
const waiterNote = ref('')

const submitHandover = async () => {
  // Validate amounts
  if (handoverType.value === 'cash' || handoverType.value === 'both') {
    if (cashAmount.value > shift.value.remaining_cash) {
      alert('Số tiền mặt vượt quá số tiền còn lại')
      return
    }
  }
  
  if (handoverType.value === 'transfer' || handoverType.value === 'both') {
    if (transferAmount.value > shift.value.remaining_transfer) {
      alert('Số tiền CK vượt quá số tiền đã thu')
      return
    }
  }
  
  // Create handover
  await shiftStore.createHandoverWithTransfer({
    cash_amount: handoverType.value === 'transfer' ? 0 : cashAmount.value,
    transfer_amount: handoverType.value === 'cash' ? 0 : transferAmount.value,
    handover_type: 'PARTIAL',
    waiter_note: waiterNote.value
  })
}
```

#### 2. CashierShiftView - Handover Confirmation

**Template Structure**
```vue
<template>
  <div class="handover-confirmation">
    <!-- Declared Amounts -->
    <div class="declared-section">
      <h3>Số tiền khai báo</h3>
      <div class="amount-grid">
        <div v-if="handover.cash_declared_amount > 0" class="amount-card cash">
          <label>💵 Tiền mặt</label>
          <div class="amount">{{ formatPrice(handover.cash_declared_amount) }}</div>
        </div>
        <div v-if="handover.transfer_declared_amount > 0" class="amount-card transfer">
          <label>💳 Tiền CK</label>
          <div class="amount">{{ formatPrice(handover.transfer_declared_amount) }}</div>
        </div>
      </div>
    </div>

    <!-- Actual Amounts Input -->
    <div class="actual-section">
      <h3>Số tiền thực tế (sau đối soát)</h3>
      
      <div v-if="handover.cash_declared_amount > 0" class="amount-input">
        <label>💵 Tiền mặt thực nhận (VNĐ)</label>
        <input v-model.number="actualCashAmount" type="number" />
        <p class="hint">Đếm tiền mặt thực tế</p>
      </div>

      <div v-if="handover.transfer_declared_amount > 0" class="amount-input">
        <label>💳 Tiền CK thực nhận (VNĐ)</label>
        <input v-model.number="actualTransferAmount" type="number" />
        <p class="hint">Kiểm tra tài khoản ngân hàng</p>
      </div>
    </div>

    <!-- Discrepancy Display -->
    <div v-if="hasDiscrepancy" class="discrepancy-section">
      <h3>⚠️ Chênh lệch</h3>
      <div class="discrepancy-grid">
        <div v-if="cashDiscrepancy !== 0" class="discrepancy-item">
          <label>Tiền mặt</label>
          <div :class="cashDiscrepancy > 0 ? 'overage' : 'shortage'">
            {{ formatPrice(Math.abs(cashDiscrepancy)) }}
            {{ cashDiscrepancy > 0 ? '(Thừa)' : '(Thiếu)' }}
          </div>
        </div>
        <div v-if="transferDiscrepancy !== 0" class="discrepancy-item">
          <label>Tiền CK</label>
          <div :class="transferDiscrepancy > 0 ? 'overage' : 'shortage'">
            {{ formatPrice(Math.abs(transferDiscrepancy)) }}
            {{ transferDiscrepancy > 0 ? '(Thừa)' : '(Thiếu)' }}
          </div>
        </div>
      </div>
      
      <!-- Discrepancy Reason -->
      <div class="reason-input">
        <label>Lý do chênh lệch</label>
        <textarea v-model="discrepancyReason" required></textarea>
      </div>
    </div>

    <!-- Cashier Note -->
    <div class="note-input">
      <label>Ghi chú của thu ngân</label>
      <textarea v-model="cashierNote"></textarea>
    </div>

    <!-- Action Buttons -->
    <div class="action-buttons">
      <button @click="confirmHandover" class="confirm-btn">
        ✓ Xác nhận
      </button>
      <button @click="rejectHandover" class="reject-btn">
        ✗ Từ chối
      </button>
    </div>
  </div>
</template>
```

**Component Logic**
```javascript
const actualCashAmount = ref(0)
const actualTransferAmount = ref(0)
const cashierNote = ref('')
const discrepancyReason = ref('')

const cashDiscrepancy = computed(() => {
  return actualCashAmount.value - handover.value.cash_declared_amount
})

const transferDiscrepancy = computed(() => {
  return actualTransferAmount.value - handover.value.transfer_declared_amount
})

const totalDiscrepancy = computed(() => {
  return cashDiscrepancy.value + transferDiscrepancy.value
})

const hasDiscrepancy = computed(() => {
  return totalDiscrepancy.value !== 0
})

const confirmHandover = async () => {
  if (hasDiscrepancy.value && !discrepancyReason.value) {
    alert('Vui lòng nhập lý do chênh lệch')
    return
  }
  
  await cashierStore.confirmHandoverWithDualAmounts(handover.value.id, {
    actual_cash_amount: actualCashAmount.value,
    actual_transfer_amount: actualTransferAmount.value,
    status: 'CONFIRMED',
    cashier_note: cashierNote.value,
    discrepancy_reason: discrepancyReason.value
  })
}
```

## Data Models

### Database Schema Changes

#### 1. CashHandover Collection (Extended)

```javascript
{
  _id: ObjectId,
  waiter_shift_id: ObjectId,
  cashier_shift_id: ObjectId,
  waiter_id: ObjectId,
  waiter_name: String,
  cashier_id: ObjectId,
  cashier_name: String,
  
  // NEW: Separate cash amounts
  cash_declared_amount: Number,      // Tiền mặt khai báo
  cash_actual_amount: Number,        // Tiền mặt thực nhận
  cash_discrepancy: Number,          // Chênh lệch tiền mặt
  
  // NEW: Separate transfer amounts
  transfer_declared_amount: Number,  // Tiền CK khai báo
  transfer_actual_amount: Number,    // Tiền CK thực nhận
  transfer_discrepancy: Number,      // Chênh lệch tiền CK
  
  // DEPRECATED: Keep for backward compatibility
  declared_amount: Number,           // = cash + transfer (for old records)
  actual_amount: Number,             // = cash + transfer (for old records)
  discrepancy: Number,               // = cash + transfer (for old records)
  
  handover_type: String,             // PARTIAL | FULL | END_SHIFT
  status: String,                    // PENDING | CONFIRMED | REJECTED | DISCREPANCY
  
  waiter_note: String,
  cashier_note: String,
  discrepancy_reason: String,
  discrepancy_responsibility: String,
  
  handover_at: Date,
  confirmed_at: Date,
  reconciled_at: Date,
  
  end_cash: Number,
  requires_approval: Boolean,
  approved_by: ObjectId,
  approved_at: Date,
  
  created_at: Date,
  updated_at: Date
}
```

**Indexes**
```javascript
// Existing indexes
{ waiter_shift_id: 1, status: 1 }
{ cashier_shift_id: 1, status: 1 }
{ status: 1, requires_approval: 1 }

// NEW: Index for transfer handovers
{ transfer_declared_amount: 1 }  // Find handovers with transfer
```

#### 2. Shift Collection (Extended)

```javascript
{
  _id: ObjectId,
  type: String,                    // MORNING | AFTERNOON | EVENING
  status: String,                  // OPEN | CLOSED
  role_type: String,               // waiter | barista
  user_id: ObjectId,
  user_name: String,
  
  // Existing cash fields
  start_cash: Number,
  end_cash: Number,
  current_cash: Number,
  handed_over_cash: Number,
  remaining_cash: Number,
  
  // NEW: Transfer tracking fields
  transfer_revenue: Number,        // Tổng tiền CK thu được
  handed_over_transfer: Number,    // Tiền CK đã bàn giao
  remaining_transfer: Number,      // Tiền CK còn lại
  
  // Existing summary fields
  total_revenue: Number,           // = cash_revenue + transfer_revenue
  total_orders: Number,
  total_discrepancy: Number,       // = cash_discrepancy + transfer_discrepancy
  handover_count: Number,
  
  started_at: Date,
  ended_at: Date,
  created_at: Date,
  updated_at: Date
}
```

**Indexes**
```javascript
// Existing indexes
{ user_id: 1, status: 1 }
{ status: 1, role_type: 1 }
{ started_at: -1 }

// NEW: Index for transfer tracking
{ transfer_revenue: 1 }
```

#### 3. CashierShift Collection (Extended)

```javascript
{
  _id: ObjectId,
  cashier_id: ObjectId,
  cashier_name: String,
  status: String,                  // OPEN | CLOSED
  
  // Existing cash fields
  start_cash: Number,
  end_cash: Number,
  received_cash: Number,
  
  // NEW: Transfer tracking fields
  received_transfer: Number,       // Tổng tiền CK nhận được
  
  // Existing summary fields
  total_discrepancy: Number,
  handover_count: Number,
  discrepancy_count: Number,
  
  started_at: Date,
  ended_at: Date,
  created_at: Date,
  updated_at: Date
}
```

### Migration Strategy

**Phase 1: Add New Fields**
```javascript
// Add transfer fields to existing documents with default values
db.shifts.updateMany(
  { transfer_revenue: { $exists: false } },
  {
    $set: {
      transfer_revenue: 0,
      handed_over_transfer: 0,
      remaining_transfer: 0
    }
  }
)

db.cash_handovers.updateMany(
  { cash_declared_amount: { $exists: false } },
  {
    $set: {
      cash_declared_amount: "$declared_amount",
      cash_actual_amount: "$actual_amount",
      cash_discrepancy: "$discrepancy",
      transfer_declared_amount: 0,
      transfer_actual_amount: 0,
      transfer_discrepancy: 0
    }
  }
)

db.cashier_shifts.updateMany(
  { received_transfer: { $exists: false } },
  {
    $set: {
      received_transfer: 0
    }
  }
)
```

**Phase 2: Recalculate Transfer Revenue**
```javascript
// For each open shift, recalculate transfer_revenue from orders
const openShifts = db.shifts.find({ status: "OPEN" })
openShifts.forEach(shift => {
  const orders = db.orders.find({ 
    shift_id: shift._id,
    status: "PAID",
    payment_method: { $in: ["TRANSFER", "QR"] }
  })
  
  const transferRevenue = orders.reduce((sum, order) => sum + order.total, 0)
  
  db.shifts.updateOne(
    { _id: shift._id },
    {
      $set: {
        transfer_revenue: transferRevenue,
        remaining_transfer: transferRevenue
      }
    }
  )
})
```



## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property Reflection

After analyzing all acceptance criteria, I identified several opportunities to consolidate redundant properties:

**Consolidations Made:**
- Properties 2.1, 2.4, and 6.4 (separate tracking) → Combined into Property 1 (comprehensive separation)
- Properties 2.2 and 1.2 (transfer revenue calculation) → Combined into Property 2
- Properties 2.5 and 1.4 (separate field storage) → Combined into Property 3
- Properties 3.3 and 4.3 (discrepancy calculation) → Combined into Property 4
- Properties 4.1, 4.4, and 4.5 (atomicity) → Combined into Property 5
- Properties 5.1, 5.2, 5.3, 5.4, 5.5 (display formatting) → Combined into Property 6
- Properties 7.1, 7.2, 7.4 (backward compatibility) → Combined into Property 7

### Core Properties

**Property 1: Separate Cash and Transfer Tracking**

*For any* shift with orders using different payment methods, the system SHALL maintain separate tracking for cash_revenue, transfer_revenue, remaining_cash, and remaining_transfer, where cash_revenue equals the sum of all CASH payments, transfer_revenue equals the sum of all TRANSFER and QR payments, and remaining amounts are calculated independently.

**Validates: Requirements 2.1, 2.4, 6.4**

---

**Property 2: Transfer Revenue Calculation**

*For any* shift with paid orders, the calculated transfer_revenue SHALL equal the sum of all order totals where payment_method is TRANSFER or QR.

**Validates: Requirements 2.2, 1.2**

---

**Property 3: Handover Record Structure**

*For any* handover created with transfer amount greater than zero, the handover record SHALL contain separate non-null fields for cash_declared_amount, transfer_declared_amount, cash_actual_amount, transfer_actual_amount, cash_discrepancy, and transfer_discrepancy.

**Validates: Requirements 2.5, 1.4**

---

**Property 4: Discrepancy Calculation**

*For any* handover confirmation with actual amounts, cash_discrepancy SHALL equal (cash_actual_amount - cash_declared_amount), transfer_discrepancy SHALL equal (transfer_actual_amount - transfer_declared_amount), and total_discrepancy SHALL equal (cash_discrepancy + transfer_discrepancy).

**Validates: Requirements 3.3, 4.3**

---

**Property 5: Atomic Confirmation**

*For any* handover confirmation operation, both cash and transfer balance updates SHALL complete successfully together, or if either update fails, the entire transaction SHALL rollback leaving all balances unchanged.

**Validates: Requirements 4.1, 4.4, 4.5**

---

**Property 6: Display Formatting Consistency**

*For any* handover displayed in the UI, if cash_declared_amount is greater than zero, it SHALL be shown with green color coding, and if transfer_declared_amount is greater than zero, it SHALL be shown with blue color coding, and both amounts SHALL be visible as separate line items.

**Validates: Requirements 5.1, 5.2, 5.3, 5.4, 5.5**

---

**Property 7: Backward Compatibility**

*For any* handover where transfer_declared_amount is zero or null, the system SHALL process it using cash-only logic, requiring only cash_actual_amount for confirmation and updating only cash balances.

**Validates: Requirements 7.1, 7.2, 7.4**

---

### Validation Properties

**Property 8: Declared Amount Validation**

*For any* handover creation request, if cash_amount is greater than shift.remaining_cash OR transfer_amount is greater than shift.remaining_transfer, the system SHALL reject the request with a validation error.

**Validates: Requirements 1.3, 6.1**

---

**Property 9: Non-Negative Transfer Amounts**

*For any* handover creation or confirmation, if transfer_declared_amount or transfer_actual_amount is less than zero, the system SHALL reject the operation with a validation error.

**Validates: Requirements 6.3**

---

**Property 10: Single Pending Handover**

*For any* shift with a handover in PENDING status, attempts to create a new handover for the same shift SHALL be rejected until the pending handover is confirmed or rejected.

**Validates: Requirements 6.2**

---

### Business Logic Properties

**Property 11: End-Shift Handover Amounts**

*For any* end-shift handover creation, the system SHALL automatically set cash_declared_amount equal to shift.remaining_cash AND transfer_declared_amount equal to shift.remaining_transfer.

**Validates: Requirements 1.5**

---

**Property 12: Manager Approval Threshold**

*For any* handover confirmation where total_discrepancy exceeds 100,000 VND OR transfer_discrepancy exceeds 50,000 VND, the system SHALL set requires_approval to true and status to DISCREPANCY.

**Validates: Requirements 8.1, 8.2**

---

**Property 13: Balance Update After Approval**

*For any* handover with requires_approval true, when a manager approves it, the system SHALL update waiter_shift.handed_over_cash by cash_actual_amount, waiter_shift.handed_over_transfer by transfer_actual_amount, cashier_shift.received_cash by cash_actual_amount, and cashier_shift.received_transfer by transfer_actual_amount, and set status to CONFIRMED.

**Validates: Requirements 8.4**

---

**Property 14: Rejection State Transition**

*For any* handover with requires_approval true, when a manager rejects it, the system SHALL set status to PENDING and leave all balance fields unchanged.

**Validates: Requirements 8.5**

---

### Data Integrity Properties

**Property 15: Audit Logging**

*For any* handover transaction (create, confirm, approve, reject), the system SHALL create a log entry containing handover_id, user_id, action_type, timestamp, and before/after state.

**Validates: Requirements 6.5**

---

**Property 16: Migration Data Integrity**

*For any* existing handover record without transfer fields, when accessed by the system, transfer_declared_amount, transfer_actual_amount, and transfer_discrepancy SHALL default to zero, and the record SHALL be processable without errors.

**Validates: Requirements 7.3, 7.5**

---

## Error Handling

### Validation Errors

**Invalid Amount Errors**
- **Scenario**: Waiter declares amount exceeding available balance
- **Response**: Return 400 Bad Request with message "Declared amount exceeds available balance"
- **Recovery**: User must enter valid amount within limits

**Negative Amount Errors**
- **Scenario**: User enters negative transfer amount
- **Response**: Return 400 Bad Request with message "Transfer amount must be non-negative"
- **Recovery**: User must enter zero or positive amount

**Pending Handover Conflict**
- **Scenario**: Waiter tries to create handover while one is pending
- **Response**: Return 409 Conflict with message "A pending handover already exists for this shift"
- **Recovery**: User must wait for cashier to confirm or cancel existing handover

### Authorization Errors

**Unauthorized Waiter**
- **Scenario**: Waiter tries to create handover for another waiter's shift
- **Response**: Return 403 Forbidden with message "You can only create handovers for your own shift"
- **Recovery**: User must use correct shift

**Unauthorized Cashier**
- **Scenario**: Cashier tries to confirm handover assigned to different cashier
- **Response**: Return 403 Forbidden with message "This handover is not assigned to you"
- **Recovery**: Correct cashier must confirm

### Business Logic Errors

**Shift Not Open**
- **Scenario**: Waiter tries to create handover for closed shift
- **Response**: Return 400 Bad Request with message "Cannot create handover for closed shift"
- **Recovery**: Shift must be reopened or use different shift

**No Active Cashier Shift**
- **Scenario**: Waiter creates handover but no cashier shift is open
- **Response**: Return 503 Service Unavailable with message "No active cashier shift found"
- **Recovery**: Cashier must start a shift first

**Discrepancy Without Reason**
- **Scenario**: Cashier confirms with discrepancy but no reason provided
- **Response**: Return 400 Bad Request with message "Discrepancy reason is required when amounts differ"
- **Recovery**: Cashier must provide reason

### System Errors

**Database Transaction Failure**
- **Scenario**: Balance update fails mid-transaction
- **Response**: Rollback all changes, return 500 Internal Server Error
- **Recovery**: Retry operation, check database connectivity

**Concurrent Modification**
- **Scenario**: Two cashiers try to confirm same handover simultaneously
- **Response**: Return 409 Conflict with message "Handover has been modified by another user"
- **Recovery**: Refresh and retry

### Error Recovery Strategies

**Automatic Retry**
- Transient network errors: Retry up to 3 times with exponential backoff
- Database connection errors: Retry with connection pool refresh

**Manual Intervention**
- Large discrepancies: Require manager approval
- Data inconsistencies: Log error and notify admin

**Graceful Degradation**
- If notification service fails: Log notification for later retry
- If audit logging fails: Complete transaction but alert admin

## Testing Strategy

### Dual Testing Approach

This feature requires both unit tests and property-based tests for comprehensive coverage:

**Unit Tests** focus on:
- Specific examples of handover scenarios
- Edge cases (zero amounts, exact threshold values)
- Error conditions and validation
- Integration between components

**Property-Based Tests** focus on:
- Universal properties across all valid inputs
- Randomized test data generation
- Comprehensive input coverage
- Invariant verification

Both testing approaches are complementary and necessary. Unit tests catch concrete bugs in specific scenarios, while property tests verify general correctness across the input space.

### Property-Based Testing Configuration

**Testing Library**: Use `gopter` for Go backend property-based testing

**Test Configuration**:
- Minimum 100 iterations per property test
- Each test must reference its design document property
- Tag format: `// Feature: bank-transfer-handover, Property N: [property title]`

**Example Property Test Structure**:
```go
func TestProperty1_SeparateCashAndTransferTracking(t *testing.T) {
    // Feature: bank-transfer-handover, Property 1: Separate Cash and Transfer Tracking
    properties := gopter.NewProperties(nil)
    
    properties.Property("cash and transfer tracked separately", 
        prop.ForAll(
            func(orders []Order) bool {
                shift := calculateShiftRevenue(orders)
                
                cashOrders := filterByPaymentMethod(orders, CASH)
                transferOrders := filterByPaymentMethod(orders, TRANSFER, QR)
                
                expectedCash := sumOrderTotals(cashOrders)
                expectedTransfer := sumOrderTotals(transferOrders)
                
                return shift.CashRevenue == expectedCash &&
                       shift.TransferRevenue == expectedTransfer &&
                       shift.TotalRevenue == expectedCash + expectedTransfer
            },
            genOrders(),
        ))
    
    properties.TestingRun(t, gopter.ConsoleReporter(false))
}
```

### Unit Test Coverage

**Handover Creation Tests**:
- Test creating cash-only handover
- Test creating transfer-only handover
- Test creating combined handover
- Test end-shift handover with both amounts
- Test validation errors (exceeding balance, negative amounts)
- Test pending handover conflict

**Handover Confirmation Tests**:
- Test confirming with exact amounts (no discrepancy)
- Test confirming with cash discrepancy
- Test confirming with transfer discrepancy
- Test confirming with both discrepancies
- Test manager approval trigger
- Test rejection flow

**Balance Update Tests**:
- Test shift balance updates after confirmation
- Test cashier shift balance updates
- Test atomic transaction rollback on failure
- Test end-shift closure with balance updates

**Backward Compatibility Tests**:
- Test cash-only handover without transfer fields
- Test displaying old handover records
- Test migration of existing data

### Integration Test Scenarios

**Complete Handover Flow**:
1. Waiter creates shift with start cash
2. Orders are paid with mixed payment methods
3. Waiter creates handover with both cash and transfer
4. Cashier confirms with actual amounts
5. Verify all balances updated correctly

**Discrepancy Approval Flow**:
1. Waiter creates handover
2. Cashier confirms with large discrepancy
3. System flags for manager approval
4. Manager approves
5. Verify balances updated after approval

**End-Shift Flow**:
1. Waiter works full shift with mixed payments
2. Waiter creates end-shift handover
3. System auto-includes all remaining amounts
4. Cashier confirms
5. Shift closes with correct final balances

### Test Data Generators

**Order Generator**:
```go
func genOrders() gopter.Gen {
    return gen.SliceOf(gen.Struct(reflect.TypeOf(Order{}), map[string]gopter.Gen{
        "Total": gen.Float64Range(1000, 1000000),
        "PaymentMethod": gen.OneConstOf(CASH, TRANSFER, QR),
        "Status": gen.OneConstOf(StatusPaid),
    }))
}
```

**Handover Generator**:
```go
func genHandover() gopter.Gen {
    return gen.Struct(reflect.TypeOf(CreateHandoverRequest{}), map[string]gopter.Gen{
        "CashAmount": gen.Float64Range(0, 10000000),
        "TransferAmount": gen.Float64Range(0, 10000000),
        "HandoverType": gen.OneConstOf(TypePartial, TypeFull, TypeEndShift),
    })
}
```

### Performance Testing

**Load Testing**:
- Test 100 concurrent handover creations
- Test 50 concurrent confirmations
- Verify no race conditions or deadlocks

**Stress Testing**:
- Test with shifts containing 1000+ orders
- Test with 100+ handovers per shift
- Verify acceptable response times (<500ms)

### Monitoring and Observability

**Metrics to Track**:
- Handover creation rate
- Confirmation latency
- Discrepancy rate (cash vs transfer)
- Manager approval rate
- Transaction rollback rate

**Alerts**:
- High discrepancy rate (>10%)
- Frequent transaction failures
- Slow confirmation times (>1s)
- Pending handovers not confirmed within 1 hour

