import api from './api'

export const printJobService = {
  async fetchPrintJobs(params = {}) {
    try {
      const response = await api.get('/manager/print-jobs', { params })
      const data = response.data
      if (Array.isArray(data)) {
        return data
      } else if (data && Array.isArray(data.data)) {
        return data.data
      } else if (data && Array.isArray(data.jobs)) {
        return data.jobs
      }
      return []
    } catch (error) {
      console.error('Error fetching print jobs:', error)
      return []
    }
  },

  async fetchPrintJobById(id) {
    const response = await api.get(`/manager/print-jobs/${id}`)
    return response.data
  },

  async fetchPendingJobs() {
    try {
      const response = await api.get('/manager/print-jobs/pending')
      const data = response.data
      if (Array.isArray(data)) {
        return data
      } else if (data && Array.isArray(data.data)) {
        return data.data
      } else if (data && Array.isArray(data.jobs)) {
        return data.jobs
      }
      return []
    } catch (error) {
      console.error('Error fetching pending jobs:', error)
      return []
    }
  },

  async fetchFailedJobs() {
    try {
      const response = await api.get('/manager/print-jobs/failed')
      const data = response.data
      if (Array.isArray(data)) {
        return data
      } else if (data && Array.isArray(data.data)) {
        return data.data
      } else if (data && Array.isArray(data.jobs)) {
        return data.jobs
      }
      return []
    } catch (error) {
      console.error('Error fetching failed jobs:', error)
      return []
    }
  },

  async retryJob(id) {
    const response = await api.post(`/manager/print-jobs/${id}/retry`)
    return response.data
  },

  async cancelJob(id) {
    await api.delete(`/manager/print-jobs/${id}`)
  },

  async reprintBill(orderId) {
    const response = await api.post(`/orders/${orderId}/reprint-bill`)
    return response.data
  },

  async reprintLabel(orderId, itemIndex = null) {
    const url = itemIndex !== null 
      ? `/orders/${orderId}/reprint-label?item_index=${itemIndex}`
      : `/orders/${orderId}/reprint-label`
    const response = await api.post(url)
    return response.data
  }
}
