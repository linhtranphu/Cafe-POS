<template>
  <transition name="slide-up">
    <div v-if="modelValue" class="fixed inset-0 bg-white z-50 overflow-hidden flex flex-col">
      <!-- Header -->
      <div class="bg-gradient-to-r from-blue-500 to-blue-600 text-white px-4 py-3 shadow-lg" 
        style="padding-top: max(1rem, env(safe-area-inset-top))">
        <div class="flex items-center justify-between">
          <button @click="handleCancel" 
            class="text-2xl p-2 active:bg-blue-600 rounded-lg touch-manipulation">
            ←
          </button>
          <h2 class="text-xl font-bold">🍽️ Tạo Order</h2>
          <div class="w-10"></div> <!-- Spacer for centering -->
        </div>
      </div>

      <!-- Category Filter Buttons -->
      <div class="bg-white border-b px-4 py-2 flex flex-wrap gap-2 flex-shrink-0">
        <button
          @click="selectedCategory = 'all'"
          :class="[
            'py-2 px-4 rounded-xl text-sm font-semibold transition-all touch-manipulation active:scale-95',
            selectedCategory === 'all'
              ? 'bg-blue-500 text-white shadow-md'
              : 'bg-gray-100 text-gray-600 active:bg-gray-200'
          ]">
          📋 Tất cả
        </button>
        <button
          v-for="cat in props.categories.filter(c => c.id !== 'all')" :key="cat.id"
          @click="selectedCategory = cat.id"
          :class="[
            'py-2 px-4 rounded-xl text-sm font-semibold transition-all touch-manipulation active:scale-95',
            selectedCategory === cat.id
              ? 'bg-blue-500 text-white shadow-md'
              : 'bg-gray-100 text-gray-600 active:bg-gray-200'
          ]">
          {{ cat.icon }} {{ cat.name }}
        </button>
      </div>

      <!-- Grouped Menu Items -->
      <div class="flex-1 overflow-y-auto bg-gray-50" 
        :style="{ paddingBottom: totalItems > 0 ? '100px' : '20px' }">
        
        <!-- Empty State -->
        <div v-if="filteredGroupedItems.length === 0" class="text-center py-20">
          <div class="text-6xl mb-3">🍽️</div>
          <p class="text-gray-500 font-medium">Không có món nào</p>
        </div>

        <!-- Category Groups -->
        <div v-for="group in filteredGroupedItems" :key="group.category.id" class="mb-6">
          <!-- Category Header (Sticky) -->
          <div class="sticky top-0 z-10 bg-gradient-to-b from-gray-50 to-gray-100 px-4 py-3 border-b-2 border-gray-200">
            <div class="flex items-center gap-2">
              <span class="text-2xl">{{ group.category.icon }}</span>
              <span class="font-bold text-gray-800 text-lg">{{ group.category.name }}</span>
              <span class="text-base text-gray-500">({{ group.items.length }})</span>
            </div>
          </div>

          <!-- Items Grid -->
          <div class="grid grid-cols-2 gap-3 px-4 py-3">
            <template v-for="item in group.items" :key="item.id">
              <!-- Single-size item -->
              <div v-if="!item.has_variants"
                :class="[
                  'bg-white rounded-2xl p-4 shadow-sm transition-all flex flex-col',
                  getItemQty(item.id) > 0 
                    ? 'border-2 border-blue-500 shadow-md' 
                    : 'border border-gray-200'
                ]">
                <!-- Item Info -->
                <div class="flex-1 mb-3">
                  <h3 class="font-bold text-gray-900 mb-2 line-clamp-2 min-h-[2.5rem] leading-tight text-base">
                    {{ item.name }}
                  </h3>
                  <p class="text-lg font-bold text-blue-600">
                    {{ formatPrice(item.price) }}
                  </p>
                </div>

                <!-- Quantity Controls -->
                <div class="mt-auto">
                  <!-- Add Button (when qty = 0) -->
                  <button v-if="getItemQty(item.id) === 0"
                    @click="addItem(item.id)"
                    class="w-full py-2.5 bg-gradient-to-r from-blue-500 to-blue-600 text-white font-bold rounded-xl active:from-blue-600 active:to-blue-700 active:scale-95 transition-all shadow-sm">
                    + Thêm
                  </button>

                  <!-- Inline Controls (when qty > 0) -->
                  <div v-else class="space-y-2">
                    <!-- Quantity Controls -->
                    <div class="flex items-center justify-between gap-2">
                      <button @click="removeItem(item.id)"
                        class="w-10 h-10 bg-white border-2 border-blue-500 text-blue-600 rounded-full text-xl font-bold active:bg-blue-50 active:scale-95 transition-all shadow-sm">
                        −
                      </button>
                      <span class="flex-1 text-center text-xl font-bold text-blue-600">
                        {{ getItemQty(item.id) }}
                      </span>
                      <button @click="addItem(item.id)"
                        class="w-10 h-10 bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-full text-xl font-bold active:from-blue-600 active:to-blue-700 active:scale-95 transition-all shadow-md">
                        +
                      </button>
                    </div>
                    
                    <!-- Sugar Level Selector -->
                    <div class="flex items-center gap-1">
                      <span class="text-xs text-gray-500 mr-1">🧊</span>
                      <button v-for="level in sugarLevels" :key="level.value"
                        @click="setSugarLevel(item.id, level.value)"
                        :class="[
                          'flex-1 py-1.5 rounded-lg text-xs font-semibold transition-all active:scale-95',
                          getSugarLevel(item.id) === level.value
                            ? 'bg-amber-500 text-white shadow-sm'
                            : 'bg-gray-100 text-gray-600 active:bg-gray-200'
                        ]">
                        {{ level.label }}
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Multi-size item (with variants) -->
              <div v-else
                :class="[
                  'bg-white rounded-2xl p-4 shadow-sm transition-all flex flex-col',
                  getItemQtyWithVariants(item.id) > 0 
                    ? 'border-2 border-blue-500 shadow-md' 
                    : 'border border-gray-200'
                ]">
                <!-- Item Info -->
                <div class="mb-3">
                  <h3 class="font-bold text-gray-900 mb-2 line-clamp-2 min-h-[2.5rem] leading-tight text-base">
                    {{ item.name }}
                  </h3>
                </div>

                <!-- Variants -->
                <div class="space-y-3 flex-1">
                  <div v-for="variant in item.variants" :key="variant.id"
                    class="space-y-1">
                    <div class="flex items-center justify-between gap-2">
                      <div class="flex-1">
                        <p class="text-sm font-semibold text-gray-700">{{ variant.name }}</p>
                        <p class="text-base font-bold text-blue-600">{{ formatPrice(variant.price) }}</p>
                      </div>
                      
                      <!-- Variant Controls -->
                      <div v-if="getItemQty(item.id, variant.id) === 0">
                        <button @click="addItem(item.id, variant.id)"
                          class="px-3 py-1.5 bg-blue-500 text-white text-xs font-bold rounded-lg active:bg-blue-600 active:scale-95 transition-all">
                          +
                        </button>
                      </div>
                      <div v-else class="flex items-center gap-1">
                        <button @click="removeItem(item.id, variant.id)"
                          class="w-7 h-7 bg-white border-2 border-blue-500 text-blue-600 rounded-full text-sm font-bold active:bg-blue-50">
                          −
                        </button>
                        <span class="w-6 text-center text-sm font-bold text-blue-600">
                          {{ getItemQty(item.id, variant.id) }}
                        </span>
                        <button @click="addItem(item.id, variant.id)"
                          class="w-7 h-7 bg-blue-500 text-white rounded-full text-sm font-bold active:bg-blue-600">
                          +
                        </button>
                      </div>
                    </div>
                    
                    <!-- Sugar Level for variant (when qty > 0) -->
                    <div v-if="getItemQty(item.id, variant.id) > 0" class="flex items-center gap-1">
                      <span class="text-xs text-gray-400">🧊</span>
                      <button v-for="level in sugarLevels" :key="level.value"
                        @click="setSugarLevel(item.id, level.value, variant.id)"
                        :class="[
                          'flex-1 py-1 rounded text-xs font-semibold transition-all active:scale-95',
                          getSugarLevel(item.id, variant.id) === level.value
                            ? 'bg-amber-500 text-white shadow-sm'
                            : 'bg-gray-100 text-gray-500 active:bg-gray-200'
                        ]">
                        {{ level.label }}
                      </button>
                    </div>
                  </div>
                </div>

                <!-- Total for this item -->
                <div v-if="getItemQtyWithVariants(item.id) > 0" 
                  class="mt-3 pt-3 border-t border-gray-200 text-center">
                  <span class="text-xs font-bold text-blue-600">
                    🛒 {{ getItemQtyWithVariants(item.id) }} món
                  </span>
                </div>
              </div>
            </template>
          </div>
        </div>
      </div>

      <!-- Floating Total Button -->
      <transition name="slide-up">
        <button v-if="totalItems > 0"
          @click="handleConfirm"
          class="fixed bottom-0 left-0 right-0 mx-4 mb-4 bg-gradient-to-r from-green-500 to-green-600 text-white rounded-2xl shadow-2xl active:scale-98 transition-all"
          style="padding-bottom: max(1rem, env(safe-area-inset-bottom))">
          <div class="px-6 py-4 flex items-center justify-between">
            <div class="flex items-center gap-3">
              <span class="text-2xl">🛒</span>
              <div class="text-left">
                <p class="text-sm font-medium opacity-90">{{ totalItems }} món</p>
                <p class="text-xl font-bold">{{ formatPrice(totalPrice) }}</p>
              </div>
            </div>
            <div class="flex items-center gap-2 text-lg font-bold">
              <span>Xác nhận</span>
              <span>→</span>
            </div>
          </div>
        </button>
      </transition>
    </div>
  </transition>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  modelValue: Boolean,
  menuItems: Array,
  categories: Array
})

