package main

import (
	"context"
	"log"
	"os/exec"
	"time"
	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain"
	"cafe-pos/backend/domain/user"
	"cafe-pos/backend/infrastructure/mongodb"
	"cafe-pos/backend/interfaces/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Kill existing processes on port 3000
	log.Println("🔄 Stopping existing processes on port 3000...")
	exec.Command("bash", "-c", "lsof -ti:3000 | xargs -r kill -9").Run()
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
	orderRepo := mongodb.NewOrderRepository(db)
	shiftRepo := mongodb.NewShiftRepository(db)
	// Cashier repositories
	cashierShiftRepo := mongodb.NewCashierShiftRepository(db)
	cashReconciliationRepo := mongodb.NewCashReconciliationRepository(db)
	paymentDiscrepancyRepo := mongodb.NewPaymentDiscrepancyRepository(db)
	paymentAuditRepo := mongodb.NewPaymentAuditRepository(db)

	// State Machine Manager
	smManager := domain.NewStateMachineManager()

	// Services
	jwtService := services.NewJWTService("your-secret-key")
	authService := services.NewAuthService(userRepo, jwtService)
	userManagementService := services.NewUserManagementService(userRepo, authService)
	orderService := services.NewOrderService(orderRepo, shiftRepo, smManager)
	shiftService := services.NewShiftService(shiftRepo, orderRepo, smManager)
	// Cashier services
	cashierShiftService := services.NewCashierShiftService(cashierShiftRepo, shiftRepo, smManager)
	cashReconciliationService := services.NewCashReconciliationService(cashReconciliationRepo, shiftRepo, orderRepo)
	paymentOversightService := services.NewPaymentOversightService(orderRepo, paymentDiscrepancyRepo, paymentAuditRepo)
	cashierReportService := services.NewCashierReportService(orderRepo, cashReconciliationRepo, shiftRepo, paymentAuditRepo)

	// Handlers
	authHandler := http.NewAuthHandler(authService)
	userManagementHandler := http.NewUserManagementHandler(userManagementService)
	orderHandler := http.NewOrderHandler(orderService, smManager)
	shiftHandler := http.NewShiftHandler(shiftService, smManager)
	// Cashier handlers
	cashierShiftHandler := http.NewCashierShiftHandler(cashierShiftService)
	cashierShiftClosureHandler := http.NewCashierShiftClosureHandler(cashierShiftService, smManager)
	cashierHandler := http.NewCashierHandler(cashReconciliationService, paymentOversightService, cashierReportService)
	// State machine handler
	stateMachineHandler := http.NewStateMachineHandler(smManager)
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
	expenseRepo := mongodb.NewExpenseRepository(db)
	expenseService := services.NewExpenseService(expenseRepo)
	expenseHandler := http.NewExpenseHandler(expenseService)

	// Auto Expense Service - wire up with other services
	autoExpenseService := services.NewAutoExpenseService(expenseService)
	ingredientService.SetAutoExpenseService(autoExpenseService)
	facilityService.SetAutoExpenseService(autoExpenseService)

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
		
		// State machine information endpoints (public for documentation)
		api.GET("/state-machines", stateMachineHandler.GetAllStateMachines)
		api.GET("/state-machines/cashier-shift", stateMachineHandler.GetCashierShiftStates)
		api.GET("/state-machines/waiter-shift", stateMachineHandler.GetWaiterShiftStates)
		api.GET("/state-machines/order", stateMachineHandler.GetOrderStates)
		
		// Protected routes
		protected := api.Group("/")
		protected.Use(http.AuthMiddleware(jwtService))
		{
			// Common routes for all authenticated users
			protected.GET("/profile", userManagementHandler.GetCurrentUser)
			protected.POST("/change-password", userManagementHandler.ChangePassword)
			
			// Shift management - available for waiter and barista only
			// Note: Cashier shifts use separate endpoints under /cashier-shifts
			shifts := protected.Group("/shifts")
			{
				shifts.POST("/start", shiftHandler.StartShift)
				shifts.POST("/:id/end", shiftHandler.EndShift)
				shifts.POST("/:id/close", shiftHandler.CloseShift)
				shifts.GET("/current", shiftHandler.GetCurrentShift)
				shifts.GET("/my", shiftHandler.GetMyShifts)
				shifts.GET("/:id", shiftHandler.GetShift)
			}
			
			// Cashier shift management - separate from waiter/barista shifts
			cashierShifts := protected.Group("/cashier-shifts")
			cashierShifts.Use(http.RequireRole(user.RoleCashier, user.RoleManager))
			{
				cashierShifts.POST("", cashierShiftHandler.StartCashierShift)
				cashierShifts.GET("/current", cashierShiftHandler.GetCurrentCashierShift)
				cashierShifts.GET("/my-shifts", cashierShiftHandler.GetMyCashierShifts)
				cashierShifts.GET("/:id", cashierShiftHandler.GetCashierShift)
				
				// Shift closure workflow
				cashierShifts.GET("/check-waiter-shifts", cashierShiftClosureHandler.CheckWaiterShifts)
				cashierShifts.POST("/:id/initiate-closure", cashierShiftClosureHandler.InitiateClosure)
				cashierShifts.POST("/:id/record-actual-cash", cashierShiftClosureHandler.RecordActualCash)
				cashierShifts.POST("/:id/document-variance", cashierShiftClosureHandler.DocumentVariance)
				cashierShifts.POST("/:id/confirm-responsibility", cashierShiftClosureHandler.ConfirmResponsibility)
				cashierShifts.POST("/:id/close", cashierShiftClosureHandler.CloseShift)
			}
			
			// Cashier shift management - manager only
			cashierShiftsManager := protected.Group("/cashier-shifts")
			cashierShiftsManager.Use(http.RequireRole(user.RoleManager))
			{
				cashierShiftsManager.GET("", cashierShiftHandler.GetAllCashierShifts)
			}
			
			// Waiter routes
			waiter := protected.Group("/waiter")
			waiter.Use(http.RequireRole(user.RoleWaiter, user.RoleCashier, user.RoleManager))
			{
				// Order management
				waiter.POST("/orders", orderHandler.CreateOrder)
				waiter.POST("/orders/:id/payment", orderHandler.CollectPayment)
				waiter.PUT("/orders/:id/edit", orderHandler.EditOrder)
				waiter.POST("/orders/:id/send", orderHandler.SendToBar)
				waiter.POST("/orders/:id/serve", orderHandler.ServeOrder)
				waiter.GET("/orders", orderHandler.GetMyOrders)
				waiter.GET("/orders/:id", orderHandler.GetOrder)
				
				// Menu (read-only)
				waiter.GET("/menu", menuHandler.GetAllMenuItems)
				
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

			// Barista routes
			barista := protected.Group("/barista")
			barista.Use(http.RequireRole(user.RoleBarista, user.RoleManager))
			{
				// View queued orders
				barista.GET("/orders/queue", orderHandler.GetQueuedOrders)
				// View my orders (in progress + ready)
				barista.GET("/orders/my", orderHandler.GetMyBaristaOrders)
				// Accept order from queue
				barista.POST("/orders/:id/accept", orderHandler.AcceptOrder)
				// Mark order as ready
				barista.POST("/orders/:id/ready", orderHandler.FinishPreparing)
				// View order details
				barista.GET("/orders/:id", orderHandler.GetOrder)
			}

			// Cashier routes
			cashier := protected.Group("/cashier")
			cashier.Use(http.RequireRole(user.RoleCashier, user.RoleManager))
			{
				// Order management
				cashier.GET("/orders", orderHandler.GetAllOrders)
				cashier.GET("/orders/:id", orderHandler.GetOrder)
				cashier.POST("/orders/:id/cancel", orderHandler.CancelOrder)
				cashier.POST("/orders/:id/refund", orderHandler.RefundPartial)
				
				// Shift management
				cashier.POST("/shifts/:id/close", shiftHandler.CloseShift)
				cashier.GET("/shifts", shiftHandler.GetAllShifts)
				cashier.GET("/shifts/:id", shiftHandler.GetShift)
				
				// Cashier-specific routes
				cashier.GET("/shifts/:id/status", cashierHandler.GetShiftStatus)
				cashier.GET("/shifts/:id/payments", cashierHandler.GetPaymentsByShift)
				cashier.POST("/discrepancies", cashierHandler.ReportDiscrepancy)
				cashier.GET("/discrepancies/pending", cashierHandler.GetPendingDiscrepancies)
				cashier.POST("/discrepancies/:id/resolve", cashierHandler.ResolveDiscrepancy)
				cashier.POST("/reconcile/cash", cashierHandler.ReconcileCash)
				cashier.POST("/orders/:id/override", cashierHandler.OverridePayment)
				cashier.POST("/orders/:id/lock", cashierHandler.LockOrder)
				cashier.GET("/reports/shift/:id", cashierHandler.GenerateShiftReport)
				cashier.GET("/reports/daily", cashierHandler.GetDailyReport)
				cashier.POST("/handover", cashierHandler.HandoverShift)
				cashier.GET("/orders/:id/audits", cashierHandler.GetOrderAudits)
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
				
				// Ingredient category routes
				manager.POST("/ingredient-categories", ingredientHandler.CreateCategory)
				manager.GET("/ingredient-categories", ingredientHandler.GetCategories)
				manager.DELETE("/ingredient-categories/:id", ingredientHandler.DeleteCategory)
				
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
				
				// Facility type and area routes
				manager.POST("/facility-types", facilityHandler.CreateFacilityType)
				manager.GET("/facility-types", facilityHandler.GetFacilityTypes)
				manager.DELETE("/facility-types/:id", facilityHandler.DeleteFacilityType)
				manager.POST("/facility-areas", facilityHandler.CreateFacilityArea)
				manager.GET("/facility-areas", facilityHandler.GetFacilityAreas)
				manager.DELETE("/facility-areas/:id", facilityHandler.DeleteFacilityArea)
				
				// Expense management routes
				manager.POST("/expenses", expenseHandler.CreateExpense)
				manager.GET("/expenses", expenseHandler.GetExpenses)
				manager.PUT("/expenses/:id", expenseHandler.UpdateExpense)
				manager.DELETE("/expenses/:id", expenseHandler.DeleteExpense)
				manager.POST("/expense-categories", expenseHandler.CreateCategory)
				manager.GET("/expense-categories", expenseHandler.GetCategories)
				manager.DELETE("/expense-categories/:id", expenseHandler.DeleteCategory)
				manager.POST("/recurring-expenses", expenseHandler.CreateRecurring)
				manager.GET("/recurring-expenses", expenseHandler.GetRecurring)
				manager.DELETE("/recurring-expenses/:id", expenseHandler.DeleteRecurring)
				manager.POST("/prepaid-expenses", expenseHandler.CreatePrepaid)
				manager.GET("/prepaid-expenses", expenseHandler.GetPrepaid)
				manager.DELETE("/prepaid-expenses/:id", expenseHandler.DeletePrepaid)
				
				// User management routes
				manager.POST("/users", userManagementHandler.CreateUser)
				manager.GET("/users", userManagementHandler.GetAllUsers)
				manager.GET("/users/active", userManagementHandler.GetActiveUsers)
				manager.GET("/users/by-role", userManagementHandler.GetUsersByRole)
				manager.GET("/users/:id", userManagementHandler.GetUser)
				manager.PUT("/users/:id", userManagementHandler.UpdateUser)
				manager.POST("/users/:id/reset-password", userManagementHandler.ResetPassword)
				manager.POST("/users/:id/toggle-status", userManagementHandler.ToggleUserStatus)
				manager.DELETE("/users/:id", userManagementHandler.DeleteUser)
				
				// Order management routes (full access)
				manager.GET("/orders", orderHandler.GetAllOrders)
				manager.GET("/orders/:id", orderHandler.GetOrder)
				manager.POST("/orders", orderHandler.CreateOrder)
				manager.POST("/orders/:id/cancel", orderHandler.CancelOrder)
				manager.POST("/orders/:id/refund", orderHandler.RefundPartial)
				manager.PUT("/orders/:id/edit", orderHandler.EditOrder)
				
				// Shift management routes
				manager.GET("/shifts", shiftHandler.GetAllShifts)
				manager.GET("/shifts/:id", shiftHandler.GetShift)
			}
		}
	}

	// Create default admin user
	// createDefaultUsers(authService, userRepo)

	log.Println("Server starting on :3000")
	r.Run(":3000")
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
		{"barista1", "barista123", user.RoleBarista, "Barista 1"},
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