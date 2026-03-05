<template>
  <div class="h-full flex flex-col bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm flex-shrink-0 px-4 py-3">
      <div class="flex items-center justify-between mb-3">
        <h2 class="text-lg font-bold text-gray-800">📝 HTML Bill Template</h2>
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
        <p>Template HTML cho bill in. Chỉnh sửa HTML/CSS để customize layout.</p>
        <p class="mt-1">Kích thước: 576px width (72mm @ 203 DPI), Margin: 20px</p>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-hidden flex">
      <!-- Left: HTML Editor -->
      <div class="flex-1 flex flex-col border-r">
        <div class="bg-gray-100 px-4 py-2 border-b">
          <h3 class="font-bold text-sm text-gray-700">📝 HTML Template</h3>
        </div>
        <div class="flex-1 overflow-y-auto p-4">
          <textarea
            v-model="htmlContent"
            class="w-full h-full font-mono text-sm border-2 border-gray-300 rounded-lg p-4 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            placeholder="Nhập HTML template..."
            @input="debouncedPreview"
          ></textarea>
        </div>
      </div>

      <!-- Right: Preview & Test -->
      <div class="flex-1 flex flex-col bg-gray-100">
        <div class="bg-gray-200 px-4 py-2 border-b flex items-center justify-between">
          <h3 class="font-bold text-sm text-gray-700">👁️ Preview & Test</h3>
          <div class="flex gap-2">
            <button
              @click="refreshPreview"
              class="text-xs bg-white px-3 py-1 rounded border hover:bg-gray-50"
            >
              🔄 Refresh
            </button>
            <button
              @click="toggleSampleData"
              class="text-xs bg-white px-3 py-1 rounded border hover:bg-gray-50"
            >
              {{ showSampleData ? '📋 Hide Data' : '📋 Show Data' }}
            </button>
          </div>
        </div>
        
        <!-- Sample Data Panel -->
        <div v-if="showSampleData" class="bg-yellow-50 border-b p-3 text-xs">
          <div class="font-bold mb-2">Dữ liệu mẫu:</div>
          <div class="space-y-1 text-gray-700">
            <div>ShopName: {{ sampleData.ShopName }}</div>
            <div>OrderNumber: {{ sampleData.OrderNumber }}</div>
            <div>Total: {{ sampleData.Total }}</div>
            <div>Items: {{ sampleData.Items.length }} món</div>
          </div>
        </div>
        
        <div class="flex-1 overflow-y-auto p-4">
          <!-- Preview Frame -->
          <div class="bg-white rounded-lg shadow-lg p-4 mb-4">
            <div class="border-2 border-gray-300 bg-white" style="width: 576px; margin: 0 auto;">
              <iframe
                ref="previewFrame"
                class="w-full border-0"
                :style="{ height: previewHeight + 'px' }"
                sandbox="allow-same-origin"
              ></iframe>
            </div>
            
            <div class="mt-4 text-xs text-gray-500 text-center">
              <p>Preview với dữ liệu mẫu</p>
              <p>Kích thước: 576px width (72mm @ 203 DPI)</p>
            </div>
          </div>
          
          <!-- Test Print Section -->
          <div class="bg-white rounded-lg shadow-lg p-4">
            <h3 class="font-bold text-gray-800 mb-3">🖨️ Test Print với Order thật</h3>
            
            <div class="space-y-3">
              <!-- Printer IP -->
              <div>
                <label class="block text-xs font-bold text-gray-700 mb-1">IP Máy in</label>
                <input
                  v-model="printerIP"
                  type="text"
                  placeholder="192.168.1.115"
                  class="w-full px-3 py-2 text-sm border-2 border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
              
              <!-- Order Selection -->
              <div>
                <label class="block text-xs font-bold text-gray-700 mb-1">Chọn Order</label>
                <div class="flex gap-2">
                  <input
                    v-model="searchQuery"
                    type="text"
                    placeholder="Tìm order..."
                    class="flex-1 px-3 py-2 text-sm border-2 border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    @input="searchOrders"
                  />
                  <button
                    @click="searchOrders"
                    class="bg-gray-500 text-white px-3 py-2 rounded-lg text-xs font-bold hover:bg-gray-600"
                  >
                    🔍
                  </button>
                </div>
                
                <!-- Orders dropdown -->
                <div v-if="orders.length > 0" class="mt-2 max-h-40 overflow-y-auto border rounded-lg">
                  <div
                    v-for="order in orders"
                    :key="order.id"
                    @click="selectedOrder = order"
                    :class="[
                      'p-2 cursor-pointer text-xs hover:bg-gray-50',
                      selectedOrder?.id === order.id ? 'bg-blue-50 border-l-4 border-blue-500' : ''
                    ]"
                  >
                    <div class="font-bold">{{ order.order_number }}</div>
                    <div class="text-gray-600">{{ formatMoney(order.total) }} VNĐ</div>
                  </div>
                </div>
                
                <!-- Selected order info -->
                <div v-if="selectedOrder" class="mt-2 p-2 bg-blue-50 rounded text-xs">
                  <div class="font-bold">✓ {{ selectedOrder.order_number }}</div>
                  <div class="text-gray-600">{{ formatMoney(selectedOrder.total) }} VNĐ - {{ selectedOrder.items?.length || 0 }} items</div>
                </div>
              </div>
              
              <!-- Action buttons -->
              <div class="flex gap-2">
                <button
                  @click="testPrint"
                  :disabled="!selectedOrder || printing"
                  class="flex-1 bg-green-500 text-white px-3 py-2 rounded-lg text-xs font-bold hover:bg-green-600 active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-1"
                >
                  <span v-if="printing" class="inline-block w-3 h-3 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
                  <span v-else>🖨️</span>
                  <span>{{ printing ? 'Đang in...' : 'Test Print' }}</span>
                </button>
                
                <button
                  @click="previewWithOrder"
                  :disabled="!selectedOrder || previewing"
                  class="flex-1 bg-purple-500 text-white px-3 py-2 rounded-lg text-xs font-bold hover:bg-purple-600 active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-1"
                >
                  <span v-if="previewing" class="inline-block w-3 h-3 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
                  <span v-else>👁️</span>
                  <span>{{ previewing ? 'Đang tạo...' : 'Preview PNG' }}</span>
                </button>
              </div>
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
import { ref, onMounted, nextTick } from 'vue'
import api from '../../services/api'

