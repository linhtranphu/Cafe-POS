# Pull-to-Refresh Implementation Guide

## 📱 Tổng quan

Pull-to-refresh là pattern cho phép người dùng kéo màn hình xuống để làm mới dữ liệu. Đây là UX pattern phổ biến trên mobile apps.

## 🎯 Đã implement

### Files created:
1. **`frontend/src/composables/usePullToRefresh.js`** - Composable xử lý logic
2. **`frontend/src/components/PullToRefresh.vue`** - UI component hiển thị indicator
3. **Example:** CashierHandoverView đã được thêm pull-to-refresh

## 🚀 Cách sử dụng

### Bước 1: Import composable và component

```vue
<script setup>
import { usePullToRefresh } from '../composables/usePullToRefresh'
import PullToRefresh from '../components/PullToRefresh.vue'
</script>
```

### Bước 2: Tạo refresh function

```javascript
const refreshData = async () => {
  // Fetch data của bạn ở đây
  await store.fetchData()
  await store.fetchMoreData()
}
```

### Bước 3: Sử dụng composable

```javascript
const { isPulling, isRefreshing, pullDistance } = usePullToRefresh(refreshData, {
  threshold: 80,        // Khoảng cách kéo để trigger refresh (px)
  resistance: 2.5,      // Độ kháng khi kéo (càng cao càng khó kéo)
  maxPullDistance: 150  // Khoảng cách kéo tối đa (px)
})
```

### Bước 4: Thêm component vào template

```vue
<template>
  <div>
    <!-- Pull to Refresh Indicator -->
    <PullToRefresh 
      :pull-distance="pullDistance" 
      :is-refreshing="isRefreshing"
      :threshold="80" />
    
    <!-- Your content here -->
    <div class="content">
      ...
    </div>
  </div>
</template>
```

## 📋 Full Example

```vue
<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Pull to Refresh Indicator -->
    <PullToRefresh 
      :pull-distance="pullDistance" 
      :is-refreshing="isRefreshing"
      :threshold="80" />
    
    <!-- Header -->
    <div class="sticky top-0 z-40 bg-white shadow-sm">
      <h1 class="text-xl font-bold">My View</h1>
    </div>

    <!-- Content -->
    <div class="px-4 py-4">
      <div v-for="item in items" :key="item.id">
        {{ item.name }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useMyStore } from '../stores/myStore'
import PullToRefresh from '../components/PullToRefresh.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'

const myStore = useMyStore()
const items = computed(() => myStore.items)

// Refresh function
const refreshData = async () => {
  await myStore.fetchItems()
}

// Pull to refresh
const { isPulling, isRefreshing, pullDistance } = usePullToRefresh(refreshData)

onMounted(async () => {
  await refreshData()
})
</script>
```

## ⚙️ Configuration Options

### usePullToRefresh options:

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `threshold` | Number | 80 | Khoảng cách kéo (px) để trigger refresh |
| `resistance` | Number | 2.5 | Độ kháng khi kéo (càng cao càng khó) |
| `maxPullDistance` | Number | 150 | Khoảng cách kéo tối đa (px) |

### PullToRefresh props:

| Prop | Type | Required | Description |
|------|------|----------|-------------|
| `pullDistance` | Number | Yes | Khoảng cách đã kéo (từ composable) |
| `isRefreshing` | Boolean | Yes | Trạng thái đang refresh (từ composable) |
| `threshold` | Number | No | Ngưỡng để hiển thị "Thả để làm mới" |

## 🎨 UI States

Pull-to-refresh component có 3 trạng thái:

1. **Pulling (⬇️)** - Đang kéo nhưng chưa đủ threshold
   - Text: "Kéo xuống để làm mới"
   
2. **Ready (🎯)** - Đã kéo đủ threshold, sẵn sàng refresh
   - Text: "Thả để làm mới"
   
3. **Refreshing (🔄)** - Đang refresh data
   - Text: "Đang tải..."
   - Icon quay tròn

## 🔧 How it works

### Touch Events Flow:

```
1. touchstart
   ↓
   Check if at top of page (scrollTop === 0)
   ↓
   Record startY position

2. touchmove
   ↓
   Calculate distance = currentY - startY
   ↓
   Apply resistance: pullDistance = distance / resistance
   ↓
   Update UI indicator

3. touchend
   ↓
   Check if pullDistance >= threshold
   ↓
   If yes: trigger refresh
   If no: reset pullDistance
```

