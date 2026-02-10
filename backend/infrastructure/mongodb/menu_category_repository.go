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

type MenuCategoryRepository struct {
	collection *mongo.Collection
}

func NewMenuCategoryRepository(db *mongo.Database) *MenuCategoryRepository {
	return &MenuCategoryRepository{
		collection: db.Collection("menu_categories"),
	}
}

func (r *MenuCategoryRepository) Create(ctx context.Context, category *menu.MenuCategory) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, category)
	if err != nil {
		return err
	}
	category.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *MenuCategoryRepository) FindAll(ctx context.Context) ([]*menu.MenuCategory, error) {
	opts := options.Find().SetSort(bson.D{{"name", 1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*menu.MenuCategory
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *MenuCategoryRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*menu.MenuCategory, error) {
	var category menu.MenuCategory
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *MenuCategoryRepository) FindByName(ctx context.Context, name string) (*menu.MenuCategory, error) {
	var category menu.MenuCategory
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *MenuCategoryRepository) Update(ctx context.Context, id primitive.ObjectID, category *menu.MenuCategory) error {
	category.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": category},
	)
	return err
}

func (r *MenuCategoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
