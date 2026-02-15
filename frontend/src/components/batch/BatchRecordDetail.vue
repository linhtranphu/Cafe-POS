<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3" style="padding-top: max(0.75rem, env(safe-area-inset-top))">
        <div class="flex items-center justify-between">
          <button @click="goBack" class="text-2xl text-gray-600">←</button>
          <h1 class="text-xl font-bold text-gray-800">Chi Tiết Batch</h1>
          <div class="w-8"></div>
        </div>
      </div>
    </div>

    <!-- Scrollable Content -->
    <div class="flex-1 overflow-y-auto">
      <div class="max-w-md mx-auto px-4 py-6 space-y-5">
        <!-- Loading State -->
        <div v-if="loading" class="text-center py-16">
          <div class="inline-block w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
          <p class="text-gray-500 mt-4">Đang tải...</p>
        </div>

        <!-- Error State -->
        <ErrorState
          v-else-if="error"
          icon="⚠️"
          title="Không thể tải batch record"
          :message="error"
          :retryable="true"
          :onRetry="loadData"
          showGoBack
          goBackRoute="/batch/records"
        />

        <!-- Batch Record Details -->
        <div v-else-if="record">
        <!-- Status Card -->
        <div :class="getRecordColorClass(record)" class="rounded-xl p-4 border-2">
          <div class="flex justify-between items-start mb-3">
            <div>
              <h2 class="text-2xl font-bold">{{ record.batch_name }}</h2>
              <span :class="getStatusBadgeClass(record)" class="inline-block text-xs px-3 py-1 rounded-full font-medium mt-2">
                {{ getStatusText(record.status) }}
              </span>
            </div>
          </div>

          <!-- Timeline Visualization -->
          <div class="mt-4 relative">
            <div class="flex items-center justify-between text-xs">
              <div class="text-center flex-1">
                <div class="text-2xl mb-1">🏭</div>
                <div class="font-medium">Chế biến</div>
                <div class="text-gray-600">{{ formatDate(record.prepared_at) }}</div>
              </div>
              
              <div class="flex-1 h-1 bg-gray-300 mx-2"></div>
              
              <div class="text-center flex-1">
                <div class="text-2xl mb-1">{{ record.status === 'depleted' ? '✅' : '📦' }}</div>
                <div class="font-medium">{{ record.status === 'depleted' ? 'Đã dùng hết' : 'Đang dùng' }}</div>
                <div class="text-gray-600">{{ getUsagePercentage }}%</div>
              </div>
              
              <div class="flex-1 h-1 bg-gray-300 mx-2"></div>
              
              <div class="text-center flex-1">
                <div class="text-2xl mb-1">{{ record.status === 'expired' ? '❌' : '⏰' }}</div>
                <div class="font-medium">Hết hạn</div>
                <div :class="getExpiryClass(record)" class="font-medium">
                  {{ formatDate(record.expires_at) }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Quantity Info -->
        <div class="bg-white rounded-xl p-4 border-2 border-gray-200">
          <div class="flex items-center gap-2 mb-3">
            <span class="text-xl">📊</span>
            <span class="text-sm font-semibold text-gray-800">Thông Tin Số Lượng</span>
          </div>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between items-center">
              <span class="text-gray-600">Đã sản xuất:</span>
              <span class="font-bold">{{ record.quantity_produced }} {{ record.unit }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-gray-600">Còn lại:</span>
              <span class="font-bold text-blue-600">{{ record.quantity_remaining }} {{ record.unit }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-gray-600">Đã sử dụng:</span>
              <span class="font-bold text-green-600">
                {{ (record.quantity_produced - record.quantity_remaining).toFixed(2) }} {{ record.unit }}
              </span>
            </div>
          </div>
        </div>

        <!-- Cost Breakdown -->
        <div class="bg-white rounded-xl p-4 border-2 border-gray-200">
          <div class="flex items-center gap-2 mb-3">
            <span class="text-xl">💰</span>
            <span class="text-sm font-semibold text-gray-800">Chi Phí</span>
          </div>
          <div class="space-y-2 text-sm mb-3">
            <div class="flex justify-between items-center">
              <span class="text-gray-600">Tổng chi phí:</span>
              <span class="font-bold text-blue-600">{{ formatCurrency(record.total_cost) }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-gray-600">Chi phí / {{ record.unit }}:</span>
              <span class="font-bold">{{ formatCurrency(record.cost_per_unit) }}</span>
            </div>
          </div>

          <!-- Ingredients Used -->
          <div v-if="record.ingredients_used && record.ingredients_used.length > 0" class="pt-3 border-t">
            <h4 class="text-xs font-semibold text-gray-700 mb-2">Nguyên liệu đã sử dụng:</h4>
            <div class="space-y-2">
              <div 
                v-for="(ingredient, index) in record.ingredients_used" 
                :key="index"
                class="bg-gray-50 rounded-lg p-2">
                <div class="flex justify-between items-start text-xs">
                  <div>
                    <div class="font-medium">{{ ingredient.ingredient_name }}</div>
                    <div class="text-gray-600">
                      {{ ingredient.quantity }} {{ ingredient.unit }} × {{ formatCurrency(ingredient.cost_per_unit) }}
                    </div>
                  </div>
                  <div class="font-bold text-blue-600">
                    {{ formatCurrency(ingredient.total_cost) }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Other Info -->
        <div class="bg-white rounded-xl p-4 border-2 border-gray-200">
          <div class="flex items-center gap-2 mb-3">
            <span class="text-xl">ℹ️</span>
            <span class="text-sm font-semibold text-gray-800">Thông Tin Khác</span>
          </div>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between items-center">
              <span class="text-gray-600">👤 Người chế biến:</span>
              <span class="font-medium">{{ preparedByName }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-gray-600">📅 Thời gian chế biến:</span>
              <span class="font-medium">{{ formatDateTime(record.prepared_at) }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-gray-600">⏰ Thời gian hết hạn:</span>
              <span :class="getExpiryClass(record)" class="font-medium">
                {{ formatDateTime(record.expires_at) }}
              </span>
            </div>
          </div>
        </div>

        <!-- Usage History -->
        <div class="bg-white rounded-xl p-4 border-2 border-gray-200">
          <div class="flex items-center gap-2 mb-3">
            <span class="text-xl">📝</span>
            <span class="text-sm font-semibold text-gray-800">Lịch Sử Sử Dụng</span>
          </div>
          
          <!-- Loading Usage History -->
          <div v-if="loadingUsage" class="text-center py-4">
            <div class="inline-block w-6 h-6 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
          </div>

          <!-- No Usage -->
          <div v-else-if="!usageHistory || usageHistory.length === 0" class="text-center py-4 text-gray-500 text-sm">
            Chưa có lịch sử sử dụng
          </div>

          <!-- Usage List -->
          <div v-else class="space-y-2">
            <div 
              v-for="(usage, index) in usageHistory" 
              :key="index"
              class="bg-gray-50 rounded-lg p-3">
              <div class="flex justify-between items-start mb-2">
                <div class="flex-1">
                  <div class="font-medium text-sm">{{ usage.menu_item_name }}</div>
                  <div class="text-xs text-gray-600">{{ formatDateTime(usage.used_at) }}</div>
                </div>
                <div class="text-right">
                  <div class="font-bold text-sm text-blue-600">
                    {{ usage.quantity_used }} {{ usage.unit }}
                  </div>
                  <div class="text-xs text-gray-600">
                    {{ formatCurrency(usage.total_cost) }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Warning Message -->
        <div v-if="record.quantity_remaining < record.quantity_produced" class="bg-yellow-50 border-2 border-yellow-300 rounded-xl p-4">
          <div class="flex items-start gap-2">
            <span class="text-xl">⚠️</span>
            <p class="text-sm text-yellow-800 flex-1">
              Không thể xóa batch đã được sử dụng một phần
            </p>
          </div>
        </div>

          <!-- Spacer for bottom buttons -->
          <div class="h-24"></div>
        </div>
      </div>
    </div>

    <!-- Fixed Footer with Actions -->
    <div class="flex-shrink-0 bg-white border-t">
      <div class="max-w-md mx-auto px-4 py-4 flex gap-3 pb-safe">
        <button 
          v-if="record && record.status === 'available'"
          @click="markExpired"
          class="flex-1 bg-yellow-500 text-white py-4 rounded-xl font-medium text-base active:bg-yellow-600 flex items-center justify-center gap-2">
          <span>⚠️</span>
          <span>Đánh dấu hết hạn</span>
        </button>
        <button 
          v-if="record"
          @click="deleteBatch"
          :disabled="record.quantity_remaining < record.quantity_produced"
          class="flex-1 bg-red-500 text-white py-4 rounded-xl font-medium text-base active:bg-red-600 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2">
          <span>🗑️</span>
          <span>Xóa</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useBatchRecordStore } from '../../stores/batchRecord'
import api from '../../services/api'
import ErrorState from './ErrorState.vue'

const router = useRouter()
const route = useRoute()
const batchRecordStore = useBatchRecordStore()

const loading = ref(false)
const loadingUsage = ref(false)
const error = ref(null)
const usageHistory = ref([])

const record = computed(() => batchRecordStore.currentRecord)

const preparedByName = computed(() => {
  if (!record.value) return ''
  
  // Use prepared_by_name if available, otherwise fall back to prepared_by
  return record.value.prepared_by_name || record.value.prepared_by
})

const getUsagePercentage = computed(() => {
  if (!record.value) return 0
  const used = record.value.quantity_produced - record.value.quantity_remaining
  return ((used / record.value.quantity_produced) * 100).toFixed(0)
})

const getRecordColorClass = (record) => {
  if (record.status === 'expired') return 'bg-red-50 border-2 border-red-200'
  if (record.status === 'depleted') return 'bg-gray-100 border-2 border-gray-300'
  
  const now = new Date()
  const expiresAt = new Date(record.expires_at)
  const hoursUntilExpiry = (expiresAt - now) / (1000 * 60 * 60)
  
  if (hoursUntilExpiry <= 4) return 'bg-yellow-50 border-2 border-yellow-300'
  
  return 'bg-green-50 border-2 border-green-200'
}

const getStatusBadgeClass = (record) => {
  if (record.status === 'expired') return 'bg-red-500 text-white'
  if (record.status === 'depleted') return 'bg-gray-500 text-white'
  return 'bg-green-500 text-white'
}

const getStatusText = (status) => {
  const statusMap = {
    available: 'Khả dụng',
    expired: 'Hết hạn',
    depleted: 'Đã hết'
  }
  return statusMap[status] || status
}

const getExpiryClass = (record) => {
  if (record.status === 'expired') return 'text-red-600 font-bold'
  
  const now = new Date()
  const expiresAt = new Date(record.expires_at)
  const hoursUntilExpiry = (expiresAt - now) / (1000 * 60 * 60)
  
  if (hoursUntilExpiry <= 4) return 'text-yellow-600 font-bold'
  
  return 'text-gray-700'
}

const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  return new Intl.DateTimeFormat('vi-VN', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  }).format(date)
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

const loadUsageHistory = async () => {
  if (!route.params.id) return
  
  loadingUsage.value = true
  
  try {
    const response = await api.get(`/batch-usage/history?batch_record_id=${route.params.id}`)
    usageHistory.value = response.data.data || []
  } catch (err) {
    console.error('Error loading usage history:', err)
    usageHistory.value = []
  } finally {
    loadingUsage.value = false
  }
}

const loadData = async () => {
  if (!route.params.id) {
    error.value = 'ID batch không hợp lệ'
    return
  }
  
  loading.value = true
  error.value = null
  
  try {
    await batchRecordStore.fetchRecordById(route.params.id)
    await loadUsageHistory()
  } catch (err) {
    error.value = batchRecordStore.error || err.message || 'Lỗi tải batch record'
  } finally {
    loading.value = false
  }
}

const markExpired = async () => {
  if (!record.value) return
  
  if (!confirm(`Đánh dấu batch "${record.value.batch_name}" là hết hạn?`)) return
  
  try {
    await batchRecordStore.markAsExpired(record.value.id)
    alert('Đã đánh dấu hết hạn thành công')
  } catch (err) {
    alert(batchRecordStore.error || 'Lỗi đánh dấu hết hạn')
  }
}

const deleteBatch = async () => {
  if (!record.value) return
  
  if (record.value.quantity_remaining < record.value.quantity_produced) {
    alert('Không thể xóa batch đã được sử dụng một phần')
    return
  }
  
  if (!confirm(`Xóa batch record "${record.value.batch_name}"? Hành động này không thể hoàn tác.`)) return
  
  try {
    await batchRecordStore.deleteRecord(record.value.id)
    alert('Đã xóa thành công')
    router.push('/batch/records')
  } catch (err) {
    alert(batchRecordStore.error || 'Lỗi xóa batch record')
  }
}

const goBack = () => {
  router.back()
}

onMounted(() => {
  loadData()
})
</script>
