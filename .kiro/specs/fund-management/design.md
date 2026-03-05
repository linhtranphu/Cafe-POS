# Fund Management - Design Document

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     Frontend (Vue.js)                        │
│  ┌────────────────────────────────────────────────────────┐ │
│  │           FundManagementView.vue                       │ │
│  │  ┌──────────────┐  ┌──────────────┐  ┌─────────────┐ │ │
│  │  │ Balance Card │  │ Transaction  │  │   Filters   │ │ │
│  │  └──────────────┘  │     List     │  └─────────────┘ │ │
│  │                    └──────────────┘                    │ │
│  │  ┌──────────────┐  ┌──────────────┐                  │ │
│  │  │Deposit Modal │  │Withdraw Modal│                  │ │
│  │  └──────────────┘  └──────────────┘                  │ │
│  └────────────────────────────────────────────────────────┘ │
│                           │                                  │
│                    fundService.js                            │
└───────────────────────────┼──────────────────────────────────┘
                            │ HTTP/REST
┌───────────────────────────┼──────────────────────────────────┐
│                     Backend (Go)                             │
│  ┌────────────────────────┴────────────────────────────────┐│
│  │         FundHandler (HTTP Layer)                        ││
│  │  - GetBalance()                                         ││
│  │  - GetTransactions()                                    ││
│  │  - Deposit()                                            ││
│  │  - Withdraw()                                           ││
│  └────────────────────────┬────────────────────────────────┘│
│                            │                                 │
│  ┌────────────────────────┴────────────────────────────────┐│
│  │         FundService (Business Logic)                    ││
│  │  - CalculateBalance()                                   ││
│  │  - GetTransactionHistory()                              ││
│  │  - CreateDeposit()                                      ││
│  │  - CreateWithdrawal()                                   ││
│  └────────────────────────┬────────────────────────────────┘│
│                            │                                 │
│  ┌────────────────────────┴────────────────────────────────┐│
│  │         FundTransactionRepository                       ││
│  │         CashHandoverRepository                          ││
│  │         FundHandoverRepository                          ││
│  │         CashierShiftRepository                          ││
│  └────────────────────────┬────────────────────────────────┘│
└────────────────────────────┼─────────────────────────────────┘
                             │
                    ┌────────┴────────┐
                    │    MongoDB      │
                    │  - fund_trans   │
                    │  - cash_hand    │
                    │  - fund_hand    │
                    │  - cashier_sh   │
                    └─────────────────┘
```

## Data Model

### FundTransaction (New Collection)

```go
type FundTransaction struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    Type        string            // "deposit" | "withdrawal"
    CashAmount  float64           
    TransferAmount float64        
    TotalAmount float64           
    Reason      string            // Required, min 10 chars
    PerformedBy primitive.ObjectID // User ID
    PerformedByName string        
    PerformedByRole string        
    Timestamp   time.Time         
    
    // Balance tracking (optional, for audit)
    BalanceBefore *FundBalance   `bson:"balance_before,omitempty"`
    BalanceAfter  *FundBalance   `bson:"balance_after,omitempty"`
    
    CreatedAt   time.Time         
    UpdatedAt   time.Time         
}

