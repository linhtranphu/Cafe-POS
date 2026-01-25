package mongodb

import (
	"context"
	"time"
	"cafe-pos/backend/domain/menu"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MenuRepository struct {
	collection *mongo.Collection
}

func NewMenuRepository(db *mongo.Database) *MenuRepository {
	return &MenuRepository{
		collection: db.Collection("menu_items"),
	}
}

func (r *MenuRepository) Create(ctx context.Context, item *menu.MenuItem) error {
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, item)
	return err
}

func (r *MenuRepository) FindAll(ctx context.Context) ([]*menu.MenuItem, error) {
	opts := options.Find().SetSort(bson.D{{"category", 1}, {"name", 1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*menu.MenuItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MenuRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*menu.MenuItem, error) {
	var item menu.MenuItem
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *MenuRepository) Update(ctx context.Context, id primitive.ObjectID, item *menu.MenuItem) error {
	item.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": item})
	return err
}

func (r *MenuRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}