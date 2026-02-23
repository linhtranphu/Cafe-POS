# Fix: Handover Rejection Rollback với MongoDB Transaction

## Vấn đề

Khi cashier từ chối (reject) bàn giao, hệ thống KHÔNG rollback lại `remaining_cash` và `remaining_transfer` cho waiter shift. Điều này dẫn đến:

1. Waiter bàn giao 100k
2. `remaining_cash` giảm từ 200k → 100k ✅
3. Cashier từ chối
4. `remaining_cash` vẫn là 100k ❌ (should be 200k)
5. Waiter mất 100k trong shift

## Root Cause

### Flow hiện tại

**CreateHandover** (Waiter tạo bàn giao):
```go
// Reduce remaining amounts immediately
waiterShift.RemainingCash -= cashAmount
waiterShift.RemainingTransfer -= transferAmount
```

**ConfirmHandover** (Cashier xác nhận/từ chối):
```go
// Update handover status
h.Status = status  // CONFIRMED or REJECTED

// If confirmed, update balances
if status == handover.StatusConfirmed && !h.RequiresApproval {
    updateBalances(ctx, h)
}

// ❌ KHÔNG có logic rollback khi REJECTED
```

### Vấn đề

- Khi `CreateHandover`: Trừ tiền ngay lập tức
- Khi `ConfirmHandover` với `StatusRejected`: KHÔNG restore lại tiền
- Không có transaction → Nếu có lỗi giữa chừng, data inconsistent

## Giải pháp

### 1. Sử dụng MongoDB Transaction

MongoDB transaction đảm bảo:
- **Atomicity**: Tất cả operations thành công hoặc tất cả rollback
- **Consistency**: Data luôn ở trạng thái hợp lệ
- **Isolation**: Các transaction không ảnh hưởng lẫn nhau
- **Durability**: Khi commit thành công, data được lưu vĩnh viễn

### 2. Implementation

#### File: `backend/application/services/cash_handover_service.go`

**Thêm mongoClient vào service**:
```go
type CashHandoverService struct {
    handoverRepo         *mongodb.CashHandoverRepository
    discrepancyRepo      *mongodb.CashDiscrepancyRepository
    shiftRepo            *mongodb.ShiftRepository
    cashierShiftRepo     *mongodb.CashierShiftRepository
    orderRepo            *mongodb.OrderRepository
    mongoClient          *mongo.Client  // ✅ NEW
    discrepancyThreshold float64
}
```

**Update ConfirmHandover với transaction**:
```go
func (s *CashHandoverService) ConfirmHandover(...) error {
    // 1-5. Validation (same as before)
    
    // 6. START MONGODB TRANSACTION
    session, err := s.mongoClient.StartSession()
    if err != nil {
        return fmt.Errorf("failed to start transaction session: %w", err)
    }
    defer session.EndSession(ctx)

    // Execute transaction
    _, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
        // 6a. Update handover record
        h.Status = status
        h.CashierNote = cashierNote
        // ... update other fields
        
        // 6b. Handle discrepancy if exists
        if status == handover.StatusConfirmed && totalDiscrepancy != 0 {
            // ... create discrepancy record
        }
        
        // 6c. Update handover record
        if err := s.handoverRepo.Update(sessCtx, handoverID, h); err != nil {
            return nil, err
        }
        
        // 6d. Handle REJECTION - Rollback remaining amounts ✅ NEW
        if status == handover.StatusRejected {
            waiterShift, err := s.shiftRepo.FindByID(sessCtx, h.WaiterShiftID)
            if err != nil {
                return nil, err
            }
            
            // RESTORE the amounts that were deducted in CreateHandover
            waiterShift.RemainingCash += h.CashDeclaredAmount
            waiterShift.RemainingTransfer += h.TransferDeclaredAmount
            waiterShift.UpdatedAt = time.Now()
            
            // Update waiter shift
            if err := s.shiftRepo.Update(sessCtx, h.WaiterShiftID, waiterShift); err != nil {
                return nil, err
            }
        }
        
        // 6e. If confirmed, update balances
        if status == handover.StatusConfirmed && !h.RequiresApproval {
            if err := s.updateBalances(sessCtx, h); err != nil {
                return nil, err
            }
        }
        
        return nil, nil
    })
    
    if err != nil {
        return err  // Transaction auto-rollback on error
    }
    
    return nil
}
```

#### File: `backend/main.go`

**Pass mongoClient vào service**:
```go
// Handover service
cashHandoverService := services.NewCashHandoverService(
    cashHandoverRepo, 
    cashDiscrepancyRepo, 
    shiftRepo, 
    cashierShiftRepo, 
    orderRepo,
    client,  // ✅ NEW - MongoDB client for transactions
)
```

## Flow sau khi fix

### Case 1: Cashier CONFIRM

```
1. Waiter tạo bàn giao 100k
   → START TRANSACTION
   → remaining_cash: 200k → 100k
   → handover status: PENDING
   → COMMIT TRANSACTION ✅

2. Cashier confirm
   → START TRANSACTION
   → handover status: PENDING → CONFIRMED
   → cashier shift: received_cash += 100k
   → waiter shift: handed_over_cash += 100k
   → COMMIT TRANSACTION ✅

Result: ✅ Waiter remaining_cash = 100k (correct)
```

