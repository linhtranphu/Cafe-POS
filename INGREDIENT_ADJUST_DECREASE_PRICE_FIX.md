# Fix: Đơn Giá Bị Tính Lại Khi Adjust Giảm

## Vấn Đề

Khi sử dụng function "Điều chỉnh" (adjust) để **GIẢM** tồn kho (ví dụ: từ 10kg → 8kg), có 2 vấn đề:

1. **UI hiển thị sai**: Vẫn hiển thị "Đơn giá được tính" mặc dù đang giảm
2. **Backend cập nhật sai**: Đơn giá bị thay đổi mặc dù không nên

### Ví Dụ Lỗi

```
Tồn kho hiện tại: 10 kg @ 50,000đ/kg
Điều chỉnh thành: 8 kg (giảm 2kg)

❌ Kết quả SAI: 
- UI hiển thị "Đơn giá được tính"
- Backend cập nhật đơn giá mới

✅ Kết quả ĐÚNG: 
- UI KHÔNG hiển thị phần nhập giá
- Đơn giá giữ nguyên 50,000đ/kg
```

## Nguyên Nhân

### 1. UI Hiển thị Sai
```vue
<!-- Hiển thị khi adjustData.quantity > 0 && adjustTotalPrice > 0 -->
<div v-if="adjustData.quantity > 0 && adjustTotalPrice > 0">
  <div>Đơn giá được tính:</div>
</div>
```

**Vấn đề**: 
- `adjustData.quantity` là **target quantity** (8kg), không phải diff (-2kg)
- Điều kiện `> 0` luôn đúng cả khi giảm
- Không kiểm tra xem có đang TĂNG hay GIẢM

### 2. Backend Nhận Giá Sai
```javascript
// Frontend gửi
const data = {
  new_quantity: 8,
  cost_per_unit: 45000,  // ❌ Giá trị cũ còn lại từ lần trước!
  reason: "Kiểm kê"
}
```

**Vấn đề**:
- Khi type = 'adjust', không có UI để nhập giá
- Nhưng `adjustData.cost_per_unit` vẫn còn giá trị cũ
- Backend nhận được giá → tính lại (SAI!)

## Giải Pháp

### 1. Chỉ Hiển Thị UI Nhập Giá Khi TĂNG

**Thêm điều kiện kiểm tra tăng/giảm:**
```vue
<!-- Chỉ hiển thị khi type = 'adjust' VÀ số lượng TĂNG -->
<div v-if="adjustData.type === 'adjust' && adjustData.quantity > (currentIngredient?.quantity || 0)">
  <!-- UI nhập giá -->
</div>
```

### 2. Reset Giá Khi Chuyển Sang Giảm

**Thêm function `onAdjustQuantityChange`:**
```javascript
const onAdjustQuantityChange = () => {
  // Recalculate if in total price mode
  if (adjustPriceMode.value === 'total') {
    calculateAdjustUnitPrice()
  }
  
  // Reset price when switching from increase to decrease
  if (adjustData.value.type === 'adjust') {
    const isIncreasing = adjustData.value.quantity > (currentIngredient.value?.quantity || 0)
    if (!isIncreasing) {
      // If decreasing, clear price inputs
      adjustData.value.cost_per_unit = 0
      adjustTotalPrice.value = 0
    }
  }
}
```

### 3. Tính Giá Dựa Trên Phần TĂNG

**Cập nhật `calculateAdjustUnitPrice`:**
```javascript
const calculateAdjustUnitPrice = () => {
  if (adjustData.value.type === 'adjust') {
    // For adjust type, calculate based on the INCREASE amount
    const quantityDiff = adjustData.value.quantity - (currentIngredient.value?.quantity || 0)
    if (quantityDiff > 0 && adjustTotalPrice.value > 0) {
      adjustData.value.cost_per_unit = adjustTotalPrice.value / quantityDiff
    } else {
      adjustData.value.cost_per_unit = 0
    }
  } else {
    // For add/remove type, calculate normally
    if (adjustData.value.quantity > 0 && adjustTotalPrice.value > 0) {
      adjustData.value.cost_per_unit = adjustTotalPrice.value / adjustData.value.quantity
    } else {
      adjustData.value.cost_per_unit = 0
    }
  }
}
```

## UI Mới

### Khi Điều Chỉnh TĂNG (10kg → 12kg)

