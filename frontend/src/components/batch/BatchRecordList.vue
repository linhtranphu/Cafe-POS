<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3" style="padding-top: max(0.75rem, env(safe-area-inset-top))">
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-3">
            <button @click="$router.push('/batch')" class="text-gray-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            <h1 class="text-xl font-bold text-gray-800">📦 Batch Records</h1>
          </div>
          <button 
            @click="showFilters = !showFilters"
            class="text-2xl">
            🔍
          </button>
        </div>
        
        <!-- Filters Panel -->
        <div v-if="showFilters" class="space-y-2 mb-3">
          <!-- Batch Type Filter -->
          <select 
            v-model="localFilters.batch_definition_id"
            @change="applyFilters"
            class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg">
            <option value="">Tất cả loại batch</option>
            <option 
              v-for="def in definitions" 
              :key="def.id" 
              :value="def.id">
              {{ def.name }}
            </option>
          </select>

          <!-- Status Filter -->
          <select 
            v-model="localFilters.status"
            @change="applyFilters"
            class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg">
            <option value="">Tất cả trạng thái</option>
            <option value="available">Khả dụng</option>
            <option value="expired">Hết hạn</option>
            <option value="depleted">Đã hết</option>
          </select>

          <!-- Date Range -->
          <div class="grid grid-cols-2 gap-2">
            <input 
              v-model="localFilters.from_date"
              @change="applyFilters"
              type="date"
              class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg"
            />
            <input 
              v-model="localFilters.to_date"
              @change="applyFilters"
              type="date"
              class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg"
            />
          </div>

          <!-- Clear Filters -->
          <button 
            @click="clearFilters"
            class="w-full bg-gray-200 text-gray-700 py-2 rounded-lg text-xs font-bold">
            Xóa bộ lọc
          </button>
        </div>

        <!-- Sort Options -->
        <div class="flex gap-2 overflow-x-auto pb-2">
          <button 
            @click="setSort('expires_at')"
            :class="sortBy === 'expires_at' ? 'bg-blue-500 text-white' : 'bg-white text-gray-700'"
            class="px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap border">
            ⏰ Hết hạn
          </button>
          <button 
            @click="setSort('prepared_at')"
            :class="sortBy === 'prepared_at' ? 'bg-blue-500 text-white' : 'bg-white text-gray-700'"
            class="px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap border">
            📅 Chế biến
          </button>
          <button 
            @click="setSort('quantity_remaining')"
            :class="sortBy === 'quantity_remaining' ? 'bg-blue-500 text-white' : 'bg-white text-gray-700'"
            class="px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap border">
            📊 Số lượng
          </button>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24">
      <!-- Create Button -->
      <button 
        @click="$router.push('/batch/records/create')"
        class="w-full bg-gradient-to-br from-green-500 to-emerald-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform mb-4">
        <div class="text-2xl mb-1">➕</div>
        <div class="text-sm font-bold">Ghi Nhận Batch Mới</div>
      </button>

      <!-- Loading State -->
      <LoadingSkeleton 
        v-if="loading && records.length === 0"
        type="list"
        :rows="5"
      />

      <!-- Error State -->
      <ErrorState
        v-else-if="error && records.length === 0"
        :icon="errorIcon"
        :title="errorTitle"
        :message="error"
        :retryable="isRetryable"
        :onRetry="() => batchRecordStore.fetchRecords()"
        showGoBack
        goBackRoute="/batch"
      />

      <!-- Empty State -->
      <EmptyState
        v-else-if="!loading && records.length === 0"
        icon="📦"
        title="Chưa có batch record"
        description="Bạn chưa ghi nhận batch nào. Hãy tạo batch record đầu tiên!"
        actionLabel="Ghi nhận batch mới"
        actionIcon="➕"
        :onAction="() => $router.push('/batch/records/create')"
      />

      <!-- Records List -->
      <div v-if="loading" class="text-center py-16">
        <div class="inline-block w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
        <p class="text-gray-500 mt-4">Đang tải...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="sortedRecords.length === 0" class="text-center py-16">
        <div class="text-6xl mb-4">📭</div>
        <p class="text-gray-500">Không có batch record nào</p>
      </div>

      <!-- Records List -->
      <div v-else class="space-y-3">
        <div 
          v-for="record in paginatedRecords" 
          :key="record.id"
          :class="getRecordColorClass(record)"
          class="rounded-2xl p-4 shadow-sm">
          
          <!-- Header -->
          <div class="flex justify-between items-start mb-3">
            <div>
              <h3 class="font-bold text-lg">{{ record.batch_name }}</h3>
              <div class="flex items-center gap-2 mt-1">
                <span :class="getStatusBadgeClass(record)" class="text-xs px-2 py-1 rounded-full font-medium">
                  {{ getStatusText(record) }}
                </span>
              </div>
            </div>
          </div>

          <!-- Info -->
          <div class="mb-3 space-y-2 text-sm">
            <div class="flex justify-between items-center">
              <span class="text-gray-600">📊 Còn lại:</span>
              <span class="font-bold">{{ record.quantity_remaining }} {{ record.unit }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-gray-600">🏭 Đã sản xuất:</span>
              <span class="font-medium">{{ record.quantity_produced }} {{ record.unit }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-gray-600">⏰ Hết hạn:</span>
              <span :class="getExpiryClass(record)" class="font-medium">
                {{ formatDateTime(record.expires_at) }}
              </span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-gray-600">👤 Người chế biến:</span>
              <span class="font-medium">{{ record.prepared_by_name || record.prepared_by }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-gray-600">💰 Chi phí:</span>
              <span class="font-medium">{{ formatCurrency(record.total_cost) }}</span>
            </div>
          </div>

          <!-- Actions -->
          <div class="grid grid-cols-3 gap-2 pt-3 border-t">
            <button 
              @click="viewDetail(record)"
              :class="getActionButtonColors('view')"
              class="py-2 rounded-lg text-xs font-bold">
              👁️ Xem
            </button>
            <button 
              v-if="record.status === 'available'"
              @click="markExpired(record)"
              :class="getActionButtonColors('expire')"
              class="py-2 rounded-lg text-xs font-bold">
              ⚠️ Hết hạn
            </button>
            <button 
              @click="deleteRecord(record)"
              :disabled="record.quantity_remaining < record.quantity_produced"
              :class="getActionButtonColors('delete')"
              class="py-2 rounded-lg text-xs font-bold disabled:opacity-50 disabled:cursor-not-allowed">
              🗑️ Xóa
            </button>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="mt-6 flex justify-center items-center gap-2">
        <button 
          @click="prevPage"
          :disabled="currentPage === 1"
          class="px-4 py-2 bg-white rounded-lg border disabled:opacity-50">
          ←
        </button>
        <span class="text-sm text-gray-600">
          Trang {{ currentPage }} / {{ totalPages }}
        </span>
        <button 
          @click="nextPage"
          :disabled="currentPage === totalPages"
          class="px-4 py-2 bg-white rounded-lg border disabled:opacity-50">
          →
        </button>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useBatchRecordStore } from '../../stores/batchRecord'
import { useBatchDefinitionStore } from '../../stores/batchDefinition'
import { useBatchColors } from '../../composables/useBatchColors'
import { useBatchErrors } from '../../composables/useBatchErrors'
import BottomNav from '../BottomNav.vue'
import ErrorState from './ErrorState.vue'
import EmptyState from './EmptyState.vue'
import LoadingSkeleton from './LoadingSkeleton.vue'

const router = useRouter()
const batchRecordStore = useBatchRecordStore()
const batchDefinitionStore = useBatchDefinitionStore()
const { getBatchRecordColors, getExpiryTextColor, getActionButtonColors } = useBatchColors()
const { parseError, getErrorIcon } = useBatchErrors()

const showFilters = ref(false)
const sortBy = ref('expires_at')
const sortOrder = ref('asc')

const localFilters = ref({
  batch_definition_id: '',
  status: '',
  from_date: '',
  to_date: ''
})

const loading = computed(() => batchRecordStore.loading)
const records = computed(() => batchRecordStore.records)
const definitions = computed(() => batchDefinitionStore.definitions)
const currentPage = computed(() => batchRecordStore.pagination.page)
const totalPages = computed(() => Math.ceil(batchRecordStore.pagination.total / batchRecordStore.pagination.limit))

// Error handling
const error = computed(() => batchRecordStore.error)
const parsedError = computed(() => {
  if (!error.value) return null
  return parseError({ response: { data: { message: error.value } } })
})
const errorIcon = computed(() => parsedError.value ? getErrorIcon(parsedError.value.type) : '❌')
const errorTitle = computed(() => {
  if (!parsedError.value) return 'Lỗi'
  const titles = {
    network: 'Lỗi kết nối',
    validation: 'Dữ liệu không hợp lệ',
    permission: 'Không có quyền',
    not_found: 'Không tìm thấy',
    server: 'Lỗi máy chủ',
    unknown: 'Lỗi không xác định'
  }
  return titles[parsedError.value.type] || 'Lỗi'
})
const isRetryable = computed(() => {
  if (!parsedError.value) return true
  return ['network', 'server'].includes(parsedError.value.type)
})


// Sort records
const sortedRecords = computed(() => {
  const sorted = [...records.value]
  
  sorted.sort((a, b) => {
    let aVal, bVal
    
    if (sortBy.value === 'expires_at') {
      aVal = new Date(a.expires_at).getTime()
      bVal = new Date(b.expires_at).getTime()
    } else if (sortBy.value === 'prepared_at') {
      aVal = new Date(a.prepared_at).getTime()
      bVal = new Date(b.prepared_at).getTime()
    } else if (sortBy.value === 'quantity_remaining') {
      aVal = a.quantity_remaining
      bVal = b.quantity_remaining
    }
    
    return sortOrder.value === 'asc' ? aVal - bVal : bVal - aVal
  })
  
  return sorted
})

// Paginated records (client-side pagination for sorted results)
const paginatedRecords = computed(() => {
  return sortedRecords.value
})

const getRecordColorClass = (record) => {
  const colors = getBatchRecordColors(record)
  return `${colors.background} border-2 ${colors.border}`
}

const getStatusBadgeClass = (record) => {
  const colors = getBatchRecordColors(record)
  return colors.badge
}

const getStatusText = (record) => {
  const colors = getBatchRecordColors(record)
  return colors.statusText
}

const getExpiryClass = (record) => {
  return getExpiryTextColor(record.expires_at)
}

const formatDateTime = (dateStr) => {
  const date = new Date(dateStr)
  return new Intl.DateTimeFormat('vi-VN', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

const formatCurrency = (value) => {
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND'
  }).format(value)
}

const setSort = (field) => {
  if (sortBy.value === field) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortBy.value = field
    sortOrder.value = field === 'expires_at' ? 'asc' : 'desc'
  }
}

const applyFilters = () => {
  const filters = {}
  
  if (localFilters.value.batch_definition_id) {
    filters.batch_definition_id = localFilters.value.batch_definition_id
  }
  if (localFilters.value.status) {
    filters.status = localFilters.value.status
  }
  if (localFilters.value.from_date) {
    filters.from_date = new Date(localFilters.value.from_date).toISOString()
  }
  if (localFilters.value.to_date) {
    filters.to_date = new Date(localFilters.value.to_date).toISOString()
  }
  
  batchRecordStore.setFilters(filters)
  batchRecordStore.fetchRecords()
}

const clearFilters = () => {
  localFilters.value = {
    batch_definition_id: '',
    status: '',
    from_date: '',
    to_date: ''
  }
  batchRecordStore.clearFilters()
  batchRecordStore.fetchRecords()
}

const viewDetail = (record) => {
  router.push(`/batch/records/${record.id}`)
}

const markExpired = async (record) => {
  if (!confirm(`Đánh dấu batch "${record.batch_name}" là hết hạn?`)) return
  
  try {
    await batchRecordStore.markAsExpired(record.id)
    alert('Đã đánh dấu hết hạn thành công')
  } catch (error) {
    alert(batchRecordStore.error || 'Lỗi đánh dấu hết hạn')
  }
}

const deleteRecord = async (record) => {
  if (record.quantity_remaining < record.quantity_produced) {
    alert('Không thể xóa batch đã được sử dụng một phần')
    return
  }
  
  if (!confirm(`Xóa batch record "${record.batch_name}"?`)) return
  
  try {
    await batchRecordStore.deleteRecord(record.id)
    alert('Đã xóa thành công')
  } catch (error) {
    alert(batchRecordStore.error || 'Lỗi xóa batch record')
  }
}

const prevPage = () => {
  if (currentPage.value > 1) {
    batchRecordStore.setPage(currentPage.value - 1)
    batchRecordStore.fetchRecords()
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    batchRecordStore.setPage(currentPage.value + 1)
    batchRecordStore.fetchRecords()
  }
}

onMounted(() => {
  batchRecordStore.fetchRecords()
  batchDefinitionStore.fetchDefinitions()
})
</script>
