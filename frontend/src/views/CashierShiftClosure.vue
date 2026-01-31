<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="sticky top-0 z-40 bg-white shadow-sm">
      <div class="px-4 py-4">
        <div class="flex items-center gap-3">
          <button
            @click="goBack"
            class="p-2 text-gray-600 active:scale-95 transition-transform"
          >
            ← Quay lại
          </button>
          <div>
            <h1 class="text-2xl font-bold text-gray-800">🔒 Đóng ca thu ngân</h1>
            <p class="text-sm text-gray-600">Quy trình đóng ca</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="px-4 py-4 pb-24">
      <!-- Loading -->
      <div v-if="loading" class="text-center py-12">
        <div class="animate-spin text-5xl mb-3">⏳</div>
        <p class="text-gray-600">Đang tải...</p>
      </div>

      <!-- Error -->
      <div v-else-if="error" class="bg-red-50 border-2 border-red-200 rounded-2xl p-4 mb-4">
        <div class="flex items-start gap-3">
          <span class="text-2xl">⚠️</span>
          <div>
            <p class="font-medium text-red-800">Lỗi</p>
            <p class="text-sm text-red-600">{{ error }}</p>
          </div>
        </div>
      </div>

      <!-- Shift Info -->
      <div v-else-if="shift" class="space-y-4">
        <!-- Shift Details Card -->
        <div class="bg-gradient-to-r from-yellow-500 to-orange-500 text-white rounded-2xl p-6 shadow-lg">
          <h2 class="text-lg font-bold mb-4">📊 Thông tin ca làm</h2>
          <div class="grid grid-cols-2 gap-3">
            <div class="bg-white/20 rounded-xl p-3 backdrop-blur-sm">
              <p class="text-xs opacity-90">Thu ngân</p>
              <p class="text-sm font-bold">{{ shift.cashier_name }}</p>
            </div>
            <div class="bg-white/20 rounded-xl p-3 backdrop-blur-sm">
              <p class="text-xs opacity-90">Trạng thái</p>
              <p class="text-sm font-bold">{{ getStatusText(shift.status) }}</p>
            </div>
            <div class="bg-white/20 rounded-xl p-3 backdrop-blur-sm">
              <p class="text-xs opacity-90">Tiền đầu ca</p>
              <p class="text-sm font-bold">{{ formatPrice(shift.starting_float) }}</p>
            </div>
            <div class="bg-white/20 rounded-xl p-3 backdrop-blur-sm">
              <p class="text-xs opacity-90">Tiền hệ thống</p>
              <p class="text-sm font-bold">{{ formatPrice(shift.system_cash) }}</p>
            </div>
          </div>
        </div>

        <!-- Step 1: Initiate Closure -->
        <div v-if="shift.status === 'OPEN'" class="bg-white rounded-2xl p-6 shadow-sm">
          <h3 class="text-lg font-bold text-gray-800 mb-4">Bước 1: Bắt đầu đóng ca</h3>
          <p class="text-sm text-gray-600 mb-4">
            Xác nhận bắt đầu quy trình đóng ca. Sau khi bắt đầu, bạn cần hoàn thành tất cả các bước.
          </p>
          <button
            @click="initiateClosure"
            :disabled="processing"
            class="w-full py-4 bg-yellow-500 text-white rounded-xl font-bold active:scale-95 transition-transform disabled:opacity-50"
          >
            {{ processing ? 'Đang xử lý...' : '▶️ Bắt đầu đóng ca' }}
          </button>
        </div>

        <!-- Step 2: Record Actual Cash -->
        <div v-if="shift.status === 'CLOSURE_INITIATED' && !shift.actual_cash" class="bg-white rounded-2xl p-6 shadow-sm">
          <h3 class="text-lg font-bold text-gray-800 mb-4">Bước 2: Nhập tiền thực tế</h3>
          <p class="text-sm text-gray-600 mb-4">
            Đếm tiền mặt trong két và nhập số tiền thực tế.
          </p>
          
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Tiền mặt thực tế (VNĐ) <span class="text-red-500">*</span>
            </label>
            <input
              v-model.number="actualCash"
              type="number"
              step="1000"
              min="0"
              class="w-full border-2 border-gray-300 rounded-xl px-4 py-3 text-base focus:outline-none focus:border-yellow-500"
              placeholder="Nhập số tiền thực tế"
            />
          </div>

          <button
            @click="recordActualCash"
            :disabled="!actualCash || processing"
            class="w-full py-4 bg-green-500 text-white rounded-xl font-bold active:scale-95 transition-transform disabled:opacity-50"
          >
            {{ processing ? 'Đang xử lý...' : '✓ Xác nhận tiền mặt' }}
          </button>
        </div>

        <!-- Step 3: Document Variance (if needed) -->
        <div v-if="shift.variance && shift.variance.amount !== 0 && !shift.variance.reason" class="bg-white rounded-2xl p-6 shadow-sm">
          <h3 class="text-lg font-bold text-gray-800 mb-4">Bước 3: Giải trình chênh lệch</h3>
          
          <div class="bg-yellow-50 border-2 border-yellow-300 rounded-xl p-4 mb-4">
            <div class="flex items-center justify-between">
              <span class="text-sm font-medium text-yellow-800">Chênh lệch:</span>
              <span :class="getVarianceClass(shift.variance.amount)" class="text-lg font-bold">
                {{ formatPrice(shift.variance.amount) }}
              </span>
            </div>
          </div>

          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Lý do <span class="text-red-500">*</span>
              </label>
              <select
                v-model="varianceReason"
                class="w-full border-2 border-gray-300 rounded-xl px-4 py-3 text-base focus:outline-none focus:border-yellow-500"
              >
                <option value="">-- Chọn lý do --</option>
                <option value="COUNTING_ERROR">Lỗi đếm tiền</option>
                <option value="UNRECORDED_SALE">Bán hàng chưa ghi nhận</option>
                <option value="THEFT">Mất cắp</option>
                <option value="CHANGE_ERROR">Lỗi trả tiền thừa</option>
                <option value="SYSTEM_ERROR">Lỗi hệ thống</option>
                <option value="OTHER">Khác</option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Ghi chú chi tiết <span class="text-red-500">*</span>
              </label>
              <textarea
                v-model="varianceNotes"
                rows="4"
                class="w-full border-2 border-gray-300 rounded-xl px-4 py-3 text-base focus:outline-none focus:border-yellow-500"
                placeholder="Mô tả chi tiết về chênh lệch (tối thiểu 10 ký tự)"
              ></textarea>
              <p class="text-xs text-gray-500 mt-1">{{ varianceNotes.length }} / 10 ký tự</p>
            </div>

            <button
              @click="documentVariance"
              :disabled="!varianceReason || varianceNotes.length < 10 || processing"
              class="w-full py-4 bg-orange-500 text-white rounded-xl font-bold active:scale-95 transition-transform disabled:opacity-50"
            >
              {{ processing ? 'Đang xử lý...' : '📝 Ghi nhận giải trình' }}
            </button>
          </div>
        </div>

        <!-- Step 4: Confirm Responsibility -->
        <div v-if="canConfirmResponsibility" class="bg-white rounded-2xl p-6 shadow-sm">
          <h3 class="text-lg font-bold text-gray-800 mb-4">Bước 4: Xác nhận trách nhiệm</h3>
          
          <div class="bg-blue-50 border-2 border-blue-200 rounded-xl p-4 mb-4">
            <p class="text-sm text-blue-800">
              Tôi xác nhận rằng tôi đã kiểm tra kỹ lưỡng và chịu trách nhiệm về số liệu tài chính trong ca làm việc này.
            </p>
          </div>

          <button
            @click="confirmResponsibility"
            :disabled="processing"
            class="w-full py-4 bg-blue-500 text-white rounded-xl font-bold active:scale-95 transition-transform disabled:opacity-50"
          >
            {{ processing ? 'Đang xử lý...' : '✓ Tôi xác nhận' }}
          </button>
        </div>

        <!-- Step 5: Close Shift -->
        <div v-if="canCloseShift" class="bg-white rounded-2xl p-6 shadow-sm">
          <h3 class="text-lg font-bold text-gray-800 mb-4">Bước 5: Hoàn tất đóng ca</h3>
          
          <!-- Waiter Shifts Warning -->
          <div v-if="waiterShiftsStatus && !waiterShiftsStatus.can_close" class="bg-red-50 border-2 border-red-300 rounded-xl p-4 mb-4">
            <div class="flex items-start gap-3 mb-3">
              <span class="text-2xl">⚠️</span>
              <div>
                <p class="font-bold text-red-800 mb-2">Không thể đóng ca!</p>
                <p class="text-sm text-red-700 mb-2">
                  Còn {{ waiterShiftsStatus.open_count }} ca waiter đang mở:
                </p>
                <ul class="text-sm text-red-700 list-disc list-inside">
                  <li v-for="openShift in waiterShiftsStatus.open_shifts" :key="openShift.id">
                    {{ openShift.user_name }} ({{ openShift.role_type }})
                  </li>
                </ul>
              </div>
            </div>
            <p class="text-xs text-red-600">
              Vui lòng đóng tất cả ca waiter trước khi đóng ca thu ngân.
            </p>
          </div>
          
          <div class="bg-green-50 border-2 border-green-200 rounded-xl p-4 mb-4">
            <p class="text-sm text-green-800 mb-2">
              ✓ Tất cả các bước đã hoàn thành. Bạn có thể đóng ca ngay bây giờ.
            </p>
            <p class="text-xs text-green-700">
              Lưu ý: Sau khi đóng ca, bạn không thể thay đổi thông tin.
            </p>
          </div>

          <button
            @click="closeShift"
            :disabled="processing"
            class="w-full py-4 bg-red-500 text-white rounded-xl font-bold active:scale-95 transition-transform disabled:opacity-50"
          >
            {{ processing ? 'Đang xử lý...' : '🔒 Đóng ca' }}
          </button>
        </div>

        <!-- Completed -->
        <div v-if="shift.status === 'CLOSED'" class="bg-white rounded-2xl p-6 shadow-sm text-center">
          <div class="text-6xl mb-4">✅</div>
          <h3 class="text-xl font-bold text-gray-800 mb-2">Ca làm đã đóng</h3>
          <p class="text-sm text-gray-600 mb-4">
            Ca làm việc đã được đóng thành công vào {{ formatDateTime(shift.end_time) }}
          </p>
          <button
            @click="goBack"
            class="px-6 py-3 bg-yellow-500 text-white rounded-xl font-bold active:scale-95 transition-transform"
          >
            Quay lại Dashboard
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import cashierShiftService from '../services/cashierShift'

