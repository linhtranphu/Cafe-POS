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
          <h1 class="text-xl font-bold text-gray-800">🍽️ Quản lý Menu</h1>
        </div>
        
        <!-- Search Bar -->
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Tìm kiếm món..."
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24">
      <!-- Loading & Error States -->
      <div v-if="loading" class="text-center py-16">
        <div class="text-6xl mb-4">⏳</div>
        <p class="text-gray-500">Đang tải...</p>
      </div>
      
      <div v-else-if="error" class="text-center py-16">
        <div class="text-6xl mb-4">⚠️</div>
        <p class="text-red-600">{{ error }}</p>
      </div>

      <!-- Quick Actions -->
      <div v-else class="mb-4">
        <h2 class="text-sm font-bold text-gray-800 mb-2">⚡ Thao tác nhanh</h2>
        <div class="grid grid-cols-2 gap-2">
          <button @click="openCreateModal"
            class="bg-gradient-to-br from-blue-500 to-cyan-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform">
            <div class="text-2xl mb-1">➕</div>
            <div class="text-sm font-bold">Thêm món</div>
          </button>
          <button @click="showCategoryModal = true"
            class="bg-gradient-to-br from-purple-500 to-pink-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform">
            <div class="text-2xl mb-1">📁</div>
            <div class="text-sm font-bold">Danh mục</div>
          </button>
        </div>
      </div>

      <!-- Menu Items by Category -->
      <div v-if="!loading && !error" class="space-y-4">
        <div v-for="category in filteredGroupedItems" :key="category.name" class="bg-white rounded-2xl p-4 shadow-sm">
          <div class="flex items-center gap-3 mb-4 pb-3 border-b-2 border-blue-500">
            <div class="w-10 h-10 rounded-xl flex items-center justify-center text-2xl" :class="getCategoryColor(category.name)">
              {{ getCategoryIcon(category.name) }}
            </div>
            <div>
              <h3 class="text-lg font-bold text-gray-800">{{ category.name }}</h3>
              <p class="text-xs text-gray-500">{{ category.items.length }} món</p>
            </div>
          </div>
          
          <div class="space-y-3">
            <div v-for="item in category.items" :key="item.id" 
              @click="viewDetails(item)"
              class="rounded-xl p-4 bg-gray-50 active:scale-98 transition-transform">
              
              <!-- Item Header -->
              <div class="flex justify-between items-start mb-3">
                <div class="flex-1 min-w-0">
                  <h4 class="font-bold text-lg text-gray-900 truncate">{{ item.name }}</h4>
                  <p class="text-sm text-gray-600 line-clamp-2">{{ item.description || 'Chưa có mô tả' }}</p>
                </div>
                <span class="ml-3 px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap flex-shrink-0" 
                  :class="item.available ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'">
                  {{ item.available ? '✅ Có' : '❌ Hết' }}
                </span>
              </div>

              <!-- Price - Single Size -->
              <div v-if="!item.has_variants" class="bg-white rounded-lg p-3 mb-3">
                <div class="text-2xl font-bold text-green-600">{{ formatPrice(item.price) }}</div>
                <div class="text-xs text-gray-500">Giá bán</div>
              </div>

              <!-- Variants - Multi Size -->
              <div v-else class="bg-white rounded-lg p-3 mb-3">
                <div class="text-xs font-semibold text-gray-700 mb-2">📏 Kích cỡ:</div>
                <div class="space-y-2">
                  <div v-for="variant in item.variants" :key="variant.id" 
                    class="flex justify-between items-center p-2 bg-gray-50 rounded-lg">
                    <div class="flex items-center gap-2">
                      <span class="font-medium text-gray-900">{{ variant.name }}</span>
                      <span v-if="variant.is_default" class="text-xs bg-blue-100 text-blue-700 px-2 py-0.5 rounded-full">
                        Mặc định
                      </span>
                    </div>
                    <span class="text-lg font-bold text-green-600">{{ formatPrice(variant.price) }}</span>
                  </div>
                </div>
              </div>

              <!-- Ingredients (if any) -->
              <div v-if="item.ingredients && item.ingredients.length > 0" class="bg-white rounded-lg p-3 mb-3">
                <div class="text-sm font-semibold text-gray-700 mb-2">🥘 Nguyên liệu:</div>
                <ul class="text-xs text-gray-600 space-y-1">
                  <li v-for="ingredient in item.ingredients" :key="ingredient.name">
                    • {{ ingredient.name }}: {{ ingredient.quantity }} {{ ingredient.unit }}
                  </li>
                </ul>
              </div>

              <!-- Quick Actions -->
              <div class="grid grid-cols-3 gap-2 pt-3 border-t">
                <button @click.stop="editItem(item)"
                  class="bg-yellow-500 text-white py-2 rounded-lg text-sm font-medium active:bg-yellow-600">
                  ✏️ Sửa
                </button>
                <button @click.stop="toggleAvailable(item)"
                  :class="item.available ? 'bg-gray-500 active:bg-gray-600' : 'bg-green-500 active:bg-green-600'"
                  class="text-white py-2 rounded-lg text-sm font-medium">
                  {{ item.available ? '🙈 Ẩn' : '👁️ Hiện' }}
                </button>
                <button @click.stop="deleteItem(item.id)"
                  class="bg-red-500 text-white py-2 rounded-lg text-sm font-medium active:bg-red-600">
                  🗑️ Xóa
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-if="filteredGroupedItems.length === 0" class="text-center py-16">
          <div class="text-6xl mb-4">📭</div>
          <p class="text-gray-500">{{ searchQuery ? 'Không tìm thấy món nào' : 'Chưa có món nào' }}</p>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />

    <!-- Category Management Modal - Mobile Optimized -->
    <transition name="slide-up">
      <div v-if="showCategoryModal" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full h-[85vh] flex flex-col">
          <!-- Fixed Header -->
          <div class="flex-shrink-0 bg-white px-4 py-4 border-b flex justify-between items-center rounded-t-3xl">
            <h3 class="text-lg font-bold">📁 Quản lý Danh mục</h3>
            <button @click="showCategoryModal = false" class="text-2xl text-gray-400">×</button>
          </div>
          
          <!-- Scrollable Content -->
          <div class="flex-1 overflow-y-auto px-4 py-4">
            <!-- Add New Category -->
            <div class="bg-gray-50 rounded-xl p-4 mb-4 flex-shrink-0">
              <h4 class="font-semibold text-gray-800 mb-3">Thêm danh mục mới</h4>
              <form @submit.prevent="addCategory" class="flex flex-col sm:flex-row gap-2">
                <input v-model="categoryForm.name" type="text" required placeholder="VD: Cà phê, Trà..." 
                  class="flex-1 px-4 py-3 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500" />
                <button type="submit" class="bg-purple-500 text-white px-6 py-3 rounded-lg font-medium text-base active:bg-purple-600 whitespace-nowrap">
                  Thêm
                </button>
              </form>
            </div>

            <!-- Category List -->
            <div class="space-y-3 pb-4">
              <!-- Debug info -->
              <div v-if="menuCategories.length === 0" class="text-center py-8 text-gray-500">
                <p>Chưa có danh mục nào</p>
                <p class="text-xs mt-2">Loading: {{ categoriesLoading }}</p>
              </div>
              
              <div v-for="(cat, idx) in menuCategories" :key="cat.id"
                class="bg-white border border-gray-200 rounded-xl p-4 flex items-center justify-between">
                <div class="flex items-center gap-3 flex-1 min-w-0">
                  <div class="w-12 h-12 rounded-lg flex items-center justify-center text-2xl flex-shrink-0" :class="getCategoryColor(cat.name)">
                    {{ getCategoryIcon(cat.name) }}
                  </div>
                  <div class="min-w-0">
                    <div class="font-medium text-gray-800 truncate">{{ cat.name }}</div>
                    <div class="text-xs text-gray-500">{{ getMenuCountByCategory(cat.name) }} món</div>
                  </div>
                </div>
                <div class="flex items-center gap-1 flex-shrink-0 ml-2">
                  <button @click="moveCategoryUp(idx)" :disabled="idx === 0"
                    class="w-8 h-8 flex items-center justify-center rounded-lg bg-gray-100 text-gray-600 disabled:opacity-30 active:bg-gray-200">
                    ↑
                  </button>
                  <button @click="moveCategoryDown(idx)" :disabled="idx === menuCategories.length - 1"
                    class="w-8 h-8 flex items-center justify-center rounded-lg bg-gray-100 text-gray-600 disabled:opacity-30 active:bg-gray-200">
                    ↓
                  </button>
                  <button @click="deleteCategory(cat.id, cat.name)" class="text-red-500 hover:text-red-700 p-2">
                    🗑️
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- Create/Edit Form Modal - Full Screen Mobile Optimized -->
    <transition name="slide-right">
      <div v-if="showMenuForm" class="fixed inset-0 bg-gray-50 z-50 flex flex-col">
        <!-- Mobile Header - Fixed with Safe Area -->
        <div class="sticky top-0 z-40 bg-gradient-to-r from-blue-500 to-cyan-500 shadow-lg flex-shrink-0">
          <div class="px-4 py-3" style="padding-top: max(0.75rem, env(safe-area-inset-top))">
            <div class="flex items-center justify-between">
              <button @click="cancelEdit" class="text-white text-2xl p-2 active:scale-95 transition-transform">
                ← 
              </button>
              <h1 class="text-lg font-bold text-white">
                {{ editingItem ? '✏️ Cập nhật món' : '➕ Thêm món mới' }}
              </h1>
              <button @click="saveItem" 
                class="bg-white text-blue-600 px-4 py-2 rounded-lg text-sm font-bold active:scale-95 transition-transform shadow-md">
                💾 Lưu
              </button>
            </div>
          </div>
        </div>

        <!-- Scrollable Content with Bottom Padding for Safe Area -->
        <div class="flex-1 overflow-y-auto" style="padding-bottom: max(1.5rem, env(safe-area-inset-bottom))">
          <form @submit.prevent="saveItem" class="p-4 space-y-4">
            
            <!-- Tên món -->
            <div class="bg-white rounded-xl p-4 shadow-sm">
              <label class="block text-sm font-bold text-gray-700 mb-2">📝 Tên món *</label>
              <input v-model="form.name" type="text" required placeholder="VD: Cà phê sữa đá"
                class="w-full px-4 py-3 text-base border-2 border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500" />
            </div>

            <!-- Danh mục -->
            <div class="bg-white rounded-xl p-4 shadow-sm">
              <label class="block text-sm font-bold text-gray-700 mb-2">📁 Danh mục *</label>
              <select v-model="form.category" required
                class="w-full px-4 py-3 text-base border-2 border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
                <option value="">Chọn danh mục</option>
                <option v-for="cat in suggestedCategories" :key="cat.id" :value="cat.name">{{ cat.name }}</option>
              </select>
            </div>

            <!-- Mô tả -->
            <div class="bg-white rounded-xl p-4 shadow-sm">
              <label class="block text-sm font-bold text-gray-700 mb-2">📄 Mô tả</label>
              <textarea v-model="form.description" rows="3" placeholder="Mô tả món ăn..."
                class="w-full px-4 py-3 text-base border-2 border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"></textarea>
            </div>

            <!-- Has Variants Toggle - Prominent -->
            <div class="bg-gradient-to-br from-purple-50 to-pink-50 rounded-xl p-4 shadow-sm border-2 border-purple-200">
              <label class="flex items-center justify-between cursor-pointer">
                <div class="flex items-center gap-3">
                  <div class="text-3xl">📏</div>
                  <div>
                    <div class="text-base font-bold text-gray-800">Có nhiều size</div>
                    <div class="text-xs text-gray-600">Size S, M, L, XL...</div>
                  </div>
                </div>
                <div class="relative">
                  <input type="checkbox" v-model="form.has_variants" @change="onHasVariantsChange"
                    class="sr-only peer" />
                  <div class="w-14 h-8 bg-gray-300 rounded-full peer peer-checked:bg-purple-500 transition-colors"></div>
                  <div class="absolute left-1 top-1 w-6 h-6 bg-white rounded-full transition-transform peer-checked:translate-x-6 shadow-md"></div>
                </div>
              </label>
            </div>

            <!-- Single Size Section -->
            <div v-if="!form.has_variants" class="space-y-4">
              <!-- Giá -->
              <div class="bg-white rounded-xl p-4 shadow-sm">
                <label class="block text-sm font-bold text-gray-700 mb-2">💰 Giá (VNĐ) *</label>
                <input v-model.number="form.price" type="number" min="0" step="1000" required placeholder="0"
                  class="w-full px-4 py-3 text-lg font-bold border-2 border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500" />
              </div>

              <!-- Nguyên liệu -->
              <div class="bg-white rounded-xl p-4 shadow-sm">
                <div class="flex items-center justify-between mb-3">
                  <label class="text-sm font-bold text-gray-700">🥘 Nguyên liệu</label>
                  <button type="button" @click="showIngredientSelector = true"
                    class="bg-blue-500 text-white px-4 py-2 rounded-lg text-sm font-bold active:scale-95 transition-transform shadow-md">
                    + Thêm
                  </button>
                </div>

                <div v-if="form.ingredients.length === 0" class="text-center py-8 text-gray-400">
                  <div class="text-4xl mb-2">🍽️</div>
                  <p class="text-sm">Chưa có nguyên liệu</p>
                </div>

                <div v-else class="space-y-3">
                  <div v-for="(ing, index) in form.ingredients" :key="index"
                    class="bg-gray-50 rounded-lg p-3 border-2 border-gray-200">
                    
                    <!-- Ingredient Header - Compact -->
                    <div class="flex justify-between items-start mb-2">
                      <div class="flex-1 min-w-0">
                        <div class="flex items-center gap-2">
                          <div class="font-bold text-sm text-gray-800 truncate">{{ ing.name || 'Không có tên' }}</div>
                          <span v-if="ing.type === 'batch'" class="text-xs bg-purple-100 text-purple-700 px-2 py-0.5 rounded-full flex-shrink-0">
                            🧪 Batch
                          </span>
                        </div>
                        <div class="text-xs text-gray-500">
                          <span v-if="ing.type === 'batch'">
                            Batch @ {{ formatPrice(ing.costPerUnit || 0) }}/{{ ing.unit }}
                          </span>
                          <span v-else>
                            Kho: {{ ing.stockUnit || ing.unit }} @ {{ formatPrice(ing.costPerUnit || 0) }}/{{ ing.stockUnit || ing.unit }}
                          </span>
                        </div>
                      </div>
                      <button type="button" @click="removeIngredient(index)"
                        class="ml-2 bg-red-500 text-white px-2 py-1 rounded text-xs font-bold active:scale-95">
                        ×
                      </button>
                    </div>

                    <!-- Recipe Unit & Quantity - Inline -->
                    <div class="grid grid-cols-2 gap-2 mb-2">
                      <div>
                        <label class="text-xs text-gray-600 block mb-1">Đơn vị:</label>
                        <select v-model="ing.unit" @change="updateRecipeUnit(index)"
                          class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded-lg">
                          <option v-for="unit in ing.compatibleUnits" :key="unit" :value="unit">
                            {{ unit }}
                          </option>
                        </select>
                      </div>
                      <div>
                        <label class="text-xs text-gray-600 block mb-1">Số lượng:</label>
                        <input v-model.number="ing.quantity" 
                          @input="updateIngredientCost(index)"
                          type="number" min="0" step="0.1" placeholder="0" required
                          class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded-lg" />
                      </div>
                    </div>
                    
                    <!-- Conversion Info -->
                    <div v-if="ing.conversionRate !== 1" 
                      class="mb-2 p-2 bg-blue-50 rounded text-xs text-blue-700">
                      <span class="font-bold">ℹ️</span> {{ getConversionExplanation(ing.stockUnit, ing.unit) }}
                    </div>
                    
                    <!-- Cost Preview - Compact -->
                    <div v-if="ing.costPerUnit > 0" class="p-2 bg-green-50 rounded-lg border border-green-200">
                      <div class="flex justify-between items-center">
                        <span class="text-xs text-green-700 font-bold">Chi phí:</span>
                        <span class="text-sm font-bold text-green-700">
                          {{ formatPrice(ing.estimatedCost) }}
                        </span>
                      </div>
                      <div v-if="ing.wastage > 0" class="text-xs text-green-600 mt-1">
                        (+ {{ ing.wastage }}% hao hụt)
                      </div>
                    </div>

                    <!-- Deduct Inventory Toggle -->
                    <div class="mt-2 flex items-center gap-1">
                      <button type="button"
                        @click="ing.deductInventory = !ing.deductInventory"
                        :class="[
                          'px-2 py-1 rounded text-xs font-medium transition-colors',
                          ing.deductInventory
                            ? 'bg-red-100 text-red-700'
                            : 'bg-gray-100 text-gray-500'
                        ]"
                        :title="ing.deductInventory ? 'Đang trừ tồn kho khi order PAID' : 'Chỉ tính cost, không trừ kho'"
                      >
                        {{ ing.deductInventory ? '📦 Trừ kho' : '💰 Chỉ cost' }}
                      </button>
                    </div>
                  </div>

                  <!-- Total Cost Summary -->
                  <div v-if="totalIngredientCost > 0" class="bg-gradient-to-r from-yellow-500 to-orange-500 rounded-xl p-4 text-white shadow-lg">
                    <div class="flex justify-between items-center">
                      <span class="text-sm font-bold">💰 Tổng chi phí nguyên liệu:</span>
                      <span class="text-xl font-bold">{{ formatPrice(totalIngredientCost) }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Multi-Size Section - Mobile Optimized -->
            <div v-else class="space-y-4">
              <!-- Variants Header with Add Button -->
              <div class="bg-gradient-to-r from-purple-500 to-pink-500 rounded-xl p-4 shadow-lg">
                <div class="flex items-center justify-between mb-3">
                  <div class="text-white">
                    <div class="text-lg font-bold">📏 Các size</div>
                    <div class="text-xs opacity-90">{{ form.variants.length }} size</div>
                  </div>
                  <button type="button" @click="addVariant"
                    class="bg-white text-purple-600 px-4 py-2 rounded-lg text-sm font-bold active:scale-95 transition-transform shadow-md">
                    + Thêm size
                  </button>
                </div>
              </div>

              <!-- Validation Errors -->
              <div v-if="variantValidationErrors.length > 0" class="bg-red-50 border-2 border-red-200 rounded-xl p-4">
                <div class="text-sm font-bold text-red-800 mb-2">⚠️ Lỗi:</div>
                <ul class="text-sm text-red-700 space-y-1">
                  <li v-for="(error, index) in variantValidationErrors" :key="index">• {{ error }}</li>
                </ul>
              </div>

              <!-- Empty State -->
              <div v-if="form.variants.length === 0" class="bg-white rounded-xl p-8 text-center shadow-sm">
                <div class="text-6xl mb-3">📏</div>
                <p class="text-gray-500 mb-4">Chưa có size nào</p>
                <button type="button" @click="addVariant"
                  class="bg-purple-500 text-white px-6 py-3 rounded-lg font-bold active:scale-95 transition-transform shadow-md">
                  + Thêm size đầu tiên
                </button>
              </div>

              <!-- Variants List - Card Style -->
              <div v-else class="space-y-4">
                <div v-for="(variant, vIndex) in form.variants" :key="vIndex"
                  class="bg-white rounded-xl shadow-md overflow-hidden border-2"
                  :class="variant.is_default ? 'border-purple-400' : 'border-gray-200'">
                  
                  <!-- Variant Header - Colorful -->
                  <div class="bg-gradient-to-r p-4"
                    :class="variant.is_default ? 'from-purple-500 to-pink-500' : 'from-gray-400 to-gray-500'">
                    <div class="flex items-center justify-between">
                      <div class="text-white">
                        <div class="text-lg font-bold">{{ variant.name || `Size ${vIndex + 1}` }}</div>
                        <div class="text-xs opacity-90">{{ formatPrice(variant.price) }}</div>
                      </div>
                      <button type="button" @click="removeVariant(vIndex)"
                        :disabled="form.variants.length === 1"
                        class="bg-white bg-opacity-20 text-white px-3 py-2 rounded-lg text-sm font-bold active:scale-95 transition-transform disabled:opacity-50 disabled:cursor-not-allowed">
                        🗑️
                      </button>
                    </div>
                  </div>

                  <!-- Variant Form - Compact -->
                  <div class="p-4 space-y-3">
                    <!-- Name & Price -->
                    <div class="space-y-3">
                      <div>
                        <label class="block text-xs font-bold text-gray-700 mb-1">Tên size *</label>
                        <input v-model="variant.name" type="text" required placeholder="Size M"
                          class="w-full px-3 py-2 text-base border-2 border-gray-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500" />
                      </div>
                      <div>
                        <label class="block text-xs font-bold text-gray-700 mb-1">Giá (VNĐ) *</label>
                        <input v-model.number="variant.price" type="number" min="0" step="1000" required placeholder="0"
                          class="w-full px-3 py-2 text-lg font-bold border-2 border-gray-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500" />
                      </div>
                    </div>

                    <!-- Is Default Toggle -->
                    <div class="bg-purple-50 rounded-lg p-3 border border-purple-200">
                      <label class="flex items-center gap-3 cursor-pointer">
                        <input type="checkbox" :checked="variant.is_default" 
                          @change="setDefaultVariant(vIndex)"
                          class="w-5 h-5 text-purple-600 rounded focus:ring-2 focus:ring-purple-500" />
                        <div class="flex-1">
                          <div class="text-sm font-bold text-gray-800">Mặc định</div>
                          <div class="text-xs text-gray-600">Size được chọn tự động</div>
                        </div>
                      </label>
                    </div>

                    <!-- Variant Ingredients -->
                    <div class="border-t-2 border-gray-100 pt-3">
                      <div class="flex items-center justify-between mb-3">
                        <label class="text-xs font-bold text-gray-700">🥘 Nguyên liệu</label>
                        <button type="button" @click="openVariantIngredientSelector(vIndex)"
                          class="bg-blue-500 text-white px-3 py-1.5 rounded-lg text-xs font-bold active:scale-95 transition-transform">
                          + Thêm
                        </button>
                      </div>

                      <div v-if="variant.ingredients.length === 0" class="text-center py-6 text-gray-400 bg-gray-50 rounded-lg">
                        <div class="text-3xl mb-1">🍽️</div>
                        <p class="text-xs">Chưa có nguyên liệu</p>
                      </div>

                      <div v-else class="space-y-2">
                        <div v-for="(ing, iIndex) in variant.ingredients" :key="iIndex"
                          class="bg-gray-50 rounded-lg p-3 border border-gray-200">
                          
                          <!-- Ingredient Header - Compact -->
                          <div class="flex justify-between items-start mb-2">
                            <div class="flex-1 min-w-0">
                              <div class="flex items-center gap-2">
                                <div class="font-bold text-sm text-gray-800 truncate">{{ ing.name || 'Không có tên' }}</div>
                                <span v-if="ing.type === 'batch'" class="text-xs bg-purple-100 text-purple-700 px-2 py-0.5 rounded-full flex-shrink-0">
                                  🧪 Batch
                                </span>
                              </div>
                              <div class="text-xs text-gray-500">
                                <span v-if="ing.type === 'batch'">
                                  Batch @ {{ formatPrice(ing.costPerUnit || 0) }}/{{ ing.unit }}
                                </span>
                                <span v-else>
                                  Kho: {{ ing.stockUnit || ing.unit }} @ {{ formatPrice(ing.costPerUnit || 0) }}/{{ ing.stockUnit || ing.unit }}
                                </span>
                              </div>
                            </div>
                            <button type="button" @click="removeVariantIngredient(vIndex, iIndex)"
                              class="ml-2 bg-red-500 text-white px-2 py-1 rounded text-xs font-bold active:scale-95">
                              ×
                            </button>
                          </div>

                          <!-- Recipe Unit & Quantity - Inline -->
                          <div class="grid grid-cols-2 gap-2 mb-2">
                            <div>
                              <label class="text-xs text-gray-600 block mb-1">Đơn vị:</label>
                              <select v-model="ing.unit" @change="updateVariantRecipeUnit(vIndex, iIndex)"
                                class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded-lg">
                                <option v-for="unit in ing.compatibleUnits" :key="unit" :value="unit">
                                  {{ unit }}
                                </option>
                              </select>
                            </div>
                            <div>
                              <label class="text-xs text-gray-600 block mb-1">Số lượng:</label>
                              <input v-model.number="ing.quantity" 
                                @input="updateVariantIngredientCost(vIndex, iIndex)"
                                type="number" min="0" step="0.1" placeholder="0" required
                                class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded-lg" />
                            </div>
                          </div>
                          
                          <!-- Conversion Info -->
                          <div v-if="ing.conversionRate !== 1" 
                            class="mb-2 p-2 bg-blue-50 rounded text-xs text-blue-700">
                            <span class="font-bold">ℹ️</span> {{ getConversionExplanation(ing.stockUnit, ing.unit) }}
                          </div>
                          
                          <!-- Cost Preview - Compact -->
                          <div v-if="ing.costPerUnit > 0" class="p-2 bg-green-50 rounded-lg border border-green-200">
                            <div class="flex justify-between items-center">
                              <span class="text-xs text-green-700 font-bold">Chi phí:</span>
                              <span class="text-sm font-bold text-green-700">
                                {{ formatPrice(ing.estimatedCost) }}
                              </span>
                            </div>
                            <div v-if="ing.wastage > 0" class="text-xs text-green-600 mt-1">
                              (+ {{ ing.wastage }}% hao hụt)
                            </div>
                          </div>

                          <!-- Deduct Inventory Toggle -->
                          <div class="mt-2 flex items-center gap-1">
                            <button type="button"
                              @click="ing.deductInventory = !ing.deductInventory"
                              :class="[
                                'px-2 py-1 rounded text-xs font-medium transition-colors',
                                ing.deductInventory
                                  ? 'bg-red-100 text-red-700'
                                  : 'bg-gray-100 text-gray-500'
                              ]"
                              :title="ing.deductInventory ? 'Đang trừ tồn kho khi order PAID' : 'Chỉ tính cost, không trừ kho'"
                            >
                              {{ ing.deductInventory ? '📦 Trừ kho' : '💰 Chỉ cost' }}
                            </button>
                          </div>
                        </div>
                        
                        <!-- Total Cost for Variant -->
                        <div v-if="getVariantTotalCost(vIndex) > 0" 
                          class="bg-purple-50 border-2 border-purple-300 rounded-lg p-3">
                          <div class="flex justify-between items-center">
                            <span class="text-sm font-bold text-purple-800">💰 Tổng chi phí:</span>
                            <span class="text-lg font-bold text-purple-900">{{ formatPrice(getVariantTotalCost(vIndex)) }}</span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Submit Button - Sticky at Bottom -->
            <div class="sticky bottom-0 left-0 right-0 bg-gradient-to-t from-gray-50 via-gray-50 to-transparent pt-4 pb-2">
              <button type="submit"
                class="w-full bg-gradient-to-r from-blue-500 to-cyan-500 text-white py-4 rounded-xl text-lg font-bold shadow-lg active:scale-98 transition-transform">
                💾 {{ editingItem ? 'Cập nhật món' : 'Thêm món mới' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>

    <!-- Ingredient/Batch Selector Modal -->
    <transition name="slide-up">
      <div v-if="showIngredientSelector" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full h-[85vh] flex flex-col">
          <!-- Fixed Header -->
          <div class="flex-shrink-0 bg-white px-4 py-4 border-b rounded-t-3xl">
            <div class="flex justify-between items-center mb-3">
              <h3 class="text-lg font-bold">{{ ingredientType === 'raw' ? '🥬 Chọn nguyên liệu' : '🧪 Chọn Batch' }}</h3>
              <button @click="showIngredientSelector = false" class="text-2xl text-gray-400">×</button>
            </div>
            
            <!-- Type Toggle -->
            <div class="flex gap-2 bg-gray-100 p-1 rounded-lg">
              <button
                @click="ingredientType = 'raw'"
                :class="[
                  'flex-1 py-2 px-4 rounded-md text-sm font-medium transition-colors',
                  ingredientType === 'raw' 
                    ? 'bg-white text-blue-600 shadow-sm' 
                    : 'text-gray-600'
                ]"
              >
                🥬 Nguyên liệu thô
              </button>
              <button
                @click="ingredientType = 'batch'"
                :class="[
                  'flex-1 py-2 px-4 rounded-md text-sm font-medium transition-colors',
                  ingredientType === 'batch' 
                    ? 'bg-white text-purple-600 shadow-sm' 
                    : 'text-gray-600'
                ]"
              >
                🧪 Batch
              </button>
            </div>
          </div>
          
          <!-- Search -->
          <div class="flex-shrink-0 px-4 py-3 border-b">
            <input
              v-model="ingredientSearchQuery"
              type="text"
              :placeholder="ingredientType === 'raw' ? 'Tìm kiếm nguyên liệu...' : 'Tìm kiếm batch...'"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
            />
          </div>
          
          <!-- Scrollable Content -->
          <div class="flex-1 overflow-y-auto px-4 py-4">
            <!-- Raw Ingredients List -->
            <div v-if="ingredientType === 'raw'">
              <div v-if="ingredientsLoading" class="text-center py-8">
                <div class="text-4xl mb-2">⏳</div>
                <p class="text-gray-500">Đang tải...</p>
              </div>
              
              <div v-else-if="filteredAvailableIngredients.length === 0" class="text-center py-8">
                <div class="text-4xl mb-2">📭</div>
                <p class="text-gray-500">Không có nguyên liệu nào</p>
              </div>
              
              <div v-else class="space-y-2">
                <button
                  v-for="ingredient in filteredAvailableIngredients"
                  :key="ingredient.id"
                  @click="selectIngredient(ingredient)"
                  :disabled="isIngredientSelected(ingredient.id)"
                  class="w-full bg-white border border-gray-200 rounded-xl p-4 text-left hover:border-blue-500 hover:bg-blue-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <div class="flex items-center justify-between">
                    <div class="flex-1">
                      <div class="font-medium text-gray-800">{{ ingredient.name }}</div>
                      <div class="text-sm text-gray-500">{{ ingredient.category }}</div>
                      <div class="text-xs text-gray-400 mt-1">
                        Tồn kho: {{ ingredient.current_stock }} {{ ingredient.unit }}
                      </div>
                    </div>
                    <div v-if="isIngredientSelected(ingredient.id)" class="text-green-500 text-xl">✓</div>
                    <div v-else class="text-blue-500 text-xl">+</div>
                  </div>
                </button>
              </div>
            </div>

            <!-- Batch List -->
            <div v-else>
              <div v-if="batchesLoading" class="text-center py-8">
                <div class="text-4xl mb-2">⏳</div>
                <p class="text-gray-500">Đang tải...</p>
              </div>
              
              <div v-else-if="filteredAvailableBatches.length === 0" class="text-center py-8">
                <div class="text-4xl mb-2">📭</div>
                <p class="text-gray-500">Không có batch nào</p>
              </div>
              
              <div v-else class="space-y-2">
                <button
                  v-for="batch in filteredAvailableBatches"
                  :key="batch.id"
                  @click="selectBatch(batch)"
                  :disabled="isBatchSelected(batch.id)"
                  class="w-full bg-white border border-gray-200 rounded-xl p-4 text-left hover:border-purple-500 hover:bg-purple-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <div class="flex items-center justify-between">
                    <div class="flex-1">
                      <div class="flex items-center gap-2">
                        <div class="font-medium text-gray-800">{{ batch.name }}</div>
                        <span class="text-xs bg-purple-100 text-purple-700 px-2 py-0.5 rounded-full">
                          Batch
                        </span>
                      </div>
                      <div class="text-sm text-gray-500 mt-1">
                        Đơn vị: {{ batch.unit }}
                      </div>
                      <div class="text-xs text-gray-400 mt-1">
                        Thời hạn: {{ batch.shelf_life_hours }} giờ
                      </div>
                    </div>
                    <div v-if="isBatchSelected(batch.id)" class="text-green-500 text-xl">✓</div>
                    <div v-else class="text-purple-500 text-xl">+</div>
                  </div>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useMenuStore } from '../stores/menu'
