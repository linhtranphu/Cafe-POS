# Batch Record Service - Unit Conversion Fix

## Vấn Đề
Lỗi "insufficient cafe: need 21.00g, have 1.00kg" khi tạo batch do backend không quy đổi đơn vị trước khi so sánh và trừ nguyên liệu.

## Nguyên Nhân
Backend đang so sánh trực tiếp giữa:
- `ing.Quantity` (số lượng trong kho, đơn vị: `ing.Unit` - ví dụ: kg)
- `ingredientCost.Quantity` (số lượng cần, đơn vị: `ingredientCost.Unit` - ví dụ: g)

Hai đơn vị này có thể khác nhau, dẫn đến so sánh sai:
```
1kg < 21g → FALSE (vì so sánh số: 1 < 21)
```

Nhưng thực tế:
```
1kg = 1000g > 21g → TRUE (đủ nguyên liệu)
```

## Ví Dụ Lỗi

### Dữ Liệu
- Ingredient: Cafe
  - Stock quantity: 1kg (đơn vị kho: kg)
  - Cost per unit: 100,000đ/kg
- Recipe: Cần 21g (đơn vị công thức: g)

### Code Cũ (SAI)
```go
if ing.Quantity < ingredientCost.Quantity {
    // 1 < 21 → TRUE → Lỗi "insufficient"
    return error("insufficient cafe: need 21.00g, have 1.00kg")
}
```

### Code Mới (ĐÚNG)
```go
conversionRate := ingredient.GetConversionRate(ing.Unit, ingredientCost.Unit)
// conversionRate = 0.001 (g→kg)
quantityInStockUnit := ingredientCost.Quantity * conversionRate
// quantityInStockUnit = 21 * 0.001 = 0.021kg

if ing.Quantity < quantityInStockUnit {
    // 1kg < 0.021kg → FALSE → Đủ nguyên liệu ✅
}
```

## Giải Pháp

### 1. Sửa Phần Kiểm Tra Availability (Dòng 85-103)

**Trước:**
```go
if ing.Quantity < ingredientCost.Quantity {
    return nil, fmt.Errorf("insufficient %s: need %.2f%s, have %.2f%s",
        ingredientCost.IngredientName,
        ingredientCost.Quantity,
        ingredientCost.Unit,
        ing.Quantity,
        ing.Unit,
    )
}
```

**Sau:**
```go
// Convert ingredientCost.Quantity from source unit to stock unit for comparison
conversionRate := ingredient.GetConversionRate(ingredient.UnitType(ing.Unit), ingredient.UnitType(ingredientCost.Unit))
quantityInStockUnit := ingredientCost.Quantity * conversionRate

if ing.Quantity < quantityInStockUnit {
    return nil, fmt.Errorf("insufficient %s: need %.2f%s (%.2f%s in stock unit), have %.2f%s",
        ingredientCost.IngredientName,
        ingredientCost.Quantity,
        ingredientCost.Unit,
        quantityInStockUnit,
        ing.Unit,
        ing.Quantity,
        ing.Unit,
    )
}
```

### 2. Sửa Phần Deduction (Dòng 115-125)

**Trước:**
```go
beforeQty := ing.Quantity

// Deduct quantity
ing.Quantity -= ingredientCost.Quantity
if ing.Quantity < 0 {
    ing.Quantity = 0
}
afterQty := ing.Quantity
```

**Sau:**
```go
beforeQty := ing.Quantity

// Convert ingredientCost.Quantity from source unit to stock unit for deduction
conversionRate := ingredient.GetConversionRate(ingredient.UnitType(ing.Unit), ingredient.UnitType(ingredientCost.Unit))
quantityToDeduct := ingredientCost.Quantity * conversionRate

// Deduct quantity in stock unit
ing.Quantity -= quantityToDeduct
if ing.Quantity < 0 {
    ing.Quantity = 0
}
afterQty := ing.Quantity
```

### 3. Cập Nhật Stock History (Dòng 140)

**Trước:**
```go
Quantity: -ingredientCost.Quantity, // Negative for deduction
```

**Sau:**
```go
Quantity: -quantityToDeduct, // Negative for deduction, in stock unit
```

## Logic Quy Đổi

### GetConversionRate Function
```go
// GetConversionRate(stockUnit, recipeUnit) returns multiplier
// Example: GetConversionRate("kg", "g") = 0.001
// Meaning: 1g = 0.001kg
```

### Conversion Examples

#### Case 1: g → kg
```
Stock: 1kg
Recipe: 21g
conversionRate = GetConversionRate("kg", "g") = 0.001
quantityInStockUnit = 21 * 0.001 = 0.021kg
Check: 1kg >= 0.021kg ✅
Deduct: 1kg - 0.021kg = 0.979kg
```

#### Case 2: ml → L
```
Stock: 2L
Recipe: 500ml
conversionRate = GetConversionRate("L", "ml") = 0.001
quantityInStockUnit = 500 * 0.001 = 0.5L
Check: 2L >= 0.5L ✅
Deduct: 2L - 0.5L = 1.5L
```

#### Case 3: Same Unit (kg → kg)
```
Stock: 5kg
Recipe: 2kg
conversionRate = GetConversionRate("kg", "kg") = 1.0
quantityInStockUnit = 2 * 1.0 = 2kg
Check: 5kg >= 2kg ✅
Deduct: 5kg - 2kg = 3kg
```

## Lợi Ích

1. **Kiểm tra chính xác**: So sánh đúng đơn vị trước khi tạo batch
2. **Trừ đúng số lượng**: Deduct đúng số lượng trong đơn vị kho
3. **Stock history chính xác**: Ghi nhận đúng số lượng thay đổi
4. **Hỗ trợ đa đơn vị**: Cho phép recipe dùng đơn vị khác với kho

## Testing Scenarios

### Scenario 1: g → kg Conversion
```
Given:
  - Ingredient: Cafe, 1kg in stock, unit: kg
  - Recipe: Need 21g, unit: g
When: Create batch
Then:
  - Check: 1kg >= 0.021kg ✅
  - Deduct: 0.021kg from stock
  - Result: 0.979kg remaining
```

### Scenario 2: ml → L Conversion
```
Given:
  - Ingredient: Milk, 2L in stock, unit: L
  - Recipe: Need 500ml, unit: ml
When: Create batch
Then:
  - Check: 2L >= 0.5L ✅
  - Deduct: 0.5L from stock
  - Result: 1.5L remaining
```

### Scenario 3: Insufficient Stock
```
Given:
  - Ingredient: Sugar, 0.01kg in stock, unit: kg
  - Recipe: Need 50g, unit: g
When: Create batch
Then:
  - Check: 0.01kg < 0.05kg ❌
  - Error: "insufficient Sugar: need 50.00g (0.05kg in stock unit), have 0.01kg"
```

## Files Modified
- `backend/application/services/batch_record_service.go`

## Dependencies
- Uses existing `ingredient.GetConversionRate()` function
- No new dependencies added

## Status
✅ Fix hoàn tất
✅ Compile thành công
🧪 Sẵn sàng test

## Related Issues
- Cùng vấn đề đã được fix trong frontend (BatchRecordForm cost calculation)
- Backend và frontend giờ đều áp dụng unit conversion đúng cách
