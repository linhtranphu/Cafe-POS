# Menu Cost Tracking Migration

This migration adds cost tracking fields to the MenuItem model to support the Menu Cost & Profit Analysis feature.

## Changes

### 1. MenuItem Model Extensions

Added the following fields to `MenuItem` struct in `backend/domain/menu/menu.go`:

- `CurrentCost` (float64): Current cost of the menu item based on ingredient costs
- `CostLastCalculatedAt` (time.Time): Timestamp of last cost calculation
- `CostStatus` (CostStatus): Status of cost calculation (FINAL, ESTIMATED, INCOMPLETE)

### 2. CostStatus Enum

Added new enum type `CostStatus` with three values:
- `FINAL`: Cost has been calculated and finalized
- `ESTIMATED`: Cost is estimated (shift not closed)
- `INCOMPLETE`: Missing ingredient cost data

### 3. Database Indexes

Created indexes on the `menu_items` collection:
- `category`: For filtering menu items by category
- `cost_status`: For filtering by cost status
- `current_cost`: For sorting by cost

## Running the Migration

### Prerequisites

- MongoDB running and accessible
- `MONGODB_URI` environment variable set (optional, defaults to `mongodb://localhost:27017`)

### Steps

1. Build the migration script:
```bash
cd backend
go build -o migrate_menu_cost ./cmd/migrate/add_menu_cost_fields.go
```

2. Run the migration:
```bash
./migrate_menu_cost
```

Or run directly with go:
```bash
go run ./cmd/migrate/add_menu_cost_fields.go
```

### What the Migration Does

1. **Updates existing menu items**: Adds the new cost tracking fields to all existing menu items with default values:
   - `current_cost`: 0.0
   - `cost_last_calculated_at`: Current timestamp
   - `cost_status`: INCOMPLETE (since costs haven't been calculated yet)

2. **Creates indexes**: Creates performance indexes on the new fields for efficient querying

### Expected Output

```
Connecting to MongoDB: mongodb://localhost:27017
✅ MongoDB connected successfully
📝 Adding cost tracking fields to existing menu items...
✅ Updated X menu items with cost tracking fields
📝 Creating indexes for menu_items collection...
✅ Indexes created successfully
🎉 Migration completed successfully!
```

## Verification

After running the migration, you can verify the changes in MongoDB:

```javascript
// Check a menu item has the new fields
db.menu_items.findOne()

// Check indexes were created
db.menu_items.getIndexes()
```

## Rollback

If you need to rollback this migration:

```javascript
// Remove the new fields from all menu items
db.menu_items.updateMany(
  {},
  {
    $unset: {
      current_cost: "",
      cost_last_calculated_at: "",
      cost_status: ""
    }
  }
)

// Drop the indexes
db.menu_items.dropIndex("category_1")
db.menu_items.dropIndex("cost_status_1")
db.menu_items.dropIndex("current_cost_1")
```

## Next Steps

After running this migration:

1. Implement the Cost Calculator Service to calculate actual costs
2. Run cost calculation for all menu items to populate `current_cost`
3. Update the frontend to display cost and profit information
