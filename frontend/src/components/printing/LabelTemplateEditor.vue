<template>
  <div class="h-full flex flex-col bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm flex-shrink-0 px-4 py-3">
      <div class="flex items-center justify-between mb-3">
        <h2 class="text-lg font-bold text-gray-800">🏷️ Label Template (TSPL)</h2>
        <div class="flex gap-2">
          <button
            @click="loadTemplate"
            class="bg-gray-500 text-white px-4 py-2 rounded-lg text-sm font-bold hover:bg-gray-600 active:scale-95 transition-transform"
          >
            🔄 Reload
          </button>
          <button
            @click="saveTemplate"
            :disabled="saving"
            class="bg-blue-500 text-white px-4 py-2 rounded-lg text-sm font-bold hover:bg-blue-600 active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
          >
            <span v-if="saving" class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
            <span>💾 Lưu Template</span>
          </button>
        </div>
      </div>
      
      <div class="text-sm text-gray-600">
        <p class="font-semibold text-blue-700">Template TSPL cho in tem. Mỗi món sẽ có một tem riêng.</p>
        <p class="mt-1">Hiển thị: Mã order, Tên món + Note, Thời gian, Tên khách</p>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-hidden flex">
      <!-- Left: TSPL Editor -->
      <div class="flex-1 flex flex-col border-r">
        <div class="bg-blue-100 px-4 py-2 border-b">
          <h3 class="font-bold text-sm text-blue-800">📝 TSPL Commands</h3>
        </div>
        <div class="flex-1 overflow-y-auto p-4">
          <textarea
            v-model="tsplContent"
            class="w-full h-full font-mono text-sm border-2 border-blue-300 rounded-lg p-4 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            placeholder="Nhập TSPL commands..."
          ></textarea>
        </div>
      </div>

      <!-- Right: Preview & Test -->
      <div class="flex-1 flex flex-col bg-gray-100">
        <div class="bg-gray-200 px-4 py-2 border-b flex items-center justify-between">
          <h3 class="font-bold text-sm text-gray-700">👁️ Preview</h3>
          <button
            @click="toggleSampleData"
            class="text-xs bg-white px-3 py-1 rounded border hover:bg-gray-50"
          >
            {{ showSampleData ? '📋 Hide Data' : '📋 Show Data' }}
          </button>
        </div>
        
        <!-- Sample Data Panel -->
        <div v-if="showSampleData" class="bg-blue-50 border-b p-3 text-xs">
          <div class="font-bold mb-2">Dữ liệu mẫu:</div>
          <div class="space-y-1 text-gray-700">
            <div>OrderNumber: {{ sampleData.OrderNumber }}</div>
            <div>ItemName: {{ sampleData.ItemName }}</div>
            <div>Note: {{ sampleData.Note }}</div>
            <div>Time: {{ sampleData.Time }}</div>
            <div>CustomerName: {{ sampleData.CustomerName }}</div>
          </div>
        </div>
        
        <div class="flex-1 overflow-y-auto p-4">
          <!-- Preview Mock Label -->
          <div class="bg-white rounded-lg shadow-lg p-4">
            <div class="border-2 border-blue-300 bg-white p-4" style="width: 400px; margin: 0 auto;">
              <div class="font-mono text-xs whitespace-pre-wrap bg-gray-50 p-4 rounded">
                {{ previewText }}
              </div>
            </div>
            
            <div class="mt-4 text-xs text-gray-500 text-center">
              <p class="text-blue-700 font-semibold">Preview với dữ liệu mẫu</p>
              <p>Kích thước: {{ labelSize }}</p>
            </div>
          </div>
          
          <!-- Test Print Section -->
          <div class="mt-4 bg-white rounded-lg shadow p-4">
            <h4 class="font-bold text-sm mb-3">🖨️ Test Print</h4>
            <div class="space-y-2">
              <!-- Printer selector from Printers List -->
              <div v-if="labelPrinters.length > 0">
                <label class="block text-xs text-gray-500 mb-1">Chọn máy in tem:</label>
                <select
                  v-model="selectedPrinterId"
                  class="w-full px-3 py-2 border rounded text-sm"
                >
                  <option v-for="p in labelPrinters" :key="p.id" :value="p.id">
                    {{ p.name }} — {{ p.ip_address }}:{{ p.port }}
                  </option>
                </select>
              </div>
              <div v-else class="text-xs text-orange-600 bg-orange-50 p-2 rounded">
                ⚠️ Chưa có máy in tem. Vào tab "Máy In" để thêm máy in loại LABEL.
              </div>
              <button
                @click="testPrint"
                :disabled="testPrinting || !selectedPrinterId"
                class="w-full bg-green-500 text-white px-4 py-2 rounded-lg text-sm font-bold hover:bg-green-600 active:scale-95 transition-transform disabled:opacity-50"
              >
                <span v-if="testPrinting">⏳ Đang in...</span>
                <span v-else>🖨️ Test Print</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Status Message -->
    <div v-if="statusMessage" :class="[
      'px-4 py-3 border-t',
      statusMessage.type === 'success' ? 'bg-green-50 text-green-800' : 'bg-red-50 text-red-800'
    ]">
      <div class="flex items-center gap-2">
        <span>{{ statusMessage.type === 'success' ? '✅' : '❌' }}</span>
        <span class="font-bold">{{ statusMessage.text }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '../../services/api'
import { usePrinterConfigStore } from '../../stores/printerConfig'

const printerStore = usePrinterConfigStore()

const tsplContent = ref('')
const saving = ref(false)
const testPrinting = ref(false)
const showSampleData = ref(false)
const statusMessage = ref(null)
const selectedPrinterId = ref(null)

const labelPrinters = computed(() => printerStore.printers.filter(p => p.type === 'LABEL' && p.is_enabled))

const sampleData = ref({
  OrderNumber: '20260313-001',
  ItemName: 'Cà phê sữa đá',
  Note: 'Ít đường',
  Time: '14:30',
  CustomerName: 'Nguyễn Văn A'
})

const labelSize = computed(() => {
  const printer = labelPrinters.value.find(p => p.id === selectedPrinterId.value)
  if (printer?.paper_width) {
    return `${printer.paper_width}mm x 30mm`
  }
  const match = tsplContent.value.match(/SIZE\s+(\d+)\s+mm,\s+(\d+)\s+mm/)
  if (match) return `${match[1]}mm x ${match[2]}mm`
  return '40mm x 30mm (default)'
})

const previewText = computed(() => {
  let preview = tsplContent.value
  
  // Replace template variables with sample data
  preview = preview.replace(/\{\{\.OrderNumber\}\}/g, sampleData.value.OrderNumber)
  preview = preview.replace(/\{\{\.ItemName\}\}/g, sampleData.value.ItemName)
  preview = preview.replace(/\{\{\.Note\}\}/g, sampleData.value.Note)
  preview = preview.replace(/\{\{\.Time\}\}/g, sampleData.value.Time)
  preview = preview.replace(/\{\{\.CustomerName\}\}/g, sampleData.value.CustomerName)
  
  // Handle conditionals (simple)
  if (sampleData.value.Note) {
    preview = preview.replace(/\{\{if \.Note\}\}([\s\S]*?)\{\{end\}\}/g, '$1')
  } else {
    preview = preview.replace(/\{\{if \.Note\}\}[\s\S]*?\{\{end\}\}/g, '')
  }
  
  return preview || 'No template loaded'
})

const toggleSampleData = () => {
  showSampleData.value = !showSampleData.value
}

const loadTemplate = async () => {
  try {
    const response = await api.get('/manager/label-templates/order-item')
    tsplContent.value = response.data.content
    
    statusMessage.value = {
      type: 'success',
      text: 'Label template loaded successfully'
    }
    setTimeout(() => statusMessage.value = null, 3000)
  } catch (error) {
    console.error('Load error:', error)
    statusMessage.value = {
      type: 'error',
      text: 'Failed to load template: ' + (error.response?.data?.error || error.message)
    }
    setTimeout(() => statusMessage.value = null, 5000)
  }
}

const saveTemplate = async () => {
  saving.value = true
  statusMessage.value = null
  
  try {
    await api.put('/manager/label-templates/order-item', {
      content: tsplContent.value
    })
    
    statusMessage.value = {
      type: 'success',
      text: 'Label template saved successfully'
    }
    
    setTimeout(() => statusMessage.value = null, 3000)
  } catch (error) {
    console.error('Save error:', error)
    statusMessage.value = {
      type: 'error',
      text: 'Failed to save template: ' + (error.response?.data?.error || error.message)
    }
  } finally {
    saving.value = false
  }
}

const testPrint = async () => {
  const printer = labelPrinters.value.find(p => p.id === selectedPrinterId.value)
  if (!printer) return

  testPrinting.value = true
  statusMessage.value = null
  
  try {
    await api.post('/manager/label-templates/test-print', {
      item_name: sampleData.value.ItemName,
      note: sampleData.value.Note,
      customer_name: sampleData.value.CustomerName,
      printer_ip: printer.ip_address,
      port: printer.port || 9100,
      paper_width: printer.paper_width || 40
    })
    
    statusMessage.value = {
      type: 'success',
      text: `Test print thành công! Kiểm tra máy in "${printer.name}".`
    }
    setTimeout(() => statusMessage.value = null, 5000)
  } catch (error) {
    console.error('Test print error:', error)
    statusMessage.value = {
      type: 'error',
      text: 'Lỗi in: ' + (error.response?.data?.error || error.message)
    }
  } finally {
    testPrinting.value = false
  }
}

onMounted(async () => {
  await Promise.all([
    loadTemplate(),
    printerStore.fetchPrinters()
  ])
  
  // Auto-select default label printer
  const defaultPrinter = printerStore.printers.find(p => p.type === 'LABEL' && p.is_default && p.is_enabled)
    || printerStore.printers.find(p => p.type === 'LABEL' && p.is_enabled)
  if (defaultPrinter) {
    selectedPrinterId.value = defaultPrinter.id
  }
})
</script>
