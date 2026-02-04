import { defineStore } from 'pinia'
import { cashierService } from '../services/cashier'

export const useCashierStore = defineStore('cashier', {
  state: () => ({
    shiftStatus: null,
    payments: [],
    discrepancies: [],
    reconciliation: null,
    reports: [],
    audits: [],
    
    // Cash handover state
    pendingHandovers: [],
    todayHandovers: [],
    discrepancyStats: null,
    discrepancyThreshold: 50000, // VND
    
    loading: false,
    handoverLoading: false,
    error: null,
    handoverError: null
  }),

  getters: {
    pendingDiscrepancies: (state) => 
      state.discrepancies.filter(d => d.status === 'PENDING'),
    
    resolvedDiscrepancies: (state) => 
      state.discrepancies.filter(d => d.status === 'RESOLVED'),
    
    cashPayments: (state) => 
      state.payments.filter(p => p.payment_method === 'CASH'),
    
    transferPayments: (state) => 
      state.payments.filter(p => p.payment_method === 'TRANSFER'),
    
    qrPayments: (state) => 
      state.payments.filter(p => p.payment_method === 'QR'),
    
    totalCashAmount: (state) => 
      state.payments
        .filter(p => p.payment_method === 'CASH')
        .reduce((sum, p) => sum + p.amount, 0),
    
    hasReconciliation: (state) => !!state.reconciliation,
    
    reconciliationStatus: (state) => state.reconciliation?.status || null,

    // Cash handover getters
    pendingHandoverCount: (state) => state.pendingHandovers.length,
    
    todayHandoverCount: (state) => state.todayHandovers.length,
    
    confirmedHandovers: (state) => 
      state.todayHandovers.filter(h => h.status === 'CONFIRMED'),
    
    discrepancyHandovers: (state) => 
      state.todayHandovers.filter(h => h.status === 'DISCREPANCY'),
    
    rejectedHandovers: (state) => 
      state.todayHandovers.filter(h => h.status === 'REJECTED'),
    
    totalHandoverAmount: (state) => 
      state.todayHandovers
        .filter(h => h.status === 'CONFIRMED')
        .reduce((sum, h) => sum + (h.actual_amount || h.requested_amount), 0),
    
    hasDiscrepancyStats: (state) => !!state.discrepancyStats
  },

  actions: {
    // FR-CASH-02: Theo dõi trạng thái ca
    async getShiftStatus(shiftId) {
      this.loading = true
      this.error = null
      try {
        const response = await cashierService.getShiftStatus(shiftId)
        this.shiftStatus = response.data
        return response.data
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to get shift status'
        throw error
      } finally {
        this.loading = false
      }
    },

    // FR-CASH-04: Giám sát thanh toán
    async getPaymentsByShift(shiftId) {
      this.loading = true
      this.error = null
      try {
        const response = await cashierService.getPaymentsByShift(shiftId)
        this.payments = response.data.payments || []
        return response.data.payments
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to get payments'
        throw error
      } finally {
        this.loading = false
      }
    },

    // FR-CASH-05: Xử lý sai lệch thanh toán
    async reportDiscrepancy(data) {
      this.loading = true
      this.error = null
      try {
        await cashierService.reportDiscrepancy(data)
        await this.getPendingDiscrepancies()
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to report discrepancy'
        throw error
      } finally {
        this.loading = false
      }
    },

    async getPendingDiscrepancies() {
      try {
        const response = await cashierService.getPendingDiscrepancies()
        this.discrepancies = response.data.discrepancies || []
        return response.data.discrepancies
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to get discrepancies'
        throw error
      }
    },

    async resolveDiscrepancy(discrepancyId) {
      this.loading = true
      this.error = null
      try {
        await cashierService.resolveDiscrepancy(discrepancyId)
        await this.getPendingDiscrepancies()
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to resolve discrepancy'
        throw error
      } finally {
        this.loading = false
      }
    },

    // FR-CASH-06: Đối soát tiền mặt
    async reconcileCash(data) {
      this.loading = true
      this.error = null
      try {
        const response = await cashierService.reconcileCash(data)
        this.reconciliation = response.data
        return response.data
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to reconcile cash'
        throw error
      } finally {
        this.loading = false
      }
    },

    // FR-CASH-08: Hủy/điều chỉnh thanh toán
    async overridePayment(orderId, reason) {
      this.loading = true
      this.error = null
      try {
        await cashierService.overridePayment(orderId, { reason })
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to override payment'
        throw error
      } finally {
        this.loading = false
      }
    },

    // FR-CASH-09: Khóa order
    async lockOrder(orderId) {
      this.loading = true
      this.error = null
      try {
        await cashierService.lockOrder(orderId)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to lock order'
        throw error
      } finally {
        this.loading = false
      }
    },

    // FR-CASH-10: Báo cáo ca
    async generateShiftReport(shiftId) {
      this.loading = true
      this.error = null
      try {
        const response = await cashierService.generateShiftReport(shiftId)
        const report = response.data
        this.reports.unshift(report)
        return report
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to generate report'
        throw error
      } finally {
        this.loading = false
      }
    },

    async getDailyReport(date) {
      this.loading = true
      this.error = null
      try {
        const response = await cashierService.getDailyReport(date)
        return response.data
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to get daily report'
        throw error
      } finally {
        this.loading = false
      }
    },

    // FR-CASH-11: Bàn giao ca
    async handoverShift(data) {
      this.loading = true
      this.error = null
      try {
        await cashierService.handoverShift(data)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to handover shift'
        throw error
      } finally {
        this.loading = false
      }
    },

    async getOrderAudits(orderId) {
      this.loading = true
      this.error = null
      try {
        const response = await cashierService.getOrderAudits(orderId)
        this.audits = response.data.audits || []
        return response.data.audits
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to get audits'
        throw error
      } finally {
        this.loading = false
      }
    },

    clearError() {
      this.error = null
    },

    reset() {
      this.shiftStatus = null
      this.payments = []
      this.discrepancies = []
      this.reconciliation = null
      this.reports = []
      this.audits = []
      this.pendingHandovers = []
      this.todayHandovers = []
      this.discrepancyStats = null
      this.loading = false
      this.handoverLoading = false
      this.error = null
      this.handoverError = null
    },

    // Cash Handover Actions

    /**
     * Fetch pending handovers for cashier
     * @returns {Promise<void>}
     */
    async fetchPendingHandovers() {
      this.handoverLoading = true
      this.handoverError = null
      try {
        const response = await cashierService.getPendingHandovers()
        this.pendingHandovers = response.handovers || []
      } catch (error) {
        this.pendingHandovers = []
        this.handoverError = error.response?.data?.error || 'Lỗi tải handovers chờ xử lý'
      } finally {
        this.handoverLoading = false
      }
    },

    /**
     * Fetch today's handovers
     * @returns {Promise<void>}
     */
    async fetchTodayHandovers() {
      this.handoverLoading = true
      this.handoverError = null
      try {
        const response = await cashierService.getTodayHandovers()
        this.todayHandovers = response.handovers || []
      } catch (error) {
        this.todayHandovers = []
        this.handoverError = error.response?.data?.error || 'Lỗi tải handovers hôm nay'
      } finally {
        this.handoverLoading = false
      }
    },

    /**
     * Quick confirm handover with exact amount
     * @param {string} handoverId - Handover ID
     * @param {string} cashierNotes - Optional notes
     * @returns {Promise<void>}
     */
    async quickConfirm(handoverId, cashierNotes = '') {
      this.handoverLoading = true
      this.handoverError = null
      try {
        await cashierService.quickConfirmHandover(handoverId, { cashier_notes: cashierNotes })
        // Refresh handovers
        await Promise.all([
          this.fetchPendingHandovers(),
          this.fetchTodayHandovers()
        ])
      } catch (error) {
        this.handoverError = error.response?.data?.error || 'Lỗi xác nhận handover'
        throw error
      } finally {
        this.handoverLoading = false
      }
    },

    /**
     * Reconcile handover with actual amount and discrepancy details
     * @param {string} handoverId - Handover ID
     * @param {Object} reconcileData - Reconciliation data
     * @returns {Promise<void>}
     */
    async reconcileHandover(handoverId, reconcileData) {
      this.handoverLoading = true
      this.handoverError = null
      try {
        await cashierService.reconcileHandover(handoverId, reconcileData)
        // Refresh handovers
        await Promise.all([
          this.fetchPendingHandovers(),
          this.fetchTodayHandovers()
        ])
      } catch (error) {
        this.handoverError = error.response?.data?.error || 'Lỗi đối soát handover'
        throw error
      } finally {
        this.handoverLoading = false
      }
    },

    /**
     * Reject handover
     * @param {string} handoverId - Handover ID
     * @param {string} reason - Rejection reason
     * @returns {Promise<void>}
     */
    async rejectHandover(handoverId, reason) {
      this.handoverLoading = true
      this.handoverError = null
      try {
        await cashierService.rejectHandover(handoverId, { reason })
        // Refresh handovers
        await Promise.all([
          this.fetchPendingHandovers(),
          this.fetchTodayHandovers()
        ])
      } catch (error) {
        this.handoverError = error.response?.data?.error || 'Lỗi từ chối handover'
        throw error
      } finally {
        this.handoverLoading = false
      }
    },

    /**
     * Get discrepancy statistics
     * @param {string} startDate - Start date (YYYY-MM-DD)
     * @param {string} endDate - End date (YYYY-MM-DD)
     * @returns {Promise<void>}
     */
    async getDiscrepancyStats(startDate, endDate) {
      this.handoverLoading = true
      this.handoverError = null
      try {
        const response = await cashierService.getDiscrepancyStats(startDate, endDate)
        this.discrepancyStats = response.stats
      } catch (error) {
        this.discrepancyStats = null
        this.handoverError = error.response?.data?.error || 'Lỗi tải thống kê chênh lệch'
      } finally {
        this.handoverLoading = false
      }
    },

    /**
     * Clear handover error message
     */
    clearHandoverError() {
      this.handoverError = null
    }
  }
})