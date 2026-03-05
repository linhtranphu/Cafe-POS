# Logic In Bill Khi Thu Tiền

## Flow Hiện Tại

### 1. User Thu Tiền (Frontend)
```
User clicks "Thu tiền" 
→ OrderView.vue: processPayment()
→ orderStore.collectPayment(orderId, paymentData)
→ API: POST /waiter/orders/:id/payment
```

### 2. Backend Xử Lý Payment
```go
// backend/interfaces/http/order_handler.go
func CollectPayment(c *gin.Context) {
    // Validate payment
    // Call service layer
    o, err := h.orderService.CollectPayment(ctx, id, &req)
    // Return updated order
}
```

### 3. Service Layer - CollectPayment
```go
// backend/application/services/order_service.go
func (s *OrderService) CollectPayment(...) {
    // 1. Validate state transition
    // 2. Add to amount_paid
    // 3. Calculate total
    // 4. If fully paid → set status = PAID
    
    // 5. AUTO-PRINT LOGIC (CHỈ KHI STATUS = PAID)
    if o.Status == order.StatusPaid && s.printService != nil {
        // Check auto-print setting
        autoPrintEnabled := true
        if s.settingsRepo != nil {
            settings, err := s.settingsRepo.FindFirst(ctx)
            if err == nil && settings != nil {
                autoPrintEnabled = settings.AutoPrintEnabled
            }
        }
        
        if autoPrintEnabled {
            // Create print jobs asynchronously
            go func() {
                s.printService.CreatePrintJobsForOrder(ctx, o)
            }()
        }
    }
}
```

## Điều Kiện Để In Bill Tự Động

### ✅ Phải Thỏa TẤT CẢ Điều Kiện Sau:

1. **Order phải FULLY PAID**
   - `amount_paid >= total_amount`
   - `amount_due <= 0`
   - Status chuyển sang `PAID`

2. **Auto-print phải ENABLED**
   - Settings: `auto_print_enabled = true`
   - Check tại: Settings > "Tự động in khi thu tiền"

3. **Print Service phải available**
   - `printService != nil`
   - Backend đã khởi tạo print service

4. **Print Bridge phải hoạt động**
   - `PRINT_BRIDGE_URL` đã config trong .env
   - Print Bridge đang chạy và accessible

## Các Trường Hợp KHÔNG In

### ❌ Trường Hợp 1: Thanh Toán Một Phần
```
Order total: 100,000 VND
Payment: 50,000 VND
→ amount_paid = 50,000
→ amount_due = 50,000
→ Status vẫn là CREATED (không chuyển sang PAID)
→ KHÔNG IN
```

**Giải pháp:** Phải thu đủ tiền để order chuyển sang PAID

### ❌ Trường Hợp 2: Auto-Print Disabled
```
Settings: auto_print_enabled = false
→ Dù order PAID nhưng auto-print bị tắt
→ KHÔNG IN
```

**Giải pháp:** 
- Vào Settings
- Bật "Tự động in khi thu tiền"
- Save

### ❌ Trường Hợp 3: Print Bridge Không Kết Nối
```
PRINT_BRIDGE_URL không đúng hoặc Print Bridge không chạy
→ Print jobs được tạo nhưng không thể in
→ Print jobs status = failed
```

**Giải pháp:**
- Check Print Bridge: `curl http://localhost:3001/health`
- Check .env: `PRINT_BRIDGE_URL=http://localhost:3001`
- Start Print Bridge nếu chưa chạy

## Cách Kiểm Tra

### 1. Check Auto-Print Setting

```bash
# Via API
curl http://localhost:3000/api/settings

# Should see:
# "auto_print_enabled": true
```

### 2. Check Order Status After Payment

```bash
# Get order
curl http://localhost:3000/api/waiter/orders/:id

# Check:
# - status: "paid"
# - amount_due: 0
# - paid_at: "2024-..."
```

### 3. Check Print Jobs Created

```bash
# List print jobs
curl http://localhost:3000/api/print-jobs?limit=10

# Should see new jobs with:
# - order_number: matching your order
# - status: "pending" or "completed"
# - created_at: recent timestamp
```

### 4. Check Backend Logs

```bash
# Watch logs
docker logs -f backend | grep -i print

# Should see:
# "INFO: Print jobs created for order ORD-XXX"
# "Processing print job..."
# "Print job completed successfully"
```

