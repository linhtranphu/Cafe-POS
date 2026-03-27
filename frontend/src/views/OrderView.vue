<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Pull to Refresh Indicator -->
    <PullToRefresh 
      :pull-distance="pullDistance" 
      :is-refreshing="isRefreshing"
      :threshold="80" />
    
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-4" style="padding-top: max(1rem, env(safe-area-inset-top))">
        <div class="flex items-center justify-between mb-4">
          <h1 class="text-2xl font-bold text-gray-900">📋 Orders</h1>
          <div class="flex gap-2">
            <!-- Merge Button (show when there are mergeable orders) -->
            <button 
              v-if="!selectionMode && mergeableOrdersCount >= 2"
              @click="enterSelectionMode"
              class="p-2.5 rounded-xl bg-purple-100 text-purple-600 active:bg-purple-200 transition-colors touch-manipulation">
              <span class="text-xl">🔗</span>
            </button>
            <button 
              @click="refreshOrders" 
              class="p-2.5 rounded-xl bg-gray-100 active:bg-gray-200 transition-colors touch-manipulation">
              <span class="text-xl">🔄</span>
            </button>
          </div>
        </div>
        
        <!-- Shift Tabs -->
        <div class="flex gap-2 mb-3">
          <button 
            @click="shiftFilter = 'current'"
            :class="[
              'flex-1 py-3 px-4 rounded-xl font-semibold text-sm transition-all touch-manipulation active:scale-98',
              shiftFilter === 'current' 
                ? 'bg-blue-500 text-white shadow-lg' 
                : 'bg-gray-100 text-gray-700 active:bg-gray-200'
            ]">
            <div class="flex items-center justify-center gap-1.5">
              <span>🔵</span>
              <span>Ca hiện tại</span>
              <span class="text-xs opacity-80">({{ currentShiftOrdersCount }})</span>
            </div>
          </button>
          <button 
            @click="shiftFilter = 'other'"
            :class="[
              'flex-1 py-3 px-4 rounded-xl font-semibold text-sm transition-all touch-manipulation active:scale-98',
              shiftFilter === 'other' 
                ? 'bg-blue-500 text-white shadow-lg' 
                : 'bg-gray-100 text-gray-700 active:bg-gray-200'
            ]">
            <div class="flex items-center justify-center gap-1.5">
              <span>📂</span>
              <span>Ca khác</span>
              <span class="text-xs opacity-80">({{ otherShiftsOrdersCount }})</span>
            </div>
          </button>
        </div>

        <!-- Payment Status Tabs -->
        <div v-if="!selectionMode" class="flex gap-2">
          <button 
            @click="paymentFilter = 'unpaid'"
            :class="[
              'flex-1 py-2.5 px-3 rounded-lg font-medium text-sm transition-all touch-manipulation active:scale-98',
              paymentFilter === 'unpaid' 
                ? 'bg-amber-500 text-white shadow-md' 
                : 'bg-gray-100 text-gray-700 active:bg-gray-200'
            ]">
            <div class="flex items-center justify-center gap-1.5">
              <span>📄</span>
              <span>Chưa thu tiền</span>
              <span class="text-xs opacity-80">({{ unpaidOrdersCount }})</span>
            </div>
          </button>
          <button 
            @click="paymentFilter = 'paid'"
            :class="[
              'flex-1 py-2.5 px-3 rounded-lg font-medium text-sm transition-all touch-manipulation active:scale-98',
              paymentFilter === 'paid' 
                ? 'bg-green-500 text-white shadow-md' 
                : 'bg-gray-100 text-gray-700 active:bg-gray-200'
            ]">
            <div class="flex items-center justify-center gap-1.5">
              <span>✅</span>
              <span>Đã thu tiền</span>
              <span class="text-xs opacity-80">({{ paidOrdersCount }})</span>
            </div>
          </button>
          <button 
            @click="paymentFilter = 'original'"
            :class="[
              'flex-1 py-2.5 px-3 rounded-lg font-medium text-sm transition-all touch-manipulation active:scale-98',
              paymentFilter === 'original' 
                ? 'bg-blue-500 text-white shadow-md' 
                : 'bg-gray-100 text-gray-700 active:bg-gray-200'
            ]">
            <div class="flex items-center justify-center gap-1.5">
              <span>📋</span>
              <span>Bill gốc</span>
              <span class="text-xs opacity-80">({{ originalOrdersCount }})</span>
            </div>
          </button>
        </div>

        <!-- Selection Mode Header -->
        <div v-else class="bg-purple-50 border border-purple-200 rounded-xl p-3">
          <div class="flex items-center justify-between mb-2">
            <span class="text-sm font-semibold text-purple-900">
              🔗 Chọn orders để gộp
            </span>
            <button 
              @click="exitSelectionMode"
              class="text-sm text-purple-600 font-medium">
              Hủy
            </button>
          </div>
          <div class="text-xs text-purple-700">
            Đã chọn: {{ selectedOrderIds.length }} orders
          </div>
        </div>
      </div>
    </div>

    <!-- Shift Warning -->
    <div v-if="!hasOpenShift" class="mx-4 mt-4 mb-4 bg-yellow-50 border-l-4 border-yellow-400 p-4 rounded-xl shadow-sm">
      <p class="text-yellow-800 text-sm font-semibold mb-2">⚠️ Chưa mở ca làm việc</p>
      <button 
        @click="$router.push('/shifts')" 
        class="bg-yellow-500 active:bg-yellow-600 text-white px-5 py-2.5 rounded-xl text-sm font-semibold shadow-md touch-manipulation active:scale-98 transition-transform">
        Mở ca ngay
      </button>
    </div>

    <!-- Orders List -->
    <div class="flex-1 overflow-y-auto px-4 py-2 pb-24">
      <div v-if="loading" class="text-center py-16">
        <div class="inline-block animate-spin rounded-full h-10 w-10 border-4 border-blue-500 border-t-transparent"></div>
        <p class="text-gray-500 text-sm mt-3">Đang tải...</p>
      </div>
      
      <div v-else-if="filteredOrders.length === 0" class="text-center py-20">
        <div class="text-7xl mb-4">📭</div>
        <p class="text-gray-500 text-base font-medium">
          {{ getEmptyMessage() }}
        </p>
      </div>
      
      <div v-else class="space-y-3">
        <div v-for="order in filteredOrders" :key="order.id" 
          @click="selectionMode ? toggleOrderSelection(order) : viewOrderDetail(order)"
          :class="[
            'bg-white rounded-2xl p-4 shadow-sm border active:scale-[0.98] active:shadow-md transition-all touch-manipulation',
            order.bill_printed ? 'border-2 border-amber-400 bg-amber-50' : 'border-gray-100',
            selectionMode && selectedOrderIds.includes(order.id) ? 'ring-2 ring-purple-500 bg-purple-50' : '',
            selectionMode && !isOrderMergeable(order) ? 'opacity-50' : ''
          ]">
          
          <!-- Selection Checkbox (in selection mode) -->
          <div v-if="selectionMode" class="absolute top-4 right-4">
            <div :class="[
              'w-6 h-6 rounded-full border-2 flex items-center justify-center',
              selectedOrderIds.includes(order.id) 
                ? 'bg-purple-500 border-purple-500' 
                : 'border-gray-300 bg-white'
            ]">
              <span v-if="selectedOrderIds.includes(order.id)" class="text-white text-sm">✓</span>
            </div>
          </div>
          
          <!-- Order Header -->
          <div class="flex justify-between items-start mb-3" :class="selectionMode ? 'pr-8' : ''">
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2">
                <h3 class="font-bold text-lg text-gray-900 truncate">{{ order.customer_name || 'Khách lẻ' }}</h3>
                <span v-if="order.bill_printed" class="text-xs bg-amber-500 text-white px-2 py-0.5 rounded-full font-semibold">
                  📄 Đã in
                </span>
              </div>
              <p class="text-xs text-gray-400 font-medium">{{ order.order_number }} · {{ formatTime(order.created_at) }}</p>
              <!-- Show shift info when viewing "other shifts" -->
              <p v-if="shiftFilter === 'other' && order.shift_id" class="text-xs text-purple-600 font-medium mt-1">
                📋 Ca: {{ order.shift_id.substring(0, 8) }}...
              </p>
            </div>
            <span :class="getStatusBadge(order.status)" class="px-3 py-1.5 rounded-full text-xs font-semibold whitespace-nowrap ml-2">
              {{ getStatusText(order.status) }}
            </span>
          </div>

          <!-- Items Summary -->
          <div class="mb-3 space-y-1">
            <div v-for="(item, idx) in order.items.slice(0, 2)" :key="idx" 
              class="flex justify-between text-sm">
              <span class="text-gray-700">
                {{ item.name }}
                <span v-if="item.variant_name" class="text-blue-600">({{ item.variant_name }})</span>
                <span class="text-gray-400"> x{{ item.quantity }}</span>
              </span>
              <span class="font-medium text-gray-900">{{ formatPrice(item.subtotal) }}</span>
            </div>
            <p v-if="order.items.length > 2" class="text-xs text-gray-400">
              +{{ order.items.length - 2 }} món khác...
            </p>
          </div>

          <!-- Total -->
          <div class="pt-3 border-t space-y-2">
            <div class="flex justify-between items-center">
              <span class="text-sm font-medium text-gray-600">Tổng cộng</span>
              <span class="text-lg font-bold text-green-600">{{ formatPrice(order.total) }}</span>
            </div>
            
            <!-- Payment Method -->
            <div v-if="order.payment_method" class="flex items-center justify-between">
              <span class="text-xs text-gray-500">Thanh toán:</span>
              <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs font-medium"
                :class="order.payment_method === 'cash' ? 'bg-green-100 text-green-700' : 'bg-blue-100 text-blue-700'">
                <span v-if="order.payment_method === 'cash'">💵 Tiền mặt</span>
                <span v-else-if="order.payment_method === 'transfer'">💳 CK</span>
                <span v-else>{{ order.payment_method }}</span>
              </span>
            </div>
          </div>

          <!-- Quick Actions -->
          <div class="mt-3 flex gap-2">
            <button v-if="isStatus(order, ORDER_STATUS.CREATED) && !order.bill_printed" 
              @click.stop="printTempBill(order.id)"
              class="flex-1 bg-amber-500 active:bg-amber-600 text-white py-3 rounded-xl text-sm font-semibold shadow-md touch-manipulation active:scale-98 transition-all">
              📄 In bill tạm
            </button>
            <button v-if="isStatus(order, ORDER_STATUS.CREATED)" 
              @click.stop="quickPayment(order)"
              class="flex-1 bg-green-500 active:bg-green-600 text-white py-3 rounded-xl text-sm font-semibold shadow-md touch-manipulation active:scale-98 transition-all">
              💰 Thu tiền
            </button>
            <button v-if="isStatus(order, ORDER_STATUS.PAID) && order.amount_due <= 0" 
              @click.stop="sendToBar(order.id)"
              class="flex-1 bg-blue-500 active:bg-blue-600 text-white py-3 rounded-xl text-sm font-semibold shadow-md touch-manipulation active:scale-98 transition-all">
              🍹 Gửi bar
            </button>
            <button v-if="isStatus(order, ORDER_STATUS.READY)" 
              @click.stop="serveOrder(order.id)"
              class="flex-1 bg-purple-500 active:bg-purple-600 text-white py-3 rounded-xl text-sm font-semibold shadow-md touch-manipulation active:scale-98 transition-all">
              🎉 Giao khách
            </button>
            <button v-if="isAnyStatus(order, [ORDER_STATUS.QUEUED, ORDER_STATUS.IN_PROGRESS])" 
              class="flex-1 bg-gray-200 text-gray-500 py-3 rounded-xl text-sm font-semibold cursor-not-allowed">
              ⏳ Đang pha...
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />

    <!-- Selection Mode Bottom Bar -->
    <div v-if="selectionMode" class="fixed bottom-0 left-0 right-0 bg-white border-t-2 border-purple-200 p-4 z-40 shadow-lg">
      <div class="flex gap-3">
        <button
          @click="exitSelectionMode"
          class="flex-1 px-4 py-3 text-gray-700 bg-gray-100 rounded-xl font-semibold active:bg-gray-200 transition-colors">
          Hủy
        </button>
        <button
          @click="openMergeModal"
          :disabled="selectedOrderIds.length < 2"
          :class="[
            'flex-1 px-4 py-3 text-white rounded-xl font-semibold transition-colors',
            selectedOrderIds.length >= 2
              ? 'bg-purple-500 active:bg-purple-600'
              : 'bg-gray-300 cursor-not-allowed'
          ]">
          🔗 Gộp {{ selectedOrderIds.length }} orders
        </button>
      </div>
    </div>

    <!-- Merge Bills Modal -->
    <MergeBillsModal
      :show="showMergeModal"
      :selected-orders="selectedOrdersForMerge"
      @close="closeMergeModal"
      @merged="handleMerged"
    />

    <!-- Order Detail Modal -->
    <transition name="slide-up">
      <div v-if="selectedOrder" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full max-h-[85vh] overflow-y-auto">
          <div class="sticky top-0 bg-white px-4 py-4 border-b flex justify-between items-center">
            <h3 class="text-lg font-bold">Chi tiết Order</h3>
            <button @click="selectedOrder = null" class="text-2xl text-gray-400">×</button>
          </div>
          
          <div class="px-4 py-4">
            <!-- Order Info -->
            <div class="mb-4">
              <h4 class="text-2xl font-bold mb-1">{{ selectedOrder.order_number }}</h4>
              <p class="text-gray-600">{{ selectedOrder.customer_name || 'Khách lẻ' }}</p>
              <p class="text-sm text-gray-400">{{ formatDate(selectedOrder.created_at) }}</p>
              <span :class="getStatusBadge(selectedOrder.status)" 
                class="inline-block mt-2 px-3 py-1 rounded-full text-xs font-medium">
                {{ getStatusText(selectedOrder.status) }}
              </span>
            </div>

            <!-- Items -->
            <div class="mb-4">
              <h5 class="font-bold mb-2">Món đã order</h5>
              <div class="space-y-2">
                <div v-for="item in selectedOrder.items" :key="item.menu_item_id" 
                  class="flex justify-between bg-gray-50 p-3 rounded-lg">
                  <div>
                    <div class="font-medium">
                      {{ item.name }}
                      <span v-if="item.variant_name" class="text-gray-500">({{ item.variant_name }})</span>
                    </div>
                    <div class="text-sm text-gray-500">{{ formatPrice(item.price) }} x {{ item.quantity }}</div>
                  </div>
                  <div class="font-bold">{{ formatPrice(item.subtotal) }}</div>
                </div>
              </div>
            </div>

            <!-- Note -->
            <div v-if="selectedOrder.note" class="mb-4 p-3 bg-yellow-50 rounded-lg">
              <p class="text-sm text-gray-700">📝 {{ selectedOrder.note }}</p>
            </div>

            <!-- Total -->
            <div class="mb-4 p-4 bg-gray-50 rounded-lg">
              <div class="flex justify-between text-lg font-bold">
                <span>Tổng cộng</span>
                <span class="text-green-600">{{ formatPrice(selectedOrder.total) }}</span>
              </div>
              <div v-if="selectedOrder.amount_paid > 0" class="flex justify-between text-sm text-gray-600 mt-1">
                <span>Đã thu</span>
                <span>{{ formatPrice(selectedOrder.amount_paid) }}</span>
              </div>
              
              <!-- Payment Method -->
              <div v-if="selectedOrder.payment_method" class="mt-3 pt-3 border-t border-gray-200">
                <div class="flex items-center gap-2">
                  <span class="text-sm text-gray-600">Phương thức thanh toán:</span>
                  <span class="inline-flex items-center gap-1 px-2 py-1 rounded-lg text-sm font-medium"
                    :class="selectedOrder.payment_method === 'cash' ? 'bg-green-100 text-green-700' : 'bg-blue-100 text-blue-700'">
                    <span v-if="selectedOrder.payment_method === 'cash'">💵 Tiền mặt</span>
                    <span v-else-if="selectedOrder.payment_method === 'transfer'">💳 Chuyển khoản</span>
                    <span v-else>{{ selectedOrder.payment_method }}</span>
                  </span>
                </div>
              </div>
            </div>

            <!-- Actions -->
            <div class="space-y-2">
              <button v-if="isStatus(selectedOrder, ORDER_STATUS.CREATED) && !selectedOrder.bill_printed" 
                @click="printTempBill(selectedOrder.id)"
                class="w-full bg-amber-500 text-white py-3 rounded-xl font-medium active:bg-amber-600">
                📄 In bill tạm
              </button>
              <button v-if="isStatus(selectedOrder, ORDER_STATUS.CREATED)" 
                @click="showPaymentModal(selectedOrder)"
                class="w-full bg-green-500 text-white py-3 rounded-xl font-medium active:bg-green-600">
                💰 Thu tiền
              </button>
              <button v-if="isStatus(selectedOrder, ORDER_STATUS.CREATED)" 
                @click="editOrder(selectedOrder)"
                class="w-full bg-blue-500 text-white py-3 rounded-xl font-medium active:bg-blue-600">
                ✏️ Chỉnh sửa
              </button>
              <button v-if="isStatus(selectedOrder, ORDER_STATUS.CREATED)" 
                @click="confirmCancelOrder(selectedOrder)"
                class="w-full bg-red-500 text-white py-3 rounded-xl font-medium active:bg-red-600">
                ❌ Hủy order
              </button>
              <button v-if="isStatus(selectedOrder, ORDER_STATUS.PAID) && selectedOrder.amount_due <= 0" 
                @click="sendToBar(selectedOrder.id)"
                class="w-full bg-blue-500 text-white py-3 rounded-xl font-medium active:bg-blue-600">
                🍹 Gửi quầy bar
              </button>
              <button v-if="isStatus(selectedOrder, ORDER_STATUS.READY)" 
                @click="serveOrder(selectedOrder.id)"
                class="w-full bg-purple-500 text-white py-3 rounded-xl font-medium active:bg-purple-600">
                🎉 Giao cho khách
              </button>
              <div v-if="isAnyStatus(selectedOrder, [ORDER_STATUS.QUEUED, ORDER_STATUS.IN_PROGRESS])" 
                class="w-full bg-gray-100 text-gray-600 py-3 rounded-xl font-medium text-center">
                ⏳ Barista đang pha chế...
              </div>
            </div>

            <!-- Reprint Section (Only for Cashier and Manager) -->
            <div v-if="canReprint" class="mt-6 pt-4 border-t">
              <h5 class="font-bold mb-3 text-gray-700">🖨️ In lại</h5>
              <div class="space-y-2">
                <!-- Reprint Bill Button -->
                <button
                  @click="handleReprintBill(selectedOrder.id)"
                  :disabled="reprintingBill"
                  class="w-full bg-purple-500 text-white py-3 rounded-xl font-medium active:bg-purple-600 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                >
                  <span v-if="reprintingBill" class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
                  <span v-else>🧾</span>
                  <span>{{ reprintingBill ? 'Đang in...' : 'In lại Bill' }}</span>
                </button>

                <!-- Reprint Labels for Each Item -->
                <div class="space-y-2">
                  <div
                    v-for="(item, index) in selectedOrder.items"
                    :key="index"
                    class="flex gap-2"
                  >
                    <button
                      @click="handleReprintLabel(selectedOrder.id, index)"
                      :disabled="reprintingLabel === index"
                      class="flex-1 bg-orange-500 text-white py-2 rounded-lg text-sm font-medium active:bg-orange-600 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                    >
                      <span v-if="reprintingLabel === index" class="inline-block w-3 h-3 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
                      <span v-else>🏷️</span>
                      <span>{{ reprintingLabel === index ? 'Đang in...' : `In tem: ${item.name}` }}</span>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- Payment Modal -->
    <transition name="slide-up">
      <div v-if="showPayment" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6">
          <h3 class="text-xl font-bold mb-4">💰 Thu tiền</h3>
          
          <div class="mb-4 p-4 bg-gray-50 rounded-lg">
            <div class="flex justify-between mb-2">
              <span class="text-gray-600">Tổng tiền</span>
              <span class="text-xl font-bold text-green-600">{{ formatPrice(paymentOrder?.total) }}</span>
            </div>
            <div v-if="paymentOrder?.amount_paid > 0" class="flex justify-between text-sm text-gray-600">
              <span>Đã thu</span>
              <span>{{ formatPrice(paymentOrder?.amount_paid) }}</span>
            </div>
          </div>

          <div class="mb-4">
            <label class="block text-sm font-medium mb-2">Số tiền thu</label>
            <input v-model.number="paymentAmount" 
              type="number" 
              step="1000"
              class="w-full px-4 py-3 text-lg font-bold border rounded-lg focus:ring-2 focus:ring-green-500">
          </div>

          <div class="mb-4">
            <label class="block text-sm font-medium mb-2">Phương thức</label>
            <div class="grid grid-cols-2 gap-2">
              <button v-for="method in paymentMethods" :key="method.value"
                @click="paymentMethod = method.value"
                :class="[
                  'py-3 rounded-lg font-medium',
                  paymentMethod === method.value 
                    ? 'bg-green-500 text-white' 
                    : 'bg-gray-100 text-gray-700'
                ]">
                {{ method.icon }} {{ method.label }}
              </button>
            </div>
          </div>

          <div class="flex gap-2">
            <button @click="showPayment = false" 
              class="flex-1 bg-gray-200 text-gray-700 py-3 rounded-xl font-medium">
              Hủy
            </button>
            <button @click="processPayment" 
              class="flex-1 bg-green-500 text-white py-3 rounded-xl font-medium">
              Xác nhận
            </button>
          </div>
        </div>
      </div>
    </transition>

    <!-- Cancel Order Modal -->
    <transition name="slide-up">
      <div v-if="showCancelModal" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6">
          <h3 class="text-xl font-bold mb-4 text-red-600">❌ Hủy Order</h3>
          
          <div class="mb-4 p-4 bg-red-50 rounded-lg border border-red-200">
            <p class="text-sm text-red-800 mb-2">
              <strong>Cảnh báo:</strong> Bạn đang hủy order này. Hành động này không thể hoàn tác.
            </p>
            <div class="text-sm text-gray-700">
              <p><strong>Order:</strong> {{ cancelingOrder?.order_number }}</p>
              <p><strong>Khách:</strong> {{ cancelingOrder?.customer_name || 'Khách lẻ' }}</p>
              <p><strong>Tổng tiền:</strong> {{ formatPrice(cancelingOrder?.total) }}</p>
            </div>
          </div>

          <div class="mb-4">
            <label class="block text-sm font-medium mb-2 text-gray-700">
              Lý do hủy <span class="text-red-500">*</span>
            </label>
            <textarea 
              v-model="cancelReason" 
              rows="3"
              placeholder="Nhập lý do hủy order (bắt buộc)..."
              class="w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-red-500 resize-none"
            ></textarea>
          </div>

          <div class="flex gap-2">
            <button @click="closeCancelModal" 
              class="flex-1 bg-gray-200 text-gray-700 py-3 rounded-xl font-medium active:bg-gray-300">
              Quay lại
            </button>
            <button 
              @click="processCancelOrder" 
              :disabled="!cancelReason.trim()"
              class="flex-1 bg-red-500 text-white py-3 rounded-xl font-medium active:bg-red-600 disabled:opacity-50 disabled:cursor-not-allowed">
              Xác nhận hủy
            </button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useOrderStore } from '../stores/order'
