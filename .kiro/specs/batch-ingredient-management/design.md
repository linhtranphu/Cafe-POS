# Tài Liệu Thiết Kế: Quản Lý Nguyên Liệu Batch

## Tổng Quan

Hệ thống Quản Lý Nguyên Liệu Batch cung cấp khả năng theo dõi và quản lý các nguyên liệu trung gian được chế biến từ nguyên liệu thô. Hệ thống tích hợp với các hệ thống hiện có (inventory, menu, cost calculator) để cung cấp giải pháp toàn diện từ nguyên liệu thô đến sản phẩm cuối cùng.

### Mục Tiêu Thiết Kế

- **Tính toàn vẹn dữ liệu**: Đảm bảo tồn kho luôn chính xác thông qua transactions
- **Hiệu suất**: Xử lý nhanh với caching và indexing phù hợp
- **Khả năng mở rộng**: Thiết kế cho phép thêm tính năng mới dễ dàng
- **Tích hợp**: Tích hợp mượt mà với các hệ thống hiện có
- **Audit trail**: Theo dõi đầy đủ lịch sử thay đổi

## Kiến Trúc

### Kiến Trúc Tổng Thể

```
┌─────────────────────────────────────────────────────────────┐
│                        Frontend (Vue.js)                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Batch Form   │  │ Batch List   │  │ Alert Panel  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ HTTP/JSON
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Backend (Go)                            │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              HTTP Handlers Layer                      │   │
│  │  - BatchDefinitionHandler                            │   │
│  │  - BatchRecordHandler                                │   │
│  │  - BatchAlertHandler                                 │   │
│  └──────────────────────────────────────────────────────┘   │
│                              │                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │           Application Services Layer                  │   │
│  │  - BatchDefinitionService                            │   │
│  │  - BatchRecordService                                │   │
│  │  - BatchCostCalculator                               │   │
│  │  - BatchAlertService                                 │   │
│  │  - BatchUsageService (FIFO)                          │   │
│  └──────────────────────────────────────────────────────┘   │
│                              │                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Domain Layer                             │   │
│  │  - BatchDefinition (entity)                          │   │
│  │  - BatchRecord (entity)                              │   │
│  │  - ConversionRate (value object)                     │   │
│  │  - BatchUsage (value object)                         │   │
│  └──────────────────────────────────────────────────────┘   │
│                              │                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │           Repository Layer                            │   │
│  │  - BatchDefinitionRepository                         │   │
│  │  - BatchRecordRepository                             │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      MongoDB                                 │
│  - batch_definitions collection                             │
│  - batch_records collection                                 │
│  - batch_usage_logs collection                              │
└─────────────────────────────────────────────────────────────┘
                              │
                    Integration với
                              │
┌─────────────────────────────────────────────────────────────┐
│              Existing Systems                                │
│  - Inventory System (ingredient stock)                      │
│  - Menu System (recipes)                                    │
│  - Cost Calculator (ingredient costs)                       │
└─────────────────────────────────────────────────────────────┘
```

### Luồng Dữ Liệu Chính

**1. Tạo Batch Definition:**
```
User → Frontend → BatchDefinitionHandler → BatchDefinitionService 
→ Validate ingredients exist → Save to MongoDB
```

**2. Ghi Nhận Batch Preparation:**
```
User → Frontend → BatchRecordHandler → BatchRecordService
→ Calculate required ingredients (with wastage)
→ Check inventory availability
→ Start MongoDB transaction
  → Deduct ingredients from inventory
  → Create batch record
  → Calculate and store cost
  → Commit transaction
→ Return success
```

**3. Sử dụng Batch trong Order:**
```
Order Created → Menu System → BatchUsageService
→ Get available batches (FIFO, not expired)
→ Start MongoDB transaction
  → Deduct batch quantity
  → Log usage
  → Update batch record
  → Commit transaction
→ Return actual cost
```

**4. Kiểm Tra Alerts:**
```
Scheduled Job (every 5 minutes) → BatchAlertService
→ Check low stock batches
→ Check expiring batches
→ Check expired batches
→ Update alert cache
→ Notify frontend via WebSocket (optional)
```

## Các Thành Phần và Giao Diện

### 1. Domain Entities

#### BatchDefinition
```go
type BatchDefinition struct {
    ID                   primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    Name                 string               `bson:"name" json:"name"`
    Unit                 string               `bson:"unit" json:"unit"`
    ShelfLifeHours       int                  `bson:"shelf_life_hours" json:"shelf_life_hours"`
    ConversionRates      []ConversionRate     `bson:"conversion_rates" json:"conversion_rates"`
    LowStockThreshold    float64              `bson:"low_stock_threshold" json:"low_stock_threshold"`
    ExpiryWarningHours   int                  `bson:"expiry_warning_hours" json:"expiry_warning_hours"`
    CreatedAt            time.Time            `bson:"created_at" json:"created_at"`
    UpdatedAt            time.Time            `bson:"updated_at" json:"updated_at"`
}

type ConversionRate struct {
    SourceIngredientID   primitive.ObjectID   `bson:"source_ingredient_id" json:"source_ingredient_id"`
    SourceIngredientName string               `bson:"source_ingredient_name" json:"source_ingredient_name"`
    SourceQuantity       float64              `bson:"source_quantity" json:"source_quantity"`
    SourceUnit           string               `bson:"source_unit" json:"source_unit"`
    BatchQuantity        float64              `bson:"batch_quantity" json:"batch_quantity"`
    WastageRate          float64              `bson:"wastage_rate" json:"wastage_rate"` // 0.0 to 1.0
}
```

#### BatchRecord
```go
type BatchRecord struct {
    ID                   primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    BatchDefinitionID    primitive.ObjectID   `bson:"batch_definition_id" json:"batch_definition_id"`
    BatchName            string               `bson:"batch_name" json:"batch_name"`
    QuantityProduced     float64              `bson:"quantity_produced" json:"quantity_produced"`
    QuantityRemaining    float64              `bson:"quantity_remaining" json:"quantity_remaining"`
    Unit                 string               `bson:"unit" json:"unit"`
    CostPerUnit          float64              `bson:"cost_per_unit" json:"cost_per_unit"`
    TotalCost            float64              `bson:"total_cost" json:"total_cost"`
    PreparedBy           string               `bson:"prepared_by" json:"prepared_by"`
    PreparedAt           time.Time            `bson:"prepared_at" json:"prepared_at"`
    ExpiresAt            time.Time            `bson:"expires_at" json:"expires_at"`
    Status               string               `bson:"status" json:"status"` // "available", "expired", "depleted"
    IngredientsUsed      []IngredientUsage    `bson:"ingredients_used" json:"ingredients_used"`
    CreatedAt            time.Time            `bson:"created_at" json:"created_at"`
    UpdatedAt            time.Time            `bson:"updated_at" json:"updated_at"`
}

type IngredientUsage struct {
    IngredientID         primitive.ObjectID   `bson:"ingredient_id" json:"ingredient_id"`
    IngredientName       string               `bson:"ingredient_name" json:"ingredient_name"`
    Quantity             float64              `bson:"quantity" json:"quantity"`
    Unit                 string               `bson:"unit" json:"unit"`
    CostPerUnit          float64              `bson:"cost_per_unit" json:"cost_per_unit"`
    TotalCost            float64              `bson:"total_cost" json:"total_cost"`
}
```

