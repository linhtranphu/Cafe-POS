# Batch Record Form - Layout Alignment with Facility Form

## Vấn Đề
BatchRecordForm có BottomNav gây xung đột với Fixed Footer buttons, làm buttons bị che khuất.

## Giải Pháp
Điều chỉnh layout của BatchRecordForm để giống với FacilityManagementView:
- Loại bỏ BottomNav khỏi form
- Sử dụng `pb-safe` thay vì custom padding
- Giữ Fixed Footer đơn giản và clean

## Thay Đổi

### 1. Loại Bỏ BottomNav Component
**Trước:**
```vue
<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Content -->
    
    <!-- Fixed Footer -->
    <div class="flex-shrink-0 bg-white px-4 py-4 border-t flex gap-3 z-50" 
         style="padding-bottom: max(1rem, calc(env(safe-area-inset-bottom) + 5rem))">
      <!-- Buttons -->
    </div>
    
    <!-- Bottom Navigation -->
    <BottomNav />
  </div>
</template>

<script setup>
import BottomNav from '../BottomNav.vue'
</script>
```

**Sau:**
```vue
<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Content -->
    
    <!-- Fixed Footer -->
    <div class="flex-shrink-0 bg-white px-4 py-4 border-t flex gap-3 pb-safe">
      <!-- Buttons -->
    </div>
  </div>
</template>

<script setup>
// BottomNav removed
</script>
```

### 2. Đơn Giản Hóa Fixed Footer
- Xóa `z-50` (không cần vì không có BottomNav)
- Xóa custom `style="padding-bottom: ..."` 
- Sử dụng `pb-safe` class (Tailwind utility cho safe-area-inset-bottom)

## So Sánh Layout

### FacilityManagementView (Reference)
```
┌─────────────────────────┐
│   Fixed Header          │
├─────────────────────────┤
│                         │
│   Scrollable Content    │
│                         │
├─────────────────────────┤
│   Fixed Footer          │
│   [Hủy]  [Thêm mới]     │
│   pb-safe               │
└─────────────────────────┘
```

### BatchRecordForm (Sau khi fix)
```
┌─────────────────────────┐
│   Fixed Header          │
├─────────────────────────┤
│                         │
│   Scrollable Content    │
│   - Batch Selector      │
│   - Batch Counter       │
│   - Ingredients         │
│   - Cost Summary        │
│                         │
├─────────────────────────┤
│   Fixed Footer          │
│   [Hủy]  [Ghi Nhận]     │
│   pb-safe               │
└─────────────────────────┘
```

## Lý Do Loại Bỏ BottomNav

1. **Xung đột UI**: BottomNav che khuất Fixed Footer buttons
2. **Không nhất quán**: Các form khác (Facility, Ingredient) không có BottomNav
3. **UX tốt hơn**: Form có buttons riêng, không cần navigation
4. **Đơn giản hơn**: Ít component, ít complexity

## Pattern Chung Cho Forms

Tất cả forms tạo/sửa nên follow pattern này:
```vue
<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Fixed Header -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3">
        <div class="flex items-center justify-between">
          <button @click="goBack">←</button>
          <h1>Title</h1>
          <div class="w-8"></div>
        </div>
      </div>
    </div>

    <!-- Scrollable Content -->
    <div class="flex-1 overflow-y-auto px-4 py-6 space-y-5">
      <!-- Form fields -->
      
      <!-- Spacer for bottom buttons -->
      <div class="h-24"></div>
    </div>

    <!-- Fixed Footer -->
    <div class="flex-shrink-0 bg-white px-4 py-4 border-t flex gap-3 pb-safe">
      <button @click="goBack">Hủy</button>
      <button @click="submit">Lưu</button>
    </div>
  </div>
</template>
```

## Lợi Ích

1. **Nhất quán**: Tất cả forms có cùng layout pattern
2. **Đơn giản**: Không cần quản lý z-index phức tạp
3. **Responsive**: `pb-safe` tự động xử lý safe-area
4. **Clean**: Không có component thừa

## Files Modified
- `frontend/src/components/batch/BatchRecordForm.vue`

## Testing Checklist
- [ ] Buttons "Hủy" và "Ghi Nhận" hiển thị đầy đủ
- [ ] Buttons không bị che khuất
- [ ] Buttons có thể click được
- [ ] Padding đủ trên iPhone có notch
- [ ] Padding đủ trên Android
- [ ] Layout giống FacilityManagementView

## Status
✅ Fix hoàn tất
✅ Không có lỗi diagnostic
✅ Layout nhất quán với Facility form
🧪 Sẵn sàng test
