package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/infrastructure/mongodb"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "cafe_pos"
	}

	fmt.Printf("🔗 Connecting to MongoDB: %s/%s\n", mongoURI, dbName)

	// Connect to MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(dbName)
	shiftRepo := mongodb.NewShiftRepository(db)

	fmt.Println("\n🧪 Test: Payment Shift Update")
	fmt.Println("==============================\n")

	// Create a test shift
	userID := primitive.NewObjectID()
	shift := &order.Shift{
		Type:              order.ShiftMorning,
		Status:            order.ShiftOpen,
		RoleType:          order.RoleWaiter,
		UserID:            userID,
		UserName:          "Test Waiter",
		StartCash:         100000,
		CurrentCash:       100000,
		RemainingCash:     100000,
		TransferRevenue:   0,
		RemainingTransfer: 0,
		TotalRevenue:      0,
		StartedAt:         time.Now(),
	}

	if err := shiftRepo.Create(ctx, shift); err != nil {
		log.Fatalf("Failed to create shift: %v", err)
	}

	fmt.Printf("✅ Created shift: %s\n", shift.ID.Hex())
	fmt.Printf("   StartCash: %.0f\n", shift.StartCash)
	fmt.Printf("   CurrentCash: %.0f\n", shift.CurrentCash)
	fmt.Printf("   RemainingCash: %.0f\n", shift.RemainingCash)
	fmt.Printf("   TransferRevenue: %.0f\n", shift.TransferRevenue)
	fmt.Printf("   RemainingTransfer: %.0f\n", shift.RemainingTransfer)
	fmt.Printf("   TotalRevenue: %.0f\n\n", shift.TotalRevenue)

	// Simulate CASH payment
	fmt.Println("💰 Simulating CASH payment of 50,000 VND...")
	cashAmount := 50000.0

	// Read shift
	shift, err = shiftRepo.FindByID(ctx, shift.ID)
	if err != nil {
		log.Fatalf("Failed to find shift: %v", err)
	}

	fmt.Printf("   BEFORE: CurrentCash=%.0f, RemainingCash=%.0f, TotalRevenue=%.0f\n",
		shift.CurrentCash, shift.RemainingCash, shift.TotalRevenue)

	// Update shift
	shift.CurrentCash += cashAmount
	shift.RemainingCash += cashAmount
	shift.TotalRevenue += cashAmount

	fmt.Printf("   AFTER (in memory): CurrentCash=%.0f, RemainingCash=%.0f, TotalRevenue=%.0f\n",
		shift.CurrentCash, shift.RemainingCash, shift.TotalRevenue)

	// Save to DB
	if err := shiftRepo.Update(ctx, shift.ID, shift); err != nil {
		log.Fatalf("Failed to update shift: %v", err)
	}

	fmt.Println("   ✅ Updated shift in DB")

	// Verify by reading back
	verifyShift, err := shiftRepo.FindByID(ctx, shift.ID)
	if err != nil {
		log.Fatalf("Failed to verify shift: %v", err)
	}

	fmt.Printf("   VERIFY (from DB): CurrentCash=%.0f, RemainingCash=%.0f, TotalRevenue=%.0f\n\n",
		verifyShift.CurrentCash, verifyShift.RemainingCash, verifyShift.TotalRevenue)

	// Simulate TRANSFER payment
	fmt.Println("🏦 Simulating TRANSFER payment of 75,000 VND...")
	transferAmount := 75000.0

	// Read shift again
	shift, err = shiftRepo.FindByID(ctx, shift.ID)
	if err != nil {
		log.Fatalf("Failed to find shift: %v", err)
	}

	fmt.Printf("   BEFORE: TransferRevenue=%.0f, RemainingTransfer=%.0f, TotalRevenue=%.0f\n",
		shift.TransferRevenue, shift.RemainingTransfer, shift.TotalRevenue)

	// Update shift
	shift.TransferRevenue += transferAmount
	shift.RemainingTransfer += transferAmount
	shift.TotalRevenue += transferAmount

	fmt.Printf("   AFTER (in memory): TransferRevenue=%.0f, RemainingTransfer=%.0f, TotalRevenue=%.0f\n",
		shift.TransferRevenue, shift.RemainingTransfer, shift.TotalRevenue)

	// Save to DB
	if err := shiftRepo.Update(ctx, shift.ID, shift); err != nil {
		log.Fatalf("Failed to update shift: %v", err)
	}

	fmt.Println("   ✅ Updated shift in DB")

	// Verify by reading back
	verifyShift, err = shiftRepo.FindByID(ctx, shift.ID)
	if err != nil {
		log.Fatalf("Failed to verify shift: %v", err)
	}

	fmt.Printf("   VERIFY (from DB): TransferRevenue=%.0f, RemainingTransfer=%.0f, TotalRevenue=%.0f\n\n",
		verifyShift.TransferRevenue, verifyShift.RemainingTransfer, verifyShift.TotalRevenue)

	// Final summary
	fmt.Println("📊 Final Shift State:")
	fmt.Println("====================")
	fmt.Printf("CurrentCash: %.0f (expected: 150,000)\n", verifyShift.CurrentCash)
	fmt.Printf("RemainingCash: %.0f (expected: 150,000)\n", verifyShift.RemainingCash)
	fmt.Printf("TransferRevenue: %.0f (expected: 75,000)\n", verifyShift.TransferRevenue)
	fmt.Printf("RemainingTransfer: %.0f (expected: 75,000)\n", verifyShift.RemainingTransfer)
	fmt.Printf("TotalRevenue: %.0f (expected: 125,000)\n\n", verifyShift.TotalRevenue)

	// Check results
	success := true
	if verifyShift.CurrentCash != 150000 {
		fmt.Printf("❌ CurrentCash mismatch: got %.0f, expected 150000\n", verifyShift.CurrentCash)
		success = false
	}
	if verifyShift.TransferRevenue != 75000 {
		fmt.Printf("❌ TransferRevenue mismatch: got %.0f, expected 75000\n", verifyShift.TransferRevenue)
		success = false
	}
	if verifyShift.TotalRevenue != 125000 {
		fmt.Printf("❌ TotalRevenue mismatch: got %.0f, expected 125000\n", verifyShift.TotalRevenue)
		success = false
	}

	if success {
		fmt.Println("✅ All tests PASSED!")
	} else {
		fmt.Println("❌ Some tests FAILED!")
		os.Exit(1)
	}
}
