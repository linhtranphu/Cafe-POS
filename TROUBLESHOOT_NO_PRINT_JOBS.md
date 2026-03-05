# Troubleshoot: Không có Print Jobs khi Collect Payment

## Checklist Debug

### 1. Kiểm tra code đã deploy chưa
```bash
./check-backend-version.sh
```

Nếu chưa deploy → Deploy code mới theo `DEPLOY_READY.md`

### 2. Kiểm tra Auto-Print có bật không
```bash
docker exec cafe-pos-mongodb mongosh cafe_pos_db --eval 'db.shop_settings.findOne({}, {auto_print_enabled: 1})'
```

Nếu `auto_print_enabled: false` → Vào Settings bật lại

### 3. Kiểm tra có Default Printer không
```bash
docker exec cafe-pos-mongodb mongosh cafe_pos_db --eval '
db.printer_configs.find({
  is_default: true,
  is_enabled: true
}, {name: 1, type: 1, ip_address: 1, port: 1})'
```

Phải có ít nhất:
- 1 BILL printer (default, enabled)
- 1 LABEL printer (default, enabled)

Nếu không có → Vào #/printer-management để cấu hình

### 4. Xem backend logs khi collect payment
```bash
# Xóa logs cũ
docker logs cafe-pos-backend 2>&1 > /dev/null

# Collect payment một order

# Xem logs mới
docker logs cafe-pos-backend 2>&1 | grep -i "print\|payment"
```

Expect thấy:
```
INFO: Print jobs created for order 20260301-XXXXXX
[PRINT] Creating print jobs for order 20260301-XXXXXX
[PRINT] Using HTML for print bridge - order_id=xxx, bridge_url=https://print.tacafe.store
[PRINT] Bill job created - job_id=xxx
```

### 5. Kiểm tra có lỗi không
```bash
docker logs cafe-pos-backend 2>&1 | grep -i "error\|failed" | tail -20
```

## Các lỗi thường gặp

### Lỗi 1: "failed to get default bill template"
**Nguyên nhân:** Không có template trong database (nhưng không cần nữa vì dùng HTML file)

**Giải pháp:** Code mới không cần template từ database nữa. Đảm bảo đã deploy code mới.

### Lỗi 2: "no default bill printer configured"
**Nguyên nhân:** Không có printer mặc định

**Giải pháp:**
```bash
# Kiểm tra printers
docker exec cafe-pos-mongodb mongosh cafe_pos_db --eval 'db.printer_configs.find({})'

# Set một printer làm default
docker exec cafe-pos-mongodb mongosh cafe_pos_db --eval '
db.printer_configs.updateOne(
  {type: "BILL"},
  {$set: {is_default: true, is_enabled: true}}
)'
```

### Lỗi 3: "Auto-print disabled"
**Nguyên nhân:** Auto-print bị tắt

**Giải pháp:**
```bash
docker exec cafe-pos-mongodb mongosh cafe_pos_db --eval '
db.shop_settings.updateOne(
  {},
  {$set: {auto_print_enabled: true}}
)'
```

### Lỗi 4: Không thấy log gì cả
**Nguyên nhân:** Code chưa được deploy hoặc order không chuyển sang PAID

**Giải pháp:**
1. Kiểm tra order status:
```bash
docker exec cafe-pos-mongodb mongosh cafe_pos_db --eval '
db.orders.find({}).sort({created_at: -1}).limit(1).forEach(function(o) {
  print("Order: " + o.order_number);
  print("Status: " + o.status);
  print("Paid At: " + o.paid_at);
})'
```

2. Nếu status không phải "PAID" → Collect payment chưa thành công

3. Nếu status là "PAID" nhưng không có log → Code chưa deploy

## Debug Flow

```
1. Check backend version
   ↓
2. Check auto-print enabled
   ↓
3. Check default printers exist
   ↓
4. Collect payment một order
   ↓
5. Check backend logs
   ↓
6. Check print jobs in database
```

## Scripts hỗ trợ

- `./check-backend-version.sh` - Kiểm tra code đã deploy chưa
- `./debug-collect-payment.sh` - Debug toàn bộ flow
- `./test-full-print-flow.sh` - Test tất cả components
- `./diagnose-print-jobs.sh` - Xem chi tiết print jobs

## Nếu vẫn không được

1. Restart backend:
```bash
docker restart cafe-pos-backend
docker logs -f cafe-pos-backend
```

2. Kiểm tra code trong container:
```bash
docker exec cafe-pos-backend ls -la /root/application/services/templates/
docker exec cafe-pos-backend cat /root/application/services/print_service.go | grep "Using HTML for print bridge"
```

3. Test trực tiếp API:
```bash
# Get một order ID
ORDER_ID=$(docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval 'db.orders.findOne({status: "SERVED"})._id.toString()')

# Collect payment
curl -X POST http://localhost:8080/api/waiter/orders/$ORDER_ID/payment \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 50000,
    "payment_method": "CASH",
    "collector_id": "xxx",
    "collector_name": "Test"
  }'
```

4. Xem response và logs ngay sau đó
