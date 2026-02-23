# Nil Pointer Dereference Fix

## ❌ Lỗi

```
runtime error: invalid memory address or nil pointer dereference
/backend/application/services/cashier_shift_service.go:745
if err := shift.DocumentVariance(*varianceReason, *varianceNotes, userID, deviceID, now)
```

## 🔍 Nguyên Nhân

Khi đóng ca **không có chênh lệch** (actual_cash = expected_cash):
- Frontend gửi: `variance_reason: null`, `variance_notes: null`
- Backend nhận: `varianceReason = nil`, `varianceNotes = nil`
- Code cố dereference `*varianceNotes` → **PANIC!**

## ✅ Giải Pháp

Thêm kiểm tra nil trước khi dereference:

```go
// 12. Document variance in shift if needed
if shift.Variance != nil && shift.Variance.RequiresDocumentation() {
    // ✅ THÊM KIỂM TRA NÀY
    if varianceReason == nil || varianceNotes == nil {
        return nil, errors.New("variance requires documentation: reason and notes are required")
    }
    if err := shift.DocumentVariance(*varianceReason, *varianceNotes, userID, deviceID, now); err != nil {
        return nil, fmt.Errorf("failed to document variance in shift: %w", err)
    }
}
```

## 📊 Luồng Xử Lý

### Trường hợp 1: Không có chênh lệch
```
actual_cash = 2,000,000
expected_cash = 2,000,000
variance = 0

→ shift.Variance = nil hoặc !RequiresDocumentation()
→ Bỏ qua DocumentVariance
→ ✅ OK
```

### Trường hợp 2: Có chênh lệch nhưng thiếu documentation
```
actual_cash = 1,995,000
expected_cash = 2,000,000
variance = -5,000

varianceReason = nil
varianceNotes = nil

→ shift.Variance.RequiresDocumentation() = true
→ Kiểm tra nil → ❌ Error: "variance requires documentation"
```

### Trường hợp 3: Có chênh lệch và có documentation
```
actual_cash = 1,995,000
expected_cash = 2,000,000
variance = -5,000

varianceReason = "COUNTING_ERROR"
varianceNotes = "Đếm nhầm tờ 50k"

→ shift.Variance.RequiresDocumentation() = true
→ Kiểm tra nil → ✅ Pass
→ DocumentVariance() → ✅ OK
```

## 🚀 Khởi Động Lại Backend

```bash
cd backend
# Ctrl+C để dừng backend hiện tại
go run main.go
```

## 🧪 Test

### Test 1: Không có chênh lệch
```json
{
  "actual_cash": 2000000,
  "variance_reason": null,
  "variance_notes": null
}
```
**Kết quả**: ✅ Thành công

### Test 2: Có chênh lệch với documentation
```json
{
  "actual_cash": 1995000,
  "variance_reason": "COUNTING_ERROR",
  "variance_notes": "Đếm nhầm tờ 50k thành 100k"
}
```
**Kết quả**: ✅ Thành công

### Test 3: Có chênh lệch nhưng thiếu documentation
```json
{
  "actual_cash": 1995000,
  "variance_reason": null,
  "variance_notes": null
}
```
**Kết quả**: ❌ Error: "variance requires documentation"

## 📝 Files Đã Sửa

- ✅ `backend/application/services/cashier_shift_service.go` (line 745)

## ✅ Checklist

- [x] Thêm nil check trước khi dereference
- [x] Error message rõ ràng
- [ ] Restart backend
- [ ] Test trên trình duyệt

---

**Làm ngay**: Restart backend và test lại!

```bash
cd backend
go run main.go
```
