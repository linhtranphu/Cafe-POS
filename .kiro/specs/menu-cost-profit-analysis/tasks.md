# Implementation Plan: Menu Cost & Profit Analysis

## Overview

Implementation plan cho tính năng Menu Cost & Profit Analysis, bao gồm:
- Backend: Go services, API endpoints, database schema changes
- Frontend: Vue.js components, views, API integration
- Testing: Unit tests và property-based tests cho correctness properties

Implementation được chia thành các phases để đảm bảo incremental progress và early validation.

## Tasks

- [ ] 1. Backend Foundation - Data Models và Database Schema
  - [x] 1.1 Extend MenuItem model với cost tracking fields
    - Add `current_cost`, `cost_last_calculated_at`, `cost_status` fields to MenuItem struct
    - Add `CostStatus` enum type (FINAL, ESTIMATED, INCOMPLETE)
    - Update MongoDB schema and create indexes
    - _Requirements: 1.1, 1.2, 1.5_
  
  - [x] 1.2 Create OrderItem collection và model
    - Create separate `order_items` collection schema
    - Define OrderItem struct with `accounting_cost`, `cost_calculated_at`, `cost_status` fields
    - Create indexes for efficient querying (order_id, menu_item_id, cost_status)
    - _Requirements: 5.1, 5.2, 5.3_
  
  - [x] 1.3 Extend Ingredient model với conversion và wastage fields
    - Add `conversion_rate` (default: 1.0) and `wastage_percentage` (default: 0.0) fields
    - Update existing ingredients with default values
    - _Requirements: 10.1, 10.2_
  
  - [x] 1.4 Create OperatingExpense model và collection
    - Define OperatingExpense struct with period and expense breakdown fields
    - Create `operating_expenses` collection with indexes
    - Define OperatingExpenseRequest DTO for API validation
    - _Requirements: 6.5.1, 6.5.2_
  
  - [x] 1.5 Extend ShopSettings model với low_margin_threshold
    - Add `low_margin_threshold` field (default: 20.0)
    - Update settings migration script
    - _Requirements: 3.2, 3.3_

- [ ] 2. Backend Core Services - Cost Calculator Service
  - [x] 2.1 Implement CalculateMenuItemCost method
    - Fetch menu item ingredients
    - Calculate cost using formula: sum(quantity * cost_per_unit * conversion_rate * (1 + wastage_percentage/100))
    - Handle missing cost_per_unit (mark as INCOMPLETE)
    - Round to 2 decimal places
    - _Requirements: 1.1, 1.2, 1.5, 1.7_
  
  - [x] 2.2 Write property test for cost calculation formula
    - *Property 1: Cost Calculation Formula*
    - *Validates: Requirements 1.1, 1.2, 1.7, 10.1, 10.2, 10.4*
  
  - [x] 2.3 Implement CalculateAllMenuItemCosts method
    - Batch fetch all menu items
    - Calculate cost for each item
    - Update menu_items collection with current_cost
    - Return summary statistics
    - _Requirements: 1.1, 1.2_
  
  - [x] 2.4 Implement CalculateShiftOrderCosts method
    - Fetch all orders in shift
    - For each order item, calculate accounting_cost using current ingredient costs
    - Save to order_items collection with cost_status = FINAL
    - Handle incomplete costs gracefully
    - _Requirements: 5.2, 5.3, 5.4_
  
  - [x] 2.5 Write property test for shift closure cost calculation
    - **Property 5: Shift Closure Cost Calculation**
    - **Validates: Requirements 5.2, 5.3**
  
  - [x] 2.6 Implement QueueCostRecalculation method
    - Find all menu items using specific ingredient
    - Queue background job for each menu item
    - Use Go channels
    - _Requirements: 1.3, 9.1_
  
  - [x] 2.7 Write property test for background job queuing
    - **Property 3: Background Job Queuing on Ingredient Update**
    - **Validates: Requirements 1.3, 9.1**

