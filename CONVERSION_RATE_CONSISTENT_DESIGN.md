# Conversion Rate & Wastage - Consistent Design

## 🎯 Use Case: Sữa (Lít → Milliliter)

### Scenario
```
Kho: Mua sữa 10 lít @ 50,000 VND/lít
Công thức: Cà phê sữa đá cần 150ml sữa
Wastage: 5% (đổ tràn, hỏng)

Câu hỏi: Chi phí sữa cho 1 ly cà phê sữa đá là bao nhiêu?
```

### Phân tích Logic

#### 1. Conversion Rate
```
Kho: Lít (L)
Công thức: Milliliter (ml)
Quy đổi: 1L = 1000ml → 1ml = 0.001L

Conversion Rate = 0.001
```

**Ý nghĩa**: Để tính cost, ta cần quy đổi 150ml về lít
```
150ml × 0.001 = 0.15L
```

#### 2. Cost Calculation
```
Base Cost = Quantity × Cost_Per_Unit × Conversion_Rate
          = 150ml × 50,000 VND/L × 0.001
          = 0.15L × 50,000 VND/L
          = 7,500 VND
```

#### 3. Wastage Percentage
```
Wastage = 5% (do đổ tràn, hỏng trong quá trình pha chế)

Final Cost = Base Cost × (1 + Wastage/100)
           = 7,500 × (1 + 5/100)
           = 7,500 × 1.05
           = 7,875 VND
```

### Công thức tổng quát
```
Cost = Quantity × Cost_Per_Unit × Conversion_Rate × (1 + Wastage/100)
     = 150 × 50,000 × 0.001 × 1.05
     = 7,875 VND
```

## 🏗️ Consistent Design Across Modules

### Module 1: Ingredient Management

#### Data Model
```go
// backend/domain/ingredient/ingredient.go
type Ingredient struct {
    ID                primitive.ObjectID `json:"id"`
    Name              string             `json:"name"`
    Unit              UnitType           `json:"unit"`              // "L" (stock unit)
    CostPerUnit       float64            `json:"cost_per_unit"`     // 50,000 VND/L
    ConversionRate    float64            `json:"conversion_rate"`   // NOT STORED HERE
    WastagePercentage float64            `json:"wastage_percentage"` // 5.0
    CurrentStock      float64            `json:"current_stock"`     // 10 L
}
```

**Key Decision**: `conversion_rate` KHÔNG lưu trong Ingredient
- ❌ Sai: Lưu conversion_rate = 0.001 trong ingredient
- ✅ Đúng: Tính conversion_rate động dựa trên stock_unit và recipe_unit

**Lý do**:
1. Conversion rate phụ thuộc vào recipe unit (ml, L, cl...)
2. Cùng 1 ingredient có thể dùng nhiều recipe units khác nhau
3. Tránh inconsistency khi có nhiều món dùng cùng ingredient với units khác nhau

#### UI - Ingredient Form
```vue
<template>
  <div class="ingredient-form">
    <!-- Stock Unit (fixed for inventory) -->
    <div>
      <label>Đơn vị kho *</label>
      <select v-model="form.unit">
        <option value="L">Lít (L)</option>
        <option value="ml">Milliliter (ml)</option>
        <option value="kg">Kilogram (kg)</option>
        <option value="g">Gram (g)</option>
      </select>
    </div>
    
    <!-- Cost Per Unit -->
    <div>
      <label>Giá vốn *</label>
      <input v-model.number="form.cost_per_unit" type="number" />
      <span>VND/{{ form.unit }}</span>
    </div>
    
    <!-- Wastage Percentage -->
    <div>
      <label>Hao hụt (%) - Tùy chọn</label>
      <input v-model.number="form.wastage_percentage" type="number" min="0" max="100" step="0.1" />
      <p class="hint">VD: 5% cho đồ lỏng (đổ tràn), 10% cho rau củ (cắt tỉa)</p>
    </div>
    
    <!-- Info Box -->
    <div class="info-box">
      <h4>ℹ️ Quy đổi đơn vị</h4>
      <p>Hệ thống sẽ tự động quy đổi khi bạn tạo công thức món.</p>
      <p>Bạn chỉ cần chọn đơn vị kho và giá vốn.</p>
    </div>
  </div>
</template>
```

