<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Pull to Refresh Indicator -->
    <PullToRefresh 
      :pull-distance="pullDistance" 
      :is-refreshing="isRefreshing"
      :threshold="80" />
    
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3" style="padding-top: max(0.75rem, env(safe-area-inset-top))">
        <h1 class="text-xl font-bold text-gray-800">⏰ Ca làm việc</h1>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24">
      <!-- Current Shift -->
      <div v-if="currentShift" class="bg-gradient-to-r from-blue-500 to-purple-500 text-white rounded-2xl p-6 mb-4 shadow-lg">
        <div class="flex justify-between items-start mb-4">
          <div>
            <h3 class="text-2xl font-bold">Ca đang mở</h3>
            <p class="text-blue-100">{{ getShiftTypeText(currentShift.type) }}</p>
            <p v-if="currentShift.role_type" class="text-sm text-blue-100 mt-1">
              {{ getRoleTypeText(currentShift.role_type) }}
            </p>
          </div>
          <span class="bg-white text-blue-600 px-4 py-2 rounded-full font-bold text-sm">ĐANG MỞ</span>
        </div>
        
        <div class="grid grid-cols-2 gap-3 mb-4">
          <div class="bg-white bg-opacity-20 rounded-xl p-3">
            <p class="text-sm text-blue-100">Bắt đầu</p>
            <p class="font-bold">{{ formatTime(currentShift.started_at) }}</p>
          </div>
          <div class="bg-white bg-opacity-20 rounded-xl p-3">
            <p class="text-sm text-blue-100">Tiền đầu ca</p>
            <p class="font-bold">{{ formatPrice(currentShift.start_cash) }}</p>
          </div>
        </div>

        <!-- Cash Status for Waiter -->
        <div v-if="isWaiter" class="grid grid-cols-3 gap-3 mb-4">
          <div class="bg-white bg-opacity-20 rounded-xl p-3">
            <p class="text-sm text-blue-100">Tiền hiện có</p>
            <p class="font-bold">{{ formatPrice(currentShift.remaining_cash || currentShift.current_cash || 0) }}</p>
          </div>
          <div class="bg-white bg-opacity-20 rounded-xl p-3">
            <p class="text-sm text-blue-100">Đã bàn giao</p>
            <p class="font-bold">{{ formatPrice(currentShift.handed_over_cash || 0) }}</p>
          </div>
          <div class="bg-white bg-opacity-20 rounded-xl p-3">
            <p class="text-sm text-blue-100">Tổng thu</p>
            <p class="font-bold">{{ formatPrice(currentShift.total_collected || 0) }}</p>
          </div>
        </div>

        <!-- Pending Handover Status -->
        <div v-if="isWaiter && pendingHandover" class="bg-yellow-500 bg-opacity-20 rounded-xl p-3 mb-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-yellow-100">🕐 Đang chờ xác nhận bàn giao</p>
              <p class="font-bold">{{ formatPrice(pendingHandover.declared_amount) }}</p>
              <p class="text-xs text-yellow-200">{{ getHandoverTypeText(pendingHandover.handover_type) }}</p>
            </div>
            <button @click="cancelHandover(pendingHandover.id)" 
              class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded-lg text-sm">
              Hủy
            </button>
          </div>
        </div>

        <!-- Action Buttons for Waiter -->
        <div v-if="isWaiter" class="space-y-2">
          <!-- Partial Handover Button -->
          <button v-if="(currentShift.remaining_cash || currentShift.current_cash || 0) > 0 && !pendingHandover" 
            @click="showPartialHandoverForm = true"
            class="w-full bg-yellow-500 hover:bg-yellow-600 text-white px-4 py-3 rounded-xl font-bold active:scale-95 transition-transform">
            💰 Bàn giao một phần
          </button>
          
          <!-- Handover and End Shift Button -->
          <button v-if="(currentShift.remaining_cash || currentShift.current_cash || 0) > 0 && !pendingHandover"
            @click="showHandoverEndShiftForm = true"
            class="w-full bg-orange-500 hover:bg-orange-600 text-white px-4 py-3 rounded-xl font-bold active:scale-95 transition-transform">
            🏁 Bàn giao và đóng ca
          </button>
          
          <!-- Regular End Shift Button (only when no remaining cash) -->
          <button v-if="(currentShift.remaining_cash || currentShift.current_cash || 0) === 0 && !pendingHandover"
            @click="showEndShiftForm = true" 
            class="w-full bg-white text-blue-600 hover:bg-blue-50 px-4 py-3 rounded-xl font-bold active:scale-95 transition-transform">
            Kết thúc ca
          </button>
          
          <!-- Disabled state when pending -->
          <div v-if="pendingHandover" class="w-full bg-gray-400 text-gray-200 px-4 py-3 rounded-xl font-bold text-center">
            Chờ cashier xác nhận...
          </div>
        </div>
        
        <!-- Action Buttons for Non-Waiter -->
        <div v-else>
          <button @click="showEndShiftForm = true" 
            class="w-full bg-white text-blue-600 hover:bg-blue-50 px-4 py-3 rounded-xl font-bold active:scale-95 transition-transform">
            Kết thúc ca
          </button>
        </div>
      </div>

      <!-- Handover History Section for Waiter -->
      <div v-if="isWaiter && handoverHistory.length > 0" class="bg-white rounded-2xl p-6 shadow-sm mb-4">
        <h3 class="text-xl font-bold mb-4">📋 Lịch sử bàn giao</h3>
        <div class="space-y-3">
          <div v-for="handover in handoverHistory" :key="handover.id" 
            class="border rounded-xl p-4">
            <div class="flex justify-between items-start mb-2">
              <div>
                <p class="font-bold">{{ formatPrice(handover.declared_amount) }}</p>
                <p class="text-sm text-gray-500">{{ formatDate(handover.handover_at) }}</p>
                <p class="text-xs text-blue-600">{{ getHandoverTypeText(handover.handover_type) }}</p>
              </div>
              <span :class="getHandoverStatusClass(handover.status)"
                class="px-3 py-1 rounded-full text-xs font-medium">
                {{ getHandoverStatusText(handover.status) }}
              </span>
            </div>
            <div v-if="handover.waiter_note" class="text-sm text-gray-600 mb-2">
              <strong>Ghi chú:</strong> {{ handover.waiter_note }}
            </div>
            <div v-if="handover.cashier_note" class="text-sm text-green-600">
              <strong>Phản hồi cashier:</strong> {{ handover.cashier_note }}
            </div>
            <div v-if="handover.discrepancy && handover.discrepancy !== 0" class="text-sm mt-2 p-2 rounded" 
              :class="handover.discrepancy > 0 ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'">
              <strong>Chênh lệch:</strong> {{ handover.discrepancy > 0 ? '+' : '' }}{{ formatPrice(handover.discrepancy) }}
            </div>
          </div>
        </div>
      </div>

      <!-- Start Shift (Hidden for Cashier) -->
      <div v-else-if="!isCashier" class="bg-white rounded-2xl p-6 mb-4 shadow-sm">
        <h3 class="text-xl font-bold mb-4">Mở ca làm việc</h3>
        <form @submit.prevent="startShift" class="space-y-4">
          <div>
            <label class="block text-sm font-medium mb-2">Chọn ca *</label>
            <select v-model="startForm.type" required 
              class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent">
              <option value="">-- Chọn ca --</option>
              <option value="MORNING">☀️ Ca sáng (7:00 - 12:00)</option>
              <option value="AFTERNOON">🌤️ Ca chiều (12:00 - 18:00)</option>
              <option value="EVENING">🌙 Ca tối (18:00 - 22:00)</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium mb-2">Tiền đầu ca (VNĐ) *</label>
            <input v-model.number="startForm.start_cash" type="number" min="0" step="1000" required 
              class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent">
          </div>
          <button type="submit" 
            class="w-full bg-blue-500 hover:bg-blue-600 text-white px-4 py-3 rounded-xl font-bold active:scale-95 transition-transform">
            Mở ca
          </button>
        </form>
      </div>

      <!-- Shift History -->
      <div class="bg-white rounded-2xl p-6 shadow-sm">
        <h3 class="text-xl font-bold mb-4">Lịch sử ca làm việc</h3>
        
        <div v-if="loading" class="text-center py-10">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
        </div>
        
        <div v-else-if="shifts.length === 0" class="text-center py-10">
          <div class="text-4xl mb-2">📭</div>
          <p class="text-gray-500">Chưa có ca làm việc nào</p>
        </div>
        
        <div v-else class="space-y-3">
          <div v-for="shift in shifts" :key="shift.id" 
            class="border rounded-xl p-4 active:scale-98 transition-transform">
            <div class="flex justify-between items-start mb-3">
              <div>
                <h4 class="font-bold text-lg">{{ getShiftTypeText(shift.type) }}</h4>
                <p class="text-sm text-gray-500">{{ formatDate(shift.started_at) }}</p>
                <p v-if="shift.role_type" class="text-xs text-blue-600 font-medium mt-1">
                  {{ getRoleTypeText(shift.role_type) }}
                </p>
              </div>
              <span :class="shift.status === SHIFT_STATUS.OPEN ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'"
                class="px-3 py-1 rounded-full text-xs font-medium">
                {{ shift.status === SHIFT_STATUS.OPEN ? 'Đang mở' : 'Đã đóng' }}
              </span>
            </div>

            <div class="grid grid-cols-2 gap-3 text-sm">
              <div class="bg-gray-50 rounded-lg p-3">
                <p class="text-gray-500 text-xs">Tiền đầu ca</p>
                <p class="font-bold">{{ formatPrice(shift.start_cash) }}</p>
              </div>
              <div v-if="shift.status === 'CLOSED'" class="bg-gray-50 rounded-lg p-3">
                <p class="text-gray-500 text-xs">Tiền cuối ca</p>
                <p class="font-bold">{{ formatPrice(shift.end_cash) }}</p>
              </div>
              <div v-if="shift.status === 'CLOSED'" class="bg-green-50 rounded-lg p-3">
                <p class="text-gray-500 text-xs">Doanh thu</p>
                <p class="font-bold text-green-600">{{ formatPrice(shift.total_revenue) }}</p>
              </div>
              <div v-if="shift.status === 'CLOSED'" class="bg-blue-50 rounded-lg p-3">
                <p class="text-gray-500 text-xs">Số order</p>
                <p class="font-bold text-blue-600">{{ shift.total_orders }}</p>
              </div>
            </div>

            <button v-if="isCashier && shift.status === 'OPEN' && shift.id !== currentShift?.id" 
              @click="showCloseShiftForm(shift)"
              class="mt-3 w-full bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-xl text-sm font-medium active:scale-95 transition-transform">
              Chốt ca
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />

    <!-- End Shift Modal -->
    <transition name="slide-up">
      <div v-if="showEndShiftForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6">
          <h3 class="text-xl font-bold mb-4">Kết thúc ca</h3>
          <form @submit.prevent="endShift" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Tiền cuối ca (VNĐ) *</label>
              <input v-model.number="endForm.end_cash" type="number" min="0" step="1000" required 
                class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-blue-500">
            </div>
            <div class="bg-blue-50 p-4 rounded-xl">
              <p class="text-sm text-gray-600">Tiền đầu ca</p>
              <p class="font-bold text-2xl text-blue-600">{{ formatPrice(currentShift?.start_cash) }}</p>
            </div>
            <div class="flex gap-2">
              <button type="button" @click="showEndShiftForm = false" 
                class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
                Hủy
              </button>
              <button type="submit" 
                class="flex-1 bg-blue-500 hover:bg-blue-600 text-white px-4 py-3 rounded-xl font-medium">
                Kết thúc
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>

    <!-- Close Shift Modal -->
    <transition name="slide-up">
      <div v-if="showCloseForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6">
          <h3 class="text-xl font-bold mb-4">Chốt ca</h3>
          <p class="text-sm text-gray-600 mb-4">Chốt ca sẽ khóa tất cả orders trong ca này</p>
          <form @submit.prevent="closeShift" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Tiền cuối ca (VNĐ) *</label>
              <input v-model.number="closeForm.end_cash" type="number" min="0" step="1000" required 
                class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-red-500">
            </div>
            <div class="flex gap-2">
              <button type="button" @click="showCloseForm = false" 
                class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
                Hủy
              </button>
              <button type="submit" 
                class="flex-1 bg-red-500 hover:bg-red-600 text-white px-4 py-3 rounded-xl font-medium">
                Chốt ca
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>

    <!-- Partial Handover Modal -->
    <transition name="slide-up">
      <div v-if="showPartialHandoverForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6">
          <h3 class="text-xl font-bold mb-4">💰 Bàn giao một phần tiền</h3>
          
          <!-- Current Cash Info -->
          <div class="bg-blue-50 p-4 rounded-xl mb-4">
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Tiền hiện có</span>
              <span class="font-bold text-2xl text-blue-600">{{ formatPrice(currentShift?.remaining_cash || currentShift?.current_cash || 0) }}</span>
            </div>
          </div>
          
          <form @submit.prevent="createPartialHandover" class="space-y-4">
            <!-- Amount Input -->
            <div>
              <label class="block text-sm font-medium mb-2">Số tiền bàn giao (VNĐ) *</label>
              <input v-model.number="partialHandoverForm.declared_amount" 
                type="number" 
                :max="currentShift?.remaining_cash || currentShift?.current_cash || 0"
                min="1000" 
                step="1000" 
                required 
                class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-yellow-500">
            </div>
            
            <!-- Note -->
            <div>
              <label class="block text-sm font-medium mb-2">Ghi chú (tùy chọn)</label>
              <textarea v-model="partialHandoverForm.waiter_note" 
                rows="3" 
                class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-yellow-500"
                placeholder="Ghi chú về việc bàn giao..."></textarea>
            </div>
            
            <!-- Action Buttons -->
            <div class="flex gap-2">
              <button type="button" @click="showPartialHandoverForm = false" 
                class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
                Hủy
              </button>
              <button type="submit" 
                class="flex-1 bg-yellow-500 hover:bg-yellow-600 text-white px-4 py-3 rounded-xl font-medium">
                Bàn giao
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>

    <!-- Handover and End Shift Modal -->
    <transition name="slide-up">
      <div v-if="showHandoverEndShiftForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6">
          <h3 class="text-xl font-bold mb-4">🏁 Bàn giao toàn bộ và đóng ca</h3>
          
          <!-- Warning Notice -->
          <div class="bg-orange-50 border-l-4 border-orange-400 p-4 mb-4">
            <div class="flex">
              <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-orange-400" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                </svg>
              </div>
              <div class="ml-3">
                <p class="text-sm text-orange-700">
                  <strong>Lưu ý:</strong> Thao tác này sẽ bàn giao toàn bộ tiền còn lại và tự động đóng ca sau khi cashier xác nhận.
                </p>
              </div>
            </div>
          </div>
          
          <!-- Cash Summary -->
          <div class="bg-orange-50 p-4 rounded-xl mb-4">
            <div class="space-y-2">
              <div class="flex justify-between items-center">
                <span class="text-sm text-gray-600">Tiền sẽ bàn giao</span>
                <span class="font-bold text-2xl text-orange-600">{{ formatPrice(currentShift?.remaining_cash || currentShift?.current_cash || 0) }}</span>
              </div>
              <div class="flex justify-between items-center text-sm">
                <span class="text-gray-500">Tiền cuối ca</span>
                <span class="font-medium">{{ formatPrice(handoverEndShiftForm.end_cash) }}</span>
              </div>
            </div>
          </div>
          
          <form @submit.prevent="createHandoverAndEndShift" class="space-y-4">
            <!-- End Cash Input -->
            <div>
              <label class="block text-sm font-medium mb-2">Tiền cuối ca (VNĐ) *</label>
              <input v-model.number="handoverEndShiftForm.end_cash" 
                type="number" 
                min="0" 
                step="1000" 
                required 
                class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-orange-500">
              <p class="text-xs text-gray-500 mt-1">Tiền còn lại sau khi bàn giao (thường là 0)</p>
            </div>
            
            <!-- Note -->
            <div>
              <label class="block text-sm font-medium mb-2">Ghi chú (tùy chọn)</label>
              <textarea v-model="handoverEndShiftForm.waiter_note" 
                rows="3" 
                class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-orange-500"
                placeholder="Ghi chú về việc bàn giao và đóng ca..."></textarea>
            </div>
            
            <!-- Action Buttons -->
            <div class="flex gap-2">
              <button type="button" @click="showHandoverEndShiftForm = false" 
                class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
                Hủy
              </button>
              <button type="submit" 
                class="flex-1 bg-orange-500 hover:bg-orange-600 text-white px-4 py-3 rounded-xl font-medium">
                Bàn giao và đóng ca
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
import { useShiftStore } from '../stores/shift'
import { useAuthStore } from '../stores/auth'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'
import { USER_ROLES } from '../constants/user'
import { SHIFT_STATUS } from '../constants/shift'

