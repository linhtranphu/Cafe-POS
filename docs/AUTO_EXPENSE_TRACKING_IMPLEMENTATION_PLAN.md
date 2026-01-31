# Auto Expense Tracking - Implementation Plan

## Overview
Tự động tạo expense records khi có các hoạt động mua sắm/nhập kho để theo dõi chi phí một cách tự động và chính xác.

## Architecture Design

### Service Layer Architecture
```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                         │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  IngredientService          FacilityService                  │
│         │                          │                          │
│         └──────────┬───────────────┘                          │
│                    │                                          │
│                    ▼                                          │
│         ┌──────────────────────┐                             │
│         │ AutoExpenseService   │  ◄── Central Service        │
│         └──────────────────────┘                             │
│                    │                                          │
│                    ▼                                          │
│         ┌──────────────────────┐                             │
│         │   ExpenseService     │                             │
│         └──────────────────────┘                             │
│                    │                                          │
└────────────────────┼──────────────────────────────────────────┘
                     │
                     ▼
         ┌──────────────────────┐
         │ ExpenseRepository    │
         └──────────────────────┘
                     │
                     ▼
                 MongoDB
```

## Implementation Tasks Breakdown

### Phase 1: Domain & Data Model (2-3 hours)

#### Task 1.1: Update Expense Domain Model
**File:** `backend/domain/expense/expense.go`

**Changes:**
```go
// Add source tracking fields
type Expense struct {
    // ... existing fields
    
    // Source tracking
    SourceType string             `bson:"source_type,omitempty" json:"source_type,omitempty"`
    SourceID   primitive.ObjectID `bson:"source_id,omitempty" json:"source_id,omitempty"`
}

// Add source type constants
const (
    SourceTypeIngredient  = "ingredient"
    SourceTypeFacility    = "facility"
    SourceTypeMaintenance = "maintenance"
    SourceTypeManual      = "manual"  // For manually created expenses
)
```

**Acceptance Criteria:**
- [ ] SourceType field added with proper BSON/JSON tags
- [ ] SourceID field added with proper BSON/JSON tags
- [ ] Constants defined for all source types
- [ ] Backward compatible (existing expenses work without these fields)

---

#### Task 1.2: Create Category Constants
**File:** `backend/domain/expense/category.go` (new file)

**Content:**
```go
package expense

// Default category names
const (
    CategoryIngredient = "Nguyên liệu"
    CategoryFacility   = "Cơ sở vật chất"
    CategoryMaintenance = "Bảo trì"
    CategoryUtility    = "Tiện ích"
    CategorySalary     = "Nhân sự"
    CategoryMarketing  = "Marketing"
    CategoryOther      = "Khác"
)

// GetDefaultCategories returns list of default categories
func GetDefaultCategories() []string {
    return []string{
        CategoryIngredient,
        CategoryFacility,
        CategoryMaintenance,
        CategoryUtility,
        CategorySalary,
        CategoryMarketing,
        CategoryOther,
    }
}
```

**Acceptance Criteria:**
- [ ] All default categories defined
- [ ] Helper function to get category list
- [ ] Vietnamese names used consistently

---

### Phase 2: Auto Expense Service (4-5 hours)

#### Task 2.1: Create AutoExpenseService
**File:** `backend/application/services/auto_expense_service.go` (new file)

**Structure:**
```go
package services

import (
    "context"
    "fmt"
    "time"
    "cafe-pos/backend/domain/expense"
    "cafe-pos/backend/domain/ingredient"
    "cafe-pos/backend/domain/facility"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type AutoExpenseService struct {
    expenseService *ExpenseService
    categoryCache  map[string]primitive.ObjectID  // Cache for category IDs
}

func NewAutoExpenseService(expenseService *ExpenseService) *AutoExpenseService {
    return &AutoExpenseService{
        expenseService: expenseService,
        categoryCache:  make(map[string]primitive.ObjectID),
    }
}

// GetOrCreateCategory gets category ID by name, creates if not exists
func (s *AutoExpenseService) GetOrCreateCategory(ctx context.Context, categoryName string) (primitive.ObjectID, error)

// TrackIngredientPurchase creates expense for ingredient purchase
func (s *AutoExpenseService) TrackIngredientPurchase(ctx context.Context, ing *ingredient.Ingredient, quantity float64) error

// TrackFacilityPurchase creates expense for facility purchase
func (s *AutoExpenseService) TrackFacilityPurchase(ctx context.Context, fac *facility.Facility) error

// TrackMaintenance creates expense for maintenance
func (s *AutoExpenseService) TrackMaintenance(ctx context.Context, facilityID primitive.ObjectID, facilityName string, cost float64, maintenanceDate time.Time, notes string) error
```

