# Fund System — Implementation Tasks

**Last Updated:** 2026-03-11  
**Spec:** FUND_SYSTEM_REQUIREMENTS.md  
**Tech Spec:** FUND_SYSTEM_SPEC.md

---

## Status Legend
- ✅ Done
- 🔄 In Progress
- ⬜ Pending
- ❌ Blocked

---

## Phase 0: Core Infrastructure ✅ COMPLETE

| # | Task | Status | Notes |
|---|------|--------|-------|
| 0.1 | Domain model `JournalEntry` + `LedgerLine` | ✅ | `backend/domain/fund/journal_entry.go` |
| 0.2 | Fund types: 5 real + 5 external counterparts | ✅ | operating, inventory, profit, cash_drawer, waiter_float, owner, supplier, customer, cash_shortage, cash_overage |
| 0.3 | `JournalRepository` (MongoDB) | ✅ | `backend/infrastructure/mongodb/journal_repository.go` |
| 0.4 | `JournalService` với tất cả business methods | ✅ | `backend/application/services/journal_service.go` |
| 0.5 | Wire vào `main.go` | ✅ | |
| 0.6 | `FundHandler` dùng journal (xóa old fund_handler routes) | ✅ | `backend/interfaces/http/fund_handler.go` |
| 0.7 | Xóa `fund_service.go` và `fund_transaction_repository.go` | ✅ | |
| 0.8 | `frontend/src/services/journal.js` | ✅ | getJournalBalances, getJournalEntries, getJournalEntry |
| 0.9 | `frontend/src/constants/fund.js` — EVENT_TYPES, FUND_TYPES, labels, icons, colors | ✅ | Bao gồm external accounts |

---

## Phase 1: Cashier Shift ✅ COMPLETE

| # | Task | Status | Notes |
|---|------|--------|-------|
| 1.1 | `RecordCashierShiftStart` → journal entry | ✅ | DEBIT cash_drawer / CREDIT operating |
| 1.2 | `RecordCashierShiftEnd` → journal entry | ✅ | DEBIT operating / CREDIT cash_drawer |
| 1.3 | `CloseShiftWithFundHandover` — dùng journal | ✅ | `cashier_shift_service.go` |
| 1.4 | Fix variance formula: expectedCash = startingFloat + receivedCash - **distributedCash** | ✅ | Bug fix 2026-03-11 |
| 1.5 | `CashierShiftClosureV2.vue` — dùng `closeShiftWithFundHandover` | ✅ | |

---

## Phase 2: Manager Fund Operations ✅ COMPLETE

| # | Task | Status | Notes |
|---|------|--------|-------|
| 2.1 | `RecordManagerDeposit` → DEBIT fund / CREDIT owner | ✅ | |
| 2.2 | `RecordManagerWithdrawal` → DEBIT owner / CREDIT fund | ✅ | |
| 2.3 | `RecordFundTransfer` → DEBIT to / CREDIT from | ✅ | |
| 2.4 | `FundManagementView.vue` — lịch sử từ journal | ✅ | realLines() lọc external accounts |
| 2.5 | Fix `WithdrawModal` — `max` dùng `fundBalance` thay `currentBalance` prop | ✅ | Bug fix 2026-03-11 |

---

## Phase 3: Expense / Restock / Facility ✅ COMPLETE

| # | Task | Status | Notes |
|---|------|--------|-------|
| 3.1 | `CreateExpenseFromFund` → DEBIT supplier / CREDIT operating | ✅ | `fund_expense_integration_service.go` |
| 3.2 | `RestockIngredientFromFund` → DEBIT supplier / CREDIT inventory | ✅ | |
| 3.3 | `PurchaseFacilityFromFund` → DEBIT supplier / CREDIT operating | ✅ | |
| 3.4 | `expense_handler`, `ingredient_handler`, `facility_handler` trả về `journal_entry` | ✅ | |
| 3.5 | `OperatingExpenseView.vue` — list từ journal entries (event_type=expense) | ✅ | |
| 3.6 | `IngredientManagementView.vue` — balance từ journal | ✅ | |
| 3.7 | `FacilityManagementView.vue` — balance từ journal | ✅ | |
| 3.8 | `getAllBalances()` / `getBalance()` trong `fund.js` → dùng `/journal-balances` | ✅ | |

