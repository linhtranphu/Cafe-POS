<template>
  <div class="min-h-screen bg-gray-100">
    <Navigation />
    <div class="p-4">
      <!-- Header -->
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold text-gray-800">👥 Quản lý User</h2>
        <button @click="showCreateForm = true" class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg font-medium">
          + Thêm User
        </button>
      </div>

      <!-- Filter Tabs -->
      <div class="flex gap-2 mb-4 overflow-x-auto">
        <button v-for="filter in filters" :key="filter.value" @click="currentFilter = filter.value"
          :class="currentFilter === filter.value ? 'bg-blue-500 text-white' : 'bg-white text-gray-700'"
          class="px-4 py-2 rounded-lg font-medium whitespace-nowrap border">
          {{ filter.label }} ({{ getUserCountByFilter(filter.value) }})
        </button>
      </div>

      <!-- Users List -->
      <div v-if="loading" class="text-center py-10">Đang tải...</div>
      <div v-else-if="filteredUsers.length === 0" class="text-center py-10 text-gray-500">Không có user nào</div>
      <div v-else class="grid gap-4">
        <div v-for="user in filteredUsers" :key="user.id" class="bg-white rounded-xl p-4 shadow-sm">
          <div class="flex justify-between items-start mb-3">
            <div class="flex-1">
              <div class="flex items-center gap-3 mb-2">
                <h3 class="font-bold text-lg">{{ user.name }}</h3>
                <span :class="getRoleColor(user.role)" class="px-2 py-1 rounded-full text-xs font-medium">
                  {{ getRoleText(user.role) }}
                </span>
                <span :class="user.active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'" 
                  class="px-2 py-1 rounded-full text-xs font-medium">
                  {{ user.active ? 'Hoạt động' : 'Tạm khóa' }}
                </span>
              </div>
              <p class="text-sm text-gray-600 mb-1">Username: {{ user.username }}</p>
              <p class="text-sm text-gray-500">Tạo: {{ formatDate(user.created_at) }}</p>
              <p v-if="user.last_login" class="text-sm text-gray-500">Đăng nhập cuối: {{ formatDate(user.last_login) }}</p>
            </div>
            
            <!-- Actions -->
            <div class="flex gap-2">
              <button @click="showEditForm(user)" class="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded text-sm">
                ✏️ Sửa
              </button>
              <button @click="showResetPasswordForm(user)" class="bg-yellow-500 hover:bg-yellow-600 text-white px-3 py-1 rounded text-sm">
                🔑 Reset PW
              </button>
              <button @click="toggleUserStatus(user)" 
                :class="user.active ? 'bg-orange-500 hover:bg-orange-600' : 'bg-green-500 hover:bg-green-600'"
                class="text-white px-3 py-1 rounded text-sm">
                {{ user.active ? '⏸️ Khóa' : '▶️ Mở' }}
              </button>
              <button @click="showDeleteConfirm(user)" class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded text-sm">
                🗑️ Xóa
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Create User Modal -->
      <div v-if="showCreateForm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-xl p-6 w-full max-w-md">
          <h3 class="text-xl font-bold mb-4">Thêm User Mới</h3>
          <form @submit.prevent="createUser" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Username *</label>
              <input v-model="createForm.username" type="text" required minlength="3" maxlength="50"
                class="w-full p-3 border rounded-lg" placeholder="Nhập username">
            </div>
            
            <div>
              <label class="block text-sm font-medium mb-2">Mật khẩu *</label>
              <input v-model="createForm.password" type="password" required minlength="6"
                class="w-full p-3 border rounded-lg" placeholder="Nhập mật khẩu">
            </div>
            
            <div>
              <label class="block text-sm font-medium mb-2">Tên hiển thị *</label>
              <input v-model="createForm.name" type="text" required minlength="2" maxlength="100"
                class="w-full p-3 border rounded-lg" placeholder="Nhập tên hiển thị">
            </div>
            
            <div>
              <label class="block text-sm font-medium mb-2">Vai trò *</label>
              <select v-model="createForm.role" required class="w-full p-3 border rounded-lg">
                <option value="">-- Chọn vai trò --</option>
                <option value="manager">Manager</option>
                <option value="cashier">Cashier</option>
                <option value="waiter">Waiter</option>
              </select>
            </div>
            
            <div class="flex items-center">
              <input v-model="createForm.active" type="checkbox" id="active" class="mr-2">
              <label for="active" class="text-sm">Kích hoạt ngay</label>
            </div>
            
            <div class="flex gap-2">
              <button type="button" @click="showCreateForm = false" class="flex-1 bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg">
                Hủy
              </button>
              <button type="submit" class="flex-1 bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg">
                Tạo User
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Edit User Modal -->
      <div v-if="showEditModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-xl p-6 w-full max-w-md">
          <h3 class="text-xl font-bold mb-4">Chỉnh sửa User</h3>
          <form @submit.prevent="updateUser" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Username</label>
              <input :value="editForm.username" disabled class="w-full p-3 border rounded-lg bg-gray-100">
            </div>
            
            <div>
              <label class="block text-sm font-medium mb-2">Tên hiển thị *</label>
              <input v-model="editForm.name" type="text" required minlength="2" maxlength="100"
                class="w-full p-3 border rounded-lg" placeholder="Nhập tên hiển thị">
            </div>
            
            <div>
              <label class="block text-sm font-medium mb-2">Vai trò *</label>
              <select v-model="editForm.role" required class="w-full p-3 border rounded-lg">
                <option value="manager">Manager</option>
                <option value="cashier">Cashier</option>
                <option value="waiter">Waiter</option>
              </select>
            </div>
            
            <div class="flex items-center">
              <input v-model="editForm.active" type="checkbox" id="editActive" class="mr-2">
              <label for="editActive" class="text-sm">Kích hoạt</label>
            </div>
            
            <div class="flex gap-2">
              <button type="button" @click="showEditModal = false" class="flex-1 bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg">
                Hủy
              </button>
              <button type="submit" class="flex-1 bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg">
                Cập nhật
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Reset Password Modal -->
      <div v-if="showResetPassword" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-xl p-6 w-full max-w-md">
          <h3 class="text-xl font-bold mb-4">Reset Mật khẩu</h3>
          <p class="text-gray-600 mb-4">Reset mật khẩu cho: <strong>{{ selectedUser?.name }}</strong></p>
          <form @submit.prevent="resetPassword" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Mật khẩu mới *</label>
              <input v-model="resetPasswordForm.newPassword" type="password" required minlength="6"
                class="w-full p-3 border rounded-lg" placeholder="Nhập mật khẩu mới">
            </div>
            
            <div class="flex gap-2">
              <button type="button" @click="showResetPassword = false" class="flex-1 bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg">
                Hủy
              </button>
              <button type="submit" class="flex-1 bg-yellow-500 hover:bg-yellow-600 text-white px-4 py-2 rounded-lg">
                Reset Password
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Delete Confirmation Modal -->
      <div v-if="showDeleteConfirmation" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-xl p-6 w-full max-w-md">
          <h3 class="text-xl font-bold mb-4 text-red-600">Xác nhận xóa</h3>
          <p class="text-gray-600 mb-4">Bạn có chắc chắn muốn xóa user: <strong>{{ selectedUser?.name }}</strong>?</p>
          <p class="text-sm text-red-500 mb-4">⚠️ Hành động này không thể hoàn tác!</p>
          
          <div class="flex gap-2">
            <button @click="showDeleteConfirmation = false" class="flex-1 bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg">
              Hủy
            </button>
            <button @click="deleteUser" class="flex-1 bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-lg">
              Xóa User
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '../stores/user'
import Navigation from '../components/Navigation.vue'

