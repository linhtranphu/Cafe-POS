<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Mobile Header -->
    <div class="sticky top-0 z-40 bg-gradient-to-r from-purple-500 to-indigo-500 shadow-lg flex-shrink-0">
      <div class="px-4 py-4" style="padding-top: max(1rem, env(safe-area-inset-top)); padding-right: max(1rem, env(safe-area-inset-right)); padding-left: max(1rem, env(safe-area-inset-left))">
        <div class="flex items-center justify-between text-white">
          <div class="flex items-center gap-3">
            <button @click="$router.push('/dashboard')" class="text-white">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            <div>
              <h1 class="text-xl font-bold">🧪 Batch Management</h1>
              <p class="text-xs opacity-90">Nguyên liệu chế biến sẵn</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24" style="padding-bottom: max(6rem, calc(6rem + env(safe-area-inset-bottom))); padding-right: max(1rem, env(safe-area-inset-right)); padding-left: max(1rem, env(safe-area-inset-left))">
      <!-- Batch Status Widget -->
      <div class="mb-6">
        <BatchStatusWidget />
      </div>

      <!-- Quick Action Cards -->
      <div class="grid grid-cols-2 gap-3 mb-6">
        <button 
          @click="$router.push('/batch/records/create')"
          class="bg-gradient-to-br from-green-500 to-emerald-500 text-white rounded-xl p-6 shadow-md active:scale-95 transition-transform">
          <div class="text-4xl mb-2">➕</div>
          <div class="font-bold text-base">Tạo Batch</div>
          <div class="text-xs opacity-80 mt-1">Ghi nhận mới</div>
        </button>
        <button 
          @click="$router.push('/batch/definitions/create')"
          class="bg-gradient-to-br from-blue-500 to-cyan-500 text-white rounded-xl p-6 shadow-md active:scale-95 transition-transform">
          <div class="text-4xl mb-2">📋</div>
          <div class="font-bold text-base">Định Nghĩa</div>
          <div class="text-xs opacity-80 mt-1">Tạo công thức</div>
        </button>
      </div>

      <!-- Navigation Cards -->
      <div class="space-y-3 mb-6">
        <button 
          @click="$router.push('/batch/records')"
          class="w-full bg-white rounded-xl p-4 shadow-sm active:scale-98 transition-transform flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div class="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center text-2xl">
              🧪
            </div>
            <div class="text-left">
              <div class="font-bold text-gray-800">Batch Records</div>
              <div class="text-xs text-gray-500">Quản lý batch đã tạo</div>
            </div>
          </div>
          <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>

        <button 
          @click="$router.push('/batch/definitions')"
          class="w-full bg-white rounded-xl p-4 shadow-sm active:scale-98 transition-transform flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div class="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center text-2xl">
              📋
            </div>
            <div class="text-left">
              <div class="font-bold text-gray-800">Batch Definitions</div>
              <div class="text-xs text-gray-500">Công thức chế biến</div>
            </div>
          </div>
          <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>

        <button 
          @click="$router.push('/batch/alerts')"
          class="w-full bg-white rounded-xl p-4 shadow-sm active:scale-98 transition-transform flex items-center justify-between relative">
          <div class="flex items-center gap-3">
            <div class="w-12 h-12 bg-orange-100 rounded-xl flex items-center justify-center text-2xl">
              ⚠️
            </div>
            <div class="text-left">
              <div class="font-bold text-gray-800">Cảnh Báo</div>
              <div class="text-xs text-gray-500">Hết hạn & tồn kho thấp</div>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <span 
              v-if="totalAlerts > 0"
              class="bg-red-500 text-white text-xs px-2 py-1 rounded-full font-bold">
              {{ totalAlerts }}
            </span>
            <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
        </button>

        <button 
          @click="$router.push('/batch/reports')"
          class="w-full bg-white rounded-xl p-4 shadow-sm active:scale-98 transition-transform flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div class="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center text-2xl">
              📊
            </div>
            <div class="text-left">
              <div class="font-bold text-gray-800">Báo Cáo</div>
              <div class="text-xs text-gray-500">Sản xuất & sử dụng</div>
            </div>
          </div>
          <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </div>

      <!-- Info Card -->
      <div class="bg-gradient-to-br from-blue-50 to-purple-50 rounded-xl p-4 border-2 border-blue-200">
        <div class="flex items-start gap-3">
          <div class="text-2xl">💡</div>
          <div class="flex-1">
            <div class="font-bold text-gray-800 mb-1">Về Batch Management</div>
            <div class="text-sm text-gray-600 leading-relaxed">
              Quản lý nguyên liệu chế biến sẵn như sốt, nước ép, topping. Theo dõi hạn sử dụng, tồn kho và chi phí sản xuất.
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
import { useBatchAlertStore } from '../stores/batchAlert'
import BatchStatusWidget from '../components/batch/BatchStatusWidget.vue'
import BottomNav from '../components/BottomNav.vue'

const alertStore = useBatchAlertStore()

const totalAlerts = computed(() => alertStore.totalAlertCount)

onMounted(() => {
  alertStore.fetchAlerts()
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
