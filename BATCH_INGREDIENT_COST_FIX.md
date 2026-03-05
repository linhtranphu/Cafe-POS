# Fix: Batch Ingredients không tính được chi phí món

## 🔴 VẤN ĐỀ

Món "áddd" sử dụng nguyên liệu batch "cfe cot" (cà phê cốt) nhưng không tính được chi phí, hiển thị:
- Chi phí: 0 ₫
- Cost Status: ⚠ Thiếu dữ liệu

## 🔍 NGUYÊN NHÂN

Hệ thống có 2 loại nguyên liệu:

### 1. Raw Ingredients (Nguyên liệu thô)
- Lưu trong collection `ingredients`
- Có `cost_per_unit` cố định
- Ví dụ: Cà phê hạt, sữa, đường

### 2. Batch Ingredients (Nguyên liệu batch)
- Lưu trong collection `batch_records`
- Được chế biến từ nguyên liệu thô
- Chi phí được tính từ công thức
- Ví dụ: Cà phê cốt (chế biến từ cà phê hạt + nước)

**Vấn đề**: `CostCalculatorService` có logic xử lý batch ingredients:

```go
// Check if this is a batch ingredient
if menuIngredient.IsBatchIngredient() {
    // Get batch cost
    batchCost, err := s.getBatchCostPerUnit(context.Background(), *menuIngredient.BatchID)
    ...
}
```

Nhưng `batchRecRepo` KHÔNG được wire vào service trong `main.go`:

```go
// TRƯỚC (SAI):
costCalculatorService := services.NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)
// batchRecRepo không được set → s.batchRecRepo == nil
```

Khi `getBatchCostPerUnit` được gọi:

```go
func (s *CostCalculatorService) getBatchCostPerUnit(...) {
    if s.batchRecRepo == nil {
        return 0, fmt.Errorf("batch repository not configured") // ← LỖI Ở ĐÂY
    }
    ...
}
```

## ✅ GIẢI PHÁP

### 1. Tạo Adapter
Tạo `batch_record_cost_adapter.go` để convert giữa 2 interfaces:
- `BatchRecordRepository` (MongoDB layer)
- `CostCalculatorBatchRecordRepository` (Service layer)

```go
type BatchRecordCostAdapter struct {
    repo BatchRecordRepositoryForAdapter
}

func (a *BatchRecordCostAdapter) FindAvailableByDefinition(...) ([]*CostCalculatorBatchRecord, error) {
    records, err := a.repo.FindAvailableByDefinition(ctx, defID)
    // Convert batch.BatchRecord → CostCalculatorBatchRecord
    ...
}
```

### 2. Wire vào main.go
```go
costCalculatorService := services.NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

// Wire up batch repository
batchRecordCostAdapter := services.NewBatchRecordCostAdapter(batchRecordRepo)
costCalculatorService.SetBatchRecordRepository(batchRecordCostAdapter)
```

## 📊 LUỒNG TÍNH CHI PHÍ BATCH

```
Menu Item "áddd"
  ↓
  Ingredient: "cfe cot" (type: batch)
  ↓
  CostCalculatorService.calculateIngredientsCost()
  ↓
  Kiểm tra: IsBatchIngredient() → true
  ↓
  getBatchCostPerUnit(batchDefID)
  ↓
  batchRecRepo.FindAvailableByDefinition() → Tìm batch records
  ↓
  Lấy batch cũ nhất (FIFO)
  ↓
  Trả về: batch.CostPerUnit
  ↓
  Tính: quantity × costPerUnit
  ↓
  Cộng vào tổng chi phí món
```

## 🎯 KẾT QUẢ

Sau khi fix:
- ✅ Món "áddd" tính được chi phí từ batch "cfe cot"
- ✅ Chi phí = 200ml × 36đ/ml = 7,200đ (ví dụ)
- ✅ Cost Status = "✓ Chính thức"
- ✅ Tất cả món sử dụng batch ingredients đều tính được chi phí

## 📝 FILES THAY ĐỔI

1. **backend/application/services/batch_record_cost_adapter.go** (NEW)
   - Adapter để convert giữa repository interfaces

2. **backend/main.go** (MODIFIED)
   - Wire batch repository vào cost calculator service

## 🧪 TESTING

Để test:
1. Restart backend
2. Vào http://localhost:5173/#/manager/menu-costs
3. Tìm món "áddd"
4. Chi phí sẽ được tính tự động từ batch "cfe cot"

## 💡 LƯU Ý

- Batch ingredients cần có batch records available
- Nếu không có batch nào available → Cost Status = INCOMPLETE
- Hệ thống dùng FIFO (First In First Out) để lấy batch cũ nhất
- Chi phí batch được tính từ công thức nguyên liệu thô
