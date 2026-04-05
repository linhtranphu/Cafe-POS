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

    <!-- 2-column layout: Category sidebar + Menu items -->
    <div class="flex flex-1 overflow-hidden min-h-0">

      <!-- Left: Category Sidebar -->
      <div class="w-20 bg-white border-r flex-shrink-0 overflow-hidden">
        <button
          @click="selectedCategory = 'all'"
          :class="[
            'w-full py-3 px-1 flex flex-col items-center gap-1 transition-all touch-manipulation border-b border-gray-100',
            selectedCategory === 'all'
              ? 'bg-blue-50 border-l-4 border-l-blue-500'
              : 'active:bg-gray-100'
          ]">
          <span :class="['text-xs font-semibold leading-tight text-center', selectedCategory === 'all' ? 'text-blue-600' : 'text-gray-600']">Tất cả</span>
        </button>
        <button
          v-for="cat in categories.filter(c => c.id !== 'all')" :key="cat.id"
          @click="selectedCategory = cat.id"
          :class="[
            'w-full py-3 px-1 flex flex-col items-center gap-1 transition-all touch-manipulation border-b border-gray-100',
            selectedCategory === cat.id
              ? 'bg-blue-50 border-l-4 border-l-blue-500'
              : 'active:bg-gray-100'
          ]">
          <span :class="['text-xs font-semibold leading-tight text-center', selectedCategory === cat.id ? 'text-blue-600' : 'text-gray-600']">{{ cat.name }}</span>
        </button>
      </div>

      <!-- Right: Menu Items (single column list) -->
      <div class="flex-1 overflow-y-auto bg-gray-50" :style="{ paddingBottom: '100px' }">

        <!-- Filter Selected Items Button -->
        <div v-if="totalItems > 0" class="px-3 pt-2">
          <button
            @click="showSelectedOnly = !showSelectedOnly"
            :class="[
              'w-full py-2 px-3 rounded-xl text-xs font-semibold transition-all touch-manipulation active:scale-95',
              showSelectedOnly
                ? 'bg-orange-500 text-white shadow-md'
                : 'bg-orange-100 text-orange-700 active:bg-orange-200'
            ]">
            {{ showSelectedOnly ? '✓ Đang xem món đã chọn' : '🛒 Xem món đã chọn' }} ({{ totalItems }})
          </button>
        </div>

        <!-- Empty State -->
        <div v-if="filteredGroupedItems.length === 0" class="text-center py-20">
          <div class="text-6xl mb-3">🍽️</div>
          <p class="text-gray-500 font-medium">Không có món nào</p>
        </div>

        <!-- Category Groups -->
        <div v-for="group in filteredGroupedItems" :key="group.category.id" class="mb-4">
          <!-- Category Header (only show when viewing all) -->
          <div v-if="selectedCategory === 'all'"
            class="sticky top-0 z-10 bg-gray-100 px-3 py-2 border-b border-gray-200 flex items-center gap-2">
            <span class="text-xl">{{ group.category.icon }}</span>
            <span class="font-bold text-gray-700 text-sm">{{ group.category.name }}</span>
            <span class="text-xs text-gray-400">({{ group.items.length }})</span>
          </div>

          <!-- Item Rows -->
          <div class="px-3 py-2 space-y-2">
            <template v-for="item in group.items" :key="item.id">

              <!-- Single-size item row -->
              <div v-if="!item.has_variants"
                :class="[
                  'rounded-xl px-4 py-3 shadow-sm border transition-all',
                  getItemQty(item.id) > 0
                    ? 'border-orange-400 bg-orange-50 shadow-md'
                    : 'bg-white border-gray-200'
                ]">
                <div class="flex items-center gap-3">
                  <div class="flex-1 min-w-0">
                    <p class="font-bold text-gray-900 text-sm leading-tight">{{ item.name }}</p>
                    <p class="text-xs font-bold text-blue-600 mt-0.5">{{ formatPrice(item.price) }}</p>
                  </div>
                  <div class="flex items-center gap-2 flex-shrink-0">
                    <button v-if="getItemQty(item.id) > 0"
                      @click="removeItem(item.id)"
                      class="w-8 h-8 bg-white border-2 border-red-400 text-red-500 rounded-lg text-lg font-bold active:bg-red-50 active:scale-95 transition-all flex items-center justify-center">
                      −
                    </button>
                    <span v-if="getItemQty(item.id) > 0"
                      class="w-6 text-center font-bold text-blue-600 text-sm">
                      {{ getItemQty(item.id) }}
                    </span>
                    <button @click="addItem(item.id)"
                      class="w-8 h-8 bg-blue-500 text-white rounded-lg text-lg font-bold active:bg-blue-600 active:scale-95 transition-all flex items-center justify-center">
                      +
                    </button>
                    <button v-if="getItemQty(item.id) > 0"
                      @click="toggleNoteInput(item.id)"
                      :class="[
                        'w-8 h-8 rounded-lg transition-all active:scale-95 flex items-center justify-center text-sm',
                        getItemNote(item.id) || showNoteInput[item.id]
                          ? 'bg-amber-500 text-white'
                          : 'bg-gray-200 text-gray-600'
                      ]">
                      📝
                    </button>
                  </div>
                </div>
                <!-- Note input -->
                <div v-if="showNoteInput[item.id]" class="mt-2 pt-2 border-t border-gray-200">
                  <textarea v-model="itemNotes[item.id]"
                    placeholder="Ghi chú cho món này..."
                    rows="2"
                    class="w-full px-3 py-2 rounded-lg border-2 border-amber-300 bg-amber-50 text-sm focus:ring-2 focus:ring-amber-500 focus:border-amber-500 resize-none"></textarea>
                </div>
              </div>

              <!-- Multi-size item row -->
              <div v-else
                :class="[
                  'rounded-xl px-4 py-3 shadow-sm border transition-all',
                  getItemQtyWithVariants(item.id) > 0
                    ? 'border-orange-400 bg-orange-50 shadow-md'
                    : 'bg-white border-gray-200'
                ]">
                <p class="font-bold text-gray-900 text-sm mb-2">{{ item.name }}</p>
                <div class="space-y-2">
                  <div v-for="variant in item.variants" :key="variant.id">
                    <div class="flex items-center gap-3">
                      <div class="flex-1 min-w-0">
                        <span class="text-sm text-gray-700 font-medium">{{ variant.name }}</span>
                        <span class="text-xs font-bold text-blue-600 ml-2">{{ formatPrice(variant.price) }}</span>
                      </div>
                      <div class="flex items-center gap-2 flex-shrink-0">
                        <button v-if="getItemQty(item.id, variant.id) > 0"
                          @click="removeItem(item.id, variant.id)"
                          class="w-8 h-8 bg-white border-2 border-red-400 text-red-500 rounded-lg text-lg font-bold active:bg-red-50 active:scale-95 transition-all flex items-center justify-center">
                          −
                        </button>
                        <span v-if="getItemQty(item.id, variant.id) > 0"
                          class="w-6 text-center font-bold text-blue-600 text-sm">
                          {{ getItemQty(item.id, variant.id) }}
                        </span>
                        <button @click="addItem(item.id, variant.id)"
                          class="w-8 h-8 bg-blue-500 text-white rounded-lg text-lg font-bold active:bg-blue-600 active:scale-95 transition-all flex items-center justify-center">
                          +
                        </button>
                        <button v-if="getItemQty(item.id, variant.id) > 0"
                          @click="toggleNoteInput(item.id, variant.id)"
                          :class="[
                            'w-8 h-8 rounded-lg transition-all active:scale-95 flex items-center justify-center text-sm',
                            getItemNote(item.id, variant.id) || showNoteInput[`${item.id}_${variant.id}`]
                              ? 'bg-amber-500 text-white'
                              : 'bg-gray-200 text-gray-600'
                          ]">
                          📝
                        </button>
                      </div>
                    </div>
                    <!-- Note input for variant -->
                    <div v-if="showNoteInput[`${item.id}_${variant.id}`]" class="mt-2 pt-2 border-t border-gray-200">
                      <textarea v-model="itemNotes[`${item.id}_${variant.id}`]"
                        placeholder="Ghi chú cho size này..."
                        rows="2"
                        class="w-full px-3 py-2 rounded-lg border-2 border-amber-300 bg-amber-50 text-sm focus:ring-2 focus:ring-amber-500 focus:border-amber-500 resize-none"></textarea>
                    </div>
                  </div>
                </div>
              </div>

            </template>
          </div>
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
import { menuCategoryService } from '../services/menuCategory'
import BottomNav from '../components/BottomNav.vue'

