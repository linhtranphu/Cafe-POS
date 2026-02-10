package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"

	"go.mongodb.org/mongo-driver/bson"
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
	fmt.Println("  MENU COST & PROFIT ANALYSIS - MIGRATION VERIFICATION")
	fmt.Println(strings.Repeat("=", 60) + "\n")

	// Run all verification checks
	allPassed := true

	allPassed = verifyMenuItemsSchema(ctx, db) && allPassed
	allPassed = verifyIngredientsSchema(ctx, db) && allPassed
	allPassed = verifyOrderItemsCollection(ctx, db) && allPassed
	allPassed = verifyOperatingExpensesCollection(ctx, db) && allPassed
	allPassed = verifyShopSettingsSchema(ctx, db) && allPassed
	allPassed = verifyIndexes(ctx, db) && allPassed
	allPassed = verifyMenuItemCosts(ctx, db) && allPassed
	allPassed = verifyOrderItemCosts(ctx, db) && allPassed

	fmt.Println("\n" + strings.Repeat("=", 60))
	if allPassed {
		fmt.Println("  ✅ ALL VERIFICATION CHECKS PASSED")
	} else {
		fmt.Println("  ⚠️  SOME VERIFICATION CHECKS FAILED")
	}
	fmt.Println(strings.Repeat("=", 60) + "\n")

	if !allPassed {
		os.Exit(1)
	}
}

func verifyMenuItemsSchema(ctx context.Context, db *mongo.Database) bool {
	fmt.Println("📝 [1/8] Verifying menu_items schema...")
	collection := db.Collection("menu_items")

	// Check if all menu items have the new fields
	totalCount, _ := collection.CountDocuments(ctx, bson.M{})
	withCostFields, _ := collection.CountDocuments(ctx, bson.M{
		"current_cost":            bson.M{"$exists": true},
		"cost_last_calculated_at": bson.M{"$exists": true},
		"cost_status":             bson.M{"$exists": true},
	})

	if totalCount == withCostFields {
		fmt.Printf("   ✅ All %d menu items have cost tracking fields\n", totalCount)
		return true
	} else {
		fmt.Printf("   ❌ Only %d/%d menu items have cost tracking fields\n", withCostFields, totalCount)
		return false
	}
}

func verifyIngredientsSchema(ctx context.Context, db *mongo.Database) bool {
	fmt.Println("\n📝 [2/8] Verifying ingredients schema...")
	collection := db.Collection("ingredients")

	// Check if all ingredients have the new fields
	totalCount, _ := collection.CountDocuments(ctx, bson.M{})
	withConversionFields, _ := collection.CountDocuments(ctx, bson.M{
		"conversion_rate":    bson.M{"$exists": true},
		"wastage_percentage": bson.M{"$exists": true},
	})

	if totalCount == withConversionFields {
		fmt.Printf("   ✅ All %d ingredients have conversion and wastage fields\n", totalCount)
		return true
	} else {
		fmt.Printf("   ❌ Only %d/%d ingredients have conversion and wastage fields\n", withConversionFields, totalCount)
		return false
	}
}

func verifyOrderItemsCollection(ctx context.Context, db *mongo.Database) bool {
	fmt.Println("\n📝 [3/8] Verifying order_items collection...")

	// Check if collection exists
	collections, err := db.ListCollectionNames(ctx, options.ListCollections())
	if err != nil {
		fmt.Printf("   ❌ Failed to list collections: %v\n", err)
		return false
	}

	collectionExists := false
	for _, name := range collections {
		if name == "order_items" {
			collectionExists = true
			break
		}
	}

	if !collectionExists {
		fmt.Println("   ❌ order_items collection does not exist")
		return false
	}

	collection := db.Collection("order_items")
	count, _ := collection.CountDocuments(ctx, bson.M{})
	fmt.Printf("   ✅ order_items collection exists with %d documents\n", count)
	return true
}

func verifyOperatingExpensesCollection(ctx context.Context, db *mongo.Database) bool {
	fmt.Println("\n📝 [4/8] Verifying operating_expenses collection...")

	// Check if collection exists
	collections, err := db.ListCollectionNames(ctx, options.ListCollections())
	if err != nil {
		fmt.Printf("   ❌ Failed to list collections: %v\n", err)
		return false
	}

	collectionExists := false
	for _, name := range collections {
		if name == "operating_expenses" {
			collectionExists = true
			break
		}
	}

	if !collectionExists {
		fmt.Println("   ❌ operating_expenses collection does not exist")
		return false
	}

	collection := db.Collection("operating_expenses")
	count, _ := collection.CountDocuments(ctx, bson.M{})
	fmt.Printf("   ✅ operating_expenses collection exists with %d documents\n", count)
	return true
}

func verifyShopSettingsSchema(ctx context.Context, db *mongo.Database) bool {
	fmt.Println("\n📝 [5/8] Verifying shop_settings schema...")
	collection := db.Collection("shop_settings")

	// Check if shop settings have the new field
	totalCount, _ := collection.CountDocuments(ctx, bson.M{})
	withThreshold, _ := collection.CountDocuments(ctx, bson.M{
		"low_margin_threshold": bson.M{"$exists": true},
	})

	if totalCount == withThreshold {
		fmt.Printf("   ✅ All %d shop settings have low_margin_threshold field\n", totalCount)
		return true
	} else {
		fmt.Printf("   ❌ Only %d/%d shop settings have low_margin_threshold field\n", withThreshold, totalCount)
		return false
	}
}

