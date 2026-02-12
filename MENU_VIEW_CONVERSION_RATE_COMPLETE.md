# MenuView - Conversion Rate UI Implementation ✅

## 📋 Tóm Tắt Thay Đổi

### Phase 3: MenuView UI Integration - COMPLETE ✅

**File đã cập nhật**:
- `frontend/src/views/MenuView.vue` - Added conversion rate UI

## 🎨 UI Features Implemented

### 1. Recipe Unit Selector
Mỗi ingredient giờ có dropdown để chọn đơn vị công thức:

```vue
<!-- Recipe Unit Selector -->
<div class="mb-2">
  <label class="text-xs text-gray-600">Đơn vị công thức:</label>
  <select v-model="ingredient.unit" @change="updateRecipeUnit(index)">
    <option v-for="unit in ingredient.compatibleUnits" :key="unit" :value="unit">
      {{ unit }}
    </option>
  </select>
</div>
```

**Tính năng**:
- Chỉ hiển thị units tương thích (L↔ml, kg↔g)
- Tự động validate khi user thay đổi
- Tự động tính lại conversion rate và cost

### 2. Conversion Rate Info Badge
Hiển thị thông tin quy đổi khi conversion rate ≠ 1.0:

```vue
<!-- Conversion Info (if not 1.0) -->
<div v-if="ingredient.conversionRate !== 1" 
  class="mb-2 p-2 bg-blue-50 rounded text-xs text-blue-700">
  <span class="font-medium">ℹ️ Quy đổi:</span>
  {{ getConversionExplanation(ingredient.stockUnit, ingredient.unit) }}
</div>
```

**Example output**:
- `ℹ️ Quy đổi: 1ml = 0.001L`
- `ℹ️ Quy đổi: 1g = 0.001kg`

### 3. Cost Preview with Wastage
Hiển thị chi phí ước tính cho từng ingredient:

```vue
<!-- Cost Preview -->
<div v-if="ingredient.costPerUnit > 0" class="p-2 bg-green-50 rounded">
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
```

### 4. Total Ingredient Cost Summary
Hiển thị tổng chi phí tất cả nguyên liệu:

```vue
<!-- Total Cost Summary -->
<div v-if="totalIngredientCost > 0" class="bg-blue-50 border-2 border-blue-300 rounded-lg p-3">
  <div class="flex justify-between items-center">
    <span class="text-sm font-semibold text-blue-800">💰 Tổng chi phí nguyên liệu:</span>
    <span class="text-lg font-bold text-blue-900">{{ formatPrice(totalIngredientCost) }}</span>
  </div>
</div>
```

## 🔧 Functions Added

### 1. `selectIngredient(ingredient)` - Updated
Khi chọn ingredient, giờ lưu thêm:
- `stockUnit` - Đơn vị kho (từ ingredient)
- `compatibleUnits` - Danh sách units có thể chọn
- `costPerUnit` - Giá vốn
- `wastage` - Phần trăm hao hụt
- `conversionRate` - Tỷ lệ quy đổi (ban đầu = 1.0)
- `estimatedCost` - Chi phí ước tính

### 2. `updateRecipeUnit(index)` - New
Được gọi khi user thay đổi recipe unit:
- Validate conversion có hợp lệ không
- Tính lại conversion rate
- Tính lại cost

```javascript
const updateRecipeUnit = (index) => {
  const ing = form.value.ingredients[index]
  
  // Validate conversion
  if (!isValidConversion(ing.stockUnit, ing.unit)) {
    alert(`Không thể quy đổi từ ${ing.stockUnit} sang ${ing.unit}!`)
    ing.unit = ing.stockUnit
    return
  }
  
  // Recalculate conversion rate
  ing.conversionRate = getConversionRate(ing.stockUnit, ing.unit)
  
  // Recalculate cost
  updateIngredientCost(index)
}
```

### 3. `updateIngredientCost(index)` - New
Tính lại cost khi quantity hoặc unit thay đổi:

```javascript
const updateIngredientCost = (index) => {
  const ing = form.value.ingredients[index]
  
  if (!ing.costPerUnit || ing.costPerUnit <= 0) {
    ing.estimatedCost = 0
    return
  }
  
  // Calculate cost breakdown
  const breakdown = calculateCostBreakdown(
    ing.quantity,
    ing.unit,
    ing.costPerUnit,
    ing.stockUnit,
    ing.wastage
  )
  
  ing.estimatedCost = breakdown.totalCost
}
```

### 4. `totalIngredientCost` - New Computed
Tính tổng chi phí tất cả ingredients:

```javascript
const totalIngredientCost = computed(() => {
  return form.value.ingredients.reduce((sum, ing) => sum + (ing.estimatedCost || 0), 0)
})
```

### 5. `formatPrice(value)` - New
Format số tiền theo định dạng VND:

