# Ingredient Stock Operations Refactor

## Overview
Refactored ingredient stock management to have separate, clear methods for different operations with proper constant-based communication between frontend and backend.

## Backend Changes

### 1. Domain Model Updates (`backend/domain/ingredient/ingredient.go`)

#### New Stock Operation Types
```go
type StockOperationType string

const (
    StockOperationIn     StockOperationType = "in"     // Stock IN (purchase)
    StockOperationOut    StockOperationType = "out"    // Stock OUT (usage/waste)
    StockOperationAdjust StockOperationType = "adjust" // Stock ADJUST (inventory correction)
)
```

#### New Request Types

**StockInRequest** - For purchasing/receiving stock:
```go
type StockInRequest struct {
    Quantity    float64 `json:"quantity" binding:"required,gt=0"`
    CostPerUnit float64 `json:"cost_per_unit" binding:"min=0"` // Optional: if 0, use current price
    Reason      string  `json:"reason"`
    UserID      string  `json:"user_id"`
    Username    string  `json:"username"`
}
```

**StockOutRequest** - For using/removing stock:
```go
type StockOutRequest struct {
    Quantity float64 `json:"quantity" binding:"required,gt=0"`
    Reason   string  `json:"reason" binding:"required"`
    UserID   string  `json:"user_id"`
    Username string  `json:"username"`
}
```

**StockAdjustRequest** - For inventory corrections:
```go
type StockAdjustRequest struct {
    NewQuantity float64 `json:"new_quantity" binding:"required,min=0"` // Target quantity
    CostPerUnit float64 `json:"cost_per_unit" binding:"min=0"`         // Optional: if increase due to purchase
    Reason      string  `json:"reason" binding:"required"`
    UserID      string  `json:"user_id"`
    Username    string  `json:"username"`
}
```

### 2. Service Layer (`backend/application/services/ingredient.go`)

#### StockIn Method
- Adds stock (purchase/receive)
- Calculates weighted average ONLY when new price is provided and different
- Creates purchase history record
- Tracks expense automatically

```go
func (s *IngredientService) StockIn(ctx context.Context, id primitive.ObjectID, req *ingredient.StockInRequest) (*ingredient.Ingredient, error)
```

**Logic:**
- `item.Quantity += req.Quantity`
- If `req.CostPerUnit > 0` AND `req.CostPerUnit != item.CostPerUnit`:
  - Calculate weighted average: `(old_qty * old_price + new_qty * new_price) / total_qty`
- Else: Keep current price
- Create history with type `TransactionPurchase`
- Track expense if price > 0

#### StockOut Method
- Removes stock (usage/waste)
- NEVER changes price
- Creates waste history record
- NO expense tracking

```go
func (s *IngredientService) StockOut(ctx context.Context, id primitive.ObjectID, req *ingredient.StockOutRequest) (*ingredient.Ingredient, error)
```

**Logic:**
- `item.Quantity -= req.Quantity`
- If `item.Quantity < 0`: Set to 0
- Price NEVER changes
- Create history with type `TransactionWaste` and negative quantity
- NO expense tracking

#### StockAdjust Method
- Sets stock to specific quantity (inventory correction)
- Calculates difference automatically
- Optionally recalculates price if quantity increased with new price

```go
func (s *IngredientService) StockAdjust(ctx context.Context, id primitive.ObjectID, req *ingredient.StockAdjustRequest) (*ingredient.Ingredient, error)
```

**Logic:**
- Calculate `quantityDiff = req.NewQuantity - beforeQty`
- `item.Quantity = req.NewQuantity`
- If `quantityDiff > 0` AND `req.CostPerUnit > 0` AND `req.CostPerUnit != item.CostPerUnit`:
  - Calculate weighted average for the increase
- Create history with type `TransactionAdjustment`
- Track expense if quantity increased with new price

### 3. HTTP Handler (`backend/interfaces/http/ingredient_handler.go`)

