# Implementation Plan: Expense Auto-Aggregation

## Overview

This implementation plan breaks down the expense auto-aggregation feature into discrete coding tasks. The approach follows a bottom-up strategy: first extending the data model and backend API, then building the frontend components, and finally integrating everything together. Each task includes property-based tests and unit tests to ensure correctness.

## Tasks

- [ ] 1. Extend Expense model with expense_type field
  - Modify the Expense struct in Go to include expense_type field
  - Add validation constants for the five expense types (staff_salary, rent, utilities, marketing_costs, other_expenses)
  - Implement Validate() method to check expense_type is required and valid
  - Create MongoDB index on (shop_id, date, expense_type) for efficient queries
  - _Requirements: 1.1, 1.3, 2.1, 2.2, 2.5_

- [ ]* 1.1 Write property test for expense validation
  - **Property 1: Expense validation rejects missing expense_type**
  - **Property 3: Expense type validation rejects invalid values**
  - **Validates: Requirements 1.3, 2.2**

- [ ]* 1.2 Write unit tests for expense model
  - Test validation with missing expense_type
  - Test validation with invalid expense_type values
  - Test validation with valid expense_type values
  - _Requirements: 1.3, 2.2_

- [ ] 2. Implement expense aggregation service
  - [ ] 2.1 Create aggregation response models (ExpenseAggregation, ExpenseBreakdownItem, AggregationResponse)
    - Define Go structs for aggregation response structure
    - _Requirements: 3.3, 3.4_
  
  - [ ] 2.2 Implement AggregateByDateRange method in ExpenseService
    - Query expenses filtered by shop_id, date range, and expense_type existence
    - Group expenses by expense_type and calculate totals
    - Build breakdown arrays with individual expense details
    - Handle timezone considerations for date boundaries
    - _Requirements: 3.2, 3.7, 7.3, 7.4, 7.5_
  
  - [ ]* 2.3 Write property test for aggregation correctness
    - **Property 4: Aggregation correctness - sum of breakdown equals total**
    - **Property 5: Aggregation response completeness**
    - **Property 9: Date field consistency in filtering**
    - **Property 13: Total expense count accuracy**
    - **Property 14: Expenses without expense_type excluded from aggregation**
    - **Validates: Requirements 3.2, 3.3, 3.4, 3.7, 7.4, 9.4, 10.3**
  
  - [ ]* 2.4 Write unit tests for aggregation service
    - Test empty date range returns zero values
    - Test boundary dates are included (start_date and end_date)
    - Test expenses without expense_type are excluded
    - Test timezone handling for date boundaries
    - _Requirements: 3.5, 7.3, 10.3_

- [ ] 3. Create aggregation API endpoint
  - [ ] 3.1 Implement AggregateExpenses handler in ExpenseHandler
    - Parse and validate start_date and end_date query parameters
    - Validate date range (end_date not before start_date)
    - Call ExpenseService.AggregateByDateRange
    - Return JSON response with aggregation results
    - Handle errors with appropriate HTTP status codes
    - _Requirements: 3.1, 3.2, 3.6, 7.1_
  
  - [ ]* 3.2 Write property test for date range validation
    - **Property 6: Date range validation**
    - **Validates: Requirements 3.6, 7.1**
  
  - [ ]* 3.3 Write unit tests for aggregation API
    - Test successful aggregation with valid date range
    - Test error response for invalid date format
    - Test error response for end_date before start_date
    - Test error response for missing parameters
    - _Requirements: 3.1, 3.6, 7.1, 7.2_
  
  - [ ] 3.4 Register API route in main.go
    - Add GET /api/expenses/aggregate route
    - _Requirements: 3.1_

