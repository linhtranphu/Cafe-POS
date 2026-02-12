# Conversion Rate - Implementation Guide

## 🎯 Quick Start Implementation

### Step 1: Backend - Unit Conversion Table

Tạo file `backend/domain/ingredient/unit_conversion.go`:

```go
package ingredient

// GetConversionRate calculates the conversion rate from stock unit to recipe unit
// Returns the multiplier to convert recipe quantity to stock quantity
func GetConversionRate(stockUnit, recipeUnit UnitType) float64 {
    // If same unit, no conversion needed
    if stockUnit == recipeUnit {
        return 1.0
    }
    
    // Mass conversions (base: kg)
    massToKg := map[UnitType]float64{
        UnitKilogram: 1.0,
        UnitGram:     0.001, // 1g = 0.001kg
    }
    
    // Volume conversions (base: L)
    volumeToL := map[UnitType]float64{
        UnitLiter:      1.0,
        UnitMilliliter: 0.001, // 1ml = 0.001L
    }
    
    // Try mass conversion
    if stockFactor, stockOk := massToKg[stockUnit]; stockOk {
        if recipeFactor, recipeOk := massToKg[recipeUnit]; recipeOk {
            // conversion_rate = recipe_factor / stock_factor
            return recipeFactor / stockFactor
        }
    }
    
    // Try volume conversion
    if stockFactor, stockOk := volumeToL[stockUnit]; stockOk {
        if recipeFactor, recipeOk := volumeToL[recipeUnit]; recipeOk {
            return recipeFactor / stockFactor
        }
    }
    
    // No conversion found, return 1.0 (assume same unit)
    return 1.0
}

// ValidateUnitConversion checks if conversion between units is valid
func ValidateUnitConversion(stockUnit, recipeUnit UnitType) bool {
    // Same unit is always valid
    if stockUnit == recipeUnit {
        return true
    }
    
    massUnits := map[UnitType]bool{
        UnitKilogram: true,
        UnitGram:     true,
    }
    
    volumeUnits := map[UnitType]bool{
        UnitLiter:      true,
        UnitMilliliter: true,
    }
    
    // Both must be mass or both must be volume
    stockIsMass := massUnits[stockUnit]
    recipeIsMass := massUnits[recipeUnit]
    stockIsVolume := volumeUnits[stockUnit]
    recipeIsVolume := volumeUnits[recipeUnit]
    
    // Valid if both are mass or both are volume
    return (stockIsMass && recipeIsMass) || (stockIsVolume && recipeIsVolume)
}
```

### Step 2: Frontend - Composable

Tạo file `frontend/src/composables/useUnitConversion.js`:

```javascript
export function useUnitConversion() {
  // Unit categories
  const massUnits = ['kg', 'g']
  const volumeUnits = ['L', 'ml']
  const countUnits = ['piece', 'box', 'pack']
  
  // Conversion factors to base unit
  const massToKg = {
    'kg': 1.0,
    'g': 0.001
  }
  
  const volumeToL = {
    'L': 1.0,
    'ml': 0.001
  }
  
  /**
   * Calculate conversion rate from stock unit to recipe unit
   * @param {string} stockUnit - Unit used in stock/inventory
   * @param {string} recipeUnit - Unit used in recipe
   * @returns {number} - Conversion rate
   */
  const getConversionRate = (stockUnit, recipeUnit) => {
    // Same unit, no conversion
    if (stockUnit === recipeUnit) {
      return 1.0
    }
    
    // Try mass conversion
    if (massToKg[stockUnit] && massToKg[recipeUnit]) {
      return massToKg[recipeUnit] / massToKg[stockUnit]
    }
    
    // Try volume conversion
    if (volumeToL[stockUnit] && volumeToL[recipeUnit]) {
      return volumeToL[recipeUnit] / volumeToL[stockUnit]
    }
    
    // No conversion found
    return 1.0
  }
  
  /**
   * Validate if conversion between units is possible
   */
  const isValidConversion = (stockUnit, recipeUnit) => {
    if (stockUnit === recipeUnit) return true
    
    const stockIsMass = massUnits.includes(stockUnit)
    const recipeIsMass = massUnits.includes(recipeUnit)
    const stockIsVolume = volumeUnits.includes(stockUnit)
    const recipeIsVolume = volumeUnits.includes(recipeUnit)
    
    // Valid if both are same category
    return (stockIsMass && recipeIsMass) || (stockIsVolume && recipeIsVolume)
  }
  
  /**
   * Get compatible units for a given stock unit
   */
  const getCompatibleUnits = (stockUnit) => {
    if (massUnits.includes(stockUnit)) return massUnits
    if (volumeUnits.includes(stockUnit)) return volumeUnits
    if (countUnits.includes(stockUnit)) return countUnits
    return [stockUnit]
  }
  
  /**
   * Calculate estimated cost for an ingredient
   */
  const calculateCost = (quantity, recipeUnit, costPerUnit, stockUnit, wastage = 0) => {
    const conversionRate = getConversionRate(stockUnit, recipeUnit)
    const cost = quantity * costPerUnit * conversionRate * (1 + wastage / 100)
    return Math.round(cost * 100) / 100
  }
  
  /**
   * Format conversion explanation for display
   */
  const getConversionExplanation = (stockUnit, recipeUnit) => {
    const rate = getConversionRate(stockUnit, recipeUnit)
    if (rate === 1.0) {
      return 'Không cần quy đổi'
    }
    
    const inverseRate = 1 / rate
    return `1${recipeUnit} = ${inverseRate}${stockUnit}`
  }
  
  return {
    getConversionRate,
    isValidConversion,
    getCompatibleUnits,
    calculateCost,
    getConversionExplanation
  }
}
```

