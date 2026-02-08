# Ingredient Cost Adjustment - Should We Recalculate Unit Price?

## Question
When adjusting ingredient stock, should we recalculate the unit price (cost_per_unit)?

## Current Implementation

### Backend Logic
```go
// Calculate weighted average cost for stock IN (positive quantity)
if req.Quantity > 0 && costPerUnit > 0 && afterQty > 0 {
    // Weighted average: (old_qty * old_price + new_qty * new_price) / total_qty
    oldValue := beforeQty * item.CostPerUnit
    newValue := req.Quantity * costPerUnit
    item.CostPerUnit = (oldValue + newValue) / afterQty
}
```

**Current behavior:**
- ✅ Stock IN (positive quantity) with price → Recalculates weighted average
- ❌ Stock IN without price (cost_per_unit = 0) → Keeps current price
- ❌ Stock OUT (negative quantity) → Keeps current price
- ❌ Stock SET (adjust type) → Keeps current price

## Analysis by Operation Type

### 1. Stock IN (Nhập hàng) - ADD

#### Scenario A: Purchase at Same Price
```
Current: 10 kg @ 50,000đ/kg
Purchase: 5 kg @ 50,000đ/kg

Should recalculate? NO (same price)
Result: 15 kg @ 50,000đ/kg ✓
```

#### Scenario B: Purchase at Different Price
```
Current: 10 kg @ 50,000đ/kg
Purchase: 5 kg @ 60,000đ/kg

Should recalculate? YES (different price)
Calculation: (10 × 50,000 + 5 × 60,000) / 15
Result: 15 kg @ 53,333đ/kg ✓
```

#### Scenario C: Purchase Without Specifying Price
```
Current: 10 kg @ 50,000đ/kg
Purchase: 5 kg (no price entered)

Should recalculate? NO (assume same price)
Result: 15 kg @ 50,000đ/kg ✓
```

**Conclusion for Stock IN:**
- ✅ **YES, recalculate** when new price is provided
- ✅ **NO, keep current** when no price provided
- ✅ **Current implementation is CORRECT**

### 2. Stock OUT (Xuất hàng) - REMOVE

#### Scenario A: Usage/Consumption
```
Current: 10 kg @ 50,000đ/kg
Usage: 3 kg for cooking

Should recalculate? NO
Reason: We're using existing stock, not buying new
Result: 7 kg @ 50,000đ/kg ✓
```

#### Scenario B: Waste/Spoilage
```
Current: 10 kg @ 50,000đ/kg
Waste: 2 kg (expired)

Should recalculate? NO
Reason: Removing stock doesn't change the cost of remaining stock
Result: 8 kg @ 50,000đ/kg ✓
```

#### Scenario C: Theft/Loss
```
Current: 10 kg @ 50,000đ/kg
Loss: 1 kg (theft)

Should recalculate? NO
Reason: The remaining stock still cost 50,000đ/kg to acquire
Result: 9 kg @ 50,000đ/kg ✓
```

**Conclusion for Stock OUT:**
- ✅ **NO, never recalculate**
- ✅ **Keep current price** for remaining stock
- ✅ **Current implementation is CORRECT**

### 3. Stock SET (Điều chỉnh) - ADJUST

#### Scenario A: Inventory Count - Found More
```
Current: 10 kg @ 50,000đ/kg
Count: Actually 12 kg

Should recalculate? DEPENDS
- If extra stock was purchased → YES, need price
- If counting error → NO, keep current price
```

#### Scenario B: Inventory Count - Found Less
```
Current: 10 kg @ 50,000đ/kg
Count: Actually 8 kg

Should recalculate? NO
Reason: Missing stock doesn't change cost of remaining
Result: 8 kg @ 50,000đ/kg ✓
```

#### Scenario C: Correction After Error
```
Current: 10 kg @ 50,000đ/kg (wrong entry)
Correct: Should be 15 kg @ 45,000đ/kg

Should recalculate? YES
Reason: Correcting both quantity AND price
Need: Allow price input for SET operations
```

**Conclusion for Stock SET:**
- ⚠️ **COMPLEX** - depends on reason
- 🔧 **Current implementation may need improvement**
- 💡 **Should allow optional price input**

## Business Logic Principles

### Principle 1: FIFO/LIFO vs Weighted Average

**FIFO (First In, First Out):**
```
Purchase 1: 5 kg @ 40,000đ
Purchase 2: 5 kg @ 50,000đ
Use 3 kg → Use from Purchase 1 @ 40,000đ
Remaining: 2 kg @ 40,000đ + 5 kg @ 50,000đ
```
❌ Complex to track, requires batch management

**Weighted Average (Current):**
```
Purchase 1: 5 kg @ 40,000đ → Avg: 40,000đ
Purchase 2: 5 kg @ 50,000đ → Avg: 45,000đ
Use 3 kg → Remaining: 7 kg @ 45,000đ
```
✅ Simple, practical for restaurants

**Conclusion:** Weighted average is appropriate for this use case.

### Principle 2: Cost Basis

The `cost_per_unit` represents:
- **Acquisition cost** (what we paid)
- **NOT current market price**
- **NOT replacement cost**

When we remove stock:
- The remaining stock still cost the same to acquire
- We don't recalculate based on current market prices
- We maintain historical cost basis

**Conclusion:** Don't recalculate on removal.

### Principle 3: Inventory Valuation

```
Inventory Value = Quantity × Cost Per Unit

Example:
10 kg @ 50,000đ/kg = 500,000đ total value
```

When we remove 3 kg:
```
Before: 10 kg @ 50,000đ = 500,000đ
After:  7 kg @ 50,000đ = 350,000đ
```

