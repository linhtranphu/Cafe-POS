import api from './api'

/**
 * Service for managing expenses paid from fund
 * Requirements: 1.1, 4.3
 */
export const fundExpenseService = {
  /**
   * Create an expense paid from fund
   * @param {Object} expenseData - Expense details
   * @param {string} expenseData.date - Expense date (ISO format)
   * @param {string} expenseData.category_id - Category ID
   * @param {number} expenseData.amount - Expense amount
   * @param {string} expenseData.description - Expense description
   * @param {string} expenseData.vendor - Vendor name (optional)
   * @param {string} expenseData.notes - Additional notes (optional)
   * @returns {Promise<Object>} Created expense and fund transaction
   */
  async createExpenseFromFund(expenseData) {
    try {
      const response = await api.post('/manager/expenses/from-fund', expenseData)
      return response.data
    } catch (error) {
      // Handle insufficient balance error
      if (error.response?.status === 400 && error.response?.data?.error?.includes('insufficient')) {
        throw new Error('Số dư quỹ không đủ để thực hiện giao dịch này')
      }
      
      // Handle validation errors
      if (error.response?.status === 400) {
        throw new Error(error.response?.data?.error || 'Dữ liệu không hợp lệ')
      }
      
      // Handle other errors - expose actual backend error message
      throw new Error(error.response?.data?.error || 'Không thể tạo chi phí từ quỹ. Vui lòng thử lại.')
    }
  },

  /**
   * Get expenses that were paid from fund
   * @param {Object} filter - Filter options
   * @param {string} filter.from_date - Start date (ISO format)
   * @param {string} filter.to_date - End date (ISO format)
   * @param {string} filter.category_id - Category ID
   * @param {number} filter.limit - Limit results
   * @param {number} filter.offset - Offset for pagination
   * @returns {Promise<Object>} List of expenses paid from fund
   */
  async getExpensesPaidFromFund(filter = {}) {
    try {
      const params = new URLSearchParams({
        ...filter,
        paid_from_fund: 'true'
      })
      
      const response = await api.get(`/manager/expenses?${params}`)
      return response.data
    } catch (error) {
      throw new Error('Không thể tải danh sách chi phí từ quỹ. Vui lòng thử lại.')
    }
  }
}
