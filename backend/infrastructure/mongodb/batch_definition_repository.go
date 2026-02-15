package mongodb

import (
	"context"
	"time"
	"cafe-pos/backend/domain/batch"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BatchDefinitionRepository implements the batch.BatchDefinitionRepository interface
type BatchDefinitionRepository struct {
	collection *mongo.Collection
}

// NewBatchDefinitionRepository creates a new BatchDefinitionRepository
func NewBatchDefinitionRepository(db *mongo.Database) *BatchDefinitionRepository {
	return &BatchDefinitionRepository{
		collection: db.Collection("batch_definitions"),
	}
}

// Create inserts a new batch definition into the database
func (r *BatchDefinitionRepository) Create(ctx context.Context, def *batch.BatchDefinition) error {
	def.CreatedAt = time.Now()
	def.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, def)
	if err != nil {
		return err
	}
	
	// Set the ID from the inserted document
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		def.ID = oid
	}
	
	return nil
}

// Update updates an existing batch definition in the database
func (r *BatchDefinitionRepository) Update(ctx context.Context, def *batch.BatchDefinition) error {
	def.UpdatedAt = time.Now()
	
	filter := bson.M{"_id": def.ID}
	update := bson.M{"$set": def}
	
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete removes a batch definition from the database
func (r *BatchDefinitionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// FindByID retrieves a batch definition by its ID
func (r *BatchDefinitionRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*batch.BatchDefinition, error) {
	var def batch.BatchDefinition
	
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&def)
	if err != nil {
		return nil, err
	}
	
	return &def, nil
}

// FindAll retrieves batch definitions with optional filtering and pagination
func (r *BatchDefinitionRepository) FindAll(ctx context.Context, filter batch.BatchDefinitionFilter) ([]*batch.BatchDefinition, int64, error) {
	// Build query filter
	query := bson.M{}
	
	// Add search filter if provided (case-insensitive partial match)
	if filter.Search != "" {
		query["name"] = bson.M{"$regex": filter.Search, "$options": "i"}
	}
	
	// Count total documents matching the filter
	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	
	// Set up pagination
	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}}) // Sort by creation date, newest first
	
	if filter.Limit > 0 {
		opts.SetLimit(int64(filter.Limit))
		
		if filter.Page > 0 {
			skip := (filter.Page - 1) * filter.Limit
			opts.SetSkip(int64(skip))
		}
	}
	
	// Execute query
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	// Decode results
	var definitions []*batch.BatchDefinition
	if err = cursor.All(ctx, &definitions); err != nil {
		return nil, 0, err
	}
	
	return definitions, total, nil
}

// CreateIndexes creates necessary indexes for the batch_definitions collection
func (r *BatchDefinitionRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
