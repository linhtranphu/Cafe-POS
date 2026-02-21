# Batch Order Unit Conversion Fix

## Vấn Đề

Khi tạo order với batch ingredient, hệ thống báo lỗi:
```
failed to deduct batch ingredients: insufficient batch cfe cot: 
Insufficient batch quantity. Need: 200.00, Available: 1.00
```

**Nguyên nhân:**
- Batch có 1 lít (unit: "l", quantity: 1)
- Menu item cần 200ml (unit: "ml", quantity: 200)
- Hệ thống so sánh trực tiếp 200 với 1 mà không convert đơn vị
- Thực tế: 1l = 1000ml, đủ để làm 5 ly (200ml mỗi ly)

## Root Cause Analysis

### 1. BatchUsageService.UseBatch()
```go
// File: backend/application/services/batch_usage_service.go
// Line: ~100

// Calculate total available quantity
totalAvailable := 0.0
for _, b := range availableBatches {
    totalAvailable += b.QuantityRemaining  // ❌ Không convert unit
}

if totalAvailable < req.QuantityNeeded {  // ❌ So sánh sai đơn vị
    return &BatchUsageResult{
        Success: false,
        Message: fmt.Sprintf("Insufficient batch quantity. Need: %.2f, Available: %.2f",
            req.QuantityNeeded, totalAvailable),
    }, nil
}
```

### 2. OrderService.deductBatchIngredients()
```go
// File: backend/application/services/order_service.go
// Line: ~590

// Calculate total quantity needed (ingredient quantity * order item quantity)
quantityNeeded := ing.Quantity * float64(item.Quantity)

// Deduct batch using BatchUsageService
req := UseBatchRequest{
    BatchDefinitionID: *ing.BatchID,
    QuantityNeeded:    quantityNeeded,  // ❌ Gửi 200ml nhưng không chỉ định unit
    OrderID:           o.ID,
    MenuItemID:        item.MenuItemID,
    MenuItemName:      item.Name,
}
```

**Vấn đề:** 
- `quantityNeeded` là 200 (ml) từ menu ingredient
- Batch có `QuantityRemaining` là 1 (l)
- Không có thông tin unit trong `UseBatchRequest`
- Không có logic convert giữa ml và l

## Giải Pháp

### Option 1: Normalize tất cả về base unit (Khuyến nghị)

Lưu tất cả số lượng theo base unit (ml cho liquid, g cho solid):

**Ưu điểm:**
- Đơn giản, không cần convert mỗi lần
- Tránh lỗi làm tròn
- Dễ so sánh và tính toán

**Nhược điểm:**
- Cần migration data hiện tại
- UI vẫn hiển thị theo unit người dùng chọn (l, kg, etc.)

### Option 2: Convert unit khi so sánh (Quick fix)

Thêm unit conversion vào `UseBatch`:

**Ưu điểm:**
- Không cần migration
- Giữ nguyên cách lưu hiện tại

**Nhược điểm:**
- Phức tạp hơn
- Cần convert mỗi lần so sánh
- Có thể có lỗi làm tròn

## Implementation Plan

### Phase 1: Quick Fix (Option 2)

1. **Thêm Unit vào UseBatchRequest**
```go
type UseBatchRequest struct {
    BatchDefinitionID primitive.ObjectID
    QuantityNeeded    float64
    Unit              string  // ← Thêm field này
    OrderID           primitive.ObjectID
    MenuItemID        primitive.ObjectID
    MenuItemName      string
}
```

2. **Thêm Unit Conversion Helper**
```go
// File: backend/domain/batch/unit_conversion.go

func ConvertQuantity(quantity float64, fromUnit, toUnit string) (float64, error) {
    if fromUnit == toUnit {
        return quantity, nil
    }
    
    // Volume conversions
    volumeUnits := map[string]float64{
        "ml": 1,
        "l":  1000,
    }
    
    // Weight conversions
    weightUnits := map[string]float64{
        "g":  1,
        "kg": 1000,
    }
    
    // Try volume conversion
    if fromBase, ok := volumeUnits[fromUnit]; ok {
        if toBase, ok := volumeUnits[toUnit]; ok {
            return quantity * fromBase / toBase, nil
        }
    }
    
    // Try weight conversion
    if fromBase, ok := weightUnits[fromUnit]; ok {
        if toBase, ok := weightUnits[toUnit]; ok {
            return quantity * fromBase / toBase, nil
        }
    }
    
    return 0, fmt.Errorf("cannot convert from %s to %s", fromUnit, toUnit)
}
```

