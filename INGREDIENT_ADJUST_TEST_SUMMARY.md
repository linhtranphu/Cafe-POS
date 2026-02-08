# Tóm Tắt: Fix Điều Chỉnh Giảm Không Thay Đổi Đơn Giá

## ✅ Trạng Thái: HOÀN THÀNH

**Ngày**: 2026-02-07  
**Files thay đổi**: `frontend/src/views/IngredientManagementView.vue`

---

## 🎯 Vấn Đề Đã Fix

### Vấn Đề 1: Adjust Giảm Bị Tính Lại Đơn Giá
**Mô tả**: Khi điều chỉnh giảm số lượng (10kg → 8kg), đơn giá bị tính lại và lưu xuống database.

**Nguyên nhân**: Frontend luôn gửi `cost_per_unit` xuống backend, kể cả khi giảm số lượng.

**Giải pháp**: Chỉ gửi `cost_per_unit` khi:
- Adjust TĂNG số lượng VÀ
- User đã nhập giá mới (> 0)

### Vấn Đề 2: Form "Sửa" Cho Phép Thay Đổi Tồn Kho
**Mô tả**: User click nút "Sửa" (✏️) có thể thay đổi số lượng và giá, gây nhầm lẫn với chức năng "Điều chỉnh".

**Giải pháp**: 
- Disable field số lượng khi edit
- Ẩn section nhập giá khi edit
- Hiển thị thông tin read-only với cảnh báo

### Vấn Đề 3: UI Phức Tạp
**Mô tả**: Có quá nhiều nút (6 nút: Nhập, Xuất, Điều chỉnh, Lịch sử, Sửa, Xóa).

**Giải pháp**: Bỏ nút "Nhập nhanh" và "Xuất nhanh", chỉ giữ 4 nút:
- 📦 Điều chỉnh
- 📊 Lịch sử
- ✏️ Sửa
- 🗑️ Xóa

---

## 🔧 Code Changes

### 1. Frontend Logic - `adjustStock()` Function (Line 1489-1507)

```javascript
else if (adjustData.value.type === ADJUSTMENT_TYPES.ADJUST) {
  const isIncrease = adjustData.value.quantity > currentIngredient.value.quantity
  
  console.log('=== ADJUST DEBUG ===')
  console.log('Current quantity:', currentIngredient.value.quantity)
  console.log('New quantity:', adjustData.value.quantity)
  console.log('Is increase?', isIncrease)
  console.log('cost_per_unit before:', adjustData.value.cost_per_unit)
  
  const data = {
    new_quantity: adjustData.value.quantity,
    // Only send cost_per_unit if:
    // 1. It's an increase AND
    // 2. User has entered a new price (not 0 or empty)
    cost_per_unit: (isIncrease && adjustData.value.cost_per_unit > 0) ? adjustData.value.cost_per_unit : 0,
    reason: adjustData.value.reason
  }
  
  console.log('Data to send:', data)
  console.log('===================')
  
  await ingredientStore.stockAdjust(currentIngredient.value.id, data)
}
```

**Key Points**:
- ✅ Kiểm tra `isIncrease` để xác định tăng hay giảm
- ✅ Chỉ gửi `cost_per_unit` khi `isIncrease && cost_per_unit > 0`
- ✅ Có console.log để debug (sẽ xóa sau khi test xong)

### 2. Edit Form - Disable Quantity Field (Line ~310)

```vue
<input v-model.number="formData.quantity" type="number" 
  :disabled="isEditing"
  :class="isEditing ? 'bg-gray-100 cursor-not-allowed' : ''"
  class="..." />
```

### 3. Edit Form - Hide Price Section (Line ~253)

```vue
<div v-if="!isEditing" class="bg-gray-50 rounded-xl p-4 space-y-4">
  <h3>💰 Thông tin giá</h3>
  <!-- Price inputs -->
</div>
```

### 4. Edit Form - Show Read-only Info (Line ~301)

```vue
<div v-else class="bg-blue-50 rounded-xl p-4 border-2 border-blue-200">
  <h3>📊 Thông tin hiện tại (chỉ xem)</h3>
  <div class="grid grid-cols-2 gap-3">
    <div>
      <p>Tồn kho</p>
      <p>{{ formData.quantity }} {{ formData.unit }}</p>
    </div>
    <div>
      <p>Đơn giá</p>
      <p>{{ formatCurrency(formData.cost_per_unit) }}</p>
    </div>
  </div>
  <p class="text-orange-600">
    ⚠️ Để thay đổi tồn kho hoặc giá, vui lòng sử dụng chức năng "Điều chỉnh"
  </p>
</div>
```

