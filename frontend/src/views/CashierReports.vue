<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm">
      <div class="px-4 py-4">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-800">📊 Báo cáo</h1>
            <p class="text-sm text-gray-600">Thu ngân & doanh thu</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="px-4 py-4 pb-24">
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

      <!-- Report Generation Cards -->
      <div class="space-y-3 mb-4">
        <!-- Shift Report -->
        <div class="bg-white rounded-2xl p-4 shadow-sm">
          <h3 class="font-bold text-gray-800 mb-3">📋 Báo cáo ca</h3>
          <select 
            v-model="selectedShiftForReport" 
            class="w-full border-2 border-gray-300 rounded-xl px-4 py-3 text-base mb-3 focus:outline-none focus:border-blue-500"
          >
            <option value="">-- Chọn ca --</option>
            <option v-for="shift in availableShifts" :key="shift.id" :value="shift.id">
              {{ getShiftTypeText(shift.type) }} - {{ formatDate(shift.started_at) }} - {{ shift.user_name }}
            </option>
          </select>
          <button
            @click="generateShiftReport"
            :disabled="!selectedShiftForReport || loading"
            class="w-full py-3 bg-blue-500 text-white rounded-xl font-medium active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ loading ? '⏳ Đang tạo...' : '✓ Tạo báo cáo ca' }}
          </button>
        </div>

        <!-- Daily Report -->
        <div class="bg-white rounded-2xl p-4 shadow-sm">
          <h3 class="font-bold text-gray-800 mb-3">📅 Báo cáo ngày</h3>
          <input
            v-model="selectedDate"
            type="date"
            class="w-full border-2 border-gray-300 rounded-xl px-3 py-2 text-sm mb-3 focus:outline-none focus:border-green-500 appearance-none"
          />
          <button
            @click="generateDailyReport"
            :disabled="!selectedDate || loading"
            class="w-full py-3 bg-green-500 text-white rounded-xl font-medium active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ loading ? '⏳ Đang tạo...' : '✓ Tạo báo cáo ngày' }}
          </button>
        </div>

        <!-- Shift Handover -->
        <div class="bg-white rounded-2xl p-4 shadow-sm">
          <h3 class="font-bold text-gray-800 mb-3">🔄 Bàn giao ca</h3>
          <input
            v-model="handoverForm.toCashierID"
            type="text"
            placeholder="ID thu ngân tiếp nhận"
            class="w-full border-2 border-gray-300 rounded-xl px-4 py-3 text-base mb-3 focus:outline-none focus:border-orange-500"
          />
          <textarea
            v-model="handoverForm.notes"
            placeholder="Ghi chú bàn giao (tùy chọn)"
            rows="2"
            class="w-full border-2 border-gray-300 rounded-xl px-4 py-3 text-base mb-3 focus:outline-none focus:border-orange-500 resize-none"
          ></textarea>
          <button
            @click="performHandover"
            :disabled="!handoverForm.toCashierID || loading"
            class="w-full py-3 bg-orange-500 text-white rounded-xl font-medium active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ loading ? '⏳ Đang xử lý...' : '✓ Xác nhận bàn giao' }}
          </button>
        </div>
      </div>

      <!-- Current Report Display -->
      <div v-if="currentReport" class="bg-white rounded-2xl p-4 shadow-sm mb-4">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-lg font-bold text-gray-800">{{ currentReport.title }}</h2>
          <button
            @click="printReport"
            class="p-2 bg-gray-100 text-gray-700 rounded-lg active:scale-95 transition-transform"
          >
            🖨️
          </button>
        </div>

        <!-- Report Content -->
        <div id="report-content" class="space-y-4">
          <!-- Header -->
          <div class="text-center border-b-2 border-gray-200 pb-4">
            <h1 class="text-xl font-bold text-gray-800">QUÁN CAFÉ</h1>
            <h2 class="text-base font-medium text-gray-700">{{ currentReport.title }}</h2>
            <p class="text-xs text-gray-500 mt-1">{{ formatDateTime(currentReport.generated_at) }}</p>
          </div>

          <!-- Summary Stats -->
          <div class="grid grid-cols-2 gap-3">
            <div class="bg-blue-50 rounded-xl p-3 text-center">
              <div class="text-2xl font-bold text-blue-600">{{ currentReport.total_orders }}</div>
              <div class="text-xs text-gray-600">Tổng đơn</div>
            </div>
            <div class="bg-green-50 rounded-xl p-3 text-center">
              <div class="text-lg font-bold text-green-600">{{ formatPrice(currentReport.total_revenue) }}</div>
              <div class="text-xs text-gray-600">Doanh thu</div>
            </div>
            <div class="bg-yellow-50 rounded-xl p-3 text-center">
              <div class="text-lg font-bold text-yellow-600">{{ formatPrice(currentReport.cash_revenue) }}</div>
              <div class="text-xs text-gray-600">💵 Tiền mặt</div>
            </div>
            <div class="bg-purple-50 rounded-xl p-3 text-center">
              <div class="text-lg font-bold text-purple-600">{{ formatPrice(currentReport.transfer_revenue + currentReport.qr_revenue) }}</div>
              <div class="text-xs text-gray-600">💳 Chuyển khoản</div>
            </div>
          </div>

          <!-- Revenue Breakdown -->
          <div class="space-y-2">
            <h3 class="font-bold text-gray-800 text-sm">Chi tiết thanh toán</h3>
            <div class="space-y-2">
              <div class="flex justify-between items-center bg-gray-50 rounded-lg p-3">
                <span class="text-sm text-gray-700">💵 Tiền mặt</span>
                <div class="text-right">
                  <div class="font-bold text-gray-800">{{ formatPrice(currentReport.cash_revenue) }}</div>
                  <div class="text-xs text-gray-500">{{ getPercentage(currentReport.cash_revenue, currentReport.total_revenue) }}%</div>
                </div>
              </div>
              <div class="flex justify-between items-center bg-gray-50 rounded-lg p-3">
                <span class="text-sm text-gray-700">💳 Chuyển khoản</span>
                <div class="text-right">
                  <div class="font-bold text-gray-800">{{ formatPrice(currentReport.transfer_revenue) }}</div>
                  <div class="text-xs text-gray-500">{{ getPercentage(currentReport.transfer_revenue, currentReport.total_revenue) }}%</div>
                </div>
              </div>
              <div class="flex justify-between items-center bg-gray-50 rounded-lg p-3">
                <span class="text-sm text-gray-700">📱 QR Code</span>
                <div class="text-right">
                  <div class="font-bold text-gray-800">{{ formatPrice(currentReport.qr_revenue) }}</div>
                  <div class="text-xs text-gray-500">{{ getPercentage(currentReport.qr_revenue, currentReport.total_revenue) }}%</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Reconciliation -->
          <div v-if="currentReport.reconciliation" class="bg-green-50 rounded-xl p-4">
            <h3 class="font-bold text-gray-800 text-sm mb-3">💰 Đối soát tiền mặt</h3>
            <div class="space-y-2">
              <div class="flex justify-between items-center">
                <span class="text-sm text-gray-600">Dự kiến:</span>
                <span class="font-medium text-gray-800">{{ formatPrice(currentReport.reconciliation.expected_cash) }}</span>
              </div>
              <div class="flex justify-between items-center">
                <span class="text-sm text-gray-600">Thực tế:</span>
                <span class="font-medium text-gray-800">{{ formatPrice(currentReport.reconciliation.actual_cash) }}</span>
              </div>
              <div class="flex justify-between items-center pt-2 border-t border-green-200">
                <span class="text-sm font-medium text-gray-700">Chênh lệch:</span>
                <span :class="getDifferenceClass(currentReport.reconciliation.difference)" class="font-bold">
                  {{ formatPrice(currentReport.reconciliation.difference) }}
                </span>
              </div>
            </div>
            <div v-if="currentReport.reconciliation.notes" class="mt-3 text-xs text-gray-600 bg-white rounded-lg p-2">
              <span class="font-medium">Ghi chú:</span> {{ currentReport.reconciliation.notes }}
            </div>
          </div>

          <!-- Audit Trail -->
          <div v-if="currentReport.audits && currentReport.audits.length > 0" class="space-y-2">
            <h3 class="font-bold text-gray-800 text-sm">📝 Nhật ký kiểm toán</h3>
            <div class="space-y-2">
              <div
                v-for="audit in currentReport.audits"
                :key="audit.id"
                class="bg-gray-50 rounded-lg p-3"
              >
                <div class="flex justify-between items-start mb-2">
                  <span :class="getAuditActionBadge(audit.action)">
                    {{ getAuditActionText(audit.action) }}
                  </span>
                  <span class="font-bold text-gray-800">{{ formatPrice(audit.amount) }}</span>
                </div>
                <div class="text-xs text-gray-600">
                  <div>Order: #{{ audit.order_id?.slice(-6) }}</div>
                  <div v-if="audit.reason">Lý do: {{ audit.reason }}</div>
                  <div class="text-gray-500 mt-1">{{ formatDateTime(audit.audited_at) }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Report History -->
      <div class="mb-4">
        <h2 class="text-lg font-bold text-gray-800 mb-3">📚 Lịch sử báo cáo</h2>
        
        <div v-if="reports.length === 0" class="text-center py-12 bg-white rounded-2xl">
          <div class="text-5xl mb-3">📭</div>
          <p class="text-gray-500">Chưa có báo cáo nào</p>
          <p class="text-sm text-gray-400 mt-1">Tạo báo cáo mới ở trên</p>
        </div>

        <div v-else class="space-y-3">
          <div
            v-for="report in reports"
            :key="report.generated_at"
            @click="viewReport(report)"
            class="bg-white rounded-xl p-4 shadow-sm active:scale-98 transition-transform"
          >
            <div class="flex justify-between items-start mb-2">
              <div>
                <h3 class="font-bold text-gray-800">{{ getReportTitle(report) }}</h3>
                <p class="text-xs text-gray-500">{{ formatDateTime(report.generated_at) }}</p>
              </div>
              <span class="text-sm text-blue-500 font-medium">Xem →</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">{{ report.total_orders }} đơn hàng</span>
              <span class="font-bold text-green-600">{{ formatPrice(report.total_revenue) }}</span>
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
import { ref, computed, onMounted } from 'vue'
import { useCashierStore } from '../stores/cashier'
import { useShiftStore } from '../stores/shift'
import BottomNav from '../components/BottomNav.vue'

const cashierStore = useCashierStore()
const shiftStore = useShiftStore()

const selectedShiftForReport = ref('')
const selectedDate = ref(new Date().toISOString().split('T')[0])
const currentReport = ref(null)
const handoverForm = ref({
  toCashierID: '',
  notes: ''
})

// Computed
const loading = computed(() => cashierStore.loading)
const error = computed(() => cashierStore.error)
const reports = computed(() => cashierStore.reports)
const availableShifts = computed(() => shiftStore.shifts)

// Methods
const generateShiftReport = async () => {
  try {
    const report = await cashierStore.generateShiftReport(selectedShiftForReport.value)
    currentReport.value = {
      ...report,
      title: `Báo cáo ca ${report.shift?.type || 'N/A'}`
    }
  } catch (error) {
    console.error('Generate shift report failed:', error)
  }
}

const generateDailyReport = async () => {
  try {
    const report = await cashierStore.getDailyReport(selectedDate.value)
    currentReport.value = {
      ...report,
      title: `Báo cáo ngày ${formatDate(selectedDate.value)}`
    }
  } catch (error) {
    console.error('Generate daily report failed:', error)
  }
}

const performHandover = async () => {
  if (!confirm('Bạn có chắc muốn bàn giao ca? Không thể hoàn tác!')) return

  try {
    await cashierStore.handoverShift({
      from_cashier_id: 'current_user',
      to_cashier_id: handoverForm.value.toCashierID,
      notes: handoverForm.value.notes
    })
    handoverForm.value = { toCashierID: '', notes: '' }
    alert('✓ Bàn giao ca thành công!')
  } catch (error) {
    console.error('Handover failed:', error)
  }
}

const viewReport = (report) => {
  currentReport.value = {
    ...report,
    title: getReportTitle(report)
  }
  // Scroll to top
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const printReport = () => {
  const printContent = document.getElementById('report-content')
  const printWindow = window.open('', '_blank')
  printWindow.document.write(`
    <html>
      <head>
        <title>${currentReport.value.title}</title>
        <style>
          body { font-family: Arial, sans-serif; margin: 20px; }
          .text-center { text-align: center; }
          .space-y-4 > * + * { margin-top: 1rem; }
          .grid { display: grid; gap: 0.75rem; }
          .grid-cols-2 { grid-template-columns: repeat(2, 1fr); }
          .p-3 { padding: 0.75rem; }
          .bg-gray-50 { background-color: #f9fafb; }
          .rounded-xl { border-radius: 0.75rem; }
          .font-bold { font-weight: bold; }
          .font-medium { font-weight: 500; }
          .text-xs { font-size: 0.75rem; }
          .text-sm { font-size: 0.875rem; }
          .text-base { font-size: 1rem; }
          .text-lg { font-size: 1.125rem; }
          .text-xl { font-size: 1.25rem; }
          .text-2xl { font-size: 1.5rem; }
          .border-b-2 { border-bottom: 2px solid #e5e7eb; }
          .pb-4 { padding-bottom: 1rem; }
          .mb-2 { margin-bottom: 0.5rem; }
          .mb-3 { margin-bottom: 0.75rem; }
          .mt-1 { margin-top: 0.25rem; }
          .mt-3 { margin-top: 0.75rem; }
        </style>
      </head>
      <body>
        ${printContent.innerHTML}
      </body>
    </html>
  `)
  printWindow.document.close()
  printWindow.print()
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

const getShiftTypeText = (type) => {
  const types = {
    MORNING: '☀️ Ca sáng',
    AFTERNOON: '🌤️ Ca chiều',
    EVENING: '🌙 Ca tối'
  }
  return types[type] || type
}

const getPercentage = (value, total) => {
  if (!total || total === 0) return 0
  return Math.round((value / total) * 100)
}

const getDifferenceClass = (difference) => {
  if (difference > 0) return 'text-green-600'
  if (difference < 0) return 'text-red-600'
  return 'text-gray-600'
}

const getAuditActionBadge = (action) => {
  const badges = {
    CANCEL: 'px-2 py-1 text-xs rounded-full bg-red-100 text-red-700 font-medium',
    REFUND: 'px-2 py-1 text-xs rounded-full bg-orange-100 text-orange-700 font-medium',
    OVERRIDE: 'px-2 py-1 text-xs rounded-full bg-yellow-100 text-yellow-700 font-medium',
    LOCK: 'px-2 py-1 text-xs rounded-full bg-gray-100 text-gray-700 font-medium'
  }
  return badges[action] || 'px-2 py-1 text-xs rounded-full bg-gray-100 text-gray-700 font-medium'
}

const getAuditActionText = (action) => {
  const texts = {
    CANCEL: '❌ Hủy',
    REFUND: '↩️ Hoàn tiền',
    OVERRIDE: '✏️ Điều chỉnh',
    LOCK: '🔒 Khóa'
  }
  return texts[action] || action
}

const getReportTitle = (report) => {
  if (report.shift) {
    return `Báo cáo ca ${report.shift.type}`
  }
  return 'Báo cáo tổng hợp'
}

// Lifecycle
onMounted(async () => {
  await shiftStore.fetchAllShifts()
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
