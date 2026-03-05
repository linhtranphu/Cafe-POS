<template>
  <div class="h-full flex flex-col bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm flex-shrink-0 px-4 py-3">
      <div class="flex items-center justify-between mb-3">
        <h2 class="text-lg font-bold text-gray-800">🧾 Visual Bill Template</h2>
        <div class="flex gap-2">
          <button
            @click="printBill"
            class="bg-green-500 text-white px-4 py-2 rounded-lg text-sm font-bold hover:bg-green-600 active:scale-95 transition-transform"
          >
            🖨️ In Bill
          </button>
          <button
            @click="resetToDefault"
            class="bg-gray-500 text-white px-4 py-2 rounded-lg text-sm font-bold hover:bg-gray-600 active:scale-95 transition-transform"
          >
            🔄 Reset
          </button>
          <button
            @click="saveTemplate"
            class="bg-blue-500 text-white px-4 py-2 rounded-lg text-sm font-bold hover:bg-blue-600 active:scale-95 transition-transform"
          >
            💾 Lưu Template
          </button>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-hidden flex">
      <!-- Left: Settings -->
      <div class="w-80 bg-white border-r overflow-y-auto p-4 space-y-4">
        <h3 class="font-bold text-gray-700 mb-3">⚙️ Cài đặt Template</h3>
        
        <div>
          <label class="block text-sm font-bold text-gray-700 mb-2">Tên Template</label>
          <input
            v-model="templateData.title"
            type="text"
            class="w-full px-3 py-2 border rounded-lg text-sm"
            placeholder="Tiệm cà phê Ông Tạ"
          />
        </div>

        <div>
          <label class="block text-sm font-bold text-gray-700 mb-2">Địa chỉ</label>
          <textarea
            v-model="templateData.address"
            rows="3"
            class="w-full px-3 py-2 border rounded-lg text-sm"
            placeholder="Đ/c: 10/8 Trần Nhật Duật..."
          ></textarea>
        </div>

        <div>
          <label class="block text-sm font-bold text-gray-700 mb-2">Số điện thoại</label>
          <input
            v-model="templateData.sdt"
            type="text"
            class="w-full px-3 py-2 border rounded-lg text-sm"
            placeholder="Hotline: 0906990602"
          />
        </div>

        <div>
          <label class="block text-sm font-bold text-gray-700 mb-2">Lời cảm ơn</label>
          <input
            v-model="templateData.thanks"
            type="text"
            class="w-full px-3 py-2 border rounded-lg text-sm"
            placeholder="Cảm ơn quý khách!"
          />
        </div>

        <div class="pt-4 border-t">
          <h4 class="font-bold text-sm text-gray-700 mb-2">📝 Dữ liệu mẫu</h4>
          <p class="text-xs text-gray-500 mb-2">Để xem preview</p>
        </div>

        <div>
          <label class="block text-sm font-bold text-gray-700 mb-2">Order Number</label>
          <input
            v-model="sampleData.orderNo"
            type="text"
            class="w-full px-3 py-2 border rounded-lg text-sm"
          />
        </div>

        <div>
          <label class="block text-sm font-bold text-gray-700 mb-2">Waiter</label>
          <input
            v-model="sampleData.waiter"
            type="text"
            class="w-full px-3 py-2 border rounded-lg text-sm"
          />
        </div>

        <div>
          <label class="block text-sm font-bold text-gray-700 mb-2">Payment Method</label>
          <input
            v-model="sampleData.paymentMethod"
            type="text"
            class="w-full px-3 py-2 border rounded-lg text-sm"
          />
        </div>
      </div>

      <!-- Right: Preview -->
      <div class="flex-1 overflow-y-auto p-8 bg-gray-100">
        <div class="max-w-2xl mx-auto">
          <div class="bg-white rounded-lg shadow-lg p-4">
            <h3 class="font-bold text-gray-700 mb-4">👁️ Preview (576px width)</h3>
            
            <!-- Bill Preview Canvas -->
            <div class="border-2 border-gray-300 bg-white overflow-hidden" style="width: 576px;">
              <canvas
                ref="billCanvas"
                :width="576"
                :height="canvasHeight"
                class="w-full"
              ></canvas>
            </div>

            <div class="mt-4 text-xs text-gray-500">
              <p>Kích thước: 576px × {{ canvasHeight }}px (72mm width @ 203 DPI)</p>
              <p>Margin: 20px</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick } from 'vue'

const billCanvas = ref(null)
const canvasHeight = ref(900)

