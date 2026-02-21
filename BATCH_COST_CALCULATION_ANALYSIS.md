# Batch Cost Calculation Analysis

## Current Logic in BatchRecordForm.vue

```javascript
const expectedCost = computed(() => {
  if (!selectedDefinition.value || batchCount.value <= 0) return 0
  
  let total = 0
  const def = selectedDefinition.value
  const quantity = totalOutput.value  // Total output quantity
  
  for (const rate of def.conversion_rates || []) {
    const ingredient = ingredients.value.find(i => i.id === rate.source_ingredient_id)
    if (ingredient && ingredient.cost_per_unit) {
      // Step 1: Calculate ratio
      const ratio = quantity / rate.batch_quantity
      
      // Step 2: Calculate base quantity needed
      const baseQuantity = rate.source_quantity * ratio
      
      // Step 3: Apply wastage
      const wastageMultiplier = 1 + (rate.wastage_rate || 0)
      const totalQuantity = baseQuantity * wastageMultiplier
      
      // Step 4: Convert to stock unit
      const conversionRate = getConversionRate(ingredient.unit, rate.source_unit)
      const quantityInStockUnit = totalQuantity * conversionRate
      
      // Step 5: Calculate cost
      total += quantityInStockUnit * ingredient.cost_per_unit
    }
  }
  
  return total
})
```

## Example Calculation

### Scenario
- **Batch Definition**: "Sữa tươi pha sẵn"
  - Output: 1L per batch
  - Ingredient: "Sữa tươi nguyên chất"
    - Recipe: 800ml (source_quantity = 800, source_unit = "ml")
    - Stock: 1L @ 50,000 VNĐ/L (ingredient.unit = "L", cost_per_unit = 50,000)
    - Wastage: 5% (wastage_rate = 0.05)

- **User Input**: Tạo 3 batch (batchCount = 3)
  - Total output: 3L (totalOutput = 3)

### Step-by-Step Calculation

#### Step 1: Calculate ratio
```
ratio = totalOutput / batch_quantity
ratio = 3L / 1L = 3
```

#### Step 2: Calculate base quantity
```
baseQuantity = source_quantity * ratio
baseQuantity = 800ml * 3 = 2400ml
```

#### Step 3: Apply wastage
```
wastageMultiplier = 1 + wastage_rate
wastageMultiplier = 1 + 0.05 = 1.05

totalQuantity = baseQuantity * wastageMultiplier
totalQuantity = 2400ml * 1.05 = 2520ml
```

#### Step 4: Convert to stock unit
```
conversionRate = getConversionRate("L", "ml")
conversionRate = 0.001

quantityInStockUnit = totalQuantity * conversionRate
quantityInStockUnit = 2520ml * 0.001 = 2.52L
```

#### Step 5: Calculate cost
```
cost = quantityInStockUnit * cost_per_unit
cost = 2.52L * 50,000 VNĐ/L = 126,000 VNĐ
```

### Expected Result
- **Chi phí dự kiến**: 126,000 VNĐ
- **Chi phí per batch**: 126,000 / 3 = 42,000 VNĐ/batch

## Potential Issues

### Issue 1: totalOutput calculation
Check if `totalOutput` is calculated correctly:
```javascript
const totalOutput = computed(() => {
  if (!selectedDefinition.value || batchCount.value <= 0) return 0
  return selectedDefinition.value.output_quantity * batchCount.value
})
```

This looks correct.

### Issue 2: Conversion rate direction
The current code uses:
```javascript
const conversionRate = getConversionRate(ingredient.unit, rate.source_unit)
```

This is CORRECT because:
- `getConversionRate(stockUnit, recipeUnit)` returns rate to convert recipe→stock
- `getConversionRate("L", "ml")` = 0.001
- 200ml × 0.001 = 0.2L ✅

### Issue 3: Multiple ingredients
If batch has multiple ingredients, the loop should sum all costs correctly.

### Issue 4: Unit mismatch
If `rate.source_unit` doesn't match any conversion table, `getConversionRate` returns 1.0, which might be wrong.

## Debugging Steps

1. **Log values in console**:
```javascript
console.log('Batch Definition:', selectedDefinition.value)
console.log('Batch Count:', batchCount.value)
console.log('Total Output:', totalOutput.value)
console.log('Conversion Rates:', def.conversion_rates)

for (const rate of def.conversion_rates || []) {
  console.log('Processing rate:', rate)
  console.log('Ingredient:', ingredient)
  console.log('Ratio:', ratio)
  console.log('Base Quantity:', baseQuantity)
  console.log('Total Quantity (with wastage):', totalQuantity)
  console.log('Conversion Rate:', conversionRate)
  console.log('Quantity in Stock Unit:', quantityInStockUnit)
  console.log('Cost:', quantityInStockUnit * ingredient.cost_per_unit)
}
```

2. **Check ingredient data**:
- Verify `ingredient.cost_per_unit` is correct
- Verify `ingredient.unit` matches expected stock unit

3. **Check conversion rate data**:
- Verify `rate.source_unit` is correct
- Verify `rate.source_quantity` is correct
- Verify `rate.batch_quantity` matches `output_quantity`

4. **Manual calculation**:
- Calculate expected cost manually
- Compare with displayed cost
- Identify which step is wrong

## Possible Fix

If the issue is that conversion is backwards, we might need to use the inverse:

```javascript
// Instead of:
const conversionRate = getConversionRate(ingredient.unit, rate.source_unit)
const quantityInStockUnit = totalQuantity * conversionRate

// Try:
const conversionRate = getConversionRate(rate.source_unit, ingredient.unit)
const quantityInStockUnit = totalQuantity / conversionRate
```

But this would only be needed if `getConversionRate` returns the inverse of what we expect.

## Testing

Create a test batch with known values:
1. Ingredient: "Sữa" - 1L @ 50,000 VNĐ
2. Batch: Uses 800ml sữa → produces 1L output
3. Create 1 batch
4. Expected cost: 800ml × 0.001 × 50,000 = 40,000 VNĐ
5. Check if displayed cost matches

If displayed cost is 40,000,000 VNĐ (1000x more), then conversion is backwards.
If displayed cost is 40 VNĐ (1000x less), then conversion is inverted.
