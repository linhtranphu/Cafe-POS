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

      <!-- Quick Access: Handover Management -->
      <div class="bg-gradient-to-r from-yellow-500 to-orange-500 text-white rounded-2xl p-6 shadow-lg mb-4">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-xl font-bold mb-2">💰 Quản lý bàn giao</h3>
            <p class="text-sm opacity-90">Xác nhận bàn giao từ phục vụ</p>
          </div>
          <button 
            @click="$router.push('/cashier/handovers')"
            class="bg-white text-orange-600 px-6 py-3 rounded-xl font-bold shadow-lg hover:shadow-xl active:scale-95 transition-all">
            Xem ngay →
          </button>
        </div>
      </div>

      <!-- Cashier Shift Manager -->
      <CashierShiftManager />

      <!-- Shift Selector - Waiter Shifts for Payment Monitoring -->
      <div class="bg-white rounded-2xl p-4 shadow-sm mb-4">
        <label class="block text-sm font-medium text-gray-700 mb-2">📅 Chọn ca phục vụ để xem thanh toán</label>
        <p class="text-xs text-gray-500 mb-2">💡 Chọn ca của waiter để xem các thanh toán trong ca đó</p>
        
        <!-- Debug info -->
        <p v-if="waiterShifts.length === 0" class="text-xs text-red-500 mb-2">
          ⚠️ Không có ca phục vụ nào. Waiter cần mở ca trước.
        </p>
        
        <select 
          v-model="selectedShift" 
          @change="loadPayments" 
          class="w-full border-2 border-gray-300 rounded-xl px-4 py-3 text-base focus:outline-none focus:border-yellow-500"
        >
          <option value="">-- Chọn ca phục vụ --</option>
          <option v-for="shift in waiterShifts" :key="shift.id" :value="shift.id">
            {{ formatDate(shift.started_at) }} - 
            {{ shift.user_name }} 
            ({{ shift.role_type }}) - 
            {{ getStatusText(shift.status) }}
          </option>
        </select>
      </div>

      <!-- Shift Status Card - Shows stats of selected waiter shift -->
      <div v-if="shiftStatus" class="bg-gradient-to-r from-yellow-500 to-orange-500 text-white rounded-2xl p-6 shadow-lg mb-4">
        <h2 class="text-lg font-bold mb-4">📊 Tổng quan ca phục vụ</h2>
        <p class="text-xs opacity-90 mb-3">Thống kê của ca đang xem</p>
        <div class="grid grid-cols-2 gap-3">
          <div class="bg-white/20 rounded-xl p-3 backdrop-blur-sm">
            <div class="text-2xl font-bold">{{ shiftStatus.total_orders }}</div>
            <div class="text-xs opacity-90">Tổng đơn</div>
          </div>
          <div class="bg-white/20 rounded-xl p-3 backdrop-blur-sm">
            <div class="text-lg font-bold">{{ formatPrice(shiftStatus.total_revenue) }}</div>
            <div class="text-xs opacity-90">Tổng doanh thu</div>
          </div>
          <div class="bg-white/20 rounded-xl p-3 backdrop-blur-sm">
            <div class="text-lg font-bold">{{ formatPrice(shiftStatus.cash_revenue) }}</div>
            <div class="text-xs opacity-90">💵 Tiền mặt</div>
          </div>
          <div class="bg-white/20 rounded-xl p-3 backdrop-blur-sm">
            <div class="text-lg font-bold">{{ formatPrice(shiftStatus.transfer_revenue + shiftStatus.qr_revenue) }}</div>
            <div class="text-xs opacity-90">💳 Chuyển khoản</div>
          </div>
        </div>
      </div>

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
      <div v-if="false" class="bg-white rounded-2xl p-4 shadow-sm mb-4">
        <h2 class="text-lg font-bold text-gray-800 mb-4">💰 Đối soát tiền mặt</h2>
        
        <div v-if="!hasReconciliation" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Tiền mặt thực tế</label>
            <input
              v-model.number="reconciliationForm.actualCash"
              type="number"
              step="1000"
              class="w-full border-2 border-gray-300 rounded-xl px-4 py-3 text-base focus:outline-none focus:border-yellow-500"
              placeholder="Nhập số tiền thực tế"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Ghi chú (tùy chọn)</label>
            <textarea
              v-model="reconciliationForm.notes"
              rows="2"
              class="w-full border-2 border-gray-300 rounded-xl px-4 py-3 text-base focus:outline-none focus:border-yellow-500"
              placeholder="Ghi chú về đối soát..."
            ></textarea>
          </div>
          <button
            @click="performReconciliation"
            :disabled="!reconciliationForm.actualCash"
            class="w-full py-4 bg-green-500 text-white rounded-xl font-bold text-base active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed"
          >
            ✓ Xác nhận đối soát
          </button>
        </div>

        <div v-else class="bg-gray-50 rounded-xl p-4">
          <div class="grid grid-cols-1 gap-3 mb-3">
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Tiền mặt dự kiến:</span>
              <span class="font-bold text-gray-800">{{ formatPrice(reconciliation.expected_cash) }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Tiền mặt thực tế:</span>
              <span class="font-bold text-gray-800">{{ formatPrice(reconciliation.actual_cash) }}</span>
            </div>
            <div class="flex justify-between items-center pt-2 border-t-2 border-gray-200">
              <span class="text-sm font-medium text-gray-700">Chênh lệch:</span>
              <span :class="getDifferenceClass(reconciliation.difference)" class="font-bold text-lg">
                {{ formatPrice(reconciliation.difference) }}
              </span>
            </div>
          </div>
          <div v-if="reconciliation.notes" class="text-sm text-gray-600 bg-white rounded-lg p-3">
            <span class="font-medium">Ghi chú:</span> {{ reconciliation.notes }}
          </div>
          <div class="mt-3 flex items-center gap-2 text-sm text-green-600">
            <span class="text-xl">✓</span>
            <span class="font-medium">Đã đối soát</span>
          </div>
        </div>
      </div>

      <!-- Payment List -->
      <div class="mb-4">
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-lg font-bold text-gray-800">💳 Danh sách thanh toán</h2>
          <span class="text-sm text-gray-600">{{ payments.length }} giao dịch</span>
        </div>

        <div v-if="payments.length === 0" class="text-center py-12 bg-white rounded-2xl">
          <div class="text-5xl mb-3">📭</div>
          <p class="text-gray-500">Chưa có thanh toán nào</p>
          <p class="text-sm text-gray-400 mt-1">Chọn ca làm việc để xem</p>
        </div>

        <div v-else class="space-y-3">
          <div
            v-for="payment in payments"
            :key="payment.order_id"
            class="bg-white rounded-xl p-4 shadow-sm active:scale-98 transition-transform"
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
                class="flex-1 py-2 bg-orange-50 text-orange-600 rounded-lg text-sm font-medium active:scale-95 transition-transform border border-orange-200"
              >
                ✏️ Điều chỉnh
              </button>
              <button
                @click="showDiscrepancyModal(payment)"
                class="flex-1 py-2 bg-yellow-50 text-yellow-600 rounded-lg text-sm font-medium active:scale-95 transition-transform border border-yellow-200"
              >
                ⚠️ Báo lỗi
              </button>
              <button
                @click="lockOrder(payment.order_id)"
                class="px-4 py-2 bg-red-50 text-red-600 rounded-lg text-sm font-medium active:scale-95 transition-transform border border-red-200"
              >
                🔒
              </button>
            </div>
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
import CashierShiftManager from '../components/CashierShiftManager.vue'
import OverridePaymentModal from '../components/OverridePaymentModal.vue'
import DiscrepancyModal from '../components/DiscrepancyModal.vue'
import { SHIFT_STATUS, CASHIER_SHIFT_STATUS, SHIFT_TYPE } from '../constants/shift'
import { ORDER_STATUS, PAYMENT_METHOD } from '../constants/order'

const cashierStore = useCashierStore()
const cashierShiftStore = useCashierShiftStore()
const shiftStore = useShiftStore()

const selectedShift = ref('')
const showOverride = ref(false)
const showDiscrepancy = ref(false)
const showDiscrepancyList = ref(false)
const selectedPayment = ref(null)
const reconciliationForm = ref({
  actualCash: null,
  notes: ''
})

// Computed
const shiftStatus = computed(() => cashierStore.shiftStatus)
const cashierShifts = computed(() => cashierShiftStore.cashierShifts)
const waiterShifts = computed(() => shiftStore.shifts) // All shifts including waiter shifts
const payments = computed(() => cashierStore.payments)
const pendingDiscrepancies = computed(() => cashierStore.pendingDiscrepancies)
const reconciliation = computed(() => cashierStore.reconciliation)
const hasReconciliation = computed(() => cashierStore.hasReconciliation)
const loading = computed(() => cashierStore.loading)
const error = computed(() => cashierStore.error)

// Methods - Define refreshData BEFORE using it in usePullToRefresh
const refreshData = async () => {
  if (selectedShift.value) {
    await Promise.all([
      cashierStore.getShiftStatus(selectedShift.value),
      cashierStore.getPaymentsByShift(selectedShift.value)
    ])
  }
  await cashierStore.getPendingDiscrepancies()
  await cashierStore.fetchPendingHandovers()
}

// Pull to refresh - Use refreshData AFTER it's defined
const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)

