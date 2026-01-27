<template>
  <div v-if="show" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg p-6 w-96">
      <h3 class="text-lg font-semibold mb-4">Báo cáo sai lệch</h3>
      
      <div class="mb-4">
        <p class="text-sm text-gray-600 mb-2">Order: {{ payment?.order_id }}</p>
        <p class="text-sm text-gray-600 mb-2">Bàn: {{ payment?.table_name }}</p>
        <p class="text-sm text-gray-600 mb-4">Số tiền gốc: {{ formatCurrency(payment?.amount) }}</p>
      </div>

      <div class="mb-4">
        <label class="block text-sm font-medium text-gray-700 mb-2">Lý do sai lệch</label>
        <textarea
          v-model="reason"
          class="w-full border border-gray-300 rounded-lg px-3 py-2 h-20 resize-none"
          placeholder="Mô tả sai lệch..."
          required
        ></textarea>
      </div>

      <div class="mb-4">
        <label class="block text-sm font-medium text-gray-700 mb-2">Số tiền sai lệch</label>
        <input
          v-model.number="amount"
          type="number"
          step="0.01"
          class="w-full border border-gray-300 rounded-lg px-3 py-2"
          placeholder="0.00"
          required
        />
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
          :disabled="!reason.trim() || !amount"
          class="px-4 py-2 bg-yellow-600 text-white rounded-lg hover:bg-yellow-700 disabled:opacity-50"
        >
          Báo cáo
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, watch } from 'vue'

export default {
  name: 'DiscrepancyModal',
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
    const amount = ref(null)

    watch(() => props.show, (newVal) => {
      if (newVal) {
        reason.value = ''
        amount.value = null
      }
    })

    const close = () => {
      emit('close')
    }

    const confirm = () => {
      if (reason.value.trim() && amount.value) {
        emit('confirm', {
          reason: reason.value.trim(),
          amount: amount.value
        })
        reason.value = ''
        amount.value = null
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
      amount,
      close,
      confirm,
      formatCurrency
    }
  }
}
</script>