# OrderHandler State Machine Integration - Complete ✅

## Tổng Quan

Đã tích hợp thành công State Machine validation vào OrderHandler. Tất cả 9 methods quan trọng giờ đây đều validate state transitions trước khi thực hiện actions.

## Changes Made

### 1. Updated OrderHandler Struct

**Before**:
```go
type OrderHandler struct {
    orderService *services.OrderService
}
```

**After**:
```go
type OrderHandler struct {
    orderService        *services.OrderService
    stateMachineManager *domain.StateMachineManager  // ✅ Added
}
```

### 2. Updated Constructor

**Before**:
```go
func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
    return &OrderHandler{orderService: orderService}
}
```

**After**:
```go
func NewOrderHandler(
    orderService *services.OrderService,
    stateMachineManager *domain.StateMachineManager,  // ✅ Added parameter
) *OrderHandler {
    return &OrderHandler{
        orderService:        orderService,
        stateMachineManager: stateMachineManager,
    }
}
```

### 3. Updated main.go

**Before**:
```go
orderHandler := http.NewOrderHandler(orderService)
```

**After**:
```go
orderHandler := http.NewOrderHandler(orderService, smManager)  // ✅ Inject smManager
```

## Methods Updated with State Machine Validation

### ✅ 1. CollectPayment()

**Validates**: `EventPayOrder`

**Flow**:
1. Get order from database
2. Validate: Can transition from current state with EventPayOrder?
3. If valid → Proceed with payment
4. If invalid → Return error with next_action and can_cancel

**Error Response**:
```json
{
  "error": "invalid transition: cannot apply event 'PAY_ORDER' in state 'PAID'",
  "next_action": "Send to bar",
  "can_cancel": true
}
```

### ✅ 2. EditOrder()

**Validates**: `CanModifyOrder()`

**Flow**:
1. Get order from database
2. Check: Can modify order in current state?
3. If yes → Proceed with edit
4. If no → Return error with status and next_action

**Error Response**:
```json
{
  "error": "cannot modify order in current state",
  "status": "QUEUED",
  "next_action": "Start preparing"
}
```

### ✅ 3. RefundPartial()

**Validates**: `EventRefundOrder`

**Flow**:
1. Get order from database
2. Validate: Can transition with EventRefundOrder?
3. If valid → Proceed with refund
4. If invalid → Return error with can_refund status

**Error Response**:
```json
{
  "error": "invalid transition: cannot apply event 'REFUND_ORDER' in state 'CREATED'",
  "status": "CREATED",
  "can_refund": false,
  "next_action": "Payment required"
}
```

### ✅ 4. SendToBar()

**Validates**: `EventSendToBar`

**Flow**:
1. Get order from database
2. Validate: Can transition with EventSendToBar?
3. Check business rules (order not empty)
4. If valid → Send to bar
5. If invalid → Return error with next_action

**Error Response**:
```json
{
  "error": "cannot send empty order to bar",
  "status": "PAID",
  "next_action": "Send to bar"
}
```

### ✅ 5. AcceptOrder()

**Validates**: `EventStartPreparing`

**Flow**:
1. Get order from database
2. Validate: Can transition with EventStartPreparing?
3. If valid → Barista accepts order
4. If invalid → Return error with next_action

**Error Response**:
```json
{
  "error": "invalid transition: cannot apply event 'START_PREPARING' in state 'PAID'",
  "status": "PAID",
  "next_action": "Send to bar"
}
```

### ✅ 6. FinishPreparing()

**Validates**: `EventMarkReady`

**Flow**:
1. Get order from database
2. Validate: Can transition with EventMarkReady?
3. If valid → Mark as ready
4. If invalid → Return error with progress and next_action

**Error Response**:
```json
{
  "error": "invalid transition: cannot apply event 'MARK_READY' in state 'QUEUED'",
  "status": "QUEUED",
  "progress": 40,
  "next_action": "Start preparing"
}
```

### ✅ 7. ServeOrder()

**Validates**: `EventServeOrder`

**Flow**:
1. Get order from database
2. Validate: Can transition with EventServeOrder?
3. If valid → Serve to customer
4. If invalid → Return error with next_action

**Error Response**:
```json
{
  "error": "invalid transition: cannot apply event 'SERVE_ORDER' in state 'IN_PROGRESS'",
  "status": "IN_PROGRESS",
  "next_action": "Mark as ready"
}
```

### ✅ 8. CancelOrder()

**Validates**: `EventCancelOrder`

**Flow**:
1. Get order from database
2. Validate: Can transition with EventCancelOrder?
3. Check: Order not in SERVED or LOCKED state
4. If valid → Cancel order
5. If invalid → Return error with can_cancel status

**Error Response**:
```json
{
  "error": "cannot cancel order in SERVED status",
  "status": "SERVED",
  "can_cancel": false,
  "next_action": "Order completed"
}
```

### ✅ 9. LockOrder()

**Validates**: `EventLockOrder`

**Flow**:
1. Get order from database
2. Validate: Can transition with EventLockOrder?
3. Check: Order must be in SERVED state
4. If valid → Lock order
5. If invalid → Return error with can_lock status

