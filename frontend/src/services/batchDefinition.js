import api from './api'

export const batchDefinitionService = {
  async getDefinitions() {
    try {
      const response = await api.get('/manager/batch-definitions')
      // Ensure we always return an array
      const data = response.data
      if (Array.isArray(data)) {
        return data
      } else if (data && Array.isArray(data.data)) {
        return data.data
      } else if (data && Array.isArray(data.definitions)) {
        return data.definitions
      }
      return []
    } catch (error) {
      console.error('Error fetching batch definitions:', error)
      return []
    }
  },

  async getDefinitionById(id) {
    const response = await api.get(`/manager/batch-definitions/${id}`)
    return response.data
  },

  async createDefinition(definition) {
    const response = await api.post('/manager/batch-definitions', definition)
    return response.data
  },

  async updateDefinition(id, definition) {
    const response = await api.put(`/manager/batch-definitions/${id}`, definition)
    return response.data
  },

  async deleteDefinition(id) {
    await api.delete(`/manager/batch-definitions/${id}`)
  }
}
