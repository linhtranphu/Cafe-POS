<template>
  <div class="fixed inset-0 bg-black/50 flex items-end sm:items-center justify-center z-50" @click.self="$emit('close')">
    <div class="bg-white rounded-t-3xl sm:rounded-2xl w-full sm:max-w-md max-h-[90vh] overflow-y-auto animate-slide-up">
      <!-- Header -->
      <div class="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
        <h2 class="text-xl font-bold text-gray-800">📤 Rút tiền từ quỹ</h2>
        <button @click="$emit('close')" class="p-2 hover:bg-gray-100 rounded-lg transition-colors">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Form -->
      <form @submit.prevent="handleSubmit" class="p-6 space-y-4">
        <!-- Fund Type Selector -->
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-2">Rút từ quỹ</label>
          <select
            v-model="form.fund_type"
            @change="loadFundBalance"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-red-500 focus:border-transparent"
          >
            <option v-for="key in fundTypeKeys" :key="key" :value="key">
              {{ FUND_TYPE_ICONS[key] }} {{ FUND_TYPE_LABELS[key] }}
            </option>
          </select>
        </div>

        <!-- Fund Balance Display -->
        <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
          <div class="text-sm font-semibold text-yellow-800 mb-2">Số dư quỹ được chọn:</div>
          <div v-if="loadingFundBalance" class="animate-pulse space-y-1">
            <div class="h-4 bg-yellow-200 rounded w-40"></div>
            <div class="h-4 bg-yellow-200 rounded w-40"></div>
          </div>
          <div v-else class="space-y-1 text-sm text-yellow-700">
            <div>💵 Tiền mặt: {{ formatCurrency(fundBalance?.cash || 0) }}</div>
            <div>💳 Chuyển khoản: {{ formatCurrency(fundBalance?.transfer || 0) }}</div>
          </div>
        </div>

        <!-- Money Type Selector -->
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-2">Loại tiền</label>
          <div class="grid grid-cols-3 gap-2">
            <button
              type="button"
              @click="form.money_type = MONEY_TYPES.CASH"
              :class="[
                'py-2 px-3 rounded-lg font-medium transition-colors',
                form.money_type === MONEY_TYPES.CASH
                  ? 'bg-green-500 text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              ]"
            >
              💵 Tiền mặt
            </button>
            <button
              type="button"
              @click="form.money_type = MONEY_TYPES.TRANSFER"
              :class="[
                'py-2 px-3 rounded-lg font-medium transition-colors',
                form.money_type === MONEY_TYPES.TRANSFER
                  ? 'bg-blue-500 text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              ]"
            >
              💳 CK
            </button>
            <button
              type="button"
              @click="form.money_type = MONEY_TYPES.BOTH"
              :class="[
                'py-2 px-3 rounded-lg font-medium transition-colors',
                form.money_type === MONEY_TYPES.BOTH
                  ? 'bg-orange-500 text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              ]"
            >
              💰 Cả hai
            </button>
          </div>
        </div>

        <!-- Cash Amount -->
        <div v-if="form.money_type === MONEY_TYPES.CASH || form.money_type === MONEY_TYPES.BOTH">
          <label class="block text-sm font-semibold text-gray-700 mb-2">
            💵 Số tiền mặt
          </label>
          <input
            v-model.number="form.cash_amount"
            type="number"
            min="0"
            :max="currentBalance?.cash || 0"
            step="1000"
            placeholder="0"
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent text-lg"
            required
          />
          <div v-if="form.cash_amount > (fundBalance?.cash || 0)" class="text-xs text-red-500 mt-1">
            Số dư không đủ!
          </div>
        </div>

        <!-- Transfer Amount -->
        <div v-if="form.money_type === MONEY_TYPES.TRANSFER || form.money_type === MONEY_TYPES.BOTH">
          <label class="block text-sm font-semibold text-gray-700 mb-2">
            💳 Số tiền chuyển khoản
          </label>
          <input
            v-model.number="form.transfer_amount"
            type="number"
            min="0"
            :max="currentBalance?.transfer || 0"
            step="1000"
            placeholder="0"
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent text-lg"
            required
          />
          <div v-if="form.transfer_amount > (fundBalance?.transfer || 0)" class="text-xs text-red-500 mt-1">
            Số dư không đủ!
          </div>
        </div>

        <!-- Total Display -->
        <div v-if="totalAmount > 0" class="bg-red-50 rounded-lg p-4">
          <div class="text-sm text-red-600 mb-1">Tổng rút</div>
          <div class="text-2xl font-bold text-red-700">{{ formatCurrency(totalAmount) }}</div>
        </div>

        <!-- Reason -->
        <div>
          <label class="block text-sm font-semibold text-gray-700 mb-2">
            Lý do <span class="text-red-500">*</span>
          </label>
          <textarea
            v-model="form.reason"
            rows="3"
            :placeholder="`Nhập lý do rút tiền (tối thiểu ${VALIDATION.MIN_REASON_LENGTH} ký tự)`"
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent resize-none"
            required
            :minlength="VALIDATION.MIN_REASON_LENGTH"
          ></textarea>
          <div class="text-xs text-gray-500 mt-1">
            {{ form.reason.length }}/{{ VALIDATION.MIN_REASON_LENGTH }} ký tự
          </div>
        </div>

        <!-- Error Message -->
        <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-3">
          <div class="text-sm text-red-600">{{ error }}</div>
        </div>

        <!-- Actions -->
        <div class="flex gap-3 pt-2">
          <button
            type="button"
            @click="$emit('close')"
            class="flex-1 py-3 px-4 bg-gray-200 hover:bg-gray-300 text-gray-700 font-semibold rounded-lg transition-colors"
          >
            Hủy
          </button>
          <button
            type="submit"
            :disabled="submitting || !isValid"
            class="flex-1 py-3 px-4 bg-red-500 hover:bg-red-600 text-white font-semibold rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ submitting ? 'Đang xử lý...' : 'Xác nhận' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import * as fundService from '../../services/fund'
