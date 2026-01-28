import { defineStore } from 'pinia'
import { orderService } from '../services/order'

export const useOrderStore = defineStore('order', {
  state: () => ({
    orders: [],
    currentOrder: null,
    loading: false,
    error: null
  }),

  actions: {
    async fetchOrders() {
      this.loading = true
      this.error = null
      try {
        this.orders = await orderService.getMyOrders() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải orders'
        this.orders = []
      } finally {
        this.loading = false
      }
    },

    async fetchAllOrders() {
      this.loading = true
      this.error = null
      try {
        this.orders = await orderService.getAllOrders() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải orders'
        this.orders = []
      } finally {
        this.loading = false
      }
    },

    async createOrder(orderData) {
      this.error = null
      try {
        const order = await orderService.createOrder(orderData)
        this.orders.unshift(order)
        return order
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tạo order'
        throw error
      }
    },

    async collectPayment(id, paymentData) {
      this.error = null
      try {
        const order = await orderService.collectPayment(id, paymentData)
        this.updateOrderInList(order)
        return order
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi thu tiền'
        throw error
      }
    },

    async editOrder(id, orderData) {
      this.error = null
      try {
        const response = await orderService.editOrder(id, orderData)
        this.updateOrderInList(response.order)
        return response
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi chỉnh sửa order'
        throw error
      }
    },

    async sendToBar(id) {
      this.error = null
      try {
        const order = await orderService.sendToBar(id)
        this.updateOrderInList(order)
        return order
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi gửi quầy bar'
        throw error
      }
    },

    async serveOrder(id) {
      this.error = null
      try {
        const order = await orderService.serveOrder(id)
        this.updateOrderInList(order)
        return order
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi phục vụ'
        throw error
      }
    },

    async cancelOrder(id, reason) {
      this.error = null
      try {
        const order = await orderService.cancelOrder(id, reason)
        this.updateOrderInList(order)
        return order
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi hủy order'
        throw error
      }
    },

    async refundPartial(id, amount, reason) {
      this.error = null
      try {
        const order = await orderService.refundPartial(id, amount, reason)
        this.updateOrderInList(order)
        return order
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi hoàn tiền'
        throw error
      }
    },

    updateOrderInList(order) {
      const index = this.orders.findIndex(o => o.id === order.id)
      if (index !== -1) {
        this.orders[index] = order
      }
    },

    setCurrentOrder(order) {
      this.currentOrder = order
    },

    clearCurrentOrder() {
      this.currentOrder = null
    }
  },

  getters: {
    ordersByStatus: (state) => (status) => {
      return state.orders.filter(o => o.status === status)
    },

    createdOrders: (state) => {
      return state.orders.filter(o => o.status === 'CREATED')
    },

    paidOrders: (state) => {
      return state.orders.filter(o => o.status === 'PAID')
    },

    inProgressOrders: (state) => {
      return state.orders.filter(o => o.status === 'IN_PROGRESS')
    },

    servedOrders: (state) => {
      return state.orders.filter(o => o.status === 'SERVED')
    }
  }
})
