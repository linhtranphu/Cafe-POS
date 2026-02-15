/**
 * Batch Color Coding System
 * 
 * Requirement 9.3: Use colors to distinguish status (green: sufficient, yellow: low, red: expiring)
 * 
 * This composable provides consistent color coding across all batch components:
 * - Batch status (available, expiring, expired, depleted)
 * - Stock levels (sufficient, low, critical)
 * - Alert types (low_stock, expiring, expired)
 */

export function useBatchColors() {
  /**
   * Get color classes for batch record based on status and expiry
   * @param {Object} record - Batch record object
   * @returns {Object} Color classes for different elements
   */
  const getBatchRecordColors = (record) => {
    if (!record) return getDefaultColors()

    // Expired batch - Red
    if (record.status === 'expired') {
      return {
        background: 'bg-red-50',
        border: 'border-red-200',
        text: 'text-red-900',
        badge: 'bg-red-500 text-white',
        icon: '🔴',
        statusText: 'Đã hết hạn'
      }
    }

    // Depleted batch - Gray
    if (record.status === 'depleted' || record.quantity_remaining <= 0) {
      return {
        background: 'bg-gray-100',
        border: 'border-gray-300',
        text: 'text-gray-700',
        badge: 'bg-gray-500 text-white',
        icon: '⚫',
        statusText: 'Đã hết'
      }
    }

    // Check expiry time for available batches
    if (record.expires_at) {
      const now = new Date()
      const expiresAt = new Date(record.expires_at)
      const hoursUntilExpiry = (expiresAt - now) / (1000 * 60 * 60)

      // Expiring soon (within 4 hours) - Yellow/Orange
      if (hoursUntilExpiry <= 4 && hoursUntilExpiry > 0) {
        return {
          background: 'bg-yellow-50',
          border: 'border-yellow-300',
          text: 'text-yellow-900',
          badge: 'bg-yellow-500 text-white',
          icon: '🟡',
          statusText: 'Sắp hết hạn'
        }
      }

      // Expiring within 24 hours - Orange
      if (hoursUntilExpiry <= 24 && hoursUntilExpiry > 4) {
        return {
          background: 'bg-orange-50',
          border: 'border-orange-300',
          text: 'text-orange-900',
          badge: 'bg-orange-500 text-white',
          icon: '🟠',
          statusText: 'Cần chú ý'
        }
      }
    }

    // Available and fresh - Green
    return {
      background: 'bg-green-50',
      border: 'border-green-200',
      text: 'text-green-900',
      badge: 'bg-green-500 text-white',
      icon: '🟢',
      statusText: 'Khả dụng'
    }
  }

  /**
   * Get color classes for stock level
   * @param {number} current - Current stock quantity
   * @param {number} threshold - Low stock threshold
   * @returns {Object} Color classes
   */
  const getStockLevelColors = (current, threshold) => {
    const ratio = threshold > 0 ? (current / threshold) : 1

    // Critical - below threshold (Red)
    if (current <= threshold) {
      return {
        background: 'bg-red-50',
        border: 'border-red-300',
        text: 'text-red-700',
        badge: 'bg-red-500 text-white',
        icon: '⚠️',
        level: 'critical'
      }
    }

    // Low - 1-2x threshold (Yellow)
    if (ratio <= 2) {
      return {
        background: 'bg-yellow-50',
        border: 'border-yellow-300',
        text: 'text-yellow-700',
        badge: 'bg-yellow-500 text-white',
        icon: '⚡',
        level: 'low'
      }
    }

    // Sufficient - above 2x threshold (Green)
    return {
      background: 'bg-green-50',
      border: 'border-green-300',
      text: 'text-green-700',
      badge: 'bg-green-500 text-white',
      icon: '✓',
      level: 'sufficient'
    }
  }

  /**
   * Get color classes for alert type
   * @param {string} type - Alert type (low_stock, expiring, expired)
   * @returns {Object} Color classes
   */
  const getAlertColors = (type) => {
    switch (type) {
      case 'expired':
        return {
          background: 'bg-red-50',
          border: 'border-red-500',
          text: 'text-red-900',
          badge: 'bg-red-500 text-white',
          button: 'bg-red-500 hover:bg-red-600 active:bg-red-600',
          icon: '🔴',
          title: 'Đã Hết Hạn'
        }
      
      case 'expiring':
        return {
          background: 'bg-yellow-50',
          border: 'border-yellow-500',
          text: 'text-yellow-900',
          badge: 'bg-yellow-500 text-white',
          button: 'bg-yellow-500 hover:bg-yellow-600 active:bg-yellow-600',
          icon: '🟡',
          title: 'Sắp Hết Hạn'
        }
      
      case 'low_stock':
        return {
          background: 'bg-orange-50',
          border: 'border-orange-500',
          text: 'text-orange-900',
          badge: 'bg-orange-500 text-white',
          button: 'bg-orange-500 hover:bg-orange-600 active:bg-orange-600',
          icon: '🟠',
          title: 'Tồn Kho Thấp'
        }
      
      default:
        return getDefaultColors()
    }
  }

  /**
   * Get color for expiry time display
   * @param {string} expiresAt - ISO date string
   * @returns {string} Text color class
   */
  const getExpiryTextColor = (expiresAt) => {
    if (!expiresAt) return 'text-gray-700'

    const now = new Date()
    const expires = new Date(expiresAt)
    const hoursUntilExpiry = (expires - now) / (1000 * 60 * 60)

    if (hoursUntilExpiry <= 0) return 'text-red-600 font-bold'
    if (hoursUntilExpiry <= 4) return 'text-yellow-600 font-bold'
    if (hoursUntilExpiry <= 24) return 'text-orange-600'
    return 'text-gray-700'
  }

  /**
   * Get color for quantity display based on percentage
   * @param {number} percentage - Percentage value (0-100)
   * @returns {string} Text color class
   */
  const getQuantityPercentageColor = (percentage) => {
    if (percentage <= 25) return 'text-red-600'
    if (percentage <= 50) return 'text-orange-600'
    if (percentage <= 75) return 'text-yellow-600'
    return 'text-green-600'
  }

  /**
   * Get default color classes
   * @returns {Object} Default color classes
   */
  const getDefaultColors = () => ({
    background: 'bg-gray-50',
    border: 'border-gray-300',
    text: 'text-gray-700',
    badge: 'bg-gray-500 text-white',
    button: 'bg-gray-500 hover:bg-gray-600 active:bg-gray-600',
    icon: '⚪',
    statusText: 'Không xác định'
  })

  /**
   * Get button color classes for actions
   * @param {string} action - Action type (delete, expire, view, create)
   * @returns {string} Button color classes
   */
  const getActionButtonColors = (action) => {
    switch (action) {
      case 'delete':
        return 'bg-red-500 text-white active:bg-red-600 hover:bg-red-600'
      
      case 'expire':
      case 'warning':
        return 'bg-yellow-500 text-white active:bg-yellow-600 hover:bg-yellow-600'
      
      case 'create':
      case 'save':
      case 'confirm':
        return 'bg-green-500 text-white active:bg-green-600 hover:bg-green-600'
      
      case 'view':
      case 'edit':
        return 'bg-blue-500 text-white active:bg-blue-600 hover:bg-blue-600'
      
      case 'cancel':
        return 'bg-gray-500 text-white active:bg-gray-600 hover:bg-gray-600'
      
      default:
        return 'bg-gray-500 text-white active:bg-gray-600 hover:bg-gray-600'
    }
  }

  return {
    getBatchRecordColors,
    getStockLevelColors,
    getAlertColors,
    getExpiryTextColor,
    getQuantityPercentageColor,
    getActionButtonColors,
    getDefaultColors
  }
}
