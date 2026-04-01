package mongodb

import (
	"context"
	"time"

	"cafe-pos/backend/domain/schedule"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ─── ShiftTemplateRepository ────────────────────────────────────────────────

type ShiftTemplateRepository struct {
	col *mongo.Collection
}

func NewShiftTemplateRepository(db *mongo.Database) *ShiftTemplateRepository {
	return &ShiftTemplateRepository{col: db.Collection("shift_templates")}
}

func (r *ShiftTemplateRepository) Create(ctx context.Context, t *schedule.ShiftTemplate) error {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	result, err := r.col.InsertOne(ctx, t)
	if err != nil {
		return err
	}
	t.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ShiftTemplateRepository) FindAll(ctx context.Context) ([]*schedule.ShiftTemplate, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})
	cursor, err := r.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var list []*schedule.ShiftTemplate
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ShiftTemplateRepository) FindActive(ctx context.Context) ([]*schedule.ShiftTemplate, error) {
	opts := options.Find().SetSort(bson.D{{Key: "start_time", Value: 1}})
	cursor, err := r.col.Find(ctx, bson.M{"is_active": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var list []*schedule.ShiftTemplate
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ShiftTemplateRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*schedule.ShiftTemplate, error) {
	var t schedule.ShiftTemplate
	if err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *ShiftTemplateRepository) Update(ctx context.Context, id primitive.ObjectID, t *schedule.ShiftTemplate) error {
	t.UpdatedAt = time.Now()
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": t})
	return err
}

func (r *ShiftTemplateRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// ─── SchedulePeriodRepository ────────────────────────────────────────────────

type SchedulePeriodRepository struct {
	col *mongo.Collection
}

func NewSchedulePeriodRepository(db *mongo.Database) *SchedulePeriodRepository {
	return &SchedulePeriodRepository{col: db.Collection("schedule_periods")}
}

func (r *SchedulePeriodRepository) Create(ctx context.Context, p *schedule.SchedulePeriod) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	result, err := r.col.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	p.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *SchedulePeriodRepository) FindAll(ctx context.Context) ([]*schedule.SchedulePeriod, error) {
	opts := options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}})
	cursor, err := r.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var list []*schedule.SchedulePeriod
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *SchedulePeriodRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*schedule.SchedulePeriod, error) {
	var p schedule.SchedulePeriod
	if err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *SchedulePeriodRepository) FindOpenOrPublished(ctx context.Context) (*schedule.SchedulePeriod, error) {
	var p schedule.SchedulePeriod
	// Bao gồm cả "closed" vì manager cần tạo giờ công sau khi đóng đăng ký
	// (trước khi publish lịch chính thức).
	err := r.col.FindOne(ctx, bson.M{
		"status": bson.M{"$in": []schedule.PeriodStatus{
			schedule.PeriodStatusOpen,
			schedule.PeriodStatusClosed,
			schedule.PeriodStatusPublished,
		}},
	}, options.FindOne().SetSort(bson.D{{Key: "start_date", Value: -1}})).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *SchedulePeriodRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status schedule.PeriodStatus) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"status":     status,
		"updated_at": time.Now(),
	}})
	return err
}

// ─── ScheduleSlotRepository ───────────────────────────────────────────────────

type ScheduleSlotRepository struct {
	col *mongo.Collection
}

func NewScheduleSlotRepository(db *mongo.Database) *ScheduleSlotRepository {
	return &ScheduleSlotRepository{col: db.Collection("schedule_slots")}
}

