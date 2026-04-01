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
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-800">👋 Xin chào</h1>
            <p class="text-sm text-gray-600">{{ user?.name }}</p>
          </div>
          <div class="text-right">
            <p class="text-xs text-gray-500">{{ currentDate }}</p>
            <p class="text-sm font-medium text-gray-700">{{ currentTime }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24">
      <!-- Manager Dashboard (No Shift Concept) -->
      <div v-if="user?.role === USER_ROLES.MANAGER">
        <!-- Welcome Card -->
        <div class="bg-gradient-to-r from-blue-500 to-purple-500 rounded-2xl p-6 text-white shadow-lg mb-6">
          <h2 class="text-2xl font-bold mb-2">🎯 Quản lý hệ thống</h2>
          <p class="text-sm opacity-90">Truy cập nhanh các chức năng quản lý</p>
        </div>

        <!-- Quick Stats for Manager -->
        <div class="grid grid-cols-2 gap-3 mb-6">
          <div class="bg-white rounded-2xl p-4 shadow-sm">
            <div class="text-3xl mb-2">📋</div>
            <div class="text-2xl font-bold text-gray-800">{{ todayOrders }}</div>
            <div class="text-xs text-gray-500">Orders hôm nay</div>
          </div>
          <div class="bg-white rounded-2xl p-4 shadow-sm">
            <div class="text-3xl mb-2">💰</div>
            <div class="text-lg font-bold text-green-600">{{ formatPrice(todayRevenue) }}</div>
            <div class="text-xs text-gray-500">Doanh thu hôm nay</div>
          </div>
        </div>

        <!-- 📊 Báo cáo & Phân tích -->
        <div class="mb-6">
          <h2 class="text-lg font-bold text-gray-800 mb-3 flex items-center gap-2">
            <span>📊</span>
            <span>Báo cáo & Phân tích</span>
          </h2>
          <div class="grid grid-cols-2 gap-3">
            <button @click="$router.push('/manager/profit-analysis')" 
              class="bg-gradient-to-br from-indigo-500 to-purple-600 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">📈</div>
              <div class="font-bold text-base">Phân tích lợi nhuận</div>
              <div class="text-xs opacity-80 mt-1">Theo category & vận hành</div>
            </button>
            <button @click="$router.push('/manager/menu-costs')" 
              class="bg-gradient-to-br from-orange-500 to-red-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">💰</div>
              <div class="font-bold text-base">Chi phí món</div>
              <div class="text-xs opacity-80 mt-1">Giá vốn & lợi nhuận</div>
            </button>
            <button @click="$router.push('/manager/cost-analysis')" 
              class="bg-gradient-to-br from-teal-500 to-green-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">📊</div>
              <div class="font-bold text-base">Chi phí theo size</div>
              <div class="text-xs opacity-80 mt-1">Phân tích variants</div>
            </button>
            <button @click="$router.push('/cashier/reports')" 
              class="bg-gradient-to-br from-blue-500 to-cyan-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">📊</div>
              <div class="font-bold text-base">Báo cáo thu ngân</div>
              <div class="text-xs opacity-80 mt-1">Doanh thu & thanh toán</div>
            </button>
            <button @click="$router.push('/manager/shifts')" 
              class="bg-gradient-to-br from-purple-500 to-pink-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">⏰</div>
              <div class="font-bold text-base">Quản lý ca</div>
              <div class="text-xs opacity-80 mt-1">Theo dõi ca làm việc</div>
            </button>
          </div>
        </div>

        <!-- 🍽️ Menu & Nguyên liệu -->
        <div class="mb-6">
          <h2 class="text-lg font-bold text-gray-800 mb-3 flex items-center gap-2">
            <span>🍽️</span>
            <span>Menu & Nguyên liệu</span>
          </h2>
          <div class="grid grid-cols-2 gap-3">
            <button @click="$router.push('/menu')" 
              class="bg-gradient-to-br from-amber-500 to-orange-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">🍽️</div>
              <div class="font-bold text-base">Menu</div>
              <div class="text-xs opacity-80 mt-1">Quản lý thực đơn</div>
            </button>
            <button @click="$router.push('/ingredients')" 
              class="bg-gradient-to-br from-green-500 to-emerald-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">🥬</div>
              <div class="font-bold text-base">Nguyên liệu</div>
              <div class="text-xs opacity-80 mt-1">Kho & giá vốn</div>
            </button>
            <button @click="$router.push('/batch')"
              class="bg-gradient-to-br from-purple-500 to-indigo-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">🧪</div>
              <div class="font-bold text-base">Batch</div>
              <div class="text-xs opacity-80 mt-1">Nguyên liệu chế biến</div>
            </button>
            <button @click="$router.push('/manager/inventory-stats')"
              class="bg-gradient-to-br from-teal-500 to-cyan-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">📦</div>
              <div class="font-bold text-base">Tồn kho</div>
              <div class="text-xs opacity-80 mt-1">Thống kê tồn kho</div>
            </button>
          </div>
        </div>

        <!-- 💸 Chi phí & Tài sản -->
        <div class="mb-6">
          <h2 class="text-lg font-bold text-gray-800 mb-3 flex items-center gap-2">
            <span>💸</span>
            <span>Chi phí & Tài sản</span>
          </h2>
          <div class="grid grid-cols-2 gap-3">
            <button @click="$router.push('/manager/fund')" 
              class="bg-gradient-to-br from-orange-500 to-yellow-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">💰</div>
              <div class="font-bold text-base">Quỹ tiền</div>
              <div class="text-xs opacity-80 mt-1">Quản lý quỹ tiền mặt</div>
            </button>
            <button @click="$router.push('/expenses')" 
              class="bg-gradient-to-br from-pink-500 to-rose-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">💸</div>
              <div class="font-bold text-base">Chi phí</div>
              <div class="text-xs opacity-80 mt-1">Theo dõi chi tiêu</div>
            </button>
            <button @click="$router.push('/facilities')"
              class="bg-gradient-to-br from-cyan-500 to-blue-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">🏢</div>
              <div class="font-bold text-base">Cơ sở vật chất</div>
              <div class="text-xs opacity-80 mt-1">Tài sản & bảo trì</div>
            </button>
            <button @click="$router.push('/operating-expenses')"
              class="bg-gradient-to-br from-orange-500 to-red-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">⚡</div>
              <div class="font-bold text-base">Chi phí vận hành</div>
              <div class="text-xs opacity-80 mt-1">Trừ từ quỹ vận hành</div>
            </button>
          </div>
        </div>

        <!-- 👥 Nhân sự & Hệ thống -->
        <div class="mb-6">
          <h2 class="text-lg font-bold text-gray-800 mb-3 flex items-center gap-2">
            <span>👥</span>
            <span>Nhân sự & Hệ thống</span>
          </h2>
          <div class="grid grid-cols-2 gap-3">
            <button @click="$router.push('/users')"
              class="bg-gradient-to-br from-violet-500 to-purple-600 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">👥</div>
              <div class="font-bold text-base">Nhân viên</div>
              <div class="text-xs opacity-80 mt-1">Quản lý tài khoản</div>
            </button>
            <button @click="$router.push('/manager/schedule')"
              class="bg-gradient-to-br from-indigo-500 to-purple-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">📅</div>
              <div class="font-bold text-base">Lịch ca</div>
              <div class="text-xs opacity-80 mt-1">Đăng ký ca làm việc</div>
            </button>
            <button @click="$router.push('/print-management')" 
              class="bg-gradient-to-br from-teal-500 to-cyan-600 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">🖨️</div>
              <div class="font-bold text-base">In ấn</div>
              <div class="text-xs opacity-80 mt-1">Máy in & templates</div>
            </button>
            <button @click="$router.push('/settings')" 
              class="bg-gradient-to-br from-gray-500 to-gray-600 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
              <div class="text-4xl mb-2">⚙️</div>
              <div class="font-bold text-base">Cài đặt</div>
              <div class="text-xs opacity-80 mt-1">Cấu hình hệ thống</div>
            </button>
          </div>
        </div>
      </div>

      <!-- Non-Manager Dashboard (With Shift Concept) -->
      <div v-else>
        <!-- Shift Status (for Waiter and Barista only, not Cashier) -->
        <div v-if="!isCashier">
          <div v-if="hasOpenShift" class="bg-gradient-to-r from-green-500 to-emerald-500 rounded-2xl p-4 text-white shadow-lg mb-4">
            <div class="flex items-center justify-between mb-2">
              <div class="flex items-center gap-2">
                <span class="text-2xl">✅</span>
                <span class="font-bold">Ca đang mở</span>
              </div>
              <span class="text-sm opacity-90">{{ shiftDuration }}</span>
            </div>
            <p class="text-sm opacity-90">Bắt đầu: {{ formatTime(currentShift?.started_at) }}</p>
          </div>
          <div v-else class="bg-gradient-to-r from-orange-500 to-red-500 rounded-2xl p-4 text-white shadow-lg mb-4">
            <div class="flex items-center justify-between mb-2">
              <div class="flex items-center gap-2">
                <span class="text-2xl">⚠️</span>
                <span class="font-bold">Chưa mở ca</span>
              </div>
            </div>
            <button @click="$router.push('/shifts')" 
              class="mt-2 bg-white text-orange-600 px-4 py-2 rounded-lg font-medium text-sm">
              Mở ca ngay
            </button>
          </div>
        </div>

        <!-- Barista Dashboard -->
        <div v-if="isBarista">
          <!-- Current Shift Info -->
          <div v-if="hasOpenShift" class="mb-4 bg-gradient-to-r from-blue-500 to-purple-500 text-white rounded-2xl p-4 shadow-lg">
            <div class="flex items-center justify-between mb-2">
              <div>
                <h3 class="font-bold text-lg">Ca làm việc</h3>
                <p class="text-sm opacity-90">{{ getShiftTypeText(currentShift.type) }}</p>
              </div>
              <div class="text-right">
                <p class="text-xs opacity-75">Thời gian</p>
                <p class="font-bold">{{ shiftDuration }}</p>
              </div>
            </div>
            <div class="text-xs opacity-90">
              Bắt đầu: {{ formatTime(currentShift.started_at) }}
            </div>
          </div>

          <!-- Barista Stats -->
          <div class="grid grid-cols-2 gap-3 mb-4">
            <div class="bg-white rounded-2xl p-4 shadow-sm">
              <div class="text-3xl mb-2">⏳</div>
              <div class="text-2xl font-bold text-yellow-600">{{ queuedOrders }}</div>
              <div class="text-xs text-gray-500">Chờ pha chế</div>
            </div>
            <div class="bg-white rounded-2xl p-4 shadow-sm">
              <div class="text-3xl mb-2">🍹</div>
              <div class="text-2xl font-bold text-blue-600">{{ inProgressOrders }}</div>
              <div class="text-xs text-gray-500">Đang pha (ca này)</div>
            </div>
            <div class="bg-white rounded-2xl p-4 shadow-sm">
              <div class="text-3xl mb-2">✅</div>
              <div class="text-2xl font-bold text-green-600">{{ readyOrders }}</div>
              <div class="text-xs text-gray-500">Sẵn sàng (ca này)</div>
            </div>
            <div class="bg-white rounded-2xl p-4 shadow-sm">
              <div class="text-3xl mb-2">🎯</div>
              <div class="text-2xl font-bold text-purple-600">{{ todayCompleted }}</div>
              <div class="text-xs text-gray-500">Hoàn tất (ca này)</div>
            </div>
          </div>

          <!-- Quick Actions for Barista -->
          <div class="mb-4">
            <h2 class="text-lg font-bold text-gray-800 mb-3">⚡ Thao tác nhanh</h2>
            <div class="grid grid-cols-2 gap-3">
              <button @click="$router.push('/barista')" 
                class="bg-gradient-to-br from-blue-500 to-blue-600 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
                <div class="text-4xl mb-2">🍹</div>
                <div class="font-bold">Pha chế</div>
              </button>
              <button @click="$router.push('/shifts')" 
                class="bg-gradient-to-br from-purple-500 to-purple-600 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
                <div class="text-4xl mb-2">⏰</div>
                <div class="font-bold">Ca làm</div>
              </button>
            </div>
          </div>

          <!-- Working Orders Preview -->
          <div v-if="myWorkingOrders.length > 0" class="mb-4">
            <div class="flex items-center justify-between mb-3">
              <h2 class="text-lg font-bold text-gray-800">🔥 Đang pha chế</h2>
              <button @click="$router.push('/barista')" class="text-sm text-blue-500 font-medium">
                Xem tất cả →
              </button>
            </div>
            <div class="space-y-3">
              <div v-for="order in myWorkingOrders.slice(0, 3)" :key="order.id"
                @click="$router.push('/barista')"
                class="bg-white rounded-xl p-4 shadow-sm active:scale-98 transition-transform border-l-4 border-blue-500">
                <div class="flex justify-between items-start mb-2">
                  <div>
                    <h3 class="font-bold">{{ order.order_number }}</h3>
                    <p class="text-sm text-gray-600">{{ order.items?.length || 0 }} món</p>
                  </div>
                  <span class="bg-blue-100 text-blue-800 px-2 py-1 rounded-full text-xs font-medium">
                    Đang pha
                  </span>
                </div>
                <div class="text-sm text-gray-500">
                  Bắt đầu: {{ formatTime(order.accepted_at) }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Waiter/Manager/Cashier Dashboard -->
        <div v-else>
          <!-- Cashier Dashboard -->
          <div v-if="isCashier">
          <!-- Cashier Stats -->
          <div class="grid grid-cols-2 gap-3 mb-4">
            <div class="bg-white rounded-2xl p-4 shadow-sm">
              <div class="text-3xl mb-2">📋</div>
              <div class="text-2xl font-bold text-gray-800">{{ todayOrders }}</div>
              <div class="text-xs text-gray-500">Orders hôm nay</div>
            </div>
            <div class="bg-white rounded-2xl p-4 shadow-sm">
              <div class="text-3xl mb-2">💰</div>
              <div class="text-lg font-bold text-green-600">{{ formatPrice(todayRevenue) }}</div>
              <div class="text-xs text-gray-500">Doanh thu hôm nay</div>
            </div>
          </div>

          <!-- Financial Summary (from CashierShiftManager) -->
          <CashierShiftManager />
          </div>

          <!-- Waiter/Manager Dashboard -->
          <div v-else-if="!isCashier">
            <!-- Quick Stats -->
            <div class="grid grid-cols-2 gap-3 mb-4">
              <div class="bg-white rounded-2xl p-4 shadow-sm">
                <div class="text-3xl mb-2">📋</div>
                <div class="text-2xl font-bold text-gray-800">{{ todayOrders }}</div>
                <div class="text-xs text-gray-500">Orders hôm nay</div>
              </div>
              <div class="bg-white rounded-2xl p-4 shadow-sm">
                <div class="text-3xl mb-2">💰</div>
                <div class="text-lg font-bold text-green-600">{{ formatPrice(todayRevenue) }}</div>
                <div class="text-xs text-gray-500">Doanh thu</div>
              </div>
              <div class="bg-white rounded-2xl p-4 shadow-sm">
                <div class="text-3xl mb-2">🍹</div>
                <div class="text-2xl font-bold text-blue-600">{{ inProgressOrders }}</div>
                <div class="text-xs text-gray-500">Đang pha chế</div>
              </div>
              <div class="bg-white rounded-2xl p-4 shadow-sm">
                <div class="text-3xl mb-2">⏳</div>
                <div class="text-2xl font-bold text-orange-600">{{ pendingOrders }}</div>
                <div class="text-xs text-gray-500">Chờ thanh toán</div>
              </div>
            </div>

            <!-- Quick Actions -->
            <div class="mb-4">
              <h2 class="text-lg font-bold text-gray-800 mb-3">⚡ Thao tác nhanh</h2>
              <div class="grid grid-cols-2 gap-3">
                <button @click="$router.push('/orders')" 
                  class="bg-gradient-to-br from-blue-500 to-blue-600 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
                  <div class="text-4xl mb-2">📋</div>
                  <div class="font-bold">Orders</div>
                </button>
                <button @click="$router.push('/shifts')" 
                  class="bg-gradient-to-br from-purple-500 to-purple-600 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
                  <div class="text-4xl mb-2">⏰</div>
                  <div class="font-bold">Ca làm</div>
                </button>
                
                <!-- Manager Actions -->
                <button v-if="user?.role === USER_ROLES.MANAGER" @click="$router.push('/menu')" 
                  class="bg-gradient-to-br from-orange-500 to-red-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
                  <div class="text-4xl mb-2">🍽️</div>
                  <div class="font-bold">Menu</div>
                </button>
                <button v-if="user?.role === USER_ROLES.MANAGER" @click="$router.push('/ingredients')" 
                  class="bg-gradient-to-br from-green-500 to-emerald-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
                  <div class="text-4xl mb-2">🥬</div>
                  <div class="font-bold">Nguyên liệu</div>
                </button>
                <button v-if="user?.role === USER_ROLES.MANAGER" @click="$router.push('/facilities')" 
                  class="bg-gradient-to-br from-cyan-500 to-blue-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
                  <div class="text-4xl mb-2">🏢</div>
                  <div class="font-bold">Cơ sở</div>
                </button>
                <button v-if="user?.role === USER_ROLES.MANAGER" @click="$router.push('/expenses')"
                  class="bg-gradient-to-br from-pink-500 to-purple-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
                  <div class="text-4xl mb-2">💸</div>
                  <div class="font-bold">Chi phí</div>
                </button>
                <button v-if="user?.role === USER_ROLES.MANAGER" @click="$router.push('/manager/fund')"
                  class="bg-gradient-to-br from-orange-500 to-yellow-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
                  <div class="text-4xl mb-2">💰</div>
                  <div class="font-bold">Quỹ tiền</div>
                </button>
                <button v-if="user?.role === USER_ROLES.MANAGER" @click="$router.push('/manager/inventory-stats')"
                  class="bg-gradient-to-br from-emerald-500 to-teal-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
                  <div class="text-4xl mb-2">📦</div>
                  <div class="font-bold">Tồn kho</div>
                </button>
                <!-- Cashier Actions -->
                <button v-if="isCashier" @click="$router.push('/cashier')" 
                  class="bg-gradient-to-br from-yellow-500 to-orange-500 text-white rounded-2xl p-6 shadow-lg active:scale-95 transition-transform">
                  <div class="text-4xl mb-2">💵</div>
                  <div class="font-bold">Thu ngân</div>
                </button>
              </div>
            </div>

          </div>
        </div>
      </div>
    </div>

    <!-- Quick Operating Expense Modal -->
    <transition name="slide-up">
      <div v-if="showQuickExpenseModal" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white w-full max-h-[85vh] flex flex-col rounded-t-3xl">
          <div class="sticky top-0 bg-white px-4 py-3 border-b flex items-center justify-between rounded-t-3xl">
            <button @click="showQuickExpenseModal = false" class="text-2xl text-gray-600">×</button>
            <h2 class="text-base font-bold text-gray-800">⚡ Chi phí vận hành</h2>
            <div class="w-6"></div>
          </div>

          <div class="flex-1 overflow-y-auto px-4 py-4 space-y-4">
            <!-- Operating fund balance -->
            <div class="bg-blue-50 border border-blue-200 rounded-xl p-3">
              <div class="text-xs font-semibold text-blue-700 mb-1">⚙️ Quỹ vận hành (Operating Fund)</div>
              <div v-if="quickExpenseFundLoading" class="text-xs text-gray-500">Đang tải...</div>
              <div v-else-if="quickExpenseFundBalance" class="text-sm text-blue-800">
                💵 Tiền mặt: <span class="font-bold">{{ formatCurrency(quickExpenseFundBalance.cash || 0) }}</span>
                &nbsp;|&nbsp;
                💳 CK: <span class="font-bold">{{ formatCurrency(quickExpenseFundBalance.transfer || 0) }}</span>
              </div>
            </div>

            <!-- Description -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Mô tả chi phí *</label>
              <input v-model="quickExpenseData.description" type="text"
                placeholder="VD: Tiền điện, tiền nước, vật tư..."
                class="w-full px-4 py-3 text-base border-2 border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500" />
            </div>

            <!-- Amount -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Số tiền *</label>
              <input v-model.number="quickExpenseData.amount" type="number" min="0" step="1000"
                placeholder="0"
                class="w-full px-4 py-4 text-xl font-bold border-2 border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500" />
            </div>

            <!-- Money type -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Thanh toán bằng</label>
              <div class="grid grid-cols-2 gap-2">
                <button type="button" @click="quickExpenseData.money_type = 'cash'"
                  :class="['py-3 rounded-xl font-semibold text-sm transition-colors',
                    quickExpenseData.money_type === 'cash' ? 'bg-green-500 text-white' : 'bg-gray-100 text-gray-700']">
                  💵 Tiền mặt
                </button>
                <button type="button" @click="quickExpenseData.money_type = 'transfer'"
                  :class="['py-3 rounded-xl font-semibold text-sm transition-colors',
                    quickExpenseData.money_type === 'transfer' ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-700']">
                  💳 Chuyển khoản
                </button>
              </div>
            </div>

            <!-- Notes -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Ghi chú</label>
              <input v-model="quickExpenseData.notes" type="text"
                placeholder="Ghi chú thêm (tùy chọn)"
                class="w-full px-4 py-3 text-base border border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500" />
            </div>

            <!-- Insufficient balance warning -->
            <div v-if="quickExpenseFundBalance && quickExpenseData.amount > 0">
              <div v-if="quickExpenseData.money_type === 'cash' && quickExpenseData.amount > (quickExpenseFundBalance.cash || 0)"
                class="bg-red-50 border border-red-300 rounded-xl p-3 text-xs text-red-700 font-semibold">
                ⚠️ Tiền mặt quỹ vận hành không đủ! Cần: {{ formatCurrency(quickExpenseData.amount) }} • Có: {{ formatCurrency(quickExpenseFundBalance.cash || 0) }}
              </div>
              <div v-else-if="quickExpenseData.money_type === 'transfer' && quickExpenseData.amount > (quickExpenseFundBalance.transfer || 0)"
                class="bg-red-50 border border-red-300 rounded-xl p-3 text-xs text-red-700 font-semibold">
                ⚠️ Tiền CK quỹ vận hành không đủ! Cần: {{ formatCurrency(quickExpenseData.amount) }} • Có: {{ formatCurrency(quickExpenseFundBalance.transfer || 0) }}
              </div>
            </div>

            <!-- Error -->
            <div v-if="quickExpenseError" class="bg-red-50 border border-red-200 rounded-xl p-3 text-sm text-red-600">
              {{ quickExpenseError }}
            </div>
          </div>

          <div class="flex-shrink-0 bg-white px-4 py-3 border-t flex gap-3 pb-safe">
            <button @click="showQuickExpenseModal = false"
              class="flex-1 bg-gray-200 text-gray-700 py-4 rounded-xl font-bold text-base">
              Hủy
            </button>
            <button @click="submitQuickExpense"
              :disabled="quickExpenseSubmitting || !quickExpenseData.amount || !quickExpenseData.description.trim()"
              class="flex-1 bg-orange-500 text-white py-4 rounded-xl font-bold text-base disabled:opacity-50 flex items-center justify-center gap-2">
              <span v-if="quickExpenseSubmitting" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
              {{ quickExpenseSubmitting ? 'Đang xử lý...' : '⚡ Xác nhận' }}
            </button>
          </div>
        </div>
      </div>
    </transition>

    <!-- Bottom Navigation -->
    <BottomNav />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useShiftStore } from '../stores/shift'
