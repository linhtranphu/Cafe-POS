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

	// Create operating_expenses collection
	fmt.Println("\n=== Creating operating_expenses collection ===")
	
	// Check if collection already exists
	collections, err := db.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		log.Fatalf("Failed to list collections: %v", err)
	}

	collectionExists := false
	for _, name := range collections {
		if name == "operating_expenses" {
			collectionExists = true
			break
		}
	}

	if collectionExists {
		fmt.Println("✓ Collection 'operating_expenses' already exists")
	} else {
		if err := db.CreateCollection(ctx, "operating_expenses"); err != nil {
			log.Fatalf("Failed to create collection: %v", err)
		}
		fmt.Println("✓ Created collection 'operating_expenses'")
	}

	// Create indexes
	fmt.Println("\n=== Creating indexes ===")
	repo := mongodb.NewOperatingExpenseRepository(db)
	if err := repo.CreateIndexes(ctx); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}
	fmt.Println("✓ Created indexes:")
	fmt.Println("  - period_range_idx (period_start, period_end)")
	fmt.Println("  - period_start_idx (period_start)")

	fmt.Println("\n=== Migration completed successfully ===")
	fmt.Println("\nCollection schema:")
	fmt.Println("  - _id: ObjectId")
	fmt.Println("  - period_start: Date")
	fmt.Println("  - period_end: Date")
	fmt.Println("  - staff_salary: Number")
	fmt.Println("  - rent: Number")
	fmt.Println("  - utilities: Number")
	fmt.Println("  - marketing_costs: Number")
	fmt.Println("  - other_expenses: Number")
	fmt.Println("  - total_expenses: Number (auto-calculated)")
	fmt.Println("  - created_at: Date")
	fmt.Println("  - updated_at: Date")
}
