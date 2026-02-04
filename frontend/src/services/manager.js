import api from './api'

/**
 * Manager Service
 * 
 * Handles API calls for manager-specific functionality including discrepancy approvals.
 */

export const managerService = {
  /**
   * Get pending approvals for manager
   * @returns {Promise<Object>} Pending approvals response
   */
  getPendingApprovals: () => api.get('/manager/cash-handovers/pending-approval'),

  /**
   * Approve or reject a discrepancy
   * @param {string} handoverId - Handover ID
   * @param {Object} data - Approval data
   * @returns {Promise<Object>} Success response
   */
  approveDiscrepancy: (handoverId, data) => 
    api.post(`/manager/cash-handovers/${handoverId}/approve`, data),

  /**
   * Get discrepancy statistics for managers
   * @param {string} startDate - Start date (YYYY-MM-DD)
   * @param {string} endDate - End date (YYYY-MM-DD)
   * @returns {Promise<Object>} Statistics response
   */
  getDiscrepancyStats: (startDate, endDate) => {
    const params = new URLSearchParams()
    if (startDate) params.append('start_date', startDate)
    if (endDate) params.append('end_date', endDate)
    return api.get(`/manager/discrepancies/stats?${params.toString()}`)
  }
}