import { useIngredientStore } from '../stores/ingredient'
import { useBatchDefinitionStore } from '../stores/batchDefinition'
import { useBatchRecordStore } from '../stores/batchRecord'
import { useUnitConversion } from '../composables/useUnitConversion'
import { menuCategoryService } from '../services/menuCategory'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'

const menuStore = useMenuStore()
const ingredientStore = useIngredientStore()
const batchDefinitionStore = useBatchDefinitionStore()
const batchRecordStore = useBatchRecordStore()

// Unit conversion composable
const { 
  getConversionRate, 
  isValidConversion, 
  getCompatibleUnits,
  calculateCostBreakdown,
  getConversionExplanation 
} = useUnitConversion()

const searchQuery = ref('')
const showMenuForm = ref(false)
const showCategoryModal = ref(false)
const showIngredientSelector = ref(false)
const ingredientSearchQuery = ref('')
const ingredientType = ref('raw') // 'raw' or 'batch'
const editingItem = ref(null)
const form = ref({
  name: '',
  category: '',
  price: 0,
  description: '',
  ingredients: [],
  has_variants: false,
  variants: []
})

const categoryForm = ref({
  name: ''
})

// Variant validation errors
const variantValidationErrors = ref([])

// Current variant index for ingredient selector
const currentVariantIndex = ref(null)

