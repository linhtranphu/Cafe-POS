<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3" style="padding-top: max(0.75rem, env(safe-area-inset-top))">
        <div class="flex items-center justify-between">
          <button @click="goBack" class="text-2xl text-gray-600">←</button>
          <h1 class="text-xl font-bold text-gray-800">➕ Ghi Nhận Batch</h1>
          <div class="w-8"></div>
        </div>
      </div>
    </div>

    <!-- Scrollable Content -->
    <div class="flex-1 overflow-y-auto px-4 py-6 space-y-5">
      <!-- Batch Definition Selector -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Loại Batch *</label>
        <select 
          v-model="formData.batch_definition_id"
          @change="onBatchDefinitionChange"
          :disabled="isLoadingDefinitions"
          class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 disabled:opacity-50">
          <option value="">{{ isLoadingDefinitions ? 'Đang tải...' : 'Chọn loại batch' }}</option>
          <option 
            v-for="def in definitions" 
            :key="def?.id || Math.random()" 
            :value="def?.id">
            {{ def?.name || 'N/A' }} ({{ def?.unit || '' }})
          </option>
        </select>
      </div>

      <!-- Batch Count Selector -->
      <div v-if="selectedDefinition" class="bg-gradient-to-br from-purple-50 to-pink-50 border-2 border-purple-300 rounded-xl p-4">
        <div class="flex items-center gap-2 mb-3">
          <span class="text-purple-600 text-xl">🧪</span>
          <span class="text-sm font-semibold text-purple-800">Số Lượng Batch Sản Xuất</span>
        </div>
        
        <!-- Batch Counter -->
        <div class="flex items-center justify-center gap-4 mb-4">
          <button 
            @click="decrementBatchCount"
            :disabled="batchCount <= 1"
            class="w-12 h-12 bg-white border-2 border-purple-300 rounded-full text-2xl font-bold text-purple-600 active:bg-purple-100 disabled:opacity-30 disabled:cursor-not-allowed">
            −
          </button>
          
          <div class="text-center">
            <div class="text-4xl font-bold text-purple-700">{{ batchCount }}</div>
            <div class="text-xs text-purple-600 font-medium">batch</div>
          </div>
          
          <button 
            @click="incrementBatchCount"
            class="w-12 h-12 bg-white border-2 border-purple-300 rounded-full text-2xl font-bold text-purple-600 active:bg-purple-100">
            +
          </button>
        </div>

        <!-- Output Display -->
        <div class="bg-white rounded-lg p-3 border border-purple-200">
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-600">Tổng thành phẩm:</span>
            <div class="text-right">
              <div class="text-2xl font-bold text-purple-700">
                {{ totalOutput.toFixed(2) }}
              </div>
              <div class="text-xs text-purple-600 font-medium">
                {{ selectedDefinition.unit }}
              </div>
            </div>
          </div>
          
          <div class="mt-2 pt-2 border-t border-purple-100 text-xs text-gray-500">
            {{ batchCount }} batch × {{ batchOutputQuantity.toFixed(2) }} {{ selectedDefinition.unit }}/batch
          </div>
        </div>
      </div>

      <!-- Required Ingredients Display -->
      <div v-if="selectedDefinition && batchCount > 0" class="bg-blue-50 border-2 border-blue-300 rounded-xl p-4">
        <div class="flex items-center gap-2 mb-3">
          <span class="text-blue-600 text-xl">📦</span>
          <span class="text-sm font-semibold text-blue-800">Nguyên Liệu Cần Thiết</span>
        </div>
        
        <div class="space-y-2">
          <div 
            v-for="(ingredient, index) in requiredIngredients" 
            :key="index"
            class="flex justify-between items-center text-sm">
            <span class="text-gray-700">{{ ingredient.name }}</span>
            <span class="font-bold text-blue-700">
              {{ ingredient.quantity.toFixed(2) }} {{ ingredient.unit }}
            </span>
          </div>
        </div>

        <div class="mt-3 pt-3 border-t border-blue-200">
          <div class="flex justify-between items-center">
            <span class="text-sm font-semibold text-blue-800">💰 Chi Phí Dự Kiến:</span>
            <span class="text-lg font-bold text-blue-700">
              {{ formatCurrency(expectedCost) }}
            </span>
          </div>
          <div class="text-xs text-gray-600 mt-1">
            {{ formatCurrency(expectedCostPerUnit) }} / {{ selectedDefinition?.unit }}
          </div>
        </div>
      </div>

      <!-- Error Display -->
      <InlineError
        v-if="error"
        :message="error"
        :show="!!error"
        :showRetry="isRetryable"
        :onRetry="handleRetry"
        @dismiss="error = null"
      />


      <!-- Success Display -->
      <div v-if="successMessage" class="bg-green-50 border-2 border-green-300 rounded-xl p-4">
        <div class="flex items-center gap-2 mb-2">
          <span class="text-green-600 text-xl">✅</span>
          <span class="text-sm font-semibold text-green-800">{{ successMessage }}</span>
        </div>
        <button 
          @click="viewCreatedBatch"
          class="w-full mt-3 bg-green-500 text-white py-2 rounded-lg text-sm font-bold active:bg-green-600">
          Xem Chi Tiết Batch
        </button>
      </div>

      <!-- Spacer for bottom buttons -->
      <div class="h-24"></div>
    </div>

    <!-- Fixed Footer -->
    <div class="flex-shrink-0 bg-white px-4 py-4 border-t flex gap-3 pb-safe">
      <button 
        @click="goBack" 
        :disabled="isSubmitting"
        class="flex-1 bg-gray-200 text-gray-700 py-4 rounded-xl font-medium text-base active:bg-gray-300 disabled:opacity-50">
        Hủy
      </button>
      <button 
        @click="showConfirmation = true" 
        :disabled="isSubmitting || !isValid"
        class="flex-1 bg-blue-500 text-white py-4 rounded-xl font-medium text-base active:bg-blue-600 disabled:opacity-50">
        Ghi Nhận
      </button>
    </div>

    <!-- Confirmation Dialog -->
    <div v-if="showConfirmation" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl p-6 max-w-sm w-full">
        <h3 class="text-lg font-bold text-gray-800 mb-4">Xác Nhận Ghi Nhận Batch</h3>
        
        <div class="space-y-3 mb-6">
          <div class="flex justify-between text-sm">
            <span class="text-gray-600">Loại batch:</span>
            <span class="font-medium">{{ selectedDefinition?.name }}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-gray-600">Số lượng batch:</span>
            <span class="font-medium">{{ batchCount }} batch</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-gray-600">Tổng thành phẩm:</span>
            <span class="font-medium">{{ totalOutput.toFixed(2) }} {{ selectedDefinition?.unit }}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-gray-600">Chi phí:</span>
            <span class="font-bold text-blue-600">{{ formatCurrency(expectedCost) }}</span>
          </div>
        </div>

        <div class="bg-yellow-50 border border-yellow-300 rounded-lg p-3 mb-6">
          <p class="text-xs text-yellow-800">
            ⚠️ Nguyên liệu sẽ được trừ tự động khỏi kho. Hành động này không thể hoàn tác.
          </p>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <button 
            @click="showConfirmation = false"
            :disabled="isSubmitting"
            class="bg-gray-200 text-gray-700 py-3 rounded-xl font-medium active:bg-gray-300 disabled:opacity-50">
            Hủy
          </button>
          <button 
            @click="submitForm"
            :disabled="isSubmitting"
            class="bg-blue-500 text-white py-3 rounded-xl font-medium active:bg-blue-600 disabled:opacity-50 flex items-center justify-center gap-2">
            <span v-if="isSubmitting" class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
            <span>{{ isSubmitting ? 'Đang xử lý...' : 'Xác nhận' }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useBatchRecordStore } from '../../stores/batchRecord'
