# Implementation Plan: Bank Transfer Handover

## Overview

Tính năng này mở rộng hệ thống bàn giao tiền hiện tại để hỗ trợ cả tiền chuyển khoản bên cạnh tiền mặt. Implementation sẽ được thực hiện theo các bước:

1. Mở rộng domain models với các trường transfer
2. Cập nhật service layer để xử lý dual amounts
3. Cập nhật API endpoints
4. Cập nhật frontend components
5. Viết tests (property-based và unit tests)

## Tasks

- [x] 1. Extend Domain Models with Transfer Fields
  - Modify `backend/domain/handover/cash_handover.go` to add transfer amount fields
  - Modify `backend/domain/order/shift.go` to add transfer tracking fields
  - Modify `backend/domain/cashier/cashier_shift.go` to add received_transfer field
  - Add helper methods for calculating totals and checking transfer presence
  - _Requirements: 2.1, 2.5, 6.4_

- [ ]* 1.1 Write property test for separate field storage
  - **Property 3: Handover Record Structure**
  - **Validates: Requirements 2.5, 1.4**

- [x] 2. Update Repository Layer for Transfer Fields
  - Modify `backend/infrastructure/mongodb/cash_handover_repository.go` to handle new fields
  - Modify `backend/infrastructure/mongodb/shift_repository.go` to handle transfer fields
  - Add indexes for transfer_declared_amount field
  - Ensure backward compatibility when reading old records
  - _Requirements: 7.3, 7.5_

- [ ]* 2.1 Write property test for backward compatibility
  - **Property 16: Migration Data Integrity**
  - **Validates: Requirements 7.3, 7.5**

- [x] 3. Extend Service Layer with Dual Amount Logic
  - [x] 3.1 Add CreateHandoverWithTransfer method to CashHandoverService
    - Validate cash and transfer amounts against shift balances
    - Create handover record with separate amounts
    - Update shift remaining_cash and remaining_transfer
    - _Requirements: 1.1, 1.3, 1.4_

  - [ ]* 3.2 Write property test for declared amount validation
    - **Property 8: Declared Amount Validation**
    - **Validates: Requirements 1.3, 6.1**

  - [x] 3.3 Add ConfirmHandoverWithDualAmounts method to CashHandoverService
    - Accept separate actual_cash and actual_transfer amounts
    - Calculate separate discrepancies for cash and transfer
    - Check manager approval thresholds
    - Update balances atomically
    - _Requirements: 3.2, 3.3, 4.1, 4.4_

  - [ ]* 3.4 Write property test for discrepancy calculation
    - **Property 4: Discrepancy Calculation**
    - **Validates: Requirements 3.3, 4.3**

  - [ ]* 3.5 Write property test for atomic confirmation
    - **Property 5: Atomic Confirmation**
    - **Validates: Requirements 4.1, 4.4, 4.5**

  - [x] 3.6 Update UpdateDualBalances method
    - Update waiter shift: handed_over_cash and handed_over_transfer
    - Update cashier shift: received_cash and received_transfer
    - Ensure atomic transaction for both updates
    - _Requirements: 4.4, 4.5_

  - [x] 3.7 Add CalculateTransferRevenue method to ShiftService
    - Sum all TRANSFER and QR payment orders
    - Update shift.transfer_revenue field
    - Calculate remaining_transfer
    - _Requirements: 2.2_

  - [ ]* 3.8 Write property test for transfer revenue calculation
    - **Property 2: Transfer Revenue Calculation**
    - **Validates: Requirements 2.2, 1.2**

- [-] 4. Update API Endpoints
  - [x] 4.1 Modify POST /api/waiter/shifts/:id/handover endpoint
    - Accept cash_amount and transfer_amount in request body
    - Call CreateHandoverWithTransfer service method
    - Return handover with separate amounts
    - _Requirements: 1.1, 1.4_

  - [ ] 4.2 Modify POST /api/waiter/shifts/:id/handover-and-end endpoint
    - Auto-calculate both remaining_cash and remaining_transfer
    - Create end-shift handover with both amounts
    - _Requirements: 1.5_

  - [ ]* 4.3 Write property test for end-shift handover amounts
    - **Property 11: End-Shift Handover Amounts**
    - **Validates: Requirements 1.5**

  - [x] 4.4 Modify PATCH /api/cashier/handovers/:id/confirm endpoint
    - Accept actual_cash_amount and actual_transfer_amount
    - Call ConfirmHandoverWithDualAmounts service method
    - Return updated handover with discrepancies
    - _Requirements: 3.2, 3.3_

  - [ ] 4.5 Update GET /api/cashier/shifts/:id/status endpoint
    - Include transfer_revenue in response
    - Include remaining_transfer in response
    - _Requirements: 2.1_

