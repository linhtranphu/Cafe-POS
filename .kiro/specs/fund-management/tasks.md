# Fund Management - Implementation Tasks

## Phase 1: Backend Foundation

### 1.1 Create Domain Model
- [ ] 1.1.1 Create `backend/domain/fund/fund_transaction.go`
  - FundTransaction struct with all fields
  - FundBalance struct
  - TransactionType constants (deposit, withdrawal)
  - NewFundTransaction constructor
  - Validate() method
- [ ] 1.1.2 Create value objects in `backend/domain/fund/value_objects.go`
  - TransactionType enum
  - Validation helpers

### 1.2 Create Repository
- [ ] 1.2.1 Create `backend/infrastructure/mongodb/fund_transaction_repository.go`
  - Create() - insert new transaction
  - FindByID() - get transaction by ID
  - FindByDateRange() - get transactions with filters
  - FindAll() - get all transactions with pagination
  - Count() - count transactions with filters
- [ ] 1.2.2 Add MongoDB indexes
  - Index on timestamp (desc)
  - Index on type
  - Compound index on (timestamp, type)

### 1.3 Create Service Layer
- [ ] 1.3.1 Create `backend/application/services/fund_service.go`
  - CalculateCurrentBalance() - aggregate from all sources
  - CalculateTodayBalance() - opening, inflow, outflow
  - GetTransactionHistory() - with filters and pagination
  - CreateDeposit() - with transaction
  - CreateWithdrawal() - with balance validation and transaction
  - GetTransactionDetail() - by ID
- [ ] 1.3.2 Add balance calculation logic
  - Aggregate from fund_handovers (inflow)
  - Aggregate from cash_handovers (inflow)
  - Aggregate from fund_transactions (deposit/withdrawal)
  - Subtract starting_float from cashier_shifts (outflow)

### 1.4 Create HTTP Handlers
- [ ] 1.4.1 Create `backend/interfaces/http/fund_handler.go`
  - GetBalance() - GET /api/manager/fund/balance
  - GetTransactions() - GET /api/manager/fund/transactions
  - Deposit() - POST /api/manager/fund/deposit
  - Withdraw() - POST /api/manager/fund/withdraw
  - GetTransactionDetail() - GET /api/manager/fund/transactions/:id
- [ ] 1.4.2 Create request/response DTOs
  - BalanceResponse
  - TransactionHistoryResponse
  - DepositRequest/Response
  - WithdrawRequest/Response
  - TransactionDetailResponse
- [ ] 1.4.3 Add validation
  - Amount > 0
  - Reason min 10 characters
  - Withdrawal <= current balance

### 1.5 Register Routes
- [ ] 1.5.1 Update `backend/main.go`
  - Initialize FundTransactionRepository
  - Initialize FundService
  - Initialize FundHandler
  - Register routes under /api/manager/fund
  - Add manager role authorization middleware

## Phase 2: Transaction History Aggregation

### 2.1 Implement Transaction History
- [ ] 2.1.1 Create aggregation pipeline in FundService
  - Aggregate fund_handovers with type "fund_handover"
  - Aggregate cash_handovers with type "cash_handover"
  - Aggregate fund_transactions with type "deposit"/"withdrawal"
  - Aggregate cashier_shifts starting_float with type "starting_float"
  - Sort by timestamp desc
  - Apply filters (type, money_type, date_range)
- [ ] 2.1.2 Add pagination support
  - Limit and offset parameters
  - Return total count
- [ ] 2.1.3 Add metadata enrichment
  - Include shift_id, order_id where applicable
  - Include variance info for fund_handovers
  - Include discrepancy info for cash_handovers

## Phase 3: Deposit & Withdrawal with Transactions

### 3.1 Implement Deposit
- [ ] 3.1.1 Add CreateDeposit in FundService
  - Validate input (amount > 0, reason length)
  - Create FundTransaction record
  - Calculate balance before/after (optional)
  - Use MongoDB transaction for atomicity
  - Add audit logging
- [ ] 3.1.2 Test deposit flow
  - Unit tests for validation
  - Integration test for database