```
┌─────────────────────────────────────────┐
│ ⚠️ Số lượng tăng - Có nhập thêm hàng?  │
│                                         │
│ Nếu tăng do nhập thêm hàng với giá mới,│
│ hãy nhập giá bên dưới. Nếu chỉ điều    │
│ chỉnh số liệu, để trống.                │
│                                         │
│ 💰 Giá nhập thêm (tùy chọn)            │
│ [Nhập tổng giá]                         │
│                                         │
│ Tổng giá mua thêm (nếu có)             │
│ ┌─────────────────────────────────┐    │
│ │ Để trống nếu không mua thêm     │    │
│ └─────────────────────────────────┘    │
│                                         │
│ ✅ Đơn giá được tính cho phần tăng:    │
│    60,000đ/kg                           │
│    = 120,000đ ÷ 2 kg (phần tăng)       │
└─────────────────────────────────────────┘
```

### Khi Điều Chỉnh GIẢM (10kg → 8kg)

```
┌─────────────────────────────────────────┐
│ Số lượng *                              │
│ ┌─────────────────────────────────┐    │
│ │ 8                               │    │
│ └─────────────────────────────────┘    │
│                                         │
│ Lý do *                                 │
│ ┌─────────────────────────────────┐    │
│ │ Kiểm kê - thiếu hụt             │    │
│ └─────────────────────────────────┘    │
│                                         │
│ ℹ️ Tồn kho sau điều chỉnh:             │
│    8 kg                                 │
└─────────────────────────────────────────┘

❌ KHÔNG có phần nhập giá
```

## Test Cases

### Test 1: Adjust Giảm - Không Nhập Giá
```
Hiện tại: 10 kg @ 50,000đ/kg
Điều chỉnh: 8 kg

UI:
- ❌ KHÔNG hiển thị phần nhập giá
- ✅ Chỉ có số lượng và lý do

Request:
POST /manager/ingredients/:id/stock-adjust
{
  "new_quantity": 8,
  "cost_per_unit": 0,  // ✅ Tự động = 0
  "reason": "Kiểm kê - thiếu hụt"
}

Backend Logic:
- quantityDiff = 8 - 10 = -2 (giảm)
- quantityDiff < 0 → KHÔNG tính lại giá
- item.CostPerUnit giữ nguyên

✅ Kết quả: 8 kg @ 50,000đ/kg (giá KHÔNG đổi)
```

### Test 2: Adjust Tăng - Không Nhập Giá
```
Hiện tại: 10 kg @ 50,000đ/kg
Điều chỉnh: 12 kg (không nhập giá)

UI:
- ✅ Hiển thị phần nhập giá (tùy chọn)
- ✅ User để trống

Request:
POST /manager/ingredients/:id/stock-adjust
{
  "new_quantity": 12,
  "cost_per_unit": 0,  // ✅ User không nhập
  "reason": "Kiểm kê - tìm thấy thêm"
}

Backend Logic:
- quantityDiff = 12 - 10 = +2 (tăng)
- req.CostPerUnit = 0 → KHÔNG tính lại giá
- item.CostPerUnit giữ nguyên

✅ Kết quả: 12 kg @ 50,000đ/kg (giá KHÔNG đổi)
```

### Test 3: Adjust Tăng - Có Nhập Giá Mới
```
Hiện tại: 10 kg @ 50,000đ/kg
Điều chỉnh: 12 kg + nhập giá 60,000đ/kg

UI:
- ✅ Hiển thị phần nhập giá
- ✅ User nhập tổng giá: 120,000đ cho 2kg tăng thêm

Calculation:
- Phần tăng: 12 - 10 = 2 kg
- Đơn giá: 120,000 ÷ 2 = 60,000đ/kg

Request:
POST /manager/ingredients/:id/stock-adjust
{
  "new_quantity": 12,
  "cost_per_unit": 60000,  // ✅ Giá mới
  "reason": "Kiểm kê + nhập thêm"
}

Backend Logic:
- quantityDiff = 12 - 10 = +2 (tăng)
- req.CostPerUnit = 60,000 (khác 50,000)
- Tính weighted average:
  - oldValue = 10 × 50,000 = 500,000
  - newValue = 2 × 60,000 = 120,000
  - new_price = 620,000 ÷ 12 = 51,667đ/kg

✅ Kết quả: 12 kg @ 51,667đ/kg (weighted average)
```

