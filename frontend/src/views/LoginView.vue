<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-600 to-purple-700 p-5">
    <div class="bg-white p-8 lg:p-10 rounded-xl shadow-2xl w-full max-w-md">
      <div class="text-center mb-8">
        <h1 class="text-3xl lg:text-4xl font-bold text-gray-800 mb-2">☕ Café POS</h1>
        <p class="text-gray-600 text-sm">Hệ thống quản lý quán cà phê</p>
      </div>

      <form @submit.prevent="handleLogin" class="space-y-5">
        <div>
          <label class="block mb-2 text-gray-700 font-medium">Tên đăng nhập</label>
          <input 
            v-model="username" 
            type="text" 
            required 
            placeholder="Nhập tên đăng nhập"
            :disabled="loading"
            class="w-full p-3 border-2 border-gray-300 rounded-lg text-base transition-colors focus:outline-none focus:border-blue-600 disabled:bg-gray-100 disabled:cursor-not-allowed"
          />
        </div>

        <div>
          <label class="block mb-2 text-gray-700 font-medium">Mật khẩu</label>
          <input 
            v-model="password" 
            type="password" 
            required 
            placeholder="Nhập mật khẩu"
            :disabled="loading"
            class="w-full p-3 border-2 border-gray-300 rounded-lg text-base transition-colors focus:outline-none focus:border-blue-600 disabled:bg-gray-100 disabled:cursor-not-allowed"
          />
        </div>

        <button 
          type="submit" 
          :disabled="loading" 
          class="w-full p-4 bg-blue-600 text-white border-none rounded-lg text-base font-semibold cursor-pointer transition-colors hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
        >
          {{ loading ? 'Đang đăng nhập...' : 'Đăng nhập' }}
        </button>

        <div v-if="error" class="text-red-600 text-center mt-4 p-3 bg-red-50 rounded-lg border border-red-200">
          {{ error }}
        </div>
      </form>

      <div class="mt-8 pt-5 border-t border-gray-200">
        <h3 class="text-gray-600 text-sm mb-3">Tài khoản demo:</h3>
        <div 
          class="p-2 bg-gray-100 my-1 rounded-lg cursor-pointer text-xs transition-colors hover:bg-gray-200" 
          @click="quickLogin('admin', 'admin123')"
        >
          <strong class="text-gray-700">Manager:</strong> admin / admin123
        </div>
        <div 
          class="p-2 bg-gray-100 my-1 rounded-lg cursor-pointer text-xs transition-colors hover:bg-gray-200" 
          @click="quickLogin('waiter1', 'waiter123')"
        >
          <strong class="text-gray-700">Waiter:</strong> waiter1 / waiter123
        </div>
        <div 
          class="p-2 bg-blue-50 my-1 rounded-lg cursor-pointer text-xs transition-colors hover:bg-blue-100 border border-blue-200" 
          @click="quickLogin('barista1', 'barista123')"
        >
          <strong class="text-blue-700">🍹 Barista:</strong> barista1 / barista123
        </div>
        <div 
          class="p-2 bg-yellow-50 my-1 rounded-lg cursor-pointer text-xs transition-colors hover:bg-yellow-100 border border-yellow-200" 
          @click="quickLogin('cashier1', 'cashier123')"
        >
          <strong class="text-yellow-700">💵 Cashier:</strong> cashier1 / cashier123
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  loading.value = true
  error.value = ''
  
  const success = await authStore.login({
    username: username.value,
    password: password.value
  })
  
  if (success) {
    router.push('/dashboard')
  } else {
    error.value = authStore.error
  }
  
  loading.value = false
}

const quickLogin = (user, pass) => {
  username.value = user
  password.value = pass
  handleLogin()
}
</script>

<style scoped>
/* Tailwind handles all styling */
</style>