# Menu Size Variants - Implementation Plan (No Migration Needed)

## 🎉 Tình huống Lý tưởng

**Status**: Chưa có menu món trong database
**Advantage**: 
- ✅ Không cần migration
- ✅ Không có data cũ phải xử lý
- ✅ Không có order history phải maintain
- ✅ Có thể implement clean architecture từ đầu

**Risk Level**: LOW → MEDIUM (chỉ còn implementation complexity)

---

## 📋 Simplified Implementation Plan

### Phase 1: Backend Domain Model (2 days)

#### 1.1 Update Menu Domain
**File**: `backend/domain/menu/menu.go`

**Changes**:
```go
// Add MenuItemVariant struct
type MenuItemVariant struct {
    ID          string               `bson:"id" json:"id"`
    Name        string               `bson:"name" json:"name"`
    Price       float64              `bson:"price" json:"price"`
    Ingredients []Ingredient         `bson:"ingredients" json:"ingredients"`
    Available   bool                 `bson:"available" json:"available"`
    IsDefault   bool                 `bson:"is_default" json:"is_default"`
    
    // Cost tracking per variant
    CurrentCost          float64    `bson:"current_cost" json:"current_cost"`
    CostLastCalculatedAt time.Time  `bson:"cost_last_calculated_at" json:"cost_last_calculated_at"`
    CostStatus           CostStatus `bson:"cost_status" json:"cost_status"`
}

// Update MenuItem struct
type MenuItem struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name        string             `bson:"name" json:"name"`
    Category    string             `bson:"category" json:"category"`
    Description string             `bson:"description" json:"description"`
    Available   bool               `bson:"available" json:"available"`
    
    // NEW: Variants support
    HasVariants bool               `bson:"has_variants" json:"has_variants"`
    Variants    []MenuItemVariant  `bson:"variants,omitempty" json:"variants,omitempty"`
    
    // For backward compatibility (single-size items)
    Price       float64      `bson:"price,omitempty" json:"price,omitempty"`
    Ingredients []Ingredient `bson:"ingredients,omitempty" json:"ingredients,omitempty"`
    CurrentCost float64      `bson:"current_cost,omitempty" json:"current_cost,omitempty"`
    CostLastCalculatedAt time.Time  `bson:"cost_last_calculated_at,omitempty" json:"cost_last_calculated_at,omitempty"`
    CostStatus  CostStatus   `bson:"cost_status,omitempty" json:"cost_status,omitempty"`
    
    CreatedAt   time.Time `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

// Add helper methods
func (m *MenuItem) GetDefaultVariant() *MenuItemVariant
func (m *MenuItem) GetVariantByID(variantID string) *MenuItemVariant
func (m *MenuItem) GetPrice(variantID string) float64
func (m *MenuItem) GetIngredients(variantID string) []Ingredient
func (m *MenuItem) Validate() error
```

**Tasks**:
- [ ] Add MenuItemVariant struct
- [ ] Update MenuItem struct
- [ ] Add helper methods
- [ ] Add validation logic
- [ ] Write unit tests

**Estimated Time**: 4 hours

#### 1.2 Update Request/Response DTOs
**File**: `backend/domain/menu/menu.go`

```go
type CreateMenuItemRequest struct {
    Name        string              `json:"name" binding:"required"`
    Category    string              `json:"category" binding:"required"`
    Description string              `json:"description"`
    
    // For single-size items
    Price       float64      `json:"price"`
    Ingredients []Ingredient `json:"ingredients"`
    
    // For multi-size items
    HasVariants bool               `json:"has_variants"`
    Variants    []MenuItemVariant  `json:"variants"`
}

