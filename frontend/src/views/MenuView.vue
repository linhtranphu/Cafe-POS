<template>
  <div class="menu-management">
    <Navigation />
    <div class="content">
      <div class="header">
      <h2>🍽️ Quản lý Menu</h2>
      <button @click="showCreateForm = true" class="btn-primary">+ Thêm món mới</button>
    </div>

    <div v-if="loading" class="loading">Đang tải...</div>
    <div v-if="error" class="error">{{ error }}</div>

    <div class="menu-grid">
      <div v-for="category in groupedItems" :key="category.name" class="category-section">
        <h3 class="category-title">{{ category.name }}</h3>
        <div class="category-items">
          <div v-for="item in category.items" :key="item.id" class="menu-card">
            <div class="menu-info">
              <h4>{{ item.name }}</h4>
              <p class="description">{{ item.description }}</p>
              <div class="price">{{ formatPrice(item.price) }}</div>
              <div class="ingredients" v-if="item.ingredients && item.ingredients.length > 0">
                <strong>Nguyên liệu:</strong>
                <ul>
                  <li v-for="ingredient in item.ingredients" :key="ingredient.name">
                    {{ ingredient.name }}: {{ ingredient.quantity }} {{ ingredient.unit }}
                  </li>
                </ul>
              </div>
              <div class="status" :class="{ available: item.available, unavailable: !item.available }">
                {{ item.available ? 'Có sẵn' : 'Hết hàng' }}
              </div>
            </div>
            <div class="menu-actions">
              <button @click="editItem(item)" class="btn-edit" title="Sửa món">
                📝 Sửa
              </button>
              <button @click="toggleAvailable(item)" 
                      :class="item.available ? 'btn-hide' : 'btn-show'" 
                      :title="item.available ? 'Ẩn món' : 'Hiện món'">
                {{ item.available ? '🙈 Ẩn' : '👁️ Hiện' }}
              </button>
              <button @click="deleteItem(item.id)" class="btn-delete" title="Xóa món">
                🗑️ Xóa
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Form Modal -->
    <div v-if="showCreateForm || editingItem" class="modal">
      <div class="modal-content">
        <h3>{{ editingItem ? 'Sửa món' : 'Thêm món mới' }}</h3>
        <form @submit.prevent="saveItem">
          <div class="form-group">
            <label>Tên món *</label>
            <input v-model="form.name" type="text" required placeholder="Nhập tên món" />
          </div>
          <div class="form-group">
            <label>Danh mục *</label>
            <select v-model="form.category" required>
              <option value="">Chọn danh mục</option>
              <option value="Cà phê">Cà phê</option>
              <option value="Trà">Trà</option>
              <option value="Nước ép">Nước ép</option>
              <option value="Bánh ngọt">Bánh ngọt</option>
              <option value="Món nhẹ">Món nhẹ</option>
            </select>
          </div>
          <div class="form-group">
            <label>Giá (VNĐ) *</label>
            <input v-model.number="form.price" type="number" min="0" step="1000" required placeholder="0" />
          </div>
          <div class="form-group">
            <label>Mô tả</label>
            <textarea v-model="form.description" rows="3" placeholder="Mô tả món ăn..."></textarea>
          </div>
          <div class="form-group">
            <label>Nguyên liệu</label>
            <div class="ingredients-section">
              <div v-if="form.ingredients.length === 0" class="no-ingredients">
                <p>Chưa có nguyên liệu nào</p>
              </div>
              <div v-for="(ingredient, index) in form.ingredients" :key="index" class="ingredient-row">
                <input v-model="ingredient.name" placeholder="Tên nguyên liệu" class="ingredient-name" required />
                <input v-model.number="ingredient.quantity" type="number" min="0" step="0.1" placeholder="Số lượng" class="ingredient-quantity" required />
                <input v-model="ingredient.unit" placeholder="Đơn vị (g, ml, cái...)" class="ingredient-unit" required />
                <button type="button" @click="removeIngredient(index)" class="btn-remove" title="Xóa nguyên liệu">×</button>
              </div>
              <button type="button" @click="addIngredient" class="btn-add-ingredient">
                + Thêm nguyên liệu
              </button>
            </div>
          </div>
          <div class="form-actions">
            <button type="button" @click="cancelEdit" class="btn-cancel">Hủy</button>
            <button type="submit" class="btn-save" :disabled="!form.name || !form.category || form.price <= 0">
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
const editingItem = ref(null)
const form = ref({
  name: '',
  category: '',
  price: 0,
  description: '',
  ingredients: []
})

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
</script>

