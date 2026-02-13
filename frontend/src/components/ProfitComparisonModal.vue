<template>
  <div v-if="isOpen" 
    class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 p-4"
    @click.self="$emit('close')">
    <div class="bg-white rounded-2xl shadow-2xl w-full max-w-4xl max-h-[90vh] overflow-hidden flex flex-col">
      <!-- Header -->
      <div class="px-6 py-4 border-b bg-gradient-to-r from-purple-500 to-blue-500 text-white flex-shrink-0">
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-bold">📊 So sánh lợi nhuận theo size</h2>
          <button @click="$emit('close')" 
            class="text-white hover:bg-white hover:bg-opacity-20 rounded-full p-2 transition-colors">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <p v-if="menuItem" class="text-sm mt-1 text-white text-opacity-90">{{ menuItem.name }}</p>
      </div>

      <!-- Content -->
      <div class="flex-1 overflow-y-auto p-6">
        <!-- Loading State -->
        <div v-if="loading" class="flex items-center justify-center py-12">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="text-center py-12">
          <div class="text-6xl mb-4">❌</div>
          <p class="text-red-500 font-medium">{{ error }}</p>
        </div>

        <!-- No Variants -->
        <div v-else-if="!menuItem || !menuItem.has_variants || !menuItem.variants || menuItem.variants.length === 0" 
          class="text-center py-12">
          <div class="text-6xl mb-4">📭</div>
          <p class="text-gray-600 font-medium">Món này không có nhiều size</p>
        </div>

        <!-- Profit Comparison -->
        <div v-else class="space-y-6">
          <!-- Summary Stats -->
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div class="bg-blue-50 rounded-lg p-4">
              <div class="text-xs text-blue-600 font-medium mb-1">Tổng số size</div>
              <div class="text-2xl font-bold text-blue-700">{{ menuItem.variants.length }}</div>
            </div>
            <div class="bg-green-50 rounded-lg p-4">
              <div class="text-xs text-green-600 font-medium mb-1">Size lời nhất</div>
              <div class="text-2xl font-bold text-green-700">{{ mostProfitableVariant?.name || 'N/A' }}</div>
            </div>
            <div class="bg-purple-50 rounded-lg p-4">
              <div class="text-xs text-purple-600 font-medium mb-1">Lợi nhuận cao nhất</div>
              <div class="text-2xl font-bold text-purple-700">{{ formatPrice(maxProfit) }}</div>
            </div>
            <div class="bg-orange-50 rounded-lg p-4">
              <div class="text-xs text-orange-600 font-medium mb-1">Tỷ suất LN cao nhất</div>
              <div class="text-2xl font-bold text-orange-700">{{ formatPercentage(maxProfitMargin) }}</div>
            </div>
          </div>

          <!-- Variants Comparison Table -->
          <div class="overflow-x-auto">
            <table class="w-full border-collapse">
              <thead>
                <tr class="bg-gray-100">
                  <th class="px-4 py-3 text-left text-sm font-bold text-gray-700 border-b-2 border-gray-300">Size</th>
                  <th class="px-4 py-3 text-right text-sm font-bold text-gray-700 border-b-2 border-gray-300">Giá bán</th>
                  <th class="px-4 py-3 text-right text-sm font-bold text-gray-700 border-b-2 border-gray-300">Chi phí</th>
                  <th class="px-4 py-3 text-right text-sm font-bold text-gray-700 border-b-2 border-gray-300">Lợi nhuận</th>
                  <th class="px-4 py-3 text-right text-sm font-bold text-gray-700 border-b-2 border-gray-300">Tỷ suất LN</th>
                  <th class="px-4 py-3 text-center text-sm font-bold text-gray-700 border-b-2 border-gray-300">Trạng thái</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="variant in sortedVariants" :key="variant.id"
                  :class="[
                    'border-b border-gray-200 hover:bg-gray-50 transition-colors',
                    variant.id === mostProfitableVariant?.id ? 'bg-green-50' : ''
                  ]">
                  <td class="px-4 py-3">
                    <div class="flex items-center gap-2">
                      <span class="font-bold text-gray-800">{{ variant.name }}</span>
                      <span v-if="variant.is_default" 
                        class="px-2 py-0.5 rounded text-[10px] font-medium bg-purple-500 text-white">
                        Mặc định
                      </span>
                      <span v-if="variant.id === mostProfitableVariant?.id"
                        class="px-2 py-0.5 rounded text-[10px] font-medium bg-green-500 text-white">
                        🏆 Lời nhất
                      </span>
                    </div>
                  </td>
                  <td class="px-4 py-3 text-right">
                    <span class="font-bold text-blue-600">{{ formatPrice(variant.price) }}</span>
                  </td>
                  <td class="px-4 py-3 text-right">
                    <span class="font-bold text-orange-600">{{ formatPrice(variant.current_cost) }}</span>
                  </td>
                  <td class="px-4 py-3 text-right">
                    <span class="font-bold" :class="getProfitColor(calculateProfit(variant))">
                      {{ formatPrice(calculateProfit(variant)) }}
                    </span>
                  </td>
                  <td class="px-4 py-3 text-right">
                    <span class="font-bold" :class="getProfitMarginColor(calculateProfitMargin(variant))">
                      {{ formatPercentage(calculateProfitMargin(variant)) }}
                    </span>
                  </td>
                  <td class="px-4 py-3 text-center">
                    <span :class="getCostStatusClass(variant.cost_status)"
                      class="px-2 py-1 rounded text-xs font-medium inline-block">
                      {{ getCostStatusLabel(variant.cost_status) }}
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Cost Difference Analysis -->
          <div class="bg-gradient-to-br from-blue-50 to-purple-50 rounded-xl p-6">
            <h3 class="text-lg font-bold text-gray-800 mb-4">📈 Phân tích chênh lệch</h3>
            
            <div class="space-y-3">
              <div v-for="(comparison, index) in costComparisons" :key="index"
                class="bg-white rounded-lg p-4 shadow-sm">
                <div class="flex items-center justify-between mb-2">
                  <span class="font-medium text-gray-700">
                    {{ comparison.from }} → {{ comparison.to }}
                  </span>
                  <span class="text-xs text-gray-500">So sánh</span>
                </div>
                
                <div class="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <div class="text-xs text-gray-500 mb-1">Chênh lệch giá</div>
                    <div class="font-bold" :class="comparison.priceDiff > 0 ? 'text-blue-600' : 'text-gray-600'">
                      {{ comparison.priceDiff > 0 ? '+' : '' }}{{ formatPrice(comparison.priceDiff) }}
                    </div>
                  </div>
                  <div>
                    <div class="text-xs text-gray-500 mb-1">Chênh lệch chi phí</div>
                    <div class="font-bold" :class="comparison.costDiff > 0 ? 'text-orange-600' : 'text-gray-600'">
                      {{ comparison.costDiff > 0 ? '+' : '' }}{{ formatPrice(comparison.costDiff) }}
                    </div>
                  </div>
                  <div>
                    <div class="text-xs text-gray-500 mb-1">Chênh lệch lợi nhuận</div>
                    <div class="font-bold" :class="getProfitColor(comparison.profitDiff)">
                      {{ comparison.profitDiff > 0 ? '+' : '' }}{{ formatPrice(comparison.profitDiff) }}
                    </div>
                  </div>
                  <div>
                    <div class="text-xs text-gray-500 mb-1">Chênh lệch tỷ suất LN</div>
                    <div class="font-bold" :class="getProfitMarginColor(comparison.marginDiff)">
                      {{ comparison.marginDiff > 0 ? '+' : '' }}{{ formatPercentage(comparison.marginDiff) }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Insights -->
          <div class="bg-yellow-50 border-l-4 border-yellow-400 rounded-lg p-4">
            <h3 class="text-sm font-bold text-yellow-800 mb-2">💡 Nhận xét</h3>
            <ul class="text-sm text-yellow-700 space-y-1">
              <li v-for="insight in insights" :key="insight">• {{ insight }}</li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="px-6 py-4 border-t bg-gray-50 flex justify-end flex-shrink-0">
        <button @click="$emit('close')"
          class="px-6 py-2 bg-gray-500 text-white rounded-lg font-medium hover:bg-gray-600 active:bg-gray-700 transition-colors">
          Đóng
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { menuService } from '../services/menu'
import { formatPrice, formatPercentage } from '../utils/formatters'

const props = defineProps({
  isOpen: {
    type: Boolean,
    required: true
  },
  menuItemId: {
    type: String,
    default: null
  }
})

defineEmits(['close'])

// State
const loading = ref(false)
const error = ref(null)
const menuItem = ref(null)

// Computed
const sortedVariants = computed(() => {
  if (!menuItem.value?.variants) return []
  
  // Sort by profit (descending)
  return [...menuItem.value.variants].sort((a, b) => {
    const profitA = calculateProfit(a)
    const profitB = calculateProfit(b)
    return profitB - profitA
  })
})

const mostProfitableVariant = computed(() => {
  if (!sortedVariants.value.length) return null
  return sortedVariants.value[0]
})

const maxProfit = computed(() => {
  if (!mostProfitableVariant.value) return 0
  return calculateProfit(mostProfitableVariant.value)
})

const maxProfitMargin = computed(() => {
  if (!mostProfitableVariant.value) return 0
  return calculateProfitMargin(mostProfitableVariant.value)
})

const costComparisons = computed(() => {
  if (!menuItem.value?.variants || menuItem.value.variants.length < 2) return []
  
  const comparisons = []
  const variants = menuItem.value.variants
  
  // Compare consecutive sizes
  for (let i = 0; i < variants.length - 1; i++) {
    const from = variants[i]
    const to = variants[i + 1]
    
    comparisons.push({
      from: from.name,
      to: to.name,
      priceDiff: to.price - from.price,
      costDiff: to.current_cost - from.current_cost,
      profitDiff: calculateProfit(to) - calculateProfit(from),
      marginDiff: calculateProfitMargin(to) - calculateProfitMargin(from)
    })
  }
  
  return comparisons
})

const insights = computed(() => {
  if (!menuItem.value?.variants || menuItem.value.variants.length === 0) return []
  
  const insights = []
  const variants = menuItem.value.variants
  
  // Find most profitable
  if (mostProfitableVariant.value) {
    insights.push(`Size "${mostProfitableVariant.value.name}" có lợi nhuận cao nhất: ${formatPrice(maxProfit.value)}`)
  }
  
  // Check if all variants are profitable
  const unprofitableVariants = variants.filter(v => calculateProfit(v) < 0)
  if (unprofitableVariants.length > 0) {
    insights.push(`⚠️ ${unprofitableVariants.length} size đang bị lỗ: ${unprofitableVariants.map(v => v.name).join(', ')}`)
  } else {
    insights.push(`✅ Tất cả các size đều có lợi nhuận`)
  }
  
  // Check profit margin consistency
  const margins = variants.map(v => calculateProfitMargin(v))
  const avgMargin = margins.reduce((a, b) => a + b, 0) / margins.length
  const marginVariance = margins.some(m => Math.abs(m - avgMargin) > 10)
  
  if (marginVariance) {
    insights.push(`Tỷ suất lợi nhuận giữa các size chênh lệch lớn (>10%)`)
  } else {
    insights.push(`Tỷ suất lợi nhuận giữa các size tương đối đồng đều`)
  }
  
  return insights
})

// Helper functions
const calculateProfit = (variant) => {
  return variant.price - variant.current_cost
}

const calculateProfitMargin = (variant) => {
  if (!variant.price || variant.price === 0) return 0
  return ((variant.price - variant.current_cost) / variant.price) * 100
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

const getCostStatusLabel = (status) => {
  const labels = {
    'FINAL': '✓ Chính thức',
    'ESTIMATED': '~ Ước tính',
    'INCOMPLETE': '⚠ Thiếu'
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

// Fetch menu item data
const fetchMenuItem = async () => {
  if (!props.menuItemId) return
  
  loading.value = true
  error.value = null
  
  try {
    const response = await menuService.getMenuItem(props.menuItemId)
    menuItem.value = response.data
  } catch (err) {
    console.error('Error fetching menu item:', err)
    error.value = err.response?.data?.error || 'Không thể tải dữ liệu món'
  } finally {
    loading.value = false
  }
}

// Watch for modal open
watch(() => props.isOpen, (isOpen) => {
  if (isOpen && props.menuItemId) {
    fetchMenuItem()
  }
}, { immediate: true })

// Watch for menuItemId changes
watch(() => props.menuItemId, (newId) => {
  if (props.isOpen && newId) {
    fetchMenuItem()
  }
})
</script>
