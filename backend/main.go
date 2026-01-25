package main

import (
	"context"
	"log"
	"os/exec"
	"time"
	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/user"
	"cafe-pos/backend/infrastructure/mongodb"
	"cafe-pos/backend/interfaces/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Kill existing processes on port 8080
	log.Println("🔄 Stopping existing processes on port 8080...")
	exec.Command("bash", "-c", "lsof -ti:8080 | xargs -r kill -9").Run()
	exec.Command("pkill", "-f", "cafe-pos-server").Run()
	exec.Command("pkill", "-f", "go run main.go").Run()
	time.Sleep(2 * time.Second)

	// MongoDB connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database("cafe_pos")

	// Repositories
	userRepo := mongodb.NewUserRepository(db)

	// Services
	jwtService := services.NewJWTService("your-secret-key")
	authService := services.NewAuthService(userRepo, jwtService)

	// Handlers
	authHandler := http.NewAuthHandler(authService)
	menuRepo := mongodb.NewMenuRepository(db)
	menuService := services.NewMenuService(menuRepo)
	menuHandler := http.NewMenuHandler(menuService)
	ingredientRepo := mongodb.NewIngredientRepository(db)
	stockHistoryRepo := mongodb.NewStockHistoryRepository(db)
	ingredientService := services.NewIngredientService(ingredientRepo, stockHistoryRepo)
	ingredientHandler := http.NewIngredientHandler(ingredientService)
	facilityRepo := mongodb.NewFacilityRepository(db)
	facilityService := services.NewFacilityService(facilityRepo)
	facilityHandler := http.NewFacilityHandler(facilityService)

	// Router
	r := gin.Default()
	
	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Routes
	api := r.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		
		// Protected routes
		protected := api.Group("/")
		protected.Use(http.AuthMiddleware(jwtService))
		{
			// Waiter routes
			waiter := protected.Group("/waiter")
			waiter.Use(http.RequireRole(user.RoleWaiter, user.RoleCashier, user.RoleManager))
			{
				waiter.GET("/profile", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "waiter access"})
				})
				// Staff can view ingredients
				waiter.GET("/ingredients", ingredientHandler.GetAllIngredients)
				// Staff can view facilities and report issues
				waiter.GET("/facilities", facilityHandler.GetAllFacilities)
				waiter.GET("/facilities/search", facilityHandler.SearchFacilities)
				waiter.POST("/issues", facilityHandler.CreateIssueReport)
			}

			// Manager routes
			manager := protected.Group("/manager")
			manager.Use(http.RequireRole(user.RoleManager))
			{
				manager.GET("/reports", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "manager access"})
				})
				
				// Menu management routes
				manager.POST("/menu", menuHandler.CreateMenuItem)
				manager.GET("/menu", menuHandler.GetAllMenuItems)
				manager.GET("/menu/:id", menuHandler.GetMenuItem)
				manager.PUT("/menu/:id", menuHandler.UpdateMenuItem)
				manager.DELETE("/menu/:id", menuHandler.DeleteMenuItem)
				
				// Ingredient management routes
				manager.POST("/ingredients", ingredientHandler.CreateIngredient)
				manager.GET("/ingredients", ingredientHandler.GetAllIngredients)
				manager.GET("/ingredients/low-stock", ingredientHandler.GetLowStock)
				manager.GET("/ingredients/:id", ingredientHandler.GetIngredient)
				manager.GET("/ingredients/:id/history", ingredientHandler.GetStockHistory)
				manager.PUT("/ingredients/:id", ingredientHandler.UpdateIngredient)
				manager.DELETE("/ingredients/:id", ingredientHandler.DeleteIngredient)
				manager.POST("/ingredients/:id/adjust", ingredientHandler.AdjustStock)
				
				// Facility management routes
				manager.GET("/facilities", facilityHandler.GetAllFacilities)
				manager.GET("/facilities/search", facilityHandler.SearchFacilities)
				manager.GET("/facilities/:id", facilityHandler.GetFacility)
				manager.POST("/facilities", facilityHandler.CreateFacility)
				manager.PUT("/facilities/:id", facilityHandler.UpdateFacility)
				manager.DELETE("/facilities/:id", facilityHandler.DeleteFacility)
				manager.GET("/facilities/:id/history", facilityHandler.GetFacilityHistory)
				manager.GET("/facilities/:id/next-maintenance", facilityHandler.GetNextMaintenanceDate)
				manager.GET("/facilities/:id/maintenance-stats", facilityHandler.GetMaintenanceStats)
				manager.GET("/facilities/:id/status-history", facilityHandler.GetStatusHistory)
				manager.GET("/facilities/history", facilityHandler.GetHistoryWithFilter)
				manager.GET("/facilities/:id/maintenance", facilityHandler.GetMaintenanceHistory)
				manager.POST("/maintenance", facilityHandler.CreateMaintenanceRecord)
				manager.GET("/maintenance/scheduled", facilityHandler.GetScheduledMaintenance)
				manager.GET("/maintenance/due", facilityHandler.GetMaintenanceDue)
				manager.GET("/issues", facilityHandler.GetIssueReports)
				manager.POST("/issues", facilityHandler.CreateIssueReport)
			}
		}
	}

	// Create default admin user
	// createDefaultUsers(authService, userRepo)

	log.Println("Server starting on :8080")
	r.Run(":8080")
}

func createDefaultUsers(authService *services.AuthService, userRepo *mongodb.UserRepository) {
	ctx := context.Background()
	
	// Check if admin exists
	if _, err := userRepo.FindByUsername(ctx, "admin"); err == nil {
		return
	}

	// Create default users
	users := []struct {
		username string
		password string
		role     user.Role
		name     string
	}{
		{"admin", "admin123", user.RoleManager, "Administrator"},
		{"waiter1", "waiter123", user.RoleWaiter, "Waiter 1"},
		{"cashier1", "cashier123", user.RoleCashier, "Cashier 1"},
	}

	for _, u := range users {
		hashedPassword, _ := authService.HashPassword(u.password)
		newUser := &user.User{
			Username:  u.username,
			Password:  hashedPassword,
			Role:      u.role,
			Name:      u.name,
			Active:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		userRepo.Create(ctx, newUser)
		log.Printf("Created user: %s (role: %s)", u.username, u.role)
	}
}