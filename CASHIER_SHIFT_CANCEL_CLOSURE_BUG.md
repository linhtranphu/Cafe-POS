# Bug: Ca thu ngân tự động đóng khi bấm "Quay lại" sau khi bắt đầu đóng ca

## Mô tả vấn đề

Khi người dùng:
1. Vào trang đóng ca thu ngân: `http://localhost:5173/#/cashier/shift-closure/699c269a1ba7fb0030879c87`
2. Bấm "Bắt đầu đóng ca"
3. Bấm "Quay lại" (không tiếp tục thực hiện do 1 ca waiter còn đang mở)

Kết quả: Ca thu ngân tự động chuyển sang trạng thái `CLOSURE_INITIATED` và không thể hoàn tác.

## Nguyên nhân

### 1. Frontend Logic (CashierShiftClosure.vue)

Khi bấm "Bắt đầu đóng ca":
```javascript
const initiateClosure = async () => {
  processing.value = true
  error.value = null
  
  try {
    await cashierShiftService.initiateClosure(shift.value.id)
    await loadShift()  // Reload shift data - status now is CLOSURE_INITIATED
  } catch (err) {
    error.value = err.response?.data?.error || 'Không thể bắt đầu đóng ca'
  } finally {
    processing.value = false
  }
}
```

Khi bấm "Quay lại":
```javascript
const goBack = () => {
  router.push('/cashier')  // Chỉ navigate, KHÔNG hoàn tác trạng thái
}
```

### 2. Backend State Machine

Backend đã có logic để cancel closure:

**State Machine (cashier_shift_state_machine.go)**:
```go
sm.transitions[CashierShiftClosureInitiated] = map[ShiftEvent]CashierShiftStatus{
    EventCloseShift:    CashierShiftClosed,
    EventCancelClosure: CashierShiftOpen, // ✅ Có logic để cancel
}

func (sm *ShiftStateMachine) CanCancelClosure(shift *CashierShift) bool {
    // Can only cancel if in CLOSURE_INITIATED state and no critical steps completed
    if shift.Status != CashierShiftClosureInitiated {
        return false
    }
    // Cannot cancel if actual cash has been recorded
    return shift.ActualCash == nil
}
```

**Nhưng KHÔNG có API endpoint và handler method** để thực hiện cancel closure!

File `cashier_shift_closure_handler.go` có các methods:
- ✅ InitiateClosure
- ✅ RecordActualCash
- ✅ DocumentVariance
- ✅ ConfirmResponsibility
- ✅ CloseShift
- ❌ CancelClosure (THIẾU!)

## Giải pháp

### Option 1: Thêm API Cancel Closure (Khuyến nghị)

#### Backend

1. **Thêm method CancelClosure vào handler**:

```go
// File: backend/interfaces/http/cashier_shift_closure_handler.go

// CancelClosure cancels the shift closure process and returns to OPEN status
func (h *CashierShiftClosureHandler) CancelClosure(c *gin.Context) {
	shiftID := c.Param("id")
	
	shiftObjID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift ID"})
		return
	}

	// Get user_id from JWT token
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system"
	}
	deviceID := "web"

	// Get the shift
	shift, err := h.cashierShiftService.GetCashierShift(c.Request.Context(), shiftObjID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shift not found"})
		return
	}

	// Validate state transition using state machine
	err = h.stateMachineManager.ValidateCashierShiftTransition(shift, cashier.EventCancelClosure)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     err.Error(),
			"next_step": h.stateMachineManager.GetCashierShiftNextStep(shift),
		})
		return
	}

	// Check if can cancel (no critical steps completed)
	if !h.stateMachineManager.CanCancelCashierShiftClosure(shift) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot cancel closure: critical steps have been completed",
		})
		return
	}

	// Cancel closure
	err = shift.CancelClosure(userID, deviceID, time.Now())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the shift
	err = h.cashierShiftService.SaveCashierShift(c.Request.Context(), shift)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save shift"})
		return
	}

	c.JSON(http.StatusOK, shift)
}
```

2. **Thêm method CancelClosure vào domain model**:

