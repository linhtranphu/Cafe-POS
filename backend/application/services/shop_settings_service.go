package services

import (
	"context"
	"fmt"

	"cafe-pos/backend/domain/settings"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ShopSettingsRepository defines the interface for shop settings data access
type ShopSettingsRepository interface {
	CreateSettings(ctx context.Context, shopSettings *settings.ShopSettings) error
	GetSettingsByID(ctx context.Context, id primitive.ObjectID) (*settings.ShopSettings, error)
	FindFirst(ctx context.Context) (*settings.ShopSettings, error)
	UpdateSettings(ctx context.Context, shopSettings *settings.ShopSettings) error
	Update(ctx context.Context, id primitive.ObjectID, shopSettings *settings.ShopSettings) error
}

// ShopSettingsService handles shop settings business logic
type ShopSettingsService struct {
	repo ShopSettingsRepository
}

// NewShopSettingsService creates a new shop settings service
func NewShopSettingsService(repo ShopSettingsRepository) *ShopSettingsService {
	return &ShopSettingsService{
		repo: repo,
	}
}

// GetSettings retrieves the shop settings (there should only be one)
func (s *ShopSettingsService) GetSettings(ctx context.Context) (*settings.ShopSettings, error) {
	shopSettings, err := s.repo.FindFirst(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop settings: %w", err)
	}
	return shopSettings, nil
}

// UpdateSettings updates the shop settings
func (s *ShopSettingsService) UpdateSettings(ctx context.Context, req *UpdateSettingsRequest) (*settings.ShopSettings, error) {
	// Get existing settings
	shopSettings, err := s.repo.FindFirst(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop settings: %w", err)
	}

	// Update fields if provided
	if req.ShopName != nil {
		shopSettings.ShopName = *req.ShopName
	}

	if req.LowMarginThreshold != nil {
		if err := shopSettings.UpdateLowMarginThreshold(*req.LowMarginThreshold); err != nil {
			return nil, fmt.Errorf("failed to update low margin threshold: %w", err)
		}
	}

	// Save to database
	if err := s.repo.UpdateSettings(ctx, shopSettings); err != nil {
		return nil, fmt.Errorf("failed to update shop settings: %w", err)
	}

	return shopSettings, nil
}

// UpdateSettingsRequest represents the request to update shop settings
type UpdateSettingsRequest struct {
	ShopName           *string  `json:"shop_name"`
	LowMarginThreshold *float64 `json:"low_margin_threshold"`
}
