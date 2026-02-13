<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Cost Breakdown Modal -->
    <CostBreakdownModal 
      :is-open="showCostBreakdownModal"
      :menu-item-id="selectedMenuItemId"
      :default-variant-id="selectedDefaultVariantId"
      @close="closeCostBreakdownModal" />
    
    <!-- Profit Comparison Modal -->
    <ProfitComparisonModal
      :is-open="showProfitComparisonModal"
      :menu-item-id="selectedMenuItemId"
      @close="closeProfitComparisonModal" />
    
    <!-- Pull to Refresh Indicator -->
    <PullToRefresh 
      :pull-distance="pullDistance" 
      :is-refreshing="isRefreshing"
      :threshold="80" />
    
    <!-- Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3 md:px-6 lg:px-8" style="padding-top: max(0.75rem, env(safe-area-inset-top))">
        <div class="flex items-center justify-between mb-3">
          <h1 class="text-xl md:text-2xl font-bold text-gray-800">📊 Phân tích chi phí theo size</h1>
        </div>
        
        <!-- Search Bar -->
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Tìm kiếm món..."
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm md:text-base mb-2"
        />
        
        <!-- Filters -->
        <div class="flex gap-2 overflow-x-auto pb-2 md:flex-wrap md:overflow-visible">
          <!-- Cost Status Filter -->
          <button @click="costStatusFilter = ''" 
            :class="costStatusFilter === '' ? 'bg-blue-500 text-white' : 'bg-white text-gray-700 border border-gray-300'"
            class="px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap flex-shrink-0">
            Tất cả
          </button>
          <button @click="costStatusFilter = 'FINAL'" 
            :class="costStatusFilter === 'FINAL' ? 'bg-green-500 text-white' : 'bg-white text-gray-700 border border-gray-300'"
            class="px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap flex-shrink-0">
            ✓ Chính thức
          </button>
          <button @click="costStatusFilter = 'ESTIMATED'" 
            :class="costStatusFilter === 'ESTIMATED' ? 'bg-yellow-500 text-white' : 'bg-white text-gray-700 border border-gray-300'"
            class="px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap flex-shrink-0">
            ~ Ước tính
          </button>
          <button @click="costStatusFilter = 'INCOMPLETE'" 
            :class="costStatusFilter === 'INCOMPLETE' ? 'bg-red-500 text-white' : 'bg-white text-gray-700 border border-gray-300'"
            class="px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap flex-shrink-0">
            ⚠ Thiếu dữ liệu
          </button>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24 md:px-6 lg:px-8 max-w-7xl mx-auto w-full">
      <!-- Loading State -->
      <div v-if="loading" class="space-y-3">
        <SkeletonLoader v-for="i in 5" :key="i" type="card" />
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
        <!-- Empty State -->
        <div v-if="filteredMenuItems.length === 0" class="text-center py-16">
          <div class="text-6xl mb-4">📭</div>
          <p class="text-gray-700 font-medium mb-2">Không tìm thấy món nào</p>
          <p class="text-gray-500 text-sm">
            {{ searchQuery || costStatusFilter ? 'Thử thay đổi bộ lọc hoặc tìm kiếm' : 'Chưa có món nào trong hệ thống' }}
          </p>
        </div>

        <!-- Menu Items List -->
        <div v-else class="space-y-4">
          <div v-for="item in filteredMenuItems" :key="item.id"
            class="bg-white rounded-2xl p-4 shadow-sm">
            
            <!-- Item Header -->
            <div class="flex justify-between items-start mb-3 pb-3 border-b">
              <div class="flex-1">
                <h3 class="font-bold text-lg mb-1">{{ item.name }}</h3>
                <span class="px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-600">
                  {{ item.category }}
                </span>
              </div>
            </div>

            <!-- Single-Size Item -->
            <div v-if="!item.has_variants" class="space-y-2">
              <div class="flex items-center justify-between text-sm">
                <span class="text-gray-600">Loại:</span>
                <span class="font-medium">Một size</span>
              </div>
              
              <div class="grid grid-cols-2 gap-3 p-3 bg-gray-50 rounded-lg">
                <div>
                  <div class="text-xs text-gray-500 mb-1">Giá bán</div>
                  <div class="font-bold text-blue-600">{{ formatPrice(item.price) }}</div>
                </div>
                <div>
                  <div class="text-xs text-gray-500 mb-1">Chi phí</div>
                  <div class="font-bold text-orange-600">{{ formatPrice(item.current_cost) }}</div>
                </div>
                <div>
                  <div class="text-xs text-gray-500 mb-1">Lợi nhuận</div>
                  <div class="font-bold" :class="getProfitColor(item.price - item.current_cost)">
                    {{ formatPrice(item.price - item.current_cost) }}
                  </div>
                </div>
                <div>
                  <div class="text-xs text-gray-500 mb-1">Tỷ suất LN</div>
                  <div class="font-bold" :class="getProfitMarginColor(calculateProfitMargin(item.price, item.current_cost))">
                    {{ formatPercentage(calculateProfitMargin(item.price, item.current_cost)) }}
                  </div>
                </div>
              </div>

              <div class="flex items-center justify-between text-sm pt-2">
                <span 
                  :class="getCostStatusClass(item.cost_status)"
                  class="px-2 py-1 rounded text-xs font-medium">
                  {{ getCostStatusLabel(item.cost_status) }}
                </span>
                <span v-if="item.cost_last_calculated_at" class="text-xs text-gray-500">
                  Cập nhật: {{ formatDate(item.cost_last_calculated_at) }}
                </span>
              </div>

              <!-- View Cost Breakdown Button -->
              <button @click="openCostBreakdownModal(item.id, null)"
                class="w-full mt-3 px-4 py-2 bg-blue-500 text-white rounded-lg text-sm font-medium active:bg-blue-600 hover:bg-blue-600 transition-colors">
                📊 Xem chi tiết chi phí
              </button>
            </div>

            <!-- Multi-Size Item with Variants -->
            <div v-else class="space-y-3">
              <div class="flex items-center justify-between text-sm mb-2">
                <span class="text-gray-600">Loại:</span>
                <span class="font-medium text-purple-600">Nhiều size ({{ item.variants.length }} variants)</span>
              </div>

              <!-- Variants List -->
              <div class="space-y-2">
                <div v-for="variant in item.variants" :key="variant.id"
                  class="p-3 rounded-lg border-2"
                  :class="variant.is_default ? 'border-purple-300 bg-purple-50' : 'border-gray-200 bg-gray-50'">
                  
                  <!-- Variant Header -->
                  <div class="flex items-center justify-between mb-2">
                    <div class="flex items-center gap-2">
                      <span class="font-bold text-base">{{ variant.name }}</span>
                      <span v-if="variant.is_default" 
                        class="px-2 py-0.5 rounded text-[10px] font-medium bg-purple-500 text-white">
                        Mặc định
                      </span>
                    </div>
                    <span 
                      :class="getCostStatusClass(variant.cost_status)"
                      class="px-2 py-0.5 rounded text-[10px] font-medium">
                      {{ getCostStatusLabel(variant.cost_status) }}
                    </span>
                  </div>

                  <!-- Variant Metrics -->
                  <div class="grid grid-cols-2 gap-2 text-sm">
                    <div>
                      <div class="text-xs text-gray-500">Giá bán</div>
                      <div class="font-bold text-blue-600">{{ formatPrice(variant.price) }}</div>
                    </div>
                    <div>
                      <div class="text-xs text-gray-500">Chi phí</div>
                      <div class="font-bold text-orange-600">{{ formatPrice(variant.current_cost) }}</div>
                    </div>
                    <div>
                      <div class="text-xs text-gray-500">Lợi nhuận</div>
                      <div class="font-bold" :class="getProfitColor(variant.price - variant.current_cost)">
                        {{ formatPrice(variant.price - variant.current_cost) }}
                      </div>
                    </div>
                    <div>
                      <div class="text-xs text-gray-500">Tỷ suất LN</div>
                      <div class="font-bold" :class="getProfitMarginColor(calculateProfitMargin(variant.price, variant.current_cost))">
                        {{ formatPercentage(calculateProfitMargin(variant.price, variant.current_cost)) }}
                      </div>
                    </div>
                  </div>

                  <!-- Variant Last Updated -->
                  <div v-if="variant.cost_last_calculated_at" class="text-xs text-gray-500 mt-2 text-right">
                    Cập nhật: {{ formatDate(variant.cost_last_calculated_at) }}
                  </div>

                  <!-- View Cost Breakdown Button -->
                  <button @click="openCostBreakdownModal(item.id, getDefaultVariantId(item))"
                    class="w-full mt-3 px-3 py-2 bg-blue-500 text-white rounded-lg text-xs font-medium active:bg-blue-600 hover:bg-blue-600 transition-colors">
                    📊 Xem chi tiết chi phí
                  </button>
                </div>
              </div>

              <!-- Compare Variants Button -->
              <button @click="openProfitComparisonModal(item.id)"
                class="w-full mt-3 px-4 py-2 bg-gradient-to-r from-purple-500 to-blue-500 text-white rounded-lg text-sm font-medium active:opacity-90 hover:opacity-90 transition-opacity">
                📊 So sánh lợi nhuận các size
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import SkeletonLoader from '../components/SkeletonLoader.vue'
import CostBreakdownModal from '../components/CostBreakdownModal.vue'
import ProfitComparisonModal from '../components/ProfitComparisonModal.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'
import { menuService } from '../services/menu'
import { formatPrice, formatPercentage } from '../utils/formatters'

