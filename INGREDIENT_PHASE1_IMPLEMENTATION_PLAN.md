# Ingredient Management - Phase 1 Implementation Plan

## Overview
Phase 1 focuses on quick wins that provide immediate value with minimal effort.

## Sprint 1: Quick Stock Operations (Week 1)

### US-1.1: Quick Stock IN Button ✅ (Already Implemented)
Current "Điều chỉnh" button already provides this functionality.

**Improvements Needed:**
- Make it more prominent
- Add quick access from card
- Simplify modal for common case

### US-1.2: Quick Stock OUT Button (NEW)
**Implementation:**

1. **Add Quick OUT Button to Card**
   ```vue
   <button @click="quickStockOut(ingredient)" 
     class="bg-orange-500 text-white">
     ➖ Xuất
   </button>
   ```

2. **Create Quick OUT Modal**
   - Simpler than full adjust modal
   - Only quantity and reason
   - Predefined reasons dropdown
   - No price input (already purchased)

3. **Backend Support**
   - Use existing adjust stock endpoint
   - Type: "remove"
   - No cost_per_unit needed

**Files to Modify:**
- `frontend/src/views/IngredientManagementView.vue`
- Add `showQuickOutModal` state
- Add `quickStockOut()` function
- Add quick out modal template

**Effort:** 2 hours

### US-1.3: Quick View Stock Level (NEW)
**Implementation:**

1. **Add Progress Bar to Card**
   ```vue
   <div class="w-full bg-gray-200 rounded-full h-2">
     <div class="bg-green-500 h-2 rounded-full" 
       :style="{ width: stockPercentage(ingredient) + '%' }">
     </div>
   </div>
   ```

2. **Calculate Stock Percentage**
   ```javascript
   const stockPercentage = (ingredient) => {
     const target = ingredient.min_stock * 3 // 3x min is "full"
     return Math.min((ingredient.quantity / target) * 100, 100)
   }
   ```

3. **Color Coding**
   - Green: > 2x min_stock
   - Yellow: 1-2x min_stock
   - Red: < min_stock

**Files to Modify:**
- `frontend/src/views/IngredientManagementView.vue`
- Add progress bar to ingredient card
- Add `stockPercentage()` function
- Add `getStockColor()` function

**Effort:** 1 hour

## Sprint 2: Filtering & Alerts (Week 2)

### US-2.1: Filter by Category (NEW)
**Implementation:**

1. **Add Category Filter Chips**
   ```vue
   <div class="flex gap-2 overflow-x-auto pb-2">
     <button v-for="cat in categories" 
       @click="filterCategory = cat"
       :class="filterCategory === cat ? 'bg-blue-500' : 'bg-gray-200'">
       {{ cat }}
     </button>
   </div>
   ```

2. **Update Filtered List**
   ```javascript
   const filteredIngredients = computed(() => {
     let filtered = ingredients.value
     
     if (filterCategory.value) {
       filtered = filtered.filter(i => i.category === filterCategory.value)
     }
     
     if (searchQuery.value) {
       filtered = filtered.filter(i => 
         i.name.toLowerCase().includes(searchQuery.value.toLowerCase())
       )
     }
     
     return filtered
   })
   ```

**Files to Modify:**
- `frontend/src/views/IngredientManagementView.vue`
- Add category filter chips
- Add `filterCategory` state
- Update `filteredIngredients` computed

**Effort:** 2 hours

### US-2.2: Filter by Stock Status (NEW)
**Implementation:**

1. **Add Status Filter Buttons**
   ```vue
   <div class="grid grid-cols-4 gap-2">
     <button @click="filterStatus = 'all'">Tất cả</button>
     <button @click="filterStatus = 'in_stock'">Đủ hàng</button>
     <button @click="filterStatus = 'low'">Sắp hết</button>
     <button @click="filterStatus = 'out'">Hết hàng</button>
   </div>
   ```

2. **Update Filter Logic**
   ```javascript
   if (filterStatus.value === 'in_stock') {
     filtered = filtered.filter(i => i.quantity > i.min_stock)
   } else if (filterStatus.value === 'low') {
     filtered = filtered.filter(i => i.quantity > 0 && i.quantity <= i.min_stock)
   } else if (filterStatus.value === 'out') {
     filtered = filtered.filter(i => i.quantity === 0)
   }
   ```

**Files to Modify:**
- `frontend/src/views/IngredientManagementView.vue`
- Add status filter buttons
- Add `filterStatus` state
- Update `filteredIngredients` computed

