<template>
  <nav class="bg-white shadow-lg">
    <!-- Top bar with logo and user info -->
    <div class="px-4 py-3 flex items-center justify-between border-b border-gray-200">
      <div class="flex items-center">
        <h1 class="text-xl font-bold text-blue-600">☕ Café POS</h1>
      </div>
      <div class="flex items-center gap-3">
        <router-link to="/profile" class="text-sm text-gray-600 hover:text-blue-600 hidden sm:block">
          👤 {{ userName }}
        </router-link>
        <button @click="logout" class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded-md text-sm transition-colors">
          Đăng xuất
        </button>
      </div>
    </div>

    <!-- Card-based navigation -->
    <div class="p-4">
      <!-- Manager Navigation (5-menu layout) -->
      <div v-if="userRole === 'manager'" class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-4 max-w-6xl mx-auto">
        <!-- Dashboard -->
        <router-link to="/dashboard" @click="handleNavClick" 
          class="bg-gradient-to-br from-blue-500 to-blue-600 hover:from-blue-600 hover:to-blue-700 text-white rounded-xl p-6 flex flex-col items-center justify-center min-h-[120px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
          <div class="text-3xl mb-2">🏠</div>
          <span class="text-base font-medium text-center">Dashboard</span>
        </router-link>

        <!-- Shift Management -->
        <router-link to="/manager/shifts" @click="handleNavClick"
          class="bg-gradient-to-br from-green-500 to-green-600 hover:from-green-600 hover:to-green-700 text-white rounded-xl p-6 flex flex-col items-center justify-center min-h-[120px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
          <div class="text-3xl mb-2">⏰</div>
          <span class="text-base font-medium text-center">Quản lý ca</span>
        </router-link>

        <!-- Reports -->
        <router-link to="/cashier/reports" @click="handleNavClick"
          class="bg-gradient-to-br from-cyan-500 to-cyan-600 hover:from-cyan-600 hover:to-cyan-700 text-white rounded-xl p-6 flex flex-col items-center justify-center min-h-[120px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
          <div class="text-3xl mb-2">📊</div>
          <span class="text-base font-medium text-center">Báo cáo</span>
        </router-link>

        <!-- Discrepancy Management -->
        <router-link to="/manager/discrepancies" @click="handleNavClick"
          class="bg-gradient-to-br from-orange-500 to-red-600 hover:from-orange-600 hover:to-red-700 text-white rounded-xl p-6 flex flex-col items-center justify-center min-h-[120px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
          <div class="text-3xl mb-2">⚖️</div>
          <span class="text-base font-medium text-center">Phê duyệt chênh lệch</span>
        </router-link>

        <!-- Users Management -->
        <router-link to="/users" @click="handleNavClick"
          class="bg-gradient-to-br from-indigo-500 to-indigo-600 hover:from-indigo-600 hover:to-indigo-700 text-white rounded-xl p-6 flex flex-col items-center justify-center min-h-[120px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
          <div class="text-3xl mb-2">👥</div>
          <span class="text-base font-medium text-center">Nhân viên</span>
        </router-link>

        <!-- Profile -->
        <router-link to="/profile" @click="handleNavClick"
          class="bg-gradient-to-br from-purple-500 to-purple-600 hover:from-purple-600 hover:to-purple-700 text-white rounded-xl p-6 flex flex-col items-center justify-center min-h-[120px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
          <div class="text-3xl mb-2">👤</div>
          <span class="text-base font-medium text-center">Cá nhân</span>
        </router-link>
      </div>

      <!-- Non-Manager Navigation (Original layout) -->
      <div v-else class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-3">
        <!-- Dashboard -->
        <router-link to="/dashboard" @click="handleNavClick" 
          class="bg-gradient-to-br from-blue-500 to-blue-600 hover:from-blue-600 hover:to-blue-700 text-white rounded-xl p-4 flex flex-col items-center justify-center min-h-[100px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
          <div class="text-2xl mb-2">🏠</div>
          <span class="text-sm font-medium text-center">Dashboard</span>
        </router-link>

        <!-- Shift -->
        <router-link to="/shifts" @click="handleNavClick"
          class="bg-gradient-to-br from-green-500 to-green-600 hover:from-green-600 hover:to-green-700 text-white rounded-xl p-4 flex flex-col items-center justify-center min-h-[100px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
          <div class="text-2xl mb-2">⏰</div>
          <span class="text-sm font-medium text-center">Ca làm việc</span>
        </router-link>

        <!-- Orders -->
        <router-link to="/orders" @click="handleNavClick"
          class="bg-gradient-to-br from-orange-500 to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white rounded-xl p-4 flex flex-col items-center justify-center min-h-[100px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
          <div class="text-2xl mb-2">📋</div>
          <span class="text-sm font-medium text-center">Orders</span>
        </router-link>

        <!-- Cashier Dashboard (Cashier only) -->
        <router-link v-if="userRole === 'cashier'" to="/cashier" @click="handleNavClick"
          class="bg-gradient-to-br from-pink-500 to-pink-600 hover:from-pink-600 hover:to-pink-700 text-white rounded-xl p-4 flex flex-col items-center justify-center min-h-[100px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
          <div class="text-2xl mb-2">💵</div>
          <span class="text-sm font-medium text-center">Thu ngân</span>
        </router-link>

        <!-- Cashier Reports (Cashier only) -->
        <router-link v-if="userRole === 'cashier'" to="/cashier/reports" @click="handleNavClick"
          class="bg-gradient-to-br from-cyan-500 to-cyan-600 hover:from-cyan-600 hover:to-cyan-700 text-white rounded-xl p-4 flex flex-col items-center justify-center min-h-[100px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
          <div class="text-2xl mb-2">📊</div>
          <span class="text-sm font-medium text-center">Báo cáo</span>
        </router-link>

        <!-- Cash Handovers (Cashier only) -->
        <router-link v-if="userRole === 'cashier'" to="/cashier/handovers" @click="handleNavClick"
          class="bg-gradient-to-br from-yellow-500 to-yellow-600 hover:from-yellow-600 hover:to-yellow-700 text-white rounded-xl p-4 flex flex-col items-center justify-center min-h-[100px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
          <div class="text-2xl mb-2">💰</div>
          <span class="text-sm font-medium text-center">Bàn giao tiền</span>
        </router-link>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const userRole = computed(() => {
  const role = authStore.user?.role
  console.log('🔍 Navigation - User:', authStore.user)
  console.log('🔍 Navigation - Role:', role)
  return role
})

const userName = computed(() => authStore.user?.name)

// Debug on mount
onMounted(() => {
  console.log('=== Navigation Component Mounted ===')
  console.log('Auth Store User:', authStore.user)
  console.log('User Role:', userRole.value)
  console.log('Is Manager:', userRole.value === 'manager')
  console.log('====================================')
})

// Watch for role changes
watch(userRole, (newRole, oldRole) => {
  console.log('🔄 Role changed from', oldRole, 'to', newRole)
})

const handleNavClick = () => {
  // Optional: Add any navigation handling logic here
}

const logout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.router-link-active {
  @apply ring-2 ring-white ring-opacity-50;
}
</style>