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

type CashHandoverRepository struct {
	collection *mongo.Collection
}

func NewCashHandoverRepository(db *mongo.Database) *CashHandoverRepository {
	return &CashHandoverRepository{
		collection: db.Collection("cash_handovers"),
	}
}

// Create creates a new cash handover record
func (r *CashHandoverRepository) Create(ctx context.Context, handover *handover.CashHandover) error {
	if handover.CreatedAt.IsZero() {
		handover.CreatedAt = time.Now()
	}
	if handover.UpdatedAt.IsZero() {
		handover.UpdatedAt = time.Now()
	}

	result, err := r.collection.InsertOne(ctx, handover)
	if err == nil {
		handover.ID = result.InsertedID.(primitive.ObjectID)
	}
	return err
}

// FindByID finds a handover by its ID
func (r *CashHandoverRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*handover.CashHandover, error) {
	var handover handover.CashHandover
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&handover)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &handover, nil
}

// Update updates an existing handover
func (r *CashHandoverRepository) Update(ctx context.Context, id primitive.ObjectID, handover *handover.CashHandover) error {
	handover.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": handover})
	return err
}

// FindByWaiterShift finds all handovers for a specific waiter shift
func (r *CashHandoverRepository) FindByWaiterShift(ctx context.Context, shiftID primitive.ObjectID) ([]*handover.CashHandover, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"waiter_shift_id": shiftID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var handovers []*handover.CashHandover
	for cursor.Next(ctx) {
		var h handover.CashHandover
		if err := cursor.Decode(&h); err != nil {
			return nil, err
		}
		handovers = append(handovers, &h)
	}
	return handovers, cursor.Err()
}

// FindByCashierShift finds all handovers for a specific cashier shift
func (r *CashHandoverRepository) FindByCashierShift(ctx context.Context, shiftID primitive.ObjectID) ([]*handover.CashHandover, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"cashier_shift_id": shiftID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var handovers []*handover.CashHandover
	for cursor.Next(ctx) {
		var h handover.CashHandover
		if err := cursor.Decode(&h); err != nil {
			return nil, err
		}
		handovers = append(handovers, &h)
	}
	return handovers, cursor.Err()
}

// FindPendingByCashier finds all pending handovers for a specific cashier
func (r *CashHandoverRepository) FindPendingByCashier(ctx context.Context, cashierID primitive.ObjectID) ([]*handover.CashHandover, error) {
	filter := bson.M{
		"status": handover.HandoverStatusPending,
		"$or": []bson.M{
			{"cashier_id": cashierID},
			{"cashier_id": bson.M{"$exists": false}}, // Unassigned handovers
		},
	}

	opts := options.Find().SetSort(bson.D{{"requested_at", 1}}) // Oldest first
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var handovers []*handover.CashHandover
	for cursor.Next(ctx) {
		var h handover.CashHandover
		if err := cursor.Decode(&h); err != nil {
			return nil, err
		}
		handovers = append(handovers, &h)
	}
	return handovers, cursor.Err()
}

// FindByDateRange finds handovers within a date range
func (r *CashHandoverRepository) FindByDateRange(ctx context.Context, start, end time.Time) ([]*handover.CashHandover, error) {
	filter := bson.M{
		"requested_at": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	opts := options.Find().SetSort(bson.D{{"requested_at", -1}}) // Newest first
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var handovers []*handover.CashHandover
	for cursor.Next(ctx) {
		var h handover.CashHandover
		if err := cursor.Decode(&h); err != nil {
			return nil, err
		}
		handovers = append(handovers, &h)
	}
	return handovers, cursor.Err()
}

// FindWithDiscrepancies finds all handovers that have discrepancies
func (r *CashHandoverRepository) FindWithDiscrepancies(ctx context.Context) ([]*handover.CashHandover, error) {
	filter := bson.M{
		"status": handover.HandoverStatusDiscrepancy,
	}

	opts := options.Find().SetSort(bson.D{{"requested_at", -1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var handovers []*handover.CashHandover
	for cursor.Next(ctx) {
		var h handover.CashHandover
		if err := cursor.Decode(&h); err != nil {
			return nil, err
		}
		handovers = append(handovers, &h)
	}
	return handovers, cursor.Err()
}

// FindRequiringApproval finds handovers that require manager approval
func (r *CashHandoverRepository) FindRequiringApproval(ctx context.Context) ([]*handover.CashHandover, error) {
	filter := bson.M{
		"requires_manager_approval": true,
		"manager_approved":          bson.M{"$exists": false},
	}

	opts := options.Find().SetSort(bson.D{{"requested_at", 1}}) // Oldest first
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var handovers []*handover.CashHandover
	for cursor.Next(ctx) {
		var h handover.CashHandover
		if err := cursor.Decode(&h); err != nil {
			return nil, err
		}
		handovers = append(handovers, &h)
	}
	return handovers, cursor.Err()
}

// FindTodayHandovers finds all handovers for today
func (r *CashHandoverRepository) FindTodayHandovers(ctx context.Context) ([]*handover.CashHandover, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	return r.FindByDateRange(ctx, startOfDay, endOfDay)
}

// Delete deletes a handover (for cancellation)
func (r *CashHandoverRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// FindPendingByWaiterShift finds pending handover for a specific waiter shift
func (r *CashHandoverRepository) FindPendingByWaiterShift(ctx context.Context, shiftID primitive.ObjectID) (*handover.CashHandover, error) {
	filter := bson.M{
		"waiter_shift_id": shiftID,
		"status":          handover.HandoverStatusPending,
	}

	var handover handover.CashHandover
	err := r.collection.FindOne(ctx, filter).Decode(&handover)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &handover, nil
}

// CountByStatus counts handovers by status
func (r *CashHandoverRepository) CountByStatus(ctx context.Context, status handover.HandoverStatus) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"status": status})
}

// GetHandoverStats gets statistics about handovers for a date range
func (r *CashHandoverRepository) GetHandoverStats(ctx context.Context, start, end time.Time) (*HandoverStats, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"requested_at": bson.M{
					"$gte": start,
					"$lte": end,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": nil,
				"total_handovers": bson.M{"$sum": 1},
				"total_amount": bson.M{"$sum": "$requested_amount"},
				"confirmed_count": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{"$eq": []interface{}{"$status", handover.HandoverStatusConfirmed}},
							1,
							0,
						},
					},
				},
				"discrepancy_count": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{"$eq": []interface{}{"$status", handover.HandoverStatusDiscrepancy}},
							1,
							0,
						},
					},
				},
				"total_discrepancy": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{"$ne": []interface{}{"$discrepancy_amount", nil}},
							"$discrepancy_amount",
							0,
						},
					},
				},
			},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []HandoverStats
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &HandoverStats{}, nil
	}

	return &results[0], nil
}

// HandoverStats represents statistics about handovers
type HandoverStats struct {
	TotalHandovers    int     `bson:"total_handovers" json:"total_handovers"`
	TotalAmount       float64 `bson:"total_amount" json:"total_amount"`
	ConfirmedCount    int     `bson:"confirmed_count" json:"confirmed_count"`
	DiscrepancyCount  int     `bson:"discrepancy_count" json:"discrepancy_count"`
	TotalDiscrepancy  float64 `bson:"total_discrepancy" json:"total_discrepancy"`
}