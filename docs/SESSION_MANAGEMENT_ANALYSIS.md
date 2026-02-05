# 🔐 Phân Tích Session Management & Auto-Expire

## 📋 Tình Trạng Hiện Tại

### Backend (Go)

**JWT Service** (`backend/application/services/jwt.go`)
```go
// Token expires sau 24 giờ
ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
```

**Cách hoạt động:**
1. User login → Backend tạo JWT token
2. Token có `ExpiresAt` = 24 giờ từ lúc tạo
3. Token được gửi về frontend
4. Backend validate token mỗi request
5. Nếu token expired → trả về 401 Unauthorized

**Ưu điểm:**
- ✅ Stateless (không cần lưu session trên server)
- ✅ Scalable (dễ scale horizontal)
- ✅ Token tự động expire sau 24h

**Nhược điểm:**
- ❌ Không thể revoke token trước khi expire
- ❌ Không có refresh token
- ❌ Không track activity (idle timeout)
- ❌ Token còn valid ngay cả khi user logout

---

### Frontend (Vue.js)

**Auth Store** (`frontend/src/stores/auth.js`)
```javascript
// Lưu token vào localStorage
localStorage.setItem('token', response.token)
localStorage.setItem('user', JSON.stringify(response.user))
```

**Cách hoạt động:**
1. Login thành công → Lưu token vào localStorage
2. App reload → Restore token từ localStorage
3. Mỗi API request → Gửi token trong header
4. Nếu API trả về 401 → Logout và redirect to login

**Ưu điểm:**
- ✅ Persistent login (token còn sau khi đóng browser)
- ✅ Tự động restore session khi reload

**Nhược điểm:**
- ❌ Token không tự động expire trên frontend
- ❌ Không có idle timeout
- ❌ Không có activity tracking
- ❌ Token có thể bị đánh cắp từ localStorage

---

## 🎯 Yêu Cầu Mới: Auto-Expire Session

### Các Loại Timeout

1. **Absolute Timeout** - Token expire sau X giờ (đã có)
2. **Idle Timeout** - Session expire sau X phút không hoạt động (chưa có)
3. **Refresh Token** - Gia hạn token trước khi expire (chưa có)

---

## 💡 Giải Pháp Đề Xuất

### Option 1: JWT với Idle Timeout (Recommended)

**Ưu điểm:**
- Giữ được stateless architecture
- Dễ implement
- Phù hợp với hệ thống hiện tại

**Cách hoạt động:**
1. Backend: Token expire sau 24h (absolute)
2. Frontend: Track last activity time
3. Frontend: Auto logout sau 30 phút không hoạt động (idle)
4. Frontend: Có thể refresh token trước khi expire

**Implementation:**

#### Backend Changes

```go
// jwt.go - Thêm refresh token
func (j *JWTService) GenerateTokenPair(u *user.User) (*TokenPair, error) {
    // Access token - expire sau 1 giờ
    accessClaims := Claims{
        UserID:   u.ID.Hex(),
        Username: u.Username,
        Role:     u.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessTokenString, err := accessToken.SignedString([]byte(j.secretKey))
    if err != nil {
        return nil, err
    }
    
    // Refresh token - expire sau 7 ngày
    refreshClaims := Claims{
        UserID:   u.ID.Hex(),
        Username: u.Username,
        Role:     u.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshTokenString, err := refreshToken.SignedString([]byte(j.secretKey))
    if err != nil {
        return nil, err
    }
    
    return &TokenPair{
        AccessToken:  accessTokenString,
        RefreshToken: refreshTokenString,
        ExpiresIn:    3600, // 1 hour in seconds
    }, nil
}

type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int    `json:"expires_in"`
}
```

```go
// auth.go - Update login response
func (a *AuthService) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
    // ... existing validation ...
    
    tokenPair, err := a.jwtService.GenerateTokenPair(u)
    if err != nil {
        return nil, err
    }
    
    return &user.LoginResponse{
        AccessToken:  tokenPair.AccessToken,
        RefreshToken: tokenPair.RefreshToken,
        ExpiresIn:    tokenPair.ExpiresIn,
        User:         *u,
    }, nil
}

// Thêm endpoint refresh token
func (a *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
    claims, err := a.jwtService.ValidateToken(refreshToken)
    if err != nil {
        return nil, errors.New("invalid refresh token")
    }
    
    // Get user from database
    userID, _ := primitive.ObjectIDFromHex(claims.UserID)
    u, err := a.userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, errors.New("user not found")
    }
    
    if !u.Active {
        return nil, errors.New("user is inactive")
    }
    
    // Generate new token pair
    return a.jwtService.GenerateTokenPair(u)
}
```

