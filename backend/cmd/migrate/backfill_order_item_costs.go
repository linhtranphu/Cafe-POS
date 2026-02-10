package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/infrastructure/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "cafe_pos"
	}

	log.Printf("🔌 Connecting to MongoDB: %s", mongoURI)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ Failed to ping MongoDB: %v", err)
	}

	db := client.Database(dbName)
	log.Printf("✅ Connected to database: %s\n", dbName)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("  BACKFILL ACCOUNTING_COST FOR HISTORICAL ORDERS")
	fmt.Println(strings.Repeat("=", 60) + "\n")

	// Initialize repositories and services
	menuRepo := mongodb.NewMenuRepository(db)
	ingredientRepo := mongodb.NewIngredientRepository(db)
	orderRepo := mongodb.NewOrderRepository(db)
	orderItemRepo := mongodb.NewOrderItemRepository(db)
	costCalculator := services.NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	// Get all closed shifts
	fmt.Println("📝 Fetching closed shifts...")
	shifts, err := getClosedShifts(ctx, db)
	if err != nil {
		log.Fatalf("❌ Failed to fetch shifts: %v", err)
	}
	fmt.Printf("   Found %d closed shifts\n\n", len(shifts))

	if len(shifts) == 0 {
		fmt.Println("ℹ️  No closed shifts found. Nothing to backfill.")
		fmt.Println()
		return
	}

	// Statistics
	stats := struct {
		TotalShifts       int
		TotalOrders       int
		TotalItems        int
		CalculatedItems   int
		IncompleteItems   int
		Errors            int
		ShiftsProcessed   int
	}{}

	stats.TotalShifts = len(shifts)

	// Process each shift
	fmt.Println("🔄 Processing shifts and calculating costs...")
	for i, shift := range shifts {
		stats.ShiftsProcessed++

		// Get all orders for this shift
		orders, err := orderRepo.FindByShiftID(ctx, shift.ID)
		if err != nil {
			log.Printf("   ⚠️  [%d/%d] Error fetching orders for shift %s: %v", 
				i+1, len(shifts), shift.ID.Hex(), err)
			stats.Errors++
			continue
		}

		if len(orders) == 0 {
			continue
		}

		stats.TotalOrders += len(orders)

		// Calculate costs for all orders in this shift
		// Note: We mark these as ESTIMATED since they weren't calculated at actual shift closure time
		for _, ord := range orders {
			for _, item := range ord.Items {
				stats.TotalItems++

				// Calculate cost for this menu item
				costResult, err := costCalculator.CalculateMenuItemCost(ctx, item.MenuItemID)
				if err != nil {
					log.Printf("   ⚠️  Error calculating cost for item '%s' in order %s: %v", 
						item.Name, ord.ID.Hex(), err)
					stats.Errors++
					continue
				}

				// Create order item with accounting cost
				orderItem := &order.OrderItemWithCost{
					ID:                 primitive.NewObjectID(),
					OrderID:            ord.ID,
					MenuItemID:         item.MenuItemID,
					Name:               item.Name,
					Price:              item.Price,
					Quantity:           item.Quantity,
					Note:               item.Note,
					Subtotal:           item.Subtotal,
					AccountingCost:     costResult.CurrentCost * float64(item.Quantity),
					CostCalculatedAt:   time.Now(),
					CostStatus:         order.CostStatusEstimated, // Mark as ESTIMATED for backfilled data
					CreatedAt:          ord.CreatedAt,
				}

				// If cost is incomplete, mark it as such
				if costResult.CostStatus == menu.CostStatusIncomplete {
					orderItem.CostStatus = order.CostStatusIncomplete
					stats.IncompleteItems++
				} else {
					stats.CalculatedItems++
				}

				// Insert into order_items collection
				_, err = db.Collection("order_items").InsertOne(ctx, orderItem)
				if err != nil {
					log.Printf("   ⚠️  Error inserting order item: %v", err)
					stats.Errors++
					continue
				}
			}
		}

		// Log progress every 10 shifts
		if (i+1)%10 == 0 || i+1 == len(shifts) {
			fmt.Printf("   Progress: %d/%d shifts processed (%d orders, %d items)\n", 
				i+1, len(shifts), stats.TotalOrders, stats.TotalItems)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("  ✅ BACKFILL COMPLETED")
	fmt.Println(strings.Repeat("=", 60) + "\n")

	// Display summary
	displaySummary(ctx, db, stats)
}

func getClosedShifts(ctx context.Context, db *mongo.Database) ([]struct {
	ID     primitive.ObjectID `bson:"_id"`
	Status string             `bson:"status"`
}, error) {
	collection := db.Collection("shifts")
	
	// Find all closed shifts (both cashier and waiter shifts)
	filter := bson.M{
		"status": "CLOSED",
	}
	
	cursor, err := collection.Find(ctx, filter, options.Find().SetSort(bson.M{"created_at": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var shifts []struct {
		ID     primitive.ObjectID `bson:"_id"`
		Status string             `bson:"status"`
	}
	if err := cursor.All(ctx, &shifts); err != nil {
		return nil, err
	}
	
	return shifts, nil
}

func displaySummary(ctx context.Context, db *mongo.Database, stats struct {
	TotalShifts     int
	TotalOrders     int
	TotalItems      int
	CalculatedItems int
	IncompleteItems int
	Errors          int
	ShiftsProcessed int
}) {
	fmt.Println("📊 Backfill Summary:")
	fmt.Println()
	fmt.Printf("   Total closed shifts:        %d\n", stats.TotalShifts)
	fmt.Printf("   Shifts processed:           %d\n", stats.ShiftsProcessed)
	fmt.Printf("   Total orders:               %d\n", stats.TotalOrders)
	fmt.Printf("   Total order items:          %d\n", stats.TotalItems)
	fmt.Printf("   ✅ Successfully calculated:  %d\n", stats.CalculatedItems)
	fmt.Printf("   ⚠️  Incomplete (missing cost): %d\n", stats.IncompleteItems)
	fmt.Printf("   ❌ Errors:                   %d\n", stats.Errors)
	fmt.Println()

	// Query database for verification
	collection := db.Collection("order_items")
	
	totalCount, _ := collection.CountDocuments(ctx, bson.M{})
	estimatedCount, _ := collection.CountDocuments(ctx, bson.M{"cost_status": order.CostStatusEstimated})
	incompleteCount, _ := collection.CountDocuments(ctx, bson.M{"cost_status": order.CostStatusIncomplete})
	
	fmt.Println("📈 Database Verification:")
	fmt.Printf("   Total order items:      %d\n", totalCount)
	fmt.Printf("   ESTIMATED status:       %d items\n", estimatedCount)
	fmt.Printf("   INCOMPLETE status:      %d items\n", incompleteCount)
	fmt.Println()

	fmt.Println("ℹ️  Important Notes:")
	fmt.Println("   • All backfilled costs are marked as ESTIMATED (not FINAL)")
	fmt.Println("   • ESTIMATED means the cost was calculated using current ingredient prices,")
	fmt.Println("     not the actual prices at the time of shift closure")
	fmt.Println("   • Future shift closures will use FINAL status for accurate accounting")
	fmt.Println()

	if stats.IncompleteItems > 0 {
		fmt.Println("⚠️  Action Required:")
		fmt.Printf("   %d order items have INCOMPLETE status due to missing ingredient costs.\n", stats.IncompleteItems)
		fmt.Println("   These items will not be included in profit reports.")
		fmt.Println("   Please update ingredient cost_per_unit values if accurate historical")
		fmt.Println("   profit analysis is needed.")
		fmt.Println()
	}

	fmt.Println("📝 Next Steps:")
	fmt.Println("   1. Run task 19.4 to verify migration completeness")
	fmt.Println("   2. Test the profit analysis features in the manager interface")
	fmt.Println("   3. Future shift closures will automatically calculate FINAL costs")
	fmt.Println()
}
