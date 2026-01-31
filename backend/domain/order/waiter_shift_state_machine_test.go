package order

import (
	"testing"
	"time"
)

func TestShiftStateMachine_ValidTransitions(t *testing.T) {
	sm := NewShiftStateMachine()

	tests := []struct {
		name      string
		fromState ShiftStatus
		event     ShiftEvent
		wantState ShiftStatus
		wantErr   bool
	}{
		// Valid transitions
		{"Start shift", ShiftOpen, EventStartShift, ShiftOpen, false},
		{"End shift", ShiftOpen, EventEndShift, ShiftClosed, false},
		
		// Invalid transitions
		{"Cannot end closed shift", ShiftClosed, EventEndShift, "", true},
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

func TestShiftStateMachine_ValidateShiftStart(t *testing.T) {
	sm := NewShiftStateMachine()

	tests := []struct {
		name          string
		existingShift *Shift
		wantErr       bool
	}{
		{
			name:          "No existing shift",
			existingShift: nil,
			wantErr:       false,
		},
		{
			name: "Existing closed shift",
			existingShift: &Shift{
				Status: ShiftClosed,
			},
			wantErr: false,
		},
		{
			name: "Existing open shift",
			existingShift: &Shift{
				Status: ShiftOpen,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sm.ValidateShiftStart(tt.existingShift)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateShiftStart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShiftStateMachine_ValidateShiftEnd(t *testing.T) {
	sm := NewShiftStateMachine()

	tests := []struct {
		name    string
		shift   *Shift
		wantErr bool
	}{
		{
			name: "Can end open shift",
			shift: &Shift{
				Status: ShiftOpen,
			},
			wantErr: false,
		},
		{
			name: "Cannot end closed shift",
			shift: &Shift{
				Status: ShiftClosed,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sm.ValidateShiftEnd(tt.shift)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateShiftEnd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShiftStateMachine_CanStartShift(t *testing.T) {
	sm := NewShiftStateMachine()

	tests := []struct {
		name          string
		existingShift *Shift
		want          bool
	}{
		{"No existing shift", nil, true},
		{"Existing closed shift", &Shift{Status: ShiftClosed}, true},
		{"Existing open shift", &Shift{Status: ShiftOpen}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sm.CanStartShift(tt.existingShift); got != tt.want {
				t.Errorf("CanStartShift() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShiftStateMachine_GetShiftDuration(t *testing.T) {
	sm := NewShiftStateMachine()
	now := time.Now()

	tests := []struct {
		name  string
		shift *Shift
		want  float64
	}{
		{
			name: "Open shift 1 hour",
			shift: &Shift{
				Status:    ShiftOpen,
				StartedAt: now.Add(-1 * time.Hour),
			},
			want: 1.0,
		},
		{
			name: "Closed shift 2 hours",
			shift: &Shift{
				Status:    ShiftClosed,
				StartedAt: now.Add(-2 * time.Hour),
				EndedAt:   &now,
			},
			want: 2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sm.GetShiftDuration(tt.shift)
			// Allow 0.1 hour tolerance for test execution time
			if got < tt.want-0.1 || got > tt.want+0.1 {
				t.Errorf("GetShiftDuration() = %v, want ~%v", got, tt.want)
			}
		})
	}
}

func TestShiftStateMachine_IsTerminalState(t *testing.T) {
	sm := NewShiftStateMachine()

	tests := []struct {
		name   string
		status ShiftStatus
		want   bool
	}{
		{"Open not terminal", ShiftOpen, false},
		{"Closed terminal", ShiftClosed, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sm.IsTerminalState(tt.status); got != tt.want {
				t.Errorf("IsTerminalState() = %v, want %v", got, tt.want)
			}
		})
	}
}
