<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Pull to Refresh Indicator -->
    <PullToRefresh 
      :pull-distance="pullDistance" 
      :is-refreshing="isRefreshing"
      :threshold="80" />
    
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-4" style="padding-top: max(1rem, env(safe-area-inset-top))">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-800">💵 Thu ngân</h1>
            <p class="text-sm text-gray-600">Giám sát & đối soát</p>
          </div>
          <button
            @click="refreshData"
            :disabled="loading"
            class="p-3 bg-yellow-500 text-white rounded-xl shadow-lg active:scale-95 transition-transform disabled:opacity-50"
          >
            <span class="text-xl" :class="{ 'animate-spin': loading }">🔄</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24">
      <!-- Error Alert -->
      <div v-if="error" class="bg-red-50 border-2 border-red-200 rounded-2xl p-4 mb-4">
        <div class="flex items-start justify-between">
          <div class="flex items-start gap-3">
            <span class="text-2xl">⚠️</span>
            <div>
              <p class="font-medium text-red-800">Lỗi</p>
              <p class="text-sm text-red-600">{{ error }}</p>
            </div>
          </div>
          <button @click="clearError" class="text-red-600 text-xl font-bold">×</button>
        </div>
      </div>

      <!-- Current Shift Info & Close Button -->
      <div class="mb-4">
        <!-- Loading State -->
        <div v-if="cashierShiftStore.loading" class="text-center py-4 bg-white rounded-2xl">
          <div class="animate-spin text-3xl">⏳</div>
          <p class="text-sm text-gray-600 mt-2">Đang tải...</p>
        </div>

        <!-- No Shift - Show Start Button -->
        <div v-else-if="!cashierShiftStore.hasOpenCashierShift" class="bg-white rounded-2xl p-4 shadow-sm">
          <div class="bg-gray-50 rounded-xl p-4 text-center">
            <div class="text-4xl mb-2">💼</div>
            <p class="text-sm text-gray-600 mb-4">Chưa có ca thu ngân nào đang mở</p>
            
            <button
              @click="showStartModal = true"
              class="w-full py-3 bg-yellow-500 text-white rounded-xl font-bold text-base active:scale-95 transition-transform"
            >
              ➕ Bắt đầu ca thu ngân
            </button>
          </div>
        </div>

        <!-- Has Open Shift - Show Current Shift Info -->
        <div v-else class="space-y-3">
          <!-- Current Shift Info with Managed Funds -->
          <div class="bg-gradient-to-r from-yellow-500 to-orange-500 text-white rounded-xl p-4 shadow-lg">
            <!-- Shift Header -->
            <div class="flex items-center justify-between mb-3">
              <div>
                <p class="text-xs opacity-90">Ca hiện tại</p>
                <p class="text-lg font-bold">{{ cashierShiftStore.currentCashierShift?.cashier_name }}</p>
              </div>
              <div class="bg-white/20 rounded-lg px-3 py-1 backdrop-blur-sm">
                <p class="text-xs font-medium">{{ getStatusText(cashierShiftStore.currentCashierShift?.status) }}</p>
              </div>
            </div>
            
            <!-- Shift Basic Info -->
            <div class="grid grid-cols-2 gap-3 mb-4">
              <div class="bg-white/20 rounded-lg p-2 backdrop-blur-sm">
                <p class="text-xs opacity-90">Bắt đầu</p>
                <p class="text-sm font-bold">{{ formatTime(cashierShiftStore.currentCashierShift?.start_time) }}</p>
              </div>
              <div class="bg-white/20 rounded-lg p-2 backdrop-blur-sm">
                <p class="text-xs opacity-90">Tiền đầu ca</p>
                <p class="text-sm font-bold">{{ formatPrice(cashierShiftStore.currentCashierShift?.starting_float) }}</p>
              </div>
            </div>

            <!-- Managed Funds Section (Integrated) -->
            <div v-if="managedFunds" class="border-t border-white/30 pt-3">
              <div class="flex items-center gap-2 mb-3">
                <span class="text-xl">💰</span>
                <h3 class="text-sm font-bold opacity-90">Tiền đang quản lý</h3>
              </div>
              
              <!-- Funds Grid -->
              <div class="grid grid-cols-2 gap-2 mb-3">
                <!-- Received Cash -->
                <div class="bg-white/25 rounded-lg p-2 backdrop-blur-sm">
                  <div class="flex items-center gap-1 mb-1">
                    <span class="text-base">💵</span>
                    <p class="text-xs opacity-90">Tiền mặt</p>
                  </div>
                  <p class="text-sm font-bold">{{ formatPrice(managedFunds.received_cash) }}</p>
                  <p class="text-xs opacity-75">Đã nhận</p>
                </div>
                
                <!-- Received Transfer -->
                <div class="bg-white/25 rounded-lg p-2 backdrop-blur-sm">
                  <div class="flex items-center gap-1 mb-1">
                    <span class="text-base">💳</span>
                    <p class="text-xs opacity-90">Tiền CK</p>
                  </div>
                  <p class="text-sm font-bold">{{ formatPrice(managedFunds.received_transfer) }}</p>
                  <p class="text-xs opacity-75">Đã nhận</p>
                </div>
              </div>
              
              <!-- Total Managed Funds -->
              <div class="bg-white/30 rounded-lg p-2 backdrop-blur-sm mb-2">
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-1">
                    <span class="text-base">📊</span>
                    <p class="text-xs font-semibold">Tổng cộng</p>
                  </div>
                  <p class="text-lg font-bold">{{ formatPrice(managedFunds.total_managed_funds) }}</p>
                </div>
              </div>
              
              <!-- Responsibility Warning -->
              <div class="bg-white/20 rounded-lg p-2 backdrop-blur-sm">
                <div class="flex items-start gap-2">
                  <span class="text-sm">⚠️</span>
                  <div class="flex-1">
                    <p class="text-xs font-semibold mb-0.5">Bạn chịu trách nhiệm số tiền này</p>
                    <p class="text-xs opacity-90">Khi đóng ca, bạn cần bàn giao lại về quỹ</p>
                  </div>
                </div>
              </div>
              
              <!-- Handover Count -->
              <div v-if="managedFunds.handover_count > 0" class="mt-2 text-center">
                <p class="text-xs opacity-75">
                  Đã nhận {{ managedFunds.handover_count }} lần bàn giao
                </p>
              </div>
            </div>
          </div>

          <!-- Close Shift Button -->
          <button
            v-if="canCloseShift"
            @click="goToShiftClosure"
            class="w-full py-3 bg-red-500 text-white rounded-xl font-bold text-base active:scale-95 transition-transform flex items-center justify-center gap-2 shadow-lg"
          >
            <span>🔒</span>
            <span>Đóng ca thu ngân</span>
          </button>
        </div>
      </div>

      <!-- Start Shift Modal -->
      <div v-if="showStartModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-2xl p-6 max-w-md w-full">
          <h3 class="text-xl font-bold text-gray-800 mb-4">Bắt đầu ca thu ngân</h3>
          
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Tiền đầu ca (VNĐ) <span class="text-red-500">*</span>
            </label>
            <input
              v-model.number="startingFloat"
              type="number"
              step="1000"
              min="0"
              class="w-full border-2 border-gray-300 rounded-xl px-4 py-3 text-base focus:outline-none focus:border-yellow-500"
              placeholder="Nhập số tiền đầu ca"
              @keyup.enter="handleStartShift"
            />
            <p v-if="startingFloatError" class="text-xs text-red-500 mt-1">{{ startingFloatError }}</p>
          </div>

          <div class="flex gap-3">
            <button
              @click="showStartModal = false"
              class="flex-1 py-3 bg-gray-200 text-gray-800 rounded-xl font-medium active:scale-95 transition-transform"
            >
              Hủy
            </button>
            <button
              @click="handleStartShift"
              :disabled="startingLoading"
              class="flex-1 py-3 bg-yellow-500 text-white rounded-xl font-bold active:scale-95 transition-transform disabled:opacity-50"
            >
              {{ startingLoading ? 'Đang xử lý...' : 'Bắt đầu' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Handover Management Button -->
      <button 
        @click="$router.push('/cashier/handovers')"
        class="w-full bg-gradient-to-r from-yellow-500 to-orange-500 text-white rounded-xl p-4 shadow-lg active:scale-95 transition-all text-left mb-4">
        <div class="flex items-center gap-3">
          <div class="text-3xl">💰</div>
          <div class="flex-1">
            <p class="text-base font-bold">Quản lý bàn giao</p>
            <p class="text-xs opacity-90 mt-1">Xác nhận từ phục vụ</p>
          </div>
          <div class="text-2xl opacity-75">→</div>
        </div>
      </button>

      <!-- Pending Discrepancies Alert -->
      <div v-if="pendingDiscrepancies.length > 0" class="bg-yellow-50 border-2 border-yellow-300 rounded-2xl p-4 mb-4">
        <div class="flex items-center gap-3 mb-3">
          <span class="text-2xl">⚠️</span>
          <div>
            <h3 class="font-bold text-yellow-800">Sai lệch cần xử lý</h3>
            <p class="text-sm text-yellow-700">{{ pendingDiscrepancies.length }} sai lệch đang chờ</p>
          </div>
        </div>
        <button 
          @click="showDiscrepancyList = !showDiscrepancyList"
          class="text-sm text-yellow-700 font-medium underline"
        >
          {{ showDiscrepancyList ? 'Ẩn' : 'Xem chi tiết' }} →
        </button>
      </div>

      <!-- Discrepancy List (Expandable) -->
      <div v-if="showDiscrepancyList && pendingDiscrepancies.length > 0" class="space-y-3 mb-4">
        <div
          v-for="discrepancy in pendingDiscrepancies"
          :key="discrepancy.id"
          class="bg-white rounded-xl p-4 shadow-sm border-l-4 border-yellow-500"
        >
          <div class="flex justify-between items-start mb-2">
            <div>
              <h4 class="font-bold text-gray-800">Order #{{ discrepancy.order_id?.slice(-6) }}</h4>
              <p class="text-sm text-gray-600">{{ discrepancy.reason }}</p>
            </div>
            <span class="text-lg font-bold text-yellow-600">{{ formatPrice(discrepancy.amount) }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-xs text-gray-500">{{ formatDateTime(discrepancy.reported_at) }}</span>
            <button
              @click="resolveDiscrepancy(discrepancy.id)"
              class="px-4 py-2 bg-green-500 text-white rounded-lg text-sm font-medium active:scale-95 transition-transform"
            >
              ✓ Giải quyết
            </button>
          </div>
        </div>
      </div>

      <!-- Cash Reconciliation - DISABLED in Distribution Model -->
      <!-- 
        In distribution model:
        - Waiter collects cash and holds it in their shift
        - Cashier receives cash via handover process
        - Cashier does NOT reconcile waiter's cash
        - Reconciliation happens during handover confirmation
        
        This section is hidden because:
        1. Selected shift is a WAITER shift (not cashier's own shift)
        2. Cashier cannot reconcile cash they don't physically hold
        3. Cash reconciliation happens in handover flow instead
      -->
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useCashierStore } from '../stores/cashier'
import { useCashierShiftStore } from '../stores/cashierShift'
import { useShiftStore } from '../stores/shift'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'
import { printJobService } from '../services/printJob'
import { SHIFT_STATUS, CASHIER_SHIFT_STATUS, SHIFT_TYPE } from '../constants/shift'

const cashierStore = useCashierStore()
const cashierShiftStore = useCashierShiftStore()
const shiftStore = useShiftStore()

const showDiscrepancyList = ref(false)

// Managed funds state
const managedFunds = ref(null)
const loadingManagedFunds = ref(false)

// Cashier shift management state
const showStartModal = ref(false)
const startingFloat = ref(1000000) // Default 1,000,000 VND
const startingFloatError = ref(null)
const startingLoading = ref(false)

// Computed
const cashierShifts = computed(() => cashierShiftStore.cashierShifts)
const pendingDiscrepancies = computed(() => cashierStore.pendingDiscrepancies)
const loading = computed(() => cashierStore.loading)
const error = computed(() => cashierStore.error)

const canCloseShift = computed(() => {
  return cashierShiftStore.currentCashierShift && cashierShiftStore.currentCashierShift.status === 'OPEN'
})

// Reprint state
const reprintingBill = ref(null) // stores order_id being reprinted
const reprintingLabel = ref(null) // stores {orderId, itemIndex}

// Methods - Define refreshData BEFORE using it in usePullToRefresh
const refreshData = async () => {
  await cashierShiftStore.fetchCurrentCashierShift()
  await cashierShiftStore.fetchMyCashierShifts()
  await cashierStore.getPendingDiscrepancies()
  await cashierStore.fetchPendingHandovers()
  // Fetch managed funds if there's an open shift
  if (cashierShiftStore.hasOpenCashierShift) {
    await fetchManagedFunds()
  }
}

// Fetch managed funds for current cashier shift
const fetchManagedFunds = async () => {
  if (!cashierShiftStore.currentCashierShift?.id) {
    managedFunds.value = null
    return
  }
  
  loadingManagedFunds.value = true
  try {
    const response = await cashierShiftStore.getManagedFunds(cashierShiftStore.currentCashierShift.id)
    managedFunds.value = response
  } catch (error) {
    console.error('Failed to fetch managed funds:', error)
    // Don't show error to user, just log it
    managedFunds.value = null
  } finally {
    loadingManagedFunds.value = false
  }
}

// Cashier shift methods
const handleStartShift = async () => {
  startingFloatError.value = null
  
  if (!startingFloat.value || startingFloat.value < 0) {
    startingFloatError.value = 'Vui lòng nhập số tiền hợp lệ'
    return
  }
  
  startingLoading.value = true
  try {
    await cashierShiftStore.startCashierShift(startingFloat.value)
    showStartModal.value = false
    startingFloat.value = 1000000 // Reset to default
    await refreshData()
  } catch (error) {
    startingFloatError.value = error.response?.data?.error || 'Không thể bắt đầu ca'
  } finally {
    startingLoading.value = false
  }
}

const goToShiftClosure = () => {
  if (cashierShiftStore.currentCashierShift?.id) {
    window.location.href = `/#/cashier/shift-closure/${cashierShiftStore.currentCashierShift.id}`
  }
}

// Pull to refresh - Use refreshData AFTER it's defined
const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)

const resolveDiscrepancy = async (discrepancyId) => {
  if (confirm('Xác nhận đã giải quyết sai lệch này?')) {
    try {
      await cashierStore.resolveDiscrepancy(discrepancyId)
      await refreshData()
    } catch (error) {
      console.error('Resolve discrepancy failed:', error)
    }
  }
}

const clearError = () => {
  cashierStore.clearError()
}

// Utility functions
const formatPrice = (amount) => {
  if (!amount && amount !== 0) return '0₫'
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND',
    maximumFractionDigits: 0
  }).format(amount)
}