type UpdateMenuItemRequest struct {
    Name        string              `json:"name"`
    Category    string              `json:"category"`
    Description string              `json:"description"`
    Available   *bool               `json:"available"`
    
    Price       float64      `json:"price"`
    Ingredients []Ingredient `json:"ingredients"`
    
    HasVariants *bool              `json:"has_variants"`
    Variants    []MenuItemVariant  `json:"variants"`
}
```

**Tasks**:
- [ ] Update CreateMenuItemRequest
- [ ] Update UpdateMenuItemRequest
- [ ] Add validation tags

**Estimated Time**: 1 hour

#### 1.3 Update Order Domain
**File**: `backend/domain/order/order.go`

```go
type OrderItem struct {
    MenuItemID  primitive.ObjectID `bson:"menu_item_id" json:"menu_item_id"`
    VariantID   string             `bson:"variant_id,omitempty" json:"variant_id,omitempty"` // NEW
    Name        string             `bson:"name" json:"name"`
    VariantName string             `bson:"variant_name,omitempty" json:"variant_name,omitempty"` // NEW
    Price       float64            `bson:"price" json:"price"`
    Quantity    int                `bson:"quantity" json:"quantity"`
    Note        string             `bson:"note,omitempty" json:"note,omitempty"`
    Subtotal    float64            `bson:"subtotal" json:"subtotal"`
}
```

**Tasks**:
- [ ] Add VariantID field
- [ ] Add VariantName field
- [ ] Update validation

**Estimated Time**: 1 hour

---

### Phase 2: Backend Service Layer (2 days)

#### 2.1 Update Menu Service
**File**: `backend/application/services/menu.go`

**Changes**:
```go
func (s *MenuService) CreateMenuItem(ctx context.Context, req *menu.CreateMenuItemRequest) (*menu.MenuItem, error) {
    item := &menu.MenuItem{
        Name:        req.Name,
        Category:    req.Category,
        Description: req.Description,
        Available:   true,
        HasVariants: req.HasVariants,
    }
    
    if req.HasVariants {
        // Validate variants
        if len(req.Variants) == 0 {
            return nil, errors.New("variants required when has_variants=true")
        }
        
        // Ensure exactly one default
        defaultCount := 0
        for _, v := range req.Variants {
            if v.IsDefault {
                defaultCount++
            }
        }
        if defaultCount != 1 {
            return nil, errors.New("must have exactly one default variant")
        }
        
        item.Variants = req.Variants
    } else {
        // Single-size item
        item.Price = req.Price
        item.Ingredients = req.Ingredients
    }
    
    // Validate
    if err := item.Validate(); err != nil {
        return nil, err
    }
    
    err := s.menuRepo.Create(ctx, item)
    if err != nil {
        return nil, err
    }
    
    return item, nil
}
```

**Tasks**:
- [ ] Update CreateMenuItem
- [ ] Update UpdateMenuItem
- [ ] Add variant validation
- [ ] Write unit tests

**Estimated Time**: 4 hours

#### 2.2 Update Order Service
**File**: `backend/application/services/order_service.go`

**Changes**:
```go
func (s *OrderService) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.Order, error) {
    // ... existing code ...
    
    // Validate and populate order items
    for i, item := range req.Items {
        menuItem, err := s.menuRepo.FindByID(ctx, item.MenuItemID)
        if err != nil {
            return nil, err
        }
        
        // Get price based on variant
        var price float64
        var variantName string
        
        if menuItem.HasVariants {
            // Must have variant_id
            if item.VariantID == "" {
                return nil, errors.New("variant_id required for multi-size item")
            }
            
            variant := menuItem.GetVariantByID(item.VariantID)
            if variant == nil {
                return nil, errors.New("invalid variant_id")
            }
            
            price = variant.Price
            variantName = variant.Name
        } else {
            // Single-size item
            price = menuItem.Price
        }
        
        req.Items[i].Price = price
        req.Items[i].VariantName = variantName
        req.Items[i].Subtotal = price * float64(item.Quantity)
    }
    
    // ... rest of code ...
}
```

**Tasks**:
- [ ] Update CreateOrder to handle variants
- [ ] Update EditOrder to handle variants
- [ ] Add variant validation
- [ ] Write unit tests

**Estimated Time**: 4 hours

#### 2.3 Update Cost Calculator
**File**: `backend/application/services/cost_calculator_service.go`

**Changes**:
```go
func (s *CostCalculatorService) CalculateMenuItemCost(ctx context.Context, menuItemID primitive.ObjectID) error {
    menuItem, err := s.menuRepo.FindByID(ctx, menuItemID)
    if err != nil {
        return err
    }
    
    if menuItem.HasVariants {
        // Calculate cost for each variant
        for i := range menuItem.Variants {
            cost, status := s.calculateIngredientsCost(ctx, menuItem.Variants[i].Ingredients)
            menuItem.Variants[i].CurrentCost = cost
            menuItem.Variants[i].CostStatus = status
            menuItem.Variants[i].CostLastCalculatedAt = time.Now()
        }
    } else {
        // Single-size item
        cost, status := s.calculateIngredientsCost(ctx, menuItem.Ingredients)
        menuItem.CurrentCost = cost
        menuItem.CostStatus = status
        menuItem.CostLastCalculatedAt = time.Now()
    }
    
    return s.menuRepo.Update(ctx, menuItemID, menuItem)
}
```

**Tasks**:
- [ ] Update CalculateMenuItemCost
- [ ] Handle variant costs
- [ ] Write unit tests

**Estimated Time**: 2 hours

---

### Phase 3: Backend API Layer (1 day)

#### 3.1 Update Menu Handler
**File**: `backend/interfaces/http/menu_handler.go`

**No major changes needed** - handlers already use DTOs

**Tasks**:
- [ ] Test API endpoints
- [ ] Update API documentation
- [ ] Add example requests

**Estimated Time**: 2 hours

#### 3.2 Update Order Handler
**File**: `backend/interfaces/http/order_handler.go`

**No major changes needed** - handlers already use DTOs

**Tasks**:
- [ ] Test API endpoints
- [ ] Update API documentation

**Estimated Time**: 1 hour

---

### Phase 4: Frontend Data Layer (1 day)

#### 4.1 Update Menu Store
**File**: `frontend/src/stores/menu.js`

**No changes needed** - store already handles API responses

**Tasks**:
- [ ] Verify store handles variants
- [ ] Add helper methods if needed

**Estimated Time**: 1 hour

#### 4.2 Update Order Store
**File**: `frontend/src/stores/order.js`

**Changes**:
```javascript
// Update addItem to handle variants
addItem(menuItem, variant = null) {
  const item = {
    menu_item_id: menuItem.id,
    name: menuItem.name,
    price: variant ? variant.price : menuItem.price,
    quantity: 1
  }
  
  if (variant) {
    item.variant_id = variant.id
    item.variant_name = variant.name
  }
  
  this.items.push(item)
}
```

**Tasks**:
- [ ] Update addItem method
- [ ] Update editItem method
- [ ] Handle variant selection

**Estimated Time**: 2 hours

---

### Phase 5: Frontend UI (3 days)

#### 5.1 Update MenuView Display
**File**: `frontend/src/views/MenuView.vue`

**Changes**:
```vue
<template>
  <div v-for="item in menuItems" :key="item.id" class="menu-item">
    <h3>{{ item.name }}</h3>
    <p>{{ item.description }}</p>
    
    <!-- Single-size item -->
    <div v-if="!item.has_variants" class="single-size">
      <div class="price">{{ formatPrice(item.price) }}</div>
      <button @click="addToOrder(item)">Thêm</button>
    </div>
    
    <!-- Multi-size item -->
    <div v-else class="variants">
      <div v-for="variant in item.variants" :key="variant.id" class="variant-option">
        <div class="variant-info">
          <span class="variant-name">{{ variant.name }}</span>
          <span class="variant-price">{{ formatPrice(variant.price) }}</span>
        </div>
        <button @click="addToOrder(item, variant)" class="add-btn">
          Thêm
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
const addToOrder = (item, variant = null) => {
  orderStore.addItem(item, variant)
  alert(`Đã thêm ${item.name}${variant ? ` (${variant.name})` : ''}`)
}
</script>
```

**Tasks**:
- [ ] Update display logic
- [ ] Add variant buttons
- [ ] Style variants section
- [ ] Test responsive design

**Estimated Time**: 4 hours

#### 5.2 Update MenuView Create/Edit Form
**File**: `frontend/src/views/MenuView.vue`

**Changes**:
```vue
<template>
  <form @submit.prevent="saveItem">
    <input v-model="form.name" placeholder="Tên món" required />
    <select v-model="form.category" required>...</select>
    <textarea v-model="form.description">...</textarea>
    
    <!-- Toggle variants -->
    <div class="form-group">
      <label>
        <input type="checkbox" v-model="form.has_variants" />
        Món có nhiều size
      </label>
    </div>
    
    <!-- Single-size fields -->
    <div v-if="!form.has_variants" class="single-size-form">
      <div class="form-group">
        <label>Giá (VNĐ) *</label>
        <input v-model.number="form.price" type="number" required />
      </div>
      
      <div class="form-group">
        <label>Nguyên liệu</label>
        <!-- Existing ingredient selector -->
      </div>
    </div>
    
    <!-- Multi-size fields -->
    <div v-else class="variants-form">
      <div v-for="(variant, index) in form.variants" :key="index" class="variant-form-item">
        <div class="variant-header">
          <h4>Size {{ index + 1 }}</h4>
          <button type="button" @click="removeVariant(index)" class="btn-remove">
            🗑️
          </button>
        </div>
        
        <div class="form-group">
          <label>ID (M, L, XL) *</label>
          <input v-model="variant.id" required placeholder="M" />
        </div>
        
        <div class="form-group">
          <label>Tên hiển thị *</label>
          <input v-model="variant.name" required placeholder="Size M" />
        </div>
        
        <div class="form-group">
          <label>Giá (VNĐ) *</label>
          <input v-model.number="variant.price" type="number" required />
        </div>
        
        <div class="form-group">
          <label>
            <input type="checkbox" v-model="variant.is_default" />
            Mặc định
          </label>
        </div>
        
        <div class="form-group">
          <label>Nguyên liệu</label>
          <!-- Ingredient selector for this variant -->
          <button type="button" @click="showIngredientSelector(index)">
            Chọn nguyên liệu
          </button>
          <div v-if="variant.ingredients.length > 0">
            <div v-for="ing in variant.ingredients" :key="ing.name">
              {{ ing.name }}: {{ ing.quantity }} {{ ing.unit }}
            </div>
          </div>
        </div>
      </div>
      
      <button type="button" @click="addVariant" class="btn-add-variant">
        + Thêm size
      </button>
    </div>
    
    <div class="form-actions">
      <button type="button" @click="cancelEdit">Hủy</button>
      <button type="submit">Lưu</button>
    </div>
  </form>
