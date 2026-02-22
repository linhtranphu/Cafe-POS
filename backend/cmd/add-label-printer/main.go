package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cafe-pos/backend/domain/printing"
	"cafe-pos/backend/infrastructure/mongodb"

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
	printerRepo := mongodb.NewPrinterConfigRepository(db)

	ctx = context.Background()

	// Get existing BILL printer to copy settings
	billPrinter, err := printerRepo.FindDefault(ctx, printing.PrinterTypeBill)
	if err != nil || billPrinter == nil {
		log.Fatal("No default BILL printer found. Please configure BILL printer first.")
	}

	fmt.Printf("Found BILL printer: %s (IP: %s:%d)\n", billPrinter.Name, billPrinter.IPAddress, billPrinter.Port)
	fmt.Println("Creating LABEL printer with same settings...")

	// Create LABEL printer with same IP/Port
	labelPrinter := &printing.PrinterConfig{
		Name:           billPrinter.Name + " (Label)",
		Type:           printing.PrinterTypeLabel,
		ConnectionType: billPrinter.ConnectionType,
		IPAddress:      billPrinter.IPAddress,
		Port:           billPrinter.Port,
		PaperWidth:     58, // Labels usually use smaller paper
		IsDefault:      true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := printerRepo.Create(ctx, labelPrinter); err != nil {
		log.Fatalf("Failed to create LABEL printer: %v", err)
	}

	fmt.Printf("✅ Created LABEL printer: %s\n", labelPrinter.Name)
	fmt.Println("\n✅ Auto-print should now work!")
	fmt.Println("\nNote: Both BILL and LABEL will print to the same physical printer.")
	fmt.Println("If you have separate printers, update the LABEL printer IP in Print Management.")
}
