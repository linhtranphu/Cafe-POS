package mongodb

import (
	"context"
	"time"
	"cafe-pos/backend/domain/expense"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExpenseRepository struct {
	expenses   *mongo.Collection
	categories *mongo.Collection
	recurring  *mongo.Collection
	prepaid    *mongo.Collection
}

func NewExpenseRepository(db *mongo.Database) *ExpenseRepository {
	return &ExpenseRepository{
		expenses:   db.Collection("expenses"),
		categories: db.Collection("expense_categories"),
		recurring:  db.Collection("recurring_expenses"),
		prepaid:    db.Collection("prepaid_expenses"),
	}
}

func (r *ExpenseRepository) CreateExpense(ctx context.Context, e *expense.Expense) error {
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
	result, err := r.expenses.InsertOne(ctx, e)
	if err == nil {
		e.ID = result.InsertedID.(primitive.ObjectID)
	}
	return err
}

func (r *ExpenseRepository) GetExpenses(ctx context.Context, filter bson.M) ([]expense.Expense, error) {
	cursor, err := r.expenses.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var expenses []expense.Expense
	if err = cursor.All(ctx, &expenses); err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *ExpenseRepository) UpdateExpense(ctx context.Context, id primitive.ObjectID, e *expense.Expense) error {
	e.UpdatedAt = time.Now()
	_, err := r.expenses.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": e})
	return err
}

func (r *ExpenseRepository) DeleteExpense(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.expenses.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *ExpenseRepository) CreateCategory(ctx context.Context, c *expense.Category) error {
	c.CreatedAt = time.Now()
	result, err := r.categories.InsertOne(ctx, c)
	if err == nil {
		c.ID = result.InsertedID.(primitive.ObjectID)
	}
	return err
}

func (r *ExpenseRepository) GetCategories(ctx context.Context) ([]expense.Category, error) {
	cursor, err := r.categories.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var categories []expense.Category
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *ExpenseRepository) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.categories.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *ExpenseRepository) CreateRecurring(ctx context.Context, re *expense.RecurringExpense) error {
	re.CreatedAt = time.Now()
	result, err := r.recurring.InsertOne(ctx, re)
	if err == nil {
		re.ID = result.InsertedID.(primitive.ObjectID)
	}
	return err
}

func (r *ExpenseRepository) GetRecurring(ctx context.Context) ([]expense.RecurringExpense, error) {
	cursor, err := r.recurring.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var recurring []expense.RecurringExpense
	if err = cursor.All(ctx, &recurring); err != nil {
		return nil, err
	}
	return recurring, nil
}

func (r *ExpenseRepository) DeleteRecurring(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.recurring.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *ExpenseRepository) CreatePrepaid(ctx context.Context, pe *expense.PrepaidExpense) error {
	pe.CreatedAt = time.Now()
	result, err := r.prepaid.InsertOne(ctx, pe)
	if err == nil {
		pe.ID = result.InsertedID.(primitive.ObjectID)
	}
	return err
}

func (r *ExpenseRepository) GetPrepaid(ctx context.Context) ([]expense.PrepaidExpense, error) {
	cursor, err := r.prepaid.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var prepaid []expense.PrepaidExpense
	if err = cursor.All(ctx, &prepaid); err != nil {
		return nil, err
	}
	return prepaid, nil
}

func (r *ExpenseRepository) DeletePrepaid(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.prepaid.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
