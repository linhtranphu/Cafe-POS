# Batch Definition Form - Unit Selector Enhancement

## Tóm tắt
Cải thiện form "Định nghĩa batch" tại `http://localhost:5173/#/batch` để cho phép chọn đơn vị nguồn (source unit) từ danh sách các đơn vị tương thích, tương tự như trong form "Tạo menu mới".

## Vấn đề trước đây
- Đơn vị nguồn (source_unit) là readonly và tự động điền từ đơn vị kho của nguyên liệu
- Không thể chọn đơn vị khác (ví dụ: chọn "g" khi kho dùng "kg")
- Không hiển thị thông tin quy đổi đơn vị
- Trải nghiệm người dùng kém hơn so với form menu

## Giải pháp

### 1. Import Unit Conversion Composable
```javascript
import { useUnitConversion } from '../../composables/useUnitConversion'
const { getCompatibleUnits, getConversionRate, getConversionExplanation } = useUnitConversion()
```

### 2. Thay đổi Input thành Select Dropdown
**Trước:**
```vue
<input 
  v-model="rate.source_unit" 
  type="text" 
  :readonly="!!rate.source_ingredient_id"
  class="... bg-gray-100" 
/>
```

**Sau:**
```vue
<select 
  v-model="rate.source_unit" 
  @change="updateConversionRate(index)"
  :disabled="!rate.source_ingredient_id"
  class="..."
  :class="!rate.source_ingredient_id ? 'bg-gray-100' : 'bg-white'">
  <option value="">Chọn đơn vị</option>
  <option 
    v-for="unit in rate.compatibleUnits || []" 
    :key="unit" 
    :value="unit">
    {{ unit }}
  </option>
</select>
```

### 3. Thêm Conversion Info Display
```vue
<!-- Conversion Info -->
<div v-if="rate.source_ingredient_id && rate.source_unit && rate.conversionRate !== 1" 
  class="mb-2 p-2 bg-blue-50 rounded-lg text-xs text-blue-700">
  <span class="font-bold">ℹ️</span> 
  {{ getConversionExplanation(getSelectedIngredient(rate.source_ingredient_id)?.unit, rate.source_unit) }}
</div>
```

### 4. Cập nhật Data Structure
Thêm 2 field mới cho mỗi conversion rate:
```javascript
{
  source_ingredient_id: '',
  source_ingredient_name: '',
  source_quantity: 0,
  source_unit: '',
  batch_quantity: 0,
  wastage_rate_percent: 0,
  compatibleUnits: [],      // NEW: Danh sách đơn vị tương thích
  conversionRate: 1         // NEW: Tỷ lệ quy đổi
}
```

### 5. Thêm Helper Functions

#### updateIngredientName (cập nhật)
```javascript
const updateIngredientName = (index) => {
  const rate = formData.value.conversion_rates[index]
  const ingredient = ingredients.value.find(i => i.id === rate.source_ingredient_id)
  if (ingredient) {
    rate.source_ingredient_name = ingredient.name
    rate.source_unit = ingredient.unit
    rate.compatibleUnits = getCompatibleUnits(ingredient.unit)  // NEW
    rate.conversionRate = 1  // NEW
  }
}
```

#### updateConversionRate (mới)
```javascript
const updateConversionRate = (index) => {
  const rate = formData.value.conversion_rates[index]
  const ingredient = ingredients.value.find(i => i.id === rate.source_ingredient_id)
  if (ingredient && rate.source_unit) {
    rate.conversionRate = getConversionRate(ingredient.unit, rate.source_unit)
  }
}
```

## Tính năng mới

### 1. Chọn đơn vị nguồn linh hoạt
- Dropdown hiển thị các đơn vị tương thích
- Ví dụ: Nếu nguyên liệu dùng "kg", có thể chọn "kg" hoặc "g"
- Ví dụ: Nếu nguyên liệu dùng "L", có thể chọn "L" hoặc "ml"

### 2. Hiển thị thông tin quy đổi
- Khi chọn đơn vị khác với đơn vị kho
- Hiển thị công thức quy đổi (ví dụ: "1g = 0.001kg")
- Màu xanh dương để dễ nhận biết

### 3. Tự động tính toán
- Tỷ lệ quy đổi được tính tự động
- Cập nhật khi thay đổi đơn vị
- Sử dụng trong tính toán chi phí