// State
const loading = ref(false)
const error = ref(null)
const searchQuery = ref('')
const costStatusFilter = ref('')
const menuItems = ref([])
const showCostBreakdownModal = ref(false)
const showProfitComparisonModal = ref(false)
const selectedMenuItemId = ref(null)
const selectedDefaultVariantId = ref(null)

// Computed
const filteredMenuItems = computed(() => {
  let filtered = menuItems.value
  
  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(item => 
      item.name?.toLowerCase().includes(query) ||
      item.category?.toLowerCase().includes(query)
    )
  }
  
  // Filter by cost status
  if (costStatusFilter.value) {
    filtered = filtered.filter(item => {
      if (!item.has_variants) {
        // Single-size item
        return item.cost_status === costStatusFilter.value
      } else {
        // Multi-size item - check if any variant matches
        return item.variants?.some(v => v.cost_status === costStatusFilter.value)
      }
    })
  }
  
  return filtered
})

// Helper functions
const getCostStatusLabel = (status) => {
  const labels = {
    'FINAL': '✓ Chính thức',
    'ESTIMATED': '~ Ước tính',
    'INCOMPLETE': '⚠ Thiếu dữ liệu'
  }
  return labels[status] || status || 'Chưa tính'
}

const getCostStatusClass = (status) => {
  const classes = {
    'FINAL': 'bg-green-100 text-green-700',
    'ESTIMATED': 'bg-yellow-100 text-yellow-700',
    'INCOMPLETE': 'bg-red-100 text-red-700'
  }
  return classes[status] || 'bg-gray-100 text-gray-700'
}

