# Requirement: Tính giá thành menu item với tuỳ chọn trừ tồn kho

**Ngày tạo:** 2026-03-12
**Trạng thái:** Draft

---

## 1. Bối cảnh & Vấn đề

### Hành vi hiện tại

Khi thêm nguyên liệu vào menu item, hệ thống phân biệt 2 loại:
- `ingredient_type = "batch"` → khi có order, hệ thống **tự động trừ** `BatchRecord.QuantityRemaining`
- `ingredient_type = "raw"` → khi order, **không trừ** gì; raw ingredient chỉ bị trừ (`ingredient.Quantity`) khi **sản xuất batch**. Raw ingredient có `CostPerUnit` riêng để tính cost.

### Vấn đề

Manager muốn kiểm soát **từng nguyên liệu** xem có nên trừ tồn kho khi order hay không:
- Một số nguyên liệu muốn track tồn kho đầy đủ (trừ khi order)
- Một số nguyên liệu chỉ muốn **tính giá thành**, không muốn trừ kho

Yêu cầu này áp dụng cho **cả raw lẫn batch ingredient**.

### Giải pháp

Thêm **checkbox "Trừ tồn kho"** cho từng nguyên liệu trong form thêm/sửa menu item. Mọi ingredient đều được dùng để tính cost; việc trừ kho là tuỳ chọn riêng cho từng ingredient.

---

## 2. Thay đổi Data Model

### 2.1 Thêm field `deduct_inventory` vào `menu.Ingredient`

**File:** `backend/domain/menu/menu.go`

```go
type Ingredient struct {
    Name            string              `bson:"name" json:"name"`
    Quantity        float64             `bson:"quantity" json:"quantity"`
    Unit            ingredient.UnitType `bson:"unit" json:"unit"`
    IngredientType  string              `bson:"ingredient_type" json:"ingredient_type"` // "raw" | "batch"
    BatchID         *primitive.ObjectID `bson:"batch_id,omitempty" json:"batch_id,omitempty"`
    IngredientID    *primitive.ObjectID `bson:"ingredient_id,omitempty" json:"ingredient_id,omitempty"` // Reference cho raw

    // MỚI: nếu false → nguyên liệu chỉ dùng để tính cost, không trừ tồn kho khi order
    DeductInventory bool `bson:"deduct_inventory" json:"deduct_inventory"`
}
```

### 2.2 Ma trận hành vi theo loại và checkbox

| ingredient_type | DeductInventory | Khi order | Cost calculation |
|---|---|---|---|
| `batch` | `true` | Trừ `BatchRecord.QuantityRemaining` (FIFO) | Dùng BatchRecord/BatchDefinition cost |
| `batch` | `false` | **Không trừ gì** | Dùng BatchRecord/BatchDefinition cost |
| `raw` | `true` | **Trừ `ingredient.Quantity`** (MỚI) | Dùng `ingredient.CostPerUnit` |
| `raw` | `false` | Không trừ gì | Dùng `ingredient.CostPerUnit` |

**Backward compatibility:** Field `deduct_inventory` không tồn tại trong document cũ → default `true` (giữ nguyên hành vi cũ cho cả raw lẫn batch).

> ⚠️ Lưu ý: Hiện tại raw ingredient KHÔNG bị trừ khi order. Thêm `deduct_inventory=true` cho raw ingredient là **hành vi MỚI**. Cần cân nhắc kỹ trước khi bật cho raw ingredient trên các món đang chạy.

---

## 3. Thay đổi Backend

### 3.1 Rename và mở rộng hàm deduction — `order_service.go`

#### Thời điểm trừ tồn kho: **khi order chuyển sang PAID**

> ⚠️ **Thay đổi timing quan trọng:**
> - Hiện tại: `deductBatchIngredients()` được gọi ngay trong `CreateOrder()` (khi tạo order)
> - Yêu cầu mới: deduction phải xảy ra khi order **đã qua PAID** — tức là trong `CollectPayment()` ngay sau khi `o.IsFullyPaid() == true` và `o.Status = order.StatusPaid`
>
> Lý do: Order có thể bị huỷ sau khi tạo nhưng trước khi thanh toán. Trừ kho sớm tạo ra sai lệch tồn kho cho các order bị huỷ.

