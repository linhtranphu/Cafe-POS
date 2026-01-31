# Fix Navigation - Quick Steps

## ✅ Đã Thêm Debug Info

Tôi đã thêm debug banner vào Navigation component. Bây giờ làm theo các bước sau:

## Bước 1: Restart Frontend

```bash
# Stop frontend (Ctrl+C nếu đang chạy)
cd frontend

# Start lại
npm run dev
```

## Bước 2: Clear Browser Cache

```
1. Mở browser
2. Nhấn F12 (mở DevTools)
3. Right-click vào nút Refresh
4. Chọn "Empty Cache and Hard Reload"
```

## Bước 3: Clear LocalStorage

Trong Console (F12), chạy:

```javascript
localStorage.clear()
location.reload()
```

## Bước 4: Login Lại

```
Username: admin
Password: admin123
```

## Bước 5: Xem Debug Banner

Sau khi login, bạn sẽ thấy một banner màu vàng ở trên cùng:

```
🐛 Debug: Role = manager | User = Administrator | Is Manager = YES ✅
```

### Nếu thấy:
- ✅ **Role = manager** và **Is Manager = YES** → Menu sẽ hiển thị
- ❌ **Role = undefined** hoặc khác → Có vấn đề với login

## Bước 6: Check Console

Mở Console (F12) và xem logs:

```
=== Navigation Component Mounted ===
Auth Store User: {username: "admin", role: "manager", ...}
User Role: manager
Is Manager: true
====================================
```

## Nếu Vẫn Không Thấy Menu

### Option A: Check Backend Response

1. Mở Network tab (F12)
2. Login lại
3. Tìm request `POST /api/login`
4. Xem Response:

```json
{
  "user": {
    "username": "admin",
    "role": "manager",  // ← Phải là "manager"
    "name": "Administrator"
  },
  "token": "..."
}
```

### Option B: Force Show (Testing)

Tạm thời sửa Navigation.vue:

```vue
<!-- Tìm dòng này -->
<template v-if="userRole === 'manager'">

<!-- Đổi thành (CHỈ ĐỂ TEST) -->
<template v-if="true">
```

Nếu menu hiển thị → Vấn đề là userRole không đúng
Nếu vẫn không hiển thị → Vấn đề khác

### Option C: Check Main.js

File `frontend/src/main.js` phải có:

```javascript
import { useAuthStore } from './stores/auth'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// ← Phải có dòng này
const authStore = useAuthStore()
authStore.initAuth()

app.mount('#app')
```

## Expected Result

Sau khi login với admin/admin123, bạn sẽ thấy:

1. ✅ Debug banner: `Role = manager | Is Manager = YES ✅`
2. ✅ 5 manager cards:
   - 👥 Quản lý User
   - 🍽️ Menu
   - 🥬 Nguyên liệu ← **MỚI**
   - 🏢 Cơ sở vật chất ← **MỚI**
   - 💰 Chi phí

## Remove Debug Banner (Sau khi fix)

Khi đã hoạt động, xóa debug banner trong Navigation.vue:

```vue
<!-- Xóa phần này -->
<div class="bg-yellow-100 border-b border-yellow-300 px-4 py-2 text-xs">
  <strong>🐛 Debug:</strong> ...
</div>
```

Và xóa console.log trong script:

```javascript
// Xóa các dòng console.log
console.log('🔍 Navigation - User:', authStore.user)
console.log('🔍 Navigation - Role:', role)
```

## Liên Hệ

Nếu vẫn không hoạt động:
1. Chụp màn hình debug banner
2. Copy console logs
3. Copy Network response của /api/login
4. Gửi cho tôi để debug tiếp

