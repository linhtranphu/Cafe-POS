<template>
  <div v-if="show" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg p-6 w-96">
      <h3 class="text-lg font-semibold mb-4">Điều chỉnh thanh toán</h3>
      
      <div class="mb-4">
        <p class="text-sm text-gray-600 mb-2">Order: {{ payment?.order_id }}</p>
        <p class="text-sm text-gray-600 mb-2">Bàn: {{ payment?.table_name }}</p>
        <p class="text-sm text-gray-600 mb-4">Số tiền: {{ formatCurrency(payment?.amount) }}</p>
      </div>

      <div class="mb-4">
        <label class="block text-sm font-medium text-gray-700 mb-2">Lý do điều chỉnh</label>
        <textarea
          v-model="reason"
          class="w-full border border-gray-300 rounded-lg px-3 py-2 h-20 resize-none"
          placeholder="Nhập lý do điều chỉnh..."
          required
        ></textarea>
      </div>

      <div class="flex justify-end space-x-3">
        <button
          @click="close"
          class="px-4 py-2 text-gray-600 border border-gray-300 rounded-lg hover:bg-gray-50"
        >
          Hủy
        </button>
        <button
          @click="confirm"
          :disabled="!reason.trim()"
          class="px-4 py-2 bg-orange-600 text-white rounded-lg hover:bg-orange-700 disabled:opacity-50"
        >
          Xác nhận
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, watch } from 'vue'

export default {
  name: 'OverridePaymentModal',
  props: {
    show: {
      type: Boolean,
      default: false
    },
    payment: {
      type: Object,
      default: null
    }
  },
  emits: ['close', 'confirm'],
  setup(props, { emit }) {
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

    const formatCurrency = (amount) => {
      return new Intl.NumberFormat('vi-VN', {
        style: 'currency',
        currency: 'VND'
      }).format(amount)
    }

    return {
      reason,
      close,
      confirm,
      formatCurrency
    }
  }
}
</script>