<style scoped>
.menu-management {
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

.btn-primary {
  background: #667eea;
  color: white;
  border: none;
  padding: 12px 20px;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
}

.btn-primary:hover {
  background: #5a6fd8;
}

.loading, .error {
  text-align: center;
  padding: 20px;
  margin: 20px 0;
}

.error {
  color: #e74c3c;
  background: #fdf2f2;
  border-radius: 8px;
}

.menu-grid {
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
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 15px;
}

.menu-card {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 15px;
  border: 1px solid #e9ecef;
}

.menu-info h4 {
  color: #333;
  margin: 0 0 8px 0;
  font-size: 16px;
}

.category {
  color: #667eea;
  font-weight: 600;
  font-size: 14px;
  margin: 0 0 8px 0;
  display: none;
}

.description {
  color: #666;
  font-size: 14px;
  margin: 0 0 12px 0;
}

.price {
  font-size: 18px;
  font-weight: bold;
  color: #27ae60;
  margin-bottom: 8px;
}

.status {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 15px;
}

.status.available {
  background: #d4edda;
  color: #155724;
}

.status.unavailable {
  background: #f8d7da;
  color: #721c24;
}

.menu-actions {
  display: flex;
  gap: 8px;
  margin-top: 10px;
}

.menu-actions button {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 600;
  transition: all 0.2s;
  flex: 1;
}

.btn-edit {
  background: #ffc107;
  color: #212529;
}

.btn-edit:hover {
  background: #e0a800;
  transform: translateY(-1px);
}

.btn-hide {
  background: #6c757d;
  color: white;
}

.btn-hide:hover {
  background: #545b62;
  transform: translateY(-1px);
}

.btn-show {
  background: #28a745;
  color: white;
}

.btn-show:hover {
  background: #218838;
  transform: translateY(-1px);
}

.btn-delete {
  background: #dc3545;
  color: white;
}

.btn-delete:hover {
  background: #c82333;
  transform: translateY(-1px);
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
}

.modal-content h3 {
  margin: 0 0 20px 0;
  color: #333;
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
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 10px;
  border: 2px solid #e1e5e9;
  border-radius: 6px;
  font-size: 14px;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #667eea;
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

.ingredients {
  margin: 10px 0;
  font-size: 12px;
}

.ingredients ul {
  margin: 5px 0 0 15px;
  padding: 0;
}

.ingredients li {
  margin: 2px 0;
  color: #666;
}

.ingredients-section {
  border: 1px solid #e1e5e9;
  border-radius: 6px;
  padding: 15px;
  background: #f8f9fa;
}

.ingredient-row {
  display: flex;
  gap: 8px;
  margin-bottom: 10px;
  align-items: center;
}

.ingredient-name {
  flex: 2;
  padding: 6px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.ingredient-quantity {
  flex: 1;
  padding: 6px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.ingredient-unit {
  flex: 1;
  padding: 6px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.btn-remove {
  background: #dc3545;
  color: white;
  border: none;
  padding: 6px 10px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-add-ingredient {
  background: #17a2b8;
  color: white;
  border: none;
  padding: 8px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  margin-top: 5px;
  width: 100%;
}

.btn-add-ingredient:hover {
  background: #138496;
}

.no-ingredients {
  text-align: center;
  color: #6c757d;
  font-style: italic;
  padding: 20px;
}

.btn-save:disabled {
  background: #6c757d;
  cursor: not-allowed;
  opacity: 0.6;
}

.form-group input:invalid {
  border-color: #dc3545;
}

.form-group select:invalid {
  border-color: #dc3545;
}
</style>