import { useOrderStore } from '../stores/order'
import { useBaristaStore } from '../stores/barista'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import CashierShiftManager from '../components/CashierShiftManager.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'
import { USER_ROLES } from '../constants/user'
import { ORDER_STATUS } from '../constants/order'
import { getAllBalances } from '../services/fund'
import { fundExpenseService } from '../services/fundExpenseService'

const authStore = useAuthStore()
const shiftStore = useShiftStore()
const orderStore = useOrderStore()
const baristaStore = useBaristaStore()

const currentTime = ref('')
const currentDate = ref('')
let timeInterval = null

// Computed
const user = computed(() => authStore.user)
const isBarista = computed(() => authStore.user?.role === USER_ROLES.BARISTA)
const hasOpenShift = computed(() => shiftStore.hasOpenShift)
const currentShift = computed(() => shiftStore.currentShift)
const orders = computed(() => {
  if (isBarista.value) {
    // Barista: combine queued + my orders
    return [
      ...baristaStore.queuedOrders,
      ...baristaStore.inProgressOrders,
      ...baristaStore.readyOrders
    ]
  }
  return orderStore.orders
})
const isCashier = computed(() => authStore.user?.role === USER_ROLES.CASHIER)

// Barista-specific stats - filtered by current shift
const queuedOrders = computed(() => {
  if (!isBarista.value) return 0
  // Queue orders are not assigned yet, show all
  return baristaStore.queueCount
})

