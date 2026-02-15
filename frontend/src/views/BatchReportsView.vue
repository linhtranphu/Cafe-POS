<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3" style="padding-top: max(0.75rem, env(safe-area-inset-top))">
        <div class="flex items-center gap-3">
          <button @click="$router.push('/batch')" class="text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <div>
            <h1 class="text-xl font-bold text-gray-800">📊 Báo cáo Batch</h1>
            <p class="text-xs text-gray-500">Thống kê và phân tích</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24 space-y-4">
      <BatchProductionReport />
      <BatchWastageReport />
      <BatchUsageReport />
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />
  </div>
</template>

<script setup>
import BatchProductionReport from '../components/batch/BatchProductionReport.vue'
import BatchWastageReport from '../components/batch/BatchWastageReport.vue'
import BatchUsageReport from '../components/batch/BatchUsageReport.vue'
import BottomNav from '../components/BottomNav.vue'
</script>
