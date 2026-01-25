<template>
  <div class="login-container">
    <div class="login-card">
      <div class="logo">
        <h1>☕ Café POS</h1>
        <p>Hệ thống quản lý quán cà phê</p>
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label>Tên đăng nhập</label>
          <input 
            v-model="username" 
            type="text" 
            required 
            placeholder="Nhập tên đăng nhập"
            :disabled="loading"
          />
        </div>

        <div class="form-group">
          <label>Mật khẩu</label>
          <input 
            v-model="password" 
            type="password" 
            required 
            placeholder="Nhập mật khẩu"
            :disabled="loading"
          />
        </div>

        <button type="submit" :disabled="loading" class="login-btn">
          {{ loading ? 'Đang đăng nhập...' : 'Đăng nhập' }}
        </button>

        <div v-if="error" class="error">{{ error }}</div>
      </form>

      <div class="demo-accounts">
        <h3>Tài khoản demo:</h3>
        <div class="demo-item" @click="quickLogin('admin', 'admin123')">
          <strong>Manager:</strong> admin / admin123
        </div>
        <div class="demo-item" @click="quickLogin('waiter1', 'waiter123')">
          <strong>Waiter:</strong> waiter1 / waiter123
        </div>
        <div class="demo-item" @click="quickLogin('cashier1', 'cashier123')">
          <strong>Cashier:</strong> cashier1 / cashier123
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
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-card {
  background: white;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
  width: 100%;
  max-width: 400px;
}

.logo {
  text-align: center;
  margin-bottom: 30px;
}

.logo h1 {
  color: #333;
  margin-bottom: 8px;
  font-size: 2.2em;
}

.logo p {
  color: #666;
  font-size: 0.9em;
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

.form-group input {
  width: 100%;
  padding: 12px;
  border: 2px solid #e1e5e9;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.3s;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
}

.form-group input:disabled {
  background: #f5f5f5;
  cursor: not-allowed;
}

.login-btn {
  width: 100%;
  padding: 14px;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.3s;
}

.login-btn:hover:not(:disabled) {
  background: #5a6fd8;
}

.login-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.error {
  color: #e74c3c;
  text-align: center;
  margin-top: 15px;
  padding: 10px;
  background: #fdf2f2;
  border-radius: 6px;
  border: 1px solid #fecaca;
}

.demo-accounts {
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.demo-accounts h3 {
  color: #666;
  font-size: 14px;
  margin-bottom: 10px;
}

.demo-item {
  padding: 8px 12px;
  background: #f8f9fa;
  margin: 5px 0;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  transition: background 0.2s;
}

.demo-item:hover {
  background: #e9ecef;
}

.demo-item strong {
  color: #495057;
}
</style>