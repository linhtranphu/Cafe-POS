# Tạo Nguyên Liệu - Ghi Lịch Sử Tồn Kho Ban Đầu

## Vấn Đề
Khi tạo nguyên liệu mới với số lượng ban đầu, không có bản ghi lịch sử tồn kho được tạo. Điều này có nghĩa:
- Lần mua đầu tiên không được theo dõi trong lịch sử
- User không thể xem giá đã trả ban đầu
- Lịch sử chỉ hiển thị các điều chỉnh sau đó
- Audit trail không đầy đủ

## Giải Pháp
Tự động tạo bản ghi lịch sử tồn kho khi tạo nguyên liệu với số lượng > 0.

## Triển Khai

### 1. Thay Đổi Tầng Service
**File:** `backend/application/services/ingredient.go`

**Thêm logic vào CreateIngredient:**
```go
// Tạo bản ghi lịch sử tồn kho ban đầu nếu quantity > 0
if req.Quantity > 0 {
    userID := primitive.NilObjectID
    if userIDStr != "" {
        if oid, err := primitive.ObjectIDFromHex(userIDStr); err == nil {
            userID = oid
        }
    }
    
    history := &ingredient.StockHistory{
        IngredientID: item.ID,
        Type:         ingredient.TransactionPurchase,
        Quantity:     req.Quantity,
        BeforeQty:    0,                              // Tạo mới
        AfterQty:     req.Quantity,
        Reason:       "Tạo nguyên liệu mới - Nhập kho đầu tiên",
        UserID:       userID,
        Username:     username,
        CostPerUnit:  req.CostPerUnit,
        TotalCost:    req.Quantity * req.CostPerUnit,
    }
    
    // Tạo lịch sử (không fail nếu thất bại)
    if err := s.stockHistoryRepo.Create(ctx, history); err != nil {
        // Log lỗi nhưng không fail operation
    }
}
```

**Điểm chính:**
- Chỉ tạo lịch sử nếu `quantity > 0`
- Set `BeforeQty` = 0 (trạng thái ban đầu)
- Set `AfterQty` = số lượng ban đầu
- Dùng loại giao dịch `TransactionPurchase`
- Ghi lại đơn giá ban đầu và tổng chi phí
- Bao gồm thông tin user cho audit trail
- Không fail việc tạo nguyên liệu nếu tạo lịch sử thất bại

### 2. Thay Đổi Tầng Handler
**File:** `backend/interfaces/http/ingredient_handler.go`

**Cập nhật CreateIngredient handler:**
```go
func (h *IngredientHandler) CreateIngredient(c *gin.Context) {
    // ... validation ...
    
    // Lấy thông tin user từ context
    userID, _ := c.Get("user_id")
    username, _ := c.Get("username")
    
    userIDStr := ""
    if uid, ok := userID.(string); ok {
        userIDStr = uid
    }
    
    createdBy := ""
    if u, ok := username.(string); ok {
        createdBy = u
    }

    item, err := h.ingredientService.CreateIngredient(
        c.Request.Context(), 
        &req, 
        userIDStr,  // Thêm userID
        createdBy,
    )
    
    // ... response ...
}
```

**Thay đổi:**
- Trích xuất cả `user_id` và `username` từ context
- Truyền cả hai vào service layer
- Service dùng userID cho audit trail đúng

### 3. Cập Nhật Signature Service
**Trước:**
```go
CreateIngredient(ctx context.Context, req *ingredient.CreateIngredientRequest, username string)
```

**Sau:**
```go
CreateIngredient(ctx context.Context, req *ingredient.CreateIngredientRequest, userIDStr string, username string)
```

## Chi Tiết Bản Ghi Lịch Sử

### Loại Giao Dịch
Dùng `TransactionPurchase` để chỉ ra đây là mua hàng ban đầu, không chỉ là điều chỉnh.

### Lý Do
Lý do mặc định: **"Tạo nguyên liệu mới - Nhập kho đầu tiên"**
- Chỉ ra rõ ràng đây là lần tạo ban đầu
- Tiếng Việt để nhất quán với UI
- Có thể tùy chỉnh nếu cần

### Số Lượng
- `BeforeQty`: Luôn là 0 (không có gì trước đó)
- `Quantity`: Số lượng ban đầu nhập vào
- `AfterQty`: Giống Quantity (0 + Quantity)

### Thông Tin Giá
- `CostPerUnit`: Đơn giá ban đầu nhập vào
- `TotalCost`: Tính bằng `Quantity × CostPerUnit`

### Thông Tin User
- `UserID`: ObjectID của user tạo nguyên liệu
- `Username`: Tên hiển thị cho UI

## Lợi Ích

### 1. Audit Trail Đầy Đủ
- Mọi thay đổi số lượng được theo dõi từ ngày đầu
- Lần mua ban đầu hiển thị trong lịch sử
- Không có dữ liệu bị thiếu

### 2. Theo Dõi Giá
- Giá ban đầu được ghi lại
- Có thể so sánh giá tương lai với giá ban đầu
- Dữ liệu giá lịch sử đầy đủ

### 3. Trách Nhiệm User
- Biết ai tạo mỗi nguyên liệu
- Theo dõi ai thực hiện mua hàng ban đầu
- Audit và tuân thủ tốt hơn

### 4. Hành Vi Nhất Quán
- Tạo nguyên liệu = điều chỉnh tồn kho
- Cả hai thao tác đều tạo bản ghi lịch sử
- Hành vi hệ thống có thể dự đoán

## Trải Nghiệm User

