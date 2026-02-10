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

	// Check if collection exists
	log.Println("\n📊 Verifying order_items collection...")
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
		log.Println("❌ order_items collection does not exist")
		return
	}
	log.Println("✅ order_items collection exists")

	// Check indexes
	log.Println("\n📊 Verifying indexes...")
	collection := db.Collection("order_items")
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to list indexes: %v", err)
	}
	defer cursor.Close(ctx)

	var indexes []bson.M
	if err = cursor.All(ctx, &indexes); err != nil {
		log.Fatalf("❌ Failed to decode indexes: %v", err)
	}

	log.Printf("Found %d indexes:\n", len(indexes))
	for _, index := range indexes {
		name := index["name"]
		keys := index["key"]
		fmt.Printf("  - %s: %v\n", name, keys)
	}

	// Check document count
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatalf("❌ Failed to count documents: %v", err)
	}
	log.Printf("\n📊 Document count: %d\n", count)

	log.Println("\n✅ Verification completed successfully!")
}
