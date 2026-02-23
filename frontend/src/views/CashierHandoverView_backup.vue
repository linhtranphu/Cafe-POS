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
              <div class="text-right">
                <!-- Display separate amounts if available -->
                <div v-if="handover.cash_declared_amount > 0 || handover.transfer_declared_amount > 0" class="space-y-1">
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
                <!-- Fallback to old format -->
                <p v-else class="text-2xl font-bold text-green-600">{{ formatPrice(handover.declared_amount) }}</p>
                
                <p v-if="handover.handover_type === 'END_SHIFT'" class="text-sm text-gray-500 mt-1">
                  Tiền cuối ca: {{ formatPrice(handover.end_cash || 0) }}
                </p>
              </div>
            </div>
            
            <!-- Shift Cash/Transfer Warning (in list) -->
            <div v-if="hasShiftCashMismatch(handover)"
              class="mb-3 p-3 rounded-lg border-2 bg-yellow-50 border-yellow-300">
              <div class="flex items-start gap-2">
                <span class="text-lg">⚠️</span>
                <div class="flex-1">
                  <!-- New format: separate cash and transfer -->
                  <template v-if="handover.cash_declared_amount > 0 || handover.transfer_declared_amount > 0">
                    <!-- Cash warning -->
                    <div v-if="handover.cash_declared_amount > 0" class="mb-2">
                      <p class="text-xs font-medium"
                        :class="handover.cash_declared_amount > getShiftInfo(handover)?.remaining_cash ? 'text-orange-800' : 'text-yellow-800'">
                        💵 {{ handover.cash_declared_amount > getShiftInfo(handover)?.remaining_cash 
                          ? 'Tiền mặt: Khai báo nhiều hơn tiền còn lại' 
                          : 'Tiền mặt: Khai báo ít hơn tiền còn lại' }}
                      </p>
                      <p class="text-xs mt-1"
                        :class="handover.cash_declared_amount > getShiftInfo(handover)?.remaining_cash ? 'text-orange-600' : 'text-yellow-600'">
                        Còn lại: {{ formatPrice(getShiftInfo(handover)?.remaining_cash || 0) }} | 
                        Chênh: {{ formatPrice(Math.abs(handover.cash_declared_amount - (getShiftInfo(handover)?.remaining_cash || 0))) }}
                      </p>
                    </div>
                    
                    <!-- Transfer warning -->
                    <div v-if="handover.transfer_declared_amount > 0">
                      <p class="text-xs font-medium"
                        :class="handover.transfer_declared_amount > getShiftInfo(handover)?.remaining_transfer ? 'text-orange-800' : 'text-yellow-800'">
                        💳 {{ handover.transfer_declared_amount > getShiftInfo(handover)?.remaining_transfer 
                          ? 'Tiền CK: Khai báo nhiều hơn tiền còn lại' 
                          : 'Tiền CK: Khai báo ít hơn tiền còn lại' }}
                      </p>
                      <p class="text-xs mt-1"
                        :class="handover.transfer_declared_amount > getShiftInfo(handover)?.remaining_transfer ? 'text-orange-600' : 'text-yellow-600'">
                        Còn lại: {{ formatPrice(getShiftInfo(handover)?.remaining_transfer || 0) }} | 
                        Chênh: {{ formatPrice(Math.abs(handover.transfer_declared_amount - (getShiftInfo(handover)?.remaining_transfer || 0))) }}
                      </p>
                    </div>
                  </template>
                  
                  <!-- Old format: single declared_amount -->
                  <template v-else>
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
                  </template>
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
              <div class="text-right">
                <!-- Display separate amounts if available -->
                <div v-if="handover.cash_declared_amount > 0 || handover.transfer_declared_amount > 0" class="space-y-1">
                  <p v-if="handover.cash_declared_amount > 0" class="font-bold text-green-600">
                    💵 {{ formatPrice(handover.cash_declared_amount) }}
                  </p>
                  <p v-if="handover.transfer_declared_amount > 0" class="font-bold text-blue-600">
                    💳 {{ formatPrice(handover.transfer_declared_amount) }}
                  </p>
                </div>
                <!-- Fallback to old format -->
                <p v-else class="font-bold text-lg">{{ formatPrice(handover.declared_amount) }}</p>
                
                <span :class="getHandoverStatusClass(handover.status)"
                  class="inline-block px-2 py-1 rounded-full text-xs font-medium mt-1">
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
            
            <!-- Display separate amounts if available -->
            <template v-if="selectedHandover?.cash_declared_amount > 0 || selectedHandover?.transfer_declared_amount > 0">
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
            </template>
            
            <!-- Fallback to old format -->
            <div v-else class="flex justify-between items-center mb-2">
              <span class="text-sm text-gray-600">Số tiền khai báo</span>
              <span class="font-bold text-lg">{{ formatPrice(selectedHandover?.declared_amount || 0) }}</span>
            </div>
            
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Loại</span>
              <span class="text-sm">{{ getHandoverTypeText(selectedHandover?.handover_type) }}</span>
            </div>
            
            <!-- Shift Cash/Transfer Warning -->
            <div v-if="shiftCashWarning" class="mt-3 p-3 rounded-lg border-2 bg-yellow-50 border-yellow-300">
              <div class="flex items-start gap-2">
                <span class="text-xl">⚠️</span>
                <div class="flex-1">
                  <p class="text-sm font-medium text-yellow-800">
                    {{ shiftCashWarning.message }}
                  </p>
                  <div class="text-xs mt-2 space-y-1">
                    <p v-if="shiftCashWarning.cashRemaining !== undefined" class="text-yellow-600">
                      💵 Tiền mặt còn lại: {{ formatPrice(shiftCashWarning.cashRemaining) }}
                    </p>
                    <p v-if="shiftCashWarning.transferRemaining !== undefined" class="text-yellow-600">
                      💳 Tiền CK còn lại: {{ formatPrice(shiftCashWarning.transferRemaining) }}
                    </p>
                  </div>
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

