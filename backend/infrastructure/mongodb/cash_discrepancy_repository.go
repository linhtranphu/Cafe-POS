package mongodb

import (
	"context"
	"time"

	"cafe-pos/backend/domain/handover"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CashDiscrepancyRepository struct {
	collection *mongo.Collection
}

func NewCashDiscrepancyRepository(db *mongo.Database) *CashDiscrepancyRepository {
	return &CashDiscrepancyRepository{
		collection: db.Collection("cash_discrepancies"),
	}
}

// Create creates a new cash discrepancy record
func (r *CashDiscrepancyRepository) Create(ctx context.Context, discrepancy *handover.CashDiscrepancy) error {
	if discrepancy.CreatedAt.IsZero() {
		discrepancy.CreatedAt = time.Now()
	}
	if discrepancy.UpdatedAt.IsZero() {
		discrepancy.UpdatedAt = time.Now()
	}

	result, err := r.collection.InsertOne(ctx, discrepancy)
	if err == nil {
		discrepancy.ID = result.InsertedID.(primitive.ObjectID)
	}
	return err
}

// FindByID finds a discrepancy by its ID
func (r *CashDiscrepancyRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*handover.CashDiscrepancy, error) {
	var discrepancy handover.CashDiscrepancy
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&discrepancy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &discrepancy, nil
}

// Update updates an existing discrepancy
func (r *CashDiscrepancyRepository) Update(ctx context.Context, id primitive.ObjectID, discrepancy *handover.CashDiscrepancy) error {
	discrepancy.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": discrepancy})
	return err
}

// FindByHandoverID finds a discrepancy by handover ID
func (r *CashDiscrepancyRepository) FindByHandoverID(ctx context.Context, handoverID primitive.ObjectID) (*handover.CashDiscrepancy, error) {
	var discrepancy handover.CashDiscrepancy
	err := r.collection.FindOne(ctx, bson.M{"handover_id": handoverID}).Decode(&discrepancy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &discrepancy, nil
}

// FindPendingResolution finds all discrepancies pending resolution
func (r *CashDiscrepancyRepository) FindPendingResolution(ctx context.Context) ([]*handover.CashDiscrepancy, error) {
	filter := bson.M{
		"status": handover.DiscrepancyStatusPending,
	}

	opts := options.Find().SetSort(bson.D{{"occurred_at", 1}}) // Oldest first
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var discrepancies []*handover.CashDiscrepancy
	for cursor.Next(ctx) {
		var d handover.CashDiscrepancy
		if err := cursor.Decode(&d); err != nil {
			return nil, err
		}
		discrepancies = append(discrepancies, &d)
	}
	return discrepancies, cursor.Err()
}

// FindRequiringApproval finds all discrepancies requiring manager approval
func (r *CashDiscrepancyRepository) FindRequiringApproval(ctx context.Context) ([]*handover.CashDiscrepancy, error) {
	filter := bson.M{
		"status": handover.DiscrepancyStatusEscalated,
	}

	opts := options.Find().SetSort(bson.D{{"occurred_at", 1}}) // Oldest first
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var discrepancies []*handover.CashDiscrepancy
	for cursor.Next(ctx) {
		var d handover.CashDiscrepancy
		if err := cursor.Decode(&d); err != nil {
			return nil, err
		}
		discrepancies = append(discrepancies, &d)
	}
	return discrepancies, cursor.Err()
}

// FindByDateRange finds discrepancies within a date range
func (r *CashDiscrepancyRepository) FindByDateRange(ctx context.Context, start, end time.Time) ([]*handover.CashDiscrepancy, error) {
	filter := bson.M{
		"occurred_at": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	opts := options.Find().SetSort(bson.D{{"occurred_at", -1}}) // Newest first
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var discrepancies []*handover.CashDiscrepancy
	for cursor.Next(ctx) {
		var d handover.CashDiscrepancy
		if err := cursor.Decode(&d); err != nil {
			return nil, err
		}
		discrepancies = append(discrepancies, &d)
	}
	return discrepancies, cursor.Err()
}

// FindByWaiter finds all discrepancies for a specific waiter
func (r *CashDiscrepancyRepository) FindByWaiter(ctx context.Context, waiterID primitive.ObjectID) ([]*handover.CashDiscrepancy, error) {
	filter := bson.M{
		"waiter_id": waiterID,
	}

	opts := options.Find().SetSort(bson.D{{"occurred_at", -1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var discrepancies []*handover.CashDiscrepancy
	for cursor.Next(ctx) {
		var d handover.CashDiscrepancy
		if err := cursor.Decode(&d); err != nil {
			return nil, err
		}
		discrepancies = append(discrepancies, &d)
	}
	return discrepancies, cursor.Err()
}

// FindByCashier finds all discrepancies for a specific cashier
func (r *CashDiscrepancyRepository) FindByCashier(ctx context.Context, cashierID primitive.ObjectID) ([]*handover.CashDiscrepancy, error) {
	filter := bson.M{
		"cashier_id": cashierID,
	}

	opts := options.Find().SetSort(bson.D{{"occurred_at", -1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var discrepancies []*handover.CashDiscrepancy
	for cursor.Next(ctx) {
		var d handover.CashDiscrepancy
		if err := cursor.Decode(&d); err != nil {
			return nil, err
		}
		discrepancies = append(discrepancies, &d)
	}
	return discrepancies, cursor.Err()
}

// GetDiscrepancyStats gets statistics about discrepancies for a date range
func (r *CashDiscrepancyRepository) GetDiscrepancyStats(ctx context.Context, start, end time.Time) (*handover.DiscrepancyStats, error) {
	filter := bson.M{
		"occurred_at": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	stats := handover.NewDiscrepancyStats()
	for cursor.Next(ctx) {
		var d handover.CashDiscrepancy
		if err := cursor.Decode(&d); err != nil {
			return nil, err
		}
		stats.AddDiscrepancy(&d)
	}

	return stats, cursor.Err()
}

// CountByType counts discrepancies by type
func (r *CashDiscrepancyRepository) CountByType(ctx context.Context, discrepancyType handover.DiscrepancyType) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"type": discrepancyType})
}

// CountByStatus counts discrepancies by status
func (r *CashDiscrepancyRepository) CountByStatus(ctx context.Context, status handover.DiscrepancyStatus) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"status": status})
}

// FindTodayDiscrepancies finds all discrepancies for today
func (r *CashDiscrepancyRepository) FindTodayDiscrepancies(ctx context.Context) ([]*handover.CashDiscrepancy, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	return r.FindByDateRange(ctx, startOfDay, endOfDay)
}

// Delete deletes a discrepancy record
func (r *CashDiscrepancyRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}