---

## Phase 4: Waiter Shift & Handover ⬜ PENDING

| # | Task | Status | Priority | Notes |
|---|------|--------|----------|-------|
| 4.1 | `RecordWaiterShiftStart` trong `JournalService` | ⬜ | HIGH | DEBIT waiter_float / CREDIT cash_drawer |
| 4.2 | Wire `RecordWaiterShiftStart` vào `ShiftService.StartShift` | ⬜ | HIGH | Thay thế old fund_transaction call |
| 4.3 | `RecordWaiterHandover` trong `JournalService` — 3 cases (exact/shortage/overage) | ⬜ | HIGH | Dùng cash_shortage / cash_overage accounts |
| 4.4 | Wire `RecordWaiterHandover` vào `CashHandoverService.ConfirmHandover` | ⬜ | HIGH | Atomic với handover confirmation |
| 4.5 | Frontend `CashierHandoverView.vue` — hiển thị shortage/overage rõ ràng | ⬜ | MEDIUM | |
| 4.6 | Xóa old `RecordWaiterStartFloat` và handover journal calls cũ | ⬜ | LOW | Cleanup sau khi 4.1–4.4 xong |

**Accounting logic cho 4.3:**
```
shortage = expectedCash - actualCash   (khi actualCash < expectedCash)
overage  = actualCash - expectedCash   (khi actualCash > expectedCash)

Exact:    DEBIT cash_drawer +actual    CREDIT waiter_float +actual
Shortage: DEBIT cash_drawer +actual
          DEBIT cash_shortage +shortage
          CREDIT waiter_float +expected
Overage:  DEBIT cash_drawer +actual
          CREDIT cash_overage +overage
          CREDIT waiter_float +expected
```

---

## Phase 5: Order Revenue (Planned)

| # | Task | Status | Priority | Notes |
|---|------|--------|----------|-------|
| 5.1 | `RecordOrderPayment` trong `JournalService` | ⬜ | LOW | DEBIT cash_drawer / CREDIT customer |
| 5.2 | Wire vào `OrderService.CollectPayment` | ⬜ | LOW | |
| 5.3 | Revenue dashboard từ journal | ⬜ | LOW | Tổng doanh thu = Σ CREDIT customer |

---

## Phase 6: Reporting & Analytics

| # | Task | Status | Priority | Notes |
|---|------|--------|----------|-------|
| 6.1 | Daily P&L report từ journal | ⬜ | MEDIUM | Revenue - Expenses per day |
| 6.2 | Fund flow report theo khoảng thời gian | ⬜ | MEDIUM | |
| 6.3 | Export journal to CSV/Excel | ⬜ | LOW | |
| 6.4 | Shortage/Overage summary report | ⬜ | LOW | cash_shortage, cash_overage totals |

---

## Known Bugs Fixed

| Date | Bug | Fix |
|------|-----|-----|
| 2026-03-11 | `GET /fund/balances` 404 — old endpoint | `getAllBalances()` → `/journal-balances` |
| 2026-03-11 | `WithdrawModal` báo "value must be 0" | `max` attribute dùng `fundBalance` thay `currentBalance` prop |
| 2026-03-11 | Cashier close báo "variance requires documentation" dù variance=0 | expectedCash thiếu `- distributedCash` |
| 2026-03-11 | `OperatingExpenseView` crash "Cannot read properties of null (reading 'find')" | `categories.value = ... || []` |

---

## Immediate Next Steps (Priority Order)

1. **[Phase 4.1–4.4]** Implement waiter handover với shortage/overage journal entries
2. **[Phase 4.5]** Update `CashierHandoverView` để confirm/display variance rõ ràng
3. **[Phase 5.1–5.2]** Record order revenue vào journal (cash_drawer / customer)
4. **[Phase 6.1]** Daily P&L dashboard
