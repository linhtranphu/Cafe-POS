package cashier

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FundHandover represents a handover of funds from cashier back to the fund
// when closing a cashier shift. This is the final step in the cashier's
// responsibility chain.
//
// The cashier receives cash and transfer amounts from waiters during their shift.
// When closing the shift, they must hand over all managed funds back to the fund.
// This record documents that handover, including any variance between expected
// and actual cash amounts.
type FundHandover struct {
	// ID is the unique identifier for the fund handover
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// CashierShiftID links this handover to the cashier shift
	CashierShiftID primitive.ObjectID `json:"cashier_shift_id" bson:"cashier_shift_id"`

	// CashierID is the cashier who is handing over the funds
	CashierID primitive.ObjectID `json:"cashier_id" bson:"cashier_id"`

	// CashierName for display purposes
	CashierName string `json:"cashier_name" bson:"cashier_name"`

	// CashAmount is the actual cash amount handed over
	CashAmount float64 `json:"cash_amount" bson:"cash_amount"`

	// TransferAmount is the transfer amount recorded (no physical handover)
	TransferAmount float64 `json:"transfer_amount" bson:"transfer_amount"`

	// TotalAmount is the sum of cash and transfer
	TotalAmount float64 `json:"total_amount" bson:"total_amount"`

	// ExpectedCash is the theoretical cash amount (starting_float + received_cash)
	ExpectedCash float64 `json:"expected_cash" bson:"expected_cash"`

	// VarianceAmount is the difference between actual and expected cash
	VarianceAmount float64 `json:"variance_amount" bson:"variance_amount"`

	// VarianceReason is the selected reason for variance (if any)
	VarianceReason *VarianceReason `json:"variance_reason,omitempty" bson:"variance_reason,omitempty"`

	// VarianceNotes provides detailed explanation for variance
	VarianceNotes string `json:"variance_notes,omitempty" bson:"variance_notes,omitempty"`

	// ReceiverID is the person receiving the funds (nullable for future use)
	// Currently null, indicating handover to general fund
	// Future: can be set to a manager's ID
	ReceiverID *primitive.ObjectID `json:"receiver_id,omitempty" bson:"receiver_id,omitempty"`

	// ReceiverName for display purposes (nullable)
	ReceiverName string `json:"receiver_name,omitempty" bson:"receiver_name,omitempty"`

	// HandoverAt is when the handover occurred
	HandoverAt time.Time `json:"handover_at" bson:"handover_at"`

	// CreatedAt is when the record was created
	CreatedAt time.Time `json:"created_at" bson:"created_at"`

	// UpdatedAt is when the record was last updated
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// NewFundHandover creates a new fund handover record with the specified parameters.
// The variance is automatically calculated as (cashAmount - expectedCash).
//
// Parameters:
//   - cashierShiftID: The ID of the cashier shift being closed
//   - cashierID: The ID of the cashier handing over funds
//   - cashierName: The name of the cashier for display
//   - cashAmount: The actual cash amount being handed over
//   - transferAmount: The transfer amount being recorded
//   - expectedCash: The theoretical cash amount (starting_float + received_cash)
//
// Returns:
//   - *FundHandover: A new fund handover record
func NewFundHandover(
	cashierShiftID, cashierID primitive.ObjectID,
	cashierName string,
	cashAmount, transferAmount, expectedCash float64,
) *FundHandover {
	now := time.Now()
	variance := cashAmount - expectedCash

	return &FundHandover{
		CashierShiftID: cashierShiftID,
		CashierID:      cashierID,
		CashierName:    cashierName,
		CashAmount:     cashAmount,
		TransferAmount: transferAmount,
		TotalAmount:    cashAmount + transferAmount,
		ExpectedCash:   expectedCash,
		VarianceAmount: variance,
		HandoverAt:     now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// HasVariance returns true if there is a variance between expected and actual cash.
// A variance exists when the actual cash amount differs from the expected amount.
// Uses a tolerance of 0.01 (1 cent) to handle floating-point precision issues.
func (fh *FundHandover) HasVariance() bool {
	const tolerance = 0.01
	if fh.VarianceAmount < 0 {
		return -fh.VarianceAmount >= tolerance
	}
	return fh.VarianceAmount >= tolerance
}

// RequiresDocumentation returns true if the variance requires documentation.
// Currently, any non-zero variance requires documentation.
func (fh *FundHandover) RequiresDocumentation() bool {
	return fh.HasVariance()
}

// DocumentVariance adds variance documentation to the fund handover.
// This method validates that:
//   - A variance exists (non-zero)
//   - The notes meet the minimum length requirement (10 characters)
//
// Parameters:
//   - reason: The selected variance reason
//   - notes: Detailed explanation (minimum 10 characters)
//
// Returns:
//   - error: if validation fails
func (fh *FundHandover) DocumentVariance(reason VarianceReason, notes string) error {
	if !fh.HasVariance() {
		return errors.New("no variance to document")
	}

	if len(notes) < 10 {
		return errors.New("variance notes must be at least 10 characters")
	}

	fh.VarianceReason = &reason
	fh.VarianceNotes = notes
	fh.UpdatedAt = time.Now()

	return nil
}

// SetReceiver sets the receiver information for the fund handover.
// This is an extension point for future functionality where funds can be
// handed over to a specific person (e.g., manager) rather than the general fund.
//
// Parameters:
//   - receiverID: The ID of the person receiving the funds
//   - receiverName: The name of the receiver for display
func (fh *FundHandover) SetReceiver(receiverID primitive.ObjectID, receiverName string) {
	fh.ReceiverID = &receiverID
	fh.ReceiverName = receiverName
	fh.UpdatedAt = time.Now()
}

// Validate performs validation on the fund handover record.
// This method checks:
//   - Cash amount is non-negative
//   - Transfer amount is non-negative
//   - If variance exists, it must be documented
//
// Returns:
//   - error: if validation fails
func (fh *FundHandover) Validate() error {
	if fh.CashAmount < 0 {
		return errors.New("cash amount must be non-negative")
	}

	if fh.TransferAmount < 0 {
		return errors.New("transfer amount must be non-negative")
	}

	if fh.RequiresDocumentation() {
		if fh.VarianceReason == nil {
			return errors.New("variance reason is required when variance exists")
		}
		if fh.VarianceNotes == "" {
			return errors.New("variance notes are required when variance exists")
		}
		if len(fh.VarianceNotes) < 10 {
			return errors.New("variance notes must be at least 10 characters")
		}
	}

	return nil
}

// IsShortage returns true if the variance represents a cash shortage (negative variance)
func (fh *FundHandover) IsShortage() bool {
	return fh.VarianceAmount < 0
}

// IsOverage returns true if the variance represents a cash overage (positive variance)
func (fh *FundHandover) IsOverage() bool {
	return fh.VarianceAmount > 0
}

// AbsoluteVariance returns the absolute value of the variance
func (fh *FundHandover) AbsoluteVariance() float64 {
	if fh.VarianceAmount < 0 {
		return -fh.VarianceAmount
	}
	return fh.VarianceAmount
}
