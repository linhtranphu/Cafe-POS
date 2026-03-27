<template>
  <div class="bg-white rounded-2xl p-4 shadow-sm mb-4">
    <h2 class="text-lg font-bold text-gray-800 mb-4">📅 Ca thu ngân</h2>
    
    <!-- Loading State -->
    <div v-if="loading" class="text-center py-4">
      <div class="animate-spin text-3xl">⏳</div>
      <p class="text-sm text-gray-600 mt-2">Đang tải...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-red-50 border-2 border-red-200 rounded-xl p-4">
      <div class="flex items-start gap-3">
        <span class="text-2xl">⚠️</span>
        <div>
          <p class="font-medium text-red-800">Lỗi</p>
          <p class="text-sm text-red-600">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- No Shift - Show Start Button -->
    <div v-else-if="!hasOpenShift" class="space-y-4">
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

    <!-- Has Open Shift - Show Details -->
    <div v-else class="space-y-3">
      <!-- Financial Summary -->
      <div class="bg-white border-2 border-gray-200 rounded-xl p-4">
        <h3 class="text-sm font-bold text-gray-800 mb-3">💰 Tổng hợp tài chính</h3>

        <!-- Loading State -->
        <div v-if="loadingHandovers" class="text-center py-4">
          <div class="inline-block animate-spin rounded-full h-6 w-6 border-b-2 border-blue-500"></div>
          <p class="text-xs text-gray-500 mt-2">Đang tải...</p>
        </div>

        <!-- Financial Summary Cards -->
        <div v-else class="space-y-3">
          <!-- 1. Already Handed Over Section -->
          <div class="bg-gradient-to-br from-green-500 to-emerald-600 text-white rounded-xl p-4 shadow-md">
            <div class="flex items-center justify-between mb-3">
              <div>
                <p class="text-sm font-bold opacity-90">✅ Đã bàn giao</p>
                <p class="text-xs opacity-75">{{ openWaiterShiftsCount }} ca đang mở</p>
              </div>
              <p class="text-2xl font-bold">{{ formatPriceShort(handedOverSummary.total) }}</p>
            </div>
            <div class="grid grid-cols-2 gap-2">
              <div class="bg-white/20 rounded-lg p-2 backdrop-blur-sm">
                <p class="text-xs opacity-90">💵 Tiền mặt</p>
                <p class="text-base font-bold">{{ formatPrice(handedOverSummary.cash) }}</p>
              </div>
              <div class="bg-white/20 rounded-lg p-2 backdrop-blur-sm">
                <p class="text-xs opacity-90">💳 Chuyển khoản</p>
                <p class="text-base font-bold">{{ formatPrice(handedOverSummary.transfer) }}</p>
              </div>
            </div>
          </div>

          <!-- 2. Not Handed Over Yet Section -->
          <div class="bg-gradient-to-br from-blue-500 to-indigo-600 text-white rounded-xl p-4 shadow-md">
            <div class="flex items-center justify-between mb-3">
              <div>
                <p class="text-sm font-bold opacity-90">💼 Chưa bàn giao</p>
                <p class="text-xs opacity-75">Còn lại trong ca</p>
              </div>
              <p class="text-2xl font-bold">{{ formatPriceShort(notHandedOverSummary.total) }}</p>
            </div>
            <div class="grid grid-cols-2 gap-2">
              <div class="bg-white/20 rounded-lg p-2 backdrop-blur-sm">
                <p class="text-xs opacity-90">💵 Tiền mặt</p>
                <p class="text-base font-bold">{{ formatPrice(notHandedOverSummary.cash) }}</p>
              </div>
              <div class="bg-white/20 rounded-lg p-2 backdrop-blur-sm">
                <p class="text-xs opacity-90">💳 Chuyển khoản</p>
                <p class="text-base font-bold">{{ formatPrice(notHandedOverSummary.transfer) }}</p>
              </div>
            </div>
          </div>

          <!-- Discrepancy Alert (if any) -->
          <div v-if="totalDiscrepancy !== 0" class="bg-red-50 border-2 border-red-200 rounded-lg p-3">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <span class="text-red-600 text-lg">⚠️</span>
                <div>
                  <p class="text-xs font-bold text-red-800">CHÊNH LỆCH</p>
                  <p class="text-xs text-red-600">Cần kiểm tra</p>
                </div>
              </div>
              <span class="text-lg font-bold text-red-600">{{ formatPrice(totalDiscrepancy) }}</span>
            </div>
          </div>

          <!-- Breakdown by Waiter -->
          <div class="pt-2 border-t-2 border-gray-200">
            <h4 class="text-xs font-bold text-gray-700 mb-2">👥 Chi tiết theo nhân viên</h4>
            
            <!-- No waiters -->
            <div v-if="waiterFinancialBreakdown.length === 0" class="text-center py-4 bg-gray-50 rounded-lg">
              <p class="text-xs text-gray-500">Không có ca đang mở</p>
            </div>

            <!-- Waiter Cards -->
            <div v-else class="space-y-2">
              <div 
                v-for="waiter in waiterFinancialBreakdown" 
                :key="waiter.userId"
                class="bg-gradient-to-r from-blue-500 to-purple-500 text-white rounded-xl p-3 shadow-md"
              >
                <!-- Waiter Header -->
                <div class="flex items-center justify-between mb-3">
                  <div>
                    <p class="font-bold text-base">{{ waiter.userName }}</p>
                    <p class="text-xs opacity-90">{{ waiter.shiftType ? getShiftTypeText(waiter.shiftType) : 'Nhiều ca' }}</p>
                  </div>
                  <button
                    @click="toggleWaiterDetails(waiter.userId)"
                    class="bg-white/20 rounded-lg px-2 py-1 text-xs font-medium backdrop-blur-sm"
                  >
                    {{ expandedWaiters.includes(waiter.userId) ? '▼' : '▶' }}
                  </button>
                </div>

                <!-- 3 Sections: Revenue, Remaining, Handed Over -->
                <div class="space-y-2">
                  <!-- Revenue Section -->
                  <div class="bg-white/20 rounded-lg p-2 backdrop-blur-sm">
                    <p class="text-xs opacity-90 mb-1.5">📊 Doanh thu</p>
                    <div class="grid grid-cols-2 gap-2">
                      <div>
                        <p class="text-xs opacity-90">💵 Mặt</p>
                        <p class="text-sm font-bold">{{ formatPrice(waiter.totalCash) }}</p>
                      </div>
                      <div>
                        <p class="text-xs opacity-90">💳 CK</p>
                        <p class="text-sm font-bold">{{ formatPrice(waiter.totalTransfer) }}</p>
                      </div>
                    </div>
                  </div>

                  <!-- Remaining (Not Handed Over) -->
                  <div class="bg-white/20 rounded-lg p-2 backdrop-blur-sm">
                    <p class="text-xs opacity-90 mb-1.5">💰 Chưa bàn giao</p>
                    <div class="grid grid-cols-2 gap-2">
                      <div>
                        <p class="text-xs opacity-90">💵 Mặt</p>
                        <p class="text-sm font-bold text-green-300">{{ formatPrice(waiter.remainingCash) }}</p>
                      </div>
                      <div>
                        <p class="text-xs opacity-90">💳 CK</p>
                        <p class="text-sm font-bold text-blue-300">{{ formatPrice(waiter.remainingTransfer) }}</p>
                      </div>
                    </div>
                  </div>

                  <!-- Handed Over (Confirmed + Pending) -->
                  <div class="bg-white/20 rounded-lg p-2 backdrop-blur-sm">
                    <p class="text-xs opacity-90 mb-1.5">✅ Đã bàn giao</p>
                    <div class="grid grid-cols-2 gap-2">
                      <div>
                        <p class="text-xs opacity-90">💵 Mặt</p>
                        <p class="text-sm font-bold">{{ formatPrice(waiter.handedOverCash) }}</p>
                      </div>
                      <div>
                        <p class="text-xs opacity-90">💳 CK</p>
                        <p class="text-sm font-bold">{{ formatPrice(waiter.handedOverTransfer) }}</p>
                      </div>
                    </div>
                  </div>

                  <!-- Pending Handover Status -->
                  <div v-if="waiter.pendingCount > 0" class="bg-yellow-500/30 rounded-lg p-2 backdrop-blur-sm border border-yellow-300/50">
                    <div class="flex items-center justify-between">
                      <p class="text-xs font-medium">⏳ Chờ xác nhận</p>
                      <p class="text-xs font-bold">{{ waiter.pendingCount }} lần</p>
                    </div>
                  </div>

                  <!-- Discrepancy Warning -->
                  <div v-if="waiter.hasDiscrepancy" class="bg-red-500/30 rounded-lg p-2 backdrop-blur-sm border border-red-300/50">
                    <div class="flex items-center justify-between">
                      <p class="text-xs font-medium">⚠️ Chênh lệch</p>
                      <p class="text-xs font-bold">{{ formatPrice(waiter.totalDiscrepancy) }}</p>
                    </div>
                  </div>
                </div>

                <!-- Detailed Handovers (Expandable) -->
                <div v-if="expandedWaiters.includes(waiter.userId)" class="mt-3 pt-3 border-t border-white/30 space-y-2">
                  <p class="text-xs font-bold opacity-90 mb-2">📋 Lịch sử bàn giao</p>
                  <div 
                    v-for="handover in waiter.handovers" 
                    :key="handover.id"
                    class="bg-white/10 rounded-lg p-2 text-xs backdrop-blur-sm"
                  >
                    <div class="flex justify-between items-start mb-1.5">
                      <div>
                        <p class="font-medium">{{ formatTime(handover.created_at) }}</p>
                        <p class="text-xs opacity-75">{{ getShiftTypeText(handover.shift_type) }}</p>
                      </div>
                      <span :class="getHandoverStatusBadge(handover.status)">
                        {{ getHandoverStatusText(handover.status) }}
                      </span>
                    </div>
                    <div class="grid grid-cols-2 gap-1.5">
                      <div>
                        <span class="opacity-75">💵:</span>
                        <span class="font-bold ml-1">{{ formatPrice(handover.cash_amount) }}</span>
                      </div>
                      <div>
                        <span class="opacity-75">💳:</span>
                        <span class="font-bold ml-1">{{ formatPrice(handover.transfer_amount) }}</span>
                      </div>
                    </div>
                    <div v-if="handover.discrepancy_amount !== 0" class="mt-1.5 pt-1.5 border-t border-white/20">
                      <span class="opacity-75">⚠️ Chênh lệch:</span>
                      <span class="font-bold ml-1" :class="handover.discrepancy_amount > 0 ? 'text-green-300' : 'text-red-300'">
                        {{ formatPrice(handover.discrepancy_amount) }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Start Shift Modal -->
    <div v-if="showStartModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-2xl p-6 max-w-md w-full">
        <h3 class="text-xl font-bold text-gray-800 mb-4">Bắt đầu ca thu ngân</h3>
        
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Tiền đầu ca (VNĐ)
            <span class="text-gray-400 text-xs font-normal ml-1">(để trống nếu không có)</span>
          </label>
          <input
            v-model.number="startingFloat"
            type="number"
            step="1000"
            min="0"
            class="w-full border-2 border-gray-300 rounded-xl px-4 py-3 text-base focus:outline-none focus:border-yellow-500"
            :class="{ 'border-red-500': startingFloatError }"
            placeholder="0 (không bắt buộc)"
            @input="validateStartingFloat"
          />
          <p v-if="startingFloatError" class="text-sm text-red-600 mt-1">{{ startingFloatError }}</p>
        </div>

        <div class="flex gap-3">
          <button
            @click="showStartModal = false"
            class="flex-1 py-3 bg-gray-200 text-gray-800 rounded-xl font-bold active:scale-95 transition-transform"
          >
            Hủy
          </button>
          <button
            @click="startShift"
            :disabled="!canStartShift || startingLoading"
            class="flex-1 py-3 bg-yellow-500 text-white rounded-xl font-bold active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ startingLoading ? 'Đang xử lý...' : 'Bắt đầu' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCashierShiftStore } from '../stores/cashierShift'
import { useCashierStore } from '../stores/cashier'
import { useShiftStore } from '../stores/shift'

const router = useRouter()
const cashierShiftStore = useCashierShiftStore()
const cashierStore = useCashierStore()
const shiftStore = useShiftStore()

// State
const showStartModal = ref(false)
const startingFloat = ref(null) // Optional starting float, defaults to 0
const startingFloatError = ref(null)
const startingLoading = ref(false)
const loadingHandovers = ref(false)
const expandedWaiters = ref([]) // Array of waiter user IDs that are expanded
const openWaiterShifts = ref([]) // Open waiter shifts to calculate "not handed over yet"

// Computed
const loading = computed(() => cashierShiftStore.loading)
const error = computed(() => cashierShiftStore.error)
const currentShift = computed(() => cashierShiftStore.currentCashierShift)
const hasOpenShift = computed(() => cashierShiftStore.hasOpenCashierShift)
const todayHandovers = computed(() => cashierStore.todayHandovers)
const pendingHandovers = computed(() => cashierStore.pendingHandovers)

const canCloseShift = computed(() => {
  return currentShift.value && currentShift.value.status === 'OPEN'
})

const canStartShift = computed(() => {
  return !startingFloatError.value
})

// Confirmed handovers summary
const confirmedSummary = computed(() => {
  const summary = { cash: 0, transfer: 0, total: 0 }
  todayHandovers.value
    .filter(h => h.status === 'CONFIRMED')
    .forEach(handover => {
      summary.cash += handover.cash_amount || 0
      summary.transfer += handover.transfer_amount || 0
    })
  summary.total = summary.cash + summary.transfer
  return summary
})

const confirmedHandoversCount = computed(() => {
  return todayHandovers.value.filter(h => h.status === 'CONFIRMED').length
})

// Pending handovers summary
const pendingSummary = computed(() => {
  const summary = { cash: 0, transfer: 0, total: 0 }
  pendingHandovers.value.forEach(handover => {
    summary.cash += handover.cash_amount || 0
    summary.transfer += handover.transfer_amount || 0
  })
  summary.total = summary.cash + summary.transfer
  return summary
})

const pendingHandoversCount = computed(() => {
  return pendingHandovers.value.length
})

// Total discrepancy across all handovers
const totalDiscrepancy = computed(() => {
  let total = 0
  todayHandovers.value.forEach(handover => {
    total += handover.discrepancy_amount || 0
  })
  return total
})

// Already handed over summary (from open waiter shifts)
const handedOverSummary = computed(() => {
  const summary = { cash: 0, transfer: 0, total: 0 }
  openWaiterShifts.value.forEach(shift => {
    summary.cash += shift.handed_over_cash || 0
    summary.transfer += shift.handed_over_transfer || 0
  })
  summary.total = summary.cash + summary.transfer
  return summary
})

// Not handed over yet summary (from open waiter shifts)
const notHandedOverSummary = computed(() => {
  const summary = { cash: 0, transfer: 0, total: 0 }
  openWaiterShifts.value.forEach(shift => {
    summary.cash += shift.remaining_cash || 0
    summary.transfer += shift.remaining_transfer || 0
  })
  summary.total = summary.cash + summary.transfer
  return summary
})

const openWaiterShiftsCount = computed(() => {
  return openWaiterShifts.value.length
})

// Waiter financial breakdown (combines handovers + open shifts)
// Only show waiters with OPEN shifts
const waiterFinancialBreakdown = computed(() => {
  const breakdown = {}
  
  // Start with open waiter shifts (this is the primary filter)
  openWaiterShifts.value.forEach(shift => {
    const userId = shift.user_id
    const userName = shift.user_name
    
    breakdown[userId] = {
      userId,
      userName,
      shiftId: shift.id,
      shiftType: shift.type,
      handovers: [],
      totalCash: shift.current_cash || 0,
      totalTransfer: shift.transfer_revenue || 0,
      handedOverCash: shift.handed_over_cash || 0,
      handedOverTransfer: shift.handed_over_transfer || 0,
      remainingCash: shift.remaining_cash || 0,
      remainingTransfer: shift.remaining_transfer || 0,
      totalDiscrepancy: 0,
      hasDiscrepancy: false,
      pendingCount: 0
    }
  })
  
  // Add handover data only for waiters with open shifts
  todayHandovers.value.forEach(handover => {
    const userId = handover.waiter_id
    
    // Only add handover if this waiter has an open shift
    if (breakdown[userId]) {
      breakdown[userId].handovers.push(handover)
      breakdown[userId].totalDiscrepancy += handover.discrepancy_amount || 0
      
      if (handover.discrepancy_amount !== 0) {
        breakdown[userId].hasDiscrepancy = true
      }
      
      if (handover.status === 'PENDING') {
        breakdown[userId].pendingCount++
      }
    }
  })
  
  // Convert to array and sort by total amount descending
  return Object.values(breakdown).sort((a, b) => {
    const totalA = a.totalCash + a.totalTransfer
    const totalB = b.totalCash + b.totalTransfer
    return totalB - totalA
  })
})

// Group handovers by waiter
const handoversByWaiter = computed(() => {
  const grouped = {}
  
  todayHandovers.value.forEach(handover => {
    const userId = handover.waiter_id
    const userName = handover.waiter_name
    
    if (!grouped[userId]) {
      grouped[userId] = {
        userId,
        userName,
        handovers: [],
        totalCash: 0,
        totalTransfer: 0,
        totalDiscrepancy: 0,
        handoverCount: 0,
        hasDiscrepancy: false
      }
    }
    
    grouped[userId].handovers.push(handover)
    grouped[userId].totalCash += handover.cash_amount || 0
    grouped[userId].totalTransfer += handover.transfer_amount || 0
    grouped[userId].totalDiscrepancy += handover.discrepancy_amount || 0
    grouped[userId].handoverCount++
    
    if (handover.discrepancy_amount !== 0) {
      grouped[userId].hasDiscrepancy = true
    }
  })
  
  // Convert to array and sort by total amount descending
  return Object.values(grouped).sort((a, b) => {
    const totalA = a.totalCash + a.totalTransfer
    const totalB = b.totalCash + b.totalTransfer
    return totalB - totalA
  })
})

// Total summary across all waiters
const totalHandoverSummary = computed(() => {
  const summary = {
    cash: 0,
    transfer: 0,
    total: 0
  }
  
  handoversByWaiter.value.forEach(waiter => {
    summary.cash += waiter.totalCash
    summary.transfer += waiter.totalTransfer
  })
  
  summary.total = summary.cash + summary.transfer
  return summary
})

// Methods

// Methods
const validateStartingFloat = () => {
  startingFloatError.value = null

  const val = startingFloat.value
  if (val !== null && val !== '' && val < 0) {
    startingFloatError.value = 'Số tiền không được âm'
    return false
  }

  return true
}

const startShift = async () => {
  if (!validateStartingFloat()) return
  
  startingLoading.value = true
  
  try {
    await cashierShiftStore.startCashierShift(startingFloat.value || 0)
    showStartModal.value = false
    startingFloat.value = null // Reset
  } catch (error) {
    console.error('Failed to start cashier shift:', error)
  } finally {
    startingLoading.value = false
  }
}

const goToShiftClosure = () => {
  if (currentShift.value?.id) {
    router.push(`/cashier/shift-closure/${currentShift.value.id}`)
  }
}

const loadHandoverSummary = async () => {
  loadingHandovers.value = true
  try {
    await Promise.all([
      cashierStore.fetchTodayHandovers(),
      cashierStore.fetchPendingHandovers(),
      loadOpenWaiterShifts()
    ])
  } catch (error) {
    console.error('Load handover summary error:', error)
  } finally {
    loadingHandovers.value = false
  }
}

const loadOpenWaiterShifts = async () => {
  try {
    // Fetch all shifts using shift store
    await shiftStore.fetchAllShifts()
    
    // Filter for open waiter shifts
    openWaiterShifts.value = shiftStore.shifts.filter(s => 
      s.status === 'OPEN' && s.role_type === 'waiter'
    )
    
    console.log('=== Open Waiter Shifts Debug ===')
    console.log('Total open waiter shifts:', openWaiterShifts.value.length)
    openWaiterShifts.value.forEach(shift => {
      console.log(`Shift ${shift.id} (${shift.user_name}):`, {
        current_cash: shift.current_cash,
        transfer_revenue: shift.transfer_revenue,
        handed_over_cash: shift.handed_over_cash,
        handed_over_transfer: shift.handed_over_transfer,
        remaining_cash: shift.remaining_cash,
        remaining_transfer: shift.remaining_transfer
      })
    })
  } catch (error) {
    console.error('Load open waiter shifts error:', error)
  }
}

const toggleWaiterDetails = (userId) => {
  const index = expandedWaiters.value.indexOf(userId)
  if (index > -1) {
    expandedWaiters.value.splice(index, 1)
  } else {
    expandedWaiters.value.push(userId)
  }
}

const formatPrice = (amount) => {
  if (!amount && amount !== 0) return '0₫'
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND',
    maximumFractionDigits: 0
  }).format(amount)
}

