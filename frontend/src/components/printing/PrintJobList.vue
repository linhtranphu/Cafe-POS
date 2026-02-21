<template>
  <div class="h-full flex flex-col bg-gray-50">
    <!-- Header with Filters -->
    <div class="bg-white shadow-sm flex-shrink-0 px-4 py-3">
      <h2 class="text-lg font-bold text-gray-800 mb-3">📄 Print Jobs</h2>
      
      <!-- Status Filter Tabs -->
      <div class="flex gap-2 overflow-x-auto pb-2">
        <button
          v-for="filter in statusFilters"
          :key="filter.value"
          @click="selectedStatus = filter.value"
          :class="[
            'px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-colors',
            selectedStatus === filter.value
              ? 'bg-blue-500 text-white'
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          ]"
        >
          {{ filter.label }}
          <span v-if="filter.count !== undefined" class="ml-1">
            ({{ filter.count }})
          </span>
        </button>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4">
      <!-- Loading State -->
      <div v-if="loading" class="text-center py-16">
        <div class="inline-block w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
        <p class="text-gray-500 mt-4">Đang tải...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-16">
        <div class="text-6xl mb-4">⚠️</div>
        <p class="text-red-500 mb-4">{{ error }}</p>
        <button
          @click="refreshJobs"
          class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
        >
          Thử lại
        </button>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredJobs.length === 0" class="text-center py-16">
        <div class="text-6xl mb-4">📭</div>
        <p class="text-gray-500">Không có print job nào</p>
      </div>

      <!-- Jobs List -->
      <div v-else class="space-y-3">
        <div
          v-for="job in filteredJobs"
          :key="job.id"
          class="bg-white rounded-2xl p-4 shadow-sm"
        >
          <!-- Header -->
          <div class="flex justify-between items-start mb-3">
            <div>
              <div class="flex items-center gap-2">
                <span class="text-2xl">{{ getJobIcon(job.type) }}</span>
                <h3 class="font-bold text-lg">{{ getJobTypeLabel(job.type) }}</h3>
              </div>
              <p class="text-sm text-gray-600">Order: {{ job.order_number }}</p>
            </div>
            <span
              :class="[
                'px-3 py-1 rounded-full text-xs font-bold',
                getStatusClass(job.status)
              ]"
            >
              {{ getStatusLabel(job.status) }}
            </span>
          </div>

          <!-- Info -->
          <div class="mb-3 space-y-2 text-sm">
            <div class="flex justify-between items-center">
              <span class="text-gray-600">🕐 Tạo lúc:</span>
              <span class="font-medium">{{ formatDateTime(job.created_at) }}</span>
            </div>
            <div v-if="job.printed_at" class="flex justify-between items-center">
              <span class="text-gray-600">✅ In lúc:</span>
              <span class="font-medium">{{ formatDateTime(job.printed_at) }}</span>
            </div>
            <div v-if="job.retry_count > 0" class="flex justify-between items-center">
              <span class="text-gray-600">🔄 Số lần thử:</span>
              <span class="font-medium">{{ job.retry_count }}/{{ job.max_retries }}</span>
            </div>
            <div v-if="job.error_msg" class="mt-2 p-2 bg-red-50 rounded-lg">
              <p class="text-xs text-red-600">
                <span class="font-bold">Lỗi:</span> {{ job.error_msg }}
              </p>
            </div>
          </div>

          <!-- Actions -->
          <div
            v-if="job.status === 'FAILED' || job.status === 'PENDING'"
            class="grid grid-cols-2 gap-2 pt-3 border-t"
          >
            <button
              v-if="job.status === 'FAILED'"
              @click="handleRetry(job.id)"
              :disabled="retrying"
              class="bg-green-500 text-white py-2 rounded-lg text-xs font-bold hover:bg-green-600 active:scale-95 transition-transform flex items-center justify-center gap-1 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span>🔄</span>
              <span>Thử lại</span>
            </button>
            <button
              v-if="job.status === 'PENDING'"
              @click="handleCancel(job.id)"
              :disabled="canceling"
              class="bg-red-500 text-white py-2 rounded-lg text-xs font-bold hover:bg-red-600 active:scale-95 transition-transform flex items-center justify-center gap-1 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span>❌</span>
              <span>Hủy</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { usePrintJobStore } from '../../stores/printJob'
