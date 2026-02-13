<template>
  <div v-if="isOpen" 
    class="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black bg-opacity-50"
    @click.self="closeModal">
    
    <!-- Modal Container -->
    <div class="bg-white w-full sm:max-w-2xl sm:rounded-2xl rounded-t-2xl max-h-[90vh] overflow-hidden flex flex-col"
      style="padding-bottom: env(safe-area-inset-bottom)">
      
      <!-- Header -->
      <div class="sticky top-0 bg-white border-b px-4 py-4 flex items-center justify-between">
        <h2 class="text-lg font-bold text-gray-800">📊 Chi tiết chi phí</h2>
        <button @click="closeModal" 
          class="p-2 hover:bg-gray-100 rounded-full transition-colors">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Content -->
      <div class="flex-1 overflow-y-auto px-4 py-4">
        <!-- Loading State -->
        <div v-if="loading" class="space-y-3">
          <div class="animate-pulse">
            <div class="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
            <div class="h-4 bg-gray-200 rounded w-1/2"></div>
          </div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="text-center py-8">
          <div class="text-4xl mb-2">❌</div>
          <p class="text-red-500 font-medium">{{ error }}</p>
        </div>

        <!-- Content -->
        <div v-else-if="costData">
          <!-- Menu Item Name -->
          <div class="mb-4 pb-3 border-b">
            <h3 class="text-xl font-bold text-gray-800">{{ costData.menu_item_name }}</h3>
            <span class="text-sm text-gray-500">
              {{ costData.has_variants ? 'Nhiều size' : 'Một size' }}
            </span>
          </div>

          <!-- Single-Size Item -->
          <div v-if="!costData.has_variants">
            <div class="mb-4 p-3 bg-blue-50 rounded-lg">
              <div class="flex justify-between items-center mb-2">
                <span class="text-sm text-gray-600">Giá bán:</span>
                <span class="font-bold text-blue-600">{{ formatPrice(costData.price) }}</span>
              </div>
              <div class="flex justify-between items-center mb-2">
                <span class="text-sm text-gray-600">Tổng chi phí:</span>
                <span class="font-bold text-orange-600">{{ formatPrice(costData.total_cost) }}</span>
              </div>
              <div class="flex justify-between items-center">
                <span class="text-sm text-gray-600">Trạng thái:</span>
                <span 
                  :class="getCostStatusClass(costData.cost_status)"
                  class="px-2 py-0.5 rounded text-xs font-medium">
                  {{ getCostStatusLabel(costData.cost_status) }}
                </span>
              </div>
            </div>

            <!-- Ingredients Breakdown -->
            <div class="space-y-2">
              <h4 class="font-bold text-gray-700 mb-2">Nguyên liệu:</h4>
              <div v-for="(ing, index) in costData.ingredients" :key="index"
                class="p-3 bg-gray-50 rounded-lg">
                <div class="flex justify-between items-start mb-2">
                  <span class="font-medium text-gray-800">{{ ing.name }}</span>
                  <span class="font-bold text-orange-600">{{ formatPrice(ing.total_cost) }}</span>
                </div>
                
                <!-- Formula Breakdown -->
                <div class="text-xs text-gray-600 space-y-1 mt-2 p-2 bg-white rounded border border-gray-200">
                  <div class="flex justify-between">
                    <span>Số lượng:</span>
                    <span class="font-mono">{{ ing.quantity }} {{ ing.unit }}</span>
                  </div>
                  <div class="flex justify-between">
                    <span>Giá/đơn vị:</span>
                    <span class="font-mono">{{ formatPrice(ing.cost_per_unit) }}</span>
                  </div>
                  <div class="flex justify-between">
                    <span>Tỷ lệ chuyển đổi:</span>
                    <span class="font-mono">{{ ing.conversion_rate.toFixed(4) }}</span>
                  </div>
                  <div class="flex justify-between">
                    <span>Hao hụt:</span>
                    <span class="font-mono">{{ ing.wastage_percentage }}%</span>
                  </div>
                  <div class="pt-2 mt-2 border-t border-gray-300">
                    <div class="text-[10px] text-gray-500 mb-1">Công thức:</div>
                    <div class="font-mono text-[10px] break-all">
                      {{ ing.quantity }} × {{ formatPrice(ing.cost_per_unit) }} × {{ ing.conversion_rate.toFixed(4) }} × (1 + {{ ing.wastage_percentage }}/100)
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Multi-Size Item with Variants -->
          <div v-else class="space-y-4">
            <div v-for="variant in costData.variants" :key="variant.variant_id"
              class="border-2 rounded-xl p-4"
              :class="variant.variant_id === defaultVariantId ? 'border-purple-300 bg-purple-50' : 'border-gray-200'">
              
              <!-- Variant Header -->
              <div class="flex items-center justify-between mb-3 pb-2 border-b">
                <div>
                  <h4 class="font-bold text-lg">{{ variant.variant_name }}</h4>
                  <span v-if="variant.variant_id === defaultVariantId"
                    class="text-xs px-2 py-0.5 bg-purple-500 text-white rounded">
                    Mặc định
                  </span>
                </div>
                <span 
                  :class="getCostStatusClass(variant.cost_status)"
                  class="px-2 py-1 rounded text-xs font-medium">
                  {{ getCostStatusLabel(variant.cost_status) }}
                </span>
              </div>

              <!-- Variant Summary -->
              <div class="mb-3 p-3 bg-blue-50 rounded-lg">
                <div class="flex justify-between items-center mb-2">
                  <span class="text-sm text-gray-600">Giá bán:</span>
                  <span class="font-bold text-blue-600">{{ formatPrice(variant.price) }}</span>
                </div>
                <div class="flex justify-between items-center">
                  <span class="text-sm text-gray-600">Tổng chi phí:</span>
                  <span class="font-bold text-orange-600">{{ formatPrice(variant.total_cost) }}</span>
                </div>
              </div>

              <!-- Variant Ingredients -->
              <div class="space-y-2">
                <h5 class="font-bold text-sm text-gray-700 mb-2">Nguyên liệu:</h5>
                <div v-for="(ing, index) in variant.ingredients" :key="index"
                  class="p-3 bg-white rounded-lg border border-gray-200">
                  <div class="flex justify-between items-start mb-2">
                    <span class="font-medium text-gray-800 text-sm">{{ ing.name }}</span>
                    <span class="font-bold text-orange-600 text-sm">{{ formatPrice(ing.total_cost) }}</span>
                  </div>
                  
                  <!-- Formula Breakdown -->
                  <div class="text-xs text-gray-600 space-y-1 mt-2 p-2 bg-gray-50 rounded">
                    <div class="flex justify-between">
                      <span>Số lượng:</span>
                      <span class="font-mono">{{ ing.quantity }} {{ ing.unit }}</span>
                    </div>
                    <div class="flex justify-between">
                      <span>Giá/đơn vị:</span>
                      <span class="font-mono">{{ formatPrice(ing.cost_per_unit) }}</span>
                    </div>
                    <div class="flex justify-between">
                      <span>Tỷ lệ chuyển đổi:</span>
                      <span class="font-mono">{{ ing.conversion_rate.toFixed(4) }}</span>
                    </div>
                    <div class="flex justify-between">
                      <span>Hao hụt:</span>
                      <span class="font-mono">{{ ing.wastage_percentage }}%</span>
                    </div>
                    <div class="pt-2 mt-2 border-t border-gray-200">
                      <div class="text-[10px] text-gray-500 mb-1">Công thức:</div>
                      <div class="font-mono text-[10px] break-all">
                        {{ ing.quantity }} × {{ formatPrice(ing.cost_per_unit) }} × {{ ing.conversion_rate.toFixed(4) }} × (1 + {{ ing.wastage_percentage }}/100)
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="sticky bottom-0 bg-white border-t px-4 py-3">
        <button @click="closeModal"
          class="w-full bg-blue-500 text-white py-3 rounded-lg font-medium active:bg-blue-600 hover:bg-blue-600 transition-colors">
          Đóng
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { menuService } from '../services/menu'
import { formatPrice } from '../utils/formatters'

// Props
const props = defineProps({
  isOpen: {
    type: Boolean,
    required: true
  },
  menuItemId: {
    type: String,
    default: null
  },
  defaultVariantId: {
    type: String,
    default: null
  }
})

// Emits
const emit = defineEmits(['close'])

// State
const loading = ref(false)
const error = ref(null)
const costData = ref(null)

// Methods
const closeModal = () => {
  emit('close')
}

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

const fetchCostBreakdown = async () => {
  if (!props.menuItemId) return
  
  loading.value = true
  error.value = null
  costData.value = null
  
  try {
    const response = await menuService.getCostBreakdown(props.menuItemId)
    costData.value = response.data
  } catch (err) {
    console.error('Error fetching cost breakdown:', err)
    error.value = err.response?.data?.error || 'Không thể tải chi tiết chi phí'
  } finally {
    loading.value = false
  }
}

// Watch for modal open and menuItemId changes
watch(() => [props.isOpen, props.menuItemId], ([isOpen, menuItemId]) => {
  if (isOpen && menuItemId) {
    fetchCostBreakdown()
  }
}, { immediate: true })
</script>

<style scoped>
/* Smooth animations */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
