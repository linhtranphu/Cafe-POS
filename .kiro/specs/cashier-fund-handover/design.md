# Design Document - Cashier Fund Handover

## Overview

Tính năng này cho phép Cashier xem và quản lý số tiền đang chịu trách nhiệm, và handover lại số tiền này về quỹ khi đóng ca. Thiết kế tập trung vào:

1. **Hiển thị rõ ràng**: Cashier luôn biết mình đang quản lý bao nhiêu tiền
2. **Quy trình đơn giản**: Handover về quỹ là một phần tự nhiên của việc đóng ca
3. **Tính toàn vẹn**: Sử dụng transaction để đảm bảo atomicity
4. **Khả năng mở rộng**: Thiết kế sẵn để sau này có thể chỉ định người nhận (manager)

## Architecture

### System Components

```
┌─────────────────────────────────────────────────────────────┐
│                     Frontend Layer                          │
├─────────────────────────────────────────────────────────────┤
│  CashierDashboard                                           │
│  - Display managed funds section                            │
│  - Show received cash + transfer                            │
│  - Warning about responsibility                             │
│                                                             │
│  CashierShiftClosureV2                                      │
│  - Display managed funds summary                            │
│  - Cash counting interface                                  │
│  - Variance documentation                                   │
│  - Fund handover confirmation                               │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                     API Layer                               │
├─────────────────────────────────────────────────────────────┤
│  GET /api/cashier/shifts/:id/managed-funds                  │
│  - Get current managed funds summary                        │
│                                                             │
│  POST /api/cashier/shifts/:id/close                         │
│  - Close shift with fund handover                           │
│  - Includes actual_cash, variance, notes                    │
│                                                             │
│  GET /api/cashier/fund-handovers                            │
│  - Query fund handover history                              │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                  Service Layer                              │
├─────────────────────────────────────────────────────────────┤
│  CashierShiftService                                        │
│  - GetManagedFunds()                                        │
│  - CloseShiftWithFundHandover()                             │
│  - CreateFundHandoverRecord()                               │
│  - CalculateVariance()                                      │
│  - ValidateVarianceDocumentation()                          │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   Domain Layer                              │
├─────────────────────────────────────────────────────────────┤
│  CashierShift (Extended)                                    │
│  - ReceivedCash                                             │
│  - ReceivedTransfer                                         │
│  - TotalManagedFunds()                                      │
│                                                             │
│  FundHandover (NEW)                                         │
│  - CashierShiftID                                           │
│  - CashAmount                                               │
│  - TransferAmount                                           │
│  - VarianceAmount                                           │
│  - ReceiverID (nullable)                                    │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                  Repository Layer                           │
├─────────────────────────────────────────────────────────────┤
│  MongoDB Collections:                                       │
│  - cashier_shifts (existing, with received amounts)         │
│  - fund_handovers (NEW)                                     │
└─────────────────────────────────────────────────────────────┘
```

## Data Models

### FundHandover (NEW Domain Model)

