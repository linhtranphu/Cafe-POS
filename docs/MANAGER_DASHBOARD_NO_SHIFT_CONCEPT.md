# Manager Dashboard - Loại bỏ khái niệm ca làm

## Tóm tắt
Đã cập nhật DashboardView và Navigation để loại bỏ hoàn toàn khái niệm ca làm (shift) cho role Manager.

## 1. Thay đổi trong DashboardView.vue

### Template - Quick Stats cho Manager
**Trước:**
- 📋 Orders hôm nay
- 💰 Doanh thu hôm nay
- 👥 **Nhân viên đang làm** (openShiftsCount) ❌
- 🍹 Đang xử lý

**Sau:**
- 📋 Orders hôm nay
- 💰 Doanh thu hôm nay
- ✅ **Hoàn tất hôm nay** (completedOrders) ✅
- 🍹 **Đang xử lý** (pendingOrders - tất cả orders chưa hoàn tất) ✅

### Script - Computed Properties

#### Thêm mới:
```javascript
const completedOrders = computed(() => {
  const today = new Date().toDateString()
  return orders.value.filter(o => 
    new Date(o.created_at).toDateString() === today && o.status === 'SERVED'
  ).length
})
```

#### Cập nhật:
```javascript
const pendingOrders = computed(() => {
  // For manager: show all orders that are not completed or cancelled
  if (user.value?.role === 'manager') {
    return orders.value.filter(o => 
      o.status !== 'SERVED' && o.status !== 'CANCELLED'
    ).length
  }
  // For others: show only created orders
  return orders.value.filter(o => o.status === 'CREATED').length
})
```

### Script - Data Loading

**Trước:**
```javascript
if (user.value?.role === 'manager') {
  await Promise.all([
    orderStore.fetchOrders(),
    shiftStore.fetchAllShifts() // ❌ Vẫn fetch shift data
  ])
  return
}
```

**Sau:**
```javascript
if (user.value?.role === 'manager') {
  await orderStore.fetchOrders() // ✅ Chỉ fetch orders
  return
}
```

## 2. Thay đổi trong Navigation.vue

### Manager Navigation
**Trước (5 menu):**
1. 🏠 Dashboard
2. ⏰ Quản lý ca ❌
3. 👥 Nhân viên
4. 📊 Báo cáo
5. 👤 Cá nhân

**Sau (4 menu):**
1. 🏠 Dashboard
2. 📊 Báo cáo
3. 👥 Nhân viên
4. 👤 Cá nhân

### Layout
- Grid: `grid-cols-2 sm:grid-cols-4` (thay vì `grid-cols-2 sm:grid-cols-3 lg:grid-cols-5`)
- Max width: `max-w-4xl` (thay vì `max-w-6xl`)

## 3. Thay đổi trong BottomNav.vue

### Manager Bottom Navigation
**Trước:**
- 🏠 Trang chủ
- 💰 Thu ngân
- 📋 Orders
- ⏰ Ca làm ❌
- 👤 Cá nhân

**Sau:**
- 🏠 Dashboard
- 📊 Báo cáo
- 👥 Nhân viên
- 👤 Cá nhân

### Logic
```javascript
// Manager navigation (4 items)
if (role === 'manager') {
  return [
    { path: '/dashboard', icon: '🏠', label: 'Dashboard' },
    { path: '/cashier/reports', icon: '📊', label: 'Báo cáo' },
    { path: '/users', icon: '👥', label: 'Nhân viên' },
    { path: '/profile', icon: '👤', label: 'Cá nhân' }
  ]
}
```

## Phân biệt rõ ràng

### Manager (Không có ca làm)
- ✅ Không cần mở/đóng ca
- ✅ Xem tất cả orders trong hệ thống
- ✅ Thống kê theo ngày (không theo ca)
- ✅ Không có menu "Ca làm" hoặc "Quản lý ca"
- ✅ Navigation đơn giản: Dashboard, Báo cáo, Nhân viên, Cá nhân

### Nhân viên (Có ca làm)
- ✅ Phải mở ca trước khi làm việc
- ✅ Xem orders trong ca của mình
- ✅ Thống kê theo ca làm
- ✅ Menu "Ca làm việc" để mở/đóng ca

## Kết quả
Manager giờ đây có:
- Dashboard riêng biệt, tập trung vào quản lý tổng thể
- Navigation gọn gàng với 4 menu chính
- Không bị ràng buộc bởi khái niệm ca làm
- Trải nghiệm người dùng tối ưu cho vai trò quản lý