**Error Response**:
```json
{
  "error": "invalid transition: cannot apply event 'LOCK_ORDER' in state 'READY'",
  "status": "READY",
  "can_lock": false,
  "next_action": "Serve to customer"
}
```

## State Machine Validation Matrix

| Method | Event | Valid From States | Business Rules |
|--------|-------|-------------------|----------------|
| CollectPayment | PAY_ORDER | CREATED | Total > 0 |
| EditOrder | - | CREATED, PAID | CanModifyOrder() |
| RefundPartial | REFUND_ORDER | PAID, SERVED | Has payment method |
| SendToBar | SEND_TO_BAR | PAID | Order not empty |
| AcceptOrder | START_PREPARING | QUEUED | - |
| FinishPreparing | MARK_READY | IN_PROGRESS | - |
| ServeOrder | SERVE_ORDER | READY | - |
| CancelOrder | CANCEL_ORDER | CREATED, PAID, QUEUED, IN_PROGRESS | Not SERVED/LOCKED |
| LockOrder | LOCK_ORDER | SERVED | - |

## Benefits Achieved

### 1. ✅ Prevent Invalid State Transitions

**Before**: Could potentially send unpaid order to bar
```
Order (CREATED) → SendToBar() → ❌ May succeed incorrectly
```

**After**: State machine prevents invalid transition
```
Order (CREATED) → SendToBar() → ✅ Blocked with clear error
Error: "invalid transition: cannot apply event 'SEND_TO_BAR' in state 'CREATED'"
Next action: "Payment required"
```

### 2. ✅ Clear Error Messages

**Before**:
```json
{
  "error": "order not paid"
}
```

**After**:
```json
{
  "error": "invalid transition: cannot apply event 'SEND_TO_BAR' in state 'CREATED'",
  "status": "CREATED",
  "next_action": "Payment required"
}
```

### 3. ✅ Consistent Validation

All order transitions now go through the same validation logic:
- Centralized in State Machine
- No duplicate validation code
- Easy to maintain and extend

### 4. ✅ Better UX

Frontend can use the additional information:
- `next_action` - Show what user should do next
- `can_cancel` - Enable/disable cancel button
- `can_refund` - Enable/disable refund button
- `progress` - Show order progress (0-100%)

## Testing Scenarios

### Scenario 1: Try to Send Unpaid Order to Bar

```bash
# Create order
POST /api/waiter/orders
{
  "table_number": 5,
  "items": [...]
}
# Response: Order created with status CREATED

# Try to send to bar WITHOUT payment
POST /api/waiter/orders/{id}/send
# ✅ BLOCKED by state machine
# Response:
{
  "error": "invalid transition: cannot apply event 'SEND_TO_BAR' in state 'CREATED'",
  "status": "CREATED",
  "next_action": "Payment required"
}
```

### Scenario 2: Try to Cancel Served Order

```bash
# Order is already served
GET /api/waiter/orders/{id}
# Response: { "status": "SERVED", ... }

# Try to cancel
POST /api/cashier/orders/{id}/cancel
# ✅ BLOCKED by state machine
# Response:
{
  "error": "cannot cancel order in SERVED status",
  "status": "SERVED",
  "can_cancel": false,
  "next_action": "Order completed"
}
```

### Scenario 3: Try to Edit Order After Sent to Bar

```bash
# Order is in queue
GET /api/waiter/orders/{id}
# Response: { "status": "QUEUED", ... }

# Try to edit
PUT /api/waiter/orders/{id}/edit
# ✅ BLOCKED by state machine
# Response:
{
  "error": "cannot modify order in current state",
  "status": "QUEUED",
  "next_action": "Start preparing"
}
```

## Compilation Status

✅ **Backend compiled successfully**
```bash
cd backend && go build -o cafe-pos-server
# Exit Code: 0
```

✅ **No diagnostics errors**
```
backend/interfaces/http/order_handler.go: No diagnostics found
backend/main.go: No diagnostics found
```

## Files Modified

- ✅ `backend/interfaces/http/order_handler.go` - Added state machine validation to 9 methods
- ✅ `backend/main.go` - Updated OrderHandler initialization to inject smManager

## Summary

**Status**: ✅ **COMPLETE**

**Methods Updated**: 9/9 (100%)
- ✅ CollectPayment
- ✅ EditOrder
- ✅ RefundPartial
- ✅ SendToBar
- ✅ AcceptOrder
- ✅ FinishPreparing
- ✅ ServeOrder
- ✅ CancelOrder
- ✅ LockOrder

**Benefits**:
- ✅ All order transitions validated
- ✅ Clear error messages with guidance
- ✅ Consistent validation logic
- ✅ Better UX with next_action hints
- ✅ Prevent invalid state transitions

**Next Steps**:
- 🟡 Integrate state machine into ShiftHandler (3 methods)
- 🟢 Add comprehensive tests
- 🟢 Frontend integration to use next_action hints

OrderHandler is now fully integrated with State Machine! 🎉
