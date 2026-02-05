# Cash Handover - Fix Reject Validation Error

## 🐛 Vấn đề

Khi cashier click nút "Từ chối" bàn giao, hệ thống báo lỗi:

```
Lỗi: Key: 'ConfirmHandoverRequest.ActualAmount' Error:Field validation for 'ActualAmount' failed on the 'required' tag
```

## 🔍 Nguyên nhân

### Backend validation sai:

```go
// ❌ SAI - actual_amount là required cho cả CONFIRMED và REJECTED
type ConfirmHandoverRequest struct {
    ActualAmount float64 `json:"actual_amount" binding:"required,gte=0"`
    Status       HandoverStatus `json:"status" binding:"required"`
    // ...
}
```

**Vấn đề:**
- Khi REJECT, không cần `actual_amount` (cashier không nhận tiền)
- Nhưng validation yêu cầu field này là `required`
- Frontend không gửi `actual_amount` khi reject → Lỗi validation

## ✅ Giải pháp

### 1. Bỏ `required` tag trong domain model

**File:** `backend/domain/handover/cash_handover.go`

```go
// ✅ ĐÚNG - actual_amount không bắt buộc
type ConfirmHandoverRequest struct {
    ActualAmount              float64            `json:"actual_amount" binding:"gte=0"`  // ← Bỏ required
    Status                    HandoverStatus     `json:"status" binding:"required"`
    CashierNote               string             `json:"cashier_note"`
    DiscrepancyReason         string             `json:"discrepancy_reason"`
    DiscrepancyResponsibility ResponsibilityType `json:"discrepancy_responsibility"`
}
```

### 2. Thêm validation logic trong service

**File:** `backend/application/services/cash_handover_service.go`

```go
func (s *CashHandoverService) ConfirmHandoverWithReconciliation(...) error {
    // ... existing code ...

    // ✅ Validate actual_amount chỉ khi CONFIRMED
    if req.Status == handover.StatusConfirmed && req.ActualAmount == 0 {
        return errors.New("actual_amount is required when confirming handover")
    }

    // ✅ Calculate discrepancy chỉ khi CONFIRMED
    var discrepancy float64
    if req.Status == handover.StatusConfirmed {
        discrepancy = req.ActualAmount - h.DeclaredAmount
    }

    // ✅ Update fields dựa trên status
    h.Status = req.Status
    h.CashierNote = req.CashierNote
    h.ConfirmedAt = &now
    h.UpdatedAt = now

    if req.Status == handover.StatusConfirmed {
        h.ActualAmount = req.ActualAmount
        h.Discrepancy = discrepancy
        h.ReconciledAt = &now
    }

    // ✅ Handle discrepancy chỉ khi CONFIRMED
    if req.Status == handover.StatusConfirmed && h.HasDiscrepancy() {
        // ... discrepancy logic ...
    }

    // ... rest of code ...
}
```

### 3. Cập nhật frontend validation

**File:** `frontend/src/views/CashierHandoverView.vue`

```javascript
const confirmHandover = async () => {
  try {
    const data = {
      status: confirmAction.value,
      cashier_note: confirmForm.value.cashier_note
    }
    
    // ✅ Chỉ gửi actual_amount khi CONFIRMED
    if (confirmAction.value === 'CONFIRMED') {
      if (!confirmForm.value.actual_amount || confirmForm.value.actual_amount === 0) {
        alert('Vui lòng nhập số tiền thực nhận')
        return
      }
      
      data.actual_amount = confirmForm.value.actual_amount
      
      // Add discrepancy info if exists
      if (hasDiscrepancy.value) {
        if (!confirmForm.value.discrepancy_reason || !confirmForm.value.discrepancy_responsibility) {
          alert('Vui lòng nhập đầy đủ thông tin chênh lệch')
          return
        }
        data.discrepancy_reason = confirmForm.value.discrepancy_reason
        data.discrepancy_responsibility = confirmForm.value.discrepancy_responsibility
      }
    } else {
      // ✅ Khi REJECTED, chỉ cần cashier_note
      if (!confirmForm.value.cashier_note || confirmForm.value.cashier_note.trim() === '') {
        alert('Vui lòng nhập lý do từ chối')
        return
      }
    }
    
    await cashierStore.confirmHandover(selectedHandover.value.id, data)
    // ... rest of code ...
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}
```

## 📊 Flow sau khi fix

### Scenario 1: CONFIRMED (Xác nhận)

```
1. Cashier nhập actual_amount
   ↓
2. Frontend validation:
   ├─ actual_amount > 0? ✅
   ├─ Có discrepancy? → Yêu cầu reason + responsibility
   └─ Gửi: { status: "CONFIRMED", actual_amount, cashier_note, ... }
   ↓
3. Backend validation:
   ├─ actual_amount > 0? ✅
   ├─ Calculate discrepancy
   ├─ Update actual_amount, discrepancy
   └─ Create discrepancy record (nếu có)
   ↓
4. Success!
```

