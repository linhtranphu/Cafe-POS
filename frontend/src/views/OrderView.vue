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
          <button 
            @click="refreshOrders" 
            class="p-2.5 rounded-xl bg-gray-100 active:bg-gray-200 transition-colors touch-manipulation">
            <span class="text-xl">🔄</span>
          </button>
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
        
        <!-- Status Filter Pills -->
        <div class="flex gap-2 overflow-x-auto pb-2 scrollbar-hide -mx-4 px-4">
          <button v-for="status in statuses" :key="status.value" 
            @click="filterStatus = status.value"
            :class="[
              'px-4 py-2.5 rounded-full text-sm font-semibold whitespace-nowrap transition-all touch-manipulation active:scale-95',
              filterStatus === status.value 
                ? 'bg-green-500 text-white shadow-lg' 
                : 'bg-gray-100 text-gray-700 active:bg-gray-200'
            ]">
            <span class="mr-1">{{ status.icon }}</span>
            <span>{{ status.label }}</span>
            <span class="ml-1.5 text-xs opacity-80">({{ getOrderCountByStatus(status.value) }})</span>
          </button>
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
          {{ shiftFilter === 'current' ? 'Chưa có order nào trong ca này' : 'Không có order từ ca khác' }}
        </p>
      </div>
      
      <div v-else class="space-y-3">
        <div v-for="order in filteredOrders" :key="order.id" 
          @click="viewOrderDetail(order)"
          class="bg-white rounded-2xl p-4 shadow-sm border border-gray-100 active:scale-[0.98] active:shadow-md transition-all touch-manipulation">
          
          <!-- Order Header -->
          <div class="flex justify-between items-start mb-3">
            <div class="flex-1 min-w-0">
              <h3 class="font-bold text-lg text-gray-900 truncate">{{ order.order_number }}</h3>
              <p class="text-sm text-gray-600 font-medium">{{ order.customer_name || 'Khách lẻ' }}</p>
              <p class="text-xs text-gray-400 mt-0.5">{{ formatTime(order.created_at) }}</p>
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

    <!-- Floating Action Button -->
    <button v-if="hasOpenShift" 
      @click="startNewOrder"
      class="fixed bottom-24 right-5 w-16 h-16 bg-blue-500 active:bg-blue-600 text-white rounded-full shadow-2xl flex items-center justify-center text-3xl touch-manipulation active:scale-95 transition-all z-30">
      ➕
    </button>

    <!-- Bottom Navigation -->
    <BottomNav />

    <!-- Create Order - Full Screen -->
    <transition name="slide-up">
      <div v-if="showCreateOrder" class="fixed inset-0 bg-white z-50 overflow-hidden flex flex-col">
        <!-- Header -->
        <div class="bg-blue-500 text-white px-4 py-4 flex items-center justify-between shadow-lg" style="padding-top: max(1rem, env(safe-area-inset-top))">
          <button @click="cancelCreateOrder" class="text-2xl p-2 active:bg-blue-600 rounded-lg touch-manipulation">←</button>
          <h2 class="text-xl font-bold">Tạo Order Mới</h2>
          <button @click="confirmOrder" :disabled="cart.length === 0" 
            class="text-sm font-semibold px-4 py-2.5 bg-white text-blue-500 rounded-xl disabled:opacity-50 disabled:cursor-not-allowed active:scale-95 transition-all touch-manipulation shadow-md">
            Xác nhận
          </button>
        </div>

        <!-- Category Tabs -->
        <div class="flex gap-2 px-4 py-3 overflow-x-auto bg-white border-b scrollbar-hide -mx-4 px-4">
          <button v-for="cat in categories" :key="cat.id"
            @click="selectedCategory = cat.id"
            :class="[
              'px-4 py-2.5 rounded-full text-sm font-semibold whitespace-nowrap touch-manipulation active:scale-95 transition-all',
              selectedCategory === cat.id 
                ? 'bg-blue-500 text-white shadow-lg' 
                : 'bg-gray-100 text-gray-700 active:bg-gray-200'
            ]">
            <span class="mr-1">{{ cat.icon }}</span>
            <span>{{ cat.name }}</span>
          </button>
        </div>

        <!-- Menu Items Grid -->
        <div class="flex-1 overflow-y-auto px-4 py-4 bg-gray-50">
          <div class="grid grid-cols-2 gap-3">
            <template v-for="item in filteredMenuItems" :key="item.id">
              <!-- Single-size item - direct add -->
              <button v-if="!item.has_variants"
                @click="addToCart(item)"
                class="bg-white rounded-2xl p-4 shadow-sm border border-gray-100 active:scale-95 active:shadow-md transition-all text-left touch-manipulation h-full flex flex-col">
                <div class="flex-1">
                  <div class="font-semibold text-gray-900 mb-2 line-clamp-2 min-h-[2.5rem]">{{ item.name }}</div>
                  <div class="text-base font-bold text-blue-600">{{ formatPrice(item.price) }}</div>
                </div>
                <div v-if="getCartItemQty(item.id) > 0" 
                  class="mt-3 bg-blue-500 text-white text-xs font-semibold px-3 py-1.5 rounded-full inline-flex items-center justify-center gap-1 shadow-md">
                  <span>🛒</span>
                  <span>{{ getCartItemQty(item.id) }} món</span>
                </div>
              </button>

              <!-- Multi-size item - show variants -->
              <div v-else
                class="bg-white rounded-2xl p-4 shadow-sm border border-gray-100 text-left flex flex-col h-full">
                <div class="font-semibold text-gray-900 mb-3 line-clamp-2 min-h-[2.5rem]">{{ item.name }}</div>
                <div class="space-y-2 flex-1">
                  <button v-for="variant in item.variants" :key="variant.id"
                    @click="addToCart(item, variant)"
                    class="w-full flex justify-between items-center px-3 py-2.5 bg-gray-50 rounded-xl active:bg-blue-50 active:scale-95 transition-all touch-manipulation border border-gray-100">
                    <span class="text-xs font-semibold text-gray-700">{{ variant.name }}</span>
                    <span class="text-sm font-bold text-blue-600">{{ formatPrice(variant.price) }}</span>
                  </button>
                </div>
                <div v-if="getCartItemQtyWithVariants(item.id) > 0" 
                  class="mt-3 bg-blue-500 text-white text-xs font-semibold px-3 py-1.5 rounded-full inline-flex items-center justify-center gap-1 shadow-md">
                  <span>🛒</span>
                  <span>{{ getCartItemQtyWithVariants(item.id) }} món</span>
                </div>
              </div>
            </template>
          </div>
        </div>

        <!-- Cart Summary - Fixed Bottom -->
        <div v-if="cart.length > 0" class="bg-white border-t shadow-lg">
          <div class="px-4 py-3">
            <!-- Cart Items -->
            <div class="max-h-32 overflow-y-auto mb-3 space-y-2">
              <div v-for="(item, idx) in cart" :key="idx" 
                class="flex items-center gap-3 bg-gray-50 rounded-lg p-2">
                <span class="flex-1 text-sm font-medium">
                  {{ item.name }}
                  <span v-if="item.variant_name" class="text-gray-500">({{ item.variant_name }})</span>
                </span>
                <div class="flex items-center gap-2">
                  <button @click="decreaseQty(idx)" 
                    class="w-8 h-8 bg-gray-200 rounded-full text-lg font-bold active:bg-gray-300">
                    −
                  </button>
                  <span class="w-8 text-center font-bold">{{ item.quantity }}</span>
                  <button @click="increaseQty(idx)" 
                    class="w-8 h-8 bg-blue-500 text-white rounded-full text-lg font-bold active:bg-blue-600">
                    +
                  </button>
                </div>
                <button @click="removeFromCart(idx)" class="text-red-500 text-xl">×</button>
              </div>
            </div>
            
            <!-- Total -->
            <div class="flex justify-between items-center">
              <span class="text-gray-600">Tổng cộng</span>
              <span class="text-2xl font-bold text-green-600">{{ formatPrice(cartTotal) }}</span>
            </div>
          </div>
        </div>
      </div>
    </transition>

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
            <div class="grid grid-cols-3 gap-2">
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

    <!-- Customer Name Modal -->
    <transition name="slide-up">
      <div v-if="showCustomerNameModal" class="fixed inset-0 bg-black bg-opacity-50 z-[60] flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6 pb-8">
          <h3 class="text-xl font-bold mb-2 text-gray-900">Tên khách hàng</h3>
          <p class="text-sm text-gray-600 mb-4">Nhập tên khách hàng để dễ dàng quản lý order (tùy chọn)</p>
          
          <!-- Prominent Customer Name Input -->
          <div class="mb-6 p-4 bg-gradient-to-br from-amber-50 via-yellow-50 to-orange-50 rounded-2xl border-2 border-amber-200 shadow-sm">
            <label class="block text-sm font-bold text-amber-900 mb-3 flex items-center gap-2">
              <span class="text-2xl">👤</span>
              <span>Tên khách hàng</span>
            </label>
            <input v-model="customerName" 
              type="text" 
              placeholder="Nhập tên khách hàng..."
              autofocus
              class="w-full px-4 py-4 rounded-xl border-2 border-amber-300 bg-white focus:ring-2 focus:ring-amber-500 focus:border-amber-500 text-lg font-medium placeholder-amber-400 touch-manipulation transition-all shadow-sm">
          </div>

          <!-- Order Summary -->
          <div class="mb-6 p-4 bg-gray-50 rounded-xl">
            <div class="flex justify-between items-center mb-2">
              <span class="text-sm text-gray-600">Số món</span>
              <span class="font-semibold text-gray-900">{{ cart.length }} món</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Tổng tiền</span>
              <span class="text-xl font-bold text-green-600">{{ formatPrice(cartTotal) }}</span>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex gap-3">
            <button @click="showCustomerNameModal = false" 
              class="flex-1 bg-gray-200 text-gray-700 py-4 rounded-xl font-semibold touch-manipulation active:scale-98 transition-all">
              Quay lại
            </button>
            <button @click="finalizeOrder" 
              class="flex-1 bg-blue-500 active:bg-blue-600 text-white py-4 rounded-xl font-semibold shadow-lg touch-manipulation active:scale-98 transition-all">
              Tạo Order
            </button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useOrderStore, cartHelpers } from '../stores/order'
