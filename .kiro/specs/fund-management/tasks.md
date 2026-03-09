# Quản lý Quỹ Tiền — Implementation Tasks

## Phase 1: Backend Domain & Repository

### Task 1.1 — Cập nhật domain `fund_transaction.go`
**File:** `backend/domain/fund/fund_transaction.go`

- [ ] Thêm type `FundType` + 4 constants (operating, inventory, profit, cash_drawer)
- [ ] Thêm type `SourceType` + 5 constants (expense, ingredient, handover, manual, fund_transfer)
- [ ] Thêm TransactionType constants: `fund_handover`, `fund_transfer`, `starting_float`
- [ ] Thêm fields vào struct `FundTransaction`: `FundType`, `SourceType`, `SourceID`, `Description`, `LinkedTxID`
- [ ] Cập nhật `Validate()` cho phép các type mới
- [ ] Thêm helper methods: `SetSource()`, `SetFundType()`, `LinkTo()`

---

### Task 1.2 — Cập nhật repository `fund_transaction_repository.go`
**File:** `backend/infrastructure/mongodb/fund_transaction_repository.go`

- [ ] Thêm method `FindByFilter(ctx, filter bson.M, limit, offset int)`
- [ ] Thêm method `CountByFilter(ctx, filter bson.M)`
- [ ] Thêm indexes: `fund_type`, compound `fund_type+timestamp`, sparse `source_type+source_id`

---

## Phase 2: Backend Service

### Task 2.1 — Thêm balance methods vào `fund_service.go`
**File:** `backend/application/services/fund_service.go`

- [ ] Thêm structs: `AllFundBalances`, `FundTransferRequest`, `FundTransferResult`, `DailyReport`
- [ ] Implement `CalculateBalanceByFundType(ctx, fundType)` với backward compat cho "operating"
- [ ] Implement `GetAllFundBalances(ctx)` gọi 4 lần CalculateBalanceByFundType + tổng
- [ ] Implement `TransferBetweenFunds(ctx, req)` atomic (MongoDB session + pre-generated IDs)
- [ ] Implement `GetDailyReport(ctx, date)` (opening = trước ngày, inflow/outflow = trong ngày)
- [ ] Cập nhật `CreateDeposit()` thêm variadic `fundType ...fund.FundType` (default "operating")
- [ ] Cập nhật `CreateWithdrawal()` thêm variadic `fundType ...fund.FundType` (default "operating")

---

## Phase 3: Backend HTTP Handler & Routes

### Task 3.1 — Cập nhật `fund_handler.go`
**File:** `backend/interfaces/http/fund_handler.go`

- [ ] Thêm `FundType string` vào `DepositRequest` và `WithdrawRequest`
- [ ] Cập nhật `Deposit()` và `Withdraw()` để truyền fund_type xuống service
- [ ] Thêm handler `GetAllFundBalances()` → `GET /balances`
- [ ] Thêm `FundTransferRequest` struct
- [ ] Thêm handler `Transfer()` → `POST /transfer` (validate, call TransferBetweenFunds)
- [ ] Thêm handler `GetDailyReport()` → `GET /daily-report?date=YYYY-MM-DD`
- [ ] Cập nhật `GetTransactions()` để nhận thêm query params: `fund_type`, `source_type`

### Task 3.2 — Thêm routes trong `main.go`
**File:** `backend/main.go`

- [ ] `GET /api/manager/fund/balances` → `fundHandler.GetAllFundBalances`
- [ ] `POST /api/manager/fund/transfer` → `fundHandler.Transfer`
- [ ] `GET /api/manager/fund/daily-report` → `fundHandler.GetDailyReport`

---

## Phase 4: Frontend Constants & Services

### Task 4.1 — Cập nhật `constants/fund.js`
**File:** `frontend/src/constants/fund.js`

- [ ] Thêm `FUND_TYPES`, `FUND_TYPE_LABELS`, `FUND_TYPE_ICONS`, `FUND_TYPE_COLORS`
- [ ] Thêm helper `getFundTypeLabel(type)`, `getFundTypeIcon(type)`
- [ ] Thêm `FUND_TYPE_FILTER_OPTIONS` (all + 4 quỹ)
- [ ] Cập nhật `TRANSACTION_TYPE_FILTER_OPTIONS` thêm `fund_transfer`
- [ ] Cập nhật `SOURCE_TYPE_LABELS` thêm `handover`, `fund_transfer`
- [ ] Thêm `SOURCE_TYPE_FILTER_OPTIONS` (all + 5 nguồn)

### Task 4.2 — Cập nhật `services/fund.js`
**File:** `frontend/src/services/fund.js`

