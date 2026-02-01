package main

import (
	"context"
	"log"
	"time"
	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/user"
	"cafe-pos/backend/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MongoDB connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
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

	// Seed only admin user
	seedAdminUser(ctx, authService, userRepo)

	log.Println("✅ Seed admin user completed!")
}

func seedAdminUser(ctx context.Context, authService *services.AuthService, userRepo *mongodb.UserRepository) {
	// Check if admin exists
	if _, err := userRepo.FindByUsername(ctx, "admin"); err == nil {
		log.Println("⚠️  Admin user already exists, skipping...")
		return
	}

	// Create admin user
	hashedPassword, _ := authService.HashPassword("admin123")
	newUser := &user.User{
		Username:  "admin",
		Password:  hashedPassword,
		Role:      user.RoleManager,
		Name:      "Administrator",
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Create(ctx, newUser)
	log.Printf("✅ Created user: admin (role: %s)", user.RoleManager)
}
