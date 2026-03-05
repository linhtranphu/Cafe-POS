# Design Document: Fund-Expense Integration

## Overview

Tính năng Fund-Expense Integration tích hợp hệ thống quản lý quỹ (Fund Manager) với các module chi tiêu hiện có (Expense, Ingredient) trong hệ thống Cafe POS. Tính năng này cho phép Manager chi tiền trực tiếp từ quỹ khi tạo các giao dịch chi tiêu, đồng thời tự động ghi nhận fund transaction và cập nhật số dư quỹ.

### Goals

1. Cho phép Manager chi tiền từ quỹ khi tạo expense hoặc mua ingredient
2. Tự động tạo fund withdrawal transaction và cập nhật số dư quỹ
3. Duy trì tính toàn vẹn dữ liệu với bidirectional linking giữa fund transactions và source records
4. Đảm bảo transaction atomicity để tránh trạng thái không nhất quán
5. Cung cấp audit trail đầy đủ cho mọi giao dịch chi tiêu từ quỹ
6. Ngăn chặn double spending và chi vượt quỹ

### Non-Goals

- Requirement 3 (Facility maintenance integration) - optional feature, không implement trong phase này
- Requirement 5 (Fund usage reporting) - optional feature, không implement trong phase này
- Tự động reconciliation giữa fund balance và accounting records
- Multi-currency support
- Approval workflow cho fund withdrawals

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Frontend (Vue.js)                     │
│  ┌──────────────────┐  ┌──────────────────┐                │
│  │ Expense Form     │  │ Ingredient Form  │                │
│  │ + Paid from Fund │  │ + Paid from Fund │                │
│  └──────────────────┘  └──────────────────┘                │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    Backend API Layer (Go)                    │
│  ┌──────────────────┐  ┌──────────────────┐                │
│  │ Expense Handler  │  │ Ingredient       │                │
│  │                  │  │ Handler          │                │
│  └──────────────────┘  └──────────────────┘                │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   Application Services                       │
│  ┌──────────────────────────────────────────────────────┐  │
│  │         FundExpenseIntegrationService                 │  │
│  │  - CreateExpenseFromFund()                           │  │
│  │  - RestockIngredientFromFund()                       │  │
│  │  - ValidateFundBalance()                             │  │
│  │  - RollbackTransaction()                             │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
                ┌───────────┴───────────┐
                ▼                       ▼
┌──────────────────────┐    ┌──────────────────────┐
│  Expense Repository  │    │  Fund Repository     │
│  - Create()          │    │  - GetBalance()      │
│  - Update()          │    │  - CreateWithdrawal()│
│  - GetByID()         │    │  - UpdateBalance()   │
└──────────────────────┘    └──────────────────────┘
                │                       │
                └───────────┬───────────┘
                            ▼
                ┌──────────────────────┐
                │   MongoDB Database   │
                │  - expenses          │
                │  - fund_transactions │
                │  - ingredients       │
                └──────────────────────┘
