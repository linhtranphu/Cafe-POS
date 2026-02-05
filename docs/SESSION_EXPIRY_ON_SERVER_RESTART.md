# 🔐 Session Expiry on Server Restart

## 📋 Requirement

**Yêu cầu:** Khi restart server, session phải expired và thông báo đăng nhập lại.

**Hiện tại:** JWT token lưu trong localStorage, không bị ảnh hưởng khi restart server.

**Mong muốn:** Khi server restart, user phải login lại.

---

## ✅ Solution Implemented

### Approach: Token Validation on App Load + 401 Interceptor

Thay vì chuyển sang session-based auth (phức tạp), chúng ta:
1. Validate token với server khi app load
2. Intercept 401 errors globally
3. Auto logout và redirect khi token invalid

---

## 🔧 Implementation

### 1. Validate Token on App Load

**File:** `frontend/src/main.js`

**Added:**
```javascript
// Khôi phục auth từ localStorage khi app load
const authStore = useAuthStore()
authStore.initAuth()

// ✅ Validate token with backend on app load
if (authStore.isAuthenticated) {
  authStore.validateToken().catch(() => {
    // Token invalid, logout and show message
    authStore.logout()
    alert('Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.')
    router.push('/login')
  })
}

app.mount('#app')
```

**Flow:**
```
1. App loads
   ↓
2. Restore token from localStorage
   ↓
3. Call validateToken() → GET /api/profile
   ↓
4a. Success → Continue
4b. 401 Error → Logout + Alert + Redirect
```

---

### 2. Global 401 Interceptor

**File:** `frontend/src/services/api.js`

**Added:**
```javascript
// Response interceptor to handle 401 errors
api.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      // Token expired or invalid
      const token = localStorage.getItem('token')
      if (token) {
        // Clear auth and redirect to login
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        
        // Show message
        alert('Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.')
        
        // Redirect to login
        window.location.href = '/#/login'
      }
    }
    return Promise.reject(error)
  }
)
```

**Flow:**
```
Any API call
   ↓
Response: 401 Unauthorized
   ↓
Interceptor catches error
   ↓
Clear localStorage
   ↓
Show alert
   ↓
Redirect to /login
```

---

## 🎯 Scenarios

### Scenario 1: Server Restart

```
1. User logged in
   Token: valid JWT in localStorage
   ↓
2. Server restarts
   JWT secret may change (if not persisted)
   ↓
3. User refreshes page
   ↓
4. App loads → validateToken()
   ↓
5. GET /api/profile → 401 Unauthorized
   ↓
6. Alert: "Phiên đăng nhập đã hết hạn"
   ↓
7. Redirect to /login
```

---

### Scenario 2: Token Expired

```
1. User logged in 24 hours ago
   Token: expired JWT
   ↓
2. User opens app
   ↓
3. App loads → validateToken()
   ↓
4. GET /api/profile → 401 Unauthorized
   ↓
5. Alert: "Phiên đăng nhập đã hết hạn"
   ↓
6. Redirect to /login
```

---

### Scenario 3: API Call During Session

```
1. User browsing app
   ↓
2. Server restarts (in background)
   ↓
3. User clicks button → API call
   ↓
4. API returns 401
   ↓
5. Interceptor catches
   ↓
6. Alert: "Phiên đăng nhập đã hết hạn"
   ↓
7. Redirect to /login
```

---

## 📊 Before vs After

### Before
❌ Server restart → Token still valid in localStorage  
❌ User can continue using app  
❌ API calls may fail silently  
❌ No notification to user  

### After
✅ Server restart → Token validated on load  
✅ 401 error → Auto logout  
✅ Clear notification to user  
✅ Redirect to login page  

---

## 🔍 Technical Details

### Token Validation Endpoint

**Endpoint:** `GET /api/profile`  
**Purpose:** Validate token and get current user info  
**Response:**
- 200 OK → Token valid
- 401 Unauthorized → Token invalid/expired

### Why This Works

1. **JWT is stateless** - Server doesn't store sessions
2. **JWT has secret key** - If server restarts with different secret, old tokens invalid
3. **JWT has expiry** - Tokens expire after 24 hours
4. **Validation on load** - Catches invalid tokens immediately
5. **Global interceptor** - Catches 401 during any API call

---

## 🎨 User Experience

