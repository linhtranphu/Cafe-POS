<template>
  <div class="min-h-screen bg-gray-100">
    <Navigation />
    <div class="p-4">
      <!-- Header -->
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold text-gray-800">📋 Quản lý Order</h2>
        <button @click="showCreateForm = true" class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg font-medium">
          + Tạo Order
        </button>
      </div>

      <!-- Shift Warning -->
      <div v-if="!hasOpenShift" class="bg-yellow-50 border-l-4 border-yellow-400 p-4 mb-4">
        <p class="text-yellow-700">⚠️ Bạn chưa mở ca. Vui lòng mở ca trước khi tạo order.</p>
        <button @click="$router.push('/shifts')" class="mt-2 bg-yellow-500 text-white px-4 py-2 rounded">Mở ca ngay</button>
      </div>

      <!-- Status Tabs -->
      <div class="flex gap-2 mb-4 overflow-x-auto">
        <button v-for="status in statuses" :key="status.value" @click="filterStatus = status.value"
          :class="filterStatus === status.value ? 'bg-blue-500 text-white' : 'bg-white text-gray-700'"
          class="px-4 py-2 rounded-lg font-medium whitespace-nowrap">
          {{ status.label }} ({{ getOrderCountByStatus(status.value) }})
        </button>
      </div>

      <!-- Orders List -->
      <div v-if="loading" class="text-center py-10">Đang tải...</div>
      <div v-else-if="filteredOrders.length === 0" class="text-center py-10 text-gray-500">Không có order nào</div>
      <div v-else class="grid gap-4">
        <div v-for="order in filteredOrders" :key="order.id" class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex justify-between items-start mb-3">
            <div>
              <h3 class="font-bold text-lg">{{ order.table_name }}</h3>
              <p class="text-sm text-gray-500">{{ formatDate(order.created_at) }}</p>
            </div>
            <span :class="getStatusColor(order.status)" class="px-3 py-1 rounded-full text-xs font-medium">
              {{ getStatusText(order.status) }}
            </span>
          </div>

          <!-- Items -->
          <div class="mb-3 space-y-1">
            <div v-for="item in order.items" :key="item.menu_item_id" class="flex justify-between text-sm">
              <span>{{ item.name }} x{{ item.quantity }}</span>
              <span class="font-medium">{{ formatPrice(item.subtotal) }}</span>
            </div>
          </div>

          <!-- Total -->
          <div class="border-t pt-2 mb-3">
            <div class="flex justify-between font-bold text-lg">
              <span>Tổng cộng:</span>
              <span class="text-green-600">{{ formatPrice(order.total) }}</span>
            </div>
          </div>

          <!-- Actions -->
          <div class="grid grid-cols-2 gap-2">
            <button v-if="order.status === 'CREATED'" @click="confirmOrder(order)" class="bg-yellow-500 hover:bg-yellow-600 text-white px-3 py-2 rounded-lg text-sm">
              Xác nhận
            </button>
            <button v-if="order.status === 'UNPAID'" @click="showPaymentForm(order)" class="bg-green-500 hover:bg-green-600 text-white px-3 py-2 rounded-lg text-sm">
              💰 Thu tiền
            </button>
            <button v-if="order.status === 'PAID'" @click="sendToKitchen(order.id)" class="bg-blue-500 hover:bg-blue-600 text-white px-3 py-2 rounded-lg text-sm">
              🍳 Gửi pha chế
            </button>
            <button v-if="order.status === 'IN_PROGRESS'" @click="serveOrder(order.id)" class="bg-purple-500 hover:bg-purple-600 text-white px-3 py-2 rounded-lg text-sm">
              ✅ Đã phục vụ
            </button>
            <button v-if="isCashier && ['UNPAID', 'PAID', 'IN_PROGRESS'].includes(order.status)" @click="showRefundForm(order)" class="bg-red-500 hover:bg-red-600 text-white px-3 py-2 rounded-lg text-sm">
              Hoàn tiền
            </button>
          </div>
        </div>
      </div>

      <!-- Create Order Modal -->
      <div v-if="showCreateForm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-xl p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
          <h3 class="text-xl font-bold mb-4">Tạo Order Mới</h3>
          <form @submit.prevent="createOrder" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Chọn bàn *</label>
              <select v-model="form.table_id" required class="w-full p-3 border rounded-lg">
                <option value="">-- Chọn bàn --</option>
                <option v-for="table in emptyTables" :key="table.id" :value="table.id">
                  {{ table.name }} ({{ table.capacity }} chỗ)
                </option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-medium mb-2">Chọn món</label>
              <div class="grid grid-cols-2 gap-2 max-h-60 overflow-y-auto">
                <button v-for="item in menuItems" :key="item.id" type="button" @click="addItem(item)"
                  class="p-3 border rounded-lg hover:bg-blue-50 text-left">
                  <div class="font-medium">{{ item.name }}</div>
                  <div class="text-sm text-gray-500">{{ formatPrice(item.price) }}</div>
                </button>
              </div>
            </div>

            <div v-if="form.items.length > 0">
              <label class="block text-sm font-medium mb-2">Món đã chọn</label>
              <div class="space-y-2">
                <div v-for="(item, index) in form.items" :key="index" class="flex items-center gap-2 p-2 bg-gray-50 rounded">
                  <span class="flex-1">{{ item.name }}</span>
                  <input v-model.number="item.quantity" type="number" min="1" class="w-16 p-1 border rounded text-center">
                  <button type="button" @click="removeItem(index)" class="text-red-500 hover:text-red-700">✕</button>
                </div>
              </div>
            </div>

            <div class="flex gap-2">
              <button type="button" @click="showCreateForm = false" class="flex-1 bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg">
                Hủy
              </button>
              <button type="submit" :disabled="!form.table_id || form.items.length === 0" class="flex-1 bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg">
                Tạo Order
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Payment Modal -->
      <div v-if="showPayment" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-xl p-6 w-full max-w-md">
          <h3 class="text-xl font-bold mb-4">Thu tiền</h3>
          <div class="mb-4">
            <p class="text-2xl font-bold text-green-600">{{ formatPrice(selectedOrder?.total) }}</p>
          </div>
          <div class="space-y-3">
            <button @click="payOrder('CASH')" class="w-full bg-green-500 hover:bg-green-600 text-white px-4 py-3 rounded-lg font-medium">
              💵 Tiền mặt
            </button>
            <button @click="payOrder('QR')" class="w-full bg-blue-500 hover:bg-blue-600 text-white px-4 py-3 rounded-lg font-medium">
              📱 QR Code
            </button>
            <button @click="payOrder('TRANSFER')" class="w-full bg-purple-500 hover:bg-purple-600 text-white px-4 py-3 rounded-lg font-medium">
              🏦 Chuyển khoản
            </button>
            <button @click="showPayment = false" class="w-full bg-gray-500 hover:bg-gray-600 text-white px-4 py-3 rounded-lg font-medium">
              Hủy
            </button>
          </div>
        </div>
      </div>

      <!-- Refund Modal -->
      <div v-if="showRefund" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-xl p-6 w-full max-w-md">
          <h3 class="text-xl font-bold mb-4">Hoàn tiền</h3>
          <form @submit.prevent="refundOrder" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Lý do *</label>
              <textarea v-model="refundReason" required rows="3" class="w-full p-3 border rounded-lg" placeholder="Nhập lý do hoàn tiền..."></textarea>
            </div>
            <div class="flex gap-2">
              <button type="button" @click="showRefund = false" class="flex-1 bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg">
                Hủy
              </button>
              <button type="submit" class="flex-1 bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-lg">
                Xác nhận hoàn tiền
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useOrderStore } from '../stores/order'
import { useTableStore } from '../stores/table'
import { useShiftStore } from '../stores/shift'
import { useMenuStore } from '../stores/menu'
import { useAuthStore } from '../stores/auth'
import Navigation from '../components/Navigation.vue'