```go
package cashier

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// FundHandover represents a handover of funds from cashier back to the fund
// when closing a cashier shift. This is the final step in the cashier's
// responsibility chain.
type FundHandover struct {
    // ID is the unique identifier for the fund handover
    ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    
    // CashierShiftID links this handover to the cashier shift
    CashierShiftID primitive.ObjectID `json:"cashier_shift_id" bson:"cashier_shift_id"`
    
    // CashierID is the cashier who is handing over the funds
    CashierID primitive.ObjectID `json:"cashier_id" bson:"cashier_id"`
    
    // CashierName for display purposes
    CashierName string `json:"cashier_name" bson:"cashier_name"`
    
    // CashAmount is the actual cash amount handed over
    CashAmount float64 `json:"cash_amount" bson:"cash_amount"`
    
    // TransferAmount is the transfer amount recorded (no physical handover)
    TransferAmount float64 `json:"transfer_amount" bson:"transfer_amount"`
    
    // TotalAmount is the sum of cash and transfer
    TotalAmount float64 `json:"total_amount" bson:"total_amount"`
    
    // ExpectedCash is the theoretical cash amount (starting_float + received_cash)
    ExpectedCash float64 `json:"expected_cash" bson:"expected_cash"`
    
    // VarianceAmount is the difference between actual and expected cash
    VarianceAmount float64 `json:"variance_amount" bson:"variance_amount"`
    
    // VarianceReason is the selected reason for variance (if any)
    VarianceReason *VarianceReason `json:"variance_reason,omitempty" bson:"variance_reason,omitempty"`
    
    // VarianceNotes provides detailed explanation for variance
    VarianceNotes string `json:"variance_notes,omitempty" bson:"variance_notes,omitempty"`
    
    // ReceiverID is the person receiving the funds (nullable for future use)
    ReceiverID *primitive.ObjectID `json:"receiver_id,omitempty" bson:"receiver_id,omitempty"`
    
    // ReceiverName for display purposes (nullable)
    ReceiverName string `json:"receiver_name,omitempty" bson:"receiver_name,omitempty"`
    
    // HandoverAt is when the handover occurred
    HandoverAt time.Time `json:"handover_at" bson:"handover_at"`
    
    // CreatedAt is when the record was created
    CreatedAt time.Time `json:"created_at" bson:"created_at"`
    
    // UpdatedAt is when the record was last updated
    UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// NewFundHandover creates a new fund handover record
func NewFundHandover(
    cashierShiftID, cashierID primitive.ObjectID,
    cashierName string,
    cashAmount, transferAmount, expectedCash float64,
) *FundHandover {
    now := time.Now()
    variance := cashAmount - expectedCash
    
    return &FundHandover{
        CashierShiftID: cashierShiftID,
        CashierID:      cashierID,
        CashierName:    cashierName,
        CashAmount:     cashAmount,
        TransferAmount: transferAmount,
        TotalAmount:    cashAmount + transferAmount,
        ExpectedCash:   expectedCash,
        VarianceAmount: variance,
        HandoverAt:     now,
        CreatedAt:      now,
        UpdatedAt:      now,
    }
}

// HasVariance returns true if there is a variance
func (fh *FundHandover) HasVariance() bool {
    return fh.VarianceAmount != 0
}

// DocumentVariance adds variance documentation
func (fh *FundHandover) DocumentVariance(reason VarianceReason, notes string) error {
    if !fh.HasVariance() {
        return errors.New("no variance to document")
    }
    
    if len(notes) < 10 {
        return errors.New("variance notes must be at least 10 characters")
    }
    
    fh.VarianceReason = &reason
    fh.VarianceNotes = notes
    fh.UpdatedAt = time.Now()
    
    return nil
}

// SetReceiver sets the receiver information (for future use)
func (fh *FundHandover) SetReceiver(receiverID primitive.ObjectID, receiverName string) {
    fh.ReceiverID = &receiverID
    fh.ReceiverName = receiverName
    fh.UpdatedAt = time.Now()
}
```

### Database Schema

#### fund_handovers Collection (NEW)

```javascript
{
  _id: ObjectId,
  cashier_shift_id: ObjectId,
  cashier_id: ObjectId,
  cashier_name: String,
  
  cash_amount: Number,           // Actual cash handed over
  transfer_amount: Number,       // Transfer amount recorded
  total_amount: Number,          // cash + transfer
  
  expected_cash: Number,         // starting_float + received_cash
  variance_amount: Number,       // actual - expected
  variance_reason: String,       // Optional: reason enum
  variance_notes: String,        // Optional: detailed notes
  
  receiver_id: ObjectId,         // Optional: for future use
  receiver_name: String,         // Optional: for future use
  
  handover_at: Date,
  created_at: Date,
  updated_at: Date
}
```

**Indexes:**
```javascript
{ cashier_shift_id: 1 }        // Unique - one handover per shift
{ cashier_id: 1, handover_at: -1 }  // Query by cashier
{ handover_at: -1 }            // Query by date
{ variance_amount: 1 }         // Find handovers with variance
```

## UI Design

### 1. CashierDashboard - Managed Funds Section

```
┌─────────────────────────────────────────────────────────────┐
│ 💵 Thu ngân                                                 │
│ Giám sát & đối soát                                         │
│                                                             │
│ [Current Shift Info Card]                                   │
│                                                             │
│ ┌─────────────────────────────────────────────────────────┐│
│ │ 💰 Tiền đang quản lý                                    ││
│ │                                                         ││
│ │ ┌─────────────────┬─────────────────┐                  ││
│ │ │ 💵 Tiền mặt     │ 💳 Tiền CK      │                  ││
│ │ │ 1,500,000₫      │ 800,000₫        │                  ││
│ │ │ Đã nhận         │ Đã nhận         │                  ││
│ │ └─────────────────┴─────────────────┘                  ││
│ │                                                         ││
│ │ ┌───────────────────────────────────┐                  ││
│ │ │ 📊 Tổng cộng                      │                  ││
│ │ │ 2,300,000₫                        │                  ││
│ │ └───────────────────────────────────┘                  ││
│ │                                                         ││
│ │ ⚠️ Bạn chịu trách nhiệm trên số tiền này              ││
│ │ Khi đóng ca, bạn cần bàn giao lại về quỹ              ││
│ └─────────────────────────────────────────────────────────┘│
│                                                             │
│ [Handover Management Button]                                │
│ [Pending Discrepancies]                                     │
└─────────────────────────────────────────────────────────────┘
```

