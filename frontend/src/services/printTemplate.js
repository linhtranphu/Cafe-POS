import api from './api'

export const printTemplateService = {
  async fetchTemplates(params = {}) {
    try {
      const response = await api.get('/manager/print-templates', { params })
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
    const response = await api.get(`/manager/print-templates/${id}`)
    return response.data
  },

  async createTemplate(template) {
    const response = await api.post('/manager/print-templates', template)
    return response.data
  },

  async updateTemplate(id, template) {
    const response = await api.put(`/manager/print-templates/${id}`, template)
    return response.data
  },

  async deleteTemplate(id) {
    await api.delete(`/manager/print-templates/${id}`)
  },

  async previewTemplate(id, sampleData = {}) {
    // Fetch the template first to get its content and type
    const templateResponse = await api.get(`/manager/print-templates/${id}`)
    const template = templateResponse.data.template || templateResponse.data
    
    // Send content and type in the request body as required by backend
    const previewRequest = {
      content: template.content,
      type: template.type,
      ...sampleData
    }
    
    const response = await api.post(`/manager/print-templates/${id}/preview`, previewRequest)
    return response.data
  }
}
