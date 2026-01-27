<template>
  <div class="min-h-screen bg-gray-100">
    <Navigation />
    <div class="p-6 space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold text-gray-900">Báo cáo Thu ngân</h1>
    </div>

    <!-- Error Alert -->
    <div v-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
      {{ error }}
      <button @click="clearError" class="float-right font-bold">&times;</button>
    </div>

    <!-- Report Generation -->
    <div class="bg-white rounded-lg shadow p-6">
      <h2 class="text-lg font-semibold mb-4">Tạo báo cáo</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <!-- Shift Report -->
        <div class="border rounded-lg p-4">
          <h3 class="font-medium mb-3">Báo cáo ca</h3>
          <select v-model="selectedShiftForReport" class="w-full border rounded px-3 py-2 mb-3">
            <option value="">Chọn ca</option>
            <option v-for="shift in availableShifts" :key="shift.id" :value="shift.id">
              {{ shift.type }} - {{ formatDate(shift.start_time) }}
            </option>
          </select>
          <button
            @click="generateShiftReport"
            :disabled="!selectedShiftForReport || loading"
            class="w-full px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
          >
            Tạo báo cáo ca
          </button>
        </div>

        <!-- Daily Report -->
        <div class="border rounded-lg p-4">
          <h3 class="font-medium mb-3">Báo cáo ngày</h3>
          <input
            v-model="selectedDate"
            type="date"
            class="w-full border rounded px-3 py-2 mb-3"
          />
          <button
            @click="generateDailyReport"
            :disabled="!selectedDate || loading"
            class="w-full px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 disabled:opacity-50"
          >
            Tạo báo cáo ngày
          </button>
        </div>

        <!-- Shift Handover -->
        <div class="border rounded-lg p-4">
          <h3 class="font-medium mb-3">Bàn giao ca</h3>
          <input
            v-model="handoverForm.toCashierID"
            type="text"
            placeholder="ID thu ngân tiếp nhận"
            class="w-full border rounded px-3 py-2 mb-2"
          />
          <textarea
            v-model="handoverForm.notes"
            placeholder="Ghi chú bàn giao"
            class="w-full border rounded px-3 py-2 mb-3 h-20 resize-none"
          ></textarea>
          <button
            @click="performHandover"
            :disabled="!handoverForm.toCashierID || loading"
            class="w-full px-4 py-2 bg-orange-600 text-white rounded hover:bg-orange-700 disabled:opacity-50"
          >
            Bàn giao ca
          </button>
        </div>
      </div>
    </div>

    <!-- Current Report Display -->
    <div v-if="currentReport" class="bg-white rounded-lg shadow p-6">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-lg font-semibold">{{ currentReport.title }}</h2>
        <button
          @click="printReport"
          class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
        >
          <i class="fas fa-print mr-2"></i>In báo cáo
        </button>
      </div>

      <!-- Report Content -->
      <div id="report-content" class="space-y-6">
        <!-- Header -->
        <div class="text-center border-b pb-4">
          <h1 class="text-xl font-bold">QUÁN CAFÉ</h1>
          <h2 class="text-lg">{{ currentReport.title }}</h2>
          <p class="text-sm text-gray-600">Ngày tạo: {{ formatDateTime(currentReport.generated_at) }}</p>
        </div>

        <!-- Summary -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div class="text-center p-4 bg-blue-50 rounded">
            <div class="text-2xl font-bold text-blue-600">{{ currentReport.total_orders }}</div>
            <div class="text-sm text-gray-600">Tổng đơn hàng</div>
          </div>
          <div class="text-center p-4 bg-green-50 rounded">
            <div class="text-2xl font-bold text-green-600">{{ formatCurrency(currentReport.total_revenue) }}</div>
            <div class="text-sm text-gray-600">Tổng doanh thu</div>
          </div>
          <div class="text-center p-4 bg-yellow-50 rounded">
            <div class="text-2xl font-bold text-yellow-600">{{ formatCurrency(currentReport.cash_revenue) }}</div>
            <div class="text-sm text-gray-600">Tiền mặt</div>
          </div>
          <div class="text-center p-4 bg-purple-50 rounded">
            <div class="text-2xl font-bold text-purple-600">{{ formatCurrency(currentReport.transfer_revenue + currentReport.qr_revenue) }}</div>
            <div class="text-sm text-gray-600">Chuyển khoản</div>
          </div>
        </div>

        <!-- Revenue Breakdown -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="bg-gray-50 p-4 rounded">
            <h3 class="font-medium mb-2">Tiền mặt</h3>
            <div class="text-lg font-bold">{{ formatCurrency(currentReport.cash_revenue) }}</div>
            <div class="text-sm text-gray-600">{{ getPercentage(currentReport.cash_revenue, currentReport.total_revenue) }}%</div>
          </div>
          <div class="bg-gray-50 p-4 rounded">
            <h3 class="font-medium mb-2">Chuyển khoản</h3>
            <div class="text-lg font-bold">{{ formatCurrency(currentReport.transfer_revenue) }}</div>
            <div class="text-sm text-gray-600">{{ getPercentage(currentReport.transfer_revenue, currentReport.total_revenue) }}%</div>
          </div>
          <div class="bg-gray-50 p-4 rounded">
            <h3 class="font-medium mb-2">QR Code</h3>
            <div class="text-lg font-bold">{{ formatCurrency(currentReport.qr_revenue) }}</div>
            <div class="text-sm text-gray-600">{{ getPercentage(currentReport.qr_revenue, currentReport.total_revenue) }}%</div>
          </div>
        </div>

        <!-- Reconciliation -->
        <div v-if="currentReport.reconciliation" class="bg-gray-50 p-4 rounded">
          <h3 class="font-medium mb-3">Đối soát tiền mặt</h3>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <div class="text-sm text-gray-600">Tiền mặt dự kiến</div>
              <div class="font-medium">{{ formatCurrency(currentReport.reconciliation.expected_cash) }}</div>
            </div>
            <div>
              <div class="text-sm text-gray-600">Tiền mặt thực tế</div>
              <div class="font-medium">{{ formatCurrency(currentReport.reconciliation.actual_cash) }}</div>
            </div>
            <div>
              <div class="text-sm text-gray-600">Chênh lệch</div>
              <div :class="getDifferenceClass(currentReport.reconciliation.difference)">
                {{ formatCurrency(currentReport.reconciliation.difference) }}
                <span class="text-xs ml-1">({{ currentReport.reconciliation.status }})</span>
              </div>
            </div>
          </div>
          <div v-if="currentReport.reconciliation.notes" class="mt-2 text-sm text-gray-600">
            Ghi chú: {{ currentReport.reconciliation.notes }}
          </div>
        </div>

        <!-- Audit Trail -->
        <div v-if="currentReport.audits && currentReport.audits.length > 0" class="space-y-3">
          <h3 class="font-medium">Nhật ký kiểm toán</h3>
          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Thời gian</th>
                  <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Hành động</th>
                  <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Order ID</th>
                  <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Lý do</th>
                  <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Số tiền</th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="audit in currentReport.audits" :key="audit.id">
                  <td class="px-4 py-2 text-sm">{{ formatDateTime(audit.audited_at) }}</td>
                  <td class="px-4 py-2 text-sm">
                    <span :class="getAuditActionClass(audit.action)">
                      {{ getAuditActionText(audit.action) }}
                    </span>
                  </td>
                  <td class="px-4 py-2 text-sm font-mono">{{ audit.order_id }}</td>
                  <td class="px-4 py-2 text-sm">{{ audit.reason }}</td>
                  <td class="px-4 py-2 text-sm font-medium">{{ formatCurrency(audit.amount) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- Report History -->
    <div class="bg-white rounded-lg shadow p-6">
      <h2 class="text-lg font-semibold mb-4">Lịch sử báo cáo</h2>
      <div v-if="reports.length > 0" class="space-y-3">
        <div
          v-for="report in reports"
          :key="report.generated_at"
          class="border rounded-lg p-4 hover:bg-gray-50 cursor-pointer"
          @click="viewReport(report)"
        >
          <div class="flex justify-between items-center">
            <div>
              <div class="font-medium">{{ getReportTitle(report) }}</div>
              <div class="text-sm text-gray-600">{{ formatDateTime(report.generated_at) }}</div>
            </div>
            <div class="text-right">
              <div class="font-medium">{{ formatCurrency(report.total_revenue) }}</div>
              <div class="text-sm text-gray-600">{{ report.total_orders }} đơn</div>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="text-center py-8 text-gray-500">
        Chưa có báo cáo nào
      </div>
    </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useCashierStore } from '../stores/cashier'
import { useShiftStore } from '../stores/shift'
import Navigation from '../components/Navigation.vue'

export default {
  name: 'CashierReports',
  components: {
    Navigation
  },
  setup() {
    const cashierStore = useCashierStore()
    const shiftStore = useShiftStore()

    const selectedShiftForReport = ref('')
    const selectedDate = ref(new Date().toISOString().split('T')[0])
    const currentReport = ref(null)
    const handoverForm = ref({
      toCashierID: '',
      notes: ''
    })

    const loading = computed(() => cashierStore.loading)
    const error = computed(() => cashierStore.error)
    const reports = computed(() => cashierStore.reports)
    const availableShifts = computed(() => shiftStore.shifts)

    onMounted(async () => {
      await shiftStore.fetchAllShifts()
    })

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
          title: `Báo cáo ngày ${selectedDate.value}`
        }
      } catch (error) {
        console.error('Generate daily report failed:', error)
      }
    }

    const performHandover = async () => {
      if (!confirm('Bạn có chắc muốn bàn giao ca?')) return

      try {
        await cashierStore.handoverShift({
          from_cashier_id: 'current_user', // Should get from auth store
          to_cashier_id: handoverForm.value.toCashierID,
          notes: handoverForm.value.notes
        })
        handoverForm.value = { toCashierID: '', notes: '' }
        alert('Bàn giao ca thành công!')
      } catch (error) {
        console.error('Handover failed:', error)
      }
    }

    const viewReport = (report) => {
      currentReport.value = {
        ...report,
        title: getReportTitle(report)
      }
    }

    const printReport = () => {
      const printContent = document.getElementById('report-content')
      const printWindow = window.open('', '_blank')
      printWindow.document.write(`
        <html>
          <head>
            <title>Báo cáo</title>
            <style>
              body { font-family: Arial, sans-serif; margin: 20px; }
              .text-center { text-align: center; }
              .grid { display: grid; gap: 1rem; }
              .grid-cols-2 { grid-template-columns: repeat(2, 1fr); }
              .grid-cols-3 { grid-template-columns: repeat(3, 1fr); }
              .grid-cols-4 { grid-template-columns: repeat(4, 1fr); }
              .p-4 { padding: 1rem; }
              .bg-gray-50 { background-color: #f9fafb; }
              .rounded { border-radius: 0.375rem; }
              .font-bold { font-weight: bold; }
              .font-medium { font-weight: 500; }
              .text-sm { font-size: 0.875rem; }
              .text-lg { font-size: 1.125rem; }
              .text-xl { font-size: 1.25rem; }
              .text-2xl { font-size: 1.5rem; }
              .border-b { border-bottom: 1px solid #e5e7eb; }
              .pb-4 { padding-bottom: 1rem; }
              .mb-2 { margin-bottom: 0.5rem; }
              .mb-3 { margin-bottom: 0.75rem; }
              .mb-4 { margin-bottom: 1rem; }
              .space-y-3 > * + * { margin-top: 0.75rem; }
              .space-y-6 > * + * { margin-top: 1.5rem; }
              table { width: 100%; border-collapse: collapse; }
              th, td { padding: 0.5rem; border: 1px solid #e5e7eb; text-align: left; }
              th { background-color: #f9fafb; font-weight: 500; }
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

    const getPercentage = (value, total) => {
      if (total === 0) return 0
      return Math.round((value / total) * 100)
    }

    const getDifferenceClass = (difference) => {
      if (difference > 0) return 'font-medium text-green-600'
      if (difference < 0) return 'font-medium text-red-600'
      return 'font-medium text-gray-600'
    }

    const getAuditActionClass = (action) => {
      const classes = {
        'CANCEL': 'px-2 py-1 text-xs rounded bg-red-100 text-red-800',
        'REFUND': 'px-2 py-1 text-xs rounded bg-orange-100 text-orange-800',
        'OVERRIDE': 'px-2 py-1 text-xs rounded bg-yellow-100 text-yellow-800',
        'LOCK': 'px-2 py-1 text-xs rounded bg-gray-100 text-gray-800'
      }
      return classes[action] || 'px-2 py-1 text-xs rounded bg-gray-100 text-gray-800'
    }

    const getAuditActionText = (action) => {
      const texts = {
        'CANCEL': 'Hủy',
        'REFUND': 'Hoàn tiền',
        'OVERRIDE': 'Điều chỉnh',
        'LOCK': 'Khóa'
      }
      return texts[action] || action
    }

    const getReportTitle = (report) => {
      if (report.shift) {
        return `Báo cáo ca ${report.shift.type}`
      }
      return 'Báo cáo tổng hợp'
    }

    return {
      selectedShiftForReport,
      selectedDate,
      currentReport,
      handoverForm,
      loading,
      error,
      reports,
      availableShifts,
      generateShiftReport,
      generateDailyReport,
      performHandover,
      viewReport,
      printReport,
      clearError,
      formatCurrency,
      formatDate,
      formatDateTime,
      getPercentage,
      getDifferenceClass,
      getAuditActionClass,
      getAuditActionText,
      getReportTitle
    }
  }
}
</script>