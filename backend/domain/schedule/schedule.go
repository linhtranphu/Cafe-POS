package schedule

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ShiftTemplate defines a reusable shift pattern
type ShiftTemplate struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	StartTime string             `bson:"start_time" json:"start_time"` // "07:00"
	EndTime   string             `bson:"end_time" json:"end_time"`     // "15:00"
	Roles     []string           `bson:"roles" json:"roles"`           // ["all"] or ["barista","waiter",...]
	MaxSlots  int                `bson:"max_slots" json:"max_slots"`
	IsActive  bool               `bson:"is_active" json:"is_active"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// PeriodStatus represents the state of a schedule period
type PeriodStatus string

const (
	PeriodStatusDraft     PeriodStatus = "draft"
	PeriodStatusOpen      PeriodStatus = "open"
	PeriodStatusClosed    PeriodStatus = "closed"
	PeriodStatusPublished PeriodStatus = "published"
	PeriodStatusCancelled PeriodStatus = "cancelled"
)

// SchedulePeriod represents a shift registration period (e.g., one week)
type SchedulePeriod struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	StartDate        time.Time          `bson:"start_date" json:"start_date"`
	EndDate          time.Time          `bson:"end_date" json:"end_date"`
	RegisterDeadline time.Time          `bson:"register_deadline" json:"register_deadline"`
	Status           PeriodStatus       `bson:"status" json:"status"`
	CreatedBy        primitive.ObjectID `bson:"created_by" json:"created_by"`
	CreatedByName    string             `bson:"created_by_name" json:"created_by_name"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

// ScheduleSlot represents a specific shift on a specific date within a period
type ScheduleSlot struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PeriodID     primitive.ObjectID `bson:"period_id" json:"period_id"`
	TemplateID   primitive.ObjectID `bson:"template_id" json:"template_id"`
	Date         time.Time          `bson:"date" json:"date"`
	TemplateName string             `bson:"template_name" json:"template_name"`
	StartTime    string             `bson:"start_time" json:"start_time"`
	EndTime      string             `bson:"end_time" json:"end_time"`
	Roles        []string           `bson:"roles" json:"roles"`
	MaxSlots     int                `bson:"max_slots" json:"max_slots"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

// RegistrationStatus represents the status of a shift registration
type RegistrationStatus string

const (
	RegistrationStatusConfirmed RegistrationStatus = "confirmed"
	RegistrationStatusCancelled RegistrationStatus = "cancelled"
)

// ShiftRegistration represents an employee's registration for a slot
type ShiftRegistration struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SlotID       primitive.ObjectID `bson:"slot_id" json:"slot_id"`
	PeriodID     primitive.ObjectID `bson:"period_id" json:"period_id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	UserName     string             `bson:"user_name" json:"user_name"`
	UserRole     string             `bson:"user_role" json:"user_role"`
	Status       RegistrationStatus `bson:"status" json:"status"`
	RegisteredAt time.Time          `bson:"registered_at" json:"registered_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// SlotWithRegistrations combines a slot with its registration data
type SlotWithRegistrations struct {
	ScheduleSlot
	RegisteredCount int                 `json:"registered_count"`
	Registrations   []ShiftRegistration `json:"registrations,omitempty"`
	MyRegistration  *ShiftRegistration  `json:"my_registration,omitempty"`
}

// PeriodDetail combines a period with its slots and registrations
type PeriodDetail struct {
	SchedulePeriod
	Slots []SlotWithRegistrations `json:"slots"`
}
