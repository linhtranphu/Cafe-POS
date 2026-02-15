/**
 * Batch Error Handling Composable
 * Provides standardized error messages and error handling utilities
 */

/**
 * Error message catalog for batch operations
 */
export const ERROR_MESSAGES = {
  // Network errors
  NETWORK_ERROR: 'Không thể kết nối đến máy chủ. Vui lòng kiểm tra kết nối mạng.',
  TIMEOUT_ERROR: 'Yêu cầu hết thời gian chờ. Vui lòng thử lại.',
  
  // Batch Definition errors
  DEFINITION_FETCH_ERROR: 'Không thể tải danh sách định nghĩa batch',
  DEFINITION_CREATE_ERROR: 'Không thể tạo định nghĩa batch',
  DEFINITION_UPDATE_ERROR: 'Không thể cập nhật định nghĩa batch',
  DEFINITION_DELETE_ERROR: 'Không thể xóa định nghĩa batch',
  DEFINITION_NOT_FOUND: 'Không tìm thấy định nghĩa batch',
  
  // Batch Record errors
  RECORD_FETCH_ERROR: 'Không thể tải danh sách batch record',
  RECORD_CREATE_ERROR: 'Không thể ghi nhận batch',
  RECORD_UPDATE_ERROR: 'Không thể cập nhật batch record',
  RECORD_DELETE_ERROR: 'Không thể xóa batch record',
  RECORD_NOT_FOUND: 'Không tìm thấy batch record',
  RECORD_MARK_EXPIRED_ERROR: 'Không thể đánh dấu hết hạn',
  
  // Insufficient resources
  INSUFFICIENT_INGREDIENTS: 'Không đủ nguyên liệu để chế biến batch',
  INSUFFICIENT_BATCH: 'Không đủ batch khả dụng',
  
  // Validation errors
  INVALID_QUANTITY: 'Số lượng không hợp lệ',
  INVALID_DATE: 'Ngày tháng không hợp lệ',
  INVALID_INPUT: 'Dữ liệu nhập vào không hợp lệ',
  REQUIRED_FIELD: 'Vui lòng điền đầy đủ thông tin bắt buộc',
  
  // Alert errors
  ALERT_FETCH_ERROR: 'Không thể tải cảnh báo',
  
  // Report errors
  REPORT_FETCH_ERROR: 'Không thể tải báo cáo',
  REPORT_EXPORT_ERROR: 'Không thể xuất báo cáo',
  
  // Permission errors
  PERMISSION_DENIED: 'Bạn không có quyền thực hiện thao tác này',
  UNAUTHORIZED: 'Vui lòng đăng nhập để tiếp tục',
  
  // Generic errors
  UNKNOWN_ERROR: 'Đã xảy ra lỗi không xác định',
  SERVER_ERROR: 'Lỗi máy chủ. Vui lòng thử lại sau.',
}

/**
 * Error types for categorization
 */
export const ERROR_TYPES = {
  NETWORK: 'network',
  VALIDATION: 'validation',
  PERMISSION: 'permission',
  NOT_FOUND: 'not_found',
  SERVER: 'server',
  UNKNOWN: 'unknown'
}

/**
 * Parse error from API response
 * @param {Error} error - Error object from API call
 * @returns {Object} Parsed error with type and message
 */
export function parseError(error) {
  // Network errors
  if (!error.response) {
    return {
      type: ERROR_TYPES.NETWORK,
      message: ERROR_MESSAGES.NETWORK_ERROR,
      originalError: error
    }
  }

  const status = error.response?.status
  const data = error.response?.data

  // Permission errors
  if (status === 401) {
    return {
      type: ERROR_TYPES.PERMISSION,
      message: ERROR_MESSAGES.UNAUTHORIZED,
      originalError: error
    }
  }

  if (status === 403) {
    return {
      type: ERROR_TYPES.PERMISSION,
      message: ERROR_MESSAGES.PERMISSION_DENIED,
      originalError: error
    }
  }

  // Not found errors
  if (status === 404) {
    return {
      type: ERROR_TYPES.NOT_FOUND,
      message: data?.message || ERROR_MESSAGES.RECORD_NOT_FOUND,
      originalError: error
    }
  }

  // Validation errors
  if (status === 400) {
    return {
      type: ERROR_TYPES.VALIDATION,
      message: data?.message || data?.error || ERROR_MESSAGES.INVALID_INPUT,
      details: data?.details,
      originalError: error
    }
  }

  // Server errors
  if (status >= 500) {
    return {
      type: ERROR_TYPES.SERVER,
      message: ERROR_MESSAGES.SERVER_ERROR,
      originalError: error
    }
  }

  // Unknown errors
  return {
    type: ERROR_TYPES.UNKNOWN,
    message: data?.message || data?.error || ERROR_MESSAGES.UNKNOWN_ERROR,
    originalError: error
  }
}

/**
 * Get user-friendly error message
 * @param {Error} error - Error object
 * @param {string} defaultMessage - Default message if parsing fails
 * @returns {string} User-friendly error message
 */
export function getErrorMessage(error, defaultMessage = ERROR_MESSAGES.UNKNOWN_ERROR) {
  if (typeof error === 'string') {
    return error
  }

  const parsed = parseError(error)
  return parsed.message || defaultMessage
}

/**
 * Check if error is retryable
 * @param {Error} error - Error object
 * @returns {boolean} True if error is retryable
 */
export function isRetryableError(error) {
  const parsed = parseError(error)
  return [ERROR_TYPES.NETWORK, ERROR_TYPES.SERVER].includes(parsed.type)
}

/**
 * Get error icon based on error type
 * @param {string} errorType - Error type from ERROR_TYPES
 * @returns {string} Emoji icon for error type
 */
export function getErrorIcon(errorType) {
  const icons = {
    [ERROR_TYPES.NETWORK]: '📡',
    [ERROR_TYPES.VALIDATION]: '⚠️',
    [ERROR_TYPES.PERMISSION]: '🔒',
    [ERROR_TYPES.NOT_FOUND]: '🔍',
    [ERROR_TYPES.SERVER]: '🔧',
    [ERROR_TYPES.UNKNOWN]: '❌'
  }
  return icons[errorType] || '❌'
}

/**
 * Composable for batch error handling
 */
export function useBatchErrors() {
  /**
   * Handle error and return formatted error object
   * @param {Error} error - Error to handle
   * @param {string} context - Context where error occurred
   * @returns {Object} Formatted error object
   */
  const handleError = (error, context = '') => {
    const parsed = parseError(error)
    
    // Log error for debugging
    console.error(`[Batch Error${context ? ` - ${context}` : ''}]:`, {
      type: parsed.type,
      message: parsed.message,
      details: parsed.details,
      originalError: parsed.originalError
    })

    return {
      ...parsed,
      context,
      icon: getErrorIcon(parsed.type),
      isRetryable: isRetryableError(error)
    }
  }

  /**
   * Create error handler for specific operation
   * @param {string} operation - Operation name
   * @returns {Function} Error handler function
   */
  const createErrorHandler = (operation) => {
    return (error) => handleError(error, operation)
  }

  return {
    ERROR_MESSAGES,
    ERROR_TYPES,
    parseError,
    getErrorMessage,
    isRetryableError,
    getErrorIcon,
    handleError,
    createErrorHandler
  }
}
