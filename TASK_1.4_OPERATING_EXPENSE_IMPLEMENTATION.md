# Task 1.4 Implementation Summary: OperatingExpense Model và Collection

## Overview

Implemented the OperatingExpense domain model, MongoDB repository, and migration script for tracking operating expenses (staff salary, rent, utilities, marketing costs, etc.) used in operating profit analysis.

## Implementation Details

### 1. Domain Model

**File**: `backend/domain/expense/operating_expense.go`

Created three main structures:

#### OperatingExpense
- Main model for storing operating expenses for a specific period
- Fields:
  - `period_start`, `period_end`: Date range for the expense period
  - `staff_salary`: Total staff salary for the period
  - `rent`: Rent expense
  - `utilities`: Utilities (electricity, water, etc.)
  - `marketing_costs`: Marketing and advertising costs
  - `other_expenses`: Miscellaneous expenses
  - `total_expenses`: Auto-calculated sum of all expenses
  - `created_at`, `updated_at`: Timestamps

#### OperatingExpenseRequest
- DTO for API validation
- All expense fields have `binding:"min=0"` validation
- Date fields use ISO 8601 string format with `binding:"required"`

#### AllocatedExpense
- Represents daily allocated expenses from a monthly/period expense
- Used when viewing daily reports but expenses are entered monthly
- Includes `expense_allocated` flag and `allocation_note` for transparency
- References original expense via `source_expense_id`

**Key Method**:
```go
func (oe *OperatingExpense) CalculateTotalExpenses()
```
- Automatically calculates total_expenses from all expense categories
- Called before saving to ensure consistency

### 2. MongoDB Repository

**File**: `backend/infrastructure/mongodb/operating_expense_repository.go`

Implemented comprehensive CRUD operations:

#### Core Methods
- `Create()`: Insert new operating expense
- `Update()`: Update existing expense by ID
- `FindByID()`: Retrieve expense by ID
- `Delete()`: Delete expense by ID
- `Upsert()`: Create or update expense for a specific period

#### Query Methods
- `FindByPeriod(startDate, endDate)`: Find expenses that overlap with date range
  - Uses overlap condition: `expense.period_start <= endDate AND expense.period_end >= startDate`
- `FindForDate(date)`: Find expense that contains a specific date
- `FindAll(startDate, endDate)`: Retrieve all expenses with optional date filter

#### Index Creation
- `CreateIndexes()`: Creates two indexes:
  - `period_range_idx`: Compound index on (period_start, period_end)
  - `period_start_idx`: Single field index on period_start

### 3. Migration Script

**File**: `backend/cmd/migrate/create_operating_expenses_collection.go`

Migration script that:
1. Connects to MongoDB using environment variables
2. Creates `operating_expenses` collection if it doesn't exist
3. Creates indexes for efficient querying
4. Provides detailed output of migration steps

**Environment Variables**:
- `MONGODB_URI`: MongoDB connection string (defaults to `mongodb://localhost:27017`)
- `MONGODB_DATABASE`: Database name (defaults to `cafe_pos`)

**Usage**:
```bash
cd backend
go run cmd/migrate/create_operating_expenses_collection.go
```

### 4. Documentation

**File**: `backend/cmd/migrate/README_OPERATING_EXPENSES.md`

Comprehensive documentation including:
- Schema definition
- Running instructions
- Verification steps
- Usage examples (MongoDB queries)
- Rollback procedure
- Related files reference

### 5. Unit Tests

**File**: `backend/domain/expense/operating_expense_test.go`

Test coverage:
- `TestOperatingExpense_CalculateTotalExpenses`: Tests total calculation with various scenarios
  - All expense types
  - Zero values
  - Partial expenses
  - Decimal values (with floating point tolerance)
- `TestAllocatedExpense_Structure`: Verifies AllocatedExpense structure
- `TestOperatingExpenseRequest_Validation`: Validates request structure

**Test Results**: All tests passing ✅

## Database Schema

