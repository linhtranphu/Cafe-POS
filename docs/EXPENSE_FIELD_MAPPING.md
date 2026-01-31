# Expense Field Mapping - Backend ↔ Frontend

This document ensures field name synchronization between backend and frontend for expense management.

## Expense Fields

### Backend: `backend/domain/expense/expense.go`
```go
type Expense struct {
    ID            primitive.ObjectID  // _id in MongoDB
    Date          time.Time          // date
    CategoryID    primitive.ObjectID  // category_id
    Amount        float64            // amount
    Description   string             // description
    PaymentMethod string             // payment_method
    Vendor        string             // vendor (optional)
    Notes         string             // notes (optional)
    CreatedAt     time.Time          // created_at
    UpdatedAt     time.Time          // updated_at
}
```

### Frontend: `frontend/src/views/ExpenseManagementView.vue`
```javascript
formData = {
    description: '',      // ✓ matches backend
    category_id: '',      // ✓ matches backend
    amount: 0,           // ✓ matches backend
    date: '',            // ✓ matches backend
    payment_method: '',  // ✓ matches backend
    vendor: '',          // ✓ matches backend
    notes: ''            // ✓ matches backend
}
```

## Category Fields

### Backend: `backend/domain/expense/expense.go`
```go
type Category struct {
    ID        primitive.ObjectID  // _id in MongoDB
    Name      string             // name
    CreatedAt time.Time          // created_at
}
```

### Frontend: `frontend/src/views/ExpenseManagementView.vue`
```javascript
category = {
    id: '',      // ✓ matches backend (from _id)
    name: ''     // ✓ matches backend
}
```

## Payment Methods

### Backend: `backend/domain/expense/expense.go`
```go
const (
    PaymentMethodCash = "cash"  // Tiền mặt
    PaymentMethodBank = "bank"  // Chuyển khoản
    PaymentMethodCard = "card"  // Thẻ
)
```

Stored as string in database:
- `"cash"` - Tiền mặt
- `"bank"` - Chuyển khoản
- `"card"` - Thẻ

### Frontend: `frontend/src/constants/expense.js`
```javascript
export const PAYMENT_METHODS = {
  CASH: 'cash',    // ✓ matches backend
  BANK: 'bank',    // ✓ matches backend
  CARD: 'card'     // ✓ matches backend
}
```

## Recurring Expense Fields

### Backend: `backend/domain/expense/expense.go`
```go
type RecurringExpense struct {
    ID          primitive.ObjectID  // _id in MongoDB
    CategoryID  primitive.ObjectID  // category_id
    Amount      float64            // amount
    Description string             // description
    Frequency   string             // frequency
    StartDate   time.Time          // start_date
    NextDue     time.Time          // next_due
    Active      bool               // active
    CreatedAt   time.Time          // created_at
}
```

### Frontend: Field names match backend
- `category_id` ✓
- `amount` ✓
- `description` ✓
- `frequency` ✓
- `start_date` ✓
- `next_due` ✓
- `active` ✓

## Recurring Frequencies

### Backend: Stored as string
- `"daily"` - Hàng ngày
- `"weekly"` - Hàng tuần
- `"monthly"` - Hàng tháng
- `"quarterly"` - Hàng quý
- `"yearly"` - Hàng năm

### Frontend: `frontend/src/constants/expense.js`
```javascript
export const RECURRING_FREQUENCIES = {
  DAILY: 'daily',          // ✓ matches backend
  WEEKLY: 'weekly',        // ✓ matches backend
  MONTHLY: 'monthly',      // ✓ matches backend
  QUARTERLY: 'quarterly',  // ✓ matches backend
  YEARLY: 'yearly'         // ✓ matches backend
}
```

## Prepaid Expense Fields

### Backend: `backend/domain/expense/expense.go`
```go
type PrepaidExpense struct {
    ID          primitive.ObjectID  // _id in MongoDB
    CategoryID  primitive.ObjectID  // category_id
    TotalAmount float64            // total_amount
    Description string             // description
    StartDate   time.Time          // start_date
    EndDate     time.Time          // end_date
    CreatedAt   time.Time          // created_at
}
```

### Frontend: Field names match backend
- `category_id` ✓
- `total_amount` ✓
- `description` ✓
- `start_date` ✓
- `end_date` ✓

## API Endpoints

### Expense Management
- `POST /api/manager/expenses` - Create expense
- `GET /api/manager/expenses` - Get expenses (with filters)
- `PUT /api/manager/expenses/:id` - Update expense
- `DELETE /api/manager/expenses/:id` - Delete expense

### Category Management
- `POST /api/manager/expense-categories` - Create category
- `GET /api/manager/expense-categories` - Get all categories
- `DELETE /api/manager/expense-categories/:id` - Delete category

### Recurring Expenses
- `POST /api/manager/recurring-expenses` - Create recurring expense
- `GET /api/manager/recurring-expenses` - Get all recurring expenses
- `DELETE /api/manager/recurring-expenses/:id` - Delete recurring expense

### Prepaid Expenses
- `POST /api/manager/prepaid-expenses` - Create prepaid expense
- `GET /api/manager/prepaid-expenses` - Get all prepaid expenses
- `DELETE /api/manager/prepaid-expenses/:id` - Delete prepaid expense

## Constants File Location

**Frontend:** `frontend/src/constants/expense.js`

This file contains:
- Payment method constants and options
- Recurring frequency constants and options
- Helper functions for label conversion

## Status: ✅ SYNCHRONIZED

All field names between backend and frontend are properly synchronized.
Last verified: 2026-01-31
