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

type PaymentAuditRepository struct {
	collection *mongo.Collection
}

func NewPaymentAuditRepository(db *mongo.Database) *PaymentAuditRepository {
	return &PaymentAuditRepository{
		collection: db.Collection("payment_audits"),
	}
}

func (r *PaymentAuditRepository) Create(audit *cashier.PaymentAudit) error {
	audit.CreatedAt = time.Now()

	result, err := r.collection.InsertOne(context.Background(), audit)
	if err != nil {
		return err
	}

	audit.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *PaymentAuditRepository) FindByOrderID(orderID string) ([]*cashier.PaymentAudit, error) {
	opts := options.Find().SetSort(bson.D{{Key: "audited_at", Value: -1}})
	cursor, err := r.collection.Find(context.Background(), bson.M{"order_id": orderID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var audits []*cashier.PaymentAudit
	for cursor.Next(context.Background()) {
		var audit cashier.PaymentAudit
		if err := cursor.Decode(&audit); err != nil {
			return nil, err
		}
		audits = append(audits, &audit)
	}

	return audits, nil
}

func (r *PaymentAuditRepository) FindByCashierID(cashierID string) ([]*cashier.PaymentAudit, error) {
	opts := options.Find().SetSort(bson.D{{Key: "audited_at", Value: -1}})
	cursor, err := r.collection.Find(context.Background(), bson.M{"cashier_id": cashierID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var audits []*cashier.PaymentAudit
	for cursor.Next(context.Background()) {
		var audit cashier.PaymentAudit
		if err := cursor.Decode(&audit); err != nil {
			return nil, err
		}
		audits = append(audits, &audit)
	}

	return audits, nil
}

func (r *PaymentAuditRepository) FindByDateRange(start, end time.Time) ([]*cashier.PaymentAudit, error) {
	opts := options.Find().SetSort(bson.D{{Key: "audited_at", Value: -1}})
	cursor, err := r.collection.Find(context.Background(), bson.M{
		"audited_at": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var audits []*cashier.PaymentAudit
	for cursor.Next(context.Background()) {
		var audit cashier.PaymentAudit
		if err := cursor.Decode(&audit); err != nil {
			return nil, err
		}
		audits = append(audits, &audit)
	}

	return audits, nil
}

func (r *PaymentAuditRepository) FindByAction(action string) ([]*cashier.PaymentAudit, error) {
	opts := options.Find().SetSort(bson.D{{Key: "audited_at", Value: -1}})
	cursor, err := r.collection.Find(context.Background(), bson.M{"action": action}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var audits []*cashier.PaymentAudit
	for cursor.Next(context.Background()) {
		var audit cashier.PaymentAudit
		if err := cursor.Decode(&audit); err != nil {
			return nil, err
		}
		audits = append(audits, &audit)
	}

	return audits, nil
}