#### Frontend Changes

```javascript
// stores/auth.js - Thêm idle timeout tracking
import { defineStore } from 'pinia'
import { login as authLogin } from '../services/auth'
import api from '../services/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    accessToken: null,
    refreshToken: null,
    expiresIn: null,
    isAuthenticated: false,
    loading: false,
    error: null,
    lastActivityTime: null,
    idleTimeoutMinutes: 30, // Configurable
    activityCheckInterval: null,
    tokenRefreshInterval: null
  }),

  actions: {
    async login(credentials) {
      this.loading = true
      this.error = null
      try {
        const response = await authLogin(credentials)
        
        if (response && response.user && response.access_token) {
          this.user = response.user
          this.accessToken = response.access_token
          this.refreshToken = response.refresh_token
          this.expiresIn = response.expires_in
          this.isAuthenticated = true
          this.lastActivityTime = Date.now()
          
          // Lưu vào localStorage
          localStorage.setItem('accessToken', response.access_token)
          localStorage.setItem('refreshToken', response.refresh_token)
          localStorage.setItem('user', JSON.stringify(response.user))
          localStorage.setItem('expiresIn', response.expires_in)
          localStorage.setItem('lastActivityTime', Date.now())
          
          // Set token cho API requests
          api.defaults.headers.common['Authorization'] = `Bearer ${response.access_token}`
          
          // Start activity tracking
          this.startActivityTracking()
          
          // Start token refresh timer
          this.startTokenRefresh()
          
          return true
        } else {
          throw new Error('Invalid response format')
        }
      } catch (error) {
        console.error('Login error:', error)
        this.error = error.response?.data?.error || error.message || 'Đăng nhập thất bại'
        return false
      } finally {
        this.loading = false
      }
    },

    // Track user activity
    updateActivity() {
      this.lastActivityTime = Date.now()
      localStorage.setItem('lastActivityTime', Date.now())
    },

    // Start tracking user activity
    startActivityTracking() {
      // Clear existing interval
      if (this.activityCheckInterval) {
        clearInterval(this.activityCheckInterval)
      }

      // Track user interactions
      const events = ['mousedown', 'keydown', 'scroll', 'touchstart', 'click']
      events.forEach(event => {
        window.addEventListener(event, this.updateActivity)
      })

      // Check idle timeout every minute
      this.activityCheckInterval = setInterval(() => {
        this.checkIdleTimeout()
      }, 60000) // Check every 1 minute
    },

    // Check if user has been idle too long
    checkIdleTimeout() {
      if (!this.isAuthenticated) return

      const now = Date.now()
      const lastActivity = this.lastActivityTime || parseInt(localStorage.getItem('lastActivityTime'))
      const idleTime = (now - lastActivity) / 1000 / 60 // minutes

      if (idleTime >= this.idleTimeoutMinutes) {
        console.log(`Session expired due to inactivity (${idleTime.toFixed(1)} minutes)`)
        this.logout('Session expired due to inactivity')
      }
    },

    // Start automatic token refresh
    startTokenRefresh() {
      // Clear existing interval
      if (this.tokenRefreshInterval) {
        clearInterval(this.tokenRefreshInterval)
      }

      // Refresh token 5 minutes before expiry
      const refreshTime = (this.expiresIn - 300) * 1000 // 5 minutes before expiry
      
      this.tokenRefreshInterval = setTimeout(async () => {
        await this.refreshAccessToken()
      }, refreshTime)
    },

    // Refresh access token using refresh token
    async refreshAccessToken() {
      if (!this.refreshToken) {
        this.logout('No refresh token available')
        return false
      }

      try {
        const response = await api.post('/refresh-token', {
          refresh_token: this.refreshToken
        })

        if (response.data && response.data.access_token) {
          this.accessToken = response.data.access_token
          this.refreshToken = response.data.refresh_token
          this.expiresIn = response.data.expires_in

          // Update localStorage
          localStorage.setItem('accessToken', response.data.access_token)
          localStorage.setItem('refreshToken', response.data.refresh_token)
          localStorage.setItem('expiresIn', response.data.expires_in)

          // Update API header
          api.defaults.headers.common['Authorization'] = `Bearer ${response.data.access_token}`

          // Restart refresh timer
          this.startTokenRefresh()

          console.log('Token refreshed successfully')
          return true
        }
      } catch (error) {
        console.error('Token refresh failed:', error)
        this.logout('Session expired')
        return false
      }
    },

    logout(reason = null) {
      // Clear intervals
      if (this.activityCheckInterval) {
        clearInterval(this.activityCheckInterval)
      }
      if (this.tokenRefreshInterval) {
        clearTimeout(this.tokenRefreshInterval)
      }

      // Remove event listeners
      const events = ['mousedown', 'keydown', 'scroll', 'touchstart', 'click']
      events.forEach(event => {
        window.removeEventListener(event, this.updateActivity)
      })

      // Clear state
      this.user = null
      this.accessToken = null
      this.refreshToken = null
      this.expiresIn = null
      this.isAuthenticated = false
      this.lastActivityTime = null

      // Clear localStorage
      localStorage.removeItem('accessToken')
      localStorage.removeItem('refreshToken')
      localStorage.removeItem('user')
      localStorage.removeItem('expiresIn')
      localStorage.removeItem('lastActivityTime')

      // Clear API header
      delete api.defaults.headers.common['Authorization']

      // Show reason if provided
      if (reason) {
        console.log('Logout reason:', reason)
        // You can show a toast/notification here
      }
    },

    // Restore auth from localStorage
    initAuth() {
      const accessToken = localStorage.getItem('accessToken')
      const refreshToken = localStorage.getItem('refreshToken')
      const user = localStorage.getItem('user')
      const expiresIn = localStorage.getItem('expiresIn')
      const lastActivityTime = localStorage.getItem('lastActivityTime')

      if (accessToken && refreshToken && user) {
        try {
          this.accessToken = accessToken
          this.refreshToken = refreshToken
          this.user = JSON.parse(user)
          this.expiresIn = parseInt(expiresIn)
          this.lastActivityTime = parseInt(lastActivityTime)
          this.isAuthenticated = true

          // Set token cho API requests
          api.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`

          // Check if session is still valid
          this.checkIdleTimeout()

          // Start tracking if still authenticated
          if (this.isAuthenticated) {
            this.startActivityTracking()
            this.startTokenRefresh()
          }
        } catch (error) {
          console.error('Error restoring auth:', error)
          this.logout()
        }
      }
    },

    // Validate token with backend
    async validateToken() {
      if (!this.accessToken) return false

      try {
        const response = await api.get('/profile')
        if (response.data) {
          this.user = response.data
          localStorage.setItem('user', JSON.stringify(response.data))
          return true
        }
        return false
      } catch (error) {
        console.error('Token validation failed:', error)
        
        // Try to refresh token if validation fails
        if (error.response?.status === 401) {
          return await this.refreshAccessToken()
        }
        
        this.logout('Session expired')
        return false
      }
    }
  }
})
```

```javascript
// services/api.js - Thêm interceptor để auto-refresh token
import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
api.interceptors.request.use(config => {
  const token = localStorage.getItem('accessToken')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor - Auto refresh token on 401
api.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config

    // If 401 and not already retried
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      const authStore = useAuthStore()
      const refreshed = await authStore.refreshAccessToken()

      if (refreshed) {
        // Retry original request with new token
        originalRequest.headers.Authorization = `Bearer ${authStore.accessToken}`
        return api(originalRequest)
      }
    }

    return Promise.reject(error)
  }
)

