import { defineStore } from 'pinia'
import cashierShiftService from '../services/cashierShift'

/**
 * Cashier Shift Store
 * 
 * Manages state for cashier shifts.
 * This is separate from the regular shift store to maintain clear separation
 * between cashier shifts and waiter/barista shifts.
 */

export const useCashierShiftStore = defineStore('cashierShift', {
  state: () => ({
    // Current open cashier shift for the authenticated cashier
    currentCashierShift: null,
    
    // List of all cashier shifts (for managers or history view)
    cashierShifts: [],
    
    // Alias for consistency with other stores
    shifts: [],
    
    // Loading state
    loading: false,
    
    // Error message
    error: null
  }),

  getters: {
    /**
     * Check if the cashier has an open shift
     * @returns {boolean} True if there's an open cashier shift
     */
    hasOpenCashierShift: (state) => {
      return state.currentCashierShift && state.currentCashierShift.status === 'OPEN'
    },

    /**
     * Check if the cashier can start a new shift
     * @returns {boolean} True if no open shift exists
     */
    canStartCashierShift: (state) => {
      return !state.currentCashierShift || state.currentCashierShift.status === 'CLOSED'
    },

    /**
     * Check if the current shift is in closure process
     * @returns {boolean} True if shift is in closure initiated state
     */
    isClosureInitiated: (state) => {
      return state.currentCashierShift && state.currentCashierShift.status === 'CLOSURE_INITIATED'
    },

    /**
     * Check if the current shift is closed
     * @returns {boolean} True if shift is closed
     */
    isClosed: (state) => {
      return state.currentCashierShift && state.currentCashierShift.status === 'CLOSED'
    },

    /**
     * Get the current shift ID
     * @returns {string|null} The current shift ID or null
     */
    currentShiftId: (state) => {
      return state.currentCashierShift?.id || null
    }
  },

  actions: {
    /**
     * Start a new cashier shift
     * @param {number} startingFloat - Initial cash amount in the drawer
     * @returns {Promise<Object>} The created shift
     */
    async startCashierShift(startingFloat) {
      this.loading = true
      this.error = null
      
      try {
        const shift = await cashierShiftService.startCashierShift(startingFloat)
        this.currentCashierShift = shift
        return shift
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to start cashier shift'
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Fetch the current open cashier shift
     * @returns {Promise<Object|null>} The current shift or null
     */
    async fetchCurrentCashierShift() {
      this.loading = true
      this.error = null
      
      try {
        const shift = await cashierShiftService.getCurrentCashierShift()
        this.currentCashierShift = shift
        return shift
      } catch (error) {
        this.currentCashierShift = null
        
        // Don't set error if no shift found (404 is expected when no shift is open)
        // Only set error for other types of failures
        if (error.response?.status !== 404) {
          this.error = error.response?.data?.error || 'Failed to fetch current cashier shift'
          console.error('Error fetching current cashier shift:', error)
        }
        
        // Don't throw error for 404, just return null
        if (error.response?.status === 404) {
          return null
        }
        
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Fetch all cashier shifts (manager only)
     * @returns {Promise<Array>} List of all cashier shifts
     */
    async fetchAllCashierShifts() {
      this.loading = true
      this.error = null
      
      try {
        const shifts = await cashierShiftService.getAllCashierShifts()
        this.cashierShifts = shifts
        this.shifts = shifts // Update alias
        return shifts
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to fetch cashier shifts'
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Alias for fetchAllCashierShifts for consistency
     * @returns {Promise<Array>} List of all cashier shifts
     */
    async fetchAllShifts() {
      return this.fetchAllCashierShifts()
    },

    /**
     * Fetch all cashier shifts for the authenticated cashier
     * @returns {Promise<Array>} List of cashier shifts
     */
    async fetchMyCashierShifts() {
      this.loading = true
      this.error = null
      
      try {
        const shifts = await cashierShiftService.getMyCashierShifts()
        this.cashierShifts = shifts
        return shifts
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to fetch my cashier shifts'
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Fetch a specific cashier shift by ID
     * @param {string} shiftId - The cashier shift ID
     * @returns {Promise<Object>} The cashier shift
     */
    async fetchCashierShift(shiftId) {
      this.loading = true
      this.error = null
      
      try {
        const shift = await cashierShiftService.getCashierShift(shiftId)
        
        // Update current shift if it matches
        if (this.currentCashierShift?.id === shiftId) {
          this.currentCashierShift = shift
        }
        
        return shift
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to fetch cashier shift'
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Clear the current cashier shift
     * Useful after closing a shift
     */
    clearCurrentShift() {
      this.currentCashierShift = null
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
      this.currentCashierShift = null
      this.cashierShifts = []
      this.shifts = []
      this.loading = false
      this.error = null
    }
  }
})
