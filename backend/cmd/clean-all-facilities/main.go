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
	// Connect to MongoDB
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("cafe_pos")

	fmt.Println("=== Cleaning All Facilities and Maintenance Data ===\n")

	// Collections to clean
	collections := []struct {
		name        string
		description string
	}{
		{"maintenance_records", "Maintenance Records"},
		{"facility_history", "Facility History"},
		{"issue_reports", "Issue Reports"},
		{"scheduled_maintenance", "Scheduled Maintenance"},
		{"facilities", "Facilities"},
	}

	// Count and delete from each collection
	for _, col := range collections {
		collection := db.Collection(col.name)

		// Count documents
		count, err := collection.CountDocuments(context.Background(), bson.M{})
		if err != nil {
			log.Printf("Error counting %s: %v\n", col.description, err)
			continue
		}

		if count == 0 {
			fmt.Printf("✓ %s: Already empty (0 documents)\n", col.description)
			continue
		}

		// Delete all documents
		result, err := collection.DeleteMany(context.Background(), bson.M{})
		if err != nil {
			log.Printf("Error deleting %s: %v\n", col.description, err)
			continue
		}

		fmt.Printf("✅ %s: Deleted %d documents\n", col.description, result.DeletedCount)
	}

	fmt.Println("\n=== Cleanup Complete ===")
	fmt.Println("All facilities and related data have been removed.")
	fmt.Println("You can now start fresh with new facility data.")
}
