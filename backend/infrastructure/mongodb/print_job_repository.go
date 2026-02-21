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

// PrintJobRepository implements the printing.PrintJobRepository interface
type PrintJobRepository struct {
	collection *mongo.Collection
}

// NewPrintJobRepository creates a new PrintJobRepository
func NewPrintJobRepository(db *mongo.Database) *PrintJobRepository {
	return &PrintJobRepository{
		collection: db.Collection("print_jobs"),
	}
}

// Create inserts a new print job into the database
func (r *PrintJobRepository) Create(ctx context.Context, job *printing.PrintJob) error {
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, job)
	if err != nil {
		return err
	}

	// Set the ID from the inserted document
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		job.ID = oid
	}

	return nil
}

// FindByID retrieves a print job by its ID
func (r *PrintJobRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*printing.PrintJob, error) {
	var job printing.PrintJob

	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

// FindByOrderID retrieves all print jobs for a specific order
func (r *PrintJobRepository) FindByOrderID(ctx context.Context, orderID primitive.ObjectID) ([]*printing.PrintJob, error) {
	filter := bson.M{"order_id": orderID}
	opts := options.Find().SetSort(bson.D{{"created_at", 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []*printing.PrintJob
	if err = cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

// FindPending retrieves pending print jobs ordered by created_at
func (r *PrintJobRepository) FindPending(ctx context.Context, limit int) ([]*printing.PrintJob, error) {
	filter := bson.M{"status": printing.PrintJobStatusPending}
	opts := options.Find().SetSort(bson.D{{"created_at", 1}})

	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []*printing.PrintJob
	if err = cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

// FindFailed retrieves all failed print jobs
func (r *PrintJobRepository) FindFailed(ctx context.Context) ([]*printing.PrintJob, error) {
	filter := bson.M{"status": printing.PrintJobStatusFailed}
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []*printing.PrintJob
	if err = cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

// UpdateStatus updates the status and error message of a print job
func (r *PrintJobRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status printing.PrintJobStatus, errorMsg string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"error_msg":  errorMsg,
			"updated_at": time.Now(),
		},
	}

	// If status is COMPLETED, set printed_at timestamp
	if status == printing.PrintJobStatusCompleted {
		now := time.Now()
		update = bson.M{
			"$set": bson.M{
				"status":     status,
				"error_msg":  errorMsg,
				"updated_at": time.Now(),
				"printed_at": now,
			},
		}
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

// IncrementRetry increments the retry count by 1
func (r *PrintJobRepository) IncrementRetry(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{"retry_count": 1},
		"$set": bson.M{"updated_at": time.Now()},
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

// Delete removes a print job from the database
func (r *PrintJobRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// DeleteOldCompleted deletes completed print jobs older than the specified time
func (r *PrintJobRepository) DeleteOldCompleted(ctx context.Context, olderThan time.Time) error {
	filter := bson.M{
		"status":     printing.PrintJobStatusCompleted,
		"created_at": bson.M{"$lt": olderThan},
	}

	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

// CreateIndexes creates necessary indexes for the print_jobs collection
func (r *PrintJobRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "order_id", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "created_at", Value: 1},
			},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(604800), // 7 days TTL
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
