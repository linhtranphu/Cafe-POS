<!--
  CashierShiftView Component
  
  This component provides a redesigned shift management interface specifically for Cashier role.
  
  Features:
  - Tab navigation to switch between Waiter and Barista shifts
  - Date filter to view shifts from specific dates
  - Shift selector dropdown to choose a shift
  - Shift summary cards with financial/order statistics
  - Payment list with actions (override, discrepancy, lock) for waiter shifts only
  - Close shift functionality for open shifts
  
  Note: This component is only rendered when the user has Cashier role.
  Waiter and Barista users continue to use the original ShiftView interface.
-->
<template>
  <div>
    <!-- Tab Navigation -->
    <div class="mb-4">
      <div class="flex gap-2">
        <button 
          @click="onTabChange('waiter')"
          :class="activeTab === 'waiter' ? 'bg-blue-500 text-white shadow-lg' : 'bg-gray-100 text-gray-700 active:bg-gray-200'"
          class="flex-1 px-4 py-3.5 rounded-xl font-semibold transition-all touch-manipulation active:scale-98">
          <div class="flex items-center justify-center gap-2">
            <span>🍽️</span>
            <span>Ca phục vụ</span>
          </div>
        </button>
        <button 
          @click="onTabChange('barista')"
          :class="activeTab === 'barista' ? 'bg-blue-500 text-white shadow-lg' : 'bg-gray-100 text-gray-700 active:bg-gray-200'"
          class="flex-1 px-4 py-3.5 rounded-xl font-semibold transition-all touch-manipulation active:scale-98">
          <div class="flex items-center justify-center gap-2">
            <span>🍹</span>
            <span>Ca pha chế</span>
          </div>
        </button>
      </div>
    </div>

    <!-- Unified Shift Selector -->
    <div class="bg-white rounded-2xl p-4 shadow-lg border border-gray-100 mb-4">
      <div class="flex items-center justify-between mb-3">
        <label class="block text-sm font-bold text-gray-900 flex items-center gap-2">
          <span class="text-lg">📋</span>
          <span>Chọn ca làm việc</span>
        </label>
        <span v-if="allShifts.length > 0" class="text-xs font-semibold text-gray-500 bg-gray-100 px-2.5 py-1 rounded-full">
          {{ allShifts.length }} ca
        </span>
      </div>

      <!-- Quick Filter: Open/Closed -->
      <div v-if="allShifts.length > 0" class="flex gap-2 mb-3">
        <button
          @click="statusFilter = 'all'"
          :class="[
            'flex-1 py-2.5 px-3 rounded-lg text-xs font-semibold transition-all touch-manipulation',
            statusFilter === 'all' ? 'bg-blue-500 text-white shadow-md' : 'bg-gray-100 text-gray-600 active:bg-gray-200'
          ]"
        >
          Tất cả ({{ allShifts.length }})
        </button>
        <button
          @click="statusFilter = 'open'"
          :class="[
            'flex-1 py-2.5 px-3 rounded-lg text-xs font-semibold transition-all touch-manipulation',
            statusFilter === 'open' ? 'bg-green-500 text-white shadow-md' : 'bg-gray-100 text-gray-600 active:bg-gray-200'
          ]"
        >
          🟢 Mở ({{ allOpenShiftsCount }})
        </button>
        <button
          @click="statusFilter = 'closed'"
          :class="[
            'flex-1 py-2.5 px-3 rounded-lg text-xs font-semibold transition-all touch-manipulation',
            statusFilter === 'closed' ? 'bg-gray-500 text-white shadow-md' : 'bg-gray-100 text-gray-600 active:bg-gray-200'
          ]"
        >
          🔴 Đóng ({{ allClosedShiftsCount }})
        </button>
      </div>
      
      <!-- No shifts message -->
      <div v-if="allShifts.length === 0" class="text-center py-8 bg-gray-50 rounded-xl border-2 border-dashed border-gray-200">
        <div class="text-4xl mb-2">📭</div>
        <p class="text-sm font-medium text-gray-600 mb-1">Chưa có ca làm việc</p>
        <p class="text-xs text-gray-500">Ca làm việc sẽ hiển thị ở đây</p>
      </div>

      <!-- No shifts after filter -->
      <div v-else-if="displayedAllShifts.length === 0" class="text-center py-8 bg-gray-50 rounded-xl border-2 border-dashed border-gray-200">
        <div class="text-4xl mb-2">🔍</div>
        <p class="text-sm font-medium text-gray-600">Không có ca {{ statusFilter === 'open' ? 'đang mở' : 'đã đóng' }}</p>
      </div>

      <!-- Unified Dropdown: Date + Shift -->
      <select 
        v-else
        v-model="selectedShiftId"
        @change="onShiftSelect"
        class="w-full border-2 border-gray-200 rounded-xl px-4 py-3.5 text-base font-medium focus:outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-500 touch-manipulation transition-all bg-white"
      >
        <option value="" disabled>-- Chọn ca để xem chi tiết --</option>
        <optgroup 
          v-for="(shiftsInDate, date) in groupedShifts" 
          :key="date" 
          :label="formatDateLabel(date)"
        >
          <option 
            v-for="shift in shiftsInDate" 
            :key="shift.id" 
            :value="shift.id"
          >
            {{ getShiftTypeIcon(shift.type) }} {{ shift.user_name }} • {{ getShiftTypeTextShort(shift.type) }} • {{ formatTime(shift.started_at) }} {{ shift.status === SHIFT_STATUS.OPEN ? '🟢' : '🔴' }}
          </option>
        </optgroup>
      </select>
    </div>

    <!-- Loading State -->
    <div v-if="loadingShiftDetails" class="text-center py-10">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
      <p class="text-gray-500 mt-2">Đang tải...</p>
    </div>

    <!-- Shift Summary Card - Waiter -->
    <div v-else-if="selectedShift && isWaiterShift && shiftStatus" 
      class="bg-gradient-to-r from-blue-500 to-purple-500 text-white rounded-2xl p-5 shadow-xl border border-blue-400 mb-4">
      <div class="flex justify-between items-start mb-4">
        <div>
          <h3 class="text-2xl font-bold">Ca đang xem</h3>
          <p class="text-blue-100 font-medium">{{ getShiftTypeText(selectedShift.type) }}</p>
          <p class="text-sm text-blue-100 mt-1">{{ selectedShift.user_name }}</p>
        </div>
        <span :class="selectedShift.status === SHIFT_STATUS.OPEN ? 'bg-white text-blue-600' : 'bg-white/30 text-white'" 
          class="px-4 py-2 rounded-full font-bold text-sm shadow-md">
          {{ selectedShift.status === SHIFT_STATUS.OPEN ? 'ĐANG MỞ' : 'ĐÃ ĐÓNG' }}
        </span>
      </div>
      
      <div class="grid grid-cols-2 gap-3 mb-4">
        <div class="bg-white bg-opacity-20 rounded-xl p-3">
          <p class="text-sm text-blue-100">Bắt đầu</p>
          <p class="font-bold">{{ formatTime(selectedShift.started_at) }}</p>
        </div>
        <div class="bg-white bg-opacity-20 rounded-xl p-3">
          <p class="text-sm text-blue-100">Tiền đầu ca</p>
          <p class="font-bold">{{ formatPrice(selectedShift.start_cash) }}</p>
        </div>
      </div>

      <!-- Cash Status for Waiter -->
      <div class="space-y-3 mb-4">
        <!-- Revenue Section -->
        <div class="bg-white bg-opacity-20 rounded-xl p-3">
          <p class="text-xs text-blue-100 mb-2">📊 Doanh thu</p>
          <div class="grid grid-cols-2 gap-2">
            <div>
              <p class="text-xs text-blue-100">💵 Tiền mặt</p>
              <p class="font-bold">{{ formatPrice(selectedShift.current_cash || 0) }}</p>
            </div>
            <div>
              <p class="text-xs text-blue-100">💳 Tiền CK</p>
              <p class="font-bold">{{ formatPrice(selectedShift.transfer_revenue || 0) }}</p>
            </div>
          </div>
        </div>
        
        <!-- Remaining Section -->
        <div class="bg-white bg-opacity-20 rounded-xl p-3">
          <p class="text-xs text-blue-100 mb-2">💰 Còn lại (chưa bàn giao)</p>
          <div class="grid grid-cols-2 gap-2">
            <div>
              <p class="text-xs text-blue-100">💵 Tiền mặt</p>
              <p class="font-bold text-green-300">{{ formatPrice(selectedShift.remaining_cash || 0) }}</p>
            </div>
            <div>
              <p class="text-xs text-blue-100">💳 Tiền CK</p>
              <p class="font-bold text-blue-300">{{ formatPrice(selectedShift.remaining_transfer || 0) }}</p>
            </div>
          </div>
        </div>
        
        <!-- Handed Over Section -->
        <div class="bg-white bg-opacity-20 rounded-xl p-3">
          <p class="text-xs text-blue-100 mb-2">✅ Đã bàn giao</p>
          <div class="grid grid-cols-2 gap-2">
            <div>
              <p class="text-xs text-blue-100">💵 Tiền mặt</p>
              <p class="font-bold">{{ formatPrice(selectedShift.handed_over_cash || 0) }}</p>
            </div>
            <div>
              <p class="text-xs text-blue-100">💳 Tiền CK</p>
              <p class="font-bold">{{ formatPrice(selectedShift.handed_over_transfer || 0) }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Stats Grid -->
      <div class="grid grid-cols-2 gap-3 mb-4">
        <div class="bg-white bg-opacity-20 rounded-xl p-3">
          <div class="text-xl sm:text-2xl font-bold">{{ shiftStatus.total_orders }}</div>
          <div class="text-xs opacity-90">Tổng đơn</div>
        </div>
        <div class="bg-white bg-opacity-20 rounded-xl p-3">
          <div class="text-base sm:text-lg font-bold break-words">{{ formatPrice(shiftStatus.total_revenue) }}</div>
          <div class="text-xs opacity-90">Tổng doanh thu</div>
        </div>
      </div>

      <!-- Close Shift Button -->
      <button v-if="isShiftOpen" 
        @click="showCloseShiftModal"
        class="w-full bg-red-500 active:bg-red-600 text-white px-4 py-3.5 rounded-xl font-bold shadow-lg touch-manipulation active:scale-98 transition-all">
        🔒 Chốt ca
      </button>
    </div>

    <!-- Shift Summary Card - Barista -->
    <div v-else-if="selectedShift && !isWaiterShift && shiftStatus" 
      class="bg-gradient-to-r from-purple-500 to-indigo-500 text-white rounded-2xl p-6 shadow-lg mb-4">
      <h2 class="text-lg font-bold mb-2">📊 Tổng quan ca pha chế</h2>
      <div class="flex items-center justify-between mb-4">
        <div>
          <p class="text-2xl font-bold">{{ selectedShift.user_name }}</p>
          <p class="text-sm opacity-90">{{ getShiftTypeText(selectedShift.type) }}</p>
          <p class="text-xs opacity-75 mt-1">
            {{ formatTime(selectedShift.started_at) }}
            <span v-if="selectedShift.ended_at"> - {{ formatTime(selectedShift.ended_at) }}</span>
          </p>
        </div>
        <span class="bg-white/30 backdrop-blur-sm px-3 py-1 rounded-full text-xs font-medium">
          {{ selectedShift.status === SHIFT_STATUS.OPEN ? '🟢 Đang mở' : '🔴 Đã đóng' }}
        </span>
      </div>

      <!-- Stats Grid -->
      <div class="grid grid-cols-2 gap-3 mb-4">
        <div class="bg-white/20 rounded-xl p-3 backdrop-blur-sm">
          <div class="text-xl sm:text-2xl font-bold">{{ shiftStatus.total_orders }}</div>
          <div class="text-xs opacity-90">Tổng đơn</div>
        </div>
        <div class="bg-white/20 rounded-xl p-3 backdrop-blur-sm">
          <div class="text-xl sm:text-2xl font-bold">{{ shiftStatus.in_progress_count || 0 }}</div>
          <div class="text-xs opacity-90">Đang pha chế</div>
        </div>
      </div>

      <!-- Close Shift Button -->
      <button v-if="isShiftOpen" 
        @click="showCloseShiftModal"
        class="w-full bg-red-500 hover:bg-red-600 text-white px-4 py-3 rounded-xl font-bold active:scale-95 transition-transform">
        🔒 Chốt ca
      </button>
    </div>

    <!-- Payment List (Waiter shifts only) -->
    <div v-if="selectedShift && isWaiterShift" class="bg-white rounded-2xl p-4 shadow-sm mb-4">
      <div class="flex items-center justify-between mb-3">
        <h2 class="text-lg font-bold text-gray-800">💳 Danh sách thanh toán</h2>
        <span class="text-sm text-gray-600">{{ payments.length }} giao dịch</span>
      </div>

      <div v-if="loadingPayments" class="text-center py-10">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
      </div>

      <div v-else-if="payments.length === 0" class="text-center py-10">
        <div class="text-5xl mb-3">📭</div>
        <p class="text-gray-500">Chưa có thanh toán nào</p>
      </div>

      <div v-else class="space-y-3">
        <div
          v-for="payment in payments"
          :key="payment.order_id"
          class="bg-white border rounded-xl p-4 shadow-sm"
        >
          <!-- Header -->
          <div class="flex justify-between items-start mb-3">
            <div>
              <h3 class="font-bold text-gray-800">{{ payment.customer_name || 'Khách lẻ' }}</h3>
              <p class="text-xs text-gray-500">{{ formatDateTime(payment.paid_at) }}</p>
            </div>
            <div class="text-right">
              <div class="text-lg font-bold text-green-600">{{ formatPrice(payment.amount) }}</div>
              <span :class="getPaymentMethodBadge(payment.payment_method)">
                {{ getPaymentMethodText(payment.payment_method) }}
              </span>
            </div>
          </div>

          <!-- Status -->
          <div class="mb-3">
            <span :class="getStatusBadge(payment.status)">
              {{ getStatusText(payment.status) }}
            </span>
          </div>

          <!-- Actions -->
          <div class="flex gap-2">
            <button
              @click="showOverrideModal(payment)"
              class="flex-1 py-2.5 bg-orange-50 text-orange-600 rounded-lg text-xs sm:text-sm font-medium active:scale-95 transition-transform border border-orange-200 min-h-[44px]"
            >
              ✏️ Điều chỉnh
            </button>
            <button
              @click="showDiscrepancyModal(payment)"
              class="flex-1 py-2.5 bg-yellow-50 text-yellow-600 rounded-lg text-xs sm:text-sm font-medium active:scale-95 transition-transform border border-yellow-200 min-h-[44px]"
            >
              ⚠️ Báo lỗi
            </button>
            <button
              @click="lockOrder(payment.order_id)"
              class="flex-1 py-2.5 bg-red-50 text-red-600 rounded-lg text-xs sm:text-sm font-medium active:scale-95 transition-transform border border-red-200 min-h-[44px]"
            >
              🔒 Khóa
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modals -->
    <OverridePaymentModal
      :show="showOverride"
      :payment="selectedPayment"
      @close="showOverride = false"
      @confirm="handleOverridePayment"
    />

    <DiscrepancyModal
      :show="showDiscrepancy"
      :payment="selectedPayment"
      @close="showDiscrepancy = false"
      @confirm="handleReportDiscrepancy"
    />

    <!-- Close Shift Modal -->
    <transition name="slide-up">
      <div v-if="showCloseForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6">
          <h3 class="text-xl font-bold mb-4">Chốt ca</h3>
          <p class="text-sm text-gray-600 mb-4">Chốt ca sẽ khóa tất cả orders trong ca này</p>
          <form @submit.prevent="closeShift" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Tiền cuối ca (VNĐ) *</label>
              <input v-model.number="closeForm.end_cash" type="number" min="0" step="1000" required 
                class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-red-500">
            </div>
            <div class="flex gap-2">
              <button type="button" @click="showCloseForm = false" 
                class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
                Hủy
              </button>
              <button type="submit" 
                class="flex-1 bg-red-500 hover:bg-red-600 text-white px-4 py-3 rounded-xl font-medium">
                Chốt ca
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useShiftStore } from '../stores/shift'
import { useCashierStore } from '../stores/cashier'
import OverridePaymentModal from './OverridePaymentModal.vue'
import DiscrepancyModal from './DiscrepancyModal.vue'
import { SHIFT_STATUS } from '../constants/shift'
import { ORDER_STATUS, PAYMENT_METHOD } from '../constants/order'