const shiftStore = useShiftStore()
const authStore = useAuthStore()

const showEndShiftForm = ref(false)
const showCloseForm = ref(false)
const selectedShift = ref(null)

// Handover refs
const showPartialHandoverForm = ref(false)
const showHandoverEndShiftForm = ref(false)
const pendingHandover = ref(null)
const handoverHistory = ref([])

const startForm = ref({
  type: '',
  start_cash: 0
})

const endForm = ref({
  end_cash: 0
})

const closeForm = ref({
  end_cash: 0
})

const partialHandoverForm = ref({
  declared_amount: 0,
  waiter_note: ''
})

const handoverEndShiftForm = ref({
  end_cash: 0,
  waiter_note: ''
})

const loading = computed(() => shiftStore.loading)
const currentShift = computed(() => shiftStore.currentShift)
const shifts = computed(() => shiftStore.shifts)
const isCashier = computed(() => authStore.user?.role === USER_ROLES.CASHIER || authStore.user?.role === USER_ROLES.MANAGER)
const isWaiter = computed(() => authStore.user?.role === USER_ROLES.WAITER)

// Refresh data function
const refreshData = async () => {
  await shiftStore.fetchCurrentShift()
  if (isWaiter.value && currentShift.value) {
    await fetchHandoverData()
  }
  if (isCashier.value) {
    await shiftStore.fetchAllShifts()
  } else {
    await shiftStore.fetchMyShifts()
  }
}

