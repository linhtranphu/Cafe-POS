package cashier

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CashierShiftStatus represents the current status of a cashier shift
type CashierShiftStatus string

const (
	// CashierShiftOpen indicates the shift is currently active
	CashierShiftOpen CashierShiftStatus = "OPEN"
	// CashierShiftClosureInitiated indicates the closure process has started
	CashierShiftClosureInitiated CashierShiftStatus = "CLOSURE_INITIATED"
	// CashierShiftClosed indicates the shift has been finalized and locked
	CashierShiftClosed CashierShiftStatus = "CLOSED"
)

// CashierShift is the aggregate root managing the cashier shift lifecycle and business rules.
// It encapsulates all shift data, financial information, and enforces business rules
// for the shift closure process.
//
// This is separate from waiter/barista shifts to maintain clear separation of concerns.
// Cashier shifts have complex closure workflows with cash reconciliation, variance handling,
// and audit trails.
type CashierShift struct {
	// ID is the unique identifier for the cashier shift
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// CashierID is the unique identifier of the cashier assigned to this shift
	CashierID primitive.ObjectID `json:"cashier_id" bson:"cashier_id"`

	// CashierName is the name of the cashier for display purposes
	CashierName string `json:"cashier_name" bson:"cashier_name"`

	// StartTime is when the shift began
	StartTime time.Time `json:"start_time" bson:"start_time"`

	// EndTime is when the shift was closed (nil if still open)
	EndTime *time.Time `json:"end_time,omitempty" bson:"end_time,omitempty"`

	// Status represents the current state of the shift
	Status CashierShiftStatus `json:"status" bson:"status"`

	// StartingFloat is the initial cash amount in the drawer at shift start
	StartingFloat float64 `json:"starting_float" bson:"starting_float"`

	// SystemCash is the theoretical cash amount calculated by the POS based on transactions
	SystemCash float64 `json:"system_cash" bson:"system_cash"`

	// ActualCash is the physical cash amount counted by the cashier (nil until recorded)
	ActualCash *float64 `json:"actual_cash,omitempty" bson:"actual_cash,omitempty"`

	// Variance contains the variance calculation and documentation (nil until calculated)
	Variance *Variance `json:"variance,omitempty" bson:"variance,omitempty"`

	// Confirmation contains the cashier's responsibility confirmation (nil until confirmed)
	Confirmation *ResponsibilityConfirmation `json:"confirmation,omitempty" bson:"confirmation,omitempty"`

	// AuditLog maintains an immutable record of all actions performed during the shift
	AuditLog []AuditLogEntry `json:"audit_log" bson:"audit_log"`

	// Cash handover tracking fields
	ReceivedCash     float64 `bson:"received_cash" json:"received_cash"`         // Total cash received from handovers
	TotalDiscrepancy float64 `bson:"total_discrepancy" json:"total_discrepancy"` // Total discrepancy from handovers
	HandoverCount    int     `bson:"handover_count" json:"handover_count"`       // Number of handovers received
	DiscrepancyCount int     `bson:"discrepancy_count" json:"discrepancy_count"` // Number of handovers with discrepancies

	// CreatedAt is when the shift record was created
	CreatedAt time.Time `json:"created_at" bson:"created_at"`

	// UpdatedAt is when the shift record was last updated
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// NewCashierShift creates a new CashierShift with the specified parameters.
// The shift is initialized in the Open status with an empty audit log.
//
// Parameters:
//   - cashierID: Unique identifier of the cashier
//   - cashierName: Name of the cashier for display
//   - startingFloat: Initial cash amount in the drawer
//
// Returns:
//   - *CashierShift: A new cashier shift in Open status
func NewCashierShift(cashierID primitive.ObjectID, cashierName string, startingFloat float64) *CashierShift {
	now := time.Now()
	return &CashierShift{
		CashierID:     cashierID,
		CashierName:   cashierName,
		StartTime:     now,
		Status:        CashierShiftOpen,
		StartingFloat: startingFloat,
		SystemCash:    startingFloat, // Initially equals starting float
		AuditLog:      []AuditLogEntry{},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// InitiateClosure starts the shift closure process by validating the current status
// and transitioning to ClosureInitiated state. This method enforces the state machine
// transition rules and adds an audit log entry for traceability.
//
// The method validates that:
//   - The shift is currently in Open status
//   - All audit log entry parameters are valid (userID, deviceID, timestamp)
//
// On success:
//   - Transitions status from Open to ClosureInitiated
//   - Adds an audit log entry with action "closure_initiated"
//
// Returns:
//   - error: if the shift is not in Open status or audit log entry creation fails
func (cs *CashierShift) InitiateClosure(userID, deviceID string, timestamp time.Time) error {
	// Validate current status is Open
	if cs.Status != CashierShiftOpen {
		return errors.New("cannot initiate closure: shift status must be Open")
	}

	// Create audit log entry first to validate parameters before changing state
	auditEntry, err := NewAuditLogEntry("closure_initiated", userID, deviceID, timestamp, nil)
	if err != nil {
		return err
	}

	// Transition status to ClosureInitiated
	cs.Status = CashierShiftClosureInitiated

	// Add audit log entry
	cs.AuditLog = append(cs.AuditLog, *auditEntry)
	cs.UpdatedAt = timestamp

	return nil
}

// RecordActualCash records the physical cash amount counted by the cashier and calculates
// the variance between actual cash and system cash. This method validates the input,
// creates a Variance value object, and adds an audit log entry for traceability.
//
// The method validates that:
//   - The actual cash amount is non-negative
//   - The actual cash has at most 2 decimal places
//   - All audit log entry parameters are valid (userID, deviceID, timestamp)
//
// On success:
//   - Stores the actual cash amount
//   - Creates and stores a Variance value object
//   - Adds an audit log entry with action "actual_cash_recorded"
//   - Returns the created Variance object
//
// Parameters:
//   - actualCash: The physical cash amount counted by the cashier
//   - userID: The unique identifier of the user recording the cash
//   - deviceID: The unique identifier of the device used
//   - timestamp: The time when the cash was recorded
//
// Returns:
//   - *Variance: The calculated variance object
//   - error: if validation fails or audit log entry creation fails
func (cs *CashierShift) RecordActualCash(actualCash float64, userID, deviceID string, timestamp time.Time) (*Variance, error) {
	// Validate actual cash is non-negative
	if actualCash < 0 {
		return nil, errors.New("actual cash must be non-negative")
	}

	// Validate decimal precision (2 places)
	// Multiply by 100 and check if it's an integer
	scaled := actualCash * 100
	if scaled != float64(int64(scaled)) {
		return nil, errors.New("actual cash must have at most 2 decimal places")
	}

	// Create Variance value object
	variance := NewVariance(cs.SystemCash, actualCash)

	// Create audit log entry with actual cash data
	auditData := map[string]interface{}{
		"actual_cash": actualCash,
	}
	auditEntry, err := NewAuditLogEntry("actual_cash_recorded", userID, deviceID, timestamp, auditData)
	if err != nil {
		return nil, err
	}

	// Store actual cash and variance
	cs.ActualCash = &actualCash
	cs.Variance = variance

	// Add audit log entry
	cs.AuditLog = append(cs.AuditLog, *auditEntry)
	cs.UpdatedAt = timestamp

	return variance, nil
}

// DocumentVariance records the reason and notes for a variance.
// This method validates that a variance exists and is non-zero, then delegates
// to the Variance value object's Document method, and adds an audit log entry.
//
// The method validates that:
//   - A variance has been calculated (Variance is not nil)
//   - The variance is non-zero (requires documentation)
//   - All audit log entry parameters are valid (userID, deviceID, timestamp)
//   - The notes meet the minimum length requirement (10 characters, validated by Variance.Document)
//
// On success:
//   - Calls Variance.Document() to record reason and notes
//   - Adds an audit log entry with action "variance_documented"
//
// Parameters:
//   - reason: The selected variance reason
//   - notes: Detailed explanation (minimum 10 characters)
//   - userID: The unique identifier of the user documenting the variance
//   - deviceID: The unique identifier of the device used
//   - timestamp: The time when the variance was documented
//
// Returns:
//   - error: if validation fails or audit log entry creation fails
func (cs *CashierShift) DocumentVariance(reason VarianceReason, notes string, userID, deviceID string, timestamp time.Time) error {
	// Validate variance exists
	if cs.Variance == nil {
		return errors.New("cannot document variance: no variance has been calculated")
	}

	// Validate variance is non-zero (requires documentation)
	if !cs.Variance.RequiresDocumentation() {
		return errors.New("cannot document variance: variance is zero and does not require documentation")
	}

	// Create audit log entry first to validate parameters before modifying variance
	auditEntry, err := NewAuditLogEntry("variance_documented", userID, deviceID, timestamp, nil)
	if err != nil {
		return err
	}

	// Call variance.Document() method to validate and record documentation
	if err := cs.Variance.Document(reason, notes); err != nil {
		return err
	}

	// Add audit log entry
	cs.AuditLog = append(cs.AuditLog, *auditEntry)
	cs.UpdatedAt = timestamp

	return nil
}

// ConfirmResponsibility records the cashier's formal confirmation of responsibility
// for the shift's financial data. This method creates a ResponsibilityConfirmation
// value object, stores it in the aggregate, and adds an audit log entry for traceability.
//
// The method validates that:
//   - All parameters are valid (userID, deviceID, timestamp) via ResponsibilityConfirmation constructor
//   - All audit log entry parameters are valid (userID, deviceID, timestamp)
//
// On success:
//   - Creates and stores a ResponsibilityConfirmation value object
//   - Adds an audit log entry with action "responsibility_confirmed"
//
// Parameters:
//   - userID: The unique identifier of the cashier confirming responsibility
//   - deviceID: The unique identifier of the device used for confirmation
//   - timestamp: The time when the confirmation was made
//
// Returns:
//   - error: if validation fails or audit log entry creation fails
func (cs *CashierShift) ConfirmResponsibility(userID, deviceID string, timestamp time.Time) error {
	// Create ResponsibilityConfirmation value object (validates parameters)
	confirmation, err := NewResponsibilityConfirmation(userID, deviceID, timestamp)
	if err != nil {
		return err
	}

	// Create audit log entry to validate parameters before storing confirmation
	auditEntry, err := NewAuditLogEntry("responsibility_confirmed", userID, deviceID, timestamp, nil)
	if err != nil {
		return err
	}

	// Store confirmation in aggregate
	cs.Confirmation = confirmation

	// Add audit log entry
	cs.AuditLog = append(cs.AuditLog, *auditEntry)
	cs.UpdatedAt = timestamp

	return nil
}

// CanClose validates all preconditions required to close the shift.
// This method checks that all necessary steps have been completed before
// allowing the shift to be closed.
//
// The method validates that:
//   - The shift status is ClosureInitiated
//   - The cashier has confirmed responsibility (Confirmation is not nil)
//   - If variance is non-zero, it has been documented
//
// Returns:
//   - error: with a descriptive message if any validation fails, nil if all checks pass
func (cs *CashierShift) CanClose() error {
	// Check status is ClosureInitiated
	if cs.Status != CashierShiftClosureInitiated {
		return errors.New("cannot close shift: status must be ClosureInitiated")
	}

	// Check confirmation exists
	if cs.Confirmation == nil {
		return errors.New("cannot close shift: responsibility confirmation is required")
	}

	// Check variance is documented if non-zero
	if cs.Variance != nil && cs.Variance.RequiresDocumentation() {
		if cs.Variance.Reason == nil || cs.Variance.Notes == "" {
			return errors.New("cannot close shift: variance must be documented with reason and notes")
		}
	}

	return nil
}

// Close finalizes the shift by validating all preconditions, setting the end time,
// transitioning to Closed status, and adding an audit log entry. Once closed, the
// shift data becomes immutable.
//
// The method validates that:
//   - All preconditions are met via CanClose()
//   - All audit log entry parameters are valid (userID, deviceID, timestamp)
//
// On success:
//   - Sets the EndTime to the provided timestamp
//   - Transitions status to Closed
//   - Adds an audit log entry with action "shift_closed"
//
// Parameters:
//   - userID: The unique identifier of the user closing the shift
//   - deviceID: The unique identifier of the device used
//   - timestamp: The time when the shift was closed (becomes EndTime)
//
// Returns:
//   - error: if validation fails or audit log entry creation fails
func (cs *CashierShift) Close(userID, deviceID string, timestamp time.Time) error {
	// Validate all preconditions via CanClose()
	if err := cs.CanClose(); err != nil {
		return err
	}

	// Create audit log entry to validate parameters before changing state
	auditEntry, err := NewAuditLogEntry("shift_closed", userID, deviceID, timestamp, nil)
	if err != nil {
		return err
	}

	// Set EndTime timestamp
	cs.EndTime = &timestamp

	// Transition status to Closed
	cs.Status = CashierShiftClosed

	// Add audit log entry
	cs.AuditLog = append(cs.AuditLog, *auditEntry)
	cs.UpdatedAt = timestamp

	return nil
}

// UpdateSystemCash updates the system cash amount based on transactions.
// This should be called when calculating the expected cash from orders and handovers.
//
// Parameters:
//   - systemCash: The calculated system cash amount
func (cs *CashierShift) UpdateSystemCash(systemCash float64) {
	cs.SystemCash = systemCash
	cs.UpdatedAt = time.Now()
}

// UpdateCashAfterHandover updates the cash amounts after receiving a handover
func (cs *CashierShift) UpdateCashAfterHandover(receivedAmount, discrepancyAmount float64, hasDiscrepancy bool) {
	cs.ReceivedCash += receivedAmount
	cs.TotalDiscrepancy += discrepancyAmount
	cs.HandoverCount++
	
	if hasDiscrepancy {
		cs.DiscrepancyCount++
	}
	
	// Update system cash to include received amount
	cs.SystemCash += receivedAmount
	cs.UpdatedAt = time.Now()
}
