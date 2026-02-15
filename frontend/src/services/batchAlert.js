import api from './api'

/**
 * Batch Alert Service
 * Handles API calls for batch alert operations
 */
class BatchAlertService {
  /**
   * Fetch all batch alerts (low stock, expiring, expired)
   * @returns {Promise<Object>} Alert data with low_stock, expiring, expired arrays and last_checked timestamp
   */
  async fetchAlerts() {
    try {
      const response = await api.get('/batch-alerts')
      return response.data
    } catch (error) {
      throw this._handleError(error, 'Không thể tải cảnh báo batch')
    }
  }

  /**
   * Poll alerts at regular intervals
   * @param {Function} callback - Callback function to handle alert data
   * @param {number} interval - Polling interval in milliseconds (default: 5 minutes)
   * @returns {number} Interval ID for clearing
   */
  startPolling(callback, interval = 300000) {
    // Fetch immediately
    this.fetchAlerts()
      .then(callback)
      .catch(error => console.error('Error fetching alerts:', error))

    // Then poll at intervals
    return setInterval(() => {
      this.fetchAlerts()
        .then(callback)
        .catch(error => console.error('Error fetching alerts:', error))
    }, interval)
  }

  /**
   * Stop polling alerts
   * @param {number} intervalId - Interval ID returned by startPolling
   */
  stopPolling(intervalId) {
    if (intervalId) {
      clearInterval(intervalId)
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

export default new BatchAlertService()