#### BatchUsageLog
```go
type BatchUsageLog struct {
    ID                   primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    BatchRecordID        primitive.ObjectID   `bson:"batch_record_id" json:"batch_record_id"`
    BatchName            string               `bson:"batch_name" json:"batch_name"`
    OrderID              primitive.ObjectID   `bson:"order_id" json:"order_id"`
    MenuItemID           primitive.ObjectID   `bson:"menu_item_id" json:"menu_item_id"`
    MenuItemName         string               `bson:"menu_item_name" json:"menu_item_name"`
    QuantityUsed         float64              `bson:"quantity_used" json:"quantity_used"`
    Unit                 string               `bson:"unit" json:"unit"`
    CostPerUnit          float64              `bson:"cost_per_unit" json:"cost_per_unit"`
    TotalCost            float64              `bson:"total_cost" json:"total_cost"`
    UsedAt               time.Time            `bson:"used_at" json:"used_at"`
}
```

### 2. Service Interfaces

#### BatchDefinitionService
```go
type BatchDefinitionService interface {
    Create(ctx context.Context, def *BatchDefinition) error
    Update(ctx context.Context, id primitive.ObjectID, def *BatchDefinition) error
    Delete(ctx context.Context, id primitive.ObjectID) error
    GetByID(ctx context.Context, id primitive.ObjectID) (*BatchDefinition, error)
    List(ctx context.Context, filter BatchDefinitionFilter) ([]*BatchDefinition, error)
    ValidateConversionRates(ctx context.Context, rates []ConversionRate) error
}
```

#### BatchRecordService
```go
type BatchRecordService interface {
    CreateBatch(ctx context.Context, req CreateBatchRequest) (*BatchRecord, error)
    GetByID(ctx context.Context, id primitive.ObjectID) (*BatchRecord, error)
    List(ctx context.Context, filter BatchRecordFilter) ([]*BatchRecord, error)
    UpdateQuantity(ctx context.Context, id primitive.ObjectID, newQuantity float64) error
    MarkAsExpired(ctx context.Context, id primitive.ObjectID) error
    Delete(ctx context.Context, id primitive.ObjectID) error
    GetAvailableBatches(ctx context.Context, batchDefID primitive.ObjectID) ([]*BatchRecord, error)
}

type CreateBatchRequest struct {
    BatchDefinitionID    primitive.ObjectID
    QuantityProduced     float64
    PreparedBy           string
}
```

#### BatchCostCalculator
```go
type BatchCostCalculator interface {
    CalculateBatchCost(ctx context.Context, def *BatchDefinition, quantity float64) (CostBreakdown, error)
}

type CostBreakdown struct {
    TotalCost            float64
    CostPerUnit          float64
    IngredientCosts      []IngredientCost
}

type IngredientCost struct {
    IngredientID         primitive.ObjectID
    IngredientName       string
    Quantity             float64
    Unit                 string
    CostPerUnit          float64
    TotalCost            float64
}
```

#### BatchUsageService
```go
type BatchUsageService interface {
    UseBatch(ctx context.Context, req UseBatchRequest) (*BatchUsageResult, error)
    GetUsageHistory(ctx context.Context, filter UsageFilter) ([]*BatchUsageLog, error)
}

type UseBatchRequest struct {
    BatchDefinitionID    primitive.ObjectID
    QuantityNeeded       float64
    OrderID              primitive.ObjectID
    MenuItemID           primitive.ObjectID
    MenuItemName         string
}

type BatchUsageResult struct {
    Success              bool
    BatchesUsed          []BatchUsageDetail
    TotalCost            float64
    Message              string
}

type BatchUsageDetail struct {
    BatchRecordID        primitive.ObjectID
    QuantityUsed         float64
    CostPerUnit          float64
}
```

#### BatchAlertService
```go
type BatchAlertService interface {
    GetAlerts(ctx context.Context) (*BatchAlerts, error)
    CheckLowStock(ctx context.Context) ([]*LowStockAlert, error)
    CheckExpiring(ctx context.Context) ([]*ExpiringAlert, error)
    CheckExpired(ctx context.Context) ([]*ExpiredAlert, error)
}

type BatchAlerts struct {
    LowStock             []*LowStockAlert
    Expiring             []*ExpiringAlert
    Expired              []*ExpiredAlert
    LastChecked          time.Time
}

type LowStockAlert struct {
    BatchDefinitionID    primitive.ObjectID
    BatchName            string
    CurrentStock         float64
    Threshold            float64
    Unit                 string
}

type ExpiringAlert struct {
    BatchRecordID        primitive.ObjectID
    BatchName            string
    QuantityRemaining    float64
    Unit                 string
    ExpiresAt            time.Time
    HoursUntilExpiry     int
}

type ExpiredAlert struct {
    BatchRecordID        primitive.ObjectID
    BatchName            string
    QuantityWasted       float64
    Unit                 string
    CostWasted           float64
    ExpiredAt            time.Time
}
```

### 3. Repository Interfaces

#### BatchDefinitionRepository
```go
type BatchDefinitionRepository interface {
    Create(ctx context.Context, def *BatchDefinition) error
    Update(ctx context.Context, def *BatchDefinition) error
    Delete(ctx context.Context, id primitive.ObjectID) error
    FindByID(ctx context.Context, id primitive.ObjectID) (*BatchDefinition, error)
    FindAll(ctx context.Context, filter BatchDefinitionFilter) ([]*BatchDefinition, error)
}
```

#### BatchRecordRepository
```go
type BatchRecordRepository interface {
    Create(ctx context.Context, record *BatchRecord) error
    Update(ctx context.Context, record *BatchRecord) error
    Delete(ctx context.Context, id primitive.ObjectID) error
    FindByID(ctx context.Context, id primitive.ObjectID) (*BatchRecord, error)
    FindAll(ctx context.Context, filter BatchRecordFilter) ([]*BatchRecord, error)
    FindAvailableByDefinition(ctx context.Context, defID primitive.ObjectID) ([]*BatchRecord, error)
    UpdateQuantity(ctx context.Context, id primitive.ObjectID, newQuantity float64) error
    GetTotalAvailableQuantity(ctx context.Context, defID primitive.ObjectID) (float64, error)
}
```