const items = computed(() => menuStore.items)
const loading = computed(() => menuStore.loading)
const error = computed(() => menuStore.error)

// Ingredients
const availableIngredients = computed(() => ingredientStore.items)
const ingredientsLoading = computed(() => ingredientStore.loading)

// Batches
const availableBatchDefinitions = computed(() => batchDefinitionStore.definitions)
const availableBatchRecords = computed(() => batchRecordStore.records)
const batchesLoading = computed(() => batchDefinitionStore.loading || batchRecordStore.loading)

// Filter available batches based on search and already selected
const filteredAvailableBatches = computed(() => {
  // Use batch definitions instead of batch records
  let filtered = availableBatchDefinitions.value || []
  
  if (ingredientSearchQuery.value) {
    const query = ingredientSearchQuery.value.toLowerCase()
    filtered = filtered.filter(batch => 
      batch.name?.toLowerCase().includes(query)
    )
  }
  
  return filtered
})

// Menu categories from API
const menuCategories = ref([])
const categoriesLoading = ref(false)

// Suggested categories for dropdown (uses API categories)
const suggestedCategories = computed(() => {
  return menuCategories.value
})

// Filter available ingredients based on search and already selected
const filteredAvailableIngredients = computed(() => {
  let filtered = availableIngredients.value || []
  
  if (ingredientSearchQuery.value) {
    const query = ingredientSearchQuery.value.toLowerCase()
    filtered = filtered.filter(ing => 
      ing.name?.toLowerCase().includes(query) ||
      ing.category?.toLowerCase().includes(query)
    )
  }
  
  return filtered
})