import { useShiftStore } from '../stores/shift'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import MergeBillsModal from '../components/MergeBillsModal.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'
import { printJobService } from '../services/printJob'
import { 
  ORDER_STATUS, 
  PAYMENT_METHOD, 
  PAYMENT_METHOD_DISPLAY,
  ORDER_STATUS_DISPLAY,
  STATUS_FILTER_OPTIONS
} from '../constants/order'

const router = useRouter()
const orderStore = useOrderStore()
const shiftStore = useShiftStore()
const authStore = useAuthStore()

// State
const shiftFilter = ref('current') // 'current' or 'other'
const paymentFilter = ref('unpaid') // 'unpaid' or 'paid'
const filterStatus = ref('ALL')
const selectedOrder = ref(null)
const showPayment = ref(false)
const paymentOrder = ref(null)
const paymentAmount = ref(0)
const paymentMethod = ref(PAYMENT_METHOD.CASH)

// Merge Bills State
const selectionMode = ref(false)
const selectedOrderIds = ref([])
const showMergeModal = ref(false)

// Cancel Order State
const showCancelModal = ref(false)
const cancelReason = ref('')
const cancelingOrder = ref(null)

// Reprint State
const reprintingBill = ref(false)
const reprintingLabel = ref(null)

// Check if user can reprint (only cashier and manager)
const canReprint = computed(() => {
  const role = authStore.user?.role
  return role === 'cashier' || role === 'manager'
})

