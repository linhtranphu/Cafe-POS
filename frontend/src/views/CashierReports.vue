<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Pull to Refresh Indicator -->
    <PullToRefresh 
      :pull-distance="pullDistance" 
      :is-refreshing="isRefreshing"
      :threshold="80" />
    
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3" style="padding-top: max(0.75rem, env(safe-area-inset-top))">
        <div class="flex items-center justify-between mb-2">
          <div>
            <h1 class="text-xl font-bold text-gray-800">📊 Báo cáo</h1>
            <p class="text-xs text-gray-600">Thu ngân & doanh thu</p>
          </div>
          <button v-if="currentReport" @click="clearCurrentReport"
            class="text-blue-500 text-sm font-medium">
            ← Quay lại
          </button>
        </div>
        
        <!-- Quick Stats (when no report shown) -->
        <div v-if="!currentReport && quickStats" class="grid grid-cols-3 gap-2 mt-2">
          <div class="bg-blue-50 rounded-lg p-2 text-center">
            <div class="text-sm font-bold text-blue-600">{{ quickStats.totalShifts }}</div>
            <div class="text-[10px] text-gray-600">Ca làm</div>
          </div>
          <div class="bg-green-50 rounded-lg p-2 text-center">
            <div class="text-xs font-bold text-green-600">{{ formatCompactPrice(quickStats.totalRevenue) }}</div>
            <div class="text-[10px] text-gray-600">Doanh thu</div>
          </div>
          <div class="bg-purple-50 rounded-lg p-2 text-center">
            <div class="text-sm font-bold text-purple-600">{{ quickStats.totalOrders }}</div>
            <div class="text-[10px] text-gray-600">Đơn hàng</div>
          </div>
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

      <!-- Report Generation Cards (only show when no current report) -->
      <div v-if="!currentReport" class="space-y-3 mb-4">
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
      </div>

      <!-- Current Report Display -->
      <div v-if="currentReport" class="bg-white rounded-2xl p-4 shadow-sm mb-4">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-base font-bold text-gray-800">{{ currentReport.title }}</h2>
          <button
            @click="printReport"
            class="px-3 py-2 bg-blue-500 text-white text-sm rounded-lg active:scale-95 transition-transform"
          >
            🖨️ In
          </button>
        </div>

        <!-- Report Content -->
        <div id="report-content" class="space-y-4">
          <!-- Header -->
          <div class="text-center border-b-2 border-gray-200 pb-3">
            <h1 class="text-lg font-bold text-gray-800">QUÁN CAFÉ</h1>
            <h2 class="text-sm font-medium text-gray-700">{{ currentReport.title }}</h2>
            <p class="text-xs text-gray-500 mt-1">{{ formatDateTime(currentReport.generated_at) }}</p>
          </div>

          <!-- Summary Stats -->
          <div class="grid grid-cols-2 gap-2">
            <div class="bg-blue-50 rounded-xl p-3 text-center">
              <div class="text-xl font-bold text-blue-600">{{ currentReport.total_orders }}</div>
              <div class="text-xs text-gray-600">Tổng đơn</div>
            </div>
            <div class="bg-green-50 rounded-xl p-3 text-center">
              <div class="text-sm font-bold text-green-600">{{ formatCompactPrice(currentReport.total_revenue) }}</div>
              <div class="text-xs text-gray-600">Doanh thu</div>
            </div>
            <div class="bg-yellow-50 rounded-xl p-3 text-center">
              <div class="text-sm font-bold text-yellow-600">{{ formatCompactPrice(currentReport.cash_revenue) }}</div>
              <div class="text-xs text-gray-600">💵 Tiền mặt</div>
            </div>
            <div class="bg-purple-50 rounded-xl p-3 text-center">
              <div class="text-sm font-bold text-purple-600">{{ formatCompactPrice(currentReport.transfer_revenue + currentReport.qr_revenue) }}</div>
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
                  <div class="font-bold text-gray-800 text-sm">{{ formatCompactPrice(currentReport.cash_revenue) }}</div>
                  <div class="text-xs text-gray-500">{{ getPercentage(currentReport.cash_revenue, currentReport.total_revenue) }}%</div>
                </div>
              </div>
              <div class="flex justify-between items-center bg-gray-50 rounded-lg p-3">
                <span class="text-sm text-gray-700">💳 Chuyển khoản</span>
                <div class="text-right">
                  <div class="font-bold text-gray-800 text-sm">{{ formatCompactPrice(currentReport.transfer_revenue) }}</div>
                  <div class="text-xs text-gray-500">{{ getPercentage(currentReport.transfer_revenue, currentReport.total_revenue) }}%</div>
                </div>
              </div>
              <div class="flex justify-between items-center bg-gray-50 rounded-lg p-3">
                <span class="text-sm text-gray-700">📱 QR Code</span>
                <div class="text-right">
                  <div class="font-bold text-gray-800 text-sm">{{ formatCompactPrice(currentReport.qr_revenue) }}</div>
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
                <span class="font-medium text-gray-800 text-sm">{{ formatCompactPrice(currentReport.reconciliation.expected_cash) }}</span>
              </div>
              <div class="flex justify-between items-center">
                <span class="text-sm text-gray-600">Thực tế:</span>
                <span class="font-medium text-gray-800 text-sm">{{ formatCompactPrice(currentReport.reconciliation.actual_cash) }}</span>
              </div>
              <div class="flex justify-between items-center pt-2 border-t border-green-200">
                <span class="text-sm font-medium text-gray-700">Chênh lệch:</span>
                <span :class="getDifferenceClass(currentReport.reconciliation.difference)" class="font-bold text-sm">
                  {{ formatCompactPrice(currentReport.reconciliation.difference) }}
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
                  <span class="font-bold text-gray-800 text-sm">{{ formatCompactPrice(audit.amount) }}</span>
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

      <!-- Report History (only show when no current report) -->
      <div v-if="!currentReport" class="mb-4">
        <h2 class="text-base font-bold text-gray-800 mb-3">📚 Lịch sử báo cáo</h2>
        
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
                <h3 class="font-bold text-gray-800 text-sm">{{ getReportTitle(report) }}</h3>
                <p class="text-xs text-gray-500">{{ formatDateTime(report.generated_at) }}</p>
              </div>
              <span class="text-sm text-blue-500 font-medium">Xem →</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">{{ report.total_orders }} đơn hàng</span>
              <span class="font-bold text-green-600 text-sm">{{ formatCompactPrice(report.total_revenue) }}</span>
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
import PullToRefresh from '../components/PullToRefresh.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'

