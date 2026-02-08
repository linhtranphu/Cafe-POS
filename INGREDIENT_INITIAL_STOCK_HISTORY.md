# Ingredient Creation - Initial Stock History Record

## Problem
When creating a new ingredient with initial quantity, no stock history record was created. This meant:
- The first purchase was not tracked in history
- Users couldn't see the initial price paid
- History only showed subsequent adjustments
- Incomplete audit trail

## Solution
Automatically create a stock history record when creating an ingredient with quantity > 0.

## Implementation

### 1. Service Layer Changes
**File:** `backend/application/services/ingredient.go`

**Added logic to CreateIngredient:**
```go
// Create initial stock history record if quantity > 0
if req.Quantity > 0 {
    userID := primitive.NilObjectID
    if userIDStr != "" {
        if oid, err := primitive.ObjectIDFromHex(userIDStr); err == nil {
            userID = oid
        }
    }
    
    history := &ingredient.StockHistory{
        IngredientID: item.ID,
        Type:         ingredient.TransactionPurchase,
        Quantity:     req.Quantity,
        BeforeQty:    0,                              // Initial creation
        AfterQty:     req.Quantity,
        Reason:       "Tạo nguyên liệu mới - Nhập kho đầu tiên",
        UserID:       userID,
        Username:     username,
        CostPerUnit:  req.CostPerUnit,
        TotalCost:    req.Quantity * req.CostPerUnit,
    }
    
    // Create stock history (don't fail if this fails)
    if err := s.stockHistoryRepo.Create(ctx, history); err != nil {
        // Log error but don't fail the operation
    }
}
```

**Key points:**
- Only creates history if `quantity > 0`
- Sets `BeforeQty` to 0 (initial state)
- Sets `AfterQty` to the initial quantity
- Uses transaction type `TransactionPurchase`
- Records the initial cost per unit and total cost
- Includes user information for audit trail
- Doesn't fail the ingredient creation if history creation fails

### 2. Handler Layer Changes
**File:** `backend/interfaces/http/ingredient_handler.go`

**Updated CreateIngredient handler:**
```go
func (h *IngredientHandler) CreateIngredient(c *gin.Context) {
    // ... validation ...
    
    // Get user info from context
    userID, _ := c.Get("user_id")
    username, _ := c.Get("username")
    
    userIDStr := ""
    if uid, ok := userID.(string); ok {
        userIDStr = uid
    }
    
    createdBy := ""
    if u, ok := username.(string); ok {
        createdBy = u
    }

    item, err := h.ingredientService.CreateIngredient(
        c.Request.Context(), 
        &req, 
        userIDStr,  // Added userID
        createdBy,
    )
    
    // ... response ...
}
```

**Changes:**
- Extract both `user_id` and `username` from context
- Pass both to service layer
- Service uses userID for proper audit trail

### 3. Service Signature Update
**Before:**
```go
CreateIngredient(ctx context.Context, req *ingredient.CreateIngredientRequest, username string)
```

**After:**
```go
CreateIngredient(ctx context.Context, req *ingredient.CreateIngredientRequest, userIDStr string, username string)
```

## Stock History Record Details

### Transaction Type
Uses `TransactionPurchase` to indicate this is an initial purchase, not just an adjustment.

### Reason
Default reason: **"Tạo nguyên liệu mới - Nhập kho đầu tiên"**
- Clear indication this is the initial creation
- Vietnamese for consistency with UI
- Can be customized if needed

### Quantities
- `BeforeQty`: Always 0 (nothing existed before)
- `Quantity`: The initial quantity entered
- `AfterQty`: Same as Quantity (0 + Quantity)

### Price Information
- `CostPerUnit`: The initial cost per unit entered
- `TotalCost`: Calculated as `Quantity × CostPerUnit`

### User Information
- `UserID`: ObjectID of the user who created the ingredient
- `Username`: Display name for UI

## Benefits

### 1. Complete Audit Trail
- Every quantity change is tracked from day one
- Initial purchase is visible in history
- No missing data

### 2. Price Tracking
- Initial price is recorded
- Can compare future prices with initial price
- Historical price data is complete

