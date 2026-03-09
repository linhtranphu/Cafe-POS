<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Header -->
    <div class="bg-gradient-to-r from-blue-500 to-blue-600 text-white px-4 py-3 shadow-lg flex-shrink-0"
      style="padding-top: max(1rem, env(safe-area-inset-top))">
      <div class="flex items-center justify-between gap-3">
        <h2 class="text-xl font-bold">🍽️ Tạo Order</h2>
        <!-- Total Summary (when has items) -->
        <button
          v-if="totalItems > 0"
          @click="showCustomerNameModal = true"
          class="flex items-center gap-2 px-4 py-2 bg-green-500 rounded-xl shadow-lg active:bg-green-600 transition-all touch-manipulation active:scale-95">
          <div class="text-left">
            <p class="text-xs font-medium opacity-90">{{ totalItems }} món</p>
            <p class="text-sm font-bold">{{ formatPrice(totalPrice) }}</p>
          </div>
          <span class="text-lg">→</span>
        </button>
      </div>
    </div>

    <!-- Category Filter Buttons -->
    <div class="bg-white border-b px-4 py-2 flex-shrink-0">
      <div class="flex flex-wrap gap-2 mb-2">
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
          v-for="cat in categories.filter(c => c.id !== 'all')" :key="cat.id"
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
      
      <!-- Filter Selected Items Button -->
      <button
        v-if="totalItems > 0"
        @click="showSelectedOnly = !showSelectedOnly"
        :class="[
          'w-full py-2 px-4 rounded-xl text-sm font-semibold transition-all touch-manipulation active:scale-95',
          showSelectedOnly
            ? 'bg-orange-500 text-white shadow-md'
            : 'bg-orange-100 text-orange-700 active:bg-orange-200'
        ]">
        {{ showSelectedOnly ? '✓ Đang xem món đã chọn' : '🛒 Xem món đã chọn' }} ({{ totalItems }})
      </button>
    </div>

    <!-- Menu Items -->
    <div class="flex-1 overflow-y-auto bg-gray-50" 
      :style="{ paddingBottom: '100px' }">
      
      <!-- Empty State -->
      <div v-if="filteredGroupedItems.length === 0" class="text-center py-20">
        <div class="text-6xl mb-3">🍽️</div>
        <p class="text-gray-500 font-medium">Không có món nào</p>
      </div>

      <!-- Category Groups -->
      <div v-for="group in filteredGroupedItems" :key="group.category.id" class="mb-6">
        <!-- Category Header -->
        <div class="sticky top-0 z-10 bg-gradient-to-b from-gray-50 to-gray-100 px-4 py-3 border-b-2 border-gray-200">
          <div class="flex items-center gap-2">
            <span class="text-3xl">{{ group.category.icon }}</span>
            <span class="font-bold text-gray-800 text-xl">{{ group.category.name }}</span>
            <span class="text-lg text-gray-500">({{ group.items.length }})</span>
          </div>
        </div>

        <!-- Items Grid -->
        <div class="grid grid-cols-2 gap-3 px-4 py-3">
          <template v-for="item in group.items" :key="item.id">
            <!-- Single-size item -->
            <div v-if="!item.has_variants"
              :class="[
                'rounded-2xl p-4 shadow-sm transition-all flex flex-col',
                getItemQty(item.id) > 0 
                  ? 'border-2 border-orange-400 bg-orange-100 shadow-md' 
                  : 'bg-white border border-gray-200'
              ]">
              <div class="flex-1 mb-3">
                <h3 class="font-bold text-gray-900 mb-2 line-clamp-2 min-h-[2.5rem] leading-tight text-lg">
                  {{ item.name }}
                </h3>
                <p class="text-xs font-bold text-blue-600">
                  {{ formatPrice(item.price) }}
                </p>
              </div>

              <div class="mt-auto">
                <button v-if="getItemQty(item.id) === 0"
                  @click="addItem(item.id)"
                  class="w-full py-2.5 bg-gradient-to-r from-blue-500 to-blue-600 text-white font-bold rounded-xl active:from-blue-600 active:to-blue-700 active:scale-95 transition-all shadow-sm text-2xl">
                  +
                </button>

                <div v-else class="space-y-2">
                  <!-- +/- buttons right under price -->
                  <div class="flex gap-2">
                    <button @click="addItem(item.id)"
                      class="flex-1 py-2 bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-lg text-xl font-bold active:from-blue-600 active:to-blue-700 active:scale-95 transition-all shadow-md">
                      +
                    </button>
                    <button @click="removeItem(item.id)"
                      class="flex-1 py-2 bg-white border-2 border-red-500 text-red-600 rounded-lg text-xl font-bold active:bg-red-50 active:scale-95 transition-all shadow-sm">
                      −
                    </button>
                  </div>
                  
                  <!-- Quantity display with note icon -->
                  <div class="flex items-center gap-2">
                    <div class="flex-1 text-center py-1 bg-blue-50 rounded-lg">
                      <span class="text-base font-bold text-blue-600">
                        Số lượng: {{ getItemQty(item.id) }}
                      </span>
                    </div>
                    <button @click="toggleNoteInput(item.id)"
                      :class="[
                        'p-2 rounded-lg transition-all active:scale-95',
                        getItemNote(item.id) || showNoteInput[item.id]
                          ? 'bg-amber-500 text-white'
                          : 'bg-gray-200 text-gray-600'
                      ]">
                      📝
                    </button>
                  </div>
                  
                  <!-- Note input -->
                  <div v-if="showNoteInput[item.id]" class="pt-2 border-t border-gray-200">
                    <textarea v-model="itemNotes[item.id]"
                      placeholder="Ghi chú cho món này..."
                      rows="2"
                      class="w-full px-3 py-2 rounded-lg border-2 border-amber-300 bg-amber-50 text-sm focus:ring-2 focus:ring-amber-500 focus:border-amber-500 resize-none"></textarea>
                  </div>
                </div>
              </div>
            </div>

            <!-- Multi-size item -->
            <div v-else
              :class="[
                'rounded-2xl p-4 shadow-sm transition-all flex flex-col',
                getItemQtyWithVariants(item.id) > 0 
                  ? 'border-2 border-orange-400 bg-orange-100 shadow-md' 
                  : 'bg-white border border-gray-200'
              ]">
              <div class="mb-3">
                <h3 class="font-bold text-gray-900 mb-2 line-clamp-2 min-h-[2.5rem] leading-tight text-lg">
                  {{ item.name }}
                </h3>
              </div>

              <div class="space-y-3 flex-1">
                <div v-for="variant in item.variants" :key="variant.id" class="space-y-2">
                  <div class="flex items-center justify-between gap-2">
                    <div class="flex-1">
                      <p class="text-base font-semibold text-gray-700">{{ variant.name }}</p>
                      <p class="text-xs font-bold text-blue-600">{{ formatPrice(variant.price) }}</p>
                    </div>
                    
                    <div v-if="getItemQty(item.id, variant.id) === 0">
                      <button @click="addItem(item.id, variant.id)"
                        class="px-4 py-2 bg-blue-500 text-white text-xl font-bold rounded-lg active:bg-blue-600 active:scale-95 transition-all">
                        +
                      </button>
                    </div>
                    <div v-else class="flex items-center gap-2">
                      <span class="text-base font-bold text-blue-600">
                        {{ getItemQty(item.id, variant.id) }}
                      </span>
                      <button @click="toggleNoteInput(item.id, variant.id)"
                        :class="[
                          'p-1.5 rounded-lg transition-all active:scale-95 text-sm',
                          getItemNote(item.id, variant.id) || showNoteInput[`${item.id}_${variant.id}`]
                            ? 'bg-amber-500 text-white'
                            : 'bg-gray-200 text-gray-600'
                        ]">
                        📝
                      </button>
                    </div>
                  </div>
                  
                  <!-- Note input for variant -->
                  <div v-if="showNoteInput[`${item.id}_${variant.id}`]" class="pt-2 border-t border-gray-200">
                    <textarea v-model="itemNotes[`${item.id}_${variant.id}`]"
                      placeholder="Ghi chú cho size này..."
                      rows="2"
                      class="w-full px-3 py-2 rounded-lg border-2 border-amber-300 bg-amber-50 text-sm focus:ring-2 focus:ring-amber-500 focus:border-amber-500 resize-none"></textarea>
                  </div>
                  
                  <!-- +/- buttons for variant (only show when qty > 0) -->
                  <div v-if="getItemQty(item.id, variant.id) > 0" class="flex gap-2">
                    <button @click="addItem(item.id, variant.id)"
                      class="flex-1 py-2 bg-blue-500 text-white rounded-lg text-xl font-bold active:bg-blue-600 active:scale-95 transition-all">
                      +
                    </button>
                    <button @click="removeItem(item.id, variant.id)"
                      class="flex-1 py-2 bg-white border-2 border-red-500 text-red-600 rounded-lg text-xl font-bold active:bg-red-50 active:scale-95 transition-all">
                      −
                    </button>
                  </div>
                </div>
              </div>

              <div v-if="getItemQtyWithVariants(item.id) > 0" 
                class="mt-3 pt-3 border-t border-gray-200">
                <div class="text-center mb-2">
                  <span class="text-xs font-bold text-blue-600">
                    🛒 {{ getItemQtyWithVariants(item.id) }} món
                  </span>
                </div>
                
                <!-- Note input at bottom for multi-variant item -->
                <div v-if="showNoteInput[item.id]" class="mt-2">
                  <textarea v-model="itemNotes[item.id]"
                    placeholder="Ghi chú chung cho món này..."
                    rows="2"
                    class="w-full px-3 py-2 rounded-lg border-2 border-amber-300 bg-amber-50 text-sm focus:ring-2 focus:ring-amber-500 focus:border-amber-500 resize-none"></textarea>
                </div>
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />

    <!-- Customer Name Modal -->
    <transition name="slide-up">
      <div v-if="showCustomerNameModal" class="fixed inset-0 bg-black/60 z-[60] flex items-end">
        <div class="bg-white rounded-t-3xl w-full max-h-[90vh] overflow-y-auto flex flex-col">

          <!-- Customer Name Section (hero) -->
          <div :class="[
            'px-4 pt-5 pb-4 transition-all',
            customerName ? 'bg-green-50' : 'bg-amber-50'
          ]">
            <!-- Label -->
            <div class="flex items-center gap-2 mb-3">
              <span class="text-2xl">👤</span>
              <div>
                <p class="font-bold text-gray-800 text-base leading-tight">Tên khách hàng</p>
                <p v-if="!customerName" class="text-amber-600 text-xs font-medium">⚠️ Chưa điền — nhớ nhập để phân biệt order!</p>
                <p v-else class="text-green-600 text-xs font-medium">✓ Đã điền</p>
              </div>
            </div>

            <!-- Input -->
            <div :class="[
              'flex items-center gap-3 rounded-2xl px-4 py-3 border-2 transition-all',
              customerName
                ? 'bg-white border-green-400 shadow-md'
                : 'bg-white border-amber-400 shadow-lg customer-name-pulse'
            ]">
              <input
                v-model="customerName"
                ref="customerNameInput"
                type="text"
                placeholder="Nhập tên khách..."
                :class="[
                  'flex-1 text-lg font-semibold outline-none bg-transparent placeholder-gray-300',
                  customerName ? 'text-gray-900' : 'text-gray-900'
                ]"
              />
              <button v-if="customerName"
                @click="customerName = ''"
                class="w-7 h-7 rounded-full bg-gray-200 text-gray-500 active:bg-gray-300 flex items-center justify-center text-sm font-bold touch-manipulation flex-shrink-0">
                ×
              </button>
            </div>

          </div>

          <div class="px-4 py-3 flex-1">
            <!-- Order Review -->
            <div class="mb-3 bg-gray-50 rounded-xl p-3 border border-gray-200">
              <div class="flex items-center justify-between mb-2">
                <span class="text-sm font-bold text-gray-700">🛒 Chi tiết order</span>
                <span class="text-xs text-blue-600 font-semibold">{{ totalItems }} món</span>
              </div>

              <div class="space-y-1.5 max-h-44 overflow-y-auto">
                <div v-for="item in orderItems" :key="item.key"
                  class="bg-white rounded-lg p-2">
                  <div class="flex justify-between items-start gap-2">
                    <div class="flex-1 min-w-0">
                      <p class="font-semibold text-gray-900 text-sm truncate">{{ item.name }}</p>
                      <p v-if="item.variant_name" class="text-xs text-gray-500">{{ item.variant_name }}</p>
                      <p v-if="getItemNote(item.menu_item_id, item.variant_id)" class="text-xs text-amber-700 mt-0.5">
                        📝 {{ getItemNote(item.menu_item_id, item.variant_id) }}
                      </p>
                    </div>
                    <div class="text-right flex-shrink-0">
                      <p class="text-xs font-bold text-blue-600">x{{ item.quantity }}</p>
                      <p class="text-xs font-bold text-gray-900">{{ formatPrice(item.price * item.quantity) }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Total -->
            <div class="mb-4 px-4 py-3 bg-green-50 rounded-xl border border-green-200 flex justify-between items-center">
              <span class="font-bold text-green-900">💰 Tổng tiền</span>
              <span class="text-xl font-bold text-green-600">{{ formatPrice(totalPrice) }}</span>
            </div>

            <!-- Action Buttons -->
            <div class="flex gap-2">
              <button @click="showCustomerNameModal = false"
                :disabled="isCreatingOrder"
                class="flex-1 bg-gray-100 text-gray-600 py-3 rounded-xl font-semibold text-sm touch-manipulation active:scale-98 transition-all disabled:opacity-50">
                Quay lại
              </button>
              <button @click="finalizeOrder"
                :disabled="isCreatingOrder"
                :class="[
                  'flex-2 py-3 rounded-xl font-bold text-sm shadow-lg touch-manipulation active:scale-98 transition-all disabled:opacity-50 text-white px-6',
                  customerName ? 'bg-blue-500 active:bg-blue-600' : 'bg-blue-400 active:bg-blue-500'
                ]">
                {{ isCreatingOrder ? 'Đang tạo...' : (customerName ? `✓ Tạo Order cho ${customerName}` : '✓ Tạo Order (không tên)') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useOrderStore } from '../stores/order'
import { useShiftStore } from '../stores/shift'
import { useMenuStore } from '../stores/menu'
import BottomNav from '../components/BottomNav.vue'

const router = useRouter()
const orderStore = useOrderStore()
const shiftStore = useShiftStore()
const menuStore = useMenuStore()

// State
const cart = ref({})
const selectedCategory = ref('all')
const showSelectedOnly = ref(false)
const itemNotes = ref({})
const showNoteInput = ref({})
const customerName = ref('')
const orderNote = ref('')
const showCustomerNameModal = ref(false)
const isCreatingOrder = ref(false)
const customerNameInput = ref(null)

// Auto-focus customer name input when modal opens
watch(showCustomerNameModal, async (val) => {
  if (val) {
    await nextTick()
    customerNameInput.value?.focus()
  }
})

// Sugar levels - DISABLED (bỏ mục Đường)
const sugarLevels = [
  // { value: 25, label: '25%' },
  // { value: 50, label: '50%' },
  // { value: 100, label: '100%' }
]

// Computed
const menuItems = computed(() => menuStore.items)

const categories = computed(() => {
  const allCategory = { id: 'all', name: 'Tất cả', icon: '📋' }
  
  if (!menuItems.value || menuItems.value.length === 0) {
    return [allCategory]
  }
  
  const uniqueCategories = [...new Set(menuItems.value.map(item => item.category))]
  
  const categoryIcons = {
    'Cà phê': '☕',
    'Trà': '🍵',
    'Nước ép': '🧃',
    'Bánh ngọt': '🍰',
    'Món nhẹ': '🍴',
    'Sinh tố': '🥤',
    'Đồ uống khác': '🥛'
  }
  
  const menuCategories = uniqueCategories.map(cat => ({
    id: cat,
    name: cat,
    icon: categoryIcons[cat] || '🍽️'
  }))
  
  return [allCategory, ...menuCategories]
})

const groupedItems = computed(() => {
  const groups = {}
  
  menuItems.value.forEach(item => {
    const categoryId = item.category
    if (!groups[categoryId]) {
      const category = categories.value.find(c => c.id === categoryId)
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
  let groups = selectedCategory.value === 'all' ? groupedItems.value : groupedItems.value.filter(g => g.category.id === selectedCategory.value)
  
  // Filter by selected items if showSelectedOnly is true
  if (showSelectedOnly.value) {
    groups = groups.map(group => ({
      ...group,
      items: group.items.filter(item => {
        if (!item.has_variants) {
          return getItemQty(item.id) > 0
        } else {
          return getItemQtyWithVariants(item.id) > 0
        }
      })
    })).filter(group => group.items.length > 0)
  }
  
  return groups
})

const totalItems = computed(() => {
  return Object.values(cart.value).reduce((sum, qty) => sum + qty, 0)
})

const totalPrice = computed(() => {
  return Object.entries(cart.value).reduce((sum, [key, qty]) => {
    const [itemId, variantId] = key.split('_')
    const item = menuItems.value.find(i => i.id === itemId)
    
    if (!item) return sum
    
    if (variantId) {
      const variant = item.variants?.find(v => v.id === variantId)
      return sum + (variant?.price || 0) * qty
    }
    
    return sum + (item.price || 0) * qty
  }, 0)
})

const orderItems = computed(() => {
  const items = []
  
  Object.entries(cart.value).forEach(([key, quantity]) => {
    const [itemId, variantId] = key.split('_')
    const item = menuItems.value.find(i => i.id === itemId)
    
    if (!item) return
    
    if (variantId) {
      const variant = item.variants?.find(v => v.id === variantId)
      if (variant) {
        items.push({
          key,
          menu_item_id: itemId,
          variant_id: variantId,
          name: item.name,
          variant_name: variant.name,
          price: variant.price,
          quantity
        })
      }
    } else {
      items.push({
        key,
        menu_item_id: itemId,
        variant_id: null,
        name: item.name,
        variant_name: null,
        price: item.price,
        quantity
      })
    }
  })
  
  return items
})

// Methods
const getItemQty = (itemId, variantId = null) => {
  const key = variantId ? `${itemId}_${variantId}` : itemId
  return cart.value[key] || 0
}

const addItem = (itemId, variantId = null) => {
  const key = variantId ? `${itemId}_${variantId}` : itemId
  cart.value[key] = (cart.value[key] || 0) + 1
}

const removeItem = (itemId, variantId = null) => {
  const key = variantId ? `${itemId}_${variantId}` : itemId
  if (cart.value[key] > 1) {
    cart.value[key]--
  } else {
    delete cart.value[key]
    delete itemNotes.value[key]
    delete showNoteInput.value[key]
  }
}

const getItemQtyWithVariants = (itemId) => {
  return Object.keys(cart.value)
    .filter(key => key.startsWith(itemId))
    .reduce((sum, key) => sum + cart.value[key], 0)
}

const toggleNoteInput = (itemId, variantId = null) => {
  const key = variantId ? `${itemId}_${variantId}` : itemId
  showNoteInput.value[key] = !showNoteInput.value[key]
}

const getItemNote = (itemId, variantId = null) => {
  const key = variantId ? `${itemId}_${variantId}` : itemId
  return itemNotes.value[key] || ''
}

const setItemNote = (itemId, note, variantId = null) => {
  const key = variantId ? `${itemId}_${variantId}` : itemId
  itemNotes.value[key] = note
}

const formatPrice = (price) => {
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND'
  }).format(price)
}

const resetCart = () => {
  if (confirm('Bạn có chắc muốn xóa hết giỏ hàng?')) {
    cart.value = {}
    itemNotes.value = {}
    selectedCategory.value = 'all'
  }
}

const finalizeOrder = async () => {
  if (isCreatingOrder.value) return
  
  if (!shiftStore.hasOpenShift) {
    alert('Vui lòng mở ca làm việc trước')
    router.push('/shifts')
    return
  }
  
  try {
    isCreatingOrder.value = true
    
    // Convert cart to order items
    const cartArray = []
    Object.entries(cart.value).forEach(([key, quantity]) => {
      const [itemId, variantId] = key.split('_')
      const item = menuItems.value.find(i => i.id === itemId)
      
      if (item) {
        const itemNote = itemNotes.value[key] || ''
        const note = itemNote
        
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
    
    const orderData = {
      customer_name: customerName.value || '',
      items: cartArray,
      note: orderNote.value || '',
      shift_id: shiftStore.currentShift.id
    }
    
    await orderStore.createOrder(orderData)
    
    // Reset and navigate
    cart.value = {}
    itemNotes.value = {}
    showNoteInput.value = {}
    customerName.value = ''
    orderNote.value = ''
    selectedCategory.value = 'all'
    showCustomerNameModal.value = false
    
    alert('✅ Tạo order thành công!')
    router.push('/orders')
  } catch (error) {
    alert('❌ Lỗi: ' + (error.response?.data?.error || error.message))
  } finally {
    isCreatingOrder.value = false
  }
}

// Lifecycle
onMounted(async () => {
  await Promise.all([
    shiftStore.fetchCurrentShift(),
    menuStore.fetchMenuItems()
  ])
  
  if (!shiftStore.hasOpenShift) {
    alert('Vui lòng mở ca làm việc trước')
    router.push('/shifts')
  }
})
</script>

<style scoped>
.touch-manipulation {
  touch-action: manipulation;
  -webkit-tap-highlight-color: transparent;
}

.scrollbar-hide::-webkit-scrollbar {
  display: none;
}

.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.active\:scale-95:active {
  transform: scale(0.95);
}

.active\:scale-98:active {
  transform: scale(0.98);
}

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

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.25rem;
}

@keyframes pulse-border {
  0%, 100% { border-color: #f59e0b; box-shadow: 0 0 0 0 rgba(245, 158, 11, 0.4); }
  50% { border-color: #d97706; box-shadow: 0 0 0 6px rgba(245, 158, 11, 0); }
}

.customer-name-pulse {
  animation: pulse-border 2s ease-in-out infinite;
}
</style>
