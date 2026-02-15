<template>
  <div :class="cardClasses" class="rounded-xl p-4 border-l-4">
    <!-- Header -->
    <div class="flex items-start justify-between mb-2">
      <div class="flex items-center gap-2">
        <span class="text-2xl">{{ icon }}</span>
        <div>
          <h4 class="font-bold text-sm">{{ alert.batch_name }}</h4>
          <p v-if="subtitle" class="text-xs opacity-75">{{ subtitle }}</p>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="space-y-2 text-sm">
      <!-- Expired Alert -->
      <template v-if="type === 'expired'">
        <div class="flex justify-between">
          <span class="opacity-75">Số lượng lãng phí:</span>
          <span class="font-bold">{{ alert.quantity_wasted }} {{ alert.unit }}</span>
        </div>
        <div class="flex justify-between">
          <span class="opacity-75">Chi phí lãng phí:</span>
          <span class="font-bold text-red-600">{{ formatCurrency(alert.cost_wasted) }}</span>
        </div>
        <div class="flex justify-between">
          <span class="opacity-75">Hết hạn lúc:</span>
          <span class="font-medium">{{ formatDateTime(alert.expired_at) }}</span>
        </div>
      </template>

      <!-- Expiring Alert -->
      <template v-else-if="type === 'expiring'">
        <div class="flex justify-between">
          <span class="opacity-75">Số lượng còn lại:</span>
          <span class="font-bold">{{ alert.quantity_remaining }} {{ alert.unit }}</span>
        </div>
        <div class="flex justify-between">
          <span class="opacity-75">Hết hạn sau:</span>
          <span class="font-bold text-yellow-600">{{ alert.hours_until_expiry }}h</span>
        </div>
        <div class="flex justify-between">
          <span class="opacity-75">Hết hạn lúc:</span>
          <span class="font-medium">{{ formatDateTime(alert.expires_at) }}</span>
        </div>
      </template>

      <!-- Low Stock Alert -->
      <template v-else-if="type === 'low_stock'">
        <div class="flex justify-between">
          <span class="opacity-75">Tồn kho hiện tại:</span>
          <span class="font-bold">{{ alert.current_stock }} {{ alert.unit }}</span>
        </div>
        <div class="flex justify-between">
          <span class="opacity-75">Ngưỡng cảnh báo:</span>
          <span class="font-medium">{{ alert.threshold }} {{ alert.unit }}</span>
        </div>
        <div class="flex justify-between">
          <span class="opacity-75">Tỷ lệ:</span>
          <span :class="stockRatioClass" class="font-bold">
            {{ stockRatioPercent }}%
          </span>
        </div>
      </template>
    </div>

    <!-- Actions -->
    <div v-if="showActions" class="mt-3 pt-3 border-t flex gap-2">
      <button 
        v-if="type === 'low_stock'"
        @click="handleCreateBatch"
        class="flex-1 bg-blue-500 text-white py-2 rounded-lg text-xs font-bold active:bg-blue-600">
        Chế Biến Batch
      </button>
      <button 
        v-if="type === 'expiring'"
        @click="handleViewBatch"
        class="flex-1 bg-yellow-500 text-white py-2 rounded-lg text-xs font-bold active:bg-yellow-600">
        Xem Chi Tiết
      </button>
      <button 
        v-if="type === 'expired'"
        @click="handleMarkExpired"
        class="flex-1 bg-red-500 text-white py-2 rounded-lg text-xs font-bold active:bg-red-600">
        Xác Nhận Hết Hạn
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  alert: {
    type: Object,
    required: true
  },
  type: {
    type: String,
    required: true,
    validator: (value) => ['expired', 'expiring', 'low_stock'].includes(value)
  },
  showActions: {
    type: Boolean,
    default: true
  }
})

const router = useRouter()

const cardClasses = computed(() => {
  const baseClasses = 'transition-all'
  
  switch (props.type) {
    case 'expired':
      return `${baseClasses} bg-red-50 border-red-500 text-red-900`
    case 'expiring':
      return `${baseClasses} bg-yellow-50 border-yellow-500 text-yellow-900`
    case 'low_stock':
      return `${baseClasses} bg-orange-50 border-orange-500 text-orange-900`
    default:
      return baseClasses
  }
})

const icon = computed(() => {
  switch (props.type) {
    case 'expired':
      return '🚫'
    case 'expiring':
      return '⏰'
    case 'low_stock':
      return '📉'
    default:
      return '⚠️'
  }
})

const subtitle = computed(() => {
  switch (props.type) {
    case 'expired':
      return 'Không thể sử dụng'
    case 'expiring':
      return 'Cần sử dụng sớm'
    case 'low_stock':
      return 'Cần chế biến thêm'
    default:
      return ''
  }
})

const stockRatioPercent = computed(() => {
  if (props.type !== 'low_stock') return 0
  const ratio = (props.alert.current_stock / props.alert.threshold) * 100
  return Math.round(ratio)
})

const stockRatioClass = computed(() => {
  const percent = stockRatioPercent.value
  if (percent <= 25) return 'text-red-600'
  if (percent <= 50) return 'text-orange-600'
  if (percent <= 75) return 'text-yellow-600'
  return 'text-green-600'
})

const formatCurrency = (value) => {
  if (!value) return '0đ'
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND'
  }).format(value)
}

const formatDateTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  return date.toLocaleString('vi-VN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const handleCreateBatch = () => {
  // Navigate to batch record creation with pre-filled batch definition
  router.push({
    name: 'batch-record-create',
    query: { batch_definition_id: props.alert.batch_definition_id }
  })
}

const handleViewBatch = () => {
  // Navigate to batch record detail
  router.push({
    name: 'batch-record-detail',
    params: { id: props.alert.batch_record_id }
  })
}

const handleMarkExpired = () => {
  // Navigate to batch record detail for confirmation
  router.push({
    name: 'batch-record-detail',
    params: { id: props.alert.batch_record_id }
  })
}
</script>
