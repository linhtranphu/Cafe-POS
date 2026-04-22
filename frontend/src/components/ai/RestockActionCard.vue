<template>
  <div class="border-2 border-blue-400 rounded-2xl overflow-hidden bg-white shadow-sm w-full max-w-sm">
    <div class="bg-blue-50 px-4 py-3 flex items-center gap-2 border-b border-blue-200">
      <span class="text-base">🔄</span>
      <span class="font-semibold text-sm text-blue-900">Nhập kho nguyên liệu</span>
    </div>
    <div class="px-4 py-3 flex flex-col gap-2">
      <div v-if="!editing">
        <FieldRow label="Tên" :value="fields.ingredient_name" />
        <FieldRow label="Tồn kho hiện tại" :value="`${fields.current_stock} ${fields.unit}`" />
        <FieldRow label="Số lượng nhập" :value="`${fields.quantity} ${fields.unit}`" />
        <FieldRow label="Đơn giá" :value="formatPrice(fields.cost_per_unit)" />
        <FieldRow label="Phương thức" :value="fields.money_type === 'cash' ? '💵 Tiền mặt' : '🏦 Chuyển khoản'" />
        <FieldRow v-if="fields.reason" label="Lý do" :value="fields.reason" />
      </div>
      <div v-else class="flex flex-col gap-2">
        <label class="text-xs text-gray-500">Số lượng nhập</label>
        <input v-model.number="edit.quantity" type="number" min="0" class="border rounded-lg px-3 py-2 text-sm w-full" />
        <label class="text-xs text-gray-500">Đơn giá (VNĐ)</label>
        <input v-model.number="edit.cost_per_unit" type="number" min="0" class="border rounded-lg px-3 py-2 text-sm w-full" />
        <label class="text-xs text-gray-500">Phương thức</label>
        <div class="flex gap-2">
          <button @click="edit.money_type = 'cash'"
            :class="edit.money_type === 'cash' ? 'bg-green-500 text-white' : 'bg-gray-100 text-gray-700'"
            class="flex-1 py-2 rounded-lg text-sm font-medium">💵 Tiền mặt</button>
          <button @click="edit.money_type = 'transfer'"
            :class="edit.money_type === 'transfer' ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-700'"
            class="flex-1 py-2 rounded-lg text-sm font-medium">🏦 CK</button>
        </div>
        <label class="text-xs text-gray-500">Lý do (tuỳ chọn)</label>
        <input v-model="edit.reason" class="border rounded-lg px-3 py-2 text-sm w-full" />
      </div>
      <p v-if="error" class="text-xs text-red-600 mt-1">{{ error }}</p>
    </div>
    <div v-if="!confirmed" class="flex border-t border-gray-100">
      <button @click="toggleEdit" class="flex-1 py-3 text-sm text-gray-600 font-medium border-r border-gray-100 active:bg-gray-50">
        {{ editing ? 'Xem lại' : '✏️ Sửa' }}
      </button>
      <button @click="confirm" :disabled="loading" class="flex-1 py-3 text-sm text-white font-semibold bg-blue-500 active:bg-blue-600 disabled:bg-gray-300">
        {{ loading ? 'Đang lưu...' : '✓ Xác nhận' }}
      </button>
    </div>
    <div v-else class="px-4 py-3 text-sm text-green-700 font-medium text-center bg-green-50">✅ Đã nhập kho thành công</div>
  </div>
</template>
<script setup>
import { ref, reactive } from 'vue'
import { fundIngredientService } from '../../services/fundIngredientService'
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
    await fundIngredientService.restockIngredientFromFund(data.ingredient_id, {
      quantity: data.quantity,
      cost_per_unit: data.cost_per_unit,
      money_type: data.money_type,
      reason: data.reason || '',
    })
    confirmed.value = true
    emit('confirmed', { type: 'restock_ingredient', name: data.ingredient_name })
  } catch (e) {
    error.value = e.message || 'Không thể nhập kho'
  } finally {
    loading.value = false
  }
}
</script>
