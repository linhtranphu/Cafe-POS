# Ingredient Adjust Stock - Total Price Input Mode

## Summary
Added total price input mode to the stock adjustment modal, making it consistent with the create ingredient flow and more user-friendly.

## Changes Made

### 1. Added Price Input Mode Toggle
- **Default mode**: Total price (more intuitive for users)
- **Alternative mode**: Unit price (for advanced users)
- Toggle button to switch between modes

### 2. Total Price Mode Features
```javascript
// When user enters total price:
- Input: Total amount paid (e.g., 50,000 VND for 2 kg)
- System calculates: Unit price = Total / Quantity
- Shows: Calculated unit price with formula
- Reference: Current unit price for comparison
```

### 3. Unit Price Mode Features
```javascript
// When user enters unit price:
- Input: Price per unit (e.g., 25,000 VND/kg)
- System uses: Entered unit price directly
- Shows: Current price as placeholder
- Note: Leave as 0 to use current price
```

### 4. UI Improvements
- **Visual distinction**: Green border for total price, gray for unit price
- **Real-time calculation**: Shows calculated unit price immediately
- **Context help**: Displays current price for reference
- **Smart defaults**: Clears opposite field when switching modes

## User Flow

### Scenario 1: User knows total price (most common)
1. Click "Điều chỉnh" on ingredient
2. Select "Nhập thêm" (Add stock)
3. Enter quantity: `2 kg`
4. Enter total price: `50,000 VND`
5. System shows: "Đơn giá được tính: 25,000 ₫/kg"
6. Auto-expense created: 50,000 VND

### Scenario 2: User knows unit price
1. Click toggle: "Nhập đơn giá"
2. Enter quantity: `2 kg`
3. Enter unit price: `25,000 VND/kg`
4. System calculates: Total = 50,000 VND
5. Auto-expense created: 50,000 VND

### Scenario 3: Price unchanged
1. Enter quantity only
2. Leave price as 0 (or don't enter)
3. System uses current ingredient price
4. Auto-expense uses current price

## Technical Implementation

### New State Variables
```javascript
const adjustPriceMode = ref('total') // 'total' or 'unit'
const adjustTotalPrice = ref(0)
```

### New Functions
```javascript
toggleAdjustPriceMode() // Switch between modes
calculateAdjustUnitPrice() // Total → Unit price
```

### Reactive Calculation
```javascript
// When quantity or total price changes:
@input="adjustPriceMode === 'total' && calculateAdjustUnitPrice()"

// Formula:
unit_price = total_price / quantity
```

## Benefits

### 1. Consistency
- Same UX as create ingredient modal
- Users learn once, use everywhere

### 2. User-Friendly
- Most users know total paid, not unit price
- No manual calculation needed
- Reduces input errors

### 3. Flexibility
- Power users can still enter unit price directly
- Easy toggle between modes
- Preserves existing functionality

### 4. Transparency
- Shows calculation formula
- Displays current price for reference
- Clear visual feedback

## Auto-Expense Integration

The auto-expense calculation works with both modes:

```javascript
effectiveAdjustPrice = adjustData.cost_per_unit > 0 
  ? adjustData.cost_per_unit 
  : currentIngredient.cost_per_unit

expenseAmount = quantity × effectiveAdjustPrice
```

**Result**: Whether user enters total or unit price, the expense is recorded correctly.

## Testing Checklist

- [ ] Total price mode: Enter total, verify unit price calculation
- [ ] Unit price mode: Enter unit price, verify expense calculation
- [ ] Toggle between modes: Verify fields clear/recalculate
- [ ] Leave price as 0: Verify uses current price
- [ ] Change quantity: Verify unit price recalculates in total mode
- [ ] Auto-expense: Verify correct amount recorded
- [ ] Price update: Verify ingredient cost_per_unit updates if new price entered

## Files Modified

1. `frontend/src/views/IngredientManagementView.vue`
   - Added `adjustPriceMode` and `adjustTotalPrice` state
   - Added `toggleAdjustPriceMode()` function
   - Added `calculateAdjustUnitPrice()` function
   - Updated adjust modal UI with mode toggle
   - Updated quantity input to trigger recalculation

## Next Steps

Consider applying the same pattern to:
- Facility stock adjustments (if applicable)
- Any other inventory management features
- Expense entry forms (reverse: enter total, calculate breakdown)
