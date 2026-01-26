import api from './api'

export const expenseService = {
  // FR-EX-01: Get expenses
  async getExpenses(filters = {}) {
    const params = new URLSearchParams(filters)
    const response = await api.get(`/manager/expenses?${params}`)
    return response.data
  },

  // FR-EX-03: Create expense
  async createExpense(expense) {
    const response = await api.post('/manager/expenses', expense)
    return response.data
  },

  // Update expense
  async updateExpense(id, expense) {
    const response = await api.put(`/manager/expenses/${id}`, expense)
    return response.data
  },

  // Delete expense
  async deleteExpense(id) {
    await api.delete(`/manager/expenses/${id}`)
  },

  // FR-EX-02: Category management
  async getCategories() {
    const response = await api.get('/manager/expense-categories')
    return response.data
  },

  async createCategory(category) {
    const response = await api.post('/manager/expense-categories', category)
    return response.data
  },

  async deleteCategory(id) {
    await api.delete(`/manager/expense-categories/${id}`)
  },

  // FR-EX-04: Recurring expenses
  async getRecurringExpenses() {
    const response = await api.get('/manager/recurring-expenses')
    return response.data
  },

  async createRecurringExpense(recurring) {
    const response = await api.post('/manager/recurring-expenses', recurring)
    return response.data
  },

  async deleteRecurringExpense(id) {
    await api.delete(`/manager/recurring-expenses/${id}`)
  },

  async getRecurringReminders() {
    const response = await api.get('/manager/recurring-expenses/reminders')
    return response.data
  },

  // FR-EX-05: Prepaid expenses
  async getPrepaidExpenses() {
    const response = await api.get('/manager/prepaid-expenses')
    return response.data
  },

  async createPrepaidExpense(prepaid) {
    const response = await api.post('/manager/prepaid-expenses', prepaid)
    return response.data
  },

  async deletePrepaidExpense(id) {
    await api.delete(`/manager/prepaid-expenses/${id}`)
  },

  // FR-EX-06: Reports
  async getExpenseReport(filters = {}) {
    const params = new URLSearchParams(filters)
    const response = await api.get(`/manager/expenses/report?${params}`)
    return response.data
  }
}
