<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white border-b border-gray-200 px-4 py-3 flex items-center gap-3 sticky top-0 z-10">
      <button @click="$router.back()" class="p-2 rounded-lg hover:bg-gray-100 active:bg-gray-200 transition-colors">
        <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      <h1 class="text-lg font-bold text-gray-800">Nhập nhanh nguyên liệu</h1>
    </div>

    <div class="p-4 space-y-3">
      <!-- Search + toggle -->
      <div class="flex gap-2">
        <input
          v-model="search"
          type="text"
          placeholder="Tìm nguyên liệu..."
          class="flex-1 px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-lime-500 focus:border-transparent text-base"
        />
        <button
          @click="toggleShowAll"
          :class="[
            'px-3 py-2 rounded-xl text-sm font-medium border transition-colors whitespace-nowrap',
            showAll
              ? 'bg-lime-500 text-white border-lime-500'
              : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
          ]"
        >
          Hiện tất cả
        </button>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="text-center py-8 text-gray-400">Đang tải...</div>

      <!-- Error -->
      <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-xl p-4 text-red-700 text-sm">
        {{ error }}
      </div>

      <!-- Empty -->
      <div v-else-if="displayList.length === 0" class="text-center py-8 text-gray-400 text-sm">
        {{ search ? 'Không tìm thấy nguyên liệu phù hợp' : 'Chưa có nguyên liệu nào được nhập trong 30 ngày qua' }}
      </div>

      <!-- Ingredient list -->
      <div v-else class="space-y-2">
        <div
          v-for="ing in displayList"
          :key="ing.id"
          class="bg-white rounded-xl border border-gray-200 overflow-hidden shadow-sm"
        >
          <!-- Collapsed row -->
          <button
            @click="selectIngredient(ing)"
            class="w-full px-4 py-3 flex items-center justify-between text-left active:bg-gray-50 transition-colors"
          >
            <div>
              <div class="font-semibold text-gray-800">{{ ing.name }}</div>
              <div class="text-xs text-gray-500 mt-0.5">
                Tồn: {{ ing.quantity }} {{ ing.unit }}
                <span v-if="ing.last_restock" class="ml-2 text-lime-600">
                  • Lần cuối: {{ ing.last_restock.quantity }} {{ ing.unit }}
                </span>
              </div>
            </div>
            <svg
              :class="['w-5 h-5 text-gray-400 transition-transform', activeId === ing.id ? 'rotate-180' : '']"
              fill="none" stroke="currentColor" viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </button>

          <!-- Expanded form -->
          <div v-if="activeId === ing.id" class="border-t border-gray-100 px-4 pb-4 pt-3 space-y-4">
            <!-- Quantity counter -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Số lượng ({{ ing.unit }})</label>
              <div class="flex items-center gap-2">
                <button
                  @click="decrement"
                  class="w-12 h-12 rounded-xl bg-gray-100 text-gray-700 text-xl font-bold flex items-center justify-center active:bg-gray-200 transition-colors"
                >−</button>
                <input
                  v-model.number="quantity"
                  type="number"
                  min="0"
                  step="any"
                  class="flex-1 h-12 text-center text-lg font-bold border-2 border-gray-300 rounded-xl focus:ring-2 focus:ring-lime-500 focus:border-lime-500"
                />
                <button
                  @click="increment"
                  class="w-12 h-12 rounded-xl bg-lime-500 text-white text-xl font-bold flex items-center justify-center active:bg-lime-600 transition-colors"
                >+</button>
              </div>
            </div>

            <!-- Cost per unit -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Đơn giá (₫ / {{ ing.unit }})</label>
              <input
                :value="costPerUnit"
                @input="onCostPerUnitInput"
                type="number"
                min="0"
                class="w-full h-12 px-4 text-base border-2 border-gray-300 rounded-xl focus:ring-2 focus:ring-lime-500 focus:border-lime-500"
              />
            </div>

            <!-- Total cost (editable — back-calculates cost per unit) -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Tổng tiền (₫)</label>
              <input
                :value="totalCostInput"
                @input="onTotalCostInput"
                type="number"
                min="0"
                class="w-full h-12 px-4 text-base border-2 border-lime-300 rounded-xl focus:ring-2 focus:ring-lime-500 focus:border-lime-500 bg-lime-50"
              />
            </div>

            <!-- Payment method -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Hình thức thanh toán</label>
              <div class="flex gap-2">
                <button
                  v-for="opt in paymentOptions"
                  :key="opt.value"
                  @click="moneyType = opt.value"
                  :class="[
                    'flex-1 py-3 rounded-xl font-medium text-sm border-2 transition-colors',
                    moneyType === opt.value
                      ? 'bg-lime-500 border-lime-500 text-white'
                      : 'bg-white border-gray-300 text-gray-700'
                  ]"
                >
                  {{ opt.label }}
                </button>
              </div>
            </div>

            <!-- Submit -->
            <button
              @click="submit(ing)"
              :disabled="submitting || quantity <= 0"
              class="w-full py-4 bg-lime-500 text-white rounded-xl font-bold text-base active:bg-lime-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {{ submitting ? 'Đang lưu...' : 'Lưu' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Toast -->
    <div
      v-if="toast.show"
      :class="[
        'fixed bottom-6 left-4 right-4 p-4 rounded-xl text-white font-medium text-center shadow-lg transition-all z-50',
        toast.success ? 'bg-green-500' : 'bg-red-500'
      ]"
    >
      {{ toast.message }}
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ingredientService } from '../services/ingredient'
import { fundIngredientService } from '../services/fundIngredientService'

