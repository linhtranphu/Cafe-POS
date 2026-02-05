package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"cafe-pos/backend/domain/handover"
)

// CashHandoverRepository handles persistence for cash handovers
type CashHandoverRepository struct {
	collection *mongo.Collection
}

// NewCashHandoverRepository creates a new cash handover repository
func NewCashHandoverRepository(db *mongo.Database) *CashHandoverRepository {
	return &CashHandoverRepository{
		collection: db.Collection("cash_handovers"),
	}
}

// Create creates a new cash handover record
func (r *CashHandoverRepository) Create(ctx context.Context, h *handover.CashHandover) error {
	if h.ID.IsZero() {
		h.ID = primitive.NewObjectID()
	}
	h.CreatedAt = time.Now()
	h.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, h)
	return err
}

// FindByID finds a handover by ID
func (r *CashHandoverRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*handover.CashHandover, error) {
	var h handover.CashHandover
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

// Update updates a handover record
func (r *CashHandoverRepository) Update(ctx context.Context, id primitive.ObjectID, h *handover.CashHandover) error {
	h.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": h},
	)
	return err
}

// Delete deletes a handover (only if PENDING)
func (r *CashHandoverRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(
		ctx,
		bson.M{
			"_id":    id,
			"status": handover.StatusPending,
		},
	)
	return err
}

// FindByWaiterShift finds all handovers for a waiter shift
func (r *CashHandoverRepository) FindByWaiterShift(ctx context.Context, shiftID primitive.ObjectID) ([]*handover.CashHandover, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{"waiter_shift_id": shiftID},
		options.Find().SetSort(bson.D{{Key: "handover_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var handovers []*handover.CashHandover
	if err = cursor.All(ctx, &handovers); err != nil {
		return nil, err
	}
	return handovers, nil
}

// FindPendingByWaiterShift finds pending handover for a waiter shift
func (r *CashHandoverRepository) FindPendingByWaiterShift(ctx context.Context, shiftID primitive.ObjectID) (*handover.CashHandover, error) {
	var h handover.CashHandover
	err := r.collection.FindOne(
		ctx,
		bson.M{
			"waiter_shift_id": shiftID,
			"status":          handover.StatusPending,
		},
	).Decode(&h)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &h, nil
}

// FindByCashierShift finds all handovers for a cashier shift
func (r *CashHandoverRepository) FindByCashierShift(ctx context.Context, shiftID primitive.ObjectID) ([]*handover.CashHandover, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{"cashier_shift_id": shiftID},
		options.Find().SetSort(bson.D{{Key: "handover_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var handovers []*handover.CashHandover
	if err = cursor.All(ctx, &handovers); err != nil {
		return nil, err
	}
	return handovers, nil
}

// FindPendingByCashier finds all pending handovers for a cashier
func (r *CashHandoverRepository) FindPendingByCashier(ctx context.Context, cashierID primitive.ObjectID) ([]*handover.CashHandover, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"cashier_id": cashierID,
			"status":     handover.StatusPending,
		},
		options.Find().SetSort(bson.D{{Key: "handover_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var handovers []*handover.CashHandover
	if err = cursor.All(ctx, &handovers); err != nil {
		return nil, err
	}
	return handovers, nil
}

// FindAllPending finds all pending handovers (for any cashier)
func (r *CashHandoverRepository) FindAllPending(ctx context.Context) ([]*handover.CashHandover, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{"status": handover.StatusPending},
		options.Find().SetSort(bson.D{{Key: "handover_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var handovers []*handover.CashHandover
	if err = cursor.All(ctx, &handovers); err != nil {
		return nil, err
	}
	return handovers, nil
}

// FindByDateRange finds handovers within a date range
func (r *CashHandoverRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*handover.CashHandover, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"handover_at": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		},
		options.Find().SetSort(bson.D{{Key: "handover_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var handovers []*handover.CashHandover
	if err = cursor.All(ctx, &handovers); err != nil {
		return nil, err
	}
	return handovers, nil
}

// FindTodayByCashier finds today's handovers for a cashier
func (r *CashHandoverRepository) FindTodayByCashier(ctx context.Context, cashierID primitive.ObjectID) ([]*handover.CashHandover, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"cashier_id": cashierID,
			"handover_at": bson.M{
				"$gte": startOfDay,
				"$lt":  endOfDay,
			},
		},
		options.Find().SetSort(bson.D{{Key: "handover_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var handovers []*handover.CashHandover
	if err = cursor.All(ctx, &handovers); err != nil {
		return nil, err
	}
	return handovers, nil
}

// FindWithDiscrepancies finds handovers with discrepancies
func (r *CashHandoverRepository) FindWithDiscrepancies(ctx context.Context) ([]*handover.CashHandover, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"discrepancy": bson.M{"$ne": 0},
		},
		options.Find().SetSort(bson.D{{Key: "handover_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var handovers []*handover.CashHandover
	if err = cursor.All(ctx, &handovers); err != nil {
		return nil, err
	}
	return handovers, nil
}

// FindRequiringApproval finds handovers requiring manager approval
func (r *CashHandoverRepository) FindRequiringApproval(ctx context.Context) ([]*handover.CashHandover, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"requires_approval": true,
			"status":            handover.StatusDiscrepancy,
		},
		options.Find().SetSort(bson.D{{Key: "handover_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var handovers []*handover.CashHandover
	if err = cursor.All(ctx, &handovers); err != nil {
		return nil, err
	}
	return handovers, nil
}

// CreateIndexes creates necessary indexes for the collection
func (r *CashHandoverRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "waiter_shift_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "cashier_shift_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "handover_at", Value: -1}},
		},
		{
			Keys: bson.D{
				{Key: "cashier_id", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "requires_approval", Value: 1},
				{Key: "status", Value: 1},
			},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
