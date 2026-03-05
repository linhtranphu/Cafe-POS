# Implementation Plan: Fund-Expense Integration

## Overview

Tính năng này tích hợp hệ thống quản lý quỹ với các module chi tiêu (Expense và Ingredient). Implementation sẽ được chia thành các phases tuần tự: Database schema, Backend domain models và services, API endpoints, Frontend UI, và Testing.

Scope: Requirements 1, 2, 4, 6, 7, 8, 9, 10, 11 (loại bỏ Requirement 3 về Facility và Requirement 5 về Báo cáo).

## Tasks

- [ ] 1. Database Schema và Migration
  - [ ] 1.1 Tạo migration script cho expenses collection
    - Thêm field `paid_from_fund` (boolean)
    - Thêm field `fund_transaction_id` (ObjectId)
    - Tạo index cho `paid_from_fund` và `fund_transaction_id`
    - _Requirements: 8.1, 8.2, 8.6_

  - [ ] 1.2 Tạo migration script cho fund_transactions collection
    - Thêm field `source_type` (string: "expense", "ingredient")
    - Thêm field `source_id` (ObjectId)
    - Tạo unique compound index cho `source_type` + `source_id` (sparse)
    - _Requirements: 8.3, 8.4, 8.5_

  - [ ] 1.3 Tạo ingredient_restock_history collection
    - Định nghĩa schema với các field: ingredient_id, quantity, cost_per_unit, total_cost
    - Thêm field `paid_from_fund`, `expense_id`, `fund_transaction_id`
    - Thêm field audit: performed_by, performed_by_name, reason, created_at
    - Tạo index cho `ingredient_id` + `created_at` và `fund_transaction_id`
    - _Requirements: 2.1, 2.3, 2.4, 2.5, 2.6, 2.7_

- [ ] 2. Backend Domain Models
  - [ ] 2.1 Cập nhật Expense domain model
    - Thêm field `PaidFromFund bool`
    - Thêm field `FundTransactionID primitive.ObjectID`
    - Thêm validation logic cho fund payment consistency
    - _Requirements: 8.1, 8.2, 10.1_

  - [ ] 2.2 Cập nhật FundTransaction domain model
    - Thêm field `SourceType string`
    - Thêm field `SourceID primitive.ObjectID`
    - Thêm validation cho source_type values ("expense", "ingredient")
    - _Requirements: 8.3, 8.4, 8.5_

  - [ ] 2.3 Tạo IngredientRestockRecord domain model
    - Định nghĩa struct với tất cả field cần thiết
    - Implement validation methods
    - _Requirements: 2.1, 2.3, 2.4, 2.5_

- [ ] 3. Backend Repositories
  - [ ] 3.1 Cập nhật ExpenseRepository
    - Thêm method `FindByFundTransactionID(ctx, fundTxID)`
    - Thêm method `FindPaidFromFund(ctx, filter)` với pagination
    - Cập nhật Create/Update methods để handle fund fields
    - _Requirements: 1.5, 4.3, 4.4_

  - [ ] 3.2 Cập nhật FundRepository
    - Thêm method `FindBySource(ctx, sourceType, sourceID)`
    - Thêm method `CreateWithdrawalWithSource(ctx, withdrawal, sourceType, sourceID)`
    - Implement duplicate prevention logic
    - _Requirements: 1.6, 8.5, 11.2_

  - [ ] 3.3 Tạo IngredientRestockRepository
    - Implement Create method với MongoDB session support
    - Implement FindByIngredientID method với pagination
    - Implement FindByFundTransactionID method
    - _Requirements: 2.5, 2.6_

- [ ] 4. Core Service: FundExpenseIntegrationService
  - [ ] 4.1 Implement ValidateFundBalance method
    - Kiểm tra fund balance đủ cho withdrawal amount
    - Return error với thông tin balance hiện tại và amount yêu cầu
    - _Requirements: 6.1, 6.2, 6.3_

  - [ ] 4.2 Implement CreateExpenseFromFund method
    - Validate fund balance trước khi tạo
    - Sử dụng MongoDB transaction để đảm bảo atomicity
    - Tạo expense record với paid_from_fund=true
    - Tạo fund withdrawal transaction với source linking
    - Update fund balance
    - Implement rollback logic nếu có lỗi
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 9.1, 9.2, 10.1_

  - [ ]* 4.3 Write property test cho CreateExpenseFromFund
    - **Property 1: Fund Balance Validation**
    - **Property 2: Withdrawal Amount Consistency**
    - **Property 3: Balance Update Invariant**
    - **Property 4: Bidirectional Linking - Expense to Fund**
    - **Property 11: Transaction Atomicity - Expense Creation Failure**
    - **Property 12: Transaction Atomicity - Fund Transaction Rollback**
    - **Validates: Requirements 1.2, 1.3, 1.4, 1.5, 1.6, 6.2, 6.3, 9.1, 9.2**

  - [ ] 4.4 Implement RestockIngredientFromFund method
    - Validate fund balance trước khi restock
    - Sử dụng MongoDB transaction cho atomicity
    - Tạo ingredient restock record
    - Tạo expense record với category "ingredient purchase"
    - Tạo fund withdrawal transaction với source linking
    - Update ingredient stock quantity
    - Update fund balance
    - Implement rollback logic cho tất cả operations
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 9.3_

  - [ ]* 4.5 Write property test cho RestockIngredientFromFund
    - **Property 5: Bidirectional Linking - Ingredient to Fund**
    - **Property 6: Expense Creation for Ingredient Restock**
    - **Property 7: Stock Quantity Update**
    - **Property 13: Transaction Atomicity - Ingredient Restock Rollback**
    - **Validates: Requirements 2.3, 2.4, 2.5, 2.6, 2.7, 9.3**

  - [ ] 4.6 Implement GetExpensesPaidFromFund method
    - Query expenses với paid_from_fund=true filter
    - Support pagination và sorting
    - _Requirements: 4.3, 4.4_

  - [ ]* 4.7 Write unit tests cho FundExpenseIntegrationService
    - Test insufficient balance scenarios
    - Test validation errors
    - Test concurrent withdrawal attempts
    - Test edge cases (zero balance, exact match)