import { useBatchDefinitionStore } from '../../stores/batchDefinition'
import { useIngredientStore } from '../../stores/ingredient'
import { useAuthStore } from '../../stores/auth'
import { useBatchErrors } from '../../composables/useBatchErrors'
import { useUnitConversion } from '../../composables/useUnitConversion'
import InlineError from './InlineError.vue'

const router = useRouter()
const batchRecordStore = useBatchRecordStore()
const batchDefinitionStore = useBatchDefinitionStore()
const ingredientStore = useIngredientStore()
const authStore = useAuthStore()
const { isRetryableError } = useBatchErrors()
const { getConversionRate } = useUnitConversion()

const formData = ref({
  batch_definition_id: '',
  quantity_produced: 0,
  prepared_by: ''
})

const batchCount = ref(1)
const isSubmitting = ref(false)
const showConfirmation = ref(false)
const error = ref(null)
const successMessage = ref(null)
const isRetryable = ref(false)
const createdBatchId = ref(null)

const definitions = computed(() => batchDefinitionStore.definitions || [])
const ingredients = computed(() => ingredientStore.items || [])
const isLoadingDefinitions = computed(() => batchDefinitionStore.loading)

const selectedDefinition = computed(() => {
  if (!formData.value.batch_definition_id) return null
  const found = definitions.value.find(d => d?.id === formData.value.batch_definition_id)
  return found || null
})