const shiftStore = useShiftStore()
const cashierStore = useCashierStore()

// Tab and selection state
const activeTab = ref('waiter') // 'waiter' | 'barista'
const selectedShiftId = ref('')
const loadingShiftDetails = ref(false)
const loadingPayments = ref(false)
const statusFilter = ref('all') // 'all' | 'open' | 'closed'

// Modal state for payment actions
const showOverride = ref(false)
const showDiscrepancy = ref(false)
const showCloseForm = ref(false)
const selectedPayment = ref(null)

const closeForm = ref({
  end_cash: 0
})

// Computed properties
const shifts = computed(() => shiftStore.shifts)
const shiftStatus = computed(() => cashierStore.shiftStatus)
const payments = computed(() => cashierStore.payments)

// Get all shifts for current tab (no date filter)
const allShifts = computed(() => {
  return shifts.value.filter(shift => {
    // Filter by role type based on active tab
    if (activeTab.value === 'waiter') {
      return shift.role_type === 'waiter'
    } else {
      return shift.role_type === 'barista'
    }
  }).sort((a, b) => {
    // Sort by date descending (newest first)
    return new Date(b.started_at) - new Date(a.started_at)
  })
})

// Apply status filter
const displayedAllShifts = computed(() => {
  if (statusFilter.value === 'all') return allShifts.value
  if (statusFilter.value === 'open') {
    return allShifts.value.filter(s => s.status === SHIFT_STATUS.OPEN)
  }
  if (statusFilter.value === 'closed') {
    return allShifts.value.filter(s => s.status === SHIFT_STATUS.CLOSED)
  }
  return allShifts.value
})

