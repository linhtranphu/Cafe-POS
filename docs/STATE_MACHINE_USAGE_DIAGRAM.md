# State Machine Usage Diagram

## Kiến Trúc Hiện Tại

```
┌─────────────────────────────────────────────────────────────┐
│                    State Machine Manager                     │
│                 (domain/state_machine_manager.go)            │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  Cashier     │  │   Waiter     │  │    Order     │     │
│  │   Shift      │  │    Shift     │  │              │     │
│  │State Machine │  │State Machine │  │State Machine │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│                                                              │
│  Methods:                                                    │
│  • ValidateCashierShiftTransition()                         │
│  • ValidateWaiterShiftTransition()                          │
│  • ValidateOrderTransition()                                │
│  • GetCashierShiftNextStep()                                │
│  • GetOrderNextAction()                                     │
│  • CanCancelOrder(), CanRefundOrder(), etc.                 │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ Injected via DI
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│   Cashier     │    │     Order     │    │     Shift     │
│    Shift      │    │    Handler    │    │    Handler    │
│   Closure     │    │               │    │               │
│   Handler     │    │               │    │               │
│               │    │               │    │               │
│   ✅ USING    │    │   ❌ NOT      │    │   ❌ NOT      │
│   STATE       │    │   USING       │    │   USING       │
│   MACHINE     │    │   STATE       │    │   STATE       │
│               │    │   MACHINE     │    │   MACHINE     │
└───────────────┘    └───────────────┘    └───────────────┘
```

## Chi Tiết Sử Dụng

### ✅ CashierShiftClosureHandler (Đã Tích Hợp)

```
┌─────────────────────────────────────────────────────────┐
│         CashierShiftClosureHandler                      │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  InitiateClosure()                                      │
│    ├─ stateMachineManager.ValidateCashierShiftTransition()
│    ├─ shift.InitiateClosure()                          │
│    └─ Save to DB                                        │
│                                                         │
│  RecordActualCash()                                     │
│    ├─ stateMachineManager.ValidateCashierShiftStep()   │
│    ├─ shift.RecordActualCash()                         │
│    └─ Save to DB                                        │
│                                                         │
│  DocumentVariance()                                     │
│    ├─ stateMachineManager.ValidateCashierShiftStep()   │
│    ├─ shift.DocumentVariance()                         │
│    └─ Save to DB                                        │
│                                                         │
│  ConfirmResponsibility()                                │
│    ├─ stateMachineManager.ValidateCashierShiftStep()   │
│    ├─ shift.ConfirmResponsibility()                    │
│    └─ Save to DB                                        │
│                                                         │
│  CloseShift()                                           │
│    ├─ stateMachineManager.ValidateCashierShiftTransition()
│    ├─ Check waiter shifts                              │
│    ├─ shift.Close()                                     │
│    └─ Save to DB                                        │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

**Status**: ✅ **HOÀN CHỈNH** - Tất cả 5 methods đều validate qua state machine

### ❌ OrderHandler (Chưa Tích Hợp)

```
┌─────────────────────────────────────────────────────────┐
│              OrderHandler                               │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  CollectPayment()                                       │
│    ├─ ❌ NO STATE VALIDATION                           │
│    ├─ orderService.CollectPayment()                    │
│    └─ Return result                                     │
│                                                         │
│  SendToBar()                                            │
│    ├─ ❌ NO STATE VALIDATION                           │
│    ├─ orderService.SendToBar()                         │
│    └─ Return result                                     │
│                                                         │
│  AcceptOrder()                                          │
│    ├─ ❌ NO STATE VALIDATION                           │
│    ├─ orderService.AcceptOrder()                       │
│    └─ Return result                                     │
│                                                         │
│  ... (6 more methods without validation)               │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

**Status**: ❌ **CHƯA TÍCH HỢP** - 9 methods không validate state

**Risk**: Có thể xảy ra invalid state transitions!

