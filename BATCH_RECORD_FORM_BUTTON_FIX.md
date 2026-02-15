# Batch Record Form - Button Visibility Fix

## Vấn Đề
Buttons "Hủy" và "Ghi Nhận" trong form ghi nhận batch bị che khuất bởi BottomNav.

## Nguyên Nhân
- Fixed Footer có `pb-safe` class nhưng không đủ padding để tránh BottomNav
- BottomNav được render sau Fixed Footer và có thể che khuất buttons
- Không có z-index để đảm bảo buttons luôn hiển thị trên cùng

## Giải Pháp

### Cập Nhật Fixed Footer
Thay đổi từ:
```html
<div class="flex-shrink-0 bg-white px-4 py-4 border-t flex gap-3 pb-safe">
```

Thành:
```html
<div class="flex-shrink-0 bg-white px-4 py-4 border-t flex gap-3 z-50" 
     style="padding-bottom: max(1rem, calc(env(safe-area-inset-bottom) + 5rem))">
```

### Thay Đổi Chi Tiết

1. **Thêm z-index**: `z-50` để đảm bảo buttons luôn hiển thị trên BottomNav
2. **Cải thiện padding**: 
   - Sử dụng inline style với `max()` function
   - Tính toán: `env(safe-area-inset-bottom) + 5rem`
   - 5rem = khoảng 80px, đủ để tránh BottomNav (thường cao ~60-70px)
   - Fallback: `1rem` nếu không có safe-area-inset

## Lợi Ích

1. **Buttons luôn hiển thị**: z-50 đảm bảo buttons không bị che
2. **Responsive với notch**: `env(safe-area-inset-bottom)` xử lý iPhone notch
3. **Đủ khoảng cách**: 5rem padding đảm bảo không bị che bởi BottomNav
4. **Tương thích**: Fallback 1rem cho browsers không hỗ trợ safe-area

## Visual Layout

```
┌─────────────────────────┐
│   Header (Fixed)        │
├─────────────────────────┤
│                         │
│   Scrollable Content    │
│                         │
│   - Batch Selector      │
│   - Batch Counter       │
│   - Ingredients         │
│   - Cost Summary        │
│                         │
├─────────────────────────┤
│   Fixed Footer (z-50)   │ ← Buttons ở đây
│   [Hủy]  [Ghi Nhận]     │
│   padding-bottom: 5rem  │ ← Đủ khoảng cách
├─────────────────────────┤
│   BottomNav (z-40)      │ ← Không che buttons
└─────────────────────────┘
```

## Testing Checklist
- [ ] Buttons "Hủy" và "Ghi Nhận" hiển thị đầy đủ
- [ ] Buttons không bị che bởi BottomNav
- [ ] Buttons có thể click được
- [ ] Padding đủ trên iPhone có notch
- [ ] Padding đủ trên Android
- [ ] Buttons disabled khi form chưa valid

## Files Modified
- `frontend/src/components/batch/BatchRecordForm.vue`

## Status
✅ Fix hoàn tất
✅ Không có lỗi diagnostic
🧪 Sẵn sàng test