// Group shifts by date for optgroup
const groupedShifts = computed(() => {
  const groups = {}
  displayedAllShifts.value.forEach(shift => {
    const date = new Date(shift.started_at).toISOString().split('T')[0]
    if (!groups[date]) {
      groups[date] = []
    }
    groups[date].push(shift)
  })
  return groups
})

// Count shifts by status
const allOpenShiftsCount = computed(() => {
  return allShifts.value.filter(s => s.status === SHIFT_STATUS.OPEN).length
})

const allClosedShiftsCount = computed(() => {
  return allShifts.value.filter(s => s.status === SHIFT_STATUS.CLOSED).length
})

const selectedShift = computed(() => {
  if (!selectedShiftId.value) return null
  return shifts.value.find(s => s.id === selectedShiftId.value)
})

const isShiftOpen = computed(() => {
  return selectedShift.value?.status === SHIFT_STATUS.OPEN
})

const isWaiterShift = computed(() => {
  return selectedShift.value?.role_type === 'waiter'
})

// Event handlers
const onTabChange = (tab) => {
  activeTab.value = tab
  selectedShiftId.value = '' // Clear selection when switching tabs
  statusFilter.value = 'all' // Reset filter when switching tabs
}

// Select shift (replaces onShiftSelect for button-based selection)
const selectShift = async (shiftId) => {
  selectedShiftId.value = shiftId
  await onShiftSelect()
}

