# Ingredient Stock History - ID Assignment Fix

## Problem
Stock history records were not being created consistently when creating new ingredients. Sometimes they appeared, sometimes they didn't.

## Root Cause
When creating an ingredient in MongoDB:
1. `InsertOne()` is called on the ingredient
2. MongoDB generates and assigns an `_id` to the document
3. However, the Go struct's `ID` field was NOT being updated
4. When creating stock history immediately after, `item.ID` was still empty/nil
5. Stock history was created with an invalid `ingredient_id`
6. When fetching history by `ingredient_id`, no records were found

## The Bug

**Before (broken):**
```go
func (r *IngredientRepository) Create(ctx context.Context, item *ingredient.Ingredient) error {
    item.CreatedAt = time.Now()
    item.UpdatedAt = time.Now()
    _, err := r.collection.InsertOne(ctx, item)  // ID not captured!
    return err
}
```

**In service:**
```go
err := s.ingredientRepo.Create(ctx, item)
// At this point, item.ID is still empty!

history := &ingredient.StockHistory{
    IngredientID: item.ID,  // This is empty/nil!
    // ...
}
s.stockHistoryRepo.Create(ctx, history)
```

## The Fix

**After (fixed):**
```go
func (r *IngredientRepository) Create(ctx context.Context, item *ingredient.Ingredient) error {
    item.CreatedAt = time.Now()
    item.UpdatedAt = time.Now()
    result, err := r.collection.InsertOne(ctx, item)
    if err != nil {
        return err
    }
    // Capture the inserted ID and update the struct
    if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
        item.ID = oid
    }
    return nil
}
```

Now when the service creates stock history, `item.ID` has the correct value.

## Why This Happened Intermittently

The bug appeared intermittent because:
1. If you created an ingredient and immediately viewed history → No records (ID was empty)
2. If you refreshed the page and viewed history → Still no records (wrong ingredient_id in DB)
3. If you adjusted stock → New history record created with correct ID
4. Now viewing history shows only the adjustment, not the initial creation

So it seemed "random" but was actually consistent - initial history was always broken.

## Impact

### Before Fix
- Initial stock history record had invalid `ingredient_id`
- Fetching history by ingredient ID returned empty array
- Users saw no history for newly created ingredients
- Only after first stock adjustment would history start appearing

### After Fix
- Initial stock history record has correct `ingredient_id`
- Fetching history returns the initial creation record
- Users see complete history from day one
- Consistent behavior every time

## Testing

### Manual Test
1. Create new ingredient with quantity > 0
2. Immediately click "Lịch sử" button
3. Should see initial purchase record
4. Verify all fields are correct

### Automated Test
Run the test script:
```bash
./test-ingredient-history-creation.sh
```

This tests:
- Creating ingredient with quantity > 0 → History created
- Creating ingredient with quantity = 0 → No history created
- Verifying all history fields are correct
- Cleanup

## Files Modified

1. `backend/infrastructure/mongodb/ingredient_repository.go`
   - Updated `Create()` method to capture and set inserted ID
   - Ensures `item.ID` is populated after insertion

## Related Code

The fix ensures this flow works correctly:

```go
// 1. Create ingredient (ID is now set)
err := s.ingredientRepo.Create(ctx, item)

// 2. Create stock history with correct ID
if req.Quantity > 0 {
    history := &ingredient.StockHistory{
        IngredientID: item.ID,  // Now has correct value!
        // ...
    }
    s.stockHistoryRepo.Create(ctx, history)
}
```

## MongoDB Behavior

Important to understand:
- MongoDB's `InsertOne()` returns `InsertResult` with `InsertedID`
- The inserted document in DB has the `_id` field
- But the Go struct passed to `InsertOne()` is NOT automatically updated
- You must manually extract and set the ID from the result

This is standard MongoDB Go driver behavior, not a bug.

## Prevention

To prevent similar issues in the future:
1. Always capture `InsertResult` from `InsertOne()`
2. Always extract and set the ID on the struct
3. Test immediately after creation, don't wait for other operations
4. Use automated tests to verify ID assignment

## Verification Checklist

After deploying this fix:
- [ ] Restart backend server
- [ ] Create new ingredient with quantity > 0
- [ ] Immediately view history
- [ ] Verify initial record appears
- [ ] Verify all fields are correct (quantity, price, reason)
- [ ] Create ingredient with quantity = 0
- [ ] Verify no history record created
- [ ] Adjust stock on both ingredients
- [ ] Verify history shows all records in correct order

## Database Cleanup (Optional)

Existing ingredients created before this fix may have orphaned stock history records with invalid ingredient_id. To clean up:

```javascript
// MongoDB shell
db.stock_history.find({ ingredient_id: { $exists: false } })
db.stock_history.find({ ingredient_id: null })
db.stock_history.find({ ingredient_id: ObjectId("000000000000000000000000") })

// Delete orphaned records
db.stock_history.deleteMany({ ingredient_id: { $exists: false } })
db.stock_history.deleteMany({ ingredient_id: null })
db.stock_history.deleteMany({ ingredient_id: ObjectId("000000000000000000000000") })
```

## Lessons Learned

1. **Test immediately**: Don't assume operations work, test right after
2. **Understand driver behavior**: Know how your database driver handles IDs
3. **Capture return values**: Don't ignore return values from insert operations
4. **Consistent patterns**: Apply same pattern across all repositories
5. **Automated tests**: Catch these issues before production
