package settings

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ShopSettingsRepository defines the interface for shop settings persistence
type ShopSettingsRepository interface {
	// GetSettings retrieves the shop settings (there should only be one document)
	GetSettings(ctx context.Context) (*ShopSettings, error)

	// GetSettingsByID retrieves shop settings by ID
	GetSettingsByID(ctx context.Context, id primitive.ObjectID) (*ShopSettings, error)

	// CreateSettings creates new shop settings
	CreateSettings(ctx context.Context, settings *ShopSettings) error

	// UpdateSettings updates existing shop settings
	UpdateSettings(ctx context.Context, settings *ShopSettings) error
}
