# Orders View - Shift Tabs Implementation

## Tổng quan
Đã thiết kế lại Orders view với hệ thống tab để phân tách orders của ca hiện tại và ca khác, giúp waiter dễ dàng tập trung vào orders của ca đang làm việc.

## Thay đổi chính

### 1. Thêm Shift Filter Tabs
- **Tab "Ca hiện tại" (🔵)**: Hiển thị orders của ca đang mở
- **Tab "Ca khác" (📂)**: Hiển thị orders từ các ca khác
- Mỗi tab hiển thị số lượng orders trong ngoặc
- Tab mặc định: "Ca hiện tại"

### 2. Logic Filter
```javascript
// Filter theo shift trước
const ordersByShift = computed(() => {
  const currentShiftId = shiftStore.currentShift?.id
  
  if (shiftFilter.value === 'current') {
    // Chỉ hiển thị orders của ca hiện tại
    if (!currentShiftId) return []
    return orders.value.filter(o => o.shift_id === currentShiftId)
  } else {
    // Hiển thị orders từ ca khác (không phải ca hiện tại)
    if (!currentShiftId) return orders.value
    return orders.value.filter(o => o.shift_id !== currentShiftId)
  }
})

// Sau đó filter theo status
const filteredOrders = computed(() => {
  if (filterStatus.value === 'ALL') return ordersByShift.value
  return ordersByShift.value.filter(o => o.status === filterStatus.value)
})
```

### 3. UI/UX Improvements

#### Shift Tabs
- Nằm giữa header và status filter pills
- Full width, 2 buttons cân đối
- Active state: blue background với shadow
- Inactive state: gray background
- Badge count với opacity 75%

#### Status Filter Pills
- Đổi màu active từ blue → green để phân biệt với shift tabs
- Vẫn giữ horizontal scroll cho mobile
- Count badge cập nhật theo shift filter hiện tại

#### Order Cards
- Khi xem "Ca khác", hiển thị thêm thông tin shift ID:
  ```
  📋 Ca: 699bda2e...
  ```
- Màu purple để phân biệt với các thông tin khác

#### Empty State
- Contextual message:
  - Ca hiện tại: "Chưa có order nào trong ca này"
  - Ca khác: "Không có order từ ca khác"

### 4. Badge Counts
```javascript
// Đếm orders cho mỗi tab
const currentShiftOrdersCount = computed(() => {
  const currentShiftId = shiftStore.currentShift?.id
  if (!currentShiftId) return 0
  return orders.value.filter(o => o.shift_id === currentShiftId).length
})

const otherShiftsOrdersCount = computed(() => {
  const currentShiftId = shiftStore.currentShift?.id
  if (!currentShiftId) return orders.value.length
  return orders.value.filter(o => o.shift_id !== currentShiftId).length
})
```

## Workflow

### Khi waiter mở Orders view:
1. Mặc định hiển thị tab "Ca hiện tại"
2. Chỉ thấy orders của ca đang làm việc
3. Có thể filter thêm theo status (ALL, CREATED, PAID, etc.)

### Khi cần xem orders cũ:
1. Click tab "Ca khác"
2. Xem tất cả orders từ các ca trước
3. Mỗi order hiển thị shift ID để phân biệt
4. Vẫn có thể filter theo status

### Khi chưa mở ca:
- Tab "Ca hiện tại" sẽ trống (0 orders)
- Tab "Ca khác" hiển thị tất cả orders có sẵn
- Vẫn hiển thị warning "Chưa mở ca làm việc"

## Lợi ích

1. **Tập trung cao hơn**: Waiter chỉ thấy orders của ca hiện tại
2. **Giảm nhiễu**: Không bị phân tâm bởi orders cũ
3. **Linh hoạt**: Vẫn có thể xem lại orders từ ca khác khi cần
4. **Trực quan**: Badge count giúp biết ngay có bao nhiêu orders
5. **Mobile-friendly**: Tab design phù hợp với màn hình nhỏ

## Files Changed
- `frontend/src/views/OrderView.vue`

## Testing Checklist
- [ ] Tab "Ca hiện tại" chỉ hiển thị orders của ca đang mở
- [ ] Tab "Ca khác" hiển thị orders từ các ca khác
- [ ] Badge count chính xác cho cả 2 tabs
- [ ] Status filter hoạt động đúng trong cả 2 tabs
- [ ] Shift ID hiển thị khi xem "Ca khác"
- [ ] Empty state message phù hợp với tab đang xem
- [ ] Khi chưa mở ca, logic vẫn hoạt động đúng
- [ ] UI responsive trên mobile
