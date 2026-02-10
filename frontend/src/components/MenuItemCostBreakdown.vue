<template>
  <div v-if="isOpen" 
    class="fixed inset-0 z-50 flex items-end justify-center bg-black bg-opacity-50"
    @click.self="handleClose">
    <div class="bg-white rounded-t-3xl w-full max-w-2xl max-h-[80vh] overflow-hidden flex flex-col animate-slide-up"
      style="padding-bottom: env(safe-area-inset-bottom)">
      <!-- Modal Header -->
      <div class="px-4 py-4 border-b flex items-center justify-between flex-shrink-0">
        <h2 class="text-lg font-bold">Chi tiết chi phí</h2>
        <button @click="handleClose" class="text-gray-500 hover:text-gray-700 text-2xl">
          ×
        </button>
      </div>
      
      <!-- Modal Content -->
      <div class="flex-1 overflow-y-auto px-4 py-4">
        <!-- Loading State -->
        <div v-if="loading" class="text-center py-8">
          <div class="text-4xl mb-2">⏳</div>
          <p class="text-gray-500">Đang tải...</p>
        </div>
        
        <!-- Error State -->
        <div v-else-if="error" class="text-center py-8">
          <div class="text-4xl mb-2">❌</div>
          <p class="text-red-500">{{ error }}</p>
        </div>
        
        <!-- Content -->
        <div v-else-if="breakdown">
          <!-- Menu Item Info -->
          <div class="bg-gray-50 rounded-xl p-4 mb-4">
            <h3 class="font-bold text-lg mb-2">{{ breakdown.menu_item?.name }}</h3>
            <div class="grid grid-cols-2 gap-3 text-sm">
              <div>
                <span class="text-gray-600">Giá bán:</span>
                <span class="font-bold ml-2">{{ formatPrice(breakdown.menu_item?.price) }}</span>
              </div>
              <div>
                <span class="text-gray-600">Chi phí:</span>
                <span class="font-bold ml-2 text-blue-600">{{ formatPrice(breakdown.total_cost) }}</span>
              </div>
            </div>
          </div>
          
          <!-- Ingredients List -->
          <div class="mb-4">
            <h4 class="font-bold mb-2">Nguyên liệu:</h4>
            <div class="space-y-2">
              <div v-for="(ingredient, index) in breakdown.ingredients" :key="index"
                class="bg-white border border-gray-200 rounded-lg p-3"
                :class="{ 'border-red-300 bg-red-50': hasIncompleteCost(ingredient) }">
                <div class="flex justify-between items-start mb-2">
                  <div class="flex-1">
                    <div class="font-medium">{{ ingredient.name }}</div>
                    <div class="text-xs text-gray-500">
                      {{ ingredient.quantity }} {{ ingredient.unit }}
                      <span v-if="hasNonDefaultConversion(ingredient)">
                        × {{ ingredient.conversion_rate }} (quy đổi)
                      </span>
                      <span v-if="hasWastage(ingredient)">
                        + {{ ingredient.wastage_percentage }}% (hao hụt)
                      </span>
                    </div>
                  </div>
                  <div class="text-right">
                    <div class="font-bold text-blue-600">{{ formatPrice(ingredient.total_cost) }}</div>
                    <div class="text-xs text-gray-500">{{ formatPrice(ingredient.cost_per_unit) }}/{{ ingredient.unit }}</div>
                  </div>
                </div>
                
                <!-- Warning for missing cost -->
                <div v-if="hasIncompleteCost(ingredient)" 
                  class="mt-2 pt-2 border-t border-red-200 text-xs text-red-600 flex items-center gap-1">
                  <span>⚠️</span>
                  <span>Thiếu giá nguyên liệu</span>
                </div>
              </div>
            </div>
          </div>
          
          <!-- Total Cost Summary -->
          <div class="bg-blue-50 rounded-xl p-4 border-2 border-blue-200">
            <div class="flex justify-between items-center">
              <span class="font-bold text-lg">Tổng chi phí:</span>
              <span class="font-bold text-xl text-blue-600">{{ formatPrice(breakdown.total_cost) }}</span>
            </div>
            
            <!-- Warning if any ingredient has incomplete cost -->
            <div v-if="hasAnyIncompleteCost" class="mt-3 pt-3 border-t border-blue-200">
              <div class="flex items-center gap-2 text-sm text-red-600">
                <span>⚠️</span>
                <span class="font-medium">Một số nguyên liệu thiếu giá, chi phí có thể không chính xác</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { menuCostService } from '../services/menuCost'
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
  }
})

// Emits
const emit = defineEmits(['close'])

// State
const loading = ref(false)
const error = ref(null)
const breakdown = ref(null)

// Computed
const hasAnyIncompleteCost = computed(() => {
  if (!breakdown.value?.ingredients) return false
  return breakdown.value.ingredients.some(ingredient => hasIncompleteCost(ingredient))
})

// Helper functions
const hasIncompleteCost = (ingredient) => {
  return !ingredient.cost_per_unit || ingredient.cost_per_unit === 0
}

const hasNonDefaultConversion = (ingredient) => {
  return ingredient.conversion_rate && ingredient.conversion_rate !== 1
}

const hasWastage = (ingredient) => {
  return ingredient.wastage_percentage && ingredient.wastage_percentage > 0
}

const handleClose = () => {
  emit('close')
}

// Fetch cost breakdown
const fetchCostBreakdown = async () => {
  if (!props.menuItemId) return
  
  loading.value = true
  error.value = null
  breakdown.value = null
  
  try {
    const response = await menuCostService.getMenuCostDetail(props.menuItemId)
    breakdown.value = response
  } catch (err) {
    console.error('Error fetching cost breakdown:', err)
    error.value = err.response?.data?.error || 'Không thể tải chi tiết chi phí'
  } finally {
    loading.value = false
  }
}

// Watch for changes in isOpen and menuItemId
watch(() => [props.isOpen, props.menuItemId], ([isOpen, menuItemId]) => {
  if (isOpen && menuItemId) {
    fetchCostBreakdown()
  }
}, { immediate: true })
</script>

<style scoped>
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
