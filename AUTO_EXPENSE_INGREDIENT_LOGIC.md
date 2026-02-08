# Auto Expense Creation for Ingredient Purchase

## Overview
When creating or restocking an ingredient, the system automatically creates an expense record to track the purchase cost.

## Flow Diagram
```
User creates/restocks ingredient
         ↓
IngredientHandler.CreateIngredient()
         ↓
IngredientService.CreateIngredient()
         ↓
Save ingredient to database
         ↓
AutoExpenseService.TrackIngredientPurchase()
         ↓
Calculate: amount = costPerUnit × quantity
         ↓
Get or create "Nguyên liệu" category
         ↓
Create expense record with:
  - Amount: calculated cost
  - Description: "Nhập nguyên liệu: {name}"
  - Source: ingredient
  - CreatedBy: username
         ↓
Done ✓
```

## Code Flow

### 1. Handler Layer
**File**: `backend/interfaces/http/ingredient_handler.go`

```go
func (h *IngredientHandler) CreateIngredient(c *gin.Context) {
    // Get username from JWT token context
    username, _ := c.Get("username")
    createdBy := ""
    if u, ok := username.(string); ok {
        createdBy = u
    }
    
    // Pass username to service
    item, err := h.ingredientService.CreateIngredient(ctx, &req, createdBy)
}
```

### 2. Service Layer
**File**: `backend/application/services/ingredient.go`

```go
func (s *IngredientService) CreateIngredient(ctx, req, username) {
    // Create ingredient
    item := &ingredient.Ingredient{...}
    s.ingredientRepo.Create(ctx, item)
    
    // Auto-track expense if configured and quantity > 0
    if s.autoExpenseService != nil && req.Quantity > 0 {
        s.autoExpenseService.TrackIngredientPurchase(ctx, item, req.Quantity, username)
    }
}
```

### 3. Auto Expense Service
**File**: `backend/application/services/auto_expense_service.go`

```go
func (s *AutoExpenseService) TrackIngredientPurchase(ctx, ing, quantity, username) {
    // Skip if no cost or quantity
    if ing.CostPerUnit <= 0 || quantity <= 0 {
        return nil
    }
    
    // Calculate amount
    amount := ing.CostPerUnit * quantity
    
    // Get or create category
    categoryID := s.GetOrCreateCategory(ctx, "Nguyên liệu")
    
    // Create expense
    exp := &expense.Expense{
        Date:          time.Now(),
        CategoryID:    categoryID,
        Amount:        amount,
        Description:   fmt.Sprintf("Nhập nguyên liệu: %s", ing.Name),
        PaymentMethod: "cash",
        Vendor:        ing.Supplier,
        Notes:         fmt.Sprintf("Số lượng: %.2f %s", quantity, ing.Unit),
        SourceType:    "ingredient",
        SourceID:      ing.ID,
        CreatedBy:     username,
    }
    
    s.expenseService.CreateExpense(ctx, exp)
}
```

## Initialization

**File**: `backend/main.go`

```go
// Create services
expenseService := services.NewExpenseService(expenseRepo)
ingredientService := services.NewIngredientService(ingredientRepo, stockHistoryRepo)

// Wire up auto expense service
autoExpenseService := services.NewAutoExpenseService(expenseService)
ingredientService.SetAutoExpenseService(autoExpenseService)
```

## Key Features

### 1. Automatic Tracking
- ✅ Triggered on ingredient creation
- ✅ Triggered on stock adjustment (positive quantity)
- ✅ No manual expense entry needed

### 2. Smart Skipping
- ✅ Skip if `costPerUnit <= 0` (free items)
- ✅ Skip if `quantity <= 0` (no purchase)
- ✅ Non-blocking: expense tracking failure doesn't fail ingredient creation

### 3. Category Management
- ✅ Auto-creates "Nguyên liệu" category if not exists
- ✅ Caches category IDs to minimize DB queries
- ✅ Thread-safe with mutex locks

### 4. Expense Details
- **Amount**: `costPerUnit × quantity`
- **Description**: "Nhập nguyên liệu: {ingredient name}"
- **Category**: "Nguyên liệu" (auto-created)
- **Payment Method**: "cash" (default)
- **Vendor**: From ingredient supplier field
- **Notes**: "Số lượng: {quantity} {unit}"
- **Source Type**: "ingredient"
- **Source ID**: Ingredient ID (for traceability)
- **Created By**: Username from JWT token

### 5. Stock Adjustment
When adjusting stock IN (positive quantity):
```go
func (s *IngredientService) AdjustStock(ctx, id, req) {
    // Update stock
    item.Quantity += req.Quantity
    
    // Track expense for stock IN
    if s.autoExpenseService != nil && req.Quantity > 0 {
        s.autoExpenseService.TrackIngredientPurchase(ctx, item, req.Quantity, req.Username)
    }
}
```

## Testing

**File**: `backend/application/services/auto_expense_service_test.go`

Tests cover:
- ✅ Track ingredient purchase
- ✅ Skip zero cost ingredient
- ✅ Skip zero quantity
- ✅ Category caching
- ✅ Concurrent access safety
- ✅ Expense record verification

## Example Scenarios

### Scenario 1: Create new ingredient
```
Input:
  Name: "Cà phê hạt"
  Quantity: 10 kg
  CostPerUnit: 200,000 VND
  Supplier: "Nhà cung cấp A"
  CreatedBy: "Admin"

Result:
  Ingredient created ✓
  Expense created ✓
    Amount: 2,000,000 VND
    Description: "Nhập nguyên liệu: Cà phê hạt"
    Category: "Nguyên liệu"
    Vendor: "Nhà cung cấp A"
    Notes: "Số lượng: 10.00 kg"
    CreatedBy: "Admin"
```

### Scenario 2: Restock ingredient
```
Input:
  Ingredient: "Đường"
  Adjust: +5 kg
  CostPerUnit: 15,000 VND
  Username: "Manager"

Result:
  Stock updated ✓
  Expense created ✓
    Amount: 75,000 VND
    Description: "Nhập nguyên liệu: Đường"
    CreatedBy: "Manager"
```

### Scenario 3: Free sample (no expense)
```
Input:
  Name: "Sample Coffee"
  Quantity: 1 kg
  CostPerUnit: 0 VND

Result:
  Ingredient created ✓
  Expense NOT created (zero cost)
```

## Benefits

1. **Automatic Tracking**: No manual expense entry needed
2. **Accurate Costing**: Expense amount matches actual purchase
3. **Traceability**: Link expense back to ingredient via SourceID
4. **Audit Trail**: Track who created the expense
5. **Consistency**: Same logic for create and restock
6. **Performance**: Category caching reduces DB queries
7. **Reliability**: Non-blocking, won't fail ingredient operations

## Related Files

- `backend/application/services/ingredient.go` - Ingredient service
- `backend/application/services/auto_expense_service.go` - Auto expense logic
- `backend/application/services/auto_expense_service_test.go` - Tests
- `backend/interfaces/http/ingredient_handler.go` - HTTP handler
- `backend/main.go` - Service initialization

---
**Status**: ✅ Implemented and tested
**Date**: 2026-02-07
**Pattern**: Auto-tracking, non-blocking, cached categories