// Get the batch output quantity from the first conversion rate
const batchOutputQuantity = computed(() => {
  if (!selectedDefinition.value || !selectedDefinition.value.conversion_rates || selectedDefinition.value.conversion_rates.length === 0) {
    return 0
  }
  // All conversion rates share the same batch_quantity
  return selectedDefinition.value.conversion_rates[0].batch_quantity || 0
})

// Calculate total output based on batch count
const totalOutput = computed(() => {
  return batchCount.value * batchOutputQuantity.value
})

const requiredIngredients = computed(() => {
  if (!selectedDefinition.value || batchCount.value <= 0) return []
  
  const def = selectedDefinition.value
  const quantity = totalOutput.value
  
  return (def.conversion_rates || []).map(rate => {
    // Calculate how much source ingredient is needed
    const ratio = quantity / rate.batch_quantity
    const baseQuantity = rate.source_quantity * ratio
    const wastageMultiplier = 1 + (rate.wastage_rate || 0)
    const totalQuantity = baseQuantity * wastageMultiplier
    
    return {
      name: rate.source_ingredient_name,
      quantity: totalQuantity,
      unit: rate.source_unit
    }
  })
})

const expectedCost = computed(() => {
  if (!selectedDefinition.value || batchCount.value <= 0) return 0
  
  let total = 0
  const def = selectedDefinition.value
  const quantity = totalOutput.value
  
  console.log('=== Batch Cost Calculation ===')
  console.log('Total Output:', quantity, def.output_unit)
  
  for (const rate of def.conversion_rates || []) {
    const ingredient = ingredients.value.find(i => i.id === rate.source_ingredient_id)
    if (ingredient && ingredient.cost_per_unit) {
      // Calculate quantity needed in source unit (recipe unit)
      const ratio = quantity / rate.batch_quantity
      const baseQuantity = rate.source_quantity * ratio
      const wastageMultiplier = 1 + (rate.wastage_rate || 0)
      const totalQuantity = baseQuantity * wastageMultiplier
      
      console.log('Ingredient:', ingredient.name)
      console.log('  Recipe needs:', rate.source_quantity, rate.source_unit)
      console.log('  Stock unit:', ingredient.unit, '@ ', ingredient.cost_per_unit, 'VNĐ')
      console.log('  Ratio:', ratio)
      console.log('  Base quantity:', baseQuantity, rate.source_unit)
      console.log('  With wastage:', totalQuantity, rate.source_unit)
      
      // Apply unit conversion: convert from source_unit (recipe) to ingredient.unit (stock)
      // We need to convert FROM recipe unit TO stock unit
      // getConversionRate(fromUnit, toUnit) where result means: 1 fromUnit = result toUnit
      // Example: getConversionRate("g", "kg") = 0.001 means 1g = 0.001kg
      // So: 200g * 0.001 = 0.2kg
      const conversionRate = getConversionRate(rate.source_unit, ingredient.unit)
      const quantityInStockUnit = totalQuantity * conversionRate
      
      console.log('  Conversion rate:', conversionRate, `(${rate.source_unit} → ${ingredient.unit})`)
      console.log('  Quantity in stock unit:', quantityInStockUnit, ingredient.unit)
      
      // Calculate cost using stock unit price
      const cost = quantityInStockUnit * ingredient.cost_per_unit
      console.log('  Cost:', cost, 'VNĐ')
      total += cost
    }
  }
  
  console.log('Total Cost:', total, 'VNĐ')
  console.log('===========================')
  
  return total
})

