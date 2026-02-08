# Fix: Không Tạo Chi Phí Khi Nhận Nguyên Liệu Tặng

## 🎯 Vấn Đề

Khi điều chỉnh TĂNG số lượng nguyên liệu mà KHÔNG nhập giá (assume là được tặng, không phải mua), hệ thống vẫn tự động tạo chi phí.

### Ví Dụ Lỗi

```
Tồn kho hiện tại: 8kg @ 25,000đ/kg
User điều chỉnh tăng: 12kg (tăng 4kg)
User KHÔNG nhập giá (để 0 - vì được tặng)

❌ Hệ thống cũ:
- Tạo expense: 4kg × 25,000đ = 100,000đ
- Lý do: Backend dùng giá hiện tại (25,000đ) để tính expense

✅ Hệ thống mới:
- KHÔNG tạo expense
- Lý do: User không nhập giá → không phải mua → không có chi phí
```

---

## 🔍 Nguyên Nhân

### Backend Logic Cũ (SAI)

**File**: `backend/application/services/ingredient.go`

#### 1. StockAdjust Method (Line ~298)
```go
// Track expense if quantity increased with new price
if s.autoExpenseService != nil && quantityDiff > 0 && costPerUnit > 0 {
    // ❌ Vấn đề: costPerUnit luôn > 0 vì nó = item.CostPerUnit khi user không nhập giá
    tempItem := *item
    tempItem.CostPerUnit = costPerUnit
    s.autoExpenseService.TrackIngredientPurchase(ctx, &tempItem, quantityDiff, req.Username)
}
```

**Vấn đề**: 
- Khi user KHÔNG nhập giá, `req.CostPerUnit = 0`
- Nhưng `costPerUnit = item.CostPerUnit` (giá hiện tại)
- Điều kiện `costPerUnit > 0` vẫn TRUE → Tạo expense SAI!

#### 2. StockIn Method (Line ~202)
```go
// Track expense for purchase
if s.autoExpenseService != nil && costPerUnit > 0 {
    // ❌ Vấn đề tương tự
    tempItem := *item
    tempItem.CostPerUnit = costPerUnit
    s.autoExpenseService.TrackIngredientPurchase(ctx, &tempItem, req.Quantity, req.Username)
}
```

---

## ✅ Giải Pháp

### Logic Đúng

**Chỉ tạo expense khi user THỰC SỰ NHẬP GIÁ MỚI** (tức là `req.CostPerUnit > 0`)

### Backend Fix

#### 1. StockAdjust Method

```go
// Determine cost per unit for this transaction
costPerUnit := item.CostPerUnit // Default to current price
userProvidedNewPrice := false   // ✅ Track if user actually provided a new price

// Only recalculate price if:
// 1. Quantity increased (positive diff)
// 2. New price provided and different from current
if quantityDiff > 0 && req.CostPerUnit > 0 && req.CostPerUnit != item.CostPerUnit && afterQty > 0 {
    // Weighted average for the increase
    oldValue := beforeQty * item.CostPerUnit
    newValue := quantityDiff * req.CostPerUnit
    item.CostPerUnit = (oldValue + newValue) / afterQty
    costPerUnit = req.CostPerUnit // Use new price for history
    userProvidedNewPrice = true   // ✅ User provided a new price
}

// ... update database and create history ...

// Track expense ONLY if:
// 1. Quantity increased AND
// 2. User provided a new price (not gifted/found)
if s.autoExpenseService != nil && quantityDiff > 0 && userProvidedNewPrice {
    // ✅ Chỉ tạo expense khi user nhập giá
    tempItem := *item
    tempItem.CostPerUnit = costPerUnit
    s.autoExpenseService.TrackIngredientPurchase(ctx, &tempItem, quantityDiff, req.Username)
}
```

#### 2. StockIn Method

