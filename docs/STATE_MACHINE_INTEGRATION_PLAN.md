# State Machine Integration Plan

## Current Status ✅

### Completed Implementation
1. **State Machines** - Fully implemented
   - `backend/domain/cashier/shift_state_machine.go` - Cashier shift state machine
   - `backend/domain/order/order_state_machine.go` - Order state machine
   - `backend/domain/order/shift_state_machine.go` - Waiter/Barista shift state machine

2. **State Machine Manager** - Centralized access
   - `backend/domain/state_machine_manager.go` - Unified interface for all state machines

3. **API Endpoints** - Public endpoints for state machine info
   - `GET /api/state-machines` - List all state machines
   - `GET /api/state-machines/cashier-shift` - Cashier shift states & transitions
   - `GET /api/state-machines/waiter-shift` - Waiter shift states & transitions
   - `GET /api/state-machines/order` - Order states & transitions

4. **Documentation** - Comprehensive documentation
   - `STATE_MACHINE_DOCUMENTATION.md` - Full documentation with examples

## Missing Integration ❌

### Handlers Need State Machine Validation

#### 1. CashierShiftClosureHandler
**File**: `backend/interfaces/http/cashier_shift_closure_handler.go`

**Current**: Direct domain method calls without state machine validation
**Needed**: Add state machine validation before each transition

**Methods to Update**:
- `InitiateClosure()` - Validate EventInitiateClosure
- `RecordActualCash()` - Validate step "record_actual_cash"
- `DocumentVariance()` - Validate step "document_variance"
- `ConfirmResponsibility()` - Validate step "confirm_responsibility"
- `CloseShift()` - Validate EventCloseShift + full workflow

#### 2. OrderHandler
**File**: `backend/interfaces/http/order_handler.go`

**Current**: Direct service calls without state machine validation
**Needed**: Add state machine validation before each order transition

**Methods to Update**:
- `CollectPayment()` - Validate EventPayOrder
- `SendToBar()` - Validate EventSendToBar
- `AcceptOrder()` - Validate EventStartPreparing
- `FinishPreparing()` - Validate EventMarkReady
- `ServeOrder()` - Validate EventServeOrder
- `CancelOrder()` - Validate EventCancelOrder
- `RefundPartial()` - Validate EventRefundOrder
- `LockOrder()` - Validate EventLockOrder
- `EditOrder()` - Validate CanModifyOrder

#### 3. ShiftHandler
**File**: `backend/interfaces/http/shift_handler.go`

**Current**: Direct service calls without state machine validation
**Needed**: Add state machine validation before shift transitions

**Methods to Update**:
- `StartShift()` - Validate EventStartShift
- `EndShift()` - Validate EventEndShift
- `CloseShift()` - Validate EventCloseShift

## Integration Strategy

### Step 1: Update Handlers to Use State Machine Manager

Each handler needs to:
1. Import state machine manager
2. Initialize state machine manager in constructor
3. Add validation before each transition
4. Return clear error messages when validation fails

### Step 2: Service Layer Integration (Optional but Recommended)

Update service layer to use state machine validation:
- `backend/application/services/cashier_shift_service.go`
- `backend/application/services/order_service.go`
- `backend/application/services/shift_service.go`

### Step 3: Frontend Integration

Create frontend service to fetch state machine information:
- `frontend/src/services/stateMachine.js`
- Use state machine info to:
  - Enable/disable buttons based on valid transitions
  - Show progress indicators
  - Display next expected action
  - Prevent invalid actions

### Step 4: Testing

Add tests for:
- State machine validation in handlers
- Invalid transition attempts
- Complete workflows
- Edge cases

## Benefits of Integration

### 1. Consistency
- All transitions validated through single source of truth
- No invalid states possible
- Business rules enforced automatically

### 2. Better Error Messages
- Clear validation errors
- Users know why action failed
- Suggest next valid action

### 3. UI/UX Improvements
- Disable invalid actions
- Show progress
- Guide users through workflows
- Prevent mistakes

### 4. Maintainability
- Logic centralized in state machines
- Easy to add new states/transitions
- Clear documentation of valid flows

## Implementation Priority

### High Priority (Do Now)
1. ✅ State machines implemented
2. ✅ State machine manager implemented
3. ❌ **Integrate into CashierShiftClosureHandler** - Critical for cashier workflow
4. ❌ **Integrate into OrderHandler** - Critical for order workflow

### Medium Priority (Do Soon)
5. ❌ Integrate into ShiftHandler - Important for shift management
6. ❌ Service layer integration - Better separation of concerns
7. ❌ Frontend state machine service - Better UX

### Low Priority (Nice to Have)
8. ❌ Event history/audit trail - For debugging and compliance
9. ❌ State transition notifications - Real-time updates
10. ❌ Analytics on state durations - Performance insights

## Next Steps

1. **Integrate state machine validation into handlers** (this is the immediate next step)
2. Test the integration with existing workflows
3. Update frontend to use state machine information
4. Add comprehensive tests
5. Document the integration patterns

## Code Examples

### Example: CashierShiftClosureHandler Integration

```go
type CashierShiftClosureHandler struct {
    cashierShiftService *services.CashierShiftService
    stateMachineManager *domain.StateMachineManager  // ADD THIS
}

func NewCashierShiftClosureHandler(
    cashierShiftService *services.CashierShiftService,
    stateMachineManager *domain.StateMachineManager,  // ADD THIS
) *CashierShiftClosureHandler {
    return &CashierShiftClosureHandler{
        cashierShiftService: cashierShiftService,
        stateMachineManager: stateMachineManager,  // ADD THIS
    }
}

func (h *CashierShiftClosureHandler) InitiateClosure(c *gin.Context) {
    // ... get shift ...
    
    // VALIDATE TRANSITION
    err := h.stateMachineManager.ValidateCashierShiftTransition(
        shift,
        cashier.EventInitiateClosure,
    )
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
            "next_step": h.stateMachineManager.GetCashierShiftNextStep(shift),
        })
        return
    }
    
    // ... proceed with domain logic ...
}
```

### Example: OrderHandler Integration

```go
func (h *OrderHandler) SendToBar(c *gin.Context) {
    // ... get order ...
    
    // VALIDATE TRANSITION
    err := h.stateMachineManager.ValidateOrderTransition(
        order,
        order.EventSendToBar,
    )
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
            "next_action": h.stateMachineManager.GetOrderNextAction(order),
            "can_cancel": h.stateMachineManager.CanCancelOrder(order),
        })
        return
    }
    
    // ... proceed with service call ...
}
```

## Conclusion

The state machine infrastructure is **complete and ready to use**. The next critical step is to **integrate state machine validation into the handlers** to enforce business rules and provide better error handling.

This integration will:
- ✅ Prevent invalid state transitions
- ✅ Provide clear error messages
- ✅ Enable better UI/UX
- ✅ Make the system more maintainable
- ✅ Enforce business rules consistently

**Ready to proceed with handler integration!**
