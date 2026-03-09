<template>
  <div class="bg-white rounded-2xl shadow-lg p-6 max-w-2xl mx-auto">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-bold text-gray-800 flex items-center gap-2">
        <span>💼</span>
        <span>Chi phí vận hành</span>
      </h2>
      <button 
        @click="handleCancel"
        class="text-gray-400 hover:text-gray-600 text-2xl leading-none"
        type="button">
        ×
      </button>
    </div>

    <!-- Form -->
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <!-- Period Date Pickers -->
      <div class="space-y-4">
        <h3 class="text-sm font-medium text-gray-700">Khoảng thời gian</h3>
        
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm text-gray-600 mb-1">Ngày bắt đầu</label>
            <input 
              v-model="formData.period_start"
              type="date"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              :class="{ 'border-red-500': errors.period_start }"
            />
            <p v-if="errors.period_start" class="text-xs text-red-500 mt-1">
              {{ errors.period_start }}
            </p>
          </div>
          
          <div>
            <label class="block text-sm text-gray-600 mb-1">Ngày kết thúc</label>
            <input 
              v-model="formData.period_end"
              type="date"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              :class="{ 'border-red-500': errors.period_end }"
            />
            <p v-if="errors.period_end" class="text-xs text-red-500 mt-1">
              {{ errors.period_end }}
            </p>
          </div>
        </div>
      </div>

      <!-- Expense Input Fields -->
      <div class="space-y-4">
        <h3 class="text-sm font-medium text-gray-700">Chi tiết chi phí</h3>
        
        <!-- Staff Salary -->
        <div>
          <label class="block text-sm text-gray-600 mb-1">
            <span class="inline-flex items-center gap-1">
              <span>👥</span>
              <span>Lương nhân viên</span>
            </span>
          </label>
          <input 
            v-model.number="formData.staff_salary"
            type="number"
            min="0"
            step="1000"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            :class="{ 'border-red-500': errors.staff_salary }"
            placeholder="0"
          />
          <p v-if="errors.staff_salary" class="text-xs text-red-500 mt-1">
            {{ errors.staff_salary }}
          </p>
        </div>

        <!-- Rent -->
        <div>
          <label class="block text-sm text-gray-600 mb-1">
            <span class="inline-flex items-center gap-1">
              <span>🏢</span>
              <span>Tiền thuê mặt bằng</span>
            </span>
          </label>
          <input 
            v-model.number="formData.rent"
            type="number"
            min="0"
            step="1000"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            :class="{ 'border-red-500': errors.rent }"
            placeholder="0"
          />
          <p v-if="errors.rent" class="text-xs text-red-500 mt-1">
            {{ errors.rent }}
          </p>
        </div>

        <!-- Utilities -->
        <div>
          <label class="block text-sm text-gray-600 mb-1">
            <span class="inline-flex items-center gap-1">
              <span>⚡</span>
              <span>Điện nước</span>
            </span>
          </label>
          <input 
            v-model.number="formData.utilities"
            type="number"
            min="0"
            step="1000"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            :class="{ 'border-red-500': errors.utilities }"
            placeholder="0"
          />
          <p v-if="errors.utilities" class="text-xs text-red-500 mt-1">
            {{ errors.utilities }}
          </p>
        </div>

        <!-- Marketing Costs -->
        <div>
          <label class="block text-sm text-gray-600 mb-1">
            <span class="inline-flex items-center gap-1">
              <span>📢</span>
              <span>Chi phí marketing</span>
            </span>
          </label>
          <input 
            v-model.number="formData.marketing_costs"
            type="number"
            min="0"
            step="1000"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            :class="{ 'border-red-500': errors.marketing_costs }"
            placeholder="0"
          />
          <p v-if="errors.marketing_costs" class="text-xs text-red-500 mt-1">
            {{ errors.marketing_costs }}
          </p>
        </div>

        <!-- Other Expenses -->
        <div>
          <label class="block text-sm text-gray-600 mb-1">
            <span class="inline-flex items-center gap-1">
              <span>📦</span>
              <span>Chi phí khác</span>
            </span>
          </label>
          <input 
            v-model.number="formData.other_expenses"
            type="number"
            min="0"
            step="1000"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            :class="{ 'border-red-500': errors.other_expenses }"
            placeholder="0"
          />
          <p v-if="errors.other_expenses" class="text-xs text-red-500 mt-1">
            {{ errors.other_expenses }}
          </p>
        </div>
      </div>

      <!-- Total Expenses Display -->
      <div class="bg-blue-50 rounded-xl p-4">
        <div class="flex justify-between items-center">
          <span class="text-sm font-medium text-gray-700">Tổng chi phí vận hành</span>
          <span class="text-xl font-bold text-blue-600">
            {{ formatPrice(totalExpenses) }}
          </span>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex gap-3 pt-4">
        <button 
          type="button"
          @click="handleCancel"
          class="flex-1 px-4 py-3 bg-gray-100 text-gray-700 rounded-lg font-medium hover:bg-gray-200 active:scale-95 transition-transform">
          Hủy
        </button>
        <button 
          type="submit"
          :disabled="saving"
          @click.prevent="handleSubmit"
          class="flex-1 px-4 py-3 bg-blue-500 text-white rounded-lg font-medium hover:bg-blue-600 active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed">
          {{ saving ? 'Đang lưu...' : 'Lưu' }}
        </button>
      </div>

      <!-- Error Message -->
      <div v-if="errors.submit" class="bg-red-50 border border-red-200 rounded-lg p-3">
        <p class="text-sm text-red-600">{{ errors.submit }}</p>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { formatPrice } from '../utils/formatters'
