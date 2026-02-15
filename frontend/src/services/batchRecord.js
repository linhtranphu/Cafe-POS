import api from './api'

/**
 * Batch Record Service
 * Handles API calls for batch record operations
 */
class BatchRecordService {
  /**
   * Fetch batch records with optional filters
   * @param {Object} filters - Filter options
   * @param {string} filters.batch_definition_id - Filter by batch definition ID
   * @param {string} filters.status - Filter by status (available, expired, depleted)
   * @param {string} filters.prepared_by - Filter by preparer
   * @param {string} filters.from_date - Filter from date (ISO format)
   * @param {string} filters.to_date - Filter to date (ISO format)
   * @param {number} filters.page - Page number
   * @param {number} filters.limit - Items per page
   * @returns {Promise<Object>} Response with data, total, page, limit
   */
  async fetchRecords(filters = {}) {
    try {
      const params = new URLSearchParams()
      
      if (filters.batch_definition_id) params.append('batch_definition_id', filters.batch_definition_id)
      if (filters.status) params.append('status', filters.status)
      if (filters.prepared_by) params.append('prepared_by', filters.prepared_by)
      if (filters.from_date) params.append('from_date', filters.from_date)
      if (filters.to_date) params.append('to_date', filters.to_date)
      if (filters.page) params.append('page', filters.page)
      if (filters.limit) params.append('limit', filters.limit)

      const response = await api.get(`/batch-records?${params.toString()}`)
      return response.data
    } catch (error) {
      throw this._handleError(error, 'Không thể tải danh sách batch records')
    }
  }

  /**
   * Fetch a single batch record by ID
   * @param {string} id - Batch record ID
   * @returns {Promise<Object>} Batch record data
   */
  async fetchRecordById(id) {
    try {
      const response = await api.get(`/batch-records/${id}`)
      return response.data
    } catch (error) {
      throw this._handleError(error, 'Không thể tải batch record')
    }
  }

  /**
   * Create a new batch record
   * @param {Object} data - Batch record data
   * @param {string} data.batch_definition_id - Batch definition ID
   * @param {number} data.quantity_produced - Quantity produced
   * @param {string} data.prepared_by - Preparer user ID
   * @returns {Promise<Object>} Created batch record
   */
  async createRecord(data) {
    try {
      const response = await api.post(`/batch-records`, data)
      return response.data
    } catch (error) {
      throw this._handleError(error, 'Không thể tạo batch record')
    }
  }

  /**
   * Update batch record quantity
   * @param {string} id - Batch record ID
   * @param {number} quantity_remaining - New quantity remaining
   * @returns {Promise<Object>} Updated batch record
   */
  async updateQuantity(id, quantity_remaining) {
    try {
      const response = await api.patch(
        `/batch-records/${id}/quantity`,
        { quantity_remaining }
      )
      return response.data
    } catch (error) {
      throw this._handleError(error, 'Không thể cập nhật số lượng batch')
    }
  }

  /**
   * Mark batch record as expired
   * @param {string} id - Batch record ID
   * @returns {Promise<Object>} Updated batch record
   */
  async markAsExpired(id) {
    try {
      const response = await api.patch(`/batch-records/${id}/expire`)
      return response.data
    } catch (error) {
      throw this._handleError(error, 'Không thể đánh dấu batch là hết hạn')
    }
  }

  /**
   * Delete a batch record
   * @param {string} id - Batch record ID
   * @returns {Promise<void>}
   */
  async deleteRecord(id) {
    try {
      await api.delete(`/batch-records/${id}`)
    } catch (error) {
      throw this._handleError(error, 'Không thể xóa batch record')
    }
  }

  /**
   * Handle API errors and return user-friendly messages
   * @private
   * @param {Error} error - Axios error
   * @param {string} defaultMessage - Default error message
   * @returns {Error} Formatted error
   */
  _handleError(error, defaultMessage) {
    if (error.response) {
      // Server responded with error
      const serverMessage = error.response.data?.error || error.response.data?.message
      
      // Map specific error messages to Vietnamese
      if (serverMessage === 'insufficient ingredients') {
        return new Error('Không đủ nguyên liệu để chế biến batch')
      }
      if (serverMessage === 'cannot delete batch that has been partially used') {
        return new Error('Không thể xóa batch đã được sử dụng một phần')
      }
      if (serverMessage === 'batch record not found') {
        return new Error('Không tìm thấy batch record')
      }
      if (serverMessage === 'invalid id' || serverMessage === 'invalid batch_definition_id') {
        return new Error('ID không hợp lệ')
      }
      
      return new Error(serverMessage || defaultMessage)
    } else if (error.request) {
      // Request made but no response
      return new Error('Không thể kết nối đến server')
    } else {
      // Something else happened
      return new Error(error.message || defaultMessage)
    }
  }
}

export default new BatchRecordService()