### 2. CashierShiftClosureV2 - Fund Handover Flow

#### Step 1: Display Managed Funds Summary

```
┌─────────────────────────────────────────────────────────────┐
│ 🔒 Đóng ca thu ngân                                         │
│                                                             │
│ Bước 1: Xem tổng quan                                       │
│                                                             │
│ ┌─────────────────────────────────────────────────────────┐│
│ │ 💰 Số tiền bạn đang quản lý                             ││
│ │                                                         ││
│ │ Tiền đầu ca:           500,000₫                         ││
│ │ Nhận từ waiter (mặt):  1,500,000₫                       ││
│ │ Nhận từ waiter (CK):   800,000₫                         ││
│ │ ─────────────────────────────────                       ││
│ │ Tổng tiền mặt lý thuyết: 2,000,000₫                     ││
│ │ Tiền CK ghi nhận:      800,000₫                         ││
│ └─────────────────────────────────────────────────────────┘│
│                                                             │
│ [Tiếp tục →]                                                │
└─────────────────────────────────────────────────────────────┘
```

#### Step 2: Count Cash

```
┌─────────────────────────────────────────────────────────────┐
│ 🔒 Đóng ca thu ngân                                         │
│                                                             │
│ Bước 2: Đếm tiền mặt                                        │
│                                                             │
│ ┌─────────────────────────────────────────────────────────┐│
│ │ 💵 Đếm tiền mặt thực tế                                 ││
│ │                                                         ││
│ │ Số tiền lý thuyết: 2,000,000₫                           ││
│ │                                                         ││
│ │ Nhập số tiền thực tế đếm được:                          ││
│ │ ┌─────────────────────────────────┐                    ││
│ │ │ [        1,995,000        ] ₫   │                    ││
│ │ └─────────────────────────────────┘                    ││
│ │                                                         ││
│ │ ⚠️ Chênh lệch: -5,000₫ (Thiếu)                         ││
│ └─────────────────────────────────────────────────────────┘│
│                                                             │
│ [← Quay lại]  [Tiếp tục →]                                 │
└─────────────────────────────────────────────────────────────┘
```

#### Step 3: Document Variance (if exists)

```
┌─────────────────────────────────────────────────────────────┐
│ 🔒 Đóng ca thu ngân                                         │
│                                                             │
│ Bước 3: Giải thích chênh lệch                               │
│                                                             │
│ ┌─────────────────────────────────────────────────────────┐│
│ │ ⚠️ Có chênh lệch: -5,000₫                               ││
│ │                                                         ││
│ │ Chọn lý do: *                                           ││
│ │ ┌─────────────────────────────────┐                    ││
│ │ │ [Lỗi đếm tiền          ▼]       │                    ││
│ │ └─────────────────────────────────┘                    ││
│ │                                                         ││
│ │ Ghi chú chi tiết: *                                     ││
│ │ ┌─────────────────────────────────┐                    ││
│ │ │ Đếm nhầm tờ 50k thành 100k      │                    ││
│ │ │                                 │                    ││
│ │ └─────────────────────────────────┘                    ││
│ │ (Tối thiểu 10 ký tự)                                    ││
│ └─────────────────────────────────────────────────────────┘│
│                                                             │
│ [← Quay lại]  [Tiếp tục →]                                 │
└─────────────────────────────────────────────────────────────┘
```

#### Step 4: Confirm Fund Handover

```
┌─────────────────────────────────────────────────────────────┐
│ 🔒 Đóng ca thu ngân                                         │
│                                                             │
│ Bước 4: Xác nhận bàn giao về quỹ                            │
│                                                             │
│ ┌─────────────────────────────────────────────────────────┐│
│ │ 📋 Tóm tắt bàn giao                                     ││
│ │                                                         ││
│ │ 💵 Tiền mặt bàn giao:    1,995,000₫                     ││
│ │ 💳 Tiền CK ghi nhận:     800,000₫                       ││
│ │ ─────────────────────────────────                       ││
│ │ 📊 Tổng cộng:            2,795,000₫                     ││
│ │                                                         ││
│ │ ⚠️ Chênh lệch:           -5,000₫                        ││
│ │ Lý do: Lỗi đếm tiền                                     ││
│ │                                                         ││
│ │ ✅ Xác nhận bàn giao về quỹ                             ││
│ │ Sau khi xác nhận, ca sẽ được đóng                       ││
│ └─────────────────────────────────────────────────────────┘│
│                                                             │
│ [← Quay lại]  [✓ Xác nhận và đóng ca]                      │
└─────────────────────────────────────────────────────────────┘
```

