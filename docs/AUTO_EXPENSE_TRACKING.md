# Auto Expense Tracking Implementation

## Overview
Automatically create expense records when purchasing ingredients or facilities, and when performing maintenance.

## Business Logic

### 1. Ingredient Purchase → Auto Expense

**Trigger Points:**
1. Create new ingredient (initial stock)
2. Adjust stock with type "IN" (nhập kho)

**Data Flow:**
```
Ingredient Creation/Stock Adjustment
  ↓
Calculate: amount = cost_per_unit × quantity_added
  ↓
Create Expense Record
  ↓
Link to "Nguyên liệu" category
```

**Expense Fields Mapping:**
```javascript
{
  date: new Date(),                              // Current date
  category_id: "INGREDIENT_CATEGORY_ID",         // "Nguyên liệu" category
  amount: cost_per_unit * quantity_added,        // Total cost
  description: `Nhập nguyên liệu: ${ingredient.name}`,
  payment_method: "cash",                        // Default
  vendor: ingredient.supplier || "",             // Supplier name
  notes: `Số lượng: ${quantity_added} ${unit}`  // Additional info
}
```

### 2. Facility Purchase → Auto Expense

**Trigger Points:**
1. Create new facility

**Data Flow:**
```
Facility Creation
  ↓
Use facility.cost as amount
  ↓
Create Expense Record
  ↓
Link to "Cơ sở vật chất" category
```

**Expense Fields Mapping:**
```javascript
{
  date: facility.purchase_date,                  // Purchase date
  category_id: "FACILITY_CATEGORY_ID",           // "Cơ sở vật chất" category
  amount: facility.cost,                         // Purchase cost
  description: `Mua thiết bị: ${facility.name}`,
  payment_method: "cash",                        // Default
  vendor: facility.supplier || "",               // Supplier name
  notes: `Loại: ${facility.type}, Khu vực: ${facility.area}`
}
```

### 3. Facility Maintenance → Auto Expense

**Trigger Points:**
1. Create maintenance record

**Data Flow:**
```
Maintenance Record Creation
  ↓
Use maintenance.cost as amount
  ↓
Create Expense Record
  ↓
Link to "Bảo trì" category
```

**Expense Fields Mapping:**
```javascript
{
  date: maintenance.maintenance_date,            // Maintenance date
  category_id: "MAINTENANCE_CATEGORY_ID",        // "Bảo trì" category
  amount: maintenance.cost,                      // Maintenance cost
  description: `Bảo trì: ${facility.name}`,
  payment_method: "cash",                        // Default
  vendor: "",                                    // Extract from notes if available
  notes: maintenance.notes                       // Maintenance notes
}
```

## Implementation Approach

### Option 1: Backend Service Layer (Recommended)
**Pros:**
- Centralized logic
- Data consistency guaranteed
- Cannot be bypassed
- Easier to maintain

**Cons:**
- Requires backend changes
- More complex testing

**Implementation:**
```go
// In ingredient service
func (s *IngredientService) CreateIngredient(ctx context.Context, ing *ingredient.Ingredient) error {
    // Create ingredient
    err := s.repo.CreateIngredient(ctx, ing)
    if err != nil {
        return err
    }
    
    // Auto-create expense
    if ing.CostPerUnit > 0 && ing.Quantity > 0 {
        expense := &expense.Expense{
            Date:          time.Now(),
            CategoryID:    s.getIngredientCategoryID(ctx),
            Amount:        ing.CostPerUnit * float64(ing.Quantity),
            Description:   fmt.Sprintf("Nhập nguyên liệu: %s", ing.Name),
            PaymentMethod: expense.PaymentMethodCash,
            Vendor:        ing.Supplier,
            Notes:         fmt.Sprintf("Số lượng: %d %s", ing.Quantity, ing.Unit),
        }
        
        if err := s.expenseService.CreateExpense(ctx, expense); err != nil {
            // Log error but don't fail the ingredient creation
            log.Printf("Failed to create auto expense: %v", err)
        }
    }
    
    return nil
}
```

### Option 2: Frontend After Success
**Pros:**
- Easier to implement
- No backend changes needed
- Quick to deploy

**Cons:**
- Can be bypassed
- Less reliable
- Duplicate logic in frontend

**Implementation:**
```javascript
// In ingredient store
async createIngredient(ingredient) {
  try {
    const newIngredient = await ingredientService.createIngredient(ingredient)
    this.ingredients.push(newIngredient)
    
    // Auto-create expense
    if (ingredient.cost_per_unit > 0 && ingredient.quantity > 0) {
      const expense = {
        date: new Date().toISOString().split('T')[0],
        category_id: await this.getIngredientCategoryId(),
        amount: ingredient.cost_per_unit * ingredient.quantity,
        description: `Nhập nguyên liệu: ${ingredient.name}`,
        payment_method: 'cash',
        vendor: ingredient.supplier || '',
        notes: `Số lượng: ${ingredient.quantity} ${ingredient.unit}`
      }
      
      await expenseStore.createExpense(expense)
    }
    
    return true
  } catch (error) {
    this.error = error.response?.data?.error || 'Lỗi tạo nguyên liệu'
    return false
  }
}
```

