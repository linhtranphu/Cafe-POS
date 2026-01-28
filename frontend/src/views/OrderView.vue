<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm">
      <div class="px-4 py-3">
        <div class="flex items-center justify-between mb-3">
          <h1 class="text-xl font-bold text-gray-800">📋 Orders</h1>
          <div class="flex gap-2">
            <button @click="refreshOrders" class="p-2 rounded-lg bg-gray-100 hover:bg-gray-200">
              🔄
            </button>
          </div>
        </div>
        
        <!-- Status Filter Pills -->
        <div class="flex gap-2 overflow-x-auto pb-2 scrollbar-hide">
          <button v-for="status in statuses" :key="status.value" 
            @click="filterStatus = status.value"
            :class="[
              'px-4 py-2 rounded-full text-sm font-medium whitespace-nowrap transition-all',
              filterStatus === status.value 
                ? 'bg-blue-500 text-white shadow-md' 
                : 'bg-gray-100 text-gray-700'
            ]">
            {{ status.icon }} {{ status.label }} 
            <span class="ml-1 text-xs opacity-75">({{ getOrderCountByStatus(status.value) }})</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Shift Warning -->
    <div v-if="!hasOpenShift" class="mx-4 mt-4 bg-yellow-50 border-l-4 border-yellow-400 p-4 rounded-lg">
      <p class="text-yellow-700 text-sm font-medium">⚠️ Chưa mở ca làm việc</p>
      <button @click="$router.push('/shifts')" class="mt-2 bg-yellow-500 text-white px-4 py-2 rounded-lg text-sm font-medium">
        Mở ca ngay
      </button>
    </div>

    <!-- Orders List -->
    <div class="px-4 py-4 pb-24">
      <div v-if="loading" class="text-center py-10">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
      </div>
      
      <div v-else-if="filteredOrders.length === 0" class="text-center py-16">
        <div class="text-6xl mb-4">📭</div>
        <p class="text-gray-500">Không có order nào</p>
      </div>
      
      <div v-else class="space-y-3">
        <div v-for="order in filteredOrders" :key="order.id" 
          @click="viewOrderDetail(order)"
          class="bg-white rounded-2xl p-4 shadow-sm active:scale-98 transition-transform">
          
          <!-- Order Header -->
          <div class="flex justify-between items-start mb-3">
            <div>
              <h3 class="font-bold text-lg">{{ order.order_number }}</h3>
              <p class="text-sm text-gray-600">{{ order.customer_name || 'Khách lẻ' }}</p>
              <p class="text-xs text-gray-400">{{ formatTime(order.created_at) }}</p>
            </div>
            <span :class="getStatusBadge(order.status)" class="px-3 py-1 rounded-full text-xs font-medium">
              {{ getStatusText(order.status) }}
            </span>
          </div>

          <!-- Items Summary -->
          <div class="mb-3 space-y-1">
            <div v-for="(item, idx) in order.items.slice(0, 2)" :key="idx" 
              class="flex justify-between text-sm">
              <span class="text-gray-700">{{ item.name }} <span class="text-gray-400">x{{ item.quantity }}</span></span>
              <span class="font-medium text-gray-900">{{ formatPrice(item.subtotal) }}</span>
            </div>
            <p v-if="order.items.length > 2" class="text-xs text-gray-400">
              +{{ order.items.length - 2 }} món khác...
            </p>
          </div>

          <!-- Total -->
          <div class="flex justify-between items-center pt-3 border-t">
            <span class="text-sm font-medium text-gray-600">Tổng cộng</span>
            <span class="text-lg font-bold text-green-600">{{ formatPrice(order.total) }}</span>
          </div>

          <!-- Quick Actions -->
          <div class="mt-3 flex gap-2">
            <button v-if="order.status === 'CREATED'" 
              @click.stop="quickPayment(order)"
              class="flex-1 bg-green-500 text-white py-2 rounded-lg text-sm font-medium active:bg-green-600">
              💰 Thu tiền
            </button>
            <button v-if="order.status === 'PAID' && order.amount_due <= 0" 
              @click.stop="sendToBar(order.id)"
              class="flex-1 bg-blue-500 text-white py-2 rounded-lg text-sm font-medium active:bg-blue-600">
              🍹 Gửi bar
            </button>
            <button v-if="order.status === 'READY'" 
              @click.stop="serveOrder(order.id)"
              class="flex-1 bg-purple-500 text-white py-2 rounded-lg text-sm font-medium active:bg-purple-600">
              🎉 Giao khách
            </button>
            <button v-if="order.status === 'QUEUED' || order.status === 'IN_PROGRESS'" 
              class="flex-1 bg-gray-300 text-gray-600 py-2 rounded-lg text-sm font-medium cursor-not-allowed">
              ⏳ Đang pha...
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Floating Action Button -->
    <button v-if="hasOpenShift" 
      @click="startNewOrder"
      class="fixed bottom-20 right-4 w-16 h-16 bg-blue-500 text-white rounded-full shadow-lg flex items-center justify-center text-2xl active:scale-95 transition-transform z-30">
      ➕
    </button>

    <!-- Bottom Navigation -->
    <BottomNav />

    <!-- Create Order - Full Screen -->
    <transition name="slide-up">
      <div v-if="showCreateOrder" class="fixed inset-0 bg-white z-50 overflow-hidden flex flex-col">
        <!-- Header -->
        <div class="bg-blue-500 text-white px-4 py-4 flex items-center justify-between">
          <button @click="cancelCreateOrder" class="text-2xl">←</button>
          <h2 class="text-lg font-bold">Tạo Order Mới</h2>
          <button @click="confirmOrder" :disabled="cart.length === 0" 
            class="text-sm font-medium px-4 py-2 bg-white text-blue-500 rounded-lg disabled:opacity-50">
            Xác nhận
          </button>
        </div>

        <!-- Customer Name -->
        <div class="px-4 py-3 bg-gray-50 border-b">
          <input v-model="customerName" 
            type="text" 
            placeholder="Tên khách hàng (tùy chọn)"
            class="w-full px-4 py-3 rounded-lg border focus:ring-2 focus:ring-blue-500 focus:border-transparent">
        </div>

        <!-- Category Tabs -->
        <div class="flex gap-2 px-4 py-3 overflow-x-auto bg-white border-b scrollbar-hide">
          <button v-for="cat in categories" :key="cat.id"
            @click="selectedCategory = cat.id"
            :class="[
              'px-4 py-2 rounded-full text-sm font-medium whitespace-nowrap',
              selectedCategory === cat.id 
                ? 'bg-blue-500 text-white' 
                : 'bg-gray-100 text-gray-700'
            ]">
            {{ cat.icon }} {{ cat.name }}
          </button>
        </div>

        <!-- Menu Items Grid -->
        <div class="flex-1 overflow-y-auto px-4 py-4">
          <div class="grid grid-cols-2 gap-3">
            <button v-for="item in filteredMenuItems" :key="item.id"
              @click="addToCart(item)"
              class="bg-white rounded-xl p-4 shadow-sm active:scale-95 transition-transform text-left">
              <div class="font-medium text-gray-900 mb-1">{{ item.name }}</div>
              <div class="text-sm font-bold text-blue-600">{{ formatPrice(item.price) }}</div>
              <div v-if="getCartItemQty(item.id) > 0" 
                class="mt-2 bg-blue-500 text-white text-xs px-2 py-1 rounded-full inline-block">
                {{ getCartItemQty(item.id) }} món
              </div>
            </button>
          </div>
        </div>

        <!-- Cart Summary - Fixed Bottom -->
        <div v-if="cart.length > 0" class="bg-white border-t shadow-lg">
          <div class="px-4 py-3">
            <!-- Cart Items -->
            <div class="max-h-32 overflow-y-auto mb-3 space-y-2">
              <div v-for="(item, idx) in cart" :key="idx" 
                class="flex items-center gap-3 bg-gray-50 rounded-lg p-2">
                <span class="flex-1 text-sm font-medium">{{ item.name }}</span>
                <div class="flex items-center gap-2">
                  <button @click="decreaseQty(idx)" 
                    class="w-8 h-8 bg-gray-200 rounded-full text-lg font-bold active:bg-gray-300">
                    −
                  </button>
                  <span class="w-8 text-center font-bold">{{ item.quantity }}</span>
                  <button @click="increaseQty(idx)" 
                    class="w-8 h-8 bg-blue-500 text-white rounded-full text-lg font-bold active:bg-blue-600">
                    +
                  </button>
                </div>
                <button @click="removeFromCart(idx)" class="text-red-500 text-xl">×</button>
              </div>
            </div>
            
            <!-- Total -->
            <div class="flex justify-between items-center">
              <span class="text-gray-600">Tổng cộng</span>
              <span class="text-2xl font-bold text-green-600">{{ formatPrice(cartTotal) }}</span>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- Order Detail Modal -->
    <transition name="slide-up">
      <div v-if="selectedOrder" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full max-h-[85vh] overflow-y-auto">
          <div class="sticky top-0 bg-white px-4 py-4 border-b flex justify-between items-center">
            <h3 class="text-lg font-bold">Chi tiết Order</h3>
            <button @click="selectedOrder = null" class="text-2xl text-gray-400">×</button>
          </div>
          
          <div class="px-4 py-4">
            <!-- Order Info -->
            <div class="mb-4">
              <h4 class="text-2xl font-bold mb-1">{{ selectedOrder.order_number }}</h4>
              <p class="text-gray-600">{{ selectedOrder.customer_name || 'Khách lẻ' }}</p>
              <p class="text-sm text-gray-400">{{ formatDate(selectedOrder.created_at) }}</p>
              <span :class="getStatusBadge(selectedOrder.status)" 
                class="inline-block mt-2 px-3 py-1 rounded-full text-xs font-medium">
                {{ getStatusText(selectedOrder.status) }}
              </span>
            </div>

            <!-- Items -->
            <div class="mb-4">
              <h5 class="font-bold mb-2">Món đã order</h5>
              <div class="space-y-2">
                <div v-for="item in selectedOrder.items" :key="item.menu_item_id" 
                  class="flex justify-between bg-gray-50 p-3 rounded-lg">
                  <div>
                    <div class="font-medium">{{ item.name }}</div>
                    <div class="text-sm text-gray-500">{{ formatPrice(item.price) }} x {{ item.quantity }}</div>
                  </div>
                  <div class="font-bold">{{ formatPrice(item.subtotal) }}</div>
                </div>
              </div>
            </div>

            <!-- Note -->
            <div v-if="selectedOrder.note" class="mb-4 p-3 bg-yellow-50 rounded-lg">
              <p class="text-sm text-gray-700">📝 {{ selectedOrder.note }}</p>
            </div>

            <!-- Total -->
            <div class="mb-4 p-4 bg-gray-50 rounded-lg">
              <div class="flex justify-between text-lg font-bold">
                <span>Tổng cộng</span>
                <span class="text-green-600">{{ formatPrice(selectedOrder.total) }}</span>
              </div>
              <div v-if="selectedOrder.amount_paid > 0" class="flex justify-between text-sm text-gray-600 mt-1">
                <span>Đã thu</span>
                <span>{{ formatPrice(selectedOrder.amount_paid) }}</span>
              </div>
            </div>

            <!-- Actions -->
            <div class="space-y-2">
              <button v-if="selectedOrder.status === 'CREATED'" 
                @click="showPaymentModal(selectedOrder)"
                class="w-full bg-green-500 text-white py-3 rounded-xl font-medium active:bg-green-600">
                💰 Thu tiền
              </button>
              <button v-if="selectedOrder.status === 'CREATED'" 
                @click="editOrder(selectedOrder)"
                class="w-full bg-blue-500 text-white py-3 rounded-xl font-medium active:bg-blue-600">
                ✏️ Chỉnh sửa
              </button>
              <button v-if="selectedOrder.status === 'PAID' && selectedOrder.amount_due <= 0" 
                @click="sendToBar(selectedOrder.id)"
                class="w-full bg-blue-500 text-white py-3 rounded-xl font-medium active:bg-blue-600">
                🍹 Gửi quầy bar
              </button>
              <button v-if="selectedOrder.status === 'READY'" 
                @click="serveOrder(selectedOrder.id)"
                class="w-full bg-purple-500 text-white py-3 rounded-xl font-medium active:bg-purple-600">
                🎉 Giao cho khách
              </button>
              <div v-if="selectedOrder.status === 'QUEUED' || selectedOrder.status === 'IN_PROGRESS'" 
                class="w-full bg-gray-100 text-gray-600 py-3 rounded-xl font-medium text-center">
                ⏳ Barista đang pha chế...
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- Payment Modal -->
    <transition name="slide-up">
      <div v-if="showPayment" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6">
          <h3 class="text-xl font-bold mb-4">💰 Thu tiền</h3>
          
          <div class="mb-4 p-4 bg-gray-50 rounded-lg">
            <div class="flex justify-between mb-2">
              <span class="text-gray-600">Tổng tiền</span>
              <span class="text-xl font-bold text-green-600">{{ formatPrice(paymentOrder?.total) }}</span>
            </div>
            <div v-if="paymentOrder?.amount_paid > 0" class="flex justify-between text-sm text-gray-600">
              <span>Đã thu</span>
              <span>{{ formatPrice(paymentOrder?.amount_paid) }}</span>
            </div>
          </div>

          <div class="mb-4">
            <label class="block text-sm font-medium mb-2">Số tiền thu</label>
            <input v-model.number="paymentAmount" 
              type="number" 
              step="1000"
              class="w-full px-4 py-3 text-lg font-bold border rounded-lg focus:ring-2 focus:ring-green-500">
          </div>

          <div class="mb-4">
            <label class="block text-sm font-medium mb-2">Phương thức</label>
            <div class="grid grid-cols-3 gap-2">
              <button v-for="method in paymentMethods" :key="method.value"
                @click="paymentMethod = method.value"
                :class="[
                  'py-3 rounded-lg font-medium',
                  paymentMethod === method.value 
                    ? 'bg-green-500 text-white' 
                    : 'bg-gray-100 text-gray-700'
                ]">
                {{ method.icon }} {{ method.label }}
              </button>
            </div>
          </div>

          <div class="flex gap-2">
            <button @click="showPayment = false" 
              class="flex-1 bg-gray-200 text-gray-700 py-3 rounded-xl font-medium">
              Hủy
            </button>
            <button @click="processPayment" 
              class="flex-1 bg-green-500 text-white py-3 rounded-xl font-medium">
              Xác nhận
            </button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useOrderStore } from '../stores/order'
