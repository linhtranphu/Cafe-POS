<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Pull to Refresh -->
    <PullToRefresh 
      :pull-distance="pullDistance" 
      :is-refreshing="isRefreshing"
      :threshold="80" />

    <!-- Header -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-4" style="padding-top: max(1rem, env(safe-area-inset-top))">
        <div class="flex items-center gap-3">
          <button
            @click="handleBack"
            class="p-2 text-gray-600 active:scale-95 transition-transform"
          >
            ← Quay lại
          </button>
          <div>
            <h1 class="text-2xl font-bold text-gray-800">🔒 Đóng ca thu ngân</h1>
            <p class="text-sm text-gray-600">Quy trình đóng ca (Frontend-driven)</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24">
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

        <!-- Managed Funds Summary Card -->
        <div v-if="managedFunds" class="bg-white rounded-2xl p-6 shadow-sm">
          <h2 class="text-lg font-bold text-gray-800 mb-4">💰 Tiền đang quản lý</h2>
          
          <!-- Starting Float -->
          <div class="mb-4">
            <div class="flex items-center justify-between py-2 border-b border-gray-200">
              <span class="text-sm text-gray-600">Tiền đầu ca</span>
              <span class="text-base font-semibold text-gray-800">{{ formatPrice(managedFunds.starting_float) }}</span>
            </div>
          </div>
          
          <!-- Received from Handovers -->
          <div class="mb-4">
            <p class="text-xs font-semibold text-gray-500 uppercase mb-2">Nhận từ waiter</p>
            <div class="space-y-2">
              <div class="flex items-center justify-between py-2 bg-green-50 rounded-lg px-3">
                <div class="flex items-center gap-2">
                  <span class="text-lg">💵</span>
                  <span class="text-sm text-green-700 font-medium">Tiền mặt</span>
                </div>
                <span class="text-base font-bold text-green-700">{{ formatPrice(managedFunds.received_cash) }}</span>
              </div>
              <div class="flex items-center justify-between py-2 bg-blue-50 rounded-lg px-3">
                <div class="flex items-center gap-2">
                  <span class="text-lg">💳</span>
                  <span class="text-sm text-blue-700 font-medium">Tiền CK</span>
                </div>
                <span class="text-base font-bold text-blue-700">{{ formatPrice(managedFunds.received_transfer) }}</span>
              </div>
            </div>
          </div>
          
          <!-- Expected Cash -->
          <div class="bg-gradient-to-r from-orange-50 to-yellow-50 rounded-xl p-4 border-2 border-orange-300 mb-4">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-xs text-gray-600 mb-1">Tổng tiền mặt lý thuyết</p>
                <p class="text-xs text-gray-500">(Đầu ca + Nhận từ waiter)</p>
              </div>
              <span class="text-xl font-bold text-orange-600">{{ formatPrice(managedFunds.expected_cash) }}</span>
            </div>
          </div>
          
          <!-- Transfer Note -->
          <div class="bg-blue-50 rounded-xl p-3 border border-blue-200">
            <div class="flex items-start gap-2">
              <span class="text-base">💳</span>
              <div class="flex-1">
                <p class="text-xs font-semibold text-blue-800 mb-1">Tiền chuyển khoản</p>
                <p class="text-xs text-blue-700">Sẽ được ghi nhận (không cần bàn giao vật lý)</p>
              </div>
            </div>
          </div>
          
          <!-- Handover Count -->
          <div v-if="managedFunds.handover_count > 0" class="mt-3 text-center">
            <p class="text-xs text-gray-500">
              Đã nhận {{ managedFunds.handover_count }} lần bàn giao từ waiter
            </p>
          </div>
        </div>

        <!-- Check Waiter Shifts First -->
        <div v-if="!waiterShiftsChecked" class="bg-white rounded-2xl p-6 shadow-sm">
          <h3 class="text-lg font-bold text-gray-800 mb-4">Bước 1: Kiểm tra ca waiter</h3>
          <p class="text-sm text-gray-600 mb-4">
            Trước khi đóng ca, cần kiểm tra tất cả ca waiter đã đóng chưa.
          </p>
          <button
            @click="checkWaiterShifts"
            :disabled="processing"
            class="w-full py-4 bg-blue-500 text-white rounded-xl font-bold active:scale-95 transition-transform disabled:opacity-50"
          >
            {{ processing ? 'Đang kiểm tra...' : '🔍 Kiểm tra ca waiter' }}
          </button>
        </div>

        <!-- Waiter Shifts Warning -->
        <div v-if="waiterShiftsStatus && !waiterShiftsStatus.can_close" class="bg-red-50 border-2 border-red-300 rounded-xl p-4">
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

        <!-- Closure Form (only if waiter shifts are closed) -->
        <div v-if="waiterShiftsStatus && waiterShiftsStatus.can_close" class="space-y-4">
          <!-- Step 2: Record Actual Cash -->
          <div class="bg-white rounded-2xl p-6 shadow-sm">
            <h3 class="text-lg font-bold text-gray-800 mb-4">Bước 2: Nhập tiền thực tế</h3>
            <p class="text-sm text-gray-600 mb-4">
              Đếm tiền mặt trong két và nhập số tiền thực tế.
            </p>
            
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Tiền mặt thực tế (VNĐ) <span class="text-red-500">*</span>
              </label>
              <input
                v-model.number="closureData.actualCash"
                @input="calculateVariance"
                type="number"
                step="1000"
                min="0"
                class="w-full border-2 border-gray-300 rounded-xl px-4 py-3 text-base focus:outline-none focus:border-yellow-500"
                placeholder="Nhập số tiền thực tế"
              />
            </div>

            <!-- Show Variance if calculated -->
            <div v-if="closureData.variance !== null" class="bg-yellow-50 border-2 border-yellow-300 rounded-xl p-4 mb-4">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-yellow-800">Chênh lệch:</span>
                <span :class="getVarianceClass(closureData.variance)" class="text-lg font-bold">
                  {{ formatPrice(closureData.variance) }}
                </span>
              </div>
            </div>
          </div>

          <!-- Step 3: Document Variance (if needed) -->
          <div v-if="needsVarianceDocumentation" class="bg-white rounded-2xl p-6 shadow-sm">
            <h3 class="text-lg font-bold text-gray-800 mb-4">Bước 3: Giải trình chênh lệch</h3>
            
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Lý do <span class="text-red-500">*</span>
                </label>
                <select
                  v-model="closureData.varianceReason"
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
                  v-model="closureData.varianceNotes"
                  rows="4"
                  class="w-full border-2 border-gray-300 rounded-xl px-4 py-3 text-base focus:outline-none focus:border-yellow-500"
                  placeholder="Mô tả chi tiết về chênh lệch (tối thiểu 10 ký tự)"
                ></textarea>
                <p class="text-xs text-gray-500 mt-1">{{ closureData.varianceNotes.length }} / 10 ký tự</p>
              </div>
            </div>
          </div>

          <!-- Step 4: Confirmation Summary -->
          <div v-if="canComplete" class="bg-gradient-to-r from-orange-50 to-yellow-50 rounded-2xl p-6 shadow-sm border-2 border-orange-300">
            <h3 class="text-lg font-bold text-gray-800 mb-4">📋 Xác nhận bàn giao về quỹ</h3>
            
            <!-- Summary Grid -->
            <div class="space-y-3 mb-4">
              <!-- Cash Handover -->
              <div class="bg-white rounded-xl p-4">
                <div class="flex items-center justify-between mb-2">
                  <div class="flex items-center gap-2">
                    <span class="text-xl">💵</span>
                    <span class="text-sm font-semibold text-gray-700">Tiền mặt bàn giao</span>
                  </div>
                  <span class="text-lg font-bold text-green-600">{{ formatPrice(closureData.actualCash) }}</span>
                </div>
                <p class="text-xs text-gray-500 pl-7">Số tiền thực tế đếm được</p>
              </div>
              
              <!-- Transfer Recorded -->
              <div v-if="managedFunds" class="bg-white rounded-xl p-4">
                <div class="flex items-center justify-between mb-2">
                  <div class="flex items-center gap-2">
                    <span class="text-xl">💳</span>
                    <span class="text-sm font-semibold text-gray-700">Tiền CK ghi nhận</span>
                  </div>
                  <span class="text-lg font-bold text-blue-600">{{ formatPrice(managedFunds.received_transfer) }}</span>
                </div>
                <p class="text-xs text-gray-500 pl-7">Không cần bàn giao vật lý</p>
              </div>
              
              <!-- Total -->
              <div class="bg-gradient-to-r from-orange-100 to-yellow-100 rounded-xl p-4 border-2 border-orange-400">
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-2">
                    <span class="text-xl">📊</span>
                    <span class="text-sm font-bold text-gray-800">Tổng cộng</span>
                  </div>
                  <span class="text-xl font-bold text-orange-600">
                    {{ formatPrice((closureData.actualCash || 0) + (managedFunds?.received_transfer || 0)) }}
                  </span>
                </div>
              </div>
              
              <!-- Variance Display -->
              <div v-if="closureData.variance !== 0" class="bg-yellow-50 rounded-xl p-4 border-2 border-yellow-400">
                <div class="flex items-center justify-between mb-2">
                  <div class="flex items-center gap-2">
                    <span class="text-xl">⚠️</span>
                    <span class="text-sm font-semibold text-yellow-800">Chênh lệch</span>
                  </div>
                  <span :class="getVarianceClass(closureData.variance)" class="text-lg font-bold">
                    {{ formatPrice(closureData.variance) }}
                  </span>
                </div>
                <div v-if="closureData.varianceReason" class="pl-7">
                  <p class="text-xs text-yellow-700 font-medium">Lý do: {{ closureData.varianceReason }}</p>
                  <p class="text-xs text-yellow-600 mt-1">{{ closureData.varianceNotes }}</p>
                </div>
              </div>
            </div>
            
            <!-- Confirmation Message -->
            <div class="bg-orange-100 rounded-xl p-4 border border-orange-300">
              <div class="flex items-start gap-2">
                <span class="text-base">✅</span>
                <div class="flex-1">
                  <p class="text-xs font-semibold text-orange-900 mb-1">Xác nhận bàn giao về quỹ</p>
                  <p class="text-xs text-orange-800">Sau khi xác nhận, ca sẽ được đóng và tiền sẽ được ghi nhận đã bàn giao về quỹ</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Complete Button -->
          <div class="bg-white rounded-2xl p-6 shadow-sm">
            <button
              @click="completeClosure"
              :disabled="!canComplete || processing"
              class="w-full py-4 bg-green-500 text-white rounded-xl font-bold active:scale-95 transition-transform disabled:opacity-50"
            >
              {{ processing ? 'Đang xử lý...' : '✅ Hoàn tất đóng ca' }}
            </button>
            <p class="text-xs text-gray-500 mt-2 text-center">
              💡 Tất cả dữ liệu sẽ được lưu trong một giao dịch duy nhất
            </p>
          </div>
        </div>

        <!-- Completed -->
        <div v-if="shift.status === CASHIER_SHIFT_STATUS.CLOSED" class="bg-white rounded-2xl p-6 shadow-sm text-center">
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

    <!-- Bottom Navigation -->
    <BottomNav />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import cashierShiftService from '../services/cashierShift'
