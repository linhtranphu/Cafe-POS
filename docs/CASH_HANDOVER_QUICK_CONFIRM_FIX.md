# Cash Handover - Fix Quick Confirm 500 Error

## 🐛 Vấn đề

Khi cashier click nút "Quick Confirm" (✅ hoặc ❌) trong danh sách pending handovers, hệ thống báo lỗi:

```
POST /api/cash-handovers/69837e580bac52676af309d0/quick-confirm
Status: 500 (Internal Server Error)
```

## 🔍 Nguyên nhân

### Code sai trong QuickConfirm handler:

```go
// ❌ SAI - GetPendingHandover nhận shiftID, không phải handoverID
h_obj, err := h.handoverService.GetPendingHandover(c.Request.Context(), handoverOID)
```

**Vấn đề:**
- `GetPendingHandover(shiftID)` - Tìm pending handover của một shift
- Nhưng đang truyền `handoverID` → Không tìm thấy → Error

### Signature của GetPendingHandover:

```go
// GetPendingHandover gets pending handover for a shift
func (s *CashHandoverService) GetPendingHandover(
    ctx context.Context, 
    shiftID primitive.ObjectID  // ← Nhận shiftID, không phải handoverID
) (*handover.CashHandover, error) {
    return s.handoverRepo.FindPendingByWaiterShift(ctx, shiftID)
}
```

## ✅ Giải pháp

### 1. Thêm method mới: GetHandoverByID

**File:** `backend/application/services/cash_handover_service.go`

```go
// GetHandoverByID gets handover by ID
func (s *CashHandoverService) GetHandoverByID(
    ctx context.Context, 
    handoverID primitive.ObjectID
) (*handover.CashHandover, error) {
    return s.handoverRepo.FindByID(ctx, handoverID)
}
```

### 2. Sửa QuickConfirm handler

**File:** `backend/interfaces/http/cash_handover_handler.go`

```go
// QuickConfirm quickly confirms or rejects without detailed reconciliation (cashier)
func (h *CashHandoverHandler) QuickConfirm(c *gin.Context) {
    handoverID := c.Param("id")
    handoverOID, err := primitive.ObjectIDFromHex(handoverID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid handover ID"})
        return
    }

    var req struct {
        Status handover.HandoverStatus `json:"status" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID, _ := c.Get("user_id")

    // ✅ ĐÚNG - Get handover by ID
    h_obj, err := h.handoverService.GetHandoverByID(c.Request.Context(), handoverOID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // For quick confirm, assume declared = actual (no discrepancy)
    confirmReq := &handover.ConfirmHandoverRequest{
        ActualAmount: h_obj.DeclaredAmount,
        Status:       req.Status,
        CashierNote:  "Quick confirm",
    }

    err = h.handoverService.ConfirmHandoverWithReconciliation(
        c.Request.Context(),
        handoverOID,
        confirmReq,
        userID.(string),
    )
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "handover confirmed successfully"})
}
```

## 📊 Flow sau khi fix

### Quick Confirm Flow:

```
1. Cashier click ✅ hoặc ❌ trong danh sách
   ↓
2. Frontend gửi:
   POST /api/cash-handovers/{handover_id}/quick-confirm
   Body: { "status": "CONFIRMED" hoặc "REJECTED" }
   ↓
3. Backend QuickConfirm handler:
   ├─ Parse handover_id từ URL
   ├─ ✅ GetHandoverByID(handover_id) → Lấy handover object
   ├─ Tạo ConfirmHandoverRequest:
   │  ├─ actual_amount = declared_amount (giả định không chênh lệch)
   │  ├─ status = CONFIRMED/REJECTED
   │  └─ cashier_note = "Quick confirm"
   └─ ConfirmHandoverWithReconciliation()
   ↓
4. Success!
   ├─ Status updated
   ├─ Cash amounts updated (nếu CONFIRMED)
   └─ Frontend refresh danh sách
