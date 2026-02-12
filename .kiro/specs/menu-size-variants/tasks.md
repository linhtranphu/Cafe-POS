# Menu Size Variants - Task List

## Phase 1: Backend Domain Layer (Days 1-2)

### 1.1 Update Menu Domain Model
- [ ] 1.1.1 Add MenuItemVariant struct to `backend/domain/menu/menu.go`
  - [ ] Add ID, Name, Price, Ingredients fields
  - [ ] Add Available, IsDefault fields
  - [ ] Add cost tracking fields (CurrentCost, CostStatus, CostLastCalculatedAt)
- [ ] 1.1.2 Update MenuItem struct
  - [ ] Add HasVariants boolean field
  - [ ] Add Variants array field
  - [ ] Keep existing Price, Ingredients fields for backward compatibility
- [ ] 1.1.3 Add MenuItem helper methods
  - [ ] Implement GetDefaultVariant()
  - [ ] Implement GetVariantByID(variantID string)
  - [ ] Implement GetPrice(variantID string)
  - [ ] Implement GetIngredients(variantID string)
- [ ] 1.1.4 Add MenuItem validation
  - [ ] Implement Validate() method
  - [ ] Validate has_variants = true requires variants
  - [ ] Validate exactly one default variant
  - [ ] Validate unique variant IDs
  - [ ] Validate variant prices > 0
  - [ ] Validate single-size price > 0
- [ ] 1.1.5 Write unit tests for MenuItem
  - [ ] Test GetDefaultVariant() - found/not found
  - [ ] Test GetVariantByID() - found/not found
  - [ ] Test GetPrice() - single-size
  - [ ] Test GetPrice() - with variant
  - [ ] Test GetPrice() - with invalid variant (fallback to default)
  - [ ] Test Validate() - valid single-size
  - [ ] Test Validate() - valid multi-size
  - [ ] Test Validate() - no variants when has_variants=true
  - [ ] Test Validate() - no default variant
  - [ ] Test Validate() - multiple default variants
  - [ ] Test Validate() - duplicate variant IDs
  - [ ] Test Validate() - invalid variant price

### 1.2 Update Request/Response DTOs
- [ ] 1.2.1 Update CreateMenuItemRequest in `backend/domain/menu/menu.go`
  - [ ] Add HasVariants field
  - [ ] Add Variants array field
  - [ ] Keep existing Price, Ingredients fields
- [ ] 1.2.2 Update UpdateMenuItemRequest
  - [ ] Add HasVariants pointer field
  - [ ] Add Variants array field
  - [ ] Keep existing Price, Ingredients fields

### 1.3 Update Order Domain Model
- [ ] 1.3.1 Update OrderItem struct in `backend/domain/order/order.go`
  - [ ] Add VariantID string field (optional)
  - [ ] Add VariantName string field (optional)
- [ ] 1.3.2 Write unit tests for OrderItem
  - [ ] Test OrderItem with variant
  - [ ] Test OrderItem without variant

## Phase 2: Backend Service Layer (Days 3-4)

### 2.1 Update Menu Service
- [ ] 2.1.1 Update CreateMenuItem in `backend/application/services/menu.go`
  - [ ] Handle has_variants = false (single-size)
  - [ ] Handle has_variants = true (multi-size)
  - [ ] Validate variants if has_variants = true
  - [ ] Ensure exactly one default variant
  - [ ] Call item.Validate()
- [ ] 2.1.2 Update UpdateMenuItem
  - [ ] Handle toggling has_variants
  - [ ] Update variants array
  - [ ] Validate changes
- [ ] 2.1.3 Write unit tests for MenuService
  - [ ] Test CreateMenuItem - single-size
  - [ ] Test CreateMenuItem - multi-size with valid variants
  - [ ] Test CreateMenuItem - multi-size with no variants (error)
  - [ ] Test CreateMenuItem - multi-size with no default (error)
  - [ ] Test CreateMenuItem - multi-size with multiple defaults (error)
  - [ ] Test UpdateMenuItem - single to multi-size
  - [ ] Test UpdateMenuItem - multi to single-size
  - [ ] Test UpdateMenuItem - update variants