const cashierStore = useCashierStore()
const shiftStore = useShiftStore()

const selectedShiftForReport = ref('')
const selectedDate = ref(new Date().toISOString().split('T')[0])
const currentReport = ref(null)

// Computed
const loading = computed(() => cashierStore.loading)
const error = computed(() => cashierStore.error)
const reports = computed(() => cashierStore.reports)
const availableShifts = computed(() => shiftStore.shifts)

// Quick stats for header (when no report shown)
const quickStats = computed(() => {
  if (reports.value.length === 0) return null
  
  const totalShifts = reports.value.filter(r => r.shift).length
  const totalRevenue = reports.value.reduce((sum, r) => sum + (r.total_revenue || 0), 0)
  const totalOrders = reports.value.reduce((sum, r) => sum + (r.total_orders || 0), 0)
  
  return {
    totalShifts,
    totalRevenue,
    totalOrders
  }
})

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

const viewReport = (report) => {
  currentReport.value = {
    ...report,
    title: getReportTitle(report)
  }
  // Scroll to top of container
  const container = document.querySelector('.overflow-y-auto')
  if (container) {
    container.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

const clearCurrentReport = () => {
  currentReport.value = null
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
const formatCompactPrice = (value) => {
  if (value === undefined || value === null || isNaN(value)) {
    return '0đ'
  }
  
  // For millions (triệu)
  if (value >= 1000000) {
    const millions = value / 1000000
    // If it's a whole number of millions
    if (millions % 1 === 0) {
      return `${millions}tr`
    }
    // If it has decimals, show 1 decimal place
    return `${millions.toFixed(1)}tr`
  }
  
  // For thousands (nghìn)
  if (value >= 1000) {
    const thousands = value / 1000
    // If it's a whole number of thousands
    if (thousands % 1 === 0) {
      return `${thousands}k`
    }
    // If it has decimals, show 1 decimal place
    return `${thousands.toFixed(1)}k`
  }
  
  // For small numbers, show as is
  return `${value}đ`
}

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

// Refresh data function
const refreshData = async () => {
  await shiftStore.fetchAllShifts()
}

// Pull to refresh
const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)

// Lifecycle
onMounted(async () => {
  await refreshData()
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
