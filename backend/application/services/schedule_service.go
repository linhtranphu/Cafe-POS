package services

import (
	"context"
	"fmt"
	"time"

	"cafe-pos/backend/domain/schedule"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ─── Repository interfaces ────────────────────────────────────────────────────

type ShiftTemplateRepository interface {
	Create(ctx context.Context, t *schedule.ShiftTemplate) error
	FindAll(ctx context.Context) ([]*schedule.ShiftTemplate, error)
	FindActive(ctx context.Context) ([]*schedule.ShiftTemplate, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*schedule.ShiftTemplate, error)
	Update(ctx context.Context, id primitive.ObjectID, t *schedule.ShiftTemplate) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type SchedulePeriodRepository interface {
	Create(ctx context.Context, p *schedule.SchedulePeriod) error
	FindAll(ctx context.Context) ([]*schedule.SchedulePeriod, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*schedule.SchedulePeriod, error)
	FindOpenOrPublished(ctx context.Context) (*schedule.SchedulePeriod, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status schedule.PeriodStatus) error
}

type ScheduleSlotRepository interface {
	Create(ctx context.Context, s *schedule.ScheduleSlot) error
	FindByPeriod(ctx context.Context, periodID primitive.ObjectID) ([]*schedule.ScheduleSlot, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*schedule.ScheduleSlot, error)
	FindByDate(ctx context.Context, date time.Time) ([]*schedule.ScheduleSlot, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type ShiftRegistrationRepository interface {
	Create(ctx context.Context, reg *schedule.ShiftRegistration) error
	FindBySlot(ctx context.Context, slotID primitive.ObjectID) ([]*schedule.ShiftRegistration, error)
	FindByUserAndPeriod(ctx context.Context, userID, periodID primitive.ObjectID) ([]*schedule.ShiftRegistration, error)
	FindByUserAndSlot(ctx context.Context, userID, slotID primitive.ObjectID) (*schedule.ShiftRegistration, error)
	CountBySlot(ctx context.Context, slotID primitive.ObjectID) (int, error)
	Cancel(ctx context.Context, id primitive.ObjectID) error
}

// ─── ScheduleService ──────────────────────────────────────────────────────────

type ScheduleService struct {
	templates     ShiftTemplateRepository
	periods       SchedulePeriodRepository
	slots         ScheduleSlotRepository
	registrations ShiftRegistrationRepository
}

func NewScheduleService(
	templates ShiftTemplateRepository,
	periods SchedulePeriodRepository,
	slots ScheduleSlotRepository,
	registrations ShiftRegistrationRepository,
) *ScheduleService {
	return &ScheduleService{templates, periods, slots, registrations}
}

// ─── Template methods ─────────────────────────────────────────────────────────

func (s *ScheduleService) CreateTemplate(ctx context.Context, t *schedule.ShiftTemplate) error {
	if t.Name == "" {
		return fmt.Errorf("tên ca không được để trống")
	}
	if t.StartTime == "" || t.EndTime == "" {
		return fmt.Errorf("giờ bắt đầu và kết thúc không được để trống")
	}
	if t.MaxSlots <= 0 {
		t.MaxSlots = 10
	}
	if len(t.Roles) == 0 {
		t.Roles = []string{"all"}
	}
	t.IsActive = true
	return s.templates.Create(ctx, t)
}

func (s *ScheduleService) GetTemplates(ctx context.Context) ([]*schedule.ShiftTemplate, error) {
	return s.templates.FindAll(ctx)
}

func (s *ScheduleService) UpdateTemplate(ctx context.Context, id primitive.ObjectID, t *schedule.ShiftTemplate) error {
	existing, err := s.templates.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("không tìm thấy ca")
	}
	t.ID = existing.ID
	t.CreatedAt = existing.CreatedAt
	return s.templates.Update(ctx, id, t)
}

func (s *ScheduleService) DeleteTemplate(ctx context.Context, id primitive.ObjectID) error {
	return s.templates.Delete(ctx, id)
}

// ─── Period methods ───────────────────────────────────────────────────────────

func (s *ScheduleService) CreatePeriod(ctx context.Context, p *schedule.SchedulePeriod) error {
	if p.Name == "" {
		return fmt.Errorf("tên kỳ không được để trống")
	}
	if p.StartDate.IsZero() || p.EndDate.IsZero() {
		return fmt.Errorf("ngày bắt đầu và kết thúc không được để trống")
	}
	if p.EndDate.Before(p.StartDate) {
		return fmt.Errorf("ngày kết thúc phải sau ngày bắt đầu")
	}
	if p.RegisterDeadline.IsZero() {
		// Default: deadline = start of the period (staff can register up until it starts)
		p.RegisterDeadline = p.StartDate
	}
	p.Status = schedule.PeriodStatusDraft
	return s.periods.Create(ctx, p)
}

func (s *ScheduleService) GetPeriods(ctx context.Context) ([]*schedule.SchedulePeriod, error) {
	return s.periods.FindAll(ctx)
}

func (s *ScheduleService) SetPeriodStatus(ctx context.Context, id primitive.ObjectID, status schedule.PeriodStatus) error {
	p, err := s.periods.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("không tìm thấy kỳ đăng ký")
	}
	// Validate transitions
	switch status {
	case schedule.PeriodStatusOpen:
		if p.Status != schedule.PeriodStatusDraft {
			return fmt.Errorf("chỉ có thể mở kỳ ở trạng thái nháp")
		}
	case schedule.PeriodStatusClosed:
		if p.Status != schedule.PeriodStatusOpen {
			return fmt.Errorf("chỉ có thể đóng kỳ đang mở")
		}
	case schedule.PeriodStatusPublished:
		if p.Status != schedule.PeriodStatusClosed {
			return fmt.Errorf("chỉ có thể publish kỳ đã đóng")
		}
	case schedule.PeriodStatusCancelled:
		if p.Status == schedule.PeriodStatusCancelled {
			return fmt.Errorf("kỳ này đã bị huỷ rồi")
		}
	}
	return s.periods.UpdateStatus(ctx, id, status)
}

func (s *ScheduleService) GetCurrentPeriod(ctx context.Context) (*schedule.SchedulePeriod, error) {
	return s.periods.FindOpenOrPublished(ctx)
}

// ─── Slot methods ─────────────────────────────────────────────────────────────

func (s *ScheduleService) AddSlot(ctx context.Context, periodID primitive.ObjectID, templateID primitive.ObjectID, date time.Time) (*schedule.ScheduleSlot, error) {
	p, err := s.periods.FindByID(ctx, periodID)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy kỳ đăng ký")
	}
	if p.Status == schedule.PeriodStatusClosed || p.Status == schedule.PeriodStatusPublished {
		return nil, fmt.Errorf("không thể thêm ca vào kỳ đã đóng")
	}
	t, err := s.templates.FindByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy ca mẫu")
	}
	// Normalize date to midnight
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	slot := &schedule.ScheduleSlot{
		PeriodID:     periodID,
		TemplateID:   templateID,
		Date:         date,
		TemplateName: t.Name,
		StartTime:    t.StartTime,
		EndTime:      t.EndTime,
		Roles:        t.Roles,
		MaxSlots:     t.MaxSlots,
	}
	if err := s.slots.Create(ctx, slot); err != nil {
		return nil, err
	}
	return slot, nil
}

