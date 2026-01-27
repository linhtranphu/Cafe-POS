<template>
  <div class="min-h-screen bg-gray-100">
    <Navigation />
    <div class="p-6 space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold text-gray-900">Thu ngân Dashboard</h1>
      <div class="flex space-x-2">
        <button
          @click="refreshData"
          :disabled="loading"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
        >
          <i class="fas fa-sync-alt mr-2" :class="{ 'animate-spin': loading }"></i>
          Làm mới
        </button>
      </div>
    </div>

    <!-- Error Alert -->
    <div v-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
      {{ error }}
      <button @click="clearError" class="float-right font-bold">&times;</button>
    </div>

    <!-- Shift Status Card -->
    <div v-if="shiftStatus" class="bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-lg p-6">
      <h2 class="text-xl font-semibold mb-4">Trạng thái ca làm việc</h2>
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div class="text-center">
          <div class="text-2xl font-bold">{{ shiftStatus.total_orders }}</div>
          <div class="text-sm opacity-90">Tổng đơn</div>
        </div>
        <div class="text-center">
          <div class="text-2xl font-bold">{{ formatCurrency(shiftStatus.total_revenue) }}</div>
          <div class="text-sm opacity-90">Tổng doanh thu</div>
        </div>
        <div class="text-center">
          <div class="text-2xl font-bold">{{ formatCurrency(shiftStatus.cash_revenue) }}</div>
          <div class="text-sm opacity-90">Tiền mặt</div>
        </div>
        <div class="text-center">
          <div class="text-2xl font-bold">{{ formatCurrency(shiftStatus.transfer_revenue + shiftStatus.qr_revenue) }}</div>
          <div class="text-sm opacity-90">Chuyển khoản</div>
        </div>
      </div>
    </div>

    <!-- Payment Oversight Panel -->
    <div class="bg-white rounded-lg shadow p-6">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-lg font-semibold">Giám sát thanh toán</h2>
        <select v-model="selectedShift" @change="loadPayments" class="border rounded px-3 py-1">
          <option value="">Chọn ca</option>
          <option v-for="shift in availableShifts" :key="shift.id" :value="shift.id">
            {{ shift.type }} - {{ formatDate(shift.started_at) }}
          </option>
        </select>
      </div>

      <div v-if="payments.length > 0" class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Bàn</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Số tiền</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Phương thức</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Trạng thái</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Thời gian</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Thao tác</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="payment in payments" :key="payment.order_id">
              <td class="px-6 py-4 whitespace-nowrap">{{ payment.table_name }}</td>
              <td class="px-6 py-4 whitespace-nowrap font-medium">{{ formatCurrency(payment.amount) }}</td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span :class="getPaymentMethodClass(payment.payment_method)">
                  {{ getPaymentMethodText(payment.payment_method) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span :class="getStatusClass(payment.status)">
                  {{ getStatusText(payment.status) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDateTime(payment.paid_at) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm space-x-2">
                <button
                  @click="showOverrideModal(payment)"
                  class="text-orange-600 hover:text-orange-900"
                  title="Điều chỉnh"
                >
                  <i class="fas fa-edit"></i>
                </button>
                <button
                  @click="lockOrder(payment.order_id)"
                  class="text-red-600 hover:text-red-900"
                  title="Khóa"
                >
                  <i class="fas fa-lock"></i>
                </button>
                <button
                  @click="showDiscrepancyModal(payment)"
                  class="text-yellow-600 hover:text-yellow-900"
                  title="Báo cáo sai lệch"
                >
                  <i class="fas fa-exclamation-triangle"></i>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else class="text-center py-8 text-gray-500">
        Chưa có thanh toán nào
      </div>
    </div>

    <!-- Discrepancy Panel -->
    <div class="bg-white rounded-lg shadow p-6">
      <h2 class="text-lg font-semibold mb-4">Sai lệch thanh toán</h2>
      <div v-if="pendingDiscrepancies.length > 0" class="space-y-3">
        <div
          v-for="discrepancy in pendingDiscrepancies"
          :key="discrepancy.id"
          class="border border-yellow-200 bg-yellow-50 rounded-lg p-4"
        >
          <div class="flex justify-between items-start">
            <div>
              <div class="font-medium">Order: {{ discrepancy.order_id }}</div>
              <div class="text-sm text-gray-600">{{ discrepancy.reason }}</div>
              <div class="text-sm text-gray-500">Số tiền: {{ formatCurrency(discrepancy.amount) }}</div>
              <div class="text-xs text-gray-400">{{ formatDateTime(discrepancy.reported_at) }}</div>
            </div>
            <button
              @click="resolveDiscrepancy(discrepancy.id)"
              class="px-3 py-1 bg-green-600 text-white text-sm rounded hover:bg-green-700"
            >
              Giải quyết
            </button>
          </div>
        </div>
      </div>
      <div v-else class="text-center py-4 text-gray-500">
        Không có sai lệch nào
      </div>
    </div>

    <!-- Cash Reconciliation -->
    <div class="bg-white rounded-lg shadow p-6">
      <h2 class="text-lg font-semibold mb-4">Đối soát tiền mặt</h2>
      <div v-if="!hasReconciliation && selectedShift" class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Tiền mặt thực tế</label>
            <input
              v-model.number="reconciliationForm.actualCash"
              type="number"
              step="0.01"
              class="w-full border border-gray-300 rounded-lg px-3 py-2"
              placeholder="0.00"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Ghi chú</label>
            <input
              v-model="reconciliationForm.notes"
              type="text"
              class="w-full border border-gray-300 rounded-lg px-3 py-2"
              placeholder="Ghi chú (tùy chọn)"
            />
          </div>
        </div>
        <button
          @click="performReconciliation"
          :disabled="!reconciliationForm.actualCash"
          class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50"
        >
          Đối soát
        </button>
      </div>
      <div v-else-if="reconciliation" class="bg-gray-50 rounded-lg p-4">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <div class="text-sm text-gray-600">Tiền mặt dự kiến</div>
            <div class="font-medium">{{ formatCurrency(reconciliation.expected_cash) }}</div>
          </div>
          <div>
            <div class="text-sm text-gray-600">Tiền mặt thực tế</div>
            <div class="font-medium">{{ formatCurrency(reconciliation.actual_cash) }}</div>
          </div>
          <div>
            <div class="text-sm text-gray-600">Chênh lệch</div>
            <div :class="getDifferenceClass(reconciliation.difference)">
              {{ formatCurrency(reconciliation.difference) }}
            </div>
          </div>
        </div>
        <div v-if="reconciliation.notes" class="mt-2 text-sm text-gray-600">
          Ghi chú: {{ reconciliation.notes }}
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
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useCashierStore } from '../stores/cashier'
import { useShiftStore } from '../stores/shift'
import Navigation from '../components/Navigation.vue'
import OverridePaymentModal from '../components/OverridePaymentModal.vue'
import DiscrepancyModal from '../components/DiscrepancyModal.vue'

export default {
  name: 'CashierDashboard',
  components: {
    Navigation,
    OverridePaymentModal,
    DiscrepancyModal
  },
  setup() {
    const cashierStore = useCashierStore()
    const shiftStore = useShiftStore()

    const selectedShift = ref('')
    const showOverride = ref(false)
    const showDiscrepancy = ref(false)
    const selectedPayment = ref(null)
    const reconciliationForm = ref({
      actualCash: null,
      notes: ''
    })

    const shiftStatus = computed(() => cashierStore.shiftStatus)
    const payments = computed(() => cashierStore.payments)
    const pendingDiscrepancies = computed(() => cashierStore.pendingDiscrepancies)
    const reconciliation = computed(() => cashierStore.reconciliation)
    const hasReconciliation = computed(() => cashierStore.hasReconciliation)
    const loading = computed(() => cashierStore.loading)
    const error = computed(() => cashierStore.error)
    const availableShifts = computed(() => shiftStore.shifts)

    onMounted(async () => {
      await shiftStore.fetchAllShifts()
      await cashierStore.getPendingDiscrepancies()
    })

    const refreshData = async () => {
      if (selectedShift.value) {
        await Promise.all([
          cashierStore.getShiftStatus(selectedShift.value),
          cashierStore.getPaymentsByShift(selectedShift.value)
        ])
      }
      await cashierStore.getPendingDiscrepancies()
    }

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
      } catch (error) {
        console.error('Report discrepancy failed:', error)
      }
    }

    const lockOrder = async (orderId) => {
      if (confirm('Bạn có chắc muốn khóa order này?')) {
        try {
          await cashierStore.lockOrder(orderId)
          await refreshData()
        } catch (error) {
          console.error('Lock order failed:', error)
        }
      }
    }

    const resolveDiscrepancy = async (discrepancyId) => {
      try {
        await cashierStore.resolveDiscrepancy(discrepancyId)
      } catch (error) {
        console.error('Resolve discrepancy failed:', error)
      }
    }

    const performReconciliation = async () => {
      try {
        await cashierStore.reconcileCash({
          shift_id: selectedShift.value,
          actual_cash: reconciliationForm.value.actualCash,
          notes: reconciliationForm.value.notes
        })
        reconciliationForm.value = { actualCash: null, notes: '' }
      } catch (error) {
        console.error('Reconciliation failed:', error)
      }
    }

    const clearError = () => {
      cashierStore.clearError()
    }

    // Utility functions
    const formatCurrency = (amount) => {
      return new Intl.NumberFormat('vi-VN', {
        style: 'currency',
        currency: 'VND'
      }).format(amount)
    }

    const formatDate = (date) => {
      if (!date) return 'N/A'
      const d = new Date(date)
      if (isNaN(d.getTime())) return 'Invalid Date'
      return d.toLocaleDateString('vi-VN')
    }

    const formatDateTime = (date) => {
      if (!date) return 'N/A'
      const d = new Date(date)
      if (isNaN(d.getTime())) return 'Invalid Date'
      return d.toLocaleString('vi-VN')
    }

    const getPaymentMethodClass = (method) => {
      const classes = {
        'CASH': 'px-2 py-1 text-xs rounded-full bg-green-100 text-green-800',
        'TRANSFER': 'px-2 py-1 text-xs rounded-full bg-blue-100 text-blue-800',
        'QR': 'px-2 py-1 text-xs rounded-full bg-purple-100 text-purple-800'
      }
      return classes[method] || 'px-2 py-1 text-xs rounded-full bg-gray-100 text-gray-800'
    }

    const getPaymentMethodText = (method) => {
      const texts = {
        'CASH': 'Tiền mặt',
        'TRANSFER': 'Chuyển khoản',
        'QR': 'QR Code'
      }
      return texts[method] || method
    }

    const getStatusClass = (status) => {
      const classes = {
        'PAID': 'px-2 py-1 text-xs rounded-full bg-green-100 text-green-800',
        'IN_PROGRESS': 'px-2 py-1 text-xs rounded-full bg-yellow-100 text-yellow-800',
        'SERVED': 'px-2 py-1 text-xs rounded-full bg-blue-100 text-blue-800',
        'LOCKED': 'px-2 py-1 text-xs rounded-full bg-red-100 text-red-800'
      }
      return classes[status] || 'px-2 py-1 text-xs rounded-full bg-gray-100 text-gray-800'
    }

    const getStatusText = (status) => {
      const texts = {
        'PAID': 'Đã thanh toán',
        'IN_PROGRESS': 'Đang pha chế',
        'SERVED': 'Đã phục vụ',
        'LOCKED': 'Đã khóa'
      }
      return texts[status] || status
    }

    const getDifferenceClass = (difference) => {
      if (difference > 0) return 'font-medium text-green-600'
      if (difference < 0) return 'font-medium text-red-600'
      return 'font-medium text-gray-600'
    }

    return {
      selectedShift,
      showOverride,
      showDiscrepancy,
      selectedPayment,
      reconciliationForm,
      shiftStatus,
      payments,
      pendingDiscrepancies,
      reconciliation,
      hasReconciliation,
      loading,
      error,
      availableShifts,
      refreshData,
      loadPayments,
      showOverrideModal,
      showDiscrepancyModal,
      handleOverridePayment,
      handleReportDiscrepancy,
      lockOrder,
      resolveDiscrepancy,
      performReconciliation,
      clearError,
      formatCurrency,
      formatDate,
      formatDateTime,
      getPaymentMethodClass,
      getPaymentMethodText,
      getStatusClass,
      getStatusText,
      getDifferenceClass
    }
  }
}
</script>