const templateData = ref({
  title: 'Tiệm cà phê Ông Tạ',
  address: 'Đ/c: 10/8 Trần Nhật Duật, P. Tân Định, Quận 1, HCM',
  sdt: 'Hotline: 0906990602',
  thanks: 'Cảm ơn quý khách!'
})

const sampleData = ref({
  orderNo: 'Order: 20260222-095703-168',
  waiter: 'Waiter: Waiter1',
  paymentMethod: 'Thanh Toán: Tiền mặt',
  createdDate: 'Ngày tạo: 26/02/2026 12:45 AM',
  items: [
    { stt: 1, name: 'Cà phê sữa đá', quantity: 2, unitPrice: 25000, total: 50000 },
    { stt: 2, name: 'Trà đào cam sả', quantity: 1, unitPrice: 35000, total: 35000 },
    { stt: 3, name: 'Bánh mì', quantity: 1, unitPrice: 20000, total: 20000 }
  ]
})

const ImageWidthPixels = 576
const Margin = 20

const formatMoney = (amount) => {
  if (amount >= 1000000) {
    const millions = Math.floor(amount / 1000000)
    const thousands = Math.floor((amount % 1000000) / 1000)
    const ones = amount % 1000
    return `${millions},${thousands.toString().padStart(3, '0')},${ones.toString().padStart(3, '0')}`
  } else if (amount >= 1000) {
    const thousands = Math.floor(amount / 1000)
    const ones = amount % 1000
    return `${thousands},${ones.toString().padStart(3, '0')}`
  }
  return amount.toString()
}

const wrapText = (ctx, text, maxWidth) => {
  const words = text.split(' ')
  const lines = []
  let currentLine = ''
  
  for (const word of words) {
    const testLine = currentLine ? `${currentLine} ${word}` : word
    const metrics = ctx.measureText(testLine)
    
    if (metrics.width > maxWidth && currentLine) {
      lines.push(currentLine)
      currentLine = word
    } else {
      currentLine = testLine
    }
  }
  
  if (currentLine) {
    lines.push(currentLine)
  }
  
  return lines
}

