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
        <h1 class="text-2xl font-bold text-gray-900 mb-4">⏰ Ca làm việc</h1>
        
        <!-- Shift Tabs (always show for waiter/barista) -->
        <div v-if="!isCashier" class="flex gap-2">
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
            </div>
          </button>
          <button 
            @click="shiftFilter = 'history'"
            :class="[
              'flex-1 py-3 px-4 rounded-xl font-semibold text-sm transition-all touch-manipulation active:scale-98',
              shiftFilter === 'history' 
                ? 'bg-blue-500 text-white shadow-lg' 
                : 'bg-gray-100 text-gray-700 active:bg-gray-200'
            ]">
            <div class="flex items-center justify-center gap-1.5">
              <span>📂</span>
              <span>Ca cũ</span>
              <span class="text-xs opacity-80">({{ oldShiftsCount }})</span>
            </div>
          </button>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-2 pb-24">
      <!-- CASHIER VIEW (NEW COMPONENT) -->
      <CashierShiftView v-if="isCashier" />

      <!-- WAITER/BARISTA VIEW (UNCHANGED) -->
      <div v-else>
      <!-- Tab: Ca hiện tại -->
      <div v-if="shiftFilter === 'current'">
        <!-- Current Shift Info (when shift is open) -->
        <div v-if="currentShift" class="bg-gradient-to-r from-blue-500 to-purple-500 text-white rounded-2xl p-5 mb-4 shadow-xl border border-blue-400">
        <div class="flex justify-between items-start mb-4">
          <div>
            <h3 class="text-2xl font-bold">Ca đang mở</h3>
            <p class="text-blue-100 font-medium">{{ getShiftTypeText(currentShift.type) }}</p>
            <p v-if="currentShift.role_type" class="text-sm text-blue-100 mt-1">
              {{ getRoleTypeText(currentShift.role_type) }}
            </p>
          </div>
          <span class="bg-white text-blue-600 px-4 py-2 rounded-full font-bold text-sm shadow-md">ĐANG MỞ</span>
        </div>
        
        <div class="grid grid-cols-2 gap-3 mb-4">
          <div class="bg-white bg-opacity-20 rounded-xl p-3">
            <p class="text-sm text-blue-100">Bắt đầu</p>
            <p class="font-bold">{{ formatTime(currentShift.started_at) }}</p>
          </div>
          <div class="bg-white bg-opacity-20 rounded-xl p-3">
            <p class="text-sm text-blue-100">Tiền đầu ca</p>
            <p class="font-bold">{{ formatPrice(currentShift.start_cash) }}</p>
          </div>
        </div>

        <!-- Cash Status for Waiter -->
        <div v-if="isWaiter" class="space-y-3 mb-4">
          <!-- Revenue Section -->
          <div class="bg-white bg-opacity-20 rounded-xl p-3">
            <p class="text-xs text-blue-100 mb-2">📊 Doanh thu</p>
            <div class="grid grid-cols-2 gap-2">
              <div>
                <p class="text-xs text-blue-100">💵 Tiền mặt</p>
                <p class="font-bold">{{ formatPrice(currentShift.current_cash || 0) }}</p>
              </div>
              <div>
                <p class="text-xs text-blue-100">💳 Tiền CK</p>
                <p class="font-bold">{{ formatPrice(currentShift.transfer_revenue || 0) }}</p>
              </div>
            </div>
          </div>
          
          <!-- Remaining Section -->
          <div class="bg-white bg-opacity-20 rounded-xl p-3">
            <p class="text-xs text-blue-100 mb-2">💰 Còn lại (chưa bàn giao)</p>
            <div class="grid grid-cols-2 gap-2">
              <div>
                <p class="text-xs text-blue-100">💵 Tiền mặt</p>
                <p class="font-bold text-green-300">{{ formatPrice(currentShift.remaining_cash || 0) }}</p>
              </div>
              <div>
                <p class="text-xs text-blue-100">💳 Tiền CK</p>
                <p class="font-bold text-blue-300">{{ formatPrice(currentShift.remaining_transfer || 0) }}</p>
              </div>
            </div>
          </div>
          
          <!-- Handed Over Section -->
          <div class="bg-white bg-opacity-20 rounded-xl p-3">
            <p class="text-xs text-blue-100 mb-2">✅ Đã bàn giao</p>
            <div class="grid grid-cols-2 gap-2">
              <div>
                <p class="text-xs text-blue-100">💵 Tiền mặt</p>
                <p class="font-bold">{{ formatPrice(currentShift.handed_over_cash || 0) }}</p>
              </div>
              <div>
                <p class="text-xs text-blue-100">💳 Tiền CK</p>
                <p class="font-bold">{{ formatPrice(currentShift.handed_over_transfer || 0) }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Pending Handover Status -->
        <div v-if="isWaiter && pendingHandover" class="bg-yellow-500 bg-opacity-20 rounded-xl p-3 mb-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-yellow-100">🕐 Đang chờ xác nhận bàn giao</p>
              <p class="font-bold">{{ formatPrice(pendingHandover.declared_amount) }}</p>
              <p class="text-xs text-yellow-200">{{ getHandoverTypeText(pendingHandover.handover_type) }}</p>
            </div>
            <button @click="cancelHandover(pendingHandover.id)" 
              class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded-lg text-sm">
              Hủy
            </button>
          </div>
        </div>

        <!-- Action Buttons for Waiter -->
        <div v-if="isWaiter" class="space-y-2.5">
          <!-- Partial Handover Button -->
          <button v-if="(currentShift.remaining_cash > 0 || currentShift.remaining_transfer > 0) && !pendingHandover" 
            @click="showPartialHandoverForm = true"
            class="w-full bg-yellow-500 active:bg-yellow-600 text-white px-4 py-3.5 rounded-xl font-bold shadow-lg touch-manipulation active:scale-98 transition-all">
            💰 Bàn giao một phần
          </button>
          
          <!-- Handover and End Shift Button -->
          <button v-if="(currentShift.remaining_cash > 0 || currentShift.remaining_transfer > 0) && !pendingHandover"
            @click="checkPendingOrdersThenShowHandoverEnd"
            class="w-full bg-orange-500 active:bg-orange-600 text-white px-4 py-3.5 rounded-xl font-bold shadow-lg touch-manipulation active:scale-98 transition-all">
            🏁 Bàn giao và đóng ca
          </button>

          <!-- Regular End Shift Button (only when no remaining cash/transfer) -->
          <button v-if="currentShift.remaining_cash === 0 && currentShift.remaining_transfer === 0 && !pendingHandover"
            @click="checkPendingOrdersThenShowEndShift"
            class="w-full bg-white text-blue-600 active:bg-blue-50 px-4 py-3.5 rounded-xl font-bold shadow-lg touch-manipulation active:scale-98 transition-all">
            Kết thúc ca
          </button>
          
          <!-- Disabled state when pending -->
          <div v-if="pendingHandover" class="w-full bg-gray-300 text-gray-500 px-4 py-3.5 rounded-xl font-bold text-center">
            Chờ cashier xác nhận...
          </div>
        </div>
        
        <!-- Action Buttons for Non-Waiter -->
        <div v-else>
          <button @click="showEndShiftForm = true" 
            class="w-full bg-white text-blue-600 active:bg-blue-50 px-4 py-3.5 rounded-xl font-bold shadow-lg touch-manipulation active:scale-98 transition-all">
            Kết thúc ca
          </button>
        </div>
      </div>

      <!-- Handover History Section for Waiter (in current tab) -->
      <div v-if="isWaiter && handoverHistory.length > 0 && shiftFilter === 'current'" class="bg-white rounded-2xl p-6 shadow-sm mb-4">
        <h3 class="text-xl font-bold mb-4">📋 Lịch sử bàn giao</h3>
        <div class="space-y-3">
          <div v-for="handover in handoverHistory" :key="handover.id" 
            class="border rounded-xl p-4">
            <div class="flex justify-between items-start mb-2">
              <div>
                <!-- Display separate amounts if available -->
                <div v-if="handover.cash_declared_amount > 0 || handover.transfer_declared_amount > 0" class="space-y-1">
                  <p v-if="handover.cash_declared_amount > 0" class="font-bold text-green-600">
                    💵 {{ formatPrice(handover.cash_declared_amount) }}
                  </p>
                  <p v-if="handover.transfer_declared_amount > 0" class="font-bold text-blue-600">
                    💳 {{ formatPrice(handover.transfer_declared_amount) }}
                  </p>
                  <p v-if="handover.cash_declared_amount > 0 && handover.transfer_declared_amount > 0" class="text-sm text-gray-500">
                    Tổng: {{ formatPrice(handover.cash_declared_amount + handover.transfer_declared_amount) }}
                  </p>
                </div>
                <!-- Fallback to old format -->
                <p v-else class="font-bold">{{ formatPrice(handover.declared_amount) }}</p>
                
                <p class="text-sm text-gray-500">{{ formatDate(handover.handover_at) }}</p>
                <p class="text-xs text-blue-600">{{ getHandoverTypeText(handover.handover_type) }}</p>
              </div>
              <span :class="getHandoverStatusClass(handover.status)"
                class="px-3 py-1 rounded-full text-xs font-medium">
                {{ getHandoverStatusText(handover.status) }}
              </span>
            </div>
            <div v-if="handover.waiter_note" class="text-sm text-gray-600 mb-2">
              <strong>Ghi chú:</strong> {{ handover.waiter_note }}
            </div>
            <div v-if="handover.cashier_note" class="text-sm text-green-600">
              <strong>Phản hồi cashier:</strong> {{ handover.cashier_note }}
            </div>
            <!-- Display separate discrepancies -->
            <div v-if="(handover.cash_discrepancy && handover.cash_discrepancy !== 0) || (handover.transfer_discrepancy && handover.transfer_discrepancy !== 0)" 
              class="text-sm mt-2 space-y-1">
              <div v-if="handover.cash_discrepancy && handover.cash_discrepancy !== 0" 
                class="p-2 rounded" 
                :class="handover.cash_discrepancy > 0 ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'">
                <strong>💵 Chênh lệch tiền mặt:</strong> {{ handover.cash_discrepancy > 0 ? '+' : '' }}{{ formatPrice(handover.cash_discrepancy) }}
              </div>
              <div v-if="handover.transfer_discrepancy && handover.transfer_discrepancy !== 0" 
                class="p-2 rounded" 
                :class="handover.transfer_discrepancy > 0 ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'">
                <strong>💳 Chênh lệch tiền CK:</strong> {{ handover.transfer_discrepancy > 0 ? '+' : '' }}{{ formatPrice(handover.transfer_discrepancy) }}
              </div>
            </div>
            <!-- Fallback to old discrepancy format -->
            <div v-else-if="handover.discrepancy && handover.discrepancy !== 0" class="text-sm mt-2 p-2 rounded" 
              :class="handover.discrepancy > 0 ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'">
              <strong>Chênh lệch:</strong> {{ handover.discrepancy > 0 ? '+' : '' }}{{ formatPrice(handover.discrepancy) }}
            </div>
          </div>
        </div>
      </div>

      <!-- Start Shift Form (when no current shift, in current tab) -->
      <div v-if="!currentShift && shiftFilter === 'current'" class="bg-white rounded-2xl p-5 mb-4 shadow-lg border border-gray-200">
        <h3 class="text-xl font-bold mb-4 text-gray-900">Mở ca làm việc</h3>
        <form @submit.prevent="startShift" class="space-y-4">
          <div>
            <label class="block text-sm font-semibold mb-2 text-gray-700">Chọn ca *</label>
            <select v-model="startForm.type" required 
              class="w-full p-3.5 border-2 border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent text-base touch-manipulation">
              <option value="">-- Chọn ca --</option>
              <option value="MORNING">☀️ Ca sáng (7:00 - 12:00)</option>
              <option value="AFTERNOON">🌤️ Ca chiều (12:00 - 18:00)</option>
              <option value="EVENING">🌙 Ca tối (18:00 - 22:00)</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-semibold mb-2 text-gray-700">Tiền đầu ca (VNĐ) *</label>
            <input v-model.number="startForm.start_cash" type="number" min="0" step="1000" required 
              class="w-full p-3.5 border-2 border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent text-base touch-manipulation">
          </div>
          <button type="submit" 
            class="w-full bg-blue-500 active:bg-blue-600 text-white px-4 py-3.5 rounded-xl font-bold shadow-lg touch-manipulation active:scale-98 transition-all">
            Mở ca
          </button>
        </form>
      </div>
      </div>
      <!-- End of Tab: Ca hiện tại -->

      <!-- Tab: Ca cũ (Shift History) -->
      <div v-if="shiftFilter === 'history'" class="bg-white rounded-2xl p-5 shadow-lg border border-gray-200">
        <h3 class="text-xl font-bold mb-4 text-gray-900">Lịch sử ca làm việc</h3>
        
        <div v-if="loading" class="text-center py-16">
          <div class="inline-block animate-spin rounded-full h-10 w-10 border-4 border-blue-500 border-t-transparent"></div>
          <p class="text-gray-500 text-sm mt-3">Đang tải...</p>
        </div>
        
        <div v-else-if="filteredShifts.length === 0" class="text-center py-16">
          <div class="text-6xl mb-3">📭</div>
          <p class="text-gray-500 font-medium">Chưa có ca làm việc cũ</p>
        </div>
        
        <div v-else class="space-y-3">
          <div v-for="shift in filteredShifts" :key="shift.id" 
            class="border-2 border-gray-100 rounded-xl p-4 active:scale-[0.98] active:shadow-md transition-all touch-manipulation">
            <div class="flex justify-between items-start mb-3">
              <div>
                <h4 class="font-bold text-lg">{{ getShiftTypeText(shift.type) }}</h4>
                <p class="text-sm text-gray-500">{{ formatDate(shift.started_at) }}</p>
                <p v-if="shift.role_type" class="text-xs text-blue-600 font-medium mt-1">
                  {{ getRoleTypeText(shift.role_type) }}
                </p>
              </div>
              <span :class="shift.status === SHIFT_STATUS.OPEN ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'"
                class="px-3 py-1 rounded-full text-xs font-medium">
                {{ shift.status === SHIFT_STATUS.OPEN ? 'Đang mở' : 'Đã đóng' }}
              </span>
            </div>

            <div class="grid grid-cols-2 gap-3 text-sm">
              <div class="bg-gray-50 rounded-lg p-3">
                <p class="text-gray-500 text-xs">Tiền đầu ca</p>
                <p class="font-bold">{{ formatPrice(shift.start_cash) }}</p>
              </div>
              <div v-if="shift.status === 'CLOSED'" class="bg-gray-50 rounded-lg p-3">
                <p class="text-gray-500 text-xs">Tiền cuối ca</p>
                <p class="font-bold">{{ formatPrice(shift.end_cash) }}</p>
              </div>
              <div v-if="shift.status === 'CLOSED'" class="bg-green-50 rounded-lg p-3">
                <p class="text-gray-500 text-xs">Doanh thu</p>
                <p class="font-bold text-green-600">{{ formatPrice(shift.total_revenue) }}</p>
              </div>
              <div v-if="shift.status === 'CLOSED'" class="bg-blue-50 rounded-lg p-3">
                <p class="text-gray-500 text-xs">Số order</p>
                <p class="font-bold text-blue-600">{{ shift.total_orders }}</p>
              </div>
            </div>

            <button v-if="isCashier && shift.status === 'OPEN' && shift.id !== currentShift?.id" 
              @click="showCloseShiftForm(shift)"
              class="mt-3 w-full bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-xl text-sm font-medium active:scale-95 transition-transform">
              Chốt ca
            </button>
          </div>
        </div>
      </div>
      </div>
      <!-- End of Waiter/Barista View -->
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />

    <!-- Pending Orders Warning Modal -->
    <transition name="slide-up">
      <div v-if="showPendingOrdersWarning" class="fixed inset-0 bg-black bg-opacity-60 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6 max-h-[85vh] overflow-y-auto">
          <div class="flex items-center gap-3 mb-4">
            <span class="text-3xl">⚠️</span>
            <div>
              <h3 class="text-xl font-bold text-red-600">Còn order chưa hoàn thành</h3>
              <p class="text-sm text-gray-500">Đóng ca sẽ hủy tất cả các order này</p>
            </div>
          </div>

          <div class="bg-red-50 border border-red-200 rounded-xl p-4 mb-5">
            <p class="text-sm font-semibold text-red-700 mb-3">
              {{ pendingOrdersList.length }} order đang pending sẽ bị hủy:
            </p>
            <div class="space-y-2 max-h-48 overflow-y-auto">
              <div v-for="o in pendingOrdersList" :key="o.id"
                class="flex justify-between items-center bg-white rounded-lg p-3 border border-red-100">
                <div>
                  <p class="font-semibold text-gray-800">#{{ o.order_number }}</p>
                  <p v-if="o.customer_name" class="text-xs text-gray-500">{{ o.customer_name }}</p>
                  <span class="text-xs px-2 py-0.5 rounded-full font-medium"
                    :class="{
                      'bg-gray-100 text-gray-700': o.status === 'CREATED',
                      'bg-blue-100 text-blue-700': o.status === 'PAID',
                      'bg-yellow-100 text-yellow-700': o.status === 'QUEUED',
                      'bg-orange-100 text-orange-700': o.status === 'IN_PROGRESS',
                      'bg-green-100 text-green-700': o.status === 'READY',
                    }">
                    {{ getOrderStatusText(o.status) }}
                  </span>
                </div>
                <p class="font-bold text-gray-800">{{ formatPrice(o.total) }}</p>
              </div>
            </div>
          </div>

          <p class="text-sm text-gray-600 mb-5">
            Bạn có chắc chắn muốn <strong class="text-red-600">hủy {{ pendingOrdersList.length }} order</strong> và tiếp tục đóng ca?
          </p>

          <div class="flex gap-3">
            <button @click="showPendingOrdersWarning = false; pendingOrdersAction = null"
              class="flex-1 bg-gray-200 text-gray-700 px-4 py-3.5 rounded-xl font-bold">
              Quay lại
            </button>
            <button @click="confirmPendingOrdersAndProceed"
              class="flex-1 bg-red-500 active:bg-red-600 text-white px-4 py-3.5 rounded-xl font-bold">
              Xác nhận hủy & đóng ca
            </button>
          </div>
        </div>
      </div>
    </transition>

    <!-- End Shift Modal -->
    <transition name="slide-up">
      <div v-if="showEndShiftForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6">
          <h3 class="text-xl font-bold mb-4">Kết thúc ca</h3>
          <form @submit.prevent="endShift" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Tiền cuối ca (VNĐ) *</label>
              <input v-model.number="endForm.end_cash" type="number" min="0" step="1000" required 
                class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-blue-500">
            </div>
            <div class="bg-blue-50 p-4 rounded-xl">
              <p class="text-sm text-gray-600">Tiền đầu ca</p>
              <p class="font-bold text-2xl text-blue-600">{{ formatPrice(currentShift?.start_cash) }}</p>
            </div>
            <div class="flex gap-2">
              <button type="button" @click="showEndShiftForm = false" 
                class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
                Hủy
              </button>
              <button type="submit" 
                class="flex-1 bg-blue-500 hover:bg-blue-600 text-white px-4 py-3 rounded-xl font-medium">
                Kết thúc
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>

    <!-- Close Shift Modal -->
    <transition name="slide-up">
      <div v-if="showCloseForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6">
          <h3 class="text-xl font-bold mb-4">Chốt ca</h3>
          <p class="text-sm text-gray-600 mb-4">Chốt ca sẽ khóa tất cả orders trong ca này</p>
          <form @submit.prevent="closeShift" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Tiền cuối ca (VNĐ) *</label>
              <input v-model.number="closeForm.end_cash" type="number" min="0" step="1000" required 
                class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-red-500">
            </div>
            <div class="flex gap-2">
              <button type="button" @click="showCloseForm = false" 
                class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
                Hủy
              </button>
              <button type="submit" 
                class="flex-1 bg-red-500 hover:bg-red-600 text-white px-4 py-3 rounded-xl font-medium">
                Chốt ca
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>

    <!-- Partial Handover Modal -->
    <transition name="slide-up">
      <div v-if="showPartialHandoverForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6 max-h-[90vh] overflow-y-auto">
          <h3 class="text-xl font-bold mb-4">💰 Bàn giao tiền</h3>
          
          <!-- Current Balance Info -->
          <div class="bg-blue-50 p-4 rounded-xl mb-4">
            <p class="text-sm font-medium text-gray-700 mb-3">💰 Số dư hiện tại</p>
            <div class="space-y-2">
              <div class="flex justify-between items-center">
                <span class="text-sm text-gray-600">💵 Tiền mặt còn lại</span>
                <span class="font-bold text-lg text-green-600">{{ formatPrice(currentShift?.remaining_cash || 0) }}</span>
              </div>
              <div class="flex justify-between items-center">
                <span class="text-sm text-gray-600">💳 Tiền CK còn lại</span>
                <span class="font-bold text-lg text-blue-600">{{ formatPrice(currentShift?.remaining_transfer || 0) }}</span>
              </div>
            </div>
          </div>
          
          <form @submit.prevent="createPartialHandover" class="space-y-4">
            <!-- Handover Type Selection -->
            <div>
              <label class="block text-sm font-medium mb-2">Chọn loại bàn giao *</label>
              <div class="grid grid-cols-3 gap-2">
                <button type="button" 
                  @click="partialHandoverForm.handover_type = 'cash'"
                  :class="partialHandoverForm.handover_type === 'cash' ? 'bg-green-500 text-white' : 'bg-gray-100 text-gray-700'"
                  class="px-4 py-3 rounded-xl font-medium transition-colors">
                  💵 Tiền mặt
                </button>
                <button type="button" 
                  @click="partialHandoverForm.handover_type = 'transfer'"
                  :class="partialHandoverForm.handover_type === 'transfer' ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-700'"
                  class="px-4 py-3 rounded-xl font-medium transition-colors">
                  💳 Tiền CK
                </button>
                <button type="button" 
                  @click="partialHandoverForm.handover_type = 'both'"
                  :class="partialHandoverForm.handover_type === 'both' ? 'bg-purple-500 text-white' : 'bg-gray-100 text-gray-700'"
                  class="px-4 py-3 rounded-xl font-medium transition-colors">
                  💰 Cả hai
                </button>
              </div>
            </div>

            <!-- Cash Amount Input -->
            <div v-if="partialHandoverForm.handover_type === 'cash' || partialHandoverForm.handover_type === 'both'">
              <label class="block text-sm font-medium mb-2">💵 Số tiền mặt bàn giao (VNĐ) *</label>
              <input v-model.number="partialHandoverForm.cash_amount" 
                type="number" 
                :max="currentShift?.remaining_cash || 0"
                min="0" 
                step="1000" 
                :required="partialHandoverForm.handover_type === 'cash' || partialHandoverForm.handover_type === 'both'"
                class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-green-500">
              <p class="text-xs text-gray-500 mt-1">Tối đa: {{ formatPrice(currentShift?.remaining_cash || 0) }}</p>
            </div>

            <!-- Transfer Amount Input -->
            <div v-if="partialHandoverForm.handover_type === 'transfer' || partialHandoverForm.handover_type === 'both'">
              <label class="block text-sm font-medium mb-2">💳 Số tiền CK bàn giao (VNĐ) *</label>
              <input v-model.number="partialHandoverForm.transfer_amount" 
                type="number" 
                :max="currentShift?.remaining_transfer || 0"
                min="0" 
                step="1000" 
                :required="partialHandoverForm.handover_type === 'transfer' || partialHandoverForm.handover_type === 'both'"
                class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-blue-500">
              <p class="text-xs text-gray-500 mt-1">Tối đa: {{ formatPrice(currentShift?.remaining_transfer || 0) }}</p>
            </div>

            <!-- Total Display -->
            <div v-if="partialHandoverForm.handover_type === 'both'" class="bg-purple-50 p-4 rounded-xl">
              <div class="flex justify-between items-center">
                <span class="text-sm text-gray-600">Tổng bàn giao</span>
                <span class="font-bold text-2xl text-purple-600">
                  {{ formatPrice((partialHandoverForm.cash_amount || 0) + (partialHandoverForm.transfer_amount || 0)) }}
                </span>
              </div>
            </div>
            
            <!-- Note -->
            <div>
              <label class="block text-sm font-medium mb-2">Ghi chú (tùy chọn)</label>
              <textarea v-model="partialHandoverForm.waiter_note" 
                rows="3" 
                class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-yellow-500"
                placeholder="Ghi chú về việc bàn giao..."></textarea>
            </div>
            
            <!-- Action Buttons -->
            <div class="flex gap-2">
              <button type="button" @click="showPartialHandoverForm = false" 
                class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
                Hủy
              </button>
              <button type="submit" 
                class="flex-1 bg-yellow-500 hover:bg-yellow-600 text-white px-4 py-3 rounded-xl font-medium">
                Bàn giao
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>

    <!-- Handover and End Shift Modal -->
    <transition name="slide-up">
      <div v-if="showHandoverEndShiftForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6">
          <h3 class="text-xl font-bold mb-4">🏁 Bàn giao toàn bộ và đóng ca</h3>
          
          <!-- Warning Notice -->
          <div class="bg-orange-50 border-l-4 border-orange-400 p-4 mb-4">
            <div class="flex">
              <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-orange-400" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                </svg>
              </div>
              <div class="ml-3">
                <p class="text-sm text-orange-700">
                  <strong>Lưu ý:</strong> Thao tác này sẽ bàn giao toàn bộ tiền mặt và tiền CK còn lại, sau đó tự động đóng ca khi cashier xác nhận.
                </p>
              </div>
            </div>
          </div>
          
          <!-- Money Summary -->
          <div class="bg-orange-50 p-4 rounded-xl mb-4">
            <p class="text-sm font-medium text-gray-700 mb-3">💰 Số tiền sẽ bàn giao</p>
            <div class="space-y-2">
              <div v-if="(currentShift?.remaining_cash || 0) > 0" class="flex justify-between items-center">
                <span class="text-sm text-gray-600">💵 Tiền mặt</span>
                <span class="font-bold text-lg text-green-600">{{ formatPrice(currentShift?.remaining_cash || 0) }}</span>
              </div>
              <div v-if="(currentShift?.remaining_transfer || 0) > 0" class="flex justify-between items-center">
                <span class="text-sm text-gray-600">💳 Tiền CK</span>
                <span class="font-bold text-lg text-blue-600">{{ formatPrice(currentShift?.remaining_transfer || 0) }}</span>
              </div>
              <div v-if="(currentShift?.remaining_cash || 0) > 0 && (currentShift?.remaining_transfer || 0) > 0" 
                class="flex justify-between items-center pt-2 border-t border-orange-200">
                <span class="text-sm font-medium text-gray-700">Tổng cộng</span>
                <span class="font-bold text-xl text-orange-600">
                  {{ formatPrice((currentShift?.remaining_cash || 0) + (currentShift?.remaining_transfer || 0)) }}
                </span>
              </div>
              <div class="flex justify-between items-center text-sm pt-2 border-t border-orange-200">
                <span class="text-gray-500">Tiền cuối ca</span>
                <span class="font-medium">{{ formatPrice(handoverEndShiftForm.end_cash) }}</span>
              </div>
            </div>
          </div>
          
          <form @submit.prevent="createHandoverAndEndShift" class="space-y-4">
            <!-- End Cash Input -->
            <div>
              <label class="block text-sm font-medium mb-2">Tiền cuối ca (VNĐ) *</label>
              <input v-model.number="handoverEndShiftForm.end_cash" 
                type="number" 
                min="0" 
                step="1000" 
                required 
                class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-orange-500">
              <p class="text-xs text-gray-500 mt-1">Tiền còn lại sau khi bàn giao (thường là 0)</p>
            </div>
            
            <!-- Note -->
            <div>
              <label class="block text-sm font-medium mb-2">Ghi chú (tùy chọn)</label>
              <textarea v-model="handoverEndShiftForm.waiter_note" 
                rows="3" 
                class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-orange-500"
                placeholder="Ghi chú về việc bàn giao và đóng ca..."></textarea>
            </div>
            
            <!-- Action Buttons -->
            <div class="flex gap-2">
              <button type="button" @click="showHandoverEndShiftForm = false" 
                class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
                Hủy
              </button>
              <button type="submit" 
                class="flex-1 bg-orange-500 hover:bg-orange-600 text-white px-4 py-3 rounded-xl font-medium">
                Bàn giao và đóng ca
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useShiftStore } from '../stores/shift'
import { useAuthStore } from '../stores/auth'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import CashierShiftView from '../components/CashierShiftView.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'
import { USER_ROLES } from '../constants/user'
import { SHIFT_STATUS } from '../constants/shift'

