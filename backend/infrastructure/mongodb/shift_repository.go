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

type ShiftRepository struct {
	collection *mongo.Collection
}

func NewShiftRepository(db *mongo.Database) *ShiftRepository {
	return &ShiftRepository{
		collection: db.Collection("shifts"),
	}
}

func (r *ShiftRepository) Create(ctx context.Context, s *order.Shift) error {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, s)
	if err != nil {
		return err
	}
	s.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ShiftRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*order.Shift, error) {
	var s order.Shift
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ShiftRepository) Update(ctx context.Context, id primitive.ObjectID, s *order.Shift) error {
	s.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": s})
	return err
}

func (r *ShiftRepository) FindOpenShiftByWaiter(ctx context.Context, waiterID primitive.ObjectID) (*order.Shift, error) {
	var s order.Shift
	err := r.collection.FindOne(ctx, bson.M{
		"waiter_id": waiterID,
		"status":    order.ShiftOpen,
	}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ShiftRepository) FindOpenShifts(ctx context.Context) ([]*order.Shift, error) {
	opts := options.Find().SetSort(bson.D{{"started_at", -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"status": order.ShiftOpen}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var shifts []*order.Shift
	if err = cursor.All(ctx, &shifts); err != nil {
		return nil, err
	}
	return shifts, nil
}

func (r *ShiftRepository) FindAll(ctx context.Context) ([]*order.Shift, error) {
	opts := options.Find().SetSort(bson.D{{"started_at", -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var shifts []*order.Shift
	if err = cursor.All(ctx, &shifts); err != nil {
		return nil, err
	}
	return shifts, nil
}

func (r *ShiftRepository) FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Shift, error) {
	opts := options.Find().SetSort(bson.D{{"started_at", -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"waiter_id": waiterID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var shifts []*order.Shift
	if err = cursor.All(ctx, &shifts); err != nil {
		return nil, err
	}
	return shifts, nil
}

func (r *ShiftRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*order.Shift, error) {
	opts := options.Find().SetSort(bson.D{{"started_at", -1}})
	cursor, err := r.collection.Find(ctx, bson.M{
		"started_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var shifts []*order.Shift
	if err = cursor.All(ctx, &shifts); err != nil {
		return nil, err
	}
	return shifts, nil
}
