package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/expense"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock repository for testing
type mockOperatingExpenseRepository struct {
	expenses []*expense.OperatingExpense
}

func (m *mockOperatingExpenseRepository) Create(ctx context.Context, operatingExpense *expense.OperatingExpense) error {
	operatingExpense.ID = primitive.NewObjectID()
	operatingExpense.CreatedAt = time.Now()
	operatingExpense.UpdatedAt = time.Now()
	m.expenses = append(m.expenses, operatingExpense)
	return nil
}

func (m *mockOperatingExpenseRepository) Update(ctx context.Context, id primitive.ObjectID, operatingExpense *expense.OperatingExpense) error {
	for i, exp := range m.expenses {
		if exp.ID == id {
			operatingExpense.ID = id
			operatingExpense.CreatedAt = exp.CreatedAt
			operatingExpense.UpdatedAt = time.Now()
			m.expenses[i] = operatingExpense
			return nil
		}
	}
	return nil
}

func (m *mockOperatingExpenseRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*expense.OperatingExpense, error) {
	for _, exp := range m.expenses {
		if exp.ID == id {
			return exp, nil
		}
	}
	return nil, nil
}

func (m *mockOperatingExpenseRepository) FindByPeriod(ctx context.Context, startDate, endDate time.Time) ([]*expense.OperatingExpense, error) {
	var result []*expense.OperatingExpense
	for _, exp := range m.expenses {
		if !exp.PeriodStart.After(endDate) && !exp.PeriodEnd.Before(startDate) {
			result = append(result, exp)
		}
	}
	return result, nil
}

func (m *mockOperatingExpenseRepository) FindForDate(ctx context.Context, date time.Time) (*expense.OperatingExpense, error) {
	for _, exp := range m.expenses {
		if !date.Before(exp.PeriodStart) && !date.After(exp.PeriodEnd) {
			return exp, nil
		}
	}
	return nil, nil
}

func (m *mockOperatingExpenseRepository) FindAll(ctx context.Context, startDate, endDate *time.Time) ([]*expense.OperatingExpense, error) {
	if startDate == nil || endDate == nil {
		return m.expenses, nil
	}
	return m.FindByPeriod(ctx, *startDate, *endDate)
}

func (m *mockOperatingExpenseRepository) Upsert(ctx context.Context, operatingExpense *expense.OperatingExpense) (*expense.OperatingExpense, error) {
	// Check if expense exists for this period
	for _, exp := range m.expenses {
		if exp.PeriodStart.Equal(operatingExpense.PeriodStart) && exp.PeriodEnd.Equal(operatingExpense.PeriodEnd) {
			// Update existing
			m.Update(ctx, exp.ID, operatingExpense)
			return operatingExpense, nil
		}
	}
	// Create new
	m.Create(ctx, operatingExpense)
	return operatingExpense, nil
}

func (m *mockOperatingExpenseRepository) FindByDateRange(ctx context.Context, start, end time.Time) ([]*expense.OperatingExpense, error) {
	return m.FindByPeriod(ctx, start, end)
}

func (m *mockOperatingExpenseRepository) FindByDate(ctx context.Context, date time.Time) (*expense.OperatingExpense, error) {
	return m.FindForDate(ctx, date)
}

// Test POST /api/operating-expenses - successful creation
func TestCreateOperatingExpense_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockRepo := &mockOperatingExpenseRepository{expenses: []*expense.OperatingExpense{}}
	service := services.NewOperatingExpenseService(mockRepo)
	handler := NewOperatingExpenseHandler(service)

	// Create request
	reqBody := expense.OperatingExpenseRequest{
		PeriodStart:    "2024-01-01",
		PeriodEnd:      "2024-01-31",
		StaffSalary:    2000000,
		Rent:           1000000,
		Utilities:      500000,
		MarketingCosts: 300000,
		OtherExpenses:  200000,
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/operating-expenses", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateOperatingExpense(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response expense.OperatingExpense
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 2000000.0, response.StaffSalary)
	assert.Equal(t, 1000000.0, response.Rent)
	assert.Equal(t, 500000.0, response.Utilities)
	assert.Equal(t, 300000.0, response.MarketingCosts)
	assert.Equal(t, 200000.0, response.OtherExpenses)
	assert.Equal(t, 4000000.0, response.TotalExpenses) // Auto-calculated
}

