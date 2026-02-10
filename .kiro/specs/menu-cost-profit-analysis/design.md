# Design Document: Menu Cost & Profit Analysis

## Overview

Hệ thống Menu Cost & Profit Analysis cung cấp khả năng tính toán chi phí nguyên liệu (cost of goods sold - COGS) và phân tích lợi nhuận cho từng menu item trong hệ thống POS cafe. Tính năng này giúp manager:

- Theo dõi giá vốn hiện tại (current_cost) của mỗi món dựa trên cost_per_unit của ingredients
- Tính toán và lưu trữ giá vốn chính thức (accounting_cost) khi kết ca để phục vụ báo cáo kế toán
- Phân tích profit margin và absolute profit cho từng món
- Phát hiện các món bán lỗ hoặc có lợi nhuận thấp
- Xem báo cáo lợi nhuận theo category và theo thời gian
- Phân tích operating profit sau khi trừ chi phí vận hành

Hệ thống sử dụng kiến trúc phân tầng với:
- **Backend**: Go với MongoDB, xử lý business logic và data persistence
- **Frontend**: Vue.js với responsive design cho cả desktop và mobile
- **Real-time updates**: Asynchronous cost recalculation với eventual consistency

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Frontend (Vue.js)                        │
│  ┌──────────────────┐  ┌──────────────────┐                │
│  │  Manager View    │  │  Cost Analysis   │                │
│  │  - Menu Cost     │  │  - Profit Report │                │
│  │  - Warnings      │  │  - Category View │                │
│  └──────────────────┘  └──────────────────┘                │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ HTTP/REST API
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    Backend API Layer (Go)                    │
│  ┌──────────────────┐  ┌──────────────────┐                │
│  │  Menu Handler    │  │  Report Handler  │                │
│  └──────────────────┘  └──────────────────┘                │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   Service Layer (Go)                         │
│  ┌──────────────────┐  ┌──────────────────┐                │
│  │ Cost Calculator  │  │ Profit Analyzer  │                │
│  │ Service          │  │ Service          │                │
│  └──────────────────┘  └──────────────────┘                │
│  ┌──────────────────┐  ┌──────────────────┐                │
│  │ Cost Recalc      │  │ Operating Expense│                │
│  │ Service          │  │ Service          │                │
│  └──────────────────┘  └──────────────────┘                │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   Data Layer (MongoDB)                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  menu_items  │  │  order_items │  │  ingredients │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│  ┌──────────────┐  ┌──────────────┐                        │
│  │  shifts      │  │  operating_  │                        │
│  │              │  │  expenses    │                        │
│  └──────────────┘  └──────────────┘                        │
└─────────────────────────────────────────────────────────────┘
```

### Key Design Decisions


**1. Dual Cost Model (Current vs Accounting)**

Hệ thống duy trì hai loại cost:
- **current_cost**: Giá vốn real-time, tính từ cost_per_unit hiện tại của ingredients. Được cập nhật asynchronously khi ingredient costs thay đổi. Dùng cho pricing decisions.
- **accounting_cost**: Giá vốn chính thức được lưu khi kết ca (shift closure). Immutable sau khi shift đóng. Dùng cho báo cáo profit/loss và accounting.

Lý do: Phù hợp với quy trình kế toán cafe, cost phản ánh giá vốn tại thời điểm chốt sổ, không bị ảnh hưởng bởi thay đổi giá nguyên liệu sau đó.

**2. Asynchronous Cost Recalculation**

Khi ingredient cost_per_unit thay đổi, hệ thống queue background job để recalculate current_cost cho tất cả menu items bị ảnh hưởng. Eventual consistency được chấp nhận để tránh treo hệ thống.

**3. Cost Status Tracking**

Mỗi cost calculation được đánh dấu với status:
- **FINAL**: Cost đã được tính và lưu chính thức
- **ESTIMATED**: Cost tạm tính từ current_cost (shift chưa đóng)
- **INCOMPLETE**: Thiếu giá nguyên liệu, không thể tính cost chính xác

**4. Shift Closure as Cost Snapshot**

Cost calculation được trigger khi manager kết ca (shift closure). Tại thời điểm này, hệ thống:
- Lấy cost_per_unit hiện tại của tất cả ingredients
- Tính accounting_cost cho tất cả order items trong shift
- Lưu cost với timestamp và status = FINAL
- Accounting_cost trở thành immutable sau khi lưu

**5. Operating Expense Allocation**

Manager nhập operating expenses theo kỳ (daily/monthly). Khi xem báo cáo với granularity khác (e.g., xem daily report nhưng expenses nhập monthly), hệ thống tự động phân bổ proportionally với indicator rõ ràng.

## Components and Interfaces

### Backend Components

#### 1. Cost Calculator Service

**Responsibility**: Tính toán cost cho menu items dựa trên ingredients

**Interface**:
```go
type CostCalculatorService interface {
    // Calculate current cost for a single menu item
    CalculateMenuItemCost(ctx context.Context, menuItemID primitive.ObjectID) (*MenuItemCost, error)
    
    // Calculate current cost for all menu items
    CalculateAllMenuItemCosts(ctx context.Context) ([]MenuItemCost, error)
    
    // Calculate accounting cost for orders in a shift (called during shift closure)
    CalculateShiftOrderCosts(ctx context.Context, shiftID primitive.ObjectID) error
    
    // Queue background job to recalculate costs for menu items using specific ingredient
    QueueCostRecalculation(ctx context.Context, ingredientID primitive.ObjectID) error
}
```

**Key Methods**:
- `CalculateMenuItemCost`: Tính current_cost cho một menu item
  - Lấy danh sách ingredients từ menu item
  - Với mỗi ingredient, lấy cost_per_unit hiện tại
  - Áp dụng conversion_rate và wastage_percentage nếu có
  - Tính tổng: `sum(quantity * conversion_rate * cost_per_unit * (1 + wastage_percentage/100))`
  - Xác định cost_status (FINAL, INCOMPLETE)
  - Round về 2 decimal places


- `CalculateShiftOrderCosts`: Tính accounting_cost cho tất cả orders trong shift khi kết ca
  - Lấy tất cả orders trong shift
  - Với mỗi order item, tính cost dựa trên menu item recipe và cost_per_unit tại thời điểm hiện tại
  - Lưu accounting_cost vào order_items collection
  - Đánh dấu cost_status = FINAL
  - Lưu timestamp cost_calculated_at

- `QueueCostRecalculation`: Queue background job khi ingredient cost thay đổi
  - Tìm tất cả menu items sử dụng ingredient này
  - Queue job để recalculate current_cost cho từng menu item
  - Job chạy asynchronously với eventual consistency

#### 2. Profit Analyzer Service

**Responsibility**: Phân tích profit margin và absolute profit

**Interface**:
```go
type ProfitAnalyzerService interface {
    // Calculate profit metrics for a menu item
    CalculateMenuItemProfit(ctx context.Context, menuItemID primitive.ObjectID) (*MenuItemProfit, error)
    
    // Get profit analysis for all menu items
    GetAllMenuItemProfits(ctx context.Context, filter ProfitFilter) ([]MenuItemProfit, error)
    
    // Get category-level profit analysis
    GetCategoryProfits(ctx context.Context, dateRange DateRange) ([]CategoryProfit, error)
    
    // Get operating profit analysis
    GetOperatingProfit(ctx context.Context, dateRange DateRange) (*OperatingProfitReport, error)
    
    // Detect loss and low margin items
    DetectWarnings(ctx context.Context, lowMarginThreshold float64) (*ProfitWarnings, error)
}
```

**Key Methods**:
- `CalculateMenuItemProfit`: Tính profit metrics cho một menu item
  - Lấy current_cost và price
  - Tính profit_margin = ((price - cost) / price) * 100
  - Tính absolute_profit = price - cost
  - Xử lý edge cases (price = 0, cost > price)
  - Round về 2 decimal places

- `GetCategoryProfits`: Tính profit theo category
  - Lấy tất cả orders trong date range
  - Group by category
  - Tính total_revenue, total_cost (dùng accounting_cost), total_profit
  - Tính average_profit_margin = (total_profit / total_revenue) * 100

- `GetOperatingProfit`: Tính operating profit sau khi trừ chi phí vận hành
  - Tính gross_profit từ orders (revenue - COGS)
  - Lấy operating expenses trong date range
  - Phân bổ expenses nếu cần (monthly → daily)
  - Tính operating_profit = gross_profit - total_expenses
  - Tính operating_profit_margin = (operating_profit / revenue) * 100

#### 3. Cost Recalculation Service

**Responsibility**: Xử lý background jobs cho cost recalculation

**Interface**:
```go
type CostRecalculationService interface {
    // Process queued cost recalculation jobs
    ProcessRecalculationQueue(ctx context.Context) error
    
    // Check recalculation status
    GetRecalculationStatus(ctx context.Context) (*RecalculationStatus, error)
}
```

**Implementation Strategy**:
- Sử dụng Go channels hoặc message queue (e.g., Redis) để queue jobs
- Worker pool xử lý jobs concurrently
- Batch processing để tối ưu database queries
- Timeout 5 seconds cho up to 1000 menu items


#### 4. Operating Expense Service

**Responsibility**: Quản lý operating expenses và phân bổ chi phí

**Interface**:
```go
type OperatingExpenseService interface {
    // Create or update operating expense for a period
    UpsertOperatingExpense(ctx context.Context, req *OperatingExpenseRequest) (*OperatingExpense, error)
    
    // Get operating expense for a specific date
    GetOperatingExpenseForDate(ctx context.Context, date time.Time) (*OperatingExpense, error)
    
    // Get operating expenses for a date range
    GetOperatingExpenses(ctx context.Context, dateRange DateRange) ([]OperatingExpense, error)
    
    // Allocate monthly expense to daily
    AllocateDailyExpense(ctx context.Context, monthlyExpense *OperatingExpense, targetDate time.Time) (*AllocatedExpense, error)
}
```

### Frontend Components

#### 1. MenuCostView Component

**Responsibility**: Hiển thị bảng cost và profit cho tất cả menu items

**Props**: None (fetches data internally)

**State**:
```typescript
interface MenuCostViewState {
  menuItems: MenuItemCost[]
  loading: boolean
  error: string | null
  filterCategory: string | null
  sortBy: 'profit_margin' | 'absolute_profit' | 'name'
  sortOrder: 'asc' | 'desc'
  lowMarginThreshold: number
  recalculationStatus: RecalculationStatus | null
}
```

**Key Features**:
- Table hiển thị: name, current_cost, price, profit_margin, absolute_profit, cost_status
- Color coding: green (profitable), yellow (low margin), red (loss), gray (incomplete)
- Filter by category
- Sort by profit_margin, absolute_profit, name
- Summary statistics: total items, loss count, low margin count, average profit_margin
- Recalculation indicator khi costs đang được update
- Click vào row để xem ingredient cost breakdown

#### 2. ProfitAnalysisView Component

**Responsibility**: Hiển thị báo cáo profit theo category và operating profit

**Props**: None

**State**:
```typescript
interface ProfitAnalysisViewState {
  dateRange: DateRange
  categoryProfits: CategoryProfit[]
  operatingProfit: OperatingProfitReport | null
  loading: boolean
  error: string | null
  viewMode: 'category' | 'operating'
}
```

**Key Features**:
- Date range picker (daily, weekly, monthly)
- Category profit table: category, revenue, cost, profit, margin
- Operating profit breakdown: gross profit, expenses by type, operating profit
- Charts: profit trend over time, category comparison
- Export to CSV


#### 3. OperatingExpenseForm Component

**Responsibility**: Form nhập operating expenses

**Props**:
```typescript
interface OperatingExpenseFormProps {
  initialData?: OperatingExpense
  onSave: (expense: OperatingExpense) => void
  onCancel: () => void
}
```

**State**:
```typescript
interface OperatingExpenseFormState {
  periodStart: Date
  periodEnd: Date
  staffSalary: number
  rent: number
  utilities: number
  marketingCosts: number
  otherExpenses: number
  saving: boolean
  errors: Record<string, string>
}
```

**Key Features**:
- Date range picker cho period
- Input fields cho từng loại expense
- Auto-calculate total expenses
- Validation: period_start <= period_end, all amounts >= 0
- Save/Cancel buttons

#### 4. MenuItemCostBreakdown Component

**Responsibility**: Hiển thị chi tiết cost breakdown cho một menu item

**Props**:
```typescript
interface MenuItemCostBreakdownProps {
  menuItemId: string
  onClose: () => void
}
```

**State**:
```typescript
interface MenuItemCostBreakdownState {
  menuItem: MenuItem
  ingredients: IngredientCostDetail[]
  totalCost: number
  loading: boolean
}

