<template>
  <div class="h-full flex flex-col bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm flex-shrink-0 px-4 py-3 border-b">
      <h2 class="text-lg font-bold text-gray-800">⚙️ Cài Đặt In Ấn</h2>
      <p class="text-sm text-gray-600 mt-1">Cấu hình thông tin quán và tùy chỉnh mẫu in</p>
    </div>

    <!-- Loading State -->
    <div v-if="loading && !settings" class="flex-1 flex items-center justify-center">
      <div class="text-center">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
        <p class="mt-4 text-gray-600">Đang tải cài đặt...</p>
      </div>
    </div>

    <!-- Error State - Show Create Form -->
    <div v-else-if="error && !settings" class="flex-1 overflow-y-auto">
      <div class="max-w-4xl mx-auto p-6">
        <!-- Info Banner -->
        <div class="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
          <div class="flex items-start gap-3">
            <span class="text-2xl">ℹ️</span>
            <div>
              <h3 class="font-semibold text-blue-900 mb-1">Chưa có cài đặt</h3>
              <p class="text-sm text-blue-700">
                Hệ thống chưa có cài đặt in ấn. Vui lòng điền thông tin bên dưới để tạo cài đặt mới.
              </p>
            </div>
          </div>
        </div>
        
        <!-- Create Form -->
        <form @submit.prevent="handleCreate" class="space-y-6">
          <div class="bg-white rounded-lg shadow p-6">
            <h3 class="text-lg font-semibold text-gray-800 mb-4">📋 Thông Tin Quán</h3>
            
            <div class="space-y-4">
              <!-- Shop Name -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Tên Quán <span class="text-red-500">*</span>
                </label>
                <input
                  v-model="formData.shop_name"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Nhập tên quán"
                />
              </div>

              <!-- Shop Address -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Địa Chỉ
                </label>
                <input
                  v-model="formData.shop_address"
                  type="text"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Nhập địa chỉ quán"
                />
              </div>

              <!-- Shop Phone -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Số Điện Thoại
                </label>
                <input
                  v-model="formData.shop_phone"
                  type="tel"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Nhập số điện thoại"
                />
              </div>

              <!-- Custom Message -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Lời Cảm Ơn / Thông Điệp
                </label>
                <textarea
                  v-model="formData.custom_message"
                  rows="3"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Cảm ơn quý khách! Hẹn gặp lại..."
                ></textarea>
              </div>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex justify-end gap-3 bg-white rounded-lg shadow p-4">
            <button
              type="submit"
              class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="loading"
            >
              <span v-if="loading">Đang tạo...</span>
              <span v-else>✨ Tạo cài đặt</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Form -->
    <div v-else-if="settings" class="flex-1 overflow-y-auto">
      <div class="max-w-4xl mx-auto p-6">
        <form @submit.prevent="handleSubmit" class="space-y-6">
          <!-- Shop Information Section -->
          <div class="bg-white rounded-lg shadow p-6">
            <h3 class="text-lg font-semibold text-gray-800 mb-4">📋 Thông Tin Quán</h3>
            
            <div class="space-y-4">
              <!-- Shop Name -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Tên Quán <span class="text-red-500">*</span>
                </label>
                <input
                  v-model="formData.shop_name"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Nhập tên quán"
                />
              </div>

              <!-- Shop Address -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Địa Chỉ
                </label>
                <input
                  v-model="formData.shop_address"
                  type="text"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Nhập địa chỉ quán"
                />
                <div class="mt-2">
                  <label class="inline-flex items-center">
                    <input
                      v-model="formData.show_address"
                      type="checkbox"
                      class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                    />
                    <span class="ml-2 text-sm text-gray-700">Hiển thị trên bill</span>
                  </label>
                </div>
              </div>

              <!-- Shop Phone -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Số Điện Thoại
                </label>
                <input
                  v-model="formData.shop_phone"
                  type="tel"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Nhập số điện thoại"
                />
                <div class="mt-2">
                  <label class="inline-flex items-center">
                    <input
                      v-model="formData.show_phone"
                      type="checkbox"
                      class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                    />
                    <span class="ml-2 text-sm text-gray-700">Hiển thị trên bill</span>
                  </label>
                </div>
              </div>

              <!-- Logo URL -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Logo URL
                </label>
                <input
                  v-model="formData.logo_url"
                  type="url"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="https://example.com/logo.png"
                />
                <p class="mt-1 text-xs text-gray-500">URL của logo quán (nếu có)</p>
                <div class="mt-2">
                  <label class="inline-flex items-center">
                    <input
                      v-model="formData.show_logo"
                      type="checkbox"
                      class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                    />
                    <span class="ml-2 text-sm text-gray-700">Hiển thị logo trên bill</span>
                  </label>
                </div>
              </div>

              <!-- Custom Message -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Lời Cảm Ơn / Thông Điệp
                </label>
                <textarea
                  v-model="formData.custom_message"
                  rows="3"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Cảm ơn quý khách! Hẹn gặp lại..."
                ></textarea>
                <div class="mt-2">
                  <label class="inline-flex items-center">
                    <input
                      v-model="formData.show_custom_message"
                      type="checkbox"
                      class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                    />
                    <span class="ml-2 text-sm text-gray-700">Hiển thị trên bill</span>
                  </label>
                </div>
              </div>
            </div>
          </div>

          <!-- Print Configuration Section -->
          <div class="bg-white rounded-lg shadow p-6">
            <h3 class="text-lg font-semibold text-gray-800 mb-4">🖨️ Cấu Hình In</h3>
            
            <div class="space-y-4">
              <!-- Print Bridge URL -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Print Bridge URL <span class="text-red-500">*</span>
                </label>
                <input
                  v-model="formData.print_bridge_url"
                  type="url"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-sm"
                  placeholder="http://192.168.1.100:3001"
                />
                <p class="mt-1 text-xs text-gray-500">
                  URL của Print Bridge chạy trên máy local (ví dụ: http://192.168.1.100:3001)
                </p>
                <div class="mt-2 flex items-center gap-2">
                  <button
                    type="button"
                    @click="testPrintBridge"
                    class="text-xs px-3 py-1 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded border border-gray-300"
                    :disabled="!formData.print_bridge_url || testingBridge"
                  >
                    <span v-if="testingBridge">⏳ Đang kiểm tra...</span>
                    <span v-else>🔍 Kiểm tra kết nối</span>
                  </button>
                  <span v-if="bridgeStatus === 'success'" class="text-xs text-green-600">
                    ✅ Kết nối thành công
                  </span>
                  <span v-else-if="bridgeStatus === 'error'" class="text-xs text-red-600">
                    ❌ Không thể kết nối
                  </span>
                </div>
              </div>

              <!-- Auto Print -->
              <div class="border-t pt-4">
                <label class="inline-flex items-center">
                  <input
                    v-model="formData.auto_print_enabled"
                    type="checkbox"
                    class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                  <span class="ml-2 text-sm font-medium text-gray-700">
                    Tự động in bill và tem khi tạo đơn hàng
                  </span>
                </label>
                <p class="mt-1 ml-6 text-xs text-gray-500">
                  Nếu tắt, bạn có thể in thủ công từ chi tiết đơn hàng
                </p>
              </div>

              <!-- Bill Paper Width -->
              <div class="border-t pt-4">
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Khổ Giấy Bill
                </label>
                <select
                  v-model.number="formData.paper_width"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  <option :value="40">40mm (Tem nhỏ)</option>
                  <option :value="58">58mm (Phổ biến)</option>
                  <option :value="80">80mm (Lớn)</option>
                </select>
                <p class="mt-1 text-xs text-gray-500">
                  Chọn khổ giấy phù hợp với máy in bill của bạn
                </p>
              </div>

              <!-- Test Print Button -->
              <div class="border-t pt-4">
                <button
                  type="button"
                  @click="testPrintTemplate"
                  class="w-full px-4 py-3 bg-green-50 hover:bg-green-100 text-green-700 rounded-lg border-2 border-green-200 font-medium transition-colors"
                  :disabled="!formData.print_bridge_url || testPrinting"
                >
                  <span v-if="testPrinting">⏳ Đang in thử...</span>
                  <span v-else>🖨️ In Thử Bill Mẫu</span>
                </button>
                <p class="mt-2 text-xs text-gray-500 text-center">
                  In một bill test để kiểm tra cấu hình
                </p>
              </div>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex justify-end gap-3 bg-white rounded-lg shadow p-4">
            <button
              type="button"
              @click="resetForm"
              class="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50"
              :disabled="loading"
            >
              Đặt lại
            </button>
            <button
              type="submit"
              class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="loading"
            >
              <span v-if="loading">Đang lưu...</span>
              <span v-else>💾 Lưu cài đặt</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Success Toast -->
    <div
      v-if="showSuccess"
      class="fixed bottom-4 right-4 bg-green-500 text-white px-6 py-3 rounded-lg shadow-lg flex items-center gap-2 animate-fade-in"
    >
      <span>✅</span>
      <span>Đã lưu cài đặt thành công!</span>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useShopSettingsStore } from '../../stores/shopSettings'