const emit = defineEmits(['update:modelValue', 'confirm'])

// Cart state - using object for better performance
const cart = ref({})
const selectedCategory = ref('all')

// Sugar level state - DISABLED (bỏ mục Đường)
const sugarLevels = [
  // { value: 25, label: '25%' },
  // { value: 50, label: '50%' },
  // { value: 100, label: '100%' }
]
const itemSugarLevels = ref({}) // { itemId: level, itemId_variantId: level }

// Helper functions
const getItemQty = (itemId, variantId = null) => {
  const key = variantId ? `${itemId}_${variantId}` : itemId
  return cart.value[key] || 0
}

const addItem = (itemId, variantId = null) => {
  const key = variantId ? `${itemId}_${variantId}` : itemId
  cart.value[key] = (cart.value[key] || 0) + 1
  
  // Set default sugar level to 100% if not set
  if (!itemSugarLevels.value[key]) {
    itemSugarLevels.value[key] = 100
  }
}

const removeItem = (itemId, variantId = null) => {
  const key = variantId ? `${itemId}_${variantId}` : itemId
  if (cart.value[key] > 1) {
    cart.value[key]--
  } else {
    delete cart.value[key]
    delete itemSugarLevels.value[key] // Also delete sugar level
  }
}

const getItemQtyWithVariants = (itemId) => {
  return Object.keys(cart.value)
    .filter(key => key.startsWith(itemId))
    .reduce((sum, key) => sum + cart.value[key], 0)
}

