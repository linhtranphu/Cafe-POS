# Ingredient Weighted Average Cost Implementation

## Overview
Implemented weighted average cost tracking for ingredients to handle price changes over time. Each stock adjustment can have its own price, and the system automatically calculates the weighted average.

## Problem Solved
**Before:**
```
Month 1: Buy 5kg coffee @ 200k/kg
Month 2: Buy 3kg coffee @ 250k/kg (price increased!)
→ System only stored one price (200k)
→ Expense tracking was incorrect
```

**After:**
```
Month 1: Buy 5kg @ 200k/kg
  - Stock: 5kg
  - Avg cost: 200k/kg
  - Expense: 1,000k ✓

Month 2: Buy 3kg @ 250k/kg
  - Stock: 8kg
  - Avg cost: 218.75k/kg (weighted average)
  - Expense: 750k ✓
  
Total expenses: 1,750k ✓
```

## Implementation

### 1. Backend Changes

#### Domain Models

**StockHistory** - Added price fields:
```go
type StockHistory struct {
    // ... existing fields
    CostPerUnit float64 `bson:"cost_per_unit" json:"cost_per_unit"` // Price at transaction time
    TotalCost   float64 `bson:"total_cost" json:"total_cost"`       // Total cost for transaction
}
```

**StockAdjustmentRequest** - Added optional price:
```go
type StockAdjustmentRequest struct {
    Quantity    float64 `json:"quantity" binding:"required"`
    Reason      string  `json:"reason" binding:"required"`
    CostPerUnit float64 `json:"cost_per_unit"` // Optional: price for this purchase
    UserID      string  `json:"user_id"`
    Username    string  `json:"username"`
}
```

#### Service Logic

**Weighted Average Calculation:**
```go
func (s *IngredientService) AdjustStock(ctx, id, req) {
    item, _ := s.ingredientRepo.FindByID(ctx, id)
    
    beforeQty := item.Quantity
    afterQty := beforeQty + req.Quantity
    
    // Determine cost for this transaction
    costPerUnit := req.CostPerUnit
    if costPerUnit <= 0 {
        costPerUnit = item.CostPerUnit // Use current if not provided
    }
    
    // Calculate weighted average for stock IN
    if req.Quantity > 0 && costPerUnit > 0 && afterQty > 0 {
        oldValue := beforeQty * item.CostPerUnit
        newValue := req.Quantity * costPerUnit
        item.CostPerUnit = (oldValue + newValue) / afterQty
    }
    
    // Save history with transaction price
    history := &ingredient.StockHistory{
        // ... other fields
        CostPerUnit: costPerUnit,
        TotalCost:   req.Quantity * costPerUnit,
    }
    
    // Track expense with actual transaction price
    tempItem := *item
    tempItem.CostPerUnit = costPerUnit
    s.autoExpenseService.TrackIngredientPurchase(ctx, &tempItem, req.Quantity, req.Username)
}
```

### 2. Frontend Changes

#### Adjust Stock Modal

**Added price input field:**
```vue
<div v-if="adjustData.type === 'add'">
  <label>Giá lần này (nếu khác giá hiện tại)</label>
  <input v-model.number="adjustData.cost_per_unit" 
    :placeholder="`Giá hiện tại: ${formatCurrency(currentIngredient.cost_per_unit)}`" />
  <p class="text-xs">
    💡 Để trống nếu giá không đổi
  </p>
</div>
```

**Dynamic expense calculation:**
```javascript
const effectiveAdjustPrice = computed(() => {
  // Use provided price, or fall back to current price
  return adjustData.value.cost_per_unit > 0 
    ? adjustData.value.cost_per_unit 
    : currentIngredient.value?.cost_per_unit || 0
})

const adjustExpenseAmount = computed(() => {
  if (adjustData.value.type !== 'add') return 0
  return adjustData.value.quantity * effectiveAdjustPrice.value
})
```

#### Stock History Display

**Shows price information:**
```vue
<div v-if="record.cost_per_unit > 0" class="bg-green-50 rounded-lg p-2">
  <div>Đơn giá: {{ formatCurrency(record.cost_per_unit) }}/{{ unit }}</div>
  <div>Tổng chi phí: {{ formatCurrency(record.total_cost) }}</div>
</div>
```

