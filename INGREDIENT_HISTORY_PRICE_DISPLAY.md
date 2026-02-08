# Ingredient Stock History - Enhanced Price Display

## Summary
Improved the stock history modal to prominently display purchase prices and cost information, making it easier to track ingredient pricing over time.

## Key Improvements

### 1. Enhanced Header
**Before:**
- Simple title: "Lịch sử tồn kho"
- Ingredient name only

**After:**
- Clear title: "Lịch sử nhập hàng" (Purchase History)
- Ingredient name with emphasis
- Current price displayed for reference
```vue
<p class="text-xs text-gray-500 text-center">
  Giá hiện tại: <span class="font-semibold text-green-600">
    {{ formatCurrency(currentIngredient?.cost_per_unit) }}/{{ currentIngredient?.unit }}
  </span>
</p>
```

### 2. Visual Distinction by Transaction Type
**Color-coded cards:**
- Green border/background: Stock IN (purchases)
- Red border/background: Stock OUT (usage/waste)

```vue
:class="record.quantity > 0 ? 'border-green-200 bg-green-50' : 'border-red-200 bg-red-50'"
```

### 3. Prominent Price Information Section
For purchase transactions (quantity > 0), display a highlighted price card:

**Features:**
- Gradient green background (from-green-100 to-emerald-100)
- Green border for emphasis
- Money emoji (💰) header
- "THÔNG TIN GIÁ" label in uppercase

**Information displayed:**
1. **Unit Price**: Price per unit for this transaction
2. **Total Cost**: Total amount paid
3. **Calculation**: Shows the math (quantity × unit price)
4. **Price Comparison**: Compares with current price (↑ or ↓)

```vue
<div class="bg-gradient-to-br from-green-100 to-emerald-100 border-2 border-green-300 rounded-xl p-3">
  <!-- Unit Price -->
  <div class="flex items-center justify-between bg-white rounded-lg px-3 py-2">
    <span class="text-xs text-gray-600">Đơn giá lần này:</span>
    <span class="font-bold text-green-700">
      {{ formatCurrency(record.cost_per_unit) }}/{{ unit }}
    </span>
  </div>
  
  <!-- Total Cost -->
  <div class="flex items-center justify-between bg-white rounded-lg px-3 py-2">
    <span class="text-xs text-gray-600">Tổng chi phí:</span>
    <span class="font-bold text-green-800 text-base">
      {{ formatCurrency(record.total_cost) }}
    </span>
  </div>
  
  <!-- Calculation -->
  <div class="text-xs text-green-700 text-center bg-green-50 rounded-lg py-1">
    = {{ quantity }} {{ unit }} × {{ formatCurrency(cost_per_unit) }}
  </div>
  
  <!-- Price Comparison -->
  <div class="flex items-center justify-center gap-2 text-xs pt-1">
    <span class="text-gray-600">So với giá hiện tại:</span>
    <span class="text-red-600 font-semibold">↑ 5,000 ₫</span>
  </div>
</div>
```

### 4. Clear Transaction Header
**Improved layout:**
- Type badge on the left
- Large, bold quantity with +/- sign
- Color-coded (green for +, red for -)
- Before → After quantities on the right

```vue
<div class="flex items-center justify-between mb-3">
  <div class="flex items-center gap-2">
    <span class="badge">Nhập thêm</span>
    <span class="text-lg font-bold text-green-700">+5 kg</span>
  </div>
  <div class="text-right">
    <p class="text-xs text-gray-500">Trước → Sau</p>
    <p class="text-sm font-bold">10 → 15</p>
  </div>
</div>
```

### 5. Reason Display
Moved to a separate white card for better readability:
```vue
<div class="mb-3 bg-white rounded-lg p-3">
  <p class="text-sm text-gray-700">
    <span class="font-semibold">Lý do:</span> {{ record.reason }}
  </p>
</div>
```

### 6. No-Price Indicator for Removals
For stock OUT transactions, show a clear message:
```vue
<div class="bg-red-100 border border-red-200 rounded-lg p-2">
  <p class="text-xs text-red-700 text-center">
    ⚠️ Xuất kho - Không có thông tin giá
  </p>
</div>
```

