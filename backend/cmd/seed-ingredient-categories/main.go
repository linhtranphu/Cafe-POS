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

type IngredientCategory struct {
	Name      string    `bson:"name"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

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
	collection := db.Collection("ingredient_categories")

	fmt.Println("🌱 Seeding ingredient categories...")
	fmt.Println()

	defaultCategories := []string{
		"Cà phê",
		"Trà",
		"Sữa",
		"Đường",
		"Syrup",
		"Topping",
		"Khác",
	}

	now := time.Now()
	for _, catName := range defaultCategories {
		// Check if category already exists
		count, err := collection.CountDocuments(ctx, bson.M{"name": catName})
		if err != nil {
			log.Printf("Error checking category %s: %v\n", catName, err)
			continue
		}

		if count > 0 {
			fmt.Printf("⏭️  Category '%s' already exists, skipping\n", catName)
			continue
		}

		cat := IngredientCategory{
			Name:      catName,
			CreatedAt: now,
			UpdatedAt: now,
		}

		_, err = collection.InsertOne(ctx, cat)
		if err != nil {
			log.Printf("Error inserting category %s: %v\n", catName, err)
			continue
		}

		fmt.Printf("✅ Created category: %s\n", catName)
	}

	fmt.Println()
	fmt.Println("🎉 Ingredient categories seeded successfully!")
}
