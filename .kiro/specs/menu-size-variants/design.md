# Menu Size Variants - Design Document

## 1. Architecture Overview

### 1.1 System Components

```
┌─────────────────────────────────────────────────────────────┐
│                         Frontend                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  MenuView    │  │  WaiterView  │  │ CashierView  │      │
│  │  (Manager)   │  │  (Ordering)  │  │  (Payment)   │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                  │                  │              │
│  ┌──────▼──────────────────▼──────────────────▼───────┐    │
│  │              Stores (Pinia)                         │    │
│  │  - menuStore    - orderStore                        │    │
│  └──────┬──────────────────────────────────────────────┘    │
└─────────┼─────────────────────────────────────────────────┘
          │ HTTP/REST
┌─────────▼─────────────────────────────────────────────────┐
│                      Backend API                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │ MenuHandler  │  │ OrderHandler │  │ CostHandler  │     │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘     │
│         │                  │                  │             │
│  ┌──────▼──────────────────▼──────────────────▼───────┐   │
│  │              Service Layer                          │   │
│  │  - MenuService  - OrderService  - CostCalculator   │   │
│  └──────┬──────────────────────────────────────────────┘   │
│         │                                                   │
│  ┌──────▼──────────────────────────────────────────────┐   │
│  │              Domain Layer                            │   │
│  │  - MenuItem  - MenuItemVariant  - OrderItem         │   │
│  └──────┬──────────────────────────────────────────────┘   │
│         │                                                   │
│  ┌──────▼──────────────────────────────────────────────┐   │
│  │           Repository Layer                           │   │
│  │  - MenuRepository  - OrderRepository                 │   │
│  └──────┬──────────────────────────────────────────────┘   │
└─────────┼─────────────────────────────────────────────────┘
          │
┌─────────▼─────────────────────────────────────────────────┐
│                      MongoDB                                │
│  - menu_items collection                                    │
│  - orders collection                                        │
└─────────────────────────────────────────────────────────────┘
```

## 2. Data Model Design

### 2.1 MenuItem Schema

```go
type MenuItem struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name        string             `bson:"name" json:"name"`
    Category    string             `bson:"category" json:"category"`
    Description string             `bson:"description" json:"description"`
    Available   bool               `bson:"available" json:"available"`
    
    // Variants support
    HasVariants bool               `bson:"has_variants" json:"has_variants"`
    Variants    []MenuItemVariant  `bson:"variants,omitempty" json:"variants,omitempty"`
    
    // Backward compatibility (single-size items)
    Price       float64      `bson:"price,omitempty" json:"price,omitempty"`
    Ingredients []Ingredient `bson:"ingredients,omitempty" json:"ingredients,omitempty"`
    CurrentCost float64      `bson:"current_cost,omitempty" json:"current_cost,omitempty"`
    CostLastCalculatedAt time.Time  `bson:"cost_last_calculated_at,omitempty" json:"cost_last_calculated_at,omitempty"`
    CostStatus  CostStatus   `bson:"cost_status,omitempty" json:"cost_status,omitempty"`
    
    CreatedAt   time.Time `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

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
```

### 2.2 OrderItem Schema

```go
type OrderItem struct {
    MenuItemID  primitive.ObjectID `bson:"menu_item_id" json:"menu_item_id"`
    VariantID   string             `bson:"variant_id,omitempty" json:"variant_id,omitempty"`
    Name        string             `bson:"name" json:"name"`
    VariantName string             `bson:"variant_name,omitempty" json:"variant_name,omitempty"`
    Price       float64            `bson:"price" json:"price"`
    Quantity    int                `bson:"quantity" json:"quantity"`
    Note        string             `bson:"note,omitempty" json:"note,omitempty"`
    Subtotal    float64            `bson:"subtotal" json:"subtotal"`
}
```

### 2.3 Database Indexes

```javascript
// menu_items collection
db.menu_items.createIndex({ "name": 1 })
db.menu_items.createIndex({ "category": 1 })
db.menu_items.createIndex({ "has_variants": 1 })
db.menu_items.createIndex({ "variants.id": 1 })
db.menu_items.createIndex({ "available": 1 })

