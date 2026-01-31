<template>
  <div class="bg-white rounded-2xl p-4 shadow-sm mb-4">
    <h2 class="text-lg font-bold text-gray-800 mb-4">📅 Ca thu ngân</h2>
    
    <!-- Loading State -->
    <div v-if="loading" class="text-center py-4">
      <div class="animate-spin text-3xl">⏳</div>
      <p class="text-sm text-gray-600 mt-2">Đang tải...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-red-50 border-2 border-red-200 rounded-xl p-4">
      <div class="flex items-start gap-3">
        <span class="text-2xl">⚠️</span>
        <div>
          <p class="font-medium text-red-800">Lỗi</p>
          <p class="text-sm text-red-600">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- No Shift - Show Start Button -->
    <div v-else-if="!hasOpenShift" class="space-y-4">
      <div class="bg-gray-50 rounded-xl p-4 text-center">
        <div class="text-4xl mb-2">💼</div>
        <p class="text-sm text-gray-600 mb-4">Chưa có ca thu ngân nào đang mở</p>
        
        <button
          @click="showStartModal = true"
          class="w-full py-3 bg-yellow-500 text-white rounded-xl font-bold text-base active:scale-95 transition-transform"
        >
          ➕ Bắt đầu ca thu ngân
        </button>
      </div>
    </div>

    <!-- Has Open Shift - Show Details -->
    <div v-else class="space-y-3">
      <div class="bg-gradient-to-r from-yellow-500 to-orange-500 text-white rounded-xl p-4">
        <div class="flex items-center justify-between mb-3">
          <div>
            <p class="text-xs opacity-90">Ca hiện tại</p>
            <p class="text-lg font-bold">{{ currentShift.cashier_name }}</p>
          </div>
          <div class="bg-white/20 rounded-lg px-3 py-1 backdrop-blur-sm">
            <p class="text-xs font-medium">{{ getStatusText(currentShift.status) }}</p>
          </div>
        </div>
        
        <div class="grid grid-cols-2 gap-3">
          <div class="bg-white/20 rounded-lg p-2 backdrop-blur-sm">
            <p class="text-xs opacity-90">Bắt đầu</p>
            <p class="text-sm font-bold">{{ formatTime(currentShift.start_time) }}</p>
          </div>
          <div class="bg-white/20 rounded-lg p-2 backdrop-blur-sm">
            <p class="text-xs opacity-90">Tiền đầu ca</p>
            <p class="text-sm font-bold">{{ formatPrice(currentShift.starting_float) }}</p>
          </div>
        </div>
      </div>

      <!-- Close Shift Button -->
      <button
        v-if="canCloseShift"
        @click="goToShiftClosure"
        class="w-full py-3 bg-red-500 text-white rounded-xl font-bold text-base active:scale-95 transition-transform flex items-center justify-center gap-2"
      >
        <span>🔒</span>
        <span>Đóng ca thu ngân</span>
      </button>
    </div>

    <!-- Start Shift Modal -->
    <div v-if="showStartModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-2xl p-6 max-w-md w-full">
        <h3 class="text-xl font-bold text-gray-800 mb-4">Bắt đầu ca thu ngân</h3>
        
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Tiền đầu ca (VNĐ) <span class="text-red-500">*</span>
          </label>
          <input
            v-model.number="startingFloat"
            type="number"
            step="1000"
            min="0"
            class="w-full border-2 border-gray-300 rounded-xl px-4 py-3 text-base focus:outline-none focus:border-yellow-500"
            :class="{ 'border-red-500': startingFloatError }"
            placeholder="Nhập số tiền đầu ca"
            @input="validateStartingFloat"
          />
          <p v-if="startingFloatError" class="text-sm text-red-600 mt-1">{{ startingFloatError }}</p>
        </div>

        <div class="flex gap-3">
          <button
            @click="showStartModal = false"
            class="flex-1 py-3 bg-gray-200 text-gray-800 rounded-xl font-bold active:scale-95 transition-transform"
          >
            Hủy
          </button>
          <button
            @click="startShift"
            :disabled="!canStartShift || startingLoading"
            class="flex-1 py-3 bg-yellow-500 text-white rounded-xl font-bold active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ startingLoading ? 'Đang xử lý...' : 'Bắt đầu' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCashierShiftStore } from '../stores/cashierShift'

const router = useRouter()
const cashierShiftStore = useCashierShiftStore()

// State
const showStartModal = ref(false)
const startingFloat = ref(null)
const startingFloatError = ref(null)
const startingLoading = ref(false)

// Computed
const loading = computed(() => cashierShiftStore.loading)
const error = computed(() => cashierShiftStore.error)
const currentShift = computed(() => cashierShiftStore.currentCashierShift)
const hasOpenShift = computed(() => cashierShiftStore.hasOpenCashierShift)

const canCloseShift = computed(() => {
  return currentShift.value && currentShift.value.status === 'OPEN'
})

const canStartShift = computed(() => {
  return startingFloat.value !== null && 
         startingFloat.value >= 0 && 
         !startingFloatError.value
})

// Methods
const validateStartingFloat = () => {
  startingFloatError.value = null
  
  if (startingFloat.value === null || startingFloat.value === '') {
    startingFloatError.value = 'Vui lòng nhập tiền đầu ca'
    return false
  }
  
  if (startingFloat.value < 0) {
    startingFloatError.value = 'Số tiền không được âm'
    return false
  }
  
  return true
}

const startShift = async () => {
  if (!validateStartingFloat()) return
  
  startingLoading.value = true
  
  try {
    await cashierShiftStore.startCashierShift(startingFloat.value)
    showStartModal.value = false
    startingFloat.value = null
  } catch (error) {
    console.error('Failed to start cashier shift:', error)
  } finally {
    startingLoading.value = false
  }
}

const goToShiftClosure = () => {
  if (currentShift.value?.id) {
    router.push(`/cashier/shift-closure/${currentShift.value.id}`)
  }
}

const formatPrice = (amount) => {
  if (!amount && amount !== 0) return '0₫'
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND',
    maximumFractionDigits: 0
  }).format(amount)
}

const formatTime = (date) => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleTimeString('vi-VN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getStatusText = (status) => {
  const statusMap = {
    'OPEN': '🟢 Đang mở',
    'CLOSURE_INITIATED': '🟡 Đang đóng',
    'CLOSED': '🔴 Đã đóng'
  }
  return statusMap[status] || status
}

// Lifecycle
onMounted(async () => {
  try {
    await cashierShiftStore.fetchCurrentCashierShift()
  } catch (error) {
    // Silently ignore 404 errors (no shift found is normal)
    if (error.response?.status !== 404) {
      console.error('Failed to fetch current cashier shift:', error)
    }
  }
})
</script>

<style scoped>
.active\:scale-95:active {
  transform: scale(0.95);
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.animate-spin {
  animation: spin 1s linear infinite;
}
</style>
