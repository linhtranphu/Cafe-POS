# Conversion Rate - Cách Sử Dụng Hiện Tại

## 📊 Tổng Quan

Hiện tại, `conversion_rate` đang được **LƯU TRONG DATABASE** ở bảng `ingredients` và được sử dụng trong cost analysis theo cách sau:

## 🔍 Cách Sử Dụng Hiện Tại

### 1. Lưu Trữ Dữ Liệu

**File**: `backend/domain/ingredient/ingredient.go`

```go
type Ingredient struct {
    ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name              string             `bson:"name" json:"name"`
    Unit              UnitType           `bson:"unit" json:"unit"`              // "L" (stock unit)
    CostPerUnit       float64            `bson:"cost_per_unit" json:"cost_per_unit"`
    ConversionRate    float64            `bson:"conversion_rate" json:"conversion_rate"`       // ⚠️ LƯU TRONG DB
    WastagePercentage float64            `bson:"wastage_percentage" json:"wastage_percentage"`
    CurrentStock      float64            `bson:"quantity" json:"quantity"`
}
```

**Giá trị mặc định**: `1.0` (nếu không được set)

### 2. Công Thức Tính Cost

**File**: `backend/application/services/cost_calculator_service.go`

```go
// Formula: quantity * cost_per_unit * conversion_rate * (1 + wastage_percentage/100)
ingredientCost := menuIngredient.Quantity * ing.CostPerUnit * conversionRate * (1 + wastagePercentage/100)
```

**Ví dụ cụ thể**:
```
Menu Item: Cà phê sữa đá
Ingredient: Sữa tươi
- Quantity (công thức): 150 ml
- CostPerUnit (kho): 50,000 VND/L
- ConversionRate (DB): 0.001  ← ⚠️ LƯU SẴN TRONG DB
- WastagePercentage: 5%

Cost = 150 × 50,000 × 0.001 × (1 + 5/100)
     = 150 × 50,000 × 0.001 × 1.05
     = 7,875 VND
```

### 3. Các Hàm Sử Dụng Conversion Rate

#### a) `CalculateMenuItemCost()` - Tính cost cho 1 món
```go
func (s *CostCalculatorService) CalculateMenuItemCost(ctx context.Context, menuItemID primitive.ObjectID) (*MenuItemCostResult, error) {
    // ...
    
    // Get conversion rate (default 1.0 if not set)
    conversionRate := ing.ConversionRate
    if conversionRate <= 0 {
        conversionRate = 1.0
    }
    
    // Calculate cost
    ingredientCost := menuIngredient.Quantity * ing.CostPerUnit * conversionRate * (1 + wastagePercentage/100)
    totalCost += ingredientCost
    
    // ...
}
```

#### b) `calculateCostForMenuItem()` - Helper cho batch calculation
```go
func (s *CostCalculatorService) calculateCostForMenuItem(menuItem *menu.MenuItem, ingredientMap map[string]*ingredient.Ingredient) *MenuItemCostResult {
    // ...
    
    conversionRate := ing.ConversionRate
    if conversionRate <= 0 {
        conversionRate = 1.0
    }
    
    ingredientCost := menuIngredient.Quantity * ing.CostPerUnit * conversionRate * (1 + wastagePercentage/100)
    totalCost += ingredientCost
    
    // ...
}
```

#### c) `CalculateMenuItemCostDetail()` - Tính cost với breakdown chi tiết
```go
func (s *CostCalculatorService) CalculateMenuItemCostDetail(ctx context.Context, menuItemID primitive.ObjectID) (*MenuItemCostDetail, error) {
    // ...
    
    conversionRate := ing.ConversionRate
    if conversionRate <= 0 {
        conversionRate = 1.0
    }
    
    ingredientCost := menuIngredient.Quantity * ing.CostPerUnit * conversionRate * (1 + wastagePercentage/100)
    
    detail.ConversionRate = conversionRate  // ← Trả về trong response
    detail.TotalCost = ingredientCost
    
    // ...
}
```

### 4. API Response Structure

**File**: `backend/interfaces/http/menu_cost_handler.go`

```go
type IngredientCostDetailResponse struct {
    Name              string  `json:"name"`
    Quantity          float64 `json:"quantity"`
    Unit              string  `json:"unit"`
    CostPerUnit       float64 `json:"cost_per_unit"`
    ConversionRate    float64 `json:"conversion_rate"`    // ← Trả về cho frontend
    WastagePercentage float64 `json:"wastage_percentage"`
    TotalCost         float64 `json:"total_cost"`
}
```

