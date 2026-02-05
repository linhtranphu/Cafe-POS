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

// CashDiscrepancyRepository handles persistence for cash discrepancies
type CashDiscrepancyRepository struct {
	collection *mongo.Collection
}

// NewCashDiscrepancyRepository creates a new cash discrepancy repository
func NewCashDiscrepancyRepository(db *mongo.Database) *CashDiscrepancyRepository {
	return &CashDiscrepancyRepository{
		collection: db.Collection("cash_discrepancies"),
	}
}

// Create creates a new discrepancy record
func (r *CashDiscrepancyRepository) Create(ctx context.Context, d *handover.CashDiscrepancy) error {
	if d.ID.IsZero() {
		d.ID = primitive.NewObjectID()
	}
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, d)
	return err
}

// FindByID finds a discrepancy by ID
func (r *CashDiscrepancyRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*handover.CashDiscrepancy, error) {
	var d handover.CashDiscrepancy
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Update updates a discrepancy record
func (r *CashDiscrepancyRepository) Update(ctx context.Context, id primitive.ObjectID, d *handover.CashDiscrepancy) error {
	d.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": d},
	)
	return err
}

// FindByHandoverID finds a discrepancy by handover ID
func (r *CashDiscrepancyRepository) FindByHandoverID(ctx context.Context, handoverID primitive.ObjectID) (*handover.CashDiscrepancy, error) {
	var d handover.CashDiscrepancy
	err := r.collection.FindOne(ctx, bson.M{"handover_id": handoverID}).Decode(&d)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

// FindPendingResolution finds discrepancies pending resolution
func (r *CashDiscrepancyRepository) FindPendingResolution(ctx context.Context) ([]*handover.CashDiscrepancy, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{"resolution_status": "PENDING"},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var discrepancies []*handover.CashDiscrepancy
	if err = cursor.All(ctx, &discrepancies); err != nil {
		return nil, err
	}
	return discrepancies, nil
}

// FindRequiringApproval finds discrepancies requiring manager approval
func (r *CashDiscrepancyRepository) FindRequiringApproval(ctx context.Context) ([]*handover.CashDiscrepancy, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"requires_manager_approval": true,
			"manager_approved":          false,
		},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var discrepancies []*handover.CashDiscrepancy
	if err = cursor.All(ctx, &discrepancies); err != nil {
		return nil, err
	}
	return discrepancies, nil
}

// FindByDateRange finds discrepancies within a date range
func (r *CashDiscrepancyRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*handover.CashDiscrepancy, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"created_at": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var discrepancies []*handover.CashDiscrepancy
	if err = cursor.All(ctx, &discrepancies); err != nil {
		return nil, err
	}
	return discrepancies, nil
}

// FindByWaiterShift finds discrepancies for a waiter shift
func (r *CashDiscrepancyRepository) FindByWaiterShift(ctx context.Context, shiftID primitive.ObjectID) ([]*handover.CashDiscrepancy, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{"waiter_shift_id": shiftID},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var discrepancies []*handover.CashDiscrepancy
	if err = cursor.All(ctx, &discrepancies); err != nil {
		return nil, err
	}
	return discrepancies, nil
}

// FindByCashierShift finds discrepancies for a cashier shift
func (r *CashDiscrepancyRepository) FindByCashierShift(ctx context.Context, shiftID primitive.ObjectID) ([]*handover.CashDiscrepancy, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{"cashier_shift_id": shiftID},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var discrepancies []*handover.CashDiscrepancy
	if err = cursor.All(ctx, &discrepancies); err != nil {
		return nil, err
	}
	return discrepancies, nil
}

// CreateIndexes creates necessary indexes for the collection
func (r *CashDiscrepancyRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "handover_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "resolution_status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "requires_manager_approval", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "waiter_shift_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "cashier_shift_id", Value: 1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