const inProgressOrders = computed(() => {
  if (!isBarista.value) return 0
  
  if (!currentShift.value) {
    // No shift: show all in progress
    return baristaStore.inProgressCount
  }
  
  // With shift: filter by shift time
  return baristaStore.inProgressOrders.filter(o => 
    isInCurrentShift(o.accepted_at)
  ).length
})

const readyOrders = computed(() => {
  if (!isBarista.value) return 0
  
  if (!currentShift.value) {
    // No shift: show all ready
    return baristaStore.readyCount
  }
  
  // With shift: filter by shift time
  return baristaStore.readyOrders.filter(o => 
    isInCurrentShift(o.ready_at)
  ).length
})

const todayCompleted = computed(() => {
  if (!isBarista.value) return 0
  
  if (!currentShift.value) {
    // No shift: count today's served orders
    const today = new Date().toDateString()
    return baristaStore.servedOrders.filter(o => 
      o.ready_at && new Date(o.ready_at).toDateString() === today
    ).length
  }
  
  // With shift: count served orders in current shift
  return baristaStore.servedOrders.filter(o => 
    isInCurrentShift(o.ready_at)
  ).length
})

// Helper function to check if timestamp is in current shift
const isInCurrentShift = (timestamp) => {
  if (!timestamp || !currentShift.value) return false
  
  const orderTime = new Date(timestamp)
  const shiftStart = new Date(currentShift.value.started_at)
  const shiftEnd = currentShift.value.ended_at ? new Date(currentShift.value.ended_at) : new Date()
  
  return orderTime >= shiftStart && orderTime <= shiftEnd
}