const formatDate = (date) => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleDateString('vi-VN', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  })
}

const formatDateTime = (date) => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleString('vi-VN', {
    day: '2-digit',
    month: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatTime = (date) => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleTimeString('vi-VN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getStatusText = (status) => {
  // Handle both shift status and order status using constants
  const statusMap = {
    // Shift statuses
    [SHIFT_STATUS.OPEN]: '🟢 Đang mở',
    [SHIFT_STATUS.CLOSED]: '🔴 Đã đóng',
    // Cashier Shift statuses
    [CASHIER_SHIFT_STATUS.OPEN]: '🟢 Đang mở',
    [CASHIER_SHIFT_STATUS.CLOSURE_INITIATED]: '🟡 Đang đóng',
    [CASHIER_SHIFT_STATUS.CLOSED]: '🔴 Đã đóng'
  }
  return statusMap[status] || status
}

const getShiftTypeText = (type) => {
  const types = {
    [SHIFT_TYPE.MORNING]: '☀️ Ca sáng',
    [SHIFT_TYPE.AFTERNOON]: '🌤️ Ca chiều',
    [SHIFT_TYPE.EVENING]: '🌙 Ca tối'
  }
  return types[type] || type
}

// Reprint functions
const handleReprintBill = async (orderId) => {
  if (reprintingBill.value) return
  
  reprintingBill.value = orderId
  try {
    await printJobService.reprintBill(orderId)
    alert('✅ Đã gửi lệnh in lại bill')
  } catch (error) {
    console.error('Reprint bill error:', error)
    alert('❌ Lỗi in lại bill: ' + (error.response?.data?.error || error.message))
  } finally {
    reprintingBill.value = null
  }
}

const handleReprintLabel = async (orderId, itemIndex) => {
  if (reprintingLabel.value) return
  
  reprintingLabel.value = { orderId, itemIndex }
  try {
    await printJobService.reprintLabel(orderId, itemIndex)
    alert('✅ Đã gửi lệnh in lại tem')
  } catch (error) {
    console.error('Reprint label error:', error)
    alert('❌ Lỗi in lại tem: ' + (error.response?.data?.error || error.message))
  } finally {
    reprintingLabel.value = null
  }
}

// Lifecycle
onMounted(async () => {
  // Fetch current shift first to set currentCashierShift (needed for hasOpenCashierShift check)
  await cashierShiftStore.fetchCurrentCashierShift()
  await cashierShiftStore.fetchMyCashierShifts()
  await cashierStore.getPendingDiscrepancies()
  await cashierStore.fetchPendingHandovers()

  // Fetch managed funds if there's an open shift
  if (cashierShiftStore.hasOpenCashierShift) {
    await fetchManagedFunds()
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

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.animate-spin {
  animation: spin 1s linear infinite;
}
</style>
