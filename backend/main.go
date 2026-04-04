package main

import (
	"context"
	"log"
	"os"
	"time"
	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain"
	"cafe-pos/backend/domain/user"
	"cafe-pos/backend/infrastructure/mongodb"
	"cafe-pos/backend/infrastructure/printbridge"
	"cafe-pos/backend/infrastructure/websocket"
	"cafe-pos/backend/interfaces/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MongoDB connection
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	log.Printf("Connecting to MongoDB: %s", mongoURI)
	
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("❌ MongoDB connection failed: %v", err)
	}
	defer client.Disconnect(context.TODO())

	// Verify MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ MongoDB ping failed: %v", err)
	}
	log.Println("✅ MongoDB connected successfully")

	db := client.Database("cafe_pos")

	// Repositories
	userRepo := mongodb.NewUserRepository(db)
	orderRepo := mongodb.NewOrderRepository(db)
	shiftRepo := mongodb.NewShiftRepository(db)
	menuRepo := mongodb.NewMenuRepository(db)
	ingredientRepo := mongodb.NewIngredientRepository(db)
	orderItemRepo := mongodb.NewOrderItemRepository(db)
	shopSettingsRepo := mongodb.NewShopSettingsRepository(db)
	operatingExpenseRepo := mongodb.NewOperatingExpenseRepository(db)
	// Cashier repositories
	cashierShiftRepo := mongodb.NewCashierShiftRepository(db)

	cashReconciliationRepo := mongodb.NewCashReconciliationRepository(db)
	paymentDiscrepancyRepo := mongodb.NewPaymentDiscrepancyRepository(db)
	paymentAuditRepo := mongodb.NewPaymentAuditRepository(db)
	// Handover repositories
	cashHandoverRepo := mongodb.NewCashHandoverRepository(db)
	cashDiscrepancyRepo := mongodb.NewCashDiscrepancyRepository(db)
	// Batch repositories
	batchDefinitionRepo := mongodb.NewBatchDefinitionRepository(db)
	batchRecordRepo := mongodb.NewBatchRecordRepository(db)
	batchUsageLogRepo := mongodb.NewBatchUsageLogRepository(db)
	// Stock history repository (needed for batch services)
	stockHistoryRepo := mongodb.NewStockHistoryRepository(db)
	// Print repositories
	printJobRepo := mongodb.NewPrintJobRepository(db)
	printerConfigRepo := mongodb.NewPrinterConfigRepository(db)
	printTemplateRepo := mongodb.NewPrintTemplateRepository(db)
	printNotificationRepo := mongodb.NewPrintNotificationRepository(db)
	// Fund repositories
	journalRepo := mongodb.NewJournalRepository(db)

	// State Machine Manager
	smManager := domain.NewStateMachineManager()

	// WebSocket Hub - Initialize and start (keep for backward compatibility)
	wsHub := websocket.NewHub()
	go wsHub.Run()
	log.Println("✅ WebSocket hub started")
	
	// Socket.IO Server - Initialize with standard library
	socketIOServer, err := websocket.NewSocketIOServer()
	if err != nil {
		log.Fatalf("❌ Failed to create Socket.IO server: %v", err)
	}
	go socketIOServer.Serve()
	defer socketIOServer.Close()
	log.Println("✅ Socket.IO server started")
	
	// WebSocket Broadcaster
	wsBroadcaster := websocket.NewBroadcaster(wsHub)
	wsBroadcaster.SetPrinterRepository(printerConfigRepo)
	wsBroadcaster.SetSocketIOServer(socketIOServer)

	// Services
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key"
		log.Println("⚠️  JWT_SECRET not set, using default (not recommended for production)")
	}
	jwtService := services.NewJWTService(jwtSecret)
	authService := services.NewAuthService(userRepo, jwtService)
	userManagementService := services.NewUserManagementService(userRepo, authService)
	// Batch services - initialize before OrderService
	batchCostCalculator := services.NewBatchCostCalculator(ingredientRepo)
	batchDefinitionService := services.NewBatchDefinitionService(batchDefinitionRepo, ingredientRepo)
	batchRecordService := services.NewBatchRecordService(batchRecordRepo, batchDefinitionRepo, ingredientRepo, stockHistoryRepo, userRepo, batchCostCalculator, client)
	batchUsageService := services.NewBatchUsageService(batchRecordRepo, batchUsageLogRepo)
	batchAlertService := services.NewBatchAlertService(batchDefinitionRepo, batchRecordRepo)
	batchReportService := services.NewBatchReportService(batchRecordRepo, batchUsageLogRepo, batchDefinitionRepo)
	// Print services - now with WebSocket broadcaster
	templateRenderer := services.NewTemplateRenderer()
	
	printService := services.NewPrintService(services.PrintServiceConfig{
		PrintJobRepo:      printJobRepo,
		PrinterConfigRepo: printerConfigRepo,
		TemplateRepo:      printTemplateRepo,
		TemplateRenderer:  templateRenderer,
		OrderRepo:         orderRepo,
		ShopSettingsRepo:  shopSettingsRepo,
		WSBroadcaster:     wsBroadcaster,
	})
	printerManager := services.NewPrinterManager()
	
	// Initialize print bridge client from shop settings
	// This allows auto-print to work via print bridge on EC2
	go func() {
		// Wait a bit for services to initialize
		time.Sleep(2 * time.Second)
		
		ctx := context.Background()
		settings, err := shopSettingsRepo.GetSettings(ctx)
		if err == nil && settings != nil && settings.PrintBridgeURL != "" {
			log.Printf("[PRINT BRIDGE] Configuring print bridge: %s", settings.PrintBridgeURL)
			bridgeClient := printbridge.NewClient(settings.PrintBridgeURL, 30*time.Second)
			
			// Test connection
			if bridgeClient.IsAvailable() {
				printerManager.SetPrintBridgeClient(bridgeClient)
				log.Printf("[PRINT BRIDGE] ✅ Print bridge configured and available")
			} else {
				log.Printf("[PRINT BRIDGE] ⚠️  Print bridge configured but not available (will use direct connection)")
			}
		} else {
			log.Printf("[PRINT BRIDGE] No print bridge URL configured (will use direct connection)")
		}
	}()
	
	printNotificationService := services.NewPrintNotificationService(printNotificationRepo)
	printWorker := services.NewPrintWorker(services.PrintWorkerConfig{
		PrintJobRepo:        printJobRepo,
		PrinterConfigRepo:   printerConfigRepo,
		PrinterManager:      printerManager,
		NotificationService: printNotificationService,
		PollInterval:        10 * time.Second,
	})
	printCleanupJob := services.NewPrintCleanupJob(services.PrintCleanupJobConfig{
		PrintJobRepo:     printJobRepo,
		NotificationRepo: printNotificationRepo,
		Interval:         24 * time.Hour,
		RetentionDays:    7,
	})
	// OrderService now includes batchUsageService
	orderService := services.NewOrderService(orderRepo, shiftRepo, menuRepo, smManager, batchUsageService)
	// Wire up print service for auto-printing
	orderService.SetPrintService(printService)
	orderService.SetSettingsRepository(shopSettingsRepo)
	// Fund service (double-entry journal as single source of truth)
	journalService := services.NewJournalService(journalRepo, client)
	// Wire up journal service into order service for order payment journal entries
	orderService.SetJournalService(journalService)
	orderService.SetMongoClient(client)
	shiftService := services.NewShiftService(shiftRepo, orderRepo, smManager, journalService, cashierShiftRepo, client)
	// Cashier services
	cashierShiftService := services.NewCashierShiftService(cashierShiftRepo, shiftRepo, smManager, journalService, client)
	cashReconciliationService := services.NewCashReconciliationService(cashReconciliationRepo, shiftRepo, orderRepo)
	paymentOversightService := services.NewPaymentOversightService(orderRepo, paymentDiscrepancyRepo, paymentAuditRepo)
	cashierReportService := services.NewCashierReportService(orderRepo, cashReconciliationRepo, shiftRepo, paymentAuditRepo, cashHandoverRepo)
	// Handover service
	cashHandoverService := services.NewCashHandoverService(cashHandoverRepo, cashDiscrepancyRepo, shiftRepo, cashierShiftRepo, orderRepo, client, journalService)

	// Create monitoring service for metrics and alerts
	monitoringService := services.NewMonitoringService()
	
	costCalculatorService := services.NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)
	// Wire up batch repository for batch ingredient cost calculation
	batchRecordCostAdapter := services.NewBatchRecordCostAdapter(batchRecordRepo)
	costCalculatorService.SetBatchRecordRepository(batchRecordCostAdapter)
	// Wire up batch definition repository for cost interpolation when no batch records exist
	costCalculatorService.SetBatchDefinitionRepository(batchDefinitionRepo)
	profitAnalyzerService := services.NewProfitAnalyzerService(menuRepo, orderItemRepo, shopSettingsRepo, operatingExpenseRepo)
	costRecalculationService := services.NewCostRecalculationService(costCalculatorService, menuRepo, 4, 1000) // 4 workers, queue size 1000
	operatingExpenseService := services.NewOperatingExpenseService(operatingExpenseRepo)
	shopSettingsService := services.NewShopSettingsService(shopSettingsRepo)

	// Wire up cost calculator with recalculation service and monitoring
	costCalculatorService.SetCostRecalculationService(costRecalculationService)
	costCalculatorService.SetMonitoringService(monitoringService)
	costRecalculationService.SetMonitoringService(monitoringService)

	// Wire up batch record service with cost recalculation and menu repo
	// This enables auto-recalculation of menu costs when new batches are created
	batchRecordService.SetCostRecalculationService(costRecalculationService)
	batchRecordService.SetMenuRepository(menuRepo)

	// Start cost recalculation worker pool
	// Requirements: 9.1, 9.2, 9.3
	log.Println("Starting cost recalculation worker pool...")
	costRecalculationService.Start()
	log.Println("✅ Cost recalculation worker pool started")

	// Start print worker
	// Requirements: 1.8, 1.9, 2.9, 2.10
	log.Println("Starting print worker...")
	go printWorker.Start(context.Background())
	log.Println("✅ Print worker started")

	// Start print cleanup job
	// Requirements: 4.7
	log.Println("Starting print cleanup job...")
	go printCleanupJob.Start(context.Background())
	log.Println("✅ Print cleanup job started")

	// Setup graceful shutdown for worker pool
	defer func() {
		log.Println("Stopping cost recalculation worker pool...")
		costRecalculationService.Stop()
		log.Println("✅ Cost recalculation worker pool stopped")
		
		log.Println("Stopping print worker...")
		printWorker.Stop()
		log.Println("✅ Print worker stopped")
		
		log.Println("Stopping print cleanup job...")
		printCleanupJob.Stop()
		log.Println("✅ Print cleanup job stopped")
	}()

	// Handlers
	authHandler := http.NewAuthHandler(authService)
	userManagementHandler := http.NewUserManagementHandler(userManagementService)
	orderHandler := http.NewOrderHandler(orderService, smManager, printService)
	shiftHandler := http.NewShiftHandler(shiftService, smManager, costCalculatorService)
	// Cashier handlers
	cashierShiftHandler := http.NewCashierShiftHandler(cashierShiftService)
	cashierShiftClosureHandler := http.NewCashierShiftClosureHandler(cashierShiftService, smManager)
	cashierHandler := http.NewCashierHandler(cashReconciliationService, paymentOversightService, cashierReportService)
	// Handover handler
	cashHandoverHandler := http.NewCashHandoverHandler(cashHandoverService)
	// Fund handler
	fundHandler := http.NewFundHandler(journalService)
	// Batch handlers
	batchDefinitionHandler := http.NewBatchDefinitionHandler(batchDefinitionService)
	batchRecordHandler := http.NewBatchRecordHandler(batchRecordService)
	batchUsageHandler := http.NewBatchUsageHandler(batchUsageService)
	batchAlertHandler := http.NewBatchAlertHandler(batchAlertService)
	batchReportHandler := http.NewBatchReportHandler(batchReportService)
	// Print handlers - now with WebSocket broadcaster
	printJobHandler := http.NewPrintJobHandler(printService, printJobRepo, wsBroadcaster)
	printerConfigHandler := http.NewPrinterConfigHandler(printerConfigRepo, printerManager, shopSettingsRepo)
	printTemplateHandler := http.NewPrintTemplateHandler(printTemplateRepo, templateRenderer, shopSettingsRepo)
	shopSettingsHandler := http.NewShopSettingsHandler(shopSettingsRepo)
	logoUploadHandler := http.NewLogoUploadHandler(shopSettingsRepo, "./uploads/logos")
	
	// Visual print handler (DISABLED - moved to print bridge)
	// var visualPrintHandler *http.VisualPrintHandler
	
	// Print Bridge Client (NEW - Dynamic from settings)
	// Handler will fetch print_bridge_url from shop settings at runtime
	var htmlTemplateHandlerBridge *http.HTMLTemplateHandlerBridge
	templatePath := "./application/services/templates/bill_template_optimized.html"
	htmlTemplateHandlerBridge = http.NewHTMLTemplateHandlerBridge(
		orderRepo,
		shopSettingsRepo,
		templatePath,
	)
	log.Println("✅ HTML template handler (bridge) initialized - will use print_bridge_url from settings")
	
	// Label template handler
	var labelTemplateHandler *http.LabelTemplateHandler
	labelTemplatePath := "./application/services/templates/label_template.tspl"
	labelTemplateHandler = http.NewLabelTemplateHandler(
		orderRepo,
		shopSettingsRepo,
		labelTemplatePath,
	)
	log.Println("✅ Label template handler initialized")
	
	// Chromedp print handler (DISABLED - Chromium removed from Docker image)
	// Kept for reference only. Use print bridge instead.
	/*
	chromedpPrintHandler, err := http.NewChromedpPrintHandler(orderRepo, shopSettingsRepo)
	if err != nil {
		log.Printf("Warning: Failed to create Chromedp print handler: %v", err)
		chromedpPrintHandler = nil
	} else {
		log.Println("✅ Chromedp print handler initialized (DEPRECATED)")
		// Cleanup on exit
		defer func() {
			if chromedpPrintHandler != nil {
				chromedpPrintHandler.Close()
				log.Println("✅ Chromedp print handler closed")
			}
		}()
	}
	*/
	log.Println("ℹ️  Chromedp disabled - using print bridge for HTML rendering")
	
	// HTML template handler (OLD - using chromedp, DISABLED)
	/*
	var htmlTemplateHandler *http.HTMLTemplateHandler
	if chromedpPrintHandler != nil && htmlTemplateHandlerBridge == nil {
		templatePath := "./application/services/templates/bill_template_optimized.html"
		htmlTemplateHandler = http.NewHTMLTemplateHandler(
			chromedpPrintHandler.GetRenderer(),
			orderRepo,
			shopSettingsRepo,
			templatePath,
		)
		log.Println("✅ HTML template handler (chromedp fallback) initialized")
	}
	*/
	
	// State machine handler
	stateMachineHandler := http.NewStateMachineHandler(smManager)
	menuService := services.NewMenuService(menuRepo)
	menuHandler := http.NewMenuHandler(menuService)
	
	// Menu category service and handler
	menuCategoryRepo := mongodb.NewMenuCategoryRepository(db)
	menuCategoryService := services.NewMenuCategoryService(menuCategoryRepo, menuRepo)
	menuCategoryHandler := http.NewMenuCategoryHandler(menuCategoryService)
	
	// Ingredient service (stockHistoryRepo already initialized above)
	ingredientRestockRepo := mongodb.NewIngredientRestockRepository(db)
	ingredientService := services.NewIngredientService(ingredientRepo, stockHistoryRepo)
	orderService.SetIngredientService(ingredientService)
	ingredientHandler := http.NewIngredientHandler(ingredientService)
	facilityRepo := mongodb.NewFacilityRepository(db)
	facilityService := services.NewFacilityService(facilityRepo)
	facilityHandler := http.NewFacilityHandler(facilityService)
	expenseRepo := mongodb.NewExpenseRepository(db)
	expenseService := services.NewExpenseService(expenseRepo)
	fundExpenseIntegrationService := services.NewFundExpenseIntegrationService(expenseRepo, ingredientRepo, ingredientRestockRepo, facilityRepo, journalService, client)
	expenseHandler := http.NewExpenseHandler(expenseService, fundExpenseIntegrationService)
	ingredientHandler.SetFundExpenseIntegrationService(fundExpenseIntegrationService)
	facilityHandler.SetFundExpenseIntegrationService(fundExpenseIntegrationService)

	attendanceRepo := mongodb.NewHoursEntryRepository(db)
	attendanceService := services.NewAttendanceService(attendanceRepo)
	attendanceHandler := http.NewAttendanceHandler(attendanceService)

	// Schedule service + handler
	scheduleTemplateRepo := mongodb.NewShiftTemplateRepository(db)
	schedulePeriodRepo := mongodb.NewSchedulePeriodRepository(db)
	scheduleSlotRepo := mongodb.NewScheduleSlotRepository(db)
	shiftRegistrationRepo := mongodb.NewShiftRegistrationRepository(db)
	weeklyShiftTemplateRepo := mongodb.NewWeeklyShiftTemplateRepository(db)
	scheduleService := services.NewScheduleService(scheduleTemplateRepo, schedulePeriodRepo, scheduleSlotRepo, shiftRegistrationRepo, weeklyShiftTemplateRepo)
	scheduleHandler := http.NewScheduleHandler(scheduleService)

	// Menu cost and profit analysis handlers
	menuCostHandler := http.NewMenuCostHandler(profitAnalyzerService, costCalculatorService, costRecalculationService, menuRepo)
	profitAnalysisHandler := http.NewProfitAnalysisHandler(profitAnalyzerService)
	operatingExpenseHandler := http.NewOperatingExpenseHandler(operatingExpenseService)
	settingsHandler := http.NewSettingsHandler(shopSettingsService)
	monitoringHandler := http.NewMonitoringHandler(monitoringService)

	// Auto Expense Service - wire up with other services
	autoExpenseService := services.NewAutoExpenseService(expenseService)
	ingredientService.SetAutoExpenseService(autoExpenseService)
	facilityService.SetAutoExpenseService(autoExpenseService)

	// Wire up ingredient service with cost calculator for recalculation triggers
	// Requirements: 1.3, 9.1
	ingredientService.SetCostCalculatorService(costCalculatorService)

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

	// Serve static files for uploads
	r.Static("/uploads", "./uploads")
	
	// Serve templates directory for HTML template preview
	r.Static("/templates", "./templates")

	// Socket.IO endpoint (using standard library)
	r.GET("/socket.io/*any", gin.WrapH(socketIOServer))
	r.POST("/socket.io/*any", gin.WrapH(socketIOServer))
	log.Println("✅ Socket.IO endpoint registered at /socket.io/")

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
				shifts.GET("/:id/pending-orders", shiftHandler.GetPendingOrders)

				// Cash handover routes (waiter)
				shifts.POST("/:id/handover", cashHandoverHandler.CreateHandover)
				// shifts.POST("/:id/handover-and-end", cashHandoverHandler.CreateHandoverAndEndShift) // Not needed - use regular handover
				shifts.GET("/:id/pending-handover", cashHandoverHandler.GetPendingHandover)
				shifts.GET("/:id/handovers", cashHandoverHandler.GetHandoverHistory)
			}
			
			// Cash handover management
			cashHandovers := protected.Group("/cash-handovers")
			{
				// Cashier routes
				cashHandovers.GET("/pending", http.RequireRole(user.RoleCashier, user.RoleManager), cashHandoverHandler.GetPendingHandovers)
				cashHandovers.GET("/all-pending", http.RequireRole(user.RoleCashier, user.RoleManager), cashHandoverHandler.GetAllPendingHandovers)
				cashHandovers.GET("/today", http.RequireRole(user.RoleCashier, user.RoleManager), cashHandoverHandler.GetTodayHandovers)
				cashHandovers.POST("/:id/confirm", http.RequireRole(user.RoleCashier, user.RoleManager), cashHandoverHandler.ConfirmHandover)
				// cashHandovers.POST("/:id/quick-confirm", http.RequireRole(user.RoleCashier, user.RoleManager), cashHandoverHandler.QuickConfirm) // Not needed
				
				// Waiter can cancel their own pending handover
				cashHandovers.DELETE("/:id", cashHandoverHandler.CancelHandover)
				
				// Manager routes - TODO: implement later
				// cashHandovers.GET("/pending-approval", http.RequireRole(user.RoleManager), cashHandoverHandler.GetPendingApprovals)
				// cashHandovers.POST("/:id/approve", http.RequireRole(user.RoleManager), cashHandoverHandler.ApproveDiscrepancy)
				// cashHandovers.GET("/discrepancy-stats", http.RequireRole(user.RoleManager), cashHandoverHandler.GetDiscrepancyStats)
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
				cashierShifts.POST("/:id/cancel-closure", cashierShiftClosureHandler.CancelClosure)
				cashierShifts.POST("/:id/record-actual-cash", cashierShiftClosureHandler.RecordActualCash)
				cashierShifts.POST("/:id/document-variance", cashierShiftClosureHandler.DocumentVariance)
				cashierShifts.POST("/:id/confirm-responsibility", cashierShiftClosureHandler.ConfirmResponsibility)
				cashierShifts.POST("/:id/close", cashierShiftClosureHandler.CloseShift)
				// Frontend-driven: Complete entire closure in one transaction
				cashierShifts.POST("/:id/complete-closure", cashierShiftClosureHandler.CompleteClosure)
				
				// Fund handover - NEW endpoints for Phase 4
				cashierShifts.GET("/:id/managed-funds", cashierShiftClosureHandler.GetManagedFunds)
				cashierShifts.POST("/:id/close-with-fund-handover", cashierShiftClosureHandler.CloseShiftWithFundHandover)
			}
			
			// Cashier shift management - manager only
			cashierShiftsManager := protected.Group("/cashier-shifts")
			cashierShiftsManager.Use(http.RequireRole(user.RoleManager))
			{
				cashierShiftsManager.GET("", cashierShiftHandler.GetAllCashierShifts)
			}
			
			// Batch record routes (barista and manager)
			batchRecords := protected.Group("/batch-records")
			batchRecords.Use(http.RequireRole(user.RoleBarista, user.RoleManager))
			{
				batchRecords.POST("", batchRecordHandler.CreateBatchRecord)
				batchRecords.GET("", batchRecordHandler.GetBatchRecords)
				batchRecords.GET("/:id", batchRecordHandler.GetBatchRecord)
				batchRecords.PATCH("/:id/quantity", batchRecordHandler.UpdateBatchQuantity)
				batchRecords.PATCH("/:id/expire", batchRecordHandler.MarkBatchExpired)
				batchRecords.DELETE("/:id", batchRecordHandler.DeleteBatchRecord)
			}
			
			// Batch usage routes (accessible to all authenticated users for order processing)
			batchUsage := protected.Group("/batch-usage")
			{
				batchUsage.POST("", batchUsageHandler.UseBatch)
				batchUsage.GET("/history", batchUsageHandler.GetUsageHistory)
			}
			
			// Batch alert routes (barista and manager)
			batchAlerts := protected.Group("/batch-alerts")
			batchAlerts.Use(http.RequireRole(user.RoleBarista, user.RoleManager))
			{
				batchAlerts.GET("", batchAlertHandler.GetAlerts)
			}
			
			// Batch report routes (manager only)
			batchReports := protected.Group("/batch-reports")
			batchReports.Use(http.RequireRole(user.RoleManager))
			{
				batchReports.GET("/production", batchReportHandler.GetProductionReport)
				batchReports.GET("/wastage", batchReportHandler.GetWastageReport)
				batchReports.GET("/usage", batchReportHandler.GetUsageReport)
			}
			
			// Print template routes (read access for all authenticated users)
			protected.GET("/print-templates", printTemplateHandler.ListTemplates)
			protected.GET("/print-templates/:id", printTemplateHandler.GetTemplate)
			
			// Settings routes (read access for all authenticated users)
			protected.GET("/settings", settingsHandler.GetSettings)
			protected.PATCH("/settings", settingsHandler.UpdateSettings)
			
			// Reprint routes (only for cashier and manager)
			reprintGroup := protected.Group("/orders/:id")
			reprintGroup.Use(http.RequireRole(user.RoleCashier, user.RoleManager))
			{
				reprintGroup.POST("/reprint-bill", orderHandler.ReprintBill)
				reprintGroup.POST("/reprint-label", orderHandler.ReprintLabel)
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
				waiter.POST("/orders/:id/cancel", orderHandler.CancelOrder)
	waiter.POST("/orders/:id/print-temp-bill", orderHandler.PrintTemporaryBill)
				waiter.POST("/orders/merge", orderHandler.MergeOrders)
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
				
				cashier.POST("/orders/:id/reprint-bill", orderHandler.ReprintBill)
				cashier.POST("/orders/:id/reprint-label", orderHandler.ReprintLabel)
				
				// Shift management
				cashier.POST("/shifts/:id/close", shiftHandler.CloseShift)
				cashier.GET("/shifts", shiftHandler.GetAllShifts)
				cashier.GET("/shifts/:id", shiftHandler.GetShift)
				
				// Cashier-specific routes
				cashier.GET("/shifts/:id/status", cashierHandler.GetShiftStatus)
				cashier.GET("/shifts/:id/payments", cashierHandler.GetPaymentsByShift)
				cashier.GET("/shifts/:id/orders", cashierHandler.GetOrdersByShift)
				cashier.POST("/discrepancies", cashierHandler.ReportDiscrepancy)
				cashier.GET("/discrepancies/pending", cashierHandler.GetPendingDiscrepancies)
				cashier.POST("/discrepancies/:id/resolve", cashierHandler.ResolveDiscrepancy)
				cashier.POST("/reconcile/cash", cashierHandler.ReconcileCash)
				cashier.POST("/orders/:id/override", cashierHandler.OverridePayment)
				cashier.POST("/orders/:id/lock", cashierHandler.LockOrder)
				cashier.GET("/reports/shift/:id", cashierHandler.GenerateShiftReport)
				cashier.GET("/reports/daily", cashierHandler.GetDailyReport)
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
				
				// Menu category routes
				manager.POST("/menu-categories", menuCategoryHandler.CreateCategory)
				manager.GET("/menu-categories", menuCategoryHandler.GetAllCategories)
				manager.GET("/menu-categories/:id", menuCategoryHandler.GetCategory)
				manager.PUT("/menu-categories/:id", menuCategoryHandler.UpdateCategory)
				manager.DELETE("/menu-categories/:id", menuCategoryHandler.DeleteCategory)
				
				// Menu cost and profit analysis routes
				manager.GET("/menu/costs", menuCostHandler.GetMenuCosts)
				manager.GET("/menu/costs/:id", menuCostHandler.GetMenuCostDetail)
				manager.GET("/menu/warnings", menuCostHandler.GetMenuWarnings)
				manager.POST("/menu/costs/recalculate-all", menuCostHandler.RecalculateAllCosts)
				
				// Variant-aware cost analysis routes (Task 6.3)
				manager.GET("/menu/:id/cost-breakdown", menuCostHandler.GetCostBreakdown)
				manager.GET("/menu/:id/profit-analysis", menuCostHandler.GetProfitAnalysis)
				manager.POST("/menu/:id/calculate-cost", menuCostHandler.CalculateCost)
				
				// Profit analysis routes
				manager.GET("/reports/category-profit", profitAnalysisHandler.GetCategoryProfit)
				manager.GET("/reports/operating-profit", profitAnalysisHandler.GetOperatingProfit)
				
				// Operating expense routes
				manager.POST("/operating-expenses", operatingExpenseHandler.CreateOperatingExpense)
				manager.GET("/operating-expenses", operatingExpenseHandler.GetOperatingExpenses)
				
				// Settings routes
				manager.GET("/settings", settingsHandler.GetSettings)
				manager.PATCH("/settings", settingsHandler.UpdateSettings)
				
				// Monitoring and metrics routes
				manager.GET("/monitoring/metrics", monitoringHandler.GetMetrics)
				manager.GET("/monitoring/alerts", monitoringHandler.GetAlerts)
				manager.GET("/monitoring/metrics/aggregated", monitoringHandler.GetAggregatedMetrics)
				manager.GET("/monitoring/health", monitoringHandler.GetHealthStatus)
				
				// Ingredient management routes
				manager.POST("/ingredients", ingredientHandler.CreateIngredient)
				manager.GET("/ingredients", ingredientHandler.GetAllIngredients)
				manager.GET("/ingredients/low-stock", ingredientHandler.GetLowStock)
				manager.GET("/ingredients/summary", ingredientHandler.GetStockSummary)
				manager.GET("/ingredients/:id", ingredientHandler.GetIngredient)
				manager.GET("/ingredients/:id/history", ingredientHandler.GetStockHistory)
				manager.PUT("/ingredients/:id", ingredientHandler.UpdateIngredient)
				manager.DELETE("/ingredients/:id", ingredientHandler.DeleteIngredient)
				
				// Stock operations
				manager.POST("/ingredients/:id/stock-in", ingredientHandler.StockIn)       // Add stock (purchase)
				manager.POST("/ingredients/:id/stock-out", ingredientHandler.StockOut)     // Remove stock (usage/waste)
				manager.POST("/ingredients/:id/stock-adjust", ingredientHandler.StockAdjust) // Set to specific quantity
				manager.POST("/ingredients/:id/adjust", ingredientHandler.AdjustStock)     // Legacy endpoint
				manager.POST("/ingredients/:id/restock/from-fund", ingredientHandler.RestockFromFund) // Restock paid from fund
				manager.GET("/ingredients/:id/restock-history", ingredientHandler.GetRestockHistory)  // Restock history
				
				// Ingredient category routes
				manager.POST("/ingredient-categories", ingredientHandler.CreateCategory)
				manager.GET("/ingredient-categories", ingredientHandler.GetCategories)
				manager.DELETE("/ingredient-categories/:id", ingredientHandler.DeleteCategory)
				
				// Batch definition routes (manager only)
				manager.POST("/batch-definitions", batchDefinitionHandler.CreateBatchDefinition)
				manager.GET("/batch-definitions", batchDefinitionHandler.GetBatchDefinitions)
				manager.GET("/batch-definitions/:id", batchDefinitionHandler.GetBatchDefinition)
				manager.PUT("/batch-definitions/:id", batchDefinitionHandler.UpdateBatchDefinition)
				manager.DELETE("/batch-definitions/:id", batchDefinitionHandler.DeleteBatchDefinition)
				
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
				manager.POST("/expenses/from-fund", expenseHandler.CreateExpenseFromFund) // Operating expense from fund
				manager.GET("/expenses", expenseHandler.GetExpenses)
				manager.PUT("/expenses/:id", expenseHandler.UpdateExpense)
				manager.DELETE("/expenses/:id", expenseHandler.DeleteExpense)
				manager.POST("/expense-categories", expenseHandler.CreateCategory)
				
				// Fund management routes (journal / double-entry source of truth)
				manager.GET("/fund/journal-balances", fundHandler.GetJournalFundBalances)
				manager.GET("/fund/journal", fundHandler.GetJournalEntries)
				manager.GET("/fund/journal/:id", fundHandler.GetJournalEntry)
				manager.POST("/fund/deposit", fundHandler.Deposit)
				manager.POST("/fund/withdraw", fundHandler.Withdraw)
				manager.POST("/fund/transfer", fundHandler.Transfer)
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
				
				manager.POST("/orders/:id/reprint-bill", orderHandler.ReprintBill)
				manager.POST("/orders/:id/reprint-label", orderHandler.ReprintLabel)
				
				// Print job management routes
				manager.GET("/print-jobs", printJobHandler.ListPrintJobs)
				manager.GET("/print-jobs/pending", printJobHandler.GetPendingJobs)
				manager.GET("/print-jobs/failed", printJobHandler.GetFailedJobs)
				manager.GET("/print-jobs/:id", printJobHandler.GetPrintJob)
				manager.POST("/print-jobs/:id/retry", printJobHandler.RetryPrintJob)
				manager.DELETE("/print-jobs/:id", printJobHandler.CancelPrintJob)
				manager.PUT("/print-jobs/:id/status", printJobHandler.UpdatePrintJobStatus) // For local print bridge
				
				// Printer configuration routes
				manager.GET("/printers", printerConfigHandler.ListPrinters)
				manager.GET("/printers/:id", printerConfigHandler.GetPrinter)
				manager.POST("/printers", printerConfigHandler.CreatePrinter)
				manager.PUT("/printers/:id", printerConfigHandler.UpdatePrinter)
				manager.DELETE("/printers/:id", printerConfigHandler.DeletePrinter)
				manager.POST("/printers/:id/test", printerConfigHandler.TestConnection)
				
				// Print template routes
				manager.GET("/print-templates", printTemplateHandler.ListTemplates)
				manager.GET("/print-templates/:id", printTemplateHandler.GetTemplate)
				manager.POST("/print-templates", printTemplateHandler.CreateTemplate)
				manager.PUT("/print-templates/:id", printTemplateHandler.UpdateTemplate)
				manager.DELETE("/print-templates/:id", printTemplateHandler.DeleteTemplate)
				manager.POST("/print-templates/:id/preview", printTemplateHandler.PreviewTemplate)
				
				// Visual print routes - DISABLED (moved to print bridge)
				// if visualPrintHandler != nil {
				// 	manager.POST("/visual-print/bill", visualPrintHandler.PrintVisualBill)
				// 	manager.GET("/visual-print/preview/:order_id", visualPrintHandler.PreviewVisualBill)
				// }
				
				// Chromedp print routes - DISABLED (UI removed)
				// if chromedpPrintHandler != nil {
				// 	manager.POST("/chromedp-print/bill", chromedpPrintHandler.PrintChromedpBill)
				// 	manager.GET("/chromedp-print/preview/:order_id", chromedpPrintHandler.PreviewChromedpBill)
				// }
				
				// HTML template management routes
				// Using print bridge (URL from settings)
				if htmlTemplateHandlerBridge != nil {
					manager.GET("/html-templates/bill", htmlTemplateHandlerBridge.GetHTMLTemplate)
					manager.PUT("/html-templates/bill", htmlTemplateHandlerBridge.UpdateHTMLTemplate)
					manager.POST("/html-templates/test-print", htmlTemplateHandlerBridge.TestPrintHTMLTemplate)
					manager.POST("/html-templates/preview", htmlTemplateHandlerBridge.PreviewHTMLTemplate)
					log.Println("✅ HTML template routes registered (using print bridge from settings)")
				}
				
				// Label template management routes
				if labelTemplateHandler != nil {
					manager.GET("/label-templates/order-item", labelTemplateHandler.GetLabelTemplate)
					manager.PUT("/label-templates/order-item", labelTemplateHandler.UpdateLabelTemplate)
					manager.POST("/label-templates/test-print", labelTemplateHandler.TestPrintLabel)
					log.Println("✅ Label template routes registered")
				}
				
				// Shop settings routes
				manager.GET("/shop-settings", shopSettingsHandler.GetSettings)
				manager.POST("/shop-settings", shopSettingsHandler.CreateSettings)
				manager.PUT("/shop-settings/:id", shopSettingsHandler.UpdateSettings)
				
				// Print bridge test route
				manager.POST("/print-bridge/test", shopSettingsHandler.TestPrintBridge)
				
				// Logo upload routes
				manager.POST("/settings/logo", logoUploadHandler.UploadLogo)
				manager.DELETE("/settings/logo", logoUploadHandler.DeleteLogo)
				
				// Shift management routes
				manager.GET("/shifts", shiftHandler.GetAllShifts)
				manager.GET("/shifts/:id", shiftHandler.GetShift)

				// Schedule management routes (manager only)
				manager.POST("/schedule/templates", scheduleHandler.CreateTemplate)
				manager.GET("/schedule/templates", scheduleHandler.GetTemplates)
				manager.PUT("/schedule/templates/:id", scheduleHandler.UpdateTemplate)
				manager.DELETE("/schedule/templates/:id", scheduleHandler.DeleteTemplate)

				manager.POST("/schedule/weekly-templates", scheduleHandler.CreateWeeklyTemplate)
				manager.GET("/schedule/weekly-templates", scheduleHandler.GetWeeklyTemplates)
				manager.PUT("/schedule/weekly-templates/:id", scheduleHandler.UpdateWeeklyTemplate)
				manager.DELETE("/schedule/weekly-templates/:id", scheduleHandler.DeleteWeeklyTemplate)

				manager.POST("/schedule/periods", scheduleHandler.CreatePeriod)
				manager.GET("/schedule/periods", scheduleHandler.GetPeriods)
				manager.PATCH("/schedule/periods/:id/status", scheduleHandler.SetPeriodStatus)
				manager.GET("/schedule/periods/:id", scheduleHandler.GetPeriodDetail)
				manager.POST("/schedule/periods/:id/slots", scheduleHandler.AddSlot)
				manager.DELETE("/schedule/periods/:id/slots/:slotId", scheduleHandler.RemoveSlot)
				manager.DELETE("/schedule/registrations/:regId", scheduleHandler.CancelRegistration)
				manager.GET("/schedule/registrations-by-date", scheduleHandler.GetRegistrationsByDate)
				manager.POST("/schedule/registrations-by-period", scheduleHandler.GetRegistrationsByPeriod)

				// Attendance (actual work hours) routes
				manager.POST("/attendance", attendanceHandler.AddEntry)
				manager.POST("/attendance/bulk", attendanceHandler.BulkAddEntries)
				manager.GET("/attendance", attendanceHandler.GetEntries)
				manager.GET("/attendance/summary", attendanceHandler.GetStaffSummary)
				manager.PATCH("/attendance/:id", attendanceHandler.UpdateEntry)
				manager.DELETE("/attendance/:id", attendanceHandler.DeleteEntry)
			}
		}

		// Schedule routes for all authenticated staff (view + register)
		scheduleStaff := protected.Group("/schedule")
		{
			scheduleStaff.GET("/current", scheduleHandler.GetCurrentPeriod)
			scheduleStaff.GET("/my", scheduleHandler.GetMySchedule)
			scheduleStaff.POST("/slots/:slotId/register", scheduleHandler.Register)
			scheduleStaff.DELETE("/slots/:slotId/register", scheduleHandler.Unregister)
		}
	}

	// Create default admin user
	createDefaultUsers(authService, userRepo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on :%s", port)
	r.Run(":" + port)
}

func createDefaultUsers(authService *services.AuthService, userRepo *mongodb.UserRepository) {
	ctx := context.Background()
	
	// Check if admin exists
	if _, err := userRepo.FindByUsername(ctx, "admin"); err == nil {
		return
	}

	// Create only admin user
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
	log.Printf("Created user: admin (role: %s)", user.RoleManager)
}

// Reprint routes (only for cashier and manager)