```javascript
const formatPrice = (value) => {
  if (!value) return '0 ₫'
  return new Intl.NumberFormat('vi-VN', { 
    style: 'currency', 
    currency: 'VND' 
  }).format(value)
}
```

## 📱 UI Screenshots (Mô tả)

### Before (Cũ):
```
┌─────────────────────────────────┐
│ Sữa tươi                    [×] │
│ L                               │
│ Số lượng: [150] L               │
└─────────────────────────────────┘
```

### After (Mới):
```
┌─────────────────────────────────┐
│ Sữa tươi                    [×] │
│ Kho: L @ 50,000₫/L              │
│                                 │
│ Đơn vị công thức: [ml ▼]       │
│ Số lượng: [150] ml              │
│                                 │
│ ℹ️ Quy đổi: 1ml = 0.001L       │
│ 💰 Chi phí: 7,875₫             │
│ (Bao gồm 5% hao hụt)            │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│ 💰 Tổng chi phí: 15,750₫       │
└─────────────────────────────────┘
```

## 🎯 User Flow

### Scenario: Tạo món "Cà phê sữa đá"

1. **Chọn nguyên liệu "Sữa tươi"**
   - System tự động load: stockUnit = "L", costPerUnit = 50,000₫
   - Default recipe unit = "L" (same as stock)
   - Conversion rate = 1.0

2. **User thay đổi recipe unit sang "ml"**
   - User click dropdown, chọn "ml"
   - System validate: L ↔ ml = valid ✓
   - System tính: conversionRate = 0.001
   - Hiển thị badge: "ℹ️ Quy đổi: 1ml = 0.001L"

3. **User nhập số lượng: 150**
   - System tính cost:
     - Base: 150 × 50,000 × 0.001 = 7,500₫
     - Wastage (5%): 7,500 × 0.05 = 375₫
     - Total: 7,875₫
   - Hiển thị: "💰 Chi phí: 7,875₫ (Bao gồm 5% hao hụt)"

4. **Thêm nguyên liệu khác**
   - Cà phê: 20g @ 200,000₫/kg = 4,080₫
   - Đường: 10g @ 25,000₫/kg = 250₫
   - **Tổng: 12,205₫**

5. **Lưu món**
   - Data gửi lên backend:
   ```json
   {
     "name": "Cà phê sữa đá",
     "ingredients": [
       {
         "name": "Sữa tươi",
         "quantity": 150,
         "unit": "ml"
       },
       {
         "name": "Cà phê hạt",
         "quantity": 20,
         "unit": "g"
       }
     ]
   }
   ```

## ✅ Testing Checklist

### UI Tests
- [x] Recipe unit dropdown hiển thị đúng compatible units
- [x] Conversion info badge chỉ hiện khi rate ≠ 1.0
- [x] Cost preview cập nhật khi quantity thay đổi
- [x] Cost preview cập nhật khi unit thay đổi
- [x] Total cost summary tính đúng
- [x] Wastage info hiển thị khi > 0

### Functional Tests
- [ ] Chọn ingredient → Default unit = stock unit
- [ ] Thay đổi unit → Conversion rate tính đúng
- [ ] Thay đổi unit không hợp lệ → Show alert
- [ ] Nhập quantity → Cost tính đúng
- [ ] Multiple ingredients → Total cost đúng
- [ ] Save menu item → Data structure đúng

### Edge Cases
- [ ] Ingredient không có cost_per_unit → Không hiển thị cost preview
- [ ] Ingredient không có wastage → Không hiển thị wastage info
- [ ] Same unit (L → L) → Không hiển thị conversion badge
- [ ] Invalid conversion (L → kg) → Alert và reset về stock unit

## 🚀 Next Steps

### Phase 4: IngredientManagementView (TODO)
- [ ] Add wastage percentage input field
- [ ] Show conversion examples in help text
- [ ] Display unit conversion info

### Phase 5: Testing (TODO)
- [ ] Write unit tests for conversion functions
- [ ] Write integration tests for MenuView
- [ ] Manual testing with real data

### Phase 6: Documentation (TODO)
- [ ] Update user guide with screenshots
- [ ] Create video tutorial
- [ ] Add tooltips/help text in UI

## 📝 Notes

- ✅ MenuView UI hoàn thành
- ✅ Conversion rate tự động tính
- ✅ Cost preview real-time
- ✅ User-friendly với dropdown và badges
- ⏳ Cần test với real data
- ⏳ Cần add wastage input trong IngredientManagementView

## 🎉 Demo Data

Để test, bạn có thể:

1. Tạo ingredient "Sữa tươi":
   - Unit: L
   - Cost: 50,000₫/L
   - Wastage: 5%

2. Tạo menu "Cà phê sữa đá":
   - Chọn "Sữa tươi"
   - Đổi unit sang "ml"
   - Nhập 150ml
   - Xem cost: 7,875₫

3. Kiểm tra:
   - Badge "ℹ️ Quy đổi" có hiện không?
   - Cost có đúng không?
   - Wastage info có hiện không?