// Data
const statuses = STATUS_FILTER_OPTIONS

const paymentMethods = PAYMENT_METHOD_DISPLAY

// Computed
const loading = computed(() => orderStore.loading)
const orders = computed(() => orderStore.orders)
const hasOpenShift = computed(() => shiftStore.hasOpenShift)

// Filter orders by shift first
const ordersByShift = computed(() => {
  const currentShiftId = shiftStore.currentShift?.id
  
  if (shiftFilter.value === 'current') {
    // Show orders from current shift only
    if (!currentShiftId) return []
    return orders.value.filter(o => o.shift_id === currentShiftId)
  } else {
    // Show orders from other shifts (not current shift)
    if (!currentShiftId) return orders.value // If no current shift, show all
    return orders.value.filter(o => o.shift_id !== currentShiftId)
  }
})

// Then filter by payment status
const ordersByPaymentStatus = computed(() => {
  if (paymentFilter.value === 'unpaid') {
    // Chưa thu tiền = tất cả orders chưa thanh toán (status = CREATED)
    return ordersByShift.value.filter(o => o.status === ORDER_STATUS.CREATED)
  } else if (paymentFilter.value === 'paid') {
    // Đã thu tiền = status = PAID
    return ordersByShift.value.filter(o => o.status === ORDER_STATUS.PAID)
  } else if (paymentFilter.value === 'original') {
    // Bill gốc = orders không có merged_into_id (chưa bị gộp vào bill khác)
    return ordersByShift.value.filter(o => !o.merged_into_id)
  }
  return ordersByShift.value
})

