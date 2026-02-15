# Backend Batch Logic Analysis

## Câu hỏi
Backend có đang đáp ứng logic cho phép chọn batch definition trực tiếp (không cần batch record) khi tạo menu?

## Trả lời: CÓ và KHÔNG

### Backend HỖ TRỢ lưu batch definition trong menu
✅ **CÓ** - Backend cho phép lưu batch definition ID trong menu:

```go
type Ingredient struct {
    Name           string               `bson:"name" json:"name"`
    Quantity       float64              `bson:"quantity" json:"quantity"`
    Unit           ingredient.UnitType  `bson:"unit" json:"unit"`
    
    // Batch support
    IngredientType string               `bson:"ingredient_type" json:"ingredient_type"` // "raw" or "batch"
    BatchID        *primitive.ObjectID  `bson:"batch_id,omitempty" json:"batch_id,omitempty"` // Reference to batch_definition_id
}
```

- Menu có thể lưu `batch_definition_id` trong field `BatchID`
- Không cần có batch record để **TẠO** menu
- Frontend có thể chọn batch definition và lưu vào menu

### Backend YÊU CẦU batch records để tính cost
❌ **KHÔNG** - Backend cần batch records để tính cost:

```go
func (s *CostCalculatorService) getBatchCostPerUnit(ctx context.Context, batchDefID primitive.ObjectID) (float64, error) {
    // Get available batches (sorted by expiry date - FIFO)
    batches, err := s.batchRecRepo.FindAvailableByDefinition(ctx, batchDefID)
    if err != nil {
        return 0, fmt.Errorf("failed to find available batches: %w", err)
    }
    
    // If no batches available, return error
    if len(batches) == 0 {
        return 0, fmt.Errorf("no available batches")
    }
    
    // Return cost per unit from the first (oldest) batch (FIFO)
    return batches[0].CostPerUnit, nil
}
```

Logic:
1. Tìm batch records có sẵn từ batch definition
2. Nếu không có batch record → lỗi "no available batches"
3. Lấy cost từ batch record cũ nhất (FIFO)

## Workflow hiện tại

### Tạo Menu (✅ Hoạt động)
```
1. User chọn batch definition
2. Frontend gửi: { ingredient_type: "batch", batch_id: "def_123" }
3. Backend lưu vào menu
4. ✅ Menu được tạo thành công
```

### Tính Cost (❌ Cần batch records)
```
1. Backend gọi CalculateMenuItemCost()
2. Tìm batch records từ batch_definition_id
3. Nếu không có batch records:
   ❌ Lỗi: "no available batches"
   ❌ Cost = incomplete
4. Nếu có batch records:
   ✅ Lấy cost từ batch record cũ nhất
   ✅ Cost được tính
```

## Tại sao thiết kế như vậy?

### Lý do hợp lý
1. **Cost thực tế**: Batch cost phụ thuộc vào giá nguyên liệu tại thời điểm sản xuất
2. **FIFO**: Sử dụng batch cũ nhất trước (First In First Out)
3. **Tồn kho**: Chỉ tính cost khi có batch thực sự available
4. **Chính xác**: Cost phản ánh giá trị thực của batch đã sản xuất

### Ví dụ
```
Batch Definition: "Cà phê Concentrate"
- Công thức: 100g cà phê hạt → 500ml concentrate

Batch Record #1 (sản xuất 01/02):
- Giá cà phê: 200,000 VNĐ/kg
- Cost per unit: 400 VNĐ/ml

Batch Record #2 (sản xuất 10/02):
- Giá cà phê: 250,000 VNĐ/kg (tăng giá)
- Cost per unit: 500 VNĐ/ml

Menu Item: "Cà phê sữa đá"
- Sử dụng: 50ml Cà phê Concentrate
- Cost = 50ml × 400 VNĐ/ml = 20,000 VNĐ (từ batch #1 - FIFO)
```

## Giải pháp

### Option 1: Giữ nguyên logic (Khuyến nghị)
**Ưu điểm:**
- ✅ Cost chính xác dựa trên batch thực tế
- ✅ Phản ánh tồn kho thực
- ✅ FIFO logic đúng với thực tế

**Nhược điểm:**
- ❌ Phải sản xuất batch trước khi có cost
- ❌ Menu mới không có cost ngay

**Workflow:**
```
1. Tạo batch definition
2. Tạo menu (chọn batch definition)
3. Sản xuất batch record
4. Cost được tính tự động
```

### Option 2: Tính cost ước tính từ definition
**Ưu điểm:**
- ✅ Menu có cost ngay lập tức
- ✅ Không cần sản xuất batch trước

**Nhược điểm:**
- ❌ Cost không chính xác (ước tính)
- ❌ Không phản ánh giá thực tế
- ❌ Cần thêm logic phức tạp

**Cần implement:**
```go
func (s *CostCalculatorService) estimateBatchCostFromDefinition(ctx context.Context, batchDefID primitive.ObjectID) (float64, error) {
    // 1. Get batch definition
    // 2. Get conversion rates
    // 3. Calculate cost from source ingredients
    // 4. Return estimated cost
}
```

### Option 3: Hybrid approach
**Logic:**
- Nếu có batch records → dùng cost thực tế (FIFO)
- Nếu không có → tính cost ước tính từ definition

**Ưu điểm:**
- ✅ Linh hoạt
- ✅ Menu luôn có cost

**Nhược điểm:**
- ❌ Phức tạp hơn
- ❌ Cost có thể không nhất quán

## Kết luận

### Backend hiện tại
✅ **Hỗ trợ** lưu batch definition trong menu
❌ **Yêu cầu** batch records để tính cost

### Khuyến nghị
**Giữ nguyên logic hiện tại** vì:
1. Cost chính xác hơn
2. Phản ánh thực tế kinh doanh
3. FIFO logic đúng
4. Đơn giản hơn

### Hướng dẫn sử dụng
```
Bước 1: Tạo Batch Definition
- Định nghĩa công thức batch
- Lưu vào hệ thống

Bước 2: Tạo Menu
- Chọn batch definition
- Menu được lưu (cost = incomplete)

Bước 3: Sản xuất Batch
- Tạo batch record từ definition
- Cost được tính tự động

Bước 4: Bán hàng
- Menu có cost chính xác
- Profit được tính đúng
```

## Files liên quan

### Backend
- `backend/domain/menu/menu.go` - Menu ingredient structure
- `backend/application/services/cost_calculator_service.go` - Cost calculation logic
- `backend/application/services/batch_cost_calculator.go` - Batch cost calculator

### Frontend
- `frontend/src/views/MenuView.vue` - Menu creation UI
- `frontend/src/stores/batchDefinition.js` - Batch definition store
- `frontend/src/stores/batchRecord.js` - Batch record store

## Status

✅ **Backend đã hỗ trợ đầy đủ** - Logic hiện tại là đúng và hợp lý
📝 **Không cần thay đổi** - Workflow yêu cầu sản xuất batch trước khi có cost chính xác
