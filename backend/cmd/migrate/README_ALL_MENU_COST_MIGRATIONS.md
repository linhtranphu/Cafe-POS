# Menu Cost & Profit Analysis - Complete Schema Migration

This migration script runs all necessary schema changes for the Menu Cost & Profit Analysis feature.

## What This Migration Does

### 1. Menu Items Collection
- Adds `current_cost` field (default: 0.0)
- Adds `cost_last_calculated_at` field (default: current timestamp)
- Adds `cost_status` field (default: "INCOMPLETE")
- Creates indexes: category, cost_status, current_cost

### 2. Ingredients Collection
- Adds `conversion_rate` field (default: 1.0)
- Adds `wastage_percentage` field (default: 0.0)

### 3. Order Items Collection
- Creates new `order_items` collection
- Schema includes: order_id, menu_item_id, name, price, quantity, note, subtotal, accounting_cost, cost_calculated_at, cost_status, created_at
- Creates indexes: order_id, menu_item_id, cost_status, cost_calculated_at, compound index (order_id + menu_item_id)

### 4. Operating Expenses Collection
- Creates new `operating_expenses` collection
- Schema includes: period_start, period_end, staff_salary, rent, utilities, marketing_costs, other_expenses, total_expenses, created_at, updated_at
- Creates indexes: compound index (period_start + period_end), period_start

### 5. Shop Settings Collection
- Adds `low_margin_threshold` field (default: 20.0)

## How to Run

### Prerequisites
- MongoDB running and accessible
- Environment variables set (optional):
  - `MONGODB_URI` (default: mongodb://localhost:27017)
  - `MONGODB_DATABASE` (default: cafe_pos)

### Build and Run
```bash
cd backend/cmd/migrate
go build -o run_all_menu_cost_migrations run_all_menu_cost_migrations.go
./run_all_menu_cost_migrations
```

### Or Run Directly
```bash
cd backend
go run cmd/migrate/run_all_menu_cost_migrations.go
```

## Expected Output

```
🔌 Connecting to MongoDB: mongodb://localhost:27017
✅ Connected to database: cafe_pos

============================================================
  MENU COST & PROFIT ANALYSIS - SCHEMA MIGRATION
============================================================

📝 [1/5] Migrating menu_items collection...
   ✅ Added cost tracking fields to 50 menu items
   Fields added: current_cost, cost_last_calculated_at, cost_status

📝 [2/5] Migrating ingredients collection...
   ✅ Added conversion and wastage fields to 30 ingredients
   Fields added: conversion_rate (default: 1.0), wastage_percentage (default: 0.0)

📝 [3/5] Creating order_items collection...
   ✅ Created order_items collection
   Schema: order_id, menu_item_id, name, price, quantity, note, subtotal,
           accounting_cost, cost_calculated_at, cost_status, created_at

📝 [4/5] Creating operating_expenses collection...
   ✅ Created operating_expenses collection
   Schema: period_start, period_end, staff_salary, rent, utilities,
           marketing_costs, other_expenses, total_expenses, created_at, updated_at

📝 [5/5] Migrating shop_settings collection...
   ✅ Added low_margin_threshold to 1 shop settings
   Field added: low_margin_threshold (default: 20.0)

📝 Creating indexes for all collections...

   Creating indexes for menu_items...
   ✅ menu_items indexes:
      - idx_category: category
      - idx_cost_status: cost_status
      - idx_current_cost: current_cost

   Creating indexes for order_items...
   ✅ order_items indexes:
      - idx_order_id: order_id
      - idx_menu_item_id: menu_item_id
      - idx_cost_status: cost_status
      - idx_cost_calculated_at: cost_calculated_at
      - idx_order_menu_item: order_id + menu_item_id (compound)

   Creating indexes for operating_expenses...
   ✅ operating_expenses indexes:
      - idx_period_range: period_start + period_end (compound)
      - idx_period_start: period_start

============================================================
  ✅ ALL MIGRATIONS COMPLETED SUCCESSFULLY
============================================================

📊 Migration Summary:
   • menu_items: 50 documents
   • ingredients: 30 documents
   • order_items: 0 documents
   • operating_expenses: 0 documents
   • shop_settings: 1 documents

📝 Next Steps:
   1. Run task 19.2 to backfill current_cost for existing menu items
   2. Run task 19.3 to backfill accounting_cost for historical orders
   3. Run task 19.4 to verify migration completeness
```

## Safety

- This migration is **idempotent** - it can be run multiple times safely
- Existing data is preserved
- Only adds new fields with default values
- Does not modify or delete existing data

## Rollback

If you need to rollback:
1. The old schema fields remain intact
2. You can remove the new fields manually if needed:
```javascript
// In MongoDB shell
db.menu_items.updateMany({}, {$unset: {current_cost: "", cost_last_calculated_at: "", cost_status: ""}})
db.ingredients.updateMany({}, {$unset: {conversion_rate: "", wastage_percentage: ""}})
db.shop_settings.updateMany({}, {$unset: {low_margin_threshold: ""}})
db.order_items.drop()
db.operating_expenses.drop()
```

## Related Migrations

This is part of task 19.1. After running this migration:
- Run `backfill_menu_item_costs.go` (task 19.2) to calculate current_cost for existing menu items
- Run `backfill_order_item_costs.go` (task 19.3) to calculate accounting_cost for historical orders
- Run `verify_menu_cost_migration.go` (task 19.4) to verify all migrations completed successfully

## Requirements Validated

- Requirement 1.1: Menu item cost tracking fields
- Requirement 5.1: Order items collection for accounting cost
- Requirement 6.5.1: Operating expenses collection
- Requirement 10.1, 10.2: Ingredient conversion and wastage fields
- Requirement 3.3: Shop settings low margin threshold