- [ ] Thêm `getAllBalances()` → `GET /manager/fund/balances`
- [ ] Thêm `transferBetweenFunds(data)` → `POST /manager/fund/transfer`
- [ ] Thêm `getDailyReport(date)` → `GET /manager/fund/daily-report`
- [ ] Cập nhật `getTransactions(filters)` thêm `fund_type` và `source_type` params

---

## Phase 5: Frontend View

### Task 5.1 — Cập nhật `FundManagementView.vue` (script)
**File:** `frontend/src/views/FundManagementView.vue`

- [ ] Import thêm: `FUND_TYPES`, `FUND_TYPE_LABELS`, `FUND_TYPE_ICONS`, `FUND_TYPE_COLORS`, `FUND_TYPE_FILTER_OPTIONS`, `SOURCE_TYPE_FILTER_OPTIONS`, `getAllBalances`, `transferBetweenFunds`
- [ ] Thêm refs: `allBalances`, `loadingAllBalances`, `showTransferModal`, `transferring`, `transferForm`
- [ ] Thêm computed `fundTypeKeys`
- [ ] Thêm `filters.fund_type = 'all'` và `filters.source_type = 'all'`
- [ ] Implement `loadAllBalances()` → gọi `getAllBalances()`
- [ ] Implement `submitTransfer()` → validate + `transferBetweenFunds()` + refresh
- [ ] Cập nhật `refreshData()` thêm `loadAllBalances()`

### Task 5.2 — Cập nhật `FundManagementView.vue` (template)
**File:** `frontend/src/views/FundManagementView.vue`

- [ ] Thay Balance Card đơn → Total Summary Card + lưới 2×2 của 4 thẻ quỹ (loading skeleton khi chưa có data)
- [ ] Thay 2 nút action → 3 nút: Nạp / Rút / Chuyển quỹ
- [ ] Thêm Transfer Modal (from/to fund selector, money type, amount, reason, source balance preview)
- [ ] Thêm CSS `slide-up` transition cho modal
- [ ] Thêm filter `fund_type` dropdown trong phần Filters
- [ ] Thêm filter `source_type` dropdown trong phần Filters
- [ ] Cập nhật transaction list: thêm badge `fund_type`, badge `source_type`
- [ ] Cập nhật transaction list: thêm "Xem record gốc" link nếu có source_type/source_id

### Task 5.3 — Cập nhật modal `WithdrawModal.vue`
**File:** `frontend/src/components/fund/WithdrawModal.vue`

- [ ] Thêm dropdown chọn quỹ (fund_type)
- [ ] Hiển thị số dư của quỹ được chọn (gọi API balance theo fund_type)
- [ ] Validate amount ≤ số dư quỹ được chọn

### Task 5.4 — Cập nhật modal `DepositModal.vue`
**File:** `frontend/src/components/fund/DepositModal.vue`

- [ ] Thêm dropdown chọn quỹ (fund_type), default "operating"

---

## Phase 6: Integration & Kết nối bàn giao ca

### Task 6.1 — Tự động tạo FundTransaction khi cashier handover
**File:** `backend/application/services/cash_reconciliation_service.go` hoặc `cashier_shift_service.go`

- [ ] Khi cashier kết ca thành công và bàn giao → gọi `fundService.RecordHandover()`
- [ ] `RecordHandover()` tạo FundTransaction(type=fund_handover, fund_type=cash_drawer, source_type=handover, source_id=shiftID)

---

## Thứ tự ưu tiên implement

```
1.1 → 1.2 → 2.1 → 3.1 → 3.2    (backend, test API trước)
       ↓
4.1 → 4.2 → 5.1 → 5.2           (frontend core)
       ↓
5.3 → 5.4 → 6.1                  (polish + integration)
```

## File checklist

| File | Thay đổi |
|------|----------|
| `backend/domain/fund/fund_transaction.go` | Thêm FundType, SourceType, fields, helpers |
| `backend/infrastructure/mongodb/fund_transaction_repository.go` | FindByFilter, CountByFilter, indexes |
| `backend/application/services/fund_service.go` | GetAllFundBalances, TransferBetweenFunds, GetDailyReport, update Deposit/Withdraw |
| `backend/interfaces/http/fund_handler.go` | GetAllFundBalances, Transfer, GetDailyReport handlers |
| `backend/main.go` | 3 routes mới |
| `frontend/src/constants/fund.js` | FUND_TYPES, FUND_TYPE_*, SOURCE_TYPE_*, filter options |
| `frontend/src/services/fund.js` | getAllBalances, transferBetweenFunds, getDailyReport |
| `frontend/src/views/FundManagementView.vue` | 4-fund UI, transfer modal, filters |
| `frontend/src/components/fund/DepositModal.vue` | fund_type selector |
| `frontend/src/components/fund/WithdrawModal.vue` | fund_type selector + balance per fund |