interface IngredientCostDetail {
  name: string
  quantity: number
  unit: string
  costPerUnit: number
  conversionRate: number
  wastagePercentage: number
  totalCost: number
}
```

**Key Features**:
- Modal/drawer hiển thị ingredient list
- Table: ingredient name, quantity, unit, cost per unit, total cost
- Hiển thị conversion rate và wastage percentage nếu có
- Total cost summary
- Warning nếu ingredient thiếu cost_per_unit

## Data Models

### Backend Data Models

#### MenuItem (Extended)

```go
type MenuItem struct {
    ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name                 string             `bson:"name" json:"name"`
    Price                float64            `bson:"price" json:"price"`
    Category             string             `bson:"category" json:"category"`
    Description          string             `bson:"description" json:"description"`
    Ingredients          []Ingredient       `bson:"ingredients" json:"ingredients"`
    Available            bool               `bson:"available" json:"available"`
    
    // NEW: Cost tracking fields
    CurrentCost          float64            `bson:"current_cost" json:"current_cost"`
    CostLastCalculatedAt time.Time          `bson:"cost_last_calculated_at" json:"cost_last_calculated_at"`
    CostStatus           CostStatus         `bson:"cost_status" json:"cost_status"`
    
    CreatedAt            time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt            time.Time          `bson:"updated_at" json:"updated_at"`
}

