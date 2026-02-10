package main

import (
	"context"
	"log"
	"os"
	"time"

	"cafe-pos/backend/infrastructure/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MongoDB connection
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	log.Printf("Connecting to MongoDB: %s", mongoURI)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("❌ MongoDB connection failed: %v", err)
	}
	defer client.Disconnect(context.TODO())

	// Verify MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ MongoDB ping failed: %v", err)
	}
	log.Println("✅ MongoDB connected successfully")

	db := client.Database("cafe_pos")

	// Step 1: Create order_items collection if it doesn't exist
	log.Println("📝 Creating order_items collection...")
	collections, err := db.ListCollectionNames(ctx, options.ListCollections())
	if err != nil {
		log.Fatalf("❌ Failed to list collections: %v", err)
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
			log.Fatalf("❌ Failed to create order_items collection: %v", err)
		}
		log.Println("✅ order_items collection created")
	} else {
		log.Println("ℹ️  order_items collection already exists")
	}

	// Step 2: Create indexes
	log.Println("📝 Creating indexes for order_items collection...")
	orderItemRepo := mongodb.NewOrderItemRepository(db)
	if err := orderItemRepo.CreateIndexes(ctx); err != nil {
		log.Fatalf("❌ Failed to create indexes: %v", err)
	}
	log.Println("✅ Indexes created successfully")

	// Step 3: Display collection info
	log.Println("\n📊 Collection Information:")
	log.Println("Collection: order_items")
	log.Println("Indexes:")
	log.Println("  - idx_order_id: order_id")
	log.Println("  - idx_menu_item_id: menu_item_id")
	log.Println("  - idx_cost_status: cost_status")
	log.Println("  - idx_cost_calculated_at: cost_calculated_at")
	log.Println("  - idx_order_menu_item: order_id + menu_item_id (compound)")

	log.Println("\n🎉 Migration completed successfully!")
	log.Println("\nℹ️  Note: This collection will be populated when shifts are closed.")
	log.Println("   Order items will have accounting_cost calculated at shift closure time.")
}
