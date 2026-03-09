<template>
  <div v-if="show" class="fixed inset-0 z-50 overflow-y-auto" @click.self="$emit('close')">
    <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:p-0">
      <!-- Background overlay -->
      <div class="fixed inset-0 transition-opacity bg-gray-500 bg-opacity-75" @click="$emit('close')"></div>

      <!-- Modal panel -->
      <div class="relative inline-block w-full max-w-lg p-6 my-8 overflow-hidden text-left align-middle transition-all transform bg-white shadow-xl rounded-2xl">
        <!-- Header -->
        <div class="mb-4">
          <h3 class="text-xl font-bold text-gray-900">🔗 Gộp Bill</h3>
          <p class="mt-1 text-sm text-gray-500">Gộp {{ selectedOrders.length }} orders thành một bill</p>
        </div>

        <!-- Orders to merge -->
        <div class="mb-4 space-y-2">
          <p class="text-sm font-semibold text-gray-700">Các orders sẽ gộp:</p>
          <div class="max-h-48 overflow-y-auto space-y-2">
            <div 
              v-for="order in selectedOrders" 
              :key="order.id"
              class="p-3 bg-gray-50 rounded-lg border border-gray-200">
              <div class="flex items-center justify-between mb-1">
                <span class="font-semibold text-gray-900">#{{ order.order_number }}</span>
                <span :class="[
                  'px-2 py-0.5 text-xs font-semibold rounded-full',
                  order.status === 'PAID' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
                ]">
                  {{ order.status === 'PAID' ? '💰 Đã thu' : '📄 Chưa thu' }}
                </span>
              </div>
              <div class="text-sm text-gray-600">
                {{ order.items.length }} món • {{ formatCurrency(order.total) }}
                <span v-if="order.amount_paid > 0" class="text-green-600">
                  (Đã thu: {{ formatCurrency(order.amount_paid) }})
                </span>
              </div>
              <div v-if="order.customer_name" class="text-xs text-gray-500 mt-1">
                👤 {{ order.customer_name }}
              </div>
            </div>
          </div>
        </div>

        <!-- Summary -->
        <div class="mb-4 p-4 bg-blue-50 rounded-lg border border-blue-200">
          <div class="space-y-2">
            <div class="flex justify-between text-sm">
              <span class="text-gray-700">Tổng số món:</span>
              <span class="font-semibold text-gray-900">{{ totalItems }} món</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-700">Tổng tiền:</span>
              <span class="font-semibold text-gray-900">{{ formatCurrency(totalAmount) }}</span>
            </div>
            <div v-if="totalPaid > 0" class="flex justify-between text-sm">
              <span class="text-green-700">Đã thu:</span>
              <span class="font-semibold text-green-700">{{ formatCurrency(totalPaid) }}</span>
            </div>
            <div class="flex justify-between text-base pt-2 border-t border-blue-300">
              <span class="font-semibold text-gray-900">Còn lại:</span>
              <span class="font-bold text-blue-600">{{ formatCurrency(totalDue) }}</span>
            </div>
          </div>
        </div>

        <!-- Customer name input -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Tên khách hàng
          </label>
          <input
            v-model="customerName"
            type="text"
            placeholder="Nhập tên khách (tùy chọn)"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        <!-- Note input -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Ghi chú
          </label>
          <textarea
            v-model="note"
            rows="2"
            placeholder="Ghi chú cho order mới (tùy chọn)"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
          ></textarea>
        </div>

        <!-- Actions -->
        <div class="flex gap-3">
          <button
            @click="$emit('close')"
            class="flex-1 px-4 py-3 text-gray-700 bg-gray-100 rounded-xl font-semibold hover:bg-gray-200 transition-colors"
          >
            Hủy
          </button>
          <button
            @click="handleMerge"
            :disabled="loading"
            class="flex-1 px-4 py-3 text-white bg-blue-500 rounded-xl font-semibold hover:bg-blue-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="loading">Đang gộp...</span>
            <span v-else>✅ Xác nhận gộp</span>
          </button>
        </div>

        <!-- Error message -->
        <div v-if="error" class="mt-4 p-3 bg-red-50 border border-red-200 rounded-lg">
          <p class="text-sm text-red-800">❌ {{ error }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useOrderStore } from '../stores/order'

const props = defineProps({
  show: Boolean,
  selectedOrders: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['close', 'merged'])

const orderStore = useOrderStore()

const customerName = ref('')
const note = ref('')
const loading = ref(false)
const error = ref(null)

// Watch for selectedOrders changes to set default customer name
watch(() => props.selectedOrders, (newOrders) => {
  if (newOrders.length > 0 && newOrders[0].customer_name && !customerName.value) {
    customerName.value = newOrders[0].customer_name
  }
}, { immediate: true })

const totalItems = computed(() => {
  return props.selectedOrders.reduce((sum, order) => sum + order.items.length, 0)
})

const totalAmount = computed(() => {
  return props.selectedOrders.reduce((sum, order) => sum + order.total, 0)
})

const totalPaid = computed(() => {
  return props.selectedOrders.reduce((sum, order) => sum + order.amount_paid, 0)
})

const totalDue = computed(() => {
  return totalAmount.value - totalPaid.value
})

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND'
  }).format(amount)
}

const handleMerge = async () => {
  if (loading.value) return

  loading.value = true
  error.value = null

  try {
    const orderIds = props.selectedOrders.map(o => o.id)
    const response = await orderStore.mergeOrders(orderIds, customerName.value, note.value)
    
    emit('merged', response)
    emit('close')
  } catch (err) {
    error.value = err.response?.data?.error || err.message || 'Lỗi gộp bill'
  } finally {
    loading.value = false
  }
}
</script>
