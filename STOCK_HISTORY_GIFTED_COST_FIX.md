# Fix: Lịch Sử Ghi Sai Chi Phí Khi Được Tặng

## 🎯 Vấn Đề

Khi điều chỉnh tăng số lượng mà KHÔNG nhập giá (được tặng), lịch sử vẫn ghi chi phí dựa trên giá hiện tại.

### Use Case Lỗi

```
Hiện tại: 1kg cà phê @ 100,000đ/kg
User điều chỉnh: Tăng 2kg (được tặng, không nhập giá)

❌ Lịch sử ghi SAI:
- Đơn giá: 100,000đ
- Tổng chi phí: 2kg × 100,000đ = 200,000đ

✅ Lịch sử đúng phải là:
- Đơn giá: 0đ
- Tổng chi phí: 0đ
```

---

## 🔍 Nguyên Nhân

### Backend Logic Cũ (SAI)

**File**: `backend/application/services/ingredient.go`

#### StockAdjust Method
```go
// Determine cost per unit for this transaction
costPerUnit := item.CostPerUnit // ❌ Default to current price (100k)

// ... logic ...

// Create stock history record
history := &ingredient.StockHistory{
    CostPerUnit:  costPerUnit,              // ❌ 100k (giá hiện tại)
    TotalCost:    quantityDiff * costPerUnit, // ❌ 2kg × 100k = 200k
}
```

**Vấn đề**: Khi user không nhập giá, `costPerUnit` vẫn = `item.CostPerUnit` (giá hiện tại), nên lịch sử ghi sai.

#### StockIn Method
```go
costPerUnit := req.CostPerUnit
if costPerUnit <= 0 {
    costPerUnit = item.CostPerUnit  // ❌ Dùng giá hiện tại
}

history := &ingredient.StockHistory{
    CostPerUnit:  costPerUnit,              // ❌ 100k
    TotalCost:    req.Quantity * costPerUnit, // ❌ 2kg × 100k = 200k
}
```

---

## ✅ Giải Pháp

### Logic Đúng

**Lịch sử phải ghi đúng chi phí thực tế**:
- Nếu user NHẬP GIÁ → Ghi giá đó
- Nếu user KHÔNG NHẬP GIÁ → Ghi 0 (vì được tặng, không có chi phí)

### Backend Fix

#### 1. StockAdjust Method

```go
// Determine cost per unit for this transaction
costPerUnit := float64(0)       // ✅ Default to 0 (no cost)
userProvidedNewPrice := false

// Only set costPerUnit if user provided a new price
if quantityDiff > 0 && req.CostPerUnit > 0 && req.CostPerUnit != item.CostPerUnit {
    // Weighted average for the increase
    oldValue := beforeQty * item.CostPerUnit
    newValue := quantityDiff * req.CostPerUnit
    item.CostPerUnit = (oldValue + newValue) / afterQty
    costPerUnit = req.CostPerUnit // ✅ Use new price for history
    userProvidedNewPrice = true
}

// Create stock history record
// For history: only record cost if user provided a new price
// If gifted/found (no price), record as 0 cost
history := &ingredient.StockHistory{
    CostPerUnit:  costPerUnit,  // ✅ 0 if gifted, req.CostPerUnit if purchased
    TotalCost:    quantityDiff * costPerUnit,  // ✅ 0 if gifted
}
```

#### 2. StockIn Method

```go
// Determine cost per unit for this transaction
costPerUnit := req.CostPerUnit  // ✅ Use provided price, or 0 if not provided
userProvidedNewPrice := req.CostPerUnit > 0

// Calculate weighted average cost ONLY when new price is provided
if req.CostPerUnit > 0 && req.CostPerUnit != item.CostPerUnit {
    oldValue := beforeQty * item.CostPerUnit
    newValue := req.Quantity * req.CostPerUnit
    item.CostPerUnit = (oldValue + newValue) / afterQty
}

// Create stock history record
// For history: only record cost if user provided a price
// If gifted/found (no price), record as 0 cost
history := &ingredient.StockHistory{
    CostPerUnit:  costPerUnit,  // ✅ 0 if gifted, req.CostPerUnit if purchased
    TotalCost:    req.Quantity * costPerUnit,  // ✅ 0 if gifted
}
```

---

## 📊 Logic Tổng Hợp

### Lịch Sử Ghi Nhận

| Thao tác | Nhập giá? | req.CostPerUnit | costPerUnit (history) | TotalCost (history) |
|----------|-----------|-----------------|----------------------|---------------------|
| Adjust tăng | ❌ Không | 0 | 0 | 0 |
| Adjust tăng | ✅ Có (30k) | 30000 | 30000 | qty × 30000 |
| Stock in | ❌ Không | 0 | 0 | 0 |
| Stock in | ✅ Có (30k) | 30000 | 30000 | qty × 30000 |

### Use Case Examples

#### Use Case 1: Được Tặng
```
Hiện tại: 1kg @ 100k
Điều chỉnh: +2kg, không nhập giá
Lý do: "Được tặng từ nhà cung cấp"

Lịch sử ghi:
- Số lượng: +2kg
- Đơn giá: 0đ ✅
- Tổng chi phí: 0đ ✅
- Lý do: "Được tặng từ nhà cung cấp"
```

