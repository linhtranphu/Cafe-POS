# Menu Size Variants - Implementation Plan

## Overview

This implementation plan converts the menu size variants design into actionable coding tasks. The feature allows menu items to have multiple size variants (e.g., M, L, XL) with different prices and ingredient quantities, while maintaining full backward compatibility with single-size items.

**Key Principles**:
- Additive changes only - no breaking changes
- Backward compatibility with existing single-size items
- Clear separation between single-size and multi-size modes
- Comprehensive testing at every step
- Cost tracking per variant for profit analysis

## Tasks

### 1. Backend Domain Layer - MenuItem Model

- [x] 1.1 Update MenuItem struct with variant support
  - Add `HasVariants bool` field (default: false)
  - Add `Variants []MenuItemVariant` field (optional, omitempty)
  - Keep existing `Price`, `Ingredients`, `CurrentCost` fields for backward compatibility
  - Add MenuItemVariant struct with ID, Name, Price, Ingredients, Available, IsDefault, CurrentCost, CostStatus, CostLastCalculatedAt
  - _Requirements: FR-1.1, FR-1.2, FR-1.3, FR-1.4, FR-2.1-FR-2.7_

- [x] 1.2 Implement MenuItem helper methods
  - Implement `GetDefaultVariant()` - returns first variant with IsDefault=true
  - Implement `GetVariantByID(variantID string)` - returns variant or nil
  - Implement `GetPrice(variantID string)` - handles both single and multi-size
  - Implement `GetIngredients(variantID string)` - handles both single and multi-size
  - _Requirements: FR-5.1, FR-5.2, FR-5.4_

- [x] 1.3 Implement MenuItem validation logic
  - Implement `Validate()` method with comprehensive rules
  - If has_variants=true: require at least 1 variant, exactly 1 default, unique IDs
  - If has_variants=true: ensure Price and Ingredients are NOT set (prevent ambiguity)
  - If has_variants=false: require Price > 0, ensure Variants is empty
  - Validate variant prices > 0 and variant IDs are non-empty
  - _Requirements: FR-4.1-FR-4.7_

- [x]* 1.4 Write unit tests for MenuItem domain model
  - Test GetDefaultVariant() with valid/invalid data
  - Test GetVariantByID() with found/not found cases
  - Test GetPrice() for both single-size and multi-size items
  - Test Validate() for all validation rules (no variants, no default, duplicates, ambiguous states)
  - Test backward compatibility - single-size items work as before
  - _Requirements: FR-1.1-FR-4.7_

### 2. Backend Domain Layer - OrderItem Model

- [x] 2.1 Update OrderItem struct with variant support
  - Add `VariantID string` field (optional)
  - Add `VariantName string` field (optional)
  - Keep existing Price, Quantity, Subtotal fields
  - _Requirements: FR-3.1, FR-3.2, FR-3.3_

- [x]* 2.2 Write unit tests for OrderItem
  - Test OrderItem with variant (new functionality)
  - Test OrderItem without variant (backward compatible)
  - Test JSON marshaling/unmarshaling
  - Test subtotal calculation with variants
  - _Requirements: FR-3.1-FR-3.3_

### 3. Backend Service Layer - MenuService

- [x] 3.1 Update MenuService.CreateMenuItem
  - Handle has_variants=false: set Price and Ingredients (existing behavior)
  - Handle has_variants=true: set Variants array, leave Price/Ingredients empty
  - Call item.Validate() before saving
  - Return clear error messages for validation failures
  - _Requirements: FR-4.1-FR-4.7, AC-1.4, AC-2.7_

- [x] 3.2 Update MenuService.UpdateMenuItem
  - Handle toggling has_variants between true/false
  - When changing false→true: clear Price and Ingredients
  - When changing true→false: clear Variants
  - Validate all changes before saving
  - _Requirements: AC-3.2, AC-3.3, AC-3.5_

- [x]* 3.3 Write unit tests for MenuService
  - Test CreateMenuItem for single-size (backward compatible)
  - Test CreateMenuItem for multi-size with valid variants
  - Test CreateMenuItem with invalid data (no variants, no default, duplicates)
  - Test UpdateMenuItem toggling between single and multi-size
  - Test UpdateMenuItem updating variants
  - _Requirements: FR-4.1-FR-4.7, AC-1.4, AC-2.7, AC-3.5_