const ingredients = ref([])
const allIngredients = ref([])
const allFetched = ref(false)
const showAll = ref(false)
const search = ref('')
const loading = ref(false)
const error = ref(null)

const activeId = ref(null)
const step = ref(1)
const quantity = ref(0)
const costPerUnit = ref(0)
const totalCostInput = ref(0)
const moneyType = ref('cash')
const submitting = ref(false)

const toast = ref({ show: false, success: true, message: '' })

const paymentOptions = [
  { value: 'cash', label: 'Tiền mặt' },
  { value: 'transfer', label: 'Chuyển khoản' }
]

const displayList = computed(() => {
  const source = showAll.value ? allIngredients.value : ingredients.value
  const q = search.value.toLowerCase().trim()
  return q ? source.filter(i => i.name.toLowerCase().includes(q)) : source
})

async function fetchRecent() {
  loading.value = true
  error.value = null
  try {
    ingredients.value = await ingredientService.getRecentRestocked(30)
  } catch {
    error.value = 'Không thể tải danh sách. Vui lòng thử lại.'
  } finally {
    loading.value = false
  }
}

async function fetchAll() {
  if (allFetched.value) return
  try {
    const data = await ingredientService.getIngredients()
    allIngredients.value = [...data].sort((a, b) => {
      const ta = a.created_at ? new Date(a.created_at).getTime() : 0
      const tb = b.created_at ? new Date(b.created_at).getTime() : 0
      return tb - ta
    })
    allFetched.value = true
  } catch {
    // Non-fatal: keep allIngredients empty
  }
}

async function toggleShowAll() {
  showAll.value = !showAll.value
  if (showAll.value) await fetchAll()
}

onMounted(fetchRecent)

function selectIngredient(ing) {
  if (activeId.value === ing.id) {
    activeId.value = null
    return
  }
  activeId.value = ing.id
  step.value = ing.last_restock?.quantity ?? 1
  quantity.value = 0
  costPerUnit.value = ing.last_restock?.cost_per_unit ?? ing.cost_per_unit ?? 0
  totalCostInput.value = 0
  moneyType.value = 'cash'
}

function syncTotal() {
  totalCostInput.value = +(quantity.value * costPerUnit.value).toFixed(0)
}

function onCostPerUnitInput(e) {
  costPerUnit.value = Math.max(0, Number(e.target.value) || 0)
  syncTotal()
}

function onTotalCostInput(e) {
  const total = Math.max(0, Number(e.target.value) || 0)
  totalCostInput.value = total
  if (quantity.value > 0) {
    costPerUnit.value = Math.round(total / quantity.value)
  }
}

function increment() {
  quantity.value = Math.max(0, +(quantity.value + step.value).toFixed(6))
  syncTotal()
}

function decrement() {
  quantity.value = Math.max(0, +(quantity.value - step.value).toFixed(6))
  syncTotal()
}

async function submit(ing) {
  if (quantity.value <= 0) {
    showToast('Vui lòng nhập số lượng hợp lệ', false)
    return
  }
  submitting.value = true
  try {
    await fundIngredientService.restockIngredientFromFund(ing.id, {
      quantity: quantity.value,
      cost_per_unit: costPerUnit.value,
      total_cost: totalCostInput.value || quantity.value * costPerUnit.value,
      reason: 'Nhập nhanh',
      money_type: moneyType.value
    })
    showToast(`Đã nhập ${quantity.value} ${ing.unit} ${ing.name}`, true)
    activeId.value = null
    await fetchRecent()
  } catch (e) {
    showToast(e?.response?.data?.error || e.message || 'Có lỗi xảy ra', false)
  } finally {
    submitting.value = false
  }
}

function formatCurrency(value) {
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(value || 0)
}

function showToast(message, success) {
  toast.value = { show: true, success, message }
  setTimeout(() => { toast.value.show = false }, 2500)
}
</script>