```go
// Determine cost per unit for this transaction
costPerUnit := req.CostPerUnit
userProvidedNewPrice := req.CostPerUnit > 0 // ✅ Track if user provided a new price

if costPerUnit <= 0 {
    // If no cost provided, use current cost for history
    costPerUnit = item.CostPerUnit
}

// Calculate weighted average cost ONLY when new price is provided and different
if req.CostPerUnit > 0 && req.CostPerUnit != item.CostPerUnit && afterQty > 0 {
    oldValue := beforeQty * item.CostPerUnit
    newValue := req.Quantity * req.CostPerUnit
    item.CostPerUnit = (oldValue + newValue) / afterQty
}

// ... update database and create history ...

// Track expense ONLY if user provided a new price (actual purchase, not gifted)
if s.autoExpenseService != nil && userProvidedNewPrice {
    // ✅ Chỉ tạo expense khi user nhập giá
    tempItem := *item
    tempItem.CostPerUnit = req.CostPerUnit // Use the price user provided
    s.autoExpenseService.TrackIngredientPurchase(ctx, &tempItem, req.Quantity, req.Username)
}
```

---

## 📊 Logic Tổng Hợp

### Adjust Tăng

| User nhập giá | req.CostPerUnit | userProvidedNewPrice | Tạo expense? | Lý do |
|---------------|-----------------|----------------------|--------------|-------|
| Có (30,000đ) | 30000 | true | ✅ CÓ | Mua hàng với giá mới |
| Không (để 0) | 0 | false | ❌ KHÔNG | Được tặng/tìm thấy |

### Stock In

| User nhập giá | req.CostPerUnit | userProvidedNewPrice | Tạo expense? | Lý do |
|---------------|-----------------|----------------------|--------------|-------|
| Có (30,000đ) | 30000 | true | ✅ CÓ | Mua hàng |
| Không (để 0) | 0 | false | ❌ KHÔNG | Được tặng |

### Stock Out

| Thao tác | Tạo expense? | Lý do |
|----------|--------------|-------|
| Xuất kho | ❌ KHÔNG | Xuất kho không phải chi phí |

---

## 🧪 Test Cases

### ✅ Test Case 1: Adjust Tăng KHÔNG Nhập Giá - KHÔNG Tạo Expense

**Bước thực hiện**:
1. Nguyên liệu: Đường (8kg @ 25,000đ/kg)
2. Click "📦 Điều chỉnh" → "Điều chỉnh"
3. Nhập số lượng: 12kg (tăng 4kg)
4. **KHÔNG nhập giá** (để 0)
5. Lý do: "Được tặng từ nhà cung cấp"
6. Click "Xác nhận"

**Kết quả mong đợi**:
- ✅ Tồn kho: 12kg
- ✅ Đơn giá: 25,000đ/kg (không đổi)
- ✅ **KHÔNG có expense mới trong danh sách chi phí**
- ✅ Lịch sử: Có record "Điều chỉnh +4kg" với lý do "Được tặng..."

**Kiểm tra**:
```bash
# Xem danh sách expense
curl http://localhost:3000/api/manager/expenses | jq '.[] | select(.category == "Nguyên liệu")'

# Không có expense mới với amount = 100,000đ (4kg × 25,000đ)
```

### ✅ Test Case 2: Adjust Tăng CÓ Nhập Giá - CÓ Tạo Expense

**Bước thực hiện**:
1. Nguyên liệu: Đường (8kg @ 25,000đ/kg)
2. Click "📦 Điều chỉnh" → "Điều chỉnh"
3. Nhập số lượng: 12kg (tăng 4kg)
4. **Nhập giá: 30,000đ**
5. Lý do: "Mua thêm với giá mới"
6. Click "Xác nhận"

**Kết quả mong đợi**:
- ✅ Tồn kho: 12kg
- ✅ Đơn giá: ~26,667đ/kg (weighted average)
- ✅ **CÓ expense mới**: 4kg × 30,000đ = 120,000đ
- ✅ Expense category: "Nguyên liệu"
- ✅ Expense payment method: "Tiền mặt"