const store = useShopSettingsStore()

const loading = computed(() => store.loading)
const error = computed(() => store.error)
const settings = computed(() => store.settings)

const showSuccess = ref(false)
const testingBridge = ref(false)
const bridgeStatus = ref(null) // 'success', 'error', or null
const testPrinting = ref(false)
const testPrintResult = ref(null) // 'success', 'error', or null

const formData = reactive({
  shop_name: '',
  shop_address: '',
  shop_phone: '',
  logo_url: '',
  custom_message: '',
  print_bridge_url: '',
  show_logo: true,
  show_address: true,
  show_phone: true,
  show_custom_message: true,
  auto_print_enabled: true,
  paper_width: 58
})

onMounted(async () => {
  await loadSettings()
})

async function loadSettings() {
  try {
    await store.fetchSettings()
    if (settings.value) {
      Object.assign(formData, {
        shop_name: settings.value.shop_name || '',
        shop_address: settings.value.shop_address || '',
        shop_phone: settings.value.shop_phone || '',
        logo_url: settings.value.logo_url || '',
        custom_message: settings.value.custom_message || '',
        print_bridge_url: settings.value.print_bridge_url || '',
        show_logo: settings.value.show_logo !== false,
        show_address: settings.value.show_address !== false,
        show_phone: settings.value.show_phone !== false,
        show_custom_message: settings.value.show_custom_message !== false,
        auto_print_enabled: settings.value.auto_print_enabled !== false,
        paper_width: settings.value.paper_width || 58
      })
    }
  } catch (err) {
    console.error('Failed to load settings:', err)
  }
}