### 4. Backend Service Layer - OrderService

- [x] 4.1 Update OrderService.CreateOrder with variant support
  - Check menuItem.HasVariants to determine flow
  - If has_variants=false: use menuItem.Price (existing behavior, variant_id optional/ignored)
  - If has_variants=true: require variant_id, validate it exists, use variant.Price
  - Populate variant_name in order item
  - Return clear error if variant_id missing or invalid for multi-size item
  - _Requirements: FR-5.1-FR-5.5, AC-5.1-AC-5.3, AC-6.1-AC-6.5_

- [x]* 4.2 Write unit tests for OrderService
  - Test CreateOrder with single-size item (backward compatible)
  - Test CreateOrder with multi-size item and valid variant_id
  - Test CreateOrder with multi-size item without variant_id (error)
  - Test CreateOrder with invalid variant_id (error)
  - Test CreateOrder with mixed single and multi-size items
  - _Requirements: FR-5.1-FR-5.5, AC-5.1-AC-6.5_

### 5. Backend Service Layer - CostCalculatorService

- [x] 5.1 Update CostCalculatorService.CalculateMenuItemCost
  - Check menuItem.HasVariants to determine flow
  - If has_variants=false: calculate cost from menuItem.Ingredients (existing behavior)
  - If has_variants=true: loop through variants, calculate cost per variant
  - Update variant.CurrentCost, variant.CostStatus, variant.CostLastCalculatedAt
  - Clear menuItem.CurrentCost when has_variants=true (prevent confusion)
  - Reuse existing calculateIngredientsCost() method
  - _Requirements: FR-6.1-FR-6.10, AC-11.1-AC-11.7_

- [x]* 5.2 Write unit tests for CostCalculatorService
  - Test CalculateMenuItemCost for single-size (backward compatible)
  - Test CalculateMenuItemCost for multi-size with variants
  - Test cost calculation with conversion rates and wastage
  - Verify formula: quantity × cost_per_unit × conversion_rate × (1 + wastage/100)
  - Verify each variant has independent cost
  - _Requirements: FR-6.1-FR-6.10, AC-11.1-AC-11.7_

### 6. Backend API Layer - Handlers

- [x] 6.1 Verify MenuHandler endpoints accept variant fields
  - Ensure CreateMenuItem accepts has_variants and variants fields
  - Ensure UpdateMenuItem accepts has_variants and variants fields
  - Return 400 for validation errors with clear messages
  - Return 201/200 for successful operations
  - _Requirements: FR-7.1-FR-7.5_

- [x] 6.2 Verify OrderHandler accepts variant_id in order items
  - Ensure CreateOrder accepts variant_id in order items
  - Return 400 for missing/invalid variant_id with clear messages
  - Return 201 for successful order creation
  - _Requirements: FR-8.1-FR-8.4_

- [x] 6.3 Implement cost analysis API endpoints
  - Implement GET /api/menu/:id/cost-breakdown - detailed cost breakdown per variant
  - Response includes: ingredient costs, conversion rates, wastage, total cost per variant
  - Implement GET /api/menu/:id/profit-analysis - profit analysis per variant
  - Response includes: price, cost, profit, profit margin % per variant
  - Implement POST /api/menu/:id/calculate-cost - trigger cost calculation
  - _Requirements: FR-7.6, FR-9.1-FR-9.4, AC-10.1-AC-10.5, AC-12.1-AC-12.4_

- [x]* 6.4 Write API integration tests for menu endpoints
  - Test POST /api/menu for single-size (201 Created)
  - Test POST /api/menu for multi-size (201 Created)
  - Test POST /api/menu with invalid data (400 Bad Request)
  - Test GET /api/menu - list items with variants and costs (200 OK)
  - Test PUT /api/menu/:id - toggle single to multi (200 OK)
  - Test DELETE /api/menu/:id for both types (200 OK)
  - _Requirements: FR-7.1-FR-7.5_

- [x]* 6.5 Write API integration tests for order endpoints
  - Test POST /api/orders with single-size items (201 Created)
  - Test POST /api/orders with multi-size items and variant_id (201 Created)
  - Test POST /api/orders with mixed items (201 Created)
  - Test POST /api/orders missing variant_id for multi-size (400 Bad Request)
  - Test POST /api/orders with invalid variant_id (400 Bad Request)
  - _Requirements: FR-8.1-FR-8.4_

