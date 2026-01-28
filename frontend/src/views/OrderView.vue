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
        <div v-for="order in filteredOrders" :key="order.id" 
          :class="['bg-white rounded-xl p-4 shadow-sm transition-all duration-500', 
                   order.refund_amount > 0 ? 'ring-2 ring-orange-200 bg-gradient-to-r from-white to-orange-50' : '']">
          <div class="flex justify-between items-start mb-3">
            <div>
              <h3 class="font-bold text-lg">{{ order.order_number }}</h3>
              <p class="text-sm text-gray-600">{{ order.customer_name || 'Khách lẻ' }}</p>
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

          <!-- Refund Info -->
          <div v-if="order.refund_amount > 0" class="mb-2 p-3 bg-gradient-to-r from-orange-50 to-yellow-50 rounded-lg border-l-4 border-orange-400 shadow-sm">
            <div class="flex items-center gap-2 text-orange-700 font-medium">
              <span class="text-lg">💸</span>
              <span>Đã hoàn tiền: <span class="font-bold text-orange-800">{{ formatPrice(order.refund_amount) }}</span></span>
            </div>
            <div v-if="order.refund_reason" class="text-sm text-orange-600 mt-1 ml-6">
              {{ order.refund_reason }}
            </div>
          </div>

          <!-- Payment Status -->
          <div v-if="order.status === 'PAID' && order.amount_due > 0" class="mb-2 p-3 bg-gradient-to-r from-yellow-50 to-amber-50 rounded-lg border-l-4 border-yellow-400 shadow-sm">
            <div class="flex items-center gap-2 text-yellow-700 font-medium">
              <span class="text-lg">⚠️</span>
              <span>Còn thiếu: <span class="font-bold text-yellow-800">{{ formatPrice(order.amount_due) }}</span></span>
            </div>
          </div>

          <!-- Total -->
          <div class="border-t pt-2 mb-3">
            <div class="flex justify-between font-bold text-lg">
              <span>Tổng cộng:</span>
              <span class="text-green-600">{{ formatPrice(order.total) }}</span>
            </div>
            <div v-if="order.amount_paid > 0" class="flex justify-between text-sm text-gray-600">
              <span>Đã thu:</span>
              <span>{{ formatPrice(order.amount_paid) }}</span>
            </div>
            <div v-if="order.refund_amount > 0" class="flex justify-between text-sm text-orange-600">
              <span>Đã hoàn:</span>
              <span>-{{ formatPrice(order.refund_amount) }}</span>
            </div>
          </div>

          <!-- Actions -->
          <div class="grid grid-cols-2 gap-2">
            <button v-if="order.status === 'CREATED'" @click="showPaymentForm(order)" class="bg-green-500 hover:bg-green-600 text-white px-3 py-2 rounded-lg text-sm">
              💰 Thu tiền
            </button>
            <button v-if="order.status === 'CREATED'" @click="showEditForm(order)" class="bg-blue-500 hover:bg-blue-600 text-white px-3 py-2 rounded-lg text-sm">
              ✏️ Chỉnh sửa
            </button>
            <button v-if="order.status === 'PAID' && order.amount_due > 0" @click="showPaymentForm(order)" class="bg-yellow-500 hover:bg-yellow-600 text-white px-3 py-2 rounded-lg text-sm">
              💰 Thu thêm
            </button>
            <button v-if="order.status === 'PAID' && order.amount_due <= 0" @click="sendToBar(order.id)" class="bg-blue-500 hover:bg-blue-600 text-white px-3 py-2 rounded-lg text-sm">
              🍹 Gửi quầy bar
            </button>
            <button v-if="order.status === 'PAID'" @click="showEditForm(order)" class="bg-purple-500 hover:bg-purple-600 text-white px-3 py-2 rounded-lg text-sm">
              ✏️ Chỉnh sửa
            </button>
            <button v-if="order.status === 'IN_PROGRESS'" @click="serveOrder(order.id)" class="bg-green-500 hover:bg-green-600 text-white px-3 py-2 rounded-lg text-sm">
              ✅ Đã phục vụ
            </button>
            <button v-if="isCashier && ['CREATED', 'PAID'].includes(order.status)" @click="showCancelForm(order)" class="bg-red-500 hover:bg-red-600 text-white px-3 py-2 rounded-lg text-sm">
              ❌ Hủy order
            </button>
            <button v-if="isCashier && order.status === 'PAID' && order.amount_paid > 0" @click="showRefundForm(order)" class="bg-orange-500 hover:bg-orange-600 text-white px-3 py-2 rounded-lg text-sm">
              💸 Hoàn tiền
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
              <label class="block text-sm font-medium mb-2">Tên khách hàng</label>
              <input v-model="form.customer_name" type="text" class="w-full p-3 border rounded-lg" placeholder="Nhập tên khách hàng (tùy chọn)">
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

            <div>
              <label class="block text-sm font-medium mb-2">Ghi chú</label>
              <textarea v-model="form.note" rows="2" class="w-full p-3 border rounded-lg" placeholder="Ghi chú đặc biệt..."></textarea>
            </div>

            <div class="flex gap-2">
              <button type="button" @click="showCreateForm = false" class="flex-1 bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg">
                Hủy
              </button>
              <button type="submit" :disabled="form.items.length === 0" class="flex-1 bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg">
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
            <p class="text-lg">Tổng tiền: <span class="font-bold text-green-600">{{ formatPrice(selectedOrder?.total) }}</span></p>
            <p v-if="selectedOrder?.amount_paid > 0" class="text-sm text-gray-600">Đã thu: {{ formatPrice(selectedOrder?.amount_paid) }}</p>
            <p v-if="selectedOrder?.amount_due > 0" class="text-sm text-red-600">Còn thiếu: {{ formatPrice(selectedOrder?.amount_due) }}</p>
          </div>
          
          <form @submit.prevent="collectPayment" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Số tiền thu *</label>
              <input v-model.number="paymentForm.amount" type="number" step="0.01" min="0" required 
                class="w-full p-3 border rounded-lg" placeholder="Nhập số tiền">
            </div>
            
            <div>
              <label class="block text-sm font-medium mb-2">Phương thức thanh toán *</label>
              <div class="space-y-2">
                <label class="flex items-center">
                  <input v-model="paymentForm.payment_method" type="radio" value="CASH" class="mr-2">
                  💵 Tiền mặt
                </label>
                <label class="flex items-center">
                  <input v-model="paymentForm.payment_method" type="radio" value="QR" class="mr-2">
                  📱 QR Code
                </label>
                <label class="flex items-center">
                  <input v-model="paymentForm.payment_method" type="radio" value="TRANSFER" class="mr-2">
                  🏦 Chuyển khoản
                </label>
              </div>
            </div>
            
            <div class="flex gap-2">
              <button type="button" @click="showPayment = false" class="flex-1 bg-gray-500 hover:bg-gray-600 text-white px-4 py-3 rounded-lg font-medium">
                Hủy
              </button>
              <button type="submit" class="flex-1 bg-green-500 hover:bg-green-600 text-white px-4 py-3 rounded-lg font-medium">
                Thu tiền
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Edit Order Modal -->
      <div v-if="showEdit" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-xl p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
          <h3 class="text-xl font-bold mb-4">Chỉnh sửa Order</h3>
          <form @submit.prevent="editOrder" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Chọn món</label>
              <div class="grid grid-cols-2 gap-2 max-h-60 overflow-y-auto">
                <button v-for="item in menuItems" :key="item.id" type="button" @click="addItemToEdit(item)"
                  class="p-3 border rounded-lg hover:bg-blue-50 text-left">
                  <div class="font-medium">{{ item.name }}</div>
                  <div class="text-sm text-gray-500">{{ formatPrice(item.price) }}</div>
                </button>
              </div>
            </div>

            <div v-if="editForm.items.length > 0">
              <label class="block text-sm font-medium mb-2">Món đã chọn</label>
              <div class="space-y-2">
                <div v-for="(item, index) in editForm.items" :key="index" class="flex items-center gap-2 p-2 bg-gray-50 rounded">
                  <span class="flex-1">{{ item.name }}</span>
                  <input v-model.number="item.quantity" type="number" min="1" class="w-16 p-1 border rounded text-center">
                  <button type="button" @click="removeItemFromEdit(index)" class="text-red-500 hover:text-red-700">✕</button>
                </div>
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium mb-2">Giảm giá</label>
              <input v-model.number="editForm.discount" type="number" step="0.01" min="0" class="w-full p-3 border rounded-lg" placeholder="Số tiền giảm giá">
            </div>

            <div>
              <label class="block text-sm font-medium mb-2">Ghi chú</label>
              <textarea v-model="editForm.note" rows="2" class="w-full p-3 border rounded-lg" placeholder="Ghi chú đặc biệt..."></textarea>
            </div>

            <div class="flex gap-2">
              <button type="button" @click="showEdit = false" class="flex-1 bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg">
                Hủy
              </button>
              <button type="submit" :disabled="editForm.items.length === 0" class="flex-1 bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg">
                Cập nhật Order
              </button>
            </div>
            
            <!-- Warning about refund -->
            <div class="mt-2 p-2 bg-yellow-50 border border-yellow-200 rounded text-xs text-yellow-700">
              ⚠️ Lưu ý: Nếu tổng tiền mới thấp hơn số tiền đã thu, hệ thống sẽ tự động hoàn tiền phần chênh lệch.
            </div>
          </form>
        </div>
      </div>

      <!-- Cancel Modal -->
      <div v-if="showCancel" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-xl p-6 w-full max-w-md">
          <h3 class="text-xl font-bold mb-4">Hủy Order</h3>
          <form @submit.prevent="cancelOrder" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Lý do hủy *</label>
              <textarea v-model="cancelReason" required rows="3" class="w-full p-3 border rounded-lg" placeholder="Nhập lý do hủy order..."></textarea>
            </div>
            <div class="flex gap-2">
              <button type="button" @click="showCancel = false" class="flex-1 bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg">
                Hủy
              </button>
              <button type="submit" class="flex-1 bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-lg">
                Xác nhận hủy
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Refund Modal -->
      <div v-if="showRefund" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-xl p-6 w-full max-w-md">
          <h3 class="text-xl font-bold mb-4">Hoàn tiền</h3>
          <div class="mb-4">
            <p class="text-sm text-gray-600">Đã thu: {{ formatPrice(selectedOrder?.amount_paid) }}</p>
          </div>
          <form @submit.prevent="refundPartial" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Số tiền hoàn *</label>
              <input v-model.number="refundForm.amount" type="number" step="0.01" min="0" :max="selectedOrder?.amount_paid" required 
                class="w-full p-3 border rounded-lg" placeholder="Nhập số tiền hoàn">
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Lý do hoàn tiền *</label>
              <textarea v-model="refundForm.reason" required rows="3" class="w-full p-3 border rounded-lg" placeholder="Nhập lý do hoàn tiền..."></textarea>
            </div>
            <div class="flex gap-2">
              <button type="button" @click="showRefund = false" class="flex-1 bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg">
                Hủy
              </button>
              <button type="submit" class="flex-1 bg-orange-500 hover:bg-orange-600 text-white px-4 py-2 rounded-lg">
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
import { useShiftStore } from '../stores/shift'
import { useMenuStore } from '../stores/menu'
import { useAuthStore } from '../stores/auth'
import Navigation from '../components/Navigation.vue'

