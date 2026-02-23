package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"cafe-pos/backend/domain/cashier"
)

// FundHandoverRepository handles persistence of fund handover records
type FundHandoverRepository struct {
	collection *mongo.Collection
}

// NewFundHandoverRepository creates a new fund handover repository
func NewFundHandoverRepository(db *mongo.Database) *FundHandoverRepository {
	collection := db.Collection("fund_handovers")

	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Unique index on cashier_shift_id (one handover per shift)
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "cashier_shift_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	// Index for querying by cashier
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "cashier_id", Value: 1},
			{Key: "handover_at", Value: -1},
		},
	})

	// Index for querying by date
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "handover_at", Value: -1}},
	})

	// Index for finding handovers with variance
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "variance_amount", Value: 1}},
	})

	return &FundHandoverRepository{
		collection: collection,
	}
}

// Create creates a new fund handover record
func (r *FundHandoverRepository) Create(ctx context.Context, handover *cashier.FundHandover) error {
	if handover == nil {
		return errors.New("handover cannot be nil")
	}

	// Validate before creating
	if err := handover.Validate(); err != nil {
		return err
	}

	result, err := r.collection.InsertOne(ctx, handover)
	if err != nil {
		// Check for duplicate key error (cashier_shift_id already exists)
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("fund handover already exists for this cashier shift")
		}
		return err
	}

	handover.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds a fund handover by its ID
func (r *FundHandoverRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*cashier.FundHandover, error) {
	var handover cashier.FundHandover
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&handover)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("fund handover not found")
		}
		return nil, err
	}
	return &handover, nil
}

// FindByCashierShift finds a fund handover by cashier shift ID
func (r *FundHandoverRepository) FindByCashierShift(ctx context.Context, cashierShiftID primitive.ObjectID) (*cashier.FundHandover, error) {
	var handover cashier.FundHandover
	err := r.collection.FindOne(ctx, bson.M{"cashier_shift_id": cashierShiftID}).Decode(&handover)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Not found is not an error for this query
		}
		return nil, err
	}
	return &handover, nil
}

// FindByCashier finds fund handovers by cashier ID with pagination
func (r *FundHandoverRepository) FindByCashier(
	ctx context.Context,
	cashierID primitive.ObjectID,
	page, pageSize int,
) ([]*cashier.FundHandover, int64, error) {
	filter := bson.M{"cashier_id": cashierID}

	// Count total
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Find with pagination
	opts := options.Find().
		SetSort(bson.D{{Key: "handover_at", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var handovers []*cashier.FundHandover
	if err := cursor.All(ctx, &handovers); err != nil {
		return nil, 0, err
	}

	return handovers, total, nil
}

// FindByDateRange finds fund handovers within a date range with pagination
func (r *FundHandoverRepository) FindByDateRange(
	ctx context.Context,
	startDate, endDate time.Time,
	page, pageSize int,
) ([]*cashier.FundHandover, int64, error) {
	filter := bson.M{
		"handover_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	// Count total
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Find with pagination
	opts := options.Find().
		SetSort(bson.D{{Key: "handover_at", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var handovers []*cashier.FundHandover
	if err := cursor.All(ctx, &handovers); err != nil {
		return nil, 0, err
	}

	return handovers, total, nil
}

// FindByCashierAndDateRange finds fund handovers by cashier and date range
func (r *FundHandoverRepository) FindByCashierAndDateRange(
	ctx context.Context,
	cashierID primitive.ObjectID,
	startDate, endDate time.Time,
	page, pageSize int,
) ([]*cashier.FundHandover, int64, error) {
	filter := bson.M{
		"cashier_id": cashierID,
		"handover_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	// Count total
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Find with pagination
	opts := options.Find().
		SetSort(bson.D{{Key: "handover_at", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var handovers []*cashier.FundHandover
	if err := cursor.All(ctx, &handovers); err != nil {
		return nil, 0, err
	}

	return handovers, total, nil
}

// FindWithVariance finds fund handovers that have variance
func (r *FundHandoverRepository) FindWithVariance(
	ctx context.Context,
	page, pageSize int,
) ([]*cashier.FundHandover, int64, error) {
	filter := bson.M{
		"variance_amount": bson.M{"$ne": 0},
	}

	// Count total
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Find with pagination
	opts := options.Find().
		SetSort(bson.D{{Key: "handover_at", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var handovers []*cashier.FundHandover
	if err := cursor.All(ctx, &handovers); err != nil {
		return nil, 0, err
	}

	return handovers, total, nil
}

// Update updates a fund handover record
func (r *FundHandoverRepository) Update(ctx context.Context, id primitive.ObjectID, handover *cashier.FundHandover) error {
	if handover == nil {
		return errors.New("handover cannot be nil")
	}

	// Validate before updating
	if err := handover.Validate(); err != nil {
		return err
	}

	handover.UpdatedAt = time.Now()

	result, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": id},
		handover,
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("fund handover not found")
	}

	return nil
}

// Delete deletes a fund handover record (use with caution - for testing only)
func (r *FundHandoverRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("fund handover not found")
	}

	return nil
}

// GetStatistics returns statistics about fund handovers
func (r *FundHandoverRepository) GetStatistics(
	ctx context.Context,
	startDate, endDate time.Time,
) (map[string]interface{}, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"handover_at": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": nil,
			"total_handovers": bson.M{"$sum": 1},
			"total_cash": bson.M{"$sum": "$cash_amount"},
			"total_transfer": bson.M{"$sum": "$transfer_amount"},
			"total_variance": bson.M{"$sum": "$variance_amount"},
			"handovers_with_variance": bson.M{
				"$sum": bson.M{
					"$cond": bson.A{
						bson.M{"$ne": bson.A{"$variance_amount", 0}},
						1,
						0,
					},
				},
			},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return map[string]interface{}{
			"total_handovers":          0,
			"total_cash":               0,
			"total_transfer":           0,
			"total_variance":           0,
			"handovers_with_variance":  0,
		}, nil
	}

	return results[0], nil
}
