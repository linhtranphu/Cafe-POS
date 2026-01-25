package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"cafe-pos/backend/domain/facility"
	"cafe-pos/backend/infrastructure/mongodb"
)

type FacilityService struct {
	repo *mongodb.FacilityRepository
}

func NewFacilityService(repo *mongodb.FacilityRepository) *FacilityService {
	return &FacilityService{repo: repo}
}

func (s *FacilityService) GetAllFacilities(ctx context.Context) ([]facility.Facility, error) {
	return s.repo.GetAll(ctx)
}

func (s *FacilityService) GetFacilityByID(ctx context.Context, id primitive.ObjectID) (*facility.Facility, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *FacilityService) CreateFacility(ctx context.Context, f *facility.Facility, userID primitive.ObjectID, username string) error {
	if f.Name == "" || f.Type == "" || f.Area == "" {
		return errors.New("tên, loại và khu vực là bắt buộc")
	}
	
	if f.Quantity <= 0 {
		return errors.New("số lượng phải lớn hơn 0")
	}
	
	if f.Status == "" {
		f.Status = facility.StatusInUse
	}
	
	if f.PurchaseDate.IsZero() {
		f.PurchaseDate = time.Now()
	}
	
	err := s.repo.Create(ctx, f)
	if err != nil {
		return err
	}
	
	// Create history record
	history := &facility.FacilityHistory{
		FacilityID:  f.ID,
		Action:      facility.ActionCreated,
		Description: "Tạo mới tài sản: " + f.Name,
		NewValue:    f,
		UserID:      userID,
		Username:    username,
	}
	
	return s.repo.CreateHistory(ctx, history)
}

func (s *FacilityService) UpdateFacility(ctx context.Context, id primitive.ObjectID, f *facility.Facility, userID primitive.ObjectID, username string) error {
	// Get old facility for history
	oldFacility, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	err = s.repo.Update(ctx, id, f)
	if err != nil {
		return err
	}
	
	// Determine what changed and create appropriate history
	var description string
	var action string
	
	if oldFacility.Status != f.Status {
		action = facility.ActionStatusChange
		description = fmt.Sprintf("Thay đổi trạng thái từ '%s' thành '%s'", oldFacility.Status, f.Status)
	} else if oldFacility.Area != f.Area {
		action = facility.ActionMoved
		description = fmt.Sprintf("Di chuyển từ '%s' đến '%s'", oldFacility.Area, f.Area)
	} else if oldFacility.Quantity != f.Quantity {
		action = facility.ActionQuantityChange
		description = fmt.Sprintf("Thay đổi số lượng từ %d thành %d", oldFacility.Quantity, f.Quantity)
	} else {
		action = facility.ActionUpdated
		description = "Cập nhật thông tin tài sản: " + f.Name
	}
	
	// Create history record
	history := &facility.FacilityHistory{
		FacilityID:  id,
		Action:      action,
		Description: description,
		OldValue:    oldFacility,
		NewValue:    f,
		UserID:      userID,
		Username:    username,
	}
	
	return s.repo.CreateHistory(ctx, history)
}

func (s *FacilityService) DeleteFacility(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID, username string) error {
	// Check if facility has maintenance history
	maintenanceHistory, err := s.repo.GetMaintenanceHistory(ctx, id)
	if err != nil {
		return err
	}
	
	if len(maintenanceHistory) > 0 {
		return errors.New("không thể xóa tài sản đã có lịch sử bảo trì")
	}
	
	// Get facility for history
	f, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	
	// Create history record
	history := &facility.FacilityHistory{
		FacilityID: id,
		Action:     "deleted",
		OldValue:   f,
		UserID:     userID,
		Username:   username,
	}
	
	return s.repo.CreateHistory(ctx, history)
}

func (s *FacilityService) CreateMaintenanceRecord(ctx context.Context, record *facility.MaintenanceRecord) error {
	if record.FacilityID.IsZero() || record.Description == "" {
		return errors.New("tài sản và mô tả là bắt buộc")
	}
	
	if record.Date.IsZero() {
		record.Date = time.Now()
	}
	
	err := s.repo.CreateMaintenanceRecord(ctx, record)
	if err != nil {
		return err
	}
	
	// Create history record
	history := &facility.FacilityHistory{
		FacilityID: record.FacilityID,
		Action:     "maintenance",
		NewValue:   record,
		UserID:     record.UserID,
		Username:   record.Username,
	}
	
	return s.repo.CreateHistory(ctx, history)
}

func (s *FacilityService) GetMaintenanceHistory(ctx context.Context, facilityID primitive.ObjectID) ([]facility.MaintenanceRecord, error) {
	return s.repo.GetMaintenanceHistory(ctx, facilityID)
}

func (s *FacilityService) CreateIssueReport(ctx context.Context, report *facility.IssueReport) error {
	if report.FacilityID.IsZero() || report.Description == "" {
		return errors.New("tài sản và mô tả sự cố là bắt buộc")
	}
	
	if report.Severity == "" {
		report.Severity = "medium"
	}
	
	return s.repo.CreateIssueReport(ctx, report)
}

func (s *FacilityService) GetIssueReports(ctx context.Context) ([]facility.IssueReport, error) {
	return s.repo.GetIssueReports(ctx)
}

func (s *FacilityService) GetScheduledMaintenance(ctx context.Context) ([]facility.ScheduledMaintenance, error) {
	return s.repo.GetScheduledMaintenance(ctx)
}

func (s *FacilityService) GetNextMaintenanceDate(ctx context.Context, facilityID primitive.ObjectID) (map[string]interface{}, error) {
	return s.repo.GetNextMaintenanceDate(ctx, facilityID)
}

func (s *FacilityService) GetMaintenanceDue(ctx context.Context) ([]facility.ScheduledMaintenance, error) {
	return s.repo.GetMaintenanceDue(ctx)
}

func (s *FacilityService) GetMaintenanceStats(ctx context.Context, facilityID primitive.ObjectID) (map[string]interface{}, error) {
	return s.repo.GetMaintenanceStats(ctx, facilityID)
}

func (s *FacilityService) GetStatusHistory(ctx context.Context, facilityID primitive.ObjectID) ([]facility.FacilityHistory, error) {
	return s.repo.GetStatusHistory(ctx, facilityID)
}

func (s *FacilityService) GetFacilityHistory(ctx context.Context, facilityID primitive.ObjectID) ([]facility.FacilityHistory, error) {
	return s.repo.GetHistory(ctx, facilityID)
}

func (s *FacilityService) GetHistoryWithFilter(ctx context.Context, filter facility.FacilityHistoryFilter) ([]facility.FacilityHistory, error) {
	return s.repo.GetHistoryWithFilter(ctx, filter)
}

func (s *FacilityService) SearchFacilities(ctx context.Context, filter facility.FacilityFilter) ([]facility.Facility, int64, error) {
	facilities, err := s.repo.GetWithFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	
	total, err := s.repo.CountWithFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	
	return facilities, total, nil
}