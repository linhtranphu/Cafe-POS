<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm">
      <div class="px-4 py-3">
        <h1 class="text-xl font-bold text-gray-800">⚖️ Phê duyệt chênh lệch</h1>
      </div>
    </div>

    <!-- Content -->
    <div class="px-4 py-4 pb-24">
      <!-- Error Alert -->
      <div v-if="approvalError" class="bg-red-50 border-2 border-red-200 rounded-2xl p-4 mb-4">
        <div class="flex items-start justify-between">
          <div class="flex items-start gap-3">
            <span class="text-2xl">⚠️</span>
            <div>
              <p class="font-medium text-red-800">Lỗi</p>
              <p class="text-sm text-red-600">{{ approvalError }}</p>
            </div>
          </div>
          <button @click="clearApprovalError" class="text-red-600 text-xl font-bold">×</button>
        </div>
      </div>

      <!-- Summary Cards -->
      <div class="grid grid-cols-2 gap-4 mb-6">
        <div class="bg-gradient-to-r from-orange-400 to-red-500 text-white rounded-2xl p-4 shadow-lg">
          <div class="text-center">
            <p class="text-2xl font-bold">{{ pendingApprovalCount }}</p>
            <p class="text-sm text-orange-100">Chờ phê duyệt</p>
          </div>
        </div>
        <div class="bg-gradient-to-r from-blue-400 to-purple-500 text-white rounded-2xl p-4 shadow-lg">
          <div class="text-center">
            <p class="text-2xl font-bold">{{ formatPrice(totalDiscrepancyAmount) }}</p>
            <p class="text-sm text-blue-100">Tổng chênh lệch</p>
          </div>
        </div>
      </div>

      <!-- Pending Approvals Section -->
      <div class="bg-white rounded-2xl p-6 mb-4 shadow-sm">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-xl font-bold">🕐 Chờ phê duyệt</h3>
          <button @click="refreshApprovals" :disabled="approvalLoading"
            class="p-2 bg-blue-500 text-white rounded-lg active:scale-95 transition-transform disabled:opacity-50">
            <span class="text-sm" :class="{ 'animate-spin': approvalLoading }">🔄</span>
          </button>
        </div>
        
        <div v-if="approvalLoading" class="text-center py-10">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
        </div>
        
        <div v-else-if="pendingApprovals.length === 0" class="text-center py-10">
          <div class="text-4xl mb-2">✅</div>
          <p class="text-gray-500">Không có yêu cầu nào</p>
        </div>
        
        <div v-else class="space-y-4">
          <div v-for="handover in pendingApprovals" :key="handover.id" 
            class="border-2 border-red-200 rounded-xl p-4 bg-red-50">
            <div class="flex justify-between items-start mb-3">
              <div>
                <h4 class="font-bold text-lg">{{ handover.waiter_name }}</h4>
                <p class="text-sm text-gray-600">→ {{ handover.cashier_name }}</p>
                <p class="text-xs text-gray-500">{{ formatDateTime(handover.requested_at) }}</p>
              </div>
              <div class="text-right">
                <p class="text-lg font-bold text-red-600">
                  {{ formatPrice(Math.abs(handover.discrepancy_amount)) }}
                </p>
                <span class="bg-red-100 text-red-800 px-2 py-1 rounded-full text-xs font-medium">
                  {{ handover.discrepancy_amount > 0 ? 'Thừa tiền' : 'Thiếu tiền' }}
                </span>
              </div>
            </div>

            <!-- Handover Details -->
            <div class="grid grid-cols-2 gap-3 mb-3 text-sm">
              <div class="bg-white p-3 rounded-lg">
                <p class="text-gray-500 text-xs">Yêu cầu</p>
                <p class="font-bold">{{ formatPrice(handover.requested_amount) }}</p>
              </div>
              <div class="bg-white p-3 rounded-lg">
                <p class="text-gray-500 text-xs">Thực tế</p>
                <p class="font-bold">{{ formatPrice(handover.actual_amount) }}</p>
              </div>
            </div>

            <!-- Discrepancy Details -->
            <div v-if="handover.discrepancy_reason" class="bg-yellow-50 border border-yellow-200 rounded-lg p-3 mb-3">
              <p class="text-sm font-medium text-yellow-800">Lý do chênh lệch:</p>
              <p class="text-sm text-yellow-700">{{ handover.discrepancy_reason }}</p>
              <p v-if="handover.responsibility" class="text-xs text-yellow-600 mt-1">
                Trách nhiệm: {{ getResponsibilityText(handover.responsibility) }}
              </p>
            </div>

            <!-- Notes -->
            <div v-if="handover.waiter_notes" class="bg-blue-50 p-3 rounded-lg mb-3">
              <p class="text-xs text-blue-700">💬 Waiter: {{ handover.waiter_notes }}</p>
            </div>
            <div v-if="handover.cashier_notes" class="bg-green-50 p-3 rounded-lg mb-3">
              <p class="text-xs text-green-700">💬 Cashier: {{ handover.cashier_notes }}</p>
            </div>

            <!-- Action Buttons -->
            <div class="grid grid-cols-2 gap-2">
              <button @click="showApprovalModal(handover, true)" :disabled="approvalLoading"
                class="bg-green-500 hover:bg-green-600 text-white px-4 py-3 rounded-xl font-bold text-sm active:scale-95 transition-transform disabled:opacity-50">
                ✅ Phê duyệt
              </button>
              <button @click="showApprovalModal(handover, false)" :disabled="approvalLoading"
                class="bg-red-500 hover:bg-red-600 text-white px-4 py-3 rounded-xl font-bold text-sm active:scale-95 transition-transform disabled:opacity-50">
                ❌ Từ chối
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Discrepancy Statistics -->
      <div class="bg-white rounded-2xl p-6 shadow-sm">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-xl font-bold">📊 Thống kê chênh lệch</h3>
          <div class="flex gap-2">
            <select v-model="statsDateRange" @change="loadDiscrepancyStats" 
              class="text-sm border rounded-lg px-3 py-1">
              <option value="today">Hôm nay</option>
              <option value="week">Tuần này</option>
              <option value="month">Tháng này</option>
              <option value="custom">Tùy chọn</option>
            </select>
          </div>
        </div>

        <div v-if="loading" class="text-center py-6">
          <div class="inline-block animate-spin rounded-full h-6 w-6 border-b-2 border-blue-500"></div>
        </div>

        <div v-else-if="hasDiscrepancyStats" class="space-y-4">
          <!-- Overview Cards -->
          <div class="grid grid-cols-2 gap-4">
            <div class="bg-red-50 border border-red-200 rounded-xl p-4">
              <div class="text-center">
                <p class="text-2xl font-bold text-red-600">{{ discrepancyStats.total_shortages }}</p>
                <p class="text-sm text-red-500">Thiếu tiền</p>
                <p class="text-xs text-gray-500">{{ formatPrice(discrepancyStats.total_shortage_amount) }}</p>
              </div>
            </div>
            <div class="bg-green-50 border border-green-200 rounded-xl p-4">
              <div class="text-center">
                <p class="text-2xl font-bold text-green-600">{{ discrepancyStats.total_overages }}</p>
                <p class="text-sm text-green-500">Thừa tiền</p>
                <p class="text-xs text-gray-500">{{ formatPrice(discrepancyStats.total_overage_amount) }}</p>
              </div>
            </div>
          </div>

          <!-- Status Breakdown -->
          <div class="grid grid-cols-3 gap-3">
            <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-3 text-center">
              <p class="text-lg font-bold text-yellow-600">{{ discrepancyStats.pending_count }}</p>
              <p class="text-xs text-yellow-500">Chờ xử lý</p>
            </div>
            <div class="bg-green-50 border border-green-200 rounded-lg p-3 text-center">
              <p class="text-lg font-bold text-green-600">{{ discrepancyStats.resolved_count }}</p>
              <p class="text-xs text-green-500">Đã giải quyết</p>
            </div>
            <div class="bg-red-50 border border-red-200 rounded-lg p-3 text-center">
              <p class="text-lg font-bold text-red-600">{{ discrepancyStats.escalated_count }}</p>
              <p class="text-xs text-red-500">Đã leo thang</p>
            </div>
          </div>

          <!-- Net Discrepancy -->
          <div class="bg-gray-50 border border-gray-200 rounded-xl p-4">
            <div class="text-center">
              <p class="text-sm text-gray-600">Chênh lệch ròng</p>
              <p :class="getNetDiscrepancyClass(discrepancyStats.net_discrepancy)" 
                class="text-2xl font-bold">
                {{ formatPrice(discrepancyStats.net_discrepancy) }}
              </p>
              <p class="text-xs text-gray-500">
                {{ discrepancyStats.net_discrepancy > 0 ? 'Thừa tổng cộng' : 'Thiếu tổng cộng' }}
              </p>
            </div>
          </div>
        </div>

        <div v-else class="text-center py-6">
          <div class="text-3xl mb-2">📊</div>
          <p class="text-gray-500">Chưa có dữ liệu thống kê</p>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />

    <!-- Approval Modal -->
    <transition name="slide-up">
      <div v-if="showApprovalForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6">
          <h3 class="text-xl font-bold mb-4">
            {{ approvalDecision ? '✅ Phê duyệt' : '❌ Từ chối' }} chênh lệch
          </h3>
          
          <div v-if="selectedHandover" class="mb-4">
            <div class="bg-gray-50 p-4 rounded-xl">
              <p class="text-sm text-gray-600">{{ selectedHandover.waiter_name }} → {{ selectedHandover.cashier_name }}</p>
              <p class="text-lg font-bold">
                Chênh lệch: {{ formatPrice(Math.abs(selectedHandover.discrepancy_amount)) }}
                ({{ selectedHandover.discrepancy_amount > 0 ? 'Thừa' : 'Thiếu' }})
              </p>
            </div>
          </div>

          <form @submit.prevent="approveDiscrepancy" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">
                {{ approvalDecision ? 'Lý do phê duyệt' : 'Lý do từ chối' }} *
              </label>
              <textarea v-model="approvalForm.manager_notes" 
                class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500" 
                rows="4" :placeholder="approvalDecision ? 'Nhập lý do phê duyệt...' : 'Nhập lý do từ chối...'" required></textarea>
            </div>

            <div v-if="approvalDecision" class="bg-green-50 border border-green-200 rounded-xl p-4">
              <p class="text-sm text-green-700">
                ✅ Phê duyệt chênh lệch này sẽ xác nhận việc bàn giao và đóng yêu cầu.
              </p>
            </div>

            <div v-else class="bg-red-50 border border-red-200 rounded-xl p-4">
              <p class="text-sm text-red-700">
                ❌ Từ chối chênh lệch này sẽ yêu cầu xử lý lại hoặc điều tra thêm.
              </p>
            </div>

            <div class="flex gap-2">
              <button type="button" @click="showApprovalForm = false" 
                class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
                Hủy
              </button>
              <button type="submit" :disabled="approvalLoading"
                :class="approvalDecision ? 'bg-green-500 hover:bg-green-600' : 'bg-red-500 hover:bg-red-600'"
                class="flex-1 text-white px-4 py-3 rounded-xl font-medium disabled:opacity-50">
                {{ approvalLoading ? 'Đang xử lý...' : (approvalDecision ? 'Phê duyệt' : 'Từ chối') }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>

    <!-- Custom Date Range Modal -->
    <transition name="slide-up">
      <div v-if="showDateRangeForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6">
          <h3 class="text-xl font-bold mb-4">📅 Chọn khoảng thời gian</h3>
          <form @submit.prevent="applyCustomDateRange" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Từ ngày</label>
              <input v-model="customDateRange.start" type="date" required 
                class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500">
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Đến ngày</label>
              <input v-model="customDateRange.end" type="date" required 
                class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500">
            </div>
            <div class="flex gap-2">
              <button type="button" @click="showDateRangeForm = false" 
                class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
                Hủy
              </button>
              <button type="submit" 
                class="flex-1 bg-blue-500 hover:bg-blue-600 text-white px-4 py-3 rounded-xl font-medium">
                Áp dụng
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useManagerStore } from '../stores/manager'
import BottomNav from '../components/BottomNav.vue'