## Mô Hình Dữ Liệu

### MongoDB Collections

#### batch_definitions
```javascript
{
  _id: ObjectId,
  name: String,
  unit: String,
  shelf_life_hours: Number,
  conversion_rates: [
    {
      source_ingredient_id: ObjectId,
      source_ingredient_name: String,
      source_quantity: Number,
      source_unit: String,
      batch_quantity: Number,
      wastage_rate: Number
    }
  ],
  low_stock_threshold: Number,
  expiry_warning_hours: Number,
  created_at: ISODate,
  updated_at: ISODate
}

// Indexes
db.batch_definitions.createIndex({ name: 1 })
db.batch_definitions.createIndex({ created_at: -1 })
```

#### batch_records
```javascript
{
  _id: ObjectId,
  batch_definition_id: ObjectId,
  batch_name: String,
  quantity_produced: Number,
  quantity_remaining: Number,
  unit: String,
  cost_per_unit: Number,
  total_cost: Number,
  prepared_by: String,
  prepared_at: ISODate,
  expires_at: ISODate,
  status: String, // "available", "expired", "depleted"
  ingredients_used: [
    {
      ingredient_id: ObjectId,
      ingredient_name: String,
      quantity: Number,
      unit: String,
      cost_per_unit: Number,
      total_cost: Number
    }
  ],
  created_at: ISODate,
  updated_at: ISODate
}

// Indexes
db.batch_records.createIndex({ batch_definition_id: 1, expires_at: 1 })
db.batch_records.createIndex({ status: 1, expires_at: 1 })
db.batch_records.createIndex({ expires_at: 1 })
db.batch_records.createIndex({ prepared_at: -1 })
```

#### batch_usage_logs
```javascript
{
  _id: ObjectId,
  batch_record_id: ObjectId,
  batch_name: String,
  order_id: ObjectId,
  menu_item_id: ObjectId,
  menu_item_name: String,
  quantity_used: Number,
  unit: String,
  cost_per_unit: Number,
  total_cost: Number,
  used_at: ISODate
}

// Indexes
db.batch_usage_logs.createIndex({ batch_record_id: 1, used_at: -1 })
db.batch_usage_logs.createIndex({ order_id: 1 })
db.batch_usage_logs.createIndex({ menu_item_id: 1, used_at: -1 })
db.batch_usage_logs.createIndex({ used_at: -1 })
```

### Quan Hệ Dữ Liệu

```
batch_definitions (1) ──→ (N) batch_records
                                    │
                                    │ (1)
                                    ↓
                                   (N)
                            batch_usage_logs
                                    │
                                    │
                                    ↓
                            orders (existing)
                            menu_items (existing)

batch_definitions.conversion_rates ──→ ingredients (existing)
batch_records.ingredients_used ──→ ingredients (existing)
```

## Thuộc Tính Đúng Đắn (Correctness Properties)

Thuộc tính đúng đắn (property) là một đặc điểm hoặc hành vi phải đúng trong tất cả các lần thực thi hợp lệ của hệ thống - về cơ bản là một tuyên bố chính thức về những gì hệ thống nên làm. Các thuộc tính này đóng vai trò là cầu nối giữa đặc tả có thể đọc được bởi con người và các đảm bảo tính đúng đắn có thể xác minh được bằng máy.


### Property 1: Inventory Conservation (Bảo Toàn Tồn Kho)

**Validates: Requirements 2.2, 2.3, 5.2**

**Mô tả:** Tổng số lượng nguyên liệu trong hệ thống (bao gồm cả nguyên liệu thô và batch) phải được bảo toàn. Khi tạo batch, nguyên liệu thô giảm đúng bằng số lượng được sử dụng (bao gồm wastage). Khi sử dụng batch, số lượng batch giảm đúng bằng số lượng được dùng.

**Formal Statement:**
```
∀ batch_creation:
  source_ingredient_before - source_ingredient_after = 
    (batch_quantity / conversion_rate) * (1 + wastage_rate)

∀ batch_usage:
  batch_quantity_before - batch_quantity_after = quantity_used
```

**Test Strategy:**
- Generate random batch creation operations
- Track ingredient quantities before and after
- Verify conservation equation holds
- Test with various wastage rates (0% to 50%)

### Property 2: Cost Accuracy (Chính Xác Chi Phí)

**Validates: Requirements 3.1, 3.2, 3.5**

**Mô tả:** Chi phí batch phải được tính chính xác từ chi phí nguyên liệu nguồn, bao gồm cả wastage. Chi phí này phải được lưu trữ tại thời điểm tạo batch và không thay đổi khi chi phí nguyên liệu thay đổi.

**Formal Statement:**
```
∀ batch_record:
  batch_record.total_cost = Σ(ingredient.quantity * ingredient.cost_per_unit)
  batch_record.cost_per_unit = batch_record.total_cost / batch_record.quantity_produced
  
∀ ingredient in batch_record.ingredients_used:
  ingredient.quantity = base_quantity * (1 + wastage_rate)
```

**Test Strategy:**
- Create batches with known ingredient costs
- Verify total cost calculation
- Verify cost per unit calculation
- Test that stored cost doesn't change when ingredient prices change
- Test with multiple ingredients and different wastage rates

### Property 3: FIFO Ordering (Thứ Tự FIFO)

**Validates: Requirements 5.3, 5.4**

**Mô tả:** Khi sử dụng batch, hệ thống phải ưu tiên batch cũ nhất (sắp hết hạn nhất) trong số các batch khả dụng. Điều này đảm bảo giảm thiểu lãng phí do hết hạn.

**Formal Statement:**
```
∀ batch_usage:
  LET available_batches = batches WHERE status = "available" AND expires_at > now
  LET sorted_batches = SORT available_batches BY expires_at ASC
  THEN first_used_batch = sorted_batches[0]
```

**Test Strategy:**
- Create multiple batches with different expiry times
- Trigger batch usage
- Verify that batch with earliest expiry is used first
- Test with partial usage (not enough in first batch)
- Verify correct handling when multiple batches needed

### Property 4: Expiry Enforcement (Ép Buộc Hết Hạn)

**Validates: Requirements 4.3, 5.5**

**Mô tả:** Batch đã hết hạn không được phép sử dụng. Hệ thống phải từ chối mọi thao tác sử dụng batch đã hết hạn.

**Formal Statement:**
```
∀ batch_record:
  IF now > batch_record.expires_at
  THEN batch_record.status = "expired"
  AND batch_record CANNOT be used in any order

∀ batch_usage_attempt:
  IF batch.expires_at <= now
  THEN usage_attempt MUST fail
```

**Test Strategy:**
- Create batches with short shelf life
- Wait for expiry or mock time
- Attempt to use expired batch
- Verify rejection with appropriate error
- Test boundary conditions (exactly at expiry time)

