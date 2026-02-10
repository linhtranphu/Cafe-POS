package mongodb

import (
	"context"
	"time"

	"cafe-pos/backend/domain/expense"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OperatingExpenseRepository struct {
	collection *mongo.Collection
}

func NewOperatingExpenseRepository(db *mongo.Database) *OperatingExpenseRepository {
	return &OperatingExpenseRepository{
		collection: db.Collection("operating_expenses"),
	}
}

// CreateIndexes creates the necessary indexes for the operating_expenses collection
func (r *OperatingExpenseRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "period_start", Value: 1},
				{Key: "period_end", Value: 1},
			},
			Options: options.Index().SetName("period_range_idx"),
		},
		{
			Keys: bson.D{
				{Key: "period_start", Value: 1},
			},
			Options: options.Index().SetName("period_start_idx"),
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// Create creates a new operating expense record
func (r *OperatingExpenseRepository) Create(ctx context.Context, operatingExpense *expense.OperatingExpense) error {
	operatingExpense.ID = primitive.NewObjectID()
	operatingExpense.CreatedAt = time.Now()
	operatingExpense.UpdatedAt = time.Now()
	operatingExpense.CalculateTotalExpenses()

	_, err := r.collection.InsertOne(ctx, operatingExpense)
	return err
}

// Update updates an existing operating expense record
func (r *OperatingExpenseRepository) Update(ctx context.Context, id primitive.ObjectID, operatingExpense *expense.OperatingExpense) error {
	operatingExpense.UpdatedAt = time.Now()
	operatingExpense.CalculateTotalExpenses()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"period_start":    operatingExpense.PeriodStart,
			"period_end":      operatingExpense.PeriodEnd,
			"staff_salary":    operatingExpense.StaffSalary,
			"rent":            operatingExpense.Rent,
			"utilities":       operatingExpense.Utilities,
			"marketing_costs": operatingExpense.MarketingCosts,
			"other_expenses":  operatingExpense.OtherExpenses,
			"total_expenses":  operatingExpense.TotalExpenses,
			"updated_at":      operatingExpense.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// FindByID finds an operating expense by ID
func (r *OperatingExpenseRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*expense.OperatingExpense, error) {
	var operatingExpense expense.OperatingExpense
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&operatingExpense)
	if err != nil {
		return nil, err
	}
	return &operatingExpense, nil
}

// FindByPeriod finds operating expenses that overlap with the given date range
func (r *OperatingExpenseRepository) FindByPeriod(ctx context.Context, startDate, endDate time.Time) ([]*expense.OperatingExpense, error) {
	// Find expenses where the period overlaps with the query range
	// Overlap condition: expense.period_start <= endDate AND expense.period_end >= startDate
	filter := bson.M{
		"period_start": bson.M{"$lte": endDate},
		"period_end":   bson.M{"$gte": startDate},
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "period_start", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var expenses []*expense.OperatingExpense
	if err := cursor.All(ctx, &expenses); err != nil {
		return nil, err
	}

	return expenses, nil
}

// FindForDate finds the operating expense that contains the given date
func (r *OperatingExpenseRepository) FindForDate(ctx context.Context, date time.Time) (*expense.OperatingExpense, error) {
	filter := bson.M{
		"period_start": bson.M{"$lte": date},
		"period_end":   bson.M{"$gte": date},
	}

	var operatingExpense expense.OperatingExpense
	err := r.collection.FindOne(ctx, filter).Decode(&operatingExpense)
	if err != nil {
		return nil, err
	}

	return &operatingExpense, nil
}

// FindAll retrieves all operating expenses with optional date range filter
func (r *OperatingExpenseRepository) FindAll(ctx context.Context, startDate, endDate *time.Time) ([]*expense.OperatingExpense, error) {
	filter := bson.M{}

	if startDate != nil && endDate != nil {
		filter = bson.M{
			"period_start": bson.M{"$lte": *endDate},
			"period_end":   bson.M{"$gte": *startDate},
		}
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "period_start", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var expenses []*expense.OperatingExpense
	if err := cursor.All(ctx, &expenses); err != nil {
		return nil, err
	}

	return expenses, nil
}

// Delete deletes an operating expense by ID
func (r *OperatingExpenseRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Upsert creates or updates an operating expense for a specific period
// If an expense already exists for the exact same period, it updates it; otherwise creates new
func (r *OperatingExpenseRepository) Upsert(ctx context.Context, operatingExpense *expense.OperatingExpense) (*expense.OperatingExpense, error) {
	// Check if an expense exists for this exact period
	filter := bson.M{
		"period_start": operatingExpense.PeriodStart,
		"period_end":   operatingExpense.PeriodEnd,
	}

	var existing expense.OperatingExpense
	err := r.collection.FindOne(ctx, filter).Decode(&existing)

	if err == mongo.ErrNoDocuments {
		// Create new
		if err := r.Create(ctx, operatingExpense); err != nil {
			return nil, err
		}
		return operatingExpense, nil
	} else if err != nil {
		return nil, err
	}

	// Update existing
	if err := r.Update(ctx, existing.ID, operatingExpense); err != nil {
		return nil, err
	}

	operatingExpense.ID = existing.ID
	operatingExpense.CreatedAt = existing.CreatedAt
	return operatingExpense, nil
}

// FindByDateRange is an alias for FindByPeriod for backward compatibility
func (r *OperatingExpenseRepository) FindByDateRange(ctx context.Context, start, end time.Time) ([]*expense.OperatingExpense, error) {
	return r.FindByPeriod(ctx, start, end)
}

// FindByDate is an alias for FindForDate for backward compatibility
func (r *OperatingExpenseRepository) FindByDate(ctx context.Context, date time.Time) (*expense.OperatingExpense, error) {
	return r.FindForDate(ctx, date)
}
