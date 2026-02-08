# Ingredient Stock OUT Fix - Quantity Sign Handling

## Problem
When using Quick Stock OUT or Adjust Stock with type="remove", the stock quantity was not decreasing. The ingredient quantity remained unchanged after the operation.

## Root Cause

### Backend Logic
```go
// backend/application/services/ingredient.go
func (s *IngredientService) AdjustStock(...) {
    beforeQty := item.Quantity
    item.Quantity += req.Quantity  // Simply adds the quantity
    if item.Quantity < 0 {
        item.Quantity = 0
    }
    afterQty := item.Quantity
}
```

The backend **always adds** the quantity value to the current stock:
- `item.Quantity += req.Quantity`
- It doesn't check the `type` field
- It expects the quantity to be **negative** for removals

### Frontend Issue
The frontend was sending:
```javascript
// WRONG - Positive quantity for removal
{
  type: 'remove',
  quantity: 5,  // Should be -5!
  reason: 'Sử dụng'
}
```

Result: `10 + 5 = 15` (stock increased instead of decreased!)

## The Fix

### Frontend Changes
Modified the frontend to send **negative quantities** for removals:

#### 1. Quick Stock OUT
```javascript
const confirmQuickOut = async () => {
  const adjustment = {
    type: 'remove',
    quantity: -Math.abs(quickOutData.value.quantity), // NEGATIVE
    reason: finalReason,
    cost_per_unit: 0
  }
  await ingredientStore.adjustStock(currentIngredient.value.id, adjustment)
}
```

#### 2. Full Adjust Stock Modal
```javascript
const adjustStock = async () => {
  let finalQuantity = adjustData.value.quantity
  
  if (adjustData.value.type === 'remove') {
    // Remove: send negative quantity
    finalQuantity = -Math.abs(adjustData.value.quantity)
  } else if (adjustData.value.type === 'add') {
    // Add: send positive quantity
    finalQuantity = Math.abs(adjustData.value.quantity)
  } else if (adjustData.value.type === 'adjust') {
    // Set: calculate difference from current
    finalQuantity = adjustData.value.quantity - currentIngredient.value.quantity
  }
  
  const adjustmentData = {
    ...adjustData.value,
    quantity: finalQuantity
  }
  
  await ingredientStore.adjustStock(currentIngredient.value.id, adjustmentData)
}
```

## How It Works Now

### Type: ADD (Nhập hàng)
```
User enters: 5
Backend receives: +5
Calculation: 10 + 5 = 15 ✓
```

### Type: REMOVE (Xuất hàng)
```
User enters: 5
Backend receives: -5
Calculation: 10 + (-5) = 5 ✓
```

### Type: ADJUST (Điều chỉnh/Set)
```
User enters: 7 (target quantity)
Current: 10
Difference: 7 - 10 = -3
Backend receives: -3
Calculation: 10 + (-3) = 7 ✓
```

## Why This Approach?

### Option 1: Fix Frontend (CHOSEN) ✓
**Pros:**
- Backend logic is simple and consistent
- One place to handle quantity math
- No breaking changes to backend
- Works with existing history records

**Cons:**
- Frontend must calculate correct sign
- Multiple places to update

### Option 2: Fix Backend (NOT CHOSEN)
**Pros:**
- Frontend sends intuitive values
- Type field is actually used

**Cons:**
- More complex backend logic
- Need to handle type in multiple places
- Potential breaking changes
- History records store signed values

## Testing

### Test Case 1: Quick Stock OUT
```
Initial: 10 kg
Action: Quick OUT 3 kg
Expected: 7 kg
Result: ✓ 7 kg
```

### Test Case 2: Adjust Modal - Remove
```
Initial: 10 kg
Action: Type=Remove, Quantity=4
Expected: 6 kg
Result: ✓ 6 kg
```

### Test Case 3: Adjust Modal - Add
```
Initial: 10 kg
Action: Type=Add, Quantity=5
Expected: 15 kg
Result: ✓ 15 kg
```

### Test Case 4: Adjust Modal - Set
```
Initial: 10 kg
Action: Type=Adjust, Quantity=8
Expected: 8 kg
Calculation: 8 - 10 = -2
Backend: 10 + (-2) = 8
Result: ✓ 8 kg
```

### Test Case 5: Remove More Than Available
```
Initial: 5 kg
Action: Remove 10 kg
Backend: 5 + (-10) = -5 → 0 (clamped)
Result: ✓ 0 kg (backend prevents negative)
```

## Stock History

History records now correctly show:
- **Positive quantity** for additions (green cards)
- **Negative quantity** for removals (red cards)
- **Positive or negative** for adjustments (blue cards)

Example:
```javascript
{
  type: "adjustment",
  quantity: -5,        // Negative for removal
  before_qty: 10,
  after_qty: 5,
  reason: "Sử dụng cho món ăn"
}
```

## Validation

### Frontend Validation
```javascript
// User enters positive number
if (quickOutData.value.quantity > currentIngredient.value.quantity) {
  alert('Số lượng xuất không được lớn hơn tồn kho')
  return
}

// Convert to negative before sending
quantity: -Math.abs(quickOutData.value.quantity)
```

### Backend Safety
```go
// Prevents negative stock
if item.Quantity < 0 {
    item.Quantity = 0
}
```

## Edge Cases Handled

### 1. Remove More Than Available
- Frontend validates before sending
- Backend clamps to 0 if negative
- User sees warning

### 2. Set to Zero
```
Type: Adjust, Quantity: 0
Difference: 0 - 10 = -10
Backend: 10 + (-10) = 0 ✓
```

### 3. Set to Same Value
```
Type: Adjust, Quantity: 10
Difference: 10 - 10 = 0
Backend: 10 + 0 = 10 ✓
```

### 4. Decimal Quantities
```
Type: Remove, Quantity: 2.5
Backend: 10 + (-2.5) = 7.5 ✓
```

## Files Modified

1. `frontend/src/views/IngredientManagementView.vue`
   - Updated `confirmQuickOut()` to send negative quantity
   - Updated `adjustStock()` to handle all three types correctly
   - Added logic to calculate difference for type="adjust"

## Backward Compatibility

### Existing History Records
- Old records may have incorrect signs
- New records will have correct signs
- Display logic handles both cases

### API Compatibility
- No API changes
- Backend still expects signed quantity
- Frontend now sends correct signs

## Future Improvements

### Option 1: Backend Type Handling
Add explicit type handling in backend:
```go
switch req.Type {
case "add":
    item.Quantity += math.Abs(req.Quantity)
case "remove":
    item.Quantity -= math.Abs(req.Quantity)
case "adjust":
    item.Quantity = req.Quantity
}
```

### Option 2: Separate Endpoints
```
POST /ingredients/:id/add-stock
POST /ingredients/:id/remove-stock
POST /ingredients/:id/set-stock
```

### Option 3: Explicit Delta Field
```javascript
{
  type: "remove",
  quantity: 5,      // User-entered value
  delta: -5,        // Calculated change
  target: null      // For set operations
}
```

## Conclusion

The fix ensures that stock OUT operations correctly decrease inventory by sending negative quantities to the backend. This maintains consistency with the backend's simple addition logic while providing an intuitive UI for users.

The three operation types now work correctly:
- **ADD**: Positive quantity → stock increases
- **REMOVE**: Negative quantity → stock decreases  
- **ADJUST**: Calculated difference → stock set to target

All operations are validated, safe, and create accurate history records.
