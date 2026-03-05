# Print Bridge URL Configuration - Implementation Complete

## Overview
Implemented dynamic print bridge URL configuration that allows users to configure the local print bridge URL through the UI instead of environment variables.

## Architecture
- **Frontend/Backend**: EC2 (tacafe.store)
- **Print Bridge**: Local machine at cafe (same network as printer)
- **Configuration**: Stored in backend `shop_settings.print_bridge_url`

## Changes Made

### 1. Backend (Already Completed)
- ✅ Added `PrintBridgeURL` field to `ShopSettings` struct
- ✅ Updated `UpdateSettingsRequest` to include `print_bridge_url`
- ✅ Updated `UpdateSettings` handler to save `print_bridge_url`

### 2. Frontend Utilities
**Created: `frontend/src/utils/env.js`**
- Utility functions for environment variable access
- `getApiUrl()`: Get API URL with fallback
- `getPrintBridgeUrl()`: Get print bridge URL with fallback
- `isDevelopment()`, `isProduction()`: Environment checks

### 3. Local Print Service
**Updated: `frontend/src/services/localPrint.js`**
- ✅ Reads print bridge URL from localStorage (priority 1) or environment (priority 2)
- ✅ Added `updateBridgeUrl(newUrl)`: Dynamically update bridge URL
- ✅ Added `getBridgeUrl()`: Get current bridge URL
- ✅ URL changes trigger availability recheck

### 4. Shop Settings Form
**Updated: `frontend/src/components/printing/ShopSettingsForm.vue`**
- ✅ Added input field for `print_bridge_url`
- ✅ Field includes helpful placeholder and description
- ✅ Form loads `print_bridge_url` from settings
- ✅ Form saves `print_bridge_url` to backend
- ✅ On save, updates `localPrintService` with new URL

### 5. App Initialization
**Updated: `frontend/src/main.js`**
- ✅ Fetches shop settings on app load (when authenticated)
- ✅ Initializes print bridge URL from settings
- ✅ Calls `localPrintService.updateBridgeUrl()` with configured URL

### 6. Database Migration
**Updated: `migrate-v2.0-mongodb.sh`**
- ✅ Creates `shop_settings` with default `print_bridge_url: "http://localhost:3001"`
- ✅ Updates existing settings to add `print_bridge_url` if missing
- ✅ Migration is idempotent and safe for production

## Configuration Flow

### Initial Setup
1. Run migration script to create/update shop_settings
2. Default print_bridge_url is set to `http://localhost:3001`

### User Configuration
1. User navigates to Print Management → Shop Settings
2. User enters print bridge URL (e.g., `http://192.168.1.100:3001`)
3. User clicks "Lưu cài đặt" (Save Settings)
4. URL is saved to backend database
5. `localPrintService` is updated immediately with new URL
6. Print bridge availability is rechecked

### App Load
1. User opens app and logs in
2. App fetches shop settings from backend
3. If `print_bridge_url` exists, updates `localPrintService`
4. URL is stored in localStorage for persistence
5. Print bridge health checks begin

## Priority Order
The print bridge URL is resolved in this order:
1. **Configured URL** (from shop_settings, stored in localStorage)
2. **Environment Variable** (`VITE_PRINT_BRIDGE_URL`)
3. **Default Fallback** (`http://localhost:3001`)

## Testing Checklist

### Development
- [ ] Run migration script: `sudo bash migrate-v2.0-mongodb.sh`
- [ ] Start frontend: `npm run dev`
- [ ] Login and navigate to Print Management → Shop Settings
- [ ] Verify print_bridge_url field is visible
- [ ] Enter a test URL (e.g., `http://192.168.1.50:3001`)
- [ ] Save settings
- [ ] Check browser console for: `[LocalPrint] Bridge URL updated: ...`
- [ ] Reload page
- [ ] Verify URL persists after reload

### Production Deployment
- [ ] Run migration on production: `sudo bash migrate-v2.0-mongodb.sh`
- [ ] Rebuild frontend: `npm run build`
- [ ] Deploy to EC2
- [ ] Login to production app
- [ ] Configure actual print bridge URL (local IP at cafe)
- [ ] Test print functionality

## Files Modified
1. `frontend/src/utils/env.js` (created)
2. `frontend/src/services/localPrint.js` (updated)
3. `frontend/src/components/printing/ShopSettingsForm.vue` (updated)
4. `frontend/src/main.js` (updated)
5. `migrate-v2.0-mongodb.sh` (updated)

## Benefits
- ✅ No need to rebuild/redeploy to change print bridge URL
- ✅ URL can be configured through UI by non-technical users
- ✅ Centralized configuration in backend database
- ✅ Automatic URL updates without page reload
- ✅ Fallback to environment variables for development
- ✅ URL persists across page reloads via localStorage

## Next Steps
1. Test the complete flow in development
2. Run migration script on production
3. Deploy updated frontend to EC2
4. Configure actual print bridge URL at cafe
5. Test printing with real printer

## Notes
- Print bridge must be running on local machine at cafe
- URL should be accessible from browser (same network or tunnel)
- Example URLs:
  - Development: `http://localhost:3001`
  - Production (local network): `http://192.168.1.100:3001`
  - Production (tunnel): `https://print-bridge.ngrok.io`
