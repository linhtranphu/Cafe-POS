# State Machine Documentation

## Tổng quan

Hệ thống sử dụng State Machine pattern để quản lý trạng thái và luồng nghiệp vụ của:
1. **Cashier Shifts** - Ca làm việc thu ngân
2. **Waiter/Barista Shifts** - Ca làm việc phục vụ/pha chế
3. **Orders** - Đơn hàng

## 1. Cashier Shift State Machine

### States (Trạng thái)
```
OPEN → CLOSURE_INITIATED → CLOSED
```

- **OPEN**: Ca đang mở, thu ngân đang làm việc
- **CLOSURE_INITIATED**: Đã bắt đầu quy trình đóng ca
- **CLOSED**: Ca đã đóng (terminal state)

### Events (Sự kiện)
- `INITIATE_CLOSURE`: Bắt đầu đóng ca
- `RECORD_ACTUAL_CASH`: Ghi nhận tiền mặt thực tế
- `DOCUMENT_VARIANCE`: Giải trình chênh lệch
- `CONFIRM_RESPONSIBILITY`: Xác nhận trách nhiệm
- `CLOSE_SHIFT`: Đóng ca
- `CANCEL_CLOSURE`: Hủy quy trình đóng ca

### Transitions (Chuyển đổi)
```
OPEN + INITIATE_CLOSURE → CLOSURE_INITIATED
CLOSURE_INITIATED + CLOSE_SHIFT → CLOSED
CLOSURE_INITIATED + CANCEL_CLOSURE → OPEN (chỉ khi chưa record actual cash)
```

### Business Rules
1. Chỉ có thể initiate closure khi ở trạng thái OPEN
2. Phải record actual cash trước khi document variance
3. Nếu có variance (≠ 0), phải document variance trước khi confirm responsibility
4. Phải confirm responsibility trước khi close shift
5. Chỉ có thể cancel closure nếu chưa record actual cash
6. Không thể đóng ca cashier khi còn ca waiter đang mở

### Validation Methods
- `ValidateShiftWorkflow()`: Validate toàn bộ workflow
- `ValidateRecordActualCash()`: Validate bước ghi tiền
- `ValidateDocumentVariance()`: Validate bước giải trình
- `ValidateConfirmResponsibility()`: Validate bước xác nhận
- `CanCancelClosure()`: Kiểm tra có thể hủy không

---

## 2. Order State Machine

### States (Trạng thái)
```
CREATED → PAID → QUEUED → IN_PROGRESS → READY → SERVED → LOCKED
                    ↓          ↓           ↓        ↓
                CANCELLED  CANCELLED  CANCELLED  REFUNDED
```

- **CREATED**: Order mới tạo, chưa thanh toán
- **PAID**: Đã thanh toán, chưa gửi cho barista
- **QUEUED**: Đã gửi cho barista, chờ nhận
- **IN_PROGRESS**: Barista đang pha chế
- **READY**: Đã pha xong, chờ giao
- **SERVED**: Đã giao cho khách
- **CANCELLED**: Đã hủy (terminal state)
- **REFUNDED**: Đã hoàn tiền (terminal state)
- **LOCKED**: Đã chốt ca (terminal state)

### Events (Sự kiện)
- `CREATE_ORDER`: Tạo order
- `PAY_ORDER`: Thanh toán
- `SEND_TO_BAR`: Gửi cho barista
- `START_PREPARING`: Bắt đầu pha chế
- `MARK_READY`: Đánh dấu sẵn sàng
- `SERVE_ORDER`: Giao cho khách
- `CANCEL_ORDER`: Hủy order
- `REFUND_ORDER`: Hoàn tiền
- `LOCK_ORDER`: Khóa order (khi chốt ca)

### Transitions (Chuyển đổi)
```
CREATED + PAY_ORDER → PAID
CREATED + CANCEL_ORDER → CANCELLED

PAID + SEND_TO_BAR → QUEUED
PAID + CANCEL_ORDER → CANCELLED
PAID + REFUND_ORDER → REFUNDED

QUEUED + START_PREPARING → IN_PROGRESS
QUEUED + CANCEL_ORDER → CANCELLED

IN_PROGRESS + MARK_READY → READY
IN_PROGRESS + CANCEL_ORDER → CANCELLED

READY + SERVE_ORDER → SERVED

SERVED + LOCK_ORDER → LOCKED
SERVED + REFUND_ORDER → REFUNDED
```

### Business Rules
1. Không thể pay order với total ≤ 0
2. Không thể send to bar nếu order rỗng (no items)
3. Không thể refund nếu không có payment method
4. Không thể cancel order đã SERVED hoặc LOCKED
5. Chỉ có thể modify order ở trạng thái CREATED hoặc PAID
6. Chỉ có thể lock order ở trạng thái SERVED

### Helper Methods
- `CanCancel()`: Kiểm tra có thể hủy
- `CanRefund()`: Kiểm tra có thể hoàn tiền
- `CanModifyOrder()`: Kiểm tra có thể sửa
- `CanLockOrder()`: Kiểm tra có thể khóa
- `GetOrderProgress()`: Lấy % tiến độ (0-100)

---

## 3. Waiter/Barista Shift State Machine

### States (Trạng thái)
```
OPEN → CLOSED
```

- **OPEN**: Ca đang mở
- **CLOSED**: Ca đã đóng (terminal state)

### Events (Sự kiện)
- `START_SHIFT`: Bắt đầu ca
- `END_SHIFT`: Kết thúc ca
- `CLOSE_SHIFT`: Đóng ca