</template>

<script setup>
const form = ref({
  name: '',
  category: '',
  description: '',
  has_variants: false,
  price: 0,
  ingredients: [],
  variants: []
})

const addVariant = () => {
  form.value.variants.push({
    id: '',
    name: '',
    price: 0,
    ingredients: [],
    available: true,
    is_default: form.value.variants.length === 0 // First is default
  })
}

const removeVariant = (index) => {
  form.value.variants.splice(index, 1)
  // Ensure at least one default
  if (form.value.variants.length > 0 && !form.value.variants.some(v => v.is_default)) {
    form.value.variants[0].is_default = true
  }
}

const saveItem = async () => {
  // Validation
  if (form.value.has_variants) {
    if (form.value.variants.length === 0) {
      alert('Phải có ít nhất 1 size')
      return
    }
    
    const defaultCount = form.value.variants.filter(v => v.is_default).length
    if (defaultCount !== 1) {
      alert('Phải có đúng 1 size mặc định')
      return
    }
  }
  
  // Save
  if (editingItem.value) {
    await menuStore.updateMenuItem(editingItem.value.id, form.value)
  } else {
    await menuStore.createMenuItem(form.value)
  }
  
  cancelEdit()
}
</script>
```

**Tasks**:
- [ ] Add has_variants checkbox
- [ ] Create variants form section
- [ ] Add/remove variant buttons
- [ ] Ingredient selector per variant
- [ ] Validation logic
- [ ] Style the form
- [ ] Test on mobile

**Estimated Time**: 8 hours

#### 5.3 Update Order Views
**File**: `frontend/src/views/WaiterView.vue`, `CashierView.vue`

**Changes**:
```vue
<template>
  <!-- Order item display -->
  <div v-for="item in order.items" :key="item.id">
    <div class="item-name">
      {{ item.name }}
      <span v-if="item.variant_name" class="variant-badge">
        ({{ item.variant_name }})
      </span>
    </div>
    <div class="item-price">{{ formatPrice(item.price) }}</div>
    <div class="item-quantity">x{{ item.quantity }}</div>
  </div>