const shiftStore = useShiftStore()
const authStore = useAuthStore()

// State
const shiftFilter = ref('current') // 'current' or 'history'
const showEndShiftForm = ref(false)
const showCloseForm = ref(false)
const selectedShift = ref(null)

// Handover refs
const showPartialHandoverForm = ref(false)
const showHandoverEndShiftForm = ref(false)
const pendingHandover = ref(null)
const handoverHistory = ref([])

// Pending orders warning refs
const showPendingOrdersWarning = ref(false)
const pendingOrdersList = ref([])
const pendingOrdersAction = ref(null) // 'endShift' | 'handoverEnd'

const startForm = ref({
  type: '',
  start_cash: 0
})

const endForm = ref({
  end_cash: 0
})

const closeForm = ref({
  end_cash: 0
})

const partialHandoverForm = ref({
  handover_type: 'both', // 'cash' | 'transfer' | 'both'
  cash_amount: 0,
  transfer_amount: 0,
  waiter_note: ''
})

const handoverEndShiftForm = ref({
  end_cash: 0,
  waiter_note: ''
})

const loading = computed(() => shiftStore.loading)
const currentShift = computed(() => shiftStore.currentShift)
const shifts = computed(() => shiftStore.shifts)
const isCashier = computed(() => authStore.user?.role === USER_ROLES.CASHIER || authStore.user?.role === USER_ROLES.MANAGER)
const isWaiter = computed(() => authStore.user?.role === USER_ROLES.WAITER)

