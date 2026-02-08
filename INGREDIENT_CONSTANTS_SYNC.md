# Ingredient Constants Synchronization

## Overview
This document ensures frontend and backend constants are synchronized to prevent communication errors.

## Constants Mapping

### 1. Unit Types

**Backend** (`backend/domain/ingredient/ingredient.go`):
```go
type UnitType string

const (
    UnitKilogram   UnitType = "kg"
    UnitGram       UnitType = "g"
    UnitLiter      UnitType = "L"
    UnitMilliliter UnitType = "ml"
    UnitPiece      UnitType = "piece"
    UnitBox        UnitType = "box"
    UnitPack       UnitType = "pack"
)
```

**Frontend** (`frontend/src/constants/ingredient.js`):
```javascript
export const INGREDIENT_UNITS = {
  KILOGRAM: 'kg',
  GRAM: 'g',
  LITER: 'L',
  MILLILITER: 'ml',
  PIECE: 'piece',
  BOX: 'box',
  PACK: 'pack'
}
```

✅ **Status**: SYNCED

---

### 2. Stock Operation Types

**Backend** (`backend/domain/ingredient/ingredient.go`):
```go
type StockOperationType string

const (
    StockOperationIn     StockOperationType = "in"     // Stock IN (purchase)
    StockOperationOut    StockOperationType = "out"    // Stock OUT (usage/waste)
    StockOperationAdjust StockOperationType = "adjust" // Stock ADJUST (inventory correction)
)
```

**Frontend** (`frontend/src/constants/ingredient.js`):
```javascript
export const STOCK_OPERATIONS = {
  IN: 'in',         // Stock IN (purchase/receive)
  OUT: 'out',       // Stock OUT (usage/waste)
  ADJUST: 'adjust'  // Stock ADJUST (inventory correction)
}
```

✅ **Status**: SYNCED

**API Endpoints**:
- `POST /manager/ingredients/:id/stock-in` → Uses `StockOperationIn`
- `POST /manager/ingredients/:id/stock-out` → Uses `StockOperationOut`
- `POST /manager/ingredients/:id/stock-adjust` → Uses `StockOperationAdjust`

---

### 3. Transaction Types (History)

**Backend** (`backend/domain/ingredient/stock_history.go`):
```go
type TransactionType string

const (
    TransactionAdjustment TransactionType = "adjustment"
    TransactionOrder      TransactionType = "order"
    TransactionPurchase   TransactionType = "purchase"
    TransactionWaste      TransactionType = "waste"
)
```

**Frontend** (`frontend/src/constants/ingredient.js`):
```javascript
export const TRANSACTION_TYPES = {
  ADJUSTMENT: 'adjustment',
  ORDER: 'order',
  PURCHASE: 'purchase',
  WASTE: 'waste'
}
```

✅ **Status**: SYNCED

**Usage Mapping**:
- Stock IN → Creates `TransactionPurchase` history
- Stock OUT → Creates `TransactionWaste` history
- Stock ADJUST → Creates `TransactionAdjustment` history
- Order usage → Creates `TransactionOrder` history

---

### 4. Legacy Adjustment Types (Deprecated)

**Frontend Only** (`frontend/src/constants/ingredient.js`):
```javascript
// Legacy - for backward compatibility only
export const ADJUSTMENT_TYPES = {
  ADD: 'add',
  REMOVE: 'remove',
  ADJUST: 'adjust'
}
```

⚠️ **Status**: DEPRECATED - Use `STOCK_OPERATIONS` instead

**Migration Path**:
- `ADJUSTMENT_TYPES.ADD` → `STOCK_OPERATIONS.IN`
- `ADJUSTMENT_TYPES.REMOVE` → `STOCK_OPERATIONS.OUT`
- `ADJUSTMENT_TYPES.ADJUST` → `STOCK_OPERATIONS.ADJUST`

---

## Request/Response Structures

### Stock IN Request

**Frontend**:
```javascript
const data = {
  quantity: 10,           // float64, required, > 0
  cost_per_unit: 50000,   // float64, optional, >= 0 (0 = use current)
  reason: "Nhập kho"      // string, optional
}
await ingredientService.stockIn(ingredientId, data)
```

