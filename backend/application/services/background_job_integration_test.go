package services

import (
	"context"
	"testing"
	"time"

	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Feature: menu-cost-profit-analysis, Integration Test: Background Job Processing
// **Validates: Requirements 1.3, 9.1, 9.3**
//
// This test verifies the complete end-to-end flow:
// 1. Update ingredient cost
// 2. Verify job is queued
// 3. Wait for job completion
// 4. Verify menu item costs are updated
func TestBackgroundJobProcessing_Integration(t *testing.T) {
	ctx := context.Background()

	// Setup repositories
	menuRepo := &mockMenuRepository{menuItems: make(map[primitive.ObjectID]*menu.MenuItem)}
	ingredientRepo := &mockIngredientRepository{ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient)}
	orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
	orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}
	stockHistoryRepo := &mockStockHistoryRepository{histories: make(map[primitive.ObjectID][]*ingredient.StockHistory)}

	// Create test ingredient
	testIngredient := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Coffee Beans",
		CostPerUnit:       100.0,
		Unit:              "kg",
		Quantity:          50.0,
		ConversionRate:    1.0,
		WastagePercentage: 0.0,
	}
	ingredientRepo.ingredients[testIngredient.ID] = testIngredient

	// Create menu items that use this ingredient
	menuItem1 := &menu.MenuItem{
		ID:    primitive.NewObjectID(),
		Name:  "Espresso",
		Price: 50000,
		Ingredients: []menu.Ingredient{
			{Name: "Coffee Beans", Quantity: 0.02, Unit: ingredient.UnitKilogram},
		},
		CurrentCost: 0, // Not yet calculated
	}
	menuItem2 := &menu.MenuItem{
		ID:    primitive.NewObjectID(),
		Name:  "Cappuccino",
		Price: 60000,
		Ingredients: []menu.Ingredient{
			{Name: "Coffee Beans", Quantity: 0.03, Unit: ingredient.UnitKilogram},
		},
		CurrentCost: 0, // Not yet calculated
	}
	menuRepo.menuItems[menuItem1.ID] = menuItem1
	menuRepo.menuItems[menuItem2.ID] = menuItem2

	// Create services
	costCalculatorService := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)
	recalcService := NewCostRecalculationService(costCalculatorService, menuRepo, 2, 100) // 2 workers
	costCalculatorService.SetCostRecalculationService(recalcService)
	
	ingredientService := NewIngredientService(ingredientRepo, stockHistoryRepo)
	ingredientService.SetCostCalculatorService(costCalculatorService)

	// Start the recalculation worker pool
	recalcService.Start()
	defer recalcService.Stop()

	// Step 1: Update ingredient cost (this should trigger background job)
	newCost := 150.0 // Update cost from 100 to 150
	_, err := ingredientService.UpdateIngredient(ctx, testIngredient.ID, &ingredient.UpdateIngredientRequest{
		Name:        testIngredient.Name,
		CostPerUnit: &newCost,
		Unit:        testIngredient.Unit,
		Quantity:    &testIngredient.Quantity,
	})
	if err != nil {
		t.Fatalf("Failed to update ingredient: %v", err)
	}

	// Step 2: Verify job is queued (check status immediately)
	// Give a moment for the goroutine to queue the job
	time.Sleep(100 * time.Millisecond)
	
	status, err := recalcService.GetRecalculationStatus(ctx)
	if err != nil {
		t.Fatalf("Failed to get recalculation status: %v", err)
	}
	
	// At this point, jobs should be queued or processing
	if status.QueuedItems == 0 && status.ProcessedItems == 0 {
		t.Errorf("Expected jobs to be queued or processing, but got queued=%d, processed=%d", 
			status.QueuedItems, status.ProcessedItems)
	}

	// Step 3: Wait for job completion
	maxWaitTime := 5 * time.Second
	startTime := time.Now()
	jobsCompleted := false
	
	for {
		status, err := recalcService.GetRecalculationStatus(ctx)
		if err != nil {
			t.Fatalf("Failed to get recalculation status: %v", err)
		}
		
		// Jobs are complete when queue is empty
		if status.QueuedItems == 0 {
			jobsCompleted = true
			break
		}
		
		if time.Since(startTime) > maxWaitTime {
			t.Fatalf("Timeout waiting for jobs to complete. Status: queued=%d, processed=%d", 
				status.QueuedItems, status.ProcessedItems)
		}
		
		time.Sleep(50 * time.Millisecond)
	}
	
	if !jobsCompleted {
		t.Fatal("Jobs did not complete within timeout")
	}

	// Give workers a moment to finish writing
	time.Sleep(100 * time.Millisecond)

	// Step 4: Verify menu item costs are updated
	updatedMenuItem1, err := menuRepo.FindByID(ctx, menuItem1.ID)
	if err != nil {
		t.Fatalf("Failed to find menu item 1: %v", err)
	}
	
	updatedMenuItem2, err := menuRepo.FindByID(ctx, menuItem2.ID)
	if err != nil {
		t.Fatalf("Failed to find menu item 2: %v", err)
	}

	// Verify costs were calculated
	if updatedMenuItem1.CurrentCost == 0 {
		t.Errorf("Menu item 1 cost was not updated (still 0)")
	}
	if updatedMenuItem2.CurrentCost == 0 {
		t.Errorf("Menu item 2 cost was not updated (still 0)")
	}

	// Verify costs are correct based on new ingredient cost (150)
	// Espresso: 0.02 kg * 150 = 3.0
	expectedCost1 := 3.0
	if updatedMenuItem1.CurrentCost != expectedCost1 {
		t.Errorf("Menu item 1 cost incorrect: expected %.2f, got %.2f", 
			expectedCost1, updatedMenuItem1.CurrentCost)
	}

	// Cappuccino: 0.03 kg * 150 = 4.5
	expectedCost2 := 4.5
	if updatedMenuItem2.CurrentCost != expectedCost2 {
		t.Errorf("Menu item 2 cost incorrect: expected %.2f, got %.2f", 
			expectedCost2, updatedMenuItem2.CurrentCost)
	}

	// Verify cost status is FINAL
	if updatedMenuItem1.CostStatus != "FINAL" {
		t.Errorf("Menu item 1 cost status incorrect: expected FINAL, got %s", 
			updatedMenuItem1.CostStatus)
	}
	if updatedMenuItem2.CostStatus != "FINAL" {
		t.Errorf("Menu item 2 cost status incorrect: expected FINAL, got %s", 
			updatedMenuItem2.CostStatus)
	}

	t.Logf("Integration test passed: %d menu items updated successfully", 2)
}