// Warning for declared amount vs shift cash/transfer
const shiftCashWarning = computed(() => {
  if (!selectedHandover.value) return null
  
  const handover = selectedHandover.value
  const shift = shiftsMap.value[handover.waiter_shift_id]
  
  if (!shift) return null
  
  const warnings = []
  
  // Determine if using new format (separate amounts) or old format
  const usingNewFormat = (handover.cash_declared_amount > 0 || handover.transfer_declared_amount > 0)
  
  // Check cash mismatch
  const cashDeclared = usingNewFormat ? (handover.cash_declared_amount || 0) : (handover.declared_amount || 0)
  const cashRemaining = shift.remaining_cash || 0
  if (cashDeclared > 0 && cashDeclared !== cashRemaining) {
    const cashDiff = cashDeclared - cashRemaining
    warnings.push(cashDiff > 0 
      ? `Tiền mặt khai báo nhiều hơn (${formatPrice(Math.abs(cashDiff))})`
      : `Tiền mặt khai báo ít hơn (${formatPrice(Math.abs(cashDiff))})`)
  }
  
  // Check transfer mismatch (only for new format)
  if (usingNewFormat) {
    const transferDeclared = handover.transfer_declared_amount || 0
    const transferRemaining = shift.remaining_transfer || 0
    if (transferDeclared > 0 && transferDeclared !== transferRemaining) {
      const transferDiff = transferDeclared - transferRemaining
      warnings.push(transferDiff > 0 
        ? `Tiền CK khai báo nhiều hơn (${formatPrice(Math.abs(transferDiff))})`
        : `Tiền CK khai báo ít hơn (${formatPrice(Math.abs(transferDiff))})`)
    }
  }
  
  if (warnings.length > 0) {
    return {
      message: warnings.join(' | '),
      cashRemaining: cashDeclared > 0 ? cashRemaining : undefined,
      transferRemaining: (usingNewFormat && handover.transfer_declared_amount > 0) ? transferRemaining : undefined
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
  if (!shift) return false
  
  // For backward compatibility, check if using new format or old format
  const usingNewFormat = (handover.cash_declared_amount > 0 || handover.transfer_declared_amount > 0)
  
  if (usingNewFormat) {
    // New format: check each type separately
    const cashDeclared = handover.cash_declared_amount || 0
    const cashRemaining = shift.remaining_cash || 0
    const hasCashMismatch = cashDeclared > 0 && cashDeclared !== cashRemaining
    
    const transferDeclared = handover.transfer_declared_amount || 0
    const transferRemaining = shift.remaining_transfer || 0
    const hasTransferMismatch = transferDeclared > 0 && transferDeclared !== transferRemaining
    
    return hasCashMismatch || hasTransferMismatch
  } else {
    // Old format: compare declared_amount with remaining_cash
    const declaredAmount = handover.declared_amount || 0
    const cashRemaining = shift.remaining_cash || 0
    return declaredAmount > 0 && declaredAmount !== cashRemaining
  }
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