const htmlContent = ref('')
const previewFrame = ref(null)
const previewHeight = ref(800)
const saving = ref(false)
const showSampleData = ref(false)
const statusMessage = ref(null)
const printerIP = ref('192.168.1.115')
const searchQuery = ref('')
const orders = ref([])
const selectedOrder = ref(null)
const printing = ref(false)
const previewing = ref(false)

const sampleData = ref({
  ShopName: 'Tiệm cà phê Ông Tạ',
  ShopAddress: 'Đ/c: 10/8 Trần Nhật Duật, P. Tân Định, Quận 1, HCM',
  ShopPhone: '0906990602',
  ShowLogo: false,
  ShowAddress: true,
  ShowPhone: true,
  ShowCustomMessage: true,
  CustomMessage: 'Cảm ơn quý khách!',
  LogoBase64: '',
  OrderNumber: '20260222-095703-168',
  WaiterName: 'Waiter1',
  PaymentMethod: 'Tiền mặt',
  CreatedDate: '26/02/2026 12:45 PM',
  Items: [
    { Name: 'Cà phê sữa đá', VariantName: '', Quantity: 2, Price: 25000, Subtotal: 50000 },
    { Name: 'Trà đào cam sả', VariantName: '', Quantity: 1, Price: 35000, Subtotal: 35000 },
    { Name: 'Bánh mì', VariantName: '', Quantity: 1, Price: 20000, Subtotal: 20000 }
  ],
  Total: 105000
})

let debounceTimer = null

const debouncedPreview = () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    updatePreview()
  }, 500)
}

