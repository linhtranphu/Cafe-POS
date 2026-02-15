package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cafe-pos/backend/infrastructure/mongodb"

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

	log.Printf("Connecting to MongoDB: %s", mongoURI)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	db := client.Database(dbName)
	fmt.Printf("Connected to database: %s\n", dbName)

	// List of collections to create
	collections := []string{
		"batch_definitions",
		"batch_records",
		"batch_usage_logs",
	}

	// Create collections
	fmt.Println("\n=== Creating batch collections ===")
	
	existingCollections, err := db.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		log.Fatalf("Failed to list collections: %v", err)
	}

	existingMap := make(map[string]bool)
	for _, name := range existingCollections {
		existingMap[name] = true
	}

	for _, collName := range collections {
		if existingMap[collName] {
			fmt.Printf("✓ Collection '%s' already exists\n", collName)
		} else {
			if err := db.CreateCollection(ctx, collName); err != nil {
				log.Fatalf("Failed to create collection '%s': %v", collName, err)
			}
			fmt.Printf("✓ Created collection '%s'\n", collName)
		}
	}

	// Create indexes for batch_definitions
	fmt.Println("\n=== Creating indexes for batch_definitions ===")
	batchDefRepo := mongodb.NewBatchDefinitionRepository(db)
	if err := batchDefRepo.CreateIndexes(ctx); err != nil {
		log.Fatalf("Failed to create indexes for batch_definitions: %v", err)
	}
	fmt.Println("✓ Created indexes:")
	fmt.Println("  - name_idx (name)")
	fmt.Println("  - created_at_idx (created_at DESC)")

	// Create indexes for batch_records
	fmt.Println("\n=== Creating indexes for batch_records ===")
	batchRecordRepo := mongodb.NewBatchRecordRepository(db)
	if err := batchRecordRepo.CreateIndexes(ctx); err != nil {
		log.Fatalf("Failed to create indexes for batch_records: %v", err)
	}
	fmt.Println("✓ Created indexes:")
	fmt.Println("  - batch_def_expires_idx (batch_definition_id, expires_at)")
	fmt.Println("  - status_expires_idx (status, expires_at)")
	fmt.Println("  - expires_at_idx (expires_at)")
	fmt.Println("  - prepared_at_idx (prepared_at DESC)")
	fmt.Println("  - prepared_by_date_idx (prepared_by, prepared_at DESC)")

	// Create indexes for batch_usage_logs
	fmt.Println("\n=== Creating indexes for batch_usage_logs ===")
	batchUsageLogRepo := mongodb.NewBatchUsageLogRepository(db)
	if err := batchUsageLogRepo.CreateIndexes(ctx); err != nil {
		log.Fatalf("Failed to create indexes for batch_usage_logs: %v", err)
	}
	fmt.Println("✓ Created indexes:")
	fmt.Println("  - batch_record_used_idx (batch_record_id, used_at DESC)")
	fmt.Println("  - order_id_idx (order_id)")
	fmt.Println("  - menu_item_used_idx (menu_item_id, used_at DESC)")
	fmt.Println("  - used_at_idx (used_at DESC)")

	fmt.Println("\n=== Migration completed successfully ===")
	
	// Print schema information
	fmt.Println("\n=== Collection Schemas ===")
	
	fmt.Println("\nbatch_definitions:")
	fmt.Println("  - _id: ObjectId")
	fmt.Println("  - name: String")
	fmt.Println("  - unit: String")
	fmt.Println("  - shelf_life_hours: Number")
	fmt.Println("  - conversion_rates: Array")
	fmt.Println("    - source_ingredient_id: ObjectId")
	fmt.Println("    - source_ingredient_name: String")
	fmt.Println("    - source_quantity: Number")
	fmt.Println("    - source_unit: String")
	fmt.Println("    - batch_quantity: Number")
	fmt.Println("    - wastage_rate: Number")
	fmt.Println("  - low_stock_threshold: Number")
	fmt.Println("  - expiry_warning_hours: Number")
	fmt.Println("  - created_at: Date")
	fmt.Println("  - updated_at: Date")

	fmt.Println("\nbatch_records:")
	fmt.Println("  - _id: ObjectId")
	fmt.Println("  - batch_definition_id: ObjectId")
	fmt.Println("  - batch_name: String")
	fmt.Println("  - quantity_produced: Number")
	fmt.Println("  - quantity_remaining: Number")
	fmt.Println("  - unit: String")
	fmt.Println("  - cost_per_unit: Number")
	fmt.Println("  - total_cost: Number")
	fmt.Println("  - prepared_by: String")
	fmt.Println("  - prepared_at: Date")
	fmt.Println("  - expires_at: Date")
	fmt.Println("  - status: String (available|expired|depleted)")
	fmt.Println("  - ingredients_used: Array")
	fmt.Println("    - ingredient_id: ObjectId")
	fmt.Println("    - ingredient_name: String")
	fmt.Println("    - quantity: Number")
	fmt.Println("    - unit: String")
	fmt.Println("    - cost_per_unit: Number")
	fmt.Println("    - total_cost: Number")
	fmt.Println("  - created_at: Date")
	fmt.Println("  - updated_at: Date")

	fmt.Println("\nbatch_usage_logs:")
	fmt.Println("  - _id: ObjectId")
	fmt.Println("  - batch_record_id: ObjectId")
	fmt.Println("  - batch_name: String")
	fmt.Println("  - order_id: ObjectId")
	fmt.Println("  - menu_item_id: ObjectId")
	fmt.Println("  - menu_item_name: String")
	fmt.Println("  - quantity_used: Number")
	fmt.Println("  - unit: String")
	fmt.Println("  - cost_per_unit: Number")
	fmt.Println("  - total_cost: Number")
	fmt.Println("  - used_at: Date")

	fmt.Println("\n=== Usage ===")
	fmt.Println("To run this migration:")
	fmt.Println("  go run backend/cmd/migrate/create_batch_collections.go")
	fmt.Println("\nOr with custom MongoDB URI:")
	fmt.Println("  MONGODB_URI=mongodb://localhost:27017 go run backend/cmd/migrate/create_batch_collections.go")
}