const orderStore = useOrderStore()
const shiftStore = useShiftStore()
const menuStore = useMenuStore()
const authStore = useAuthStore()

const showCreateForm = ref(false)
const showPayment = ref(false)
const showEdit = ref(false)
const showCancel = ref(false)
const showRefund = ref(false)
const selectedOrder = ref(null)
const cancelReason = ref('')
const filterStatus = ref('ALL')

const form = ref({
  customer_name: '',
  items: [],
  note: '',
  shift_id: ''
})

const paymentForm = ref({
  amount: 0,
  payment_method: 'CASH'
})

const editForm = ref({
  items: [],
  discount: 0,
  note: ''
})

const refundForm = ref({
  amount: 0,
  reason: ''
})

const statuses = [
  { value: 'ALL', label: 'Tất cả' },
  { value: 'CREATED', label: 'Mới tạo' },
  { value: 'PAID', label: 'Đã thanh toán' },
  { value: 'IN_PROGRESS', label: 'Đang pha chế' },
  { value: 'SERVED', label: 'Đã phục vụ' },
  { value: 'CANCELLED', label: 'Đã hủy' }
]

const loading = computed(() => orderStore.loading)
const orders = computed(() => orderStore.orders)
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
    form.value = { customer_name: '', items: [], note: '', shift_id: '' }
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

