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
          <h1 class="text-xl md:text-2xl font-bold text-gray-800">📊 Phân tích lợi nhuận</h1>
        </div>
        
        <!-- View Mode Toggle -->
        <div class="flex gap-2 mb-3">
          <button 
            @click="viewMode = 'category'"
            :class="viewMode === 'category' ? 'bg-blue-500 text-white' : 'bg-white text-gray-700 border border-gray-300'"
            class="flex-1 px-4 py-2 rounded-lg text-sm font-medium active:scale-95 hover:bg-blue-600 hover:text-white transition-all">
            📁 Theo danh mục
          </button>
          <button 
            @click="viewMode = 'operating'"
            :class="viewMode === 'operating' ? 'bg-blue-500 text-white' : 'bg-white text-gray-700 border border-gray-300'"
            class="flex-1 px-4 py-2 rounded-lg text-sm font-medium active:scale-95 hover:bg-blue-600 hover:text-white transition-all">
            💼 Lợi nhuận vận hành
          </button>
        </div>
        
        <!-- Date Range Picker -->
        <div class="space-y-2">
          <!-- Preset Buttons -->
          <div class="flex gap-2 overflow-x-auto pb-2 md:flex-wrap md:overflow-visible">
            <button 
              v-for="preset in datePresets" 
              :key="preset.value"
              @click="selectDatePreset(preset.value)"
              :class="selectedPreset === preset.value ? 'bg-blue-500 text-white' : 'bg-white text-gray-700 border border-gray-300'"
              class="px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap active:scale-95 hover:bg-blue-600 hover:text-white transition-all flex-shrink-0">
              {{ preset.label }}
            </button>
          </div>
          
          <!-- Custom Date Range -->
          <div class="flex gap-2 items-center flex-col sm:flex-row">
            <input 
              v-model="dateRange.start"
              type="date"
              @change="onDateChange"
              class="w-full sm:flex-1 px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
            />
            <span class="text-gray-500 hidden sm:inline">→</span>
            <input 
              v-model="dateRange.end"
              type="date"
              @change="onDateChange"
              class="w-full sm:flex-1 px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24 md:px-6 lg:px-8 max-w-7xl mx-auto w-full">
      <!-- Loading State with Skeletons -->
      <div v-if="loading">
        <div class="bg-blue-50 rounded-xl p-4 mb-4">
          <div class="h-4 bg-blue-200 rounded w-32 mb-2 animate-pulse"></div>
          <div class="h-4 bg-blue-200 rounded w-48 animate-pulse"></div>
        </div>
        <div class="space-y-4">
          <SkeletonLoader v-for="i in 3" :key="i" type="profit-section" />
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

      <!-- Content Views -->
      <div v-else>
        <!-- Category Profit View -->
        <div v-if="viewMode === 'category'">
          <CategoryProfitView 
            :date-range="dateRange"
            :category-profits="categoryProfits" />
        </div>

        <!-- Operating Profit View -->
        <div v-if="viewMode === 'operating'">
          <OperatingProfitView 
            :operating-profit="operatingProfit" />
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import CategoryProfitView from '../components/CategoryProfitView.vue'
import OperatingProfitView from '../components/OperatingProfitView.vue'
import SkeletonLoader from '../components/SkeletonLoader.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'
import { profitAnalysisService } from '../services/profitAnalysis'

// State
const loading = ref(false)
const error = ref(null)
const viewMode = ref('category')
const selectedPreset = ref('today')
const dateRange = ref({
  start: '',
  end: ''
})
const categoryProfits = ref([])
const operatingProfit = ref(null)

// Date presets
const datePresets = [
  { label: 'Hôm nay', value: 'today' },
  { label: 'Tuần này', value: 'this_week' },
  { label: 'Tháng này', value: 'this_month' }
]

// Functions
const getDateRangeForPreset = (preset) => {
  const today = new Date()
  const year = today.getFullYear()
  const month = String(today.getMonth() + 1).padStart(2, '0')
  const day = String(today.getDate()).padStart(2, '0')
  
  switch (preset) {
    case 'today':
      return {
        start: `${year}-${month}-${day}`,
        end: `${year}-${month}-${day}`
      }
    case 'this_week': {
      const dayOfWeek = today.getDay()
      const monday = new Date(today)
      monday.setDate(today.getDate() - (dayOfWeek === 0 ? 6 : dayOfWeek - 1))
      const sunday = new Date(monday)
      sunday.setDate(monday.getDate() + 6)
      
      return {
        start: monday.toISOString().split('T')[0],
        end: sunday.toISOString().split('T')[0]
      }
    }
    case 'this_month': {
      const firstDay = new Date(year, today.getMonth(), 1)
      const lastDay = new Date(year, today.getMonth() + 1, 0)
      
      return {
        start: firstDay.toISOString().split('T')[0],
        end: lastDay.toISOString().split('T')[0]
      }
    }
    default:
      return {
        start: `${year}-${month}-${day}`,
        end: `${year}-${month}-${day}`
      }
  }
}

const selectDatePreset = (preset) => {
  selectedPreset.value = preset
  dateRange.value = getDateRangeForPreset(preset)
  fetchData()
}

const onDateChange = () => {
  selectedPreset.value = null
  fetchData()
}

// Fetch data function
const fetchData = async () => {
  if (!dateRange.value.start || !dateRange.value.end) {
    return
  }
  
  loading.value = true
  error.value = null
  
  try {
    if (viewMode.value === 'category') {
      const response = await profitAnalysisService.getCategoryProfit(dateRange.value)
      categoryProfits.value = response.categories || []
    } else if (viewMode.value === 'operating') {
      const response = await profitAnalysisService.getOperatingProfit(dateRange.value)
      operatingProfit.value = response
    }
  } catch (err) {
    console.error('Error fetching profit analysis:', err)
    error.value = err.response?.data?.error || 'Không thể tải dữ liệu phân tích lợi nhuận'
  } finally {
    loading.value = false
  }
}

// Pull to refresh
const { pullDistance, isRefreshing } = usePullToRefresh(fetchData)

// Watch view mode changes
watch(viewMode, () => {
  fetchData()
})

// Lifecycle
onMounted(() => {
  // Initialize with today's date
  selectDatePreset('today')
})
</script>

<style scoped>
.active\:scale-95:active {
  transform: scale(0.95);
}
</style>