type CostStatus string

const (
    CostStatusFinal      CostStatus = "FINAL"
    CostStatusEstimated  CostStatus = "ESTIMATED"
    CostStatusIncomplete CostStatus = "INCOMPLETE"
)
```


#### OrderItem (Extended)

```go
type OrderItem struct {
    MenuItemID         primitive.ObjectID `bson:"menu_item_id" json:"menu_item_id"`
    Name               string             `bson:"name" json:"name"`
    Price              float64            `bson:"price" json:"price"`
    Quantity           int                `bson:"quantity" json:"quantity"`
    Note               string             `bson:"note,omitempty" json:"note,omitempty"`
    Subtotal           float64            `bson:"subtotal" json:"subtotal"`
    
    // NEW: Accounting cost fields
    AccountingCost     float64            `bson:"accounting_cost" json:"accounting_cost"`
    CostCalculatedAt   time.Time          `bson:"cost_calculated_at" json:"cost_calculated_at"`
    CostStatus         CostStatus         `bson:"cost_status" json:"cost_status"`
}
```

**Note**: OrderItem không có collection riêng, nó là embedded document trong Order. Để lưu accounting_cost, ta có 2 options:
- **Option A**: Lưu trong Order.Items array (embedded) - Đơn giản nhưng khó query
- **Option B**: Tạo collection order_items riêng với order_id reference - Dễ query và aggregate

**Decision**: Sử dụng Option B (separate collection) để dễ dàng query và aggregate cho profit reports.

#### Ingredient (Extended)

```go
type Ingredient struct {
    ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name             string             `bson:"name" json:"name"`
    Category         string             `bson:"category" json:"category"`
    Unit             UnitType           `bson:"unit" json:"unit"`
    Quantity         float64            `bson:"quantity" json:"quantity"`
    MinStock         float64            `bson:"min_stock" json:"min_stock"`
    CostPerUnit      float64            `bson:"cost_per_unit" json:"cost_per_unit"`
    Supplier         string             `bson:"supplier" json:"supplier"`
    
    // NEW: Unit conversion and wastage
    ConversionRate   float64            `bson:"conversion_rate" json:"conversion_rate"`     // Default: 1.0
    WastagePercentage float64           `bson:"wastage_percentage" json:"wastage_percentage"` // Default: 0.0
    
    CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}
