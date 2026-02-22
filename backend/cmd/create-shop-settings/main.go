package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cafe-pos/backend/domain/settings"
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
	settingsRepo := mongodb.NewShopSettingsRepository(db)

	ctx = context.Background()

	// Check if settings already exist
	existing, err := settingsRepo.GetSettings(ctx)
	if err == nil && existing != nil {
		fmt.Printf("✅ Shop settings already exist\n")
		fmt.Printf("   Shop Name: %s\n", existing.ShopName)
		fmt.Printf("   Auto Print: %v\n", existing.AutoPrintEnabled)
		return
	}

	// Create default shop settings
	defaultSettings := settings.NewShopSettings("Cafe POS")
	
	// Set additional fields
	defaultSettings.ShopAddress = "123 Main Street"
	defaultSettings.ShopPhone = "0123-456-789"
	defaultSettings.CustomMessage = "Cảm ơn quý khách! Hẹn gặp lại!"
	defaultSettings.SetAutoPrintEnabled(true)

	if err := settingsRepo.CreateSettings(ctx, defaultSettings); err != nil {
		log.Fatalf("Failed to create shop settings: %v", err)
	}

	fmt.Println("✅ Created default shop settings:")
	fmt.Printf("   Shop Name: %s\n", defaultSettings.ShopName)
	fmt.Printf("   Address: %s\n", defaultSettings.ShopAddress)
	fmt.Printf("   Phone: %s\n", defaultSettings.ShopPhone)
	fmt.Printf("   Auto Print: %v\n", defaultSettings.AutoPrintEnabled)
	fmt.Println()
	fmt.Println("You can update these settings from the frontend Settings page.")
}