### Module 2: Menu Management

#### Data Model
```go
// backend/domain/menu/menu.go
type MenuItem struct {
    ID          primitive.ObjectID `json:"id"`
    Name        string             `json:"name"`
    Price       float64            `json:"price"`
    Ingredients []MenuIngredient   `json:"ingredients"`
    CurrentCost float64            `json:"current_cost"` // Calculated
}

type MenuIngredient struct {
    Name     string   `json:"name"`      // "Sữa tươi"
    Quantity float64  `json:"quantity"`  // 150
    Unit     UnitType `json:"unit"`      // "ml" (recipe unit)
}
```

**Key Decision**: Lưu `unit` trong MenuIngredient (recipe unit)
- ✅ Đúng: Mỗi món có thể dùng unit khác nhau cho cùng ingredient
- ✅ Flexible: Món A dùng "ml", Món B dùng "L" cho cùng sữa

#### UI - Menu Form (Ingredient Selector)
```vue
<template>
  <div class="ingredient-selector">
    <!-- Selected Ingredient Display -->
    <div v-for="(ing, index) in form.ingredients" :key="index" class="ingredient-item">
      <!-- Header -->
      <div class="header">
        <div class="info">
          <span class="name">{{ ing.name }}</span>
          <span class="stock-info">Kho: {{ ing.stockUnit }} @ {{ formatPrice(ing.costPerUnit) }}/{{ ing.stockUnit }}</span>
        </div>
        <button @click="removeIngredient(index)">×</button>
      </div>
      
      <!-- Recipe Unit Selector -->
      <div class="recipe-unit">
        <label>Đơn vị công thức:</label>
        <select v-model="ing.unit" @change="updateConversion(index)">
          <option v-for="unit in ing.compatibleUnits" :key="unit" :value="unit">
            {{ unit }}
          </option>
        </select>
      </div>
      
      <!-- Quantity Input -->
      <div class="quantity">
        <label>Số lượng:</label>
        <input v-model.number="ing.quantity" @input="updateCost(index)" type="number" step="0.1" />
        <span>{{ ing.unit }}</span>
      </div>
      
      <!-- Conversion Info (Auto-calculated, Read-only) -->
      <div v-if="ing.conversionRate !== 1" class="conversion-info">
        <span class="icon">ℹ️</span>
        <span class="text">Quy đổi: 1{{ ing.unit }} = {{ 1/ing.conversionRate }}{{ ing.stockUnit }}</span>
      </div>
      
      <!-- Cost Preview -->
      <div class="cost-preview">
        <div class="row">
          <span>Chi phí cơ bản:</span>
          <span>{{ formatPrice(ing.baseCost) }}</span>
        </div>
        <div v-if="ing.wastage > 0" class="row wastage">
          <span>Hao hụt ({{ ing.wastage }}%):</span>
          <span>+{{ formatPrice(ing.wastageCost) }}</span>
        </div>
        <div class="row total">
          <span>Tổng chi phí:</span>
          <span class="amount">{{ formatPrice(ing.totalCost) }}</span>
        </div>
      </div>
    </div>
    
    <!-- Total Cost Summary -->
    <div class="total-summary">
      <span>Tổng chi phí nguyên liệu:</span>
      <span class="amount">{{ formatPrice(totalIngredientCost) }}</span>
    </div>
  </div>
</template>

<script setup>
import { useUnitConversion } from '@/composables/useUnitConversion'

const { getConversionRate, getCompatibleUnits, isValidConversion } = useUnitConversion()

// When selecting ingredient from list
const selectIngredient = (ingredient) => {
  const compatibleUnits = getCompatibleUnits(ingredient.unit)
  const defaultRecipeUnit = ingredient.unit // Start with same unit
  
  form.value.ingredients.push({
    id: ingredient.id,
    name: ingredient.name,
    quantity: 1,
    unit: defaultRecipeUnit,           // Recipe unit (can be changed)
    stockUnit: ingredient.unit,        // Stock unit (fixed)
    compatibleUnits: compatibleUnits,  // Available units for dropdown
    costPerUnit: ingredient.cost_per_unit,
    wastage: ingredient.wastage_percentage || 0,
    conversionRate: 1.0,               // Initially 1.0 (same unit)
    baseCost: 0,
    wastageCost: 0,
    totalCost: 0
  })
  
  updateCost(form.value.ingredients.length - 1)
}

// Update conversion rate when recipe unit changes
const updateConversion = (index) => {
  const ing = form.value.ingredients[index]
  
  // Validate conversion
  if (!isValidConversion(ing.stockUnit, ing.unit)) {
    alert(`Không thể quy đổi từ ${ing.stockUnit} sang ${ing.unit}!`)
    ing.unit = ing.stockUnit
    return
  }
  
  // Calculate new conversion rate
  ing.conversionRate = getConversionRate(ing.stockUnit, ing.unit)
  
  // Recalculate cost
  updateCost(index)
}

// Calculate cost for ingredient
const updateCost = (index) => {
  const ing = form.value.ingredients[index]
  
  // Base cost = quantity × cost_per_unit × conversion_rate
  ing.baseCost = ing.quantity * ing.costPerUnit * ing.conversionRate
  ing.baseCost = Math.round(ing.baseCost * 100) / 100
  
  // Wastage cost = base_cost × (wastage / 100)
  ing.wastageCost = ing.baseCost * (ing.wastage / 100)
  ing.wastageCost = Math.round(ing.wastageCost * 100) / 100
  
  // Total cost = base_cost + wastage_cost
  ing.totalCost = ing.baseCost + ing.wastageCost
  ing.totalCost = Math.round(ing.totalCost * 100) / 100
}

// Total cost for all ingredients
const totalIngredientCost = computed(() => {
  return form.value.ingredients.reduce((sum, ing) => sum + ing.totalCost, 0)
})
</script>
```