</template>
```

**Tasks**:
- [ ] Display variant name in orders
- [ ] Update receipt display
- [ ] Test order flow

**Estimated Time**: 2 hours

---

### Phase 6: Testing (2 days)

#### 6.1 Backend Unit Tests
**Files**: `*_test.go`

**Test Cases**:
```go
// Menu domain tests
- TestMenuItem_GetDefaultVariant
- TestMenuItem_GetVariantByID
- TestMenuItem_GetPrice_SingleSize
- TestMenuItem_GetPrice_WithVariant
- TestMenuItem_Validate_NoVariants
- TestMenuItem_Validate_NoDefault
- TestMenuItem_Validate_MultipleDefaults
- TestMenuItem_Validate_DuplicateVariantIDs

// Service tests
- TestMenuService_CreateMenuItem_SingleSize
- TestMenuService_CreateMenuItem_WithVariants
- TestMenuService_CreateMenuItem_InvalidVariants
- TestOrderService_CreateOrder_WithVariant
- TestOrderService_CreateOrder_MissingVariantID
- TestCostCalculator_CalculateVariantCosts
```

**Tasks**:
- [ ] Write unit tests
- [ ] Achieve 80%+ coverage
- [ ] Test edge cases

**Estimated Time**: 8 hours

#### 6.2 Frontend Unit Tests
**Files**: `*.spec.js`

**Test Cases**:
```javascript
// Component tests
- MenuView displays single-size items
- MenuView displays multi-size items
- MenuView form toggles variants
- MenuView form validates variants
- OrderStore adds item with variant
```

**Tasks**:
- [ ] Write component tests
- [ ] Test user interactions
- [ ] Test validation

**Estimated Time**: 4 hours

#### 6.3 Integration Tests
**Manual Testing**:

**Test Scenarios**:
1. Create single-size item → Display → Order
2. Create multi-size item → Display → Order each size
3. Edit item: single → multi-size
4. Edit item: multi → single-size
5. Delete item with variants
6. Order with variant → Receipt display
7. Cost calculation per variant

**Tasks**:
- [ ] Manual testing
- [ ] Fix bugs
- [ ] Document issues

**Estimated Time**: 4 hours

---

## 📊 Timeline Summary

| Phase | Duration | Tasks |
|-------|----------|-------|
| Phase 1: Backend Domain | 2 days | Domain models, DTOs, validation |
| Phase 2: Backend Service | 2 days | Service layer, business logic |
| Phase 3: Backend API | 1 day | API handlers, documentation |
| Phase 4: Frontend Data | 1 day | Stores, data management |
| Phase 5: Frontend UI | 3 days | Views, forms, styling |
| Phase 6: Testing | 2 days | Unit, integration tests |
| **Total** | **11 days** | **~2 weeks** |

---

## 🎯 Success Criteria

### Must Have
- [ ] Can create single-size items (backward compatible)
- [ ] Can create multi-size items with variants
- [ ] Can order items with variant selection
- [ ] Variants display correctly in UI
- [ ] Cost calculation works per variant
- [ ] All tests passing

### Nice to Have
- [ ] Quick size templates (M/L/XL preset)
- [ ] Bulk edit variants
- [ ] Variant analytics
- [ ] Mobile-optimized variant selector

---

## 🚀 Getting Started

### Step 1: Backend Foundation (Day 1-2)
```bash
# Create feature branch
git checkout -b feature/menu-size-variants