## Các đơn vị được hỗ trợ

### Khối lượng (Mass)
- kg (Kilogram)
- g (Gram)
- Quy đổi: 1kg = 1000g

### Thể tích (Volume)
- L (Lít)
- ml (Milliliter)
- Quy đổi: 1L = 1000ml

### Đếm (Count)
- piece (Cái)
- box (Hộp)
- pack (Gói)
- Không quy đổi giữa các đơn vị này

## Ví dụ sử dụng

### Ví dụ 1: Tạo batch cà phê concentrate
1. Chọn nguyên liệu: "Cà phê hạt" (kho: kg)
2. Chọn đơn vị nguồn: "g" (từ dropdown)
3. Nhập số lượng: 100g
4. Hệ thống hiển thị: "ℹ️ 1g = 0.001kg"
5. Tính toán chi phí dựa trên quy đổi

### Ví dụ 2: Tạo batch sữa tươi
1. Chọn nguyên liệu: "Sữa tươi" (kho: L)
2. Chọn đơn vị nguồn: "ml" (từ dropdown)
3. Nhập số lượng: 500ml
4. Hệ thống hiển thị: "ℹ️ 1ml = 0.001L"
5. Tính toán chi phí dựa trên quy đổi

## Files thay đổi

### frontend/src/components/batch/BatchDefinitionForm.vue
1. Import useUnitConversion composable
2. Thay input thành select dropdown
3. Thêm conversion info display
4. Cập nhật updateIngredientName function
5. Thêm updateConversionRate function
6. Cập nhật addConversionRate function
7. Cập nhật watch function cho edit mode

## Lợi ích

### Cho người dùng
1. ✅ Linh hoạt hơn trong việc chọn đơn vị
2. ✅ Dễ dàng làm việc với đơn vị nhỏ hơn (g, ml)
3. ✅ Hiểu rõ cách quy đổi đơn vị
4. ✅ Trải nghiệm nhất quán với form menu

### Cho hệ thống
1. ✅ Tính toán chi phí chính xác hơn
2. ✅ Hỗ trợ nhiều đơn vị đo lường
3. ✅ Code dễ bảo trì và mở rộng
4. ✅ Tái sử dụng logic từ useUnitConversion

## Testing Checklist

- [x] No syntax errors
- [ ] Dropdown hiển thị đúng các đơn vị tương thích
- [ ] Chọn đơn vị cập nhật conversionRate
- [ ] Conversion info hiển thị đúng
- [ ] Chi phí ước tính tính toán chính xác
- [ ] Edit mode load đúng compatibleUnits
- [ ] Validation hoạt động với đơn vị mới
- [ ] Save batch definition thành công

## Hướng dẫn test

1. Mở `http://localhost:5173/#/batch`
2. Click "➕ Tạo Batch Mới"
3. Điền tên batch và các thông tin cơ bản
4. Click "+ Thêm" trong phần "📦 Nguyên Liệu Nguồn"
5. Chọn một nguyên liệu (ví dụ: Cà phê hạt - kg)
6. Kiểm tra dropdown "Đơn vị nguồn" hiển thị: kg, g
7. Chọn "g" từ dropdown
8. Kiểm tra hiển thị: "ℹ️ 1g = 0.001kg"
9. Nhập số lượng và kiểm tra chi phí ước tính
10. Lưu batch definition và verify

## Ghi chú kỹ thuật

### Unit Conversion Logic
- Sử dụng `useUnitConversion` composable
- Hỗ trợ quy đổi trong cùng category (mass-to-mass, volume-to-volume)
- Không hỗ trợ quy đổi giữa các category khác nhau
- Conversion rate được tính dựa trên base unit (kg cho mass, L cho volume)

### Data Flow
1. User chọn ingredient → Load compatibleUnits
2. User chọn unit → Calculate conversionRate
3. User nhập quantity → Calculate estimated cost
4. User save → Send to backend với source_unit đã chọn

### Backend Compatibility
- Backend đã hỗ trợ source_unit field
- Không cần thay đổi API
- Frontend chỉ cần gửi đúng source_unit đã chọn

## Status

✅ **COMPLETE** - Tính năng đã được implement và không có lỗi syntax