// Filter shifts - exclude current shift from history
const filteredShifts = computed(() => {
  // When viewing history, show all shifts except current one
  if (shiftFilter.value === 'history') {
    if (!currentShift.value) {
      // No current shift, show all shifts
      return shifts.value
    }
    // Exclude the current shift from history
    return shifts.value.filter(s => s.id !== currentShift.value.id)
  }
  
  // When viewing current tab, don't show any shifts in the list
  // (current shift is shown separately above)
  return []
})

// Count old shifts (excluding current)
const oldShiftsCount = computed(() => {
  if (!currentShift.value) return shifts.value.length
  return shifts.value.filter(s => s.id !== currentShift.value.id).length
})

// Refresh data function
const refreshData = async () => {
  await shiftStore.fetchCurrentShift()
  if (isWaiter.value && currentShift.value) {
    await fetchHandoverData()
  }
  if (isCashier.value) {
    await shiftStore.fetchAllShifts()
  } else {
    await shiftStore.fetchMyShifts()
  }
}

// Pull to refresh
const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)

onMounted(async () => {
  await refreshData()
})

const fetchHandoverData = async () => {
  if (!currentShift.value?.id) return
  try {
    pendingHandover.value = await shiftStore.getPendingHandover(currentShift.value.id)
    console.log('Pending handover:', pendingHandover.value)
    handoverHistory.value = await shiftStore.getHandoverHistory(currentShift.value.id)
  } catch (error) {
    console.error('Error fetching handover data:', error)
  }
}

