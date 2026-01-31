<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm">
      <div class="px-4 py-3">
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
        
        <!-- Source Type Filter -->
        <div class="mt-2 flex gap-2 overflow-x-auto pb-2">
          <button @click="sourceFilter = ''" 
            :class="sourceFilter === '' ? 'bg-purple-500 text-white' : 'bg-white text-gray-700 border border-gray-300'"
            class="px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap">
            Tất cả
          </button>
          <button @click="sourceFilter = 'manual'" 
            :class="sourceFilter === 'manual' ? 'bg-purple-500 text-white' : 'bg-white text-gray-700 border border-gray-300'"
            class="px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap">
            ✍️ Thủ công
          </button>
          <button @click="sourceFilter = 'ingredient'" 
            :class="sourceFilter === 'ingredient' ? 'bg-green-500 text-white' : 'bg-white text-gray-700 border border-gray-300'"
            class="px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap">
            🥬 Nguyên liệu
          </button>
          <button @click="sourceFilter = 'facility'" 
            :class="sourceFilter === 'facility' ? 'bg-blue-500 text-white' : 'bg-white text-gray-700 border border-gray-300'"
            class="px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap">
            🏢 Cơ sở vật chất
          </button>
          <button @click="sourceFilter = 'maintenance'" 
            :class="sourceFilter === 'maintenance' ? 'bg-orange-500 text-white' : 'bg-white text-gray-700 border border-gray-300'"
            class="px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap">
            🔧 Bảo trì
          </button>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="px-4 py-4 pb-24">
      <!-- Stats Cards - Single Row -->
      <div class="bg-gradient-to-br from-purple-500 to-pink-500 rounded-xl p-4 mb-4 text-white shadow-lg">
        <div class="text-xs opacity-90 mb-2">Tổng quan chi phí</div>
        <div class="grid grid-cols-4 gap-1.5">
          <div class="text-center">
            <div class="text-lg font-bold">{{ expenses.length }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap">Tổng</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold">{{ formatPrice(totalThisMonth) }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap">Tháng này</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold">{{ recurringCount }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap">Định kỳ</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold">{{ categories.length }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap">Danh mục</div>
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
            <div class="text-sm font-bold">Quản lý danh mục</div>
          </button>
          <button @click="showRecurringModal = true"
            class="bg-gradient-to-br from-green-500 to-emerald-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform">
            <div class="text-2xl mb-1">🔁</div>
            <div class="text-sm font-bold">Chi phí định kỳ</div>
          </button>
          <button @click="showReportModal = true"
            class="bg-gradient-to-br from-orange-500 to-red-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform">
            <div class="text-2xl mb-1">📊</div>
            <div class="text-sm font-bold">Báo cáo</div>
          </button>
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
            @click="viewDetails(expense)"
            class="bg-white rounded-2xl p-4 shadow-sm active:scale-98 transition-transform">
            
            <!-- Expense Header -->
            <div class="flex justify-between items-start mb-3">
              <div>
                <div class="flex items-center gap-2 mb-1">
                  <h3 class="font-bold text-lg">{{ expense.description }}</h3>
                  <span v-if="expense.source_type && expense.source_type !== 'manual'" 
                    :class="getSourceTypeBadgeClass(expense.source_type)"
                    class="px-2 py-0.5 rounded text-[10px] font-medium">
                    {{ getSourceTypeLabel(expense.source_type) }}
                  </span>
                </div>
                <p class="text-sm text-gray-600">{{ getCategoryName(expense.category_id) }}</p>
                <p class="text-xs text-gray-400">{{ formatDate(expense.date) }}</p>
              </div>
              <div class="text-right">
                <div class="text-lg font-bold text-red-600">-{{ formatPrice(expense.amount) }}</div>
                <span class="px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                  {{ getPaymentMethodLabel(expense.payment_method) }}
                </span>
              </div>
            </div>

            <!-- Expense Info -->
            <div v-if="expense.vendor" class="mb-3 text-sm">
              <div class="flex items-center gap-2 text-gray-600">
                <span>🏪</span>
                <span>{{ expense.vendor }}</span>
              </div>
            </div>

            <!-- Quick Actions -->
            <div class="flex gap-2 pt-3 border-t">
              <button @click.stop="openEditModal(expense)"
                class="flex-1 bg-blue-500 text-white py-2 rounded-lg text-sm font-medium active:bg-blue-600">
                ✏️ Sửa
              </button>
              <button @click.stop="deleteExpense(expense)"
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

    <!-- Category Management Modal -->
    <transition name="slide-up">
      <div v-if="showCategoryModal" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full max-h-[85vh] overflow-y-auto">
          <div class="sticky top-0 bg-white px-4 py-4 border-b flex justify-between items-center">
            <h3 class="text-lg font-bold">📁 Quản lý danh mục</h3>
            <button @click="showCategoryModal = false" class="text-2xl text-gray-400">×</button>
          </div>
          
          <div class="px-4 py-4">
            <!-- Add New Category -->
            <div class="bg-gray-50 rounded-xl p-4 mb-4">
              <h4 class="font-semibold text-gray-800 mb-3">Thêm danh mục mới</h4>
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
                class="bg-white border border-gray-200 rounded-xl p-4 flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-lg bg-purple-100 text-purple-600 flex items-center justify-center text-xl">
                    💰
                  </div>
                  <div>
                    <div class="font-medium text-gray-800">{{ cat.name }}</div>
                    <div class="text-xs text-gray-500">{{ getCategoryCount(cat.id) }} chi phí</div>
                  </div>
                </div>
                <button @click="deleteCategory(cat.id)" class="text-red-500 hover:text-red-700 p-2">
                  🗑️
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- Create/Edit Modal -->
    <transition name="slide-up">
      <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full max-h-[90vh] overflow-y-auto">
          <div class="sticky top-0 bg-white px-4 py-4 border-b flex justify-between items-center">
            <h3 class="text-lg font-bold">{{ isEditing ? 'Cập nhật chi phí' : 'Thêm chi phí mới' }}</h3>
            <button @click="closeModal" class="text-2xl text-gray-400">×</button>
          </div>
          
          <div class="px-4 py-4 space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Mô tả *</label>
              <input v-model="formData.description" type="text" 
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500" />
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Danh mục *</label>
                <select v-model="formData.category_id" 
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500">
                  <option value="">Chọn danh mục</option>
                  <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Số tiền *</label>
                <input v-model.number="formData.amount" type="number" 
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500" />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Ngày *</label>
                <input v-model="formData.date" type="date" 
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Thanh toán *</label>
                <select v-model="formData.payment_method" 
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500">
                  <option v-for="option in PAYMENT_METHOD_OPTIONS" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </option>
                </select>
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Nhà cung cấp</label>
              <input v-model="formData.vendor" type="text" 
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500" />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Ghi chú</label>
              <textarea v-model="formData.notes" rows="3" 
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"></textarea>
            </div>

            <div class="flex gap-3 pt-4">
              <button @click="closeModal" 
                class="flex-1 bg-gray-200 text-gray-700 py-3 rounded-xl font-medium active:bg-gray-300">
                Hủy
              </button>
              <button @click="saveExpense" 
                class="flex-1 bg-blue-500 text-white py-3 rounded-xl font-medium active:bg-blue-600">
                {{ isEditing ? 'Cập nhật' : 'Thêm mới' }}
              </button>
            </div>
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
import { formatDate, formatPrice } from '../utils/formatters'
import { PAYMENT_METHODS, PAYMENT_METHOD_OPTIONS, getPaymentMethodLabel } from '../constants/expense'

const expenseStore = useExpenseStore()

const searchQuery = ref('')
const sourceFilter = ref('')
const showModal = ref(false)
const showCategoryModal = ref(false)
const showRecurringModal = ref(false)
const showReportModal = ref(false)
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

const filteredExpenses = computed(() => {
  let filtered = expenses.value
  
  // Filter by source type
  if (sourceFilter.value) {
    filtered = filtered.filter(e => {
      if (sourceFilter.value === 'manual') {
        return !e.source_type || e.source_type === 'manual'
      }
      return e.source_type === sourceFilter.value
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
  
  return filtered
})

const totalThisMonth = computed(() => {
  const now = new Date()
  const thisMonth = now.getMonth()
  const thisYear = now.getFullYear()
  
  return expenses.value
    .filter(e => {
      const expenseDate = new Date(e.date)
      return expenseDate.getMonth() === thisMonth && expenseDate.getFullYear() === thisYear
    })
    .reduce((sum, e) => sum + e.amount, 0)
})

const recurringCount = computed(() => {
  return expenseStore.recurringExpenses?.length || 0
})

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
  showModal.value = true
}

const openEditModal = (expense) => {
  isEditing.value = true
  currentExpense.value = expense
  formData.value = {
    ...expense,
    date: new Date(expense.date).toISOString().split('T')[0]
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  isEditing.value = false
  currentExpense.value = null
}

const saveExpense = async () => {
  try {
    if (isEditing.value) {
      await expenseStore.updateExpense(currentExpense.value.id, formData.value)
      alert('Cập nhật chi phí thành công')
    } else {
      await expenseStore.createExpense(formData.value)
      alert('Thêm chi phí thành công')
    }
    closeModal()
  } catch (error) {
    console.error('Error saving expense:', error)
    alert('Có lỗi xảy ra khi lưu chi phí')
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

const viewDetails = (expense) => {
  openEditModal(expense)
}

onMounted(async () => {
  await expenseStore.fetchCategories()
  await expenseStore.fetchExpenses()
  await expenseStore.fetchRecurringExpenses()
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
</style>
