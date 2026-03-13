# Implementation Spec: Deduct Inventory Toggle + Batch Cost Interpolation

**Requirement:** `REQUIREMENTS_STANDALONE_COST_CALCULATOR.md`
**Ngày tạo:** 2026-03-12

---

## Tổng quan

Triển khai 2 nhóm tính năng độc lập:

**Nhóm A — Deduct Inventory Toggle:**
Thêm checkbox "Trừ tồn kho" cho từng ingredient trong menu item. Cả raw lẫn batch có thể bật/tắt độc lập. Inventory deduction chuyển từ `CreateOrder` → `CollectPayment` (khi PAID).

**Nhóm B — Batch Cost Interpolation:**
Khi batch ingredient không có BatchRecord available, thay vì INCOMPLETE → nội suy từ BatchDefinition → ESTIMATED. MenuCostView hiển thị badge ~ ESTIMATED.

---

## TASK A-1: Data model — thêm `deduct_inventory` và `ingredient_id`

**File:** `backend/domain/menu/menu.go`

**Thay đổi struct `Ingredient`:**
```go
type Ingredient struct {
    Name           string              `bson:"name" json:"name"`
    Quantity       float64             `bson:"quantity" json:"quantity"`
    Unit           ingredient.UnitType `bson:"unit" json:"unit"`
    IngredientType string              `bson:"ingredient_type" json:"ingredient_type"`
    BatchID        *primitive.ObjectID `bson:"batch_id,omitempty" json:"batch_id,omitempty"`

    // THÊM MỚI:
    IngredientID    *primitive.ObjectID `bson:"ingredient_id,omitempty" json:"ingredient_id,omitempty"`
    DeductInventory bool               `bson:"deduct_inventory" json:"deduct_inventory"`
}
```

**Backward compatibility trong `GetIngredients()` hoặc bất kỳ nơi nào đọc ingredient từ DB:**
Không cần code xử lý đặc biệt — BSON unmarshal field `bool` mà không có trong document sẽ ra `false`. Xem task A-6 (migration) để set đúng default.

**Validate:** Cập nhật `ValidateIngredients()` nếu có — không cần validate `ingredient_id` là bắt buộc (backward compat).

---

## TASK A-2: IngredientService — thêm `DeductStock` và `RestoreStock`

**File:** `backend/application/services/ingredient.go`

**2 method mới:**

```go
// DeductStock trừ ingredient.Quantity khi order PAID, ghi StockHistory
// unit: đơn vị trong recipe (có thể khác stock unit → cần convertUnit)
func (s *IngredientService) DeductStock(
    ctx context.Context,
    ingredientID primitive.ObjectID,
    quantity float64,
    unit ingredient.UnitType,
    orderID primitive.ObjectID,
    orderNumber string,
) error {
    ing, err := s.ingredientRepo.FindByID(ctx, ingredientID)
    // convert unit → stock unit
    convRate := ingredient.GetConversionRate(ing.Unit, unit)
    qtyInStock := quantity * convRate

    beforeQty := ing.Quantity
    ing.Quantity -= qtyInStock
    if ing.Quantity < 0 { ing.Quantity = 0 }

    err = s.ingredientRepo.Update(ctx, ingredientID, ing)
    // ghi StockHistory: type=TransactionOrderDeduct, reason="Order #ORDER_NUMBER"
}

// RestoreStock hoàn trả ingredient.Quantity nếu cần rollback
func (s *IngredientService) RestoreStock(
    ctx context.Context,
    ingredientID primitive.ObjectID,
    quantity float64,
    unit ingredient.UnitType,
    orderID primitive.ObjectID,
) error
```

**StockHistory transaction type mới:** Thêm constant `TransactionOrderDeduct` vào `domain/ingredient/ingredient.go` nếu chưa có.

---

## TASK A-3: OrderService — chuyển deduction sang PAID, hỗ trợ raw

**File:** `backend/application/services/order_service.go`

### A-3a: Xoá deduction khỏi `CreateOrder()`

```go
// XOÁ đoạn này:
if s.batchUsageService != nil {
    batchCost, err := s.deductBatchIngredients(ctx, o)
    ...
}
```

### A-3b: Thêm deduction vào `CollectPayment()` tại điểm PAID

```go
if o.IsFullyPaid() {
    o.Status = order.StatusPaid
    o.PaidAt = &now

    // Trừ tồn kho khi order PAID
    if err := s.deductIngredients(ctx, o); err != nil {
        log.Printf("WARNING: inventory deduction failed for paid order %s: %v", o.ID.Hex(), err)
        // Không rollback payment — chỉ log
    }
}
```

### A-3c: Đổi tên + mở rộng hàm