New endpoints:
- `POST /manager/ingredients/:id/stock-in` → `StockIn()`
- `POST /manager/ingredients/:id/stock-out` → `StockOut()`
- `POST /manager/ingredients/:id/stock-adjust` → `StockAdjust()`
- `POST /manager/ingredients/:id/adjust` → `AdjustStock()` (legacy)

### 4. Routes (`backend/main.go`)

```go
// Stock operations
manager.POST("/ingredients/:id/stock-in", ingredientHandler.StockIn)
manager.POST("/ingredients/:id/stock-out", ingredientHandler.StockOut)
manager.POST("/ingredients/:id/stock-adjust", ingredientHandler.StockAdjust)
manager.POST("/ingredients/:id/adjust", ingredientHandler.AdjustStock) // Legacy
```

## Frontend Changes

### 1. Constants (`frontend/src/constants/ingredient.js`)

#### New Stock Operation Constants
```javascript
export const STOCK_OPERATIONS = {
  IN: 'in',         // Stock IN (purchase/receive)
  OUT: 'out',       // Stock OUT (usage/waste)
  ADJUST: 'adjust'  // Stock ADJUST (inventory correction)
}

export const TRANSACTION_TYPES = {
  ADJUSTMENT: 'adjustment',
  ORDER: 'order',
  PURCHASE: 'purchase',
  WASTE: 'waste'
}
```

### 2. Service Layer (`frontend/src/services/ingredient.js`)

New methods:
```javascript
async stockIn(id, data) {
  // data: { quantity, cost_per_unit, reason }
  const response = await api.post(`/manager/ingredients/${id}/stock-in`, data)
  return response.data
}

async stockOut(id, data) {
  // data: { quantity, reason }
  const response = await api.post(`/manager/ingredients/${id}/stock-out`, data)
  return response.data
}

async stockAdjust(id, data) {
  // data: { new_quantity, cost_per_unit, reason }
  const response = await api.post(`/manager/ingredients/${id}/stock-adjust`, data)
  return response.data
}
```

### 3. Store Layer (`frontend/src/stores/ingredient.js`)

New actions:
```javascript
async stockIn(id, data) {
  this.error = null
  try {
    const updatedItem = await ingredientService.stockIn(id, data)
    // Update local state
    return true
  } catch (error) {
    this.error = error.response?.data?.error || 'Lỗi nhập kho'
    return false
  }
}

async stockOut(id, data) { /* similar */ }
async stockAdjust(id, data) { /* similar */ }
```

### 4. View Layer (`frontend/src/views/IngredientManagementView.vue`)

**Quick Stock IN:**
```javascript
const confirmQuickIn = async () => {
  const data = {
    quantity: quickInData.value.quantity,
    cost_per_unit: quickInData.value.cost_per_unit || 0,
    reason: 'Nhập kho'
  }
  await ingredientStore.stockIn(currentIngredient.value.id, data)
}
```

**Quick Stock OUT:**
```javascript
const confirmQuickOut = async () => {
  const data = {
    quantity: quickOutData.value.quantity,
    reason: quickOutData.value.reason
  }
  await ingredientStore.stockOut(currentIngredient.value.id, data)
}
```

**Stock Adjust:**
```javascript
const adjustStock = async () => {
  const data = {
    new_quantity: adjustData.value.quantity,
    cost_per_unit: adjustData.value.cost_per_unit || 0,
    reason: adjustData.value.reason
  }
  await ingredientStore.stockAdjust(currentIngredient.value.id, data)
}
```

## Key Improvements

### 1. Clear Separation of Concerns
- **Stock IN**: Purchase/receive inventory
- **Stock OUT**: Use/waste inventory
- **Stock ADJUST**: Correct inventory count

### 2. Simplified Logic
- Frontend doesn't calculate quantity differences
- Backend handles all business logic
- No more confusion about positive/negative quantities

### 3. Proper Price Handling
- **Stock IN**: Can update price (weighted average)
- **Stock OUT**: Never updates price
- **Stock ADJUST**: Can update price if increase due to purchase

### 4. Consistent Communication
- All constants defined in both frontend and backend
- Type-safe operations
- Clear API contracts

