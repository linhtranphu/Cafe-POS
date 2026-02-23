<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Pull to Refresh Indicator -->
    <PullToRefresh 
      :pull-distance="pullDistance" 
      :is-refreshing="isRefreshing"
      :threshold="80" />
    
    <!-- Mobile Header -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3" style="padding-top: max(0.75rem, env(safe-area-inset-top))">
        <h1 class="text-xl font-bold text-gray-800">💰 Quản lý bàn giao</h1>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24">
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
                <div class="flex gap-2 mt-1">
                  <span :class="getHandoverTypeClass(handover.handover_type)"
                    class="inline-block px-2 py-1 rounded-full text-xs font-medium">
                    {{ getHandoverTypeText(handover.handover_type) }}
                  </span>
                  <!-- Payment type badge -->
                  <span v-if="handover.cash_declared_amount > 0 && handover.transfer_declared_amount > 0"
                    class="inline-block px-2 py-1 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                    💰 Cả hai
                  </span>
                  <span v-else-if="handover.transfer_declared_amount > 0"
                    class="inline-block px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                    💳 Chuyển khoản
                  </span>
                  <span v-else-if="handover.cash_declared_amount > 0"
                    class="inline-block px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                    💵 Tiền mặt
                  </span>
                </div>
              </div>
              <div class="text-right space-y-1">
                <p v-if="handover.cash_declared_amount > 0" class="text-lg font-bold text-green-600">
                  💵 {{ formatPrice(handover.cash_declared_amount) }}
                </p>
                <p v-if="handover.transfer_declared_amount > 0" class="text-lg font-bold text-blue-600">
                  💳 {{ formatPrice(handover.transfer_declared_amount) }}
                </p>
                <p v-if="handover.cash_declared_amount > 0 && handover.transfer_declared_amount > 0" class="text-xs text-gray-500">
                  Tổng: {{ formatPrice(handover.cash_declared_amount + handover.transfer_declared_amount) }}
                </p>
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
                <div class="flex gap-2 mt-1">
                  <span :class="getHandoverTypeClass(handover.handover_type)"
                    class="inline-block px-2 py-1 rounded-full text-xs font-medium">
                    {{ getHandoverTypeText(handover.handover_type) }}
                  </span>
                  <!-- Payment type badge -->
                  <span v-if="handover.cash_declared_amount > 0 && handover.transfer_declared_amount > 0"
                    class="inline-block px-2 py-1 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                    💰 Cả hai
                  </span>
                  <span v-else-if="handover.transfer_declared_amount > 0"
                    class="inline-block px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                    💳 CK
                  </span>
                  <span v-else-if="handover.cash_declared_amount > 0"
                    class="inline-block px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                    💵 Mặt
                  </span>
                </div>
              </div>
              <div class="text-right space-y-1">
                <p v-if="handover.cash_declared_amount > 0" class="font-bold text-green-600">
                  💵 {{ formatPrice(handover.cash_declared_amount) }}
                </p>
                <p v-if="handover.transfer_declared_amount > 0" class="font-bold text-blue-600">
                  💳 {{ formatPrice(handover.transfer_declared_amount) }}
                </p>
                <span :class="getHandoverStatusClass(handover.status)"
                  class="inline-block px-2 py-1 rounded-full text-xs font-medium mt-1">
                  {{ getHandoverStatusText(handover.status) }}
                </span>
              </div>
            </div>
            
            <div v-if="handover.cashier_note" class="text-sm text-gray-600 mt-2">
              <strong>Ghi chú của bạn:</strong> {{ handover.cashier_note }}
            </div>
            
            <!-- Display separate discrepancies -->
            <div v-if="(handover.cash_discrepancy && handover.cash_discrepancy !== 0) || (handover.transfer_discrepancy && handover.transfer_discrepancy !== 0)" 
              class="text-sm mt-2 space-y-1">
              <div v-if="handover.cash_discrepancy && handover.cash_discrepancy !== 0" 
                class="p-2 rounded" 
                :class="handover.cash_discrepancy > 0 ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'">
                <strong>💵 Chênh lệch tiền mặt:</strong> {{ handover.cash_discrepancy > 0 ? '+' : '' }}{{ formatPrice(handover.cash_discrepancy) }}
              </div>
              <div v-if="handover.transfer_discrepancy && handover.transfer_discrepancy !== 0" 
                class="p-2 rounded" 
                :class="handover.transfer_discrepancy > 0 ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'">
                <strong>💳 Chênh lệch tiền CK:</strong> {{ handover.transfer_discrepancy > 0 ? '+' : '' }}{{ formatPrice(handover.transfer_discrepancy) }}
              </div>
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
        <div class="bg-white rounded-t-3xl w-full p-6 max-h-[90vh] overflow-y-auto">
          <h3 class="text-xl font-bold mb-4">
            {{ confirmAction === 'CONFIRMED' ? '✅ Xác nhận bàn giao' : '❌ Từ chối bàn giao' }}
          </h3>
          
          <!-- Handover Summary -->
          <div class="bg-gray-50 p-4 rounded-xl mb-4">
            <div class="flex justify-between items-center mb-2">
              <span class="text-sm text-gray-600">Waiter</span>
              <span class="font-medium">{{ selectedHandover?.waiter_name }}</span>
            </div>
            
            <div v-if="selectedHandover?.cash_declared_amount > 0" class="flex justify-between items-center mb-2">
              <span class="text-sm text-gray-600">💵 Tiền mặt khai báo</span>
              <span class="font-bold text-lg text-green-600">{{ formatPrice(selectedHandover?.cash_declared_amount || 0) }}</span>
            </div>
            <div v-if="selectedHandover?.transfer_declared_amount > 0" class="flex justify-between items-center mb-2">
              <span class="text-sm text-gray-600">💳 Tiền CK khai báo</span>
              <span class="font-bold text-lg text-blue-600">{{ formatPrice(selectedHandover?.transfer_declared_amount || 0) }}</span>
            </div>
            <div v-if="selectedHandover?.cash_declared_amount > 0 && selectedHandover?.transfer_declared_amount > 0" class="flex justify-between items-center mb-2 pt-2 border-t">
              <span class="text-sm text-gray-600">Tổng</span>
              <span class="font-bold text-lg">{{ formatPrice((selectedHandover?.cash_declared_amount || 0) + (selectedHandover?.transfer_declared_amount || 0)) }}</span>
            </div>
            
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Loại</span>
              <span class="text-sm">{{ getHandoverTypeText(selectedHandover?.handover_type) }}</span>
            </div>
          </div>
          
          <form @submit.prevent="confirmHandover" class="space-y-4">
            <!-- Actual Amounts (only for CONFIRMED) -->
            <div v-if="confirmAction === 'CONFIRMED'">
              <!-- Cash Actual Amount -->
              <div v-if="selectedHandover?.cash_declared_amount > 0" class="mb-4">
                <label class="block text-sm font-medium mb-2">💵 Số tiền mặt thực nhận (VNĐ) *</label>
                <input v-model.number="confirmForm.actual_cash_amount" 
                  type="number" 
                  min="0" 
                  step="1000" 
                  required 
                  class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-green-500">
                <p class="text-xs text-gray-500 mt-1">Đếm tiền mặt thực tế nhận được</p>
              </div>
              
              <!-- Transfer Actual Amount -->
              <div v-if="selectedHandover?.transfer_declared_amount > 0" class="mb-4">
                <label class="block text-sm font-medium mb-2">💳 Số tiền CK thực nhận (VNĐ) *</label>
                <input v-model.number="confirmForm.actual_transfer_amount" 
                  type="number" 
                  min="0" 
                  step="1000" 
                  required 
                  class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-blue-500">
                <p class="text-xs text-gray-500 mt-1">Kiểm tra số dư tài khoản ngân hàng</p>
              </div>
            </div>
            
            <!-- Discrepancy Warning -->
            <div v-if="hasDiscrepancy && confirmAction === 'CONFIRMED'" 
              class="p-4 rounded-xl border-2 bg-orange-50 border-orange-300">
              
              <div class="flex items-start gap-3 mb-3">
                <span class="text-2xl">⚠️</span>
                <div class="flex-1">
                  <h4 class="font-bold text-orange-800">Có chênh lệch!</h4>
                  
                  <!-- Cash Discrepancy -->
                  <div v-if="cashDiscrepancy !== 0" class="mt-2 p-2 rounded"
                    :class="cashDiscrepancy > 0 ? 'bg-green-50' : 'bg-red-50'">
                    <p class="text-sm font-medium"
                      :class="cashDiscrepancy > 0 ? 'text-green-800' : 'text-red-800'">
                      💵 {{ cashDiscrepancy > 0 ? 'Thừa' : 'Thiếu' }} tiền mặt: {{ formatPrice(Math.abs(cashDiscrepancy)) }}
                    </p>
                  </div>
                  
                  <!-- Transfer Discrepancy -->
                  <div v-if="transferDiscrepancy !== 0" class="mt-2 p-2 rounded"
                    :class="transferDiscrepancy > 0 ? 'bg-green-50' : 'bg-red-50'">
                    <p class="text-sm font-medium"
                      :class="transferDiscrepancy > 0 ? 'text-green-800' : 'text-red-800'">
                      💳 {{ transferDiscrepancy > 0 ? 'Thừa' : 'Thiếu' }} tiền CK: {{ formatPrice(Math.abs(transferDiscrepancy)) }}
                    </p>
                  </div>
                  
                  <p v-if="requiresManagerApproval" class="text-sm mt-2 font-medium text-orange-700 bg-orange-50 p-2 rounded">
                    🔔 Chênh lệch lớn - Cần manager phê duyệt
                  </p>
                </div>
              </div>
              
              <!-- Discrepancy Reason (Required) -->
              <div class="mb-3">
                <label class="block text-sm font-medium mb-2">Lý do chênh lệch *</label>
                <textarea v-model="confirmForm.discrepancy_reason" 
                  required
                  rows="2" 
                  class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-orange-500"
                  placeholder="Giải thích nguyên nhân chênh lệch..."></textarea>
              </div>
              
              <!-- Discrepancy Responsibility (Required) -->
              <div>
                <label class="block text-sm font-medium mb-2">Trách nhiệm *</label>
                <select v-model="confirmForm.discrepancy_responsibility" 
                  required
                  class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-orange-500">
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
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'