### Module 3: Cost Calculator Service

#### Backend Implementation
```go
// backend/application/services/cost_calculator_service.go

func (s *CostCalculatorService) CalculateMenuItemCost(ctx context.Context, menuItemID primitive.ObjectID) (*MenuItemCostResult, error) {
    // 1. Fetch menu item
    menuItem, err := s.menuRepo.FindByID(ctx, menuItemID)
    if err != nil {
        return nil, err
    }
    
    // 2. Fetch all ingredients (for lookup)
    allIngredients, err := s.ingredientRepo.FindAll(ctx)
    if err != nil {
        return nil, err
    }
    
    // Build lookup map
    ingredientMap := make(map[string]*ingredient.Ingredient)
    for _, ing := range allIngredients {
        ingredientMap[ing.Name] = ing
    }
    
    // 3. Calculate cost for each menu ingredient
    var totalCost float64
    var missingIngredients []string
    
    for _, menuIng := range menuItem.Ingredients {
        // Find ingredient in database
        dbIng, exists := ingredientMap[menuIng.Name]
        if !exists || dbIng.CostPerUnit <= 0 {
            missingIngredients = append(missingIngredients, menuIng.Name)
            continue
        }
        
        // Calculate conversion rate dynamically
        // stockUnit = dbIng.Unit (e.g., "L")
        // recipeUnit = menuIng.Unit (e.g., "ml")
        conversionRate := ingredient.GetConversionRate(dbIng.Unit, menuIng.Unit)
        
        // Get wastage (default 0 if not set)
        wastage := dbIng.WastagePercentage
        if wastage < 0 {
            wastage = 0.0
        }
        
        // Calculate cost
        // Formula: quantity × cost_per_unit × conversion_rate × (1 + wastage/100)
        cost := menuIng.Quantity * dbIng.CostPerUnit * conversionRate * (1 + wastage/100)
        totalCost += cost
    }
    
    // 4. Round to 2 decimal places
    totalCost = math.Round(totalCost * 100) / 100
    
    // 5. Return result
    return &MenuItemCostResult{
        MenuItemID:         menuItemID,
        CurrentCost:        totalCost,
        CostStatus:         determineCostStatus(missingIngredients),
        MissingIngredients: missingIngredients,
    }, nil
}
```

### Module 4: Cost Analysis & Reports