- [ ]* 4.6 Write unit tests for API endpoints
  - Test handover creation with cash only, transfer only, and both
  - Test validation errors for invalid amounts
  - Test end-shift handover auto-calculation
  - Test confirmation with dual amounts
  - _Requirements: 1.1, 1.3, 1.5, 3.2_

- [ ] 5. Checkpoint - Ensure backend tests pass
  - Run all property-based tests
  - Run all unit tests
  - Verify API endpoints work correctly
  - Ask the user if questions arise

- [-] 6. Update Frontend - ShiftView (Waiter)
  - [x] 6.1 Add handover type selection UI
    - Add buttons for "Cash only", "Transfer only", "Both"
    - Show/hide amount inputs based on selection
    - _Requirements: 1.1_

  - [x] 6.2 Display transfer revenue in shift summary
    - Add transfer_revenue display with blue badge
    - Show remaining_transfer amount
    - _Requirements: 1.2, 5.1_

  - [x] 6.3 Add transfer amount input field
    - Add input for transfer_amount
    - Validate against remaining_transfer
    - Show validation errors
    - _Requirements: 1.3_

  - [x] 6.4 Update handover submission logic
    - Send cash_amount and transfer_amount to API
    - Handle validation errors
    - Update shift data after successful handover
    - _Requirements: 1.4_

  - [x] 6.5 Update handover history display
    - Show both cash and transfer amounts for each handover
    - Use color coding (green for cash, blue for transfer)
    - _Requirements: 5.2_

- [-] 7. Update Frontend - CashierShiftView (Cashier)
  - [ ] 7.1 Update pending handover display
    - Show declared cash amount with green badge
    - Show declared transfer amount with blue badge
    - Distinguish between cash-only, transfer-only, and combined handovers
    - _Requirements: 3.1, 5.3_

  - [ ] 7.2 Add dual amount confirmation inputs
    - Add input for actual_cash_amount
    - Add input for actual_transfer_amount
    - Show hints for reconciliation process
    - _Requirements: 3.2_

  - [ ] 7.3 Add discrepancy display
    - Calculate and show cash_discrepancy
    - Calculate and show transfer_discrepancy
    - Show total_discrepancy
    - Use color coding for shortage/overage
    - _Requirements: 3.3, 5.5_

  - [ ] 7.4 Update confirmation submission logic
    - Send actual_cash_amount and actual_transfer_amount to API
    - Handle manager approval requirement
    - Show success/error messages
    - _Requirements: 3.2, 3.5_

  - [ ] 7.5 Update shift summary display
    - Show cash revenue and transfer revenue separately
    - Use consistent color coding
    - _Requirements: 5.1, 5.4_

- [ ]* 7.6 Write unit tests for frontend components
  - Test handover type selection
  - Test amount validation
  - Test discrepancy calculation display
  - Test color coding consistency
  - _Requirements: 1.1, 1.3, 3.3, 5.4_

- [ ] 8. Update Manager Approval Logic
  - [ ] 8.1 Update approval threshold checks
    - Check total_discrepancy > 100,000 VND
    - Check transfer_discrepancy > 50,000 VND
    - Set requires_approval flag accordingly
    - _Requirements: 8.1, 8.2_

  - [ ]* 8.2 Write property test for manager approval threshold
    - **Property 12: Manager Approval Threshold**
    - **Validates: Requirements 8.1, 8.2**

  - [ ] 8.3 Update manager approval endpoint
    - Update shift balances after approval
    - Handle both cash and transfer amounts
    - _Requirements: 8.4_

  - [ ]* 8.4 Write property test for balance update after approval
    - **Property 13: Balance Update After Approval**
    - **Validates: Requirements 8.4**

  - [ ] 8.5 Update manager rejection endpoint
    - Return handover to pending status
    - Do not update balances
    - _Requirements: 8.5_

  - [ ]* 8.6 Write property test for rejection state transition
    - **Property 14: Rejection State Transition**
    - **Validates: Requirements 8.5**