// Then filter by status (if needed)
const filteredOrders = computed(() => {
  if (filterStatus.value === 'ALL') return ordersByPaymentStatus.value
  return ordersByPaymentStatus.value.filter(o => o.status === filterStatus.value)
})

// Count orders for each tab
const currentShiftOrdersCount = computed(() => {
  const currentShiftId = shiftStore.currentShift?.id
  if (!currentShiftId) return 0
  return orders.value.filter(o => o.shift_id === currentShiftId).length
})

const otherShiftsOrdersCount = computed(() => {
  const currentShiftId = shiftStore.currentShift?.id
  if (!currentShiftId) return orders.value.length
  return orders.value.filter(o => o.shift_id !== currentShiftId).length
})

// Count unpaid and paid orders in current shift filter
const unpaidOrdersCount = computed(() => {
  const currentShiftId = shiftStore.currentShift?.id
  let shiftOrders = []
  
  if (shiftFilter.value === 'current') {
    if (!currentShiftId) return 0
    shiftOrders = orders.value.filter(o => o.shift_id === currentShiftId)
  } else {
    if (!currentShiftId) {
      shiftOrders = orders.value
    } else {
      shiftOrders = orders.value.filter(o => o.shift_id !== currentShiftId)
    }
  }
  
  // Chưa thu tiền = tất cả orders có status = CREATED
  return shiftOrders.filter(o => o.status === ORDER_STATUS.CREATED).length
})

