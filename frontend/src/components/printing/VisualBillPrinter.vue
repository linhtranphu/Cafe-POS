<template>
  <div class="h-full flex flex-col bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm flex-shrink-0 px-4 py-3 border-b">
      <h2 class="text-lg font-bold text-gray-800">🖨️ In Bill (Visual Template)</h2>
      <p class="text-sm text-gray-600 mt-1">In bill với template chính xác như preview.go</p>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto p-6">
      <div class="max-w-4xl mx-auto space-y-6">
        <!-- Printer Settings -->
        <div class="bg-white rounded-lg shadow-sm p-6">
          <h3 class="font-bold text-gray-800 mb-4">⚙️ Cài đặt máy in</h3>
          
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-bold text-gray-700 mb-2">IP Máy in</label>
              <input
                v-model="printerIP"
                type="text"
                placeholder="192.168.1.115"
                class="w-full px-4 py-2 border-2 border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              />
              <p class="text-xs text-gray-500 mt-1">Nhập IP của máy in Zywell ZY303 (port 9100)</p>
            </div>
          </div>
        </div>

        <!-- Order Selection -->
        <div class="bg-white rounded-lg shadow-sm p-6">
          <h3 class="font-bold text-gray-800 mb-4">📋 Chọn Order để in</h3>
          
          <div class="space-y-4">
            <!-- Search -->
            <div>
              <label class="block text-sm font-bold text-gray-700 mb-2">Tìm Order</label>
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Nhập số order hoặc tên khách..."
                class="w-full px-4 py-2 border-2 border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                @input="searchOrders"
              />
            </div>

            <!-- Loading -->
            <div v-if="loading" class="text-center py-8">
              <div class="inline-block w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
              <p class="text-gray-500 mt-2">Đang tải...</p>
            </div>

            <!-- Orders List -->
            <div v-else-if="orders.length > 0" class="space-y-2 max-h-96 overflow-y-auto">
              <div
                v-for="order in orders"
                :key="order.id"
                @click="selectedOrder = order"
                :class="[
                  'p-4 border-2 rounded-lg cursor-pointer transition-all',
                  selectedOrder?.id === order.id
                    ? 'border-blue-500 bg-blue-50'
                    : 'border-gray-200 hover:border-gray-300'
                ]"
              >
                <div class="flex justify-between items-start">
                  <div>
                    <div class="font-bold text-gray-800">{{ order.order_number }}</div>
                    <div class="text-sm text-gray-600 mt-1">
                      Waiter: {{ order.waiter_name || 'N/A' }}
                    </div>
                    <div class="text-sm text-gray-600">
                      {{ formatMoney(order.total) }} VNĐ
                    </div>
                  </div>
                  <div class="text-right">
                    <div class="text-xs text-gray-500">
                      {{ formatDate(order.created_at) }}
                    </div>
                    <div class="text-xs font-medium mt-1" :class="getStatusColor(order.status)">
                      {{ order.status }}
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Empty State -->
            <div v-else class="text-center py-8 text-gray-500">
              <div class="text-4xl mb-2">📋</div>
              <p>Không tìm thấy order nào</p>
            </div>
          </div>
        </div>

        <!-- Selected Order Preview -->
        <div v-if="selectedOrder" class="bg-white rounded-lg shadow-sm p-6">
          <h3 class="font-bold text-gray-800 mb-4">👁️ Order đã chọn</h3>
          
          <div class="space-y-3">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <div class="text-sm text-gray-600">Order Number</div>
                <div class="font-bold">{{ selectedOrder.order_number }}</div>
              </div>
              <div>
                <div class="text-sm text-gray-600">Waiter</div>
                <div class="font-bold">{{ selectedOrder.waiter_name || 'N/A' }}</div>
              </div>
              <div>
                <div class="text-sm text-gray-600">Payment Method</div>
                <div class="font-bold">{{ selectedOrder.payment_method }}</div>
              </div>
              <div>
                <div class="text-sm text-gray-600">Total</div>
                <div class="font-bold text-green-600">{{ formatMoney(selectedOrder.total) }} VNĐ</div>
              </div>
            </div>

            <div>
              <div class="text-sm text-gray-600 mb-2">Items ({{ selectedOrder.items?.length || 0 }})</div>
              <div class="space-y-1">
                <div
                  v-for="(item, idx) in selectedOrder.items"
                  :key="idx"
                  class="text-sm flex justify-between"
                >
                  <span>{{ item.quantity }}x {{ item.name }} {{ item.variant_name ? `(${item.variant_name})` : '' }}</span>
                  <span class="font-medium">{{ formatMoney(item.subtotal) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div v-if="selectedOrder" class="flex gap-3">
          <button
            @click="printVisualBill"
            :disabled="printing"
            class="flex-1 bg-green-500 text-white px-6 py-3 rounded-lg font-bold hover:bg-green-600 active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <span v-if="printing" class="inline-block w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
            <span v-else>🖨️</span>
            <span>{{ printing ? 'Đang in...' : 'In Bill' }}</span>
          </button>
          
          <button
            @click="previewBill"
            :disabled="previewing"
            class="flex-1 bg-purple-500 text-white px-6 py-3 rounded-lg font-bold hover:bg-purple-600 active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <span v-if="previewing" class="inline-block w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
            <span v-else>👁️</span>
            <span>{{ previewing ? 'Đang tạo...' : 'Preview' }}</span>
          </button>
        </div>

        <!-- Result Message -->
        <div v-if="resultMessage" :class="[
          'p-4 rounded-lg',
          resultMessage.type === 'success' ? 'bg-green-50 text-green-800' : 'bg-red-50 text-red-800'
        ]">
          <div class="flex items-center gap-2">
            <span>{{ resultMessage.type === 'success' ? '✅' : '❌' }}</span>
            <span class="font-bold">{{ resultMessage.text }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../../services/api'

const printerIP = ref('192.168.1.115')
const searchQuery = ref('')
const orders = ref([])
const selectedOrder = ref(null)
const loading = ref(false)
const printing = ref(false)
const previewing = ref(false)
const resultMessage = ref(null)

const formatMoney = (amount) => {
  if (!amount) return '0'
  return amount.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  const date = new Date(dateString)
  return date.toLocaleString('vi-VN', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getStatusColor = (status) => {
  const colors = {
    'COMPLETED': 'text-green-600',
    'PENDING': 'text-yellow-600',
    'CANCELLED': 'text-red-600',
    'REFUNDED': 'text-gray-600'
  }
  return colors[status] || 'text-gray-600'
}

const searchOrders = async () => {
  loading.value = true
  resultMessage.value = null
  
  try {
    const params = {}
    if (searchQuery.value) {
      params.search = searchQuery.value
    }
    
    const response = await api.get('/manager/orders', { params })
    orders.value = response.data.orders || response.data || []
  } catch (error) {
    console.error('Failed to search orders:', error)
    resultMessage.value = {
      type: 'error',
      text: 'Lỗi tải danh sách orders'
    }
  } finally {
    loading.value = false
  }
}

const printVisualBill = async () => {
  if (!selectedOrder.value) return
  if (!printerIP.value) {
    alert('Vui lòng nhập IP máy in')
    return
  }
  
  printing.value = true
  resultMessage.value = null
  
  try {
    const response = await api.post('/manager/visual-print/bill', {
      order_id: selectedOrder.value.id,
      printer_ip: printerIP.value
    })
    
    resultMessage.value = {
      type: 'success',
      text: `In thành công! Order: ${response.data.order_number}`
    }
  } catch (error) {
    console.error('Failed to print:', error)
    resultMessage.value = {
      type: 'error',
      text: error.response?.data?.error || 'Lỗi khi in bill'
    }
  } finally {
    printing.value = false
  }
}

const previewBill = async () => {
  if (!selectedOrder.value) return
  
  previewing.value = true
  resultMessage.value = null
  
  try {
    const response = await api.get(`/manager/visual-print/preview/${selectedOrder.value.id}`)
    
    resultMessage.value = {
      type: 'success',
      text: `Preview đã tạo: ${response.data.filename}`
    }
  } catch (error) {
    console.error('Failed to preview:', error)
    resultMessage.value = {
      type: 'error',
      text: error.response?.data?.error || 'Lỗi khi tạo preview'
    }
  } finally {
    previewing.value = false
  }
}

onMounted(() => {
  searchOrders()
})
</script>
