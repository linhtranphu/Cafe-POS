import api from './api'

/**
 * Shift Service
 * 
 * Handles API calls for waiter and barista shift management.
 * Note: Cashier shifts are handled separately in cashierShift.js service.
 */

export const shiftService = {
  /**
   * Start a new shift (waiter or barista only)
   * Note: Cashiers should use cashierShift.startCashierShift() instead
   * @param {Object} shiftData - Shift data including type and start_cash
   * @returns {Promise<Object>} The created shift
   */
  async startShift(shiftData) {
    const response = await api.post('/shifts/start', shiftData)
    return response.data
  },

  /**
   * End a shift (waiter or barista)
   * @param {string} id - Shift ID
   * @param {number} endCash - Final cash amount
   * @returns {Promise<Object>} The updated shift
   */
  async endShift(id, endCash) {
    const response = await api.post(`/shifts/${id}/end`, { end_cash: endCash })
    return response.data
  },

  /**
   * Close a shift and lock orders (waiter or barista)
   * @param {string} id - Shift ID
   * @param {number} endCash - Final cash amount
   * @returns {Promise<Object>} The closed shift
   */
  async closeShift(id, endCash) {
    const response = await api.post(`/shifts/${id}/close`, { end_cash: endCash })
    return response.data
  },

  /**
   * Get the current open shift for the authenticated user
   * Note: This returns waiter/barista shifts only, not cashier shifts
   * @returns {Promise<Object>} The current shift
   */
  async getCurrentShift() {
    const response = await api.get('/shifts/current')
    return response.data
  },

  /**
   * Get all shifts for the authenticated user
   * Note: This returns waiter/barista shifts only, not cashier shifts
   * @returns {Promise<Array>} List of shifts
   */
  async getMyShifts() {
    const response = await api.get('/shifts/my')
    return response.data
  },

  /**
   * Get all shifts (manager only)
   * Note: This returns waiter/barista shifts only
   * For cashier shifts, use cashierShift.getAllCashierShifts()
   * @returns {Promise<Array>} List of all shifts
   */
  async getAllShifts() {
    const response = await api.get('/shifts')
    return response.data
  },

  /**
   * Get a specific shift by ID
   * @param {string} id - Shift ID
   * @returns {Promise<Object>} The shift
   */
  async getShift(id) {
    const response = await api.get(`/shifts/${id}`)
    return response.data
  }
}