## API Design

### 1. Get Managed Funds

```
GET /api/cashier/shifts/:id/managed-funds

Response:
{
  "cashier_shift_id": "...",
  "starting_float": 500000,
  "received_cash": 1500000,
  "received_transfer": 800000,
  "total_managed_funds": 2300000,
  "expected_cash": 2000000,
  "handover_count": 5
}
```

### 2. Close Shift with Fund Handover

```
POST /api/cashier/shifts/:id/close

Request:
{
  "actual_cash": 1995000,
  "variance_reason": "COUNTING_ERROR",
  "variance_notes": "Đếm nhầm tờ 50k thành 100k",
  "receiver_id": null  // Optional, for future use
}

Response:
{
  "cashier_shift": { ... },
  "fund_handover": {
    "id": "...",
    "cash_amount": 1995000,
    "transfer_amount": 800000,
    "total_amount": 2795000,
    "variance_amount": -5000,
    "variance_reason": "COUNTING_ERROR",
    "variance_notes": "Đếm nhầm tờ 50k thành 100k",
    "handover_at": "2024-01-15T18:30:00Z"
  }
}
```

### 3. Get Fund Handover History

```
GET /api/cashier/fund-handovers?cashier_id=...&from=...&to=...

Response:
{
  "handovers": [
    {
      "id": "...",
      "cashier_name": "Nguyễn Văn A",
      "cash_amount": 1995000,
      "transfer_amount": 800000,
      "total_amount": 2795000,
      "variance_amount": -5000,
      "handover_at": "2024-01-15T18:30:00Z"
    },
    ...
  ],
  "total": 10,
  "page": 1,
  "page_size": 20
}
```

## Service Implementation

### CashierShiftService Extension

```go
// GetManagedFunds returns the current managed funds for a cashier shift
func (s *CashierShiftService) GetManagedFunds(
    ctx context.Context,
    shiftID primitive.ObjectID,
) (*ManagedFundsSummary, error) {
    shift, err := s.cashierShiftRepo.FindByID(ctx, shiftID)
    if err != nil {
        return nil, err
    }
    
    return &ManagedFundsSummary{
        CashierShiftID:     shift.ID,
        StartingFloat:      shift.StartingFloat,
        ReceivedCash:       shift.ReceivedCash,
        ReceivedTransfer:   shift.ReceivedTransfer,
        TotalManagedFunds:  shift.ReceivedCash + shift.ReceivedTransfer,
        ExpectedCash:       shift.StartingFloat + shift.ReceivedCash,
        HandoverCount:      shift.HandoverCount,
    }, nil
}

// CloseShiftWithFundHandover closes a cashier shift and creates fund handover record
func (s *CashierShiftService) CloseShiftWithFundHandover(
    ctx context.Context,
    shiftID primitive.ObjectID,
    actualCash float64,
    varianceReason *cashier.VarianceReason,
    varianceNotes string,
    receiverID *primitive.ObjectID,
    userID, deviceID string,
) (*cashier.CashierShift, *cashier.FundHandover, error) {
    // Start transaction
    session, err := s.mongoClient.StartSession()
    if err != nil {
        return nil, nil, err
    }
    defer session.EndSession(ctx)
    
    var resultShift *cashier.CashierShift
    var resultHandover *cashier.FundHandover
    
    _, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
        // 1. Get cashier shift
        shift, err := s.cashierShiftRepo.FindByID(sessCtx, shiftID)
        if err != nil {
            return nil, err
        }
        
        // 2. Calculate expected cash
        expectedCash := shift.StartingFloat + shift.ReceivedCash
        
        // 3. Create fund handover record
        handover := cashier.NewFundHandover(
            shift.ID,
            shift.CashierID,
            shift.CashierName,
            actualCash,
            shift.ReceivedTransfer,
            expectedCash,
        )
        
        // 4. Document variance if exists
        if handover.HasVariance() {
            if varianceReason == nil || varianceNotes == "" {
                return nil, errors.New("variance documentation required")
            }
            if err := handover.DocumentVariance(*varianceReason, varianceNotes); err != nil {
                return nil, err
            }
        }
        
        // 5. Set receiver if provided
        if receiverID != nil {
            // TODO: Get receiver name from user service
            handover.SetReceiver(*receiverID, "Manager Name")
        }
        
        // 6. Save fund handover
        if err := s.fundHandoverRepo.Create(sessCtx, handover); err != nil {
            return nil, err
        }
        
        // 7. Record actual cash in shift (for variance tracking)
        if err := shift.RecordActualCash(actualCash, userID, deviceID, time.Now()); err != nil {
            return nil, err
        }
        
        // 8. Close the shift
        if err := shift.Close(userID, deviceID, time.Now()); err != nil {
            return nil, err
        }
        
        // 9. Save shift
        if err := s.cashierShiftRepo.Save(sessCtx, shift); err != nil {
            return nil, err
        }
        
        resultShift = shift
        resultHandover = handover
        return nil, nil
    })
    
    if err != nil {
        return nil, nil, err
    }
    
    return resultShift, resultHandover, nil
}
```

