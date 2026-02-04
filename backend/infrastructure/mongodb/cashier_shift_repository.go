package mongodb

import (
	"context"
	"time"

	"cafe-pos/backend/domain/cashier"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CashierShiftRepository handles persistence of cashier shifts in MongoDB.
// This is separate from the regular shift repository to maintain clear separation
// between cashier shifts and waiter/barista shifts.
type CashierShiftRepository struct {
	collection *mongo.Collection
}

// NewCashierShiftRepository creates a new repository for cashier shifts.
// It uses a separate collection "cashier_shifts" to avoid confusion with waiter shifts.
func NewCashierShiftRepository(db *mongo.Database) *CashierShiftRepository {
	collection := db.Collection("cashier_shifts")
	
	// Create indexes for better query performance
	ctx := context.Background()
	
	// Index on cashier_id and start_time for finding shifts by cashier
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "cashier_id", Value: 1},
			{Key: "start_time", Value: -1},
		},
	})
	
	// Index on status for finding open/closed shifts
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "status", Value: 1}},
	})
	
	// Index on end_time for reporting
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "end_time", Value: -1}},
	})
	
	return &CashierShiftRepository{
		collection: collection,
	}
}

// Create inserts a new cashier shift into the database.
func (r *CashierShiftRepository) Create(ctx context.Context, shift *cashier.CashierShift) error {
	shift.CreatedAt = time.Now()
	shift.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, shift)
	if err != nil {
		return err
	}
	
	shift.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID retrieves a cashier shift by its ID.
func (r *CashierShiftRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*cashier.CashierShift, error) {
	var shift cashier.CashierShift
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&shift)
	if err != nil {
		return nil, err
	}
	return &shift, nil
}

// Save updates an existing cashier shift in the database.
// This is used during the shift closure workflow to persist state changes.
func (r *CashierShiftRepository) Save(ctx context.Context, shift *cashier.CashierShift) error {
	shift.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": shift.ID},
		bson.M{"$set": shift},
	)
	return err
}

// Update updates an existing cashier shift by ID.
// This method follows the same pattern as other repositories.
func (r *CashierShiftRepository) Update(ctx context.Context, id primitive.ObjectID, shift *cashier.CashierShift) error {
	shift.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": shift},
	)
	return err
}

// FindOpenByCashier finds an open cashier shift for a specific cashier.
// Returns nil if no open shift is found.
func (r *CashierShiftRepository) FindOpenByCashier(ctx context.Context, cashierID primitive.ObjectID) (*cashier.CashierShift, error) {
	var shift cashier.CashierShift
	err := r.collection.FindOne(ctx, bson.M{
		"cashier_id": cashierID,
		"status":     cashier.CashierShiftOpen,
	}).Decode(&shift)
	
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &shift, nil
}

// FindAll retrieves all cashier shifts, sorted by start time (newest first).
func (r *CashierShiftRepository) FindAll(ctx context.Context) ([]*cashier.CashierShift, error) {
	opts := options.Find().SetSort(bson.D{{Key: "start_time", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var shifts []*cashier.CashierShift
	if err = cursor.All(ctx, &shifts); err != nil {
		return nil, err
	}
	return shifts, nil
}

// FindByCashierID retrieves all shifts for a specific cashier, sorted by start time (newest first).
func (r *CashierShiftRepository) FindByCashierID(ctx context.Context, cashierID primitive.ObjectID) ([]*cashier.CashierShift, error) {
	opts := options.Find().SetSort(bson.D{{Key: "start_time", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"cashier_id": cashierID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var shifts []*cashier.CashierShift
	if err = cursor.All(ctx, &shifts); err != nil {
		return nil, err
	}
	return shifts, nil
}

// FindByStatus retrieves all cashier shifts with a specific status.
func (r *CashierShiftRepository) FindByStatus(ctx context.Context, status cashier.CashierShiftStatus) ([]*cashier.CashierShift, error) {
	opts := options.Find().SetSort(bson.D{{Key: "start_time", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"status": status}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var shifts []*cashier.CashierShift
	if err = cursor.All(ctx, &shifts); err != nil {
		return nil, err
	}
	return shifts, nil
}

// FindByDateRange retrieves cashier shifts within a date range.
func (r *CashierShiftRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*cashier.CashierShift, error) {
	opts := options.Find().SetSort(bson.D{{Key: "start_time", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{
		"start_time": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var shifts []*cashier.CashierShift
	if err = cursor.All(ctx, &shifts); err != nil {
		return nil, err
	}
	return shifts, nil
}
