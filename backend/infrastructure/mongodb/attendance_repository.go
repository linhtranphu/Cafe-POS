package mongodb

import (
	"context"
	"time"

	"cafe-pos/backend/domain/attendance"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HoursEntryRepository struct {
	col *mongo.Collection
}

func NewHoursEntryRepository(db *mongo.Database) *HoursEntryRepository {
	return &HoursEntryRepository{col: db.Collection("hours_entries")}
}

func (r *HoursEntryRepository) Create(ctx context.Context, e *attendance.HoursEntry) error {
	e.CreatedAt = time.Now()
	result, err := r.col.InsertOne(ctx, e)
	if err != nil {
		return err
	}
	e.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *HoursEntryRepository) FindByDateRange(ctx context.Context, from, to time.Time) ([]*attendance.HoursEntry, error) {
	return r.findRange(ctx, from, to, "")
}

func (r *HoursEntryRepository) FindByStaffAndDateRange(ctx context.Context, staffName string, from, to time.Time) ([]*attendance.HoursEntry, error) {
	return r.findRange(ctx, from, to, staffName)
}

func (r *HoursEntryRepository) findRange(ctx context.Context, from, to time.Time, staffName string) ([]*attendance.HoursEntry, error) {
	filter := bson.M{"date": bson.M{"$gte": from, "$lte": to}}
	if staffName != "" {
		filter["staff_name"] = staffName
	}
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: 1}, {Key: "staff_name", Value: 1}})
	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var list []*attendance.HoursEntry
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *HoursEntryRepository) Update(ctx context.Context, id primitive.ObjectID, hours float64, note string) error {
	_, err := r.col.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"hours": hours, "note": note}},
	)
	return err
}

func (r *HoursEntryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
