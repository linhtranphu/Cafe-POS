# Design Document: Menu Cost & Profit Analysis

## Overview

Hệ thống Menu Cost & Profit Analysis cung cấp khả năng tính toán và phân tích chi phí (cost) và lợi nhuận (profit) cho từng menu item dựa trên ingredients và cost_per_unit. Hệ thống hỗ trợ hai loại cost:

- **Current Cost**: Giá vốn hiện tại tính từ cost_per_unit hiện tại của ingredients, dùng cho pricing decisions
- **Accounting Cost**: Giá vốn chính thức được lưu khi kết ca (shift closure), dùng cho báo cáo profit/loss

Thiết kế này đảm bảo:
- Tính toán cost chính xác với unit conversion và wastage factor
- Asynchronous processing để tránh blocking UI
- Immutable historical data sau khi shift đóng
- Clear status indicators (FINAL, ESTIMATED, INCOMPLETE)

## Architecture

### High-Level Components

```
┌─────────────────────────────────────────────────────────────┐
│                     Manager Interface                        │
│  (View cost, profit, warnings, history, export reports)     │
└────────────────┬────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────┐
│                    API Layer (Go)                            │
│  - GET /api/menu-items/cost-analysis                        │
│  - GET /api/menu-items/:id/cost-history                     │
│  - GET /api/categories/profit-analysis                      │
└────────────────┬────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────┐
│                  Service Layer                               │
│  ┌──────────────────┐  ┌──────────────────┐                │
│  │ Cost Calculator  │  │ Profit Analyzer  │                │
│  │ Service          │  │ Service          │                │
│  └──────────────────┘  └──────────────────┘                │
│  ┌──────────────────┐  ┌──────────────────┐                │
│  │ Shift Closure    │  │ Background Job   │                │
│  │ Service          │  │ Worker Pool      │                │
│  └──────────────────┘  └──────────────────┘                │
└────────────────┬────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────┐
│                  Data Layer (MongoDB)                        │
│  - menu_items (with current_cost fields)                    │
│  - order_items (with accounting_cost fields)                │
│  - ingredients (with cost_per_unit)                         │
│  - shifts                                                    │
└─────────────────────────────────────────────────────────────┘
```

### Data Flow

**1. Current Cost Calculation (Real-time)**
```
Ingredient cost_per_unit updated
    ↓
Queue background job
    ↓
Worker pool picks up job
    ↓
Calculate current_cost for affected menu items
    ↓
Update menu_items collection
    ↓
Notify UI (refresh or indicator)
```

**2. Accounting Cost Calculation (Shift Closure)**
```
Manager closes shift
    ↓
Fetch all orders in shift
    ↓
For each order_item:
    - Calculate cost from current ingredient cost_per_unit
    - Store accounting_cost in order_items
    - Set cost_status = FINAL
    ↓
Lock shift
    ↓
Generate profit/loss report
```

## Components and Interfaces

### 1. Cost Calculator Service

**Responsibilities:**
- Calculate current_cost for menu items
- Handle unit conversion and wastage
- Manage cost_status (FINAL, ESTIMATED, INCOMPLETE)
- Queue recalculation jobs when ingredients change

**Interface:**
```go
type CostCalculatorService interface {
    // Calculate current cost for a single menu item
    CalculateCurrentCost(menuItemID primitive.ObjectID) (*CostResult, error)
    
    // Calculate current cost for multiple menu items (batch)
    CalculateBatchCurrentCost(menuItemIDs []primitive.ObjectID) ([]CostResult, error)
    
    // Queue recalculation job when ingredient cost changes
    QueueRecalculation(ingredientID primitive.ObjectID) error
    
    // Calculate accounting cost for order items (used during shift closure)
    CalculateAccountingCost(orderItems []OrderItem) ([]OrderItemCost, error)
}

type CostResult struct {
    MenuItemID      primitive.ObjectID
    CurrentCost     float64
    CostStatus      CostStatus
    CalculatedAt    time.Time
    IncompleteItems []string // Names of ingredients with null cost_per_unit
}

type CostStatus string
const (
    CostStatusFinal      CostStatus = "FINAL"      // Shift closed or no ingredients
    CostStatusEstimated  CostStatus = "ESTIMATED"  // Shift not closed
    CostStatusIncomplete CostStatus = "INCOMPLETE" // Missing ingredient costs
)

type OrderItemCost struct {
    OrderItemID     primitive.ObjectID
    AccountingCost  float64
    CostStatus      CostStatus
    CalculatedAt    time.Time
}
```