func (r *ScheduleSlotRepository) Create(ctx context.Context, s *schedule.ScheduleSlot) error {
	s.CreatedAt = time.Now()
	result, err := r.col.InsertOne(ctx, s)
	if err != nil {
		return err
	}
	s.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ScheduleSlotRepository) FindByDate(ctx context.Context, date time.Time) ([]*schedule.ScheduleSlot, error) {
	from := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	opts := options.Find().SetSort(bson.D{{Key: "start_time", Value: 1}})
	cursor, err := r.col.Find(ctx, bson.M{
		"date": bson.M{"$gte": from, "$lt": to},
	}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var list []*schedule.ScheduleSlot
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ScheduleSlotRepository) FindByDateRange(ctx context.Context, from, to time.Time) ([]*schedule.ScheduleSlot, error) {
	fromDay := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	toDay := time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 999999999, time.UTC)
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: 1}, {Key: "start_time", Value: 1}})
	cursor, err := r.col.Find(ctx, bson.M{
		"date": bson.M{"$gte": fromDay, "$lte": toDay},
	}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var list []*schedule.ScheduleSlot
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ScheduleSlotRepository) FindByPeriod(ctx context.Context, periodID primitive.ObjectID) ([]*schedule.ScheduleSlot, error) {
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: 1}, {Key: "start_time", Value: 1}})
	cursor, err := r.col.Find(ctx, bson.M{"period_id": periodID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var list []*schedule.ScheduleSlot
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ScheduleSlotRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*schedule.ScheduleSlot, error) {
	var s schedule.ScheduleSlot
	if err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ScheduleSlotRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// ─── ShiftRegistrationRepository ─────────────────────────────────────────────

type ShiftRegistrationRepository struct {
	col *mongo.Collection
}

func NewShiftRegistrationRepository(db *mongo.Database) *ShiftRegistrationRepository {
	return &ShiftRegistrationRepository{col: db.Collection("shift_registrations")}
}

func (r *ShiftRegistrationRepository) Create(ctx context.Context, reg *schedule.ShiftRegistration) error {
	reg.RegisteredAt = time.Now()
	reg.UpdatedAt = time.Now()
	result, err := r.col.InsertOne(ctx, reg)
	if err != nil {
		return err
	}
	reg.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ShiftRegistrationRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*schedule.ShiftRegistration, error) {
	var reg schedule.ShiftRegistration
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&reg)
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func (r *ShiftRegistrationRepository) FindBySlot(ctx context.Context, slotID primitive.ObjectID) ([]*schedule.ShiftRegistration, error) {
	cursor, err := r.col.Find(ctx, bson.M{
		"slot_id": slotID,
		"status":  schedule.RegistrationStatusConfirmed,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var list []*schedule.ShiftRegistration
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ShiftRegistrationRepository) FindByUserAndPeriod(ctx context.Context, userID, periodID primitive.ObjectID) ([]*schedule.ShiftRegistration, error) {
	cursor, err := r.col.Find(ctx, bson.M{
		"user_id":   userID,
		"period_id": periodID,
		"status":    schedule.RegistrationStatusConfirmed,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var list []*schedule.ShiftRegistration
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ShiftRegistrationRepository) FindByUserAndSlot(ctx context.Context, userID, slotID primitive.ObjectID) (*schedule.ShiftRegistration, error) {
	var reg schedule.ShiftRegistration
	err := r.col.FindOne(ctx, bson.M{
		"user_id": userID,
		"slot_id": slotID,
		"status":  schedule.RegistrationStatusConfirmed,
	}).Decode(&reg)
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func (r *ShiftRegistrationRepository) CountBySlot(ctx context.Context, slotID primitive.ObjectID) (int, error) {
	count, err := r.col.CountDocuments(ctx, bson.M{
		"slot_id": slotID,
		"status":  schedule.RegistrationStatusConfirmed,
	})
	return int(count), err
}

func (r *ShiftRegistrationRepository) Cancel(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"status":     schedule.RegistrationStatusCancelled,
		"updated_at": time.Now(),
	}})
	return err
}

// ─── WeeklyShiftTemplateRepository ───────────────────────────────────────────

type WeeklyShiftTemplateRepository struct {
	col *mongo.Collection
}

func NewWeeklyShiftTemplateRepository(db *mongo.Database) *WeeklyShiftTemplateRepository {
	return &WeeklyShiftTemplateRepository{col: db.Collection("weekly_shift_templates")}
}

func (r *WeeklyShiftTemplateRepository) Create(ctx context.Context, wt *schedule.WeeklyShiftTemplate) error {
	wt.CreatedAt = time.Now()
	wt.UpdatedAt = time.Now()
	result, err := r.col.InsertOne(ctx, wt)
	if err != nil {
		return err
	}
	wt.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *WeeklyShiftTemplateRepository) FindAll(ctx context.Context) ([]*schedule.WeeklyShiftTemplate, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})
	cursor, err := r.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var list []*schedule.WeeklyShiftTemplate
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *WeeklyShiftTemplateRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*schedule.WeeklyShiftTemplate, error) {
	var wt schedule.WeeklyShiftTemplate
	if err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&wt); err != nil {
		return nil, err
	}
	return &wt, nil
}

func (r *WeeklyShiftTemplateRepository) Update(ctx context.Context, id primitive.ObjectID, wt *schedule.WeeklyShiftTemplate) error {
	wt.UpdatedAt = time.Now()
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": wt})
	return err
}

func (r *WeeklyShiftTemplateRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
