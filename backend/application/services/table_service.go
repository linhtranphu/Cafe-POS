package services

import (
	"context"
	"cafe-pos/backend/domain/order"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TableRepository interface {
	Create(ctx context.Context, t *order.Table) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*order.Table, error)
	Update(ctx context.Context, id primitive.ObjectID, t *order.Table) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindAll(ctx context.Context) ([]*order.Table, error)
	FindByStatus(ctx context.Context, status order.TableStatus) ([]*order.Table, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status order.TableStatus) error
}

type TableService struct {
	tableRepo TableRepository
}

func NewTableService(tableRepo TableRepository) *TableService {
	return &TableService{
		tableRepo: tableRepo,
	}
}

func (s *TableService) CreateTable(ctx context.Context, req *order.CreateTableRequest) (*order.Table, error) {
	t := &order.Table{
		Name:     req.Name,
		Capacity: req.Capacity,
		Area:     req.Area,
		Status:   order.TableEmpty,
	}

	if err := s.tableRepo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TableService) UpdateTable(ctx context.Context, id primitive.ObjectID, req *order.UpdateTableRequest) (*order.Table, error) {
	t, err := s.tableRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		t.Name = req.Name
	}
	if req.Capacity > 0 {
		t.Capacity = req.Capacity
	}
	if req.Status != "" {
		t.Status = req.Status
	}
	if req.Area != "" {
		t.Area = req.Area
	}

	if err := s.tableRepo.Update(ctx, id, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TableService) DeleteTable(ctx context.Context, id primitive.ObjectID) error {
	return s.tableRepo.Delete(ctx, id)
}

func (s *TableService) GetAllTables(ctx context.Context) ([]*order.Table, error) {
	return s.tableRepo.FindAll(ctx)
}

func (s *TableService) GetTable(ctx context.Context, id primitive.ObjectID) (*order.Table, error) {
	return s.tableRepo.FindByID(ctx, id)
}

func (s *TableService) GetTablesByStatus(ctx context.Context, status order.TableStatus) ([]*order.Table, error) {
	return s.tableRepo.FindByStatus(ctx, status)
}

func (s *TableService) UpdateTableStatus(ctx context.Context, id primitive.ObjectID, status order.TableStatus) error {
	return s.tableRepo.UpdateStatus(ctx, id, status)
}
