<template>
  <div class="min-h-screen bg-gray-100">
    <Navigation />
    <div class="p-4">
      <div class="flex flex-col lg:flex-row justify-between items-center mb-6">
        <h2 class="text-xl lg:text-2xl font-semibold text-gray-800 mb-4 lg:mb-0">
          🍽️ Quản lý Menu
        </h2>
        <div class="flex flex-wrap gap-2">
          <button @click="showCategoryForm = true" class="bg-purple-500 hover:bg-purple-600 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors">
            📁 Quản lý danh mục
          </button>
          <button @click="showCreateForm = true" class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors">
            + Thêm món mới
          </button>
        </div>
      </div>

      <div v-if="loading" class="text-center py-10 text-gray-600 text-lg">Đang tải...</div>
      <div v-if="error" class="text-center py-10 text-red-600 bg-red-50 border border-red-200 rounded-lg">{{ error }}</div>

      <div class="grid grid-cols-1 gap-4">
        <div v-for="category in groupedItems" :key="category.name" class="bg-white rounded-xl p-4 shadow-sm">
          <h3 class="text-lg font-bold text-gray-800 mb-4 pb-2 border-b-2 border-blue-500">{{ category.name }}</h3>
          <div class="space-y-3">
            <div v-for="item in category.items" :key="item.id" class="rounded-xl p-4 bg-gray-50">
              <div class="flex items-center justify-between mb-3">
                <div class="flex items-center space-x-3">
                  <div class="w-12 h-12 rounded-xl flex items-center justify-center text-2xl" :class="getCategoryColor(item.category)">
                    {{ getCategoryIcon(item.category) }}
                  </div>
                  <div>
                    <h4 class="font-bold text-gray-800">{{ item.name }}</h4>
                    <p class="text-sm text-gray-500">{{ item.description || 'Chưa có mô tả' }}</p>
                  </div>
                </div>
                <div class="text-right">
                  <span class="px-3 py-1 rounded-full text-xs font-medium" :class="item.available ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'">
                    {{ item.available ? 'Có sẵn' : 'Hết hàng' }}
                  </span>
                </div>
              </div>

              <div class="grid grid-cols-1 gap-3 mb-4">
                <div class="bg-white rounded-lg p-3">
                  <div class="text-2xl font-bold text-green-600">{{ formatPrice(item.price) }}</div>
                  <div class="text-xs text-gray-500">Giá bán</div>
                </div>
                <div v-if="item.ingredients && item.ingredients.length > 0" class="bg-white rounded-lg p-3">
                  <div class="text-sm font-semibold text-gray-700 mb-2">Nguyên liệu:</div>
                  <ul class="text-xs text-gray-600 space-y-1">
                    <li v-for="ingredient in item.ingredients" :key="ingredient.name">
                      • {{ ingredient.name }}: {{ ingredient.quantity }} {{ ingredient.unit }}
                    </li>
                  </ul>
                </div>
              </div>

              <div class="grid grid-cols-3 gap-2">
                <button @click="editItem(item)" class="bg-yellow-500 hover:bg-yellow-600 text-white px-3 py-2 rounded-lg text-sm font-medium transition-colors">
                  📝 Sửa
                </button>
                <button @click="toggleAvailable(item)" :class="item.available ? 'bg-gray-500 hover:bg-gray-600' : 'bg-green-500 hover:bg-green-600'" class="text-white px-3 py-2 rounded-lg text-sm font-medium transition-colors">
                  {{ item.available ? '🙈 Ẩn' : '👁️ Hiện' }}
                </button>
                <button @click="deleteItem(item.id)" class="bg-red-500 hover:bg-red-600 text-white px-3 py-2 rounded-lg text-sm font-medium transition-colors">
                  🗑️ Xóa
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Category Management Modal -->
      <div v-if="showCategoryForm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-xl p-6 w-full max-w-md max-h-[80vh] overflow-y-auto">
          <h3 class="text-xl font-bold text-gray-800 mb-4">📁 Quản lý Danh mục Menu</h3>
          
          <div class="bg-gray-50 rounded-lg p-4 mb-4">
            <h4 class="font-semibold text-gray-800 mb-3">Thêm danh mục mới</h4>
            <form @submit.prevent="addCategory">
              <input v-model="categoryForm.name" type="text" required placeholder="Ví dụ: Cà phê" class="w-full p-3 border border-gray-300 rounded-lg mb-3" />
              <button type="submit" class="w-full bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg font-medium">+ Thêm danh mục</button>
            </form>
          </div>

          <div class="space-y-2 max-h-96 overflow-y-auto">
            <div v-for="cat in menuCategories" :key="cat.id" class="bg-white border border-gray-200 rounded-lg p-3 flex items-center justify-between">
              <div class="flex items-center space-x-3">
                <div class="w-10 h-10 rounded-lg flex items-center justify-center text-xl" :class="getCategoryColor(cat.name)">
                  {{ getCategoryIcon(cat.name) }}
                </div>
                <div>
                  <div class="font-medium text-gray-800">{{ cat.name }}</div>
                  <div class="text-xs text-gray-500">{{ getMenuCountByCategory(cat.name) }} món</div>
                </div>
              </div>
              <button @click="deleteCategory(cat.id, cat.name)" class="text-red-500 hover:text-red-700 p-2">
                🗑️
              </button>
            </div>
          </div>

          <div class="mt-4">
            <button type="button" @click="showCategoryForm = false" class="w-full bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg font-medium">Đóng</button>
          </div>
        </div>
      </div>

      <!-- Create/Edit Form Modal -->
      <div v-if="showCreateForm || editingItem" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-xl p-6 w-full max-w-md max-h-[80vh] overflow-y-auto">
          <h3 class="text-xl font-bold text-gray-800 mb-4">{{ editingItem ? 'Sửa món' : 'Thêm món mới' }}</h3>
          <form @submit.prevent="saveItem" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Tên món *</label>
              <input v-model="form.name" type="text" required placeholder="Nhập tên món" class="w-full p-3 border border-gray-300 rounded-lg" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Danh mục *</label>
              <select v-model="form.category" required class="w-full p-3 border border-gray-300 rounded-lg">
                <option value="">Chọn danh mục</option>
                <option v-for="cat in menuCategories" :key="cat.id" :value="cat.name">{{ cat.name }}</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Giá (VNĐ) *</label>
              <input v-model.number="form.price" type="number" min="0" step="1000" required placeholder="0" class="w-full p-3 border border-gray-300 rounded-lg" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Mô tả</label>
              <textarea v-model="form.description" rows="3" placeholder="Mô tả món ăn..." class="w-full p-3 border border-gray-300 rounded-lg"></textarea>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Nguyên liệu</label>
              <div class="border border-gray-300 rounded-lg p-3 bg-gray-50">
                <div v-if="form.ingredients.length === 0" class="text-center text-gray-500 italic py-4">
                  Chưa có nguyên liệu nào
                </div>
                <div v-for="(ingredient, index) in form.ingredients" :key="index" class="flex gap-2 mb-2">
                  <input v-model="ingredient.name" placeholder="Tên" class="flex-1 p-2 border border-gray-300 rounded" required />
                  <input v-model.number="ingredient.quantity" type="number" min="0" step="0.1" placeholder="SL" class="w-20 p-2 border border-gray-300 rounded" required />
                  <input v-model="ingredient.unit" placeholder="Đơn vị" class="w-20 p-2 border border-gray-300 rounded" required />
                  <button type="button" @click="removeIngredient(index)" class="bg-red-500 text-white px-3 rounded hover:bg-red-600">×</button>
                </div>
                <button type="button" @click="addIngredient" class="w-full bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg mt-2">
                  + Thêm nguyên liệu
                </button>
              </div>
            </div>
            <div class="flex gap-2">
              <button type="button" @click="cancelEdit" class="flex-1 bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg font-medium">Hủy</button>
              <button type="submit" class="flex-1 bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded-lg font-medium" :disabled="!form.name || !form.category || form.price <= 0">
                {{ editingItem ? 'Cập nhật' : 'Thêm món' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useMenuStore } from '../stores/menu'
import Navigation from '../components/Navigation.vue'

const menuStore = useMenuStore()

const showCreateForm = ref(false)
const showCategoryForm = ref(false)
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

const menuCategories = ref([
  { id: '1', name: 'Cà phê' },
  { id: '2', name: 'Trà' },
  { id: '3', name: 'Nước ép' },
  { id: '4', name: 'Bánh ngọt' },
  { id: '5', name: 'Món nhẹ' }
])

const items = computed(() => menuStore.items)
const loading = computed(() => menuStore.loading)
const error = computed(() => menuStore.error)

const groupedItems = computed(() => {
  const groups = {}
  items.value.forEach(item => {
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

onMounted(() => {
  menuStore.fetchMenuItems()
})

const formatPrice = (price) => {
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND'
  }).format(price)
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
  showCreateForm.value = false
}

const cancelEdit = () => {
  showCreateForm.value = false
  editingItem.value = null
  form.value = { name: '', category: '', price: 0, description: '', ingredients: [] }
}

const addIngredient = () => {
  form.value.ingredients.push({ name: '', quantity: 0, unit: '' })
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
    if (!success && menuStore.error) {
      alert('Lỗi: ' + menuStore.error)
    }
  }
}

const addCategory = () => {
  if (!categoryForm.value.name) return
  
  menuCategories.value.push({
    id: Date.now().toString(),
    name: categoryForm.value.name
  })
  
  categoryForm.value = { name: '' }
}

const deleteCategory = (id, name) => {
  const hasMenuItems = items.value.some(item => item.category === name)
  
  if (hasMenuItems) {
    alert('Không thể xóa danh mục đã có món!')
    return
  }
  
  if (confirm(`Bạn có chắc muốn xóa danh mục "${name}"?`)) {
    menuCategories.value = menuCategories.value.filter(c => c.id !== id)
  }
}

const getMenuCountByCategory = (categoryName) => {
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
button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
