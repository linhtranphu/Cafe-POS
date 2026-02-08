# Fix: Form "Sửa" Không Được Thay Đổi Tồn Kho và Giá

## Vấn Đề Phát Hiện

User click nút "Sửa" (✏️) ở danh sách nguyên liệu, thay đổi số lượng từ 10kg xuống 8kg, và đơn giá bị tính lại khi save.

### Nguyên Nhân

Form "Sửa" (Edit Ingredient) cho phép thay đổi:
- ✅ Tên nguyên liệu
- ✅ Danh mục
- ✅ Đơn vị
- ❌ **Số lượng tồn kho** (KHÔNG NÊN)
- ❌ **Đơn giá** (KHÔNG NÊN)
- ✅ Mức tối thiểu

**Vấn đề**: Số lượng và giá chỉ nên thay đổi qua chức năng "Điều chỉnh tồn kho", không phải qua form "Sửa".

## Giải Pháp

### 1. Disable Field "Số lượng" Khi Edit

```vue
<!-- Trước -->
<input v-model.number="formData.quantity" type="number" />

<!-- Sau -->
<input v-model.number="formData.quantity" type="number" 
  :disabled="isEditing"
  :class="isEditing ? 'bg-gray-100 cursor-not-allowed' : ''" />
```

**Kết quả**: Khi edit, field số lượng bị disable (màu xám, không thể thay đổi).

### 2. Ẩn Section "Thông tin giá" Khi Edit

```vue
<!-- Trước -->
<div class="bg-gray-50 rounded-xl p-4 space-y-4">
  <h3>💰 Thông tin giá</h3>
  <!-- Price inputs -->
</div>

<!-- Sau -->
<div v-if="!isEditing" class="bg-gray-50 rounded-xl p-4 space-y-4">
  <h3>💰 Thông tin giá</h3>
  <!-- Price inputs -->
</div>
```

**Kết quả**: Khi edit, section nhập giá bị ẩn hoàn toàn.

### 3. Hiển thị Thông Tin Read-Only Khi Edit