**Di chuyển lời gọi từ `CreateOrder()` sang `CollectPayment()`:**

```go
// Trong CollectPayment():
if o.IsFullyPaid() {
    o.Status = order.StatusPaid
    o.PaidAt = &now

    // Trừ tồn kho SAU KHI order được xác nhận thanh toán
    if err := s.deductIngredients(ctx, o); err != nil {
        // Rollback payment? Hoặc log warning và tiếp tục?
        // → Khuyến nghị: log error, không rollback payment (tránh UX xấu)
        log.Printf("WARNING: Failed to deduct inventory for paid order %s: %v", o.ID.Hex(), err)
    }
}
```

Đổi tên `deductBatchIngredients()` → `deductIngredients()` để phản ánh rằng cả raw lẫn batch đều có thể bị trừ:

```go
// deductIngredients deducts inventory for all ingredients that have DeductInventory=true
// Called when order status transitions to PAID
func (s *OrderService) deductIngredients(ctx context.Context, o *order.Order) (float64, error) {
    totalCost := 0.0

    for _, item := range o.Items {
        menuItem, err := s.menuRepo.FindByID(ctx, item.MenuItemID)
        if err != nil { return 0, err }

        ingredients := menuItem.GetIngredients(item.VariantID)

        for _, ing := range ingredients {
            // Bỏ qua nếu không tick "trừ tồn kho"
            if !ing.DeductInventory {
                continue
            }

            quantityNeeded := ing.Quantity * float64(item.Quantity)

            if ing.IsBatchIngredient() && ing.BatchID != nil {
                // Trừ BatchRecord (hành vi hiện tại)
                req := UseBatchRequest{ ... }
                result, err := s.batchUsageService.UseBatch(ctx, req)
                totalCost += result.TotalCost

            } else if ing.IsRawIngredient() && ing.IngredientID != nil {
                // MỚI: Trừ ingredient.Quantity trực tiếp
                err := s.ingredientService.DeductStock(ctx, *ing.IngredientID, quantityNeeded, ing.Unit, o.ID)
                // ...
            }
        }
    }
    return totalCost, nil
}
```

**Rollback:** Khi order thất bại:
- Batch ingredient: `batchUsageService.RollbackBatchUsage()` (đã có)
- Raw ingredient: `ingredientService.RestoreStock()` (cần thêm mới)

### 3.2 Thêm `DeductStock` / `RestoreStock` vào IngredientService

Để hỗ trợ trừ raw ingredient khi order và rollback khi thất bại:

```go
// DeductStock trừ ingredient.Quantity khi order, ghi StockHistory
func (s *IngredientService) DeductStock(
    ctx context.Context,
    ingredientID primitive.ObjectID,
    quantity float64,
    unit ingredient.UnitType,
    orderID primitive.ObjectID,
) error

// RestoreStock hoàn trả ingredient.Quantity khi order bị rollback
func (s *IngredientService) RestoreStock(
    ctx context.Context,
    ingredientID primitive.ObjectID,
    quantity float64,
    unit ingredient.UnitType,
    orderID primitive.ObjectID,
) error
```

`StockHistory` ghi lại lý do: `"Trừ kho theo order #ORDER_NUMBER"` / `"Hoàn trả kho — order rollback"`.

### 3.3 Cost calculation — `cost_calculator_service.go`

**Không thay đổi.** Hàm `calculateIngredientsCost()` dùng **tất cả ingredient**, bất kể `DeductInventory`:

```go
// DeductInventory không ảnh hưởng đến cost calculation
for _, ing := range ingredients {
    // tính cost như hiện tại...
}
```

### 3.4 Nội suy batch cost từ BatchDefinition (ESTIMATED)

Khi ingredient `type=batch` nhưng **không có BatchRecord available** → nội suy từ BatchDefinition:

```go
// Trong getBatchCostPerUnit():
// 1. BatchRecord available → CostPerUnit (FINAL)
// 2. Không có BatchRecord → interpolateBatchCostFromDefinition() → (ESTIMATED)
// 3. Không có BatchDefinition → (INCOMPLETE)

func (s *CostCalculatorService) interpolateBatchCostFromDefinition(
    ctx context.Context,
    batchDefID primitive.ObjectID,
) (costPerUnit float64, status CostStatus, err error)
```

