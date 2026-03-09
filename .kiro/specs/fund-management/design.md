# Quản lý Quỹ Tiền — Design Document

## Chiến lược: Logical Separation (Option B)

Không tách thành collection riêng. Thêm field `fund_type` vào `FundTransaction` hiện có.
- Backward compat: transaction cũ không có `fund_type` → mặc định `operating`
- Balance tính bằng cách filter transactions theo `fund_type`
- Không cần migration data lớn

---

## Backend

### 1. Domain: `backend/domain/fund/fund_transaction.go`

**Thêm types mới:**
```go
type FundType string
const (
    FundTypeOperating   FundType = "operating"
    FundTypeInventory   FundType = "inventory"
    FundTypeProfit      FundType = "profit"
    FundTypeCashDrawer  FundType = "cash_drawer"
)

type SourceType string
const (
    SourceTypeExpense      SourceType = "expense"
    SourceTypeIngredient   SourceType = "ingredient"
    SourceTypeHandover     SourceType = "handover"
    SourceTypeManual       SourceType = "manual"
    SourceTypeFundTransfer SourceType = "fund_transfer"
)
```

**Thêm TransactionType:**
```go
TransactionTypeFundHandover TransactionType = "fund_handover"
TransactionTypeFundTransfer TransactionType = "fund_transfer"
TransactionTypeStartingFloat TransactionType = "starting_float"
```

**Thêm fields vào FundTransaction struct:**
```go
FundType    FundType   `bson:"fund_type,omitempty" json:"fund_type,omitempty"`
SourceType  SourceType `bson:"source_type,omitempty" json:"source_type,omitempty"`
SourceID    primitive.ObjectID `bson:"source_id,omitempty" json:"source_id,omitempty"`
Description string     `bson:"description,omitempty" json:"description,omitempty"`
LinkedTxID  primitive.ObjectID `bson:"linked_tx_id,omitempty" json:"linked_tx_id,omitempty"` // cho fund_transfer pairs
```

**Cập nhật `Validate()`:**
- Cho phép `fund_handover`, `fund_transfer`, `starting_float` ngoài `deposit`/`withdrawal`

**Thêm helper:**
```go
func (ft *FundTransaction) SetSource(sourceType SourceType, sourceID primitive.ObjectID) error
func (ft *FundTransaction) SetFundType(fundType FundType)
func (ft *FundTransaction) LinkTo(linkedTxID primitive.ObjectID)
```

---

### 2. Repository: `backend/infrastructure/mongodb/fund_transaction_repository.go`

**Thêm method:**
```go
// FindByFilter: query linh hoạt bằng BSON filter, dùng cho balance per fund_type
func (r *FundTransactionRepository) FindByFilter(
    ctx context.Context,
    filter bson.M,
    limit, offset int,
) ([]*fund.FundTransaction, error)

// CountByFilter: đếm cho pagination
func (r *FundTransactionRepository) CountByFilter(
    ctx context.Context,
    filter bson.M,
) (int64, error)
```

**Thêm index:**
```go
// Index on fund_type
collection.Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys: bson.D{{Key: "fund_type", Value: 1}},
})
// Compound: fund_type + timestamp
collection.Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys: bson.D{
        {Key: "fund_type", Value: 1},
        {Key: "timestamp", Value: -1},
    },
})
// Index on source_type + source_id (sparse)
collection.Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys: bson.D{
        {Key: "source_type", Value: 1},
        {Key: "source_id", Value: 1},
    },
    Options: options.Index().SetSparse(true),
})
```

---

### 3. Service: `backend/application/services/fund_service.go`

#### Cấu trúc dữ liệu mới

