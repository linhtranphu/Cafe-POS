package mongodb

import (
	"context"
	"time"
	"cafe-pos/backend/domain/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TableRepository struct {
	collection *mongo.Collection
}

func NewTableRepository(db *mongo.Database) *TableRepository {
	return &TableRepository{
		collection: db.Collection("tables"),
	}
}

func (r *TableRepository) Create(ctx context.Context, t *order.Table) error {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, t)
	if err != nil {
		return err
	}
	t.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *TableRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*order.Table, error) {
	var t order.Table
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TableRepository) Update(ctx context.Context, id primitive.ObjectID, t *order.Table) error {
	t.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": t})
	return err
}

func (r *TableRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *TableRepository) FindAll(ctx context.Context) ([]*order.Table, error) {
	opts := options.Find().SetSort(bson.D{{"name", 1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tables []*order.Table
	if err = cursor.All(ctx, &tables); err != nil {
		return nil, err
	}
	return tables, nil
}

func (r *TableRepository) FindByStatus(ctx context.Context, status order.TableStatus) ([]*order.Table, error) {
	opts := options.Find().SetSort(bson.D{{"name", 1}})
	cursor, err := r.collection.Find(ctx, bson.M{"status": status}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tables []*order.Table
	if err = cursor.All(ctx, &tables); err != nil {
		return nil, err
	}
	return tables, nil
}

func (r *TableRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status order.TableStatus) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}},
	)
	return err
}
