/**
 * Environment variable utilities
 * Provides safe access to Vite environment variables with fallbacks
 */

/**
 * Get API URL from environment
 * @returns {string}
 */
export function getApiUrl() {
  return import.meta.env.VITE_API_URL || 'http://localhost:3000'
}

/**
 * Get Print Bridge URL from environment
 * @returns {string}
 */
export function getPrintBridgeUrl() {
  return import.meta.env.VITE_PRINT_BRIDGE_URL || 'http://localhost:3001'
}

/**
 * Check if running in development mode
 * @returns {boolean}
 */
export function isDevelopment() {
  return import.meta.env.DEV
}

/**
 * Check if running in production mode
 * @returns {boolean}
 */
export function isProduction() {
  return import.meta.env.PROD
}
