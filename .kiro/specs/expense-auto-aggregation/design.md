# Design Document: Expense Auto-Aggregation

## Overview

The expense auto-aggregation feature extends the existing Expense Management system to support categorization and automatic aggregation of expenses into Operating Expense summaries. The design introduces an expense_type field to the Expense model, a new aggregation API endpoint, and a UI component in Shop Settings for selecting date ranges, viewing aggregated results with breakdowns, and saving to Operating Expenses.

The system follows a three-layer architecture: frontend (Vue.js), backend API (Go), and data layer (MongoDB). The design emphasizes data integrity, traceability, and user experience with clear feedback and validation.

## Architecture

### System Components

```
┌─────────────────────────────────────────────────────────────┐
│                        Frontend (Vue.js)                     │
├─────────────────────────────────────────────────────────────┤
│  - ExpenseForm (modified with expense_type selector)        │
│  - OperatingExpenseAggregator (new component)               │
│  - AggregationResultsView (new component)                   │
│  - BreakdownDetailModal (new component)                     │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ HTTP/REST
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Backend API (Go)                        │
├─────────────────────────────────────────────────────────────┤
│  - ExpenseHandler (modified)                                │
│  - ExpenseAggregationHandler (new)                          │
│  - ExpenseService (modified)                                │
│  - ShopSettingsService (modified)                           │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ MongoDB Driver
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Data Layer (MongoDB)                    │
├─────────────────────────────────────────────────────────────┤
│  - expenses collection (modified schema)                    │
│  - shop_settings collection (existing)                      │
└─────────────────────────────────────────────────────────────┘
```

### Data Flow

1. **Expense Creation Flow**:
   - User selects expense_type in ExpenseForm
   - Frontend validates expense_type is selected
   - Backend validates expense_type against allowed values
   - Expense saved to MongoDB with expense_type field

2. **Aggregation Flow**:
   - User selects date range in OperatingExpenseAggregator
   - Frontend calls GET /api/expenses/aggregate with date range
   - Backend queries expenses collection filtered by date range
   - Backend groups expenses by expense_type and calculates totals
   - Backend returns aggregated totals with breakdown details
   - Frontend displays results with editable amounts

3. **Save to Operating Expenses Flow**:
   - User reviews/edits aggregated amounts
   - User clicks "Save to Operating Expenses"
   - Frontend calls PUT /api/shop-settings/operating-expenses
   - Backend updates shop_settings document
   - Frontend displays success confirmation

## Components and Interfaces

### Backend Components

#### 1. Modified Expense Model

```go
type Expense struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    ShopID      primitive.ObjectID `bson:"shop_id" json:"shop_id"`
    Amount      float64            `bson:"amount" json:"amount"`
    Description string             `bson:"description" json:"description"`
    Date        time.Time          `bson:"date" json:"date"`
    ExpenseType string             `bson:"expense_type" json:"expense_type"` // NEW FIELD
    CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// Valid expense types
const (
    ExpenseTypeStaffSalary    = "staff_salary"
    ExpenseTypeRent           = "rent"
    ExpenseTypeUtilities      = "utilities"
    ExpenseTypeMarketingCosts = "marketing_costs"
    ExpenseTypeOtherExpenses  = "other_expenses"
)

func (e *Expense) Validate() error {
    validTypes := []string{
        ExpenseTypeStaffSalary,
        ExpenseTypeRent,
        ExpenseTypeUtilities,
        ExpenseTypeMarketingCosts,
        ExpenseTypeOtherExpenses,
    }
    
    if e.ExpenseType == "" {
        return errors.New("expense_type is required")
    }
    
    isValid := false
    for _, vt := range validTypes {
        if e.ExpenseType == vt {
            isValid = true
            break
        }
    }
    
    if !isValid {
        return errors.New("invalid expense_type")
    }
    
    return nil
}
```

#### 2. Aggregation Response Model

```go
type ExpenseAggregation struct {
    ExpenseType string                `json:"expense_type"`
    Total       float64               `json:"total"`
    Count       int                   `json:"count"`
    Breakdown   []ExpenseBreakdownItem `json:"breakdown"`
}

type ExpenseBreakdownItem struct {
    ID          string    `json:"id"`
    Date        time.Time `json:"date"`
    Description string    `json:"description"`
    Amount      float64   `json:"amount"`
}

type AggregationResponse struct {
    StartDate    time.Time             `json:"start_date"`
    EndDate      time.Time             `json:"end_date"`
    Aggregations []ExpenseAggregation  `json:"aggregations"`
    TotalExpenses int                  `json:"total_expenses"`
}
```

