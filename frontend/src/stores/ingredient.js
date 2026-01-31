import { defineStore } from 'pinia'
import { ingredientService } from '../services/ingredient'

export const useIngredientStore = defineStore('ingredient', {
  state: () => ({
    items: [],
    lowStockItems: [],
    categories: [],
    loading: false,
    error: null
  }),

  actions: {
    async fetchIngredients() {
      this.loading = true
      this.error = null
      try {
        this.items = await ingredientService.getIngredients() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải nguyên liệu'
        this.items = []
      } finally {
        this.loading = false
      }
    },

    async createIngredient(ingredient) {
      this.error = null
      try {
        const newItem = await ingredientService.createIngredient(ingredient)
        this.items.push(newItem)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tạo nguyên liệu'
        return false
      }
    },

    async updateIngredient(id, ingredient) {
      this.error = null
      try {
        const updatedItem = await ingredientService.updateIngredient(id, ingredient)
        const index = this.items.findIndex(i => i.id === id)
        if (index !== -1) {
          this.items[index] = updatedItem
        }
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi cập nhật nguyên liệu'
        return false
      }
    },

    async deleteIngredient(id) {
      this.error = null
      try {
        await ingredientService.deleteIngredient(id)
        this.items = this.items.filter(i => i.id !== id)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi xóa nguyên liệu'
        return false
      }
    },

    async fetchLowStock() {
      try {
        this.lowStockItems = await ingredientService.getLowStockItems() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải nguyên liệu sắp hết'
        this.lowStockItems = []
      }
    },

    async fetchLowStockItems() {
      try {
        return await ingredientService.getLowStockItems() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải nguyên liệu sắp hết'
        return []
      }
    },

    async fetchStockHistory(id) {
      try {
        return await ingredientService.getStockHistory(id) || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải lịch sử tồn kho'
        return []
      }
    },

    async adjustStock(id, adjustment) {
      this.error = null
      try {
        const updatedItem = await ingredientService.adjustStock(id, adjustment)
        const index = this.items.findIndex(i => i.id === id)
        if (index !== -1) {
          this.items[index] = updatedItem
        }
        await this.fetchLowStock()
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi điều chỉnh tồn kho'
        return false
      }
    },

    // Category actions
    async fetchCategories() {
      this.error = null
      try {
        this.categories = await ingredientService.getCategories() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải danh mục'
        this.categories = []
      }
    },

    async createCategory(category) {
      this.error = null
      try {
        const newCategory = await ingredientService.createCategory(category)
        this.categories.push(newCategory)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tạo danh mục'
        return false
      }
    },

    async deleteCategory(id) {
      this.error = null
      try {
        await ingredientService.deleteCategory(id)
        this.categories = this.categories.filter(c => c.id !== id)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi xóa danh mục'
        return false
      }
    }
  }
})
