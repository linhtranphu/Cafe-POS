// +build ignore

package main

import (
	"context"
	"fmt"
	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock repositories for testing
type mockMenuRepository struct {
	menuItems map[primitive.ObjectID]*menu.MenuItem
}

func (m *mockMenuRepository) Create(ctx context.Context, item *menu.MenuItem) error {
	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}
	m.menuItems[item.ID] = item
	return nil
}

func (m *mockMenuRepository) FindAll(ctx context.Context) ([]*menu.MenuItem, error) {
	items := make([]*menu.MenuItem, 0, len(m.menuItems))
	for _, item := range m.menuItems {
		items = append(items, item)
	}
	return items, nil
}

func (m *mockMenuRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*menu.MenuItem, error) {
	item, exists := m.menuItems[id]
	if !exists {
		return nil, fmt.Errorf("menu item not found")
	}
	return item, nil
}

func (m *mockMenuRepository) Update(ctx context.Context, id primitive.ObjectID, item *menu.MenuItem) error {
	m.menuItems[id] = item
	return nil
}

func (m *mockMenuRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	delete(m.menuItems, id)
	return nil
}

type mockIngredientRepository struct {
	ingredients map[string]*ingredient.Ingredient
}

func (m *mockIngredientRepository) FindAll(ctx context.Context) ([]*ingredient.Ingredient, error) {
	items := make([]*ingredient.Ingredient, 0, len(m.ingredients))
	for _, item := range m.ingredients {
		items = append(items, item)
	}
	return items, nil
}

func (m *mockIngredientRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*ingredient.Ingredient, error) {
	for _, item := range m.ingredients {
		if item.ID == id {
			return item, nil
		}
	}
	return nil, fmt.Errorf("ingredient not found")
}

func (m *mockIngredientRepository) Create(ctx context.Context, item *ingredient.Ingredient) error {
	return nil
}

func (m *mockIngredientRepository) Update(ctx context.Context, id primitive.ObjectID, item *ingredient.Ingredient) error {
	return nil
}

func (m *mockIngredientRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func (m *mockIngredientRepository) FindLowStock(ctx context.Context) ([]*ingredient.Ingredient, error) {
	return nil, nil
}

func (m *mockIngredientRepository) CreateCategory(ctx context.Context, cat *ingredient.IngredientCategory) error {
	return nil
}

func (m *mockIngredientRepository) GetCategories(ctx context.Context) ([]ingredient.IngredientCategory, error) {
	return nil, nil
}

func (m *mockIngredientRepository) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func main() {
	// Setup test data
	menuRepo := &mockMenuRepository{
		menuItems: make(map[primitive.ObjectID]*menu.MenuItem),
	}

	ingredientRepo := &mockIngredientRepository{
		ingredients: map[string]*ingredient.Ingredient{
			"Espresso": {
				ID:                primitive.NewObjectID(),
				Name:              "Espresso",
				CostPerUnit:       200.0,
				ConversionRate:    1.0,
				WastagePercentage: 5.0,
			},
			"Milk": {
				ID:                primitive.NewObjectID(),
				Name:              "Milk",
				CostPerUnit:       50.0,
				ConversionRate:    1.0,
				WastagePercentage: 10.0,
			},
		},
	}

	// Create service
	service := services.NewCostCalculatorService(menuRepo, ingredientRepo)

	// Create a test menu item
	menuItem := &menu.MenuItem{
		ID:   primitive.NewObjectID(),
		Name: "Cappuccino",
		Ingredients: []menu.Ingredient{
			{Name: "Espresso", Quantity: 30, Unit: "ml"},
			{Name: "Milk", Quantity: 150, Unit: "ml"},
		},
	}
	menuRepo.menuItems[menuItem.ID] = menuItem

	// Test 1: Basic calculation
	fmt.Println("Test 1: Basic cost calculation")
	result, err := service.CalculateMenuItemCost(context.Background(), menuItem.ID)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("✅ Menu Item: %s\n", menuItem.Name)
	fmt.Printf("   Current Cost: %.2f\n", result.CurrentCost)
	fmt.Printf("   Cost Status: %s\n", result.CostStatus)
	fmt.Printf("   Expected: 14550.00 (30*200*1.0*1.05 + 150*50*1.0*1.10)\n")
	if result.CurrentCost == 14550.0 && result.CostStatus == menu.CostStatusFinal {
		fmt.Println("   ✅ PASS")
	} else {
		fmt.Println("   ❌ FAIL")
	}

	// Test 2: Missing ingredient cost
	fmt.Println("\nTest 2: Missing ingredient cost")
	ingredientRepo.ingredients["Chocolate"] = &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Chocolate",
		CostPerUnit:       0.0, // Missing cost
		ConversionRate:    1.0,
		WastagePercentage: 0.0,
	}
	
	menuItem2 := &menu.MenuItem{
		ID:   primitive.NewObjectID(),
		Name: "Mocha",
		Ingredients: []menu.Ingredient{
			{Name: "Espresso", Quantity: 30, Unit: "ml"},
			{Name: "Chocolate", Quantity: 20, Unit: "g"},
		},
	}
	menuRepo.menuItems[menuItem2.ID] = menuItem2

	result2, err := service.CalculateMenuItemCost(context.Background(), menuItem2.ID)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("✅ Menu Item: %s\n", menuItem2.Name)
	fmt.Printf("   Current Cost: %.2f\n", result2.CurrentCost)
	fmt.Printf("   Cost Status: %s\n", result2.CostStatus)
	fmt.Printf("   Missing Ingredients: %v\n", result2.MissingIngredients)
	fmt.Printf("   Expected: 6300.00 with INCOMPLETE status\n")
	if result2.CurrentCost == 6300.0 && result2.CostStatus == menu.CostStatusIncomplete && len(result2.MissingIngredients) == 1 {
		fmt.Println("   ✅ PASS")
	} else {
		fmt.Println("   ❌ FAIL")
	}

	// Test 3: No ingredients
	fmt.Println("\nTest 3: No ingredients")
	menuItem3 := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Service Item",
		Ingredients: []menu.Ingredient{},
	}
	menuRepo.menuItems[menuItem3.ID] = menuItem3

	result3, err := service.CalculateMenuItemCost(context.Background(), menuItem3.ID)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("✅ Menu Item: %s\n", menuItem3.Name)
	fmt.Printf("   Current Cost: %.2f\n", result3.CurrentCost)
	fmt.Printf("   Cost Status: %s\n", result3.CostStatus)
	fmt.Printf("   Expected: 0.00 with FINAL status\n")
	if result3.CurrentCost == 0.0 && result3.CostStatus == menu.CostStatusFinal {
		fmt.Println("   ✅ PASS")
	} else {
		fmt.Println("   ❌ FAIL")
	}

	fmt.Println("\n✅ All tests completed!")
}