const cashierStore = useCashierStore()

// Pull to refresh
const refreshData = async () => {
  await Promise.all([
    cashierStore.fetchPendingHandovers(),
    cashierStore.fetchTodayHandovers()
  ])
}

const { isPulling, isRefreshing, pullDistance } = usePullToRefresh(refreshData)

const showConfirmForm = ref(false)
const selectedHandover = ref(null)
const confirmAction = ref('')
const confirmForm = ref({
  actual_cash_amount: 0,
  actual_transfer_amount: 0,
  cashier_note: '',
  discrepancy_reason: '',
  discrepancy_responsibility: ''
})

const loading = computed(() => cashierStore.loading)
const pendingHandovers = computed(() => cashierStore.pendingHandovers)
const todayHandovers = computed(() => cashierStore.todayHandovers)

// Discrepancy calculations
const cashDiscrepancy = computed(() => {
  if (!selectedHandover.value) return 0
  return (confirmForm.value.actual_cash_amount || 0) - (selectedHandover.value.cash_declared_amount || 0)
})

const transferDiscrepancy = computed(() => {
  if (!selectedHandover.value) return 0
  return (confirmForm.value.actual_transfer_amount || 0) - (selectedHandover.value.transfer_declared_amount || 0)
})

const totalDiscrepancy = computed(() => {
  return cashDiscrepancy.value + transferDiscrepancy.value
})

