import { defineStore } from 'pinia'
import batchRecordService from '../services/batchRecord'

/**
 * Batch Record Store
 * Manages batch record state and operations
 */
export const useBatchRecordStore = defineStore('batchRecord', {
  state: () => ({
    // Batch records list
    records: [],
    
    // Current selected record
    currentRecord: null,
    
    // Loading states
    loading: false,
    
    // Error state
    error: null,
    
    // Filters
    filters: {
      batch_definition_id: null,
      status: null,
      prepared_by: null,
      from_date: null,
      to_date: null
    },
    
    // Pagination
    pagination: {
      page: 1,
      limit: 20,
      total: 0
    }
  }),

  getters: {
    /**
     * Get available batch records (not expired or depleted)
     */
    availableRecords: (state) => {
      return state.records.filter(record => record.status === 'available')
    },

    /**
     * Get expired batch records
     */
    expiredRecords: (state) => {
      return state.records.filter(record => record.status === 'expired')
    },

    /**
     * Get depleted batch records
     */
    depletedRecords: (state) => {
      return state.records.filter(record => record.status === 'depleted')
    },

    /**
     * Check if there are any records
     */
    hasRecords: (state) => {
      return state.records.length > 0
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
    }
  },

  actions: {
    /**
     * Fetch batch records with current filters and pagination
     */
    async fetchRecords(customFilters = null) {
      this.loading = true
      this.error = null

      try {
        const filters = customFilters || {
          ...this.filters,
          page: this.pagination.page,
          limit: this.pagination.limit
        }

        const response = await batchRecordService.fetchRecords(filters)
        
        this.records = response.data || []
        this.pagination.total = response.total || 0
        this.pagination.page = parseInt(response.page) || 1
        this.pagination.limit = parseInt(response.limit) || 20
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Fetch a single batch record by ID
     * @param {string} id - Batch record ID
     */
    async fetchRecordById(id) {
      this.loading = true
      this.error = null

      try {
        const record = await batchRecordService.fetchRecordById(id)
        this.currentRecord = record
        
        // Update in records list if exists
        const index = this.records.findIndex(r => r.id === id)
        if (index !== -1) {
          this.records[index] = record
        }
        
        return record
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Create a new batch record
     * @param {Object} data - Batch record data
     */
    async createRecord(data) {
      this.loading = true
      this.error = null

      try {
        const record = await batchRecordService.createRecord(data)
        
        // Add to beginning of records list
        this.records.unshift(record)
        this.pagination.total += 1
        
        return record
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Update batch record quantity
     * @param {string} id - Batch record ID
     * @param {number} quantity - New quantity remaining
     */
    async updateQuantity(id, quantity) {
      this.loading = true
      this.error = null

      try {
        const updatedRecord = await batchRecordService.updateQuantity(id, quantity)
        
        // Update in records list
        const index = this.records.findIndex(r => r.id === id)
        if (index !== -1) {
          this.records[index] = updatedRecord
        }
        
        // Update current record if it's the same
        if (this.currentRecord?.id === id) {
          this.currentRecord = updatedRecord
        }
        
        return updatedRecord
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Mark batch record as expired
     * @param {string} id - Batch record ID
     */
    async markAsExpired(id) {
      this.loading = true
      this.error = null

      try {
        const updatedRecord = await batchRecordService.markAsExpired(id)
        
        // Update in records list
        const index = this.records.findIndex(r => r.id === id)
        if (index !== -1) {
          this.records[index] = updatedRecord
        }
        
        // Update current record if it's the same
        if (this.currentRecord?.id === id) {
          this.currentRecord = updatedRecord
        }
        
        return updatedRecord
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Delete a batch record
     * @param {string} id - Batch record ID
     */
    async deleteRecord(id) {
      this.loading = true
      this.error = null

      try {
        await batchRecordService.deleteRecord(id)
        
        // Remove from records list
        const index = this.records.findIndex(r => r.id === id)
        if (index !== -1) {
          this.records.splice(index, 1)
          this.pagination.total -= 1
        }
        
        // Clear current record if it's the same
        if (this.currentRecord?.id === id) {
          this.currentRecord = null
        }
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * Set filters and reset pagination
     * @param {Object} newFilters - New filter values
     */
    setFilters(newFilters) {
      this.filters = {
        ...this.filters,
        ...newFilters
      }
      this.pagination.page = 1 // Reset to first page when filters change
    },

    /**
     * Clear all filters
     */
    clearFilters() {
      this.filters = {
        batch_definition_id: null,
        status: null,
        prepared_by: null,
        from_date: null,
        to_date: null
      }
      this.pagination.page = 1
    },

    /**
     * Set pagination page
     * @param {number} page - Page number
     */
    setPage(page) {
      this.pagination.page = page
    },

    /**
     * Set pagination limit
     * @param {number} limit - Items per page
     */
    setLimit(limit) {
      this.pagination.limit = limit
      this.pagination.page = 1 // Reset to first page when limit changes
    },

    /**
     * Clear error state
     */
    clearError() {
      this.error = null
    },

    /**
     * Clear current record
     */
    clearCurrentRecord() {
      this.currentRecord = null
    },

    /**
     * Reset store to initial state
     */
    reset() {
      this.records = []
      this.currentRecord = null
      this.loading = false
      this.error = null
      this.filters = {
        batch_definition_id: null,
        status: null,
        prepared_by: null,
        from_date: null,
        to_date: null
      }
      this.pagination = {
        page: 1,
        limit: 20,
        total: 0
      }
    }
  }
})
