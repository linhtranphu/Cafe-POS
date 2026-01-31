package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	maintenanceCollection := db.Collection("maintenance_records")
	facilitiesCollection := db.Collection("facilities")

	// The facility ID with incorrect maintenance records
	facilityID, _ := primitive.ObjectIDFromHex("697629c9d0ac2facdbb23baa")

	// Get facility info
	var facility bson.M
	err = facilitiesCollection.FindOne(context.Background(), bson.M{"_id": facilityID}).Decode(&facility)
	if err != nil {
		log.Fatal("Facility not found:", err)
	}

	fmt.Println("=== Cleaning Incorrect Maintenance Records ===")
	fmt.Printf("Facility: %s - %s\n", facility["name"], facility["type"])

	// Count maintenance records
	count, err := maintenanceCollection.CountDocuments(context.Background(), bson.M{"facility_id": facilityID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nFound %d maintenance records\n", count)

	// Get sample records
	cursor, err := maintenanceCollection.Find(context.Background(), bson.M{"facility_id": facilityID}, options.Find().SetLimit(3))
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	fmt.Println("\nSample maintenance records:")
	for cursor.Next(context.Background()) {
		var record bson.M
		cursor.Decode(&record)
		fmt.Printf("  - %s\n", record["description"])
	}

	fmt.Println("\n⚠️  These maintenance records are incorrect!")
	fmt.Println("They describe 'Máy pha cà phê' but are attached to 'Bàn khách 2 chỗ'")

	// Delete incorrect records
	fmt.Println("\nDeleting incorrect maintenance records...")
	result, err := maintenanceCollection.DeleteMany(context.Background(), bson.M{"facility_id": facilityID})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ Deleted %d incorrect maintenance records\n", result.DeletedCount)

	// Verify
	remaining, _ := maintenanceCollection.CountDocuments(context.Background(), bson.M{"facility_id": facilityID})
	fmt.Printf("Remaining maintenance records for this facility: %d\n", remaining)

	fmt.Println("\n=== Cleanup Complete ===")
	fmt.Println("You can now delete the facility 'Bàn khách 2 chỗ'")
}