const formatPriceShort = (amount) => {
  if (!amount && amount !== 0) return '0₫'
  if (amount >= 1000000) {
    return (amount / 1000000).toFixed(1) + 'tr'
  }
  if (amount >= 1000) {
    return (amount / 1000).toFixed(0) + 'k'
  }
  return amount + '₫'
}

const formatTime = (date) => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleTimeString('vi-VN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getStatusText = (status) => {
  const statusMap = {
    'OPEN': '🟢 Đang mở',
    'CLOSURE_INITIATED': '🟡 Đang đóng',
    'CLOSED': '🔴 Đã đóng'
  }
  return statusMap[status] || status
}

const getShiftTypeText = (type) => {
  const types = {
    'MORNING': '☀️ Sáng',
    'AFTERNOON': '🌤️ Chiều',
    'EVENING': '🌙 Tối'
  }
  return types[type] || type
}

const getHandoverStatusText = (status) => {
  const statusMap = {
    'PENDING': '⏳ Chờ',
    'CONFIRMED': '✅ OK',
    'REJECTED': '❌ Từ chối',
    'REQUIRES_APPROVAL': '⚠️ Cần duyệt'
  }
  return statusMap[status] || status
}

const getHandoverStatusBadge = (status) => {
  const badges = {
    'PENDING': 'px-1.5 py-0.5 text-xs rounded bg-yellow-100 text-yellow-700 font-medium',
    'CONFIRMED': 'px-1.5 py-0.5 text-xs rounded bg-green-100 text-green-700 font-medium',
    'REJECTED': 'px-1.5 py-0.5 text-xs rounded bg-red-100 text-red-700 font-medium',
    'REQUIRES_APPROVAL': 'px-1.5 py-0.5 text-xs rounded bg-orange-100 text-orange-700 font-medium'
  }
  return badges[status] || 'px-1.5 py-0.5 text-xs rounded bg-gray-100 text-gray-700 font-medium'
}

// Lifecycle
onMounted(async () => {
  try {
    await cashierShiftStore.fetchCurrentCashierShift()
    // Load handover summary if shift is open
    if (hasOpenShift.value) {
      await loadHandoverSummary()
    }
  } catch (error) {
    // Silently ignore 404 errors (no shift found is normal)
    if (error.response?.status !== 404) {
      console.error('Failed to fetch current cashier shift:', error)
    }
  }
})
</script>

<style scoped>
.active\:scale-95:active {
  transform: scale(0.95);
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
