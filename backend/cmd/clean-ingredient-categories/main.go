package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("cafe_pos")

	fmt.Println("🧹 Cleaning all ingredient category data...")
	fmt.Println()

	// Delete all ingredient categories
	categoriesResult, err := db.Collection("ingredient_categories").DeleteMany(ctx, map[string]interface{}{})
	if err != nil {
		log.Fatal("Failed to delete ingredient categories:", err)
	}
	fmt.Printf("🗑️  Deleted %d ingredient categories\n", categoriesResult.DeletedCount)

	fmt.Println()
	fmt.Println("✅ All ingredient category data cleaned successfully!")
}
