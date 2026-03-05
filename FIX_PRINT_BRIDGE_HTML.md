# Fix: Print Jobs sử dụng HTML và Print Bridge

## Vấn đề
1. Print jobs không được tạo vì thiếu text template
2. Backend không sử dụng print bridge URL đã cấu hình

## Giải pháp

### 1. Sửa `createBillJob` logic
**File:** `backend/application/services/print_service.go`

**Thay đổi:**
- Kiểm tra `shopSettings.PrintBridgeURL` trước
- Nếu có print bridge URL → render HTML và để print bridge xử lý
- Content type = "html" thay vì "binary" hoặc "text"
- Không cần chromedp/htmlRenderer trên backend nữa

**Logic mới:**
```
if shopSettings.PrintBridgeURL != "" {
    // Render HTML
    htmlContent = renderBillHTML(order, shopSettings)
    contentType = "html"
} else {
    // Fallback to local renderers (chromedp, htmlRenderer, text)
}
```

### 2. Thêm `renderBillHTML` function
**File:** `backend/application/services/print_service.go`

- Đọc HTML template từ `./application/services/templates/bill_template_optimized.html`
- Parse template với Go html/template
- Render với order data

### 3. Sửa `BridgePrinter.Print()` để xử lý HTML
**File:** `backend/infrastructure/printing/bridge_printer.go`

**Thay đổi:**
- Kiểm tra content type (HTML, base64, hoặc text)
- Nếu HTML → gọi `RenderAndPrint()` (print bridge sẽ convert HTML → image → ESC/POS)
- Nếu binary/text → gọi `PrintESCPOS()` như cũ

**Logic mới:**
```go
if isHTMLContent(content) {
    // Send HTML to print bridge for rendering
    bridgeClient.RenderAndPrint(ctx, content, printerIP, printerPort, 576)
} else {
    // Send ESC/POS data directly
    bridgeClient.PrintESCPOS(ctx, escposData, printerIP, printerPort)
}
```

### 4. Thêm helper functions
- `isHTMLContent()` - Kiểm tra content có phải HTML không
- `prepareBillData()` - Chuẩn bị data cho template
- `formatMoneyVN()` - Format tiền VND

## Flow hoàn chỉnh

### Khi collect payment:
```
1. CollectPayment()
2. CreatePrintJobsForOrder()
3. createBillJob()
   → Lấy shopSettings (có PrintBridgeURL)
   → Render HTML từ template
   → Tạo print job với content = HTML, contentType = "html"
   → Lưu vào database (status = PENDING)

4. Print Worker (background, mỗi 10s)
   → Lấy pending jobs
   → ProcessJob()
   → GetPrinter() → NewBridgePrinter()
   → printer.Print(htmlContent)
   → Phát hiện là HTML
   → bridgeClient.RenderAndPrint(html, printerIP, printerPort)
   → POST https://print.tacafe.store/render-and-print
   → Print Bridge render HTML → image → ESC/POS
   → Print Bridge gửi tới máy in
   → Update job status = COMPLETED
```

## Files thay đổi
1. `backend/application/services/print_service.go`
   - Sửa `createBillJob()` logic
   - Thêm `renderBillHTML()`
   - Thêm `prepareBillData()`
   - Thêm `formatMoneyVN()`
   - Thêm imports: `bytes`, `html/template`, `os`

2. `backend/infrastructure/printing/bridge_printer.go`
   - Sửa `Print()` để xử lý HTML
   - Thêm `isHTMLContent()`

## Deploy

### Build và deploy:
```bash
# Local
cd backend
docker build -t cafe-pos-backend:latest .

# Upload lên EC2 hoặc push git
git add .
git commit -m "fix: use HTML template and print bridge for auto-print"
git push

# EC2
cd ~/Cafe\ POS
git pull
cd backend
docker build -t cafe-pos-backend:latest .
docker stop cafe-pos-backend && docker rm cafe-pos-backend
docker run -d --name cafe-pos-backend --network cafe-pos-network -p 8080:8080 --env-file ../.env.ec2 --restart unless-stopped cafe-pos-backend:latest
```

## Verification

### 1. Kiểm tra backend logs
```bash
docker logs -f cafe-pos-backend | grep "PRINT\|BRIDGE"
```

Expect:
```
[PRINT] Using HTML for print bridge - order_id=xxx, bridge_url=https://print.tacafe.store
[BRIDGE PRINTER] Sending HTML to print bridge - printer: Bill Printer, IP: 192.168.x.x:9100
[BRIDGE PRINTER] HTML print successful via bridge - printer: Bill Printer
```

### 2. Test collect payment
- Tạo order mới
- Collect payment
- Kiểm tra print job được tạo với contentType = "html"
- Kiểm tra máy in có in ra không

### 3. Xem print bridge logs
```bash
docker logs -f local-print-bridge
```

Expect:
```
POST /render-and-print
Rendering HTML...
Printer IP: 192.168.x.x
Printer Port: 9100
Print successful
```

## Troubleshooting

### Nếu vẫn lỗi "no documents in result":
- Chạy `./fix-missing-templates.sh` để tạo templates (nhưng không cần nữa vì dùng HTML file)

### Nếu print bridge không available:
- Kiểm tra print bridge có chạy: `curl https://print.tacafe.store/health`
- Restart backend để nó reconnect: `docker restart cafe-pos-backend`

### Nếu không tìm thấy HTML template file:
- Kiểm tra file tồn tại: `ls backend/application/services/templates/bill_template_optimized.html`
- Đảm bảo file được copy vào Docker image (check Dockerfile)
