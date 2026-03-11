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
          <div class="flex items-center gap-2">
            <button @click="$router.back()" class="text-2xl text-gray-600">←</button>
            <h1 class="text-xl font-bold text-gray-800">⚡ Chi phí vận hành</h1>
          </div>
        </div>

        <!-- Search Bar -->
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Tìm kiếm chi phí..."
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-transparent"
        />
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24">
      <!-- Stats Cards -->
      <div class="bg-gradient-to-br from-orange-500 to-red-500 rounded-xl p-4 mb-4 text-white shadow-lg">
        <div class="flex items-center justify-between mb-2">
          <div class="text-xs opacity-90">Quỹ vận hành</div>
          <div v-if="operatingFundBalance !== null" class="text-xs opacity-90 bg-white/20 px-2 py-1 rounded-full">
            Số dư: {{ formatPrice(operatingFundBalance) }}
          </div>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div class="text-center">
            <div class="text-base font-bold leading-tight">{{ formatCompact(totalThisMonth) }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap mt-1">Tháng này</div>
          </div>
          <div class="text-center">
            <div class="text-base font-bold leading-tight">{{ formatCompact(totalAllTime) }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap mt-1">Tổng từ đầu</div>
          </div>
        </div>
      </div>

      <!-- Operating Fund Balance -->
      <div class="bg-blue-50 border border-blue-200 rounded-xl p-3 mb-4">
        <div class="text-xs font-semibold text-blue-700 mb-1">⚙️ Quỹ vận hành (Operating Fund)</div>
        <div v-if="fundLoading" class="text-xs text-gray-500">Đang tải...</div>
        <div v-else-if="operatingFundBalance !== null" class="text-sm text-blue-800">
          💵 Tiền mặt: <span class="font-bold">{{ formatPrice(operatingCashBalance || 0) }}</span>
          &nbsp;|&nbsp;
          💳 CK: <span class="font-bold">{{ formatPrice(operatingTransferBalance || 0) }}</span>
        </div>
        <div v-else class="text-xs text-red-600">Không thể tải số dư quỹ</div>
      </div>

      <!-- Quick Actions -->
      <div class="mb-4">
        <h2 class="text-sm font-bold text-gray-800 mb-2">⚡ Thao tác nhanh</h2>
        <div class="grid grid-cols-2 gap-2">
          <button @click="openCreateModal"
            class="bg-gradient-to-br from-orange-500 to-red-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform">
            <div class="text-2xl mb-1">➕</div>
            <div class="text-sm font-bold">Thêm chi phí</div>
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
                  <input v-model="newCategoryName" type="text" required placeholder="VD: Điện, Nước, Thuê mặt bằng..."
                    class="flex-1 px-4 py-3 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500" />
                  <button type="submit" class="bg-orange-500 text-white px-6 py-3 rounded-lg font-medium text-base active:bg-orange-600 whitespace-nowrap">
                    Thêm
                  </button>
                </form>
              </div>

              <!-- Category List -->
              <div class="space-y-3 pb-4">
                <div v-if="categories.length === 0" class="text-center py-8 text-gray-500">
                  <div class="text-4xl mb-2">📭</div>
                  <p>Chưa có danh mục nào</p>
                </div>
                <div v-for="cat in categories" :key="cat.id"
                  class="bg-white border border-gray-200 rounded-xl p-4 flex items-center justify-between">
                  <div class="flex items-center gap-3 flex-1 min-w-0">
                    <div class="w-12 h-12 rounded-lg bg-orange-100 text-orange-600 flex items-center justify-center text-2xl flex-shrink-0">
                      ⚡
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

      <!-- Expenses List -->
      <div class="mb-4">
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-lg font-bold text-gray-800">📋 Danh sách chi phí</h2>
        </div>

        <div v-if="filteredEntries.length === 0" class="text-center py-16">
          <div class="text-6xl mb-4">📭</div>
          <p class="text-gray-500">Không có chi phí nào</p>
        </div>

        <div v-else class="space-y-3">
          <div v-for="entry in filteredEntries" :key="entry.id"
            class="bg-white rounded-2xl p-4 shadow-sm border-l-4 border-orange-500">

            <!-- Entry Header -->
            <div class="flex justify-between items-start mb-2">
              <div class="flex-1">
                <h3 class="font-bold text-base">{{ getEntryDescription(entry) }}</h3>
                <p class="text-xs text-gray-500">{{ formatDate(entry.timestamp) }}</p>
              </div>
              <div class="text-right">
                <div class="text-lg font-bold text-red-600">-{{ formatPrice(getEntryAmount(entry)) }}</div>
                <div class="text-xs text-gray-500">{{ getEntryMoneyType(entry) === 'transfer' ? '💳 CK' : '💵 Tiền mặt' }}</div>
              </div>
            </div>

            <!-- Badges -->
            <div class="flex flex-wrap gap-1">
              <span class="bg-orange-100 text-orange-700 px-2 py-0.5 rounded text-[10px] font-medium">
                ⚡ Chi từ quỹ vận hành
              </span>
              <span class="bg-gray-100 text-gray-600 px-2 py-0.5 rounded text-[10px]">
                👤 {{ entry.performed_by_name || 'Hệ thống' }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />

    <!-- Create Expense Modal - Slide from bottom -->
    <transition name="slide-up">
      <div v-if="showCreateForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full max-h-[90vh] flex flex-col">
          <!-- Fixed Header -->
          <div class="flex-shrink-0 px-4 py-4 border-b flex justify-between items-center">
            <h3 class="text-lg font-bold">➕ Thêm chi phí vận hành</h3>
            <button @click="showCreateForm = false" class="text-2xl text-gray-400">×</button>
          </div>

          <!-- Scrollable Content -->
          <div class="flex-1 overflow-y-auto px-4 py-4 space-y-4">
            <!-- Operating Fund Balance -->
            <div class="bg-blue-50 border border-blue-200 rounded-xl p-3">
              <div class="text-xs font-semibold text-blue-700 mb-1">⚙️ Quỹ vận hành (Operating Fund)</div>
              <div v-if="fundLoading" class="text-xs text-gray-500">Đang tải...</div>
              <div v-else-if="operatingFundBalance !== null" class="text-sm text-blue-800">
                💵 Tiền mặt: <span class="font-bold">{{ formatPrice(operatingCashBalance || 0) }}</span>
                &nbsp;|&nbsp;
                💳 CK: <span class="font-bold">{{ formatPrice(operatingTransferBalance || 0) }}</span>
              </div>
              <div v-else class="text-xs text-red-600">Không thể tải số dư quỹ</div>
            </div>

            <!-- Description -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Mô tả *</label>
              <input v-model="formData.description" type="text" required placeholder="VD: Tiền điện tháng 3"
                class="w-full px-4 py-3 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500" />
            </div>

            <!-- Category & Amount -->
            <div class="grid grid-cols-1 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Danh mục <span class="text-red-500">*</span>
                </label>
                <select v-model="formData.category_id" required
                  class="w-full px-4 py-3 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500">
                  <option value="">Chọn danh mục</option>
                  <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                </select>
                <p v-if="categories.length === 0" class="text-xs text-orange-600 mt-1">
                  ⚠️ Chưa có danh mục. Tạo danh mục trước!
                </p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Số tiền (VNĐ) *</label>
                <input v-model.number="formData.amount" type="number" min="0" step="1000" required placeholder="0"
                  class="w-full px-4 py-3 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500" />
              </div>
            </div>

            <!-- Date -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Ngày *</label>
              <input v-model="formData.date" type="date" required
                class="w-full px-4 py-3 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 appearance-none" />
            </div>

            <!-- Money Type -->
            <div>
              <div class="text-sm font-medium text-gray-700 mb-2">Thanh toán từ quỹ bằng</div>
              <div class="grid grid-cols-2 gap-2">
                <button type="button" @click="formData.money_type = 'cash'"
                  :class="['py-3 rounded-xl font-semibold text-sm transition-colors',
                    formData.money_type !== 'transfer' ? 'bg-green-500 text-white' : 'bg-gray-100 text-gray-700']">
                  💵 Tiền mặt
                </button>
                <button type="button" @click="formData.money_type = 'transfer'"
                  :class="['py-3 rounded-xl font-semibold text-sm transition-colors',
                    formData.money_type === 'transfer' ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-700']">
                  💳 Chuyển khoản
                </button>
              </div>
            </div>

            <!-- Vendor -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Nhà cung cấp</label>
              <input v-model="formData.vendor" type="text" placeholder="VD: Công ty điện lực"
                class="w-full px-4 py-3 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500" />
            </div>

            <!-- Notes -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Ghi chú</label>
              <textarea v-model="formData.notes" rows="2" placeholder="Ghi chú thêm..."
                class="w-full px-4 py-3 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 resize-none"></textarea>
            </div>

            <!-- Deduction Summary -->
            <div v-if="formData.amount > 0" class="bg-red-50 rounded-xl p-3 border border-red-200">
              <div class="text-sm text-red-700 font-medium">
                ⚠️ Sẽ trừ <strong>{{ formatPrice(formData.amount) }}</strong> từ quỹ vận hành ({{ formData.money_type === 'cash' ? 'tiền mặt' : 'chuyển khoản' }})
              </div>
            </div>

            <!-- Insufficient balance warning -->
            <div v-if="isInsufficientBalance" class="bg-red-100 rounded-xl p-3 border border-red-300">
              <div class="text-sm text-red-800 font-medium">
                ❌ {{ formData.money_type === 'cash' ? 'Tiền mặt' : 'Chuyển khoản' }} không đủ! Cần {{ formatPrice(formData.amount) }}, hiện có {{ formatPrice(currentTypeBalance || 0) }}
              </div>
            </div>

            <div class="h-4"></div>
          </div>

          <!-- Fixed Footer -->
          <div class="flex-shrink-0 bg-white px-4 py-4 border-t flex gap-3 pb-safe">
            <button @click="showCreateForm = false"
              class="flex-1 bg-gray-200 text-gray-700 py-4 rounded-xl font-medium text-base active:bg-gray-300">
              Hủy
            </button>
            <button @click="saveExpense"
              :disabled="submitting || !formData.description || !formData.category_id || formData.amount <= 0 || isInsufficientBalance"
              class="flex-1 bg-orange-500 text-white py-4 rounded-xl font-medium text-base active:bg-orange-600 disabled:opacity-50 disabled:cursor-not-allowed">
              {{ submitting ? 'Đang lưu...' : 'Thêm chi phí' }}
            </button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'
import { formatDate, formatPrice } from '../utils/formatters'
import { fundExpenseService } from '../services/fundExpenseService'
import { expenseService } from '../services/expense'
import { getAllBalances } from '../services/fund'
import { getJournalEntries } from '../services/journal'

const searchQuery = ref('')
const showCreateForm = ref(false)
const showCategoryModal = ref(false)
const newCategoryName = ref('')
const submitting = ref(false)
const fundLoading = ref(false)
const operatingFundBalance = ref(null)
const operatingCashBalance = ref(null)
const operatingTransferBalance = ref(null)

// Operating expense categories — stored in backend with type="operating"
const categories = ref([])

// Journal entries for event_type=expense (source of truth)
const journalEntries = ref([])

const formData = ref({
  description: '',
  category_id: '',
  amount: 0,
  date: new Date().toISOString().split('T')[0],
  money_type: 'cash',
  vendor: '',
  notes: ''
})

// Get deducted amount from journal entry (CREDIT operating line)
const getEntryAmount = (entry) => {
  const creditLine = (entry.lines || []).find(l => l.direction === 'credit' && l.fund_type === 'operating')
  return creditLine?.total_amount || 0
}

const getEntryMoneyType = (entry) => {
  const creditLine = (entry.lines || []).find(l => l.direction === 'credit' && l.fund_type === 'operating')
  if (!creditLine) return 'cash'
  return creditLine.transfer_amount > 0 ? 'transfer' : 'cash'
}

const getEntryDescription = (entry) => {
  return entry.description?.replace(/^Chi tiêu: /, '') || entry.description || ''
}

const filteredEntries = computed(() => {
  let filtered = journalEntries.value
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(e =>
      e.description?.toLowerCase().includes(query) ||
      e.performed_by_name?.toLowerCase().includes(query)
    )
  }
  return [...filtered].sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp))
})