import PullToRefresh from '../components/PullToRefresh.vue'
import BottomNav from '../components/BottomNav.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'
import { CASHIER_SHIFT_STATUS } from '../constants/shift'

const router = useRouter()
const route = useRoute()

const shift = ref(null)
const managedFunds = ref(null)
const loading = ref(true)
const error = ref(null)
const processing = ref(false)
const waiterShiftsStatus = ref(null)
const waiterShiftsChecked = ref(false)

// Closure data (kept in frontend, not saved until complete)
const closureData = ref({
  actualCash: null,
  variance: null,
  varianceReason: '',
  varianceNotes: ''
})

// Computed
const needsVarianceDocumentation = computed(() => {
  // Use tolerance of 0.01 to handle floating-point precision
  const tolerance = 0.01
  if (closureData.value.variance === null) return false
  const absVariance = Math.abs(closureData.value.variance)
  return absVariance >= tolerance
})

const canComplete = computed(() => {
  if (!closureData.value.actualCash) return false
  if (needsVarianceDocumentation.value) {
    return closureData.value.varianceReason && closureData.value.varianceNotes.length >= 10
  }
  return true
})

// Methods
const loadShift = async () => {
  loading.value = true
  error.value = null
  
  try {
    const shiftId = route.params.id
    const response = await cashierShiftService.getCashierShift(shiftId)
    shift.value = response
    
    // Fetch managed funds
    await loadManagedFunds(shiftId)
  } catch (err) {
    error.value = err.response?.data?.error || 'Không thể tải thông tin ca làm'
  } finally {
    loading.value = false
  }
}

