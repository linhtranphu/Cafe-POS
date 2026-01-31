package order

import (
	"fmt"
)

// OrderEvent represents events that can trigger order state transitions
type OrderEvent string

const (
	// Order lifecycle events
	EventCreateOrder    OrderEvent = "CREATE_ORDER"
	EventPayOrder       OrderEvent = "PAY_ORDER"
	EventSendToBar      OrderEvent = "SEND_TO_BAR"
	EventStartPreparing OrderEvent = "START_PREPARING"
	EventMarkReady      OrderEvent = "MARK_READY"
	EventServeOrder     OrderEvent = "SERVE_ORDER"
	EventCancelOrder    OrderEvent = "CANCEL_ORDER"
	EventRefundOrder    OrderEvent = "REFUND_ORDER"
	EventLockOrder      OrderEvent = "LOCK_ORDER"
)

// OrderStateMachine manages state transitions for orders
type OrderStateMachine struct {
	// Valid transitions map: current_state -> event -> next_state
	transitions map[OrderStatus]map[OrderEvent]OrderStatus
}

// NewOrderStateMachine creates a new state machine for orders
func NewOrderStateMachine() *OrderStateMachine {
	sm := &OrderStateMachine{
		transitions: make(map[OrderStatus]map[OrderEvent]OrderStatus),
	}
	
	sm.defineTransitions()
	
	return sm
}

// defineTransitions sets up all valid state transitions
func (sm *OrderStateMachine) defineTransitions() {
	// From CREATED state (order created but not paid)
	sm.transitions[StatusCreated] = map[OrderEvent]OrderStatus{
		EventPayOrder:    StatusPaid,
		EventCancelOrder: StatusCancelled,
	}
	
	// From PAID state (paid but not sent to bar)
	sm.transitions[StatusPaid] = map[OrderEvent]OrderStatus{
		EventSendToBar:   StatusQueued,
		EventCancelOrder: StatusCancelled,
		EventRefundOrder: StatusRefunded,
	}
	
	// From QUEUED state (waiting for barista)
	sm.transitions[StatusQueued] = map[OrderEvent]OrderStatus{
		EventStartPreparing: StatusInProgress,
		EventCancelOrder:    StatusCancelled,
	}
	
	// From IN_PROGRESS state (being prepared)
	sm.transitions[StatusInProgress] = map[OrderEvent]OrderStatus{
		EventMarkReady:   StatusReady,
		EventCancelOrder: StatusCancelled,
	}
	
	// From READY state (ready to serve)
	sm.transitions[StatusReady] = map[OrderEvent]OrderStatus{
		EventServeOrder: StatusServed,
	}
	
	// From SERVED state (completed)
	sm.transitions[StatusServed] = map[OrderEvent]OrderStatus{
		EventLockOrder:   StatusLocked,
		EventRefundOrder: StatusRefunded,
	}
	
	// Terminal states (no transitions)
	sm.transitions[StatusCancelled] = map[OrderEvent]OrderStatus{}
	sm.transitions[StatusRefunded] = map[OrderEvent]OrderStatus{}
	sm.transitions[StatusLocked] = map[OrderEvent]OrderStatus{}
}

// CanTransition checks if a transition is valid
func (sm *OrderStateMachine) CanTransition(currentState OrderStatus, event OrderEvent) bool {
	if stateTransitions, exists := sm.transitions[currentState]; exists {
		_, canTransition := stateTransitions[event]
		return canTransition
	}
	return false
}

// Transition attempts to transition from current state to next state via event
func (sm *OrderStateMachine) Transition(currentState OrderStatus, event OrderEvent) (OrderStatus, error) {
	if !sm.CanTransition(currentState, event) {
		return currentState, fmt.Errorf(
			"invalid transition: cannot apply event '%s' in state '%s'",
			event,
			currentState,
		)
	}
	
	nextState := sm.transitions[currentState][event]
	return nextState, nil
}

// GetValidEvents returns all valid events for the current state
func (sm *OrderStateMachine) GetValidEvents(currentState OrderStatus) []OrderEvent {
	var events []OrderEvent
	
	if stateTransitions, exists := sm.transitions[currentState]; exists {
		for event := range stateTransitions {
			events = append(events, event)
		}
	}
	
	return events
}

// IsTerminalState checks if a state is terminal (no further transitions)
func (sm *OrderStateMachine) IsTerminalState(state OrderStatus) bool {
	transitions, exists := sm.transitions[state]
	return !exists || len(transitions) == 0
}

// CanCancel checks if an order can be cancelled from its current state
func (sm *OrderStateMachine) CanCancel(currentState OrderStatus) bool {
	return sm.CanTransition(currentState, EventCancelOrder)
}

// CanRefund checks if an order can be refunded from its current state
func (sm *OrderStateMachine) CanRefund(currentState OrderStatus) bool {
	return sm.CanTransition(currentState, EventRefundOrder)
}

// GetNextExpectedAction returns the next expected action for an order
func (sm *OrderStateMachine) GetNextExpectedAction(currentState OrderStatus) string {
	switch currentState {
	case StatusCreated:
		return "Payment required"
	case StatusPaid:
		return "Send to bar"
	case StatusQueued:
		return "Start preparing"
	case StatusInProgress:
		return "Mark as ready"
	case StatusReady:
		return "Serve to customer"
	case StatusServed:
		return "Order completed"
	case StatusCancelled:
		return "Order cancelled"
	case StatusRefunded:
		return "Order refunded"
	case StatusLocked:
		return "Order locked"
	default:
		return "Unknown state"
	}
}

// ValidateTransition validates if a transition is allowed with additional business rules
func (sm *OrderStateMachine) ValidateTransition(order *Order, event OrderEvent) error {
	// First check if the transition is valid in the state machine
	if !sm.CanTransition(order.Status, event) {
		return fmt.Errorf(
			"invalid transition: cannot apply event '%s' in state '%s'",
			event,
			order.Status,
		)
	}
	
	// Additional business rule validations
	switch event {
	case EventPayOrder:
		if order.Total <= 0 {
			return fmt.Errorf("cannot pay order with zero or negative amount")
		}
		
	case EventSendToBar:
		if len(order.Items) == 0 {
			return fmt.Errorf("cannot send empty order to bar")
		}
		
	case EventRefundOrder:
		if order.PaymentMethod == "" {
			return fmt.Errorf("cannot refund order without payment method")
		}
		
	case EventCancelOrder:
		// Check if order can still be cancelled based on status
		if order.Status == StatusServed || order.Status == StatusLocked {
			return fmt.Errorf("cannot cancel order in %s status", order.Status)
		}
	}
	
	return nil
}

// CanModifyOrder checks if an order can be modified
func (sm *OrderStateMachine) CanModifyOrder(currentState OrderStatus) bool {
	// Can only modify orders that haven't been sent to bar yet
	return currentState == StatusCreated || currentState == StatusPaid
}

// CanLockOrder checks if an order can be locked
func (sm *OrderStateMachine) CanLockOrder(currentState OrderStatus) bool {
	return currentState == StatusServed
}

// GetOrderProgress returns the progress percentage of an order
func (sm *OrderStateMachine) GetOrderProgress(currentState OrderStatus) int {
	switch currentState {
	case StatusCreated:
		return 0
	case StatusPaid:
		return 20
	case StatusQueued:
		return 40
	case StatusInProgress:
		return 60
	case StatusReady:
		return 80
	case StatusServed:
		return 100
	case StatusLocked:
		return 100
	case StatusCancelled, StatusRefunded:
		return 0
	default:
		return 0
	}
}
