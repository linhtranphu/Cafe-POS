import api from './api'

export const orderService = {
  async createOrder(order) {
    const response = await api.post('/waiter/orders', order)
    return response.data
  },

  async confirmOrder(id, discount = 0) {
    const response = await api.put(`/waiter/orders/${id}/confirm`, { discount })
    return response.data
  },

  async payOrder(id, paymentMethod) {
    const response = await api.post(`/waiter/orders/${id}/payment`, { payment_method: paymentMethod })
    return response.data
  },

  async sendToKitchen(id) {
    const response = await api.post(`/waiter/orders/${id}/send`)
    return response.data
  },

  async serveOrder(id) {
    const response = await api.post(`/waiter/orders/${id}/serve`)
    return response.data
  },

  async cancelOrder(id, reason) {
    const response = await api.post(`/cashier/orders/${id}/cancel`, { reason })
    return response.data
  },

  async refundOrder(id, reason) {
    const response = await api.post(`/cashier/orders/${id}/refund`, { reason })
    return response.data
  },

  async lockOrder(id) {
    const response = await api.post(`/cashier/orders/${id}/lock`)
    return response.data
  },

  async getMyOrders() {
    const response = await api.get('/waiter/orders')
    return response.data
  },

  async getAllOrders() {
    const response = await api.get('/cashier/orders')
    return response.data
  },

  async getOrder(id) {
    const response = await api.get(`/waiter/orders/${id}`)
    return response.data
  }
}