import { profitAnalysisService } from '../services/profitAnalysis'

// Props
const props = defineProps({
  initialData: {
    type: Object,
    default: null
  }
})

// Emits
const emit = defineEmits(['save', 'cancel'])

// State
const formData = ref({
  period_start: '',
  period_end: '',
  staff_salary: 0,
  rent: 0,
  utilities: 0,
  marketing_costs: 0,
  other_expenses: 0
})

const errors = ref({})
const saving = ref(false)

// Initialize form data if editing
if (props.initialData) {
  formData.value = {
    period_start: props.initialData.period_start || '',
    period_end: props.initialData.period_end || '',
    staff_salary: props.initialData.staff_salary || 0,
    rent: props.initialData.rent || 0,
    utilities: props.initialData.utilities || 0,
    marketing_costs: props.initialData.marketing_costs || 0,
    other_expenses: props.initialData.other_expenses || 0
  }
}

// Computed: Auto-calculate total expenses
const totalExpenses = computed(() => {
  return (
    (formData.value.staff_salary || 0) +
    (formData.value.rent || 0) +
    (formData.value.utilities || 0) +
    (formData.value.marketing_costs || 0) +
    (formData.value.other_expenses || 0)
  )
})

// Watch for changes to clear errors
watch(formData, () => {
  errors.value = {}
}, { deep: true })

// Validation function
const validateForm = () => {
  const newErrors = {}
  
  // Validate period dates
  if (!formData.value.period_start) {
    newErrors.period_start = 'Vui lòng chọn ngày bắt đầu'
  }
  
  if (!formData.value.period_end) {
    newErrors.period_end = 'Vui lòng chọn ngày kết thúc'
  }
  
  // Validate period_start <= period_end
  if (formData.value.period_start && formData.value.period_end) {
    const startDate = new Date(formData.value.period_start)
    const endDate = new Date(formData.value.period_end)
    
    if (startDate > endDate) {
      newErrors.period_end = 'Ngày kết thúc phải sau ngày bắt đầu'
    }
  }
  
  // Validate all amounts >= 0
  const amountFields = ['staff_salary', 'rent', 'utilities', 'marketing_costs', 'other_expenses']
  amountFields.forEach(field => {
    const value = formData.value[field]
    if (value < 0) {
      newErrors[field] = 'Số tiền không được âm'
    }
  })
  
  errors.value = newErrors
  return Object.keys(newErrors).length === 0
}

// Handle submit
const handleSubmit = async () => {
  console.log('[OperatingExpenseForm] Submit triggered')
  
  // Validate form
  if (!validateForm()) {
    console.log('[OperatingExpenseForm] Validation failed')
    return
  }
  
  if (saving.value) {
    console.log('[OperatingExpenseForm] Already saving, skipping')
    return
  }
  
  saving.value = true
  errors.value = {}
  
  try {
    // Prepare data for API
    const expenseData = {
      period_start: formData.value.period_start,
      period_end: formData.value.period_end,
      staff_salary: formData.value.staff_salary || 0,
      rent: formData.value.rent || 0,
      utilities: formData.value.utilities || 0,
      marketing_costs: formData.value.marketing_costs || 0,
      other_expenses: formData.value.other_expenses || 0
    }
    
    console.log('[OperatingExpenseForm] Calling API with data:', expenseData)
    
    // Call API
    const result = await profitAnalysisService.createOperatingExpense(expenseData)
    
    console.log('[OperatingExpenseForm] API success:', result)
    
    // Emit save event with result
    emit('save', result)
  } catch (err) {
    console.error('[OperatingExpenseForm] Error saving operating expense:', err)
    errors.value.submit = err.response?.data?.error || 'Không thể lưu chi phí vận hành'
  } finally {
    saving.value = false
  }
}

// Handle cancel
const handleCancel = () => {
  emit('cancel')
}
</script>

<style scoped>
.active\:scale-95:active {
  transform: scale(0.95);
}
</style>
