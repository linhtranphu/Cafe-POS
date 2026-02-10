<template>
  <div>
    <!-- Date Range Display -->
    <div class="bg-blue-50 rounded-xl p-4 mb-4">
      <div class="text-xs md:text-sm text-blue-600 font-medium mb-1">Khoảng thời gian</div>
      <div class="text-sm md:text-base text-gray-800">
        {{ formatDate(dateRange.start) }} → {{ formatDate(dateRange.end) }}
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="categoryProfits.length === 0" class="text-center py-16">
      <div class="text-6xl mb-4">📭</div>
      <p class="text-gray-700 font-medium mb-2">Không có dữ liệu</p>
      <p class="text-gray-500 text-sm">Không có đơn hàng nào trong khoảng thời gian này</p>
    </div>

    <!-- Category Profit List - Mobile/Tablet -->
    <div v-else-if="categoryProfits.length > 0" class="space-y-3 md:hidden">
      <div 
        v-for="category in categoryProfits" 
        :key="category.category"
        class="bg-white rounded-2xl p-4 shadow-sm">
        
        <!-- Category Header -->
        <div class="flex items-center justify-between mb-3 pb-3 border-b">
          <div>
            <h3 class="font-bold text-lg text-gray-800">{{ category.category }}</h3>
            <div class="text-xs text-gray-500 mt-1">
              {{ category.order_count }} đơn • {{ category.item_count }} món
            </div>
          </div>
          <div class="text-right">
            <div class="text-xs text-gray-500 mb-1">Lợi nhuận %</div>
            <div class="text-lg font-bold" :class="getProfitMarginColor(category.average_profit_margin)">
              {{ formatPercentage(category.average_profit_margin) }}
            </div>
          </div>
        </div>

        <!-- Financial Metrics -->
        <div class="grid grid-cols-3 gap-3">
          <div>
            <div class="text-xs text-gray-500 mb-1">Doanh thu</div>
            <div class="font-bold text-blue-600 text-sm">{{ formatPrice(category.total_revenue) }}</div>
          </div>
          <div>
            <div class="text-xs text-gray-500 mb-1">Chi phí</div>
            <div class="font-bold text-orange-600 text-sm">{{ formatPrice(category.total_cost) }}</div>
          </div>
          <div>
            <div class="text-xs text-gray-500 mb-1">Lợi nhuận</div>
            <div class="font-bold text-sm" :class="getProfitColor(category.total_profit)">
              {{ formatPrice(category.total_profit) }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Category Profit Table - Desktop -->
    <div v-if="categoryProfits.length > 0" class="hidden md:block bg-white rounded-xl shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b">
            <tr>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Danh mục</th>
              <th class="px-4 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">Đơn hàng</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Doanh thu</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Chi phí</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Lợi nhuận</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">LN %</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="category in categoryProfits" :key="category.category" class="hover:bg-gray-50 transition-colors">
              <td class="px-4 py-3">
                <div class="font-bold text-gray-800">{{ category.category }}</div>
                <div class="text-xs text-gray-500">{{ category.item_count }} món</div>
              </td>
              <td class="px-4 py-3 text-center text-sm text-gray-900">{{ category.order_count }}</td>
              <td class="px-4 py-3 text-right font-medium text-blue-600">{{ formatPrice(category.total_revenue) }}</td>
              <td class="px-4 py-3 text-right font-medium text-orange-600">{{ formatPrice(category.total_cost) }}</td>
              <td class="px-4 py-3 text-right font-bold" :class="getProfitColor(category.total_profit)">
                {{ formatPrice(category.total_profit) }}
              </td>
              <td class="px-4 py-3 text-right font-bold text-lg" :class="getProfitMarginColor(category.average_profit_margin)">
                {{ formatPercentage(category.average_profit_margin) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { formatPrice, formatPercentage } from '../utils/formatters'

// Props
const props = defineProps({
  dateRange: {
    type: Object,
    required: true
  },
  categoryProfits: {
    type: Array,
    required: true
  }
})

// Helper functions
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const day = String(date.getDate()).padStart(2, '0')
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const year = date.getFullYear()
  return `${day}/${month}/${year}`
}

const getProfitMarginColor = (margin) => {
  if (margin < 0) return 'text-red-600'
  if (margin < 20) return 'text-yellow-600'
  return 'text-green-600'
}

const getProfitColor = (profit) => {
  if (profit < 0) return 'text-red-600'
  if (profit === 0) return 'text-gray-600'
  return 'text-green-600'
}
</script>
