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

              <!-- Price -->
              <div class="bg-white rounded-lg p-3 mb-3">
                <div class="text-2xl font-bold text-green-600">{{ formatPrice(item.price) }}</div>
                <div class="text-xs text-gray-500">Giá bán</div>
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
              
              <div v-for="cat in menuCategories" :key="cat.id" 
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
                <button @click="deleteCategory(cat.id, cat.name)" class="text-red-500 hover:text-red-700 p-2 flex-shrink-0 ml-2">
                  🗑️
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- Create/Edit Form Modal - Slide from Right -->
    <transition name="slide-right">
      <div v-if="showMenuForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-gray-50 w-full h-screen flex flex-col">
          <!-- Mobile Header - Fixed -->
          <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
            <div class="px-4 py-3">
              <div class="flex items-center justify-between">
                <button @click="cancelEdit" class="text-2xl text-gray-600">←</button>
                <h1 class="text-xl font-bold text-gray-800">{{ editingItem ? '✏️ Cập nhật món' : '➕ Thêm món mới' }}</h1>
                <div class="w-8"></div>
              </div>
            </div>
          </div>

          <!-- Scrollable Content -->
          <div class="flex-1 overflow-y-auto px-4 py-6 space-y-5">
            <!-- Tên món -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">Tên món *</label>
              <input v-model="form.name" type="text" required placeholder="VD: Cà phê sữa đá"
                class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
            </div>

            <!-- Danh mục & Giá - Responsive Grid -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">Danh mục *</label>
                <select v-model="form.category" required
                  class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent">
                  <option value="">Chọn danh mục</option>
                  <option v-for="cat in suggestedCategories" :key="cat.id" :value="cat.name">{{ cat.name }}</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-3">Giá (VNĐ) *</label>
                <input v-model.number="form.price" type="number" min="0" step="1000" required placeholder="0"
                  class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
              </div>
            </div>

            <!-- Mô tả -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">Mô tả</label>
              <textarea v-model="form.description" rows="3" placeholder="Mô tả món ăn..."
                class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"></textarea>
            </div>

            <!-- Nguyên liệu -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">🥘 Nguyên liệu</label>
              <div class="border border-gray-300 rounded-lg p-4 bg-white">
                <div v-if="form.ingredients.length === 0" class="text-center text-gray-500 italic py-8">
                  <div class="text-4xl mb-2">📋</div>
                  <p>Chưa có nguyên liệu nào</p>
                </div>
                <div v-else class="space-y-3 mb-3">
                  <div v-for="(ingredient, index) in form.ingredients" :key="index" 
                    class="bg-gray-50 rounded-lg p-3">
                    
                    <!-- Header -->
                    <div class="flex justify-between items-start mb-2">
                      <div class="flex-1">
                        <div class="font-medium text-gray-800">{{ ingredient.name }}</div>
                        <div class="text-xs text-gray-500">
                          Kho: {{ ingredient.stockUnit }} @ {{ formatPrice(ingredient.costPerUnit) }}/{{ ingredient.stockUnit }}
                        </div>
                      </div>
                      <button type="button" @click="removeIngredient(index)" 
                        class="bg-red-500 text-white px-3 py-1 rounded-lg hover:bg-red-600 flex-shrink-0 text-sm">
                        ×
                      </button>
                    </div>

                    <!-- Recipe Unit Selector -->
                    <div class="mb-2">
                      <label class="text-xs text-gray-600">Đơn vị công thức:</label>
                      <select v-model="ingredient.unit" @change="updateRecipeUnit(index)"
                        class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg">
                        <option v-for="unit in ingredient.compatibleUnits" :key="unit" :value="unit">
                          {{ unit }}
                        </option>
                      </select>
                    </div>
                    
                    <!-- Quantity Input -->
                    <div class="mb-2">
                      <label class="text-xs text-gray-600">Số lượng:</label>
                      <div class="flex gap-2 items-center">
                        <input v-model.number="ingredient.quantity" 
                          @input="updateIngredientCost(index)"
                          type="number" min="0" step="0.1" placeholder="0" required
                          class="flex-1 px-3 py-2 text-base border border-gray-300 rounded-lg" />
                        <span class="text-sm text-gray-600">{{ ingredient.unit }}</span>
                      </div>
                    </div>
                    
                    <!-- Conversion Info (if not 1.0) -->
                    <div v-if="ingredient.conversionRate !== 1" 
                      class="mb-2 p-2 bg-blue-50 rounded text-xs text-blue-700">
                      <span class="font-medium">ℹ️ Quy đổi:</span>
                      {{ getConversionExplanation(ingredient.stockUnit, ingredient.unit) }}
                    </div>
                    
                    <!-- Cost Preview -->
                    <div v-if="ingredient.costPerUnit > 0" class="p-2 bg-green-50 rounded">
                      <div class="flex justify-between items-center">
                        <span class="text-xs text-green-700">Chi phí ước tính:</span>
                        <span class="text-sm font-bold text-green-700">
                          {{ formatPrice(ingredient.estimatedCost) }}
                        </span>
                      </div>
                      <div v-if="ingredient.wastage > 0" class="text-xs text-green-600 mt-1">
                        (Bao gồm {{ ingredient.wastage }}% hao hụt)
                      </div>
                    </div>
                  </div>
                  
                  <!-- Total Cost Summary -->
                  <div v-if="totalIngredientCost > 0" class="bg-blue-50 border-2 border-blue-300 rounded-lg p-3">
                    <div class="flex justify-between items-center">
                      <span class="text-sm font-semibold text-blue-800">💰 Tổng chi phí nguyên liệu:</span>
                      <span class="text-lg font-bold text-blue-900">{{ formatPrice(totalIngredientCost) }}</span>
                    </div>
                  </div>
                </div>
                <button type="button" @click="showIngredientSelector = true" 
                  class="w-full bg-blue-500 text-white py-3 rounded-lg font-medium text-base active:bg-blue-600">
                  + Chọn nguyên liệu
                </button>
              </div>
            </div>

            <!-- Spacer for bottom buttons -->
            <div class="h-24"></div>
          </div>

          <!-- Fixed Footer -->
          <div class="flex-shrink-0 bg-white px-4 py-4 border-t flex gap-3 pb-safe">
            <button @click="cancelEdit" 
              class="flex-1 bg-gray-200 text-gray-700 py-4 rounded-xl font-medium text-base active:bg-gray-300 transition-colors">
              Hủy
            </button>
            <button @click="saveItem" :disabled="!form.name || !form.category || form.price <= 0"
              class="flex-1 bg-green-500 text-white py-4 rounded-xl font-medium text-base active:bg-green-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
              {{ editingItem ? 'Cập nhật' : 'Thêm món' }}
            </button>
          </div>
        </div>
      </div>
    </transition>

    <!-- Ingredient Selector Modal -->
    <transition name="slide-up">
      <div v-if="showIngredientSelector" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full h-[85vh] flex flex-col">
          <!-- Fixed Header -->
          <div class="flex-shrink-0 bg-white px-4 py-4 border-b flex justify-between items-center rounded-t-3xl">
            <h3 class="text-lg font-bold">🥬 Chọn nguyên liệu</h3>
            <button @click="showIngredientSelector = false" class="text-2xl text-gray-400">×</button>
          </div>
          
          <!-- Search -->
          <div class="flex-shrink-0 px-4 py-3 border-b">
            <input
              v-model="ingredientSearchQuery"
              type="text"
              placeholder="Tìm kiếm nguyên liệu..."
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
            />
          </div>
          
          <!-- Scrollable Content -->
          <div class="flex-1 overflow-y-auto px-4 py-4">
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
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useMenuStore } from '../stores/menu'
import { useIngredientStore } from '../stores/ingredient'
import { useUnitConversion } from '../composables/useUnitConversion'
import { menuCategoryService } from '../services/menuCategory'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'

