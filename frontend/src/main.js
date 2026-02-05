import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'
import './style.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// Khôi phục auth từ localStorage khi app load
const authStore = useAuthStore()
authStore.initAuth()

// Validate token with backend on app load
if (authStore.isAuthenticated) {
  authStore.validateToken().catch(() => {
    // Token invalid, logout and show message
    authStore.logout()
    alert('Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.')
    router.push('/login')
  })
}

app.mount('#app')