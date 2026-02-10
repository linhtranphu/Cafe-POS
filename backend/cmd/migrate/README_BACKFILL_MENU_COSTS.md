# Backfill Current Cost for Menu Items

This script calculates and backfills the `current_cost` field for all existing menu items in the database.

## What This Script Does

For each menu item:
1. Fetches the menu item's ingredients
2. Calculates the current cost using the formula:
   ```
   cost = sum(quantity * cost_per_unit * conversion_rate * (1 + wastage_percentage/100))
   ```
3. Determines the cost_status:
   - **FINAL**: All ingredients have valid cost_per_unit values (or no ingredients)
   - **INCOMPLETE**: One or more ingredients are missing cost_per_unit
4. Updates the menu item with:
   - `current_cost`: Calculated cost value
   - `cost_last_calculated_at`: Current timestamp
   - `cost_status`: FINAL or INCOMPLETE

## Prerequisites

- Task 19.1 must be completed (schema migration)
- MongoDB running and accessible
- Ingredients should have `cost_per_unit` values set (optional but recommended)

## How to Run

### Build and Run
```bash
cd backend/cmd/migrate
go build -o backfill_menu_item_costs backfill_menu_item_costs.go
./backfill_menu_item_costs
```

### Or Run Directly
```bash
cd backend
go run cmd/migrate/backfill_menu_item_costs.go
```

## Expected Output

```
🔌 Connecting to MongoDB: mongodb://localhost:27017
✅ Connected to database: cafe_pos

============================================================
  BACKFILL CURRENT_COST FOR MENU ITEMS
============================================================

📝 Fetching all menu items...
   Found 50 menu items

🔄 Calculating costs for menu items...
   Progress: 10/50 items processed
   Progress: 20/50 items processed
   Progress: 30/50 items processed
   Progress: 40/50 items processed
   Progress: 50/50 items processed

============================================================
  ✅ BACKFILL COMPLETED
============================================================

📊 Backfill Summary:
   Total menu items:           50
   ✅ Successfully calculated:  42
   📦 No ingredients (cost=0):  3
   ⚠️  Incomplete (missing cost): 5
   ❌ Errors:                   0

📈 Database Verification:
   FINAL status:       45 items
   INCOMPLETE status:  5 items

⚠️  Action Required:
   5 menu items have INCOMPLETE status due to missing ingredient costs.
   Please update ingredient cost_per_unit values and re-run this script.

   Examples of incomplete items:
      - Special Latte (Coffee)
      - Premium Tea (Tea)
      - Smoothie Bowl (Food)

📝 Next Steps:
   1. Update missing ingredient costs in the system
   2. Re-run this backfill script to recalculate incomplete items
   3. Run task 19.3 to backfill accounting_cost for historical orders
```

## Handling Incomplete Items

If some menu items have INCOMPLETE status:

1. **Identify missing ingredient costs:**
   ```javascript
   // In MongoDB shell
   db.menu_items.find({cost_status: "INCOMPLETE"}, {name: 1, ingredients: 1})
   ```

2. **Update ingredient costs:**
   - Use the ingredient management interface
   - Or update directly in MongoDB:
   ```javascript
   db.ingredients.updateOne(
     {name: "Espresso"},
     {$set: {cost_per_unit: 200}}
   )
   ```

3. **Re-run this backfill script:**
   ```bash
   ./backfill_menu_item_costs
   ```

## Safety

- This script is **idempotent** - it can be run multiple times safely
- Recalculates costs based on current ingredient prices
- Does not modify ingredient data
- Does not affect historical accounting_cost data

## Performance

- Processes items sequentially to avoid overwhelming the database
- Progress updates every 10 items
- Typical performance: ~100 items per second

## Validation

After running this script, verify:
1. All menu items have `current_cost` field set
2. All menu items have `cost_status` field set
3. Items with complete ingredient data have status = FINAL
4. Items with missing ingredient costs have status = INCOMPLETE

Run the verification script (task 19.4) to confirm:
```bash
go run cmd/migrate/verify_menu_cost_migration.go
```

## Requirements Validated

- Requirement 1.1: Calculate current_cost for menu items
- Requirement 1.2: Use current cost_per_unit values
- Requirement 1.5: Mark items with missing costs as INCOMPLETE
- Requirement 1.7: Round costs to 2 decimal places
- Requirement 10.1, 10.2, 10.4: Apply conversion_rate and wastage_percentage
