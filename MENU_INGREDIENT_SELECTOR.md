# Menu - Ingredient Selector Integration

## Tóm tắt
Đã cập nhật MenuView để chọn nguyên liệu từ danh sách có sẵn thay vì nhập thủ công, giúp việc thêm nguyên liệu vào món ăn thuận tiện và chính xác hơn.

## Thay đổi

### Frontend: `frontend/src/views/MenuView.vue`

#### 1. Import thêm Ingredient Store
```javascript
import { useIngredientStore } from '../stores/ingredient'
const ingredientStore = useIngredientStore()
```

#### 2. Thêm State mới
- `showIngredientSelector` - Hiển thị modal chọn nguyên liệu
- `ingredientSearchQuery` - Tìm kiếm nguyên liệu trong modal
- `availableIngredients` - Danh sách nguyên liệu từ store
- `ingredientsLoading` - Loading state cho ingredients

#### 3. Computed Properties mới
- `filteredAvailableIngredients` - Lọc nguyên liệu theo search query

#### 4. Functions mới
- `selectIngredient(ingredient)` - Chọn nguyên liệu từ danh sách
- `isIngredientSelected(ingredientId)` - Kiểm tra nguyên liệu đã được chọn chưa

#### 5. UI Changes

**Trước (Nhập thủ công):**
```
📋 Nguyên liệu
┌─────────────────────────────────┐
│ Tên: [_____________]      [×]   │
│ Số lượng: [___] Đơn vị: [___]   │
└─────────────────────────────────┘
[+ Thêm nguyên liệu]
```

**Sau (Chọn từ danh sách):**
```
📋 Nguyên liệu
┌─────────────────────────────────┐
│ Cà phê                          │
│ gram                       [×]  │
│ Số lượng: [100] gram            │
└─────────────────────────────────┘
[+ Chọn nguyên liệu]
  ↓
┌─────────────────────────────────┐
│ 🥬 Chọn nguyên liệu        [×]  │
├─────────────────────────────────┤
│ [Tìm kiếm...]                   │
├─────────────────────────────────┤
│ ┌─────────────────────────────┐ │
│ │ Cà phê                  [+] │ │
│ │ Nguyên liệu khô             │ │
│ │ Tồn kho: 5000 gram          │ │
│ └─────────────────────────────┘ │
│ ┌─────────────────────────────┐ │
│ │ Sữa tươi                [✓] │ │
│ │ Nguyên liệu tươi sống       │ │
│ │ Tồn kho: 20 lít             │ │
│ └─────────────────────────────┘ │
└─────────────────────────────────┘
```

## Tính năng

### Ingredient Selector Modal
- ✅ Hiển thị danh sách tất cả nguyên liệu có sẵn
- ✅ Tìm kiếm nguyên liệu theo tên hoặc danh mục
- ✅ Hiển thị thông tin: tên, danh mục, tồn kho, đơn vị
- ✅ Đánh dấu nguyên liệu đã được chọn (✓)
- ✅ Disable nguyên liệu đã chọn để tránh trùng lặp
- ✅ Tự động đóng modal sau khi chọn

### Ingredient Display
- ✅ Hiển thị tên nguyên liệu (từ database)
- ✅ Hiển thị đơn vị (từ database)
- ✅ Chỉ cần nhập số lượng
- ✅ Xóa nguyên liệu dễ dàng

### Data Flow
```
1. User clicks "Chọn nguyên liệu"
   ↓
2. Modal opens with all ingredients
   ↓
3. User searches/selects ingredient
   ↓
4. Ingredient added to form with:
   - id (from database)
   - name (from database)
   - unit (from database)
   - quantity (default: 1, user can edit)
   ↓
5. Modal closes
   ↓
6. User edits quantity if needed
   ↓
7. Save menu item with ingredients
```

## Lợi ích

### 1. Chính xác hơn
- Tên nguyên liệu chuẩn từ database
- Đơn vị đúng, không bị nhầm lẫn
- Có ID để link với ingredient management