- [-] 3. Backend Core Services - Profit Analyzer Service
  - [x] 3.1 Implement CalculateMenuItemProfit method
    - Calculate profit_margin = ((price - cost) / price) * 100
    - Calculate absolute_profit = price - cost
    - Handle edge cases (price = 0, cost > price)
    - Round to 2 decimal places
    - _Requirements: 2.1, 2.5, 2.6_
  
  - [x] 3.2 Write property test for profit calculations
    - **Property 2: Profit Calculations**
    - **Validates: Requirements 2.1, 2.5, 2.6**
  
  - [x] 3.3 Implement DetectWarningStatus method
    - Check if cost > price → mark as "loss"
    - Check if profit_margin < threshold → mark as "low_margin"
    - Otherwise mark as "none"
    - _Requirements: 3.1, 3.2, 3.6_
  
  - [x] 3.4 Write property tests for warning detection
    - **Property 10: Loss Detection**
    - **Property 11: Low Margin Detection**
    - **Property 12: Warning Status Transitions**
    - **Validates: Requirements 3.1, 3.2, 3.6**
  
  - [x] 3.5 Implement GetAllMenuItemProfits method
    - Fetch all menu items with costs
    - Calculate profit metrics for each
    - Apply filters (category, sort)
    - Return with summary statistics
    - _Requirements: 2.1, 2.9, 4.3, 4.4_
  
  - [x] 3.6 Write property tests for filtering and sorting
    - **Property 14: Category Filtering**
    - **Property 15: Profit Margin Sorting**
    - **Validates: Requirements 4.3, 4.4**
  
  - [x] 3.7 Implement GetCategoryProfits method
    - Aggregate orders by category within date range
    - Calculate total_revenue, total_cost (using accounting_cost), total_profit
    - Calculate average_profit_margin
    - _Requirements: 6.1, 6.2, 6.3, 6.4_
  
  - [x] 3.8 Write property test for category profit aggregation
    - **Property 7: Category Profit Aggregation**
    - **Validates: Requirements 6.1, 6.2, 6.3**
  
  - [x] 3.9 Implement GetOperatingProfit method
    - Calculate gross_profit from orders
    - Fetch operating expenses for date range
    - Allocate expenses if needed (monthly → daily)
    - Calculate operating_profit = gross_profit - expenses
    - _Requirements: 6.5.1, 6.5.3, 6.5.4, 6.5.6_
  
  - [x] 3.10 Write property test for operating profit calculations
    - **Property 8: Operating Profit Calculations**
    - **Validates: Requirements 6.5.1, 6.5.3, 6.5.4**

- [x] 4. Backend Supporting Services
  - [x] 4.1 Implement CostRecalculationService
    - Create worker pool for processing jobs
    - Implement ProcessRecalculationQueue method
    - Implement GetRecalculationStatus method
    - Handle timeouts and retries
    - _Requirements: 9.2, 9.3, 9.4_
  
  - [x] 4.2 Implement OperatingExpenseService
    - Implement UpsertOperatingExpense method with validation
    - Implement GetOperatingExpenseForDate method
    - Implement GetOperatingExpenses method with date range filter
    - Implement AllocateDailyExpense method for proportional allocation
    - _Requirements: 6.5.2, 6.5.7, 6.5.8_
  
  - [x] 4.3 Write property test for expense allocation
    - **Property 9: Expense Allocation**
    - **Validates: Requirements 6.5.8**
  
  - [x] 4.4 Write unit tests for edge cases
    - Test incomplete ingredient data handling
    - Test zero price edge case
    - Test negative profit scenarios
    - Test empty date ranges
    - _Requirements: 1.5, 1.6, 2.9_

- [x] 5. Checkpoint - Backend Core Logic Complete
  - Ensure all backend services compile and unit tests pass
  - Verify cost calculation logic with sample data
  - Ask the user if questions arise