const paidOrdersCount = computed(() => {
  const currentShiftId = shiftStore.currentShift?.id
  let shiftOrders = []
  
  if (shiftFilter.value === 'current') {
    if (!currentShiftId) return 0
    shiftOrders = orders.value.filter(o => o.shift_id === currentShiftId)
  } else {
    if (!currentShiftId) {
      shiftOrders = orders.value
    } else {
      shiftOrders = orders.value.filter(o => o.shift_id !== currentShiftId)
    }
  }
  
  return shiftOrders.filter(o => o.status === ORDER_STATUS.PAID).length
})

// Count original orders (not merged into another order)
const originalOrdersCount = computed(() => {
  const currentShiftId = shiftStore.currentShift?.id
  let shiftOrders = []
  
  if (shiftFilter.value === 'current') {
    if (!currentShiftId) return 0
    shiftOrders = orders.value.filter(o => o.shift_id === currentShiftId)
  } else {
    if (!currentShiftId) {
      shiftOrders = orders.value
    } else {
      shiftOrders = orders.value.filter(o => o.shift_id !== currentShiftId)
    }
  }
  
  // Bill gốc = orders không có merged_into_id
  return shiftOrders.filter(o => !o.merged_into_id).length
})

// Helper to check order status
const isStatus = (order, status) => order.status === status
const isAnyStatus = (order, statuses) => statuses.includes(order.status)