const orderStore = useOrderStore()
const tableStore = useTableStore()
const shiftStore = useShiftStore()
const menuStore = useMenuStore()
const authStore = useAuthStore()

const showCreateForm = ref(false)
const showPayment = ref(false)
const showRefund = ref(false)
const selectedOrder = ref(null)
const refundReason = ref('')
const filterStatus = ref('ALL')

const form = ref({
  table_id: '',
  items: [],
  shift_id: ''
})

const statuses = [
  { value: 'ALL', label: 'Tất cả' },
  { value: 'CREATED', label: 'Mới tạo' },
  { value: 'UNPAID', label: 'Chưa thanh toán' },
  { value: 'PAID', label: 'Đã thanh toán' },
  { value: 'IN_PROGRESS', label: 'Đang pha chế' },
  { value: 'SERVED', label: 'Đã phục vụ' }
]

const loading = computed(() => orderStore.loading)
const orders = computed(() => orderStore.orders)
const emptyTables = computed(() => tableStore.emptyTables)
const menuItems = computed(() => menuStore.items)
const hasOpenShift = computed(() => shiftStore.hasOpenShift)
const isCashier = computed(() => authStore.user?.role === 'cashier' || authStore.user?.role === 'manager')

const filteredOrders = computed(() => {
  if (filterStatus.value === 'ALL') return orders.value
  return orders.value.filter(o => o.status === filterStatus.value)
})