```go
// RENAME deductBatchIngredients → deductIngredients
// Thêm xử lý raw ingredient với DeductInventory=true
func (s *OrderService) deductIngredients(ctx context.Context, o *order.Order) error {
    for _, item := range o.Items {
        menuItem, _ := s.menuRepo.FindByID(ctx, item.MenuItemID)
        ingredients := menuItem.GetIngredients(item.VariantID)

        for _, ing := range ingredients {
            if !ing.DeductInventory {
                continue // bỏ qua nếu không tick "trừ tồn kho"
            }

            qty := ing.Quantity * float64(item.Quantity)

            switch {
            case ing.IsBatchIngredient() && ing.BatchID != nil:
                // batch: dùng batchUsageService.UseBatch() (hiện tại)
                req := UseBatchRequest{
                    BatchDefinitionID: *ing.BatchID,
                    QuantityNeeded:    qty,
                    Unit:              string(ing.Unit),
                    OrderID:           o.ID,
                    MenuItemID:        item.MenuItemID,
                    MenuItemName:      item.Name,
                }
                s.batchUsageService.UseBatch(ctx, req)

            case ing.IsRawIngredient() && ing.IngredientID != nil:
                // raw: dùng ingredientService.DeductStock() (MỚI)
                s.ingredientService.DeductStock(ctx, *ing.IngredientID, qty, ing.Unit, o.ID, o.OrderNumber)
            }
        }
    }
    return nil
}
```

**Inject `ingredientService` vào `OrderService`:** Thêm field và constructor param nếu chưa có.

### A-3d: Xoá rollback batch khỏi CreateOrder

Rollback trong `CreateOrder()` gắn liền với deduction cũ → xoá đoạn `RollbackBatchUsage` trong `CreateOrder()`.

---

## TASK A-4: Migration script

**File mới:** `scripts/migrate-deduct-inventory.sh` (hoặc Go script)

**Logic:**
```js
// batch ingredients (ingredient_type = "batch"):
db.menu_items.updateMany(
  { "ingredients.ingredient_type": "batch", "ingredients.deduct_inventory": { $exists: false } },
  { $set: { "ingredients.$[elem].deduct_inventory": true } },
  { arrayFilters: [{ "elem.ingredient_type": "batch" }] }
)

// raw ingredients (ingredient_type = "raw" hoặc ""):
db.menu_items.updateMany(
  { "ingredients.deduct_inventory": { $exists: false } },
  { $set: { "ingredients.$[elem].deduct_inventory": false } },
  { arrayFilters: [{ "elem.deduct_inventory": { $exists: false } }] }
)

// Lặp tương tự cho variants[].ingredients
```

**Kết quả mong đợi:**
- Batch ingredient cũ: `deduct_inventory = true` (giữ nguyên hành vi)
- Raw ingredient cũ: `deduct_inventory = false` (giữ nguyên hành vi — không trừ)

---

## TASK A-5: Frontend — MenuView toggle "Trừ tồn kho"

**File:** `frontend/src/views/MenuView.vue`

### A-5a: Cập nhật `selectIngredient()` — thêm `deductInventory` default

```js
// Khi thêm raw ingredient:
form.value.ingredients.push({
    id: ingredient.id,
    type: 'raw',
    name: ingredient.name,
    quantity: 1,
    unit: recipeUnit,
    deductInventory: true,  // MỚI: default BẬT
    // ... các field hiện có
})

// Khi thêm batch ingredient (selectBatch):
variant.ingredients.push({
    // ... fields hiện có
    deductInventory: true,  // MỚI: default BẬT
})
```

### A-5b: Cập nhật `editItem()` — load `deduct_inventory` từ DB

```js
// Khi parse ingredient từ response:
{
    deductInventory: ing.deduct_inventory ?? true,  // default true nếu không có
    // ... các field hiện có
}
```

### A-5c: Thêm toggle UI trong ingredient list

Trong template, mỗi hàng ingredient (cả single-size lẫn variant):

```html
<!-- Toggle Trừ tồn kho -->
<div class="flex items-center gap-1">
  <button
    @click="ing.deductInventory = !ing.deductInventory"
    :class="[
      'px-2 py-1 rounded text-xs font-medium transition-colors',
      ing.deductInventory
        ? 'bg-red-100 text-red-700'
        : 'bg-gray-100 text-gray-500'
    ]"
    :title="ing.deductInventory ? 'Đang trừ tồn kho khi order' : 'Chỉ tính cost, không trừ kho'"
  >
    {{ ing.deductInventory ? '📦 Trừ kho' : '💰 Cost' }}
  </button>
</div>
```

### A-5d: Cập nhật `saveItem()` — gửi `ingredient_id` và `deduct_inventory`

```js
ingredients: form.value.ingredients.map(ing => ({
    name: ing.name,
    quantity: ing.quantity,
    unit: ing.unit,
    deduct_inventory: ing.deductInventory,
    ...(ing.type === 'batch' && {
        ingredient_type: 'batch',
        batch_id: ing.batch_definition_id || ing.id
    }),
    ...(ing.type === 'raw' && {
        ingredient_type: 'raw',
        ingredient_id: ing.id   // MỚI: lưu ingredient_id cho raw
    })
}))
```

Áp dụng tương tự cho **variant ingredients**.

---

## TASK B-1: CostCalculatorService — thêm batch cost interpolation

**File:** `backend/application/services/cost_calculator_service.go`

### B-1a: Thêm hàm `interpolateBatchCostFromDefinition()`

