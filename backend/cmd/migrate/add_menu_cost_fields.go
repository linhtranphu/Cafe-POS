package main

import (
	"context"
	"log"
	"os"
	"time"

	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/infrastructure/mongodb"

	"go.mongodb.org/mongo-driver/bson"
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
	menuRepo := mongodb.NewMenuRepository(db)

	// Step 1: Add new fields to existing menu items
	log.Println("📝 Adding cost tracking fields to existing menu items...")
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
		log.Fatalf("❌ Failed to update menu items: %v", err)
	}
	log.Printf("✅ Updated %d menu items with cost tracking fields", result.ModifiedCount)

	// Step 2: Create indexes
	log.Println("📝 Creating indexes for menu_items collection...")
	if err := menuRepo.CreateIndexes(ctx); err != nil {
		log.Fatalf("❌ Failed to create indexes: %v", err)
	}
	log.Println("✅ Indexes created successfully")

	log.Println("🎉 Migration completed successfully!")
}