```

#### OperatingExpense (NEW)

```go
type OperatingExpense struct {
    ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    PeriodStart    time.Time          `bson:"period_start" json:"period_start"`
    PeriodEnd      time.Time          `bson:"period_end" json:"period_end"`
    StaffSalary    float64            `bson:"staff_salary" json:"staff_salary"`
    Rent           float64            `bson:"rent" json:"rent"`
    Utilities      float64            `bson:"utilities" json:"utilities"`
    MarketingCosts float64            `bson:"marketing_costs" json:"marketing_costs"`
    OtherExpenses  float64            `bson:"other_expenses" json:"other_expenses"`
    TotalExpenses  float64            `bson:"total_expenses" json:"total_expenses"` // Auto-calculated
    CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

type OperatingExpenseRequest struct {
    PeriodStart    string  `json:"period_start" binding:"required"` // ISO 8601 date
    PeriodEnd      string  `json:"period_end" binding:"required"`
    StaffSalary    float64 `json:"staff_salary" binding:"min=0"`
    Rent           float64 `json:"rent" binding:"min=0"`
    Utilities      float64 `json:"utilities" binding:"min=0"`
    MarketingCosts float64 `json:"marketing_costs" binding:"min=0"`
    OtherExpenses  float64 `json:"other_expenses" binding:"min=0"`
}
```


#### ShopSettings (Extended)

```go
type ShopSettings struct {
    ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    ShopName            string             `bson:"shop_name" json:"shop_name"`
    
    // NEW: Profit analysis settings
    LowMarginThreshold  float64            `bson:"low_margin_threshold" json:"low_margin_threshold"` // Default: 20.0
    
    CreatedAt           time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt           time.Time          `bson:"updated_at" json:"updated_at"`
}
```

#### Response DTOs

```go
// MenuItemCost - Response for menu item cost calculation
type MenuItemCost struct {
    MenuItemID           primitive.ObjectID `json:"menu_item_id"`
    Name                 string             `json:"name"`
    Category             string             `json:"category"`
    Price                float64            `json:"price"`
    CurrentCost          float64            `json:"current_cost"`
    ProfitMargin         float64            `json:"profit_margin"`
    AbsoluteProfit       float64            `json:"absolute_profit"`
    CostStatus           CostStatus         `json:"cost_status"`
    CostLastCalculatedAt time.Time          `json:"cost_last_calculated_at"`
    WarningStatus        WarningStatus      `json:"warning_status"` // "none", "low_margin", "loss"
}

type WarningStatus string

const (
    WarningNone      WarningStatus = "none"
    WarningLowMargin WarningStatus = "low_margin"
    WarningLoss      WarningStatus = "loss"
)

// CategoryProfit - Response for category-level profit analysis
type CategoryProfit struct {
    Category            string    `json:"category"`
    TotalRevenue        float64   `json:"total_revenue"`
    TotalCost           float64   `json:"total_cost"`
    TotalProfit         float64   `json:"total_profit"`
    AverageProfitMargin float64   `json:"average_profit_margin"`
    OrderCount          int       `json:"order_count"`
    ItemCount           int       `json:"item_count"`
}

// OperatingProfitReport - Response for operating profit analysis
type OperatingProfitReport struct {
    DateRange           DateRange          `json:"date_range"`
    TotalRevenue        float64            `json:"total_revenue"`
    TotalCOGS           float64            `json:"total_cogs"`
    GrossProfit         float64            `json:"gross_profit"`
    GrossProfitMargin   float64            `json:"gross_profit_margin"`
    
    // Operating expenses breakdown
    StaffSalary         float64            `json:"staff_salary"`
    Rent                float64            `json:"rent"`
    Utilities           float64            `json:"utilities"`
    MarketingCosts      float64            `json:"marketing_costs"`
    OtherExpenses       float64            `json:"other_expenses"`
    TotalExpenses       float64            `json:"total_expenses"`
    
    OperatingProfit     float64            `json:"operating_profit"`
    OperatingProfitMargin float64          `json:"operating_profit_margin"`
    
    ExpenseAllocated    bool               `json:"expense_allocated"` // true if expenses were allocated from monthly to daily
    AllocationNote      string             `json:"allocation_note,omitempty"`
}

// ProfitWarnings - Response for loss and low margin detection
type ProfitWarnings struct {
    LossItems       []MenuItemCost `json:"loss_items"`
    LowMarginItems  []MenuItemCost `json:"low_margin_items"`
    LossCount       int            `json:"loss_count"`
    LowMarginCount  int            `json:"low_margin_count"`
    Threshold       float64        `json:"threshold"`
}

// RecalculationStatus - Status of background cost recalculation
type RecalculationStatus struct {
    InProgress      bool      `json:"in_progress"`
    QueuedItems     int       `json:"queued_items"`
    ProcessedItems  int       `json:"processed_items"`
    LastUpdated     time.Time `json:"last_updated"`
}
```


### Frontend Data Models

```typescript
// MenuItem with cost information
interface MenuItemCost {
  menu_item_id: string
  name: string
  category: string
  price: number
  current_cost: number
  profit_margin: number
  absolute_profit: number
  cost_status: 'FINAL' | 'ESTIMATED' | 'INCOMPLETE'
  cost_last_calculated_at: string
  warning_status: 'none' | 'low_margin' | 'loss'
}

// Category profit analysis
interface CategoryProfit {
  category: string
  total_revenue: number
  total_cost: number
  total_profit: number
  average_profit_margin: number
  order_count: number
  item_count: number
}

// Operating profit report
interface OperatingProfitReport {
  date_range: DateRange
  total_revenue: number
  total_cogs: number
  gross_profit: number
  gross_profit_margin: number
  staff_salary: number
  rent: number
  utilities: number
  marketing_costs: number
  other_expenses: number
  total_expenses: number
  operating_profit: number
  operating_profit_margin: number
  expense_allocated: boolean
  allocation_note?: string
}

// Operating expense form data
interface OperatingExpense {
  id?: string
  period_start: string
  period_end: string
  staff_salary: number
  rent: number
  utilities: number
  marketing_costs: number
  other_expenses: number
  total_expenses: number
}

// Date range filter
interface DateRange {
  start: string // ISO 8601 date
  end: string   // ISO 8601 date
}

// Profit filter
interface ProfitFilter {
  category?: string
  sort_by?: 'profit_margin' | 'absolute_profit' | 'name'
  sort_order?: 'asc' | 'desc'
}
```

## Database Schema Changes

### Collection: menu_items

**New Fields**:
```javascript
{
  // ... existing fields ...
  current_cost: Number,              // Giá vốn hiện tại
  cost_last_calculated_at: Date,     // Timestamp tính cost
  cost_status: String                // "FINAL" | "ESTIMATED" | "INCOMPLETE"
}
```

**Indexes**:
```javascript
db.menu_items.createIndex({ category: 1 })
db.menu_items.createIndex({ cost_status: 1 })
db.menu_items.createIndex({ current_cost: 1 })
```

### Collection: order_items (NEW)

**Purpose**: Separate collection để lưu order items với accounting cost, dễ query và aggregate

**Schema**:
```javascript
{
  _id: ObjectId,
  order_id: ObjectId,              // Reference to orders collection
  menu_item_id: ObjectId,          // Reference to menu_items collection
  name: String,
  price: Number,
  quantity: Number,
  note: String,
  subtotal: Number,
  
  // Cost tracking
  accounting_cost: Number,         // Giá vốn chính thức (per item)
  cost_calculated_at: Date,        // Timestamp tính cost
  cost_status: String,             // "FINAL" | "ESTIMATED" | "INCOMPLETE"
  
  created_at: Date
}
```

**Indexes**:
```javascript
db.order_items.createIndex({ order_id: 1 })
db.order_items.createIndex({ menu_item_id: 1 })
db.order_items.createIndex({ cost_status: 1 })
db.order_items.createIndex({ cost_calculated_at: 1 })
```


### Collection: ingredients

**New Fields**:
```javascript
{
  // ... existing fields ...
  conversion_rate: Number,           // Tỷ lệ quy đổi đơn vị (default: 1.0)
  wastage_percentage: Number         // Tỷ lệ hao hụt % (default: 0.0)
}
```

### Collection: operating_expenses (NEW)

**Schema**:
```javascript
{
  _id: ObjectId,
  period_start: Date,                // Ngày bắt đầu kỳ
  period_end: Date,                  // Ngày kết thúc kỳ
  staff_salary: Number,              // Lương nhân viên
  rent: Number,                      // Tiền thuê mặt bằng
  utilities: Number,                 // Điện nước
  marketing_costs: Number,           // Chi phí marketing
  other_expenses: Number,            // Chi phí khác
  total_expenses: Number,            // Tổng chi phí (auto-calculated)
  created_at: Date,
  updated_at: Date
}
```

**Indexes**:
```javascript
db.operating_expenses.createIndex({ period_start: 1, period_end: 1 })
db.operating_expenses.createIndex({ period_start: 1 })
```

### Collection: shop_settings

**New Fields**:
```javascript
{
  // ... existing fields ...
  low_margin_threshold: Number       // Ngưỡng cảnh báo lợi nhuận thấp (default: 20.0)
}
```

## API Endpoints

### Menu Cost Endpoints

#### GET /api/menu/costs

**Description**: Lấy danh sách cost và profit cho tất cả menu items

**Query Parameters**:
- `category` (optional): Filter by category
- `sort_by` (optional): "profit_margin" | "absolute_profit" | "name"
- `sort_order` (optional): "asc" | "desc"

**Response**:
```json
{
  "items": [
    {
      "menu_item_id": "...",
      "name": "Cappuccino",
      "category": "Coffee",
      "price": 45000,
      "current_cost": 15000,
      "profit_margin": 66.67,
      "absolute_profit": 30000,
      "cost_status": "FINAL",
      "cost_last_calculated_at": "2024-01-15T10:30:00Z",
      "warning_status": "none"
    }
  ],
  "summary": {
    "total_items": 50,
    "loss_count": 2,
    "low_margin_count": 5,
    "average_profit_margin": 55.5
  },
  "recalculation_status": {
    "in_progress": false,
    "queued_items": 0,
    "processed_items": 50,
    "last_updated": "2024-01-15T10:30:00Z"
  }
}
```

#### GET /api/menu/costs/:id

**Description**: Lấy cost detail và ingredient breakdown cho một menu item

**Response**:
```json
{
  "menu_item": {
    "id": "...",
    "name": "Cappuccino",
    "price": 45000,
    "current_cost": 15000,
    "cost_status": "FINAL"
  },
  "ingredients": [
    {
      "name": "Espresso",
      "quantity": 30,
      "unit": "ml",
      "cost_per_unit": 200,
      "conversion_rate": 1.0,
      "wastage_percentage": 5.0,
      "total_cost": 6300
    },
    {
      "name": "Milk",
      "quantity": 150,
      "unit": "ml",
      "cost_per_unit": 50,
      "conversion_rate": 1.0,
      "wastage_percentage": 10.0,
      "total_cost": 8250
    }
  ],
  "total_cost": 15000
}
```


#### GET /api/menu/warnings

**Description**: Lấy danh sách menu items có warning (loss hoặc low margin)

**Query Parameters**:
- `threshold` (optional): Custom low margin threshold (default: from shop settings)

**Response**:
```json
{
  "loss_items": [
    {
      "menu_item_id": "...",
      "name": "Special Promo Coffee",
      "price": 20000,
      "current_cost": 25000,
      "profit_margin": -25.0,
      "absolute_profit": -5000,
      "warning_status": "loss"
    }
  ],
  "low_margin_items": [
    {
      "menu_item_id": "...",
      "name": "Basic Coffee",
      "price": 30000,
      "current_cost": 25000,
      "profit_margin": 16.67,
      "absolute_profit": 5000,
      "warning_status": "low_margin"
    }
  ],
  "loss_count": 1,
  "low_margin_count": 1,
  "threshold": 20.0
}
```

### Profit Analysis Endpoints

#### GET /api/reports/category-profit

**Description**: Lấy báo cáo profit theo category

**Query Parameters**:
- `start_date` (required): ISO 8601 date
- `end_date` (required): ISO 8601 date

**Response**:
```json
{
  "date_range": {
    "start": "2024-01-01",
    "end": "2024-01-31"
  },
  "categories": [
    {
      "category": "Coffee",
      "total_revenue": 5000000,
      "total_cost": 1500000,
      "total_profit": 3500000,
      "average_profit_margin": 70.0,
      "order_count": 150,
      "item_count": 200
    },
    {
      "category": "Tea",
      "total_revenue": 2000000,
      "total_cost": 800000,
      "total_profit": 1200000,
      "average_profit_margin": 60.0,
      "order_count": 80,
      "item_count": 100
    }
  ]
}
```

#### GET /api/reports/operating-profit

**Description**: Lấy báo cáo operating profit

**Query Parameters**:
- `start_date` (required): ISO 8601 date
- `end_date` (required): ISO 8601 date

**Response**:
```json
{
  "date_range": {
    "start": "2024-01-01",
    "end": "2024-01-31"
  },
  "total_revenue": 10000000,
  "total_cogs": 3000000,
  "gross_profit": 7000000,
  "gross_profit_margin": 70.0,
  "staff_salary": 2000000,
  "rent": 1000000,
  "utilities": 500000,
  "marketing_costs": 300000,
  "other_expenses": 200000,
  "total_expenses": 4000000,
  "operating_profit": 3000000,
  "operating_profit_margin": 30.0,
  "expense_allocated": false
}
```

### Operating Expense Endpoints

#### POST /api/operating-expenses

**Description**: Tạo hoặc update operating expense cho một kỳ

**Request Body**:
```json
{
  "period_start": "2024-01-01",
  "period_end": "2024-01-31",
  "staff_salary": 2000000,
  "rent": 1000000,
  "utilities": 500000,
  "marketing_costs": 300000,
  "other_expenses": 200000
}
```

**Response**:
```json
{
  "id": "...",
  "period_start": "2024-01-01",
  "period_end": "2024-01-31",
  "staff_salary": 2000000,
  "rent": 1000000,
  "utilities": 500000,
  "marketing_costs": 300000,
  "other_expenses": 200000,
  "total_expenses": 4000000,
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```


#### GET /api/operating-expenses

**Description**: Lấy danh sách operating expenses

**Query Parameters**:
- `start_date` (optional): Filter from date
- `end_date` (optional): Filter to date

**Response**:
```json
{
  "expenses": [
    {
      "id": "...",
      "period_start": "2024-01-01",
      "period_end": "2024-01-31",
      "staff_salary": 2000000,
      "rent": 1000000,
      "utilities": 500000,
      "marketing_costs": 300000,
      "other_expenses": 200000,
      "total_expenses": 4000000,
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ]
}
```

### Shift Closure Endpoint (Modified)

#### POST /api/shifts/:id/close

**Description**: Kết ca và trigger cost calculation cho tất cả orders trong shift

**Request Body**:
```json
{
  "end_cash": 5000000
}
```

**Response**:
```json
{
  "shift": {
    "id": "...",
    "status": "CLOSED",
    "ended_at": "2024-01-15T18:00:00Z",
    "total_revenue": 2000000,
    "total_orders": 50
  },
  "cost_calculation": {
    "total_orders_processed": 50,
    "total_items_processed": 120,
    "total_cogs": 600000,
    "items_with_incomplete_cost": 5
  }
}
```

**Backend Logic**:
1. Close shift (existing logic)
2. Trigger `CostCalculatorService.CalculateShiftOrderCosts(shiftID)`
3. Return shift data + cost calculation summary

### Settings Endpoint (Modified)

#### PATCH /api/settings

**Description**: Update shop settings including low margin threshold

**Request Body**:
```json
{
  "low_margin_threshold": 25.0
}
```

**Response**:
```json
{
  "id": "...",
  "shop_name": "My Cafe",
  "low_margin_threshold": 25.0,
  "updated_at": "2024-01-15T10:00:00Z"
}
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

Before writing correctness properties, I need to analyze each acceptance criterion from the requirements document to determine if it's testable as a property, example, or edge case.


### Property Reflection

After analyzing all acceptance criteria, I identified the following redundancies:

**Redundant Properties:**
1. Property 1.1 and 1.2 can be combined - both test the cost calculation formula
2. Property 1.3 and 9.1 are duplicates - both test background job queuing when ingredient cost changes
3. Property 2.1 and 2.5 can be combined - profit margin calculation and rounding
4. Property 2.6 can be combined with 2.1 - both are profit calculations
5. Property 5.2 and 5.3 can be combined - both test shift closure cost calculation
6. Property 6.1, 6.2, and 6.3 can be combined - all test category profit aggregation
7. Property 6.5.1, 6.5.3, and 6.5.4 can be combined - all test operating profit calculations
8. Property 10.1, 10.2, and 10.4 can be combined - all test the complete cost formula with conversion and wastage

**Consolidated Properties:**
- Combine 1.1 + 1.2 → Property 1: Cost calculation formula
- Remove 9.1 (duplicate of 1.3)
- Combine 2.1 + 2.5 + 2.6 → Property 2: Profit calculations
- Combine 5.2 + 5.3 → Property 5: Shift closure cost calculation
- Combine 6.1 + 6.2 + 6.3 → Property 6: Category profit aggregation
- Combine 6.5.1 + 6.5.3 + 6.5.4 → Property 7: Operating profit calculations
- Combine 10.1 + 10.2 + 10.4 → Property 10: Cost formula with conversion and wastage

### Correctness Properties


**Property 1: Cost Calculation Formula**

*For any* menu item with ingredients that have valid cost_per_unit values, calculating the current_cost should produce a result equal to the sum of (ingredient.quantity * ingredient.cost_per_unit * conversion_rate * (1 + wastage_percentage/100)) for all ingredients, rounded to 2 decimal places.

**Validates: Requirements 1.1, 1.2, 1.7, 10.1, 10.2, 10.4**

---

**Property 2: Profit Calculations**

*For any* menu item with valid cost and price values, the profit_margin should equal ((price - cost) / price) * 100 rounded to 2 decimal places, and the absolute_profit should equal (price - cost).

**Validates: Requirements 2.1, 2.5, 2.6**

---

**Property 3: Background Job Queuing on Ingredient Update**

*For any* ingredient, when its cost_per_unit is updated, the system should queue a background job to recalculate current_cost for all menu items that use that ingredient.

**Validates: Requirements 1.3, 9.1**

---

**Property 4: Incomplete Cost Status**

*For any* menu item where at least one ingredient has null or undefined cost_per_unit, the cost_status should be marked as "INCOMPLETE" and the item should not be included in profit calculations or reports.

**Validates: Requirements 1.5, 1.6, 2.9**

---

**Property 5: Shift Closure Cost Calculation**

*For any* shift, when the shift is closed, the system should calculate accounting_cost for all order items in that shift using the current cost_per_unit values at the time of closure, using the same calculation method as current menu item cost, and mark the cost_status as "FINAL".

**Validates: Requirements 5.2, 5.3**

---

**Property 6: Accounting Cost Immutability**

*For any* order item with cost_status = "FINAL" (from a closed shift), when ingredient cost_per_unit values change, the accounting_cost should remain unchanged.

**Validates: Requirements 5.8, 9.6**

---

**Property 7: Category Profit Aggregation**

*For any* category and date range, the category profit should be calculated as: total_revenue = sum of all order item revenues, total_cost = sum of all order item accounting_costs (not current_costs), total_profit = total_revenue - total_cost, and average_profit_margin = (total_profit / total_revenue) * 100.

**Validates: Requirements 6.1, 6.2, 6.3**

---

**Property 8: Operating Profit Calculations**

*For any* date range with orders and operating expenses, the operating_profit should equal gross_profit - total_expenses, where gross_profit = total_revenue - total_cogs, and operating_profit_margin = (operating_profit / total_revenue) * 100.

**Validates: Requirements 6.5.1, 6.5.3, 6.5.4**

---

**Property 9: Expense Allocation**

*For any* monthly operating expense and target date within that month, the allocated daily expense should equal monthly_expense / days_in_month, with expense_allocated flag set to true.

**Validates: Requirements 6.5.8**

---

**Property 10: Loss Detection**

*For any* menu item where cost exceeds price, the warning_status should be marked as "loss".

**Validates: Requirements 3.1**

---

**Property 11: Low Margin Detection**

*For any* menu item where profit_margin is below the configured low_margin_threshold and cost does not exceed price, the warning_status should be marked as "low_margin".

**Validates: Requirements 3.2**

---

**Property 12: Warning Status Transitions**

*For any* menu item, when cost or price changes such that the profit_margin crosses the low_margin_threshold or the cost-price relationship changes, the warning_status should update immediately to reflect the new state.

**Validates: Requirements 3.6**

---

**Property 13: Warning Count Aggregation**

*For any* set of menu items, the loss_count should equal the number of items with warning_status = "loss", and the low_margin_count should equal the number of items with warning_status = "low_margin".

**Validates: Requirements 3.5**

---

**Property 14: Category Filtering**

*For any* category filter, the API should return only menu items that belong to that category.

**Validates: Requirements 4.3**

---

**Property 15: Profit Margin Sorting**

*For any* sort order (ascending or descending), the API should return menu items sorted by profit_margin in the specified order.

**Validates: Requirements 4.4**

---

**Property 16: Date Range Filtering**

*For any* date range, the category profit analysis and operating profit analysis should include only orders where the order date falls within the specified range.

**Validates: Requirements 6.4, 6.5.6**

---

**Property 17: Batch Recalculation Optimization**

*For any* batch update of multiple ingredients, the system should recalculate affected menu items once after all updates complete, not once per ingredient update.

**Validates: Requirements 9.2**

---

**Property 18: Recipe Change Immutability**

*For any* menu item with historical accounting_cost data in closed shifts, when the recipe is modified, the historical accounting_cost values should remain unchanged.

**Validates: Requirements 8.6**

---

**Property 19: Summary Statistics Calculation**

*For any* set of menu items, the summary statistics should correctly calculate: total_items = count of all items, loss_count = count of items with cost > price, low_margin_count = count of items with profit_margin < threshold, average_profit_margin = average of all profit_margins (excluding items with price <= 0).

**Validates: Requirements 7.4**


## Error Handling

### Cost Calculation Errors

**Incomplete Ingredient Data**:
- **Error**: Ingredient missing cost_per_unit
- **Handling**: Mark menu item with cost_status = "INCOMPLETE", display warning "⚠ Thiếu giá nguyên liệu", exclude from profit reports
- **User Action**: Manager needs to update ingredient cost_per_unit

**Invalid Calculation Inputs**:
- **Error**: Negative quantity, negative cost_per_unit, invalid conversion_rate
- **Handling**: Log error, return error response with descriptive message, do not save invalid data
- **User Action**: Fix input data

**Background Job Failures**:
- **Error**: Cost recalculation job fails (database error, timeout)
- **Handling**: Log error, retry with exponential backoff (max 3 retries), mark recalculation_status as failed
- **User Action**: System admin investigates logs, manual retry if needed

### Shift Closure Errors

**Cost Calculation Timeout**:
- **Error**: Shift has too many orders, cost calculation exceeds timeout
- **Handling**: Process in batches, continue with next batch, log progress
- **User Action**: None, system handles automatically

**Partial Cost Calculation**:
- **Error**: Some order items have incomplete ingredient data
- **Handling**: Mark those items with cost_status = "INCOMPLETE", continue with other items, return summary with incomplete count
- **User Action**: Manager reviews incomplete items, updates ingredient costs

### API Errors

**Invalid Date Range**:
- **Error**: start_date > end_date, invalid date format
- **Handling**: Return 400 Bad Request with error message
- **User Action**: Fix date range

**Missing Operating Expenses**:
- **Error**: No operating expense data for requested period
- **Handling**: Return gross_profit only, set expense_allocated = false, include note "Chưa nhập chi phí vận hành"
- **User Action**: Manager inputs operating expenses

**Database Connection Errors**:
- **Error**: MongoDB connection lost, query timeout
- **Handling**: Return 503 Service Unavailable, retry with exponential backoff
- **User Action**: Wait and retry, contact system admin if persists

### Data Consistency Errors

**Orphaned Order Items**:
- **Error**: Order item references non-existent menu item
- **Handling**: Use stored name and price from order item, mark cost_status = "INCOMPLETE"
- **User Action**: None, historical data preserved

**Missing Ingredient in Recipe**:
- **Error**: Menu item recipe references ingredient that no longer exists
- **Handling**: Skip that ingredient in cost calculation, mark cost_status = "INCOMPLETE", log warning
- **User Action**: Manager updates recipe or restores ingredient

## Testing Strategy

### Dual Testing Approach

The testing strategy combines **unit tests** and **property-based tests** to ensure comprehensive coverage:

- **Unit tests**: Verify specific examples, edge cases, error conditions, and integration points
- **Property tests**: Verify universal properties across all inputs through randomization

Both approaches are complementary and necessary for comprehensive coverage. Unit tests catch concrete bugs in specific scenarios, while property tests verify general correctness across a wide range of inputs.

### Unit Testing

**Focus Areas**:
- Specific examples demonstrating correct behavior
- Edge cases: empty ingredients, zero price, negative profit, break-even items
- Error conditions: missing data, invalid inputs, database errors
- Integration points: API endpoints, database operations, service interactions

**Example Unit Tests**:
```go
// Example: Cost calculation for a specific menu item
func TestCalculateMenuItemCost_Cappuccino(t *testing.T) {
    // Given a Cappuccino with known ingredients
    menuItem := &MenuItem{
        Name: "Cappuccino",
        Ingredients: []Ingredient{
            {Name: "Espresso", Quantity: 30, CostPerUnit: 200},
            {Name: "Milk", Quantity: 150, CostPerUnit: 50},
        },
    }
    
    // When calculating cost
    cost, err := calculator.CalculateMenuItemCost(menuItem)
    
    // Then cost should be 30*200 + 150*50 = 13500
    assert.NoError(t, err)
    assert.Equal(t, 13500.0, cost)
}

// Example: Edge case - menu item with no ingredients
func TestCalculateMenuItemCost_NoIngredients(t *testing.T) {
    menuItem := &MenuItem{Name: "Service", Ingredients: []Ingredient{}}
    cost, status, err := calculator.CalculateMenuItemCost(menuItem)
    
    assert.NoError(t, err)
    assert.Equal(t, 0.0, cost)
    assert.Equal(t, CostStatusFinal, status)
}

// Example: Error condition - missing ingredient cost
func TestCalculateMenuItemCost_MissingCost(t *testing.T) {
    menuItem := &MenuItem{
        Name: "Latte",
        Ingredients: []Ingredient{
            {Name: "Espresso", Quantity: 30, CostPerUnit: 0}, // Missing cost
        },
    }
    
    cost, status, err := calculator.CalculateMenuItemCost(menuItem)
    
    assert.NoError(t, err)
    assert.Equal(t, CostStatusIncomplete, status)
}
```


### Property-Based Testing

**Configuration**:
- Use **testify/assert** for Go testing framework
- Use **gopter** library for property-based testing in Go
- Minimum **100 iterations** per property test
- Each property test must reference its design document property

**Tag Format**:
```go
// Feature: menu-cost-profit-analysis, Property 1: Cost calculation formula
```

**Property Test Examples**:

```go
// Feature: menu-cost-profit-analysis, Property 1: Cost calculation formula
func TestProperty_CostCalculationFormula(t *testing.T) {
    properties := gopter.NewProperties(nil)
    
    properties.Property("Cost equals sum of ingredient costs", prop.ForAll(
        func(ingredients []Ingredient) bool {
            // Generate random menu item with random ingredients
            menuItem := &MenuItem{
                Name: "Random Item",
                Ingredients: ingredients,
            }
            
            // Calculate cost
            cost, status, _ := calculator.CalculateMenuItemCost(menuItem)
            
            // Skip if incomplete
            if status == CostStatusIncomplete {
                return true
            }
            
            // Calculate expected cost manually
            expectedCost := 0.0
            for _, ing := range ingredients {
                conversionRate := ing.ConversionRate
                if conversionRate == 0 {
                    conversionRate = 1.0
                }
                wastage := 1.0 + (ing.WastagePercentage / 100.0)
                expectedCost += ing.Quantity * ing.CostPerUnit * conversionRate * wastage
            }
            expectedCost = math.Round(expectedCost*100) / 100 // Round to 2 decimals
            
            // Verify cost matches expected
            return math.Abs(cost - expectedCost) < 0.01
        },
        genIngredients(), // Generator for random ingredients
    ))
    
    properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 2: Profit calculations
func TestProperty_ProfitCalculations(t *testing.T) {
    properties := gopter.NewProperties(nil)
    
    properties.Property("Profit margin and absolute profit are correct", prop.ForAll(
        func(price, cost float64) bool {
            // Skip invalid cases
            if price <= 0 {
                return true
            }
            
            // Calculate profit metrics
            profitMargin := ((price - cost) / price) * 100
            profitMargin = math.Round(profitMargin*100) / 100
            absoluteProfit := price - cost
            
            // Calculate using service
            result := analyzer.CalculateProfit(price, cost)
            
            // Verify both metrics
            return math.Abs(result.ProfitMargin - profitMargin) < 0.01 &&
                   math.Abs(result.AbsoluteProfit - absoluteProfit) < 0.01
        },
        gen.Float64Range(1, 100000),  // price
        gen.Float64Range(0, 100000),  // cost
    ))
    
    properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 6: Accounting cost immutability
func TestProperty_AccountingCostImmutability(t *testing.T) {
    properties := gopter.NewProperties(nil)
    
    properties.Property("Accounting cost doesn't change when ingredient cost changes", prop.ForAll(
        func(orderItem OrderItem, newIngredientCost float64) bool {
            // Given an order item with FINAL cost status
            orderItem.CostStatus = CostStatusFinal
            originalCost := orderItem.AccountingCost
            
            // When ingredient cost changes
            ingredient := getIngredient(orderItem.MenuItemID)
            ingredient.CostPerUnit = newIngredientCost
            updateIngredient(ingredient)
            
            // Then accounting cost should remain unchanged
            updatedOrderItem := getOrderItem(orderItem.ID)
            return updatedOrderItem.AccountingCost == originalCost
        },
        genOrderItem(),
        gen.Float64Range(1, 10000),
    ))
    
    properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 10: Loss detection
func TestProperty_LossDetection(t *testing.T) {
    properties := gopter.NewProperties(nil)
    
    properties.Property("Items with cost > price are marked as loss", prop.ForAll(
        func(price, cost float64) bool {
            menuItem := &MenuItemCost{
                Price: price,
                CurrentCost: cost,
            }
            
            warningStatus := analyzer.DetectWarningStatus(menuItem, 20.0)
            
            if cost > price {
                return warningStatus == WarningLoss
            }
            return true
        },
        gen.Float64Range(1, 100000),
        gen.Float64Range(1, 100000),
    ))
    
    properties.TestingRun(t, gopter.ConsoleReporter(false))
}
```

**Generators**:
```go
// Generator for random ingredients with valid cost_per_unit
func genIngredients() gopter.Gen {
    return gen.SliceOf(gen.Struct(reflect.TypeOf(Ingredient{}), map[string]gopter.Gen{
        "Name":              gen.Identifier(),
        "Quantity":          gen.Float64Range(1, 1000),
        "CostPerUnit":       gen.Float64Range(1, 10000),
        "ConversionRate":    gen.Float64Range(0.1, 10),
        "WastagePercentage": gen.Float64Range(0, 50),
    }))
}

// Generator for random order items
func genOrderItem() gopter.Gen {
    return gen.Struct(reflect.TypeOf(OrderItem{}), map[string]gopter.Gen{
        "MenuItemID":     genObjectID(),
        "Name":           gen.Identifier(),
        "Price":          gen.Float64Range(1, 100000),
        "Quantity":       gen.IntRange(1, 10),
        "AccountingCost": gen.Float64Range(1, 50000),
        "CostStatus":     gen.Const(CostStatusFinal),
    })
}
```

### Integration Testing

**Focus Areas**:
- End-to-end API flows
- Database operations and transactions
- Background job processing
- Shift closure workflow

**Example Integration Tests**:
```go
func TestIntegration_ShiftClosureCostCalculation(t *testing.T) {
    // Setup: Create shift with orders
    shift := createTestShift()
    orders := createTestOrders(shift.ID, 10)
    
    // Execute: Close shift
    response := closeShift(shift.ID)
    
    // Verify: All order items have accounting_cost
    for _, order := range orders {
        items := getOrderItems(order.ID)
        for _, item := range items {
            assert.NotZero(t, item.AccountingCost)
            assert.Equal(t, CostStatusFinal, item.CostStatus)
            assert.NotZero(t, item.CostCalculatedAt)
        }
    }
    
    // Verify: Cost calculation summary
    assert.Equal(t, 10, response.CostCalculation.TotalOrdersProcessed)
}
```

### Frontend Testing

**Unit Tests (Vue Test Utils + Vitest)**:
- Component rendering and props
- User interactions (click, input)
- State management
- API integration (mocked)

**Example Frontend Tests**:
```typescript
describe('MenuCostView', () => {
  it('displays menu items with cost information', async () => {
    const wrapper = mount(MenuCostView, {
      global: {
        plugins: [createTestingPinia()]
      }
    })
    
    await wrapper.vm.$nextTick()
    
    expect(wrapper.find('.menu-item').exists()).toBe(true)
    expect(wrapper.find('.current-cost').text()).toContain('15,000')
    expect(wrapper.find('.profit-margin').text()).toContain('66.67%')
  })
  
  it('filters items by category', async () => {
    const wrapper = mount(MenuCostView)
    
    await wrapper.find('[data-category="Coffee"]').trigger('click')
    
    const items = wrapper.findAll('.menu-item')
    items.forEach(item => {
      expect(item.find('.category').text()).toBe('Coffee')
    })
  })
  
  it('displays warning for loss items', async () => {
    const wrapper = mount(MenuCostView, {
      props: {
        items: [
          { name: 'Loss Item', cost: 50000, price: 30000, warning_status: 'loss' }
        ]
      }
    })
    
    expect(wrapper.find('.warning-loss').exists()).toBe(true)
    expect(wrapper.find('.warning-loss').classes()).toContain('text-red-500')
  })
})
```

### Performance Testing

**Load Testing**:
- Cost recalculation for 1000+ menu items
- Shift closure with 500+ orders
- Category profit aggregation for 10,000+ orders

**Benchmarks**:
```go
func BenchmarkCostCalculation(b *testing.B) {
    menuItems := generateMenuItems(1000)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        for _, item := range menuItems {
            calculator.CalculateMenuItemCost(item)
        }
    }
}

func BenchmarkShiftClosure(b *testing.B) {
    shift := createTestShift()
    orders := createTestOrders(shift.ID, 500)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        calculator.CalculateShiftOrderCosts(shift.ID)
    }
}
```

## Performance Optimization

### Database Optimization

**Indexes**:
- `menu_items`: (category), (cost_status), (current_cost)
- `order_items`: (order_id), (menu_item_id), (cost_status), (cost_calculated_at)
- `operating_expenses`: (period_start, period_end), (period_start)

**Query Optimization**:
- Use aggregation pipeline for category profit calculations
- Batch fetch ingredients for multiple menu items
- Use projection to fetch only required fields

**Caching Strategy**:
- Cache menu item costs for 5 minutes (invalidate on ingredient update)
- Cache category profit reports for 1 hour (invalidate on shift closure)
- Use Redis for distributed caching

### Background Job Optimization

**Batch Processing**:
- Process cost recalculation in batches of 100 menu items
- Use goroutines for parallel processing (worker pool pattern)
- Implement exponential backoff for retries

**Queue Management**:
- Use Redis or Go channels for job queue
- Deduplicate jobs (if same ingredient updated multiple times)
- Priority queue: shift closure > ingredient update

### Frontend Optimization

**Data Loading**:
- Lazy load menu item details (cost breakdown)
- Paginate menu cost table (50 items per page)
- Debounce filter and sort operations

**State Management**:
- Use Pinia for centralized state
- Cache API responses in store
- Implement optimistic updates for better UX

**UI Performance**:
- Virtual scrolling for large tables
- Memoize computed properties
- Use v-show instead of v-if for frequently toggled elements

## Security Considerations

**Authorization**:
- Only managers can view cost and profit reports
- Only managers can input operating expenses
- Only managers can configure low margin threshold

**Data Validation**:
- Validate all numeric inputs (cost, price, expenses) >= 0
- Validate date ranges (start <= end)
- Sanitize user inputs to prevent injection attacks

**Audit Logging**:
- Log all operating expense changes with user ID and timestamp
- Log all low margin threshold changes
- Log shift closure cost calculations

## Migration Strategy

**Phase 1: Schema Migration**
- Add new fields to menu_items collection
- Create order_items collection
- Create operating_expenses collection
- Add indexes

**Phase 2: Data Migration**
- Calculate current_cost for all existing menu items
- Set cost_status = FINAL for items without ingredients
- Set cost_status = INCOMPLETE for items with missing ingredient costs

**Phase 3: Backfill Historical Data**
- For closed shifts, calculate accounting_cost using current ingredient costs
- Mark as cost_status = ESTIMATED (not FINAL since not calculated at shift closure time)
- Add note indicating backfilled data

**Phase 4: Feature Rollout**
- Deploy backend API endpoints
- Deploy frontend manager views
- Enable background job processing
- Monitor performance and errors

**Rollback Plan**:
- Keep old schema fields intact during migration
- Feature flag to disable new cost calculation
- Ability to revert to old cost calculation method

## Monitoring and Observability

**Metrics**:
- Cost recalculation job duration and success rate
- Shift closure cost calculation duration
- API response times for cost and profit endpoints
- Number of menu items with INCOMPLETE status

**Alerts**:
- Cost recalculation job failures
- Shift closure cost calculation timeout
- High number of INCOMPLETE menu items (> 10%)
- API error rate > 5%

**Logging**:
- Log all cost calculations with input/output
- Log all background job executions
- Log all API errors with stack traces
- Log all data validation failures

## Future Enhancements (Out of Scope for Phase 1)

**Combo Items**:
- Support menu items composed of other menu items
- Recursive cost calculation for nested combos

**Tax and Fees**:
- Calculate VAT, service fees, delivery fees
- Net profit after taxes

**Multi-location Cost**:
- Different ingredient costs per location
- Location-specific profit analysis

**Predictive Analytics**:
- Forecast ingredient cost trends
- Recommend optimal pricing based on target margin
- Identify seasonal cost variations

**Advanced Reporting**:
- Profit by time of day, day of week
- Customer segment profitability
- Menu item performance matrix (profit vs popularity)

