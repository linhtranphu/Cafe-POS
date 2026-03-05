# Phân Tích: Tự Động Tổng Hợp Chi Phí Vận Hành

## 📋 Tổng Quan

Tính năng cho phép tự động tổng hợp chi phí từ phần **Chi phí** (Expenses) sang **Chi phí vận hành** (Operating Expenses) trong Settings theo khoảng thời gian được định nghĩa.

## 🎯 Mục Tiêu

- Khi nhập expense, bắt buộc chọn loại chi phí vận hành (operating_type)
- Trong Settings, có thể tự động tổng hợp expenses theo khoảng thời gian
- Giảm thiểu nhập tay 2 lần (expenses + operating expenses)
- Đảm bảo consistency giữa 2 hệ thống

## 🏗️ Kiến Trúc Hiện Tại

### Expenses (Chi phí thường)
- **Location**: ExpenseManagementView (`frontend/src/views/ExpenseManagementView.vue`)
- **Purpose**: Quản lý chi phí chi tiết hàng ngày
- **Fields**: date, category_id, amount, description, payment_method, vendor, notes
- **Source tracking**: ingredient, facility, maintenance, manual

### Operating Expenses (Chi phí vận hành)
- **Location**: SettingsView → OperatingExpenseForm
- **Purpose**: Báo cáo lợi nhuận (profit analysis)
- **Fields**: period_start, period_end, staff_salary, rent, utilities, marketing_costs, other_expenses
- **Feature**: Phân bổ tự động (AllocateDailyExpense) khi xem báo cáo


## 🔄 Giải Pháp: Thêm Operating Type vào Expense

### 1. Backend Changes

#### A. Domain Model - `backend/domain/expense/expense.go`

**Thêm field mới:**
```go
type Expense struct {
    ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Date          time.Time          `bson:"date" json:"date"`
    CategoryID    primitive.ObjectID `bson:"category_id" json:"category_id"`
    Amount        float64            `bson:"amount" json:"amount"`
    Description   string             `bson:"description" json:"description"`
    PaymentMethod string             `bson:"payment_method" json:"payment_method"`
    Vendor        string             `bson:"vendor,omitempty" json:"vendor,omitempty"`
    Notes         string             `bson:"notes,omitempty" json:"notes,omitempty"`
    
    // NEW: Operating expense type (REQUIRED)
    OperatingType string             `bson:"operating_type" json:"operating_type"`
    
    // Source tracking
    SourceType string             `bson:"source_type,omitempty" json:"source_type,omitempty"`
    SourceID   primitive.ObjectID `bson:"source_id,omitempty" json:"source_id,omitempty"`
    
    CreatedBy string    `bson:"created_by" json:"created_by"`
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// Operating expense type constants
const (
    OperatingTypeStaffSalary    = "staff_salary"
    OperatingTypeRent           = "rent"
    OperatingTypeUtilities      = "utilities"
    OperatingTypeMarketingCosts = "marketing_costs"
    OperatingTypeOther          = "other_expenses"
)
```


#### B. Service Layer - `backend/application/services/operating_expense_service.go`

