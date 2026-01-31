package order

import (
	"fmt"
)

// ShiftEvent represents events that can trigger shift state transitions
type ShiftEvent string

const (
	// Shift lifecycle events
	EventStartShift ShiftEvent = "START_SHIFT"
	EventEndShift   ShiftEvent = "END_SHIFT"
	EventCloseShift ShiftEvent = "CLOSE_SHIFT"
)

// ShiftStateMachine manages state transitions for waiter/barista shifts
type ShiftStateMachine struct {
	// Valid transitions map: current_state -> event -> next_state
	transitions map[ShiftStatus]map[ShiftEvent]ShiftStatus
}

// NewShiftStateMachine creates a new state machine for shifts
func NewShiftStateMachine() *ShiftStateMachine {
	sm := &ShiftStateMachine{
		transitions: make(map[ShiftStatus]map[ShiftEvent]ShiftStatus),
	}
	
	sm.defineTransitions()
	
	return sm
}

// defineTransitions sets up all valid state transitions
func (sm *ShiftStateMachine) defineTransitions() {
	// From OPEN state
	sm.transitions[ShiftOpen] = map[ShiftEvent]ShiftStatus{
		EventEndShift: ShiftClosed,
	}
	
	// CLOSED is a terminal state
	sm.transitions[ShiftClosed] = map[ShiftEvent]ShiftStatus{}
}

// CanTransition checks if a transition is valid
func (sm *ShiftStateMachine) CanTransition(currentState ShiftStatus, event ShiftEvent) bool {
	if stateTransitions, exists := sm.transitions[currentState]; exists {
		_, canTransition := stateTransitions[event]
		return canTransition
	}
	return false
}

// Transition attempts to transition from current state to next state via event
func (sm *ShiftStateMachine) Transition(currentState ShiftStatus, event ShiftEvent) (ShiftStatus, error) {
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
func (sm *ShiftStateMachine) GetValidEvents(currentState ShiftStatus) []ShiftEvent {
	var events []ShiftEvent
	
	if stateTransitions, exists := sm.transitions[currentState]; exists {
		for event := range stateTransitions {
			events = append(events, event)
		}
	}
	
	return events
}

// IsTerminalState checks if a state is terminal (no further transitions)
func (sm *ShiftStateMachine) IsTerminalState(state ShiftStatus) bool {
	transitions, exists := sm.transitions[state]
	return !exists || len(transitions) == 0
}

// ValidateShiftEnd validates if a shift can be ended
func (sm *ShiftStateMachine) ValidateShiftEnd(shift *Shift) error {
	if shift.Status != ShiftOpen {
		return fmt.Errorf("can only end shifts in OPEN status")
	}
	
	// Additional business rules can be added here
	// For example: check if all orders are completed
	
	return nil
}

// CanStartShift checks if a new shift can be started
func (sm *ShiftStateMachine) CanStartShift(existingShift *Shift) bool {
	// Can start a new shift if no existing shift or existing shift is closed
	return existingShift == nil || existingShift.Status == ShiftClosed
}

// ValidateShiftStart validates if a shift can be started
func (sm *ShiftStateMachine) ValidateShiftStart(existingShift *Shift) error {
	if existingShift != nil && existingShift.Status == ShiftOpen {
		return fmt.Errorf("cannot start new shift: user already has an open shift")
	}
	return nil
}

// GetShiftDuration calculates shift duration in hours
func (sm *ShiftStateMachine) GetShiftDuration(shift *Shift) float64 {
	if shift.EndedAt == nil {
		// Shift is still open
		return 0
	}
	
	duration := shift.EndedAt.Sub(shift.StartedAt)
	return duration.Hours()
}

// GetNextExpectedAction returns the next expected action for a shift
func (sm *ShiftStateMachine) GetNextExpectedAction(currentState ShiftStatus) string {
	switch currentState {
	case ShiftOpen:
		return "Shift is active"
	case ShiftClosed:
		return "Shift is closed"
	default:
		return "Unknown state"
	}
}
