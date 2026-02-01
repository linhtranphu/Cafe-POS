import { defineStore } from 'pinia'
import { login as authLogin } from '../services/auth'
import api from '../services/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: null,
    isAuthenticated: false,
    loading: false,
    error: null
  }),

  getters: {
    isManager: (state) => state.user?.role === 'manager',
    isWaiter: (state) => state.user?.role === 'waiter',
    isCashier: (state) => state.user?.role === 'cashier',
    isBarista: (state) => state.user?.role === 'barista'
  },

  actions: {
    async login(credentials) {
      this.loading = true
      this.error = null
      try {
        const response = await authLogin(credentials)
        
        if (response && response.user && response.token) {
          this.user = response.user
          this.token = response.token
          this.isAuthenticated = true
          
          // Lưu vào localStorage
          localStorage.setItem('token', response.token)
          localStorage.setItem('user', JSON.stringify(response.user))
          
          // Set token cho API requests
          api.defaults.headers.common['Authorization'] = `Bearer ${response.token}`
          
          return true
        } else {
          throw new Error('Invalid response format')
        }
      } catch (error) {
        console.error('Login error:', error)
        this.error = error.response?.data?.error || error.message || 'Đăng nhập thất bại'
        return false
      } finally {
        this.loading = false
      }
    },

    setUser(user) {
      this.user = user
      this.isAuthenticated = !!user
    },

    setToken(token) {
      this.token = token
    },

    logout() {
      this.user = null
      this.token = null
      this.isAuthenticated = false
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      delete api.defaults.headers.common['Authorization']
    },

    // Khôi phục auth từ localStorage khi app load
    initAuth() {
      const token = localStorage.getItem('token')
      const user = localStorage.getItem('user')
      
      if (token && user) {
        try {
          this.token = token
          this.user = JSON.parse(user)
          this.isAuthenticated = true
          
          // Set token cho API requests
          api.defaults.headers.common['Authorization'] = `Bearer ${token}`
        } catch (error) {
          console.error('Error restoring auth:', error)
          this.logout()
        }
      }
    },

    // Validate token với backend (optional)
    async validateToken() {
      if (!this.token) return false
      
      try {
        const response = await api.get('/profile')
        if (response.data) {
          this.user = response.data
          localStorage.setItem('user', JSON.stringify(response.data))
          return true
        }
        return false
      } catch (error) {
        console.error('Token validation failed:', error)
        this.logout()
        return false
      }
    }
  }
})
