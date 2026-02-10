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
        <div class="flex items-center justify-between mb-3">
          <h1 class="text-xl font-bold text-gray-800">💰 Chi phí</h1>
        </div>
        
        <!-- Search Bar -->
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Tìm kiếm chi phí..."
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
        
        <!-- Creator Filter -->
        <div class="mt-2 flex gap-2 overflow-x-auto pb-2">
          <button @click="creatorFilter = ''" 
            :class="creatorFilter === '' ? 'bg-purple-500 text-white' : 'bg-white text-gray-700 border border-gray-300'"
            class="px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap">
            👥 Tất cả
          </button>
          <button v-for="creator in uniqueCreators" :key="creator"
            @click="creatorFilter = creator" 
            :class="creatorFilter === creator ? 'bg-purple-500 text-white' : 'bg-white text-gray-700 border border-gray-300'"
            class="px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap">
            👤 {{ creator }}
          </button>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24">
      <!-- Stats Cards -->
      <div class="bg-gradient-to-br from-purple-500 to-pink-500 rounded-xl p-4 mb-4 text-white shadow-lg">
        <div class="flex items-center justify-between mb-2">
          <div class="text-xs opacity-90">Chi phí</div>
          <div v-if="creatorFilter" class="text-xs opacity-90 bg-white/20 px-2 py-1 rounded-full">
            👤 {{ creatorFilter }}
          </div>
        </div>
        <div class="grid grid-cols-3 gap-3">
          <div class="text-center">
            <div class="text-base font-bold leading-tight">{{ formatCompactPrice(totalAllTime) }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap mt-1">Tổng từ đầu</div>
          </div>
          <div class="text-center">
            <div class="text-base font-bold leading-tight">{{ formatCompactPrice(totalThisMonth) }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap mt-1">Tháng này</div>
          </div>
          <div class="text-center">
            <div class="text-base font-bold leading-tight">{{ recurringCount }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap mt-1">Định kỳ</div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="mb-4">
        <h2 class="text-sm font-bold text-gray-800 mb-2">⚡ Thao tác nhanh</h2>
        <div class="grid grid-cols-2 gap-2">
          <button @click="openCreateModal"
            class="bg-gradient-to-br from-blue-500 to-cyan-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform">
            <div class="text-2xl mb-1">➕</div>
            <div class="text-sm font-bold">Tạo chi phí</div>
          </button>
          <button @click="showCategoryModal = true"
            class="bg-gradient-to-br from-purple-500 to-pink-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform">
            <div class="text-2xl mb-1">📁</div>
            <div class="text-sm font-bold">Danh mục</div>
          </button>
        </div>
      </div>

      <!-- Category Management Modal -->
      <transition name="slide-up">
        <div v-if="showCategoryModal" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
          <div class="bg-white rounded-t-3xl w-full h-[85vh] flex flex-col">
            <!-- Fixed Header -->
            <div class="flex-shrink-0 bg-white px-4 py-4 border-b flex justify-between items-center rounded-t-3xl">
              <h3 class="text-lg font-bold">📁 Quản lý Danh mục</h3>
              <button @click="showCategoryModal = false" class="text-2xl text-gray-400">×</button>
            </div>
            
            <!-- Scrollable Content -->
            <div class="flex-1 overflow-y-auto px-4 py-4">
              <!-- Add New Category -->
              <div class="bg-gray-50 rounded-xl p-4 mb-4 flex-shrink-0">
                <h4 class="font-semibold text-gray-800 mb-3">Thêm danh mục mới</h4>
                <form @submit.prevent="addCategory" class="flex flex-col sm:flex-row gap-2">
                  <input v-model="newCategoryName" type="text" required placeholder="VD: Tiền điện, Tiền nước..." 
                    class="flex-1 px-4 py-3 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500" />
                  <button type="submit" class="bg-purple-500 text-white px-6 py-3 rounded-lg font-medium text-base active:bg-purple-600 whitespace-nowrap">
                    Thêm
                  </button>
                </form>
              </div>

              <!-- Category List -->
              <div class="space-y-3 pb-4">
                <div v-for="cat in categories" :key="cat.id" 
                  class="bg-white border border-gray-200 rounded-xl p-4 flex items-center justify-between">
                  <div class="flex items-center gap-3 flex-1 min-w-0">
                    <div class="w-12 h-12 rounded-lg bg-purple-100 text-purple-600 flex items-center justify-center text-2xl flex-shrink-0">
                      💰
                    </div>
                    <div class="min-w-0">
                      <div class="font-medium text-gray-800 truncate">{{ cat.name }}</div>
                      <div class="text-xs text-gray-500">{{ getCategoryCount(cat.id) }} chi phí</div>
                    </div>
                  </div>
                  <button @click="deleteCategory(cat.id)" class="text-red-500 hover:text-red-700 p-2 flex-shrink-0 ml-2">
                    🗑️
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </transition>

      <!-- Category Management Form (old inline version - removed) -->
      <div v-if="false" class="bg-white rounded-2xl p-4 mb-4 shadow-md border-2 border-purple-200">
        <h3 class="text-lg font-bold mb-4">📁 Quản lý danh mục</h3>
        
        <!-- Add New Category -->
        <div class="mb-4">
          <div class="flex gap-2">
            <input v-model="newCategoryName" type="text" placeholder="Tên danh mục..." 
              class="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500" />
            <button @click="addCategory" class="bg-purple-500 text-white px-6 py-3 rounded-lg font-medium active:bg-purple-600">
              Thêm
            </button>
          </div>
        </div>

        <!-- Category List -->
        <div class="space-y-2">
          <div v-for="cat in categories" :key="cat.id" 
            class="bg-gray-50 border border-gray-200 rounded-xl p-3 flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-lg bg-purple-100 text-purple-600 flex items-center justify-center text-lg">
                💰
              </div>
              <div>
                <div class="font-medium text-gray-800 text-sm">{{ cat.name }}</div>
                <div class="text-xs text-gray-500">{{ getCategoryCount(cat.id) }} chi phí</div>
              </div>
            </div>
            <button @click="deleteCategory(cat.id)" class="text-red-500 hover:text-red-700 p-2">
              🗑️
            </button>
          </div>
        </div>
      </div>

      <!-- Expenses List -->
      <div class="mb-4">
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-lg font-bold text-gray-800">📋 Danh sách chi phí</h2>
        </div>
        
        <div v-if="filteredExpenses.length === 0" class="text-center py-16">
          <div class="text-6xl mb-4">📭</div>
          <p class="text-gray-500">Không có chi phí nào</p>
        </div>
        
        <div v-else class="space-y-3">
          <div v-for="expense in filteredExpenses" :key="expense.id"
            class="bg-white rounded-2xl p-4 shadow-sm border-l-4 border-purple-500">
            
            <!-- Expense Header -->
            <div class="flex justify-between items-start mb-3">
              <div class="flex-1">
                <div class="flex items-center gap-2 mb-1">
                  <h3 class="font-bold text-lg">{{ expense.description }}</h3>
                  <span v-if="expense.source_type && expense.source_type !== 'manual'" 
                    :class="getSourceTypeBadgeClass(expense.source_type)"
                    class="px-2 py-0.5 rounded text-[10px] font-medium">
                    {{ getSourceTypeLabel(expense.source_type) }}
                  </span>
                </div>
                <p class="text-sm text-gray-600">{{ getCategoryName(expense.category_id) }}</p>
              </div>
              <div class="text-right">
                <div class="text-lg font-bold text-red-600">-{{ formatPrice(expense.amount) }}</div>
              </div>
            </div>

            <!-- Expense Details Grid -->
            <div class="grid grid-cols-2 gap-3 mb-3 text-sm">
              <!-- Date & Payment -->
              <div class="flex items-center gap-2 text-gray-600">
                <span>📅</span>
                <span>{{ formatDate(expense.date) }}</span>
              </div>
              <div class="flex items-center gap-2 text-gray-600">
                <span>💳</span>
                <span class="px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                  {{ getPaymentMethodLabel(expense.payment_method) }}
                </span>
              </div>

              <!-- Vendor -->
              <div v-if="expense.vendor" class="flex items-center gap-2 text-gray-600">
                <span>🏪</span>
                <span>{{ expense.vendor }}</span>
              </div>

              <!-- Creator -->
              <div class="flex items-center gap-2 text-gray-600">
                <span>👤</span>
                <span class="font-medium">{{ expense.created_by || 'Hệ thống' }}</span>
              </div>
            </div>

            <!-- Notes -->
            <div v-if="expense.notes" class="mb-3 p-2 bg-gray-50 rounded-lg text-sm text-gray-600 border-l-2 border-gray-300">
              <span class="text-xs font-semibold text-gray-500">Ghi chú:</span>
              <p>{{ expense.notes }}</p>
            </div>

            <!-- Quick Actions -->
            <div class="flex gap-2 pt-3 border-t">
              <button @click="openEditModal(expense)"
                class="flex-1 bg-blue-500 text-white py-2 rounded-lg text-sm font-medium active:bg-blue-600">
                ✏️ Sửa
              </button>
              <button @click="deleteExpense(expense)"
                class="flex-1 bg-red-500 text-white py-2 rounded-lg text-sm font-medium active:bg-red-600">
                🗑️ Xóa
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />

    <!-- Create/Edit Expense Form Modal - Slide from Right -->
    <transition name="slide-right">
      <div v-if="showCreateForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-gray-50 w-full h-screen flex flex-col">
          <!-- Mobile Header - Fixed -->
          <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
            <div class="px-4 py-3">
              <div class="flex items-center justify-between">
                <button @click="cancelEdit" class="text-2xl text-gray-600">←</button>
                <h1 class="text-xl font-bold text-gray-800">{{ isEditing ? '✏️ Cập nhật chi phí' : '➕ Thêm chi phí mới' }}</h1>
                <div class="w-8"></div>
              </div>
            </div>
          </div>

          <!-- Scrollable Content -->
          <div class="flex-1 overflow-y-auto px-4 py-6 space-y-5">
            <!-- Mô tả -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">Mô tả *</label>
              <input v-model="formData.description" type="text" required placeholder="VD: Tiền điện tháng 1"
                class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
            </div>

            <!-- Danh mục & Số tiền -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">
                  Danh mục <span class="text-red-500">*</span>
                </label>
                <select v-model="formData.category_id" required
                  :class="{'border-red-500': !formData.category_id}"
                  class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent">
                  <option value="">Chọn danh mục</option>
                  <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                </select>
                <p v-if="!formData.category_id && categories.length === 0" class="text-xs text-orange-600 mt-1">
                  ⚠️ Chưa có danh mục. Tạo danh mục trước!
                </p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">Số tiền (VNĐ) *</label>
                <input v-model.number="formData.amount" type="number" min="0" step="1000" required placeholder="0"
                  class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
              </div>
            </div>

            <!-- Ngày & Thanh toán -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">Ngày *</label>
                <input v-model="formData.date" type="date" required
                  class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent appearance-none" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">Thanh toán *</label>
                <select v-model="formData.payment_method" required
                  class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent">
                  <option v-for="option in PAYMENT_METHOD_OPTIONS" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </option>
                </select>
              </div>
            </div>

            <!-- Nhà cung cấp -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">Nhà cung cấp</label>
              <input v-model="formData.vendor" type="text" placeholder="VD: Công ty điện lực"
                class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
            </div>

            <!-- Ghi chú -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">Ghi chú</label>
              <textarea v-model="formData.notes" rows="3" placeholder="Ghi chú thêm..."
                class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"></textarea>
            </div>

            <!-- Spacer for bottom buttons -->
            <div class="h-24"></div>
          </div>

          <!-- Fixed Footer -->
          <div class="flex-shrink-0 bg-white px-4 py-4 border-t flex gap-3 pb-safe">
            <button @click="cancelEdit" 
              class="flex-1 bg-gray-200 text-gray-700 py-4 rounded-xl font-medium text-base active:bg-gray-300 transition-colors">
              Hủy
            </button>
            <button @click="saveExpense" :disabled="!formData.description || !formData.category_id || formData.amount <= 0"
              class="flex-1 bg-green-500 text-white py-4 rounded-xl font-medium text-base active:bg-green-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
              {{ isEditing ? 'Cập nhật' : 'Thêm chi phí' }}
            </button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useExpenseStore } from '../stores/expense'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'
import { formatDate, formatPrice } from '../utils/formatters'
import { PAYMENT_METHODS, PAYMENT_METHOD_OPTIONS, getPaymentMethodLabel } from '../constants/expense'

const expenseStore = useExpenseStore()

const searchQuery = ref('')
const creatorFilter = ref('')
const showCreateForm = ref(false)
const showCategoryModal = ref(false)
const isEditing = ref(false)
const currentExpense = ref(null)
const newCategoryName = ref('')

const formData = ref({
  description: '',
  category_id: '',
  amount: 0,
  date: new Date().toISOString().split('T')[0],
  payment_method: PAYMENT_METHODS.CASH,
  vendor: '',
  notes: ''
})

const expenses = computed(() => expenseStore.expenses || [])
const categories = computed(() => expenseStore.categories || [])

// Get unique creators from expenses
const uniqueCreators = computed(() => {
  const creators = expenses.value
    .map(e => e.created_by || 'Hệ thống')
    .filter((value, index, self) => self.indexOf(value) === index)
  return creators.sort()
})

const filteredExpenses = computed(() => {
  let filtered = expenses.value
  
  // Filter by creator
  if (creatorFilter.value) {
    filtered = filtered.filter(e => {
      const creator = e.created_by || 'Hệ thống'
      return creator === creatorFilter.value
    })
  }
  
  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(e => 
      e.description?.toLowerCase().includes(query) ||
      e.vendor?.toLowerCase().includes(query)
    )
  }
  
  // Sort by date (newest first)
  return [...filtered].sort((a, b) => {
    const dateA = new Date(a.date || a.created_at || 0)
    const dateB = new Date(b.date || b.created_at || 0)
    return dateB - dateA // Newest first
  })
})

const totalThisMonth = computed(() => {
  const now = new Date()
  const thisMonth = now.getMonth()
  const thisYear = now.getFullYear()
  
  let filtered = expenses.value
  
  // Filter by creator if selected
  if (creatorFilter.value) {
    filtered = filtered.filter(e => {
      const creator = e.created_by || 'Hệ thống'
      return creator === creatorFilter.value
    })
  }
  
  return filtered
    .filter(e => {
      const expenseDate = new Date(e.date)
      return expenseDate.getMonth() === thisMonth && expenseDate.getFullYear() === thisYear
    })
    .reduce((sum, e) => sum + e.amount, 0)
})

const totalAllTime = computed(() => {
  let filtered = expenses.value
  
  // Filter by creator if selected
  if (creatorFilter.value) {
    filtered = filtered.filter(e => {
      const creator = e.created_by || 'Hệ thống'
      return creator === creatorFilter.value
    })
  }
  
  return filtered.reduce((sum, e) => sum + e.amount, 0)
})

const recurringCount = computed(() => {
  return expenseStore.recurringExpenses?.length || 0
})

// Format price in compact form - always show in thousands (÷1000)
const formatCompactPrice = (value) => {
  if (value === undefined || value === null || isNaN(value)) {
    return '0k'
  }
  
  // Always divide by 1000 to show in thousands
  const thousands = value / 1000
  
  // If it's a whole number of thousands
  if (thousands % 1 === 0) {
    return `${thousands}k`
  }
  
  // If it has decimals, show 1 decimal place
  return `${thousands.toFixed(1)}k`
}

const getCategoryName = (categoryId) => {
  const category = categories.value.find(c => c.id === categoryId)
  return category ? category.name : 'Không xác định'
}

const getCategoryCount = (categoryId) => {
  return expenses.value.filter(e => e.category_id === categoryId).length
}

const getSourceTypeLabel = (sourceType) => {
  const labels = {
    ingredient: '🥬 Tự động',
    facility: '🏢 Tự động',
    maintenance: '🔧 Tự động',
    manual: '✍️ Thủ công'
  }
  return labels[sourceType] || ''
}

const getSourceTypeBadgeClass = (sourceType) => {
  const classes = {
    ingredient: 'bg-green-100 text-green-700',
    facility: 'bg-blue-100 text-blue-700',
    maintenance: 'bg-orange-100 text-orange-700',
    manual: 'bg-gray-100 text-gray-700'
  }
  return classes[sourceType] || 'bg-gray-100 text-gray-700'
}

const openCreateModal = () => {
  cancelEdit()
  showCreateForm.value = true
}

const toggleCategoryForm = () => {
  showCategoryModal.value = !showCategoryModal.value
}

const openEditModal = (expense) => {
  isEditing.value = true
  currentExpense.value = expense
  formData.value = {
    ...expense,
    date: new Date(expense.date).toISOString().split('T')[0]
  }
  showCreateForm.value = true
}

const cancelEdit = () => {
  showCreateForm.value = false
  isEditing.value = false
  currentExpense.value = null
  formData.value = {
    description: '',
    category_id: '',
    amount: 0,
    date: new Date().toISOString().split('T')[0],
    payment_method: PAYMENT_METHODS.CASH,
    vendor: '',
    notes: ''
  }
}

const saveExpense = async () => {
  try {
    // Validate required fields
    if (!formData.value.description?.trim()) {
      alert('Vui lòng nhập mô tả')
      return
    }
    if (!formData.value.category_id) {
      alert('Vui lòng chọn danh mục')
      return
    }
    if (!formData.value.amount || formData.value.amount <= 0) {
      alert('Vui lòng nhập số tiền hợp lệ')
      return
    }
    if (!formData.value.date) {
      alert('Vui lòng chọn ngày')
      return
    }
    
    // Prepare data - convert date to ISO format
    const dataToSend = { ...formData.value }
    if (dataToSend.date) {
      dataToSend.date = dataToSend.date + 'T00:00:00Z'
    } else {
      delete dataToSend.date
    }
    
    console.log('Sending expense data:', dataToSend)
    
    if (isEditing.value) {
      await expenseStore.updateExpense(currentExpense.value.id, dataToSend)
      alert('Cập nhật chi phí thành công')
    } else {
      await expenseStore.createExpense(dataToSend)
      alert('Thêm chi phí thành công')
    }
    showCreateForm.value = false
    cancelEdit()
  } catch (error) {
    console.error('Error saving expense:', error)
    console.error('Error response:', error.response?.data)
    alert(`Có lỗi xảy ra: ${error.response?.data?.error || error.message}`)
  }
}

const deleteExpense = async (expense) => {
  if (confirm(`Bạn có chắc muốn xóa chi phí "${expense.description}"?`)) {
    try {
      await expenseStore.deleteExpense(expense.id)
      alert('Xóa chi phí thành công')
    } catch (error) {
      console.error('Error deleting expense:', error)
      alert('Có lỗi xảy ra khi xóa chi phí')
    }
  }
}

const addCategory = async () => {
  if (!newCategoryName.value.trim()) return
  try {
    await expenseStore.createCategory({ name: newCategoryName.value.trim() })
    newCategoryName.value = ''
  } catch (error) {
    console.error('Error adding category:', error)
    alert('Có lỗi xảy ra khi thêm danh mục')
  }
}

const deleteCategory = async (categoryId) => {
  const hasExpenses = expenses.value.some(e => e.category_id === categoryId)
  if (hasExpenses) {
    alert('Không thể xóa danh mục đã có chi phí!')
    return
  }
  
  if (confirm('Bạn có chắc muốn xóa danh mục này?')) {
    try {
      await expenseStore.deleteCategory(categoryId)
    } catch (error) {
      console.error('Error deleting category:', error)
      alert('Có lỗi xảy ra khi xóa danh mục')
    }
  }
}

// Refresh data function
const refreshData = async () => {
  await expenseStore.fetchCategories()
  await expenseStore.fetchExpenses()
  await expenseStore.fetchRecurringExpenses()
}

// Pull to refresh
const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)

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

.pb-safe {
  padding-bottom: max(1rem, env(safe-area-inset-bottom));
}

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