import { useShiftStore } from '../stores/shift'
import { useMenuStore } from '../stores/menu'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
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
const menuStore = useMenuStore()
const authStore = useAuthStore()

// State
const shiftFilter = ref('current') // 'current' or 'other'
const filterStatus = ref('ALL')
const showCreateOrder = ref(false)
const showCustomerNameModal = ref(false) // Modal for customer name
const selectedOrder = ref(null)
const showPayment = ref(false)
const paymentOrder = ref(null)
const paymentAmount = ref(0)
const paymentMethod = ref(PAYMENT_METHOD.CASH)

// Reprint State
const reprintingBill = ref(false)
const reprintingLabel = ref(null)

// Check if user can reprint (only cashier and manager)
const canReprint = computed(() => {
  const role = authStore.user?.role
  return role === 'cashier' || role === 'manager'
})

// Create Order State
const customerName = ref('')
const selectedCategory = ref('all')
const cart = ref([])

// Data
const statuses = STATUS_FILTER_OPTIONS

// Get unique categories from menu items dynamically
const categories = computed(() => {
  const allCategory = { id: 'all', name: 'Tất cả', icon: '📋' }
  
  if (!menuItems.value || menuItems.value.length === 0) {
    return [allCategory]
  }
  
  // Get unique categories from menu items
  const uniqueCategories = [...new Set(menuItems.value.map(item => item.category))]
  
  // Map categories with icons
  const categoryIcons = {
    'Cà phê': '☕',
    'Trà': '🍵',
    'Nước ép': '🧃',
    'Bánh ngọt': '🍰',
    'Món nhẹ': '🍴',
    'Sinh tố': '🥤',
    'Đồ uống khác': '🥛'
  }
  
  const menuCategories = uniqueCategories.map(cat => ({
    id: cat,
    name: cat,
    icon: categoryIcons[cat] || '🍽️'
  }))
  
  return [allCategory, ...menuCategories]
})