### Property 5: Alert Correctness (Đúng Đắn Cảnh Báo)

**Validates: Requirements 4.1, 4.2, 4.5**

**Mô tả:** Cảnh báo phải được hiển thị chính xác dựa trên ngưỡng đã định nghĩa. Low stock alert khi tổng số lượng khả dụng <= threshold. Expiry alert khi thời gian còn lại <= warning hours.

**Formal Statement:**
```
∀ batch_definition:
  LET total_available = Σ(batch_records WHERE status = "available").quantity_remaining
  IF total_available <= batch_definition.low_stock_threshold
  THEN low_stock_alert MUST be shown

∀ batch_record WHERE status = "available":
  LET hours_until_expiry = (expires_at - now) / 3600
  IF hours_until_expiry <= expiry_warning_hours
  THEN expiry_alert MUST be shown
```

**Test Strategy:**
- Create batch definitions with specific thresholds
- Create batch records with varying quantities and expiry times
- Query alert service
- Verify alerts appear when conditions met
- Verify alerts don't appear when conditions not met
- Test edge cases (exactly at threshold)

### Property 6: Transaction Atomicity (Tính Nguyên Tử Giao Dịch)

**Validates: Requirements 2.2, 2.6, 8.7**

**Mô tả:** Các thao tác tạo batch và sử dụng batch phải là atomic. Nếu bất kỳ bước nào thất bại, toàn bộ thao tác phải rollback, không để lại trạng thái không nhất quán.

**Formal Statement:**
```
∀ batch_creation_transaction:
  (deduct_ingredients AND create_batch_record AND calculate_cost) 
  OR (rollback_all AND return_error)
  
∀ batch_usage_transaction:
  (deduct_batch_quantity AND log_usage AND update_status)
  OR (rollback_all AND return_error)
```

**Test Strategy:**
- Simulate failures at different points in transaction
- Verify complete rollback
- Check database state remains consistent
- Test concurrent transactions
- Verify no partial updates

### Property 7: Quantity Non-Negativity (Không Âm Số Lượng)

**Validates: Requirements 2.5, 5.5**

**Mô tả:** Số lượng nguyên liệu và batch không bao giờ được âm. Hệ thống phải từ chối mọi thao tác dẫn đến số lượng âm.

**Formal Statement:**
```
∀ ingredient, ∀ time:
  ingredient.quantity >= 0

∀ batch_record, ∀ time:
  batch_record.quantity_remaining >= 0

∀ operation:
  IF operation would result in quantity < 0
  THEN operation MUST be rejected
```

**Test Strategy:**
- Attempt to create batch with insufficient ingredients
- Attempt to use more batch than available
- Verify rejection with appropriate error
- Test boundary conditions (exactly zero)
- Test concurrent operations that might cause race conditions

## API Endpoints

### Batch Definition Endpoints

#### POST /api/batch-definitions
Tạo batch definition mới

**Request:**
```json
{
  "name": "Cà Phê Concentrate",
  "unit": "ml",
  "shelf_life_hours": 24,
  "conversion_rates": [
    {
      "source_ingredient_id": "507f1f77bcf86cd799439011",
      "source_quantity": 100,
      "source_unit": "g",
      "batch_quantity": 500,
      "wastage_rate": 0.1
    }
  ],
  "low_stock_threshold": 200,
  "expiry_warning_hours": 4
}
```

**Response (201):**
```json
{
  "id": "507f1f77bcf86cd799439012",
  "name": "Cà Phê Concentrate",
  "unit": "ml",
  "shelf_life_hours": 24,
  "conversion_rates": [...],
  "low_stock_threshold": 200,
  "expiry_warning_hours": 4,
  "created_at": "2026-02-13T10:00:00Z",
  "updated_at": "2026-02-13T10:00:00Z"
}
```

#### GET /api/batch-definitions
Lấy danh sách batch definitions

**Query Parameters:**
- `page` (optional): Số trang (default: 1)
- `limit` (optional): Số lượng mỗi trang (default: 20)
- `search` (optional): Tìm kiếm theo tên

**Response (200):**
```json
{
  "data": [...],
  "total": 10,
  "page": 1,
  "limit": 20
}
```

#### GET /api/batch-definitions/:id
Lấy chi tiết batch definition

**Response (200):**
```json
{
  "id": "507f1f77bcf86cd799439012",
  "name": "Cà Phê Concentrate",
  ...
}
```

#### PUT /api/batch-definitions/:id
Cập nhật batch definition

**Request:** Giống POST

**Response (200):** Giống POST

#### DELETE /api/batch-definitions/:id
Xóa batch definition

**Response (204):** No content

### Batch Record Endpoints

#### POST /api/batch-records
Tạo batch record mới (ghi nhận chế biến batch)

**Request:**
```json
{
  "batch_definition_id": "507f1f77bcf86cd799439012",
  "quantity_produced": 500,
  "prepared_by": "user_id_123"
}
```

**Response (201):**
```json
{
  "id": "507f1f77bcf86cd799439013",
  "batch_definition_id": "507f1f77bcf86cd799439012",
  "batch_name": "Cà Phê Concentrate",
  "quantity_produced": 500,
  "quantity_remaining": 500,
  "unit": "ml",
  "cost_per_unit": 0.15,
  "total_cost": 75.0,
  "prepared_by": "user_id_123",
  "prepared_at": "2026-02-13T10:00:00Z",
  "expires_at": "2026-02-14T10:00:00Z",
  "status": "available",
  "ingredients_used": [
    {
      "ingredient_id": "507f1f77bcf86cd799439011",
      "ingredient_name": "Hạt Cà Phê",
      "quantity": 110,
      "unit": "g",
      "cost_per_unit": 0.68,
      "total_cost": 75.0
    }
  ],
  "created_at": "2026-02-13T10:00:00Z",
  "updated_at": "2026-02-13T10:00:00Z"
}
```

**Error Response (400):**
```json
{
  "error": "insufficient_ingredients",
  "message": "Không đủ Hạt Cà Phê. Cần: 110g, Có: 50g"
}
```

#### GET /api/batch-records
Lấy danh sách batch records

**Query Parameters:**
- `batch_definition_id` (optional): Lọc theo batch definition
- `status` (optional): Lọc theo status (available, expired, depleted)
- `prepared_by` (optional): Lọc theo người chế biến
- `from_date` (optional): Từ ngày
- `to_date` (optional): Đến ngày
- `page` (optional): Số trang
- `limit` (optional): Số lượng mỗi trang

**Response (200):**
```json
{
  "data": [...],
  "total": 50,
  "page": 1,
  "limit": 20
}
```

#### GET /api/batch-records/:id
Lấy chi tiết batch record

**Response (200):** Giống POST response

#### PATCH /api/batch-records/:id/quantity
Cập nhật số lượng batch record