func (s *ScheduleService) RemoveSlot(ctx context.Context, slotID primitive.ObjectID) error {
	return s.slots.Delete(ctx, slotID)
}

func (s *ScheduleService) GetPeriodDetail(ctx context.Context, periodID primitive.ObjectID, userID *primitive.ObjectID) (*schedule.PeriodDetail, error) {
	p, err := s.periods.FindByID(ctx, periodID)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy kỳ đăng ký")
	}
	slotList, err := s.slots.FindByPeriod(ctx, periodID)
	if err != nil {
		return nil, err
	}

	var slotsWithReg []schedule.SlotWithRegistrations
	for _, sl := range slotList {
		regs, err := s.registrations.FindBySlot(ctx, sl.ID)
		if err != nil {
			return nil, err
		}
		regList := make([]schedule.ShiftRegistration, 0, len(regs))
		for _, r := range regs {
			regList = append(regList, *r)
		}
		swr := schedule.SlotWithRegistrations{
			ScheduleSlot:    *sl,
			RegisteredCount: len(regList),
			Registrations:   regList,
		}
		if userID != nil {
			for _, r := range regList {
				if r.UserID == *userID {
					copy := r
					swr.MyRegistration = &copy
					break
				}
			}
		}
		slotsWithReg = append(slotsWithReg, swr)
	}
	if slotsWithReg == nil {
		slotsWithReg = []schedule.SlotWithRegistrations{}
	}

	return &schedule.PeriodDetail{
		SchedulePeriod: *p,
		Slots:          slotsWithReg,
	}, nil
}

// ─── Registration methods ─────────────────────────────────────────────────────