const loadManagedFunds = async (shiftId) => {
  try {
    const funds = await cashierShiftService.getManagedFunds(shiftId)
    managedFunds.value = funds
  } catch (err) {
    console.error('Failed to load managed funds:', err)
    // Don't show error to user, just log it
  }
}

const calculateVariance = () => {
  if (closureData.value.actualCash !== null && managedFunds.value) {
    // Use expected_cash from managed funds (starting_float + received_cash)
    closureData.value.variance = closureData.value.actualCash - managedFunds.value.expected_cash
  }
}

const checkWaiterShifts = async () => {
  processing.value = true
  error.value = null
  
  try {
    const status = await cashierShiftService.checkWaiterShifts()
    waiterShiftsStatus.value = status
    waiterShiftsChecked.value = true
  } catch (err) {
    error.value = err.response?.data?.error || 'Không thể kiểm tra ca waiter'
  } finally {
    processing.value = false
  }
}

const completeClosure = async () => {
  if (!confirm('Bạn có chắc muốn hoàn tất đóng ca?\n\nTất cả dữ liệu sẽ được lưu, tiền sẽ được bàn giao về quỹ, và ca sẽ được đóng.')) {
    return
  }
  
  processing.value = true
  error.value = null
  
  try {
    const payload = {
      actual_cash: closureData.value.actualCash
    }
    
    // Add variance documentation if needed
    if (needsVarianceDocumentation.value) {
      payload.variance_reason = closureData.value.varianceReason
      payload.variance_notes = closureData.value.varianceNotes
    }
    
    // Use new API that creates fund handover
    const response = await cashierShiftService.closeShiftWithFundHandover(shift.value.id, payload)
    
    // Response includes both shift and fund_handover
    console.log('Closure complete:', response)
    
    await loadShift()
    
    alert('✅ Đã đóng ca và bàn giao về quỹ thành công!')
  } catch (err) {
    error.value = err.response?.data?.error || 'Không thể đóng ca'
  } finally {
    processing.value = false
  }
}

const handleBack = () => {
  if (closureData.value.actualCash !== null) {
    if (!confirm('Bạn có chắc muốn quay lại?\n\n⚠️ Tất cả dữ liệu đã nhập sẽ bị mất (chưa lưu vào hệ thống).')) {
      return
    }
  }
  goBack()
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
    [CASHIER_SHIFT_STATUS.OPEN]: '🟢 Đang mở',
    [CASHIER_SHIFT_STATUS.CLOSURE_INITIATED]: '🟡 Đang đóng',
    [CASHIER_SHIFT_STATUS.CLOSED]: '🔴 Đã đóng'
  }
  return statusMap[status] || status
}

const getVarianceClass = (amount) => {
  if (amount > 0) return 'text-green-600'
  if (amount < 0) return 'text-red-600'
  return 'text-gray-600'
}

// Lifecycle
const refreshData = async () => {
  await loadShift()
}

const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)

onMounted(async () => {
  await refreshData()
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
