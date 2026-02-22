package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cafe-pos/backend/infrastructure/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Get MongoDB URI from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin"
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Get database
	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "cafe_pos"
	}
	db := client.Database(dbName)

	// Create repository
	printJobRepo := mongodb.NewPrintJobRepository(db)

	ctx = context.Background()

	fmt.Println("=== Recent Print Jobs (Last 10) ===")
	fmt.Println()

	// Get recent print jobs
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(10)
	cursor, err := db.Collection("print_jobs").Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Fatalf("Failed to fetch print jobs: %v", err)
	}
	defer cursor.Close(ctx)

	count := 0
	for cursor.Next(ctx) {
		var job map[string]interface{}
		if err := cursor.Decode(&job); err != nil {
			continue
		}
		count++

		fmt.Printf("%d. Job ID: %v\n", count, job["_id"])
		fmt.Printf("   Order: %v\n", job["order_number"])
		fmt.Printf("   Type: %v\n", job["type"])
		fmt.Printf("   Status: %v\n", job["status"])
		fmt.Printf("   Printer: %v\n", job["printer_name"])
		if errorMsg, ok := job["error_message"].(string); ok && errorMsg != "" {
			fmt.Printf("   Error: %v\n", errorMsg)
		}
		fmt.Printf("   Created: %v\n", job["created_at"])
		fmt.Println()
	}

	if count == 0 {
		fmt.Println("❌ No print jobs found!")
		fmt.Println()
		fmt.Println("This means print jobs are NOT being created when payment is collected.")
		fmt.Println()
		fmt.Println("Possible reasons:")
		fmt.Println("1. Backend not restarted after adding LABEL printer")
		fmt.Println("2. Order status not changing to PAID")
		fmt.Println("3. Print service not initialized properly")
		fmt.Println()
		fmt.Println("Check backend logs for errors starting with [PRINT]")
	} else {
		fmt.Printf("Found %d print jobs\n", count)
		
		// Check for failed jobs
		failedCount, _ := db.Collection("print_jobs").CountDocuments(ctx, bson.M{"status": "FAILED"})
		if failedCount > 0 {
			fmt.Printf("\n⚠️  Warning: %d failed print jobs\n", failedCount)
		}
	}

	// Check pending jobs
	pendingJobs, err := printJobRepo.FindPending(ctx, 100)
	if err == nil && len(pendingJobs) > 0 {
		fmt.Printf("\n⚠️  %d jobs still PENDING (not processed by print worker)\n", len(pendingJobs))
		fmt.Println("Check if Print Bridge is running!")
	}
}
