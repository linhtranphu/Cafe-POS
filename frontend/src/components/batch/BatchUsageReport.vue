<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3" style="padding-top: max(0.75rem, env(safe-area-inset-top))">
        <h1 class="text-xl font-bold text-gray-800">📈 Báo Cáo Sử Dụng</h1>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24">
      <!-- Date Range Picker -->
      <div class="bg-white rounded-2xl p-4 shadow-sm mb-4">
        <h3 class="font-bold mb-3">Chọn Khoảng Thời Gian</h3>
        <div class="space-y-3">
          <div>
            <label class="text-sm text-gray-600 block mb-1">Từ ngày</label>
            <input
              v-model="fromDate"
              type="date"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="text-sm text-gray-600 block mb-1">Đến ngày</label>
            <input
              v-model="toDate"
              type="date"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <button
            @click="fetchReport"
            :disabled="!fromDate || !toDate || loading"
            class="w-full bg-blue-500 text-white py-3 rounded-lg font-bold active:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed">
            {{ loading ? 'Đang tải...' : 'Xem Báo Cáo' }}
          </button>
        </div>
      </div>

      <!-- Error State -->
      <InlineError
        v-if="error"
        :message="error"
        :show="!!error"
        :showRetry="true"
        :onRetry="fetchReport"
        @dismiss="() => reportStore.clearError()"
      />

      <!-- Report Content -->
      <div v-if="report" class="space-y-4">
        <!-- Summary Cards -->
        <div class="grid grid-cols-2 gap-3">
          <div class="bg-gradient-to-br from-blue-500 to-cyan-500 rounded-2xl p-4 text-white shadow-md">
            <div class="text-2xl font-bold">{{ report.total_usage_count || 0 }}</div>
            <div class="text-xs opacity-90">Lần Sử Dụng</div>
          </div>
          <div class="bg-gradient-to-br from-green-500 to-emerald-500 rounded-2xl p-4 text-white shadow-md">
            <div class="text-2xl font-bold">{{ formatNumber(report.total_quantity_used) }}</div>
            <div class="text-xs opacity-90">Tổng Số Lượng</div>
          </div>
          <div class="bg-gradient-to-br from-purple-500 to-pink-500 rounded-2xl p-4 text-white shadow-md col-span-2">
            <div class="text-2xl font-bold">{{ formatCurrency(report.total_cost) }}</div>
            <div class="text-xs opacity-90">Tổng Chi Phí Sử Dụng</div>
          </div>
        </div>

        <!-- By Menu Item -->
        <div v-if="report.by_menu_item?.length" class="bg-white rounded-2xl p-4 shadow-sm">
          <h3 class="font-bold mb-3">Theo Món Ăn</h3>
          <div class="space-y-3">
            <div
              v-for="item in report.by_menu_item"
              :key="item.menu_item_name"
              class="border-l-4 border-blue-500 pl-3 py-2">
              <div class="font-bold text-sm">{{ item.menu_item_name }}</div>
              <div class="grid grid-cols-3 gap-2 mt-1 text-xs text-gray-600">
                <div>
                  <span class="block">Số lần</span>
                  <span class="font-bold text-gray-800">{{ item.count }}</span>
                </div>
                <div>
                  <span class="block">Số lượng</span>
                  <span class="font-bold text-gray-800">{{ formatNumber(item.quantity_used) }}</span>
                </div>
                <div>
                  <span class="block">Chi phí</span>
                  <span class="font-bold text-gray-800">{{ formatCurrency(item.total_cost) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Most Used Batches -->
        <div v-if="report.most_used_batches?.length" class="bg-white rounded-2xl p-4 shadow-sm">
          <h3 class="font-bold mb-3">🏆 Batch Được Dùng Nhiều Nhất</h3>
          <div class="space-y-3">
            <div
              v-for="(item, index) in report.most_used_batches"
              :key="item.batch_name"
              class="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
              <div class="text-2xl font-bold text-gray-400">{{ index + 1 }}</div>
              <div class="flex-1">
                <div class="font-bold text-sm">{{ item.batch_name }}</div>
                <div class="text-xs text-gray-600">
                  {{ item.count }} lần • {{ formatNumber(item.quantity_used) }} {{ item.unit }}
                </div>
              </div>
              <div class="text-right">
                <div class="font-bold text-sm">{{ formatCurrency(item.total_cost) }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- By Batch Type -->
        <div v-if="report.by_batch_type?.length" class="bg-white rounded-2xl p-4 shadow-sm">
          <h3 class="font-bold mb-3">Theo Loại Batch</h3>
          <div class="space-y-3">
            <div
              v-for="item in report.by_batch_type"
              :key="item.batch_name"
              class="border-l-4 border-green-500 pl-3 py-2">
              <div class="font-bold text-sm">{{ item.batch_name }}</div>
              <div class="grid grid-cols-3 gap-2 mt-1 text-xs text-gray-600">
                <div>
                  <span class="block">Số lần</span>
                  <span class="font-bold text-gray-800">{{ item.count }}</span>
                </div>
                <div>
                  <span class="block">Số lượng</span>
                  <span class="font-bold text-gray-800">{{ formatNumber(item.quantity_used) }}</span>
                </div>
                <div>
                  <span class="block">Chi phí</span>
                  <span class="font-bold text-gray-800">{{ formatCurrency(item.total_cost) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Export Button -->
        <button
          @click="exportReport"
          :disabled="loading"
          class="w-full bg-green-500 text-white py-3 rounded-lg font-bold active:bg-green-600 disabled:opacity-50 flex items-center justify-center gap-2">
          <span>📥</span>
          <span>Xuất CSV</span>
        </button>
      </div>

      <!-- Empty State -->
      <div v-else-if="!loading && fromDate && toDate" class="text-center py-16">
        <div class="text-6xl mb-4">📈</div>
        <p class="text-gray-500">Chọn khoảng thời gian và nhấn "Xem Báo Cáo"</p>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useBatchReportStore } from '../../stores/batchReport'
import BottomNav from '../BottomNav.vue'
import InlineError from './InlineError.vue'

const reportStore = useBatchReportStore()

const fromDate = ref('')
const toDate = ref('')

const loading = computed(() => reportStore.loading)
const error = computed(() => reportStore.error)
const report = computed(() => reportStore.usageReport)

const fetchReport = async () => {
  if (!fromDate.value || !toDate.value) return

  try {
    await reportStore.fetchUsageReport({
      from_date: new Date(fromDate.value).toISOString(),
      to_date: new Date(toDate.value).toISOString()
    })
  } catch (err) {
    console.error('Failed to fetch report:', err)
  }
}

const exportReport = async () => {
  try {
    await reportStore.exportReport('usage')
    alert('Đã xuất báo cáo thành công')
  } catch (err) {
    alert('Lỗi xuất báo cáo: ' + err.message)
  }
}

const formatNumber = (value) => {
  if (!value) return '0'
  return new Intl.NumberFormat('vi-VN').format(value)
}

const formatCurrency = (value) => {
  if (!value) return '0đ'
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND'
  }).format(value)
}
</script>
