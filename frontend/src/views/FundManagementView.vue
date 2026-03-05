<template>
  <div class="min-h-screen bg-gray-50 pb-24">
    <!-- Header -->
    <div class="bg-gradient-to-r from-orange-500 to-yellow-500 text-white p-4 sticky top-0 z-10 shadow-md">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <button @click="goBack" class="p-2 hover:bg-white/20 rounded-lg transition-colors">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <h1 class="text-xl font-bold">💰 Quản lý quỹ tiền</h1>
        </div>
        <button @click="refreshData" class="p-2 hover:bg-white/20 rounded-lg transition-colors">
          <svg class="w-6 h-6" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Pull to refresh indicator -->
    <div v-if="refreshing" class="text-center py-2 bg-blue-50">
      <span class="text-sm text-blue-600">Đang làm mới...</span>
    </div>

    <!-- Balance Card -->
    <div class="p-4">
      <div class="bg-gradient-to-br from-orange-400 to-yellow-400 rounded-2xl p-6 text-white shadow-lg">
        <div class="text-sm opacity-90 mb-2">Số dư hiện tại</div>
        <div class="space-y-3">
          <div class="flex justify-between items-center">
            <span class="text-base">💵 Tiền mặt</span>
            <span class="text-2xl font-bold">{{ formatCurrency(balance.current_balance?.cash || 0) }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-base">💳 Chuyển khoản</span>
            <span class="text-2xl font-bold">{{ formatCurrency(balance.current_balance?.transfer || 0) }}</span>
          </div>
          <div class="border-t border-white/30 pt-3 mt-3">
            <div class="flex justify-between items-center">
              <span class="text-lg font-semibold">📊 Tổng cộng</span>
              <span class="text-3xl font-bold">{{ formatCurrency(balance.current_balance?.total || 0) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Today Summary -->
    <div class="px-4 pb-4">
      <div class="bg-white rounded-xl p-4 shadow">
        <div class="text-sm font-semibold text-gray-700 mb-3">Tổng quan hôm nay</div>
        <div class="grid grid-cols-2 gap-3 mb-3">
          <div class="bg-green-50 rounded-lg p-3">
            <div class="text-xs text-green-600 mb-1">📥 Thu vào</div>
            <div class="text-lg font-bold text-green-700">
              {{ formatCurrency(balance.today_summary?.total_inflow?.total || 0) }}
            </div>
          </div>
          <div class="bg-red-50 rounded-lg p-3">
            <div class="text-xs text-red-600 mb-1">📤 Chi ra</div>
            <div class="text-lg font-bold text-red-700">
              {{ formatCurrency(balance.today_summary?.total_outflow?.total || 0) }}
            </div>
          </div>
        </div>
        
        <!-- Source breakdown -->
        <div class="border-t border-gray-100 pt-3 mt-3">
          <div class="text-xs font-semibold text-gray-600 mb-2">📊 Phân tích nguồn</div>
          <div class="space-y-2 text-xs">
            <div class="flex justify-between items-center">
              <span class="text-gray-600">💰 Từ doanh thu</span>
              <span class="font-semibold text-blue-600">{{ formatCurrency(todayRevenueFunds) }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-gray-600">💵 Bổ sung thêm</span>
              <span class="font-semibold text-green-600">{{ formatCurrency(todayDeposits) }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-gray-600">📤 Rút ra</span>
              <span class="font-semibold text-red-600">{{ formatCurrency(todayWithdrawals) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Action Buttons -->
    <div class="px-4 pb-4">
      <div class="grid grid-cols-2 gap-3">
        <button
          @click="showDepositModal = true"
          class="bg-green-500 hover:bg-green-600 text-white font-semibold py-3 px-4 rounded-xl shadow-md transition-colors flex items-center justify-center gap-2"
        >
          <span class="text-xl">📥</span>
          <span>Thêm tiền</span>
        </button>
        <button
          @click="showWithdrawModal = true"
          class="bg-red-500 hover:bg-red-600 text-white font-semibold py-3 px-4 rounded-xl shadow-md transition-colors flex items-center justify-center gap-2"
        >
          <span class="text-xl">📤</span>
          <span>Rút tiền</span>
        </button>
      </div>
    </div>

    <!-- Filters -->
    <div class="px-4 pb-4">
      <div class="bg-white rounded-xl p-4 shadow">
        <div class="grid grid-cols-2 gap-3">
          <select
            v-model="filters.type"
            @change="loadTransactions"
            class="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-orange-500 focus:border-transparent"
          >
            <option
              v-for="option in TRANSACTION_TYPE_FILTER_OPTIONS"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </option>
          </select>
          <select
            v-model="filters.money_type"
            @change="loadTransactions"
            class="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-orange-500 focus:border-transparent"
          >
            <option
              v-for="option in MONEY_TYPE_FILTER_OPTIONS"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </option>
          </select>
        </div>
      </div>
    </div>

    <!-- Transaction List -->
    <div class="px-4 pb-4">
      <div class="text-sm font-semibold text-gray-700 mb-3">Lịch sử giao dịch</div>
      
      <!-- Loading skeleton -->
      <div v-if="loadingTransactions" class="space-y-3">
        <div v-for="i in 3" :key="i" class="bg-white rounded-xl p-4 shadow animate-pulse">
          <div class="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
          <div class="h-3 bg-gray-200 rounded w-1/2"></div>
        </div>
      </div>

      <!-- Empty state -->
      <div v-else-if="transactions.length === 0" class="bg-white rounded-xl p-8 shadow text-center">
        <div class="text-4xl mb-2">📋</div>
        <div class="text-gray-500">Chưa có giao dịch nào</div>
      </div>

      <!-- Transaction cards -->
      <div v-else class="space-y-3">
        <div
          v-for="tx in transactions"
          :key="tx.id"
          @click="viewTransactionDetail(tx)"
          class="bg-white rounded-xl p-4 shadow hover:shadow-md transition-shadow cursor-pointer"
        >
          <div class="flex items-start justify-between mb-2">
            <div class="flex items-center gap-2 flex-1">
              <span class="text-2xl">{{ getTransactionTypeIcon(tx.type) }}</span>
              <div class="flex-1">
                <div class="flex items-center gap-2 flex-wrap">
                  <span class="font-semibold text-gray-800">{{ getTransactionTypeLabel(tx.type) }}</span>
                  <span 
                    class="text-xs px-2 py-0.5 rounded-full font-medium"
                    :class="getSourceBadgeClass(tx.type)"
                  >
                    {{ getSourceLabel(tx.type) }}
                  </span>
                </div>
                <div class="text-xs text-gray-500">{{ tx.performed_by }} • {{ tx.performed_by_role }}</div>
              </div>
            </div>
            <div class="text-right">
              <div
                class="text-lg font-bold"
                :class="getAmountColor(tx.type)"
              >
                {{ getAmountPrefix(tx.type) }}{{ formatCurrency(tx.total_amount) }}
              </div>
              <div class="text-xs text-gray-500">{{ formatDateTime(tx.timestamp) }}</div>
            </div>
          </div>
          <div class="text-sm text-gray-600 line-clamp-2">{{ tx.description }}</div>
          
          <!-- Money breakdown -->
          <div v-if="tx.cash_amount > 0 || tx.transfer_amount > 0" class="flex gap-3 mt-2 text-xs">
            <span v-if="tx.cash_amount > 0" class="text-green-600">
              💵 {{ formatCurrency(tx.cash_amount) }}
            </span>
            <span v-if="tx.transfer_amount > 0" class="text-blue-600">
              💳 {{ formatCurrency(tx.transfer_amount) }}
            </span>
          </div>
        </div>
      </div>

      <!-- Load more -->
      <div v-if="hasMore" class="mt-4 text-center">
        <button
          @click="loadMore"
          :disabled="loadingMore"
          class="px-6 py-2 bg-gray-200 hover:bg-gray-300 text-gray-700 rounded-lg transition-colors disabled:opacity-50"
        >
          {{ loadingMore ? 'Đang tải...' : 'Xem thêm' }}
        </button>
      </div>
    </div>

    <!-- Deposit Modal -->
    <DepositModal
      v-if="showDepositModal"
      @close="showDepositModal = false"
      @success="handleDepositSuccess"
    />

    <!-- Withdraw Modal -->
    <WithdrawModal
      v-if="showWithdrawModal"
      :current-balance="balance.current_balance"
      @close="showWithdrawModal = false"
      @success="handleWithdrawSuccess"
    />

    <!-- Transaction Detail Modal -->
    <TransactionDetailModal
      v-if="showDetailModal"
      :transaction="selectedTransaction"
      @close="showDetailModal = false"
    />

    <!-- Bottom Navigation -->
    <BottomNav />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import * as fundService from '../services/fund'
import BottomNav from '../components/BottomNav.vue'
import DepositModal from '../components/fund/DepositModal.vue'
import WithdrawModal from '../components/fund/WithdrawModal.vue'
import TransactionDetailModal from '../components/fund/TransactionDetailModal.vue'
import {
  formatCurrency,
  formatDateTime,
  getTransactionTypeIcon,
  getTransactionTypeLabel,
  getAmountColor,
  getAmountPrefix,
  TRANSACTION_TYPE_FILTER_OPTIONS,
  MONEY_TYPE_FILTER_OPTIONS
} from '../constants/fund'

const router = useRouter()

// State
const loading = ref(false)
const refreshing = ref(false)
const loadingTransactions = ref(false)
const loadingMore = ref(false)
const balance = ref({})
const transactions = ref([])
const filters = ref({
  type: 'all',
  money_type: 'all',
  limit: 20,
  offset: 0
})
const hasMore = ref(false)
const showDepositModal = ref(false)
const showWithdrawModal = ref(false)
const showDetailModal = ref(false)
const selectedTransaction = ref(null)

// Computed - Today breakdown by source
const todayRevenueFunds = computed(() => {
  const today = new Date().toDateString()
  return transactions.value
    .filter(tx => 
      new Date(tx.timestamp).toDateString() === today && 
      (tx.type === 'fund_handover' || tx.type === 'cash_handover')
    )
    .reduce((sum, tx) => sum + tx.total_amount, 0)
})

const todayDeposits = computed(() => {
  const today = new Date().toDateString()
  return transactions.value
    .filter(tx => 
      new Date(tx.timestamp).toDateString() === today && 
      tx.type === 'deposit'
    )
    .reduce((sum, tx) => sum + tx.total_amount, 0)
})

const todayWithdrawals = computed(() => {
  const today = new Date().toDateString()
  return transactions.value
    .filter(tx => 
      new Date(tx.timestamp).toDateString() === today && 
      tx.type === 'withdrawal'
    )
    .reduce((sum, tx) => sum + tx.total_amount, 0)
})

// Methods
const goBack = () => {
  router.push('/dashboard')
}

const refreshData = async () => {
  loading.value = true
  try {
    await Promise.all([
      loadBalance(),
      loadTransactions()
    ])
  } finally {
    loading.value = false
  }
}

const loadBalance = async () => {
  try {
    balance.value = await fundService.getBalance()
  } catch (error) {
    console.error('Failed to load balance:', error)
    alert('Không thể tải số dư quỹ')
  }
}

const loadTransactions = async () => {
  loadingTransactions.value = true
  filters.value.offset = 0
  try {
    const response = await fundService.getTransactions(filters.value)
    transactions.value = response.transactions || []
    hasMore.value = response.total > filters.value.limit
  } catch (error) {
    console.error('Failed to load transactions:', error)
    alert('Không thể tải lịch sử giao dịch')
  } finally {
    loadingTransactions.value = false
  }
}

const loadMore = async () => {
  loadingMore.value = true
  filters.value.offset += filters.value.limit
  try {
    const response = await fundService.getTransactions(filters.value)
    transactions.value.push(...(response.transactions || []))
    hasMore.value = filters.value.offset + filters.value.limit < response.total
  } catch (error) {
    console.error('Failed to load more transactions:', error)
  } finally {
    loadingMore.value = false
  }
}

const handleDepositSuccess = () => {
  showDepositModal.value = false
  refreshData()
}

const handleWithdrawSuccess = () => {
  showWithdrawModal.value = false
  refreshData()
}

const viewTransactionDetail = (transaction) => {
  selectedTransaction.value = transaction
  showDetailModal.value = true
}

// Helper functions for source classification
const getSourceLabel = (type) => {
  const labels = {
    deposit: 'Bổ sung',
    withdrawal: 'Rút ra',
    fund_handover: 'Doanh thu',
    cash_handover: 'Doanh thu'
  }
  return labels[type] || 'Khác'
}

const getSourceBadgeClass = (type) => {
  const classes = {
    deposit: 'bg-green-100 text-green-700',
    withdrawal: 'bg-red-100 text-red-700',
    fund_handover: 'bg-blue-100 text-blue-700',
    cash_handover: 'bg-blue-100 text-blue-700'
  }
  return classes[type] || 'bg-gray-100 text-gray-700'
}

// Lifecycle
onMounted(() => {
  refreshData()
})
</script>
