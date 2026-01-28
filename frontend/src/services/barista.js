import api from './api'

export const baristaService = {
  // Get queued orders (waiting for barista)
  async getQueuedOrders() {
    const response = await api.get('/barista/orders/queue')
    return response.data
  },

  // Get my orders (in progress + ready)
  async getMyOrders() {
    const response = await api.get('/barista/orders/my')
    return response.data
  },

  // Accept order from queue
  async acceptOrder(id) {
    const response = await api.post(`/barista/orders/${id}/accept`)
    return response.data
  },

  // Mark order as ready
  async markReady(id) {
    const response = await api.post(`/barista/orders/${id}/ready`)
    return response.data
  },

  // Get order details
  async getOrder(id) {
    const response = await api.get(`/barista/orders/${id}`)
    return response.data
  }
}
