package services

import (
	"context"
	"fmt"
	"time"

	"cafe-pos/backend/domain/attendance"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HoursEntryRepository interface {
	Create(ctx context.Context, e *attendance.HoursEntry) error
	FindByDateRange(ctx context.Context, from, to time.Time) ([]*attendance.HoursEntry, error)
	Update(ctx context.Context, id primitive.ObjectID, hours float64, note string) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type AttendanceService struct {
	repo HoursEntryRepository
}

func NewAttendanceService(repo HoursEntryRepository) *AttendanceService {
	return &AttendanceService{repo: repo}
}

func (s *AttendanceService) AddEntry(ctx context.Context, e *attendance.HoursEntry) error {
	if e.StaffName == "" {
		return fmt.Errorf("tên nhân viên không được để trống")
	}
	if e.Hours == 0 {
		return fmt.Errorf("số giờ không được bằng 0")
	}
	if e.Date.IsZero() {
		return fmt.Errorf("ngày không được để trống")
	}
	e.Date = time.Date(e.Date.Year(), e.Date.Month(), e.Date.Day(), 0, 0, 0, 0, time.UTC)
	return s.repo.Create(ctx, e)
}

func (s *AttendanceService) GetEntries(ctx context.Context, from, to time.Time) ([]*attendance.HoursEntry, error) {
	return s.repo.FindByDateRange(ctx, from, to)
}

func (s *AttendanceService) UpdateEntry(ctx context.Context, id primitive.ObjectID, hours float64, note string) error {
	if hours == 0 {
		return fmt.Errorf("số giờ không được bằng 0")
	}
	return s.repo.Update(ctx, id, hours, note)
}

func (s *AttendanceService) BulkAddEntries(ctx context.Context, entries []*attendance.HoursEntry) error {
	for _, e := range entries {
		if err := s.AddEntry(ctx, e); err != nil {
			return err
		}
	}
	return nil
}

func (s *AttendanceService) DeleteEntry(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}
