# ✅ Print Bridge URL - Configurable from UI

## Thay đổi

Đã chuyển `PRINT_BRIDGE_URL` từ environment variable sang configurable field trong Settings UI.

## Trước đây

```bash
# Backend .env
PRINT_BRIDGE_URL=http://192.168.1.100:3001
```

**Vấn đề:**
- Phải SSH vào EC2 để thay đổi
- Phải restart backend sau khi thay đổi
- Không user-friendly

## Bây giờ

```
Settings UI → Print Configuration → Print Bridge URL
```

**Ưu điểm:**
- ✅ Config từ UI (http://localhost:5173/#/print-management)
- ✅ Không cần SSH vào EC2
- ✅ Không cần restart backend
- ✅ Test connection ngay trong UI
- ✅ Dynamic - fetch từ database mỗi lần sử dụng

## Files đã thay đổi

### 1. Frontend

**frontend/src/components/printing/ShopSettingsForm.vue:**
- ✅ Added `print_bridge_url` field
- ✅ Added "Test Connection" button
- ✅ Shows connection status (success/error)
- ✅ Validates URL format

**UI Fields:**
```vue
Print Bridge URL: [http://192.168.1.100:3001]
                  [🔍 Kiểm tra kết nối]
                  ✅ Kết nối thành công
```

### 2. Backend

**backend/domain/settings/shop_settings.go:**
- ✅ Field already exists: `PrintBridgeURL string`

**backend/interfaces/http/shop_settings_handler.go:**
- ✅ Already handles `print_bridge_url` in Create/Update

**backend/interfaces/http/html_template_handler_bridge.go:**
- ✅ Removed fixed `printBridgeClient` from struct
- ✅ Added `getPrintBridgeClient()` method - creates client dynamically from settings
- ✅ Updated `TestPrintHTMLTemplate()` - fetches URL from settings
- ✅ Updated `PreviewHTMLTemplate()` - fetches URL from settings

**backend/main.go:**
- ✅ Removed `PRINT_BRIDGE_URL` environment variable check
- ✅ Simplified initialization - no connection test at startup
- ✅ Handler always registered (will check settings at runtime)

## Flow hoạt động

### 1. User configures Print Bridge URL

```
1. User opens Settings → Print Management
2. Enters Print Bridge URL: http://192.168.1.100:3001
3. Clicks "Test Connection"
   → Frontend calls: GET http://192.168.1.100:3001/health
   → Shows: ✅ Kết nối thành công
4. Clicks "Save"
   → POST /api/manager/shop-settings
   → Saves to MongoDB
```

### 2. User tests HTML template printing

```
1. User clicks "Test Print" in Settings
2. Frontend calls: POST /api/manager/html-templates/test-print
3. Backend handler:
   a. Fetches shop_settings from MongoDB
   b. Gets print_bridge_url from settings
   c. Creates print bridge client dynamically
   d. Checks if print bridge is available
   e. Renders HTML → ESC/POS
   f. Sends to printer
```

### 3. Dynamic client creation

```go
// html_template_handler_bridge.go

func (h *HTMLTemplateHandlerBridge) getPrintBridgeClient(ctx context.Context) (*printbridge.Client, error) {
    // Fetch settings from database
    shopSettings, err := h.shopSettingsRepo.GetSettings(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get shop settings: %w", err)
    }

    // Check if URL is configured
    if shopSettings.PrintBridgeURL == "" {
        return nil, fmt.Errorf("print bridge URL not configured in settings")
    }

    // Create client with URL from settings
    return printbridge.NewClient(shopSettings.PrintBridgeURL, 30*time.Second), nil
}
```

## Testing

### 1. Test connection from UI

```
1. Go to: http://localhost:5173/#/print-management
2. Enter Print Bridge URL: http://192.168.1.100:3001
3. Click "Kiểm tra kết nối"
4. Should see: ✅ Kết nối thành công
```

### 2. Test HTML template printing

```
1. Make sure Print Bridge is running on local machine
2. Configure Print Bridge URL in Settings
3. Click "Test Print"
4. Should print successfully
```

### 3. Test error handling

```
# Wrong URL
Print Bridge URL: http://192.168.1.999:3001
→ Click Test → ❌ Không thể kết nối

# Empty URL
Print Bridge URL: (empty)
→ Click Test Print → Error: "Print bridge not configured"

# Print Bridge not running
Print Bridge URL: http://192.168.1.100:3001 (but service stopped)
→ Click Test Print → Error: "Print bridge is not available"
```

## Migration

### Existing deployments with PRINT_BRIDGE_URL in .env

**Option 1: Keep environment variable (backward compatible)**
```bash
# Backend will still work with PRINT_BRIDGE_URL in .env
# But it's ignored - settings take priority
```

**Option 2: Migrate to settings**
```bash
# 1. Get current PRINT_BRIDGE_URL
echo $PRINT_BRIDGE_URL

# 2. Login to UI and set in Settings
# 3. Remove from .env (optional)
sed -i '/PRINT_BRIDGE_URL/d' .env
```

### New deployments

```bash
# No need to set PRINT_BRIDGE_URL in .env
# Just configure in UI after deployment
```

## Deployment Steps

### 1. Deploy Backend

```bash
# Build new image
cd backend
docker build -t linhtranphu/cafe-pos-backend:latest .
docker push linhtranphu/cafe-pos-backend:latest

# On EC2
ssh ubuntu@tacafe.store
cd ~/cafe-pos
docker-compose pull backend
docker-compose up -d backend

# Check logs
docker logs -f cafe-pos-backend
# Should see: "✅ HTML template handler (bridge) initialized - will use print_bridge_url from settings"
```

### 2. Deploy Frontend

```bash
# Build new image
cd frontend
docker build -t linhtranphu/cafe-pos-frontend:latest .
docker push linhtranphu/cafe-pos-frontend:latest

# On EC2
docker-compose pull frontend
docker-compose up -d frontend
```

### 3. Configure Print Bridge URL

```bash
# 1. Start Print Bridge on local machine
cd local-print-bridge
npm start

# 2. Get local machine IP
ifconfig | grep "inet " | grep -v 127.0.0.1
# Example: 192.168.1.100

# 3. Open browser
https://tacafe.store

# 4. Login → Settings → Print Management

# 5. Enter Print Bridge URL
http://192.168.1.100:3001

# 6. Click "Test Connection"
# Should see: ✅ Kết nối thành công

# 7. Click "Save"
```

## Benefits

### For Users
- ✅ Easy to configure from UI
- ✅ No need to SSH into server
- ✅ Test connection before saving
- ✅ Visual feedback (success/error)

### For Developers
- ✅ No hardcoded URLs
- ✅ Dynamic configuration
- ✅ Better error messages
- ✅ Easier to debug

### For Operations
- ✅ No server restart needed
- ✅ Can change URL anytime
- ✅ Centralized configuration
- ✅ Audit trail in database

## Troubleshooting

### "Print bridge not configured"

**Cause:** `print_bridge_url` not set in settings

**Fix:**
1. Go to Settings → Print Management
2. Enter Print Bridge URL
3. Click Save

### "Print bridge is not available"

**Cause:** Print Bridge service not running or unreachable

**Fix:**
1. Check Print Bridge is running: `curl http://192.168.1.100:3001/health`
2. Check firewall allows port 3001
3. Check IP address is correct
4. Check network connectivity

### 404 on /api/manager/html-templates/bill

**Cause:** Backend not updated or not started properly

**Fix:**
1. Check backend logs: `docker logs cafe-pos-backend`
2. Should see: "✅ HTML template handler (bridge) initialized"
3. If not, rebuild and restart backend

## Summary

✅ **Removed:**
- PRINT_BRIDGE_URL environment variable dependency
- Fixed print bridge client at startup
- Connection test at startup

✅ **Added:**
- Print Bridge URL field in Settings UI
- Test Connection button
- Dynamic client creation from settings
- Better error messages

✅ **Benefits:**
- User-friendly configuration
- No server restart needed
- Real-time connection testing
- Centralized settings management

---

**Status:** Ready for deployment  
**UI:** http://localhost:5173/#/print-management  
**Field:** Print Bridge URL (with test button)
