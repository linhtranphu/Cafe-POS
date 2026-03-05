# Order UI/UX Redesign - Waiter Screen (Grouped by Category)

## Design Philosophy

**Core Principle:** Minimize taps, maximize visibility, natural grouping

### User Flow
```
1. Open create order → See ALL items grouped by category
2. Scroll to find category → Items are grouped with headers
3. Tap + on item → Quantity increases (shown on item)
4. Tap - to decrease
5. See floating total button (always visible)
6. Tap total button → Confirm order
```

## UI Components

### 1. Header (Fixed Top)
```
┌─────────────────────────────────────┐
│ ← Back    🍽️ Tạo Order             │
│                                     │
│ 🔍 Search...                    [×] │
└─────────────────────────────────────┘
```

### 2. Grouped Menu Items (Scrollable)
```
┌─────────────────────────────────────┐
│ ☕ CÀ PHÊ                           │ ← Category Header
├─────────────────────────────────────┤
│ ┌──────────┐  ┌──────────┐         │
│ │Cà phê đá │  │Bạc xỉu   │         │
│ │25,000đ   │  │30,000đ   │         │
│ │[+ Thêm]  │  │[-] 2 [+] │         │
│ └──────────┘  └──────────┘         │
├─────────────────────────────────────┤
│ 🍵 TRÀ                              │ ← Category Header
├─────────────────────────────────────┤
│ ┌──────────┐  ┌──────────┐         │
│ │Trà đào   │  │Trà sữa   │         │
│ │35,000đ   │  │40,000đ   │         │
│ │[+ Thêm]  │  │[+ Thêm]  │         │
│ └──────────┘  └──────────┘         │
└─────────────────────────────────────┘
```

### 3. Menu Item Card (Enhanced)
```
┌──────────────────────┐
│ Cà phê sữa đá        │
│ 25,000đ              │
│                      │
│ [  -  ]  2  [  +  ]  │ ← Inline controls
└──────────────────────┘
```

When quantity = 0:
```
┌──────────────────────┐
│ Cà phê sữa đá        │
│ 25,000đ              │
│                      │
│     [  + Thêm  ]     │ ← Single add button
└──────────────────────┘
```

### 4. Floating Total Button (Fixed Bottom)
```
┌─────────────────────────────────────┐
│  🛒 3 món • 75,000đ    [Xác nhận] │
└─────────────────────────────────────┘
```

## Key Features

### ✅ Grouped by Category
- **Category headers** - Clear visual separation
- **Sticky headers** - Category name stays visible while scrolling
- **All items visible** - No need to switch tabs
- **Natural flow** - Scroll through all categories

### ✅ Inline Quantity Controls
- **Add button** when qty = 0
- **+/- controls** when qty > 0
- **Quantity badge** always visible
- **Smooth animations** for feedback

### ✅ No Separate Cart View
- Cart info in floating button
- Tap to see summary before confirm
- Less screen space used
- Cleaner interface

### ✅ Search Across All Categories
- Search filters items across all categories
- Matching items shown with their category
- Easy to find specific items

### ✅ Visual Feedback
- Selected items have colored border
- Quantity badge with accent color
- Smooth scale animations on tap
- Clear visual hierarchy

## Layout Structure

```
┌─────────────────────────────────────┐
│ Header (Search)                     │ ← Fixed
├─────────────────────────────────────┤
│                                     │
│ ☕ CÀ PHÊ                           │ ← Sticky Category Header
│ ┌──────┐ ┌──────┐ ┌──────┐        │
│ │Item 1│ │Item 2│ │Item 3│        │
│ └──────┘ └──────┘ └──────┘        │
│                                     │
│ 🍵 TRÀ                              │ ← Sticky Category Header
│ ┌──────┐ ┌──────┐                  │
│ │Item 4│ │Item 5│                  │
│ └──────┘ └──────┘                  │
│                                     │
│ 🥤 NƯỚC NGỌT                        │ ← Sticky Category Header
│ ┌──────┐ ┌──────┐ ┌──────┐        │
│ │Item 6│ │Item 7│ │Item 8│        │
│ └──────┘ └──────┘ └──────┘        │
│                                     │
├─────────────────────────────────────┤
│ 🛒 3 món • 75,000đ  [Xác nhận]    │ ← Fixed
└─────────────────────────────────────┘
```

