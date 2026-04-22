<template>
  <div class="border-2 border-amber-400 rounded-2xl overflow-hidden bg-white shadow-sm w-full max-w-sm">
    <div class="bg-amber-50 px-4 py-3 flex items-center gap-2 border-b border-amber-200">
      <span class="text-base">📦</span>
      <span class="font-semibold text-sm text-amber-900">Thêm nguyên liệu mới</span>
    </div>
    <div class="px-4 py-3 flex flex-col gap-2">
      <div v-if="!editing">
        <FieldRow label="Tên" :value="fields.name" />
        <FieldRow label="Số lượng" :value="`${fields.quantity} ${fields.unit}`" />
        <FieldRow label="Đơn giá" :value="formatPrice(fields.cost_per_unit)" />
        <FieldRow label="Tổng giá trị" :value="formatPrice(fields.quantity * fields.cost_per_unit)" :highlight="true" />
      </div>
      <div v-else class="flex flex-col gap-2">
        <label class="text-xs text-gray-500">Tên</label>
        <input v-model="edit.name" class="border rounded-lg px-3 py-2 text-sm w-full" />
        <label class="text-xs text-gray-500">Số lượng</label>
        <div class="flex gap-2">
          <input v-model.number="edit.quantity" type="number" min="0" class="border rounded-lg px-3 py-2 text-sm flex-1" />
          <input v-model="edit.unit" class="border rounded-lg px-3 py-2 text-sm w-20" placeholder="đơn vị" />
        </div>
        <label class="text-xs text-gray-500">Đơn giá (VNĐ)</label>
        <input v-model.number="edit.cost_per_unit" type="number" min="0" class="border rounded-lg px-3 py-2 text-sm w-full" />
      </div>
      <p v-if="error" class="text-xs text-red-600 mt-1">{{ error }}</p>
    </div>
    <div v-if="!confirmed" class="flex border-t border-gray-100">
      <button @click="toggleEdit" class="flex-1 py-3 text-sm text-gray-600 font-medium border-r border-gray-100 active:bg-gray-50">
        {{ editing ? 'Xem lại' : '✏️ Sửa' }}
      </button>
      <button @click="confirm" :disabled="loading" class="flex-1 py-3 text-sm text-white font-semibold bg-amber-500 active:bg-amber-600 disabled:bg-gray-300">
        {{ loading ? 'Đang lưu...' : '✓ Xác nhận' }}
      </button>
    </div>
    <div v-else class="px-4 py-3 text-sm text-green-700 font-medium text-center bg-green-50">✅ Đã thêm thành công</div>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
import { ingredientService } from '../../services/ingredient'
import FieldRow from './FieldRow.vue'

const props = defineProps({ fields: { type: Object, required: true } })
const emit = defineEmits(['confirmed'])

const editing = ref(false)
const loading = ref(false)
const confirmed = ref(false)
const error = ref('')
const edit = reactive({ ...props.fields })

function toggleEdit() {
  if (!editing.value) Object.assign(edit, props.fields)
  editing.value = !editing.value
}
function formatPrice(v) {
  return new Intl.NumberFormat('vi-VN').format(v || 0) + 'đ'
}
async function confirm() {
  error.value = ''
  loading.value = true
  try {
    const data = editing.value ? edit : props.fields
    await ingredientService.createIngredient({
      name: data.name,
      quantity: data.quantity,
      unit: data.unit,
      cost_per_unit: data.cost_per_unit,
    })
    confirmed.value = true
    emit('confirmed', { type: 'add_ingredient', name: data.name })
  } catch (e) {
    error.value = e.message || 'Không thể thêm nguyên liệu'
  } finally {
    loading.value = false
  }
}
</script>
