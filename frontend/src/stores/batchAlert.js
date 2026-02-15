import { defineStore } from 'pinia'
import batchAlertService from '../services/batchAlert'

/**
 * Batch Alert Store
 * Manages batch alert state and auto-refresh functionality
 */
export const useBatchAlertStore = defineStore('batchAlert', {
  state: () => ({
    // Alert data
    alerts: {
      low_stock: [],
      expiring: [],
      expired: [],
      last_checked: null
    },
    
    // Loading state
    loading: false,
    
    // Error state
    error: null,
    
    // Auto-refresh interval ID
    refreshIntervalId: null,
    
    // Auto-refresh enabled flag
    autoRefreshEnabled: false
  }),

  getters: {
    /**
     * Get total count of all alerts
     */
    totalAlertCount: (state) => {
      return (
        state.alerts.low_stock.length +
        state.alerts.expiring.length +
        state.alerts.expired.length
      )
    },

    /**
     * Get low stock alert count
     */
    lowStockCount: (state) => {
      return state.alerts.low_stock.length
    },

    /**
     * Get expiring alert count
     */
    expiringCount: (state) => {
      return state.alerts.expiring.length
    },

    /**
     * Get expired alert count
     */
    expiredCount: (state) => {
      return state.alerts.expired.length
    },

    /**
     * Check if there are any alerts
     */
    hasAlerts: (state) => {
      return (
        state.alerts.low_stock.length > 0 ||
        state.alerts.expiring.length > 0 ||
        state.alerts.expired.length > 0
      )
    },

    /**
     * Check if currently loading
     */
    isLoading: (state) => {
      return state.loading
    },

    /**
     * Get current error message
     */
    errorMessage: (state) => {
      return state.error
    },

    /**
     * Get last checked timestamp
     */
    lastChecked: (state) => {
      return state.alerts.last_checked
    },

    /**
     * Check if auto-refresh is active
     */
    isAutoRefreshActive: (state) => {
      return state.autoRefreshEnabled && state.refreshIntervalId !== null
    }
  },

  actions: {
    /**
     * Fetch alerts from API
     */
    async fetchAlerts() {
      this.loading = true
      this.error = null

      try {
        const data = await batchAlertService.fetchAlerts()
        
        this.alerts = {
          low_stock: data.low_stock || [],
          expiring: data.expiring || [],
          expired: data.expired || [],
          last_checked: data.last_checked || new Date().toISOString()
        }
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Start auto-refresh with specified interval
     * @param {number} interval - Refresh interval in milliseconds (default: 5 minutes)
     */
    startAutoRefresh(interval = 300000) {
      // Stop existing refresh if any
      this.stopAutoRefresh()

      // Start new refresh interval
      this.refreshIntervalId = batchAlertService.startPolling(
        (data) => {
          this.alerts = {
            low_stock: data.low_stock || [],
            expiring: data.expiring || [],
            expired: data.expired || [],
            last_checked: data.last_checked || new Date().toISOString()
          }
        },
        interval
      )

      this.autoRefreshEnabled = true
    },

    /**
     * Stop auto-refresh
     */
    stopAutoRefresh() {
      if (this.refreshIntervalId) {
        batchAlertService.stopPolling(this.refreshIntervalId)
        this.refreshIntervalId = null
      }
      this.autoRefreshEnabled = false
    },

    /**
     * Get alerts grouped by type
     * @returns {Object} Alerts grouped by type with counts
     */
    getGroupedAlerts() {
      return {
        low_stock: {
          alerts: this.alerts.low_stock,
          count: this.alerts.low_stock.length
        },
        expiring: {
          alerts: this.alerts.expiring,
          count: this.alerts.expiring.length
        },
        expired: {
          alerts: this.alerts.expired,
          count: this.alerts.expired.length
        }
      }
    },

    /**
     * Get badge counts for each alert type
     * @returns {Object} Badge counts
     */
    getBadgeCounts() {
      return {
        low_stock: this.alerts.low_stock.length,
        expiring: this.alerts.expiring.length,
        expired: this.alerts.expired.length,
        total: this.totalAlertCount
      }
    },

    /**
     * Clear error state
     */
    clearError() {
      this.error = null
    },

    /**
     * Reset store to initial state
     */
    reset() {
      this.stopAutoRefresh()
      this.alerts = {
        low_stock: [],
        expiring: [],
        expired: [],
        last_checked: null
      }
      this.loading = false
      this.error = null
      this.refreshIntervalId = null
      this.autoRefreshEnabled = false
    }
  }
})