```

### Component Responsibilities

#### FundExpenseIntegrationService
- Orchestrates the creation of expenses/ingredient restocks paid from fund
- Validates fund balance before processing
- Ensures transaction atomicity across multiple operations
- Handles rollback in case of failures
- Maintains bidirectional linking between records

#### Expense Module
- Manages expense records with fund integration fields
- Validates payment method consistency
- Prevents modification of fund-linked expenses
- Provides filtering for fund-paid expenses

#### Fund Manager
- Manages fund transactions and balance
- Creates withdrawal transactions for fund-paid operations
- Records audit trail information
- Validates source_type and prevents duplicate transactions

#### Ingredient Module
- Manages ingredient restock operations
- Integrates with fund manager for fund-paid restocks
- Updates stock quantities atomically with fund transactions

## Components and Interfaces

### Domain Models

#### Enhanced Expense Model
```go
type Expense struct {
    ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Date          time.Time          `bson:"date" json:"date"`
    CategoryID    primitive.ObjectID `bson:"category_id" json:"category_id"`
    Amount        float64            `bson:"amount" json:"amount"`
    Description   string             `bson:"description" json:"description"`
    PaymentMethod string             `bson:"payment_method" json:"payment_method"`
    Vendor        string             `bson:"vendor,omitempty" json:"vendor,omitempty"`
    Notes         string             `bson:"notes,omitempty" json:"notes,omitempty"`
    
    // Source tracking
    SourceType string             `bson:"source_type,omitempty" json:"source_type,omitempty"`
    SourceID   primitive.ObjectID `bson:"source_id,omitempty" json:"source_id,omitempty"`
    
    // NEW: Fund integration fields
    PaidFromFund      bool               `bson:"paid_from_fund" json:"paid_from_fund"`
    FundTransactionID primitive.ObjectID `bson:"fund_transaction_id,omitempty" json:"fund_transaction_id,omitempty"`
    
    // Audit fields
    CreatedBy string    `bson:"created_by" json:"created_by"`
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
```

#### Enhanced FundTransaction Model
```go
type FundTransaction struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Type            TransactionType    `bson:"type" json:"type"`
    CashAmount      float64            `bson:"cash_amount" json:"cash_amount"`
    TransferAmount  float64            `bson:"transfer_amount" json:"transfer_amount"`
    TotalAmount     float64            `bson:"total_amount" json:"total_amount"`
    Reason          string             `bson:"reason" json:"reason"`
    
    // NEW: Source tracking fields
    SourceType string             `bson:"source_type,omitempty" json:"source_type,omitempty"`
    SourceID   primitive.ObjectID `bson:"source_id,omitempty" json:"source_id,omitempty"`
    
    // Audit fields
    PerformedBy     primitive.ObjectID `bson:"performed_by" json:"performed_by"`
    PerformedByName string             `bson:"performed_by_name" json:"performed_by_name"`
    PerformedByRole string             `bson:"performed_by_role" json:"performed_by_role"`
    Timestamp       time.Time          `bson:"timestamp" json:"timestamp"`
    BalanceBefore   *FundBalance       `bson:"balance_before,omitempty" json:"balance_before,omitempty"`
    BalanceAfter    *FundBalance       `bson:"balance_after,omitempty" json:"balance_after,omitempty"`
    CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}
```

#### Ingredient Restock Record
```go
type IngredientRestockRecord struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    IngredientID    primitive.ObjectID `bson:"ingredient_id" json:"ingredient_id"`
    Quantity        float64            `bson:"quantity" json:"quantity"`
    CostPerUnit     float64            `bson:"cost_per_unit" json:"cost_per_unit"`
    TotalCost       float64            `bson:"total_cost" json:"total_cost"`
    
    // NEW: Fund integration fields
    PaidFromFund      bool               `bson:"paid_from_fund" json:"paid_from_fund"`
    ExpenseID         primitive.ObjectID `bson:"expense_id,omitempty" json:"expense_id,omitempty"`
    FundTransactionID primitive.ObjectID `bson:"fund_transaction_id,omitempty" json:"fund_transaction_id,omitempty"`
    
    // Audit fields
    PerformedBy     string    `bson:"performed_by" json:"performed_by"`
    PerformedByName string    `bson:"performed_by_name" json:"performed_by_name"`
    Reason          string    `bson:"reason" json:"reason"`
    CreatedAt       time.Time `bson:"created_at" json:"created_at"`
}
```

### Service Interfaces

#### FundExpenseIntegrationService Interface
```go
type FundExpenseIntegrationService interface {
    // Create expense paid from fund with atomicity guarantee
    CreateExpenseFromFund(ctx context.Context, req CreateExpenseFromFundRequest) (*ExpenseFromFundResult, error)
    
    // Restock ingredient paid from fund with atomicity guarantee
    RestockIngredientFromFund(ctx context.Context, req RestockFromFundRequest) (*RestockFromFundResult, error)
    
    // Validate fund balance before operation
    ValidateFundBalance(ctx context.Context, requiredAmount float64) error
    
    // Get expenses paid from fund with filtering
    GetExpensesPaidFromFund(ctx context.Context, filter ExpenseFilter) ([]Expense, error)
}

type CreateExpenseFromFundRequest struct {
    Date          time.Time
    CategoryID    primitive.ObjectID
    Amount        float64
    Description   string
    Vendor        string
    Notes         string
    UserID        primitive.ObjectID
    UserName      string
    UserRole      string
}

type RestockFromFundRequest struct {
    IngredientID primitive.ObjectID
    Quantity     float64
    CostPerUnit  float64
    Reason       string
    UserID       primitive.ObjectID
    UserName     string
    UserRole     string
}

type ExpenseFromFundResult struct {
    Expense         *Expense
    FundTransaction *FundTransaction
}

