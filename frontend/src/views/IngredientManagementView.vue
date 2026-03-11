<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Pull to Refresh Indicator -->
    <PullToRefresh 
      :pull-distance="pullDistance" 
      :is-refreshing="isRefreshing"
      :threshold="80" />
    
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3" style="padding-top: max(0.75rem, env(safe-area-inset-top))">
        <div class="flex items-center justify-between mb-3">
          <h1 class="text-xl font-bold text-gray-800">🥬 Nguyên liệu</h1>
        </div>
        
        <!-- Search Bar -->
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Tìm kiếm nguyên liệu..."
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24">
      <!-- Stats Cards - Single Row -->
      <div class="bg-gradient-to-br from-blue-500 to-purple-500 rounded-xl p-4 mb-4 text-white shadow-lg">
        <div class="text-xs opacity-90 mb-2">Tổng quan</div>
        <div class="grid grid-cols-4 gap-1.5">
          <div class="text-center">
            <div class="text-lg font-bold">{{ ingredients.length }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap">Tổng</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold">{{ inStockCount }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap">Đủ hàng</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold">{{ lowStockCount }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap">Sắp hết</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold">{{ outOfStockCount }}</div>
            <div class="text-[10px] opacity-90 whitespace-nowrap">Hết hàng</div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="mb-4">
        <h2 class="text-sm font-bold text-gray-800 mb-2">⚡ Thao tác nhanh</h2>
        <div class="grid grid-cols-2 gap-2">
          <button @click="openCreateModal"
            class="bg-gradient-to-br from-blue-500 to-cyan-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform">
            <div class="text-2xl mb-1">➕</div>
            <div class="text-sm font-bold">Tạo nguyên liệu</div>
          </button>
          <button @click="showCategoryModal = true"
            class="bg-gradient-to-br from-purple-500 to-pink-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform">
            <div class="text-2xl mb-1">📁</div>
            <div class="text-sm font-bold">Quản lý danh mục</div>
          </button>
          <button @click="showLowStock"
            class="bg-gradient-to-br from-yellow-500 to-orange-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform">
            <div class="text-2xl mb-1">⚠️</div>
            <div class="text-sm font-bold">Sắp hết hàng</div>
          </button>
          <button @click="showStockHistory"
            class="bg-gradient-to-br from-green-500 to-emerald-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform">
            <div class="text-2xl mb-1">📊</div>
            <div class="text-sm font-bold">Lịch sử nhập</div>
          </button>
        </div>
      </div>

      <!-- Ingredients List -->
      <div class="mb-4">
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-lg font-bold text-gray-800">📋 Danh sách nguyên liệu</h2>
        </div>
        
        <div v-if="filteredIngredients.length === 0" class="text-center py-16">
          <div class="text-6xl mb-4">📭</div>
          <p class="text-gray-500">Không có nguyên liệu nào</p>
        </div>
        
        <div v-else class="space-y-3">
          <div v-for="ingredient in filteredIngredients" :key="ingredient.id"
            class="bg-white rounded-2xl p-4 shadow-sm">
            
            <!-- Ingredient Header -->
            <div class="flex justify-between items-start mb-3">
              <div>
                <h3 class="font-bold text-lg">{{ ingredient.name }}</h3>
                <p class="text-sm text-gray-600">{{ ingredient.category }}</p>
                <p class="text-xs text-gray-400" v-if="ingredient.supplier">📦 {{ ingredient.supplier }}</p>
              </div>
              <span :class="getStockStatusClass(ingredient)" class="px-3 py-1 rounded-full text-xs font-medium">
                {{ getStockStatusText(ingredient) }}
              </span>
            </div>

            <!-- Ingredient Info -->
            <div class="mb-3 space-y-2 text-sm">
              <div class="flex justify-between items-center">
                <span class="text-gray-600">📦 Tồn kho:</span>
                <span class="font-bold text-lg">{{ ingredient.quantity }} {{ ingredient.unit }}</span>
              </div>
              <div class="flex justify-between items-center">
                <span class="text-gray-600">⚠️ Tối thiểu:</span>
                <span class="font-medium">{{ ingredient.min_stock }} {{ ingredient.unit }}</span>
              </div>
              <div class="flex justify-between items-center bg-green-50 rounded-lg px-3 py-2 -mx-1">
                <span class="text-gray-700 font-medium">💵 Đơn giá:</span>
                <span class="font-bold text-green-600">{{ formatCurrency(ingredient.cost_per_unit) }}/{{ ingredient.unit }}</span>
              </div>
            </div>

            <!-- Quick Actions - Simplified -->
            <div class="grid grid-cols-4 gap-1.5 pt-3 border-t">
              <button @click="openAdjustModal(ingredient)"
                :disabled="isAdjusting"
                class="bg-yellow-500 text-white py-2 rounded-lg text-xs font-bold active:bg-yellow-600 disabled:opacity-50 disabled:cursor-not-allowed flex flex-col items-center justify-center">
                <span class="text-base">📦</span>
                <span class="text-[10px]">Điều chỉnh</span>
              </button>
              <button @click="viewHistory(ingredient)"
                class="bg-purple-500 text-white py-2 rounded-lg text-xs font-bold active:bg-purple-600 flex flex-col items-center justify-center">
                <span class="text-base">📊</span>
                <span class="text-[10px]">Lịch sử</span>
              </button>
              <button @click="openEditModal(ingredient)"
                :disabled="isSubmitting"
                class="bg-blue-500 text-white py-2 rounded-lg text-xs font-bold active:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed flex flex-col items-center justify-center">
                <span class="text-base">✏️</span>
                <span class="text-[10px]">Sửa</span>
              </button>
              <button @click="deleteIngredient(ingredient)"
                :disabled="isDeleting"
                class="bg-red-500 text-white py-2 rounded-lg text-xs font-bold active:bg-red-600 disabled:opacity-50 disabled:cursor-not-allowed flex flex-col items-center justify-center">
                <span class="text-base">🗑️</span>
                <span class="text-[10px]">Xóa</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />

    <!-- Category Management Modal -->
    <transition name="slide-up">
      <div v-if="showCategoryModal" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full max-h-[85vh] overflow-y-auto">
          <div class="sticky top-0 bg-white px-4 py-4 border-b flex justify-between items-center">
            <h3 class="text-lg font-bold">📁 Quản lý danh mục</h3>
            <button @click="showCategoryModal = false" class="text-2xl text-gray-400">×</button>
          </div>
          
          <div class="px-4 py-4">
            <!-- Add New Category -->
            <div class="bg-gray-50 rounded-xl p-4 mb-4">
              <h4 class="font-semibold text-gray-800 mb-3">Thêm danh mục mới</h4>
              <div class="flex gap-2">
                <input v-model="newCategoryName" type="text" placeholder="Tên danh mục..." 
                  class="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500" />
                <button @click="addCategory" class="bg-purple-500 text-white px-6 py-3 rounded-lg font-medium active:bg-purple-600">
                  Thêm
                </button>
              </div>
            </div>

            <!-- Category List -->
            <div class="space-y-2">
              <div v-for="cat in ingredientStore.categories" :key="cat.id" 
                class="bg-white border border-gray-200 rounded-xl p-4 flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-lg bg-purple-100 text-purple-600 flex items-center justify-center text-xl">
                    📦
                  </div>
                  <div>
                    <div class="font-medium text-gray-800">{{ cat.name }}</div>
                    <div class="text-xs text-gray-500">{{ getCategoryCount(cat.name) }} nguyên liệu</div>
                  </div>
                </div>
                <button @click="deleteCategory(cat.name)" class="text-red-500 hover:text-red-700 p-2">
                  🗑️
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- Create/Edit Modal - Slide from Right -->
    <transition name="slide-right">
      <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-gray-50 w-full h-screen flex flex-col">
          <!-- Mobile Header - Fixed -->
          <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
            <div class="px-4 py-3">
              <div class="flex items-center justify-between">
                <button @click="closeModal" class="text-2xl text-gray-600">←</button>
                <h1 class="text-xl font-bold text-gray-800">{{ isEditing ? '✏️ Cập nhật nguyên liệu' : '➕ Thêm nguyên liệu mới' }}</h1>
                <div class="w-8"></div>
              </div>
            </div>
          </div>

          <!-- Scrollable Content -->
          <div class="flex-1 overflow-y-auto px-4 py-6 space-y-5">
            <!-- Tên nguyên liệu -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">Tên nguyên liệu *</label>
              <input v-model="formData.name" type="text" 
                class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
            </div>

            <!-- Danh mục & Đơn vị - Responsive Grid -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">Danh mục *</label>
                <select v-model="formData.category" 
                  class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent">
                  <option value="">Chọn danh mục</option>
                  <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">Đơn vị *</label>
                <select v-model="formData.unit" 
                  class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent">
                  <option value="">Chọn đơn vị</option>
                  <option v-for="unit in INGREDIENT_UNIT_OPTIONS" :key="unit.value" :value="unit.value">
                    {{ unit.label }}
                  </option>
                </select>
              </div>
            </div>

            <!-- Số lượng nhập & Tối thiểu -->
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">
                  Số lượng nhập *
                  <span class="text-xs text-gray-500 font-normal">({{ formData.unit || 'đơn vị' }})</span>
                  <span v-if="isEditing" class="text-xs text-orange-600 font-normal ml-2">⚠️ Dùng "Điều chỉnh" để thay đổi tồn kho</span>
                </label>
                <input v-model.number="formData.quantity" type="number" step="0.01" min="0"
                  :disabled="isEditing"
                  @input="priceInputMode === 'total' && calculateUnitPrice()"
                  placeholder="VD: 2"
                  :class="isEditing ? 'bg-gray-100 cursor-not-allowed' : ''"
                  class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
                <p v-if="isEditing" class="text-xs text-gray-500 mt-2">
                  💡 Để thay đổi tồn kho, vui lòng sử dụng nút "Điều chỉnh" ở danh sách nguyên liệu
                </p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">
                  Mức tối thiểu *
                  <span class="text-xs text-gray-500 font-normal">(cảnh báo)</span>
                </label>
                <input v-model.number="formData.min_stock" type="number" step="0.01" min="0"
                  placeholder="VD: 1"
                  class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
              </div>
            </div>

            <!-- Pricing Section -->
            <div v-if="!isEditing" class="bg-gray-50 rounded-xl p-4 space-y-4">
              <div class="flex items-center justify-between mb-2">
                <h3 class="text-sm font-semibold text-gray-800">💰 Thông tin giá</h3>
                <button @click="togglePriceInputMode" type="button"
                  class="text-xs text-blue-600 underline">
                  {{ priceInputMode === 'unit' ? 'Nhập tổng giá' : 'Nhập đơn giá' }}
                </button>
              </div>

              <!-- Mode 1: Enter Unit Price (default) -->
              <div v-if="priceInputMode === 'unit'" class="space-y-3">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Đơn giá *
                    <span class="text-xs text-gray-500 font-normal">(Giá cho 1 {{ formData.unit || 'đơn vị' }})</span>
                  </label>
                  <input v-model.number="formData.cost_per_unit" type="number" step="1000" min="0"
                    placeholder="VD: 5000 (cho 1 chai)"
                    class="w-full px-4 py-4 text-base border-2 border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500" />
                  <p class="text-xs text-gray-500 mt-2">💡 Nhập giá tiền cho 1 {{ formData.unit || 'đơn vị' }}</p>
                </div>
              </div>

              <!-- Mode 2: Enter Total Price -->
              <div v-else class="space-y-3">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Tổng giá mua *
                    <span class="text-xs text-gray-500 font-normal">(Tổng tiền cho {{ formData.quantity || 0 }} {{ formData.unit || 'đơn vị' }})</span>
                  </label>
                  <input v-model.number="totalPriceInput" type="number" step="1000" min="0"
                    @input="calculateUnitPrice"
                    placeholder="VD: 10000 (cho 2 chai)"
                    class="w-full px-4 py-4 text-base border-2 border-green-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-green-500" />
                  <p class="text-xs text-gray-500 mt-2">💡 Nhập tổng tiền đã trả</p>
                </div>
                
                <!-- Show calculated unit price -->
                <div v-if="formData.quantity > 0 && totalPriceInput > 0" class="bg-blue-50 rounded-lg p-3">
                  <div class="text-xs text-gray-600 mb-1">Đơn giá được tính:</div>
                  <div class="text-base font-bold text-blue-700">
                    {{ formatCurrency(formData.cost_per_unit) }}/{{ formData.unit }}
                  </div>
                  <div class="text-xs text-gray-500 mt-1">
                    = {{ formatCurrency(totalPriceInput) }} ÷ {{ formData.quantity }} {{ formData.unit }}
                  </div>
                </div>
              </div>
            </div>

            <!-- Read-only info when editing -->
            <div v-else class="bg-blue-50 rounded-xl p-4 border-2 border-blue-200">
              <h3 class="text-sm font-semibold text-gray-800 mb-3">📊 Thông tin hiện tại (chỉ xem)</h3>
              <div class="grid grid-cols-2 gap-3">
                <div class="bg-white rounded-lg p-3">
                  <p class="text-xs text-gray-600 mb-1">Tồn kho</p>
                  <p class="text-base font-bold text-blue-700">{{ formData.quantity }} {{ formData.unit }}</p>
                </div>
                <div class="bg-white rounded-lg p-3">
                  <p class="text-xs text-gray-600 mb-1">Đơn giá</p>
                  <p class="text-base font-bold text-green-700">{{ formatCurrency(formData.cost_per_unit) }}</p>
                </div>
              </div>
              <p class="text-xs text-orange-600 mt-3">
                ⚠️ Để thay đổi tồn kho hoặc giá, vui lòng sử dụng chức năng "Điều chỉnh"
              </p>
            </div>

            <!-- Nhà cung cấp -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">Nhà cung cấp</label>
              <input v-model="formData.supplier" type="text" 
                class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
            </div>

            <!-- Ghi chú -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">Ghi chú</label>
              <textarea v-model="formData.notes" rows="3" 
                class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"></textarea>
            </div>

            <!-- Ngày tạo (only when creating) -->
            <div v-if="!isEditing">
              <label class="block text-sm font-medium text-gray-700 mb-3">
                Ngày tạo
                <span class="text-xs text-gray-500 font-normal">(mặc định: hôm nay)</span>
              </label>
              <input v-model="formData.created_date" type="datetime-local"
                class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
              <p class="text-xs text-gray-500 mt-2">
                💡 Để trống để sử dụng thời gian hiện tại
              </p>
            </div>

            <!-- Fund deduction section (only when creating with cost > 0) -->
            <div v-if="!isEditing && formData.quantity > 0 && formData.cost_per_unit > 0" class="space-y-3">
              <!-- Inventory fund balance -->
              <div class="bg-purple-50 border border-purple-200 rounded-xl p-3">
                <div class="text-xs font-semibold text-purple-700 mb-1">📦 Quỹ hàng hóa (Inventory Fund)</div>
                <div v-if="fundBalanceLoading" class="text-xs text-gray-500">Đang tải...</div>
                <div v-else-if="fundBalance" class="text-sm text-purple-800">
                  💵 Tiền mặt: <span class="font-bold">{{ formatCurrency(fundBalance.cash || 0) }}</span>
                  &nbsp;|&nbsp;
                  💳 CK: <span class="font-bold">{{ formatCurrency(fundBalance.transfer || 0) }}</span>
                </div>
              </div>

              <!-- Money type selector -->
              <div>
                <div class="text-sm font-medium text-gray-700 mb-2">Thanh toán từ quỹ bằng</div>
                <div class="grid grid-cols-2 gap-2">
                  <button type="button" @click="formData.money_type = 'cash'"
                    :class="['py-3 rounded-xl font-semibold text-sm transition-colors',
                      formData.money_type !== 'transfer' ? 'bg-green-500 text-white' : 'bg-gray-100 text-gray-700']">
                    💵 Tiền mặt
                  </button>
                  <button type="button" @click="formData.money_type = 'transfer'"
                    :class="['py-3 rounded-xl font-semibold text-sm transition-colors',
                      formData.money_type === 'transfer' ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-700']">
                    💳 Chuyển khoản
                  </button>
                </div>
              </div>

              <!-- Deduction summary -->
              <div class="bg-orange-50 border-2 border-orange-300 rounded-xl p-3">
                <div class="flex items-center gap-2 mb-1">
                  <span class="text-lg">📦</span>
                  <span class="text-sm font-semibold text-orange-800">Tự động trừ Quỹ hàng hóa</span>
                </div>
                <div class="text-xl font-bold text-orange-700">{{ formatCurrency(totalCost) }}</div>
                <div class="text-xs text-gray-500 mt-1">
                  {{ formData.quantity }} {{ formData.unit }} × {{ formatCurrency(formData.cost_per_unit) }}
                </div>
                <div v-if="formData.money_type !== 'transfer' && totalCost > (fundBalance?.cash || 0)"
                  class="text-xs text-red-600 mt-1 font-semibold">⚠️ Tiền mặt quỹ không đủ!</div>
                <div v-if="formData.money_type === 'transfer' && totalCost > (fundBalance?.transfer || 0)"
                  class="text-xs text-red-600 mt-1 font-semibold">⚠️ Tiền chuyển khoản quỹ không đủ!</div>
              </div>
            </div>

            <!-- Spacer for bottom buttons -->
            <div class="h-24"></div>
          </div>

          <!-- Fixed Footer -->
          <div class="flex-shrink-0 bg-white px-4 py-4 border-t flex gap-3 pb-safe">
            <button @click="closeModal" 
              :disabled="isSubmitting"
              class="flex-1 bg-gray-200 text-gray-700 py-4 rounded-xl font-medium text-base active:bg-gray-300 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
              Hủy
            </button>
            <button @click="saveIngredient" 
              :disabled="isSubmitting"
              class="flex-1 bg-blue-500 text-white py-4 rounded-xl font-medium text-base active:bg-blue-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2">
              <span v-if="isSubmitting" class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
              <span>{{ isSubmitting ? 'Đang lưu...' : (isEditing ? 'Cập nhật' : 'Thêm mới') }}</span>
            </button>
          </div>
        </div>
      </div>
    </transition>

      <!-- Quick Stock IN Modal - Simplified -->
      <transition name="slide-up">
        <div v-if="showQuickInModal" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
          <div class="bg-white w-full max-h-[80vh] flex flex-col rounded-t-3xl">
            <!-- Header -->
            <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0 rounded-t-3xl">
              <div class="px-4 py-3">
                <div class="flex items-center justify-between mb-2">
                  <button @click="closeQuickInModal" class="text-2xl text-gray-600 -ml-2">×</button>
                  <h1 class="text-base font-bold text-gray-800">➕ Nhập kho nhanh</h1>
                  <div class="w-6"></div>
                </div>
                <p class="text-sm text-gray-700 text-center font-semibold">{{ currentIngredient?.name }}</p>
                <p class="text-xs text-gray-500 text-center">
                  Tồn: <span class="font-semibold">{{ currentIngredient?.quantity }} {{ currentIngredient?.unit }}</span>
                </p>
              </div>
            </div>

            <!-- Content -->
            <div class="flex-1 overflow-y-auto px-4 py-4 space-y-4">
              <!-- Số lượng -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Số lượng nhập *
                  <span class="text-xs text-gray-500 font-normal">({{ currentIngredient?.unit }})</span>
                </label>
                <input v-model.number="quickInData.quantity" type="number" step="0.01" min="0.01"
                  @input="quickInPriceMode === 'total' && calculateQuickInUnitPrice()"
                  placeholder="VD: 5"
                  class="w-full px-4 py-4 text-lg font-bold border-2 border-gray-300 rounded-xl focus:ring-2 focus:ring-green-500 focus:border-green-500" />
              </div>

              <!-- Giá -->
              <div class="bg-green-50 rounded-xl p-4 space-y-3">
                <div class="flex items-center justify-between">
                  <h3 class="text-sm font-semibold text-gray-800">💰 Giá nhập</h3>
                  <button @click="toggleQuickInPriceMode" type="button"
                    class="text-xs text-green-600 underline font-medium">
                    {{ quickInPriceMode === 'unit' ? 'Nhập tổng giá' : 'Nhập đơn giá' }}
                  </button>
                </div>

                <!-- Unit Price Mode -->
                <div v-if="quickInPriceMode === 'unit'">
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Đơn giá
                    <span class="text-xs text-gray-500 font-normal">({{ currentIngredient?.unit }})</span>
                  </label>
                  <input v-model.number="quickInData.cost_per_unit" type="number" step="1000" min="0"
                    :placeholder="`Hiện tại: ${formatCurrency(currentIngredient?.cost_per_unit || 0)}`"
                    class="w-full px-4 py-4 text-lg font-bold border-2 border-green-300 rounded-xl focus:ring-2 focus:ring-green-500" />
                  <p class="text-xs text-gray-500 mt-2">
                    💡 Để 0 nếu giá không đổi
                  </p>
                </div>

                <!-- Total Price Mode -->
                <div v-else>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Tổng giá mua
                    <span class="text-xs text-gray-500 font-normal">({{ quickInData.quantity || 0 }} {{ currentIngredient?.unit }})</span>
                  </label>
                  <input v-model.number="quickInTotalPrice" type="number" step="1000" min="0"
                    @input="calculateQuickInUnitPrice"
                    placeholder="VD: 100000"
                    class="w-full px-4 py-4 text-lg font-bold border-2 border-green-300 rounded-xl focus:ring-2 focus:ring-green-500" />
                  
                  <div v-if="quickInData.quantity > 0 && quickInTotalPrice > 0" class="bg-white rounded-lg p-3 mt-2">
                    <div class="text-xs text-gray-600 mb-1">Đơn giá được tính:</div>
                    <div class="text-base font-bold text-green-700">
                      {{ formatCurrency(quickInData.cost_per_unit) }}/{{ currentIngredient?.unit }}
                    </div>
                  </div>
                </div>
              </div>

              <!-- Preview -->
              <div class="bg-blue-50 rounded-xl p-3">
                <div class="text-xs text-gray-600 mb-1">Tồn kho sau nhập:</div>
                <div class="text-xl font-bold text-blue-900">
                  {{ (currentIngredient?.quantity || 0) + (quickInData.quantity || 0) }} {{ currentIngredient?.unit }}
                </div>
              </div>

              <!-- Fund deduction indicator -->
              <div v-if="quickInData.quantity > 0 && quickInExpenseAmount > 0" class="space-y-3">
                <!-- Inventory fund balance -->
                <div class="bg-purple-50 border border-purple-200 rounded-xl p-3">
                  <div class="text-xs font-semibold text-purple-700 mb-1">Quỹ hàng hóa (Inventory Fund)</div>
                  <div v-if="fundBalanceLoading" class="text-xs text-gray-500">Đang tải...</div>
                  <div v-else class="text-sm text-purple-800">
                    💵 Tiền mặt: <span class="font-bold">{{ formatCurrency(fundBalance?.cash || 0) }}</span>
                    &nbsp;|&nbsp;
                    💳 CK: <span class="font-bold">{{ formatCurrency(fundBalance?.transfer || 0) }}</span>
                  </div>
                </div>

                <!-- Money type selector -->
                <div>
                  <div class="text-xs font-medium text-gray-700 mb-2">Thanh toán bằng</div>
                  <div class="grid grid-cols-2 gap-2">
                    <button type="button" @click="quickInData.money_type = 'cash'"
                      :class="['py-2 px-3 rounded-lg text-sm font-medium transition-colors',
                        quickInData.money_type === 'cash' ? 'bg-green-500 text-white' : 'bg-gray-100 text-gray-700']">
                      💵 Tiền mặt
                    </button>
                    <button type="button" @click="quickInData.money_type = 'transfer'"
                      :class="['py-2 px-3 rounded-lg text-sm font-medium transition-colors',
                        quickInData.money_type === 'transfer' ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-700']">
                      💳 Chuyển khoản
                    </button>
                  </div>
                </div>

                <!-- Deduction summary -->
                <div class="bg-orange-50 border-2 border-orange-300 rounded-xl p-3">
                  <div class="flex items-center gap-2 mb-1">
                    <span class="text-orange-600 text-lg">📦</span>
                    <span class="text-sm font-semibold text-orange-800">Tự động trừ Quỹ hàng hóa</span>
                  </div>
                  <div class="text-lg font-bold text-orange-700">
                    {{ formatCurrency(quickInExpenseAmount) }}
                  </div>
                  <div v-if="quickInData.money_type === 'cash' && quickInExpenseAmount > (fundBalance?.cash || 0)"
                    class="text-xs text-red-600 mt-1 font-semibold">
                    ⚠️ Tiền mặt quỹ không đủ!
                  </div>
                  <div v-if="quickInData.money_type === 'transfer' && quickInExpenseAmount > (fundBalance?.transfer || 0)"
                    class="text-xs text-red-600 mt-1 font-semibold">
                    ⚠️ Tiền chuyển khoản quỹ không đủ!
                  </div>
                </div>
              </div>
            </div>

            <!-- Footer -->
            <div class="flex-shrink-0 bg-white px-4 py-3 border-t flex gap-3 pb-safe">
              <button @click="closeQuickInModal" 
                :disabled="isAdjusting"
                class="flex-1 bg-gray-200 text-gray-700 py-4 rounded-xl font-bold text-base active:bg-gray-300 disabled:opacity-50">
                Hủy
              </button>
              <button @click="confirmQuickIn" 
                :disabled="isAdjusting || !quickInData.quantity || quickInData.quantity <= 0"
                class="flex-1 bg-green-500 text-white py-4 rounded-xl font-bold text-base active:bg-green-600 disabled:opacity-50 flex items-center justify-center gap-2">
                <span v-if="isAdjusting" class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
                <span>{{ isAdjusting ? 'Đang xử lý...' : '✓ Xác nhận' }}</span>
              </button>
            </div>
          </div>
        </div>
      </transition>

      <!-- Quick Stock OUT Modal - Simplified -->
      <transition name="slide-up">
        <div v-if="showQuickOutModal" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
          <div class="bg-white w-full max-h-[75vh] flex flex-col rounded-t-3xl">
            <!-- Header -->
            <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0 rounded-t-3xl">
              <div class="px-4 py-3">
                <div class="flex items-center justify-between mb-2">
                  <button @click="closeQuickOutModal" class="text-2xl text-gray-600 -ml-2">×</button>
                  <h1 class="text-base font-bold text-gray-800">➖ Xuất kho nhanh</h1>
                  <div class="w-6"></div>
                </div>
                <p class="text-sm text-gray-700 text-center font-semibold">{{ currentIngredient?.name }}</p>
                <p class="text-xs text-gray-500 text-center">
                  Tồn: <span class="font-semibold">{{ currentIngredient?.quantity }} {{ currentIngredient?.unit }}</span>
                </p>
              </div>
            </div>

            <!-- Content -->
            <div class="flex-1 overflow-y-auto px-4 py-4 space-y-4">
              <!-- Số lượng -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Số lượng xuất *
                  <span class="text-xs text-gray-500 font-normal">({{ currentIngredient?.unit }})</span>
                </label>
                <input v-model.number="quickOutData.quantity" type="number" step="0.01" min="0.01"
                  :max="currentIngredient?.quantity"
                  placeholder="VD: 2"
                  class="w-full px-4 py-4 text-lg font-bold border-2 border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-orange-500" />
                <p class="text-xs text-gray-500 mt-1">
                  Tối đa: {{ currentIngredient?.quantity }} {{ currentIngredient?.unit }}
                </p>
              </div>

              <!-- Lý do -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Lý do *</label>
                <select v-model="quickOutData.reason" 
                  class="w-full px-4 py-4 text-base border-2 border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500 mb-2">
                  <option value="">Chọn lý do...</option>
                  <option value="Sử dụng cho món ăn">🍽️ Sử dụng cho món ăn</option>
                  <option value="Hỏng/Hư">❌ Hỏng/Hư</option>
                  <option value="Hết hạn">⏰ Hết hạn</option>
                  <option value="Thất thoát">📉 Thất thoát</option>
                  <option value="Kiểm kê">📋 Kiểm kê</option>
                  <option value="custom">✏️ Lý do khác...</option>
                </select>
                
                <input v-if="quickOutData.reason === 'custom'" 
                  v-model="quickOutData.customReason" 
                  type="text"
                  placeholder="Nhập lý do..."
                  class="w-full px-4 py-3 text-base border-2 border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500" />
              </div>

              <!-- Preview -->
              <div class="bg-orange-50 rounded-xl p-3">
                <div class="text-xs text-gray-600 mb-1">Tồn kho sau xuất:</div>
                <div class="text-xl font-bold" 
                  :class="(currentIngredient?.quantity || 0) - (quickOutData.quantity || 0) < (currentIngredient?.min_stock || 0) ? 'text-red-600' : 'text-orange-900'">
                  {{ Math.max(0, (currentIngredient?.quantity || 0) - (quickOutData.quantity || 0)) }} {{ currentIngredient?.unit }}
                </div>
                <p v-if="(currentIngredient?.quantity || 0) - (quickOutData.quantity || 0) < (currentIngredient?.min_stock || 0)" 
                  class="text-xs text-red-600 mt-1">
                  ⚠️ Dưới mức tối thiểu!
                </p>
              </div>
            </div>

            <!-- Footer -->
            <div class="flex-shrink-0 bg-white px-4 py-3 border-t flex gap-3 pb-safe">
              <button @click="closeQuickOutModal" 
                :disabled="isAdjusting"
                class="flex-1 bg-gray-200 text-gray-700 py-4 rounded-xl font-bold text-base active:bg-gray-300 disabled:opacity-50">
                Hủy
              </button>
              <button @click="confirmQuickOut" 
                :disabled="isAdjusting || !quickOutData.quantity || quickOutData.quantity <= 0 || !quickOutData.reason"
                class="flex-1 bg-orange-500 text-white py-4 rounded-xl font-bold text-base active:bg-orange-600 disabled:opacity-50 flex items-center justify-center gap-2">
                <span v-if="isAdjusting" class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
                <span>{{ isAdjusting ? 'Đang xử lý...' : '✓ Xác nhận' }}</span>
              </button>
            </div>
          </div>
        </div>
      </transition>

      <!-- Adjust Stock Modal - Mobile First -->
      <transition name="slide-up">
        <div v-if="showAdjustModal" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
          <div class="bg-gray-50 w-full h-[90vh] flex flex-col rounded-t-3xl">
            <!-- Mobile Header - Fixed -->
            <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0 rounded-t-3xl">
              <div class="px-4 py-4">
                <div class="flex items-center justify-between mb-3">
                  <button @click="closeAdjustModal" class="text-2xl text-gray-600">×</button>
                  <h1 class="text-xl font-bold text-gray-800">📦 Điều chỉnh tồn kho</h1>
                  <div class="w-8"></div>
                </div>
                
                <!-- Current Info Card -->
                <div class="bg-blue-50 rounded-xl p-3">
                  <p class="text-sm text-gray-700">
                    <span class="font-semibold">{{ currentIngredient?.name }}</span>
                  </p>
                  <p class="text-xs text-gray-600 mt-1">
                    Tồn kho: <span class="font-bold text-blue-700">{{ currentIngredient?.quantity }} {{ currentIngredient?.unit }}</span>
                  </p>
                </div>
              </div>
            </div>

            <!-- Scrollable Content -->
            <div class="flex-1 overflow-y-auto px-4 py-4 space-y-5">
              <!-- Paid from Fund Checkbox (only for ADD type) -->
              <div v-if="adjustData.type === ADJUSTMENT_TYPES.ADD" class="bg-purple-50 rounded-xl p-4 border-2 border-purple-300">
                <label class="flex items-start gap-3 cursor-pointer">
                  <input 
                    type="checkbox" 
                    v-model="adjustData.paidFromFund"
                    @change="onPaidFromFundChange"
                    class="mt-1 w-5 h-5 text-purple-600 rounded focus:ring-2 focus:ring-purple-500" />
                  <div class="flex-1">
                    <div class="font-semibold text-gray-800 mb-1">💰 Chi từ quỹ</div>
                    <p class="text-xs text-gray-600">
                      Tự động trừ tiền từ quỹ và ghi nhận giao dịch
                    </p>
                    <!-- Fund Balance Display -->
                    <div v-if="adjustData.paidFromFund" class="mt-3 bg-white rounded-lg p-3 border border-purple-300">
                      <div v-if="fundBalanceLoading" class="text-xs text-gray-500">
                        Đang tải số dư quỹ...
                      </div>
                      <div v-else-if="fundBalanceError" class="text-xs text-red-600">
                        ⚠️ {{ fundBalanceError }}
                      </div>
                      <div v-else>
                        <div class="text-xs text-gray-600 mb-1">Quỹ hàng hóa (Inventory Fund):</div>
                        <div class="text-lg font-bold text-purple-700">
                          {{ formatCurrency(fundBalance?.total || 0) }}
                        </div>
                        <div class="text-xs text-gray-500 mt-1">
                          Tiền mặt: {{ formatCurrency(fundBalance?.cash || 0) }} • 
                          Chuyển khoản: {{ formatCurrency(fundBalance?.transfer || 0) }}
                        </div>
                        
                        <!-- Money Type Selector -->
                        <div class="mt-3 pt-3 border-t border-purple-200">
                          <label class="block text-xs font-medium text-gray-700 mb-2">Trừ từ *</label>
                          <div class="grid grid-cols-2 gap-2">
                            <button type="button"
                              @click="adjustData.money_type = 'cash'"
                              :class="adjustData.money_type === 'cash' ? 'bg-green-500 text-white' : 'bg-gray-100 text-gray-700'"
                              class="py-2 px-3 rounded-lg text-sm font-medium transition-colors">
                              💵 Tiền mặt
                            </button>
                            <button type="button"
                              @click="adjustData.money_type = 'transfer'"
                              :class="adjustData.money_type === 'transfer' ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-700'"
                              class="py-2 px-3 rounded-lg text-sm font-medium transition-colors">
                              💳 Chuyển khoản
                            </button>
                          </div>
                        </div>
                      </div>
                    </div>
                    <!-- Insufficient Balance Warning -->
                    <div v-if="adjustData.paidFromFund && !fundBalanceLoading && !fundBalanceError && adjustExpenseAmount > 0" 
                      class="mt-2">
                      <!-- Check specific money type balance -->
                      <div v-if="adjustData.money_type === 'cash' && adjustExpenseAmount > (fundBalance?.cash || 0)"
                        class="bg-red-50 border border-red-300 rounded-lg p-2">
                        <div class="flex items-start gap-2">
                          <span class="text-red-600 text-lg">⚠️</span>
                          <div class="flex-1">
                            <p class="text-xs font-semibold text-red-800">Tiền mặt trong quỹ không đủ!</p>
                            <p class="text-xs text-red-600 mt-1">
                              Cần: {{ formatCurrency(adjustExpenseAmount) }} • 
                              Có: {{ formatCurrency(fundBalance?.cash || 0) }}
                            </p>
                          </div>
                        </div>
                      </div>
                      <div v-else-if="adjustData.money_type === 'transfer' && adjustExpenseAmount > (fundBalance?.transfer || 0)"
                        class="bg-red-50 border border-red-300 rounded-lg p-2">
                        <div class="flex items-start gap-2">
                          <span class="text-red-600 text-lg">⚠️</span>
                          <div class="flex-1">
                            <p class="text-xs font-semibold text-red-800">Tiền chuyển khoản trong quỹ không đủ!</p>
                            <p class="text-xs text-red-600 mt-1">
                              Cần: {{ formatCurrency(adjustExpenseAmount) }} • 
                              Có: {{ formatCurrency(fundBalance?.transfer || 0) }}
                            </p>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </label>
              </div>
              <!-- Loại điều chỉnh -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">Loại điều chỉnh *</label>
                <select v-model="adjustData.type" 
                  class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500">
                  <option v-for="type in ADJUSTMENT_TYPE_OPTIONS" :key="type.value" :value="type.value">
                    {{ type.label }}
                  </option>
                </select>
              </div>

              <!-- Số lượng -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">
                  Số lượng *
                  <span class="text-xs text-gray-500 font-normal">({{ currentIngredient?.unit }})</span>
                </label>
                <input v-model.number="adjustData.quantity" type="number" step="0.01" 
                  @input="onAdjustQuantityChange"
                  placeholder="VD: 2"
                  class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500" />
              </div>

              <!-- Price input for stock IN (type = 'add') -->
              <div v-if="adjustData.type === ADJUSTMENT_TYPES.ADD" class="bg-gray-50 rounded-xl p-4 space-y-4">
                <div class="flex items-center justify-between mb-2">
                  <h3 class="text-sm font-semibold text-gray-800">💰 Giá nhập lần này</h3>
                  <button @click="toggleAdjustPriceMode" type="button"
                    class="text-xs text-blue-600 underline">
                    {{ adjustPriceMode === 'unit' ? 'Nhập tổng giá' : 'Nhập đơn giá' }}
                  </button>
                </div>

                <!-- Mode 1: Enter Unit Price -->
                <div v-if="adjustPriceMode === 'unit'" class="space-y-3">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      Đơn giá
                      <span class="text-xs text-gray-500 font-normal">(Giá cho 1 {{ currentIngredient?.unit }})</span>
                    </label>
                    <input v-model.number="adjustData.cost_per_unit" type="number" step="1000" min="0"
                      @input="calculateAdjustExpense"
                      :placeholder="`Giá hiện tại: ${formatCurrency(currentIngredient?.cost_per_unit || 0)}`"
                      class="w-full px-4 py-4 text-base border-2 border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500" />
                    <p class="text-xs text-gray-500 mt-2">
                      💡 Để trống (0) nếu giá không đổi. Giá hiện tại: {{ formatCurrency(currentIngredient?.cost_per_unit || 0) }}/{{ currentIngredient?.unit }}
                    </p>
                  </div>
                </div>

                <!-- Mode 2: Enter Total Price -->
                <div v-else class="space-y-3">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      Tổng giá mua
                      <span class="text-xs text-gray-500 font-normal">(Tổng tiền cho {{ adjustData.quantity || 0 }} {{ currentIngredient?.unit }})</span>
                    </label>
                    <input v-model.number="adjustTotalPrice" type="number" step="1000" min="0"
                      @input="calculateAdjustUnitPrice"
                      placeholder="VD: 50000 (cho 2 kg)"
                      class="w-full px-4 py-4 text-base border-2 border-green-300 rounded-lg focus:ring-2 focus:ring-green-500" />
                    <p class="text-xs text-gray-500 mt-2">💡 Nhập tổng tiền đã trả cho lần nhập này</p>
                  </div>
                  
                  <!-- Show calculated unit price -->
                  <div v-if="adjustData.quantity > 0 && adjustTotalPrice > 0" class="bg-blue-50 rounded-lg p-3">
                    <div class="text-xs text-gray-600 mb-1">Đơn giá được tính:</div>
                    <div class="text-base font-bold text-blue-700">
                      {{ formatCurrency(adjustData.cost_per_unit) }}/{{ currentIngredient?.unit }}
                    </div>
                    <div class="text-xs text-gray-500 mt-1">
                      = {{ formatCurrency(adjustTotalPrice) }} ÷ {{ adjustData.quantity }} {{ currentIngredient?.unit }}
                    </div>
                  </div>

                  <!-- Show current price for reference -->
                  <div class="bg-gray-100 rounded-lg p-2">
                    <p class="text-xs text-gray-600">
                      Giá hiện tại: <span class="font-semibold">{{ formatCurrency(currentIngredient?.cost_per_unit || 0) }}/{{ currentIngredient?.unit }}</span>
                    </p>
                    <p class="text-xs text-gray-500 mt-1">💡 Để trống (0) nếu giá không đổi</p>
                  </div>
                </div>
              </div>

              <!-- Price input for stock ADJUST (type = 'adjust') - Only show if quantity increases -->
              <div v-if="adjustData.type === ADJUSTMENT_TYPES.ADJUST && adjustData.quantity > (currentIngredient?.quantity || 0)" 
                class="bg-yellow-50 rounded-xl p-4 space-y-4 border-2 border-yellow-300">
                <div class="flex items-start gap-2 mb-2">
                  <span class="text-yellow-600 text-xl">⚠️</span>
                  <div class="flex-1">
                    <h3 class="text-sm font-semibold text-gray-800 mb-1">Số lượng tăng - Có nhập thêm hàng không?</h3>
                    <p class="text-xs text-gray-600">
                      Nếu tăng do nhập thêm hàng với giá mới, hãy nhập giá bên dưới. Nếu chỉ điều chỉnh số liệu, để trống.
                    </p>
                  </div>
                </div>

                <div class="flex items-center justify-between mb-2">
                  <h3 class="text-sm font-semibold text-gray-800">💰 Giá nhập thêm (tùy chọn)</h3>
                  <button @click="toggleAdjustPriceMode" type="button"
                    class="text-xs text-yellow-700 underline">
                    {{ adjustPriceMode === 'unit' ? 'Nhập tổng giá' : 'Nhập đơn giá' }}
                  </button>
                </div>

                <!-- Mode 1: Enter Unit Price -->
                <div v-if="adjustPriceMode === 'unit'" class="space-y-3">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      Đơn giá mới
                      <span class="text-xs text-gray-500 font-normal">(nếu có)</span>
                    </label>
                    <input v-model.number="adjustData.cost_per_unit" type="number" step="1000" min="0"
                      @input="calculateAdjustExpense"
                      :placeholder="`Để 0 nếu không đổi giá`"
                      class="w-full px-4 py-4 text-base border-2 border-yellow-300 rounded-lg focus:ring-2 focus:ring-yellow-500" />
                    <p class="text-xs text-gray-500 mt-2">
                      💡 Giá hiện tại: {{ formatCurrency(currentIngredient?.cost_per_unit || 0) }}/{{ currentIngredient?.unit }}
                    </p>
                  </div>
                </div>

                <!-- Mode 2: Enter Total Price -->
                <div v-else class="space-y-3">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      Tổng giá mua thêm
                      <span class="text-xs text-gray-500 font-normal">(nếu có)</span>
                    </label>
                    <input v-model.number="adjustTotalPrice" type="number" step="1000" min="0"
                      @input="calculateAdjustUnitPrice"
                      placeholder="Để trống nếu không mua thêm"
                      class="w-full px-4 py-4 text-base border-2 border-yellow-300 rounded-lg focus:ring-2 focus:ring-yellow-500" />
                  </div>
                  
                  <!-- Show calculated unit price -->
                  <div v-if="adjustData.quantity > (currentIngredient?.quantity || 0) && adjustTotalPrice > 0" class="bg-white rounded-lg p-3 border border-yellow-300">
                    <div class="text-xs text-gray-600 mb-1">Đơn giá được tính cho phần tăng thêm:</div>
                    <div class="text-base font-bold text-yellow-700">
                      {{ formatCurrency(adjustData.cost_per_unit) }}/{{ currentIngredient?.unit }}
                    </div>
                    <div class="text-xs text-gray-500 mt-1">
                      = {{ formatCurrency(adjustTotalPrice) }} ÷ {{ (adjustData.quantity - (currentIngredient?.quantity || 0)) }} {{ currentIngredient?.unit }} (phần tăng)
                    </div>
                  </div>
                </div>
              </div>

              <!-- Lý do -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">Lý do *</label>
                <textarea v-model="adjustData.reason" rows="3" 
                  placeholder="VD: Nhập hàng từ nhà cung cấp"
                  class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 resize-none"></textarea>
              </div>

              <!-- Preview -->
              <div class="bg-blue-50 rounded-xl p-4">
                <p class="text-sm text-blue-800 mb-1">Tồn kho sau điều chỉnh:</p>
                <p class="text-xl font-bold text-blue-900">
                  {{ calculateNewQuantity() }} {{ currentIngredient?.unit }}
                </p>
              </div>

              <!-- Auto-Expense Indicator for Stock IN -->
              <div v-if="adjustData.type === ADJUSTMENT_TYPES.ADD && adjustData.quantity > 0 && adjustExpenseAmount > 0" 
                :class="adjustData.paidFromFund ? 'bg-purple-50 border-purple-300' : 'bg-green-50 border-green-300'"
                class="border-2 rounded-xl p-4">
                <div class="flex items-start gap-3">
                  <span :class="adjustData.paidFromFund ? 'text-purple-600' : 'text-green-600'" class="text-2xl flex-shrink-0">✅</span>
                  <div class="flex-1 min-w-0">
                    <p :class="adjustData.paidFromFund ? 'text-purple-800' : 'text-green-800'" class="text-sm font-semibold mb-2">
                      {{ adjustData.paidFromFund ? 'Chi phí sẽ được trừ từ quỹ' : 'Tự động ghi nhận chi phí' }}
                    </p>
                    <div class="bg-white rounded-lg p-3 mb-2">
                      <div class="text-xs text-gray-600 mb-1">Chi phí nhập thêm:</div>
                      <div :class="adjustData.paidFromFund ? 'text-purple-700' : 'text-green-700'" class="text-base font-bold">
                        {{ formatCurrency(adjustExpenseAmount) }}
                      </div>
                      <div class="text-xs text-gray-500 mt-1">
                        = {{ adjustData.quantity }} {{ currentIngredient.unit }} × {{ formatCurrency(effectiveAdjustPrice) }}/{{ currentIngredient.unit }}
                      </div>
                    </div>
                    <p :class="adjustData.paidFromFund ? 'text-purple-600' : 'text-green-600'" class="text-xs">
                      📝 Danh mục: "Nguyên liệu" • 💵 {{ adjustData.paidFromFund ? 'Nguồn: Quỹ' : 'Phương thức: Tiền mặt' }}
                    </p>
                  </div>
                </div>
              </div>

              <!-- Spacer for bottom buttons -->
              <div class="h-24"></div>
            </div>

            <!-- Fixed Footer -->
            <div class="flex-shrink-0 bg-white px-4 py-4 border-t flex gap-3 pb-safe">
              <button @click="closeAdjustModal" 
                :disabled="isAdjusting"
                class="flex-1 bg-gray-200 text-gray-700 py-4 rounded-xl font-medium text-base active:bg-gray-300 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
                Hủy
              </button>
              <button @click="adjustStock" 
                :disabled="isAdjusting"
                class="flex-1 bg-blue-500 text-white py-4 rounded-xl font-medium text-base active:bg-blue-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2">
                <span v-if="isAdjusting" class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
                <span>{{ isAdjusting ? 'Đang xử lý...' : 'Xác nhận' }}</span>
              </button>
            </div>
          </div>
        </div>
      </transition>

      <!-- Stock History Modal - Mobile First Compact -->
      <transition name="slide-up">
        <div v-if="showHistoryModal" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
          <div class="bg-white w-full h-[90vh] flex flex-col rounded-t-3xl">
            <!-- Mobile Header - Fixed -->
            <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0 rounded-t-3xl">
              <div class="px-4 py-3">
                <div class="flex items-center justify-between mb-2">
                  <button @click="closeHistoryModal" class="text-2xl text-gray-600 -ml-2">×</button>
                  <h1 class="text-base font-bold text-gray-800">📊 Lịch sử nhập hàng</h1>
                  <div class="w-6"></div>
                </div>
                <p class="text-sm text-gray-700 text-center font-semibold">{{ currentIngredient?.name }}</p>
                <p class="text-xs text-gray-500 text-center">
                  Giá hiện tại: <span class="font-semibold text-green-600">{{ formatCurrency(currentIngredient?.cost_per_unit || 0) }}</span>
                </p>
              </div>
            </div>

            <!-- Scrollable Content -->
            <div class="flex-1 overflow-y-auto px-3 py-3">
              <div v-if="stockHistory.length === 0" class="text-center py-16">
                <div class="text-5xl mb-3">📭</div>
                <p class="text-sm text-gray-500">Chưa có lịch sử nhập hàng</p>
              </div>

              <div v-else class="space-y-2 pb-6">
                <div v-for="record in stockHistory" :key="record.id" 
                  class="rounded-xl p-3 shadow-sm border"
                  :class="record.quantity > 0 ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'">
                  
                  <!-- Compact Header: Type + Quantity + Date -->
                  <div class="flex items-center justify-between mb-2">
                    <div class="flex items-center gap-2">
                      <span :class="getAdjustmentTypeClass(record.type)" 
                        class="px-2 py-0.5 text-[10px] font-bold rounded-full">
                        {{ getAdjustmentTypeText(record.type) }}
                      </span>
                      <span class="text-base font-bold" :class="record.quantity > 0 ? 'text-green-700' : 'text-red-700'">
                        {{ record.quantity > 0 ? '+' : '' }}{{ record.quantity }} {{ currentIngredient?.unit }}
                      </span>
                    </div>
                    <span class="text-[10px] text-gray-500">
                      {{ formatCompactDate(record.created_at) }}
                    </span>
                  </div>

                  <!-- Fund Payment Badge -->
                  <div v-if="record.paid_from_fund" class="mb-2">
                    <span class="inline-flex items-center gap-1 px-2 py-1 bg-purple-100 text-purple-700 rounded-full text-[10px] font-bold">
                      💰 Chi từ quỹ
                    </span>
                  </div>

                  <!-- Price Info - Compact for purchases -->
                  <div v-if="record.cost_per_unit > 0 && record.quantity > 0" 
                    :class="record.paid_from_fund ? 'border-purple-300 bg-purple-50' : 'border-green-300 bg-white'"
                    class="rounded-lg p-2 mb-2 border">
                    <div class="flex items-center justify-between mb-1">
                      <span class="text-[10px] text-gray-600">Đơn giá:</span>
                      <span :class="record.paid_from_fund ? 'text-purple-700' : 'text-green-700'" class="text-xs font-bold">
                        {{ formatCurrency(record.cost_per_unit) }}
                      </span>
                    </div>
                    <div class="flex items-center justify-between">
                      <span class="text-[10px] text-gray-600">Tổng:</span>
                      <span :class="record.paid_from_fund ? 'text-purple-800' : 'text-green-800'" class="text-sm font-bold">
                        {{ formatCurrency(record.total_cost || 0) }}
                      </span>
                    </div>
                    <!-- Price comparison - compact -->
                    <div v-if="currentIngredient?.cost_per_unit !== record.cost_per_unit" 
                      :class="record.paid_from_fund ? 'border-purple-200' : 'border-green-200'"
                      class="flex items-center justify-center gap-1 text-[10px] mt-1 pt-1 border-t">
                      <span :class="record.cost_per_unit > currentIngredient?.cost_per_unit ? 'text-red-600' : 'text-green-600'" 
                        class="font-semibold">
                        {{ record.cost_per_unit > currentIngredient?.cost_per_unit ? '↑' : '↓' }}
                        {{ formatCurrency(Math.abs(record.cost_per_unit - currentIngredient?.cost_per_unit)) }}
                      </span>
                      <span class="text-gray-500">vs hiện tại</span>
                    </div>
                  </div>

                  <!-- Fund Transaction Link -->
                  <div v-if="record.paid_from_fund && (record.fund_transaction_id || record.expense_id)" 
                    class="mb-2 bg-purple-50 rounded-lg p-2 border border-purple-200">
                    <div class="text-[10px] text-gray-600 mb-1">Liên kết:</div>
                    <div class="flex flex-wrap gap-1">
                      <a v-if="record.expense_id" 
                        :href="`#/expenses/${record.expense_id}`"
                        class="inline-flex items-center gap-1 px-2 py-0.5 bg-white border border-purple-300 rounded text-[10px] text-purple-700 hover:bg-purple-100">
                        📝 Chi phí
                      </a>
                      <a v-if="record.fund_transaction_id" 
                        :href="`#/fund/transactions/${record.fund_transaction_id}`"
                        class="inline-flex items-center gap-1 px-2 py-0.5 bg-white border border-purple-300 rounded text-[10px] text-purple-700 hover:bg-purple-100">
                        💰 Giao dịch quỹ
                      </a>
                    </div>
                  </div>

                  <!-- Reason - Compact -->
                  <div class="text-xs text-gray-600 mb-1.5 line-clamp-2">
                    {{ record.reason }}
                  </div>
                  
                  <!-- Footer: User + Stock Change -->
                  <div class="flex items-center justify-between text-[10px] text-gray-500 pt-1.5 border-t border-gray-200">
                    <span class="font-medium">👤 {{ record.username }}</span>
                    <span class="font-semibold text-gray-700">
                      {{ record.before_qty }} → {{ record.after_qty }}
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Fixed Footer -->
            <div class="flex-shrink-0 bg-white px-4 py-3 border-t pb-safe">
              <button @click="closeHistoryModal" 
                class="w-full bg-gray-600 text-white py-3 rounded-xl font-medium text-sm active:bg-gray-700 transition-colors">
                Đóng
              </button>
            </div>
          </div>
        </div>
      </transition>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useIngredientStore } from '../stores/ingredient'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'
