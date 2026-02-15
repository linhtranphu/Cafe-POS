import api from './api'

/**
 * Batch Report Service
 * Handles API calls for batch report operations
 */
class BatchReportService {
  /**
   * Fetch production report
   * @param {Object} params - Report parameters
   * @param {string} params.from_date - Start date (ISO format)
   * @param {string} params.to_date - End date (ISO format)
   * @param {string} params.batch_definition_id - Optional: Filter by batch definition ID
   * @param {string} params.prepared_by - Optional: Filter by preparer
   * @returns {Promise<Object>} Production report data
   */
  async fetchProductionReport(params) {
    try {
      const queryParams = new URLSearchParams()
      
      if (params.from_date) queryParams.append('from_date', params.from_date)
      if (params.to_date) queryParams.append('to_date', params.to_date)
      if (params.batch_definition_id) queryParams.append('batch_definition_id', params.batch_definition_id)
      if (params.prepared_by) queryParams.append('prepared_by', params.prepared_by)

      const response = await api.get(`/batch-reports/production?${queryParams.toString()}`)
      return response.data
    } catch (error) {
      throw this._handleError(error, 'Không thể tải báo cáo sản xuất')
    }
  }

  /**
   * Fetch wastage report
   * @param {Object} params - Report parameters
   * @param {string} params.from_date - Start date (ISO format)
   * @param {string} params.to_date - End date (ISO format)
   * @param {string} params.batch_definition_id - Optional: Filter by batch definition ID
   * @returns {Promise<Object>} Wastage report data
   */
  async fetchWastageReport(params) {
    try {
      const queryParams = new URLSearchParams()
      
      if (params.from_date) queryParams.append('from_date', params.from_date)
      if (params.to_date) queryParams.append('to_date', params.to_date)
      if (params.batch_definition_id) queryParams.append('batch_definition_id', params.batch_definition_id)

      const response = await api.get(`/batch-reports/wastage?${queryParams.toString()}`)
      return response.data
    } catch (error) {
      throw this._handleError(error, 'Không thể tải báo cáo lãng phí')
    }
  }

  /**
   * Fetch usage report
   * @param {Object} params - Report parameters
   * @param {string} params.from_date - Start date (ISO format)
   * @param {string} params.to_date - End date (ISO format)
   * @param {string} params.batch_definition_id - Optional: Filter by batch definition ID
   * @param {string} params.menu_item_id - Optional: Filter by menu item ID
   * @returns {Promise<Object>} Usage report data
   */
  async fetchUsageReport(params) {
    try {
      const queryParams = new URLSearchParams()
      
      if (params.from_date) queryParams.append('from_date', params.from_date)
      if (params.to_date) queryParams.append('to_date', params.to_date)
      if (params.batch_definition_id) queryParams.append('batch_definition_id', params.batch_definition_id)
      if (params.menu_item_id) queryParams.append('menu_item_id', params.menu_item_id)

      const response = await api.get(`/batch-reports/usage?${queryParams.toString()}`)
      return response.data
    } catch (error) {
      throw this._handleError(error, 'Không thể tải báo cáo sử dụng')
    }
  }

  /**
   * Export report to CSV
   * @param {string} reportType - Type of report ('production', 'wastage', 'usage')
   * @param {Object} params - Report parameters (same as fetch methods)
   * @returns {Promise<Blob>} CSV file blob
   */
  async exportReport(reportType, params) {
    try {
      const queryParams = new URLSearchParams()
      
      if (params.from_date) queryParams.append('from_date', params.from_date)
      if (params.to_date) queryParams.append('to_date', params.to_date)
      if (params.batch_definition_id) queryParams.append('batch_definition_id', params.batch_definition_id)
      if (params.prepared_by) queryParams.append('prepared_by', params.prepared_by)
      if (params.menu_item_id) queryParams.append('menu_item_id', params.menu_item_id)

      const response = await api.get(
        `/batch-reports/${reportType}?${queryParams.toString()}`,
        {
          headers: {
            'Accept': 'text/csv'
          },
          responseType: 'blob'
        }
      )
      
      return response.data
    } catch (error) {
      throw this._handleError(error, 'Không thể xuất báo cáo')
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
      if (serverMessage === 'invalid date range') {
        return new Error('Khoảng thời gian không hợp lệ')
      }
      if (serverMessage === 'from_date and to_date are required') {
        return new Error('Vui lòng chọn khoảng thời gian')
      }
      if (serverMessage === 'invalid batch_definition_id') {
        return new Error('ID định nghĩa batch không hợp lệ')
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

export default new BatchReportService()
