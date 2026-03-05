<template>
  <div class="fixed inset-0 bg-black/50 flex items-end sm:items-center justify-center z-50" @click.self="$emit('close')">
    <div class="bg-white rounded-t-3xl sm:rounded-2xl w-full sm:max-w-md max-h-[90vh] overflow-y-auto animate-slide-up">
      <!-- Header -->
      <div class="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
        <h2 class="text-xl font-bold text-gray-800">Chi tiết giao dịch</h2>
        <button @click="$emit('close')" class="p-2 hover:bg-gray-100 rounded-lg transition-colors">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Content -->
      <div class="p-6 space-y-4">
        <!-- Transaction Type -->
        <div class="text-center py-4">
          <div class="text-6xl mb-2">{{ getTransactionTypeIcon(transaction.type) }}</div>
          <div class="text-2xl font-bold text-gray-800">{{ getTransactionTypeLabel(transaction.type) }}</div>
        </div>

        <!-- Amount -->
        <div class="bg-gray-50 rounded-xl p-4">
          <div class="text-center">
            <div class="text-sm text-gray-600 mb-1">Số tiền</div>
            <div
              class="text-3xl font-bold"
              :class="getAmountColor(transaction.type)"
            >
              {{ getAmountPrefix(transaction.type) }}{{ formatCurrency(transaction.total_amount) }}
            </div>
          </div>
          
          <!-- Money breakdown -->
          <div v-if="transaction.cash_amount > 0 || transaction.transfer_amount > 0" class="mt-4 pt-4 border-t border-gray-200">
            <div class="grid grid-cols-2 gap-3">
              <div v-if="transaction.cash_amount > 0" class="text-center">
                <div class="text-xs text-gray-600 mb-1">💵 Tiền mặt</div>
                <div class="text-lg font-semibold text-green-600">
                  {{ formatCurrency(transaction.cash_amount) }}
                </div>
              </div>
              <div v-if="transaction.transfer_amount > 0" class="text-center">
                <div class="text-xs text-gray-600 mb-1">💳 Chuyển khoản</div>
                <div class="text-lg font-semibold text-blue-600">
                  {{ formatCurrency(transaction.transfer_amount) }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Description -->
        <div>
          <div class="text-sm font-semibold text-gray-700 mb-2">Mô tả</div>
          <div class="bg-gray-50 rounded-lg p-3 text-gray-800">
            {{ transaction.description }}
          </div>
        </div>

        <!-- Info Grid -->
        <div class="space-y-3">
          <div class="flex justify-between py-2 border-b border-gray-100">
            <span class="text-sm text-gray-600">Người thực hiện</span>
            <span class="text-sm font-semibold text-gray-800">{{ transaction.performed_by }}</span>
          </div>
          <div class="flex justify-between py-2 border-b border-gray-100">
            <span class="text-sm text-gray-600">Vai trò</span>
            <span class="text-sm font-semibold text-gray-800">{{ getRoleLabel(transaction.performed_by_role) }}</span>
          </div>
          <div class="flex justify-between py-2 border-b border-gray-100">
            <span class="text-sm text-gray-600">Thời gian</span>
            <span class="text-sm font-semibold text-gray-800">{{ formatFullDateTime(transaction.timestamp) }}</span>
          </div>
          <div class="flex justify-between py-2 border-b border-gray-100">
            <span class="text-sm text-gray-600">Mã giao dịch</span>
            <button
              @click="copyTransactionId"
              class="text-sm font-mono text-blue-600 hover:text-blue-700 flex items-center gap-1"
            >
              {{ transaction.id.substring(0, 8) }}...
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Metadata (if exists) -->
        <div v-if="transaction.metadata && Object.keys(transaction.metadata).length > 0" class="bg-blue-50 rounded-lg p-4">
          <div class="text-sm font-semibold text-blue-800 mb-2">Thông tin bổ sung</div>
          <div class="space-y-1 text-sm text-blue-700">
            <div v-if="transaction.metadata.cashier_shift_id">
              Ca thu ngân: {{ transaction.metadata.cashier_shift_id.substring(0, 8) }}...
            </div>
            <div v-if="transaction.metadata.variance_amount !== undefined && transaction.metadata.variance_amount !== 0">
              Chênh lệch: {{ formatCurrency(transaction.metadata.variance_amount) }}
            </div>
          </div>
        </div>

        <!-- Close Button -->
        <button
          @click="$emit('close')"
          class="w-full py-3 px-4 bg-gray-200 hover:bg-gray-300 text-gray-700 font-semibold rounded-lg transition-colors"
        >
          Đóng
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue'
import {
  formatCurrency,
  formatFullDateTime,
  getTransactionTypeIcon,
  getTransactionTypeLabel,
  getAmountColor,
  getAmountPrefix,
  getRoleLabel
} from '../../constants/fund'

const props = defineProps({
  transaction: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close'])

// Methods
const copyTransactionId = () => {
  navigator.clipboard.writeText(props.transaction.id)
  alert('Đã sao chép mã giao dịch')
}
</script>

<style scoped>
@keyframes slide-up {
  from {
    transform: translateY(100%);
  }
  to {
    transform: translateY(0);
  }
}

.animate-slide-up {
  animation: slide-up 0.3s ease-out;
}
</style>
