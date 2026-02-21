package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("cafe_pos")
	collection := db.Collection("shop_settings")

	// Update all shop_settings documents that don't have the new print fields
	filter := bson.M{
		"$or": []bson.M{
			{"shop_address": bson.M{"$exists": false}},
			{"shop_phone": bson.M{"$exists": false}},
			{"paper_width": bson.M{"$exists": false}},
			{"label_size": bson.M{"$exists": false}},
		},
	}

	update := bson.M{
		"$set": bson.M{
			"shop_address":        "",
			"shop_phone":          "",
			"logo_url":            "",
			"custom_message":      "",
			"paper_width":         80,      // Default: 80mm
			"label_size":          "60x40", // Default: 60x40mm
			"show_logo":           true,
			"show_address":        true,
			"show_phone":          true,
			"show_custom_message": true,
			"updated_at":          time.Now(),
		},
	}

	result, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatalf("Failed to update shop_settings: %v", err)
	}

	fmt.Printf("✅ Migration completed successfully\n")
	fmt.Printf("   Updated %d shop_settings documents with print configuration fields\n", result.ModifiedCount)
}
