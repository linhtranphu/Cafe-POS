# Shop Settings Migration

## Overview

This migration creates the `shop_settings` collection and initializes it with default values for the Menu Cost & Profit Analysis feature.

## What it does

1. **Creates shop_settings collection** if it doesn't exist
2. **Adds default shop settings** with:
   - `shop_name`: "My Cafe" (default name)
   - `low_margin_threshold`: 20.0 (default 20% threshold for low margin warnings)
   - `created_at`: Current timestamp
   - `updated_at`: Current timestamp

3. **Updates existing settings** if they already exist:
   - Adds `low_margin_threshold` field if missing (default: 20.0)
   - Updates `updated_at` timestamp

## Schema

```javascript
{
  _id: ObjectId,
  shop_name: String,              // Shop name
  low_margin_threshold: Number,   // Threshold for low margin warning (default: 20.0)
  created_at: Date,
  updated_at: Date
}
```

## Usage

### Build the migration

```bash
cd backend/cmd/migrate
go build -o create_shop_settings create_shop_settings.go
```

### Run the migration

```bash
# Using default MongoDB URI (mongodb://localhost:27017)
./create_shop_settings

# Using custom MongoDB URI
MONGODB_URI="mongodb://localhost:27017" ./create_shop_settings
```

### Expected Output

```
=== Shop Settings Migration ===
Creating shop_settings collection with default values...
No existing settings found. Creating default shop settings...
✓ Created default shop settings
  - Shop Name: My Cafe
  - Low Margin Threshold: 20.0%

=== Migration Complete ===
Shop settings collection is ready

Default values:
  - low_margin_threshold: 20.0 (20%)

You can update these values via the settings API endpoint
```

## Verification

After running the migration, verify the settings were created:

```bash
# Connect to MongoDB
mongosh cafe_pos

# Check shop_settings collection
db.shop_settings.find().pretty()
```

Expected result:
```javascript
{
  _id: ObjectId("..."),
  shop_name: "My Cafe",
  low_margin_threshold: 20.0,
  created_at: ISODate("2024-01-15T10:00:00.000Z"),
  updated_at: ISODate("2024-01-15T10:00:00.000Z")
}
```

## Low Margin Threshold

The `low_margin_threshold` field is used by the profit analysis feature to detect menu items with low profit margins:

- **Default value**: 20.0 (20%)
- **Purpose**: Menu items with profit margin below this threshold will be flagged with a "low_margin" warning
- **Configurable**: Managers can update this value via the settings API endpoint

### Example Usage

- If `low_margin_threshold = 20.0`:
  - Menu item with 25% profit margin → No warning
  - Menu item with 18% profit margin → Low margin warning (yellow)
  - Menu item with -5% profit margin → Loss warning (red)

## Rollback

If you need to remove the shop_settings collection:

```bash
mongosh cafe_pos --eval "db.shop_settings.drop()"
```

## Notes

- This migration is **idempotent** - safe to run multiple times
- If settings already exist, it only adds missing fields
- There should typically be only one document in the shop_settings collection
- The low_margin_threshold can be updated later via the PATCH /api/settings endpoint

## Related Files

- Domain model: `backend/domain/settings/shop_settings.go`
- Repository: `backend/infrastructure/mongodb/shop_settings_repository.go`
- Migration script: `backend/cmd/migrate/create_shop_settings.go`

## Requirements

This migration implements:
- **Requirement 3.2**: System SHALL allow managers to configure the low_margin_threshold value
- **Requirement 3.3**: System SHALL allow managers to configure the low_margin_threshold value per shop
