# Sửa lỗi: Yêu cầu giải trình khi không có chênh lệch

## Vấn đề
Khi cashier bàn giao đầy đủ (actual_cash = expected_cash), hệ thống vẫn yêu cầu giải trình chênh lệch với lỗi:
```
variance requires documentation: reason and notes are required
```

## Nguyên nhân
1. **Floating-point precision issue**: Khi tính `variance = actualCash - expectedCash`, kết quả có thể là số rất nhỏ (ví dụ: 0.0000000001) thay vì chính xác 0
2. **Duplicate variance check**: Code check variance 2 lần:
   - Lần 1: Trong `FundHandover` (đúng)
   - Lần 2: Trong `CashierShift` (thừa, gây lỗi)

## Giải pháp

### 1. Thêm tolerance threshold cho floating-point comparison

#### Backend: `backend/domain/cashier/fund_handover.go`
```go
func (fh *FundHandover) HasVariance() bool {
	const tolerance = 0.01
	if fh.VarianceAmount < 0 {
		return -fh.VarianceAmount >= tolerance
	}
	return fh.VarianceAmount >= tolerance
}
```

#### Backend: `backend/domain/cashier/value_objects.go`
```go
func (v *Variance) RequiresDocumentation() bool {
	const tolerance = 0.01
	if v.Amount < 0 {
		return -v.Amount >= tolerance
	}
	return v.Amount >= tolerance
}
```

#### Frontend: `frontend/src/views/CashierShiftClosureV2.vue`
```javascript
const needsVarianceDocumentation = computed(() => {
  const tolerance = 0.01
  if (closureData.value.variance === null) return false
  const absVariance = Math.abs(closureData.value.variance)
  return absVariance >= tolerance
})
```

### 2. Xóa duplicate variance check

#### Backend: `backend/application/services/cashier_shift_service.go`

**TRƯỚC (SAI - check 2 lần):**
```go
// 6. Document variance if exists (FundHandover)
if handover.HasVariance() {
    if varianceReason == nil || varianceNotes == nil {
        return nil, errors.New("variance requires documentation...")
    }
    ...
}

// 12. Document variance in shift if needed (CashierShift) ❌ THỪA
if shift.Variance != nil && shift.Variance.RequiresDocumentation() {
    if varianceReason == nil || varianceNotes == nil {
        return nil, errors.New("variance requires documentation...")
    }
    ...
}
```

**SAU (ĐÚNG - check 1 lần, sync nếu có):**
```go
// 6. Document variance if exists (FundHandover) - CHECK DUY NHẤT
if handover.HasVariance() {
    if varianceReason == nil || varianceNotes == nil {
        return nil, errors.New("variance requires documentation...")
    }
    ...
}

// 12. Sync variance documentation from FundHandover to CashierShift
if handover.HasVariance() && handover.VarianceReason != nil {
    // Chỉ sync, không check lại
    shift.DocumentVariance(*handover.VarianceReason, handover.VarianceNotes, ...)
}
```

## Kết quả
- Chênh lệch < 0.01 (1 xu): Không yêu cầu giải trình ✅
- Chênh lệch >= 0.01: Yêu cầu giải trình (check 1 lần duy nhất)
- Variance documentation được sync từ FundHandover sang CashierShift

## Files changed
- `backend/domain/cashier/fund_handover.go` - HasVariance() với tolerance
- `backend/domain/cashier/value_objects.go` - RequiresDocumentation() với tolerance
- `backend/application/services/cashier_shift_service.go` - Xóa duplicate check, chỉ sync
- `frontend/src/views/CashierShiftClosureV2.vue` - needsVarianceDocumentation với tolerance

