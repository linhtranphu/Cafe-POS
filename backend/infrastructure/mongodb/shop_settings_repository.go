package mongodb

import (
	"context"
	"time"

	"cafe-pos/backend/domain/settings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ShopSettingsRepository handles shop settings persistence
type ShopSettingsRepository struct {
	collection *mongo.Collection
}

// NewShopSettingsRepository creates a new shop settings repository
func NewShopSettingsRepository(db *mongo.Database) *ShopSettingsRepository {
	return &ShopSettingsRepository{
		collection: db.Collection("shop_settings"),
	}
}

// GetSettings retrieves the shop settings (there should only be one document)
func (r *ShopSettingsRepository) GetSettings(ctx context.Context) (*settings.ShopSettings, error) {
	var shopSettings settings.ShopSettings
	err := r.collection.FindOne(ctx, bson.M{}).Decode(&shopSettings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, settings.ErrSettingsNotFound
		}
		return nil, err
	}
	return &shopSettings, nil
}

// FindFirst retrieves the first shop settings document (alias for GetSettings)
func (r *ShopSettingsRepository) FindFirst(ctx context.Context) (*settings.ShopSettings, error) {
	return r.GetSettings(ctx)
}

// GetSettingsByID retrieves shop settings by ID
func (r *ShopSettingsRepository) GetSettingsByID(ctx context.Context, id primitive.ObjectID) (*settings.ShopSettings, error) {
	var shopSettings settings.ShopSettings
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&shopSettings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, settings.ErrSettingsNotFound
		}
		return nil, err
	}
	return &shopSettings, nil
}

// CreateSettings creates new shop settings
func (r *ShopSettingsRepository) CreateSettings(ctx context.Context, shopSettings *settings.ShopSettings) error {
	shopSettings.CreatedAt = time.Now()
	shopSettings.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, shopSettings)
	if err != nil {
		return err
	}
	
	shopSettings.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// UpdateSettings updates existing shop settings
func (r *ShopSettingsRepository) UpdateSettings(ctx context.Context, shopSettings *settings.ShopSettings) error {
	shopSettings.UpdatedAt = time.Now()
	
	filter := bson.M{"_id": shopSettings.ID}
	update := bson.M{
		"$set": bson.M{
			"shop_name":            shopSettings.ShopName,
			"shop_address":         shopSettings.ShopAddress,
			"shop_phone":           shopSettings.ShopPhone,
			"logo_url":             shopSettings.LogoURL,
			"custom_message":       shopSettings.CustomMessage,
			"show_logo":            shopSettings.ShowLogo,
			"show_address":         shopSettings.ShowAddress,
			"show_phone":           shopSettings.ShowPhone,
			"show_custom_message":  shopSettings.ShowCustomMessage,
			"low_margin_threshold": shopSettings.LowMarginThreshold,
			"auto_print_enabled":   shopSettings.AutoPrintEnabled,
			"updated_at":           shopSettings.UpdatedAt,
		},
	}
	
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	
	if result.MatchedCount == 0 {
		return settings.ErrSettingsNotFound
	}
	
	return nil
}

// Update updates shop settings by ID (alias for UpdateSettings with ID parameter)
func (r *ShopSettingsRepository) Update(ctx context.Context, id primitive.ObjectID, shopSettings *settings.ShopSettings) error {
	shopSettings.ID = id
	return r.UpdateSettings(ctx, shopSettings)
}

// UpdateLowMarginThreshold updates only the low margin threshold
func (r *ShopSettingsRepository) UpdateLowMarginThreshold(ctx context.Context, id primitive.ObjectID, threshold float64) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"low_margin_threshold": threshold,
			"updated_at":           time.Now(),
		},
	}
	
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	
	if result.MatchedCount == 0 {
		return settings.ErrSettingsNotFound
	}
	
	return nil
}

// EnsureIndexes creates necessary indexes for shop_settings collection
func (r *ShopSettingsRepository) EnsureIndexes(ctx context.Context) error {
	// No specific indexes needed for shop_settings as there's typically only one document
	return nil
}
