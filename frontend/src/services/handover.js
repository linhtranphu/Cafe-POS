import api from './api'

/**
 * Cash Handover Service
 * Handles API calls for waiter-to-cashier cash handover functionality
 */
export const handoverService = {
  // ==================== WAITER ENDPOINTS ====================
  
  /**
   * Create a partial or full handover request
   * @param {string} shiftId - Waiter shift ID
   * @param {Object} data - { declared_amount, handover_type, waiter_note }
   * @returns {Promise} Handover object
   */
  createHandover: (shiftId, data) => 
    api.post(`/shifts/${shiftId}/handover`, data),
  
  /**
   * Create handover and end shift in one action
   * @param {string} shiftId - Waiter shift ID
   * @param {Object} data - { declared_amount, waiter_note, end_cash }
   * @returns {Promise} Handover object
   */
  createHandoverAndEndShift: (shiftId, data) => 
    api.post(`/shifts/${shiftId}/handover-and-end`, data),
  
  /**
   * Get pending handover for a shift
   * @param {string} shiftId - Waiter shift ID
   * @returns {Promise} Pending handover or null
   */
  getPendingHandover: (shiftId) => 
    api.get(`/shifts/${shiftId}/pending-handover`),
  
  /**
   * Get handover history for a shift
   * @param {string} shiftId - Waiter shift ID
   * @returns {Promise} Array of handovers
   */
  getHandoverHistory: (shiftId) => 
    api.get(`/shifts/${shiftId}/handovers`),
  
  /**
   * Cancel a pending handover
   * @param {string} handoverId - Handover ID
   * @returns {Promise} Success message
   */
  cancelHandover: (handoverId) => 
    api.delete(`/cash-handovers/${handoverId}`),
  
  // ==================== CASHIER ENDPOINTS ====================
  
  /**
   * Get all pending handovers for current cashier
   * @returns {Promise} Array of pending handovers
   */
  getPendingHandovers: () => 
    api.get('/cash-handovers/pending'),
  
  /**
   * Get today's handovers
   * @returns {Promise} Array of today's handovers
   */
  getTodayHandovers: () => 
    api.get('/cash-handovers/today'),
  
  /**
   * Confirm or reject a handover with reconciliation
   * @param {string} handoverId - Handover ID
   * @param {Object} data - { actual_amount, status, cashier_note, discrepancy_reason, discrepancy_responsibility }
   * @returns {Promise} Success message
   */
  confirmHandover: (handoverId, data) => 
    api.post(`/cash-handovers/${handoverId}/confirm`, data),
  
  /**
   * Quick confirm/reject without detailed reconciliation
   * @param {string} handoverId - Handover ID
   * @param {string} status - 'CONFIRMED' or 'REJECTED'
   * @returns {Promise} Success message
   */
  quickConfirm: (handoverId, status) => 
    api.post(`/cash-handovers/${handoverId}/quick-confirm`, { status }),
  
  // ==================== MANAGER ENDPOINTS ====================
  
  /**
   * Get handovers requiring manager approval (large discrepancies)
   * @returns {Promise} Array of handovers requiring approval
   */
  getPendingApprovals: () => 
    api.get('/cash-handovers/pending-approval'),
  
  /**
   * Approve or reject a discrepancy
   * @param {string} handoverId - Handover ID
   * @param {Object} data - { approved, manager_note }
   * @returns {Promise} Success message
   */
  approveDiscrepancy: (handoverId, data) => 
    api.post(`/cash-handovers/${handoverId}/approve`, data),
  
  /**
   * Get discrepancy statistics for a date range
   * @param {string} startDate - Start date (YYYY-MM-DD)
   * @param {string} endDate - End date (YYYY-MM-DD)
   * @returns {Promise} Discrepancy statistics
   */
  getDiscrepancyStats: (startDate, endDate) => 
    api.get(`/cash-handovers/discrepancy-stats?start=${startDate}&end=${endDate}`)
}
