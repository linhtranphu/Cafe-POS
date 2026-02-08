# Fix: Đơn Giá Không Nên Tính Lại Khi Điều Chỉnh Tồn Kho

## Vấn Đề

Khi sử dụng function "Điều chỉnh" (adjust) để set tồn kho về một số lượng cụ thể (ví dụ: từ 10kg → 12kg), đơn giá đang bị tính lại theo weighted average, mặc dù không nên.

### Ví Dụ Lỗi

```
Tồn kho hiện tại: 10 kg @ 50,000đ/kg
Điều chỉnh thành: 12 kg (không nhập giá mới)

❌ Kết quả SAI: Đơn giá bị tính lại
✅ Kết quả ĐÚNG: Đơn giá giữ nguyên 50,000đ/kg
```

## Nguyên Nhân

Frontend đang sử dụng **legacy method** `adjustStock` với logic cũ:
- Khi type = 'adjust', frontend tính: `quantity = newQty - currentQty`
- Ví dụ: 12 - 10 = +2
- Backend nhận `quantity = +2` và `cost_per_unit = 0`
- Backend thấy quantity > 0 → nghĩ là "stock IN" → tính lại giá (SAI!)

## Giải Pháp

Chuyển sang sử dụng **API mới** với 3 methods riêng biệt:

### 1. Stock IN (`/stock-in`)
- Dùng khi: Nhập hàng mới
- Request: `{ quantity, cost_per_unit, reason }`
- Logic: Tính weighted average nếu giá khác

### 2. Stock OUT (`/stock-out`)
- Dùng khi: Xuất hàng, sử dụng, hỏng
- Request: `{ quantity, reason }`
- Logic: KHÔNG BAO GIỜ tính lại giá

### 3. Stock ADJUST (`/stock-adjust`)
- Dùng khi: Kiểm kê, điều chỉnh số lượng
- Request: `{ new_quantity, cost_per_unit, reason }`
- Logic: Chỉ tính lại giá nếu tăng + có giá mới

## Thay Đổi Code

### Frontend: Quick Stock IN

**Trước:**
```javascript
const adjustment = {
  type: 'add',
  quantity: quickInData.value.quantity,
  reason: 'Nhập kho',
  cost_per_unit: quickInData.value.cost_per_unit || 0
}
await ingredientStore.adjustStock(currentIngredient.value.id, adjustment)
```

**Sau:**
```javascript
const data = {
  quantity: quickInData.value.quantity,
  cost_per_unit: quickInData.value.cost_per_unit || 0,
  reason: 'Nhập kho'
}
await ingredientStore.stockIn(currentIngredient.value.id, data)
```

### Frontend: Quick Stock OUT

**Trước:**
```javascript
const adjustment = {
  type: 'remove',
  quantity: -Math.abs(quickOutData.value.quantity),
  reason: finalReason,
  cost_per_unit: 0
}
await ingredientStore.adjustStock(currentIngredient.value.id, adjustment)
```

**Sau:**
```javascript
const data = {
  quantity: quickOutData.value.quantity,
  reason: finalReason
}
await ingredientStore.stockOut(currentIngredient.value.id, data)
```

### Frontend: Stock Adjust Modal

**Trước:**
```javascript
let finalQuantity = adjustData.value.quantity

if (adjustData.value.type === 'remove') {
  finalQuantity = -Math.abs(adjustData.value.quantity)
} else if (adjustData.value.type === 'add') {
  finalQuantity = Math.abs(adjustData.value.quantity)
} else if (adjustData.value.type === 'adjust') {
  finalQuantity = adjustData.value.quantity - currentIngredient.value.quantity
}

const adjustmentData = {
  ...adjustData.value,
  quantity: finalQuantity
}
await ingredientStore.adjustStock(currentIngredient.value.id, adjustmentData)
```

**Sau:**
```javascript
if (adjustData.value.type === 'add') {
  // Stock IN
  const data = {
    quantity: adjustData.value.quantity,
    cost_per_unit: adjustData.value.cost_per_unit || 0,
    reason: adjustData.value.reason
  }
  await ingredientStore.stockIn(currentIngredient.value.id, data)
} else if (adjustData.value.type === 'remove') {
  // Stock OUT
  const data = {
    quantity: adjustData.value.quantity,
    reason: adjustData.value.reason
  }
  await ingredientStore.stockOut(currentIngredient.value.id, data)
} else if (adjustData.value.type === 'adjust') {
  // Stock ADJUST
  const data = {
    new_quantity: adjustData.value.quantity,
    cost_per_unit: adjustData.value.cost_per_unit || 0,
    reason: adjustData.value.reason
  }
  await ingredientStore.stockAdjust(currentIngredient.value.id, data)
}
```

## Backend Logic

### StockAdjust Method
```go
func (s *IngredientService) StockAdjust(ctx context.Context, id primitive.ObjectID, req *ingredient.StockAdjustRequest) (*ingredient.Ingredient, error) {
    // Calculate difference
    quantityDiff := req.NewQuantity - beforeQty
    item.Quantity = req.NewQuantity
    
    // Only recalculate price if:
    // 1. Quantity increased (positive diff)
    // 2. New price provided and different from current
    if quantityDiff > 0 && req.CostPerUnit > 0 && req.CostPerUnit != item.CostPerUnit && afterQty > 0 {
        // Weighted average for the increase
        oldValue := beforeQty * item.CostPerUnit
        newValue := quantityDiff * req.CostPerUnit
        item.CostPerUnit = (oldValue + newValue) / afterQty
    }
    // Otherwise: Keep current price
}
```

## Test Cases

