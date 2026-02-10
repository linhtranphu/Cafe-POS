package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/infrastructure/mongodb"

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

	// Run all migrations
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("  MENU COST & PROFIT ANALYSIS - SCHEMA MIGRATION")
	fmt.Println(strings.Repeat("=", 60) + "\n")

	// Migration 1: Add cost tracking fields to menu_items
	if err := migrateMenuItems(ctx, db); err != nil {
		log.Fatalf("❌ Menu items migration failed: %v", err)
	}

	// Migration 2: Add conversion and wastage fields to ingredients
	if err := migrateIngredients(ctx, db); err != nil {
		log.Fatalf("❌ Ingredients migration failed: %v", err)
	}

	// Migration 3: Create order_items collection
	if err := createOrderItemsCollection(ctx, db); err != nil {
		log.Fatalf("❌ Order items collection creation failed: %v", err)
	}

	// Migration 4: Create operating_expenses collection
	if err := createOperatingExpensesCollection(ctx, db); err != nil {
		log.Fatalf("❌ Operating expenses collection creation failed: %v", err)
	}

	// Migration 5: Add low_margin_threshold to shop_settings
	if err := migrateShopSettings(ctx, db); err != nil {
		log.Fatalf("❌ Shop settings migration failed: %v", err)
	}

	// Create all indexes
	if err := createAllIndexes(ctx, db); err != nil {
		log.Fatalf("❌ Index creation failed: %v", err)
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("  ✅ ALL MIGRATIONS COMPLETED SUCCESSFULLY")
	fmt.Println(strings.Repeat("=", 60) + "\n")

	// Display summary
	displayMigrationSummary(ctx, db)
}

func migrateMenuItems(ctx context.Context, db *mongo.Database) error {
	fmt.Println("📝 [1/5] Migrating menu_items collection...")
	collection := db.Collection("menu_items")

	// Update all menu items that don't have the new fields
	filter := bson.M{
		"$or": []bson.M{
			{"current_cost": bson.M{"$exists": false}},
			{"cost_last_calculated_at": bson.M{"$exists": false}},
			{"cost_status": bson.M{"$exists": false}},
		},
	}

	update := bson.M{
		"$set": bson.M{
			"current_cost":            0.0,
			"cost_last_calculated_at": time.Now(),
			"cost_status":             menu.CostStatusIncomplete,
		},
	}

	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update menu items: %w", err)
	}

	fmt.Printf("   ✅ Added cost tracking fields to %d menu items\n", result.ModifiedCount)
	fmt.Println("   Fields added: current_cost, cost_last_calculated_at, cost_status")
	return nil
}

func migrateIngredients(ctx context.Context, db *mongo.Database) error {
	fmt.Println("\n📝 [2/5] Migrating ingredients collection...")
	collection := db.Collection("ingredients")

	// Update all ingredients that don't have the new fields
	filter := bson.M{
		"$or": []bson.M{
			{"conversion_rate": bson.M{"$exists": false}},
			{"wastage_percentage": bson.M{"$exists": false}},
		},
	}

	update := bson.M{
		"$set": bson.M{
			"conversion_rate":    1.0,
			"wastage_percentage": 0.0,
		},
	}

	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update ingredients: %w", err)
	}

	fmt.Printf("   ✅ Added conversion and wastage fields to %d ingredients\n", result.ModifiedCount)
	fmt.Println("   Fields added: conversion_rate (default: 1.0), wastage_percentage (default: 0.0)")
	return nil
}

func createOrderItemsCollection(ctx context.Context, db *mongo.Database) error {
	fmt.Println("\n📝 [3/5] Creating order_items collection...")

	// Check if collection already exists
	collections, err := db.ListCollectionNames(ctx, options.ListCollections())
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	collectionExists := false
	for _, name := range collections {
		if name == "order_items" {
			collectionExists = true
			break
		}
	}

	if !collectionExists {
		err = db.CreateCollection(ctx, "order_items")
		if err != nil {
			return fmt.Errorf("failed to create order_items collection: %w", err)
		}
		fmt.Println("   ✅ Created order_items collection")
	} else {
		fmt.Println("   ℹ️  order_items collection already exists")
	}

	fmt.Println("   Schema: order_id, menu_item_id, name, price, quantity, note, subtotal,")
	fmt.Println("           accounting_cost, cost_calculated_at, cost_status, created_at")
	return nil
}