import { useShiftStore } from '../stores/shift'
import { useMenuStore } from '../stores/menu'
import { useRouter } from 'vue-router'
import BottomNav from '../components/BottomNav.vue'

const router = useRouter()
const orderStore = useOrderStore()
const shiftStore = useShiftStore()
const menuStore = useMenuStore()

// State
const filterStatus = ref('ALL')
const showCreateOrder = ref(false)
const selectedOrder = ref(null)
const showPayment = ref(false)
const paymentOrder = ref(null)
const paymentAmount = ref(0)
const paymentMethod = ref('CASH')

// Create Order State
const customerName = ref('')
const selectedCategory = ref('all')
const cart = ref([])

// Data
const statuses = [
  { value: 'ALL', label: 'Tất cả', icon: '📋' },
  { value: 'CREATED', label: 'Mới', icon: '🆕' },
  { value: 'PAID', label: 'Đã thu', icon: '💰' },
  { value: 'QUEUED', label: 'Chờ pha', icon: '⏳' },
  { value: 'IN_PROGRESS', label: 'Đang pha', icon: '🍹' },
  { value: 'READY', label: 'Sẵn sàng', icon: '✅' },
  { value: 'SERVED', label: 'Hoàn tất', icon: '🎉' }
]

const categories = [
  { id: 'all', name: 'Tất cả', icon: '📋' },
  { id: 'coffee', name: 'Cà phê', icon: '☕' },
  { id: 'tea', name: 'Trà', icon: '🍵' },
  { id: 'juice', name: 'Nước ép', icon: '🧃' },
  { id: 'food', name: 'Đồ ăn', icon: '🍰' }
]

