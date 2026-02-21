import api from './api'

/**
 * Shop Settings API Service
 */

/**
 * Fetch shop settings
 * @returns {Promise<Object>} Shop settings
 */
export async function fetchShopSettings() {
  const response = await api.get('/shop-settings')
  return response.data
}

/**
 * Update shop settings
 * @param {string} id - Settings ID
 * @param {Object} settings - Settings data
 * @returns {Promise<Object>} Updated settings
 */
export async function updateShopSettings(id, settings) {
  const response = await api.put(`/shop-settings/${id}`, settings)
  return response.data
}