#### 3. Aggregation API Handler

```go
// GET /api/expenses/aggregate?start_date=2024-01-01&end_date=2024-01-31
func (h *ExpenseHandler) AggregateExpenses(c *gin.Context) {
    startDateStr := c.Query("start_date")
    endDateStr := c.Query("end_date")
    
    // Parse and validate dates
    startDate, err := time.Parse("2006-01-02", startDateStr)
    if err != nil {
        c.JSON(400, gin.H{"error": "invalid start_date format"})
        return
    }
    
    endDate, err := time.Parse("2006-01-02", endDateStr)
    if err != nil {
        c.JSON(400, gin.H{"error": "invalid end_date format"})
        return
    }
    
    if endDate.Before(startDate) {
        c.JSON(400, gin.H{"error": "end_date must be after start_date"})
        return
    }
    
    // Get shop ID from context
    shopID := c.GetString("shop_id")
    
    // Call service to aggregate
    result, err := h.expenseService.AggregateByDateRange(shopID, startDate, endDate)
    if err != nil {
        c.JSON(500, gin.H{"error": "failed to aggregate expenses"})
        return
    }
    
    c.JSON(200, result)
}
```

#### 4. Expense Service Aggregation Method

```go
func (s *ExpenseService) AggregateByDateRange(shopID string, startDate, endDate time.Time) (*AggregationResponse, error) {
    // Set time to start and end of day for inclusive filtering
    startOfDay := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
    endOfDay := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
    
    // Query expenses in date range with expense_type
    filter := bson.M{
        "shop_id": shopID,
        "date": bson.M{
            "$gte": startOfDay,
            "$lte": endOfDay,
        },
        "expense_type": bson.M{"$exists": true, "$ne": ""},
    }
    
    expenses, err := s.expenseRepo.Find(filter)
    if err != nil {
        return nil, err
    }
    
    // Group by expense_type
    aggregationMap := make(map[string]*ExpenseAggregation)
    expenseTypes := []string{
        ExpenseTypeStaffSalary,
        ExpenseTypeRent,
        ExpenseTypeUtilities,
        ExpenseTypeMarketingCosts,
        ExpenseTypeOtherExpenses,
    }
    
    // Initialize all expense types with zero values
    for _, et := range expenseTypes {
        aggregationMap[et] = &ExpenseAggregation{
            ExpenseType: et,
            Total:       0,
            Count:       0,
            Breakdown:   []ExpenseBreakdownItem{},
        }
    }
    
    // Aggregate expenses
    for _, expense := range expenses {
        agg := aggregationMap[expense.ExpenseType]
        agg.Total += expense.Amount
        agg.Count++
        agg.Breakdown = append(agg.Breakdown, ExpenseBreakdownItem{
            ID:          expense.ID.Hex(),
            Date:        expense.Date,
            Description: expense.Description,
            Amount:      expense.Amount,
        })
    }
    
    // Convert map to slice
    aggregations := make([]ExpenseAggregation, 0, len(expenseTypes))
    for _, et := range expenseTypes {
        aggregations = append(aggregations, *aggregationMap[et])
    }
    
    return &AggregationResponse{
        StartDate:     startDate,
        EndDate:       endDate,
        Aggregations:  aggregations,
        TotalExpenses: len(expenses),
    }, nil
}
```

#### 5. Shop Settings Update Handler

```go
type UpdateOperatingExpensesRequest struct {
    StaffSalary    float64 `json:"staff_salary"`
    Rent           float64 `json:"rent"`
    Utilities      float64 `json:"utilities"`
    MarketingCosts float64 `json:"marketing_costs"`
    OtherExpenses  float64 `json:"other_expenses"`
}

// PUT /api/shop-settings/operating-expenses
func (h *ShopSettingsHandler) UpdateOperatingExpenses(c *gin.Context) {
    var req UpdateOperatingExpensesRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid request body"})
        return
    }
    
    shopID := c.GetString("shop_id")
    
    // Update shop settings
    err := h.settingsService.UpdateOperatingExpenses(shopID, req)
    if err != nil {
        c.JSON(500, gin.H{"error": "failed to update operating expenses"})
        return
    }
    
    c.JSON(200, gin.H{"message": "operating expenses updated successfully"})
}
```

