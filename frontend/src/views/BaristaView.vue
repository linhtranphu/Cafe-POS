<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm">
      <div class="px-4 py-3">
        <div class="flex items-center justify-between mb-3">
          <h1 class="text-xl font-bold text-gray-800">🍹 Barista</h1>
          <button @click="refreshAll" class="p-2 rounded-lg bg-gray-100 hover:bg-gray-200">
            🔄
          </button>
        </div>
        
        <!-- Tab Filter -->
        <div class="flex gap-2">
          <button 
            @click="activeTab = 'queue'"
            :class="[
              'flex-1 px-4 py-2 rounded-xl text-sm font-medium transition-all',
              activeTab === 'queue' 
                ? 'bg-blue-500 text-white shadow-md' 
                : 'bg-gray-100 text-gray-700'
            ]">
            ⏳ Queue <span class="ml-1">({{ queueCount }})</span>
          </button>
          <button 
            @click="activeTab = 'working'"
            :class="[
              'flex-1 px-4 py-2 rounded-xl text-sm font-medium transition-all',
              activeTab === 'working' 
                ? 'bg-blue-500 text-white shadow-md' 
                : 'bg-gray-100 text-gray-700'
            ]">
            🍹 Đang pha <span class="ml-1">({{ inProgressCount }})</span>
          </button>
          <button 
            @click="activeTab = 'ready'"
            :class="[
              'flex-1 px-4 py-2 rounded-xl text-sm font-medium transition-all',
              activeTab === 'ready' 
                ? 'bg-blue-500 text-white shadow-md' 
                : 'bg-gray-100 text-gray-700'
            ]">
            ✅ Sẵn sàng <span class="ml-1">({{ readyCount }})</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="px-4 py-4 pb-24">
      <!-- Shift Warning Banner -->
      <div v-if="!hasOpenShift" class="mb-4 bg-gradient-to-r from-orange-500 to-red-500 text-white rounded-2xl p-4 shadow-lg">
        <div class="flex items-start gap-3">
          <div class="text-3xl">⚠️</div>
          <div class="flex-1">
            <h3 class="font-bold text-lg mb-1">Chưa mở ca làm việc</h3>
            <p class="text-sm opacity-90 mb-3">Bạn cần mở ca trước khi nhận order từ queue</p>
            <button @click="$router.push('/shifts')"
              class="bg-white text-orange-600 px-4 py-2 rounded-lg font-medium text-sm active:scale-95 transition-transform">
              Mở ca ngay →
            </button>
          </div>
        </div>
      </div>

      <!-- Queue Tab -->
      <div v-if="activeTab === 'queue'">
        <div v-if="loading" class="text-center py-10">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
        </div>
        
        <div v-else-if="queuedOrders.length === 0" class="text-center py-16">
          <div class="text-6xl mb-4">🎉</div>
          <p class="text-gray-500">Không có order nào trong queue</p>
        </div>
        
        <div v-else class="space-y-3">
          <div v-for="order in queuedOrders" :key="order.id"
            class="bg-white rounded-2xl p-4 shadow-sm">
            
            <!-- Order Header -->
            <div class="flex justify-between items-start mb-3">
              <div>
                <h3 class="font-bold text-lg">{{ order.order_number }}</h3>
                <p class="text-sm text-gray-600">{{ order.customer_name || 'Khách lẻ' }}</p>
                <p class="text-xs text-gray-400">{{ formatTime(order.queued_at) }}</p>
              </div>
              <span class="bg-yellow-100 text-yellow-800 px-3 py-1 rounded-full text-xs font-medium">
                ⏳ Chờ pha
              </span>
            </div>

            <!-- Items -->
            <div class="mb-3 space-y-2">
              <div v-for="(item, idx) in order.items" :key="idx" 
                class="flex justify-between bg-gray-50 p-3 rounded-lg">
                <div>
                  <div class="font-medium">{{ item.name }}</div>
                  <div class="text-sm text-gray-500">x{{ item.quantity }}</div>
                </div>
              </div>
            </div>

            <!-- Note -->
            <div v-if="order.note" class="mb-3 p-3 bg-yellow-50 rounded-lg">
              <p class="text-sm text-gray-700">📝 {{ order.note }}</p>
            </div>

            <!-- Action -->
            <button @click="acceptOrder(order.id)"
              :disabled="!hasOpenShift"
              :class="[
                'w-full py-3 rounded-xl font-bold transition-all',
                hasOpenShift 
                  ? 'bg-blue-500 text-white active:scale-95' 
                  : 'bg-gray-300 text-gray-500 cursor-not-allowed'
              ]">
              {{ hasOpenShift ? '👍 Nhận order' : '🔒 Cần mở ca' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Working Tab -->
      <div v-if="activeTab === 'working'">
        <div v-if="loading" class="text-center py-10">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
        </div>
        
        <div v-else-if="inProgressOrders.length === 0" class="text-center py-16">
          <div class="text-6xl mb-4">☕</div>
          <p class="text-gray-500">Chưa có order nào đang pha</p>
        </div>
        
        <div v-else class="space-y-3">
          <div v-for="order in inProgressOrders" :key="order.id"
            class="bg-white rounded-2xl p-4 shadow-sm border-2 border-blue-200">
            
            <!-- Order Header -->
            <div class="flex justify-between items-start mb-3">
              <div>
                <h3 class="font-bold text-lg">{{ order.order_number }}</h3>
                <p class="text-sm text-gray-600">{{ order.customer_name || 'Khách lẻ' }}</p>
                <p class="text-xs text-gray-400">Bắt đầu: {{ formatTime(order.accepted_at) }}</p>
                <p class="text-xs text-blue-600 font-medium">⏱️ {{ getWorkingTime(order.accepted_at) }}</p>
              </div>
              <span class="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-xs font-medium">
                🍹 Đang pha
              </span>
            </div>

            <!-- Items -->
            <div class="mb-3 space-y-2">
              <div v-for="(item, idx) in order.items" :key="idx" 
                class="flex justify-between bg-blue-50 p-3 rounded-lg">
                <div>
                  <div class="font-medium">{{ item.name }}</div>
                  <div class="text-sm text-gray-500">x{{ item.quantity }}</div>
                </div>
              </div>
            </div>

            <!-- Note -->
            <div v-if="order.note" class="mb-3 p-3 bg-yellow-50 rounded-lg">
              <p class="text-sm text-gray-700">📝 {{ order.note }}</p>
            </div>

            <!-- Action -->
            <button @click="markReady(order.id)"
              class="w-full bg-green-500 text-white py-3 rounded-xl font-bold active:scale-95 transition-transform">
              ✅ Hoàn tất
            </button>
          </div>
        </div>
      </div>

      <!-- Ready Tab -->
      <div v-if="activeTab === 'ready'">
        <div v-if="loading" class="text-center py-10">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
        </div>
        
        <div v-else-if="readyOrders.length === 0" class="text-center py-16">
          <div class="text-6xl mb-4">✨</div>
          <p class="text-gray-500">Chưa có order nào sẵn sàng</p>
        </div>
        
        <div v-else class="space-y-3">
          <div v-for="order in readyOrders" :key="order.id"
            class="bg-white rounded-2xl p-4 shadow-sm border-2 border-green-200">
            
            <!-- Order Header -->
            <div class="flex justify-between items-start mb-3">
              <div>
                <h3 class="font-bold text-lg">{{ order.order_number }}</h3>
                <p class="text-sm text-gray-600">{{ order.customer_name || 'Khách lẻ' }}</p>
                <p class="text-xs text-gray-400">Hoàn tất: {{ formatTime(order.ready_at) }}</p>
                <p class="text-xs text-green-600 font-medium">⏱️ Chờ giao: {{ getWaitingTime(order.ready_at) }}</p>
              </div>
              <span class="bg-green-100 text-green-800 px-3 py-1 rounded-full text-xs font-medium">
                ✅ Sẵn sàng
              </span>
            </div>

            <!-- Items -->
            <div class="mb-3 space-y-2">
              <div v-for="(item, idx) in order.items" :key="idx" 
                class="flex justify-between bg-green-50 p-3 rounded-lg">
                <div>
                  <div class="font-medium">{{ item.name }}</div>
                  <div class="text-sm text-gray-500">x{{ item.quantity }}</div>
                </div>
              </div>
            </div>

            <!-- Note -->
            <div v-if="order.note" class="mb-3 p-3 bg-yellow-50 rounded-lg">
              <p class="text-sm text-gray-700">📝 {{ order.note }}</p>
            </div>

            <!-- Info -->
            <div class="p-3 bg-green-50 rounded-lg text-center">
              <p class="text-sm text-green-700 font-medium">
                🎉 Chờ waiter giao cho khách
              </p>
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
import { useRouter } from 'vue-router'
import { useBaristaStore } from '../stores/barista'
import { useShiftStore } from '../stores/shift'
import BottomNav from '../components/BottomNav.vue'

const router = useRouter()
const baristaStore = useBaristaStore()
const shiftStore = useShiftStore()
const activeTab = ref('queue')
let refreshInterval = null

// Computed
const loading = computed(() => baristaStore.loading)
const queuedOrders = computed(() => baristaStore.queuedOrders)
const inProgressOrders = computed(() => baristaStore.inProgressOrders)
const readyOrders = computed(() => baristaStore.readyOrders)
const queueCount = computed(() => baristaStore.queueCount)
const inProgressCount = computed(() => baristaStore.inProgressCount)
const readyCount = computed(() => baristaStore.readyCount)
const hasOpenShift = computed(() => shiftStore.hasOpenShift)

// Methods
const refreshAll = async () => {
  await Promise.all([
    baristaStore.fetchQueuedOrders(),
    baristaStore.fetchMyOrders(),
    shiftStore.fetchCurrentShift()
  ])
}

const acceptOrder = async (id) => {
  // Check shift before accepting
  if (!hasOpenShift.value) {
    if (confirm('Bạn chưa mở ca làm việc. Bạn có muốn mở ca ngay không?')) {
      router.push('/shifts')
    }
    return
  }

  try {
    await baristaStore.acceptOrder(id)
    // Switch to working tab
    activeTab.value = 'working'
  } catch (error) {
    const errorMsg = error.response?.data?.error || error.message
    
    // Handle specific error for shift requirement
    if (errorMsg.includes('shift')) {
      alert('⚠️ Bạn phải mở ca trước khi nhận order.\n\nVui lòng vào "Ca làm việc" để mở ca.')
      if (confirm('Chuyển đến trang Ca làm việc?')) {
        router.push('/shifts')
      }
    } else {
      alert('Lỗi: ' + errorMsg)
    }
  }
}

const markReady = async (id) => {
  try {
    await baristaStore.markReady(id)
    // Switch to ready tab
    activeTab.value = 'ready'
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const formatTime = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleTimeString('vi-VN', { 
    hour: '2-digit', 
    minute: '2-digit' 
  })
}

const getWorkingTime = (startTime) => {
  if (!startTime) return ''
  const start = new Date(startTime)
  const now = new Date()
  const diff = now - start
  const minutes = Math.floor(diff / 60000)
  const seconds = Math.floor((diff % 60000) / 1000)
  return `${minutes}m ${seconds}s`
}

const getWaitingTime = (readyTime) => {
  if (!readyTime) return ''
  const ready = new Date(readyTime)
  const now = new Date()
  const diff = now - ready
  const minutes = Math.floor(diff / 60000)
  return `${minutes} phút`
}

// Lifecycle
onMounted(async () => {
  await refreshAll()
  
  // Auto refresh every 10 seconds
  refreshInterval = setInterval(refreshAll, 10000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<style scoped>
.active\:scale-95:active {
  transform: scale(0.95);
}
</style>
