import { defineStore } from 'pinia'
import { menuService } from '../services/menu'

export const useMenuStore = defineStore('menu', {
  state: () => ({
    items: [],
    loading: false,
    error: null
  }),

  getters: {
    // Get variant by ID from a menu item
    getVariantById: (state) => (menuItemId, variantId) => {
      const item = state.items.find(i => i.id === menuItemId)
      if (!item || !item.has_variants || !item.variants) {
        return null
      }
      return item.variants.find(v => v.id === variantId) || null
    },

    // Get default variant from a menu item
    getDefaultVariant: (state) => (menuItemId) => {
      const item = state.items.find(i => i.id === menuItemId)
      if (!item || !item.has_variants || !item.variants) {
        return null
      }
      const defaultVariant = item.variants.find(v => v.is_default)
      // Fallback to first variant if no default found
      return defaultVariant || (item.variants.length > 0 ? item.variants[0] : null)
    }
  },

  actions: {
    async fetchMenuItems() {
      this.loading = true
      this.error = null
      
      try {
        this.items = await menuService.getMenuItems()
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải menu'
      } finally {
        this.loading = false
      }
    },

    async createMenuItem(item) {
      try {
        const newItem = await menuService.createMenuItem(item)
        this.items.push(newItem)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tạo món'
        return false
      }
    },

    async updateMenuItem(id, item) {
      this.error = null
      try {
        const updatedItem = await menuService.updateMenuItem(id, item)
        const index = this.items.findIndex(i => i.id === id)
        if (index !== -1) {
          this.items[index] = updatedItem
        }
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi cập nhật món'
        console.error('Update error:', error)
        return false
      }
    },

    async deleteMenuItem(id) {
      this.error = null
      try {
        await menuService.deleteMenuItem(id)
        this.items = this.items.filter(i => i.id !== id)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi xóa món'
        console.error('Delete error:', error)
        return false
      }
    }
  }
})