const startShift = async () => {
  try {
    await shiftStore.startShift(startForm.value)
    startForm.value = { type: '', start_cash: 0 }
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

const endShift = async () => {
  try {
    await shiftStore.endShift(currentShift.value.id, endForm.value.end_cash)
    showEndShiftForm.value = false
    endForm.value = { end_cash: 0 }
    await shiftStore.fetchMyShifts()
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

const showCloseShiftForm = (shift) => {
  selectedShift.value = shift
  showCloseForm.value = true
}

const closeShift = async () => {
  try {
    await shiftStore.closeShift(selectedShift.value.id, closeForm.value.end_cash)
    showCloseForm.value = false
    selectedShift.value = null
    closeForm.value = { end_cash: 0 }
    await shiftStore.fetchAllShifts()
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

// Handover functions
const createPartialHandover = async () => {
  try {
    const handoverData = {
      handover_type: 'PARTIAL',
      waiter_note: partialHandoverForm.value.waiter_note
    }

    // Add amounts based on handover type
    if (partialHandoverForm.value.handover_type === 'cash') {
      handoverData.cash_amount = partialHandoverForm.value.cash_amount
      handoverData.transfer_amount = 0
    } else if (partialHandoverForm.value.handover_type === 'transfer') {
      handoverData.cash_amount = 0
      handoverData.transfer_amount = partialHandoverForm.value.transfer_amount
    } else { // both
      handoverData.cash_amount = partialHandoverForm.value.cash_amount || 0
      handoverData.transfer_amount = partialHandoverForm.value.transfer_amount || 0
    }

    // Validate at least one amount
    if (handoverData.cash_amount === 0 && handoverData.transfer_amount === 0) {
      alert('Vui lòng nhập ít nhất một loại tiền để bàn giao')
      return
    }
    
    await shiftStore.createCashHandover(currentShift.value.id, handoverData)
    showPartialHandoverForm.value = false
    partialHandoverForm.value = { 
      handover_type: 'both',
      cash_amount: 0,
      transfer_amount: 0,
      waiter_note: '' 
    }
    
    // Refresh data
    await shiftStore.fetchCurrentShift()
    await fetchHandoverData()
    
    alert('Đã gửi yêu cầu bàn giao. Chờ thu ngân xác nhận.')
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

const createHandoverAndEndShift = async () => {
  try {
    const handoverData = {
      cash_amount: currentShift.value?.remaining_cash || 0,
      transfer_amount: currentShift.value?.remaining_transfer || 0,
      handover_type: 'END_SHIFT',
      waiter_note: handoverEndShiftForm.value.waiter_note,
      end_cash: handoverEndShiftForm.value.end_cash
    }
    
    // Validate at least one amount
    if (handoverData.cash_amount === 0 && handoverData.transfer_amount === 0) {
      alert('Không có tiền để bàn giao')
      return
    }
    
    await shiftStore.createCashHandover(currentShift.value.id, handoverData)
    showHandoverEndShiftForm.value = false
    handoverEndShiftForm.value = { end_cash: 0, waiter_note: '' }
    
    // Refresh data
    await shiftStore.fetchCurrentShift()
    await fetchHandoverData()
    
    alert('Đã gửi yêu cầu bàn giao toàn bộ và đóng ca. Chờ thu ngân xác nhận.')
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

const cancelHandover = async (handoverId) => {
  // Validate handover ID
  if (!handoverId) {
    alert('Lỗi: Không tìm thấy ID bàn giao')
    console.error('Invalid handover ID:', handoverId)
    return
  }
  
  if (confirm('Bạn có chắc muốn hủy yêu cầu bàn giao này?')) {
    try {
      await shiftStore.cancelHandover(handoverId)
      await fetchHandoverData()
      alert('Đã hủy yêu cầu bàn giao!')
    } catch (error) {
      const errorMsg = error.response?.data?.error || error.message
      alert('Lỗi: ' + errorMsg)
      console.error('Cancel handover error:', error)
    }
  }
}

// Pending orders check helpers
const checkPendingOrdersThenShowEndShift = async () => {
  await checkPendingOrdersAndShow('endShift')
}

const checkPendingOrdersThenShowHandoverEnd = async () => {
  await checkPendingOrdersAndShow('handoverEnd')
}

const checkPendingOrdersAndShow = async (action) => {
  if (!currentShift.value?.id) return
  try {
    const result = await shiftStore.getPendingOrders(currentShift.value.id)
    if (result.count > 0) {
      pendingOrdersList.value = result.pending_orders
      pendingOrdersAction.value = action
      showPendingOrdersWarning.value = true
    } else {
      // No pending orders - proceed directly
      if (action === 'endShift') showEndShiftForm.value = true
      else showHandoverEndShiftForm.value = true
    }
  } catch (error) {
    // On error, just proceed without the check
    if (action === 'endShift') showEndShiftForm.value = true
    else showHandoverEndShiftForm.value = true
  }
}

const confirmPendingOrdersAndProceed = () => {
  showPendingOrdersWarning.value = false
  if (pendingOrdersAction.value === 'endShift') showEndShiftForm.value = true
  else if (pendingOrdersAction.value === 'handoverEnd') showHandoverEndShiftForm.value = true
  pendingOrdersAction.value = null
  pendingOrdersList.value = []
}

const getOrderStatusText = (status) => {
  const statuses = {
    CREATED: 'Chưa thanh toán',
    PAID: 'Đã thanh toán',
    QUEUED: 'Chờ pha chế',
    IN_PROGRESS: 'Đang pha chế',
    READY: 'Sẵn sàng giao'
  }
  return statuses[status] || status
}

const getShiftTypeText = (type) => {
  const types = {
    MORNING: '☀️ Ca sáng',
    AFTERNOON: '🌤️ Ca chiều',
    EVENING: '🌙 Ca tối'
  }
  return types[type] || type
}

const getRoleTypeText = (roleType) => {
  const roles = {
    waiter: '👨‍💼 Phục vụ',
    barista: '🍹 Pha chế',
    cashier: '💰 Thu ngân'
  }
  return roles[roleType] || roleType
}

const getHandoverTypeText = (type) => {
  const types = {
    'PARTIAL': 'Một phần',
    'FULL': 'Toàn bộ',
    'END_SHIFT': 'Toàn bộ + Đóng ca'
  }
  return types[type] || type
}

const getHandoverStatusText = (status) => {
  const statuses = {
    'PENDING': 'Chờ xác nhận',
    'CONFIRMED': 'Đã xác nhận',
    'REJECTED': 'Đã từ chối',
    'DISCREPANCY': 'Có chênh lệch'
  }
  return statuses[status] || status
}

const getHandoverStatusClass = (status) => {
  const classes = {
    'PENDING': 'bg-yellow-100 text-yellow-800',
    'CONFIRMED': 'bg-green-100 text-green-800',
    'REJECTED': 'bg-red-100 text-red-800',
    'DISCREPANCY': 'bg-orange-100 text-orange-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const formatPrice = (price) => {
  return new Intl.NumberFormat('vi-VN', { 
    style: 'currency', 
    currency: 'VND',
    maximumFractionDigits: 0
  }).format(price)
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('vi-VN')
}

const formatTime = (date) => {
  return new Date(date).toLocaleTimeString('vi-VN', { 
    hour: '2-digit', 
    minute: '2-digit' 
  })
}
</script>

<style scoped>
.touch-manipulation {
  touch-action: manipulation;
  -webkit-tap-highlight-color: transparent;
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
</style>