- [x]* 6.6 Write API integration tests for cost analysis endpoints
  - Test GET /api/menu/:id/cost-breakdown for single-size item
  - Test GET /api/menu/:id/cost-breakdown for multi-size item (per variant)
  - Test GET /api/menu/:id/profit-analysis for both types
  - Test POST /api/menu/:id/calculate-cost triggers recalculation
  - Verify response times < 500ms for cost-breakdown and profit-analysis
  - _Requirements: FR-9.1-FR-9.4, NFR-1.5, NFR-1.6_

### 7. Checkpoint - Backend Complete

- [x] 7.1 Run all backend tests and verify coverage
  - Run `go test ./... -v -cover`
  - Verify test coverage > 80%
  - Verify all tests passing
  - Fix any failing tests
  - _Requirements: NFR-7.1, NFR-7.2_

### 8. Frontend Data Layer - Stores

- [x] 8.1 Verify menuStore handles variant fields
  - Ensure createMenuItem accepts has_variants and variants
  - Ensure updateMenuItem accepts has_variants and variants
  - Ensure items are stored with variants
  - Add helper methods: getVariantById, getDefaultVariant
  - _Requirements: AC-1.1-AC-4.5_

- [x] 8.2 Update orderStore.addItem to handle variants
  - Accept optional variant parameter
  - If variant provided: add variant_id and variant_name, use variant.price
  - If no variant: use item.price (backward compatible)
  - _Requirements: AC-5.1-AC-6.6_

- [ ]* 8.3 Write unit tests for stores
  - Test menuStore createMenuItem for both types
  - Test menuStore updateMenuItem toggling modes
  - Test orderStore addItem with/without variant
  - Test orderStore calculateTotal with mixed items
  - _Requirements: AC-1.1-AC-6.6_

### 9. Frontend UI - MenuView Display

- [x] 9.1 Update MenuView template for variant display
  - Add conditional rendering for single-size items (show price, add button)
  - Add conditional rendering for multi-size items (show variants list)
  - Display each variant with name, price, and add button
  - Mobile-friendly layout with clear visual separation
  - _Requirements: AC-4.1-AC-4.3, NFR-3.2, NFR-3.5_

- [x] 9.2 Update MenuView script to handle variant selection
  - Update addToOrder method to accept variant parameter
  - Handle single-size item click (immediate add)
  - Handle multi-size variant click (add with variant)
  - _Requirements: AC-5.1-AC-6.3_

- [ ]* 9.3 Write unit tests for MenuView component
  - Test single-size item rendering
  - Test multi-size item rendering with variants
  - Test variant selection and add to order
  - Test mobile responsiveness
  - _Requirements: AC-4.1-AC-6.3, NFR-3.2, NFR-3.5_

### 10. Frontend UI - MenuForm (Create/Edit)

- [x] 10.1 Add has_variants toggle to form
  - Add checkbox "Món có nhiều size"
  - Toggle between single-size and multi-size sections
  - _Requirements: AC-1.1, AC-2.1_

- [x] 10.2 Implement single-size form section
  - Show price input when has_variants=false
  - Show ingredients selector
  - _Requirements: AC-1.2, AC-1.3_

- [x] 10.3 Implement multi-size form section
  - Show variants array form when has_variants=true
  - Each variant has: ID, Name, Price inputs
  - Each variant has: Is Default checkbox, Ingredients selector, Remove button
  - Add "Thêm size" button
  - _Requirements: AC-2.2-AC-2.6_

- [x] 10.4 Implement form logic and validation
  - Implement addVariant() and removeVariant() methods
  - Ensure at least 1 variant when has_variants=true
  - Ensure exactly 1 default variant
  - Validate variant IDs are unique
  - Show validation errors clearly
  - _Requirements: FR-4.1-FR-4.7, AC-2.3-AC-2.6_

- [x] 10.5 Update saveItem method
  - Validate before save
  - Handle single-size save (has_variants=false)
  - Handle multi-size save (has_variants=true)
  - _Requirements: AC-1.4, AC-2.7, AC-3.5_

- [ ]* 10.6 Write unit tests for MenuForm component
  - Test form toggle between single and multi-size
  - Test variant add/remove functionality
  - Test validation rules
  - Test save for both types
  - _Requirements: AC-1.1-AC-3.6, FR-4.1-FR-4.7_

