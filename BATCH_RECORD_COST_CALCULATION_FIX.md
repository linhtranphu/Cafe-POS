# Batch Record Form - Cost Calculation Fix

## Vấn Đề
Phần "Nguyên Liệu Cần Thiết" trong form ghi nhận batch tính toán sai chi phí dự kiến do không áp dụng quy đổi đơn vị.

## Nguyên Nhân
Code cũ tính chi phí trực tiếp bằng cách nhân số lượng với `cost_per_unit`, nhưng không xét đến việc:
- `source_unit` (đơn vị trong công thức) có thể khác với `ingredient.unit` (đơn vị trong kho)
- Cần phải quy đổi từ đơn vị công thức sang đơn vị kho trước khi nhân với giá

### Ví Dụ Lỗi
**Dữ liệu:**
- Nguyên liệu: Sữa tươi, giá 10,000đ/L (đơn vị kho: L)
- Công thức batch: cần 100ml (đơn vị công thức: ml)
- Wastage: 10%

**Tính toán SAI (code cũ):**
```
totalQuantity = 100 * 1.1 = 110ml
cost = 110 * 10,000 = 1,100,000đ ❌ SAI!
```

**Tính toán ĐÚNG (code mới):**
```
totalQuantity = 100 * 1.1 = 110ml
conversionRate = getConversionRate("L", "ml") = 0.001
quantityInStockUnit = 110 * 0.001 = 0.11L
cost = 0.11 * 10,000 = 1,100đ ✅ ĐÚNG!
```

## Giải Pháp

### 1. Import useUnitConversion
```javascript
import { useUnitConversion } from '../../composables/useUnitConversion'
const { getConversionRate } = useUnitConversion()
```

### 2. Cập Nhật expectedCost Computed Property
```javascript
const expectedCost = computed(() => {
  if (!selectedDefinition.value || batchCount.value <= 0) return 0
  
  let total = 0
  const def = selectedDefinition.value
  const quantity = totalOutput.value
  
  for (const rate of def.conversion_rates || []) {
    const ingredient = ingredients.value.find(i => i.id === rate.source_ingredient_id)
    if (ingredient && ingredient.cost_per_unit) {
      // Calculate quantity needed in source unit
      const ratio = quantity / rate.batch_quantity
      const baseQuantity = rate.source_quantity * ratio
      const wastageMultiplier = 1 + (rate.wastage_rate || 0)
      const totalQuantity = baseQuantity * wastageMultiplier
      
      // Apply unit conversion: convert from source_unit to ingredient stock unit
      const conversionRate = getConversionRate(ingredient.unit, rate.source_unit)
      const quantityInStockUnit = totalQuantity * conversionRate
      
      // Calculate cost using stock unit price
      total += quantityInStockUnit * ingredient.cost_per_unit
    }
  }
  
  return total
})
```

## Logic Tính Toán Chi Phí

### Bước 1: Tính Số Lượng Nguyên Liệu Cần
```javascript
ratio = totalOutput / batch_quantity
baseQuantity = source_quantity * ratio
wastageMultiplier = 1 + wastage_rate
totalQuantity = baseQuantity * wastageMultiplier
```

### Bước 2: Quy Đổi Sang Đơn Vị Kho
```javascript
conversionRate = getConversionRate(stockUnit, recipeUnit)
quantityInStockUnit = totalQuantity * conversionRate
```

### Bước 3: Tính Chi Phí
```javascript
cost = quantityInStockUnit * cost_per_unit
```

## Ví Dụ Cụ Thể

### Case 1: ml → L
```
Ingredient: Sữa tươi
- Stock unit: L
- Cost per unit: 10,000đ/L

Recipe: 100ml
Wastage: 10%

Calculation:
1. totalQuantity = 100 * 1.1 = 110ml
2. conversionRate = 0.001 (ml→L)
3. quantityInStockUnit = 110 * 0.001 = 0.11L
4. cost = 0.11 * 10,000 = 1,100đ ✅
```

### Case 2: g → kg
```
Ingredient: Đường
- Stock unit: kg
- Cost per unit: 20,000đ/kg

Recipe: 50g
Wastage: 5%

Calculation:
1. totalQuantity = 50 * 1.05 = 52.5g
2. conversionRate = 0.001 (g→kg)
3. quantityInStockUnit = 52.5 * 0.001 = 0.0525kg
4. cost = 0.0525 * 20,000 = 1,050đ ✅
```

### Case 3: Same Unit (L → L)
```
Ingredient: Nước lọc
- Stock unit: L
- Cost per unit: 5,000đ/L

Recipe: 2L
Wastage: 0%

Calculation:
1. totalQuantity = 2 * 1.0 = 2L
2. conversionRate = 1.0 (L→L, no conversion)
3. quantityInStockUnit = 2 * 1.0 = 2L
4. cost = 2 * 5,000 = 10,000đ ✅
```

## Files Modified
- `frontend/src/components/batch/BatchRecordForm.vue`

## Testing Checklist
- [ ] Test với nguyên liệu có đơn vị ml, công thức dùng L
- [ ] Test với nguyên liệu có đơn vị kg, công thức dùng g
- [ ] Test với nguyên liệu và công thức cùng đơn vị
- [ ] Test với wastage rate khác 0
- [ ] Test với nhiều nguyên liệu trong 1 batch
- [ ] Verify chi phí hiển thị khớp với tính toán thủ công

## Status
✅ Fix hoàn tất
✅ Không có lỗi diagnostic
🧪 Sẵn sàng test

## Related Issues
- Cùng lỗi đã được fix trong `BatchDefinitionForm.vue` trước đó
- Logic tính toán giống với menu cost calculation