// Compound indexes
db.menu_items.createIndex({ "category": 1, "available": 1 })
db.menu_items.createIndex({ "has_variants": 1, "price": 1 })
```

### 2.4 Example Documents

**Single-Size Item**:
```json
{
  "_id": "507f1f77bcf86cd799439011",
  "name": "Bánh mì thịt",
  "category": "Món ăn",
  "description": "Bánh mì Việt Nam truyền thống",
  "available": true,
  "has_variants": false,
  "price": 20000,
  "ingredients": [
    { "name": "Bánh mì", "quantity": 1, "unit": "cái" },
    { "name": "Thịt", "quantity": 50, "unit": "g" }
  ],
  "current_cost": 12000,
  "cost_status": "FINAL",
  "cost_last_calculated_at": "2026-02-13T10:00:00Z",
  "created_at": "2026-02-13T09:00:00Z",
  "updated_at": "2026-02-13T10:00:00Z"
}
```

**Multi-Size Item**:
```json
{
  "_id": "507f1f77bcf86cd799439012",
  "name": "Cà phê sữa đá",
  "category": "Cà phê",
  "description": "Cà phê phin truyền thống với sữa đá",
  "available": true,
  "has_variants": true,
  "variants": [
    {
      "id": "M",
      "name": "Size M",
      "price": 25000,
      "ingredients": [
        { "name": "Cà phê", "quantity": 20, "unit": "g" },
        { "name": "Sữa đặc", "quantity": 30, "unit": "ml" }
      ],
      "available": true,
      "is_default": true,
      "current_cost": 15000,
      "cost_status": "FINAL",
      "cost_last_calculated_at": "2026-02-13T10:00:00Z"
    },
    {
      "id": "L",
      "name": "Size L",
      "price": 30000,
      "ingredients": [
        { "name": "Cà phê", "quantity": 30, "unit": "g" },
        { "name": "Sữa đặc", "quantity": 45, "unit": "ml" }
      ],
      "available": true,
      "is_default": false,
      "current_cost": 22000,
      "cost_status": "FINAL",
      "cost_last_calculated_at": "2026-02-13T10:00:00Z"
    }
  ],
  "created_at": "2026-02-13T09:00:00Z",
  "updated_at": "2026-02-13T10:00:00Z"
}
```

## 3. API Design

### 3.1 Menu Endpoints

**Create Menu Item**
```
POST /api/menu
Content-Type: application/json

// Single-size
{
  "name": "Bánh mì thịt",
  "category": "Món ăn",
  "description": "...",
  "has_variants": false,
  "price": 20000,
  "ingredients": [...]
}

// Multi-size
{
  "name": "Cà phê sữa đá",
  "category": "Cà phê",
  "description": "...",
  "has_variants": true,
  "variants": [
    {
      "id": "M",
      "name": "Size M",
      "price": 25000,
      "ingredients": [...],
      "is_default": true
    },
    {
      "id": "L",
      "name": "Size L",
      "price": 30000,
      "ingredients": [...],
      "is_default": false
    }
  ]
}

Response: 201 Created
{
  "id": "...",
  "name": "...",
  ...
}
```

**Get All Menu Items**
```
GET /api/menu

Response: 200 OK
{
  "data": [
    {
      "id": "...",
      "name": "Bánh mì thịt",
      "has_variants": false,
      "price": 20000,
      ...
    },
    {
      "id": "...",
      "name": "Cà phê sữa đá",
      "has_variants": true,
      "variants": [...]
    }
  ]
}
```

**Update Menu Item**
```
PUT /api/menu/:id
Content-Type: application/json

{
  "name": "...",
  "has_variants": true,
  "variants": [...]
}

Response: 200 OK
```

### 3.2 Order Endpoints

**Create Order**
```
POST /api/orders
Content-Type: application/json

{
  "customer_name": "Khách 1",
  "items": [
    {
      "menu_item_id": "...",
      "quantity": 1
      // No variant_id for single-size
    },
    {
      "menu_item_id": "...",
      "variant_id": "L",  // Required for multi-size
      "quantity": 2
    }
  ]
}

Response: 201 Created
{
  "id": "...",
  "items": [
    {
      "menu_item_id": "...",
      "name": "Bánh mì thịt",
      "price": 20000,
      "quantity": 1,
      "subtotal": 20000
    },
    {
      "menu_item_id": "...",
      "variant_id": "L",
      "name": "Cà phê sữa đá",
      "variant_name": "Size L",
      "price": 30000,
      "quantity": 2,
      "subtotal": 60000
    }
  ],
  "total": 80000
}
```

## 4. Business Logic Design

### 4.1 MenuItem Helper Methods

```go
// Get default variant
func (m *MenuItem) GetDefaultVariant() *MenuItemVariant {
    if !m.HasVariants {
        return nil
    }
    for i := range m.Variants {
        if m.Variants[i].IsDefault {
            return &m.Variants[i]
        }
    }
    // Fallback to first variant
    if len(m.Variants) > 0 {
        return &m.Variants[0]
    }
    return nil
}

// Get variant by ID
func (m *MenuItem) GetVariantByID(variantID string) *MenuItemVariant {
    if !m.HasVariants {
        return nil
    }
    for i := range m.Variants {
        if m.Variants[i].ID == variantID {
            return &m.Variants[i]
        }
    }
    return nil
}