type RestockFromFundResult struct {
    RestockRecord   *IngredientRestockRecord
    Expense         *Expense
    FundTransaction *FundTransaction
    UpdatedStock    float64
}
```

### API Endpoints

#### Expense Endpoints
```
POST   /api/expenses/from-fund          - Create expense paid from fund
GET    /api/expenses?paid_from_fund=true - List expenses paid from fund
GET    /api/expenses/:id                 - Get expense detail (includes fund transaction link)
```

#### Ingredient Endpoints
```
POST   /api/ingredients/:id/restock/from-fund - Restock ingredient paid from fund
GET    /api/ingredients/:id/restock-history    - Get restock history (includes fund info)
```

#### Fund Endpoints
```
GET    /api/fund/balance                       - Get current fund balance
GET    /api/fund/transactions                  - List fund transactions (includes source links)
GET    /api/fund/transactions/:id              - Get fund transaction detail
```

## Data Models

### Database Schema Changes

#### expenses Collection
```javascript
{
  _id: ObjectId,
  date: ISODate,
  category_id: ObjectId,
  amount: Number,
  description: String,
  payment_method: String,  // "cash", "bank", "card", "fund"
  vendor: String,
  notes: String,
  source_type: String,     // "ingredient", "facility", "manual"
  source_id: ObjectId,
  
  // NEW FIELDS
  paid_from_fund: Boolean,      // Flag indicating payment from fund
  fund_transaction_id: ObjectId, // Reference to fund_transactions
  
  created_by: String,
  created_at: ISODate,
  updated_at: ISODate
}
```

#### fund_transactions Collection
```javascript
{
  _id: ObjectId,
  type: String,  // "deposit", "withdrawal"
  cash_amount: Number,
  transfer_amount: Number,
  total_amount: Number,
  reason: String,
  
  // NEW FIELDS
  source_type: String,  // "expense", "ingredient", null for manual transactions
  source_id: ObjectId,  // Reference to source record
  
  performed_by: ObjectId,
  performed_by_name: String,
  performed_by_role: String,
  timestamp: ISODate,
  balance_before: {
    cash: Number,
    transfer: Number,
    total: Number
  },
  balance_after: {
    cash: Number,
    transfer: Number,
    total: Number
  },
  created_at: ISODate,
  updated_at: ISODate
}
```

#### ingredient_restock_history Collection (NEW)
```javascript
{
  _id: ObjectId,
  ingredient_id: ObjectId,
  quantity: Number,
  cost_per_unit: Number,
  total_cost: Number,
  
  // Fund integration fields
  paid_from_fund: Boolean,
  expense_id: ObjectId,          // Reference to expenses
  fund_transaction_id: ObjectId, // Reference to fund_transactions
  
  performed_by: String,
  performed_by_name: String,
  reason: String,
  created_at: ISODate
}
```

### Database Indexes

```javascript
// expenses collection
db.expenses.createIndex({ "paid_from_fund": 1, "created_at": -1 })
db.expenses.createIndex({ "fund_transaction_id": 1 })

// fund_transactions collection
db.fund_transactions.createIndex({ "source_type": 1, "source_id": 1 }, { unique: true, sparse: true })
db.fund_transactions.createIndex({ "timestamp": -1 })

// ingredient_restock_history collection
db.ingredient_restock_history.createIndex({ "ingredient_id": 1, "created_at": -1 })
db.ingredient_restock_history.createIndex({ "fund_transaction_id": 1 })
```

### Data Relationships

```
Expense (paid_from_fund=true)
    ├─> fund_transaction_id ──> FundTransaction
    └─< source_id ──────────────┘ (bidirectional)

IngredientRestockRecord (paid_from_fund=true)
    ├─> expense_id ──────────> Expense
    ├─> fund_transaction_id ──> FundTransaction
    └─< source_id ──────────────┘ (bidirectional)

FundTransaction
    └─> source_type + source_id ──> Expense | IngredientRestockRecord
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Fund Balance Validation

*For any* withdrawal request (expense or ingredient restock), if the requested amount exceeds the current fund balance, the system SHALL reject the request and return an error indicating the current balance and requested amount.

**Validates: Requirements 1.2, 2.2, 6.2, 6.3**

### Property 2: Withdrawal Amount Consistency

*For any* expense or ingredient restock paid from fund, the created fund withdrawal transaction SHALL have a total_amount exactly equal to the expense amount or restock cost.

**Validates: Requirements 1.3, 2.4**

### Property 3: Balance Update Invariant

*For any* fund withdrawal transaction, the balance_after SHALL equal balance_before minus the withdrawal total_amount (balance_after.total = balance_before.total - total_amount).

**Validates: Requirements 1.4**