```vue
<!-- Thêm section mới -->
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

**Kết quả**: Khi edit, hiển thị thông tin tồn kho và giá ở chế độ read-only với cảnh báo.

## So Sánh Trước và Sau

### Trước Fix

**Form "Sửa"** (Edit):
- ✏️ Tên: Có thể sửa
- ✏️ Danh mục: Có thể sửa
- ✏️ Đơn vị: Có thể sửa
- ✏️ **Số lượng: Có thể sửa** ❌
- ✏️ **Đơn giá: Có thể sửa** ❌
- ✏️ Mức tối thiểu: Có thể sửa

**Vấn đề**: User có thể vô tình thay đổi tồn kho và giá qua form "Sửa".

### Sau Fix

**Form "Sửa"** (Edit):
- ✏️ Tên: Có thể sửa
- ✏️ Danh mục: Có thể sửa
- ✏️ Đơn vị: Có thể sửa
- 👁️ **Số lượng: Chỉ xem (disabled)** ✅
- 👁️ **Đơn giá: Chỉ xem (ẩn input, hiển thị info)** ✅
- ✏️ Mức tối thiểu: Có thể sửa

**Lợi ích**: 
- User không thể vô tình thay đổi tồn kho và giá
- Rõ ràng phải dùng "Điều chỉnh" để thay đổi tồn kho
- Tránh nhầm lẫn giữa "Sửa" và "Điều chỉnh"

## Phân Biệt 2 Chức Năng

### Nút "Sửa" (✏️)
**Mục đích**: Sửa thông tin cơ bản của nguyên liệu
**Có thể thay đổi**:
- Tên nguyên liệu
- Danh mục
- Đơn vị
- Mức tối thiểu

**KHÔNG thể thay đổi**:
- Số lượng tồn kho
- Đơn giá

### Nút "Điều chỉnh" (📦)
**Mục đích**: Điều chỉnh tồn kho và giá
**Có thể thay đổi**:
- Số lượng tồn kho (nhập/xuất/điều chỉnh)
- Đơn giá (khi nhập hàng mới)

**Có tracking**:
- Lịch sử nhập/xuất
- Lý do điều chỉnh
- Người thực hiện

## UI Changes

### Form "Sửa" - Trước Fix
```
┌─────────────────────────────────┐
│ ✏️ Sửa nguyên liệu              │
├─────────────────────────────────┤
│ Tên: [Đường]                    │
│ Danh mục: [Nguyên liệu khô]     │
│ Đơn vị: [kg]                    │
│                                 │
│ 💰 Thông tin giá                │
│ Đơn giá: [25000] ← Có thể sửa  │
│                                 │
│ Số lượng: [10] ← Có thể sửa ❌  │
│ Mức tối thiểu: [5]              │
│                                 │
│ [Hủy]  [Lưu]                    │
└─────────────────────────────────┘
```

### Form "Sửa" - Sau Fix
```
┌─────────────────────────────────┐
│ ✏️ Sửa nguyên liệu              │
├─────────────────────────────────┤
│ Tên: [Đường]                    │
│ Danh mục: [Nguyên liệu khô]     │
│ Đơn vị: [kg]                    │
│                                 │
│ 📊 Thông tin hiện tại (chỉ xem) │
│ ┌──────────┬──────────┐         │
│ │ Tồn kho  │ Đơn giá  │         │
│ │ 10 kg    │ 25,000đ  │         │
│ └──────────┴──────────┘         │
│ ⚠️ Dùng "Điều chỉnh" để thay đổi│
│                                 │
│ Số lượng: [10] (disabled) ✅    │
│ ⚠️ Dùng "Điều chỉnh" để thay đổi│
│                                 │
│ Mức tối thiểu: [5]              │
│                                 │
│ [Hủy]  [Lưu]                    │
└─────────────────────────────────┘
```

## Testing

### Test Case 1: Edit Tên Nguyên Liệu
1. Click nút "Sửa" (✏️)
2. Thay đổi tên: "Đường" → "Đường trắng"
3. Click "Lưu"
4. **Kết quả**: Tên được cập nhật, tồn kho và giá không đổi ✅

### Test Case 2: Không Thể Sửa Tồn Kho
1. Click nút "Sửa" (✏️)
2. Thử thay đổi số lượng
3. **Kết quả**: Field bị disable, không thể thay đổi ✅

### Test Case 3: Không Thể Sửa Giá
1. Click nút "Sửa" (✏️)
2. Tìm field nhập giá
3. **Kết quả**: Không có field nhập giá, chỉ hiển thị thông tin read-only ✅

### Test Case 4: Thay Đổi Tồn Kho Qua "Điều chỉnh"
1. Click nút "Điều chỉnh" (📦)
2. Chọn loại: "Điều chỉnh"
3. Nhập số lượng mới: 8kg
4. Nhập lý do
5. Click "Xác nhận"
6. **Kết quả**: Tồn kho được cập nhật, có lịch sử tracking ✅

## Files Thay Đổi

### `frontend/src/views/IngredientManagementView.vue`

**1. Field Số lượng** (Line ~310):
- Thêm `:disabled="isEditing"`
- Thêm `:class="isEditing ? 'bg-gray-100 cursor-not-allowed' : ''"`
- Thêm warning text khi edit

**2. Section Thông tin giá** (Line ~253):
- Thêm `v-if="!isEditing"` để ẩn khi edit

**3. Section Read-only** (Line ~301):
- Thêm `v-else` để hiển thị khi edit
- Hiển thị tồn kho và giá ở chế độ read-only
- Thêm warning text

## Kết Luận

✅ **Fix hoàn thành**: Form "Sửa" giờ chỉ cho phép sửa thông tin cơ bản

✅ **Tách biệt rõ ràng**: "Sửa" vs "Điều chỉnh"

✅ **Tránh nhầm lẫn**: User không thể vô tình thay đổi tồn kho/giá qua form "Sửa"

✅ **UX tốt hơn**: Hiển thị thông tin read-only với cảnh báo rõ ràng

---

**Ngày fix**: 2026-02-07  
**File**: `frontend/src/views/IngredientManagementView.vue`  
**Chức năng**: Edit Ingredient Form

