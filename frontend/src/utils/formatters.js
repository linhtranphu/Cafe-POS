/**
 * Utility functions for data formatting and transformation
 * Centralized to ensure consistency across the application
 */

/**
 * Convert a date string or Date object to ISO format with timezone
 * Backend expects: "2024-01-15T00:00:00Z"
 * 
 * @param {string|Date|null} date - Date to convert
 * @param {boolean} includeTime - Whether to include time (default: false, sets to 00:00:00)
 * @returns {string|null} ISO formatted date string or null
 */
export function toISODate(date, includeTime = false) {
  if (!date) return null
  
  try {
    let dateObj
    
    // Handle string input (e.g., "2026-01-31" from date input)
    if (typeof date === 'string') {
      // If it's already in ISO format, return as is
      if (date.includes('T') && date.includes('Z')) {
        return date
      }
      // Convert date-only string to ISO format
      dateObj = new Date(date + 'T00:00:00Z')
    } 
    // Handle Date object
    else if (date instanceof Date) {
      dateObj = date
    } 
    // Invalid input
    else {
      return null
    }
    
    // Check if date is valid
    if (isNaN(dateObj.getTime())) {
      return null
    }
    
    // Return ISO string
    if (includeTime) {
      return dateObj.toISOString()
    } else {
      // Set time to 00:00:00 UTC
      const year = dateObj.getUTCFullYear()
      const month = String(dateObj.getUTCMonth() + 1).padStart(2, '0')
      const day = String(dateObj.getUTCDate()).padStart(2, '0')
      return `${year}-${month}-${day}T00:00:00Z`
    }
  } catch (error) {
    console.error('Error converting date to ISO format:', error)
    return null
  }
}

/**
 * Convert ISO date string to local date string for date input
 * Converts: "2024-01-15T00:00:00Z" -> "2024-01-15"
 * 
 * @param {string|null} isoDate - ISO date string
 * @returns {string} Date string in YYYY-MM-DD format
 */
export function fromISODate(isoDate) {
  if (!isoDate) return ''
  
  try {
    const date = new Date(isoDate)
    if (isNaN(date.getTime())) return ''
    
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    
    return `${year}-${month}-${day}`
  } catch (error) {
    console.error('Error converting ISO date to local format:', error)
    return ''
  }
}

/**
 * Format date for display in Vietnamese locale
 * 
 * @param {string|Date|null} date - Date to format
 * @param {object} options - Intl.DateTimeFormat options
 * @returns {string} Formatted date string
 */
export function formatDate(date, options = {}) {
  if (!date) return 'N/A'
  
  try {
    const dateObj = typeof date === 'string' ? new Date(date) : date
    if (isNaN(dateObj.getTime())) return 'N/A'
    
    const defaultOptions = {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      ...options
    }
    
    return new Intl.DateTimeFormat('vi-VN', defaultOptions).format(dateObj)
  } catch (error) {
    console.error('Error formatting date:', error)
    return 'N/A'
  }
}

/**
 * Format date and time for display in Vietnamese locale
 * 
 * @param {string|Date|null} date - Date to format
 * @returns {string} Formatted date and time string
 */
export function formatDateTime(date) {
  return formatDate(date, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

/**
 * Format price in Vietnamese currency
 * 
 * @param {number|null} price - Price to format
 * @param {boolean} showSymbol - Whether to show currency symbol (default: true)
 * @returns {string} Formatted price string
 */
export function formatPrice(price, showSymbol = true) {
  if (price === null || price === undefined) return showSymbol ? '0 ₫' : '0'
  
  try {
    const formatted = new Intl.NumberFormat('vi-VN', {
      style: showSymbol ? 'currency' : 'decimal',
      currency: 'VND',
      maximumFractionDigits: 0
    }).format(price)
    
    return formatted
  } catch (error) {
    console.error('Error formatting price:', error)
    return showSymbol ? '0 ₫' : '0'
  }
}

/**
 * Format number with thousand separators
 * 
 * @param {number|null} num - Number to format
 * @returns {string} Formatted number string
 */
export function formatNumber(num) {
  if (num === null || num === undefined) return '0'
  
  try {
    return new Intl.NumberFormat('vi-VN').format(num)
  } catch (error) {
    console.error('Error formatting number:', error)
    return '0'
  }
}

/**
 * Sanitize form data before sending to backend
 * - Converts dates to ISO format
 * - Ensures numbers are valid
 * - Removes empty strings (converts to null or default values)
 * 
 * @param {object} data - Form data to sanitize
 * @param {object} schema - Schema defining field types and defaults
 * @returns {object} Sanitized data
 */
export function sanitizeFormData(data, schema = {}) {
  const sanitized = {}
  
  for (const [key, value] of Object.entries(data)) {
    const fieldSchema = schema[key] || {}
    const fieldType = fieldSchema.type || 'string'
    const defaultValue = fieldSchema.default
    
    // Handle different field types
    switch (fieldType) {
      case 'date':
        sanitized[key] = value ? toISODate(value) : (defaultValue || null)
        break
        
      case 'datetime':
        sanitized[key] = value ? toISODate(value, true) : (defaultValue || null)
        break
        
      case 'number':
        sanitized[key] = value !== '' && value !== null && value !== undefined 
          ? Number(value) 
          : (defaultValue !== undefined ? defaultValue : 0)
        break
        
      case 'string':
        sanitized[key] = value !== null && value !== undefined 
          ? String(value) 
          : (defaultValue !== undefined ? defaultValue : '')
        break
        
      case 'boolean':
        sanitized[key] = Boolean(value)
        break
        
      default:
        sanitized[key] = value !== undefined ? value : defaultValue
    }
  }
  
  return sanitized
}

/**
 * Parse backend data for frontend display
 * - Converts ISO dates to local format for date inputs
 * 
 * @param {object} data - Backend data to parse
 * @param {object} schema - Schema defining field types
 * @returns {object} Parsed data
 */
export function parseBackendData(data, schema = {}) {
  const parsed = { ...data }
  
  for (const [key, value] of Object.entries(data)) {
    const fieldSchema = schema[key] || {}
    const fieldType = fieldSchema.type || 'string'
    
    // Handle different field types
    switch (fieldType) {
      case 'date':
      case 'datetime':
        parsed[key] = value ? fromISODate(value) : ''
        break
        
      default:
        parsed[key] = value
    }
  }
  
  return parsed
}

/**
 * Validate required fields
 * 
 * @param {object} data - Data to validate
 * @param {array} requiredFields - Array of required field names
 * @returns {object} { valid: boolean, errors: array }
 */
export function validateRequired(data, requiredFields = []) {
  const errors = []
  
  for (const field of requiredFields) {
    const value = data[field]
    if (value === null || value === undefined || value === '') {
      errors.push(`Trường "${field}" là bắt buộc`)
    }
  }
  
  return {
    valid: errors.length === 0,
    errors
  }
}

/**
 * Deep clone an object
 * 
 * @param {object} obj - Object to clone
 * @returns {object} Cloned object
 */
export function deepClone(obj) {
  if (obj === null || typeof obj !== 'object') return obj
  
  try {
    return JSON.parse(JSON.stringify(obj))
  } catch (error) {
    console.error('Error cloning object:', error)
    return obj
  }
}