const paymentMethods = [
  { value: 'CASH', label: 'Tiền mặt', icon: '💵' },
  { value: 'QR', label: 'QR', icon: '📱' },
  { value: 'TRANSFER', label: 'CK', icon: '🏦' }
]

// Computed
const loading = computed(() => orderStore.loading)
const orders = computed(() => orderStore.orders)
const menuItems = computed(() => menuStore.items)
const hasOpenShift = computed(() => shiftStore.hasOpenShift)

const filteredOrders = computed(() => {
  if (filterStatus.value === 'ALL') return orders.value
  return orders.value.filter(o => o.status === filterStatus.value)
})

const filteredMenuItems = computed(() => {
  if (selectedCategory.value === 'all') return menuItems.value
  return menuItems.value.filter(item => item.category === selectedCategory.value)
})

const cartTotal = computed(() => {
  return cart.value.reduce((sum, item) => sum + (item.price * item.quantity), 0)
})

// Methods
const refreshOrders = async () => {
  await orderStore.fetchOrders()
}

const getOrderCountByStatus = (status) => {
  if (status === 'ALL') return orders.value.length
  return orders.value.filter(o => o.status === status).length
}

const getStatusBadge = (status) => {
  const badges = {
    CREATED: 'bg-gray-100 text-gray-800',
    PAID: 'bg-green-100 text-green-800',
    QUEUED: 'bg-yellow-100 text-yellow-800',
    IN_PROGRESS: 'bg-blue-100 text-blue-800',
    READY: 'bg-purple-100 text-purple-800',
    SERVED: 'bg-green-100 text-green-800',
    CANCELLED: 'bg-red-100 text-red-800'
  }
  return badges[status] || 'bg-gray-100 text-gray-800'
}

