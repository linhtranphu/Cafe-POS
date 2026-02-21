import api from './api'

export const printTemplateService = {
  async fetchTemplates(params = {}) {
    try {
      const response = await api.get('/print-templates', { params })
      const data = response.data
      if (Array.isArray(data)) {
        return data
      } else if (data && Array.isArray(data.data)) {
        return data.data
      } else if (data && Array.isArray(data.templates)) {
        return data.templates
      }
      return []
    } catch (error) {
      console.error('Error fetching print templates:', error)
      return []
    }
  },

  async fetchTemplateById(id) {
    const response = await api.get(`/print-templates/${id}`)
    return response.data
  },

  async createTemplate(template) {
    const response = await api.post('/print-templates', template)
    return response.data
  },

  async updateTemplate(id, template) {
    const response = await api.put(`/print-templates/${id}`, template)
    return response.data
  },

  async deleteTemplate(id) {
    await api.delete(`/print-templates/${id}`)
  },

  async previewTemplate(id, sampleData = {}) {
    const response = await api.post(`/print-templates/${id}/preview`, sampleData)
    return response.data
  }
}