const groupedItems = computed(() => {
  const groups = {}
  if (!items.value || !Array.isArray(items.value)) {
    return []
  }
  
  // Sort items by created_at (newest first)
  const sortedItems = [...items.value].sort((a, b) => {
    const dateA = new Date(a.created_at || 0)
    const dateB = new Date(b.created_at || 0)
    return dateB - dateA // Newest first
  })
  
  sortedItems.forEach(item => {
    if (!groups[item.category]) {
      groups[item.category] = {
        name: item.category,
        items: []
      }
    }
    groups[item.category].items.push(item)
  })
  return Object.values(groups)
})

const filteredGroupedItems = computed(() => {
  if (!searchQuery.value) return groupedItems.value || []
  
  const query = searchQuery.value.toLowerCase()
  const grouped = groupedItems.value || []
  return grouped
    .map(category => ({
      ...category,
      items: (category.items || []).filter(item => 
        item.name?.toLowerCase().includes(query) ||
        item.description?.toLowerCase().includes(query) ||
        item.category?.toLowerCase().includes(query)
      )
    }))
    .filter(category => category.items.length > 0)
})

// Fetch categories from API
const fetchCategories = async () => {
  try {
    categoriesLoading.value = true
    const response = await menuCategoryService.getCategories()
    console.log('Raw API response:', response)
    
    // Handle both cases: direct array or wrapped in { data: [...] }
    if (Array.isArray(response)) {
      menuCategories.value = response
    } else if (response && Array.isArray(response.data)) {
      menuCategories.value = response.data
    } else {
      menuCategories.value = []
    }
    
    console.log('Parsed categories:', menuCategories.value)
  } catch (err) {
    console.error('Failed to fetch categories:', err)
    menuCategories.value = []
  } finally {
    categoriesLoading.value = false
  }
}

