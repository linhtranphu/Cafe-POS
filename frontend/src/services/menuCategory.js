import api from './api'

export const menuCategoryService = {
  async getCategories() {
    const response = await api.get('/manager/menu-categories')
    return response.data
  },

  async createCategory(category) {
    const response = await api.post('/manager/menu-categories', category)
    return response.data
  },

  async updateCategory(id, category) {
    const response = await api.put(`/manager/menu-categories/${id}`, category)
    return response.data
  },

  async deleteCategory(id) {
    const response = await api.delete(`/manager/menu-categories/${id}`)
    return response.data
  },

  async reorderCategories(orders) {
    // orders: [{id, sort_order}, ...]
    const response = await api.put('/manager/menu-categories/reorder', { orders })
    return response.data
  },

  async getWaiterCategories() {
    const response = await api.get('/waiter/menu-categories')
    return response.data
  }
}
