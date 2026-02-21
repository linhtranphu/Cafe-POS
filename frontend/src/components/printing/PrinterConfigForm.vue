<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
    <div class="bg-gray-50 w-full h-screen flex flex-col">
      <!-- Mobile Header - Fixed -->
      <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
        <div class="px-4 py-3">
          <div class="flex items-center justify-between">
            <button @click="$emit('close')" class="text-2xl text-gray-600">←</button>
            <h1 class="text-xl font-bold text-gray-800">
              {{ mode === 'edit' ? '✏️ Cập nhật Máy In' : '➕ Thêm Máy In' }}
            </h1>
            <div class="w-8"></div>
          </div>
        </div>
      </div>

      <!-- Scrollable Content -->
      <div class="flex-1 overflow-y-auto px-4 py-6 space-y-5">
        
        <!-- SECTION 1: THÔNG TIN CƠ BẢN -->
        <div class="bg-gradient-to-br from-blue-50 to-cyan-50 rounded-2xl p-4 border-2 border-blue-200">
          <div class="flex items-center gap-2 mb-4">
            <div class="text-2xl">🖨️</div>
            <h2 class="text-lg font-bold text-blue-900">Thông Tin Cơ Bản</h2>
          </div>

          <div class="space-y-4">
            <!-- Tên máy in -->
            <div>
              <label class="block text-sm font-bold text-gray-700 mb-2">📝 Tên Máy In *</label>
              <input 
                v-model="formData.name" 
                type="text" 
                placeholder="VD: Máy in Bill 1"
                class="w-full px-4 py-3 border-2 border-blue-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white" 
              />
            </div>

            <!-- Loại máy in -->
            <div>
              <label class="block text-sm font-bold text-gray-700 mb-2">🏷️ Loại Máy In *</label>
              <select 
                v-model="formData.type" 
                class="w-full px-4 py-3 border-2 border-blue-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white">
                <option value="">Chọn loại</option>
                <option value="BILL">🧾 Bill (Hóa đơn)</option>
                <option value="LABEL">🏷️ Label (Tem nhãn)</option>
              </select>
            </div>

            <!-- Paper Width (chỉ cho BILL) -->
            <div v-if="formData.type === 'BILL'">
              <label class="block text-sm font-bold text-gray-700 mb-2">📏 Khổ Giấy *</label>
              <select 
                v-model.number="formData.paper_width" 
                class="w-full px-4 py-3 border-2 border-blue-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white">
                <option :value="58">58mm (Nhỏ)</option>
                <option :value="80">80mm (Lớn)</option>
              </select>
            </div>
          </div>
        </div>

        <!-- SECTION 2: KẾT NỐI -->
        <div class="bg-gradient-to-br from-purple-50 to-pink-50 rounded-2xl p-4 border-2 border-purple-200">
          <div class="flex items-center gap-2 mb-4">
            <div class="text-2xl">🔌</div>
            <h2 class="text-lg font-bold text-purple-900">Cấu Hình Kết Nối</h2>
          </div>

          <div class="space-y-4">
            <!-- Connection Type -->
            <div>
              <label class="block text-sm font-bold text-gray-700 mb-2">🔗 Loại Kết Nối *</label>
              <select 
                v-model="formData.connection_type" 
                class="w-full px-4 py-3 border-2 border-purple-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 bg-white">
                <option value="">Chọn loại kết nối</option>
                <option value="NETWORK">🌐 Network (TCP/IP)</option>
                <option value="USB">🔌 USB</option>
              </select>
            </div>

            <!-- Network Connection Details -->
            <div v-if="formData.connection_type === 'NETWORK'" class="space-y-3">
              <div>
                <label class="block text-sm font-bold text-gray-700 mb-2">🌐 Địa Chỉ IP *</label>
                <input 
                  v-model="formData.ip_address" 
                  type="text" 
                  placeholder="VD: 192.168.1.100"
                  class="w-full px-4 py-3 border-2 border-purple-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 bg-white" 
                />
              </div>
              <div>
                <label class="block text-sm font-bold text-gray-700 mb-2">🔢 Port *</label>
                <input 
                  v-model.number="formData.port" 
                  type="number" 
                  placeholder="VD: 9100"
                  class="w-full px-4 py-3 border-2 border-purple-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 bg-white" 
                />
              </div>
            </div>

            <!-- USB Connection Details -->
            <div v-if="formData.connection_type === 'USB'">
              <label class="block text-sm font-bold text-gray-700 mb-2">🔌 USB Path *</label>
              <input 
                v-model="formData.usb_path" 
                type="text" 
                placeholder="VD: /dev/usb/lp0"
                class="w-full px-4 py-3 border-2 border-purple-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 bg-white" 
              />
              <p class="text-xs text-gray-500 mt-1">Đường dẫn thiết bị USB trên hệ thống</p>
            </div>

            <!-- Test Connection Button -->
            <button
              v-if="formData.connection_type"
              @click="handleTestConnection"
              :disabled="testingConnection || !canTestConnection"
              type="button"
              class="w-full bg-purple-600 text-white py-3 rounded-lg font-bold hover:bg-purple-700 active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              <span v-if="testingConnection" class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
              <span v-else>🔍</span>
              <span>{{ testingConnection ? 'Đang kiểm tra...' : 'Kiểm tra kết nối' }}</span>
            </button>

            <!-- Test Result -->
            <div v-if="testResult" :class="[
              'p-3 rounded-lg text-sm',
              testResult.success ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
            ]">
              <div class="flex items-center gap-2">
                <span>{{ testResult.success ? '✅' : '❌' }}</span>
                <span class="font-bold">{{ testResult.message }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- SECTION 3: CÀI ĐẶT -->
        <div class="bg-gradient-to-br from-green-50 to-emerald-50 rounded-2xl p-4 border-2 border-green-200">
          <div class="flex items-center gap-2 mb-4">
            <div class="text-2xl">⚙️</div>
            <h2 class="text-lg font-bold text-green-900">Cài Đặt</h2>
          </div>

          <div class="space-y-4">
            <!-- Set as Default -->
            <label class="flex items-center gap-3 p-3 bg-white rounded-lg border-2 border-green-200 cursor-pointer">
              <input 
                v-model="formData.is_default" 
                type="checkbox"
                class="w-5 h-5 text-green-600 rounded focus:ring-2 focus:ring-green-500"
              />
              <div>
                <div class="font-bold text-gray-800">⭐ Đặt làm máy in mặc định</div>
                <div class="text-xs text-gray-600">Tự động sử dụng máy in này cho loại {{ formData.type === 'BILL' ? 'bill' : 'tem' }}</div>
              </div>
            </label>

            <!-- Enable/Disable -->
            <label class="flex items-center gap-3 p-3 bg-white rounded-lg border-2 border-green-200 cursor-pointer">
              <input 
                v-model="formData.is_enabled" 
                type="checkbox"
                class="w-5 h-5 text-green-600 rounded focus:ring-2 focus:ring-green-500"
              />
              <div>
                <div class="font-bold text-gray-800">✅ Kích hoạt máy in</div>
                <div class="text-xs text-gray-600">Cho phép sử dụng máy in này</div>
              </div>
            </label>
          </div>
        </div>

        <!-- Error Message -->
        <div v-if="error" class="bg-red-50 border-2 border-red-200 rounded-lg p-4">
          <div class="flex items-center gap-2 text-red-800">
            <span class="text-xl">⚠️</span>
            <span class="font-bold">{{ error }}</span>
          </div>
        </div>
      </div>

      <!-- Bottom Action Buttons -->
      <div class="sticky bottom-0 bg-white border-t-2 border-gray-200 p-4 flex-shrink-0">
        <div class="grid grid-cols-2 gap-3">
          <button 
            @click="$emit('close')"
            type="button"
            class="bg-gray-200 text-gray-800 py-3 rounded-lg font-bold hover:bg-gray-300 active:scale-95 transition-transform">
            Hủy
          </button>
          <button 
            @click="handleSave"
            :disabled="saving || !isValid"
            class="bg-blue-600 text-white py-3 rounded-lg font-bold hover:bg-blue-700 active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2">
            <span v-if="saving" class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
            <span>{{ saving ? 'Đang lưu...' : 'Lưu' }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { usePrinterConfigStore } from '../../stores/printerConfig'

const props = defineProps({
  printer: {
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

const printerStore = usePrinterConfigStore()

const formData = ref({
  name: '',
  type: '',
  connection_type: '',
  ip_address: '',
  port: 9100,
  usb_path: '',
  paper_width: 80,
  is_default: false,
  is_enabled: true
})

const saving = ref(false)
const testingConnection = ref(false)
const testResult = ref(null)
const error = ref(null)

// Initialize form data if editing
if (props.mode === 'edit' && props.printer) {
  formData.value = {
    id: props.printer.id,
    name: props.printer.name || '',
    type: props.printer.type || '',
    connection_type: props.printer.connection_type || '',
    ip_address: props.printer.ip_address || '',
    port: props.printer.port || 9100,
    usb_path: props.printer.usb_path || '',
    paper_width: props.printer.paper_width || 80,
    is_default: props.printer.is_default || false,
    is_enabled: props.printer.is_enabled !== undefined ? props.printer.is_enabled : true
  }
}

const isValid = computed(() => {
  if (!formData.value.name || !formData.value.type || !formData.value.connection_type) {
    return false
  }

  if (formData.value.connection_type === 'NETWORK') {
    return formData.value.ip_address && formData.value.port
  }

  if (formData.value.connection_type === 'USB') {
    return formData.value.usb_path
  }

  return false
})

const canTestConnection = computed(() => {
  if (formData.value.connection_type === 'NETWORK') {
    return formData.value.ip_address && formData.value.port
  }
  if (formData.value.connection_type === 'USB') {
    return formData.value.usb_path
  }
  return false
})

// Clear test result when connection details change
watch(() => [formData.value.connection_type, formData.value.ip_address, formData.value.port, formData.value.usb_path], () => {
  testResult.value = null
  printerStore.clearTestResult()
})

const handleTestConnection = async () => {
  if (!canTestConnection.value) return

  // For new printers, we need to save first before testing
  if (!formData.value.id) {
    error.value = 'Vui lòng lưu cấu hình trước khi kiểm tra kết nối'
    return
  }

  testingConnection.value = true
  testResult.value = null
  error.value = null

  try {
    const result = await printerStore.testConnection(formData.value.id)
    testResult.value = result
  } catch (err) {
    testResult.value = {
      success: false,
      message: err.message || 'Lỗi kiểm tra kết nối'
    }
  } finally {
    testingConnection.value = false
  }
}

const handleSave = async () => {
  if (!isValid.value || saving.value) return

  saving.value = true
  error.value = null

  try {
    const printerData = { ...formData.value }
    
    // Clean up data based on connection type
    if (printerData.connection_type === 'NETWORK') {
      delete printerData.usb_path
    } else if (printerData.connection_type === 'USB') {
      delete printerData.ip_address
      delete printerData.port
    }

    // Clean up paper_width for LABEL type
    if (printerData.type === 'LABEL') {
      delete printerData.paper_width
    }

    await printerStore.savePrinter(printerData)
    emit('saved')
  } catch (err) {
    error.value = printerStore.error || 'Lỗi lưu cấu hình máy in'
  } finally {
    saving.value = false
  }
}
</script>