### 5. Run Diagnostic Script

```bash
./diagnose-auto-print-on-payment.sh
```

## Debug Steps

### Step 1: Verify Auto-Print Enabled

```bash
# Check setting
curl http://localhost:3000/api/settings | grep auto_print

# If disabled, enable it:
curl -X PUT http://localhost:3000/api/settings \
  -H "Content-Type: application/json" \
  -d '{"auto_print_enabled": true}'
```

### Step 2: Test Full Payment Flow

```bash
# 1. Create order
# 2. Collect FULL payment (not partial)
# 3. Watch logs
docker logs -f backend

# Should see:
# "💰 [PAYMENT] Received - ..."
# "INFO: Print jobs created for order ORD-XXX"
```

### Step 3: Check Print Jobs

```bash
# List recent print jobs
curl http://localhost:3000/api/print-jobs?limit=5

# Check status:
# - pending: waiting to be processed
# - processing: being processed
# - completed: successfully printed
# - failed: error occurred
```

### Step 4: Check Print Worker

```bash
# Print worker logs
docker logs backend 2>&1 | grep -i "print worker\|processing print"

# Should see:
# "Print worker started"
# "Processing print job..."
```

### Step 5: Test Print Bridge

```bash
# Health check
curl http://localhost:3001/health

# Test print
curl -X POST http://localhost:3001/test-connection \
  -H "Content-Type: application/json" \
  -d '{"printerIP": "192.168.1.115", "printerPort": 9100}'
```

## Common Issues

### Issue 1: "Auto-print disabled"

**Symptom:** Logs show "Auto-print disabled, skipping print jobs"

**Solution:**
```bash
# Enable in Settings UI
# Or via API:
curl -X PUT http://localhost:3000/api/settings \
  -H "Content-Type: application/json" \
  -d '{"auto_print_enabled": true}'
```

### Issue 2: "Order not fully paid"

**Symptom:** Order status không chuyển sang PAID

**Check:**
```bash
# Get order
curl http://localhost:3000/api/waiter/orders/:id

# Verify:
# amount_paid >= total_amount
# amount_due <= 0
```

**Solution:** Thu đủ tiền (không thanh toán một phần)

### Issue 3: "Print jobs created but not printing"

**Symptom:** Print jobs có status "pending" hoặc "failed"

**Check:**
```bash
# 1. Print Bridge connection
curl http://localhost:3001/health

# 2. Print worker logs
docker logs backend | grep -i "print worker"

# 3. Print job errors
curl http://localhost:3000/api/print-jobs | grep failed
```

**Solution:**
- Restart Print Bridge
- Check PRINT_BRIDGE_URL in .env
- Check printer IP/port settings

### Issue 4: "No print jobs created"

**Symptom:** Không thấy print jobs sau khi thu tiền

**Possible causes:**
1. Auto-print disabled
2. Order không fully paid
3. Print service không khởi tạo

**Check:**
```bash
# 1. Auto-print setting
curl http://localhost:3000/api/settings | grep auto_print

# 2. Order status
curl http://localhost:3000/api/waiter/orders/:id | grep status

# 3. Backend logs
docker logs backend | grep -i "print service"
```

## Test Checklist

- [ ] Auto-print enabled in Settings
- [ ] PRINT_BRIDGE_URL configured in .env
- [ ] Print Bridge running and healthy
- [ ] Create test order
- [ ] Collect FULL payment (not partial)
- [ ] Order status changes to PAID
- [ ] Print jobs created (check API)
- [ ] Print worker processes jobs (check logs)
- [ ] Bill prints successfully

## Code References

- Frontend: `frontend/src/views/OrderView.vue` - processPayment()
- API Handler: `backend/interfaces/http/order_handler.go` - CollectPayment()
- Service: `backend/application/services/order_service.go` - CollectPayment()
- Print Service: `backend/application/services/print_service.go` - CreatePrintJobsForOrder()
- Settings: `backend/domain/settings/settings.go` - AutoPrintEnabled field

---

**TL;DR:** 
1. Bật "Tự động in khi thu tiền" trong Settings
2. Thu ĐỦ tiền (không thanh toán một phần)
3. Check Print Bridge đang chạy
4. Order phải chuyển sang status PAID thì mới in
