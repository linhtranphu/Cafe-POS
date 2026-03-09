package fund

import "errors"

// ValidateTransactionType validates if the transaction type is valid
func ValidateTransactionType(t TransactionType) error {
	switch t {
	case TransactionTypeDeposit, TransactionTypeWithdrawal,
		TransactionTypeFundHandover, TransactionTypeStartingFloat:
		return nil
	default:
		return errors.New("invalid transaction type")
	}
}

// ValidateAmount validates if the amount is valid (greater than 0)
func ValidateAmount(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	return nil
}

// ValidateReason validates if the reason meets minimum length requirement
func ValidateReason(reason string) error {
	if len(reason) < 10 {
		return errors.New("reason must be at least 10 characters")
	}
	return nil
}

// CalculateTotalAmount calculates the total amount from cash and transfer
func CalculateTotalAmount(cashAmount, transferAmount float64) float64 {
	return cashAmount + transferAmount
}