function resetForm() {
  if (settings.value) {
    Object.assign(formData, {
      shop_name: settings.value.shop_name || '',
      shop_address: settings.value.shop_address || '',
      shop_phone: settings.value.shop_phone || '',
      logo_url: settings.value.logo_url || '',
      custom_message: settings.value.custom_message || '',
      print_bridge_url: settings.value.print_bridge_url || '',
      show_logo: settings.value.show_logo !== false,
      show_address: settings.value.show_address !== false,
      show_phone: settings.value.show_phone !== false,
      show_custom_message: settings.value.show_custom_message !== false,
      auto_print_enabled: settings.value.auto_print_enabled !== false,
      paper_width: settings.value.paper_width || 58
    })
  }
  bridgeStatus.value = null
}

async function testPrintBridge() {
  if (!formData.print_bridge_url) return
  
  testingBridge.value = true
  bridgeStatus.value = null
  
  try {
    // Call backend API instead of direct fetch to avoid CORS
    const token = localStorage.getItem('token')
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    
    const response = await fetch(`${apiUrl}/api/manager/print-bridge/test`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        bridge_url: formData.print_bridge_url
      })
    })
    
    const data = await response.json()
    
    if (response.ok && data.success) {
      bridgeStatus.value = 'success'
    } else {
      bridgeStatus.value = 'error'
    }
  } catch (err) {
    console.error('Print bridge test failed:', err)
    bridgeStatus.value = 'error'
  } finally {
    testingBridge.value = false
  }
}

async function testPrintTemplate() {
  const printerIP = prompt('Nhập IP máy in (ví dụ: 192.168.1.100):')
  if (!printerIP) return
  
  testPrinting.value = true
  testPrintResult.value = null
  
  try {
    const token = localStorage.getItem('token')
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    
    const response = await fetch(`${apiUrl}/api/manager/html-templates/test-print`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        use_test_data: true,
        printer_ip: printerIP
      })
    })
    
    const data = await response.json()
    
    if (response.ok && data.success) {
      testPrintResult.value = 'success'
      alert('✅ In thử thành công!\n\nĐơn hàng test: ' + data.order_number)
    } else {
      testPrintResult.value = 'error'
      alert('❌ Lỗi: ' + (data.error || 'Không thể in'))
    }
  } catch (err) {
    console.error('Test print failed:', err)
    testPrintResult.value = 'error'
    alert('❌ Lỗi: ' + err.message)
  } finally {
    testPrinting.value = false
  }
}

async function handleSubmit() {
  try {
    await store.updateSettings(formData)
    
    // Show success message
    showSuccess.value = true
    setTimeout(() => {
      showSuccess.value = false
    }, 3000)
  } catch (err) {
    console.error('Failed to update settings:', err)
    alert('Lỗi: ' + (err.response?.data?.error || 'Không thể lưu cài đặt'))
  }
}

async function handleCreate() {
  try {
    await store.createSettings(formData)
    
    // Reload settings
    await loadSettings()
    
    // Show success message
    showSuccess.value = true
    setTimeout(() => {
      showSuccess.value = false
    }, 3000)
  } catch (err) {
    console.error('Failed to create settings:', err)
    alert('Lỗi: ' + (err.response?.data?.error || 'Không thể tạo cài đặt'))
  }
}
</script>

<style scoped>
@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in {
  animation: fade-in 0.3s ease-out;
}
</style>
