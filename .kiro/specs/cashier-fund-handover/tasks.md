# Implementation Tasks - Cashier Fund Handover

## Overview
Implement tính năng cho phép Cashier xem số tiền đang quản lý và handover về quỹ khi đóng ca.

## Task Breakdown

### Phase 1: Backend Foundation

#### Task 1.1: Create FundHandover Domain Model
**File**: `backend/domain/cashier/fund_handover.go`

- [ ] Create FundHandover struct with all fields
- [ ] Implement NewFundHandover constructor
- [ ] Implement HasVariance() method
- [ ] Implement DocumentVariance() method
- [ ] Implement SetReceiver() method
- [ ] Add validation logic
- [ ] Write unit tests

**Estimated Time**: 2 hours

---

#### Task 1.2: Create FundHandoverRepository
**File**: `backend/infrastructure/mongodb/fund_handover_repository.go`

- [ ] Create FundHandoverRepository struct
- [ ] Implement Create() method
- [ ] Implement FindByID() method
- [ ] Implement FindByCashierShift() method
- [ ] Implement FindByCashier() method with pagination
- [ ] Implement FindByDateRange() method
- [ ] Create indexes
- [ ] Write unit tests

**Estimated Time**: 3 hours

---

#### Task 1.3: Extend CashierShiftService
**File**: `backend/application/services/cashier_shift_service.go`

- [ ] Add GetManagedFunds() method
- [ ] Add CloseShiftWithFundHandover() method
- [ ] Implement transaction logic
- [ ] Add variance validation
- [ ] Add audit logging
- [ ] Write unit tests
- [ ] Write integration tests

**Estimated Time**: 4 hours

---

#### Task 1.4: Create API Endpoints
**Files**: 
- `backend/api/handlers/cashier_handler.go`
- `backend/api/routes/cashier_routes.go`

- [ ] GET /api/cashier/shifts/:id/managed-funds
- [ ] POST /api/cashier/shifts/:id/close (extend existing)
- [ ] GET /api/cashier/fund-handovers
- [ ] Add request/response DTOs
- [ ] Add validation middleware
- [ ] Add error handling
- [ ] Write API tests

**Estimated Time**: 3 hours

---

### Phase 2: Frontend - Dashboard Display

#### Task 2.1: Add Managed Funds Section to CashierDashboard
**File**: `frontend/src/views/CashierDashboard.vue`

- [ ] Create managed funds computed property
- [ ] Add API call to fetch managed funds
- [ ] Design managed funds card component
- [ ] Add received cash display
- [ ] Add received transfer display
- [ ] Add total managed funds display
- [ ] Add responsibility warning message
- [ ] Style with gradient and colors
- [ ] Add loading state
- [ ] Add error handling
- [ ] Test responsive design

**Estimated Time**: 3 hours

---

### Phase 3: Frontend - Closure Flow

#### Task 3.1: Extend CashierShiftClosureV2 - Step 1
**File**: `frontend/src/views/CashierShiftClosureV2.vue`

- [ ] Add managed funds summary display
- [ ] Show starting float
- [ ] Show received cash from handovers
- [ ] Show received transfer from handovers
- [ ] Calculate expected cash
- [ ] Style summary card
- [ ] Add navigation to next step

**Estimated Time**: 2 hours

---

#### Task 3.2: Extend CashierShiftClosureV2 - Step 2
**File**: `frontend/src/views/CashierShiftClosureV2.vue`

- [ ] Add cash counting interface
- [ ] Display expected cash amount
- [ ] Add input for actual cash
- [ ] Calculate variance in real-time
- [ ] Display variance with color coding
- [ ] Add validation
- [ ] Add navigation buttons

**Estimated Time**: 2 hours

---

#### Task 3.3: Extend CashierShiftClosureV2 - Step 3
**File**: `frontend/src/views/CashierShiftClosureV2.vue`