const myWorkingOrders = computed(() => {
  if (!isBarista.value) return []
  
  // Use barista store directly
  return baristaStore.inProgressOrders
    .filter(o => {
      if (!currentShift.value) return true
      return isInCurrentShift(o.accepted_at)
    })
    .sort((a, b) => new Date(a.accepted_at) - new Date(b.accepted_at))
})

// Waiter/Manager stats
const todayOrders = computed(() => {
  const today = new Date().toDateString()
  return orders.value.filter(o => 
    new Date(o.created_at).toDateString() === today
  ).length
})

const todayRevenue = computed(() => {
  const today = new Date().toDateString()
  return orders.value
    .filter(o => new Date(o.created_at).toDateString() === today && o.status !== 'CANCELLED')
    .reduce((sum, o) => sum + o.total, 0)
})

const completedOrders = computed(() => {
  const today = new Date().toDateString()
  return orders.value.filter(o => 
    new Date(o.created_at).toDateString() === today && o.status === ORDER_STATUS.SERVED
  ).length
})

const pendingOrders = computed(() => {
  // For manager: show all orders that are not completed or cancelled
  if (user.value?.role === USER_ROLES.MANAGER) {
    return orders.value.filter(o => 
      o.status !== ORDER_STATUS.SERVED && o.status !== ORDER_STATUS.CANCELLED
    ).length
  }
  // For others: show only created orders
  return orders.value.filter(o => o.status === ORDER_STATUS.CREATED).length
})