- [x] 6. Backend API Endpoints - Menu Cost APIs
  - [x] 6.1 Implement GET /api/menu/costs endpoint
    - Create handler with query params (category, sort_by, sort_order)
    - Call GetAllMenuItemProfits service
    - Return items with summary and recalculation_status
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 7.4_
  
  - [x] 6.2 Implement GET /api/menu/costs/:id endpoint
    - Fetch menu item with ingredient breakdown
    - Calculate cost detail for each ingredient
    - Return menu item, ingredients array, and total_cost
    - _Requirements: 8.1, 8.2, 8.3_
  
  - [x] 6.3 Implement GET /api/menu/warnings endpoint
    - Accept optional threshold query param
    - Call DetectWarnings service method
    - Return loss_items, low_margin_items, and counts
    - _Requirements: 3.3, 3.4, 3.5_
  
  - [x] 6.4 Write integration tests for menu cost APIs
    - Test GET /api/menu/costs with various filters
    - Test GET /api/menu/costs/:id with valid and invalid IDs
    - Test GET /api/menu/warnings with custom threshold
    - _Requirements: 4.1, 4.2, 8.1_

- [x] 7. Backend API Endpoints - Profit Analysis APIs
  - [x] 7.1 Implement GET /api/reports/category-profit endpoint
    - Accept start_date and end_date query params
    - Validate date range
    - Call GetCategoryProfits service
    - Return categories array with date_range
    - _Requirements: 6.1, 6.4, 7.1_
  
  - [x] 7.2 Write property test for date range filtering
    - **Property 16: Date Range Filtering**
    - **Validates: Requirements 6.4, 6.5.6**
  
  - [x] 7.3 Implement GET /api/reports/operating-profit endpoint
    - Accept start_date and end_date query params
    - Call GetOperatingProfit service
    - Handle missing expenses gracefully
    - Return full operating profit breakdown
    - _Requirements: 6.5.1, 6.5.6, 6.5.9_
  
  - [x] 7.4 Write integration tests for profit analysis APIs
    - Test category profit with various date ranges
    - Test operating profit with and without expenses
    - Test expense allocation scenarios
    - _Requirements: 6.1, 6.5.1_

- [x] 8. Backend API Endpoints - Operating Expense APIs
  - [x] 8.1 Implement POST /api/operating-expenses endpoint
    - Validate request body (dates, amounts >= 0)
    - Call UpsertOperatingExpense service
    - Auto-calculate total_expenses
    - Return created/updated expense
    - _Requirements: 6.5.2, 6.5.7_
  
  - [x] 8.2 Implement GET /api/operating-expenses endpoint
    - Accept optional start_date and end_date query params
    - Call GetOperatingExpenses service
    - Return expenses array
    - _Requirements: 6.5.7_
  
  - [x] 8.3 Write unit tests for operating expense APIs
    - Test validation errors (invalid dates, negative amounts)
    - Test upsert behavior (create vs update)
    - Test date range filtering
    - _Requirements: 6.5.2, 6.5.7_

- [x] 9. Backend API Endpoints - Modified Endpoints
  - [x] 9.1 Modify POST /api/shifts/:id/close endpoint
    - After closing shift, call CalculateShiftOrderCosts
    - Return shift data + cost_calculation summary
    - Handle cost calculation errors gracefully
    - _Requirements: 5.1, 5.2, 5.5_
  
  - [x] 9.2 Write property test for accounting cost immutability
    - **Property 6: Accounting Cost Immutability**
    - **Validates: Requirements 5.8, 9.6**
  
  - [x] 9.3 Modify PATCH /api/settings endpoint
    - Accept low_margin_threshold in request body
    - Validate threshold >= 0
    - Update shop settings
    - _Requirements: 3.3_
  
  - [x] 9.4 Write integration test for shift closure workflow
    - Create shift with orders
    - Close shift
    - Verify all order items have accounting_cost
    - Verify cost_status = FINAL
    - _Requirements: 5.1, 5.2, 5.3_

- [x] 10. Checkpoint - Backend APIs Complete
  - Ensure all API endpoints work correctly
  - Test with Postman or curl
  - Verify database operations
  - Ask the user if questions arise

