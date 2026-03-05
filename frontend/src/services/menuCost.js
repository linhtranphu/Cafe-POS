import api from './api.js'

/**
 * @typedef {import('./types/menuCost').MenuCostsResponse} MenuCostsResponse
 * @typedef {import('./types/menuCost').MenuItemCostBreakdown} MenuItemCostBreakdown
 * @typedef {import('./types/menuCost').ProfitWarnings} ProfitWarnings
 * @typedef {import('./types/menuCost').ProfitFilter} ProfitFilter
 */

/**
 * Menu Cost API Service
 * Provides methods for retrieving menu item cost and profit information
 */
export const menuCostService = {
  /**
   * Get menu costs with optional filtering and sorting
   * @param {ProfitFilter} [filter={}] - Filter options
   * @returns {Promise<MenuCostsResponse>} Menu costs response
   */
  async getMenuCosts(filter = {}) {
    const params = new URLSearchParams()
    
    if (filter.category) {
      params.append('category', filter.category)
    }
    if (filter.sort_by) {
      params.append('sort_by', filter.sort_by)
    }
    if (filter.sort_order) {
      params.append('sort_order', filter.sort_order)
    }
    
    const queryString = params.toString()
    const url = queryString ? `/manager/menu/costs?${queryString}` : '/manager/menu/costs'
    
    const response = await api.get(url)
    return response.data
  },

  /**
   * Get detailed cost breakdown for a specific menu item
   * @param {string} id - Menu item ID
   * @returns {Promise<MenuItemCostBreakdown>} Menu item cost breakdown
   */
  async getMenuCostDetail(id) {
    const response = await api.get(`/manager/menu/costs/${id}`)
    return response.data
  },

  /**
   * Get menu items with warnings (loss or low margin)
   * @param {number} [threshold] - Optional custom low margin threshold
   * @returns {Promise<ProfitWarnings>} Profit warnings response
   */
  async getMenuWarnings(threshold) {
    const params = new URLSearchParams()
    
    if (threshold !== undefined && threshold !== null) {
      params.append('threshold', threshold.toString())
    }
    
    const queryString = params.toString()
    const url = queryString ? `/manager/menu/warnings?${queryString}` : '/manager/menu/warnings'
    
    const response = await api.get(url)
    return response.data
  },

  /**
   * Recalculate costs for all menu items
   * Queues background jobs to recalculate costs based on current ingredient/batch prices
   * @returns {Promise<{message: string, total_items: number, queued_count: number, failed_count: number}>}
   */
  async recalculateAllCosts() {
    const response = await api.post('/manager/menu/costs/recalculate-all')
    return response.data
  }
}