- [ ] 4. Implement operating expenses update API
  - [ ] 4.1 Create UpdateOperatingExpensesRequest struct
    - Define request model with five expense type fields
    - _Requirements: 6.2_
  
  - [ ] 4.2 Implement UpdateOperatingExpenses handler in ShopSettingsHandler
    - Parse request body with expense amounts
    - Call ShopSettingsService to update operating expenses
    - Return success or error response
    - _Requirements: 4.10, 6.1, 6.3, 6.4_
  
  - [ ] 4.3 Implement UpdateOperatingExpenses method in ShopSettingsService
    - Update shop_settings document with new operating expense values
    - Map expense_type values to correct shop_settings fields
    - _Requirements: 6.2_
  
  - [ ]* 4.4 Write property test for operating expense save
    - **Property 7: Operating expense save with correct field mapping**
    - **Property 2: Expense persistence round-trip** (for shop settings)
    - **Validates: Requirements 4.10, 6.1, 6.2**
  
  - [ ]* 4.5 Write unit tests for operating expense update
    - Test successful update with valid amounts
    - Test error handling for invalid shop_id
    - Test error handling for database failures
    - _Requirements: 6.3, 6.4_
  
  - [ ] 4.6 Register API route in main.go
    - Add PUT /api/shop-settings/operating-expenses route
    - _Requirements: 4.10_

- [ ] 5. Checkpoint - Backend API complete
  - Ensure all backend tests pass
  - Manually test aggregation API with curl or Postman
  - Verify database indexes are created
  - Ask the user if questions arise

- [ ] 6. Modify ExpenseForm component to include expense_type selector
  - [ ] 6.1 Add expense_type field to expense data model in component
    - Add expense_type to reactive data
    - _Requirements: 1.1_
  
  - [ ] 6.2 Add expense_type select dropdown to template
    - Create select element with five options
    - Use Vietnamese labels for options
    - Mark field as required
    - _Requirements: 1.1, 1.2, 1.5_
  
  - [ ] 6.3 Add client-side validation for expense_type
    - Validate expense_type is selected before save
    - Display error message if not selected
    - _Requirements: 1.3_
  
  - [ ]* 6.4 Write unit tests for ExpenseForm
    - Test expense_type field is rendered
    - Test five options are present with correct values
    - Test Vietnamese labels are displayed
    - Test validation error when expense_type is empty
    - _Requirements: 1.1, 1.2, 1.3, 1.5_

- [ ] 7. Create OperatingExpenseAggregator component
  - [ ] 7.1 Create component file and basic structure
    - Set up Vue component with template, script, and style
    - Add reactive data for startDate, endDate, loading, errors
    - _Requirements: 4.1, 4.2_
  
  - [ ] 7.2 Implement date range selector UI
    - Add start_date and end_date input fields
    - Add "Tổng hợp" (Aggregate) button
    - _Requirements: 4.2, 4.3_
  
  - [ ] 7.3 Implement date range validation
    - Validate dates are not empty
    - Validate end_date is not before start_date
    - Display validation error messages
    - _Requirements: 7.1, 7.2_
  
  - [ ] 7.4 Implement aggregate() method
    - Call GET /api/expenses/aggregate with date range
    - Handle loading state
    - Handle API errors and display error messages
    - Store aggregation results in component state
    - Initialize editableAmounts with aggregated totals
    - _Requirements: 4.4, 8.1, 8.3_
  
  - [ ]* 7.5 Write unit tests for OperatingExpenseAggregator
    - Test date range validation
    - Test API call with correct parameters
    - Test error handling for API failures
    - Test loading state during aggregation
    - _Requirements: 4.4, 7.1, 7.2, 8.1, 9.1_

- [ ] 8. Create AggregationResultsView component
  - [ ] 8.1 Create component file and basic structure
    - Set up Vue component to receive results as prop
    - _Requirements: 4.5_
  
  - [ ] 8.2 Implement results table UI
    - Display table with five rows (one per expense_type)
    - Show Vietnamese labels for expense types
    - Display count and total for each expense_type
    - Add editable input fields for amounts
    - Add "Xem chi tiết" (View details) button for each row
    - Display date range and total expense count in header
    - _Requirements: 4.5, 4.8, 9.3, 9.4_
  
  - [ ] 8.3 Implement amount editing functionality
    - Bind input fields to editableAmounts
    - Update totals when amounts are edited
    - _Requirements: 4.8, 6.5_
  
  - [ ] 8.4 Implement "Save to Operating Expenses" button
    - Add button to trigger save operation
    - Emit save event to parent component
    - _Requirements: 4.9_
  
  - [ ] 8.5 Add expense type label mapping utility
    - Create method to map expense_type codes to Vietnamese labels
    - _Requirements: 4.5_
  
  - [ ]* 8.6 Write unit tests for AggregationResultsView
    - Test table renders with five rows
    - Test Vietnamese labels are displayed correctly
    - Test editable inputs are bound to amounts
    - Test "Save" button emits event
    - Test total calculations
    - _Requirements: 4.5, 4.8, 4.9, 9.3, 9.4_