type FundBalance struct {
    Cash     float64
    Transfer float64
    Total    float64
}
```

### Transaction History View Model

```go
type TransactionHistoryItem struct {
    ID              string
    Type            string // "deposit", "withdrawal", "cash_handover", "fund_handover", "starting_float"
    CashAmount      float64
    TransferAmount  float64
    TotalAmount     float64
    Description     string
    PerformedBy     string
    PerformedByRole string
    Timestamp       time.Time
    
    // Type-specific metadata
    Metadata map[string]interface{}
}
```

## API Design

### GET /api/manager/fund/balance

**Response:**
```json
{
  "current_balance": {
    "cash": 5000000,
    "transfer": 3000000,
    "total": 8000000
  },
  "today_summary": {
    "opening_balance": {
      "cash": 2000000,
      "transfer": 1000000,
      "total": 3000000
    },
    "total_inflow": {
      "cash": 4000000,
      "transfer": 2500000,
      "total": 6500000
    },
    "total_outflow": {
      "cash": 1000000,
      "transfer": 500000,
      "total": 1500000
    }
  },
  "last_updated": "2026-02-24T10:30:00Z"
}
```

### GET /api/manager/fund/transactions

**Query Parameters:**
- `type`: "deposit" | "withdrawal" | "handover" | "all" (default: "all")
- `money_type`: "cash" | "transfer" | "all" (default: "all")
- `from_date`: ISO date (default: today 00:00)
- `to_date`: ISO date (default: now)
- `limit`: int (default: 50, max: 200)
- `offset`: int (default: 0)

**Response:**
```json
{
  "transactions": [
    {
      "id": "...",
      "type": "fund_handover",
      "cash_amount": 250000,
      "transfer_amount": 150000,
      "total_amount": 400000,
      "description": "Bàn giao từ ca thu ngân #1234",
      "performed_by": "Nguyễn Văn A",
      "performed_by_role": "cashier",
      "timestamp": "2026-02-24T10:00:00Z",
      "metadata": {
        "cashier_shift_id": "...",
        "variance_amount": 0
      }
    },
    {
      "id": "...",
      "type": "withdrawal",
      "cash_amount": 500000,
      "transfer_amount": 0,
      "total_amount": 500000,
      "description": "Rút tiền mua nguyên liệu",
      "performed_by": "Manager",
      "performed_by_role": "manager",
      "timestamp": "2026-02-24T09:00:00Z",
      "metadata": {}
    }
  ],
  "total": 45,
  "limit": 50,
  "offset": 0
}
```

### POST /api/manager/fund/deposit

**Request:**
```json
{
  "cash_amount": 1000000,
  "transfer_amount": 0,
  "reason": "Bổ sung vốn đầu ngày"
}
```

**Response:**
```json
{
  "transaction": {
    "id": "...",
    "type": "deposit",
    "cash_amount": 1000000,
    "transfer_amount": 0,
    "total_amount": 1000000,
    "reason": "Bổ sung vốn đầu ngày",
    "performed_by": "Manager",
    "timestamp": "2026-02-24T08:00:00Z"
  },
  "new_balance": {
    "cash": 3000000,
    "transfer": 1000000,
    "total": 4000000
  }
}
```

### POST /api/manager/fund/withdraw

**Request:**
```json
{
  "cash_amount": 500000,
  "transfer_amount": 0,
  "reason": "Mua nguyên liệu"
}
```

**Response:**
```json
{
  "transaction": {
    "id": "...",
    "type": "withdrawal",
    "cash_amount": 500000,
    "transfer_amount": 0,
    "total_amount": 500000,
    "reason": "Mua nguyên liệu",
    "performed_by": "Manager",
    "timestamp": "2026-02-24T09:00:00Z"
  },
  "new_balance": {
    "cash": 2500000,
    "transfer": 1000000,
    "total": 3500000
  }
}
```

## UI Design

### Layout Structure

```
┌─────────────────────────────────────┐
│ Header: 💰 Quản lý quỹ tiền         │
├─────────────────────────────────────┤
│ ┌─────────────────────────────────┐ │
│ │   Current Balance Card          │ │
│ │   💵 Cash: 5,000,000            │ │
│ │   💳 Transfer: 3,000,000        │ │
│ │   📊 Total: 8,000,000           │ │
│ └─────────────────────────────────┘ │
│                                     │
│ ┌─────────────────────────────────┐ │
│ │   Today Summary                 │ │
│ │   📥 Inflow: +6,500,000         │ │
│ │   📤 Outflow: -1,500,000        │ │
│ └─────────────────────────────────┘ │
│                                     │
│ [📥 Thêm tiền] [📤 Rút tiền]       │
│                                     │
│ ┌─────────────────────────────────┐ │
│ │   Filters                       │ │
│ │   [Hôm nay ▼] [Tất cả ▼]       │ │
│ └─────────────────────────────────┘ │
│                                     │
│ ┌─────────────────────────────────┐ │
│ │ 📥 Bàn giao từ cashier          │ │
│ │ +400,000 | 10:00 | Nguyễn Văn A│ │
│ └─────────────────────────────────┘ │
│ ┌─────────────────────────────────┐ │
│ │ 📤 Rút tiền                     │ │
│ │ -500,000 | 09:00 | Manager     │ │
│ └─────────────────────────────────┘ │
│                                     │
│ Bottom Navigation                   │
└─────────────────────────────────────┘
```

### Color Scheme
- Primary: Orange/Yellow (fund theme)
- Success/Inflow: Green
- Danger/Outflow: Red
- Info: Blue (transfer)
- Cash: Green tones
- Transfer: Blue tones

## Implementation Phases

### Phase 1: Backend Foundation
1. Create FundTransaction domain model
2. Create FundTransactionRepository
3. Create FundService with balance calculation
4. Create HTTP handlers
5. Add routes with manager auth

### Phase 2: Transaction History
1. Implement GetTransactionHistory with filters
2. Aggregate data from multiple collections
3. Add pagination support

### Phase 3: Deposit/Withdrawal
1. Implement CreateDeposit
2. Implement CreateWithdrawal with balance validation
3. Add transaction atomicity

### Phase 4: Frontend View
1. Create FundManagementView
2. Create fundService
3. Implement balance display
4. Implement transaction list
5. Add filters

### Phase 5: Deposit/Withdrawal UI
1. Create DepositModal
2. Create WithdrawalModal
3. Add form validation
4. Integrate with backend

### Phase 6: Polish & Testing
1. Add pull-to-refresh
2. Add loading states
3. Add error handling
4. Mobile optimization
5. Testing

## Security Considerations

1. **Authorization**: Only manager role can access
2. **Validation**: 
   - Amount > 0
   - Withdrawal <= current balance
   - Reason min 10 characters
3. **Audit**: Log all transactions with user info
4. **Atomicity**: Use MongoDB transactions for balance updates
5. **Rate Limiting**: Prevent abuse of deposit/withdrawal endpoints

## Performance Considerations

1. **Indexing**:
   - `fund_transactions`: timestamp (desc), type
   - Compound index: (timestamp, type, performed_by)
2. **Caching**: Cache current balance (invalidate on transaction)
3. **Pagination**: Limit transaction list to 50-200 items
4. **Aggregation**: Use MongoDB aggregation pipeline for balance calculation
