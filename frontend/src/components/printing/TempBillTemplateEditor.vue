<template>
  <div class="h-full flex flex-col bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm flex-shrink-0 px-4 py-3">
      <div class="flex items-center justify-between mb-3">
        <h2 class="text-lg font-bold text-gray-800">📄 Temp Bill Template</h2>
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
            class="bg-amber-500 text-white px-4 py-2 rounded-lg text-sm font-bold hover:bg-amber-600 active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
          >
            <span v-if="saving" class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
            <span>💾 Lưu Template</span>
          </button>
        </div>
      </div>
      
      <div class="text-sm text-gray-600">
        <p class="font-semibold text-amber-700">Template HTML cho bill tạm (chưa thanh toán). Có watermark và thông báo chưa thanh toán.</p>
        <p class="mt-1">Kích thước: 576px width (72mm @ 203 DPI), Margin: 20px</p>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-hidden flex">
      <!-- Left: HTML Editor -->
      <div class="flex-1 flex flex-col border-r">
        <div class="bg-amber-100 px-4 py-2 border-b">
          <h3 class="font-bold text-sm text-amber-800">📝 Temp Bill HTML Template</h3>
        </div>
        <div class="flex-1 overflow-y-auto p-4">
          <textarea
            v-model="htmlContent"
            class="w-full h-full font-mono text-sm border-2 border-amber-300 rounded-lg p-4 focus:ring-2 focus:ring-amber-500 focus:border-amber-500"
            placeholder="Nhập HTML template cho bill tạm..."
            @input="debouncedPreview"
          ></textarea>
        </div>
      </div>

      <!-- Right: Preview -->
      <div class="flex-1 flex flex-col bg-gray-100">
        <div class="bg-gray-200 px-4 py-2 border-b flex items-center justify-between">
          <h3 class="font-bold text-sm text-gray-700">👁️ Preview</h3>
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
            <div class="text-amber-700 font-semibold">Status: CHƯA THANH TOÁN</div>
          </div>
        </div>
        
        <div class="flex-1 overflow-y-auto p-4">
          <!-- Preview Frame -->
          <div class="bg-white rounded-lg shadow-lg p-4">
            <div class="border-2 border-amber-300 bg-white" style="width: 576px; margin: 0 auto;">
              <iframe
                ref="previewFrame"
                class="w-full border-0"
                :style="{ height: previewHeight + 'px' }"
                sandbox="allow-same-origin"
              ></iframe>
            </div>
            
            <div class="mt-4 text-xs text-gray-500 text-center">
              <p class="text-amber-700 font-semibold">Preview Bill Tạm với dữ liệu mẫu</p>
              <p>Kích thước: 576px width (72mm @ 203 DPI)</p>
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
  CreatedDate: '26/02/2026 12:45 PM',
  Items: [
    { STT: 1, Name: 'Cà phê sữa đá', Quantity: 2, UnitPrice: '25,000', Total: '50,000', Note: '50%' },
    { STT: 2, Name: 'Trà đào cam sả', Quantity: 1, UnitPrice: '35,000', Total: '35,000', Note: '' },
    { STT: 3, Name: 'Bánh mì', Quantity: 1, UnitPrice: '20,000', Total: '20,000', Note: '' }
  ],
  Total: '105,000'
})

let debounceTimer = null

const debouncedPreview = () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    updatePreview()
  }, 500)
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
  processed = processed.replace(/\{\{\.CreatedDate\}\}/g, sampleData.value.CreatedDate)
  processed = processed.replace(/\{\{\.CustomMessage\}\}/g, sampleData.value.CustomMessage)
  processed = processed.replace(/\{\{\.Total\}\}/g, sampleData.value.Total)
  processed = processed.replace(/\{\{\.LogoBase64\}\}/g, sampleData.value.LogoBase64)
  
  // Handle conditionals
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
    
    sampleData.value.Items.forEach((item) => {
      let itemHtml = itemTemplate
      itemHtml = itemHtml.replace(/\{\{\.STT\}\}/g, item.STT.toString())
      itemHtml = itemHtml.replace(/\{\{\.Name\}\}/g, item.Name)
      itemHtml = itemHtml.replace(/\{\{\.Quantity\}\}/g, item.Quantity.toString())
      itemHtml = itemHtml.replace(/\{\{\.UnitPrice\}\}/g, item.UnitPrice)
      itemHtml = itemHtml.replace(/\{\{\.Total\}\}/g, item.Total)
      
      // Handle note conditional
      if (item.Note) {
        itemHtml = itemHtml.replace(/\{\{if \.Note\}\}([\s\S]*?)\{\{end\}\}/g, (match, content) => {
          return content.replace(/\{\{\.Note\}\}/g, item.Note)
        })
      } else {
        itemHtml = itemHtml.replace(/\{\{if \.Note\}\}[\s\S]*?\{\{end\}\}/g, '')
      }
      
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
    const response = await api.get('/manager/html-templates/temp-bill')
    htmlContent.value = response.data.content
    
    // Load shop settings to get logo
    await loadShopSettings()
    
    updatePreview()
    statusMessage.value = {
      type: 'success',
      text: 'Temp bill template loaded successfully'
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
        const logoResponse = await fetch(settings.logo_url)
        if (!logoResponse.ok) {
          throw new Error(`Failed to fetch logo: ${logoResponse.status}`)
        }
        
        const logoBlob = await logoResponse.blob()
        const reader = new FileReader()
        
        reader.onloadend = () => {
          sampleData.value.LogoBase64 = reader.result
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
    await api.put('/manager/html-templates/temp-bill', {
      content: htmlContent.value
    })
    
    statusMessage.value = {
      type: 'success',
      text: 'Temp bill template saved successfully'
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

onMounted(() => {
  loadTemplate()
})
</script>