- [ ] 5. Checkpoint - Backend Core Logic Complete
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 6. API Endpoints - Expense
  - [ ] 6.1 Implement POST /api/expenses/from-fund endpoint
    - Parse request body với expense details
    - Extract user info từ JWT token
    - Call FundExpenseIntegrationService.CreateExpenseFromFund
    - Return expense và fund transaction details
    - Handle errors (insufficient balance, validation errors)
    - _Requirements: 1.1, 1.2, 1.3, 6.4_

  - [ ] 6.2 Cập nhật GET /api/expenses endpoint
    - Thêm query parameter `paid_from_fund` filter
    - Return expenses với fund transaction info nếu có
    - _Requirements: 4.3, 4.4_

  - [ ] 6.3 Cập nhật GET /api/expenses/:id endpoint
    - Include fund_transaction_id trong response
    - Populate fund transaction details nếu có
    - _Requirements: 4.2_

  - [ ]* 6.4 Write integration tests cho expense endpoints
    - Test POST /api/expenses/from-fund với valid data
    - Test POST /api/expenses/from-fund với insufficient balance
    - Test GET /api/expenses với paid_from_fund filter
    - Test GET /api/expenses/:id với fund-paid expense

- [ ] 7. API Endpoints - Ingredient
  - [ ] 7.1 Implement POST /api/ingredients/:id/restock/from-fund endpoint
    - Parse request body với restock details (quantity, cost_per_unit, reason)
    - Extract user info từ JWT token
    - Call FundExpenseIntegrationService.RestockIngredientFromFund
    - Return restock record, expense, và fund transaction details
    - Handle errors (insufficient balance, validation errors)
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 6.5_

  - [ ] 7.2 Implement GET /api/ingredients/:id/restock-history endpoint
    - Query ingredient_restock_history collection
    - Include fund transaction info cho fund-paid restocks
    - Support pagination
    - _Requirements: 2.5, 2.6, 2.7_

  - [ ]* 7.3 Write integration tests cho ingredient endpoints
    - Test POST /api/ingredients/:id/restock/from-fund với valid data
    - Test POST /api/ingredients/:id/restock/from-fund với insufficient balance
    - Test GET /api/ingredients/:id/restock-history

- [ ] 8. API Endpoints - Fund
  - [ ] 8.1 Cập nhật GET /api/fund/transactions endpoint
    - Include source_type và source_id trong response
    - Support filter by source_type
    - Populate source record details (expense/ingredient info)
    - _Requirements: 7.7, 8.3, 8.4_

  - [ ] 8.2 Cập nhật GET /api/fund/transactions/:id endpoint
    - Include full source record details
    - Provide navigation link to source record
    - _Requirements: 7.7, 7.8_

  - [ ]* 8.3 Write integration tests cho fund endpoints
    - Test GET /api/fund/transactions với source filtering
    - Test GET /api/fund/transactions/:id với source details

- [ ] 9. Checkpoint - Backend API Complete
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 10. Frontend Services
  - [ ] 10.1 Tạo fundExpenseService.js
    - Implement `createExpenseFromFund(expenseData)` API call
    - Implement `getExpensesPaidFromFund(filter)` API call
    - Handle API errors và return user-friendly messages
    - _Requirements: 1.1, 4.3_

  - [ ] 10.2 Tạo fundIngredientService.js
    - Implement `restockIngredientFromFund(ingredientId, restockData)` API call
    - Implement `getRestockHistory(ingredientId)` API call
    - Handle API errors và return user-friendly messages
    - _Requirements: 2.1, 2.5_

  - [ ] 10.3 Cập nhật fundService.js
    - Cập nhật `getFundTransactions()` để include source info
    - Cập nhật `getFundTransactionDetail(id)` để include source details
    - _Requirements: 7.7, 7.8_

