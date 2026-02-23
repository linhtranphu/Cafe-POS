# Troubleshooting 404 API Errors

## Problem

Getting 404 errors when calling API endpoints from the frontend:
```
GET http://localhost:5173/api/cashier-shifts/.../managed-funds 404 (Not Found)
```

## Quick Fix

**RESTART THE FRONTEND DEV SERVER**

```bash
# Stop current server (Ctrl+C in the terminal running npm run dev)
# Then restart:
cd frontend
npm run dev
```

The Vite configuration was updated to proxy to port 8080, but Vite needs to be restarted for changes to take effect.

---

## Detailed Troubleshooting

### Step 1: Verify Backend is Running

```bash
# Check if backend is running
curl http://localhost:8080/api/v1/health
# or
curl http://localhost:8080/api/health

# If you get connection refused, start the backend:
cd backend
go run main.go
```

### Step 2: Check Backend Port

Look at the backend terminal output:
```
Server running on :8080  ✅ Correct
Server running on :3000  ⚠️  Update vite.config.js to use port 3000
```

### Step 3: Verify Vite Proxy Configuration

Check `frontend/vite.config.js`:
```javascript
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',  // Should match backend port
      changeOrigin: true
    }
  }
}
```

### Step 4: Restart Frontend Dev Server

**IMPORTANT**: Vite only reads config on startup!

```bash
# Stop the dev server (Ctrl+C)
cd frontend
npm run dev
```

### Step 5: Clear Browser Cache

```bash
# In browser:
# 1. Open DevTools (F12)
# 2. Right-click refresh button
# 3. Select "Empty Cache and Hard Reload"
```

### Step 6: Verify API Calls

1. Open browser DevTools (F12)
2. Go to Network tab
3. Navigate to Cashier Dashboard
4. Look for API calls

**Before Fix**:
```
Request URL: http://localhost:5173/api/cashier-shifts/.../managed-funds
Status: 404 Not Found
```

**After Fix**:
```
Request URL: http://localhost:5173/api/v1/cashier-shifts/.../managed-funds
Status: 200 OK (proxied to http://localhost:8080)
```

---

## Common Issues

### Issue 1: Wrong API Version

**Symptom**: 404 errors even after restart

**Check**: API routes might use `/api/v1/` prefix

**Solution**: Check backend routes in `backend/main.go`:
```go
// If routes are defined as:
v1 := router.Group("/api/v1")

// Then frontend should call:
GET /api/v1/cashier-shifts/:id/managed-funds
```

**Fix**: Update frontend service to include `/v1`:
```javascript
// frontend/src/services/api.js
const api = axios.create({
  baseURL: '/api/v1',  // Add /v1 if backend uses it
})
```

### Issue 2: CORS Errors

**Symptom**: CORS policy errors in console

**Solution**: Backend should allow CORS from frontend origin:
```go
// backend/main.go
router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
}))
```

### Issue 3: Authentication Errors

**Symptom**: 401 Unauthorized

**Solution**: 
1. Check JWT token in localStorage
2. Login again to get fresh token
3. Verify token is being sent in Authorization header

```javascript
// Check in browser console:
localStorage.getItem('token')

// Should see token value
// If null, login again
```

### Issue 4: Route Not Found

**Symptom**: 404 for specific endpoint

**Check**: Verify route exists in backend:
```bash
# Search for route definition
cd backend
grep -r "managed-funds" .
```

**Verify**: Route should be registered in `backend/main.go`:
```go
cashierShifts.GET("/:id/managed-funds", cashierShiftClosureHandler.GetManagedFunds)
```

---

## Verification Checklist

After applying fixes:

- [ ] Backend running on correct port
- [ ] Frontend dev server restarted
- [ ] Browser cache cleared
- [ ] Network tab shows correct URLs
- [ ] API calls return 200 OK
- [ ] No CORS errors
- [ ] JWT token present
- [ ] Data displays correctly

---

## Still Not Working?

### Check Backend Logs

```bash
# Backend terminal should show:
[GIN] 2024/01/15 - 10:30:00 | 200 |    5.123ms |       127.0.0.1 | GET      "/api/v1/cashier-shifts/699c72a8b0e2b7ba12423012/managed-funds"
```

If you see 404 in backend logs, the route doesn't exist.
If you don't see the request at all, proxy isn't working.

### Test Backend Directly

```bash
# Get JWT token from browser localStorage
TOKEN="your_token_here"

# Test endpoint directly
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/cashier-shifts/699c72a8b0e2b7ba12423012/managed-funds

# Should return JSON data
```

### Check Network Tab Details

In browser DevTools > Network:
1. Click on failed request
2. Check "Headers" tab
3. Verify:
   - Request URL
   - Request Method
   - Status Code
   - Response

---

## Quick Commands

```bash
# Restart everything
cd backend && go run main.go &
cd frontend && npm run dev

# Check ports
lsof -i :8080  # Backend
lsof -i :5173  # Frontend

# Test API directly
curl http://localhost:8080/api/v1/health

# Check Vite config
cat frontend/vite.config.js | grep target
```

---

## Summary

Most 404 errors are caused by:
1. ❌ Frontend dev server not restarted after config change
2. ❌ Wrong backend port in proxy config
3. ❌ Missing `/v1` in API base URL
4. ❌ Backend not running
5. ❌ Route not registered in backend

**Solution**: Restart frontend dev server after fixing vite.config.js!