- [ ] Add variance documentation form
- [ ] Add variance reason dropdown
- [ ] Add variance notes textarea
- [ ] Add validation (min 10 chars)
- [ ] Show/hide based on variance
- [ ] Style form
- [ ] Add navigation buttons

**Estimated Time**: 2 hours

---

#### Task 3.4: Extend CashierShiftClosureV2 - Step 4
**File**: `frontend/src/views/CashierShiftClosureV2.vue`

- [ ] Add final confirmation screen
- [ ] Display handover summary
- [ ] Show cash amount
- [ ] Show transfer amount
- [ ] Show total amount
- [ ] Show variance (if any)
- [ ] Add confirmation button
- [ ] Integrate with API
- [ ] Handle success/error
- [ ] Redirect after success

**Estimated Time**: 3 hours

---

### Phase 4: Testing & Refinement

#### Task 4.1: Backend Testing
- [ ] Write unit tests for FundHandover domain
- [ ] Write unit tests for repository
- [ ] Write unit tests for service
- [ ] Write integration tests for API
- [ ] Test transaction rollback scenarios
- [ ] Test variance validation
- [ ] Test error handling
- [ ] Achieve >80% code coverage

**Estimated Time**: 4 hours

---

#### Task 4.2: Frontend Testing
- [ ] Test dashboard display
- [ ] Test closure flow (happy path)
- [ ] Test closure with variance
- [ ] Test validation errors
- [ ] Test API error handling
- [ ] Test responsive design
- [ ] Test on mobile devices
- [ ] Cross-browser testing

**Estimated Time**: 3 hours

---

#### Task 4.3: E2E Testing
- [ ] Create test script for complete flow
- [ ] Test: Start shift → Receive handovers → Close shift
- [ ] Test: Closure with no variance
- [ ] Test: Closure with variance
- [ ] Test: Closure with large variance
- [ ] Test: Transaction rollback on error
- [ ] Document test scenarios
- [ ] Create test data fixtures

**Estimated Time**: 3 hours

---

### Phase 5: Documentation & Deployment

#### Task 5.1: Documentation
- [ ] Update API documentation
- [ ] Create user guide for cashiers
- [ ] Document closure process
- [ ] Create training materials
- [ ] Update system architecture docs
- [ ] Document database schema
- [ ] Create troubleshooting guide

**Estimated Time**: 3 hours

---

#### Task 5.2: Deployment
- [ ] Create database migration script
- [ ] Create indexes
- [ ] Deploy backend changes
- [ ] Deploy frontend changes
- [ ] Verify in staging environment
- [ ] Monitor for errors
- [ ] Deploy to production
- [ ] Post-deployment verification

**Estimated Time**: 2 hours

---

## Total Estimated Time: 42 hours (~5-6 days)

## Dependencies
- MongoDB connection
- Existing CashierShift domain model
- Existing CashierShiftService
- Existing CashierDashboard component
- Existing CashierShiftClosureV2 component

## Success Criteria
- [ ] Cashiers can see managed funds in dashboard
- [ ] Cashiers can close shift with fund handover
- [ ] Variance is calculated correctly
- [ ] Variance documentation is enforced
- [ ] Fund handover records are created
- [ ] Transaction atomicity is maintained
- [ ] All tests pass
- [ ] Documentation is complete
- [ ] Feature is deployed to production

## Risks & Mitigation
1. **Risk**: Transaction timeout during closure
   - **Mitigation**: Set appropriate timeout, optimize queries

2. **Risk**: Data inconsistency if transaction fails
   - **Mitigation**: Thorough testing of rollback scenarios

3. **Risk**: UI confusion during multi-step closure
   - **Mitigation**: Clear visual indicators, progress bar

4. **Risk**: Performance impact on dashboard
   - **Mitigation**: Cache managed funds data, optimize queries

## Notes
- Keep receiver_id nullable for future extension
- Design API to be backward compatible
- Ensure mobile-friendly UI
- Add comprehensive error messages
- Log all fund handover operations for audit
