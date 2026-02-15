<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3" style="padding-top: max(0.75rem, env(safe-area-inset-top))">
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-3">
            <button @click="$router.push('/batch')" class="text-gray-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            <h1 class="text-xl font-bold text-gray-800">🧪 Định Nghĩa Batch</h1>
          </div>
        </div>
        
        <!-- Search Bar -->
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Tìm kiếm batch..."
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24">
      <!-- Create Button -->
      <button 
        @click="openCreateModal"
        class="w-full bg-gradient-to-br from-blue-500 to-cyan-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform mb-4">
        <div class="text-2xl mb-1">➕</div>
        <div class="text-sm font-bold">Tạo Batch Definition Mới</div>
      </button>

      <!-- Loading State -->
      <div v-if="loading" class="text-center py-16">
        <div class="inline-block w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
        <p class="text-gray-500 mt-4">Đang tải...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredDefinitions.length === 0" class="text-center py-16">
        <div class="text-6xl mb-4">📭</div>
        <p class="text-gray-500">Không có batch definition nào</p>
      </div>

      <!-- Definitions List -->
      <div v-else class="space-y-3">
        <div 
          v-for="definition in filteredDefinitions" 
          :key="definition?.id || Math.random()"
          class="bg-white rounded-2xl p-4 shadow-sm">
          
          <!-- Header -->
          <div class="flex justify-between items-start mb-3">
            <div>
              <h3 class="font-bold text-lg">{{ definition?.name || 'N/A' }}</h3>
              <p class="text-sm text-gray-600">{{ definition?.unit || '' }}</p>
            </div>
          </div>

          <!-- Info -->
          <div class="mb-3 space-y-2 text-sm">
            <div class="flex justify-between items-center">
              <span class="text-gray-600">⏱️ Thời hạn:</span>
              <span class="font-medium">{{ definition?.shelf_life_hours || 0 }}h</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-gray-600">⚠️ Ngưỡng thấp:</span>
              <span class="font-medium">{{ definition?.low_stock_threshold || 0 }} {{ definition?.unit || '' }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-gray-600">🔔 Cảnh báo hết hạn:</span>
              <span class="font-medium">{{ definition?.expiry_warning_hours || 0 }}h trước</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-gray-600">📦 Nguyên liệu:</span>
              <span class="font-medium">{{ definition?.conversion_rates?.length || 0 }} loại</span>
            </div>
          </div>

          <!-- Actions -->
          <div class="grid grid-cols-2 gap-2 pt-3 border-t">
            <button 
              @click="openEditModal(definition)"
              :disabled="!definition?.id"
              class="bg-blue-500 text-white py-2 rounded-lg text-xs font-bold active:bg-blue-600 flex items-center justify-center gap-1 disabled:opacity-50 disabled:cursor-not-allowed">
              <span>✏️</span>
              <span>Sửa</span>
            </button>
            <button 
              @click="deleteDefinition(definition)"
              :disabled="!definition?.id"
              class="bg-red-500 text-white py-2 rounded-lg text-xs font-bold active:bg-red-600 flex items-center justify-center gap-1 disabled:opacity-50 disabled:cursor-not-allowed">
              <span>🗑️</span>
              <span>Xóa</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />

    <!-- Create/Edit Modal -->
    <BatchDefinitionForm
      v-if="showModal"
      :definition="selectedDefinition"
      :mode="modalMode"
      @close="closeModal"
      @saved="handleSaved"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useBatchDefinitionStore } from '../../stores/batchDefinition'
import BottomNav from '../BottomNav.vue'
import BatchDefinitionForm from './BatchDefinitionForm.vue'

const batchStore = useBatchDefinitionStore()

const searchQuery = ref('')
const showModal = ref(false)
const selectedDefinition = ref(null)
const modalMode = ref('create')

const loading = computed(() => batchStore.loading)
const definitions = computed(() => batchStore.definitions)

const filteredDefinitions = computed(() => {
  if (!searchQuery.value) return definitions.value || []
  
  const query = searchQuery.value.toLowerCase()
  return (definitions.value || []).filter(d => 
    d && d.name?.toLowerCase().includes(query)
  )
})

const openCreateModal = () => {
  selectedDefinition.value = null
  modalMode.value = 'create'
  showModal.value = true
}

const openEditModal = (definition) => {
  selectedDefinition.value = definition
  modalMode.value = 'edit'
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  selectedDefinition.value = null
}

const handleSaved = () => {
  closeModal()
  batchStore.fetchDefinitions()
}

const deleteDefinition = async (definition) => {
  if (!definition || !definition.id) {
    alert('Không thể xóa: Batch definition không hợp lệ')
    return
  }
  
  if (!confirm(`Xóa batch definition "${definition.name}"?`)) return
  
  const success = await batchStore.deleteDefinition(definition.id)
  if (success) {
    alert('Đã xóa thành công')
  } else {
    alert(batchStore.error || 'Lỗi xóa batch definition')
  }
}

onMounted(() => {
  batchStore.fetchDefinitions()
})
</script>
