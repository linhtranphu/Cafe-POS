<template>
  <div class="h-full flex flex-col bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm flex-shrink-0 px-4 py-3">
      <div class="flex items-center justify-between mb-3">
        <h2 class="text-lg font-bold text-gray-800">🖨️ Máy In</h2>
        <button
          @click="openCreateModal"
          class="bg-blue-500 text-white px-4 py-2 rounded-lg text-sm font-bold hover:bg-blue-600 active:scale-95 transition-transform"
        >
          ➕ Thêm
        </button>
      </div>

      <!-- Type Filter -->
      <div class="flex gap-2">
        <button
          v-for="filter in typeFilters"
          :key="filter.value"
          @click="selectedType = filter.value"
          :class="[
            'px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-colors',
            selectedType === filter.value
              ? 'bg-blue-500 text-white'
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          ]"
        >
          {{ filter.label }}
        </button>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4">
      <!-- Loading State -->
      <div v-if="loading" class="text-center py-16">
        <div class="inline-block w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
        <p class="text-gray-500 mt-4">Đang tải...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-16">
        <div class="text-6xl mb-4">⚠️</div>
        <p class="text-red-500 mb-4">{{ error }}</p>
        <button
          @click="refreshPrinters"
          class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
        >
          Thử lại
        </button>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredPrinters.length === 0" class="text-center py-16">
        <div class="text-6xl mb-4">🖨️</div>
        <p class="text-gray-500 mb-2">Chưa có máy in nào</p>
        <button
          @click="openCreateModal"
          class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 text-sm"
        >
          Thêm máy in đầu tiên
        </button>
      </div>

      <!-- Printers List -->
      <div v-else class="space-y-3">
        <div
          v-for="printer in filteredPrinters"
          :key="printer.id"
          class="bg-white rounded-2xl p-4 shadow-sm"
        >
          <!-- Header -->
          <div class="flex justify-between items-start mb-3">
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-1">
                <span class="text-2xl">{{ getPrinterIcon(printer.type) }}</span>
                <h3 class="font-bold text-lg">{{ printer.name }}</h3>
                <span
                  v-if="printer.is_default"
                  class="px-2 py-1 bg-yellow-100 text-yellow-800 text-xs font-bold rounded-full"
                >
                  ⭐ Mặc định
                </span>
              </div>
              <p class="text-sm text-gray-600">{{ getPrinterTypeLabel(printer.type) }}</p>
            </div>
            <div class="flex items-center gap-2">
              <!-- Status Indicator -->
              <div
                :class="[
                  'w-3 h-3 rounded-full',
                  printer.is_enabled ? 'bg-green-500' : 'bg-gray-400'
                ]"
                :title="printer.is_enabled ? 'Đang hoạt động' : 'Đã tắt'"
              ></div>
            </div>
          </div>

          <!-- Connection Info -->
          <div class="mb-3 space-y-2 text-sm">
            <div class="flex justify-between items-center">
              <span class="text-gray-600">🔗 Kết nối:</span>
              <span class="font-medium">{{ getConnectionLabel(printer.connection_type) }}</span>
            </div>
            <div v-if="printer.connection_type === 'NETWORK'" class="flex justify-between items-center">
              <span class="text-gray-600">🌐 Địa chỉ:</span>
              <span class="font-medium">{{ printer.ip_address }}:{{ printer.port }}</span>
            </div>
            <div v-if="printer.connection_type === 'USB'" class="flex justify-between items-center">
              <span class="text-gray-600">🔌 USB:</span>
              <span class="font-medium text-xs">{{ printer.usb_path }}</span>
            </div>
            <div v-if="printer.type === 'BILL'" class="flex justify-between items-center">
              <span class="text-gray-600">📏 Khổ giấy:</span>
              <span class="font-medium">{{ printer.paper_width }}mm</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-gray-600">⚙️ Trạng thái:</span>
              <span :class="[
                'font-medium',
                printer.is_enabled ? 'text-green-600' : 'text-gray-500'
              ]">
                {{ printer.is_enabled ? 'Đang hoạt động' : 'Đã tắt' }}
              </span>
            </div>
          </div>

          <!-- Test Result (if testing this printer) -->
          <div
            v-if="testingPrinterId === printer.id && testResult"
            :class="[
              'mb-3 p-2 rounded-lg text-xs',
              testResult.success ? 'bg-green-50 text-green-800' : 'bg-red-50 text-red-800'
            ]"
          >
            <div class="flex items-center gap-2">
              <span>{{ testResult.success ? '✅' : '❌' }}</span>
              <span class="font-bold">{{ testResult.message }}</span>
            </div>
          </div>

          <!-- Actions -->
          <div class="grid grid-cols-3 gap-2 pt-3 border-t">
            <button
              @click="handleTestConnection(printer.id)"
              :disabled="testingPrinterId === printer.id"
              class="bg-purple-500 text-white py-2 rounded-lg text-xs font-bold hover:bg-purple-600 active:scale-95 transition-transform flex items-center justify-center gap-1 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="testingPrinterId === printer.id" class="inline-block w-3 h-3 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
              <span v-else>🔍</span>
              <span>{{ testingPrinterId === printer.id ? 'Test...' : 'Test' }}</span>
            </button>
            <button
              @click="openEditModal(printer)"
              class="bg-blue-500 text-white py-2 rounded-lg text-xs font-bold hover:bg-blue-600 active:scale-95 transition-transform flex items-center justify-center gap-1"
            >
              <span>✏️</span>
              <span>Sửa</span>
            </button>
            <button
              @click="handleDelete(printer)"
              class="bg-red-500 text-white py-2 rounded-lg text-xs font-bold hover:bg-red-600 active:scale-95 transition-transform flex items-center justify-center gap-1"
            >
              <span>🗑️</span>
              <span>Xóa</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <PrinterConfigForm
      v-if="showModal"
      :printer="selectedPrinter"
      :mode="modalMode"
      @close="closeModal"
      @saved="handleSaved"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { usePrinterConfigStore } from '../../stores/printerConfig'