const hasDiscrepancy = computed(() => totalDiscrepancy.value !== 0)

const requiresManagerApproval = computed(() => {
  return Math.abs(totalDiscrepancy.value) > 100000
})

onMounted(async () => {
  await refreshData()
})

const showConfirmModal = (handover, action) => {
  selectedHandover.value = handover
  confirmAction.value = action
  confirmForm.value = {
    actual_cash_amount: handover.cash_declared_amount || 0, // Default to declared
    actual_transfer_amount: handover.transfer_declared_amount || 0, // Default to declared
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
      cashier_note: confirmForm.value.cashier_note,
      actual_cash_amount: confirmForm.value.actual_cash_amount || 0,
      actual_transfer_amount: confirmForm.value.actual_transfer_amount || 0
    }
    
    // Validate for CONFIRMED
    if (confirmAction.value === 'CONFIRMED') {
      // Check if at least one actual amount is provided
      if (selectedHandover.value.cash_declared_amount > 0 && !confirmForm.value.actual_cash_amount) {
        alert('Vui lòng nhập số tiền mặt thực nhận')
        return
      }
      if (selectedHandover.value.transfer_declared_amount > 0 && !confirmForm.value.actual_transfer_amount) {
        alert('Vui lòng nhập số tiền CK thực nhận')
        return
      }
      
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
      actual_cash_amount: 0,
      actual_transfer_amount: 0,
      cashier_note: '',
      discrepancy_reason: '',
      discrepancy_responsibility: ''
    }
    
    // Refresh data
    await refreshData()
    
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