### Resistance Calculation:

```javascript
// Without resistance: Pull 100px → Indicator moves 100px
// With resistance 2.5: Pull 100px → Indicator moves 40px

pullDistance = actualPullDistance / resistance
```

Điều này tạo cảm giác tự nhiên, không quá nhạy.

## 📱 Mobile-First Design

Pull-to-refresh được thiết kế cho mobile:
- Chỉ hoạt động với touch events
- Chỉ trigger khi ở đầu trang (scrollTop === 0)
- Có resistance để tránh trigger nhầm
- Smooth animation

## 🎯 Best Practices

### 1. Chỉ dùng cho views có data cần refresh thường xuyên
```javascript
// ✅ Good: List views, dashboards
- CashierHandoverView (pending handovers)
- OrderView (orders list)
- DashboardView (stats)

// ❌ Bad: Static pages
- ProfileView
- LoginView
```

### 2. Refresh function nên nhanh
```javascript
// ✅ Good: Parallel fetching
const refreshData = async () => {
  await Promise.all([
    store.fetchPendingHandovers(),
    store.fetchTodayHandovers()
  ])
}

// ❌ Bad: Sequential fetching
const refreshData = async () => {
  await store.fetchPendingHandovers()
  await store.fetchTodayHandovers()
}
```

### 3. Handle errors gracefully
```javascript
const refreshData = async () => {
  try {
    await store.fetchData()
  } catch (error) {
    console.error('Refresh failed:', error)
    // Show error toast/notification
  }
}
```

### 4. Đặt threshold hợp lý
```javascript
// ✅ Good: 60-100px (comfortable pull distance)
const { ... } = usePullToRefresh(refreshData, { threshold: 80 })

// ❌ Bad: Too small (trigger too easily)
const { ... } = usePullToRefresh(refreshData, { threshold: 20 })

// ❌ Bad: Too large (hard to trigger)
const { ... } = usePullToRefresh(refreshData, { threshold: 200 })
```

## 🐛 Troubleshooting

### Issue: Pull-to-refresh không hoạt động

**Nguyên nhân:**
- Không ở đầu trang (scrollTop > 0)
- Có element khác đang handle touch events
- Browser không hỗ trợ touch events

**Giải pháp:**
```javascript
// Check if at top
console.log('scrollTop:', window.pageYOffset)

// Check touch events
document.addEventListener('touchstart', (e) => {
  console.log('Touch detected:', e.touches[0])
})
```

### Issue: Trigger quá nhạy

**Giải pháp:** Tăng resistance hoặc threshold
```javascript
const { ... } = usePullToRefresh(refreshData, {
  threshold: 100,    // Increase from 80
  resistance: 3.5    // Increase from 2.5
})
```

### Issue: Animation không smooth

**Giải pháp:** Thêm CSS transition
```css
.pull-indicator {
  transition: transform 0.2s ease-out;
}
```

## 🔄 Alternative: Button Refresh

Nếu không muốn dùng pull-to-refresh, có thể dùng button:

```vue
<template>
  <div class="sticky top-0 z-40 bg-white shadow-sm">
    <div class="px-4 py-3 flex justify-between items-center">
      <h1 class="text-xl font-bold">My View</h1>
      <button @click="refreshData" :disabled="loading"
        class="p-2 rounded-lg bg-gray-100 hover:bg-gray-200">
        <span :class="{ 'animate-spin': loading }">🔄</span>
      </button>
    </div>
  </div>
</template>

<script setup>
const loading = ref(false)

const refreshData = async () => {
  loading.value = true
  try {
    await store.fetchData()
  } finally {
    loading.value = false
  }
}
</script>
```

## 📚 References

- [MDN Touch Events](https://developer.mozilla.org/en-US/docs/Web/API/Touch_events)
- [Vue 3 Composables](https://vuejs.org/guide/reusability/composables.html)
- [Mobile UX Patterns](https://www.nngroup.com/articles/pull-to-refresh/)

---

**Implemented in:** CashierHandoverView
**Status:** ✅ Ready to use
**Mobile-friendly:** Yes
**Browser support:** Modern browsers with touch events