const totalThisMonth = computed(() => {
  const now = new Date()
  return journalEntries.value
    .filter(e => {
      const d = new Date(e.timestamp)
      return d.getMonth() === now.getMonth() && d.getFullYear() === now.getFullYear()
    })
    .reduce((sum, e) => sum + getEntryAmount(e), 0)
})

const totalAllTime = computed(() => {
  return journalEntries.value.reduce((sum, e) => sum + getEntryAmount(e), 0)
})

const currentTypeBalance = computed(() => {
  if (formData.value.money_type === 'transfer') return operatingTransferBalance.value
  return operatingCashBalance.value
})

const isInsufficientBalance = computed(() => {
  if (!formData.value.amount || currentTypeBalance.value === null) return false
  return formData.value.amount > currentTypeBalance.value
})

const formatCompact = (value) => {
  if (!value) return '0k'
  const thousands = value / 1000
  return thousands % 1 === 0 ? `${thousands}k` : `${thousands.toFixed(1)}k`
}

const getCategoryCount = (categoryId) => {
  // count by matching description patterns since category isn't stored in journal
  return journalEntries.value.filter(e => e.category_id === categoryId).length
}

const fetchFundBalance = async () => {
  try {
    fundLoading.value = true
    const data = await getAllBalances()
    const operating = data?.['operating']
    operatingCashBalance.value = operating?.cash || 0
    operatingTransferBalance.value = operating?.transfer || 0
    operatingFundBalance.value = (operating?.cash || 0) + (operating?.transfer || 0)
  } catch (error) {
    console.error('Error fetching fund balance:', error)
    operatingFundBalance.value = null
    operatingCashBalance.value = null
    operatingTransferBalance.value = null
  } finally {
    fundLoading.value = false
  }
}

