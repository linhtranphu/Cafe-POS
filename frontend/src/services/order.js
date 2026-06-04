import api from './api'

export const orderService = {
  async createOrder(order) {
    const response = await api.post('/waiter/orders', order)
    return response.data
  },

  async collectPayment(id, paymentData) {
    const response = await api.post(`/waiter/orders/${id}/payment`, paymentData)
    return response.data
  },

  async editOrder(id, orderData) {
    const response = await api.put(`/waiter/orders/${id}/edit`, orderData)
    return response.data
  },

  async sendToBar(id) {
    const response = await api.post(`/waiter/orders/${id}/send`)
    return response.data
  },

  async serveOrder(id) {
    const response = await api.post(`/waiter/orders/${id}/serve`)
    return response.data
  },

  async cancelOrder(id, reason) {
    const response = await api.post(`/waiter/orders/${id}/cancel`, { reason })
    return response.data
  },

  async changePaymentMethod(id, paymentMethod) {
    const response = await api.patch(`/waiter/orders/${id}/payment-method`, { payment_method: paymentMethod })
    return response.data
  },

  async renameOrder(id, customerName) {
    const response = await api.patch(`/waiter/orders/${id}/rename`, { customer_name: customerName })
    return response.data
  },

  async printTemporaryBill(id) {
    const response = await api.post(`/waiter/orders/${id}/print-temp-bill`)
    return response.data
  },

  async mergeOrders(orderIds, customerName, note) {
    const response = await api.post('/waiter/orders/merge', {
      order_ids: orderIds,
      customer_name: customerName,
      note: note
    })
    return response.data
  },

  async refundPartial(id, amount, reason) {
    const response = await api.post(`/cashier/orders/${id}/refund`, { amount, reason })
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
