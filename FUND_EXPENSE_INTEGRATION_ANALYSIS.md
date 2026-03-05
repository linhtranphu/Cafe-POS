# Phân tích: Tích hợp Fund với Expense/Ingredient/Facility

## Yêu cầu

Khi tạo expense, mua ingredient, hoặc sửa chữa facility, hệ thống cần:
1. Ghi nhận là tiền được chi từ quỹ (fund)
2. Tự động tạo fund withdrawal transaction
3. Cập nhật số dư quỹ
4. Liên kết expense với fund transaction để audit

## Cấu trúc hiện tại

### 1. Fund Transaction
```go
type FundTransaction struct {
    ID              primitive.ObjectID
    Type            TransactionType     // "deposit" | "withdrawal"
    CashAmount      float64
    TransferAmount  float64
    TotalAmount     float64
    Reason          string
    PerformedBy     primitive.ObjectID
    PerformedByName string
    PerformedByRole string
    Timestamp       time.Time
    BalanceBefore   *FundBalance
    BalanceAfter    *FundBalance
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### 2. Expense
```go
type Expense struct {
    ID            primitive.ObjectID
    Date          time.Time
    CategoryID    primitive.ObjectID
    Amount        float64
    Description   string
    PaymentMethod string              // "cash" | "bank" | "card"
    Vendor        string
    Notes         string
    
    // Source tracking
    SourceType    string              // "ingredient" | "facility" | "manual"
    SourceID      primitive.ObjectID
    
    CreatedBy     string
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

### 3. Ingredient
```go
type Ingredient struct {
    ID                primitive.ObjectID
    Name              string
    Category          string
    Unit              string
    CurrentStock      float64
    MinimumStock      float64
    UnitCost          float64
    Supplier          string
    LastRestockDate   time.Time
    LastRestockAmount float64
    CreatedAt         time.Time
    UpdatedAt         time.Time
}
```

### 4. Facility
```go
type Facility struct {
    ID          primitive.ObjectID
    Name        string
    Type        string
    Area        string
    Status      string
    Description string
    // ... other fields
}

type FacilityHistory struct {
    ID          primitive.ObjectID
    FacilityID  primitive.ObjectID
    Action      string
    Description string
    Cost        float64
    PerformedBy string
    Timestamp   time.Time
}
```

## Thiết kế giải pháp

### Option 1: Thêm FundTransactionID vào Expense (Khuyến nghị)

**Ưu điểm:**
- Đơn giản, rõ ràng
- Dễ audit: Từ expense → fund transaction
- Không thay đổi nhiều code hiện tại

**Cấu trúc mới:**
```go
type Expense struct {
    // ... existing fields ...
    
    // Fund integration - NEW
    FundTransactionID *primitive.ObjectID `bson:"fund_transaction_id,omitempty" json:"fund_transaction_id,omitempty"`
    PaidFromFund      bool                `bson:"paid_from_fund" json:"paid_from_fund"` // true if paid from fund
}
```

**Luồng xử lý:**
```
1. User tạo expense (amount: 500,000đ, payment_method: "cash")
2. Backend:
   a. Tạo Expense record
   b. Nếu paid_from_fund = true:
      - Tạo FundTransaction (type: "withdrawal", amount: 500,000)
      - Cập nhật fund balance
      - Lưu fund_transaction_id vào expense
   c. Return expense với fund_transaction_id
```

### Option 2: Thêm SourceType vào FundTransaction

**Ưu điểm:**
- Audit từ fund transaction → expense
- Biết được fund được dùng cho mục đích gì

**Cấu trúc mới:**
```go
type FundTransaction struct {
    // ... existing fields ...
    
    // Source tracking - NEW
    SourceType string             `bson:"source_type,omitempty" json:"source_type,omitempty"` // "expense" | "ingredient" | "facility"
    SourceID   primitive.ObjectID `bson:"source_id,omitempty" json:"source_id,omitempty"`
}
```

### Option 3: Kết hợp cả 2 (Tốt nhất)

Thêm cả 2 chiều liên kết để audit dễ dàng:
- Expense → FundTransaction (fund_transaction_id)
- FundTransaction → Expense (source_type, source_id)

## Implementation Plan

### Phase 1: Backend Changes

#### 1.1. Update Domain Models

**File: `backend/domain/expense/expense.go`**
```go
type Expense struct {
    // ... existing fields ...
    
    // Fund integration
    FundTransactionID *primitive.ObjectID `bson:"fund_transaction_id,omitempty" json:"fund_transaction_id,omitempty"`
    PaidFromFund      bool                `bson:"paid_from_fund" json:"paid_from_fund"`
}
```

**File: `backend/domain/fund/fund_transaction.go`**
```go
type FundTransaction struct {
    // ... existing fields ...
    
    // Source tracking
    SourceType string             `bson:"source_type,omitempty" json:"source_type,omitempty"`
    SourceID   primitive.ObjectID `bson:"source_id,omitempty" json:"source_id,omitempty"`
}

// Source type constants
const (
    SourceTypeExpense    = "expense"
    SourceTypeIngredient = "ingredient"
    SourceTypeFacility   = "facility"
)
```

#### 1.2. Update Services

**File: `backend/application/services/expense_service.go`**
```go
func (s *ExpenseService) CreateExpense(ctx context.Context, req *CreateExpenseRequest) (*Expense, error) {
    // 1. Create expense
    expense := &Expense{
        // ... set fields ...
        PaidFromFund: req.PaidFromFund,
    }
    
    // 2. If paid from fund, create fund transaction
    if req.PaidFromFund {
        fundTx, err := s.fundService.CreateWithdrawal(ctx, &FundWithdrawalRequest{
            CashAmount:     req.Amount, // or split between cash/transfer
            TransferAmount: 0,
            Reason:         fmt.Sprintf("Chi phí: %s", req.Description),
            SourceType:     fund.SourceTypeExpense,
            SourceID:       expense.ID,
        })
        if err != nil {
            return nil, fmt.Errorf("failed to create fund withdrawal: %w", err)
        }
        
        expense.FundTransactionID = &fundTx.ID
    }
    
    // 3. Save expense
    if err := s.expenseRepo.Create(ctx, expense); err != nil {
        // Rollback fund transaction if needed
        return nil, err
    }
    
    return expense, nil
}
```

**Similar for Ingredient and Facility:**
- When restocking ingredient → create fund withdrawal
- When repairing facility → create fund withdrawal

#### 1.3. Update API Handlers

**File: `backend/interfaces/http/expense_handler.go`**
```go
type CreateExpenseRequest struct {
    Date          string  `json:"date" binding:"required"`
    CategoryID    string  `json:"category_id" binding:"required"`
    Amount        float64 `json:"amount" binding:"required,gt=0"`
    Description   string  `json:"description" binding:"required"`
    PaymentMethod string  `json:"payment_method" binding:"required"`
    Vendor        string  `json:"vendor"`
    Notes         string  `json:"notes"`
    
    // Fund integration - NEW
    PaidFromFund  bool    `json:"paid_from_fund"` // Default false
}
```

### Phase 2: Frontend Changes

#### 2.1. Update Expense Form

**Add checkbox:**
```vue
<div class="form-group">
  <label>
    <input type="checkbox" v-model="form.paid_from_fund" />
    Chi từ quỹ
  </label>
  <p class="help-text">Nếu chọn, số tiền sẽ được trừ từ quỹ</p>
</div>
```

#### 2.2. Show Fund Transaction Link

**In expense detail:**
```vue
<div v-if="expense.fund_transaction_id" class="fund-link">
  <span>💰 Đã chi từ quỹ</span>
  <router-link :to="`/fund/transactions/${expense.fund_transaction_id}`">
    Xem giao dịch
  </router-link>
</div>
```

### Phase 3: Database Migration

**Create migration script:**
```javascript
// Add new fields to existing expenses
db.expenses.updateMany(
  {},
  {
    $set: {
      paid_from_fund: false,
      fund_transaction_id: null
    }
  }
)

// Add source tracking to fund transactions
db.fund_transactions.updateMany(
  {},
  {
    $set: {
      source_type: null,
      source_id: null
    }
  }
)
```

## Use Cases

### Use Case 1: Tạo expense thủ công

**Scenario:** Manager tạo expense cho tiền điện

**Steps:**
1. Manager vào "Expenses" → "Create New"
2. Nhập:
   - Category: Utilities
   - Amount: 500,000đ
   - Description: Tiền điện tháng 3
   - Payment method: Cash
   - ✅ Paid from fund: TRUE
3. Submit
4. System:
   - Tạo expense record
   - Tạo fund withdrawal (500,000đ)
   - Cập nhật fund balance: -500,000đ
   - Link expense ↔ fund transaction

### Use Case 2: Mua ingredient

**Scenario:** Manager mua thêm cà phê

**Steps:**
1. Manager vào "Ingredients" → Select "Cà phê" → "Restock"
2. Nhập:
   - Quantity: 10kg
   - Unit cost: 200,000đ/kg
   - Total: 2,000,000đ
   - ✅ Paid from fund: TRUE
3. Submit
4. System:
   - Cập nhật ingredient stock
   - Tạo expense (category: "Ingredient Purchase")
   - Tạo fund withdrawal (2,000,000đ)
   - Link expense ↔ fund transaction

### Use Case 3: Sửa chữa facility

**Scenario:** Sửa máy pha cà phê

**Steps:**
1. Manager vào "Facilities" → Select "Máy pha cà phê" → "Report Issue"
2. Nhập:
   - Issue: Hỏng motor
   - Cost: 1,500,000đ
   - ✅ Paid from fund: TRUE
3. Submit
4. System:
   - Tạo facility history record
   - Tạo expense (category: "Maintenance")
   - Tạo fund withdrawal (1,500,000đ)
   - Link expense ↔ fund transaction

## Reporting & Analytics

### Fund Usage Report

**Query:** Xem fund được dùng cho mục đích gì

```javascript
// Aggregate fund withdrawals by source type
db.fund_transactions.aggregate([
  { $match: { type: "withdrawal" } },
  { $group: {
      _id: "$source_type",
      total: { $sum: "$total_amount" },
      count: { $sum: 1 }
  }}
])

// Result:
// { _id: "expense", total: 5000000, count: 10 }
// { _id: "ingredient", total: 3000000, count: 5 }
// { _id: "facility", total: 2000000, count: 3 }
```

### Expense with Fund Tracking

**Query:** Xem expenses nào được chi từ quỹ

```javascript
db.expenses.find({ paid_from_fund: true })
```

## Validation Rules

### 1. Fund Balance Check
```go
// Before creating withdrawal, check if fund has enough balance
if fundBalance.Total < withdrawalAmount {
    return errors.New("insufficient fund balance")
}
```

### 2. Payment Method Consistency
```go
// If paid_from_fund = true, payment_method should be "cash" or "bank"
if req.PaidFromFund && req.PaymentMethod == "card" {
    return errors.New("card payment cannot be from fund")
}
```

### 3. Prevent Double Spending
```go
// Expense can only have one fund_transaction_id
if expense.FundTransactionID != nil {
    return errors.New("expense already linked to fund transaction")
}
```

## Security & Permissions

### Who can create fund withdrawals?
- **Manager**: Full access
- **Cashier**: Can create for specific categories
- **Waiter/Barista**: No access

### Audit Trail
- All fund transactions logged with:
  - Who performed (user_id, name, role)
  - When (timestamp)
  - Why (reason)
  - Source (expense_id, ingredient_id, facility_id)

## Testing Checklist

- [ ] Create expense with paid_from_fund = true
- [ ] Verify fund withdrawal created
- [ ] Verify fund balance updated
- [ ] Verify expense.fund_transaction_id set
- [ ] Verify fund_transaction.source_id set
- [ ] Create expense with paid_from_fund = false
- [ ] Verify no fund transaction created
- [ ] Restock ingredient with paid_from_fund = true
- [ ] Verify fund withdrawal for ingredient
- [ ] Repair facility with paid_from_fund = true
- [ ] Verify fund withdrawal for facility
- [ ] Check insufficient fund balance error
- [ ] Check fund usage report
- [ ] Check expense audit trail

## Migration Strategy

### Step 1: Add new fields (backward compatible)
- Add `paid_from_fund` (default false)
- Add `fund_transaction_id` (nullable)
- Add `source_type`, `source_id` to fund_transaction

### Step 2: Update services
- Implement fund withdrawal logic
- Add validation

### Step 3: Update UI
- Add checkbox to forms
- Show fund transaction links

### Step 4: Test thoroughly
- Test all scenarios
- Verify data integrity

### Step 5: Deploy
- Deploy backend first
- Deploy frontend
- Monitor for issues

## Future Enhancements

### 1. Automatic Fund Allocation
- Auto-allocate fund for recurring expenses
- Budget planning per category

### 2. Fund Transfer Between Categories
- Move fund from one category to another
- Track fund allocation history

### 3. Fund Approval Workflow
- Large withdrawals require approval
- Multi-level approval for different amounts

### 4. Fund Reconciliation
- Daily/weekly fund reconciliation
- Compare expected vs actual balance

---

**Status:** 📋 Analysis Complete - Ready for Implementation
**Priority:** High
**Estimated Effort:** 3-5 days
