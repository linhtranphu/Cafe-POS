import api from './api'

export const menuService = {
  async getMenuItems() {
    const response = await api.get('/waiter/menu')
    return response.data
  },

  async createMenuItem(item) {
    const response = await api.post('/manager/menu', item)
    return response.data
  },

  async updateMenuItem(id, item) {
    const response = await api.put(`/manager/menu/${id}`, item)
    return response.data
  },

  async deleteMenuItem(id) {
    const response = await api.delete(`/manager/menu/${id}`)
    return response.data
  }
}