## Advantages of Grouped Layout

### vs Tab-based Navigation
| Aspect | Tabs | Grouped |
|--------|------|---------|
| **Visibility** | One category at a time | All categories visible |
| **Navigation** | Tap to switch | Natural scroll |
| **Context** | Lost when switching | Always visible |
| **Speed** | Slower (tap + wait) | Faster (just scroll) |
| **Discovery** | Limited | Better |

### User Benefits
1. **See everything** - No hidden items
2. **Natural browsing** - Scroll like a menu
3. **Quick comparison** - See items from different categories
4. **Less taps** - No tab switching needed
5. **Better memory** - Visual landmarks (category headers)

## Implementation Notes

### Category Header Component
```vue
<div class="category-header sticky top-0 z-10">
  <h3 class="category-title">
    <span class="icon">{{ category.icon }}</span>
    <span class="name">{{ category.name }}</span>
    <span class="count">({{ itemCount }})</span>
  </h3>
</div>
```

### Grouped Items Structure
```javascript
// Group items by category
const groupedItems = computed(() => {
  const groups = {}
  
  filteredMenuItems.value.forEach(item => {
    const categoryId = item.category
    if (!groups[categoryId]) {
      groups[categoryId] = {
        category: categories.value.find(c => c.id === categoryId),
        items: []
      }
    }
    groups[categoryId].items.push(item)
  })
  
  return Object.values(groups)
})
```

### Render Template
```vue
<div class="menu-container">
  <div v-for="group in groupedItems" :key="group.category.id" class="category-group">
    <!-- Category Header (Sticky) -->
    <div class="category-header">
      <span>{{ group.category.icon }}</span>
      <span>{{ group.category.name }}</span>
      <span>({{ group.items.length }})</span>
    </div>
    
    <!-- Items Grid -->
    <div class="items-grid">
      <MenuItem 
        v-for="item in group.items" 
        :key="item.id"
        :item="item"
        :quantity="getItemQty(item.id)"
        @add="addItem(item.id)"
        @remove="removeItem(item.id)"
      />
    </div>
  </div>
</div>
```

## Styling

### Category Header
```css
.category-header {
  position: sticky;
  top: 0;
  z-index: 10;
  background: linear-gradient(to bottom, #f9fafb 0%, #f3f4f6 100%);
  padding: 12px 16px;
  border-bottom: 2px solid #e5e7eb;
  font-weight: 700;
  font-size: 16px;
  color: #374151;
  display: flex;
  align-items: center;
  gap: 8px;
}

.category-header .icon {
  font-size: 20px;
}

.category-header .count {
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
}
```

### Category Group
```css
.category-group {
  margin-bottom: 24px;
}

.category-group:last-child {
  margin-bottom: 80px; /* Space for floating button */
}

.items-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  padding: 16px;
}
```

## Search Behavior

When user searches:
1. Filter items across all categories
2. Show only categories that have matching items
3. Highlight matching text (optional)
4. Keep category grouping

```javascript
const filteredMenuItems = computed(() => {
  let items = menuItems.value
  
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    items = items.filter(item => 
      item.name.toLowerCase().includes(query)
    )
  }
  
  return items
})

const groupedItems = computed(() => {
  // Group filtered items by category
  // Only show categories that have items
})
```

## Performance Optimization

### Virtual Scrolling (Optional)
For menus with 100+ items:
```javascript
import { useVirtualList } from '@vueuse/core'

const { list, containerProps, wrapperProps } = useVirtualList(
  groupedItems,
  { itemHeight: 200 }
)
```

### Lazy Loading Images (Future)
```vue
<img 
  :src="item.image" 
  loading="lazy"
  class="item-image"
/>
```

## Success Metrics

- ⏱️ Time to find item: < 5 seconds
- 👆 Taps to create order: < 10
- 📊 Items per screen: 4-6 visible
- 🔄 Scroll performance: 60fps
- 😊 User satisfaction: High

