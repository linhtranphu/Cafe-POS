package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Get MongoDB URI from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database("cafe_pos")
	collection := db.Collection("shop_settings")

	// Update all shop_settings documents that don't have auto_print_enabled field
	filter := bson.M{
		"auto_print_enabled": bson.M{"$exists": false},
	}
	update := bson.M{
		"$set": bson.M{
			"auto_print_enabled": true, // Default to enabled
		},
	}

	result, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatalf("Failed to update shop_settings: %v", err)
	}

	fmt.Printf("✅ Migration completed successfully\n")
	fmt.Printf("   Updated %d shop_settings documents with auto_print_enabled field\n", result.ModifiedCount)
}