### 5. Better History Tracking
- Proper transaction types (purchase, waste, adjustment)
- Accurate price information
- Clear audit trail

## Migration Path

### Phase 1: Backend (Completed)
✅ Add new domain types
✅ Add new service methods
✅ Add new HTTP handlers
✅ Add new routes
✅ Keep legacy endpoint for compatibility

### Phase 2: Frontend (Next)
- [ ] Update view to use new methods
- [ ] Simplify quick actions
- [ ] Remove quantity calculation logic
- [ ] Test all operations

### Phase 3: Testing
- [ ] Test Stock IN with same price
- [ ] Test Stock IN with different price
- [ ] Test Stock OUT
- [ ] Test Stock ADJUST (increase)
- [ ] Test Stock ADJUST (decrease)
- [ ] Verify expense tracking
- [ ] Verify history records

### Phase 4: Cleanup (Future)
- [ ] Remove legacy adjustStock method
- [ ] Remove old adjustment type constants
- [ ] Update documentation

## API Examples

### Stock IN (Purchase)
```bash
POST /manager/ingredients/:id/stock-in
{
  "quantity": 10,
  "cost_per_unit": 50000,  // Optional: 0 = use current price
  "reason": "Nhập từ nhà cung cấp ABC"
}
```

### Stock OUT (Usage/Waste)
```bash
POST /manager/ingredients/:id/stock-out
{
  "quantity": 5,
  "reason": "Sử dụng cho món ăn"
}
```

### Stock ADJUST (Inventory Correction)
```bash
POST /manager/ingredients/:id/stock-adjust
{
  "new_quantity": 12,      // Target quantity
  "cost_per_unit": 55000,  // Optional: only if increase due to purchase
  "reason": "Kiểm kê định kỳ"
}
```

## Business Rules

### Price Recalculation Rules

| Operation | Quantity Change | Price Provided | Price Different | Recalculate? |
|-----------|----------------|----------------|-----------------|--------------|
| Stock IN | +10 | Yes (50k) | Yes | ✅ YES |
| Stock IN | +10 | Yes (40k) | No (same) | ❌ NO |
| Stock IN | +10 | No (0) | N/A | ❌ NO |
| Stock OUT | -5 | N/A | N/A | ❌ NO |
| Stock ADJUST | +2 | Yes (50k) | Yes | ✅ YES |
| Stock ADJUST | +2 | No (0) | N/A | ❌ NO |
| Stock ADJUST | -3 | N/A | N/A | ❌ NO |

### Weighted Average Formula
```
new_price = (old_qty * old_price + new_qty * new_price) / total_qty
```

Example:
```
Current: 10 kg @ 40,000đ/kg
Purchase: 5 kg @ 50,000đ/kg

new_price = (10 * 40,000 + 5 * 50,000) / 15
          = (400,000 + 250,000) / 15
          = 650,000 / 15
          = 43,333đ/kg
```

## Testing Checklist

### Stock IN Tests
- [ ] Purchase with new higher price → price increases
- [ ] Purchase with new lower price → price decreases
- [ ] Purchase with same price → price unchanged
- [ ] Purchase without price → uses current price
- [ ] Expense is tracked correctly
- [ ] History shows purchase type

### Stock OUT Tests
- [ ] Remove stock → quantity decreases
- [ ] Remove more than available → quantity = 0
- [ ] Price never changes
- [ ] No expense tracked
- [ ] History shows waste type

### Stock ADJUST Tests
- [ ] Increase without price → price unchanged
- [ ] Increase with new price → weighted average
- [ ] Decrease → price unchanged
- [ ] Set to exact quantity → correct
- [ ] Expense tracked only for increases with price
- [ ] History shows adjustment type

## Summary

This refactor provides:
1. ✅ Clear, separate methods for each operation
2. ✅ Simplified frontend logic
3. ✅ Proper constant-based communication
4. ✅ Correct price calculation logic
5. ✅ Better history tracking
6. ✅ Automatic expense tracking
7. ✅ Backward compatibility

The system now has a clean, maintainable architecture that follows business rules correctly and provides a better developer experience.