const calculateProfitMargin = (price, cost) => {
  if (!price || price === 0) return 0
  return ((price - cost) / price) * 100
}

const getProfitColor = (profit) => {
  if (profit < 0) return 'text-red-600'
  if (profit === 0) return 'text-gray-600'
  return 'text-green-600'
}

const getProfitMarginColor = (margin) => {
  if (margin < 0) return 'text-red-600'
  if (margin < 20) return 'text-yellow-600'
  return 'text-green-600'
}

const formatDate = (dateString) => {
  if (!dateString) return 'Chưa có'
  const date = new Date(dateString)
  return date.toLocaleString('vi-VN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getDefaultVariantId = (item) => {
  if (!item.has_variants || !item.variants) return null
  const defaultVariant = item.variants.find(v => v.is_default)
  return defaultVariant ? defaultVariant.id : null
}

const openCostBreakdownModal = (menuItemId, defaultVariantId) => {
  selectedMenuItemId.value = menuItemId
  selectedDefaultVariantId.value = defaultVariantId
  showCostBreakdownModal.value = true
}

const closeCostBreakdownModal = () => {
  showCostBreakdownModal.value = false
  selectedMenuItemId.value = null
  selectedDefaultVariantId.value = null
}

const openProfitComparisonModal = (menuItemId) => {
  selectedMenuItemId.value = menuItemId
  showProfitComparisonModal.value = true
}

const closeProfitComparisonModal = () => {
  showProfitComparisonModal.value = false
  selectedMenuItemId.value = null
}

// Fetch data function
const fetchData = async () => {
  loading.value = true
  error.value = null
  
  try {
    const response = await menuService.getMenuItems()
    menuItems.value = response.data || []
  } catch (err) {
    console.error('Error fetching menu items:', err)
    error.value = err.response?.data?.error || 'Không thể tải dữ liệu menu'
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

.active\:bg-blue-600:active {
  background-color: #2563eb;
}
</style>
