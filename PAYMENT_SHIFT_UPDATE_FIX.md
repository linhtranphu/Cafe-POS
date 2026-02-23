# Fix: Current Cash và Transfer Revenue không được cập nhật khi có Order

## Vấn đề

Khi tạo order và thanh toán (dù là tiền mặt hay chuyển khoản), các giá trị `current_cash` và `transfer_revenue` trong shift không được cập nhật đúng cách.

### Nguyên nhân

Trong file `backend/application/services/order_service.go`, hàm `CollectPayment()` chỉ cập nhật:
- `current_cash` và `remaining_cash` khi thanh toán bằng **TIỀN MẶT** (`PaymentCash`)
- **KHÔNG** cập nhật `transfer_revenue` và `remaining_transfer` khi thanh toán bằng **CHUYỂN KHOẢN** (`PaymentTransfer`) hoặc **QR** (`PaymentQR`)

### Code cũ (có lỗi)

```go
// Update shift cash if payment is cash and order has shift_id
if req.PaymentMethod == order.PaymentCash && !o.ShiftID.IsZero() {
    shift, err := s.shiftRepo.FindByID(ctx, o.ShiftID)
    if err == nil && shift != nil {
        // Add cash to shift
        shift.RemainingCash += req.Amount
        shift.CurrentCash += req.Amount
        shift.TotalRevenue += req.Amount
        
        // Update shift
        s.shiftRepo.Update(ctx, o.ShiftID, shift)
    }
}
```

**Vấn đề:** Chỉ xử lý khi `PaymentMethod == PaymentCash`, bỏ qua các phương thức thanh toán khác.

## Giải pháp

Cập nhật logic để xử lý TẤT CẢ các phương thức thanh toán với debug logging chi tiết:

### Code mới (đã sửa)

```go
// Update shift revenue based on payment method
if !o.ShiftID.IsZero() {
    fmt.Printf("DEBUG: Payment received - ShiftID: %s, Method: %s, Amount: %.2f\n", 
        o.ShiftID.Hex(), req.PaymentMethod, req.Amount)
    
    shift, err := s.shiftRepo.FindByID(ctx, o.ShiftID)
    if err != nil {
        fmt.Printf("ERROR: Failed to find shift: %v\n", err)
    } else if shift == nil {
        fmt.Printf("ERROR: Shift is nil\n")
    } else {
        fmt.Printf("DEBUG: Found shift - ID: %s\n", shift.ID.Hex())
        fmt.Printf("DEBUG: BEFORE - CurrentCash: %.2f, RemainingCash: %.2f, TransferRevenue: %.2f, RemainingTransfer: %.2f, TotalRevenue: %.2f\n",
            shift.CurrentCash, shift.RemainingCash, shift.TransferRevenue, shift.RemainingTransfer, shift.TotalRevenue)
        
        // Update total revenue for all payment methods
        shift.TotalRevenue += req.Amount
        
        // Update specific payment method fields
        if req.PaymentMethod == order.PaymentCash {
            fmt.Printf("DEBUG: Processing CASH payment\n")
            shift.RemainingCash += req.Amount
            shift.CurrentCash += req.Amount
        } else if req.PaymentMethod == order.PaymentTransfer || req.PaymentMethod == order.PaymentQR {
            fmt.Printf("DEBUG: Processing TRANSFER/QR payment\n")
            shift.TransferRevenue += req.Amount
            shift.RemainingTransfer += req.Amount
        } else {
            fmt.Printf("WARNING: Unknown payment method: %s\n", req.PaymentMethod)
        }
        
        fmt.Printf("DEBUG: AFTER - CurrentCash: %.2f, RemainingCash: %.2f, TransferRevenue: %.2f, RemainingTransfer: %.2f, TotalRevenue: %.2f\n",
            shift.CurrentCash, shift.RemainingCash, shift.TransferRevenue, shift.RemainingTransfer, shift.TotalRevenue)
        
        // Update shift
        if err := s.shiftRepo.Update(ctx, o.ShiftID, shift); err != nil {
            fmt.Printf("ERROR: Failed to update shift in DB: %v\n", err)
        } else {
            fmt.Printf("DEBUG: Shift updated successfully in DB\n")
            
            // Verify the update by reading back
            verifyShift, verifyErr := s.shiftRepo.FindByID(ctx, o.ShiftID)
            if verifyErr == nil && verifyShift != nil {
                fmt.Printf("DEBUG: VERIFY - CurrentCash: %.2f, RemainingCash: %.2f, TransferRevenue: %.2f, RemainingTransfer: %.2f, TotalRevenue: %.2f\n",
                    verifyShift.CurrentCash, verifyShift.RemainingCash, verifyShift.TransferRevenue, verifyShift.RemainingTransfer, verifyShift.TotalRevenue)
            } else {
                fmt.Printf("ERROR: Failed to verify shift update: %v\n", verifyErr)
            }
        }
    }
}
```

### Thay đổi chính