**Backend**:
```go
type StockInRequest struct {
    Quantity    float64 `json:"quantity" binding:"required,gt=0"`
    CostPerUnit float64 `json:"cost_per_unit" binding:"min=0"`
    Reason      string  `json:"reason"`
    UserID      string  `json:"user_id"`      // Auto-filled by middleware
    Username    string  `json:"username"`     // Auto-filled by middleware
}
```

✅ **Validation**:
- `quantity` must be > 0
- `cost_per_unit` must be >= 0
- `reason` is optional
- `user_id` and `username` are auto-filled from JWT token

---

### Stock OUT Request

**Frontend**:
```javascript
const data = {
  quantity: 5,                      // float64, required, > 0
  reason: "Sử dụng cho món ăn"      // string, required
}
await ingredientService.stockOut(ingredientId, data)
```

**Backend**:
```go
type StockOutRequest struct {
    Quantity float64 `json:"quantity" binding:"required,gt=0"`
    Reason   string  `json:"reason" binding:"required"`
    UserID   string  `json:"user_id"`      // Auto-filled
    Username string  `json:"username"`     // Auto-filled
}
```

✅ **Validation**:
- `quantity` must be > 0
- `reason` is required
- NO `cost_per_unit` field (price never changes on stock out)

---

### Stock ADJUST Request

**Frontend**:
```javascript
const data = {
  new_quantity: 12,       // float64, required, >= 0 (target quantity)
  cost_per_unit: 60000,   // float64, optional, >= 0 (only if increase with new price)
  reason: "Kiểm kê"       // string, required
}
await ingredientService.stockAdjust(ingredientId, data)
```

**Backend**:
```go
type StockAdjustRequest struct {
    NewQuantity float64 `json:"new_quantity" binding:"required,min=0"`
    CostPerUnit float64 `json:"cost_per_unit" binding:"min=0"`
    Reason      string  `json:"reason" binding:"required"`
    UserID      string  `json:"user_id"`      // Auto-filled
    Username    string  `json:"username"`     // Auto-filled
}
```

✅ **Validation**:
- `new_quantity` must be >= 0 (target quantity, not diff)
- `cost_per_unit` is optional (0 = keep current price)
- `reason` is required
- Backend calculates diff: `quantityDiff = new_quantity - current_quantity`

---

## Business Logic Rules

### Price Recalculation Matrix

| Operation | Quantity Change | Price Provided | Price Different | Recalculate? |
|-----------|----------------|----------------|-----------------|--------------|
| Stock IN | +10 | Yes (50k) | Yes | ✅ YES |
| Stock IN | +10 | Yes (40k) | No (same) | ❌ NO |
| Stock IN | +10 | No (0) | N/A | ❌ NO |
| Stock OUT | -5 | N/A | N/A | ❌ NO |
| Stock ADJUST | +2 | Yes (50k) | Yes | ✅ YES |
| Stock ADJUST | +2 | Yes (40k) | No (same) | ❌ NO |
| Stock ADJUST | +2 | No (0) | N/A | ❌ NO |
| Stock ADJUST | -3 | N/A | N/A | ❌ NO |

### Weighted Average Formula

```
new_price = (old_qty × old_price + new_qty × new_price) / total_qty
```

**Example**:
```
Current: 10 kg @ 40,000đ/kg
Purchase: 5 kg @ 50,000đ/kg

new_price = (10 × 40,000 + 5 × 50,000) / 15
          = (400,000 + 250,000) / 15
          = 43,333đ/kg
```

---

## Frontend Usage

### Import Constants

```javascript
import {
  INGREDIENT_UNITS,
  STOCK_OPERATIONS,
  TRANSACTION_TYPES,
  getStockStatus,
  getStockStatusClass,
  getStockStatusText,
  getAdjustmentTypeClass,
  getAdjustmentTypeText
} from '@/constants/ingredient'
```

### Use in Components