**Example Response**:
```json
{
  "menu_item_id": "...",
  "menu_item_name": "Cà phê sữa đá",
  "ingredients": [
    {
      "name": "Sữa tươi",
      "quantity": 150,
      "unit": "ml",
      "cost_per_unit": 50000,
      "conversion_rate": 0.001,    ← Frontend nhận được giá trị này
      "wastage_percentage": 5.0,
      "total_cost": 7875
    }
  ]
}
```

## ⚠️ VẤN ĐỀ VỚI CÁCH HIỆN TẠI

### Vấn đề 1: Inconsistency
```
Ingredient: Sữa tươi
- Stock Unit: L
- ConversionRate trong DB: 0.001

Menu Item A: Cà phê sữa đá
- Ingredient: Sữa tươi 150ml
- ConversionRate: 0.001 ✓ (ml → L)

Menu Item B: Sinh tố
- Ingredient: Sữa tươi 0.2L
- ConversionRate: 0.001 ✗ (WRONG! Should be 1.0 for L → L)
```

**Vấn đề**: Cùng 1 ingredient nhưng có thể dùng nhiều recipe units khác nhau (ml, L, cl...). Lưu 1 conversion_rate cố định trong DB sẽ không đúng cho tất cả trường hợp.

### Vấn đề 2: Không Linh Hoạt
```
Nếu muốn thay đổi recipe unit từ ml → cl:
- Phải update conversion_rate trong ingredient DB
- Ảnh hưởng đến TẤT CẢ các món đang dùng ingredient này
- Không thể có 2 món dùng 2 units khác nhau
```

### Vấn đề 3: Dễ Nhầm Lẫn
```
User tạo ingredient:
- Unit: L
- ConversionRate: ??? (User phải tự tính và nhập)

Nếu user nhập sai → Tất cả cost calculations đều sai!
```

## ✅ GIẢI PHÁP ĐỀ XUẤT

### Thay đổi 1: KHÔNG lưu conversion_rate trong Ingredient
```go
type Ingredient struct {
    ID                primitive.ObjectID
    Name              string
    Unit              UnitType           // Stock unit (L, kg, ...)
    CostPerUnit       float64
    WastagePercentage float64
    // ❌ REMOVE: ConversionRate float64
}
```

### Thay đổi 2: LƯU recipe unit trong MenuIngredient
```go
type MenuIngredient struct {
    Name     string
    Quantity float64
    Unit     UnitType  // ← Recipe unit (ml, g, L, kg, ...)
}
```

### Thay đổi 3: TÍNH conversion_rate động
```go
// In cost_calculator_service.go
func (s *CostCalculatorService) CalculateMenuItemCost(...) {
    // ...
    
    // Get stock unit from ingredient
    stockUnit := ing.Unit  // "L"
    
    // Get recipe unit from menu ingredient
    recipeUnit := menuIngredient.Unit  // "ml"
    
    // Calculate conversion rate dynamically
    conversionRate := ingredient.GetConversionRate(stockUnit, recipeUnit)  // 0.001
    
    // Calculate cost
    ingredientCost := menuIngredient.Quantity * ing.CostPerUnit * conversionRate * (1 + wastagePercentage/100)
    
    // ...
}
```

## 📈 So Sánh

| Aspect | Cách Hiện Tại | Cách Đề Xuất |
|--------|---------------|--------------|
| **Lưu trữ** | ConversionRate trong Ingredient DB | KHÔNG lưu, tính động |
| **Linh hoạt** | ❌ 1 ingredient = 1 conversion rate | ✅ Mỗi món có thể dùng unit khác nhau |
| **Chính xác** | ⚠️ Dễ sai nếu user nhập sai | ✅ Tự động tính, không thể sai |
| **Bảo trì** | ❌ Phải update DB khi đổi unit | ✅ Không cần update gì |
| **Consistency** | ❌ Không đảm bảo | ✅ Luôn đúng |

## 🎯 Kết Luận

**Hiện tại**: Conversion rate được **LƯU TRONG DATABASE** và sử dụng trực tiếp trong công thức tính cost.

**Vấn đề**: Không linh hoạt, dễ sai, không consistent.

**Giải pháp**: Chuyển sang **TÍNH ĐỘNG** dựa trên stock_unit và recipe_unit.

**Ưu điểm**:
- ✅ Tự động, chính xác
- ✅ Linh hoạt (mỗi món dùng unit riêng)
- ✅ Dễ bảo trì
- ✅ Consistent across all modules
