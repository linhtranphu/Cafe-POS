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

// Helper functions for cart management with variant support
// These are exported separately as they don't need store state
export const cartHelpers = {
  // Create cart item from menu item and optional variant
  createCartItem(menuItem, variant = null) {
    const item = {
      menu_item_id: menuItem.id,
      name: menuItem.name,
      quantity: 1
    }

    if (variant) {
      // Multi-size item with variant
      item.variant_id = variant.id
      item.variant_name = variant.name
      item.price = variant.price
    } else {
      // Single-size item (backward compatible)
      item.price = menuItem.price
    }

    return item
  },

  // Check if two cart items are the same (including variant)
  isSameCartItem(item1, item2) {
    if (item1.menu_item_id !== item2.menu_item_id) {
      return false
    }
    // For multi-size items, variant_id must match
    if (item1.variant_id || item2.variant_id) {
      return item1.variant_id === item2.variant_id
    }
    // For single-size items, just menu_item_id match is enough
    return true
  }
}