// Refresh data function
const refreshData = async () => {
  await Promise.all([
    menuStore.fetchMenuItems(),
    fetchCategories(),
    ingredientStore.fetchIngredients(),
    batchDefinitionStore.fetchDefinitions(),
    batchRecordStore.fetchRecords()
  ])
}

// Pull to refresh
const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)

onMounted(() => {
  refreshData()
})

const formatPrice = (price) => {
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND'
  }).format(price)
}

const openCreateModal = () => {
  form.value = {
    name: '',
    category: '',
    price: 0,
    description: '',
    ingredients: [],
    has_variants: false,
    variants: []
  }
  editingItem.value = null
  variantValidationErrors.value = []
  showMenuForm.value = true
}

const viewDetails = (item) => {
  editItem(item)
}

const editItem = (item) => {
  editingItem.value = item
  
  // Prepare ingredients for single-size items
  const preparedIngredients = item.ingredients ? item.ingredients.map(ing => {
    // Check if this is a batch ingredient
    const isBatch = ing.ingredient_type === 'batch' || ing.type === 'batch'
    
    if (isBatch) {
      // Handle batch ingredient
      const batchData = availableBatchDefinitions.value.find(b => b.id === ing.batch_id || b.name === ing.name)
      if (batchData) {
        const batchCostPerUnit = calculateBatchCostPerUnit(batchData)
        const compatibleUnits = getCompatibleUnits(batchData.unit)
        const conversionRate = getConversionRate(batchData.unit, ing.unit)
        
        // Calculate estimated cost
        const breakdown = calculateCostBreakdown(
          ing.quantity,
          ing.unit,
          batchCostPerUnit,
          batchData.unit,
          0 // Batches don't have wastage
        )
        
        return {
          id: batchData.id,
          batch_definition_id: batchData.id,
          type: 'batch',
          name: ing.name,
          quantity: ing.quantity,
          unit: ing.unit,
          stockUnit: batchData.unit,
          compatibleUnits: compatibleUnits,
          costPerUnit: batchCostPerUnit,
          wastage: 0,
          conversionRate: conversionRate,
          estimatedCost: breakdown.totalCost,
          deductInventory: ing.deduct_inventory ?? true
        }
      }
    } else {
      // Handle raw ingredient
      const ingredientData = availableIngredients.value.find(i => i.name === ing.name || i.id === ing.ingredient_id)
      if (ingredientData) {
        const compatibleUnits = getCompatibleUnits(ingredientData.unit)
        const conversionRate = getConversionRate(ingredientData.unit, ing.unit)

        // Calculate estimated cost
        const breakdown = calculateCostBreakdown(
          ing.quantity,
          ing.unit,
          ingredientData.cost_per_unit || 0,
          ingredientData.unit,
          ingredientData.wastage_percentage || 0
        )

        return {
          id: ingredientData.id,
          type: 'raw',
          name: ing.name,
          quantity: ing.quantity,
          unit: ing.unit,
          stockUnit: ingredientData.unit,
          compatibleUnits: compatibleUnits,
          costPerUnit: ingredientData.cost_per_unit || 0,
          wastage: ingredientData.wastage_percentage || 0,
          conversionRate: conversionRate,
          estimatedCost: breakdown.totalCost,
          deductInventory: ing.deduct_inventory ?? true
        }
      }
    }
    return ing
  }) : []

  // Prepare variants with enriched ingredient data
  const preparedVariants = item.variants ? item.variants.map(variant => {
    const enrichedIngredients = variant.ingredients ? variant.ingredients.map(ing => {
      // Check if this is a batch ingredient
      const isBatch = ing.ingredient_type === 'batch' || ing.type === 'batch'

      if (isBatch) {
        // Handle batch ingredient
        const batchData = availableBatchDefinitions.value.find(b => b.id === ing.batch_id || b.name === ing.name)
        if (batchData) {
          const batchCostPerUnit = calculateBatchCostPerUnit(batchData)
          const compatibleUnits = getCompatibleUnits(batchData.unit)
          const conversionRate = getConversionRate(batchData.unit, ing.unit)

          // Calculate estimated cost
          const breakdown = calculateCostBreakdown(
            ing.quantity,
            ing.unit,
            batchCostPerUnit,
            batchData.unit,
            0 // Batches don't have wastage
          )

          return {
            id: batchData.id,
            batch_definition_id: batchData.id,
            type: 'batch',
            name: ing.name,
            quantity: ing.quantity,
            unit: ing.unit,
            stockUnit: batchData.unit,
            compatibleUnits: compatibleUnits,
            costPerUnit: batchCostPerUnit,
            wastage: 0,
            conversionRate: conversionRate,
            estimatedCost: breakdown.totalCost,
            deductInventory: ing.deduct_inventory ?? true
          }
        }
      } else {
        // Handle raw ingredient
        const ingredientData = availableIngredients.value.find(i => i.name === ing.name || i.id === ing.ingredient_id)
        if (ingredientData) {
          const compatibleUnits = getCompatibleUnits(ingredientData.unit)
          const conversionRate = getConversionRate(ingredientData.unit, ing.unit)

          // Calculate estimated cost
          const breakdown = calculateCostBreakdown(
            ing.quantity,
            ing.unit,
            ingredientData.cost_per_unit || 0,
            ingredientData.unit,
            ingredientData.wastage_percentage || 0
          )

          return {
            id: ingredientData.id,
            type: 'raw',
            name: ing.name,
            quantity: ing.quantity,
            unit: ing.unit,
            stockUnit: ingredientData.unit,
            compatibleUnits: compatibleUnits,
            costPerUnit: ingredientData.cost_per_unit || 0,
            wastage: ingredientData.wastage_percentage || 0,
            conversionRate: conversionRate,
            estimatedCost: breakdown.totalCost,
            deductInventory: ing.deduct_inventory ?? true
          }
        }
      }
      return ing
    }) : []
    
    return {
      ...variant,
      ingredients: enrichedIngredients
    }
  }) : []
  
  form.value = {
    name: item.name,
    category: item.category,
    price: item.price || 0,
    description: item.description,
    ingredients: preparedIngredients,
    has_variants: item.has_variants || false,
    variants: preparedVariants
  }
  variantValidationErrors.value = []
  showMenuForm.value = true
}

