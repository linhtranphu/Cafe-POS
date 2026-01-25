package mongodb

import (
	"context"
	"time"
	"cafe-pos/backend/domain/ingredient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IngredientRepository struct {
	collection *mongo.Collection
}

func NewIngredientRepository(db *mongo.Database) *IngredientRepository {
	return &IngredientRepository{
		collection: db.Collection("ingredients"),
	}
}

func (r *IngredientRepository) Create(ctx context.Context, item *ingredient.Ingredient) error {
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, item)
	return err
}

func (r *IngredientRepository) FindAll(ctx context.Context) ([]*ingredient.Ingredient, error) {
	opts := options.Find().SetSort(bson.D{{"category", 1}, {"name", 1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*ingredient.Ingredient
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *IngredientRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*ingredient.Ingredient, error) {
	var item ingredient.Ingredient
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *IngredientRepository) Update(ctx context.Context, id primitive.ObjectID, item *ingredient.Ingredient) error {
	item.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": item})
	return err
}

func (r *IngredientRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *IngredientRepository) FindLowStock(ctx context.Context) ([]*ingredient.Ingredient, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"$expr": bson.M{"$lte": []interface{}{"$quantity", "$min_stock"}}}},
		{"$sort": bson.M{"category": 1, "name": 1}},
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*ingredient.Ingredient
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}