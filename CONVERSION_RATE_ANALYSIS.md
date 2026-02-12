# Phân tích Conversion Rate & Implementation Guide

## 📊 Phân tích Conversion Rate

### 1. Vấn đề cần giải quyết

**Scenario thực tế**:
```
Quản lý kho: Mua sữa theo lít (L)
Công thức món: Dùng sữa theo ml
Vấn đề: Làm sao tính cost chính xác?
```

**Không có conversion rate**:
```
❌ Sai: 150ml sữa * 50,000 VND/L = 7,500,000 VND (sai hoàn toàn!)
```

**Có conversion rate**:
```
✅ Đúng: 150ml * 50,000 VND/L * (1/1000) = 7,500 VND
```

### 2. Cách hoạt động

#### Công thức
```
Cost = Quantity × Cost_Per_Unit × Conversion_Rate × (1 + Wastage/100)
```

#### Ý nghĩa từng thành phần

**Quantity**: Số lượng trong công thức món
- VD: 150 (ml trong công thức)

**Cost_Per_Unit**: Giá vốn trên đơn vị kho
- VD: 50,000 VND/lít

**Conversion_Rate**: Tỷ lệ quy đổi từ đơn vị kho → đơn vị công thức
- VD: 1/1000 (1L = 1000ml, nên 1ml = 1/1000 L)

**Wastage**: Tỷ lệ hao hụt %
- VD: 10% (tính thêm 10% do đổ tràn, hỏng)

### 3. Các trường hợp sử dụng

#### Case 1: Đơn vị giống nhau
```
Kho: gram (g)
Công thức: gram (g)
Conversion Rate: 1.0 (không cần quy đổi)
```

#### Case 2: Kg → gram
```
Kho: Kilogram (kg)
Công thức: gram (g)
Conversion Rate: 1000
Giải thích: 1kg = 1000g, nên 1g = 1/1000 kg

Tính toán:
- Cần: 20g cà phê
- Giá kho: 200,000 VND/kg
- Cost = 20g × 200,000 VND/kg × (1/1000) = 4,000 VND
```

#### Case 3: Lít → ml
```
Kho: Lít (L)
Công thức: Milliliter (ml)
Conversion Rate: 1000
Giải thích: 1L = 1000ml, nên 1ml = 1/1000 L

Tính toán:
- Cần: 100ml sữa
- Giá kho: 50,000 VND/L
- Cost = 100ml × 50,000 VND/L × (1/1000) = 5,000 VND
```

#### Case 4: Box → piece
```
Kho: Box (hộp)
Công thức: Piece (cái)
Conversion Rate: 12 (nếu 1 box = 12 pieces)
Giải thích: 1 box = 12 pieces, nên 1 piece = 1/12 box

Tính toán:
- Cần: 2 pieces bánh
- Giá kho: 120,000 VND/box
- Cost = 2 pieces × 120,000 VND/box × (1/12) = 20,000 VND
```

## 🎯 Gợi ý Implementation

### Option 1: Tự động tính Conversion Rate (Recommended)

#### Backend: Thêm Unit Conversion Logic

```go
// backend/domain/ingredient/unit_conversion.go
package ingredient

type UnitConversion struct {
    FromUnit UnitType
    ToUnit   UnitType
    Factor   float64
}

var conversionTable = map[string]float64{
    // Mass conversions
    "kg->g":   1000,
    "g->kg":   0.001,
    "kg->mg":  1000000,
    "mg->kg":  0.000001,
    
    // Volume conversions
    "L->ml":   1000,
    "ml->L":   0.001,
    "L->cl":   100,
    "cl->L":   0.01,
    
    // Same unit
    "kg->kg":  1,
    "g->g":    1,
    "L->L":    1,
    "ml->ml":  1,
}

func GetConversionFactor(fromUnit, toUnit UnitType) float64 {
    key := string(fromUnit) + "->" + string(toUnit)
    if factor, exists := conversionTable[key]; exists {
        return factor
    }
    // If no conversion found, assume same unit
    return 1.0
}

func CalculateConversionRate(stockUnit, recipeUnit UnitType) float64 {
    // Conversion rate = 1 / conversion factor
    // Because we want: recipe_quantity * stock_price * conversion_rate
    factor := GetConversionFactor(stockUnit, recipeUnit)
    if factor == 0 {
        return 1.0
    }
    return 1.0 / factor
}
```

#### Frontend: Auto-calculate khi chọn ingredient