### 3.2 Implement Withdrawal
- [ ] 3.2.1 Add CreateWithdrawal in FundService
  - Validate input (amount > 0, reason length)
  - Calculate current balance
  - Validate withdrawal <= balance
  - Create FundTransaction record
  - Use MongoDB transaction for atomicity
  - Add audit logging
- [ ] 3.2.2 Test withdrawal flow
  - Unit tests for validation
  - Test insufficient balance error
  - Integration test for database

## Phase 4: Frontend Service & Store

### 4.1 Create Frontend Service
- [ ] 4.1.1 Create `frontend/src/services/fund.js`
  - getBalance() - fetch current balance
  - getTransactions(filters) - fetch transaction history
  - deposit(data) - create deposit
  - withdraw(data) - create withdrawal
  - getTransactionDetail(id) - fetch detail
- [ ] 4.1.2 Add error handling
  - Handle 400, 401, 403, 500 errors
  - Format error messages

### 4.2 Create Pinia Store (Optional)
- [ ] 4.2.1 Create `frontend/src/stores/fund.js`
  - State: balance, transactions, loading, error
  - Actions: fetchBalance, fetchTransactions, deposit, withdraw
  - Getters: formattedBalance, filteredTransactions

## Phase 5: Frontend View - Balance & List

### 5.1 Create Main View
- [ ] 5.1.1 Create `frontend/src/views/FundManagementView.vue`
  - Mobile-first layout with pull-to-refresh
  - Header with title "💰 Quản lý quỹ tiền"
  - Balance card section
  - Today summary section
  - Action buttons (Deposit/Withdraw)
  - Filters section
  - Transaction list section
  - Bottom navigation
  - Loading states
  - Error handling

### 5.2 Create Balance Card Component
- [ ] 5.2.1 Create `frontend/src/components/fund/FundBalanceCard.vue`
  - Display current cash balance (green)
  - Display current transfer balance (blue)
  - Display total balance
  - Gradient background (orange/yellow theme)
  - Auto-refresh on mount

### 5.3 Create Today Summary Component
- [ ] 5.3.1 Create `frontend/src/components/fund/FundTodaySummary.vue`
  - Display opening balance
  - Display total inflow (green, 📥)
  - Display total outflow (red, 📤)
  - Net change indicator

### 5.4 Create Transaction List Component
- [ ] 5.4.1 Create `frontend/src/components/fund/FundTransactionList.vue`
  - Display transactions in cards
  - Show transaction type icon
  - Show amount with color coding
  - Show timestamp and performer
  - Click to view detail
  - Empty state
  - Loading skeleton

### 5.5 Create Filters Component
- [ ] 5.5.1 Create `frontend/src/components/fund/FundFilters.vue`
  - Date filter dropdown (Today, Yesterday, This Week, This Month, Custom)
  - Transaction type filter (All, Deposit, Withdrawal, Handover)
  - Money type filter (All, Cash, Transfer)
  - Apply filters button

## Phase 6: Frontend Modals

### 6.1 Create Deposit Modal
- [ ] 6.1.1 Create `frontend/src/components/fund/DepositModal.vue`
  - Form with money type selector (Cash/Transfer/Both)
  - Cash amount input (if selected)
  - Transfer amount input (if selected)
  - Reason textarea (required, min 10 chars)
  - Validation
  - Submit button
  - Cancel button
  - Success/error feedback

### 6.2 Create Withdrawal Modal
- [ ] 6.2.1 Create `frontend/src/components/fund/WithdrawalModal.vue`
  - Form with money type selector (Cash/Transfer/Both)
  - Cash amount input (if selected)
  - Transfer amount input (if selected)
  - Show current balance warning
  - Reason textarea (required, min 10 chars)
  - Validation (amount <= balance)
  - Submit button
  - Cancel button
  - Success/error feedback

### 6.3 Create Transaction Detail Modal
- [ ] 6.3.1 Create `frontend/src/components/fund/TransactionDetailModal.vue`
  - Display all transaction fields
  - Show metadata (shift_id, variance, etc.)
  - Show balance before/after (if available)
  - Close button
  - Copy transaction ID button

## Phase 7: Integration & Navigation