### Property 4: Bidirectional Linking - Expense to Fund

*For any* expense with paid_from_fund=true, the expense SHALL have a populated fund_transaction_id, and the referenced fund transaction SHALL have source_type="expense" and source_id equal to the expense ID.

**Validates: Requirements 1.5, 1.6**

### Property 5: Bidirectional Linking - Ingredient to Fund

*For any* ingredient restock paid from fund, the restock record SHALL have a populated fund_transaction_id, and the referenced fund transaction SHALL have source_type="ingredient" and source_id equal to the restock record ID.

**Validates: Requirements 2.6, 2.7**

### Property 6: Expense Creation for Ingredient Restock

*For any* ingredient restock paid from fund, an expense record SHALL be created with category "ingredient purchase" and amount equal to the restock total cost.

**Validates: Requirements 2.3**

### Property 7: Stock Quantity Update

*For any* ingredient restock operation, the ingredient's stock quantity after the operation SHALL equal the quantity before plus the restock quantity.

**Validates: Requirements 2.5**

### Property 8: Fund-Paid Expense Filter

*For any* expense list filtered by paid_from_fund=true, all returned expenses SHALL have paid_from_fund=true and fund_transaction_id populated.

**Validates: Requirements 4.3**

### Property 9: Audit Trail Completeness

*For any* fund transaction created from an expense or ingredient restock, the transaction SHALL have all audit fields populated: performed_by, performed_by_name, performed_by_role, timestamp, reason, balance_before, and balance_after.

**Validates: Requirements 7.1, 7.2, 7.3, 7.4, 7.5, 7.6**

### Property 10: Source Type Validation

*For any* fund transaction with source_type populated, the source_type SHALL be one of the allowed values ("expense", "ingredient"), and source_id SHALL also be populated.

**Validates: Requirements 8.3, 8.5**

### Property 11: Transaction Atomicity - Expense Creation Failure

*For any* expense creation from fund operation, if the fund transaction creation fails, no expense record SHALL be created in the database.

**Validates: Requirements 9.1**

### Property 12: Transaction Atomicity - Fund Transaction Rollback

*For any* expense creation from fund operation, if the expense creation fails after the fund transaction is created, the fund transaction SHALL be deleted and the fund balance SHALL be restored to its original value.

**Validates: Requirements 9.2**

### Property 13: Transaction Atomicity - Ingredient Restock Rollback

*For any* ingredient restock from fund operation, if any step fails (stock update, expense creation, or fund transaction creation), all changes SHALL be rolled back: no restock record, no expense, no fund transaction, and stock quantity unchanged.

**Validates: Requirements 9.3**

### Property 14: Payment Method Consistency

*For any* expense with paid_from_fund=true, the payment_method SHALL be set to "fund", and for any fund transaction with source_type="expense", the referenced expense SHALL have payment_method="fund".

**Validates: Requirements 10.1, 10.3**

### Property 15: Immutability - Fund Transaction Link

*For any* expense that has fund_transaction_id populated, any attempt to modify the paid_from_fund flag or fund_transaction_id SHALL be rejected.

**Validates: Requirements 11.1**

### Property 16: Uniqueness - No Duplicate Fund Transactions

*For any* combination of source_type and source_id, at most one fund transaction SHALL exist with that combination.

**Validates: Requirements 11.2**

### Property 17: Immutability - Expense Amount

*For any* expense with paid_from_fund=true, any attempt to modify the amount field SHALL be rejected.

**Validates: Requirements 11.3**

## Error Handling

### Error Types

```go
var (
    ErrInsufficientFundBalance = errors.New("insufficient fund balance")
    ErrInvalidSourceType       = errors.New("invalid source type")
    ErrDuplicateFundTransaction = errors.New("duplicate fund transaction for source")
    ErrPaymentMethodMismatch   = errors.New("payment method does not match fund payment")
    ErrCannotModifyFundExpense = errors.New("cannot modify expense paid from fund")
    ErrTransactionRollbackFailed = errors.New("transaction rollback failed")
)
```

### Error Handling Strategy

#### Insufficient Balance
```go
if requiredAmount > currentBalance.Total {
    return nil, fmt.Errorf("%w: required=%.2f, available=%.2f", 
        ErrInsufficientFundBalance, requiredAmount, currentBalance.Total)
}
```