- [x] 11. Frontend Foundation - API Client và Types
  - [x] 11.1 Create TypeScript types for cost and profit data
    - Define MenuItemCost, CategoryProfit, OperatingProfitReport interfaces
    - Define OperatingExpense, DateRange, ProfitFilter interfaces
    - Define RecalculationStatus, WarningStatus types
    - _Requirements: 4.1, 6.1, 6.5.1_
  
  - [x] 11.2 Create API client methods for menu cost endpoints
    - Implement getMenuCosts(filter: ProfitFilter)
    - Implement getMenuCostDetail(id: string)
    - Implement getMenuWarnings(threshold?: number)
    - _Requirements: 4.1, 8.1, 3.3_
  
  - [x] 11.3 Create API client methods for profit analysis endpoints
    - Implement getCategoryProfit(dateRange: DateRange)
    - Implement getOperatingProfit(dateRange: DateRange)
    - _Requirements: 6.1, 6.5.1_
  
  - [x] 11.4 Create API client methods for operating expense endpoints
    - Implement createOperatingExpense(data: OperatingExpense)
    - Implement getOperatingExpenses(dateRange?: DateRange)
    - _Requirements: 6.5.2, 6.5.7_

- [x] 12. Frontend Components - MenuCostView
  - [x] 12.1 Create MenuCostView component structure
    - Setup component with data fetching on mount
    - Implement loading and error states
    - Create table layout for menu items
    - _Requirements: 4.1, 7.1_
  
  - [x] 12.2 Implement menu cost table with columns
    - Display name, category, price, current_cost, profit_margin, absolute_profit
    - Implement color coding (green/yellow/red/gray based on warning_status)
    - Add cost_status indicator
    - _Requirements: 4.1, 7.2, 7.3_
  
  - [x] 12.3 Implement filtering and sorting
    - Add category filter dropdown
    - Add sort by dropdown (profit_margin, absolute_profit, name)
    - Add sort order toggle (asc/desc)
    - _Requirements: 4.3, 4.4_
  
  - [x] 12.4 Implement summary statistics section
    - Display total_items, loss_count, low_margin_count
    - Display average_profit_margin
    - Add recalculation status indicator
    - _Requirements: 7.4, 9.5_
  
  - [x] 12.5 Implement row click to show cost breakdown
    - Open MenuItemCostBreakdown component in modal/drawer
    - Pass menu_item_id as prop
    - _Requirements: 8.1_
  
  - [x] 12.6 Write unit tests for MenuCostView
    - Test component rendering with mock data
    - Test filtering by category
    - Test sorting functionality
    - Test warning color coding
    - _Requirements: 4.1, 4.3, 4.4_

- [x] 13. Frontend Components - MenuItemCostBreakdown
  - [x] 13.1 Create MenuItemCostBreakdown component
    - Setup modal/drawer layout
    - Fetch cost detail on mount
    - Display loading state
    - _Requirements: 8.1_
  
  - [x] 13.2 Implement ingredient breakdown table
    - Display columns: name, quantity, unit, cost_per_unit, total_cost
    - Show conversion_rate and wastage_percentage if non-default
    - Highlight ingredients with missing cost_per_unit
    - _Requirements: 8.2, 8.3, 8.4_
  
  - [x] 13.3 Implement total cost summary
    - Display total_cost at bottom
    - Show warning if any ingredient has incomplete cost
    - _Requirements: 8.3_
  
  - [x] 13.4 Write unit tests for MenuItemCostBreakdown
    - Test rendering with complete data
    - Test warning display for missing costs
    - Test conversion rate and wastage display
    - _Requirements: 8.1, 8.2, 8.3_

- [x] 14. Frontend Components - ProfitAnalysisView
  - [x] 14.1 Create ProfitAnalysisView component structure
    - Setup component with date range picker
    - Implement view mode toggle (category vs operating)
    - Add loading and error states
    - _Requirements: 6.1, 6.5.1, 7.1_
  
  - [x] 14.2 Implement category profit view
    - Create table with columns: category, revenue, cost, profit, margin
    - Display order_count and item_count
    - Add date range display
    - _Requirements: 6.1, 6.4, 7.1_
  
  - [x] 14.3 Implement operating profit view
    - Display gross profit section (revenue, COGS, gross profit, margin)
    - Display expenses breakdown (staff, rent, utilities, marketing, other)
    - Display operating profit section (total expenses, operating profit, margin)
    - Show expense_allocated indicator and note if applicable
    - _Requirements: 6.5.1, 6.5.3, 6.5.4, 6.5.9_
  
  - [x] 14.4 Implement date range picker
    - Add preset options (today, this week, this month)
    - Add custom date range selector
    - Trigger data refresh on date change
    - _Requirements: 6.4, 6.5.6_
  
  - [x] 14.5 Write unit tests for ProfitAnalysisView
    - Test category view rendering
    - Test operating profit view rendering
    - Test date range filtering
    - Test view mode toggle
    - _Requirements: 6.1, 6.5.1_