**Cost Calculation Algorithm:**
```
For each ingredient in menu item:
    1. Get ingredient from database by name
    2. If ingredient.cost_per_unit is null:
        - Mark cost_status = INCOMPLETE
        - Add to incomplete_items list
        - Skip this ingredient
    3. Apply conversion_rate (if defined, else 1.0)
    4. Apply wastage_percentage (if defined, else 0%)
    5. Calculate: ingredient_cost = (quantity * conversion_rate * cost_per_unit) * (1 + wastage_percentage/100)
    6. Add to total_cost

If incomplete_items is not empty:
    - Return cost_status = INCOMPLETE
    - Do not include in profit calculations
Else:
    - Return cost_status = FINAL or ESTIMATED (depending on context)
    - Round to 2 decimal places
```

### 2. Profit Analyzer Service

**Responsibilities:**
- Calculate profit_margin and absolute_profit
- Detect loss and low_margin warnings
- Aggregate profit by category
- Generate profit reports

**Interface:**
```go
type ProfitAnalyzerService interface {
    // Analyze profit for a single menu item
    AnalyzeMenuItem(menuItemID primitive.ObjectID) (*ProfitAnalysis, error)
    
    // Analyze profit for all menu items
    AnalyzeAllMenuItems(filters ProfitFilters) ([]ProfitAnalysis, error)
    
    // Analyze profit by category
    AnalyzeByCategory(dateRange DateRange) ([]CategoryProfit, error)
    
    // Get summary statistics
    GetSummaryStats() (*ProfitSummary, error)
}

type ProfitAnalysis struct {
    MenuItemID      primitive.ObjectID
    Name            string
    Category        string
    Price           float64
    CurrentCost     float64
    ProfitMargin    float64  // Percentage
    AbsoluteProfit  float64  // Cash amount
    WarningStatus   WarningStatus
    CostStatus      CostStatus
}

type WarningStatus string
const (
    WarningNone      WarningStatus = "NONE"       // Profit margin >= 20%
    WarningLowMargin WarningStatus = "LOW_MARGIN" // 0% <= margin < 20%
    WarningLoss      WarningStatus = "LOSS"       // margin < 0%
)

type CategoryProfit struct {
    Category           string
    TotalRevenue       float64
    TotalCost          float64
    TotalProfit        float64
    AverageProfitMargin float64
    OrderCount         int
}

type ProfitSummary struct {
    TotalItems         int
    ItemsWithLoss      int
    ItemsWithLowMargin int
    ItemsIncomplete    int
    AverageProfitMargin float64
}
```

**Profit Calculation Algorithm:**
```
1. Get menu item price and current_cost
2. If cost_status == INCOMPLETE:
    - Return profit_margin = "N/A"
    - Return warning_status = NONE
    - Skip from calculations
3. If price <= 0 (promotional/gifted):
    - Return profit_margin = "N/A"
    - Skip from average calculations
4. Calculate absolute_profit = price - cost
5. Calculate profit_margin = ((price - cost) / price) * 100
6. Round to 2 decimal places
7. Determine warning_status:
    - If profit_margin < 0: LOSS
    - If 0 <= profit_margin < 20: LOW_MARGIN
    - Else: NONE
8. Return result
```

### 3. Shift Closure Service

**Responsibilities:**
- Calculate and store accounting_cost when shift closes
- Lock orders in the shift
- Generate shift profit report

**Interface:**
```go
type ShiftClosureService interface {
    // Close shift and calculate accounting costs
    CloseShift(shiftID primitive.ObjectID) (*ShiftClosureResult, error)
    
    // Get shift profit report
    GetShiftProfitReport(shiftID primitive.ObjectID) (*ShiftProfitReport, error)
}

type ShiftClosureResult struct {
    ShiftID          primitive.ObjectID
    OrdersProcessed  int
    ItemsProcessed   int
    TotalRevenue     float64
    TotalCost        float64
    TotalProfit      float64
    ClosedAt         time.Time
}

type ShiftProfitReport struct {
    ShiftID             primitive.ObjectID
    StartTime           time.Time
    EndTime             time.Time
    TotalRevenue        float64
    TotalCost           float64
    TotalProfit         float64
    ProfitMargin        float64
    OrderCount          int
    CategoryBreakdown   []CategoryProfit
}
```

**Shift Closure Algorithm:**
```
1. Fetch all orders in shift with status != CANCELLED
2. For each order:
    a. For each order_item:
        - Get menu item by menu_item_id
        - Calculate accounting_cost using current ingredient cost_per_unit
        - Store accounting_cost in order_items
        - Store cost_calculated_at = now
        - Set cost_status = FINAL
    b. Update order status to LOCKED
    c. Set locked_at = now
3. Calculate shift totals:
    - total_revenue = sum of order.total
    - total_cost = sum of order_item.accounting_cost * quantity
    - total_profit = total_revenue - total_cost
4. Return ShiftClosureResult
```

### 4. Background Job Worker Pool

**Responsibilities:**
- Process cost recalculation jobs asynchronously
- Manage worker goroutines
- Handle job queue

