package domain

import (
	"fmt"
	
	"cafe-pos/backend/domain/cashier"
	"cafe-pos/backend/domain/order"
)

// StateMachineManager provides centralized access to all state machines
type StateMachineManager struct {
	CashierShiftSM *cashier.ShiftStateMachine
	WaiterShiftSM  *order.ShiftStateMachine
	OrderSM        *order.OrderStateMachine
}

// NewStateMachineManager creates a new centralized state machine manager
func NewStateMachineManager() *StateMachineManager {
	return &StateMachineManager{
		CashierShiftSM: cashier.NewShiftStateMachine(),
		WaiterShiftSM:  order.NewShiftStateMachine(),
		OrderSM:        order.NewOrderStateMachine(),
	}
}

// GetCashierShiftStateMachine returns the cashier shift state machine
func (m *StateMachineManager) GetCashierShiftStateMachine() *cashier.ShiftStateMachine {
	return m.CashierShiftSM
}

// GetWaiterShiftStateMachine returns the waiter/barista shift state machine
func (m *StateMachineManager) GetWaiterShiftStateMachine() *order.ShiftStateMachine {
	return m.WaiterShiftSM
}

// GetOrderStateMachine returns the order state machine
func (m *StateMachineManager) GetOrderStateMachine() *order.OrderStateMachine {
	return m.OrderSM
}

// ValidateCashierShiftTransition validates a cashier shift state transition
func (m *StateMachineManager) ValidateCashierShiftTransition(
	shift *cashier.CashierShift,
	event cashier.ShiftEvent,
) error {
	// Check if transition is valid
	if !m.CashierShiftSM.CanTransition(shift.Status, event) {
		_, err := m.CashierShiftSM.Transition(shift.Status, event)
		return err
	}
	
	// Validate workflow if closing
	if event == cashier.EventCloseShift {
		return m.CashierShiftSM.ValidateShiftWorkflow(shift)
	}
	
	return nil
}

// ValidateWaiterShiftTransition validates a waiter shift state transition
func (m *StateMachineManager) ValidateWaiterShiftTransition(
	shift *order.Shift,
	event order.ShiftEvent,
) error {
	// Check if transition is valid
	if !m.WaiterShiftSM.CanTransition(shift.Status, event) {
		_, err := m.WaiterShiftSM.Transition(shift.Status, event)
		return err
	}
	
	// Validate shift end
	if event == order.EventEndShift {
		return m.WaiterShiftSM.ValidateShiftEnd(shift)
	}
	
	return nil
}

// ValidateOrderTransition validates an order state transition
func (m *StateMachineManager) ValidateOrderTransition(
	ord *order.Order,
	event order.OrderEvent,
) error {
	return m.OrderSM.ValidateTransition(ord, event)
}

// GetCashierShiftNextStep returns the next required step for a cashier shift
func (m *StateMachineManager) GetCashierShiftNextStep(shift *cashier.CashierShift) string {
	return m.CashierShiftSM.GetNextRequiredStep(shift)
}

// GetOrderNextAction returns the next expected action for an order
func (m *StateMachineManager) GetOrderNextAction(ord *order.Order) string {
	return m.OrderSM.GetNextExpectedAction(ord.Status)
}

// CanCancelOrder checks if an order can be cancelled
func (m *StateMachineManager) CanCancelOrder(ord *order.Order) bool {
	return m.OrderSM.CanCancel(ord.Status)
}

// CanRefundOrder checks if an order can be refunded
func (m *StateMachineManager) CanRefundOrder(ord *order.Order) bool {
	return m.OrderSM.CanRefund(ord.Status)
}

// ValidateCashierShiftStep validates a specific step in cashier shift closure
func (m *StateMachineManager) ValidateCashierShiftStep(
	shift *cashier.CashierShift,
	step string,
) error {
	switch step {
	case "record_actual_cash":
		return m.CashierShiftSM.ValidateRecordActualCash(shift)
	case "document_variance":
		return m.CashierShiftSM.ValidateDocumentVariance(shift)
	case "confirm_responsibility":
		return m.CashierShiftSM.ValidateConfirmResponsibility(shift)
	default:
		return fmt.Errorf("unknown step: %s", step)
	}
}

// CanCancelCashierShiftClosure checks if cashier shift closure can be cancelled
func (m *StateMachineManager) CanCancelCashierShiftClosure(shift *cashier.CashierShift) bool {
	return m.CashierShiftSM.CanCancelClosure(shift)
}

// CanModifyOrder checks if an order can be modified
func (m *StateMachineManager) CanModifyOrder(ord *order.Order) bool {
	return m.OrderSM.CanModifyOrder(ord.Status)
}

// CanLockOrder checks if an order can be locked
func (m *StateMachineManager) CanLockOrder(ord *order.Order) bool {
	return m.OrderSM.CanLockOrder(ord.Status)
}

// GetOrderProgress returns the progress percentage of an order
func (m *StateMachineManager) GetOrderProgress(ord *order.Order) int {
	return m.OrderSM.GetOrderProgress(ord.Status)
}

// ValidateWaiterShiftStart validates if a waiter shift can be started
func (m *StateMachineManager) ValidateWaiterShiftStart(existingShift *order.Shift) error {
	return m.WaiterShiftSM.ValidateShiftStart(existingShift)
}

// CanStartWaiterShift checks if a new waiter shift can be started
func (m *StateMachineManager) CanStartWaiterShift(existingShift *order.Shift) bool {
	return m.WaiterShiftSM.CanStartShift(existingShift)
}

// GetWaiterShiftDuration calculates waiter shift duration
func (m *StateMachineManager) GetWaiterShiftDuration(shift *order.Shift) float64 {
	return m.WaiterShiftSM.GetShiftDuration(shift)
}

// IsCashierShiftTerminal checks if cashier shift is in terminal state
func (m *StateMachineManager) IsCashierShiftTerminal(shift *cashier.CashierShift) bool {
	return m.CashierShiftSM.IsTerminalState(shift.Status)
}

// IsOrderTerminal checks if order is in terminal state
func (m *StateMachineManager) IsOrderTerminal(ord *order.Order) bool {
	return m.OrderSM.IsTerminalState(ord.Status)
}

// IsWaiterShiftTerminal checks if waiter shift is in terminal state
func (m *StateMachineManager) IsWaiterShiftTerminal(shift *order.Shift) bool {
	return m.WaiterShiftSM.IsTerminalState(shift.Status)
}
