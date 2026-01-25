package facility

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Facility struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Type        string             `json:"type" bson:"type"`
	Area        string             `json:"area" bson:"area"`
	Quantity    int                `json:"quantity" bson:"quantity"`
	Status      string             `json:"status" bson:"status"`
	PurchaseDate time.Time         `json:"purchase_date" bson:"purchase_date"`
	Cost        float64            `json:"cost" bson:"cost"`
	Supplier    string             `json:"supplier" bson:"supplier"`
	Notes       string             `json:"notes" bson:"notes"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type MaintenanceRecord struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FacilityID  primitive.ObjectID `json:"facility_id" bson:"facility_id"`
	Type        string             `json:"type" bson:"type"` // "scheduled", "emergency"
	Description string             `json:"description" bson:"description"`
	Cost        float64            `json:"cost" bson:"cost"`
	Vendor      string             `json:"vendor" bson:"vendor"`
	Date        time.Time          `json:"date" bson:"date"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Username    string             `json:"username" bson:"username"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type IssueReport struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FacilityID  primitive.ObjectID `json:"facility_id" bson:"facility_id"`
	Description string             `json:"description" bson:"description"`
	Severity    string             `json:"severity" bson:"severity"` // "low", "medium", "high", "critical"
	Status      string             `json:"status" bson:"status"`     // "open", "in_progress", "resolved"
	ReportedBy  primitive.ObjectID `json:"reported_by" bson:"reported_by"`
	Username    string             `json:"username" bson:"username"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	ResolvedAt  *time.Time         `json:"resolved_at,omitempty" bson:"resolved_at,omitempty"`
}

type ScheduledMaintenance struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FacilityID   primitive.ObjectID `json:"facility_id" bson:"facility_id"`
	FacilityName string             `json:"facility_name" bson:"facility_name"`
	Type         string             `json:"type" bson:"type"` // "scheduled", "preventive", "corrective"
	Description  string             `json:"description" bson:"description"`
	ScheduledDate time.Time         `json:"scheduled_date" bson:"scheduled_date"`
	Status       string             `json:"status" bson:"status"` // "pending", "in_progress", "completed", "cancelled"
	AssignedTo   string             `json:"assigned_to" bson:"assigned_to"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

type FacilityHistory struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FacilityID  primitive.ObjectID `json:"facility_id" bson:"facility_id"`
	Action      string             `json:"action" bson:"action"`
	Description string             `json:"description" bson:"description"`
	OldValue    interface{}        `json:"old_value,omitempty" bson:"old_value,omitempty"`
	NewValue    interface{}        `json:"new_value,omitempty" bson:"new_value,omitempty"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Username    string             `json:"username" bson:"username"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type FacilityHistoryFilter struct {
	FacilityID *primitive.ObjectID `json:"facility_id,omitempty"`
	Action     *string             `json:"action,omitempty"`
	DateFrom   *time.Time          `json:"date_from,omitempty"`
	DateTo     *time.Time          `json:"date_to,omitempty"`
	Limit      int                 `json:"limit,omitempty"`
	Offset     int                 `json:"offset,omitempty"`
}

type FacilityFilter struct {
	Name   *string `json:"name,omitempty"`
	Type   *string `json:"type,omitempty"`
	Area   *string `json:"area,omitempty"`
	Status *string `json:"status,omitempty"`
	Limit  int     `json:"limit,omitempty"`
	Offset int     `json:"offset,omitempty"`
}

// Constants for facility history actions
const (
	ActionCreated        = "created"
	ActionUpdated        = "updated"
	ActionMaintenance    = "maintenance"
	ActionMoved          = "moved"
	ActionDisposed       = "disposed"
	ActionStatusChange   = "status_change"
	ActionQuantityChange = "quantity_change"
)

// Constants for facility types
const (
	TypeFurniture = "Bàn ghế"
	TypeMachine   = "Máy móc"
	TypeUtensil   = "Dụng cụ"
	TypeElectric  = "Điện tử"
	TypeOther     = "Khác"
)

// Constants for facility status
const (
	StatusInUse     = "Đang sử dụng"
	StatusBroken    = "Hỏng"
	StatusRepairing = "Đang sửa"
	StatusInactive  = "Ngừng sử dụng"
	StatusDisposed  = "Thanh lý"
)

// Constants for areas
const (
	AreaDiningRoom = "Phòng khách"
	AreaKitchen    = "Bếp"
	AreaCounter    = "Quầy bar"
	AreaStorage    = "Kho"
	AreaOffice     = "Văn phòng"
	AreaOther      = "Khác"
)