import { usePrinterStatus } from '../../composables/usePrinterStatus'
import PrinterConfigForm from './PrinterConfigForm.vue'

const printerStore = usePrinterConfigStore()

// Setup printer status monitoring
const { printerStatuses, getPrinterStatus } = usePrinterStatus()

const selectedType = ref('all')
const showModal = ref(false)
const selectedPrinter = ref(null)
const modalMode = ref('create')
const testingPrinterId = ref(null)
const testResult = ref(null)

const typeFilters = [
  { label: 'Tất cả', value: 'all' },
  { label: '🧾 Bill', value: 'BILL' },
  { label: '🏷️ Tem', value: 'LABEL' }
]

const loading = computed(() => printerStore.loading)
const error = computed(() => printerStore.error)

const filteredPrinters = computed(() => {
  if (selectedType.value === 'all') {
    return printerStore.printers
  }
  return printerStore.printersByType(selectedType.value)
})

const getPrinterIcon = (type) => {
  return type === 'BILL' ? '🧾' : '🏷️'
}

const getPrinterTypeLabel = (type) => {
  return type === 'BILL' ? 'Máy in Bill' : 'Máy in Tem'
}

const getConnectionLabel = (type) => {
  return type === 'NETWORK' ? 'Network (TCP/IP)' : 'USB'
}

const openCreateModal = () => {
  selectedPrinter.value = null
  modalMode.value = 'create'
  showModal.value = true
}

const openEditModal = (printer) => {
  selectedPrinter.value = printer
  modalMode.value = 'edit'
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  selectedPrinter.value = null
}

const handleSaved = () => {
  closeModal()
  refreshPrinters()
}

const refreshPrinters = async () => {
  await printerStore.fetchPrinters()
}

const handleTestConnection = async (printerId) => {
  testingPrinterId.value = printerId
  testResult.value = null

  try {
    const result = await printerStore.testConnection(printerId)
    testResult.value = result

    // Clear test result after 5 seconds
    setTimeout(() => {
      if (testingPrinterId.value === printerId) {
        testResult.value = null
        testingPrinterId.value = null
      }
    }, 5000)
  } catch (err) {
    testResult.value = {
      success: false,
      message: err.message || 'Lỗi kiểm tra kết nối'
    }
  } finally {
    testingPrinterId.value = null
  }
}

const handleDelete = async (printer) => {
  if (!confirm(`Xóa máy in "${printer.name}"?`)) return

  const success = await printerStore.deletePrinter(printer.id)
  if (success) {
    alert('Đã xóa máy in')
  } else {
    alert(printerStore.error || 'Lỗi xóa máy in')
  }
}

onMounted(() => {
  refreshPrinters()
})
</script>