**Thêm method tổng hợp:**
```go
// AggregateFromExpenses tổng hợp expenses thành operating expense
func (s *OperatingExpenseService) AggregateFromExpenses(
    ctx context.Context, 
    startDate, endDate time.Time,
) (*expense.OperatingExpense, *AggregationBreakdown, error) {
    
    // Lấy tất cả expenses trong khoảng thời gian
    filter := bson.M{
        "date": bson.M{
            "$gte": startDate,
            "$lte": endDate,
        },
    }
    
    expenses, err := s.expenseRepo.GetExpenses(ctx, filter)
    if err != nil {
        return nil, nil, err
    }
    
    // Tổng hợp theo operating_type
    operatingExpense := &expense.OperatingExpense{
        PeriodStart: startDate,
        PeriodEnd:   endDate,
    }
    
    breakdown := &AggregationBreakdown{
        StaffSalary:    &TypeBreakdown{ExpenseIDs: []primitive.ObjectID{}},
        Rent:           &TypeBreakdown{ExpenseIDs: []primitive.ObjectID{}},
        Utilities:      &TypeBreakdown{ExpenseIDs: []primitive.ObjectID{}},
        MarketingCosts: &TypeBreakdown{ExpenseIDs: []primitive.ObjectID{}},
        OtherExpenses:  &TypeBreakdown{ExpenseIDs: []primitive.ObjectID{}},
    }
    
    for _, exp := range expenses {
        switch exp.OperatingType {
        case expense.OperatingTypeStaffSalary:
            operatingExpense.StaffSalary += exp.Amount
            breakdown.StaffSalary.Count++
            breakdown.StaffSalary.ExpenseIDs = append(breakdown.StaffSalary.ExpenseIDs, exp.ID)
        case expense.OperatingTypeRent:
            operatingExpense.Rent += exp.Amount
            breakdown.Rent.Count++
            breakdown.Rent.ExpenseIDs = append(breakdown.Rent.ExpenseIDs, exp.ID)
        case expense.OperatingTypeUtilities:
            operatingExpense.Utilities += exp.Amount
            breakdown.Utilities.Count++
            breakdown.Utilities.ExpenseIDs = append(breakdown.Utilities.ExpenseIDs, exp.ID)
        case expense.OperatingTypeMarketingCosts:
            operatingExpense.MarketingCosts += exp.Amount
            breakdown.MarketingCosts.Count++
            breakdown.MarketingCosts.ExpenseIDs = append(breakdown.MarketingCosts.ExpenseIDs, exp.ID)
        case expense.OperatingTypeOther:
            operatingExpense.OtherExpenses += exp.Amount
            breakdown.OtherExpenses.Count++
            breakdown.OtherExpenses.ExpenseIDs = append(breakdown.OtherExpenses.ExpenseIDs, exp.ID)
        }
    }
    
    operatingExpense.CalculateTotalExpenses()
    return operatingExpense, breakdown, nil
}

type AggregationBreakdown struct {
    StaffSalary    *TypeBreakdown `json:"staff_salary"`
    Rent           *TypeBreakdown `json:"rent"`
    Utilities      *TypeBreakdown `json:"utilities"`
    MarketingCosts *TypeBreakdown `json:"marketing_costs"`
    OtherExpenses  *TypeBreakdown `json:"other_expenses"`
}

type TypeBreakdown struct {
    Count      int                   `json:"count"`
    ExpenseIDs []primitive.ObjectID  `json:"expense_ids"`
}
```


#### C. HTTP Handler - `backend/interfaces/http/operating_expense_handler.go`

**Thêm endpoint mới:**
```go
// AggregateFromExpenses handles GET /api/operating-expenses/aggregate
// Tổng hợp expenses thành operating expense theo khoảng thời gian
func (h *OperatingExpenseHandler) AggregateFromExpenses(c *gin.Context) {
    startDateStr := c.Query("start_date")
    endDateStr := c.Query("end_date")
    
    if startDateStr == "" || endDateStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "start_date and end_date are required",
        })
        return
    }
    
    startDate, err := time.Parse("2006-01-02", startDateStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid start_date format. Expected YYYY-MM-DD",
        })
        return
    }
    
    endDate, err := time.Parse("2006-01-02", endDateStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid end_date format. Expected YYYY-MM-DD",
        })
        return
    }
    
    aggregated, breakdown, err := h.operatingExpenseService.AggregateFromExpenses(
        c.Request.Context(), 
        startDate, 
        endDate,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to aggregate expenses",
            "details": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "operating_expense": aggregated,
        "breakdown": breakdown,
    })
}
```

**Đăng ký route trong `backend/main.go`:**
```go
// Operating expenses
operatingExpenseHandler := http.NewOperatingExpenseHandler(operatingExpenseService)
api.POST("/operating-expenses", operatingExpenseHandler.CreateOperatingExpense)
api.GET("/operating-expenses", operatingExpenseHandler.GetOperatingExpenses)
api.GET("/operating-expenses/aggregate", operatingExpenseHandler.AggregateFromExpenses) // NEW
```

