package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/infrastructure/mongodb"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGO_URI")
	}
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
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

	fmt.Println("\n📋 List of Open Shifts")
	fmt.Println("======================\n")

	// Find all open shifts
	openShifts, err := shiftRepo.FindOpenShifts(ctx)
	if err != nil {
		log.Fatalf("Failed to find open shifts: %v", err)
	}

	if len(openShifts) == 0 {
		fmt.Println("✅ No open shifts found. All shifts are closed.")
		return
	}

	fmt.Printf("Found %d open shift(s):\n\n", len(openShifts))

	// Display each shift
	for i, shift := range openShifts {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("[%d] Shift ID: %s\n", i+1, shift.ID.Hex())
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("👤 User: %s (ID: %s)\n", shift.UserName, shift.UserID.Hex())
		fmt.Printf("🏷️  Role: %s\n", shift.RoleType)
		fmt.Printf("⏰ Type: %s\n", shift.Type)
		fmt.Printf("📅 Started: %s\n", shift.StartedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("📊 Status: %s\n", shift.Status)
		fmt.Println()

		// Cash information
		fmt.Printf("💰 Cash Information:\n")
		fmt.Printf("   Start Cash: %.0f VND\n", shift.StartCash)
		fmt.Printf("   Current Cash: %.0f VND\n", shift.CurrentCash)
		fmt.Printf("   Remaining Cash: %.0f VND\n", shift.RemainingCash)
		fmt.Printf("   Handed Over Cash: %.0f VND\n", shift.HandedOverCash)
		fmt.Println()

		// Transfer information
		fmt.Printf("💳 Transfer Information:\n")
		fmt.Printf("   Transfer Revenue: %.0f VND\n", shift.TransferRevenue)
		fmt.Printf("   Remaining Transfer: %.0f VND\n", shift.RemainingTransfer)
		fmt.Printf("   Handed Over Transfer: %.0f VND\n", shift.HandedOverTransfer)
		fmt.Println()

		// Get orders for this shift
		orders, err := orderRepo.FindByShiftID(ctx, shift.ID)
		if err != nil {
			fmt.Printf("   ⚠️  Warning: Failed to get orders: %v\n", err)
		} else {
			totalRevenue := 0.0
			cashRevenue := 0.0
			transferRevenue := 0.0
			paidOrders := 0

			for _, o := range orders {
				if o.Status == order.StatusPaid || o.Status == order.StatusInProgress || o.Status == order.StatusServed {
					totalRevenue += o.Total
					paidOrders++

					if o.PaymentMethod == order.PaymentCash {
						cashRevenue += o.Total
					} else if o.PaymentMethod == order.PaymentTransfer || o.PaymentMethod == order.PaymentQR {
						transferRevenue += o.Total
					}
				}
			}

			fmt.Printf("📦 Orders:\n")
			fmt.Printf("   Total Orders: %d\n", len(orders))
			fmt.Printf("   Paid Orders: %d\n", paidOrders)
			fmt.Printf("   Total Revenue: %.0f VND\n", totalRevenue)
			fmt.Printf("   Cash Revenue: %.0f VND\n", cashRevenue)
			fmt.Printf("   Transfer Revenue: %.0f VND\n", transferRevenue)
		}

		fmt.Println()
	}

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Total: %d open shift(s)\n", len(openShifts))
}