3. **Update UseBatch() để convert unit**
```go
func (s *BatchUsageService) UseBatch(ctx context.Context, req UseBatchRequest) (*BatchUsageResult, error) {
    // ... existing code ...
    
    // Get batch unit from first batch
    batchUnit := availableBatches[0].Unit
    
    // Convert quantity needed to batch unit
    quantityNeededInBatchUnit, err := batch.ConvertQuantity(req.QuantityNeeded, req.Unit, batchUnit)
    if err != nil {
        return nil, fmt.Errorf("failed to convert units: %w", err)
    }
    
    // Calculate total available quantity (already in batch unit)
    totalAvailable := 0.0
    for _, b := range availableBatches {
        totalAvailable += b.QuantityRemaining
    }
    
    if totalAvailable < quantityNeededInBatchUnit {
        return &BatchUsageResult{
            Success: false,
            Message: fmt.Sprintf("Insufficient batch quantity. Need: %.2f%s (%.2f%s), Available: %.2f%s",
                req.QuantityNeeded, req.Unit,
                quantityNeededInBatchUnit, batchUnit,
                totalAvailable, batchUnit),
        }, nil
    }
    
    // Use quantityNeededInBatchUnit for deduction
    remainingNeeded := quantityNeededInBatchUnit
    // ... rest of the code ...
}
```

4. **Update OrderService để truyền unit**
```go
func (s *OrderService) deductBatchIngredients(ctx context.Context, o *order.Order) (float64, error) {
    // ... existing code ...
    
    // Deduct batch using BatchUsageService
    req := UseBatchRequest{
        BatchDefinitionID: *ing.BatchID,
        QuantityNeeded:    quantityNeeded,
        Unit:              ing.Unit,  // ← Thêm unit từ ingredient
        OrderID:           o.ID,
        MenuItemID:        item.MenuItemID,
        MenuItemName:      item.Name,
    }
    
    // ... rest of the code ...
}
```

### Phase 2: Long-term Solution (Option 1)

Sau khi Phase 1 stable, có thể migrate sang base unit:

1. Migration script để convert tất cả batch records sang base unit
2. Update UI để convert khi hiển thị
3. Simplify logic trong UseBatch (không cần convert nữa)

## Testing

### Test Case 1: Same Unit
- Batch: 500ml
- Menu: 200ml
- Expected: Success, 300ml remaining

### Test Case 2: Different Unit (l → ml)
- Batch: 1l (1000ml)
- Menu: 200ml
- Expected: Success, 800ml remaining (0.8l)

### Test Case 3: Different Unit (kg → g)
- Batch: 1kg (1000g)
- Menu: 100g
- Expected: Success, 900g remaining (0.9kg)

### Test Case 4: Insufficient After Conversion
- Batch: 0.1l (100ml)
- Menu: 200ml
- Expected: Error "Insufficient batch quantity. Need: 200ml (0.2l), Available: 100ml (0.1l)"

## Files to Modify

1. `backend/domain/batch/unit_conversion.go` (NEW)
2. `backend/application/services/batch_usage_service.go`
3. `backend/application/services/order_service.go`
4. `backend/application/services/batch_usage_service_test.go` (add tests)

## Rollout Plan

1. Implement Phase 1 (unit conversion)
2. Test thoroughly with different unit combinations
3. Deploy to staging
4. Monitor for issues
5. Deploy to production
6. (Optional) Plan Phase 2 migration

## Notes

- Conversion chỉ support volume (ml, l) và weight (g, kg) units
- Nếu cần thêm units khác (oz, lb, etc.), thêm vào conversion map
- Cần handle edge cases: invalid units, incompatible unit types (volume vs weight)
