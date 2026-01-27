package mongodb

import (
	"context"
	"time"

	"cafe-pos/backend/domain/cashier"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CashReconciliationRepository struct {
	collection *mongo.Collection
}

func NewCashReconciliationRepository(db *mongo.Database) *CashReconciliationRepository {
	return &CashReconciliationRepository{
		collection: db.Collection("cash_reconciliations"),
	}
}

func (r *CashReconciliationRepository) Create(reconciliation *cashier.CashReconciliation) error {
	reconciliation.CreatedAt = time.Now()
	reconciliation.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(context.Background(), reconciliation)
	if err != nil {
		return err
	}

	reconciliation.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *CashReconciliationRepository) FindByShiftID(shiftID string) (*cashier.CashReconciliation, error) {
	var reconciliation cashier.CashReconciliation
	err := r.collection.FindOne(context.Background(), bson.M{"shift_id": shiftID}).Decode(&reconciliation)
	if err != nil {
		return nil, err
	}
	return &reconciliation, nil
}

func (r *CashReconciliationRepository) FindByCashierID(cashierID string) ([]*cashier.CashReconciliation, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"cashier_id": cashierID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var reconciliations []*cashier.CashReconciliation
	for cursor.Next(context.Background()) {
		var reconciliation cashier.CashReconciliation
		if err := cursor.Decode(&reconciliation); err != nil {
			return nil, err
		}
		reconciliations = append(reconciliations, &reconciliation)
	}

	return reconciliations, nil
}

func (r *CashReconciliationRepository) Update(reconciliation *cashier.CashReconciliation) error {
	reconciliation.UpdatedAt = time.Now()

	objectID, err := primitive.ObjectIDFromHex(reconciliation.ID)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": reconciliation},
	)
	return err
}

func (r *CashReconciliationRepository) FindByDateRange(start, end time.Time) ([]*cashier.CashReconciliation, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{
		"reconciliation_at": bson.M{
			"$gte": start,
			"$lte": end,
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var reconciliations []*cashier.CashReconciliation
	for cursor.Next(context.Background()) {
		var reconciliation cashier.CashReconciliation
		if err := cursor.Decode(&reconciliation); err != nil {
			return nil, err
		}
		reconciliations = append(reconciliations, &reconciliation)
	}

	return reconciliations, nil
}