**Request:**
```json
{
  "quantity_remaining": 450
}
```

**Response (200):** Batch record đã cập nhật

#### PATCH /api/batch-records/:id/expire
Đánh dấu batch record là hết hạn

**Response (200):** Batch record với status = "expired"

#### DELETE /api/batch-records/:id
Xóa batch record (chỉ nếu chưa được sử dụng)

**Response (204):** No content

**Error Response (400):**
```json
{
  "error": "batch_already_used",
  "message": "Không thể xóa batch đã được sử dụng"
}
```

### Batch Usage Endpoints

#### POST /api/batch-usage
Sử dụng batch (được gọi từ order system)

**Request:**
```json
{
  "batch_definition_id": "507f1f77bcf86cd799439012",
  "quantity_needed": 30,
  "order_id": "507f1f77bcf86cd799439014",
  "menu_item_id": "507f1f77bcf86cd799439015",
  "menu_item_name": "Cà Phê Đen"
}
```

**Response (200):**
```json
{
  "success": true,
  "batches_used": [
    {
      "batch_record_id": "507f1f77bcf86cd799439013",
      "quantity_used": 30,
      "cost_per_unit": 0.15
    }
  ],
  "total_cost": 4.5,
  "message": "Đã sử dụng batch thành công"
}
```

**Error Response (400):**
```json
{
  "success": false,
  "message": "Không đủ batch khả dụng. Cần: 30ml, Có: 20ml"
}
```

#### GET /api/batch-usage/history
Lấy lịch sử sử dụng batch

**Query Parameters:**
- `batch_record_id` (optional)
- `order_id` (optional)
- `menu_item_id` (optional)
- `from_date` (optional)
- `to_date` (optional)
- `page` (optional)
- `limit` (optional)

**Response (200):**
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439016",
      "batch_record_id": "507f1f77bcf86cd799439013",
      "batch_name": "Cà Phê Concentrate",
      "order_id": "507f1f77bcf86cd799439014",
      "menu_item_id": "507f1f77bcf86cd799439015",
      "menu_item_name": "Cà Phê Đen",
      "quantity_used": 30,
      "unit": "ml",
      "cost_per_unit": 0.15,
      "total_cost": 4.5,
      "used_at": "2026-02-13T11:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "limit": 20
}
```

### Alert Endpoints

#### GET /api/batch-alerts
Lấy tất cả cảnh báo

**Response (200):**
```json
{
  "low_stock": [
    {
      "batch_definition_id": "507f1f77bcf86cd799439012",
      "batch_name": "Cà Phê Concentrate",
      "current_stock": 150,
      "threshold": 200,
      "unit": "ml"
    }
  ],
  "expiring": [
    {
      "batch_record_id": "507f1f77bcf86cd799439013",
      "batch_name": "Cà Phê Concentrate",
      "quantity_remaining": 200,
      "unit": "ml",
      "expires_at": "2026-02-13T14:00:00Z",
      "hours_until_expiry": 3
    }
  ],
  "expired": [
    {
      "batch_record_id": "507f1f77bcf86cd799439017",
      "batch_name": "Trà Đen Concentrate",
      "quantity_wasted": 100,
      "unit": "ml",
      "cost_wasted": 10.0,
      "expired_at": "2026-02-13T09:00:00Z"
    }
  ],
  "last_checked": "2026-02-13T11:00:00Z"
}
```

### Report Endpoints

#### GET /api/batch-reports/production
Báo cáo sản xuất batch

**Query Parameters:**
- `from_date` (required)
- `to_date` (required)
- `batch_definition_id` (optional)
- `prepared_by` (optional)

**Response (200):**
```json
{
  "total_batches_produced": 50,
  "total_quantity_produced": 25000,
  "total_cost": 3750.0,
  "by_batch_type": [
    {
      "batch_name": "Cà Phê Concentrate",
      "count": 30,
      "total_quantity": 15000,
      "total_cost": 2250.0
    }
  ],
  "by_preparer": [
    {
      "preparer": "user_id_123",
      "count": 25,
      "total_quantity": 12500
    }
  ]
}
```

#### GET /api/batch-reports/wastage
Báo cáo lãng phí

**Query Parameters:**
- `from_date` (required)
- `to_date` (required)
- `batch_definition_id` (optional)

**Response (200):**
```json
{
  "total_expired_batches": 5,
  "total_quantity_wasted": 500,
  "total_cost_wasted": 75.0,
  "wastage_by_type": [
    {
      "batch_name": "Cà Phê Concentrate",
      "expired_count": 3,
      "quantity_wasted": 300,
      "cost_wasted": 45.0
    }
  ]
}
```

#### GET /api/batch-reports/usage
Báo cáo sử dụng batch

**Query Parameters:**
- `from_date` (required)
- `to_date` (required)
- `batch_definition_id` (optional)
- `menu_item_id` (optional)

**Response (200):**
```json
{
  "total_usage_count": 200,
  "total_quantity_used": 6000,
  "total_cost": 900.0,
  "by_menu_item": [
    {
      "menu_item_name": "Cà Phê Đen",
      "usage_count": 100,
      "quantity_used": 3000,
      "cost": 450.0
    }
  ],
  "usage_trend": [
    {
      "date": "2026-02-13",
      "quantity_used": 500,
      "cost": 75.0
    }
  ]
}
```


## Frontend Components

### 1. Batch Definition Management

#### BatchDefinitionList.vue
Hiển thị danh sách batch definitions

**Features:**
- Table view với columns: Name, Unit, Shelf Life, Low Stock Threshold, Actions
- Search bar để tìm kiếm theo tên
- Button "Tạo Batch Definition Mới"
- Actions: Edit, Delete, View Details
- Responsive design cho mobile

**Props:**
```typescript
interface Props {
  // No props, fetches data internally
}
```

**Emits:**
```typescript
interface Emits {
  (e: 'edit', id: string): void
  (e: 'delete', id: string): void
}
```

#### BatchDefinitionForm.vue
Form tạo/chỉnh sửa batch definition

**Features:**
- Input fields: Name, Unit, Shelf Life Hours
- Dynamic list of conversion rates với add/remove buttons
- Ingredient selector (autocomplete)
- Input cho source quantity, batch quantity, wastage rate
- Low stock threshold và expiry warning hours
- Validation cho tất cả fields
- Preview chi phí dự kiến

**Props:**
```typescript
interface Props {
  batchDefinition?: BatchDefinition // For edit mode
  mode: 'create' | 'edit'
}
```

**Emits:**
```typescript
interface Emits {
  (e: 'submit', data: BatchDefinitionFormData): void
  (e: 'cancel'): void
}
```

### 2. Batch Record Management

#### BatchRecordList.vue
Hiển thị danh sách batch records

**Features:**
- Table view với columns: Batch Name, Quantity Remaining, Status, Expires At, Prepared By, Actions
- Filter by: Batch Type, Status, Date Range, Preparer
- Sort by: Expiry Date (default), Prepared Date, Quantity
- Color coding: Green (available, plenty), Yellow (low stock), Red (expiring soon), Gray (expired)
- Actions: View Details, Mark as Expired, Delete
- Pagination
- Responsive design

**Props:**
```typescript
interface Props {
  filters?: BatchRecordFilters
}
```

#### BatchRecordForm.vue
Form ghi nhận chế biến batch

**Features:**
- Batch definition selector (dropdown)
- Quantity produced input
- Auto-display: Required ingredients, Expected cost
- Confirmation dialog với cost breakdown
- Success message với batch details
- Error handling cho insufficient ingredients

**Props:**
```typescript
interface Props {
  // No props
}
```

**Emits:**
```typescript
interface Emits {
  (e: 'success', batchRecord: BatchRecord): void
  (e: 'cancel'): void
}
```

#### BatchRecordDetail.vue
Hiển thị chi tiết batch record

**Features:**
- Display all batch record information
- Ingredients used breakdown
- Cost breakdown
- Usage history (list of orders that used this batch)
- Timeline: Prepared → Used → Expired
- Actions: Mark as Expired, Delete (if not used)

**Props:**
```typescript
interface Props {
  batchRecordId: string
}
```

### 3. Alert Components

#### BatchAlertPanel.vue
Hiển thị tất cả cảnh báo batch

**Features:**
- Three sections: Low Stock, Expiring Soon, Expired
- Badge với số lượng cảnh báo
- Expandable sections
- Click to view details
- Auto-refresh every 5 minutes
- Sound/notification khi có cảnh báo mới (optional)

**Props:**
```typescript
interface Props {
  autoRefresh?: boolean // default: true
  refreshInterval?: number // default: 300000 (5 minutes)
}
```

#### BatchAlertCard.vue
Card hiển thị một cảnh báo

**Features:**
- Icon và color theo loại cảnh báo
- Batch name và thông tin chính
- Action button: "Chế Biến Thêm" (low stock), "Sử Dụng Ngay" (expiring)
- Dismiss button (for expired alerts)

**Props:**
```typescript
interface Props {
  alert: LowStockAlert | ExpiringAlert | ExpiredAlert
  type: 'low_stock' | 'expiring' | 'expired'
}
```

### 4. Report Components

#### BatchProductionReport.vue
Báo cáo sản xuất batch

**Features:**
- Date range picker
- Filter by batch type, preparer
- Summary cards: Total Batches, Total Quantity, Total Cost
- Chart: Production trend over time
- Table: Breakdown by batch type
- Export to CSV

**Props:**
```typescript
interface Props {
  // No props
}
```

#### BatchWastageReport.vue
Báo cáo lãng phí

**Features:**
- Date range picker
- Filter by batch type
- Summary cards: Total Expired, Quantity Wasted, Cost Wasted
- Chart: Wastage trend over time
- Table: Breakdown by batch type
- Recommendations để giảm lãng phí

**Props:**
```typescript
interface Props {
  // No props
}
```

#### BatchUsageReport.vue
Báo cáo sử dụng batch

**Features:**
- Date range picker
- Filter by batch type, menu item
- Summary cards: Total Usage, Total Quantity, Total Cost
- Chart: Usage trend over time
- Table: Breakdown by menu item
- Most used batches ranking

**Props:**
```typescript
interface Props {
  // No props
}
```

### 5. Integration Components

#### MenuRecipeEditor.vue (Enhancement)
Thêm khả năng chọn batch trong recipe

**New Features:**
- Toggle để chọn giữa "Nguyên Liệu Thô" và "Batch"
- Batch selector khi chọn "Batch"
- Display available batch quantity
- Warning nếu batch không đủ
- Cost calculation từ batch

### 6. Dashboard Widget

#### BatchStatusWidget.vue
Widget hiển thị trạng thái batch trên dashboard

**Features:**
- Summary: Total Batches, Available Quantity
- Alert count badges
- Quick links: Create Batch, View Alerts, View Reports
- Mini chart: Usage trend (last 7 days)
- Compact design cho dashboard

**Props:**
```typescript
interface Props {
  compact?: boolean // default: false
}
```

## State Management (Pinia Stores)

### useBatchDefinitionStore
```typescript
interface BatchDefinitionStore {
  definitions: BatchDefinition[]
  loading: boolean
  error: string | null
  
