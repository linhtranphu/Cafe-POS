# Order UI with Inline Controls & Grouped Categories

## Overview
Redesign order creation UI with:
1. **Grouped menu items** by category (no tabs)
2. **Inline +/- controls** on each item
3. **Floating total button** for confirmation
4. **No separate cart view**

## Goals
- Show all items at once, grouped by category
- Reduce taps needed to create an order
- Make quantity adjustment more intuitive
- Natural scrolling experience like a menu

## Design Principles
1. **Grouped layout** - All categories visible, scroll to browse
2. **Sticky headers** - Category names stay visible
3. **Inline controls** - +/- buttons directly on items
4. **Floating total** - Always visible, tap to confirm

## User Stories

### Story 1: Browse All Categories
**As a** waiter  
**I want to** see all menu items grouped by category  
**So that** I can browse naturally without switching tabs

**Acceptance Criteria:**
- [ ] All categories shown in one scrollable view
- [ ] Category headers are sticky while scrolling
- [ ] Items grouped under their category
- [ ] Smooth scrolling between categories

### Story 2: Quick Add Items
**As a** waiter  
**I want to** quickly add items with + button  
**So that** I can create orders faster

**Acceptance Criteria:**
- [ ] Each menu item shows a "+ Thêm" button when qty = 0
- [ ] Tapping + adds 1 to quantity
- [ ] Item shows inline +/- controls when qty > 0
- [ ] Quantity is displayed between +/- buttons

### Story 2: Adjust Quantity
**As a** waiter  
**I want to** adjust quantity inline  
**So that** I don't need to open a separate cart view

**Acceptance Criteria:**
- [ ] + button increases quantity
- [ ] - button decreases quantity
- [ ] When qty reaches 0, item returns to "+ Thêm" state
- [ ] Smooth animations for feedback

### Story 3: See Total
**As a** waiter  
**I want to** see total items and price at all times  
**So that** I know the order status

**Acceptance Criteria:**
- [ ] Floating button at bottom shows total items and price
- [ ] Button is always visible while scrolling
- [ ] Button is disabled when cart is empty
- [ ] Tapping button confirms order

## Technical Design

### State Management
```javascript
// Simplified cart structure
const cart = ref({})  // { itemId: quantity } or { `${itemId}_${variantId}`: quantity }

// Helper functions
const getItemQty = (itemId, variantId = null) => {
  const key = variantId ? `${itemId}_${variantId}` : itemId
  return cart.value[key] || 0
}

const addItem = (itemId, variantId = null) => {
  const key = variantId ? `${itemId}_${variantId}` : itemId
  cart.value[key] = (cart.value[key] || 0) + 1
}

const removeItem = (itemId, variantId = null) => {
  const key = variantId ? `${itemId}_${variantId}` : itemId
  if (cart.value[key] > 1) {
    cart.value[key]--
  } else {
    delete cart.value[key]
  }
}

// Group items by category
const groupedItems = computed(() => {
  const groups = {}
  
  filteredMenuItems.value.forEach(item => {
    const categoryId = item.category
    if (!groups[categoryId]) {
      const category = categories.value.find(c => c.id === categoryId)
      groups[categoryId] = {
        category: category,
        items: []
      }
    }
    groups[categoryId].items.push(item)
  })
  
  return Object.values(groups).filter(g => g.items.length > 0)
})
```

### Component Structure
```
CreateOrderView
├── Header (title, back, search)
├── GroupedMenuList (scrollable)
│   └── CategoryGroup (for each category)
│       ├── CategoryHeader (sticky)
│       └── MenuItemsGrid
│           └── MenuItem
│               ├── ItemInfo (name, price)
│               └── QuantityControls
│                   ├── AddButton (when qty = 0)
│                   └── InlineControls (when qty > 0)
│                       ├── MinusButton
│                       ├── QuantityDisplay
│                       └── PlusButton
└── FloatingTotalButton
```

## UI Components

### MenuItem Card
```vue
<div class="menu-item-card">
  <!-- Item Info -->
  <div class="item-info">
    <h3>{{ item.name }}</h3>
    <p class="price">{{ formatPrice(item.price) }}</p>
  </div>
  
  <!-- Quantity Controls -->
  <div class="quantity-controls">
    <!-- When qty = 0 -->
    <button v-if="qty === 0" @click="addItem">
      + Thêm
    </button>
    
    <!-- When qty > 0 -->
    <div v-else class="inline-controls">
      <button @click="removeItem">-</button>
      <span class="qty">{{ qty }}</span>
      <button @click="addItem">+</button>
    </div>
  </div>
</div>
```

### Floating Total Button
```vue
<button 
  v-if="totalItems > 0"
  @click="confirmOrder"
  class="floating-total-button">
  <div class="total-info">
    <span>🛒 {{ totalItems }} món</span>
    <span>{{ formatPrice(totalPrice) }}</span>
  </div>
  <span class="confirm-text">Xác nhận →</span>
</button>
```

## Styling

### Colors
- Primary: `#3B82F6` (blue-500)
- Success: `#10B981` (green-500)
- Border: `#E5E7EB` (gray-200)
- Selected Border: `#3B82F6` (blue-500)

### Spacing
- Card padding: `16px`
- Grid gap: `12px`
- Button height: `44px` (min touch target)

### Animations
```css
.menu-item-card {
  transition: all 0.2s ease;
}

.menu-item-card.has-quantity {
  border-color: #3B82F6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.1);
}

.quantity-controls button {
  transition: transform 0.1s ease;
}

.quantity-controls button:active {
  transform: scale(0.95);
}
```

## Implementation Tasks

### Phase 1: State Management
- [ ] Refactor cart to use object structure
- [ ] Implement getItemQty helper
- [ ] Implement addItem helper
- [ ] Implement removeItem helper
- [ ] Add computed for totalItems
- [ ] Add computed for totalPrice

### Phase 2: UI Components
- [ ] Update MenuItem card layout
- [ ] Add inline quantity controls
- [ ] Style selected state
- [ ] Add animations
- [ ] Create FloatingTotalButton component

### Phase 3: Remove Old Cart
- [ ] Remove cart list section
- [ ] Remove cart-related functions
- [ ] Update confirmOrder flow
- [ ] Test all scenarios

### Phase 4: Polish
- [ ] Add haptic feedback (if supported)
- [ ] Optimize animations
- [ ] Test on different screen sizes
- [ ] Add loading states

## Testing Scenarios

### Test 1: Add Single Item
1. Open create order
2. Tap "+ Thêm" on an item
3. Verify qty shows as 1
4. Verify +/- controls appear
5. Verify floating button shows "1 món"

### Test 2: Adjust Quantity
1. Add item (qty = 1)
2. Tap + button
3. Verify qty increases to 2
4. Tap - button
5. Verify qty decreases to 1
6. Tap - again
7. Verify item returns to "+ Thêm" state

### Test 3: Multiple Items
1. Add multiple different items
2. Verify each shows correct quantity
3. Verify total is correct
4. Tap floating button
5. Verify order is created correctly

### Test 4: Variants
1. Select item with variants
2. Tap variant option
3. Verify qty increases for that variant
4. Verify other variants remain independent

## Success Criteria
- ✅ Taps reduced by ~40%
- ✅ No separate cart view
- ✅ Floating total always visible
- ✅ Smooth animations
- ✅ Clear visual feedback
- ✅ Works on all screen sizes

## References
- Design doc: `ORDER_UI_REDESIGN.md`
- Current implementation: `frontend/src/views/OrderView.vue`
