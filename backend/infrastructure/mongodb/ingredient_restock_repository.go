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

type IngredientRestockRepository struct {
	collection *mongo.Collection
}

func NewIngredientRestockRepository(db *mongo.Database) *IngredientRestockRepository {
	collection := db.Collection("ingredient_restock_history")
	
	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Compound index on ingredient_id and created_at (descending)
	// For querying restock history by ingredient
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "ingredient_id", Value: 1},
			{Key: "created_at", Value: -1},
		},
	})
	
	// Index on fund_transaction_id for reverse lookup
	// Requirement 2.6: Find restock by fund transaction ID
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "fund_transaction_id", Value: 1}},
	})
	
	return &IngredientRestockRepository{
		collection: collection,
	}
}

// Create inserts a new ingredient restock record
// Requirement 2.5: Create restock record with MongoDB session support
func (r *IngredientRestockRepository) Create(ctx context.Context, record *ingredient.IngredientRestockRecord) error {
	if record.ID.IsZero() {
		record.ID = primitive.NewObjectID()
	}
	
	if record.CreatedAt.IsZero() {
		record.CreatedAt = time.Now()
	}
	
	// Validate the record before insertion
	if err := record.Validate(); err != nil {
		return err
	}
	
	_, err := r.collection.InsertOne(ctx, record)
	return err
}

// FindByIngredientID retrieves restock records for a specific ingredient with pagination
// Requirement 2.5: Find restocks for ingredient with pagination
func (r *IngredientRestockRepository) FindByIngredientID(
	ctx context.Context,
	ingredientID primitive.ObjectID,
	limit, offset int,
) ([]*ingredient.IngredientRestockRecord, error) {
	filter := bson.M{"ingredient_id": ingredientID}
	
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var records []*ingredient.IngredientRestockRecord
	if err := cursor.All(ctx, &records); err != nil {
		return nil, err
	}
	
	return records, nil
}

// FindByFundTransactionID finds a restock record by fund transaction ID
// Requirement 2.6: Find restock by fund transaction ID for bidirectional linking
func (r *IngredientRestockRepository) FindByFundTransactionID(
	ctx context.Context,
	fundTxID primitive.ObjectID,
) (*ingredient.IngredientRestockRecord, error) {
	var record ingredient.IngredientRestockRecord
	err := r.collection.FindOne(ctx, bson.M{"fund_transaction_id": fundTxID}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

// FindByID retrieves a restock record by ID
func (r *IngredientRestockRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*ingredient.IngredientRestockRecord, error) {
	var record ingredient.IngredientRestockRecord
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

// CountByIngredientID counts restock records for a specific ingredient
func (r *IngredientRestockRepository) CountByIngredientID(ctx context.Context, ingredientID primitive.ObjectID) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"ingredient_id": ingredientID})
}

// FindLatestByIngredientSince returns the most recent restock record per ingredient
// for restocks created on or after `since`. Results are sorted by created_at DESC.
func (r *IngredientRestockRepository) FindLatestByIngredientSince(
	ctx context.Context,
	since time.Time,
) ([]*ingredient.IngredientRestockRecord, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"created_at": bson.M{"$gte": since},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$ingredient_id"},
			{Key: "ingredient_id", Value: bson.M{"$first": "$ingredient_id"}},
			{Key: "quantity", Value: bson.M{"$first": "$quantity"}},
			{Key: "cost_per_unit", Value: bson.M{"$first": "$cost_per_unit"}},
			{Key: "total_cost", Value: bson.M{"$first": "$total_cost"}},
			{Key: "performed_by", Value: bson.M{"$first": "$performed_by"}},
			{Key: "performed_by_name", Value: bson.M{"$first": "$performed_by_name"}},
			{Key: "reason", Value: bson.M{"$first": "$reason"}},
			{Key: "created_at", Value: bson.M{"$first": "$created_at"}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rawResults []bson.M
	if err := cursor.All(ctx, &rawResults); err != nil {
		return nil, err
	}

	records := make([]*ingredient.IngredientRestockRecord, 0, len(rawResults))
	for _, raw := range rawResults {
		rec := &ingredient.IngredientRestockRecord{}
		if id, ok := raw["ingredient_id"].(primitive.ObjectID); ok {
			rec.IngredientID = id
		}
		if q, ok := raw["quantity"].(float64); ok {
			rec.Quantity = q
		}
		if c, ok := raw["cost_per_unit"].(float64); ok {
			rec.CostPerUnit = c
		}
		if t, ok := raw["total_cost"].(float64); ok {
			rec.TotalCost = t
		}
		if pb, ok := raw["performed_by"].(string); ok {
			rec.PerformedBy = pb
		}
		if pbn, ok := raw["performed_by_name"].(string); ok {
			rec.PerformedByName = pbn
		}
		if reason, ok := raw["reason"].(string); ok {
			rec.Reason = reason
		}
		if ca, ok := raw["created_at"].(primitive.DateTime); ok {
			rec.CreatedAt = ca.Time()
		}
		records = append(records, rec)
	}

	return records, nil
}
