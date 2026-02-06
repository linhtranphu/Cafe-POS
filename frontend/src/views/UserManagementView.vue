<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Pull to Refresh Indicator -->
    <PullToRefresh 
      :pull-distance="pullDistance" 
      :is-refreshing="isRefreshing"
      :threshold="80" />
    
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3" style="padding-top: max(0.75rem, env(safe-area-inset-top))">
        <h1 class="text-xl font-bold text-gray-800">👥 Quản lý User</h1>
        <input v-model="searchQuery" type="text" placeholder="Tìm kiếm user..." class="w-full px-4 py-2 border border-gray-300 rounded-lg mt-2 focus:ring-2 focus:ring-blue-500" />
      </div>
    </div>

    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24">
      <div class="bg-gradient-to-br from-purple-500 to-pink-500 rounded-xl p-4 mb-4 text-white shadow-lg">
        <div class="text-xs opacity-90 mb-2">Tổng quan nhân viên</div>
        <div class="grid grid-cols-4 gap-1.5">
          <div class="text-center">
            <div class="text-lg font-bold">{{ users.length }}</div>
            <div class="text-[10px] opacity-90">Tổng</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold">{{ activeCount }}</div>
            <div class="text-[10px] opacity-90">Hoạt động</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold">{{ managerCount }}</div>
            <div class="text-[10px] opacity-90">Manager</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold">{{ cashierCount }}</div>
            <div class="text-[10px] opacity-90">Cashier</div>
          </div>
        </div>
      </div>

      <div class="mb-4">
        <h2 class="text-sm font-bold text-gray-800 mb-2">⚡ Thao tác nhanh</h2>
        <div class="grid grid-cols-2 gap-2">
          <button @click="openCreateModal" class="bg-gradient-to-br from-blue-500 to-cyan-500 text-white rounded-xl p-4 shadow-md">
            <div class="text-2xl mb-1">➕</div>
            <div class="text-sm font-bold">Tạo User</div>
          </button>
          <button @click="currentFilter = 'ALL'" class="bg-gradient-to-br from-purple-500 to-pink-500 text-white rounded-xl p-4 shadow-md">
            <div class="text-2xl mb-1">🔍</div>
            <div class="text-sm font-bold">Tất cả</div>
          </button>
        </div>
      </div>

      <div v-if="loading" class="text-center py-10">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
      </div>
      <div v-else-if="filteredUsers.length === 0" class="text-center py-16">
        <div class="text-6xl mb-4">👥</div>
        <p class="text-gray-500">Không có user nào</p>
      </div>
      <div v-else class="space-y-3">
        <div v-for="user in filteredUsers" :key="user.id" class="bg-white rounded-2xl p-4 shadow-sm">
          <div class="flex justify-between items-start mb-3">
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-2">
                <h3 class="font-bold text-lg">{{ user.name }}</h3>
                <span :class="getRoleColor(user.role)" class="px-2 py-0.5 rounded-full text-xs font-medium">{{ getRoleText(user.role) }}</span>
              </div>
              <p class="text-sm text-gray-600">@{{ user.username }}</p>
            </div>
            <span :class="user.active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'" class="px-3 py-1 rounded-full text-xs font-medium">
              {{ user.active ? '✅ Hoạt động' : '🔒 Khóa' }}
            </span>
          </div>
          <div class="flex gap-2 pt-3 border-t">
            <button @click="showEditForm(user)" class="flex-1 bg-blue-500 text-white py-2 rounded-lg text-sm font-medium">✏️ Sửa</button>
            <button @click="showResetPasswordForm(user)" class="flex-1 bg-yellow-500 text-white py-2 rounded-lg text-sm font-medium">🔑 Reset</button>
            <button @click="toggleUserStatus(user)" :class="user.active ? 'bg-orange-500' : 'bg-green-500'" class="flex-1 text-white py-2 rounded-lg text-sm font-medium">{{ user.active ? '⏸️' : '▶️' }}</button>
            <button @click="showDeleteConfirm(user)" class="flex-1 bg-red-500 text-white py-2 rounded-lg text-sm font-medium">🗑️</button>
          </div>
        </div>
      </div>
    </div>

    <BottomNav />

    <transition name="slide-right">
      <div v-if="showCreateForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white w-full h-screen flex flex-col">
          <div class="px-4 py-3 flex items-center justify-between border-b">
            <button @click="showCreateForm = false" class="text-2xl">←</button>
            <h1 class="text-xl font-bold">➕ Tạo User</h1>
            <div class="w-8"></div>
          </div>
          <div class="flex-1 overflow-y-auto px-4 py-6 space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Username *</label>
              <input v-model="createForm.username" type="text" required class="w-full px-4 py-3 border rounded-lg" placeholder="Nhập username">
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Mật khẩu *</label>
              <input v-model="createForm.password" type="password" required class="w-full px-4 py-3 border rounded-lg" placeholder="Nhập mật khẩu">
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Tên hiển thị *</label>
              <input v-model="createForm.name" type="text" required class="w-full px-4 py-3 border rounded-lg" placeholder="Nhập tên">
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Vai trò *</label>
              <select v-model="createForm.role" required class="w-full px-4 py-3 border rounded-lg">
                <option value="">Chọn vai trò</option>
                <option value="manager">Manager</option>
                <option value="cashier">Cashier</option>
                <option value="waiter">Waiter</option>
              </select>
            </div>
            <div class="h-20"></div>
          </div>
          <div class="flex gap-3 px-4 py-4 border-t">
            <button @click="showCreateForm = false" class="flex-1 bg-gray-200 py-3 rounded-lg">Hủy</button>
            <button @click="createUser" class="flex-1 bg-blue-500 text-white py-3 rounded-lg">Tạo</button>
          </div>
        </div>
      </div>
    </transition>

    <transition name="slide-right">
      <div v-if="showEditModal" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white w-full h-screen flex flex-col">
          <div class="px-4 py-3 flex items-center justify-between border-b">
            <button @click="showEditModal = false" class="text-2xl">←</button>
            <h1 class="text-xl font-bold">✏️ Sửa User</h1>
            <div class="w-8"></div>
          </div>
          <div class="flex-1 overflow-y-auto px-4 py-6 space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">Tên hiển thị *</label>
              <input v-model="editForm.name" type="text" required class="w-full px-4 py-3 border rounded-lg" placeholder="Nhập tên">
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Vai trò *</label>
              <select v-model="editForm.role" required class="w-full px-4 py-3 border rounded-lg">
                <option value="manager">Manager</option>
                <option value="cashier">Cashier</option>
                <option value="waiter">Waiter</option>
              </select>
            </div>
            <div class="h-20"></div>
          </div>
          <div class="flex gap-3 px-4 py-4 border-t">
            <button @click="showEditModal = false" class="flex-1 bg-gray-200 py-3 rounded-lg">Hủy</button>
            <button @click="updateUser" class="flex-1 bg-blue-500 text-white py-3 rounded-lg">Cập nhật</button>
          </div>
        </div>
      </div>
    </transition>

    <transition name="slide-right">
      <div v-if="showResetPassword" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white w-full h-screen flex flex-col">
          <div class="px-4 py-3 flex items-center justify-between border-b">
            <button @click="showResetPassword = false" class="text-2xl">←</button>
            <h1 class="text-xl font-bold">🔑 Reset Password</h1>
            <div class="w-8"></div>
          </div>
          <div class="flex-1 overflow-y-auto px-4 py-6 space-y-4">
            <div class="bg-gray-50 p-4 rounded-lg">
              <p class="text-sm text-gray-600">Reset cho: <strong>{{ selectedUser?.name }}</strong></p>
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Mật khẩu mới *</label>
              <input v-model="resetPasswordForm.newPassword" type="password" required class="w-full px-4 py-3 border rounded-lg" placeholder="Nhập mật khẩu mới">
            </div>
            <div class="h-20"></div>
          </div>
          <div class="flex gap-3 px-4 py-4 border-t">
            <button @click="showResetPassword = false" class="flex-1 bg-gray-200 py-3 rounded-lg">Hủy</button>
            <button @click="resetPassword" class="flex-1 bg-yellow-500 text-white py-3 rounded-lg">Reset</button>
          </div>
        </div>
      </div>
    </transition>

    <transition name="fade">
      <div v-if="showDeleteConfirmation" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
        <div class="bg-white rounded-2xl p-6 w-full max-w-sm mx-4">
          <h3 class="text-xl font-bold mb-4 text-red-600">⚠️ Xác nhận xóa</h3>
          <p class="text-gray-600 mb-4">Xóa user: <strong>{{ selectedUser?.name }}</strong>?</p>
          <div class="flex gap-3">
            <button @click="showDeleteConfirmation = false" class="flex-1 bg-gray-200 py-3 rounded-lg">Hủy</button>
            <button @click="deleteUser" class="flex-1 bg-red-500 text-white py-3 rounded-lg">Xóa</button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '../stores/user'