### 2.2 Update Order Service
- [ ] 2.2.1 Update CreateOrder in `backend/application/services/order_service.go`
  - [ ] For each order item, get menu item
  - [ ] Check if menu item has variants
  - [ ] If has_variants, require variant_id
  - [ ] If has_variants, validate variant_id exists
  - [ ] If has_variants, get price from variant
  - [ ] If single-size, get price from item
  - [ ] Populate variant_name in order item
  - [ ] Calculate subtotal
- [ ] 2.2.2 Update EditOrder
  - [ ] Handle variant_id in edited items
  - [ ] Validate variant_id if provided
- [ ] 2.2.3 Write unit tests for OrderService
  - [ ] Test CreateOrder - single-size item
  - [ ] Test CreateOrder - multi-size item with variant_id
  - [ ] Test CreateOrder - multi-size item without variant_id (error)
  - [ ] Test CreateOrder - multi-size item with invalid variant_id (error)
  - [ ] Test CreateOrder - unavailable variant (error)
  - [ ] Test EditOrder - with variant changes

### 2.3 Update Cost Calculator Service
- [ ] 2.3.1 Update CalculateMenuItemCost in `backend/application/services/cost_calculator_service.go`
  - [ ] Check if menu item has variants
  - [ ] If has_variants, loop through variants
  - [ ] Calculate cost per variant
  - [ ] Update variant cost fields
  - [ ] If single-size, calculate cost normally
- [ ] 2.3.2 Write unit tests for CostCalculatorService
  - [ ] Test CalculateMenuItemCost - single-size
  - [ ] Test CalculateMenuItemCost - multi-size with variants
  - [ ] Test cost calculation with different ingredient quantities per variant

## Phase 3: Backend API Layer (Day 5)

### 3.1 Update Menu Handler
- [ ] 3.1.1 Verify CreateMenuItem handler in `backend/interfaces/http/menu_handler.go`
  - [ ] Ensure it accepts has_variants and variants fields
  - [ ] Ensure validation errors are returned properly
- [ ] 3.1.2 Verify UpdateMenuItem handler
  - [ ] Ensure it accepts has_variants and variants fields
- [ ] 3.1.3 Test API endpoints manually
  - [ ] POST /api/menu - single-size item
  - [ ] POST /api/menu - multi-size item
  - [ ] GET /api/menu - list items with variants
  - [ ] GET /api/menu/:id - get item with variants
  - [ ] PUT /api/menu/:id - update item with variants
  - [ ] DELETE /api/menu/:id - delete item with variants

### 3.2 Update Order Handler
- [ ] 3.2.1 Verify CreateOrder handler in `backend/interfaces/http/order_handler.go`
  - [ ] Ensure it accepts variant_id in order items
  - [ ] Ensure validation errors are returned properly
- [ ] 3.2.2 Test API endpoints manually
  - [ ] POST /api/orders - with single-size items
  - [ ] POST /api/orders - with multi-size items and variant_id
  - [ ] POST /api/orders - with missing variant_id (error)
  - [ ] POST /api/orders - with invalid variant_id (error)

### 3.3 API Documentation
- [ ] 3.3.1 Update API documentation
  - [ ] Document has_variants field
  - [ ] Document variants array structure
  - [ ] Document variant_id in order items
  - [ ] Add example requests/responses

## Phase 4: Frontend Data Layer (Day 6)

### 4.1 Update Menu Store
- [ ] 4.1.1 Verify menuStore in `frontend/src/stores/menu.js`
  - [ ] Ensure createMenuItem accepts has_variants and variants
  - [ ] Ensure updateMenuItem accepts has_variants and variants
  - [ ] Ensure items are stored with variants
- [ ] 4.1.2 Add helper methods if needed
  - [ ] getVariantById(itemId, variantId)
  - [ ] getDefaultVariant(itemId)

### 4.2 Update Order Store
- [ ] 4.2.1 Update addItem in `frontend/src/stores/order.js`
  - [ ] Accept optional variant parameter
  - [ ] If variant provided, add variant_id and variant_name
  - [ ] Use variant.price if variant provided
  - [ ] Use item.price if no variant
- [ ] 4.2.2 Update editItem
  - [ ] Handle variant changes
- [ ] 4.2.3 Write tests for order store
  - [ ] Test addItem - single-size
  - [ ] Test addItem - with variant

## Phase 5: Frontend UI (Days 7-9)

