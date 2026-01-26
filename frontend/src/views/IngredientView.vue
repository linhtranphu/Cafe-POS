<template>
  <div class="min-h-screen bg-gray-100">
    <Navigation />
    <div class="p-4">
      <div class="flex flex-col lg:flex-row justify-between items-center mb-6">
        <h2 class="text-xl lg:text-2xl font-semibold text-gray-800 mb-4 lg:mb-0">
          🥬 Quản lý Nguyên liệu
        </h2>
        <div class="flex flex-wrap gap-2">
          <button @click="showCategoryForm = true" class="bg-purple-500 hover:bg-purple-600 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors">
            📁 Quản lý danh mục
          </button>
          <button v-if="isManager" @click="showCreateForm = true" class="btn-primary text-sm px-4 py-2">
            + Thêm nguyên liệu
          </button>
        </div>
      </div>

      <div v-if="loading" class="text-center py-10 text-gray-600 text-lg">Đang tải...</div>
      <div v-if="error" class="text-center py-10 text-red-600 bg-red-50 border border-red-200 rounded-lg">{{ error }}</div>

      <!-- Filters Section -->
      <div class="bg-white rounded-xl p-4 mb-4 shadow-sm">
        <div class="grid grid-cols-1 gap-3">
          <input v-model="searchQuery" type="text" placeholder="Tìm theo tên nguyên liệu..." class="w-full p-3 border border-gray-300 rounded-lg text-base focus:ring-2 focus:ring-blue-500" />
          <select v-model="filterCategory" class="p-3 border border-gray-300 rounded-lg text-base focus:ring-2 focus:ring-blue-500">
            <option value="">Tất cả danh mục</option>
            <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
          </select>
          <select v-model="filterStatus" class="p-3 border border-gray-300 rounded-lg text-base focus:ring-2 focus:ring-blue-500">
            <option value="">Tất cả trạng thái</option>
            <option value="in-stock">Còn hàng</option>
            <option value="low-stock">Sắp hết</option>
            <option value="out-of-stock">Hết hàng</option>
          </select>
        </div>
      </div>

      <!-- Low Stock Alert -->
      <div v-if="lowStockItems && lowStockItems.length > 0" class="bg-yellow-50 border border-yellow-200 rounded-xl p-4 mb-4">
        <h3 class="text-yellow-800 font-semibold mb-2">⚠️ Nguyên liệu sắp hết ({{ lowStockItems.length }} món)</h3>
        <div class="flex flex-wrap gap-2">
          <span v-for="item in lowStockItems" :key="item.id" class="bg-yellow-200 text-yellow-800 px-2 py-1 rounded text-xs">
            {{ item.name }}: {{ item.quantity }} {{ item.unit }}
          </span>
        </div>
      </div>

      <div class="grid grid-cols-1 gap-4">
        <div v-for="category in groupedItems" :key="category.name" class="bg-white rounded-xl p-4 shadow-sm">
          <h3 class="text-lg font-bold text-gray-800 mb-4 pb-2 border-b-2 border-blue-500">{{ category.name }}</h3>
          <div class="space-y-3">
            <div v-for="item in category.items" :key="item.id" class="rounded-xl p-4" :class="isLowStock(item) ? 'bg-yellow-50 border-2 border-yellow-400' : 'bg-gray-50'">
              <!-- Ingredient Header -->
              <div class="flex items-center justify-between mb-3">
                <div class="flex items-center space-x-3">
                  <div class="w-12 h-12 rounded-xl flex items-center justify-center text-2xl" :class="getCategoryColor(item.category)">
                    {{ getCategoryIcon(item.category) }}
                  </div>
                  <div>
                    <h4 class="font-bold text-gray-800 flex items-center gap-2">
                      {{ item.name }}
                      <span v-if="isLowStock(item)" class="text-yellow-600 text-lg">⚠️</span>
                    </h4>
                    <p class="text-sm text-gray-500">{{ item.supplier || 'Chưa có nhà cung cấp' }}</p>
                  </div>
                </div>
                <div class="text-right">
                  <span class="px-3 py-1 rounded-full text-xs font-medium" :class="getStockBadge(item)">
                    {{ getStockStatus(item) }}
                  </span>
                </div>
              </div>

              <!-- Stock Info Grid -->
              <div class="grid grid-cols-2 gap-3 mb-4">
                <div class="bg-white rounded-lg p-3 text-center">
                  <div class="text-2xl font-bold" :class="item.quantity <= item.min_stock ? 'text-red-600' : 'text-green-600'">
                    {{ item.quantity }}
                  </div>
                  <div class="text-xs text-gray-500">{{ item.unit }} hiện tại</div>
                </div>
                <div class="bg-white rounded-lg p-3 text-center">
                  <div class="text-lg font-semibold text-gray-600">{{ item.min_stock }}</div>
                  <div class="text-xs text-gray-500">{{ item.unit }} tối thiểu</div>
                </div>
                <div class="bg-white rounded-lg p-3 text-center col-span-2" v-if="item.cost_per_unit">
                  <div class="text-lg font-semibold text-blue-600">{{ formatPrice(item.cost_per_unit) }}</div>
                  <div class="text-xs text-gray-500">Giá/{{ item.unit }}</div>
                </div>
              </div>

              <!-- Quick Actions -->
              <div class="grid grid-cols-2 gap-2 mb-3">
                <button @click="showHistory(item)" class="bg-purple-500 hover:bg-purple-600 text-white px-4 py-3 rounded-xl text-sm font-medium transition-colors flex items-center justify-center space-x-2">
                  <span>📈</span>
                  <span>Lịch sử</span>
                </button>
                <button v-if="isManager" @click="showAdjustStock(item)" class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-3 rounded-xl text-sm font-medium transition-colors flex items-center justify-center space-x-2">
                  <span>📦</span>
                  <span>Điều chỉnh</span>
                </button>
              </div>

              <!-- More Actions Toggle -->
              <button @click="item.showMore = !item.showMore" class="w-full py-2 text-blue-600 text-sm font-medium border border-blue-200 rounded-lg hover:bg-blue-50 transition-colors">
                {{ item.showMore ? '▲ Ẩn bớt' : '▼ Thêm tùy chọn' }}
              </button>

              <!-- Extended Actions -->
              <div v-if="item.showMore" class="grid grid-cols-2 gap-2 mt-3">
                <button v-if="isManager" @click="editItem(item)" class="bg-yellow-500 hover:bg-yellow-600 text-white px-3 py-2 rounded-lg text-sm font-medium transition-colors">
                  📝 Sửa
                </button>
                <button v-if="isManager" @click="deleteItem(item.id)" class="bg-red-500 hover:bg-red-600 text-white px-3 py-2 rounded-lg text-sm font-medium transition-colors">
                  🗑️ Xóa
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Category Management Modal -->
      <div v-if="showCategoryForm" class="modal">
        <div class="modal-content">
          <h3>📁 Quản lý Danh mục Nguyên liệu</h3>
          
          <!-- Add New Category -->
          <div class="bg-gray-50 rounded-lg p-4 mb-4">
            <h4 class="font-semibold text-gray-800 mb-3">Thêm danh mục mới</h4>
            <form @submit.prevent="addCategory">
              <div class="form-group">
                <label>Tên danh mục *</label>
                <input v-model="categoryForm.name" type="text" required placeholder="Ví dụ: Cà phê" />
              </div>
              <button type="submit" class="btn-primary w-full">+ Thêm danh mục</button>
            </form>
          </div>

          <!-- Category List -->
          <div class="space-y-2 max-h-96 overflow-y-auto">
            <div v-for="cat in ingredientCategories" :key="cat.id" class="bg-white border border-gray-200 rounded-lg p-3 flex items-center justify-between">
              <div class="flex items-center space-x-3">
                <div class="w-10 h-10 rounded-lg flex items-center justify-center text-xl" :class="getCategoryColor(cat.name)">
                  {{ getCategoryIcon(cat.name) }}
                </div>
                <div>
                  <div class="font-medium text-gray-800">{{ cat.name }}</div>
                  <div class="text-xs text-gray-500">{{ getIngredientCountByCategory(cat.name) }} nguyên liệu</div>
                </div>
              </div>
              <button @click="deleteCategory(cat.id, cat.name)" class="text-red-500 hover:text-red-700 p-2">
                🗑️
              </button>
            </div>
          </div>

          <div class="form-actions mt-4">
            <button type="button" @click="showCategoryForm = false" class="btn-cancel">Đóng</button>
          </div>
        </div>
      </div>

      <!-- Create/Edit Form Modal -->
    <div v-if="showCreateForm || editingItem" class="modal">
      <div class="modal-content">
        <h3>{{ editingItem ? 'Sửa nguyên liệu' : 'Thêm nguyên liệu mới' }}</h3>
        <form @submit.prevent="saveItem">
          <div class="form-group">
            <label>Tên nguyên liệu *</label>
            <input v-model="form.name" type="text" required placeholder="Ví dụ: Cà phê bột" />
          </div>
          <div class="form-group">
            <label>Danh mục *</label>
            <select v-model="form.category" required>
              <option value="">Chọn danh mục</option>
              <option v-for="cat in ingredientCategories" :key="cat.id" :value="cat.name">{{ cat.name }}</option>
            </select>
          </div>
          <div class="form-group">
            <label>Đơn vị *</label>
            <select v-model="form.unit" required>
              <option value="">Chọn đơn vị</option>
              <option value="kg">Kilogram (kg)</option>
              <option value="g">Gram (g)</option>
              <option value="L">Liter (L)</option>
              <option value="ml">Milliliter (ml)</option>
              <option value="piece">Cái</option>
              <option value="box">Hộp</option>
              <option value="pack">Gói</option>
            </select>
          </div>
          <div class="form-group">
            <label>Số lượng hiện tại *</label>
            <input v-model.number="form.quantity" type="number" min="0" step="0.1" required />
          </div>
          <div class="form-group">
            <label>Tồn kho tối thiểu *</label>
            <input v-model.number="form.min_stock" type="number" min="0" step="0.1" required />
          </div>
          <div class="form-group">
            <label>Giá mỗi đơn vị (VNĐ)</label>
            <input v-model.number="form.cost_per_unit" type="number" min="0" step="100" />
          </div>
          <div class="form-group">
            <label>Nhà cung cấp</label>
            <input v-model="form.supplier" type="text" placeholder="Tên nhà cung cấp" />
          </div>
          <div class="form-actions">
            <button type="button" @click="cancelEdit" class="btn-cancel">Hủy</button>
            <button type="submit" class="btn-save">{{ editingItem ? 'Cập nhật' : 'Thêm' }}</button>
          </div>
        </form>
      </div>
      </div>

      <!-- Stock History Modal -->
    <div v-if="historyItem" class="modal">
      <div class="modal-content history-modal">
        <h3>📈 Lịch sử biến động: {{ historyItem.name }}</h3>
        <div v-if="!stockHistories || stockHistories.length === 0" class="no-history">
          <p>Chưa có lịch sử biến động</p>
        </div>
        <div v-else class="history-list">
          <div v-for="history in stockHistories" :key="history.id" class="history-item">
            <div class="history-header">
              <span class="history-type" :class="getHistoryTypeClass(history.type)">
                {{ getHistoryTypeText(history.type) }}
              </span>
              <span class="history-date">{{ formatDate(history.created_at) }}</span>
            </div>
            <div class="history-details">
              <div class="quantity-change">
                <span class="before">{{ history.before_qty }} {{ historyItem.unit }}</span>
                <span class="arrow">→</span>
                <span class="after">{{ history.after_qty }} {{ historyItem.unit }}</span>
                <span class="change" :class="{ positive: history.quantity > 0, negative: history.quantity < 0 }">
                  ({{ history.quantity > 0 ? '+' : '' }}{{ history.quantity }})
                </span>
              </div>
              <div class="history-reason">Lý do: {{ history.reason }}</div>
              <div class="history-user">Người thực hiện: {{ history.username }}</div>
            </div>
          </div>
        </div>
        <div class="form-actions">
          <button type="button" @click="closeHistory" class="btn-cancel">Đóng</button>
        </div>
      </div>
      </div>

      <!-- Stock Adjustment Modal -->
    <div v-if="adjustingItem" class="modal">
      <div class="modal-content">
        <h3>📦 Điều chỉnh tồn kho: {{ adjustingItem.name }}</h3>
        <p>Hiện tại: <strong>{{ adjustingItem.quantity }} {{ adjustingItem.unit }}</strong></p>
        <form @submit.prevent="saveStockAdjustment">
          <div class="form-group">
            <label>Số lượng thay đổi *</label>
            <input v-model.number="stockForm.quantity" type="number" step="0.1" required 
                   placeholder="Nhập số dương để thêm, số âm để trừ" />
            <small>Ví dụ: +10 để thêm 10 đơn vị, -5 để trừ 5 đơn vị</small>
          </div>
          <div class="form-group">
            <label>Lý do *</label>
            <select v-model="stockForm.reason" required>
              <option value="">Chọn lý do</option>
              <option value="Nhập hàng">Nhập hàng</option>
              <option value="Hỏng hóc">Hỏng hóc</option>
              <option value="Kiểm kê">Kiểm kê</option>
              <option value="Khác">Khác</option>
            </select>
          </div>
          <div class="form-group" v-if="stockForm.reason === 'Khác'">
            <label>Lý do cụ thể</label>
            <input v-model="stockForm.custom_reason" type="text" placeholder="Nhập lý do cụ thể" />
          </div>
          <div class="preview-result">
            <strong>Sau điều chỉnh: {{ (adjustingItem.quantity + (stockForm.quantity || 0)).toFixed(1) }} {{ adjustingItem.unit }}</strong>
          </div>
          <div class="form-actions">
            <button type="button" @click="closeStockAdjustment" class="btn-cancel">Hủy</button>
            <button type="submit" class="btn-save">Lưu thay đổi</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useIngredientStore } from '../stores/ingredient'