**Implementation Pattern:**
```go
type JobQueue struct {
    jobs    chan Job
    workers int
    wg      sync.WaitGroup
}

type Job struct {
    Type        JobType
    IngredientID primitive.ObjectID
    MenuItemIDs  []primitive.ObjectID
}

type JobType string
const (
    JobTypeRecalculateCost JobType = "RECALCULATE_COST"
)

func NewJobQueue(workers int) *JobQueue {
    jq := &JobQueue{
        jobs:    make(chan Job, 100), // Buffer size 100
        workers: workers,
    }
    jq.start()
    return jq
}

func (jq *JobQueue) start() {
    for i := 0; i < jq.workers; i++ {
        jq.wg.Add(1)
        go jq.worker()
    }
}

func (jq *JobQueue) worker() {
    defer jq.wg.Done()
    for job := range jq.jobs {
        jq.processJob(job)
    }
}

func (jq *JobQueue) processJob(job Job) {
    switch job.Type {
    case JobTypeRecalculateCost:
        // Find all menu items using this ingredient
        menuItems := findMenuItemsByIngredient(job.IngredientID)
        // Recalculate current_cost for each
        for _, menuItem := range menuItems {
            costResult := calculateCurrentCost(menuItem)
            updateMenuItemCost(menuItem.ID, costResult)
        }
    }
}

func (jq *JobQueue) Enqueue(job Job) {
    jq.jobs <- job
}
```

**Worker Pool Configuration:**
- Number of workers: 5 (configurable)
- Job queue buffer: 100
- Timeout per job: 10 seconds
- Eventual consistency: UI may show stale data briefly

## Data Models

### MenuItem (Extended)

```go
type MenuItem struct {
    ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name                string             `bson:"name" json:"name"`
    Price               float64            `bson:"price" json:"price"`
    Category            string             `bson:"category" json:"category"`
    Description         string             `bson:"description" json:"description"`
    Ingredients         []Ingredient       `bson:"ingredients" json:"ingredients"`
    Available           bool               `bson:"available" json:"available"`
    
    // NEW: Cost tracking fields
    CurrentCost         float64            `bson:"current_cost" json:"current_cost"`
    CostLastCalculatedAt time.Time         `bson:"cost_last_calculated_at" json:"cost_last_calculated_at"`
    CostStatus          CostStatus         `bson:"cost_status" json:"cost_status"`
    IncompleteIngredients []string         `bson:"incomplete_ingredients,omitempty" json:"incomplete_ingredients,omitempty"`
    
    CreatedAt           time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt           time.Time          `bson:"updated_at" json:"updated_at"`
}
```

### OrderItem (Extended)

```go
type OrderItem struct {
    MenuItemID      primitive.ObjectID `bson:"menu_item_id" json:"menu_item_id"`
    Name            string             `bson:"name" json:"name"`
    Price           float64            `bson:"price" json:"price"`
    Quantity        int                `bson:"quantity" json:"quantity"`
    Note            string             `bson:"note,omitempty" json:"note,omitempty"`
    Subtotal        float64            `bson:"subtotal" json:"subtotal"`
    
    // NEW: Accounting cost fields
    AccountingCost  float64            `bson:"accounting_cost,omitempty" json:"accounting_cost,omitempty"`
    CostCalculatedAt *time.Time        `bson:"cost_calculated_at,omitempty" json:"cost_calculated_at,omitempty"`
    CostStatus      CostStatus         `bson:"cost_status,omitempty" json:"cost_status,omitempty"`
}
```

### Ingredient (Extended)

```go
type Ingredient struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name            string             `bson:"name" json:"name"`
    Category        string             `bson:"category" json:"category"`
    Unit            UnitType           `bson:"unit" json:"unit"`
    Quantity        float64            `bson:"quantity" json:"quantity"`
    MinStock        float64            `bson:"min_stock" json:"min_stock"`
    CostPerUnit     float64            `bson:"cost_per_unit" json:"cost_per_unit"`
    Supplier        string             `bson:"supplier" json:"supplier"`
    
    // NEW: Unit conversion and wastage
    ConversionRate  float64            `bson:"conversion_rate,omitempty" json:"conversion_rate,omitempty"` // Default: 1.0
    WastagePercentage float64          `bson:"wastage_percentage,omitempty" json:"wastage_percentage,omitempty"` // Default: 0%
    
    CostUpdatedAt   time.Time          `bson:"cost_updated_at" json:"cost_updated_at"`
    CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}
```

### Database Indexes

```javascript
// menu_items collection
db.menu_items.createIndex({ "cost_status": 1 })
db.menu_items.createIndex({ "category": 1 })
db.menu_items.createIndex({ "current_cost": 1 })

// order_items (embedded in orders)
db.orders.createIndex({ "shift_id": 1, "status": 1 })
db.orders.createIndex({ "items.cost_status": 1 })

// ingredients
db.ingredients.createIndex({ "name": 1 }, { unique: true })
db.ingredients.createIndex({ "cost_updated_at": 1 })
```

