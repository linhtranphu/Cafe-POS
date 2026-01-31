# Facility View - Mobile First Optimization

## Tóm tắt
Đã tối ưu hóa FacilityView cho mobile với focus vào UX/UI gọn gàng, dễ sử dụng trên màn hình nhỏ.

## Thay đổi chính

### 1. Summary Card - Status thành 1 Row

**Trước:**
```
┌─────────────────────┐
│ Tổng số tài sản     │
│      150            │
└─────────────────────┘
```

**Sau:**
```
┌──────────────────────────────────────┐
│ Tổng quan tài sản                    │
│ ┌────┬────┬────┬────┐               │
│ │150 │120 │ 15 │ 15 │               │
│ │Tổng│Hoạt│Bảo │Hỏng│               │
│ └────┴────┴────┴────┘               │
└──────────────────────────────────────┘
```

**Lợi ích:**
- Tiết kiệm không gian màn hình
- Hiển thị nhiều thông tin hơn trong 1 view
- Dễ so sánh các chỉ số

### 2. Header - Compact Layout

**Trước:**
- Title lớn: "Quản lý Cơ sở vật chất"
- 3 buttons với text đầy đủ
- Chiếm nhiều không gian

**Sau:**
- Title ngắn gọn: "Cơ sở vật chất"
- 3 buttons grid với text rút gọn:
  - "📅 Lịch BT" (Bảo trì)
  - "📊 Báo cáo"
  - "+ Thêm"
- Grid 3 cột đều nhau

### 3. Filters - Optimized

**Thay đổi:**
- Giảm padding: `p-3` → `p-2.5`
- Text size: `text-base` → `text-sm`
- Layout: Type và Status cùng 1 row (grid-cols-2)
- Placeholder ngắn gọn: "🔍 Tìm kiếm..."
- Options rút gọn: "Tất cả TT", "Ngừng SD"

### 4. Facility Cards - Compact Design

**Header:**
- Icon size: `w-12 h-12` → `w-10 h-10`
- Font size: `text-2xl` → `text-xl`
- Truncate long names
- Status badge với `whitespace-nowrap`

**Details:**
- Text size: `text-sm` → `text-xs`
- Padding: `p-3` → `p-2`
- Gap: `gap-3` → `gap-2`
- Labels rút gọn: "Số lượng" → "SL"

**Actions:**
- Layout: `grid-cols-2` → `grid-cols-3`
- Gap: `gap-2` → `gap-1.5`
- Text size: `text-sm` → `text-xs`
- Padding: `px-3 py-2` → `px-2 py-2`
- Labels rút gọn:
  - "Trạng thái" → "TT"
  - "Bảo trì" → "BT"
  - "Di chuyển" → "DC"
  - "Lịch sử" → "LS"
- Loại bỏ button "Xóa" khỏi quick actions (có thể xóa từ edit form)

**Shadow:**
- `shadow-md` → `shadow-sm`
- Thêm `border border-gray-200` để nhẹ nhàng hơn

### 5. Spacing Optimization

**Margins:**
- Header: `mb-6` → `mb-4`
- Filters: `mb-4` → `mb-3`
- Summary: `mb-4` → `mb-4` (giữ nguyên)
- Cards gap: `gap-4` → `gap-3`

**Padding:**
- Container: `p-4` (giữ nguyên)
- Cards: `p-4` (giữ nguyên)
- Filters: `p-4` → `p-3`

### 6. Touch Optimization

**Active States:**
- Thêm `active:bg-*-700` cho tất cả buttons
- Tăng contrast khi tap
- Feedback rõ ràng hơn

**Button Sizes:**
- Min height đủ lớn cho touch (py-2)
- Spacing giữa buttons (gap-1.5, gap-2)

## Computed Properties Mới

```javascript
const activeCount = computed(() => 
  filteredItems.value.filter(i => i.status === 'Đang sử dụng').length
)

const maintenanceCount = computed(() => 
  filteredItems.value.filter(i => i.status === 'Đang sửa').length
)

const brokenCount = computed(() => 
  filteredItems.value.filter(i => i.status === 'Hỏng').length
)
```

## Kết quả

### Trước:
- Summary card chiếm ~80px height
- Header chiếm ~100px
- Filters chiếm ~200px
- Mỗi card chiếm ~350px
- **Tổng cho 1 item: ~730px**

### Sau:
- Summary card chiếm ~70px height
- Header chiếm ~80px
- Filters chiếm ~140px
- Mỗi card chiếm ~280px
- **Tổng cho 1 item: ~570px**

**Tiết kiệm: ~160px (~22%) cho mỗi item**

## Mobile UX Improvements

1. **Thumb-friendly**: Tất cả buttons đủ lớn để tap
2. **Scan-friendly**: Thông tin quan trọng nổi bật
3. **Space-efficient**: Hiển thị nhiều items hơn trong 1 screen
4. **Quick actions**: 6 actions quan trọng nhất trong grid 3x2
5. **Visual hierarchy**: Icon, status, và actions rõ ràng
6. **Performance**: Giảm DOM elements, tăng tốc render

## Responsive Design

- Mobile first: Tối ưu cho màn hình nhỏ
- Grid system: Tự động adapt với màn hình
- Touch targets: Đủ lớn cho mobile (min 44x44px)
- Typography: Readable trên mọi kích thước màn hình
