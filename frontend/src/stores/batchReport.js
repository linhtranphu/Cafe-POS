import { defineStore } from 'pinia'
import batchReportService from '../services/batchReport'

/**
 * Batch Report Store
 * Manages batch report state and operations
 */
export const useBatchReportStore = defineStore('batchReport', {
  state: () => ({
    // Report data
    productionReport: null,
    wastageReport: null,
    usageReport: null,
    
    // Loading states
    loading: false,
    
    // Error state
    error: null,
    
    // Current report parameters
    reportParams: {
      from_date: null,
      to_date: null,
      batch_definition_id: null,
      prepared_by: null,
      menu_item_id: null
    }
  }),

  getters: {
    /**
     * Check if production report is loaded
     */
    hasProductionReport: (state) => {
      return state.productionReport !== null
    },

    /**
     * Check if wastage report is loaded
     */
    hasWastageReport: (state) => {
      return state.wastageReport !== null
    },

    /**
     * Check if usage report is loaded
     */
    hasUsageReport: (state) => {
      return state.usageReport !== null
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
     * Get total batches produced from production report
     */
    totalBatchesProduced: (state) => {
      return state.productionReport?.total_batches_produced || 0
    },

    /**
     * Get total quantity produced from production report
     */
    totalQuantityProduced: (state) => {
      return state.productionReport?.total_quantity_produced || 0
    },

    /**
     * Get total production cost from production report
     */
    totalProductionCost: (state) => {
      return state.productionReport?.total_cost || 0
    },

    /**
     * Get total expired batches from wastage report
     */
    totalExpiredBatches: (state) => {
      return state.wastageReport?.total_expired_batches || 0
    },

    /**
     * Get total quantity wasted from wastage report
     */
    totalQuantityWasted: (state) => {
      return state.wastageReport?.total_quantity_wasted || 0
    },

    /**
     * Get total cost wasted from wastage report
     */
    totalCostWasted: (state) => {
      return state.wastageReport?.total_cost_wasted || 0
    },

    /**
     * Get total usage count from usage report
     */
    totalUsageCount: (state) => {
      return state.usageReport?.total_usage_count || 0
    },

    /**
     * Get total quantity used from usage report
     */
    totalQuantityUsed: (state) => {
      return state.usageReport?.total_quantity_used || 0
    },

    /**
     * Get total usage cost from usage report
     */
    totalUsageCost: (state) => {
      return state.usageReport?.total_cost || 0
    }
  },

  actions: {
    /**
     * Fetch production report
     * @param {Object} params - Report parameters
     * @param {string} params.from_date - Start date (ISO format)
     * @param {string} params.to_date - End date (ISO format)
     * @param {string} params.batch_definition_id - Optional: Filter by batch definition ID
     * @param {string} params.prepared_by - Optional: Filter by preparer
     */
    async fetchProductionReport(params) {
      this.loading = true
      this.error = null

      try {
        // Validate required parameters
        if (!params.from_date || !params.to_date) {
          throw new Error('Vui lòng chọn khoảng thời gian')
        }

        this.reportParams = {
          ...this.reportParams,
          from_date: params.from_date,
          to_date: params.to_date,
          batch_definition_id: params.batch_definition_id || null,
          prepared_by: params.prepared_by || null
        }

        const report = await batchReportService.fetchProductionReport(params)
        this.productionReport = report
        
        return report
      } catch (error) {
        this.error = error.message
        this.productionReport = null
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Fetch wastage report
     * @param {Object} params - Report parameters
     * @param {string} params.from_date - Start date (ISO format)
     * @param {string} params.to_date - End date (ISO format)
     * @param {string} params.batch_definition_id - Optional: Filter by batch definition ID
     */
    async fetchWastageReport(params) {
      this.loading = true
      this.error = null

      try {
        // Validate required parameters
        if (!params.from_date || !params.to_date) {
          throw new Error('Vui lòng chọn khoảng thời gian')
        }

        this.reportParams = {
          ...this.reportParams,
          from_date: params.from_date,
          to_date: params.to_date,
          batch_definition_id: params.batch_definition_id || null
        }

        const report = await batchReportService.fetchWastageReport(params)
        this.wastageReport = report
        
        return report
      } catch (error) {
        this.error = error.message
        this.wastageReport = null
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Fetch usage report
     * @param {Object} params - Report parameters
     * @param {string} params.from_date - Start date (ISO format)
     * @param {string} params.to_date - End date (ISO format)
     * @param {string} params.batch_definition_id - Optional: Filter by batch definition ID
     * @param {string} params.menu_item_id - Optional: Filter by menu item ID
     */
    async fetchUsageReport(params) {
      this.loading = true
      this.error = null

      try {
        // Validate required parameters
        if (!params.from_date || !params.to_date) {
          throw new Error('Vui lòng chọn khoảng thời gian')
        }

        this.reportParams = {
          ...this.reportParams,
          from_date: params.from_date,
          to_date: params.to_date,
          batch_definition_id: params.batch_definition_id || null,
          menu_item_id: params.menu_item_id || null
        }

        const report = await batchReportService.fetchUsageReport(params)
        this.usageReport = report
        
        return report
      } catch (error) {
        this.error = error.message
        this.usageReport = null
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Export report to CSV
     * @param {string} reportType - Type of report ('production', 'wastage', 'usage')
     * @param {Object} params - Optional: Report parameters (uses current reportParams if not provided)
     */
    async exportReport(reportType, params = null) {
      this.loading = true
      this.error = null

      try {
        // Use provided params or current reportParams
        const exportParams = params || this.reportParams

        // Validate required parameters
        if (!exportParams.from_date || !exportParams.to_date) {
          throw new Error('Vui lòng chọn khoảng thời gian trước khi xuất báo cáo')
        }

        // Validate report type
        if (!['production', 'wastage', 'usage'].includes(reportType)) {
          throw new Error('Loại báo cáo không hợp lệ')
        }

        const blob = await batchReportService.exportReport(reportType, exportParams)
        
        // Create download link
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        
        // Generate filename with date range
        const fromDate = exportParams.from_date.split('T')[0]
        const toDate = exportParams.to_date.split('T')[0]
        const reportTypeVi = {
          production: 'san-xuat',
          wastage: 'lang-phi',
          usage: 'su-dung'
        }[reportType]
        
        link.download = `bao-cao-${reportTypeVi}-${fromDate}-${toDate}.csv`
        
        // Trigger download
        document.body.appendChild(link)
        link.click()
        
        // Cleanup
        document.body.removeChild(link)
        window.URL.revokeObjectURL(url)
        
        return true
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Set report parameters
     * @param {Object} params - Report parameters
     */
    setReportParams(params) {
      this.reportParams = {
        ...this.reportParams,
        ...params
      }
    },

    /**
     * Clear report parameters
     */
    clearReportParams() {
      this.reportParams = {
        from_date: null,
        to_date: null,
        batch_definition_id: null,
        prepared_by: null,
        menu_item_id: null
      }
    },

    /**
     * Clear all reports
     */
    clearReports() {
      this.productionReport = null
      this.wastageReport = null
      this.usageReport = null
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
      this.productionReport = null
      this.wastageReport = null
      this.usageReport = null
      this.loading = false
      this.error = null
      this.reportParams = {
        from_date: null,
        to_date: null,
        batch_definition_id: null,
        prepared_by: null,
        menu_item_id: null
      }
    }
  }
})
