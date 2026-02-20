# Phân tích: Cho phép chuyển đổi đơn vị cho Batch Ingredients

## 1. Vấn đề hiện tại

Khi thêm batch ingredient vào menu, hệ thống hiện tại:
- Chỉ cho phép sử dụng đơn vị gốc của batch (ví dụ: nếu batch có đơn vị "L", chỉ có thể dùng "L")
- Không cho phép chuyển đổi sang đơn vị khác (ví dụ: từ "L" sang "ml")

### Code hiện tại (MenuView.vue, dòng 1280):
```javascript
compatibleUnits: [batch.unit], // Batches typically use single unit
```

## 2. Yêu cầu mới

Cho phép chuyển đổi đơn vị cho batch ingredients, ví dụ:
- Batch có đơn vị "L" → Có thể chọn "ml" khi pha chế
- Batch có đơn vị "kg" → Có thể chọn "g" khi pha chế

## 3. Giải pháp đề xuất

### 3.1. Sử dụng `getCompatibleUnits()` cho batch

Thay vì hardcode `[batch.unit]`, sử dụng function `getCompatibleUnits()` đã có sẵn trong `useUnitConversion` composable.

### 3.2. Thay đổi code

**File: `frontend/src/views/MenuView.vue`**

#### Thay đổi 1: Function `selectBatch()` - Variant ingredients (dòng ~1280)

**Trước:**
```javascript
variant.ingredients.push({
  id: batch.id,
  batch_definition_id: batch.id,
  type: 'batch',
  name: batch.name,
  quantity: 1,
  unit: batch.unit,
  stockUnit: batch.unit,
  compatibleUnits: [batch.unit], // ❌ Chỉ cho phép 1 đơn vị
  costPerUnit: batchCostPerUnit,
  wastage: 0,
  conversionRate: 1,
  estimatedCost: batchCostPerUnit
})
```

**Sau:**
```javascript
// Get compatible units for batch
const compatibleUnits = getCompatibleUnits(batch.unit)

variant.ingredients.push({
  id: batch.id,
  batch_definition_id: batch.id,
  type: 'batch',
  name: batch.name,
  quantity: 1,
  unit: batch.unit,
  stockUnit: batch.unit,
  compatibleUnits: compatibleUnits, // ✅ Cho phép chuyển đổi đơn vị
  costPerUnit: batchCostPerUnit,
  wastage: 0,
  conversionRate: 1,
  estimatedCost: batchCostPerUnit
})
```

#### Thay đổi 2: Function `selectBatch()` - Single-size ingredients (dòng ~1310)

**Trước:**
```javascript
form.value.ingredients.push({
  id: batch.id,
  batch_definition_id: batch.id,
  type: 'batch',
  name: batch.name,
  quantity: 1,
  unit: batch.unit,
  stockUnit: batch.unit,
  compatibleUnits: [batch.unit], // ❌ Chỉ cho phép 1 đơn vị
  costPerUnit: batchCostPerUnit,
  wastage: 0,
  conversionRate: 1,
  estimatedCost: batchCostPerUnit
})
```

**Sau:**
```javascript
// Get compatible units for batch
const compatibleUnits = getCompatibleUnits(batch.unit)

form.value.ingredients.push({
  id: batch.id,
  batch_definition_id: batch.id,
  type: 'batch',
  name: batch.name,
  quantity: 1,
  unit: batch.unit,
  stockUnit: batch.unit,
  compatibleUnits: compatibleUnits, // ✅ Cho phép chuyển đổi đơn vị
  costPerUnit: batchCostPerUnit,
  wastage: 0,
  conversionRate: 1,
  estimatedCost: batchCostPerUnit
})
```

## 4. Logic chuyển đổi đơn vị

### 4.1. `getCompatibleUnits()` đã có sẵn

Function này đã được implement trong `useUnitConversion.js`:

```javascript
const getCompatibleUnits = (stockUnit) => {
  if (massUnits.includes(stockUnit)) return massUnits      // ['kg', 'g']
  if (volumeUnits.includes(stockUnit)) return volumeUnits  // ['L', 'ml']
  if (countUnits.includes(stockUnit)) return countUnits    // ['piece', 'box', 'pack']
  return [stockUnit]
}
```

### 4.2. Ví dụ hoạt động

**Batch có đơn vị "L":**
- `getCompatibleUnits("L")` → `["L", "ml"]`
- User có thể chọn "L" hoặc "ml" trong dropdown

