package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"cafe-pos/backend/domain/facility"
)

type FacilityRepository struct {
	db                *mongo.Database
	facilities        *mongo.Collection
	maintenanceRecords *mongo.Collection
	issueReports      *mongo.Collection
	facilityHistory   *mongo.Collection
}

func NewFacilityRepository(db *mongo.Database) *FacilityRepository {
	return &FacilityRepository{
		db:                db,
		facilities:        db.Collection("facilities"),
		maintenanceRecords: db.Collection("maintenance_records"),
		issueReports:      db.Collection("issue_reports"),
		facilityHistory:   db.Collection("facility_history"),
	}
}

func (r *FacilityRepository) GetAll(ctx context.Context) ([]facility.Facility, error) {
	cursor, err := r.facilities.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var facilities []facility.Facility
	if err = cursor.All(ctx, &facilities); err != nil {
		return nil, err
	}
	return facilities, nil
}

func (r *FacilityRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*facility.Facility, error) {
	var f facility.Facility
	err := r.facilities.FindOne(ctx, bson.M{"_id": id}).Decode(&f)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FacilityRepository) Create(ctx context.Context, f *facility.Facility) error {
	f.ID = primitive.NewObjectID()
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
	
	_, err := r.facilities.InsertOne(ctx, f)
	return err
}

func (r *FacilityRepository) Update(ctx context.Context, id primitive.ObjectID, f *facility.Facility) error {
	f.UpdatedAt = time.Now()
	
	update := bson.M{
		"$set": bson.M{
			"name":          f.Name,
			"type":          f.Type,
			"area":          f.Area,
			"quantity":      f.Quantity,
			"status":        f.Status,
			"purchase_date": f.PurchaseDate,
			"cost":          f.Cost,
			"supplier":      f.Supplier,
			"notes":         f.Notes,
			"updated_at":    f.UpdatedAt,
		},
	}
	
	_, err := r.facilities.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *FacilityRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.facilities.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *FacilityRepository) CreateMaintenanceRecord(ctx context.Context, record *facility.MaintenanceRecord) error {
	record.ID = primitive.NewObjectID()
	record.CreatedAt = time.Now()
	
	_, err := r.maintenanceRecords.InsertOne(ctx, record)
	return err
}