### 11. Frontend UI - Order Views

- [x] 11.1 Update WaiterView to display variant names
  - Display format: "Item Name (Variant Name)" for multi-size
  - Display format: "Item Name" for single-size
  - _Requirements: AC-7.1, AC-7.2_

- [x] 11.2 Update CashierView to display variant names
  - Show variant names in order list
  - Show variant names in receipt
  - _Requirements: AC-8.1, AC-8.3_

- [x] 11.3 Update BaristaView to display variant names
  - Show variant names in order queue
  - Show size clearly (M/L/XL) for preparation
  - Show ingredients list matching variant
  - _Requirements: AC-9.1, AC-9.2, AC-9.3_

- [ ]* 11.4 Write unit tests for order view components
  - Test WaiterView displays variants correctly
  - Test CashierView displays variants in receipt
  - Test BaristaView displays variant details
  - _Requirements: AC-7.1-AC-9.3_

### 12. Frontend UI - Cost Analysis Views

- [x] 12.1 Create CostAnalysisView component for managers
  - Display menu items with cost data per variant
  - Show current_cost, cost_status, cost_last_calculated_at for each variant
  - Show profit margin (price - cost) per variant
  - Allow filtering by cost_status (FINAL/ESTIMATED/INCOMPLETE)
  - _Requirements: AC-10.1-AC-10.3, AC-12.1_

- [x] 12.2 Implement cost breakdown modal
  - Show detailed ingredient costs per variant
  - Display conversion rates and wastage percentages
  - Show formula breakdown: quantity × cost_per_unit × conversion_rate × (1 + wastage/100)
  - Display total cost calculation
  - _Requirements: AC-10.4, AC-11.1-AC-11.5_

- [x] 12.3 Implement profit comparison view
  - Display all variants side-by-side with costs
  - Show cost difference between sizes
  - Show profit margin difference between sizes
  - Highlight most profitable variant
  - _Requirements: AC-12.2-AC-12.4_

- [x]* 12.4 Write unit tests for cost analysis components
  - Test CostAnalysisView renders correctly with variant data
  - Test cost breakdown modal displays formula correctly
  - Test profit comparison calculations
  - Test filtering by cost_status
  - _Requirements: AC-10.1-AC-12.4_

### 13. Checkpoint - Frontend Complete

- [x] 13.1 Run all frontend tests
  - Run `npm run test`
  - Verify all tests passing
  - Fix any failing tests
  - _Requirements: NFR-7.1, NFR-7.2_

### 14. Integration Testing

- [x] 14.1 Test complete single-size item flow
  - Create single-size item via API
  - Display in MenuView
  - Add to order (no variant_id)
  - Verify order has correct price
  - Calculate cost
  - Display in receipt
  - _Requirements: AC-1.1-AC-1.5, AC-5.1-AC-5.4_

- [x] 14.2 Test complete multi-size item flow
  - Create multi-size item via API
  - Display in MenuView with variants
  - Add to order with variant_id
  - Verify order has variant_name and correct price
  - Calculate cost per variant
  - Display in receipt with variant name
  - _Requirements: AC-2.1-AC-2.8, AC-6.1-AC-6.6_

- [x] 14.3 Test mixed orders
  - Create order with both single-size and multi-size items
  - Verify all items priced correctly
  - Verify total calculation correct
  - Verify receipt displays correctly
  - _Requirements: AC-5.1-AC-6.6, AC-8.1-AC-8.4_

- [x] 14.4 Test toggling between modes
  - Create single-size item, edit to multi-size
  - Verify old price cleared, variants saved
  - Edit back to single-size
  - Verify variants cleared, price saved
  - _Requirements: AC-3.1-AC-3.6_

- [x] 14.5 Test cost analysis flow
  - Create multi-size item with variants
  - Calculate costs for all variants
  - View cost breakdown per variant
  - View profit analysis comparing variants
  - Verify cost_status updates correctly (FINAL/INCOMPLETE)
  - Update ingredient price, verify costs recalculate
  - _Requirements: AC-10.1-AC-12.4, FR-6.6_

### 15. Error Handling Testing