- [ ] 9. Add Audit Logging
  - [ ] 9.1 Create audit log structure
    - Define AuditLog model with required fields
    - Create audit_logs collection
    - _Requirements: 6.5_

  - [ ] 9.2 Add logging to handover operations
    - Log handover creation with before/after state
    - Log confirmation with amounts and discrepancies
    - Log approval/rejection with manager info
    - _Requirements: 6.5_

  - [ ]* 9.3 Write property test for audit logging
    - **Property 15: Audit Logging**
    - **Validates: Requirements 6.5**

- [ ] 10. Add Validation and Error Handling
  - [ ] 10.1 Add non-negative transfer amount validation
    - Validate in request structs
    - Return appropriate error messages
    - _Requirements: 6.3_

  - [ ]* 10.2 Write property test for non-negative amounts
    - **Property 9: Non-Negative Transfer Amounts**
    - **Validates: Requirements 6.3**

  - [ ] 10.3 Add single pending handover validation
    - Check for existing pending handover before creation
    - Return conflict error if exists
    - _Requirements: 6.2_

  - [ ]* 10.4 Write property test for single pending handover
    - **Property 10: Single Pending Handover**
    - **Validates: Requirements 6.2**

  - [ ] 10.5 Add comprehensive error handling
    - Handle validation errors with clear messages
    - Handle authorization errors
    - Handle business logic errors
    - Implement transaction rollback on failures
    - _Requirements: 4.5_

- [ ]* 10.6 Write unit tests for error scenarios
  - Test invalid amount errors
  - Test pending handover conflict
  - Test authorization errors
  - Test transaction rollback
  - _Requirements: 1.3, 4.5, 6.2, 6.3_

- [ ] 11. Ensure Backward Compatibility
  - [ ] 11.1 Add default value handling for old records
    - Set transfer amounts to zero when missing
    - Handle null transfer fields gracefully
    - _Requirements: 7.2, 7.3_

  - [ ]* 11.2 Write property test for backward compatibility
    - **Property 7: Backward Compatibility**
    - **Validates: Requirements 7.1, 7.2, 7.4**

  - [ ] 11.3 Test cash-only handover flow
    - Verify cash-only handovers work without transfer fields
    - Verify confirmation works with only cash amount
    - _Requirements: 7.1, 7.4_

- [ ] 12. Integration Testing
  - [ ] 12.1 Test complete handover flow
    - Create shift with mixed payment orders
    - Create handover with both amounts
    - Confirm with actual amounts
    - Verify all balances updated correctly
    - _Requirements: 1.1, 1.4, 3.2, 4.4_

  - [ ] 12.2 Test discrepancy approval flow
    - Create handover with large discrepancy
    - Verify manager approval required
    - Manager approves
    - Verify balances updated after approval
    - _Requirements: 3.5, 8.1, 8.4_

  - [ ] 12.3 Test end-shift flow
    - Create shift with mixed payments
    - Create end-shift handover
    - Verify auto-calculation of both amounts
    - Confirm and verify shift closes correctly
    - _Requirements: 1.5_

  - [ ] 12.4 Test backward compatibility
    - Create cash-only handover
    - Verify it works without transfer fields
    - Display old handover records
    - Verify no errors
    - _Requirements: 7.1, 7.3, 7.4_

- [ ] 13. Final Checkpoint - Complete Testing
  - Run all property-based tests (minimum 100 iterations each)
  - Run all unit tests
  - Run all integration tests
  - Verify all acceptance criteria met
  - Ask the user if questions arise

- [ ] 14. Documentation and Deployment
  - Update API documentation with new fields
  - Update user guide with transfer handover instructions
  - Create deployment checklist
  - Prepare rollback plan
  - _Requirements: All_

## Notes

- Tasks marked with `*` are optional test tasks and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Property tests validate universal correctness properties with minimum 100 iterations
- Unit tests validate specific examples and edge cases
- Integration tests validate end-to-end flows
- Checkpoints ensure incremental validation
- Backward compatibility is critical - old handovers must continue to work
- New fields will have default values (0) for existing records, handled by repository layer

