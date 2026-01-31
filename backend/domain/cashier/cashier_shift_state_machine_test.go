package cashier

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCashierShiftStateMachine_ValidTransitions(t *testing.T) {
	sm := NewShiftStateMachine()

	tests := []struct {
		name      string
		fromState CashierShiftStatus
		event     ShiftEvent
		wantState CashierShiftStatus
		wantErr   bool
	}{
		// Valid transitions
		{"Initiate closure", CashierShiftOpen, EventInitiateClosure, CashierShiftClosureInitiated, false},
		{"Close shift", CashierShiftClosureInitiated, EventCloseShift, CashierShiftClosed, false},
		{"Cancel closure", CashierShiftClosureInitiated, EventCancelClosure, CashierShiftOpen, false},
		
		// Invalid transitions
		{"Cannot close from open", CashierShiftOpen, EventCloseShift, "", true},
		{"Cannot cancel from closed", CashierShiftClosed, EventCancelClosure, "", true},
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

func TestCashierShiftStateMachine_ValidateRecordActualCash(t *testing.T) {
	sm := NewShiftStateMachine()

	tests := []struct {
		name    string
		shift   *CashierShift
		wantErr bool
	}{
		{
			name: "Valid actual cash recorded",
			shift: &CashierShift{
				Status:     CashierShiftClosureInitiated,
				ActualCash: 1000.0,
			},
			wantErr: false,
		},
		{
			name: "Actual cash not recorded",
			shift: &CashierShift{
				Status:     CashierShiftClosureInitiated,
				ActualCash: 0,
			},
			wantErr: true,
		},
		{
			name: "Wrong status",
			shift: &CashierShift{
				Status:     CashierShiftOpen,
				ActualCash: 1000.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sm.ValidateRecordActualCash(tt.shift)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRecordActualCash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCashierShiftStateMachine_ValidateDocumentVariance(t *testing.T) {
	sm := NewShiftStateMachine()

	tests := []struct {
		name    string
		shift   *CashierShift
		wantErr bool
	}{
		{
			name: "No variance - no documentation needed",
			shift: &CashierShift{
				Status:       CashierShiftClosureInitiated,
				ActualCash:   1000.0,
				ExpectedCash: 1000.0,
			},
			wantErr: false,
		},
		{
			name: "Variance documented",
			shift: &CashierShift{
				Status:       CashierShiftClosureInitiated,
				ActualCash:   900.0,
				ExpectedCash: 1000.0,
				Variance: &Variance{
					Amount: -100.0,
					Reason: "Missing cash",
				},
			},
			wantErr: false,
		},
		{
			name: "Variance not documented",
			shift: &CashierShift{
				Status:       CashierShiftClosureInitiated,
				ActualCash:   900.0,
				ExpectedCash: 1000.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sm.ValidateDocumentVariance(tt.shift)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDocumentVariance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCashierShiftStateMachine_ValidateConfirmResponsibility(t *testing.T) {
	sm := NewShiftStateMachine()

	tests := []struct {
		name    string
		shift   *CashierShift
		wantErr bool
	}{
		{
			name: "Responsibility confirmed",
			shift: &CashierShift{
				Status: CashierShiftClosureInitiated,
				ResponsibilityConfirmation: &ResponsibilityConfirmation{
					Confirmed:   true,
					ConfirmedAt: time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Responsibility not confirmed",
			shift: &CashierShift{
				Status: CashierShiftClosureInitiated,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sm.ValidateConfirmResponsibility(tt.shift)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfirmResponsibility() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCashierShiftStateMachine_ValidateShiftWorkflow(t *testing.T) {
	sm := NewShiftStateMachine()
	now := time.Now()

	tests := []struct {
		name    string
		shift   *CashierShift
		wantErr bool
	}{
		{
			name: "Complete workflow",
			shift: &CashierShift{
				Status:       CashierShiftClosureInitiated,
				ActualCash:   1000.0,
				ExpectedCash: 1000.0,
				ResponsibilityConfirmation: &ResponsibilityConfirmation{
					Confirmed:   true,
					ConfirmedAt: now,
				},
			},
			wantErr: false,
		},
		{
			name: "Missing actual cash",
			shift: &CashierShift{
				Status:       CashierShiftClosureInitiated,
				ActualCash:   0,
				ExpectedCash: 1000.0,
			},
			wantErr: true,
		},
		{
			name: "Missing responsibility confirmation",
			shift: &CashierShift{
				Status:       CashierShiftClosureInitiated,
				ActualCash:   1000.0,
				ExpectedCash: 1000.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sm.ValidateShiftWorkflow(tt.shift)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateShiftWorkflow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCashierShiftStateMachine_CanCancelClosure(t *testing.T) {
	sm := NewShiftStateMachine()

	tests := []struct {
		name  string
		shift *CashierShift
		want  bool
	}{
		{
			name: "Can cancel before recording cash",
			shift: &CashierShift{
				Status:     CashierShiftClosureInitiated,
				ActualCash: 0,
			},
			want: true,
		},
		{
			name: "Cannot cancel after recording cash",
			shift: &CashierShift{
				Status:     CashierShiftClosureInitiated,
				ActualCash: 1000.0,
			},
			want: false,
		},
		{
			name: "Cannot cancel closed shift",
			shift: &CashierShift{
				Status: CashierShiftClosed,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sm.CanCancelClosure(tt.shift); got != tt.want {
				t.Errorf("CanCancelClosure() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCashierShiftStateMachine_GetNextRequiredStep(t *testing.T) {
	sm := NewShiftStateMachine()
	now := time.Now()

	tests := []struct {
		name  string
		shift *CashierShift
		want  string
	}{
		{
			name: "Need to record actual cash",
			shift: &CashierShift{
				Status:     CashierShiftClosureInitiated,
				ActualCash: 0,
			},
			want: "record_actual_cash",
		},
		{
			name: "Need to document variance",
			shift: &CashierShift{
				Status:       CashierShiftClosureInitiated,
				ActualCash:   900.0,
				ExpectedCash: 1000.0,
			},
			want: "document_variance",
		},
		{
			name: "Need to confirm responsibility",
			shift: &CashierShift{
				Status:       CashierShiftClosureInitiated,
				ActualCash:   1000.0,
				ExpectedCash: 1000.0,
			},
			want: "confirm_responsibility",
		},
		{
			name: "Ready to close",
			shift: &CashierShift{
				Status:       CashierShiftClosureInitiated,
				ActualCash:   1000.0,
				ExpectedCash: 1000.0,
				ResponsibilityConfirmation: &ResponsibilityConfirmation{
					Confirmed:   true,
					ConfirmedAt: now,
				},
			},
			want: "close_shift",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sm.GetNextRequiredStep(tt.shift); got != tt.want {
				t.Errorf("GetNextRequiredStep() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCashierShiftStateMachine_IsTerminalState(t *testing.T) {
	sm := NewShiftStateMachine()

	tests := []struct {
		name   string
		status CashierShiftStatus
		want   bool
	}{
		{"Open not terminal", CashierShiftOpen, false},
		{"Closure initiated not terminal", CashierShiftClosureInitiated, false},
		{"Closed terminal", CashierShiftClosed, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sm.IsTerminalState(tt.status); got != tt.want {
				t.Errorf("IsTerminalState() = %v, want %v", got, tt.want)
			}
		})
	}
}
