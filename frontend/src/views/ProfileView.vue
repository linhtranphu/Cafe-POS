<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm">
      <div class="px-4 py-3">
        <h1 class="text-xl font-bold text-gray-800">👤 Cá nhân</h1>
      </div>
    </div>

    <!-- Content -->
    <div class="px-4 py-4 pb-24">
      <div v-if="loading" class="text-center py-10">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
      </div>

      <div v-else-if="currentUser" class="space-y-4">
        <!-- Profile Card -->
        <div class="bg-gradient-to-br from-blue-500 to-purple-500 rounded-2xl p-6 text-white shadow-lg">
          <div class="text-center mb-4">
            <div class="w-20 h-20 bg-white rounded-full mx-auto mb-3 flex items-center justify-center text-4xl">
              👤
            </div>
            <h2 class="text-2xl font-bold">{{ currentUser.name }}</h2>
            <p class="text-blue-100">@{{ currentUser.username }}</p>
          </div>
          
          <div class="flex justify-center gap-2">
            <span :class="getRoleColor(currentUser.role)" 
              class="px-4 py-2 rounded-full text-sm font-medium">
              {{ getRoleText(currentUser.role) }}
            </span>
            <span :class="currentUser.active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'" 
              class="px-4 py-2 rounded-full text-sm font-medium">
              {{ currentUser.active ? 'Hoạt động' : 'Tạm khóa' }}
            </span>
          </div>
        </div>

        <!-- Info Card -->
        <div class="bg-white rounded-2xl p-6 shadow-sm">
          <h3 class="text-lg font-bold mb-4">Thông tin tài khoản</h3>
          
          <div class="space-y-3">
            <div class="flex items-center justify-between p-3 bg-gray-50 rounded-xl">
              <div>
                <p class="text-xs text-gray-500">Username</p>
                <p class="font-medium">{{ currentUser.username }}</p>
              </div>
              <span class="text-2xl">👤</span>
            </div>
            
            <div class="flex items-center justify-between p-3 bg-gray-50 rounded-xl">
              <div>
                <p class="text-xs text-gray-500">Tên hiển thị</p>
                <p class="font-medium">{{ currentUser.name }}</p>
              </div>
              <span class="text-2xl">✏️</span>
            </div>
            
            <div class="flex items-center justify-between p-3 bg-gray-50 rounded-xl">
              <div>
                <p class="text-xs text-gray-500">Ngày tạo</p>
                <p class="font-medium">{{ formatDate(currentUser.created_at) }}</p>
              </div>
              <span class="text-2xl">📅</span>
            </div>
            
            <div class="flex items-center justify-between p-3 bg-gray-50 rounded-xl">
              <div>
                <p class="text-xs text-gray-500">Đăng nhập cuối</p>
                <p class="font-medium">
                  {{ currentUser.last_login ? formatDate(currentUser.last_login) : 'Chưa có' }}
                </p>
              </div>
              <span class="text-2xl">🕐</span>
            </div>
          </div>
        </div>

        <!-- Stats Card -->
        <div class="bg-white rounded-2xl p-6 shadow-sm">
          <h3 class="text-lg font-bold mb-4">Thống kê hoạt động</h3>
          <div class="grid grid-cols-3 gap-3">
            <div class="text-center p-4 bg-blue-50 rounded-xl">
              <div class="text-2xl mb-1">📋</div>
              <div class="text-xl font-bold text-blue-600">--</div>
              <div class="text-xs text-gray-600">Orders</div>
            </div>
            <div class="text-center p-4 bg-green-50 rounded-xl">
              <div class="text-2xl mb-1">⏰</div>
              <div class="text-xl font-bold text-green-600">--</div>
              <div class="text-xs text-gray-600">Ca làm</div>
            </div>
            <div class="text-center p-4 bg-purple-50 rounded-xl">
              <div class="text-2xl mb-1">💰</div>
              <div class="text-xl font-bold text-purple-600">--</div>
              <div class="text-xs text-gray-600">Doanh thu</div>
            </div>
          </div>
          <p class="text-xs text-gray-500 text-center mt-3">
            * Thống kê chi tiết sẽ được cập nhật
          </p>
        </div>

        <!-- Change Password Card -->
        <div class="bg-white rounded-2xl p-6 shadow-sm">
          <h3 class="text-lg font-bold mb-4">🔒 Đổi mật khẩu</h3>
          
          <button @click="showPasswordForm = !showPasswordForm"
            class="w-full bg-blue-500 text-white py-3 rounded-xl font-medium active:scale-95 transition-transform">
            {{ showPasswordForm ? 'Ẩn form' : 'Đổi mật khẩu' }}
          </button>

          <form v-if="showPasswordForm" @submit.prevent="changePassword" class="space-y-4 mt-4">
            <div>
              <label class="block text-sm font-medium mb-2">Mật khẩu hiện tại *</label>
              <input v-model="passwordForm.currentPassword" type="password" required
                class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Nhập mật khẩu hiện tại">
            </div>
            
            <div>
              <label class="block text-sm font-medium mb-2">Mật khẩu mới *</label>
              <input v-model="passwordForm.newPassword" type="password" required minlength="6"
                class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Tối thiểu 6 ký tự">
            </div>
            
            <div>
              <label class="block text-sm font-medium mb-2">Xác nhận mật khẩu *</label>
              <input v-model="passwordForm.confirmPassword" type="password" required
                class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Nhập lại mật khẩu mới">
              <div v-if="passwordForm.newPassword && passwordForm.confirmPassword && passwordForm.newPassword !== passwordForm.confirmPassword" 
                class="text-red-500 text-sm mt-1">
                ⚠️ Mật khẩu xác nhận không khớp
              </div>
            </div>
            
            <div class="flex gap-2">
              <button type="button" @click="resetPasswordForm" 
                class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
                Hủy
              </button>
              <button type="submit" :disabled="!isPasswordFormValid" 
                class="flex-1 bg-green-500 hover:bg-green-600 text-white px-4 py-3 rounded-xl font-medium disabled:opacity-50 disabled:cursor-not-allowed">
                Xác nhận
              </button>
            </div>
          </form>
        </div>

        <!-- Logout Button -->
        <button @click="logout" 
          class="w-full bg-red-500 hover:bg-red-600 text-white py-3 rounded-xl font-bold active:scale-95 transition-transform">
          🚪 Đăng xuất
        </button>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '../stores/user'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import BottomNav from '../components/BottomNav.vue'

const userStore = useUserStore()
const authStore = useAuthStore()
const router = useRouter()

const showPasswordForm = ref(false)

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
    alert('✅ Đổi mật khẩu thành công!')
    resetPasswordForm()
    showPasswordForm.value = false
  } catch (error) {
    alert('❌ Lỗi: ' + error.message)
  }
}

const resetPasswordForm = () => {
  passwordForm.value = {
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  }
}

const logout = async () => {
  if (confirm('Bạn có chắc muốn đăng xuất?')) {
    await authStore.logout()
    router.push('/login')
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
  return new Date(date).toLocaleDateString('vi-VN', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<style scoped>
.active\:scale-95:active {
  transform: scale(0.95);
}

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