**Batch có đơn vị "kg":**
- `getCompatibleUnits("kg")` → `["kg", "g"]`
- User có thể chọn "kg" hoặc "g" trong dropdown

### 4.3. Tính toán chi phí tự động

Khi user thay đổi đơn vị, function `updateVariantRecipeUnit()` hoặc `updateRecipeUnit()` sẽ:
1. Tính conversion rate mới: `getConversionRate(stockUnit, newUnit)`
2. Cập nhật `ing.conversionRate`
3. Tính lại chi phí: `updateVariantIngredientCost()` hoặc `updateIngredientCost()`

**Ví dụ:**
- Batch: 1L = 50,000 VNĐ
- User chọn 200ml cho recipe
- Conversion: `getConversionRate("L", "ml")` = 0.001
- Chi phí: 200ml × 0.001 × 50,000 = 10,000 VNĐ

## 5. Kiểm tra validation

### 5.1. Frontend validation

Function `isValidConversion()` đã có sẵn để kiểm tra:
```javascript
const isValidConversion = (stockUnit, recipeUnit) => {
  if (stockUnit === recipeUnit) return true
  
  const stockIsMass = massUnits.includes(stockUnit)
  const recipeIsMass = massUnits.includes(recipeUnit)
  const stockIsVolume = volumeUnits.includes(stockUnit)
  const recipeIsVolume = volumeUnits.includes(recipeUnit)
  
  // Valid if both are same category
  return (stockIsMass && recipeIsMass) || (stockIsVolume && recipeIsVolume)
}
```

### 5.2. Backend validation

Backend không cần thay đổi vì:
- Backend chỉ lưu `quantity` và `unit` trong recipe
- Backend không validate conversion rate
- Conversion chỉ dùng để tính chi phí ở frontend

## 6. Testing scenarios

### Scenario 1: Batch "Sữa tươi" - 1L
1. Tạo batch definition: "Sữa tươi", unit = "L"
2. Thêm vào menu item
3. Kiểm tra dropdown đơn vị: phải có ["L", "ml"]
4. Chọn "ml", nhập 200ml
5. Kiểm tra chi phí tính đúng

### Scenario 2: Batch "Bột cà phê" - 1kg
1. Tạo batch definition: "Bột cà phê", unit = "kg"
2. Thêm vào menu item
3. Kiểm tra dropdown đơn vị: phải có ["kg", "g"]
4. Chọn "g", nhập 18g
5. Kiểm tra chi phí tính đúng

### Scenario 3: Batch "Syrup" - piece
1. Tạo batch definition: "Syrup", unit = "piece"
2. Thêm vào menu item
3. Kiểm tra dropdown đơn vị: chỉ có ["piece"] (không có conversion)
4. Nhập 1 piece
5. Kiểm tra chi phí tính đúng

## 7. Lợi ích

### 7.1. Linh hoạt hơn
- Batch có thể được đo bằng đơn vị lớn (L, kg) để dễ quản lý
- Recipe có thể dùng đơn vị nhỏ (ml, g) để chính xác hơn

### 7.2. Tính toán chính xác
- Conversion rate tự động
- Chi phí được tính đúng dựa trên conversion

### 7.3. UX tốt hơn
- User không cần tính toán thủ công
- Hiển thị conversion explanation: "1L = 1000ml"

## 8. Implementation steps

1. ✅ Phân tích vấn đề và giải pháp (file này)
2. ⏳ Thay đổi code trong `selectBatch()` function (2 chỗ)
3. ⏳ Test với batch có đơn vị "L" → chọn "ml"
4. ⏳ Test với batch có đơn vị "kg" → chọn "g"
5. ⏳ Verify chi phí tính đúng
6. ⏳ Test với multi-size items
7. ⏳ Test với single-size items

## 9. Code changes summary

**File cần sửa:** `frontend/src/views/MenuView.vue`

**Số dòng cần sửa:** 2 chỗ (dòng ~1280 và ~1310)

**Thay đổi:**
```javascript
// Thay vì:
compatibleUnits: [batch.unit],

// Dùng:
const compatibleUnits = getCompatibleUnits(batch.unit)
// ...
compatibleUnits: compatibleUnits,
```

**Không cần thay đổi:**
- Backend code
- Database schema
- API endpoints
- Other frontend files

## 10. Kết luận

Giải pháp rất đơn giản: chỉ cần sử dụng function `getCompatibleUnits()` đã có sẵn thay vì hardcode `[batch.unit]`. Tất cả logic conversion và cost calculation đã được implement sẵn, chỉ cần "mở khóa" tính năng này cho batch ingredients.
