# Menu Cost Auto-Recalculation Implementation

## 🎯 Mục tiêu

Implement 2 tính năng để cập nhật chi phí món tự động và thủ công:

1. **Option 2**: Auto-recalculate khi tạo batch mới
2. **Option 3**: Nút "Tính lại tất cả" trong UI

## 📊 Vấn đề ban đầu

### Chi phí KHÔNG được tính real-time

```
User mở Menu Costs view
    ↓
GET /api/manager/menu/costs
    ↓
ĐỌC menu_items.current_cost từ database
    ↓
KHÔNG tính lại chi phí
```

**Kết quả**: Khi tạo batch mới với chi phí khác, menu cost vẫn hiển thị giá trị cũ.

### Ví dụ thực tế:

```
10:00 - Tạo batch "cfe cot": 36 đ/ml
        Tính cost món "áddd": 7,200 đ
        → menu_items.current_cost = 7,200 đ

14:00 - Giá cà phê tăng
        Tạo batch mới: 45 đ/ml

15:00 - User mở Menu Costs view
        → Vẫn hiển thị: 7,200 đ ❌
        → Thực tế nếu tính lại: 9,000 đ
```

---

## ✅ OPTION 2: Auto-recalculate khi tạo batch

### Backend Changes

#### 1. Thêm method mới vào MenuRepository

**File**: `backend/application/services/menu.go`

```go
type MenuRepository interface {
    // ... existing methods ...
    FindByBatchDefinitionID(ctx context.Context, batchDefID primitive.ObjectID) ([]*menu.MenuItem, error)
}
```

#### 2. Implement trong MongoDB Repository

**File**: `backend/infrastructure/mongodb/menu_repository.go`

```go
// FindByBatchDefinitionID finds all menu items that use a specific batch definition
// Searches in both single-size ingredients and multi-size variant ingredients
func (r *MenuRepository) FindByBatchDefinitionID(ctx context.Context, batchDefID primitive.ObjectID) ([]*menu.MenuItem, error) {
    filter := bson.M{
        "$or": []bson.M{
            // Single-size items with batch ingredient
            {"ingredients.batch_id": batchDefID},
            // Multi-size items with batch ingredient in variants
            {"variants.ingredients.batch_id": batchDefID},
        },
    }
    
    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var items []*menu.MenuItem
    if err = cursor.All(ctx, &items); err != nil {
        return nil, err
    }
    return items, nil
}
```

#### 3. Update BatchRecordService

**File**: `backend/application/services/batch_record_service.go`

**Thêm dependencies**:
```go
type BatchRecordService struct {
    // ... existing fields ...
    costRecalcService *CostRecalculationService
    menuRepo          MenuRepository
}

// Setter methods
func (s *BatchRecordService) SetCostRecalculationService(costRecalcService *CostRecalculationService) {
    s.costRecalcService = costRecalcService
}

func (s *BatchRecordService) SetMenuRepository(menuRepo MenuRepository) {
    s.menuRepo = menuRepo
}
```

**Thêm logic trong CreateBatch**:
```go
func (s *BatchRecordService) CreateBatch(...) (*batch.BatchRecord, error) {
    // ... existing batch creation logic ...
    
    // After successful batch creation, queue cost recalculation
    if s.costRecalcService != nil && s.menuRepo != nil {
        go s.queueMenuCostRecalculation(context.Background(), req.BatchDefinitionID)
    }

    return batchRecord, nil
}

// queueMenuCostRecalculation finds all menu items using this batch and queues cost recalculation
func (s *BatchRecordService) queueMenuCostRecalculation(ctx context.Context, batchDefID primitive.ObjectID) {
    // Find all menu items that use this batch definition
    menuItems, err := s.menuRepo.FindByBatchDefinitionID(ctx, batchDefID)
    if err != nil {
        fmt.Printf("Warning: failed to find menu items for batch definition %s: %v\n", batchDefID.Hex(), err)
        return
    }

    // Queue recalculation for each menu item
    for _, menuItem := range menuItems {
        err := s.costRecalcService.QueueRecalculation(menuItem.ID)
        if err != nil {
            fmt.Printf("Warning: failed to queue recalculation for menu item %s: %v\n", menuItem.ID.Hex(), err)
        }
    }

    if len(menuItems) > 0 {
        fmt.Printf("Info: queued cost recalculation for %d menu items using batch definition %s\n", len(menuItems), batchDefID.Hex())
    }
}
```