- [ ] 11. Frontend UI - Expense Form
  - [ ] 11.1 Cập nhật CreateExpenseModal.vue hoặc ExpenseForm.vue
    - Thêm checkbox "Chi từ quỹ" (Paid from Fund)
    - Khi checkbox được chọn, disable payment method selection và set to "fund"
    - Hiển thị current fund balance khi checkbox được chọn
    - Validate fund balance đủ trước khi submit
    - Call fundExpenseService.createExpenseFromFund khi checkbox được chọn
    - Hiển thị error message nếu insufficient balance
    - _Requirements: 1.1, 1.2, 6.3, 6.4, 10.1, 10.2_

  - [ ] 11.2 Cập nhật ExpenseList.vue
    - Thêm visual indicator (icon/badge) cho expenses paid from fund
    - Thêm filter option "Chi từ quỹ"
    - _Requirements: 4.1, 4.3_

  - [ ] 11.3 Cập nhật ExpenseDetail.vue
    - Hiển thị fund transaction link nếu expense paid from fund
    - Implement navigation to fund transaction detail
    - Disable edit amount nếu expense paid from fund
    - _Requirements: 4.2, 11.3_

- [ ] 12. Frontend UI - Ingredient Form
  - [ ] 12.1 Cập nhật IngredientRestockModal.vue hoặc IngredientForm.vue
    - Thêm checkbox "Chi từ quỹ" (Paid from Fund)
    - Hiển thị current fund balance khi checkbox được chọn
    - Calculate total cost (quantity × cost_per_unit) và validate với fund balance
    - Call fundIngredientService.restockIngredientFromFund khi checkbox được chọn
    - Hiển thị error message nếu insufficient balance
    - _Requirements: 2.1, 2.2, 6.5_

  - [ ] 12.2 Tạo IngredientRestockHistory.vue component
    - Hiển thị restock history cho ingredient
    - Show fund transaction info cho fund-paid restocks
    - Include link to expense và fund transaction
    - _Requirements: 2.5, 2.6, 2.7_

- [ ] 13. Frontend UI - Fund Transaction Views
  - [ ] 13.1 Cập nhật FundTransactionList.vue
    - Hiển thị source type và description cho mỗi transaction
    - Thêm filter by source_type
    - Show clickable link to source record
    - _Requirements: 7.7, 8.3, 8.4_

  - [ ] 13.2 Cập nhật FundTransactionDetail.vue
    - Hiển thị full source record details
    - Implement navigation button to source record (expense/ingredient)
    - Show all audit trail information
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 7.7, 7.8_

- [ ] 14. Checkpoint - Frontend UI Complete
  - Ensure all UI components work correctly, ask the user if questions arise.

- [ ] 15. Property-Based Tests
  - [ ]* 15.1 Write property test for Payment Method Consistency
    - **Property 14: Payment Method Consistency**
    - **Validates: Requirements 10.1, 10.3**

  - [ ]* 15.2 Write property test for Immutability Constraints
    - **Property 15: Immutability - Fund Transaction Link**
    - **Property 17: Immutability - Expense Amount**
    - **Validates: Requirements 11.1, 11.3**

  - [ ]* 15.3 Write property test for Uniqueness Constraint
    - **Property 16: Uniqueness - No Duplicate Fund Transactions**
    - **Validates: Requirements 11.2**

  - [ ]* 15.4 Write property test for Fund-Paid Expense Filter
    - **Property 8: Fund-Paid Expense Filter**
    - **Validates: Requirements 4.3**

  - [ ]* 15.5 Write property test for Audit Trail Completeness
    - **Property 9: Audit Trail Completeness**
    - **Validates: Requirements 7.1, 7.2, 7.3, 7.4, 7.5, 7.6**

  - [ ]* 15.6 Write property test for Source Type Validation
    - **Property 10: Source Type Validation**
    - **Validates: Requirements 8.3, 8.5**

- [ ] 16. End-to-End Testing và Integration
  - [ ]* 16.1 Write E2E test cho expense from fund flow
    - Test complete flow: create expense from fund → verify expense → verify fund transaction → verify balance update

  - [ ]* 16.2 Write E2E test cho ingredient restock from fund flow
    - Test complete flow: restock from fund → verify restock record → verify expense → verify fund transaction → verify stock update

  - [ ]* 16.3 Write E2E test cho insufficient balance scenarios
    - Test rejection khi fund balance không đủ
    - Verify error messages hiển thị đúng

  - [ ]* 16.4 Write E2E test cho immutability constraints
    - Test không thể edit expense amount khi paid from fund
    - Test không thể tạo duplicate fund transaction

- [ ] 17. Final Checkpoint
  - Ensure all tests pass, verify all requirements are met, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional testing tasks và có thể skip để faster MVP
- Mỗi task reference specific requirements để dễ traceability
- MongoDB transactions được sử dụng để đảm bảo atomicity (Requirements 9.1, 9.2, 9.3)
- Property-based tests sử dụng `gopter` library với minimum 100 iterations
- Frontend sử dụng Vue.js 3 với Composition API
- Backend sử dụng Go với clean architecture pattern
- Tất cả API endpoints yêu cầu Manager role authentication