// Test that multiple ingredient updates are handled correctly
func TestBackgroundJobProcessing_MultipleUpdates(t *testing.T) {
	ctx := context.Background()

	// Setup repositories
	menuRepo := &mockMenuRepository{menuItems: make(map[primitive.ObjectID]*menu.MenuItem)}
	ingredientRepo := &mockIngredientRepository{ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient)}
	orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
	orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}
	stockHistoryRepo := &mockStockHistoryRepository{histories: make(map[primitive.ObjectID][]*ingredient.StockHistory)}

	// Create two ingredients
	ingredient1 := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Coffee",
		CostPerUnit:       100.0,
		Unit:              "kg",
		Quantity:          50.0,
		ConversionRate:    1.0,
		WastagePercentage: 0.0,
	}
	ingredient2 := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Milk",
		CostPerUnit:       50.0,
		Unit:              "liter",
		Quantity:          30.0,
		ConversionRate:    1.0,
		WastagePercentage: 0.0,
	}
	ingredientRepo.ingredients[ingredient1.ID] = ingredient1
	ingredientRepo.ingredients[ingredient2.ID] = ingredient2

	// Create menu item that uses both ingredients
	menuItem := &menu.MenuItem{
		ID:    primitive.NewObjectID(),
		Name:  "Latte",
		Price: 70000,
		Ingredients: []menu.Ingredient{
			{Name: "Coffee", Quantity: 0.02, Unit: ingredient.UnitKilogram},
			{Name: "Milk", Quantity: 0.2, Unit: ingredient.UnitLiter},
		},
		CurrentCost: 0,
	}
	menuRepo.menuItems[menuItem.ID] = menuItem

	// Create services
	costCalculatorService := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)
	recalcService := NewCostRecalculationService(costCalculatorService, menuRepo, 2, 100)
	costCalculatorService.SetCostRecalculationService(recalcService)
	
	ingredientService := NewIngredientService(ingredientRepo, stockHistoryRepo)
	ingredientService.SetCostCalculatorService(costCalculatorService)

	recalcService.Start()
	defer recalcService.Stop()

	// Update both ingredients
	newCost1 := 150.0
	_, err := ingredientService.UpdateIngredient(ctx, ingredient1.ID, &ingredient.UpdateIngredientRequest{
		Name:        ingredient1.Name,
		CostPerUnit: &newCost1,
		Unit:        ingredient1.Unit,
		Quantity:    &ingredient1.Quantity,
	})
	if err != nil {
		t.Fatalf("Failed to update ingredient 1: %v", err)
	}

	newCost2 := 60.0
	_, err = ingredientService.UpdateIngredient(ctx, ingredient2.ID, &ingredient.UpdateIngredientRequest{
		Name:        ingredient2.Name,
		CostPerUnit: &newCost2,
		Unit:        ingredient2.Unit,
		Quantity:    &ingredient2.Quantity,
	})
	if err != nil {
		t.Fatalf("Failed to update ingredient 2: %v", err)
	}

	// Wait for all jobs to complete
	maxWaitTime := 5 * time.Second
	startTime := time.Now()
	
	for {
		status, _ := recalcService.GetRecalculationStatus(ctx)
		if status.QueuedItems == 0 {
			break
		}
		if time.Since(startTime) > maxWaitTime {
			t.Fatalf("Timeout waiting for jobs to complete")
		}
		time.Sleep(50 * time.Millisecond)
	}

	time.Sleep(100 * time.Millisecond)

	// Verify menu item cost reflects both ingredient updates
	updatedMenuItem, err := menuRepo.FindByID(ctx, menuItem.ID)
	if err != nil {
		t.Fatalf("Failed to find menu item: %v", err)
	}

	// Expected: (0.02 * 150) + (0.2 * 60) = 3.0 + 12.0 = 15.0
	expectedCost := 15.0
	if updatedMenuItem.CurrentCost != expectedCost {
		t.Errorf("Menu item cost incorrect: expected %.2f, got %.2f", 
			expectedCost, updatedMenuItem.CurrentCost)
	}

	t.Logf("Multiple updates test passed: final cost = %.2f", updatedMenuItem.CurrentCost)
}

// Mock stock history repository for integration tests
type mockStockHistoryRepository struct {
	histories map[primitive.ObjectID][]*ingredient.StockHistory
}

func (m *mockStockHistoryRepository) Create(ctx context.Context, history *ingredient.StockHistory) error {
	if history.ID.IsZero() {
		history.ID = primitive.NewObjectID()
	}
	m.histories[history.IngredientID] = append(m.histories[history.IngredientID], history)
	return nil
}

func (m *mockStockHistoryRepository) FindByIngredientID(ctx context.Context, ingredientID primitive.ObjectID) ([]*ingredient.StockHistory, error) {
	return m.histories[ingredientID], nil
}
