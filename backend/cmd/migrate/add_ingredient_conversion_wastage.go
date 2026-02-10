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
	
	// Replace 'mongodb' hostname with 'localhost' for local execution
	// This handles the case where MONGODB_URI is set for Docker but we're running locally
	if mongoURI == "mongodb://admin:password123@mongodb:27017" {
		mongoURI = "mongodb://admin:password123@localhost:27017/?authSource=admin"
	} else if mongoURI == "mongodb://admin:password123@localhost:27017" {
		mongoURI = "mongodb://admin:password123@localhost:27017/?authSource=admin"
	}

	fmt.Println("Connecting to MongoDB...")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(ctx)
	
	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	
	fmt.Println("Connected to MongoDB successfully!")

	// Get database
	db := client.Database("cafe_pos")
	ingredientsCollection := db.Collection("ingredients")

	fmt.Println("Starting migration: Add conversion_rate and wastage_percentage to ingredients...")

	// Update all ingredients that don't have conversion_rate or wastage_percentage
	filter := bson.M{
		"$or": []bson.M{
			{"conversion_rate": bson.M{"$exists": false}},
			{"wastage_percentage": bson.M{"$exists": false}},
		},
	}

	update := bson.M{
		"$set": bson.M{
			"conversion_rate":    1.0,
			"wastage_percentage": 0.0,
		},
	}

	result, err := ingredientsCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Fatal("Failed to update ingredients:", err)
	}

	fmt.Printf("Migration completed successfully!\n")
	fmt.Printf("- Matched documents: %d\n", result.MatchedCount)
	fmt.Printf("- Modified documents: %d\n", result.ModifiedCount)

	// Verify the migration
	fmt.Println("\nVerifying migration...")
	
	// Count ingredients with conversion_rate = 1.0
	countWithConversion, err := ingredientsCollection.CountDocuments(ctx, bson.M{"conversion_rate": 1.0})
	if err != nil {
		log.Fatal("Failed to count ingredients with conversion_rate:", err)
	}
	
	// Count ingredients with wastage_percentage = 0.0
	countWithWastage, err := ingredientsCollection.CountDocuments(ctx, bson.M{"wastage_percentage": 0.0})
	if err != nil {
		log.Fatal("Failed to count ingredients with wastage_percentage:", err)
	}
	
	// Count total ingredients
	totalCount, err := ingredientsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to count total ingredients:", err)
	}

	fmt.Printf("- Total ingredients: %d\n", totalCount)
	fmt.Printf("- Ingredients with conversion_rate = 1.0: %d\n", countWithConversion)
	fmt.Printf("- Ingredients with wastage_percentage = 0.0: %d\n", countWithWastage)
	
	if countWithConversion == totalCount && countWithWastage == totalCount {
		fmt.Println("\n✓ Migration verified successfully! All ingredients have default values.")
	} else {
		fmt.Println("\n⚠ Warning: Some ingredients may not have been updated correctly.")
	}
}
