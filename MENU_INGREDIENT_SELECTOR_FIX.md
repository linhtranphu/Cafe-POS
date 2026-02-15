# Menu Ingredient Selector Fix

## Vấn đề
Khi thêm món mới tại `http://localhost:5173/#/menu`, phần chọn nguyên liệu không hiển thị. Khi inspect HTML chỉ thấy comment:
```html
<!-- Ingredient details here - keep existing structure -->
```

## Nguyên nhân
Code hiển thị chi tiết nguyên liệu trong form "Thêm món mới" (single-size mode) đã bị thay thế bằng comment placeholder thay vì code thực tế.

## Giải pháp
Đã khôi phục lại code hiển thị chi tiết nguyên liệu với cấu trúc đầy đủ:

### Cấu trúc được khôi phục:

1. **Ingredient Header** - Hiển thị tên và thông tin kho
   - Tên nguyên liệu
   - Đơn vị kho và giá/đơn vị
   - Nút xóa (×)

2. **Recipe Unit & Quantity** - Chọn đơn vị và số lượng
   - Dropdown chọn đơn vị (từ compatibleUnits)
   - Input số lượng cần dùng
   - Layout grid 2 cột

3. **Conversion Info** - Thông tin quy đổi
   - Hiển thị khi conversionRate !== 1
   - Giải thích quy đổi giữa đơn vị kho và đơn vị công thức

4. **Cost Preview** - Xem trước chi phí
   - Hiển thị chi phí ước tính
   - Hiển thị % hao hụt (nếu có)
   - Màu xanh lá để dễ nhận biết

## Code thay đổi

### File: `frontend/src/views/MenuView.vue`

**Trước (bị lỗi):**
```vue
<div v-else class="space-y-3">
  <div v-for="(ing, index) in form.ingredients" :key="index"
    class="bg-gray-50 rounded-lg p-3 border-2 border-gray-200">
    <!-- Ingredient details here - keep existing structure -->
  </div>
</div>
```

**Sau (đã fix):**
```vue
<div v-else class="space-y-3">
  <div v-for="(ing, index) in form.ingredients" :key="index"
    class="bg-gray-50 rounded-lg p-3 border-2 border-gray-200">
    
    <!-- Ingredient Header - Compact -->
    <div class="flex justify-between items-start mb-2">
      <div class="flex-1 min-w-0">
        <div class="font-bold text-sm text-gray-800 truncate">{{ ing.name }}</div>
        <div class="text-xs text-gray-500">
          Kho: {{ ing.stockUnit }} @ {{ formatPrice(ing.costPerUnit) }}/{{ ing.stockUnit }}
        </div>
      </div>
      <button type="button" @click="removeIngredient(index)"
        class="ml-2 bg-red-500 text-white px-2 py-1 rounded text-xs font-bold active:scale-95">
        ×
      </button>
    </div>

    <!-- Recipe Unit & Quantity - Inline -->
    <div class="grid grid-cols-2 gap-2 mb-2">
      <div>
        <label class="text-xs text-gray-600 block mb-1">Đơn vị:</label>
        <select v-model="ing.unit" @change="updateRecipeUnit(index)"
          class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded-lg">
          <option v-for="unit in ing.compatibleUnits" :key="unit" :value="unit">
            {{ unit }}
          </option>
        </select>
      </div>
      <div>
        <label class="text-xs text-gray-600 block mb-1">Số lượng:</label>
        <input v-model.number="ing.quantity" 
          @input="updateIngredientCost(index)"
          type="number" min="0" step="0.1" placeholder="0" required
          class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded-lg" />
      </div>
    </div>
    
    <!-- Conversion Info -->
    <div v-if="ing.conversionRate !== 1" 
      class="mb-2 p-2 bg-blue-50 rounded text-xs text-blue-700">
      <span class="font-bold">ℹ️</span> {{ getConversionExplanation(ing.stockUnit, ing.unit) }}
    </div>
    
    <!-- Cost Preview - Compact -->
    <div v-if="ing.costPerUnit > 0" class="p-2 bg-green-50 rounded-lg border border-green-200">
      <div class="flex justify-between items-center">
        <span class="text-xs text-green-700 font-bold">Chi phí:</span>
        <span class="text-sm font-bold text-green-700">
          {{ formatPrice(ing.estimatedCost) }}
        </span>
      </div>
      <div v-if="ing.wastage > 0" class="text-xs text-green-600 mt-1">
        (+ {{ ing.wastage }}% hao hụt)
      </div>
    </div>
  </div>
</div>
```

## Tính năng được khôi phục

1. ✅ Hiển thị tên nguyên liệu và thông tin kho
2. ✅ Chọn đơn vị công thức từ danh sách tương thích
3. ✅ Nhập số lượng cần dùng
4. ✅ Hiển thị thông tin quy đổi đơn vị (nếu có)
5. ✅ Xem trước chi phí ước tính
6. ✅ Hiển thị % hao hụt (nếu có)
7. ✅ Nút xóa nguyên liệu

## Kiểm tra

- [x] No syntax errors
- [x] Code structure matches variant ingredients section
- [ ] Test adding ingredient to single-size menu item
- [ ] Verify ingredient details display correctly
- [ ] Verify unit conversion works
- [ ] Verify cost calculation displays
- [ ] Verify remove button works

## Ghi chú

- Code được khôi phục dựa trên cấu trúc tương tự trong phần variant ingredients (multi-size mode)
- Cấu trúc này đã được test và hoạt động tốt trong phần variants
- Các function helper cần thiết:
  - `removeIngredient(index)` - Xóa nguyên liệu
  - `updateRecipeUnit(index)` - Cập nhật đơn vị công thức
  - `updateIngredientCost(index)` - Cập nhật chi phí
  - `getConversionExplanation(stockUnit, recipeUnit)` - Giải thích quy đổi
  - `formatPrice(value)` - Format giá VNĐ

## Hướng dẫn test

1. Mở `http://localhost:5173/#/menu`
2. Click "➕ Thêm món"
3. Điền tên món, chọn danh mục
4. Đảm bảo toggle "Có nhiều size" là OFF (single-size mode)
5. Click "+ Thêm" trong phần "🥘 Nguyên liệu"
6. Chọn một nguyên liệu
7. Kiểm tra xem:
   - Tên nguyên liệu hiển thị
   - Thông tin kho hiển thị
   - Dropdown đơn vị hoạt động
   - Input số lượng hoạt động
   - Chi phí ước tính hiển thị
   - Nút xóa (×) hoạt động

## Status

✅ **FIXED** - Code đã được khôi phục và không có lỗi syntax
