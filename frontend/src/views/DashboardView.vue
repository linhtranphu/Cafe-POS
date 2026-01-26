<template>
  <div class="min-h-screen bg-gray-100">
    <Navigation />
    <div class="p-4">
      <!-- Welcome Header -->
      <div class="bg-white rounded-xl p-6 mb-6 shadow-sm">
        <h2 class="text-2xl font-bold text-gray-800 mb-2 text-center">Chào mừng đến với Café POS</h2>
        <p class="text-gray-600 text-center mb-4">
          Bạn đã đăng nhập với quyền {{ user?.role === 'manager' ? 'Quản lý' : user?.role === 'waiter' ? 'Nhân viên' : 'Thu ngân' }}
        </p>
        <div class="flex justify-center">
          <span class="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm font-medium">
            {{ userName }}
          </span>
        </div>
      </div>

      <!-- Quick Actions Grid -->
      <div class="grid grid-cols-2 gap-4 mb-6">
        <!-- Menu Management -->
        <div v-if="user?.role === 'manager'" @click="$router.push('/menu')" 
             class="bg-gradient-to-br from-orange-500 to-red-500 rounded-xl p-6 text-white cursor-pointer transform hover:scale-105 transition-all duration-200 shadow-lg">
          <div class="text-3xl mb-3">🍽️</div>
          <h3 class="font-bold text-lg mb-1">Menu</h3>
          <p class="text-sm opacity-90">Quản lý thực đơn</p>
        </div>

        <!-- Ingredients -->
        <div v-if="user?.role === 'manager'" @click="$router.push('/ingredients')" 
             class="bg-gradient-to-br from-green-500 to-emerald-500 rounded-xl p-6 text-white cursor-pointer transform hover:scale-105 transition-all duration-200 shadow-lg">
          <div class="text-3xl mb-3">🥬</div>
          <h3 class="font-bold text-lg mb-1">Nguyên liệu</h3>
          <p class="text-sm opacity-90">Quản lý kho</p>
        </div>

        <!-- Facilities -->
        <div v-if="user?.role === 'manager'" @click="$router.push('/facilities')" 
             class="bg-gradient-to-br from-blue-500 to-indigo-500 rounded-xl p-6 text-white cursor-pointer transform hover:scale-105 transition-all duration-200 shadow-lg">
          <div class="text-3xl mb-3">🏢</div>
          <h3 class="font-bold text-lg mb-1">Cơ sở vật chất</h3>
          <p class="text-sm opacity-90">Quản lý tài sản</p>
        </div>

        <!-- Expenses -->
        <div v-if="user?.role === 'manager'" @click="$router.push('/expenses')" 
             class="bg-gradient-to-br from-purple-500 to-pink-500 rounded-xl p-6 text-white cursor-pointer transform hover:scale-105 transition-all duration-200 shadow-lg">
          <div class="text-3xl mb-3">💰</div>
          <h3 class="font-bold text-lg mb-1">Chi phí</h3>
          <p class="text-sm opacity-90">Quản lý tài chính</p>
        </div>

        <!-- Orders (for all roles) -->
        <div @click="$router.push('/orders')" 
             class="bg-gradient-to-br from-teal-500 to-cyan-500 rounded-xl p-6 text-white cursor-pointer transform hover:scale-105 transition-all duration-200 shadow-lg">
          <div class="text-3xl mb-3">📋</div>
          <h3 class="font-bold text-lg mb-1">Đơn hàng</h3>
          <p class="text-sm opacity-90">Quản lý order</p>
        </div>

        <!-- Reports (Manager only) -->
        <div v-if="user?.role === 'manager'" @click="$router.push('/reports')" 
             class="bg-gradient-to-br from-yellow-500 to-orange-500 rounded-xl p-6 text-white cursor-pointer transform hover:scale-105 transition-all duration-200 shadow-lg">
          <div class="text-3xl mb-3">📊</div>
          <h3 class="font-bold text-lg mb-1">Báo cáo</h3>
          <p class="text-sm opacity-90">Thống kê doanh thu</p>
        </div>
      </div>

      <!-- Role Permissions Card -->
      <div class="bg-white rounded-xl p-6 shadow-sm">
        <h3 class="text-lg font-semibold text-gray-800 mb-4 text-center">Quyền hạn của bạn</h3>
        <div class="space-y-3">
          <div v-for="permission in permissions" :key="permission" 
               class="flex items-center p-3 bg-gray-50 rounded-lg">
            <div class="w-6 h-6 bg-green-500 rounded-full flex items-center justify-center mr-3">
              <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
              </svg>
            </div>
            <span class="text-gray-700 text-sm">{{ permission }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import Navigation from '../components/Navigation.vue'

const authStore = useAuthStore()

const user = computed(() => authStore.user)
const userName = computed(() => authStore.user?.name || 'User')

const permissions = computed(() => {
  switch (user.value?.role) {
    case 'manager':
      return [
        'Xem tất cả order',
        'Xem báo cáo doanh thu',
        'In lại bill',
        'Chỉnh sửa/hủy order đã thanh toán',
        'Quản lý menu & giá',
        'Quản lý bàn',
        'Quản lý user'
      ]
    case 'waiter':
      return [
        'Tạo order mới',
        'Nhập món, số lượng',
        'Gắn bàn',
        'Xem & thông báo tổng tiền',
        'Chọn phương thức thanh toán',
        'Xác nhận đã thu tiền',
        'In bill'
      ]
    case 'cashier':
      return [
        'Xem order đã tạo',
        'Thu tiền',
        'Xác nhận thanh toán',
        'In/in lại bill'
      ]
    default:
      return []
  }
})
</script>

<style scoped>
/* Tailwind handles all styling */
</style>