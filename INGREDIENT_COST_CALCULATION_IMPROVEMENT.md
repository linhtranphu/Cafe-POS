# Ingredient Cost Calculation UI Improvement

## Overview
Improved the ingredient management UI to clearly show how cost is calculated and how expense is automatically tracked.

## Problem
The previous UI was not clear about:
- What "cost_per_unit" means (price per unit)
- How total expense is calculated (quantity × cost_per_unit)
- When expense is automatically created

## Solution
Enhanced UI with clear labels, calculations, and visual indicators.

## Changes Made

### 1. Create/Edit Ingredient Form

#### Field Labels Improved
**Before:**
- "Số lượng" (ambiguous)
- "Giá/Đơn vị" (unclear)

**After:**
- "Số lượng nhập" with unit hint: `(kg)`
- "Đơn giá" with unit price hint: `(VND/kg)`
- Added placeholder examples: "VD: 10", "VD: 200000"

#### Added Total Cost Display
New blue card showing:
```
💰 Tổng chi phí: 2,000,000₫
= 10 kg × 200,000₫/kg
```

This appears when both quantity and cost_per_unit are entered.

#### Enhanced Auto-Expense Indicator
**Before:**
- Simple green box with total amount

**After:**
- Prominent green card with:
  - Clear title: "Tự động ghi nhận chi phí"
  - White box showing expense amount
  - Calculation breakdown: `(10 kg × 200,000₫)`
  - Category and payment method info

### 2. Ingredient List Display

#### Improved Info Layout
**Before:**
```
Tồn kho: 10 kg
Tối thiểu: 5 kg
Đơn giá: 200,000₫
```

**After:**
```
📦 Tồn kho: 10 kg (bold, larger)
⚠️ Tối thiểu: 5 kg
💵 Đơn giá: 200,000₫/kg (highlighted in green box)
```

### 3. Stock Adjustment Modal

#### Enhanced Auto-Expense Indicator
When adjusting stock IN (adding quantity):
- Shows calculated expense for the added quantity
- Formula: `added_quantity × cost_per_unit`
- Example: `5 kg × 200,000₫/kg = 1,000,000₫`

### 4. Added Computed Property

```javascript
const totalCost = computed(() => {
  const quantity = formData.value.quantity || 0
  const costPerUnit = formData.value.cost_per_unit || 0
  return quantity * costPerUnit
})
```

## Visual Improvements

### Color Coding
- **Blue cards**: Calculation displays (informational)
- **Green cards**: Auto-expense indicators (action will be taken)
- **Icons**: 💰 for cost, 📝 for category, 💵 for payment method

### Typography
- **Bold** for important numbers
- **Larger font** for total amounts
- **Smaller gray text** for explanations

### Layout
- Responsive grid for form fields
- Clear visual hierarchy
- Adequate spacing between sections

## User Flow Example

### Creating New Ingredient

1. **Enter basic info:**
   - Name: "Cà phê hạt"
   - Category: "Nguyên liệu"
   - Unit: "kg"

2. **Enter quantity and pricing:**
   - Số lượng nhập: `10` kg
   - Đơn giá: `200,000` VND/kg

3. **See calculations:**
   - Blue card shows: "💰 Tổng chi phí: 2,000,000₫"
   - Formula: "= 10 kg × 200,000₫/kg"

4. **See auto-expense indicator:**
   - Green card confirms expense will be created
   - Shows exact amount: 2,000,000₫
   - Shows category and payment method

5. **Submit:**
   - Ingredient created ✓
   - Expense automatically created ✓

### Adjusting Stock

1. **Select ingredient** from list
2. **Click "📦 Điều chỉnh"**
3. **Choose type:** "Nhập thêm" (add)
4. **Enter quantity:** `5` kg
5. **See auto-expense:**
   - Green card shows: "Chi phí nhập thêm: 1,000,000₫"
   - Formula: "= 5 kg × 200,000₫/kg"
6. **Submit:**
   - Stock updated ✓
   - Expense created ✓

## Backend Logic (Unchanged)

The backend logic remains the same:
```go
amount := ing.CostPerUnit * quantity
```

This UI improvement only makes the calculation more visible and understandable to users.

## Benefits

1. **Clarity**: Users understand what they're entering
2. **Transparency**: Calculation is visible before submission
3. **Confidence**: Users know expense will be tracked correctly
4. **Education**: New users learn the system quickly
5. **Error Prevention**: Clear labels reduce input mistakes

## Files Modified

- `frontend/src/views/IngredientManagementView.vue`

## Testing Checklist

- [ ] Create ingredient with quantity and cost_per_unit
- [ ] Verify total cost calculation displays correctly
- [ ] Verify auto-expense indicator shows correct amount
- [ ] Adjust stock IN and verify expense indicator
- [ ] Check responsive layout on mobile
- [ ] Verify unit labels update dynamically
- [ ] Test with different units (kg, L, piece, etc.)
- [ ] Verify expense is created with correct amount

## Formula Reference

```
Total Cost = Quantity × Cost Per Unit

Examples:
- 10 kg × 200,000₫/kg = 2,000,000₫
- 5 L × 50,000₫/L = 250,000₫
- 20 piece × 5,000₫/piece = 100,000₫
```

---
**Status**: ✅ Complete
**Date**: 2026-02-07
**Impact**: UI clarity improvement, no backend changes