const router = useRouter()
const route = useRoute()

const shift = ref(null)
const loading = ref(true)
const error = ref(null)
const processing = ref(false)
const waiterShiftsStatus = ref(null)

// Form data
const actualCash = ref(null)
const varianceReason = ref('')
const varianceNotes = ref('')

// Computed
const canConfirmResponsibility = computed(() => {
  if (!shift.value) return false
  if (shift.value.status !== 'CLOSURE_INITIATED') return false
  if (!shift.value.actual_cash) return false
  if (shift.value.confirmation) return false
  
  // If there's variance, it must be documented
  if (shift.value.variance && shift.value.variance.amount !== 0) {
    return shift.value.variance.reason !== null
  }
  
  return true
})

const canCloseShift = computed(() => {
  if (!shift.value) return false
  if (shift.value.status !== 'CLOSURE_INITIATED') return false
  return shift.value.confirmation !== null
})

// Methods
const loadShift = async () => {
  loading.value = true
  error.value = null
  
  try {
    const shiftId = route.params.id
    const response = await cashierShiftService.getCashierShift(shiftId)
    shift.value = response
  } catch (err) {
    error.value = err.response?.data?.error || 'Không thể tải thông tin ca làm'
  } finally {
    loading.value = false
  }
}

const initiateClosure = async () => {
  processing.value = true
  error.value = null
  
  try {
    await cashierShiftService.initiateClosure(shift.value.id)
    await loadShift()
  } catch (err) {
    error.value = err.response?.data?.error || 'Không thể bắt đầu đóng ca'
  } finally {
    processing.value = false
  }
}

