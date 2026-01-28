import { defineStore } from 'pinia'
import { userService } from '../services/user'

export const useUserStore = defineStore('user', {
  state: () => ({
    users: [],
    currentUser: null,
    loading: false,
    error: null
  }),

  actions: {
    async fetchUsers() {
      this.loading = true
      this.error = null
      try {
        this.users = await userService.getAllUsers() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải danh sách user'
        this.users = []
      } finally {
        this.loading = false
      }
    },

    async createUser(userData) {
      this.error = null
      try {
        const user = await userService.createUser(userData)
        this.users.unshift(user)
        return user
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tạo user'
        throw error
      }
    },

    async updateUser(id, userData) {
      this.error = null
      try {
        const user = await userService.updateUser(id, userData)
        this.updateUserInList(user)
        return user
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi cập nhật user'
        throw error
      }
    },

    async resetPassword(id, newPassword) {
      this.error = null
      try {
        await userService.resetPassword(id, newPassword)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi reset password'
        throw error
      }
    },

    async toggleUserStatus(id) {
      this.error = null
      try {
        const user = await userService.toggleUserStatus(id)
        this.updateUserInList(user)
        return user
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi thay đổi trạng thái user'
        throw error
      }
    },

    async deleteUser(id) {
      this.error = null
      try {
        await userService.deleteUser(id)
        this.users = this.users.filter(u => u.id !== id)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi xóa user'
        throw error
      }
    },

    async fetchUsersByRole(role) {
      this.loading = true
      this.error = null
      try {
        this.users = await userService.getUsersByRole(role) || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải users theo role'
        this.users = []
      } finally {
        this.loading = false
      }
    },

    async fetchActiveUsers() {
      this.loading = true
      this.error = null
      try {
        this.users = await userService.getActiveUsers() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải active users'
        this.users = []
      } finally {
        this.loading = false
      }
    },

    async fetchCurrentUser() {
      this.error = null
      try {
        this.currentUser = await userService.getCurrentUser()
        return this.currentUser
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải thông tin user'
        throw error
      }
    },

    async changePassword(currentPassword, newPassword) {
      this.error = null
      try {
        await userService.changePassword(currentPassword, newPassword)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi đổi mật khẩu'
        throw error
      }
    },

    updateUserInList(user) {
      const index = this.users.findIndex(u => u.id === user.id)
      if (index !== -1) {
        this.users[index] = user
      }
    },

    setCurrentUser(user) {
      this.currentUser = user
    },

    clearError() {
      this.error = null
    }
  },

  getters: {
    usersByRole: (state) => (role) => {
      return state.users.filter(u => u.role === role)
    },

    activeUsers: (state) => {
      return state.users.filter(u => u.active)
    },

    inactiveUsers: (state) => {
      return state.users.filter(u => !u.active)
    },

    managers: (state) => {
      return state.users.filter(u => u.role === 'manager')
    },

    cashiers: (state) => {
      return state.users.filter(u => u.role === 'cashier')
    },

    waiters: (state) => {
      return state.users.filter(u => u.role === 'waiter')
    }
  }
})