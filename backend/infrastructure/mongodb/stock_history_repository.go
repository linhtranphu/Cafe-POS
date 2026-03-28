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
	opts := options.Find().SetSort(bson.D{{"created_at", -1}}).SetLimit(20)
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

func (r *StockHistoryRepository) GetSummaryByDateRange(ctx context.Context, from, to time.Time) ([]*ingredient.IngredientStockSummary, error) {
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "created_at", Value: bson.D{{Key: "$gte", Value: from}, {Key: "$lte", Value: to}}},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$ingredient_id"},
			{Key: "consumed", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{
				bson.D{{Key: "$eq", Value: bson.A{"$type", "order"}}},
				bson.D{{Key: "$abs", Value: "$quantity"}},
				0,
			}}}}}},
			{Key: "restocked", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{
				bson.D{{Key: "$eq", Value: bson.A{"$type", "purchase"}}},
				"$quantity",
				0,
			}}}}}},
			{Key: "wasted", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{
				bson.D{{Key: "$eq", Value: bson.A{"$type", "waste"}}},
				bson.D{{Key: "$abs", Value: "$quantity"}},
				0,
			}}}}}},
			{Key: "adjusted", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{
				bson.D{{Key: "$eq", Value: bson.A{"$type", "adjustment"}}},
				"$quantity",
				0,
			}}}}}},
			{Key: "order_count", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{
				bson.D{{Key: "$eq", Value: bson.A{"$type", "order"}}},
				1,
				0,
			}}}}}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var summaries []*ingredient.IngredientStockSummary
	if err = cursor.All(ctx, &summaries); err != nil {
		return nil, err
	}
	return summaries, nil
}