### Test 1: Điều Chỉnh Tăng Không Có Giá
```
Hiện tại: 10 kg @ 50,000đ/kg
Điều chỉnh: 12 kg (không nhập giá)

Request:
POST /manager/ingredients/:id/stock-adjust
{
  "new_quantity": 12,
  "cost_per_unit": 0,
  "reason": "Kiểm kê"
}

✅ Kết quả: 12 kg @ 50,000đ/kg (giá KHÔNG đổi)
```

### Test 2: Điều Chỉnh Tăng Có Giá Mới
```
Hiện tại: 10 kg @ 50,000đ/kg
Điều chỉnh: 12 kg @ 60,000đ/kg (nhập giá mới)

Request:
POST /manager/ingredients/:id/stock-adjust
{
  "new_quantity": 12,
  "cost_per_unit": 60000,
  "reason": "Kiểm kê + nhập thêm"
}

✅ Kết quả: 12 kg @ 51,667đ/kg (weighted average)
Tính: (10 × 50,000 + 2 × 60,000) / 12 = 51,667
```

### Test 3: Điều Chỉnh Giảm
```
Hiện tại: 10 kg @ 50,000đ/kg
Điều chỉnh: 8 kg

Request:
POST /manager/ingredients/:id/stock-adjust
{
  "new_quantity": 8,
  "cost_per_unit": 0,
  "reason": "Kiểm kê - thiếu hụt"
}

✅ Kết quả: 8 kg @ 50,000đ/kg (giá KHÔNG đổi)
```

### Test 4: Nhập Hàng Không Có Giá
```
Hiện tại: 10 kg @ 50,000đ/kg
Nhập thêm: 5 kg (không nhập giá)

Request:
POST /manager/ingredients/:id/stock-in
{
  "quantity": 5,
  "cost_per_unit": 0,
  "reason": "Nhập kho"
}

✅ Kết quả: 15 kg @ 50,000đ/kg (giá KHÔNG đổi)
```

### Test 5: Nhập Hàng Có Giá Mới
```
Hiện tại: 10 kg @ 50,000đ/kg
Nhập thêm: 5 kg @ 60,000đ/kg

Request:
POST /manager/ingredients/:id/stock-in
{
  "quantity": 5,
  "cost_per_unit": 60000,
  "reason": "Nhập kho"
}

✅ Kết quả: 15 kg @ 53,333đ/kg (weighted average)
Tính: (10 × 50,000 + 5 × 60,000) / 15 = 53,333
```

### Test 6: Xuất Hàng
```
Hiện tại: 10 kg @ 50,000đ/kg
Xuất: 3 kg

Request:
POST /manager/ingredients/:id/stock-out
{
  "quantity": 3,
  "reason": "Sử dụng cho món ăn"
}

✅ Kết quả: 7 kg @ 50,000đ/kg (giá KHÔNG BAO GIỜ đổi)
```

## Quy Tắc Tính Giá

| Operation | Tăng/Giảm | Có Giá Mới? | Giá Khác? | Tính Lại? |
|-----------|-----------|-------------|-----------|-----------|
| Stock IN | Tăng | Có | Có | ✅ YES |
| Stock IN | Tăng | Có | Không | ❌ NO |
| Stock IN | Tăng | Không | N/A | ❌ NO |
| Stock OUT | Giảm | N/A | N/A | ❌ NO |
| Stock ADJUST | Tăng | Có | Có | ✅ YES |
| Stock ADJUST | Tăng | Có | Không | ❌ NO |
| Stock ADJUST | Tăng | Không | N/A | ❌ NO |
| Stock ADJUST | Giảm | N/A | N/A | ❌ NO |

## Lợi Ích

### 1. Logic Rõ Ràng
- Mỗi operation có method riêng
- Không còn nhầm lẫn về quantity âm/dương
- Backend xử lý toàn bộ business logic

### 2. Frontend Đơn Giản
- Không cần tính toán quantity difference
- Không cần lo về dấu âm/dương
- Chỉ cần gọi đúng method

### 3. Đúng Business Rules
- Stock IN: Có thể thay đổi giá
- Stock OUT: Không bao giờ thay đổi giá
- Stock ADJUST: Chỉ thay đổi giá khi cần

### 4. Dễ Maintain
- Code rõ ràng, dễ đọc
- Dễ test từng operation
- Dễ mở rộng sau này

## Files Đã Thay Đổi

### Backend
- ✅ `backend/domain/ingredient/ingredient.go` - Thêm request types mới
- ✅ `backend/application/services/ingredient.go` - Thêm 3 methods mới
- ✅ `backend/interfaces/http/ingredient_handler.go` - Thêm 3 handlers mới
- ✅ `backend/main.go` - Thêm 3 routes mới

### Frontend
- ✅ `frontend/src/constants/ingredient.js` - Thêm constants mới
- ✅ `frontend/src/services/ingredient.js` - Thêm 3 service methods
- ✅ `frontend/src/stores/ingredient.js` - Thêm 3 store actions
- ✅ `frontend/src/views/IngredientManagementView.vue` - Cập nhật sử dụng API mới

## Kết Luận

Fix này giải quyết vấn đề đơn giá bị tính lại không đúng lúc bằng cách:
1. Tách riêng 3 operations: IN, OUT, ADJUST
2. Mỗi operation có logic riêng rõ ràng
3. Frontend đơn giản hơn, chỉ gọi đúng method
4. Backend xử lý toàn bộ business logic

Giờ đây:
- ✅ Nhập hàng với giá mới → Tính weighted average
- ✅ Nhập hàng không có giá → Giữ nguyên giá
- ✅ Xuất hàng → Không bao giờ đổi giá
- ✅ Điều chỉnh không có giá → Giữ nguyên giá
- ✅ Điều chỉnh có giá mới (nếu tăng) → Tính weighted average