const renderBill = async () => {
  await nextTick()
  
  if (!billCanvas.value) return
  
  const canvas = billCanvas.value
  const ctx = canvas.getContext('2d')
  
  // Clear canvas
  ctx.fillStyle = '#FFFFFF'
  ctx.fillRect(0, 0, ImageWidthPixels, canvas.height)
  ctx.fillStyle = '#000000'
  
  let y = 20
  
  // Logo placeholder (200px width)
  ctx.strokeStyle = '#CCCCCC'
  ctx.strokeRect(Margin + 20, y, 200, 100)
  ctx.font = '12px Arial'
  ctx.fillStyle = '#999999'
  ctx.fillText('Logo (200px)', Margin + 90, y + 50)
  ctx.fillStyle = '#000000'
  
  const logoHeight = 100
  
  // Title (right side of logo)
  const textX = Margin + 280
  ctx.font = 'bold 25px Arial'
  const maxWidthTitle = ImageWidthPixels - Margin - 210
  const titleLines = wrapText(ctx, templateData.value.title, maxWidthTitle)
  let textY = y + 20
  
  titleLines.forEach((line, i) => {
    ctx.fillText(line, textX, textY + i * 18)
    // Fake bold
    ctx.fillText(line, textX + 2, textY + i * 18)
  })
  
  // Address (below title)
  ctx.font = '16px Arial'
  const maxWidthAddress = ImageWidthPixels - Margin - 360
  const addressLines = wrapText(ctx, templateData.value.address, maxWidthAddress)
  textY = y + 50
  
  addressLines.forEach((line, i) => {
    ctx.fillText(line, Margin + 285, textY + i * 18)
  })
  
  // Phone (below address)
  ctx.font = '16px Arial'
  const sdtLines = wrapText(ctx, templateData.value.sdt, maxWidthAddress)
  textY = y + 87
  
  sdtLines.forEach((line, i) => {
    ctx.fillText(line, Margin + 285, textY + i * 18)
  })
  
  // Update y after logo
  y += Math.max(logoHeight, 60) + 45
  
  // Title "HÓA ĐƠN THANH TOÁN"
  ctx.font = 'bold 34px Arial'
  const billTitle = 'HÓA ĐƠN THANH TOÁN'
  const titleWidth = ctx.measureText(billTitle).width
  ctx.fillText(billTitle, (ImageWidthPixels - titleWidth) / 2, y)
  y += 40
  
  // Order info
  ctx.font = '16px Arial'
  ctx.fillText(sampleData.value.orderNo, Margin + 10, y)
  y += 20
  
  ctx.fillText(sampleData.value.waiter, Margin + 10, y)
  y += 20
  
  ctx.fillText(sampleData.value.paymentMethod, Margin + 10, y)
  y += 20
  
  ctx.fillText(sampleData.value.createdDate, Margin + 10, y)
  y += 20
  
  // Line
  ctx.beginPath()
  ctx.moveTo(Margin, y)
  ctx.lineTo(ImageWidthPixels - Margin, y)
  ctx.stroke()
  y += 25
  
  // Table header
  ctx.font = '17px Arial'
  const colX = [
    Margin + 10,
    Margin + 50,
    Margin + 290,
    Margin + 340,
    Margin + 450
  ]
  
  ctx.fillText('STT', colX[0], y)
  ctx.fillText('Tên món', colX[1], y)
  ctx.fillText('SL', colX[2], y)
  ctx.fillText('Đơn giá', colX[3], y)
  ctx.fillText('Thành tiền', colX[4], y)
  y += 8
  
  // Line
  ctx.beginPath()
  ctx.moveTo(Margin, y)
  ctx.lineTo(ImageWidthPixels - Margin, y)
  ctx.stroke()
  y += 20
  
  // Items
  let totalAmount = 0
  sampleData.value.items.forEach(item => {
    ctx.fillText(item.stt.toString(), colX[0], y)
    ctx.fillText(item.name, colX[1], y)
    ctx.fillText(item.quantity.toString(), colX[2], y)
    ctx.fillText(formatMoney(item.unitPrice), colX[3], y)
    
    const totalStr = formatMoney(item.total)
    const totalWidth = ctx.measureText(totalStr).width
    ctx.fillText(totalStr, ImageWidthPixels - Margin - totalWidth - 10, y)
    
    totalAmount += item.total
    y += 28
  })
  
  y += 5
  
  // Line
  ctx.beginPath()
  ctx.moveTo(Margin, y)
  ctx.lineTo(ImageWidthPixels - Margin, y)
  ctx.stroke()
  y += 30
  
  // Total
  ctx.font = 'bold 24px Arial'
  ctx.fillText('TỔNG TIỀN:', colX[2] - 50, y)
  const totalStr = formatMoney(totalAmount)
  const totalWidth = ctx.measureText(totalStr).width
  ctx.fillText(totalStr, ImageWidthPixels - Margin - totalWidth - 10, y)
  y += 30
  
  // Line
  ctx.beginPath()
  ctx.moveTo(Margin, y)
  ctx.lineTo(ImageWidthPixels - Margin, y)
  ctx.stroke()
  y += 30
  
  // Thanks
  ctx.font = '22px Arial'
  const thanksWidth = ctx.measureText(templateData.value.thanks).width
  ctx.fillText(templateData.value.thanks, (ImageWidthPixels - thanksWidth) / 2, y)
  y += 30
  
  // Update canvas height
  canvasHeight.value = y + 10
}

const resetToDefault = () => {
  templateData.value = {
    title: 'Tiệm cà phê Ông Tạ',
    address: 'Đ/c: 10/8 Trần Nhật Duật, P. Tân Định, Quận 1, HCM',
    sdt: 'Hotline: 0906990602',
    thanks: 'Cảm ơn quý khách!'
  }
  renderBill()
}

const saveTemplate = () => {
  // TODO: Implement save to backend
  alert('Chức năng lưu template sẽ được implement sau')
}

const printBill = async () => {
  if (!sampleData.value.orderNo) {
    alert('Vui lòng nhập Order Number')
    return
  }
  
  const printerIP = prompt('Nhập IP máy in (VD: 192.168.1.115):', '192.168.1.115')
  if (!printerIP) return
  
  try {
    // Extract order ID from order number or use a test order
    // For now, we'll need to get the actual order ID
    const orderIdInput = prompt('Nhập Order ID (ObjectID):', '')
    if (!orderIdInput) return
    
    const response = await fetch('http://localhost:8080/api/manager/visual-print/bill', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        order_id: orderIdInput,
        printer_ip: printerIP
      })
    })
    
    const result = await response.json()
    
    if (response.ok) {
      alert(`✅ In thành công! Order: ${result.order_number}`)
    } else {
      alert(`❌ Lỗi: ${result.error}`)
    }
  } catch (error) {
    alert(`❌ Lỗi kết nối: ${error.message}`)
  }
}

// Watch for changes and re-render
watch([templateData, sampleData], () => {
  renderBill()
}, { deep: true })

onMounted(() => {
  renderBill()
})
</script>
