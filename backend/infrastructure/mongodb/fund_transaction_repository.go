package mongodb

import (
	"context"
	"time"

	"cafe-pos/backend/domain/fund"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FundTransactionRepository struct {
	collection *mongo.Collection
}

func NewFundTransactionRepository(db *mongo.Database) *FundTransactionRepository {
	collection := db.Collection("fund_transactions")
	
	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Index on timestamp (descending)
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "timestamp", Value: -1}},
	})
	
	// Index on type
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "type", Value: 1}},
	})
	
	// Compound index on timestamp and type
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "timestamp", Value: -1},
			{Key: "type", Value: 1},
		},
	})

	// Index on fund_type
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "fund_type", Value: 1}},
	})

	// Compound index on fund_type + timestamp
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "fund_type", Value: 1},
			{Key: "timestamp", Value: -1},
		},
	})

	// Sparse index on source_type + source_id
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "source_type", Value: 1},
			{Key: "source_id", Value: 1},
		},
		Options: options.Index().SetSparse(true),
	})

	return &FundTransactionRepository{
		collection: collection,
	}
}

// Create inserts a new fund transaction
func (r *FundTransactionRepository) Create(ctx context.Context, transaction *fund.FundTransaction) error {
	if transaction.ID.IsZero() {
		transaction.ID = primitive.NewObjectID()
	}
	
	_, err := r.collection.InsertOne(ctx, transaction)
	return err
}

// FindByID retrieves a fund transaction by ID
func (r *FundTransactionRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*fund.FundTransaction, error) {
	var transaction fund.FundTransaction
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&transaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &transaction, nil
}

// FindByDateRange retrieves fund transactions within a date range with filters
func (r *FundTransactionRepository) FindByDateRange(
	ctx context.Context,
	fromDate, toDate time.Time,
	transactionType *fund.TransactionType,
	limit, offset int,
) ([]*fund.FundTransaction, error) {
	filter := bson.M{
		"timestamp": bson.M{
			"$gte": fromDate,
			"$lte": toDate,
		},
	}
	
	if transactionType != nil {
		filter["type"] = *transactionType
	}
	
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var transactions []*fund.FundTransaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}
	
	return transactions, nil
}

// FindAll retrieves all fund transactions with pagination
func (r *FundTransactionRepository) FindAll(ctx context.Context, limit, offset int) ([]*fund.FundTransaction, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))
	
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var transactions []*fund.FundTransaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}
	
	return transactions, nil
}

// FindByFilter retrieves fund transactions matching an arbitrary BSON filter
func (r *FundTransactionRepository) FindByFilter(ctx context.Context, filter bson.M, limit, offset int) ([]*fund.FundTransaction, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*fund.FundTransaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}
	return transactions, nil
}

// CountByFilter counts fund transactions matching an arbitrary BSON filter
func (r *FundTransactionRepository) CountByFilter(ctx context.Context, filter bson.M) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}

// Count counts fund transactions matching the filter
func (r *FundTransactionRepository) Count(
	ctx context.Context,
	fromDate, toDate time.Time,
	transactionType *fund.TransactionType,
) (int64, error) {
	filter := bson.M{
		"timestamp": bson.M{
			"$gte": fromDate,
			"$lte": toDate,
		},
	}
	
	if transactionType != nil {
		filter["type"] = *transactionType
	}
	
	return r.collection.CountDocuments(ctx, filter)
}