func verifyIndexes(ctx context.Context, db *mongo.Database) bool {
	fmt.Println("\n📝 [6/8] Verifying indexes...")
	allPassed := true

	// Verify menu_items indexes
	menuIndexes := []string{"idx_category", "idx_cost_status", "idx_current_cost"}
	if !verifyCollectionIndexes(ctx, db, "menu_items", menuIndexes) {
		allPassed = false
	}

	// Verify order_items indexes
	orderItemIndexes := []string{"idx_order_id", "idx_menu_item_id", "idx_cost_status", "idx_cost_calculated_at", "idx_order_menu_item"}
	if !verifyCollectionIndexes(ctx, db, "order_items", orderItemIndexes) {
		allPassed = false
	}

	// Verify operating_expenses indexes
	expenseIndexes := []string{"idx_period_range", "idx_period_start"}
	if !verifyCollectionIndexes(ctx, db, "operating_expenses", expenseIndexes) {
		allPassed = false
	}

	if allPassed {
		fmt.Println("   ✅ All required indexes exist")
	}
	return allPassed
}

func verifyCollectionIndexes(ctx context.Context, db *mongo.Database, collectionName string, expectedIndexes []string) bool {
	collection := db.Collection(collectionName)
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		fmt.Printf("   ⚠️  Failed to list indexes for %s: %v\n", collectionName, err)
		return false
	}
	defer cursor.Close(ctx)

	var indexes []bson.M
	if err := cursor.All(ctx, &indexes); err != nil {
		fmt.Printf("   ⚠️  Failed to decode indexes for %s: %v\n", collectionName, err)
		return false
	}

	indexNames := make(map[string]bool)
	for _, index := range indexes {
		if name, ok := index["name"].(string); ok {
			indexNames[name] = true
		}
	}

	allFound := true
	for _, expectedIndex := range expectedIndexes {
		if !indexNames[expectedIndex] {
			fmt.Printf("   ⚠️  Missing index '%s' on %s\n", expectedIndex, collectionName)
			allFound = false
		}
	}

	return allFound
}

func verifyMenuItemCosts(ctx context.Context, db *mongo.Database) bool {
	fmt.Println("\n📝 [7/8] Verifying menu item costs...")
	collection := db.Collection("menu_items")

	// Count items by cost_status
	totalCount, _ := collection.CountDocuments(ctx, bson.M{})
	finalCount, _ := collection.CountDocuments(ctx, bson.M{"cost_status": menu.CostStatusFinal})
	incompleteCount, _ := collection.CountDocuments(ctx, bson.M{"cost_status": menu.CostStatusIncomplete})

	fmt.Printf("   Total menu items:       %d\n", totalCount)
	fmt.Printf("   FINAL status:           %d\n", finalCount)
	fmt.Printf("   INCOMPLETE status:      %d\n", incompleteCount)

	// Check if all items have a cost_status
	if totalCount == finalCount+incompleteCount {
		fmt.Println("   ✅ All menu items have valid cost_status")
		
		// Sample a few items to verify cost calculation
		cursor, err := collection.Find(ctx, bson.M{"cost_status": menu.CostStatusFinal}, 
			options.Find().SetLimit(3))
		if err == nil {
			defer cursor.Close(ctx)
			var items []struct {
				Name        string  `bson:"name"`
				CurrentCost float64 `bson:"current_cost"`
			}
			if err := cursor.All(ctx, &items); err == nil && len(items) > 0 {
				fmt.Println("   Sample items with calculated costs:")
				for _, item := range items {
					fmt.Printf("      - %s: %.2f VND\n", item.Name, item.CurrentCost)
				}
			}
		}
		
		return true
	} else {
		fmt.Printf("   ❌ Some menu items have invalid cost_status\n")
		return false
	}
}

func verifyOrderItemCosts(ctx context.Context, db *mongo.Database) bool {
	fmt.Println("\n📝 [8/8] Verifying order item costs...")
	collection := db.Collection("order_items")

	// Count items by cost_status
	totalCount, _ := collection.CountDocuments(ctx, bson.M{})
	
	if totalCount == 0 {
		fmt.Println("   ℹ️  No order items found (no closed shifts or backfill not run)")
		return true
	}

	finalCount, _ := collection.CountDocuments(ctx, bson.M{"cost_status": order.CostStatusFinal})
	estimatedCount, _ := collection.CountDocuments(ctx, bson.M{"cost_status": order.CostStatusEstimated})
	incompleteCount, _ := collection.CountDocuments(ctx, bson.M{"cost_status": order.CostStatusIncomplete})

	fmt.Printf("   Total order items:      %d\n", totalCount)
	fmt.Printf("   FINAL status:           %d\n", finalCount)
	fmt.Printf("   ESTIMATED status:       %d\n", estimatedCount)
	fmt.Printf("   INCOMPLETE status:      %d\n", incompleteCount)

	// Check if all items have a cost_status
	if totalCount == finalCount+estimatedCount+incompleteCount {
		fmt.Println("   ✅ All order items have valid cost_status")
		
		// Sample a few items to verify cost calculation
		cursor, err := collection.Find(ctx, bson.M{}, 
			options.Find().SetLimit(3))
		if err == nil {
			defer cursor.Close(ctx)
			var items []struct {
				Name           string  `bson:"name"`
				AccountingCost float64 `bson:"accounting_cost"`
				CostStatus     string  `bson:"cost_status"`
			}
			if err := cursor.All(ctx, &items); err == nil && len(items) > 0 {
				fmt.Println("   Sample order items with calculated costs:")
				for _, item := range items {
					fmt.Printf("      - %s: %.2f VND (%s)\n", item.Name, item.AccountingCost, item.CostStatus)
				}
			}
		}
		
		return true
	} else {
		fmt.Printf("   ❌ Some order items have invalid cost_status\n")
		return false
	}
}
