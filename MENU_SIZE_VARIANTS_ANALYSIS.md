# Menu Size Variants - Phân tích & Thiết kế

## 📋 Yêu cầu

**Use Case**: Một món có thể có nhiều kích cỡ khác nhau
- Ví dụ: Cà phê sữa đá có size M, L, XL
- Mỗi size có giá khác nhau
- Mỗi size có công thức (ingredients) khác nhau

## 🔍 Phân tích Hiện trạng

### Cấu trúc hiện tại

**Backend - MenuItem**:
```go
type MenuItem struct {
    ID          primitive.ObjectID
    Name        string      // "Cà phê sữa đá"
    Price       float64     // 25000
    Category    string
    Description string
    Ingredients []Ingredient
    Available   bool
    // ... cost tracking fields
}
```

**Vấn đề**:
- ❌ 1 MenuItem = 1 giá cố định
- ❌ Không có khái niệm "size" hoặc "variant"
- ❌ Muốn có nhiều size phải tạo nhiều MenuItem riêng biệt:
  - "Cà phê sữa đá - M"
  - "Cà phê sữa đá - L"
  - "Cà phê sữa đá - XL"

**Hậu quả**:
- 😞 Menu dài, khó quản lý
- 😞 Trùng lặp thông tin (name, category, description)
- 😞 Khó maintain (update description phải sửa 3 món)
- 😞 UI không thân thiện (khách phải chọn từ 3 món riêng biệt)

## 🎯 Giải pháp

### Option 1: Menu Item với Size Variants (Recommended)

**Concept**: 1 MenuItem có nhiều variants, mỗi variant là 1 size

**Data Structure**:
```go
// New: Size variant
type MenuItemVariant struct {
    ID          string       `bson:"id" json:"id"`           // "M", "L", "XL"
    Name        string       `bson:"name" json:"name"`       // "Size M", "Size L", "Size XL"
    Price       float64      `bson:"price" json:"price"`     // 25000, 30000, 35000
    Ingredients []Ingredient `bson:"ingredients" json:"ingredients"`
    Available   bool         `bson:"available" json:"available"`
    IsDefault   bool         `bson:"is_default" json:"is_default"` // Default size
}

// Updated: MenuItem with variants
type MenuItem struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name        string             `bson:"name" json:"name"`        // "Cà phê sữa đá"
    Category    string             `bson:"category" json:"category"`
    Description string             `bson:"description" json:"description"`
    Available   bool               `bson:"available" json:"available"`
    
    // NEW: Variants
    HasVariants bool               `bson:"has_variants" json:"has_variants"`
    Variants    []MenuItemVariant  `bson:"variants,omitempty" json:"variants,omitempty"`
    
    // For backward compatibility (single-size items)
    Price       float64      `bson:"price,omitempty" json:"price,omitempty"`
    Ingredients []Ingredient `bson:"ingredients,omitempty" json:"ingredients,omitempty"`
    
    // Cost tracking
    CurrentCost          float64    `bson:"current_cost" json:"current_cost"`
    CostLastCalculatedAt time.Time  `bson:"cost_last_calculated_at" json:"cost_last_calculated_at"`
    CostStatus           CostStatus `bson:"cost_status" json:"cost_status"`
    
    CreatedAt   time.Time `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
```

**Ưu điểm**:
- ✅ 1 MenuItem cho tất cả sizes
- ✅ Không trùng lặp thông tin
- ✅ Dễ maintain
- ✅ UI thân thiện (chọn món → chọn size)
- ✅ Backward compatible (món không có size vẫn hoạt động)

**Nhược điểm**:
- ⚠️ Phức tạp hơn về code
- ⚠️ Cần migration cho data cũ
- ⚠️ Order phải lưu thêm variant_id

### Option 2: Separate Menu Items (Current)

**Concept**: Mỗi size là 1 MenuItem riêng biệt

**Data Structure**: Giữ nguyên như hiện tại

**Ưu điểm**:
- ✅ Đơn giản, không cần thay đổi code
- ✅ Không cần migration
- ✅ Order structure không đổi

**Nhược điểm**:
- ❌ Menu dài, khó quản lý
- ❌ Trùng lặp thông tin
- ❌ Khó maintain
- ❌ UI không thân thiện

### Option 3: Menu Groups (Alternative)

**Concept**: Nhóm các MenuItem liên quan lại

**Data Structure**:
```go
type MenuGroup struct {
    ID          primitive.ObjectID
    Name        string      // "Cà phê sữa đá"
    Category    string
    Description string
    Items       []primitive.ObjectID // IDs of menu items
}

