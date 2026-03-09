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

    <!-- Total Balance Summary Card -->
    <div class="p-4">
      <div class="bg-gradient-to-br from-orange-400 to-yellow-400 rounded-2xl p-5 text-white shadow-lg">
        <div class="text-sm opacity-90 mb-1">Tổng tất cả quỹ</div>
        <div v-if="loadingAllBalances" class="animate-pulse">
          <div class="h-8 bg-white/30 rounded w-48 mb-2"></div>
          <div class="flex gap-4">
            <div class="h-5 bg-white/20 rounded w-28"></div>
            <div class="h-5 bg-white/20 rounded w-28"></div>
          </div>
        </div>
        <div v-else>
          <div class="text-3xl font-bold mb-2">{{ formatCurrency(allBalances?.total?.total || 0) }}</div>
          <div class="flex gap-4 text-sm opacity-90">
            <span>💵 TM: {{ formatCurrency(allBalances?.total?.cash || 0) }}</span>
            <span>💳 CK: {{ formatCurrency(allBalances?.total?.transfer || 0) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 4 Fund Balance Cards -->
    <div class="px-4 pb-4">
      <div class="grid grid-cols-2 gap-3">
        <template v-if="loadingAllBalances">
          <div v-for="i in 4" :key="i" class="bg-white rounded-xl p-4 shadow animate-pulse">
            <div class="h-4 bg-gray-200 rounded w-3/4 mb-3"></div>
            <div class="h-6 bg-gray-200 rounded w-1/2 mb-2"></div>
            <div class="h-3 bg-gray-100 rounded w-full"></div>
          </div>
        </template>
        <template v-else>
          <div
            v-for="key in fundTypeKeys"
            :key="key"
            class="bg-white rounded-xl p-4 shadow border-l-4 cursor-pointer hover:shadow-md transition-shadow"
            :class="`border-${FUND_TYPE_COLORS[key]}-400`"
            @click="filterByFund(key)"
          >
            <div class="flex items-center gap-2 mb-2">
              <span class="text-lg">{{ FUND_TYPE_ICONS[key] }}</span>
              <span class="text-xs font-semibold text-gray-600">{{ FUND_TYPE_LABELS[key] }}</span>
            </div>
            <div class="text-lg font-bold text-gray-800">
              {{ formatCurrency(allBalances?.[key]?.total || 0) }}
            </div>
            <div class="flex gap-2 mt-1 text-xs text-gray-500">
              <span>TM: {{ formatCurrency(allBalances?.[key]?.cash || 0) }}</span>
              <span>CK: {{ formatCurrency(allBalances?.[key]?.transfer || 0) }}</span>
            </div>
          </div>
        </template>
      </div>
    </div>

    <!-- Action Buttons -->
    <div class="px-4 pb-4">
      <div class="grid grid-cols-3 gap-3">
        <button
          @click="showDepositModal = true"
          class="bg-green-500 hover:bg-green-600 text-white font-semibold py-3 px-3 rounded-xl shadow-md transition-colors flex flex-col items-center justify-center gap-1"
        >
          <span class="text-xl">📥</span>
          <span class="text-xs">Thêm tiền</span>
        </button>
        <button
          @click="showWithdrawModal = true"
          class="bg-red-500 hover:bg-red-600 text-white font-semibold py-3 px-3 rounded-xl shadow-md transition-colors flex flex-col items-center justify-center gap-1"
        >
          <span class="text-xl">📤</span>
          <span class="text-xs">Rút tiền</span>
        </button>
        <button
          @click="showTransferModal = true"
          class="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-3 px-3 rounded-xl shadow-md transition-colors flex flex-col items-center justify-center gap-1"
        >
          <span class="text-xl">↔️</span>
          <span class="text-xs">Chuyển quỹ</span>
        </button>
      </div>
    </div>

    <!-- Filters -->
    <div class="px-4 pb-4">
      <div class="bg-white rounded-xl p-4 shadow space-y-3">
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
        <div class="grid grid-cols-2 gap-3">
          <select
            v-model="filters.fund_type"
            @change="loadTransactions"
            class="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-orange-500 focus:border-transparent"
          >
            <option
              v-for="option in FUND_TYPE_FILTER_OPTIONS"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </option>
          </select>
          <select
            v-model="filters.source_type"
            @change="loadTransactions"
            class="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-orange-500 focus:border-transparent"
          >
            <option
              v-for="option in SOURCE_TYPE_FILTER_OPTIONS"
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
                <div class="flex items-center gap-1.5 flex-wrap">
                  <span class="font-semibold text-gray-800">{{ getTransactionTypeLabel(tx.type) }}</span>
                  <!-- Fund type badge -->
                  <span
                    v-if="tx.fund_type"
                    class="text-xs px-1.5 py-0.5 rounded-full font-medium"
                    :class="`bg-${FUND_TYPE_COLORS[tx.fund_type] || 'gray'}-100 text-${FUND_TYPE_COLORS[tx.fund_type] || 'gray'}-700`"
                  >
                    {{ FUND_TYPE_ICONS[tx.fund_type] }} {{ FUND_TYPE_LABELS[tx.fund_type] || tx.fund_type }}
                  </span>
                  <!-- Source type badge -->
                  <span
                    v-if="tx.source_type"
                    class="text-xs px-1.5 py-0.5 rounded-full bg-gray-100 text-gray-600 font-medium"
                  >
                    {{ SOURCE_TYPE_ICONS[tx.source_type] }} {{ SOURCE_TYPE_LABELS[tx.source_type] || tx.source_type }}
                  </span>
                </div>
                <div class="text-xs text-gray-500">{{ tx.performed_by_name }} • {{ tx.performed_by_role }}</div>
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
          <div class="text-sm text-gray-600 line-clamp-2">{{ tx.description || tx.reason }}</div>

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

    <!-- Transfer Modal -->
    <div
      v-if="showTransferModal"
      class="fixed inset-0 bg-black/50 z-50 flex items-end justify-center"
      @click.self="showTransferModal = false"
    >
      <div class="bg-white w-full max-w-lg rounded-t-2xl p-6 slide-up">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-bold text-gray-800">↔️ Chuyển tiền giữa quỹ</h2>
          <button @click="showTransferModal = false" class="p-2 hover:bg-gray-100 rounded-lg">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="space-y-4">
          <!-- From fund -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Từ quỹ</label>
            <select
              v-model="transferForm.from_fund_type"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500"
            >
              <option v-for="key in fundTypeKeys" :key="key" :value="key">
                {{ FUND_TYPE_ICONS[key] }} {{ FUND_TYPE_LABELS[key] }}
                ({{ formatCurrency(allBalances?.[key]?.total || 0) }})
              </option>
            </select>
          </div>

          <!-- To fund -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Đến quỹ</label>
            <select
              v-model="transferForm.to_fund_type"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500"
            >
              <option
                v-for="key in fundTypeKeys"
                :key="key"
                :value="key"
                :disabled="key === transferForm.from_fund_type"
              >
                {{ FUND_TYPE_ICONS[key] }} {{ FUND_TYPE_LABELS[key] }}
              </option>
            </select>
          </div>

          <!-- Money type -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Loại tiền</label>
            <div class="grid grid-cols-2 gap-2">
              <button
                @click="transferForm.money_type = 'cash'"
                :class="transferForm.money_type === 'cash' ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-700'"
                class="py-2 rounded-lg text-sm font-medium transition-colors"
              >
                💵 Tiền mặt
              </button>
              <button
                @click="transferForm.money_type = 'transfer'"
                :class="transferForm.money_type === 'transfer' ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-700'"
                class="py-2 rounded-lg text-sm font-medium transition-colors"
              >
                💳 Chuyển khoản
              </button>
            </div>
          </div>

          <!-- Amount -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Số tiền</label>
            <input
              v-model.number="transferForm.amount"
              type="number"
              min="0"
              placeholder="Nhập số tiền..."
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500"
            />
            <!-- Source fund balance hint -->
            <div class="text-xs text-gray-500 mt-1">
              Số dư quỹ nguồn: {{ formatCurrency(sourceBalance) }}
            </div>
          </div>

          <!-- Reason -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Lý do (tối thiểu 10 ký tự)</label>
            <textarea
              v-model="transferForm.reason"
              rows="2"
              placeholder="Nhập lý do chuyển quỹ..."
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 resize-none"
            ></textarea>
          </div>

          <!-- Submit -->
          <button
            @click="submitTransfer"
            :disabled="transferring || !isTransferValid"
            class="w-full py-3 bg-blue-500 hover:bg-blue-600 disabled:bg-gray-300 text-white font-semibold rounded-xl transition-colors"
          >
            {{ transferring ? 'Đang xử lý...' : 'Xác nhận chuyển quỹ' }}
          </button>
        </div>
      </div>
    </div>

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
  MONEY_TYPE_FILTER_OPTIONS,
  FUND_TYPES,
  FUND_TYPE_LABELS,
  FUND_TYPE_ICONS,
  FUND_TYPE_COLORS,
  FUND_TYPE_FILTER_OPTIONS,
  SOURCE_TYPE_LABELS,
  SOURCE_TYPE_ICONS,
  SOURCE_TYPE_FILTER_OPTIONS
} from '../constants/fund'

const router = useRouter()

// State
const loading = ref(false)
const refreshing = ref(false)
const loadingTransactions = ref(false)
const loadingMore = ref(false)
const loadingAllBalances = ref(false)
const balance = ref({})
const allBalances = ref(null)
const transactions = ref([])
const filters = ref({
  type: 'all',
  money_type: 'all',
  fund_type: 'all',
  source_type: 'all',
  limit: 20,
  offset: 0
})
const hasMore = ref(false)
const showDepositModal = ref(false)
const showWithdrawModal = ref(false)
const showTransferModal = ref(false)
const showDetailModal = ref(false)
const selectedTransaction = ref(null)
const transferring = ref(false)
const transferForm = ref({
  from_fund_type: FUND_TYPES.OPERATING,
  to_fund_type: FUND_TYPES.INVENTORY,
  money_type: 'cash',
  amount: null,
  reason: ''
})

// Computed
const fundTypeKeys = computed(() => [
  FUND_TYPES.OPERATING,
  FUND_TYPES.INVENTORY,
  FUND_TYPES.PROFIT,
  FUND_TYPES.CASH_DRAWER
])

const sourceBalance = computed(() => {
  if (!allBalances.value || !transferForm.value.from_fund_type) return 0
  const b = allBalances.value[transferForm.value.from_fund_type]
  if (!b) return 0
  return transferForm.value.money_type === 'cash' ? b.cash : b.transfer
})

const isTransferValid = computed(() => {
  const f = transferForm.value
  return f.from_fund_type &&
    f.to_fund_type &&
    f.from_fund_type !== f.to_fund_type &&
    f.amount > 0 &&
    f.amount <= sourceBalance.value &&
    f.reason.length >= 10
})

// Methods
const goBack = () => {
  router.push('/dashboard')
}

const filterByFund = (fundType) => {
  filters.value.fund_type = fundType
  loadTransactions()
}

const refreshData = async () => {
  loading.value = true
  try {
    await Promise.all([
      loadBalance(),
      loadAllBalances(),
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
  }
}

const loadAllBalances = async () => {
  loadingAllBalances.value = true
  try {
    const data = await fundService.getAllBalances()
    allBalances.value = data
  } catch (error) {
    console.error('Failed to load all balances:', error)
  } finally {
    loadingAllBalances.value = false
  }
}

const loadTransactions = async () => {
  loadingTransactions.value = true
  filters.value.offset = 0
  try {
    const params = { ...filters.value }
    if (params.fund_type === 'all') delete params.fund_type
    if (params.source_type === 'all') delete params.source_type
    const response = await fundService.getTransactions(params)
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
    const params = { ...filters.value }
    if (params.fund_type === 'all') delete params.fund_type
    if (params.source_type === 'all') delete params.source_type
    const response = await fundService.getTransactions(params)
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

const submitTransfer = async () => {
  if (!isTransferValid.value) return
  transferring.value = true
  try {
    const f = transferForm.value
    const payload = {
      from_fund_type: f.from_fund_type,
      to_fund_type: f.to_fund_type,
      cash_amount: f.money_type === 'cash' ? f.amount : 0,
      transfer_amount: f.money_type === 'transfer' ? f.amount : 0,
      reason: f.reason
    }
    await fundService.transferBetweenFunds(payload)
    showTransferModal.value = false
    transferForm.value = {
      from_fund_type: FUND_TYPES.OPERATING,
      to_fund_type: FUND_TYPES.INVENTORY,
      money_type: 'cash',
      amount: null,
      reason: ''
    }
    await refreshData()
  } catch (error) {
    console.error('Transfer failed:', error)
    alert('Chuyển quỹ thất bại: ' + (error.response?.data?.error || error.message))
  } finally {
    transferring.value = false
  }
}

const viewTransactionDetail = (transaction) => {
  selectedTransaction.value = transaction
  showDetailModal.value = true
}

// Lifecycle
onMounted(() => {
  refreshData()
})
</script>

<style scoped>
.slide-up {
  animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
  from {
    transform: translateY(100%);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}
</style>
