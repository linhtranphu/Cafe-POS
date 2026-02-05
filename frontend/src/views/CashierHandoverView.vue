<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Pull to Refresh Indicator -->
    <PullToRefresh 
      :pull-distance="pullDistance" 
      :is-refreshing="isRefreshing"
      :threshold="80" />
    
    <!-- Mobile Header -->
    <div class="sticky top-0 z-40 bg-white shadow-sm">
      <div class="px-4 py-3">
        <h1 class="text-xl font-bold text-gray-800">💰 Quản lý bàn giao</h1>
      </div>
    </div>

    <!-- Content -->
    <div class="px-4 py-4 pb-24">
      <!-- Pending Handovers -->
      <div class="bg-white rounded-2xl p-6 shadow-sm mb-4">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-xl font-bold">🕐 Chờ xác nhận</h3>
          <span class="bg-red-100 text-red-800 px-3 py-1 rounded-full text-sm font-medium">
            {{ pendingHandovers.length }}
          </span>
        </div>
        
        <div v-if="loading" class="text-center py-10">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
        </div>
        
        <div v-else-if="pendingHandovers.length === 0" class="text-center py-10">
          <div class="text-4xl mb-2">✅</div>
          <p class="text-gray-500">Không có yêu cầu bàn giao nào</p>
        </div>
        
        <div v-else class="space-y-4">
          <div v-for="handover in pendingHandovers" :key="handover.id" 
            class="border-2 border-yellow-200 rounded-xl p-4 bg-yellow-50">
            
            <!-- Handover Header -->
            <div class="flex justify-between items-start mb-3">
              <div>
                <h4 class="font-bold text-lg">{{ handover.waiter_name }}</h4>
                <p class="text-sm text-gray-500">{{ formatDate(handover.handover_at) }}</p>
                <span :class="getHandoverTypeClass(handover.handover_type)"
                  class="inline-block px-2 py-1 rounded-full text-xs font-medium mt-1">
                  {{ getHandoverTypeText(handover.handover_type) }}
                </span>
              </div>
              <div class="text-right">
                <p class="text-2xl font-bold text-green-600">{{ formatPrice(handover.declared_amount) }}</p>
                <p v-if="handover.handover_type === 'END_SHIFT'" class="text-sm text-gray-500">
                  Tiền cuối ca: {{ formatPrice(handover.end_cash || 0) }}
                </p>
              </div>
            </div>
            
            <!-- Shift Cash Warning (in list) -->
            <div v-if="hasShiftCashMismatch(handover)"
              class="mb-3 p-3 rounded-lg border-2"
              :class="handover.declared_amount > getShiftInfo(handover)?.remaining_cash ? 'bg-orange-50 border-orange-300' : 'bg-yellow-50 border-yellow-300'">
              <div class="flex items-start gap-2">
                <span class="text-lg">⚠️</span>
                <div class="flex-1">
                  <p class="text-xs font-medium"
                    :class="handover.declared_amount > getShiftInfo(handover)?.remaining_cash ? 'text-orange-800' : 'text-yellow-800'">
                    {{ handover.declared_amount > getShiftInfo(handover)?.remaining_cash 
                      ? 'Khai báo nhiều hơn tiền còn lại trong ca' 
                      : 'Khai báo ít hơn tiền còn lại trong ca' }}
                  </p>
                  <p class="text-xs mt-1"
                    :class="handover.declared_amount > getShiftInfo(handover)?.remaining_cash ? 'text-orange-600' : 'text-yellow-600'">
                    Tiền còn lại: {{ formatPrice(getShiftInfo(handover)?.remaining_cash || 0) }} | 
                    Chênh: {{ formatPrice(Math.abs(handover.declared_amount - (getShiftInfo(handover)?.remaining_cash || 0))) }}
                  </p>
                </div>
              </div>
            </div>
            
            <!-- Waiter Note -->
            <div v-if="handover.waiter_note" class="bg-blue-50 p-3 rounded-lg mb-3">
              <p class="text-sm text-blue-800">
                <strong>Ghi chú từ waiter:</strong><br>
                {{ handover.waiter_note }}
              </p>
            </div>
            
            <!-- Action Buttons -->
            <div class="flex gap-2">
              <button @click="showConfirmModal(handover, 'CONFIRMED')"
                class="flex-1 bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded-xl font-medium">
                ✅ Xác nhận
              </button>
              <button @click="showConfirmModal(handover, 'REJECTED')"
                class="flex-1 bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-xl font-medium">
                ❌ Từ chối
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Today's Handovers -->
      <div class="bg-white rounded-2xl p-6 shadow-sm">
        <h3 class="text-xl font-bold mb-4">📋 Bàn giao hôm nay</h3>
        
        <div v-if="todayHandovers.length === 0" class="text-center py-10">
          <div class="text-4xl mb-2">📭</div>
          <p class="text-gray-500">Chưa có bàn giao nào hôm nay</p>
        </div>
        
        <div v-else class="space-y-3">
          <div v-for="handover in todayHandovers" :key="handover.id" 
            class="border rounded-xl p-4">
            <div class="flex justify-between items-start mb-2">
              <div>
                <h4 class="font-bold">{{ handover.waiter_name }}</h4>
                <p class="text-sm text-gray-500">{{ formatTime(handover.handover_at) }}</p>
                <span :class="getHandoverTypeClass(handover.handover_type)"
                  class="inline-block px-2 py-1 rounded-full text-xs font-medium mt-1">
                  {{ getHandoverTypeText(handover.handover_type) }}
                </span>
              </div>
              <div class="text-right">
                <p class="font-bold text-lg">{{ formatPrice(handover.declared_amount) }}</p>
                <span :class="getHandoverStatusClass(handover.status)"
                  class="px-2 py-1 rounded-full text-xs font-medium">
                  {{ getHandoverStatusText(handover.status) }}
                </span>
              </div>
            </div>
            
            <div v-if="handover.cashier_note" class="text-sm text-gray-600 mt-2">
              <strong>Ghi chú của bạn:</strong> {{ handover.cashier_note }}
            </div>
            
            <div v-if="handover.discrepancy && handover.discrepancy !== 0" class="text-sm mt-2 p-2 rounded" 
              :class="handover.discrepancy > 0 ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'">
              <strong>Chênh lệch:</strong> {{ handover.discrepancy > 0 ? '+' : '' }}{{ formatPrice(handover.discrepancy) }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />

    <!-- Confirm Modal -->
    <transition name="slide-up">
      <div v-if="showConfirmForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6">
          <h3 class="text-xl font-bold mb-4">
            {{ confirmAction === 'CONFIRMED' ? '✅ Xác nhận bàn giao' : '❌ Từ chối bàn giao' }}
          </h3>
          
          <!-- Handover Summary -->
          <div class="bg-gray-50 p-4 rounded-xl mb-4">
            <div class="flex justify-between items-center mb-2">
              <span class="text-sm text-gray-600">Waiter</span>
              <span class="font-medium">{{ selectedHandover?.waiter_name }}</span>
            </div>
            <div class="flex justify-between items-center mb-2">
              <span class="text-sm text-gray-600">Số tiền khai báo</span>
              <span class="font-bold text-lg">{{ formatPrice(selectedHandover?.declared_amount || 0) }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Loại</span>
              <span class="text-sm">{{ getHandoverTypeText(selectedHandover?.handover_type) }}</span>
            </div>
            
            <!-- Shift Cash Warning -->
            <div v-if="shiftCashWarning" class="mt-3 p-3 rounded-lg border-2"
              :class="shiftCashWarning.type === 'OVER_DECLARED' ? 'bg-orange-50 border-orange-300' : 'bg-yellow-50 border-yellow-300'">
              <div class="flex items-start gap-2">
                <span class="text-xl">⚠️</span>
                <div class="flex-1">
                  <p class="text-sm font-medium" 
                    :class="shiftCashWarning.type === 'OVER_DECLARED' ? 'text-orange-800' : 'text-yellow-800'">
                    {{ shiftCashWarning.message }}
                  </p>
                  <p class="text-xs mt-1" 
                    :class="shiftCashWarning.type === 'OVER_DECLARED' ? 'text-orange-600' : 'text-yellow-600'">
                    Tiền còn lại trong ca: {{ formatPrice(shiftCashWarning.remainingCash) }}
                  </p>
                </div>
              </div>
            </div>
          </div>
          
          <form @submit.prevent="confirmHandover" class="space-y-4">
            <!-- Actual Amount (only for CONFIRMED) -->
            <div v-if="confirmAction === 'CONFIRMED'">
              <label class="block text-sm font-medium mb-2">Số tiền thực nhận (VNĐ) *</label>
              <input v-model.number="confirmForm.actual_amount" 
                type="number" 
                min="0" 
                step="1000" 
                required 
                class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-blue-500">
            </div>
            
            <!-- Discrepancy Warning -->
            <div v-if="hasDiscrepancy && confirmAction === 'CONFIRMED'" 
              class="p-4 rounded-xl border-2"
              :class="discrepancy > 0 ? 'bg-green-50 border-green-300' : 'bg-red-50 border-red-300'">
              
              <div class="flex items-start gap-3 mb-3">
                <span class="text-2xl">{{ discrepancy > 0 ? '📈' : '📉' }}</span>
                <div class="flex-1">
                  <h4 class="font-bold" :class="discrepancy > 0 ? 'text-green-800' : 'text-red-800'">
                    {{ discrepancy > 0 ? '⚠️ Thừa tiền' : '⚠️ Thiếu tiền' }}
                  </h4>
                  <p class="text-sm mt-1" :class="discrepancy > 0 ? 'text-green-700' : 'text-red-700'">
                    Chênh lệch: <strong>{{ formatPrice(Math.abs(discrepancy)) }}</strong>
                  </p>
                  <p v-if="requiresManagerApproval" class="text-sm mt-2 font-medium text-orange-700 bg-orange-50 p-2 rounded">
                    🔔 Chênh lệch lớn hơn 100,000₫ - Cần manager phê duyệt
                  </p>
                </div>
              </div>
              
              <!-- Discrepancy Reason (Required) -->
              <div class="mb-3">
                <label class="block text-sm font-medium mb-2">Lý do chênh lệch *</label>
                <textarea v-model="confirmForm.discrepancy_reason" 
                  required
                  rows="2" 
                  class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500"
                  placeholder="Giải thích nguyên nhân chênh lệch..."></textarea>
              </div>
              
              <!-- Discrepancy Responsibility (Required) -->
              <div>
                <label class="block text-sm font-medium mb-2">Trách nhiệm *</label>
                <select v-model="confirmForm.discrepancy_responsibility" 
                  required
                  class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500">
                  <option value="">-- Chọn người chịu trách nhiệm --</option>
                  <option value="WAITER">Waiter</option>
                  <option value="CASHIER">Cashier (Tôi)</option>
                  <option value="CUSTOMER">Khách hàng</option>
                  <option value="SYSTEM">Hệ thống</option>
                  <option value="UNKNOWN">Chưa rõ</option>
                </select>
              </div>
            </div>
            
            <!-- Cashier Note -->
            <div>
              <label class="block text-sm font-medium mb-2">
                {{ confirmAction === 'CONFIRMED' ? 'Ghi chú xác nhận' : 'Lý do từ chối' }}
                {{ confirmAction === 'REJECTED' ? ' *' : '' }}
              </label>
              <textarea v-model="confirmForm.cashier_note" 
                :required="confirmAction === 'REJECTED'"
                rows="3" 
                class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500"
                :placeholder="confirmAction === 'CONFIRMED' ? 'Ghi chú về việc nhận tiền...' : 'Lý do từ chối bàn giao...'"></textarea>
            </div>
            
            <!-- Action Buttons -->
            <div class="flex gap-2">
              <button type="button" @click="showConfirmForm = false" 
                class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
                Hủy
              </button>
              <button type="submit" 
                :class="[
                  'flex-1 px-4 py-3 rounded-xl font-medium',
                  confirmAction === 'CONFIRMED' 
                    ? 'bg-green-500 hover:bg-green-600 text-white' 
                    : 'bg-red-500 hover:bg-red-600 text-white'
                ]">
                {{ confirmAction === 'CONFIRMED' ? 'Xác nhận' : 'Từ chối' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useCashierStore } from '../stores/cashier'
import { useAuthStore } from '../stores/auth'
import { shiftService } from '../services/shift'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'

const cashierStore = useCashierStore()
const authStore = useAuthStore()

// Pull to refresh
const refreshData = async () => {
  await Promise.all([
    cashierStore.fetchPendingHandovers(),
    cashierStore.fetchTodayHandovers()
  ])
  
  // Fetch shift info for all handovers
  const allHandovers = [...pendingHandovers.value, ...todayHandovers.value]
  const shiftIds = [...new Set(allHandovers.map(h => h.waiter_shift_id))]
  
  for (const shiftId of shiftIds) {
    try {
      const shift = await shiftService.getShift(shiftId)
      shiftsMap.value[shiftId] = shift
    } catch (error) {
      console.error(`Failed to fetch shift ${shiftId}:`, error)
    }
  }
}

const { isPulling, isRefreshing, pullDistance } = usePullToRefresh(refreshData)

const showConfirmForm = ref(false)
const selectedHandover = ref(null)
const confirmAction = ref('')
const confirmForm = ref({
  actual_amount: 0,
  cashier_note: '',
  discrepancy_reason: '',
  discrepancy_responsibility: ''
})

const loading = computed(() => cashierStore.loading)
const pendingHandovers = computed(() => cashierStore.pendingHandovers)
const todayHandovers = computed(() => cashierStore.todayHandovers)

// Discrepancy calculations
const discrepancy = computed(() => {
  if (!selectedHandover.value || !confirmForm.value.actual_amount) return 0
  return confirmForm.value.actual_amount - selectedHandover.value.declared_amount
})

const hasDiscrepancy = computed(() => discrepancy.value !== 0)

const discrepancyType = computed(() => {
  if (discrepancy.value < 0) return 'SHORTAGE' // Thiếu
  if (discrepancy.value > 0) return 'OVERAGE'  // Thừa
  return null
})

const requiresManagerApproval = computed(() => {
  return Math.abs(discrepancy.value) > 100000
})

// Map to store shift info by shift_id
const shiftsMap = ref({})

// Warning for declared amount vs shift cash
const shiftCashWarning = computed(() => {
  if (!selectedHandover.value) return null
  
  const handover = selectedHandover.value
  const shift = shiftsMap.value[handover.waiter_shift_id]
  
  // Check if declared amount matches expected cash from shift
  if (shift && shift.remaining_cash !== undefined) {
    const difference = handover.declared_amount - shift.remaining_cash
    
    if (difference !== 0) {
      return {
        type: difference > 0 ? 'OVER_DECLARED' : 'UNDER_DECLARED',
        difference: Math.abs(difference),
        remainingCash: shift.remaining_cash,
        message: difference > 0 
          ? `Waiter khai báo nhiều hơn tiền còn lại trong ca (${formatPrice(Math.abs(difference))})`
          : `Waiter khai báo ít hơn tiền còn lại trong ca (${formatPrice(Math.abs(difference))})`
      }
    }
  }
  
  return null
})

// Helper to get shift info for a handover
const getShiftInfo = (handover) => {
  return shiftsMap.value[handover.waiter_shift_id]
}

// Helper to check if handover has shift cash mismatch
const hasShiftCashMismatch = (handover) => {
  const shift = shiftsMap.value[handover.waiter_shift_id]
  return shift && shift.remaining_cash !== undefined && handover.declared_amount !== shift.remaining_cash
}

onMounted(async () => {
  await refreshData()
})

const showConfirmModal = (handover, action) => {
  selectedHandover.value = handover
  confirmAction.value = action
  confirmForm.value = {
    actual_amount: handover.declared_amount, // Default to declared amount
    cashier_note: '',
    discrepancy_reason: '',
    discrepancy_responsibility: ''
  }
  showConfirmForm.value = true
}

const confirmHandover = async () => {
  try {
    const data = {
      status: confirmAction.value,
      cashier_note: confirmForm.value.cashier_note
    }
    
    // Add actual_amount only for CONFIRMED
    if (confirmAction.value === 'CONFIRMED') {
      if (!confirmForm.value.actual_amount || confirmForm.value.actual_amount === 0) {
        alert('Vui lòng nhập số tiền thực nhận')
        return
      }
      
      data.actual_amount = confirmForm.value.actual_amount
      
      // Add discrepancy info if exists
      if (hasDiscrepancy.value) {
        if (!confirmForm.value.discrepancy_reason || !confirmForm.value.discrepancy_responsibility) {
          alert('Vui lòng nhập đầy đủ thông tin chênh lệch')
          return
        }
        data.discrepancy_reason = confirmForm.value.discrepancy_reason
        data.discrepancy_responsibility = confirmForm.value.discrepancy_responsibility
      }
    } else {
      // For REJECTED, cashier_note is required
      if (!confirmForm.value.cashier_note || confirmForm.value.cashier_note.trim() === '') {
        alert('Vui lòng nhập lý do từ chối')
        return
      }
    }
    
    await cashierStore.confirmHandover(selectedHandover.value.id, data)
    
    showConfirmForm.value = false
    selectedHandover.value = null
    confirmForm.value = { 
      actual_amount: 0, 
      cashier_note: '',
      discrepancy_reason: '',
      discrepancy_responsibility: ''
    }
    
    // Refresh data
    await cashierStore.fetchPendingHandovers()
    await cashierStore.fetchTodayHandovers()
    
    const message = confirmAction.value === 'CONFIRMED' 
      ? 'Đã xác nhận bàn giao thành công!' 
      : 'Đã từ chối bàn giao!'
    alert(message)
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

// Helper functions
const formatPrice = (price) => {
  return new Intl.NumberFormat('vi-VN', { 
    style: 'currency', 
    currency: 'VND',
    maximumFractionDigits: 0
  }).format(price)
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('vi-VN')
}

const formatTime = (date) => {
  return new Date(date).toLocaleTimeString('vi-VN', { 
    hour: '2-digit', 
    minute: '2-digit' 
  })
}

const getHandoverTypeText = (type) => {
  const types = {
    'PARTIAL': 'Một phần',
    'FULL': 'Toàn bộ',
    'END_SHIFT': 'Toàn bộ + Đóng ca'
  }
  return types[type] || type
}

const getHandoverTypeClass = (type) => {
  const classes = {
    'PARTIAL': 'bg-yellow-100 text-yellow-800',
    'FULL': 'bg-blue-100 text-blue-800',
    'END_SHIFT': 'bg-orange-100 text-orange-800'
  }
  return classes[type] || 'bg-gray-100 text-gray-800'
}

const getHandoverStatusText = (status) => {
  const statuses = {
    'PENDING': 'Chờ xác nhận',
    'CONFIRMED': 'Đã xác nhận',
    'REJECTED': 'Đã từ chối',
    'DISCREPANCY': 'Có chênh lệch'
  }
  return statuses[status] || status
}

const getHandoverStatusClass = (status) => {
  const classes = {
    'PENDING': 'bg-yellow-100 text-yellow-800',
    'CONFIRMED': 'bg-green-100 text-green-800',
    'REJECTED': 'bg-red-100 text-red-800',
    'DISCREPANCY': 'bg-orange-100 text-orange-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}
</script>

<style scoped>
.active\:scale-95:active {
  transform: scale(0.95);
}

.active\:scale-98:active {
  transform: scale(0.98);
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: transform 0.3s ease;
}

.slide-up-enter-from {
  transform: translateY(100%);
}

.slide-up-leave-to {
  transform: translateY(100%);
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
