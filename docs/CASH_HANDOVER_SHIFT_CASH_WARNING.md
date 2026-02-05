# Cash Handover - Shift Cash Mismatch Warning ⚠️

## 🎯 Vấn đề

Khi waiter tạo handover request, họ có thể khai báo số tiền **không khớp** với số tiền còn lại thực tế trong ca (`remaining_cash`). Cashier cần được cảnh báo về điều này **trước khi** xác nhận bàn giao.

### Ví dụ:

```
Waiter Shift:
- Start Cash: 500,000₫
- Cash Revenue: 200,000₫
- Handed Over: 100,000₫
- Remaining Cash: 600,000₫ (500k + 200k - 100k)

Waiter tạo handover:
- Declared Amount: 400,000₫ ❌ (Không khớp với 600k)

→ Cashier cần được cảnh báo!
```

## 🔍 Nguyên nhân

1. **Backend không trả về shift info** - API `/cash-handovers/pending` chỉ trả về handover data, không có thông tin về shift
2. **Frontend không so sánh** - Không có logic để so sánh `declared_amount` với `shift_remaining_cash`
3. **Không có cảnh báo UI** - Không có component hiển thị cảnh báo khi phát hiện mismatch

## ✅ Giải pháp

### 1. Backend Changes

#### File: `backend/application/services/cash_handover_service.go`

**Thêm DTO mới:**

```go
// HandoverWithShiftInfo represents a handover with additional shift information
type HandoverWithShiftInfo struct {
	*handover.CashHandover
	ShiftRemainingCash float64 `json:"shift_remaining_cash"`
	ShiftCurrentCash   float64 `json:"shift_current_cash"`
	ShiftCashRevenue   float64 `json:"shift_cash_revenue"`
}
```

**Thêm methods mới:**

```go
// GetPendingByCashierWithShiftInfo gets pending handovers with shift information
func (s *CashHandoverService) GetPendingByCashierWithShiftInfo(ctx context.Context, cashierID primitive.ObjectID) ([]*HandoverWithShiftInfo, error) {
	handovers, err := s.handoverRepo.FindPendingByCashier(ctx, cashierID)
	if err != nil {
		return nil, err
	}

	result := make([]*HandoverWithShiftInfo, 0, len(handovers))
	for _, h := range handovers {
		// Get waiter shift info
		shift, err := s.shiftRepo.FindByID(ctx, h.WaiterShiftID)
		if err != nil {
			// If shift not found, still include handover but with zero values
			result = append(result, &HandoverWithShiftInfo{
				CashHandover:       h,
				ShiftRemainingCash: 0,
				ShiftCurrentCash:   0,
				ShiftCashRevenue:   0,
			})
			continue
		}

		result = append(result, &HandoverWithShiftInfo{
			CashHandover:       h,
			ShiftRemainingCash: shift.RemainingCash,
			ShiftCurrentCash:   shift.CurrentCash,
			ShiftCashRevenue:   shift.CashRevenue,
		})
	}

	return result, nil
}

// GetAllPendingWithShiftInfo gets all pending handovers with shift information
func (s *CashHandoverService) GetAllPendingWithShiftInfo(ctx context.Context) ([]*HandoverWithShiftInfo, error) {
	// Similar implementation
}
```

#### File: `backend/interfaces/http/cash_handover_handler.go`

**Cập nhật handlers:**

```go
// GetPendingHandovers gets all pending handovers for current cashier
func (h *CashHandoverHandler) GetPendingHandovers(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// ✅ Use new method with shift info
	handovers, err := h.handoverService.GetPendingByCashierWithShiftInfo(c.Request.Context(), userOID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, handovers)
}
```

### 2. Frontend Changes

#### File: `frontend/src/views/CashierHandoverView.vue`

**Thêm computed property:**

```javascript
// Warning for declared amount vs shift cash
const shiftCashWarning = computed(() => {
  if (!selectedHandover.value) return null
  
  const handover = selectedHandover.value
  
  // Check if declared amount matches expected cash from shift
  if (handover.shift_remaining_cash !== undefined) {
    const difference = handover.declared_amount - handover.shift_remaining_cash
    
    if (difference !== 0) {
      return {
        type: difference > 0 ? 'OVER_DECLARED' : 'UNDER_DECLARED',
        difference: Math.abs(difference),
        message: difference > 0 
          ? `Waiter khai báo nhiều hơn tiền còn lại trong ca (${formatPrice(Math.abs(difference))})`
          : `Waiter khai báo ít hơn tiền còn lại trong ca (${formatPrice(Math.abs(difference))})`
      }
    }
  }
  
  return null
})
```