// Load shift details and payments when a shift is selected
const onShiftSelect = async () => {
  if (!selectedShiftId.value) return
  
  loadingShiftDetails.value = true
  try {
    const statusData = await cashierStore.getShiftStatus(selectedShiftId.value)
    console.log('Shift status data:', statusData)
    console.log('Cash revenue:', statusData.cash_revenue)
    console.log('Transfer revenue:', statusData.transfer_revenue)
    console.log('QR revenue:', statusData.qr_revenue)
    
    // Only load payments for waiter shifts
    if (isWaiterShift.value) {
      loadingPayments.value = true
      await cashierStore.getPaymentsByShift(selectedShiftId.value)
      loadingPayments.value = false
    }
  } catch (error) {
    console.error('Error loading shift details:', error)
    alert('Lỗi khi tải thông tin ca: ' + error.message)
  } finally {
    loadingShiftDetails.value = false
  }
}

const showCloseShiftModal = () => {
  showCloseForm.value = true
}

// Close shift handler
const closeShift = async () => {
  try {
    await shiftStore.closeShift(selectedShift.value.id, closeForm.value.end_cash)
    showCloseForm.value = false
    closeForm.value = { end_cash: 0 }
    selectedShiftId.value = ''
    await shiftStore.fetchAllShifts()
    alert('✅ Đã chốt ca thành công')
  } catch (error) {
    console.error('Close shift error:', error)
    alert('❌ Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

// Payment action handlers
const showOverrideModal = (payment) => {
  selectedPayment.value = payment
  showOverride.value = true
}

const showDiscrepancyModal = (payment) => {
  selectedPayment.value = payment
  showDiscrepancy.value = true
}

const handleOverridePayment = async (reason) => {
  try {
    await cashierStore.overridePayment(selectedPayment.value.order_id, reason)
    showOverride.value = false
    await cashierStore.getPaymentsByShift(selectedShiftId.value)
    alert('✅ Đã điều chỉnh thanh toán')
  } catch (error) {
    console.error('Override failed:', error)
    alert('❌ Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

const handleReportDiscrepancy = async (data) => {
  try {
    await cashierStore.reportDiscrepancy({
      order_id: selectedPayment.value.order_id,
      ...data
    })
    showDiscrepancy.value = false
    alert('✅ Đã báo lỗi')
  } catch (error) {
    console.error('Report discrepancy failed:', error)
    alert('❌ Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

const lockOrder = async (orderId) => {
  if (confirm('Bạn có chắc muốn khóa order này? Không thể hoàn tác!')) {
    try {
      await cashierStore.lockOrder(orderId)
      await cashierStore.getPaymentsByShift(selectedShiftId.value)
      alert('✅ Đã khóa order')
    } catch (error) {
      console.error('Lock order failed:', error)
      alert('❌ Lỗi: ' + (error.response?.data?.error || error.message))
    }
  }
}

// Helper functions
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

const formatDateLabel = (dateStr) => {
  const date = new Date(dateStr)
  const today = new Date()
  const yesterday = new Date(today)
  yesterday.setDate(yesterday.getDate() - 1)
  
  const dateOnly = dateStr
  const todayStr = today.toISOString().split('T')[0]
  const yesterdayStr = yesterday.toISOString().split('T')[0]
  
  if (dateOnly === todayStr) {
    return '📅 Hôm nay - ' + date.toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit' })
  } else if (dateOnly === yesterdayStr) {
    return '📅 Hôm qua - ' + date.toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit' })
  } else {
    return '📅 ' + date.toLocaleDateString('vi-VN', { 
      weekday: 'short',
      day: '2-digit', 
      month: '2-digit',
      year: 'numeric'
    })
  }
}

const formatTime = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleTimeString('vi-VN', { 
    hour: '2-digit', 
    minute: '2-digit' 
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

const getShiftTypeText = (type) => {
  const types = {
    MORNING: '☀️ Ca sáng',
    AFTERNOON: '🌤️ Ca chiều',
    EVENING: '🌙 Ca tối'
  }
  return types[type] || type
}

const getShiftTypeTextShort = (type) => {
  const types = {
    MORNING: 'Sáng',
    AFTERNOON: 'Chiều',
    EVENING: 'Tối'
  }
  return types[type] || type
}

const getShiftTypeIcon = (type) => {
  const icons = {
    MORNING: '☀️',
    AFTERNOON: '🌤️',
    EVENING: '🌙'
  }
  return icons[type] || '⏰'
}

const getPaymentMethodBadge = (method) => {
  const badges = {
    CASH: 'px-2 py-1 text-xs rounded-full bg-green-100 text-green-700 font-medium',
    TRANSFER: 'px-2 py-1 text-xs rounded-full bg-blue-100 text-blue-700 font-medium',
    QR: 'px-2 py-1 text-xs rounded-full bg-purple-100 text-purple-700 font-medium'
  }
  return badges[method] || 'px-2 py-1 text-xs rounded-full bg-gray-100 text-gray-700 font-medium'
}

const getPaymentMethodText = (method) => {
  const texts = {
    CASH: '💵 Tiền mặt',
    TRANSFER: '💳 CK',
    QR: '📱 QR'
  }
  return texts[method] || method
}

const getStatusBadge = (status) => {
  const badges = {
    PAID: 'px-3 py-1 text-xs rounded-full bg-green-100 text-green-700 font-medium inline-block',
    QUEUED: 'px-3 py-1 text-xs rounded-full bg-yellow-100 text-yellow-700 font-medium inline-block',
    IN_PROGRESS: 'px-3 py-1 text-xs rounded-full bg-blue-100 text-blue-700 font-medium inline-block',
    READY: 'px-3 py-1 text-xs rounded-full bg-purple-100 text-purple-700 font-medium inline-block',
    SERVED: 'px-3 py-1 text-xs rounded-full bg-green-100 text-green-700 font-medium inline-block',
    LOCKED: 'px-3 py-1 text-xs rounded-full bg-red-100 text-red-700 font-medium inline-block'
  }
  return badges[status] || 'px-3 py-1 text-xs rounded-full bg-gray-100 text-gray-700 font-medium inline-block'
}

const getStatusText = (status) => {
  const texts = {
    PAID: '✓ Đã thu',
    QUEUED: '⏳ Chờ pha',
    IN_PROGRESS: '🍹 Đang pha',
    READY: '✅ Sẵn sàng',
    SERVED: '🎯 Hoàn tất',
    LOCKED: '🔒 Đã khóa'
  }
  return texts[status] || status
}

// Lifecycle
onMounted(async () => {
  await shiftStore.fetchAllShifts()
})
</script>

<style scoped>
.touch-manipulation {
  touch-action: manipulation;
  -webkit-tap-highlight-color: transparent;
}

.active\:scale-95:active {
  transform: scale(0.95);
}

.active\:scale-98:active {
  transform: scale(0.98);
}

.scrollbar-thin {
  scrollbar-width: thin;
}

.scrollbar-thin::-webkit-scrollbar {
  width: 6px;
}

.scrollbar-thin::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 10px;
}

.scrollbar-thin::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 10px;
}

.scrollbar-thin::-webkit-scrollbar-thumb:hover {
  background: #555;
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