#### 4. Wire dependencies trong main.go

**File**: `backend/main.go`

```go
// Wire up batch record service with cost recalculation and menu repo
// This enables auto-recalculation of menu costs when new batches are created
batchRecordService.SetCostRecalculationService(costRecalculationService)
batchRecordService.SetMenuRepository(menuRepo)
```

### Luồng hoạt động

```
Barista tạo batch mới
    ↓
BatchRecordService.CreateBatch()
    ↓
1. Tính cost batch (BatchCostCalculator)
2. Lưu batch record với cost_per_unit
3. Commit transaction
    ↓
4. Trigger background job (goroutine)
    ↓
queueMenuCostRecalculation()
    ↓
5. Tìm tất cả menu items dùng batch này
   menuRepo.FindByBatchDefinitionID()
    ↓
6. Queue recalculation cho từng món
   costRecalcService.QueueRecalculation()
    ↓
7. Background workers xử lý queue
   CostCalculatorService.CalculateMenuItemCost()
    ↓
8. Update menu_items.current_cost
```

### Ưu điểm

- ✅ Tự động, không cần user action
- ✅ Chạy background, không block batch creation
- ✅ Chỉ tính lại món bị ảnh hưởng (hiệu quả)
- ✅ Sử dụng worker pool có sẵn

---

## ✅ OPTION 3: Nút "Tính lại tất cả"

### Backend Changes

#### 1. Thêm API endpoint

**File**: `backend/main.go`

```go
// Menu cost and profit analysis routes
manager.GET("/menu/costs", menuCostHandler.GetMenuCosts)
manager.GET("/menu/costs/:id", menuCostHandler.GetMenuCostDetail)
manager.GET("/menu/warnings", menuCostHandler.GetMenuWarnings)
manager.POST("/menu/costs/recalculate-all", menuCostHandler.RecalculateAllCosts)  // NEW
```

#### 2. Update MenuCostHandler

**File**: `backend/interfaces/http/menu_cost_handler.go`

**Thêm MenuRepository dependency**:
```go
type MenuCostHandler struct {
    profitAnalyzer       *services.ProfitAnalyzerService
    costCalculator       *services.CostCalculatorService
    recalculationService *services.CostRecalculationService
    menuRepo             services.MenuRepository  // NEW
}

func NewMenuCostHandler(
    profitAnalyzer *services.ProfitAnalyzerService,
    costCalculator *services.CostCalculatorService,
    recalculationService *services.CostRecalculationService,
    menuRepo services.MenuRepository,  // NEW
) *MenuCostHandler {
    return &MenuCostHandler{
        profitAnalyzer:       profitAnalyzer,
        costCalculator:       costCalculator,
        recalculationService: recalculationService,
        menuRepo:             menuRepo,
    }
}
```

**Thêm handler method**:
```go
// RecalculateAllCosts handles POST /api/menu/costs/recalculate-all
// Queues cost recalculation for all menu items
func (h *MenuCostHandler) RecalculateAllCosts(c *gin.Context) {
    // Get all menu items
    menuItems, err := h.menuRepo.FindAll(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch menu items"})
        return
    }

    // Queue recalculation for each menu item
    queuedCount := 0
    failedCount := 0
    
    for _, menuItem := range menuItems {
        err := h.recalculationService.QueueRecalculation(menuItem.ID)
        if err != nil {
            failedCount++
        } else {
            queuedCount++
        }
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message":      "cost recalculation queued for all menu items",
        "total_items":  len(menuItems),
        "queued_count": queuedCount,
        "failed_count": failedCount,
    })
}
```

#### 3. Wire dependency trong main.go

```go
// Menu cost and profit analysis handlers
menuCostHandler := http.NewMenuCostHandler(profitAnalyzerService, costCalculatorService, costRecalculationService, menuRepo)
```

### Frontend Changes

#### 1. Thêm service method

**File**: `frontend/src/services/menuCost.js`

