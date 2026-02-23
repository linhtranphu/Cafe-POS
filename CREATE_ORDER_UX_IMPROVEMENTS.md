# Create Order UX Improvements

## Tổng quan
Đã cải thiện UI/UX cho màn hình tạo order, tập trung vào 2 vấn đề chính:
1. Làm nổi bật phần tên khách hàng để tránh bị bỏ qua
2. Thiết kế lại menu items cho đồng đều và gọn gàng hơn

## Vấn đề trước đây

### 1. Tên khách hàng dễ bị missed
- Input field nhỏ, màu xám nhạt
- Không có label rõ ràng
- Nằm trong vùng bg-gray-50 không nổi bật
- Placeholder text không đủ thu hút sự chú ý

### 2. Menu items không đồng đều
- Single-size items và multi-size items có chiều cao khác nhau
- Layout không consistent
- Badge "X món" có style khác nhau
- Không có min-height cho tên món

## Giải pháp

### 1. Customer Name Section - Prominent Design

#### Before:
```vue
<div class="px-4 py-3 bg-gray-50 border-b">
  <input v-model="customerName" 
    type="text" 
    placeholder="Tên khách hàng (tùy chọn)"
    class="w-full px-4 py-3 rounded-lg border">
</div>
```

#### After:
```vue
<div class="px-4 py-4 bg-gradient-to-b from-blue-50 to-white border-b-2 border-blue-100">
  <label class="block text-sm font-bold text-gray-700 mb-2 flex items-center gap-2">
    <span class="text-lg">👤</span>
    <span>Tên khách hàng</span>
    <span class="text-xs font-normal text-gray-500">(tùy chọn)</span>
  </label>
  <input v-model="customerName" 
    type="text" 
    placeholder="Nhập tên khách hàng..."
    class="w-full px-4 py-3.5 rounded-xl border-2 border-blue-200 bg-white 
           focus:ring-2 focus:ring-blue-500 focus:border-blue-500 
           text-base font-medium placeholder-gray-400">
</div>
```

**Improvements:**
- ✅ Gradient background (blue-50 → white) để nổi bật
- ✅ Label rõ ràng với icon 👤
- ✅ Border-2 thay vì border-1 (dễ nhìn hơn)
- ✅ Border color blue-200 (nổi bật hơn gray)
- ✅ Padding tăng (py-3 → py-3.5)
- ✅ Border radius lớn hơn (lg → xl)
- ✅ Font-medium cho input text
- ✅ Placeholder rõ ràng hơn
- ✅ Focus state với ring-2

### 2. Menu Items - Uniform Design

#### Key Changes:

**A. Consistent Height với Flexbox:**
```vue
<!-- Single-size item -->
<button class="h-full flex flex-col">
  <div class="flex-1">
    <div class="min-h-[2.5rem] line-clamp-2">{{ item.name }}</div>
    <div>{{ formatPrice(item.price) }}</div>
  </div>
  <div v-if="hasInCart">Badge</div>
</button>

<!-- Multi-size item -->
<div class="flex flex-col h-full">
  <div class="min-h-[2.5rem] line-clamp-2">{{ item.name }}</div>
  <div class="flex-1">Variants</div>
  <div v-if="hasInCart">Badge</div>
</div>
```

**Benefits:**
- `h-full` - card chiều cao 100% của grid cell
- `flex flex-col` - layout dọc
- `flex-1` - phần giữa tự động fill space
- `min-h-[2.5rem]` - tên món có chiều cao tối thiểu
- `line-clamp-2` - giới hạn 2 dòng, tự động ...

**B. Line Clamp CSS:**
```css
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.25rem;
}
```

**Benefits:**
- Tên món dài tự động cắt sau 2 dòng
- Thêm "..." tự động
- Chiều cao consistent (2.5rem = 2 lines)

**C. Improved Card Styling:**
```vue
class="bg-white rounded-2xl p-4 shadow-sm border border-gray-100 
       active:scale-95 active:shadow-md transition-all"
```

**Improvements:**
- Border radius từ xl → 2xl (softer)
- Thêm border subtle
- Shadow tăng khi active
- Background gray-50 cho grid container

**D. Uniform Badge Design:**
```vue
<div class="mt-3 bg-blue-500 text-white text-xs font-semibold 
            px-3 py-1.5 rounded-full inline-flex items-center 
            justify-center gap-1 shadow-md">
  <span>🛒</span>
  <span>{{ qty }} món</span>
</div>
```

**Improvements:**
- Consistent spacing (mt-3)
- Icon 🛒 thay vì chỉ text
- inline-flex với gap-1
- Shadow-md để nổi bật
- Font-semibold

**E. Variant Buttons:**
```vue
<button class="w-full flex justify-between items-center 
               px-3 py-2.5 bg-gray-50 rounded-xl 
               active:bg-blue-50 active:scale-95 
               border border-gray-100">
  <span class="text-xs font-semibold">{{ variant.name }}</span>
  <span class="text-sm font-bold">{{ formatPrice(variant.price) }}</span>
</button>
```