### Transitions (Chuyển đổi)
```
OPEN + END_SHIFT → CLOSED
```

### Business Rules
1. Chỉ có thể end shift khi ở trạng thái OPEN
2. Không thể start shift mới nếu đã có shift OPEN
3. Một user chỉ có thể có 1 shift OPEN tại một thời điểm

### Helper Methods
- `ValidateShiftStart()`: Validate khi start shift
- `ValidateShiftEnd()`: Validate khi end shift
- `CanStartShift()`: Kiểm tra có thể start
- `GetShiftDuration()`: Tính thời gian làm việc (hours)

---

## 4. State Machine Manager

### Centralized Management
`StateMachineManager` cung cấp interface thống nhất để:
- Truy cập tất cả state machines
- Validate transitions
- Kiểm tra business rules
- Helper methods cho UI

### Key Methods

#### Cashier Shift
```go
ValidateCashierShiftTransition(shift, event) error
ValidateCashierShiftStep(shift, step) error
CanCancelCashierShiftClosure(shift) bool
GetCashierShiftNextStep(shift) string
IsCashierShiftTerminal(shift) bool
```

#### Order
```go
ValidateOrderTransition(order, event) error
CanCancelOrder(order) bool
CanRefundOrder(order) bool
CanModifyOrder(order) bool
CanLockOrder(order) bool
GetOrderProgress(order) int
GetOrderNextAction(order) string
IsOrderTerminal(order) bool
```

#### Waiter Shift
```go
ValidateWaiterShiftTransition(shift, event) error
ValidateWaiterShiftStart(existingShift) error
CanStartWaiterShift(existingShift) bool
GetWaiterShiftDuration(shift) float64
IsWaiterShiftTerminal(shift) bool
```

---

## 5. API Endpoints

### State Machine Information (Public)
```
GET /api/state-machines
GET /api/state-machines/cashier-shift
GET /api/state-machines/waiter-shift
GET /api/state-machines/order
```

Response format:
```json
{
  "states": ["OPEN", "CLOSURE_INITIATED", "CLOSED"],
  "events": ["INITIATE_CLOSURE", "CLOSE_SHIFT", ...],
  "transitions": {
    "OPEN": ["INITIATE_CLOSURE"],
    "CLOSURE_INITIATED": ["CLOSE_SHIFT", "CANCEL_CLOSURE"],
    "CLOSED": []
  }
}
```

---

## 6. Usage Examples

### Example 1: Validate Cashier Shift Closure
```go
smManager := domain.NewStateMachineManager()

// Check if can close
err := smManager.ValidateCashierShiftTransition(shift, cashier.EventCloseShift)
if err != nil {
    // Cannot close, show error
    return err
}

// Get next required step
nextStep := smManager.GetCashierShiftNextStep(shift)
// Returns: "Record actual cash" or "Document variance" etc.
```

### Example 2: Validate Order Transition
```go
// Check if can send to bar
err := smManager.ValidateOrderTransition(order, order.EventSendToBar)
if err != nil {
    return err
}

// Check if can cancel
canCancel := smManager.CanCancelOrder(order)

// Get progress
progress := smManager.GetOrderProgress(order) // 0-100
```

### Example 3: Check Waiter Shift
```go
// Check if can start new shift
canStart := smManager.CanStartWaiterShift(existingShift)

// Validate before starting
err := smManager.ValidateWaiterShiftStart(existingShift)
```

---

## 7. Benefits

### 1. **Tính nhất quán**
- Tất cả transitions đều được validate
- Không thể có invalid state
- Business rules được enforce tự động

### 2. **Dễ maintain**
- Logic tập trung ở một nơi
- Dễ thêm states/events mới
- Dễ test

### 3. **Clear documentation**
- State diagram rõ ràng
- Valid transitions được định nghĩa rõ
- Business rules được document

### 4. **Error prevention**
- Validate trước khi transition
- Không thể skip steps
- Audit trail đầy đủ

### 5. **UI/UX tốt hơn**
- Biết được next step
- Disable buttons không hợp lệ
- Progress tracking

---

## 8. Testing

### Unit Tests
Mỗi state machine nên có tests cho:
- Valid transitions
- Invalid transitions
- Business rule validation
- Edge cases

### Integration Tests
- End-to-end workflows
- Multi-step processes
- Error handling

---

## 9. Future Enhancements

### Planned Features
1. **Event History**: Log tất cả transitions
2. **Rollback**: Undo transitions (nếu cần)
3. **Notifications**: Trigger events khi state change
4. **Analytics**: Track state durations, bottlenecks
5. **Visualization**: Generate state diagrams tự động

### Extensibility
- Dễ thêm states mới
- Dễ thêm business rules
- Dễ integrate với external systems

---

## 10. Best Practices

### DO ✅
- Luôn validate trước khi transition
- Sử dụng StateMachineManager thay vì direct access
- Log tất cả state changes
- Handle errors gracefully
- Document business rules

### DON'T ❌
- Không skip validation
- Không modify state trực tiếp
- Không hardcode state strings
- Không assume current state
- Không ignore errors

---

## Conclusion

State Machine pattern giúp:
- ✅ Quản lý trạng thái chặt chẽ
- ✅ Enforce business rules
- ✅ Prevent invalid states
- ✅ Clear documentation
- ✅ Easy to maintain and extend

Backend đã implement đầy đủ và sẵn sàng sử dụng!