### Scenario 2: REJECTED (Từ chối)

```
1. Cashier nhập lý do từ chối (cashier_note)
   ↓
2. Frontend validation:
   ├─ cashier_note không rỗng? ✅
   └─ Gửi: { status: "REJECTED", cashier_note }
       (KHÔNG gửi actual_amount)
   ↓
3. Backend validation:
   ├─ Status = REJECTED? ✅
   ├─ KHÔNG validate actual_amount
   ├─ KHÔNG calculate discrepancy
   └─ Chỉ update status + cashier_note
   ↓
4. Success!
```

## 🧪 Test Cases

### Test 1: Reject với lý do hợp lệ ✅

```
Input:
- Status: REJECTED
- Cashier Note: "Số tiền không đúng, yêu cầu waiter kiểm tra lại"
- Actual Amount: (không gửi)

Expected:
- ✅ Success
- Status = REJECTED
- CashierNote được lưu
- ActualAmount = 0
- Discrepancy = 0
```

### Test 2: Reject không có lý do ❌

```
Input:
- Status: REJECTED
- Cashier Note: ""
- Actual Amount: (không gửi)

Expected:
- ❌ Frontend validation error
- Alert: "Vui lòng nhập lý do từ chối"
```

### Test 3: Confirm với actual_amount = 0 ❌

```
Input:
- Status: CONFIRMED
- Actual Amount: 0
- Cashier Note: "OK"

Expected:
- ❌ Frontend validation error
- Alert: "Vui lòng nhập số tiền thực nhận"
```

### Test 4: Confirm với actual_amount hợp lệ ✅

```
Input:
- Status: CONFIRMED
- Actual Amount: 100000
- Cashier Note: "Đã nhận đủ"

Expected:
- ✅ Success
- Status = CONFIRMED
- ActualAmount = 100000
- Discrepancy = 0 (nếu declared = 100k)
```

## 📝 Validation Rules Summary

| Field | CONFIRMED | REJECTED |
|-------|-----------|----------|
| **status** | Required | Required |
| **actual_amount** | Required, > 0 | Not required, not sent |
| **cashier_note** | Optional | Required |
| **discrepancy_reason** | Required if discrepancy | Not applicable |
| **discrepancy_responsibility** | Required if discrepancy | Not applicable |

## 🔧 Files Changed

1. **backend/domain/handover/cash_handover.go**
   - Bỏ `required` tag từ `ActualAmount`

2. **backend/application/services/cash_handover_service.go**
   - Thêm validation logic cho CONFIRMED
   - Chỉ process discrepancy khi CONFIRMED
   - Chỉ update actual_amount khi CONFIRMED

3. **frontend/src/views/CashierHandoverView.vue**
   - Thêm validation riêng cho CONFIRMED và REJECTED
   - Không gửi actual_amount khi REJECTED
   - Validate cashier_note khi REJECTED

## ✅ Checklist

- [x] Bỏ `required` tag từ ActualAmount
- [x] Thêm validation logic trong service
- [x] Cập nhật frontend validation
- [x] Test REJECTED flow
- [x] Test CONFIRMED flow
- [x] Verify no syntax errors
- [x] Document changes

## 🚀 Deployment

```bash
# Rebuild backend (BẮT BUỘC)
cd backend
go build -o cafe-pos-server

# Restart services
docker-compose restart backend

# Frontend tự động reload
```

## 📚 Related Issues

- Vấn đề tương tự có thể xảy ra với các API khác có conditional validation
- Nên review tất cả request structs để đảm bảo validation đúng

## 💡 Best Practice

**Nguyên tắc:** Validation nên phụ thuộc vào context/state

```go
// ❌ SAI - Validation cứng nhắc
type Request struct {
    Field1 string `binding:"required"`  // Always required
    Field2 int    `binding:"required"`  // Always required
}

// ✅ ĐÚNG - Validation linh hoạt
type Request struct {
    Field1 string `binding:""`  // Validate in business logic
    Field2 int    `binding:""`  // Validate based on context
}

func ProcessRequest(req *Request) error {
    if req.Action == "CREATE" && req.Field1 == "" {
        return errors.New("field1 is required for CREATE")
    }
    if req.Action == "UPDATE" && req.Field2 == 0 {
        return errors.New("field2 is required for UPDATE")
    }
    // ...
}
```

---

**Status:** ✅ FIXED  
**Date:** 2026-02-05  
**Impact:** Critical - Blocking reject functionality
