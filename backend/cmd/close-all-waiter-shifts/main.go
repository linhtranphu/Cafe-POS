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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load .env from parent directory (project root)
	if err := godotenv.Load("../.env"); err != nil {
		// Try current directory
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: No .env file found")
		}
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGO_URI")
	}
	if mongoURI == "" {
		log.Fatal("MONGODB_URI not found in environment")
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = os.Getenv("DB_NAME")
	}
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
	orderRepo := mongodb.NewOrderRepository(db)

	fmt.Println("\n🔄 Closing All Open Waiter Shifts")
	fmt.Println("==================================\n")

	// Find all open shifts
	openShifts, err := shiftRepo.FindOpenShifts(ctx)
	if err != nil {
		log.Fatalf("Failed to find open shifts: %v", err)
	}

	if len(openShifts) == 0 {
		fmt.Println("✅ No open shifts found. All shifts are already closed.")
		return
	}

	fmt.Printf("Found %d open shift(s)\n\n", len(openShifts))

	// Close each shift
	for i, shift := range openShifts {
		fmt.Printf("[%d/%d] Closing shift: %s\n", i+1, len(openShifts), shift.ID.Hex())
		fmt.Printf("   User: %s (%s)\n", shift.UserName, shift.RoleType)
		fmt.Printf("   Type: %s\n", shift.Type)
		fmt.Printf("   Started: %s\n", shift.StartedAt.Format("2006-01-02 15:04:05"))

		// Get orders for this shift
		orders, err := orderRepo.FindByShiftID(ctx, shift.ID)
		if err != nil {
			fmt.Printf("   ⚠️  Warning: Failed to get orders: %v\n", err)
			orders = []*order.Order{} // Continue with empty orders
		}

		// Calculate revenue
		totalRevenue := 0.0
		cashRevenue := 0.0
		transferRevenue := 0.0
		totalOrders := 0

		for _, o := range orders {
			if o.Status == order.StatusPaid || o.Status == order.StatusInProgress || o.Status == order.StatusServed {
				totalRevenue += o.Total
				totalOrders++

				if o.PaymentMethod == order.PaymentCash {
					cashRevenue += o.Total
				} else if o.PaymentMethod == order.PaymentTransfer || o.PaymentMethod == order.PaymentQR {
					transferRevenue += o.Total
				}
			}
		}

		// Update shift
		now := time.Now()
		shift.Status = order.ShiftClosed
		shift.EndedAt = &now
		shift.TotalRevenue = totalRevenue
		shift.TotalOrders = totalOrders
		shift.EndCash = shift.CurrentCash // Use current cash as end cash
		shift.UpdatedAt = now

		// Update in database
		if err := shiftRepo.Update(ctx, shift.ID, shift); err != nil {
			fmt.Printf("   ❌ Failed to close shift: %v\n", err)
			continue
		}

		// Lock completed orders
		lockedCount := 0
		for _, o := range orders {
			if o.Status == order.StatusServed || o.Status == order.StatusCancelled {
				o.Status = order.StatusLocked
				o.LockedAt = &now
				if err := orderRepo.Update(ctx, o.ID, o); err != nil {
					fmt.Printf("   ⚠️  Warning: Failed to lock order %s: %v\n", o.OrderNumber, err)
				} else {
					lockedCount++
				}
			}
		}

		fmt.Printf("   ✅ Shift closed successfully\n")
		fmt.Printf("   📊 Summary:\n")
		fmt.Printf("      - Total Orders: %d\n", totalOrders)
		fmt.Printf("      - Total Revenue: %.0f VND\n", totalRevenue)
		fmt.Printf("      - Cash Revenue: %.0f VND\n", cashRevenue)
		fmt.Printf("      - Transfer Revenue: %.0f VND\n", transferRevenue)
		fmt.Printf("      - Current Cash: %.0f VND\n", shift.CurrentCash)
		fmt.Printf("      - Remaining Cash: %.0f VND\n", shift.RemainingCash)
		fmt.Printf("      - Remaining Transfer: %.0f VND\n", shift.RemainingTransfer)
		fmt.Printf("      - Locked Orders: %d\n", lockedCount)
		fmt.Println()
	}

	fmt.Printf("✅ Successfully closed %d shift(s)\n", len(openShifts))
}
