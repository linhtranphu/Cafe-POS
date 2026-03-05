# Chẩn Đoán: Tại Sao Không Tự Động In Bill?

## Kết Quả Kiểm Tra

### ✅ Những gì ĐÃ OK:
1. **auto_print_enabled = true** - Cài đặt đã bật
2. **BILL template** - Có template mặc định
3. **LABEL template** - Có template mặc định  
4. **Backend đang chạy** - Service hoạt động
5. **Print jobs được tạo** - Hệ thống đã tạo print jobs trong database

### ❌ VẤN ĐỀ CHÍNH:
**KHÔNG CÓ PRINTER NÀO ĐƯỢC CẤU HÌNH!**

- ❌ BILL printer: MISSING
- ❌ LABEL printer: MISSING

## Giải Thích

Khi bạn:
1. Tạo order
2. Collect payment (tiền mặt)
3. Order status → PAID

Backend **ĐÃ** tạo print jobs (thấy trong database):
```
Order: 20260227-010641-363 | Type: BILL | Status: COMPLETED
Order: 20260227-010641-363 | Type: LABEL | Status: FAILED
```

Nhưng:
- BILL job → COMPLETED (có thể in được vì có printer config cũ?)
- LABEL job → FAILED (không có printer)

## Nguyên Nhân

Trong code `CollectPayment` (backend/application/services/order_service.go):

```go
if o.Status == order.StatusPaid && s.printService != nil {
    // Check auto-print setting
    autoPrintEnabled := true
    if s.settingsRepo != nil {
        settings, err := s.settingsRepo.FindFirst(ctx)
        if err == nil && settings != nil {
            autoPrintEnabled = settings.AutoPrintEnabled  // ✅ TRUE
        }
    }

    if autoPrintEnabled {  // ✅ PASS
        go func() {
            // Tạo print jobs
            if err := s.printService.CreatePrintJobsForOrder(printCtx, o); err != nil {
                fmt.Printf("ERROR: Failed to create print jobs: %v\n", err)
            } else {
                fmt.Printf("INFO: Print jobs created for order %s\n", o.OrderNumber)
            }
        }()
    }
}
```

Print jobs **ĐÃ ĐƯỢC TẠO** nhưng không có printer để gửi đến!

## Giải Pháp

### Cách 1: Thêm Printer Qua Script (NHANH)

```bash
chmod +x add-printer-quick.sh
./add-printer-quick.sh
```

Script sẽ hỏi:
- Tên máy in (ví dụ: Máy in Bill)
- IP Address (ví dụ: 192.168.1.115)
- Port (mặc định: 9100)

Sau đó tự động tạo 2 printers:
- BILL printer
- LABEL printer (cùng máy in)

### Cách 2: Thêm Printer Qua UI

1. Mở frontend
2. Vào **Print Management** → **Máy In**
3. Click **Thêm Máy In**
4. Nhập thông tin:
   - Type: BILL
   - Name: Máy in Bill
   - Connection: Network
   - IP: 192.168.1.115
   - Port: 9100
   - Enabled: ✅
   - Default: ✅
5. Lưu
6. Lặp lại cho LABEL printer

### Cách 3: Thêm Trực Tiếp Vào MongoDB

```javascript
use cafe_pos

// BILL printer
db.printers.insertOne({
  name: "Máy in Bill",
  type: "BILL",
  connection_type: "NETWORK",
  ip_address: "192.168.1.115",
  port: 9100,
  paper_width: 72,
  enabled: true,
  is_default: true,
  created_at: new Date(),
  updated_at: new Date()
})

// LABEL printer
db.printers.insertOne({
  name: "Máy in Label",
  type: "LABEL",
  connection_type: "NETWORK",
  ip_address: "192.168.1.115",
  port: 9100,
  paper_width: 58,
  enabled: true,
  is_default: true,
  created_at: new Date(),
  updated_at: new Date()
})
```

## Kiểm Tra Sau Khi Sửa

```bash
./test-auto-print-flow.sh
```

Kết quả mong đợi:
```
✓ auto_print_enabled = true
✓ BILL printer configured
✓ LABEL printer configured
✓ BILL template configured
✓ LABEL template configured
✓ Backend đang chạy

✓ TẤT CẢ CẤU HÌNH OK!
```

## Test Auto Print

1. Tạo order mới
2. Collect payment (tiền mặt hoặc chuyển khoản)
3. Order status → PAID
4. Bill và label sẽ **TỰ ĐỘNG IN**

## API Flow

```
Frontend: POST /api/waiter/orders/:id/payment
    ↓
Backend: OrderHandler.CollectPayment()
    ↓
Backend: OrderService.CollectPayment()
    ↓
Check: auto_print_enabled = true? ✅
    ↓
Check: order.Status = PAID? ✅
    ↓
Create Print Jobs (async)
    ↓
PrintService.CreatePrintJobsForOrder()
    ↓
Create 1 BILL job + N LABEL jobs
    ↓
Print Worker picks up jobs
    ↓
Send to printer via Print Bridge
    ↓
✅ PRINTED!
```

## Tóm Tắt

**Vấn đề:** Không có printer nào được cấu hình trong database

**Giải pháp:** Thêm BILL printer và LABEL printer

**Sau khi sửa:** Auto-print sẽ hoạt động ngay lập tức!