const paymentMethods = PAYMENT_METHOD_DISPLAY

// Computed
const loading = computed(() => orderStore.loading)
const orders = computed(() => orderStore.orders)
const menuItems = computed(() => menuStore.items)
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

// Then filter by status
const filteredOrders = computed(() => {
  if (filterStatus.value === 'ALL') return ordersByShift.value
  return ordersByShift.value.filter(o => o.status === filterStatus.value)
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

const filteredMenuItems = computed(() => {
  if (selectedCategory.value === 'all') return menuItems.value
  return menuItems.value.filter(item => item.category === selectedCategory.value)
})

const cartTotal = computed(() => {
  return cart.value.reduce((sum, item) => sum + (item.price * item.quantity), 0)
})

// Helper to check order status
const isStatus = (order, status) => order.status === status
const isAnyStatus = (order, statuses) => statuses.includes(order.status)

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

const viewOrderDetail = (order) => {
  selectedOrder.value = order
}

const startNewOrder = () => {
  cart.value = []
  customerName.value = ''
  selectedCategory.value = 'all'
  showCreateOrder.value = true
}

const cancelCreateOrder = () => {
  if (cart.value.length > 0) {
    if (!confirm('Bạn có chắc muốn hủy order này?')) return
  }
  showCreateOrder.value = false
  cart.value = []
  customerName.value = ''
}

const addToCart = (item, variant = null) => {
  // Use cartHelpers to create cart item
  const cartItem = cartHelpers.createCartItem(item, variant)
  
  // Check if item already exists in cart (considering variant)
  const existing = cart.value.find(i => cartHelpers.isSameCartItem(i, cartItem))
  
  if (existing) {
    existing.quantity++
  } else {
    cart.value.push(cartItem)
  }
}

const getCartItemQty = (itemId) => {
  const item = cart.value.find(i => i.menu_item_id === itemId && !i.variant_id)
  return item ? item.quantity : 0
}

const getCartItemQtyWithVariants = (itemId) => {
  return cart.value
    .filter(i => i.menu_item_id === itemId)
    .reduce((sum, i) => sum + i.quantity, 0)
}

const increaseQty = (index) => {
  cart.value[index].quantity++
}

const decreaseQty = (index) => {
  if (cart.value[index].quantity > 1) {
    cart.value[index].quantity--
  } else {
    removeFromCart(index)
  }
}

const removeFromCart = (index) => {
  cart.value.splice(index, 1)
}

const confirmOrder = () => {
  // Show customer name modal before creating order
  showCustomerNameModal.value = true
}

const finalizeOrder = async () => {
  try {
    const orderData = {
      customer_name: customerName.value || '',
      items: cart.value,
      note: '',
      shift_id: shiftStore.currentShift.id
    }
    await orderStore.createOrder(orderData)
    
    // Close all modals and reset
    showCustomerNameModal.value = false
    showCreateOrder.value = false
    cart.value = []
    customerName.value = ''
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
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

// Lifecycle
onMounted(async () => {
  await Promise.all([
    shiftStore.fetchCurrentShift(),
    orderStore.fetchOrders(),
    menuStore.fetchMenuItems()
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
