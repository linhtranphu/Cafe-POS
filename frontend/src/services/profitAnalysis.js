import api from './api.js'

/**
 * @typedef {import('./types/menuCost').DateRange} DateRange
 * @typedef {import('./types/menuCost').CategoryProfitResponse} CategoryProfitResponse
 * @typedef {import('./types/menuCost').OperatingProfitReport} OperatingProfitReport
 * @typedef {import('./types/menuCost').OperatingExpense} OperatingExpense
 * @typedef {import('./types/menuCost').OperatingExpensesResponse} OperatingExpensesResponse
 */

/**
 * Profit Analysis API Service
 * Provides methods for retrieving profit analysis reports
 */
export const profitAnalysisService = {
  /**
   * Get category-level profit analysis for a date range
   * @param {DateRange} dateRange - Date range for the report
   * @returns {Promise<CategoryProfitResponse>} Category profit response
   */
  async getCategoryProfit(dateRange) {
    const params = new URLSearchParams()
    params.append('start_date', dateRange.start)
    params.append('end_date', dateRange.end)
    
    const response = await api.get(`/reports/category-profit?${params.toString()}`)
    return response.data
  },

  /**
   * Get operating profit analysis for a date range
   * @param {DateRange} dateRange - Date range for the report
   * @returns {Promise<OperatingProfitReport>} Operating profit report
   */
  async getOperatingProfit(dateRange) {
    const params = new URLSearchParams()
    params.append('start_date', dateRange.start)
    params.append('end_date', dateRange.end)
    
    const response = await api.get(`/reports/operating-profit?${params.toString()}`)
    return response.data
  },

  /**
   * Create or update operating expense
   * @param {OperatingExpense} data - Operating expense data
   * @returns {Promise<OperatingExpense>} Created/updated operating expense
   */
  async createOperatingExpense(data) {
    const response = await api.post('/manager/expenses', data)
    return response.data
  },

  /**
   * Get operating expenses for a date range
   * @param {DateRange} [dateRange] - Optional date range filter
   * @returns {Promise<OperatingExpensesResponse>} Operating expenses response
   */
  async getOperatingExpenses(dateRange) {
    if (dateRange) {
      const params = new URLSearchParams()
      params.append('start_date', dateRange.start)
      params.append('end_date', dateRange.end)
      
      const response = await api.get(`/manager/expenses?${params.toString()}`)
      return response.data
    }
    
    const response = await api.get('/manager/expenses')
    return response.data
  }
}