### Frontend Components

#### 1. Modified ExpenseForm Component

```vue
<template>
  <div class="expense-form">
    <!-- Existing fields: amount, description, date -->
    
    <!-- NEW: Expense Type Selector -->
    <div class="form-group">
      <label for="expense-type">Loại chi phí *</label>
      <select 
        id="expense-type" 
        v-model="expense.expense_type" 
        required
        class="form-control"
      >
        <option value="">-- Chọn loại chi phí --</option>
        <option value="staff_salary">Lương nhân viên</option>
        <option value="rent">Thuê mặt bằng</option>
        <option value="utilities">Điện nước</option>
        <option value="marketing_costs">Marketing</option>
        <option value="other_expenses">Khác</option>
      </select>
      <span v-if="errors.expense_type" class="error">{{ errors.expense_type }}</span>
    </div>
    
    <button @click="saveExpense">Lưu</button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      expense: {
        amount: 0,
        description: '',
        date: new Date(),
        expense_type: '' // NEW FIELD
      },
      errors: {}
    }
  },
  methods: {
    validateForm() {
      this.errors = {}
      
      if (!this.expense.expense_type) {
        this.errors.expense_type = 'Vui lòng chọn loại chi phí'
        return false
      }
      
      // Other validations...
      return true
    },
    
    async saveExpense() {
      if (!this.validateForm()) {
        return
      }
      
      try {
        await this.$api.post('/expenses', this.expense)
        this.$emit('saved')
      } catch (error) {
        this.$toast.error('Không thể lưu chi phí')
      }
    }
  }
}
</script>
```

#### 2. OperatingExpenseAggregator Component

```vue
<template>
  <div class="operating-expense-aggregator">
    <h3>Tự động tổng hợp từ Chi phí</h3>
    
    <!-- Date Range Selector -->
    <div class="date-range-selector">
      <div class="form-group">
        <label>Từ ngày</label>
        <input type="date" v-model="startDate" class="form-control" />
      </div>
      
      <div class="form-group">
        <label>Đến ngày</label>
        <input type="date" v-model="endDate" class="form-control" />
      </div>
      
      <button @click="aggregate" :disabled="loading" class="btn-primary">
        {{ loading ? 'Đang tổng hợp...' : 'Tổng hợp' }}
      </button>
    </div>
    
    <!-- Validation Error -->
    <div v-if="validationError" class="alert alert-error">
      {{ validationError }}
    </div>
    
    <!-- Aggregation Results -->
    <AggregationResultsView 
      v-if="aggregationResults"
      :results="aggregationResults"
      :editable-amounts="editableAmounts"
      @save="saveToOperatingExpenses"
      @view-breakdown="showBreakdown"
    />
    
    <!-- Breakdown Modal -->
    <BreakdownDetailModal
      v-if="selectedBreakdown"
      :breakdown="selectedBreakdown"
      @close="selectedBreakdown = null"
    />
  </div>
</template>

<script>
export default {
  data() {
    return {
      startDate: '',
      endDate: '',
      loading: false,
      validationError: '',
      aggregationResults: null,
      editableAmounts: {},
      selectedBreakdown: null
    }
  },
  
  methods: {
    validateDateRange() {
      if (!this.startDate || !this.endDate) {
        this.validationError = 'Vui lòng chọn khoảng thời gian'
        return false
      }
      
      if (new Date(this.endDate) < new Date(this.startDate)) {
        this.validationError = 'Ngày kết thúc phải sau ngày bắt đầu'
        return false
      }
      
      this.validationError = ''
      return true
    },
    
    async aggregate() {
      if (!this.validateDateRange()) {
        return
      }
      
      this.loading = true
      
      try {
        const response = await this.$api.get('/expenses/aggregate', {
          params: {
            start_date: this.startDate,
            end_date: this.endDate
          }
        })
        
        this.aggregationResults = response.data
        
        // Initialize editable amounts with aggregated totals
        this.editableAmounts = {}
        response.data.aggregations.forEach(agg => {
          this.editableAmounts[agg.expense_type] = agg.total
        })
        
        if (response.data.total_expenses === 0) {
          this.$toast.info('Không tìm thấy chi phí trong khoảng thời gian này')
        }
      } catch (error) {
        this.$toast.error('Không thể tổng hợp chi phí')
      } finally {
        this.loading = false
      }
    },
    
    showBreakdown(expenseType) {
      const agg = this.aggregationResults.aggregations.find(
        a => a.expense_type === expenseType
      )
      this.selectedBreakdown = agg
    },
    
    async saveToOperatingExpenses() {
      try {
        await this.$api.put('/shop-settings/operating-expenses', {
          staff_salary: this.editableAmounts.staff_salary || 0,
          rent: this.editableAmounts.rent || 0,
          utilities: this.editableAmounts.utilities || 0,
          marketing_costs: this.editableAmounts.marketing_costs || 0,
          other_expenses: this.editableAmounts.other_expenses || 0
        })
        
        this.$toast.success('Đã cập nhật chi phí hoạt động')
        this.$emit('saved')
      } catch (error) {
        this.$toast.error('Không thể lưu chi phí hoạt động')
      }
    }
  }
}
</script>
```

