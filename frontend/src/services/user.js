import api from './api'

export const userService = {
  // User Management (Manager only)
  async createUser(userData) {
    const response = await api.post('/manager/users', userData)
    return response.data
  },

  async getAllUsers() {
    const response = await api.get('/manager/users')
    return response.data
  },

  async getUser(id) {
    const response = await api.get(`/manager/users/${id}`)
    return response.data
  },

  async updateUser(id, userData) {
    const response = await api.put(`/manager/users/${id}`, userData)
    return response.data
  },

  async resetPassword(id, newPassword) {
    const response = await api.post(`/manager/users/${id}/reset-password`, {
      new_password: newPassword
    })
    return response.data
  },

  async toggleUserStatus(id) {
    const response = await api.post(`/manager/users/${id}/toggle-status`)
    return response.data
  },

  async deleteUser(id) {
    const response = await api.delete(`/manager/users/${id}`)
    return response.data
  },

  async getUsersByRole(role) {
    const response = await api.get(`/manager/users/by-role?role=${role}`)
    return response.data
  },

  async getActiveUsers() {
    const response = await api.get('/manager/users/active')
    return response.data
  },

  // Profile Management (All users)
  async getCurrentUser() {
    const response = await api.get('/profile')
    return response.data
  },

  async changePassword(currentPassword, newPassword) {
    const response = await api.post('/change-password', {
      current_password: currentPassword,
      new_password: newPassword
    })
    return response.data
  }
}