#### Transaction Rollback
```go
// Use MongoDB session for transaction
session, err := client.StartSession()
if err != nil {
    return nil, err
}
defer session.EndSession(ctx)

err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
    // Start transaction
    if err := session.StartTransaction(); err != nil {
        return err
    }
    
    // Perform operations
    expense, err := createExpense(sc, req)
    if err != nil {
        session.AbortTransaction(sc)
        return err
    }
    
    fundTx, err := createFundTransaction(sc, expense)
    if err != nil {
        session.AbortTransaction(sc)
        return err
    }
    
    // Commit transaction
    return session.CommitTransaction(sc)
})
```

#### Validation Errors
- Return clear error messages with context
- Include current state information (e.g., current balance)
- Log validation failures for audit purposes

#### Duplicate Prevention
```go
// Check for existing fund transaction with same source
existing, err := fundRepo.FindBySource(ctx, sourceType, sourceID)
if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
    return nil, err
}
if existing != nil {
    return nil, ErrDuplicateFundTransaction
}
```

## Testing Strategy

### Unit Testing

Unit tests will focus on specific examples, edge cases, and error conditions:

1. **Validation Tests**
   - Test insufficient balance rejection with specific amounts
   - Test invalid source type rejection
   - Test payment method mismatch detection

2. **Edge Cases**
   - Zero balance scenarios
   - Exact balance match scenarios
   - Concurrent withdrawal attempts
   - Very large amounts (boundary testing)

3. **Error Handling**
   - Test each error type is returned correctly
   - Test error messages contain required information
   - Test rollback behavior on specific failure points

4. **Integration Tests**
   - Test complete flow from API to database
   - Test MongoDB transaction behavior
   - Test concurrent operations with locking

### Property-Based Testing

Property tests will verify universal properties across all inputs using a property-based testing library (e.g., `gopter` for Go). Each test will run minimum 100 iterations with randomized inputs.

**Test Configuration:**
```go
parameters := gopter.DefaultTestParameters()
parameters.MinSuccessfulTests = 100
```

**Property Test Examples:**

1. **Balance Update Invariant (Property 3)**
```go
// Feature: fund-expense-integration, Property 3: Balance update invariant
properties.Property("balance after withdrawal equals balance before minus amount", 
    prop.ForAll(
        func(initialBalance, withdrawalAmount float64) bool {
            // Given: initial balance and withdrawal amount
            // When: withdrawal is processed
            // Then: balance_after = balance_before - withdrawal_amount
        },
        gen.Float64Range(100, 10000),  // initial balance
        gen.Float64Range(1, 100),      // withdrawal amount
    ))
```

2. **Bidirectional Linking (Property 4)**
```go
// Feature: fund-expense-integration, Property 4: Bidirectional linking expense to fund
properties.Property("expense and fund transaction are bidirectionally linked",
    prop.ForAll(
        func(expenseAmount float64, description string) bool {
            // Given: expense paid from fund
            // When: expense is created
            // Then: expense.fund_transaction_id points to fund transaction
            //   AND fund_transaction.source_type = "expense"
            //   AND fund_transaction.source_id = expense.id
        },
        gen.Float64Range(1, 1000),
        gen.AlphaString(),
    ))
```

3. **Transaction Atomicity (Property 11)**
```go
// Feature: fund-expense-integration, Property 11: Transaction atomicity on failure
properties.Property("no expense created if fund transaction fails",
    prop.ForAll(
        func(expenseData ExpenseData) bool {
            // Given: expense creation request
            // When: fund transaction creation fails (simulated)
            // Then: no expense record exists in database
        },
        genExpenseData(),
    ))
```

4. **Uniqueness Constraint (Property 16)**
```go
// Feature: fund-expense-integration, Property 16: No duplicate fund transactions
properties.Property("only one fund transaction per source",
    prop.ForAll(
        func(sourceType string, sourceID primitive.ObjectID) bool {
            // Given: existing fund transaction with source_type and source_id
            // When: attempt to create another with same source_type and source_id
            // Then: creation is rejected with error
        },
        gen.OneConstOf("expense", "ingredient"),
        genObjectID(),
    ))
```

### Test Coverage Goals

- Unit test coverage: >80% for service layer
- Property test coverage: All 17 correctness properties
- Integration test coverage: All API endpoints
- Error path coverage: All error types and rollback scenarios

### Testing Tools

- **Unit Testing**: Go standard `testing` package
- **Property-Based Testing**: `gopter` library
- **Mocking**: `gomock` for repository mocks
- **Database Testing**: `testcontainers-go` for MongoDB integration tests
- **API Testing**: `httptest` for HTTP handler tests