```go
type FundTypeBalance struct {
    FundType fund.FundType `json:"fund_type"`
    Balance  fund.FundBalance `json:"balance"`
}

type AllFundBalances struct {
    Operating   fund.FundBalance `json:"operating"`
    Inventory   fund.FundBalance `json:"inventory"`
    Profit      fund.FundBalance `json:"profit"`
    CashDrawer  fund.FundBalance `json:"cash_drawer"`
    Total       fund.FundBalance `json:"total"`
}

type DailyReport struct {
    Date           time.Time          `json:"date"`
    OpeningBalance AllFundBalances    `json:"opening_balance"`
    ClosingBalance AllFundBalances    `json:"closing_balance"`
    TotalInflow    fund.FundBalance   `json:"total_inflow"`
    TotalOutflow   fund.FundBalance   `json:"total_outflow"`
    InflowBySource map[string]fund.FundBalance `json:"inflow_by_source"`
    OutflowBySource map[string]fund.FundBalance `json:"outflow_by_source"`
}

type FundTransferRequest struct {
    FromFundType    fund.FundType
    ToFundType      fund.FundType
    CashAmount      float64
    TransferAmount  float64
    Reason          string
    PerformedBy     primitive.ObjectID
    PerformedByName string
    PerformedByRole string
}

type FundTransferResult struct {
    WithdrawalTx *fund.FundTransaction
    DepositTx    *fund.FundTransaction
    FromBalance  *fund.FundBalance
    ToBalance    *fund.FundBalance
}
```

#### Methods mới

```go
// CalculateBalanceByFundType: tính số dư theo từng quỹ
// Backward compat: operating bao gồm cả transactions không có fund_type
func (s *FundService) CalculateBalanceByFundType(ctx context.Context, fundType fund.FundType) (*fund.FundBalance, error)

// GetAllFundBalances: trả về balance của 4 quỹ + tổng
func (s *FundService) GetAllFundBalances(ctx context.Context) (*AllFundBalances, error)

// TransferBetweenFunds: atomic transfer
// 1. Validate source balance đủ
// 2. Pre-generate 2 ObjectIDs
// 3. Tạo withdrawal từ from-fund (LinkedTxID = depositID)
// 4. Tạo deposit vào to-fund (LinkedTxID = withdrawalID)
// 5. Dùng MongoDB session để đảm bảo atomic
func (s *FundService) TransferBetweenFunds(ctx context.Context, req FundTransferRequest) (*FundTransferResult, error)

// GetDailyReport: báo cáo ngày
// Opening balance = tổng tất cả transactions có timestamp < start of day
// Inflow/outflow = transactions trong ngày, nhóm theo source_type
func (s *FundService) GetDailyReport(ctx context.Context, date time.Time) (*DailyReport, error)
```

#### Cập nhật methods hiện có

```go
// CreateDeposit: thêm tham số fund_type (optional, default "operating")
func (s *FundService) CreateDeposit(
    ctx context.Context,
    cashAmount, transferAmount float64,
    reason string,
    performedBy primitive.ObjectID,
    performedByName, performedByRole string,
    fundType ...fund.FundType, // variadic, default "operating"
) (*fund.FundTransaction, *fund.FundBalance, error)

// CreateWithdrawal: thêm tham số fund_type (optional, default "operating")
func (s *FundService) CreateWithdrawal(
    ctx context.Context,
    cashAmount, transferAmount float64,
    reason string,
    performedBy primitive.ObjectID,
    performedByName, performedByRole string,
    fundType ...fund.FundType, // variadic, default "operating"
) (*fund.FundTransaction, *fund.FundBalance, error)

// CalculateCurrentBalance: vẫn giữ (tính tổng tất cả quỹ, cho backward compat)
```

**Balance calculation logic cho `CalculateBalanceByFundType`:**
```
filter = {fund_type: ft}
// Backward compat: operating bao gồm cả empty/missing
if ft == "operating" {
    filter = {$or: [{fund_type: "operating"}, {fund_type: ""}, {fund_type: {$exists: false}}]}
}

balance = 0
for each transaction matching filter:
    if type in [deposit, fund_handover, starting_float] → add to balance
    if type in [withdrawal] → subtract from balance
    if type == fund_transfer:
        if this tx is deposit side → add
        if this tx is withdrawal side → subtract
```

---

### 4. HTTP Handler: `backend/interfaces/http/fund_handler.go`

**Cập nhật requests:**
```go
type DepositRequest struct {
    CashAmount     float64 `json:"cash_amount"`
    TransferAmount float64 `json:"transfer_amount"`
    Reason         string  `json:"reason"`
    FundType       string  `json:"fund_type"` // optional, default "operating"
}

type WithdrawRequest struct {
    CashAmount     float64 `json:"cash_amount"`
    TransferAmount float64 `json:"transfer_amount"`
    Reason         string  `json:"reason"`
    FundType       string  `json:"fund_type"` // optional, default "operating"
}

type FundTransferRequest struct {
    FromFundType   string  `json:"from_fund_type"`
    ToFundType     string  `json:"to_fund_type"`
    CashAmount     float64 `json:"cash_amount"`
    TransferAmount float64 `json:"transfer_amount"`
    Reason         string  `json:"reason"`
}
```