### 7.1 Add Navigation
- [x] 7.1.1 Update `frontend/src/views/DashboardView.vue`
  - Add "💰 Quỹ tiền" button in manager section
  - Route to /manager/fund
- [x] 7.1.2 Update router
  - Add route for /manager/fund
  - Add manager role guard

### 7.2 Add Constants
- [x] 7.2.1 Create `frontend/src/constants/fund.js`
  - Transaction type constants
  - Money type constants
  - Filter options
  - Display text mappings
  - Helper functions (formatCurrency, formatDateTime, etc.)
  - All components updated to use constants

## Phase 8: Polish & Testing

### 8.1 UI Polish
- [ ] 8.1.1 Add pull-to-refresh functionality
- [ ] 8.1.2 Add loading skeletons
- [ ] 8.1.3 Add empty states with illustrations
- [ ] 8.1.4 Add success/error toasts
- [ ] 8.1.5 Optimize for mobile (touch targets, spacing)
- [ ] 8.1.6 Add animations (slide-up modals, fade transitions)

### 8.2 Error Handling
- [ ] 8.2.1 Handle network errors gracefully
- [ ] 8.2.2 Add retry mechanisms
- [ ] 8.2.3 Show user-friendly error messages
- [ ] 8.2.4 Add error logging

### 8.3 Testing
- [ ] 8.3.1 Manual testing
  - Test deposit flow (cash, transfer, both)
  - Test withdrawal flow (cash, transfer, both)
  - Test withdrawal with insufficient balance
  - Test filters (date, type, money type)
  - Test pagination
  - Test transaction detail view
  - Test on mobile devices
- [ ] 8.3.2 Create test script
  - Script to create test transactions
  - Script to verify balance calculations
- [ ] 8.3.3 Integration testing
  - Test with real fund_handovers
  - Test with real cash_handovers
  - Verify balance accuracy

### 8.4 Documentation
- [ ] 8.4.1 Create user guide
  - How to view balance
  - How to deposit money
  - How to withdraw money
  - How to view transaction history
- [ ] 8.4.2 Create API documentation
  - Document all endpoints
  - Add request/response examples
- [ ] 8.4.3 Create developer notes
  - Balance calculation logic
  - Transaction aggregation logic
  - Future enhancement ideas

## Phase 9: Performance Optimization

### 9.1 Backend Optimization
- [ ] 9.1.1 Add caching for current balance
  - Cache with 5-minute TTL
  - Invalidate on deposit/withdrawal
- [ ] 9.1.2 Optimize aggregation queries
  - Use MongoDB aggregation pipeline
  - Add proper indexes
- [ ] 9.1.3 Add query result caching
  - Cache transaction list with filters
  - Invalidate on new transactions

### 9.2 Frontend Optimization
- [ ] 9.2.1 Implement virtual scrolling for long lists
- [ ] 9.2.2 Add debouncing for filter changes
- [ ] 9.2.3 Lazy load transaction details
- [ ] 9.2.4 Optimize re-renders with memo/computed

## Phase 10: Security & Audit

### 10.1 Security Hardening
- [ ] 10.1.1 Add rate limiting for deposit/withdrawal endpoints
- [ ] 10.1.2 Add CSRF protection
- [ ] 10.1.3 Validate all inputs server-side
- [ ] 10.1.4 Add request logging

### 10.2 Audit Trail
- [ ] 10.2.1 Log all deposit/withdrawal operations
- [ ] 10.2.2 Include user info, IP, timestamp
- [ ] 10.2.3 Add audit log viewer (future)

## Notes

### Dependencies
- Phase 4-6 depend on Phase 1-3 (backend must be complete)
- Phase 7 depends on Phase 4-6 (frontend components must exist)
- Phase 8-10 can be done in parallel after Phase 7

### Estimated Time
- Phase 1-3: 2-3 days (backend)
- Phase 4-6: 2-3 days (frontend)
- Phase 7: 0.5 day (integration)
- Phase 8-10: 1-2 days (polish & testing)
- Total: 5-8 days

### Priority
- P0 (Must Have): Phase 1-7
- P1 (Should Have): Phase 8
- P2 (Nice to Have): Phase 9-10
