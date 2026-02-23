# Fix: Shift Revenue Calculation Logic

## Vấn đề

Trong view `http://localhost:5173/#/shifts`, hệ thống chỉ tính tiền cho orders có status `PAID`, `IN_PROGRESS`, và `SERVED`. Logic này không chính xác vì bỏ sót các status khác đã thanh toán.

## Order Status Flow

```
CREATED (chưa thanh toán)
  ↓ [Thanh toán]
PAID (đã thanh toán, chưa giao cho pha chế)
  ↓ [Gửi cho barista]
QUEUED (đã gửi cho barista, chờ nhận)
  ↓ [Barista nhận]
IN_PROGRESS (đang pha)
  ↓ [Pha xong]
READY (pha xong, chờ giao)
  ↓ [Giao cho khách]
SERVED (đã giao)
  ↓ [Chốt ca]
LOCKED (đã chốt ca)

Nhánh phụ:
CANCELLED (đã hủy - có thể từ bất kỳ status nào)
REFUNDED (đã hoàn tiền - từ PAID trở đi)
```

## Logic hiện tại (SAI)

**File**: `backend/application/services/shift_service.go`

```go
for _, o := range orders {
    // ❌ Chỉ tính 3 status
    if o.Status == order.StatusPaid || 
       o.Status == order.StatusInProgress || 
       o.Status == order.StatusServed {
        if o.PaymentMethod == order.PaymentCash {
            cashRevenue += o.Total
        } else if o.PaymentMethod == order.PaymentTransfer || 
                  o.PaymentMethod == order.PaymentQR {
            transferRevenue += o.Total
        }
    }
}
```

**Vấn đề**:
- ❌ Bỏ sót `QUEUED` - Order đã thanh toán, đang chờ barista nhận
- ❌ Bỏ sót `READY` - Order đã pha xong, chờ giao
- ❌ Bỏ sót `LOCKED` - Order đã chốt ca

## Logic đúng

Nên tính tiền cho TẤT CẢ orders đã thanh toán, tức là LOẠI TRỪ các status chưa thanh toán hoặc đã hủy/hoàn tiền:

```go
for _, o := range orders {
    // ✅ Loại trừ các status KHÔNG tính tiền
    if o.Status == order.StatusCreated ||    // Chưa thanh toán
       o.Status == order.StatusCancelled ||  // Đã hủy
       o.Status == order.StatusRefunded {    // Đã hoàn tiền
        continue // Skip order này
    }
    
    // Tính tiền cho tất cả status còn lại
    if o.PaymentMethod == order.PaymentCash {
        cashRevenue += o.Total
    } else if o.PaymentMethod == order.PaymentTransfer || 
              o.PaymentMethod == order.PaymentQR {
        transferRevenue += o.Total
    }
}
```

## Phân tích từng status

| Status | Đã thanh toán? | Nên tính? | Logic hiện tại | Logic đúng |
|--------|---------------|-----------|----------------|------------|
| CREATED | ❌ Chưa | ❌ Không | ❌ Không tính | ✅ Không tính |
| PAID | ✅ Rồi | ✅ Có | ✅ Tính | ✅ Tính |
| QUEUED | ✅ Rồi | ✅ Có | ❌ KHÔNG tính | ✅ Tính |
| IN_PROGRESS | ✅ Rồi | ✅ Có | ✅ Tính | ✅ Tính |
| READY | ✅ Rồi | ✅ Có | ❌ KHÔNG tính | ✅ Tính |
| SERVED | ✅ Rồi | ✅ Có | ✅ Tính | ✅ Tính |
| LOCKED | ✅ Rồi | ✅ Có | ❌ KHÔNG tính | ✅ Tính |
| CANCELLED | ❌ Đã hủy | ❌ Không | ❌ Không tính | ✅ Không tính |
| REFUNDED | ❌ Đã hoàn | ❌ Không | ❌ Không tính | ✅ Không tính |

## Impact Analysis

### Trường hợp bị ảnh hưởng:

1. **Order ở status QUEUED**:
   - Waiter đã thu tiền và gửi order cho barista
   - Barista chưa nhận (vẫn ở QUEUED)
   - Waiter xem shift → Tiền KHÔNG được tính ❌

2. **Order ở status READY**:
   - Barista đã pha xong
   - Chưa giao cho khách (vẫn ở READY)
   - Waiter xem shift → Tiền KHÔNG được tính ❌

3. **Order ở status LOCKED**:
   - Order đã chốt ca
   - Waiter xem shift cũ → Tiền KHÔNG được tính ❌

### Ví dụ cụ thể:

```
Waiter có 5 orders trong ca:
1. Order A: PAID - 50k → ✅ Tính (hiện tại)
2. Order B: QUEUED - 30k → ❌ KHÔNG tính (SAI!)
3. Order C: IN_PROGRESS - 40k → ✅ Tính (hiện tại)
4. Order D: READY - 35k → ❌ KHÔNG tính (SAI!)
5. Order E: SERVED - 45k → ✅ Tính (hiện tại)

Tổng thực tế: 200k
Tổng hiển thị: 135k (thiếu 65k!)
```

## Giải pháp

### Option 1: Whitelist (Liệt kê status cần tính) - KHÔNG khuyến nghị

```go
if o.Status == order.StatusPaid || 
   o.Status == order.StatusQueued ||
   o.Status == order.StatusInProgress || 
   o.Status == order.StatusReady ||
   o.Status == order.StatusServed ||
   o.Status == order.StatusLocked {
    // Tính tiền
}
```

**Nhược điểm**: Nếu thêm status mới trong tương lai, dễ quên update

### Option 2: Blacklist (Loại trừ status KHÔNG tính) - KHUYẾN NGHỊ ✅

```go
// Skip orders that should NOT be counted
if o.Status == order.StatusCreated ||    // Not paid yet
   o.Status == order.StatusCancelled ||  // Cancelled
   o.Status == order.StatusRefunded {    // Refunded
    continue
}

// Count all other statuses (paid orders)
if o.PaymentMethod == order.PaymentCash {
    cashRevenue += o.Total
} else if o.PaymentMethod == order.PaymentTransfer || 
          o.PaymentMethod == order.PaymentQR {
    transferRevenue += o.Total
}
```

**Ưu điểm**: 
- Rõ ràng hơn - chỉ loại trừ những gì KHÔNG tính
- An toàn hơn - status mới sẽ tự động được tính (default behavior đúng)
- Dễ maintain hơn

## Implementation

### File cần sửa:

`backend/application/services/shift_service.go` - Method `CalculateTransferRevenue`

### Code change:

```go
// Calculate cash and transfer revenue separately
cashRevenue := 0.0
transferRevenue := 0.0
for _, o := range orders {
    // Skip orders that should NOT be counted in revenue
    if o.Status == order.StatusCreated ||    // Not paid yet
       o.Status == order.StatusCancelled ||  // Cancelled before/after payment
       o.Status == order.StatusRefunded {    // Payment refunded
        continue
    }
    
    // Count all other statuses (PAID, QUEUED, IN_PROGRESS, READY, SERVED, LOCKED)
    if o.PaymentMethod == order.PaymentCash {
        cashRevenue += o.Total
    } else if o.PaymentMethod == order.PaymentTransfer || 
              o.PaymentMethod == order.PaymentQR {
        transferRevenue += o.Total
    }
}
```

## Testing

### Test Case 1: Order flow bình thường
```
1. Create order (CREATED) → Revenue: 0
2. Pay order (PAID) → Revenue: +50k ✅
3. Send to barista (QUEUED) → Revenue: 50k ✅ (không mất)
4. Barista accept (IN_PROGRESS) → Revenue: 50k ✅
5. Barista done (READY) → Revenue: 50k ✅
6. Serve customer (SERVED) → Revenue: 50k ✅
7. Close shift (LOCKED) → Revenue: 50k ✅
```

### Test Case 2: Order bị hủy
```
1. Create order (CREATED) → Revenue: 0
2. Cancel (CANCELLED) → Revenue: 0 ✅
```

### Test Case 3: Order hoàn tiền
```
1. Create order (CREATED) → Revenue: 0
2. Pay order (PAID) → Revenue: +50k
3. Refund (REFUNDED) → Revenue: 0 ✅ (trừ lại)
```

### Test Case 4: Mixed orders
```
Order A: CREATED - 10k → Không tính
Order B: PAID - 20k → Tính
Order C: QUEUED - 30k → Tính (FIX!)
Order D: IN_PROGRESS - 40k → Tính
Order E: READY - 50k → Tính (FIX!)
Order F: SERVED - 60k → Tính
Order G: CANCELLED - 70k → Không tính
Order H: REFUNDED - 80k → Không tính

Expected: 20 + 30 + 40 + 50 + 60 = 200k
```

## Deployment

1. ✅ Phân tích vấn đề
2. ⏳ Update code
3. ⏳ Backend restart
4. ⏳ Test với data thực tế
5. ⏳ Verify shift revenue calculations

## Notes

- Logic này ảnh hưởng đến:
  - Shift view (waiter/barista)
  - Shift closure
  - Revenue reports
  - Cash handover calculations
  
- Cần kiểm tra xem có nơi nào khác dùng logic tương tự không

- Có thể cần recalculate revenue cho các shifts cũ nếu data bị sai