// Cashier-specific stats
const shiftRevenue = computed(() => {
  if (!currentShift.value) return 0
  
  const shiftStart = new Date(currentShift.value.started_at)
  const shiftEnd = currentShift.value.ended_at ? new Date(currentShift.value.ended_at) : new Date()
  
  return orders.value
    .filter(o => {
      if (o.status === ORDER_STATUS.CANCELLED) return false
      const orderTime = new Date(o.created_at)
      return orderTime >= shiftStart && orderTime <= shiftEnd
    })
    .reduce((sum, o) => sum + o.total, 0)
})

const openShiftsCount = computed(() => {
  return shiftStore.openShifts.length
})

const openShifts = computed(() => {
  return shiftStore.openShifts
    .sort((a, b) => new Date(b.started_at) - new Date(a.started_at))
})

// Cashier-specific: Calculate orders and revenue per shift
const shiftsWithStats = computed(() => {
  if (!isCashier.value) return []
  
  return openShifts.value.map(shift => {
    const shiftStart = new Date(shift.started_at)
    const shiftEnd = shift.ended_at ? new Date(shift.ended_at) : new Date()
    
    // Filter orders for this shift
    const shiftOrders = orders.value.filter(o => {
      const orderTime = new Date(o.created_at)
      return orderTime >= shiftStart && orderTime <= shiftEnd
    })
    
    // Calculate revenue by payment method (exclude cancelled orders)
    const paidOrders = shiftOrders.filter(o => o.status !== ORDER_STATUS.CANCELLED)
    const cashRevenue = paidOrders
      .filter(o => o.payment_method === 'CASH')
      .reduce((sum, o) => sum + o.total, 0)
    const transferRevenue = paidOrders
      .filter(o => o.payment_method === 'TRANSFER' || o.payment_method === 'QR')
      .reduce((sum, o) => sum + o.total, 0)
    const totalRevenue = cashRevenue + transferRevenue
    
    // For barista: count orders in progress
    const inProgressCount = shift.role_type === 'barista' 
      ? shiftOrders.filter(o => o.status === ORDER_STATUS.IN_PROGRESS).length 
      : 0
    
    return {
      ...shift,
      orderCount: shiftOrders.length,
      totalRevenue: totalRevenue,
      cashRevenue: cashRevenue,
      transferRevenue: transferRevenue,
      inProgressCount: inProgressCount
    }
  })
})