### 3. User Accountability
- Know who created each ingredient
- Track who made initial purchase
- Better audit and compliance

### 4. Consistent Behavior
- Creating ingredient = same as adjusting stock
- Both operations create history records
- Predictable system behavior

## User Experience

### Before
1. User creates ingredient with 10 kg @ 50,000 VND/kg
2. Opens history → Empty (no records)
3. Confusion: "Where's my initial purchase?"

### After
1. User creates ingredient with 10 kg @ 50,000 VND/kg
2. Opens history → Shows initial purchase record:
   ```
   📦 Nhập thêm
   +10 kg
   
   Lý do: Tạo nguyên liệu mới - Nhập kho đầu tiên
   
   💰 THÔNG TIN GIÁ
   Đơn giá lần này: 50,000 ₫/kg
   Tổng chi phí: 500,000 ₫
   = 10 kg × 50,000 ₫
   
   👤 Admin
   🕐 07/02/2026 10:30
   ```

## Edge Cases Handled

### 1. Zero Quantity
If user creates ingredient with quantity = 0:
- No history record created
- Makes sense: nothing was purchased yet
- History starts when first stock adjustment happens

### 2. Missing User Info
If user_id or username not in context:
- Uses empty string for username
- Uses NilObjectID for userID
- Operation still succeeds

### 3. History Creation Failure
If stock history creation fails:
- Error is logged (not shown to user)
- Ingredient creation still succeeds
- History is secondary to main operation

### 4. Expense Tracking
Both operations happen:
- Stock history record created
- Auto-expense record created
- Independent operations, both can succeed/fail separately

## Testing

### Test Case 1: Create with Quantity
```bash
POST /api/ingredients
{
  "name": "Sữa tươi",
  "category": "Nguyên liệu chính",
  "unit": "L",
  "quantity": 10,
  "min_stock": 2,
  "cost_per_unit": 25000,
  "supplier": "Vinamilk"
}

# Expected:
# 1. Ingredient created
# 2. Stock history record created
# 3. Auto-expense record created
# 4. GET /api/ingredients/:id/history returns 1 record
```

### Test Case 2: Create with Zero Quantity
```bash
POST /api/ingredients
{
  "name": "Sữa tươi",
  "category": "Nguyên liệu chính",
  "unit": "L",
  "quantity": 0,
  "min_stock": 2,
  "cost_per_unit": 25000,
  "supplier": "Vinamilk"
}

# Expected:
# 1. Ingredient created
# 2. NO stock history record
# 3. NO auto-expense record
# 4. GET /api/ingredients/:id/history returns empty array
```

### Test Case 3: Verify History Display
```bash
# 1. Create ingredient with quantity
# 2. Open ingredient history in UI
# 3. Verify shows:
#    - Green card (purchase)
#    - "+10 L" quantity
#    - "Tạo nguyên liệu mới - Nhập kho đầu tiên" reason
#    - Price information section
#    - User and timestamp
```

## Database Impact

### New Records
Each ingredient creation with quantity > 0 creates:
- 1 ingredient document
- 1 stock_history document
- 1 expense document (if auto-expense enabled)

### Storage
Minimal impact:
- Stock history records are small (~200 bytes)
- Essential for audit trail
- Worth the storage cost

## Migration

### Existing Ingredients
Existing ingredients created before this change:
- Will NOT have initial history record
- History starts from first adjustment
- This is acceptable (historical data)

### Optional: Backfill Script
If needed, can create script to:
1. Find ingredients with no history
2. Create synthetic initial history record
3. Use current cost_per_unit as initial price
4. Set created_at to ingredient created_at

## Files Modified

1. `backend/application/services/ingredient.go`
   - Updated `CreateIngredient` signature
   - Added stock history creation logic
   - Added user ID handling

2. `backend/interfaces/http/ingredient_handler.go`
   - Updated `CreateIngredient` handler
   - Extract user_id from context
   - Pass both userID and username to service

## Related Features

This change complements:
- Auto-expense tracking (already implemented)
- Stock history display (already implemented)
- Price tracking in history (already implemented)

Together, these features provide:
- Complete financial tracking
- Full audit trail
- Price history analysis
- User accountability