func (s *ScheduleService) Register(ctx context.Context, slotID primitive.ObjectID, userID primitive.ObjectID, userName, userRole string) (*schedule.ShiftRegistration, error) {
	slot, err := s.slots.FindByID(ctx, slotID)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy ca")
	}

	p, err := s.periods.FindByID(ctx, slot.PeriodID)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy kỳ đăng ký")
	}
	if p.Status != schedule.PeriodStatusOpen {
		return nil, fmt.Errorf("kỳ đăng ký chưa mở hoặc đã đóng")
	}
	if time.Now().After(p.RegisterDeadline) {
		return nil, fmt.Errorf("đã quá hạn đăng ký")
	}

	// Check role eligibility
	if !roleAllowed(userRole, slot.Roles) {
		return nil, fmt.Errorf("vai trò của bạn không được đăng ký ca này")
	}

	// Check already registered
	existing, _ := s.registrations.FindByUserAndSlot(ctx, userID, slotID)
	if existing != nil {
		return nil, fmt.Errorf("bạn đã đăng ký ca này rồi")
	}

	// Check slot capacity
	count, err := s.registrations.CountBySlot(ctx, slotID)
	if err != nil {
		return nil, err
	}
	if slot.MaxSlots > 0 && count >= slot.MaxSlots {
		return nil, fmt.Errorf("ca này đã đủ người (%d/%d)", count, slot.MaxSlots)
	}

	reg := &schedule.ShiftRegistration{
		SlotID:   slotID,
		PeriodID: slot.PeriodID,
		UserID:   userID,
		UserName: userName,
		UserRole: userRole,
		Status:   schedule.RegistrationStatusConfirmed,
	}
	if err := s.registrations.Create(ctx, reg); err != nil {
		return nil, err
	}
	return reg, nil
}

func (s *ScheduleService) Unregister(ctx context.Context, slotID primitive.ObjectID, userID primitive.ObjectID) error {
	slot, err := s.slots.FindByID(ctx, slotID)
	if err != nil {
		return fmt.Errorf("không tìm thấy ca")
	}
	p, err := s.periods.FindByID(ctx, slot.PeriodID)
	if err != nil {
		return fmt.Errorf("không tìm thấy kỳ đăng ký")
	}
	if p.Status != schedule.PeriodStatusOpen {
		return fmt.Errorf("chỉ có thể hủy đăng ký khi kỳ đang mở")
	}
	if time.Now().After(p.RegisterDeadline) {
		return fmt.Errorf("đã quá hạn hủy đăng ký")
	}

	reg, err := s.registrations.FindByUserAndSlot(ctx, userID, slotID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("bạn chưa đăng ký ca này")
		}
		return err
	}
	return s.registrations.Cancel(ctx, reg.ID)
}

// CancelRegistration allows a manager to cancel any registration by its ID
func (s *ScheduleService) CancelRegistration(ctx context.Context, regID primitive.ObjectID) error {
	return s.registrations.Cancel(ctx, regID)
}

func (s *ScheduleService) GetMySchedule(ctx context.Context, userID primitive.ObjectID) (*schedule.PeriodDetail, error) {
	p, err := s.periods.FindOpenOrPublished(ctx)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return s.GetPeriodDetail(ctx, p.ID, &userID)
}

// ─── ScheduledStaff ───────────────────────────────────────────────────────────

type ScheduledStaff struct {
	StaffName      string  `json:"staff_name"`
	TemplateName   string  `json:"template_name"`
	StartTime      string  `json:"start_time"`
	EndTime        string  `json:"end_time"`
	ScheduledHours float64 `json:"scheduled_hours"`
}

func (s *ScheduleService) GetRegistrationsByDate(ctx context.Context, date time.Time) ([]ScheduledStaff, error) {
	slots, err := s.slots.FindByDate(ctx, date)
	if err != nil {
		return nil, err
	}
	var result []ScheduledStaff
	for _, sl := range slots {
		regs, err := s.registrations.FindBySlot(ctx, sl.ID)
		if err != nil {
			continue
		}
		hours := calcShiftHours(sl.StartTime, sl.EndTime)
		for _, reg := range regs {
			result = append(result, ScheduledStaff{
				StaffName:      reg.UserName,
				TemplateName:   sl.TemplateName,
				StartTime:      sl.StartTime,
				EndTime:        sl.EndTime,
				ScheduledHours: hours,
			})
		}
	}
	if result == nil {
		result = []ScheduledStaff{}
	}
	return result, nil
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func calcShiftHours(start, end string) float64 {
	var sh, sm, eh, em int
	fmt.Sscanf(start, "%d:%d", &sh, &sm)
	fmt.Sscanf(end, "%d:%d", &eh, &em)
	startMins := sh*60 + sm
	endMins := eh*60 + em
	if endMins <= startMins {
		endMins += 24 * 60
	}
	return float64(endMins-startMins) / 60
}

func roleAllowed(userRole string, slotRoles []string) bool {
	for _, r := range slotRoles {
		if r == "all" || r == userRole {
			return true
		}
	}
	return false
}