- [x] 15.1 Test validation errors
  - Create multi-size without variants (should fail)
  - Create multi-size without default variant (should fail)
  - Create multi-size with duplicate variant IDs (should fail)
  - Create item with both price and variants (should fail)
  - Order multi-size without variant_id (should fail with clear error)
  - Order with invalid variant_id (should fail with clear error)
  - _Requirements: FR-4.1-FR-4.7, FR-8.3, FR-8.4_

- [x] 15.2 Test edge cases
  - Item with 1 variant
  - Item with 10 variants (max)
  - Toggle has_variants multiple times
  - Delete default variant (should reassign or fail)
  - _Requirements: NFR-2.1_

- [x] 15.3 Test cost calculation edge cases
  - Calculate cost with missing ingredient data (should set INCOMPLETE status)
  - Calculate cost with zero wastage percentage
  - Calculate cost with high wastage percentage (>50%)
  - Calculate cost with unit conversion (g to kg, ml to L)
  - Verify cost accuracy to 2 decimal places
  - _Requirements: FR-6.4-FR-6.9, NFR-5.4-NFR-5.7_

### 16. Backward Compatibility Testing

- [x] 16.1 Verify single-size items work exactly as before
  - Create single-size item
  - Verify has_variants=false
  - Verify price field populated
  - Order without variant_id
  - Calculate cost
  - Verify no regressions
  - _Requirements: FR-1.5, FR-5.1, FR-6.1_

- [ ] 16.2 Verify no breaking changes to existing features
  - Test menu categories still work
  - Test ingredient management still works
  - Test order editing still works
  - Test payment processing still works
  - Test shift management still works
  - _Requirements: NFR-5.1, NFR-5.2_

### 17. Performance Testing

- [ ] 17.1 Test API response times
  - GET /api/menu with 10-15 items < 500ms
  - POST /api/menu < 1s
  - POST /api/orders < 1s
  - POST /api/menu/:id/calculate-cost < 2s
  - GET /api/menu/:id/cost-breakdown < 500ms
  - GET /api/menu/:id/profit-analysis < 500ms
  - _Requirements: NFR-1.1, NFR-1.2, NFR-1.3, NFR-1.4, NFR-1.5, NFR-1.6_

- [ ] 17.2 Test with realistic data volume
  - Create 10-15 single-size items
  - Create 10-15 multi-size items (2-3 variants each)
  - Verify menu list loads quickly
  - Create 20-30 orders
  - Verify order creation works smoothly
  - Calculate costs for all items (100+ ingredients total)
  - _Requirements: NFR-2.1, NFR-2.2, NFR-2.4_

### 18. Documentation and Deployment

- [ ] 18.1 Update API documentation
  - Document variant fields in MenuItem schema
  - Document variant_id in OrderItem schema
  - Document cost analysis endpoints
  - Add request/response examples for all endpoints
  - _Requirements: NFR-7.3_

- [ ] 18.2 Update user documentation
  - Create guide for creating multi-size items
  - Create guide for ordering with variants
  - Create guide for viewing cost analysis
  - Add screenshots and examples
  - _Requirements: NFR-7.3_

- [ ] 18.3 Deploy to staging
  - Deploy backend and frontend to staging
  - Run smoke tests on staging
  - Verify all features work
  - Fix any issues found
  - _Requirements: All NFRs_

- [ ] 18.4 Deploy to production
  - Final review and team approval
  - Deploy to production
  - Monitor error logs and performance
  - Verify production deployment
  - _Requirements: All NFRs_

- [ ] 18.5 Post-deployment monitoring
  - Monitor for 24 hours
  - Check error rates and performance metrics
  - Collect user feedback
  - Address any issues immediately
  - _Requirements: NFR-5.1, NFR-6.1, NFR-6.2_

## Notes

- Tasks marked with `*` are optional test tasks that can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Backward compatibility is verified at every step
- All validation errors must have clear, user-friendly messages
- Cost analysis features are fully integrated for profit visibility

## Success Criteria

- All acceptance criteria met (AC-1.1 through AC-12.4)
- All functional requirements implemented (FR-1.1 through FR-9.4)
- All non-functional requirements met (NFR-1.1 through NFR-7.4)
- Test coverage > 80%
- Zero critical bugs
- No regressions in existing features
- Positive user feedback
- Cost analysis features fully functional
- Profit margins visible for all variants
- Cost calculation accuracy: 100% (no rounding errors beyond 2 decimals)