## Implementation Steps

### Phase 1: Backend Foundation
1. Create FundHandover domain model
2. Create FundHandoverRepository
3. Extend CashierShiftService with fund handover methods
4. Add API endpoints
5. Write unit tests

### Phase 2: Frontend - Dashboard
1. Add managed funds section to CashierDashboard
2. Fetch and display received cash/transfer
3. Add warning message about responsibility
4. Style with gradient cards

### Phase 3: Frontend - Closure Flow
1. Extend CashierShiftClosureV2 with fund handover steps
2. Add managed funds summary display
3. Add cash counting interface
4. Add variance documentation form
5. Add final confirmation screen
6. Integrate with backend API

### Phase 4: Testing & Refinement
1. Manual testing of complete flow
2. Test variance scenarios
3. Test transaction rollback
4. Performance testing
5. UI/UX refinement

### Phase 5: Documentation
1. Update user guide
2. Create training materials
3. Document API
4. Update system architecture docs

## Testing Strategy

### Unit Tests
- FundHandover domain model methods
- Service layer business logic
- Variance calculation
- Validation rules

### Integration Tests
- API endpoints
- Database operations
- Transaction atomicity
- Error handling

### E2E Tests
- Complete closure flow
- Variance documentation
- Fund handover creation
- Dashboard display

## Migration Strategy

### Database Migration
```javascript
// No migration needed for existing data
// New collection will be created on first use
db.createCollection("fund_handovers")
db.fund_handovers.createIndex({ cashier_shift_id: 1 }, { unique: true })
db.fund_handovers.createIndex({ cashier_id: 1, handover_at: -1 })
db.fund_handovers.createIndex({ handover_at: -1 })
```

### Backward Compatibility
- Existing cashier shifts without fund handover records will work normally
- New closure flow will create fund handover records going forward
- Old shifts can be queried without issues

## Security Considerations

1. **Authentication**: Only authenticated cashiers can access managed funds
2. **Authorization**: Cashiers can only view their own managed funds
3. **Audit Trail**: All fund handovers are logged with user and device info
4. **Data Integrity**: Use transactions to prevent partial updates
5. **Input Validation**: Validate all amounts and variance documentation

## Performance Considerations

1. **Caching**: Cache managed funds summary for dashboard
2. **Indexing**: Proper indexes on fund_handovers collection
3. **Transaction Timeout**: Set reasonable timeout for closure transaction
4. **Query Optimization**: Use projection to fetch only needed fields
5. **Pagination**: Implement pagination for handover history

## Monitoring and Alerts

1. **Metrics to Track**:
   - Number of fund handovers per day
   - Average variance amount
   - Percentage of handovers with variance
   - Closure transaction duration

2. **Alerts**:
   - Large variance amounts (> threshold)
   - Failed closure transactions
   - Missing fund handover records
   - Unusual patterns in variance reasons

## Future Enhancements

1. **Receiver Selection**: Allow specifying manager as receiver
2. **Manager Approval**: Require approval for large variances
3. **Denomination Breakdown**: Help cashiers count by denomination
4. **Photo Evidence**: Allow uploading photos of cash count
5. **Integration**: Connect with accounting system
6. **Analytics**: Dashboard showing variance trends
7. **Notifications**: Alert manager when handover occurs
8. **Multi-Currency**: Support multiple currencies