### 7. Enhanced Metadata Footer
```vue
<div class="flex items-center justify-between text-xs text-gray-500 pt-2 border-t">
  <div class="flex items-center gap-1">
    <span>👤</span>
    <span class="font-medium">{{ username }}</span>
  </div>
  <div class="flex items-center gap-1">
    <span>🕐</span>
    <span>{{ formatDateTime(created_at) }}</span>
  </div>
</div>
```

## Visual Hierarchy

### Priority 1: Transaction Type & Amount
- Large, bold quantity with color
- Clear badge for transaction type

### Priority 2: Price Information (for purchases)
- Highlighted card with gradient background
- Unit price and total cost prominently displayed
- Calculation breakdown for transparency

### Priority 3: Context
- Reason in separate card
- Before/After quantities
- User and timestamp

## Use Cases

### Use Case 1: Track Price Changes
**Scenario:** Manager wants to see if ingredient prices are increasing

**Solution:**
1. Open ingredient history
2. See current price at top
3. Scroll through purchase records
4. Each record shows:
   - Price at time of purchase
   - Comparison with current price (↑ or ↓)
   - Visual indicator (red for increase, green for decrease)

### Use Case 2: Verify Purchase Cost
**Scenario:** Manager wants to verify a recent purchase

**Solution:**
1. Find the purchase record (green card)
2. See highlighted price section:
   - Unit price paid
   - Total cost
   - Calculation breakdown
3. Compare with receipt

### Use Case 3: Analyze Cost Trends
**Scenario:** Manager wants to understand cost trends

**Solution:**
1. Scroll through history chronologically
2. Each purchase shows:
   - Date and time
   - Quantity purchased
   - Price paid
   - Comparison with current price
3. Identify patterns (seasonal changes, supplier changes, etc.)

## Mobile Optimization

### Touch-Friendly
- Large cards with good spacing
- Clear visual separation
- Easy to scroll

### Information Density
- Balanced: Not too crowded, not too sparse
- Important info (price) is prominent
- Secondary info (metadata) is smaller but readable

### Color Coding
- Green: Positive (purchases, stock in)
- Red: Negative (usage, stock out)
- Intuitive and consistent

## Benefits

### 1. Price Transparency
- Clear display of historical prices
- Easy to track price changes
- Comparison with current price

### 2. Cost Verification
- Total cost clearly shown
- Calculation breakdown provided
- Easy to verify against receipts

### 3. Better Decision Making
- Historical price data visible
- Trends easy to identify
- Informed purchasing decisions

### 4. Professional Appearance
- Clean, modern design
- Color-coded for clarity
- Gradient effects for emphasis

## Technical Details

### Conditional Rendering
```javascript
// Show price info only for purchases with price data
v-if="record.cost_per_unit > 0 && record.quantity > 0"

// Show no-price indicator for removals
v-else-if="record.quantity < 0"
```

### Price Comparison Logic
```javascript
// Calculate difference
const priceDiff = Math.abs(record.cost_per_unit - currentIngredient?.cost_per_unit)

// Determine direction
const isIncrease = record.cost_per_unit > currentIngredient?.cost_per_unit

// Show arrow
{{ isIncrease ? '↑' : '↓' }} {{ formatCurrency(priceDiff) }}
```

### Color Classes
```javascript
// Dynamic color based on price change
:class="record.cost_per_unit > currentIngredient?.cost_per_unit 
  ? 'text-red-600'    // Price increased
  : 'text-green-600'" // Price decreased
```

## Future Enhancements

Consider adding:
1. **Price chart**: Visual graph of price changes over time
2. **Average price**: Calculate and display average purchase price
3. **Price alerts**: Notify when price changes significantly
4. **Export**: Export history to CSV/Excel
5. **Filters**: Filter by date range, price range, user
6. **Statistics**: Min/max/average prices, total spent

## Files Modified

1. `frontend/src/views/IngredientManagementView.vue`
   - Enhanced stock history modal UI
   - Added price information section
   - Added color coding for transaction types
   - Added price comparison logic
   - Improved layout and spacing