- [ ] 9. Create BreakdownDetailModal component
  - [ ] 9.1 Create modal component structure
    - Set up modal overlay and content
    - Add close button
    - _Requirements: 4.6, 5.1_
  
  - [ ] 9.2 Implement breakdown table UI
    - Display table with columns: date, description, amount
    - Show expense_type label in modal header
    - Display sum of breakdown items in footer
    - Add link to view full expense record
    - _Requirements: 4.7, 5.2, 5.3, 5.4, 5.5_
  
  - [ ] 9.3 Add date and currency formatting
    - Format dates in Vietnamese locale
    - Format amounts as VND currency
    - _Requirements: 5.3_
  
  - [ ]* 9.4 Write unit tests for BreakdownDetailModal
    - Test modal renders with breakdown data
    - Test table displays all breakdown items
    - Test sum calculation matches total
    - Test close button emits close event
    - Test link to expense record is correct
    - _Requirements: 4.7, 5.2, 5.3, 5.4, 5.5_

- [ ] 10. Integrate components in OperatingExpenseAggregator
  - [ ] 10.1 Add AggregationResultsView to template
    - Conditionally render when aggregationResults exists
    - Pass results and editableAmounts as props
    - Listen to save and view-breakdown events
    - _Requirements: 4.4, 4.5_
  
  - [ ] 10.2 Add BreakdownDetailModal to template
    - Conditionally render when selectedBreakdown exists
    - Pass breakdown data as prop
    - Listen to close event
    - _Requirements: 4.6, 5.1_
  
  - [ ] 10.3 Implement showBreakdown method
    - Find aggregation by expense_type
    - Set selectedBreakdown to show modal
    - _Requirements: 4.6, 5.2_
  
  - [ ] 10.4 Implement saveToOperatingExpenses method
    - Call PUT /api/shop-settings/operating-expenses
    - Send editableAmounts (preserving edited values)
    - Display success toast on success
    - Display error toast on failure
    - Preserve aggregated data on failure for retry
    - Emit saved event to parent on success
    - _Requirements: 4.10, 6.1, 6.3, 6.4, 6.6, 8.4_
  
  - [ ]* 10.5 Write property test for edited values preservation
    - **Property 8: Edited values preserved for saving**
    - **Validates: Requirements 6.6**
  
  - [ ]* 10.6 Write integration tests for full aggregation flow
    - Test complete flow: select dates → aggregate → view breakdown → edit → save
    - Test state preservation on save failure
    - Test retry after failure
    - _Requirements: 4.4, 4.10, 6.6, 8.4_

- [ ] 11. Add OperatingExpenseAggregator to Settings page
  - [ ] 11.1 Import and register OperatingExpenseAggregator component
    - Add component to Settings/Operating Expenses view
    - _Requirements: 4.1_
  
  - [ ] 11.2 Add "Auto-aggregate from Expenses" button or section
    - Add button to trigger aggregator display
    - Or embed aggregator directly in the page
    - _Requirements: 4.1_
  
  - [ ] 11.3 Handle saved event from aggregator
    - Refresh operating expenses display after save
    - _Requirements: 6.3_
  
  - [ ]* 11.4 Write end-to-end test for Settings integration
    - Test navigation to Settings → Operating Expenses
    - Test aggregator button is visible
    - Test complete aggregation and save flow
    - _Requirements: 4.1, 4.10_

