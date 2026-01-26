import { defineStore } from 'pinia'
import { expenseService } from '../services/expense'

export const useExpenseStore = defineStore('expense', {
  state: () => ({
    expenses: [],
    categories: [],
    recurringExpenses: [],
    prepaidExpenses: [],
    loading: false,
    error: null
  }),

  actions: {
    async fetchExpenses(filters = {}) {
      this.loading = true
      this.error = null
      try {
        this.expenses = await expenseService.getExpenses(filters)
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải chi phí'
        this.expenses = []
      } finally {
        this.loading = false
      }
    },

    async createExpense(expense) {
      this.error = null
      try {
        const newExpense = await expenseService.createExpense(expense)
        this.expenses.push(newExpense)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tạo chi phí'
        return false
      }
    },

    async updateExpense(id, expense) {
      this.error = null
      try {
        const updated = await expenseService.updateExpense(id, expense)
        const index = this.expenses.findIndex(e => e.id === id)
        if (index !== -1) {
          this.expenses[index] = updated
        }
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi cập nhật chi phí'
        return false
      }
    },

    async deleteExpense(id) {
      this.error = null
      try {
        await expenseService.deleteExpense(id)
        this.expenses = this.expenses.filter(e => e.id !== id)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi xóa chi phí'
        return false
      }
    },

    async fetchCategories() {
      try {
        this.categories = await expenseService.getCategories()
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải danh mục'
        this.categories = []
      }
    },

    async createCategory(category) {
      this.error = null
      try {
        const newCategory = await expenseService.createCategory(category)
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
        await expenseService.deleteCategory(id)
        this.categories = this.categories.filter(c => c.id !== id)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi xóa danh mục'
        return false
      }
    },

    async fetchRecurringExpenses() {
      try {
        this.recurringExpenses = await expenseService.getRecurringExpenses()
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải chi phí định kỳ'
        this.recurringExpenses = []
      }
    },

    async createRecurringExpense(recurring) {
      this.error = null
      try {
        const newRecurring = await expenseService.createRecurringExpense(recurring)
        this.recurringExpenses.push(newRecurring)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tạo chi phí định kỳ'
        return false
      }
    },

    async deleteRecurringExpense(id) {
      this.error = null
      try {
        await expenseService.deleteRecurringExpense(id)
        this.recurringExpenses = this.recurringExpenses.filter(r => r.id !== id)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi xóa chi phí định kỳ'
        return false
      }
    },

    async fetchPrepaidExpenses() {
      try {
        this.prepaidExpenses = await expenseService.getPrepaidExpenses()
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải chi phí trả trước'
        this.prepaidExpenses = []
      }
    },

    async createPrepaidExpense(prepaid) {
      this.error = null
      try {
        const newPrepaid = await expenseService.createPrepaidExpense(prepaid)
        this.prepaidExpenses.push(newPrepaid)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tạo chi phí trả trước'
        return false
      }
    },

    async deletePrepaidExpense(id) {
      this.error = null
      try {
        await expenseService.deletePrepaidExpense(id)
        this.prepaidExpenses = this.prepaidExpenses.filter(p => p.id !== id)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi xóa chi phí trả trước'
        return false
      }
    },

    async getExpenseReport(filters = {}) {
      try {
        return await expenseService.getExpenseReport(filters)
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải báo cáo'
        return null
      }
    }
  }
})