const showPaymentForm = (order) => {
  selectedOrder.value = order
  paymentForm.value.amount = order.amount_due || order.total
  paymentForm.value.payment_method = 'CASH'
  showPayment.value = true
}

const collectPayment = async () => {
  try {
    await orderStore.collectPayment(selectedOrder.value.id, paymentForm.value)
    showPayment.value = false
    selectedOrder.value = null
    paymentForm.value = { amount: 0, payment_method: 'CASH' }
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const showEditForm = (order) => {
  selectedOrder.value = order
  editForm.value = {
    items: [...order.items],
    discount: order.discount || 0,
    note: order.note || ''
  }
  showEdit.value = true
}

const addItemToEdit = (item) => {
  const existing = editForm.value.items.find(i => i.menu_item_id === item.id)
  if (existing) {
    existing.quantity++
  } else {
    editForm.value.items.push({
      menu_item_id: item.id,
      name: item.name,
      price: item.price,
      quantity: 1
    })
  }
}

const removeItemFromEdit = (index) => {
  editForm.value.items.splice(index, 1)
}

const editOrder = async () => {
  try {
    const response = await orderStore.editOrder(selectedOrder.value.id, editForm.value)
    showEdit.value = false
    selectedOrder.value = null
    editForm.value = { items: [], discount: 0, note: '' }
    
    // Refresh orders to show updated status
    await orderStore.fetchOrders()
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const sendToBar = async (id) => {
  try {
    await orderStore.sendToBar(id)
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

const showCancelForm = (order) => {
  selectedOrder.value = order
  cancelReason.value = ''
  showCancel.value = true
}

const cancelOrder = async () => {
  try {
    await orderStore.cancelOrder(selectedOrder.value.id, cancelReason.value)
    showCancel.value = false
    selectedOrder.value = null
    cancelReason.value = ''
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const showRefundForm = (order) => {
  selectedOrder.value = order
  refundForm.value = { amount: 0, reason: '' }
  showRefund.value = true
}

const refundPartial = async () => {
  try {
    await orderStore.refundPartial(selectedOrder.value.id, refundForm.value.amount, refundForm.value.reason)
    showRefund.value = false
    selectedOrder.value = null
    refundForm.value = { amount: 0, reason: '' }
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
    PAID: 'bg-green-100 text-green-800',
    IN_PROGRESS: 'bg-blue-100 text-blue-800',
    SERVED: 'bg-purple-100 text-purple-800',
    CANCELLED: 'bg-red-100 text-red-800',
    LOCKED: 'bg-gray-200 text-gray-600'
  }
  return colors[status] || 'bg-gray-100 text-gray-800'
}

const getStatusText = (status) => {
  const texts = {
    CREATED: 'Mới tạo',
    PAID: 'Đã thanh toán',
    IN_PROGRESS: 'Đang pha chế',
    SERVED: 'Đã phục vụ',
    CANCELLED: 'Đã hủy',
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