## Usage Examples

### Example 1: Price Unchanged

```
Adjust Stock:
- Quantity: +3kg
- Price: (leave empty)
- Reason: "Nhập thêm hàng"

Result:
- Uses current price: 200k/kg
- Expense: 600k
- Avg cost: unchanged (200k/kg)
```

### Example 2: Price Changed

```
Current:
- Stock: 5kg
- Avg cost: 200k/kg

Adjust Stock:
- Quantity: +3kg
- Price: 250k/kg ← Enter new price
- Reason: "Nhập thêm, giá tăng"

Result:
- Stock: 8kg
- Avg cost: 218.75k/kg
  = (5kg × 200k + 3kg × 250k) / 8kg
  = (1,000k + 750k) / 8kg
  = 218.75k/kg
- Expense: 750k ✓
```

### Example 3: Multiple Purchases

```
Purchase 1:
- 5kg @ 200k/kg
- Stock: 5kg, Avg: 200k/kg
- Expense: 1,000k

Purchase 2:
- 3kg @ 250k/kg
- Stock: 8kg, Avg: 218.75k/kg
- Expense: 750k

Purchase 3:
- 2kg @ 220k/kg
- Stock: 10kg, Avg: 219k/kg
  = (8kg × 218.75k + 2kg × 220k) / 10kg
  = (1,750k + 440k) / 10kg
  = 219k/kg
- Expense: 440k

Total expenses: 2,190k ✓
```

## Weighted Average Formula

```
New Avg Cost = (Old Qty × Old Cost + New Qty × New Cost) / Total Qty

Example:
Old: 5kg @ 200k/kg = 1,000k
New: 3kg @ 250k/kg = 750k
Total: 8kg

Avg = (1,000k + 750k) / 8kg = 218.75k/kg
```

## Benefits

1. **Accurate Costing**: Tracks actual purchase prices
2. **Historical Data**: Full price history in stock records
3. **Flexible**: Can use current price or enter new price
4. **Automatic**: Weighted average calculated automatically
5. **Expense Tracking**: Expenses match actual costs
6. **Accounting Compliant**: Follows weighted average cost method

## Stock History Record Example

```json
{
  "id": "...",
  "ingredient_id": "...",
  "type": "adjustment",
  "quantity": 3,
  "before_qty": 5,
  "after_qty": 8,
  "reason": "Nhập thêm hàng, giá tăng",
  "cost_per_unit": 250000,
  "total_cost": 750000,
  "username": "Manager",
  "created_at": "2026-02-07T10:00:00Z"
}
```

## UI Flow

### Adjust Stock with New Price

1. **Click "📦 Điều chỉnh"** on ingredient
2. **Select type:** "Nhập thêm"
3. **Enter quantity:** 3kg
4. **Enter new price:** 250,000 (or leave empty to use current)
5. **See calculation:**
   - Green card shows: "Chi phí: 750,000₫"
   - Formula: "= 3kg × 250,000₫/kg"
6. **Enter reason:** "Nhập thêm, giá tăng"
7. **Submit**
   - Stock updated: 8kg
   - Avg cost updated: 218.75k/kg
   - Expense created: 750k
   - History saved with price

### View History

1. **Click "📊 Lịch sử"** on ingredient
2. **See all transactions** with:
   - Date and user
   - Quantity change
   - Before/after stock
   - **Price at that time** (green box)
   - **Total cost** for that transaction

## Testing Checklist

- [ ] Adjust stock without entering price (uses current)
- [ ] Adjust stock with new price (calculates weighted avg)
- [ ] Multiple adjustments with different prices
- [ ] Verify expense amounts are correct
- [ ] Check stock history shows prices
- [ ] Verify weighted average calculation
- [ ] Test with decimal quantities
- [ ] Test with large price differences
- [ ] Verify stock OUT doesn't affect avg cost

## Files Modified

**Backend:**
- `backend/domain/ingredient/stock_history.go`
- `backend/domain/ingredient/ingredient.go`
- `backend/application/services/ingredient.go`

**Frontend:**
- `frontend/src/views/IngredientManagementView.vue`

---
**Status**: ✅ Complete
**Date**: 2026-02-07
**Method**: Weighted Average Cost
**Accounting Standard**: Compliant
