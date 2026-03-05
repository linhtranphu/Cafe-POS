package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/infrastructure/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("Testing Chromedp HTML capture...")
	
	// Connect to MongoDB
	mongoURI := "mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin"
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())
	
	// Get shop settings
	settingsRepo := mongodb.NewShopSettingsRepository(client.Database("cafe_pos"))
	shopSettings, err := settingsRepo.GetSettings(context.Background())
	if err != nil {
		log.Fatal("Failed to get shop settings:", err)
	}
	
	// Get an order
	orderRepo := mongodb.NewOrderRepository(client.Database("cafe_pos"))
	var orders []interface{}
	cursor, err := orderRepo.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal("Failed to query orders:", err)
	}
	defer cursor.Close(context.Background())
	
	if err = cursor.All(context.Background(), &orders); err != nil || len(orders) == 0 {
		log.Fatal("Failed to get orders or no orders found")
	}
	
	// Get first order
	var testOrder interface{}
	err = orderRepo.Collection.FindOne(context.Background(), bson.M{}).Decode(&testOrder)
	if err != nil {
		log.Fatal("Failed to get order:", err)
	}
	
	// For now, just test if we can create the renderer
	renderer, err := services.NewChromedpBillRendererOptimized()
	if err != nil {
		log.Fatal("Failed to create renderer:", err)
	}
	defer renderer.Close()
	
	fmt.Println("✓ Chromedp renderer created successfully")
	fmt.Println("✓ Chrome context initialized")
	fmt.Println("\nTo test actual capture, use the preview endpoint:")
	fmt.Println("  POST /api/manager/html-templates/preview")
	fmt.Println("  with order_id in body")
	fmt.Println("\nOr use the SavePreviewImage method in code")
}

