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
    
    // Loading state
    loading: false,
    
    // Error message
    error: null
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
      this.loading = false
      this.error = null
    }
  }
})

