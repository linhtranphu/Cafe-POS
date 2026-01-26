package services

import (
	"context"
	"cafe-pos/backend/domain/expense"
	"cafe-pos/backend/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExpenseService struct {
	repo *mongodb.ExpenseRepository
}

func NewExpenseService(repo *mongodb.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

func (s *ExpenseService) CreateExpense(ctx context.Context, e *expense.Expense) error {
	return s.repo.CreateExpense(ctx, e)
}

func (s *ExpenseService) GetExpenses(ctx context.Context, filter bson.M) ([]expense.Expense, error) {
	return s.repo.GetExpenses(ctx, filter)
}

func (s *ExpenseService) UpdateExpense(ctx context.Context, id primitive.ObjectID, e *expense.Expense) error {
	return s.repo.UpdateExpense(ctx, id, e)
}

func (s *ExpenseService) DeleteExpense(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.DeleteExpense(ctx, id)
}

func (s *ExpenseService) CreateCategory(ctx context.Context, c *expense.Category) error {
	return s.repo.CreateCategory(ctx, c)
}

func (s *ExpenseService) GetCategories(ctx context.Context) ([]expense.Category, error) {
	return s.repo.GetCategories(ctx)
}

func (s *ExpenseService) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.DeleteCategory(ctx, id)
}

func (s *ExpenseService) CreateRecurring(ctx context.Context, re *expense.RecurringExpense) error {
	return s.repo.CreateRecurring(ctx, re)
}

func (s *ExpenseService) GetRecurring(ctx context.Context) ([]expense.RecurringExpense, error) {
	return s.repo.GetRecurring(ctx)
}

func (s *ExpenseService) DeleteRecurring(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.DeleteRecurring(ctx, id)
}

func (s *ExpenseService) CreatePrepaid(ctx context.Context, pe *expense.PrepaidExpense) error {
	return s.repo.CreatePrepaid(ctx, pe)
}

func (s *ExpenseService) GetPrepaid(ctx context.Context) ([]expense.PrepaidExpense, error) {
	return s.repo.GetPrepaid(ctx)
}

func (s *ExpenseService) DeletePrepaid(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.DeletePrepaid(ctx, id)
}