**Improvements:**
- Padding tăng (py-2.5)
- Border radius lớn hơn (lg → xl)
- Thêm border
- Active state với bg-blue-50
- Font-semibold cho variant name

### 3. Header Improvements

```vue
<div class="bg-blue-500 text-white px-4 py-4 shadow-lg" 
     style="padding-top: max(1rem, env(safe-area-inset-top))">
  <button class="text-2xl p-2 active:bg-blue-600 rounded-lg">←</button>
  <h2 class="text-xl font-bold">Tạo Order Mới</h2>
  <button class="px-4 py-2.5 bg-white text-blue-500 rounded-xl 
                 font-semibold shadow-md active:scale-95">
    Xác nhận
  </button>
</div>
```

**Improvements:**
- Safe area inset support
- Shadow-lg cho depth
- Back button có padding và active state
- Title lớn hơn (lg → xl)
- Confirm button có shadow và animation

### 4. Category Tabs

```vue
<button :class="[
  'px-4 py-2.5 rounded-full text-sm font-semibold 
   touch-manipulation active:scale-95',
  selectedCategory === cat.id 
    ? 'bg-blue-500 text-white shadow-lg' 
    : 'bg-gray-100 text-gray-700 active:bg-gray-200'
]">
  <span class="mr-1">{{ cat.icon }}</span>
  <span>{{ cat.name }}</span>
</button>
```

**Improvements:**
- Padding tăng (py-2 → py-2.5)
- Font-semibold
- Shadow-lg khi active
- Active state cho inactive tabs
- Icon và text tách riêng với spacing

## Visual Comparison

### Customer Name Section

**Before:**
```
┌─────────────────────────────────┐
│ [Tên khách hàng (tùy chọn)   ] │ ← Nhỏ, xám, dễ bỏ qua
└─────────────────────────────────┘
```

**After:**
```
┌─────────────────────────────────┐
│ 👤 Tên khách hàng (tùy chọn)   │ ← Label rõ ràng
│ ┌─────────────────────────────┐ │
│ │ Nhập tên khách hàng...      │ │ ← Nổi bật, border xanh
│ └─────────────────────────────┘ │
└─────────────────────────────────┘
```

### Menu Items Grid

**Before:**
```
┌──────────┐ ┌──────────┐
│ Cà phê   │ │ Trà      │
│ 25,000đ  │ │ S 15,000 │ ← Chiều cao khác nhau
│          │ │ M 20,000 │
│ 2 món    │ │ L 25,000 │
└──────────┘ │ 3 món    │
             └──────────┘
```

**After:**
```
┌──────────┐ ┌──────────┐
│ Cà phê   │ │ Trà      │
│ 25,000đ  │ │ S 15,000 │
│          │ │ M 20,000 │ ← Chiều cao đồng đều
│          │ │ L 25,000 │
│ 🛒 2 món │ │ 🛒 3 món │
└──────────┘ └──────────┘
```

## Benefits

### 1. Customer Name
- ✅ Không thể bỏ qua - gradient background nổi bật
- ✅ Label rõ ràng với icon
- ✅ Visual hierarchy tốt hơn
- ✅ Focus state rõ ràng

### 2. Menu Items
- ✅ Chiều cao đồng đều - professional look
- ✅ Tên món không bị overflow
- ✅ Layout consistent
- ✅ Badge position consistent
- ✅ Dễ scan và chọn món

### 3. Overall UX
- ✅ Touch targets đủ lớn (44px+)
- ✅ Visual feedback rõ ràng
- ✅ Spacing consistent
- ✅ Mobile-optimized

## Technical Details

### Flexbox Layout
```vue
<!-- Parent -->
<div class="grid grid-cols-2 gap-3">
  
  <!-- Child - Single size -->
  <button class="h-full flex flex-col">
    <div class="flex-1">Content</div>
    <div>Badge</div>
  </button>
  
  <!-- Child - Multi size -->
  <div class="flex flex-col h-full">
    <div>Title</div>
    <div class="flex-1">Variants</div>
    <div>Badge</div>
  </div>
</div>
```

### Line Clamp
- Webkit-specific CSS
- Works on all modern browsers
- Fallback: overflow hidden
- 2 lines max with ellipsis

### Min Height
- `min-h-[2.5rem]` = 40px
- 2 lines × 1.25rem line-height
- Ensures consistent spacing

## Files Changed
- `frontend/src/views/OrderView.vue`

## Testing Checklist
- [ ] Customer name section nổi bật, không thể bỏ qua
- [ ] Label và placeholder rõ ràng
- [ ] Focus state hoạt động tốt
- [ ] Menu items có chiều cao đồng đều
- [ ] Tên món dài tự động cắt với "..."
- [ ] Badge position consistent
- [ ] Single-size và multi-size items aligned
- [ ] Touch targets đủ lớn (44px+)
- [ ] Active states hoạt động mượt
- [ ] Grid responsive trên các màn hình
- [ ] Scroll smooth
- [ ] Category tabs hoạt động tốt
