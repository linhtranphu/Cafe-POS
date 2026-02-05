# Pull-to-Refresh Quick Implementation Guide

## 🚀 Quick Steps (5 phút/view)

### Step 1: Import (thêm vào script section)

```javascript
import PullToRefresh from '../components/PullToRefresh.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'
```

### Step 2: Create refresh function

```javascript
// Tìm onMounted và copy logic fetch data
const refreshData = async () => {
  // Copy tất cả fetch calls từ onMounted
  await store.fetchData()
  await store.fetchMoreData()
}
```

### Step 3: Use composable

```javascript
// Thêm sau khi define refreshData
const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)
```

### Step 4: Update onMounted

```javascript
onMounted(async () => {
  await refreshData() // Gọi refreshData thay vì duplicate code
})
```

### Step 5: Add component to template

```vue
<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Pull to Refresh Indicator - ADD THIS -->
    <PullToRefresh 
      :pull-distance="pullDistance" 
      :is-refreshing="isRefreshing"
      :threshold="80" />
    
    <!-- Rest of template -->
    ...
  </div>
</template>
```

## 📋 Complete Example

### Before:
```vue
<template>
  <div class="min-h-screen bg-gray-50">
    <div class="sticky top-0">
      <h1>My View</h1>
    </div>
    <!-- content -->
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useMyStore } from '../stores/myStore'

const myStore = useMyStore()

onMounted(async () => {
  await myStore.fetchData()
})
</script>
```

### After:
```vue
<template>
  <div class="min-h-screen bg-gray-50">
    <!-- ADD THIS -->
    <PullToRefresh 
      :pull-distance="pullDistance" 
      :is-refreshing="isRefreshing"
      :threshold="80" />
    
    <div class="sticky top-0">
      <h1>My View</h1>
    </div>
    <!-- content -->
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useMyStore } from '../stores/myStore'
// ADD THESE IMPORTS
import PullToRefresh from '../components/PullToRefresh.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'

const myStore = useMyStore()

// ADD THIS FUNCTION
const refreshData = async () => {
  await myStore.fetchData()
}

// ADD THIS COMPOSABLE
const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)

// UPDATE onMounted
onMounted(async () => {
  await refreshData() // Call refreshData instead
})
</script>
```

## 📝 Views Checklist

### Priority 1 (Data changes frequently):
- [x] CashierHandoverView.vue ✅ DONE
- [ ] OrderView.vue
- [ ] CashierDashboard.vue
- [ ] ShiftView.vue
- [ ] BaristaView.vue

### Priority 2 (Moderate frequency):
- [ ] DashboardView.vue
- [ ] ManagerShiftView.vue
- [ ] ExpenseManagementView.vue
- [ ] IngredientManagementView.vue
- [ ] FacilityManagementView.vue
- [ ] CashierReports.vue

### Skip (Static or form-based):
- ProfileView.vue (user profile - rarely changes)
- LoginView.vue (no data to refresh)
- MenuView.vue (menu items - rarely change)
- UserManagementView.vue (admin only, low frequency)
- FacilityAddEditView.vue (form view)
- CashierShiftClosure.vue (one-time action)

## 🎯 View-Specific Examples

### OrderView.vue
```javascript
const refreshData = async () => {
  await Promise.all([
    shiftStore.fetchCurrentShift(),
    orderStore.fetchOrders(),
    menuStore.fetchMenuItems()
  ])
}
```

### CashierDashboard.vue
```javascript
const refreshData = async () => {
  if (selectedShift.value) {
    await Promise.all([
      cashierStore.getShiftStatus(selectedShift.value),
      cashierStore.getPaymentsByShift(selectedShift.value)
    ])
  }
  await cashierStore.getPendingDiscrepancies()
}
```

### ShiftView.vue
```javascript
const refreshData = async () => {
  await shiftStore.fetchCurrentShift()
  if (isWaiter.value && currentShift.value) {
    await Promise.all([
      fetchPendingHandover(),
      fetchHandoverHistory()
    ])
  }
  await shiftStore.fetchMyShifts()
}
```

### BaristaView.vue
```javascript
const refreshData = async () => {
  await Promise.all([
    baristaStore.fetchOrders(),
    shiftStore.fetchCurrentShift()
  ])
}
```

### DashboardView.vue
```javascript
const refreshData = async () => {
  await Promise.all([
    shiftStore.fetchCurrentShift(),
    orderStore.fetchOrders(),
    baristaStore.fetchOrders()
  ])
}
```

### ManagerShiftView.vue
```javascript
const refreshData = async () => {
  await Promise.all([
    shiftStore.fetchAllShifts(),
    cashierShiftStore.fetchAllCashierShifts()
  ])
}
```

### ExpenseManagementView.vue
```javascript
const refreshData = async () => {
  await Promise.all([
    expenseStore.fetchCategories(),
    expenseStore.fetchExpenses()
  ])
}
```

### IngredientManagementView.vue
```javascript
const refreshData = async () => {
  await Promise.all([
    ingredientStore.fetchCategories(),
    ingredientStore.fetchIngredients()
  ])
}
```

### FacilityManagementView.vue
```javascript
const refreshData = async () => {
  await Promise.all([
    facilityStore.fetchFacilities(),
    facilityStore.fetchFacilityTypes()
  ])
}
```

## ⚡ Pro Tips

### 1. Use Promise.all for parallel fetching
```javascript
// ✅ Good - Fast
const refreshData = async () => {
  await Promise.all([
    store1.fetch(),
    store2.fetch()
  ])
}

// ❌ Bad - Slow
const refreshData = async () => {
  await store1.fetch()
  await store2.fetch()
}
```

### 2. Handle conditional fetching
```javascript
const refreshData = async () => {
  await store.fetchData()
  
  // Only fetch if condition met
  if (someCondition.value) {
    await store.fetchMore()
  }
}
```

### 3. Reuse existing refresh functions
```javascript
// If view already has refreshOrders(), refreshData(), etc.
// Just use it directly!
const { pullDistance, isRefreshing } = usePullToRefresh(refreshOrders)
```

### 4. Add error handling
```javascript
const refreshData = async () => {
  try {
    await store.fetchData()
  } catch (error) {
    console.error('Refresh failed:', error)
    // Optionally show toast/notification
  }
}
```

## 🐛 Common Issues

### Issue: "pullDistance is not defined"
**Solution:** Make sure you destructure from usePullToRefresh:
```javascript
const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)
```

### Issue: Pull-to-refresh not working
**Solution:** Check if PullToRefresh component is INSIDE the main div:
```vue
<template>
  <div class="min-h-screen">
    <PullToRefresh ... /> <!-- Must be here -->
    ...
  </div>
</template>
```

### Issue: Refresh called twice
**Solution:** Make sure onMounted calls refreshData, not duplicate fetch:
```javascript
// ✅ Good
onMounted(async () => {
  await refreshData()
})

// ❌ Bad - will fetch twice
onMounted(async () => {
  await refreshData()
  await store.fetchData() // Duplicate!
})
```

## 📚 Full Documentation

See `docs/PULL_TO_REFRESH_IMPLEMENTATION.md` for complete details.

---

**Time estimate:** 5 minutes per view
**Difficulty:** Easy
**Status:** Template ready, implement as needed