- [x] 15. Frontend Components - OperatingExpenseForm
  - [x] 15.1 Create OperatingExpenseForm component
    - Setup form with period date pickers
    - Add input fields for all expense types
    - Implement auto-calculate total_expenses
    - _Requirements: 6.5.2_
  
  - [x] 15.2 Implement form validation
    - Validate period_start <= period_end
    - Validate all amounts >= 0
    - Display validation errors
    - _Requirements: 6.5.2_
  
  - [x] 15.3 Implement save and cancel actions
    - Call createOperatingExpense API on save
    - Handle success and error responses
    - Emit events to parent component
    - _Requirements: 6.5.2_
  
  - [x] 15.4 Write unit tests for OperatingExpenseForm
    - Test form validation
    - Test total calculation
    - Test save action
    - Test error handling
    - _Requirements: 6.5.2_

- [x] 16. Frontend Integration - Navigation và Routes
  - [x] 16.1 Add menu cost routes to router
    - Add /manager/menu-costs route for MenuCostView
    - Add /manager/profit-analysis route for ProfitAnalysisView
    - Add route guards for manager role
    - _Requirements: 4.1, 6.1_
  
  - [x] 16.2 Add navigation items to manager menu
    - Add "Chi phí & Lợi nhuận" menu item
    - Add sub-items for "Chi phí món" and "Phân tích lợi nhuận"
    - Update BottomNav for mobile
    - _Requirements: 4.1, 6.1_
  
  - [x] 16.3 Add operating expense management to settings
    - Add "Chi phí vận hành" section in settings view
    - Integrate OperatingExpenseForm component
    - Display list of existing expenses
    - _Requirements: 6.5.2, 6.5.7_

- [x] 17. Frontend Polish - Responsive Design và UX
  - [x] 17.1 Implement responsive design for MenuCostView
    - Optimize table layout for mobile (card view)
    - Adjust filter and sort controls for small screens
    - Test on various screen sizes
    - _Requirements: 4.1, 7.1_
  
  - [x] 17.2 Implement responsive design for ProfitAnalysisView
    - Optimize tables for mobile
    - Stack sections vertically on small screens
    - Adjust date picker for mobile
    - _Requirements: 6.1, 6.5.1_
  
  - [x] 17.3 Add loading skeletons and empty states
    - Add skeleton loaders for tables
    - Add empty state messages when no data
    - Add error state with retry button
    - _Requirements: 4.1, 6.1_
  
  - [x] 17.4 Implement number formatting
    - Format currency values with thousand separators
    - Format percentages with 2 decimal places
    - Use Vietnamese locale formatting
    - _Requirements: 4.1, 6.1, 7.2_

- [x] 18. Checkpoint - Frontend Complete
  - Ensure all components render correctly
  - Test user flows end-to-end
  - Verify responsive design on mobile and desktop
  - Ask the user if questions arise

- [x] 19. Data Migration và Backfill
  - [x] 19.1 Create migration script for schema changes
    - Add new fields to menu_items collection
    - Create order_items collection
    - Create operating_expenses collection
    - Add indexes
    - _Requirements: 1.1, 5.1, 6.5.1_
  
  - [x] 19.2 Backfill current_cost for existing menu items
    - Calculate current_cost for all menu items
    - Set cost_status based on ingredient data completeness
    - Update cost_last_calculated_at
    - _Requirements: 1.1, 1.2_
  
  - [x] 19.3 Backfill accounting_cost for historical orders
    - For closed shifts, calculate accounting_cost using current ingredient costs
    - Mark as cost_status = ESTIMATED (not FINAL)
    - Add note indicating backfilled data
    - _Requirements: 5.1, 5.2_
  
  - [x] 19.4 Write migration verification tests
    - Verify all menu items have current_cost
    - Verify all historical orders have accounting_cost
    - Verify indexes are created
    - _Requirements: 1.1, 5.1_

