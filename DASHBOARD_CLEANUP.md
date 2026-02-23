# Dashboard Cleanup - Xóa Section "Orders gần đây"

## VẤN ĐỀ

Dashboard (http://localhost:5173/#/dashboard) có section "Orders gần đây" hiển thị 3 orders mới nhất.

**Lý do xóa:**
- Thông tin này đã có đầy đủ trong trang Orders (http://localhost:5173/#/orders)
- Duplicate content không cần thiết
- Dashboard nên tập trung vào thông tin tổng quan (stats, shift info)

## GIẢI PHÁP

Xóa toàn bộ section "Orders gần đây" khỏi Dashboard.

### Code đã xóa

#### 1. Template - Section UI
```vue
<!-- Recent Orders -->
<div class="mb-4">
  <div class="flex items-center justify-between mb-3">
    <h2 class="text-lg font-bold text-gray-800">🕐 Orders gần đây</h2>
    <button @click="$router.push('/orders')" class="text-sm text-blue-500 font-medium">
      Xem tất cả →
    </button>
  </div>
  <div v-if="recentOrders.length === 0" class="text-center py-8 text-gray-500">
    <div class="text-4xl mb-2">📭</div>
    <p>Chưa có order nào</p>
  </div>
  <div v-else class="space-y-3">
    <div v-for="order in recentOrders.slice(0, 3)" :key="order.id"
      @click="$router.push('/orders')"
      class="bg-white rounded-xl p-4 shadow-sm active:scale-98 transition-transform">
      <!-- Order details -->
    </div>
  </div>
</div>
```

#### 2. Script - Computed property
```js
const recentOrders = computed(() => {
  return [...orders.value].sort((a, b) => 
    new Date(b.created_at) - new Date(a.created_at)
  )
})
```

## DASHBOARD SAU KHI CLEANUP

Dashboard bây giờ chỉ hiển thị:

### Cho Waiter/Barista:
1. **Header** - Chào mừng + thời gian
2. **Shift Info** - Thông tin ca làm việc hiện tại
3. **Quick Stats** - Thống kê nhanh (nếu có)
4. **Quick Actions** - Các nút action nhanh

### Cho Cashier:
1. **Header** - Chào mừng + thời gian
2. **Cashier Stats** - Thống kê thu ngân
3. **Pending Handovers** - Bàn giao chờ xác nhận

## LỢI ÍCH

✅ Dashboard gọn gàng hơn, tập trung vào thông tin quan trọng
✅ Không duplicate content với trang Orders
✅ Giảm số lượng API calls không cần thiết
✅ User muốn xem orders → vào trang Orders (có đầy đủ filter, search, etc.)
✅ Dashboard load nhanh hơn

## NAVIGATION

User muốn xem orders:
1. Click vào tab "Orders" ở bottom navigation
2. Hoặc click vào quick action "Tạo order mới" (nếu có)

## FILES THAY ĐỔI

- `frontend/src/views/DashboardView.vue`
  - Xóa section "Recent Orders" (template)
  - Xóa computed `recentOrders` (script)
  - Giữ nguyên các phần khác (shift info, stats, etc.)

## LƯU Ý

Nếu sau này muốn thêm lại thông tin orders trong dashboard, nên:
- Chỉ hiển thị số lượng (count) theo status
- Không hiển thị list chi tiết
- Ví dụ: "5 orders đang chờ", "3 orders đang làm"