const getStatusText = (status) => {
  const texts = {
    CREATED: 'Mới tạo',
    PAID: 'Đã thanh toán',
    QUEUED: 'Chờ pha chế',
    IN_PROGRESS: 'Đang pha chế',
    READY: 'Sẵn sàng',
    SERVED: 'Đã phục vụ',
    CANCELLED: 'Đã hủy'
  }
  return texts[status] || status
}

const formatPrice = (price) => {
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(price)
}

const formatTime = (date) => {
  const d = new Date(date)
  return d.toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' })
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('vi-VN')
}

const viewOrderDetail = (order) => {
  selectedOrder.value = order
}

const startNewOrder = () => {
  cart.value = []
  customerName.value = ''
  selectedCategory.value = 'all'
  showCreateOrder.value = true
}

const cancelCreateOrder = () => {
  if (cart.value.length > 0) {
    if (!confirm('Bạn có chắc muốn hủy order này?')) return
  }
  showCreateOrder.value = false
  cart.value = []
  customerName.value = ''
}

const addToCart = (item) => {
  const existing = cart.value.find(i => i.menu_item_id === item.id)
  if (existing) {
    existing.quantity++
  } else {
    cart.value.push({
      menu_item_id: item.id,
      name: item.name,
      price: item.price,
      quantity: 1
    })
  }
}