// Merge Bills Computed Properties
const mergeableOrders = computed(() => {
  // Only orders in current shift that are CREATED (unpaid) can be merged
  const currentShiftId = shiftStore.currentShift?.id
  if (!currentShiftId) return []
  
  return orders.value.filter(o => 
    o.shift_id === currentShiftId &&
    o.status === ORDER_STATUS.CREATED
  )
})

const mergeableOrdersCount = computed(() => mergeableOrders.value.length)

const selectedOrdersForMerge = computed(() => {
  return orders.value.filter(o => selectedOrderIds.value.includes(o.id))
})

const isOrderMergeable = (order) => {
  return mergeableOrders.value.some(o => o.id === order.id)
}

// Methods
const refreshOrders = async () => {
  await orderStore.fetchOrders()
}

// Pull to refresh
const { pullDistance, isRefreshing } = usePullToRefresh(refreshOrders)

const getOrderCountByStatus = (status) => {
  if (status === 'ALL') return ordersByShift.value.length
  return ordersByShift.value.filter(o => o.status === status).length
}

const getStatusBadge = (status) => {
  return ORDER_STATUS_DISPLAY[status]?.badge || 'bg-gray-100 text-gray-800'
}

const getStatusText = (status) => {
  return ORDER_STATUS_DISPLAY[status]?.label || status
}

