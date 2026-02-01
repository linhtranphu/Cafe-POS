package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"cafe-pos/backend/domain/expense"
	"cafe-pos/backend/domain/facility"
	"cafe-pos/backend/domain/ingredient"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AutoExpenseService handles automatic expense tracking for purchases and maintenance
type AutoExpenseService struct {
	expenseService *ExpenseService
	categoryCache  map[string]primitive.ObjectID
	cacheMutex     sync.RWMutex
}

// NewAutoExpenseService creates a new AutoExpenseService instance
func NewAutoExpenseService(expenseService *ExpenseService) *AutoExpenseService {
	return &AutoExpenseService{
		expenseService: expenseService,
		categoryCache:  make(map[string]primitive.ObjectID),
	}
}

// GetOrCreateCategory gets category ID by name, creates if not exists
// Uses caching to minimize database queries
func (s *AutoExpenseService) GetOrCreateCategory(ctx context.Context, categoryName string) (primitive.ObjectID, error) {
	// Check cache first (read lock)
	s.cacheMutex.RLock()
	if id, exists := s.categoryCache[categoryName]; exists {
		s.cacheMutex.RUnlock()
		return id, nil
	}
	s.cacheMutex.RUnlock()

	// Not in cache, query database
	categories, err := s.expenseService.GetCategories(ctx)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to get categories: %w", err)
	}

	// Look for category in database results
	for _, cat := range categories {
		if cat.Name == categoryName {
			// Update cache (write lock)
			s.cacheMutex.Lock()
			s.categoryCache[categoryName] = cat.ID
			s.cacheMutex.Unlock()
			return cat.ID, nil
		}
	}

	// Category doesn't exist, create it
	newCategory := &expense.Category{
		Name: categoryName,
	}

	if err := s.expenseService.CreateCategory(ctx, newCategory); err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to create category: %w", err)
	}

	// Update cache (write lock)
	s.cacheMutex.Lock()
	s.categoryCache[categoryName] = newCategory.ID
	s.cacheMutex.Unlock()

	log.Printf("[AutoExpense] Created new category: %s (ID: %s)", categoryName, newCategory.ID.Hex())
	return newCategory.ID, nil
}

// TrackIngredientPurchase creates an expense record for ingredient purchase
// This is called when creating a new ingredient or adjusting stock IN
func (s *AutoExpenseService) TrackIngredientPurchase(ctx context.Context, ing *ingredient.Ingredient, quantity float64, username string) error {
	// Skip if no cost or quantity
	if ing.CostPerUnit <= 0 || quantity <= 0 {
		log.Printf("[AutoExpense] Skipping ingredient purchase tracking: zero cost or quantity (ingredient: %s)", ing.Name)
		return nil
	}

	// Calculate total amount
	amount := ing.CostPerUnit * quantity

	// Get or create category
	categoryID, err := s.GetOrCreateCategory(ctx, expense.CategoryIngredient)
	if err != nil {
		log.Printf("[AutoExpense] Failed to get/create category for ingredient: %v", err)
		return err
	}

	// Create expense record
	exp := &expense.Expense{
		Date:          time.Now(),
		CategoryID:    categoryID,
		Amount:        amount,
		Description:   fmt.Sprintf("Nhập nguyên liệu: %s", ing.Name),
		PaymentMethod: expense.PaymentMethodCash, // Default to cash
		Vendor:        ing.Supplier,
		Notes:         fmt.Sprintf("Số lượng: %.2f %s", quantity, ing.Unit),
		SourceType:    expense.SourceTypeIngredient,
		SourceID:      ing.ID,
		CreatedBy:     username, // Set to ingredient creator
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.expenseService.CreateExpense(ctx, exp); err != nil {
		log.Printf("[AutoExpense] Failed to create expense for ingredient purchase: %v", err)
		return err
	}

	log.Printf("[AutoExpense] Tracked ingredient purchase: %s (%.2f %s) - Amount: %.2f VND", 
		ing.Name, quantity, ing.Unit, amount)
	return nil
}

// TrackFacilityPurchase creates an expense record for facility purchase
// This is called when creating a new facility
func (s *AutoExpenseService) TrackFacilityPurchase(ctx context.Context, fac *facility.Facility, username string) error {
	// Skip if no cost
	if fac.Cost <= 0 {
		log.Printf("[AutoExpense] Skipping facility purchase tracking: zero cost (facility: %s)", fac.Name)
		return nil
	}

	// Get or create category
	categoryID, err := s.GetOrCreateCategory(ctx, expense.CategoryFacility)
	if err != nil {
		log.Printf("[AutoExpense] Failed to get/create category for facility: %v", err)
		return err
	}

	// Create expense record
	exp := &expense.Expense{
		Date:          fac.PurchaseDate,
		CategoryID:    categoryID,
		Amount:        fac.Cost,
		Description:   fmt.Sprintf("Mua thiết bị: %s", fac.Name),
		PaymentMethod: expense.PaymentMethodCash, // Default to cash
		Vendor:        fac.Supplier,
		Notes:         fmt.Sprintf("Loại: %s, Khu vực: %s, Số lượng: %d", fac.Type, fac.Area, fac.Quantity),
		SourceType:    expense.SourceTypeFacility,
		SourceID:      fac.ID,
		CreatedBy:     username, // Set to facility creator
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.expenseService.CreateExpense(ctx, exp); err != nil {
		log.Printf("[AutoExpense] Failed to create expense for facility purchase: %v", err)
		return err
	}

	log.Printf("[AutoExpense] Tracked facility purchase: %s - Amount: %.2f VND", fac.Name, fac.Cost)
	return nil
}

// TrackMaintenance creates an expense record for facility maintenance
// This is called when creating a maintenance record
func (s *AutoExpenseService) TrackMaintenance(ctx context.Context, facilityID primitive.ObjectID, facilityName string, cost float64, maintenanceDate time.Time, notes string, username string) error {
	// Skip if no cost
	if cost <= 0 {
		log.Printf("[AutoExpense] Skipping maintenance tracking: zero cost (facility: %s)", facilityName)
		return nil
	}

	// Get or create category
	categoryID, err := s.GetOrCreateCategory(ctx, expense.CategoryMaintenance)
	if err != nil {
		log.Printf("[AutoExpense] Failed to get/create category for maintenance: %v", err)
		return err
	}

	// Create expense record
	exp := &expense.Expense{
		Date:          maintenanceDate,
		CategoryID:    categoryID,
		Amount:        cost,
		Description:   fmt.Sprintf("Bảo trì: %s", facilityName),
		PaymentMethod: expense.PaymentMethodCash, // Default to cash
		Vendor:        "",                         // Vendor info might be in notes
		Notes:         notes,
		SourceType:    expense.SourceTypeMaintenance,
		SourceID:      facilityID,
		CreatedBy:     username, // Set to maintenance creator
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.expenseService.CreateExpense(ctx, exp); err != nil {
		log.Printf("[AutoExpense] Failed to create expense for maintenance: %v", err)
		return err
	}

	log.Printf("[AutoExpense] Tracked maintenance: %s - Amount: %.2f VND", facilityName, cost)
	return nil
}

// ClearCache clears the category cache
// Useful for testing or when categories are modified externally
func (s *AutoExpenseService) ClearCache() {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()
	s.categoryCache = make(map[string]primitive.ObjectID)
	log.Println("[AutoExpense] Category cache cleared")
}

// GetCacheSize returns the number of cached categories
// Useful for monitoring and debugging
func (s *AutoExpenseService) GetCacheSize() int {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()
	return len(s.categoryCache)
}
