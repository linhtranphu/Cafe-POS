import { defineStore } from 'pinia'
import { managerService } from '../services/manager'

/**
 * Manager Store
 * 
 * Manages state for manager-specific functionality including discrepancy approvals.
 */

export const useManagerStore = defineStore('manager', {
  state: () => ({
    // Discrepancy approval state
    pendingApprovals: [],
    discrepancyStats: null,
    
    // Loading state
    loading: false,
    approvalLoading: false,
    
    // Error messages
    error: null,
    approvalError: null
  }),

  getters: {
    /**
     * Get count of pending approvals
     * @returns {number} Number of pending approvals
     */
    pendingApprovalCount: (state) => state.pendingApprovals.length,

    /**
     * Get large discrepancies (requiring approval)
     * @returns {Array} List of large discrepancies
     */
    largeDiscrepancies: (state) => 
      state.pendingApprovals.filter(h => h.requires_manager_approval),

    /**
     * Get shortage discrepancies
     * @returns {Array} List of shortage discrepancies
     */
    shortageDiscrepancies: (state) => 
      state.pendingApprovals.filter(h => h.discrepancy_amount < 0),

    /**
     * Get overage discrepancies
     * @returns {Array} List of overage discrepancies
     */
    overageDiscrepancies: (state) => 
      state.pendingApprovals.filter(h => h.discrepancy_amount > 0),

    /**
     * Check if there are stats available
     * @returns {boolean} True if stats are loaded
     */
    hasDiscrepancyStats: (state) => !!state.discrepancyStats,

    /**
     * Get total discrepancy amount
     * @returns {number} Total discrepancy amount
     */
    totalDiscrepancyAmount: (state) => {
      if (!state.discrepancyStats) return 0
      return state.discrepancyStats.net_discrepancy || 0
    }
  },

  actions: {
    /**
     * Fetch pending approvals
     * @returns {Promise<void>}
     */
    async fetchPendingApprovals() {
      this.approvalLoading = true
      this.approvalError = null
      try {
        const response = await managerService.getPendingApprovals()
        this.pendingApprovals = response.handovers || []
      } catch (error) {
        this.pendingApprovals = []
        this.approvalError = error.response?.data?.error || 'Lỗi tải danh sách phê duyệt'
      } finally {
        this.approvalLoading = false
      }
    },

    /**
     * Approve or reject a discrepancy
     * @param {string} handoverId - Handover ID
     * @param {boolean} approved - Whether to approve or reject
     * @param {string} managerNotes - Manager notes
     * @returns {Promise<void>}
     */
    async approveDiscrepancy(handoverId, approved, managerNotes) {
      this.approvalLoading = true
      this.approvalError = null
      try {
        await managerService.approveDiscrepancy(handoverId, {
          approved,
          manager_notes: managerNotes
        })
        // Refresh pending approvals
        await this.fetchPendingApprovals()
      } catch (error) {
        this.approvalError = error.response?.data?.error || 'Lỗi phê duyệt chênh lệch'
        throw error
      } finally {
        this.approvalLoading = false
      }
    },

    /**
     * Get discrepancy statistics for managers
     * @param {string} startDate - Start date (YYYY-MM-DD)
     * @param {string} endDate - End date (YYYY-MM-DD)
     * @returns {Promise<void>}
     */
    async getDiscrepancyStats(startDate, endDate) {
      this.loading = true
      this.error = null
      try {
        const response = await managerService.getDiscrepancyStats(startDate, endDate)
        this.discrepancyStats = response.stats
      } catch (error) {
        this.discrepancyStats = null
        this.error = error.response?.data?.error || 'Lỗi tải thống kê chênh lệch'
      } finally {
        this.loading = false
      }
    },

    /**
     * Clear approval error message
     */
    clearApprovalError() {
      this.approvalError = null
    },

    /**
     * Clear general error message
     */
    clearError() {
      this.error = null
    },

    /**
     * Reset the entire store state
     */
    reset() {
      this.pendingApprovals = []
      this.discrepancyStats = null
      this.loading = false
      this.approvalLoading = false
      this.error = null
      this.approvalError = null
    }
  }
})