#### API Response Structure
```json
{
  "menu_item_id": "...",
  "menu_item_name": "Cà phê sữa đá",
  "current_cost": 15750,
  "ingredients": [
    {
      "name": "Cà phê",
      "quantity": 20,
      "recipe_unit": "g",
      "stock_unit": "kg",
      "cost_per_unit": 200000,
      "conversion_rate": 0.001,
      "wastage_percentage": 2.0,
      "base_cost": 4000,
      "wastage_cost": 80,
      "total_cost": 4080
    },
    {
      "name": "Sữa tươi",
      "quantity": 150,
      "recipe_unit": "ml",
      "stock_unit": "L",
      "cost_per_unit": 50000,
      "conversion_rate": 0.001,
      "wastage_percentage": 5.0,
      "base_cost": 7500,
      "wastage_cost": 375,
      "total_cost": 7875
    },
    {
      "name": "Đường",
      "quantity": 15,
      "recipe_unit": "g",
      "stock_unit": "kg",
      "cost_per_unit": 25000,
      "conversion_rate": 0.001,
      "wastage_percentage": 0.0,
      "base_cost": 375,
      "wastage_cost": 0,
      "total_cost": 375
    }
  ],
  "total_base_cost": 11875,
  "total_wastage_cost": 455,
  "total_cost": 12330
}
```

#### Frontend Display
```vue
<template>
  <div class="cost-breakdown">
    <h3>Chi phí nguyên liệu - {{ menuItem.name }}</h3>
    
    <table class="ingredient-table">
      <thead>
        <tr>
          <th>Nguyên liệu</th>
          <th>Số lượng</th>
          <th>Đơn giá</th>
          <th>Quy đổi</th>
          <th>Hao hụt</th>
          <th>Chi phí</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="ing in costDetail.ingredients" :key="ing.name">
          <td>{{ ing.name }}</td>
          <td>{{ ing.quantity }} {{ ing.recipe_unit }}</td>
          <td>{{ formatPrice(ing.cost_per_unit) }}/{{ ing.stock_unit }}</td>
          <td>
            <span v-if="ing.conversion_rate !== 1" class="badge">
              {{ ing.conversion_rate }}
            </span>
            <span v-else>-</span>
          </td>
          <td>
            <span v-if="ing.wastage_percentage > 0" class="badge warning">
              {{ ing.wastage_percentage }}%
            </span>
            <span v-else>-</span>
          </td>
          <td class="cost">
            <div class="base">{{ formatPrice(ing.base_cost) }}</div>
            <div v-if="ing.wastage_cost > 0" class="wastage">
              +{{ formatPrice(ing.wastage_cost) }}
            </div>
            <div class="total">{{ formatPrice(ing.total_cost) }}</div>
          </td>
        </tr>
      </tbody>
      <tfoot>
        <tr class="summary">
          <td colspan="5">Tổng chi phí nguyên liệu:</td>
          <td class="total-cost">{{ formatPrice(costDetail.total_cost) }}</td>
        </tr>
      </tfoot>
    </table>
  </div>
</template>
```

## 🔄 Data Flow - Complete Picture

