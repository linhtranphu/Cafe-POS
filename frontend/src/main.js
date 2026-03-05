import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'
import { useShopSettingsStore } from './stores/shopSettings'
// WebSocket disabled - using HTTP only for print bridge
// import { websocketService } from './services/websocket'
import { localPrintService } from './services/localPrint'
import './style.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// Khôi phục auth từ localStorage khi app load
const authStore = useAuthStore()
authStore.initAuth()

// Initialize shop settings and print bridge URL
const shopSettingsStore = useShopSettingsStore()
if (authStore.isAuthenticated) {
  // Load shop settings to get print bridge URL
  shopSettingsStore.fetchSettings()
    .then(() => {
      const settings = shopSettingsStore.settings
      if (settings?.print_bridge_url) {
        localPrintService.updateBridgeUrl(settings.print_bridge_url)
      }
    })
    .catch(err => {
      console.warn('[App] Failed to load shop settings:', err)
    })
}

// Validate token with backend on app load
if (authStore.isAuthenticated) {
  authStore.validateToken().catch(() => {
    // Token invalid, logout and show message
    authStore.logout()
    alert('Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.')
    router.push('/login')
  })

  // WebSocket disabled - using HTTP only for print bridge
  // const token = localStorage.getItem('token')
  // if (token) {
  //   websocketService.connect(token)
  // }
}

// WebSocket disabled - using HTTP only for print bridge
// authStore.$subscribe((mutation, state) => {
//   if (state.isAuthenticated && state.token) {
//     websocketService.connect(state.token)
//   } else {
//     websocketService.disconnect()
//   }
// })

app.mount('#app')