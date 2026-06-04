<template>
  <div class="border-2 border-purple-400 rounded-2xl overflow-hidden bg-white shadow-sm w-full max-w-sm">
    <div class="bg-purple-50 px-4 py-3 flex items-center gap-2 border-b border-purple-200">
      <span class="text-base">💸</span>
      <span class="font-semibold text-sm text-purple-900">Ghi nhận chi phí</span>
    </div>
    <div class="px-4 py-3 flex flex-col gap-2">
      <div v-if="!editing">
        <FieldRow label="Mô tả" :value="fields.description" />
        <FieldRow label="Số tiền" :value="formatPrice(fields.amount)" :highlight="true" />
        <FieldRow label="Phương thức" :value="fields.money_type === 'cash' ? '💵 Tiền mặt' : '🏦 Chuyển khoản'" />
        <FieldRow label="Danh mục" :value="categoryName" />
        <FieldRow label="Ngày" :value="fields.date || today" />
      </div>
      <div v-else class="flex flex-col gap-2">
        <label class="text-xs text-gray-500">Mô tả</label>
        <input v-model="edit.description" class="border rounded-lg px-3 py-2 text-sm w-full" />
        <label class="text-xs text-gray-500">Số tiền (VNĐ)</label>
        <input v-model.number="edit.amount" type="number" min="0" class="border rounded-lg px-3 py-2 text-sm w-full" />
        <label class="text-xs text-gray-500">Phương thức</label>
        <div class="flex gap-2">
          <button @click="edit.money_type = 'cash'"
            :class="edit.money_type === 'cash' ? 'bg-green-500 text-white' : 'bg-gray-100 text-gray-700'"
            class="flex-1 py-2 rounded-lg text-sm font-medium">💵 Tiền mặt</button>
          <button @click="edit.money_type = 'transfer'"
            :class="edit.money_type === 'transfer' ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-700'"
            class="flex-1 py-2 rounded-lg text-sm font-medium">🏦 CK</button>
        </div>
        <label class="text-xs text-gray-500">Danh mục</label>
        <select v-model="edit.category_id" class="border rounded-lg px-3 py-2 text-sm w-full">
          <option value="">-- Chọn danh mục --</option>
          <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
        </select>
        <label class="text-xs text-gray-500">Ngày</label>
        <input v-model="edit.date" type="date" class="border rounded-lg px-3 py-2 text-sm w-full" />
      </div>
      <p v-if="error" class="text-xs text-red-600 mt-1">{{ error }}</p>
    </div>
    <div v-if="!confirmed" class="flex border-t border-gray-100">
      <button @click="toggleEdit" class="flex-1 py-3 text-sm text-gray-600 font-medium border-r border-gray-100 active:bg-gray-50">
        {{ editing ? 'Xem lại' : '✏️ Sửa' }}
      </button>
      <button @click="confirm" :disabled="loading" class="flex-1 py-3 text-sm text-white font-semibold bg-purple-500 active:bg-purple-600 disabled:bg-gray-300">
        {{ loading ? 'Đang lưu...' : '✓ Xác nhận' }}
      </button>
    </div>
    <div v-else class="px-4 py-3 text-sm text-green-700 font-medium text-center bg-green-50">✅ Đã ghi nhận chi phí thành công</div>
  </div>
</template>
<script setup>
import { ref, reactive, computed } from 'vue'
import { fundExpenseService } from '../../services/fundExpenseService'
import FieldRow from './FieldRow.vue'

const props = defineProps({
  fields: { type: Object, required: true },
  categories: { type: Array, default: () => [] },
})
const emit = defineEmits(['confirmed'])

const editing = ref(false)
const loading = ref(false)
const confirmed = ref(false)
const error = ref('')
const edit = reactive({ ...props.fields })
const today = new Date().toISOString().slice(0, 10)

const categoryName = computed(() => {
  const cat = props.categories.find(c => c.id === props.fields.category_id)
  return cat ? cat.name : '—'
})
function toggleEdit() {
  if (!editing.value) Object.assign(edit, props.fields)
  editing.value = !editing.value
}
function formatPrice(v) {
  return new Intl.NumberFormat('vi-VN').format(v || 0) + 'đ'
}
async function confirm() {
  error.value = ''
  const data = editing.value ? edit : props.fields
  if (!data.category_id) {
    if (!editing.value) { editing.value = true; Object.assign(edit, props.fields) }
    error.value = 'Vui lòng chọn danh mục chi phí'
    return
  }
  loading.value = true
  try {
    await fundExpenseService.createExpenseFromFund({
      description: data.description,
      amount: data.amount,
      money_type: data.money_type,
      category_id: data.category_id,
      date: data.date || today,
    })
    confirmed.value = true
    emit('confirmed', { type: 'add_expense', description: data.description })
  } catch (e) {
    error.value = e.message || 'Không thể ghi nhận chi phí'
  } finally {
    loading.value = false
  }
}
</script>
