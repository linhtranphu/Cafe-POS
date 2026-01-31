package services

import (
	"context"
	"testing"
	"time"

	"cafe-pos/backend/domain/facility"
	"cafe-pos/backend/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Test database setup
func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	dbName := "cafe_pos_test_" + primitive.NewObjectID().Hex()
	db := client.Database(dbName)

	cleanup := func() {
		db.Drop(ctx)
		client.Disconnect(ctx)
	}

	return db, cleanup
}

func TestCreateFacility(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewFacilityRepository(db)
	service := NewFacilityService(repo)
	ctx := context.Background()
	userID := primitive.NewObjectID()
	username := "testuser"

	t.Run("Success - Create valid facility", func(t *testing.T) {
		f := &facility.Facility{
			Name:     "Test Machine",
			Type:     "Equipment",
			Area:     "Kitchen",
			Quantity: 1,
			Status:   facility.StatusInUse,
		}

		err := service.CreateFacility(ctx, f, userID, username)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if f.ID.IsZero() {
			t.Error("Expected facility ID to be set")
		}

		// Verify history was created
		history, _ := repo.GetHistory(ctx, f.ID)
		if len(history) != 1 {
			t.Errorf("Expected 1 history record, got %d", len(history))
		}
		if history[0].Action != facility.ActionCreated {
			t.Errorf("Expected action 'created', got %s", history[0].Action)
		}
	})

	t.Run("Error - Missing required fields", func(t *testing.T) {
		f := &facility.Facility{
			Name: "Test",
			// Missing Type and Area
		}

		err := service.CreateFacility(ctx, f, userID, username)
		if err == nil {
			t.Error("Expected error for missing required fields")
		}
	})

	t.Run("Error - Invalid quantity", func(t *testing.T) {
		f := &facility.Facility{
			Name:     "Test",
			Type:     "Equipment",
			Area:     "Kitchen",
			Quantity: 0, // Invalid
		}

		err := service.CreateFacility(ctx, f, userID, username)
		if err == nil {
			t.Error("Expected error for invalid quantity")
		}
	})
}

func TestDeleteFacility(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewFacilityRepository(db)
	service := NewFacilityService(repo)
	ctx := context.Background()
	userID := primitive.NewObjectID()
	username := "testuser"

	t.Run("Success - Delete facility without maintenance history", func(t *testing.T) {
		f := &facility.Facility{
			Name:     "Test Machine",
			Type:     "Equipment",
			Area:     "Kitchen",
			Quantity: 1,
			Status:   facility.StatusInUse,
		}
		service.CreateFacility(ctx, f, userID, username)

		err := service.DeleteFacility(ctx, f.ID, userID, username)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify facility was deleted
		_, err = repo.GetByID(ctx, f.ID)
		if err == nil {
			t.Error("Expected facility to be deleted")
		}
	})

	t.Run("Error - Cannot delete facility with maintenance history", func(t *testing.T) {
		f := &facility.Facility{
			Name:     "Test Machine 2",
			Type:     "Equipment",
			Area:     "Kitchen",
			Quantity: 1,
			Status:   facility.StatusInUse,
		}
		service.CreateFacility(ctx, f, userID, username)

		// Add maintenance record
		maintenance := &facility.MaintenanceRecord{
			FacilityID:  f.ID,
			Type:        "scheduled",
			Description: "Test maintenance",
			Date:        time.Now(),
			UserID:      userID,
			Username:    username,
		}
		repo.CreateMaintenanceRecord(ctx, maintenance)

		err := service.DeleteFacility(ctx, f.ID, userID, username)
		if err == nil {
			t.Error("Expected error when deleting facility with maintenance history")
		}
		if err.Error() != "không thể xóa tài sản đã có lịch sử bảo trì" {
			t.Errorf("Expected specific error message, got %v", err)
		}
	})
}

func TestCreateMaintenanceRecord(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewFacilityRepository(db)
	service := NewFacilityService(repo)
	ctx := context.Background()

	// Create facility first
	f := &facility.Facility{
		Name:     "Test Machine",
		Type:     "Equipment",
		Area:     "Kitchen",
		Quantity: 1,
		Status:   facility.StatusInUse,
	}
	userID := primitive.NewObjectID()
	service.CreateFacility(ctx, f, userID, "testuser")

	t.Run("Success - Create maintenance record", func(t *testing.T) {
		record := &facility.MaintenanceRecord{
			FacilityID:  f.ID,
			Type:        "scheduled",
			Description: "Regular maintenance",
			Cost:        100000,
			Date:        time.Now(),
			UserID:      userID,
			Username:    "testuser",
		}

		err := service.CreateMaintenanceRecord(ctx, record)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if record.ID.IsZero() {
			t.Error("Expected maintenance record ID to be set")
		}
	})

	t.Run("Error - Missing required fields", func(t *testing.T) {
		record := &facility.MaintenanceRecord{
			// Missing FacilityID and Description
			Type: "scheduled",
		}

		err := service.CreateMaintenanceRecord(ctx, record)
		if err == nil {
			t.Error("Expected error for missing required fields")
		}
	})
}