#### 3. AggregationResultsView Component

```vue
<template>
  <div class="aggregation-results">
    <div class="results-header">
      <h4>Kết quả tổng hợp</h4>
      <p>Từ {{ formatDate(results.start_date) }} đến {{ formatDate(results.end_date) }}</p>
      <p>Tổng số chi phí: {{ results.total_expenses }}</p>
    </div>
    
    <table class="results-table">
      <thead>
        <tr>
          <th>Loại chi phí</th>
          <th>Số lượng</th>
          <th>Tổng tiền</th>
          <th>Chỉnh sửa</th>
          <th>Chi tiết</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="agg in results.aggregations" :key="agg.expense_type">
          <td>{{ getExpenseTypeLabel(agg.expense_type) }}</td>
          <td>{{ agg.count }}</td>
          <td>{{ formatCurrency(agg.total) }}</td>
          <td>
            <input 
              type="number" 
              v-model.number="editableAmounts[agg.expense_type]"
              class="amount-input"
            />
          </td>
          <td>
            <button 
              @click="$emit('view-breakdown', agg.expense_type)"
              :disabled="agg.count === 0"
              class="btn-link"
            >
              Xem chi tiết
            </button>
          </td>
        </tr>
      </tbody>
      <tfoot>
        <tr>
          <td colspan="2"><strong>Tổng cộng</strong></td>
          <td><strong>{{ formatCurrency(totalAmount) }}</strong></td>
          <td><strong>{{ formatCurrency(editedTotalAmount) }}</strong></td>
          <td></td>
        </tr>
      </tfoot>
    </table>
    
    <div class="actions">
      <button @click="$emit('save')" class="btn-primary">
        Lưu vào Chi phí hoạt động
      </button>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    results: {
      type: Object,
      required: true
    },
    editableAmounts: {
      type: Object,
      required: true
    }
  },
  
  computed: {
    totalAmount() {
      return this.results.aggregations.reduce((sum, agg) => sum + agg.total, 0)
    },
    
    editedTotalAmount() {
      return Object.values(this.editableAmounts).reduce((sum, val) => sum + val, 0)
    }
  },
  
  methods: {
    getExpenseTypeLabel(type) {
      const labels = {
        staff_salary: 'Lương nhân viên',
        rent: 'Thuê mặt bằng',
        utilities: 'Điện nước',
        marketing_costs: 'Marketing',
        other_expenses: 'Khác'
      }
      return labels[type] || type
    },
    
    formatCurrency(amount) {
      return new Intl.NumberFormat('vi-VN', {
        style: 'currency',
        currency: 'VND'
      }).format(amount)
    },
    
    formatDate(dateStr) {
      return new Date(dateStr).toLocaleDateString('vi-VN')
    }
  }
}
</script>
```

#### 4. BreakdownDetailModal Component