```javascript
/**
 * Recalculate costs for all menu items
 * Queues background jobs to recalculate costs based on current ingredient/batch prices
 * @returns {Promise<{message: string, total_items: number, queued_count: number, failed_count: number}>}
 */
async recalculateAllCosts() {
  const response = await api.post('/manager/menu/costs/recalculate-all')
  return response.data
}
```

#### 2. Thêm UI button

**File**: `frontend/src/views/MenuCostView.vue`

**Template**:
```vue
<div class="flex items-center justify-between mb-3">
  <h1 class="text-xl md:text-2xl font-bold text-gray-800">💰 Chi phí món</h1>
  
  <!-- Desktop: View Toggle + Recalculate Button -->
  <div class="hidden md:flex gap-2">
    <button @click="recalculateAllCosts" 
      :disabled="isRecalculating"
      :class="isRecalculating ? 'bg-gray-300 cursor-not-allowed' : 'bg-green-500 hover:bg-green-600'"
      class="px-3 py-1 rounded-lg text-sm font-medium text-white flex items-center gap-1">
      <span v-if="isRecalculating" class="animate-spin">⏳</span>
      <span v-else>🔄</span>
      <span>{{ isRecalculating ? 'Đang tính...' : 'Tính lại tất cả' }}</span>
    </button>
    <!-- ... view toggle buttons ... -->
  </div>
  
  <!-- Mobile: Recalculate Button Only -->
  <div class="md:hidden">
    <button @click="recalculateAllCosts" 
      :disabled="isRecalculating"
      class="px-3 py-1 rounded-lg text-sm font-medium text-white bg-green-500">
      <span v-if="isRecalculating" class="animate-spin">⏳</span>
      <span v-else>🔄</span>
    </button>
  </div>
</div>
```

**Script**:
```javascript
// State
const isRecalculating = ref(false)

// Recalculate all costs function
const recalculateAllCosts = async () => {
  if (isRecalculating.value) return
  
  isRecalculating.value = true
  
  try {
    const response = await menuCostService.recalculateAllCosts()
    
    // Show success message
    alert(`✅ Đã gửi yêu cầu tính lại chi phí cho ${response.total_items} món.\n\nĐang xử lý trong background...`)
    
    // Refresh data after a short delay
    setTimeout(() => {
      fetchData()
    }, 1000)
  } catch (err) {
    console.error('Error recalculating costs:', err)
    alert('❌ Lỗi khi tính lại chi phí: ' + (err.response?.data?.error || err.message))
  } finally {
    isRecalculating.value = false
  }
}
```

### Luồng hoạt động

```
User click "Tính lại tất cả"
    ↓
POST /api/manager/menu/costs/recalculate-all
    ↓
MenuCostHandler.RecalculateAllCosts()
    ↓
1. Lấy tất cả menu items
   menuRepo.FindAll()
    ↓
2. Queue recalculation cho từng món
   costRecalcService.QueueRecalculation()
    ↓
3. Trả về response ngay lập tức
   {total_items: 50, queued_count: 50, failed_count: 0}
    ↓
4. Background workers xử lý queue
   CostCalculatorService.CalculateMenuItemCost()
    ↓
5. Update menu_items.current_cost
    ↓
6. Frontend refresh sau 1 giây
   Hiển thị chi phí mới
```

### Ưu điểm

- ✅ User có control khi nào tính lại
- ✅ Đơn giản, dễ hiểu
- ✅ Không block UI (background processing)
- ✅ Hiển thị progress qua recalculation_status

---

## 🎨 UI/UX Improvements

### Desktop View

