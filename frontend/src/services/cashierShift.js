import api from './api'

/**
 * Cashier Shift Service
 * 
 * Handles API calls for cashier shift management.
 * This is separate from the regular shift service to maintain clear separation
 * between cashier shifts and waiter/barista shifts.
 */

export default {
  /**
   * Start a new cashier shift
   * @param {number} startingFloat - Initial cash amount in the drawer
   * @returns {Promise<Object>} The created cashier shift
   */
  async startCashierShift(startingFloat) {
    const response = await api.post('/cashier-shifts', {
      starting_float: startingFloat
    })
    return response.data
  },

  /**
   * Get the current open cashier shift for the authenticated cashier
   * @returns {Promise<Object|null>} The current cashier shift or null if none exists
   */
  async getCurrentCashierShift() {
    try {
      const response = await api.get('/cashier-shifts/current')
      return response.data
    } catch (error) {
      // Return null if no shift found (404)
      if (error.response?.status === 404) {
        return null
      }
      throw error
    }
  },

  /**
   * Get all cashier shifts (manager only)
   * @returns {Promise<Array>} List of all cashier shifts
   */
  async getAllCashierShifts() {
    const response = await api.get('/cashier-shifts')
    return response.data
  },

  /**
   * Get all cashier shifts for the authenticated cashier
   * @returns {Promise<Array>} List of cashier shifts
   */
  async getMyCashierShifts() {
    const response = await api.get('/cashier-shifts/my-shifts')
    return response.data
  },

  /**
   * Get a specific cashier shift by ID
   * @param {string} shiftId - The cashier shift ID
   * @returns {Promise<Object>} The cashier shift
   */
  async getCashierShift(shiftId) {
    const response = await api.get(`/cashier-shifts/${shiftId}`)
    return response.data
  },

  /**
   * Initiate shift closure
   * @param {string} shiftId - The cashier shift ID
   * @returns {Promise<Object>} The updated shift
   */
  async initiateClosure(shiftId) {
    const response = await api.post(`/cashier-shifts/${shiftId}/initiate-closure`, {})
    return response.data
  },

  /**
   * Record actual cash counted
   * @param {string} shiftId - The cashier shift ID
   * @param {number} actualCash - The actual cash amount
   * @returns {Promise<Object>} The updated shift with variance
   */
  async recordActualCash(shiftId, actualCash) {
    const response = await api.post(`/cashier-shifts/${shiftId}/record-actual-cash`, {
      actual_cash: actualCash
    })
    return response.data
  },

  /**
   * Document variance reason and notes
   * @param {string} shiftId - The cashier shift ID
   * @param {Object} data - Variance documentation data
   * @param {string} data.reason - Variance reason
   * @param {string} data.notes - Detailed notes
   * @returns {Promise<Object>} The updated shift
   */
  async documentVariance(shiftId, data) {
    const response = await api.post(`/cashier-shifts/${shiftId}/document-variance`, data)
    return response.data
  },

  /**
   * Confirm cashier's responsibility
   * @param {string} shiftId - The cashier shift ID
   * @returns {Promise<Object>} The updated shift
   */
  async confirmResponsibility(shiftId) {
    const response = await api.post(`/cashier-shifts/${shiftId}/confirm-responsibility`, {})
    return response.data
  },

  /**
   * Close the cashier shift
   * @param {string} shiftId - The cashier shift ID
   * @returns {Promise<Object>} The closed shift
   */
  async closeShift(shiftId) {
    const response = await api.post(`/cashier-shifts/${shiftId}/close`, {})
    return response.data
  },

  /**
   * Check if there are any open waiter shifts
   * @returns {Promise<Object>} Status of waiter shifts
   */
  async checkWaiterShifts() {
    const response = await api.get('/cashier-shifts/check-waiter-shifts')
    return response.data
  }
}