```javascript
// frontend/src/composables/useUnitConversion.js
export function useUnitConversion() {
  const conversionTable = {
    // Mass
    'kg->g': 1000,
    'g->kg': 0.001,
    
    // Volume
    'L->ml': 1000,
    'ml->L': 0.001,
    
    // Same unit
    'kg->kg': 1,
    'g->g': 1,
    'L->L': 1,
    'ml->ml': 1,
  }
  
  const getConversionFactor = (fromUnit, toUnit) => {
    const key = `${fromUnit}->${toUnit}`
    return conversionTable[key] || 1.0
  }
  
  const calculateConversionRate = (stockUnit, recipeUnit) => {
    const factor = getConversionFactor(stockUnit, recipeUnit)
    return factor === 0 ? 1.0 : 1.0 / factor
  }
  
  const estimateCost = (quantity, recipeUnit, costPerUnit, stockUnit, wastage = 0) => {
    const conversionRate = calculateConversionRate(stockUnit, recipeUnit)
    const cost = quantity * costPerUnit * conversionRate * (1 + wastage / 100)
    return Math.round(cost * 100) / 100
  }
  
  return {
    getConversionFactor,
    calculateConversionRate,
    estimateCost
  }
}
```

#### UI Implementation trong MenuView

```vue
<!-- Khi chọn ingredient -->
<div v-for="(ingredient, index) in form.ingredients" :key="index">
  <div class="ingredient-item">
    <div class="ingredient-info">
      <span class="name">{{ ingredient.name }}</span>
      <span class="stock-unit">Kho: {{ ingredient.stockUnit }}</span>
    </div>
    
    <!-- Recipe Unit Selector -->
    <div class="recipe-unit">
      <label>Đơn vị công thức:</label>
      <select v-model="ingredient.recipeUnit" @change="updateConversion(index)">
        <option value="g">gram (g)</option>
        <option value="kg">kilogram (kg)</option>
        <option value="ml">milliliter (ml)</option>
        <option value="L">liter (L)</option>
      </select>
    </div>
    
    <!-- Quantity Input -->
    <div class="quantity">
      <label>Số lượng:</label>
      <input v-model.number="ingredient.quantity" type="number" step="0.1" />
      <span>{{ ingredient.recipeUnit }}</span>
    </div>
    
    <!-- Auto-calculated Conversion Rate (read-only, shown for transparency) -->
    <div class="conversion-info" v-if="ingredient.conversionRate !== 1">
      <span class="label">Quy đổi:</span>
      <span class="value">{{ ingredient.conversionRate }}</span>
      <span class="explanation">
        (1{{ ingredient.recipeUnit }} = {{ 1/ingredient.conversionRate }}{{ ingredient.stockUnit }})
      </span>
    </div>
    
    <!-- Estimated Cost Preview -->
    <div class="cost-preview">
      <span class="label">Chi phí ước tính:</span>
      <span class="value">{{ formatPrice(ingredient.estimatedCost) }}</span>
    </div>
  </div>
</div>
```

```javascript
// Script
import { useUnitConversion } from '@/composables/useUnitConversion'

const { calculateConversionRate, estimateCost } = useUnitConversion()

const selectIngredient = (ingredient) => {
  // Add ingredient with auto-calculated conversion
  const recipeUnit = ingredient.unit // Default to stock unit
  const conversionRate = calculateConversionRate(ingredient.unit, recipeUnit)
  
  form.value.ingredients.push({
    id: ingredient.id,
    name: ingredient.name,
    quantity: 1,
    stockUnit: ingredient.unit,
    recipeUnit: recipeUnit,
    conversionRate: conversionRate,
    costPerUnit: ingredient.cost_per_unit,
    wastage: ingredient.wastage_percentage || 0,
    estimatedCost: estimateCost(1, recipeUnit, ingredient.cost_per_unit, ingredient.unit, ingredient.wastage_percentage)
  })
}

const updateConversion = (index) => {
  const ing = form.value.ingredients[index]
  ing.conversionRate = calculateConversionRate(ing.stockUnit, ing.recipeUnit)
  ing.estimatedCost = estimateCost(
    ing.quantity, 
    ing.recipeUnit, 
    ing.costPerUnit, 
    ing.stockUnit, 
    ing.wastage
  )
}
```

### Option 2: Manual Input (Simpler, less user-friendly)

```vue
<!-- Cho phép user nhập conversion rate thủ công -->
<div class="conversion-manual">
  <label>Tỷ lệ quy đổi:</label>
  <input v-model.number="ingredient.conversionRate" type="number" step="0.001" />
  <span class="hint">
    VD: Kho (kg) → Công thức (g) = 1000
  </span>
</div>
```

### Option 3: Hybrid Approach (Best UX)