onMounted(async () => {
  await Promise.all([
    shiftStore.fetchCurrentShift(),
    orderStore.fetchOrders(),
    tableStore.fetchTables(),
    menuStore.fetchMenuItems()
  ])
})

const addItem = (item) => {
  const existing = form.value.items.find(i => i.menu_item_id === item.id)
  if (existing) {
    existing.quantity++
  } else {
    form.value.items.push({
      menu_item_id: item.id,
      name: item.name,
      price: item.price,
      quantity: 1
    })
  }
}

const removeItem = (index) => {
  form.value.items.splice(index, 1)
}

const createOrder = async () => {
  try {
    form.value.shift_id = shiftStore.currentShift.id
    await orderStore.createOrder(form.value)
    showCreateForm.value = false
    form.value = { table_id: '', items: [], shift_id: '' }
    await tableStore.fetchTables()
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

const confirmOrder = async (order) => {
  try {
    await orderStore.confirmOrder(order.id, 0)
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const showPaymentForm = (order) => {
  selectedOrder.value = order
  showPayment.value = true
}

const payOrder = async (method) => {
  try {
    await orderStore.payOrder(selectedOrder.value.id, method)
    showPayment.value = false
    selectedOrder.value = null
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const sendToKitchen = async (id) => {
  try {
    await orderStore.sendToKitchen(id)
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const serveOrder = async (id) => {
  try {
    await orderStore.serveOrder(id)
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const showRefundForm = (order) => {
  selectedOrder.value = order
  showRefund.value = true
}

const refundOrder = async () => {
  try {
    await orderStore.refundOrder(selectedOrder.value.id, refundReason.value)
    showRefund.value = false
    selectedOrder.value = null
    refundReason.value = ''
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const getOrderCountByStatus = (status) => {
  if (status === 'ALL') return orders.value.length
  return orders.value.filter(o => o.status === status).length
}

const getStatusColor = (status) => {
  const colors = {
    CREATED: 'bg-gray-100 text-gray-800',
    UNPAID: 'bg-yellow-100 text-yellow-800',
    PAID: 'bg-green-100 text-green-800',
    IN_PROGRESS: 'bg-blue-100 text-blue-800',
    SERVED: 'bg-purple-100 text-purple-800',
    CANCELLED: 'bg-red-100 text-red-800',
    REFUNDED: 'bg-orange-100 text-orange-800'
  }
  return colors[status] || 'bg-gray-100 text-gray-800'
}

const getStatusText = (status) => {
  const texts = {
    CREATED: 'Mới tạo',
    UNPAID: 'Chưa thanh toán',
    PAID: 'Đã thanh toán',
    IN_PROGRESS: 'Đang pha chế',
    SERVED: 'Đã phục vụ',
    CANCELLED: 'Đã hủy',
    REFUNDED: 'Đã hoàn tiền',
    LOCKED: 'Đã khóa'
  }
  return texts[status] || status
}

const formatPrice = (price) => {
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(price)
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('vi-VN')
}
</script>

<style scoped>
button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