```
┌─────────────────────────────────────────────────────────┐
│ 💰 Chi phí món          [🔄 Tính lại tất cả] [📱 Thẻ] [📊 Bảng] │
├─────────────────────────────────────────────────────────┤
│ [Tìm kiếm món...]                                       │
│ [📁 Tất cả] [☕ Coffee] [🍰 Cake] ...                   │
│ Sắp xếp: [Lợi nhuận %] [↓ Giảm]                       │
├─────────────────────────────────────────────────────────┤
│ ┌─────────────────────────────────────────────────────┐ │
│ │ Tổng quan chi phí món                               │ │
│ │ ⏳ Đang cập nhật... (Đã xử lý: 25/50 món)          │ │
│ │ [50] Tổng món  [2] Bán lỗ  [5] LN thấp  [45%] LN TB│ │
│ └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

### Mobile View

```
┌───────────────────────────┐
│ 💰 Chi phí món        [🔄] │
├───────────────────────────┤
│ [Tìm kiếm món...]         │
│ [📁 Tất cả] [☕] [🍰] ... │
├───────────────────────────┤
│ ⏳ Đang cập nhật...       │
│ Đã xử lý: 25/50 món       │
└───────────────────────────┘
```

---

## 📊 Monitoring & Status

### Recalculation Status

View hiển thị status real-time:

```javascript
{
  in_progress: true,
  queued_items: 50,
  processed_items: 25,
  last_updated: "2024-02-24T10:30:00Z"
}
```

**UI hiển thị**:
- ⏳ Icon spinning khi đang xử lý
- Progress: "Đã xử lý: 25/50 món"
- Auto-refresh mỗi vài giây

---

## 🔄 Workflow tổng hợp

### Scenario 1: Tạo batch mới

```
1. Barista tạo batch "cfe cot" mới (cost = 45 đ/ml)
   ↓
2. BatchRecordService tự động queue recalculation
   → Tìm món "áddd" dùng batch này
   → Queue recalculation
   ↓
3. Background worker tính lại
   → current_cost: 7,200 đ → 9,000 đ
   ↓
4. User mở Menu Costs view
   → Hiển thị: 9,000 đ ✅
```

### Scenario 2: User muốn refresh tất cả

```
1. User click "Tính lại tất cả"
   ↓
2. Frontend gọi API
   POST /api/manager/menu/costs/recalculate-all
   ↓
3. Backend queue 50 món
   → Response: {total_items: 50, queued_count: 50}
   ↓
4. Alert: "Đã gửi yêu cầu tính lại 50 món"
   ↓
5. Background workers xử lý
   → Status: in_progress = true
   → processed_items tăng dần
   ↓
6. Frontend auto-refresh
   → Hiển thị progress
   → Sau khi xong: in_progress = false
```

---

## ✅ Testing Checklist

### Backend Tests

- [ ] `FindByBatchDefinitionID` tìm đúng menu items
  - Single-size items với batch ingredient
  - Multi-size items với batch ingredient trong variants
- [ ] `CreateBatch` trigger recalculation
  - Queue được tạo cho menu items đúng
  - Không block batch creation
- [ ] `RecalculateAllCosts` endpoint
  - Queue tất cả menu items
  - Trả về đúng count
  - Handle errors gracefully

### Frontend Tests

- [ ] Button "Tính lại tất cả" hoạt động
  - Disable khi đang xử lý
  - Hiển thị loading state
  - Show success/error message
- [ ] UI responsive
  - Desktop: Full button với text
  - Mobile: Icon only
- [ ] Status indicator
  - Hiển thị khi đang xử lý
  - Update progress real-time

### Integration Tests

- [ ] End-to-end flow
  1. Tạo batch mới
  2. Verify menu cost được update
  3. Click "Tính lại tất cả"
  4. Verify tất cả costs được update

---

## 🎯 Kết quả

### Trước khi implement

❌ Chi phí không được cập nhật khi tạo batch mới  
❌ User phải manually trigger recalculation cho từng món  
❌ Không có cách nào refresh tất cả costs  

### Sau khi implement

✅ Chi phí tự động update khi tạo batch (Option 2)  
✅ User có thể refresh tất cả costs bằng 1 click (Option 3)  
✅ Background processing không block UI  
✅ Status indicator cho user biết progress  
✅ Responsive design cho mobile & desktop  

---

## 📝 Notes

### Performance Considerations

- Background processing không block batch creation
- Worker pool giới hạn concurrent jobs (4 workers)
- Queue size: 1000 items
- Goroutine cho async processing

### Error Handling

- Batch creation không fail nếu recalculation fail
- Log warnings cho failed queue operations
- Frontend hiển thị error messages
- Graceful degradation

### Future Enhancements

- [ ] Progress bar thay vì text counter
- [ ] Notification khi recalculation hoàn thành
- [ ] Selective recalculation (chỉ batch ingredients)
- [ ] Schedule auto-recalculation (daily/weekly)
- [ ] Cost history tracking (daily snapshots)

