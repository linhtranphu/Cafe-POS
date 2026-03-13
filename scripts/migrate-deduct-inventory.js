/**
 * Migration: Set deduct_inventory on existing menu item ingredients
 *
 * Rules:
 *   - batch ingredients (ingredient_type = "batch")  → deduct_inventory = true  (preserve existing behavior)
 *   - raw ingredients (ingredient_type = "raw" | "") → deduct_inventory = false (preserve existing behavior — raw was never deducted on order)
 *
 * Run with:
 *   mongosh <connection-string> migrate-deduct-inventory.js
 *   # or inside mongosh shell:
 *   load("scripts/migrate-deduct-inventory.js")
 */

// ─── Top-level ingredients (single-size items) ────────────────────────────────

// 1. batch top-level ingredients → true
var r1 = db.menu_items.updateMany(
  {
    "ingredients": { $elemMatch: { "ingredient_type": "batch", "deduct_inventory": { $exists: false } } }
  },
  { $set: { "ingredients.$[elem].deduct_inventory": true } },
  { arrayFilters: [{ "elem.ingredient_type": "batch", "elem.deduct_inventory": { $exists: false } }] }
)
print("Top-level batch → true:", r1.modifiedCount, "documents updated")

// 2. raw / unset top-level ingredients → false
var r2 = db.menu_items.updateMany(
  {
    "ingredients": { $elemMatch: { "deduct_inventory": { $exists: false } } }
  },
  { $set: { "ingredients.$[elem].deduct_inventory": false } },
  { arrayFilters: [{ "elem.deduct_inventory": { $exists: false } }] }
)
print("Top-level raw → false:", r2.modifiedCount, "documents updated")

// ─── Variant ingredients (multi-size items) ───────────────────────────────────

// 3. batch variant ingredients → true
var r3 = db.menu_items.updateMany(
  {
    "variants": { $elemMatch: { "ingredients": { $elemMatch: { "ingredient_type": "batch", "deduct_inventory": { $exists: false } } } } }
  },
  { $set: { "variants.$[].ingredients.$[elem].deduct_inventory": true } },
  { arrayFilters: [{ "elem.ingredient_type": "batch", "elem.deduct_inventory": { $exists: false } }] }
)
print("Variant batch → true:", r3.modifiedCount, "documents updated")

// 4. raw / unset variant ingredients → false
var r4 = db.menu_items.updateMany(
  {
    "variants": { $elemMatch: { "ingredients": { $elemMatch: { "deduct_inventory": { $exists: false } } } } }
  },
  { $set: { "variants.$[].ingredients.$[elem].deduct_inventory": false } },
  { arrayFilters: [{ "elem.deduct_inventory": { $exists: false } }] }
)
print("Variant raw → false:", r4.modifiedCount, "documents updated")

print("Migration complete.")
