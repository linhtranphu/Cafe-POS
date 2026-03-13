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
          <!-- Menu Item Name -->
          <div class="bg-gray-50 rounded-xl p-4 mb-4">
            <h3 class="font-bold text-lg">{{ breakdown.menu_item_name }}</h3>
          </div>

          <!-- Multi-variant: tab selector -->
          <div v-if="breakdown.has_variants && breakdown.variants?.length">
            <!-- Variant Tabs -->
            <div class="flex gap-2 mb-4 overflow-x-auto pb-1">
              <button
                v-for="v in breakdown.variants"
                :key="v.variant_id"
                @click="selectedVariantId = v.variant_id"
                class="flex-shrink-0 px-3 py-1.5 rounded-full text-sm font-medium border transition-colors"
                :class="selectedVariantId === v.variant_id
                  ? 'bg-purple-600 text-white border-purple-600'
                  : 'bg-white text-gray-600 border-gray-300 hover:border-purple-400'">
                {{ v.variant_name }}
              </button>
            </div>

            <!-- Selected Variant Breakdown -->
            <div v-if="selectedVariant">
              <!-- Price / Cost row -->
              <div class="grid grid-cols-2 gap-3 text-sm mb-4">
                <div>
                  <span class="text-gray-600">Giá bán:</span>
                  <span class="font-bold ml-2">{{ formatPrice(selectedVariant.price) }}</span>
                </div>
                <div>
                  <span class="text-gray-600">Chi phí:</span>
                  <span class="font-bold ml-2 text-blue-600">{{ formatPrice(selectedVariant.total_cost) }}</span>
                </div>
              </div>

              <IngredientList :ingredients="selectedVariant.ingredients" />

              <!-- Total -->
              <div class="bg-blue-50 rounded-xl p-4 border-2 border-blue-200 mt-4">
                <div class="flex justify-between items-center">
                  <span class="font-bold text-lg">Tổng chi phí:</span>
                  <span class="font-bold text-xl text-blue-600">{{ formatPrice(selectedVariant.total_cost) }}</span>
                </div>
                <div v-if="variantHasIncompleteCost(selectedVariant)" class="mt-3 pt-3 border-t border-blue-200">
                  <div class="flex items-center gap-2 text-sm text-red-600">
                    <span>⚠️</span>
                    <span class="font-medium">Một số nguyên liệu thiếu giá, chi phí có thể không chính xác</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Single-size item -->
          <div v-else>
            <div class="grid grid-cols-2 gap-3 text-sm mb-4">
              <div>
                <span class="text-gray-600">Giá bán:</span>
                <span class="font-bold ml-2">{{ formatPrice(breakdown.price) }}</span>
              </div>
              <div>
                <span class="text-gray-600">Chi phí:</span>
                <span class="font-bold ml-2 text-blue-600">{{ formatPrice(breakdown.total_cost) }}</span>
              </div>
            </div>

            <IngredientList :ingredients="breakdown.ingredients" />

            <!-- Total -->
            <div class="bg-blue-50 rounded-xl p-4 border-2 border-blue-200 mt-4">
              <div class="flex justify-between items-center">
                <span class="font-bold text-lg">Tổng chi phí:</span>
                <span class="font-bold text-xl text-blue-600">{{ formatPrice(breakdown.total_cost) }}</span>
              </div>
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
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { menuCostService } from '../services/menuCost'
import { formatPrice } from '../utils/formatters'

// Sub-component for ingredient list (defined inline)
const IngredientList = {
  name: 'IngredientList',
  props: { ingredients: Array },
  setup(props) {
    const hasIncompleteCost = (ing) => !ing.cost_per_unit || ing.cost_per_unit === 0
    const hasNonDefaultConversion = (ing) => ing.conversion_rate && ing.conversion_rate !== 1
    const hasWastage = (ing) => ing.wastage_percentage && ing.wastage_percentage > 0
    return { hasIncompleteCost, hasNonDefaultConversion, hasWastage, formatPrice }
  },
  template: `
    <div class="mb-4">
      <h4 class="font-bold mb-2">Nguyên liệu:</h4>
      <div class="space-y-2">
        <div v-for="(ingredient, index) in ingredients" :key="index"
          class="bg-white border border-gray-200 rounded-lg p-3"
          :class="{ 'border-red-300 bg-red-50': hasIncompleteCost(ingredient) }">
          <div class="flex justify-between items-start mb-2">
            <div class="flex-1">
              <div class="flex items-center gap-2 flex-wrap">
                <div class="font-medium">{{ ingredient.name }}</div>
                <span v-if="ingredient.deduct_inventory" class="text-xs bg-red-100 text-red-700 px-1.5 py-0.5 rounded">📦 Trừ kho</span>
                <span v-else class="text-xs bg-gray-100 text-gray-500 px-1.5 py-0.5 rounded">💰 Chỉ cost</span>
              </div>
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
          <div v-if="hasIncompleteCost(ingredient)"
            class="mt-2 pt-2 border-t border-red-200 text-xs text-red-600 flex items-center gap-1">
            <span>⚠️</span>
            <span>Thiếu giá nguyên liệu</span>
          </div>
        </div>
      </div>
    </div>
  `
}

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
const selectedVariantId = ref(null)

// Computed
const selectedVariant = computed(() => {
  if (!breakdown.value?.variants?.length) return null
  return breakdown.value.variants.find(v => v.variant_id === selectedVariantId.value)
    ?? breakdown.value.variants[0]
})

const hasAnyIncompleteCost = computed(() => {
  if (!breakdown.value?.ingredients) return false
  return breakdown.value.ingredients.some(ing => !ing.cost_per_unit || ing.cost_per_unit === 0)
})

const variantHasIncompleteCost = (variant) => {
  return variant?.ingredients?.some(ing => !ing.cost_per_unit || ing.cost_per_unit === 0)
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
  selectedVariantId.value = null

  try {
    const response = await menuCostService.getMenuCostDetail(props.menuItemId)
    breakdown.value = response
    // Pre-select first variant
    if (response.has_variants && response.variants?.length) {
      selectedVariantId.value = response.variants[0].variant_id
    }
  } catch (err) {
    console.error('Error fetching cost breakdown:', err)
    error.value = err.response?.data?.error || 'Không thể tải chi tiết chi phí'
  } finally {
    loading.value = false
  }
}

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
