package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get MongoDB URI from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(ctx)

	// Get database
	db := client.Database("cafe_pos")
	collection := db.Collection("shop_settings")

	fmt.Println("🔧 Removing paper_width and label_size from shop_settings...")
	fmt.Println("These fields should be configured per printer, not globally.")

	// Remove paper_width and label_size fields from all documents
	update := bson.M{
		"$unset": bson.M{
			"paper_width": "",
			"label_size":  "",
		},
	}

	result, err := collection.UpdateMany(ctx, bson.M{}, update)
	if err != nil {
		log.Fatal("Failed to update shop_settings:", err)
	}

	fmt.Printf("✅ Updated %d document(s)\n", result.ModifiedCount)
	fmt.Println("\n📝 Note: Paper width is now configured per printer in printer_configs collection")
	fmt.Println("   Each printer (BILL/LABEL) has its own paper_width setting")
}
