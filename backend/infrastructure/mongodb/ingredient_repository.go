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
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	// Only set CreatedAt if not already set (allow custom creation date)
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now()
	}
	// Always update UpdatedAt
	item.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return err
	}
	// Set the ID from the inserted document
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		item.ID = oid
	}
	return nil
}

func (r *IngredientRepository) FindAll(ctx context.Context) ([]*ingredient.Ingredient, error) {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	opts := options.Find().SetSort(bson.D{{"category", 1}, {"name", 1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		if IsCollectionNotFoundError(err) {
			return []*ingredient.Ingredient{}, nil
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*ingredient.Ingredient
	if err = cursor.All(ctx, &items); err != nil {
		if IsCollectionNotFoundError(err) {
			return []*ingredient.Ingredient{}, nil
		}
		return nil, err
	}
	
	if items == nil {
		items = []*ingredient.Ingredient{}
	}
	return items, nil
}

func (r *IngredientRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*ingredient.Ingredient, error) {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	var item ingredient.Ingredient
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *IngredientRepository) Update(ctx context.Context, id primitive.ObjectID, item *ingredient.Ingredient) error {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	item.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": item})
	return err
}

func (r *IngredientRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *IngredientRepository) FindLowStock(ctx context.Context) ([]*ingredient.Ingredient, error) {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	pipeline := []bson.M{
		{"$match": bson.M{"$expr": bson.M{"$lte": []interface{}{"$quantity", "$min_stock"}}}},
		{"$sort": bson.M{"category": 1, "name": 1}},
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		if IsCollectionNotFoundError(err) {
			return []*ingredient.Ingredient{}, nil
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*ingredient.Ingredient
	if err = cursor.All(ctx, &items); err != nil {
		if IsCollectionNotFoundError(err) {
			return []*ingredient.Ingredient{}, nil
		}
		return nil, err
	}
	
	if items == nil {
		items = []*ingredient.Ingredient{}
	}
	return items, nil
}

// Category methods
func (r *IngredientRepository) CreateCategory(ctx context.Context, cat *ingredient.IngredientCategory) error {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	categoryCollection := r.collection.Database().Collection("ingredient_categories")
	cat.CreatedAt = time.Now()
	cat.UpdatedAt = time.Now()
	result, err := categoryCollection.InsertOne(ctx, cat)
	if err != nil {
		return err
	}
	cat.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *IngredientRepository) GetCategories(ctx context.Context) ([]ingredient.IngredientCategory, error) {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	categoryCollection := r.collection.Database().Collection("ingredient_categories")
	opts := options.Find().SetSort(bson.D{{"name", 1}})
	cursor, err := categoryCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		if IsCollectionNotFoundError(err) {
			return []ingredient.IngredientCategory{}, nil
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []ingredient.IngredientCategory
	if err = cursor.All(ctx, &categories); err != nil {
		if IsCollectionNotFoundError(err) {
			return []ingredient.IngredientCategory{}, nil
		}
		return nil, err
	}
	
	if categories == nil {
		categories = []ingredient.IngredientCategory{}
	}
	return categories, nil
}

func (r *IngredientRepository) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	categoryCollection := r.collection.Database().Collection("ingredient_categories")
	
	// Check if any ingredients use this category
	var cat ingredient.IngredientCategory
	err := categoryCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&cat)
	if err != nil {
		return err
	}
	
	count, err := r.collection.CountDocuments(ctx, bson.M{"category": cat.Name})
	if err != nil {
		return err
	}
	
	if count > 0 {
		return mongo.ErrNoDocuments // Return error if category is in use
	}
	
	_, err = categoryCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
