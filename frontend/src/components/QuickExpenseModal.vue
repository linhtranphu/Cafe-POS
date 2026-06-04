<template>
  <transition name="slide-up">
    <div v-if="show" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end" @click.self="$emit('close')">
      <div class="bg-white rounded-t-3xl w-full mb-20" style="max-height: 90vh; overflow-y: auto;">
        <!-- Header -->
        <div class="flex items-center justify-between px-5 pt-5 pb-3">
          <h3 class="text-xl font-bold">⚡ Chi phí nhanh</h3>
          <button @click="$emit('close')" class="text-gray-400 text-3xl leading-none">&times;</button>
        </div>

        <div class="px-5 pb-6 space-y-5">
          <!-- Mô tả -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Mô tả *</label>
            <input
              v-model="form.description"
              type="text"
              placeholder="VD: Mua đá, tiền điện..."
              class="w-full px-4 py-3 text-base border border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-400 focus:border-transparent"
            />
          </div>

          <!-- Số tiền -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Số tiền *</label>
            <div class="relative">
              <input
                v-model.number="form.amount"
                type="number"
                inputmode="numeric"
                min="0"
                step="1000"
                placeholder="0"
                class="w-full px-4 py-4 pr-10 text-2xl font-bold border border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-400 focus:border-transparent"
              />
              <span class="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 font-medium">₫</span>
            </div>
            <div class="flex gap-2 mt-2 flex-wrap">
              <button v-for="amt in quickAmounts" :key="amt"
                @click="form.amount = amt"
                class="px-3 py-1.5 text-sm bg-orange-50 text-orange-600 rounded-lg border border-orange-200 active:bg-orange-100">
                {{ formatShort(amt) }}
              </button>
            </div>
          </div>

          <!-- Hình thức thanh toán -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Hình thức *</label>
            <div class="grid grid-cols-2 gap-3">
              <button
                @click="form.money_type = 'cash'"
                :class="['py-4 rounded-xl border-2 font-semibold text-base transition-all',
                  form.money_type === 'cash'
                    ? 'border-green-500 bg-green-50 text-green-700'
                    : 'border-gray-200 bg-white text-gray-600 active:bg-gray-50']">
                💵 Tiền mặt
              </button>
              <button
                @click="form.money_type = 'transfer'"
                :class="['py-4 rounded-xl border-2 font-semibold text-base transition-all',
                  form.money_type === 'transfer'
                    ? 'border-blue-500 bg-blue-50 text-blue-700'
                    : 'border-gray-200 bg-white text-gray-600 active:bg-gray-50']">
                💳 Chuyển khoản
              </button>
            </div>
          </div>

          <!-- Quỹ -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Trừ quỹ *</label>
            <div class="grid grid-cols-2 gap-3">
              <button
                @click="form.fund_type = 'operating'"
                :class="['py-4 rounded-xl border-2 font-semibold text-base transition-all',
                  form.fund_type === 'operating'
                    ? 'border-purple-500 bg-purple-50 text-purple-700'
                    : 'border-gray-200 bg-white text-gray-600 active:bg-gray-50']">
                ⚙️ Vận hành
              </button>
              <button
                @click="form.fund_type = 'inventory'"
                :class="['py-4 rounded-xl border-2 font-semibold text-base transition-all',
                  form.fund_type === 'inventory'
                    ? 'border-emerald-500 bg-emerald-50 text-emerald-700'
                    : 'border-gray-200 bg-white text-gray-600 active:bg-gray-50']">
                📦 Hàng hóa
              </button>
            </div>
          </div>

          <!-- Danh mục -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Danh mục *</label>
            <select
              v-model="form.category_id"
              class="w-full px-4 py-3 text-base border border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-400">
              <option value="">Chọn danh mục</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
          </div>

          <!-- Ghi chú -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Ghi chú <span class="text-gray-400 font-normal">(tuỳ chọn)</span></label>
            <input
              v-model="form.notes"
              type="text"
              placeholder="Ghi chú thêm..."
              class="w-full px-4 py-3 text-base border border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-400"
            />
          </div>

          <!-- Error -->
          <p v-if="error" class="text-sm text-red-600 bg-red-50 px-4 py-3 rounded-xl">{{ error }}</p>

          <!-- Submit -->
          <button
            @click="submit"
            :disabled="!isValid || saving"
            :class="['w-full py-4 rounded-xl font-bold text-lg transition-all',
              isValid && !saving
                ? 'bg-orange-500 text-white active:bg-orange-600'
                : 'bg-gray-200 text-gray-400 cursor-not-allowed']">
            {{ saving ? 'Đang lưu...' : '✅ Lưu chi phí' }}
          </button>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { expenseService } from '../services/expense'
import { fundExpenseService } from '../services/fundExpenseService'

const props = defineProps({ show: Boolean })
const emit = defineEmits(['close', 'saved'])

const categories = ref([])
const saving = ref(false)
const error = ref('')
const quickAmounts = [10000, 20000, 50000, 100000, 200000, 500000]

const newForm = () => ({
  description: '',
  amount: null,
  money_type: 'cash',
  fund_type: 'operating',
  category_id: '',
  notes: '',
  date: new Date().toISOString().split('T')[0],
})

const form = ref(newForm())

const isValid = computed(() =>
  form.value.description.trim() &&
  form.value.amount > 0 &&
  form.value.category_id
)

const formatShort = (n) => {
  if (n >= 1000000) return (n / 1000000) + 'tr'
  if (n >= 1000) return (n / 1000) + 'k'
  return n
}

watch(() => props.show, async (val) => {
  if (val) {
    form.value = newForm()
    error.value = ''
    try {
      categories.value = await expenseService.getCategories()
    } catch {
      categories.value = []
    }
  }
})

const submit = async () => {
  if (!isValid.value) return
  saving.value = true
  error.value = ''
  try {
    await fundExpenseService.createExpenseFromFund({
      date: form.value.date + 'T00:00:00Z',
      category_id: form.value.category_id,
      amount: form.value.amount,
      description: form.value.description,
      notes: form.value.notes,
      money_type: form.value.money_type,
      fund_type: form.value.fund_type,
    })
    emit('saved')
    emit('close')
  } catch (e) {
    error.value = e.message || 'Có lỗi xảy ra'
  } finally {
    saving.value = false
  }
}
</script>