import { formatCurrency, MONEY_TYPES, VALIDATION, FUND_TYPES, FUND_TYPE_LABELS, FUND_TYPE_ICONS } from '../../constants/fund'

const props = defineProps({
  currentBalance: {
    type: Object,
    default: () => ({ cash: 0, transfer: 0, total: 0 })
  }
})

const emit = defineEmits(['close', 'success'])

const fundTypeKeys = [FUND_TYPES.OPERATING, FUND_TYPES.INVENTORY, FUND_TYPES.PROFIT, FUND_TYPES.CASH_DRAWER]

// State
const form = ref({
  fund_type: FUND_TYPES.OPERATING,
  money_type: MONEY_TYPES.CASH,
  cash_amount: 0,
  transfer_amount: 0,
  reason: ''
})
const submitting = ref(false)
const error = ref('')
const fundBalance = ref({ cash: 0, transfer: 0, total: 0 })
const loadingFundBalance = ref(false)

const loadFundBalance = async () => {
  loadingFundBalance.value = true
  try {
    const data = await fundService.getAllBalances()
    fundBalance.value = data?.[form.value.fund_type] || { cash: 0, transfer: 0, total: 0 }
  } catch (e) {
    console.error('Failed to load fund balance:', e)
  } finally {
    loadingFundBalance.value = false
  }
}

// Computed
const totalAmount = computed(() => {
  let total = 0
  if (form.value.money_type === MONEY_TYPES.CASH || form.value.money_type === MONEY_TYPES.BOTH) {
    total += form.value.cash_amount || 0
  }
  if (form.value.money_type === MONEY_TYPES.TRANSFER || form.value.money_type === MONEY_TYPES.BOTH) {
    total += form.value.transfer_amount || 0
  }
  return total
})

const isValid = computed(() => {
  if (totalAmount.value <= 0 || form.value.reason.length < VALIDATION.MIN_REASON_LENGTH) return false

  if (form.value.money_type === MONEY_TYPES.CASH || form.value.money_type === MONEY_TYPES.BOTH) {
    if (form.value.cash_amount > (fundBalance.value?.cash || 0)) return false
  }
  if (form.value.money_type === MONEY_TYPES.TRANSFER || form.value.money_type === MONEY_TYPES.BOTH) {
    if (form.value.transfer_amount > (fundBalance.value?.transfer || 0)) return false
  }

  return true
})

onMounted(() => {
  loadFundBalance()
})

// Methods
const handleSubmit = async () => {
  if (!isValid.value) return

  submitting.value = true
  error.value = ''

  try {
    const data = {
      fund_type: form.value.fund_type,
      cash_amount: form.value.money_type === MONEY_TYPES.CASH || form.value.money_type === MONEY_TYPES.BOTH ? form.value.cash_amount : 0,
      transfer_amount: form.value.money_type === MONEY_TYPES.TRANSFER || form.value.money_type === MONEY_TYPES.BOTH ? form.value.transfer_amount : 0,
      reason: form.value.reason
    }

    await fundService.createWithdrawal(data)
    emit('success')
  } catch (err) {
    console.error('Withdrawal failed:', err)
    error.value = err.response?.data?.error || 'Không thể rút tiền. Vui lòng thử lại.'
  } finally {
    submitting.value = false
  }
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