import { usePrintJobWebSocket } from '../../composables/usePrintJobWebSocket'
import { useNotifications } from '../../composables/useNotifications'

const printJobStore = usePrintJobStore()

// Setup WebSocket listeners for real-time updates
const { setupListeners, cleanupListeners } = usePrintJobWebSocket()

// Setup notifications
const { showSuccess, showError } = useNotifications()

const selectedStatus = ref('all')
const retrying = ref(false)
const canceling = ref(false)
let refreshInterval = null

const statusFilters = computed(() => [
  { label: 'Tất cả', value: 'all' },
  { label: 'Đang chờ', value: 'PENDING', count: printJobStore.getPendingCount },
  { label: 'Thất bại', value: 'FAILED', count: printJobStore.getFailedCount },
  { label: 'Đang in', value: 'PRINTING' },
  { label: 'Hoàn thành', value: 'COMPLETED' }
])

const loading = computed(() => printJobStore.loading)
const error = computed(() => printJobStore.error)

const filteredJobs = computed(() => {
  if (selectedStatus.value === 'all') {
    return printJobStore.printJobs
  }
  if (selectedStatus.value === 'PENDING') {
    return printJobStore.pendingJobs
  }
  if (selectedStatus.value === 'FAILED') {
    return printJobStore.failedJobs
  }
  return printJobStore.jobsByStatus(selectedStatus.value)
})

const getJobIcon = (type) => {
  return type === 'BILL' ? '🧾' : '🏷️'
}

const getJobTypeLabel = (type) => {
  return type === 'BILL' ? 'Bill' : 'Tem'
}

const getStatusLabel = (status) => {
  const labels = {
    PENDING: 'Đang chờ',
    PRINTING: 'Đang in',
    COMPLETED: 'Hoàn thành',
    FAILED: 'Thất bại'
  }
  return labels[status] || status
}

const getStatusClass = (status) => {
  const classes = {
    PENDING: 'bg-yellow-100 text-yellow-800',
    PRINTING: 'bg-blue-100 text-blue-800',
    COMPLETED: 'bg-green-100 text-green-800',
    FAILED: 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const formatDateTime = (dateString) => {
  if (!dateString) return 'N/A'
  const date = new Date(dateString)
  return date.toLocaleString('vi-VN', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const refreshJobs = async () => {
  await printJobStore.fetchJobs()
  await printJobStore.fetchPendingJobs()
  await printJobStore.fetchFailedJobs()
}

const handleRetry = async (jobId) => {
  if (retrying.value) return
  
  retrying.value = true
  try {
    const success = await printJobStore.retryJob(jobId)
    if (success) {
      showSuccess('Đã gửi lại print job')
    } else {
      showError(printJobStore.error || 'Lỗi retry print job')
    }
  } finally {
    retrying.value = false
  }
}

const handleCancel = async (jobId) => {
  if (canceling.value) return
  
  if (!confirm('Bạn có chắc muốn hủy print job này?')) return
  
  canceling.value = true
  try {
    const success = await printJobStore.cancelJob(jobId)
    if (success) {
      showSuccess('Đã hủy print job')
    } else {
      showError(printJobStore.error || 'Lỗi hủy print job')
    }
  } finally {
    canceling.value = false
  }
}

onMounted(() => {
  refreshJobs()
  
  // Auto-refresh every 10 seconds for real-time updates
  refreshInterval = setInterval(() => {
    refreshJobs()
  }, 10000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>
