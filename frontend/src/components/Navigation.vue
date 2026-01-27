<template>
  <nav class="bg-white shadow-lg">
    <!-- Top bar with logo and user info -->
    <div class="px-4 py-3 flex items-center justify-between border-b border-gray-200">
      <div class="flex items-center">
        <h1 class="text-xl font-bold text-blue-600">☕ Café POS</h1>
      </div>
      <div class="flex items-center gap-3">
        <span class="text-sm text-gray-600 hidden sm:block">{{ userName }}</span>
        <button @click="logout" class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded-md text-sm transition-colors">
          Đăng xuất
        </button>
      </div>
    </div>

    <!-- Card-based navigation -->
    <div class="p-4">
      <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-3">
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

        <!-- Tables -->
        <router-link to="/tables" @click="handleNavClick"
          class="bg-gradient-to-br from-purple-500 to-purple-600 hover:from-purple-600 hover:to-purple-700 text-white rounded-xl p-4 flex flex-col items-center justify-center min-h-[100px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
          <div class="text-2xl mb-2">🪑</div>
          <span class="text-sm font-medium text-center">Bàn</span>
        </router-link>

        <!-- Cashier Dashboard (Cashier & Manager only) -->
        <router-link v-if="['cashier', 'manager'].includes(userRole)" to="/cashier" @click="handleNavClick"
          class="bg-gradient-to-br from-pink-500 to-pink-600 hover:from-pink-600 hover:to-pink-700 text-white rounded-xl p-4 flex flex-col items-center justify-center min-h-[100px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
          <div class="text-2xl mb-2">💵</div>
          <span class="text-sm font-medium text-center">Thu ngân</span>
        </router-link>

        <!-- Cashier Reports (Cashier & Manager only) -->
        <router-link v-if="['cashier', 'manager'].includes(userRole)" to="/cashier/reports" @click="handleNavClick"
          class="bg-gradient-to-br from-cyan-500 to-cyan-600 hover:from-cyan-600 hover:to-cyan-700 text-white rounded-xl p-4 flex flex-col items-center justify-center min-h-[100px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
          <div class="text-2xl mb-2">📊</div>
          <span class="text-sm font-medium text-center">Báo cáo</span>
        </router-link>

        <!-- Manager only cards -->
        <template v-if="userRole === 'manager'">
          <!-- Menu -->
          <router-link to="/menu" @click="handleNavClick"
            class="bg-gradient-to-br from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 text-white rounded-xl p-4 flex flex-col items-center justify-center min-h-[100px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
            <div class="text-2xl mb-2">🍽️</div>
            <span class="text-sm font-medium text-center">Menu</span>
          </router-link>

          <!-- Ingredients -->
          <router-link to="/ingredients" @click="handleNavClick"
            class="bg-gradient-to-br from-teal-500 to-teal-600 hover:from-teal-600 hover:to-teal-700 text-white rounded-xl p-4 flex flex-col items-center justify-center min-h-[100px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
            <div class="text-2xl mb-2">🥬</div>
            <span class="text-sm font-medium text-center">Nguyên liệu</span>
          </router-link>

          <!-- Facilities -->
          <router-link to="/facilities" @click="handleNavClick"
            class="bg-gradient-to-br from-indigo-500 to-indigo-600 hover:from-indigo-600 hover:to-indigo-700 text-white rounded-xl p-4 flex flex-col items-center justify-center min-h-[100px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
            <div class="text-2xl mb-2">🏢</div>
            <span class="text-sm font-medium text-center">Cơ sở vật chất</span>
          </router-link>

          <!-- Expenses -->
          <router-link to="/expenses" @click="handleNavClick"
            class="bg-gradient-to-br from-yellow-500 to-yellow-600 hover:from-yellow-600 hover:to-yellow-700 text-white rounded-xl p-4 flex flex-col items-center justify-center min-h-[100px] shadow-md hover:shadow-lg transform hover:scale-105 transition-all duration-200">
            <div class="text-2xl mb-2">💰</div>
            <span class="text-sm font-medium text-center">Chi phí</span>
          </router-link>
        </template>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const userRole = computed(() => authStore.user?.role)
const userName = computed(() => authStore.user?.name)

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