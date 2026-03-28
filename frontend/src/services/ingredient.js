import api from './api'

const getAuthHeader = () => {
  const token = localStorage.getItem('token')
  return { Authorization: `Bearer ${token}` }
}

export const ingredientService = {
  async getIngredients() {
    const response = await api.get('/manager/ingredients')
    return response.data
  },

  async createIngredient(ingredient) {
    const response = await api.post('/manager/ingredients', ingredient)
    return response.data
  },

  async updateIngredient(id, ingredient) {
    const response = await api.put(`/manager/ingredients/${id}`, ingredient)
    return response.data
  },

  async deleteIngredient(id) {
    await api.delete(`/manager/ingredients/${id}`)
  },

  async getLowStockItems() {
    const response = await api.get('/manager/ingredients/low-stock')
    return response.data
  },

  async getStockHistory(id) {
    const response = await api.get(`/manager/ingredients/${id}/history`)
    return response.data
  },

  // New stock operation methods
  async stockIn(id, data) {
    // data: { quantity, cost_per_unit, reason }
    const response = await api.post(`/manager/ingredients/${id}/stock-in`, data)
    return response.data
  },

  async stockOut(id, data) {
    // data: { quantity, reason }
    const response = await api.post(`/manager/ingredients/${id}/stock-out`, data)
    return response.data
  },

  async stockAdjust(id, data) {
    // data: { new_quantity, cost_per_unit, reason }
    const response = await api.post(`/manager/ingredients/${id}/stock-adjust`, data)
    return response.data
  },

  // Legacy method for backward compatibility
  async adjustStock(id, adjustment) {
    const response = await api.post(`/manager/ingredients/${id}/adjust`, adjustment)
    return response.data
  },

  // Category methods
  async getCategories() {
    const response = await api.get('/manager/ingredient-categories')
    return response.data
  },

  async createCategory(category) {
    const response = await api.post('/manager/ingredient-categories', category)
    return response.data
  },

  async deleteCategory(id) {
    await api.delete(`/manager/ingredient-categories/${id}`)
  },

  async getStockSummary(from, to) {
    const params = {}
    if (from) params.from = from
    if (to) params.to = to
    const response = await api.get('/manager/ingredients/summary', { params })
    return response.data
  }
}