const managerStore = useManagerStore()

// Form states
const showApprovalForm = ref(false)
const showDateRangeForm = ref(false)
const selectedHandover = ref(null)
const approvalDecision = ref(true) // true = approve, false = reject
const statsDateRange = ref('month')

const approvalForm = ref({
  manager_notes: ''
})

const customDateRange = ref({
  start: '',
  end: ''
})

// Computed properties
const pendingApprovals = computed(() => managerStore.pendingApprovals)
const discrepancyStats = computed(() => managerStore.discrepancyStats)
const pendingApprovalCount = computed(() => managerStore.pendingApprovalCount)
const totalDiscrepancyAmount = computed(() => managerStore.totalDiscrepancyAmount)
const hasDiscrepancyStats = computed(() => managerStore.hasDiscrepancyStats)
const approvalLoading = computed(() => managerStore.approvalLoading)
const loading = computed(() => managerStore.loading)
const approvalError = computed(() => managerStore.approvalError)

// Methods
const refreshApprovals = async () => {
  await Promise.all([
    managerStore.fetchPendingApprovals(),
    loadDiscrepancyStats()
  ])
}

const showApprovalModal = (handover, approved) => {
  selectedHandover.value = handover
  approvalDecision.value = approved
  approvalForm.value = { manager_notes: '' }
  showApprovalForm.value = true
}