// MenuItem giữ nguyên
```

**Ưu điểm**:
- ✅ Không thay đổi MenuItem structure
- ✅ Linh hoạt (có thể group bất kỳ món nào)

**Nhược điểm**:
- ❌ Vẫn trùng lặp thông tin
- ❌ Phức tạp hơn (2 collections)
- ❌ Khó maintain consistency

## 🏆 Recommendation: Option 1 - Size Variants

### Lý do chọn Option 1

1. **User Experience tốt nhất**
   - Khách chọn món → chọn size (natural flow)
   - Menu gọn gàng, dễ browse

2. **Maintainability**
   - Update description/category 1 lần cho tất cả sizes
   - Thêm/xóa size dễ dàng

3. **Scalability**
   - Có thể mở rộng thêm attributes khác (toppings, customizations)
   - Phù hợp với F&B business model

4. **Data Integrity**
   - Không trùng lặp
   - Single source of truth

## 📐 Thiết kế Chi tiết

### 1. Backend Changes

#### Domain Model

**File**: `backend/domain/menu/menu.go`

```go
package menu

import (
	"time"
	"cafe-pos/backend/domain/ingredient"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MenuItemVariant represents a size/variant of a menu item
type MenuItemVariant struct {
	ID          string               `bson:"id" json:"id"`                     // "M", "L", "XL"
	Name        string               `bson:"name" json:"name"`                 // "Size M", "Size L"
	Price       float64              `bson:"price" json:"price"`
	Ingredients []Ingredient         `bson:"ingredients" json:"ingredients"`
	Available   bool                 `bson:"available" json:"available"`
	IsDefault   bool                 `bson:"is_default" json:"is_default"`
	
	// Cost tracking per variant
	CurrentCost          float64    `bson:"current_cost" json:"current_cost"`
	CostLastCalculatedAt time.Time  `bson:"cost_last_calculated_at" json:"cost_last_calculated_at"`
	CostStatus           CostStatus `bson:"cost_status" json:"cost_status"`
}

type MenuItem struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Category    string             `bson:"category" json:"category"`
	Description string             `bson:"description" json:"description"`
	Available   bool               `bson:"available" json:"available"`
	
	// Variants support
	HasVariants bool               `bson:"has_variants" json:"has_variants"`
	Variants    []MenuItemVariant  `bson:"variants,omitempty" json:"variants,omitempty"`
	
	// Backward compatibility (for single-size items)
	Price       float64      `bson:"price,omitempty" json:"price,omitempty"`
	Ingredients []Ingredient `bson:"ingredients,omitempty" json:"ingredients,omitempty"`
	CurrentCost float64      `bson:"current_cost,omitempty" json:"current_cost,omitempty"`
	CostLastCalculatedAt time.Time  `bson:"cost_last_calculated_at,omitempty" json:"cost_last_calculated_at,omitempty"`
	CostStatus  CostStatus   `bson:"cost_status,omitempty" json:"cost_status,omitempty"`
	
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

// Helper methods
func (m *MenuItem) GetDefaultVariant() *MenuItemVariant {
	if !m.HasVariants {
		return nil
	}
	for i := range m.Variants {
		if m.Variants[i].IsDefault {
			return &m.Variants[i]
		}
	}
	// Return first variant if no default
	if len(m.Variants) > 0 {
		return &m.Variants[0]
	}
	return nil
}

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

func (m *MenuItem) GetPrice(variantID string) float64 {
	if m.HasVariants {
		variant := m.GetVariantByID(variantID)
		if variant != nil {
			return variant.Price
		}
		// Fallback to default variant
		defaultVariant := m.GetDefaultVariant()
		if defaultVariant != nil {
			return defaultVariant.Price
		}
		return 0
	}
	return m.Price
}

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

#### API Requests

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
	
	// For single-size items
	Price       float64      `json:"price"`
	Ingredients []Ingredient `json:"ingredients"`
	
	// For multi-size items
	HasVariants *bool              `json:"has_variants"`
	Variants    []MenuItemVariant  `json:"variants"`
}
```

### 2. Order Changes

**File**: `backend/domain/order/order.go`

```go
type OrderItem struct {
	MenuItemID  primitive.ObjectID `bson:"menu_item_id" json:"menu_item_id"`
	VariantID   string             `bson:"variant_id,omitempty" json:"variant_id,omitempty"` // NEW
	Name        string             `bson:"name" json:"name"`
	VariantName string             `bson:"variant_name,omitempty" json:"variant_name,omitempty"` // NEW: "Size M"
	Price       float64            `bson:"price" json:"price"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	Note        string             `bson:"note,omitempty" json:"note,omitempty"`
	Subtotal    float64            `bson:"subtotal" json:"subtotal"`
}
```

**Display Name**:
- Single-size: "Cà phê sữa đá"
- Multi-size: "Cà phê sữa đá (Size M)"

### 3. Frontend Changes

#### MenuView.vue

**Display**:
```vue
<template>
  <div v-for="item in menuItems" :key="item.id">
    <div class="menu-item">
      <h3>{{ item.name }}</h3>
      <p>{{ item.description }}</p>
      
      <!-- Single-size item -->
      <div v-if="!item.has_variants">
        <div class="price">{{ formatPrice(item.price) }}</div>
        <button @click="addToOrder(item)">Thêm</button>
      </div>
      
      <!-- Multi-size item -->
      <div v-else class="variants">
        <div v-for="variant in item.variants" :key="variant.id" class="variant">
          <div class="variant-info">
            <span class="variant-name">{{ variant.name }}</span>
            <span class="variant-price">{{ formatPrice(variant.price) }}</span>
          </div>
          <button @click="addToOrder(item, variant)">Thêm</button>
        </div>
      </div>
    </div>
  </div>
</template>
```

**Create/Edit Form**:
```vue
<template>
  <form @submit.prevent="saveMenuItem">
    <input v-model="form.name" placeholder="Tên món" />
    <select v-model="form.category">...</select>
    <textarea v-model="form.description">...</textarea>
    
    <!-- Toggle variants -->
    <label>
      <input type="checkbox" v-model="form.has_variants" />
      Món có nhiều size
    </label>
    
    <!-- Single-size fields -->
    <div v-if="!form.has_variants">
      <input v-model.number="form.price" placeholder="Giá" />
      <!-- Ingredients selector -->
    </div>
    
    <!-- Multi-size fields -->
    <div v-else>
      <div v-for="(variant, index) in form.variants" :key="index" class="variant-form">
        <input v-model="variant.id" placeholder="ID (M, L, XL)" />
        <input v-model="variant.name" placeholder="Tên (Size M)" />
        <input v-model.number="variant.price" placeholder="Giá" />
        <label>
          <input type="checkbox" v-model="variant.is_default" />
          Mặc định
        </label>
        <!-- Ingredients selector for this variant -->
        <button @click="removeVariant(index)">Xóa</button>
      </div>
      <button @click="addVariant">+ Thêm size</button>
    </div>
    
    <button type="submit">Lưu</button>
  </form>
</template>
```

### 4. Migration Strategy

#### Step 1: Update Schema (Backward Compatible)

```javascript
// MongoDB migration
db.menu_items.updateMany(
  { has_variants: { $exists: false } },
  { 
    $set: { 
      has_variants: false 
    } 
  }
)
```

#### Step 2: Convert Existing Multi-Size Items (Optional)

```javascript
// Example: Convert "Cà phê sữa đá - M/L/XL" to variants
// Manual process or script
```

#### Step 3: Deploy Backend

- New code supports both single-size and multi-size
- Old data still works (has_variants = false)

#### Step 4: Deploy Frontend

- UI shows variants if has_variants = true
- Falls back to single-size display if false

## 📊 Ví dụ Thực tế

### Ví dụ 1: Cà phê sữa đá (Multi-size)

**Database**:
```json
{
  "_id": "...",
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
      "is_default": true
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
      "is_default": false
    },
    {
      "id": "XL",
      "name": "Size XL",
      "price": 35000,
      "ingredients": [
        { "name": "Cà phê", "quantity": 40, "unit": "g" },
        { "name": "Sữa đặc", "quantity": 60, "unit": "ml" }
      ],
      "available": true,
      "is_default": false
    }
  ]
}
```

**UI Display**:
```
┌─────────────────────────────────────┐
│ ☕ Cà phê sữa đá                    │
│ Cà phê phin truyền thống với sữa đá │
│                                     │
│ ┌─────────────────────────────────┐ │
│ │ Size M          25,000₫  [Thêm] │ │
│ └─────────────────────────────────┘ │
│ ┌─────────────────────────────────┐ │
│ │ Size L          30,000₫  [Thêm] │ │
│ └─────────────────────────────────┘ │
│ ┌─────────────────────────────────┐ │
│ │ Size XL         35,000₫  [Thêm] │ │
│ └─────────────────────────────────┘ │
└─────────────────────────────────────┘
```

**Order**:
```json
{
  "items": [
    {
      "menu_item_id": "...",
      "variant_id": "L",
      "name": "Cà phê sữa đá",
      "variant_name": "Size L",
      "price": 30000,
      "quantity": 2,
      "subtotal": 60000
    }
  ]
}
```

**Receipt Display**:
```
Cà phê sữa đá (Size L) x2    60,000₫
```

### Ví dụ 2: Bánh mì (Single-size)

**Database**:
```json
{
  "_id": "...",
  "name": "Bánh mì thịt",
  "category": "Món ăn",
  "description": "Bánh mì Việt Nam truyền thống",
  "available": true,
  "has_variants": false,
  "price": 20000,
  "ingredients": [
    { "name": "Bánh mì", "quantity": 1, "unit": "cái" },
    { "name": "Thịt", "quantity": 50, "unit": "g" }
  ]
}
```

**UI Display**:
```
┌─────────────────────────────────────┐
│ 🥖 Bánh mì thịt                     │
│ Bánh mì Việt Nam truyền thống       │
│                                     │
│ 20,000₫                    [Thêm]  │
└─────────────────────────────────────┘
```

## 🎯 Implementation Plan

### Phase 1: Backend Foundation (Week 1)
- [ ] Update domain models (MenuItem, MenuItemVariant)
- [ ] Update repository methods
- [ ] Update service layer
- [ ] Update API handlers
- [ ] Write unit tests

### Phase 2: Order Integration (Week 1)
- [ ] Update OrderItem to support variant_id
- [ ] Update order service to handle variants
- [ ] Update cost calculation for variants
- [ ] Write integration tests

### Phase 3: Frontend UI (Week 2)
- [ ] Update MenuView to display variants
- [ ] Update create/edit form for variants
- [ ] Update order creation to select variants
- [ ] Update order display to show variant names

### Phase 4: Migration & Testing (Week 2)
- [ ] Create migration scripts
- [ ] Test with real data
- [ ] User acceptance testing
- [ ] Performance testing

### Phase 5: Deployment (Week 3)
- [ ] Deploy to staging
- [ ] Final testing
- [ ] Deploy to production
- [ ] Monitor and fix issues

## 🚨 Risks & Mitigation

### Risk 1: Breaking Changes
**Mitigation**: Backward compatibility - old data still works

### Risk 2: Complex UI
**Mitigation**: Progressive disclosure - simple items stay simple

### Risk 3: Performance
**Mitigation**: Index on has_variants, lazy load variants

### Risk 4: Data Migration
**Mitigation**: Optional migration, can keep old structure

## 📈 Benefits

### Business
- ✅ Better menu organization
- ✅ Easier to add new sizes
- ✅ Better customer experience
- ✅ More accurate inventory tracking per size

### Technical
- ✅ Less data duplication
- ✅ Easier maintenance
- ✅ Scalable for future features (toppings, customizations)
- ✅ Better data integrity

### User Experience
- ✅ Cleaner menu display
- ✅ Natural selection flow (item → size)
- ✅ Clear pricing per size
- ✅ Better mobile experience

## 🎉 Conclusion

**Recommendation**: Implement Option 1 - Size Variants

**Timeline**: 3 weeks

**Effort**: Medium

**Impact**: High

**Priority**: Medium (nice-to-have, not critical)

Bạn có muốn tôi bắt đầu implement không?