const getCartItemQty = (itemId) => {
  const item = cart.value.find(i => i.menu_item_id === itemId)
  return item ? item.quantity : 0
}

const increaseQty = (index) => {
  cart.value[index].quantity++
}

const decreaseQty = (index) => {
  if (cart.value[index].quantity > 1) {
    cart.value[index].quantity--
  } else {
    removeFromCart(index)
  }
}

const removeFromCart = (index) => {
  cart.value.splice(index, 1)
}

const confirmOrder = async () => {
  try {
    const orderData = {
      customer_name: customerName.value || '',
      items: cart.value,
      note: '',
      shift_id: shiftStore.currentShift.id
    }
    await orderStore.createOrder(orderData)
    showCreateOrder.value = false
    cart.value = []
    customerName.value = ''
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

const quickPayment = (order) => {
  paymentOrder.value = order
  paymentAmount.value = order.amount_due || order.total
  paymentMethod.value = 'CASH'
  showPayment.value = true
  selectedOrder.value = null
}

const showPaymentModal = (order) => {
  paymentOrder.value = order
  paymentAmount.value = order.amount_due || order.total
  paymentMethod.value = 'CASH'
  showPayment.value = true
  selectedOrder.value = null
}

const processPayment = async () => {
  try {
    await orderStore.collectPayment(paymentOrder.value.id, {
      amount: paymentAmount.value,
      payment_method: paymentMethod.value
    })
    showPayment.value = false
    paymentOrder.value = null
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const sendToBar = async (orderId) => {
  try {
    await orderStore.sendToBar(orderId)
    selectedOrder.value = null
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const serveOrder = async (orderId) => {
  try {
    await orderStore.serveOrder(orderId)
    selectedOrder.value = null
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const editOrder = (order) => {
  // TODO: Implement edit order functionality
  alert('Chức năng chỉnh sửa order đang được phát triển')
}

// Lifecycle
onMounted(async () => {
  await Promise.all([
    shiftStore.fetchCurrentShift(),
    orderStore.fetchOrders(),
    menuStore.fetchMenuItems()
  ])
})
</script>

<style scoped>
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
  transition: transform 0.3s ease;
}

.slide-up-enter-from {
  transform: translateY(100%);
}

.slide-up-leave-to {
  transform: translateY(100%);
}
</style>
