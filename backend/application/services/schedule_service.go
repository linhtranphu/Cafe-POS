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
	FindByDateRange(ctx context.Context, from, to time.Time) ([]*schedule.ScheduleSlot, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type ShiftRegistrationRepository interface {
	Create(ctx context.Context, reg *schedule.ShiftRegistration) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*schedule.ShiftRegistration, error)
	FindBySlot(ctx context.Context, slotID primitive.ObjectID) ([]*schedule.ShiftRegistration, error)
	FindByUserAndPeriod(ctx context.Context, userID, periodID primitive.ObjectID) ([]*schedule.ShiftRegistration, error)
	FindByUserAndSlot(ctx context.Context, userID, slotID primitive.ObjectID) (*schedule.ShiftRegistration, error)
	CountBySlot(ctx context.Context, slotID primitive.ObjectID) (int, error)
	Cancel(ctx context.Context, id primitive.ObjectID) error
}

type WeeklyShiftTemplateRepository interface {
	Create(ctx context.Context, wt *schedule.WeeklyShiftTemplate) error
	FindAll(ctx context.Context) ([]*schedule.WeeklyShiftTemplate, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*schedule.WeeklyShiftTemplate, error)
	Update(ctx context.Context, id primitive.ObjectID, wt *schedule.WeeklyShiftTemplate) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// ─── ScheduleService ──────────────────────────────────────────────────────────

type ScheduleService struct {
	templates        ShiftTemplateRepository
	periods          SchedulePeriodRepository
	slots            ScheduleSlotRepository
	registrations    ShiftRegistrationRepository
	weeklyTemplates  WeeklyShiftTemplateRepository
}

func NewScheduleService(
	templates ShiftTemplateRepository,
	periods SchedulePeriodRepository,
	slots ScheduleSlotRepository,
	registrations ShiftRegistrationRepository,
	weeklyTemplates WeeklyShiftTemplateRepository,
) *ScheduleService {
	return &ScheduleService{templates, periods, slots, registrations, weeklyTemplates}
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

// ─── WeeklyShiftTemplate methods ──────────────────────────────────────────────

func (s *ScheduleService) CreateWeeklyTemplate(ctx context.Context, wt *schedule.WeeklyShiftTemplate) error {
	if wt.Name == "" {
		return fmt.Errorf("tên template tuần không được để trống")
	}
	wt.IsActive = true
	return s.weeklyTemplates.Create(ctx, wt)
}

func (s *ScheduleService) GetWeeklyTemplates(ctx context.Context) ([]*schedule.WeeklyShiftTemplate, error) {
	return s.weeklyTemplates.FindAll(ctx)
}

func (s *ScheduleService) UpdateWeeklyTemplate(ctx context.Context, id primitive.ObjectID, wt *schedule.WeeklyShiftTemplate) error {
	existing, err := s.weeklyTemplates.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("không tìm thấy template tuần")
	}
	wt.ID = existing.ID
	wt.CreatedAt = existing.CreatedAt
	return s.weeklyTemplates.Update(ctx, id, wt)
}

func (s *ScheduleService) DeleteWeeklyTemplate(ctx context.Context, id primitive.ObjectID) error {
	return s.weeklyTemplates.Delete(ctx, id)
}

// applyWeeklyTemplate creates slots for every day in [startDate, endDate] according to the weekly template.
func (s *ScheduleService) applyWeeklyTemplate(ctx context.Context, periodID primitive.ObjectID, weeklyTemplateID primitive.ObjectID, startDate, endDate time.Time) error {
	wt, err := s.weeklyTemplates.FindByID(ctx, weeklyTemplateID)
	if err != nil {
		return fmt.Errorf("không tìm thấy template tuần")
	}
	dayMap := make(map[int][]primitive.ObjectID)
	for _, d := range wt.Days {
		if len(d.TemplateIDs) > 0 {
			dayMap[d.DayOfWeek] = d.TemplateIDs
		}
	}
	cur := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, time.UTC)
	for !cur.After(end) {
		dow := int(cur.Weekday()) // 0=Sun, 1=Mon, ..., 6=Sat
		if templateIDs, ok := dayMap[dow]; ok {
			for _, templateID := range templateIDs {
				if _, err := s.AddSlot(ctx, periodID, templateID, cur); err != nil {
					// skip individual slot errors (e.g. duplicate), continue
					_ = err
				}
			}
		}
		cur = cur.AddDate(0, 0, 1)
	}
	return nil
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

// CreatePeriodWithWeeklyTemplate creates a period and optionally auto-generates slots from a weekly template.
func (s *ScheduleService) CreatePeriodWithWeeklyTemplate(ctx context.Context, p *schedule.SchedulePeriod, weeklyTemplateID *primitive.ObjectID) error {
	if err := s.CreatePeriod(ctx, p); err != nil {
		return err
	}
	if weeklyTemplateID != nil {
		if err := s.applyWeeklyTemplate(ctx, p.ID, *weeklyTemplateID, p.StartDate, p.EndDate); err != nil {
			return err
		}
	}
	return nil
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

	// Rule 1 & 2: lấy tất cả đăng ký hiện tại của user trong kỳ
	userRegs, err := s.registrations.FindByUserAndPeriod(ctx, userID, slot.PeriodID)
	if err != nil {
		return nil, err
	}
	newSlotHours := calcShiftHours(slot.StartTime, slot.EndTime)
	totalHours := newSlotHours
	for _, reg := range userRegs {
		if reg.Status != schedule.RegistrationStatusConfirmed {
			continue
		}
		existingSlot, err := s.slots.FindByID(ctx, reg.SlotID)
		if err != nil {
			continue
		}
		// Rule 1: không cho phép trùng giờ trong cùng ngày
		if existingSlot.Date.Equal(slot.Date) {
			if timesOverlap(existingSlot.StartTime, existingSlot.EndTime, slot.StartTime, slot.EndTime) {
				return nil, fmt.Errorf("ca này bị trùng giờ với ca đã đăng ký: %s (%s–%s)", existingSlot.TemplateName, existingSlot.StartTime, existingSlot.EndTime)
			}
		}
		totalHours += calcShiftHours(existingSlot.StartTime, existingSlot.EndTime)
	}
	// Rule 2: kiểm tra tổng giờ không vượt quá MaxHoursPercent
	if p.MaxHoursPercent > 0 {
		// Tổng số giờ của kỳ (tính cả ngày EndDate)
		periodHours := p.EndDate.Sub(p.StartDate).Hours() + 24
		maxHours := p.MaxHoursPercent / 100 * periodHours
		if totalHours > maxHours {
			return nil, fmt.Errorf("tổng giờ đăng ký (%.1f giờ) vượt quá giới hạn %.0f%% tổng giờ trong kỳ (tối đa %.1f giờ)", totalHours, p.MaxHoursPercent, maxHours)
		}
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
	if err := s.registrations.Cancel(ctx, reg.ID); err != nil {
		return err
	}
	count, err := s.registrations.CountBySlot(ctx, slotID)
	if err == nil && count == 0 {
		_ = s.slots.Delete(ctx, slotID)
	}
	return nil
}

// CancelRegistration allows a manager to cancel any registration by its ID.
// If no confirmed registrations remain for the slot, the slot is also deleted.
func (s *ScheduleService) CancelRegistration(ctx context.Context, regID primitive.ObjectID) error {
	reg, err := s.registrations.FindByID(ctx, regID)
	if err != nil {
		return err
	}
	if err := s.registrations.Cancel(ctx, regID); err != nil {
		return err
	}
	count, err := s.registrations.CountBySlot(ctx, reg.SlotID)
	if err == nil && count == 0 {
		_ = s.slots.Delete(ctx, reg.SlotID)
	}
	return nil
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

type ScheduledStaffWithDate struct {
	Date           string  `json:"date"`
	StaffName      string  `json:"staff_name"`
	TemplateName   string  `json:"template_name"`
	StartTime      string  `json:"start_time"`
	EndTime        string  `json:"end_time"`
	ScheduledHours float64 `json:"scheduled_hours"`
}

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

func (s *ScheduleService) GetRegistrationsByPeriod(ctx context.Context, periodID primitive.ObjectID) ([]ScheduledStaffWithDate, error) {
	slots, err := s.slots.FindByPeriod(ctx, periodID)
	if err != nil {
		return nil, err
	}
	var result []ScheduledStaffWithDate
	for _, sl := range slots {
		regs, err := s.registrations.FindBySlot(ctx, sl.ID)
		if err != nil {
			continue
		}
		hours := calcShiftHours(sl.StartTime, sl.EndTime)
		dateStr := sl.Date.Format("2006-01-02")
		for _, reg := range regs {
			result = append(result, ScheduledStaffWithDate{
				Date:           dateStr,
				StaffName:      reg.UserName,
				TemplateName:   sl.TemplateName,
				StartTime:      sl.StartTime,
				EndTime:        sl.EndTime,
				ScheduledHours: hours,
			})
		}
	}
	if result == nil {
		result = []ScheduledStaffWithDate{}
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

// timesOverlap kiểm tra hai khoảng thời gian (HH:MM) có bị chồng chéo không.
// Hỗ trợ ca qua đêm (end <= start).
func timesOverlap(start1, end1, start2, end2 string) bool {
	toMins := func(t string) int {
		var h, m int
		fmt.Sscanf(t, "%d:%d", &h, &m)
		return h*60 + m
	}
	s1, e1 := toMins(start1), toMins(end1)
	s2, e2 := toMins(start2), toMins(end2)
	if e1 <= s1 {
		e1 += 24 * 60 // ca qua đêm
	}
	if e2 <= s2 {
		e2 += 24 * 60
	}
	return s1 < e2 && s2 < e1
}

func roleAllowed(userRole string, slotRoles []string) bool {
	for _, r := range slotRoles {
		if r == "all" || r == userRole {
			return true
		}
	}
	return false
}
