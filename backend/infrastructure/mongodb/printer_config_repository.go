package mongodb

import (
	"context"
	"time"

	"cafe-pos/backend/domain/printing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PrinterConfigRepository implements the printing.PrinterConfigRepository interface
type PrinterConfigRepository struct {
	collection *mongo.Collection
}

// NewPrinterConfigRepository creates a new PrinterConfigRepository
func NewPrinterConfigRepository(db *mongo.Database) *PrinterConfigRepository {
	return &PrinterConfigRepository{
		collection: db.Collection("printer_configs"),
	}
}

// Create inserts a new printer configuration into the database
func (r *PrinterConfigRepository) Create(ctx context.Context, config *printing.PrinterConfig) error {
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, config)
	if err != nil {
		return err
	}

	// Set the ID from the inserted document
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		config.ID = oid
	}

	return nil
}

// FindByID retrieves a printer configuration by its ID
func (r *PrinterConfigRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*printing.PrinterConfig, error) {
	var config printing.PrinterConfig

	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// FindAll retrieves all printer configurations
func (r *PrinterConfigRepository) FindAll(ctx context.Context) ([]*printing.PrinterConfig, error) {
	opts := options.Find().SetSort(bson.D{{"created_at", 1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var configs []*printing.PrinterConfig
	if err = cursor.All(ctx, &configs); err != nil {
		return nil, err
	}

	return configs, nil
}

// FindByType retrieves all printer configurations of a specific type
func (r *PrinterConfigRepository) FindByType(ctx context.Context, printerType printing.PrinterType) ([]*printing.PrinterConfig, error) {
	filter := bson.M{"type": printerType}
	opts := options.Find().SetSort(bson.D{{"created_at", 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var configs []*printing.PrinterConfig
	if err = cursor.All(ctx, &configs); err != nil {
		return nil, err
	}

	return configs, nil
}

// FindDefault retrieves the default printer for a specific type
func (r *PrinterConfigRepository) FindDefault(ctx context.Context, printerType printing.PrinterType) (*printing.PrinterConfig, error) {
	var config printing.PrinterConfig

	filter := bson.M{
		"type":       printerType,
		"is_default": true,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Update updates an existing printer configuration
func (r *PrinterConfigRepository) Update(ctx context.Context, config *printing.PrinterConfig) error {
	config.UpdatedAt = time.Now()

	filter := bson.M{"_id": config.ID}
	update := bson.M{
		"$set": bson.M{
			"name":            config.Name,
			"type":            config.Type,
			"connection_type": config.ConnectionType,
			"ip_address":      config.IPAddress,
			"port":            config.Port,
			"usb_path":        config.USBPath,
			"paper_width":     config.PaperWidth,
			"is_default":      config.IsDefault,
			"is_enabled":      config.IsEnabled,
			"updated_at":      config.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// Delete removes a printer configuration from the database
func (r *PrinterConfigRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// CreateIndexes creates necessary indexes for the printer_configs collection
func (r *PrinterConfigRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
				{Key: "is_default", Value: 1},
			},
		},
		{
			Keys: bson.D{{Key: "is_enabled", Value: 1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