const approveDiscrepancy = async () => {
  try {
    await managerStore.approveDiscrepancy(
      selectedHandover.value.id,
      approvalDecision.value,
      approvalForm.value.manager_notes
    )
    showApprovalForm.value = false
    selectedHandover.value = null
    await refreshApprovals()
  } catch (error) {
    console.error('Approval failed:', error)
  }
}

const loadDiscrepancyStats = async () => {
  const { startDate, endDate } = getDateRange()
  await managerStore.getDiscrepancyStats(startDate, endDate)
}

const getDateRange = () => {
  const now = new Date()
  let startDate, endDate

  switch (statsDateRange.value) {
    case 'today':
      startDate = endDate = now.toISOString().split('T')[0]
      break
    case 'week':
      const weekStart = new Date(now)
      weekStart.setDate(now.getDate() - now.getDay())
      startDate = weekStart.toISOString().split('T')[0]
      endDate = now.toISOString().split('T')[0]
      break
    case 'month':
      startDate = new Date(now.getFullYear(), now.getMonth(), 1).toISOString().split('T')[0]
      endDate = now.toISOString().split('T')[0]
      break
    case 'custom':
      if (!showDateRangeForm.value) {
        showDateRangeForm.value = true
        return { startDate: null, endDate: null }
      }
      startDate = customDateRange.value.start
      endDate = customDateRange.value.end
      break
    default:
      startDate = endDate = now.toISOString().split('T')[0]
  }

  return { startDate, endDate }
}