const expectedCostPerUnit = computed(() => {
  if (!totalOutput.value) return 0
  return expectedCost.value / totalOutput.value
})

const isValid = computed(() => {
  return formData.value.batch_definition_id && 
         batchCount.value > 0
})

const incrementBatchCount = () => {
  batchCount.value++
  updateQuantityProduced()
}

const decrementBatchCount = () => {
  if (batchCount.value > 1) {
    batchCount.value--
    updateQuantityProduced()
  }
}

const updateQuantityProduced = () => {
  formData.value.quantity_produced = totalOutput.value
}

const onBatchDefinitionChange = () => {
  error.value = null
  successMessage.value = null
  batchCount.value = 1
  updateQuantityProduced()
}

const formatCurrency = (value) => {
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND'
  }).format(value)
}

const submitForm = async () => {
  if (!isValid.value) return

  isSubmitting.value = true
  error.value = null
  successMessage.value = null

  try {
    // Get current user ID
    const userId = authStore.user?.id || 'unknown'
    
    // Ensure quantity_produced is up to date
    updateQuantityProduced()
    
    const payload = {
      batch_definition_id: formData.value.batch_definition_id,
      quantity_produced: formData.value.quantity_produced,
      prepared_by: userId
    }

    const createdRecord = await batchRecordStore.createRecord(payload)
    
    createdBatchId.value = createdRecord.id
    successMessage.value = `Đã ghi nhận ${batchCount.value} batch thành công! Batch sẽ hết hạn vào ${formatDateTime(createdRecord.expires_at)}`
    
    // Reset form
    formData.value = {
      batch_definition_id: '',
      quantity_produced: 0,
      prepared_by: ''
    }
    batchCount.value = 1
    
    showConfirmation.value = false
  } catch (err) {
    error.value = batchRecordStore.error || err.message || 'Lỗi ghi nhận batch'
    isRetryable.value = isRetryableError(err)
    showConfirmation.value = false
  } finally {
    isSubmitting.value = false
  }
}

const handleRetry = async () => {
  if (successMessage.value) {
    // If there's a success message, retry means submit again
    await submitForm()
  } else {
    // Otherwise, retry loading data
    error.value = null
    try {
      await Promise.all([
        batchDefinitionStore.fetchDefinitions(),
        ingredientStore.fetchIngredients()
      ])
      
      if (batchDefinitionStore.error) {
        error.value = batchDefinitionStore.error
        isRetryable.value = true
      }
    } catch (err) {
      error.value = 'Không thể tải dữ liệu. Vui lòng thử lại.'
      isRetryable.value = true
    }
  }
}

const formatDateTime = (dateStr) => {
  const date = new Date(dateStr)
  return new Intl.DateTimeFormat('vi-VN', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

const viewCreatedBatch = () => {
  if (createdBatchId.value) {
    router.push(`/batch/records/${createdBatchId.value}`)
  }
}

const goBack = () => {
  router.back()
}

onMounted(async () => {
  try {
    await Promise.all([
      batchDefinitionStore.fetchDefinitions(),
      ingredientStore.fetchIngredients()
    ])
    
    // Check if definitions loaded successfully
    if (batchDefinitionStore.error) {
      error.value = batchDefinitionStore.error
      isRetryable.value = true
    }
  } catch (err) {
    error.value = 'Không thể tải dữ liệu. Vui lòng thử lại.'
    isRetryable.value = true
  }
})
</script>