const getRoleTypeText = (roleType) => {
  const roles = {
    waiter: '🍽️ Phục vụ',
    barista: '🍹 Pha chế',
    cashier: '💵 Thu ngân'
  }
  return roles[roleType] || roleType
}

const shiftDuration = computed(() => {
  if (!currentShift.value?.started_at) return ''
  const start = new Date(currentShift.value.started_at)
  const now = new Date()
  const diff = now - start
  const hours = Math.floor(diff / 3600000)
  const minutes = Math.floor((diff % 3600000) / 60000)
  return `${hours}h ${minutes}m`
})

// Methods
const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('vi-VN', { 
    hour: '2-digit', 
    minute: '2-digit' 
  })
  currentDate.value = now.toLocaleDateString('vi-VN', { 
    weekday: 'long',
    day: 'numeric',
    month: 'long'
  })
}

const formatPrice = (price) => {
  return new Intl.NumberFormat('vi-VN', { 
    style: 'currency', 
    currency: 'VND',
    maximumFractionDigits: 0
  }).format(price)
}

const formatTime = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleTimeString('vi-VN', { 
    hour: '2-digit', 
    minute: '2-digit' 
  })
}

const getShiftTypeText = (type) => {
  const types = {
    MORNING: '☀️ Ca sáng',
    AFTERNOON: '🌤️ Ca chiều',
    EVENING: '🌙 Ca tối'
  }
  return types[type] || type
}