**Thêm cảnh báo trong danh sách:**

```vue
<!-- Shift Cash Warning (in list) -->
<div v-if="handover.shift_remaining_cash !== undefined && handover.declared_amount !== handover.shift_remaining_cash"
  class="mb-3 p-3 rounded-lg border-2"
  :class="handover.declared_amount > handover.shift_remaining_cash ? 'bg-orange-50 border-orange-300' : 'bg-yellow-50 border-yellow-300'">
  <div class="flex items-start gap-2">
    <span class="text-lg">⚠️</span>
    <div class="flex-1">
      <p class="text-xs font-medium"
        :class="handover.declared_amount > handover.shift_remaining_cash ? 'text-orange-800' : 'text-yellow-800'">
        {{ handover.declared_amount > handover.shift_remaining_cash 
          ? 'Khai báo nhiều hơn tiền còn lại trong ca' 
          : 'Khai báo ít hơn tiền còn lại trong ca' }}
      </p>
      <p class="text-xs mt-1"
        :class="handover.declared_amount > handover.shift_remaining_cash ? 'text-orange-600' : 'text-yellow-600'">
        Tiền còn lại: {{ formatPrice(handover.shift_remaining_cash) }} | 
        Chênh: {{ formatPrice(Math.abs(handover.declared_amount - handover.shift_remaining_cash)) }}
      </p>
    </div>
  </div>
</div>
```

**Thêm cảnh báo trong modal:**

```vue
<!-- Shift Cash Warning -->
<div v-if="shiftCashWarning" class="mt-3 p-3 rounded-lg border-2"
  :class="shiftCashWarning.type === 'OVER_DECLARED' ? 'bg-orange-50 border-orange-300' : 'bg-yellow-50 border-yellow-300'">
  <div class="flex items-start gap-2">
    <span class="text-xl">⚠️</span>
    <div class="flex-1">
      <p class="text-sm font-medium" 
        :class="shiftCashWarning.type === 'OVER_DECLARED' ? 'text-orange-800' : 'text-yellow-800'">
        {{ shiftCashWarning.message }}
      </p>
      <p class="text-xs mt-1" 
        :class="shiftCashWarning.type === 'OVER_DECLARED' ? 'text-orange-600' : 'text-yellow-600'">
        Tiền còn lại trong ca: {{ formatPrice(selectedHandover?.shift_remaining_cash || 0) }}
      </p>
    </div>
  </div>
</div>
```

## 🎨 UI/UX

### Màu sắc cảnh báo:

| Tình huống | Màu | Ý nghĩa |
|------------|-----|---------|
| **OVER_DECLARED** | 🟠 Cam | Waiter khai báo > Remaining Cash (Nghi ngờ gian lận) |
| **UNDER_DECLARED** | 🟡 Vàng | Waiter khai báo < Remaining Cash (Có thể quên tiền) |

### Vị trí hiển thị:

1. **Trong danh sách pending handovers** - Cảnh báo nhỏ dưới thông tin handover
2. **Trong modal xác nhận** - Cảnh báo lớn trong phần summary

## 📊 Scenarios

### Scenario 1: Khai báo đúng ✅

```
Shift Remaining Cash: 600,000₫
Declared Amount: 600,000₫
→ Không có cảnh báo
```

### Scenario 2: Khai báo ít hơn 🟡

```
Shift Remaining Cash: 600,000₫
Declared Amount: 400,000₫
Difference: -200,000₫

→ Cảnh báo vàng:
"Waiter khai báo ít hơn tiền còn lại trong ca (200,000₫)"
"Tiền còn lại: 600,000₫ | Chênh: 200,000₫"

Lý do có thể:
- Waiter quên một phần tiền
- Waiter giữ lại tiền để tiếp tục làm việc
- Handover type = PARTIAL (bàn giao một phần)
```

### Scenario 3: Khai báo nhiều hơn 🟠

```
Shift Remaining Cash: 600,000₫
Declared Amount: 800,000₫
Difference: +200,000₫

→ Cảnh báo cam:
"Waiter khai báo nhiều hơn tiền còn lại trong ca (200,000₫)"
"Tiền còn lại: 600,000₫ | Chênh: 200,000₫"

Lý do có thể:
- Waiter nhầm lẫn
- Waiter cố tình khai báo sai
- Lỗi hệ thống tracking cash
- Waiter nhận tiền từ nguồn khác (không hợp lệ)

⚠️ CẦN KIỂM TRA KỸ!
```

## 🔍 Cashier Actions

### Khi thấy cảnh báo UNDER_DECLARED (Vàng):

