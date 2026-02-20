# Fix: Lưu ingredient_type và batch_id khi save menu item

## Vấn đề

Khi tạo món mới với batch ingredient và chọn đơn vị chuyển đổi (ví dụ: batch "L" chọn "ml"), sau khi lưu và click "Sửa" lại:
- ✅ Raw ingredients: hiển thị đúng đơn vị
- ❌ Batch ingredients: không hiển thị đơn vị (dropdown trống)

## Nguyên nhân

Function `saveItem()` khi chuẩn bị data gửi lên API chỉ gửi 3 fields:
```javascript
ingredients: v.ingredients.map(ing => ({
  name: ing.name,
  quantity: ing.quantity,
  unit: ing.unit
  // ❌ Thiếu ingredient_type và batch_id
}))
```

Khi backend lưu vào database, không có `ingredient_type: "batch"` và `batch_id`, nên khi load lại:
- Backend trả về ingredient không có `ingredient_type` → mặc định là "raw"
- Frontend function `editItem()` tìm trong `availableIngredients` (raw) → không tìm thấy
- Không có `compatibleUnits` → dropdown trống

## Giải pháp

Sửa function `saveItem()` để gửi đầy đủ fields cho batch ingredients.

### 1. Sửa cho Multi-size items (variants)

**File:** `frontend/src/views/MenuView.vue`

**Trước:**
```javascript
ingredients: v.ingredients.map(ing => ({
  name: ing.name,
  quantity: ing.quantity,
  unit: ing.unit
}))
```

**Sau:**
```javascript
ingredients: v.ingredients.map(ing => ({
  name: ing.name,
  quantity: ing.quantity,
  unit: ing.unit,
  // Include batch fields if this is a batch ingredient
  ...(ing.type === 'batch' && {
    ingredient_type: 'batch',
    batch_id: ing.batch_definition_id || ing.id
  })
}))
```

### 2. Sửa cho Single-size items

**Trước:**
```javascript
itemData.ingredients = form.value.ingredients.map(ing => ({
  name: ing.name,
  quantity: ing.quantity,
  unit: ing.unit
}))
```

**Sau:**
```javascript
itemData.ingredients = form.value.ingredients.map(ing => ({
  name: ing.name,
  quantity: ing.quantity,
  unit: ing.unit,
  // Include batch fields if this is a batch ingredient
  ...(ing.type === 'batch' && {
    ingredient_type: 'batch',
    batch_id: ing.batch_definition_id || ing.id
  })
}))
```

## Cách hoạt động

### Spread operator với conditional
```javascript
...(ing.type === 'batch' && {
  ingredient_type: 'batch',
  batch_id: ing.batch_definition_id || ing.id
})
```

- Nếu `ing.type === 'batch'` → spread object `{ ingredient_type: 'batch', batch_id: ... }`
- Nếu không phải batch → spread `false` (không có effect)

### Kết quả

**Raw ingredient:**
```json
{
  "name": "Đường",
  "quantity": 10,
  "unit": "g"
}
```

**Batch ingredient:**
```json
{
  "name": "Sữa tươi",
  "quantity": 200,
  "unit": "ml",
  "ingredient_type": "batch",
  "batch_id": "507f1f77bcf86cd799439011"
}
```

## Flow hoàn chỉnh

### 1. Tạo món mới với batch
```
User chọn batch "Sữa tươi" (1L)
  ↓
selectBatch() thêm vào form với:
  - type: 'batch'
  - batch_definition_id: <id>
  - unit: 'L'
  - compatibleUnits: ['L', 'ml']
  ↓
User chọn đơn vị 'ml', nhập 200ml
  ↓
Click "Lưu"
  ↓
saveItem() gửi lên API:
  {
    name: "Sữa tươi",
    quantity: 200,
    unit: "ml",
    ingredient_type: "batch",  ✅ Có field này
    batch_id: <id>              ✅ Có field này
  }
  ↓
Backend lưu vào database
```

### 2. Edit món đã lưu
```
User click "Sửa"
  ↓
Backend trả về ingredient:
  {
    name: "Sữa tươi",
    quantity: 200,
    unit: "ml",
    ingredient_type: "batch",  ✅ Có field này
    batch_id: <id>              ✅ Có field này
  }
  ↓
editItem() kiểm tra ingredient_type === 'batch'
  ↓
Tìm batch trong availableBatchDefinitions
  ↓
Tính compatibleUnits = ['L', 'ml']
  ↓
Form hiển thị đúng:
  - Dropdown: ['L', 'ml']
  - Selected: 'ml'
  - Quantity: 200
```

## Testing

### Test case 1: Tạo mới với batch
1. Tạo batch "Sữa tươi" với đơn vị "L"
2. Tạo món mới, thêm batch "Sữa tươi"
3. Chọn "ml", nhập 200ml
4. Lưu món
5. Mở browser DevTools → Network tab
6. ✅ Kiểm tra request payload có `ingredient_type: "batch"` và `batch_id`

### Test case 2: Edit món đã lưu
1. Tiếp test case 1
2. Click "Sửa" món vừa tạo
3. ✅ Kiểm tra dropdown đơn vị có ["L", "ml"]
4. ✅ Kiểm tra đang chọn "ml"
5. ✅ Kiểm tra số lượng là 200

### Test case 3: Mix raw và batch
1. Tạo món với:
   - Raw: "Đường" 10g
   - Batch: "Sữa tươi" 200ml
2. Lưu món
3. Click "Sửa"
4. ✅ Cả 2 ingredients đều hiển thị đúng dropdown

### Test case 4: Multi-size với batch
1. Tạo món 2 size:
   - Size M: Batch "Bột cà phê" 18g
   - Size L: Batch "Bột cà phê" 25g
2. Lưu món
3. Click "Sửa"
4. ✅ Cả 2 size đều hiển thị đúng dropdown

## Các fix liên quan

### Fix 1: Cho phép chuyển đổi đơn vị khi thêm batch
- File: `MenuView.vue`, function `selectBatch()`
- Thay đổi: `compatibleUnits: [batch.unit]` → `compatibleUnits: getCompatibleUnits(batch.unit)`

### Fix 2: Load đúng compatibleUnits khi edit
- File: `MenuView.vue`, function `editItem()`
- Thay đổi: Thêm logic xử lý batch ingredients

### Fix 3: Lưu đúng ingredient_type và batch_id (FIX HIỆN TẠI)
- File: `MenuView.vue`, function `saveItem()`
- Thay đổi: Gửi `ingredient_type` và `batch_id` lên API

## Backend compatibility

Backend đã support đầy đủ:
```go
type Ingredient struct {
    Name           string               `bson:"name" json:"name"`
    Quantity       float64              `bson:"quantity" json:"quantity"`
    Unit           ingredient.UnitType  `bson:"unit" json:"unit"`
    IngredientType string               `bson:"ingredient_type" json:"ingredient_type"`
    BatchID        *primitive.ObjectID  `bson:"batch_id,omitempty" json:"batch_id,omitempty"`
}
```

## Kết luận

Với 3 fixes:
1. ✅ Cho phép chuyển đổi đơn vị khi thêm batch
2. ✅ Load đúng compatibleUnits khi edit
3. ✅ Lưu đúng ingredient_type và batch_id

Bây giờ batch ingredients hoạt động hoàn toàn giống raw ingredients về mặt unit conversion!