const formatMoney = (amount) => {
  if (!amount) return '0'
  return amount.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

const formatMoneyTemplate = (amount) => {
  const amountInt = Math.floor(amount)
  if (amountInt >= 1000000) {
    const millions = Math.floor(amountInt / 1000000)
    const thousands = Math.floor((amountInt % 1000000) / 1000)
    const ones = amountInt % 1000
    return `${millions},${thousands.toString().padStart(3, '0')},${ones.toString().padStart(3, '0')}`
  } else if (amountInt >= 1000) {
    const thousands = Math.floor(amountInt / 1000)
    const ones = amountInt % 1000
    return `${thousands},${ones.toString().padStart(3, '0')}`
  }
  return amountInt.toString()
}

const processTemplate = (html) => {
  // Replace Go template variables with actual data
  let processed = html
  
  // Simple replacements
  processed = processed.replace(/\{\{\.ShopName\}\}/g, sampleData.value.ShopName)
  processed = processed.replace(/\{\{\.ShopAddress\}\}/g, sampleData.value.ShopAddress)
  processed = processed.replace(/\{\{\.ShopPhone\}\}/g, sampleData.value.ShopPhone)
  processed = processed.replace(/\{\{\.OrderNumber\}\}/g, sampleData.value.OrderNumber)
  processed = processed.replace(/\{\{\.WaiterName\}\}/g, sampleData.value.WaiterName)
  processed = processed.replace(/\{\{\.PaymentMethod\}\}/g, sampleData.value.PaymentMethod)
  processed = processed.replace(/\{\{\.CreatedDate\}\}/g, sampleData.value.CreatedDate)
  processed = processed.replace(/\{\{\.CustomMessage\}\}/g, sampleData.value.CustomMessage)
  processed = processed.replace(/\{\{\.Total\}\}/g, formatMoneyTemplate(sampleData.value.Total))
  processed = processed.replace(/\{\{\.LogoBase64\}\}/g, sampleData.value.LogoBase64)
  
  // Handle conditionals - show logo if we have it
  if (sampleData.value.ShowLogo && sampleData.value.LogoBase64) {
    processed = processed.replace(/\{\{if \.ShowLogo\}\}([\s\S]*?)\{\{end\}\}/g, '$1')
  } else {
    processed = processed.replace(/\{\{if \.ShowLogo\}\}[\s\S]*?\{\{end\}\}/g, '')
  }
  
  processed = processed.replace(/\{\{if \.ShowAddress\}\}([\s\S]*?)\{\{end\}\}/g, '$1')
  processed = processed.replace(/\{\{if \.ShowPhone\}\}([\s\S]*?)\{\{end\}\}/g, '$1')
  processed = processed.replace(/\{\{if \.ShowCustomMsg\}\}([\s\S]*?)\{\{end\}\}/g, '$1')
  
  // Handle range loop for items
  const rangeMatch = processed.match(/\{\{range \.Items\}\}([\s\S]*?)\{\{end\}\}/g)
  if (rangeMatch) {
    const itemTemplate = rangeMatch[0].replace(/\{\{range \.Items\}\}/, '').replace(/\{\{end\}\}/, '')
    let itemsHtml = ''
    
    sampleData.value.Items.forEach((item, index) => {
      let itemHtml = itemTemplate
      itemHtml = itemHtml.replace(/\{\{\.STT\}\}/g, (index + 1).toString())
      itemHtml = itemHtml.replace(/\{\{\.Name\}\}/g, item.Name + (item.VariantName ? ` (${item.VariantName})` : ''))
      itemHtml = itemHtml.replace(/\{\{\.Quantity\}\}/g, item.Quantity.toString())
      itemHtml = itemHtml.replace(/\{\{\.UnitPrice\}\}/g, formatMoneyTemplate(item.Price))
      itemHtml = itemHtml.replace(/\{\{\.Total\}\}/g, formatMoneyTemplate(item.Subtotal))
      itemsHtml += itemHtml
    })
    
    processed = processed.replace(/\{\{range \.Items\}\}[\s\S]*?\{\{end\}\}/, itemsHtml)
  }
  
  return processed
}

const updatePreview = async () => {
  if (!previewFrame.value) return
  
  try {
    const processedHtml = processTemplate(htmlContent.value)
    const iframe = previewFrame.value
    const doc = iframe.contentDocument || iframe.contentWindow.document
    
    doc.open()
    doc.write(processedHtml)
    doc.close()
    
    // Update iframe height after content loads
    await nextTick()
    setTimeout(() => {
      const contentHeight = doc.body.scrollHeight
      previewHeight.value = Math.max(contentHeight, 800)
    }, 100)
  } catch (error) {
    console.error('Preview error:', error)
  }
}

const refreshPreview = () => {
  updatePreview()
}

const toggleSampleData = () => {
  showSampleData.value = !showSampleData.value
}

const loadTemplate = async () => {
  try {
    const response = await api.get('/manager/html-templates/bill')
    htmlContent.value = response.data.content
    
    // Load shop settings to get logo
    await loadShopSettings()
    
    updatePreview()
    statusMessage.value = {
      type: 'success',
      text: 'Template loaded successfully'
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

const loadShopSettings = async () => {
  try {
    const response = await api.get('/manager/shop-settings')
    const settings = response.data
    
    // Update sample data with real shop settings
    sampleData.value.ShopName = settings.shop_name || sampleData.value.ShopName
    sampleData.value.ShopAddress = settings.shop_address || sampleData.value.ShopAddress
    sampleData.value.ShopPhone = settings.shop_phone || sampleData.value.ShopPhone
    sampleData.value.ShowLogo = settings.show_logo || false
    sampleData.value.ShowAddress = settings.show_address !== false
    sampleData.value.ShowPhone = settings.show_phone !== false
    sampleData.value.CustomMessage = settings.custom_message || sampleData.value.CustomMessage
    
    // Load logo as base64 if available
    if (settings.show_logo && settings.logo_url) {
      try {
        // Fetch logo from backend (will be proxied by Vite)
        const logoResponse = await fetch(settings.logo_url)
        if (!logoResponse.ok) {
          throw new Error(`Failed to fetch logo: ${logoResponse.status}`)
        }
        
        const logoBlob = await logoResponse.blob()
        const reader = new FileReader()
        
        reader.onloadend = () => {
          sampleData.value.LogoBase64 = reader.result
          console.log('Logo loaded successfully, base64 length:', reader.result.length)
          updatePreview()
        }
        
        reader.readAsDataURL(logoBlob)
      } catch (logoError) {
        console.error('Failed to load logo:', logoError)
      }
    }
  } catch (error) {
    console.error('Failed to load shop settings:', error)
  }
}

const saveTemplate = async () => {
  saving.value = true
  statusMessage.value = null
  
  try {
    await api.put('/manager/html-templates/bill', {
      content: htmlContent.value
    })
    
    statusMessage.value = {
      type: 'success',
      text: 'Template saved successfully'
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

const searchOrders = async () => {
  try {
    const params = {}
    if (searchQuery.value) {
      params.search = searchQuery.value
    }
    
    const response = await api.get('/manager/orders', { params })
    orders.value = response.data.orders || response.data || []
  } catch (error) {
    console.error('Failed to search orders:', error)
  }
}

const testPrint = async () => {
  if (!selectedOrder.value) return
  if (!printerIP.value) {
    alert('Vui lòng nhập IP máy in')
    return
  }
  
  printing.value = true
  statusMessage.value = null
  
  try {
    const response = await api.post('/manager/html-templates/test-print', {
      order_id: selectedOrder.value.id,
      printer_ip: printerIP.value
    })
    
    statusMessage.value = {
      type: 'success',
      text: `Test print thành công! Order: ${response.data.order_number}`
    }
  } catch (error) {
    console.error('Failed to print:', error)
    statusMessage.value = {
      type: 'error',
      text: error.response?.data?.error || 'Lỗi khi test print'
    }
  } finally {
    printing.value = false
  }
}

const previewWithOrder = async () => {
  if (!selectedOrder.value) return
  
  previewing.value = true
  statusMessage.value = null
  
  try {
    const response = await api.post('/manager/html-templates/preview', {
      order_id: selectedOrder.value.id
    })
    
    statusMessage.value = {
      type: 'success',
      text: `Preview đã tạo: ${response.data.filename}`
    }
  } catch (error) {
    console.error('Failed to preview:', error)
    statusMessage.value = {
      type: 'error',
      text: error.response?.data?.error || 'Lỗi khi tạo preview'
    }
  } finally {
    previewing.value = false
  }
}

onMounted(() => {
  loadTemplate()
  searchOrders()
})
</script>