```vue
<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h3>Chi tiết: {{ getExpenseTypeLabel(breakdown.expense_type) }}</h3>
        <button @click="$emit('close')" class="close-btn">&times;</button>
      </div>
      
      <div class="modal-body">
        <table class="breakdown-table">
          <thead>
            <tr>
              <th>Ngày</th>
              <th>Mô tả</th>
              <th>Số tiền</th>
              <th>Hành động</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in breakdown.breakdown" :key="item.id">
              <td>{{ formatDate(item.date) }}</td>
              <td>{{ item.description }}</td>
              <td>{{ formatCurrency(item.amount) }}</td>
              <td>
                <a :href="`/expenses/${item.id}`" target="_blank" class="btn-link">
                  Xem
                </a>
              </td>
            </tr>
          </tbody>
          <tfoot>
            <tr>
              <td colspan="2"><strong>Tổng cộng</strong></td>
              <td><strong>{{ formatCurrency(totalBreakdown) }}</strong></td>
              <td></td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    breakdown: {
      type: Object,
      required: true
    }
  },
  
  computed: {
    totalBreakdown() {
      return this.breakdown.breakdown.reduce((sum, item) => sum + item.amount, 0)
    }
  },
  
  methods: {
    getExpenseTypeLabel(type) {
      const labels = {
        staff_salary: 'Lương nhân viên',
        rent: 'Thuê mặt bằng',
        utilities: 'Điện nước',
        marketing_costs: 'Marketing',
        other_expenses: 'Khác'
      }
      return labels[type] || type
    },
    
    formatCurrency(amount) {
      return new Intl.NumberFormat('vi-VN', {
        style: 'currency',
        currency: 'VND'
      }).format(amount)
    },
    
    formatDate(dateStr) {
      return new Date(dateStr).toLocaleDateString('vi-VN')
    }
  }
}
</script>
```

## Data Models

### MongoDB Schema Changes

#### Expenses Collection

```javascript
{
  _id: ObjectId,
  shop_id: ObjectId,
  amount: Number,
  description: String,
  date: ISODate,
  expense_type: String,  // NEW: one of the five predefined values
  created_at: ISODate,
  updated_at: ISODate
}

// Index for efficient aggregation queries
db.expenses.createIndex({ shop_id: 1, date: 1, expense_type: 1 })
```

#### Shop Settings Collection (Existing)

```javascript
{
  _id: ObjectId,
  shop_id: ObjectId,
  operating_expenses: {
    staff_salary: Number,
    rent: Number,
    utilities: Number,
    marketing_costs: Number,
    other_expenses: Number
  },
  // ... other settings
}
```

### Data Migration Strategy

For existing expense records without expense_type:

1. **Backward Compatibility**: The system will handle expenses without expense_type gracefully
2. **Aggregation Exclusion**: Expenses without expense_type will be excluded from aggregation
3. **Optional Migration Script**: Provide a script to set expense_type to "other_expenses" for existing records
4. **UI Indication**: Show a badge or indicator for expenses missing expense_type in the expense list

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*


### Property 1: Expense validation rejects missing expense_type

*For any* expense record without an expense_type field, attempting to save it should result in a validation error and the expense should not be persisted to the database.

**Validates: Requirements 1.3**

### Property 2: Expense persistence round-trip

*For any* valid expense with a valid expense_type, saving the expense then retrieving it should return an expense with the same expense_type value.

**Validates: Requirements 1.4, 2.3, 2.4**

### Property 3: Expense type validation rejects invalid values

*For any* expense with an expense_type value not in the set {staff_salary, rent, utilities, marketing_costs, other_expenses}, the system should reject the expense with a validation error.

**Validates: Requirements 2.2**

### Property 4: Aggregation correctness - sum of breakdown equals total

*For any* date range and any set of expenses, the sum of all breakdown item amounts for each expense_type should exactly equal the aggregated total for that expense_type, and the sum of all individual expense amounts in the date range with a given expense_type should equal the aggregated total for that expense_type.

**Validates: Requirements 3.2, 3.7, 5.5**

### Property 5: Aggregation response completeness

*For any* aggregation request, the response should include exactly five expense_type entries (one for each: staff_salary, rent, utilities, marketing_costs, other_expenses), and each entry should include a breakdown array where every item contains id, date, description, and amount fields.

**Validates: Requirements 3.3, 3.4**

### Property 6: Date range validation

*For any* date range where end_date is before start_date, the aggregation API should return an error response and not perform aggregation.

**Validates: Requirements 3.6, 7.1**

### Property 7: Operating expense save with correct field mapping

*For any* set of aggregated or edited expense amounts, saving to operating expenses should correctly map each expense_type to its corresponding shop_settings field (staff_salary → staff_salary, rent → rent, utilities → utilities, marketing_costs → marketing_costs, other_expenses → other_expenses), and retrieving the shop_settings should return the saved values.

**Validates: Requirements 4.10, 6.1, 6.2**

### Property 8: Edited values preserved for saving

*For any* aggregated result where a user edits one or more expense_type amounts, the save operation should use the edited values rather than the originally calculated aggregated values.

