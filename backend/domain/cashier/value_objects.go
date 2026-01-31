package cashier

import (
	"errors"
	"time"
)

// VarianceReason represents the reason for a cash variance
type VarianceReason string

const (
	VarianceReasonCountingError    VarianceReason = "COUNTING_ERROR"
	VarianceReasonUnrecordedSale   VarianceReason = "UNRECORDED_SALE"
	VarianceReasonTheft            VarianceReason = "THEFT"
	VarianceReasonChangeError      VarianceReason = "CHANGE_ERROR"
	VarianceReasonSystemError      VarianceReason = "SYSTEM_ERROR"
	VarianceReasonOther            VarianceReason = "OTHER"
)

// Variance represents the difference between system cash and actual cash
type Variance struct {
	SystemCash float64         `json:"system_cash" bson:"system_cash"`
	ActualCash float64         `json:"actual_cash" bson:"actual_cash"`
	Amount     float64         `json:"amount" bson:"amount"`
	Reason     *VarianceReason `json:"reason,omitempty" bson:"reason,omitempty"`
	Notes      string          `json:"notes,omitempty" bson:"notes,omitempty"`
}

// NewVariance creates a new Variance value object
func NewVariance(systemCash, actualCash float64) *Variance {
	return &Variance{
		SystemCash: systemCash,
		ActualCash: actualCash,
		Amount:     actualCash - systemCash,
	}
}

// RequiresDocumentation returns true if the variance is non-zero
func (v *Variance) RequiresDocumentation() bool {
	return v.Amount != 0
}

// Document records the reason and notes for the variance
func (v *Variance) Document(reason VarianceReason, notes string) error {
	if len(notes) < 10 {
		return errors.New("variance notes must be at least 10 characters")
	}
	v.Reason = &reason
	v.Notes = notes
	return nil
}

// ResponsibilityConfirmation represents the cashier's confirmation of responsibility
type ResponsibilityConfirmation struct {
	UserID    string    `json:"user_id" bson:"user_id"`
	DeviceID  string    `json:"device_id" bson:"device_id"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

// NewResponsibilityConfirmation creates a new ResponsibilityConfirmation value object
func NewResponsibilityConfirmation(userID, deviceID string, timestamp time.Time) (*ResponsibilityConfirmation, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if deviceID == "" {
		return nil, errors.New("device ID is required")
	}
	if timestamp.IsZero() {
		return nil, errors.New("timestamp is required")
	}
	return &ResponsibilityConfirmation{
		UserID:    userID,
		DeviceID:  deviceID,
		Timestamp: timestamp,
	}, nil
}

// AuditLogEntry represents a single entry in the audit log
type AuditLogEntry struct {
	Action    string                 `json:"action" bson:"action"`
	UserID    string                 `json:"user_id" bson:"user_id"`
	DeviceID  string                 `json:"device_id" bson:"device_id"`
	Timestamp time.Time              `json:"timestamp" bson:"timestamp"`
	Data      map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"`
}

// NewAuditLogEntry creates a new AuditLogEntry value object
func NewAuditLogEntry(action, userID, deviceID string, timestamp time.Time, data map[string]interface{}) (*AuditLogEntry, error) {
	if action == "" {
		return nil, errors.New("action is required")
	}
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if deviceID == "" {
		return nil, errors.New("device ID is required")
	}
	if timestamp.IsZero() {
		return nil, errors.New("timestamp is required")
	}
	return &AuditLogEntry{
		Action:    action,
		UserID:    userID,
		DeviceID:  deviceID,
		Timestamp: timestamp,
		Data:      data,
	}, nil
}