```
┌─────────────────────────────────────────────────────────────┐
│ 1. INGREDIENT MANAGEMENT                                    │
├─────────────────────────────────────────────────────────────┤
│ Input:                                                      │
│   - Name: "Sữa tươi"                                        │
│   - Stock Unit: "L"                                         │
│   - Cost Per Unit: 50,000 VND/L                            │
│   - Wastage: 5%                                             │
│   - Current Stock: 10 L                                     │
│                                                             │
│ Stored in DB:                                               │
│   {                                                         │
│     "name": "Sữa tươi",                                     │
│     "unit": "L",                                            │
│     "cost_per_unit": 50000,                                 │
│     "wastage_percentage": 5.0,                              │
│     "current_stock": 10                                     │
│   }                                                         │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. MENU CREATION                                            │
├─────────────────────────────────────────────────────────────┤
│ User Action:                                                │
│   1. Select ingredient: "Sữa tươi"                          │
│   2. Choose recipe unit: "ml" (from dropdown)               │
│   3. Enter quantity: 150                                    │
│                                                             │
│ Frontend Calculation:                                       │
│   - Stock Unit: "L" (from ingredient)                       │
│   - Recipe Unit: "ml" (user selected)                       │
│   - Conversion Rate: 0.001 (auto-calculated)               │
│   - Base Cost: 150 × 50,000 × 0.001 = 7,500               │
│   - Wastage Cost: 7,500 × 0.05 = 375                       │
│   - Total Cost: 7,875 VND                                   │
│                                                             │
│ Stored in DB:                                               │
│   {                                                         │
│     "name": "Cà phê sữa đá",                                │
│     "ingredients": [                                        │
│       {                                                     │
│         "name": "Sữa tươi",                                 │
│         "quantity": 150,                                    │
│         "unit": "ml"                                        │
│       }                                                     │
│     ]                                                       │
│   }                                                         │
│                                                             │
│ Note: conversion_rate NOT stored, calculated dynamically    │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. COST CALCULATION (Backend)                               │
├─────────────────────────────────────────────────────────────┤
│ When calculating cost:                                      │
│   1. Fetch menu item: "Cà phê sữa đá"                       │
│   2. Fetch ingredient: "Sữa tươi"                           │
│   3. Calculate conversion rate:                             │
│      GetConversionRate("L", "ml") = 0.001                   │
│   4. Calculate cost:                                        │
│      150 × 50,000 × 0.001 × 1.05 = 7,875 VND              │
│   5. Return result                                          │
│                                                             │
│ Result:                                                     │
│   {                                                         │
│     "menu_item_id": "...",                                  │
│     "current_cost": 7875,                                   │
│     "ingredients": [                                        │
│       {                                                     │
│         "name": "Sữa tươi",                                 │
│         "quantity": 150,                                    │
│         "recipe_unit": "ml",                                │
│         "stock_unit": "L",                                  │
│         "conversion_rate": 0.001,                           │
│         "wastage_percentage": 5.0,                          │
│         "total_cost": 7875                                  │
│       }                                                     │
│     ]                                                       │
│   }                                                         │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. COST ANALYSIS & REPORTS                                  │
├─────────────────────────────────────────────────────────────┤
│ Display:                                                    │
│   - Ingredient breakdown with conversion info               │
│   - Base cost vs wastage cost                               │
│   - Total cost per menu item                                │
│   - Profit margin analysis                                  │
└─────────────────────────────────────────────────────────────┘
```

## ✅ Consistency Checklist

### Data Storage
- [x] Ingredient: Store `unit` (stock unit), `wastage_percentage`
- [x] Ingredient: DO NOT store `conversion_rate`
- [x] Menu Item: Store `unit` (recipe unit) per ingredient
- [x] Menu Item: DO NOT store `conversion_rate`

### Calculation
- [x] Calculate `conversion_rate` dynamically based on stock_unit and recipe_unit
- [x] Use same formula across all modules: `quantity × cost × conversion × (1 + wastage/100)`
- [x] Round to 2 decimal places consistently

### UI/UX
- [x] Show conversion rate when != 1.0
- [x] Show wastage cost separately for transparency
- [x] Display cost preview in real-time
- [x] Validate unit compatibility

### API
- [x] Return conversion_rate in cost detail responses
- [x] Include both stock_unit and recipe_unit in responses
- [x] Separate base_cost and wastage_cost in breakdown

## 🎯 Implementation Priority

### Phase 1: Core (Must Have)
1. ✅ Add `wastage_percentage` to Ingredient model
2. ✅ Add `unit` to MenuIngredient model
3. ✅ Implement `GetConversionRate()` function
4. ✅ Update cost calculation formula

### Phase 2: UI (Should Have)
1. ✅ Recipe unit selector in menu form
2. ✅ Cost preview with breakdown
3. ✅ Conversion rate display
4. ✅ Wastage percentage input

### Phase 3: Enhancement (Nice to Have)
1. ⭐ Auto-suggest compatible units
2. ⭐ Validation for unit compatibility
3. ⭐ Cost comparison (with/without wastage)
4. ⭐ Historical cost tracking

Đây là design hoàn chỉnh và consistent cho toàn bộ hệ thống!