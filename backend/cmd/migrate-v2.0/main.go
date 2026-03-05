package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cafe-pos/backend/domain/settings"
	"cafe-pos/backend/infrastructure/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("🚀 Cafe POS - Migration to v2.0")
	fmt.Println("========================================")
	fmt.Println()

	// Get MongoDB URI from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin"
		log.Println("⚠️  Using default MongoDB URI")
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ MongoDB ping failed: %v", err)
	}
	log.Println("✅ Connected to MongoDB")

	// Get database
	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "cafe_pos"
	}
	db := client.Database(dbName)
	log.Printf("📦 Database: %s\n", dbName)
	fmt.Println()

	ctx = context.Background()

	// Run migrations
	fmt.Println("🔄 Running migrations...")
	fmt.Println()

	// Migration 1: Create shop_settings if not exists
	if err := migrateShopSettings(ctx, db); err != nil {
		log.Fatalf("❌ Shop settings migration failed: %v", err)
	}

	// Migration 2: Add indexes for new collections
	if err := createIndexes(ctx, db); err != nil {
		log.Fatalf("❌ Index creation failed: %v", err)
	}

	// Migration 3: Update existing orders with new fields (if needed)
	if err := migrateOrders(ctx, db); err != nil {
		log.Fatalf("❌ Orders migration failed: %v", err)
	}

	// Migration 4: Ensure print collections exist
	if err := ensurePrintCollections(ctx, db); err != nil {
		log.Fatalf("❌ Print collections migration failed: %v", err)
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("✅ Migration to v2.0 completed successfully!")
	fmt.Println("========================================")
}

// migrateShopSettings creates default shop settings if not exists
func migrateShopSettings(ctx context.Context, db *mongo.Database) error {
	fmt.Println("1️⃣  Migrating shop_settings...")

	settingsRepo := mongodb.NewShopSettingsRepository(db)

	// Check if settings already exist
	existing, err := settingsRepo.GetSettings(ctx)
	if err == nil && existing != nil {
		fmt.Printf("   ✅ Shop settings already exist (ID: %s)\n", existing.ID.Hex())
		return nil
	}

	// Create default shop settings
	defaultSettings := settings.NewShopSettings("Cafe POS")
	defaultSettings.ShopAddress = "123 Main Street"
	defaultSettings.ShopPhone = "0123-456-789"
	defaultSettings.CustomMessage = "Cảm ơn quý khách! Hẹn gặp lại!"
	defaultSettings.SetAutoPrintEnabled(true)
	defaultSettings.SetFieldVisibility(false, true, true, true) // logo=false initially

	if err := settingsRepo.CreateSettings(ctx, defaultSettings); err != nil {
		return fmt.Errorf("failed to create shop settings: %w", err)
	}

	fmt.Printf("   ✅ Created default shop settings (ID: %s)\n", defaultSettings.ID.Hex())
	return nil
}

// createIndexes creates indexes for collections
func createIndexes(ctx context.Context, db *mongo.Database) error {
	fmt.Println("2️⃣  Creating indexes...")

	// Print jobs indexes
	printJobsCol := db.Collection("print_jobs")
	_, err := printJobsCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "status", Value: 1}, {Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "order_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "printer_id", Value: 1}},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create print_jobs indexes: %w", err)
	}
	fmt.Println("   ✅ Created print_jobs indexes")

	// Printer configs indexes
	printerConfigsCol := db.Collection("printer_configs")
	_, err = printerConfigsCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "type", Value: 1}, {Key: "is_default", Value: 1}},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create printer_configs indexes: %w", err)
	}
	fmt.Println("   ✅ Created printer_configs indexes")

	// Print templates indexes
	printTemplatesCol := db.Collection("print_templates")
	_, err = printTemplatesCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "type", Value: 1}, {Key: "is_default", Value: 1}},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create print_templates indexes: %w", err)
	}
	fmt.Println("   ✅ Created print_templates indexes")

	return nil
}

// migrateOrders updates existing orders with new fields if needed
func migrateOrders(ctx context.Context, db *mongo.Database) error {
	fmt.Println("3️⃣  Migrating orders...")

	ordersCol := db.Collection("orders")

	// Count orders that need migration (orders without print-related fields)
	// In v2.0, we don't add new fields to orders, so this is a no-op
	// But we keep it for future migrations

	count, err := ordersCol.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to count orders: %w", err)
	}

	fmt.Printf("   ✅ Verified %d orders (no changes needed)\n", count)
	return nil
}

// ensurePrintCollections ensures print-related collections exist
func ensurePrintCollections(ctx context.Context, db *mongo.Database) error {
	fmt.Println("4️⃣  Ensuring print collections exist...")

	collections := []string{
		"print_jobs",
		"printer_configs",
		"print_templates",
		"print_notifications",
	}

	existingCollections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	existingMap := make(map[string]bool)
	for _, name := range existingCollections {
		existingMap[name] = true
	}

	for _, colName := range collections {
		if existingMap[colName] {
			fmt.Printf("   ✅ Collection '%s' already exists\n", colName)
		} else {
			// Create collection
			if err := db.CreateCollection(ctx, colName); err != nil {
				return fmt.Errorf("failed to create collection %s: %w", colName, err)
			}
			fmt.Printf("   ✅ Created collection '%s'\n", colName)
		}
	}

	return nil
}