---

**Status:** Ready for implementation
**Priority:** High
**Estimated effort:** 6-8 hours

## Color Scheme

```css
Primary: Blue (#3B82F6)
Success: Green (#10B981)
Accent: Orange (#F59E0B)
Background: Gray-50 (#F9FAFB)
Border: Gray-200 (#E5E7EB)
```

## Interaction States

### Item Card States
1. **Default:** White background, gray border
2. **Has Quantity:** Blue border, blue badge
3. **Active (tap):** Scale 0.98, shadow increase
4. **Disabled:** Opacity 0.5, no interaction

### Button States
1. **Add (+):** Blue background, white text
2. **Remove (-):** White background, blue border
3. **Active:** Darker shade, scale 0.95
4. **Disabled:** Gray, opacity 0.5

## Responsive Behavior

### Mobile (< 640px)
- 2-column grid
- Compact spacing
- Larger touch targets

### Tablet (≥ 640px)
- 3-column grid
- More spacing
- Larger cards

## Accessibility

- ✅ Min touch target: 44x44px
- ✅ High contrast ratios
- ✅ Clear focus states
- ✅ Semantic HTML
- ✅ ARIA labels

## Performance

- ✅ Virtual scrolling for large menus
- ✅ Debounced search
- ✅ Optimized re-renders
- ✅ Smooth 60fps animations

## Implementation Notes

### State Management
```javascript
// Simple cart state
const cart = ref({})  // { itemId: quantity }

// Add item
const addItem = (itemId) => {
  cart.value[itemId] = (cart.value[itemId] || 0) + 1
}

// Remove item
const removeItem = (itemId) => {
  if (cart.value[itemId] > 1) {
    cart.value[itemId]--
  } else {
    delete cart.value[itemId]
  }
}

// Get quantity
const getQty = (itemId) => cart.value[itemId] || 0

// Total items
const totalItems = computed(() => 
  Object.values(cart.value).reduce((sum, qty) => sum + qty, 0)
)

// Total price
const totalPrice = computed(() => {
  return Object.entries(cart.value).reduce((sum, [id, qty]) => {
    const item = menuItems.value.find(i => i.id === id)
    return sum + (item?.price || 0) * qty
  }, 0)
})
```

### Component Structure
```
OrderView.vue
├── Header (search, back, total)
├── CategoryTabs (sticky)
├── MenuGrid
│   └── MenuItem (with inline controls)
└── FloatingTotalButton
```

## User Testing Scenarios

### Scenario 1: Quick Order
```
1. Open create order
2. Tap "Cà phê" category
3. Tap + on "Cà phê sữa đá" (2 times)
4. Tap + on "Bạc xỉu" (1 time)
5. Tap floating total button
6. Confirm
```
**Expected:** 5 taps total, < 10 seconds

### Scenario 2: Complex Order
```
1. Open create order
2. Search "trà"
3. Tap + on multiple items
4. Switch to "Nước" category
5. Add more items
6. Review total
7. Confirm
```
**Expected:** Smooth, no confusion

### Scenario 3: Modify Order
```
1. Add items
2. Tap - to decrease
3. Tap - again to remove
4. Add different items
5. Confirm
```
**Expected:** Clear feedback, no errors

## Comparison: Old vs New

### Old Design
- ❌ Click item → Add to cart
- ❌ Cart list at bottom (takes space)
- ❌ Scroll to see cart
- ❌ Multiple steps to adjust quantity
- ❌ ~8-10 taps for simple order

### New Design
- ✅ Inline +/- controls
- ✅ No separate cart view
- ✅ Floating total (always visible)
- ✅ Direct quantity adjustment
- ✅ ~4-5 taps for simple order

## Success Metrics

- ⏱️ Time to create order: < 30 seconds
- 👆 Taps per order: < 10
- 😊 User satisfaction: High
- 🐛 Error rate: < 5%
- 🔄 Return rate: Low

---

**Status:** Ready for implementation
**Priority:** High
**Estimated effort:** 4-6 hours
