<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Pull to Refresh Indicator -->
    <PullToRefresh 
      :pull-distance="pullDistance" 
      :is-refreshing="isRefreshing"
      :threshold="80" />
    
    <!-- Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3 md:px-6 lg:px-8" style="padding-top: max(0.75rem, env(safe-area-inset-top))">
        <div class="flex items-center justify-between mb-3">
          <h1 class="text-xl md:text-2xl font-bold text-gray-800">💰 Chi phí món</h1>
          <!-- Desktop: View Toggle -->
          <div class="hidden md:flex gap-2">
            <button @click="viewMode = 'card'" 
              :class="viewMode === 'card' ? 'bg-blue-500 text-white' : 'bg-white text-gray-700 border border-gray-300'"
              class="px-3 py-1 rounded-lg text-sm font-medium">
              📱 Thẻ
            </button>
            <button @click="viewMode = 'table'" 
              :class="viewMode === 'table' ? 'bg-blue-500 text-white' : 'bg-white text-gray-700 border border-gray-300'"
              class="px-3 py-1 rounded-lg text-sm font-medium">
              📊 Bảng
            </button>
          </div>
        </div>
        
        <!-- Search Bar -->
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Tìm kiếm món..."
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm md:text-base"
        />
        
        <!-- Filters - Responsive -->
        <div class="mt-2 flex gap-2 overflow-x-auto pb-2 md:flex-wrap md:overflow-visible">
          <!-- Category Filter -->
          <button @click="categoryFilter = ''" 
            :class="categoryFilter === '' ? 'bg-blue-500 text-white' : 'bg-white text-gray-700 border border-gray-300'"
            class="px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap flex-shrink-0">
            📁 Tất cả
          </button>
          <button v-for="category in uniqueCategories" :key="category"
            @click="categoryFilter = category" 
            :class="categoryFilter === category ? 'bg-blue-500 text-white' : 'bg-white text-gray-700 border border-gray-300'"
            class="px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap flex-shrink-0">
            {{ category }}
          </button>
        </div>
        
        <!-- Sort Controls - Responsive -->
        <div class="mt-2 flex gap-2 items-center">
          <span class="text-xs text-gray-600 font-medium hidden sm:inline">Sắp xếp:</span>
          <select v-model="sortBy" 
            class="flex-1 md:flex-none md:min-w-[180px] px-3 py-1 text-xs border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500">
            <option value="profit_margin">Lợi nhuận %</option>
            <option value="absolute_profit">Lợi nhuận tiền</option>
            <option value="name">Tên món</option>
          </select>
          <button @click="toggleSortOrder"
            class="px-3 py-1 bg-white border border-gray-300 rounded-lg text-xs font-medium active:bg-gray-100 flex-shrink-0">
            {{ sortOrder === 'desc' ? '↓ Giảm' : '↑ Tăng' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24 md:px-6 lg:px-8 max-w-7xl mx-auto w-full">
      <!-- Loading State with Skeletons -->
      <div v-if="loading">
        <!-- Summary Skeleton -->
        <SkeletonLoader type="summary" />
        
        <!-- Card Skeletons -->
        <div v-if="viewMode === 'card'" class="space-y-3">
          <SkeletonLoader v-for="i in 5" :key="i" type="card" />
        </div>
        
        <!-- Table Skeleton -->
        <div v-else class="bg-white rounded-xl shadow-sm overflow-hidden">
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead class="bg-gray-50 border-b">
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tên món</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Danh mục</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Giá bán</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Chi phí</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">LN %</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">LN tiền</th>
                  <th class="px-4 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">Trạng thái</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <SkeletonLoader v-for="i in 8" :key="i" type="table-row" />
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-16">
        <div class="text-6xl mb-4">❌</div>
        <p class="text-red-500 mb-2 font-medium">Lỗi tải dữ liệu</p>
        <p class="text-gray-600 mb-4 text-sm">{{ error }}</p>
        <button @click="fetchData" 
          class="bg-blue-500 text-white px-6 py-2 rounded-lg font-medium active:bg-blue-600 hover:bg-blue-600 transition-colors">
          🔄 Thử lại
        </button>
      </div>

      <!-- Content -->
      <div v-else>
        <!-- Summary Statistics -->
        <div v-if="summary" class="bg-gradient-to-br from-blue-500 to-purple-500 rounded-xl p-4 md:p-6 mb-4 text-white shadow-lg">
          <div class="flex items-center justify-between mb-3">
            <div class="text-sm md:text-base opacity-90">Tổng quan chi phí món</div>
            <div v-if="recalculationStatus?.in_progress" class="text-xs opacity-90 bg-white/20 px-2 py-1 rounded-full flex items-center gap-1">
              <span class="animate-spin">⏳</span>
              <span>Đang cập nhật...</span>
            </div>
          </div>
          <div class="grid grid-cols-4 gap-3 md:gap-6">
            <div class="text-center">
              <div class="text-lg md:text-2xl font-bold leading-tight">{{ summary.total_items || 0 }}</div>
              <div class="text-xs md:text-sm opacity-90 whitespace-nowrap mt-1">Tổng món</div>
            </div>
            <div class="text-center">
              <div class="text-lg md:text-2xl font-bold leading-tight text-red-200">{{ summary.loss_count || 0 }}</div>
              <div class="text-xs md:text-sm opacity-90 whitespace-nowrap mt-1">Bán lỗ</div>
            </div>
            <div class="text-center">
              <div class="text-lg md:text-2xl font-bold leading-tight text-yellow-200">{{ summary.low_margin_count || 0 }}</div>
              <div class="text-xs md:text-sm opacity-90 whitespace-nowrap mt-1">LN thấp</div>
            </div>
            <div class="text-center">
              <div class="text-lg md:text-2xl font-bold leading-tight">{{ formatPercentage(summary.average_profit_margin) }}</div>
              <div class="text-xs md:text-sm opacity-90 whitespace-nowrap mt-1">LN TB</div>
            </div>
          </div>
          
          <!-- Recalculation Status Details -->
          <div v-if="recalculationStatus && recalculationStatus.in_progress" class="mt-3 pt-3 border-t border-white/20">
            <div class="text-xs md:text-sm opacity-90">
              Đã xử lý: {{ recalculationStatus.processed_items || 0 }} / {{ recalculationStatus.queued_items || 0 }} món
            </div>
          </div>
        </div>
        
        <!-- Empty State -->
        <div v-if="filteredMenuItems.length === 0" class="text-center py-16">
          <div class="text-6xl mb-4">📭</div>
          <p class="text-gray-700 font-medium mb-2">Không tìm thấy món nào</p>
          <p class="text-gray-500 text-sm">
            {{ searchQuery || categoryFilter ? 'Thử thay đổi bộ lọc hoặc tìm kiếm' : 'Chưa có món nào trong hệ thống' }}
          </p>
        </div>

        <!-- Card View (Mobile & Optional Desktop) -->
        <div v-else-if="viewMode === 'card'" class="space-y-3">
          <div v-for="item in filteredMenuItems" :key="item.menu_item_id"
            @click="openCostBreakdown(item)"
            class="bg-white rounded-2xl p-4 shadow-sm cursor-pointer active:scale-98 hover:shadow-md transition-all"
            :class="getItemBorderClass(item.warning_status)">
            
            <!-- Item Header -->
            <div class="flex justify-between items-start mb-3">
              <div class="flex-1">
                <div class="flex items-center gap-2 mb-1 flex-wrap">
                  <h3 class="font-bold text-lg">{{ item.name }}</h3>
                  <span class="px-2 py-0.5 rounded text-[10px] font-medium bg-gray-100 text-gray-600">
                    {{ item.category }}
                  </span>
                </div>
                <div class="flex items-center gap-2 flex-wrap">
                  <span class="text-sm text-gray-600">Giá bán: {{ formatPrice(item.price) }}</span>
                  <span v-if="item.cost_status" 
                    :class="getCostStatusClass(item.cost_status)"
                    class="px-2 py-0.5 rounded text-[10px] font-medium">
                    {{ getCostStatusLabel(item.cost_status) }}
                  </span>
                </div>
              </div>
              <div class="text-right ml-2">
                <div class="text-sm text-gray-500 mb-1">Chi phí</div>
                <div class="text-lg font-bold text-blue-600">{{ formatPrice(item.current_cost) }}</div>
              </div>
            </div>

            <!-- Profit Metrics -->
            <div class="grid grid-cols-2 gap-3 pt-3 border-t">
              <div>
                <div class="text-xs text-gray-500 mb-1">Lợi nhuận %</div>
                <div class="font-bold" :class="getProfitMarginColor(item.profit_margin, item.warning_status)">
                  {{ formatPercentage(item.profit_margin) }}
                </div>
              </div>
              <div>
                <div class="text-xs text-gray-500 mb-1">Lợi nhuận tiền</div>
                <div class="font-bold" :class="getAbsoluteProfitColor(item.absolute_profit)">
                  {{ formatPrice(item.absolute_profit) }}
                </div>
              </div>
            </div>

            <!-- Warning Indicator -->
            <div v-if="item.warning_status !== 'none'" class="mt-3 pt-3 border-t">
              <div class="flex items-center gap-2 text-sm" :class="getWarningTextClass(item.warning_status)">
                <span>{{ getWarningIcon(item.warning_status) }}</span>
                <span class="font-medium">{{ getWarningMessage(item.warning_status) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Table View (Desktop) -->
        <div v-else-if="viewMode === 'table'" class="bg-white rounded-xl shadow-sm overflow-hidden">
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead class="bg-gray-50 border-b">
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tên món</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Danh mục</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Giá bán</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Chi phí</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">LN %</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">LN tiền</th>
                  <th class="px-4 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">Trạng thái</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <tr v-for="item in filteredMenuItems" :key="item.menu_item_id"
                  @click="openCostBreakdown(item)"
                  class="hover:bg-gray-50 cursor-pointer transition-colors"
                  :class="getTableRowClass(item.warning_status)">
                  <td class="px-4 py-3 text-sm font-medium text-gray-900">{{ item.name }}</td>
                  <td class="px-4 py-3 text-sm text-gray-600">
                    <span class="px-2 py-1 rounded-full text-xs bg-gray-100">{{ item.category }}</span>
                  </td>
                  <td class="px-4 py-3 text-sm text-right text-gray-900">{{ formatPrice(item.price) }}</td>
                  <td class="px-4 py-3 text-sm text-right font-medium text-blue-600">{{ formatPrice(item.current_cost) }}</td>
                  <td class="px-4 py-3 text-sm text-right font-bold" :class="getProfitMarginColor(item.profit_margin, item.warning_status)">
                    {{ formatPercentage(item.profit_margin) }}
                  </td>
                  <td class="px-4 py-3 text-sm text-right font-bold" :class="getAbsoluteProfitColor(item.absolute_profit)">
                    {{ formatPrice(item.absolute_profit) }}
                  </td>
                  <td class="px-4 py-3 text-center">
                    <span v-if="item.cost_status" 
                      :class="getCostStatusClass(item.cost_status)"
                      class="px-2 py-1 rounded text-xs font-medium inline-block">
                      {{ getCostStatusLabel(item.cost_status) }}
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />
    
    <!-- Cost Breakdown Modal Component -->
    <MenuItemCostBreakdown 
      :is-open="showCostBreakdown"
      :menu-item-id="selectedMenuItemId"
      @close="closeCostBreakdown" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import MenuItemCostBreakdown from '../components/MenuItemCostBreakdown.vue'
import SkeletonLoader from '../components/SkeletonLoader.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'
import { menuCostService } from '../services/menuCost'
import { formatPrice, formatPercentage } from '../utils/formatters'

// State
const loading = ref(false)
const error = ref(null)
const searchQuery = ref('')
const categoryFilter = ref('')
const sortBy = ref('profit_margin')
const sortOrder = ref('desc')
const menuItems = ref([])
const summary = ref(null)
const recalculationStatus = ref(null)
const viewMode = ref('card') // 'card' or 'table'

// Cost breakdown modal state
const showCostBreakdown = ref(false)
const selectedMenuItemId = ref(null)

// Computed
const uniqueCategories = computed(() => {
  const categories = menuItems.value
    .map(item => item.category)
    .filter((value, index, self) => value && self.indexOf(value) === index)
  return categories.sort()
})

const filteredMenuItems = computed(() => {
  let filtered = menuItems.value
  
  // Filter by category
  if (categoryFilter.value) {
    filtered = filtered.filter(item => item.category === categoryFilter.value)
  }
  
  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(item => 
      item.name?.toLowerCase().includes(query) ||
      item.category?.toLowerCase().includes(query)
    )
  }
  
  // Sort
  const sorted = [...filtered].sort((a, b) => {
    let compareValue = 0
    
    if (sortBy.value === 'profit_margin') {
      compareValue = (a.profit_margin || 0) - (b.profit_margin || 0)
    } else if (sortBy.value === 'absolute_profit') {
      compareValue = (a.absolute_profit || 0) - (b.absolute_profit || 0)
    } else if (sortBy.value === 'name') {
      compareValue = (a.name || '').localeCompare(b.name || '')
    }
    
    return sortOrder.value === 'desc' ? -compareValue : compareValue
  })
  
  return sorted
})

// Functions
const toggleSortOrder = () => {
  sortOrder.value = sortOrder.value === 'desc' ? 'asc' : 'desc'
}

// Helper functions
const getCostStatusLabel = (status) => {
  const labels = {
    'FINAL': '✓ Chính thức',
    'ESTIMATED': '~ Ước tính',
    'INCOMPLETE': '⚠ Thiếu dữ liệu'
  }
  return labels[status] || status
}

const getCostStatusClass = (status) => {
  const classes = {
    'FINAL': 'bg-green-100 text-green-700',
    'ESTIMATED': 'bg-yellow-100 text-yellow-700',
    'INCOMPLETE': 'bg-red-100 text-red-700'
  }
  return classes[status] || 'bg-gray-100 text-gray-700'
}

const getItemBorderClass = (warningStatus) => {
  const classes = {
    'loss': 'border-l-4 border-red-500',
    'low_margin': 'border-l-4 border-yellow-500',
    'none': 'border-l-4 border-green-500'
  }
  return classes[warningStatus] || ''
}

const getProfitMarginColor = (margin, warningStatus) => {
  if (warningStatus === 'loss') return 'text-red-600'
  if (warningStatus === 'low_margin') return 'text-yellow-600'
  return 'text-green-600'
}

const getAbsoluteProfitColor = (profit) => {
  if (profit < 0) return 'text-red-600'
  if (profit === 0) return 'text-gray-600'
  return 'text-green-600'
}

const getWarningIcon = (warningStatus) => {
  const icons = {
    'loss': '🔴',
    'low_margin': '⚠️'
  }
  return icons[warningStatus] || ''
}

const getWarningMessage = (warningStatus) => {
  const messages = {
    'loss': 'Bán lỗ - Chi phí cao hơn giá bán',
    'low_margin': 'Lợi nhuận thấp - Cần xem xét điều chỉnh'
  }
  return messages[warningStatus] || ''
}

const getWarningTextClass = (warningStatus) => {
  const classes = {
    'loss': 'text-red-600',
    'low_margin': 'text-yellow-600'
  }
  return classes[warningStatus] || ''
}

const getTableRowClass = (warningStatus) => {
  const classes = {
    'loss': 'border-l-4 border-red-500',
    'low_margin': 'border-l-4 border-yellow-500',
    'none': 'border-l-4 border-green-500'
  }
  return classes[warningStatus] || ''
}

const openCostBreakdown = (item) => {
  selectedMenuItemId.value = item.menu_item_id
  showCostBreakdown.value = true
}

const closeCostBreakdown = () => {
  showCostBreakdown.value = false
  selectedMenuItemId.value = null
}

// Fetch data function
const fetchData = async () => {
  loading.value = true
  error.value = null
  
  try {
    const response = await menuCostService.getMenuCosts({})
    menuItems.value = response.items || []
    summary.value = response.summary || null
    recalculationStatus.value = response.recalculation_status || null
  } catch (err) {
    console.error('Error fetching menu costs:', err)
    error.value = err.response?.data?.error || 'Không thể tải dữ liệu chi phí món'
  } finally {
    loading.value = false
  }
}

// Pull to refresh
const { pullDistance, isRefreshing } = usePullToRefresh(fetchData)

// Lifecycle
onMounted(async () => {
  await fetchData()
})
</script>

<style scoped>
.active\:scale-95:active {
  transform: scale(0.95);
}

.active\:scale-98:active {
  transform: scale(0.98);
}

/* Responsive table */
@media (max-width: 768px) {
  table {
    font-size: 0.875rem;
  }
}

@keyframes slide-up {
  from {
    transform: translateY(100%);
  }
  to {
    transform: translateY(0);
  }
}

.animate-slide-up {
  animation: slide-up 0.3s ease-out;
}
</style>
