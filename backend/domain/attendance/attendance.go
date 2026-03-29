package attendance

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HoursEntry records an hours adjustment for a staff member on a specific day
type HoursEntry struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Date          time.Time          `bson:"date" json:"date"`
	StaffName     string             `bson:"staff_name" json:"staff_name"`
	Hours         float64            `bson:"hours" json:"hours"` // positive=add, negative=deduct
	Note          string             `bson:"note" json:"note"`
	CreatedBy     primitive.ObjectID `bson:"created_by" json:"created_by"`
	CreatedByName string             `bson:"created_by_name" json:"created_by_name"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}
