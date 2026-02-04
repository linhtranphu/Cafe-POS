<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm">
      <div class="px-4 py-3">
        <h1 class="text-xl font-bold text-gray-800">⏰ Ca làm việc</h1>
      </div>
    </div>

    <!-- Content -->
    <div class="px-4 py-4 pb-24">
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

        <!-- Cash Status Display -->
        <div class="grid grid-cols-3 gap-2 mb-4">
          <div class="bg-white bg-opacity-20 rounded-xl p-3">
            <p class="text-xs text-blue-100">Tiền hiện có</p>
            <p class="font-bold text-sm">{{ formatPrice(currentShift.current_cash || currentShift.start_cash) }}</p>
          </div>
          <div class="bg-white bg-opacity-20 rounded-xl p-3">
            <p class="text-xs text-blue-100">Đã bàn giao</p>
            <p class="font-bold text-sm">{{ formatPrice(currentShift.handed_over_cash || 0) }}</p>
          </div>
          <div class="bg-white bg-opacity-20 rounded-xl p-3">
            <p class="text-xs text-blue-100">Còn lại</p>
            <p class="font-bold text-sm">{{ formatPrice(availableCash) }}</p>
          </div>
        </div>

        <!-- Pending Handover Status -->
        <div v-if="hasPendingHandover" class="bg-yellow-400 bg-opacity-20 border border-yellow-300 rounded-xl p-3 mb-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium">🕐 Đang chờ xác nhận bàn giao</p>
              <p class="text-xs text-blue-100">{{ formatPrice(pendingHandover.requested_amount) }}</p>
            </div>
            <button @click="cancelPendingHandover" 
              class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded-lg text-xs font-medium">
              Hủy
            </button>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="space-y-2">
          <div v-if="!hasPendingHandover" class="grid grid-cols-2 gap-2">
            <button @click="showPartialHandoverForm = true" 
              :disabled="availableCash <= 0"
              class="bg-white bg-opacity-20 hover:bg-opacity-30 text-white px-4 py-3 rounded-xl font-bold text-sm active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed">
              💰 Bàn giao một phần
            </button>
            <button @click="showHandoverEndShiftForm = true" 
              :disabled="availableCash <= 0"
              class="bg-white bg-opacity-20 hover:bg-opacity-30 text-white px-4 py-3 rounded-xl font-bold text-sm active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed">
              🏁 Bàn giao và đóng ca
            </button>
          </div>
          
          <button v-if="canEndShift" @click="showEndShiftForm = true" 
            class="w-full bg-white text-blue-600 hover:bg-blue-50 px-4 py-3 rounded-xl font-bold active:scale-95 transition-transform">
            Kết thúc ca
          </button>
          
          <button v-else-if="!hasPendingHandover && availableCash > 0" @click="showEndShiftForm = true" 
            class="w-full bg-white bg-opacity-20 text-white px-4 py-3 rounded-xl font-bold opacity-50 cursor-not-allowed">
            Kết thúc ca (cần bàn giao hết tiền)
          </button>
        </div>
      </div>

      <!-- Start Shift -->
      <div v-else class="bg-white rounded-2xl p-6 mb-4 shadow-sm">
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

      <!-- Handover History Section -->
      <div v-if="currentShift" class="bg-white rounded-2xl p-6 mb-4 shadow-sm">
        <h3 class="text-xl font-bold mb-4">📋 Lịch sử bàn giao</h3>
        
        <div v-if="handoverLoading" class="text-center py-6">
          <div class="inline-block animate-spin rounded-full h-6 w-6 border-b-2 border-blue-500"></div>
        </div>
        
        <div v-else-if="handoverHistory.length === 0" class="text-center py-6">
          <div class="text-3xl mb-2">📭</div>
          <p class="text-gray-500 text-sm">Chưa có bàn giao nào</p>
        </div>
        
        <div v-else class="space-y-3">
          <div v-for="handover in handoverHistory" :key="handover.id" 
            class="border rounded-xl p-4">
            <div class="flex justify-between items-start mb-2">
              <div>
                <h4 class="font-bold">{{ getHandoverTypeText(handover.type) }}</h4>
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
              <span :class="shift.status === 'OPEN' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'"
                class="px-3 py-1 rounded-full text-xs font-medium">
                {{ shift.status === 'OPEN' ? 'Đang mở' : 'Đã đóng' }}
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

    <!-- Partial Handover Modal -->
    <transition name="slide-up">
      <div v-if="showPartialHandoverForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6">
          <h3 class="text-xl font-bold mb-4">💰 Bàn giao một phần</h3>
          <form @submit.prevent="createPartialHandover" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Số tiền bàn giao (VNĐ) *</label>
              <input v-model.number="partialHandoverForm.requested_amount" 
                type="number" min="1000" :max="availableCash" step="1000" required 
                class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-blue-500">
              <p class="text-xs text-gray-500 mt-1">Tối đa: {{ formatPrice(availableCash) }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Ghi chú</label>
              <textarea v-model="partialHandoverForm.waiter_notes" 
                class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500" 
                rows="3" placeholder="Ghi chú về việc bàn giao..."></textarea>
            </div>
            <div class="flex gap-2">
              <button type="button" @click="showPartialHandoverForm = false" 
                class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
                Hủy
              </button>
              <button type="submit" :disabled="handoverLoading"
                class="flex-1 bg-blue-500 hover:bg-blue-600 text-white px-4 py-3 rounded-xl font-medium disabled:opacity-50">
                {{ handoverLoading ? 'Đang xử lý...' : 'Tạo yêu cầu' }}
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
          <h3 class="text-xl font-bold mb-4">🏁 Bàn giao và đóng ca</h3>
          <div class="bg-red-50 border border-red-200 rounded-xl p-4 mb-4">
            <p class="text-sm text-red-700">
              ⚠️ Thao tác này sẽ bàn giao toàn bộ tiền còn lại và đóng ca. 
              Không thể hoàn tác sau khi thực hiện.
            </p>
          </div>
          <form @submit.prevent="createHandoverAndEndShift" class="space-y-4">
            <div class="bg-blue-50 p-4 rounded-xl">
              <p class="text-sm text-gray-600">Số tiền sẽ bàn giao</p>
              <p class="font-bold text-2xl text-blue-600">{{ formatPrice(availableCash) }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Ghi chú</label>
              <textarea v-model="handoverEndShiftForm.waiter_notes" 
                class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500" 
                rows="3" placeholder="Ghi chú về việc bàn giao cuối ca..."></textarea>
            </div>
            <div class="flex gap-2">
              <button type="button" @click="showHandoverEndShiftForm = false" 
                class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
                Hủy
              </button>
              <button type="submit" :disabled="handoverLoading"
                class="flex-1 bg-red-500 hover:bg-red-600 text-white px-4 py-3 rounded-xl font-medium disabled:opacity-50">
                {{ handoverLoading ? 'Đang xử lý...' : 'Bàn giao và đóng ca' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>

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
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useShiftStore } from '../stores/shift'
import { useAuthStore } from '../stores/auth'
import BottomNav from '../components/BottomNav.vue'

const shiftStore = useShiftStore()
const authStore = useAuthStore()

const showEndShiftForm = ref(false)
const showCloseForm = ref(false)
const selectedShift = ref(null)

// Handover form states
const showPartialHandoverForm = ref(false)
const showHandoverEndShiftForm = ref(false)

const partialHandoverForm = ref({
  requested_amount: 0,
  waiter_notes: ''
})

const handoverEndShiftForm = ref({
  waiter_notes: ''
})

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

const loading = computed(() => shiftStore.loading)
const currentShift = computed(() => shiftStore.currentShift)
const shifts = computed(() => shiftStore.shifts)
const isCashier = computed(() => authStore.user?.role === 'cashier' || authStore.user?.role === 'manager')

// Handover computed properties
const handoverLoading = computed(() => shiftStore.handoverLoading)
const pendingHandover = computed(() => shiftStore.pendingHandover)
const handoverHistory = computed(() => shiftStore.handoverHistory)
const hasPendingHandover = computed(() => shiftStore.hasPendingHandover)
const availableCash = computed(() => shiftStore.availableCash)
const canEndShift = computed(() => shiftStore.canEndShift)

onMounted(async () => {
  await shiftStore.fetchCurrentShift()
  if (currentShift.value) {
    // Fetch handover data for current shift
    await Promise.all([
      shiftStore.fetchPendingHandover(currentShift.value.id),
      shiftStore.fetchHandoverHistory(currentShift.value.id)
    ])
  }
  if (isCashier.value) {
    await shiftStore.fetchAllShifts()
  } else {
    await shiftStore.fetchMyShifts()
  }
})

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

const formatDateTime = (date) => {
  return new Date(date).toLocaleString('vi-VN', {
    day: '2-digit',
    month: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Handover methods
const createPartialHandover = async () => {
  try {
    const handoverData = {
      type: 'PARTIAL',
      requested_amount: partialHandoverForm.value.requested_amount,
      waiter_notes: partialHandoverForm.value.waiter_notes
    }
    
    await shiftStore.createCashHandover(currentShift.value.id, handoverData)
    showPartialHandoverForm.value = false
    partialHandoverForm.value = { requested_amount: 0, waiter_notes: '' }
    
    // Refresh handover data
    await Promise.all([
      shiftStore.fetchPendingHandover(currentShift.value.id),
      shiftStore.fetchHandoverHistory(currentShift.value.id)
    ])
  } catch (error) {
    alert('Lỗi: ' + (shiftStore.handoverError || error.message))
  }
}

const createHandoverAndEndShift = async () => {
  try {
    const handoverData = {
      type: 'END_SHIFT',
      requested_amount: availableCash.value,
      waiter_notes: handoverEndShiftForm.value.waiter_notes
    }
    
    await shiftStore.createHandoverAndEndShift(currentShift.value.id, handoverData)
    showHandoverEndShiftForm.value = false
    handoverEndShiftForm.value = { waiter_notes: '' }
    
    // Refresh handover data
    await Promise.all([
      shiftStore.fetchPendingHandover(currentShift.value.id),
      shiftStore.fetchHandoverHistory(currentShift.value.id)
    ])
  } catch (error) {
    alert('Lỗi: ' + (shiftStore.handoverError || error.message))
  }
}

const cancelPendingHandover = async () => {
  if (!confirm('Bạn có chắc muốn hủy yêu cầu bàn giao này?')) return
  
  try {
    await shiftStore.cancelHandover(pendingHandover.value.id)
    
    // Refresh handover data
    await Promise.all([
      shiftStore.fetchPendingHandover(currentShift.value.id),
      shiftStore.fetchHandoverHistory(currentShift.value.id)
    ])
  } catch (error) {
    alert('Lỗi: ' + (shiftStore.handoverError || error.message))
  }
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