**Công thức nội suy:**
```
Cho mỗi nguyên liệu [i] trong BatchDefinition.ConversionRates:
  qty_per_unit  = source_quantity[i] / batch_quantity[i]
  actual_qty    = qty_per_unit × (1 + wastage_rate[i])
  stock_qty     = convertUnit(actual_qty, source_unit[i] → ingredient_stock_unit[i])
  cost_contrib  = stock_qty × ingredient[i].cost_per_unit

batch_cost_per_unit = Σ(cost_contrib)
```

---

## 4. Thay đổi Frontend — MenuView (Ingredient Input)

### 4.1 Thêm toggle "Trừ tồn kho" vào mỗi hàng nguyên liệu

Toggle hiển thị cho **cả raw lẫn batch**:

```
┌─────────────────────────────────────────────────────────────────┐
│ Nguyên liệu                                                     │
│                                                                 │
│  Tên              SL    Đơn vị    Trừ tồn kho    Xoá          │
│  ─────────────────────────────────────────────────────────────  │
│  Hạt cà phê 🧪    20    g         [✓ Có]          [x]         │  ← batch, trừ kho
│  Sữa tươi 🥬      150   ml        [✓ Có]          [x]         │  ← raw, trừ kho (mới)
│  Đường 🥬         10    g         [ Không]        [x]         │  ← raw, chỉ cost
│                                                                 │
│                                            [+ Thêm nguyên liệu] │
└─────────────────────────────────────────────────────────────────┘
```

**Hành vi:**
- **Mặc định: BẬT** cho ingredient mới — giữ backward compatibility
- Hiển thị icon loại: 🧪 batch, 🥬 raw
- Khi TẮT: hiển thị badge "💰 Chỉ cost" bên cạnh tên

**Tooltip:**
> **BẬT:** Khi có order, tồn kho nguyên liệu này sẽ tự động bị trừ.
> **TẮT:** Nguyên liệu chỉ dùng để tính giá thành, không trừ tồn kho.

### 4.2 Raw ingredient cần lưu thêm `ingredient_id`

Hiện tại khi save, raw ingredient **không** lưu `ingredient_id` vào payload. Cần bổ sung để backend có thể trừ kho:

```js
// Trong saveItem() — hiện tại:
ingredients: form.value.ingredients.map(ing => ({
    name: ing.name,
    quantity: ing.quantity,
    unit: ing.unit,
    ...(ing.type === 'batch' && { ingredient_type: 'batch', batch_id: ing.id })
}))

// Sau khi sửa:
ingredients: form.value.ingredients.map(ing => ({
    name: ing.name,
    quantity: ing.quantity,
    unit: ing.unit,
    deduct_inventory: ing.deductInventory,
    ...(ing.type === 'batch' && { ingredient_type: 'batch', batch_id: ing.id }),
    ...(ing.type === 'raw' && { ingredient_type: 'raw', ingredient_id: ing.id })
}))
```

---

## 5. Frontend — `/manager/menu-costs` (MenuCostView) là nơi phản ánh cost

MenuCostView tại route `/manager/menu-costs` là **trang duy nhất phản ánh giá thành** cho tất cả menu item. Mọi thay đổi về ingredient (thêm/sửa/xoá trong MenuView) phải được phản ánh ngay tại đây.

### 5.1 Cost từ tất cả ingredient — kể cả cost-only

MenuCostView tính và hiển thị cost từ **tất cả ingredient**, bất kể `DeductInventory`:

- `DeductInventory = true` → ingredient trừ kho khi order → **vẫn tính vào cost**
- `DeductInventory = false` → ingredient chỉ tính cost, không trừ kho → **vẫn tính vào cost**

Đây là điểm cốt lõi: **cost phản ánh toàn bộ nguyên liệu của món**, không phụ thuộc vào cơ chế quản lý kho.

### 5.2 Thêm badge ESTIMATED

Batch ingredient không có BatchRecord → nội suy → `ESTIMATED`:
- Badge vàng `~` (hiện chỉ có FINAL và INCOMPLETE)
- Tooltip: "Giá ước tính từ công thức nguyên liệu. Chưa có lô sản xuất đang hoạt động."

### 5.3 Breakdown chi tiết: phân biệt "Trừ kho" vs "Chỉ cost"