### Test 4: Chuyển Từ Tăng Sang Giảm
```
User Action:
1. Nhập số lượng: 12 (tăng)
2. Nhập giá: 60,000đ
3. Đổi ý, nhập số lượng: 8 (giảm)

Behavior:
- onAdjustQuantityChange() được gọi
- Phát hiện: 8 < 10 (giảm)
- Tự động reset:
  - adjustData.cost_per_unit = 0
  - adjustTotalPrice = 0
- UI ẩn phần nhập giá

✅ Kết quả: Giá bị xóa tự động, không gửi lên backend
```

## Backend Logic (Không Đổi)

Backend logic đã đúng từ trước:

```go
func (s *IngredientService) StockAdjust(...) {
    quantityDiff := req.NewQuantity - beforeQty
    item.Quantity = req.NewQuantity
    
    // Only recalculate if:
    // 1. Quantity increased (positive diff)
    // 2. New price provided and different
    if quantityDiff > 0 && req.CostPerUnit > 0 && req.CostPerUnit != item.CostPerUnit {
        // Weighted average
        oldValue := beforeQty * item.CostPerUnit
        newValue := quantityDiff * req.CostPerUnit
        item.CostPerUnit = (oldValue + newValue) / afterQty
    }
    // Otherwise: Keep current price
}
```

**Điều kiện tính lại giá:**
- ✅ `quantityDiff > 0` - Phải TĂNG
- ✅ `req.CostPerUnit > 0` - Phải có giá mới
- ✅ `req.CostPerUnit != item.CostPerUnit` - Giá phải KHÁC

## Quy Tắc Cuối Cùng

| Scenario | Tăng/Giảm | Nhập Giá? | UI Hiển Thị? | Backend Tính Lại? |
|----------|-----------|-----------|--------------|-------------------|
| Adjust giảm | Giảm | N/A | ❌ Ẩn | ❌ NO |
| Adjust tăng, không giá | Tăng | Không | ✅ Hiện (tùy chọn) | ❌ NO |
| Adjust tăng, giá = 0 | Tăng | 0 | ✅ Hiện (tùy chọn) | ❌ NO |
| Adjust tăng, giá mới | Tăng | Có | ✅ Hiện | ✅ YES |
| Adjust tăng, giá giống | Tăng | Giống | ✅ Hiện | ❌ NO |

## Files Đã Thay Đổi

### Frontend
- ✅ `frontend/src/views/IngredientManagementView.vue`
  - Thêm UI nhập giá cho adjust type (chỉ khi tăng)
  - Thêm `onAdjustQuantityChange()` để reset giá
  - Cập nhật `calculateAdjustUnitPrice()` để tính dựa trên phần tăng
  - Ẩn UI nhập giá khi giảm

### Backend
- ✅ Không cần thay đổi (logic đã đúng)

## Lợi Ích

### 1. UI Rõ Ràng
- Chỉ hiển thị nhập giá khi cần (tăng)
- Không gây nhầm lẫn khi giảm
- Có cảnh báo rõ ràng về mục đích

### 2. Logic Đúng
- Giảm → KHÔNG BAO GIỜ đổi giá
- Tăng không giá → Giữ nguyên giá
- Tăng có giá → Tính weighted average

### 3. UX Tốt
- Tự động reset giá khi chuyển từ tăng sang giảm
- Tính giá dựa trên phần tăng thêm (chính xác hơn)
- Có hướng dẫn rõ ràng

### 4. Tránh Lỗi
- Không thể vô tình gửi giá khi giảm
- Validation tự động
- Dữ liệu sạch

## Kết Luận

Fix này giải quyết hoàn toàn vấn đề đơn giá bị tính lại khi adjust giảm:

✅ **UI**: Chỉ hiển thị nhập giá khi tăng
✅ **Logic**: Tự động reset giá khi giảm
✅ **Calculation**: Tính dựa trên phần tăng thêm
✅ **Backend**: Không cần thay đổi (đã đúng)

Giờ đây:
- Adjust giảm → UI không có phần giá, backend không đổi giá
- Adjust tăng không giá → UI có phần giá (tùy chọn), backend giữ nguyên
- Adjust tăng có giá → UI tính đúng, backend tính weighted average
