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

type OrderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) *OrderRepository {
	return &OrderRepository{
		collection: db.Collection("orders"),
	}
}

func (r *OrderRepository) Create(ctx context.Context, o *order.Order) error {
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, o)
	if err != nil {
		return err
	}
	o.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *OrderRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
	var o order.Order
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepository) Update(ctx context.Context, id primitive.ObjectID, o *order.Order) error {
	o.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": o})
	return err
}

func (r *OrderRepository) FindByShiftID(ctx context.Context, shiftID primitive.ObjectID) ([]*order.Order, error) {
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"shift_id": shiftID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []*order.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Order, error) {
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"waiter_id": waiterID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []*order.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) FindByStatus(ctx context.Context, status order.OrderStatus) ([]*order.Order, error) {
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"status": status}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []*order.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) FindAll(ctx context.Context) ([]*order.Order, error) {
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []*order.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*order.Order, error) {
	var o order.Order
	err := r.collection.FindOne(ctx, bson.M{"order_number": orderNumber}).Decode(&o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}
