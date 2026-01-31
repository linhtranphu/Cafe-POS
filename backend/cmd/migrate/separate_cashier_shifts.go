package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cafe-pos/backend/domain/cashier"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Migration script to separate cashier shifts from waiter/barista shifts
// This script:
// 1. Finds all shifts with role_type = "cashier" in the "shifts" collection
// 2. Transforms them into CashierShift format
// 3. Inserts them into the new "cashier_shifts" collection
// 4. Deletes the old cashier shifts from "shifts" collection

func main() {
	// Connect to MongoDB (no auth for local development)
	mongoURI := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database("cafe_pos")

	log.Println("🚀 Starting cashier shift migration...")
	
	if err := migrateCashierShifts(db); err != nil {
		log.Fatal("❌ Migration failed:", err)
	}

	log.Println("✅ Migration completed successfully!")
}

func migrateCashierShifts(db *mongo.Database) error {
	ctx := context.Background()
	
	shiftsCollection := db.Collection("shifts")
	cashierShiftsCollection := db.Collection("cashier_shifts")

	// Step 1: Find all cashier shifts in the old collection
	log.Println("📊 Finding cashier shifts in 'shifts' collection...")
	
	cursor, err := shiftsCollection.Find(ctx, bson.M{
		"role_type": "cashier",
	})
	if err != nil {
		return fmt.Errorf("failed to find cashier shifts: %w", err)
	}
	defer cursor.Close(ctx)

	// Step 2: Transform and insert into new collection
	var migratedCount int
	var skippedCount int

	for cursor.Next(ctx) {
		var oldShift struct {
			ID           primitive.ObjectID `bson:"_id"`
			Type         string             `bson:"type"`
			Status       string             `bson:"status"`
			RoleType     string             `bson:"role_type"`
			UserID       primitive.ObjectID `bson:"user_id"`
			UserName     string             `bson:"user_name"`
			CashierID    primitive.ObjectID `bson:"cashier_id,omitempty"`
			CashierName  string             `bson:"cashier_name,omitempty"`
			StartCash    float64            `bson:"start_cash"`
			EndCash      float64            `bson:"end_cash"`
			TotalRevenue float64            `bson:"total_revenue"`
			TotalOrders  int                `bson:"total_orders"`
			StartedAt    time.Time          `bson:"started_at"`
			EndedAt      *time.Time         `bson:"ended_at,omitempty"`
			CreatedAt    time.Time          `bson:"created_at"`
			UpdatedAt    time.Time          `bson:"updated_at"`
		}

		if err := cursor.Decode(&oldShift); err != nil {
			log.Printf("⚠️  Failed to decode shift %s: %v", oldShift.ID.Hex(), err)
			skippedCount++
			continue
		}

		// Check if already migrated
		existingCount, err := cashierShiftsCollection.CountDocuments(ctx, bson.M{"_id": oldShift.ID})
		if err != nil {
			log.Printf("⚠️  Failed to check if shift %s exists: %v", oldShift.ID.Hex(), err)
			skippedCount++
			continue
		}
		if existingCount > 0 {
			log.Printf("⏭️  Shift %s already migrated, skipping...", oldShift.ID.Hex())
			skippedCount++
			continue
		}

		// Transform to CashierShift format
		newCashierShift := transformToCashierShift(oldShift)

		// Insert into new collection
		_, err = cashierShiftsCollection.InsertOne(ctx, newCashierShift)
		if err != nil {
			log.Printf("⚠️  Failed to insert shift %s: %v", oldShift.ID.Hex(), err)
			skippedCount++
			continue
		}

		log.Printf("✓ Migrated shift %s (%s - %s)", 
			oldShift.ID.Hex(), 
			oldShift.UserName, 
			oldShift.StartedAt.Format("2006-01-02 15:04"))
		migratedCount++
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("cursor error: %w", err)
	}

	log.Printf("📈 Migration summary: %d migrated, %d skipped", migratedCount, skippedCount)

	// Step 3: Delete old cashier shifts from shifts collection
	if migratedCount > 0 {
		log.Println("🗑️  Deleting old cashier shifts from 'shifts' collection...")
		
		deleteResult, err := shiftsCollection.DeleteMany(ctx, bson.M{
			"role_type": "cashier",
		})
		if err != nil {
			return fmt.Errorf("failed to delete old cashier shifts: %w", err)
		}

		log.Printf("✓ Deleted %d old cashier shifts", deleteResult.DeletedCount)
	}

	// Step 4: Create indexes on new collection
	log.Println("📑 Creating indexes on 'cashier_shifts' collection...")
	
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "cashier_id", Value: 1},
				{Key: "start_time", Value: -1},
			},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "end_time", Value: -1}},
		},
	}

	_, err = cashierShiftsCollection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Printf("⚠️  Warning: Failed to create indexes: %v", err)
	} else {
		log.Println("✓ Indexes created successfully")
	}

	return nil
}

func transformToCashierShift(oldShift interface{}) *cashier.CashierShift {
	// Type assertion to get the old shift data
	old := oldShift.(struct {
		ID           primitive.ObjectID `bson:"_id"`
		Type         string             `bson:"type"`
		Status       string             `bson:"status"`
		RoleType     string             `bson:"role_type"`
		UserID       primitive.ObjectID `bson:"user_id"`
		UserName     string             `bson:"user_name"`
		CashierID    primitive.ObjectID `bson:"cashier_id,omitempty"`
		CashierName  string             `bson:"cashier_name,omitempty"`
		StartCash    float64            `bson:"start_cash"`
		EndCash      float64            `bson:"end_cash"`
		TotalRevenue float64            `bson:"total_revenue"`
		TotalOrders  int                `bson:"total_orders"`
		StartedAt    time.Time          `bson:"started_at"`
		EndedAt      *time.Time         `bson:"ended_at,omitempty"`
		CreatedAt    time.Time          `bson:"created_at"`
		UpdatedAt    time.Time          `bson:"updated_at"`
	})

	// Determine cashier ID and name
	cashierID := old.UserID
	cashierName := old.UserName
	if !old.CashierID.IsZero() {
		cashierID = old.CashierID
	}
	if old.CashierName != "" {
		cashierName = old.CashierName
	}

	// Convert status
	var status cashier.CashierShiftStatus
	switch old.Status {
	case "OPEN":
		status = cashier.CashierShiftOpen
	case "CLOSED":
		status = cashier.CashierShiftClosed
	default:
		status = cashier.CashierShiftOpen
	}

	// Create new cashier shift
	newShift := &cashier.CashierShift{
		ID:            old.ID,
		CashierID:     cashierID,
		CashierName:   cashierName,
		StartTime:     old.StartedAt,
		EndTime:       old.EndedAt,
		Status:        status,
		StartingFloat: old.StartCash,
		SystemCash:    old.EndCash, // Use end_cash as system_cash approximation
		AuditLog:      []cashier.AuditLogEntry{},
		CreatedAt:     old.CreatedAt,
		UpdatedAt:     old.UpdatedAt,
	}

	// If shift is closed, set actual cash to end cash
	if status == cashier.CashierShiftClosed {
		actualCash := old.EndCash
		newShift.ActualCash = &actualCash
		
		// Create a zero variance (since we don't have historical variance data)
		variance := cashier.NewVariance(old.EndCash, old.EndCash)
		newShift.Variance = variance
	}

	return newShift
}