- [x] 20. Background Job Setup
  - [x] 20.1 Setup cost recalculation worker
    - Initialize worker pool on server start
    - Configure queue (Go channels)
    - Implement graceful shutdown
    - _Requirements: 9.1, 9.2, 9.3_
  
  - [x] 20.2 Integrate recalculation trigger on ingredient update
    - Hook into ingredient update handler
    - Call QueueCostRecalculation after successful update
    - Handle errors gracefully
    - _Requirements: 1.3, 9.1_
  
  - [x] 20.3 Write property test for batch recalculation optimization
    - **Property 17: Batch Recalculation Optimization**
    - **Validates: Requirements 9.2**
  
  - [x] 20.4 Write integration test for background job processing
    - Update ingredient cost
    - Verify job is queued
    - Wait for job completion
    - Verify menu item costs are updated
    - _Requirements: 1.3, 9.1, 9.3_

- [x] 21. Final Integration Testing
  - [x] 21.1 Test complete shift closure workflow
    - Create shift with orders
    - Close shift
    - Verify accounting_cost is calculated
    - Update ingredient costs
    - Verify accounting_cost remains unchanged
    - _Requirements: 5.1, 5.2, 5.8, 9.6_
  
  - [x] 21.2 Test complete profit analysis workflow
    - Create orders with various menu items
    - Close shift
    - View category profit report
    - Add operating expenses
    - View operating profit report
    - _Requirements: 6.1, 6.5.1_
  
  - [x] 21.3 Test warning detection workflow
    - Create menu items with various cost/price ratios
    - View menu cost list with warnings
    - Adjust low margin threshold
    - Verify warning status updates
    - _Requirements: 3.1, 3.2, 3.3, 3.6_
  
  - [x] 21.4 Write property test for recipe change immutability
    - **Property 18: Recipe Change Immutability**
    - **Validates: Requirements 8.6**
  
  - [x] 21.5 Write property test for summary statistics
    - **Property 19: Summary Statistics Calculation**
    - **Validates: Requirements 7.4**

- [ ] 22. Performance Testing và Optimization
  - [ ] 22.1 Run performance benchmarks
    - Benchmark cost calculation for 1000+ menu items
    - Benchmark shift closure with 500+ orders
    - Benchmark category profit aggregation for 10,000+ orders
    - _Requirements: 9.3_
  
  - [ ]* 22.2 Optimize slow queries
    - Add missing indexes if needed
    - Optimize aggregation pipelines
    - Implement caching where appropriate
    - _Requirements: 9.3_
  
  - [ ]* 22.3 Test frontend performance
    - Test table rendering with 100+ items
    - Test filter and sort responsiveness
    - Optimize re-renders if needed
    - _Requirements: 4.1_

- [ ] 23. Documentation và Deployment
  - [ ] 23.1 Write API documentation
    - Document all new endpoints with examples
    - Document request/response schemas
    - Document error codes and messages
    - _Requirements: 4.1, 6.1, 6.5.1_
  
  - [x] 23.2 Write user guide
    - Document how to view menu costs
    - Document how to analyze profits
    - Document how to input operating expenses
    - Document warning indicators
    - _Requirements: 4.1, 6.1, 6.5.2_
  
  - [x] 23.3 Setup monitoring and alerts
    - Configure metrics collection
    - Setup alerts for job failures
    - Setup alerts for high error rates
    - _Requirements: 9.3_
  
  - [ ] 23.4 Deploy to production
    - Run migration scripts
    - Deploy backend services
    - Deploy frontend build
    - Enable background workers
    - Monitor for errors
    - _Requirements: All_

- [ ] 24. Final Checkpoint - Feature Complete
  - Verify all requirements are implemented
  - Ensure all tests pass
  - Verify production deployment is stable
  - Ask the user if questions arise

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation at key milestones
- Property tests validate universal correctness properties (minimum 100 iterations each)
- Unit tests validate specific examples and edge cases
- Backend uses Go with MongoDB, Frontend uses Vue.js with TypeScript
- Background jobs use Go channels or Redis for queueing
- Property-based testing uses gopter library for Go