// Test POST /api/operating-expenses - validation error: invalid date format
func TestCreateOperatingExpense_InvalidDateFormat(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockRepo := &mockOperatingExpenseRepository{expenses: []*expense.OperatingExpense{}}
	service := services.NewOperatingExpenseService(mockRepo)
	handler := NewOperatingExpenseHandler(service)

	// Create request with invalid date format
	reqBody := expense.OperatingExpenseRequest{
		PeriodStart:    "01-01-2024", // Invalid format
		PeriodEnd:      "2024-01-31",
		StaffSalary:    2000000,
		Rent:           1000000,
		Utilities:      500000,
		MarketingCosts: 300000,
		OtherExpenses:  200000,
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/operating-expenses", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateOperatingExpense(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["details"], "invalid period_start date format")
}

// Test POST /api/operating-expenses - validation error: period_start > period_end
func TestCreateOperatingExpense_InvalidDateRange(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockRepo := &mockOperatingExpenseRepository{expenses: []*expense.OperatingExpense{}}
	service := services.NewOperatingExpenseService(mockRepo)
	handler := NewOperatingExpenseHandler(service)

	// Create request with invalid date range
	reqBody := expense.OperatingExpenseRequest{
		PeriodStart:    "2024-01-31",
		PeriodEnd:      "2024-01-01", // End before start
		StaffSalary:    2000000,
		Rent:           1000000,
		Utilities:      500000,
		MarketingCosts: 300000,
		OtherExpenses:  200000,
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/operating-expenses", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateOperatingExpense(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["details"], "period_start must be before or equal to period_end")
}

// Test POST /api/operating-expenses - validation error: negative amounts
func TestCreateOperatingExpense_NegativeAmounts(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockRepo := &mockOperatingExpenseRepository{expenses: []*expense.OperatingExpense{}}
	service := services.NewOperatingExpenseService(mockRepo)
	handler := NewOperatingExpenseHandler(service)

	// Create request with negative amount
	reqBody := map[string]interface{}{
		"period_start":    "2024-01-01",
		"period_end":      "2024-01-31",
		"staff_salary":    -1000000, // Negative
		"rent":            1000000,
		"utilities":       500000,
		"marketing_costs": 300000,
		"other_expenses":  200000,
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/operating-expenses", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateOperatingExpense(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test POST /api/operating-expenses - upsert behavior (update existing)
func TestCreateOperatingExpense_UpsertUpdate(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockRepo := &mockOperatingExpenseRepository{expenses: []*expense.OperatingExpense{}}
	service := services.NewOperatingExpenseService(mockRepo)
	handler := NewOperatingExpenseHandler(service)

	// Create initial expense
	reqBody1 := expense.OperatingExpenseRequest{
		PeriodStart:    "2024-01-01",
		PeriodEnd:      "2024-01-31",
		StaffSalary:    2000000,
		Rent:           1000000,
		Utilities:      500000,
		MarketingCosts: 300000,
		OtherExpenses:  200000,
	}
	body1, _ := json.Marshal(reqBody1)

	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request = httptest.NewRequest("POST", "/api/operating-expenses", bytes.NewBuffer(body1))
	c1.Request.Header.Set("Content-Type", "application/json")
	handler.CreateOperatingExpense(c1)

	// Update with same period but different amounts
	reqBody2 := expense.OperatingExpenseRequest{
		PeriodStart:    "2024-01-01",
		PeriodEnd:      "2024-01-31",
		StaffSalary:    2500000, // Updated
		Rent:           1200000, // Updated
		Utilities:      500000,
		MarketingCosts: 300000,
		OtherExpenses:  200000,
	}
	body2, _ := json.Marshal(reqBody2)

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest("POST", "/api/operating-expenses", bytes.NewBuffer(body2))
	c2.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateOperatingExpense(c2)

	// Assert
	assert.Equal(t, http.StatusOK, w2.Code)

	var response expense.OperatingExpense
	json.Unmarshal(w2.Body.Bytes(), &response)

	assert.Equal(t, 2500000.0, response.StaffSalary)
	assert.Equal(t, 1200000.0, response.Rent)
	assert.Equal(t, 4700000.0, response.TotalExpenses) // Updated total

	// Verify only one expense exists in repository
	assert.Equal(t, 1, len(mockRepo.expenses))
}

// Test GET /api/operating-expenses - no filters
func TestGetOperatingExpenses_NoFilters(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	// Create mock data
	periodStart1, _ := time.Parse("2006-01-02", "2024-01-01")
	periodEnd1, _ := time.Parse("2006-01-02", "2024-01-31")
	periodStart2, _ := time.Parse("2006-01-02", "2024-02-01")
	periodEnd2, _ := time.Parse("2006-01-02", "2024-02-29")
	
	mockRepo := &mockOperatingExpenseRepository{
		expenses: []*expense.OperatingExpense{
			{
				ID:             primitive.NewObjectID(),
				PeriodStart:    periodStart1,
				PeriodEnd:      periodEnd1,
				StaffSalary:    2000000,
				Rent:           1000000,
				Utilities:      500000,
				MarketingCosts: 300000,
				OtherExpenses:  200000,
				TotalExpenses:  4000000,
			},
			{
				ID:             primitive.NewObjectID(),
				PeriodStart:    periodStart2,
				PeriodEnd:      periodEnd2,
				StaffSalary:    2200000,
				Rent:           1000000,
				Utilities:      550000,
				MarketingCosts: 350000,
				OtherExpenses:  250000,
				TotalExpenses:  4350000,
			},
		},
	}
	
	service := services.NewOperatingExpenseService(mockRepo)
	handler := NewOperatingExpenseHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/operating-expenses", nil)

	// Execute
	handler.GetOperatingExpenses(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	expenses := response["expenses"].([]interface{})
	assert.Equal(t, 2, len(expenses))
}

// Test GET /api/operating-expenses - with date range filter
func TestGetOperatingExpenses_WithDateRange(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	// Create mock data
	periodStart1, _ := time.Parse("2006-01-02", "2024-01-01")
	periodEnd1, _ := time.Parse("2006-01-02", "2024-01-31")
	periodStart2, _ := time.Parse("2006-01-02", "2024-02-01")
	periodEnd2, _ := time.Parse("2006-01-02", "2024-02-29")
	periodStart3, _ := time.Parse("2006-01-02", "2024-03-01")
	periodEnd3, _ := time.Parse("2006-01-02", "2024-03-31")
	
	mockRepo := &mockOperatingExpenseRepository{
		expenses: []*expense.OperatingExpense{
			{
				ID:             primitive.NewObjectID(),
				PeriodStart:    periodStart1,
				PeriodEnd:      periodEnd1,
				StaffSalary:    2000000,
				TotalExpenses:  4000000,
			},
			{
				ID:             primitive.NewObjectID(),
				PeriodStart:    periodStart2,
				PeriodEnd:      periodEnd2,
				StaffSalary:    2200000,
				TotalExpenses:  4350000,
			},
			{
				ID:             primitive.NewObjectID(),
				PeriodStart:    periodStart3,
				PeriodEnd:      periodEnd3,
				StaffSalary:    2100000,
				TotalExpenses:  4200000,
			},
		},
	}
	
	service := services.NewOperatingExpenseService(mockRepo)
	handler := NewOperatingExpenseHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/operating-expenses?start_date=2024-01-15&end_date=2024-02-15", nil)

	// Execute
	handler.GetOperatingExpenses(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	expenses := response["expenses"].([]interface{})
	// Should return expenses that overlap with the date range (Jan and Feb)
	assert.Equal(t, 2, len(expenses))
}

// Test GET /api/operating-expenses - invalid date format
func TestGetOperatingExpenses_InvalidDateFormat(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockRepo := &mockOperatingExpenseRepository{expenses: []*expense.OperatingExpense{}}
	service := services.NewOperatingExpenseService(mockRepo)
	handler := NewOperatingExpenseHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/operating-expenses?start_date=01-01-2024", nil)

	// Execute
	handler.GetOperatingExpenses(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Invalid start_date format")
}
