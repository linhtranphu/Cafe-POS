<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm">
      <div class="px-4 py-3">
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
    <div class="px-4 py-4 pb-24">
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
            <div class="mb-3 space-y-1 text-sm">
              <div class="flex justify-between">
                <span class="text-gray-600">Tồn kho:</span>
                <span class="font-medium">{{ ingredient.quantity }} {{ ingredient.unit }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">Tối thiểu:</span>
                <span class="font-medium">{{ ingredient.min_stock }} {{ ingredient.unit }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">Đơn giá:</span>
                <span class="font-medium text-green-600">{{ formatCurrency(ingredient.cost_per_unit) }}</span>
              </div>
            </div>

            <!-- Quick Actions -->
            <div class="grid grid-cols-4 gap-2 pt-3 border-t">
              <button @click="openAdjustModal(ingredient)"
                class="bg-blue-500 text-white py-2 rounded-lg text-xs font-medium active:bg-blue-600">
                📦 Điều chỉnh
              </button>
              <button @click="viewHistory(ingredient)"
                class="bg-purple-500 text-white py-2 rounded-lg text-xs font-medium active:bg-purple-600">
                📊 Lịch sử
              </button>
              <button @click="openEditModal(ingredient)"
                class="bg-green-500 text-white py-2 rounded-lg text-xs font-medium active:bg-green-600">
                ✏️ Sửa
              </button>
              <button @click="deleteIngredient(ingredient)"
                class="bg-red-500 text-white py-2 rounded-lg text-xs font-medium active:bg-red-600">
                🗑️ Xóa
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

    <!-- Create/Edit Modal -->
      <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
          <h2 class="text-2xl font-bold mb-4">{{ isEditing ? 'Cập Nhật Nguyên Liệu' : 'Thêm Nguyên Liệu Mới' }}</h2>
          
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Tên Nguyên Liệu *</label>
              <input v-model="formData.name" type="text" class="w-full px-3 py-2 border border-gray-300 rounded-lg" />
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Danh Mục *</label>
                <select v-model="formData.category" class="w-full px-3 py-2 border border-gray-300 rounded-lg">
                  <option value="">Chọn danh mục</option>
                  <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Đơn Vị *</label>
                <select v-model="formData.unit" class="w-full px-3 py-2 border border-gray-300 rounded-lg">
                  <option value="">Chọn đơn vị</option>
                  <option v-for="unit in INGREDIENT_UNIT_OPTIONS" :key="unit.value" :value="unit.value">
                    {{ unit.label }}
                  </option>
                </select>
              </div>
            </div>

            <div class="grid grid-cols-3 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Số Lượng *</label>
                <input v-model.number="formData.quantity" type="number" step="0.01" class="w-full px-3 py-2 border border-gray-300 rounded-lg" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Số Lượng Tối Thiểu *</label>
                <input v-model.number="formData.min_stock" type="number" step="0.01" class="w-full px-3 py-2 border border-gray-300 rounded-lg" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Giá/Đơn Vị *</label>
                <input v-model.number="formData.cost_per_unit" type="number" step="0.01" class="w-full px-3 py-2 border border-gray-300 rounded-lg" />
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Nhà Cung Cấp</label>
              <input v-model="formData.supplier" type="text" class="w-full px-3 py-2 border border-gray-300 rounded-lg" />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Ghi Chú</label>
              <textarea v-model="formData.notes" rows="3" class="w-full px-3 py-2 border border-gray-300 rounded-lg"></textarea>
            </div>

            <!-- Auto-Expense Indicator -->
            <div v-if="!isEditing && formData.quantity > 0 && formData.cost_per_unit > 0" 
              class="bg-green-50 border border-green-200 rounded-lg p-3">
              <div class="flex items-start gap-2">
                <span class="text-green-600 text-lg">✅</span>
                <div class="flex-1">
                  <p class="text-sm font-medium text-green-800">Tự động ghi nhận chi phí</p>
                  <p class="text-xs text-green-600 mt-1">
                    Hệ thống sẽ tự động tạo chi phí: {{ formatCurrency(formData.quantity * formData.cost_per_unit) }}
                  </p>
                  <p class="text-xs text-green-600">Danh mục: Nguyên liệu</p>
                </div>
              </div>
            </div>
          </div>

          <div class="flex justify-end gap-3 mt-6">
            <button @click="closeModal" class="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50">
              Hủy
            </button>
            <button @click="saveIngredient" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
              {{ isEditing ? 'Cập Nhật' : 'Thêm Mới' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Adjust Stock Modal -->
      <div v-if="showAdjustModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-6 w-full max-w-md">
          <h2 class="text-2xl font-bold mb-4">Điều Chỉnh Tồn Kho</h2>
          
          <div class="mb-4">
            <p class="text-sm text-gray-600">Nguyên liệu: <span class="font-semibold">{{ currentIngredient?.name }}</span></p>
            <p class="text-sm text-gray-600">Tồn kho hiện tại: <span class="font-semibold">{{ currentIngredient?.quantity }} {{ currentIngredient?.unit }}</span></p>
          </div>

          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Loại Điều Chỉnh *</label>
              <select v-model="adjustData.type" class="w-full px-3 py-2 border border-gray-300 rounded-lg">
                <option v-for="type in ADJUSTMENT_TYPE_OPTIONS" :key="type.value" :value="type.value">
                  {{ type.label }}
                </option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Số Lượng *</label>
              <input v-model.number="adjustData.quantity" type="number" step="0.01" class="w-full px-3 py-2 border border-gray-300 rounded-lg" />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Lý Do *</label>
              <textarea v-model="adjustData.reason" rows="3" class="w-full px-3 py-2 border border-gray-300 rounded-lg"></textarea>
            </div>

            <div class="bg-blue-50 p-3 rounded-lg">
              <p class="text-sm text-blue-800">
                Tồn kho sau điều chỉnh: <span class="font-bold">{{ calculateNewQuantity() }} {{ currentIngredient?.unit }}</span>
              </p>
            </div>

            <!-- Auto-Expense Indicator for Stock IN -->
            <div v-if="adjustData.type === 'add' && adjustData.quantity > 0 && currentIngredient?.cost_per_unit > 0" 
              class="bg-green-50 border border-green-200 rounded-lg p-3">
              <div class="flex items-start gap-2">
                <span class="text-green-600 text-lg">✅</span>
                <div class="flex-1">
                  <p class="text-sm font-medium text-green-800">Tự động ghi nhận chi phí</p>
                  <p class="text-xs text-green-600 mt-1">
                    Hệ thống sẽ tự động tạo chi phí: {{ formatCurrency(adjustData.quantity * currentIngredient.cost_per_unit) }}
                  </p>
                  <p class="text-xs text-green-600">Danh mục: Nguyên liệu</p>
                </div>
              </div>
            </div>
          </div>

          <div class="flex justify-end gap-3 mt-6">
            <button @click="closeAdjustModal" class="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50">
              Hủy
            </button>
            <button @click="adjustStock" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
              Xác Nhận
            </button>
          </div>
        </div>
      </div>

      <!-- Stock History Modal -->
      <div v-if="showHistoryModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-6 w-full max-w-4xl max-h-[90vh] overflow-y-auto">
          <h2 class="text-2xl font-bold mb-4">Lịch Sử Tồn Kho - {{ currentIngredient?.name }}</h2>
          
          <div class="space-y-3">
            <div v-for="record in stockHistory" :key="record.id" class="border border-gray-200 rounded-lg p-4">
              <div class="flex justify-between items-start">
                <div>
                  <div class="flex items-center gap-2">
                    <span :class="getAdjustmentTypeClass(record.type)" class="px-2 py-1 text-xs font-semibold rounded-full">
                      {{ getAdjustmentTypeText(record.type) }}
                    </span>
                    <span class="text-sm font-medium text-gray-900">{{ record.quantity }} {{ record.unit }}</span>
                  </div>
                  <p class="text-sm text-gray-600 mt-1">{{ record.reason }}</p>
                  <p class="text-xs text-gray-500 mt-2">
                    {{ formatDateTime(record.created_at) }} - {{ record.created_by }}
                  </p>
                </div>
                <div class="text-right">
                  <p class="text-xs text-gray-500">Trước</p>
                  <p class="text-sm font-medium">{{ record.quantity_before }}</p>
                  <p class="text-xs text-gray-500 mt-1">Sau</p>
                  <p class="text-sm font-medium">{{ record.quantity_after }}</p>
                </div>
              </div>
            </div>
          </div>

          <div class="flex justify-end mt-6">
            <button @click="closeHistoryModal" class="px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700">
              Đóng
            </button>
          </div>
        </div>
      </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useIngredientStore } from '../stores/ingredient'
import BottomNav from '../components/BottomNav.vue'
import {
  INGREDIENT_UNIT_OPTIONS,
  ADJUSTMENT_TYPE_OPTIONS,
  getStockStatusClass,
  getStockStatusText,
  getAdjustmentTypeClass,
  getAdjustmentTypeText
} from '../constants/ingredient'

export default {
  name: 'IngredientManagementView',
  components: {
    BottomNav
  },
  setup() {
    const ingredientStore = useIngredientStore()
    
    const searchQuery = ref('')
    const showModal = ref(false)
    const showAdjustModal = ref(false)
    const showHistoryModal = ref(false)
    const showCategoryModal = ref(false)
    const isEditing = ref(false)
    const currentIngredient = ref(null)
    const newCategoryName = ref('')
    
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
      notes: ''
    })

    const adjustData = ref({
      type: 'add',
      quantity: 0,
      reason: ''
    })

    const ingredients = computed(() => ingredientStore.items || [])
    const stockHistory = ref([])

    const filteredIngredients = computed(() => {
      if (!searchQuery.value) return ingredients.value
      const query = searchQuery.value.toLowerCase()
      return ingredients.value.filter(i => 
        i.name.toLowerCase().includes(query) ||
        i.category.toLowerCase().includes(query) ||
        i.supplier?.toLowerCase().includes(query)
      )
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

    const calculateNewQuantity = () => {
      if (!currentIngredient.value) return 0
      const current = currentIngredient.value.quantity
      const adjust = adjustData.value.quantity || 0
      
      if (adjustData.value.type === 'add') return current + adjust
      if (adjustData.value.type === 'remove') return current - adjust
      return adjust
    }

    const openCreateModal = () => {
      isEditing.value = false
      currentIngredient.value = null
      formData.value = {
        name: '',
        category: '',
        unit: '',
        quantity: 0,
        min_stock: 0,
        cost_per_unit: 0,
        supplier: '',
        notes: ''
      }
      showModal.value = true
    }

    const openEditModal = (ingredient) => {
      isEditing.value = true
      currentIngredient.value = ingredient
      formData.value = { ...ingredient }
      showModal.value = true
    }

    const closeModal = () => {
      showModal.value = false
      isEditing.value = false
      currentIngredient.value = null
    }

    const openAdjustModal = (ingredient) => {
      currentIngredient.value = ingredient
      adjustData.value = {
        type: 'add',
        quantity: 0,
        reason: ''
      }
      showAdjustModal.value = true
    }

    const closeAdjustModal = () => {
      showAdjustModal.value = false
      currentIngredient.value = null
    }

    const closeHistoryModal = () => {
      showHistoryModal.value = false
      currentIngredient.value = null
    }

    const saveIngredient = async () => {
      try {
        if (isEditing.value) {
          await ingredientStore.updateIngredient(currentIngredient.value.id, formData.value)
        } else {
          await ingredientStore.createIngredient(formData.value)
        }
        closeModal()
      } catch (error) {
        console.error('Error saving ingredient:', error)
        alert('Có lỗi xảy ra khi lưu nguyên liệu')
      }
    }

    const adjustStock = async () => {
      try {
        await ingredientStore.adjustStock(currentIngredient.value.id, adjustData.value)
        closeAdjustModal()
      } catch (error) {
        console.error('Error adjusting stock:', error)
        alert('Có lỗi xảy ra khi điều chỉnh tồn kho')
      }
    }

    const deleteIngredient = async (ingredient) => {
      if (confirm(`Bạn có chắc muốn xóa nguyên liệu "${ingredient.name}"?`)) {
        try {
          await ingredientStore.deleteIngredient(ingredient.id)
        } catch (error) {
          console.error('Error deleting ingredient:', error)
          alert('Có lỗi xảy ra khi xóa nguyên liệu')
        }
      }
    }

    const viewHistory = async (ingredient) => {
      currentIngredient.value = ingredient
      stockHistory.value = await ingredientStore.fetchStockHistory(ingredient.id)
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

    onMounted(async () => {
      await ingredientStore.fetchCategories()
      await ingredientStore.fetchIngredients()
    })

    return {
      searchQuery,
      showModal,
      showAdjustModal,
      showHistoryModal,
      showCategoryModal,
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
      getStockStatusClass,
      getStockStatusText,
      getAdjustmentTypeClass,
      getAdjustmentTypeText,
      formatCurrency,
      formatDateTime,
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
      getCategoryCount
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
</style>