**Key Methods:**

1. **GetOrCreateCategory**
   - Check cache first
   - Query database if not in cache
   - Create category if not exists
   - Update cache
   - Return category ID

2. **TrackIngredientPurchase**
   - Calculate amount: `cost_per_unit * quantity`
   - Get/create "Nguyên liệu" category
   - Create expense with:
     - Description: "Nhập nguyên liệu: [name]"
     - Amount: calculated amount
     - Vendor: supplier
     - SourceType: "ingredient"
     - SourceID: ingredient ID
   - Handle errors gracefully (log but don't fail)

3. **TrackFacilityPurchase**
   - Get/create "Cơ sở vật chất" category
   - Create expense with:
     - Description: "Mua thiết bị: [name]"
     - Amount: facility cost
     - Date: purchase_date
     - Vendor: supplier
     - SourceType: "facility"
     - SourceID: facility ID

4. **TrackMaintenance**
   - Get/create "Bảo trì" category
   - Create expense with:
     - Description: "Bảo trì: [facility_name]"
     - Amount: maintenance cost
     - Date: maintenance_date
     - Notes: maintenance notes
     - SourceType: "maintenance"
     - SourceID: facility ID

**Acceptance Criteria:**
- [ ] Service created with all methods
- [ ] Category caching implemented
- [ ] Error handling doesn't fail main operations
- [ ] Logging for debugging
- [ ] Thread-safe category cache

---

#### Task 2.2: Add Unit Tests for AutoExpenseService
**File:** `backend/application/services/auto_expense_service_test.go` (new file)

**Test Cases:**
```go
func TestGetOrCreateCategory(t *testing.T)
func TestGetOrCreateCategory_Caching(t *testing.T)
func TestTrackIngredientPurchase(t *testing.T)
func TestTrackIngredientPurchase_ZeroCost(t *testing.T)
func TestTrackFacilityPurchase(t *testing.T)
func TestTrackMaintenance(t *testing.T)
func TestAutoExpense_ErrorHandling(t *testing.T)
```

**Acceptance Criteria:**
- [ ] All methods have unit tests
- [ ] Edge cases covered (zero cost, missing data)
- [ ] Error scenarios tested
- [ ] Mock dependencies properly
- [ ] Test coverage > 80%

---

### Phase 3: Integration with Existing Services (3-4 hours)

#### Task 3.1: Update IngredientService
**File:** `backend/application/services/ingredient.go`

**Changes:**
```go
type IngredientService struct {
    repo               *mongodb.IngredientRepository
    autoExpenseService *AutoExpenseService  // Add this
}

// Update constructor
func NewIngredientService(repo *mongodb.IngredientRepository, autoExpenseService *AutoExpenseService) *IngredientService {
    return &IngredientService{
        repo:               repo,
        autoExpenseService: autoExpenseService,
    }
}

// Update CreateIngredient
func (s *IngredientService) CreateIngredient(ctx context.Context, ing *ingredient.Ingredient) error {
    // Create ingredient
    err := s.repo.CreateIngredient(ctx, ing)
    if err != nil {
        return err
    }
    
    // Auto-track expense if has cost and quantity
    if ing.CostPerUnit > 0 && ing.Quantity > 0 {
        if err := s.autoExpenseService.TrackIngredientPurchase(ctx, ing, ing.Quantity); err != nil {
            // Log error but don't fail
            log.Printf("Failed to auto-track ingredient expense: %v", err)
        }
    }
    
    return nil
}

// Update AdjustStock (for stock IN)
func (s *IngredientService) AdjustStock(ctx context.Context, id primitive.ObjectID, adjustment *ingredient.StockAdjustment) error {
    // Get ingredient first
    ing, err := s.repo.GetIngredient(ctx, id)
    if err != nil {
        return err
    }
    
    // Adjust stock
    err = s.repo.AdjustStock(ctx, id, adjustment)
    if err != nil {
        return err
    }
    
    // Auto-track expense for stock IN
    if adjustment.Type == ingredient.AdjustmentTypeAdd && ing.CostPerUnit > 0 {
        if err := s.autoExpenseService.TrackIngredientPurchase(ctx, ing, adjustment.Quantity); err != nil {
            log.Printf("Failed to auto-track stock adjustment expense: %v", err)
        }
    }
    
    return nil
}
```

**Acceptance Criteria:**
- [ ] AutoExpenseService injected via constructor
- [ ] CreateIngredient calls auto-tracking
- [ ] AdjustStock calls auto-tracking for "IN" type only
- [ ] Errors logged but don't fail main operation
- [ ] Existing tests still pass

---

#### Task 3.2: Update FacilityService
**File:** `backend/application/services/facility_service.go`

**Changes:**
```go
type FacilityService struct {
    repo               *mongodb.FacilityRepository
    autoExpenseService *AutoExpenseService  // Add this
}

// Update constructor
func NewFacilityService(repo *mongodb.FacilityRepository, autoExpenseService *AutoExpenseService) *FacilityService {
    return &FacilityService{
        repo:               repo,
        autoExpenseService: autoExpenseService,
    }
}

// Update CreateFacility
func (s *FacilityService) CreateFacility(ctx context.Context, fac *facility.Facility) error {
    // Create facility
    err := s.repo.CreateFacility(ctx, fac)
    if err != nil {
        return err
    }
    
    // Auto-track expense if has cost
    if fac.Cost > 0 {
        if err := s.autoExpenseService.TrackFacilityPurchase(ctx, fac); err != nil {
            log.Printf("Failed to auto-track facility expense: %v", err)
        }
    }
    
    return nil
}

// Update CreateMaintenanceRecord
func (s *FacilityService) CreateMaintenanceRecord(ctx context.Context, record *facility.MaintenanceRecord) error {
    // Get facility first
    fac, err := s.repo.GetFacility(ctx, record.FacilityID)
    if err != nil {
        return err
    }
    
    // Create maintenance record
    err = s.repo.CreateMaintenanceRecord(ctx, record)
    if err != nil {
        return err
    }
    
    // Auto-track expense if has cost
    if record.Cost > 0 {
        if err := s.autoExpenseService.TrackMaintenance(ctx, fac.ID, fac.Name, record.Cost, record.MaintenanceDate, record.Notes); err != nil {
            log.Printf("Failed to auto-track maintenance expense: %v", err)
        }
    }
    
    return nil
}
```

**Acceptance Criteria:**
- [ ] AutoExpenseService injected via constructor
- [ ] CreateFacility calls auto-tracking
- [ ] CreateMaintenanceRecord calls auto-tracking
- [ ] Errors logged but don't fail main operation
- [ ] Existing tests still pass

---

#### Task 3.3: Update Service Initialization in main.go
**File:** `backend/main.go`

**Changes:**
```go
// Initialize services with proper dependencies
expenseRepo := mongodb.NewExpenseRepository(db)
expenseService := services.NewExpenseService(expenseRepo)

// Create auto expense service
autoExpenseService := services.NewAutoExpenseService(expenseService)

// Update ingredient service initialization
ingredientRepo := mongodb.NewIngredientRepository(db)
ingredientService := services.NewIngredientService(ingredientRepo, autoExpenseService)

// Update facility service initialization
facilityRepo := mongodb.NewFacilityRepository(db)
facilityService := services.NewFacilityService(facilityRepo, autoExpenseService)
```

**Acceptance Criteria:**
- [ ] Services initialized in correct order
- [ ] Dependencies properly injected
- [ ] No circular dependencies
- [ ] Application compiles and runs

---

### Phase 4: Database Seeding (1-2 hours)

#### Task 4.1: Create Expense Category Seed Script
**File:** `backend/cmd/seed-expense-categories/main.go` (new file)

**Content:**
```go
package main

import (
    "context"
    "log"
    "cafe-pos/backend/domain/expense"
    "cafe-pos/backend/infrastructure/mongodb"
    // ... imports
)

func main() {
    // Connect to MongoDB
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(context.Background())
    
    db := client.Database("cafe_pos")
    repo := mongodb.NewExpenseRepository(db)
    
    // Get default categories
    categories := expense.GetDefaultCategories()
    
    // Check and create each category
    for _, catName := range categories {
        // Check if exists
        existing, _ := repo.GetCategories(context.Background())
        found := false
        for _, cat := range existing {
            if cat.Name == catName {
                found = true
                break
            }
        }
        
        // Create if not exists
        if !found {
            cat := &expense.Category{
                Name: catName,
            }
            if err := repo.CreateCategory(context.Background(), cat); err != nil {
                log.Printf("Failed to create category %s: %v", catName, err)
            } else {
                log.Printf("Created category: %s", catName)
            }
        } else {
            log.Printf("Category already exists: %s", catName)
        }
    }
    
    log.Println("Expense categories seeding completed!")
}
```

**Acceptance Criteria:**
- [ ] Script creates all default categories
- [ ] Checks for existing categories before creating
- [ ] Idempotent (can run multiple times safely)
- [ ] Proper error handling and logging
- [ ] Can be run via `go run backend/cmd/seed-expense-categories/main.go`

---

#### Task 4.2: Add Seed Script to Makefile
**File:** `Makefile`

**Add:**
```makefile
.PHONY: seed-expense-categories
seed-expense-categories:
	@echo "Seeding expense categories..."
	@go run backend/cmd/seed-expense-categories/main.go
```

**Acceptance Criteria:**
- [ ] Command added to Makefile
- [ ] Can run via `make seed-expense-categories`
- [ ] Documented in README or scripts/README.md

---

### Phase 5: Testing & Validation (2-3 hours)

#### Task 5.1: Integration Tests
**File:** `backend/application/services/auto_expense_integration_test.go` (new file)

**Test Scenarios:**
```go
func TestIngredientPurchase_CreatesExpense(t *testing.T)
func TestStockAdjustmentIN_CreatesExpense(t *testing.T)
func TestStockAdjustmentOUT_DoesNotCreateExpense(t *testing.T)
func TestFacilityPurchase_CreatesExpense(t *testing.T)
func TestMaintenance_CreatesExpense(t *testing.T)
func TestExpenseHasCorrectSourceTracking(t *testing.T)
func TestExpenseAmountsAreCorrect(t *testing.T)
func TestCategoryAutoCreation(t *testing.T)
```

**Acceptance Criteria:**
- [ ] All integration tests pass
- [ ] Tests use real MongoDB (testcontainers or local)
- [ ] Tests verify expense creation
- [ ] Tests verify expense data accuracy
- [ ] Tests verify source tracking

---

#### Task 5.2: Manual Testing Checklist
**Document:** `docs/AUTO_EXPENSE_TESTING_CHECKLIST.md`

**Checklist:**
```markdown
## Ingredient Purchase Testing
- [ ] Create ingredient with cost → Expense created
- [ ] Create ingredient without cost → No expense created
- [ ] Adjust stock IN → Expense created
- [ ] Adjust stock OUT → No expense created
- [ ] Adjust stock SET → No expense created
- [ ] Verify expense amount = cost_per_unit × quantity
- [ ] Verify expense category = "Nguyên liệu"
- [ ] Verify expense vendor = ingredient supplier
- [ ] Verify source_type = "ingredient"
- [ ] Verify source_id = ingredient ID

## Facility Purchase Testing
- [ ] Create facility with cost → Expense created
- [ ] Create facility without cost → No expense created
- [ ] Verify expense amount = facility cost
- [ ] Verify expense date = purchase_date
- [ ] Verify expense category = "Cơ sở vật chất"
- [ ] Verify expense vendor = facility supplier
- [ ] Verify source_type = "facility"
- [ ] Verify source_id = facility ID

## Maintenance Testing
- [ ] Create maintenance with cost → Expense created
- [ ] Create maintenance without cost → No expense created
- [ ] Verify expense amount = maintenance cost
- [ ] Verify expense date = maintenance_date
- [ ] Verify expense category = "Bảo trì"
- [ ] Verify source_type = "maintenance"
- [ ] Verify source_id = facility ID

## Category Testing
- [ ] Categories auto-created on first use
- [ ] Categories cached properly
- [ ] Multiple operations use same category ID

## Error Handling
- [ ] Expense creation failure doesn't fail main operation
- [ ] Errors are logged properly
- [ ] System continues to work if expense service is down
```

**Acceptance Criteria:**
- [ ] All checklist items tested
- [ ] Results documented
- [ ] Bugs found and fixed

---

### Phase 6: Documentation & Deployment (1-2 hours)

#### Task 6.1: Update API Documentation
**File:** `docs/API_DOCUMENTATION.md`

**Add:**
- Expense source_type field documentation
- Expense source_id field documentation
- Auto-expense behavior documentation

**Acceptance Criteria:**
- [ ] API docs updated
- [ ] Examples provided
- [ ] Source tracking explained

---

#### Task 6.2: Create User Guide
**File:** `docs/AUTO_EXPENSE_USER_GUIDE.md`

**Content:**
- How auto-expense works
- What triggers expense creation
- How to view auto-created expenses
- How to identify auto vs manual expenses
- Troubleshooting guide

**Acceptance Criteria:**
- [ ] User guide created
- [ ] Clear explanations
- [ ] Screenshots/examples included
- [ ] Vietnamese language

---

#### Task 6.3: Deployment Steps
**File:** `docs/AUTO_EXPENSE_DEPLOYMENT.md`

**Steps:**
1. Run database migrations (if any)
2. Run seed script for categories
3. Deploy backend with new code
4. Verify auto-expense is working
5. Monitor logs for errors
6. Rollback plan if issues

**Acceptance Criteria:**
- [ ] Deployment steps documented
- [ ] Rollback plan included
- [ ] Monitoring checklist provided

---

## Summary of Files to Create/Modify

### New Files (9 files)
1. `backend/domain/expense/category.go`
2. `backend/application/services/auto_expense_service.go`
3. `backend/application/services/auto_expense_service_test.go`
4. `backend/application/services/auto_expense_integration_test.go`
5. `backend/cmd/seed-expense-categories/main.go`
6. `docs/AUTO_EXPENSE_TESTING_CHECKLIST.md`
7. `docs/AUTO_EXPENSE_USER_GUIDE.md`
8. `docs/AUTO_EXPENSE_DEPLOYMENT.md`
9. `docs/AUTO_EXPENSE_TRACKING_IMPLEMENTATION_PLAN.md` (this file)

### Modified Files (5 files)
1. `backend/domain/expense/expense.go` - Add source tracking fields
2. `backend/application/services/ingredient.go` - Integrate auto-expense
3. `backend/application/services/facility_service.go` - Integrate auto-expense
4. `backend/main.go` - Update service initialization
5. `Makefile` - Add seed command

### Total: 14 files

## Estimated Timeline

| Phase | Tasks | Estimated Time |
|-------|-------|----------------|
| Phase 1: Domain & Data Model | 2 tasks | 2-3 hours |
| Phase 2: Auto Expense Service | 2 tasks | 4-5 hours |
| Phase 3: Service Integration | 3 tasks | 3-4 hours |
| Phase 4: Database Seeding | 2 tasks | 1-2 hours |
| Phase 5: Testing & Validation | 2 tasks | 2-3 hours |
| Phase 6: Documentation & Deployment | 3 tasks | 1-2 hours |
| **Total** | **14 tasks** | **13-19 hours** |

## Risk Assessment

### High Risk
- **Service initialization order**: Must ensure ExpenseService is created before AutoExpenseService
- **Circular dependencies**: Careful with service dependencies

### Medium Risk
- **Category caching**: Thread-safety concerns
- **Error handling**: Must not fail main operations

### Low Risk
- **Database seeding**: Idempotent script
- **Testing**: Comprehensive test coverage

## Success Criteria

1. ✅ All 14 tasks completed
2. ✅ All tests passing (unit + integration)
3. ✅ Manual testing checklist completed
4. ✅ Zero impact on existing functionality
5. ✅ Auto-expenses created correctly
6. ✅ Source tracking working
7. ✅ Documentation complete
8. ✅ Successfully deployed to production

## Next Steps

1. Review this implementation plan
2. Get approval from stakeholders
3. Create tasks in project management tool
4. Assign tasks to developers
5. Begin Phase 1 implementation
6. Regular progress reviews
7. Testing and validation
8. Deployment

---

**Document Version:** 1.0  
**Last Updated:** 2026-01-31  
**Author:** Development Team  
**Status:** Ready for Review