### 5. UI Simplification - Button Layout (Line ~122-148)

```vue
<!-- Before: 6 buttons in 2 rows -->
<!-- After: 4 buttons in 1 row -->
<div class="grid grid-cols-4 gap-1.5 pt-3 border-t">
  <button @click="openAdjustModal(ingredient)">📦 Điều chỉnh</button>
  <button @click="viewHistory(ingredient)">📊 Lịch sử</button>
  <button @click="openEditModal(ingredient)">✏️ Sửa</button>
  <button @click="deleteIngredient(ingredient)">🗑️ Xóa</button>
</div>
```

---

## 🧪 Test Cases

### ✅ Test Case 1: Adjust Giảm - Không Thay Đổi Giá

**Bước thực hiện**:
1. Mở http://localhost:5173/#/ingredients
2. Chọn nguyên liệu: Đường (10kg @ 25,000đ/kg)
3. Click nút "📦 Điều chỉnh"
4. Chọn loại: "Điều chỉnh"
5. Nhập số lượng mới: 8kg
6. Không nhập giá (để 0 hoặc trống)
7. Nhập lý do: "Hỏng hóc"
8. Click "Xác nhận"

**Kết quả mong đợi**:
- ✅ Quantity = 8kg
- ✅ Cost per unit = 25,000đ/kg (KHÔNG ĐỔI)
- ✅ Console log hiển thị:
  ```
  Is increase? false
  cost_per_unit before: 0
  Data to send: { new_quantity: 8, cost_per_unit: 0, reason: "Hỏng hóc" }
  ```

### ✅ Test Case 2: Adjust Tăng Không Nhập Giá

**Bước thực hiện**:
1. Nguyên liệu: Đường (8kg @ 25,000đ/kg)
2. Click "📦 Điều chỉnh"
3. Chọn loại: "Điều chỉnh"
4. Nhập số lượng mới: 12kg
5. Không nhập giá (để 0)
6. Nhập lý do: "Tìm thấy thêm"
7. Click "Xác nhận"

**Kết quả mong đợi**:
- ✅ Quantity = 12kg
- ✅ Cost per unit = 25,000đ/kg (KHÔNG ĐỔI)
- ✅ Console log:
  ```
  Is increase? true
  cost_per_unit before: 0
  Data to send: { new_quantity: 12, cost_per_unit: 0, reason: "Tìm thấy thêm" }
  ```

### ✅ Test Case 3: Adjust Tăng Có Nhập Giá

**Bước thực hiện**:
1. Nguyên liệu: Đường (8kg @ 25,000đ/kg)
2. Click "📦 Điều chỉnh"
3. Chọn loại: "Điều chỉnh"
4. Nhập số lượng mới: 12kg
5. Nhập giá mới: 30,000đ/kg
6. Nhập lý do: "Mua thêm với giá mới"
7. Click "Xác nhận"

**Kết quả mong đợi**:
- ✅ Quantity = 12kg
- ✅ Cost per unit = weighted average (tính toán mới)
- ✅ Console log:
  ```
  Is increase? true
  cost_per_unit before: 30000
  Data to send: { new_quantity: 12, cost_per_unit: 30000, reason: "Mua thêm..." }
  ```

### ✅ Test Case 4: Form "Sửa" Không Thể Thay Đổi Tồn Kho

**Bước thực hiện**:
1. Click nút "✏️ Sửa"
2. Thử thay đổi số lượng

**Kết quả mong đợi**:
- ✅ Field số lượng bị disable (màu xám)
- ✅ Không thể click hoặc nhập
- ✅ Hiển thị warning: "Dùng 'Điều chỉnh' để thay đổi tồn kho"

### ✅ Test Case 5: Form "Sửa" Không Hiển Thị Input Giá

**Bước thực hiện**:
1. Click nút "✏️ Sửa"
2. Tìm section nhập giá

**Kết quả mong đợi**:
- ✅ Không có section "💰 Thông tin giá"
- ✅ Thay vào đó hiển thị "📊 Thông tin hiện tại (chỉ xem)"
- ✅ Hiển thị tồn kho và giá ở chế độ read-only
- ✅ Có warning: "Dùng 'Điều chỉnh' để thay đổi"

### ✅ Test Case 6: UI Chỉ Có 4 Nút