**Thêm handlers:**
```go
// GET /api/manager/fund/balances → trả về AllFundBalances
func (h *FundHandler) GetAllFundBalances(c *gin.Context)

// POST /api/manager/fund/transfer → FundTransferResult
func (h *FundHandler) Transfer(c *gin.Context)

// GET /api/manager/fund/daily-report?date=2024-01-15 → DailyReport
func (h *FundHandler) GetDailyReport(c *gin.Context)
```

**Cập nhật `GetTransactions`:**
- Thêm query param `fund_type` để filter theo quỹ
- Thêm query param `source_type` để filter theo nguồn

---

### 5. Routes: `backend/main.go`

```go
fundGroup.GET("/balance", fundHandler.GetBalance)       // existing
fundGroup.GET("/balances", fundHandler.GetAllFundBalances) // new
fundGroup.GET("/transactions", fundHandler.GetTransactions)
fundGroup.GET("/transactions/:id", fundHandler.GetTransactionDetail)
fundGroup.POST("/deposit", fundHandler.Deposit)
fundGroup.POST("/withdraw", fundHandler.Withdraw)
fundGroup.POST("/transfer", fundHandler.Transfer) // new
fundGroup.GET("/daily-report", fundHandler.GetDailyReport) // new
```

---

## Frontend

### 1. Constants: `frontend/src/constants/fund.js`

**Thêm:**
```js
export const FUND_TYPES = {
  OPERATING:   'operating',
  INVENTORY:   'inventory',
  PROFIT:      'profit',
  CASH_DRAWER: 'cash_drawer'
}

export const FUND_TYPE_LABELS = {
  operating:   'Quỹ vận hành',
  inventory:   'Quỹ nguyên liệu',
  profit:      'Quỹ lợi nhuận',
  cash_drawer: 'Ngăn kéo tiền'
}

export const FUND_TYPE_ICONS = {
  operating:   '⚙️',
  inventory:   '📦',
  profit:      '💎',
  cash_drawer: '🗄️'
}

export const FUND_TYPE_COLORS = {
  operating:   { bg: 'from-orange-400 to-yellow-400', badge: 'bg-orange-100 text-orange-700' },
  inventory:   { bg: 'from-green-400 to-teal-400',   badge: 'bg-green-100 text-green-700' },
  profit:      { bg: 'from-purple-400 to-pink-400',  badge: 'bg-purple-100 text-purple-700' },
  cash_drawer: { bg: 'from-blue-400 to-cyan-400',    badge: 'bg-blue-100 text-blue-700' }
}

export const SOURCE_TYPE_LABELS = {
  expense:       'Chi phí',
  ingredient:    'Nguyên liệu',
  handover:      'Bàn giao ca',
  manual:        'Thủ công',
  fund_transfer: 'Chuyển quỹ'
}

export const SOURCE_TYPE_ICONS = {
  expense:       '💸',
  ingredient:    '🥤',
  handover:      '🔄',
  manual:        '✍️',
  fund_transfer: '↔️'
}
```

**Thêm vào `TRANSACTION_TYPE_FILTER_OPTIONS`:**
```js
{ value: 'fund_transfer', label: 'Chuyển quỹ' }
```

**Thêm `FUND_TYPE_FILTER_OPTIONS`:**
```js
export const FUND_TYPE_FILTER_OPTIONS = [
  { value: 'all', label: 'Tất cả quỹ' },
  { value: 'operating',   label: '⚙️ Quỹ vận hành' },
  { value: 'inventory',   label: '📦 Quỹ nguyên liệu' },
  { value: 'profit',      label: '💎 Quỹ lợi nhuận' },
  { value: 'cash_drawer', label: '🗄️ Ngăn kéo tiền' },
]
```

---

### 2. Services: `frontend/src/services/fund.js`

**Thêm:**
```js
export const getAllBalances = async () => {
  const { data } = await api.get('/manager/fund/balances')
  return data
}

export const transferBetweenFunds = async (payload) => {
  const { data } = await api.post('/manager/fund/transfer', payload)
  return data
}

export const getDailyReport = async (date) => {
  const { data } = await api.get('/manager/fund/daily-report', { params: { date } })
  return data
}
```

