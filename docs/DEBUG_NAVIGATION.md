# Debug Navigation Issue

## Vấn đề
Manager (admin/admin123) không thấy menu Facility và Ingredient sau khi login.

## Các bước debug

### 1. Kiểm tra User Role trong Console

Mở browser console (F12) và chạy:

```javascript
// Check localStorage
console.log('Token:', localStorage.getItem('token'))
console.log('User:', localStorage.getItem('user'))

// Parse user
const user = JSON.parse(localStorage.getItem('user'))
console.log('User role:', user?.role)
console.log('User name:', user?.name)
```

**Expected output**:
```
User role: "manager"
User name: "Administrator"
```

### 2. Kiểm tra Auth Store

Trong console:

```javascript
// Import store (nếu có thể)
// Hoặc check Vue DevTools

// Check if Navigation component receives correct role
```

### 3. Kiểm tra Navigation Component

Thêm debug log vào Navigation.vue:

```vue
<script setup>
import { computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const userRole = computed(() => {
  console.log('🔍 Navigation - User:', authStore.user)
  console.log('🔍 Navigation - Role:', authStore.user?.role)
  return authStore.user?.role
})

const userName = computed(() => authStore.user?.name)

// Watch for changes
watch(userRole, (newRole) => {
  console.log('🔄 Role changed to:', newRole)
})

const handleNavClick = () => {
  // Optional: Add any navigation handling logic here
}

const logout = () => {
  authStore.logout()
  router.push('/login')
}
</script>
```

### 4. Force Refresh Frontend

```bash
# Stop frontend (Ctrl+C)
cd frontend

# Clear cache
rm -rf node_modules/.vite
rm -rf dist

# Restart
npm run dev
```

### 5. Hard Refresh Browser

```
1. Mở DevTools (F12)
2. Right-click vào nút Refresh
3. Chọn "Empty Cache and Hard Reload"
```

### 6. Kiểm tra Backend Response

Trong Network tab (F12):

```
1. Login với admin/admin123
2. Xem request POST /api/login
3. Check response:
   {
     "user": {
       "id": "...",
       "username": "admin",
       "role": "manager",  // ← Phải là "manager"
       "name": "Administrator"
     },
     "token": "..."
   }
```

## Quick Fix Options

### Option 1: Add Debug Info to Navigation

Thêm vào Navigation.vue (temporary):

```vue
<template>
  <nav class="bg-white shadow-lg">
    <!-- Debug info (remove after fixing) -->
    <div class="bg-yellow-100 p-2 text-xs">
      Debug: Role = {{ userRole }} | User = {{ userName }}
    </div>
    
    <!-- Rest of navigation -->
    ...
  </nav>
</template>
```

### Option 2: Force Show Menus (Testing Only)

Temporarily change condition:

```vue
<!-- Change from -->
<template v-if="userRole === 'manager'">

<!-- To (for testing) -->
<template v-if="true">
```

**⚠️ Remember to change back after testing!**

### Option 3: Check Main.js

Verify auth is initialized in `main.js`:

```javascript
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './style.css'
import { useAuthStore } from './stores/auth'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// Initialize auth from localStorage
const authStore = useAuthStore()
authStore.initAuth()

app.mount('#app')
```

## Expected Behavior

Khi login với admin/admin123:

1. ✅ Backend trả về user với role="manager"
2. ✅ Auth store lưu user vào state
3. ✅ Navigation component nhận được userRole="manager"
4. ✅ Template v-if="userRole === 'manager'" evaluates to true
5. ✅ Manager menus hiển thị:
   - 👥 Quản lý User
   - 🍽️ Menu
   - 🥬 Nguyên liệu
   - 🏢 Cơ sở vật chất
   - 💰 Chi phí

## Common Issues

### Issue 1: Role không đúng format
**Problem**: Backend trả về "Manager" thay vì "manager"
**Solution**: Check backend user model

### Issue 2: Frontend cache
**Problem**: Old code still running
**Solution**: Hard refresh + clear cache

### Issue 3: Auth not initialized
**Problem**: Auth store không load từ localStorage
**Solution**: Check main.js có gọi `authStore.initAuth()`

### Issue 4: Reactive issue
**Problem**: Computed không update
**Solution**: Restart dev server

## Test Script

Tạo file `test-navigation.html`:

```html
<!DOCTYPE html>
<html>
<head>
  <title>Test Navigation</title>
</head>
<body>
  <h1>Navigation Test</h1>
  <div id="result"></div>
  
  <script>
    // Get from localStorage
    const user = JSON.parse(localStorage.getItem('user'))
    const token = localStorage.getItem('token')
    
    const result = document.getElementById('result')
    result.innerHTML = `
      <p>Token: ${token ? 'Present' : 'Missing'}</p>
      <p>User: ${user ? JSON.stringify(user, null, 2) : 'Missing'}</p>
      <p>Role: ${user?.role}</p>
      <p>Should show manager menus: ${user?.role === 'manager' ? 'YES' : 'NO'}</p>
    `
  </script>
</body>
</html>
```

## Solution Steps

1. **Clear everything**:
   ```bash
   # Clear browser cache
   # Clear localStorage
   localStorage.clear()
   
   # Restart frontend
   cd frontend
   npm run dev
   ```

2. **Login again**:
   - Username: admin
   - Password: admin123

3. **Check console**:
   - Should see role="manager"

4. **Verify menus appear**:
   - Should see 5 manager-only cards

## If Still Not Working

Add this to Navigation.vue temporarily:

```vue
<script setup>
// ... existing code ...

// Force log on mount
import { onMounted } from 'vue'

onMounted(() => {
  console.log('=== Navigation Mounted ===')
  console.log('Auth Store:', authStore)
  console.log('User:', authStore.user)
  console.log('Role:', authStore.user?.role)
  console.log('Is Manager:', authStore.user?.role === 'manager')
  console.log('========================')
})
</script>
```

Then check console output when page loads.