const userStore = useUserStore()

const showCreateForm = ref(false)
const showEditModal = ref(false)
const showResetPassword = ref(false)
const showDeleteConfirmation = ref(false)
const selectedUser = ref(null)
const currentFilter = ref('ALL')

const createForm = ref({
  username: '',
  password: '',
  name: '',
  role: '',
  active: true
})

const editForm = ref({
  id: '',
  username: '',
  name: '',
  role: '',
  active: true
})

const resetPasswordForm = ref({
  newPassword: ''
})

const filters = [
  { value: 'ALL', label: 'Tất cả' },
  { value: 'ACTIVE', label: 'Hoạt động' },
  { value: 'INACTIVE', label: 'Tạm khóa' },
  { value: 'MANAGER', label: 'Manager' },
  { value: 'CASHIER', label: 'Cashier' },
  { value: 'WAITER', label: 'Waiter' }
]

const loading = computed(() => userStore.loading)
const users = computed(() => userStore.users)

const filteredUsers = computed(() => {
  switch (currentFilter.value) {
    case 'ACTIVE':
      return users.value.filter(u => u.active)
    case 'INACTIVE':
      return users.value.filter(u => !u.active)
    case 'MANAGER':
      return users.value.filter(u => u.role === 'manager')
    case 'CASHIER':
      return users.value.filter(u => u.role === 'cashier')
    case 'WAITER':
      return users.value.filter(u => u.role === 'waiter')
    default:
      return users.value
  }
})