const recordActualCash = async () => {
  processing.value = true
  error.value = null
  
  try {
    await cashierShiftService.recordActualCash(shift.value.id, actualCash.value)
    await loadShift()
  } catch (err) {
    error.value = err.response?.data?.error || 'Không thể ghi nhận tiền mặt'
  } finally {
    processing.value = false
  }
}

const documentVariance = async () => {
  processing.value = true
  error.value = null
  
  try {
    await cashierShiftService.documentVariance(shift.value.id, {
      reason: varianceReason.value,
      notes: varianceNotes.value
    })
    await loadShift()
  } catch (err) {
    error.value = err.response?.data?.error || 'Không thể ghi nhận giải trình'
  } finally {
    processing.value = false
  }
}

const confirmResponsibility = async () => {
  processing.value = true
  error.value = null
  
  try {
    await cashierShiftService.confirmResponsibility(shift.value.id)
    await loadShift()
  } catch (err) {
    error.value = err.response?.data?.error || 'Không thể xác nhận trách nhiệm'
  } finally {
    processing.value = false
  }
}

const closeShift = async () => {
  // Check waiter shifts first
  try {
    const status = await cashierShiftService.checkWaiterShifts()
    waiterShiftsStatus.value = status
    
    if (!status.can_close) {
      const shiftList = status.open_shifts.map(s => `- ${s.user_name} (${s.role_type})`).join('\n')
      const message = `Không thể đóng ca thu ngân!\n\nCòn ${status.open_count} ca waiter đang mở:\n${shiftList}\n\nVui lòng đóng tất cả ca waiter trước.`
      alert(message)
      return
    }
  } catch (err) {
    error.value = 'Không thể kiểm tra trạng thái ca waiter'
    return
  }
  
  if (!confirm('Bạn có chắc muốn đóng ca? Không thể hoàn tác sau khi đóng.')) {
    return
  }
  
  processing.value = true
  error.value = null
  
  try {
    await cashierShiftService.closeShift(shift.value.id)
    await loadShift()
  } catch (err) {
    error.value = err.response?.data?.error || 'Không thể đóng ca'
  } finally {
    processing.value = false
  }
}

const goBack = () => {
  router.push('/cashier')
}

// Utility functions
const formatPrice = (amount) => {
  if (!amount && amount !== 0) return '0₫'
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND',
    maximumFractionDigits: 0
  }).format(amount)
}

const formatDateTime = (date) => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleString('vi-VN', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
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

const getVarianceClass = (amount) => {
  if (amount > 0) return 'text-green-600'
  if (amount < 0) return 'text-red-600'
  return 'text-gray-600'
}

// Lifecycle
onMounted(() => {
  loadShift()
  checkWaiterShiftsStatus()
})

// Auto-check waiter shifts status
const checkWaiterShiftsStatus = async () => {
  try {
    const status = await cashierShiftService.checkWaiterShifts()
    waiterShiftsStatus.value = status
  } catch (err) {
    console.error('Failed to check waiter shifts:', err)
  }
}
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
