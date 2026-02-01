package expense

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Payment Method Constants
const (
	PaymentMethodCash = "cash"
	PaymentMethodBank = "bank"
	PaymentMethodCard = "card"
)

// Recurring Frequency Constants
const (
	FrequencyDaily     = "daily"
	FrequencyWeekly    = "weekly"
	FrequencyMonthly   = "monthly"
	FrequencyQuarterly = "quarterly"
	FrequencyYearly    = "yearly"
)

// Source Type Constants for auto-tracking
const (
	SourceTypeIngredient  = "ingredient"
	SourceTypeFacility    = "facility"
	SourceTypeMaintenance = "maintenance"
	SourceTypeManual      = "manual" // For manually created expenses
)

type Expense struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Date          time.Time          `bson:"date" json:"date"`
	CategoryID    primitive.ObjectID `bson:"category_id" json:"category_id"`
	Amount        float64            `bson:"amount" json:"amount"`
	Description   string             `bson:"description" json:"description"`
	PaymentMethod string             `bson:"payment_method" json:"payment_method"`
	Vendor        string             `bson:"vendor,omitempty" json:"vendor,omitempty"`
	Notes         string             `bson:"notes,omitempty" json:"notes,omitempty"`
	
	// Source tracking for auto-generated expenses
	SourceType string             `bson:"source_type,omitempty" json:"source_type,omitempty"`
	SourceID   primitive.ObjectID `bson:"source_id,omitempty" json:"source_id,omitempty"`
	
	// Creator tracking
	CreatedBy string    `bson:"created_by" json:"created_by"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type Category struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type RecurringExpense struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CategoryID  primitive.ObjectID `bson:"category_id" json:"category_id"`
	Amount      float64            `bson:"amount" json:"amount"`
	Description string             `bson:"description" json:"description"`
	Frequency   string             `bson:"frequency" json:"frequency"`
	StartDate   time.Time          `bson:"start_date" json:"start_date"`
	NextDue     time.Time          `bson:"next_due" json:"next_due"`
	Active      bool               `bson:"active" json:"active"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

type PrepaidExpense struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CategoryID  primitive.ObjectID `bson:"category_id" json:"category_id"`
	TotalAmount float64            `bson:"total_amount" json:"total_amount"`
	Description string             `bson:"description" json:"description"`
	StartDate   time.Time          `bson:"start_date" json:"start_date"`
	EndDate     time.Time          `bson:"end_date" json:"end_date"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}
