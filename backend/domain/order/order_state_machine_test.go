package order

import (
	"testing"
)

func TestOrderStateMachine_ValidTransitions(t *testing.T) {
	sm := NewOrderStateMachine()

	tests := []struct {
		name      string
		fromState OrderStatus
		event     OrderEvent
		wantState OrderStatus
		wantErr   bool
	}{
		// Valid transitions
		{"Pay order", StatusCreated, EventPayOrder, StatusPaid, false},
		{"Send to bar", StatusPaid, EventSendToBar, StatusQueued, false},
		{"Start preparing", StatusQueued, EventStartPreparing, StatusInProgress, false},
		{"Mark ready", StatusInProgress, EventMarkReady, StatusReady, false},
		{"Serve order", StatusReady, EventServeOrder, StatusServed, false},
		{"Lock order", StatusServed, EventLockOrder, StatusLocked, false},
		{"Cancel from created", StatusCreated, EventCancelOrder, StatusCancelled, false},
		{"Cancel from paid", StatusPaid, EventCancelOrder, StatusCancelled, false},
		{"Refund from paid", StatusPaid, EventRefundOrder, StatusRefunded, false},
		
		// Invalid transitions
		{"Cannot send unpaid", StatusCreated, EventSendToBar, "", true},
		{"Cannot pay twice", StatusPaid, EventPayOrder, "", true},
		{"Cannot cancel served", StatusServed, EventCancelOrder, "", true},
		{"Cannot modify locked", StatusLocked, EventPayOrder, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotState, err := sm.Transition(tt.fromState, tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotState != tt.wantState {
				t.Errorf("Transition() = %v, want %v", gotState, tt.wantState)
			}
		})
	}
}

func TestOrderStateMachine_ValidateTransition(t *testing.T) {
	sm := NewOrderStateMachine()

	tests := []struct {
		name    string
		order   *Order
		event   OrderEvent
		wantErr bool
	}{
		{
			name: "Valid payment",
			order: &Order{
				Status: StatusCreated,
				Total:  100.0,
			},
			event:   EventPayOrder,
			wantErr: false,
		},
		{
			name: "Cannot pay zero total",
			order: &Order{
				Status: StatusCreated,
				Total:  0,
			},
			event:   EventPayOrder,
			wantErr: true,
		},
		{
			name: "Valid send to bar",
			order: &Order{
				Status:     StatusPaid,
				Total:      100.0,
				AmountPaid: 100.0,
				Items:      []OrderItem{{Name: "Coffee", Quantity: 1, Price: 100.0}},
			},
			event:   EventSendToBar,
			wantErr: false,
		},
		{
			name: "Cannot send empty order",
			order: &Order{
				Status:     StatusPaid,
				Total:      100.0,
				AmountPaid: 100.0,
				Items:      []OrderItem{},
			},
			event:   EventSendToBar,
			wantErr: true,
		},
		{
			name: "Cannot cancel served order",
			order: &Order{
				Status: StatusServed,
			},
			event:   EventCancelOrder,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sm.ValidateTransition(tt.order, tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTransition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderStateMachine_CanModifyOrder(t *testing.T) {
	sm := NewOrderStateMachine()

	tests := []struct {
		name   string
		status OrderStatus
		want   bool
	}{
		{"Can modify created", StatusCreated, true},
		{"Can modify paid", StatusPaid, true},
		{"Cannot modify queued", StatusQueued, false},
		{"Cannot modify in progress", StatusInProgress, false},
		{"Cannot modify ready", StatusReady, false},
		{"Cannot modify served", StatusServed, false},
		{"Cannot modify locked", StatusLocked, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sm.CanModifyOrder(tt.status); got != tt.want {
				t.Errorf("CanModifyOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderStateMachine_CanCancel(t *testing.T) {
	sm := NewOrderStateMachine()

	tests := []struct {
		name   string
		status OrderStatus
		want   bool
	}{
		{"Can cancel created", StatusCreated, true},
		{"Can cancel paid", StatusPaid, true},
		{"Can cancel queued", StatusQueued, true},
		{"Can cancel in progress", StatusInProgress, true},
		{"Cannot cancel ready", StatusReady, false},
		{"Cannot cancel served", StatusServed, false},
		{"Cannot cancel locked", StatusLocked, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sm.CanCancel(tt.status); got != tt.want {
				t.Errorf("CanCancel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderStateMachine_GetOrderProgress(t *testing.T) {
	sm := NewOrderStateMachine()

	tests := []struct {
		name   string
		status OrderStatus
		want   int
	}{
		{"Created progress", StatusCreated, 0},
		{"Paid progress", StatusPaid, 20},
		{"Queued progress", StatusQueued, 40},
		{"In progress", StatusInProgress, 60},
		{"Ready progress", StatusReady, 80},
		{"Served progress", StatusServed, 100},
		{"Locked progress", StatusLocked, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sm.GetOrderProgress(tt.status); got != tt.want {
				t.Errorf("GetOrderProgress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderStateMachine_IsTerminalState(t *testing.T) {
	sm := NewOrderStateMachine()

	tests := []struct {
		name   string
		status OrderStatus
		want   bool
	}{
		{"Created not terminal", StatusCreated, false},
		{"Paid not terminal", StatusPaid, false},
		{"Served not terminal (can lock)", StatusServed, false},
		{"Locked terminal", StatusLocked, true},
		{"Cancelled terminal", StatusCancelled, true},
		{"Refunded terminal", StatusRefunded, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sm.IsTerminalState(tt.status); got != tt.want {
				t.Errorf("IsTerminalState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderStateMachine_GetNextExpectedAction(t *testing.T) {
	sm := NewOrderStateMachine()

	tests := []struct {
		name   string
		status OrderStatus
		want   string
	}{
		{"Created next", StatusCreated, "Payment required"},
		{"Paid next", StatusPaid, "Send to bar"},
		{"Queued next", StatusQueued, "Start preparing"},
		{"In progress next", StatusInProgress, "Mark as ready"},
		{"Ready next", StatusReady, "Serve to customer"},
		{"Served next", StatusServed, "Order completed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sm.GetNextExpectedAction(tt.status); got != tt.want {
				t.Errorf("GetNextExpectedAction() = %v, want %v", got, tt.want)
			}
		})
	}
}