## Required Expense Categories

Need to ensure these categories exist:
1. **Nguyên liệu** - For ingredient purchases
2. **Cơ sở vật chất** - For facility purchases
3. **Bảo trì** - For maintenance costs

### Category Seeding Script
```javascript
// backend/cmd/seed-expense-categories/main.go
const defaultCategories = [
  { name: "Nguyên liệu" },
  { name: "Cơ sở vật chất" },
  { name: "Bảo trì" },
  { name: "Nhân sự" },
  { name: "Tiện ích" },
  { name: "Marketing" },
  { name: "Khác" }
]
```

## Edge Cases to Handle

### 1. Category Not Found
- Create default categories on first run
- Or use a fallback "Khác" category

### 2. Expense Creation Fails
- Log the error
- Don't fail the main operation (ingredient/facility creation)
- Show warning to user

### 3. Duplicate Prevention
- Check if expense already exists for this ingredient/facility
- Use reference ID to link expense to source

### 4. Update vs Create
- Only create expense on CREATE, not UPDATE
- For stock adjustment, only create expense for "IN" type

### 5. Zero Cost Items
- Skip expense creation if cost is 0
- Or create with 0 amount for tracking purposes

## Database Schema Enhancement

### Add Reference Fields to Expense
```go
type Expense struct {
    // ... existing fields
    
    // Reference to source
    SourceType string             `bson:"source_type,omitempty" json:"source_type,omitempty"` // "ingredient", "facility", "maintenance"
    SourceID   primitive.ObjectID `bson:"source_id,omitempty" json:"source_id,omitempty"`     // ID of source record
}
```

This allows:
- Tracking which expense came from which source
- Preventing duplicates
- Linking back to original record
- Better reporting

## Testing Strategy

### Unit Tests
1. Test expense creation when ingredient created
2. Test expense creation when facility created
3. Test expense creation when maintenance recorded
4. Test error handling when expense creation fails
5. Test category lookup/creation

### Integration Tests
1. Create ingredient → verify expense exists
2. Create facility → verify expense exists
3. Create maintenance → verify expense exists
4. Verify expense amounts are correct
5. Verify expense categories are correct

### Manual Testing Checklist
- [ ] Create ingredient with cost → Check expense created
- [ ] Adjust stock IN → Check expense created
- [ ] Adjust stock OUT → Check no expense created
- [ ] Create facility → Check expense created
- [ ] Create maintenance → Check expense created
- [ ] Verify expense amounts match
- [ ] Verify expense categories are correct
- [ ] Test with missing categories
- [ ] Test with zero cost items

## Rollout Plan

### Phase 1: Setup (Day 1)
1. Create expense categories seed script
2. Run seed script to create default categories
3. Test category creation

### Phase 2: Backend Implementation (Day 2-3)
1. Add SourceType and SourceID to Expense model
2. Implement auto-expense in IngredientService
3. Implement auto-expense in FacilityService
4. Add unit tests
5. Add integration tests

### Phase 3: Testing (Day 4)
1. Manual testing
2. Fix bugs
3. Performance testing

### Phase 4: Deployment (Day 5)
1. Deploy backend changes
2. Monitor logs for errors
3. Verify expenses are being created

## Monitoring & Maintenance

### Metrics to Track
- Number of auto-expenses created per day
- Failed expense creation attempts
- Category distribution
- Average expense amounts

### Logs to Monitor
- Expense creation failures
- Category not found errors
- Invalid data errors

### Alerts to Set
- High failure rate (>5%)
- Missing categories
- Unusual expense amounts

## Future Enhancements

1. **Configurable Auto-Expense**
   - Allow users to enable/disable auto-expense
   - Configure which events trigger expenses

2. **Expense Templates**
   - Pre-defined templates for common expenses
   - Custom fields per category

3. **Approval Workflow**
   - Auto-expenses go to pending state
   - Manager approval required

4. **Bulk Operations**
   - Import ingredients → Create multiple expenses
   - Batch expense creation

5. **Advanced Reporting**
   - Expense by source type
   - Cost analysis per ingredient/facility
   - ROI tracking

## Recommendation

**Use Option 1 (Backend Service Layer)** for production implementation because:
1. More reliable and consistent
2. Cannot be bypassed
3. Easier to maintain long-term
4. Better for audit trail
5. Supports future enhancements

Start with basic implementation, then add enhancements based on user feedback.