const getStatusBadge = (status) => {
  const badges = {
    CREATED: 'bg-gray-100 text-gray-800',
    PAID: 'bg-green-100 text-green-800',
    QUEUED: 'bg-yellow-100 text-yellow-800',
    IN_PROGRESS: 'bg-blue-100 text-blue-800',
    READY: 'bg-purple-100 text-purple-800',
    SERVED: 'bg-green-100 text-green-800',
    CANCELLED: 'bg-red-100 text-red-800'
  }
  return badges[status] || 'bg-gray-100 text-gray-800'
}

const getStatusText = (status) => {
  const texts = {
    CREATED: 'Mới',
    PAID: 'Đã thu',
    QUEUED: 'Chờ pha',
    IN_PROGRESS: 'Đang pha',
    READY: 'Sẵn sàng',
    SERVED: 'Hoàn tất',
    CANCELLED: 'Đã hủy'
  }
  return texts[status] || status
}

// Refresh data function
const refreshData = async () => {
  // Manager doesn't need shift data
  if (user.value?.role === USER_ROLES.MANAGER) {
    await orderStore.fetchOrders()
    return
  }
  
  if (isBarista.value) {
    // Barista uses barista store
    await Promise.all([
      shiftStore.fetchCurrentShift(),
      baristaStore.fetchQueuedOrders(),
      baristaStore.fetchMyOrders()
    ])
  } else if (isCashier.value) {
    // Cashier doesn't need waiter/barista shift, only needs all shifts and orders
    await Promise.all([
      shiftStore.fetchAllShifts(),
      orderStore.fetchOrders()
    ])
  } else {
    // Other roles (waiter) use order store
    await Promise.all([
      shiftStore.fetchCurrentShift(),
      orderStore.fetchOrders()
    ])
  }
}