```

## 🔄 So sánh 2 methods

### GetPendingHandover (cho waiter):

```go
// Use case: Waiter muốn xem pending handover của shift mình
GET /shifts/{shift_id}/pending-handover

func GetPendingHandover(shiftID) {
    return FindPendingByWaiterShift(shiftID)
}
```

### GetHandoverByID (cho cashier):

```go
// Use case: Cashier muốn xem/xử lý một handover cụ thể
POST /cash-handovers/{handover_id}/quick-confirm

func GetHandoverByID(handoverID) {
    return FindByID(handoverID)
}
```

## 🧪 Test Cases

### Test 1: Quick Confirm ✅

```bash
POST /api/cash-handovers/69837e580bac52676af309d0/quick-confirm
Body: { "status": "CONFIRMED" }

Expected:
- ✅ Status 200
- Handover status = CONFIRMED
- Actual amount = Declared amount
- Cash amounts updated
```

### Test 2: Quick Reject ❌

```bash
POST /api/cash-handovers/69837e580bac52676af309d0/quick-confirm
Body: { "status": "REJECTED" }

Expected:
- ✅ Status 200
- Handover status = REJECTED
- Cashier note = "Quick confirm"
- Cash amounts NOT updated
```

### Test 3: Invalid Handover ID

```bash
POST /api/cash-handovers/invalid-id/quick-confirm
Body: { "status": "CONFIRMED" }

Expected:
- ❌ Status 400
- Error: "invalid handover ID"
```

### Test 4: Handover Not Found

```bash
POST /api/cash-handovers/000000000000000000000000/quick-confirm
Body: { "status": "CONFIRMED" }

Expected:
- ❌ Status 400
- Error: "handover not found" (from repository)
```

## 🔧 Files Changed

1. **backend/application/services/cash_handover_service.go**
   - Thêm method `GetHandoverByID(handoverID)`

2. **backend/interfaces/http/cash_handover_handler.go**
   - Sửa `QuickConfirm` handler để dùng `GetHandoverByID` thay vì `GetPendingHandover`

## ✅ Checklist

- [x] Thêm GetHandoverByID method
- [x] Sửa QuickConfirm handler
- [x] Verify no syntax errors
- [x] Test quick confirm ✅
- [x] Test quick reject ❌
- [x] Document changes

## 🚀 Deployment

```bash
# Rebuild backend (BẮT BUỘC)
cd backend
go build -o cafe-pos-server

# Restart services
docker-compose restart backend
```

## 💡 Bài học

### Vấn đề phổ biến: Nhầm lẫn giữa các ID types

```go
// ❌ SAI - Nhầm lẫn ID types
func ProcessHandover(handoverID) {
    shift := GetShift(handoverID)  // Wrong! handoverID ≠ shiftID
}

// ✅ ĐÚNG - Rõ ràng về ID types
func ProcessHandover(handoverID) {
    handover := GetHandover(handoverID)
    shift := GetShift(handover.ShiftID)
}
```

### Best Practice: Naming conventions

```go
// ✅ GOOD - Tên method rõ ràng về parameter type
func GetHandoverByID(handoverID)      // Nhận handover ID
func GetPendingHandover(shiftID)      // Nhận shift ID
func GetHandoversByShift(shiftID)     // Nhận shift ID
func GetHandoversByCashier(cashierID) // Nhận cashier ID
```

### Type Safety

```go
// ✅ BETTER - Sử dụng type aliases để tránh nhầm lẫn
type HandoverID primitive.ObjectID
type ShiftID primitive.ObjectID
type CashierID primitive.ObjectID

func GetHandoverByID(id HandoverID) (*Handover, error)
func GetPendingHandover(id ShiftID) (*Handover, error)
```

## 📚 Related Issues

- Tương tự có thể xảy ra với các API khác có multiple ID types
- Nên review tất cả handlers để đảm bảo đúng ID type

---

**Status:** ✅ FIXED  
**Date:** 2026-02-05  
**Impact:** Critical - Blocking quick confirm functionality  
**Root Cause:** Wrong method call with incorrect parameter type
