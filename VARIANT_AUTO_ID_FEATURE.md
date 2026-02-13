# Tính năng Tự động tạo ID cho Variants và Fix Ingredient Unit

## Tổng quan

Đã cập nhật giao diện tạo menu để:
1. Tự động tạo ID cho các size variants - người dùng chỉ cần nhập tên size (label) và giá
2. Fix lỗi không lưu đơn vị công thức (unit) của nguyên liệu trong variants

## Thay đổi

### 1. Auto-generate Variant ID (MenuView.vue)

**Xóa trường nhập ID**: Đã xóa input field cho variant ID khỏi form
**Tự động tạo ID**: Thêm hàm `generateVariantId()` để tạo ID unique
**Cải thiện UX**: Form giờ chỉ có 2 trường thay vì 3 (Tên size và Giá)

### 2. Fix Ingredient Unit Saving (MenuView.vue)

**Vấn đề**: Khi edit menu item có variants, các trường bổ sung của ingredient (unit, stockUnit, compatibleUnits, etc.) không được load từ backend, dẫn đến:
- Dropdown đơn vị công thức không hiển thị
- Unit không được lưu khi update

**Giải pháp**: Enrich ingredient data khi load menu item để edit:
- Tìm ingredient data từ store
- Thêm các trường: stockUnit, compatibleUnits, costPerUnit, wastage, conversionRate
- Tính toán lại conversion rate dựa trên unit đã lưu

```javascript
const editItem = (item) => {
  // ... 
  // Prepare variants with enriched ingredient data
  const preparedVariants = item.variants ? item.variants.map(variant => {
    const enrichedIngredients = variant.ingredients ? variant.ingredients.map(ing => {
      const ingredientData = ingredients.value.find(i => i.name === ing.name)
      if (ingredientData) {
        const compatibleUnits = getCompatibleUnits(ingredientData.unit)
        const conversionRate = getConversionRate(ingredientData.unit, ing.unit)
        return {
          id: ingredientData.id,
          name: ing.name,
          quantity: ing.quantity,
          unit: ing.unit,  // ✅ Preserved from backend
          stockUnit: ingredientData.unit,
          compatibleUnits: compatibleUnits,
          costPerUnit: ingredientData.cost_per_unit || 0,
          wastage: ingredientData.wastage_percentage || 0,
          conversionRate: conversionRate,
          estimatedCost: 0
        }
      }
      return ing
    }) : []
    
    return {
      ...variant,
      ingredients: enrichedIngredients
    }
  }) : []
  // ...
}
```

## Cách hoạt động

### Auto-generate Variant ID
```javascript
const generateVariantId = () => {
  // Use timestamp + random number for uniqueness
  const timestamp = Date.now().toString(36) // Convert to base36 for shorter string
  const random = Math.random().toString(36).substring(2, 7) // 5 random chars
  return `${timestamp}-${random}`
}
```

**Ví dụ ID được tạo**: `lm8x9k2-7h3j5`, `lm8x9k3-9a2b4`

### Ingredient Unit Flow
1. **Thêm nguyên liệu mới**: Unit được set từ ingredient selector
2. **Lưu vào backend**: Unit được gửi trong payload
3. **Load lại để edit**: Unit được enrich từ backend data + ingredient store
4. **Hiển thị dropdown**: Compatible units được tính toán lại
5. **Update unit**: Conversion rate được tính lại khi user thay đổi unit

## Tính năng

- ✅ ID được tạo tự động khi thêm variant mới
- ✅ ID là unique (dùng timestamp + random)
- ✅ Khi edit menu item, ID cũ được giữ nguyên
- ✅ Đơn vị công thức được lưu và hiển thị đúng
- ✅ Dropdown đơn vị hoạt động khi edit
- ✅ Conversion rate được tính toán đúng
- ✅ Backend validation vẫn hoạt động bình thường
- ✅ Không breaking change - tương thích với dữ liệu cũ

## Giao diện trước và sau

### Trước (3 trường):
```
┌─────────────────────────────────────────────┐
│ ID *          │ Tên *        │ Giá (VNĐ) *  │
│ [M____]       │ [Size M___]  │ [30000____]  │
└─────────────────────────────────────────────┘
```

### Sau (2 trường):
```
┌─────────────────────────────────────────────┐
│ Tên size *              │ Giá (VNĐ) *       │
│ [Size M______________]  │ [30000__________] │
└─────────────────────────────────────────────┘
```

## Testing

### Manual Test - Auto ID
1. Mở frontend: http://localhost:5173/#/menu
2. Nhấn "Thêm món mới"
3. Chọn "Có nhiều size"
4. Nhấn "Thêm size"
5. Chỉ cần nhập:
   - Tên size (ví dụ: "Size M")
   - Giá (ví dụ: 30000)
6. Lưu và kiểm tra - ID sẽ được tạo tự động

### Manual Test - Ingredient Unit
1. Tạo menu item multi-size với nguyên liệu
2. Chọn đơn vị công thức khác với đơn vị kho (ví dụ: kho là "kg", công thức là "g")
3. Lưu menu item
4. Edit lại menu item
5. Kiểm tra:
   - ✅ Dropdown đơn vị hiển thị đúng
   - ✅ Đơn vị đã chọn được giữ nguyên
   - ✅ Conversion rate hiển thị đúng

### API Test
```bash
# Test auto-generated IDs
./test-variant-auto-id.sh

# Test ingredient unit saving
./test-variant-ingredient-unit.sh
```

## Files thay đổi

- `frontend/src/views/MenuView.vue`: 
  - Xóa ID input field
  - Thêm auto-generation cho variant ID
  - Fix ingredient data enrichment trong editItem()
- `test-variant-auto-id.sh`: Test script cho auto ID
- `test-variant-ingredient-unit.sh`: Test script cho ingredient unit
- `VARIANT_AUTO_ID_FEATURE.md`: Tài liệu này

## Lưu ý

- ID vẫn được lưu trong database và sử dụng cho logic backend
- Người dùng không thấy ID trong giao diện
- Khi edit, ID cũ được giữ nguyên (không tạo ID mới)
- Format ID: `{timestamp-base36}-{random-5chars}`
- Ingredient unit được lưu chính xác trong backend
- Khi edit, ingredient data được enrich để hiển thị đầy đủ thông tin

## Tương thích ngược

✅ Hoàn toàn tương thích với:
- Menu items cũ (single-size)
- Menu items có variants với ID thủ công
- Menu items có ingredients với units khác nhau
- Tất cả API endpoints hiện tại
