# Fix: Batch Unit Conversion - Đơn vị bị biến mất khi edit

## Vấn đề

Sau khi tạo món mới với batch ingredient và chọn đơn vị chuyển đổi (ví dụ: batch "L" nhưng chọn "ml"), khi click "Sửa" để xem lại món, dropdown đơn vị bị biến mất hoặc chỉ hiển thị 1 đơn vị.

## Nguyên nhân

Function `editItem()` chỉ xử lý raw ingredients, không xử lý batch ingredients. Khi load món đã lưu:
- Raw ingredients: tìm trong `availableIngredients` → có `compatibleUnits`
- Batch ingredients: không tìm thấy → trả về `ing` gốc → không có `compatibleUnits`

## Giải pháp đã thực hiện

### 1. Sửa function `editItem()` - Single-size ingredients

**File:** `frontend/src/views/MenuView.vue`

**Thay đổi:** Thêm logic xử lý batch ingredients

```javascript
const preparedIngredients = item.ingredients ? item.ingredients.map(ing => {
  // Check if this is a batch ingredient
  const isBatch = ing.ingredient_type === 'batch' || ing.type === 'batch'
  
  if (isBatch) {
    // Handle batch ingredient
    const batchData = availableBatchDefinitions.value.find(b => b.id === ing.batch_id || b.name === ing.name)
    if (batchData) {
      const batchCostPerUnit = calculateBatchCostPerUnit(batchData)
      const compatibleUnits = getCompatibleUnits(batchData.unit) // ✅ Lấy compatible units
      const conversionRate = getConversionRate(batchData.unit, ing.unit)
      
      // Calculate estimated cost
      const breakdown = calculateCostBreakdown(
        ing.quantity,
        ing.unit,
        batchCostPerUnit,
        batchData.unit,
        0 // Batches don't have wastage
      )
      
      return {
        id: batchData.id,
        batch_definition_id: batchData.id,
        type: 'batch',
        name: ing.name,
        quantity: ing.quantity,
        unit: ing.unit,
        stockUnit: batchData.unit,
        compatibleUnits: compatibleUnits, // ✅ Có đầy đủ compatible units
        costPerUnit: batchCostPerUnit,
        wastage: 0,
        conversionRate: conversionRate,
        estimatedCost: breakdown.totalCost
      }
    }
  } else {
    // Handle raw ingredient (existing logic)
    // ...
  }
  return ing
}) : []
```

### 2. Sửa function `editItem()` - Variant ingredients

**Thay đổi tương tự:** Thêm logic xử lý batch ingredients cho variants

```javascript
const enrichedIngredients = variant.ingredients ? variant.ingredients.map(ing => {
  const isBatch = ing.ingredient_type === 'batch' || ing.type === 'batch'
  
  if (isBatch) {
    // Handle batch ingredient (same logic as single-size)
    const batchData = availableBatchDefinitions.value.find(b => b.id === ing.batch_id || b.name === ing.name)
    if (batchData) {
      // ... (same as above)
      return {
        // ... with compatibleUnits
      }
    }
  } else {
    // Handle raw ingredient
    // ...
  }
  return ing
}) : []
```

## Kết quả

Bây giờ khi edit món đã lưu:
- ✅ Batch ingredients có đầy đủ `compatibleUnits`
- ✅ Dropdown đơn vị hiển thị đúng (ví dụ: "L" và "ml")
- ✅ Đơn vị đã chọn được giữ nguyên
- ✅ Chi phí được tính đúng
- ✅ Conversion rate được giữ nguyên

## Testing

### Test case 1: Single-size item với batch ingredient
1. Tạo batch "Sữa tươi" với đơn vị "L"
2. Tạo món mới, thêm batch "Sữa tươi"
3. Chọn đơn vị "ml", nhập 200ml
4. Lưu món
5. Click "Sửa" để xem lại
6. ✅ Kiểm tra: Dropdown đơn vị có "L" và "ml", đang chọn "ml", số lượng là 200

### Test case 2: Multi-size item với batch ingredient
1. Tạo batch "Bột cà phê" với đơn vị "kg"
2. Tạo món mới với 2 size: M và L
3. Size M: thêm batch "Bột cà phê", chọn "g", nhập 18g
4. Size L: thêm batch "Bột cà phê", chọn "g", nhập 25g
5. Lưu món
6. Click "Sửa" để xem lại
7. ✅ Kiểm tra: Cả 2 size đều có dropdown đơn vị ["kg", "g"], đang chọn "g"

### Test case 3: Mix raw và batch ingredients
1. Tạo món với:
   - Raw ingredient: "Đường" (kg) → chọn "g", 10g
   - Batch ingredient: "Sữa tươi" (L) → chọn "ml", 200ml
2. Lưu món
3. Click "Sửa"
4. ✅ Kiểm tra: Cả 2 ingredients đều có dropdown đúng

## Các thay đổi liên quan

### Thay đổi trước đó (đã hoàn thành)
1. ✅ Cho phép chuyển đổi đơn vị khi thêm batch mới
   - File: `MenuView.vue`, function `selectBatch()`
   - Thay đổi: `compatibleUnits: [batch.unit]` → `compatibleUnits: getCompatibleUnits(batch.unit)`

### Thay đổi hiện tại (vừa hoàn thành)
2. ✅ Giữ nguyên compatible units khi edit món đã lưu
   - File: `MenuView.vue`, function `editItem()`
   - Thay đổi: Thêm logic xử lý batch ingredients

## Backend compatibility

Backend không cần thay đổi vì:
- Backend lưu `ingredient_type: "batch"` và `batch_id`
- Backend không validate conversion rate
- Frontend tự tính conversion và cost

## Notes

- Function `calculateBatchCostPerUnit()` đã tồn tại và hoạt động đúng
- Function `getCompatibleUnits()` từ `useUnitConversion` composable
- Batch ingredients được phân biệt bằng field `ingredient_type === 'batch'` hoặc `type === 'batch'`
- Tìm batch data trong `availableBatchDefinitions` (không phải `availableBatchRecords`)