**Cập nhật `getTransactions`:**
- Thêm `fund_type` vào query params

---

### 3. View: `frontend/src/views/FundManagementView.vue`

#### Template layout

```
┌─────────────────────────────────┐
│  Header: 💰 Quản lý quỹ tiền    │
├─────────────────────────────────┤
│  Total Summary Card (orange)    │
│  Cash: X | Transfer: Y | Total  │
├─────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐     │
│  │⚙️ Vận hành│  │📦 Nguyên │     │
│  │  X.XXX   │  │  X.XXX   │     │
│  └──────────┘  └──────────┘     │
│  ┌──────────┐  ┌──────────┐     │
│  │💎 Lợi nhuận│ │🗄️ Ngăn  │     │
│  │  X.XXX   │  │  X.XXX   │     │
│  └──────────┘  └──────────┘     │
├─────────────────────────────────┤
│  [📥 Nạp] [📤 Rút] [🔀 Chuyển] │
├─────────────────────────────────┤
│  Today Summary (inflow/outflow) │
├─────────────────────────────────┤
│  Filters: type | fund | source  │
├─────────────────────────────────┤
│  Transaction List               │
│  - badge nguồn, badge quỹ       │
│  - link đến record gốc          │
└─────────────────────────────────┘
```

#### Script additions

```js
// New refs
const allBalances = ref(null)         // AllFundBalances từ API
const loadingAllBalances = ref(false)
const showTransferModal = ref(false)
const transferring = ref(false)
const transferForm = ref({
  from_fund_type: '',
  to_fund_type: '',
  money_type: 'cash',
  amount: null,
  reason: ''
})

// New computed
const fundTypeKeys = computed(() => Object.values(FUND_TYPES))

// New methods
async function loadAllBalances() { ... }   // GET /balances
async function submitTransfer() { ... }    // POST /transfer, validate, refresh

// Update filters
filters.value.fund_type = 'all'  // add new filter
```

#### Transfer Modal design

```
┌────────────────────────────────┐
│ 🔀 Chuyển tiền giữa quỹ    ×  │
├────────────────────────────────┤
│ Từ quỹ:                        │
│ [⚙️ Vận hành] [📦 Nguyên liệu] │
│ [💎 Lợi nhuận][🗄️ Ngăn kéo]   │
├────────────────────────────────┤
│ Đến quỹ:                       │
│ [⚙️ Vận hành] [📦 Nguyên liệu] │
│ [💎 Lợi nhuận][🗄️ Ngăn kéo]   │
│ (disabled nếu trùng from)      │
├────────────────────────────────┤
│ Loại tiền: [💵 Mặt] [🏦 CK]   │
├────────────────────────────────┤
│ Số tiền: [__________]          │
│ Lý do:   [__________]          │
│ Số dư nguồn: 💵 X | 🏦 Y      │
├────────────────────────────────┤
│      [Hủy]   [🔀 Chuyển quỹ]  │
└────────────────────────────────┘
```

---

## Migration & Backward Compatibility

### Không cần migration data
- Transactions cũ không có `fund_type` field
- `CalculateBalanceByFundType("operating")` filter:
  ```
  {$or: [{fund_type: "operating"}, {fund_type: ""}, {fund_type: {$exists: false}}]}
  ```
- GetAllFundBalances tổng sẽ bằng GetBalance cũ

### API backward compat
- `GET /fund/balance` giữ nguyên (tổng tất cả quỹ)
- `POST /fund/deposit` giữ nguyên (default fund_type=operating)
- `POST /fund/withdraw` giữ nguyên (default fund_type=operating)

---

## Testing Plan

1. **Backend unit tests**
   - `CalculateBalanceByFundType` với data test cho từng quỹ
   - `TransferBetweenFunds` với insufficient balance → error
   - `TransferBetweenFunds` thành công → 2 transactions linked

2. **API integration tests**
   - `GET /balances` → 4 quỹ + total
   - `POST /transfer` → withdrawal + deposit atomic
   - `GET /transactions?fund_type=operating` → filter đúng

3. **Frontend E2E**
   - Load trang → hiển thị 4 thẻ quỹ
   - Chuyển quỹ: fill form → submit → balance cập nhật
   - Filter theo quỹ → chỉ thấy transaction của quỹ đó
