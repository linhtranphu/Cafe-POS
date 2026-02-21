import api from './api'

export const printerConfigService = {
  async fetchPrinters(params = {}) {
    try {
      const response = await api.get('/manager/printers', { params })
      const data = response.data
      if (Array.isArray(data)) {
        return data
      } else if (data && Array.isArray(data.data)) {
        return data.data
      } else if (data && Array.isArray(data.printers)) {
        return data.printers
      }
      return []
    } catch (error) {
      console.error('Error fetching printers:', error)
      return []
    }
  },

  async fetchPrinterById(id) {
    const response = await api.get(`/manager/printers/${id}`)
    return response.data
  },

  async createPrinter(printer) {
    const response = await api.post('/manager/printers', printer)
    return response.data
  },

  async updatePrinter(id, printer) {
    const response = await api.put(`/manager/printers/${id}`, printer)
    return response.data
  },

  async deletePrinter(id) {
    await api.delete(`/manager/printers/${id}`)
  },

  async testConnection(id) {
    const response = await api.post(`/manager/printers/${id}/test`)
    return response.data
  }
}