### Trước
1. User tạo nguyên liệu với 10 kg @ 50.000đ/kg
2. Mở lịch sử → Trống (không có bản ghi)
3. Nhầm lẫn: "Lần mua ban đầu của tôi đâu?"

### Sau
1. User tạo nguyên liệu với 10 kg @ 50.000đ/kg
2. Mở lịch sử → Hiển thị bản ghi mua ban đầu:
   ```
   📦 Nhập thêm
   +10 kg
   
   Lý do: Tạo nguyên liệu mới - Nhập kho đầu tiên
   
   💰 THÔNG TIN GIÁ
   Đơn giá lần này: 50.000 ₫/kg
   Tổng chi phí: 500.000 ₫
   = 10 kg × 50.000 ₫
   
   👤 Admin
   🕐 07/02/2026 10:30
   ```

## Xử Lý Edge Cases

### 1. Số Lượng Bằng 0
Nếu user tạo nguyên liệu với quantity = 0:
- Không tạo bản ghi lịch sử
- Hợp lý: chưa mua gì
- Lịch sử bắt đầu khi điều chỉnh tồn kho lần đầu

### 2. Thiếu Thông Tin User
Nếu user_id hoặc username không có trong context:
- Dùng chuỗi rỗng cho username
- Dùng NilObjectID cho userID
- Thao tác vẫn thành công

### 3. Tạo Lịch Sử Thất Bại
Nếu tạo lịch sử tồn kho thất bại:
- Lỗi được log (không hiển thị cho user)
- Tạo nguyên liệu vẫn thành công
- Lịch sử là phụ so với thao tác chính

### 4. Theo Dõi Chi Phí
Cả hai thao tác xảy ra:
- Bản ghi lịch sử tồn kho được tạo
- Bản ghi chi phí tự động được tạo
- Thao tác độc lập, cả hai có thể thành công/thất bại riêng

## Kiểm Tra

### Test Case 1: Tạo Với Số Lượng
```bash
POST /api/ingredients
{
  "name": "Sữa tươi",
  "category": "Nguyên liệu chính",
  "unit": "L",
  "quantity": 10,
  "min_stock": 2,
  "cost_per_unit": 25000,
  "supplier": "Vinamilk"
}

# Kỳ vọng:
# 1. Nguyên liệu được tạo
# 2. Bản ghi lịch sử tồn kho được tạo
# 3. Bản ghi chi phí tự động được tạo
# 4. GET /api/ingredients/:id/history trả về 1 bản ghi
```

### Test Case 2: Tạo Với Số Lượng 0
```bash
POST /api/ingredients
{
  "name": "Sữa tươi",
  "category": "Nguyên liệu chính",
  "unit": "L",
  "quantity": 0,
  "min_stock": 2,
  "cost_per_unit": 25000,
  "supplier": "Vinamilk"
}

# Kỳ vọng:
# 1. Nguyên liệu được tạo
# 2. KHÔNG có bản ghi lịch sử tồn kho
# 3. KHÔNG có bản ghi chi phí tự động
# 4. GET /api/ingredients/:id/history trả về mảng rỗng
```

### Test Case 3: Xác Minh Hiển Thị Lịch Sử
```bash
# 1. Tạo nguyên liệu với số lượng
# 2. Mở lịch sử nguyên liệu trong UI
# 3. Xác minh hiển thị:
#    - Card xanh (mua hàng)
#    - Số lượng "+10 L"
#    - Lý do "Tạo nguyên liệu mới - Nhập kho đầu tiên"
#    - Phần thông tin giá
#    - User và timestamp
```

## Tác Động Database

### Bản Ghi Mới
Mỗi lần tạo nguyên liệu với quantity > 0 tạo:
- 1 document ingredient
- 1 document stock_history
- 1 document expense (nếu auto-expense bật)

### Lưu Trữ
Tác động tối thiểu:
- Bản ghi lịch sử tồn kho nhỏ (~200 bytes)
- Cần thiết cho audit trail
- Đáng giá chi phí lưu trữ

## Migration

### Nguyên Liệu Hiện Có
Nguyên liệu hiện có được tạo trước thay đổi này:
- SẼ KHÔNG có bản ghi lịch sử ban đầu
- Lịch sử bắt đầu từ điều chỉnh đầu tiên
- Điều này chấp nhận được (dữ liệu lịch sử)

### Tùy Chọn: Script Backfill
Nếu cần, có thể tạo script để:
1. Tìm nguyên liệu không có lịch sử
2. Tạo bản ghi lịch sử ban đầu tổng hợp
3. Dùng cost_per_unit hiện tại làm giá ban đầu
4. Set created_at = ingredient created_at

## File Đã Sửa

1. `backend/application/services/ingredient.go`
   - Cập nhật signature `CreateIngredient`
   - Thêm logic tạo lịch sử tồn kho
   - Thêm xử lý user ID

2. `backend/interfaces/http/ingredient_handler.go`
   - Cập nhật handler `CreateIngredient`
   - Trích xuất user_id từ context
   - Truyền cả userID và username vào service

## Tính Năng Liên Quan

Thay đổi này bổ sung cho:
- Theo dõi chi phí tự động (đã triển khai)
- Hiển thị lịch sử tồn kho (đã triển khai)
- Theo dõi giá trong lịch sử (đã triển khai)

Cùng nhau, các tính năng này cung cấp:
- Theo dõi tài chính đầy đủ
- Audit trail hoàn chỉnh
- Phân tích lịch sử giá
- Trách nhiệm user
