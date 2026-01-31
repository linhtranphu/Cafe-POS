package cashier

import (
	"errors"
	"fmt"
)

// ShiftEvent represents events that can trigger state transitions
type ShiftEvent string

const (
	// Shift events
	EventInitiateClosure       ShiftEvent = "INITIATE_CLOSURE"
	EventRecordActualCash      ShiftEvent = "RECORD_ACTUAL_CASH"
	EventDocumentVariance      ShiftEvent = "DOCUMENT_VARIANCE"
	EventConfirmResponsibility ShiftEvent = "CONFIRM_RESPONSIBILITY"
	EventCloseShift            ShiftEvent = "CLOSE_SHIFT"
	EventCancelClosure         ShiftEvent = "CANCEL_CLOSURE"
)

// ShiftStateMachine manages state transitions for cashier shifts
type ShiftStateMachine struct {
	// Valid transitions map: current_state -> event -> next_state
	transitions map[CashierShiftStatus]map[ShiftEvent]CashierShiftStatus
}

// NewShiftStateMachine creates a new state machine for cashier shifts
func NewShiftStateMachine() *ShiftStateMachine {
	sm := &ShiftStateMachine{
		transitions: make(map[CashierShiftStatus]map[ShiftEvent]CashierShiftStatus),
	}
	
	// Define valid state transitions
	sm.defineTransitions()
	
	return sm
}

// defineTransitions sets up all valid state transitions
func (sm *ShiftStateMachine) defineTransitions() {
	// From OPEN state
	sm.transitions[CashierShiftOpen] = map[ShiftEvent]CashierShiftStatus{
		EventInitiateClosure: CashierShiftClosureInitiated,
	}
	
	// From CLOSURE_INITIATED state
	sm.transitions[CashierShiftClosureInitiated] = map[ShiftEvent]CashierShiftStatus{
		EventCloseShift:    CashierShiftClosed,
		EventCancelClosure: CashierShiftOpen, // Allow canceling closure
	}
	
	// CLOSED is a terminal state - no transitions allowed
	sm.transitions[CashierShiftClosed] = map[ShiftEvent]CashierShiftStatus{}
}

// CanTransition checks if a transition is valid
func (sm *ShiftStateMachine) CanTransition(currentState CashierShiftStatus, event ShiftEvent) bool {
	if stateTransitions, exists := sm.transitions[currentState]; exists {
		_, canTransition := stateTransitions[event]
		return canTransition
	}
	return false
}

// Transition attempts to transition from current state to next state via event
func (sm *ShiftStateMachine) Transition(currentState CashierShiftStatus, event ShiftEvent) (CashierShiftStatus, error) {
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
func (sm *ShiftStateMachine) GetValidEvents(currentState CashierShiftStatus) []ShiftEvent {
	var events []ShiftEvent
	
	if stateTransitions, exists := sm.transitions[currentState]; exists {
		for event := range stateTransitions {
			events = append(events, event)
		}
	}
	
	return events
}

// ValidateShiftWorkflow validates the entire shift closure workflow
func (sm *ShiftStateMachine) ValidateShiftWorkflow(shift *CashierShift) error {
	// Check if shift is in a valid state
	switch shift.Status {
	case CashierShiftOpen:
		// Open shift is valid, can initiate closure
		return nil
		
	case CashierShiftClosureInitiated:
		// Check if all required steps are completed
		if shift.ActualCash == nil {
			return errors.New("actual cash must be recorded before closing")
		}
		
		// If there's a variance, it must be documented
		if shift.Variance != nil && shift.Variance.RequiresDocumentation() {
			if shift.Variance.Reason == nil || shift.Variance.Notes == "" {
				return errors.New("variance must be documented before closing")
			}
		}
		
		// Responsibility must be confirmed
		if shift.Confirmation == nil {
			return errors.New("responsibility must be confirmed before closing")
		}
		
		return nil
		
	case CashierShiftClosed:
		// Closed shift cannot be modified
		return errors.New("shift is already closed and cannot be modified")
		
	default:
		return fmt.Errorf("unknown shift status: %s", shift.Status)
	}
}

// ValidateRecordActualCash validates if actual cash can be recorded
func (sm *ShiftStateMachine) ValidateRecordActualCash(shift *CashierShift) error {
	if shift.Status != CashierShiftClosureInitiated {
		return errors.New("can only record actual cash after initiating closure")
	}
	
	if shift.ActualCash != nil {
		return errors.New("actual cash has already been recorded")
	}
	
	return nil
}

// ValidateDocumentVariance validates if variance can be documented
func (sm *ShiftStateMachine) ValidateDocumentVariance(shift *CashierShift) error {
	if shift.Status != CashierShiftClosureInitiated {
		return errors.New("can only document variance after initiating closure")
	}
	
	if shift.ActualCash == nil {
		return errors.New("must record actual cash before documenting variance")
	}
	
	if shift.Variance == nil {
		return errors.New("no variance to document")
	}
	
	if !shift.Variance.RequiresDocumentation() {
		return errors.New("variance is zero and does not require documentation")
	}
	
	if shift.Variance.Reason != nil {
		return errors.New("variance has already been documented")
	}
	
	return nil
}

// ValidateConfirmResponsibility validates if responsibility can be confirmed
func (sm *ShiftStateMachine) ValidateConfirmResponsibility(shift *CashierShift) error {
	if shift.Status != CashierShiftClosureInitiated {
		return errors.New("can only confirm responsibility after initiating closure")
	}
	
	if shift.ActualCash == nil {
		return errors.New("must record actual cash before confirming responsibility")
	}
	
	// If there's a variance, it must be documented first
	if shift.Variance != nil && shift.Variance.RequiresDocumentation() {
		if shift.Variance.Reason == nil {
			return errors.New("must document variance before confirming responsibility")
		}
	}
	
	if shift.Confirmation != nil {
		return errors.New("responsibility has already been confirmed")
	}
	
	return nil
}

// CanCancelClosure checks if closure can be cancelled
func (sm *ShiftStateMachine) CanCancelClosure(shift *CashierShift) bool {
	// Can only cancel if in CLOSURE_INITIATED state and no critical steps completed
	if shift.Status != CashierShiftClosureInitiated {
		return false
	}
	
	// Cannot cancel if actual cash has been recorded
	if shift.ActualCash != nil {
		return false
	}
	
	return true
}

// IsTerminalState checks if a state is terminal (no further transitions)
func (sm *ShiftStateMachine) IsTerminalState(state CashierShiftStatus) bool {
	transitions, exists := sm.transitions[state]
	return !exists || len(transitions) == 0
}

// GetNextRequiredStep returns the next required step in the workflow
func (sm *ShiftStateMachine) GetNextRequiredStep(shift *CashierShift) string {
	switch shift.Status {
	case CashierShiftOpen:
		return "Initiate closure"
		
	case CashierShiftClosureInitiated:
		if shift.ActualCash == nil {
			return "Record actual cash"
		}
		
		if shift.Variance != nil && shift.Variance.RequiresDocumentation() {
			if shift.Variance.Reason == nil || shift.Variance.Notes == "" {
				return "Document variance"
			}
		}
		
		if shift.Confirmation == nil {
			return "Confirm responsibility"
		}
		
		return "Close shift"
		
	case CashierShiftClosed:
		return "Shift is closed"
		
	default:
		return "Unknown state"
	}
}