**Bước thực hiện**:
1. Xem danh sách nguyên liệu
2. Đếm số nút action

**Kết quả mong đợi**:
- ✅ Chỉ có 4 nút: 📦 Điều chỉnh, 📊 Lịch sử, ✏️ Sửa, 🗑️ Xóa
- ✅ Không có nút "Nhập nhanh" và "Xuất nhanh"
- ✅ 4 nút nằm trên 1 hàng

---

## 🔍 Backend Verification

### Backend Logic (Đã Đúng - Không Cần Sửa)

File: `backend/application/services/ingredient.go` (Line 268-275)

```go
// Only recalculate price if:
// 1. Quantity increased (positive diff)
// 2. New price provided and different from current
if quantityDiff > 0 && req.CostPerUnit > 0 && req.CostPerUnit != item.CostPerUnit && afterQty > 0 {
    // Weighted average for the increase
    oldValue := beforeQty * item.CostPerUnit
    newValue := quantityDiff * req.CostPerUnit
    item.CostPerUnit = (oldValue + newValue) / afterQty
    costPerUnit = req.CostPerUnit
}
```

**Key Points**:
- ✅ Chỉ tính lại giá khi `quantityDiff > 0` (tăng số lượng)
- ✅ Chỉ tính lại giá khi `req.CostPerUnit > 0` (có giá mới)
- ✅ Khi `req.CostPerUnit = 0`, backend giữ nguyên giá hiện tại

---

## 📋 Checklist Hoàn Thành

### Frontend Changes
- ✅ Fix `adjustStock()` function - chỉ gửi price khi increase
- ✅ Disable quantity field trong edit form
- ✅ Ẩn price input section trong edit form
- ✅ Hiển thị read-only info trong edit form
- ✅ Bỏ nút "Nhập nhanh" và "Xuất nhanh"
- ✅ Giữ 4 nút: Điều chỉnh, Lịch sử, Sửa, Xóa
- ✅ Thêm console.log debugging

### Backend Verification
- ✅ Verify backend logic đúng (không cần sửa)
- ✅ Backend chỉ tính lại giá khi `quantityDiff > 0 && req.CostPerUnit > 0`

### Documentation
- ✅ `INGREDIENT_ADJUST_DECREASE_PRICE_FIX_FINAL.md`
- ✅ `INGREDIENT_EDIT_FORM_FIX_VI.md`
- ✅ `INGREDIENT_ADJUST_TEST_SUMMARY.md` (file này)

---

## 🚀 Next Steps

### 1. Testing (Cần User Test)
- [ ] Test Case 1: Adjust giảm không thay đổi giá
- [ ] Test Case 2: Adjust tăng không nhập giá
- [ ] Test Case 3: Adjust tăng có nhập giá
- [ ] Test Case 4: Form "Sửa" không thể thay đổi tồn kho
- [ ] Test Case 5: Form "Sửa" không hiển thị input giá
- [ ] Test Case 6: UI chỉ có 4 nút

### 2. Cleanup (Sau Khi Test Xong)
- [ ] Xóa console.log debugging statements (line 1493-1503)
- [ ] Verify không có lỗi trong production

### 3. Deployment
- [ ] Build frontend: `cd frontend && npm run build`
- [ ] Test trên staging environment
- [ ] Deploy to production

---

## 📝 Notes

### Console Log Debugging
Console.log statements đã được thêm vào để debug (line 1493-1503):
```javascript
console.log('=== ADJUST DEBUG ===')
console.log('Current quantity:', currentIngredient.value.quantity)
console.log('New quantity:', adjustData.value.quantity)
console.log('Is increase?', isIncrease)
console.log('cost_per_unit before:', adjustData.value.cost_per_unit)
console.log('Data to send:', data)
console.log('===================')
```

**Cần xóa sau khi test xong!**

### Logic Summary

| Thao tác | Số lượng | Giá nhập | cost_per_unit gửi | Backend xử lý |
|----------|----------|----------|-------------------|---------------|
| Adjust giảm | 10→8 | 0 | 0 | Giữ nguyên giá |
| Adjust giảm | 10→8 | 30000 | 0 | Giữ nguyên giá |
| Adjust tăng | 8→12 | 0 | 0 | Giữ nguyên giá |
| Adjust tăng | 8→12 | 30000 | 30000 | Tính weighted avg |

---

**Status**: ✅ Code đã fix xong, chờ user test để confirm
**Date**: 2026-02-07