### Step 3: Update MenuView - Ingredient Selection

Cập nhật `frontend/src/views/MenuView.vue`:

```vue
<script setup>
import { useUnitConversion } from '@/composables/useUnitConversion'

const { 
  getConversionRate, 
  isValidConversion, 
  getCompatibleUnits,
  calculateCost,
  getConversionExplanation 
} = useUnitConversion()

// Khi select ingredient
const selectIngredient = (ingredient) => {
  if (isIngredientSelected(ingredient.id)) {
    return
  }
  
  // Get compatible units for this ingredient
  const compatibleUnits = getCompatibleUnits(ingredient.unit)
  
  // Default recipe unit = stock unit (no conversion)
  const recipeUnit = ingredient.unit
  const conversionRate = getConversionRate(ingredient.unit, recipeUnit)
  
  // Add ingredient with conversion info
  form.value.ingredients.push({
    id: ingredient.id,
    name: ingredient.name,
    quantity: 1,
    stockUnit: ingredient.unit,
    recipeUnit: recipeUnit,
    compatibleUnits: compatibleUnits,
    conversionRate: conversionRate,
    costPerUnit: ingredient.cost_per_unit,
    wastage: ingredient.wastage_percentage || 0,
    estimatedCost: calculateCost(
      1, 
      recipeUnit, 
      ingredient.cost_per_unit, 
      ingredient.unit, 
      ingredient.wastage_percentage || 0
    )
  })
  
  showIngredientSelector.value = false
  ingredientSearchQuery.value = ''
}

// Update conversion when recipe unit changes
const updateRecipeUnit = (index) => {
  const ing = form.value.ingredients[index]
  
  // Validate conversion
  if (!isValidConversion(ing.stockUnit, ing.recipeUnit)) {
    alert(`Không thể quy đổi từ ${ing.stockUnit} sang ${ing.recipeUnit}!`)
    ing.recipeUnit = ing.stockUnit // Reset to stock unit
    return
  }
  
  // Recalculate conversion rate
  ing.conversionRate = getConversionRate(ing.stockUnit, ing.recipeUnit)
  
  // Recalculate cost
  ing.estimatedCost = calculateCost(
    ing.quantity,
    ing.recipeUnit,
    ing.costPerUnit,
    ing.stockUnit,
    ing.wastage
  )
}

// Update cost when quantity changes
const updateQuantity = (index) => {
  const ing = form.value.ingredients[index]
  ing.estimatedCost = calculateCost(
    ing.quantity,
    ing.recipeUnit,
    ing.costPerUnit,
    ing.stockUnit,
    ing.wastage
  )
}
</script>

<template>
  <!-- Ingredient List in Form -->
  <div v-for="(ingredient, index) in form.ingredients" :key="index" 
    class="bg-gray-50 rounded-lg p-3">
    
    <!-- Header -->
    <div class="flex justify-between items-start mb-2">
      <div class="flex-1">
        <div class="font-medium text-gray-800">{{ ingredient.name }}</div>
        <div class="text-xs text-gray-500">
          Kho: {{ ingredient.stockUnit }} @ {{ formatPrice(ingredient.costPerUnit) }}/{{ ingredient.stockUnit }}
        </div>
      </div>
      <button type="button" @click="removeIngredient(index)" 
        class="bg-red-500 text-white px-3 py-1 rounded-lg hover:bg-red-600 flex-shrink-0 text-sm">
        ×
      </button>
    </div>
    
    <!-- Recipe Unit Selector -->
    <div class="mb-2">
      <label class="text-xs text-gray-600">Đơn vị công thức:</label>
      <select v-model="ingredient.recipeUnit" @change="updateRecipeUnit(index)"
        class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg">
        <option v-for="unit in ingredient.compatibleUnits" :key="unit" :value="unit">
          {{ unit }}
        </option>
      </select>
    </div>
    
    <!-- Quantity Input -->
    <div class="mb-2">
      <label class="text-xs text-gray-600">Số lượng:</label>
      <div class="flex gap-2 items-center">
        <input v-model.number="ingredient.quantity" 
          @input="updateQuantity(index)"
          type="number" min="0" step="0.1" placeholder="0" required
          class="flex-1 px-3 py-2 text-base border border-gray-300 rounded-lg" />
        <span class="text-sm text-gray-600">{{ ingredient.recipeUnit }}</span>
      </div>
    </div>
    
    <!-- Conversion Info (if not 1.0) -->
    <div v-if="ingredient.conversionRate !== 1" 
      class="mb-2 p-2 bg-blue-50 rounded text-xs text-blue-700">
      <span class="font-medium">ℹ️ Quy đổi:</span>
      {{ getConversionExplanation(ingredient.stockUnit, ingredient.recipeUnit) }}
    </div>
    
    <!-- Cost Preview -->
    <div class="p-2 bg-green-50 rounded">
      <div class="flex justify-between items-center">
        <span class="text-xs text-green-700">Chi phí ước tính:</span>
        <span class="text-sm font-bold text-green-700">
          {{ formatPrice(ingredient.estimatedCost) }}
        </span>
      </div>
      <div v-if="ingredient.wastage > 0" class="text-xs text-green-600 mt-1">
        (Bao gồm {{ ingredient.wastage }}% hao hụt)
      </div>
    </div>
  </div>
</template>
```