const formatPrice = (price) => {
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(price)
}

const formatTime = (date) => {
  const d = new Date(date)
  return d.toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' })
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('vi-VN')
}

const getEmptyMessage = () => {
  if (paymentFilter.value === 'unpaid') {
    return 'Chưa có order nào chưa thu tiền'
  } else if (paymentFilter.value === 'paid') {
    return 'Chưa có order nào đã thu tiền'
  } else if (paymentFilter.value === 'original') {
    return 'Chưa có bill gốc nào (chưa gộp)'
  }
  return 'Chưa có order nào'
}

const viewOrderDetail = (order) => {
  selectedOrder.value = order
}

const printTempBill = async (orderId) => {
  try {
    await orderStore.printTemporaryBill(orderId)
    alert('✅ Đã in tem thành công')
  } catch (error) {
    alert('❌ Lỗi: ' + error.message)
  }
}

const quickPayment = (order) => {
  paymentOrder.value = order
  paymentAmount.value = order.amount_due || order.total
  paymentMethod.value = PAYMENT_METHOD.CASH
  showPayment.value = true
  selectedOrder.value = null
}

const showPaymentModal = (order) => {
  paymentOrder.value = order
  paymentAmount.value = order.amount_due || order.total
  paymentMethod.value = PAYMENT_METHOD.CASH
  showPayment.value = true
  selectedOrder.value = null
}