  // Actions
  fetchDefinitions(): Promise<void>
  fetchDefinitionById(id: string): Promise<BatchDefinition>
  createDefinition(data: BatchDefinitionFormData): Promise<BatchDefinition>
  updateDefinition(id: string, data: BatchDefinitionFormData): Promise<BatchDefinition>
  deleteDefinition(id: string): Promise<void>
}
```

### useBatchRecordStore
```typescript
interface BatchRecordStore {
  records: BatchRecord[]
  currentRecord: BatchRecord | null
  loading: boolean
  error: string | null
  filters: BatchRecordFilters
  pagination: Pagination
  
  // Actions
  fetchRecords(filters?: BatchRecordFilters): Promise<void>
  fetchRecordById(id: string): Promise<BatchRecord>
  createRecord(data: CreateBatchRequest): Promise<BatchRecord>
  updateQuantity(id: string, quantity: number): Promise<BatchRecord>
  markAsExpired(id: string): Promise<BatchRecord>
  deleteRecord(id: string): Promise<void>
  setFilters(filters: BatchRecordFilters): void
}
```

### useBatchAlertStore
```typescript
interface BatchAlertStore {
  alerts: BatchAlerts
  loading: boolean
  lastChecked: Date | null
  
  // Actions
  fetchAlerts(): Promise<void>
  dismissExpiredAlert(id: string): Promise<void>
  startAutoRefresh(interval: number): void
  stopAutoRefresh(): void
}
```

### useBatchReportStore
```typescript
interface BatchReportStore {
  productionReport: ProductionReport | null
  wastageReport: WastageReport | null
  usageReport: UsageReport | null
  loading: boolean
  error: string | null
  