func createOperatingExpensesCollection(ctx context.Context, db *mongo.Database) error {
	fmt.Println("\n📝 [4/5] Creating operating_expenses collection...")

	// Check if collection already exists
	collections, err := db.ListCollectionNames(ctx, options.ListCollections())
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	collectionExists := false
	for _, name := range collections {
		if name == "operating_expenses" {
			collectionExists = true
			break
		}
	}

	if !collectionExists {
		err = db.CreateCollection(ctx, "operating_expenses")
		if err != nil {
			return fmt.Errorf("failed to create operating_expenses collection: %w", err)
		}
		fmt.Println("   ✅ Created operating_expenses collection")
	} else {
		fmt.Println("   ℹ️  operating_expenses collection already exists")
	}

	fmt.Println("   Schema: period_start, period_end, staff_salary, rent, utilities,")
	fmt.Println("           marketing_costs, other_expenses, total_expenses, created_at, updated_at")
	return nil
}

func migrateShopSettings(ctx context.Context, db *mongo.Database) error {
	fmt.Println("\n📝 [5/5] Migrating shop_settings collection...")
	collection := db.Collection("shop_settings")

	// Update all shop settings that don't have the new field
	filter := bson.M{
		"low_margin_threshold": bson.M{"$exists": false},
	}

	update := bson.M{
		"$set": bson.M{
			"low_margin_threshold": 20.0, // Default: 20%
		},
	}

	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update shop settings: %w", err)
	}

	fmt.Printf("   ✅ Added low_margin_threshold to %d shop settings\n", result.ModifiedCount)
	fmt.Println("   Field added: low_margin_threshold (default: 20.0)")
	return nil
}

func createAllIndexes(ctx context.Context, db *mongo.Database) error {
	fmt.Println("\n📝 Creating indexes for all collections...")

	// Menu items indexes
	fmt.Println("\n   Creating indexes for menu_items...")
	menuRepo := mongodb.NewMenuRepository(db)
	if err := menuRepo.CreateIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create menu_items indexes: %w", err)
	}
	fmt.Println("   ✅ menu_items indexes:")
	fmt.Println("      - idx_category: category")
	fmt.Println("      - idx_cost_status: cost_status")
	fmt.Println("      - idx_current_cost: current_cost")

	// Order items indexes
	fmt.Println("\n   Creating indexes for order_items...")
	orderItemRepo := mongodb.NewOrderItemRepository(db)
	if err := orderItemRepo.CreateIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create order_items indexes: %w", err)
	}
	fmt.Println("   ✅ order_items indexes:")
	fmt.Println("      - idx_order_id: order_id")
	fmt.Println("      - idx_menu_item_id: menu_item_id")
	fmt.Println("      - idx_cost_status: cost_status")
	fmt.Println("      - idx_cost_calculated_at: cost_calculated_at")
	fmt.Println("      - idx_order_menu_item: order_id + menu_item_id (compound)")

	// Operating expenses indexes
	fmt.Println("\n   Creating indexes for operating_expenses...")
	expenseRepo := mongodb.NewOperatingExpenseRepository(db)
	if err := expenseRepo.CreateIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create operating_expenses indexes: %w", err)
	}
	fmt.Println("   ✅ operating_expenses indexes:")
	fmt.Println("      - idx_period_range: period_start + period_end (compound)")
	fmt.Println("      - idx_period_start: period_start")

	return nil
}

func displayMigrationSummary(ctx context.Context, db *mongo.Database) {
	fmt.Println("📊 Migration Summary:")
	fmt.Println()

	// Count menu items
	menuCount, _ := db.Collection("menu_items").CountDocuments(ctx, bson.M{})
	fmt.Printf("   • menu_items: %d documents\n", menuCount)

	// Count ingredients
	ingredientCount, _ := db.Collection("ingredients").CountDocuments(ctx, bson.M{})
	fmt.Printf("   • ingredients: %d documents\n", ingredientCount)

	// Count order items
	orderItemCount, _ := db.Collection("order_items").CountDocuments(ctx, bson.M{})
	fmt.Printf("   • order_items: %d documents\n", orderItemCount)

	// Count operating expenses
	expenseCount, _ := db.Collection("operating_expenses").CountDocuments(ctx, bson.M{})
	fmt.Printf("   • operating_expenses: %d documents\n", expenseCount)

	// Count shop settings
	settingsCount, _ := db.Collection("shop_settings").CountDocuments(ctx, bson.M{})
	fmt.Printf("   • shop_settings: %d documents\n", settingsCount)

	fmt.Println()
	fmt.Println("📝 Next Steps:")
	fmt.Println("   1. Run task 19.2 to backfill current_cost for existing menu items")
	fmt.Println("   2. Run task 19.3 to backfill accounting_cost for historical orders")
	fmt.Println("   3. Run task 19.4 to verify migration completeness")
	fmt.Println()
}