const cancelEdit = () => {
  showMenuForm.value = false
  editingItem.value = null
  variantValidationErrors.value = []
  form.value = { 
    name: '', 
    category: '', 
    price: 0, 
    description: '', 
    ingredients: [],
    has_variants: false,
    variants: []
  }
}

// Handle has_variants toggle
const onHasVariantsChange = () => {
  if (form.value.has_variants) {
    // Switching to multi-size: clear price and ingredients
    form.value.price = 0
    form.value.ingredients = []
    
    // Add default variant if none exist
    if (form.value.variants.length === 0) {
      addVariant()
    }
  } else {
    // Switching to single-size: clear variants
    form.value.variants = []
  }
  variantValidationErrors.value = []
}

// Generate unique variant ID
const generateVariantId = () => {
  // Use timestamp + random number for uniqueness
  const timestamp = Date.now().toString(36) // Convert to base36 for shorter string
  const random = Math.random().toString(36).substring(2, 7) // 5 random chars
  return `${timestamp}-${random}`
}

// Add new variant
const addVariant = () => {
  const variantNumber = form.value.variants.length + 1
  const isFirst = form.value.variants.length === 0
  
  form.value.variants.push({
    id: generateVariantId(), // Auto-generate ID
    name: `Size ${variantNumber}`,
    price: 0,
    ingredients: [],
    available: true,
    is_default: isFirst // First variant is default
  })
}

// Remove variant
const removeVariant = (index) => {
  if (form.value.variants.length === 1) {
    alert('Phải có ít nhất 1 size!')
    return
  }
  
  const wasDefault = form.value.variants[index].is_default
  form.value.variants.splice(index, 1)
  
  // If removed variant was default, set first variant as default
  if (wasDefault && form.value.variants.length > 0) {
    form.value.variants[0].is_default = true
  }
}

// Set default variant
const setDefaultVariant = (index) => {
  // Unset all defaults
  form.value.variants.forEach((v, i) => {
    v.is_default = (i === index)
  })
}

// Open ingredient selector for specific variant
const openVariantIngredientSelector = (variantIndex) => {
  currentVariantIndex.value = variantIndex
  showIngredientSelector.value = true
}

// Remove ingredient from variant
const removeVariantIngredient = (variantIndex, ingredientIndex) => {
  form.value.variants[variantIndex].ingredients.splice(ingredientIndex, 1)
}

// Validate variants before save
const validateVariants = () => {
  const errors = []
  
  if (form.value.has_variants) {
    // Must have at least 1 variant
    if (form.value.variants.length === 0) {
      errors.push('Phải có ít nhất 1 size')
      return errors
    }
    
    // Check for duplicate IDs
    const ids = form.value.variants.map(v => v.id.trim()).filter(id => id)
    const uniqueIds = new Set(ids)
    if (ids.length !== uniqueIds.size) {
      errors.push('ID các size không được trùng nhau')
    }
    
    // Check each variant
    let hasDefault = false
    form.value.variants.forEach((variant, index) => {
      if (!variant.id || !variant.id.trim()) {
        errors.push(`Size ${index + 1}: ID không được để trống`)
      }
      if (!variant.name || !variant.name.trim()) {
        errors.push(`Size ${index + 1}: Tên không được để trống`)
      }
      if (!variant.price || variant.price <= 0) {
        errors.push(`Size ${index + 1}: Giá phải lớn hơn 0`)
      }
      if (variant.is_default) {
        hasDefault = true
      }
      // Validate ingredient quantities
      if (variant.ingredients && variant.ingredients.length > 0) {
        variant.ingredients.forEach((ing, iIndex) => {
          if (!ing.quantity || ing.quantity <= 0) {
            errors.push(`Size ${index + 1} - Nguyên liệu "${ing.name || iIndex + 1}": Số lượng phải lớn hơn 0`)
          }
        })
      }
    })
    
    // Must have exactly 1 default
    if (!hasDefault) {
      errors.push('Phải có 1 size mặc định')
    }
  }
  
  return errors
}

const selectIngredient = (ingredient) => {
  // Check if selecting for variant or single-size
  if (currentVariantIndex.value !== null) {
    // Adding to variant
    const variant = form.value.variants[currentVariantIndex.value]
    
    // Check if already selected for this variant
    if (variant.ingredients.some(ing => ing.id === ingredient.id && ing.type === 'raw')) {
      return
    }
    
    // Get compatible units
    const compatibleUnits = getCompatibleUnits(ingredient.unit)
    const recipeUnit = ingredient.unit
    const conversionRate = getConversionRate(ingredient.unit, recipeUnit)
    
    // Add ingredient to variant
    variant.ingredients.push({
      id: ingredient.id,
      type: 'raw', // Mark as raw ingredient
      name: ingredient.name,
      quantity: 1,
      unit: recipeUnit,
      stockUnit: ingredient.unit,
      compatibleUnits: compatibleUnits,
      costPerUnit: ingredient.cost_per_unit || 0,
      wastage: ingredient.wastage_percentage || 0,
      conversionRate: conversionRate,
      estimatedCost: 0,
      deductInventory: true
    })
    
    // Calculate initial cost
    const ingIndex = variant.ingredients.length - 1
    updateVariantIngredientCost(currentVariantIndex.value, ingIndex)
    
    // Reset and close
    currentVariantIndex.value = null
    showIngredientSelector.value = false
    ingredientSearchQuery.value = ''
    ingredientType.value = 'raw'
  } else {
    // Adding to single-size item (existing logic)
    if (isIngredientSelected(ingredient.id)) {
      return
    }
    
    const compatibleUnits = getCompatibleUnits(ingredient.unit)
    const recipeUnit = ingredient.unit
    const conversionRate = getConversionRate(ingredient.unit, recipeUnit)
    
    form.value.ingredients.push({
      id: ingredient.id,
      type: 'raw', // Mark as raw ingredient
      name: ingredient.name,
      quantity: 1,
      unit: recipeUnit,
      stockUnit: ingredient.unit,
      compatibleUnits: compatibleUnits,
      costPerUnit: ingredient.cost_per_unit || 0,
      wastage: ingredient.wastage_percentage || 0,
      conversionRate: conversionRate,
      estimatedCost: 0,
      deductInventory: true
    })
    
    updateIngredientCost(form.value.ingredients.length - 1)
    showIngredientSelector.value = false
    ingredientSearchQuery.value = ''
    ingredientType.value = 'raw'
  }
}

