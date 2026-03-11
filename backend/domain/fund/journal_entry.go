package fund

import (
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Direction indicates whether a ledger line increases or decreases a fund
type Direction string

const (
	DirectionDebit  Direction = "debit"  // Fund receives money (balance increases)
	DirectionCredit Direction = "credit" // Fund sends money (balance decreases)
)

// JournalEventType identifies the business event that triggered the journal entry
type JournalEventType string

const (
	EventCashierShiftStart JournalEventType = "cashier_shift_start"
	EventCashierShiftEnd   JournalEventType = "cashier_shift_end"
	EventWaiterShiftStart  JournalEventType = "waiter_shift_start"
	EventWaiterHandover    JournalEventType = "waiter_handover"
	EventFundTransfer      JournalEventType = "fund_transfer"
	EventManagerDeposit    JournalEventType = "manager_deposit"    // Nạp tiền từ ngoài vào quỹ
	EventManagerWithdrawal JournalEventType = "manager_withdrawal" // Rút tiền từ quỹ ra ngoài
	EventExpense           JournalEventType = "expense"            // Chi tiêu từ quỹ vận hành
	EventIngredientRestock JournalEventType = "ingredient_restock" // Mua nguyên liệu từ quỹ hàng hóa
	EventFacilityPurchase  JournalEventType = "facility_purchase"  // Mua tài sản từ quỹ hàng hóa
)

// LedgerLine is one side of a double-entry journal entry
type LedgerLine struct {
	FundType       FundType    `bson:"fund_type" json:"fund_type"`
	Direction      Direction   `bson:"direction" json:"direction"`
	CashAmount     float64     `bson:"cash_amount" json:"cash_amount"`
	TransferAmount float64     `bson:"transfer_amount" json:"transfer_amount"`
	TotalAmount    float64     `bson:"total_amount" json:"total_amount"`
	BalanceBefore  FundBalance `bson:"balance_before" json:"balance_before"`
	BalanceAfter   FundBalance `bson:"balance_after" json:"balance_after"`
}

// JournalEntry records a balanced double-entry bookkeeping transaction.
// Invariant: Σ debit TotalAmount == Σ credit TotalAmount
type JournalEntry struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EventType       JournalEventType   `bson:"event_type" json:"event_type"`
	EventID         primitive.ObjectID `bson:"event_id,omitempty" json:"event_id,omitempty"`
	Lines           []LedgerLine       `bson:"lines" json:"lines"`
	Description     string             `bson:"description" json:"description"`
	PerformedBy     primitive.ObjectID `bson:"performed_by" json:"performed_by"`
	PerformedByName string             `bson:"performed_by_name" json:"performed_by_name"`
	PerformedByRole string             `bson:"performed_by_role" json:"performed_by_role"`
	Timestamp       time.Time          `bson:"timestamp" json:"timestamp"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
}

// NewJournalEntry creates an empty JournalEntry with metadata set
func NewJournalEntry(
	eventType JournalEventType,
	eventID primitive.ObjectID,
	description string,
	performedBy primitive.ObjectID,
	performedByName, performedByRole string,
) *JournalEntry {
	now := time.Now()
	return &JournalEntry{
		EventType:       eventType,
		EventID:         eventID,
		Description:     description,
		PerformedBy:     performedBy,
		PerformedByName: performedByName,
		PerformedByRole: performedByRole,
		Timestamp:       now,
		CreatedAt:       now,
	}
}

// AddLine appends a ledger line to the journal entry
func (je *JournalEntry) AddLine(
	fundType FundType,
	direction Direction,
	cashAmount, transferAmount float64,
	balanceBefore, balanceAfter FundBalance,
) {
	je.Lines = append(je.Lines, LedgerLine{
		FundType:       fundType,
		Direction:      direction,
		CashAmount:     cashAmount,
		TransferAmount: transferAmount,
		TotalAmount:    cashAmount + transferAmount,
		BalanceBefore:  balanceBefore,
		BalanceAfter:   balanceAfter,
	})
}

// Validate checks the journal entry is balanced (Σ debit == Σ credit)
func (je *JournalEntry) Validate() error {
	if len(je.Lines) < 2 {
		return errors.New("journal entry must have at least 2 lines")
	}

	var totalDebit, totalCredit float64
	for _, line := range je.Lines {
		switch line.Direction {
		case DirectionDebit:
			totalDebit += line.TotalAmount
		case DirectionCredit:
			totalCredit += line.TotalAmount
		default:
			return fmt.Errorf("invalid direction: %s", line.Direction)
		}
	}

	// Allow 1-cent floating point tolerance
	diff := totalDebit - totalCredit
	if diff < -0.01 || diff > 0.01 {
		return fmt.Errorf("unbalanced journal entry: debit=%.2f, credit=%.2f", totalDebit, totalCredit)
	}

	return nil
}
