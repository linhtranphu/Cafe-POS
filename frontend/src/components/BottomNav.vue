<template>
  <div class="fixed bottom-0 left-0 right-0 bg-white border-t shadow-lg z-40">
    <div class="flex justify-around py-2">
      <button 
        v-for="item in navItems" 
        :key="item.path"
        @click="navigate(item.path)"
        :class="[
          'flex flex-col items-center py-2 px-3 rounded-lg transition-colors relative',
          isActive(item.path) ? 'text-blue-500' : 'text-gray-600'
        ]">
        <span class="text-2xl mb-1">{{ item.icon }}</span>
        <span class="text-xs font-medium whitespace-nowrap">{{ item.label }}</span>
        <!-- Badge for pending items -->
        <span v-if="item.badge && item.badge > 0" 
          class="absolute top-1 right-1 bg-red-500 text-white text-xs font-bold rounded-full h-5 w-5 flex items-center justify-center">
          {{ item.badge > 9 ? '9+' : item.badge }}
        </span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useCashierStore } from '../stores/cashier'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const cashierStore = useCashierStore()

const navItems = computed(() => {
  const role = authStore.user?.role
  
  // Manager navigation (5 items - essential only, rest in dashboard)
  if (role === 'manager') {
    return [
      { path: '/dashboard', icon: '🏠', label: 'Dashboard' },
      { path: '/cashier/reports', icon: '📊', label: 'Báo cáo' },
      { path: '/manager/inventory-stats', icon: '📦', label: 'Tồn kho' },
      { path: '/manager/profit-analysis', icon: '📈', label: 'Lợi nhuận' },
      { path: '/profile', icon: '👤', label: 'Cá nhân' }
    ]
  }

  // Barista navigation
  if (role === 'barista') {
    return [
      { path: '/dashboard', icon: '🏠', label: 'Trang chủ' },
      { path: '/barista', icon: '🍹', label: 'Barista' },
      { path: '/shifts', icon: '⏰', label: 'Ca làm' },
      { path: '/profile', icon: '👤', label: 'Cá nhân' }
    ]
  }

  // Cashier navigation with badge
  if (role === 'cashier') {
    const pendingCount = cashierStore.pendingHandovers?.length || 0
    return [
      { path: '/dashboard', icon: '🏠', label: 'Trang chủ' },
      { path: '/cashier', icon: '💰', label: 'Thu ngân', badge: pendingCount },
      { path: '/cashier/reports', icon: '📊', label: 'Báo cáo' },
      { path: '/shifts', icon: '⏰', label: 'Ca làm' },
      { path: '/profile', icon: '👤', label: 'Cá nhân' }
    ]
  }

  // Default navigation (waiter, etc.)
  return [
    { path: '/dashboard', icon: '🏠', label: 'Trang chủ' },
    { path: '/create-order', icon: '➕', label: 'Tạo order' },
    { path: '/orders', icon: '📋', label: 'Orders' },
    { path: '/shifts', icon: '⏰', label: 'Ca làm' },
    { path: '/profile', icon: '👤', label: 'Cá nhân' }
  ]
})

const isActive = (path) => {
  return route.path === path || route.path.startsWith(path + '/')
}

const navigate = (path) => {
  router.push(path)
}
</script>

<style scoped>
/* No safe area padding - container scroll handles spacing with pb-24 */
</style>