#### Use Case 2: Mua Hàng
```
Hiện tại: 1kg @ 100k
Điều chỉnh: +2kg, nhập giá 120k
Lý do: "Mua thêm"

Lịch sử ghi:
- Số lượng: +2kg
- Đơn giá: 120,000đ ✅
- Tổng chi phí: 240,000đ ✅
- Lý do: "Mua thêm"
```

#### Use Case 3: Tìm Thấy Thêm
```
Hiện tại: 5kg @ 50k
Điều chỉnh: +1kg, không nhập giá
Lý do: "Tìm thấy thêm trong kho"

Lịch sử ghi:
- Số lượng: +1kg
- Đơn giá: 0đ ✅
- Tổng chi phí: 0đ ✅
- Lý do: "Tìm thấy thêm trong kho"
```

---

## 🧪 Test Cases

### ✅ Test Case 1: Adjust Tăng Không Nhập Giá

**Bước thực hiện**:
1. Nguyên liệu: Cà phê (1kg @ 100,000đ)
2. Click "📦 Điều chỉnh" → "Điều chỉnh"
3. Nhập số lượng: 3kg (tăng 2kg)
4. **KHÔNG nhập giá** (để 0)
5. Lý do: "Được tặng từ nhà cung cấp"
6. Click "Xác nhận"
7. Click "📊 Lịch sử" để xem

**Kết quả mong đợi**:
- ✅ Lịch sử mới nhất:
  - Số lượng: +2kg
  - Đơn giá: **0đ** (không phải 100,000đ)
  - Tổng chi phí: **0đ** (không phải 200,000đ)
  - Lý do: "Được tặng từ nhà cung cấp"

### ✅ Test Case 2: Adjust Tăng Có Nhập Giá

**Bước thực hiện**:
1. Nguyên liệu: Cà phê (3kg @ 100,000đ)
2. Click "📦 Điều chỉnh" → "Điều chỉnh"
3. Nhập số lượng: 5kg (tăng 2kg)
4. **Nhập giá: 120,000đ**
5. Lý do: "Mua thêm"
6. Click "Xác nhận"
7. Click "📊 Lịch sử"

**Kết quả mong đợi**:
- ✅ Lịch sử mới nhất:
  - Số lượng: +2kg
  - Đơn giá: **120,000đ**
  - Tổng chi phí: **240,000đ** (2kg × 120,000đ)
  - Lý do: "Mua thêm"

### ✅ Test Case 3: Stock In Không Nhập Giá

**Bước thực hiện**:
1. Nguyên liệu: Sữa (5 chai @ 10,000đ)
2. Click "📦 Điều chỉnh" → "Nhập kho"
3. Nhập số lượng: 3 chai
4. **KHÔNG nhập giá**
5. Lý do: "Khuyến mãi"
6. Click "Xác nhận"
7. Click "📊 Lịch sử"

**Kết quả mong đợi**:
- ✅ Lịch sử mới nhất:
  - Số lượng: +3 chai
  - Đơn giá: **0đ**
  - Tổng chi phí: **0đ**
  - Lý do: "Khuyến mãi"

---

## 📁 Files Thay Đổi

### Backend
- `backend/application/services/ingredient.go`
  - **StockAdjust method** (Line ~265): Change `costPerUnit` default from `item.CostPerUnit` to `0`
  - **StockIn method** (Line ~160): Remove fallback to `item.CostPerUnit`, keep as `req.CostPerUnit` (0 if not provided)

---

## 🔄 Impact

### Trước Fix
```
Adjust tăng 2kg, không nhập giá:
Lịch sử: Đơn giá 100k, Tổng 200k ❌ SAI
Expense: Không tạo ✅ ĐÚNG (đã fix trước)
```

### Sau Fix
```
Adjust tăng 2kg, không nhập giá:
Lịch sử: Đơn giá 0đ, Tổng 0đ ✅ ĐÚNG
Expense: Không tạo ✅ ĐÚNG
```

---

## ✅ Checklist

### Backend Changes
- ✅ Fix `StockAdjust()` - default `costPerUnit = 0`
- ✅ Fix `StockIn()` - không fallback to `item.CostPerUnit`
- ✅ Restart backend server

### Testing
- ⏳ Test Case 1: Adjust tăng không nhập giá → Lịch sử ghi 0đ
- ⏳ Test Case 2: Adjust tăng có nhập giá → Lịch sử ghi đúng giá
- ⏳ Test Case 3: Stock in không nhập giá → Lịch sử ghi 0đ

---

## 🎯 Kết Luận

### Vấn Đề
Lịch sử ghi sai chi phí khi được tặng (ghi giá hiện tại thay vì 0).

### Giải Pháp
- Default `costPerUnit = 0` (không có chi phí)
- Chỉ set `costPerUnit = req.CostPerUnit` khi user thực sự nhập giá

### Kết Quả
- ✅ Lịch sử ghi đúng: 0đ khi được tặng
- ✅ Lịch sử ghi đúng: giá thực khi mua hàng
- ✅ Expense logic đúng: chỉ tạo khi mua hàng

---

**Ngày fix**: 2026-02-07  
**Backend đã restart**: ✅ YES  
**Sẵn sàng test**: ✅ YES

---

**Tạo bởi**: Kiro AI
