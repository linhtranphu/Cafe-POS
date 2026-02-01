<template>
  <div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">
    <!-- Mobile Header - Fixed -->
    <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
      <div class="px-4 py-3">
        <h1 class="text-xl font-bold text-gray-800">Profile</h1>
      </div>
    </div>

    <!-- Content - Scrollable -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-24">
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
            <span :class="getRoleColor(currentUser.role)" class="px-4 py-2 rounded-full text-sm font-medium">
              {{ getRoleText(currentUser.role) }}
            </span>
            <span :class="currentUser.active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'" class="px-4 py-2 rounded-full text-sm font-medium">
              {{ currentUser.active ? 'Active' : 'Locked' }}
            </span>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="mb-4">
          <h2 class="text-sm font-bold text-gray-800 mb-2">Quick Actions</h2>
          <div class="grid grid-cols-2 gap-2">
            <button @click="showPasswordForm = true" class="bg-gradient-to-br from-blue-500 to-cyan-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform">
              <div class="text-2xl mb-1">🔒</div>
              <div class="text-sm font-bold">Change Password</div>
            </button>
            <button @click="logout" class="bg-gradient-to-br from-red-500 to-pink-500 text-white rounded-xl p-4 shadow-md active:scale-95 transition-transform">
              <div class="text-2xl mb-1">🚪</div>
              <div class="text-sm font-bold">Logout</div>
            </button>
          </div>
        </div>

        <!-- Info Card -->
        <div class="bg-white rounded-2xl p-6 shadow-sm">
          <h3 class="text-lg font-bold mb-4">Account Info</h3>
          
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
                <p class="text-xs text-gray-500">Display Name</p>
                <p class="font-medium">{{ currentUser.name }}</p>
              </div>
              <span class="text-2xl">✏️</span>
            </div>
            
            <div class="flex items-center justify-between p-3 bg-gray-50 rounded-xl">
              <div>
                <p class="text-xs text-gray-500">Role</p>
                <p class="font-medium">{{ getRoleText(currentUser.role) }}</p>
              </div>
              <span class="text-2xl">👔</span>
            </div>
            
            <div class="flex items-center justify-between p-3 bg-gray-50 rounded-xl">
              <div>
                <p class="text-xs text-gray-500">Created</p>
                <p class="font-medium">{{ formatDate(currentUser.created_at) }}</p>
              </div>
              <span class="text-2xl">📅</span>
            </div>
            
            <div class="flex items-center justify-between p-3 bg-gray-50 rounded-xl">
              <div>
                <p class="text-xs text-gray-500">Last Login</p>
                <p class="font-medium">
                  {{ currentUser.last_login ? formatDate(currentUser.last_login) : 'Never' }}
                </p>
              </div>
              <span class="text-2xl">⏰</span>
            </div>
          </div>
        </div>

        <!-- Stats Card -->
        <div class="bg-white rounded-2xl p-6 shadow-sm">
          <h3 class="text-lg font-bold mb-4">Activity Stats</h3>
          <div class="grid grid-cols-3 gap-3">
            <div class="text-center p-4 bg-blue-50 rounded-xl">
              <div class="text-2xl mb-1">📋</div>
              <div class="text-xl font-bold text-blue-600">--</div>
              <div class="text-xs text-gray-600">Orders</div>
            </div>
            <div class="text-center p-4 bg-green-50 rounded-xl">
              <div class="text-2xl mb-1">🕐</div>
              <div class="text-xl font-bold text-green-600">--</div>
              <div class="text-xs text-gray-600">Shifts</div>
            </div>
            <div class="text-center p-4 bg-purple-50 rounded-xl">
              <div class="text-2xl mb-1">💰</div>
              <div class="text-xl font-bold text-purple-600">--</div>
              <div class="text-xs text-gray-600">Revenue</div>
            </div>
          </div>
          <p class="text-xs text-gray-500 text-center mt-3">* Detailed stats coming soon</p>
        </div>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <BottomNav />

    <!-- Change Password Modal - Slide from Right -->
    <transition name="slide-right">
      <div v-if="showPasswordForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-gray-50 w-full h-screen flex flex-col">
          <!-- Mobile Header - Fixed -->
          <div class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">
            <div class="px-4 py-3 flex items-center justify-between">
              <button @click="showPasswordForm = false" class="text-2xl text-gray-600">Back</button>
              <h1 class="text-xl font-bold text-gray-800">Change Password</h1>
              <div class="w-8"></div>
            </div>
          </div>

          <!-- Scrollable Content -->
          <div class="flex-1 overflow-y-auto px-4 py-6 space-y-5">
            <!-- Current Password -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">Current Password *</label>
              <input v-model="passwordForm.currentPassword" type="password" required
                class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Enter current password">
            </div>

            <!-- New Password -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">New Password *</label>
              <input v-model="passwordForm.newPassword" type="password" required minlength="6"
                class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Enter new password (min 6 characters)">
            </div>

            <!-- Confirm Password -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">Confirm Password *</label>
              <input v-model="passwordForm.confirmPassword" type="password" required
                class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Re-enter new password">
              <div v-if="passwordForm.newPassword && passwordForm.confirmPassword && passwordForm.newPassword !== passwordForm.confirmPassword" 
                class="text-red-500 text-sm mt-2">
                Passwords do not match
              </div>
            </div>

            <!-- Info Box -->
            <div class="bg-blue-50 border border-blue-200 rounded-xl p-4">
              <p class="text-sm text-blue-800">
                Password must be at least 6 characters
              </p>
            </div>

            <!-- Spacer -->
            <div class="h-24"></div>
          </div>

          <!-- Fixed Footer -->
          <div class="flex-shrink-0 bg-white px-4 py-4 border-t flex gap-3 pb-safe">
            <button @click="showPasswordForm = false"
              class="flex-1 bg-gray-200 text-gray-700 py-4 rounded-xl font-medium text-base active:bg-gray-300 transition-colors">
              Cancel
            </button>
            <button @click="changePassword" :disabled="!isPasswordFormValid"
              class="flex-1 bg-green-500 text-white py-4 rounded-xl font-medium text-base active:bg-green-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
              Confirm
            </button>
          </div>
        </div>
      </div>
    </transition>
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
    alert('Please check password information')
    return
  }

  try {
    await userStore.changePassword(passwordForm.value.currentPassword, passwordForm.value.newPassword)
    alert('Password changed successfully!')
    resetPasswordForm()
    showPasswordForm.value = false
  } catch (error) {
    alert('Error: ' + error.message)
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
  if (confirm('Are you sure you want to logout?')) {
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
  return new Date(date).toLocaleDateString('en-US', {
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
</style>