func (r *FacilityRepository) GetMaintenanceHistory(ctx context.Context, facilityID primitive.ObjectID) ([]facility.MaintenanceRecord, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.maintenanceRecords.Find(ctx, bson.M{"facility_id": facilityID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []facility.MaintenanceRecord
	if err = cursor.All(ctx, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func (r *FacilityRepository) CreateIssueReport(ctx context.Context, report *facility.IssueReport) error {
	report.ID = primitive.NewObjectID()
	report.CreatedAt = time.Now()
	report.Status = "open"
	
	_, err := r.issueReports.InsertOne(ctx, report)
	return err
}

func (r *FacilityRepository) GetIssueReports(ctx context.Context) ([]facility.IssueReport, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.issueReports.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reports []facility.IssueReport
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *FacilityRepository) CreateHistory(ctx context.Context, history *facility.FacilityHistory) error {
	history.ID = primitive.NewObjectID()
	history.CreatedAt = time.Now()
	
	_, err := r.facilityHistory.InsertOne(ctx, history)
	return err
}

func (r *FacilityRepository) GetHistory(ctx context.Context, facilityID primitive.ObjectID) ([]facility.FacilityHistory, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(50)
	cursor, err := r.facilityHistory.Find(ctx, bson.M{"facility_id": facilityID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var history []facility.FacilityHistory
	if err = cursor.All(ctx, &history); err != nil {
		return nil, err
	}
	return history, nil
}

func (r *FacilityRepository) GetHistoryWithFilter(ctx context.Context, filter facility.FacilityHistoryFilter) ([]facility.FacilityHistory, error) {
	mongoFilter := bson.M{}
	
	if filter.FacilityID != nil {
		mongoFilter["facility_id"] = *filter.FacilityID
	}
	if filter.Action != nil {
		mongoFilter["action"] = *filter.Action
	}
	if filter.DateFrom != nil || filter.DateTo != nil {
		dateFilter := bson.M{}
		if filter.DateFrom != nil {
			dateFilter["$gte"] = *filter.DateFrom
		}
		if filter.DateTo != nil {
			dateFilter["$lte"] = *filter.DateTo
		}
		mongoFilter["created_at"] = dateFilter
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	if filter.Limit > 0 {
		opts.SetLimit(int64(filter.Limit))
	}
	if filter.Offset > 0 {
		opts.SetSkip(int64(filter.Offset))
	}

	cursor, err := r.facilityHistory.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var history []facility.FacilityHistory
	if err = cursor.All(ctx, &history); err != nil {
		return nil, err
	}
	return history, nil
}

func (r *FacilityRepository) GetWithFilter(ctx context.Context, filter facility.FacilityFilter) ([]facility.Facility, error) {
	mongoFilter := bson.M{}
	
	if filter.Name != nil {
		mongoFilter["name"] = bson.M{"$regex": *filter.Name, "$options": "i"}
	}
	if filter.Type != nil {
		mongoFilter["type"] = *filter.Type
	}
	if filter.Area != nil {
		mongoFilter["area"] = *filter.Area
	}
	if filter.Status != nil {
		mongoFilter["status"] = *filter.Status
	}

	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	if filter.Limit > 0 {
		opts.SetLimit(int64(filter.Limit))
	}
	if filter.Offset > 0 {
		opts.SetSkip(int64(filter.Offset))
	}

	cursor, err := r.facilities.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var facilities []facility.Facility
	if err = cursor.All(ctx, &facilities); err != nil {
		return nil, err
	}
	return facilities, nil
}

func (r *FacilityRepository) CountWithFilter(ctx context.Context, filter facility.FacilityFilter) (int64, error) {
	mongoFilter := bson.M{}
	
	if filter.Name != nil {
		mongoFilter["name"] = bson.M{"$regex": *filter.Name, "$options": "i"}
	}
	if filter.Type != nil {
		mongoFilter["type"] = *filter.Type
	}
	if filter.Area != nil {
		mongoFilter["area"] = *filter.Area
	}
	if filter.Status != nil {
		mongoFilter["status"] = *filter.Status
	}

	return r.facilities.CountDocuments(ctx, mongoFilter)
}

func (r *FacilityRepository) GetScheduledMaintenance(ctx context.Context) ([]facility.ScheduledMaintenance, error) {
	opts := options.Find().SetSort(bson.D{{Key: "scheduled_date", Value: 1}})
	cursor, err := r.db.Collection("scheduled_maintenance").Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []facility.ScheduledMaintenance
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *FacilityRepository) GetNextMaintenanceDate(ctx context.Context, facilityID primitive.ObjectID) (map[string]interface{}, error) {
	filter := bson.M{
		"facility_id": facilityID,
		"status": "pending",
		"scheduled_date": bson.M{"$gte": time.Now()},
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "scheduled_date", Value: 1}})
	
	var task facility.ScheduledMaintenance
	err := r.db.Collection("scheduled_maintenance").FindOne(ctx, filter, opts).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return map[string]interface{}{"next_date": nil}, nil
		}
		return nil, err
	}
	
	return map[string]interface{}{"next_date": task.ScheduledDate}, nil
}

func (r *FacilityRepository) GetMaintenanceDue(ctx context.Context) ([]facility.ScheduledMaintenance, error) {
	filter := bson.M{
		"status": "pending",
		"scheduled_date": bson.M{"$lte": time.Now()},
	}
	opts := options.Find().SetSort(bson.D{{Key: "scheduled_date", Value: 1}})
	
	cursor, err := r.db.Collection("scheduled_maintenance").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []facility.ScheduledMaintenance
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *FacilityRepository) GetMaintenanceStats(ctx context.Context, facilityID primitive.ObjectID) (map[string]interface{}, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"facility_id": facilityID}},
		{"$group": bson.M{
			"_id": nil,
			"total": bson.M{"$sum": 1},
			"totalCost": bson.M{"$sum": "$cost"},
			"scheduledCost": bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$eq": []string{"$type", "scheduled"}}, "$cost", 0}}},
			"emergencyCost": bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$eq": []string{"$type", "emergency"}}, "$cost", 0}}},
		}},
	}
	
	cursor, err := r.maintenanceRecords.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	
	if len(results) == 0 {
		return map[string]interface{}{
			"total": 0,
			"totalCost": 0,
			"avgInterval": 0,
			"scheduledCost": 0,
			"emergencyCost": 0,
		}, nil
	}
	
	stats := results[0]
	stats["avgInterval"] = 30 // Default average interval in days
	return stats, nil
}

func (r *FacilityRepository) GetStatusHistory(ctx context.Context, facilityID primitive.ObjectID) ([]facility.FacilityHistory, error) {
	filter := bson.M{
		"facility_id": facilityID,
		"action": "status_change",
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(50)
	
	cursor, err := r.facilityHistory.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var history []facility.FacilityHistory
	if err = cursor.All(ctx, &history); err != nil {
		return nil, err
	}
	return history, nil
}

// FacilityType and FacilityArea management

func (r *FacilityRepository) CreateFacilityType(ctx context.Context, ft *facility.FacilityType) error {
	ft.CreatedAt = time.Now()
	result, err := r.db.Collection("facility_types").InsertOne(ctx, ft)
	if err != nil {
		return err
	}
	ft.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *FacilityRepository) GetFacilityTypes(ctx context.Context) ([]facility.FacilityType, error) {
	cursor, err := r.db.Collection("facility_types").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var types []facility.FacilityType
	if err = cursor.All(ctx, &types); err != nil {
		return nil, err
	}
	return types, nil
}

func (r *FacilityRepository) DeleteFacilityType(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.db.Collection("facility_types").DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *FacilityRepository) CreateFacilityArea(ctx context.Context, fa *facility.FacilityArea) error {
	fa.CreatedAt = time.Now()
	result, err := r.db.Collection("facility_areas").InsertOne(ctx, fa)
	if err != nil {
		return err
	}
	fa.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *FacilityRepository) GetFacilityAreas(ctx context.Context) ([]facility.FacilityArea, error) {
	cursor, err := r.db.Collection("facility_areas").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var areas []facility.FacilityArea
	if err = cursor.All(ctx, &areas); err != nil {
		return nil, err
	}
	return areas, nil
}

func (r *FacilityRepository) DeleteFacilityArea(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.db.Collection("facility_areas").DeleteOne(ctx, bson.M{"_id": id})
	return err
}
