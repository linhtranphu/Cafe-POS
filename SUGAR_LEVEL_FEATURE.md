# Tính năng Chọn độ ngọt (Sugar Level)

## Tổng quan

Thêm tính năng cho phép waiter chọn độ ngọt cho từng món khi tạo order. Hệ thống hỗ trợ 3 mức:
- **25%** - Ít đường
- **50%** - Vừa đường  
- **100%** - Đường bình thường (mặc định)

## Design Philosophy

### Inline, không modal
- 3 nút nhỏ hiển thị ngay dưới quantity controls
- Không cần mở modal riêng
- Tap trực tiếp để chọn
- Nhanh, tiện lợi, ít thao tác

### Visual Design
- **Icon**: 🧊 (ice cube) để biểu thị độ ngọt
- **Màu sắc**: 
  - Selected: Vàng (`bg-amber-500`)
  - Unselected: Xám (`bg-gray-100`)
- **Layout**: 3 nút ngang, kích thước bằng nhau
- **Default**: 100% (không cần chọn nếu muốn bình thường)

## UI/UX

### Single-size items
```
┌─────────────────────┐
│  Cà phê sữa         │
│  25,000đ            │
│                     │
│  [−]  [2]  [+]     │  ← Quantity controls
│  🧊 [25%][50%][100%]│  ← Sugar level (100% highlighted)
└─────────────────────┘
```

### Multi-size items (variants)
```
┌─────────────────────┐
│  Trà sữa            │
│                     │
│  Size M - 30,000đ   │
│  [−] [1] [+]        │  ← Quantity controls
│  🧊 [25%][50%][100%] │  ← Sugar level for Size M
│                     │
│  Size L - 35,000đ   │
│  [−] [2] [+]        │  ← Quantity controls
│  🧊 [25%][50%][100%] │  ← Sugar level for Size L (independent)
└─────────────────────┘
```

## Luồng sử dụng

### 1. Thêm món vào cart
- Waiter tap "+" để thêm món
- Nút +/- hiện ra
- Sugar level selector hiện ra (default 100%)

### 2. Chọn độ ngọt (optional)
- Nếu muốn 100% → Không cần làm gì
- Nếu muốn 25% → Tap nút "25%"
- Nếu muốn 50% → Tap nút "50%"
- Nút được chọn sẽ highlight màu vàng

### 3. Thay đổi độ ngọt
- Có thể tap lại để đổi mức khác
- Ví dụ: 25% → 50% → 100%

### 4. Xác nhận order
- Tap "Xác nhận" ở floating button
- Nếu sugar level = 100% → Không có note
- Nếu sugar level = 25% → Note: "25% đường"
- Nếu sugar level = 50% → Note: "50% đường"

## Technical Implementation

### State Management

```javascript
// Sugar levels configuration
const sugarLevels = [
  { value: 25, label: '25%' },
  { value: 50, label: '50%' },
  { value: 100, label: '100%' }
]

// State: { itemId: level, itemId_variantId: level }
const itemSugarLevels = ref({})
```

### Methods

```javascript
// Get sugar level (default 100%)
getSugarLevel(itemId, variantId = null)

// Set sugar level
setSugarLevel(itemId, level, variantId = null)
```

### Data Flow

**1. Khi thêm món:**
```javascript
addItem("abc123", null)
// → cart.value["abc123"] = 1
// → itemSugarLevels.value["abc123"] = 100 (default)
```

**2. Khi chọn độ ngọt:**
```javascript
setSugarLevel("abc123", 25)
// → itemSugarLevels.value["abc123"] = 25
```

**3. Khi confirm order:**
```javascript
handleConfirm()
// → Tạo cartArray với note:
[
  {
    menu_item_id: "abc123",
    name: "Cà phê sữa",
    price: 25000,
    quantity: 2,
    note: "25% đường" // ← Auto-generated from sugar level
  }
]
```