// Sugar level methods
const getSugarLevel = (itemId, variantId = null) => {
  const key = variantId ? `${itemId}_${variantId}` : itemId
  return itemSugarLevels.value[key] || 100 // Default 100%
}

const setSugarLevel = (itemId, level, variantId = null) => {
  const key = variantId ? `${itemId}_${variantId}` : itemId
  itemSugarLevels.value[key] = level
}

// Computed properties
const groupedItems = computed(() => {
  const groups = {}
  
  props.menuItems.forEach(item => {
    const categoryId = item.category
    if (!groups[categoryId]) {
      const category = props.categories.find(c => c.id === categoryId)
      if (category) {
        groups[categoryId] = {
          category: category,
          items: []
        }
      }
    }
    if (groups[categoryId]) {
      groups[categoryId].items.push(item)
    }
  })
  
  return Object.values(groups).filter(g => g.items.length > 0)
})

const filteredGroupedItems = computed(() => {
  if (selectedCategory.value === 'all') return groupedItems.value
  return groupedItems.value.filter(g => g.category.id === selectedCategory.value)
})

const totalItems = computed(() => {
  return Object.values(cart.value).reduce((sum, qty) => sum + qty, 0)
})

const totalPrice = computed(() => {
  return Object.entries(cart.value).reduce((sum, [key, qty]) => {
    const [itemId, variantId] = key.split('_')
    const item = props.menuItems.find(i => i.id === itemId)
    
    if (!item) return sum
    
    if (variantId) {
      const variant = item.variants?.find(v => v.id === variantId)
      return sum + (variant?.price || 0) * qty
    }
    
    return sum + (item.price || 0) * qty
  }, 0)
})

// Methods
const formatPrice = (price) => {
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND'
  }).format(price)
}

const handleCancel = () => {
  if (totalItems.value > 0) {
    if (confirm('Bạn có chắc muốn hủy order này?')) {
      cart.value = {}
      itemSugarLevels.value = {}
      emit('update:modelValue', false)
    }
  } else {
    emit('update:modelValue', false)
  }
}

const handleConfirm = () => {
  // Convert cart object to array format for API
  const cartArray = []
  
  Object.entries(cart.value).forEach(([key, quantity]) => {
    const [itemId, variantId] = key.split('_')
    const item = props.menuItems.find(i => i.id === itemId)
    
    if (item) {
      const sugarLevel = itemSugarLevels.value[key] || 100
      const note = sugarLevel === 100 ? '' : `${sugarLevel}%` // Only add note if not 100%
      
      if (variantId) {
        const variant = item.variants?.find(v => v.id === variantId)
        if (variant) {
          cartArray.push({
            menu_item_id: itemId,
            variant_id: variantId,
            name: item.name,
            variant_name: variant.name,
            price: variant.price,
            quantity: quantity,
            note: note
          })
        }
      } else {
        cartArray.push({
          menu_item_id: itemId,
          name: item.name,
          price: item.price,
          quantity: quantity,
          note: note
        })
      }
    }
  })
  
  emit('confirm', cartArray)
  cart.value = {}
  itemSugarLevels.value = {}
}

// Reset cart when modal closes
const resetCart = () => {
  cart.value = {}
  itemSugarLevels.value = {}
  selectedCategory.value = 'all'
}

defineExpose({ resetCart })
</script>

<style scoped>
.slide-up-enter-active,
.slide-up-leave-active {
  transition: transform 0.3s ease-out;
}

.slide-up-enter-from {
  transform: translateY(100%);
}

.slide-up-leave-to {
  transform: translateY(100%);
}

.scrollbar-hide::-webkit-scrollbar {
  display: none;
}

.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
