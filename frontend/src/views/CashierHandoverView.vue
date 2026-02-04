<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm">
      <div class="px-4 py-3">
        <h1 class="text-xl font-bold text-gray-800">💰 Bàn giao tiền</h1>
      </div>
    </div>

    <!-- Content -->
    <div class="px-4 py-4 pb-24">
      <!-- Error Alert -->
      <div v-if="handoverError" class="bg-red-50 border-2 border-red-200 rounded-2xl p-4 mb-4">
        <div class="flex items-start justify-between">
          <div class="flex items-start gap-3">
            <span class="text-2xl">⚠️</span>
            <div>
              <p class="font-medium text-red-800">Lỗi</p>
              <p class="text-sm text-red-600">{{ handoverError }}</p>
            </div>
          </div>
          <button @click="clearHandoverError" class="text-red-600 text-xl font-bold">×</button>
        </div>
      </div>

      <!-- Pending Handovers Section -->
      <div class="bg-white rounded-2xl p-6 mb-4 shadow-sm">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-xl font-bold">🕐 Chờ xác nhận</h3>
          <button @click="refreshHandovers" :disabled="handoverLoading"
            class="p-2 bg-blue-500 text-white rounded-lg active:scale-95 transition-transform disabled:opacity-50">
            <span class="text-sm" :class="{ 'animate-spin': handoverLoading }">🔄</span>
          </button>
        </div>
        
        <div v-if="handoverLoading" class="text-center py-10">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
        </div>
        
        <div v-else-if="pendingHandovers.length === 0" class="text-center py-10">
          <div class="text-4xl mb-2">✅</div>
          <p class="text-gray-500">Không có yêu cầu nào</p>
        </div>
        
        <div v-else class="space-y-4">
          <div v-for="handover in pendingHandovers" :key="handover.id" 
            class="border-2 border-orange-200 rounded-xl p-4 bg-orange-50">
            <div class="flex justify-between items-start mb-3">
              <div>
                <h4 class="font-bold text-lg">{{ handover.waiter_name }}</h4>
                <p class="text-sm text-gray-600">{{ getHandoverTypeText(handover.type) }}</p>
                <p class="text-xs text-gray-500">{{ formatDateTime(handover.requested_at) }}</p>
              </div>
              <div class="text-right">
                <p class="text-2xl font-bold text-orange-600">{{ formatPrice(handover.requested_amount) }}</p>
                <span class="bg-yellow-100 text-yellow-800 px-2 py-1 rounded-full text-xs font-medium">
                  Chờ xử lý
                </span>
              </div>
            </div>

            <div v-if="handover.waiter_notes" class="bg-white p-3 rounded-lg mb-3">
              <p class="text-sm text-gray-700">💬 {{ handover.waiter_notes }}</p>
            </div>

            <!-- Action Buttons -->
            <div class="grid grid-cols-3 gap-2">
              <button @click="quickConfirmHandover(handover)" :disabled="handoverLoading"
                class="bg-green-500 hover:bg-green-600 text-white px-4 py-3 rounded-xl font-bold text-sm active:scale-95 transition-transform disabled:opacity-50">
                ✅ Xác nhận
              </button>
              <button @click="showReconcileModal(handover)" :disabled="handoverLoading"
                class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-3 rounded-xl font-bold text-sm active:scale-95 transition-transform disabled:opacity-50">
                ⚖️ Đối soát
              </button>
              <button @click="showRejectModal(handover)" :disabled="handoverLoading"
                class="bg-red-500 hover:bg-red-600 text-white px-4 py-3 rounded-xl font-bold text-sm active:scale-95 transition-transform disabled:opacity-50">
                ❌ Từ chối
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Today's Handovers Section -->
      <div class="bg-white rounded-2xl p-6 shadow-sm">
        <h3 class="text-xl font-bold mb-4">📅 Hôm nay</h3>
        
        <div v-if="todayHandovers.length === 0" class="text-center py-10">
          <div class="text-4xl mb-2">📭</div>
          <p class="text-gray-500">Chưa có bàn giao nào</p>
        </div>
        
        <div v-else class="space-y-3">
          <div v-for="handover in todayHandovers" :key="handover.id" 
            class="border rounded-xl p-4">
            <div class="flex justify-between items-start mb-2">
              <div>
                <h4 class="font-bold">{{ handover.waiter_name }}</h4>
                <p class="text-sm text-gray-500">{{ formatDateTime(handover.requested_at) }}</p>
              </div>
              <span :class="getHandoverStatusClass(handover.status)"
                class="px-3 py-1 rounded-full text-xs font-medium">
                {{ getHandoverStatusText(handover.status) }}
              </span>
            </div>

            <div class="grid grid-cols-2 gap-3 text-sm">
              <div class="bg-gray-50 rounded-lg p-3">
                <p class="text-gray-500 text-xs">Yêu cầu</p>
                <p class="font-bold">{{ formatPrice(handover.requested_amount) }}</p>
              </div>
              <div v-if="handover.actual_amount" class="bg-gray-50 rounded-lg p-3">
                <p class="text-gray-500 text-xs">Thực tế</p>
                <p class="font-bold">{{ formatPrice(handover.actual_amount) }}</p>
              </div>
            </div>

            <div v-if="handover.discrepancy_amount" class="mt-2 p-2 bg-yellow-50 rounded-lg">
              <p class="text-xs text-yellow-700">
                Chênh lệch: {{ formatPrice(Math.abs(handover.discrepancy_amount)) }}
                ({{ handover.discrepancy_amount > 0 ? 'Thừa' : 'Thiếu' }})
              </p>
            </div>

            <div v-if="handover.cashier_notes" class="mt-2 p-2 bg-blue-50 rounded-lg">
              <p class="text-xs text-blue-700">💬 {{ handover.cashier_notes }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />

    <!-- Reconciliation Modal -->
    <transition name="slide-up">
      <div v-if="showReconcileForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6">
          <h3 class="text-xl font-bold mb-4">⚖️ Đối soát bàn giao</h3>
          <div v-if="selectedHandover" class="mb-4">
            <div class="bg-blue-50 p-4 rounded-xl mb-4">
              <p class="text-sm text-gray-600">Từ: {{ selectedHandover.waiter_name }}</p>
              <p class="text-sm text-gray-600">Yêu cầu: {{ formatPrice(selectedHandover.requested_amount) }}</p>
            </div>
          </div>
          <form @submit.prevent="reconcileHandover" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Số tiền thực tế (VNĐ) *</label>
              <input v-model.number="reconcileForm.actual_amount" 
                type="number" min="0" step="1000" required 
                class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-blue-500">
            </div>
            
            <div v-if="calculatedDiscrepancy !== 0" class="bg-yellow-50 border border-yellow-200 rounded-xl p-4">
              <p class="text-sm font-medium text-yellow-800">
                Chênh lệch: {{ formatPrice(Math.abs(calculatedDiscrepancy)) }}
                ({{ calculatedDiscrepancy > 0 ? 'Thừa tiền' : 'Thiếu tiền' }})
              </p>
              <div class="mt-3 space-y-3">
                <div>
                  <label class="block text-sm font-medium mb-2">Lý do chênh lệch *</label>
                  <select v-model="reconcileForm.discrepancy_reason" required 
                    class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-yellow-500">
                    <option value="">-- Chọn lý do --</option>
                    <option value="Thiếu tiền lẻ">Thiếu tiền lẻ</option>
                    <option value="Nhầm lẫn khi đếm">Nhầm lẫn khi đếm</option>
                    <option value="Khách trả thiếu">Khách trả thiếu</option>
                    <option value="Lỗi hệ thống">Lỗi hệ thống</option>
                    <option value="Khác">Khác</option>
                  </select>
                </div>
                <div>
                  <label class="block text-sm font-medium mb-2">Trách nhiệm</label>
                  <select v-model="reconcileForm.responsibility" required 
                    class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-yellow-500">
                    <option value="">-- Chọn trách nhiệm --</option>
                    <option value="WAITER">Phục vụ</option>
                    <option value="CASHIER">Thu ngân</option>
                    <option value="SYSTEM">Hệ thống</option>
                    <option value="UNKNOWN">Chưa rõ</option>
                  </select>
                </div>
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium mb-2">Ghi chú</label>
              <textarea v-model="reconcileForm.cashier_notes" 
                class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500" 
                rows="3" placeholder="Ghi chú về việc đối soát..."></textarea>
            </div>

            <div class="flex gap-2">
              <button type="button" @click="showReconcileForm = false" 
                class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
                Hủy
              </button>
              <button type="submit" :disabled="handoverLoading"
                class="flex-1 bg-blue-500 hover:bg-blue-600 text-white px-4 py-3 rounded-xl font-medium disabled:opacity-50">
                {{ handoverLoading ? 'Đang xử lý...' : 'Xác nhận đối soát' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>

    <!-- Reject Modal -->
    <transition name="slide-up">
      <div v-if="showRejectForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6">
          <h3 class="text-xl font-bold mb-4">❌ Từ chối bàn giao</h3>
          <div v-if="selectedHandover" class="mb-4">
            <div class="bg-red-50 p-4 rounded-xl">
              <p class="text-sm text-gray-600">Từ: {{ selectedHandover.waiter_name }}</p>
              <p class="text-lg font-bold">{{ formatPrice(selectedHandover.requested_amount) }}</p>
            </div>
          </div>
          <form @submit.prevent="rejectHandover" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Lý do từ chối *</label>
              <textarea v-model="rejectForm.reason" 
                class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-red-500" 
                rows="4" placeholder="Nhập lý do từ chối..." required></textarea>
            </div>
            <div class="flex gap-2">
              <button type="button" @click="showRejectForm = false" 
                class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
                Hủy
              </button>
              <button type="submit" :disabled="handoverLoading"
                class="flex-1 bg-red-500 hover:bg-red-600 text-white px-4 py-3 rounded-xl font-medium disabled:opacity-50">
                {{ handoverLoading ? 'Đang xử lý...' : 'Từ chối' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useCashierStore } from '../stores/cashier'
import BottomNav from '../components/BottomNav.vue'

const cashierStore = useCashierStore()

// Form states
const showReconcileForm = ref(false)
const showRejectForm = ref(false)
const selectedHandover = ref(null)

const reconcileForm = ref({
  actual_amount: 0,
  discrepancy_reason: '',
  responsibility: '',
  cashier_notes: ''
})

const rejectForm = ref({
  reason: ''
})

// Computed properties
const pendingHandovers = computed(() => cashierStore.pendingHandovers)
const todayHandovers = computed(() => cashierStore.todayHandovers)
const handoverLoading = computed(() => cashierStore.handoverLoading)
const handoverError = computed(() => cashierStore.handoverError)

const calculatedDiscrepancy = computed(() => {
  if (!selectedHandover.value || !reconcileForm.value.actual_amount) return 0
  return reconcileForm.value.actual_amount - selectedHandover.value.requested_amount
})

// Methods
const refreshHandovers = async () => {
  await Promise.all([
    cashierStore.fetchPendingHandovers(),
    cashierStore.fetchTodayHandovers()
  ])
}

const quickConfirmHandover = async (handover) => {
  if (confirm(`Xác nhận bàn giao ${formatPrice(handover.requested_amount)} từ ${handover.waiter_name}?`)) {
    try {
      await cashierStore.quickConfirm(handover.id, 'Xác nhận nhanh - số tiền chính xác')
      await refreshHandovers()
    } catch (error) {
      console.error('Quick confirm failed:', error)
    }
  }
}

const showReconcileModal = (handover) => {
  selectedHandover.value = handover
  reconcileForm.value = {
    actual_amount: handover.requested_amount,
    discrepancy_reason: '',
    responsibility: '',
    cashier_notes: ''
  }
  showReconcileForm.value = true
}

const reconcileHandover = async () => {
  try {
    await cashierStore.reconcileHandover(selectedHandover.value.id, reconcileForm.value)
    showReconcileForm.value = false
    selectedHandover.value = null
    await refreshHandovers()
  } catch (error) {
    console.error('Reconcile failed:', error)
  }
}

const showRejectModal = (handover) => {
  selectedHandover.value = handover
  rejectForm.value = { reason: '' }
  showRejectForm.value = true
}

const rejectHandover = async () => {
  try {
    await cashierStore.rejectHandover(selectedHandover.value.id, rejectForm.value.reason)
    showRejectForm.value = false
    selectedHandover.value = null
    await refreshHandovers()
  } catch (error) {
    console.error('Reject failed:', error)
  }
}

const clearHandoverError = () => {
  cashierStore.clearHandoverError()
}

// Utility functions
const formatPrice = (amount) => {
  return new Intl.NumberFormat('vi-VN', { 
    style: 'currency', 
    currency: 'VND',
    maximumFractionDigits: 0
  }).format(amount)
}

const formatDateTime = (date) => {
  return new Date(date).toLocaleString('vi-VN', {
    day: '2-digit',
    month: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getHandoverTypeText = (type) => {
  const types = {
    PARTIAL: '💰 Bàn giao một phần',
    FULL: '💰 Bàn giao toàn bộ',
    END_SHIFT: '🏁 Bàn giao cuối ca'
  }
  return types[type] || type
}

const getHandoverStatusText = (status) => {
  const statuses = {
    PENDING: 'Chờ xác nhận',
    CONFIRMED: 'Đã xác nhận',
    REJECTED: 'Đã từ chối',
    DISCREPANCY: 'Có chênh lệch'
  }
  return statuses[status] || status
}

const getHandoverStatusClass = (status) => {
  const classes = {
    PENDING: 'bg-yellow-100 text-yellow-800',
    CONFIRMED: 'bg-green-100 text-green-800',
    REJECTED: 'bg-red-100 text-red-800',
    DISCREPANCY: 'bg-orange-100 text-orange-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

// Lifecycle
onMounted(async () => {
  await refreshHandovers()
  
  // Set up auto-refresh every 15 seconds for handover view
  const refreshInterval = setInterval(async () => {
    if (!document.hidden) { // Only refresh when tab is visible
      await refreshHandovers()
    }
  }, 15000)
  
  // Cleanup interval on unmount
  onUnmounted(() => {
    clearInterval(refreshInterval)
  })
})
</script>

<style scoped>
.active\:scale-95:active {
  transform: scale(0.95);
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
</style>