```javascript
// Stock IN
const handleStockIn = async () => {
  const data = {
    quantity: formData.quantity,
    cost_per_unit: formData.cost_per_unit || 0,
    reason: formData.reason || 'Nhập kho'
  }
  await ingredientStore.stockIn(ingredient.id, data)
}

// Stock OUT
const handleStockOut = async () => {
  const data = {
    quantity: formData.quantity,
    reason: formData.reason  // Required!
  }
  await ingredientStore.stockOut(ingredient.id, data)
}

// Stock ADJUST
const handleStockAdjust = async () => {
  const data = {
    new_quantity: formData.targetQuantity,  // Target, not diff!
    cost_per_unit: formData.cost_per_unit || 0,
    reason: formData.reason  // Required!
  }
  await ingredientStore.stockAdjust(ingredient.id, data)
}
```

### Display History

```javascript
// Get transaction type text
const typeText = getAdjustmentTypeText(history.type)
// Returns: "Nhập Hàng", "Xuất Hàng", "Điều Chỉnh", "Sử Dụng"

// Get transaction type class
const typeClass = getAdjustmentTypeClass(history.type)
// Returns: "bg-green-100 text-green-800", etc.
```

---

## Backend Usage

### Service Layer

```go
// Stock IN
func (s *IngredientService) StockIn(ctx context.Context, id primitive.ObjectID, req *ingredient.StockInRequest) (*ingredient.Ingredient, error) {
    // Validate
    if req.Quantity <= 0 {
        return nil, errors.New("quantity must be positive")
    }
    
    // Calculate weighted average if new price provided
    if req.CostPerUnit > 0 && req.CostPerUnit != item.CostPerUnit {
        // Weighted average logic
    }
    
    // Create history with TransactionPurchase
    history := &ingredient.StockHistory{
        Type: ingredient.TransactionPurchase,
        // ...
    }
}

// Stock OUT
func (s *IngredientService) StockOut(ctx context.Context, id primitive.ObjectID, req *ingredient.StockOutRequest) (*ingredient.Ingredient, error) {
    // Validate
    if req.Quantity <= 0 {
        return nil, errors.New("quantity must be positive")
    }
    if req.Reason == "" {
        return nil, errors.New("reason is required")
    }
    
    // Price NEVER changes
    // Create history with TransactionWaste
    history := &ingredient.StockHistory{
        Type: ingredient.TransactionWaste,
        Quantity: -req.Quantity,  // Negative for removal
        // ...
    }
}

// Stock ADJUST
func (s *IngredientService) StockAdjust(ctx context.Context, id primitive.ObjectID, req *ingredient.StockAdjustRequest) (*ingredient.Ingredient, error) {
    // Calculate diff
    quantityDiff := req.NewQuantity - item.Quantity
    
    // Only recalculate if increase with new price
    if quantityDiff > 0 && req.CostPerUnit > 0 && req.CostPerUnit != item.CostPerUnit {
        // Weighted average logic
    }
    
    // Create history with TransactionAdjustment
    history := &ingredient.StockHistory{
        Type: ingredient.TransactionAdjustment,
        Quantity: quantityDiff,  // Can be positive or negative
        // ...
    }
}
```

---

## Validation Rules

### Frontend Validation

```javascript
// Stock IN
if (!data.quantity || data.quantity <= 0) {
  alert('Số lượng phải lớn hơn 0')
  return
}

// Stock OUT
if (!data.quantity || data.quantity <= 0) {
  alert('Số lượng phải lớn hơn 0')
  return
}
if (!data.reason || data.reason.trim() === '') {
  alert('Vui lòng nhập lý do')
  return
}
if (data.quantity > currentQuantity) {
  alert('Số lượng xuất không được lớn hơn tồn kho')
  return
}

// Stock ADJUST
if (data.new_quantity < 0) {
  alert('Số lượng không được âm')
  return
}
if (!data.reason || data.reason.trim() === '') {
  alert('Vui lòng nhập lý do')
  return
}
```

### Backend Validation

