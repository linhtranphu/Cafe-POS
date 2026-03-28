package fund

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TransactionType represents the type of fund transaction
type TransactionType string

const (
	TransactionTypeDeposit       TransactionType = "deposit"
	TransactionTypeWithdrawal    TransactionType = "withdrawal"
	TransactionTypeFundHandover  TransactionType = "fund_handover"
	TransactionTypeStartingFloat TransactionType = "starting_float"
)

// FundType represents which fund a transaction belongs to
type FundType string

const (
	FundTypeOperating   FundType = "operating"
	FundTypeInventory   FundType = "inventory"
	FundTypeProfit      FundType = "profit"
	FundTypeCashDrawer  FundType = "cash_drawer"
	FundTypeWaiterFloat FundType = "waiter_float" // Tiền waiter đang giữ trong ca
	// External counterpart accounts (không phải quỹ thực — dùng để cân bằng double-entry)
	FundTypeOwner        FundType = "owner"         // Chủ quán: nguồn/đích của deposit/withdrawal
	FundTypeSupplier     FundType = "supplier"      // Nhà cung cấp: đích của chi tiêu/mua nguyên liệu/tài sản
	FundTypeCustomer     FundType = "customer"      // Khách hàng: nguồn của doanh thu bán hàng
	FundTypeCashShortage FundType = "cash_shortage" // Waiter bàn giao thiếu tiền
	FundTypeCashOverage  FundType = "cash_overage"  // Waiter bàn giao thừa tiền
)

// IsExternalFund returns true for counterpart accounts that do not hold real cash.
// Balance validation is skipped for these when used as a transfer source.
func IsExternalFund(ft FundType) bool {
	switch ft {
	case FundTypeOwner, FundTypeSupplier, FundTypeCustomer, FundTypeCashShortage, FundTypeCashOverage:
		return true
	}
	return false
}

// SourceType represents the origin/source of a fund transaction
type SourceType string

const (
	SourceTypeExpense      SourceType = "expense"
	SourceTypeIngredient   SourceType = "ingredient"
	SourceTypeFacility     SourceType = "facility"
	SourceTypeHandover     SourceType = "handover"
	SourceTypeManual       SourceType = "manual"
	SourceTypeFundTransfer SourceType = "fund_transfer"
	SourceTypeWaiterStart  SourceType = "waiter_start" // Phát tiền đầu ca cho waiter
)

// FundBalance represents the balance of the fund
type FundBalance struct {
	Cash     float64 `bson:"cash" json:"cash"`
	Transfer float64 `bson:"transfer" json:"transfer"`
	Total    float64 `bson:"total" json:"total"`
}

// FundTransaction represents a deposit or withdrawal transaction
type FundTransaction struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type            TransactionType    `bson:"type" json:"type"`
	FundType        FundType           `bson:"fund_type,omitempty" json:"fund_type,omitempty"`
	CashAmount      float64            `bson:"cash_amount" json:"cash_amount"`
	TransferAmount  float64            `bson:"transfer_amount" json:"transfer_amount"`
	TotalAmount     float64            `bson:"total_amount" json:"total_amount"`
	Reason          string             `bson:"reason" json:"reason"`
	Description     string             `bson:"description,omitempty" json:"description,omitempty"`
	SourceType      SourceType         `bson:"source_type,omitempty" json:"source_type,omitempty"`
	SourceID        primitive.ObjectID `bson:"source_id,omitempty" json:"source_id,omitempty"`
	LinkedTxID      primitive.ObjectID `bson:"linked_tx_id,omitempty" json:"linked_tx_id,omitempty"`
	PerformedBy     primitive.ObjectID `bson:"performed_by" json:"performed_by"`
	PerformedByName string             `bson:"performed_by_name" json:"performed_by_name"`
	PerformedByRole string             `bson:"performed_by_role" json:"performed_by_role"`
	Timestamp       time.Time          `bson:"timestamp" json:"timestamp"`
	BalanceBefore   *FundBalance       `bson:"balance_before,omitempty" json:"balance_before,omitempty"`
	BalanceAfter    *FundBalance       `bson:"balance_after,omitempty" json:"balance_after,omitempty"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

// NewFundTransaction creates a new fund transaction
func NewFundTransaction(
	transactionType TransactionType,
	cashAmount, transferAmount float64,
	reason string,
	performedBy primitive.ObjectID,
	performedByName, performedByRole string,
) (*FundTransaction, error) {
	now := time.Now()

	transaction := &FundTransaction{
		Type:            transactionType,
		CashAmount:      cashAmount,
		TransferAmount:  transferAmount,
		TotalAmount:     cashAmount + transferAmount,
		Reason:          reason,
		PerformedBy:     performedBy,
		PerformedByName: performedByName,
		PerformedByRole: performedByRole,
		Timestamp:       now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := transaction.Validate(); err != nil {
		return nil, err
	}

	return transaction, nil
}

// Validate validates the fund transaction
func (ft *FundTransaction) Validate() error {
	validTypes := map[TransactionType]bool{
		TransactionTypeDeposit:       true,
		TransactionTypeWithdrawal:    true,
		TransactionTypeFundHandover:  true,
		TransactionTypeStartingFloat: true,
	}
	if !validTypes[ft.Type] {
		return errors.New("invalid transaction type")
	}

	if ft.CashAmount < 0 || ft.TransferAmount < 0 {
		return errors.New("amounts cannot be negative")
	}

	if ft.TotalAmount <= 0 {
		return errors.New("total amount must be greater than 0")
	}

	if len(ft.Reason) < 10 {
		return errors.New("reason must be at least 10 characters")
	}

	if ft.PerformedBy.IsZero() {
		return errors.New("performed_by is required")
	}

	if ft.PerformedByName == "" {
		return errors.New("performed_by_name is required")
	}

	if ft.PerformedByRole == "" {
		return errors.New("performed_by_role is required")
	}

	return nil
}

// SetBalances sets the before and after balances for audit purposes
func (ft *FundTransaction) SetBalances(before, after *FundBalance) {
	ft.BalanceBefore = before
	ft.BalanceAfter = after
}

// SetFundType sets which fund this transaction belongs to
func (ft *FundTransaction) SetFundType(fundType FundType) {
	ft.FundType = fundType
}

// SetSource sets the source reference for this transaction
func (ft *FundTransaction) SetSource(sourceType SourceType, sourceID primitive.ObjectID) error {
	ft.SourceType = sourceType
	ft.SourceID = sourceID
	return nil
}

// LinkTo sets the linked transaction ID (used for fund_transfer pairs)
func (ft *FundTransaction) LinkTo(linkedTxID primitive.ObjectID) {
	ft.LinkedTxID = linkedTxID
}
