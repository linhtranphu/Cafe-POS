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
	fmt.Println("  BACKFILL CURRENT_COST FOR MENU ITEMS")
	fmt.Println(strings.Repeat("=", 60) + "\n")

	// Initialize repositories and services
	menuRepo := mongodb.NewMenuRepository(db)
	ingredientRepo := mongodb.NewIngredientRepository(db)
	orderRepo := mongodb.NewOrderRepository(db)
	orderItemRepo := mongodb.NewOrderItemRepository(db)
	costCalculator := services.NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	// Get all menu items
	fmt.Println("📝 Fetching all menu items...")
	menuItems, err := menuRepo.FindAll(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to fetch menu items: %v", err)
	}
	fmt.Printf("   Found %d menu items\n\n", len(menuItems))

	// Statistics
	stats := struct {
		Total      int
		Calculated int
		Incomplete int
		NoIngredients int
		Errors     int
	}{}

	// Process each menu item
	fmt.Println("🔄 Calculating costs for menu items...")
	for i, item := range menuItems {
		stats.Total++
		
		// Calculate cost
		costResult, err := costCalculator.CalculateMenuItemCost(ctx, item.ID)
		if err != nil {
			log.Printf("   ⚠️  [%d/%d] Error calculating cost for '%s': %v", 
				i+1, len(menuItems), item.Name, err)
			stats.Errors++
			continue
		}

		// Update menu item with calculated cost
		update := bson.M{
			"$set": bson.M{
				"current_cost":            costResult.CurrentCost,
				"cost_last_calculated_at": costResult.CostLastCalculatedAt,
				"cost_status":             costResult.CostStatus,
			},
		}

		_, err = db.Collection("menu_items").UpdateOne(
			ctx,
			bson.M{"_id": item.ID},
			update,
		)
		if err != nil {
			log.Printf("   ⚠️  [%d/%d] Error updating '%s': %v", 
				i+1, len(menuItems), item.Name, err)
			stats.Errors++
			continue
		}

		// Update statistics
		switch costResult.CostStatus {
		case menu.CostStatusFinal:
			if len(item.Ingredients) == 0 {
				stats.NoIngredients++
			} else {
				stats.Calculated++
			}
		case menu.CostStatusIncomplete:
			stats.Incomplete++
		}

		// Log progress every 10 items
		if (i+1)%10 == 0 || i+1 == len(menuItems) {
			fmt.Printf("   Progress: %d/%d items processed\n", i+1, len(menuItems))
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("  ✅ BACKFILL COMPLETED")
	fmt.Println(strings.Repeat("=", 60) + "\n")

	// Display summary
	displaySummary(ctx, db, stats)
}

func displaySummary(ctx context.Context, db *mongo.Database, stats struct {
	Total         int
	Calculated    int
	Incomplete    int
	NoIngredients int
	Errors        int
}) {
	fmt.Println("📊 Backfill Summary:")
	fmt.Println()
	fmt.Printf("   Total menu items:           %d\n", stats.Total)
	fmt.Printf("   ✅ Successfully calculated:  %d\n", stats.Calculated)
	fmt.Printf("   📦 No ingredients (cost=0):  %d\n", stats.NoIngredients)
	fmt.Printf("   ⚠️  Incomplete (missing cost): %d\n", stats.Incomplete)
	fmt.Printf("   ❌ Errors:                   %d\n", stats.Errors)
	fmt.Println()

	// Query database for verification
	collection := db.Collection("menu_items")
	
	finalCount, _ := collection.CountDocuments(ctx, bson.M{"cost_status": menu.CostStatusFinal})
	incompleteCount, _ := collection.CountDocuments(ctx, bson.M{"cost_status": menu.CostStatusIncomplete})
	
	fmt.Println("📈 Database Verification:")
	fmt.Printf("   FINAL status:       %d items\n", finalCount)
	fmt.Printf("   INCOMPLETE status:  %d items\n", incompleteCount)
	fmt.Println()

	if stats.Incomplete > 0 {
		fmt.Println("⚠️  Action Required:")
		fmt.Printf("   %d menu items have INCOMPLETE status due to missing ingredient costs.\n", stats.Incomplete)
		fmt.Println("   Please update ingredient cost_per_unit values and re-run this script.")
		fmt.Println()
		
		// Show some examples of incomplete items
		fmt.Println("   Examples of incomplete items:")
		cursor, err := collection.Find(ctx, bson.M{"cost_status": menu.CostStatusIncomplete}, 
			options.Find().SetLimit(5).SetProjection(bson.M{"name": 1, "category": 1}))
		if err == nil {
			defer cursor.Close(ctx)
			var items []struct {
				ID       primitive.ObjectID `bson:"_id"`
				Name     string             `bson:"name"`
				Category string             `bson:"category"`
			}
			if err := cursor.All(ctx, &items); err == nil {
				for _, item := range items {
					fmt.Printf("      - %s (%s)\n", item.Name, item.Category)
				}
			}
		}
		fmt.Println()
	}

	fmt.Println("📝 Next Steps:")
	if stats.Incomplete > 0 {
		fmt.Println("   1. Update missing ingredient costs in the system")
		fmt.Println("   2. Re-run this backfill script to recalculate incomplete items")
		fmt.Println("   3. Run task 19.3 to backfill accounting_cost for historical orders")
	} else {
		fmt.Println("   1. Run task 19.3 to backfill accounting_cost for historical orders")
		fmt.Println("   2. Run task 19.4 to verify migration completeness")
	}
	fmt.Println()
}
