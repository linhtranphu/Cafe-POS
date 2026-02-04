<template>
  <div class="fixed bottom-0 left-0 right-0 bg-white border-t shadow-lg z-40 safe-area-bottom">
    <div class="flex justify-around py-2">
      <button 
        v-for="item in navItems" 
        :key="item.path"
        @click="navigate(item.path)"
        :class="[
          'flex flex-col items-center py-2 px-4 rounded-lg transition-colors',
          isActive(item.path) ? 'text-blue-500' : 'text-gray-600'
        ]">
        <span class="text-2xl mb-1">{{ item.icon }}</span>
        <span class="text-xs font-medium">{{ item.label }}</span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const navItems = computed(() => {
  const role = authStore.user?.role
  
  // Manager navigation (5 items)
  if (role === 'manager') {
    return [
      { path: '/dashboard', icon: '🏠', label: 'Dashboard' },
      { path: '/manager/shifts', icon: '⏰', label: 'Quản lý ca' },
      { path: '/manager/discrepancies', icon: '⚖️', label: 'Phê duyệt' },
      { path: '/cashier/reports', icon: '📊', label: 'Báo cáo' },
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

  // Cashier navigation
  if (role === 'cashier') {
    return [
      { path: '/dashboard', icon: '🏠', label: 'Trang chủ' },
      { path: '/cashier', icon: '💰', label: 'Thu ngân' },
      { path: '/orders', icon: '📋', label: 'Orders' },
      { path: '/shifts', icon: '⏰', label: 'Ca làm' },
      { path: '/profile', icon: '👤', label: 'Cá nhân' }
    ]
  }

  // Default navigation (waiter, etc.)
  return [
    { path: '/dashboard', icon: '🏠', label: 'Trang chủ' },
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
.safe-area-bottom {
  padding-bottom: env(safe-area-inset-bottom);
}
</style>