```go
// Gin binding tags handle validation
type StockInRequest struct {
    Quantity    float64 `json:"quantity" binding:"required,gt=0"`     // Must be > 0
    CostPerUnit float64 `json:"cost_per_unit" binding:"min=0"`        // Must be >= 0
    Reason      string  `json:"reason"`                               // Optional
}

type StockOutRequest struct {
    Quantity float64 `json:"quantity" binding:"required,gt=0"`        // Must be > 0
    Reason   string  `json:"reason" binding:"required"`               // Required
}

type StockAdjustRequest struct {
    NewQuantity float64 `json:"new_quantity" binding:"required,min=0"` // Must be >= 0
    CostPerUnit float64 `json:"cost_per_unit" binding:"min=0"`         // Must be >= 0
    Reason      string  `json:"reason" binding:"required"`             // Required
}
```

---

## Error Handling

### Frontend

```javascript
try {
  await ingredientStore.stockIn(id, data)
  // Success
} catch (error) {
  const errorMsg = error.response?.data?.error || 'Lỗi không xác định'
  alert(errorMsg)
}
```

### Backend

```go
// Return clear error messages
if req.Quantity <= 0 {
    c.JSON(http.StatusBadRequest, gin.H{"error": "quantity must be positive"})
    return
}

if req.Reason == "" {
    c.JSON(http.StatusBadRequest, gin.H{"error": "reason is required"})
    return
}

if item.Quantity < req.Quantity {
    c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient stock"})
    return
}
```

---

## Testing Checklist

### Unit Constants Match
- [ ] `INGREDIENT_UNITS` matches `UnitType` constants
- [ ] `STOCK_OPERATIONS` matches `StockOperationType` constants
- [ ] `TRANSACTION_TYPES` matches `TransactionType` constants

### API Contracts Match
- [ ] Stock IN request structure matches
- [ ] Stock OUT request structure matches
- [ ] Stock ADJUST request structure matches
- [ ] Response structures match

### Business Logic Consistent
- [ ] Price recalculation rules match
- [ ] Weighted average formula matches
- [ ] History creation logic matches

### Validation Consistent
- [ ] Frontend validation matches backend
- [ ] Error messages are clear
- [ ] Edge cases handled

---

## Maintenance Guidelines

### When Adding New Constants

1. **Add to Backend First**:
   ```go
   // backend/domain/ingredient/ingredient.go
   const (
       UnitNewUnit UnitType = "new_unit"
   )
   ```

2. **Add to Frontend**:
   ```javascript
   // frontend/src/constants/ingredient.js
   export const INGREDIENT_UNITS = {
       NEW_UNIT: 'new_unit'
   }
   ```

3. **Update This Document**:
   - Add to constants mapping
   - Add usage examples
   - Update validation rules

4. **Test**:
   - Backend accepts new constant
   - Frontend sends new constant
   - Database stores correctly
   - UI displays correctly

### When Changing Constants

1. **Never Change Existing Values** (breaks backward compatibility)
2. **Add New Constants** instead
3. **Deprecate Old Constants** with migration path
4. **Update Documentation**

### Version Control

- Backend constants are source of truth
- Frontend must match backend exactly
- Document any deviations with clear reasons
- Keep this sync document updated

---

## Quick Reference

### Stock Operations Summary

| Operation | Endpoint | Quantity | Price | Reason | Recalculates Price? |
|-----------|----------|----------|-------|--------|---------------------|
| Stock IN | `/stock-in` | Positive | Optional | Optional | If new price provided |
| Stock OUT | `/stock-out` | Positive | N/A | Required | Never |
| Stock ADJUST | `/stock-adjust` | Target (>=0) | Optional | Required | If increase + new price |

### Transaction Types Summary

| Type | Created By | Quantity Sign | Price Changes? |
|------|-----------|---------------|----------------|
| purchase | Stock IN | Positive | Maybe |
| waste | Stock OUT | Negative | Never |
| adjustment | Stock ADJUST | Positive/Negative | Maybe |
| order | Order System | Negative | Never |

---

## Status: ✅ SYNCED

Last Updated: 2026-02-07
Last Verified: 2026-02-07

All constants are synchronized between frontend and backend.