### 5.1 Update MenuView Display
- [ ] 5.1.1 Update MenuView.vue template
  - [ ] Add conditional rendering for single-size items
  - [ ] Add conditional rendering for multi-size items
  - [ ] Display variants list for multi-size items
  - [ ] Each variant shows name and price
  - [ ] Each variant has "Thêm" button
- [ ] 5.1.2 Update MenuView.vue script
  - [ ] Update addToOrder method to accept variant
  - [ ] Handle single-size item click
  - [ ] Handle multi-size variant click
- [ ] 5.1.3 Style variants display
  - [ ] Mobile-friendly layout
  - [ ] Clear visual separation between variants
  - [ ] Responsive design

### 5.2 Update MenuView Create/Edit Form
- [ ] 5.2.1 Add has_variants checkbox to form
  - [ ] Toggle between single-size and multi-size sections
- [ ] 5.2.2 Create single-size form section
  - [ ] Price input
  - [ ] Ingredients selector (existing)
- [ ] 5.2.3 Create multi-size form section
  - [ ] Variants array form
  - [ ] Each variant has: ID, Name, Price inputs
  - [ ] Each variant has: Is Default checkbox
  - [ ] Each variant has: Ingredients selector
  - [ ] Each variant has: Remove button
  - [ ] Add Variant button
- [ ] 5.2.4 Implement form logic
  - [ ] addVariant() method
  - [ ] removeVariant(index) method
  - [ ] Ensure at least one default when adding first variant
  - [ ] Ensure at least one default when removing variants
  - [ ] Handle ingredient selection per variant
- [ ] 5.2.5 Add form validation
  - [ ] If has_variants, must have at least 1 variant
  - [ ] If has_variants, must have exactly 1 default
  - [ ] Variant IDs must be unique
  - [ ] Variant prices must be > 0
  - [ ] Show validation errors
- [ ] 5.2.6 Update saveItem method
  - [ ] Validate before save
  - [ ] Handle single-size save
  - [ ] Handle multi-size save
- [ ] 5.2.7 Style the form
  - [ ] Mobile-friendly layout
  - [ ] Clear visual hierarchy
  - [ ] Responsive design
  - [ ] Loading states
  - [ ] Error states

### 5.3 Update Order Views
- [ ] 5.3.1 Update WaiterView.vue
  - [ ] Display variant_name in order items
  - [ ] Format: "Item Name (Variant Name)"
  - [ ] Handle items without variants
- [ ] 5.3.2 Update CashierView.vue
  - [ ] Display variant_name in order items
  - [ ] Receipt shows variant names
- [ ] 5.3.3 Update BaristaView.vue (if exists)
  - [ ] Display variant_name in order queue
  - [ ] Show size clearly for preparation

### 5.4 Mobile Optimization
- [ ] 5.4.1 Test on mobile devices
  - [ ] Variants display correctly
  - [ ] Form is usable on mobile
  - [ ] Touch targets are at least 44x44px
  - [ ] No horizontal scrolling
- [ ] 5.4.2 Optimize for small screens
  - [ ] Vertical layout for variants
  - [ ] Collapsible sections if needed
  - [ ] Bottom sheet for variant selection

## Phase 6: Testing (Days 10-11)

### 6.1 Backend Unit Tests
- [ ] 6.1.1 Run all backend tests
  ```bash
  cd backend
  go test ./domain/menu/... -v
  go test ./application/services/... -v
  ```
- [ ] 6.1.2 Check test coverage
  ```bash
  go test ./... -cover
  ```
- [ ] 6.1.3 Ensure coverage > 80%
- [ ] 6.1.4 Fix failing tests

### 6.2 Frontend Unit Tests
- [ ] 6.2.1 Write component tests
  - [ ] MenuView displays single-size items
  - [ ] MenuView displays multi-size items
  - [ ] MenuView form toggles variants
  - [ ] MenuView form validates variants
  - [ ] OrderStore adds item with variant
- [ ] 6.2.2 Run frontend tests
  ```bash
  cd frontend
  npm run test
  ```
- [ ] 6.2.3 Fix failing tests

### 6.3 Integration Tests
- [ ] 6.3.1 Test complete flows
  - [ ] Create single-size item → Display → Order → Receipt
  - [ ] Create multi-size item → Display → Order each variant → Receipt
  - [ ] Edit item: single → multi-size
  - [ ] Edit item: multi → single-size
  - [ ] Delete item with variants
  - [ ] Cost calculation per variant
