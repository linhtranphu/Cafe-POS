<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Pull to Refresh Indicator -->
    <PullToRefresh 
      :pull-distance="pullDistance" 
      :is-refreshing="isRefreshing"
      :threshold="80" />
    
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm">
      <div class="px-4 py-3">
        <div class="flex items-center justify-between mb-3">
          <h1 class="text-xl font-bold text-gray-800">⏰ Quản lý ca làm việc</h1>
          <button @click="refreshData" class="p-2 rounded-lg bg-blue-500 text-white">
            🔄
          </button>
        </div>
        
        <!-- Filter Tabs -->
        <div class="flex gap-2 overflow-x-auto pb-2">
          <button @click="filterStatus = 'all'" 
            :class="filterStatus === 'all' ? 'bg-blue-500 text-white' : 'bg-gray-200 text-gray-700'"
            class="px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap">
            Tất cả
          </button>
          <button @click="filterStatus = 'OPEN'" 
            :class="filterStatus === 'OPEN' ? 'bg-green-500 text-white' : 'bg-gray-200 text-gray-700'"
            class="px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap">
            Đang mở
          </button>
          <button @click="filterStatus = 'CLOSED'" 
            :class="filterStatus === 'CLOSED' ? 'bg-gray-500 text-white' : 'bg-gray-200 text-gray-700'"
            class="px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap">
            Đã đóng
          </button>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="px-4 py-4 pb-24">
      <!-- Stats Cards - Single Row -->
      <div class="bg-gradient-to-br from-blue-500 to-purple-500 rounded-xl p-4 mb-4 text-white shadow-lg">
        <div class="text-xs opacity-90 mb-2">Tổng quan ca làm</div>
        <div class="grid grid-cols-4 gap-1.5">
          <div class="text-center">
            <div class="text-lg font-bold">{{ todayShifts.length }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap">Hôm nay</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold">{{ openWaiterShifts.length }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap">Waiter</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold">{{ openBaristaShifts.length }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap">Barista</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold">{{ openCashierShifts.length }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap">Cashier</div>
          </div>
        </div>
      </div>

      <!-- Waiter Shifts Section -->
      <div class="mb-4">
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-lg font-bold text-gray-800">🍽️ Ca Waiter</h2>
          <span class="text-sm text-gray-500">{{ filteredWaiterShifts.length }} ca</span>
        </div>
        
        <div v-if="filteredWaiterShifts.length === 0" class="text-center py-8 bg-white rounded-xl">
          <div class="text-4xl mb-2">📭</div>
          <p class="text-gray-500">Không có ca nào</p>
        </div>
        
        <div v-else class="space-y-3">
          <div v-for="shift in filteredWaiterShifts" :key="shift.id"
            @click="viewShiftDetails(shift, 'waiter')"
            class="bg-white rounded-2xl p-4 shadow-sm active:scale-98 transition-transform border-l-4 border-blue-500">
            
            <!-- Shift Header -->
            <div class="flex justify-between items-start mb-3">
              <div>
                <h3 class="font-bold text-lg">{{ shift.user_name }}</h3>
                <p class="text-sm text-gray-600">{{ getRoleTypeText(shift.role_type) }}</p>
                <p class="text-xs text-gray-400">{{ getShiftTypeText(shift.type) }}</p>
              </div>
              <span :class="getStatusClass(shift.status)" class="px-3 py-1 rounded-full text-xs font-medium">
                {{ getStatusText(shift.status) }}
              </span>
            </div>

            <!-- Shift Info -->
            <div class="grid grid-cols-2 gap-2 text-sm mb-3">
              <div>
                <span class="text-gray-500">Bắt đầu:</span>
                <p class="font-medium">{{ formatDateTime(shift.started_at) }}</p>
              </div>
              <div v-if="shift.ended_at">
                <span class="text-gray-500">Kết thúc:</span>
                <p class="font-medium">{{ formatDateTime(shift.ended_at) }}</p>
              </div>
              <div v-else>
                <span class="text-gray-500">Thời gian:</span>
                <p class="font-medium text-green-600">{{ calculateDuration(shift.started_at) }}</p>
              </div>
            </div>

            <!-- Shift Stats -->
            <div v-if="shift.status === 'CLOSED'" class="grid grid-cols-3 gap-2 pt-3 border-t">
              <div class="text-center">
                <p class="text-xs text-gray-500">Tiền đầu</p>
                <p class="font-bold text-sm">{{ formatPrice(shift.start_cash) }}</p>
              </div>
              <div class="text-center">
                <p class="text-xs text-gray-500">Tiền cuối</p>
                <p class="font-bold text-sm">{{ formatPrice(shift.end_cash) }}</p>
              </div>
              <div class="text-center">
                <p class="text-xs text-gray-500">Doanh thu</p>
                <p class="font-bold text-sm text-green-600">{{ formatPrice(shift.total_revenue) }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Barista Shifts Section -->
      <div class="mb-4">
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-lg font-bold text-gray-800">🍹 Ca Barista</h2>
          <span class="text-sm text-gray-500">{{ filteredBaristaShifts.length }} ca</span>
        </div>
        
        <div v-if="filteredBaristaShifts.length === 0" class="text-center py-8 bg-white rounded-xl">
          <div class="text-4xl mb-2">📭</div>
          <p class="text-gray-500">Không có ca nào</p>
        </div>
        
        <div v-else class="space-y-3">
          <div v-for="shift in filteredBaristaShifts" :key="shift.id"
            @click="viewShiftDetails(shift, 'barista')"
            class="bg-white rounded-2xl p-4 shadow-sm active:scale-98 transition-transform border-l-4 border-purple-500">
            
            <!-- Shift Header -->
            <div class="flex justify-between items-start mb-3">
              <div>
                <h3 class="font-bold text-lg">{{ shift.user_name }}</h3>
                <p class="text-sm text-gray-600">{{ getRoleTypeText(shift.role_type) }}</p>
                <p class="text-xs text-gray-400">{{ getShiftTypeText(shift.type) }}</p>
              </div>
              <span :class="getStatusClass(shift.status)" class="px-3 py-1 rounded-full text-xs font-medium">
                {{ getStatusText(shift.status) }}
              </span>
            </div>

            <!-- Shift Info -->
            <div class="grid grid-cols-2 gap-2 text-sm mb-3">
              <div>
                <span class="text-gray-500">Bắt đầu:</span>
                <p class="font-medium">{{ formatDateTime(shift.started_at) }}</p>
              </div>
              <div v-if="shift.ended_at">
                <span class="text-gray-500">Kết thúc:</span>
                <p class="font-medium">{{ formatDateTime(shift.ended_at) }}</p>
              </div>
              <div v-else>
                <span class="text-gray-500">Thời gian:</span>
                <p class="font-medium text-green-600">{{ calculateDuration(shift.started_at) }}</p>
              </div>
            </div>

            <!-- Shift Stats -->
            <div v-if="shift.status === 'CLOSED'" class="grid grid-cols-3 gap-2 pt-3 border-t">
              <div class="text-center">
                <p class="text-xs text-gray-500">Tiền đầu</p>
                <p class="font-bold text-sm">{{ formatPrice(shift.start_cash) }}</p>
              </div>
              <div class="text-center">
                <p class="text-xs text-gray-500">Tiền cuối</p>
                <p class="font-bold text-sm">{{ formatPrice(shift.end_cash) }}</p>
              </div>
              <div class="text-center">
                <p class="text-xs text-gray-500">Doanh thu</p>
                <p class="font-bold text-sm text-green-600">{{ formatPrice(shift.total_revenue) }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Cashier Shifts Section -->
      <div class="mb-4">
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-lg font-bold text-gray-800">💵 Ca Thu ngân</h2>
          <span class="text-sm text-gray-500">{{ filteredCashierShifts.length }} ca</span>
        </div>
        
        <div v-if="filteredCashierShifts.length === 0" class="text-center py-8 bg-white rounded-xl">
          <div class="text-4xl mb-2">📭</div>
          <p class="text-gray-500">Không có ca nào</p>
        </div>
        
        <div v-else class="space-y-3">
          <div v-for="shift in filteredCashierShifts" :key="shift.id"
            @click="viewShiftDetails(shift, 'cashier')"
            class="bg-white rounded-2xl p-4 shadow-sm active:scale-98 transition-transform border-l-4 border-yellow-500">
            
            <!-- Shift Header -->
            <div class="flex justify-between items-start mb-3">
              <div>
                <h3 class="font-bold text-lg">{{ shift.cashier_name }}</h3>
                <p class="text-sm text-gray-600">💵 Thu ngân</p>
              </div>
              <span :class="getCashierStatusClass(shift.status)" class="px-3 py-1 rounded-full text-xs font-medium">
                {{ getCashierStatusText(shift.status) }}
              </span>
            </div>

            <!-- Shift Info -->
            <div class="grid grid-cols-2 gap-2 text-sm mb-3">
              <div>
                <span class="text-gray-500">Bắt đầu:</span>
                <p class="font-medium">{{ formatDateTime(shift.opened_at) }}</p>
              </div>
              <div v-if="shift.closed_at">
                <span class="text-gray-500">Kết thúc:</span>
                <p class="font-medium">{{ formatDateTime(shift.closed_at) }}</p>
              </div>
              <div v-else>
                <span class="text-gray-500">Thời gian:</span>
                <p class="font-medium text-green-600">{{ calculateDuration(shift.opened_at) }}</p>
              </div>
            </div>

            <!-- Cashier Shift Stats -->
            <div v-if="shift.status === 'CLOSED'" class="grid grid-cols-3 gap-2 pt-3 border-t">
              <div class="text-center">
                <p class="text-xs text-gray-500">Tiền mặt</p>
                <p class="font-bold text-sm">{{ formatPrice(shift.actual_cash) }}</p>
              </div>
              <div class="text-center">
                <p class="text-xs text-gray-500">Dự kiến</p>
                <p class="font-bold text-sm">{{ formatPrice(shift.expected_cash) }}</p>
              </div>
              <div class="text-center">
                <p class="text-xs text-gray-500">Chênh lệch</p>
                <p :class="shift.variance >= 0 ? 'text-green-600' : 'text-red-600'" class="font-bold text-sm">
                  {{ shift.variance >= 0 ? '+' : '' }}{{ formatPrice(Math.abs(shift.variance)) }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />

    <!-- Shift Detail Modal -->
    <transition name="slide-up">
      <div v-if="showDetailModal" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full max-h-[85vh] overflow-y-auto">
          <div class="sticky top-0 bg-white px-4 py-4 border-b flex justify-between items-center">
            <h3 class="text-lg font-bold">Chi tiết ca làm việc</h3>
            <button @click="closeDetailModal" class="text-2xl text-gray-400">×</button>
          </div>
          
          <div class="px-4 py-4">
            <!-- Waiter Shift Details -->
            <div v-if="selectedShiftType === 'waiter' && selectedShift">
              <div class="bg-gradient-to-r from-blue-500 to-purple-500 rounded-xl p-4 text-white mb-4">
                <h4 class="font-bold text-lg mb-1">{{ selectedShift.user_name }}</h4>
                <p class="text-sm opacity-90">{{ getRoleTypeText(selectedShift.role_type) }}</p>
                <p class="text-xs opacity-75">{{ getShiftTypeText(selectedShift.type) }}</p>
              </div>

              <div class="space-y-3">
                <div class="bg-gray-50 rounded-lg p-3">
                  <p class="text-xs text-gray-500 mb-1">Trạng thái</p>
                  <span :class="getStatusClass(selectedShift.status)" class="px-3 py-1 rounded-full text-sm font-medium">
                    {{ getStatusText(selectedShift.status) }}
                  </span>
                </div>

                <div class="bg-gray-50 rounded-lg p-3">
                  <p class="text-xs text-gray-500 mb-1">Thời gian bắt đầu</p>
                  <p class="font-medium">{{ formatDateTime(selectedShift.started_at) }}</p>
                </div>

                <div v-if="selectedShift.ended_at" class="bg-gray-50 rounded-lg p-3">
                  <p class="text-xs text-gray-500 mb-1">Thời gian kết thúc</p>
                  <p class="font-medium">{{ formatDateTime(selectedShift.ended_at) }}</p>
                </div>

                <div v-if="selectedShift.status === 'CLOSED'" class="grid grid-cols-2 gap-3">
                  <div class="bg-blue-50 rounded-lg p-3">
                    <p class="text-xs text-blue-600 mb-1">Tiền đầu ca</p>
                    <p class="font-bold text-lg">{{ formatPrice(selectedShift.start_cash) }}</p>
                  </div>
                  <div class="bg-green-50 rounded-lg p-3">
                    <p class="text-xs text-green-600 mb-1">Tiền cuối ca</p>
                    <p class="font-bold text-lg">{{ formatPrice(selectedShift.end_cash) }}</p>
                  </div>
                  <div class="bg-purple-50 rounded-lg p-3 col-span-2">
                    <p class="text-xs text-purple-600 mb-1">Tổng doanh thu</p>
                    <p class="font-bold text-2xl">{{ formatPrice(selectedShift.total_revenue) }}</p>
                  </div>
                </div>
              </div>
            </div>

            <!-- Cashier Shift Details -->
            <div v-if="selectedShiftType === 'cashier' && selectedShift">
              <div class="bg-gradient-to-r from-yellow-500 to-orange-500 rounded-xl p-4 text-white mb-4">
                <h4 class="font-bold text-lg mb-1">{{ selectedShift.cashier_name }}</h4>
                <p class="text-sm opacity-90">💵 Thu ngân</p>
              </div>

              <div class="space-y-3">
                <div class="bg-gray-50 rounded-lg p-3">
                  <p class="text-xs text-gray-500 mb-1">Trạng thái</p>
                  <span :class="getCashierStatusClass(selectedShift.status)" class="px-3 py-1 rounded-full text-sm font-medium">
                    {{ getCashierStatusText(selectedShift.status) }}
                  </span>
                </div>

                <div class="bg-gray-50 rounded-lg p-3">
                  <p class="text-xs text-gray-500 mb-1">Thời gian mở ca</p>
                  <p class="font-medium">{{ formatDateTime(selectedShift.opened_at) }}</p>
                </div>

                <div v-if="selectedShift.closed_at" class="bg-gray-50 rounded-lg p-3">
                  <p class="text-xs text-gray-500 mb-1">Thời gian đóng ca</p>
                  <p class="font-medium">{{ formatDateTime(selectedShift.closed_at) }}</p>
                </div>

                <div v-if="selectedShift.status === 'CLOSED'" class="space-y-3">
                  <div class="grid grid-cols-2 gap-3">
                    <div class="bg-blue-50 rounded-lg p-3">
                      <p class="text-xs text-blue-600 mb-1">Tiền mặt thực tế</p>
                      <p class="font-bold text-lg">{{ formatPrice(selectedShift.actual_cash) }}</p>
                    </div>
                    <div class="bg-green-50 rounded-lg p-3">
                      <p class="text-xs text-green-600 mb-1">Tiền dự kiến</p>
                      <p class="font-bold text-lg">{{ formatPrice(selectedShift.expected_cash) }}</p>
                    </div>
                  </div>

                  <div :class="selectedShift.variance >= 0 ? 'bg-green-50' : 'bg-red-50'" class="rounded-lg p-3">
                    <p :class="selectedShift.variance >= 0 ? 'text-green-600' : 'text-red-600'" class="text-xs mb-1">
                      Chênh lệch
                    </p>
                    <p :class="selectedShift.variance >= 0 ? 'text-green-600' : 'text-red-600'" class="font-bold text-2xl">
                      {{ selectedShift.variance >= 0 ? '+' : '' }}{{ formatPrice(Math.abs(selectedShift.variance)) }}
                    </p>
                  </div>

                  <div v-if="selectedShift.variance_reason" class="bg-yellow-50 rounded-lg p-3">
                    <p class="text-xs text-yellow-600 mb-1">Lý do chênh lệch</p>
                    <p class="text-sm">{{ selectedShift.variance_reason }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useShiftStore } from '../stores/shift'
import { useCashierShiftStore } from '../stores/cashierShift'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'
import { formatDate, formatDateTime, formatPrice } from '../utils/formatters'

const shiftStore = useShiftStore()
const cashierShiftStore = useCashierShiftStore()

const filterStatus = ref('all')
const showDetailModal = ref(false)
const selectedShift = ref(null)
const selectedShiftType = ref(null) // 'waiter' or 'cashier'

// Computed
const waiterShifts = computed(() => shiftStore.shifts || [])
const cashierShifts = computed(() => cashierShiftStore.shifts || [])

const allShifts = computed(() => {
  const waiter = waiterShifts.value || []
  const cashier = cashierShifts.value || []
  return [...waiter, ...cashier]
})

const openWaiterShifts = computed(() => {
  const shifts = waiterShifts.value || []
  return shifts.filter(s => s && s.status === 'OPEN' && s.role_type === 'waiter')
})

const openBaristaShifts = computed(() => {
  const shifts = waiterShifts.value || []
  return shifts.filter(s => s && s.status === 'OPEN' && s.role_type === 'barista')
})

const openCashierShifts = computed(() => {
  const shifts = cashierShifts.value || []
  return shifts.filter(s => s && s.status === 'OPEN')
})

const todayShifts = computed(() => {
  const today = new Date().toDateString()
  const shifts = allShifts.value || []
  return shifts.filter(s => {
    if (!s) return false
    const shiftDate = s.started_at || s.opened_at
    return shiftDate && new Date(shiftDate).toDateString() === today
  })
})

const filteredWaiterShifts = computed(() => {
  const shifts = waiterShifts.value || []
  const waiterOnly = shifts.filter(s => s && s.role_type === 'waiter')
  if (filterStatus.value === 'all') return waiterOnly
  return waiterOnly.filter(s => s && s.status === filterStatus.value)
})

const filteredBaristaShifts = computed(() => {
  const shifts = waiterShifts.value || []
  const baristaOnly = shifts.filter(s => s && s.role_type === 'barista')
  if (filterStatus.value === 'all') return baristaOnly
  return baristaOnly.filter(s => s && s.status === filterStatus.value)
})

const filteredCashierShifts = computed(() => {
  const shifts = cashierShifts.value || []
  if (filterStatus.value === 'all') return shifts
  return shifts.filter(s => s && s.status === filterStatus.value)
})

// Methods
const refreshData = async () => {
  await Promise.all([
    shiftStore.fetchAllShifts(),
    cashierShiftStore.fetchAllShifts()
  ])
}

// Pull to refresh
const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)

const getRoleTypeText = (roleType) => {
  const roles = {
    waiter: '🍽️ Phục vụ',
    barista: '🍹 Pha chế',
    cashier: '💵 Thu ngân'
  }
  return roles[roleType] || roleType
}

const getShiftTypeText = (type) => {
  const types = {
    MORNING: '☀️ Ca sáng',
    AFTERNOON: '🌤️ Ca chiều',
    EVENING: '🌙 Ca tối'
  }
  return types[type] || type
}

const getStatusClass = (status) => {
  return status === 'OPEN' 
    ? 'bg-green-100 text-green-800' 
    : 'bg-gray-100 text-gray-800'
}

const getStatusText = (status) => {
  return status === 'OPEN' ? 'Đang mở' : 'Đã đóng'
}

const getCashierStatusClass = (status) => {
  const classes = {
    OPEN: 'bg-green-100 text-green-800',
    CLOSURE_INITIATED: 'bg-yellow-100 text-yellow-800',
    CLOSED: 'bg-gray-100 text-gray-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getCashierStatusText = (status) => {
  const texts = {
    OPEN: 'Đang mở',
    CLOSURE_INITIATED: 'Đang đóng',
    CLOSED: 'Đã đóng'
  }
  return texts[status] || status
}

const calculateDuration = (startTime) => {
  if (!startTime) return ''
  const start = new Date(startTime)
  const now = new Date()
  const diff = now - start
  const hours = Math.floor(diff / 3600000)
  const minutes = Math.floor((diff % 3600000) / 60000)
  return `${hours}h ${minutes}m`
}

const viewShiftDetails = (shift, type) => {
  selectedShift.value = shift
  selectedShiftType.value = type
  showDetailModal.value = true
}

const closeDetailModal = () => {
  showDetailModal.value = false
  selectedShift.value = null
  selectedShiftType.value = null
}

// Lifecycle
onMounted(async () => {
  await refreshData()
})
</script>

<style scoped>
.active\:scale-98:active {
  transform: scale(0.98);
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: transform 0.3s ease;
}

.slide-up-enter-from {
  transform: translateY(100%);
}

.slide-up-leave-to {
  transform: translateY(100%);
}
</style>