```go
// interpolateBatchCostFromDefinition tính CostPerUnit của batch
// từ giá nguyên liệu thô hiện tại trong BatchDefinition.ConversionRates.
// Dùng khi không có BatchRecord available.
func (s *CostCalculatorService) interpolateBatchCostFromDefinition(
    ctx context.Context,
    batchDefID primitive.ObjectID,
) (float64, CostStatus, error) {
    def, err := s.batchDefRepo.FindByID(ctx, batchDefID)
    if err != nil { return 0, CostStatusIncomplete, err }

    totalCostPerUnit := 0.0
    for _, rate := range def.ConversionRates {
        ing, err := s.ingredientRepo.FindByID(ctx, rate.SourceIngredientID)
        if err != nil || ing.CostPerUnit <= 0 {
            return 0, CostStatusIncomplete, fmt.Errorf("missing cost for %s", rate.SourceIngredientName)
        }

        qtyPerUnit := rate.SourceQuantity / rate.BatchQuantity
        actualQty  := qtyPerUnit * (1 + rate.WastageRate)
        // convert từ source_unit → ingredient stock unit
        convRate   := ingredient.GetConversionRate(ingredient.UnitType(ing.Unit), ingredient.UnitType(rate.SourceUnit))
        stockQty   := actualQty * convRate
        totalCostPerUnit += stockQty * ing.CostPerUnit
    }

    return totalCostPerUnit, CostStatusEstimated, nil
}
```

### B-1b: Cập nhật `getBatchCostPerUnit()`

```go
func (s *CostCalculatorService) getBatchCostPerUnit(ctx context.Context, batchDefID primitive.ObjectID) (float64, CostStatus, error) {
    // 1. Tìm BatchRecord available (hiện tại)
    records, err := s.batchRecordRepo.FindAvailableByDefinition(ctx, batchDefID)
    if err == nil && len(records) > 0 {
        // FIFO: lấy record cũ nhất
        return records[0].CostPerUnit, CostStatusFinal, nil
    }

    // 2. MỚI: nội suy từ BatchDefinition
    cost, status, err := s.interpolateBatchCostFromDefinition(ctx, batchDefID)
    if err != nil {
        return 0, CostStatusIncomplete, err
    }
    return cost, status, nil
}
```

**Lưu ý:** Kiểm tra signature hiện tại của `getBatchCostPerUnit()` — nếu chỉ trả `(float64, error)` thì cần thêm `CostStatus` vào return.

---

## TASK B-2: Frontend — MenuCostView badge ESTIMATED + breakdown

**File:** `frontend/src/views/MenuCostView.vue`

### B-2a: Thêm badge ESTIMATED

Tìm nơi hiển thị `cost_status` badge (hiện có FINAL ✓ và INCOMPLETE ✗), thêm case ESTIMATED:

```js
// Trong computed/method getCostStatusBadge hoặc tương đương:
const getCostStatusBadge = (status) => {
    switch (status) {
        case 'FINAL':      return { text: '✓', class: 'bg-green-100 text-green-700', title: 'Giá chính xác' }
        case 'ESTIMATED':  return { text: '~', class: 'bg-yellow-100 text-yellow-700', title: 'Giá ước tính từ công thức nguyên liệu' }  // MỚI
        case 'INCOMPLETE': return { text: '✗', class: 'bg-red-100 text-red-700', title: 'Thiếu dữ liệu giá' }
        default:           return { text: '?', class: 'bg-gray-100 text-gray-500', title: '' }
    }
}
```

### B-2b: Breakdown detail — cột "Kho"

Trong component detail (modal hoặc inline), thêm cột phân biệt:

```html
<td>
  <span v-if="ingredient.deduct_inventory" class="text-xs text-blue-600">📦 Trừ kho</span>
  <span v-else class="text-xs text-gray-400">💰 Chỉ cost</span>
</td>
```

**Lưu ý:** Backend cần trả `deduct_inventory` trong response breakdown. Kiểm tra `GET /manager/menu/costs/:id`.

---

## Thứ tự implement (dependencies)

```
A-1 (domain model)
  └─> A-2 (IngredientService.DeductStock)
  └─> A-3 (OrderService deductIngredients)
        └─> requires A-2
  └─> A-5 (Frontend MenuView)
        └─> requires A-1 (field mới)
A-4 (migration) — chạy sau A-1 deploy

B-1 (CostCalculator interpolation) — độc lập
  └─> B-2 (Frontend MenuCostView badge) — requires B-1
```

**Có thể làm song song:** A-1→A-2→A-3 và B-1→B-2 là 2 luồng độc lập.

---

## Checklist kiểm tra trước khi deploy

- [ ] Migration script đã chạy trên production DB
- [ ] Order cũ (tạo trước deploy) không bị trừ kho 2 lần
- [ ] Batch ingredient cũ (không có `deduct_inventory`) vẫn bị trừ đúng khi order PAID
- [ ] Raw ingredient cũ vẫn không bị trừ khi order (migration set false)
- [ ] Cost trong MenuCostView hiển thị đúng cho cả ingredient có/không `deduct_inventory`
- [ ] ESTIMATED badge hiển thị khi batch không có BatchRecord
- [ ] FINAL vẫn hiển thị khi batch có BatchRecord