- [ ] 6.3.2 Test error cases
  - [ ] Order multi-size without variant_id
  - [ ] Order with invalid variant_id
  - [ ] Create item with no default variant
  - [ ] Create item with duplicate variant IDs
- [ ] 6.3.3 Test edge cases
  - [ ] Item with 1 variant
  - [ ] Item with 10 variants (max)
  - [ ] Toggle has_variants multiple times
  - [ ] Delete default variant (should reassign)

### 6.4 Manual Testing
- [ ] 6.4.1 Test on development environment
  - [ ] Create test menu items
  - [ ] Create test orders
  - [ ] Verify all flows work
- [ ] 6.4.2 Test on mobile devices
  - [ ] iOS Safari
  - [ ] Android Chrome
  - [ ] Responsive design
- [ ] 6.4.3 Performance testing
  - [ ] Menu list loads < 500ms
  - [ ] Create item < 1s
  - [ ] Order creation < 1s
  - [ ] Cost calculation < 2s per item

### 6.5 Bug Fixes
- [ ] 6.5.1 Document all bugs found
- [ ] 6.5.2 Prioritize bugs (critical, high, medium, low)
- [ ] 6.5.3 Fix critical and high priority bugs
- [ ] 6.5.4 Retest after fixes

## Phase 7: Documentation & Deployment (Day 11)

### 7.1 Documentation
- [ ] 7.1.1 Update user documentation
  - [ ] How to create single-size items
  - [ ] How to create multi-size items
  - [ ] How to order items with variants
  - [ ] Screenshots and examples
- [ ] 7.1.2 Update developer documentation
  - [ ] API changes
  - [ ] Data model changes
  - [ ] Code examples
- [ ] 7.1.3 Create migration guide (if needed in future)

### 7.2 Deployment Preparation
- [ ] 7.2.1 Create deployment checklist
- [ ] 7.2.2 Backup database (even though empty)
- [ ] 7.2.3 Prepare rollback plan
- [ ] 7.2.4 Set up monitoring
  - [ ] API response times
  - [ ] Error rates
  - [ ] Menu item creation rate

### 7.3 Staging Deployment
- [ ] 7.3.1 Deploy to staging environment
  ```bash
  git checkout main
  git pull
  git merge feature/menu-size-variants
  # Deploy to staging
  ```
- [ ] 7.3.2 Run smoke tests on staging
  - [ ] Create menu items
  - [ ] Create orders
  - [ ] Verify all features work
- [ ] 7.3.3 Performance testing on staging
- [ ] 7.3.4 Fix any issues found

### 7.4 Production Deployment
- [ ] 7.4.1 Final review
  - [ ] All tests passing
  - [ ] All acceptance criteria met
  - [ ] Documentation complete
  - [ ] Team approval
- [ ] 7.4.2 Deploy to production
  ```bash
  # Deploy backend
  cd backend
  make deploy-production
  
  # Deploy frontend
  cd frontend
  npm run build
  # Deploy build
  ```
- [ ] 7.4.3 Monitor production
  - [ ] Check error logs
  - [ ] Check API response times
  - [ ] Check user activity
- [ ] 7.4.4 Verify production deployment
  - [ ] Create test menu item
  - [ ] Create test order
  - [ ] Delete test data

### 7.5 Post-Deployment
- [ ] 7.5.1 Monitor for 24 hours
  - [ ] Check error rates
  - [ ] Check performance metrics
  - [ ] Check user feedback
- [ ] 7.5.2 Address any issues immediately
- [ ] 7.5.3 Collect user feedback
- [ ] 7.5.4 Plan improvements based on feedback

## Summary

**Total Tasks**: ~150 tasks
**Estimated Time**: 11 days (2 weeks)
**Team Size**: 1-2 developers

**Critical Path**:
1. Backend Domain (Days 1-2) - Foundation
2. Backend Service (Days 3-4) - Business logic
3. Backend API (Day 5) - Endpoints
4. Frontend Data (Day 6) - State management
5. Frontend UI (Days 7-9) - User interface
6. Testing (Days 10-11) - Quality assurance
7. Deployment (Day 11) - Go live

**Success Criteria**:
- [ ] All tasks completed
- [ ] All tests passing
- [ ] All acceptance criteria met
- [ ] Zero critical bugs
- [ ] Performance requirements met
- [ ] Deployed to production
- [ ] User feedback positive