# Update domain models
# File: backend/domain/menu/menu.go
# Add MenuItemVariant struct
# Update MenuItem struct
# Add helper methods

# Run tests
go test ./backend/domain/menu/...
```

### Step 2: Service Layer (Day 3-4)
```bash
# Update services
# File: backend/application/services/menu.go
# File: backend/application/services/order_service.go
# File: backend/application/services/cost_calculator_service.go

# Run tests
go test ./backend/application/services/...
```

### Step 3: Frontend (Day 5-9)
```bash
# Update stores
# File: frontend/src/stores/menu.js
# File: frontend/src/stores/order.js

# Update views
# File: frontend/src/views/MenuView.vue
# File: frontend/src/views/WaiterView.vue

# Test locally
npm run dev
```

### Step 4: Testing (Day 10-11)
```bash
# Run all tests
go test ./...
npm run test

# Manual testing
# Create test data
# Test all flows
```

---

## 💡 Implementation Tips

### 1. Start Simple
```javascript
// First iteration: Basic variants
{
  has_variants: true,
  variants: [
    { id: "M", name: "Size M", price: 25000 },
    { id: "L", name: "Size L", price: 30000 }
  ]
}

// Later: Add ingredients per variant
// Later: Add cost tracking per variant
```

### 2. Use Defaults
```javascript
// Auto-set first variant as default
if (form.variants.length === 1) {
  form.variants[0].is_default = true
}
```

### 3. Progressive Disclosure
```vue
<!-- Hide complexity initially -->
<div v-if="!form.has_variants">
  <!-- Simple form -->
</div>

<div v-else>
  <!-- Show variants only when needed -->
</div>
```

### 4. Validation Early
```javascript
// Validate on input, not just on submit
watch(() => form.variants, (variants) => {
  // Check for duplicate IDs
  // Check for default count
  // Show warnings
}, { deep: true })
```

---

## 🎉 Benefits (No Migration Risk!)

### Immediate Benefits
- ✅ Clean implementation from start
- ✅ No legacy data to handle
- ✅ No backward compatibility concerns
- ✅ Can use best practices

### Long-term Benefits
- ✅ Scalable architecture
- ✅ Better UX from day 1
- ✅ Easier to maintain
- ✅ Ready for future features

---

## 📝 Next Steps

**Ready to start?**

1. Review this plan
2. Confirm approach
3. Start with Phase 1 (Backend Domain)
4. Iterate and test

**Questions to answer**:
- Có muốn tôi bắt đầu implement không?
- Bắt đầu từ backend hay frontend trước?
- Có cần adjust timeline không?
- Có feature nào cần prioritize không?

Bạn muốn tôi bắt đầu implement từ đâu?