// Get price (with optional variant)
func (m *MenuItem) GetPrice(variantID string) float64 {
    if m.HasVariants {
        if variantID != "" {
            variant := m.GetVariantByID(variantID)
            if variant != nil {
                return variant.Price
            }
        }
        // Fallback to default
        defaultVariant := m.GetDefaultVariant()
        if defaultVariant != nil {
            return defaultVariant.Price
        }
        return 0
    }
    return m.Price
}

// Get ingredients (with optional variant)
func (m *MenuItem) GetIngredients(variantID string) []Ingredient {
    if m.HasVariants {
        variant := m.GetVariantByID(variantID)
        if variant != nil {
            return variant.Ingredients
        }
        return nil
    }
    return m.Ingredients
}
```

### 4.2 Validation Logic

```go
func (m *MenuItem) Validate() error {
    // Basic validation
    if m.Name == "" {
        return errors.New("name is required")
    }
    if m.Category == "" {
        return errors.New("category is required")
    }
    
    if m.HasVariants {
        // Variants validation
        if len(m.Variants) == 0 {
            return errors.New("variants required when has_variants=true")
        }
        
        // Must have exactly 1 default
        defaultCount := 0
        for _, v := range m.Variants {
            if v.IsDefault {
                defaultCount++
            }
        }
        if defaultCount == 0 {
            return errors.New("must have at least one default variant")
        }
        if defaultCount > 1 {
            return errors.New("must have exactly one default variant")
        }
        
        // Variant IDs must be unique
        ids := make(map[string]bool)
        for _, v := range m.Variants {
            if v.ID == "" {
                return errors.New("variant ID is required")
            }
            if ids[v.ID] {
                return fmt.Errorf("duplicate variant ID: %s", v.ID)
            }
            ids[v.ID] = true
            
            // Variant must have valid price
            if v.Price <= 0 {
                return fmt.Errorf("variant %s price must be > 0", v.ID)
            }
        }
    } else {
        // Single-size validation
        if m.Price <= 0 {
            return errors.New("price must be > 0 for single-size item")
        }
    }
    
    return nil
}
```

### 4.3 Order Creation Logic

```go
func (s *OrderService) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.Order, error) {
    // Validate and populate order items
    for i, item := range req.Items {
        // Get menu item
        menuItem, err := s.menuRepo.FindByID(ctx, item.MenuItemID)
        if err != nil {
            return nil, fmt.Errorf("menu item not found: %w", err)
        }
        
        // Check availability
        if !menuItem.Available {
            return nil, fmt.Errorf("menu item %s is not available", menuItem.Name)
        }
        
        // Get price and validate variant
        var price float64
        var variantName string
        
        if menuItem.HasVariants {
            // Multi-size item - variant_id required
            if item.VariantID == "" {
                return nil, fmt.Errorf("variant_id required for %s", menuItem.Name)
            }
            
            variant := menuItem.GetVariantByID(item.VariantID)
            if variant == nil {
                return nil, fmt.Errorf("invalid variant_id %s for %s", item.VariantID, menuItem.Name)
            }
            
            if !variant.Available {
                return nil, fmt.Errorf("variant %s is not available", variant.Name)
            }
            
            price = variant.Price
            variantName = variant.Name
        } else {
            // Single-size item
            price = menuItem.Price
        }
        
        // Populate order item
        req.Items[i].Name = menuItem.Name
        req.Items[i].VariantName = variantName
        req.Items[i].Price = price
        req.Items[i].Subtotal = price * float64(item.Quantity)
    }
    
    // Create order
    newOrder := &order.Order{
        Items: req.Items,
        // ... rest of order creation
    }
    newOrder.CalculateTotal()
    
    err := s.orderRepo.Create(ctx, newOrder)
    if err != nil {
        return nil, err
    }
    
    return newOrder, nil
}
```

### 4.4 Cost Calculation Logic

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

func (s *CostCalculatorService) calculateIngredientsCost(ctx context.Context, ingredients []menu.Ingredient) (float64, menu.CostStatus) {
    totalCost := 0.0
    status := menu.CostStatusFinal
    
    for _, ing := range ingredients {
        // Get ingredient from database
        dbIng, err := s.ingredientRepo.FindByName(ctx, ing.Name)
        if err != nil || dbIng.CostPerUnit <= 0 {
            status = menu.CostStatusIncomplete
            continue
        }
        
        // Calculate cost with conversion rate
        conversionRate := ingredient.GetConversionRate(dbIng.Unit, ing.Unit)
        wastageMultiplier := 1.0 + (dbIng.WastagePercentage / 100.0)
        
        cost := ing.Quantity * dbIng.CostPerUnit * conversionRate * wastageMultiplier
        totalCost += cost
    }
    
    return totalCost, status
}
```

## 5. Frontend Design

### 5.1 Component Structure