  // Actions
  fetchProductionReport(params: ReportParams): Promise<void>
  fetchWastageReport(params: ReportParams): Promise<void>
  fetchUsageReport(params: ReportParams): Promise<void>
  exportReport(type: ReportType, format: 'csv' | 'pdf'): Promise<void>
}
```

## Routing

```typescript
const batchRoutes = [
  {
    path: '/batch',
    component: BatchLayout,
    meta: { requiresAuth: true, role: ['manager', 'barista'] },
    children: [
      {
        path: '',
        redirect: '/batch/records'
      },
      {
        path: 'definitions',
        name: 'BatchDefinitions',
        component: BatchDefinitionList,
        meta: { role: ['manager'] }
      },
      {
        path: 'definitions/create',
        name: 'CreateBatchDefinition',
        component: BatchDefinitionForm,
        meta: { role: ['manager'] }
      },
      {
        path: 'definitions/:id/edit',
        name: 'EditBatchDefinition',
        component: BatchDefinitionForm,
        meta: { role: ['manager'] }
      },
      {
        path: 'records',
        name: 'BatchRecords',
        component: BatchRecordList
      },
      {
        path: 'records/create',
        name: 'CreateBatchRecord',
        component: BatchRecordForm
      },
      {
        path: 'records/:id',
        name: 'BatchRecordDetail',
        component: BatchRecordDetail
      },
      {
        path: 'alerts',
        name: 'BatchAlerts',
        component: BatchAlertPanel
      },
      {
        path: 'reports',
        component: BatchReportLayout,
        meta: { role: ['manager'] },
        children: [
          {
            path: 'production',
            name: 'BatchProductionReport',
            component: BatchProductionReport
          },
          {
            path: 'wastage',
            name: 'BatchWastageReport',
            component: BatchWastageReport
          },
          {
            path: 'usage',
            name: 'BatchUsageReport',
            component: BatchUsageReport
          }
        ]
      }
    ]
  }
]
```

## Implementation Considerations

### 1. Concurrency Control

**Problem:** Multiple users có thể tạo batch hoặc sử dụng batch đồng thời, dẫn đến race conditions.

**Solution:**
- Sử dụng MongoDB transactions cho tất cả operations liên quan đến inventory
- Implement optimistic locking với version field
- Retry logic cho failed transactions
- Lock timeout để tránh deadlock

**Example:**
```go
func (s *BatchRecordService) CreateBatch(ctx context.Context, req CreateBatchRequest) (*BatchRecord, error) {
    session, err := s.mongoClient.StartSession()
    if err != nil {
        return nil, err
    }
    defer session.EndSession(ctx)
    
    var batchRecord *BatchRecord
    _, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
        // 1. Check ingredient availability
        // 2. Deduct ingredients
        // 3. Create batch record
        // 4. Calculate cost
        return nil, nil
    })
    
    return batchRecord, err
}
```

### 2. Cost Calculation Caching

**Problem:** Tính toán chi phí batch yêu cầu query nhiều ingredients, có thể chậm.

**Solution:**
- Cache ingredient costs trong memory (Redis hoặc in-memory cache)
- TTL: 5 minutes
- Invalidate cache khi ingredient cost thay đổi
- Fallback to database nếu cache miss

**Example:**
```go
func (c *BatchCostCalculator) GetIngredientCost(ctx context.Context, id primitive.ObjectID) (float64, error) {
    // Try cache first
    if cost, found := c.cache.Get(id.Hex()); found {
        return cost.(float64), nil
    }
    
    // Cache miss, query database
    cost, err := c.ingredientRepo.GetCost(ctx, id)
    if err != nil {
        return 0, err
    }
    
    // Store in cache
    c.cache.Set(id.Hex(), cost, 5*time.Minute)
    return cost, nil
}
```

### 3. FIFO Implementation

**Problem:** Cần đảm bảo batch cũ nhất được sử dụng trước.

**Solution:**
- Query batch records sorted by expires_at ASC
- Filter chỉ lấy available batches (status = "available" AND expires_at > now)
- Sử dụng từng batch theo thứ tự cho đến khi đủ quantity
- Update quantity_remaining và status trong transaction

**Example:**
```go
func (s *BatchUsageService) UseBatch(ctx context.Context, req UseBatchRequest) (*BatchUsageResult, error) {
    // Get available batches sorted by expiry (FIFO)
    batches, err := s.batchRepo.FindAvailableByDefinition(ctx, req.BatchDefinitionID)
    if err != nil {
        return nil, err
    }
    
    remainingNeeded := req.QuantityNeeded
    var batchesUsed []BatchUsageDetail
    
    for _, batch := range batches {
        if remainingNeeded <= 0 {
            break
        }
        
        quantityToUse := math.Min(batch.QuantityRemaining, remainingNeeded)
        
        // Deduct from batch
        err = s.batchRepo.UpdateQuantity(ctx, batch.ID, batch.QuantityRemaining - quantityToUse)
        if err != nil {
            return nil, err
        }
        
        // Log usage
        err = s.logUsage(ctx, batch.ID, quantityToUse, req)
        if err != nil {
            return nil, err
        }
        
        batchesUsed = append(batchesUsed, BatchUsageDetail{
            BatchRecordID: batch.ID,
            QuantityUsed: quantityToUse,
            CostPerUnit: batch.CostPerUnit,
        })
        
        remainingNeeded -= quantityToUse
    }
    
    if remainingNeeded > 0 {
        return &BatchUsageResult{
            Success: false,
            Message: fmt.Sprintf("Không đủ batch. Cần thêm: %.2f%s", remainingNeeded, batches[0].Unit),
        }, nil
    }
    
    return &BatchUsageResult{
        Success: true,
        BatchesUsed: batchesUsed,
        TotalCost: calculateTotalCost(batchesUsed),
    }, nil
}
```

### 4. Alert Checking Strategy

**Problem:** Kiểm tra alerts cho mọi batch mỗi request có thể chậm.

**Solution:**
- Background job chạy mỗi 5 phút để kiểm tra alerts
- Cache alerts trong memory hoặc Redis
- Frontend poll alerts endpoint mỗi 5 phút
- Optional: WebSocket để push alerts real-time

**Example:**
```go
func (s *BatchAlertService) StartAlertChecker(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            alerts, err := s.checkAllAlerts(ctx)
            if err != nil {
                log.Printf("Error checking alerts: %v", err)
                continue
            }
            
            // Cache alerts
            s.cache.Set("batch_alerts", alerts, 5*time.Minute)
            
            // Optional: Push to WebSocket clients
            s.notifyClients(alerts)
            
        case <-ctx.Done():
            return
        }
    }
}
```

### 5. Expiry Handling

**Problem:** Batch có thể hết hạn bất cứ lúc nào, cần cơ chế tự động.

**Solution:**
- Background job chạy mỗi giờ để mark expired batches
- Check expiry trước mỗi usage operation
- Index trên expires_at để query nhanh
- Soft delete (mark as expired) thay vì hard delete

**Example:**
```go
func (s *BatchRecordService) MarkExpiredBatches(ctx context.Context) error {
    filter := bson.M{
        "status": "available",
        "expires_at": bson.M{"$lte": time.Now()},
    }
    
    update := bson.M{
        "$set": bson.M{
            "status": "expired",
            "updated_at": time.Now(),
        },
    }
    
    _, err := s.batchRepo.UpdateMany(ctx, filter, update)
    return err
}
```

### 6. Integration với Menu System

**Problem:** Menu recipes cần hỗ trợ cả ingredients và batches.

**Solution:**
- Extend Recipe schema để hỗ trợ ingredient_type: "raw" | "batch"
- Khi ingredient_type = "batch", ingredient_id trỏ đến batch_definition_id
- Cost calculation check type và query từ đúng source
- Order processing check type và deduct từ đúng inventory

**Example Recipe Schema:**
```go
type RecipeIngredient struct {
    IngredientType  string              `bson:"ingredient_type" json:"ingredient_type"` // "raw" or "batch"
    IngredientID    primitive.ObjectID  `bson:"ingredient_id" json:"ingredient_id"`
    IngredientName  string              `bson:"ingredient_name" json:"ingredient_name"`
    Quantity        float64             `bson:"quantity" json:"quantity"`
    Unit            string              `bson:"unit" json:"unit"`
}
```

### 7. Testing Strategy

**Unit Tests:**
- Test mỗi service method với mocked dependencies
- Test cost calculation với different scenarios
- Test FIFO logic với multiple batches
- Test expiry checking logic
- Test alert generation logic

**Integration Tests:**
- Test full batch creation flow với real database
- Test batch usage flow với order integration
- Test concurrent batch creation
- Test concurrent batch usage
- Test transaction rollback scenarios

**Property-Based Tests:**
- Test inventory conservation property
- Test cost accuracy property
- Test FIFO ordering property
- Test expiry enforcement property
- Test alert correctness property
- Test transaction atomicity property
- Test quantity non-negativity property

**E2E Tests:**
- Test complete user flow: Create definition → Create batch → Use in order
- Test alert flow: Create batch → Wait for expiry → Check alert
- Test report generation với real data

### 8. Performance Optimization

**Database Indexes:**
```javascript
// batch_definitions
db.batch_definitions.createIndex({ name: 1 })
db.batch_definitions.createIndex({ created_at: -1 })

