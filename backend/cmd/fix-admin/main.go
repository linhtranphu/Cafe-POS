package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("🔧 Fixing admin user role...")
	fmt.Println()

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(ctx)

	// Get database and collection
	db := client.Database("cafe_pos")
	usersCollection := db.Collection("users")

	// Update admin user role to manager
	filter := bson.M{"username": "admin"}
	update := bson.M{"$set": bson.M{"role": "manager"}}

	result, err := usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal("Failed to update admin role:", err)
	}

	fmt.Printf("Matched: %d\n", result.MatchedCount)
	fmt.Printf("Modified: %d\n", result.ModifiedCount)
	fmt.Println()

	// Verify the change
	var admin struct {
		Username string `bson:"username"`
		Role     string `bson:"role"`
		Name     string `bson:"name"`
	}

	err = usersCollection.FindOne(ctx, filter).Decode(&admin)
	if err != nil {
		log.Fatal("Failed to find admin user:", err)
	}

	fmt.Println("Admin user:")
	fmt.Printf("  Username: %s\n", admin.Username)
	fmt.Printf("  Role: %s\n", admin.Role)
	fmt.Printf("  Name: %s\n", admin.Name)
	fmt.Println()

	if admin.Role == "manager" {
		fmt.Println("✅ Admin role successfully updated to manager")
	} else {
		fmt.Println("❌ Failed to update admin role")
	}

	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("1. Logout from the app")
	fmt.Println("2. Clear browser cache (Ctrl+Shift+Delete)")
	fmt.Println("3. Run: localStorage.clear() in console")
	fmt.Println("4. Login again with admin/admin123")
	fmt.Println("5. You should see all manager menus including:")
	fmt.Println("   - 👥 Quản lý User")
	fmt.Println("   - 🍽️ Menu")
	fmt.Println("   - 🥬 Nguyên liệu")
	fmt.Println("   - 🏢 Cơ sở vật chất")
	fmt.Println("   - 💰 Chi phí")
	fmt.Println()
}