const processPayment = async () => {
  try {
    await orderStore.collectPayment(paymentOrder.value.id, {
      amount: paymentAmount.value,
      payment_method: paymentMethod.value
    })
    showPayment.value = false
    paymentOrder.value = null
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const sendToBar = async (orderId) => {
  try {
    await orderStore.sendToBar(orderId)
    selectedOrder.value = null
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const serveOrder = async (orderId) => {
  try {
    await orderStore.serveOrder(orderId)
    selectedOrder.value = null
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const editOrder = (order) => {
  // TODO: Implement edit order functionality
  alert('Chức năng chỉnh sửa order đang được phát triển')
}

const confirmCancelOrder = (order) => {
  cancelingOrder.value = order
  cancelReason.value = ''
  showCancelModal.value = true
  selectedOrder.value = null
}

const closeCancelModal = () => {
  showCancelModal.value = false
  cancelingOrder.value = null
  cancelReason.value = ''
}

const processCancelOrder = async () => {
  if (!cancelReason.value.trim()) {
    alert('Vui lòng nhập lý do hủy order')
    return
  }

  try {
    await orderStore.cancelOrder(cancelingOrder.value.id, cancelReason.value.trim())
    showCancelModal.value = false
    cancelingOrder.value = null
    cancelReason.value = ''
    alert('✅ Đã hủy order thành công')
  } catch (error) {
    alert('❌ Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

// Reprint Functions
const handleReprintBill = async (orderId) => {
  if (reprintingBill.value) return

  reprintingBill.value = true
  try {
    await printJobService.reprintBill(orderId)
    alert('✅ Đã gửi lệnh in lại bill')
  } catch (error) {
    console.error('Reprint bill error:', error)
    alert('❌ Lỗi in lại bill: ' + (error.response?.data?.error || error.message))
  } finally {
    reprintingBill.value = false
  }
}

const handleReprintLabel = async (orderId, itemIndex) => {
  if (reprintingLabel.value !== null) return

  reprintingLabel.value = itemIndex
  try {
    await printJobService.reprintLabel(orderId, itemIndex)
    alert('✅ Đã gửi lệnh in lại tem')
  } catch (error) {
    console.error('Reprint label error:', error)
    alert('❌ Lỗi in lại tem: ' + (error.response?.data?.error || error.message))
  } finally {
    reprintingLabel.value = null
  }
}

// Merge Bills Methods
const enterSelectionMode = () => {
  selectionMode.value = true
  selectedOrderIds.value = []
}

const exitSelectionMode = () => {
  selectionMode.value = false
  selectedOrderIds.value = []
}

const toggleOrderSelection = (order) => {
  // Only allow selecting mergeable orders
  if (!isOrderMergeable(order)) return
  
  const index = selectedOrderIds.value.indexOf(order.id)
  if (index > -1) {
    selectedOrderIds.value.splice(index, 1)
  } else {
    selectedOrderIds.value.push(order.id)
  }
}

const openMergeModal = () => {
  if (selectedOrderIds.value.length < 2) {
    alert('⚠️ Vui lòng chọn ít nhất 2 orders để gộp')
    return
  }
  showMergeModal.value = true
}

const closeMergeModal = () => {
  showMergeModal.value = false
}

const handleMerged = (response) => {
  // Exit selection mode
  exitSelectionMode()
  
  // Show success message
  alert(`✅ ${response.message}`)
  
  // Refresh orders
  refreshOrders()
}

// Lifecycle
onMounted(async () => {
  await Promise.all([
    shiftStore.fetchCurrentShift(),
    orderStore.fetchOrders()
  ])
})
</script>

<style scoped>
.touch-manipulation {
  touch-action: manipulation;
  -webkit-tap-highlight-color: transparent;
}

.scrollbar-hide::-webkit-scrollbar {
  display: none;
}

.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.active\:scale-95:active {
  transform: scale(0.95);
}

.active\:scale-98:active {
  transform: scale(0.98);
}

.active\:scale-\[0\.98\]:active {
  transform: scale(0.98);
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: transform 0.3s ease;
}

.slide-up-enter-from {
  transform: translateY(100%);
}

.slide-up-leave-to {
  transform: translateY(100%);
}

/* Improve button press feedback */
button:active {
  transition: all 0.1s ease;
}

/* Smooth scrolling */
.overflow-y-auto {
  -webkit-overflow-scrolling: touch;
}

/* Line clamp for menu item names */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.25rem;
}
</style>