- [ ] 12. Checkpoint - Frontend components complete
  - Ensure all frontend tests pass
  - Manually test UI in browser
  - Verify Vietnamese labels display correctly
  - Test responsive design on mobile
  - Ask the user if questions arise

- [ ] 13. Handle existing expenses without expense_type
  - [ ] 13.1 Update expense list view to show indicator for missing expense_type
    - Add badge or icon for expenses without expense_type
    - _Requirements: 10.4_
  
  - [ ] 13.2 Ensure expense edit form allows adding expense_type to existing expenses
    - Test editing existing expense and adding expense_type
    - _Requirements: 10.2_
  
  - [ ]* 13.3 Write unit tests for backward compatibility
    - Test retrieving expense without expense_type doesn't error
    - Test editing expense to add expense_type
    - Test aggregation excludes expenses without expense_type
    - _Requirements: 10.1, 10.2, 10.3_

- [ ] 14. Add error handling and user feedback
  - [ ] 14.1 Add loading indicators during API calls
    - Show spinner during aggregation
    - Show spinner during save operation
    - _Requirements: 9.1_
  
  - [ ] 14.2 Add success and error toast notifications
    - Success message after save: "Đã cập nhật chi phí hoạt động"
    - Error message for API failures: "Không thể tổng hợp chi phí"
    - Error message for save failures: "Không thể lưu chi phí hoạt động"
    - Info message for empty results: "Không tìm thấy chi phí trong khoảng thời gian này"
    - _Requirements: 6.3, 6.4, 8.1, 8.2, 9.2_
  
  - [ ] 14.3 Add retry functionality for failed operations
    - Allow retry button for failed aggregation
    - Allow retry button for failed save
    - _Requirements: 8.3_
  
  - [ ]* 14.4 Write unit tests for error handling
    - Test error messages display correctly
    - Test retry functionality
    - Test loading indicators
    - _Requirements: 8.1, 8.3, 9.1_

- [ ] 15. Implement concurrency control for operating expense updates
  - [ ] 15.1 Add version field or timestamp to shop_settings
    - Add updated_at timestamp to shop_settings model
    - _Requirements: 8.5_
  
  - [ ] 15.2 Implement optimistic locking in update operation
    - Check updated_at before update
    - Return conflict error if version mismatch
    - _Requirements: 8.5_
  
  - [ ]* 15.3 Write property test for concurrent update safety
    - **Property 12: Concurrent update safety**
    - **Validates: Requirements 8.5**
  
  - [ ]* 15.4 Write unit tests for concurrency control
    - Test concurrent updates with version checking
    - Test conflict detection and error response
    - _Requirements: 8.5_

- [ ] 16. Final integration and testing
  - [ ] 16.1 Run all property-based tests
    - Verify all 14 properties pass with 100+ iterations
    - _All requirements_
  
  - [ ] 16.2 Run all unit tests
    - Verify backend and frontend unit tests pass
    - _All requirements_
  
  - [ ] 16.3 Perform manual end-to-end testing
    - Create expenses with all five expense_types
    - Test aggregation with various date ranges
    - Test editing amounts before save
    - Test viewing breakdown details
    - Test saving to operating expenses
    - Verify operating expenses update in shop settings
    - Test with existing expenses without expense_type
    - _All requirements_
  
  - [ ] 16.4 Test error scenarios
    - Test with invalid date ranges
    - Test with network failures
    - Test with concurrent updates
    - _Requirements: 3.6, 7.1, 8.1, 8.3, 8.5_

- [ ] 17. Final checkpoint - Feature complete
  - Ensure all tests pass
  - Verify all requirements are met
  - Review code for security and performance
  - Update documentation if needed
  - Ask the user if questions arise

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Property tests validate universal correctness properties with 100+ iterations
- Unit tests validate specific examples, edge cases, and UI interactions
- Checkpoints ensure incremental validation and allow for user feedback
- The implementation follows a bottom-up approach: data model → backend API → frontend components → integration
- Vietnamese labels and localization are integrated throughout the UI components