const loadPayments = async () => {
  if (selectedShift.value) {
    await Promise.all([
      cashierStore.getShiftStatus(selectedShift.value),
      cashierStore.getPaymentsByShift(selectedShift.value)
    ])
  }
}

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
    await refreshData()
  } catch (error) {
    console.error('Override failed:', error)
  }
}

const handleReportDiscrepancy = async (data) => {
  try {
    await cashierStore.reportDiscrepancy({
      order_id: selectedPayment.value.order_id,
      ...data
    })
    showDiscrepancy.value = false
    await refreshData()
  } catch (error) {
    console.error('Report discrepancy failed:', error)
  }
}

const lockOrder = async (orderId) => {
  if (confirm('Bạn có chắc muốn khóa order này? Không thể hoàn tác!')) {
    try {
      await cashierStore.lockOrder(orderId)
      await refreshData()
    } catch (error) {
      console.error('Lock order failed:', error)
    }
  }
}

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

const performReconciliation = async () => {
  if (!reconciliationForm.value.actualCash) return
  
  if (confirm('Xác nhận đối soát tiền mặt? Không thể thay đổi sau khi xác nhận.')) {
    try {
      await cashierStore.reconcileCash({
        shift_id: selectedShift.value,
        actual_cash: reconciliationForm.value.actualCash,
        notes: reconciliationForm.value.notes
      })
      reconciliationForm.value = { actualCash: null, notes: '' }
      await refreshData()
    } catch (error) {
      console.error('Reconciliation failed:', error)
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

const getStatusText = (status) => {
  // Handle both shift status and order status using constants
  const statusMap = {
    // Shift statuses
    [SHIFT_STATUS.OPEN]: '🟢 Đang mở',
    [SHIFT_STATUS.CLOSED]: '🔴 Đã đóng',
    // Cashier Shift statuses
    [CASHIER_SHIFT_STATUS.OPEN]: '🟢 Đang mở',
    [CASHIER_SHIFT_STATUS.CLOSURE_INITIATED]: '🟡 Đang đóng',
    [CASHIER_SHIFT_STATUS.CLOSED]: '🔴 Đã đóng',
    // Order statuses
    [ORDER_STATUS.PAID]: '✓ Đã thu',
    [ORDER_STATUS.QUEUED]: '⏳ Chờ pha',
    [ORDER_STATUS.IN_PROGRESS]: '🍹 Đang pha',
    [ORDER_STATUS.READY]: '✅ Sẵn sàng',
    [ORDER_STATUS.SERVED]: '🎯 Hoàn tất',
    [ORDER_STATUS.LOCKED]: '🔒 Đã khóa'
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

const getPaymentMethodBadge = (method) => {
  const badges = {
    [PAYMENT_METHOD.CASH]: 'px-2 py-1 text-xs rounded-full bg-green-100 text-green-700 font-medium',
    [PAYMENT_METHOD.TRANSFER]: 'px-2 py-1 text-xs rounded-full bg-blue-100 text-blue-700 font-medium',
    [PAYMENT_METHOD.QR]: 'px-2 py-1 text-xs rounded-full bg-purple-100 text-purple-700 font-medium'
  }
  return badges[method] || 'px-2 py-1 text-xs rounded-full bg-gray-100 text-gray-700 font-medium'
}

const getPaymentMethodText = (method) => {
  const texts = {
    [PAYMENT_METHOD.CASH]: '💵 Tiền mặt',
    [PAYMENT_METHOD.TRANSFER]: '💳 CK',
    [PAYMENT_METHOD.QR]: '📱 QR'
  }
  return texts[method] || method
}

const getStatusBadge = (status) => {
  const badges = {
    [ORDER_STATUS.PAID]: 'px-3 py-1 text-xs rounded-full bg-green-100 text-green-700 font-medium inline-block',
    [ORDER_STATUS.QUEUED]: 'px-3 py-1 text-xs rounded-full bg-yellow-100 text-yellow-700 font-medium inline-block',
    [ORDER_STATUS.IN_PROGRESS]: 'px-3 py-1 text-xs rounded-full bg-blue-100 text-blue-700 font-medium inline-block',
    [ORDER_STATUS.READY]: 'px-3 py-1 text-xs rounded-full bg-purple-100 text-purple-700 font-medium inline-block',
    [ORDER_STATUS.SERVED]: 'px-3 py-1 text-xs rounded-full bg-green-100 text-green-700 font-medium inline-block',
    [ORDER_STATUS.LOCKED]: 'px-3 py-1 text-xs rounded-full bg-red-100 text-red-700 font-medium inline-block'
  }
  return badges[status] || 'px-3 py-1 text-xs rounded-full bg-gray-100 text-gray-700 font-medium inline-block'
}

const getDifferenceClass = (difference) => {
  if (difference > 0) return 'text-green-600'
  if (difference < 0) return 'text-red-600'
  return 'text-gray-600'
}

// Lifecycle
onMounted(async () => {
  // Fetch cashier shifts for cashier shift manager
  await cashierShiftStore.fetchMyCashierShifts()
  // Fetch all shifts (including waiter shifts) for payment monitoring
  await shiftStore.fetchAllShifts()
  console.log('DEBUG: Waiter shifts loaded:', shiftStore.shifts)
  console.log('DEBUG: First shift structure:', shiftStore.shifts[0])
  await cashierStore.getPendingDiscrepancies()
  await cashierStore.fetchPendingHandovers()
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
