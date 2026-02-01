package main

import (
	"context"
	"log"
	"os"
	"time"
	"cafe-pos/backend/application/services"
	"cafe-pos/backend/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MongoDB connection
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://admin:password123@localhost:27017"
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database("cafe_pos")
	ctx := context.Background()

	// Repositories
	userRepo := mongodb.NewUserRepository(db)

	// Services
	authService := services.NewAuthService(userRepo, services.NewJWTService("secret"))

	// Update admin password
	updateAdminPassword(ctx, authService, userRepo)

	log.Println("✅ Admin password updated successfully!")
}

func updateAdminPassword(ctx context.Context, authService *services.AuthService, userRepo *mongodb.UserRepository) {
	// Find admin user
	admin, err := userRepo.FindByUsername(ctx, "admin")
	if err != nil {
		log.Fatalf("❌ Admin user not found: %v", err)
	}

	// Hash new password
	hashedPassword, err := authService.HashPassword("admin123")
	if err != nil {
		log.Fatalf("❌ Failed to hash password: %v", err)
	}

	// Update password
	admin.Password = hashedPassword
	admin.UpdatedAt = time.Now()

	err = userRepo.Update(ctx, admin.ID, admin)
	if err != nil {
		log.Fatalf("❌ Failed to update admin password: %v", err)
	}

	log.Printf("✅ Updated admin password to: admin123")
}
