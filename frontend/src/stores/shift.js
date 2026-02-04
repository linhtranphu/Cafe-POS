import { defineStore } from 'pinia'
import { shiftService } from '../services/shift'

/**
 * Shift Store
 * 
 * Manages state for waiter and barista shifts.
 * Note: Cashier shifts are handled separately in cashierShift.js store.
 */

export const useShiftStore = defineStore('shift', {
  state: () => ({
    // Current open shift for the authenticated user (waiter/barista only)
    currentShift: null,
    
    // List of shifts (waiter/barista only)
    shifts: [],
    
    // Cash handover state
    pendingHandover: null,
    handoverHistory: [],
    
    // Loading state
    loading: false,
    handoverLoading: false,
    
    // Error message
    error: null,
    handoverError: null
  }),

  getters: {
    /**
     * Check if the user has an open shift
     * @returns {boolean} True if there's an open waiter/barista shift
     */
    hasOpenShift: (state) => {
      return state.currentShift !== null && state.currentShift.status === 'OPEN'
    },

    /**
     * Get all open shifts
     * @returns {Array} List of open shifts
     */
    openShifts: (state) => {
      return state.shifts.filter(s => s.status === 'OPEN')
    },

    /**
     * Get all closed shifts
     * @returns {Array} List of closed shifts
     */
    closedShifts: (state) => {
      return state.shifts.filter(s => s.status === 'CLOSED')
    },

    /**
     * Get waiter shifts only
     * @returns {Array} List of waiter shifts
     */
    waiterShifts: (state) => {
      return state.shifts.filter(s => s.role_type === 'waiter')
    },

    /**
     * Get barista shifts only
     * @returns {Array} List of barista shifts
     */
    baristaShifts: (state) => {
      return state.shifts.filter(s => s.role_type === 'barista')
    },

    /**
     * Check if current shift has pending handover
     * @returns {boolean} True if there's a pending handover
     */
    hasPendingHandover: (state) => {
      return state.pendingHandover !== null && state.pendingHandover.status === 'PENDING'
    },

    /**
     * Get available cash for handover
     * @returns {number} Available cash amount
     */
    availableCash: (state) => {
      return state.currentShift?.remaining_cash || 0
    },

    /**
     * Check if shift can be ended (no pending handover and remaining cash is 0)
     * @returns {boolean} True if shift can be ended
     */
    canEndShift: (state) => {
      return !state.pendingHandover && (state.currentShift?.remaining_cash || 0) === 0
    }
  },

  actions: {
    /**
     * Start a new shift (waiter or barista only)
     * Note: Cashiers should use cashierShiftStore.startCashierShift() instead
     * @param {Object} shiftData - Shift data including type and start_cash
     * @returns {Promise<Object>} The created shift
     */
    async startShift(shiftData) {
      this.error = null
      try {
        const shift = await shiftService.startShift(shiftData)
        this.currentShift = shift
        return shift
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi mở ca'
        throw error
      }
    },

    /**
     * End a shift (waiter or barista)
     * @param {string} id - Shift ID
     * @param {number} endCash - Final cash amount
     * @returns {Promise<Object>} The updated shift
     */
    async endShift(id, endCash) {
      this.error = null
      try {
        const shift = await shiftService.endShift(id, endCash)
        this.currentShift = null
        return shift
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi kết ca'
        throw error
      }
    },

    /**
     * Close a shift and lock orders (waiter or barista)
     * @param {string} id - Shift ID
     * @param {number} endCash - Final cash amount
     * @returns {Promise<Object>} The closed shift
     */
    async closeShift(id, endCash) {
      this.error = null
      try {
        const shift = await shiftService.closeShift(id, endCash)
        return shift
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi chốt ca'
        throw error
      }
    },

    /**
     * Fetch the current open shift for the authenticated user
     * Note: This fetches waiter/barista shifts only, not cashier shifts
     * @returns {Promise<void>}
     */
    async fetchCurrentShift() {
      this.loading = true
      this.error = null
      try {
        this.currentShift = await shiftService.getCurrentShift()
      } catch (error) {
        this.currentShift = null
        if (error.response?.data?.error === 'no open shift found') {
          this.error = null
        } else {
          this.error = error.response?.data?.error || 'Lỗi tải ca'
        }
      } finally {
        this.loading = false
      }
    },

    /**
     * Fetch all shifts for the authenticated user
     * Note: This fetches waiter/barista shifts only, not cashier shifts
     * @returns {Promise<void>}
     */
    async fetchMyShifts() {
      this.loading = true
      this.error = null
      try {
        this.shifts = await shiftService.getMyShifts() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải shifts'
        this.shifts = []
      } finally {
        this.loading = false
      }
    },

    /**
     * Fetch all shifts (manager only)
     * Note: This fetches waiter/barista shifts only
     * For cashier shifts, use cashierShiftStore.fetchAllCashierShifts()
     * @returns {Promise<void>}
     */
    async fetchAllShifts() {
      this.loading = true
      this.error = null
      try {
        this.shifts = await shiftService.getAllShifts() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải shifts'
        this.shifts = []
      } finally {
        this.loading = false
      }
    },

    /**
     * Clear error message
     */
    clearError() {
      this.error = null
    },

    /**
     * Reset the entire store state
     */
    reset() {
      this.currentShift = null
      this.shifts = []
      this.pendingHandover = null
      this.handoverHistory = []
      this.loading = false
      this.handoverLoading = false
      this.error = null
      this.handoverError = null
    },

    // Cash Handover Actions

    /**
     * Create a cash handover request
     * @param {string} shiftId - Shift ID
     * @param {Object} handoverData - Handover data
     * @returns {Promise<Object>} The created handover
     */
    async createCashHandover(shiftId, handoverData) {
      this.handoverLoading = true
      this.handoverError = null
      try {
        const handover = await shiftService.createCashHandover(shiftId, handoverData)
        this.pendingHandover = handover
        // Refresh current shift to update cash amounts
        await this.fetchCurrentShift()
        return handover
      } catch (error) {
        this.handoverError = error.response?.data?.error || 'Lỗi tạo yêu cầu bàn giao'
        throw error
      } finally {
        this.handoverLoading = false
      }
    },

    /**
     * Create handover and end shift
     * @param {string} shiftId - Shift ID
     * @param {Object} handoverData - Handover data
     * @returns {Promise<Object>} The created handover
     */
    async createHandoverAndEndShift(shiftId, handoverData) {
      this.handoverLoading = true
      this.handoverError = null
      try {
        const handover = await shiftService.createHandoverAndEndShift(shiftId, handoverData)
        this.pendingHandover = handover
        // Refresh current shift
        await this.fetchCurrentShift()
        return handover
      } catch (error) {
        this.handoverError = error.response?.data?.error || 'Lỗi bàn giao và kết ca'
        throw error
      } finally {
        this.handoverLoading = false
      }
    },

    /**
     * Get pending handover for a shift
     * @param {string} shiftId - Shift ID
     * @returns {Promise<void>}
     */
    async fetchPendingHandover(shiftId) {
      this.handoverLoading = true
      this.handoverError = null
      try {
        this.pendingHandover = await shiftService.getPendingHandover(shiftId)
      } catch (error) {
        this.pendingHandover = null
        if (error.response?.status !== 404) {
          this.handoverError = error.response?.data?.error || 'Lỗi tải handover'
        }
      } finally {
        this.handoverLoading = false
      }
    },

    /**
     * Get handover history for a shift
     * @param {string} shiftId - Shift ID
     * @returns {Promise<void>}
     */
    async fetchHandoverHistory(shiftId) {
      this.handoverLoading = true
      this.handoverError = null
      try {
        const response = await shiftService.getHandoverHistory(shiftId)
        this.handoverHistory = response.handovers || []
      } catch (error) {
        this.handoverHistory = []
        this.handoverError = error.response?.data?.error || 'Lỗi tải lịch sử bàn giao'
      } finally {
        this.handoverLoading = false
      }
    },

    /**
     * Cancel a handover
     * @param {string} handoverId - Handover ID
     * @returns {Promise<void>}
     */
    async cancelHandover(handoverId) {
      this.handoverLoading = true
      this.handoverError = null
      try {
        await shiftService.cancelHandover(handoverId)
        this.pendingHandover = null
        // Refresh current shift
        await this.fetchCurrentShift()
      } catch (error) {
        this.handoverError = error.response?.data?.error || 'Lỗi hủy bàn giao'
        throw error
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

