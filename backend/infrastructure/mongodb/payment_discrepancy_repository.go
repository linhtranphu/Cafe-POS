package mongodb

import (
	"context"
	"time"

	"cafe-pos/backend/domain/cashier"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentDiscrepancyRepository struct {
	collection *mongo.Collection
}

func NewPaymentDiscrepancyRepository(db *mongo.Database) *PaymentDiscrepancyRepository {
	return &PaymentDiscrepancyRepository{
		collection: db.Collection("payment_discrepancies"),
	}
}

func (r *PaymentDiscrepancyRepository) Create(discrepancy *cashier.PaymentDiscrepancy) error {
	discrepancy.CreatedAt = time.Now()
	discrepancy.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(context.Background(), discrepancy)
	if err != nil {
		return err
	}

	discrepancy.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *PaymentDiscrepancyRepository) FindByOrderID(orderID string) ([]*cashier.PaymentDiscrepancy, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"order_id": orderID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var discrepancies []*cashier.PaymentDiscrepancy
	for cursor.Next(context.Background()) {
		var discrepancy cashier.PaymentDiscrepancy
		if err := cursor.Decode(&discrepancy); err != nil {
			return nil, err
		}
		discrepancies = append(discrepancies, &discrepancy)
	}

	return discrepancies, nil
}

func (r *PaymentDiscrepancyRepository) FindPendingDiscrepancies() ([]*cashier.PaymentDiscrepancy, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"status": cashier.DiscrepancyStatusPending})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var discrepancies []*cashier.PaymentDiscrepancy
	for cursor.Next(context.Background()) {
		var discrepancy cashier.PaymentDiscrepancy
		if err := cursor.Decode(&discrepancy); err != nil {
			return nil, err
		}
		discrepancies = append(discrepancies, &discrepancy)
	}

	return discrepancies, nil
}

func (r *PaymentDiscrepancyRepository) UpdateStatus(id string, status string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	if status == cashier.DiscrepancyStatusResolved {
		update["$set"].(bson.M)["resolved_at"] = time.Now()
	}

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		update,
	)
	return err
}

func (r *PaymentDiscrepancyRepository) FindByShiftID(shiftID string) ([]*cashier.PaymentDiscrepancy, error) {
	// Join with orders collection to find discrepancies by shift
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "orders",
				"localField":   "order_id",
				"foreignField": "_id",
				"as":           "order",
			},
		},
		{
			"$match": bson.M{
				"order.shift_id": shiftID,
			},
		},
	}

	cursor, err := r.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var discrepancies []*cashier.PaymentDiscrepancy
	for cursor.Next(context.Background()) {
		var discrepancy cashier.PaymentDiscrepancy
		if err := cursor.Decode(&discrepancy); err != nil {
			return nil, err
		}
		discrepancies = append(discrepancies, &discrepancy)
	}

	return discrepancies, nil
}