# Ingredient Conversion Rate and Wastage Percentage Migration

## Overview

This migration adds two new fields to the `ingredients` collection:
- `conversion_rate`: Default value 1.0 - Used for unit conversion in cost calculations
- `wastage_percentage`: Default value 0.0 - Used to account for wastage in cost calculations

## Purpose

These fields support the Menu Cost & Profit Analysis feature (Requirement 10.1, 10.2) by allowing accurate cost calculations that account for:
1. Unit conversions (e.g., stock in kg but recipe uses grams)
2. Wastage factors (e.g., 10% wastage means multiply cost by 1.10)

## Cost Calculation Formula

When calculating menu item cost with conversion and wastage:
```
cost = (quantity * conversion_rate * cost_per_unit) * (1 + wastage_percentage/100)
```

## Running the Migration

```bash
# From the backend directory
cd backend

# Build the migration
go build -o add_ingredient_conversion_wastage ./cmd/migrate/add_ingredient_conversion_wastage.go

# Run the migration (requires MongoDB to be running)
MONGODB_URI="mongodb://admin:password123@localhost:27017" ./add_ingredient_conversion_wastage
```

## What the Migration Does

1. Connects to MongoDB using the provided URI
2. Updates all ingredients that don't have `conversion_rate` or `wastage_percentage` fields
3. Sets default values:
   - `conversion_rate`: 1.0 (no conversion)
   - `wastage_percentage`: 0.0 (no wastage)
4. Verifies the migration by counting updated documents

## Expected Output

```
Connecting to MongoDB...
Connected to MongoDB successfully!
Starting migration: Add conversion_rate and wastage_percentage to ingredients...
Migration completed successfully!
- Matched documents: X
- Modified documents: X

Verifying migration...
- Total ingredients: X
- Ingredients with conversion_rate = 1.0: X
- Ingredients with wastage_percentage = 0.0: X

✓ Migration verified successfully! All ingredients have default values.
```

## Schema Changes

### Before
```javascript
{
  _id: ObjectId,
  name: String,
  category: String,
  unit: String,
  quantity: Number,
  min_stock: Number,
  cost_per_unit: Number,
  supplier: String,
  created_at: Date,
  updated_at: Date
}
```

### After
```javascript
{
  _id: ObjectId,
  name: String,
  category: String,
  unit: String,
  quantity: Number,
  min_stock: Number,
  cost_per_unit: Number,
  supplier: String,
  conversion_rate: Number,      // NEW - Default: 1.0
  wastage_percentage: Number,   // NEW - Default: 0.0
  created_at: Date,
  updated_at: Date
}
```

## API Changes

### CreateIngredientRequest
New optional fields:
- `conversion_rate` (optional, default: 1.0)
- `wastage_percentage` (optional, default: 0.0)

### UpdateIngredientRequest
New optional fields:
- `conversion_rate` (optional)
- `wastage_percentage` (optional)

## Example Usage

### Creating an ingredient with conversion and wastage
```json
{
  "name": "Milk",
  "category": "Dairy",
  "unit": "L",
  "quantity": 10,
  "cost_per_unit": 50000,
  "conversion_rate": 1000,
  "wastage_percentage": 10
}
```

This means:
- Stock is tracked in liters (L)
- Recipe uses milliliters (1L = 1000ml, so conversion_rate = 1000)
- 10% wastage is expected (wastage_percentage = 10)

### Cost calculation example
If a recipe uses 150ml of milk:
```
cost = (0.15 * 1000 * 50000) * (1 + 10/100)
     = (150 * 50000) * 1.10
     = 7,500,000 * 1.10
     = 8,250,000 VND
```

## Rollback

If you need to rollback this migration:

```javascript
// Connect to MongoDB
use cafe_pos

// Remove the new fields
db.ingredients.updateMany(
  {},
  {
    $unset: {
      conversion_rate: "",
      wastage_percentage: ""
    }
  }
)
```

## Related Files

- Domain model: `backend/domain/ingredient/ingredient.go`
- Service: `backend/application/services/ingredient.go`
- Migration script: `backend/cmd/migrate/add_ingredient_conversion_wastage.go`

## Requirements

- Requirement 10.1: Ingredient Unit Conversion
- Requirement 10.2: Wastage Factor
- Task 1.3: Extend Ingredient model với conversion và wastage fields