const router = useRouter()
const orderStore = useOrderStore()
const shiftStore = useShiftStore()
const menuStore = useMenuStore()

// State
const cart = ref({})
const selectedCategory = ref('all')
const categoryOrder = ref([]) // ordered list of category names from API
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
  const allCategory = { id: 'all', name: 'Tất cả' }

  if (!menuItems.value || menuItems.value.length === 0) {
    return [allCategory]
  }

  const uniqueNames = new Set(menuItems.value.map(item => item.category))

  let ordered
  if (categoryOrder.value.length > 0) {
    // Sort by manager-defined order, append any unknown categories at the end
    const known = categoryOrder.value.filter(name => uniqueNames.has(name))
    const unknown = [...uniqueNames].filter(name => !categoryOrder.value.includes(name))
    ordered = [...known, ...unknown]
  } else {
    ordered = [...uniqueNames]
  }

  return [allCategory, ...ordered.map(name => ({ id: name, name }))]
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
    menuStore.fetchMenuItems(),
    (async () => {
      try {
        const res = await menuCategoryService.getWaiterCategories()
        const cats = Array.isArray(res) ? res : (res?.data || [])
        categoryOrder.value = cats.map(c => c.name)
      } catch (e) {
        // fallback: no predefined order
      }
    })()
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
