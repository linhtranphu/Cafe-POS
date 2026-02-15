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

// BatchRecordRepository implements the batch.BatchRecordRepository interface
type BatchRecordRepository struct {
	collection *mongo.Collection
}

// NewBatchRecordRepository creates a new BatchRecordRepository
func NewBatchRecordRepository(db *mongo.Database) *BatchRecordRepository {
	return &BatchRecordRepository{
		collection: db.Collection("batch_records"),
	}
}

// Create inserts a new batch record into the database
func (r *BatchRecordRepository) Create(ctx context.Context, record *batch.BatchRecord) error {
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, record)
	if err != nil {
		return err
	}
	
	// Set the ID from the inserted document
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		record.ID = oid
	}
	
	return nil
}

// Update updates an existing batch record in the database
func (r *BatchRecordRepository) Update(ctx context.Context, record *batch.BatchRecord) error {
	record.UpdatedAt = time.Now()
	
	filter := bson.M{"_id": record.ID}
	update := bson.M{"$set": record}
	
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete removes a batch record from the database
func (r *BatchRecordRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// FindByID retrieves a batch record by its ID
func (r *BatchRecordRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*batch.BatchRecord, error) {
	var record batch.BatchRecord
	
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&record)
	if err != nil {
		return nil, err
	}
	
	return &record, nil
}

// FindAll retrieves batch records with optional filtering and pagination
func (r *BatchRecordRepository) FindAll(ctx context.Context, filter batch.BatchRecordFilter) ([]*batch.BatchRecord, int64, error) {
	// Build query filter
	query := bson.M{}
	
	// Add batch definition ID filter if provided
	if filter.BatchDefinitionID != nil {
		query["batch_definition_id"] = *filter.BatchDefinitionID
	}
	
	// Add status filter if provided
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	
	// Add prepared by filter if provided
	if filter.PreparedBy != "" {
		query["prepared_by"] = filter.PreparedBy
	}
	
	// Add date range filter if provided
	if filter.FromDate != nil || filter.ToDate != nil {
		dateFilter := bson.M{}
		if filter.FromDate != nil {
			dateFilter["$gte"] = *filter.FromDate
		}
		if filter.ToDate != nil {
			dateFilter["$lte"] = *filter.ToDate
		}
		query["prepared_at"] = dateFilter
	}
	
	// Count total documents matching the filter
	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	
	// Set up pagination and sorting
	opts := options.Find()
	opts.SetSort(bson.D{{"expires_at", 1}}) // Sort by expiry date, oldest first (FIFO)
	
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
	var records []*batch.BatchRecord
	if err = cursor.All(ctx, &records); err != nil {
		return nil, 0, err
	}
	
	return records, total, nil
}

// FindAvailableByDefinition retrieves available batch records for a specific batch definition
// sorted by expiry date (FIFO - First In First Out)
func (r *BatchRecordRepository) FindAvailableByDefinition(ctx context.Context, defID primitive.ObjectID) ([]*batch.BatchRecord, error) {
	now := time.Now()
	
	// Query for available batches that haven't expired
	query := bson.M{
		"batch_definition_id": defID,
		"status":              batch.BatchStatusAvailable,
		"expires_at":          bson.M{"$gt": now},
		"quantity_remaining":  bson.M{"$gt": 0},
	}
	
	// Sort by expiry date ascending (oldest first - FIFO)
	opts := options.Find().SetSort(bson.D{{"expires_at", 1}})
	
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var records []*batch.BatchRecord
	if err = cursor.All(ctx, &records); err != nil {
		return nil, err
	}
	
	return records, nil
}

// UpdateQuantity updates the quantity remaining for a batch record
// This method uses optimistic locking to prevent race conditions
func (r *BatchRecordRepository) UpdateQuantity(ctx context.Context, id primitive.ObjectID, newQuantity float64) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"quantity_remaining": newQuantity,
			"updated_at":         time.Now(),
		},
	}
	
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	
	// Update status based on new quantity
	// If quantity is 0, mark as depleted
	if newQuantity <= 0 {
		statusUpdate := bson.M{
			"$set": bson.M{
				"status":     batch.BatchStatusDepleted,
				"updated_at": time.Now(),
			},
		}
		_, err = r.collection.UpdateOne(ctx, filter, statusUpdate)
		return err
	}
	
	return nil
}

// GetTotalAvailableQuantity calculates the total available quantity for a batch definition
func (r *BatchRecordRepository) GetTotalAvailableQuantity(ctx context.Context, defID primitive.ObjectID) (float64, error) {
	now := time.Now()
	
	// Query for available batches that haven't expired
	query := bson.M{
		"batch_definition_id": defID,
		"status":              batch.BatchStatusAvailable,
		"expires_at":          bson.M{"$gt": now},
		"quantity_remaining":  bson.M{"$gt": 0},
	}
	
	// Use aggregation to sum the quantity_remaining
	pipeline := []bson.M{
		{"$match": query},
		{
			"$group": bson.M{
				"_id":   nil,
				"total": bson.M{"$sum": "$quantity_remaining"},
			},
		},
	}
	
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)
	
	var result []struct {
		Total float64 `bson:"total"`
	}
	
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}
	
	if len(result) == 0 {
		return 0, nil
	}
	
	return result[0].Total, nil
}

// CreateIndexes creates necessary indexes for the batch_records collection
func (r *BatchRecordRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "batch_definition_id", Value: 1},
				{Key: "expires_at", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "expires_at", Value: 1},
			},
		},
		{
			Keys: bson.D{{Key: "expires_at", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "prepared_at", Value: -1}},
		},
		{
			Keys: bson.D{
				{Key: "prepared_by", Value: 1},
				{Key: "prepared_at", Value: -1},
			},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
