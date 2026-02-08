# Fix: Adjust Giảm Không Được Thay Đổi Đơn Giá

## Vấn Đề

Khi adjust giảm số lượng nguyên liệu (ví dụ: 10kg → 8kg):
- UI tự động tính lại đơn giá ở field "Đơn giá được tính" (do computed property)
- Khi click "Cập nhật", đơn giá này bị gửi xuống backend
- Backend lưu đơn giá mới → **SAI LOGIC**

### Ví Dụ Lỗi
```
Tồn kho hiện tại: 10kg @ 25,000đ/kg
User adjust giảm xuống: 8kg
UI tự tính: 25,000đ/kg (từ computed property)
Backend nhận: cost_per_unit = 25,000
→ Đơn giá bị cập nhật (không nên!)
```

## Nguyên Nhân

Trong `adjustStock()` function, code đang gửi `cost_per_unit` xuống backend bất kể là tăng hay giảm:

```javascript
// CODE CŨ - SAI
else if (adjustData.value.type === ADJUSTMENT_TYPES.ADJUST) {
  const data = {
    new_quantity: adjustData.value.quantity,
    cost_per_unit: adjustData.value.cost_per_unit || 0,  // ← Luôn gửi!
    reason: adjustData.value.reason
  }
  await ingredientStore.stockAdjust(currentIngredient.value.id, data)
}
```

## Giải Pháp

Chỉ gửi `cost_per_unit` khi:
1. **Adjust TĂNG** số lượng (new_quantity > current_quantity) VÀ
2. **User đã nhập giá mới** (cost_per_unit > 0)

```javascript
// CODE MỚI - ĐÚNG
else if (adjustData.value.type === ADJUSTMENT_TYPES.ADJUST) {
  const isIncrease = adjustData.value.quantity > currentIngredient.value.quantity
  
  const data = {
    new_quantity: adjustData.value.quantity,
    // Only send cost_per_unit if:
    // 1. It's an increase AND
    // 2. User has entered a new price (not 0 or empty)
    cost_per_unit: (isIncrease && adjustData.value.cost_per_unit > 0) ? adjustData.value.cost_per_unit : 0,
    reason: adjustData.value.reason
  }
  await ingredientStore.stockAdjust(currentIngredient.value.id, data)
}
```

## Logic Đúng

### Adjust Giảm (10kg → 8kg)
```
Frontend gửi:
{
  new_quantity: 8,
  cost_per_unit: 0,        // ← Không gửi giá
  reason: "Hỏng hóc"
}

Backend:
- Cập nhật quantity = 8
- KHÔNG thay đổi cost_per_unit
- Giữ nguyên giá 25,000đ/kg
```

### Adjust Tăng Không Nhập Giá (8kg → 12kg)
```
Frontend gửi:
{
  new_quantity: 12,
  cost_per_unit: 0,        // ← User không nhập giá mới
  reason: "Tìm thấy thêm"
}

Backend:
- Cập nhật quantity = 12
- KHÔNG thay đổi cost_per_unit
- Giữ nguyên giá 25,000đ/kg
```

### Adjust Tăng Có Nhập Giá (8kg → 12kg @ 30,000đ/kg)
```
Frontend gửi:
{
  new_quantity: 12,
  cost_per_unit: 30000,    // ← User nhập giá mới
  reason: "Mua thêm"
}

Backend:
- Cập nhật quantity = 12
- Tính weighted average với giá mới 30,000đ/kg
- Cập nhật cost_per_unit
```

## Files Thay Đổi

### `frontend/src/views/IngredientManagementView.vue`
- **Function**: `adjustStock()`
- **Line**: ~1489-1493
- **Change**: Thêm logic kiểm tra `isIncrease` và chỉ gửi `cost_per_unit` khi cần

## Testing

### Test Case 1: Adjust Giảm
1. Nguyên liệu: Đường (10kg @ 25,000đ/kg)
2. Chọn "Điều chỉnh" → Nhập 8kg
3. Nhập lý do: "Hỏng hóc"
4. Click "Xác nhận"
5. **Kết quả mong đợi**:
   - Quantity = 8kg
   - Cost per unit = 25,000đ/kg (không đổi)

### Test Case 2: Adjust Tăng Không Nhập Giá
1. Nguyên liệu: Đường (8kg @ 25,000đ/kg)
2. Chọn "Điều chỉnh" → Nhập 12kg
3. Không nhập giá (để trống hoặc 0)
4. Nhập lý do: "Tìm thấy thêm"
5. Click "Xác nhận"
6. **Kết quả mong đợi**:
   - Quantity = 12kg
   - Cost per unit = 25,000đ/kg (không đổi)

### Test Case 3: Adjust Tăng Có Nhập Giá
1. Nguyên liệu: Đường (8kg @ 25,000đ/kg)
2. Chọn "Điều chỉnh" → Nhập 12kg
3. Nhập giá mới: 30,000đ/kg
4. Nhập lý do: "Mua thêm với giá mới"
5. Click "Xác nhận"
6. **Kết quả mong đợi**:
   - Quantity = 12kg
   - Cost per unit = weighted average (tính toán mới)

## Backend Logic (Tham Khảo)

Backend đã có logic đúng trong `StockAdjust()`:

```go
// backend/domain/ingredient/ingredient.go
func (i *Ingredient) StockAdjust(newQuantity float64, newPrice float64, reason string) error {
    if newQuantity < 0 {
        return errors.New("new quantity cannot be negative")
    }
    
    oldQuantity := i.Quantity
    
    // If increasing quantity AND new price provided
    if newQuantity > oldQuantity && newPrice > 0 && newPrice != i.CostPerUnit {
        // Calculate weighted average
        quantityIncrease := newQuantity - oldQuantity
        totalOldCost := oldQuantity * i.CostPerUnit
        totalNewCost := quantityIncrease * newPrice
        i.CostPerUnit = (totalOldCost + totalNewCost) / newQuantity
    }
    // Otherwise, keep current price
    
    i.Quantity = newQuantity
    return nil
}
```

## Kết Luận

✅ **Fix hoàn thành**: Frontend giờ chỉ gửi `cost_per_unit` khi adjust tăng và user nhập giá mới

✅ **Logic đúng**: Adjust giảm không thay đổi đơn giá

✅ **Tương thích backend**: Frontend và backend đồng bộ logic

---

**Ngày fix**: 2026-02-07  
**File**: `frontend/src/views/IngredientManagementView.vue`  
**Function**: `adjustStock()`

