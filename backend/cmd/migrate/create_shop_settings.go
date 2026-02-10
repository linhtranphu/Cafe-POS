package main

import (
	"context"
	"log"
	"os"
	"time"

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

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database("cafe_pos")
	collection := db.Collection("shop_settings")

	// Check if shop_settings already exists
	count, err := collection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		log.Fatalf("Failed to count shop_settings: %v", err)
	}

	if count > 0 {
		log.Println("✓ Shop settings already exists")
		return
	}

	// Create default shop settings
	shopSettings := bson.M{
		"shop_name":            "My Cafe",
		"low_margin_threshold": 20.0,
		"created_at":           time.Now(),
		"updated_at":           time.Now(),
	}

	result, err := collection.InsertOne(context.TODO(), shopSettings)
	if err != nil {
		log.Fatalf("Failed to create shop settings: %v", err)
	}

	log.Printf("✓ Created shop settings with ID: %v", result.InsertedID)
	log.Println("✓ Default low_margin_threshold: 20.0")
}