import Navigation from '../components/Navigation.vue'

const ingredientStore = useIngredientStore()

const showCreateForm = ref(false)
const showCategoryForm = ref(false)
const editingItem = ref(null)
const adjustingItem = ref(null)
const historyItem = ref(null)
const stockHistories = ref([])
const form = ref({
  name: '',
  category: '',
  unit: '',
  quantity: 0,
  min_stock: 0,
  cost_per_unit: 0,
  supplier: ''
})
const stockForm = ref({
  quantity: 0,
  reason: ''
})

const categoryForm = ref({
  name: ''
})

// Ingredient Categories (FR-IM-09)
const ingredientCategories = ref([
  { id: '1', name: 'Cà phê' },
  { id: '2', name: 'Trà' },
  { id: '3', name: 'Sữa' },
  { id: '4', name: 'Đường' },
  { id: '5', name: 'Trái cây' },
  { id: '6', name: 'Bánh' },
  { id: '7', name: 'Khác' }
])

const searchQuery = ref('')
const filterCategory = ref('')
const filterStatus = ref('')

const items = computed(() => ingredientStore.items || [])
const lowStockItems = computed(() => ingredientStore.lowStockItems || [])
const loading = computed(() => ingredientStore.loading)
const error = computed(() => ingredientStore.error)

