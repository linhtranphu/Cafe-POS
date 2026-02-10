# Order Items Collection Migration

## Overview

This migration creates the `order_items` collection for tracking accounting costs of menu items sold in orders. This is part of the Menu Cost & Profit Analysis feature.

## Purpose

The `order_items` collection stores individual order items with their accounting cost calculated at shift closure time. This allows for:

- Accurate profit/loss reporting based on historical costs
- Cost immutability after shift closure (accounting requirement)
- Efficient querying and aggregation for profit analysis reports
- Separation of concerns between order data and cost tracking

## Schema

```javascript
{
  _id: ObjectId,
  order_id: ObjectId,              // Reference to orders collection
  menu_item_id: ObjectId,          // Reference to menu_items collection
  name: String,                    // Menu item name (denormalized)
  price: Number,                   // Price at time of order
  quantity: Number,                // Quantity ordered
  note: String,                    // Optional note
  subtotal: Number,                // price * quantity
  
  // Accounting cost fields
  accounting_cost: Number,         // Cost per item (calculated at shift closure)
  cost_calculated_at: Date,        // Timestamp when cost was calculated
  cost_status: String,             // "FINAL" | "ESTIMATED" | "INCOMPLETE"
  
  created_at: Date                 // When the order item was created
}
```

## Indexes

The migration creates the following indexes for efficient querying:

1. **idx_order_id**: `order_id` - Find all items for a specific order
2. **idx_menu_item_id**: `menu_item_id` - Find all orders containing a specific menu item
3. **idx_cost_status**: `cost_status` - Find items by cost calculation status
4. **idx_cost_calculated_at**: `cost_calculated_at` - Query by cost calculation time
5. **idx_order_menu_item**: `order_id + menu_item_id` - Compound index for efficient lookups

## Running the Migration

### Prerequisites

- MongoDB running and accessible
- `MONGODB_URI` environment variable set (optional, defaults to `mongodb://localhost:27017`)

### Execute Migration

```bash
# From backend directory
go run ./cmd/migrate/create_order_items_collection.go
```

### Expected Output

```
Connecting to MongoDB: mongodb://localhost:27017
✅ MongoDB connected successfully
📝 Creating order_items collection...
✅ order_items collection created
📝 Creating indexes for order_items collection...
✅ Indexes created successfully

📊 Collection Information:
Collection: order_items
Indexes:
  - idx_order_id: order_id
  - idx_menu_item_id: menu_item_id
  - idx_cost_status: cost_status
  - idx_cost_calculated_at: cost_calculated_at
  - idx_order_menu_item: order_id + menu_item_id (compound)

🎉 Migration completed successfully!

ℹ️  Note: This collection will be populated when shifts are closed.
   Order items will have accounting_cost calculated at shift closure time.
```

## Data Population

The `order_items` collection will be populated automatically when:

1. **Shift Closure**: When a manager closes a shift, the system will:
   - Calculate accounting_cost for all order items in that shift
   - Use current ingredient cost_per_unit values at closure time
   - Store items in order_items collection with cost_status = "FINAL"

2. **Historical Data**: For orders before this feature was implemented:
   - Accounting_cost will be calculated using current ingredient costs
   - Items will be marked with cost_status = "ESTIMATED"
   - A note will indicate this is backfilled data

## Verification

After running the migration, verify the collection and indexes:

```bash
# Connect to MongoDB
mongosh cafe_pos

# Check collection exists
db.getCollectionNames()

# Check indexes
db.order_items.getIndexes()

# Expected output should show 6 indexes (including default _id index)
```

## Rollback

If you need to remove the collection:

```bash
mongosh cafe_pos
db.order_items.drop()
```

**Warning**: This will delete all accounting cost data. Only do this in development.

## Related Files

- **Domain Model**: `backend/domain/order/order_item.go`
- **Repository**: `backend/infrastructure/mongodb/order_item_repository.go`
- **Migration Script**: `backend/cmd/migrate/create_order_items_collection.go`

## Requirements Satisfied

This migration satisfies the following requirements from the spec:

- **Requirement 5.1**: Order items stored separately for cost tracking
- **Requirement 5.2**: Accounting cost calculated at shift closure
- **Requirement 5.3**: Cost status tracking (FINAL, ESTIMATED, INCOMPLETE)

## Next Steps

After running this migration:

1. Implement the Cost Calculator Service (Task 2.1)
2. Implement shift closure cost calculation (Task 2.4)
3. Update shift closure endpoint to populate order_items (Task 9.1)
