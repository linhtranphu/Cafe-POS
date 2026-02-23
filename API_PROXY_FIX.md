# API Proxy Configuration Fix

## Issue

Frontend was getting 404 errors when calling API endpoints:
```
GET http://localhost:5173/api/cashier-shifts/699c72a8b0e2b7ba12423012/managed-funds 404 (Not Found)
```

## Root Cause

The Vite proxy configuration in `frontend/vite.config.js` was pointing to the wrong backend port:
- **Configured**: `http://localhost:3000`
- **Actual Backend**: `http://localhost:8080`

## Fix Applied

Updated `frontend/vite.config.js`:

```javascript
server: {
  port: 5173,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',  // Changed from 3000 to 8080
      changeOrigin: true,
      rewrite: (path) => path
    }
  }
}
```

## How to Apply

### Option 1: Restart Frontend Dev Server (Recommended)

```bash
# Stop the current dev server (Ctrl+C)
# Then restart it
cd frontend
npm run dev
```

The Vite dev server needs to be restarted for the proxy configuration changes to take effect.

### Option 2: Verify Backend Port

If your backend is actually running on port 3000, you don't need to change anything. Just verify:

```bash
# Check what port your backend is running on
# Look for output like: "Server running on :8080"
```

## Verification

After restarting the frontend dev server:

1. Open browser console (F12)
2. Navigate to Cashier Dashboard
3. Check Network tab
4. API calls should now go to `http://localhost:8080/api/...`
5. Should return 200 OK instead of 404

## Testing

```bash
# 1. Make sure backend is running
cd backend
go run main.go
# Should see: Server running on :8080

# 2. Restart frontend (in another terminal)
cd frontend
npm run dev
# Should see: Local: http://localhost:5173

# 3. Test in browser
# Navigate to http://localhost:5173
# Login as cashier
# Go to Cashier Dashboard
# Should see managed funds without 404 errors
```

## Additional Notes

### Environment Variables

For production, you may want to use environment variables:

```javascript
// vite.config.js
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: process.env.VITE_API_URL || 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
```

Then create `.env.local`:
```
VITE_API_URL=http://localhost:8080
```

### Production Build

For production, the API calls should go directly to your backend server. Make sure to configure your web server (nginx, Apache, etc.) to proxy `/api` requests to your backend.

Example nginx configuration:
```nginx
location /api {
    proxy_pass http://localhost:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

## Status

✅ **FIXED** - Vite proxy now points to correct backend port (8080)

**Next Step**: Restart frontend dev server to apply changes