```vue
<div class="conversion-section">
  <!-- Auto mode (default) -->
  <div v-if="ingredient.autoConversion">
    <div class="auto-conversion">
      <span>✓ Tự động: {{ ingredient.conversionRate }}</span>
      <button @click="ingredient.autoConversion = false">Chỉnh thủ công</button>
    </div>
  </div>
  
  <!-- Manual mode -->
  <div v-else>
    <div class="manual-conversion">
      <input v-model.number="ingredient.conversionRate" type="number" step="0.001" />
      <button @click="resetToAuto(index)">Về tự động</button>
    </div>
  </div>
</div>
```

## 📝 Implementation Steps

### Phase 1: Backend Enhancement

1. **Add Unit Conversion Table**
   - Create `backend/domain/ingredient/unit_conversion.go`
   - Define conversion factors
   - Add helper functions

2. **Update Cost Calculator**
   - Modify to use conversion table
   - Add validation for unit compatibility

3. **API Enhancement**
   - Return `stock_unit` in ingredient response
   - Accept `recipe_unit` in menu item creation

### Phase 2: Frontend Implementation

1. **Create Composable**
   - `frontend/src/composables/useUnitConversion.js`
   - Implement conversion logic
   - Add cost estimation

2. **Update MenuView**
   - Add recipe unit selector
   - Auto-calculate conversion rate
   - Show cost preview
   - Add unit conversion hints

3. **Update IngredientManagementView**
   - Show conversion rate in ingredient details
   - Add wastage percentage input
   - Show cost calculation examples

### Phase 3: UX Enhancements

1. **Visual Indicators**
   - Show conversion rate badge when != 1.0
   - Highlight unit mismatches
   - Display cost breakdown

2. **Validation**
   - Warn if units are incompatible
   - Suggest correct recipe unit
   - Validate conversion rate range

3. **Help & Documentation**
   - Add tooltips explaining conversion
   - Show examples for common conversions
   - Add conversion rate calculator tool

## 🎨 UI/UX Mockup

```
┌─────────────────────────────────────────┐
│ 🥘 Nguyên liệu: Sữa tươi               │
├─────────────────────────────────────────┤
│ Kho: 5 L @ 50,000 VND/L                │
│                                         │
│ Đơn vị công thức: [ml ▼]               │
│ Số lượng: [150] ml                      │
│                                         │
│ ℹ️ Quy đổi: 1ml = 0.001L               │
│ 💰 Chi phí: 7,500 VND                  │
│                                         │
│ Hao hụt: [10]% (tùy chọn)              │
│ 💰 Chi phí cuối: 8,250 VND             │
└─────────────────────────────────────────┘
```

## ⚠️ Edge Cases & Validation

### 1. Incompatible Units
```javascript
// Không cho phép quy đổi giữa mass và volume
if (isMassUnit(stockUnit) && isVolumeUnit(recipeUnit)) {
  alert('Không thể quy đổi giữa khối lượng và thể tích!')
  return false
}
```

### 2. Zero or Negative Conversion
```javascript
if (conversionRate <= 0) {
  conversionRate = 1.0 // Fallback to default
}
```

### 3. Very Large/Small Numbers
```javascript
if (conversionRate > 1000000 || conversionRate < 0.000001) {
  alert('Tỷ lệ quy đổi không hợp lý. Vui lòng kiểm tra lại đơn vị.')
}
```

## 📊 Testing Scenarios

### Test Case 1: Same Unit
```
Stock: 100g @ 200,000 VND/kg
Recipe: 20g
Expected: 20 * 200,000 * 1 = 4,000 VND
```

### Test Case 2: Kg to Gram
```
Stock: 1kg @ 200,000 VND/kg
Recipe: 20g
Conversion: 1/1000
Expected: 20 * 200,000 * (1/1000) = 4,000 VND
```

### Test Case 3: With Wastage
```
Stock: 1L @ 50,000 VND/L
Recipe: 100ml
Conversion: 1/1000
Wastage: 10%
Expected: 100 * 50,000 * (1/1000) * 1.1 = 5,500 VND
```

## 🚀 Recommended Approach

**Tôi khuyên dùng Option 3 (Hybrid)**:

1. **Default**: Auto-calculate conversion rate
2. **Flexibility**: Cho phép override thủ công nếu cần
3. **Transparency**: Hiển thị rõ conversion rate và cost preview
4. **User-friendly**: Không yêu cầu user hiểu công thức phức tạp

**Benefits**:
- ✅ Giảm lỗi nhập liệu
- ✅ Tính toán chính xác
- ✅ Dễ sử dụng cho non-technical users
- ✅ Vẫn linh hoạt cho advanced users
- ✅ Transparent - user thấy được cách tính

**Implementation Priority**:
1. Phase 1: Backend unit conversion table
2. Phase 2: Frontend auto-calculation
3. Phase 3: UX enhancements & validation
