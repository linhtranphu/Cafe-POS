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
        <!-- Shop Settings Card -->
        <div class="bg-white rounded-2xl p-6 shadow-sm">
          <h3 class="text-lg font-bold mb-4">🏪 Cài đặt cửa hàng</h3>
          
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
              {{ savingSettings ? 'Đang lưu...' : 'Lưu cài đặt' }}
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
import api from '../services/api'

const authStore = useAuthStore()

const loading = ref(false)
const savingSettings = ref(false)
const showExpenseForm = ref(false)
const selectedExpense = ref(null)
const expenses = ref([])

const settings = ref({
  low_margin_threshold: 20.0
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

// Save settings
const saveSettings = async () => {
  try {
    savingSettings.value = true
    await api.patch('/settings', {
      low_margin_threshold: settings.value.low_margin_threshold
    })
    alert('Cài đặt đã được lưu!')
  } catch (error) {
    console.error('Error saving settings:', error)
    alert('Lỗi khi lưu cài đặt: ' + (error.response?.data?.error || error.message))
  } finally {
    savingSettings.value = false
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
