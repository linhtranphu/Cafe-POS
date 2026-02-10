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

type Ingredient struct {
	Name              string  `bson:"name"`
	ConversionRate    float64 `bson:"conversion_rate"`
	WastagePercentage float64 `bson:"wastage_percentage"`
}

func main() {
	// Get MongoDB URI from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	
	// Replace 'mongodb' hostname with 'localhost' for local execution
	if mongoURI == "mongodb://admin:password123@mongodb:27017" {
		mongoURI = "mongodb://admin:password123@localhost:27017/?authSource=admin"
	} else if mongoURI == "mongodb://admin:password123@localhost:27017" {
		mongoURI = "mongodb://admin:password123@localhost:27017/?authSource=admin"
	}

	fmt.Println("Verifying ingredient conversion_rate and wastage_percentage fields...")

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

	// Get database
	db := client.Database("cafe_pos")
	ingredientsCollection := db.Collection("ingredients")

	// Find all ingredients
	cursor, err := ingredientsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to find ingredients:", err)
	}
	defer cursor.Close(ctx)

	var ingredients []Ingredient
	if err := cursor.All(ctx, &ingredients); err != nil {
		log.Fatal("Failed to decode ingredients:", err)
	}

	fmt.Printf("\nFound %d ingredients:\n", len(ingredients))
	fmt.Println("----------------------------------------")
	
	allHaveFields := true
	for i, ing := range ingredients {
		fmt.Printf("%d. %s\n", i+1, ing.Name)
		fmt.Printf("   - conversion_rate: %.1f\n", ing.ConversionRate)
		fmt.Printf("   - wastage_percentage: %.1f\n", ing.WastagePercentage)
		
		if ing.ConversionRate == 0 && ing.WastagePercentage == 0 {
			// Check if fields exist (0 could be default or missing)
			var raw bson.M
			ingredientsCollection.FindOne(ctx, bson.M{"name": ing.Name}).Decode(&raw)
			if _, hasConversion := raw["conversion_rate"]; !hasConversion {
				fmt.Println("   ⚠️  WARNING: conversion_rate field is missing!")
				allHaveFields = false
			}
			if _, hasWastage := raw["wastage_percentage"]; !hasWastage {
				fmt.Println("   ⚠️  WARNING: wastage_percentage field is missing!")
				allHaveFields = false
			}
		}
	}
	
	fmt.Println("----------------------------------------")
	if allHaveFields {
		fmt.Println("\n✓ All ingredients have conversion_rate and wastage_percentage fields!")
		fmt.Println("✓ Task 1.3 implementation verified successfully!")
	} else {
		fmt.Println("\n✗ Some ingredients are missing the required fields!")
	}
}