export default api
```

---

### Option 2: Session-Based với Redis (Advanced)

**Ưu điểm:**
- Có thể revoke session bất cứ lúc nào
- Track activity chính xác
- Có thể force logout user
- Có thể limit concurrent sessions

**Nhược điểm:**
- Cần thêm Redis server
- Stateful (phức tạp hơn khi scale)
- Tốn thêm infrastructure cost

**Cách hoạt động:**
1. User login → Tạo session ID, lưu vào Redis
2. Session có TTL (time to live) trong Redis
3. Mỗi request → Validate session ID với Redis
4. User activity → Reset TTL trong Redis
5. Session expire → Auto delete từ Redis

---

## 📊 So Sánh Các Option

| Feature | Current | Option 1 (JWT + Idle) | Option 2 (Redis) |
|---------|---------|----------------------|------------------|
| Stateless | ✅ | ✅ | ❌ |
| Absolute Timeout | ✅ (24h) | ✅ (1h) | ✅ (configurable) |
| Idle Timeout | ❌ | ✅ (30min) | ✅ (configurable) |
| Token Refresh | ❌ | ✅ | ✅ |
| Revoke Token | ❌ | ❌ | ✅ |
| Activity Tracking | ❌ | ✅ (frontend) | ✅ (backend) |
| Complexity | Low | Medium | High |
| Infrastructure | None | None | Redis required |
| Scalability | High | High | Medium |

---

## 🎯 Recommendation

**Đề xuất: Option 1 - JWT với Idle Timeout**

**Lý do:**
1. ✅ Không cần thêm infrastructure
2. ✅ Giữ được stateless architecture
3. ✅ Dễ implement và maintain
4. ✅ Đáp ứng đủ yêu cầu security
5. ✅ Phù hợp với hệ thống hiện tại

**Cấu hình đề xuất:**
- Access token: 1 giờ
- Refresh token: 7 ngày
- Idle timeout: 30 phút
- Auto refresh: 5 phút trước khi expire

---

## 🚀 Implementation Plan

### Phase 1: Backend
1. Update JWT service để generate token pair
2. Thêm refresh token endpoint
3. Update login response structure
4. Test với Postman

### Phase 2: Frontend
1. Update auth store với idle tracking
2. Implement token refresh logic
3. Add API interceptor
4. Test idle timeout

### Phase 3: Testing
1. Test absolute timeout
2. Test idle timeout
3. Test token refresh
4. Test concurrent sessions

### Phase 4: Deployment
1. Update API documentation
2. Update user guide
3. Deploy backend
4. Deploy frontend
5. Monitor logs

---

## 📝 Configuration

```javascript
// frontend/src/config/session.js
export const SESSION_CONFIG = {
  // Access token lifetime (backend)
  ACCESS_TOKEN_LIFETIME: 3600, // 1 hour in seconds
  
  // Refresh token lifetime (backend)
  REFRESH_TOKEN_LIFETIME: 604800, // 7 days in seconds
  
  // Idle timeout (frontend)
  IDLE_TIMEOUT_MINUTES: 30,
  
  // Token refresh before expiry (frontend)
  REFRESH_BEFORE_EXPIRY_SECONDS: 300, // 5 minutes
  
  // Activity check interval (frontend)
  ACTIVITY_CHECK_INTERVAL_MS: 60000, // 1 minute
  
  // Events to track for activity
  ACTIVITY_EVENTS: ['mousedown', 'keydown', 'scroll', 'touchstart', 'click']
}
```

---

## 🔒 Security Considerations

1. **Token Storage:**
   - ✅ localStorage (current) - OK for most cases
   - ⚠️ Consider httpOnly cookies for higher security
   - ❌ Avoid sessionStorage (lost on tab close)

2. **XSS Protection:**
   - Sanitize all user inputs
   - Use Content Security Policy (CSP)
   - Validate token on every request

3. **CSRF Protection:**
   - Use CSRF tokens for state-changing operations
   - Validate origin header

4. **Token Rotation:**
   - Rotate refresh token on each use
   - Invalidate old refresh tokens

---

## 📚 References

- JWT Best Practices: https://tools.ietf.org/html/rfc8725
- OWASP Session Management: https://owasp.org/www-community/Session_Management_Cheat_Sheet
- Vue.js Authentication: https://vuejs.org/guide/best-practices/security.html

---

**Last Updated:** 2026-02-04  
**Status:** Ready for Implementation  
**Estimated Time:** 4-6 hours
