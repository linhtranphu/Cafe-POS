package http

import (
	"net/http"

	"cafe-pos/backend/domain"
	"cafe-pos/backend/domain/cashier"
	"cafe-pos/backend/domain/order"

	"github.com/gin-gonic/gin"
)

// StateMachineHandler handles HTTP requests for state machine information
type StateMachineHandler struct {
	smManager *domain.StateMachineManager
}

// NewStateMachineHandler creates a new handler
func NewStateMachineHandler(smManager *domain.StateMachineManager) *StateMachineHandler {
	return &StateMachineHandler{
		smManager: smManager,
	}
}

// GetCashierShiftStates returns all possible states and transitions for cashier shifts
func (h *StateMachineHandler) GetCashierShiftStates(c *gin.Context) {
	sm := h.smManager.GetCashierShiftStateMachine()
	
	states := []string{
		string(cashier.CashierShiftOpen),
		string(cashier.CashierShiftClosureInitiated),
		string(cashier.CashierShiftClosed),
	}
	
	events := []string{
		string(cashier.EventInitiateClosure),
		string(cashier.EventRecordActualCash),
		string(cashier.EventDocumentVariance),
		string(cashier.EventConfirmResponsibility),
		string(cashier.EventCloseShift),
		string(cashier.EventCancelClosure),
	}
	
	// Build transition map
	transitions := make(map[string][]string)
	for _, state := range states {
		status := cashier.CashierShiftStatus(state)
		validEvents := sm.GetValidEvents(status)
		
		var eventStrings []string
		for _, event := range validEvents {
			eventStrings = append(eventStrings, string(event))
		}
		transitions[state] = eventStrings
	}
	
	c.JSON(http.StatusOK, gin.H{
		"states":      states,
		"events":      events,
		"transitions": transitions,
	})
}

// GetOrderStates returns all possible states and transitions for orders
func (h *StateMachineHandler) GetOrderStates(c *gin.Context) {
	sm := h.smManager.GetOrderStateMachine()
	
	states := []string{
		string(order.StatusCreated),
		string(order.StatusPaid),
		string(order.StatusQueued),
		string(order.StatusInProgress),
		string(order.StatusReady),
		string(order.StatusServed),
		string(order.StatusCancelled),
		string(order.StatusRefunded),
		string(order.StatusLocked),
	}
	
	events := []string{
		string(order.EventCreateOrder),
		string(order.EventPayOrder),
		string(order.EventSendToBar),
		string(order.EventStartPreparing),
		string(order.EventMarkReady),
		string(order.EventServeOrder),
		string(order.EventCancelOrder),
		string(order.EventRefundOrder),
		string(order.EventLockOrder),
	}
	
	// Build transition map
	transitions := make(map[string][]string)
	for _, state := range states {
		status := order.OrderStatus(state)
		validEvents := sm.GetValidEvents(status)
		
		var eventStrings []string
		for _, event := range validEvents {
			eventStrings = append(eventStrings, string(event))
		}
		transitions[state] = eventStrings
	}
	
	c.JSON(http.StatusOK, gin.H{
		"states":      states,
		"events":      events,
		"transitions": transitions,
	})
}

// GetWaiterShiftStates returns all possible states and transitions for waiter shifts
func (h *StateMachineHandler) GetWaiterShiftStates(c *gin.Context) {
	sm := h.smManager.GetWaiterShiftStateMachine()
	
	states := []string{
		string(order.ShiftOpen),
		string(order.ShiftClosed),
	}
	
	events := []string{
		string(order.EventStartShift),
		string(order.EventEndShift),
		string(order.EventCloseShift),
	}
	
	// Build transition map
	transitions := make(map[string][]string)
	for _, state := range states {
		status := order.ShiftStatus(state)
		validEvents := sm.GetValidEvents(status)
		
		var eventStrings []string
		for _, event := range validEvents {
			eventStrings = append(eventStrings, string(event))
		}
		transitions[state] = eventStrings
	}
	
	c.JSON(http.StatusOK, gin.H{
		"states":      states,
		"events":      events,
		"transitions": transitions,
	})
}

// GetAllStateMachines returns information about all state machines
func (h *StateMachineHandler) GetAllStateMachines(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"cashier_shift": gin.H{
			"description": "State machine for cashier shift lifecycle",
			"endpoint":    "/api/state-machines/cashier-shift",
		},
		"waiter_shift": gin.H{
			"description": "State machine for waiter/barista shift lifecycle",
			"endpoint":    "/api/state-machines/waiter-shift",
		},
		"order": gin.H{
			"description": "State machine for order lifecycle",
			"endpoint":    "/api/state-machines/order",
		},
	})
}