**Validates: Requirements 6.6**

### Property 9: Date field consistency in filtering

*For any* aggregation request with a date range, all expenses included in the aggregation should have their date field (or transaction date field) within the specified range [start_date, end_date] inclusive.

**Validates: Requirements 7.4**

### Property 10: Timezone handling consistency

*For any* date range specified in a particular timezone, the aggregation should consistently interpret start_date as the beginning of that day and end_date as the end of that day in the same timezone, ensuring expenses on boundary dates are correctly included or excluded.

**Validates: Requirements 7.5**

### Property 11: State preservation on save failure

*For any* aggregation result with edited amounts, if the save operation to operating expenses fails, the aggregated data and edited amounts should remain available in the UI for retry without requiring re-aggregation.

**Validates: Requirements 8.4**

### Property 12: Concurrent update safety

*For any* two concurrent save operations to operating expenses, the final state should reflect one complete save operation without partial updates or data corruption from interleaved writes.

**Validates: Requirements 8.5**

### Property 13: Total expense count accuracy

*For any* aggregation result, the total_expenses count should equal the sum of all count values across all five expense_type aggregations.

**Validates: Requirements 9.4**

### Property 14: Expenses without expense_type excluded from aggregation

*For any* aggregation request, expenses that do not have an expense_type field or have an empty expense_type value should not appear in any breakdown array and should not contribute to any aggregated total.

**Validates: Requirements 10.3**

## Error Handling

### Validation Errors

1. **Missing expense_type**: Return 400 Bad Request with message "expense_type is required"
2. **Invalid expense_type**: Return 400 Bad Request with message "invalid expense_type, must be one of: staff_salary, rent, utilities, marketing_costs, other_expenses"
3. **Invalid date format**: Return 400 Bad Request with message "invalid date format, expected YYYY-MM-DD"
4. **Invalid date range**: Return 400 Bad Request with message "end_date must be after start_date"

### API Errors

1. **Database connection failure**: Return 500 Internal Server Error with message "database connection failed"
2. **Query timeout**: Return 504 Gateway Timeout with message "aggregation query timed out"
3. **Shop settings not found**: Return 404 Not Found with message "shop settings not found"
4. **Unauthorized access**: Return 401 Unauthorized with message "authentication required"

### Frontend Error Handling

1. **Network errors**: Display toast notification "Không thể kết nối đến server. Vui lòng thử lại."
2. **Validation errors**: Display inline error messages below form fields
3. **API errors**: Display toast notification with translated error message
4. **Retry mechanism**: Provide "Thử lại" button for failed operations

### Edge Cases

1. **Empty date range**: Return aggregation with all zero values
2. **Future date range**: Return aggregation with all zero values (no expenses yet)
3. **Very large date range**: Implement pagination or limit to prevent performance issues
4. **Expenses without expense_type**: Silently exclude from aggregation, log warning
5. **Concurrent saves**: Use optimistic locking or last-write-wins strategy

## Testing Strategy

### Dual Testing Approach

This feature requires both unit tests and property-based tests for comprehensive coverage:

- **Unit tests**: Verify specific examples, edge cases, UI interactions, and error conditions
- **Property tests**: Verify universal properties across all inputs using randomized test data

Both testing approaches are complementary and necessary. Unit tests catch concrete bugs in specific scenarios, while property tests verify general correctness across a wide range of inputs.

### Property-Based Testing