const applyCustomDateRange = () => {
  showDateRangeForm.value = false
  loadDiscrepancyStats()
}

const clearApprovalError = () => {
  managerStore.clearApprovalError()
}

// Utility functions
const formatPrice = (amount) => {
  return new Intl.NumberFormat('vi-VN', { 
    style: 'currency', 
    currency: 'VND',
    maximumFractionDigits: 0
  }).format(amount)
}

const formatDateTime = (date) => {
  return new Date(date).toLocaleString('vi-VN', {
    day: '2-digit',
    month: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getResponsibilityText = (responsibility) => {
  const responsibilities = {
    WAITER: 'Phục vụ',
    CASHIER: 'Thu ngân',
    SYSTEM: 'Hệ thống',
    UNKNOWN: 'Chưa rõ'
  }
  return responsibilities[responsibility] || responsibility
}

const getNetDiscrepancyClass = (amount) => {
  if (amount > 0) return 'text-green-600'
  if (amount < 0) return 'text-red-600'
  return 'text-gray-600'
}

// Watch for date range changes
watch(statsDateRange, (newValue) => {
  if (newValue !== 'custom') {
    loadDiscrepancyStats()
  }
})

// Lifecycle
onMounted(async () => {
  await refreshApprovals()
})
</script>

<style scoped>
.active\:scale-95:active {
  transform: scale(0.95);
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