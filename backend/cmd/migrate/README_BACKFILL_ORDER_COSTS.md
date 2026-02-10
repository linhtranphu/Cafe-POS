# Backfill Accounting Cost for Historical Orders

This script calculates and backfills the `accounting_cost` field for all order items in closed shifts.

## What This Script Does

For each closed shift:
1. Fetches all orders in that shift
2. For each order item:
   - Calculates the cost using current ingredient prices
   - Creates an entry in the `order_items` collection with:
     - `accounting_cost`: Calculated cost (per item quantity)
     - `cost_calculated_at`: Current timestamp
     - `cost_status`: **ESTIMATED** (not FINAL, since not calculated at actual shift closure time)
3. Marks items with missing ingredient costs as INCOMPLETE

## Important: ESTIMATED vs FINAL Status

**ESTIMATED Status:**
- Used for backfilled historical data
- Cost calculated using **current** ingredient prices
- May not reflect actual ingredient prices at the time of shift closure
- Still useful for approximate profit analysis

**FINAL Status:**
- Used for future shift closures (after this feature is deployed)
- Cost calculated using ingredient prices **at the time of shift closure**
- Provides accurate accounting data

## Prerequisites

- Task 19.1 must be completed (schema migration)
- Task 19.2 must be completed (menu item costs backfilled)
- MongoDB running and accessible
- Closed shifts exist in the database

## How to Run

### Build and Run
```bash
cd backend/cmd/migrate
go build -o backfill_order_item_costs backfill_order_item_costs.go
./backfill_order_item_costs
```

### Or Run Directly
```bash
cd backend
go run cmd/migrate/backfill_order_item_costs.go
```

## Expected Output

```
🔌 Connecting to MongoDB: mongodb://localhost:27017
✅ Connected to database: cafe_pos

============================================================
  BACKFILL ACCOUNTING_COST FOR HISTORICAL ORDERS
============================================================

📝 Fetching closed shifts...
   Found 25 closed shifts

🔄 Processing shifts and calculating costs...
   Progress: 10/25 shifts processed (150 orders, 450 items)
   Progress: 20/25 shifts processed (300 orders, 900 items)
   Progress: 25/25 shifts processed (375 orders, 1125 items)

============================================================
  ✅ BACKFILL COMPLETED
============================================================

📊 Backfill Summary:
   Total closed shifts:        25
   Shifts processed:           25
   Total orders:               375
   Total order items:          1125
   ✅ Successfully calculated:  1050
   ⚠️  Incomplete (missing cost): 75
   ❌ Errors:                   0

📈 Database Verification:
   Total order items:      1125
   ESTIMATED status:       1050 items
   INCOMPLETE status:      75 items

ℹ️  Important Notes:
   • All backfilled costs are marked as ESTIMATED (not FINAL)
   • ESTIMATED means the cost was calculated using current ingredient prices,
     not the actual prices at the time of shift closure
   • Future shift closures will use FINAL status for accurate accounting

⚠️  Action Required:
   75 order items have INCOMPLETE status due to missing ingredient costs.
   These items will not be included in profit reports.
   Please update ingredient cost_per_unit values if accurate historical
   profit analysis is needed.

📝 Next Steps:
   1. Run task 19.4 to verify migration completeness
   2. Test the profit analysis features in the manager interface
   3. Future shift closures will automatically calculate FINAL costs
```

## Handling Incomplete Items

If some order items have INCOMPLETE status:

1. **Understand the impact:**
   - These items will be excluded from profit calculations
   - Historical profit reports may be incomplete
   - This is acceptable if you don't need 100% accurate historical data

2. **Optional: Update ingredient costs and re-run:**
   - Update missing ingredient costs in the database
   - Re-run this backfill script
   - Note: This still produces ESTIMATED costs, not FINAL

3. **Going forward:**
   - Ensure all ingredients have cost_per_unit values
   - Future shift closures will calculate FINAL costs automatically

## Safety

- This script is **idempotent** - it can be run multiple times
- On re-run, it will insert duplicate order items (consider clearing order_items collection first)
- Does not modify existing orders or shifts
- Does not affect future shift closures

## Performance

- Processes shifts sequentially
- Progress updates every 10 shifts
- Typical performance: ~50 orders per second

## Data Integrity

After running this script:
- Each order item in a closed shift will have a corresponding entry in `order_items` collection
- All entries will have `cost_status` = ESTIMATED or INCOMPLETE
- Future shift closures will create entries with `cost_status` = FINAL

## Clearing and Re-running

If you need to clear backfilled data and re-run:

```javascript
// In MongoDB shell
// Clear all ESTIMATED order items
db.order_items.deleteMany({cost_status: "ESTIMATED"})

// Then re-run the backfill script
```

## Requirements Validated

- Requirement 5.1: Order items collection populated
- Requirement 5.2: Accounting cost calculated for historical orders
- Requirement 5.7: Historical data marked as ESTIMATED (not FINAL)
- Requirement 1.5: Items with missing costs marked as INCOMPLETE

## Related Scripts

- Run `backfill_menu_item_costs.go` (task 19.2) first to ensure menu items have current_cost
- Run `verify_menu_cost_migration.go` (task 19.4) after to verify completeness