// Calculate batch cost per unit based on conversion rates
const calculateBatchCostPerUnit = (batch) => {
  if (!batch.conversion_rates || !Array.isArray(batch.conversion_rates) || batch.conversion_rates.length === 0) {
    console.warn('⚠️ No conversion_rates found for batch:', batch.name)
    return 0
  }
  
  let totalCost = 0
  
  // Sum up costs from all source ingredients
  for (const conversionRate of batch.conversion_rates) {
    // Find the source ingredient
    const ingredientId = conversionRate.source_ingredient_id || conversionRate.ingredient_id
    const sourceIngredient = availableIngredients.value.find(ing => ing.id === ingredientId)
    
    if (!sourceIngredient) {
      console.warn(`❌ Source ingredient not found: ${ingredientId}`)
      continue
    }
    
    // Calculate cost for this ingredient
    const sourceQuantity = conversionRate.source_quantity || 0
    const sourceUnit = conversionRate.source_unit || sourceIngredient.unit
    const ingredientStockUnit = sourceIngredient.unit
    const wastageRate = conversionRate.wastage_rate || 0 // Already in decimal format (0.1 = 10%)
    const costPerUnit = sourceIngredient.cost_per_unit || 0
    
    // Convert source_quantity from sourceUnit to ingredient stock unit
    // Example: 200g → 0.2kg
    const unitConversionRate = getConversionRate(sourceUnit, ingredientStockUnit)
    const quantityInStockUnit = sourceQuantity * unitConversionRate
    
    // Calculate: (quantity_in_stock_unit * (1 + wastage_rate)) * cost_per_unit
    const ingredientCost = quantityInStockUnit * (1 + wastageRate) * costPerUnit
    
    totalCost += ingredientCost
  }
  
  // Divide by batch quantity to get cost per batch output unit
  const batchQuantity = batch.conversion_rates[0]?.batch_quantity || batch.batch_quantity || 1
  const costPerBatchUnit = totalCost / batchQuantity
  
  return costPerBatchUnit
}

const selectBatch = (batch) => {
  // batch is now a batch definition, not a batch record
  // Calculate cost per unit for this batch
  const batchCostPerUnit = calculateBatchCostPerUnit(batch)
  
  // Get compatible units for batch (allows unit conversion like L→ml, kg→g)
  const compatibleUnits = getCompatibleUnits(batch.unit)
  
  // Check if selecting for variant or single-size
  if (currentVariantIndex.value !== null) {
    // Adding to variant
    const variant = form.value.variants[currentVariantIndex.value]
    
    // Check if already selected for this variant
    if (variant.ingredients.some(ing => ing.id === batch.id && ing.type === 'batch')) {
      return
    }
    
    // Add batch definition to variant
    variant.ingredients.push({
      id: batch.id,
      batch_definition_id: batch.id,
      type: 'batch', // Mark as batch
      name: batch.name,
      quantity: 1,
      unit: batch.unit,
      stockUnit: batch.unit,
      compatibleUnits: compatibleUnits, // Allow unit conversion (e.g., L→ml, kg→g)
      costPerUnit: batchCostPerUnit, // Calculated cost per unit
      wastage: 0, // Batches don't have wastage
      conversionRate: 1,
      estimatedCost: batchCostPerUnit, // Initial cost for quantity 1
      deductInventory: true
    })
    
    // Calculate initial cost
    const ingIndex = variant.ingredients.length - 1
    updateVariantIngredientCost(currentVariantIndex.value, ingIndex)
    
    // Reset and close
    currentVariantIndex.value = null
    showIngredientSelector.value = false
    ingredientSearchQuery.value = ''
    ingredientType.value = 'raw'
  } else {
    // Adding to single-size item
    if (isBatchSelected(batch.id)) {
      return
    }
    
    form.value.ingredients.push({
      id: batch.id,
      batch_definition_id: batch.id,
      type: 'batch', // Mark as batch
      name: batch.name,
      quantity: 1,
      unit: batch.unit,
      stockUnit: batch.unit,
      compatibleUnits: compatibleUnits, // Allow unit conversion (e.g., L→ml, kg→g)
      costPerUnit: batchCostPerUnit, // Calculated cost per unit
      wastage: 0,
      conversionRate: 1,
      estimatedCost: batchCostPerUnit, // Initial cost for quantity 1
      deductInventory: true
    })
    
    updateIngredientCost(form.value.ingredients.length - 1)
    showIngredientSelector.value = false
    ingredientSearchQuery.value = ''
    ingredientType.value = 'raw'
  }
}

// Update cost for variant ingredient
const updateVariantIngredientCost = (variantIndex, ingredientIndex) => {
  const ing = form.value.variants[variantIndex].ingredients[ingredientIndex]
  
  if (!ing.costPerUnit || ing.costPerUnit <= 0) {
    ing.estimatedCost = 0
    return
  }
  
  const breakdown = calculateCostBreakdown(
    ing.quantity,
    ing.unit,
    ing.costPerUnit,
    ing.stockUnit,
    ing.wastage
  )
  
  ing.estimatedCost = breakdown.totalCost
}

// Update conversion rate when recipe unit changes for variant ingredient
const updateVariantRecipeUnit = (variantIndex, ingredientIndex) => {
  const ing = form.value.variants[variantIndex].ingredients[ingredientIndex]
  
  // Validate conversion
  if (!isValidConversion(ing.stockUnit, ing.unit)) {
    alert(`Không thể quy đổi từ ${ing.stockUnit} sang ${ing.unit}!`)
    ing.unit = ing.stockUnit // Reset to stock unit
    return
  }
  
  // Recalculate conversion rate
  ing.conversionRate = getConversionRate(ing.stockUnit, ing.unit)
  
  // Recalculate cost
  updateVariantIngredientCost(variantIndex, ingredientIndex)
}

const isIngredientSelected = (ingredientId) => {
  // If selecting for a variant, check that variant's ingredients
  if (currentVariantIndex.value !== null) {
    const variant = form.value.variants[currentVariantIndex.value]
    return variant.ingredients.some(ing => ing.id === ingredientId && ing.type === 'raw')
  }
  // Otherwise check single-size ingredients
  return form.value.ingredients.some(ing => ing.id === ingredientId && ing.type === 'raw')
}

