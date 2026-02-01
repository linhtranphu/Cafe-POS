<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3">
        <div class="flex items-center justify-between">
          <button @click="goBack" class="text-2xl text-gray-600">←</button>
          <h1 class="text-xl font-bold text-gray-800">{{ isEditing ? '✏️ Cập nhật thiết bị' : '➕ Thêm thiết bị mới' }}</h1>
          <div class="w-8"></div>
        </div>
      </div>
    </div>

    <!-- Scrollable Content -->
    <div class="flex-1 overflow-y-auto px-4 py-6 space-y-5">
      <!-- Tên thiết bị -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-3">Tên thiết bị *</label>
        <input v-model="formData.name" type="text" 
          class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
      </div>

      <!-- Loại & Số lượng - Responsive Grid -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-3">Loại *</label>
          <select v-model="formData.type" 
            class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent">
            <option value="">Chọn loại</option>
            <option v-for="cat in facilityCategories" :key="cat" :value="cat">{{ cat }}</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-3">Số lượng *</label>
          <input v-model.number="formData.quantity" type="number" 
            class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
        </div>
      </div>

      <!-- Khu vực & Trạng thái - Responsive Grid -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-3">Khu vực *</label>
          <input v-model="formData.area" type="text" placeholder="VD: Quầy bar"
            class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-3">Trạng thái *</label>
          <select v-model="formData.status" 
            class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent">
            <option v-for="option in FACILITY_STATUS_OPTIONS" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </div>
      </div>

      <!-- Ngày mua & Giá trị - Responsive Grid -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-3">Ngày mua</label>
          <input v-model="formData.purchase_date" type="date" 
            class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-3">Giá trị (VND)</label>
          <input v-model.number="formData.cost" type="number" placeholder="0"
            class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
        </div>
      </div>

      <!-- Nhà cung cấp -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-3">Nhà cung cấp</label>
        <input v-model="formData.supplier" type="text" 
          class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
      </div>

      <!-- Ghi chú -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-3">Ghi chú</label>
        <textarea v-model="formData.notes" rows="3" 
          class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"></textarea>
      </div>

      <!-- Auto-Expense Indicator -->
      <div v-if="!isEditing && formData.cost > 0" 
        class="bg-green-50 border border-green-200 rounded-lg p-4 mt-6">
        <div class="flex items-start gap-3">
          <span class="text-green-600 text-2xl flex-shrink-0">✅</span>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-green-800">Tự động ghi nhận chi phí</p>
            <p class="text-xs text-green-600 mt-2 break-words">
              Hệ thống sẽ tự động tạo chi phí: <span class="font-semibold">{{ formatPrice(formData.cost) }}</span>
            </p>
            <p class="text-xs text-green-600 mt-1">Danh mục: Cơ sở vật chất</p>
          </div>
        </div>
      </div>

      <!-- Spacer for bottom buttons -->
      <div class="h-24"></div>
    </div>

    <!-- Fixed Footer -->
    <div class="flex-shrink-0 bg-white px-4 py-4 border-t flex gap-3 pb-safe">
      <button @click="goBack" 
        class="flex-1 bg-gray-200 text-gray-700 py-4 rounded-xl font-medium text-base active:bg-gray-300 transition-colors">
        Hủy
      </button>
      <button @click="saveFacility" 
        class="flex-1 bg-blue-500 text-white py-4 rounded-xl font-medium text-base active:bg-blue-600 transition-colors">
        {{ isEditing ? 'Cập nhật' : 'Thêm mới' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useFacilityStore } from '../stores/facility'
import { 
  FACILITY_STATUS, 
  FACILITY_STATUS_OPTIONS,
  FACILITY_TYPE_OPTIONS
} from '../constants/facility'
import { 
  formatPrice,
  sanitizeFormData,
  parseBackendData
} from '../utils/formatters'

const router = useRouter()
const route = useRoute()
const facilityStore = useFacilityStore()

const isEditing = ref(false)
const currentFacility = ref(null)

const formData = ref({
  name: '',
  type: '',
  area: '',
  quantity: 1,
  status: FACILITY_STATUS.IN_USE,
  purchase_date: '',
  cost: 0,
  supplier: '',
  notes: ''
})

// Facility categories from constants + custom categories from backend
const facilityCategories = computed(() => {
  const defaultCategories = FACILITY_TYPE_OPTIONS.map(opt => opt.label)
  const backendTypes = facilityStore.types.map(t => t.name)
  return [...new Set([...defaultCategories, ...backendTypes])]
})

const goBack = async () => {
  try {
    await router.back()
  } catch (error) {
    console.error('Navigation error:', error)
    // Fallback to facilities page if back fails
    await router.push('/facilities')
  }
}

const saveFacility = async () => {
  try {
    // Sanitize form data before sending to backend
    const dataToSend = sanitizeFormData(formData.value, {
      name: { type: 'string' },
      type: { type: 'string' },
      area: { type: 'string' },
      quantity: { type: 'number', default: 1 },
      status: { type: 'string' },
      purchase_date: { type: 'date', default: new Date().toISOString() },
      cost: { type: 'number', default: 0 },
      supplier: { type: 'string', default: '' },
      notes: { type: 'string', default: '' }
    })
    
    if (isEditing.value) {
      await facilityStore.updateFacility(currentFacility.value.id, dataToSend)
      alert('Cập nhật thiết bị thành công')
    } else {
      await facilityStore.createFacility(dataToSend)
      alert('Thêm thiết bị thành công')
    }
    
    // Navigate back after successful save
    try {
      await router.back()
    } catch (error) {
      console.error('Navigation error:', error)
      await router.push('/facilities')
    }
  } catch (error) {
    console.error('Error saving facility:', error)
    console.error('Error response:', error.response?.data)
    const errorMessage = error.response?.data?.error || 'Có lỗi xảy ra khi lưu thiết bị'
    alert(errorMessage)
  }
}

onMounted(async () => {
  // Check if we're editing an existing facility
  if (route.params.id) {
    isEditing.value = true
    // Get facility from route params or fetch it
    const facilityId = route.params.id
    const facility = facilityStore.items.find(f => f.id === facilityId)
    
    if (facility) {
      currentFacility.value = facility
      // Parse backend data for form display (converts ISO dates to local format)
      formData.value = parseBackendData({ ...facility }, {
        purchase_date: { type: 'date' }
      })
    }
  }
  
  // Ensure facility types are loaded
  if (facilityStore.types.length === 0) {
    await facilityStore.fetchFacilityTypes()
  }
})
</script>

<style scoped>
.pb-safe {
  padding-bottom: max(1rem, env(safe-area-inset-bottom));
}
</style>