**Library Selection**: Use a property-based testing library appropriate for the language:
- **Go backend**: Use [gopter](https://github.com/leanovate/gopter) or [rapid](https://github.com/flyingmutant/rapid)
- **JavaScript frontend**: Use [fast-check](https://github.com/dubzzz/fast-check)

**Configuration**:
- Each property test MUST run minimum 100 iterations
- Each test MUST include a comment tag referencing the design property
- Tag format: `// Feature: expense-auto-aggregation, Property {number}: {property_text}`

**Property Test Implementation**:

Each correctness property listed above should be implemented as a single property-based test. For example:

```go
// Feature: expense-auto-aggregation, Property 4: Aggregation correctness - sum of breakdown equals total
func TestProperty_AggregationSumEqualsBreakdown(t *testing.T) {
    parameters := gopter.DefaultTestParameters()
    parameters.MinSuccessfulTests = 100
    
    properties := gopter.NewProperties(parameters)
    
    properties.Property("sum of breakdown equals aggregated total", prop.ForAll(
        func(expenses []Expense, startDate, endDate time.Time) bool {
            // Generate random expenses with random expense_types
            // Call aggregation service
            result, _ := service.AggregateByDateRange(shopID, startDate, endDate)
            
            // For each expense_type, verify sum of breakdown equals total
            for _, agg := range result.Aggregations {
                breakdownSum := 0.0
                for _, item := range agg.Breakdown {
                    breakdownSum += item.Amount
                }
                
                if math.Abs(breakdownSum - agg.Total) > 0.01 {
                    return false
                }
            }
            
            return true
        },
        genExpenseSlice(),
        genDate(),
        genDate(),
    ))
    
    properties.TestingRun(t)
}
```

### Unit Testing

**Focus Areas**:
- Specific examples demonstrating correct behavior
- Edge cases: empty date ranges, boundary dates, missing expense_type
- Error conditions: invalid dates, network failures, validation errors
- UI interactions: button clicks, form submissions, modal displays
- Integration points: API calls, database operations

**Unit Test Balance**:
- Avoid writing too many unit tests for scenarios covered by property tests
- Focus unit tests on concrete examples and integration scenarios
- Use unit tests for UI component rendering and user interactions
- Use unit tests for specific error handling paths

**Example Unit Tests**:

```go
func TestExpenseValidation_MissingExpenseType(t *testing.T) {
    expense := Expense{
        Amount: 100000,
        Description: "Test expense",
        Date: time.Now(),
        // expense_type is missing
    }
    
    err := expense.Validate()
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "expense_type is required")
}

func TestAggregation_EmptyDateRange(t *testing.T) {
    startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
    
    // No expenses in database
    result, err := service.AggregateByDateRange(shopID, startDate, endDate)
    
    assert.NoError(t, err)
    assert.Equal(t, 0, result.TotalExpenses)
    assert.Equal(t, 5, len(result.Aggregations))
    
    for _, agg := range result.Aggregations {
        assert.Equal(t, 0.0, agg.Total)
        assert.Equal(t, 0, agg.Count)
    }
}
```

### Integration Testing

1. **End-to-end flow**: Create expense → Aggregate → Save to operating expenses
2. **API integration**: Test all API endpoints with real database
3. **Frontend integration**: Test Vue components with mock API responses
4. **Database integration**: Test MongoDB queries and indexes

### Manual Testing Checklist

1. Create expenses with all five expense_types
2. Verify Vietnamese labels display correctly
3. Test aggregation with various date ranges
4. Verify breakdown details match individual expenses
5. Edit aggregated amounts and verify edited values are saved
6. Test with empty date ranges
7. Test with invalid date ranges
8. Test concurrent saves from multiple users
9. Test with existing expenses without expense_type
10. Verify operating expenses update correctly in shop settings

## Implementation Notes

### Performance Considerations

1. **Index Strategy**: Create compound index on (shop_id, date, expense_type) for efficient aggregation queries
2. **Query Optimization**: Use MongoDB aggregation pipeline for server-side grouping
3. **Caching**: Consider caching aggregation results for frequently accessed date ranges
4. **Pagination**: For very large result sets, implement pagination in breakdown views

### Security Considerations

1. **Authorization**: Verify user has permission to view expenses and modify shop settings
2. **Input Validation**: Sanitize all user inputs to prevent injection attacks
3. **Rate Limiting**: Implement rate limiting on aggregation API to prevent abuse
4. **Audit Logging**: Log all changes to operating expenses for audit trail

### Deployment Strategy

1. **Database Migration**: Add expense_type field to expenses collection schema
2. **Index Creation**: Create indexes before deploying code changes
3. **Backward Compatibility**: Ensure existing code handles expenses without expense_type
4. **Feature Flag**: Consider using feature flag for gradual rollout
5. **Monitoring**: Add metrics for aggregation API performance and error rates

### Future Enhancements

1. **Custom Date Ranges**: Add preset options (This Month, Last Month, This Quarter)
2. **Export Functionality**: Allow exporting aggregation results to CSV/Excel
3. **Scheduled Aggregation**: Automatically aggregate expenses at month-end
4. **Multi-Currency Support**: Handle expenses in different currencies
5. **Expense Categories**: Allow custom expense categories beyond the five predefined types
6. **Approval Workflow**: Add approval step before saving to operating expenses