const openCreateModal = async () => {
  formData.value = {
    description: '',
    category_id: '',
    amount: 0,
    date: new Date().toISOString().split('T')[0],
    money_type: 'cash',
    vendor: '',
    notes: ''
  }
  showCreateForm.value = true
  await fetchFundBalance()
}

const fetchJournalEntries = async () => {
  try {
    const res = await getJournalEntries({ event_type: 'expense', limit: 200 })
    journalEntries.value = res.entries || res || []
  } catch (error) {
    console.error('Error fetching journal entries:', error)
    journalEntries.value = []
  }
}

const saveExpense = async () => {
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

  try {
    submitting.value = true
    await fundExpenseService.createExpenseFromFund({
      description: formData.value.description,
      category_id: formData.value.category_id,
      amount: formData.value.amount,
      date: formData.value.date + 'T00:00:00Z',
      money_type: formData.value.money_type,
      vendor: formData.value.vendor || undefined,
      notes: formData.value.notes || undefined
    })
    alert('Thêm chi phí vận hành thành công')
    showCreateForm.value = false
    await Promise.all([fetchJournalEntries(), fetchFundBalance()])
  } catch (error) {
    console.error('Error saving expense:', error)
    alert(error.message || 'Có lỗi xảy ra')
  } finally {
    submitting.value = false
  }
}