**Effort:** 1 hour

### US-3.1: Low Stock Banner (NEW)
**Implementation:**

1. **Add Banner Component**
   ```vue
   <div v-if="lowStockCount > 0" 
     class="bg-yellow-100 border-l-4 border-yellow-500 p-3 mb-4">
     <div class="flex items-center justify-between">
       <div>
         <p class="font-bold text-yellow-800">
           ⚠️ {{ lowStockCount }} nguyên liệu sắp hết
         </p>
         <p class="text-sm text-yellow-700">
           {{ topLowStockItems }}
         </p>
       </div>
       <button @click="showLowStock" 
         class="text-yellow-800 underline">
         Xem
       </button>
     </div>
   </div>
   ```

2. **Calculate Low Stock Items**
   ```javascript
   const topLowStockItems = computed(() => {
     const low = ingredients.value
       .filter(i => i.quantity > 0 && i.quantity <= i.min_stock)
       .slice(0, 3)
       .map(i => i.name)
     return low.join(', ')
   })
   ```

**Files to Modify:**
- `frontend/src/views/IngredientManagementView.vue`
- Add low stock banner
- Add `topLowStockItems` computed
- Update `showLowStock()` to set filter

**Effort:** 1 hour

## Implementation Order

### Day 1-2: Quick Stock Operations
1. ✅ Review existing adjust stock functionality
2. 🔨 Add quick stock OUT button and modal
3. 🔨 Add stock level progress bars
4. ✅ Test on mobile devices

### Day 3-4: Category Filtering
1. 🔨 Add category filter chips
2. 🔨 Update filtered list logic
3. 🔨 Add category counts
4. ✅ Test filtering combinations

### Day 5: Status Filtering
1. 🔨 Add status filter buttons
2. 🔨 Update filtered list logic
3. 🔨 Add status counts
4. ✅ Test all filter combinations

### Day 6-7: Low Stock Banner
1. 🔨 Add banner component
2. 🔨 Calculate low stock items
3. 🔨 Add click to filter
4. ✅ Test banner behavior
5. ✅ Final testing and polish

## Testing Checklist

### Quick Stock OUT
- [ ] Button appears on all ingredient cards
- [ ] Modal opens with correct ingredient
- [ ] Predefined reasons work
- [ ] Custom reason can be entered
- [ ] Stock decreases correctly
- [ ] History record created
- [ ] No expense created
- [ ] Success feedback shown
- [ ] Modal closes after success

### Stock Level Progress Bar
- [ ] Shows on all ingredient cards
- [ ] Color changes based on stock level
- [ ] Percentage calculated correctly
- [ ] Responsive on mobile
- [ ] Updates after stock changes

### Category Filter
- [ ] All categories shown
- [ ] "Tất cả" shows all ingredients
- [ ] Clicking category filters list
- [ ] Active category highlighted
- [ ] Count shown per category
- [ ] Works with search
- [ ] Persists during session

### Status Filter
- [ ] All status buttons shown
- [ ] Clicking status filters list
- [ ] Active status highlighted
- [ ] Count shown per status
- [ ] Works with category filter
- [ ] Works with search
- [ ] Persists during session

### Low Stock Banner
- [ ] Shows when low stock items exist
- [ ] Shows correct count
- [ ] Shows top 3 items
- [ ] Clicking "Xem" filters list
- [ ] Dismissible (optional)
- [ ] Reappears on refresh
- [ ] Updates when stock changes

## Success Criteria

### Performance
- All operations complete in < 2 seconds
- No UI lag when filtering
- Smooth animations

### UX
- Intuitive button placement
- Clear visual feedback
- Mobile-friendly touch targets
- Consistent with existing design

### Functionality
- All filters work independently
- Filters can be combined
- Stock operations are accurate
- History is recorded correctly

## Rollout Plan

### Week 1
- Deploy quick stock operations
- Monitor for issues
- Gather user feedback

### Week 2
- Deploy filtering features
- Deploy low stock banner
- Monitor usage metrics
- Iterate based on feedback

## Metrics to Track

### Usage
- Number of quick stock OUT operations per day
- Filter usage frequency
- Low stock banner click-through rate

### Performance
- Time to complete stock operation
- Filter response time
- Page load time

### Business Impact
- Reduction in stock-outs
- Improvement in inventory accuracy
- Time saved on inventory management

## Next Phase Preview

After Phase 1 completion, Phase 2 will focus on:
- Sort options
- Recipe integration
- Reorder suggestions
- Enhanced analytics

Estimated start: Week 3
