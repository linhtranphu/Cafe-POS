<template>
  <div>
    <!-- Empty State -->
    <div v-if="!operatingProfit" class="text-center py-16">
      <div class="text-6xl mb-4">📭</div>
      <p class="text-gray-700 font-medium mb-2">Không có dữ liệu</p>
      <p class="text-gray-500 text-sm">Không có dữ liệu lợi nhuận trong khoảng thời gian này</p>
    </div>

    <div v-else class="space-y-4">
      <!-- Date Range Display -->
      <div class="bg-blue-50 rounded-xl p-4">
        <div class="text-xs md:text-sm text-blue-600 font-medium mb-1">Khoảng thời gian</div>
        <div class="text-sm md:text-base text-gray-800">
          {{ formatDate(operatingProfit.date_range?.start) }} → {{ formatDate(operatingProfit.date_range?.end) }}
        </div>
      </div>

      <!-- Gross Profit Section -->
      <div class="bg-white rounded-2xl p-4 md:p-6 shadow-sm">
        <h3 class="font-bold text-lg md:text-xl text-gray-800 mb-3 md:mb-4 flex items-center gap-2">
          <span>💰</span>
          <span>Lợi nhuận gộp</span>
        </h3>
        
        <div class="space-y-3">
          <div class="flex justify-between items-center">
            <span class="text-sm md:text-base text-gray-600">Doanh thu</span>
            <span class="font-bold text-sm md:text-base text-blue-600">{{ formatPrice(operatingProfit.total_revenue) }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-sm md:text-base text-gray-600">Chi phí hàng bán (COGS)</span>
            <span class="font-bold text-sm md:text-base text-orange-600">{{ formatPrice(operatingProfit.total_cogs) }}</span>
          </div>
          <div class="pt-3 border-t flex justify-between items-center">
            <span class="text-sm md:text-base font-medium text-gray-800">Lợi nhuận gộp</span>
            <div class="text-right">
              <div class="font-bold text-lg md:text-xl" :class="getProfitColor(operatingProfit.gross_profit)">
                {{ formatPrice(operatingProfit.gross_profit) }}
              </div>
              <div class="text-xs md:text-sm" :class="getProfitMarginColor(operatingProfit.gross_profit_margin)">
                {{ formatPercentage(operatingProfit.gross_profit_margin) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Operating Expenses Section -->
      <div class="bg-white rounded-2xl p-4 md:p-6 shadow-sm">
        <div class="flex items-center justify-between mb-3 md:mb-4 flex-wrap gap-2">
          <h3 class="font-bold text-lg md:text-xl text-gray-800 flex items-center gap-2">
            <span>💼</span>
            <span>Chi phí vận hành</span>
          </h3>
          <!-- Expense Allocated Indicator -->
          <div v-if="operatingProfit.expense_allocated" 
            class="text-xs bg-yellow-100 text-yellow-700 px-2 py-1 rounded-full">
            ⚠️ Phân bổ
          </div>
        </div>
        
        <div class="space-y-2">
          <div class="flex justify-between items-center">
            <span class="text-sm md:text-base text-gray-600">👥 Lương nhân viên</span>
            <span class="font-medium text-sm md:text-base text-gray-800">{{ formatPrice(operatingProfit.staff_salary) }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-sm md:text-base text-gray-600">🏢 Tiền thuê mặt bằng</span>
            <span class="font-medium text-sm md:text-base text-gray-800">{{ formatPrice(operatingProfit.rent) }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-sm md:text-base text-gray-600">⚡ Điện nước</span>
            <span class="font-medium text-sm md:text-base text-gray-800">{{ formatPrice(operatingProfit.utilities) }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-sm md:text-base text-gray-600">📢 Marketing</span>
            <span class="font-medium text-sm md:text-base text-gray-800">{{ formatPrice(operatingProfit.marketing_costs) }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-sm md:text-base text-gray-600">📦 Chi phí khác</span>
            <span class="font-medium text-sm md:text-base text-gray-800">{{ formatPrice(operatingProfit.other_expenses) }}</span>
          </div>
          <div class="pt-3 border-t flex justify-between items-center">
            <span class="text-sm md:text-base font-medium text-gray-800">Tổng chi phí vận hành</span>
            <span class="font-bold text-sm md:text-base text-red-600">{{ formatPrice(operatingProfit.total_expenses) }}</span>
          </div>
        </div>

        <!-- Allocation Note -->
        <div v-if="operatingProfit.allocation_note" class="mt-3 pt-3 border-t">
          <div class="text-xs md:text-sm text-yellow-700 bg-yellow-50 p-2 md:p-3 rounded-lg">
            ℹ️ {{ operatingProfit.allocation_note }}
          </div>
        </div>
      </div>

      <!-- Operating Profit Section -->
      <div class="bg-gradient-to-br from-purple-500 to-blue-500 rounded-2xl p-4 md:p-6 shadow-lg text-white">
        <h3 class="font-bold text-lg md:text-xl mb-3 md:mb-4 flex items-center gap-2">
          <span>📈</span>
          <span>Lợi nhuận vận hành</span>
        </h3>
        
        <div class="space-y-2">
          <div class="flex justify-between items-center text-sm md:text-base opacity-90">
            <span>Lợi nhuận gộp</span>
            <span>{{ formatPrice(operatingProfit.gross_profit) }}</span>
          </div>
          <div class="flex justify-between items-center text-sm md:text-base opacity-90">
            <span>Chi phí vận hành</span>
            <span>- {{ formatPrice(operatingProfit.total_expenses) }}</span>
          </div>
          <div class="pt-3 border-t border-white/20 flex justify-between items-center">
            <span class="font-medium text-sm md:text-base">Lợi nhuận vận hành</span>
            <div class="text-right">
              <div class="font-bold text-xl md:text-3xl">
                {{ formatPrice(operatingProfit.operating_profit) }}
              </div>
              <div class="text-sm md:text-base opacity-90">
                {{ formatPercentage(operatingProfit.operating_profit_margin) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- No Expenses Warning -->
      <div v-if="operatingProfit.total_expenses === 0" class="bg-yellow-50 rounded-xl p-4">
        <div class="flex items-start gap-3">
          <span class="text-2xl flex-shrink-0">⚠️</span>
          <div>
            <div class="font-medium text-yellow-800 mb-1 text-sm md:text-base">Chưa nhập chi phí vận hành</div>
            <div class="text-xs md:text-sm text-yellow-700">
              Hiện tại chỉ hiển thị lợi nhuận gộp. Vui lòng nhập chi phí vận hành để xem lợi nhuận vận hành chính xác.
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { formatPrice, formatPercentage } from '../utils/formatters'

// Props
const props = defineProps({
  operatingProfit: {
    type: Object,
    default: null
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

const getProfitColor = (profit) => {
  if (profit < 0) return 'text-red-600'
  if (profit === 0) return 'text-gray-600'
  return 'text-green-600'
}

const getProfitMarginColor = (margin) => {
  if (margin < 0) return 'text-red-600'
  if (margin < 20) return 'text-yellow-600'
  return 'text-green-600'
}
</script>