import BottomNav from '../components/BottomNav.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'

const userStore = useUserStore()
const searchQuery = ref('')
const showCreateForm = ref(false)
const showEditModal = ref(false)
const showResetPassword = ref(false)
const showDeleteConfirmation = ref(false)
const selectedUser = ref(null)
const currentFilter = ref('ALL')

const createForm = ref({ username: '', password: '', name: '', role: '', active: true })
const editForm = ref({ id: '', username: '', name: '', role: '', active: true })
const resetPasswordForm = ref({ newPassword: '' })

const loading = computed(() => userStore.loading)
const users = computed(() => userStore.users || [])

const filteredUsers = computed(() => {
  let filtered = users.value
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    filtered = filtered.filter(u => u.name?.toLowerCase().includes(q) || u.username?.toLowerCase().includes(q))
  }
  return filtered
})

const activeCount = computed(() => users.value.filter(u => u.active).length)
const managerCount = computed(() => users.value.filter(u => u.role === 'manager').length)
const cashierCount = computed(() => users.value.filter(u => u.role === 'cashier').length)

const getRoleColor = (role) => {
  const colors = { manager: 'bg-purple-100 text-purple-800', cashier: 'bg-blue-100 text-blue-800', waiter: 'bg-green-100 text-green-800' }
  return colors[role] || 'bg-gray-100 text-gray-800'
}

