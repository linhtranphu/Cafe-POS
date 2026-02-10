# Verify Menu Cost & Profit Analysis Migration

This script verifies that all migration steps have been completed successfully.

## What This Script Verifies

### 1. Menu Items Schema (Task 19.1)
- ✅ All menu items have `current_cost` field
- ✅ All menu items have `cost_last_calculated_at` field
- ✅ All menu items have `cost_status` field

### 2. Ingredients Schema (Task 19.1)
- ✅ All ingredients have `conversion_rate` field
- ✅ All ingredients have `wastage_percentage` field

### 3. Order Items Collection (Task 19.1)
- ✅ `order_items` collection exists
- ✅ Collection has documents (if backfill was run)

### 4. Operating Expenses Collection (Task 19.1)
- ✅ `operating_expenses` collection exists

### 5. Shop Settings Schema (Task 19.1)
- ✅ All shop settings have `low_margin_threshold` field

### 6. Indexes (Task 19.1)
- ✅ menu_items: idx_category, idx_cost_status, idx_current_cost
- ✅ order_items: idx_order_id, idx_menu_item_id, idx_cost_status, idx_cost_calculated_at, idx_order_menu_item
- ✅ operating_expenses: idx_period_range, idx_period_start

### 7. Menu Item Costs (Task 19.2)
- ✅ All menu items have valid cost_status (FINAL or INCOMPLETE)
- ✅ Sample items show calculated costs

### 8. Order Item Costs (Task 19.3)
- ✅ All order items have valid cost_status (FINAL, ESTIMATED, or INCOMPLETE)
- ✅ Sample items show calculated accounting costs

## How to Run

### Build and Run
```bash
cd backend/cmd/migrate
go build -o verify_menu_cost_migration verify_menu_cost_migration.go
./verify_menu_cost_migration
```

### Or Run Directly
```bash
cd backend
go run cmd/migrate/verify_menu_cost_migration.go
```

## Expected Output (All Passed)

```
🔌 Connecting to MongoDB: mongodb://localhost:27017
✅ Connected to database: cafe_pos

============================================================
  MENU COST & PROFIT ANALYSIS - MIGRATION VERIFICATION
============================================================

📝 [1/8] Verifying menu_items schema...
   ✅ All 50 menu items have cost tracking fields

📝 [2/8] Verifying ingredients schema...
   ✅ All 30 ingredients have conversion and wastage fields

📝 [3/8] Verifying order_items collection...
   ✅ order_items collection exists with 1125 documents

📝 [4/8] Verifying operating_expenses collection...
   ✅ operating_expenses collection exists with 0 documents

📝 [5/8] Verifying shop_settings schema...
   ✅ All 1 shop settings have low_margin_threshold field

📝 [6/8] Verifying indexes...
   ✅ All required indexes exist

📝 [7/8] Verifying menu item costs...
   Total menu items:       50
   FINAL status:           45
   INCOMPLETE status:      5
   ✅ All menu items have valid cost_status
   Sample items with calculated costs:
      - Cappuccino: 15000.00 VND
      - Latte: 18000.00 VND
      - Espresso: 12000.00 VND

📝 [8/8] Verifying order item costs...
   Total order items:      1125
   FINAL status:           0
   ESTIMATED status:       1050
   INCOMPLETE status:      75
   ✅ All order items have valid cost_status
   Sample order items with calculated costs:
      - Cappuccino: 15000.00 VND (ESTIMATED)
      - Latte: 18000.00 VND (ESTIMATED)
      - Espresso: 12000.00 VND (ESTIMATED)

============================================================
  ✅ ALL VERIFICATION CHECKS PASSED
============================================================
```

## Expected Output (Some Failed)

```
🔌 Connecting to MongoDB: mongodb://localhost:27017
✅ Connected to database: cafe_pos

============================================================
  MENU COST & PROFIT ANALYSIS - MIGRATION VERIFICATION
============================================================

📝 [1/8] Verifying menu_items schema...
   ❌ Only 30/50 menu items have cost tracking fields

📝 [2/8] Verifying ingredients schema...
   ✅ All 30 ingredients have conversion and wastage fields

...

============================================================
  ⚠️  SOME VERIFICATION CHECKS FAILED
============================================================
```

## Troubleshooting

### If menu_items schema check fails:
```bash
# Re-run the schema migration
go run cmd/migrate/run_all_menu_cost_migrations.go
```

### If menu item costs check fails:
```bash
# Re-run the cost backfill
go run cmd/migrate/backfill_menu_item_costs.go
```

### If order item costs check fails:
```bash
# Re-run the order cost backfill
go run cmd/migrate/backfill_order_item_costs.go
```

### If indexes check fails:
```bash
# Re-run the schema migration (includes index creation)
go run cmd/migrate/run_all_menu_cost_migrations.go
```

## Exit Codes

- **0**: All verification checks passed
- **1**: One or more verification checks failed

## Use in CI/CD

This script can be used in CI/CD pipelines to verify migrations:

```bash
# In your deployment script
go run cmd/migrate/verify_menu_cost_migration.go
if [ $? -ne 0 ]; then
  echo "Migration verification failed!"
  exit 1
fi
```

## Requirements Validated

- Requirement 1.1: Menu items have cost tracking fields
- Requirement 5.1: Order items collection exists
- Requirement 6.5.1: Operating expenses collection exists
- Requirement 10.1, 10.2: Ingredients have conversion and wastage fields
- Requirement 3.3: Shop settings have low margin threshold

## Related Scripts

This verification script should be run after:
1. `run_all_menu_cost_migrations.go` (task 19.1)
2. `backfill_menu_item_costs.go` (task 19.2)
3. `backfill_order_item_costs.go` (task 19.3)