const categories = computed(() => {
  const itemsArray = items.value || []
  const cats = [...new Set(itemsArray.map(item => item?.category).filter(Boolean))]
  return cats.sort()
})

const filteredItems = computed(() => {
  let filtered = items.value || []
  
  // Search filter
  if (searchQuery.value) {
    filtered = filtered.filter(item => 
      item?.name?.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }
  
  // Category filter
  if (filterCategory.value) {
    filtered = filtered.filter(item => item?.category === filterCategory.value)
  }
  
  // Status filter
  if (filterStatus.value) {
    filtered = filtered.filter(item => {
      if (!item || typeof item.quantity !== 'number' || typeof item.min_stock !== 'number') return false
      
      if (filterStatus.value === 'low-stock') {
        return item.quantity <= item.min_stock && item.quantity > 0
      } else if (filterStatus.value === 'out-of-stock') {
        return item.quantity === 0
      } else if (filterStatus.value === 'in-stock') {
        return item.quantity > item.min_stock
      }
      return true
    })
  }
  
  return filtered
})

const groupedItems = computed(() => {
  const groups = {}
  const filtered = filteredItems.value || []
  filtered.forEach(item => {
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

onMounted(async () => {
  try {
    await ingredientStore.fetchIngredients()
    await ingredientStore.fetchLowStock()
  } catch (error) {
    console.error('Error loading ingredients:', error)
  }
})

const formatPrice = (price) => {
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND'
  }).format(price)
}

const isLowStock = (item) => {
  return item.quantity <= item.min_stock
}

const editItem = (item) => {
  editingItem.value = item
  form.value = {
    name: item.name,
    category: item.category,
    unit: item.unit,
    quantity: item.quantity,
    min_stock: item.min_stock,
    cost_per_unit: item.cost_per_unit,
    supplier: item.supplier
  }
  showCreateForm.value = false
}

const cancelEdit = () => {
  showCreateForm.value = false
  editingItem.value = null
  form.value = { name: '', category: '', unit: '', quantity: 0, min_stock: 0, cost_per_unit: 0, supplier: '' }
}

const saveItem = async () => {
  try {
    let success = false
    
    if (editingItem.value) {
      success = await ingredientStore.updateIngredient(editingItem.value.id, form.value)
    } else {
      success = await ingredientStore.createIngredient(form.value)
    }
    
    if (success) {
      cancelEdit()
      await ingredientStore.fetchLowStock()
    }
  } catch (error) {
    console.error('Error saving ingredient:', error)
  }
}

const showAdjustStock = (item) => {
  adjustingItem.value = item
  stockForm.value = { quantity: 0, reason: '' }
}

const cancelStockAdjustment = () => {
  adjustingItem.value = null
  stockForm.value = { quantity: 0, reason: '' }
}

const saveStockAdjustment = async () => {
  try {
    console.log('Saving stock adjustment:', adjustingItem.value.id, stockForm.value)
    const success = await ingredientStore.adjustStock(adjustingItem.value.id, stockForm.value)
    if (success) {
      console.log('Stock adjustment successful')
      cancelStockAdjustment()
      await ingredientStore.fetchLowStock()
    } else {
      console.error('Stock adjustment failed:', ingredientStore.error)
      alert('Lỗi: ' + ingredientStore.error)
    }
  } catch (error) {
    console.error('Error adjusting stock:', error)
    alert('Lỗi điều chỉnh tồn kho: ' + error.message)
  }
}

const deleteItem = async (id) => {
  if (confirm('Bạn có chắc muốn xóa nguyên liệu này? Hành động này không thể hoàn tác.')) {
    try {
      const success = await ingredientStore.deleteIngredient(id)
      if (success) {
        await ingredientStore.fetchLowStock()
      }
    } catch (error) {
      console.error('Error deleting ingredient:', error)
      alert('Lỗi xóa nguyên liệu: ' + error.message)
    }
  }
}

const showHistory = async (item) => {
  try {
    historyItem.value = item
    stockHistories.value = await ingredientStore.fetchStockHistory(item.id) || []
  } catch (error) {
    console.error('Error fetching stock history:', error)
    stockHistories.value = []
  }
}

const closeHistory = () => {
  historyItem.value = null
  stockHistories.value = []
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleString('vi-VN')
}

const getHistoryTypeText = (type) => {
  const types = {
    'adjustment': 'Điều chỉnh',
    'order': 'Sử dụng',
    'purchase': 'Nhập hàng',
    'waste': 'Hỏng hóc'
  }
  return types[type] || type
}

const getHistoryTypeClass = (type) => {
  const classes = {
    'adjustment': 'type-adjustment',
    'order': 'type-order',
    'purchase': 'type-purchase',
    'waste': 'type-waste'
  }
  return classes[type] || 'type-default'
}

const getCategoryIcon = (category) => {
  const iconMap = {
    'Cà phê': '☕',
    'Trà': '🍵',
    'Sữa': '🥛',
    'Đường': '🍯',
    'Trái cây': '🍎',
    'Bánh': '🍰',
    'Khác': '📦'
  }
  return iconMap[category] || '📦'
}

const getCategoryColor = (category) => {
  const colorMap = {
    'Cà phê': 'bg-amber-100 text-amber-600',
    'Trà': 'bg-green-100 text-green-600',
    'Sữa': 'bg-blue-100 text-blue-600',
    'Đường': 'bg-pink-100 text-pink-600',
    'Trái cây': 'bg-orange-100 text-orange-600',
    'Bánh': 'bg-yellow-100 text-yellow-600',
    'Khác': 'bg-gray-100 text-gray-600'
  }
  return colorMap[category] || 'bg-gray-100 text-gray-600'
}

const getStockStatus = (item) => {
  if (item.quantity === 0) return 'Hết hàng'
  if (item.quantity <= item.min_stock) return 'Sắp hết'
  return 'Còn hàng'
}

const getStockBadge = (item) => {
  if (item.quantity === 0) return 'bg-red-100 text-red-800'
  if (item.quantity <= item.min_stock) return 'bg-yellow-100 text-yellow-800'
  return 'bg-green-100 text-green-800'
}

const isManager = computed(() => {
  // Add your role check logic here
  return true // For now, assume manager role
})

const addCategory = () => {
  if (!categoryForm.value.name) return
  
  ingredientCategories.value.push({
    id: Date.now().toString(),
    name: categoryForm.value.name
  })
  
  categoryForm.value = { name: '' }
}

const deleteCategory = (id, name) => {
  const hasIngredients = items.value.some(item => item.category === name)
  
  if (hasIngredients) {
    alert('Không thể xóa danh mục đã có nguyên liệu!')
    return
  }
  
  if (confirm(`Bạn có chắc muốn xóa danh mục "${name}"?`)) {
    ingredientCategories.value = ingredientCategories.value.filter(c => c.id !== id)
  }
}

const getIngredientCountByCategory = (categoryName) => {
  return items.value.filter(item => item.category === categoryName).length
}
</script>

<style scoped>
.ingredient-management {
  min-height: 100vh;
  background: #f5f6fa;
}

.content {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.header h2 {
  color: #333;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.btn-primary, .btn-warning {
  border: none;
  padding: 12px 20px;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
}

.btn-primary {
  background: #667eea;
  color: white;
}

.btn-warning {
  background: #ffc107;
  color: #212529;
}

.low-stock-alert {
  background: #fff3cd;
  border: 1px solid #ffeaa7;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 20px;
}

.low-stock-alert h3 {
  color: #856404;
  margin: 0 0 10px 0;
}

.low-stock-items {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.low-stock-item {
  background: #ffeaa7;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  color: #856404;
}

.ingredient-grid {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.category-section {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.category-title {
  color: #333;
  margin: 0 0 20px 0;
  padding-bottom: 10px;
  border-bottom: 2px solid #667eea;
  font-size: 20px;
}

.category-items {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 15px;
}

.ingredient-card {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 15px;
  border: 1px solid #e9ecef;
}

.ingredient-info h4 {
  color: #333;
  margin: 0 0 10px 0;
  font-size: 16px;
}

.stock-info {
  margin-bottom: 8px;
}

.current-stock {
  font-size: 18px;
  color: #28a745;
}

.current-stock.low-stock {
  color: #dc3545;
}

.min-stock {
  font-size: 12px;
  color: #6c757d;
}

.cost-info, .supplier {
  font-size: 12px;
  color: #666;
  margin: 4px 0;
}

.ingredient-actions {
  display: flex;
  gap: 8px;
  margin-top: 10px;
}

.ingredient-actions button {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 600;
  flex: 1;
}

.btn-history {
  background: #6f42c1;
  color: white;
}

.btn-adjust {
  background: #17a2b8;
  color: white;
}

.btn-edit {
  background: #ffc107;
  color: #212529;
}

.btn-delete {
  background: #dc3545;
  color: white;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 30px;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  overflow-y: auto;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #333;
  font-weight: 500;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 10px;
  border: 2px solid #e1e5e9;
  border-radius: 6px;
  font-size: 14px;
}

.form-group small {
  color: #6c757d;
  font-size: 12px;
}

.preview-result {
  background: #e7f3ff;
  padding: 10px;
  border-radius: 6px;
  margin: 15px 0;
  text-align: center;
  color: #0066cc;
}

.form-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

.btn-cancel {
  background: #6c757d;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
}

.btn-save {
  background: #28a745;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
}

.filters-section {
  background: white;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.filters-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr;
  gap: 15px;
  align-items: end;
}

.search-input, .filter-select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

@media (max-width: 768px) {
  .filters-row {
    grid-template-columns: 1fr;
  }
}

.history-modal {
  max-width: 600px;
}

.no-history {
  text-align: center;
  color: #6c757d;
  font-style: italic;
  padding: 40px 20px;
}

.history-list {
  max-height: 400px;
  overflow-y: auto;
}

.history-item {
  border: 1px solid #e9ecef;
  border-radius: 6px;
  margin-bottom: 10px;
  padding: 15px;
  background: #f8f9fa;
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.history-type {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
}

.type-adjustment { background: #e7f3ff; color: #0066cc; }
.type-order { background: #fff3cd; color: #856404; }
.type-purchase { background: #d4edda; color: #155724; }
.type-waste { background: #f8d7da; color: #721c24; }

.history-date {
  font-size: 12px;
  color: #6c757d;
}

.quantity-change {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 5px;
  font-weight: 600;
}

.before, .after {
  background: #e9ecef;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 14px;
}

.arrow {
  color: #6c757d;
}

.change.positive {
  color: #28a745;
}

.change.negative {
  color: #dc3545;
}

.history-reason, .history-user {
  font-size: 12px;
  color: #666;
  margin: 2px 0;
}
</style>