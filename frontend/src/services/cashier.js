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
  getOrderAudits: (orderId) => api.get(`/cashier/orders/${orderId}/audits`)
}