# Fix: Cannot read properties of undefined (reading 'createCartItem')

## Vấn đề

Khi tạo order (thêm món vào giỏ hàng), xuất hiện lỗi:

```
TypeError: Cannot read properties of undefined (reading 'createCartItem')
at addToCart (OrderView.vue)
```

## Nguyên nhân

Trong file `frontend/src/stores/order.js`, `helpers` object được định nghĩa như một property của store definition:

```javascript
export const useOrderStore = defineStore('order', {
  state: () => ({ ... }),
  actions: { ... },
  getters: { ... },
  helpers: {  // ❌ Không hợp lệ trong Pinia
    createCartItem() { ... },
    isSameCartItem() { ... }
  }
})
```

Trong Pinia, chỉ có 3 properties hợp lệ cho store definition:
- `state`
- `actions`
- `getters`

Property `helpers` không được Pinia nhận diện, nên `orderStore.helpers` là `undefined`.

## Giải pháp

Export `helpers` như một object riêng biệt, không phải là part của store definition.

### 1. Sửa `frontend/src/stores/order.js`

**Trước:**
```javascript
export const useOrderStore = defineStore('order', {
  // ...
  helpers: {
    createCartItem() { ... },
    isSameCartItem() { ... }
  }
})
```

**Sau:**
```javascript
export const useOrderStore = defineStore('order', {
  state: () => ({ ... }),
  actions: { ... },
  getters: { ... }
})

// Export helpers separately
export const cartHelpers = {
  createCartItem(menuItem, variant = null) {
    const item = {
      menu_item_id: menuItem.id,
      name: menuItem.name,
      quantity: 1
    }

    if (variant) {
      // Multi-size item with variant
      item.variant_id = variant.id
      item.variant_name = variant.name
      item.price = variant.price
    } else {
      // Single-size item (backward compatible)
      item.price = menuItem.price
    }

    return item
  },

  isSameCartItem(item1, item2) {
    if (item1.menu_item_id !== item2.menu_item_id) {
      return false
    }
    // For multi-size items, variant_id must match
    if (item1.variant_id || item2.variant_id) {
      return item1.variant_id === item2.variant_id
    }
    // For single-size items, just menu_item_id match is enough
    return true
  }
}
```

### 2. Sửa `frontend/src/views/OrderView.vue`

**Import cartHelpers:**
```javascript
import { useOrderStore, cartHelpers } from '../stores/order'
```

**Sử dụng cartHelpers:**
```javascript
const addToCart = (item, variant = null) => {
  // Use cartHelpers instead of orderStore.helpers
  const cartItem = cartHelpers.createCartItem(item, variant)
  const existing = cart.value.find(i => cartHelpers.isSameCartItem(i, cartItem))
  
  if (existing) {
    existing.quantity++
  } else {
    cart.value.push(cartItem)
  }
}
```

## Tại sao export riêng?

### Lý do 1: Pinia store structure
Pinia chỉ hỗ trợ 3 properties: `state`, `actions`, `getters`. Bất kỳ property nào khác sẽ bị ignore.

### Lý do 2: Helpers không cần store state
Các helper functions này là pure functions, không cần access vào store state. Export riêng giúp code rõ ràng hơn.

### Lý do 3: Reusability
Export riêng cho phép import và sử dụng helpers ở bất kỳ đâu mà không cần khởi tạo store.

## Testing

### Test case 1: Thêm single-size item vào cart
1. Mở OrderView
2. Click vào món không có variants (ví dụ: "Cà phê đen")
3. ✅ Món được thêm vào giỏ hàng
4. ✅ Không có lỗi console

### Test case 2: Thêm multi-size item vào cart
1. Mở OrderView
2. Click vào món có variants (ví dụ: "Cà phê sữa")
3. Chọn size (ví dụ: "Size M")
4. ✅ Món với size được thêm vào giỏ hàng
5. ✅ Hiển thị đúng tên và giá của variant

### Test case 3: Thêm cùng món nhiều lần
1. Click vào món "Cà phê đen" 3 lần
2. ✅ Quantity tăng lên 3 (không tạo 3 items riêng)

### Test case 4: Thêm cùng món nhưng khác size
1. Click "Cà phê sữa" → Size M
2. Click "Cà phê sữa" → Size L
3. ✅ Có 2 items riêng trong cart (vì khác variant_id)

## Kết quả

- ✅ Lỗi `Cannot read properties of undefined` đã được fix
- ✅ Có thể thêm món vào giỏ hàng
- ✅ Single-size items hoạt động đúng
- ✅ Multi-size items (variants) hoạt động đúng
- ✅ Logic kiểm tra duplicate items hoạt động đúng

## Code structure sau khi fix

```
frontend/src/stores/order.js
├── useOrderStore (Pinia store)
│   ├── state
│   ├── actions
│   └── getters
└── cartHelpers (exported object)
    ├── createCartItem()
    └── isSameCartItem()

frontend/src/views/OrderView.vue
├── import { useOrderStore, cartHelpers }
└── addToCart() uses cartHelpers
```

## Notes

- Helpers không cần access store state nên export riêng là best practice
- Nếu cần access store state, nên move vào `actions` thay vì `helpers`
- Pattern này có thể áp dụng cho các stores khác nếu cần