**Kiểm tra**:
```bash
# Xem expense mới nhất
curl http://localhost:3000/api/manager/expenses | jq '.[0]'

# Phải có:
# - amount: 120000
# - category: "Nguyên liệu"
# - description: "Mua Đường (4kg)"
```

### ✅ Test Case 3: Stock In KHÔNG Nhập Giá - KHÔNG Tạo Expense

**Bước thực hiện**:
1. Nguyên liệu: Sữa (5 chai @ 10,000đ/chai)
2. Click "📦 Điều chỉnh" → "Nhập kho"
3. Nhập số lượng: 3 chai
4. **KHÔNG nhập giá** (để 0)
5. Lý do: "Khuyến mãi từ nhà cung cấp"
6. Click "Xác nhận"

**Kết quả mong đợi**:
- ✅ Tồn kho: 8 chai
- ✅ Đơn giá: 10,000đ/chai (không đổi)
- ✅ **KHÔNG có expense mới**

### ✅ Test Case 4: Stock In CÓ Nhập Giá - CÓ Tạo Expense

**Bước thực hiện**:
1. Nguyên liệu: Sữa (5 chai @ 10,000đ/chai)
2. Click "📦 Điều chỉnh" → "Nhập kho"
3. Nhập số lượng: 3 chai
4. **Nhập giá: 12,000đ/chai**
5. Lý do: "Mua thêm"
6. Click "Xác nhận"

**Kết quả mong đợi**:
- ✅ Tồn kho: 8 chai
- ✅ Đơn giá: ~10,750đ/chai (weighted average)
- ✅ **CÓ expense mới**: 3 × 12,000đ = 36,000đ

---

## 📁 Files Thay Đổi

### Backend
- `backend/application/services/ingredient.go`
  - **StockAdjust method** (Line ~265-305): Thêm `userProvidedNewPrice` flag
  - **StockIn method** (Line ~160-210): Thêm `userProvidedNewPrice` flag

### Frontend
- `frontend/src/views/IngredientManagementView.vue`
  - **adjustStock function** (Line ~1489-1507): Đã fix trước đó (chỉ gửi cost_per_unit khi increase)

---

## ✅ Checklist

### Backend Changes
- ✅ Fix `StockAdjust` method - thêm `userProvidedNewPrice` flag
- ✅ Fix `StockIn` method - thêm `userProvidedNewPrice` flag
- ✅ Restart backend server

### Frontend (Đã Fix Trước)
- ✅ Fix `adjustStock()` function - chỉ gửi price khi increase

### Testing
- ⏳ Test Case 1: Adjust tăng không nhập giá → Không tạo expense
- ⏳ Test Case 2: Adjust tăng có nhập giá → Tạo expense
- ⏳ Test Case 3: Stock in không nhập giá → Không tạo expense
- ⏳ Test Case 4: Stock in có nhập giá → Tạo expense

---

## 🎯 Kết Luận

### Trước Fix
- ❌ Adjust tăng không nhập giá → Vẫn tạo expense (SAI)
- ❌ Stock in không nhập giá → Vẫn tạo expense (SAI)

### Sau Fix
- ✅ Adjust tăng không nhập giá → KHÔNG tạo expense (ĐÚNG)
- ✅ Adjust tăng có nhập giá → Tạo expense (ĐÚNG)
- ✅ Stock in không nhập giá → KHÔNG tạo expense (ĐÚNG)
- ✅ Stock in có nhập giá → Tạo expense (ĐÚNG)

### Use Cases
1. **Mua hàng**: Nhập giá → Tạo expense ✅
2. **Được tặng**: Không nhập giá → Không tạo expense ✅
3. **Tìm thấy thêm**: Không nhập giá → Không tạo expense ✅
4. **Khuyến mãi**: Không nhập giá → Không tạo expense ✅

---

**Ngày fix**: 2026-02-07  
**Backend đã restart**: ✅ YES  
**Sẵn sàng test**: ✅ YES

---

**Tạo bởi**: Kiro AI
