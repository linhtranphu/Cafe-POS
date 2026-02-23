# Shift View - Tabs Implementation (Always Visible)

## Tổng quan
Đã thiết kế lại Shift view với hệ thống tab luôn hiển thị để tách ca hiện tại và ca cũ. Tabs không phụ thuộc vào việc có ca đang mở hay không.

## Thay đổi chính

### 1. Tabs Luôn Hiển Thị
- **Tab "Ca hiện tại" (🔵)**: 
  - Khi chưa mở ca: hiển thị form "Mở ca làm việc"
  - Khi đã mở ca: hiển thị thông tin ca đang mở + lịch sử bàn giao
- **Tab "Ca cũ" (📂)**: 
  - Hiển thị lịch sử các ca đã đóng
  - Badge count hiển thị số lượng ca cũ
- Tabs luôn hiển thị cho waiter/barista (không phụ thuộc vào currentShift)

### 2. Logic Filter
```javascript
// State
const shiftFilter = ref('current') // 'current' or 'history'

// Filter shifts - exclude current shift from history
const filteredShifts = computed(() => {
  // When viewing history, show all shifts except current one
  if (shiftFilter.value === 'history') {
    if (!currentShift.value) {
      // No current shift, show all shifts
      return shifts.value
    }
    // Exclude the current shift from history
    return shifts.value.filter(s => s.id !== currentShift.value.id)
  }
  
  // When viewing current tab, don't show any shifts in the list
  // (current shift is shown separately above)
  return []
})

// Count old shifts (excluding current)
const oldShiftsCount = computed(() => {
  if (!currentShift.value) return shifts.value.length
  return shifts.value.filter(s => s.id !== currentShift.value.id).length
})
```

### 3. Cấu trúc Tab

#### Tab "Ca hiện tại" (shiftFilter === 'current')
```vue
<div v-if="shiftFilter === 'current'">
  <!-- Current Shift Info (when shift is open) -->
  <div v-if="currentShift" class="...">
    <!-- Thông tin ca đang mở -->
    <!-- Lịch sử bàn giao -->
    <!-- Các nút action -->
  </div>
  
  <!-- Start Shift Form (when no current shift) -->
  <div v-if="!currentShift" class="...">
    <h3>Mở ca làm việc</h3>
    <!-- Form mở ca -->
  </div>
</div>
```

#### Tab "Ca cũ" (shiftFilter === 'history')
```vue
<div v-if="shiftFilter === 'history'" class="...">
  <h3>Lịch sử ca làm việc</h3>
  <!-- Danh sách các ca đã đóng -->
  <!-- Không bao gồm ca hiện tại -->
</div>
```

### 4. UI/UX Improvements

#### Tabs Design
- Luôn hiển thị cho waiter/barista
- Full width, 2 buttons cân đối
- Active state: blue background với shadow
- Inactive state: gray background
- Badge count cho "Ca cũ" tab

#### Empty State Messages
- Tab "Ca cũ": "Chưa có ca làm việc cũ"
- Không còn contextual message phức tạp

#### Section Headers
- Tab "Ca cũ": luôn là "Lịch sử ca làm việc"

## Workflow

### Khi chưa mở ca:
1. Tabs vẫn hiển thị (mặc định "Ca hiện tại")
2. Tab "Ca hiện tại":
   - Hiển thị form "Mở ca làm việc"
   - Không hiển thị thông tin ca nào
3. Tab "Ca cũ":
   - Hiển thị tất cả shifts trong lịch sử (nếu có)

### Khi đã mở ca:
1. Tabs vẫn hiển thị (mặc định "Ca hiện tại")
2. Tab "Ca hiện tại":
   - Hiển thị thông tin ca đang mở
   - Hiển thị lịch sử bàn giao (nếu là waiter)
   - Các nút action (bàn giao, đóng ca)
   - KHÔNG hiển thị form "Mở ca làm việc"
3. Tab "Ca cũ":
   - Hiển thị danh sách các ca đã đóng
   - Không bao gồm ca hiện tại
   - Badge count hiển thị số lượng

### Đặc biệt cho Cashier:
- Cashier có component riêng (`CashierShiftView`)
- Không hiển thị tabs
- Vẫn có thể chốt ca của waiter/barista

## Lợi ích

1. **Nhất quán**: Tabs luôn hiển thị, không thay đổi layout khi mở/đóng ca
2. **Trực quan**: Dễ hiểu - "Ca hiện tại" là nơi làm việc với ca hiện tại hoặc mở ca mới
3. **Tách bạch rõ ràng**: Ca hiện tại và ca cũ hoàn toàn tách biệt
4. **UX tốt hơn**: Không có sự thay đổi đột ngột trong UI
5. **Badge count**: Luôn biết có bao nhiêu ca cũ
6. **Mobile-friendly**: Tab design phù hợp với màn hình nhỏ

## So sánh với phiên bản trước

### Trước (tabs chỉ hiển thị khi có ca):
- ❌ Layout thay đổi khi mở ca (tabs xuất hiện đột ngột)
- ❌ Khi chưa mở ca, form "Mở ca" và lịch sử ca trộn lẫn
- ❌ Không nhất quán

### Sau (tabs luôn hiển thị):
- ✅ Layout nhất quán, tabs luôn ở đó
- ✅ Tách bạch rõ ràng: "Ca hiện tại" vs "Ca cũ"
- ✅ Form "Mở ca" nằm trong tab "Ca hiện tại" - logic hơn
- ✅ Không có sự thay đổi đột ngột trong UI

## Files Changed
- `frontend/src/views/ShiftView.vue`

## Testing Checklist
- [ ] Tabs luôn hiển thị cho waiter/barista (cả khi chưa mở ca)
- [ ] Tab "Ca hiện tại" - chưa mở ca: hiển thị form "Mở ca làm việc"
- [ ] Tab "Ca hiện tại" - đã mở ca: hiển thị thông tin ca + lịch sử bàn giao
- [ ] Tab "Ca cũ" - chưa mở ca: hiển thị tất cả shifts
- [ ] Tab "Ca cũ" - đã mở ca: hiển thị các ca cũ (không bao gồm ca hiện tại)
- [ ] Badge count chính xác
- [ ] Empty state message: "Chưa có ca làm việc cũ"
- [ ] Cashier view không bị ảnh hưởng (không hiển thị tabs)
- [ ] UI responsive trên mobile
- [ ] Pull to refresh hoạt động đúng
- [ ] Chuyển tab mượt mà, không lag
