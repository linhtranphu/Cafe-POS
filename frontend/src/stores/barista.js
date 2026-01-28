import { defineStore } from 'pinia'
import { baristaService } from '../services/barista'

export const useBaristaStore = defineStore('barista', {
  state: () => ({
    queuedOrders: [],
    myOrders: [],
    currentOrder: null,
    loading: false,
    error: null
  }),

  actions: {
    async fetchQueuedOrders() {
      this.loading = true
      this.error = null
      try {
        this.queuedOrders = await baristaService.getQueuedOrders() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải queue'
        this.queuedOrders = []
      } finally {
        this.loading = false
      }
    },

    async fetchMyOrders() {
      this.loading = true
      this.error = null
      try {
        this.myOrders = await baristaService.getMyOrders() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải orders'
        this.myOrders = []
      } finally {
        this.loading = false
      }
    },

    async acceptOrder(id) {
      this.error = null
      try {
        const order = await baristaService.acceptOrder(id)
        // Remove from queue
        this.queuedOrders = this.queuedOrders.filter(o => o.id !== id)
        // Add to my orders
        this.myOrders.unshift(order)
        return order
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi nhận order'
        throw error
      }
    },

    async markReady(id) {
      this.error = null
      try {
        const order = await baristaService.markReady(id)
        // Update in my orders
        const index = this.myOrders.findIndex(o => o.id === id)
        if (index !== -1) {
          this.myOrders[index] = order
        }
        return order
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi đánh dấu hoàn tất'
        throw error
      }
    },

    async getOrder(id) {
      this.error = null
      try {
        this.currentOrder = await baristaService.getOrder(id)
        return this.currentOrder
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải order'
        throw error
      }
    },

    clearCurrentOrder() {
      this.currentOrder = null
    }
  },

  getters: {
    inProgressOrders: (state) => {
      return state.myOrders.filter(o => o.status === 'IN_PROGRESS')
    },

    readyOrders: (state) => {
      return state.myOrders.filter(o => o.status === 'READY')
    },

    servedOrders: (state) => {
      return state.myOrders.filter(o => o.status === 'SERVED')
    },

    queueCount: (state) => {
      return state.queuedOrders.length
    },

    inProgressCount: (state) => {
      return state.myOrders.filter(o => o.status === 'IN_PROGRESS').length
    },

    readyCount: (state) => {
      return state.myOrders.filter(o => o.status === 'READY').length
    },

    servedCount: (state) => {
      return state.myOrders.filter(o => o.status === 'SERVED').length
    }
  }
})
