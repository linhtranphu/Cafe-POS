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

	fmt.Println("🧹 Cleaning all ingredient data...")
	fmt.Println()

	// Delete all ingredients
	ingredientsResult, err := db.Collection("ingredients").DeleteMany(ctx, map[string]interface{}{})
	if err != nil {
		log.Fatal("Failed to delete ingredients:", err)
	}
	fmt.Printf("🗑️  Deleted %d ingredients\n", ingredientsResult.DeletedCount)

	// Delete all stock history
	historyResult, err := db.Collection("stock_history").DeleteMany(ctx, map[string]interface{}{})
	if err != nil {
		log.Fatal("Failed to delete stock history:", err)
	}
	fmt.Printf("🗑️  Deleted %d stock history records\n", historyResult.DeletedCount)

	fmt.Println()
	fmt.Println("✅ All ingredient data cleaned successfully!")
}
