<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
    <div class="bg-gray-50 w-full h-screen flex flex-col">
      <!-- Mobile Header - Fixed -->
      <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
        <div class="px-4 py-3">
          <div class="flex items-center justify-between">
            <button @click="$emit('close')" class="text-2xl text-gray-600">←</button>
            <h1 class="text-xl font-bold text-gray-800">
              {{ mode === 'edit' ? '✏️ Cập nhật Batch' : '➕ Tạo Batch Mới' }}
            </h1>
            <div class="w-8"></div>
          </div>
        </div>
      </div>

      <!-- Scrollable Content -->
      <div class="flex-1 overflow-y-auto px-4 py-6 space-y-5">
        
        <!-- SECTION 1: THÔNG TIN THÀNH PHẨM BATCH -->
        <div class="bg-gradient-to-br from-purple-50 to-pink-50 rounded-2xl p-4 border-2 border-purple-200">
          <div class="flex items-center gap-2 mb-4">
            <div class="text-2xl">🧪</div>
            <h2 class="text-lg font-bold text-purple-900">Thông Tin Thành Phẩm Batch</h2>
          </div>

          <div class="space-y-4">
            <!-- Tên batch -->
            <div>
              <label class="block text-sm font-bold text-gray-700 mb-2">📝 Tên Batch *</label>
              <input 
                v-model="formData.name" 
                type="text" 
                placeholder="VD: Cà Phê Concentrate"
                class="w-full px-4 py-3 border-2 border-purple-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 bg-white" 
              />
            </div>

            <!-- Khối lượng & Đơn vị thành phẩm -->
            <div class="bg-white rounded-xl p-4 border-2 border-purple-300">
              <div class="text-sm font-bold text-purple-800 mb-3">📦 Khối Lượng Thành Phẩm</div>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="block text-xs font-medium text-gray-600 mb-2">Số lượng *</label>
                  <input 
                    v-model.number="batchOutputQuantity" 
                    type="number" 
                    min="0"
                    step="0.01"
                    placeholder="VD: 500"
                    class="w-full px-4 py-3 text-lg font-bold border-2 border-purple-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 bg-white" 
                  />
                </div>
                <div>
                  <label class="block text-xs font-medium text-gray-600 mb-2">Đơn vị *</label>
                  <select 
                    v-model="formData.unit" 
                    class="w-full px-4 py-3 text-lg font-bold border-2 border-purple-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 bg-white">
                    <option value="">Chọn đơn vị</option>
                    <optgroup label="Khối lượng">
                      <option value="g">g (gram)</option>
                      <option value="kg">kg (kilogram)</option>
                    </optgroup>
                    <optgroup label="Thể tích">
                      <option value="ml">ml (milliliter)</option>
                      <option value="L">L (liter)</option>
                    </optgroup>
                    <optgroup label="Khác">
                      <option value="cái">cái</option>
                      <option value="phần">phần</option>
                    </optgroup>
                  </select>
                </div>
              </div>
              <p class="text-xs text-gray-500 mt-2">Khối lượng batch sau khi chế biến từ tất cả nguyên liệu</p>
            </div>

            <!-- Thời hạn sử dụng -->
            <div>
              <label class="block text-sm font-bold text-gray-700 mb-2">⏰ Thời Hạn Sử Dụng (giờ) *</label>
              <input 
                v-model.number="formData.shelf_life_hours" 
                type="number" 
                min="1"
                placeholder="VD: 24"
                class="w-full px-4 py-3 border-2 border-purple-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 bg-white" 
              />
              <p class="text-xs text-gray-500 mt-1">Batch có thể sử dụng trong bao lâu sau khi chế biến</p>
            </div>

            <!-- Ngưỡng cảnh báo -->
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-bold text-gray-700 mb-2">⚠️ Tồn kho thấp *</label>
                <input 
                  v-model.number="formData.low_stock_threshold" 
                  type="number" 
                  min="0"
                  step="0.01"
                  placeholder="VD: 200"
                  class="w-full px-3 py-2 text-sm border-2 border-purple-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 bg-white" 
                />
              </div>
              <div>
                <label class="block text-xs font-bold text-gray-700 mb-2">🔔 Cảnh báo hết hạn (giờ) *</label>
                <input 
                  v-model.number="formData.expiry_warning_hours" 
                  type="number" 
                  min="1"
                  placeholder="VD: 4"
                  class="w-full px-3 py-2 text-sm border-2 border-purple-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 bg-white" 
                />
              </div>
            </div>
          </div>
        </div>

        <!-- SECTION 2: NGUYÊN LIỆU NGUỒN (NHIỀU NGUYÊN LIỆU) -->
        <div class="bg-gradient-to-br from-green-50 to-emerald-50 rounded-2xl p-4 border-2 border-green-200">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <div class="text-2xl">🥬</div>
              <div>
                <h2 class="text-lg font-bold text-green-900">Nguyên Liệu Nguồn</h2>
                <p class="text-xs text-green-700">Các nguyên liệu cần để tạo batch</p>
              </div>
            </div>
            <button 
              @click="addConversionRate"
              type="button"
              class="bg-green-600 text-white px-4 py-2 rounded-lg text-sm font-bold active:bg-green-700 shadow-md flex-shrink-0">
              + Thêm
            </button>
          </div>

          <!-- Loading ingredients -->
          <div v-if="loadingIngredients" class="text-center py-8 bg-white rounded-xl">
            <div class="inline-block w-6 h-6 border-2 border-green-500 border-t-transparent rounded-full animate-spin"></div>
            <p class="text-xs text-gray-500 mt-2">Đang tải nguyên liệu...</p>
          </div>

          <!-- No ingredients available -->
          <div v-else-if="ingredients.length === 0" class="text-center py-8 bg-white rounded-xl">
            <div class="text-4xl mb-2">⚠️</div>
            <p class="text-sm font-medium text-red-600">Không có nguyên liệu nào</p>
            <p class="text-xs text-gray-500 mt-1">Vui lòng thêm nguyên liệu trước</p>
          </div>

          <!-- No conversion rates added yet -->
          <div v-else-if="formData.conversion_rates.length === 0" class="text-center py-8 bg-white rounded-xl">
            <div class="text-4xl mb-2">🍽️</div>
            <p class="text-sm text-gray-600">Chưa có nguyên liệu nào</p>
            <p class="text-xs text-gray-500 mt-1">Nhấn "Thêm" để thêm nguyên liệu nguồn</p>
          </div>

          <!-- Conversion rates list -->
          <div v-else class="space-y-3">
            <div 
              v-for="(rate, index) in formData.conversion_rates" 
              :key="index"
              class="bg-white rounded-xl p-4 shadow-sm border-2 border-green-100">
              
              <!-- Header with ingredient number -->
              <div class="flex items-center justify-between mb-3 pb-2 border-b-2 border-green-100">
                <div class="flex items-center gap-2">
                  <div class="w-8 h-8 bg-green-500 text-white rounded-full flex items-center justify-center font-bold text-sm">
                    {{ index + 1 }}
                  </div>
                  <span class="font-bold text-gray-800">Nguyên liệu #{{ index + 1 }}</span>
                </div>
                <button 
                  @click="removeConversionRate(index)"
                  type="button"
                  class="bg-red-500 text-white px-3 py-1 rounded-lg text-xs font-bold active:bg-red-600">
                  🗑️ Xóa
                </button>
              </div>

              <div class="space-y-3">
                <!-- Ingredient Selector -->
                <div>
                  <label class="block text-xs font-bold text-gray-700 mb-2">🥘 Chọn Nguyên Liệu *</label>
                  <select 
                    v-model="rate.source_ingredient_id"
                    @change="updateIngredientName(index)"
                    class="w-full px-3 py-3 text-sm border-2 border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-green-500 bg-white font-medium">
                    <option value="">-- Chọn nguyên liệu --</option>
                    <option 
                      v-for="ingredient in ingredients" 
                      :key="ingredient.id" 
                      :value="ingredient.id">
                      {{ ingredient.name }} ({{ ingredient.unit }})
                    </option>
                  </select>
                </div>

                <!-- Selected ingredient info -->
                <div v-if="getSelectedIngredient(rate.source_ingredient_id)" class="bg-blue-50 border-2 border-blue-200 rounded-lg p-3">
                  <div class="text-xs font-bold text-blue-800 mb-2">📊 Thông tin nguyên liệu:</div>
                  <div class="space-y-1 text-xs">
                    <div class="flex justify-between">
                      <span class="text-gray-600">Tồn kho:</span>
                      <span class="font-bold text-gray-800">{{ getSelectedIngredient(rate.source_ingredient_id).quantity || 0 }} {{ getSelectedIngredient(rate.source_ingredient_id).unit }}</span>
                    </div>
                    <div class="flex justify-between">
                      <span class="text-gray-600">Giá:</span>
                      <span class="font-bold text-gray-800">{{ formatCurrency(getSelectedIngredient(rate.source_ingredient_id).cost_per_unit || 0) }}/{{ getSelectedIngredient(rate.source_ingredient_id).unit }}</span>
                    </div>
                  </div>
                </div>

                <!-- Source Quantity Section -->
                <div class="bg-gray-50 rounded-lg p-3 border border-gray-200">
                  <div class="text-xs font-bold text-gray-700 mb-3">📦 Số Lượng Cần Dùng</div>
                  
                  <div class="grid grid-cols-2 gap-3">
                    <div>
                      <label class="block text-xs font-medium text-gray-600 mb-1">Số lượng *</label>
                      <input 
                        v-model.number="rate.source_quantity" 
                        type="number" 
                        min="0"
                        step="0.01"
                        placeholder="VD: 100"
                        class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 bg-white font-medium" 
                      />
                    </div>
                    <div>
                      <label class="block text-xs font-medium text-gray-600 mb-1">Đơn vị *</label>
                      <select 
                        v-model="rate.source_unit" 
                        @change="updateConversionRate(index)"
                        :disabled="!rate.source_ingredient_id"
                        class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 font-medium"
                        :class="!rate.source_ingredient_id ? 'bg-gray-100 text-gray-400' : 'bg-white'">
                        <option value="">Chọn đơn vị</option>
                        <option 
                          v-for="unit in rate.compatibleUnits || []" 
                          :key="unit" 
                          :value="unit">
                          {{ unit }}
                        </option>
                      </select>
                    </div>
                  </div>

                  <!-- Conversion Info -->
                  <div v-if="rate.source_ingredient_id && rate.source_unit && rate.conversionRate !== 1" 
                    class="mt-2 p-2 bg-blue-50 rounded-lg text-xs text-blue-700 border border-blue-200">
                    <span class="font-bold">ℹ️ Quy đổi:</span> {{ getConversionExplanation(getSelectedIngredient(rate.source_ingredient_id)?.unit, rate.source_unit) }}
                  </div>
                </div>

                <!-- Wastage Rate -->
                <div class="bg-orange-50 rounded-lg p-3 border border-orange-200">
                  <label class="block text-xs font-bold text-orange-800 mb-2">⚠️ Tỷ Lệ Hao Hụt (%)</label>
                  <input 
                    v-model.number="rate.wastage_rate_percent" 
                    type="number" 
                    min="0"
                    max="100"
                    step="0.1"
                    placeholder="VD: 10"
                    class="w-full px-3 py-2 text-sm border-2 border-orange-200 rounded-lg focus:ring-2 focus:ring-orange-500 bg-white font-medium" 
                  />
                  <p class="text-xs text-gray-600 mt-1">Hao hụt khi chế biến nguyên liệu này</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- SECTION 3: CÔNG THỨC CHUYỂN ĐỔI -->
        <div v-if="formData.conversion_rates.length > 0 && batchOutputQuantity > 0" 
          class="bg-gradient-to-br from-blue-50 to-indigo-50 rounded-2xl p-4 border-2 border-blue-300">
          <div class="flex items-center gap-2 mb-4">
            <span class="text-2xl">📐</span>
            <span class="text-lg font-bold text-blue-900">Công Thức Chuyển Đổi</span>
          </div>

          <div class="bg-white rounded-xl p-4">
            <!-- Input ingredients -->
            <div class="mb-4">
              <div class="text-xs font-bold text-gray-600 mb-3">NGUYÊN LIỆU ĐẦU VÀO:</div>
              <div class="space-y-2">
                <div v-for="(rate, index) in formData.conversion_rates" :key="index" 
                  class="flex items-center gap-2 bg-green-50 px-3 py-2 rounded-lg border border-green-200">
                  <span class="w-6 h-6 bg-green-500 text-white rounded-full flex items-center justify-center text-xs font-bold flex-shrink-0">
                    {{ index + 1 }}
                  </span>
                  <span class="text-sm font-medium text-gray-800 flex-1">
                    {{ getSelectedIngredient(rate.source_ingredient_id)?.name || 'Chưa chọn' }}
                  </span>
                  <span class="text-sm font-bold text-green-700 whitespace-nowrap">
                    {{ rate.source_quantity || 0 }} {{ rate.source_unit || '?' }}
                  </span>
                  <span v-if="rate.wastage_rate_percent > 0" class="text-xs text-orange-600 whitespace-nowrap">
                    (+{{ rate.wastage_rate_percent }}%)
                  </span>
                </div>
              </div>
            </div>

            <!-- Arrow -->
            <div class="text-center my-3">
              <div class="inline-flex items-center gap-2 text-gray-400">
                <div class="h-px w-16 bg-gray-300"></div>
                <span class="text-2xl">⬇️</span>
                <div class="h-px w-16 bg-gray-300"></div>
              </div>
              <div class="text-xs text-gray-500 mt-1">Chế biến thành</div>
            </div>

            <!-- Output batch -->
            <div>
              <div class="text-xs font-bold text-gray-600 mb-3">THÀNH PHẨM BATCH:</div>
              <div class="bg-purple-50 px-4 py-3 rounded-lg border-2 border-purple-300">
                <div class="flex items-center justify-between">
                  <span class="text-sm font-medium text-gray-800">{{ formData.name || 'Batch' }}</span>
                  <span class="text-lg font-bold text-purple-700">
                    {{ batchOutputQuantity }} {{ formData.unit || '?' }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- SECTION 4: TÓM TẮT CHI PHÍ -->
        <div v-if="estimatedCost > 0 && batchOutputQuantity > 0" class="bg-gradient-to-br from-yellow-50 to-orange-50 rounded-2xl p-4 border-2 border-yellow-300 shadow-lg">
          <div class="flex items-center gap-2 mb-3">
            <span class="text-2xl">💰</span>
            <span class="text-lg font-bold text-yellow-900">Tóm Tắt Chi Phí</span>
          </div>
          
          <div class="bg-white rounded-xl p-4 mb-3">
            <div class="flex justify-between items-center mb-2">
              <span class="text-sm text-gray-600">Chi phí nguyên liệu:</span>
              <span class="text-lg font-bold text-gray-800">{{ formatCurrency(estimatedCost) }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Khối lượng batch:</span>
              <span class="text-lg font-bold text-purple-600">
                {{ batchOutputQuantity }} {{ formData.unit }}
              </span>
            </div>
          </div>

          <div class="bg-gradient-to-r from-green-500 to-emerald-500 rounded-xl p-4 text-white">
            <div class="text-xs font-medium mb-1">Chi phí trên 1 {{ formData.unit }}:</div>
            <div class="text-2xl font-bold">
              {{ formatCurrency(estimatedCost / batchOutputQuantity) }}
            </div>
          </div>
        </div>

        <!-- Spacer for bottom buttons -->
        <div class="h-24"></div>
      </div>

      <!-- Fixed Footer -->
      <div class="flex-shrink-0 bg-white px-4 py-4 border-t flex gap-3 pb-safe">
        <button 
          @click="$emit('close')" 
          :disabled="isSubmitting"
          class="flex-1 bg-gray-200 text-gray-700 py-4 rounded-xl font-medium text-base active:bg-gray-300 disabled:opacity-50">
          Hủy
        </button>
        <button 
          @click="save" 
          :disabled="isSubmitting || !isValid"
          class="flex-1 bg-blue-500 text-white py-4 rounded-xl font-medium text-base active:bg-blue-600 disabled:opacity-50 flex items-center justify-center gap-2">
          <span v-if="isSubmitting" class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
          <span>{{ isSubmitting ? 'Đang lưu...' : (mode === 'edit' ? 'Cập nhật' : 'Tạo mới') }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useBatchDefinitionStore } from '../../stores/batchDefinition'
import { useIngredientStore } from '../../stores/ingredient'
import { useUnitConversion } from '../../composables/useUnitConversion'

const props = defineProps({
  definition: {
    type: Object,
    default: null
  },
  mode: {
    type: String,
    default: 'create',
    validator: (value) => ['create', 'edit'].includes(value)
  }
})

const emit = defineEmits(['close', 'saved'])

const batchStore = useBatchDefinitionStore()
const ingredientStore = useIngredientStore()
const { getCompatibleUnits, getConversionRate, getConversionExplanation } = useUnitConversion()

const isSubmitting = ref(false)
const batchOutputQuantity = ref(0)

const formData = ref({
  name: '',
  unit: '',
  shelf_life_hours: 24,
  low_stock_threshold: 0,
  expiry_warning_hours: 4,
  conversion_rates: []
})

const ingredients = computed(() => ingredientStore.items || [])
const loadingIngredients = computed(() => ingredientStore.loading)

const isValid = computed(() => {
  return formData.value.name && 
         formData.value.unit && 
         formData.value.shelf_life_hours > 0 &&
         formData.value.low_stock_threshold >= 0 &&
         formData.value.expiry_warning_hours > 0 &&
         batchOutputQuantity.value > 0 &&
         formData.value.conversion_rates.length > 0 &&
         formData.value.conversion_rates.every(r => 
           r.source_ingredient_id && 
           r.source_quantity > 0 &&
           r.source_unit
         )
})

const estimatedCost = computed(() => {
  let total = 0
  for (const rate of formData.value.conversion_rates) {
    const ingredient = ingredients.value.find(i => i.id === rate.source_ingredient_id)
    if (ingredient && ingredient.cost_per_unit && rate.source_unit) {
      // We need to convert source_quantity (in source_unit) to ingredient stock unit
      // Example: 200g → kg, we need to know how many kg is 200g
      // getConversionRate(stockUnit, recipeUnit) gives us the multiplier
      // But we're going from recipe → stock, so we need the inverse relationship
      
      // Get the conversion rate: how to convert FROM source_unit TO stock unit
      const conversionRate = getConversionRate(rate.source_unit, ingredient.unit)
      
      // Calculate quantity in stock unit
      const quantityInStockUnit = rate.source_quantity * conversionRate
      
      // Apply wastage
      const wastageMultiplier = 1 + ((rate.wastage_rate_percent || 0) / 100)
      const quantityNeeded = quantityInStockUnit * wastageMultiplier
      
      // Calculate cost
      total += quantityNeeded * ingredient.cost_per_unit
    }
  }
  return total
})

const addConversionRate = () => {
  formData.value.conversion_rates.push({
    source_ingredient_id: '',
    source_ingredient_name: '',
    source_quantity: 0,
    source_unit: '',
    batch_quantity: 0, // Will be set to batchOutputQuantity when saving
    wastage_rate_percent: 0,
    compatibleUnits: [],
    conversionRate: 1
  })
}

const removeConversionRate = (index) => {
  formData.value.conversion_rates.splice(index, 1)
}

const updateIngredientName = (index) => {
  const rate = formData.value.conversion_rates[index]
  const ingredient = ingredients.value.find(i => i.id === rate.source_ingredient_id)
  if (ingredient) {
    rate.source_ingredient_name = ingredient.name
    rate.source_unit = ingredient.unit
    rate.compatibleUnits = getCompatibleUnits(ingredient.unit)
    rate.conversionRate = 1 // Default to 1 when first selected
  }
}

const updateConversionRate = (index) => {
  const rate = formData.value.conversion_rates[index]
  const ingredient = ingredients.value.find(i => i.id === rate.source_ingredient_id)
  if (ingredient && rate.source_unit) {
    rate.conversionRate = getConversionRate(ingredient.unit, rate.source_unit)
  }
}

const getSelectedIngredient = (ingredientId) => {
  if (!ingredientId) return null
  return ingredients.value.find(i => i.id === ingredientId)
}

const formatCurrency = (value) => {
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND'
  }).format(value)
}

