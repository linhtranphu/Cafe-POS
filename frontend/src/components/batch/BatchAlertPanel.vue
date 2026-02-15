<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3" style="padding-top: max(0.75rem, env(safe-area-inset-top))">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
            <button @click="$router.push('/batch')" class="text-gray-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            <h1 class="text-xl font-bold text-gray-800">🔔 Cảnh Báo Batch</h1>
          </div>
          <button 
            @click="refreshAlerts"
            :disabled="loading"
            class="p-2 rounded-lg bg-blue-500 text-white active:bg-blue-600 disabled:opacity-50">
            <span class="text-lg">🔄</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24">
      <!-- Loading State -->
      <div v-if="loading" class="text-center py-16">
        <div class="inline-block w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
        <p class="text-gray-500 mt-4">Đang tải cảnh báo...</p>
      </div>

      <!-- Error State -->
      <ErrorState
        v-else-if="error"
        icon="⚠️"
        title="Không thể tải cảnh báo"
        :message="error"
        :retryable="true"
        :onRetry="loadAlerts"
        size="medium"
        variant="inline"
      />

      <!-- No Alerts -->
      <div v-else-if="!hasAlerts" class="text-center py-16">
        <div class="text-6xl mb-4">✅</div>
        <p class="text-gray-500">Không có cảnh báo nào</p>
        <p class="text-sm text-gray-400 mt-2">Tất cả batch đều ổn định</p>
      </div>

      <!-- Alerts Sections -->
      <div v-else class="space-y-4">
        <!-- Expired Alerts -->
        <div v-if="expiredCount > 0" class="bg-white rounded-2xl shadow-sm overflow-hidden">
          <button 
            @click="toggleSection('expired')"
            class="w-full px-4 py-3 flex items-center justify-between bg-red-50 active:bg-red-100">
            <div class="flex items-center gap-3">
              <span class="text-2xl">🚫</span>
              <div class="text-left">
                <h3 class="font-bold text-red-700">Đã Hết Hạn</h3>
                <p class="text-xs text-red-600">Batch không thể sử dụng</p>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <span class="bg-red-500 text-white text-xs font-bold px-2 py-1 rounded-full">
                {{ expiredCount }}
              </span>
              <span class="text-red-700">{{ expandedSections.expired ? '▼' : '▶' }}</span>
            </div>
          </button>
          
          <div v-if="expandedSections.expired" class="p-4 space-y-3">
            <BatchAlertCard
              v-for="alert in alerts.expired"
              :key="alert.batch_record_id"
              :alert="alert"
              type="expired"
            />
          </div>
        </div>

        <!-- Expiring Soon Alerts -->
        <div v-if="expiringCount > 0" class="bg-white rounded-2xl shadow-sm overflow-hidden">
          <button 
            @click="toggleSection('expiring')"
            class="w-full px-4 py-3 flex items-center justify-between bg-yellow-50 active:bg-yellow-100">
            <div class="flex items-center gap-3">
              <span class="text-2xl">⏰</span>
              <div class="text-left">
                <h3 class="font-bold text-yellow-700">Sắp Hết Hạn</h3>
                <p class="text-xs text-yellow-600">Cần sử dụng sớm</p>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <span class="bg-yellow-500 text-white text-xs font-bold px-2 py-1 rounded-full">
                {{ expiringCount }}
              </span>
              <span class="text-yellow-700">{{ expandedSections.expiring ? '▼' : '▶' }}</span>
            </div>
          </button>
          
          <div v-if="expandedSections.expiring" class="p-4 space-y-3">
            <BatchAlertCard
              v-for="alert in alerts.expiring"
              :key="alert.batch_record_id"
              :alert="alert"
              type="expiring"
            />
          </div>
        </div>

        <!-- Low Stock Alerts -->
        <div v-if="lowStockCount > 0" class="bg-white rounded-2xl shadow-sm overflow-hidden">
          <button 
            @click="toggleSection('lowStock')"
            class="w-full px-4 py-3 flex items-center justify-between bg-orange-50 active:bg-orange-100">
            <div class="flex items-center gap-3">
              <span class="text-2xl">📉</span>
              <div class="text-left">
                <h3 class="font-bold text-orange-700">Tồn Kho Thấp</h3>
                <p class="text-xs text-orange-600">Cần chế biến thêm</p>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <span class="bg-orange-500 text-white text-xs font-bold px-2 py-1 rounded-full">
                {{ lowStockCount }}
              </span>
              <span class="text-orange-700">{{ expandedSections.lowStock ? '▼' : '▶' }}</span>
            </div>
          </button>
          
          <div v-if="expandedSections.lowStock" class="p-4 space-y-3">
            <BatchAlertCard
              v-for="alert in alerts.low_stock"
              :key="alert.batch_definition_id"
              :alert="alert"
              type="low_stock"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useBatchAlertStore } from '../../stores/batchAlert'
import BottomNav from '../BottomNav.vue'
import BatchAlertCard from './BatchAlertCard.vue'
import ErrorState from './ErrorState.vue'

const alertStore = useBatchAlertStore()

const expandedSections = ref({
  expired: true,
  expiring: true,
  lowStock: true
})

const loading = computed(() => alertStore.loading)
const error = computed(() => alertStore.error)
const alerts = computed(() => alertStore.alerts)
const hasAlerts = computed(() => alertStore.hasAlerts)
const lowStockCount = computed(() => alertStore.lowStockCount)
const expiringCount = computed(() => alertStore.expiringCount)
const expiredCount = computed(() => alertStore.expiredCount)
const lastChecked = computed(() => alertStore.lastChecked)

const toggleSection = (section) => {
  expandedSections.value[section] = !expandedSections.value[section]
}

const loadAlerts = async () => {
  try {
    await alertStore.fetchAlerts()
  } catch (err) {
    console.error('Failed to load alerts:', err)
  }
}

const refreshAlerts = async () => {
  await loadAlerts()
}

const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  const now = new Date()
  const diff = Math.floor((now - date) / 1000) // seconds
  
  if (diff < 60) return 'Vừa xong'
  if (diff < 3600) return `${Math.floor(diff / 60)} phút trước`
  if (diff < 86400) return `${Math.floor(diff / 3600)} giờ trước`
  return date.toLocaleDateString('vi-VN')
}

onMounted(() => {
  // Fetch alerts on mount
  alertStore.fetchAlerts()
  
  // Start auto-refresh every 5 minutes
  alertStore.startAutoRefresh(300000)
})

onUnmounted(() => {
  // Stop auto-refresh when component unmounts
  alertStore.stopAutoRefresh()
})
</script>