const isBatchSelected = (batchId) => {
  // If selecting for a variant, check that variant's ingredients
  if (currentVariantIndex.value !== null) {
    const variant = form.value.variants[currentVariantIndex.value]
    return variant.ingredients.some(ing => ing.id === batchId && ing.type === 'batch')
  }
  // Otherwise check single-size ingredients
  return form.value.ingredients.some(ing => ing.id === batchId && ing.type === 'batch')
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('vi-VN', { 
    day: '2-digit', 
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Update conversion rate when recipe unit changes
const updateRecipeUnit = (index) => {
  const ing = form.value.ingredients[index]
  
  // Validate conversion
  if (!isValidConversion(ing.stockUnit, ing.unit)) {
    alert(`Không thể quy đổi từ ${ing.stockUnit} sang ${ing.unit}!`)
    ing.unit = ing.stockUnit // Reset to stock unit
    return
  }
  
  // Recalculate conversion rate
  ing.conversionRate = getConversionRate(ing.stockUnit, ing.unit)
  
  // Recalculate cost
  updateIngredientCost(index)
}

// Update cost when quantity changes
const updateIngredientCost = (index) => {
  const ing = form.value.ingredients[index]
  
  if (!ing.costPerUnit || ing.costPerUnit <= 0) {
    ing.estimatedCost = 0
    return
  }
  
  // Calculate cost breakdown
  const breakdown = calculateCostBreakdown(
    ing.quantity,
    ing.unit,
    ing.costPerUnit,
    ing.stockUnit,
    ing.wastage
  )
  
  ing.estimatedCost = breakdown.totalCost
}

// Calculate total ingredient cost
const totalIngredientCost = computed(() => {
  return form.value.ingredients.reduce((sum, ing) => sum + (ing.estimatedCost || 0), 0)
})

// Calculate total ingredient cost for a specific variant
const getVariantTotalCost = (variantIndex) => {
  const variant = form.value.variants[variantIndex]
  if (!variant || !variant.ingredients) {
    return 0
  }
  return variant.ingredients.reduce((sum, ing) => sum + (ing.estimatedCost || 0), 0)
}

// Check if form can be saved
const canSave = computed(() => {
  if (!form.value.name || !form.value.category) {
    return false
  }
  
  if (form.value.has_variants) {
    // Multi-size: must have at least 1 variant with valid data
    if (form.value.variants.length === 0) {
      return false
    }
    // Check if all variants have required fields
    return form.value.variants.every(v => 
      v.id && v.id.trim() && 
      v.name && v.name.trim() && 
      v.price > 0
    )
  } else {
    // Single-size: must have price > 0
    return form.value.price > 0
  }
})

const addIngredient = () => {
  showIngredientSelector.value = true
}

const removeIngredient = (index) => {
  form.value.ingredients.splice(index, 1)
}

const saveItem = async () => {
  // Validate variants if multi-size
  if (form.value.has_variants) {
    const errors = validateVariants()
    if (errors.length > 0) {
      variantValidationErrors.value = errors
      alert('Vui lòng kiểm tra lại thông tin các size!')
      return
    }
  }
  
  // Clear validation errors
  variantValidationErrors.value = []
  
  // Prepare data for API
  const itemData = {
    name: form.value.name,
    category: form.value.category,
    description: form.value.description,
    has_variants: form.value.has_variants
  }
  
  if (form.value.has_variants) {
    // Multi-size: send variants
    itemData.variants = form.value.variants.map(v => ({
      id: v.id.trim(),
      name: v.name.trim(),
      price: v.price,
      ingredients: v.ingredients.map(ing => ({
        name: ing.name,
        quantity: ing.quantity,
        unit: ing.unit,
        deduct_inventory: ing.deductInventory ?? true,
        ...(ing.type === 'batch' && {
          ingredient_type: 'batch',
          batch_id: ing.batch_definition_id || ing.id
        }),
        ...(ing.type === 'raw' && {
          ingredient_type: 'raw',
          ingredient_id: ing.id
        })
      })),
      available: v.available !== false,
      is_default: v.is_default || false
    }))
    
    // Debug log
    console.log('Saving multi-size item:', itemData)
  } else {
    // Single-size: send price and ingredients
    itemData.price = form.value.price
    itemData.ingredients = form.value.ingredients.map(ing => ({
      name: ing.name,
      quantity: ing.quantity,
      unit: ing.unit,
      deduct_inventory: ing.deductInventory ?? true,
      ...(ing.type === 'batch' && {
        ingredient_type: 'batch',
        batch_id: ing.batch_definition_id || ing.id
      }),
      ...(ing.type === 'raw' && {
        ingredient_type: 'raw',
        ingredient_id: ing.id
      })
    }))
    
    // Debug log
    console.log('Saving single-size item:', itemData)
  }
  
  let success = false
  
  if (editingItem.value) {
    success = await menuStore.updateMenuItem(editingItem.value.id, itemData)
  } else {
    success = await menuStore.createMenuItem(itemData)
  }
  
  if (success) {
    const wasEditing = !!editingItem.value
    cancelEdit()
    alert(wasEditing ? 'Cập nhật món thành công' : 'Thêm món thành công')
  } else {
    alert('Lỗi: ' + (menuStore.error || 'Có lỗi xảy ra'))
  }
}

const toggleAvailable = async (item) => {
  const success = await menuStore.updateMenuItem(item.id, { available: !item.available })
  if (!success && menuStore.error) {
    alert('Lỗi: ' + menuStore.error)
  }
}

const deleteItem = async (id) => {
  if (confirm('Bạn có chắc muốn xóa món này? Hành động này không thể hoàn tác.')) {
    const success = await menuStore.deleteMenuItem(id)
    if (success) {
      alert('Xóa món thành công')
    } else {
      alert('Lỗi: ' + (menuStore.error || 'Có lỗi xảy ra'))
    }
  }
}

const addCategory = async () => {
  if (!categoryForm.value.name) return
  
  // Check if category already exists
  if (menuCategories.value.some(c => c.name.toLowerCase() === categoryForm.value.name.toLowerCase())) {
    alert('Danh mục đã tồn tại!')
    return
  }
  
  try {
    await menuCategoryService.createCategory({ name: categoryForm.value.name })
    const categoryName = categoryForm.value.name
    categoryForm.value = { name: '' }
    // Refresh categories first
    await fetchCategories()
    // Then show success message
    alert(`Danh mục "${categoryName}" đã được tạo thành công`)
  } catch (err) {
    console.error('Failed to create category:', err)
    alert('Lỗi: Không thể tạo danh mục. ' + (err.response?.data?.error || err.message))
  }
}

const moveCategoryUp = async (idx) => {
  if (idx === 0) return
  const cats = [...menuCategories.value]
  ;[cats[idx - 1], cats[idx]] = [cats[idx], cats[idx - 1]]
  menuCategories.value = cats
  await saveCategoryOrder()
}

const moveCategoryDown = async (idx) => {
  if (idx === menuCategories.value.length - 1) return
  const cats = [...menuCategories.value]
  ;[cats[idx], cats[idx + 1]] = [cats[idx + 1], cats[idx]]
  menuCategories.value = cats
  await saveCategoryOrder()
}

const saveCategoryOrder = async () => {
  try {
    const orders = menuCategories.value.map((cat, i) => ({ id: cat.id, sort_order: i }))
    await menuCategoryService.reorderCategories(orders)
  } catch (err) {
    console.error('Failed to save category order:', err)
  }
}

const deleteCategory = async (id, name) => {
  const hasMenuItems = items.value.some(item => item.category === name)
  
  if (hasMenuItems) {
    alert('Không thể xóa danh mục đã có món! Vui lòng xóa tất cả món trong danh mục trước.')
    return
  }
  
  if (!confirm(`Bạn có chắc muốn xóa danh mục "${name}"?`)) {
    return
  }
  
  try {
    await menuCategoryService.deleteCategory(id)
    // Refresh categories first
    await fetchCategories()
    // Then show success message
    alert(`Danh mục "${name}" đã được xóa`)
  } catch (err) {
    console.error('Failed to delete category:', err)
    alert('Lỗi: Không thể xóa danh mục. ' + (err.response?.data?.error || err.message))
  }
}

const getMenuCountByCategory = (categoryName) => {
  if (!items.value || !Array.isArray(items.value)) {
    return 0
  }
  return items.value.filter(item => item.category === categoryName).length
}

const getCategoryIcon = (category) => {
  const iconMap = {
    'Cà phê': '☕',
    'Trà': '🍵',
    'Nước ép': '🧃',
    'Bánh ngọt': '🍰',
    'Món nhẹ': '🍴',
    'Sinh tố': '🥤',
    'Đồ uống khác': '🥛'
  }
  return iconMap[category] || '🍽️'
}

const getCategoryColor = (category) => {
  const colorMap = {
    'Cà phê': 'bg-amber-100 text-amber-600',
    'Trà': 'bg-green-100 text-green-600',
    'Nước ép': 'bg-orange-100 text-orange-600',
    'Bánh ngọt': 'bg-pink-100 text-pink-600',
    'Món nhẹ': 'bg-blue-100 text-blue-600',
    'Sinh tố': 'bg-purple-100 text-purple-600',
    'Đồ uống khác': 'bg-gray-100 text-gray-600'
  }
  return colorMap[category] || 'bg-gray-100 text-gray-600'
}
</script>

<style scoped>
.active\:scale-95:active {
  transform: scale(0.95);
}

.active\:scale-98:active {
  transform: scale(0.98);
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
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

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
