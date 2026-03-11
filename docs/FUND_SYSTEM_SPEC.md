# Fund System — Technical Specification

**Version:** 2.0  
**Last Updated:** 2026-03-11

---

## 1. Architecture

```
Frontend (Vue.js)
  └── services/journal.js          getJournalBalances, getJournalEntries, getJournalEntry
  └── services/fund.js             getAllBalances, getBalance (→ /journal-balances)
  └── constants/fund.js            FUND_TYPES, EVENT_TYPES, labels, icons, colors

Backend
  └── domain/fund/
  │     ├── journal_entry.go       JournalEntry, LedgerLine, JournalEventType, Direction
  │     └── fund_transaction.go    FundType (10 types), SourceType, FundBalance
  │
  └── infrastructure/mongodb/
  │     └── journal_repository.go  Create, FindByID, List, GetFundBalance, GetAllFundBalances
  │
  └── application/services/
  │     ├── journal_service.go     All journal recording methods
  │     └── fund_expense_integration_service.go  Atomic fund+expense operations
  │
  └── interfaces/http/
        └── fund_handler.go        Deposit, Withdraw, Transfer, GetJournalFundBalances, GetJournalEntries, GetJournalEntry
```

---

## 2. Domain Models

### 2.1 JournalEntry
```go
type JournalEntry struct {
    ID              primitive.ObjectID // _id
    EventType       JournalEventType   // cashier_shift_start | cashier_shift_end | waiter_shift_start | waiter_handover | fund_transfer | manager_deposit | manager_withdrawal | expense | ingredient_restock | facility_purchase
    EventID         primitive.ObjectID // ID của entity liên quan (shift, expense, restock...)
    Lines           []LedgerLine       // ≥ 2 lines, Σ DEBIT == Σ CREDIT
    Description     string
    PerformedBy     primitive.ObjectID
    PerformedByName string
    PerformedByRole string
    Timestamp       time.Time
    CreatedAt       time.Time
}
```

### 2.2 LedgerLine
```go
type LedgerLine struct {
    FundType       FundType    // operating | inventory | profit | cash_drawer | waiter_float | owner | supplier | customer | cash_shortage | cash_overage
    Direction      Direction   // debit | credit
    CashAmount     float64
    TransferAmount float64
    TotalAmount    float64     // = CashAmount + TransferAmount
    BalanceBefore  FundBalance // snapshot trước transaction
    BalanceAfter   FundBalance // snapshot sau transaction
}
```

### 2.3 FundBalance
```go
type FundBalance struct {
    Cash     float64
    Transfer float64
    Total    float64
}
```

---

## 3. MongoDB Collection: `journal_entries`

### Indexes
```js
{ event_type: 1, timestamp: -1 }   // Filter by event type
{ "lines.fund_type": 1, timestamp: -1 }  // Filter by fund
{ event_id: 1 }                     // Lookup by source
{ timestamp: -1 }                   // Chronological listing
```

### Balance Aggregation Pipeline
```js
// GetFundBalance(fundType)
[
  { $unwind: "$lines" },
  { $match: { "lines.fund_type": fundType } },
  { $group: {
    _id: "$lines.direction",
    totalCash:     { $sum: "$lines.cash_amount" },
    totalTransfer: { $sum: "$lines.transfer_amount" },
    totalAmount:   { $sum: "$lines.total_amount" }
  }},
  // balance = debit - credit
]
```

---

## 4. JournalService API

```go
// Balance queries
GetFundBalance(ctx, fundType) (*FundBalance, error)
GetAllFundBalances(ctx) (map[FundType]*FundBalance, error)

// Manager operations
RecordManagerDeposit(ctx, fundType, cashAmt, transferAmt, reason, by, name, role) (*JournalEntry, error)
RecordManagerWithdrawal(ctx, fundType, cashAmt, transferAmt, reason, by, name, role) (*JournalEntry, error)
RecordFundTransfer(ctx, fromFund, toFund, cashAmt, transferAmt, reason, by, name, role) (*JournalEntry, error)

// Cashier shift
RecordCashierShiftStart(ctx, startingFloat, cashierShiftID, cashierID, cashierName) (*JournalEntry, error)
RecordCashierShiftEnd(ctx, cashAmt, transferAmt, cashierShiftID, cashierID, cashierName) (*JournalEntry, error)

// Waiter shift
RecordWaiterShiftStart(ctx, amount, shiftID, waiterID, waiterName) (*JournalEntry, error)  // TODO
RecordWaiterHandover(ctx, actualCash, expectedCash, shiftID, handoverID, waiterID, waiterName) (*JournalEntry, error)  // TODO (handles shortage/overage)

// Fund-paid expenses (used inside MongoDB session)
RecordFundWithdrawalInSession(
    ctx SessionContext,
    eventType JournalEventType,
    fundType FundType,
    counterpart FundType,   // supplier for purchases, owner for withdrawals
    cashAmt, transferAmt float64,
    description string,
    eventID ObjectID,
    performedBy ObjectID, name, role string,
) (*JournalEntry, error)

// Listing
ListEntries(ctx, JournalFilter) ([]*JournalEntry, int64, error)
GetEntry(ctx, id) (*JournalEntry, error)
```

