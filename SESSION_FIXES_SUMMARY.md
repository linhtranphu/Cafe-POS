# Session Fixes Summary

## 1. Batch Unit Conversion - Cho phép chuyển đổi đơn vị

### Vấn đề
Batch ingredients chỉ cho phép dùng đơn vị gốc (ví dụ: batch "L" chỉ dùng được "L", không dùng được "ml").

### Giải pháp
Sửa function `selectBatch()` trong `MenuView.vue` để sử dụng `getCompatibleUnits()`:

**File:** `frontend/src/views/MenuView.vue`
```javascript
// Thay vì:
compatibleUnits: [batch.unit]

// Dùng:
const compatibleUnits = getCompatibleUnits(batch.unit)
compatibleUnits: compatibleUnits
```

### Kết quả
- ✅ Batch "L" → có thể chọn "L" hoặc "ml"
- ✅ Batch "kg" → có thể chọn "kg" hoặc "g"
- ✅ Chi phí tự động tính đúng

---

## 2. Batch Ingredient Edit - Giữ đơn vị khi edit

### Vấn đề
Sau khi lưu món với batch ingredient, click "Sửa" thì dropdown đơn vị bị biến mất.

### Nguyên nhân
Function `editItem()` chỉ xử lý raw ingredients, không xử lý batch ingredients.

### Giải pháp
Sửa function `editItem()` để xử lý cả batch ingredients:

**File:** `frontend/src/views/MenuView.vue`
```javascript
const editItem = (item) => {
  // ...
  const preparedIngredients = item.ingredients ? item.ingredients.map(ing => {
    const isBatch = ing.ingredient_type === 'batch' || ing.type === 'batch'
    
    if (isBatch) {
      // Tìm batch data và tính compatibleUnits
      const batchData = availableBatchDefinitions.value.find(...)
      const compatibleUnits = getCompatibleUnits(batchData.unit)
      // ...
    } else {
      // Xử lý raw ingredient
    }
  })
}
```

### Kết quả
- ✅ Dropdown đơn vị hiển thị đầy đủ khi edit
- ✅ Đơn vị đã chọn được giữ nguyên

---

## 3. Batch Ingredient Save - Lưu ingredient_type và batch_id

### Vấn đề
Khi save menu item, không gửi `ingredient_type` và `batch_id` lên backend → khi load lại không biết đây là batch ingredient.

### Giải pháp
Sửa function `saveItem()` để gửi đầy đủ fields:

**File:** `frontend/src/views/MenuView.vue`
```javascript
ingredients: v.ingredients.map(ing => ({
  name: ing.name,
  quantity: ing.quantity,
  unit: ing.unit,
  // Thêm batch fields
  ...(ing.type === 'batch' && {
    ingredient_type: 'batch',
    batch_id: ing.batch_definition_id || ing.id
  })
}))
```

### Kết quả
- ✅ Backend lưu đúng `ingredient_type: "batch"` và `batch_id`
- ✅ Khi load lại, frontend biết đây là batch ingredient

---

## 4. Order Cart Helpers - Fix createCartItem error

### Vấn đề
Lỗi khi thêm món vào giỏ hàng:
```
TypeError: Cannot read properties of undefined (reading 'createCartItem')
```

### Nguyên nhân
Trong Pinia store, `helpers` object không phải là property hợp lệ → `orderStore.helpers` là `undefined`.

### Giải pháp
Export `cartHelpers` riêng biệt, không phải part của store:

**File:** `frontend/src/stores/order.js`
```javascript
export const useOrderStore = defineStore('order', {
  state: () => ({ ... }),
  actions: { ... },
  getters: { ... }
})

// Export helpers separately
export const cartHelpers = {
  createCartItem() { ... },
  isSameCartItem() { ... }
}
```

**File:** `frontend/src/views/OrderView.vue`
```javascript
import { useOrderStore, cartHelpers } from '../stores/order'

// Sử dụng
const cartItem = cartHelpers.createCartItem(item, variant)
```

### Kết quả
- ✅ Có thể thêm món vào giỏ hàng
- ✅ Single-size và multi-size items đều hoạt động

---

## 5. User Role Barista - Cho phép tạo user barista

### Vấn đề
Không thể tạo user với role "barista" - backend trả về 400 Bad Request.

### Nguyên nhân
Function `isValidRole()` trong backend chỉ cho phép 3 roles: Manager, Cashier, Waiter - thiếu Barista.

### Giải pháp
Thêm `RoleBarista` vào danh sách valid roles:

**File:** `backend/application/services/user_management_service.go`
```go
func isValidRole(role user.Role) bool {
    validRoles := []user.Role{
        user.RoleManager, 
        user.RoleCashier, 
        user.RoleWaiter, 
        user.RoleBarista  // ✅ Thêm dòng này
    }
    // ...
}
```

### Kết quả
- ✅ Có thể tạo user với role "barista"
- ✅ Frontend đã có đầy đủ UI cho barista role

---

## 6. Order Category Tabs - Map đúng với danh mục menu

### Vấn đề
Category tabs trong OrderView bị hardcode với các id như 'coffee', 'tea', 'juice' - không khớp với categories thực tế trong database (ví dụ: "Cà phê", "Trà").

### Giải pháp
Lấy categories động từ menu items thay vì hardcode:

**File:** `frontend/src/views/OrderView.vue`
```javascript
// Thay vì hardcode:
const categories = [
  { id: 'all', name: 'Tất cả', icon: '📋' },
  { id: 'coffee', name: 'Cà phê', icon: '☕' },
  // ...
]

// Dùng computed để lấy từ menu items:
const categories = computed(() => {
  const allCategory = { id: 'all', name: 'Tất cả', icon: '📋' }
  
  // Get unique categories from menu items
  const uniqueCategories = [...new Set(menuItems.value.map(item => item.category))]
  
  const menuCategories = uniqueCategories.map(cat => ({
    id: cat,
    name: cat,
    icon: categoryIcons[cat] || '🍽️'
  }))
  
  return [allCategory, ...menuCategories]
})
```

### Kết quả
- ✅ Category tabs hiển thị đúng các danh mục có trong menu
- ✅ Tự động cập nhật khi thêm/xóa danh mục
- ✅ Filter menu items hoạt động đúng

---

## Files Changed

### Frontend
1. `frontend/src/views/MenuView.vue` - 3 fixes (batch unit conversion, edit, save)
2. `frontend/src/stores/order.js` - 1 fix (cart helpers)
3. `frontend/src/views/OrderView.vue` - 2 fixes (cart helpers usage, category tabs)

### Backend
1. `backend/application/services/user_management_service.go` - 1 fix (barista role)

---

## Testing Checklist

### Batch Unit Conversion
- [ ] Tạo batch "Sữa tươi" (1L)
- [ ] Thêm vào menu, chọn "ml", nhập 200ml
- [ ] Lưu món
- [ ] Click "Sửa" → Kiểm tra dropdown có ["L", "ml"], đang chọn "ml"

### Order Cart
- [ ] Thêm món single-size vào giỏ hàng
- [ ] Thêm món multi-size vào giỏ hàng
- [ ] Kiểm tra không có lỗi console

### User Barista
- [ ] Vào /users
- [ ] Tạo user mới với role "Barista"
- [ ] Kiểm tra user được tạo thành công

### Category Tabs
- [ ] Vào /orders
- [ ] Kiểm tra category tabs hiển thị đúng danh mục
- [ ] Click vào từng category → menu items được filter đúng

---

## Notes

- Tất cả fixes đều backward compatible
- Không cần migration database
- Frontend và backend có thể deploy độc lập
- Cần restart backend sau khi sửa user_management_service.go
