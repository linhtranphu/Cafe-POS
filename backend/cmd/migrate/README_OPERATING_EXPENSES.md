# Operating Expenses Collection Migration

## Overview

This migration creates the `operating_expenses` collection for tracking operating expenses (staff salary, rent, utilities, marketing costs, etc.) used in operating profit analysis.

## What This Migration Does

1. **Creates Collection**: Creates the `operating_expenses` collection if it doesn't exist
2. **Creates Indexes**:
   - `period_range_idx`: Compound index on (period_start, period_end) for efficient period overlap queries
   - `period_start_idx`: Single field index on period_start for sorting and filtering

## Schema

```javascript
{
  _id: ObjectId,
  period_start: Date,              // Start date of the expense period
  period_end: Date,                // End date of the expense period
  staff_salary: Number,            // Total staff salary for the period
  rent: Number,                    // Rent expense for the period
  utilities: Number,               // Utilities (electricity, water, etc.)
  marketing_costs: Number,         // Marketing and advertising costs
  other_expenses: Number,          // Other miscellaneous expenses
  total_expenses: Number,          // Auto-calculated sum of all expenses
  created_at: Date,
  updated_at: Date
}
```

## Running the Migration

### Prerequisites

- Go 1.21 or higher
- MongoDB connection configured in `.env` file
- Required environment variables:
  - `MONGODB_URI`: MongoDB connection string
  - `MONGODB_DATABASE`: Database name (defaults to "cafe_pos")

### Steps

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Run the migration:
   ```bash
   go run cmd/migrate/create_operating_expenses_collection.go
   ```

3. Expected output:
   ```
   Connected to database: cafe_pos
   
   === Creating operating_expenses collection ===
   ✓ Created collection 'operating_expenses'
   
   === Creating indexes ===
   ✓ Created indexes:
     - period_range_idx (period_start, period_end)
     - period_start_idx (period_start)
   
   === Migration completed successfully ===
   ```

## Verification

After running the migration, verify the collection and indexes:

```bash
# Connect to MongoDB
mongosh "your-mongodb-uri"

# Switch to database
use cafe_pos

# Check if collection exists
show collections

# Check indexes
db.operating_expenses.getIndexes()
```

Expected indexes:
```javascript
[
  { v: 2, key: { _id: 1 }, name: '_id_' },
  { v: 2, key: { period_start: 1, period_end: 1 }, name: 'period_range_idx' },
  { v: 2, key: { period_start: 1 }, name: 'period_start_idx' }
]
```

## Usage Examples

### Creating an Operating Expense

```javascript
db.operating_expenses.insertOne({
  period_start: ISODate("2024-01-01T00:00:00Z"),
  period_end: ISODate("2024-01-31T23:59:59Z"),
  staff_salary: 2000000,
  rent: 1000000,
  utilities: 500000,
  marketing_costs: 300000,
  other_expenses: 200000,
  total_expenses: 4000000,
  created_at: new Date(),
  updated_at: new Date()
})
```

### Querying Expenses for a Date Range

```javascript
// Find expenses that overlap with January 2024
db.operating_expenses.find({
  period_start: { $lte: ISODate("2024-01-31T23:59:59Z") },
  period_end: { $gte: ISODate("2024-01-01T00:00:00Z") }
})
```

### Finding Expense for a Specific Date

```javascript
// Find expense that contains January 15, 2024
db.operating_expenses.findOne({
  period_start: { $lte: ISODate("2024-01-15T00:00:00Z") },
  period_end: { $gte: ISODate("2024-01-15T00:00:00Z") }
})
```

## Rollback

If you need to rollback this migration:

```bash
mongosh "your-mongodb-uri"
use cafe_pos
db.operating_expenses.drop()
```

**Warning**: This will delete all operating expense data. Make sure to backup first if you have production data.

## Related Files

- **Domain Model**: `backend/domain/expense/operating_expense.go`
- **Repository**: `backend/infrastructure/mongodb/operating_expense_repository.go`
- **Migration Script**: `backend/cmd/migrate/create_operating_expenses_collection.go`

## Requirements

This migration implements:
- **Requirement 6.5.1**: Operating profit calculation infrastructure
- **Requirement 6.5.2**: Operating expense input and storage

## Next Steps

After running this migration:
1. Implement the OperatingExpenseService for business logic
2. Create API endpoints for CRUD operations
3. Build frontend forms for expense input
4. Integrate with profit analysis reports