1. ✅ Hỏi waiter: "Bạn có chắc chỉ bàn giao số tiền này?"
2. ✅ Kiểm tra xem có phải handover type = PARTIAL không
3. ✅ Xác nhận waiter có giữ lại tiền để tiếp tục làm việc
4. ✅ Nếu OK → Xác nhận bình thường

### Khi thấy cảnh báo OVER_DECLARED (Cam):

1. ⚠️ **DỪNG LẠI** - Không xác nhận ngay
2. ⚠️ Gọi waiter đến kiểm tra
3. ⚠️ Đếm lại tiền cẩn thận
4. ⚠️ Kiểm tra lịch sử orders trong ca
5. ⚠️ Nếu không khớp → Từ chối và báo manager
6. ⚠️ Nếu khớp → Có thể lỗi hệ thống, báo manager

## 🔄 Flow hoàn chỉnh

```
1. Waiter tạo handover request
   ├─ Declared Amount: 400,000₫
   └─ Shift Remaining Cash: 600,000₫
   
2. Backend lưu handover
   └─ Status: PENDING
   
3. Cashier vào /cashier/handovers
   └─ Backend trả về handover + shift info
   
4. Frontend hiển thị danh sách
   ├─ Tính: difference = 400k - 600k = -200k
   ├─ Type: UNDER_DECLARED
   └─ 🟡 Hiển thị cảnh báo vàng
   
5. Cashier click "Xác nhận"
   └─ Modal mở ra với cảnh báo lớn hơn
   
6. Cashier kiểm tra với waiter
   └─ Xác nhận waiter giữ lại 200k để tiếp tục
   
7. Cashier nhập actual_amount: 400,000₫
   └─ Không có discrepancy (declared = actual)
   
8. Cashier xác nhận
   └─ Handover completed
```

## 🧪 Testing

### Test Case 1: Khai báo đúng
```bash
# Shift remaining: 600k
# Declared: 600k
# Expected: No warning
```

### Test Case 2: Khai báo ít (PARTIAL)
```bash
# Shift remaining: 600k
# Declared: 400k (PARTIAL handover)
# Expected: Yellow warning
```

### Test Case 3: Khai báo nhiều
```bash
# Shift remaining: 600k
# Declared: 800k
# Expected: Orange warning (ALERT!)
```

### Test Case 4: END_SHIFT type
```bash
# Shift remaining: 600k
# Declared: 600k (END_SHIFT)
# Expected: No warning (should match)
```

## 📝 Notes

### Khi nào cảnh báo là bình thường?

- **PARTIAL handover** - Waiter chỉ bàn giao một phần, giữ lại tiền để tiếp tục
- **Multiple handovers** - Waiter đã bàn giao nhiều lần trong ca

### Khi nào cảnh báo là nghiêm trọng?

- **END_SHIFT handover** - Phải bàn giao toàn bộ, không được chênh lệch
- **OVER_DECLARED** - Khai báo nhiều hơn remaining cash (nghi ngờ gian lận)
- **Large difference** - Chênh lệch quá lớn (> 100k)

## 🔮 Future Enhancements

1. **Auto-suggest declared amount** - Gợi ý số tiền = remaining cash
2. **Handover history** - Hiển thị lịch sử handover trước đó của waiter
3. **Cash flow visualization** - Biểu đồ dòng tiền trong ca
4. **Manager approval** - Yêu cầu manager approval cho OVER_DECLARED
5. **Audit log** - Log chi tiết mọi cảnh báo và actions

## ✅ Checklist

- [x] Thêm HandoverWithShiftInfo DTO
- [x] Thêm GetPendingByCashierWithShiftInfo method
- [x] Thêm GetAllPendingWithShiftInfo method
- [x] Cập nhật handlers
- [x] Thêm shiftCashWarning computed property
- [x] Thêm cảnh báo trong danh sách
- [x] Thêm cảnh báo trong modal
- [x] Test với các scenarios
- [x] Viết documentation

## 🚀 Deployment

```bash
# Rebuild backend
cd backend
go build -o cafe-pos-server

# Restart backend
docker-compose restart backend

# Frontend không cần rebuild (chỉ cần refresh)
```

## 📚 Related Documents

- [CASH_HANDOVER_DISCREPANCY_FIX.md](./CASH_HANDOVER_DISCREPANCY_FIX.md) - Discrepancy warning giữa declared và actual
- [CASH_HANDOVER_COMPLETE_SUMMARY.md](./CASH_HANDOVER_COMPLETE_SUMMARY.md) - Tổng quan tính năng
