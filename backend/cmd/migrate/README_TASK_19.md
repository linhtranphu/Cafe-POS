# Task 19: Data Migration và Backfill - Quick Start Guide

This guide provides a quick overview of the migration process for the Menu Cost & Profit Analysis feature.

## Quick Start

Run all migrations in order:

```bash
cd backend

# 1. Schema migration (adds fields, creates collections, creates indexes)
go run cmd/migrate/run_all_menu_cost_migrations.go

# 2. Backfill menu item costs
go run cmd/migrate/backfill_menu_item_costs.go

# 3. Backfill order item costs (skip if no historical data)
go run cmd/migrate/backfill_order_item_costs.go

# 4. Verify migration
go run cmd/migrate/verify_menu_cost_migration.go
```

## What Gets Migrated

### Collections Modified
- ✅ `menu_items` - adds cost tracking fields
- ✅ `ingredients` - adds conversion and wastage fields
- ✅ `shop_settings` - adds low margin threshold

### Collections Created
- ✅ `order_items` - for accounting cost tracking
- ✅ `operating_expenses` - for operating profit analysis

### Indexes Created
- ✅ menu_items: category, cost_status, current_cost
- ✅ order_items: order_id, menu_item_id, cost_status, cost_calculated_at, compound
- ✅ operating_expenses: period_range, period_start

## Migration Scripts

| Script | Purpose | Documentation |
|--------|---------|---------------|
| `run_all_menu_cost_migrations.go` | Complete schema migration | [README_ALL_MENU_COST_MIGRATIONS.md](README_ALL_MENU_COST_MIGRATIONS.md) |
| `backfill_menu_item_costs.go` | Calculate current_cost for menu items | [README_BACKFILL_MENU_COSTS.md](README_BACKFILL_MENU_COSTS.md) |
| `backfill_order_item_costs.go` | Calculate accounting_cost for orders | [README_BACKFILL_ORDER_COSTS.md](README_BACKFILL_ORDER_COSTS.md) |
| `verify_menu_cost_migration.go` | Verify migration completeness | [README_VERIFY_MIGRATION.md](README_VERIFY_MIGRATION.md) |

## Expected Timeline

For a typical cafe database:
- Schema migration: ~5 seconds
- Menu item backfill (50 items): ~1 second
- Order item backfill (1000 orders): ~20 seconds
- Verification: ~2 seconds

**Total: ~30 seconds**

## Safety

✅ All migrations are **idempotent** - safe to run multiple times
✅ All migrations are **non-destructive** - no data is deleted
✅ All migrations are **additive** - only new fields/collections added
✅ **No downtime required** - application continues to work during migration

## Rollback

If needed, you can rollback by removing the new fields:

```javascript
// In MongoDB shell
db.menu_items.updateMany({}, {$unset: {current_cost: "", cost_last_calculated_at: "", cost_status: ""}})
db.ingredients.updateMany({}, {$unset: {conversion_rate: "", wastage_percentage: ""}})
db.shop_settings.updateMany({}, {$unset: {low_margin_threshold: ""}})
db.order_items.drop()
db.operating_expenses.drop()
```

## Troubleshooting

### Problem: Menu items have INCOMPLETE status
**Solution**: Update ingredient cost_per_unit values and re-run backfill

### Problem: Verification fails
**Solution**: Check which step failed and re-run that specific migration script

### Problem: Order items collection is empty
**Solution**: This is normal if you have no closed shifts. Run backfill after closing some shifts.

## Environment Variables

Optional environment variables:
- `MONGODB_URI` - MongoDB connection string (default: mongodb://localhost:27017)
- `MONGODB_DATABASE` - Database name (default: cafe_pos)

## Production Deployment

For production deployment:

```bash
# 1. Backup database first
mongodump --uri="mongodb://..." --out=/backup/before-migration

# 2. Run migrations
go run cmd/migrate/run_all_menu_cost_migrations.go
go run cmd/migrate/backfill_menu_item_costs.go
go run cmd/migrate/backfill_order_item_costs.go

# 3. Verify
go run cmd/migrate/verify_menu_cost_migration.go

# 4. If verification passes, deploy new application code
# 5. If verification fails, investigate and fix before deploying
```

## CI/CD Integration

Add to your deployment pipeline:

```bash
#!/bin/bash
set -e

echo "Running migrations..."
go run cmd/migrate/run_all_menu_cost_migrations.go
go run cmd/migrate/backfill_menu_item_costs.go
go run cmd/migrate/backfill_order_item_costs.go

echo "Verifying migrations..."
go run cmd/migrate/verify_menu_cost_migration.go

if [ $? -eq 0 ]; then
  echo "✅ Migration successful"
else
  echo "❌ Migration failed"
  exit 1
fi
```

## Support

For detailed information about each migration step, see:
- [TASK_19_IMPLEMENTATION_SUMMARY.md](TASK_19_IMPLEMENTATION_SUMMARY.md) - Complete implementation details
- Individual README files for each script (linked in table above)

## Requirements

- Go 1.16 or higher
- MongoDB 4.0 or higher
- Access to the database
- Sufficient disk space for new collections

## Post-Migration

After migration is complete:
1. ✅ Menu cost features will be available in the manager interface
2. ✅ Profit analysis reports will work with historical data
3. ✅ Future shift closures will automatically calculate FINAL costs
4. ✅ Operating expenses can be entered for profit analysis

## Questions?

See [TASK_19_IMPLEMENTATION_SUMMARY.md](TASK_19_IMPLEMENTATION_SUMMARY.md) for comprehensive documentation.