// Pull to refresh
const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)

onMounted(async () => {
  await refreshData()
})

const fetchHandoverData = async () => {
  if (!currentShift.value?.id) return
  try {
    pendingHandover.value = await shiftStore.getPendingHandover(currentShift.value.id)
    console.log('Pending handover:', pendingHandover.value)
    handoverHistory.value = await shiftStore.getHandoverHistory(currentShift.value.id)
  } catch (error) {
    console.error('Error fetching handover data:', error)
  }
}

const startShift = async () => {
  try {
    await shiftStore.startShift(startForm.value)
    startForm.value = { type: '', start_cash: 0 }
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

const endShift = async () => {
  try {
    await shiftStore.endShift(currentShift.value.id, endForm.value.end_cash)
    showEndShiftForm.value = false
    endForm.value = { end_cash: 0 }
    await shiftStore.fetchMyShifts()
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

const showCloseShiftForm = (shift) => {
  selectedShift.value = shift
  showCloseForm.value = true
}

const closeShift = async () => {
  try {
    await shiftStore.closeShift(selectedShift.value.id, closeForm.value.end_cash)
    showCloseForm.value = false
    selectedShift.value = null
    closeForm.value = { end_cash: 0 }
    await shiftStore.fetchAllShifts()
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

// Handover functions
const createPartialHandover = async () => {
  try {
    const handoverData = {
      declared_amount: partialHandoverForm.value.declared_amount,
      handover_type: 'PARTIAL',
      waiter_note: partialHandoverForm.value.waiter_note
    }
    
    await shiftStore.createCashHandover(currentShift.value.id, handoverData)
    showPartialHandoverForm.value = false
    partialHandoverForm.value = { declared_amount: 0, waiter_note: '' }
    
    // Refresh data
    await shiftStore.fetchCurrentShift()
    await fetchHandoverData()
    
    alert('Đã gửi yêu cầu bàn giao một phần tiền. Chờ thu ngân xác nhận.')
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

const createHandoverAndEndShift = async () => {
  try {
    const handoverData = {
      declared_amount: currentShift.value?.remaining_cash || currentShift.value?.current_cash || 0,
      waiter_note: handoverEndShiftForm.value.waiter_note,
      end_cash: handoverEndShiftForm.value.end_cash
    }
    
    await shiftStore.createHandoverAndEndShift(currentShift.value.id, handoverData)
    showHandoverEndShiftForm.value = false
    handoverEndShiftForm.value = { end_cash: 0, waiter_note: '' }
    
    // Refresh data
    await shiftStore.fetchCurrentShift()
    await fetchHandoverData()
    
    alert('Đã gửi yêu cầu bàn giao toàn bộ và đóng ca. Chờ thu ngân xác nhận.')
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

const cancelHandover = async (handoverId) => {
  // Validate handover ID
  if (!handoverId) {
    alert('Lỗi: Không tìm thấy ID bàn giao')
    console.error('Invalid handover ID:', handoverId)
    return
  }
  
  if (confirm('Bạn có chắc muốn hủy yêu cầu bàn giao này?')) {
    try {
      await shiftStore.cancelHandover(handoverId)
      await fetchHandoverData()
      alert('Đã hủy yêu cầu bàn giao!')
    } catch (error) {
      const errorMsg = error.response?.data?.error || error.message
      alert('Lỗi: ' + errorMsg)
      console.error('Cancel handover error:', error)
    }
  }
}

const getShiftTypeText = (type) => {
  const types = {
    MORNING: '☀️ Ca sáng',
    AFTERNOON: '🌤️ Ca chiều',
    EVENING: '🌙 Ca tối'
  }
  return types[type] || type
}

const getRoleTypeText = (roleType) => {
  const roles = {
    waiter: '👨‍💼 Phục vụ',
    barista: '🍹 Pha chế',
    cashier: '💰 Thu ngân'
  }
  return roles[roleType] || roleType
}

const getHandoverTypeText = (type) => {
  const types = {
    'PARTIAL': 'Một phần',
    'FULL': 'Toàn bộ',
    'END_SHIFT': 'Toàn bộ + Đóng ca'
  }
  return types[type] || type
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
</style>