import {
  INGREDIENT_UNIT_OPTIONS,
  ADJUSTMENT_TYPE_OPTIONS,
  ADJUSTMENT_TYPES,
  getStockStatusClass,
  getStockStatusText,
  getAdjustmentTypeClass,
  getAdjustmentTypeText
} from '../constants/ingredient'
import { getAllBalances } from '../services/fund'
import { fundIngredientService } from '../services/fundIngredientService'

export default {
  name: 'IngredientManagementView',
  components: {
    BottomNav,
    PullToRefresh
  },
  setup() {
    const ingredientStore = useIngredientStore()
    
    const searchQuery = ref('')
    const showModal = ref(false)
    const showAdjustModal = ref(false)
    const showHistoryModal = ref(false)
    const showCategoryModal = ref(false)
    const showQuickInModal = ref(false)
    const showQuickOutModal = ref(false)
    const isEditing = ref(false)
    const currentIngredient = ref(null)
    const newCategoryName = ref('')
    
    // Loading states to prevent double submission
    const isSubmitting = ref(false)
    const isDeleting = ref(false)
    const isAdjusting = ref(false)

    // Quick IN data
    const quickInData = ref({
      quantity: 0,
      cost_per_unit: 0
    })
    const quickInPriceMode = ref('total')
    const quickInTotalPrice = ref(0)

    // Quick OUT data
    const quickOutData = ref({
      quantity: 0,
      reason: '',
      customReason: ''
    })

    const toggleQuickInPriceMode = () => {
      quickInPriceMode.value = quickInPriceMode.value === 'unit' ? 'total' : 'unit'
      if (quickInPriceMode.value === 'unit') {
        quickInTotalPrice.value = 0
      } else {
        if (quickInData.value.quantity > 0 && quickInData.value.cost_per_unit > 0) {
          quickInTotalPrice.value = quickInData.value.quantity * quickInData.value.cost_per_unit
        }
      }
    }

    const calculateQuickInUnitPrice = () => {
      if (quickInData.value.quantity > 0 && quickInTotalPrice.value > 0) {
        quickInData.value.cost_per_unit = quickInTotalPrice.value / quickInData.value.quantity
      } else {
        quickInData.value.cost_per_unit = 0
      }
    }

    const quickInExpenseAmount = computed(() => {
      const qty = quickInData.value.quantity || 0
      const price = quickInData.value.cost_per_unit || currentIngredient.value?.cost_per_unit || 0
      return qty * price
    })

    const quickStockIn = (ingredient) => {
      currentIngredient.value = ingredient
      quickInData.value = {
        quantity: 0,
        cost_per_unit: 0,
        money_type: 'cash'
      }
      quickInPriceMode.value = 'total'
      quickInTotalPrice.value = 0
      showQuickInModal.value = true
      fetchFundBalance()
    }

    const closeQuickInModal = () => {
      showQuickInModal.value = false
      currentIngredient.value = null
      quickInData.value = { quantity: 0, cost_per_unit: 0 }
      quickInTotalPrice.value = 0
    }

    const confirmQuickIn = async () => {
      if (isAdjusting.value) return
      if (!quickInData.value.quantity || quickInData.value.quantity <= 0) {
        alert('Vui lòng nhập số lượng hợp lệ')
        return
      }

      isAdjusting.value = true
      try {
        const costPerUnit = quickInData.value.cost_per_unit || 0
        if (costPerUnit > 0) {
          // Auto-deduct from inventory fund
          const data = {
            quantity: quickInData.value.quantity,
            cost_per_unit: costPerUnit,
            reason: 'Nhập kho',
            money_type: quickInData.value.money_type || 'cash'
          }
          await fundIngredientService.restockIngredientFromFund(currentIngredient.value.id, data)
        } else {
          // No cost — just update stock without fund deduction
          const data = {
            quantity: quickInData.value.quantity,
            cost_per_unit: 0,
            reason: 'Nhập kho'
          }
          await ingredientStore.stockIn(currentIngredient.value.id, data)
        }
        closeQuickInModal()
      } catch (error) {
        console.error('Error quick stock in:', error)
        alert(error.message || 'Có lỗi xảy ra khi nhập kho')
      } finally {
        isAdjusting.value = false
      }
    }

    const quickStockOut = (ingredient) => {
      currentIngredient.value = ingredient
      quickOutData.value = {
        quantity: 0,
        reason: '',
        customReason: ''
      }
      showQuickOutModal.value = true
    }

    const closeQuickOutModal = () => {
      showQuickOutModal.value = false
      currentIngredient.value = null
      quickOutData.value = { quantity: 0, reason: '', customReason: '' }
    }

    const confirmQuickOut = async () => {
      if (isAdjusting.value) return
      if (!quickOutData.value.quantity || quickOutData.value.quantity <= 0) {
        alert('Vui lòng nhập số lượng hợp lệ')
        return
      }
      if (!quickOutData.value.reason) {
        alert('Vui lòng chọn lý do')
        return
      }
      if (quickOutData.value.quantity > currentIngredient.value.quantity) {
        alert('Số lượng xuất không được lớn hơn tồn kho')
        return
      }

      isAdjusting.value = true
      try {
        const finalReason = quickOutData.value.reason === 'custom' 
          ? quickOutData.value.customReason 
          : quickOutData.value.reason

        const data = {
          quantity: quickOutData.value.quantity,
          reason: finalReason
        }
        
        await ingredientStore.stockOut(currentIngredient.value.id, data)
        closeQuickOutModal()
      } catch (error) {
        console.error('Error quick stock out:', error)
        alert('Có lỗi xảy ra khi xuất kho')
      } finally {
        isAdjusting.value = false
      }
    }
    
    // Categories from store
    const categories = computed(() => ingredientStore.categories.map(c => c.name))
    
    const addCategory = async () => {
      if (!newCategoryName.value.trim()) return
      if (categories.value.includes(newCategoryName.value.trim())) {
        alert('Danh mục đã tồn tại!')
        return
      }
      const success = await ingredientStore.createCategory({ name: newCategoryName.value.trim() })
      if (success) {
        newCategoryName.value = ''
      } else {
        alert(ingredientStore.error || 'Lỗi tạo danh mục')
      }
    }
    
    const deleteCategory = async (categoryName) => {
      const hasIngredients = ingredients.value.some(item => item.category === categoryName)
      if (hasIngredients) {
        alert('Không thể xóa danh mục đã có nguyên liệu!')
        return
      }
      if (confirm(`Bạn có chắc muốn xóa danh mục "${categoryName}"?`)) {
        const category = ingredientStore.categories.find(c => c.name === categoryName)
        if (category) {
          const success = await ingredientStore.deleteCategory(category.id)
          if (!success) {
            alert(ingredientStore.error || 'Lỗi xóa danh mục')
          }
        }
      }
    }
    
    const getCategoryCount = (categoryName) => {
      return ingredients.value.filter(item => item.category === categoryName).length
    }
    
    const formData = ref({
      name: '',
      category: '',
      unit: '',
      quantity: 0,
      min_stock: 0,
      cost_per_unit: 0,
      supplier: '',
      notes: '',
      created_date: '' // Will be set to current datetime when opening create modal
    })

    // Price input mode: 'unit' (enter unit price) or 'total' (enter total price)
    const priceInputMode = ref('total') // Default to total price mode (more intuitive)
    const totalPriceInput = ref(0)

    const togglePriceInputMode = () => {
      priceInputMode.value = priceInputMode.value === 'unit' ? 'total' : 'unit'
      // Clear the other input when switching
      if (priceInputMode.value === 'unit') {
        totalPriceInput.value = 0
      } else {
        // When switching to total mode, calculate total from unit price
        if (formData.value.quantity > 0 && formData.value.cost_per_unit > 0) {
          totalPriceInput.value = formData.value.quantity * formData.value.cost_per_unit
        }
      }
    }

    const calculateUnitPrice = () => {
      if (formData.value.quantity > 0 && totalPriceInput.value > 0) {
        formData.value.cost_per_unit = totalPriceInput.value / formData.value.quantity
      } else {
        formData.value.cost_per_unit = 0
      }
    }

    const adjustData = ref({
      type: ADJUSTMENT_TYPES.ADD,
      quantity: 0,
      reason: '',
      cost_per_unit: 0, // Price for this transaction (optional)
      paidFromFund: false, // NEW: Flag for fund payment
      money_type: 'cash' // NEW: Money type for fund payment (cash or transfer)
    })

    // Fund balance state
    const fundBalance = ref(null)
    const fundBalanceLoading = ref(false)
    const fundBalanceError = ref(null)

    // Price input mode for adjust modal
    const adjustPriceMode = ref('total') // Default to total price mode
    const adjustTotalPrice = ref(0)

    const toggleAdjustPriceMode = () => {
      adjustPriceMode.value = adjustPriceMode.value === 'unit' ? 'total' : 'unit'
      // Clear the other input when switching
      if (adjustPriceMode.value === 'unit') {
        adjustTotalPrice.value = 0
      } else {
        // When switching to total mode, calculate total from unit price
        if (adjustData.value.quantity > 0 && adjustData.value.cost_per_unit > 0) {
          adjustTotalPrice.value = adjustData.value.quantity * adjustData.value.cost_per_unit
        }
      }
    }

    const calculateAdjustUnitPrice = () => {
      if (adjustData.value.type === ADJUSTMENT_TYPES.ADJUST) {
        // For adjust type, calculate based on the INCREASE amount
        const quantityDiff = adjustData.value.quantity - (currentIngredient.value?.quantity || 0)
        if (quantityDiff > 0 && adjustTotalPrice.value > 0) {
          adjustData.value.cost_per_unit = adjustTotalPrice.value / quantityDiff
        } else {
          adjustData.value.cost_per_unit = 0
        }
      } else {
        // For add/remove type, calculate normally
        if (adjustData.value.quantity > 0 && adjustTotalPrice.value > 0) {
          adjustData.value.cost_per_unit = adjustTotalPrice.value / adjustData.value.quantity
        } else {
          adjustData.value.cost_per_unit = 0
        }
      }
    }

    const onAdjustQuantityChange = () => {
      // When quantity changes, recalculate if in total price mode
      if (adjustPriceMode.value === 'total') {
        calculateAdjustUnitPrice()
      }
      
      // Reset price inputs when switching from increase to decrease (or vice versa)
      if (adjustData.value.type === ADJUSTMENT_TYPES.ADJUST) {
        const isIncreasing = adjustData.value.quantity > (currentIngredient.value?.quantity || 0)
        if (!isIncreasing) {
          // If decreasing, clear price inputs
          adjustData.value.cost_per_unit = 0
          adjustTotalPrice.value = 0
        }
      }
    }

    // Computed properties for adjust expense calculation
    const effectiveAdjustPrice = computed(() => {
      // Use provided price, or fall back to current ingredient price
      if (adjustData.value.cost_per_unit > 0) {
        return adjustData.value.cost_per_unit
      }
      return currentIngredient.value?.cost_per_unit || 0
    })

    const adjustExpenseAmount = computed(() => {
      if (adjustData.value.type !== 'add' || adjustData.value.quantity <= 0) {
        return 0
      }
      return adjustData.value.quantity * effectiveAdjustPrice.value
    })

    const calculateAdjustExpense = () => {
      // Trigger reactivity for expense calculation
      // The computed properties will handle the actual calculation
    }

    const ingredients = computed(() => ingredientStore.items || [])
    const stockHistory = ref([])

    // Calculate total cost for expense
    const totalCost = computed(() => {
      const quantity = formData.value.quantity || 0
      const costPerUnit = formData.value.cost_per_unit || 0
      return quantity * costPerUnit
    })

    const filteredIngredients = computed(() => {
      let filtered = ingredients.value
      
      // Filter by search query
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        filtered = filtered.filter(i => 
          i.name.toLowerCase().includes(query) ||
          i.category.toLowerCase().includes(query) ||
          i.supplier?.toLowerCase().includes(query)
        )
      }
      
      // Sort by created_at (newest first)
      return [...filtered].sort((a, b) => {
        const dateA = new Date(a.created_at || 0)
        const dateB = new Date(b.created_at || 0)
        return dateB - dateA // Newest first
      })
    })

    const inStockCount = computed(() => 
      ingredients.value.filter(i => i.quantity > i.min_stock).length
    )
    const lowStockCount = computed(() => 
      ingredients.value.filter(i => i.quantity > 0 && i.quantity <= i.min_stock).length
    )
    const outOfStockCount = computed(() => 
      ingredients.value.filter(i => i.quantity === 0).length
    )

    const formatCurrency = (value) => {
      if (value === undefined || value === null || isNaN(value)) {
        return '0 ₫'
      }
      return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(value)
    }

    const formatDateTime = (date) => {
      if (!date) return 'N/A'
      return new Date(date).toLocaleString('vi-VN')
    }

    const formatCompactDate = (date) => {
      if (!date) return 'N/A'
      const d = new Date(date)
      const now = new Date()
      const diffMs = now - d
      const diffMins = Math.floor(diffMs / 60000)
      const diffHours = Math.floor(diffMs / 3600000)
      const diffDays = Math.floor(diffMs / 86400000)
      
      // Less than 1 hour: show minutes
      if (diffMins < 60) {
        return `${diffMins} phút trước`
      }
      // Less than 24 hours: show hours
      if (diffHours < 24) {
        return `${diffHours} giờ trước`
      }
      // Less than 7 days: show days
      if (diffDays < 7) {
        return `${diffDays} ngày trước`
      }
      // More than 7 days: show full date and time
      return d.toLocaleString('vi-VN', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    const calculateNewQuantity = () => {
      if (!currentIngredient.value) return 0
      const current = currentIngredient.value.quantity
      const adjust = adjustData.value.quantity || 0
      
      if (adjustData.value.type === ADJUSTMENT_TYPES.ADD) return current + adjust
      if (adjustData.value.type === ADJUSTMENT_TYPES.REMOVE) return current - adjust
      return adjust
    }

    const openCreateModal = () => {
      isEditing.value = false
      currentIngredient.value = null
      
      // Set default created_date to current datetime in local timezone
      const now = new Date()
      const year = now.getFullYear()
      const month = String(now.getMonth() + 1).padStart(2, '0')
      const day = String(now.getDate()).padStart(2, '0')
      const hours = String(now.getHours()).padStart(2, '0')
      const minutes = String(now.getMinutes()).padStart(2, '0')
      const defaultDateTime = `${year}-${month}-${day}T${hours}:${minutes}`
      
      formData.value = {
        name: '',
        category: '',
        unit: '',
        quantity: 0,
        min_stock: 0,
        cost_per_unit: 0,
        supplier: '',
        notes: '',
        created_date: defaultDateTime,
        money_type: 'cash'
      }
      priceInputMode.value = 'total' // Default to total price mode
      totalPriceInput.value = 0
      showModal.value = true
      fetchFundBalance()
    }

    const openEditModal = (ingredient) => {
      isEditing.value = true
      currentIngredient.value = ingredient
      formData.value = { ...ingredient }
      // When editing, calculate total price from existing data
      if (ingredient.quantity > 0 && ingredient.cost_per_unit > 0) {
        totalPriceInput.value = ingredient.quantity * ingredient.cost_per_unit
      }
      priceInputMode.value = 'total' // Default to total price mode
      showModal.value = true
    }

    const closeModal = () => {
      showModal.value = false
      isEditing.value = false
      currentIngredient.value = null
      isSubmitting.value = false // Reset loading state
    }

    const openAdjustModal = (ingredient) => {
      currentIngredient.value = ingredient
      adjustData.value = {
        type: ADJUSTMENT_TYPES.ADD,
        quantity: 0,
        reason: '',
        cost_per_unit: 0, // Will use current price if left as 0
        paidFromFund: false,
        money_type: 'cash' // Default to cash
      }
      adjustPriceMode.value = 'total' // Default to total price mode
      adjustTotalPrice.value = 0
      isAdjusting.value = false // Reset loading state
      
      // Reset fund balance state
      fundBalance.value = null
      fundBalanceError.value = null
      
      showAdjustModal.value = true
    }

    // Fetch inventory fund balance
    const fetchFundBalance = async () => {
      fundBalanceLoading.value = true
      fundBalanceError.value = null
      try {
        const data = await getAllBalances()
        // Get inventory fund balance specifically
        fundBalance.value = data?.['inventory'] || { cash: 0, transfer: 0, total: 0 }
      } catch (error) {
        console.error('Error fetching fund balance:', error)
        fundBalanceError.value = 'Không thể tải số dư quỹ hàng hóa. Vui lòng thử lại.'
      } finally {
        fundBalanceLoading.value = false
      }
    }

    // Handle paid from fund checkbox change
    const onPaidFromFundChange = () => {
      if (adjustData.value.paidFromFund) {
        fetchFundBalance()
      }
    }

    const closeAdjustModal = () => {
      showAdjustModal.value = false
      currentIngredient.value = null
      isAdjusting.value = false // Reset loading state
    }

    const closeHistoryModal = () => {
      showHistoryModal.value = false
      currentIngredient.value = null
    }

    const saveIngredient = async () => {
      // Prevent double submission
      if (isSubmitting.value) return
      
      // Validation
      if (!formData.value.name || !formData.value.category || !formData.value.unit) {
        alert('Vui lòng điền đầy đủ thông tin bắt buộc')
        return
      }
      
      if (formData.value.quantity < 0 || formData.value.min_stock < 0 || formData.value.cost_per_unit < 0) {
        alert('Số lượng và giá không được âm')
        return
      }
      
      isSubmitting.value = true
      try {
        if (isEditing.value) {
          await ingredientStore.updateIngredient(currentIngredient.value.id, formData.value)
        } else {
          const initialQty = formData.value.quantity || 0
          const costPerUnit = formData.value.cost_per_unit || 0
          // Create ingredient with quantity=0 first, then restock from fund if needed
          const payload = { ...formData.value, quantity: 0 }
          const newIngredient = await ingredientStore.createIngredient(payload)
          if (!newIngredient) throw new Error('Không thể tạo nguyên liệu')

          if (initialQty > 0 && costPerUnit > 0 && newIngredient.id) {
            // Deduct from inventory fund
            await fundIngredientService.restockIngredientFromFund(newIngredient.id, {
              quantity: initialQty,
              cost_per_unit: costPerUnit,
              reason: 'Nhập kho ban đầu',
              money_type: formData.value.money_type || 'cash'
            })
          } else if (initialQty > 0) {
            // No cost → just add stock without fund deduction
            await ingredientStore.stockIn(newIngredient.id, {
              quantity: initialQty,
              cost_per_unit: costPerUnit,
              reason: 'Nhập kho ban đầu'
            })
          }
        }
        closeModal()
      } catch (error) {
        console.error('Error saving ingredient:', error)
        alert(error.message || 'Có lỗi xảy ra khi lưu nguyên liệu')
      } finally {
        isSubmitting.value = false
      }
    }

    const adjustStock = async () => {
      // Prevent double submission
      if (isAdjusting.value) return
      
      // Validation
      if (!adjustData.value.quantity || adjustData.value.quantity < 0) {
        alert('Vui lòng nhập số lượng hợp lệ')
        return
      }
      
      // Check if new quantity would be negative
      const newQty = calculateNewQuantity()
      if (newQty < 0) {
        alert('Số lượng sau điều chỉnh không được âm')
        return
      }

      // Validate fund balance if paying from fund
      if (adjustData.value.paidFromFund && adjustData.value.type === ADJUSTMENT_TYPES.ADD) {
        if (!fundBalance.value) {
          alert('Không thể tải số dư quỹ. Vui lòng thử lại.')
          return
        }
        
        const totalCost = adjustExpenseAmount.value
        if (totalCost > fundBalance.value.total) {
          alert(`Số dư quỹ không đủ!\nCần: ${formatCurrency(totalCost)}\nCó: ${formatCurrency(fundBalance.value.total)}`)
          return
        }
      }
      
      isAdjusting.value = true
      try {
        // Use new API methods based on operation type
        if (adjustData.value.type === ADJUSTMENT_TYPES.ADD) {
          // Stock IN - check if paid from fund
          if (adjustData.value.paidFromFund) {
            // Validate money_type selection
            if (!adjustData.value.money_type) {
              alert('Vui lòng chọn loại tiền (Tiền mặt hoặc Chuyển khoản)')
              return
            }
            
            // Validate balance for selected money type
            const totalCost = adjustData.value.quantity * (adjustData.value.cost_per_unit || currentIngredient.value.cost_per_unit || 0)
            if (adjustData.value.money_type === 'cash' && totalCost > (fundBalance.value?.cash || 0)) {
              alert(`Tiền mặt trong quỹ không đủ!\nCần: ${formatCurrency(totalCost)}\nCó: ${formatCurrency(fundBalance.value?.cash || 0)}`)
              return
            }
            if (adjustData.value.money_type === 'transfer' && totalCost > (fundBalance.value?.transfer || 0)) {
              alert(`Tiền chuyển khoản trong quỹ không đủ!\nCần: ${formatCurrency(totalCost)}\nCó: ${formatCurrency(fundBalance.value?.transfer || 0)}`)
              return
            }
            
            // Use fund ingredient service
            const data = {
              quantity: adjustData.value.quantity,
              cost_per_unit: adjustData.value.cost_per_unit || currentIngredient.value.cost_per_unit || 0,
              reason: adjustData.value.reason,
              money_type: adjustData.value.money_type // Add money_type
            }
            await fundIngredientService.restockIngredientFromFund(currentIngredient.value.id, data)
          } else {
            // Regular stock in
            const data = {
              quantity: adjustData.value.quantity,
              cost_per_unit: adjustData.value.cost_per_unit || 0,
              reason: adjustData.value.reason
            }
            await ingredientStore.stockIn(currentIngredient.value.id, data)
          }
        } else if (adjustData.value.type === ADJUSTMENT_TYPES.REMOVE) {
          // Stock OUT
          const data = {
            quantity: adjustData.value.quantity,
            reason: adjustData.value.reason
          }
          await ingredientStore.stockOut(currentIngredient.value.id, data)
        } else if (adjustData.value.type === ADJUSTMENT_TYPES.ADJUST) {
          // Stock ADJUST - set to specific quantity
          const isIncrease = adjustData.value.quantity > currentIngredient.value.quantity
          
          console.log('=== ADJUST DEBUG ===')
          console.log('Current quantity:', currentIngredient.value.quantity)
          console.log('New quantity:', adjustData.value.quantity)
          console.log('Is increase?', isIncrease)
          console.log('cost_per_unit before:', adjustData.value.cost_per_unit)
          
          const data = {
            new_quantity: adjustData.value.quantity,
            // Only send cost_per_unit if:
            // 1. It's an increase AND
            // 2. User has entered a new price (not 0 or empty)
            cost_per_unit: (isIncrease && adjustData.value.cost_per_unit > 0) ? adjustData.value.cost_per_unit : 0,
            reason: adjustData.value.reason
          }
          
          console.log('Data to send:', data)
          console.log('===================')
          
          await ingredientStore.stockAdjust(currentIngredient.value.id, data)
        }
        
        closeAdjustModal()
      } catch (error) {
        console.error('Error adjusting stock:', error)
        alert(error.message || 'Có lỗi xảy ra khi điều chỉnh tồn kho')
      } finally {
        isAdjusting.value = false
      }
    }

    const deleteIngredient = async (ingredient) => {
      // Prevent double deletion
      if (isDeleting.value) return
      
      if (confirm(`Bạn có chắc muốn xóa nguyên liệu "${ingredient.name}"?`)) {
        isDeleting.value = true
        try {
          await ingredientStore.deleteIngredient(ingredient.id)
        } catch (error) {
          console.error('Error deleting ingredient:', error)
          alert('Có lỗi xảy ra khi xóa nguyên liệu')
        } finally {
          isDeleting.value = false
        }
      }
    }

    const viewHistory = async (ingredient) => {
      currentIngredient.value = ingredient
      try {
        // Try to use the new fund-aware API first
        const response = await fundIngredientService.getRestockHistory(ingredient.id)
        stockHistory.value = response.history || response.records || []
      } catch (error) {
        console.error('Error fetching fund restock history, falling back to regular history:', error)
        // Fallback to regular stock history if fund API fails
        stockHistory.value = await ingredientStore.fetchStockHistory(ingredient.id)
      }
      showHistoryModal.value = true
    }

    const showLowStock = async () => {
      await ingredientStore.fetchLowStock()
      searchQuery.value = '' // Clear search to show filtered results
    }

    const showStockHistory = () => {
      // TODO: Implement stock history view
      alert('Chức năng lịch sử nhập kho đang được phát triển')
    }

    // Refresh data function
    const refreshData = async () => {
      await ingredientStore.fetchCategories()
      await ingredientStore.fetchIngredients()
    }

    // Pull to refresh
    const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)

    onMounted(async () => {
      await refreshData()
    })

    return {
      searchQuery,
      showModal,
      showAdjustModal,
      showHistoryModal,
      showCategoryModal,
      showQuickInModal,
      showQuickOutModal,
      isEditing,
      currentIngredient,
      formData,
      adjustData,
      ingredients,
      filteredIngredients,
      stockHistory,
      inStockCount,
      lowStockCount,
      outOfStockCount,
      categories,
      newCategoryName,
      ingredientStore,
      INGREDIENT_UNIT_OPTIONS,
      ADJUSTMENT_TYPE_OPTIONS,
      ADJUSTMENT_TYPES,
      priceInputMode,
      totalPriceInput,
      totalCost,
      effectiveAdjustPrice,
      adjustExpenseAmount,
      adjustPriceMode,
      adjustTotalPrice,
      fundBalance,
      fundBalanceLoading,
      fundBalanceError,
      togglePriceInputMode,
      calculateUnitPrice,
      toggleAdjustPriceMode,
      calculateAdjustUnitPrice,
      onAdjustQuantityChange,
      calculateAdjustExpense,
      fetchFundBalance,
      onPaidFromFundChange,
      getStockStatusClass,
      getStockStatusText,
      getAdjustmentTypeClass,
      getAdjustmentTypeText,
      formatCurrency,
      formatDateTime,
      formatCompactDate,
      calculateNewQuantity,
      openCreateModal,
      openEditModal,
      closeModal,
      openAdjustModal,
      closeAdjustModal,
      closeHistoryModal,
      saveIngredient,
      adjustStock,
      deleteIngredient,
      viewHistory,
      showLowStock,
      addCategory,
      deleteCategory,
      getCategoryCount,
      pullDistance,
      isRefreshing,
      isSubmitting,
      isDeleting,
      isAdjusting,
      quickInData,
      quickInPriceMode,
      quickInTotalPrice,
      quickInExpenseAmount,
      quickOutData,
      toggleQuickInPriceMode,
      calculateQuickInUnitPrice,
      quickStockIn,
      quickStockOut,
      closeQuickInModal,
      closeQuickOutModal,
      confirmQuickIn,
      confirmQuickOut
    }
  }
}
</script>

<style scoped>
.active\:scale-95:active {
  transform: scale(0.95);
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

.slide-right-enter-active,
.slide-right-leave-active {
  transition: transform 0.3s ease;
}

.slide-right-enter-from {
  transform: translateX(100%);
}

.slide-right-leave-to {
  transform: translateX(100%);
}

.pb-safe {
  padding-bottom: max(1rem, env(safe-area-inset-bottom));
}

/* Spinner animation */
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.animate-spin {
  animation: spin 1s linear infinite;
}

/* Line clamp for compact text */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
