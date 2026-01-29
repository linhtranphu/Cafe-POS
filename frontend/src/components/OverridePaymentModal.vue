<template>
  <div v-if="show" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
    <div class="bg-white rounded-2xl p-6 w-full max-w-md">
      <h3 class="text-xl font-bold text-gray-800 mb-4">✏️ Điều chỉnh thanh toán</h3>
      
      <div class="bg-gray-50 rounded-xl p-4 mb-4 space-y-2">
        <div class="flex justify-between items-center">
          <span class="text-sm text-gray-600">Order:</span>
          <span class="font-mono text-sm font-medium">#{{ payment?.order_id?.slice(-6) }}</span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-sm text-gray-600">Khách hàng:</span>
          <span class="font-medium">{{ payment?.customer_name || 'Khách lẻ' }}</span>
        </div>
        <div class="flex justify-between items-center pt-2 border-t border-gray-200">
          <span class="text-sm text-gray-600">Số tiền:</span>
          <span class="font-bold text-green-600">{{ formatPrice(payment?.amount) }}</span>
        </div>
      </div>

      <div class="mb-6">
        <label class="block text-sm font-medium text-gray-700 mb-2">Lý do điều chỉnh *</label>
        <textarea
          v-model="reason"
          class="w-full border-2 border-gray-300 rounded-xl px-4 py-3 text-base resize-none focus:outline-none focus:border-orange-500"
          rows="4"
          placeholder="Nhập lý do điều chỉnh thanh toán..."
          required
        ></textarea>
      </div>

      <div class="flex gap-3">
        <button
          @click="close"
          class="flex-1 py-3 text-gray-700 bg-gray-100 rounded-xl font-medium active:scale-95 transition-transform"
        >
          Hủy
        </button>
        <button
          @click="confirm"
          :disabled="!reason.trim()"
          class="flex-1 py-3 bg-orange-500 text-white rounded-xl font-medium active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Xác nhận
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  payment: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'confirm'])

const reason = ref('')

watch(() => props.show, (newVal) => {
  if (newVal) {
    reason.value = ''
  }
})

const close = () => {
  emit('close')
}

const confirm = () => {
  if (reason.value.trim()) {
    emit('confirm', reason.value.trim())
    reason.value = ''
  }
}

const formatPrice = (amount) => {
  if (!amount && amount !== 0) return '0₫'
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND',
    maximumFractionDigits: 0
  }).format(amount)
}
</script>

<style scoped>
.active\:scale-95:active {
  transform: scale(0.95);
}
</style>