### Step 4: Update Ingredient Management

Thêm conversion rate và wastage vào ingredient form:

```vue
<!-- frontend/src/views/IngredientManagementView.vue -->
<template>
  <!-- In Create/Edit Form -->
  
  <!-- Wastage Percentage -->
  <div>
    <label class="block text-sm font-medium text-gray-700 mb-3">
      Hao hụt (%) - Tùy chọn
    </label>
    <input v-model.number="form.wastage_percentage" 
      type="number" min="0" max="100" step="0.1" placeholder="0"
      class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg" />
    <p class="text-xs text-gray-500 mt-1">
      VD: 10% cho rau củ (do cắt tỉa), 5% cho thịt cá
    </p>
  </div>
  
  <!-- Conversion Rate Info (Read-only, calculated automatically) -->
  <div class="p-4 bg-blue-50 rounded-lg">
    <h4 class="font-medium text-blue-900 mb-2">ℹ️ Quy đổi đơn vị</h4>
    <p class="text-sm text-blue-700">
      Hệ thống sẽ tự động quy đổi khi bạn tạo công thức món.
    </p>
    <div class="mt-2 text-xs text-blue-600">
      <div>• Kg ↔ gram: Tự động</div>
      <div>• Lít ↔ ml: Tự động</div>
      <div>• Cùng đơn vị: Không cần quy đổi</div>
    </div>
  </div>
</template>
```

## 🎨 Visual Examples

### Example 1: Cà phê (Kg → gram)

```
┌─────────────────────────────────────────┐
│ 🥘 Cà phê                               │
├─────────────────────────────────────────┤
│ Kho: kg @ 200,000 VND/kg                │
│                                         │
│ Đơn vị công thức: [gram ▼]             │
│ Số lượng: [20] gram                     │
│                                         │
│ ℹ️ Quy đổi: 1g = 0.001kg               │
│ 💰 Chi phí: 4,000 VND                  │
└─────────────────────────────────────────┘
```

### Example 2: Sữa (L → ml)

```
┌─────────────────────────────────────────┐
│ 🥛 Sữa tươi                             │
├─────────────────────────────────────────┤
│ Kho: L @ 50,000 VND/L                   │
│                                         │
│ Đơn vị công thức: [ml ▼]               │
│ Số lượng: [150] ml                      │
│                                         │
│ ℹ️ Quy đổi: 1ml = 0.001L               │
│ 💰 Chi phí: 7,500 VND                  │
│                                         │
│ Hao hụt: 10%                            │
│ 💰 Chi phí cuối: 8,250 VND             │
└─────────────────────────────────────────┘
```

## ✅ Testing Checklist

- [ ] Same unit (kg → kg): Conversion rate = 1.0
- [ ] Kg to gram: Conversion rate = 0.001
- [ ] Gram to kg: Conversion rate = 1000
- [ ] L to ml: Conversion rate = 0.001
- [ ] ml to L: Conversion rate = 1000
- [ ] Invalid conversion (kg → L): Show error
- [ ] Cost calculation with conversion
- [ ] Cost calculation with wastage
- [ ] Cost preview updates on quantity change
- [ ] Cost preview updates on unit change

## 🚀 Deployment Steps

1. **Backend**: Deploy unit conversion logic
2. **Frontend**: Deploy composable and updated views
3. **Test**: Verify conversion calculations
4. **Document**: Update user guide with examples
5. **Train**: Show users how to use conversion feature

## 📚 User Documentation

Tạo help section trong app:

```markdown
# Hướng dẫn Quy đổi Đơn vị

## Khi nào cần quy đổi?

Khi đơn vị lưu kho khác với đơn vị trong công thức:
- Kho: Kilogram (kg) → Công thức: gram (g)
- Kho: Lít (L) → Công thức: Milliliter (ml)

## Hệ thống tự động quy đổi

Bạn chỉ cần:
1. Chọn nguyên liệu
2. Chọn đơn vị công thức
3. Nhập số lượng
4. Hệ thống tự tính chi phí!

## Ví dụ

**Cà phê sữa đá**:
- Cà phê: 20g (kho: kg)
- Sữa: 100ml (kho: L)
- Đường: 10g (kho: kg)

Hệ thống tự động quy đổi và tính chi phí chính xác!
```

Đây là implementation guide đầy đủ với code examples cụ thể!