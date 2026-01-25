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

type StockHistoryRepository struct {
	collection *mongo.Collection
}

func NewStockHistoryRepository(db *mongo.Database) *StockHistoryRepository {
	return &StockHistoryRepository{
		collection: db.Collection("stock_history"),
	}
}

func (r *StockHistoryRepository) Create(ctx context.Context, history *ingredient.StockHistory) error {
	history.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, history)
	return err
}

func (r *StockHistoryRepository) FindByIngredientID(ctx context.Context, ingredientID primitive.ObjectID) ([]*ingredient.StockHistory, error) {
	opts := options.Find().SetSort(bson.D{{"created_at", -1}}).SetLimit(50)
	cursor, err := r.collection.Find(ctx, bson.M{"ingredient_id": ingredientID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var histories []*ingredient.StockHistory
	if err = cursor.All(ctx, &histories); err != nil {
		return nil, err
	}
	return histories, nil
}