const menuStore = useMenuStore()
const ingredientStore = useIngredientStore()

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
const editingItem = ref(null)
const form = ref({
  name: '',
  category: '',
  price: 0,
  description: '',
  ingredients: []
})

const categoryForm = ref({
  name: ''
})

const items = computed(() => menuStore.items)
const loading = computed(() => menuStore.loading)
const error = computed(() => menuStore.error)

// Ingredients
const availableIngredients = computed(() => ingredientStore.items)
const ingredientsLoading = computed(() => ingredientStore.loading)

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
    ingredientStore.fetchIngredients()
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
    ingredients: []
  }
  editingItem.value = null
  showMenuForm.value = true
}

const viewDetails = (item) => {
  editItem(item)
}

const editItem = (item) => {
  editingItem.value = item
  form.value = {
    name: item.name,
    category: item.category,
    price: item.price,
    description: item.description,
    ingredients: item.ingredients ? [...item.ingredients] : []
  }
  showMenuForm.value = true
}

const cancelEdit = () => {
  showMenuForm.value = false
  editingItem.value = null
  form.value = { name: '', category: '', price: 0, description: '', ingredients: [] }
}

const selectIngredient = (ingredient) => {
  // Check if already selected
  if (isIngredientSelected(ingredient.id)) {
    return
  }
  
  // Get compatible units for this ingredient
  const compatibleUnits = getCompatibleUnits(ingredient.unit)
  
  // Default recipe unit = stock unit (no conversion initially)
  const recipeUnit = ingredient.unit
  const conversionRate = getConversionRate(ingredient.unit, recipeUnit)
  
  // Add ingredient with conversion info
  form.value.ingredients.push({
    id: ingredient.id,
    name: ingredient.name,
    quantity: 1,
    unit: recipeUnit,                    // Recipe unit (can be changed by user)
    stockUnit: ingredient.unit,          // Stock unit (fixed, from ingredient)
    compatibleUnits: compatibleUnits,    // Available units for dropdown
    costPerUnit: ingredient.cost_per_unit || 0,
    wastage: ingredient.wastage_percentage || 0,
    conversionRate: conversionRate,      // Initially 1.0 (same unit)
    estimatedCost: 0                     // Will be calculated
  })
  
  // Calculate initial cost
  updateIngredientCost(form.value.ingredients.length - 1)
  
  // Close modal
  showIngredientSelector.value = false
  ingredientSearchQuery.value = ''
}

const isIngredientSelected = (ingredientId) => {
  return form.value.ingredients.some(ing => ing.id === ingredientId)
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

const addIngredient = () => {
  showIngredientSelector.value = true
}

const removeIngredient = (index) => {
  form.value.ingredients.splice(index, 1)
}

const saveItem = async () => {
  let success = false
  
  if (editingItem.value) {
    success = await menuStore.updateMenuItem(editingItem.value.id, form.value)
  } else {
    success = await menuStore.createMenuItem(form.value)
  }
  
  if (success) {
    cancelEdit()
    alert(editingItem.value ? 'Cập nhật món thành công' : 'Thêm món thành công')
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
