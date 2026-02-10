import api from './api.js'

/**
 * @typedef {import('./types/menuCost').OperatingExpense} OperatingExpense
 * @typedef {import('./types/menuCost').OperatingExpensesResponse} OperatingExpensesResponse
 * @typedef {import('./types/menuCost').DateRange} DateRange
 */

/**
 * Operating Expense API Service
 * Provides methods for managing operating expenses
 */
export const operatingExpenseService = {
  /**
   * Create or update an operating expense
   * @param {OperatingExpense} data - Operating expense data
   * @returns {Promise<OperatingExpense>} Created/updated operating expense
   */
  async createOperatingExpense(data) {
    const response = await api.post('/operating-expenses', data)
    return response.data
  },

  /**
   * Get operating expenses with optional date range filter
   * @param {DateRange} [dateRange] - Optional date range filter
   * @returns {Promise<OperatingExpensesResponse>} Operating expenses response
   */
  async getOperatingExpenses(dateRange) {
    if (dateRange) {
      const params = new URLSearchParams()
      params.append('start_date', dateRange.start)
      params.append('end_date', dateRange.end)
      
      const response = await api.get(`/operating-expenses?${params.toString()}`)
      return response.data
    }
    
    const response = await api.get('/operating-expenses')
    return response.data
  }
}
