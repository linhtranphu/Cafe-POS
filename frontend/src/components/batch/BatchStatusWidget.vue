<template>
  <div class="bg-white rounded-2xl p-4 shadow-sm">
    <!-- Header -->
    <div class="flex items-center justify-between mb-3">
      <h3 class="font-bold text-lg">🧪 Batch Status</h3>
      <router-link 
        to="/batch/records"
        class="text-xs text-blue-500 font-bold">
        Xem Tất Cả →
      </router-link>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-8">
      <div class="inline-block w-6 h-6 border-3 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <!-- Content -->
    <div v-else class="space-y-3">
      <!-- Alert Badges -->
      <div class="grid grid-cols-3 gap-2">
        <div class="bg-red-50 rounded-lg p-2 text-center">
          <div class="text-xl font-bold text-red-600">{{ expiredCount }}</div>
          <div class="text-xs text-red-600">Hết Hạn</div>
        </div>
        <div class="bg-yellow-50 rounded-lg p-2 text-center">
          <div class="text-xl font-bold text-yellow-600">{{ expiringCount }}</div>
          <div class="text-xs text-yellow-600">Sắp Hết</div>
        </div>
        <div class="bg-orange-50 rounded-lg p-2 text-center">
          <div class="text-xl font-bold text-orange-600">{{ lowStockCount }}</div>
          <div class="text-xs text-orange-600">Tồn Thấp</div>
        </div>
      </div>

      <!-- Summary Stats -->
      <div class="border-t pt-3 space-y-2 text-sm">
        <div class="flex justify-between">
          <span class="text-gray-600">Tổng batch khả dụng:</span>
          <span class="font-bold">{{ totalAvailable }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-600">Batch definitions:</span>
          <span class="font-bold">{{ totalDefinitions }}</span>
        </div>
      </div>

      <!-- Quick Links -->
      <div class="grid grid-cols-2 gap-2 pt-3 border-t">
        <router-link
          to="/batch/records/create"
          class="bg-blue-500 text-white py-2 rounded-lg text-xs font-bold text-center active:bg-blue-600">
          ➕ Tạo Batch
        </router-link>
        <router-link
          to="/batch/alerts"
          class="bg-orange-500 text-white py-2 rounded-lg text-xs font-bold text-center active:bg-orange-600 relative">
          🔔 Cảnh Báo
          <span 
            v-if="totalAlertCount > 0"
            class="absolute -top-1 -right-1 bg-red-500 text-white text-xs w-5 h-5 rounded-full flex items-center justify-center">
            {{ totalAlertCount }}
          </span>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useBatchAlertStore } from '../../stores/batchAlert'
import { useBatchRecordStore } from '../../stores/batchRecord'
import { useBatchDefinitionStore } from '../../stores/batchDefinition'

const alertStore = useBatchAlertStore()
const recordStore = useBatchRecordStore()
const definitionStore = useBatchDefinitionStore()

const loading = ref(true)

const expiredCount = computed(() => alertStore.expiredCount)
const expiringCount = computed(() => alertStore.expiringCount)
const lowStockCount = computed(() => alertStore.lowStockCount)
const totalAlertCount = computed(() => alertStore.totalAlertCount)

const totalAvailable = computed(() => {
  return recordStore.records.filter(r => r.status === 'available').length
})

const totalDefinitions = computed(() => {
  return definitionStore.definitions.length
})

onMounted(async () => {
  loading.value = true
  try {
    // Fetch data in parallel
    await Promise.all([
      alertStore.fetchAlerts(),
      recordStore.fetchRecords({ status: 'available', limit: 100 }),
      definitionStore.fetchDefinitions()
    ])
  } catch (err) {
    console.error('Failed to load batch status:', err)
  } finally {
    loading.value = false
  }
})
</script>