// batch_records
db.batch_records.createIndex({ batch_definition_id: 1, expires_at: 1 })
db.batch_records.createIndex({ status: 1, expires_at: 1 })
db.batch_records.createIndex({ expires_at: 1 })
db.batch_records.createIndex({ prepared_at: -1 })
db.batch_records.createIndex({ prepared_by: 1, prepared_at: -1 })

// batch_usage_logs
db.batch_usage_logs.createIndex({ batch_record_id: 1, used_at: -1 })
db.batch_usage_logs.createIndex({ order_id: 1 })
db.batch_usage_logs.createIndex({ menu_item_id: 1, used_at: -1 })
db.batch_usage_logs.createIndex({ used_at: -1 })
```

**Query Optimization:**
- Use projection để chỉ lấy fields cần thiết
- Limit results với pagination
- Use aggregation pipeline cho complex reports
- Cache frequently accessed data

**Frontend Optimization:**
- Lazy load components
- Virtual scrolling cho long lists
- Debounce search inputs
- Cache API responses trong Pinia store
- Optimistic UI updates

## Security Considerations

### 1. Authentication & Authorization
- Tất cả endpoints yêu cầu authentication
- Role-based access control:
  - Manager: Full access
  - Barista: Create batch records, view alerts, view reports (read-only)
  - Waiter: No access
- JWT token validation cho mọi request

### 2. Input Validation
- Validate tất cả input data trên backend
- Sanitize user inputs để tránh injection attacks
- Validate số lượng > 0
- Validate dates (expiry > now)
- Validate ingredient IDs tồn tại

### 3. Data Integrity
- Use transactions để đảm bảo atomicity
- Validate foreign keys (ingredient_id, batch_definition_id)
- Prevent negative quantities
- Audit log cho tất cả thay đổi quan trọng

### 4. Rate Limiting
- Limit số lượng batch creation requests per user per minute
- Limit API calls để tránh abuse
- Implement exponential backoff cho failed requests

## Monitoring & Logging

### 1. Metrics to Track
- Batch creation rate (per hour/day)
- Batch usage rate
- Wastage rate (expired batches)
- Average batch shelf life utilization
- API response times
- Error rates

### 2. Logging
- Log tất cả batch creation với details
- Log tất cả batch usage
- Log tất cả errors với stack trace
- Log transaction failures
- Log concurrent access conflicts

### 3. Alerts
- Alert khi wastage rate > threshold (e.g., 10%)
- Alert khi API error rate > threshold
- Alert khi database connection issues
- Alert khi background jobs fail

## Migration Strategy

### Phase 1: Backend Implementation (Week 1-2)
1. Implement domain entities
2. Implement repositories
3. Implement services
4. Implement HTTP handlers
5. Write unit tests
6. Write integration tests

### Phase 2: Frontend Implementation (Week 2-3)
1. Implement Pinia stores
2. Implement components
3. Implement routing
4. Integrate with backend APIs
5. Write component tests

### Phase 3: Integration (Week 3-4)
1. Integrate với menu system
2. Integrate với order system
3. Integrate với cost calculator
4. Test end-to-end flows
5. Write property-based tests

### Phase 4: Testing & Refinement (Week 4)
1. User acceptance testing
2. Performance testing
3. Security testing
4. Bug fixes
5. Documentation

### Phase 5: Deployment (Week 5)
1. Deploy to staging
2. Final testing
3. Deploy to production
4. Monitor metrics
5. Gather feedback

## Future Enhancements

### 1. Batch Templates
- Save common batch recipes as templates
- Quick create từ template
- Template versioning

### 2. Predictive Analytics
- Predict batch demand based on historical data
- Suggest optimal batch quantities
- Forecast wastage

### 3. Mobile App
- Dedicated mobile app cho baristas
- Quick batch creation
- Push notifications cho alerts
- Barcode scanning cho ingredients

### 4. Multi-Location Support
- Track batches per location
- Transfer batches between locations
- Location-specific reports

### 5. Quality Control
- Add quality rating cho batches
- Track quality issues
- Quality-based FIFO (use lower quality first)

### 6. Batch Variants
- Support multiple variants của same batch (e.g., strong vs. mild coffee concentrate)
- Different conversion rates per variant
- Variant-specific pricing

## Conclusion

Hệ thống Quản Lý Nguyên Liệu Batch cung cấp giải pháp toàn diện để theo dõi và quản lý nguyên liệu trung gian từ chế biến đến sử dụng. Với thiết kế tập trung vào tính toàn vẹn dữ liệu, hiệu suất, và khả năng mở rộng, hệ thống sẽ giúp giảm lãng phí, tối ưu chi phí, và cải thiện quy trình vận hành.

Các correctness properties được định nghĩa rõ ràng sẽ được kiểm tra thông qua property-based testing để đảm bảo hệ thống hoạt động đúng trong mọi trường hợp. Integration với các hệ thống hiện có được thiết kế cẩn thận để đảm bảo tính nhất quán và dễ bảo trì.
