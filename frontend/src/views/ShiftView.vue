<template>
  <div class="min-h-screen bg-gray-100">
    <Navigation />
    <div class="p-4">
      <h2 class="text-2xl font-bold text-gray-800 mb-6">⏰ Quản lý Ca làm việc</h2>

      <!-- Current Shift -->
      <div v-if="currentShift" class="bg-gradient-to-r from-blue-500 to-purple-500 text-white rounded-xl p-6 mb-6 shadow-lg">
        <div class="flex justify-between items-start mb-4">
          <div>
            <h3 class="text-2xl font-bold">Ca đang mở</h3>
            <p class="text-blue-100">{{ getShiftTypeText(currentShift.type) }}</p>
          </div>
          <span class="bg-white text-blue-600 px-4 py-2 rounded-full font-bold">ĐANG MỞ</span>
        </div>
        
        <div class="grid grid-cols-2 gap-4 mb-4">
          <div class="bg-white bg-opacity-20 rounded-lg p-3">
            <p class="text-sm text-blue-100">Bắt đầu</p>
            <p class="font-bold">{{ formatTime(currentShift.started_at) }}</p>
          </div>
          <div class="bg-white bg-opacity-20 rounded-lg p-3">
            <p class="text-sm text-blue-100">Tiền đầu ca</p>
            <p class="font-bold">{{ formatPrice(currentShift.start_cash) }}</p>
          </div>
        </div>

        <button @click="showEndShiftForm = true" class="w-full bg-white text-blue-600 hover:bg-blue-50 px-4 py-3 rounded-lg font-bold">
          Kết thúc ca
        </button>
      </div>

      <!-- Start Shift -->
      <div v-else class="bg-white rounded-xl p-6 mb-6 shadow-sm">
        <h3 class="text-xl font-bold mb-4">Mở ca làm việc</h3>
        <form @submit.prevent="startShift" class="space-y-4">
          <div>
            <label class="block text-sm font-medium mb-2">Chọn ca *</label>
            <select v-model="startForm.type" required class="w-full p-3 border rounded-lg">
              <option value="">-- Chọn ca --</option>
              <option value="MORNING">Ca sáng (7:00 - 12:00)</option>
              <option value="AFTERNOON">Ca chiều (12:00 - 18:00)</option>
              <option value="EVENING">Ca tối (18:00 - 22:00)</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium mb-2">Tiền đầu ca (VNĐ) *</label>
            <input v-model.number="startForm.start_cash" type="number" min="0" step="1000" required class="w-full p-3 border rounded-lg">
          </div>
          <button type="submit" class="w-full bg-blue-500 hover:bg-blue-600 text-white px-4 py-3 rounded-lg font-bold">
            Mở ca
          </button>
        </form>
      </div>

      <!-- Shift History -->
      <div class="bg-white rounded-xl p-6 shadow-sm">
        <h3 class="text-xl font-bold mb-4">Lịch sử ca làm việc</h3>
        
        <div v-if="loading" class="text-center py-10">Đang tải...</div>
        <div v-else-if="shifts.length === 0" class="text-center py-10 text-gray-500">Chưa có ca làm việc nào</div>
        <div v-else class="space-y-3">
          <div v-for="shift in shifts" :key="shift.id" class="border rounded-lg p-4">
            <div class="flex justify-between items-start mb-3">
              <div>
                <h4 class="font-bold">{{ getShiftTypeText(shift.type) }}</h4>
                <p class="text-sm text-gray-500">{{ formatDate(shift.started_at) }}</p>
              </div>
              <span :class="shift.status === 'OPEN' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'"
                class="px-3 py-1 rounded-full text-xs font-medium">
                {{ shift.status === 'OPEN' ? 'Đang mở' : 'Đã đóng' }}
              </span>
            </div>

            <div class="grid grid-cols-2 gap-3 text-sm">
              <div>
                <p class="text-gray-500">Tiền đầu ca</p>
                <p class="font-medium">{{ formatPrice(shift.start_cash) }}</p>
              </div>
              <div v-if="shift.status === 'CLOSED'">
                <p class="text-gray-500">Tiền cuối ca</p>
                <p class="font-medium">{{ formatPrice(shift.end_cash) }}</p>
              </div>
              <div v-if="shift.status === 'CLOSED'">
                <p class="text-gray-500">Doanh thu</p>
                <p class="font-medium text-green-600">{{ formatPrice(shift.total_revenue) }}</p>
              </div>
              <div v-if="shift.status === 'CLOSED'">
                <p class="text-gray-500">Số order</p>
                <p class="font-medium">{{ shift.total_orders }}</p>
              </div>
            </div>

            <button v-if="isCashier && shift.status === 'OPEN' && shift.id !== currentShift?.id" 
              @click="showCloseShiftForm(shift)"
              class="mt-3 w-full bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-lg text-sm font-medium">
              Chốt ca
            </button>
          </div>
        </div>
      </div>

      <!-- End Shift Modal -->
      <div v-if="showEndShiftForm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-xl p-6 w-full max-w-md">
          <h3 class="text-xl font-bold mb-4">Kết thúc ca</h3>
          <form @submit.prevent="endShift" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Tiền cuối ca (VNĐ) *</label>
              <input v-model.number="endForm.end_cash" type="number" min="0" step="1000" required class="w-full p-3 border rounded-lg">
            </div>
            <div class="bg-blue-50 p-3 rounded-lg">
              <p class="text-sm text-gray-600">Tiền đầu ca</p>
              <p class="font-bold text-lg">{{ formatPrice(currentShift?.start_cash) }}</p>
            </div>
            <div class="flex gap-2">
              <button type="button" @click="showEndShiftForm = false" class="flex-1 bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg">
                Hủy
              </button>
              <button type="submit" class="flex-1 bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg">
                Kết thúc
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Close Shift Modal -->
      <div v-if="showCloseForm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-xl p-6 w-full max-w-md">
          <h3 class="text-xl font-bold mb-4">Chốt ca</h3>
          <p class="text-sm text-gray-600 mb-4">Chốt ca sẽ khóa tất cả orders trong ca này</p>
          <form @submit.prevent="closeShift" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Tiền cuối ca (VNĐ) *</label>
              <input v-model.number="closeForm.end_cash" type="number" min="0" step="1000" required class="w-full p-3 border rounded-lg">
            </div>
            <div class="flex gap-2">
              <button type="button" @click="showCloseForm = false" class="flex-1 bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg">
                Hủy
              </button>
              <button type="submit" class="flex-1 bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-lg">
                Chốt ca
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useShiftStore } from '../stores/shift'
import { useAuthStore } from '../stores/auth'
import Navigation from '../components/Navigation.vue'

const shiftStore = useShiftStore()
const authStore = useAuthStore()

const showEndShiftForm = ref(false)
const showCloseForm = ref(false)
const selectedShift = ref(null)

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

onMounted(async () => {
  await shiftStore.fetchCurrentShift()
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
    MORNING: 'Ca sáng',
    AFTERNOON: 'Ca chiều',
    EVENING: 'Ca tối'
  }
  return types[type] || type
}

const formatPrice = (price) => {
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(price)
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('vi-VN')
}

const formatTime = (date) => {
  return new Date(date).toLocaleTimeString('vi-VN')
}
</script>

<style scoped>
button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
