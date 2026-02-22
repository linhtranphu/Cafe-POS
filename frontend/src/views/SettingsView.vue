<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Pull to Refresh -->
    <PullToRefresh 
      :pull-distance="pullDistance" 
      :is-refreshing="isRefreshing"
      :threshold="80" />

    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3" style="padding-top: max(0.75rem, env(safe-area-inset-top))">
        <h1 class="text-xl font-bold text-gray-800">⚙️ Cài đặt</h1>
      </div>
    </div>

    <!-- Content - Scrollable -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24">
      <div v-if="loading" class="text-center py-10">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
      </div>

      <div v-else class="space-y-4">
        <!-- Shop Information Card -->
        <div class="bg-white rounded-2xl p-6 shadow-sm">
          <h3 class="text-lg font-bold mb-4">🏪 Thông tin cửa hàng</h3>
          
          <div class="space-y-4">
            <!-- Shop Name -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Tên cửa hàng *
              </label>
              <input 
                v-model="shopSettings.shop_name" 
                type="text" 
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Cafe POS">
            </div>

            <!-- Shop Address -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Địa chỉ
              </label>
              <input 
                v-model="shopSettings.shop_address" 
                type="text" 
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="123 Main Street">
            </div>

            <!-- Shop Phone -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Số điện thoại
              </label>
              <input 
                v-model="shopSettings.shop_phone" 
                type="text" 
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="0123-456-789">
            </div>

            <!-- Custom Message -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Tin nhắn tùy chỉnh (hiển thị trên hóa đơn)
              </label>
              <textarea 
                v-model="shopSettings.custom_message" 
                rows="3"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Cảm ơn quý khách! Hẹn gặp lại!"></textarea>
            </div>

            <!-- Display Options -->
            <div class="space-y-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Hiển thị trên hóa đơn
              </label>
              
              <label class="flex items-center space-x-3 cursor-pointer">
                <input 
                  v-model="shopSettings.show_address" 
                  type="checkbox" 
                  class="w-5 h-5 text-blue-500 rounded focus:ring-2 focus:ring-blue-500">
                <span class="text-sm text-gray-700">Hiển thị địa chỉ</span>
              </label>

              <label class="flex items-center space-x-3 cursor-pointer">
                <input 
                  v-model="shopSettings.show_phone" 
                  type="checkbox" 
                  class="w-5 h-5 text-blue-500 rounded focus:ring-2 focus:ring-blue-500">
                <span class="text-sm text-gray-700">Hiển thị số điện thoại</span>
              </label>

              <label class="flex items-center space-x-3 cursor-pointer">
                <input 
                  v-model="shopSettings.show_custom_message" 
                  type="checkbox" 
                  class="w-5 h-5 text-blue-500 rounded focus:ring-2 focus:ring-blue-500">
                <span class="text-sm text-gray-700">Hiển thị tin nhắn tùy chỉnh</span>
              </label>
            </div>

            <!-- Auto Print -->
            <div class="border-t pt-4">
              <label class="flex items-center space-x-3 cursor-pointer">
                <input 
                  v-model="shopSettings.auto_print_enabled" 
                  type="checkbox" 
                  class="w-5 h-5 text-blue-500 rounded focus:ring-2 focus:ring-blue-500">
                <div>
                  <span class="text-sm font-medium text-gray-700">Tự động in khi thu tiền</span>
                  <p class="text-xs text-gray-500">Tự động in hóa đơn và nhãn món khi waiter thu tiền</p>
                </div>
              </label>
            </div>

            <!-- Save Button -->
            <button 
              @click="saveShopSettings" 
              :disabled="savingShopSettings"
              class="w-full bg-blue-500 text-white py-3 rounded-lg font-medium hover:bg-blue-600 active:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
              {{ savingShopSettings ? 'Đang lưu...' : 'Lưu thông tin cửa hàng' }}
            </button>
          </div>
        </div>

        <!-- Shop Settings Card -->
        <div class="bg-white rounded-2xl p-6 shadow-sm">
          <h3 class="text-lg font-bold mb-4">⚙️ Cài đặt hệ thống</h3>
          
          <div class="space-y-4">
            <!-- Low Margin Threshold -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Ngưỡng cảnh báo lợi nhuận thấp (%)
              </label>
              <input 
                v-model.number="settings.low_margin_threshold" 
                type="number" 
                min="0" 
                max="100"
                step="0.1"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="20.0">
              <p class="text-xs text-gray-500 mt-1">
                Món có lợi nhuận dưới ngưỡng này sẽ được đánh dấu cảnh báo màu vàng
              </p>
            </div>

            <!-- Save Button -->
            <button 
              @click="saveSettings" 
              :disabled="savingSettings"
              class="w-full bg-blue-500 text-white py-3 rounded-lg font-medium hover:bg-blue-600 active:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
              {{ savingSettings ? 'Đang lưu...' : 'Lưu cài đặt hệ thống' }}
            </button>
          </div>
        </div>

        <!-- Operating Expenses Section -->
        <div class="bg-white rounded-2xl p-6 shadow-sm">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-bold">💸 Chi phí vận hành</h3>
            <button 
              @click="showExpenseForm = true"
              class="bg-green-500 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-green-600 active:bg-green-700 transition-colors">
              + Thêm mới
            </button>
          </div>

          <!-- Expenses List -->
          <div v-if="expenses.length === 0" class="text-center py-8 text-gray-500">
            <div class="text-4xl mb-2">📭</div>
            <p>Chưa có chi phí vận hành nào</p>
            <p class="text-sm mt-1">Nhấn "Thêm mới" để bắt đầu</p>
          </div>

          <div v-else class="space-y-3">
            <div 
              v-for="expense in expenses" 
              :key="expense.id"
              @click="editExpense(expense)"
              class="bg-gray-50 rounded-xl p-4 active:scale-98 transition-transform cursor-pointer">
              <div class="flex justify-between items-start mb-2">
                <div>
                  <h4 class="font-bold">{{ formatDateRange(expense.period_start, expense.period_end) }}</h4>
                  <p class="text-sm text-gray-600">Tổng: {{ formatPrice(expense.total_expenses) }}</p>
                </div>
                <span class="bg-blue-100 text-blue-800 px-2 py-1 rounded-full text-xs font-medium">
                  {{ getPeriodType(expense.period_start, expense.period_end) }}
                </span>
              </div>
              <div class="grid grid-cols-2 gap-2 text-xs text-gray-600">
                <div>💼 Lương: {{ formatPrice(expense.staff_salary) }}</div>
                <div>🏢 Thuê MB: {{ formatPrice(expense.rent) }}</div>
                <div>💡 Điện nước: {{ formatPrice(expense.utilities) }}</div>
                <div>📢 Marketing: {{ formatPrice(expense.marketing_costs) }}</div>
                <div v-if="expense.other_expenses > 0" class="col-span-2">
                  📦 Khác: {{ formatPrice(expense.other_expenses) }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />

    <!-- Operating Expense Form Modal -->
    <transition name="slide-right">
      <div v-if="showExpenseForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-gray-50 w-full h-screen flex flex-col">
          <OperatingExpenseForm
            :initial-data="selectedExpense"
            @save="handleExpenseSave"
            @cancel="closeExpenseForm" />
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import OperatingExpenseForm from '../components/OperatingExpenseForm.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'
import { profitAnalysisService } from '../services/profitAnalysis'
import { shopSettingsService } from '../services/shopSettings'
import api from '../services/api'

const authStore = useAuthStore()

const loading = ref(false)
const savingSettings = ref(false)
const savingShopSettings = ref(false)
const showExpenseForm = ref(false)
const selectedExpense = ref(null)
const expenses = ref([])

const settings = ref({
  low_margin_threshold: 20.0
})

const shopSettings = ref({
  id: null,
  shop_name: '',
  shop_address: '',
  shop_phone: '',
  custom_message: '',
  show_address: true,
  show_phone: true,
  show_custom_message: true,
  auto_print_enabled: true
})

// Fetch settings
const fetchSettings = async () => {
  try {
    loading.value = true
    const response = await api.get('/settings')
    if (response.data) {
      settings.value.low_margin_threshold = response.data.low_margin_threshold || 20.0
    }
  } catch (error) {
    console.error('Error fetching settings:', error)
  } finally {
    loading.value = false
  }
}

// Fetch shop settings
const fetchShopSettings = async () => {
  try {
    const response = await shopSettingsService.getSettings()
    if (response) {
      shopSettings.value = {
        id: response.id,
        shop_name: response.shop_name || '',
        shop_address: response.shop_address || '',
        shop_phone: response.shop_phone || '',
        custom_message: response.custom_message || '',
        show_address: response.show_address !== false,
        show_phone: response.show_phone !== false,
        show_custom_message: response.show_custom_message !== false,
        auto_print_enabled: response.auto_print_enabled !== false
      }
    }
  } catch (error) {
    console.error('Error fetching shop settings:', error)
  }
}

// Save settings
const saveSettings = async () => {
  try {
    savingSettings.value = true
    await api.patch('/settings', {
      low_margin_threshold: settings.value.low_margin_threshold
    })
    alert('Cài đặt hệ thống đã được lưu!')
  } catch (error) {
    console.error('Error saving settings:', error)
    alert('Lỗi khi lưu cài đặt: ' + (error.response?.data?.error || error.message))
  } finally {
    savingSettings.value = false
  }
}

// Save shop settings
const saveShopSettings = async () => {
  try {
    if (!shopSettings.value.shop_name) {
      alert('Vui lòng nhập tên cửa hàng')
      return
    }

    savingShopSettings.value = true
    await shopSettingsService.updateSettings(shopSettings.value.id, {
      shop_name: shopSettings.value.shop_name,
      shop_address: shopSettings.value.shop_address,
      shop_phone: shopSettings.value.shop_phone,
      custom_message: shopSettings.value.custom_message,
      show_address: shopSettings.value.show_address,
      show_phone: shopSettings.value.show_phone,
      show_custom_message: shopSettings.value.show_custom_message,
      auto_print_enabled: shopSettings.value.auto_print_enabled
    })
    alert('Thông tin cửa hàng đã được lưu!')
  } catch (error) {
    console.error('Error saving shop settings:', error)
    alert('Lỗi khi lưu thông tin: ' + (error.response?.data?.error || error.message))
  } finally {
    savingShopSettings.value = false
  }
}

// Fetch operating expenses
const fetchExpenses = async () => {
  try {
    const data = await profitAnalysisService.getOperatingExpenses()
    expenses.value = data.sort((a, b) => 
      new Date(b.period_start) - new Date(a.period_start)
    )
  } catch (error) {
    console.error('Error fetching expenses:', error)
  }
}

// Edit expense
const editExpense = (expense) => {
  selectedExpense.value = expense
  showExpenseForm.value = true
}

// Handle expense save
const handleExpenseSave = async () => {
  closeExpenseForm()
  await fetchExpenses()
}

// Close expense form
const closeExpenseForm = () => {
  showExpenseForm.value = false
  selectedExpense.value = null
}

// Format price
const formatPrice = (price) => {
  return new Intl.NumberFormat('vi-VN', { 
    style: 'currency', 
    currency: 'VND',
    maximumFractionDigits: 0
  }).format(price)
}

// Format date range
const formatDateRange = (start, end) => {
  const startDate = new Date(start)
  const endDate = new Date(end)
  
  if (startDate.toDateString() === endDate.toDateString()) {
    return startDate.toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit', year: 'numeric' })
  }
  
  return `${startDate.toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit' })} - ${endDate.toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit', year: 'numeric' })}`
}

// Get period type
const getPeriodType = (start, end) => {
  const startDate = new Date(start)
  const endDate = new Date(end)
  const diffDays = Math.ceil((endDate - startDate) / (1000 * 60 * 60 * 24))
  
  if (diffDays === 0) return 'Ngày'
  if (diffDays <= 7) return 'Tuần'
  if (diffDays <= 31) return 'Tháng'
  return 'Kỳ'
}

// Refresh data
const refreshData = async () => {
  await Promise.all([
    fetchSettings(),
    fetchShopSettings(),
    fetchExpenses()
  ])
}

// Pull to refresh
const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)

// Lifecycle
onMounted(async () => {
  await refreshData()
})
</script>

<style scoped>
.active\:scale-95:active {
  transform: scale(0.95);
}

.active\:scale-98:active {
  transform: scale(0.98);
}

.slide-right-enter-active,
.slide-right-leave-active {
  transition: transform 0.3s ease;
}

.slide-right-enter-from {
  transform: translateX(100%);
}

.slide-right-leave-to {
  transform: translateX(100%);
}
</style>