### ❌ ShiftHandler (Chưa Tích Hợp)

```
┌─────────────────────────────────────────────────────────┐
│              ShiftHandler                               │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  StartShift()                                           │
│    ├─ ❌ NO STATE VALIDATION                           │
│    ├─ shiftService.StartShift()                        │
│    └─ Return result                                     │
│                                                         │
│  EndShift()                                             │
│    ├─ ❌ NO STATE VALIDATION                           │
│    ├─ shiftService.EndShift()                          │
│    └─ Return result                                     │
│                                                         │
│  CloseShift()                                           │
│    ├─ ❌ NO STATE VALIDATION                           │
│    ├─ shiftService.CloseShiftAndLockOrders()           │
│    └─ Return result                                     │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

**Status**: ❌ **CHƯA TÍCH HỢP** - 3 methods không validate state

**Risk**: Có thể start shift khi đã có shift open!

## So Sánh: Có vs Không Có State Machine Validation

### Scenario: User Cố Gắng Close Shift Khi Chưa Confirm Responsibility

#### ✅ Với State Machine (CashierShiftClosureHandler)

```
User → CloseShift()
         │
         ├─ stateMachineManager.ValidateCashierShiftTransition()
         │    │
         │    ├─ Check: shift.Confirmation == nil?
         │    └─ ❌ FAIL: "responsibility must be confirmed before closing"
         │
         └─ Return error + next_step: "Confirm responsibility"
```

**Result**: ✅ Prevented invalid transition, clear error message

#### ❌ Không Có State Machine (OrderHandler)

```
User → SendToBar()
         │
         ├─ ❌ NO VALIDATION
         │
         ├─ orderService.SendToBar()
         │    │
         │    └─ May succeed even if order is in wrong state!
         │
         └─ Return success (but state may be invalid)
```

**Result**: ❌ Invalid state possible, no clear error

## Dependency Injection Flow

### main.go

```go
func main() {
    // 1. Create State Machine Manager
    smManager := domain.NewStateMachineManager()
    
    // 2. Inject into handlers
    
    // ✅ CashierShiftClosureHandler - HAS state machine
    cashierShiftClosureHandler := http.NewCashierShiftClosureHandler(
        cashierShiftService,
        smManager,  // ✅ Injected
    )
    
    // ❌ OrderHandler - NO state machine
    orderHandler := http.NewOrderHandler(
        orderService,
        // ❌ Missing: smManager
    )
    
    // ❌ ShiftHandler - NO state machine
    shiftHandler := http.NewShiftHandler(
        shiftService,
        // ❌ Missing: smManager
    )
}
```

## Roadmap: Tích Hợp Đầy Đủ

### Phase 1: ✅ DONE
- ✅ Create State Machine Manager
- ✅ Integrate into CashierShiftClosureHandler
- ✅ Test and verify

### Phase 2: 🔴 TODO (High Priority)
- ⏳ Integrate into OrderHandler
- ⏳ Add validation to 9 methods
- ⏳ Test order workflows

### Phase 3: 🟡 TODO (Medium Priority)
- ⏳ Integrate into ShiftHandler
- ⏳ Add validation to 3 methods
- ⏳ Test shift workflows

### Phase 4: 🟢 TODO (Low Priority)
- ⏳ Move validation to service layer
- ⏳ Better separation of concerns
- ⏳ Add comprehensive tests

## Kết Luận

**Hiện tại**: State Machine Manager đã được tạo và hoạt động tốt, nhưng chỉ 1/3 handlers đang sử dụng.

**Cần làm**: Tích hợp vào OrderHandler và ShiftHandler để đạt được quản lý tập trung hoàn toàn.

**Lợi ích khi hoàn thành**:
- ✅ 100% state transitions được validate
- ✅ Không thể có invalid states
- ✅ Clear error messages cho tất cả handlers
- ✅ Consistent behavior across the system
