<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm">
      <div class="px-4 py-4">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-800">👋 Xin chào</h1>
            <p class="text-sm text-gray-600">{{ user?.name }}</p>
          </div>
          <div class="text-right">
            <p class="text-xs text-gray-500">{{ currentDate }}</p>
            <p class="text-sm font-medium text-gray-700">{{ currentTime }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="px-4 py-4 pb-24">
      <!-- Shift Status -->
      <div v-if="hasOpenShift" class="bg-gradient-to-r from-green-500 to-emerald-500 rounded-2xl p-4 text-white shadow-lg mb-4">
        <div class="flex items-center justify-between mb-2">
          <div class="flex items-center gap-2">
            <span class="text-2xl">✅</span>
            <span class="font-bold">Ca đang mở</span>
          </div>
          <span class="text-sm opacity-90">{{ shiftDuration }}</span>
        </div>
        <p class="text-sm opacity-90">Bắt đầu: {{ formatTime(currentShift?.started_at) }}</p>
      </div>
      <div v-else class="bg-gradient-to-r from-orange-500 to-red-500 rounded-2xl p-4 text-white shadow-lg mb-4">
        <div class="flex items-center justify-between mb-2">
          <div class="flex items-center gap-2">
            <span class="text-2xl">⚠️</span>
            <span class="font-bold">Chưa mở ca</span>
          </div>
        </div>
        <button @click="$router.push('/shifts')" 
          class="mt-2 bg-white text-orange-600 px-4 py-2 rounded-lg font-medium text-sm">
          Mở ca ngay
        </button>
      </div>

      <!-- Barista Dashboard -->
      <div v-if="isBarista">
        <!-- Current Shift Info -->
        <div v-if="hasOpenShift" class="mb-4 bg-gradient-to-r from-blue-500 to-purple-500 text-white rounded-2xl p-4 shadow-lg">
          <div class="flex items-center justify-between mb-2">
            <div>
              <h3 class="font-bold text-lg">Ca làm việc</h3>
              <p class="text-sm opacity-90">{{ getShiftTypeText(currentShift.type) }}</p>
            </div>
            <div class="text-right">
              <p class="text-xs opacity-75">Thời gian</p>
              <p class="font-bold">{{ shiftDuration }}</p>
            </div>
          </div>
          <div class="text-xs opacity-90">
            Bắt đầu: {{ formatTime(currentShift.started_at) }}
          </div>
        </div>

        <!-- Barista Stats -->
        <div class="grid grid-cols-2 gap-3 mb-4">
          <div class="bg-white rounded-2xl p-4 shadow-sm">
            <div class="text-3xl mb-2">⏳</div>
            <div class="text-2xl font-bold text-yellow-600">{{ queuedOrders }}</div>
            <div class="text-xs text-gray-500">Chờ pha chế</div>
          </div>
          <div class="bg-white rounded-2xl p-4 shadow-sm">
            <div class="text-3xl mb-2">🍹</div>
            <div class="text-2xl font-bold text-blue-600">{{ inProgressOrders }}</div>
            <div class="text-xs text-gray-500">Đang pha (ca này)</div>
          </div>
          <div class="bg-white rounded-2xl p-4 shadow-sm">
            <div class="text-3xl mb-2">✅</div>
            <div class="text-2xl font-bold text-green-600">{{ readyOrders }}</div>
            <div class="text-xs text-gray-500">Sẵn sàng (ca này)</div>
          </div>
          <div class="bg-white rounded-2xl p-4 shadow-sm">
            <div class="text-3xl mb-2">🎯</div>
            <div class="text-2xl font-bold text-purple-600">{{ todayCompleted }}</div>
            <div class="text-xs text-gray-500">Hoàn tất (ca này)</div>
          </div>
        </div>

        <!-- Quick Actions for Barista -->
        <div class="mb-4">
          <h2 class="text-lg font-bold text-gray-800 mb-3">⚡ Thao tác nhanh</h2>
          <div class="grid grid-cols-2 gap-3">
            <button @click="$router.push('/barista')" 
              class="bg-gradient-to-br from-blue-500 to-blue-600 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">🍹</div>
              <div class="font-bold">Pha chế</div>
            </button>
            <button @click="$router.push('/shifts')" 
              class="bg-gradient-to-br from-purple-500 to-purple-600 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">⏰</div>
              <div class="font-bold">Ca làm</div>
            </button>
          </div>
        </div>

        <!-- Working Orders Preview -->
        <div v-if="myWorkingOrders.length > 0" class="mb-4">
          <div class="flex items-center justify-between mb-3">
            <h2 class="text-lg font-bold text-gray-800">🔥 Đang pha chế</h2>
            <button @click="$router.push('/barista')" class="text-sm text-blue-500 font-medium">
              Xem tất cả →
            </button>
          </div>
          <div class="space-y-3">
            <div v-for="order in myWorkingOrders.slice(0, 3)" :key="order.id"
              @click="$router.push('/barista')"
              class="bg-white rounded-xl p-4 shadow-sm active:scale-98 transition-transform border-l-4 border-blue-500">
              <div class="flex justify-between items-start mb-2">
                <div>
                  <h3 class="font-bold">{{ order.order_number }}</h3>
                  <p class="text-sm text-gray-600">{{ order.items?.length || 0 }} món</p>
                </div>
                <span class="bg-blue-100 text-blue-800 px-2 py-1 rounded-full text-xs font-medium">
                  Đang pha
                </span>
              </div>
              <div class="text-sm text-gray-500">
                Bắt đầu: {{ formatTime(order.accepted_at) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Waiter/Manager/Cashier Dashboard -->
      <div v-else>
        <!-- Quick Stats -->
        <div class="grid grid-cols-2 gap-3 mb-4">
          <div class="bg-white rounded-2xl p-4 shadow-sm">
            <div class="text-3xl mb-2">📋</div>
            <div class="text-2xl font-bold text-gray-800">{{ todayOrders }}</div>
            <div class="text-xs text-gray-500">Orders hôm nay</div>
          </div>
          <div class="bg-white rounded-2xl p-4 shadow-sm">
            <div class="text-3xl mb-2">💰</div>
            <div class="text-lg font-bold text-green-600">{{ formatPrice(todayRevenue) }}</div>
            <div class="text-xs text-gray-500">Doanh thu</div>
          </div>
          <div class="bg-white rounded-2xl p-4 shadow-sm">
            <div class="text-3xl mb-2">🍹</div>
            <div class="text-2xl font-bold text-blue-600">{{ inProgressOrders }}</div>
            <div class="text-xs text-gray-500">Đang pha chế</div>
          </div>
          <div class="bg-white rounded-2xl p-4 shadow-sm">
            <div class="text-3xl mb-2">⏳</div>
            <div class="text-2xl font-bold text-orange-600">{{ pendingOrders }}</div>
            <div class="text-xs text-gray-500">Chờ thanh toán</div>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="mb-4">
          <h2 class="text-lg font-bold text-gray-800 mb-3">⚡ Thao tác nhanh</h2>
          <div class="grid grid-cols-2 gap-3">
            <button @click="$router.push('/orders')" 
              class="bg-gradient-to-br from-blue-500 to-blue-600 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">📋</div>
              <div class="font-bold">Orders</div>
            </button>
            <button @click="$router.push('/shifts')" 
              class="bg-gradient-to-br from-purple-500 to-purple-600 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">⏰</div>
              <div class="font-bold">Ca làm</div>
            </button>
            
            <!-- Manager Actions -->
            <button v-if="user?.role === 'manager'" @click="$router.push('/menu')" 
              class="bg-gradient-to-br from-orange-500 to-red-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">🍽️</div>
              <div class="font-bold">Menu</div>
            </button>
            <button v-if="user?.role === 'manager'" @click="$router.push('/ingredients')" 
              class="bg-gradient-to-br from-green-500 to-emerald-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">🥬</div>
              <div class="font-bold">Nguyên liệu</div>
            </button>
            <button v-if="user?.role === 'manager'" @click="$router.push('/facilities')" 
              class="bg-gradient-to-br from-cyan-500 to-blue-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">🏢</div>
              <div class="font-bold">Cơ sở</div>
            </button>
            <button v-if="user?.role === 'manager'" @click="$router.push('/expenses')" 
              class="bg-gradient-to-br from-pink-500 to-purple-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">💸</div>
              <div class="font-bold">Chi phí</div>
            </button>
            
            <!-- Cashier Actions -->
            <button v-if="isCashier" @click="$router.push('/cashier')" 
              class="bg-gradient-to-br from-yellow-500 to-orange-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">💵</div>
              <div class="font-bold">Thu ngân</div>
            </button>
          </div>
        </div>

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
              <div class="flex justify-between items-start mb-2">
                <div>
                  <h3 class="font-bold">{{ order.order_number }}</h3>
                  <p class="text-sm text-gray-600">{{ order.customer_name || 'Khách lẻ' }}</p>
                </div>
                <span :class="getStatusBadge(order.status)" class="px-2 py-1 rounded-full text-xs font-medium">
                  {{ getStatusText(order.status) }}
                </span>
              </div>
              <div class="flex justify-between items-center text-sm">
                <span class="text-gray-500">{{ formatTime(order.created_at) }}</span>
                <span class="font-bold text-green-600">{{ formatPrice(order.total) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useShiftStore } from '../stores/shift'
import { useOrderStore } from '../stores/order'
import { useBaristaStore } from '../stores/barista'
import BottomNav from '../components/BottomNav.vue'

const authStore = useAuthStore()
const shiftStore = useShiftStore()
const orderStore = useOrderStore()
const baristaStore = useBaristaStore()

const currentTime = ref('')
const currentDate = ref('')
let timeInterval = null

// Computed
const user = computed(() => authStore.user)
const isBarista = computed(() => authStore.user?.role === 'barista')
const hasOpenShift = computed(() => shiftStore.hasOpenShift)
const currentShift = computed(() => shiftStore.currentShift)
const orders = computed(() => {
  if (isBarista.value) {
    // Barista: combine queued + my orders
    return [
      ...baristaStore.queuedOrders,
      ...baristaStore.inProgressOrders,
      ...baristaStore.readyOrders
    ]
  }
  return orderStore.orders
})
const isCashier = computed(() => authStore.user?.role === 'cashier' || authStore.user?.role === 'manager')

const recentOrders = computed(() => {
  return [...orders.value].sort((a, b) => 
    new Date(b.created_at) - new Date(a.created_at)
  )
})

// Barista-specific stats - filtered by current shift
const queuedOrders = computed(() => {
  if (!isBarista.value) return 0
  // Queue orders are not assigned yet, show all
  return baristaStore.queueCount
})

const inProgressOrders = computed(() => {
  if (!isBarista.value) return 0
  
  if (!currentShift.value) {
    // No shift: show all in progress
    return baristaStore.inProgressCount
  }
  
  // With shift: filter by shift time
  return baristaStore.inProgressOrders.filter(o => 
    isInCurrentShift(o.accepted_at)
  ).length
})

const readyOrders = computed(() => {
  if (!isBarista.value) return 0
  
  if (!currentShift.value) {
    // No shift: show all ready
    return baristaStore.readyCount
  }
  
  // With shift: filter by shift time
  return baristaStore.readyOrders.filter(o => 
    isInCurrentShift(o.ready_at)
  ).length
})

const todayCompleted = computed(() => {
  if (!isBarista.value) return 0
  
  if (!currentShift.value) {
    // No shift: count today's served orders
    const today = new Date().toDateString()
    return baristaStore.servedOrders.filter(o => 
      o.ready_at && new Date(o.ready_at).toDateString() === today
    ).length
  }
  
  // With shift: count served orders in current shift
  return baristaStore.servedOrders.filter(o => 
    isInCurrentShift(o.ready_at)
  ).length
})

// Helper function to check if timestamp is in current shift
const isInCurrentShift = (timestamp) => {
  if (!timestamp || !currentShift.value) return false
  
  const orderTime = new Date(timestamp)
  const shiftStart = new Date(currentShift.value.started_at)
  const shiftEnd = currentShift.value.ended_at ? new Date(currentShift.value.ended_at) : new Date()
  
  return orderTime >= shiftStart && orderTime <= shiftEnd
}

const myWorkingOrders = computed(() => {
  if (!isBarista.value) return []
  
  // Use barista store directly
  return baristaStore.inProgressOrders
    .filter(o => {
      if (!currentShift.value) return true
      return isInCurrentShift(o.accepted_at)
    })
    .sort((a, b) => new Date(a.accepted_at) - new Date(b.accepted_at))
})

// Waiter/Manager stats
const todayOrders = computed(() => {
  const today = new Date().toDateString()
  return orders.value.filter(o => 
    new Date(o.created_at).toDateString() === today
  ).length
})

const todayRevenue = computed(() => {
  const today = new Date().toDateString()
  return orders.value
    .filter(o => new Date(o.created_at).toDateString() === today && o.status !== 'CANCELLED')
    .reduce((sum, o) => sum + o.total, 0)
})

const pendingOrders = computed(() => {
  return orders.value.filter(o => o.status === 'CREATED').length
})

const shiftDuration = computed(() => {
  if (!currentShift.value?.started_at) return ''
  const start = new Date(currentShift.value.started_at)
  const now = new Date()
  const diff = now - start
  const hours = Math.floor(diff / 3600000)
  const minutes = Math.floor((diff % 3600000) / 60000)
  return `${hours}h ${minutes}m`
})

// Methods
const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('vi-VN', { 
    hour: '2-digit', 
    minute: '2-digit' 
  })
  currentDate.value = now.toLocaleDateString('vi-VN', { 
    weekday: 'long',
    day: 'numeric',
    month: 'long'
  })
}

const formatPrice = (price) => {
  return new Intl.NumberFormat('vi-VN', { 
    style: 'currency', 
    currency: 'VND',
    maximumFractionDigits: 0
  }).format(price)
}

const formatTime = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleTimeString('vi-VN', { 
    hour: '2-digit', 
    minute: '2-digit' 
  })
}

const getShiftTypeText = (type) => {
  const types = {
    MORNING: '☀️ Ca sáng',
    AFTERNOON: '🌤️ Ca chiều',
    EVENING: '🌙 Ca tối'
  }
  return types[type] || type
}

const getStatusBadge = (status) => {
  const badges = {
    CREATED: 'bg-gray-100 text-gray-800',
    PAID: 'bg-green-100 text-green-800',
    QUEUED: 'bg-yellow-100 text-yellow-800',
    IN_PROGRESS: 'bg-blue-100 text-blue-800',
    READY: 'bg-purple-100 text-purple-800',
    SERVED: 'bg-green-100 text-green-800',
    CANCELLED: 'bg-red-100 text-red-800'
  }
  return badges[status] || 'bg-gray-100 text-gray-800'
}

const getStatusText = (status) => {
  const texts = {
    CREATED: 'Mới',
    PAID: 'Đã thu',
    QUEUED: 'Chờ pha',
    IN_PROGRESS: 'Đang pha',
    READY: 'Sẵn sàng',
    SERVED: 'Hoàn tất',
    CANCELLED: 'Đã hủy'
  }
  return texts[status] || status
}

// Lifecycle
onMounted(async () => {
  updateTime()
  timeInterval = setInterval(updateTime, 1000)
  
  if (isBarista.value) {
    // Barista uses barista store
    await Promise.all([
      shiftStore.fetchCurrentShift(),
      baristaStore.fetchQueuedOrders(),
      baristaStore.fetchMyOrders()
    ])
  } else {
    // Other roles use order store
    await Promise.all([
      shiftStore.fetchCurrentShift(),
      orderStore.fetchOrders()
    ])
  }
})

onUnmounted(() => {
  if (timeInterval) {
    clearInterval(timeInterval)
  }
})
</script>

<style scoped>
.active\:scale-95:active {
  transform: scale(0.95);
}

.active\:scale-98:active {
  transform: scale(0.98);
}
</style>