### 2. Thuận tiện hơn
- Không cần nhớ tên chính xác
- Tìm kiếm nhanh
- Xem được tồn kho ngay khi chọn
- Không thể chọn trùng

### 3. Dễ quản lý
- Link trực tiếp với ingredient database
- Có thể tracking ingredient usage
- Dễ dàng update khi ingredient thay đổi

## Testing Checklist

### Ingredient Selector
- [ ] Modal mở khi click "Chọn nguyên liệu"
- [ ] Hiển thị đầy đủ danh sách ingredients
- [ ] Search hoạt động đúng
- [ ] Hiển thị thông tin ingredient đầy đủ
- [ ] Chọn ingredient thêm vào form
- [ ] Ingredient đã chọn bị disable
- [ ] Modal đóng sau khi chọn

### Ingredient Management
- [ ] Hiển thị ingredient đã chọn đúng
- [ ] Edit quantity hoạt động
- [ ] Xóa ingredient hoạt động
- [ ] Lưu menu item với ingredients
- [ ] Load menu item với ingredients

### Edge Cases
- [ ] Không có ingredient nào
- [ ] Search không tìm thấy
- [ ] Chọn nhiều ingredients
- [ ] Xóa tất cả ingredients
- [ ] Edit menu item có ingredients

## Hướng dẫn sử dụng

### Thêm món mới với nguyên liệu

1. Click "Thêm món"
2. Điền tên món, danh mục, giá
3. Scroll xuống phần "🥘 Nguyên liệu"
4. Click "Chọn nguyên liệu"
5. Tìm kiếm hoặc scroll để tìm nguyên liệu
6. Click vào nguyên liệu muốn thêm
7. Nhập số lượng cần dùng
8. Lặp lại bước 4-7 để thêm nguyên liệu khác
9. Click "Thêm món" để lưu

### Sửa nguyên liệu trong món

1. Click vào món cần sửa
2. Scroll xuống phần nguyên liệu
3. Để xóa: Click nút [×] bên cạnh nguyên liệu
4. Để thêm: Click "Chọn nguyên liệu"
5. Để sửa số lượng: Nhập số mới vào ô số lượng
6. Click "Cập nhật" để lưu

## Cấu trúc dữ liệu

### Ingredient trong form
```javascript
{
  id: "507f1f77bcf86cd799439011",  // ID từ ingredient database
  name: "Cà phê",                    // Tên từ database
  quantity: 20,                      // User nhập
  unit: "gram"                       // Đơn vị từ database
}
```

### Menu Item với ingredients
```javascript
{
  name: "Cà phê sữa đá",
  category: "Cà phê",
  price: 25000,
  ingredients: [
    {
      id: "507f1f77bcf86cd799439011",
      name: "Cà phê",
      quantity: 20,
      unit: "gram"
    },
    {
      id: "507f1f77bcf86cd799439012",
      name: "Sữa tươi",
      quantity: 100,
      unit: "ml"
    }
  ]
}
```

## Files đã thay đổi

- `frontend/src/views/MenuView.vue` - Thêm ingredient selector modal và logic

## Next Steps (Optional)

### Tính năng nâng cao có thể thêm:
1. **Ingredient Categories Filter** - Lọc theo danh mục trong selector
2. **Recent Ingredients** - Hiển thị nguyên liệu hay dùng ở đầu
3. **Bulk Add** - Chọn nhiều ingredients cùng lúc
4. **Ingredient Cost Display** - Hiển thị giá nguyên liệu khi chọn
5. **Recipe Templates** - Lưu công thức để dùng lại
6. **Ingredient Suggestions** - Gợi ý nguyên liệu dựa trên category món

## Kết luận

✅ **Hoàn thành**: Ingredient selector đã được tích hợp vào MenuView
✅ **Hoàn thành**: Frontend build thành công
✅ **Cải thiện**: UX tốt hơn, chính xác hơn, dễ quản lý hơn
⏳ **Tiếp theo**: Test thực tế với dữ liệu thật