const save = async () => {
  if (!isValid.value) return

  isSubmitting.value = true
  
  try {
    // Set batch_quantity for all conversion rates to the same output quantity
    const payload = {
      ...formData.value,
      conversion_rates: formData.value.conversion_rates.map(r => ({
        ...r,
        batch_quantity: batchOutputQuantity.value, // All ingredients contribute to same batch output
        wastage_rate: (r.wastage_rate_percent || 0) / 100
      }))
    }

    let success
    if (props.mode === 'edit') {
      success = await batchStore.updateDefinition(props.definition.id, payload)
    } else {
      success = await batchStore.createDefinition(payload)
    }

    if (success) {
      emit('saved')
    } else {
      alert(batchStore.error || 'Lỗi lưu batch definition')
    }
  } catch (error) {
    alert('Lỗi: ' + (error.message || 'Không thể lưu'))
  } finally {
    isSubmitting.value = false
  }
}

// Initialize form data
watch(() => props.definition, (newDef) => {
  if (newDef && props.mode === 'edit') {
    // Get batch output quantity from first conversion rate (they should all be the same)
    batchOutputQuantity.value = newDef.conversion_rates?.[0]?.batch_quantity || 0
    
    formData.value = {
      name: newDef.name || '',
      unit: newDef.unit || '',
      shelf_life_hours: newDef.shelf_life_hours || 24,
      low_stock_threshold: newDef.low_stock_threshold || 0,
      expiry_warning_hours: newDef.expiry_warning_hours || 4,
      conversion_rates: (newDef.conversion_rates || []).map(r => {
        const ingredient = ingredients.value.find(i => i.id === r.source_ingredient_id)
        const stockUnit = ingredient?.unit || r.source_unit
        return {
          ...r,
          wastage_rate_percent: (r.wastage_rate || 0) * 100,
          compatibleUnits: getCompatibleUnits(stockUnit),
          conversionRate: getConversionRate(stockUnit, r.source_unit)
        }
      })
    }
  }
}, { immediate: true })

onMounted(() => {
  ingredientStore.fetchIngredients()
  
  // Add one conversion rate by default for create mode
  if (props.mode === 'create' && formData.value.conversion_rates.length === 0) {
    addConversionRate()
  }
})
</script>