const getRoleText = (role) => {
  const texts = { manager: 'Manager', cashier: 'Cashier', waiter: 'Waiter' }
  return texts[role] || role
}

const formatDate = (date) => new Date(date).toLocaleString('vi-VN')

const openCreateModal = () => {
  createForm.value = { username: '', password: '', name: '', role: '', active: true }
  showCreateForm.value = true
}

const createUser = async () => {
  try {
    await userStore.createUser(createForm.value)
    showCreateForm.value = false
    alert('✅ Tạo user thành công!')
    await userStore.fetchUsers()
  } catch (error) {
    alert('❌ Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

const showEditForm = (user) => {
  selectedUser.value = user
  editForm.value = { id: user.id, username: user.username, name: user.name, role: user.role, active: user.active }
  showEditModal.value = true
}

const updateUser = async () => {
  try {
    await userStore.updateUser(editForm.value.id, { name: editForm.value.name, role: editForm.value.role, active: editForm.value.active })
    showEditModal.value = false
    alert('✅ Cập nhật user thành công!')
    await userStore.fetchUsers()
  } catch (error) {
    alert('❌ Lỗi: ' + error.message)
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
    alert('✅ Reset password thành công!')
    await userStore.fetchUsers()
  } catch (error) {
    alert('❌ Lỗi: ' + error.message)
  }
}

const toggleUserStatus = async (user) => {
  try {
    await userStore.toggleUserStatus(user.id)
    alert(`✅ ${user.active ? 'Khóa' : 'Mở khóa'} user thành công!`)
    await userStore.fetchUsers()
  } catch (error) {
    alert('❌ Lỗi: ' + error.message)
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
    alert('✅ Xóa user thành công!')
    await userStore.fetchUsers()
  } catch (error) {
    alert('❌ Lỗi: ' + error.message)
  }
}

// Refresh data function
const refreshData = async () => {
  await userStore.fetchUsers()
}

// Pull to refresh
const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)

onMounted(async () => {
  await refreshData()
})
</script>

<style scoped>
.slide-right-enter-active, .slide-right-leave-active { transition: transform 0.3s ease; }
.slide-right-enter-from { transform: translateX(100%); }
.slide-right-leave-to { transform: translateX(100%); }
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
