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

// JournalRepository persists double-entry journal entries
type JournalRepository struct {
	collection *mongo.Collection
}

// NewJournalRepository creates the repo and ensures indexes exist
func NewJournalRepository(db *mongo.Database) *JournalRepository {
	coll := db.Collection("journal_entries")

	ctx := context.Background()
	coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "event_id", Value: 1}}},
		{Keys: bson.D{{Key: "lines.fund_type", Value: 1}, {Key: "timestamp", Value: -1}}},
		{Keys: bson.D{{Key: "event_type", Value: 1}, {Key: "timestamp", Value: -1}}},
		{Keys: bson.D{{Key: "timestamp", Value: -1}}},
	})

	return &JournalRepository{collection: coll}
}

// Create inserts a journal entry. Must be called within a MongoDB session for atomicity.
func (r *JournalRepository) Create(ctx context.Context, entry *fund.JournalEntry) error {
	if entry.ID.IsZero() {
		entry.ID = primitive.NewObjectID()
	}
	_, err := r.collection.InsertOne(ctx, entry)
	return err
}

// JournalFilter holds query parameters for listing journal entries
type JournalFilter struct {
	FundType  *fund.FundType
	EventType *fund.JournalEventType
	FromDate  *time.Time
	ToDate    *time.Time
	Limit     int
	Offset    int
}

// List returns journal entries matching the filter, plus total count
func (r *JournalRepository) List(ctx context.Context, filter JournalFilter) ([]*fund.JournalEntry, int64, error) {
	query := bson.M{}

	if filter.FundType != nil {
		query["lines.fund_type"] = string(*filter.FundType)
	}
	if filter.EventType != nil {
		query["event_type"] = string(*filter.EventType)
	}
	if filter.FromDate != nil || filter.ToDate != nil {
		ts := bson.M{}
		if filter.FromDate != nil {
			ts["$gte"] = *filter.FromDate
		}
		if filter.ToDate != nil {
			ts["$lte"] = *filter.ToDate
		}
		query["timestamp"] = ts
	}

	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	limit := 20
	if filter.Limit > 0 {
		limit = filter.Limit
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(filter.Offset))

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var entries []*fund.JournalEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

// FindByID retrieves a journal entry by its ID
func (r *JournalRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*fund.JournalEntry, error) {
	var entry fund.JournalEntry
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

// GetFundBalance aggregates net balance for a specific fund from all journal entries
func (r *JournalRepository) GetFundBalance(ctx context.Context, fundType fund.FundType) (*fund.FundBalance, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$lines"}},
		{{Key: "$match", Value: bson.M{"lines.fund_type": string(fundType)}}},
		{{Key: "$group", Value: bson.M{
			"_id": nil,
			"cash": bson.M{"$sum": bson.M{"$cond": bson.A{
				bson.M{"$eq": bson.A{"$lines.direction", string(fund.DirectionDebit)}},
				"$lines.cash_amount",
				bson.M{"$multiply": bson.A{"$lines.cash_amount", -1}},
			}}},
			"transfer": bson.M{"$sum": bson.M{"$cond": bson.A{
				bson.M{"$eq": bson.A{"$lines.direction", string(fund.DirectionDebit)}},
				"$lines.transfer_amount",
				bson.M{"$multiply": bson.A{"$lines.transfer_amount", -1}},
			}}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Cash     float64 `bson:"cash"`
		Transfer float64 `bson:"transfer"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	balance := &fund.FundBalance{}
	if len(results) > 0 {
		balance.Cash = results[0].Cash
		balance.Transfer = results[0].Transfer
		balance.Total = balance.Cash + balance.Transfer
	}
	return balance, nil
}

// GetAllFundBalances returns balances for all real fund types plus audit accounts (shortage/overage)
func (r *JournalRepository) GetAllFundBalances(ctx context.Context) (map[fund.FundType]*fund.FundBalance, error) {
	allFunds := []fund.FundType{
		fund.FundTypeOperating,
		fund.FundTypeInventory,
		fund.FundTypeProfit,
		fund.FundTypeCashDrawer,
		fund.FundTypeWaiterFloat,
		fund.FundTypeCashShortage, // audit: accumulated shortages from waiter handovers
		fund.FundTypeCashOverage,  // audit: accumulated overages from waiter handovers
		// owner/supplier/customer omitted — pure double-entry counterpart accounts
	}

	result := make(map[fund.FundType]*fund.FundBalance)
	for _, ft := range allFunds {
		bal, err := r.GetFundBalance(ctx, ft)
		if err != nil {
			return nil, err
		}
		result[ft] = bal
	}
	return result, nil
}
