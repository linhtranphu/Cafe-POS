import api from './api'

/**
 * Service for managing ingredient restocks paid from fund
 * Requirements: 2.1, 2.5
 */
export const fundIngredientService = {
  /**
   * Restock ingredient with payment from fund
   * @param {string} ingredientId - Ingredient ID
   * @param {Object} restockData - Restock details
   * @param {number} restockData.quantity - Quantity to restock
   * @param {number} restockData.cost_per_unit - Cost per unit
   * @param {string} restockData.reason - Reason for restock
   * @returns {Promise<Object>} Restock record, expense, and fund transaction
   */
  async restockIngredientFromFund(ingredientId, restockData) {
    try {
      const response = await api.post(
        `/manager/ingredients/${ingredientId}/restock/from-fund`,
        restockData
      )
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
      
      // Handle not found error
      if (error.response?.status === 404) {
        throw new Error('Không tìm thấy nguyên liệu')
      }
      
      // Handle other errors
      throw new Error('Không thể nhập kho từ quỹ. Vui lòng thử lại.')
    }
  },

  /**
   * Get restock history for an ingredient
   * @param {string} ingredientId - Ingredient ID
   * @param {Object} filter - Filter options
   * @param {number} filter.limit - Limit results
   * @param {number} filter.offset - Offset for pagination
   * @returns {Promise<Object>} Restock history with fund transaction info
   */
  async getRestockHistory(ingredientId, filter = {}) {
    try {
      const params = new URLSearchParams(filter)
      
      const response = await api.get(
        `/manager/ingredients/${ingredientId}/restock-history?${params}`
      )
      return response.data
    } catch (error) {
      // Handle not found error
      if (error.response?.status === 404) {
        throw new Error('Không tìm thấy nguyên liệu')
      }
      
      // Handle other errors
      throw new Error('Không thể tải lịch sử nhập kho. Vui lòng thử lại.')
    }
  }
}