```
MenuView.vue
├── MenuItemCard (single-size)
│   ├── Name, Description
│   ├── Price
│   └── Add Button
│
└── MenuItemCard (multi-size)
    ├── Name, Description
    ├── Variants List
    │   ├── Variant Option 1 (Name, Price, Add Button)
    │   ├── Variant Option 2
    │   └── Variant Option 3
    └── ...

MenuForm (Create/Edit)
├── Basic Info (Name, Category, Description)
├── Has Variants Checkbox
├── Single-Size Section (if !has_variants)
│   ├── Price Input
│   └── Ingredients Selector
└── Variants Section (if has_variants)
    ├── Variant Form 1
    │   ├── ID, Name, Price
    │   ├── Is Default Checkbox
    │   ├── Ingredients Selector
    │   └── Remove Button
    ├── Variant Form 2
    ├── ...
    └── Add Variant Button
```

### 5.2 State Management

```javascript
// menuStore.js
export const useMenuStore = defineStore('menu', {
  state: () => ({
    items: [],
    loading: false,
    error: null
  }),
  
  actions: {
    async createMenuItem(data) {
      this.loading = true
      try {
        const response = await menuService.create(data)
        this.items.push(response.data)
        return true
      } catch (error) {
        this.error = error.message
        return false
      } finally {
        this.loading = false
      }
    },
    
    // ... other actions
  }
})

// orderStore.js
export const useOrderStore = defineStore('order', {
  state: () => ({
    currentOrder: {
      items: []
    }
  }),
  
  actions: {
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
      
      this.currentOrder.items.push(item)
    }
  }
})
```

### 5.3 UI/UX Flow

**Create Menu Item Flow**:
```
1. Click "Thêm món"
2. Enter name, category, description
3. Toggle "Món có nhiều size"
   
   If NO (single-size):
   4a. Enter price
   5a. Select ingredients
   6a. Click "Lưu"
   
   If YES (multi-size):
   4b. Click "Thêm size"
   5b. Enter variant ID, name, price
   6b. Select ingredients for variant
   7b. Mark one as default
   8b. Repeat for more sizes
   9b. Click "Lưu"
```

**Order Item Flow**:
```
Single-size:
1. Tap item → Added to order

Multi-size:
1. Tap item → Variant options appear
2. Tap size → Added to order with variant
```

## 6. Testing Strategy

### 6.1 Unit Tests

**Backend**:
- MenuItem.Validate() - all validation rules
- MenuItem.GetPrice() - with/without variant
- MenuItem.GetVariantByID() - found/not found
- OrderService.CreateOrder() - with/without variant
- CostCalculator - per variant calculation

**Frontend**:
- MenuForm validation
- Variant add/remove logic
- Order item display with variants

### 6.2 Integration Tests

- Create single-size item → Display → Order
- Create multi-size item → Display → Order each variant
- Edit item: single → multi-size
- Cost calculation for variants
- Order with invalid variant_id (error)

### 6.3 E2E Tests

- Complete flow: Create menu → Order → Payment → Receipt
- Variant selection in order flow
- Cost tracking per variant

## 7. Performance Considerations

### 7.1 Database Optimization

- Index on `has_variants` for filtering
- Index on `variants.id` for lookups
- Compound index on `category` + `available`
- Limit variants per item to 10

### 7.2 Frontend Optimization

- Lazy load variant details
- Cache menu items
- Debounce form inputs
- Optimize re-renders

### 7.3 API Optimization

- Return only necessary fields
- Paginate menu list if > 100 items
- Cache frequently accessed items
- Batch cost calculations

## 8. Security Considerations

- Validate all inputs server-side
- Prevent SQL injection (use parameterized queries)
- Sanitize user inputs
- Rate limit API endpoints
- Require authentication for menu management

## 9. Monitoring & Logging

### 9.1 Metrics to Track

- Menu item creation rate
- Variant usage distribution
- Order creation with variants
- Cost calculation performance
- API response times

### 9.2 Logging

- Log all menu item changes
- Log validation errors
- Log order creation with variants
- Log cost calculation results

## 10. Rollout Plan

### Phase 1: Backend (Days 1-5)
- Implement domain models
- Implement service layer
- Implement API endpoints
- Write unit tests

### Phase 2: Frontend (Days 6-9)
- Implement menu form
- Implement menu display
- Implement order flow
- Write component tests

### Phase 3: Testing (Days 10-11)
- Integration testing
- E2E testing
- Bug fixes
- Performance testing

### Phase 4: Deployment (Day 11)
- Deploy to staging
- Final testing
- Deploy to production
- Monitor

## 11. Success Criteria

- All acceptance criteria met
- Unit test coverage > 80%
- All integration tests passing
- Performance requirements met
- Zero critical bugs
- Positive user feedback
