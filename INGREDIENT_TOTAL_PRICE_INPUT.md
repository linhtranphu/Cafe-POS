# Ingredient Total Price Input Implementation

## Overview
Implemented dual-mode price input for ingredients: users can enter either **total price** (default) or **unit price**, making data entry more intuitive.

## Problem Solved
Users typically know the **total amount paid** rather than the **unit price**:
- "Bought 2 bottles of water for 10,000 VND" (not "5,000 VND per bottle")
- "Bought 0.5kg coffee for 200,000 VND" (not "400,000 VND per kg")

## Solution
Added toggle between two input modes with automatic calculation.

## Implementation

### 1. Two Input Modes

#### Mode 1: Total Price (Default) ✅
```
User enters:
- Quantity: 2 chai
- Total price: 10,000 VND

System calculates:
- Unit price: 10,000 ÷ 2 = 5,000 VND/chai
- Expense: 10,000 VND
```

#### Mode 2: Unit Price
```
User enters:
- Quantity: 2 chai  
- Unit price: 5,000 VND/chai

System calculates:
- Total: 2 × 5,000 = 10,000 VND
- Expense: 10,000 VND
```

### 2. UI Features

#### Toggle Button
- Located in pricing section header
- Text changes based on current mode:
  - In unit mode: "Nhập tổng giá"
  - In total mode: "Nhập đơn giá"

#### Visual Differentiation
- **Total price input**: Green border (primary mode)
- **Unit price input**: Gray border (alternative mode)

#### Real-time Calculation Display
Shows calculated value in blue card:
```
📊 Đơn giá (tự động tính)
5,000₫/chai
= 10,000₫ ÷ 2 chai
```

### 3. Code Structure

#### New State Variables
```javascript
const priceInputMode = ref('total') // 'total' or 'unit'
const totalPriceInput = ref(0)
```

#### Toggle Function
```javascript
const togglePriceInputMode = () => {
  priceInputMode.value = priceInputMode.value === 'unit' ? 'total' : 'unit'
  // Clear the other input when switching
  if (priceInputMode.value === 'unit') {
    totalPriceInput.value = 0
  } else {
    // Calculate total from unit price when switching
    if (formData.value.quantity > 0 && formData.value.cost_per_unit > 0) {
      totalPriceInput.value = formData.value.quantity * formData.value.cost_per_unit
    }
  }
}
```

#### Calculate Unit Price
```javascript
const calculateUnitPrice = () => {
  if (formData.value.quantity > 0 && totalPriceInput.value > 0) {
    formData.value.cost_per_unit = totalPriceInput.value / formData.value.quantity
  } else {
    formData.value.cost_per_unit = 0
  }
}
```

#### Auto-recalculate on Quantity Change
```html
<input v-model.number="formData.quantity" 
  @input="priceInputMode === 'total' && calculateUnitPrice()"
  ... />
```

### 4. Backend Integration

No backend changes needed! The system still stores `cost_per_unit`:
- In total mode: calculated from total ÷ quantity
- In unit mode: entered directly

Expense calculation remains: `amount = cost_per_unit × quantity`

## User Flow Examples

### Example 1: Water Bottles (Total Price Mode)

1. **Open create form** → Default to total price mode
2. **Enter basic info:**
   - Name: "Nước suối"
   - Unit: "chai"
   - Quantity: 2
3. **Enter total price:** 10,000 VND
4. **See calculation:**
   - Blue card shows: "Đơn giá: 5,000₫/chai"
   - Formula: "= 10,000₫ ÷ 2 chai"
5. **See expense indicator:**
   - Green card: "Chi phí: 10,000₫"
6. **Submit** → Saved with cost_per_unit = 5,000

### Example 2: Coffee (Switch to Unit Price Mode)

1. **Open create form**
2. **Click toggle:** "Nhập đơn giá"
3. **Enter info:**
   - Name: "Cà phê hạt"
   - Unit: "kg"
   - Quantity: 0.5
   - Unit price: 400,000 VND/kg
4. **See total:** "Tổng chi phí: 200,000₫"
5. **Submit** → Saved with cost_per_unit = 400,000

### Example 3: Edit Existing Ingredient

1. **Click edit** on existing ingredient
2. **Form loads** with:
   - Quantity: 2
   - cost_per_unit: 5,000
   - totalPriceInput: auto-calculated to 10,000
3. **Can switch modes** to adjust either value
4. **Submit** → Updates with new cost_per_unit

## Benefits

1. **Intuitive**: Matches how users think about purchases
2. **Flexible**: Supports both input methods
3. **Accurate**: Automatic calculation prevents errors
4. **Transparent**: Shows formula for verification
5. **Consistent**: Same expense tracking regardless of mode

## Edge Cases Handled

### Division by Zero
```javascript
if (formData.value.quantity > 0 && totalPriceInput.value > 0) {
  // Calculate
} else {
  formData.value.cost_per_unit = 0
}
```

### Mode Switching
- Clears opposite input to avoid confusion
- Preserves data when switching back

### Edit Mode
- Auto-calculates total from existing data
- Allows editing in either mode

## Testing Checklist

- [ ] Create ingredient in total price mode
- [ ] Create ingredient in unit price mode
- [ ] Toggle between modes
- [ ] Change quantity in total price mode (auto-recalculates)
- [ ] Edit existing ingredient
- [ ] Verify expense amount is correct
- [ ] Test with decimal quantities (0.5 kg)
- [ ] Test with large numbers (millions)
- [ ] Test edge cases (zero quantity, zero price)

## Formula Reference

### Total Price Mode
```
cost_per_unit = total_price ÷ quantity
expense_amount = total_price
```

### Unit Price Mode
```
total_price = cost_per_unit × quantity
expense_amount = cost_per_unit × quantity
```

Both modes result in the same expense amount!

---
**Status**: ✅ Complete
**Date**: 2026-02-07
**Default Mode**: Total Price (more intuitive)
**Backend Changes**: None (UI only)
