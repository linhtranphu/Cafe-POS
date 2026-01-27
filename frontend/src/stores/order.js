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

    async confirmOrder(id, discount) {
      this.error = null
      try {
        const order = await orderService.confirmOrder(id, discount)
        this.updateOrderInList(order)
        return order
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi xác nhận order'
        throw error
      }
    },

    async payOrder(id, paymentMethod) {
      this.error = null
      try {
        const order = await orderService.payOrder(id, paymentMethod)
        this.updateOrderInList(order)
        return order
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi thanh toán'
        throw error
      }
    },

    async sendToKitchen(id) {
      this.error = null
      try {
        const order = await orderService.sendToKitchen(id)
        this.updateOrderInList(order)
        return order
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi gửi pha chế'
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

    async refundOrder(id, reason) {
      this.error = null
      try {
        const order = await orderService.refundOrder(id, reason)
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

    unpaidOrders: (state) => {
      return state.orders.filter(o => o.status === 'UNPAID')
    },

    paidOrders: (state) => {
      return state.orders.filter(o => o.status === 'PAID')
    },

    inProgressOrders: (state) => {
      return state.orders.filter(o => o.status === 'IN_PROGRESS')
    }
  }
})
