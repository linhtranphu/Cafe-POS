package mongodb

import (
	"context"
	"cafe-pos/backend/domain/batch"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BatchUsageLogRepository implements the batch.BatchUsageLogRepository interface
type BatchUsageLogRepository struct {
	collection *mongo.Collection
}

// NewBatchUsageLogRepository creates a new BatchUsageLogRepository
func NewBatchUsageLogRepository(db *mongo.Database) *BatchUsageLogRepository {
	return &BatchUsageLogRepository{
		collection: db.Collection("batch_usage_logs"),
	}
}

// Create inserts a new batch usage log into the database
func (r *BatchUsageLogRepository) Create(ctx context.Context, log *batch.BatchUsageLog) error {
	result, err := r.collection.InsertOne(ctx, log)
	if err != nil {
		return err
	}
	
	// Set the ID from the inserted document
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		log.ID = oid
	}
	
	return nil
}

// FindAll retrieves batch usage logs with optional filtering and pagination
func (r *BatchUsageLogRepository) FindAll(ctx context.Context, filter batch.BatchUsageLogFilter) ([]*batch.BatchUsageLog, int64, error) {
	// Build query filter
	query := bson.M{}
	
	// Add batch record ID filter if provided
	if filter.BatchRecordID != nil {
		query["batch_record_id"] = *filter.BatchRecordID
	}
	
	// Add order ID filter if provided
	if filter.OrderID != nil {
		query["order_id"] = *filter.OrderID
	}
	
	// Add menu item ID filter if provided
	if filter.MenuItemID != nil {
		query["menu_item_id"] = *filter.MenuItemID
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
		query["used_at"] = dateFilter
	}
	
	// Count total documents matching the filter
	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	
	// Set up pagination and sorting
	opts := options.Find()
	opts.SetSort(bson.D{{"used_at", -1}}) // Sort by usage time, newest first
	
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
	var logs []*batch.BatchUsageLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, 0, err
	}
	
	return logs, total, nil
}

// CreateIndexes creates necessary indexes for the batch_usage_logs collection
func (r *BatchUsageLogRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "batch_record_id", Value: 1},
				{Key: "used_at", Value: -1},
			},
		},
		{
			Keys: bson.D{{Key: "order_id", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "menu_item_id", Value: 1},
				{Key: "used_at", Value: -1},
			},
		},
		{
			Keys: bson.D{{Key: "used_at", Value: -1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// Delete removes a batch usage log from the database
func (r *BatchUsageLogRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