// ─── Quick Operating Expense ──────────────────────────────────────────────────

const showQuickExpenseModal = ref(false)
const quickExpenseSubmitting = ref(false)
const quickExpenseError = ref('')
const quickExpenseFundBalance = ref(null)
const quickExpenseFundLoading = ref(false)
const quickExpenseData = ref({ description: '', amount: 0, money_type: 'cash', notes: '' })

const formatCurrency = (val) => new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(val)

const openQuickExpenseModal = async () => {
  quickExpenseData.value = { description: '', amount: 0, money_type: 'cash', notes: '' }
  quickExpenseError.value = ''
  showQuickExpenseModal.value = true
  quickExpenseFundLoading.value = true
  try {
    const data = await getAllBalances()
    quickExpenseFundBalance.value = data?.['operating'] || { cash: 0, transfer: 0, total: 0 }
  } catch {
    quickExpenseFundBalance.value = null
  } finally {
    quickExpenseFundLoading.value = false
  }
}

const submitQuickExpense = async () => {
  if (quickExpenseSubmitting.value) return
  if (!quickExpenseData.value.amount || quickExpenseData.value.amount <= 0) {
    quickExpenseError.value = 'Vui lòng nhập số tiền hợp lệ'
    return
  }
  if (!quickExpenseData.value.description.trim()) {
    quickExpenseError.value = 'Vui lòng nhập mô tả chi phí'
    return
  }
  quickExpenseSubmitting.value = true
  quickExpenseError.value = ''
  try {
    await fundExpenseService.createExpenseFromFund({
      date: new Date().toISOString(),
      amount: quickExpenseData.value.amount,
      description: quickExpenseData.value.description,
      notes: quickExpenseData.value.notes,
      money_type: quickExpenseData.value.money_type
    })
    showQuickExpenseModal.value = false
  } catch (err) {
    quickExpenseError.value = err.message || 'Không thể tạo chi phí. Vui lòng thử lại.'
  } finally {
    quickExpenseSubmitting.value = false
  }
}

// Pull to refresh
const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)

// Lifecycle
onMounted(async () => {
  updateTime()
  timeInterval = setInterval(updateTime, 1000)
  await refreshData()
})

onUnmounted(() => {
  if (timeInterval) {
    clearInterval(timeInterval)
  }
})
</script>


<style scoped>
.slide-up-enter-active, .slide-up-leave-active { transition: transform 0.3s ease; }
.slide-up-enter-from, .slide-up-leave-to { transform: translateY(100%); }

.active\:scale-95:active {
  transform: scale(0.95);
}

.active\:scale-98:active {
  transform: scale(0.98);
}
</style>
