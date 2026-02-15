# Batch Collections MongoDB Indexes

This document describes the MongoDB indexes created for the batch management collections.

## Collections

### 1. batch_definitions

Stores batch definition templates that define how batches are prepared.

**Indexes:**

| Index Name | Fields | Type | Purpose |
|------------|--------|------|---------|
| name_idx | `name: 1` | Single | Fast lookup by batch name for search functionality |
| created_at_idx | `created_at: -1` | Single | Sort batch definitions by creation date (newest first) |

**Query Patterns:**
- Search by name (case-insensitive): Uses `name_idx`
- List all definitions sorted by date: Uses `created_at_idx`
- Pagination queries: Uses `created_at_idx` for sorting

### 2. batch_records

Stores individual batch instances that have been prepared.

**Indexes:**

| Index Name | Fields | Type | Purpose |
|------------|--------|------|---------|
| batch_def_expires_idx | `batch_definition_id: 1, expires_at: 1` | Compound | FIFO queries - find available batches for a definition sorted by expiry |
| status_expires_idx | `status: 1, expires_at: 1` | Compound | Find batches by status and expiry (e.g., available batches expiring soon) |
| expires_at_idx | `expires_at: 1` | Single | Background job to mark expired batches |
| prepared_at_idx | `prepared_at: -1` | Single | List batches by preparation date (newest first) |
| prepared_by_date_idx | `prepared_by: 1, prepared_at: -1` | Compound | Filter batches by preparer and sort by date |

**Query Patterns:**
- **FIFO batch selection**: Uses `batch_def_expires_idx` to find available batches for a definition, sorted by expiry (oldest first)
- **Alert generation**: Uses `status_expires_idx` to find expiring batches
- **Expiry checking**: Uses `expires_at_idx` for background job that marks expired batches
- **User activity**: Uses `prepared_by_date_idx` to show batches prepared by a specific user
- **Recent batches**: Uses `prepared_at_idx` for dashboard views

### 3. batch_usage_logs

Records when batches are used in orders for audit trail and reporting.

**Indexes:**

| Index Name | Fields | Type | Purpose |
|------------|--------|------|---------|
| batch_record_used_idx | `batch_record_id: 1, used_at: -1` | Compound | View usage history for a specific batch |
| order_id_idx | `order_id: 1` | Single | Find all batch usages for an order |
| menu_item_used_idx | `menu_item_id: 1, used_at: -1` | Compound | Analyze batch usage by menu item over time |
| used_at_idx | `used_at: -1` | Single | Recent usage history and time-based reports |

**Query Patterns:**
- **Batch usage history**: Uses `batch_record_used_idx` to show when and how a batch was used
- **Order details**: Uses `order_id_idx` to show all batches used in an order
- **Menu item analysis**: Uses `menu_item_used_idx` for reports on which menu items use which batches
- **Usage reports**: Uses `used_at_idx` for time-based usage reports

## Index Design Rationale

### FIFO (First In First Out) Support

The `batch_def_expires_idx` compound index on `(batch_definition_id, expires_at)` is critical for FIFO batch usage:

```javascript
// Query to get available batches for FIFO usage
db.batch_records.find({
  batch_definition_id: ObjectId("..."),
  status: "available",
  expires_at: { $gt: new Date() },
  quantity_remaining: { $gt: 0 }
}).sort({ expires_at: 1 })
```

This ensures batches with the earliest expiry date are used first, minimizing waste.

### Performance Considerations

1. **Compound Indexes**: Used for queries that filter and sort on multiple fields
2. **Sort Direction**: Indexes match the sort direction used in queries (-1 for descending)
3. **Selectivity**: Most selective fields are placed first in compound indexes
4. **Index Size**: Kept minimal to balance query performance with write performance

### Background Jobs

The following indexes support background jobs:

- `expires_at_idx`: Used by the job that marks expired batches (runs every hour)
- `status_expires_idx`: Used by the alert checking job (runs every 5 minutes)

## Migration

To create these collections and indexes, run:

```bash
# From the backend directory
./create_batch_collections

# Or directly with Go
go run cmd/migrate/create_batch_collections.go
```

The migration script:
1. Creates the three collections if they don't exist
2. Creates all indexes
3. Prints schema information
4. Is idempotent (safe to run multiple times)

## Monitoring

Monitor index usage with:

```javascript
// Check index usage statistics
db.batch_records.aggregate([{ $indexStats: {} }])

// Check slow queries
db.setProfilingLevel(1, { slowms: 100 })
db.system.profile.find().sort({ ts: -1 }).limit(10)
```

## Index Maintenance

MongoDB automatically maintains indexes, but consider:

1. **Rebuilding indexes** after bulk operations:
   ```javascript
   db.batch_records.reIndex()
   ```

2. **Monitoring index size**:
   ```javascript
   db.batch_records.stats().indexSizes
   ```

3. **Dropping unused indexes** if query patterns change

## Future Optimizations

Consider adding these indexes if needed:

1. **Text search**: If full-text search on batch names is required
   ```javascript
   db.batch_definitions.createIndex({ name: "text" })
   ```

2. **Geospatial**: If batch location tracking is added
   ```javascript
   db.batch_records.createIndex({ location: "2dsphere" })
   ```

3. **TTL index**: If automatic deletion of old logs is needed
   ```javascript
   db.batch_usage_logs.createIndex(
     { used_at: 1 },
     { expireAfterSeconds: 31536000 } // 1 year
   )
   ```