**4. Backend nhận:**
```json
{
  "items": [
    {
      "menu_item_id": "abc123",
      "name": "Cà phê sữa",
      "price": 25000,
      "quantity": 2,
      "note": "25% đường"
    }
  ]
}
```

## Variants Support

Mỗi variant có sugar level riêng:

```javascript
// Trà sữa - Size M: 25% đường
itemSugarLevels.value["abc123_variant-m"] = 25

// Trà sữa - Size L: 50% đường
itemSugarLevels.value["abc123_variant-l"] = 50
```

## Note Generation Logic

```javascript
const sugarLevel = itemSugarLevels.value[key] || 100
const note = sugarLevel === 100 ? '' : `${sugarLevel}% đường`
```

- **100%**: Không có note (bình thường)
- **50%**: Note = "50% đường"
- **25%**: Note = "25% đường"

## Advantages

### 1. Nhanh chóng
- Không cần mở modal
- Chỉ 1 tap để chọn
- Inline, trực quan

### 2. Trực quan
- Thấy ngay 3 lựa chọn
- Highlight rõ ràng
- Icon 🧊 dễ nhận biết

### 3. Linh hoạt
- Mỗi món có độ ngọt riêng
- Mỗi variant có độ ngọt riêng
- Có thể đổi bất cứ lúc nào

### 4. Default thông minh
- 100% là default
- Không cần chọn nếu muốn bình thường
- Giảm thao tác cho case phổ biến nhất

## Testing Checklist

- [ ] Thêm món → Sugar level selector hiện ra
- [ ] Default là 100% (highlighted)
- [ ] Tap 25% → Nút 25% highlight
- [ ] Tap 50% → Nút 50% highlight
- [ ] Tap 100% → Nút 100% highlight
- [ ] Đổi qua lại giữa các mức → Hoạt động mượt
- [ ] Xóa món → Sugar level cũng bị xóa
- [ ] Confirm với 100% → Không có note
- [ ] Confirm với 50% → Note = "50% đường"
- [ ] Confirm với 25% → Note = "25% đường"
- [ ] Variants có sugar level độc lập
- [ ] Note hiển thị đúng trong order detail
- [ ] Note được in ra trong bill

## Responsive Design

### Mobile (< 640px)
- 3 nút ngang, kích thước bằng nhau
- Text size: `text-xs` (12px)
- Padding: `py-1.5` (6px vertical)
- Touch target: Đủ lớn cho ngón tay

### Tablet/Desktop
- Tương tự mobile
- Component được thiết kế cho mobile-first

## Color Scheme

```css
/* Selected */
bg-amber-500    /* #f59e0b - Vàng đậm */
text-white      /* Trắng */
shadow-sm       /* Bóng nhẹ */

/* Unselected */
bg-gray-100     /* #f3f4f6 - Xám nhạt */
text-gray-600   /* #4b5563 - Xám đậm */

/* Active (tap) */
active:scale-95 /* Scale down khi tap */
active:bg-gray-200 /* Xám đậm hơn khi tap */
```

## Future Enhancements

### 1. Thêm mức độ ngọt
```javascript
const sugarLevels = [
  { value: 0, label: '0%' },    // Không đường
  { value: 25, label: '25%' },
  { value: 50, label: '50%' },
  { value: 75, label: '75%' },
  { value: 100, label: '100%' },
  { value: 120, label: '120%' }  // Extra sweet
]
```

### 2. Lưu preference
- Nhớ độ ngọt thường dùng của từng món
- Auto-select based on history

### 3. Customize per item
- Một số món không cần chọn độ ngọt (bánh, snack)
- Chỉ hiện cho đồ uống

### 4. Ice level
- Thêm selector cho lượng đá
- Tương tự sugar level

## Files thay đổi

1. `frontend/src/components/CreateOrderModal.vue` - Thêm UI và logic

## Backend

Backend đã hỗ trợ sẵn field `note` trong `OrderItem`. Không cần thay đổi.

---

**Ngày thực hiện:** 4 tháng 3, 2026
**Trạng thái:** ✅ Hoàn thành