1. **Loại bỏ điều kiện** `req.PaymentMethod == order.PaymentCash` ở ngoài cùng
2. **Cập nhật `TotalRevenue`** cho TẤT CẢ phương thức thanh toán
3. **Thêm xử lý cho Transfer/QR:**
   - Cập nhật `TransferRevenue` 
   - Cập nhật `RemainingTransfer`
4. **Giữ nguyên xử lý cho Cash:**
   - Cập nhật `CurrentCash`
   - Cập nhật `RemainingCash`
5. **Thêm debug logging chi tiết:**
   - Log payment method và amount
   - Log giá trị BEFORE và AFTER update
   - Verify bằng cách đọc lại từ DB
   - Log errors rõ ràng

## Các trường trong Shift

```go
type Shift struct {
    // Cash fields
    StartCash        float64 // Tiền mặt ban đầu
    CurrentCash      float64 // Tiền mặt hiện tại (tăng khi có order cash)
    RemainingCash    float64 // Tiền mặt còn lại (sau khi nộp)
    HandedOverCash   float64 // Tiền mặt đã nộp
    
    // Transfer fields
    TransferRevenue    float64 // Doanh thu chuyển khoản (tăng khi có order transfer)
    RemainingTransfer  float64 // Chuyển khoản còn lại (sau khi nộp)
    HandedOverTransfer float64 // Chuyển khoản đã nộp
    
    // Total
    TotalRevenue float64 // Tổng doanh thu (cash + transfer + qr)
}
```

## Kiểm tra

### Test 1: Direct DB Test (Recommended)

Test trực tiếp với MongoDB để loại trừ các vấn đề về API:

```bash
./test-payment-shift-update.sh
```

Test này sẽ:
- Tạo một shift mới
- Simulate cash payment 50,000 VND
- Simulate transfer payment 75,000 VND
- Verify các giá trị trong DB

### Test 2: API Test với Debug

Test qua API với menu items thực:

```bash
./test-cash-payment-debug.sh
```

Test này sẽ:
- Login và start shift
- Tạo order với menu item thực
- Thanh toán bằng CASH
- Kiểm tra shift có được cập nhật không

### Test 3: Full API Test

Test đầy đủ cả cash và transfer:

```bash
./test-payment-update.sh
```

### Xem Backend Logs

Khi chạy test, xem backend logs để thấy debug messages:

```
DEBUG: Payment received - ShiftID: xxx, Method: CASH, Amount: 50000.00
DEBUG: Found shift - ID: xxx
DEBUG: BEFORE - CurrentCash: 100000.00, RemainingCash: 100000.00, ...
DEBUG: Processing CASH payment
DEBUG: AFTER - CurrentCash: 150000.00, RemainingCash: 150000.00, ...
DEBUG: Shift updated successfully in DB
DEBUG: VERIFY - CurrentCash: 150000.00, RemainingCash: 150000.00, ...
```

### Kết quả mong đợi

Sau khi thanh toán:
- 1 order với CASH 50,000 VND
- 1 order với TRANSFER 75,000 VND

Shift sẽ có:
```
current_cash: 150,000 (100,000 start + 50,000 cash)
remaining_cash: 150,000
transfer_revenue: 75,000
remaining_transfer: 75,000
total_revenue: 125,000 (50,000 + 75,000)
```

## Troubleshooting

### Nếu vẫn không update:

1. **Kiểm tra backend logs** - Xem có DEBUG messages không
2. **Kiểm tra PaymentMethod value** - Có thể là "cash" thay vì "CASH"
3. **Kiểm tra ShiftID** - Có thể order không có shift_id
4. **Kiểm tra MongoDB** - Có thể có vấn đề với connection hoặc permissions
5. **Kiểm tra race condition** - Có thể có code khác cũng đang update shift

### Debug Commands

```bash
# Xem shift trong MongoDB
mongo cafe_pos --eval 'db.shifts.find().pretty()'

# Xem orders trong MongoDB
mongo cafe_pos --eval 'db.orders.find().pretty()'

# Xem backend logs
tail -f backend.log
```

## Files đã thay đổi

- `backend/application/services/order_service.go` - Hàm `CollectPayment()` với debug logging
- `backend/cmd/test-payment-shift-update/main.go` - Direct DB test
- `test-payment-shift-update.sh` - Script chạy direct DB test
- `test-cash-payment-debug.sh` - Script test API với debug
- `test-payment-update.sh` - Script test full API

## Tác động

- ✅ `current_cash` được cập nhật đúng khi thanh toán tiền mặt
- ✅ `transfer_revenue` được cập nhật đúng khi thanh toán chuyển khoản/QR
- ✅ `total_revenue` được cập nhật cho TẤT CẢ phương thức thanh toán
- ✅ Không ảnh hưởng đến logic hiện tại của cash handover
- ✅ Tương thích với tính năng bank transfer handover đang phát triển
- ✅ Debug logging giúp troubleshoot dễ dàng

## Lưu ý

Logic này chỉ cập nhật khi:
1. Order có `shift_id` hợp lệ
2. Shift tồn tại trong database
3. Payment được thực hiện thành công

Nếu có lỗi khi cập nhật shift, payment vẫn được thực hiện (chỉ log error), đảm bảo không làm gián đoạn quy trình bán hàng.
