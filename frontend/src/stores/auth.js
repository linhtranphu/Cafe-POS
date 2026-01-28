import { defineStore } from 'pinia'
import { login as authLogin } from '../services/auth'

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
          
          localStorage.setItem('token', response.token)
          localStorage.setItem('user', JSON.stringify(response.user))
          
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
    },

    initAuth() {
      const token = localStorage.getItem('token')
      const user = localStorage.getItem('user')
      if (token && user) {
        this.token = token
        this.user = JSON.parse(user)
        this.isAuthenticated = true
      }
    }
  }
})