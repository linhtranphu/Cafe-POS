package settings

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ShopSettings represents the shop configuration and settings
type ShopSettings struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ShopName           string             `bson:"shop_name" json:"shop_name"`
	LowMarginThreshold float64            `bson:"low_margin_threshold" json:"low_margin_threshold"` // Default: 20.0
	CreatedAt          time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at" json:"updated_at"`
}

// NewShopSettings creates a new ShopSettings with default values
func NewShopSettings(shopName string) *ShopSettings {
	now := time.Now()
	return &ShopSettings{
		ShopName:           shopName,
		LowMarginThreshold: 20.0, // Default threshold
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

// UpdateLowMarginThreshold updates the low margin threshold
func (s *ShopSettings) UpdateLowMarginThreshold(threshold float64) error {
	if threshold < 0 {
		return ErrInvalidThreshold
	}
	s.LowMarginThreshold = threshold
	s.UpdatedAt = time.Now()
	return nil
}
