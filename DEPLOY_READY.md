# ✅ Backend đã build thành công!

## Thay đổi đã hoàn tất

### 1. Print Service sử dụng HTML + Print Bridge
- `backend/application/services/print_service.go`
  - Sửa `createBillJob()` để render HTML khi có print bridge URL
  - Thêm `renderBillHTML()` để render HTML từ template
  - Thêm `prepareBillData()` để chuẩn bị data
  - Sử dụng `formatMoneyVN()` đã có sẵn

### 2. Bridge Printer xử lý HTML
- `backend/infrastructure/printing/bridge_printer.go`
  - Sửa `Print()` để phát hiện và xử lý HTML content
  - Gọi `RenderAndPrint()` cho HTML (print bridge sẽ convert)
  - Gọi `PrintESCPOS()` cho binary/text như cũ
  - Thêm `isHTMLContent()` helper
  - Sử dụng `isBase64Content()` đã có trong `helpers.go`

## Deploy lên EC2

### Option 1: Git Push (Khuyến nghị)
```bash
# Local
git add backend/
git commit -m "fix: use HTML template with print bridge for auto-print"
git push origin main

# EC2
ssh ec2-user@<EC2-IP>
cd ~/Cafe\ POS
git pull
cd backend
docker build -t cafe-pos-backend:latest .
docker stop cafe-pos-backend
docker rm cafe-pos-backend
docker run -d \
  --name cafe-pos-backend \
  --network cafe-pos-network \
  -p 8080:8080 \
  --env-file ../.env.ec2 \
  --restart unless-stopped \
  cafe-pos-backend:latest

# Xem logs
docker logs -f cafe-pos-backend
```

### Option 2: Upload Docker Image
```bash
# Local - Save image
docker save cafe-pos-backend:latest | gzip > backend-image.tar.gz

# Upload
scp backend-image.tar.gz ec2-user@<EC2-IP>:~/

# EC2 - Load và run
ssh ec2-user@<EC2-IP>
docker load < backend-image.tar.gz
docker stop cafe-pos-backend && docker rm cafe-pos-backend
docker run -d \
  --name cafe-pos-backend \
  --network cafe-pos-network \
  -p 8080:8080 \
  --env-file .env.ec2 \
  --restart unless-stopped \
  cafe-pos-backend:latest
```

## Kiểm tra sau deploy

### 1. Backend health
```bash
curl http://localhost:8080/health
```

### 2. Xem logs
```bash
docker logs -f cafe-pos-backend | grep "PRINT\|BRIDGE"
```

Expect thấy:
```
[PRINT BRIDGE] Configuring print bridge: https://print.tacafe.store
[PRINT] Using HTML for print bridge - order_id=xxx, bridge_url=https://print.tacafe.store
[BRIDGE PRINTER] Sending HTML to print bridge - printer: Bill Printer, IP: 192.168.x.x:9100
[BRIDGE PRINTER] HTML print successful via bridge
```

### 3. Test collect payment
1. Tạo order mới trên frontend
2. Collect payment
3. Kiểm tra backend logs
4. Kiểm tra máy in có in ra không

### 4. Xem print bridge logs (trên máy local)
```bash
docker logs -f local-print-bridge
```

Expect:
```
POST /render-and-print
Rendering HTML (576px width)...
Printer IP: 192.168.x.x
Printer Port: 9100
✓ Print successful
```

## Troubleshooting

### Nếu backend không start:
```bash
docker logs cafe-pos-backend
```

### Nếu không tìm thấy HTML template:
```bash
# Kiểm tra file có trong image không
docker exec cafe-pos-backend ls -la /root/application/services/templates/
```

Nếu không có, kiểm tra Dockerfile có COPY đúng không.

### Nếu vẫn không in được:
1. Chạy diagnostic scripts:
```bash
./test-full-print-flow.sh
./check-print-worker-detail.sh
```

2. Kiểm tra print bridge URL:
```bash
docker exec cafe-pos-mongodb mongosh cafe_pos_db --eval 'db.shop_settings.findOne({}, {print_bridge_url: 1})'
```

Phải là: `https://print.tacafe.store`

3. Test print bridge:
```bash
curl https://print.tacafe.store/health
```

## Flow hoàn chỉnh

```
User: Collect Payment
  ↓
Backend: CreatePrintJobsForOrder()
  ↓
Backend: createBillJob()
  → Render HTML từ template
  → Save print job (content = HTML, contentType = "html", status = PENDING)
  ↓
Print Worker (background, mỗi 10s)
  → Fetch pending jobs
  → ProcessJob()
  → GetPrinter() → BridgePrinter
  → printer.Print(htmlContent)
  → Detect HTML content
  → bridgeClient.RenderAndPrint(html, printerIP, printerPort, 576)
  ↓
Print Bridge (https://print.tacafe.store)
  → Receive HTML
  → Render HTML → Image (chromedp)
  → Convert Image → ESC/POS
  → Send to printer (192.168.x.x:9100)
  ↓
Máy in: In bill
  ↓
Backend: Update job status = COMPLETED
```

## Lưu ý

- Đảm bảo file `backend/application/services/templates/bill_template_optimized.html` tồn tại
- Print Bridge phải chạy và accessible từ EC2
- Printer config phải có IP và port đúng
- Shop settings phải có `print_bridge_url = https://print.tacafe.store`