const fetchCategories = async () => {
  try {
    categories.value = await expenseService.getOperatingCategories() || []
  } catch (error) {
    console.error('Error fetching operating categories:', error)
  }
}

const addCategory = async () => {
  const name = newCategoryName.value.trim()
  if (!name) return
  try {
    const newCat = await expenseService.createOperatingCategory(name)
    categories.value.push(newCat)
    newCategoryName.value = ''
  } catch (error) {
    console.error('Error adding category:', error)
    alert('Có lỗi xảy ra khi thêm danh mục')
  }
}

const deleteCategory = async (categoryId) => {
  const hasExpenses = journalEntries.value.some(e => e.category_id === categoryId)
  if (hasExpenses) {
    alert('Không thể xóa danh mục đã có chi phí!')
    return
  }
  if (confirm('Bạn có chắc muốn xóa danh mục này?')) {
    try {
      await expenseService.deleteCategory(categoryId)
      categories.value = categories.value.filter(c => c.id !== categoryId)
    } catch (error) {
      console.error('Error deleting category:', error)
      alert('Có lỗi xảy ra khi xóa danh mục')
    }
  }
}

const refreshData = async () => {
  await Promise.all([
    fetchJournalEntries(),
    fetchCategories(),
    fetchFundBalance()
  ])
}

const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)

onMounted(async () => {
  await refreshData()
})
</script>

<style scoped>
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

.pb-safe {
  padding-bottom: max(1rem, env(safe-area-inset-bottom));
}

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