---

## 5. JournalFilter
```go
type JournalFilter struct {
    FundType  *FundType          // filter by any line's fund_type
    EventType *JournalEventType  // filter by event_type
    FromDate  *time.Time
    ToDate    *time.Time
    Limit     int                // default 20
    Offset    int
}
```

---

## 6. HTTP Handler Request/Response

### POST /api/manager/fund/deposit
```json
Request:  { "fund_type": "operating", "cash_amount": 500000, "transfer_amount": 0, "reason": "Nạp tiền đầu tháng" }
Response: JournalEntry
```

### POST /api/manager/fund/withdraw
```json
Request:  { "fund_type": "operating", "cash_amount": 200000, "transfer_amount": 0, "reason": "Rút tiền chi tiêu" }
Response: JournalEntry
```

### POST /api/manager/fund/transfer
```json
Request:  { "from_fund_type": "operating", "to_fund_type": "inventory", "cash_amount": 300000, "transfer_amount": 0, "reason": "Chuyển quỹ mua hàng" }
Response: JournalEntry
```

### GET /api/manager/fund/journal-balances
```json
Response: {
  "operating":    { "cash": 500000, "transfer": 0, "total": 500000 },
  "inventory":    { "cash": 200000, "transfer": 100000, "total": 300000 },
  "profit":       { "cash": 0, "transfer": 0, "total": 0 },
  "cash_drawer":  { "cash": 150000, "transfer": 0, "total": 150000 },
  "waiter_float": { "cash": 50000, "transfer": 0, "total": 50000 }
}
```

### GET /api/manager/fund/journal
```json
Query params: fund_type, event_type, from_date (RFC3339), to_date (RFC3339), limit (default 20), offset
Response: { "entries": [...JournalEntry], "total": 42 }
```

---

## 7. Frontend Constants (fund.js)

```js
// Event types
EVENT_TYPES.CASHIER_SHIFT_START = 'cashier_shift_start'
// ... (10 types total)

// Fund types
FUND_TYPES.OPERATING = 'operating'
FUND_TYPES.CASH_SHORTAGE = 'cash_shortage'
FUND_TYPES.CASH_OVERAGE = 'cash_overage'
// ... (10 types total)

// Display helpers
FUND_TYPE_LABELS   // { operating: 'Quỹ vận hành', ... }
FUND_TYPE_ICONS    // { operating: '⚙️', cash_shortage: '📉', ... }
FUND_TYPE_COLORS   // { operating: 'blue', cash_shortage: 'red', cash_overage: 'emerald', ... }
EVENT_TYPE_LABELS  // { cashier_shift_start: 'Đầu ca thu ngân', ... }
INFLOW_EVENTS      // Set: manager_deposit, cashier_shift_end, waiter_handover
OUTFLOW_EVENTS     // Set: manager_withdrawal, cashier_shift_start, waiter_shift_start, expense, ...
```

---

## 8. FundManagementView — Display Logic

```js
// Filter out external counterpart accounts from main line display
const EXTERNAL_ACCOUNTS = new Set(['external', 'owner', 'supplier', 'customer', 'cash_shortage', 'cash_overage'])
const realLines = (entry) => entry.lines.filter(l => !EXTERNAL_ACCOUNTS.has(l.fund_type))

// Amount = sum of debit real lines (for non-transfer events)
const entryAmount = (entry) => realLines(entry)
  .filter(l => l.direction === 'debit')
  .reduce((sum, l) => sum + l.total_amount, 0)

// Color: green=inflow, blue=transfer, red=outflow
const entryAmountColor = (entry) =>
  INFLOW_EVENTS.has(entry.event_type) ? 'text-green-600' :
  entry.event_type === 'fund_transfer' ? 'text-blue-600' : 'text-red-600'
```

---

## 9. Waiter Handover — Pending Implementation (Phase 4)

### RecordWaiterHandover — Target Implementation
```go
func (s *JournalService) RecordWaiterHandover(
    ctx context.Context,
    actualCash float64,
    expectedCash float64,
    shiftID primitive.ObjectID,
    waiterID primitive.ObjectID,
    waiterName string,
) (*JournalEntry, error) {
    variance := actualCash - expectedCash
    // shortage = expected > actual → variance < 0
    // overage  = actual > expected → variance > 0
    // exact    = variance == 0 (within tolerance)

    if abs(variance) < 0.01 {
        // Exact: DEBIT cash_drawer / CREDIT waiter_float
    } else if variance < 0 {
        // Shortage: DEBIT cash_drawer(actual) + DEBIT cash_shortage(|variance|) / CREDIT waiter_float(expected)
    } else {
        // Overage: DEBIT cash_drawer(actual) / CREDIT cash_overage(variance) + CREDIT waiter_float(expected)
    }
}
```