Trong màn hình detail của từng món (click vào item trong MenuCostView), hiển thị rõ từng ingredient:

| Nguyên liệu | Loại | SL | Thành tiền | Kho |
|---|---|---|---|---|
| Hạt cà phê | 🧪 batch | 20g | 10,000₫ | ✓ Trừ kho |
| Sữa tươi | 🥬 raw | 150ml | 7,500₫ | ✓ Trừ kho |
| Đường | 🥬 raw | 10g | 500₫ | 💰 Chỉ cost |

Tổng cost = 17,500₫ (bao gồm cả ingredient "Chỉ cost")

### 5.4 Trigger recalculate

Khi manager lưu menu item (thêm/sửa nguyên liệu trong MenuView) → backend tự động tái tính cost và cập nhật `current_cost` + `cost_status` trên menu item → MenuCostView hiển thị giá trị mới ngay khi refresh.

---

## 6. Quy tắc CostStatus (đồng nhất toàn hệ thống)

| Status | Badge | Màu | Điều kiện |
|---|---|---|---|
| `FINAL` | ✓ | Xanh lá | Tất cả ingredient có giá; batch có BatchRecord thực tế |
| `ESTIMATED` | ~ | Vàng | Ít nhất 1 batch ingredient nội suy từ BatchDefinition |
| `INCOMPLETE` | ✗ | Đỏ | Ít nhất 1 ingredient không có giá |

---

## 7. Tác động & Không thay đổi

| Hành vi | Thay đổi? |
|---|---|
| Batch `DeductInventory=true` khi order | Không đổi — vẫn trừ BatchRecord |
| Batch `DeductInventory=false` khi order | **MỚI: bỏ qua** |
| Raw `DeductInventory=true` khi order | **MỚI: trừ `ingredient.Quantity`** |
| Raw `DeductInventory=false` khi order | Không đổi — không trừ |
| Cost calculation | Không đổi — dùng tất cả ingredient |
| Batch production deduction (raw → batch) | Không đổi |
| `CalculateShiftOrderCosts` khi đóng ca | Không đổi |
| Timing trừ kho: từ `CreateOrder()` → `CollectPayment()` khi PAID | **THAY ĐỔI** |
| Rollback batch khi order thất bại (tại CreateOrder) | **Bỏ** — deduction không còn ở CreateOrder |
| Rollback raw khi order thất bại | **MỚI: cần `RestoreStock()` nếu deduction sau PAID gặp lỗi** |

---

## 8. Migration

**Existing batch ingredients** (`ingredient_type=batch`, không có `deduct_inventory`):
→ Default `true` → không thay đổi hành vi ✓

**Existing raw ingredients** (`ingredient_type=raw` hoặc rỗng, không có `deduct_inventory`):
→ Default `true` → **THAY ĐỔI hành vi** (trước: không trừ, sau: trừ kho)

> ⚠️ **Quan trọng:** Để tránh ảnh hưởng hệ thống đang chạy, migration nên set:
> - `deduct_inventory = true` cho tất cả **batch ingredient** hiện có (giữ nguyên)
> - `deduct_inventory = false` cho tất cả **raw ingredient** hiện có (giữ nguyên — không trừ kho)
>
> Việc bật `deduct_inventory=true` cho raw ingredient là quyết định của manager, không phải default.

---

## 9. Ưu tiên triển khai

| Thứ tự | Hạng mục | File |
|---|---|---|
| 1 | Thêm `deduct_inventory`, `ingredient_id` vào `menu.Ingredient` struct | `domain/menu/menu.go` |
| 2 | Thêm `DeductStock()` / `RestoreStock()` vào IngredientService | `ingredient.go` |
| 3 | Đổi `deductBatchIngredients()` → `deductIngredients()`, xử lý cả raw | `order_service.go` |
| 4 | Backend: `interpolateBatchCostFromDefinition()` + cập nhật `getBatchCostPerUnit()` | `cost_calculator_service.go` |
| 5 | Frontend: toggle "Trừ tồn kho" + lưu `ingredient_id` + `deduct_inventory` | `MenuView.vue` |
| 6 | Frontend: badge ESTIMATED trong MenuCostView | `MenuCostView.vue` |
| 7 | Migration script: batch=true, raw=false | script |