The unit cost stays the same, only total value changes.

**Conclusion:** Removal doesn't affect unit cost.

## Edge Cases & Special Situations

### Edge Case 1: Price Increase During Shortage
```
Current: 2 kg @ 50,000đ/kg (low stock)
Purchase: 10 kg @ 70,000đ/kg (price increased!)

Weighted Average:
(2 × 50,000 + 10 × 70,000) / 12 = 66,667đ/kg

Result: 12 kg @ 66,667đ/kg
```
✅ Correctly reflects the blended cost

### Edge Case 2: Bulk Discount
```
Current: 5 kg @ 50,000đ/kg
Purchase: 20 kg @ 45,000đ/kg (bulk discount)

Weighted Average:
(5 × 50,000 + 20 × 45,000) / 25 = 46,000đ/kg

Result: 25 kg @ 46,000đ/kg
```
✅ Correctly reflects the lower average cost

### Edge Case 3: Waste After Price Increase
```
Current: 10 kg @ 60,000đ/kg (after price increase)
Waste: 3 kg (expired)

Should we reduce cost because we wasted expensive stock? NO
Remaining: 7 kg @ 60,000đ/kg

Reason: The remaining stock still cost 60,000đ/kg to acquire
```
✅ Keep current price

### Edge Case 4: Found Extra Stock (Unknown Source)
```
Current: 10 kg @ 50,000đ/kg
Found: 2 kg extra during inventory

Options:
A) Assume same cost: 12 kg @ 50,000đ/kg
B) Assume zero cost: (10 × 50,000 + 2 × 0) / 12 = 41,667đ/kg
C) Ask for price: Let user specify

Best: Option C (ask for price)
```
⚠️ Current implementation doesn't handle this well

## Recommendations

### Current Implementation Assessment

| Operation | Current Behavior | Is Correct? | Notes |
|-----------|-----------------|-------------|-------|
| Stock IN with price | Recalculates weighted avg | ✅ YES | Perfect |
| Stock IN without price | Keeps current price | ✅ YES | Good default |
| Stock OUT | Keeps current price | ✅ YES | Correct |
| Stock SET (increase) | Keeps current price | ⚠️ MAYBE | Should allow price input |
| Stock SET (decrease) | Keeps current price | ✅ YES | Correct |

### Improvements Needed

#### 1. Stock SET with Price Input (Optional)
```javascript
// When type = 'adjust' and quantity increases
if (adjustData.type === 'adjust' && newQty > currentQty) {
  // Show optional price input
  // "Nếu tăng do nhập thêm, nhập giá:"
  // If price provided, calculate weighted average
}
```

#### 2. Clear UI Guidance
```
Stock IN (➕):
- Always show price input
- Default to current price
- Calculate weighted average if different

Stock OUT (➖):
- Never show price input
- Keep current price
- No calculation needed

Stock SET (📋):
- Show optional price input
- Only if quantity increases
- Label: "Giá nhập (nếu có)"
```

#### 3. History Display
```
Show in history:
- Old price: 50,000đ/kg
- New price: 53,333đ/kg (if changed)
- Reason for change: "Nhập thêm với giá khác"
```

## Implementation Recommendations

### Keep Current Logic ✅
```go
// This is CORRECT - keep it
if req.Quantity > 0 && costPerUnit > 0 && afterQty > 0 {
    oldValue := beforeQty * item.CostPerUnit
    newValue := req.Quantity * costPerUnit
    item.CostPerUnit = (oldValue + newValue) / afterQty
}
```

### Add Optional Price for SET Operations
```go
// For type="adjust" with quantity increase
if req.Type == "adjust" && req.Quantity > 0 && costPerUnit > 0 {
    // Calculate as if it's a purchase
    oldValue := beforeQty * item.CostPerUnit
    newValue := req.Quantity * costPerUnit
    item.CostPerUnit = (oldValue + newValue) / afterQty
}
```

### Frontend Changes
```javascript
// In adjust modal, when type = 'adjust'
if (adjustData.type === 'adjust' && newQty > currentQty) {
  showPriceInput = true
  priceLabel = "Giá nhập thêm (tùy chọn)"
  priceHelp = "Chỉ nhập nếu tăng do mua thêm hàng"
}
```

## Conclusion

### Current Implementation: 90% Correct ✅

**What's working:**
- ✅ Stock IN with price → Weighted average (CORRECT)
- ✅ Stock IN without price → Keep current (CORRECT)
- ✅ Stock OUT → Keep current (CORRECT)

**What needs improvement:**
- ⚠️ Stock SET (increase) → Should allow optional price input
- ⚠️ No clear guidance for inventory corrections

### Final Answer

**Should we recalculate unit price when adjusting?**

| Operation | Recalculate? | When? |
|-----------|--------------|-------|
| Stock IN | ✅ YES | When new price provided |
| Stock OUT | ❌ NO | Never |
| Stock SET (increase) | ⚠️ OPTIONAL | If due to purchase |
| Stock SET (decrease) | ❌ NO | Never |

**Key Principle:**
> Only recalculate when **new stock is acquired at a different price**. Never recalculate when removing or correcting existing stock.

### Recommended Actions

1. ✅ **Keep current weighted average logic** for Stock IN
2. ✅ **Keep current "no recalculation" logic** for Stock OUT
3. 🔧 **Add optional price input** for Stock SET when quantity increases
4. 📝 **Add clear UI labels** to guide users
5. 📊 **Show price changes** in history when they occur

The current implementation is fundamentally sound and follows correct accounting principles. Minor improvements would make it more flexible for edge cases.
