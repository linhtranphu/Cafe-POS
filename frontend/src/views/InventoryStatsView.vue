<template>
  <div class="min-h-screen bg-gray-50 pb-20">
    <!-- Header -->
    <div class="bg-white shadow-sm sticky top-0 z-10">
      <div class="max-w-7xl mx-auto px-4 py-4 flex items-center gap-3">
        <router-link to="/dashboard" class="text-gray-500 hover:text-gray-700">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </router-link>
        <div>
          <h1 class="text-xl font-bold text-gray-800">Thống kê tồn kho</h1>
          <p class="text-sm text-gray-500">Khối lượng bán & tồn kho nguyên liệu</p>
        </div>
      </div>

      <!-- Tab switcher -->
      <div class="max-w-7xl mx-auto px-4 flex border-t border-gray-100">
        <button
          @click="activeTab = 'stats'"
          :class="[
            'px-5 py-3 text-sm font-medium border-b-2 transition-colors',
            activeTab === 'stats' ? 'border-blue-600 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700'
          ]"
        >
          Thống kê
        </button>
        <button
          @click="activeTab = 'stocktake'"
          :class="[
            'px-5 py-3 text-sm font-medium border-b-2 transition-colors',
            activeTab === 'stocktake' ? 'border-blue-600 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700'
          ]"
        >
          Kiểm kê
        </button>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 py-4 space-y-4">

      <!-- ======== STATS TAB ======== -->
      <template v-if="activeTab === 'stats'">
        <!-- Date range filter -->
        <div class="bg-white rounded-2xl p-4 shadow-sm">
          <div class="flex flex-wrap gap-2 mb-3">
            <button
              v-for="preset in datePresets"
              :key="preset.key"
              @click="applyPreset(preset)"
              :class="[
                'px-3 py-1.5 rounded-lg text-sm font-medium transition-colors',
                activePreset === preset.key
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              ]"
            >
              {{ preset.label }}
            </button>
          </div>
          <div class="flex gap-3 items-end">
            <div class="flex-1">
              <label class="block text-xs text-gray-500 mb-1">Từ ngày</label>
              <input
                v-model="fromDate"
                type="date"
                class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-500"
                @change="activePreset = 'custom'; loadData()"
              />
            </div>
            <div class="flex-1">
              <label class="block text-xs text-gray-500 mb-1">Đến ngày</label>
              <input
                v-model="toDate"
                type="date"
                class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-500"
                @change="activePreset = 'custom'; loadData()"
              />
            </div>
          </div>
        </div>

        <!-- Summary cards -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
          <div class="bg-white rounded-xl p-4 shadow-sm text-center">
            <div class="text-2xl font-bold text-gray-800">{{ items.length }}</div>
            <div class="text-xs text-gray-500 mt-1">Tổng nguyên liệu</div>
          </div>
          <div class="bg-white rounded-xl p-4 shadow-sm text-center">
            <div class="text-2xl font-bold text-green-600">{{ okCount }}</div>
            <div class="text-xs text-gray-500 mt-1">Đủ hàng</div>
          </div>
          <div class="bg-white rounded-xl p-4 shadow-sm text-center">
            <div class="text-2xl font-bold text-yellow-500">{{ lowCount }}</div>
            <div class="text-xs text-gray-500 mt-1">Sắp hết</div>
          </div>
          <div class="bg-white rounded-xl p-4 shadow-sm text-center">
            <div class="text-2xl font-bold text-red-500">{{ outCount }}</div>
            <div class="text-xs text-gray-500 mt-1">Hết hàng</div>
          </div>
        </div>

        <!-- Filter & search -->
        <div class="bg-white rounded-2xl p-4 shadow-sm flex flex-wrap gap-3">
          <input
            v-model="search"
            type="text"
            placeholder="Tìm nguyên liệu..."
            class="flex-1 min-w-[160px] border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-500"
          />
          <select
            v-model="filterStatus"
            class="border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-500"
          >
            <option value="">Tất cả trạng thái</option>
            <option value="out">Hết hàng</option>
            <option value="low">Sắp hết</option>
            <option value="ok">Đủ hàng</option>
          </select>
          <select
            v-model="filterCategory"
            class="border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-500"
          >
            <option value="">Tất cả danh mục</option>
            <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
          </select>
        </div>

        <!-- Loading -->
        <div v-if="loading" class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-10 w-10 border-4 border-blue-500 border-t-transparent"></div>
        </div>

        <!-- Table (desktop) -->
        <div v-else class="hidden md:block bg-white rounded-2xl shadow-sm overflow-hidden">
          <table class="w-full text-sm">
            <thead class="bg-gray-50 border-b border-gray-200">
              <tr>
                <th
                  v-for="col in columns"
                  :key="col.key"
                  class="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wide cursor-pointer hover:bg-gray-100 select-none"
                  @click="sortBy(col.key)"
                >
                  {{ col.label }}
                  <span v-if="sortKey === col.key" class="ml-1">{{ sortDir === 'asc' ? '↑' : '↓' }}</span>
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-if="filteredItems.length === 0">
                <td :colspan="columns.length" class="px-4 py-8 text-center text-gray-400">Không có dữ liệu</td>
              </tr>
              <tr
                v-for="item in filteredItems"
                :key="item.id"
                class="hover:bg-gray-50 transition-colors"
              >
                <td class="px-4 py-3 font-medium text-gray-800">
                  {{ item.name }}
                  <div class="text-xs text-gray-400">{{ item.category || '—' }}</div>
                </td>
                <td class="px-4 py-3 text-center">
                  <span :class="statusBadgeClass(item.status)">{{ statusLabel(item.status) }}</span>
                </td>
                <td class="px-4 py-3 text-right">
                  <span :class="item.status === 'out' ? 'text-red-600 font-semibold' : item.status === 'low' ? 'text-yellow-600 font-semibold' : 'text-gray-800'">
                    {{ fmt(item.current_qty) }}
                  </span>
                  <span class="text-gray-400 ml-1">{{ item.unit }}</span>
                </td>
                <td class="px-4 py-3 text-right text-gray-500">{{ fmt(item.min_stock) }} {{ item.unit }}</td>
                <td class="px-4 py-3 text-right text-blue-600 font-medium">
                  {{ fmt(item.consumed) }} {{ item.unit }}
                  <div v-if="item.order_count > 0" class="text-xs text-gray-400">{{ item.order_count }} đơn</div>
                </td>
                <td class="px-4 py-3 text-right text-green-600">{{ fmt(item.restocked) }} {{ item.unit }}</td>
                <td class="px-4 py-3 text-right text-orange-500">{{ fmt(item.wasted) }} {{ item.unit }}</td>
                <td class="px-4 py-3 text-right" :class="item.adjusted >= 0 ? 'text-gray-500' : 'text-red-500'">
                  {{ item.adjusted >= 0 ? '+' : '' }}{{ fmt(item.adjusted) }} {{ item.unit }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Card list (mobile) -->
        <div v-if="!loading" class="md:hidden space-y-3">
          <div v-if="filteredItems.length === 0" class="bg-white rounded-2xl p-6 text-center text-gray-400 shadow-sm">
            Không có dữ liệu
          </div>
          <div
            v-for="item in filteredItems"
            :key="item.id"
            class="bg-white rounded-2xl p-4 shadow-sm"
          >
            <div class="flex items-start justify-between mb-3">
              <div>
                <div class="font-semibold text-gray-800">{{ item.name }}</div>
                <div class="text-xs text-gray-400">{{ item.category || '—' }}</div>
              </div>
              <span :class="statusBadgeClass(item.status)">{{ statusLabel(item.status) }}</span>
            </div>

            <div class="grid grid-cols-2 gap-2 text-sm">
              <div class="bg-gray-50 rounded-lg p-2">
                <div class="text-xs text-gray-500">Tồn kho</div>
                <div :class="['font-semibold', item.status === 'out' ? 'text-red-600' : item.status === 'low' ? 'text-yellow-600' : 'text-gray-800']">
                  {{ fmt(item.current_qty) }} {{ item.unit }}
                </div>
                <div class="text-xs text-gray-400">Tối thiểu: {{ fmt(item.min_stock) }}</div>
              </div>
              <div class="bg-blue-50 rounded-lg p-2">
                <div class="text-xs text-blue-500">Đã bán (kỳ này)</div>
                <div class="font-semibold text-blue-700">{{ fmt(item.consumed) }} {{ item.unit }}</div>
                <div v-if="item.order_count > 0" class="text-xs text-blue-400">{{ item.order_count }} đơn hàng</div>
              </div>
              <div v-if="item.restocked > 0" class="bg-green-50 rounded-lg p-2">
                <div class="text-xs text-green-500">Nhập kho</div>
                <div class="font-semibold text-green-700">+{{ fmt(item.restocked) }} {{ item.unit }}</div>
              </div>
              <div v-if="item.wasted > 0" class="bg-orange-50 rounded-lg p-2">
                <div class="text-xs text-orange-500">Hao hụt</div>
                <div class="font-semibold text-orange-700">{{ fmt(item.wasted) }} {{ item.unit }}</div>
              </div>
              <div v-if="item.adjusted !== 0" class="bg-gray-50 rounded-lg p-2">
                <div class="text-xs text-gray-500">Điều chỉnh</div>
                <div :class="['font-semibold', item.adjusted >= 0 ? 'text-gray-700' : 'text-red-600']">
                  {{ item.adjusted >= 0 ? '+' : '' }}{{ fmt(item.adjusted) }} {{ item.unit }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>

      <!-- ======== STOCKTAKE TAB ======== -->
      <template v-else>
        <!-- Info banner -->
        <div class="bg-blue-50 border border-blue-200 rounded-xl px-4 py-3 text-sm text-blue-700">
          Nhập số lượng thực tế để so sánh với hệ thống. Chỉ những mục có chênh lệch mới được điều chỉnh khi bấm xác nhận.
        </div>

        <!-- Loading -->
        <div v-if="loading" class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-10 w-10 border-4 border-blue-500 border-t-transparent"></div>
        </div>

        <template v-else>
          <!-- Search filter for stocktake -->
          <div class="bg-white rounded-2xl p-4 shadow-sm flex gap-3">
            <input
              v-model="stocktakeSearch"
              type="text"
              placeholder="Tìm nguyên liệu..."
              class="flex-1 border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-500"
            />
            <select
              v-model="stocktakeCategory"
              class="border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-500"
            >
              <option value="">Tất cả danh mục</option>
              <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
            </select>
          </div>

          <!-- Stocktake table (desktop) -->
          <div class="hidden md:block bg-white rounded-2xl shadow-sm overflow-hidden">
            <table class="w-full text-sm">
              <thead class="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wide">Nguyên liệu</th>
                  <th class="px-4 py-3 text-right text-xs font-semibold text-gray-500 uppercase tracking-wide">Hệ thống</th>
                  <th class="px-4 py-3 text-right text-xs font-semibold text-gray-500 uppercase tracking-wide w-40">Thực tế</th>
                  <th class="px-4 py-3 text-right text-xs font-semibold text-gray-500 uppercase tracking-wide">Chênh lệch</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100">
                <tr v-if="filteredStocktakeItems.length === 0">
                  <td colspan="4" class="px-4 py-8 text-center text-gray-400">Không có dữ liệu</td>
                </tr>
                <tr
                  v-for="item in filteredStocktakeItems"
                  :key="item.id"
                  :class="['transition-colors', physicalCounts[item.id] !== undefined ? 'bg-yellow-50' : 'hover:bg-gray-50']"
                >
                  <td class="px-4 py-3">
                    <div class="font-medium text-gray-800">{{ item.name }}</div>
                    <div class="text-xs text-gray-400">{{ item.category || '—' }}</div>
                  </td>
                  <td class="px-4 py-3 text-right text-gray-700">
                    {{ fmt(item.current_qty) }} <span class="text-gray-400">{{ item.unit }}</span>
                  </td>
                  <td class="px-4 py-3 text-right">
                    <div class="flex items-center justify-end gap-1">
                      <input
                        :value="physicalCounts[item.id] ?? ''"
                        @input="setPhysicalCount(item.id, $event.target.value)"
                        type="number"
                        min="0"
                        step="any"
                        :placeholder="fmt(item.current_qty)"
                        class="w-28 border border-gray-300 rounded-lg px-2 py-1 text-right text-sm focus:outline-none focus:border-blue-500"
                      />
                      <span class="text-gray-400 text-xs">{{ item.unit }}</span>
                    </div>
                  </td>
                  <td class="px-4 py-3 text-right">
                    <span :class="diffClass(item)">{{ diffLabel(item) }}</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Stocktake card list (mobile) -->
          <div class="md:hidden space-y-2">
            <div v-if="filteredStocktakeItems.length === 0" class="bg-white rounded-2xl p-6 text-center text-gray-400 shadow-sm">
              Không có dữ liệu
            </div>
            <div
              v-for="item in filteredStocktakeItems"
              :key="item.id"
              :class="['rounded-2xl p-4 shadow-sm', physicalCounts[item.id] !== undefined ? 'bg-yellow-50 border border-yellow-200' : 'bg-white']"
            >
              <div class="flex items-center justify-between mb-2">
                <div>
                  <div class="font-semibold text-gray-800">{{ item.name }}</div>
                  <div class="text-xs text-gray-400">{{ item.category || '—' }}</div>
                </div>
                <span :class="diffClass(item)">{{ diffLabel(item) }}</span>
              </div>
              <div class="flex items-center gap-3">
                <div class="flex-1">
                  <div class="text-xs text-gray-500 mb-1">Hệ thống</div>
                  <div class="text-gray-700 font-medium">{{ fmt(item.current_qty) }} {{ item.unit }}</div>
                </div>
                <div class="flex-1">
                  <div class="text-xs text-gray-500 mb-1">Thực tế</div>
                  <div class="flex items-center gap-1">
                    <input
                      :value="physicalCounts[item.id] ?? ''"
                      @input="setPhysicalCount(item.id, $event.target.value)"
                      type="number"
                      min="0"
                      step="any"
                      :placeholder="fmt(item.current_qty)"
                      class="w-full border border-gray-300 rounded-lg px-2 py-1.5 text-sm focus:outline-none focus:border-blue-500"
                    />
                    <span class="text-gray-400 text-xs whitespace-nowrap">{{ item.unit }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Changed items summary + reason + submit -->
          <div class="bg-white rounded-2xl p-4 shadow-sm space-y-4">
            <div class="flex items-center justify-between">
              <span class="text-sm font-medium text-gray-700">Mục có thay đổi</span>
              <span :class="['text-sm font-bold', changedItems.length > 0 ? 'text-blue-600' : 'text-gray-400']">
                {{ changedItems.length }} / {{ items.length }}
              </span>
            </div>

            <div v-if="changedItems.length > 0" class="space-y-1">
              <div
                v-for="item in changedItems"
                :key="item.id"
                class="flex items-center justify-between text-sm py-1 border-b border-gray-50 last:border-0"
              >
                <span class="text-gray-700">{{ item.name }}</span>
                <span :class="diffClass(item)">{{ fmt(item.current_qty) }} → {{ fmt(physicalCounts[item.id]) }} {{ item.unit }}</span>
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Lý do kiểm kê</label>
              <input
                v-model="stocktakeReason"
                type="text"
                placeholder="Ví dụ: Kiểm kê cuối ngày 27/03..."
                class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-500"
              />
            </div>

            <button
              @click="submitStocktake"
              :disabled="changedItems.length === 0 || !stocktakeReason.trim() || submitting"
              :class="[
                'w-full py-3 rounded-xl font-semibold text-sm transition-all',
                changedItems.length > 0 && stocktakeReason.trim() && !submitting
                  ? 'bg-blue-600 text-white active:scale-95'
                  : 'bg-gray-200 text-gray-400 cursor-not-allowed'
              ]"
            >
              <span v-if="submitting">Đang xử lý...</span>
              <span v-else>Áp dụng điều chỉnh ({{ changedItems.length }} mục)</span>
            </button>

            <!-- Result message -->
            <div v-if="stocktakeResult" :class="['rounded-lg px-4 py-3 text-sm', stocktakeResult.success ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700']">
              {{ stocktakeResult.message }}
            </div>
          </div>
        </template>
      </template>

    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ingredientService } from '../services/ingredient'

const activeTab = ref('stats')
const loading = ref(false)
const items = ref([])

// ── Stats tab state ──
const search = ref('')
const filterStatus = ref('')
const filterCategory = ref('')
const sortKey = ref('status')
const sortDir = ref('asc')
const activePreset = ref('today')

// Date range
const today = new Date()
const fmt2 = (d) => d.toISOString().slice(0, 10)
const fromDate = ref(fmt2(today))
const toDate = ref(fmt2(today))

const datePresets = [
  { key: 'today', label: 'Hôm nay', getDates: () => ({ from: fmt2(today), to: fmt2(today) }) },
  {
    key: '7days',
    label: '7 ngày',
    getDates: () => {
      const d = new Date(today); d.setDate(d.getDate() - 6)
      return { from: fmt2(d), to: fmt2(today) }
    }
  },
  {
    key: '30days',
    label: '30 ngày',
    getDates: () => {
      const d = new Date(today); d.setDate(d.getDate() - 29)
      return { from: fmt2(d), to: fmt2(today) }
    }
  },
  {
    key: 'month',
    label: 'Tháng này',
    getDates: () => {
      const d = new Date(today.getFullYear(), today.getMonth(), 1)
      return { from: fmt2(d), to: fmt2(today) }
    }
  },
]

const columns = [
  { key: 'name', label: 'Nguyên liệu' },
  { key: 'status', label: 'Trạng thái' },
  { key: 'current_qty', label: 'Tồn kho' },
  { key: 'min_stock', label: 'Tối thiểu' },
  { key: 'consumed', label: 'Đã bán' },
  { key: 'restocked', label: 'Nhập kho' },
  { key: 'wasted', label: 'Hao hụt' },
  { key: 'adjusted', label: 'Điều chỉnh' },
]

const statusOrder = { out: 0, low: 1, ok: 2 }

const filteredItems = computed(() => {
  let list = [...items.value]

  if (search.value) {
    const q = search.value.toLowerCase()
    list = list.filter(i => i.name.toLowerCase().includes(q) || (i.category || '').toLowerCase().includes(q))
  }
  if (filterStatus.value) {
    list = list.filter(i => i.status === filterStatus.value)
  }
  if (filterCategory.value) {
    list = list.filter(i => i.category === filterCategory.value)
  }

  list.sort((a, b) => {
    let av = a[sortKey.value]
    let bv = b[sortKey.value]
    if (sortKey.value === 'status') {
      av = statusOrder[av] ?? 99
      bv = statusOrder[bv] ?? 99
    }
    if (typeof av === 'string') av = av.toLowerCase()
    if (typeof bv === 'string') bv = bv.toLowerCase()
    if (av < bv) return sortDir.value === 'asc' ? -1 : 1
    if (av > bv) return sortDir.value === 'asc' ? 1 : -1
    return 0
  })

  return list
})

const categories = computed(() => {
  const cats = new Set(items.value.map(i => i.category).filter(Boolean))
  return [...cats].sort()
})

const okCount = computed(() => items.value.filter(i => i.status === 'ok').length)
const lowCount = computed(() => items.value.filter(i => i.status === 'low').length)
const outCount = computed(() => items.value.filter(i => i.status === 'out').length)

// ── Stocktake tab state ──
const stocktakeSearch = ref('')
const stocktakeCategory = ref('')
const physicalCounts = ref({}) // id → number
const stocktakeReason = ref('')
const submitting = ref(false)
const stocktakeResult = ref(null)

const filteredStocktakeItems = computed(() => {
  let list = [...items.value].sort((a, b) => a.name.localeCompare(b.name))
  if (stocktakeSearch.value) {
    const q = stocktakeSearch.value.toLowerCase()
    list = list.filter(i => i.name.toLowerCase().includes(q) || (i.category || '').toLowerCase().includes(q))
  }
  if (stocktakeCategory.value) {
    list = list.filter(i => i.category === stocktakeCategory.value)
  }
  return list
})

const changedItems = computed(() => {
  return items.value.filter(item => {
    const val = physicalCounts.value[item.id]
    if (val === undefined || val === '') return false
    return parseFloat(val) !== parseFloat(item.current_qty)
  })
})

const setPhysicalCount = (id, value) => {
  if (value === '') {
    const counts = { ...physicalCounts.value }
    delete counts[id]
    physicalCounts.value = counts
  } else {
    physicalCounts.value = { ...physicalCounts.value, [id]: value }
  }
}

const getDiff = (item) => {
  const val = physicalCounts.value[item.id]
  if (val === undefined || val === '') return null
  return parseFloat(val) - parseFloat(item.current_qty)
}

const diffClass = (item) => {
  const diff = getDiff(item)
  if (diff === null) return 'text-gray-300 text-xs'
  if (Math.abs(diff) < 0.001) return 'inline-block px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-700'
  if (diff > 0) return 'inline-block px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-700'
  return 'inline-block px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-700'
}

const diffLabel = (item) => {
  const diff = getDiff(item)
  if (diff === null) return '—'
  if (Math.abs(diff) < 0.001) return '= Khớp'
  const sign = diff > 0 ? '+' : ''
  return `${sign}${fmt(diff)} ${item.unit}`
}

const submitStocktake = async () => {
  if (changedItems.value.length === 0 || !stocktakeReason.value.trim()) return
  submitting.value = true
  stocktakeResult.value = null

  const results = await Promise.allSettled(
    changedItems.value.map(item =>
      ingredientService.stockAdjust(item.id, {
        new_quantity: parseFloat(physicalCounts.value[item.id]),
        reason: stocktakeReason.value.trim()
      })
    )
  )

  const failed = results.filter(r => r.status === 'rejected').length
  const succeeded = results.length - failed

  submitting.value = false

  if (failed === 0) {
    stocktakeResult.value = { success: true, message: `Đã điều chỉnh ${succeeded} nguyên liệu thành công.` }
    // Clear inputs and reload
    physicalCounts.value = {}
    stocktakeReason.value = ''
    await loadData()
  } else {
    stocktakeResult.value = {
      success: false,
      message: `${succeeded} thành công, ${failed} thất bại. Vui lòng thử lại.`
    }
  }
}

// ── Shared helpers ──
const fmt = (n) => {
  if (!n && n !== 0) return '0'
  const num = parseFloat(n)
  return Number.isInteger(num) ? num.toString() : num.toFixed(2).replace(/\.?0+$/, '')
}

const statusBadgeClass = (status) => {
  if (status === 'out') return 'inline-block px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-700'
  if (status === 'low') return 'inline-block px-2 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-700'
  return 'inline-block px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-700'
}

const statusLabel = (status) => {
  if (status === 'out') return 'Hết hàng'
  if (status === 'low') return 'Sắp hết'
  return 'Đủ hàng'
}

const sortBy = (key) => {
  if (sortKey.value === key) {
    sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortKey.value = key
    sortDir.value = key === 'status' ? 'asc' : 'desc'
  }
}

const applyPreset = (preset) => {
  activePreset.value = preset.key
  const { from, to } = preset.getDates()
  fromDate.value = from
  toDate.value = to
  loadData()
}

const loadData = async () => {
  loading.value = true
  try {
    const data = await ingredientService.getStockSummary(fromDate.value, toDate.value)
    items.value = data.items || []
  } catch (err) {
    console.error('Failed to load stock summary:', err)
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>