onMounted(async () => {
  await userStore.fetchUsers()
})

const createUser = async () => {
  try {
    await userStore.createUser(createForm.value)
    showCreateForm.value = false
    createForm.value = { username: '', password: '', name: '', role: '', active: true }
    alert('Tạo user thành công!')
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

const showEditForm = (user) => {
  selectedUser.value = user
  editForm.value = {
    id: user.id,
    username: user.username,
    name: user.name,
    role: user.role,
    active: user.active
  }
  showEditModal.value = true
}

const updateUser = async () => {
  try {
    await userStore.updateUser(editForm.value.id, {
      name: editForm.value.name,
      role: editForm.value.role,
      active: editForm.value.active
    })
    showEditModal.value = false
    selectedUser.value = null
    alert('Cập nhật user thành công!')
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const showResetPasswordForm = (user) => {
  selectedUser.value = user
  resetPasswordForm.value.newPassword = ''
  showResetPassword.value = true
}

const resetPassword = async () => {
  try {
    await userStore.resetPassword(selectedUser.value.id, resetPasswordForm.value.newPassword)
    showResetPassword.value = false
    selectedUser.value = null
    resetPasswordForm.value.newPassword = ''
    alert('Reset password thành công!')
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const toggleUserStatus = async (user) => {
  try {
    await userStore.toggleUserStatus(user.id)
    alert(`${user.active ? 'Khóa' : 'Mở khóa'} user thành công!`)
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const showDeleteConfirm = (user) => {
  selectedUser.value = user
  showDeleteConfirmation.value = true
}

const deleteUser = async () => {
  try {
    await userStore.deleteUser(selectedUser.value.id)
    showDeleteConfirmation.value = false
    selectedUser.value = null
    alert('Xóa user thành công!')
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const getUserCountByFilter = (filter) => {
  switch (filter) {
    case 'ACTIVE':
      return users.value.filter(u => u.active).length
    case 'INACTIVE':
      return users.value.filter(u => !u.active).length
    case 'MANAGER':
      return users.value.filter(u => u.role === 'manager').length
    case 'CASHIER':
      return users.value.filter(u => u.role === 'cashier').length
    case 'WAITER':
      return users.value.filter(u => u.role === 'waiter').length
    default:
      return users.value.length
  }
}

const getRoleColor = (role) => {
  const colors = {
    manager: 'bg-purple-100 text-purple-800',
    cashier: 'bg-blue-100 text-blue-800',
    waiter: 'bg-green-100 text-green-800'
  }
  return colors[role] || 'bg-gray-100 text-gray-800'
}

const getRoleText = (role) => {
  const texts = {
    manager: 'Manager',
    cashier: 'Cashier',
    waiter: 'Waiter'
  }
  return texts[role] || role
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('vi-VN')
}
</script>

<style scoped>
button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>