```go
// File: backend/domain/cashier/cashier_shift.go

// CancelClosure cancels the shift closure process and returns to Open state
func (cs *CashierShift) CancelClosure(userID, deviceID string, timestamp time.Time) error {
	// Validate current status is ClosureInitiated
	if cs.Status != CashierShiftClosureInitiated {
		return fmt.Errorf("can only cancel closure from CLOSURE_INITIATED status, current: %s", cs.Status)
	}

	// Cannot cancel if actual cash has been recorded
	if cs.ActualCash != nil {
		return fmt.Errorf("cannot cancel closure: actual cash has been recorded")
	}

	// Reset to Open status
	cs.Status = CashierShiftOpen

	// Add audit log
	auditEntry := AuditLogEntry{
		Action:    "CANCEL_CLOSURE",
		UserID:    userID,
		DeviceID:  deviceID,
		Timestamp: timestamp,
		Details:   "Closure process cancelled, returned to OPEN status",
	}
	cs.AuditLog = append(cs.AuditLog, auditEntry)

	return nil
}
```

3. **Thêm route vào main.go**:

```go
// File: backend/main.go

cashierShifts.POST("/:id/initiate-closure", cashierShiftClosureHandler.InitiateClosure)
cashierShifts.POST("/:id/cancel-closure", cashierShiftClosureHandler.CancelClosure) // ← THÊM
cashierShifts.POST("/:id/record-actual-cash", cashierShiftClosureHandler.RecordActualCash)
```

#### Frontend

1. **Thêm method vào service**:

```javascript
// File: frontend/src/services/cashierShift.js

async cancelClosure(shiftId) {
  const response = await api.post(`/cashier-shifts/${shiftId}/cancel-closure`)
  return response.data
}
```

2. **Cập nhật logic goBack trong CashierShiftClosure.vue**:

```javascript
const goBack = async () => {
  // If closure has been initiated but no critical steps completed, offer to cancel
  if (shift.value?.status === CASHIER_SHIFT_STATUS.CLOSURE_INITIATED && 
      !shift.value.actual_cash) {
    
    const shouldCancel = confirm(
      'Bạn đã bắt đầu đóng ca nhưng chưa hoàn thành. ' +
      'Bạn có muốn hủy quy trình đóng ca và quay về trạng thái mở ca không?'
    )
    
    if (shouldCancel) {
      try {
        processing.value = true
        await cashierShiftService.cancelClosure(shift.value.id)
        // Success - go back
        router.push('/cashier')
      } catch (err) {
        error.value = err.response?.data?.error || 'Không thể hủy đóng ca'
        processing.value = false
      }
      return
    }
  }
  
  // Normal back navigation
  router.push('/cashier')
}
```

### Option 2: Cảnh báo người dùng (Giải pháp tạm thời)

Nếu không muốn implement cancel closure ngay, có thể thêm cảnh báo:

```javascript
const goBack = () => {
  if (shift.value?.status === CASHIER_SHIFT_STATUS.CLOSURE_INITIATED) {
    alert(
      '⚠️ Cảnh báo: Ca đang trong quá trình đóng!\n\n' +
      'Bạn đã bắt đầu đóng ca. Vui lòng hoàn thành tất cả các bước ' +
      'hoặc liên hệ quản lý để hủy quy trình đóng ca.'
    )
  }
  router.push('/cashier')
}
```

## Khuyến nghị

**Nên implement Option 1** vì:
1. Backend đã có state machine logic sẵn
2. UX tốt hơn - cho phép người dùng sửa lỗi
3. Tránh trường hợp ca bị "kẹt" ở trạng thái CLOSURE_INITIATED

## Test Cases

Sau khi fix, cần test:

1. ✅ Bắt đầu đóng ca → Bấm quay lại → Confirm cancel → Ca quay về OPEN
2. ✅ Bắt đầu đóng ca → Nhập tiền thực tế → Bấm quay lại → KHÔNG cho phép cancel
3. ✅ Bắt đầu đóng ca → Bấm quay lại → Cancel → Không confirm → Vẫn ở trang đóng ca
4. ✅ Audit log ghi nhận đúng action CANCEL_CLOSURE
