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
  }
}
