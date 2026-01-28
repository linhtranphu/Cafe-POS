<template>
  <div class="min-h-screen bg-gray-100">
    <Navigation />
    <div class="p-4">
      <!-- Header -->
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold text-gray-800">👤 Thông tin cá nhân</h2>
      </div>

      <div class="max-w-2xl mx-auto space-y-6">
        <!-- Profile Info Card -->
        <div class="bg-white rounded-xl p-6 shadow-sm">
          <h3 class="text-lg font-bold mb-4">Thông tin tài khoản</h3>
          
          <div v-if="loading" class="text-center py-4">Đang tải...</div>
          <div v-else-if="currentUser" class="space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Username</label>
                <div class="p-3 bg-gray-50 rounded-lg">{{ currentUser.username }}</div>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Tên hiển thị</label>
                <div class="p-3 bg-gray-50 rounded-lg">{{ currentUser.name }}</div>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Vai trò</label>
                <div class="p-3 bg-gray-50 rounded-lg">
                  <span :class="getRoleColor(currentUser.role)" class="px-2 py-1 rounded-full text-xs font-medium">
                    {{ getRoleText(currentUser.role) }}
                  </span>
                </div>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Trạng thái</label>
                <div class="p-3 bg-gray-50 rounded-lg">
                  <span :class="currentUser.active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'" 
                    class="px-2 py-1 rounded-full text-xs font-medium">
                    {{ currentUser.active ? 'Hoạt động' : 'Tạm khóa' }}
                  </span>
                </div>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Ngày tạo</label>
                <div class="p-3 bg-gray-50 rounded-lg">{{ formatDate(currentUser.created_at) }}</div>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Đăng nhập cuối</label>
                <div class="p-3 bg-gray-50 rounded-lg">
                  {{ currentUser.last_login ? formatDate(currentUser.last_login) : 'Chưa có' }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Change Password Card -->
        <div class="bg-white rounded-xl p-6 shadow-sm">
          <h3 class="text-lg font-bold mb-4">Đổi mật khẩu</h3>
          
          <form @submit.prevent="changePassword" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Mật khẩu hiện tại *</label>
              <input v-model="passwordForm.currentPassword" type="password" required
                class="w-full p-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="Nhập mật khẩu hiện tại">
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Mật khẩu mới *</label>
              <input v-model="passwordForm.newPassword" type="password" required minlength="6"
                class="w-full p-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="Nhập mật khẩu mới (tối thiểu 6 ký tự)">
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Xác nhận mật khẩu mới *</label>
              <input v-model="passwordForm.confirmPassword" type="password" required
                class="w-full p-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="Nhập lại mật khẩu mới">
              <div v-if="passwordForm.newPassword && passwordForm.confirmPassword && passwordForm.newPassword !== passwordForm.confirmPassword" 
                class="text-red-500 text-sm mt-1">
                Mật khẩu xác nhận không khớp
              </div>
            </div>
            
            <div class="flex gap-2">
              <button type="button" @click="resetPasswordForm" class="flex-1 bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg">
                Hủy
              </button>
              <button type="submit" :disabled="!isPasswordFormValid" 
                class="flex-1 bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed">
                Đổi mật khẩu
              </button>
            </div>
          </form>
        </div>

        <!-- Activity Summary Card (if available) -->
        <div class="bg-white rounded-xl p-6 shadow-sm">
          <h3 class="text-lg font-bold mb-4">Thống kê hoạt động</h3>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="text-center p-4 bg-blue-50 rounded-lg">
              <div class="text-2xl font-bold text-blue-600">--</div>
              <div class="text-sm text-gray-600">Orders hôm nay</div>
            </div>
            <div class="text-center p-4 bg-green-50 rounded-lg">
              <div class="text-2xl font-bold text-green-600">--</div>
              <div class="text-sm text-gray-600">Ca làm việc</div>
            </div>
            <div class="text-center p-4 bg-purple-50 rounded-lg">
              <div class="text-2xl font-bold text-purple-600">--</div>
              <div class="text-sm text-gray-600">Doanh thu</div>
            </div>
          </div>
          <p class="text-sm text-gray-500 text-center mt-4">
            * Thống kê chi tiết sẽ được cập nhật trong phiên bản tiếp theo
          </p>
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

const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const loading = computed(() => userStore.loading)
const currentUser = computed(() => userStore.currentUser)

const isPasswordFormValid = computed(() => {
  return passwordForm.value.currentPassword &&
         passwordForm.value.newPassword &&
         passwordForm.value.confirmPassword &&
         passwordForm.value.newPassword === passwordForm.value.confirmPassword &&
         passwordForm.value.newPassword.length >= 6
})

onMounted(async () => {
  try {
    await userStore.fetchCurrentUser()
  } catch (error) {
    console.error('Error loading profile:', error)
  }
})

const changePassword = async () => {
  if (!isPasswordFormValid.value) {
    alert('Vui lòng kiểm tra lại thông tin mật khẩu')
    return
  }

  try {
    await userStore.changePassword(passwordForm.value.currentPassword, passwordForm.value.newPassword)
    alert('Đổi mật khẩu thành công!')
    resetPasswordForm()
  } catch (error) {
    alert('Lỗi: ' + error.message)
  }
}

const resetPasswordForm = () => {
  passwordForm.value = {
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
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