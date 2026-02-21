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

// PrintTemplateRepository implements the printing.PrintTemplateRepository interface
type PrintTemplateRepository struct {
	collection *mongo.Collection
}

// NewPrintTemplateRepository creates a new PrintTemplateRepository
func NewPrintTemplateRepository(db *mongo.Database) *PrintTemplateRepository {
	return &PrintTemplateRepository{
		collection: db.Collection("print_templates"),
	}
}

// Create inserts a new print template into the database
func (r *PrintTemplateRepository) Create(ctx context.Context, template *printing.PrintTemplate) error {
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, template)
	if err != nil {
		return err
	}

	// Set the ID from the inserted document
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		template.ID = oid
	}

	return nil
}

// FindByID retrieves a print template by its ID
func (r *PrintTemplateRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*printing.PrintTemplate, error) {
	var template printing.PrintTemplate

	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&template)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

// FindByType retrieves all print templates of a specific type
func (r *PrintTemplateRepository) FindByType(ctx context.Context, templateType printing.TemplateType) ([]*printing.PrintTemplate, error) {
	filter := bson.M{"type": templateType}
	opts := options.Find().SetSort(bson.D{{"created_at", 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []*printing.PrintTemplate
	if err = cursor.All(ctx, &templates); err != nil {
		return nil, err
	}

	return templates, nil
}

// FindDefault retrieves the default template for a specific type
func (r *PrintTemplateRepository) FindDefault(ctx context.Context, templateType printing.TemplateType) (*printing.PrintTemplate, error) {
	var template printing.PrintTemplate

	filter := bson.M{
		"type":       templateType,
		"is_default": true,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&template)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

// Update updates an existing print template
func (r *PrintTemplateRepository) Update(ctx context.Context, template *printing.PrintTemplate) error {
	template.UpdatedAt = time.Now()

	filter := bson.M{"_id": template.ID}
	update := bson.M{
		"$set": bson.M{
			"type":       template.Type,
			"name":       template.Name,
			"content":    template.Content,
			"is_default": template.IsDefault,
			"updated_at": template.UpdatedAt,
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

// Delete removes a print template from the database
func (r *PrintTemplateRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// CreateIndexes creates necessary indexes for the print_templates collection
func (r *PrintTemplateRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
				{Key: "is_default", Value: 1},
			},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
