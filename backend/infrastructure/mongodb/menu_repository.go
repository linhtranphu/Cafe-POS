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
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	// Generate new ObjectID if not set
	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}
	
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	
	_, err := r.collection.InsertOne(ctx, item)
	return err
}

func (r *MenuRepository) FindAll(ctx context.Context) ([]*menu.MenuItem, error) {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	opts := options.Find().SetSort(bson.D{{"category", 1}, {"name", 1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		if IsCollectionNotFoundError(err) {
			return []*menu.MenuItem{}, nil
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*menu.MenuItem
	if err = cursor.All(ctx, &items); err != nil {
		if IsCollectionNotFoundError(err) {
			return []*menu.MenuItem{}, nil
		}
		return nil, err
	}
	
	if items == nil {
		items = []*menu.MenuItem{}
	}
	return items, nil
}

func (r *MenuRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*menu.MenuItem, error) {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	var item menu.MenuItem
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// FindByCategory finds all menu items in a specific category
func (r *MenuRepository) FindByCategory(ctx context.Context, category string) ([]*menu.MenuItem, error) {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	opts := options.Find().SetSort(bson.D{{"name", 1}})
	cursor, err := r.collection.Find(ctx, bson.M{"category": category}, opts)
	if err != nil {
		if IsCollectionNotFoundError(err) {
			return []*menu.MenuItem{}, nil
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*menu.MenuItem
	if err = cursor.All(ctx, &items); err != nil {
		if IsCollectionNotFoundError(err) {
			return []*menu.MenuItem{}, nil
		}
		return nil, err
	}
	
	if items == nil {
		items = []*menu.MenuItem{}
	}
	return items, nil
}

func (r *MenuRepository) Update(ctx context.Context, id primitive.ObjectID, item *menu.MenuItem) error {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	item.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": item})
	return err
}

func (r *MenuRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// FindByIngredientName finds all menu items that use a specific ingredient
func (r *MenuRepository) FindByIngredientName(ctx context.Context, ingredientName string) ([]*menu.MenuItem, error) {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	filter := bson.M{
		"ingredients.name": ingredientName,
	}
	
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		if IsCollectionNotFoundError(err) {
			return []*menu.MenuItem{}, nil
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*menu.MenuItem
	if err = cursor.All(ctx, &items); err != nil {
		if IsCollectionNotFoundError(err) {
			return []*menu.MenuItem{}, nil
		}
		return nil, err
	}
	
	if items == nil {
		items = []*menu.MenuItem{}
	}
	return items, nil
}

// FindByBatchDefinitionID finds all menu items that use a specific batch definition
// Searches in both single-size ingredients and multi-size variant ingredients
func (r *MenuRepository) FindByBatchDefinitionID(ctx context.Context, batchDefID primitive.ObjectID) ([]*menu.MenuItem, error) {
	ctx, cancel := WithQueryTimeout(ctx)
	defer cancel()
	
	filter := bson.M{
		"$or": []bson.M{
			// Single-size items with batch ingredient
			{"ingredients.batch_id": batchDefID},
			// Multi-size items with batch ingredient in variants
			{"variants.ingredients.batch_id": batchDefID},
		},
	}
	
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		if IsCollectionNotFoundError(err) {
			return []*menu.MenuItem{}, nil
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*menu.MenuItem
	if err = cursor.All(ctx, &items); err != nil {
		if IsCollectionNotFoundError(err) {
			return []*menu.MenuItem{}, nil
		}
		return nil, err
	}
	
	if items == nil {
		items = []*menu.MenuItem{}
	}
	return items, nil
}

// CreateIndexes creates necessary indexes for the menu_items collection
func (r *MenuRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "category", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "cost_status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "current_cost", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "ingredients.name", Value: 1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
