package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/user"
	"cafe-pos/backend/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Use correct MongoDB URI with auth
	uri := "mongodb://admin:108trannhatduat@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database("cafe_pos")
	ctx := context.Background()

	// Repositories
	userRepo := mongodb.NewUserRepository(db)

	// Services
	authService := services.NewAuthService(userRepo, services.NewJWTService("your-jwt-secret-key-min-32-chars-long"))

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
	err = userRepo.Create(ctx, newUser)
	if err != nil {
		log.Fatal("Failed to create user:", err)
	}
	
	fmt.Println("✅ Created admin user successfully!")
	fmt.Println("Username: admin")
	fmt.Println("Password: admin123")
}