```javascript
{
  _id: ObjectId,
  period_start: Date,              // Start date of expense period
  period_end: Date,                // End date of expense period
  staff_salary: Number,            // Staff salary for period
  rent: Number,                    // Rent expense
  utilities: Number,               // Utilities expense
  marketing_costs: Number,         // Marketing costs
  other_expenses: Number,          // Other expenses
  total_expenses: Number,          // Auto-calculated total
  created_at: Date,
  updated_at: Date
}
```

## Indexes

1. **period_range_idx**: `(period_start: 1, period_end: 1)`
   - Optimizes period overlap queries
   - Used by `FindByPeriod()` and `FindForDate()`

2. **period_start_idx**: `(period_start: 1)`
   - Optimizes sorting and filtering by start date
   - Used by `FindAll()` with sorting

## Key Features

### 1. Automatic Total Calculation
- `CalculateTotalExpenses()` method ensures total_expenses is always accurate
- Called automatically in `Create()` and `Update()` repository methods

### 2. Period Overlap Queries
- Efficient querying for expenses that overlap with a date range
- Supports both exact period match and partial overlap

### 3. Upsert Support
- `Upsert()` method allows creating or updating expenses for a period
- Prevents duplicate entries for the same period
- Preserves original `created_at` timestamp on update

### 4. Expense Allocation
- `AllocatedExpense` structure supports proportional allocation
- Enables viewing daily reports when expenses are entered monthly
- Transparent with `expense_allocated` flag and allocation note

## Requirements Satisfied

✅ **Requirement 6.5.1**: Operating profit calculation infrastructure
- Created data model for storing operating expenses
- Implemented repository for CRUD operations

✅ **Requirement 6.5.2**: Operating expense input and storage
- Defined OperatingExpenseRequest DTO with validation
- Implemented Upsert method for creating/updating expenses
- All expense fields validated as >= 0

## Files Created

1. `backend/domain/expense/operating_expense.go` - Domain model
2. `backend/infrastructure/mongodb/operating_expense_repository.go` - Repository
3. `backend/cmd/migrate/create_operating_expenses_collection.go` - Migration script
4. `backend/cmd/migrate/README_OPERATING_EXPENSES.md` - Documentation
5. `backend/domain/expense/operating_expense_test.go` - Unit tests
6. `TASK_1.4_OPERATING_EXPENSE_IMPLEMENTATION.md` - This summary

## Verification

### Build Verification
```bash
✅ go build ./domain/expense/operating_expense.go
✅ go build ./infrastructure/mongodb/operating_expense_repository.go
✅ go build ./cmd/migrate/create_operating_expenses_collection.go
```

### Test Verification
```bash
✅ go test ./domain/expense/operating_expense_test.go
   - All 3 test suites passing
   - 7 test cases total
```

## Next Steps

After this task, the following can be implemented:

1. **Task 4.2**: Implement OperatingExpenseService
   - Business logic for expense management
   - Expense allocation algorithm
   - Date range validation

2. **Task 8.1-8.2**: Create API endpoints
   - POST /api/operating-expenses (create/update)
   - GET /api/operating-expenses (list with filters)

3. **Task 15.1-15.4**: Build frontend form
   - OperatingExpenseForm component
   - Form validation
   - Date range picker

4. **Task 3.9**: Integrate with profit analysis
   - GetOperatingProfit method
   - Expense allocation logic
   - Operating profit calculation

## Notes

- The model uses float64 for monetary values (consistent with existing codebase)
- Timestamps use time.Time (UTC)
- Repository methods use context.Context for cancellation support
- All expense amounts must be >= 0 (enforced by validation tags)
- The Upsert method prevents duplicate periods by checking exact period match
- Floating point comparison in tests uses 0.01 tolerance for precision issues

## Status

✅ **Task 1.4 Complete**
- Domain model implemented
- Repository implemented with full CRUD
- Migration script created
- Documentation written
- Unit tests passing
- Code compiles successfully
