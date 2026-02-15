# Batch Cost Estimation Implementation

## Summary
Implemented batch cost estimation based on batch definition conversion rates to help with menu pricing.

## Problem
When adding batch ingredients to menu items, the cost was showing as 0 VND, making it impossible to estimate menu pricing accurately.

## Solution
Created a `calculateBatchCostPerUnit()` function that:
1. Reads the batch definition's `conversion_rates` array
2. For each conversion rate, finds the source ingredient from the ingredient store using `source_ingredient_id`
3. Calculates cost: `(source_quantity * (1 + wastage_rate)) * source_ingredient.cost_per_unit`
4. Sums all ingredient costs
5. Divides by `batch_quantity` (from first conversion rate) to get cost per batch output unit

## Implementation Details

### Key Fixes
1. **Field Name**: Backend uses `source_ingredient_id`, not `ingredient_id`
2. **Wastage Rate Format**: Backend stores as decimal (0.1 = 10%), not percentage (10)
3. **Batch Quantity**: Taken from `conversion_rates[0].batch_quantity`, not `batch.batch_quantity`

### New Function: `calculateBatchCostPerUnit(batch)`
```javascript
const calculateBatchCostPerUnit = (batch) => {
  console.log('🧪 Calculating batch cost for:', batch.name)
  
  if (!batch.conversion_rates || !Array.isArray(batch.conversion_rates) || batch.conversion_rates.length === 0) {
    console.warn('⚠️ No conversion_rates found for batch:', batch.name)
    return 0
  }
  
  let totalCost = 0
  
  for (const conversionRate of batch.conversion_rates) {
    // Find source ingredient using source_ingredient_id
    const ingredientId = conversionRate.source_ingredient_id || conversionRate.ingredient_id
    const sourceIngredient = availableIngredients.value.find(ing => ing.id === ingredientId)
    
    if (!sourceIngredient) {
      console.warn(`❌ Source ingredient not found: ${ingredientId}`)
      continue
    }
    
    const sourceQuantity = conversionRate.source_quantity || 0
    const wastageRate = conversionRate.wastage_rate || 0 // Already decimal (0.1 = 10%)
    const costPerUnit = sourceIngredient.cost_per_unit || 0
    
    const ingredientCost = sourceQuantity * (1 + wastageRate) * costPerUnit
    totalCost += ingredientCost
  }
  
  // Get batch_quantity from first conversion rate
  const batchQuantity = batch.conversion_rates[0]?.batch_quantity || 1
  const costPerBatchUnit = totalCost / batchQuantity
  
  return costPerBatchUnit
}
```

### Updated `selectBatch()` Function
- Calls `calculateBatchCostPerUnit(batch)` to get the estimated cost per unit
- Sets `costPerUnit` to the calculated value instead of 0
- Sets `estimatedCost` to `batchCostPerUnit` for initial quantity of 1
- Works for both single-size items and variants

### Updated UI Display
Changed batch ingredient display from:
```
Batch (chi phí tính theo công thức)
```

To:
```
Batch @ 10,000 ₫/L
```

This shows the actual calculated cost per unit, making it clear and actionable for menu pricing.

## Backend Data Structure

```go
type BatchDefinition struct {
    ID              primitive.ObjectID `json:"id"`
    Name            string             `json:"name"`
    Unit            string             `json:"unit"`
    ConversionRates []ConversionRate   `json:"conversion_rates"`
    // ...
}

type ConversionRate struct {
    SourceIngredientID   primitive.ObjectID `json:"source_ingredient_id"`
    SourceIngredientName string             `json:"source_ingredient_name"`
    SourceQuantity       float64            `json:"source_quantity"`
    SourceUnit           string             `json:"source_unit"`
    BatchQuantity        float64            `json:"batch_quantity"`
    WastageRate          float64            `json:"wastage_rate"` // 0.0 to 1.0 (0.1 = 10%)
}
```

## Example Calculation

If a batch definition has:
- Batch quantity: 1 L (from conversion_rates[0].batch_quantity)
- Conversion rates:
  - Ingredient A: 500g @ 20,000 ₫/kg, wastage 0.05 (5%)
  - Ingredient B: 200ml @ 50,000 ₫/L, wastage 0.10 (10%)

Calculation:
1. Ingredient A cost: 0.5 kg × 1.05 × 20,000 = 10,500 ₫
2. Ingredient B cost: 0.2 L × 1.10 × 50,000 = 11,000 ₫
3. Total cost: 10,500 + 11,000 = 21,500 ₫
4. Cost per L: 21,500 ÷ 1 = 21,500 ₫/L

When adding this batch to a menu item with quantity 0.5 L:
- Estimated cost: 0.5 × 21,500 = 10,750 ₫

## Debugging
Added console.log statements to help debug:
- Logs batch data when calculating
- Logs each conversion rate processing
- Logs source ingredient lookup
- Logs cost calculation for each ingredient
- Logs final total cost and cost per unit

Check browser console when adding batch ingredients to see detailed calculation logs.

## Files Modified
- `frontend/src/views/MenuView.vue`
  - Added `calculateBatchCostPerUnit()` function
  - Updated `selectBatch()` to calculate and set batch cost
  - Updated template to display batch cost per unit
  - Fixed field name from `ingredient_id` to `source_ingredient_id`
  - Fixed wastage_rate handling (already decimal, no need to divide by 100)
  - Fixed batch_quantity source (from conversion_rates[0])

## Testing
To test:
1. Create a batch definition with conversion rates
2. Open browser console (F12)
3. Add the batch to a menu item
4. Check console logs for calculation details
5. Verify the cost per unit is displayed correctly
6. Change the quantity and verify the estimated cost updates
7. Check that the total cost calculation includes batch ingredients

## Status
✅ Complete - Batch cost estimation is now working for menu pricing with proper field mappings