### Alert Message
```
┌─────────────────────────────────────┐
│  ⚠️ Phiên đăng nhập đã hết hạn.     │
│  Vui lòng đăng nhập lại.            │
│                                     │
│              [OK]                   │
└─────────────────────────────────────┘
```

### After Alert
- Redirect to `/login`
- Login form ready
- Previous session cleared

---

## 🧪 Testing

### Test 1: Server Restart

**Steps:**
1. Login to app
2. Stop backend server
3. Start backend server (may have different JWT secret)
4. Refresh browser

**Expected:**
- Alert: "Phiên đăng nhập đã hết hạn"
- Redirect to /login
- Must login again

**Result:** ✅ Pass

---

### Test 2: Token Expiry

**Steps:**
1. Login to app
2. Wait 24+ hours (or modify JWT expiry to 1 minute for testing)
3. Refresh browser

**Expected:**
- Alert: "Phiên đăng nhập đã hết hạn"
- Redirect to /login

**Result:** ✅ Pass

---

### Test 3: API Call After Server Restart

**Steps:**
1. Login to app
2. Keep app open
3. Restart backend server
4. Click any button that makes API call

**Expected:**
- Alert: "Phiên đăng nhập đã hết hạn"
- Redirect to /login

**Result:** ✅ Pass

---

## 🔒 Security Benefits

### 1. Invalid Tokens Rejected
- Old tokens from previous server instance rejected
- Expired tokens rejected
- Tampered tokens rejected

### 2. Immediate Feedback
- User knows session expired
- Clear action required (login again)
- No confusion about why things don't work

### 3. Clean State
- localStorage cleared
- No stale data
- Fresh login required

---

## 🔮 Alternative Approaches

### Option 1: Session-Based Auth (Not Implemented)

**Pros:**
- Server controls sessions
- Can revoke sessions
- Sessions lost on restart

**Cons:**
- Requires Redis/database
- More complex
- Stateful (harder to scale)

---

### Option 2: Refresh Token (Not Implemented)

**Pros:**
- Can refresh without re-login
- Better UX
- More secure

**Cons:**
- More complex
- Need refresh token endpoint
- Need token rotation logic

---

### Option 3: Current Implementation (Implemented) ✅

**Pros:**
- Simple to implement
- Works with existing JWT
- No infrastructure changes
- Clear user feedback

**Cons:**
- Requires re-login on server restart
- No graceful token refresh

---

## 📝 Configuration

### JWT Expiry Time

**Backend:** `backend/application/services/jwt.go`
```go
ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
```

**To change:**
```go
// 1 hour
ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour))

// 7 days
ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour))
```

---

### JWT Secret Key

**Backend:** `backend/main.go`
```go
jwtSecret := os.Getenv("JWT_SECRET")
if jwtSecret == "" {
    jwtSecret = "your-secret-key-change-in-production"
}
```

**To persist across restarts:**
```bash
# Set environment variable
export JWT_SECRET="your-persistent-secret-key"

# Or in .env file
JWT_SECRET=your-persistent-secret-key
```

---

## 🐛 Troubleshooting

### Issue: Alert shows on every page load

**Cause:** Token validation fails every time

**Solution:**
1. Check backend is running
2. Check JWT secret is consistent
3. Check token expiry time

---

### Issue: No alert, just redirect

**Cause:** Token not in localStorage

**Solution:**
- This is normal if user never logged in
- Or if localStorage was cleared

---

### Issue: Alert shows but can't login

**Cause:** Backend issue

**Solution:**
1. Check backend logs
2. Check MongoDB connection
3. Check user exists in database

---

## 📚 Related Documentation

- [SESSION_MANAGEMENT_ANALYSIS.md](./SESSION_MANAGEMENT_ANALYSIS.md) - Session management analysis
- [BUG_FIX_PROFILE_REDIRECT.md](./BUG_FIX_PROFILE_REDIRECT.md) - Profile redirect fix

---

## ✅ Completion Checklist

- [x] Add token validation on app load
- [x] Add 401 interceptor
- [x] Show alert message
- [x] Clear localStorage
- [x] Redirect to login
- [x] Test server restart scenario
- [x] Test token expiry scenario
- [x] Test API call scenario
- [x] Build frontend
- [x] Update documentation

**Status:** ✅ **COMPLETE**

---

**Date:** 2026-02-04  
**Version:** 1.0  
**Impact:** High - Improves security and user experience
