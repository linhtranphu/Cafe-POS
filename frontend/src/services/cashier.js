import api from './api'

export const cashierService = {
  // Shift Management
  getShiftStatus: (shiftId) => api.get(`/cashier/shifts/${shiftId}/status`),
  
  // Payment Oversight
  getPaymentsByShift: (shiftId) => api.get(`/cashier/shifts/${shiftId}/payments`),
  reportDiscrepancy: (data) => api.post('/cashier/discrepancies', data),
  getPendingDiscrepancies: () => api.get('/cashier/discrepancies/pending'),
  resolveDiscrepancy: (discrepancyId) => api.post(`/cashier/discrepancies/${discrepancyId}/resolve`),
  overridePayment: (orderId, data) => api.post(`/cashier/orders/${orderId}/override`, data),
  lockOrder: (orderId) => api.post(`/cashier/orders/${orderId}/lock`),
  
  // Reconciliation
  reconcileCash: (data) => api.post('/cashier/reconcile/cash', data),
  
  // Reports
  generateShiftReport: (shiftId) => api.get(`/cashier/reports/shift/${shiftId}`),
  getDailyReport: (date) => api.get(`/cashier/reports/daily?date=${date}`),
  handoverShift: (data) => api.post('/cashier/handover', data),
  
  // Audit
  getOrderAudits: (orderId) => api.get(`/cashier/orders/${orderId}/audits`),

  // Cash Handover Functions
  
  /**
   * Get pending handovers for cashier
   * @returns {Promise<Object>} Pending handovers response
   */
  getPendingHandovers: () => api.get('/cash-handovers/pending'),

  /**
   * Get today's handovers
   * @returns {Promise<Object>} Today's handovers response
   */
  getTodayHandovers: () => api.get('/cash-handovers/today'),

  /**
   * Quick confirm handover with exact amount
   * @param {string} handoverId - Handover ID
   * @param {Object} data - Confirmation data
   * @returns {Promise<Object>} Success response
   */
  quickConfirmHandover: (handoverId, data) => 
    api.post(`/cash-handovers/${handoverId}/quick-confirm`, data),

  /**
   * Reconcile handover with actual amount and discrepancy details
   * @param {string} handoverId - Handover ID
   * @param {Object} reconcileData - Reconciliation data
   * @returns {Promise<Object>} Success response
   */
  reconcileHandover: (handoverId, reconcileData) => 
    api.post(`/cash-handovers/${handoverId}/reconcile`, reconcileData),

  /**
   * Reject handover
   * @param {string} handoverId - Handover ID
   * @param {Object} data - Rejection data with reason
   * @returns {Promise<Object>} Success response
   */
  rejectHandover: (handoverId, data) => 
    api.post(`/cash-handovers/${handoverId}/reject`, data),

  /**
   * Get discrepancy statistics
   * @param {string} startDate - Start date (YYYY-MM-DD)
   * @param {string} endDate - End date (YYYY-MM-DD)
   * @returns {Promise<Object>} Statistics response
   */
  getDiscrepancyStats: (startDate, endDate) => {
    const params = new URLSearchParams()
    if (startDate) params.append('start_date', startDate)
    if (endDate) params.append('end_date', endDate)
    return api.get(`/cash-handovers/discrepancy-stats?${params.toString()}`)
  }
}