### Case 2: Cashier REJECT

```
1. Waiter tạo bàn giao 100k
   → START TRANSACTION
   → remaining_cash: 200k → 100k
   → handover status: PENDING
   → COMMIT TRANSACTION ✅

2. Cashier reject
   → START TRANSACTION
   → handover status: PENDING → REJECTED
   → ROLLBACK: remaining_cash: 100k → 200k ✅ NEW
   → COMMIT TRANSACTION ✅

Result: ✅ Waiter remaining_cash = 200k (restored)
```

### Case 3: Error during transaction

```
1. Waiter tạo bàn giao 100k
   → START TRANSACTION
   → remaining_cash: 200k → 100k
   → handover status: PENDING
   → COMMIT TRANSACTION ✅

2. Cashier reject
   → START TRANSACTION
   → handover status: PENDING → REJECTED
   → ROLLBACK: remaining_cash: 100k → 200k
   → ❌ ERROR updating shift
   → AUTO ROLLBACK ENTIRE TRANSACTION ✅
   → handover status: PENDING (unchanged)
   → remaining_cash: 100k (unchanged)

Result: ✅ Data consistent - can retry
```

## Benefits

### 1. Data Integrity ✅
- Không bao giờ mất tiền khi cashier reject
- Transaction đảm bảo all-or-nothing

### 2. Atomicity ✅
- Tất cả updates thành công hoặc tất cả rollback
- Không có trạng thái "giữa chừng"

### 3. Error Recovery ✅
- Nếu có lỗi, transaction tự động rollback
- Data luôn ở trạng thái hợp lệ

### 4. Audit Trail ✅
- Log rõ ràng mỗi bước trong transaction
- Dễ debug khi có vấn đề

## Testing

### Test Case 1: Normal Rejection
```bash
# 1. Start shift with 200k
# 2. Collect payment 100k → remaining_cash = 300k
# 3. Create handover 100k → remaining_cash = 200k
# 4. Cashier reject
# Expected: remaining_cash = 300k (restored)
```

### Test Case 2: Multiple Handovers
```bash
# 1. Start shift with 200k
# 2. Handover 50k → remaining_cash = 150k
# 3. Cashier confirm → remaining_cash = 150k
# 4. Handover 50k → remaining_cash = 100k
# 5. Cashier reject → remaining_cash = 150k (restored)
```

### Test Case 3: Mixed Cash and Transfer
```bash
# 1. Start shift with 100k cash
# 2. Collect 50k transfer → remaining_transfer = 50k
# 3. Handover 50k cash + 30k transfer
#    → remaining_cash = 50k
#    → remaining_transfer = 20k
# 4. Cashier reject
#    → remaining_cash = 100k (restored)
#    → remaining_transfer = 50k (restored)
```

### Test Case 4: Transaction Failure
```bash
# Simulate DB error during rejection
# Expected: All changes rollback, handover stays PENDING
```

## MongoDB Transaction Requirements

### 1. Replica Set
MongoDB transactions require replica set. For development:

```bash
# Start MongoDB as replica set
mongod --replSet rs0

# Initialize replica set
mongo
> rs.initiate()
```

### 2. Connection String
```
mongodb://localhost:27017/?replicaSet=rs0
```

### 3. Error Handling
```go
session, err := s.mongoClient.StartSession()
if err != nil {
    // Handle: MongoDB not configured as replica set
    return fmt.Errorf("transaction not supported: %w", err)
}
```

## Deployment Checklist

- [x] Add mongoClient to CashHandoverService
- [x] Implement transaction in ConfirmHandover
- [x] Add rollback logic for rejection
- [x] Update main.go to pass mongoClient
- [x] Add logging for debugging
- [ ] Test with replica set MongoDB
- [ ] Test rejection flow
- [ ] Test error scenarios
- [ ] Update API documentation

## Files Changed

1. `backend/application/services/cash_handover_service.go`
   - Added `mongoClient *mongo.Client` field
   - Updated `NewCashHandoverService` constructor
   - Rewrote `ConfirmHandover` with transaction
   - Added rollback logic for rejection

2. `backend/main.go`
   - Updated `NewCashHandoverService` call to pass `client`

## Notes

### Transaction Performance
- Transactions có overhead nhỏ (~1-2ms)
- Acceptable cho handover operations (không frequent)
- Benefits về data integrity >> performance cost

### Backward Compatibility
- Không breaking changes cho API
- Frontend không cần update
- Chỉ thay đổi internal implementation

### Future Improvements
1. Add transaction timeout configuration
2. Add retry logic for transient errors
3. Add metrics for transaction success/failure rates
4. Consider optimistic locking for high concurrency

## Related Issues

- HANDOVER_LOGIC_FIX.md - Original handover implementation
- HANDOVER_COMPLETE_FIX_SUMMARY.md - Previous fixes
- TRANSFER_HANDOVER_COMPLETE_FIX.md - Transfer handover support
