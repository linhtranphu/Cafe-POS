# ✅ Backend Started Successfully!

## Service Status

```
✅ MongoDB:      Running on localhost:27017 (Replica Set: rs0)
✅ Backend:      Running on localhost:3000 (PID: 50459)
✅ Frontend:     Running on localhost:5173 (PID: 50479)
✅ Print Bridge: Running on localhost:3001 (Docker)
```

## Backend Initialization Logs

```
✅ MongoDB connected successfully
✅ WebSocket hub started
✅ Socket.IO server started
✅ HTML bill renderer initialized
✅ Cost recalculation worker pool started
✅ Print worker started
✅ Print cleanup job started
✅ Chromedp print handler initialized
✅ HTML template handler initialized
```

## Registered Routes

HTML Template Management routes đã được register:

```
GET    /api/manager/html-templates/bill
PUT    /api/manager/html-templates/bill
POST   /api/manager/html-templates/test-print
POST   /api/manager/html-templates/preview
```

Chromedp Print routes:

```
POST   /api/manager/chromedp-print/bill
GET    /api/manager/chromedp-print/preview/:order_id
```

## Access URLs

### Frontend
- Local: http://localhost:5173
- LAN: http://192.168.1.8:5173

### Backend
- API: http://localhost:3000
- Health: http://localhost:3000/health

### Print Bridge
- Health: http://localhost:3001/health

### MongoDB
- URI: mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin

## Testing HTML Template Management

### 1. Open Frontend

```
http://localhost:5173/#/print-management
```

### 2. Navigate to Templates Tab

Click: **Templates** → **🌐 HTML Template**

### 3. Features Available

- ✅ Load template from backend
- ✅ Edit HTML/CSS in editor
- ✅ Live preview with sample data
- ✅ Save template (with auto backup)
- ✅ Search and select orders
- ✅ Test print with real order
- ✅ Preview PNG with real order

## Expected Behavior

### Load Template

When you open HTML Template tab:
1. Frontend sends: `GET /api/manager/html-templates/bill`
2. Backend loads: `./application/services/templates/bill_template_optimized.html`
3. Template displays in editor
4. Preview renders automatically

### Edit & Save

When you edit and save:
1. Edit HTML/CSS in left panel
2. Preview updates automatically (debounced 500ms)
3. Click "💾 Lưu Template"
4. Backend creates backup: `bill_template_optimized.html.backup`
5. Backend saves new template
6. Success message displays

### Test Print

When you test print:
1. Enter printer IP (e.g., 192.168.1.115)
2. Search for order
3. Select order from list
4. Click "🖨️ Test Print"
5. Backend:
   - Fetches order from database
   - Renders HTML template with real data
   - Chromedp captures as PNG
   - Converts to ESC/POS
   - Sends to printer
6. Success/error message displays

### Preview PNG

When you preview:
1. Select order
2. Click "👁️ Preview PNG"
3. Backend creates PNG file
4. File saved: `preview_html_template_[order_number].png`
5. Success message with filename

## Authentication

API endpoints require authentication. Frontend automatically includes auth token from login session.

If testing with curl, you need to:
1. Login first to get token
2. Include token in requests:

```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.token')

# Use token
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/manager/html-templates/bill
```

## Troubleshooting

### Issue: Template not loading

**Check:**
1. Backend logs: `tail -f backend.log`
2. Look for: "✅ HTML template handler initialized"
3. Verify file exists: `ls backend/application/services/templates/bill_template_optimized.html`

**Solution:**
```bash
# If file missing, copy from bill_template.html
cp backend/application/services/templates/bill_template.html \
   backend/application/services/templates/bill_template_optimized.html
```

### Issue: Preview not working

**Check:**
1. Browser console for errors
2. Network tab for API calls
3. Backend logs for errors

**Common causes:**
- Template syntax error
- Missing data fields
- JavaScript error in preview processing

### Issue: Test print fails

**Check:**
1. Order exists in database
2. Printer IP is correct
3. Printer is online: `ping 192.168.1.115`
4. Port 9100 is accessible: `telnet 192.168.1.115 9100`

**Backend logs:**
```bash
tail -f backend.log | grep -i "print\|error"
```

### Issue: Chromedp errors

**Error:** `chromedp: failed to allocate`

**Solution:**
```bash
# macOS
brew install chromium

# Linux
sudo apt-get install chromium-browser
```

## Logs

### View Backend Logs
```bash
tail -f backend.log
```

### View Frontend Logs
```bash
tail -f frontend.log
```

### View MongoDB Logs
```bash
docker logs cafe-pos-mongodb
```

### View Print Bridge Logs
```bash
docker logs local-print-bridge
```

## Stop Services

```bash
# Stop backend
kill 50459

# Stop frontend
kill 50479

# Stop Print Bridge
docker stop local-print-bridge

# Stop MongoDB
docker-compose -f docker-compose.replica-set.yml down
```

## Restart Services

```bash
./restart_local.sh
```

## Next Steps

1. ✅ Open frontend: http://localhost:5173/#/print-management
2. ✅ Click Templates → HTML Template
3. ✅ Edit template and see live preview
4. ✅ Save template
5. ✅ Select an order
6. ✅ Test print or preview PNG

## Summary

Everything is running successfully! You can now:

- ✅ Edit HTML templates in the UI
- ✅ See live preview
- ✅ Test print with real orders
- ✅ Generate preview PNGs
- ✅ All changes are saved with automatic backup

The 404 error is now fixed because:
1. ✅ Backend is running
2. ✅ Routes are registered
3. ✅ HTML template handler is initialized
4. ✅ Template file exists

Just open the frontend and start editing templates!
