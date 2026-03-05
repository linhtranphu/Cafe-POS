# Quick Start: HTML Template Management

## Lỗi 404 - Backend chưa chạy

Nếu bạn gặp lỗi 404 khi test HTML templates, có nghĩa là backend chưa được start hoặc routes chưa được register.

## Bước 1: Start Backend

```bash
cd backend
go run main.go
```

Hoặc build và chạy:

```bash
cd backend
go build -o cafe-pos-server main.go
./cafe-pos-server
```

## Bước 2: Kiểm tra Backend Logs

Khi backend start, bạn sẽ thấy logs:

```
✅ Chromedp print handler initialized
✅ HTML template handler initialized
```

Nếu không thấy "HTML template handler initialized", có nghĩa là:
- Chromedp renderer failed to initialize
- Template file không tồn tại

## Bước 3: Verify Routes

Test endpoint bằng curl:

```bash
# Test GET template
curl http://localhost:8080/api/manager/html-templates/bill

# Nếu cần auth token
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/manager/html-templates/bill
```

## Bước 4: Start Frontend

```bash
cd frontend
npm run dev
```

## Bước 5: Test trong Browser

1. Mở http://localhost:5173/#/print-management
2. Click tab "Templates"
3. Click "🌐 HTML Template"
4. Template sẽ tự động load

## Troubleshooting

### Backend không start

**Lỗi:** `chromedp: failed to start browser`

**Giải pháp:**
```bash
# macOS
brew install chromium

# Linux
sudo apt-get install chromium-browser

# Windows
# Download từ https://www.chromium.org/
```

### Template file không tồn tại

**Lỗi:** `Failed to read template`

**Giải pháp:**
```bash
# Kiểm tra file tồn tại
ls -la backend/application/services/templates/bill_template_optimized.html

# Nếu không có, copy từ bill_template.html
cp backend/application/services/templates/bill_template.html \
   backend/application/services/templates/bill_template_optimized.html
```

### Routes không hoạt động

**Kiểm tra:**
1. Backend có log "✅ HTML template handler initialized"?
2. Routes có được register trong main.go?
3. Auth middleware có block request không?

**Debug:**
```bash
# Check backend logs
tail -f backend/logs/app.log

# Test without auth
curl -v http://localhost:8080/api/manager/html-templates/bill
```

### CORS errors

Nếu gặp CORS errors, kiểm tra backend CORS config:

```go
// main.go
router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
}))
```

## Testing Workflow

### 1. Load Template

```bash
curl http://localhost:8080/api/manager/html-templates/bill
```

Expected response:
```json
{
  "success": true,
  "content": "<!DOCTYPE html>...",
  "path": "./application/services/templates/bill_template_optimized.html",
  "filename": "bill_template_optimized.html"
}
```

### 2. Save Template

```bash
curl -X PUT http://localhost:8080/api/manager/html-templates/bill \
  -H "Content-Type: application/json" \
  -d '{"content": "<!DOCTYPE html>..."}'
```

Expected response:
```json
{
  "success": true,
  "message": "Template saved successfully",
  "backup": "./application/services/templates/bill_template_optimized.html.backup"
}
```

### 3. Test Print

```bash
curl -X POST http://localhost:8080/api/manager/html-templates/test-print \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "507f1f77bcf86cd799439011",
    "printer_ip": "192.168.1.115"
  }'
```

Expected response:
```json
{
  "success": true,
  "message": "Test print successful",
  "order_number": "20260222-095703-168"
}
```

### 4. Preview PNG

```bash
curl -X POST http://localhost:8080/api/manager/html-templates/preview \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "507f1f77bcf86cd799439011"
  }'
```

Expected response:
```json
{
  "success": true,
  "message": "Preview created successfully",
  "filename": "preview_html_template_20260222-095703-168.png",
  "order_number": "20260222-095703-168"
}
```

## Common Issues

### Issue: "Failed to preview: AxiosError 404"

**Cause:** Backend không chạy hoặc routes không được register

**Solution:**
1. Start backend: `cd backend && go run main.go`
2. Check logs cho "✅ HTML template handler initialized"
3. Verify routes: `curl http://localhost:8080/api/manager/html-templates/bill`

### Issue: "chromedp: failed to allocate"

**Cause:** Chrome/Chromium không được cài đặt

**Solution:**
```bash
# macOS
brew install chromium

# Linux
sudo apt-get install chromium-browser
```

### Issue: "Template file not found"

**Cause:** Template path sai hoặc file không tồn tại

**Solution:**
```bash
# Check file exists
ls backend/application/services/templates/bill_template_optimized.html

# Create if missing
cp backend/application/services/templates/bill_template.html \
   backend/application/services/templates/bill_template_optimized.html
```

### Issue: "Order not found"

**Cause:** Order ID không tồn tại trong database

**Solution:**
1. Tạo order mới trong POS
2. Hoặc lấy order ID từ database:
```bash
# MongoDB
mongo cafe_pos
db.orders.find().limit(1)
```

### Issue: "Printer connection failed"

**Cause:** Printer IP sai hoặc printer offline

**Solution:**
1. Ping printer: `ping 192.168.1.115`
2. Telnet test: `telnet 192.168.1.115 9100`
3. Check printer power và network

## Development Tips

### Hot Reload Backend

Sử dụng `air` để auto-reload backend khi code thay đổi:

```bash
# Install air
go install github.com/air-verse/air@latest

# Run with air
cd backend
air
```

### Debug Mode

Enable debug logs:

```bash
# Set environment variable
export DEBUG=true

# Run backend
go run main.go
```

### Test với Mock Data

Nếu chưa có orders trong database, tạo mock data:

```bash
# Run seed script
cd backend
go run cmd/seed/main.go
```

## Next Steps

1. ✅ Start backend
2. ✅ Verify logs
3. ✅ Test endpoints
4. ✅ Open frontend
5. ✅ Edit template
6. ✅ Test print

## Support

Nếu vẫn gặp vấn đề:
1. Check backend logs
2. Check browser console
3. Test với curl
4. Verify database connection
5. Check printer connection

## Summary

Để sử dụng HTML Template Management:

```bash
# Terminal 1: Backend
cd backend
go run main.go

# Terminal 2: Frontend
cd frontend
npm run dev

# Browser
# Open http://localhost:5173/#/print-management
# Click Templates → HTML Template
```

Đảm bảo:
- ✅ Backend running (port 8080)
- ✅ Frontend running (port 5173)
- ✅ Chromium installed
- ✅ Template file exists
- ✅ Database connected
- ✅ Printer accessible (for test print)
